package models

import (
	"time"

	"github.com/google/uuid"
)

// MaintenanceRecord represents a vehicle maintenance record.
type MaintenanceRecord struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	VehicleID     uuid.UUID  `json:"vehicle_id"`
	Type          *string    `json:"type"`
	Description   *string    `json:"description"`
	Status        string     `json:"status"`
	ScheduledDate *time.Time `json:"scheduled_date"`
	CompletedDate *time.Time `json:"completed_date"`
	Cost          *float64   `json:"cost"`
	Mechanic      *string    `json:"mechanic"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
