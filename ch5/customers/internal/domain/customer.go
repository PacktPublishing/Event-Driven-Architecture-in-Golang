package domain

import (
	"github.com/stackus/errors"

	"eda-in-golang/internal/ddd"
)

const CustomerAggregate = "customers.CustomerAggregate"

type Customer struct {
	ddd.Aggregate
	Name      string
	SmsNumber string
	Enabled   bool
}

var (
	ErrNameCannotBeBlank       = errors.Wrap(errors.ErrBadRequest, "the customer name cannot be blank")
	ErrCustomerIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the customer id cannot be blank")
	ErrSmsNumberCannotBeBlank  = errors.Wrap(errors.ErrBadRequest, "the SMS number cannot be blank")
	ErrCustomerAlreadyEnabled  = errors.Wrap(errors.ErrBadRequest, "the customer is already enabled")
	ErrCustomerAlreadyDisabled = errors.Wrap(errors.ErrBadRequest, "the customer is already disabled")
	ErrCustomerNotAuthorized   = errors.Wrap(errors.ErrUnauthorized, "customer is not authorized")
)

func NewCustomer(id string) *Customer {
	return &Customer{
		Aggregate: ddd.NewAggregate(id, CustomerAggregate),
	}
}

func RegisterCustomer(id, name, smsNumber string) (*Customer, error) {
	if id == "" {
		return nil, ErrCustomerIDCannotBeBlank
	}

	if name == "" {
		return nil, ErrNameCannotBeBlank
	}

	if smsNumber == "" {
		return nil, ErrSmsNumberCannotBeBlank
	}

	customer := NewCustomer(id)
	customer.Name = name
	customer.SmsNumber = smsNumber
	customer.Enabled = true

	customer.AddEvent(CustomerRegisteredEvent, &CustomerRegistered{
		Customer: customer,
	})

	return customer, nil
}

func (Customer) Key() string { return CustomerAggregate }

func (c *Customer) Authorize( /* TODO authorize what? */ ) error {
	if !c.Enabled {
		return ErrCustomerNotAuthorized
	}

	c.AddEvent(CustomerAuthorizedEvent, &CustomerAuthorized{
		Customer: c,
	})

	return nil
}

func (c *Customer) Enable() error {
	if c.Enabled {
		return ErrCustomerAlreadyEnabled
	}

	c.Enabled = true

	c.AddEvent(CustomerEnabledEvent, &CustomerEnabled{
		Customer: c,
	})

	return nil
}

func (c *Customer) Disable() error {
	if !c.Enabled {
		return ErrCustomerAlreadyDisabled
	}

	c.Enabled = false

	c.AddEvent(CustomerDisabledEvent, &CustomerDisabled{
		Customer: c,
	})

	return nil
}
