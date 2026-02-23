package types

import "github.com/graphql-go/graphql"

// ---------------------------------------------------------------------------
// Client
// ---------------------------------------------------------------------------

// ClientType represents a business client.
var ClientType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Client",
	Fields: graphql.Fields{
		"id":                 &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"tenantId":           &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"companyName":        &graphql.Field{Type: graphql.String},
		"contactPerson":      &graphql.Field{Type: graphql.String},
		"email":              &graphql.Field{Type: graphql.String},
		"phone":              &graphql.Field{Type: graphql.String},
		"address":            &graphql.Field{Type: graphql.String},
		"industry":           &graphql.Field{Type: graphql.String},
		"totalShipments":     &graphql.Field{Type: graphql.Int},
		"totalSpent":         &graphql.Field{Type: graphql.Float},
		"satisfactionRating": &graphql.Field{Type: graphql.Float},
		"status":             &graphql.Field{Type: graphql.String},
		"createdAt":          &graphql.Field{Type: graphql.String},
		"updatedAt":          &graphql.Field{Type: graphql.String},
	},
})

// ClientInputType contains fields for creating or updating a client.
var ClientInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "ClientInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"companyName":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"contactPerson":      &graphql.InputObjectFieldConfig{Type: graphql.String},
		"email":              &graphql.InputObjectFieldConfig{Type: graphql.String},
		"phone":              &graphql.InputObjectFieldConfig{Type: graphql.String},
		"address":            &graphql.InputObjectFieldConfig{Type: graphql.String},
		"industry":           &graphql.InputObjectFieldConfig{Type: graphql.String},
		"totalShipments":     &graphql.InputObjectFieldConfig{Type: graphql.Int},
		"totalSpent":         &graphql.InputObjectFieldConfig{Type: graphql.Float},
		"satisfactionRating": &graphql.InputObjectFieldConfig{Type: graphql.Float},
		"status":             &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

// ClientConnectionType is a paginated list of clients.
var ClientConnectionType = ConnectionType("ClientConnection", ClientType)

// ---------------------------------------------------------------------------
// Feedback
// ---------------------------------------------------------------------------

// FeedbackType represents client feedback on a service.
var FeedbackType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Feedback",
	Fields: graphql.Fields{
		"id":         &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"tenantId":   &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"clientId":   &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"clientName": &graphql.Field{Type: graphql.String},
		"rating":     &graphql.Field{Type: graphql.Int},
		"comment":    &graphql.Field{Type: graphql.String},
		"category":   &graphql.Field{Type: graphql.String},
		"createdAt":  &graphql.Field{Type: graphql.String},
	},
})

// FeedbackInputType contains fields for creating feedback.
var FeedbackInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "FeedbackInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"clientId": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"rating":   &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
		"comment":  &graphql.InputObjectFieldConfig{Type: graphql.String},
		"category": &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

// FeedbackConnectionType is a paginated list of feedback entries.
var FeedbackConnectionType = ConnectionType("FeedbackConnection", FeedbackType)
