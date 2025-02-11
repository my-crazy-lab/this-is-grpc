package schema

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
	"github.com/my-crazy-lab/this-is-grpc/proto-module/proto/auth"
	"github.com/my-crazy-lab/this-is-grpc/proto-module/proto/product"
	"google.golang.org/grpc/metadata"

	"github.com/my-crazy-lab/this-is-grpc/proto-module/client"
	"github.com/my-crazy-lab/this-is-grpc/shared/constants"
	"github.com/my-crazy-lab/this-is-grpc/shared/utils"
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

var getProductResponse = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "GetProductResponse",
		Fields: graphql.Fields{
			"total": &graphql.Field{
				Type: graphql.Int,
			},
			"products": &graphql.Field{
				Type: graphql.NewList(productItemType),
			},
		},
	},
)

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
				Type: graphql.Float,
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
	"GetProduct": &graphql.Field{
		Type:        productItemType,
		Description: "Get product by ud",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ctx := params.Context.(context.Context)

			id := params.Args["id"].(int)

			return getProductById(ctx, client.AuthenticationService, int32(id)), nil
		},
	},
	"GetProducts": &graphql.Field{
		Type:        getProductResponse,
		Description: "Get all products",
		Args: graphql.FieldConfigArgument{
			"pageSize":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			"pageIndex": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ctx := params.Context.(context.Context)

			pageSize := params.Args["pageSize"].(int)
			pageIndex := params.Args["pageIndex"].(int)

			return getProducts(ctx, client.AuthenticationService, pageSize, pageIndex), nil
		},
	},
	"GetReviews": &graphql.Field{
		Type:        graphql.NewList(reviewType),
		Description: "Get all reviews",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ctx := p.Context.(context.Context)

			return getReviews(ctx, client.AuthenticationService), nil
		},
	},
	"GetCategories": &graphql.Field{
		Type:        graphql.NewList(categoryItemType),
		Description: "Get all categories",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ctx := p.Context.(context.Context)
			categories := getCategories(ctx, client.AuthenticationService).Categories

			res := make([]CategoriesResponse, len(categories))

			for _, category := range categories {
				var r CategoriesResponse

				r.Name = category.Name
				r.Id = category.Id
				r.Description = category.Description
				r.CreatedAt = utils.PbTimestampToISO(category.CreatedAt)
				r.UpdatedAt = utils.PbTimestampToISO(category.UpdatedAt)

				res = append(res, r)
			}

			return res, nil
		},
	},
}

type CategoriesResponse struct {
	Id          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

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
				Type: graphql.Float,
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
			quantity, _ := params.Args["quantity"].(int)

			rawCategories, _ := params.Args["categories"].([]interface{})

			categories := make([]int32, len(rawCategories))

			for i, v := range rawCategories {
				fmt.Printf("Index %d: Type: %T, Value: %v\n", i, v, v) // Debugging

				switch val := v.(type) {
				case int:
					categories[i] = int32(val) // Convert int to int32
				case float64:
					categories[i] = int32(val) // Convert float64 to int32 (for safety)
				default:
					fmt.Println("Warning: Unexpected category type:", v)
				}
			}

			ctx := params.Context

			res := createProduct(ctx, client.AuthenticationService, &product.CreateProductRequest{
				Name: name, Description: description, Price: price, Quantity: int32(quantity), CategoryIds: categories,
			})

			return res, nil
		},
	},
	"CreateCategories": &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name: "CategoryIdInserted",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
			},
		}),
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

			res, err := createCategories(ctx, client.AuthenticationService, &product.CreateCategoriesRequest{
				Name: name, Description: description,
			})
			if err != nil {
				log.Fatalf("AuthenticationClient.CreateCategories(_) = _, %v: ", err)
			}

			return res, nil
		},
	},
}

func getProductById(ctx context.Context, client auth.AuthClient, id int32) *product.ProductItem {
	ctx, cancel := context.WithTimeout(ctx, constants.TIMEOUT)
	defer cancel()

	resp, err := client.GetProduct(ctx, &product.GetProductRequest{ProductId: id})
	if err != nil {
		log.Fatalf("AuthenticationClient.GetProduct(_) = _, %v: ", err)
	}

	return resp
}

func getProducts(ctx context.Context, client auth.AuthClient, pageSize, pageIndex int) *product.GetProductsResponse {
	ctx, cancel := context.WithTimeout(ctx, constants.TIMEOUT)
	defer cancel()

	resp, err := client.GetProducts(ctx, &product.GetProductsRequest{Pagination: &product.Pagination{PageSize: int32(pageSize), PageIndex: int32(pageIndex)}})
	if err != nil {
		log.Fatalf("AuthenticationClient.GetProducts(_) = _, %v: ", err)
	}

	return resp
}

func getReviews(ctx context.Context, client auth.AuthClient) *product.GetReviewsResponse {
	ctx, cancel := context.WithTimeout(ctx, constants.TIMEOUT)
	defer cancel()

	resp, err := client.GetReviews(ctx, &product.GetReviewsRequest{})
	if err != nil {
		log.Fatalf("AuthenticationClient.GetReviews(_) = _, %v: ", err)
	}

	return resp
}

func getCategories(ctx context.Context, client auth.AuthClient) *product.GetCategoriesResponse {
	ctx, cancel := context.WithTimeout(ctx, constants.TIMEOUT)
	defer cancel()

	resp, err := client.GetCategories(ctx, &product.GetCategoriesRequest{})
	if err != nil {
		log.Fatalf("AuthenticationClient.GetReviews(_) = _, %v: ", err)
	}

	return resp
}

func createProduct(ctx context.Context, client auth.AuthClient, params *product.CreateProductRequest) *product.CreateProductResponse {
	ctx, cancel := context.WithTimeout(ctx, constants.TIMEOUT)
	defer cancel()

	resp, err := client.CreateProduct(ctx, params)
	if err != nil {
		log.Fatalf("AuthenticationClient.CreateProduct(_) = _, %v: ", err)
	}

	return resp
}

func createCategories(ctx context.Context, client auth.AuthClient, params *product.CreateCategoriesRequest) (*product.CreateCategoriesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.TIMEOUT)
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
