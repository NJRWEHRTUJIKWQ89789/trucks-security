package types

import "github.com/graphql-go/graphql"

// DashboardStatsType aggregates high-level KPIs for the dashboard overview.
var DashboardStatsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DashboardStats",
	Fields: graphql.Fields{
		"totalShipments":    &graphql.Field{Type: graphql.Int},
		"activeShipments":   &graphql.Field{Type: graphql.Int},
		"deliveredToday":    &graphql.Field{Type: graphql.Int},
		"pendingOrders":     &graphql.Field{Type: graphql.Int},
		"totalVehicles":     &graphql.Field{Type: graphql.Int},
		"activeVehicles":    &graphql.Field{Type: graphql.Int},
		"totalDrivers":      &graphql.Field{Type: graphql.Int},
		"availableDrivers":  &graphql.Field{Type: graphql.Int},
		"totalWarehouses":   &graphql.Field{Type: graphql.Int},
		"totalInventory":    &graphql.Field{Type: graphql.Int},
		"totalRevenue":      &graphql.Field{Type: graphql.Float},
		"monthlyRevenue":    &graphql.Field{Type: graphql.Float},
		"totalClients":      &graphql.Field{Type: graphql.Int},
		"activeVendors":     &graphql.Field{Type: graphql.Int},
	},
})

// ActivityItemType represents a single entry in the activity feed.
var ActivityItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ActivityItem",
	Fields: graphql.Fields{
		"id":         &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"userId":     &graphql.Field{Type: graphql.String},
		"action":     &graphql.Field{Type: graphql.String},
		"entityType": &graphql.Field{Type: graphql.String},
		"entityId":   &graphql.Field{Type: graphql.String},
		"details":    &graphql.Field{Type: graphql.String},
		"ipAddress":  &graphql.Field{Type: graphql.String},
		"createdAt":  &graphql.Field{Type: graphql.String},
	},
})

// PerformanceType provides time-series performance metrics for the dashboard.
var PerformanceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Performance",
	Fields: graphql.Fields{
		"onTimeDeliveryRate": &graphql.Field{Type: graphql.Float},
		"averageDeliveryTime": &graphql.Field{Type: graphql.Float},
		"customerSatisfaction": &graphql.Field{Type: graphql.Float},
		"fleetUtilization":    &graphql.Field{Type: graphql.Float},
		"warehouseUtilization": &graphql.Field{Type: graphql.Float},
		"orderFulfillmentRate": &graphql.Field{Type: graphql.Float},
	},
})
