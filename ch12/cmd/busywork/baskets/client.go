package baskets

import (
	"context"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"eda-in-golang/baskets/basketsclient"
	"eda-in-golang/baskets/basketsclient/basket"
	"eda-in-golang/baskets/basketsclient/item"
	"eda-in-golang/baskets/basketsclient/models"
)

type Client interface {
	StartBasket(ctx context.Context, customerID string) (string, error)
	CheckoutBasket(ctx context.Context, basketID, paymentID string) error
	CancelBasket(ctx context.Context, basketID string) error

	AddItem(ctx context.Context, basketID, productID string, quantity int) error
}

type client struct {
	c *basketsclient.ShoppingBaskets
}

func NewClient(transport runtime.ClientTransport) Client {
	return &client{
		c: basketsclient.New(transport, strfmt.Default),
	}
}

func (c *client) StartBasket(ctx context.Context, customerID string) (string, error) {
	resp, err := c.c.Basket.StartBasket(&basket.StartBasketParams{
		Body:    &models.BasketspbStartBasketRequest{CustomerID: customerID},
		Context: ctx,
	})
	if err != nil {
		return "", err
	}

	return resp.GetPayload().ID, nil
}

func (c *client) CheckoutBasket(ctx context.Context, basketID, paymentID string) error {
	_, err := c.c.Basket.CheckoutBasket(&basket.CheckoutBasketParams{
		Body:    &models.CheckoutBasketParamsBody{PaymentID: paymentID},
		ID:      basketID,
		Context: ctx,
	})
	return err
}

func (c *client) CancelBasket(ctx context.Context, basketID string) error {
	_, err := c.c.Basket.CancelBasket(&basket.CancelBasketParams{
		ID:      basketID,
		Context: ctx,
	})
	return err
}

func (c *client) AddItem(ctx context.Context, basketID, productID string, quantity int) error {
	_, err := c.c.Item.AddItem(&item.AddItemParams{
		Body: &models.AddItemParamsBody{
			ProductID: productID,
			Quantity:  int32(quantity),
		},
		ID:      basketID,
		Context: ctx,
	})
	return err
}
