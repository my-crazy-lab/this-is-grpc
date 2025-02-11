package server

import (
	"context"

	"github.com/my-crazy-lab/this-is-grpc/order/pg"
	"github.com/my-crazy-lab/this-is-grpc/proto-module/proto/order"
)

func (s *orderServer) AddToCart(ctx context.Context, req *order.AddToCartRequest) (*order.AddToCartResponse, error) {
	res, err := pg.AddToCart(req.ProductId, req.UserId, req.Quantity)
	if err != nil {
		return nil, err
	}

	return res, nil
}
