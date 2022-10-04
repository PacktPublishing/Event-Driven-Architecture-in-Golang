//go:build integration

package grpc

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"eda-in-golang/baskets/basketspb"
	"eda-in-golang/baskets/internal/application"
	"eda-in-golang/baskets/internal/domain"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/es"
)

type serverSuite struct {
	mocks struct {
		baskets   *domain.MockBasketRepository
		stores    *domain.MockStoreRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}
	server *grpc.Server
	client basketspb.BasketServiceClient
	suite.Suite
}

func TestServer(t *testing.T) {
	suite.Run(t, &serverSuite{})
}

func (s *serverSuite) SetupSuite()    {}
func (s *serverSuite) TearDownSuite() {}

func (s *serverSuite) SetupTest() {
	const grpcTestPort = ":10912"

	var err error
	// create server
	s.server = grpc.NewServer()
	var listener net.Listener
	listener, err = net.Listen("tcp", grpcTestPort)
	if err != nil {
		s.T().Fatal(err)
	}

	// create mocks
	s.mocks = struct {
		baskets   *domain.MockBasketRepository
		stores    *domain.MockStoreRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}{
		baskets:   domain.NewMockBasketRepository(s.T()),
		stores:    domain.NewMockStoreRepository(s.T()),
		products:  domain.NewMockProductRepository(s.T()),
		publisher: ddd.NewMockEventPublisher[ddd.Event](s.T()),
	}

	// create app
	app := application.New(s.mocks.baskets, s.mocks.stores, s.mocks.products, s.mocks.publisher)

	// register app with server
	if err = RegisterServer(app, s.server); err != nil {
		s.T().Fatal(err)
	}
	go func(listener net.Listener) {
		err := s.server.Serve(listener)
		if err != nil {
			s.T().Fatal(err)
		}
	}(listener)

	// create client
	var conn *grpc.ClientConn
	conn, err = grpc.Dial(grpcTestPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.T().Fatal(err)
	}
	s.client = basketspb.NewBasketServiceClient(conn)
}
func (s *serverSuite) TearDownTest() {
	s.server.GracefulStop()
}

func (s *serverSuite) TestBasketService_StartBasket() {
	s.mocks.baskets.On("Load", mock.Anything, mock.AnythingOfType("string")).Return(&domain.Basket{
		Aggregate: es.NewAggregate("basket-id", domain.BasketAggregate),
	}, nil)
	s.mocks.baskets.On("Save", mock.Anything, mock.AnythingOfType("*domain.Basket")).Return(nil)
	s.mocks.publisher.On("Publish", mock.Anything, mock.AnythingOfType("ddd.event")).Return(nil)

	_, err := s.client.StartBasket(context.Background(), &basketspb.StartBasketRequest{CustomerId: "customer-id"})
	s.Assert().NoError(err)
}

func (s *serverSuite) TestBasketService_CancelBasket() {
	s.mocks.baskets.On("Load", mock.Anything, "basket-id").Return(&domain.Basket{
		Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
		CustomerID: "customer-id",
		Status:     domain.BasketIsOpen,
	}, nil)
	s.mocks.baskets.On("Save", mock.Anything, mock.AnythingOfType("*domain.Basket")).Return(nil)
	s.mocks.publisher.On("Publish", mock.Anything, mock.AnythingOfType("ddd.event")).Return(nil)

	_, err := s.client.CancelBasket(context.Background(), &basketspb.CancelBasketRequest{Id: "basket-id"})
	s.Assert().NoError(err)
}

func (s *serverSuite) TestBasketService_CheckoutBasket() {
	s.mocks.baskets.On("Load", mock.Anything, "basket-id").Return(&domain.Basket{
		Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
		CustomerID: "customer-id",
		Items: map[string]domain.Item{
			"product-id": {
				StoreID:      "store-id",
				ProductID:    "product-id",
				StoreName:    "store-name",
				ProductName:  "product-name",
				ProductPrice: 1.00,
				Quantity:     1,
			},
		},
		Status: domain.BasketIsOpen,
	}, nil)
	s.mocks.baskets.On("Save", mock.Anything, mock.AnythingOfType("*domain.Basket")).Return(nil)
	s.mocks.publisher.On("Publish", mock.Anything, mock.AnythingOfType("ddd.event")).Return(nil)

	_, err := s.client.CheckoutBasket(context.Background(), &basketspb.CheckoutBasketRequest{
		Id:        "basket-id",
		PaymentId: "payment-id",
	})
	s.Assert().NoError(err)
}

func (s *serverSuite) TestBasketService_AddItem() {
	product := &domain.Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}
	store := &domain.Store{
		ID:   "store-id",
		Name: "store-name",
	}
	s.mocks.baskets.On("Load", mock.Anything, "basket-id").Return(&domain.Basket{
		Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
		CustomerID: "customer-id",
		Items: map[string]domain.Item{
			"product-id": {
				StoreID:      "store-id",
				ProductID:    "product-id",
				StoreName:    "store-name",
				ProductName:  "product-name",
				ProductPrice: 1.00,
				Quantity:     1,
			},
		},
		Status: domain.BasketIsOpen,
	}, nil)
	s.mocks.baskets.On("Save", mock.Anything, mock.AnythingOfType("*domain.Basket")).Return(nil)
	s.mocks.products.On("Find", mock.Anything, "product-id").Return(product, nil)
	s.mocks.stores.On("Find", mock.Anything, "store-id").Return(store, nil)

	_, err := s.client.AddItem(context.Background(), &basketspb.AddItemRequest{
		Id:        "basket-id",
		ProductId: "product-id",
		Quantity:  1,
	})
	s.Assert().NoError(err)
}

func (s *serverSuite) TestBasketService_RemoveItem() {
	product := &domain.Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}

	s.mocks.baskets.On("Load", mock.Anything, "basket-id").Return(&domain.Basket{
		Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
		CustomerID: "customer-id",
		Items: map[string]domain.Item{
			"product-id": {
				StoreID:      "store-id",
				ProductID:    "product-id",
				StoreName:    "store-name",
				ProductName:  "product-name",
				ProductPrice: 1.00,
				Quantity:     1,
			},
		},
		Status: domain.BasketIsOpen,
	}, nil)
	s.mocks.baskets.On("Save", mock.Anything, mock.AnythingOfType("*domain.Basket")).Return(nil)
	s.mocks.products.On("Find", mock.Anything, "product-id").Return(product, nil)

	_, err := s.client.RemoveItem(context.Background(), &basketspb.RemoveItemRequest{
		Id:        "basket-id",
		ProductId: "product-id",
		Quantity:  1,
	})
	s.Assert().NoError(err)
}

func (s *serverSuite) TestBasketService_GetBasket() {
	basket := &domain.Basket{
		Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
		CustomerID: "customer-id",
		Items: map[string]domain.Item{
			"product-id": {
				StoreID:      "store-id",
				ProductID:    "product-id",
				StoreName:    "store-name",
				ProductName:  "product-name",
				ProductPrice: 1.00,
				Quantity:     1,
			},
		},
		Status: domain.BasketIsOpen,
	}
	s.mocks.baskets.On("Load", mock.Anything, "basket-id").Return(basket, nil)

	resp, err := s.client.GetBasket(context.Background(), &basketspb.GetBasketRequest{Id: "basket-id"})
	if s.Assert().NoError(err) {
		s.Assert().Equal(basket.ID(), resp.Basket.GetId())
	}
}
