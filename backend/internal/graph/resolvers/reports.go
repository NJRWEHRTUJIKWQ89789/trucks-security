package resolvers

import (
	"cargomax-api/internal/graph/types"

	"github.com/graphql-go/graphql"
)

// ReportQueries returns GraphQL query fields for report operations.
func (r *Resolver) ReportQueries() graphql.Fields {
	return graphql.Fields{
		"revenueReport": &graphql.Field{
			Type: types.RevenueReportType,
			Args: graphql.FieldConfigArgument{
				"year": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				year := p.Args["year"].(int)
				return r.ReportRepo.GetRevenueReport(p.Context, tenantID, year)
			},
		},
		"deliveryReport": &graphql.Field{
			Type: types.DeliveryReportType,
			Args: graphql.FieldConfigArgument{
				"year": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				year := p.Args["year"].(int)
				return r.ReportRepo.GetDeliveryReport(p.Context, tenantID, year)
			},
		},
		"fleetReport": &graphql.Field{
			Type: types.FleetReportType,
			Args: graphql.FieldConfigArgument{
				"year": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tenantID, err := requireTenant(p.Context)
				if err != nil {
					return nil, err
				}
				year := p.Args["year"].(int)
				return r.ReportRepo.GetFleetReport(p.Context, tenantID, year)
			},
		},
	}
}
