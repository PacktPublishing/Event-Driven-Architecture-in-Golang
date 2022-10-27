//go:build e2e

package e2e

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/go-openapi/strfmt"
	"github.com/stackus/errors"

	"eda-in-golang/customers/customersclient"
	"eda-in-golang/customers/customersclient/customer"
	"eda-in-golang/customers/customersclient/models"
)

type customerIDKey = struct{}

type customersFeature struct {
	client *customersclient.Customers
	db     *sql.DB
}

func (c *customersFeature) init(cfg featureConfig) (err error) {
	if cfg.useMonoDB {
		c.db, err = sql.Open("pgx", "postgres://mallbots_user:mallbots_pass@localhost:5432/mallbots?sslmode=disable")
	} else {
		c.db, err = sql.Open("pgx", "postgres://customers_user:customers_pass@localhost:5432/customers?sslmode=disable&search_path=customers,public")
	}
	if err != nil {
		return
	}
	c.client = customersclient.New(cfg.transport, strfmt.Default)

	return
}

func (c *customersFeature) register(ctx *godog.ScenarioContext) {
	ctx.Step(`^I am a registered customer$`, c.iAmARegisteredCustomer)
	ctx.Step(`^I register a new customer as "([^"]*)"$`, c.iRegisterANewCustomerAs)
	ctx.Step(`^(?:I )?(?:ensure |expect )?the customer (?:was|is) created$`, c.expectTheCustomerWasCreated)
	ctx.Step(`^(?:I )?(?:ensure |expect )?a customer named "([^"]*)" (?:to )?exists?$`, c.expectACustomerNamedToExist)
	ctx.Step(`^(?:I )?(?:ensure |expect )?no customer named "([^"]*)" (?:to )?exists?$`, c.expectNoCustomerNamedToExist)

}

func (c *customersFeature) reset() {
	truncate := func(tableName string) {
		_, _ = c.db.Exec(fmt.Sprintf("TRUNCATE %s", tableName))
	}

	truncate("customers.customers")
	truncate("customers.inbox")
	truncate("customers.outbox")
}

func (c *customersFeature) iAmARegisteredCustomer(ctx context.Context) context.Context {
	resp, err := c.client.Customer.CreateCustomer(customer.NewCreateCustomerParams().WithBody(&models.CustomerspbRegisterCustomerRequest{
		Name:      "RegisteredCustomer",
		SmsNumber: "555-555-1212",
	}))
	ctx = setLastResponseAndError(ctx, resp, err)
	if err != nil {
		return ctx
	}
	return context.WithValue(ctx, customerIDKey{}, resp.Payload.ID)
}

func (c *customersFeature) iRegisterANewCustomerAs(ctx context.Context, name string) context.Context {
	resp, err := c.client.Customer.CreateCustomer(customer.NewCreateCustomerParams().WithBody(&models.CustomerspbRegisterCustomerRequest{
		Name:      name,
		SmsNumber: "555-555-1212",
	}))
	ctx = setLastResponseAndError(ctx, resp, err)
	if err != nil {
		return ctx
	}
	return context.WithValue(ctx, customerIDKey{}, resp.Payload.ID)
}

func (c *customersFeature) expectACustomerNamedToExist(name string) error {
	var customerID string
	row := c.db.QueryRow("SELECT id FROM customers.customers WHERE name = $1", name)
	err := row.Scan(&customerID)
	if err != nil {
		return errors.ErrNotFound.Msgf("the customer `%s` does not exist", name)
	}

	return nil
}

func (c *customersFeature) expectNoCustomerNamedToExist(name string) error {
	var customerID string
	row := c.db.QueryRow("SELECT id FROM customers.customers WHERE name = $1", name)
	err := row.Scan(&customerID)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}

	return errors.ErrAlreadyExists.Msgf("the customer `%s` does exist", name)
}

func (c *customersFeature) expectTheCustomerWasCreated(ctx context.Context) error {
	if err := lastResponseWas(ctx, &customer.CreateCustomerOK{}); err != nil {
		return err
	}

	return nil
}

func lastCustomerID(ctx context.Context) (string, error) {
	v := ctx.Value(customerIDKey{})
	if v == nil {
		return "", errors.ErrNotFound.Msg("no customer ID to work with")
	}
	return v.(string), nil
}
