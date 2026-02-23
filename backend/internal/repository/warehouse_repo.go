package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WarehouseRepo handles database operations for warehouses.
type WarehouseRepo struct {
	db *pgxpool.Pool
}

// NewWarehouseRepo creates a new WarehouseRepo instance.
func NewWarehouseRepo(db *pgxpool.Pool) *WarehouseRepo {
	return &WarehouseRepo{db: db}
}

// Create inserts a new warehouse.
func (r *WarehouseRepo) Create(ctx context.Context, w *models.Warehouse) error {
	w.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO warehouses (id, tenant_id, name, location, address, capacity, used_capacity, manager, phone, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())`,
		w.ID, w.TenantID, w.Name, w.Location, w.Address, w.Capacity, w.UsedCapacity, w.Manager, w.Phone, w.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to create warehouse: %w", err)
	}
	return nil
}

// GetByID retrieves a warehouse by ID within a tenant.
func (r *WarehouseRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Warehouse, error) {
	w := &models.Warehouse{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, name, location, address, capacity, used_capacity, manager, phone, status, created_at, updated_at
		 FROM warehouses WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(&w.ID, &w.TenantID, &w.Name, &w.Location, &w.Address, &w.Capacity, &w.UsedCapacity, &w.Manager, &w.Phone, &w.Status, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get warehouse by id: %w", err)
	}
	return w, nil
}

// List returns a paginated list of warehouses within a tenant.
func (r *WarehouseRepo) List(ctx context.Context, tenantID uuid.UUID, page, perPage int) ([]models.Warehouse, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM warehouses WHERE tenant_id = $1`,
		tenantID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count warehouses: %w", err)
	}

	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, name, location, address, capacity, used_capacity, manager, phone, status, created_at, updated_at
		 FROM warehouses WHERE tenant_id = $1 ORDER BY name ASC LIMIT $2 OFFSET $3`,
		tenantID, perPage, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list warehouses: %w", err)
	}
	defer rows.Close()

	var warehouses []models.Warehouse
	for rows.Next() {
		var w models.Warehouse
		if err := rows.Scan(&w.ID, &w.TenantID, &w.Name, &w.Location, &w.Address, &w.Capacity, &w.UsedCapacity, &w.Manager, &w.Phone, &w.Status, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan warehouse: %w", err)
		}
		warehouses = append(warehouses, w)
	}
	return warehouses, total, nil
}

// Update modifies an existing warehouse.
func (r *WarehouseRepo) Update(ctx context.Context, tenantID, id uuid.UUID, w *models.Warehouse) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE warehouses SET name = $1, location = $2, address = $3, capacity = $4, used_capacity = $5, manager = $6, phone = $7, status = $8, updated_at = NOW()
		 WHERE id = $9 AND tenant_id = $10`,
		w.Name, w.Location, w.Address, w.Capacity, w.UsedCapacity, w.Manager, w.Phone, w.Status, id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update warehouse: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("warehouse not found")
	}
	return nil
}

// Delete removes a warehouse by ID within a tenant.
func (r *WarehouseRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	ct, err := r.db.Exec(ctx,
		`DELETE FROM warehouses WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete warehouse: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("warehouse not found")
	}
	return nil
}

// GetCapacityStats returns aggregated capacity statistics for all warehouses within a tenant.
func (r *WarehouseRepo) GetCapacityStats(ctx context.Context, tenantID uuid.UUID) (*models.WarehouseCapacityStats, error) {
	stats := &models.WarehouseCapacityStats{}
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(capacity), 0), COALESCE(SUM(used_capacity), 0)
		 FROM warehouses WHERE tenant_id = $1 AND status = 'active'`,
		tenantID,
	).Scan(&stats.TotalCapacity, &stats.UsedCapacity)
	if err != nil {
		return nil, fmt.Errorf("failed to get capacity stats: %w", err)
	}
	stats.FreeCapacity = stats.TotalCapacity - stats.UsedCapacity
	if stats.TotalCapacity > 0 {
		stats.Utilization = float64(stats.UsedCapacity) / float64(stats.TotalCapacity) * 100
	}
	return stats, nil
}
