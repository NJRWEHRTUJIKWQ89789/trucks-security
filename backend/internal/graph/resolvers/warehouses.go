package resolvers

import (
	"fmt"

	"cargomax-api/internal/graph/types"
	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

// WarehouseQueries returns GraphQL query fields for warehouse and inventory operations.
func (r *Resolver) WarehouseQueries() graphql.Fields {
	return graphql.Fields{
		"warehouses": &graphql.Field{
			Type: types.WarehouseConnectionType,
			Args: graphql.FieldConfigArgument{
				"page":    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
				"perPage": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 20},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				page := p.Args["page"].(int)
				perPage := p.Args["perPage"].(int)

				items, total, err := r.WarehouseRepo.List(p.Context, tenantID, page, perPage)
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{
					"items":      items,
					"totalCount": total,
					"page":       page,
					"perPage":    perPage,
					"totalPages": (total + perPage - 1) / perPage,
				}, nil
			},
		},
		"warehouse": &graphql.Field{
			Type: types.WarehouseType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid warehouse id: %w", err)
				}
				return r.WarehouseRepo.GetByID(p.Context, tenantID, id)
			},
		},
		"inventoryItems": &graphql.Field{
			Type: types.InventoryItemConnectionType,
			Args: graphql.FieldConfigArgument{
				"warehouseId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"page":        &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
				"perPage":     &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 20},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				warehouseID, err := uuid.Parse(p.Args["warehouseId"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid warehouse id: %w", err)
				}
				page := p.Args["page"].(int)
				perPage := p.Args["perPage"].(int)

				items, total, err := r.InventoryRepo.List(p.Context, tenantID, &warehouseID, page, perPage)
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{
					"items":      items,
					"totalCount": total,
					"page":       page,
					"perPage":    perPage,
					"totalPages": (total + perPage - 1) / perPage,
				}, nil
			},
		},
		"lowStockItems": &graphql.Field{
			Type: graphql.NewList(types.InventoryItemType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				return r.InventoryRepo.GetLowStock(p.Context, tenantID)
			},
		},
	}
}

// WarehouseMutations returns GraphQL mutation fields for warehouse and inventory operations.
func (r *Resolver) WarehouseMutations() graphql.Fields {
	return graphql.Fields{
		"createWarehouse": &graphql.Field{
			Type: types.WarehouseType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.WarehouseInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				input := p.Args["input"].(map[string]interface{})

				w := &models.Warehouse{
					TenantID: tenantID,
					Name:     input["name"].(string),
				}
				if v, ok := input["location"].(string); ok {
					w.Location = &v
				}
				if v, ok := input["address"].(string); ok {
					w.Address = &v
				}
				if v, ok := input["capacity"].(int); ok {
					w.Capacity = v
				}
				if v, ok := input["usedCapacity"].(int); ok {
					w.UsedCapacity = v
				}
				if v, ok := input["manager"].(string); ok {
					w.Manager = &v
				}
				if v, ok := input["phone"].(string); ok {
					w.Phone = &v
				}
				if v, ok := input["status"].(string); ok {
					w.Status = v
				}

				if err := r.WarehouseRepo.Create(p.Context, w); err != nil {
					return nil, err
				}
				return w, nil
			},
		},
		"updateWarehouse": &graphql.Field{
			Type: types.WarehouseType,
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.WarehouseInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid warehouse id: %w", err)
				}
				input := p.Args["input"].(map[string]interface{})

				w := &models.Warehouse{
					Name: input["name"].(string),
				}
				if v, ok := input["location"].(string); ok {
					w.Location = &v
				}
				if v, ok := input["address"].(string); ok {
					w.Address = &v
				}
				if v, ok := input["capacity"].(int); ok {
					w.Capacity = v
				}
				if v, ok := input["usedCapacity"].(int); ok {
					w.UsedCapacity = v
				}
				if v, ok := input["manager"].(string); ok {
					w.Manager = &v
				}
				if v, ok := input["phone"].(string); ok {
					w.Phone = &v
				}
				if v, ok := input["status"].(string); ok {
					w.Status = v
				}

				if err := r.WarehouseRepo.Update(p.Context, tenantID, id, w); err != nil {
					return nil, err
				}
				return r.WarehouseRepo.GetByID(p.Context, tenantID, id)
			},
		},
		"deleteWarehouse": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid warehouse id: %w", err)
				}
				if err := r.WarehouseRepo.Delete(p.Context, tenantID, id); err != nil {
					return false, err
				}
				return true, nil
			},
		},
		"createInventoryItem": &graphql.Field{
			Type: types.InventoryItemType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.InventoryItemInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				input := p.Args["input"].(map[string]interface{})

				warehouseID, err := uuid.Parse(input["warehouseId"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid warehouse id: %w", err)
				}

				// Validate the warehouse belongs to the same tenant (prevent IDOR).
				if _, err := r.WarehouseRepo.GetByID(p.Context, tenantID, warehouseID); err != nil {
					return nil, fmt.Errorf("warehouse not found in tenant")
				}

				item := &models.InventoryItem{
					TenantID:    tenantID,
					WarehouseID: warehouseID,
					SKU:         input["sku"].(string),
				}
				if v, ok := input["name"].(string); ok {
					item.Name = &v
				}
				if v, ok := input["category"].(string); ok {
					item.Category = &v
				}
				if v, ok := input["quantity"].(int); ok {
					item.Quantity = v
				}
				if v, ok := input["minQuantity"].(int); ok {
					item.MinQuantity = v
				}
				if v, ok := input["unitPrice"].(float64); ok {
					item.UnitPrice = &v
				}
				if v, ok := input["weight"].(float64); ok {
					item.Weight = &v
				}
				if v, ok := input["status"].(string); ok {
					item.Status = v
				}

				if err := r.InventoryRepo.Create(p.Context, item); err != nil {
					return nil, err
				}
				return item, nil
			},
		},
		"updateInventoryItem": &graphql.Field{
			Type: types.InventoryItemType,
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.InventoryItemInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid inventory item id: %w", err)
				}
				input := p.Args["input"].(map[string]interface{})

				warehouseID, err := uuid.Parse(input["warehouseId"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid warehouse id: %w", err)
				}

				// Validate the warehouse belongs to the same tenant (prevent IDOR).
				if _, err := r.WarehouseRepo.GetByID(p.Context, tenantID, warehouseID); err != nil {
					return nil, fmt.Errorf("warehouse not found in tenant")
				}

				item := &models.InventoryItem{
					WarehouseID: warehouseID,
					SKU:         input["sku"].(string),
				}
				if v, ok := input["name"].(string); ok {
					item.Name = &v
				}
				if v, ok := input["category"].(string); ok {
					item.Category = &v
				}
				if v, ok := input["quantity"].(int); ok {
					item.Quantity = v
				}
				if v, ok := input["minQuantity"].(int); ok {
					item.MinQuantity = v
				}
				if v, ok := input["unitPrice"].(float64); ok {
					item.UnitPrice = &v
				}
				if v, ok := input["weight"].(float64); ok {
					item.Weight = &v
				}
				if v, ok := input["status"].(string); ok {
					item.Status = v
				}

				if err := r.InventoryRepo.Update(p.Context, tenantID, id, item); err != nil {
					return nil, err
				}
				return r.InventoryRepo.GetByID(p.Context, tenantID, id)
			},
		},
		"restockItem": &graphql.Field{
			Type: types.InventoryItemType,
			Args: graphql.FieldConfigArgument{
				"id":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"quantity": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid inventory item id: %w", err)
				}
				quantity := p.Args["quantity"].(int)
				if quantity <= 0 {
					return nil, fmt.Errorf("restock quantity must be positive")
				}

				if err := r.InventoryRepo.Restock(p.Context, tenantID, id, quantity); err != nil {
					return nil, err
				}
				return r.InventoryRepo.GetByID(p.Context, tenantID, id)
			},
		},
	}
}
