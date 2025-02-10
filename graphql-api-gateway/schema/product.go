package schema

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/graphql-go/graphql"
	authPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/auth"
	productPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/product"
	"google.golang.org/grpc/metadata"

	client "github.com/my-crazy-lab/this-is-grpc/proto-module/client"
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

			users := getProducts(ctx, client.AuthenticationService)
			return users, nil
		},
	},
	"GetReviews": &graphql.Field{
		Type:        graphql.NewList(reviewType),
		Description: "Get all reviews",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ctx := p.Context.(context.Context)

			users := getReviews(ctx, client.AuthenticationService)
			return users, nil
		},
	},
}

var test = graphql.NewObject(graphql.ObjectConfig{
	Name: "Test",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var productMutation = graphql.Fields{
	"CreateProduct": &graphql.Field{
		Type:        productItemType,
		Description: "Create new product",
		Args: graphql.FieldConfigArgument{
			"name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"description": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"price": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"quantity": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"categories": &graphql.ArgumentConfig{
				Type: graphql.NewList(graphql.Int),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			name, _ := params.Args["name"].(string)
			description, _ := params.Args["description"].(string)
			price, _ := params.Args["price"].(float64)
			quantity, _ := params.Args["quantity"].(int32)

			rawCategories, _ := params.Args["categories"].([]interface{})
			categories := make([]int32, len(rawCategories))

			for i, v := range rawCategories {
				categories[i], _ = v.(int32) // Convert interface{} to int
			}

			ctx := params.Context

			res := createProduct(ctx, client.AuthenticationService, &productPb.CreateProductRequest{
				Name: name, Description: description, Price: price, Quantity: quantity, CategoryIds: categories,
			})

			return res, nil
		},
	},
	"CreateCategories": &graphql.Field{
		Type:        test,
		Description: "Create new categories",
		Args: graphql.FieldConfigArgument{
			"name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"description": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			name, _ := params.Args["name"].(string)
			description, _ := params.Args["description"].(string)

			ctx := params.Context

			res, err := createCategories(ctx, client.AuthenticationService, &productPb.CreateCategoriesRequest{
				Name: name, Description: description,
			})
			if err != nil {
				log.Fatalf("AuthenticationClient.CreateCategories(_) = _, %v: ", err)
			}

			return res, nil
		},
	},
}

func getProducts(ctx context.Context, client authPb.AuthClient) *productPb.GetProductsResponse {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := client.GetProducts(ctx, &productPb.GetProductsRequest{})
	if err != nil {
		log.Fatalf("AuthenticationClient.GetProducts(_) = _, %v: ", err)
	}

	return resp
}

func getReviews(ctx context.Context, client authPb.AuthClient) *productPb.GetReviewsResponse {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := client.GetReviews(ctx, &productPb.GetReviewsRequest{})
	if err != nil {
		log.Fatalf("AuthenticationClient.GetReviews(_) = _, %v: ", err)
	}

	return resp
}

func createProduct(ctx context.Context, client authPb.AuthClient, params *productPb.CreateProductRequest) *productPb.CreateProductResponse {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := client.CreateProduct(ctx, params)
	if err != nil {
		log.Fatalf("AuthenticationClient.CreateProduct(_) = _, %v: ", err)
	}

	return resp
}

func createCategories(ctx context.Context, client authPb.AuthClient, params *productPb.CreateCategoriesRequest) (*productPb.CreateCategoriesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	token, ok := ctx.Value("token").(string)
	if !ok || token == "" {
		return nil, errors.New("unauthorized: missing token")
	}

	// Create gRPC metadata with the token
	md := metadata.New(map[string]string{"user-authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.CreateCategories(ctx, params)
	if err != nil {
		log.Fatalf("AuthenticationClient.CreateCategories(_) = _, %v: ", err)
	}

	return resp, nil
}
