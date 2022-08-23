package commands

import (
	"context"

	"eda-in-golang/stores/internal/domain"
)

type IncreaseProductPrice struct {
	ID    string
	Price float64
}

type IncreaseProductPriceHandler struct {
	products domain.ProductRepository
}

func NewIncreaseProductPriceHandler(products domain.ProductRepository) IncreaseProductPriceHandler {
	return IncreaseProductPriceHandler{products: products}
}

func (h IncreaseProductPriceHandler) IncreaseProductPrice(ctx context.Context, cmd IncreaseProductPrice) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = product.IncreasePrice(cmd.Price); err != nil {
		return err
	}

	return h.products.Save(ctx, product)
}
