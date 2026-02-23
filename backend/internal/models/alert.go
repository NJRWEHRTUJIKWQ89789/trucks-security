package models

import (
	"time"

	"github.com/google/uuid"
)

type Alert struct {
	ID                       uuid.UUID  `json:"id"`
	TenantID                 uuid.UUID  `json:"tenant_id"`
	DriverID                 uuid.UUID  `json:"driver_id"`
	ShiftID                  *uuid.UUID `json:"shift_id,omitempty"`
	Type                     string     `json:"type"`
	Status                   string     `json:"status"`
	StopLatitude             *float64   `json:"stop_latitude,omitempty"`
	StopLongitude            *float64   `json:"stop_longitude,omitempty"`
	StopDurationSeconds      int        `json:"stop_duration_seconds"`
	NearestZoneID            *uuid.UUID `json:"nearest_zone_id,omitempty"`
	NearestZoneDistanceM     *float64   `json:"nearest_zone_distance_meters,omitempty"`
	ManagerNotes             *string    `json:"manager_notes,omitempty"`
	TriggeredAt              time.Time  `json:"triggered_at"`
	NotifiedAt               *time.Time `json:"notified_at,omitempty"`
	AcknowledgedAt           *time.Time `json:"acknowledged_at,omitempty"`
	ResolvedAt               *time.Time `json:"resolved_at,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
}
