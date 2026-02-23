package types

import "github.com/graphql-go/graphql"

// ---------------------------------------------------------------------------
// Setting
// ---------------------------------------------------------------------------

// SettingType represents a tenant-scoped configuration key-value pair.
var SettingType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Setting",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"tenantId":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"key":       &graphql.Field{Type: graphql.String},
		"value":     &graphql.Field{Type: graphql.String},
		"category":  &graphql.Field{Type: graphql.String},
		"updatedBy": &graphql.Field{Type: graphql.String},
		"updatedAt": &graphql.Field{Type: graphql.String},
	},
})

// ---------------------------------------------------------------------------
// Role
// ---------------------------------------------------------------------------

// RoleType represents a tenant-scoped role with JSON permissions.
var RoleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Role",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"tenantId":    &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"name":        &graphql.Field{Type: graphql.String},
		"permissions": &graphql.Field{Type: graphql.String},
		"createdAt":   &graphql.Field{Type: graphql.String},
		"updatedAt":   &graphql.Field{Type: graphql.String},
	},
})

// RoleInputType contains fields for creating or updating a role.
var RoleInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "RoleInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"permissions": &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

// ---------------------------------------------------------------------------
// Notification
// ---------------------------------------------------------------------------

// NotificationType represents a user notification.
var NotificationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Notification",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"tenantId":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"userId":    &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"title":     &graphql.Field{Type: graphql.String},
		"message":   &graphql.Field{Type: graphql.String},
		"type":      &graphql.Field{Type: graphql.String},
		"read":      &graphql.Field{Type: graphql.Boolean},
		"createdAt": &graphql.Field{Type: graphql.String},
	},
})

// ---------------------------------------------------------------------------
// NotificationPreference
// ---------------------------------------------------------------------------

// NotificationPreferenceType represents a user's notification delivery preferences.
var NotificationPreferenceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NotificationPreference",
	Fields: graphql.Fields{
		"id":           &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"tenantId":     &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"userId":       &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"eventType":    &graphql.Field{Type: graphql.String},
		"emailEnabled": &graphql.Field{Type: graphql.Boolean},
		"smsEnabled":   &graphql.Field{Type: graphql.Boolean},
		"pushEnabled":  &graphql.Field{Type: graphql.Boolean},
	},
})

// NotificationConnectionType is a paginated wrapper for notifications.
var NotificationConnectionType = ConnectionType("NotificationConnection", NotificationType)

// NotificationPrefInputType contains fields for creating or updating notification preferences.
var NotificationPrefInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "NotificationPrefInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"eventType":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"emailEnabled": &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
		"smsEnabled":   &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
		"pushEnabled":  &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
	},
})
