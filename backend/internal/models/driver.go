package models

import (
	"time"

	"github.com/google/uuid"
)

// Driver represents a delivery driver assigned to a tenant.
type Driver struct {
	ID              uuid.UUID  `json:"id"`
	TenantID        uuid.UUID  `json:"tenant_id"`
	EmployeeID      string     `json:"employee_id"`
	FirstName       *string    `json:"first_name"`
	LastName        *string    `json:"last_name"`
	Email           *string    `json:"email"`
	Phone           *string    `json:"phone"`
	LicenseNumber   *string    `json:"license_number"`
	LicenseExpiry   *time.Time `json:"license_expiry"`
	Status          string     `json:"status"`
	Rating          *float64   `json:"rating"`
	TotalDeliveries int        `json:"total_deliveries"`
	VehicleID       *uuid.UUID `json:"vehicle_id"`
	PinHash         *string    `json:"pin_hash,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
