package schema

import "github.com/graphql-go/graphql"

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"GetUsers":      authQuery["GetUsers"],
		"GetProducts":   productQuery["GetProducts"],
		"GetCategories": productQuery["GetCategories"],
		"GetReviews":    productQuery["GetReviews"],
	},
})

var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"SignIn":           authMutation["SignIn"],
		"SignUp":           authMutation["SignUp"],
		"CreateProduct":    productMutation["CreateProduct"],
		"CreateCategories": productMutation["CreateCategories"],
	},
})
