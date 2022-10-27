package e2e

import (
	"github.com/cucumber/godog"
	"github.com/go-openapi/strfmt"
	"github.com/stackus/errors"

	"eda-in-golang/stores/storesclient"
	"eda-in-golang/stores/storesclient/models"
	"eda-in-golang/stores/storesclient/store"
)

type storesContext struct {
	*suiteContext
	client        *storesclient.StoreManagement
	fetchedStores bool
}

var _ featureContext = (*storesContext)(nil)

func newStoresContext(sc *suiteContext) featureContext {
	return &storesContext{
		suiteContext: sc,
		client:       storesclient.New(sc.transport, strfmt.Default),
	}
}

func (c *storesContext) register(ctx *godog.ScenarioContext) {
	ctx.Step(`^a valid store owner$`, c.aValidStoreOwner)
	ctx.Step(`^I create the store called "([^"]*)"$`, c.iCreateTheStoreCalled)
	ctx.Step(`^(?:ensure |expect )?a store called "([^"]*)" (?:to )?exists?$`, c.expectAStoreCalledToExist)
	ctx.Step(`^(?:ensure |expect )?no store called "([^"]*)" (?:to )?exists?$`, c.expectNoStoreCalledToExist)
}

func (c *storesContext) reset() {
	c.stores = make(map[string]string)
	c.fetchedStores = false
	c.truncate("stores.stores")
	c.truncate("stores.products")
	c.truncate("stores.events")
	c.truncate("stores.snapshots")
	c.truncate("stores.inbox")
	c.truncate("stores.outbox")

}

func (c *storesContext) aValidStoreOwner() {
	// noop
}

func (c *storesContext) expectAStoreCalledToExist(name string) error {
	if !c.fetchedStores {
		err := c.fetchStores()
		if err != nil {
			return err
		}
	}

	if _, exists := c.stores[name]; !exists {
		return errors.ErrNotFound.Msgf("the store `%s` does not exist", name)
	}
	return nil
}

func (c *storesContext) expectNoStoreCalledToExist(name string) error {
	if !c.fetchedStores {
		err := c.fetchStores()
		if err != nil {
			return err
		}
	}

	if _, exists := c.stores[name]; exists {
		return errors.ErrNotFound.Msgf("the store `%s` does exist", name)
	}
	return nil
}

func (c *storesContext) iCreateTheStoreCalled(name string) {
	resp, err := c.client.Store.CreateStore(store.NewCreateStoreParams().WithBody(&models.StorespbCreateStoreRequest{
		Location: "anywhere",
		Name:     name,
	}))
	if err != nil {
		c.lastErr = err
		return
	}

	c.stores[name] = resp.Payload.ID
}

func (c *storesContext) fetchStores() error {
	resp, err := c.client.Store.GetStores(store.NewGetStoresParams())
	if err != nil {
		return err
	}

	for _, s := range resp.Payload.Stores {
		c.stores[s.Name] = s.ID
	}

	c.fetchedStores = true

	return nil
}
