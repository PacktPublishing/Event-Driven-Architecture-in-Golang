package e2e

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/go-openapi/runtime/client"
	"github.com/stackus/errors"
)

type featureContext interface {
	register(ctx *godog.ScenarioContext)
	reset()
}

type suiteContext struct {
	transport *client.Runtime
	db        *sql.DB
	lastErr   error
	// global data
	customers map[string]string
	stores    map[string]string
}

type suiteConfig struct {
	paths       []string
	featureCtxs []featureContext
}

func newTestSuite(cfg suiteConfig) godog.TestSuite {
	return godog.TestSuite{
		Name: "mallbots-e2e",
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			for _, featureCtx := range cfg.featureCtxs {
				featureCtx.register(ctx)
			}
			ctx.Before(func(ctx context.Context, s *godog.Scenario) (context.Context, error) {
				for _, featureCtx := range cfg.featureCtxs {
					featureCtx.reset()
				}

				return ctx, nil
			})
		},
		Options: &godog.Options{
			Format:    "pretty",
			Paths:     cfg.paths,
			Randomize: -1,
		},
	}
}

func (c *suiteContext) register(ctx *godog.ScenarioContext) {
	ctx.Step(`^I expect the (?:request|command|query) to fail$`, c.iExpectTheCommandToFail)
	ctx.Step(`^I expect the (?:request|command|query) to succeed$`, c.iExpectTheCommandToSucceed)

	ctx.Step(`^(?:ensure |expect )?the returned error message is "([^"]*)"$`, c.theReturnedErrorMessageIs)
}

func (c *suiteContext) reset() {
	c.lastErr = nil
}

func (c *suiteContext) truncate(tableName string) {
	_, _ = c.db.Exec(fmt.Sprintf("TRUNCATE %s", tableName))
}

func (c *suiteContext) iExpectTheCommandToFail() error {
	if c.lastErr == nil {
		return errors.Wrap(errors.ErrUnknown, "expected error to not be nil")
	}
	return nil
}

func (c *suiteContext) iExpectTheCommandToSucceed() error {
	if c.lastErr != nil {
		return errors.Wrapf(c.lastErr, "expected error to be nil: got %s", c.lastErr)
	}

	return nil
}

func (c *suiteContext) theReturnedErrorMessageIs(errorMsg string) error {
	if c.lastErr == nil {
		return errors.Wrap(errors.ErrUnknown, "expected error to not be nil")
	}

	if errorMsg != c.lastErr.Error() {
		return errors.Wrapf(errors.ErrInvalidArgument, "expected: %s: got: %s", errorMsg, c.lastErr.Error())
	}

	return nil
}
