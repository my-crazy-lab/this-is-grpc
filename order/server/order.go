package server

import (
	"context"

	"github.com/my-crazy-lab/this-is-grpc/order/pg"
	"github.com/my-crazy-lab/this-is-grpc/proto-module/proto/order"
)

func (s *orderServer) AddToCart(ctx context.Context, req *order.AddToCartRequest) (*order.AddToCartResponse, error) {
	return pg.AddToCart(req)
}

func (s *orderServer) ViewCart(ctx context.Context, req *order.ViewCartRequest) (*order.ViewCartResponse, error) {
	return pg.ViewCart(req)
}

func (s *orderServer) PlaceOrder(ctx context.Context, req *order.PlaceOrderRequest) (*order.PlaceOrderResponse, error) {
	return pg.PlaceOrder(req)
}

func (s *orderServer) UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error) {
	return pg.UpdateOrderStatus(req)
}

func (s *orderServer) CancelOrder(ctx context.Context, req *order.CancelOrderRequest) (*order.CancelOrderResponse, error) {
	return pg.CancelOrder(req)
}

func (s *orderServer) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.OrderItem, error) {
	return pg.GetOrder(req)
}

func (s *orderServer) CreateShipping(ctx context.Context, req *order.CreateShippingRequest) (*order.CreateShippingResponse, error) {
	return pg.CreateShipping(req)
}
