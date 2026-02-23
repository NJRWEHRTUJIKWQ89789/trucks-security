package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user account within a tenant.
type User struct {
	ID               uuid.UUID `json:"id"`
	TenantID         uuid.UUID `json:"tenant_id"`
	Email            string    `json:"email"`
	PasswordHash     string    `json:"-"`
	FirstName        *string   `json:"first_name"`
	LastName         *string   `json:"last_name"`
	Role             string    `json:"role"`
	EmailVerified    bool      `json:"email_verified"`
	EmailVerifyToken *string   `json:"email_verify_token"`
	AvatarURL        *string   `json:"avatar_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
