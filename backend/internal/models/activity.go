package models

import (
	"time"

	"github.com/google/uuid"
)

// Activity represents an entry in the activity log.
type Activity struct {
	ID         uuid.UUID  `json:"id"`
	TenantID   uuid.UUID  `json:"tenant_id"`
	UserID     *uuid.UUID `json:"user_id"`
	Action     *string    `json:"action"`
	EntityType *string    `json:"entity_type"`
	EntityID   *uuid.UUID `json:"entity_id"`
	Details    *string    `json:"details"`
	IPAddress  *string    `json:"ip_address"`
	CreatedAt  time.Time  `json:"created_at"`
}
