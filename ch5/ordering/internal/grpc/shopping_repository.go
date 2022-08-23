package grpc

import (
	"context"

	"google.golang.org/grpc"

	"eda-in-golang/depot/depotpb"
	"eda-in-golang/ordering/internal/domain"
)

type ShoppingRepository struct {
	client depotpb.DepotServiceClient
}

var _ domain.ShoppingRepository = (*ShoppingRepository)(nil)

func NewShoppingListRepository(conn *grpc.ClientConn) ShoppingRepository {
	return ShoppingRepository{client: depotpb.NewDepotServiceClient(conn)}
}

func (r ShoppingRepository) Create(ctx context.Context, orderID string, orderItems []domain.Item) (string, error) {
	items := make([]*depotpb.OrderItem, len(orderItems))
	for i, item := range orderItems {
		items[i] = r.itemFromDomain(item)
	}

	response, err := r.client.CreateShoppingList(ctx, &depotpb.CreateShoppingListRequest{
		OrderId: orderID,
		Items:   items,
	})
	if err != nil {
		return "", err
	}

	return response.GetId(), nil
}

func (r ShoppingRepository) Cancel(ctx context.Context, shoppingID string) error {
	_, err := r.client.CancelShoppingList(ctx, &depotpb.CancelShoppingListRequest{Id: shoppingID})
	return err
}

func (r ShoppingRepository) itemFromDomain(item domain.Item) *depotpb.OrderItem {
	return &depotpb.OrderItem{
		ProductId: item.ProductID,
		StoreId:   item.StoreID,
		Quantity:  int32(item.Quantity),
	}
}
