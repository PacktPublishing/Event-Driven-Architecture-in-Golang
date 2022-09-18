package application

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"eda-in-golang/baskets/internal/domain"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/es"
)

func TestApplication_AddItem(t *testing.T) {
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

	type mocks struct {
		baskets   *domain.MockBasketRepository
		stores    *domain.MockStoreRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}
	type args struct {
		ctx context.Context
		add AddItem
	}
	tests := map[string]struct {
		args    args
		on      func(f mocks)
		wantErr bool
	}{
		"Success": {
			args: args{
				ctx: context.Background(),
				add: AddItem{
					ID:        "basket-id",
					ProductID: "product-id",
					Quantity:  1,
				},
			},
			on: func(f mocks) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      make(map[string]domain.Item),
					Status:     domain.BasketIsOpen,
				}, nil)
				f.products.On("Find", context.Background(), "product-id").Return(product, nil)
				f.stores.On("Find", context.Background(), "store-id").Return(store, nil)
				f.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(nil)
			},
		},
		"NoBasket": {
			args: args{
				ctx: context.Background(),
				add: AddItem{
					ID:        "basket-id",
					ProductID: "product-id",
					Quantity:  1,
				},
			},
			on: func(f mocks) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(nil, fmt.Errorf("no basket"))
			},
			wantErr: true,
		},
		"NoProduct": {
			args: args{
				ctx: context.Background(),
				add: AddItem{
					ID:        "basket-id",
					ProductID: "product-id",
					Quantity:  1,
				},
			},
			on: func(f mocks) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      make(map[string]domain.Item),
					Status:     domain.BasketIsOpen,
				}, nil)
				f.products.On("Find", context.Background(), "product-id").Return(nil, fmt.Errorf("no product"))
			},
			wantErr: true,
		},
		"NoStore": {
			args: args{
				ctx: context.Background(),
				add: AddItem{
					ID:        "basket-id",
					ProductID: "product-id",
					Quantity:  1,
				},
			},
			on: func(f mocks) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      make(map[string]domain.Item),
					Status:     domain.BasketIsOpen,
				}, nil)
				f.products.On("Find", context.Background(), "product-id").Return(product, nil)
				f.stores.On("Find", context.Background(), "store-id").Return(nil, fmt.Errorf("no store"))
			},
			wantErr: true,
		},
		"SaveFailed": {
			args: args{
				ctx: context.Background(),
				add: AddItem{
					ID:        "basket-id",
					ProductID: "product-id",
					Quantity:  1,
				},
			},
			on: func(f mocks) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      make(map[string]domain.Item),
					Status:     domain.BasketIsOpen,
				}, nil)
				f.products.On("Find", context.Background(), "product-id").Return(product, nil)
				f.stores.On("Find", context.Background(), "store-id").Return(store, nil)
				f.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(fmt.Errorf("save failed"))
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			m := mocks{
				baskets:   domain.NewMockBasketRepository(t),
				stores:    domain.NewMockStoreRepository(t),
				products:  domain.NewMockProductRepository(t),
				publisher: ddd.NewMockEventPublisher[ddd.Event](t),
			}
			a := New(m.baskets, m.stores, m.products, m.publisher)
			if tt.on != nil {
				tt.on(m)
			}

			if err := a.AddItem(tt.args.ctx, tt.args.add); (err != nil) != tt.wantErr {
				t.Errorf("AddItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_CancelBasket(t *testing.T) {
	type fields struct {
		baskets   *domain.MockBasketRepository
		stores    *domain.MockStoreRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}
	type args struct {
		ctx    context.Context
		cancel CancelBasket
	}
	tests := map[string]struct {
		args    args
		on      func(f fields)
		wantErr bool
	}{
		"Success": {
			args: args{
				ctx: context.Background(),
				cancel: CancelBasket{
					ID: "basket-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "",
					Items:      make(map[string]domain.Item),
					Status:     domain.BasketIsOpen,
				}, nil)
				f.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(nil)
				f.publisher.On("Publish", context.Background(), mock.AnythingOfType("ddd.event")).Return(nil)
			},
		},
		"NoBasket": {
			args: args{
				ctx: context.Background(),
				cancel: CancelBasket{
					ID: "basket-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(nil, fmt.Errorf("no basket"))
			},
			wantErr: true,
		},
		"SaveFailed": {
			args: args{
				ctx: context.Background(),
				cancel: CancelBasket{
					ID: "basket-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "",
					Items:      make(map[string]domain.Item),
					Status:     domain.BasketIsOpen,
				}, nil)
				f.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(fmt.Errorf("save failed"))
			},
			wantErr: true,
		},
		"PublishFailed": {
			args: args{
				ctx: context.Background(),
				cancel: CancelBasket{
					ID: "basket-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "",
					Items:      make(map[string]domain.Item),
					Status:     domain.BasketIsOpen,
				}, nil)
				f.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(nil)
				f.publisher.On("Publish", context.Background(), mock.AnythingOfType("ddd.event")).Return(fmt.Errorf("publish failed"))
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := fields{
				baskets:   domain.NewMockBasketRepository(t),
				stores:    domain.NewMockStoreRepository(t),
				products:  domain.NewMockProductRepository(t),
				publisher: ddd.NewMockEventPublisher[ddd.Event](t),
			}
			a := Application{
				baskets:   f.baskets,
				stores:    f.stores,
				products:  f.products,
				publisher: f.publisher,
			}
			if tt.on != nil {
				tt.on(f)
			}

			if err := a.CancelBasket(tt.args.ctx, tt.args.cancel); (err != nil) != tt.wantErr {
				t.Errorf("CancelBasket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_CheckoutBasket(t *testing.T) {
	store := &domain.Store{
		ID:   "store-id",
		Name: "store-name",
	}
	product := &domain.Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}
	item := domain.Item{
		StoreID:      store.ID,
		ProductID:    product.ID,
		StoreName:    store.Name,
		ProductName:  product.Name,
		ProductPrice: product.Price,
		Quantity:     10,
	}

	type fields struct {
		baskets   *domain.MockBasketRepository
		stores    *domain.MockStoreRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}
	type args struct {
		ctx      context.Context
		checkout CheckoutBasket
	}
	tests := map[string]struct {
		args    args
		on      func(f fields)
		wantErr bool
	}{
		"Success": {
			args: args{
				ctx: context.Background(),
				checkout: CheckoutBasket{
					ID:        "basket-id",
					PaymentID: "payment-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "",
					Items: map[string]domain.Item{
						product.ID: item,
					},
					Status: domain.BasketIsOpen,
				}, nil)
				f.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(nil)
				f.publisher.On("Publish", context.Background(), mock.AnythingOfType("ddd.event")).Return(nil)
			},
		},
		"MissingPaymentID": {
			args: args{
				ctx: context.Background(),
				checkout: CheckoutBasket{
					ID:        "basket-id",
					PaymentID: "",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "",
					Items: map[string]domain.Item{
						product.ID: item,
					},
					Status: domain.BasketIsOpen,
				}, nil)
			},
			wantErr: true,
		},
		"NoBasket": {
			args: args{
				ctx: context.Background(),
				checkout: CheckoutBasket{
					ID:        "basket-id",
					PaymentID: "payment-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(nil, fmt.Errorf("no basket"))
			},
			wantErr: true,
		},
		"SaveFailed": {
			args: args{
				ctx: context.Background(),
				checkout: CheckoutBasket{
					ID:        "basket-id",
					PaymentID: "payment-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "",
					Items: map[string]domain.Item{
						product.ID: item,
					},
					Status: domain.BasketIsOpen,
				}, nil)
				f.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(fmt.Errorf("save failed"))
			},
			wantErr: true,
		},
		"PublishFailed": {
			args: args{
				ctx: context.Background(),
				checkout: CheckoutBasket{
					ID:        "basket-id",
					PaymentID: "payment-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "",
					Items: map[string]domain.Item{
						product.ID: item,
					},
					Status: domain.BasketIsOpen,
				}, nil)
				f.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(nil)
				f.publisher.On("Publish", context.Background(), mock.AnythingOfType("ddd.event")).Return(fmt.Errorf("publish failed"))
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := fields{
				baskets:   domain.NewMockBasketRepository(t),
				stores:    domain.NewMockStoreRepository(t),
				products:  domain.NewMockProductRepository(t),
				publisher: ddd.NewMockEventPublisher[ddd.Event](t),
			}
			a := Application{
				baskets:   f.baskets,
				stores:    f.stores,
				products:  f.products,
				publisher: f.publisher,
			}
			if tt.on != nil {
				tt.on(f)
			}

			if err := a.CheckoutBasket(tt.args.ctx, tt.args.checkout); (err != nil) != tt.wantErr {
				t.Errorf("CheckoutBasket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_GetBasket(t *testing.T) {
	store := &domain.Store{
		ID:   "store-id",
		Name: "store-name",
	}
	product := &domain.Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}
	item := domain.Item{
		StoreID:      store.ID,
		ProductID:    product.ID,
		StoreName:    store.Name,
		ProductName:  product.Name,
		ProductPrice: product.Price,
		Quantity:     10,
	}

	type fields struct {
		baskets   *domain.MockBasketRepository
		stores    *domain.MockStoreRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}
	type args struct {
		ctx context.Context
		get GetBasket
	}
	tests := map[string]struct {
		args    args
		on      func(f fields)
		want    *domain.Basket
		wantErr bool
	}{
		"GetBasket": {
			args: args{
				ctx: context.Background(),
				get: GetBasket{
					ID: "basket-id",
				},
			},
			on: func(f fields) {
				f.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items: map[string]domain.Item{
						product.ID: item,
					},
					Status: domain.BasketIsOpen,
				}, nil)
			},
			want: &domain.Basket{
				Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]domain.Item{
					product.ID: item,
				},
				Status: domain.BasketIsOpen,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			f := fields{
				baskets:   domain.NewMockBasketRepository(t),
				stores:    domain.NewMockStoreRepository(t),
				products:  domain.NewMockProductRepository(t),
				publisher: ddd.NewMockEventPublisher[ddd.Event](t),
			}
			a := Application{
				baskets:   f.baskets,
				stores:    f.stores,
				products:  f.products,
				publisher: f.publisher,
			}
			if tt.on != nil {
				tt.on(f)
			}

			// Act
			got, err := a.GetBasket(tt.args.ctx, tt.args.get)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBasket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want.ID(), tt.want.ID())
			assert.Equal(t, tt.want.CustomerID, got.CustomerID)
			assert.Equal(t, tt.want.PaymentID, got.PaymentID)
			assert.Equal(t, tt.want.Items, got.Items)
			assert.Equal(t, tt.want.Status, got.Status)
		})
	}
}

func TestApplication_RemoveItem(t *testing.T) {
	store := &domain.Store{
		ID:   "store-id",
		Name: "store-name",
	}
	product := &domain.Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}
	item := domain.Item{
		StoreID:      store.ID,
		ProductID:    product.ID,
		StoreName:    store.Name,
		ProductName:  product.Name,
		ProductPrice: product.Price,
		Quantity:     10,
	}

	type mocks struct {
		baskets   *domain.MockBasketRepository
		stores    *domain.MockStoreRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}
	type args struct {
		ctx    context.Context
		remove RemoveItem
	}
	tests := map[string]struct {
		args    args
		on      func(m mocks)
		wantErr bool
	}{
		"Success": {
			args: args{
				ctx: context.Background(),
				remove: RemoveItem{
					ID:        "basket-id",
					ProductID: product.ID,
					Quantity:  1,
				},
			},
			on: func(m mocks) {
				m.products.On("Find", context.Background(), product.ID).Return(product, nil)
				m.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "",
					Items: map[string]domain.Item{
						product.ID: item,
					},
					Status: domain.BasketIsOpen,
				}, nil)
				m.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(nil)
			},
		},
		"NoProduct": {
			args: args{
				ctx: context.Background(),
				remove: RemoveItem{
					ID:        "basket-id",
					ProductID: product.ID,
					Quantity:  1,
				},
			},
			on: func(m mocks) {
				m.products.On("Find", context.Background(), product.ID).Return(nil, fmt.Errorf("no product"))
			},
			wantErr: true,
		},
		"NoBasket": {
			args: args{
				ctx: context.Background(),
				remove: RemoveItem{
					ID:        "basket-id",
					ProductID: product.ID,
					Quantity:  1,
				},
			},
			on: func(m mocks) {
				m.products.On("Find", context.Background(), product.ID).Return(product, nil)
				m.baskets.On("Load", context.Background(), "basket-id").Return(nil, fmt.Errorf("no basket"))
			},
			wantErr: true,
		},
		"SaveFailed": {
			args: args{
				ctx: context.Background(),
				remove: RemoveItem{
					ID:        "basket-id",
					ProductID: product.ID,
					Quantity:  1,
				},
			},
			on: func(m mocks) {
				m.products.On("Find", context.Background(), product.ID).Return(product, nil)
				m.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "",
					Items: map[string]domain.Item{
						product.ID: item,
					},
					Status: domain.BasketIsOpen,
				}, nil)
				m.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(fmt.Errorf("save failed"))
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			m := mocks{
				baskets:   domain.NewMockBasketRepository(t),
				stores:    domain.NewMockStoreRepository(t),
				products:  domain.NewMockProductRepository(t),
				publisher: ddd.NewMockEventPublisher[ddd.Event](t),
			}
			a := Application{
				baskets:   m.baskets,
				stores:    m.stores,
				products:  m.products,
				publisher: m.publisher,
			}
			if tt.on != nil {
				tt.on(m)
			}

			if err := a.RemoveItem(tt.args.ctx, tt.args.remove); (err != nil) != tt.wantErr {
				t.Errorf("RemoveItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_StartBasket(t *testing.T) {
	type mocks struct {
		baskets   *domain.MockBasketRepository
		stores    *domain.MockStoreRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}
	type args struct {
		ctx   context.Context
		start StartBasket
	}
	tests := map[string]struct {
		args    args
		on      func(m mocks)
		wantErr bool
	}{
		"Success": {
			args: args{
				ctx: context.Background(),
				start: StartBasket{
					ID:         "basket-id",
					CustomerID: "customer-id",
				},
			},
			on: func(m mocks) {
				m.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "",
					Items:      make(map[string]domain.Item),
				}, nil)
				m.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(nil)
				m.publisher.On("Publish", context.Background(), mock.AnythingOfType("ddd.event")).Return(nil)
			},
		},
		"NoBasket": {
			args: args{
				ctx: context.Background(),
				start: StartBasket{
					ID:         "basket-id",
					CustomerID: "customer-id",
				},
			},
			on: func(m mocks) {
				m.baskets.On("Load", context.Background(), "basket-id").Return(nil, fmt.Errorf("no basket"))
			},
			wantErr: true,
		},
		"SaveFailed": {
			args: args{
				ctx: context.Background(),
				start: StartBasket{
					ID:         "basket-id",
					CustomerID: "customer-id",
				},
			},
			on: func(m mocks) {
				m.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "",
					Items:      make(map[string]domain.Item),
				}, nil)
				m.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(fmt.Errorf("save failed"))
			},
			wantErr: true,
		},
		"PublishFailed": {
			args: args{
				ctx: context.Background(),
				start: StartBasket{
					ID:         "basket-id",
					CustomerID: "customer-id",
				},
			},
			on: func(m mocks) {
				m.baskets.On("Load", context.Background(), "basket-id").Return(&domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "",
					Items:      make(map[string]domain.Item),
				}, nil)
				m.baskets.On("Save", context.Background(), mock.AnythingOfType("*domain.Basket")).Return(nil)
				m.publisher.On("Publish", context.Background(), mock.AnythingOfType("ddd.event")).Return(fmt.Errorf("publish failed"))
			},
			wantErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := mocks{
				baskets:   domain.NewMockBasketRepository(t),
				stores:    domain.NewMockStoreRepository(t),
				products:  domain.NewMockProductRepository(t),
				publisher: ddd.NewMockEventPublisher[ddd.Event](t),
			}
			a := Application{
				baskets:   m.baskets,
				stores:    m.stores,
				products:  m.products,
				publisher: m.publisher,
			}
			if tc.on != nil {
				tc.on(m)
			}

			if err := a.StartBasket(tc.args.ctx, tc.args.start); (err != nil) != tc.wantErr {
				t.Errorf("StartBasket() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
