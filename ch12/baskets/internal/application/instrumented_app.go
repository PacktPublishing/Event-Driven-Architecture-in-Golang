package application

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

type instrumentedApp struct {
	App
	basketsStarted    prometheus.Counter
	basketsCheckedOut prometheus.Counter
	basketsCanceled   prometheus.Counter
}

var _ App = (*instrumentedApp)(nil)

func NewInstrumentedApp(app App, basketsStarted, basketsCheckedOut, baksetsCanceled prometheus.Counter) App {
	return instrumentedApp{
		App:               app,
		basketsStarted:    basketsStarted,
		basketsCheckedOut: basketsCheckedOut,
		basketsCanceled:   baksetsCanceled,
	}
}

func (a instrumentedApp) StartBasket(ctx context.Context, start StartBasket) error {
	err := a.App.StartBasket(ctx, start)
	if err != nil {
		return err
	}
	a.basketsStarted.Inc()
	return nil
}

func (a instrumentedApp) CheckoutBasket(ctx context.Context, checkout CheckoutBasket) error {
	err := a.App.CheckoutBasket(ctx, checkout)
	if err != nil {
		return err
	}
	a.basketsCheckedOut.Inc()
	return nil
}

func (a instrumentedApp) CancelBasket(ctx context.Context, cancel CancelBasket) error {
	err := a.App.CancelBasket(ctx, cancel)
	if err != nil {
		return err
	}
	a.basketsCanceled.Inc()
	return nil
}
