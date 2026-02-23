package resolvers

import (
	"fmt"
	"time"

	"cargomax-api/internal/graph/types"
	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

// clampPagination ensures page and perPage are within safe bounds.
func clampPagination(page, perPage int) (int, int) {
	const maxPerPage = 100
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > maxPerPage {
		perPage = maxPerPage
	}
	return page, perPage
}

// OrderQueries returns GraphQL query fields for order operations.
func (r *Resolver) OrderQueries() graphql.Fields {
	return graphql.Fields{
		"orders": &graphql.Field{
			Type: types.OrderConnectionType,
			Args: graphql.FieldConfigArgument{
				"page":    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
				"perPage": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 20},
				"status":  &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				page, perPage := clampPagination(p.Args["page"].(int), p.Args["perPage"].(int))

				var status string
				if v, ok := p.Args["status"].(string); ok {
					status = v
				}

				items, total, err := r.OrderRepo.List(p.Context, tenantID, status, page, perPage)
				if err != nil {
					return nil, err
				}

				totalPages := 0
				if perPage > 0 {
					totalPages = (total + perPage - 1) / perPage
				}
				return map[string]interface{}{
					"items":      items,
					"totalCount": total,
					"page":       page,
					"perPage":    perPage,
					"totalPages": totalPages,
				}, nil
			},
		},
		"order": &graphql.Field{
			Type: types.OrderType,
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
					return nil, fmt.Errorf("invalid order id: %w", err)
				}
				return r.OrderRepo.GetByID(p.Context, tenantID, id)
			},
		},
		"scheduledOrders": &graphql.Field{
			Type: types.OrderConnectionType,
			Args: graphql.FieldConfigArgument{
				"page":    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
				"perPage": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 20},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				page, perPage := clampPagination(p.Args["page"].(int), p.Args["perPage"].(int))

				items, total, err := r.OrderRepo.List(p.Context, tenantID, "scheduled", page, perPage)
				if err != nil {
					return nil, err
				}

				totalPages := 0
				if perPage > 0 {
					totalPages = (total + perPage - 1) / perPage
				}
				return map[string]interface{}{
					"items":      items,
					"totalCount": total,
					"page":       page,
					"perPage":    perPage,
					"totalPages": totalPages,
				}, nil
			},
		},
		"returnOrders": &graphql.Field{
			Type: types.OrderConnectionType,
			Args: graphql.FieldConfigArgument{
				"page":    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
				"perPage": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 20},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				page, perPage := clampPagination(p.Args["page"].(int), p.Args["perPage"].(int))

				items, total, err := r.OrderRepo.List(p.Context, tenantID, "returned", page, perPage)
				if err != nil {
					return nil, err
				}

				totalPages := 0
				if perPage > 0 {
					totalPages = (total + perPage - 1) / perPage
				}
				return map[string]interface{}{
					"items":      items,
					"totalCount": total,
					"page":       page,
					"perPage":    perPage,
					"totalPages": totalPages,
				}, nil
			},
		},
		"cancelledOrders": &graphql.Field{
			Type: types.OrderConnectionType,
			Args: graphql.FieldConfigArgument{
				"page":    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
				"perPage": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 20},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				page, perPage := clampPagination(p.Args["page"].(int), p.Args["perPage"].(int))

				items, total, err := r.OrderRepo.List(p.Context, tenantID, "cancelled", page, perPage)
				if err != nil {
					return nil, err
				}

				totalPages := 0
				if perPage > 0 {
					totalPages = (total + perPage - 1) / perPage
				}
				return map[string]interface{}{
					"items":      items,
					"totalCount": total,
					"page":       page,
					"perPage":    perPage,
					"totalPages": totalPages,
				}, nil
			},
		},
	}
}

// OrderMutations returns GraphQL mutation fields for order operations.
func (r *Resolver) OrderMutations() graphql.Fields {
	return graphql.Fields{
		"createOrder": &graphql.Field{
			Type: types.OrderType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.OrderInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				input := p.Args["input"].(map[string]interface{})

				o := &models.Order{
					TenantID:    tenantID,
					OrderNumber: input["orderNumber"].(string),
				}
				if v, ok := input["customerName"].(string); ok {
					o.CustomerName = &v
				}
				if v, ok := input["customerEmail"].(string); ok {
					o.CustomerEmail = &v
				}
				if v, ok := input["status"].(string); ok {
					o.Status = v
				}
				if v, ok := input["type"].(string); ok {
					o.Type = v
				}
				if v, ok := input["totalAmount"].(float64); ok {
					o.TotalAmount = &v
				}
				if v, ok := input["shipmentId"].(string); ok {
					sid, err := uuid.Parse(v)
					if err != nil {
						return nil, fmt.Errorf("invalid shipment id: %w", err)
					}
					// Validate the shipment belongs to the same tenant (prevent IDOR).
					if _, err := r.ShipmentRepo.GetByID(p.Context, tenantID, sid); err != nil {
						return nil, fmt.Errorf("shipment not found in tenant")
					}
					o.ShipmentID = &sid
				}
				if v, ok := input["scheduledDate"].(string); ok {
					t, err := time.Parse(time.RFC3339, v)
					if err != nil {
						return nil, fmt.Errorf("invalid scheduledDate format: %w", err)
					}
					o.ScheduledDate = &t
				}

				if err := r.OrderRepo.Create(p.Context, o); err != nil {
					return nil, err
				}
				return o, nil
			},
		},
		"updateOrder": &graphql.Field{
			Type: types.OrderType,
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.OrderInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid order id: %w", err)
				}
				input := p.Args["input"].(map[string]interface{})

				o := &models.Order{
					OrderNumber: input["orderNumber"].(string),
				}
				if v, ok := input["customerName"].(string); ok {
					o.CustomerName = &v
				}
				if v, ok := input["customerEmail"].(string); ok {
					o.CustomerEmail = &v
				}
				if v, ok := input["status"].(string); ok {
					o.Status = v
				}
				if v, ok := input["type"].(string); ok {
					o.Type = v
				}
				if v, ok := input["totalAmount"].(float64); ok {
					o.TotalAmount = &v
				}
				if v, ok := input["shipmentId"].(string); ok {
					sid, err := uuid.Parse(v)
					if err != nil {
						return nil, fmt.Errorf("invalid shipment id: %w", err)
					}
					// Validate the shipment belongs to the same tenant (prevent IDOR).
					if _, err := r.ShipmentRepo.GetByID(p.Context, tenantID, sid); err != nil {
						return nil, fmt.Errorf("shipment not found in tenant")
					}
					o.ShipmentID = &sid
				}
				if v, ok := input["scheduledDate"].(string); ok {
					t, err := time.Parse(time.RFC3339, v)
					if err != nil {
						return nil, fmt.Errorf("invalid scheduledDate format: %w", err)
					}
					o.ScheduledDate = &t
				}

				if err := r.OrderRepo.Update(p.Context, tenantID, id, o); err != nil {
					return nil, err
				}
				return r.OrderRepo.GetByID(p.Context, tenantID, id)
			},
		},
		"cancelOrder": &graphql.Field{
			Type: types.OrderType,
			Args: graphql.FieldConfigArgument{
				"id":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"reason": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid order id: %w", err)
				}
				reason := p.Args["reason"].(string)

				if err := r.OrderRepo.CancelOrder(p.Context, tenantID, id, reason); err != nil {
					return nil, err
				}
				return r.OrderRepo.GetByID(p.Context, tenantID, id)
			},
		},
		"returnOrder": &graphql.Field{
			Type: types.OrderType,
			Args: graphql.FieldConfigArgument{
				"id":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"reason": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid order id: %w", err)
				}
				reason := p.Args["reason"].(string)

				if err := r.OrderRepo.ReturnOrder(p.Context, tenantID, id, reason); err != nil {
					return nil, err
				}
				return r.OrderRepo.GetByID(p.Context, tenantID, id)
			},
		},
	}
}
