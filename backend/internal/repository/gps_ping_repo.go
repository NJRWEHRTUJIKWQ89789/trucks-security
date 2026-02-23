package repository

import (
	"context"
	"fmt"
	"time"

	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GPSPingRepo struct {
	db *pgxpool.Pool
}

func NewGPSPingRepo(db *pgxpool.Pool) *GPSPingRepo {
	return &GPSPingRepo{db: db}
}

// BulkInsert inserts a batch of GPS pings in a single query.
func (r *GPSPingRepo) BulkInsert(ctx context.Context, pings []models.GPSPing) (int, error) {
	if len(pings) == 0 {
		return 0, nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	count := 0
	for _, p := range pings {
		_, err := tx.Exec(ctx,
			`INSERT INTO gps_pings (tenant_id, driver_id, truck_id, shift_id, latitude, longitude, speed_kmh, heading, accuracy, battery_level, is_moving, recorded_at, received_at, is_delayed)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
			p.TenantID, p.DriverID, p.TruckID, p.ShiftID, p.Latitude, p.Longitude, p.SpeedKmh, p.Heading, p.Accuracy, p.BatteryLevel, p.IsMoving, p.RecordedAt, p.ReceivedAt, p.IsDelayed,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to insert ping: %w", err)
		}
		count++
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return count, nil
}

// GetLatestByDriver returns the most recent ping for a driver.
func (r *GPSPingRepo) GetLatestByDriver(ctx context.Context, tenantID, driverID uuid.UUID) (*models.GPSPing, error) {
	p := &models.GPSPing{}
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, driver_id, truck_id, shift_id, latitude, longitude, speed_kmh, heading, accuracy, battery_level, is_moving, recorded_at, received_at, is_delayed, created_at
		 FROM gps_pings WHERE tenant_id = $1 AND driver_id = $2 ORDER BY recorded_at DESC LIMIT 1`,
		tenantID, driverID,
	).Scan(&p.ID, &p.TenantID, &p.DriverID, &p.TruckID, &p.ShiftID, &p.Latitude, &p.Longitude, &p.SpeedKmh, &p.Heading, &p.Accuracy, &p.BatteryLevel, &p.IsMoving, &p.RecordedAt, &p.ReceivedAt, &p.IsDelayed, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetByShift returns all pings for a shift ordered by time.
func (r *GPSPingRepo) GetByShift(ctx context.Context, tenantID, shiftID uuid.UUID) ([]models.GPSPing, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, driver_id, truck_id, shift_id, latitude, longitude, speed_kmh, heading, accuracy, battery_level, is_moving, recorded_at, received_at, is_delayed, created_at
		 FROM gps_pings WHERE tenant_id = $1 AND shift_id = $2 ORDER BY recorded_at ASC`,
		tenantID, shiftID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query pings by shift: %w", err)
	}
	defer rows.Close()

	var pings []models.GPSPing
	for rows.Next() {
		var p models.GPSPing
		if err := rows.Scan(&p.ID, &p.TenantID, &p.DriverID, &p.TruckID, &p.ShiftID, &p.Latitude, &p.Longitude, &p.SpeedKmh, &p.Heading, &p.Accuracy, &p.BatteryLevel, &p.IsMoving, &p.RecordedAt, &p.ReceivedAt, &p.IsDelayed, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan ping: %w", err)
		}
		pings = append(pings, p)
	}
	return pings, nil
}

// GetLatestByTenant returns latest ping per active driver in tenant (for live dashboard).
func (r *GPSPingRepo) GetLatestByTenant(ctx context.Context, tenantID uuid.UUID) ([]models.GPSPing, error) {
	rows, err := r.db.Query(ctx,
		`SELECT DISTINCT ON (driver_id) id, tenant_id, driver_id, truck_id, shift_id, latitude, longitude, speed_kmh, heading, accuracy, battery_level, is_moving, recorded_at, received_at, is_delayed, created_at
		 FROM gps_pings WHERE tenant_id = $1 AND recorded_at > $2 ORDER BY driver_id, recorded_at DESC`,
		tenantID, time.Now().Add(-30*time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query latest pings: %w", err)
	}
	defer rows.Close()

	var pings []models.GPSPing
	for rows.Next() {
		var p models.GPSPing
		if err := rows.Scan(&p.ID, &p.TenantID, &p.DriverID, &p.TruckID, &p.ShiftID, &p.Latitude, &p.Longitude, &p.SpeedKmh, &p.Heading, &p.Accuracy, &p.BatteryLevel, &p.IsMoving, &p.RecordedAt, &p.ReceivedAt, &p.IsDelayed, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan latest ping: %w", err)
		}
		pings = append(pings, p)
	}
	return pings, nil
}

// CalculateShiftKm calculates total distance traveled in a shift using Haversine formula.
func (r *GPSPingRepo) CalculateShiftKm(ctx context.Context, tenantID, shiftID uuid.UUID) (float64, error) {
	var totalKm float64
	err := r.db.QueryRow(ctx,
		`WITH ordered AS (
			SELECT latitude, longitude,
				   LAG(latitude) OVER (ORDER BY recorded_at) AS prev_lat,
				   LAG(longitude) OVER (ORDER BY recorded_at) AS prev_lon
			FROM gps_pings WHERE tenant_id = $1 AND shift_id = $2 ORDER BY recorded_at
		)
		SELECT COALESCE(SUM(
			6371 * 2 * ASIN(SQRT(
				POWER(SIN(RADIANS(latitude - prev_lat) / 2), 2) +
				COS(RADIANS(prev_lat)) * COS(RADIANS(latitude)) *
				POWER(SIN(RADIANS(longitude - prev_lon) / 2), 2)
			))
		), 0) FROM ordered WHERE prev_lat IS NOT NULL`,
		tenantID, shiftID,
	).Scan(&totalKm)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate shift km: %w", err)
	}
	return totalKm, nil
}
