package commands

import (
	"context"

	"github.com/stackus/errors"

	"eda-in-golang/stores/internal/domain"
)

type AddProduct struct {
	ID          string
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
}

type AddProductHandler struct {
	products domain.ProductRepository
}

func NewAddProductHandler(products domain.ProductRepository) AddProductHandler {
	return AddProductHandler{
		products: products,
	}
}

func (h AddProductHandler) AddProduct(ctx context.Context, cmd AddProduct) error {
	product, err := domain.CreateProduct(cmd.ID, cmd.StoreID, cmd.Name, cmd.Description, cmd.SKU, cmd.Price)
	if err != nil {
		return errors.Wrap(err, "error adding product")
	}

	return errors.Wrap(h.products.Save(ctx, product), "error adding product")
}
