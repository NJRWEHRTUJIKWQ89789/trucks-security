package models

import (
	"time"

	"github.com/google/uuid"
)

// Order represents a customer order in the system.
type Order struct {
	ID                 uuid.UUID  `json:"id"`
	TenantID           uuid.UUID  `json:"tenant_id"`
	OrderNumber        string     `json:"order_number"`
	CustomerName       *string    `json:"customer_name"`
	CustomerEmail      *string    `json:"customer_email"`
	Status             string     `json:"status"`
	Type               string     `json:"type"`
	TotalAmount        *float64   `json:"total_amount"`
	ShipmentID         *uuid.UUID `json:"shipment_id"`
	ScheduledDate      *time.Time `json:"scheduled_date"`
	ReturnReason       *string    `json:"return_reason"`
	CancellationReason *string    `json:"cancellation_reason"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
