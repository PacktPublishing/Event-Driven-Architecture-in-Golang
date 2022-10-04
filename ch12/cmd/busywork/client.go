package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/ddosify/go-faker/faker"
	"github.com/go-openapi/runtime/client"

	"eda-in-golang/cmd/busywork/baskets"
	"eda-in-golang/cmd/busywork/customers"
	"eda-in-golang/cmd/busywork/payments"
	"eda-in-golang/cmd/busywork/stores"
)

var f = faker.NewFaker()

type busyworkClient struct {
	id        string
	log       *log.Logger
	baskets   baskets.Client
	customers customers.Client
	payments  payments.Client
	stores    stores.Client
	interval  time.Duration
}

func newBusyworkClient(id string, interval time.Duration) *busyworkClient {
	transport := client.New(*hostAddr, "/", nil)
	// transport := client.NewWithClient(*hostAddr, "/", nil, &http.Client{
	// 	Transport: otelhttp.NewTransport(http.DefaultTransport),
	// })

	return &busyworkClient{
		id:        id,
		log:       log.New(os.Stdout, fmt.Sprintf("[%s] ", id), log.Lmicroseconds|log.Lmsgprefix),
		baskets:   baskets.NewClient(transport),
		customers: customers.NewClient(transport),
		payments:  payments.NewClient(transport),
		stores:    stores.NewClient(transport),
		interval:  interval,
	}
}

func (c *busyworkClient) run(ctx context.Context) error {
	timer := time.NewTimer(0)
	for {
		select {
		case <-timer.C:
			// time to get busy
			err := c.work(context.Background())
			if err != nil {
				c.log.Println(fmt.Sprintf(`is having a hard time. Error: %s`, err.Error()))
			}
			timer.Reset(c.interval)
		case <-ctx.Done():
			c.log.Println("Quitting time")
			return nil
		default:
		}
	}
}

func (c *busyworkClient) work(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	first := rand.Intn(2)
	switch first {
	case 0:
		second := rand.Intn(4)
		switch second {
		case 0, 1, 2:
			c.log.Println("is considering browsing for new things")
			return c.justBrowsing(ctx)
		case 3:
			c.log.Println("is considering buying some things")
			return c.buyingItems(ctx)
		}
	case 1:
		second := rand.Intn(10)
		switch second {
		case 0, 5, 6, 7:
			c.log.Println("is considering registering a new account")
			return c.registerCustomer(ctx)
		case 1:
			c.log.Println("is considering setting up a new store")
			return c.setupAStore(ctx)
		case 2, 8, 9:
			c.log.Println("is considering adding new inventory")
			return c.addNewInventory(ctx)
		case 3:
			c.log.Println("is considering rebranding a store")
			return c.rebrandStore(ctx)
		case 4:
			c.log.Println("is considering updating product branding")
			return c.rebrandProduct(ctx)
		}
	}

	return nil
}

func (c *busyworkClient) pause() {
	jitter := time.Duration(rand.Int63n(int64(1500 * time.Millisecond)))
	time.Sleep(500*time.Millisecond + jitter)
}

func (c *busyworkClient) registerCustomer(ctx context.Context) error {
	_, err := c.generateCustomer(ctx)
	return err
}

func (c *busyworkClient) setupAStore(ctx context.Context) error {
	storeID, err := c.generateStore(ctx)
	if err != nil {
		return err
	}

	_, err = c.generateProducts(ctx, storeID, 2, 5)
	if err != nil {
		return err
	}

	name, err := c.stores.GetStoreName(ctx, storeID)
	if err != nil {
		return err
	}

	c.log.Println(fmt.Sprintf(`has finished setting up "%s"`, name))
	return nil
}

func (c *busyworkClient) addNewInventory(ctx context.Context) error {
	storeIDs, err := c.stores.GetStores(ctx)
	if err != nil {
		return err
	}

	if len(storeIDs) == 0 {
		c.log.Println("has no store to work with")
		return nil
	}

	storeID := storeIDs[rand.Intn(len(storeIDs))]

	_, err = c.generateProducts(ctx, storeID, 1, 3)
	if err != nil {
		return err
	}

	c.log.Println("is done adding new inventory")

	return nil
}

func (c *busyworkClient) rebrandStore(ctx context.Context) error {
	storeIDs, err := c.getStores(ctx)
	if err != nil {
		return err
	}

	if len(storeIDs) == 0 {
		c.log.Println("but has no store to work with")
		return nil
	}

	storeID := storeIDs[rand.Intn(len(storeIDs))]

	name, err := c.stores.GetStoreName(ctx, storeID)
	if err != nil {
		return err
	}

	newName := f.RandomCompanyName()

	err = c.stores.RebrandStore(ctx, storeID, newName)
	if err != nil {
		return nil
	}
	c.log.Println(fmt.Sprintf(`has rebranded the store "%s" to "%s"`, name, newName))
	return nil
}

func (c *busyworkClient) rebrandProduct(ctx context.Context) error {
	storeIDs, err := c.getStores(ctx)
	if err != nil {
		return err
	}

	if len(storeIDs) == 0 {
		c.log.Println("but has no store to work with")
		return nil
	}

	storeID := storeIDs[rand.Intn(len(storeIDs))]

	productIDs, err := c.getCatalog(ctx, storeID)
	if err != nil {
		return err
	}

	if len(productIDs) == 0 {
		c.log.Println("but has no product to work with")
		return nil
	}

	productID := productIDs[rand.Intn(len(productIDs))]

	name, _, err := c.stores.GetProductDetails(ctx, productID)
	if err != nil {
		return err
	}

	newName := f.RandomProductName()
	newDesc := f.RandomBs()

	err = c.stores.RebrandProduct(ctx, productID, newName, newDesc)
	if err != nil {
		return nil
	}
	c.log.Println(fmt.Sprintf(`has rebranded the product "%s" to "%s"`, name, newName))
	return nil
}

func (c *busyworkClient) justBrowsing(ctx context.Context) error {
	customerID, err := c.generateCustomer(ctx)
	if err != nil {
		return err
	}

	basketID, err := c.generateBasket(ctx, customerID)
	if err != nil {
		return err
	}

	storeIDs, err := c.stores.GetStores(ctx)
	if err != nil {
		return err
	}

	if len(storeIDs) == 0 {
		c.log.Println("but has no stores to shop at")
		return nil
	}

	storeCount := 1 + rand.Intn(3)

	if storeCount > len(storeIDs) {
		c.log.Println("but not enough stores are available")
		return nil
	}

	total := 0.0
	for i := 0; i < storeCount; i++ {
		// do not care about repeats
		storeID := storeIDs[rand.Intn(len(storeIDs))]
		name, err := c.stores.GetStoreName(ctx, storeID)
		if err != nil {
			return err
		}
		c.log.Println(fmt.Sprintf(`is browsing the items from "%s"`, name))
		c.pause()

		productIDs, err := c.stores.GetCatalog(ctx, storeID)
		if err != nil {
			return err
		}

		productID := productIDs[rand.Intn(len(productIDs))]

		name, price, err := c.stores.GetProductDetails(ctx, productID)
		quantity := 1 + rand.Intn(4)
		c.log.Println(fmt.Sprintf(`might buy %d "%s" for $%.2f each`, quantity, name, price))

		total += price * float64(quantity)

		err = c.addItem(ctx, basketID, productID, quantity)
		if err != nil {
			return err
		}
	}
	c.log.Println(fmt.Sprintf(`thinks $%.2f is too much`, total))
	return c.baskets.CancelBasket(ctx, basketID)

}

func (c *busyworkClient) buyingItems(ctx context.Context) error {
	customerID, err := c.generateCustomer(ctx)
	if err != nil {
		return err
	}

	basketID, err := c.generateBasket(ctx, customerID)
	if err != nil {
		return err
	}

	storeIDs, err := c.stores.GetStores(ctx)
	if err != nil {
		return err
	}

	if len(storeIDs) == 0 {
		c.log.Println("but has no stores to shop at")
		return nil
	}

	storeCount := 1 + rand.Intn(3)

	if storeCount > len(storeIDs) {
		c.log.Println("but not enough store are available")
		return nil
	}

	total := 0.0
	for i := 0; i < storeCount; i++ {
		// do not care about repeats
		storeID := storeIDs[rand.Intn(len(storeIDs))]
		name, err := c.stores.GetStoreName(ctx, storeID)
		if err != nil {
			return err
		}
		c.log.Println(fmt.Sprintf(`is browsing the items from "%s"`, name))
		c.pause()

		productIDs, err := c.stores.GetCatalog(ctx, storeID)
		if err != nil {
			return err
		}

		productID := productIDs[rand.Intn(len(productIDs))]

		name, price, err := c.stores.GetProductDetails(ctx, productID)
		quantity := 1 + rand.Intn(4)
		c.log.Println(fmt.Sprintf(`might buy %d "%s" for $%.2f each`, quantity, name, price))

		total += price * float64(quantity)

		err = c.addItem(ctx, basketID, productID, quantity)
		if err != nil {
			return err
		}
	}
	c.log.Println(fmt.Sprintf(`is OK with $%.2f`, total))

	paymentID, err := c.generatePayment(ctx, customerID, total)
	if err != nil {
		return err
	}

	return c.baskets.CheckoutBasket(ctx, basketID, paymentID)
}

// ---

func (c *busyworkClient) generateCustomer(ctx context.Context) (string, error) {
	c.pause()

	return c.customers.RegisterCustomer(ctx, f.RandomUsername(), f.RandomPhoneNumber())
}

func (c *busyworkClient) generateBasket(ctx context.Context, customerID string) (string, error) {
	c.pause()
	return c.baskets.StartBasket(ctx, customerID)
}

func (c *busyworkClient) generatePayment(ctx context.Context, customerID string, total float64) (string, error) {
	c.pause()
	return c.payments.AuthorizePayment(ctx, customerID, total)
}

func (c *busyworkClient) addItem(ctx context.Context, basketID, productID string, quantity int) error {
	c.pause()
	return c.baskets.AddItem(ctx, basketID, productID, quantity)
}

func (c *busyworkClient) getStores(ctx context.Context) ([]string, error) {
	c.pause()
	return c.stores.GetStores(ctx)
}

func (c *busyworkClient) getCatalog(ctx context.Context, storeID string) ([]string, error) {
	c.pause()
	return c.stores.GetCatalog(ctx, storeID)
}

func (c *busyworkClient) generateStore(ctx context.Context) (string, error) {
	c.pause()

	name := f.RandomCompanyName()
	department := f.RandomDepartment()
	c.log.Println(fmt.Sprintf(`is getting "%s" ready`, name))
	storeID, err := c.stores.CreateStore(ctx, name, department)
	if err != nil {
		return "", err
	}

	return storeID, err
}

func (c *busyworkClient) generateProducts(ctx context.Context, storeID string, min, max int) ([]string, error) {
	count := f.RandomIntBetween(min, max)
	productIDs := make([]string, count)
	for i := 0; i < count; i++ {
		c.pause()

		name := f.RandomProductName()
		price := float64(500+rand.Intn(700)) / 100
		c.log.Println(fmt.Sprintf(`is adding "%s" for $%.2f`, name, price))
		productID, err := c.stores.AddProduct(ctx, storeID, name, f.RandomBs(), f.RandomProductAdjective(), price)
		if err != nil {
			return nil, err
		}
		productIDs[i] = productID
	}
	return productIDs, nil
}
