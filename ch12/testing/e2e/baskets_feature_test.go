//go:build e2e

package e2e

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/go-openapi/strfmt"
	"github.com/stackus/errors"

	"eda-in-golang/baskets/basketsclient"
	"eda-in-golang/baskets/basketsclient/basket"
	"eda-in-golang/baskets/basketsclient/item"
	"eda-in-golang/baskets/basketsclient/models"
)

type basketIDKey = struct{}

type basketsFeature struct {
	client *basketsclient.ShoppingBaskets
	db     *sql.DB
}

func (c *basketsFeature) init(cfg featureConfig) (err error) {
	if cfg.useMonoDB {
		c.db, err = sql.Open("pgx", "postgres://mallbots_user:mallbots_pass@localhost:5432/mallbots?sslmode=disable")
	} else {
		c.db, err = sql.Open("pgx", "postgres://baskets_user:baskets_pass@localhost:5432/baskets?sslmode=disable&search_path=baskets,public")
	}
	if err != nil {
		return
	}
	c.client = basketsclient.New(cfg.transport, strfmt.Default)

	return
}

func (c *basketsFeature) register(ctx *godog.ScenarioContext) {
	ctx.Step(`^I start a new basket$`, c.iStartANewBasket)
	ctx.Step(`^(?:I )?(?:ensure |expect )?the basket (?:was|is) started$`, c.expectTheBasketWasStarted)

	ctx.Step(`^I add the items$`, c.iAddTheItems)
	ctx.Step(`^(?:I )?(?:ensure |expect )?the items (?:were|are) added$`, c.expectTheItemsWereAdded)
}

func (c *basketsFeature) reset() {
	truncate := func(tableName string) {
		_, _ = c.db.Exec(fmt.Sprintf("TRUNCATE %s", tableName))
	}

	truncate("baskets.events")
	truncate("baskets.snapshots")
	truncate("baskets.inbox")
	truncate("baskets.outbox")
	truncate("baskets.products_cache")
	truncate("baskets.stores_cache")
}

func (c *basketsFeature) iStartANewBasket(ctx context.Context) (context.Context, error) {
	customerID, err := lastCustomerID(ctx)
	if err != nil {
		return ctx, err
	}
	resp, err := c.client.Basket.StartBasket(basket.NewStartBasketParams().WithBody(&models.BasketspbStartBasketRequest{
		CustomerID: customerID,
	}))

	ctx = setLastResponseAndError(ctx, resp, err)
	if err != nil {
		return ctx, nil
	}
	return context.WithValue(ctx, basketIDKey{}, resp.Payload.ID), nil
}

func (c *basketsFeature) expectTheBasketWasStarted(ctx context.Context) error {
	if err := lastResponseWas(ctx, &basket.StartBasketOK{}); err != nil {
		return err
	}

	return nil
}

func (c *basketsFeature) iAddTheItems(ctx context.Context, table *godog.Table) (context.Context, error) {
	type Item struct {
		Name     string
		Quantity int
	}

	basketID, err := lastBasketID(ctx)
	if err != nil {
		return ctx, err
	}

	items, err := assist.CreateSlice(new(Item), table)
	if err != nil {
		return ctx, err
	}
	for _, i := range items.([]*Item) {
		productID := getProductID(ctx, i.Name)
		resp, err := c.client.Item.AddItem(item.NewAddItemParams().WithID(basketID).WithBody(&models.AddItemParamsBody{
			ProductID: productID,
			Quantity:  int32(i.Quantity),
		}))
		ctx = setLastResponseAndError(ctx, resp, err)
		if err != nil {
			break
		}
	}

	return ctx, nil
}

func (c *basketsFeature) expectTheItemsWereAdded(ctx context.Context) error {
	if err := lastResponseWas(ctx, &item.AddItemOK{}); err != nil {
		return err
	}

	return nil
}

func lastBasketID(ctx context.Context) (string, error) {
	v := ctx.Value(basketIDKey{})
	if v == nil {
		return "", errors.ErrNotFound.Msg("no basket ID to work with")
	}
	return v.(string), nil
}
