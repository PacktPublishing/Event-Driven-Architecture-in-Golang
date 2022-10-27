package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"eda-in-golang/baskets/basketspb"
	"eda-in-golang/baskets/internal/application"
	"eda-in-golang/baskets/internal/domain"
)

type server struct {
	app application.App
	basketspb.UnimplementedBasketServiceServer
}

var _ basketspb.BasketServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	basketspb.RegisterBasketServiceServer(registrar, server{app: app})
	return nil
}

func (s server) StartBasket(ctx context.Context, request *basketspb.StartBasketRequest) (*basketspb.StartBasketResponse, error) {
	basketID := uuid.New().String()
	err := s.app.StartBasket(ctx, application.StartBasket{
		ID:         basketID,
		CustomerID: request.GetCustomerId(),
	})

	return &basketspb.StartBasketResponse{Id: basketID}, err
}

func (s server) CancelBasket(ctx context.Context, request *basketspb.CancelBasketRequest) (*basketspb.CancelBasketResponse, error) {
	err := s.app.CancelBasket(ctx, application.CancelBasket{
		ID: request.GetId(),
	})

	return &basketspb.CancelBasketResponse{}, err
}

func (s server) CheckoutBasket(ctx context.Context, request *basketspb.CheckoutBasketRequest) (*basketspb.CheckoutBasketResponse, error) {
	err := s.app.CheckoutBasket(ctx, application.CheckoutBasket{
		ID:        request.GetId(),
		PaymentID: request.GetPaymentId(),
	})

	return &basketspb.CheckoutBasketResponse{}, err
}

func (s server) AddItem(ctx context.Context, request *basketspb.AddItemRequest) (*basketspb.AddItemResponse, error) {
	err := s.app.AddItem(ctx, application.AddItem{
		ID:        request.GetId(),
		ProductID: request.GetProductId(),
		Quantity:  int(request.GetQuantity()),
	})

	return &basketspb.AddItemResponse{}, err
}

func (s server) RemoveItem(ctx context.Context, request *basketspb.RemoveItemRequest) (*basketspb.RemoveItemResponse, error) {
	err := s.app.RemoveItem(ctx, application.RemoveItem{
		ID:        request.GetId(),
		ProductID: request.GetProductId(),
		Quantity:  int(request.GetQuantity()),
	})

	return &basketspb.RemoveItemResponse{}, err
}

func (s server) GetBasket(ctx context.Context, request *basketspb.GetBasketRequest) (*basketspb.GetBasketResponse, error) {
	basket, err := s.app.GetBasket(ctx, application.GetBasket{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &basketspb.GetBasketResponse{
		Basket: s.basketFromDomain(basket),
	}, nil
}

func (s server) basketFromDomain(basket *domain.Basket) *basketspb.Basket {
	protoBasket := &basketspb.Basket{
		Id: basket.ID(),
	}

	protoBasket.Items = make([]*basketspb.Item, 0, len(basket.Items))

	for _, item := range basket.Items {
		protoBasket.Items = append(protoBasket.Items, &basketspb.Item{
			StoreId:      item.StoreID,
			StoreName:    item.StoreName,
			ProductId:    item.ProductID,
			ProductName:  item.ProductName,
			ProductPrice: item.ProductPrice,
			Quantity:     int32(item.Quantity),
		})
	}

	return protoBasket
}
