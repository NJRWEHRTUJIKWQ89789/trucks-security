package graph

import (
	"cargomax-api/internal/graph/resolvers"

	"github.com/graphql-go/graphql"
)

// NewSchema assembles the root GraphQL schema by merging all query and
// mutation fields provided by the resolver's domain-specific methods.
func NewSchema(r *resolvers.Resolver) (graphql.Schema, error) {
	queryFields := graphql.Fields{}
	mutationFields := graphql.Fields{}

	// Merge query fields from every domain resolver.
	for k, v := range r.AuthQueries() {
		queryFields[k] = v
	}
	for k, v := range r.DashboardQueries() {
		queryFields[k] = v
	}
	for k, v := range r.ShipmentQueries() {
		queryFields[k] = v
	}
	for k, v := range r.FleetQueries() {
		queryFields[k] = v
	}
	for k, v := range r.WarehouseQueries() {
		queryFields[k] = v
	}
	for k, v := range r.OrderQueries() {
		queryFields[k] = v
	}
	for k, v := range r.VendorQueries() {
		queryFields[k] = v
	}
	for k, v := range r.ClientQueries() {
		queryFields[k] = v
	}
	for k, v := range r.ReportQueries() {
		queryFields[k] = v
	}
	for k, v := range r.SettingsQueries() {
		queryFields[k] = v
	}

	// Merge mutation fields from every domain resolver.
	for k, v := range r.AuthMutations() {
		mutationFields[k] = v
	}
	for k, v := range r.ShipmentMutations() {
		mutationFields[k] = v
	}
	for k, v := range r.FleetMutations() {
		mutationFields[k] = v
	}
	for k, v := range r.WarehouseMutations() {
		mutationFields[k] = v
	}
	for k, v := range r.OrderMutations() {
		mutationFields[k] = v
	}
	for k, v := range r.VendorMutations() {
		mutationFields[k] = v
	}
	for k, v := range r.ClientMutations() {
		mutationFields[k] = v
	}
	for k, v := range r.SettingsMutations() {
		mutationFields[k] = v
	}

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: queryFields,
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: mutationFields,
		}),
	})
}
