package types

import "github.com/graphql-go/graphql"

// OrderType represents a customer order.
var OrderType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Order",
	Fields: graphql.Fields{
		"id":                 &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"tenantId":           &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"orderNumber":        &graphql.Field{Type: graphql.String},
		"customerName":       &graphql.Field{Type: graphql.String},
		"customerEmail":      &graphql.Field{Type: graphql.String},
		"status":             &graphql.Field{Type: graphql.String},
		"type":               &graphql.Field{Type: graphql.String},
		"totalAmount":        &graphql.Field{Type: graphql.Float},
		"shipmentId":         &graphql.Field{Type: graphql.String},
		"scheduledDate":      &graphql.Field{Type: graphql.String},
		"returnReason":       &graphql.Field{Type: graphql.String},
		"cancellationReason": &graphql.Field{Type: graphql.String},
		"createdAt":          &graphql.Field{Type: graphql.String},
		"updatedAt":          &graphql.Field{Type: graphql.String},
	},
})

// OrderInputType contains fields for creating or updating an order.
var OrderInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "OrderInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"orderNumber":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"customerName":       &graphql.InputObjectFieldConfig{Type: graphql.String},
		"customerEmail":      &graphql.InputObjectFieldConfig{Type: graphql.String},
		"status":             &graphql.InputObjectFieldConfig{Type: graphql.String},
		"type":               &graphql.InputObjectFieldConfig{Type: graphql.String},
		"totalAmount":        &graphql.InputObjectFieldConfig{Type: graphql.Float},
		"shipmentId":         &graphql.InputObjectFieldConfig{Type: graphql.String},
		"scheduledDate":      &graphql.InputObjectFieldConfig{Type: graphql.String},
		"returnReason":       &graphql.InputObjectFieldConfig{Type: graphql.String},
		"cancellationReason": &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

// OrderConnectionType is a paginated list of orders.
var OrderConnectionType = ConnectionType("OrderConnection", OrderType)
