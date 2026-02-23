package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ShipmentRepo handles database operations for shipments.
type ShipmentRepo struct {
	db *pgxpool.Pool
}

// NewShipmentRepo creates a new ShipmentRepo instance.
func NewShipmentRepo(db *pgxpool.Pool) *ShipmentRepo {
	return &ShipmentRepo{db: db}
}

// Create inserts a new shipment.
func (r *ShipmentRepo) Create(ctx context.Context, s *models.Shipment) error {
	s.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO shipments (id, tenant_id, tracking_number, origin, destination, status, carrier, weight, dimensions, estimated_delivery, actual_delivery, customer_name, customer_email, notes, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW(), NOW())`,
		s.ID, s.TenantID, s.TrackingNumber, s.Origin, s.Destination, s.Status, s.Carrier, s.Weight, s.Dimensions, s.EstimatedDelivery, s.ActualDelivery, s.CustomerName, s.CustomerEmail, s.Notes,
	)
	if err != nil {
		return fmt.Errorf("failed to create shipment: %w", err)
	}
	return nil
}

// GetByID retrieves a shipment by ID within a tenant.
func (r *ShipmentRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Shipment, error) {
	s := &models.Shipment{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, tracking_number, origin, destination, status, carrier, weight, dimensions, estimated_delivery, actual_delivery, customer_name, customer_email, notes, created_at, updated_at
		 FROM shipments WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(&s.ID, &s.TenantID, &s.TrackingNumber, &s.Origin, &s.Destination, &s.Status, &s.Carrier, &s.Weight, &s.Dimensions, &s.EstimatedDelivery, &s.ActualDelivery, &s.CustomerName, &s.CustomerEmail, &s.Notes, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment by id: %w", err)
	}
	return s, nil
}

// GetByTracking retrieves a shipment by tracking number within a tenant.
func (r *ShipmentRepo) GetByTracking(ctx context.Context, tenantID uuid.UUID, trackingNumber string) (*models.Shipment, error) {
	s := &models.Shipment{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, tracking_number, origin, destination, status, carrier, weight, dimensions, estimated_delivery, actual_delivery, customer_name, customer_email, notes, created_at, updated_at
		 FROM shipments WHERE tracking_number = $1 AND tenant_id = $2`,
		trackingNumber, tenantID,
	).Scan(&s.ID, &s.TenantID, &s.TrackingNumber, &s.Origin, &s.Destination, &s.Status, &s.Carrier, &s.Weight, &s.Dimensions, &s.EstimatedDelivery, &s.ActualDelivery, &s.CustomerName, &s.CustomerEmail, &s.Notes, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment by tracking number: %w", err)
	}
	return s, nil
}

// List returns a paginated list of shipments, optionally filtered by status.
func (r *ShipmentRepo) List(ctx context.Context, tenantID uuid.UUID, status string, page, perPage int) ([]models.Shipment, int, error) {
	var total int
	if status != "" {
		err := r.db.QueryRow(ctx,
			`SELECT COUNT(*) FROM shipments WHERE tenant_id = $1 AND status = $2`,
			tenantID, status,
		).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count shipments: %w", err)
		}
	} else {
		err := r.db.QueryRow(ctx,
			`SELECT COUNT(*) FROM shipments WHERE tenant_id = $1`,
			tenantID,
		).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count shipments: %w", err)
		}
	}

	offset := (page - 1) * perPage
	var query string
	var args []interface{}
	if status != "" {
		query = `SELECT id, tenant_id, tracking_number, origin, destination, status, carrier, weight, dimensions, estimated_delivery, actual_delivery, customer_name, customer_email, notes, created_at, updated_at
				 FROM shipments WHERE tenant_id = $1 AND status = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`
		args = []interface{}{tenantID, status, perPage, offset}
	} else {
		query = `SELECT id, tenant_id, tracking_number, origin, destination, status, carrier, weight, dimensions, estimated_delivery, actual_delivery, customer_name, customer_email, notes, created_at, updated_at
				 FROM shipments WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = []interface{}{tenantID, perPage, offset}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list shipments: %w", err)
	}
	defer rows.Close()

	var shipments []models.Shipment
	for rows.Next() {
		var s models.Shipment
		if err := rows.Scan(&s.ID, &s.TenantID, &s.TrackingNumber, &s.Origin, &s.Destination, &s.Status, &s.Carrier, &s.Weight, &s.Dimensions, &s.EstimatedDelivery, &s.ActualDelivery, &s.CustomerName, &s.CustomerEmail, &s.Notes, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan shipment: %w", err)
		}
		shipments = append(shipments, s)
	}
	return shipments, total, nil
}

// GetDelayed returns shipments that are past their estimated delivery date and not yet delivered.
func (r *ShipmentRepo) GetDelayed(ctx context.Context, tenantID uuid.UUID) ([]models.Shipment, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, tracking_number, origin, destination, status, carrier, weight, dimensions, estimated_delivery, actual_delivery, customer_name, customer_email, notes, created_at, updated_at
		 FROM shipments WHERE tenant_id = $1 AND status != 'delivered' AND estimated_delivery < NOW() AND estimated_delivery IS NOT NULL
		 ORDER BY estimated_delivery ASC`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get delayed shipments: %w", err)
	}
	defer rows.Close()

	var shipments []models.Shipment
	for rows.Next() {
		var s models.Shipment
		if err := rows.Scan(&s.ID, &s.TenantID, &s.TrackingNumber, &s.Origin, &s.Destination, &s.Status, &s.Carrier, &s.Weight, &s.Dimensions, &s.EstimatedDelivery, &s.ActualDelivery, &s.CustomerName, &s.CustomerEmail, &s.Notes, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan delayed shipment: %w", err)
		}
		shipments = append(shipments, s)
	}
	return shipments, nil
}

// Update modifies an existing shipment.
func (r *ShipmentRepo) Update(ctx context.Context, tenantID, id uuid.UUID, s *models.Shipment) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE shipments SET tracking_number = $1, origin = $2, destination = $3, status = $4, carrier = $5, weight = $6, dimensions = $7, estimated_delivery = $8, actual_delivery = $9, customer_name = $10, customer_email = $11, notes = $12, updated_at = NOW()
		 WHERE id = $13 AND tenant_id = $14`,
		s.TrackingNumber, s.Origin, s.Destination, s.Status, s.Carrier, s.Weight, s.Dimensions, s.EstimatedDelivery, s.ActualDelivery, s.CustomerName, s.CustomerEmail, s.Notes, id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update shipment: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("shipment not found or access denied")
	}
	return nil
}

// Delete removes a shipment by ID within a tenant.
func (r *ShipmentRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	ct, err := r.db.Exec(ctx,
		`DELETE FROM shipments WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete shipment: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("shipment not found or access denied")
	}
	return nil
}

// CountByStatus returns a count of shipments grouped by status within a tenant.
func (r *ShipmentRepo) CountByStatus(ctx context.Context, tenantID uuid.UUID) (map[string]int, error) {
	rows, err := r.db.Query(ctx,
		`SELECT status, COUNT(*) FROM shipments WHERE tenant_id = $1 GROUP BY status`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to count shipments by status: %w", err)
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

// GetDeliveredToday returns the count of shipments delivered today within a tenant.
func (r *ShipmentRepo) GetDeliveredToday(ctx context.Context, tenantID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM shipments WHERE tenant_id = $1 AND status = 'delivered' AND actual_delivery >= CURRENT_DATE AND actual_delivery < CURRENT_DATE + INTERVAL '1 day'`,
		tenantID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count delivered today: %w", err)
	}
	return count, nil
}
