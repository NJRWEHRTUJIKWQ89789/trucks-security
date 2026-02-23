package models

import (
	"time"

	"github.com/google/uuid"
)

// Client represents a business client of the logistics company.
type Client struct {
	ID                 uuid.UUID `json:"id"`
	TenantID           uuid.UUID `json:"tenant_id"`
	CompanyName        string    `json:"company_name"`
	ContactPerson      *string   `json:"contact_person"`
	Email              *string   `json:"email"`
	Phone              *string   `json:"phone"`
	Address            *string   `json:"address"`
	Industry           *string   `json:"industry"`
	TotalShipments     int       `json:"total_shipments"`
	TotalSpent         *float64  `json:"total_spent"`
	SatisfactionRating *float64  `json:"satisfaction_rating"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
