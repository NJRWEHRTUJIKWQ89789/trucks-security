package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ActivityRepo handles database operations for the activity log.
type ActivityRepo struct {
	db *pgxpool.Pool
}

// NewActivityRepo creates a new ActivityRepo instance.
func NewActivityRepo(db *pgxpool.Pool) *ActivityRepo {
	return &ActivityRepo{db: db}
}

// Create inserts a new activity log entry.
func (r *ActivityRepo) Create(ctx context.Context, a *models.Activity) error {
	a.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO activity_log (id, tenant_id, user_id, action, entity_type, entity_id, details, ip_address, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())`,
		a.ID, a.TenantID, a.UserID, a.Action, a.EntityType, a.EntityID, a.Details, a.IPAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to create activity log: %w", err)
	}
	return nil
}

// List returns a paginated list of activity log entries within a tenant.
func (r *ActivityRepo) List(ctx context.Context, tenantID uuid.UUID, page, perPage int) ([]models.Activity, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM activity_log WHERE tenant_id = $1`,
		tenantID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count activities: %w", err)
	}

	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, user_id, action, entity_type, entity_id, details, ip_address, created_at
		 FROM activity_log WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		tenantID, perPage, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list activities: %w", err)
	}
	defer rows.Close()

	var activities []models.Activity
	for rows.Next() {
		var a models.Activity
		if err := rows.Scan(&a.ID, &a.TenantID, &a.UserID, &a.Action, &a.EntityType, &a.EntityID, &a.Details, &a.IPAddress, &a.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan activity: %w", err)
		}
		activities = append(activities, a)
	}
	return activities, total, nil
}

// ListByEntity returns a paginated list of activity log entries for a specific entity within a tenant.
func (r *ActivityRepo) ListByEntity(ctx context.Context, tenantID uuid.UUID, entityType string, entityID uuid.UUID, page, perPage int) ([]models.Activity, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM activity_log WHERE tenant_id = $1 AND entity_type = $2 AND entity_id = $3`,
		tenantID, entityType, entityID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count entity activities: %w", err)
	}

	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, user_id, action, entity_type, entity_id, details, ip_address, created_at
		 FROM activity_log WHERE tenant_id = $1 AND entity_type = $2 AND entity_id = $3 ORDER BY created_at DESC LIMIT $4 OFFSET $5`,
		tenantID, entityType, entityID, perPage, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list entity activities: %w", err)
	}
	defer rows.Close()

	var activities []models.Activity
	for rows.Next() {
		var a models.Activity
		if err := rows.Scan(&a.ID, &a.TenantID, &a.UserID, &a.Action, &a.EntityType, &a.EntityID, &a.Details, &a.IPAddress, &a.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan entity activity: %w", err)
		}
		activities = append(activities, a)
	}
	return activities, total, nil
}
