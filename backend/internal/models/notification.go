package models

import (
	"time"

	"github.com/google/uuid"
)

// Notification represents a user notification.
type Notification struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     *string   `json:"title"`
	Message   *string   `json:"message"`
	Type      *string   `json:"type"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

// NotificationPreference represents a user's notification preferences per event type.
type NotificationPreference struct {
	ID           uuid.UUID `json:"id"`
	TenantID     uuid.UUID `json:"tenant_id"`
	UserID       uuid.UUID `json:"user_id"`
	EventType    string    `json:"event_type"`
	EmailEnabled bool      `json:"email_enabled"`
	SMSEnabled   bool      `json:"sms_enabled"`
	PushEnabled  bool      `json:"push_enabled"`
}
