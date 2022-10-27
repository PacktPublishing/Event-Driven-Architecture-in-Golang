package application

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/orderingpb"
	"eda-in-golang/search/internal/models"
)

type OrderHandlers[T ddd.Event] struct {
	orders    OrderRepository
	customers CustomerRepository
	stores    StoreRepository
	products  ProductRepository
}

var _ ddd.EventHandler[ddd.Event] = (*OrderHandlers[ddd.Event])(nil)

func NewOrderHandlers(orders OrderRepository, customers CustomerRepository, stores StoreRepository, products ProductRepository) OrderHandlers[ddd.Event] {
	return OrderHandlers[ddd.Event]{
		orders:    orders,
		customers: customers,
		stores:    stores,
		products:  products,
	}
}

func (h OrderHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case orderingpb.OrderCreatedEvent:
		return h.onOrderCreated(ctx, event)
	case orderingpb.OrderReadiedEvent:
		return h.onOrderReadied(ctx, event)
	case orderingpb.OrderCanceledEvent:
		return h.onOrderCanceled(ctx, event)
	case orderingpb.OrderCompletedEvent:
		return h.onOrderCompleted(ctx, event)
	}
	return nil
}

func (h OrderHandlers[T]) onOrderCreated(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderCreated)

	customer, err := h.customers.Find(ctx, payload.CustomerId)
	if err != nil {
		return err
	}

	var total float64
	items := make([]models.Item, len(payload.GetItems()))
	seenStores := map[string]*models.Store{}
	for i, item := range payload.GetItems() {
		product, err := h.products.Find(ctx, item.GetProductId())
		if err != nil {
			return err
		}
		var store *models.Store
		var exists bool

		if store, exists = seenStores[product.StoreID]; !exists {
			store, err = h.stores.Find(ctx, product.StoreID)
			if err != nil {
				return err
			}
			seenStores[store.ID] = store
		}
		items[i] = models.Item{
			ProductID:   product.ID,
			StoreID:     store.ID,
			ProductName: product.Name,
			StoreName:   store.Name,
			Price:       item.Price,
			Quantity:    int(item.Quantity),
		}
		total += float64(item.Quantity) * item.Price
	}
	order := &models.Order{
		OrderID:      payload.GetId(),
		CustomerID:   customer.ID,
		CustomerName: customer.Name,
		Items:        items,
		Total:        total,
		Status:       "New",
	}
	return h.orders.Add(ctx, order)
}

func (h OrderHandlers[T]) onOrderReadied(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderReadied)
	return h.orders.UpdateStatus(ctx, payload.GetId(), "Ready For Pickup")
}

func (h OrderHandlers[T]) onOrderCanceled(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderCanceled)
	return h.orders.UpdateStatus(ctx, payload.GetId(), "Canceled")
}

func (h OrderHandlers[T]) onOrderCompleted(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderCompleted)
	return h.orders.UpdateStatus(ctx, payload.GetId(), "Completed")
}
