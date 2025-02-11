package server

import (
	"context"

	"github.com/my-crazy-lab/this-is-grpc/product/pg"
	productPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/product"
)

func (s *productServer) CreateProduct(ctx context.Context, req *productPb.CreateProductRequest) (*productPb.CreateProductResponse, error) {
	id, err := pg.CreateProduct(req.Name, req.Description, req.Price, &req.CategoryIds, req.Quantity)

	if err != nil {
		return nil, err
	}

	return &productPb.CreateProductResponse{Id: id}, nil
}

func (s *productServer) CreateCategories(ctx context.Context, req *productPb.CreateCategoriesRequest) (*productPb.CreateCategoriesResponse, error) {
	err := pg.InsertCategories(req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	return &productPb.CreateCategoriesResponse{}, nil
}

func (s *productServer) GetCategories(ctx context.Context, req *productPb.GetCategoriesRequest) (*productPb.GetCategoriesResponse, error) {
	categories, err := pg.GetCategories()
	if err != nil {
		return nil, err
	}

	return &productPb.GetCategoriesResponse{Categories: categories}, nil
}

func (s *productServer) GetProducts(ctx context.Context, req *productPb.GetProductsRequest) (*productPb.GetProductsResponse, error) {
	products, total, err := pg.GetProducts(req.Pagination.PageSize, req.Pagination.PageIndex)
	if err != nil {
		return nil, err
	}

	return &productPb.GetProductsResponse{Products: products, Total: total}, nil
}

func (s *productServer) GetProduct(ctx context.Context, req *productPb.GetProductRequest) (*productPb.ProductItem, error) {
	product, err := pg.GetProduct(req.ProductId)
	if err != nil {
		return nil, err
	}

	return product, nil
}
