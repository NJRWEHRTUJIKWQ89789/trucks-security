package workers

import (
	"context"
	"log"
	"time"

	"cargomax-api/internal/models"
	"cargomax-api/internal/repository"
	"cargomax-api/internal/rest"
	"cargomax-api/internal/utils"

	"github.com/google/uuid"
)

type AlertWorker struct {
	ShiftRepo *repository.ShiftRepo
	PingRepo  *repository.GPSPingRepo
	AlertRepo *repository.AlertRepo
	ZoneRepo  *repository.ZoneRepo
	WSHub     *rest.Hub
}

func NewAlertWorker(shiftRepo *repository.ShiftRepo, pingRepo *repository.GPSPingRepo, alertRepo *repository.AlertRepo, zoneRepo *repository.ZoneRepo, hub *rest.Hub) *AlertWorker {
	return &AlertWorker{
		ShiftRepo: shiftRepo,
		PingRepo:  pingRepo,
		AlertRepo: alertRepo,
		ZoneRepo:  zoneRepo,
		WSHub:     hub,
	}
}

func (w *AlertWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Println("Alert worker started (checking every 30s)")

	for {
		select {
		case <-ctx.Done():
			log.Println("Alert worker stopped")
			return
		case <-ticker.C:
			w.checkAlerts(ctx)
		}
	}
}

func (w *AlertWorker) checkAlerts(ctx context.Context) {
	shifts, err := w.ShiftRepo.GetAllActive(ctx)
	if err != nil {
		log.Printf("alert worker: failed to get active shifts: %v", err)
		return
	}

	for _, shift := range shifts {
		// Get alert config for this tenant
		config, err := w.ZoneRepo.GetAlertConfig(ctx, shift.TenantID)
		if err != nil {
			// Use defaults if no config
			config = &models.AlertConfig{
				MaxStopDurationMinutes:  5,
				AlertOnDriverOffline:    true,
				OfflineThresholdMinutes: 3,
			}
		}

		// Get latest ping for this driver
		latestPing, err := w.PingRepo.GetLatestByDriver(ctx, shift.TenantID, shift.DriverID)
		if err != nil {
			// No pings yet, check if shift is old enough to be considered offline
			if config.AlertOnDriverOffline && time.Since(shift.StartedAt) > time.Duration(config.OfflineThresholdMinutes)*time.Minute {
				w.createAlert(ctx, shift.TenantID, shift.DriverID, &shift.ID, "driver_offline", nil, nil, 0)
			}
			continue
		}

		// Check driver offline
		if config.AlertOnDriverOffline {
			timeSinceLastPing := time.Since(latestPing.RecordedAt)
			if timeSinceLastPing > time.Duration(config.OfflineThresholdMinutes)*time.Minute {
				w.createAlert(ctx, shift.TenantID, shift.DriverID, &shift.ID, "driver_offline", nil, nil, 0)
			}
		}

		// Check unauthorized stop
		if !latestPing.IsMoving {
			// Calculate how long the driver has been stopped
			// Simple check: if the latest ping is stationary and recorded_at is older than threshold
			stopDuration := time.Since(latestPing.RecordedAt)
			if stopDuration > time.Duration(config.MaxStopDurationMinutes)*time.Minute {
				// Check if within an approved zone
				zones, err := w.ZoneRepo.GetByTenant(ctx, shift.TenantID)
				if err != nil {
					continue
				}

				inApprovedZone := false
				var nearestZoneID *uuid.UUID
				var nearestDistance *float64
				minDist := float64(999999999)

				for _, zone := range zones {
					dist := utils.HaversineMeters(latestPing.Latitude, latestPing.Longitude, zone.Latitude, zone.Longitude)
					if dist < minDist {
						minDist = dist
						zid := zone.ID
						nearestZoneID = &zid
						nearestDistance = &minDist
					}
					if dist <= float64(zone.RadiusMeters) {
						inApprovedZone = true
						break
					}
				}

				if !inApprovedZone {
					lat := latestPing.Latitude
					lng := latestPing.Longitude
					alert := &models.Alert{
						TenantID:             shift.TenantID,
						DriverID:             shift.DriverID,
						ShiftID:              &shift.ID,
						Type:                 "unauthorized_stop",
						Status:               "triggered",
						StopLatitude:         &lat,
						StopLongitude:        &lng,
						StopDurationSeconds:  int(stopDuration.Seconds()),
						NearestZoneID:        nearestZoneID,
						NearestZoneDistanceM: nearestDistance,
					}
					if err := w.AlertRepo.Create(ctx, alert); err != nil {
						log.Printf("alert worker: failed to create unauthorized_stop alert: %v", err)
					} else if w.WSHub != nil {
						w.WSHub.BroadcastAlert(shift.TenantID, alert)
					}
				}
			}
		}
	}
}

func (w *AlertWorker) createAlert(ctx context.Context, tenantID, driverID uuid.UUID, shiftID *uuid.UUID, alertType string, lat, lng *float64, stopDuration int) {
	alert := &models.Alert{
		TenantID:            tenantID,
		DriverID:            driverID,
		ShiftID:             shiftID,
		Type:                alertType,
		Status:              "triggered",
		StopLatitude:        lat,
		StopLongitude:       lng,
		StopDurationSeconds: stopDuration,
	}
	if err := w.AlertRepo.Create(ctx, alert); err != nil {
		log.Printf("alert worker: failed to create %s alert: %v", alertType, err)
	} else if w.WSHub != nil {
		w.WSHub.BroadcastAlert(tenantID, alert)
	}
}
