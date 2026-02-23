package models

import (
	"time"

	"github.com/google/uuid"
)

// Setting represents a key-value configuration setting for a tenant.
type Setting struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	Key       string     `json:"key"`
	Value     *string    `json:"value"`
	Category  *string    `json:"category"`
	UpdatedBy *uuid.UUID `json:"updated_by"`
	UpdatedAt time.Time  `json:"updated_at"`
}
