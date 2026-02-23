package repository

import (
	"context"
	"errors"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrOrderNotFound is returned when an order does not exist within the tenant.
var ErrOrderNotFound = errors.New("order not found")

// OrderRepo handles database operations for orders.
type OrderRepo struct {
	db *pgxpool.Pool
}

// NewOrderRepo creates a new OrderRepo instance.
func NewOrderRepo(db *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{db: db}
}

// Create inserts a new order.
func (r *OrderRepo) Create(ctx context.Context, o *models.Order) error {
	o.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO orders (id, tenant_id, order_number, customer_name, customer_email, status, type, total_amount, shipment_id, scheduled_date, return_reason, cancellation_reason, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())`,
		o.ID, o.TenantID, o.OrderNumber, o.CustomerName, o.CustomerEmail, o.Status, o.Type, o.TotalAmount, o.ShipmentID, o.ScheduledDate, o.ReturnReason, o.CancellationReason,
	)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

// GetByID retrieves an order by ID within a tenant.
func (r *OrderRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Order, error) {
	o := &models.Order{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, order_number, customer_name, customer_email, status, type, total_amount, shipment_id, scheduled_date, return_reason, cancellation_reason, created_at, updated_at
		 FROM orders WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(&o.ID, &o.TenantID, &o.OrderNumber, &o.CustomerName, &o.CustomerEmail, &o.Status, &o.Type, &o.TotalAmount, &o.ShipmentID, &o.ScheduledDate, &o.ReturnReason, &o.CancellationReason, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by id: %w", err)
	}
	return o, nil
}

// List returns a paginated list of orders, optionally filtered by status.
func (r *OrderRepo) List(ctx context.Context, tenantID uuid.UUID, status string, page, perPage int) ([]models.Order, int, error) {
	var total int
	if status != "" {
		err := r.db.QueryRow(ctx,
			`SELECT COUNT(*) FROM orders WHERE tenant_id = $1 AND status = $2`,
			tenantID, status,
		).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count orders: %w", err)
		}
	} else {
		err := r.db.QueryRow(ctx,
			`SELECT COUNT(*) FROM orders WHERE tenant_id = $1`,
			tenantID,
		).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count orders: %w", err)
		}
	}

	offset := (page - 1) * perPage
	var query string
	var args []interface{}
	if status != "" {
		query = `SELECT id, tenant_id, order_number, customer_name, customer_email, status, type, total_amount, shipment_id, scheduled_date, return_reason, cancellation_reason, created_at, updated_at
				 FROM orders WHERE tenant_id = $1 AND status = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`
		args = []interface{}{tenantID, status, perPage, offset}
	} else {
		query = `SELECT id, tenant_id, order_number, customer_name, customer_email, status, type, total_amount, shipment_id, scheduled_date, return_reason, cancellation_reason, created_at, updated_at
				 FROM orders WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = []interface{}{tenantID, perPage, offset}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.TenantID, &o.OrderNumber, &o.CustomerName, &o.CustomerEmail, &o.Status, &o.Type, &o.TotalAmount, &o.ShipmentID, &o.ScheduledDate, &o.ReturnReason, &o.CancellationReason, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, o)
	}
	return orders, total, nil
}

// GetScheduled returns orders with a scheduled_date in the future within a tenant.
func (r *OrderRepo) GetScheduled(ctx context.Context, tenantID uuid.UUID) ([]models.Order, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, order_number, customer_name, customer_email, status, type, total_amount, shipment_id, scheduled_date, return_reason, cancellation_reason, created_at, updated_at
		 FROM orders WHERE tenant_id = $1 AND scheduled_date IS NOT NULL AND scheduled_date >= NOW() AND status NOT IN ('cancelled', 'returned')
		 ORDER BY scheduled_date ASC`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get scheduled orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.TenantID, &o.OrderNumber, &o.CustomerName, &o.CustomerEmail, &o.Status, &o.Type, &o.TotalAmount, &o.ShipmentID, &o.ScheduledDate, &o.ReturnReason, &o.CancellationReason, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan scheduled order: %w", err)
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// GetReturns returns orders with status 'returned' within a tenant.
func (r *OrderRepo) GetReturns(ctx context.Context, tenantID uuid.UUID) ([]models.Order, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, order_number, customer_name, customer_email, status, type, total_amount, shipment_id, scheduled_date, return_reason, cancellation_reason, created_at, updated_at
		 FROM orders WHERE tenant_id = $1 AND status = 'returned' ORDER BY updated_at DESC`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get returned orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.TenantID, &o.OrderNumber, &o.CustomerName, &o.CustomerEmail, &o.Status, &o.Type, &o.TotalAmount, &o.ShipmentID, &o.ScheduledDate, &o.ReturnReason, &o.CancellationReason, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan returned order: %w", err)
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// GetCancelled returns orders with status 'cancelled' within a tenant.
func (r *OrderRepo) GetCancelled(ctx context.Context, tenantID uuid.UUID) ([]models.Order, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, order_number, customer_name, customer_email, status, type, total_amount, shipment_id, scheduled_date, return_reason, cancellation_reason, created_at, updated_at
		 FROM orders WHERE tenant_id = $1 AND status = 'cancelled' ORDER BY updated_at DESC`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get cancelled orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.TenantID, &o.OrderNumber, &o.CustomerName, &o.CustomerEmail, &o.Status, &o.Type, &o.TotalAmount, &o.ShipmentID, &o.ScheduledDate, &o.ReturnReason, &o.CancellationReason, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan cancelled order: %w", err)
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// Update modifies an existing order.
func (r *OrderRepo) Update(ctx context.Context, tenantID, id uuid.UUID, o *models.Order) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE orders SET order_number = $1, customer_name = $2, customer_email = $3, status = $4, type = $5, total_amount = $6, shipment_id = $7, scheduled_date = $8, return_reason = $9, cancellation_reason = $10, updated_at = NOW()
		 WHERE id = $11 AND tenant_id = $12`,
		o.OrderNumber, o.CustomerName, o.CustomerEmail, o.Status, o.Type, o.TotalAmount, o.ShipmentID, o.ScheduledDate, o.ReturnReason, o.CancellationReason, id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrOrderNotFound
	}
	return nil
}

// CancelOrder sets an order's status to 'cancelled' with a reason.
func (r *OrderRepo) CancelOrder(ctx context.Context, tenantID, id uuid.UUID, reason string) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE orders SET status = 'cancelled', cancellation_reason = $1, updated_at = NOW() WHERE id = $2 AND tenant_id = $3`,
		reason, id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrOrderNotFound
	}
	return nil
}

// ReturnOrder sets an order's status to 'returned' with a reason.
func (r *OrderRepo) ReturnOrder(ctx context.Context, tenantID, id uuid.UUID, reason string) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE orders SET status = 'returned', return_reason = $1, updated_at = NOW() WHERE id = $2 AND tenant_id = $3`,
		reason, id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to return order: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrOrderNotFound
	}
	return nil
}

// Delete removes an order by ID within a tenant.
func (r *OrderRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	ct, err := r.db.Exec(ctx,
		`DELETE FROM orders WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrOrderNotFound
	}
	return nil
}

// CountByStatus returns a count of orders grouped by status within a tenant.
func (r *OrderRepo) CountByStatus(ctx context.Context, tenantID uuid.UUID) (map[string]int, error) {
	rows, err := r.db.Query(ctx,
		`SELECT status, COUNT(*) FROM orders WHERE tenant_id = $1 GROUP BY status`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to count orders by status: %w", err)
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
