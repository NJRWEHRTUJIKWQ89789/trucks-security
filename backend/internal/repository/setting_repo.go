package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SettingRepo handles database operations for settings.
type SettingRepo struct {
	db *pgxpool.Pool
}

// NewSettingRepo creates a new SettingRepo instance.
func NewSettingRepo(db *pgxpool.Pool) *SettingRepo {
	return &SettingRepo{db: db}
}

// Get retrieves a setting by key within a tenant.
func (r *SettingRepo) Get(ctx context.Context, tenantID uuid.UUID, key string) (*models.Setting, error) {
	s := &models.Setting{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, key, value, category, updated_by, updated_at
		 FROM settings WHERE tenant_id = $1 AND key = $2`,
		tenantID, key,
	).Scan(&s.ID, &s.TenantID, &s.Key, &s.Value, &s.Category, &s.UpdatedBy, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get setting: %w", err)
	}
	return s, nil
}

// GetByCategory retrieves all settings in a category within a tenant.
func (r *SettingRepo) GetByCategory(ctx context.Context, tenantID uuid.UUID, category string) ([]models.Setting, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, key, value, category, updated_by, updated_at
		 FROM settings WHERE tenant_id = $1 AND category = $2 ORDER BY key ASC`,
		tenantID, category,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings by category: %w", err)
	}
	defer rows.Close()

	var settings []models.Setting
	for rows.Next() {
		var s models.Setting
		if err := rows.Scan(&s.ID, &s.TenantID, &s.Key, &s.Value, &s.Category, &s.UpdatedBy, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan setting: %w", err)
		}
		settings = append(settings, s)
	}
	return settings, nil
}

// Set upserts a setting value within a tenant.
func (r *SettingRepo) Set(ctx context.Context, s *models.Setting) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO settings (id, tenant_id, key, value, category, updated_by, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, NOW())
		 ON CONFLICT (tenant_id, key)
		 DO UPDATE SET value = $4, category = $5, updated_by = $6, updated_at = NOW()`,
		uuid.New(), s.TenantID, s.Key, s.Value, s.Category, s.UpdatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to set setting: %w", err)
	}
	return nil
}
