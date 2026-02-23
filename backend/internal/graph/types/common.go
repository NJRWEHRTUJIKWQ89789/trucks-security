package types

import "github.com/graphql-go/graphql"

// ConnectionType creates a paginated connection wrapper for any GraphQL object type.
func ConnectionType(name string, itemType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: name,
		Fields: graphql.Fields{
			"items":      &graphql.Field{Type: graphql.NewList(itemType)},
			"totalCount": &graphql.Field{Type: graphql.Int},
			"page":       &graphql.Field{Type: graphql.Int},
			"perPage":    &graphql.Field{Type: graphql.Int},
			"totalPages": &graphql.Field{Type: graphql.Int},
		},
	})
}

// PaginationArgs provides standard page/perPage arguments for list queries.
var PaginationArgs = graphql.FieldConfigArgument{
	"page":    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
	"perPage": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 20},
}

// PageInfoType describes pagination metadata for a connection.
var PageInfoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PageInfo",
	Fields: graphql.Fields{
		"totalCount": &graphql.Field{Type: graphql.Int},
		"page":       &graphql.Field{Type: graphql.Int},
		"perPage":    &graphql.Field{Type: graphql.Int},
		"totalPages": &graphql.Field{Type: graphql.Int},
	},
})
