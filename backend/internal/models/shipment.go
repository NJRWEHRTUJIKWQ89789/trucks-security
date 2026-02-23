package models

import (
	"time"

	"github.com/google/uuid"
)

// Shipment represents a shipment tracked by the logistics system.
type Shipment struct {
	ID                uuid.UUID  `json:"id"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	TrackingNumber    string     `json:"tracking_number"`
	Origin            *string    `json:"origin"`
	Destination       *string    `json:"destination"`
	Status            string     `json:"status"`
	Carrier           *string    `json:"carrier"`
	Weight            *float64   `json:"weight"`
	Dimensions        *string    `json:"dimensions"`
	EstimatedDelivery *time.Time `json:"estimated_delivery"`
	ActualDelivery    *time.Time `json:"actual_delivery"`
	CustomerName      *string    `json:"customer_name"`
	CustomerEmail     *string    `json:"customer_email"`
	Notes             *string    `json:"notes"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
