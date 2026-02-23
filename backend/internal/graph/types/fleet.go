package types

import "github.com/graphql-go/graphql"

// ---------------------------------------------------------------------------
// Vehicle
// ---------------------------------------------------------------------------

// VehicleType represents a fleet vehicle.
var VehicleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Vehicle",
	Fields: graphql.Fields{
		"id":           &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"tenantId":     &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"vehicleId":    &graphql.Field{Type: graphql.String},
		"name":         &graphql.Field{Type: graphql.String},
		"type":         &graphql.Field{Type: graphql.String},
		"status":       &graphql.Field{Type: graphql.String},
		"fuelLevel":    &graphql.Field{Type: graphql.Int},
		"mileage":      &graphql.Field{Type: graphql.Int},
		"lastService":  &graphql.Field{Type: graphql.String},
		"nextService":  &graphql.Field{Type: graphql.String},
		"licensePlate": &graphql.Field{Type: graphql.String},
		"year":         &graphql.Field{Type: graphql.Int},
		"createdAt":    &graphql.Field{Type: graphql.String},
		"updatedAt":    &graphql.Field{Type: graphql.String},
	},
})

// VehicleInputType contains fields for creating or updating a vehicle.
var VehicleInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "VehicleInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"vehicleId":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"name":         &graphql.InputObjectFieldConfig{Type: graphql.String},
		"type":         &graphql.InputObjectFieldConfig{Type: graphql.String},
		"status":       &graphql.InputObjectFieldConfig{Type: graphql.String},
		"fuelLevel":    &graphql.InputObjectFieldConfig{Type: graphql.Int},
		"mileage":      &graphql.InputObjectFieldConfig{Type: graphql.Int},
		"lastService":  &graphql.InputObjectFieldConfig{Type: graphql.String},
		"nextService":  &graphql.InputObjectFieldConfig{Type: graphql.String},
		"licensePlate": &graphql.InputObjectFieldConfig{Type: graphql.String},
		"year":         &graphql.InputObjectFieldConfig{Type: graphql.Int},
	},
})

// VehicleConnectionType is a paginated list of vehicles.
var VehicleConnectionType = ConnectionType("VehicleConnection", VehicleType)

// ---------------------------------------------------------------------------
// Driver
// ---------------------------------------------------------------------------

// DriverType represents a fleet driver.
var DriverType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Driver",
	Fields: graphql.Fields{
		"id":              &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"tenantId":        &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"employeeId":      &graphql.Field{Type: graphql.String},
		"firstName":       &graphql.Field{Type: graphql.String},
		"lastName":        &graphql.Field{Type: graphql.String},
		"email":           &graphql.Field{Type: graphql.String},
		"phone":           &graphql.Field{Type: graphql.String},
		"licenseNumber":   &graphql.Field{Type: graphql.String},
		"licenseExpiry":   &graphql.Field{Type: graphql.String},
		"status":          &graphql.Field{Type: graphql.String},
		"rating":          &graphql.Field{Type: graphql.Float},
		"totalDeliveries": &graphql.Field{Type: graphql.Int},
		"vehicleId":       &graphql.Field{Type: graphql.String},
		"createdAt":       &graphql.Field{Type: graphql.String},
		"updatedAt":       &graphql.Field{Type: graphql.String},
	},
})

// DriverInputType contains fields for creating or updating a driver.
var DriverInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "DriverInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"employeeId":      &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"firstName":       &graphql.InputObjectFieldConfig{Type: graphql.String},
		"lastName":        &graphql.InputObjectFieldConfig{Type: graphql.String},
		"email":           &graphql.InputObjectFieldConfig{Type: graphql.String},
		"phone":           &graphql.InputObjectFieldConfig{Type: graphql.String},
		"licenseNumber":   &graphql.InputObjectFieldConfig{Type: graphql.String},
		"licenseExpiry":   &graphql.InputObjectFieldConfig{Type: graphql.String},
		"status":          &graphql.InputObjectFieldConfig{Type: graphql.String},
		"rating":          &graphql.InputObjectFieldConfig{Type: graphql.Float},
		"totalDeliveries": &graphql.InputObjectFieldConfig{Type: graphql.Int},
		"vehicleId":       &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

// DriverConnectionType is a paginated list of drivers.
var DriverConnectionType = ConnectionType("DriverConnection", DriverType)

// ---------------------------------------------------------------------------
// Maintenance
// ---------------------------------------------------------------------------

// MaintenanceType represents a vehicle maintenance record.
var MaintenanceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Maintenance",
	Fields: graphql.Fields{
		"id":            &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"tenantId":      &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"vehicleId":     &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"type":          &graphql.Field{Type: graphql.String},
		"description":   &graphql.Field{Type: graphql.String},
		"status":        &graphql.Field{Type: graphql.String},
		"scheduledDate": &graphql.Field{Type: graphql.String},
		"completedDate": &graphql.Field{Type: graphql.String},
		"cost":          &graphql.Field{Type: graphql.Float},
		"mechanic":      &graphql.Field{Type: graphql.String},
		"createdAt":     &graphql.Field{Type: graphql.String},
		"updatedAt":     &graphql.Field{Type: graphql.String},
	},
})

// MaintenanceInputType contains fields for creating or updating a maintenance record.
var MaintenanceInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "MaintenanceInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"vehicleId":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"type":          &graphql.InputObjectFieldConfig{Type: graphql.String},
		"description":   &graphql.InputObjectFieldConfig{Type: graphql.String},
		"status":        &graphql.InputObjectFieldConfig{Type: graphql.String},
		"scheduledDate": &graphql.InputObjectFieldConfig{Type: graphql.String},
		"completedDate": &graphql.InputObjectFieldConfig{Type: graphql.String},
		"cost":          &graphql.InputObjectFieldConfig{Type: graphql.Float},
		"mechanic":      &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

// MaintenanceConnectionType is a paginated list of maintenance records.
var MaintenanceConnectionType = ConnectionType("MaintenanceConnection", MaintenanceType)
