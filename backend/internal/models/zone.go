package models

import (
	"time"

	"github.com/google/uuid"
)

type ApprovedZone struct {
	ID           uuid.UUID `json:"id"`
	TenantID     uuid.UUID `json:"tenant_id"`
	Label        string    `json:"label"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	RadiusMeters int       `json:"radius_meters"`
	Type         string    `json:"type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AlertConfig struct {
	ID                       uuid.UUID `json:"id"`
	TenantID                 uuid.UUID `json:"tenant_id"`
	MaxStopDurationMinutes   int       `json:"max_stop_duration_minutes"`
	AlertOnDriverOffline     bool      `json:"alert_on_driver_offline"`
	OfflineThresholdMinutes  int       `json:"offline_threshold_minutes"`
	NotifyViaPush            bool      `json:"notify_via_push"`
	NotifyViaEmail           bool      `json:"notify_via_email"`
	NotifyViaSMS             bool      `json:"notify_via_sms"`
	UpdatedAt                time.Time `json:"updated_at"`
}
