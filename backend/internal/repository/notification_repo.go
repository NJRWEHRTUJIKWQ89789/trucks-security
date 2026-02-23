package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NotificationRepo handles database operations for notifications.
type NotificationRepo struct {
	db *pgxpool.Pool
}

// NewNotificationRepo creates a new NotificationRepo instance.
func NewNotificationRepo(db *pgxpool.Pool) *NotificationRepo {
	return &NotificationRepo{db: db}
}

// Create inserts a new notification.
func (r *NotificationRepo) Create(ctx context.Context, n *models.Notification) error {
	n.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO notifications (id, tenant_id, user_id, title, message, type, read, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`,
		n.ID, n.TenantID, n.UserID, n.Title, n.Message, n.Type, n.Read,
	)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	return nil
}

// List returns a paginated list of notifications for a user within a tenant.
func (r *NotificationRepo) List(ctx context.Context, tenantID, userID uuid.UUID, page, perPage int) ([]models.Notification, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM notifications WHERE tenant_id = $1 AND user_id = $2`,
		tenantID, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count notifications: %w", err)
	}

	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, user_id, title, message, type, read, created_at
		 FROM notifications WHERE tenant_id = $1 AND user_id = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`,
		tenantID, userID, perPage, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.TenantID, &n.UserID, &n.Title, &n.Message, &n.Type, &n.Read, &n.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, n)
	}
	return notifications, total, nil
}

// GetByID retrieves a notification by ID within a tenant, scoped to the owning user.
func (r *NotificationRepo) GetByID(ctx context.Context, tenantID, userID, id uuid.UUID) (*models.Notification, error) {
	n := &models.Notification{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, user_id, title, message, type, read, created_at
		 FROM notifications WHERE id = $1 AND tenant_id = $2 AND user_id = $3`,
		id, tenantID, userID,
	).Scan(&n.ID, &n.TenantID, &n.UserID, &n.Title, &n.Message, &n.Type, &n.Read, &n.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification by id: %w", err)
	}
	return n, nil
}

// MarkRead marks a notification as read, scoped to the owning user.
func (r *NotificationRepo) MarkRead(ctx context.Context, tenantID, userID, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE notifications SET read = TRUE WHERE id = $1 AND tenant_id = $2 AND user_id = $3`,
		id, tenantID, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}

// GetUnreadCount returns the count of unread notifications for a user within a tenant.
func (r *NotificationRepo) GetUnreadCount(ctx context.Context, tenantID, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM notifications WHERE tenant_id = $1 AND user_id = $2 AND read = FALSE`,
		tenantID, userID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}
	return count, nil
}

// GetPreferences returns notification preferences for a user within a tenant.
func (r *NotificationRepo) GetPreferences(ctx context.Context, tenantID, userID uuid.UUID) ([]models.NotificationPreference, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, user_id, event_type, email_enabled, sms_enabled, push_enabled
		 FROM notification_preferences WHERE tenant_id = $1 AND user_id = $2 ORDER BY event_type ASC`,
		tenantID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification preferences: %w", err)
	}
	defer rows.Close()

	var prefs []models.NotificationPreference
	for rows.Next() {
		var p models.NotificationPreference
		if err := rows.Scan(&p.ID, &p.TenantID, &p.UserID, &p.EventType, &p.EmailEnabled, &p.SMSEnabled, &p.PushEnabled); err != nil {
			return nil, fmt.Errorf("failed to scan notification preference: %w", err)
		}
		prefs = append(prefs, p)
	}
	return prefs, nil
}

// UpdatePreference upserts a notification preference for a user within a tenant.
func (r *NotificationRepo) UpdatePreference(ctx context.Context, p *models.NotificationPreference) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO notification_preferences (id, tenant_id, user_id, event_type, email_enabled, sms_enabled, push_enabled)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (tenant_id, user_id, event_type)
		 DO UPDATE SET email_enabled = $5, sms_enabled = $6, push_enabled = $7`,
		uuid.New(), p.TenantID, p.UserID, p.EventType, p.EmailEnabled, p.SMSEnabled, p.PushEnabled,
	)
	if err != nil {
		return fmt.Errorf("failed to update notification preference: %w", err)
	}
	return nil
}
