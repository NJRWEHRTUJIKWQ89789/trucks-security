package types

import "github.com/graphql-go/graphql"

// ---------------------------------------------------------------------------
// Warehouse
// ---------------------------------------------------------------------------

// WarehouseType represents a warehouse facility.
var WarehouseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Warehouse",
	Fields: graphql.Fields{
		"id":           &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"tenantId":     &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"name":         &graphql.Field{Type: graphql.String},
		"location":     &graphql.Field{Type: graphql.String},
		"address":      &graphql.Field{Type: graphql.String},
		"capacity":     &graphql.Field{Type: graphql.Int},
		"usedCapacity": &graphql.Field{Type: graphql.Int},
		"manager":      &graphql.Field{Type: graphql.String},
		"phone":        &graphql.Field{Type: graphql.String},
		"status":       &graphql.Field{Type: graphql.String},
		"createdAt":    &graphql.Field{Type: graphql.String},
		"updatedAt":    &graphql.Field{Type: graphql.String},
	},
})

// WarehouseInputType contains fields for creating or updating a warehouse.
var WarehouseInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "WarehouseInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name":         &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"location":     &graphql.InputObjectFieldConfig{Type: graphql.String},
		"address":      &graphql.InputObjectFieldConfig{Type: graphql.String},
		"capacity":     &graphql.InputObjectFieldConfig{Type: graphql.Int},
		"usedCapacity": &graphql.InputObjectFieldConfig{Type: graphql.Int},
		"manager":      &graphql.InputObjectFieldConfig{Type: graphql.String},
		"phone":        &graphql.InputObjectFieldConfig{Type: graphql.String},
		"status":       &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

// WarehouseConnectionType is a paginated list of warehouses.
var WarehouseConnectionType = ConnectionType("WarehouseConnection", WarehouseType)

// ---------------------------------------------------------------------------
// InventoryItem
// ---------------------------------------------------------------------------

// InventoryItemType represents an inventory item stored in a warehouse.
var InventoryItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "InventoryItem",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"tenantId":    &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"warehouseId": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"sku":         &graphql.Field{Type: graphql.String},
		"name":        &graphql.Field{Type: graphql.String},
		"category":    &graphql.Field{Type: graphql.String},
		"quantity":    &graphql.Field{Type: graphql.Int},
		"minQuantity": &graphql.Field{Type: graphql.Int},
		"unitPrice":   &graphql.Field{Type: graphql.Float},
		"weight":      &graphql.Field{Type: graphql.Float},
		"status":      &graphql.Field{Type: graphql.String},
		"createdAt":   &graphql.Field{Type: graphql.String},
		"updatedAt":   &graphql.Field{Type: graphql.String},
	},
})

// InventoryItemInputType contains fields for creating or updating an inventory item.
var InventoryItemInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "InventoryItemInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"warehouseId": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"sku":         &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"name":        &graphql.InputObjectFieldConfig{Type: graphql.String},
		"category":    &graphql.InputObjectFieldConfig{Type: graphql.String},
		"quantity":    &graphql.InputObjectFieldConfig{Type: graphql.Int},
		"minQuantity": &graphql.InputObjectFieldConfig{Type: graphql.Int},
		"unitPrice":   &graphql.InputObjectFieldConfig{Type: graphql.Float},
		"weight":      &graphql.InputObjectFieldConfig{Type: graphql.Float},
		"status":      &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

// InventoryItemConnectionType is a paginated list of inventory items.
var InventoryItemConnectionType = ConnectionType("InventoryItemConnection", InventoryItemType)
