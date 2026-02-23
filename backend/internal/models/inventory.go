package models

import (
	"time"

	"github.com/google/uuid"
)

// InventoryItem represents an item stored in a warehouse.
type InventoryItem struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	WarehouseID uuid.UUID `json:"warehouse_id"`
	SKU         string    `json:"sku"`
	Name        *string   `json:"name"`
	Category    *string   `json:"category"`
	Quantity    int       `json:"quantity"`
	MinQuantity int       `json:"min_quantity"`
	UnitPrice   *float64  `json:"unit_price"`
	Weight      *float64  `json:"weight"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
