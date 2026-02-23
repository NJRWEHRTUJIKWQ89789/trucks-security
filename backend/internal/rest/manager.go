package rest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cargomax-api/internal/auth"
	"cargomax-api/internal/config"
	"cargomax-api/internal/models"
	"cargomax-api/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ManagerHandler serves the manager-facing REST API consumed by the Next.js
// dashboard. Every request is scoped to a tenant via JWT claims.
type ManagerHandler struct {
	Config      *config.Config
	DriverRepo  *repository.DriverRepo
	VehicleRepo *repository.VehicleRepo
	ShiftRepo   *repository.ShiftRepo
	PingRepo    *repository.GPSPingRepo
	AlertRepo   *repository.AlertRepo
	ZoneRepo    *repository.ZoneRepo
}

// NewManagerHandler constructs a ManagerHandler with all required dependencies.
func NewManagerHandler(cfg *config.Config, driverRepo *repository.DriverRepo, vehicleRepo *repository.VehicleRepo, shiftRepo *repository.ShiftRepo, pingRepo *repository.GPSPingRepo, alertRepo *repository.AlertRepo, zoneRepo *repository.ZoneRepo) *ManagerHandler {
	return &ManagerHandler{
		Config:      cfg,
		DriverRepo:  driverRepo,
		VehicleRepo: vehicleRepo,
		ShiftRepo:   shiftRepo,
		PingRepo:    pingRepo,
		AlertRepo:   alertRepo,
		ZoneRepo:    zoneRepo,
	}
}

// Routes returns a chi.Router with all manager endpoints behind auth middleware.
func (h *ManagerHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(h.managerAuthMiddleware)

	// Live tracking
	r.Get("/tracking/live", h.GetLivePositions)

	// Active shifts
	r.Get("/shifts/active", h.GetActiveShifts)

	// Alerts
	r.Get("/alerts", h.ListAlerts)
	r.Get("/alerts/config", h.GetAlertConfig)
	r.Put("/alerts/config", h.UpdateAlertConfig)
	r.Get("/alerts/{id}", h.GetAlert)
	r.Put("/alerts/{id}/acknowledge", h.AcknowledgeAlert)
	r.Put("/alerts/{id}/resolve", h.ResolveAlert)
	r.Put("/alerts/{id}/false-alarm", h.MarkFalseAlarm)

	// Zones
	r.Get("/zones", h.ListZones)
	r.Post("/zones", h.CreateZone)
	r.Put("/zones/{id}", h.UpdateZone)
	r.Delete("/zones/{id}", h.DeleteZone)

	return r
}

// ---------------------------------------------------------------------------
// Auth middleware
// ---------------------------------------------------------------------------

// managerAuthMiddleware validates JWT from Bearer header or HttpOnly cookie.
// Drivers are rejected -- only managers/admins/dispatchers may access.
func (h *ManagerHandler) managerAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := ""

		// 1. Check Authorization Bearer header first.
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// 2. Fall back to HttpOnly cookie.
		if tokenString == "" {
			tokenString = auth.GetAccessToken(r)
		}

		if tokenString == "" {
			jsonError(w, "missing authentication token", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(h.Config.JWTPublicKey, tokenString)
		if err != nil {
			jsonError(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Reject drivers -- they must not access manager endpoints.
		if claims.Role == "driver" {
			jsonError(w, "insufficient permissions", http.StatusForbidden)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, models.CtxTenantID, claims.TenantID)
		ctx = context.WithValue(ctx, models.CtxUserID, claims.UserID)
		ctx = context.WithValue(ctx, models.CtxUserRole, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ---------------------------------------------------------------------------
// Live tracking
// ---------------------------------------------------------------------------

// GetLivePositions handles GET /api/v1/manager/tracking/live
// Returns every active driver's latest position with a status colour.
func (h *ManagerHandler) GetLivePositions(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)

	pings, err := h.PingRepo.GetLatestByTenant(r.Context(), tenantID)
	if err != nil {
		log.Printf("manager: failed to get latest pings: %v", err)
		jsonError(w, "failed to fetch live positions", http.StatusInternalServerError)
		return
	}

	// Load alert config for offline threshold (fall back to 10 min).
	offlineThreshold := 10 * time.Minute
	alertCfg, cfgErr := h.ZoneRepo.GetAlertConfig(r.Context(), tenantID)
	if cfgErr == nil && alertCfg != nil && alertCfg.OfflineThresholdMinutes > 0 {
		offlineThreshold = time.Duration(alertCfg.OfflineThresholdMinutes) * time.Minute
	}

	now := time.Now()
	positions := make([]map[string]interface{}, 0, len(pings))
	for _, p := range pings {
		// Resolve driver name.
		driverName := ""
		driver, _ := h.DriverRepo.GetByID(r.Context(), tenantID, p.DriverID)
		if driver != nil {
			if driver.FirstName != nil {
				driverName = *driver.FirstName
			}
			if driver.LastName != nil {
				if driverName != "" {
					driverName += " "
				}
				driverName += *driver.LastName
			}
		}

		// Resolve truck plate.
		truckPlate := ""
		truck, _ := h.VehicleRepo.GetByID(r.Context(), tenantID, p.TruckID)
		if truck != nil && truck.LicensePlate != nil {
			truckPlate = *truck.LicensePlate
		}

		// Determine status colour.
		age := now.Sub(p.RecordedAt)
		status := "green"
		switch {
		case age > offlineThreshold:
			status = "gray" // offline
		case !p.IsMoving && age < offlineThreshold:
			status = "yellow" // stopped but still reporting
		case p.IsMoving && age < 2*time.Minute:
			status = "green" // actively moving & fresh
		}

		// Override with red if there is an active alert for this driver.
		hasAlert, _ := h.AlertRepo.HasRecentAlert(r.Context(), tenantID, p.DriverID, "", 30)
		if hasAlert {
			status = "red"
		}

		positions = append(positions, map[string]interface{}{
			"driver_id":     p.DriverID,
			"driver_name":   driverName,
			"truck_id":      p.TruckID,
			"truck_plate":   truckPlate,
			"shift_id":      p.ShiftID,
			"latitude":      p.Latitude,
			"longitude":     p.Longitude,
			"speed_kmh":     p.SpeedKmh,
			"heading":       p.Heading,
			"battery_level": p.BatteryLevel,
			"is_moving":     p.IsMoving,
			"recorded_at":   p.RecordedAt,
			"status":        status,
		})
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{"positions": positions})
}

// ---------------------------------------------------------------------------
// Active shifts
// ---------------------------------------------------------------------------

// GetActiveShifts handles GET /api/v1/manager/shifts/active
func (h *ManagerHandler) GetActiveShifts(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)

	allShifts, err := h.ShiftRepo.GetAllActive(r.Context())
	if err != nil {
		log.Printf("manager: failed to get active shifts: %v", err)
		jsonError(w, "failed to fetch active shifts", http.StatusInternalServerError)
		return
	}

	// Filter to this tenant only (GetAllActive is cross-tenant by design for the worker).
	shifts := make([]map[string]interface{}, 0)
	for _, s := range allShifts {
		if s.TenantID != tenantID {
			continue
		}

		// Resolve driver name.
		driverName := ""
		driver, _ := h.DriverRepo.GetByID(r.Context(), tenantID, s.DriverID)
		if driver != nil {
			if driver.FirstName != nil {
				driverName = *driver.FirstName
			}
			if driver.LastName != nil {
				if driverName != "" {
					driverName += " "
				}
				driverName += *driver.LastName
			}
		}

		// Resolve truck plate.
		truckPlate := ""
		truck, _ := h.VehicleRepo.GetByID(r.Context(), tenantID, s.TruckID)
		if truck != nil && truck.LicensePlate != nil {
			truckPlate = *truck.LicensePlate
		}

		shifts = append(shifts, map[string]interface{}{
			"shift_id":    s.ID,
			"driver_id":   s.DriverID,
			"driver_name": driverName,
			"truck_id":    s.TruckID,
			"truck_plate": truckPlate,
			"started_at":  s.StartedAt,
			"total_km":    s.TotalKm,
		})
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{"shifts": shifts})
}

// ---------------------------------------------------------------------------
// Alerts
// ---------------------------------------------------------------------------

// ListAlerts handles GET /api/v1/manager/alerts
// Query params: status (optional), limit (default 50), offset (default 0)
func (h *ManagerHandler) ListAlerts(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)

	status := r.URL.Query().Get("status")
	limit := intQueryParam(r, "limit", 50)
	offset := intQueryParam(r, "offset", 0)

	if status != "" {
		alerts, total, err := h.AlertRepo.ListByStatus(r.Context(), tenantID, status, limit, offset)
		if err != nil {
			log.Printf("manager: failed to list alerts by status: %v", err)
			jsonError(w, "failed to fetch alerts", http.StatusInternalServerError)
			return
		}
		if alerts == nil {
			alerts = []models.Alert{}
		}
		jsonResponse(w, http.StatusOK, map[string]interface{}{
			"alerts": alerts,
			"total":  total,
			"limit":  limit,
			"offset": offset,
		})
		return
	}

	// No status filter -- return latest alerts up to limit.
	alerts, err := h.AlertRepo.GetByTenant(r.Context(), tenantID, limit)
	if err != nil {
		log.Printf("manager: failed to list alerts: %v", err)
		jsonError(w, "failed to fetch alerts", http.StatusInternalServerError)
		return
	}
	if alerts == nil {
		alerts = []models.Alert{}
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"alerts": alerts,
		"total":  len(alerts),
		"limit":  limit,
		"offset": offset,
	})
}

// GetAlert handles GET /api/v1/manager/alerts/{id}
func (h *ManagerHandler) GetAlert(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)

	alertID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid alert id", http.StatusBadRequest)
		return
	}

	alert, err := h.AlertRepo.GetByID(r.Context(), tenantID, alertID)
	if err != nil {
		jsonError(w, "alert not found", http.StatusNotFound)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{"alert": alert})
}

// AcknowledgeAlert handles PUT /api/v1/manager/alerts/{id}/acknowledge
func (h *ManagerHandler) AcknowledgeAlert(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)

	alertID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid alert id", http.StatusBadRequest)
		return
	}

	if err := h.AlertRepo.Acknowledge(r.Context(), tenantID, alertID); err != nil {
		log.Printf("manager: failed to acknowledge alert: %v", err)
		jsonError(w, "failed to acknowledge alert", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{"status": "acknowledged"})
}

// ResolveAlert handles PUT /api/v1/manager/alerts/{id}/resolve
func (h *ManagerHandler) ResolveAlert(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)

	alertID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid alert id", http.StatusBadRequest)
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.AlertRepo.Resolve(r.Context(), tenantID, alertID, req.Notes); err != nil {
		log.Printf("manager: failed to resolve alert: %v", err)
		jsonError(w, "failed to resolve alert", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{"status": "resolved"})
}

// MarkFalseAlarm handles PUT /api/v1/manager/alerts/{id}/false-alarm
func (h *ManagerHandler) MarkFalseAlarm(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)

	alertID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid alert id", http.StatusBadRequest)
		return
	}

	if err := h.AlertRepo.MarkFalseAlarm(r.Context(), tenantID, alertID); err != nil {
		log.Printf("manager: failed to mark false alarm: %v", err)
		jsonError(w, "failed to mark alert as false alarm", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{"status": "false_alarm"})
}

// ---------------------------------------------------------------------------
// Zones
// ---------------------------------------------------------------------------

// ListZones handles GET /api/v1/manager/zones
func (h *ManagerHandler) ListZones(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)

	zones, err := h.ZoneRepo.GetByTenant(r.Context(), tenantID)
	if err != nil {
		log.Printf("manager: failed to list zones: %v", err)
		jsonError(w, "failed to fetch zones", http.StatusInternalServerError)
		return
	}
	if zones == nil {
		zones = []models.ApprovedZone{}
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{"zones": zones})
}

// CreateZone handles POST /api/v1/manager/zones
func (h *ManagerHandler) CreateZone(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)

	var req struct {
		Label        string  `json:"label"`
		Latitude     float64 `json:"latitude"`
		Longitude    float64 `json:"longitude"`
		RadiusMeters int     `json:"radius_meters"`
		Type         string  `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Label == "" {
		jsonError(w, "label is required", http.StatusBadRequest)
		return
	}
	if req.RadiusMeters <= 0 {
		jsonError(w, "radius_meters must be positive", http.StatusBadRequest)
		return
	}
	if req.Type == "" {
		req.Type = "stop"
	}

	zone := &models.ApprovedZone{
		TenantID:     tenantID,
		Label:        req.Label,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		RadiusMeters: req.RadiusMeters,
		Type:         req.Type,
	}
	if err := h.ZoneRepo.Create(r.Context(), zone); err != nil {
		log.Printf("manager: failed to create zone: %v", err)
		jsonError(w, "failed to create zone", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]interface{}{"zone": zone})
}

// UpdateZone handles PUT /api/v1/manager/zones/{id}
func (h *ManagerHandler) UpdateZone(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)

	zoneID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid zone id", http.StatusBadRequest)
		return
	}

	var req struct {
		Label        string  `json:"label"`
		Latitude     float64 `json:"latitude"`
		Longitude    float64 `json:"longitude"`
		RadiusMeters int     `json:"radius_meters"`
		Type         string  `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	zone := &models.ApprovedZone{
		Label:        req.Label,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		RadiusMeters: req.RadiusMeters,
		Type:         req.Type,
	}
	if err := h.ZoneRepo.Update(r.Context(), tenantID, zoneID, zone); err != nil {
		log.Printf("manager: failed to update zone: %v", err)
		jsonError(w, "failed to update zone", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{"status": "updated"})
}

// DeleteZone handles DELETE /api/v1/manager/zones/{id}
func (h *ManagerHandler) DeleteZone(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)

	zoneID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid zone id", http.StatusBadRequest)
		return
	}

	if err := h.ZoneRepo.Delete(r.Context(), tenantID, zoneID); err != nil {
		log.Printf("manager: failed to delete zone: %v", err)
		jsonError(w, "failed to delete zone", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{"status": "deleted"})
}

// ---------------------------------------------------------------------------
// Alert config
// ---------------------------------------------------------------------------

// GetAlertConfig handles GET /api/v1/manager/alerts/config
func (h *ManagerHandler) GetAlertConfig(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)

	cfg, err := h.ZoneRepo.GetAlertConfig(r.Context(), tenantID)
	if err != nil {
		// Return sensible defaults if no config exists yet.
		jsonResponse(w, http.StatusOK, map[string]interface{}{
			"config": map[string]interface{}{
				"max_stop_duration_minutes":  15,
				"alert_on_driver_offline":    true,
				"offline_threshold_minutes":  10,
				"notify_via_push":            true,
				"notify_via_email":           false,
				"notify_via_sms":             false,
			},
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{"config": cfg})
}

// UpdateAlertConfig handles PUT /api/v1/manager/alerts/config
func (h *ManagerHandler) UpdateAlertConfig(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(models.CtxTenantID).(uuid.UUID)

	var req struct {
		MaxStopDurationMinutes  int  `json:"max_stop_duration_minutes"`
		AlertOnDriverOffline    bool `json:"alert_on_driver_offline"`
		OfflineThresholdMinutes int  `json:"offline_threshold_minutes"`
		NotifyViaPush           bool `json:"notify_via_push"`
		NotifyViaEmail          bool `json:"notify_via_email"`
		NotifyViaSMS            bool `json:"notify_via_sms"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.MaxStopDurationMinutes <= 0 {
		req.MaxStopDurationMinutes = 15
	}
	if req.OfflineThresholdMinutes <= 0 {
		req.OfflineThresholdMinutes = 10
	}

	cfg := &models.AlertConfig{
		TenantID:                tenantID,
		MaxStopDurationMinutes:  req.MaxStopDurationMinutes,
		AlertOnDriverOffline:    req.AlertOnDriverOffline,
		OfflineThresholdMinutes: req.OfflineThresholdMinutes,
		NotifyViaPush:           req.NotifyViaPush,
		NotifyViaEmail:          req.NotifyViaEmail,
		NotifyViaSMS:            req.NotifyViaSMS,
	}
	if err := h.ZoneRepo.CreateOrUpdateAlertConfig(r.Context(), cfg); err != nil {
		log.Printf("manager: failed to update alert config: %v", err)
		jsonError(w, "failed to update alert config", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{"config": cfg})
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// intQueryParam reads an integer query parameter with a default value.
func intQueryParam(r *http.Request, key string, defaultVal int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return defaultVal
	}
	return n
}
