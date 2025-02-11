package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/my-crazy-lab/this-is-grpc/authentication/pg"
	client "github.com/my-crazy-lab/this-is-grpc/proto-module/client"
	productPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/product"

	"google.golang.org/grpc/metadata"
)

func (s *authServer) CreateCategories(ctx context.Context, req *productPb.CreateCategoriesRequest) (*productPb.CreateCategoriesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Extract token from gRPC metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	authHeader, exists := md["user-authorization"]
	if !exists || len(authHeader) == 0 {
		return nil, errors.New("unauthorized: missing token")
	}

	token := authHeader[0]
	_, err := pg.VerifyJWT(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	resp, err := client.ProductService.CreateCategories(ctx, req)
	if err != nil {
		log.Fatalf("From ProductClient.CreateCategories(_) = _, %v: ", err)
	}

	return resp, nil
}

func (s *authServer) GetCategories(ctx context.Context, req *productPb.GetCategoriesRequest) (*productPb.GetCategoriesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	verifyToken(ctx)

	resp, err := client.ProductService.GetCategories(ctx, req)
	if err != nil {
		log.Fatalf("From ProductClient.GetCategories(_) = _, %v: ", err)
	}

	return resp, nil
}

func (s *authServer) CreateProduct(ctx context.Context, req *productPb.CreateProductRequest) (*productPb.CreateProductResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	verifyToken(ctx)

	resp, err := client.ProductService.CreateProduct(ctx, req)
	if err != nil {
		log.Fatalf("From ProductClient.CreateProduct(_) = _, %v: ", err)
	}

	return resp, nil
}

func (s *authServer) GetProducts(ctx context.Context, req *productPb.GetProductsRequest) (*productPb.GetProductsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	verifyToken(ctx)

	resp, err := client.ProductService.GetProducts(ctx, req)
	if err != nil {
		log.Fatalf("From ProductClient.GetProducts(_) = _, %v: ", err)
	}

	return resp, nil
}

func (s *authServer) GetProduct(ctx context.Context, req *productPb.GetProductRequest) (*productPb.ProductItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	verifyToken(ctx)

	resp, err := client.ProductService.GetProduct(ctx, req)
	if err != nil {
		log.Fatalf("From ProductClient.GetProduct(_) = _, %v: ", err)
	}

	return resp, nil
}
