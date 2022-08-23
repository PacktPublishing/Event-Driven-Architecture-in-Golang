package domain

type CustomerRegistered struct {
	Customer *Customer
}

func (CustomerRegistered) EventName() string { return "customers.CustomerRegistered" }

type CustomerAuthorized struct {
	Customer *Customer
}

func (CustomerAuthorized) EventName() string { return "customers.CustomerAuthorized" }

type CustomerEnabled struct {
	Customer *Customer
}

func (CustomerEnabled) EventName() string { return "customers.CustomerEnabled" }

type CustomerDisabled struct {
	Customer *Customer
}

func (CustomerDisabled) EventName() string { return "customers.CustomerDisabled" }
