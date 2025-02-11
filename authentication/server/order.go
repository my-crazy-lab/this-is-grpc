package server

import (
	"context"
	"log"
	"time"

	"github.com/my-crazy-lab/this-is-grpc/proto-module/client"
	"github.com/my-crazy-lab/this-is-grpc/proto-module/proto/order"
)

func (s *authServer) AddToCart(ctx context.Context, req *order.AddToCartRequest) (*order.AddToCartResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	verifyToken(ctx)

	resp, err := client.OrderService.AddToCart(ctx, req)
	if err != nil {
		log.Fatalf("From OrderClient.AddToCart(_) = _, %v: ", err)
	}

	return resp, nil
}
