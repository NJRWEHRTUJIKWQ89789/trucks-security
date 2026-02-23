package models

import (
	"time"

	"github.com/google/uuid"
)

// Warehouse represents a storage warehouse facility.
type Warehouse struct {
	ID           uuid.UUID `json:"id"`
	TenantID     uuid.UUID `json:"tenant_id"`
	Name         string    `json:"name"`
	Location     *string   `json:"location"`
	Address      *string   `json:"address"`
	Capacity     int       `json:"capacity"`
	UsedCapacity int       `json:"used_capacity"`
	Manager      *string   `json:"manager"`
	Phone        *string   `json:"phone"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// WarehouseCapacityStats holds aggregated capacity information.
type WarehouseCapacityStats struct {
	TotalCapacity int     `json:"total_capacity"`
	UsedCapacity  int     `json:"used_capacity"`
	FreeCapacity  int     `json:"free_capacity"`
	Utilization   float64 `json:"utilization"`
}
