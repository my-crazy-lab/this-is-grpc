package schemaOrder

import "github.com/graphql-go/graphql"

var AddToCartArgs = graphql.FieldConfigArgument{
	"user_id": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
	"product_id": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
	"quantity": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
}

var ViewCartArgs = graphql.FieldConfigArgument{
	"user_id": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
}

var PlaceOrderArgs = graphql.FieldConfigArgument{
	"user_id": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
	"cart_id": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
}

var UpdateOrderStatusArgs = graphql.FieldConfigArgument{
	"order_id": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
	"status": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
}

var CancelOrderArgs = graphql.FieldConfigArgument{
	"order_id": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
}

var GetOrderArgs = graphql.FieldConfigArgument{
	"order_id": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
}

var ShippingAddressInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "ShippingAddressInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"id": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"user_id": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"address": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"city": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"state": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"country": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"zip_code": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	},
)

var CreateShippingArgs = graphql.FieldConfigArgument{
	"order_id": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
	"address": &graphql.ArgumentConfig{
		Type: ShippingAddressInput,
	},
}
