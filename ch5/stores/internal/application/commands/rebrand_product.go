package commands

import (
	"context"

	"eda-in-golang/stores/internal/domain"
)

type RebrandProduct struct {
	ID          string
	Name        string
	Description string
}

type RebrandProductHandler struct {
	products domain.ProductRepository
}

func NewRebrandProductHandler(products domain.ProductRepository) RebrandProductHandler {
	return RebrandProductHandler{
		products: products,
	}
}

func (h RebrandProductHandler) RebrandProduct(ctx context.Context, cmd RebrandProduct) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = product.Rebrand(cmd.Name, cmd.Description); err != nil {
		return err
	}

	return h.products.Save(ctx, product)
}
