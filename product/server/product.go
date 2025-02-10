package server

import (
	"context"

	"github.com/my-crazy-lab/this-is-grpc/product/pg"
	productPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/product"
)

func (s *productServer) CreateCategories(ctx context.Context, req *productPb.CreateCategoriesRequest) (*productPb.CreateCategoriesResponse, error) {
	err := pg.InsertCategories(req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	return &productPb.CreateCategoriesResponse{}, nil
}
