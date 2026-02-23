package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DriverRepo handles database operations for drivers.
type DriverRepo struct {
	db *pgxpool.Pool
}

// NewDriverRepo creates a new DriverRepo instance.
func NewDriverRepo(db *pgxpool.Pool) *DriverRepo {
	return &DriverRepo{db: db}
}

// Create inserts a new driver.
func (r *DriverRepo) Create(ctx context.Context, d *models.Driver) error {
	d.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO drivers (id, tenant_id, employee_id, first_name, last_name, email, phone, license_number, license_expiry, status, rating, total_deliveries, vehicle_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW())`,
		d.ID, d.TenantID, d.EmployeeID, d.FirstName, d.LastName, d.Email, d.Phone, d.LicenseNumber, d.LicenseExpiry, d.Status, d.Rating, d.TotalDeliveries, d.VehicleID,
	)
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}
	return nil
}

// GetByID retrieves a driver by ID within a tenant.
func (r *DriverRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Driver, error) {
	d := &models.Driver{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, employee_id, first_name, last_name, email, phone, license_number, license_expiry, status, rating, total_deliveries, vehicle_id, created_at, updated_at
		 FROM drivers WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(&d.ID, &d.TenantID, &d.EmployeeID, &d.FirstName, &d.LastName, &d.Email, &d.Phone, &d.LicenseNumber, &d.LicenseExpiry, &d.Status, &d.Rating, &d.TotalDeliveries, &d.VehicleID, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get driver by id: %w", err)
	}
	return d, nil
}

// List returns a paginated list of drivers within a tenant.
func (r *DriverRepo) List(ctx context.Context, tenantID uuid.UUID, page, perPage int) ([]models.Driver, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM drivers WHERE tenant_id = $1`,
		tenantID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count drivers: %w", err)
	}

	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, employee_id, first_name, last_name, email, phone, license_number, license_expiry, status, rating, total_deliveries, vehicle_id, created_at, updated_at
		 FROM drivers WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		tenantID, perPage, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list drivers: %w", err)
	}
	defer rows.Close()

	var drivers []models.Driver
	for rows.Next() {
		var d models.Driver
		if err := rows.Scan(&d.ID, &d.TenantID, &d.EmployeeID, &d.FirstName, &d.LastName, &d.Email, &d.Phone, &d.LicenseNumber, &d.LicenseExpiry, &d.Status, &d.Rating, &d.TotalDeliveries, &d.VehicleID, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan driver: %w", err)
		}
		drivers = append(drivers, d)
	}
	return drivers, total, nil
}

// Update modifies an existing driver.
func (r *DriverRepo) Update(ctx context.Context, tenantID, id uuid.UUID, d *models.Driver) error {
	_, err := r.db.Exec(ctx,
		`UPDATE drivers SET employee_id = $1, first_name = $2, last_name = $3, email = $4, phone = $5, license_number = $6, license_expiry = $7, status = $8, rating = $9, total_deliveries = $10, vehicle_id = $11, updated_at = NOW()
		 WHERE id = $12 AND tenant_id = $13`,
		d.EmployeeID, d.FirstName, d.LastName, d.Email, d.Phone, d.LicenseNumber, d.LicenseExpiry, d.Status, d.Rating, d.TotalDeliveries, d.VehicleID, id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update driver: %w", err)
	}
	return nil
}

// Delete removes a driver by ID within a tenant.
func (r *DriverRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM drivers WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete driver: %w", err)
	}
	return nil
}

// GetAvailable returns all drivers with status 'available' within a tenant.
func (r *DriverRepo) GetAvailable(ctx context.Context, tenantID uuid.UUID) ([]models.Driver, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, employee_id, first_name, last_name, email, phone, license_number, license_expiry, status, rating, total_deliveries, vehicle_id, created_at, updated_at
		 FROM drivers WHERE tenant_id = $1 AND status = 'available' ORDER BY last_name ASC`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get available drivers: %w", err)
	}
	defer rows.Close()

	var drivers []models.Driver
	for rows.Next() {
		var d models.Driver
		if err := rows.Scan(&d.ID, &d.TenantID, &d.EmployeeID, &d.FirstName, &d.LastName, &d.Email, &d.Phone, &d.LicenseNumber, &d.LicenseExpiry, &d.Status, &d.Rating, &d.TotalDeliveries, &d.VehicleID, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan available driver: %w", err)
		}
		drivers = append(drivers, d)
	}
	return drivers, nil
}

// GetByPhone retrieves a driver by phone number (for mobile login).
func (r *DriverRepo) GetByPhone(ctx context.Context, phone string) (*models.Driver, error) {
	d := &models.Driver{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, employee_id, first_name, last_name, email, phone, license_number, license_expiry, status, rating, total_deliveries, vehicle_id, pin_hash, created_at, updated_at
		 FROM drivers WHERE phone = $1`,
		phone,
	).Scan(&d.ID, &d.TenantID, &d.EmployeeID, &d.FirstName, &d.LastName, &d.Email, &d.Phone, &d.LicenseNumber, &d.LicenseExpiry, &d.Status, &d.Rating, &d.TotalDeliveries, &d.VehicleID, &d.PinHash, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get driver by phone: %w", err)
	}
	return d, nil
}
