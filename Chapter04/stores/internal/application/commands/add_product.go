package commands

import (
	"context"

	"github.com/stackus/errors"

	"eda-in-golang/internal/ddd"
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
	stores          domain.StoreRepository
	products        domain.ProductRepository
	domainPublisher ddd.EventPublisher
}

func NewAddProductHandler(stores domain.StoreRepository, products domain.ProductRepository, domainPublisher ddd.EventPublisher) AddProductHandler {
	return AddProductHandler{
		stores:          stores,
		products:        products,
		domainPublisher: domainPublisher,
	}
}

func (h AddProductHandler) AddProduct(ctx context.Context, cmd AddProduct) error {
	if _, err := h.stores.Find(ctx, cmd.StoreID); err != nil {
		return errors.Wrap(err, "error adding product")
	}

	product, err := domain.CreateProduct(cmd.ID, cmd.StoreID, cmd.Name, cmd.Description, cmd.SKU, cmd.Price)
	if err != nil {
		return errors.Wrap(err, "error adding product")
	}

	if err = h.products.Save(ctx, product); err != nil {
		return errors.Wrap(err, "error adding product")
	}

	// publish domain events
	if err = h.domainPublisher.Publish(ctx, product.GetEvents()...); err != nil {
		return err
	}

	return nil
}
