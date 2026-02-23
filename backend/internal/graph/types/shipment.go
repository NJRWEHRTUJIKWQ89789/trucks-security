package types

import "github.com/graphql-go/graphql"

// ShipmentType represents a shipment entity.
var ShipmentType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Shipment",
	Fields: graphql.Fields{
		"id":                &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"tenantId":          &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"trackingNumber":    &graphql.Field{Type: graphql.String},
		"origin":            &graphql.Field{Type: graphql.String},
		"destination":       &graphql.Field{Type: graphql.String},
		"status":            &graphql.Field{Type: graphql.String},
		"carrier":           &graphql.Field{Type: graphql.String},
		"weight":            &graphql.Field{Type: graphql.Float},
		"dimensions":        &graphql.Field{Type: graphql.String},
		"estimatedDelivery": &graphql.Field{Type: graphql.String},
		"actualDelivery":    &graphql.Field{Type: graphql.String},
		"customerName":      &graphql.Field{Type: graphql.String},
		"customerEmail":     &graphql.Field{Type: graphql.String},
		"notes":             &graphql.Field{Type: graphql.String},
		"createdAt":         &graphql.Field{Type: graphql.String},
		"updatedAt":         &graphql.Field{Type: graphql.String},
	},
})

// ShipmentInputType contains fields for creating or updating a shipment.
var ShipmentInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "ShipmentInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"trackingNumber":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"origin":            &graphql.InputObjectFieldConfig{Type: graphql.String},
		"destination":       &graphql.InputObjectFieldConfig{Type: graphql.String},
		"status":            &graphql.InputObjectFieldConfig{Type: graphql.String},
		"carrier":           &graphql.InputObjectFieldConfig{Type: graphql.String},
		"weight":            &graphql.InputObjectFieldConfig{Type: graphql.Float},
		"dimensions":        &graphql.InputObjectFieldConfig{Type: graphql.String},
		"estimatedDelivery": &graphql.InputObjectFieldConfig{Type: graphql.String},
		"actualDelivery":    &graphql.InputObjectFieldConfig{Type: graphql.String},
		"customerName":      &graphql.InputObjectFieldConfig{Type: graphql.String},
		"customerEmail":     &graphql.InputObjectFieldConfig{Type: graphql.String},
		"notes":             &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

// ShipmentConnectionType is a paginated list of shipments.
var ShipmentConnectionType = ConnectionType("ShipmentConnection", ShipmentType)
