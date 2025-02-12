package schema

import (
	"github.com/graphql-go/graphql"
	schemaOrder "github.com/my-crazy-lab/this-is-grpc/graphql-api-gateway/schema/order"
)

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"GetUser":       authQuery["GetUser"],
		"GetUsers":      authQuery["GetUsers"],
		"GetProduct":    productQuery["GetProduct"],
		"GetProducts":   productQuery["GetProducts"],
		"GetCategories": productQuery["GetCategories"],
		"GetReviews":    productQuery["GetReviews"],
		"ViewCart":      schemaOrder.OrderQuery["ViewCart"],
		"GetOrder":      schemaOrder.OrderQuery["GetOrder"],
	},
})

var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"SignIn":            authMutation["SignIn"],
		"SignUp":            authMutation["SignUp"],
		"CreateProduct":     productMutation["CreateProduct"],
		"CreateCategories":  productMutation["CreateCategories"],
		"UpdateInventory":   productMutation["UpdateInventory"],
		"AddToCart":         schemaOrder.OrderMutation["AddToCart"],
		"PlaceOrder":        schemaOrder.OrderMutation["PlaceOrder"],
		"UpdateOrderStatus": schemaOrder.OrderMutation["UpdateOrderStatus"],
		"CancelOrder":       schemaOrder.OrderMutation["CancelOrder"],
		"CreateShipping":    schemaOrder.OrderMutation["CreateShipping"],
	},
})
