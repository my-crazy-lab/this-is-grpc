package schemaOrder

import "github.com/graphql-go/graphql"

var cartItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "cartItemType",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"cart_id": &graphql.Field{
			Type: graphql.Int,
		},
		"product_id": &graphql.Field{
			Type: graphql.Int,
		},
		"quantity": &graphql.Field{
			Type: graphql.Int,
		},
		"created_at": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var orderItemType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "orderItemType",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"user_id": &graphql.Field{
				Type: graphql.Int,
			},
			"cart_id": &graphql.Field{
				Type: graphql.Int,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"total": &graphql.Field{
				Type: graphql.Int,
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var shippingItemType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "shippingItemType",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"order_id": &graphql.Field{
				Type: graphql.Int,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"address": &graphql.Field{
				Type: shippingAddressType,
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var shippingAddressType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "shippingAddressType",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"user_id": &graphql.Field{
				Type: graphql.Int,
			},
			"address": &graphql.Field{
				Type: graphql.String,
			},
			"city": &graphql.Field{
				Type: graphql.String,
			},
			"country": &graphql.Field{
				Type: graphql.String,
			},
			"zip_code": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var cartType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "cartType",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"user_id": &graphql.Field{
				Type: graphql.Int,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var AddToCartResponseType = cartItemType

var ViewCartResponseType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ViewCartResponseType",
		Fields: graphql.Fields{
			"cart": &graphql.Field{
				Type: cartType,
			},
			"items": &graphql.Field{
				Type: graphql.NewList(cartItemType),
			},
		},
	},
)

var PlaceOrderResponseType = orderItemType

var UpdateOrderStatusResponseType = orderItemType

var CancelOrderResponseType = orderItemType

var GetOrderResponseType = orderItemType

var CreateShippingResponseType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "CreateShippingResponseType",
		Fields: graphql.Fields{
			"shipping": &graphql.Field{
				Type: shippingItemType,
			},
		},
	},
)
