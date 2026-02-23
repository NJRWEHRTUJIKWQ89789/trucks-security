package models

// DashboardStats holds aggregate statistics for the dashboard view.
type DashboardStats struct {
	TotalShipments    int     `json:"total_shipments"`
	ActiveShipments   int     `json:"active_shipments"`
	DeliveredToday    int     `json:"delivered_today"`
	DelayedShipments  int     `json:"delayed_shipments"`
	TotalVehicles     int     `json:"total_vehicles"`
	ActiveVehicles    int     `json:"active_vehicles"`
	TotalDrivers      int     `json:"total_drivers"`
	AvailableDrivers  int     `json:"available_drivers"`
	TotalOrders       int     `json:"total_orders"`
	PendingOrders     int     `json:"pending_orders"`
	TotalWarehouses   int     `json:"total_warehouses"`
	TotalInventory    int     `json:"total_inventory"`
	LowStockItems     int     `json:"low_stock_items"`
	TotalClients      int     `json:"total_clients"`
	TotalRevenue      float64 `json:"total_revenue"`
	MaintenanceDue    int     `json:"maintenance_due"`
}

// Performance holds performance metrics for the dashboard view.
type Performance struct {
	OnTimeDeliveryRate   float64 `json:"on_time_delivery_rate"`
	AverageDeliveryTime  float64 `json:"average_delivery_time"`
	CustomerSatisfaction float64 `json:"customer_satisfaction"`
	FleetUtilization     float64 `json:"fleet_utilization"`
	WarehouseUtilization float64 `json:"warehouse_utilization"`
	OrderFulfillmentRate float64 `json:"order_fulfillment_rate"`
}
