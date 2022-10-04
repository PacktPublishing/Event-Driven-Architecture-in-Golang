package grpc

import (
	"context"

	"google.golang.org/grpc"

	"eda-in-golang/internal/rpc"
	"eda-in-golang/search/internal/application"
	"eda-in-golang/search/internal/models"
	"eda-in-golang/stores/storespb"
)

type ProductRepository struct {
	endpoint string
}

var _ application.ProductRepository = (*ProductRepository)(nil)

func NewProductRepository(endpoint string) ProductRepository {
	return ProductRepository{
		endpoint: endpoint,
	}
}

func (r ProductRepository) Find(ctx context.Context, productID string) (product *models.Product, err error) {
	var conn *grpc.ClientConn
	conn, err = r.dial(ctx)
	if err != nil {
		return nil, err
	}

	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	resp, err := storespb.NewStoresServiceClient(conn).GetProduct(ctx, &storespb.GetProductRequest{Id: productID})
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

func (r ProductRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}
