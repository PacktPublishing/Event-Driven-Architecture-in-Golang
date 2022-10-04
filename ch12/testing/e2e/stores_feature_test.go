//go:build e2e

package e2e

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stackus/errors"

	"eda-in-golang/stores/storesclient"
	"eda-in-golang/stores/storesclient/models"
	"eda-in-golang/stores/storesclient/product"
	"eda-in-golang/stores/storesclient/store"
)

type storeIDKey struct{}
type productIDKey struct{}
type productMapKey struct{}

type storesFeature struct {
	db     *sql.DB
	client *storesclient.StoreManagement
}

var _ feature = (*storesFeature)(nil)

func (c *storesFeature) init(cfg featureConfig) (err error) {
	if cfg.useMonoDB {
		c.db, err = sql.Open("pgx", "postgres://mallbots_user:mallbots_pass@localhost:5432/mallbots?sslmode=disable")
	} else {
		c.db, err = sql.Open("pgx", "postgres://stores_user:stores_pass@localhost:5432/stores?sslmode=disable&search_path=stores,public")
	}
	if err != nil {
		return
	}
	c.client = storesclient.New(cfg.transport, strfmt.Default)

	return
}

func (c *storesFeature) register(ctx *godog.ScenarioContext) {
	ctx.Step(`^a valid store owner$`, c.noop)
	ctx.Step(`^the store (?:called )?"([^"]*)" already exists$`, c.iCreateTheStoreCalled)
	ctx.Step(`^I create (?:the|a) store called "([^"]*)"$`, c.iCreateTheStoreCalled)
	ctx.Step(`^(?:I )?(?:ensure |expect )?the store (?:was|is) created$`, c.expectTheStoreWasCreated)
	ctx.Step(`^(?:I )?(?:ensure |expect )?a store called "([^"]*)" (?:to )?exists?$`, c.expectAStoreCalledToExist)
	ctx.Step(`^(?:I )?(?:ensure |expect )?no store called "([^"]*)" (?:to )?exists?$`, c.expectNoStoreCalledToExist)

	ctx.Step(`^I create the product called "([^"]*)"$`, c.iCreateTheProductCalled)
	ctx.Step(`^I create the product called "([^"]*)" with price "([^"]*)"$`, c.iCreateTheProductCalledWithPrice)
	ctx.Step(`^(?:I )?(?:ensure |expect )?the product (?:was|is) created$`, c.expectTheProductWasCreated)
	ctx.Step(`^(?:I )?(?:ensure |expect )?a product called "([^"]*)" (?:to )?exists?$`, c.expectAProductCalledToExist)
	ctx.Step(`^(?:I )?(?:ensure |expect )?no product called "([^"]*)" (?:to )?exists?$`, c.expectNoProductCalledToExist)

	ctx.Step(`^a store has the following items$`, c.aStoreHasTheFollowingItems)
}

func (c *storesFeature) reset() {
	truncate := func(tableName string) {
		_, _ = c.db.Exec(fmt.Sprintf("TRUNCATE %s", tableName))
	}

	truncate("stores.stores")
	truncate("stores.products")
	truncate("stores.events")
	truncate("stores.snapshots")
	truncate("stores.inbox")
	truncate("stores.outbox")
}

func (c *storesFeature) noop() {
	// noop
}

func (c *storesFeature) expectAStoreCalledToExist(ctx context.Context, name string) error {
	var storeID string
	row := c.db.QueryRow("SELECT id FROM stores.stores WHERE name = $1", name)
	err := row.Scan(&storeID)
	if err != nil {
		return errors.ErrNotFound.Msgf("the store `%s` does not exist", name)
	}

	return nil
}

func (c *storesFeature) expectNoStoreCalledToExist(name string) error {
	var storeID string
	row := c.db.QueryRow("SELECT id FROM stores.stores WHERE name = $1", name)
	err := row.Scan(&storeID)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}

	return errors.ErrAlreadyExists.Msgf("the store `%s` does exist", name)
}

func (c *storesFeature) iCreateTheStoreCalled(ctx context.Context, name string) context.Context {
	resp, err := c.client.Store.CreateStore(store.NewCreateStoreParams().WithBody(&models.StorespbCreateStoreRequest{
		Location: "anywhere",
		Name:     name,
	}))
	ctx = setLastResponseAndError(ctx, resp, err)
	if err != nil {
		return ctx
	}
	return context.WithValue(ctx, storeIDKey{}, resp.Payload.ID)
}

func (c *storesFeature) expectTheStoreWasCreated(ctx context.Context) error {
	if err := lastResponseWas(ctx, &store.CreateStoreOK{}); err != nil {
		return err
	}

	return nil
}

func (c *storesFeature) expectAProductCalledToExist(ctx context.Context, name string) error {
	var productID string
	row := c.db.QueryRow("SELECT id FROM stores.products WHERE name = $1", name)
	err := row.Scan(&productID)
	if err != nil {
		return errors.ErrNotFound.Msgf("the product `%s` does not exist", name)
	}

	return nil
}

func (c *storesFeature) expectNoProductCalledToExist(name string) error {
	var productID string
	row := c.db.QueryRow("SELECT id FROM stores.products WHERE name = $1", name)
	err := row.Scan(&productID)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}

	return errors.ErrAlreadyExists.Msgf("the product `%s` does exist", name)
}

func (c *storesFeature) iCreateTheProductCalled(ctx context.Context, name string) (context.Context, error) {
	return c.iCreateTheProductCalledWithPrice(ctx, name, 9.99)
}

func (c *storesFeature) iCreateTheProductCalledWithPrice(ctx context.Context, name string, price float64) (context.Context, error) {
	storeID, err := lastStoreID(ctx)
	if err != nil {
		return ctx, err
	}
	resp, err := c.client.Product.AddProduct(product.NewAddProductParams().WithStoreID(storeID).WithBody(&models.AddProductParamsBody{
		Name:  name,
		Price: price,
	}))
	ctx = setLastResponseAndError(ctx, resp, err)
	if err != nil {
		return ctx, nil
	}
	ctx = addProduct(ctx, resp.Payload.ID, name)
	return context.WithValue(ctx, productIDKey{}, resp.Payload.ID), nil
}

func (c *storesFeature) expectTheProductWasCreated(ctx context.Context) error {
	if err := lastResponseWas(ctx, &product.AddProductOK{}); err != nil {
		return err
	}

	return nil
}

func (c *storesFeature) aStoreHasTheFollowingItems(ctx context.Context, table *godog.Table) (context.Context, error) {
	type Item struct {
		Name  string
		Price float64
	}
	ctx = c.iCreateTheStoreCalled(ctx, "AnyStore")

	if err := lastError(ctx); err != nil {
		return ctx, err
	}
	items, err := assist.CreateSlice(new(Item), table)
	if err != nil {
		return ctx, err
	}

	for _, i := range items.([]*Item) {
		ctx, err = c.iCreateTheProductCalledWithPrice(ctx, i.Name, i.Price)
		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}

func lastStoreID(ctx context.Context) (string, error) {
	v := ctx.Value(storeIDKey{})
	if v == nil {
		return "", errors.ErrNotFound.Msg("no store ID to work with")
	}
	return v.(string), nil
}

func lastProductID(ctx context.Context) (string, error) {
	v := ctx.Value(productIDKey{})
	if v == nil {
		return "", errors.ErrNotFound.Msg("no product ID to work with")
	}
	return v.(string), nil
}

func addProduct(ctx context.Context, id, name string) context.Context {
	var products map[string]string
	v := ctx.Value(productMapKey{})
	if v == nil {
		products = make(map[string]string)
		ctx = context.WithValue(ctx, productMapKey{}, products)
	} else {
		products = v.(map[string]string)
	}

	products[name] = id

	return ctx
}

func getProductID(ctx context.Context, name string) string {
	v := ctx.Value(productMapKey{})
	if v == nil {
		return uuid.NewString()
	}
	products := v.(map[string]string)

	id, exists := products[name]
	if !exists {
		return uuid.NewString()
	}

	return id
}
