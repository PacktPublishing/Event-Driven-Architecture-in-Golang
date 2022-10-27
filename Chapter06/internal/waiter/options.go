package waiter

import (
	"context"
)

type WaiterOption func(c *waiterCfg)

func ParentContext(ctx context.Context) WaiterOption {
	return func(c *waiterCfg) {
		c.parentCtx = ctx
	}
}

func CatchSignals() WaiterOption {
	return func(c *waiterCfg) {
		c.catchSignals = true
	}
}
