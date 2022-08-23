package application

import (
	"context"

	"eda-in-golang/customers/internal/domain"
	"eda-in-golang/internal/ddd"
)

type (
	RegisterCustomer struct {
		ID        string
		Name      string
		SmsNumber string
	}

	AuthorizeCustomer struct {
		ID string
	}

	GetCustomer struct {
		ID string
	}

	EnableCustomer struct {
		ID string
	}

	DisableCustomer struct {
		ID string
	}

	App interface {
		RegisterCustomer(ctx context.Context, register RegisterCustomer) error
		AuthorizeCustomer(ctx context.Context, authorize AuthorizeCustomer) error
		GetCustomer(ctx context.Context, get GetCustomer) (*domain.Customer, error)
		EnableCustomer(ctx context.Context, enable EnableCustomer) error
		DisableCustomer(ctx context.Context, disable DisableCustomer) error
	}

	Application struct {
		customers       domain.CustomerRepository
		domainPublisher ddd.EventPublisher[ddd.AggregateEvent]
	}
)

var _ App = (*Application)(nil)

func New(customers domain.CustomerRepository, domainPublisher ddd.EventPublisher[ddd.AggregateEvent]) *Application {
	return &Application{
		customers:       customers,
		domainPublisher: domainPublisher,
	}
}

func (a Application) RegisterCustomer(ctx context.Context, register RegisterCustomer) error {
	customer, err := domain.RegisterCustomer(register.ID, register.Name, register.SmsNumber)
	if err != nil {
		return err
	}

	if err = a.customers.Save(ctx, customer); err != nil {
		return err
	}

	// publish domain events
	if err = a.domainPublisher.Publish(ctx, customer.Events()...); err != nil {
		return err
	}

	return nil
}

func (a Application) AuthorizeCustomer(ctx context.Context, authorize AuthorizeCustomer) error {
	customer, err := a.customers.Find(ctx, authorize.ID)
	if err != nil {
		return err
	}

	if err = customer.Authorize(); err != nil {
		return err
	}

	// publish domain events
	if err = a.domainPublisher.Publish(ctx, customer.Events()...); err != nil {
		return err
	}

	return nil
}

func (a Application) EnableCustomer(ctx context.Context, enable EnableCustomer) error {
	customer, err := a.customers.Find(ctx, enable.ID)
	if err != nil {
		return err
	}

	if err = customer.Enable(); err != nil {
		return err
	}

	if err = a.customers.Update(ctx, customer); err != nil {
		return err
	}

	// publish domain events
	if err = a.domainPublisher.Publish(ctx, customer.Events()...); err != nil {
		return err
	}

	return nil
}

func (a Application) DisableCustomer(ctx context.Context, disable DisableCustomer) error {
	customer, err := a.customers.Find(ctx, disable.ID)
	if err != nil {
		return err
	}

	if err = customer.Disable(); err != nil {
		return err
	}

	if err = a.customers.Update(ctx, customer); err != nil {
		return err
	}

	// publish domain events
	if err = a.domainPublisher.Publish(ctx, customer.Events()...); err != nil {
		return err
	}

	return nil
}

func (a Application) GetCustomer(ctx context.Context, get GetCustomer) (*domain.Customer, error) {
	return a.customers.Find(ctx, get.ID)
}
