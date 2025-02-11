package schema

import (
	"context"
	"errors"
	"log"

	"github.com/graphql-go/graphql"
	client "github.com/my-crazy-lab/this-is-grpc/proto-module/client"
	"github.com/my-crazy-lab/this-is-grpc/proto-module/proto/auth"
	"github.com/my-crazy-lab/this-is-grpc/proto-module/proto/order"
	"github.com/my-crazy-lab/this-is-grpc/shared/constants"
	"google.golang.org/grpc/metadata"
)

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

var orderQuery = graphql.Fields{}

var orderMutation = graphql.Fields{
	"AddToCart": &graphql.Field{
		Type:        cartItemType,
		Description: "Add new product into cart",
		Args: graphql.FieldConfigArgument{
			"user_id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"product_id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"quantity": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			user_id, _ := params.Args["user_id"].(int)
			product_id, _ := params.Args["product_id"].(int)
			quantity, _ := params.Args["quantity"].(int)

			ctx := params.Context

			res, err := addToCart(ctx, client.AuthenticationService, &order.AddToCartRequest{
				UserId: int32(user_id), ProductId: int32(product_id), Quantity: int32(quantity),
			})
			if err != nil {
				log.Fatalf("AuthenticationClient.AddToCart(_) = _, %v: ", err)
			}

			return res.Item, nil
		},
	},
}

func addToCart(ctx context.Context, client auth.AuthClient, params *order.AddToCartRequest) (*order.AddToCartResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.TIMEOUT)
	defer cancel()

	token, ok := ctx.Value("token").(string)
	if !ok || token == "" {
		return nil, errors.New("unauthorized: missing token")
	}

	// Create gRPC metadata with the token
	md := metadata.New(map[string]string{"user-authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.AddToCart(ctx, params)
	if err != nil {
		log.Fatalf("AuthenticationClient.AddToCart(_) = _, %v: ", err)
	}

	return resp, nil
}
