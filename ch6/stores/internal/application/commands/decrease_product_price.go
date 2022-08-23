package commands

import (
	"context"

	"eda-in-golang/stores/internal/domain"
)

type DecreaseProductPrice struct {
	ID    string
	Price float64
}

type DecreaseProductPriceHandler struct {
	products domain.ProductRepository
}

func NewDecreaseProductPriceHandler(products domain.ProductRepository) DecreaseProductPriceHandler {
	return DecreaseProductPriceHandler{products: products}
}

func (h DecreaseProductPriceHandler) DecreaseProductPrice(ctx context.Context, cmd DecreaseProductPrice) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = product.DecreasePrice(cmd.Price); err != nil {
		return err
	}

	return h.products.Save(ctx, product)
}
