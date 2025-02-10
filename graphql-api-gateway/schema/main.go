package schema

import "github.com/graphql-go/graphql"

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"getUsers":    authQuery["getUsers"],
		"getProducts": productQuery["getProducts"],
		"getReviews":  productQuery["getReviews"],
	},
})

var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"signIn": authMutation["signIn"],
		"signUp": authMutation["signUp"],
	},
})
