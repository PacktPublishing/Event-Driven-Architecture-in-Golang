package domain

const (
	CustomerRegisteredEvent = "customers.CustomerRegistered"
	CustomerAuthorizedEvent = "customers.CustomerAuthorized"
	CustomerEnabledEvent    = "customers.CustomerEnabled"
	CustomerDisabledEvent   = "customers.CustomerDisabled"
)

type CustomerRegistered struct {
	Customer *Customer
}

func (CustomerRegistered) Key() string { return CustomerRegisteredEvent }

type CustomerAuthorized struct {
	Customer *Customer
}

func (CustomerAuthorized) Key() string { return CustomerAuthorizedEvent }

type CustomerEnabled struct {
	Customer *Customer
}

func (CustomerEnabled) Key() string { return CustomerEnabledEvent }

type CustomerDisabled struct {
	Customer *Customer
}

func (CustomerDisabled) Key() string { return CustomerDisabledEvent }
