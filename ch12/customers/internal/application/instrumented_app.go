package application

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

type instrumentedApp struct {
	App
	customersRegistered prometheus.Counter
}

var _ App = (*instrumentedApp)(nil)

func NewInstrumentedApp(app App, customersRegistered prometheus.Counter) App {
	return instrumentedApp{
		App:                 app,
		customersRegistered: customersRegistered,
	}
}

func (a instrumentedApp) RegisterCustomer(ctx context.Context, register RegisterCustomer) error {
	err := a.App.RegisterCustomer(ctx, register)
	if err != nil {
		return err
	}
	a.customersRegistered.Inc()
	return nil
}
