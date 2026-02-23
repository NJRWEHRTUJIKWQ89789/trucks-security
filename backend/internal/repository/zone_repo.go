package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ZoneRepo struct {
	db *pgxpool.Pool
}

func NewZoneRepo(db *pgxpool.Pool) *ZoneRepo {
	return &ZoneRepo{db: db}
}

func (r *ZoneRepo) GetByTenant(ctx context.Context, tenantID uuid.UUID) ([]models.ApprovedZone, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, label, latitude, longitude, radius_meters, type, created_at, updated_at
		 FROM approved_zones WHERE tenant_id = $1`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query zones: %w", err)
	}
	defer rows.Close()

	var zones []models.ApprovedZone
	for rows.Next() {
		var z models.ApprovedZone
		if err := rows.Scan(&z.ID, &z.TenantID, &z.Label, &z.Latitude, &z.Longitude, &z.RadiusMeters, &z.Type, &z.CreatedAt, &z.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan zone: %w", err)
		}
		zones = append(zones, z)
	}
	return zones, nil
}

func (r *ZoneRepo) GetAlertConfig(ctx context.Context, tenantID uuid.UUID) (*models.AlertConfig, error) {
	c := &models.AlertConfig{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, max_stop_duration_minutes, alert_on_driver_offline, offline_threshold_minutes, notify_via_push, notify_via_email, notify_via_sms, updated_at
		 FROM alert_config WHERE tenant_id = $1`,
		tenantID,
	).Scan(&c.ID, &c.TenantID, &c.MaxStopDurationMinutes, &c.AlertOnDriverOffline, &c.OfflineThresholdMinutes, &c.NotifyViaPush, &c.NotifyViaEmail, &c.NotifyViaSMS, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Create inserts a new approved zone
func (r *ZoneRepo) Create(ctx context.Context, z *models.ApprovedZone) error {
	z.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO approved_zones (id, tenant_id, label, latitude, longitude, radius_meters, type, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())`,
		z.ID, z.TenantID, z.Label, z.Latitude, z.Longitude, z.RadiusMeters, z.Type,
	)
	if err != nil {
		return fmt.Errorf("failed to create zone: %w", err)
	}
	return nil
}

// Update modifies an existing approved zone
func (r *ZoneRepo) Update(ctx context.Context, tenantID, id uuid.UUID, z *models.ApprovedZone) error {
	_, err := r.db.Exec(ctx,
		`UPDATE approved_zones SET label = $3, latitude = $4, longitude = $5, radius_meters = $6, type = $7, updated_at = NOW()
		 WHERE tenant_id = $1 AND id = $2`,
		tenantID, id, z.Label, z.Latitude, z.Longitude, z.RadiusMeters, z.Type,
	)
	if err != nil {
		return fmt.Errorf("failed to update zone: %w", err)
	}
	return nil
}

// Delete removes an approved zone
func (r *ZoneRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM approved_zones WHERE tenant_id = $1 AND id = $2`,
		tenantID, id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete zone: %w", err)
	}
	return nil
}

// GetByID retrieves a single zone by ID within a tenant
func (r *ZoneRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.ApprovedZone, error) {
	z := &models.ApprovedZone{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, label, latitude, longitude, radius_meters, type, created_at, updated_at
		 FROM approved_zones WHERE tenant_id = $1 AND id = $2`,
		tenantID, id,
	).Scan(&z.ID, &z.TenantID, &z.Label, &z.Latitude, &z.Longitude, &z.RadiusMeters, &z.Type, &z.CreatedAt, &z.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get zone by id: %w", err)
	}
	return z, nil
}

// CreateOrUpdateAlertConfig upserts the alert configuration for a tenant
func (r *ZoneRepo) CreateOrUpdateAlertConfig(ctx context.Context, c *models.AlertConfig) error {
	c.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO alert_config (id, tenant_id, max_stop_duration_minutes, alert_on_driver_offline, offline_threshold_minutes, notify_via_push, notify_via_email, notify_via_sms, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		 ON CONFLICT (tenant_id) DO UPDATE SET max_stop_duration_minutes = $3, alert_on_driver_offline = $4, offline_threshold_minutes = $5, notify_via_push = $6, notify_via_email = $7, notify_via_sms = $8, updated_at = NOW()`,
		c.ID, c.TenantID, c.MaxStopDurationMinutes, c.AlertOnDriverOffline, c.OfflineThresholdMinutes, c.NotifyViaPush, c.NotifyViaEmail, c.NotifyViaSMS,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert alert config: %w", err)
	}
	return nil
}
