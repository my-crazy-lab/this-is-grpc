package schema

import (
	"context"
	"log"
	"time"

	"github.com/graphql-go/graphql"
	authPb "github.com/my-crazy-lab/this-is-grpc/graphql-api-gateway/proto/auth"

	"github.com/my-crazy-lab/this-is-grpc/graphql-api-gateway/rpcServices"
)

var categoryItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "category",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"created_at": &graphql.Field{
			Type: graphql.String,
		},
		"updated_at": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var productItemType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Product",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"price": &graphql.Field{
				Type: graphql.Int,
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
			},
			"categories": &graphql.Field{
				Type: graphql.NewList(categoryItemType),
			},
			"quantity": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

var reviewType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Review",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"comment": &graphql.Field{
				Type: graphql.String,
			},
			"rating": &graphql.Field{
				Type: graphql.Int,
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
			},
			"user": &graphql.Field{
				Type: graphql.NewList(userType),
			},
		},
	})

var productQuery = graphql.Fields{
	"GetProducts": &graphql.Field{
		Type:        graphql.NewList(productItemType),
		Description: "Get all products",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ctx := p.Context.(context.Context)

			users := getProducts(ctx, rpcServices.AuthenticationService)
			return users, nil
		},
	},
	"GetReviews": &graphql.Field{
		Type:        graphql.NewList(reviewType),
		Description: "Get all reviews",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ctx := p.Context.(context.Context)

			users := getReviews(ctx, rpcServices.AuthenticationService)
			return users, nil
		},
	},
}

var productMutation = graphql.Fields{
	"": &graphql.Field{},
}

func getProducts(ctx context.Context, client authPb.AuthClient) *authPb.GetProductsResponse {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := client.GetProducts(ctx, &authPb.GetProductsRequest{})
	if err != nil {
		log.Fatalf("client.GetProducts(_) = _, %v: ", err)
	}

	return resp
}

func getReviews(ctx context.Context, client authPb.AuthClient) *authPb.GetReviewsResponse {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := client.GetReviews(ctx, &authPb.GetReviewsRequest{})
	if err != nil {
		log.Fatalf("client.GetReviews(_) = _, %v: ", err)
	}

	return resp
}
