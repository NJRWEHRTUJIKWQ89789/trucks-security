package models

import (
	"time"

	"github.com/google/uuid"
)

// Vendor represents a third-party vendor or supplier.
type Vendor struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	Name          string     `json:"name"`
	ContactPerson *string    `json:"contact_person"`
	Email         *string    `json:"email"`
	Phone         *string    `json:"phone"`
	Address       *string    `json:"address"`
	Category      *string    `json:"category"`
	Rating        *float64   `json:"rating"`
	ContractStart *time.Time `json:"contract_start"`
	ContractEnd   *time.Time `json:"contract_end"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
