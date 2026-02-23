package resolvers

import (
	"fmt"

	"cargomax-api/internal/graph/types"
	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

// ClientQueries returns GraphQL query fields for client and feedback operations.
func (r *Resolver) ClientQueries() graphql.Fields {
	return graphql.Fields{
		"clients": &graphql.Field{
			Type: types.ClientConnectionType,
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

				items, total, err := r.ClientRepo.List(p.Context, tenantID, page, perPage)
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
		"client": &graphql.Field{
			Type: types.ClientType,
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
					return nil, fmt.Errorf("invalid client id: %w", err)
				}
				return r.ClientRepo.GetByID(p.Context, tenantID, id)
			},
		},
		"feedbacks": &graphql.Field{
			Type: types.FeedbackConnectionType,
			Args: graphql.FieldConfigArgument{
				"page":    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
				"perPage": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 50},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				page := p.Args["page"].(int)
				perPage := p.Args["perPage"].(int)

				items, total, err := r.FeedbackRepo.ListAll(p.Context, tenantID, page, perPage)
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
		"clientFeedback": &graphql.Field{
			Type: types.FeedbackConnectionType,
			Args: graphql.FieldConfigArgument{
				"clientId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"page":     &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
				"perPage":  &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 20},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				clientID, err := uuid.Parse(p.Args["clientId"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid client id: %w", err)
				}
				page := p.Args["page"].(int)
				perPage := p.Args["perPage"].(int)

				items, total, err := r.FeedbackRepo.ListByClient(p.Context, tenantID, clientID, page, perPage)
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
	}
}

// ClientMutations returns GraphQL mutation fields for client and feedback operations.
func (r *Resolver) ClientMutations() graphql.Fields {
	return graphql.Fields{
		"createClient": &graphql.Field{
			Type: types.ClientType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.ClientInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				input := p.Args["input"].(map[string]interface{})

				c := &models.Client{
					TenantID:    tenantID,
					CompanyName: input["companyName"].(string),
				}
				if v, ok := input["contactPerson"].(string); ok {
					c.ContactPerson = &v
				}
				if v, ok := input["email"].(string); ok {
					c.Email = &v
				}
				if v, ok := input["phone"].(string); ok {
					c.Phone = &v
				}
				if v, ok := input["address"].(string); ok {
					c.Address = &v
				}
				if v, ok := input["industry"].(string); ok {
					c.Industry = &v
				}
				if v, ok := input["totalShipments"].(int); ok {
					c.TotalShipments = v
				}
				if v, ok := input["totalSpent"].(float64); ok {
					c.TotalSpent = &v
				}
				if v, ok := input["satisfactionRating"].(float64); ok {
					c.SatisfactionRating = &v
				}
				if v, ok := input["status"].(string); ok {
					c.Status = v
				}

				if err := r.ClientRepo.Create(p.Context, c); err != nil {
					return nil, err
				}
				return c, nil
			},
		},
		"updateClient": &graphql.Field{
			Type: types.ClientType,
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.ClientInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid client id: %w", err)
				}
				input := p.Args["input"].(map[string]interface{})

				c := &models.Client{
					CompanyName: input["companyName"].(string),
				}
				if v, ok := input["contactPerson"].(string); ok {
					c.ContactPerson = &v
				}
				if v, ok := input["email"].(string); ok {
					c.Email = &v
				}
				if v, ok := input["phone"].(string); ok {
					c.Phone = &v
				}
				if v, ok := input["address"].(string); ok {
					c.Address = &v
				}
				if v, ok := input["industry"].(string); ok {
					c.Industry = &v
				}
				if v, ok := input["totalShipments"].(int); ok {
					c.TotalShipments = v
				}
				if v, ok := input["totalSpent"].(float64); ok {
					c.TotalSpent = &v
				}
				if v, ok := input["satisfactionRating"].(float64); ok {
					c.SatisfactionRating = &v
				}
				if v, ok := input["status"].(string); ok {
					c.Status = v
				}

				if err := r.ClientRepo.Update(p.Context, tenantID, id, c); err != nil {
					return nil, err
				}
				return r.ClientRepo.GetByID(p.Context, tenantID, id)
			},
		},
		"deleteClient": &graphql.Field{
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
					return nil, fmt.Errorf("invalid client id: %w", err)
				}
				if err := r.ClientRepo.Delete(p.Context, tenantID, id); err != nil {
					return false, err
				}
				return true, nil
			},
		},
		"submitFeedback": &graphql.Field{
			Type: types.FeedbackType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.FeedbackInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				input := p.Args["input"].(map[string]interface{})

				clientID, err := uuid.Parse(input["clientId"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid client id: %w", err)
				}

				// Validate the client belongs to the same tenant (prevent IDOR).
				if _, err := r.ClientRepo.GetByID(p.Context, tenantID, clientID); err != nil {
					return nil, fmt.Errorf("client not found in tenant")
				}

				f := &models.Feedback{
					TenantID: tenantID,
					ClientID: clientID,
					Rating:   input["rating"].(int),
				}
				if v, ok := input["comment"].(string); ok {
					f.Comment = &v
				}
				if v, ok := input["category"].(string); ok {
					f.Category = &v
				}

				if err := r.FeedbackRepo.Create(p.Context, f); err != nil {
					return nil, err
				}
				return f, nil
			},
		},
	}
}
