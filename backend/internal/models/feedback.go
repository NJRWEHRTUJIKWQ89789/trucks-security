package models

import (
	"time"

	"github.com/google/uuid"
)

// Feedback represents client feedback on services.
type Feedback struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	ClientID   uuid.UUID `json:"client_id"`
	Rating     int       `json:"rating"`
	Comment    *string   `json:"comment"`
	Category   *string   `json:"category"`
	CreatedAt  time.Time `json:"created_at"`
	ClientName *string   `json:"client_name,omitempty"` // populated by joined queries
}
