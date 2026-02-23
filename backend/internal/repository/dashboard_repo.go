package repository

import (
	"context"
	"fmt"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DashboardRepo handles aggregate dashboard statistics.
type DashboardRepo struct {
	db *pgxpool.Pool
}

// NewDashboardRepo creates a new DashboardRepo instance.
func NewDashboardRepo(db *pgxpool.Pool) *DashboardRepo {
	return &DashboardRepo{db: db}
}

// GetStats returns aggregated dashboard statistics across tables for a tenant.
func (r *DashboardRepo) GetStats(ctx context.Context, tenantID uuid.UUID) (*models.DashboardStats, error) {
	stats := &models.DashboardStats{}

	// Shipment stats
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*),
			COALESCE(SUM(CASE WHEN status IN ('in_transit', 'pending', 'processing') THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'delivered' AND actual_delivery >= CURRENT_DATE AND actual_delivery < CURRENT_DATE + INTERVAL '1 day' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status != 'delivered' AND estimated_delivery < NOW() AND estimated_delivery IS NOT NULL THEN 1 ELSE 0 END), 0)
		 FROM shipments WHERE tenant_id = $1`,
		tenantID,
	).Scan(&stats.TotalShipments, &stats.ActiveShipments, &stats.DeliveredToday, &stats.DelayedShipments)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment stats: %w", err)
	}

	// Vehicle stats
	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*),
			COALESCE(SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END), 0)
		 FROM vehicles WHERE tenant_id = $1`,
		tenantID,
	).Scan(&stats.TotalVehicles, &stats.ActiveVehicles)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle stats: %w", err)
	}

	// Driver stats
	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*),
			COALESCE(SUM(CASE WHEN status = 'available' THEN 1 ELSE 0 END), 0)
		 FROM drivers WHERE tenant_id = $1`,
		tenantID,
	).Scan(&stats.TotalDrivers, &stats.AvailableDrivers)
	if err != nil {
		return nil, fmt.Errorf("failed to get driver stats: %w", err)
	}

	// Order stats
	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*),
			COALESCE(SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END), 0)
		 FROM orders WHERE tenant_id = $1`,
		tenantID,
	).Scan(&stats.TotalOrders, &stats.PendingOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to get order stats: %w", err)
	}

	// Warehouse stats
	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM warehouses WHERE tenant_id = $1`,
		tenantID,
	).Scan(&stats.TotalWarehouses)
	if err != nil {
		return nil, fmt.Errorf("failed to get warehouse stats: %w", err)
	}

	// Inventory stats
	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*),
			COALESCE(SUM(CASE WHEN quantity <= min_quantity THEN 1 ELSE 0 END), 0)
		 FROM inventory_items WHERE tenant_id = $1`,
		tenantID,
	).Scan(&stats.TotalInventory, &stats.LowStockItems)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory stats: %w", err)
	}

	// Client stats
	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM clients WHERE tenant_id = $1`,
		tenantID,
	).Scan(&stats.TotalClients)
	if err != nil {
		return nil, fmt.Errorf("failed to get client stats: %w", err)
	}

	// Revenue
	var totalRevenue *float64
	err = r.db.QueryRow(ctx,
		`SELECT SUM(total_amount) FROM orders WHERE tenant_id = $1 AND status NOT IN ('cancelled', 'returned')`,
		tenantID,
	).Scan(&totalRevenue)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue stats: %w", err)
	}
	if totalRevenue != nil {
		stats.TotalRevenue = *totalRevenue
	}

	// Maintenance due
	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM maintenance_records WHERE tenant_id = $1 AND status = 'scheduled' AND scheduled_date <= NOW() + INTERVAL '7 days'`,
		tenantID,
	).Scan(&stats.MaintenanceDue)
	if err != nil {
		return nil, fmt.Errorf("failed to get maintenance stats: %w", err)
	}

	return stats, nil
}

// GetPerformance returns performance metrics for a tenant's dashboard.
func (r *DashboardRepo) GetPerformance(ctx context.Context, tenantID uuid.UUID) (*models.Performance, error) {
	perf := &models.Performance{}

	// On-time delivery rate: delivered shipments where actual_delivery <= estimated_delivery
	var totalDelivered, onTime int
	err := r.db.QueryRow(ctx,
		`SELECT
			COALESCE(SUM(CASE WHEN status = 'delivered' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'delivered' AND actual_delivery <= estimated_delivery THEN 1 ELSE 0 END), 0)
		 FROM shipments WHERE tenant_id = $1 AND estimated_delivery IS NOT NULL`,
		tenantID,
	).Scan(&totalDelivered, &onTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery rate: %w", err)
	}
	if totalDelivered > 0 {
		perf.OnTimeDeliveryRate = float64(onTime) / float64(totalDelivered) * 100
	}

	// Average delivery time in hours
	var avgHours *float64
	err = r.db.QueryRow(ctx,
		`SELECT AVG(EXTRACT(EPOCH FROM (actual_delivery - created_at)) / 3600)
		 FROM shipments WHERE tenant_id = $1 AND status = 'delivered' AND actual_delivery IS NOT NULL`,
		tenantID,
	).Scan(&avgHours)
	if err != nil {
		return nil, fmt.Errorf("failed to get average delivery time: %w", err)
	}
	if avgHours != nil {
		perf.AverageDeliveryTime = *avgHours
	}

	// Fleet utilization: active vehicles / total vehicles
	var totalVehicles, activeVehicles int
	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*),
			COALESCE(SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END), 0)
		 FROM vehicles WHERE tenant_id = $1`,
		tenantID,
	).Scan(&totalVehicles, &activeVehicles)
	if err != nil {
		return nil, fmt.Errorf("failed to get fleet utilization: %w", err)
	}
	if totalVehicles > 0 {
		perf.FleetUtilization = float64(activeVehicles) / float64(totalVehicles) * 100
	}

	// Order fulfillment rate: fulfilled orders / total orders
	var totalOrders, fulfilledOrders int
	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*),
			COALESCE(SUM(CASE WHEN status IN ('delivered', 'completed', 'fulfilled') THEN 1 ELSE 0 END), 0)
		 FROM orders WHERE tenant_id = $1`,
		tenantID,
	).Scan(&totalOrders, &fulfilledOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to get order fulfillment rate: %w", err)
	}
	if totalOrders > 0 {
		perf.OrderFulfillmentRate = float64(fulfilledOrders) / float64(totalOrders) * 100
	}

	return perf, nil
}
