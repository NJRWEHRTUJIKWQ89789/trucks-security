package models

import (
	"time"

	"github.com/google/uuid"
)

type Shift struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	DriverID  uuid.UUID  `json:"driver_id"`
	TruckID   uuid.UUID  `json:"truck_id"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	Status    string     `json:"status"`
	TotalKm   float64   `json:"total_km"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
