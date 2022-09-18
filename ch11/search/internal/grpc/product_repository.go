package grpc

import (
	"context"

	"google.golang.org/grpc"

	"eda-in-golang/search/internal/application"
	"eda-in-golang/search/internal/models"
	"eda-in-golang/stores/storespb"
)

type ProductRepository struct {
	client storespb.StoresServiceClient
}

var _ application.ProductRepository = (*ProductRepository)(nil)

func NewProductRepository(conn *grpc.ClientConn) ProductRepository {
	return ProductRepository{client: storespb.NewStoresServiceClient(conn)}
}

func (r ProductRepository) Find(ctx context.Context, productID string) (*models.Product, error) {
	resp, err := r.client.GetProduct(ctx, &storespb.GetProductRequest{Id: productID})
	if err != nil {
		return nil, err
	}

	return r.productToDomain(resp.Product), nil
}

func (r ProductRepository) productToDomain(product *storespb.Product) *models.Product {
	return &models.Product{
		ID:   product.GetId(),
		Name: product.GetName(),
	}
}
