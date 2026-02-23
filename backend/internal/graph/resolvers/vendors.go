package resolvers

import (
	"fmt"
	"time"

	"cargomax-api/internal/graph/types"
	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

// VendorQueries returns GraphQL query fields for vendor operations.
func (r *Resolver) VendorQueries() graphql.Fields {
	return graphql.Fields{
		"vendors": &graphql.Field{
			Type: types.VendorConnectionType,
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

				items, total, err := r.VendorRepo.List(p.Context, tenantID, page, perPage)
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
		"vendor": &graphql.Field{
			Type: types.VendorType,
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
					return nil, fmt.Errorf("invalid vendor id: %w", err)
				}
				return r.VendorRepo.GetByID(p.Context, tenantID, id)
			},
		},
	}
}

// VendorMutations returns GraphQL mutation fields for vendor operations.
func (r *Resolver) VendorMutations() graphql.Fields {
	return graphql.Fields{
		"createVendor": &graphql.Field{
			Type: types.VendorType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.VendorInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				input := p.Args["input"].(map[string]interface{})

				v := &models.Vendor{
					TenantID: tenantID,
					Name:     input["name"].(string),
				}
				if val, ok := input["contactPerson"].(string); ok {
					v.ContactPerson = &val
				}
				if val, ok := input["email"].(string); ok {
					v.Email = &val
				}
				if val, ok := input["phone"].(string); ok {
					v.Phone = &val
				}
				if val, ok := input["address"].(string); ok {
					v.Address = &val
				}
				if val, ok := input["category"].(string); ok {
					v.Category = &val
				}
				if val, ok := input["rating"].(float64); ok {
					v.Rating = &val
				}
				if val, ok := input["contractStart"].(string); ok {
					t, err := time.Parse(time.RFC3339, val)
					if err != nil {
						return nil, fmt.Errorf("invalid contractStart format: %w", err)
					}
					v.ContractStart = &t
				}
				if val, ok := input["contractEnd"].(string); ok {
					t, err := time.Parse(time.RFC3339, val)
					if err != nil {
						return nil, fmt.Errorf("invalid contractEnd format: %w", err)
					}
					v.ContractEnd = &t
				}
				if val, ok := input["status"].(string); ok {
					v.Status = val
				}

				if err := r.VendorRepo.Create(p.Context, v); err != nil {
					return nil, err
				}
				return v, nil
			},
		},
		"updateVendor": &graphql.Field{
			Type: types.VendorType,
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.VendorInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid vendor id: %w", err)
				}
				input := p.Args["input"].(map[string]interface{})

				v := &models.Vendor{
					Name: input["name"].(string),
				}
				if val, ok := input["contactPerson"].(string); ok {
					v.ContactPerson = &val
				}
				if val, ok := input["email"].(string); ok {
					v.Email = &val
				}
				if val, ok := input["phone"].(string); ok {
					v.Phone = &val
				}
				if val, ok := input["address"].(string); ok {
					v.Address = &val
				}
				if val, ok := input["category"].(string); ok {
					v.Category = &val
				}
				if val, ok := input["rating"].(float64); ok {
					v.Rating = &val
				}
				if val, ok := input["contractStart"].(string); ok {
					t, err := time.Parse(time.RFC3339, val)
					if err != nil {
						return nil, fmt.Errorf("invalid contractStart format: %w", err)
					}
					v.ContractStart = &t
				}
				if val, ok := input["contractEnd"].(string); ok {
					t, err := time.Parse(time.RFC3339, val)
					if err != nil {
						return nil, fmt.Errorf("invalid contractEnd format: %w", err)
					}
					v.ContractEnd = &t
				}
				if val, ok := input["status"].(string); ok {
					v.Status = val
				}

				if err := r.VendorRepo.Update(p.Context, tenantID, id, v); err != nil {
					return nil, err
				}
				return r.VendorRepo.GetByID(p.Context, tenantID, id)
			},
		},
		"deleteVendor": &graphql.Field{
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
					return nil, fmt.Errorf("invalid vendor id: %w", err)
				}
				if err := r.VendorRepo.Delete(p.Context, tenantID, id); err != nil {
					return false, err
				}
				return true, nil
			},
		},
	}
}
