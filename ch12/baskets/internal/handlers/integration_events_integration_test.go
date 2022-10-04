//go:build integration

package handlers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"eda-in-golang/baskets/internal/domain"
	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/jetstream"
	"eda-in-golang/internal/registry"
	"eda-in-golang/stores/storespb"
)

const streamName = "mallbots"

type integrationEventsTestSuite struct {
	container  testcontainers.Container
	natsConn   *nats.Conn
	reg        registry.Registry
	js         nats.JetStreamContext
	publisher  am.EventPublisher
	subscriber am.MessageStream
	mocks      struct {
		products *domain.MockProductCacheRepository
		stores   *domain.MockStoreCacheRepository
	}
	suite.Suite
}

func TestIntegrationEvents(t *testing.T) {
	suite.Run(t, &integrationEventsTestSuite{})
}

func (s *integrationEventsTestSuite) SetupSuite() {
	var err error
	ctx := context.Background()
	natsReq := testcontainers.ContainerRequest{
		Image:        "nats:2-alpine",
		Hostname:     "nats",
		ExposedPorts: []string{"4222/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("4222/tcp"),
			wait.ForLog("Server is ready"),
		),
		Cmd: []string{"-js"},
	}
	s.container, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: natsReq,
		Started:          true,
	})
	if err != nil {
		s.T().Fatal(err)
	}

	s.reg = registry.New()
	if err = storespb.Registrations(s.reg); err != nil {
		s.T().Fatal(err)
	}

	endpoint, err := s.container.Endpoint(ctx, "")
	if err != nil {
		s.T().Fatal(err)
	}

	s.natsConn, err = nats.Connect(
		endpoint,
		nats.Timeout(5*time.Second),
		nats.RetryOnFailedConnect(true),
	)

	if err != nil {
		s.T().Fatal(err)
	}
	s.js, err = s.natsConn.JetStream()
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = s.js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{fmt.Sprintf("%s.>", streamName)},
	})
	if err != nil {
		s.T().Fatal(err)
	}

}

func (s *integrationEventsTestSuite) TearDownSuite() {
	s.natsConn.Close()
	if err := s.container.Terminate(context.Background()); err != nil {
		s.T().Fatal(err)
	}
}

func (s *integrationEventsTestSuite) SetupTest() {
	s.mocks = struct {
		products *domain.MockProductCacheRepository
		stores   *domain.MockStoreCacheRepository
	}{
		products: domain.NewMockProductCacheRepository(s.T()),
		stores:   domain.NewMockStoreCacheRepository(s.T()),
	}

	logger := zerolog.New(zerolog.NewConsoleWriter()).
		Level(zerolog.DebugLevel).
		With().
		Logger()

	stream := jetstream.NewStream(streamName, s.js, logger)
	s.publisher = am.NewEventPublisher(s.reg, stream)
	s.subscriber = stream
	handler := am.NewEventHandler(s.reg, integrationHandlers[ddd.Event]{
		products: s.mocks.products,
		stores:   s.mocks.stores,
	})

	if err := RegisterIntegrationEventHandlers(s.subscriber, handler); err != nil {
		s.T().Fatal(err)
	}
}

func (s *integrationEventsTestSuite) TearDownTest() {
	if err := s.subscriber.Unsubscribe(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *integrationEventsTestSuite) TestStoreAggregateChannel_StoreCreated() {
	s.wait(func(done chan struct{}) {
		s.mocks.stores.On("Add", mock.Anything, "store-id", "store-name").Return(nil).Run(func(_ mock.Arguments) {
			close(done)
		})

		_ = s.publisher.Publish(context.Background(), storespb.StoreAggregateChannel,
			ddd.NewEvent(storespb.StoreCreatedEvent, &storespb.StoreCreated{
				Id:       "store-id",
				Name:     "store-name",
				Location: "store-location",
			}),
		)
	})
}

func (s *integrationEventsTestSuite) TestStoreAggregateChannel_StoreRebranded() {
	s.wait(func(done chan struct{}) {
		s.mocks.stores.On("Rename", mock.Anything, "store-id", "store-name").Return(nil).Run(func(_ mock.Arguments) {
			close(done)
		})

		s.NoError(s.publisher.Publish(context.Background(), storespb.StoreAggregateChannel,
			ddd.NewEvent(storespb.StoreRebrandedEvent, &storespb.StoreRebranded{
				Id:   "store-id",
				Name: "store-name",
			}),
		))
	})
}

func (s *integrationEventsTestSuite) TestProductAggregateChannel_ProductAdded() {
	s.wait(func(done chan struct{}) {
		s.mocks.products.On("Add", mock.Anything, "product-id", "store-id", "product-name", 10.00).Return(nil).Run(func(_ mock.Arguments) {
			close(done)
		})

		s.NoError(s.publisher.Publish(context.Background(), storespb.ProductAggregateChannel,
			ddd.NewEvent(storespb.ProductAddedEvent, &storespb.ProductAdded{
				Id:      "product-id",
				StoreId: "store-id",
				Name:    "product-name",
				Price:   10.00,
			}),
		))
	})
}

func (s *integrationEventsTestSuite) TestProductAggregateChannel_ProductRebranded() {
	s.wait(func(done chan struct{}) {
		s.mocks.products.On("Rebrand", mock.Anything, "product-id", "product-name").Return(nil).Run(func(_ mock.Arguments) {
			close(done)
		})

		s.NoError(s.publisher.Publish(context.Background(), storespb.ProductAggregateChannel,
			ddd.NewEvent(storespb.ProductRebrandedEvent, &storespb.ProductRebranded{
				Id:   "product-id",
				Name: "product-name",
			}),
		))
	})
}

func (s *integrationEventsTestSuite) TestProductAggregateChannel_ProductPriceIncreased() {
	s.wait(func(done chan struct{}) {
		s.mocks.products.On("UpdatePrice", mock.Anything, "product-id", 1.00).Return(nil).Run(func(_ mock.Arguments) {
			close(done)
		})

		s.NoError(s.publisher.Publish(context.Background(), storespb.ProductAggregateChannel,
			ddd.NewEvent(storespb.ProductPriceIncreasedEvent, &storespb.ProductPriceChanged{
				Id:    "product-id",
				Delta: 1.00,
			}),
		))
	})
}

func (s *integrationEventsTestSuite) TestProductAggregateChannel_ProductPriceDecreased() {
	s.wait(func(done chan struct{}) {
		s.mocks.products.On("UpdatePrice", mock.Anything, "product-id", -1.00).Return(nil).Run(func(_ mock.Arguments) {
			close(done)
		})

		s.NoError(s.publisher.Publish(context.Background(), storespb.ProductAggregateChannel,
			ddd.NewEvent(storespb.ProductPriceDecreasedEvent, &storespb.ProductPriceChanged{
				Id:    "product-id",
				Delta: -1.00,
			}),
		))
	})
}

func (s *integrationEventsTestSuite) TestProductAggregateChannel_ProductRemoved() {
	s.wait(func(done chan struct{}) {
		s.mocks.products.On("Remove", mock.Anything, "product-id").Return(nil).Run(func(_ mock.Arguments) {
			close(done)
		})

		s.NoError(s.publisher.Publish(context.Background(), storespb.ProductAggregateChannel,
			ddd.NewEvent(storespb.ProductRemovedEvent, &storespb.ProductRemoved{
				Id: "product-id",
			}),
		))
	})
}

func (s *integrationEventsTestSuite) wait(aFn func(done chan struct{})) {
	done := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	aFn(done)

	select {
	case <-done:
		// time.Sleep(1 * time.Second)
	case <-ctx.Done():
		s.T().Error("test timed out")
	}
}
