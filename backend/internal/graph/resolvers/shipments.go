package resolvers

import (
	"fmt"
	"time"

	"cargomax-api/internal/graph/types"
	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

// ShipmentQueries returns the GraphQL query fields for the shipment domain.
func (r *Resolver) ShipmentQueries() graphql.Fields {
	return graphql.Fields{
		// -----------------------------------------------------------------
		// shipments (paginated, optional status filter)
		// -----------------------------------------------------------------
		"shipments": &graphql.Field{
			Type:        types.ShipmentConnectionType,
			Description: "Returns a paginated list of shipments, optionally filtered by status.",
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
				page := p.Args["page"].(int)
				perPage := p.Args["perPage"].(int)

				if page < 1 {
					page = 1
				}
				if perPage < 1 {
					perPage = 20
				}
				if perPage > 100 {
					perPage = 100
				}

				var status string
				if s, ok := p.Args["status"].(string); ok {
					status = s
				}

				items, total, err := r.ShipmentRepo.List(p.Context, tenantID, status, page, perPage)
				if err != nil {
					return nil, fmt.Errorf("failed to list shipments: %w", err)
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

		// -----------------------------------------------------------------
		// shipment (by ID)
		// -----------------------------------------------------------------
		"shipment": &graphql.Field{
			Type:        types.ShipmentType,
			Description: "Returns a single shipment by its UUID.",
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
					return nil, fmt.Errorf("invalid shipment id: %w", err)
				}

				shipment, err := r.ShipmentRepo.GetByID(p.Context, tenantID, id)
				if err != nil {
					return nil, fmt.Errorf("shipment not found: %w", err)
				}
				return shipment, nil
			},
		},

		// -----------------------------------------------------------------
		// trackShipment (by tracking number)
		// -----------------------------------------------------------------
		"trackShipment": &graphql.Field{
			Type:        types.ShipmentType,
			Description: "Returns a shipment by its tracking number.",
			Args: graphql.FieldConfigArgument{
				"trackingNumber": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				trackingNumber := p.Args["trackingNumber"].(string)

				shipment, err := r.ShipmentRepo.GetByTracking(p.Context, tenantID, trackingNumber)
				if err != nil {
					return nil, fmt.Errorf("shipment not found: %w", err)
				}
				return shipment, nil
			},
		},

		// -----------------------------------------------------------------
		// delayedShipments
		// -----------------------------------------------------------------
		"delayedShipments": &graphql.Field{
			Type:        graphql.NewList(types.ShipmentType),
			Description: "Returns shipments that have passed their estimated delivery date without being delivered.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}

				shipments, err := r.ShipmentRepo.GetDelayed(p.Context, tenantID)
				if err != nil {
					return nil, fmt.Errorf("failed to list delayed shipments: %w", err)
				}
				return shipments, nil
			},
		},
	}
}

// ShipmentMutations returns the GraphQL mutation fields for the shipment domain.
func (r *Resolver) ShipmentMutations() graphql.Fields {
	return graphql.Fields{
		// -----------------------------------------------------------------
		// createShipment
		// -----------------------------------------------------------------
		"createShipment": &graphql.Field{
			Type:        types.ShipmentType,
			Description: "Create a new shipment for the current tenant.",
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.ShipmentInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				input := p.Args["input"].(map[string]interface{})

				now := time.Now()
				shipment := &models.Shipment{
					ID:             uuid.New(),
					TenantID:       tenantID,
					TrackingNumber: input["trackingNumber"].(string),
					Status:         "pending",
					CreatedAt:      now,
					UpdatedAt:      now,
				}

				if v, ok := input["origin"].(string); ok {
					shipment.Origin = &v
				}
				if v, ok := input["destination"].(string); ok {
					shipment.Destination = &v
				}
				if v, ok := input["status"].(string); ok && v != "" {
					shipment.Status = v
				}
				if v, ok := input["carrier"].(string); ok {
					shipment.Carrier = &v
				}
				if v, ok := input["weight"].(float64); ok {
					shipment.Weight = &v
				}
				if v, ok := input["dimensions"].(string); ok {
					shipment.Dimensions = &v
				}
				if v, ok := input["estimatedDelivery"].(string); ok && v != "" {
					if t, err := time.Parse(time.RFC3339, v); err == nil {
						shipment.EstimatedDelivery = &t
					}
				}
				if v, ok := input["actualDelivery"].(string); ok && v != "" {
					if t, err := time.Parse(time.RFC3339, v); err == nil {
						shipment.ActualDelivery = &t
					}
				}
				if v, ok := input["customerName"].(string); ok {
					shipment.CustomerName = &v
				}
				if v, ok := input["customerEmail"].(string); ok {
					shipment.CustomerEmail = &v
				}
				if v, ok := input["notes"].(string); ok {
					shipment.Notes = &v
				}

				if err := r.ShipmentRepo.Create(p.Context, shipment); err != nil {
					return nil, fmt.Errorf("failed to create shipment: %w", err)
				}
				return shipment, nil
			},
		},

		// -----------------------------------------------------------------
		// updateShipment
		// -----------------------------------------------------------------
		"updateShipment": &graphql.Field{
			Type:        types.ShipmentType,
			Description: "Update an existing shipment.",
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.ShipmentInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid shipment id: %w", err)
				}

				// Fetch existing shipment to merge updates.
				shipment, err := r.ShipmentRepo.GetByID(p.Context, tenantID, id)
				if err != nil {
					return nil, fmt.Errorf("shipment not found: %w", err)
				}

				input := p.Args["input"].(map[string]interface{})

				if v, ok := input["trackingNumber"].(string); ok {
					shipment.TrackingNumber = v
				}
				if v, ok := input["origin"].(string); ok {
					shipment.Origin = &v
				}
				if v, ok := input["destination"].(string); ok {
					shipment.Destination = &v
				}
				if v, ok := input["status"].(string); ok {
					shipment.Status = v
				}
				if v, ok := input["carrier"].(string); ok {
					shipment.Carrier = &v
				}
				if v, ok := input["weight"].(float64); ok {
					shipment.Weight = &v
				}
				if v, ok := input["dimensions"].(string); ok {
					shipment.Dimensions = &v
				}
				if v, ok := input["estimatedDelivery"].(string); ok && v != "" {
					if t, err := time.Parse(time.RFC3339, v); err == nil {
						shipment.EstimatedDelivery = &t
					}
				}
				if v, ok := input["actualDelivery"].(string); ok && v != "" {
					if t, err := time.Parse(time.RFC3339, v); err == nil {
						shipment.ActualDelivery = &t
					}
				}
				if v, ok := input["customerName"].(string); ok {
					shipment.CustomerName = &v
				}
				if v, ok := input["customerEmail"].(string); ok {
					shipment.CustomerEmail = &v
				}
				if v, ok := input["notes"].(string); ok {
					shipment.Notes = &v
				}

				shipment.UpdatedAt = time.Now()

				if err := r.ShipmentRepo.Update(p.Context, tenantID, id, shipment); err != nil {
					return nil, fmt.Errorf("failed to update shipment: %w", err)
				}
				return shipment, nil
			},
		},

		// -----------------------------------------------------------------
		// deleteShipment
		// -----------------------------------------------------------------
		"deleteShipment": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Delete a shipment by ID.",
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
					return false, fmt.Errorf("invalid shipment id: %w", err)
				}

				if err := r.ShipmentRepo.Delete(p.Context, tenantID, id); err != nil {
					return false, fmt.Errorf("failed to delete shipment: %w", err)
				}
				return true, nil
			},
		},
	}
}
