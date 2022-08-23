//go:build e2e

package e2e

import (
	"database/sql"
	"testing"

	"github.com/go-openapi/runtime/client"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func TestEndToEnd(t *testing.T) {
	sc := &suiteContext{
		transport: client.New("localhost:8080", "/", nil),
	}

	var err error
	sc.db, err = sql.Open("pgx", "postgres://mallbots_user:mallbots_pass@localhost:5432/mallbots?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	suite := newTestSuite(suiteConfig{
		paths: []string{
			"features/baskets",
			"features/customers",
			"features/kiosk",
			"features/orders",
			"features/stores",
		},
		featureCtxs: []featureContext{
			sc,
			newCustomersContext(sc),
			newStoresContext(sc),
		},
	})

	if status := suite.Run(); status != 0 {
		t.Error("end to end feature test failed with status:", status)
	}
}
