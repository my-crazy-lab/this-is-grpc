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

func (s *authServer) CreateProduct(ctx context.Context, req *productPb.CreateProductRequest) (*productPb.CreateProductResponse, error) {
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

	resp, err := client.ProductService.CreateProduct(ctx, req)
	if err != nil {
		log.Fatalf("From ProductClient.CreateProduct(_) = _, %v: ", err)
	}

	return resp, nil
}
