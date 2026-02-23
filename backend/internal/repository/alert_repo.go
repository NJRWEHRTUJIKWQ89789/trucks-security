package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlertRepo struct {
	db *pgxpool.Pool
}

func NewAlertRepo(db *pgxpool.Pool) *AlertRepo {
	return &AlertRepo{db: db}
}

func (r *AlertRepo) Create(ctx context.Context, a *models.Alert) error {
	a.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO alerts (id, tenant_id, driver_id, shift_id, type, status, stop_latitude, stop_longitude, stop_duration_seconds, nearest_zone_id, nearest_zone_distance_meters, triggered_at, notified_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())`,
		a.ID, a.TenantID, a.DriverID, a.ShiftID, a.Type, a.Status, a.StopLatitude, a.StopLongitude, a.StopDurationSeconds, a.NearestZoneID, a.NearestZoneDistanceM,
	)
	if err != nil {
		return fmt.Errorf("failed to create alert: %w", err)
	}
	return nil
}

func (r *AlertRepo) GetByTenant(ctx context.Context, tenantID uuid.UUID, limit int) ([]models.Alert, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, driver_id, shift_id, type, status, stop_latitude, stop_longitude, stop_duration_seconds, nearest_zone_id, nearest_zone_distance_meters, manager_notes, triggered_at, notified_at, acknowledged_at, resolved_at, created_at
		 FROM alerts WHERE tenant_id = $1 ORDER BY triggered_at DESC LIMIT $2`,
		tenantID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query alerts: %w", err)
	}
	defer rows.Close()

	var alerts []models.Alert
	for rows.Next() {
		var a models.Alert
		if err := rows.Scan(&a.ID, &a.TenantID, &a.DriverID, &a.ShiftID, &a.Type, &a.Status, &a.StopLatitude, &a.StopLongitude, &a.StopDurationSeconds, &a.NearestZoneID, &a.NearestZoneDistanceM, &a.ManagerNotes, &a.TriggeredAt, &a.NotifiedAt, &a.AcknowledgedAt, &a.ResolvedAt, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan alert: %w", err)
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}

// GetByID retrieves a single alert by ID within a tenant
func (r *AlertRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Alert, error) {
	a := &models.Alert{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, driver_id, shift_id, type, status, stop_latitude, stop_longitude, stop_duration_seconds, nearest_zone_id, nearest_zone_distance_meters, manager_notes, triggered_at, notified_at, acknowledged_at, resolved_at, created_at
		 FROM alerts WHERE tenant_id = $1 AND id = $2`,
		tenantID, id,
	).Scan(&a.ID, &a.TenantID, &a.DriverID, &a.ShiftID, &a.Type, &a.Status, &a.StopLatitude, &a.StopLongitude, &a.StopDurationSeconds, &a.NearestZoneID, &a.NearestZoneDistanceM, &a.ManagerNotes, &a.TriggeredAt, &a.NotifiedAt, &a.AcknowledgedAt, &a.ResolvedAt, &a.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert by id: %w", err)
	}
	return a, nil
}

// ListByStatus returns alerts filtered by status for a tenant (paginated)
func (r *AlertRepo) ListByStatus(ctx context.Context, tenantID uuid.UUID, status string, limit, offset int) ([]models.Alert, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM alerts WHERE tenant_id = $1 AND status = $2`,
		tenantID, status,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count alerts by status: %w", err)
	}

	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, driver_id, shift_id, type, status, stop_latitude, stop_longitude, stop_duration_seconds, nearest_zone_id, nearest_zone_distance_meters, manager_notes, triggered_at, notified_at, acknowledged_at, resolved_at, created_at
		 FROM alerts WHERE tenant_id = $1 AND status = $2 ORDER BY triggered_at DESC LIMIT $3 OFFSET $4`,
		tenantID, status, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query alerts by status: %w", err)
	}
	defer rows.Close()

	var alerts []models.Alert
	for rows.Next() {
		var a models.Alert
		if err := rows.Scan(&a.ID, &a.TenantID, &a.DriverID, &a.ShiftID, &a.Type, &a.Status, &a.StopLatitude, &a.StopLongitude, &a.StopDurationSeconds, &a.NearestZoneID, &a.NearestZoneDistanceM, &a.ManagerNotes, &a.TriggeredAt, &a.NotifiedAt, &a.AcknowledgedAt, &a.ResolvedAt, &a.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan alert: %w", err)
		}
		alerts = append(alerts, a)
	}
	return alerts, total, nil
}

// Acknowledge sets an alert's status to 'acknowledged' with timestamp
func (r *AlertRepo) Acknowledge(ctx context.Context, tenantID, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE alerts SET status = 'acknowledged', acknowledged_at = NOW() WHERE tenant_id = $1 AND id = $2`,
		tenantID, id,
	)
	if err != nil {
		return fmt.Errorf("failed to acknowledge alert: %w", err)
	}
	return nil
}

// Resolve sets an alert's status to 'resolved' with timestamp and optional notes
func (r *AlertRepo) Resolve(ctx context.Context, tenantID, id uuid.UUID, notes string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE alerts SET status = 'resolved', resolved_at = NOW(), manager_notes = $3 WHERE tenant_id = $1 AND id = $2`,
		tenantID, id, notes,
	)
	if err != nil {
		return fmt.Errorf("failed to resolve alert: %w", err)
	}
	return nil
}

// MarkFalseAlarm sets an alert's status to 'false_alarm'
func (r *AlertRepo) MarkFalseAlarm(ctx context.Context, tenantID, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE alerts SET status = 'false_alarm' WHERE tenant_id = $1 AND id = $2`,
		tenantID, id,
	)
	if err != nil {
		return fmt.Errorf("failed to mark alert as false alarm: %w", err)
	}
	return nil
}

// HasRecentAlert checks if a similar alert already exists for this driver/shift in the last N minutes (to avoid duplicate alerts)
func (r *AlertRepo) HasRecentAlert(ctx context.Context, tenantID, driverID uuid.UUID, alertType string, withinMinutes int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM alerts WHERE tenant_id = $1 AND driver_id = $2 AND type = $3 AND triggered_at > NOW() - ($4 || ' minutes')::interval)`,
		tenantID, driverID, alertType, withinMinutes,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check for recent alert: %w", err)
	}
	return exists, nil
}
