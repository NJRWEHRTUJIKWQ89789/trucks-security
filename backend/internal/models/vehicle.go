package models

import (
	"time"

	"github.com/google/uuid"
)

// Vehicle represents a fleet vehicle tracked by the logistics system.
type Vehicle struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	VehicleID    string     `json:"vehicle_id"`
	Name         *string    `json:"name"`
	Type         *string    `json:"type"`
	Status       string     `json:"status"`
	FuelLevel    int        `json:"fuel_level"`
	Mileage      int        `json:"mileage"`
	LastService  *time.Time `json:"last_service"`
	NextService  *time.Time `json:"next_service"`
	LicensePlate *string    `json:"license_plate"`
	Year         *int       `json:"year"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
