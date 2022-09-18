package handlers

import (
	"context"

	"eda-in-golang/customers/customerspb"
	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/orderingpb"
	"eda-in-golang/search/internal/application"
	"eda-in-golang/search/internal/models"
	"eda-in-golang/stores/storespb"
)

type integrationHandlers[T ddd.Event] struct {
	orders    application.OrderRepository
	customers application.CustomerCacheRepository
	products  application.ProductCacheRepository
	stores    application.StoreCacheRepository
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(orders application.OrderRepository, customers application.CustomerCacheRepository,
	stores application.StoreCacheRepository, products application.ProductCacheRepository) ddd.EventHandler[ddd.Event] {
	return integrationHandlers[ddd.Event]{
		orders:    orders,
		customers: customers,
		products:  products,
		stores:    stores,
	}
}

func RegisterIntegrationEventHandlers(subscriber am.EventSubscriber, handlers ddd.EventHandler[ddd.Event]) (err error) {
	evtMsgHandler := am.MessageHandlerFunc[am.IncomingEventMessage](func(ctx context.Context, eventMsg am.IncomingEventMessage) error {
		return handlers.HandleEvent(ctx, eventMsg)
	})

	if _, err = subscriber.Subscribe(customerspb.CustomerAggregateChannel, evtMsgHandler, am.MessageFilter{
		customerspb.CustomerRegisteredEvent,
	}, am.GroupName("search-customers")); err != nil {
		return
	}

	if _, err = subscriber.Subscribe(orderingpb.OrderAggregateChannel, evtMsgHandler, am.MessageFilter{
		orderingpb.OrderCreatedEvent,
		orderingpb.OrderReadiedEvent,
		orderingpb.OrderCanceledEvent,
		orderingpb.OrderCompletedEvent,
	}, am.GroupName("notification-orders")); err != nil {
		return
	}

	if _, err = subscriber.Subscribe(storespb.ProductAggregateChannel, evtMsgHandler, am.MessageFilter{
		storespb.ProductAddedEvent,
		storespb.ProductRebrandedEvent,
		storespb.ProductRemovedEvent,
	}, am.GroupName("search-products")); err != nil {
		return
	}

	if _, err = subscriber.Subscribe(storespb.StoreAggregateChannel, evtMsgHandler, am.MessageFilter{
		storespb.StoreCreatedEvent,
		storespb.StoreRebrandedEvent,
	}, am.GroupName("search-stores")); err != nil {
		return
	}

	return
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case customerspb.CustomerRegisteredEvent:
		return h.onCustomerRegistered(ctx, event)
	case storespb.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case storespb.ProductRebrandedEvent:
		return h.onProductRebranded(ctx, event)
	case storespb.ProductRemovedEvent:
		return h.onProductRemoved(ctx, event)
	case storespb.StoreCreatedEvent:
		return h.onStoreCreated(ctx, event)
	case storespb.StoreRebrandedEvent:
		return h.onStoreRebranded(ctx, event)
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

func (h integrationHandlers[T]) onCustomerRegistered(ctx context.Context, event T) error {
	payload := event.Payload().(*customerspb.CustomerRegistered)
	return h.customers.Add(ctx, payload.GetId(), payload.GetName())
}

func (h integrationHandlers[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.ProductAdded)
	return h.products.Add(ctx, payload.GetId(), payload.GetStoreId(), payload.GetName())
}

func (h integrationHandlers[T]) onProductRebranded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.ProductRebranded)
	return h.products.Rebrand(ctx, payload.GetId(), payload.GetName())
}

func (h integrationHandlers[T]) onProductRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.ProductRemoved)
	return h.products.Remove(ctx, payload.GetId())
}

func (h integrationHandlers[T]) onStoreCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.StoreCreated)
	return h.stores.Add(ctx, payload.GetId(), payload.GetName())
}

func (h integrationHandlers[T]) onStoreRebranded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.StoreRebranded)
	return h.stores.Rename(ctx, payload.GetId(), payload.GetName())
}
func (h integrationHandlers[T]) onOrderCreated(ctx context.Context, event T) error {
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

func (h integrationHandlers[T]) onOrderReadied(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderReadied)
	return h.orders.UpdateStatus(ctx, payload.GetId(), "Ready For Pickup")
}

func (h integrationHandlers[T]) onOrderCanceled(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderCanceled)
	return h.orders.UpdateStatus(ctx, payload.GetId(), "Canceled")
}

func (h integrationHandlers[T]) onOrderCompleted(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderCompleted)
	return h.orders.UpdateStatus(ctx, payload.GetId(), "Completed")
}
