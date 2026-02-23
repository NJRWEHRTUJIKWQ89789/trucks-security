package types

import "github.com/graphql-go/graphql"

// VendorType represents a vendor or supplier.
var VendorType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Vendor",
	Fields: graphql.Fields{
		"id":            &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"tenantId":      &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"name":          &graphql.Field{Type: graphql.String},
		"contactPerson": &graphql.Field{Type: graphql.String},
		"email":         &graphql.Field{Type: graphql.String},
		"phone":         &graphql.Field{Type: graphql.String},
		"address":       &graphql.Field{Type: graphql.String},
		"category":      &graphql.Field{Type: graphql.String},
		"rating":        &graphql.Field{Type: graphql.Float},
		"contractStart": &graphql.Field{Type: graphql.String},
		"contractEnd":   &graphql.Field{Type: graphql.String},
		"status":        &graphql.Field{Type: graphql.String},
		"createdAt":     &graphql.Field{Type: graphql.String},
		"updatedAt":     &graphql.Field{Type: graphql.String},
	},
})

// VendorInputType contains fields for creating or updating a vendor.
var VendorInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "VendorInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name":          &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"contactPerson": &graphql.InputObjectFieldConfig{Type: graphql.String},
		"email":         &graphql.InputObjectFieldConfig{Type: graphql.String},
		"phone":         &graphql.InputObjectFieldConfig{Type: graphql.String},
		"address":       &graphql.InputObjectFieldConfig{Type: graphql.String},
		"category":      &graphql.InputObjectFieldConfig{Type: graphql.String},
		"rating":        &graphql.InputObjectFieldConfig{Type: graphql.Float},
		"contractStart": &graphql.InputObjectFieldConfig{Type: graphql.String},
		"contractEnd":   &graphql.InputObjectFieldConfig{Type: graphql.String},
		"status":        &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

// VendorConnectionType is a paginated list of vendors.
var VendorConnectionType = ConnectionType("VendorConnection", VendorType)
