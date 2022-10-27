package commands

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/internal/domain"
)

type DecreaseProductPrice struct {
	ID    string
	Price float64
}

type DecreaseProductPriceHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewDecreaseProductPriceHandler(products domain.ProductRepository, publisher ddd.EventPublisher[ddd.Event]) DecreaseProductPriceHandler {
	return DecreaseProductPriceHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h DecreaseProductPriceHandler) DecreaseProductPrice(ctx context.Context, cmd DecreaseProductPrice) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := product.DecreasePrice(cmd.Price)
	if err != nil {
		return err
	}

	err = h.products.Save(ctx, product)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
