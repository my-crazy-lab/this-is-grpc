package schemaOrder

import (
	"github.com/graphql-go/graphql"
)

var OrderQuery = graphql.Fields{
	"ViewCart": &graphql.Field{
		Description: "View cart",
		Args:        ViewCartArgs,
		Resolve:     ViewCartResolver,
		Type:        ViewCartResponseType},
	"GetOrder": &graphql.Field{
		Description: "Get order information",
		Args:        GetOrderArgs,
		Resolve:     GetOrderResolver,
		Type:        GetOrderResponseType},
}

var OrderMutation = graphql.Fields{
	"AddToCart": &graphql.Field{
		Description: "Add new product into cart",
		Args:        AddToCartArgs,
		Resolve:     AddToCartResolver,
		Type:        AddToCartResponseType,
	},
	"PlaceOrder": &graphql.Field{
		Description: "Place order",
		Args:        PlaceOrderArgs,
		Resolve:     PlaceOrderResolver,
		Type:        PlaceOrderResponseType},
	"UpdateOrderStatus": &graphql.Field{
		Description: "Update order status",
		Args:        UpdateOrderStatusArgs,
		Resolve:     UpdateOrderStatusResolver,
		Type:        UpdateOrderStatusResponseType},
	"CancelOrder": &graphql.Field{
		Description: "Cancel ordere",
		Args:        CancelOrderArgs,
		Resolve:     CancelOrderResolver,
		Type:        CancelOrderResponseType},
	"CreateShipping": &graphql.Field{
		Description: "Create shipping",
		Args:        CreateShippingArgs,
		Resolve:     CreateShippingResolver,
		Type:        CreateShippingResponseType},
}
