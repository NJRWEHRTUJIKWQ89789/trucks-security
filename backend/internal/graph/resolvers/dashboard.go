package resolvers

import (
	"fmt"

	"cargomax-api/internal/graph/types"

	"github.com/graphql-go/graphql"
)

// ActivityConnectionType is a paginated wrapper for activity items.
var activityConnectionType = types.ConnectionType("ActivityConnection", types.ActivityItemType)

// DashboardQueries returns the GraphQL query fields for the dashboard domain.
func (r *Resolver) DashboardQueries() graphql.Fields {
	return graphql.Fields{
		// -----------------------------------------------------------------
		// dashboardStats
		// -----------------------------------------------------------------
		"dashboardStats": &graphql.Field{
			Type:        types.DashboardStatsType,
			Description: "Returns aggregated KPI statistics for the tenant dashboard.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}

				stats, err := r.DashboardRepo.GetStats(p.Context, tenantID)
				if err != nil {
					return nil, fmt.Errorf("failed to fetch dashboard stats: %w", err)
				}
				return stats, nil
			},
		},

		// -----------------------------------------------------------------
		// dashboardActivity (paginated)
		// -----------------------------------------------------------------
		"dashboardActivity": &graphql.Field{
			Type:        activityConnectionType,
			Description: "Returns a paginated list of recent activity log entries.",
			Args: graphql.FieldConfigArgument{
				"page":    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
				"perPage": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 20},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				page := p.Args["page"].(int)
				perPage := p.Args["perPage"].(int)

				items, total, err := r.ActivityRepo.List(p.Context, tenantID, page, perPage)
				if err != nil {
					return nil, fmt.Errorf("failed to fetch activity: %w", err)
				}

				totalPages := 0
				if perPage > 0 {
					totalPages = (total + perPage - 1) / perPage
				}

				return map[string]interface{}{
					"items":      items,
					"totalCount": total,
					"page":       page,
					"perPage":    perPage,
					"totalPages": totalPages,
				}, nil
			},
		},

		// -----------------------------------------------------------------
		// dashboardPerformance
		// -----------------------------------------------------------------
		"dashboardPerformance": &graphql.Field{
			Type:        types.PerformanceType,
			Description: "Returns performance metrics for the tenant dashboard.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}

				// Performance is derived from dashboard stats.
				stats, err := r.DashboardRepo.GetStats(p.Context, tenantID)
				if err != nil {
					return nil, fmt.Errorf("failed to fetch performance metrics: %w", err)
				}

				// Compute on-time delivery rate: delivered today / (delivered today + delayed).
				var onTimeRate float64
				if delivered := stats.DeliveredToday + stats.DelayedShipments; delivered > 0 {
					onTimeRate = float64(stats.DeliveredToday) / float64(delivered) * 100.0
				}

				// Fleet utilization: active vehicles / total vehicles.
				var fleetUtil float64
				if stats.TotalVehicles > 0 {
					fleetUtil = float64(stats.ActiveVehicles) / float64(stats.TotalVehicles) * 100.0
				}

				// Warehouse utilization: (total inventory - low stock) / total inventory.
				var warehouseUtil float64
				if stats.TotalInventory > 0 {
					warehouseUtil = float64(stats.TotalInventory-stats.LowStockItems) / float64(stats.TotalInventory) * 100.0
				}

				// Order fulfillment rate: (total orders - pending) / total orders.
				var fulfillmentRate float64
				if stats.TotalOrders > 0 {
					fulfillmentRate = float64(stats.TotalOrders-stats.PendingOrders) / float64(stats.TotalOrders) * 100.0
				}

				return map[string]interface{}{
					"onTimeDeliveryRate":   onTimeRate,
					"averageDeliveryTime":  0.0,
					"customerSatisfaction": 0.0,
					"fleetUtilization":     fleetUtil,
					"warehouseUtilization": warehouseUtil,
					"orderFulfillmentRate": fulfillmentRate,
				}, nil
			},
		},
	}
}
