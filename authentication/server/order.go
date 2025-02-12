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

func (s *authServer) ViewCart(ctx context.Context, req *order.ViewCartRequest) (*order.ViewCartResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	verifyToken(ctx)

	resp, err := client.OrderService.ViewCart(ctx, req)
	if err != nil {
		log.Fatalf("From OrderClient.ViewCart(_) = _, %v: ", err)
	}

	return resp, nil
}

func (s *authServer) PlaceOrder(ctx context.Context, req *order.PlaceOrderRequest) (*order.PlaceOrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	verifyToken(ctx)

	resp, err := client.OrderService.PlaceOrder(ctx, req)
	if err != nil {
		log.Fatalf("From OrderClient.PlaceOrder(_) = _, %v: ", err)
	}

	return resp, nil
}

func (s *authServer) UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	verifyToken(ctx)

	resp, err := client.OrderService.UpdateOrderStatus(ctx, req)
	if err != nil {
		log.Fatalf("From OrderClient.UpdateOrderStatus(_) = _, %v: ", err)
	}

	return resp, nil
}

func (s *authServer) CancelOrder(ctx context.Context, req *order.CancelOrderRequest) (*order.CancelOrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	verifyToken(ctx)

	resp, err := client.OrderService.CancelOrder(ctx, req)
	if err != nil {
		log.Fatalf("From OrderClient.CancelOrder(_) = _, %v: ", err)
	}

	return resp, nil
}

func (s *authServer) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.OrderItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	verifyToken(ctx)

	resp, err := client.OrderService.GetOrder(ctx, req)
	if err != nil {
		log.Fatalf("From OrderClient.GetOrder(_) = _, %v: ", err)
	}

	return resp, nil
}

func (s *authServer) CreateShipping(ctx context.Context, req *order.CreateShippingRequest) (*order.CreateShippingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	verifyToken(ctx)

	resp, err := client.OrderService.CreateShipping(ctx, req)
	if err != nil {
		log.Fatalf("From OrderClient.CreateShipping(_) = _, %v: ", err)
	}

	return resp, nil
}
