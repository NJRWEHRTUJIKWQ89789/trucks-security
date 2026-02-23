package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InventoryRepo handles database operations for inventory items.
type InventoryRepo struct {
	db *pgxpool.Pool
}

// NewInventoryRepo creates a new InventoryRepo instance.
func NewInventoryRepo(db *pgxpool.Pool) *InventoryRepo {
	return &InventoryRepo{db: db}
}

// Create inserts a new inventory item.
func (r *InventoryRepo) Create(ctx context.Context, i *models.InventoryItem) error {
	i.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO inventory_items (id, tenant_id, warehouse_id, sku, name, category, quantity, min_quantity, unit_price, weight, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())`,
		i.ID, i.TenantID, i.WarehouseID, i.SKU, i.Name, i.Category, i.Quantity, i.MinQuantity, i.UnitPrice, i.Weight, i.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to create inventory item: %w", err)
	}
	return nil
}

// GetByID retrieves an inventory item by ID within a tenant.
func (r *InventoryRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.InventoryItem, error) {
	i := &models.InventoryItem{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, warehouse_id, sku, name, category, quantity, min_quantity, unit_price, weight, status, created_at, updated_at
		 FROM inventory_items WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(&i.ID, &i.TenantID, &i.WarehouseID, &i.SKU, &i.Name, &i.Category, &i.Quantity, &i.MinQuantity, &i.UnitPrice, &i.Weight, &i.Status, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory item by id: %w", err)
	}
	return i, nil
}

// GetBySKU retrieves an inventory item by SKU within a tenant.
func (r *InventoryRepo) GetBySKU(ctx context.Context, tenantID uuid.UUID, sku string) (*models.InventoryItem, error) {
	i := &models.InventoryItem{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, warehouse_id, sku, name, category, quantity, min_quantity, unit_price, weight, status, created_at, updated_at
		 FROM inventory_items WHERE sku = $1 AND tenant_id = $2`,
		sku, tenantID,
	).Scan(&i.ID, &i.TenantID, &i.WarehouseID, &i.SKU, &i.Name, &i.Category, &i.Quantity, &i.MinQuantity, &i.UnitPrice, &i.Weight, &i.Status, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory item by sku: %w", err)
	}
	return i, nil
}

// List returns a paginated list of inventory items, optionally filtered by warehouse ID.
func (r *InventoryRepo) List(ctx context.Context, tenantID uuid.UUID, warehouseID *uuid.UUID, page, perPage int) ([]models.InventoryItem, int, error) {
	var total int
	if warehouseID != nil {
		err := r.db.QueryRow(ctx,
			`SELECT COUNT(*) FROM inventory_items WHERE tenant_id = $1 AND warehouse_id = $2`,
			tenantID, *warehouseID,
		).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count inventory items: %w", err)
		}
	} else {
		err := r.db.QueryRow(ctx,
			`SELECT COUNT(*) FROM inventory_items WHERE tenant_id = $1`,
			tenantID,
		).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count inventory items: %w", err)
		}
	}

	offset := (page - 1) * perPage
	var query string
	var args []interface{}
	if warehouseID != nil {
		query = `SELECT id, tenant_id, warehouse_id, sku, name, category, quantity, min_quantity, unit_price, weight, status, created_at, updated_at
				 FROM inventory_items WHERE tenant_id = $1 AND warehouse_id = $2 ORDER BY name ASC LIMIT $3 OFFSET $4`
		args = []interface{}{tenantID, *warehouseID, perPage, offset}
	} else {
		query = `SELECT id, tenant_id, warehouse_id, sku, name, category, quantity, min_quantity, unit_price, weight, status, created_at, updated_at
				 FROM inventory_items WHERE tenant_id = $1 ORDER BY name ASC LIMIT $2 OFFSET $3`
		args = []interface{}{tenantID, perPage, offset}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list inventory items: %w", err)
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var i models.InventoryItem
		if err := rows.Scan(&i.ID, &i.TenantID, &i.WarehouseID, &i.SKU, &i.Name, &i.Category, &i.Quantity, &i.MinQuantity, &i.UnitPrice, &i.Weight, &i.Status, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan inventory item: %w", err)
		}
		items = append(items, i)
	}
	return items, total, nil
}

// GetLowStock returns inventory items where quantity is at or below min_quantity within a tenant.
func (r *InventoryRepo) GetLowStock(ctx context.Context, tenantID uuid.UUID) ([]models.InventoryItem, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, warehouse_id, sku, name, category, quantity, min_quantity, unit_price, weight, status, created_at, updated_at
		 FROM inventory_items WHERE tenant_id = $1 AND quantity <= min_quantity ORDER BY quantity ASC`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock items: %w", err)
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var i models.InventoryItem
		if err := rows.Scan(&i.ID, &i.TenantID, &i.WarehouseID, &i.SKU, &i.Name, &i.Category, &i.Quantity, &i.MinQuantity, &i.UnitPrice, &i.Weight, &i.Status, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan low stock item: %w", err)
		}
		items = append(items, i)
	}
	return items, nil
}

// Update modifies an existing inventory item.
func (r *InventoryRepo) Update(ctx context.Context, tenantID, id uuid.UUID, i *models.InventoryItem) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE inventory_items SET warehouse_id = $1, sku = $2, name = $3, category = $4, quantity = $5, min_quantity = $6, unit_price = $7, weight = $8, status = $9, updated_at = NOW()
		 WHERE id = $10 AND tenant_id = $11`,
		i.WarehouseID, i.SKU, i.Name, i.Category, i.Quantity, i.MinQuantity, i.UnitPrice, i.Weight, i.Status, id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update inventory item: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("inventory item not found")
	}
	return nil
}

// Delete removes an inventory item by ID within a tenant.
func (r *InventoryRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	ct, err := r.db.Exec(ctx,
		`DELETE FROM inventory_items WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete inventory item: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("inventory item not found")
	}
	return nil
}

// Restock increases the quantity of an inventory item.
func (r *InventoryRepo) Restock(ctx context.Context, tenantID, id uuid.UUID, quantity int) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE inventory_items SET quantity = quantity + $1, status = CASE WHEN quantity + $1 > min_quantity THEN 'in_stock' ELSE status END, updated_at = NOW()
		 WHERE id = $2 AND tenant_id = $3`,
		quantity, id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to restock inventory item: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("inventory item not found")
	}
	return nil
}
