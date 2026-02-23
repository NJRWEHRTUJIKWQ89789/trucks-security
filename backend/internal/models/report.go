package models

// MonthlyBreakdown holds data for a single month in a report.
type MonthlyBreakdown struct {
	Month string  `json:"month"`
	Value float64 `json:"value"`
	Count int     `json:"count"`
}

// RevenueReport holds revenue report data with monthly breakdowns.
type RevenueReport struct {
	TotalRevenue     float64            `json:"total_revenue"`
	MonthlyBreakdown []MonthlyBreakdown `json:"monthly_breakdown"`
}

// DeliveryReport holds delivery performance report data.
type DeliveryReport struct {
	TotalDeliveries  int                `json:"total_deliveries"`
	OnTimeDeliveries int                `json:"on_time_deliveries"`
	DelayedDeliveries int               `json:"delayed_deliveries"`
	OnTimeRate       float64            `json:"on_time_rate"`
	MonthlyBreakdown []MonthlyBreakdown `json:"monthly_breakdown"`
}

// FleetReport holds fleet utilization report data.
type FleetReport struct {
	TotalVehicles    int                `json:"total_vehicles"`
	ActiveVehicles   int                `json:"active_vehicles"`
	MaintenanceCount int                `json:"maintenance_count"`
	TotalCost        float64            `json:"total_cost"`
	MonthlyBreakdown []MonthlyBreakdown `json:"monthly_breakdown"`
}
