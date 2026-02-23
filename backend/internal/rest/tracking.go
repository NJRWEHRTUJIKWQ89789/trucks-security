package rest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"cargomax-api/internal/auth"
	"cargomax-api/internal/config"
	"cargomax-api/internal/models"
	"cargomax-api/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TrackingHandler struct {
	Config      *config.Config
	DriverRepo  *repository.DriverRepo
	VehicleRepo *repository.VehicleRepo
	ShiftRepo   *repository.ShiftRepo
	PingRepo    *repository.GPSPingRepo
	AlertRepo   *repository.AlertRepo
	ZoneRepo    *repository.ZoneRepo
	WSHub       *Hub

	// Rate limiting: last ping time per driver
	pingRateMu sync.Mutex
	pingRates  map[uuid.UUID]time.Time
}

func NewTrackingHandler(cfg *config.Config, driverRepo *repository.DriverRepo, vehicleRepo *repository.VehicleRepo, shiftRepo *repository.ShiftRepo, pingRepo *repository.GPSPingRepo, alertRepo *repository.AlertRepo, zoneRepo *repository.ZoneRepo, hub *Hub) *TrackingHandler {
	return &TrackingHandler{
		Config:      cfg,
		DriverRepo:  driverRepo,
		VehicleRepo: vehicleRepo,
		ShiftRepo:   shiftRepo,
		PingRepo:    pingRepo,
		AlertRepo:   alertRepo,
		ZoneRepo:    zoneRepo,
		WSHub:       hub,
		pingRates:   make(map[uuid.UUID]time.Time),
	}
}

func (h *TrackingHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// Public routes (no auth)
	r.Post("/auth/driver-login", h.DriverLogin)

	// Protected routes (driver auth required)
	r.Group(func(r chi.Router) {
		r.Use(h.driverAuthMiddleware)
		r.Get("/driver/trucks", h.GetTrucks)
		r.Post("/shifts/start", h.StartShift)
		r.Post("/shifts/end", h.EndShift)
		r.Post("/tracking/ping", h.ReceivePings)
		r.Post("/auth/refresh-token", h.RefreshToken)
		r.Get("/driver/active-shift", h.GetActiveShift)
	})

	return r
}

// driverAuthMiddleware validates Bearer token for driver endpoints.
func (h *TrackingHandler) driverAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			jsonError(w, "missing authorization header", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			jsonError(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(h.Config.JWTPublicKey, parts[1])
		if err != nil {
			jsonError(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, models.CtxTenantID, claims.TenantID)
		ctx = context.WithValue(ctx, models.CtxUserID, claims.UserID)
		ctx = context.WithValue(ctx, models.CtxUserRole, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// DriverLogin handles POST /api/v1/auth/driver-login
func (h *TrackingHandler) DriverLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
		Pin   string `json:"pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Phone == "" || req.Pin == "" {
		jsonError(w, "phone and pin are required", http.StatusBadRequest)
		return
	}

	// Look up driver by phone number
	driver, err := h.DriverRepo.GetByPhone(r.Context(), req.Phone)
	if err != nil {
		jsonError(w, "invalid phone or PIN", http.StatusUnauthorized)
		return
	}

	// Verify PIN
	if driver.PinHash == nil || *driver.PinHash == "" {
		jsonError(w, "driver account not configured for mobile login", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*driver.PinHash), []byte(req.Pin)); err != nil {
		jsonError(w, "invalid phone or PIN", http.StatusUnauthorized)
		return
	}

	// Build display name
	name := ""
	if driver.FirstName != nil {
		name = *driver.FirstName
	}
	if driver.LastName != nil {
		if name != "" {
			name += " "
		}
		name += *driver.LastName
	}

	email := ""
	if driver.Email != nil {
		email = *driver.Email
	}

	// Generate JWT tokens with role="driver"
	accessToken, err := auth.CreateAccessToken(h.Config.JWTPrivateKey, driver.ID, driver.TenantID, email, "driver")
	if err != nil {
		jsonError(w, "failed to create token", http.StatusInternalServerError)
		return
	}
	refreshToken, err := auth.CreateRefreshToken(h.Config.JWTPrivateKey, driver.ID, driver.TenantID)
	if err != nil {
		jsonError(w, "failed to create refresh token", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"token":         accessToken,
		"refresh_token": refreshToken,
		"driver": map[string]interface{}{
			"id":        driver.ID,
			"name":      name,
			"tenant_id": driver.TenantID,
		},
	})
}

// GetTrucks handles GET /api/v1/driver/trucks
func (h *TrackingHandler) GetTrucks(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)

	trucks, err := h.VehicleRepo.GetActive(r.Context(), tenantID)
	if err != nil {
		jsonError(w, "failed to fetch trucks", http.StatusInternalServerError)
		return
	}

	result := make([]map[string]interface{}, 0, len(trucks))
	for _, t := range trucks {
		plate := ""
		if t.LicensePlate != nil {
			plate = *t.LicensePlate
		}
		name := ""
		if t.Name != nil {
			name = *t.Name
		}
		vType := ""
		if t.Type != nil {
			vType = *t.Type
		}
		result = append(result, map[string]interface{}{
			"id":           t.ID,
			"plate_number": plate,
			"name":         name,
			"type":         vType,
		})
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{"trucks": result})
}

// StartShift handles POST /api/v1/shifts/start
func (h *TrackingHandler) StartShift(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)
	driverID := r.Context().Value(models.CtxUserID).(uuid.UUID)

	var req struct {
		TruckID string `json:"truck_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	truckUUID, err := uuid.Parse(req.TruckID)
	if err != nil {
		jsonError(w, "invalid truck_id", http.StatusBadRequest)
		return
	}

	// Validate truck belongs to the same tenant (prevent IDOR / cross-tenant access).
	if _, err := h.VehicleRepo.GetByID(r.Context(), tenantID, truckUUID); err != nil {
		jsonError(w, "truck not found in tenant", http.StatusNotFound)
		return
	}

	// Check driver doesn't already have an active shift
	existing, _ := h.ShiftRepo.GetActiveByDriver(r.Context(), tenantID, driverID)
	if existing != nil {
		jsonError(w, "driver already has an active shift", http.StatusConflict)
		return
	}

	// Check truck is not in use
	inUse, _ := h.ShiftRepo.IsTruckInUse(r.Context(), tenantID, truckUUID)
	if inUse {
		jsonError(w, "truck is already in use by another driver", http.StatusConflict)
		return
	}

	shift := &models.Shift{
		TenantID: tenantID,
		DriverID: driverID,
		TruckID:  truckUUID,
	}
	if err := h.ShiftRepo.Create(r.Context(), shift); err != nil {
		jsonError(w, "failed to create shift", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"shift_id":   shift.ID,
		"started_at": shift.StartedAt,
	})
}

// EndShift handles POST /api/v1/shifts/end
func (h *TrackingHandler) EndShift(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)
	driverID := r.Context().Value(models.CtxUserID).(uuid.UUID)

	var req struct {
		ShiftID string `json:"shift_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	shiftUUID, err := uuid.Parse(req.ShiftID)
	if err != nil {
		jsonError(w, "invalid shift_id", http.StatusBadRequest)
		return
	}

	// Calculate total km from pings (tenant-scoped)
	totalKm, _ := h.PingRepo.CalculateShiftKm(r.Context(), tenantID, shiftUUID)

	// End the shift with tenant_id + driver_id validation to prevent cross-tenant/cross-driver IDOR
	shift, err := h.ShiftRepo.End(r.Context(), tenantID, driverID, shiftUUID, totalKm)
	if err != nil {
		jsonError(w, "failed to end shift", http.StatusInternalServerError)
		return
	}

	durationMinutes := 0.0
	if shift.EndedAt != nil {
		durationMinutes = shift.EndedAt.Sub(shift.StartedAt).Minutes()
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"ended_at":               shift.EndedAt,
		"total_km":               shift.TotalKm,
		"total_duration_minutes": durationMinutes,
	})
}

// ReceivePings handles POST /api/v1/tracking/ping
func (h *TrackingHandler) ReceivePings(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)
	driverID := r.Context().Value(models.CtxUserID).(uuid.UUID)

	// Rate limiting: max 1 batch per 5 seconds per driver
	h.pingRateMu.Lock()
	lastPing, exists := h.pingRates[driverID]
	if exists && time.Since(lastPing) < 5*time.Second {
		h.pingRateMu.Unlock()
		jsonError(w, "rate limited: max 1 batch per 5 seconds", http.StatusTooManyRequests)
		return
	}
	h.pingRates[driverID] = time.Now()
	h.pingRateMu.Unlock()

	var req struct {
		Pings []struct {
			Latitude     float64   `json:"latitude"`
			Longitude    float64   `json:"longitude"`
			SpeedKmh     float64   `json:"speed_kmh"`
			Heading      int       `json:"heading"`
			Accuracy     float64   `json:"accuracy"`
			BatteryLevel int       `json:"battery_level"`
			IsMoving     bool      `json:"is_moving"`
			RecordedAt   time.Time `json:"recorded_at"`
			TruckID      string    `json:"truck_id"`
			ShiftID      string    `json:"shift_id"`
		} `json:"pings"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Pings) == 0 {
		jsonError(w, "no pings provided", http.StatusBadRequest)
		return
	}
	if len(req.Pings) > 100 {
		jsonError(w, "max 100 pings per batch", http.StatusBadRequest)
		return
	}

	// Validate that the driver has an active shift in this tenant.
	// Use the active shift as the authoritative source for truck_id and shift_id
	// to prevent spoofing of these values from the client payload.
	activeShift, err := h.ShiftRepo.GetActiveByDriver(r.Context(), tenantID, driverID)
	if err != nil || activeShift == nil {
		jsonError(w, "no active shift found for driver", http.StatusForbidden)
		return
	}

	now := time.Now()
	pings := make([]models.GPSPing, 0, len(req.Pings))
	for _, p := range req.Pings {
		// Validate
		if p.Latitude < -90 || p.Latitude > 90 || p.Longitude < -180 || p.Longitude > 180 {
			continue
		}
		if p.SpeedKmh < 0 || p.SpeedKmh > 200 {
			p.SpeedKmh = 0
		}

		isDelayed := now.Sub(p.RecordedAt) > 2*time.Minute

		// Use server-side authoritative shift_id and truck_id from the active shift,
		// ignoring client-supplied values to prevent IDOR / data injection.
		pings = append(pings, models.GPSPing{
			TenantID:     tenantID,
			DriverID:     driverID,
			TruckID:      activeShift.TruckID,
			ShiftID:      activeShift.ID,
			Latitude:     p.Latitude,
			Longitude:    p.Longitude,
			SpeedKmh:     p.SpeedKmh,
			Heading:      p.Heading,
			Accuracy:     p.Accuracy,
			BatteryLevel: p.BatteryLevel,
			IsMoving:     p.IsMoving,
			RecordedAt:   p.RecordedAt,
			ReceivedAt:   now,
			IsDelayed:    isDelayed,
		})
	}

	count, err := h.PingRepo.BulkInsert(r.Context(), pings)
	if err != nil {
		log.Printf("failed to insert pings: %v", err)
		jsonError(w, "failed to store pings", http.StatusInternalServerError)
		return
	}

	// Broadcast to WebSocket hub for live dashboard
	if h.WSHub != nil && len(pings) > 0 {
		lastPing := pings[len(pings)-1]
		// Get driver name for broadcast
		driver, _ := h.DriverRepo.GetByID(r.Context(), tenantID, driverID)
		driverName := ""
		if driver != nil {
			if driver.FirstName != nil {
				driverName = *driver.FirstName
			}
			if driver.LastName != nil {
				driverName += " " + *driver.LastName
			}
		}

		// Get truck plate
		truckPlate := ""
		truck, _ := h.VehicleRepo.GetByID(r.Context(), tenantID, lastPing.TruckID)
		if truck != nil && truck.LicensePlate != nil {
			truckPlate = *truck.LicensePlate
		}

		h.WSHub.BroadcastTracking(tenantID, map[string]interface{}{
			"driver_id":   driverID,
			"driver_name": strings.TrimSpace(driverName),
			"truck_plate": truckPlate,
			"lat":         lastPing.Latitude,
			"lng":         lastPing.Longitude,
			"speed":       lastPing.SpeedKmh,
			"heading":     lastPing.Heading,
			"battery":     lastPing.BatteryLevel,
			"is_moving":   lastPing.IsMoving,
			"timestamp":   lastPing.RecordedAt,
		})
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{"received": count})
}

// RefreshToken handles POST /api/v1/auth/refresh-token
func (h *TrackingHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(models.CtxUserID).(uuid.UUID)
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)

	driver, err := h.DriverRepo.GetByID(r.Context(), tenantID, userID)
	if err != nil {
		jsonError(w, "driver not found", http.StatusUnauthorized)
		return
	}

	email := ""
	if driver.Email != nil {
		email = *driver.Email
	}

	accessToken, err := auth.CreateAccessToken(h.Config.JWTPrivateKey, driver.ID, driver.TenantID, email, "driver")
	if err != nil {
		jsonError(w, "failed to create token", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{"token": accessToken})
}

// GetActiveShift handles GET /api/v1/driver/active-shift
func (h *TrackingHandler) GetActiveShift(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)
	driverID := r.Context().Value(models.CtxUserID).(uuid.UUID)

	shift, err := h.ShiftRepo.GetActiveByDriver(r.Context(), tenantID, driverID)
	if err != nil {
		jsonResponse(w, http.StatusOK, map[string]interface{}{"active_shift": nil})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"active_shift": map[string]interface{}{
			"shift_id":   shift.ID,
			"truck_id":   shift.TruckID,
			"started_at": shift.StartedAt,
		},
	})
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
