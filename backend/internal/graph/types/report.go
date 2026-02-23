package types

import "github.com/graphql-go/graphql"

// ---------------------------------------------------------------------------
// Revenue Report
// ---------------------------------------------------------------------------

// MonthlyRevenueDataType holds revenue figures for a single month.
var MonthlyRevenueDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "MonthlyRevenueData",
	Fields: graphql.Fields{
		"month":    &graphql.Field{Type: graphql.String},
		"revenue":  &graphql.Field{Type: graphql.Float},
		"expenses": &graphql.Field{Type: graphql.Float},
		"profit":   &graphql.Field{Type: graphql.Float},
		"orders":   &graphql.Field{Type: graphql.Int},
	},
})

// RevenueReportType aggregates revenue data across multiple months.
var RevenueReportType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RevenueReport",
	Fields: graphql.Fields{
		"totalRevenue":  &graphql.Field{Type: graphql.Float},
		"totalExpenses": &graphql.Field{Type: graphql.Float},
		"totalProfit":   &graphql.Field{Type: graphql.Float},
		"monthlyData":   &graphql.Field{Type: graphql.NewList(MonthlyRevenueDataType)},
	},
})

// ---------------------------------------------------------------------------
// Delivery Report
// ---------------------------------------------------------------------------

// MonthlyDeliveryDataType holds delivery metrics for a single month.
var MonthlyDeliveryDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "MonthlyDeliveryData",
	Fields: graphql.Fields{
		"month":          &graphql.Field{Type: graphql.String},
		"totalDeliveries": &graphql.Field{Type: graphql.Int},
		"onTime":         &graphql.Field{Type: graphql.Int},
		"late":           &graphql.Field{Type: graphql.Int},
		"onTimeRate":     &graphql.Field{Type: graphql.Float},
	},
})

// DeliveryReportType aggregates delivery performance across multiple months.
var DeliveryReportType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DeliveryReport",
	Fields: graphql.Fields{
		"totalDeliveries":     &graphql.Field{Type: graphql.Int},
		"averageOnTimeRate":   &graphql.Field{Type: graphql.Float},
		"averageDeliveryTime": &graphql.Field{Type: graphql.Float},
		"monthlyData":         &graphql.Field{Type: graphql.NewList(MonthlyDeliveryDataType)},
	},
})

// ---------------------------------------------------------------------------
// Fleet Report
// ---------------------------------------------------------------------------

// MonthlyFleetDataType holds fleet utilisation metrics for a single month.
var MonthlyFleetDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "MonthlyFleetData",
	Fields: graphql.Fields{
		"month":            &graphql.Field{Type: graphql.String},
		"activeVehicles":   &graphql.Field{Type: graphql.Int},
		"totalMileage":     &graphql.Field{Type: graphql.Int},
		"maintenanceCost":  &graphql.Field{Type: graphql.Float},
		"utilizationRate":  &graphql.Field{Type: graphql.Float},
	},
})

// TopVehicleType highlights a single vehicle's key statistics.
var TopVehicleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TopVehicle",
	Fields: graphql.Fields{
		"vehicleId":  &graphql.Field{Type: graphql.String},
		"name":       &graphql.Field{Type: graphql.String},
		"mileage":    &graphql.Field{Type: graphql.Int},
		"deliveries": &graphql.Field{Type: graphql.Int},
		"efficiency": &graphql.Field{Type: graphql.Float},
	},
})

// FleetReportType aggregates fleet performance data with monthly breakdown and top vehicles.
var FleetReportType = graphql.NewObject(graphql.ObjectConfig{
	Name: "FleetReport",
	Fields: graphql.Fields{
		"totalVehicles":      &graphql.Field{Type: graphql.Int},
		"averageUtilization": &graphql.Field{Type: graphql.Float},
		"totalMaintenanceCost": &graphql.Field{Type: graphql.Float},
		"monthlyData":        &graphql.Field{Type: graphql.NewList(MonthlyFleetDataType)},
		"topVehicles":        &graphql.Field{Type: graphql.NewList(TopVehicleType)},
	},
})
