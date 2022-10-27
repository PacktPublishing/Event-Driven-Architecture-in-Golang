package grpc

import (
	"context"

	"github.com/stackus/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"eda-in-golang/baskets/internal/domain"
	"eda-in-golang/internal/rpc"
	"eda-in-golang/stores/storespb"
)

type ProductRepository struct {
	endpoint string
}

var _ domain.ProductRepository = (*ProductRepository)(nil)

func NewProductRepository(endpoint string) ProductRepository {
	return ProductRepository{
		endpoint: endpoint,
	}
}

func (r ProductRepository) Find(ctx context.Context, productID string) (product *domain.Product, err error) {
	var conn *grpc.ClientConn
	conn, err = r.dial(ctx)
	if err != nil {
		return nil, err
	}

	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	resp, err := storespb.NewStoresServiceClient(conn).GetProduct(ctx, &storespb.GetProductRequest{
		Id: productID,
	})
	if err != nil {
		if errors.GRPCCode(err) == codes.NotFound {
			return nil, errors.ErrNotFound.Msg("product was not located")
		}
		return nil, errors.Wrap(err, "requesting product")
	}

	return r.productToDomain(resp.Product), nil
}

func (r ProductRepository) productToDomain(product *storespb.Product) *domain.Product {
	return &domain.Product{
		ID:      product.GetId(),
		StoreID: product.GetStoreId(),
		Name:    product.GetName(),
		Price:   product.GetPrice(),
	}
}

func (r ProductRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}
