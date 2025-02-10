package server

import productPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/product"

type productServer struct {
	productPb.UnimplementedProductServer
}

func NewProductServer() productPb.ProductServer {
	return &productServer{}
}
