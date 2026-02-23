package resolvers

import (
	"fmt"

	"cargomax-api/internal/graph/types"
	"cargomax-api/internal/models"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

func (r *Resolver) SettingsQueries() graphql.Fields {
	return graphql.Fields{
		"settings": &graphql.Field{
			Type: graphql.NewList(types.SettingType),
			Args: graphql.FieldConfigArgument{
				"category": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				var category string
				if v, ok := p.Args["category"].(string); ok {
					category = v
				}
				return r.SettingRepo.GetByCategory(p.Context, tenantID, category)
			},
		},
		"roles": &graphql.Field{
			Type: graphql.NewList(types.RoleType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				roles, _, err := r.RoleRepo.List(p.Context, tenantID, 1, 100)
				if err != nil {
					return nil, err
				}
				return roles, nil
			},
		},
		"notifications": &graphql.Field{
			Type: types.NotificationConnectionType,
			Args: graphql.FieldConfigArgument{
				"page":    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
				"perPage": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 20},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, userID, err := requireAuth(p.Context)
				if err != nil {
					return nil, err
				}
				page := p.Args["page"].(int)
				perPage := p.Args["perPage"].(int)
				items, total, err := r.NotificationRepo.List(p.Context, tenantID, userID, page, perPage)
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{
					"items": items, "totalCount": total,
					"page": page, "perPage": perPage,
					"totalPages": (total + perPage - 1) / perPage,
				}, nil
			},
		},
		"notificationPreferences": &graphql.Field{
			Type: graphql.NewList(types.NotificationPreferenceType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, userID, err := requireAuth(p.Context)
				if err != nil {
					return nil, err
				}
				return r.NotificationRepo.GetPreferences(p.Context, tenantID, userID)
			},
		},
	}
}

func (r *Resolver) SettingsMutations() graphql.Fields {
	return graphql.Fields{
		"updateSetting": &graphql.Field{
			Type: types.SettingType,
			Args: graphql.FieldConfigArgument{
				"key":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"value": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, userID, err := requireAuth(p.Context)
				if err != nil {
					return nil, err
				}
				key := p.Args["key"].(string)
				value := p.Args["value"].(string)
				s := &models.Setting{ID: uuid.New(), TenantID: tenantID, Key: key, Value: &value, UpdatedBy: &userID}
				if err := r.SettingRepo.Set(p.Context, s); err != nil {
					return nil, err
				}
				return r.SettingRepo.Get(p.Context, tenantID, key)
			},
		},
		"updateRole": &graphql.Field{
			Type: types.RoleType,
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.RoleInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid role id: %w", err)
				}
				input := p.Args["input"].(map[string]interface{})
				name := input["name"].(string)
				var permissions string
				if v, ok := input["permissions"].(string); ok {
					permissions = v
				}
				role := &models.Role{ID: id, TenantID: tenantID, Name: name, Permissions: permissions}
				if err := r.RoleRepo.Update(p.Context, tenantID, id, role); err != nil {
					return nil, err
				}
				return r.RoleRepo.GetByID(p.Context, tenantID, id)
			},
		},
		"updateNotificationPreference": &graphql.Field{
			Type: types.NotificationPreferenceType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.NotificationPrefInputType)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, userID, err := requireAuth(p.Context)
				if err != nil {
					return nil, err
				}
				input := p.Args["input"].(map[string]interface{})
				eventType := input["eventType"].(string)
				emailEnabled, smsEnabled, pushEnabled := true, false, false
				if v, ok := input["emailEnabled"].(bool); ok { emailEnabled = v }
				if v, ok := input["smsEnabled"].(bool); ok { smsEnabled = v }
				if v, ok := input["pushEnabled"].(bool); ok { pushEnabled = v }
				pref := &models.NotificationPreference{
					ID: uuid.New(), TenantID: tenantID, UserID: userID,
					EventType: eventType, EmailEnabled: emailEnabled,
					SMSEnabled: smsEnabled, PushEnabled: pushEnabled,
				}
				if err := r.NotificationRepo.UpdatePreference(p.Context, pref); err != nil {
					return nil, err
				}
				return pref, nil
			},
		},
		"markNotificationRead": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, userID, err := requireAuth(p.Context)
				if err != nil {
					return nil, err
				}
				id, err := uuid.Parse(p.Args["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid notification id: %w", err)
				}
				if _, err := r.NotificationRepo.GetByID(p.Context, tenantID, userID, id); err != nil {
					return nil, fmt.Errorf("notification not found: %w", err)
				}
				if err := r.NotificationRepo.MarkRead(p.Context, tenantID, userID, id); err != nil {
					return nil, err
				}
				return true, nil
			},
		},
	}
}
