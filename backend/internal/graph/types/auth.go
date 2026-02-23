package types

import "github.com/graphql-go/graphql"

// UserType represents a user account.
var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":               &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"tenantId":         &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"email":            &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"firstName":        &graphql.Field{Type: graphql.String},
		"lastName":         &graphql.Field{Type: graphql.String},
		"role":             &graphql.Field{Type: graphql.String},
		"emailVerified":    &graphql.Field{Type: graphql.Boolean},
		"avatarUrl":        &graphql.Field{Type: graphql.String},
		"createdAt":        &graphql.Field{Type: graphql.String},
		"updatedAt":        &graphql.Field{Type: graphql.String},
	},
})

// AuthPayloadType is returned after successful login or registration.
var AuthPayloadType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AuthPayload",
	Fields: graphql.Fields{
		"user":  &graphql.Field{Type: UserType},
		"token": &graphql.Field{Type: graphql.String},
	},
})

// RegisterInputType contains fields required to register a new user.
var RegisterInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "RegisterInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"email":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"password":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"firstName": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"lastName":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"tenantName": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
	},
})

// LoginInputType contains fields required to log in.
var LoginInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "LoginInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"email":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"password": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
	},
})
