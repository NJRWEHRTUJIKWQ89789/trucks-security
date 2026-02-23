package models

import (
	"time"

	"github.com/google/uuid"
)

type GPSPing struct {
	ID           int64     `json:"id"`
	TenantID     uuid.UUID `json:"tenant_id"`
	DriverID     uuid.UUID `json:"driver_id"`
	TruckID      uuid.UUID `json:"truck_id"`
	ShiftID      uuid.UUID `json:"shift_id"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	SpeedKmh     float64   `json:"speed_kmh"`
	Heading      int       `json:"heading"`
	Accuracy     float64   `json:"accuracy"`
	BatteryLevel int       `json:"battery_level"`
	IsMoving     bool      `json:"is_moving"`
	RecordedAt   time.Time `json:"recorded_at"`
	ReceivedAt   time.Time `json:"received_at"`
	IsDelayed    bool      `json:"is_delayed"`
	CreatedAt    time.Time `json:"created_at"`
}
