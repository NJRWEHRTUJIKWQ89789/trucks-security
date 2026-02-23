package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// VehicleRepo handles database operations for vehicles.
type VehicleRepo struct {
	db *pgxpool.Pool
}

// NewVehicleRepo creates a new VehicleRepo instance.
func NewVehicleRepo(db *pgxpool.Pool) *VehicleRepo {
	return &VehicleRepo{db: db}
}

// Create inserts a new vehicle.
func (r *VehicleRepo) Create(ctx context.Context, v *models.Vehicle) error {
	v.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO vehicles (id, tenant_id, vehicle_id, name, type, status, fuel_level, mileage, last_service, next_service, license_plate, year, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())`,
		v.ID, v.TenantID, v.VehicleID, v.Name, v.Type, v.Status, v.FuelLevel, v.Mileage, v.LastService, v.NextService, v.LicensePlate, v.Year,
	)
	if err != nil {
		return fmt.Errorf("failed to create vehicle: %w", err)
	}
	return nil
}

// GetByID retrieves a vehicle by ID within a tenant.
func (r *VehicleRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Vehicle, error) {
	v := &models.Vehicle{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, vehicle_id, name, type, status, fuel_level, mileage, last_service, next_service, license_plate, year, created_at, updated_at
		 FROM vehicles WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(&v.ID, &v.TenantID, &v.VehicleID, &v.Name, &v.Type, &v.Status, &v.FuelLevel, &v.Mileage, &v.LastService, &v.NextService, &v.LicensePlate, &v.Year, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle by id: %w", err)
	}
	return v, nil
}

// List returns a paginated list of vehicles, optionally filtered by status.
func (r *VehicleRepo) List(ctx context.Context, tenantID uuid.UUID, status string, page, perPage int) ([]models.Vehicle, int, error) {
	var total int
	if status != "" {
		err := r.db.QueryRow(ctx,
			`SELECT COUNT(*) FROM vehicles WHERE tenant_id = $1 AND status = $2`,
			tenantID, status,
		).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count vehicles: %w", err)
		}
	} else {
		err := r.db.QueryRow(ctx,
			`SELECT COUNT(*) FROM vehicles WHERE tenant_id = $1`,
			tenantID,
		).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count vehicles: %w", err)
		}
	}

	offset := (page - 1) * perPage
	var query string
	var args []interface{}
	if status != "" {
		query = `SELECT id, tenant_id, vehicle_id, name, type, status, fuel_level, mileage, last_service, next_service, license_plate, year, created_at, updated_at
				 FROM vehicles WHERE tenant_id = $1 AND status = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`
		args = []interface{}{tenantID, status, perPage, offset}
	} else {
		query = `SELECT id, tenant_id, vehicle_id, name, type, status, fuel_level, mileage, last_service, next_service, license_plate, year, created_at, updated_at
				 FROM vehicles WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = []interface{}{tenantID, perPage, offset}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []models.Vehicle
	for rows.Next() {
		var v models.Vehicle
		if err := rows.Scan(&v.ID, &v.TenantID, &v.VehicleID, &v.Name, &v.Type, &v.Status, &v.FuelLevel, &v.Mileage, &v.LastService, &v.NextService, &v.LicensePlate, &v.Year, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan vehicle: %w", err)
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, total, nil
}

// Update modifies an existing vehicle.
func (r *VehicleRepo) Update(ctx context.Context, tenantID, id uuid.UUID, v *models.Vehicle) error {
	_, err := r.db.Exec(ctx,
		`UPDATE vehicles SET vehicle_id = $1, name = $2, type = $3, status = $4, fuel_level = $5, mileage = $6, last_service = $7, next_service = $8, license_plate = $9, year = $10, updated_at = NOW()
		 WHERE id = $11 AND tenant_id = $12`,
		v.VehicleID, v.Name, v.Type, v.Status, v.FuelLevel, v.Mileage, v.LastService, v.NextService, v.LicensePlate, v.Year, id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}
	return nil
}

// Delete removes a vehicle by ID within a tenant.
func (r *VehicleRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM vehicles WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}
	return nil
}

// CountByStatus returns a count of vehicles grouped by status within a tenant.
func (r *VehicleRepo) CountByStatus(ctx context.Context, tenantID uuid.UUID) (map[string]int, error) {
	rows, err := r.db.Query(ctx,
		`SELECT status, COUNT(*) FROM vehicles WHERE tenant_id = $1 GROUP BY status`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to count vehicles by status: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan status count: %w", err)
		}
		counts[status] = count
	}
	return counts, nil
}

// GetActive returns all active vehicles for a tenant.
func (r *VehicleRepo) GetActive(ctx context.Context, tenantID uuid.UUID) ([]models.Vehicle, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, vehicle_id, name, type, status, fuel_level, mileage, last_service, next_service, license_plate, year, created_at, updated_at
		 FROM vehicles WHERE tenant_id = $1 AND status = 'available' ORDER BY license_plate ASC`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get active vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []models.Vehicle
	for rows.Next() {
		var v models.Vehicle
		if err := rows.Scan(&v.ID, &v.TenantID, &v.VehicleID, &v.Name, &v.Type, &v.Status, &v.FuelLevel, &v.Mileage, &v.LastService, &v.NextService, &v.LicensePlate, &v.Year, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}
