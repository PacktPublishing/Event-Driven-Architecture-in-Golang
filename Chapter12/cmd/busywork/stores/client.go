package stores

import (
	"context"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"eda-in-golang/stores/storesclient"
	"eda-in-golang/stores/storesclient/models"
	"eda-in-golang/stores/storesclient/product"
	"eda-in-golang/stores/storesclient/store"
)

type Client interface {
	CreateStore(ctx context.Context, name, location string) (string, error)
	GetStores(ctx context.Context) ([]string, error)
	GetStoreName(ctx context.Context, storeID string) (string, error)
	RebrandStore(ctx context.Context, storeID, name string) error

	AddProduct(ctx context.Context, storeID string, name, description, sku string, price float64) (string, error)
	RebrandProduct(ctx context.Context, productID, name, description string) error
	GetProductDetails(ctx context.Context, productID string) (string, float64, error)
	GetCatalog(ctx context.Context, storeID string) ([]string, error)
}

type client struct {
	c *storesclient.StoreManagement
}

func NewClient(transport runtime.ClientTransport) Client {
	return &client{c: storesclient.New(transport, strfmt.Default)}
}

func (c *client) CreateStore(ctx context.Context, name, location string) (string, error) {
	resp, err := c.c.Store.CreateStore(&store.CreateStoreParams{
		Body: &models.StorespbCreateStoreRequest{
			Location: location,
			Name:     name,
		},
		Context: ctx,
	})
	if err != nil {
		return "", err
	}

	return resp.GetPayload().ID, nil
}

func (c *client) GetStores(ctx context.Context) ([]string, error) {
	resp, err := c.c.Store.GetStores(&store.GetStoresParams{
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}

	storeIDs := make([]string, len(resp.GetPayload().Stores))
	for i, s := range resp.GetPayload().Stores {
		storeIDs[i] = s.ID
	}

	return storeIDs, nil
}

func (c *client) RebrandStore(ctx context.Context, storeID, name string) error {
	_, err := c.c.Store.RebrandStore(&store.RebrandStoreParams{
		Body:    &models.RebrandStoreParamsBody{Name: name},
		ID:      storeID,
		Context: ctx,
	})
	return err
}

func (c *client) GetStoreName(ctx context.Context, storeID string) (string, error) {
	resp, err := c.c.Store.GetStore(&store.GetStoreParams{
		ID:      storeID,
		Context: ctx,
	})
	if err != nil {
		return "", err
	}

	return resp.GetPayload().Store.Name, nil
}

func (c *client) AddProduct(ctx context.Context, storeID string, name, description, sku string, price float64) (string, error) {
	resp, err := c.c.Product.AddProduct(&product.AddProductParams{
		Body: &models.AddProductParamsBody{
			Description: description,
			Name:        name,
			Price:       price,
			Sku:         sku,
		},
		StoreID: storeID,
		Context: ctx,
	})
	if err != nil {
		return "", err
	}

	return resp.GetPayload().ID, nil
}

func (c *client) GetCatalog(ctx context.Context, storeID string) ([]string, error) {
	resp, err := c.c.Product.GetStoreProducts(&product.GetStoreProductsParams{
		StoreID: storeID,
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}

	productIDs := make([]string, len(resp.GetPayload().Products))
	for i, s := range resp.GetPayload().Products {
		productIDs[i] = s.ID
	}

	return productIDs, nil
}

func (c *client) RebrandProduct(ctx context.Context, productID, name, description string) error {
	_, err := c.c.Product.RebrandProduct(&product.RebrandProductParams{
		Body: &models.RebrandProductParamsBody{
			Name:        name,
			Description: description,
		},
		ID:      productID,
		Context: ctx,
	})
	return err
}

func (c *client) GetProductDetails(ctx context.Context, productID string) (string, float64, error) {
	resp, err := c.c.Product.GetProduct(&product.GetProductParams{
		ID:      productID,
		Context: ctx,
	})
	if err != nil {
		return "", 0, err
	}

	return resp.GetPayload().Product.Name, resp.GetPayload().Product.Price, nil
}
