package repository

import (
	"context"
	"fmt"
	"time"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShiftRepo struct {
	db *pgxpool.Pool
}

func NewShiftRepo(db *pgxpool.Pool) *ShiftRepo {
	return &ShiftRepo{db: db}
}

func (r *ShiftRepo) Create(ctx context.Context, s *models.Shift) error {
	s.ID = uuid.New()
	s.StartedAt = time.Now()
	s.Status = "active"
	_, err := r.db.Exec(ctx,
		`INSERT INTO shifts (id, tenant_id, driver_id, truck_id, started_at, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`,
		s.ID, s.TenantID, s.DriverID, s.TruckID, s.StartedAt, s.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to create shift: %w", err)
	}
	return nil
}

func (r *ShiftRepo) End(ctx context.Context, tenantID, driverID, shiftID uuid.UUID, totalKm float64) (*models.Shift, error) {
	s := &models.Shift{}
	now := time.Now()
	err := r.db.QueryRow(ctx,
		`UPDATE shifts SET ended_at = $1, status = 'completed', total_km = $2, updated_at = $1
		 WHERE id = $3 AND tenant_id = $4 AND driver_id = $5 AND status = 'active'
		 RETURNING id, tenant_id, driver_id, truck_id, started_at, ended_at, status, total_km, created_at, updated_at`,
		now, totalKm, shiftID, tenantID, driverID,
	).Scan(&s.ID, &s.TenantID, &s.DriverID, &s.TruckID, &s.StartedAt, &s.EndedAt, &s.Status, &s.TotalKm, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to end shift: %w", err)
	}
	return s, nil
}

func (r *ShiftRepo) GetByID(ctx context.Context, tenantID, shiftID uuid.UUID) (*models.Shift, error) {
	s := &models.Shift{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, driver_id, truck_id, started_at, ended_at, status, total_km, created_at, updated_at
		 FROM shifts WHERE id = $1 AND tenant_id = $2`,
		shiftID, tenantID,
	).Scan(&s.ID, &s.TenantID, &s.DriverID, &s.TruckID, &s.StartedAt, &s.EndedAt, &s.Status, &s.TotalKm, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get shift: %w", err)
	}
	return s, nil
}

func (r *ShiftRepo) GetActiveByDriver(ctx context.Context, tenantID, driverID uuid.UUID) (*models.Shift, error) {
	s := &models.Shift{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, driver_id, truck_id, started_at, ended_at, status, total_km, created_at, updated_at
		 FROM shifts WHERE tenant_id = $1 AND driver_id = $2 AND status = 'active'
		 ORDER BY started_at DESC LIMIT 1`,
		tenantID, driverID,
	).Scan(&s.ID, &s.TenantID, &s.DriverID, &s.TruckID, &s.StartedAt, &s.EndedAt, &s.Status, &s.TotalKm, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get active shift: %w", err)
	}
	return s, nil
}

func (r *ShiftRepo) GetAllActive(ctx context.Context) ([]models.Shift, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, driver_id, truck_id, started_at, ended_at, status, total_km, created_at, updated_at
		 FROM shifts WHERE status = 'active'`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query active shifts: %w", err)
	}
	defer rows.Close()

	var shifts []models.Shift
	for rows.Next() {
		var s models.Shift
		if err := rows.Scan(&s.ID, &s.TenantID, &s.DriverID, &s.TruckID, &s.StartedAt, &s.EndedAt, &s.Status, &s.TotalKm, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan shift: %w", err)
		}
		shifts = append(shifts, s)
	}
	return shifts, nil
}

func (r *ShiftRepo) IsTruckInUse(ctx context.Context, tenantID, truckID uuid.UUID) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM shifts WHERE tenant_id = $1 AND truck_id = $2 AND status = 'active'`,
		tenantID, truckID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check truck in use: %w", err)
	}
	return count > 0, nil
}
