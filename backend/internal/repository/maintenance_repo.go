package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MaintenanceRepo handles database operations for maintenance records.
type MaintenanceRepo struct {
	db *pgxpool.Pool
}

// NewMaintenanceRepo creates a new MaintenanceRepo instance.
func NewMaintenanceRepo(db *pgxpool.Pool) *MaintenanceRepo {
	return &MaintenanceRepo{db: db}
}

// Create inserts a new maintenance record.
func (r *MaintenanceRepo) Create(ctx context.Context, m *models.MaintenanceRecord) error {
	m.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO maintenance_records (id, tenant_id, vehicle_id, type, description, status, scheduled_date, completed_date, cost, mechanic, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())`,
		m.ID, m.TenantID, m.VehicleID, m.Type, m.Description, m.Status, m.ScheduledDate, m.CompletedDate, m.Cost, m.Mechanic,
	)
	if err != nil {
		return fmt.Errorf("failed to create maintenance record: %w", err)
	}
	return nil
}

// GetByID retrieves a maintenance record by ID within a tenant.
func (r *MaintenanceRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.MaintenanceRecord, error) {
	m := &models.MaintenanceRecord{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, vehicle_id, type, description, status, scheduled_date, completed_date, cost, mechanic, created_at, updated_at
		 FROM maintenance_records WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(&m.ID, &m.TenantID, &m.VehicleID, &m.Type, &m.Description, &m.Status, &m.ScheduledDate, &m.CompletedDate, &m.Cost, &m.Mechanic, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get maintenance record by id: %w", err)
	}
	return m, nil
}

// List returns a paginated list of maintenance records, optionally filtered by vehicle ID.
func (r *MaintenanceRepo) List(ctx context.Context, tenantID uuid.UUID, vehicleID *uuid.UUID, page, perPage int) ([]models.MaintenanceRecord, int, error) {
	var total int
	if vehicleID != nil {
		err := r.db.QueryRow(ctx,
			`SELECT COUNT(*) FROM maintenance_records WHERE tenant_id = $1 AND vehicle_id = $2`,
			tenantID, *vehicleID,
		).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count maintenance records: %w", err)
		}
	} else {
		err := r.db.QueryRow(ctx,
			`SELECT COUNT(*) FROM maintenance_records WHERE tenant_id = $1`,
			tenantID,
		).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count maintenance records: %w", err)
		}
	}

	offset := (page - 1) * perPage
	var query string
	var args []interface{}
	if vehicleID != nil {
		query = `SELECT id, tenant_id, vehicle_id, type, description, status, scheduled_date, completed_date, cost, mechanic, created_at, updated_at
				 FROM maintenance_records WHERE tenant_id = $1 AND vehicle_id = $2 ORDER BY scheduled_date DESC LIMIT $3 OFFSET $4`
		args = []interface{}{tenantID, *vehicleID, perPage, offset}
	} else {
		query = `SELECT id, tenant_id, vehicle_id, type, description, status, scheduled_date, completed_date, cost, mechanic, created_at, updated_at
				 FROM maintenance_records WHERE tenant_id = $1 ORDER BY scheduled_date DESC LIMIT $2 OFFSET $3`
		args = []interface{}{tenantID, perPage, offset}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list maintenance records: %w", err)
	}
	defer rows.Close()

	var records []models.MaintenanceRecord
	for rows.Next() {
		var m models.MaintenanceRecord
		if err := rows.Scan(&m.ID, &m.TenantID, &m.VehicleID, &m.Type, &m.Description, &m.Status, &m.ScheduledDate, &m.CompletedDate, &m.Cost, &m.Mechanic, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan maintenance record: %w", err)
		}
		records = append(records, m)
	}
	return records, total, nil
}

// Update modifies an existing maintenance record.
func (r *MaintenanceRepo) Update(ctx context.Context, tenantID, id uuid.UUID, m *models.MaintenanceRecord) error {
	_, err := r.db.Exec(ctx,
		`UPDATE maintenance_records SET vehicle_id = $1, type = $2, description = $3, status = $4, scheduled_date = $5, completed_date = $6, cost = $7, mechanic = $8, updated_at = NOW()
		 WHERE id = $9 AND tenant_id = $10`,
		m.VehicleID, m.Type, m.Description, m.Status, m.ScheduledDate, m.CompletedDate, m.Cost, m.Mechanic, id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update maintenance record: %w", err)
	}
	return nil
}

// GetUpcoming returns maintenance records scheduled in the next 7 days that are not completed.
func (r *MaintenanceRepo) GetUpcoming(ctx context.Context, tenantID uuid.UUID) ([]models.MaintenanceRecord, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, vehicle_id, type, description, status, scheduled_date, completed_date, cost, mechanic, created_at, updated_at
		 FROM maintenance_records WHERE tenant_id = $1 AND status = 'scheduled' AND scheduled_date >= NOW() AND scheduled_date <= NOW() + INTERVAL '7 days'
		 ORDER BY scheduled_date ASC`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming maintenance: %w", err)
	}
	defer rows.Close()

	var records []models.MaintenanceRecord
	for rows.Next() {
		var m models.MaintenanceRecord
		if err := rows.Scan(&m.ID, &m.TenantID, &m.VehicleID, &m.Type, &m.Description, &m.Status, &m.ScheduledDate, &m.CompletedDate, &m.Cost, &m.Mechanic, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan upcoming maintenance: %w", err)
		}
		records = append(records, m)
	}
	return records, nil
}

// GetOverdue returns maintenance records that are past their scheduled date and not completed.
func (r *MaintenanceRepo) GetOverdue(ctx context.Context, tenantID uuid.UUID) ([]models.MaintenanceRecord, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, vehicle_id, type, description, status, scheduled_date, completed_date, cost, mechanic, created_at, updated_at
		 FROM maintenance_records WHERE tenant_id = $1 AND status = 'scheduled' AND scheduled_date < NOW()
		 ORDER BY scheduled_date ASC`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue maintenance: %w", err)
	}
	defer rows.Close()

	var records []models.MaintenanceRecord
	for rows.Next() {
		var m models.MaintenanceRecord
		if err := rows.Scan(&m.ID, &m.TenantID, &m.VehicleID, &m.Type, &m.Description, &m.Status, &m.ScheduledDate, &m.CompletedDate, &m.Cost, &m.Mechanic, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan overdue maintenance: %w", err)
		}
		records = append(records, m)
	}
	return records, nil
}
