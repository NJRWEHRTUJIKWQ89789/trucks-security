package repository

import (
	"context"
	"fmt"
	"time"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ReportRepo handles report generation queries.
type ReportRepo struct {
	db *pgxpool.Pool
}

// NewReportRepo creates a new ReportRepo instance.
func NewReportRepo(db *pgxpool.Pool) *ReportRepo {
	return &ReportRepo{db: db}
}

// GetRevenueReport returns revenue report data with monthly breakdowns for the specified year.
func (r *ReportRepo) GetRevenueReport(ctx context.Context, tenantID uuid.UUID, year int) (*models.RevenueReport, error) {
	report := &models.RevenueReport{}
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)

	// Total revenue for the specified year
	var totalRev *float64
	err := r.db.QueryRow(ctx,
		`SELECT SUM(total_amount) FROM orders WHERE tenant_id = $1 AND status NOT IN ('cancelled', 'returned') AND created_at >= $2 AND created_at < $3`,
		tenantID, startDate, endDate,
	).Scan(&totalRev)
	if err != nil {
		return nil, fmt.Errorf("failed to get total revenue: %w", err)
	}
	if totalRev != nil {
		report.TotalRevenue = *totalRev
	}

	// Monthly breakdown for the specified year
	rows, err := r.db.Query(ctx,
		`SELECT TO_CHAR(created_at, 'YYYY-MM') AS month, COALESCE(SUM(total_amount), 0) AS value, COUNT(*) AS count
		 FROM orders WHERE tenant_id = $1 AND status NOT IN ('cancelled', 'returned') AND created_at >= $2 AND created_at < $3
		 GROUP BY month ORDER BY month ASC`,
		tenantID, startDate, endDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue monthly breakdown: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mb models.MonthlyBreakdown
		if err := rows.Scan(&mb.Month, &mb.Value, &mb.Count); err != nil {
			return nil, fmt.Errorf("failed to scan revenue breakdown: %w", err)
		}
		report.MonthlyBreakdown = append(report.MonthlyBreakdown, mb)
	}
	return report, nil
}

// GetDeliveryReport returns delivery performance data with monthly breakdowns for the specified year.
func (r *ReportRepo) GetDeliveryReport(ctx context.Context, tenantID uuid.UUID, year int) (*models.DeliveryReport, error) {
	report := &models.DeliveryReport{}
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)

	err := r.db.QueryRow(ctx,
		`SELECT
			COUNT(*),
			COALESCE(SUM(CASE WHEN actual_delivery IS NOT NULL AND estimated_delivery IS NOT NULL AND actual_delivery <= estimated_delivery THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN actual_delivery IS NOT NULL AND estimated_delivery IS NOT NULL AND actual_delivery > estimated_delivery THEN 1 ELSE 0 END), 0)
		 FROM shipments WHERE tenant_id = $1 AND status = 'delivered' AND actual_delivery >= $2 AND actual_delivery < $3`,
		tenantID, startDate, endDate,
	).Scan(&report.TotalDeliveries, &report.OnTimeDeliveries, &report.DelayedDeliveries)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery stats: %w", err)
	}

	if report.TotalDeliveries > 0 {
		report.OnTimeRate = float64(report.OnTimeDeliveries) / float64(report.TotalDeliveries) * 100
	}

	// Monthly breakdown for the specified year
	rows, err := r.db.Query(ctx,
		`SELECT TO_CHAR(actual_delivery, 'YYYY-MM') AS month, COUNT(*) AS count, 0 AS value
		 FROM shipments WHERE tenant_id = $1 AND status = 'delivered' AND actual_delivery IS NOT NULL AND actual_delivery >= $2 AND actual_delivery < $3
		 GROUP BY month ORDER BY month ASC`,
		tenantID, startDate, endDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery monthly breakdown: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mb models.MonthlyBreakdown
		if err := rows.Scan(&mb.Month, &mb.Count, &mb.Value); err != nil {
			return nil, fmt.Errorf("failed to scan delivery breakdown: %w", err)
		}
		report.MonthlyBreakdown = append(report.MonthlyBreakdown, mb)
	}
	return report, nil
}

// GetFleetReport returns fleet utilization data with monthly breakdowns for the specified year.
func (r *ReportRepo) GetFleetReport(ctx context.Context, tenantID uuid.UUID, year int) (*models.FleetReport, error) {
	report := &models.FleetReport{}
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)

	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*),
			COALESCE(SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END), 0)
		 FROM vehicles WHERE tenant_id = $1`,
		tenantID,
	).Scan(&report.TotalVehicles, &report.ActiveVehicles)
	if err != nil {
		return nil, fmt.Errorf("failed to get fleet vehicle stats: %w", err)
	}

	// Maintenance count and cost for the specified year
	var totalCost *float64
	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*), SUM(cost) FROM maintenance_records WHERE tenant_id = $1 AND scheduled_date >= $2 AND scheduled_date < $3`,
		tenantID, startDate, endDate,
	).Scan(&report.MaintenanceCount, &totalCost)
	if err != nil {
		return nil, fmt.Errorf("failed to get fleet maintenance stats: %w", err)
	}
	if totalCost != nil {
		report.TotalCost = *totalCost
	}

	// Monthly maintenance cost breakdown for the specified year
	rows, err := r.db.Query(ctx,
		`SELECT TO_CHAR(scheduled_date, 'YYYY-MM') AS month, COALESCE(SUM(cost), 0) AS value, COUNT(*) AS count
		 FROM maintenance_records WHERE tenant_id = $1 AND scheduled_date IS NOT NULL AND scheduled_date >= $2 AND scheduled_date < $3
		 GROUP BY month ORDER BY month ASC`,
		tenantID, startDate, endDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get fleet monthly breakdown: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mb models.MonthlyBreakdown
		if err := rows.Scan(&mb.Month, &mb.Value, &mb.Count); err != nil {
			return nil, fmt.Errorf("failed to scan fleet breakdown: %w", err)
		}
		report.MonthlyBreakdown = append(report.MonthlyBreakdown, mb)
	}
	return report, nil
}
