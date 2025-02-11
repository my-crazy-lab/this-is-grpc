package server

import "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/order"

type orderServer struct {
	order.UnimplementedOrderServer
}

func NewOrderServer() order.OrderServer {
	return &orderServer{}
}
