package domain

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/es"
)

func TestBasket_AddItem(t *testing.T) {
	store := &Store{
		ID:   "store-id",
		Name: "store-name",
	}
	product := &Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}

	type fields struct {
		CustomerID string
		PaymentID  string
		Items      map[string]Item
		Status     BasketStatus
	}
	type args struct {
		store    *Store
		product  *Product
		quantity int
	}
	tests := map[string]struct {
		fields  fields
		args    args
		on      func(a *es.MockAggregate)
		wantErr bool
	}{
		"OpenBasket": {
			fields: fields{
				Items:  make(map[string]Item),
				Status: BasketIsOpen,
			},
			args: args{
				store:    store,
				product:  product,
				quantity: 1,
			},
			on: func(a *es.MockAggregate) {
				a.On("AddEvent", BasketItemAddedEvent, &BasketItemAdded{
					Item: Item{
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     1,
					},
				})
			},
			wantErr: false,
		},
		"CheckedOutBasket": {
			fields: fields{
				Items:  make(map[string]Item),
				Status: BasketIsCheckedOut,
			},
			args: args{
				store:    store,
				product:  product,
				quantity: 1,
			},
			wantErr: true,
		},
		"CanceledOutBasket": {
			fields: fields{
				Items:  make(map[string]Item),
				Status: BasketIsCanceled,
			},
			args: args{
				store:    store,
				product:  product,
				quantity: 1,
			},
			wantErr: true,
		},
		"ZeroQuantity": {
			fields: fields{
				Items:  make(map[string]Item),
				Status: BasketIsCheckedOut,
			},
			args: args{
				store:    store,
				product:  product,
				quantity: 0,
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			aggregate := es.NewMockAggregate(t)
			b := &Basket{
				Aggregate:  aggregate,
				CustomerID: tt.fields.CustomerID,
				PaymentID:  tt.fields.PaymentID,
				Items:      tt.fields.Items,
				Status:     tt.fields.Status,
			}
			if tt.on != nil {
				tt.on(aggregate)
			}

			if err := b.AddItem(tt.args.store, tt.args.product, tt.args.quantity); (err != nil) != tt.wantErr {
				t.Errorf("AddItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBasket_ApplyEvent(t *testing.T) {
	store := &Store{
		ID:   "store-id",
		Name: "store-name",
	}
	product := &Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}
	product2 := &Product{
		ID:      "product-id2",
		StoreID: "store-id",
		Name:    "product-name2",
		Price:   100.00,
	}

	type fields struct {
		CustomerID string
		PaymentID  string
		Items      map[string]Item
		Status     BasketStatus
	}
	type args struct {
		event ddd.Event
	}
	tests := map[string]struct {
		fields  fields
		args    args
		want    fields
		wantErr bool
	}{
		"BasketItemAddedEvent": {
			fields: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsOpen,
			},
			args: args{
				event: ddd.NewEvent(BasketItemAddedEvent, &BasketItemAdded{
					Item: Item{
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     1,
					},
				}),
			},
			want: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: {
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     1,
					},
				},
				Status: BasketIsOpen,
			},
			wantErr: false,
		},
		"BasketItemAddedEvent.Quantity": {
			fields: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: {
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     1,
					},
				},
				Status: BasketIsOpen,
			},
			args: args{
				event: ddd.NewEvent(BasketItemAddedEvent, &BasketItemAdded{
					Item: Item{
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     1,
					},
				}),
			},
			want: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: {
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     2,
					},
				},
				Status: BasketIsOpen,
			},
			wantErr: false,
		},
		"BasketItemAddedEvent.Second": {
			fields: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: {
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     1,
					},
				},
				Status: BasketIsOpen,
			},
			args: args{
				event: ddd.NewEvent(BasketItemAddedEvent, &BasketItemAdded{
					Item: Item{
						StoreID:      store.ID,
						ProductID:    product2.ID,
						StoreName:    store.Name,
						ProductName:  product2.Name,
						ProductPrice: product2.Price,
						Quantity:     1,
					},
				}),
			},
			want: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: {
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     1,
					},
					product2.ID: {
						StoreID:      store.ID,
						ProductID:    product2.ID,
						StoreName:    store.Name,
						ProductName:  product2.Name,
						ProductPrice: product2.Price,
						Quantity:     1,
					},
				},
				Status: BasketIsOpen,
			},
			wantErr: false,
		},
		"BasketCanceledEvent": {
			fields: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsOpen,
			},
			args: args{event: ddd.NewEvent(BasketCanceledEvent, &BasketCanceled{})},
			want: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      map[string]Item{},
				Status:     BasketIsCanceled,
			},
			wantErr: false,
		},
		"BasketCanceledEvent.Cleared": {
			fields: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: {
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     1,
					},
				},
				Status: BasketIsOpen,
			},
			args: args{event: ddd.NewEvent(BasketCanceledEvent, &BasketCanceled{})},
			want: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      map[string]Item{},
				Status:     BasketIsCanceled,
			},
			wantErr: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b := &Basket{
				Aggregate:  es.NewMockAggregate(t),
				CustomerID: tt.fields.CustomerID,
				PaymentID:  tt.fields.PaymentID,
				Items:      tt.fields.Items,
				Status:     tt.fields.Status,
			}
			if err := b.ApplyEvent(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("ApplyEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, b.CustomerID, tt.want.CustomerID)
			assert.Equal(t, b.PaymentID, tt.want.PaymentID)
			assert.Equal(t, b.Items, tt.want.Items)
			assert.Equal(t, b.Status, tt.want.Status)
		})
	}
}

func TestBasket_ApplySnapshot(t *testing.T) {
	store := &Store{
		ID:   "store-id",
		Name: "store-name",
	}
	product := &Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}
	item := Item{
		StoreID:      store.ID,
		ProductID:    product.ID,
		StoreName:    store.Name,
		ProductName:  product.Name,
		ProductPrice: product.Price,
		Quantity:     1,
	}
	type fields struct {
		CustomerID string
		PaymentID  string
		Items      map[string]Item
		Status     BasketStatus
	}
	type args struct {
		snapshot es.Snapshot
	}
	tests := map[string]struct {
		fields  fields
		args    args
		want    fields
		wantErr bool
	}{
		"V1": {
			fields: fields{},
			args: args{
				snapshot: &BasketV1{
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items: map[string]Item{
						product.ID: item,
					},
					Status: BasketIsOpen,
				},
			},
			want: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: item,
				},
				Status: BasketIsOpen,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b := &Basket{
				Aggregate:  es.NewMockAggregate(t),
				CustomerID: tt.fields.CustomerID,
				PaymentID:  tt.fields.PaymentID,
				Items:      tt.fields.Items,
				Status:     tt.fields.Status,
			}
			if err := b.ApplySnapshot(tt.args.snapshot); (err != nil) != tt.wantErr {
				t.Errorf("ApplySnapshot() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, b.CustomerID, tt.want.CustomerID)
			assert.Equal(t, b.PaymentID, tt.want.PaymentID)
			assert.Equal(t, b.Items, tt.want.Items)
			assert.Equal(t, b.Status, tt.want.Status)
		})
	}
}

func TestBasket_Cancel(t *testing.T) {
	type fields struct {
		CustomerID string
		PaymentID  string
		Items      map[string]Item
		Status     BasketStatus
	}
	tests := map[string]struct {
		fields  fields
		on      func(a *es.MockAggregate)
		want    ddd.Event
		wantErr bool
	}{
		"OpenBasket": {
			fields: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsOpen,
			},
			on: func(a *es.MockAggregate) {
				a.On("AddEvent", BasketCanceledEvent, &BasketCanceled{})
			},
			want: ddd.NewEvent(BasketCanceledEvent, &Basket{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsCanceled,
			}),
		},
		"CheckedOutBasket": {
			fields: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsCheckedOut,
			},
			wantErr: true,
		},
		"CanceledBasket": {
			fields: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsCanceled,
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			aggregate := es.NewMockAggregate(t)
			b := &Basket{
				Aggregate:  aggregate,
				CustomerID: tt.fields.CustomerID,
				PaymentID:  tt.fields.PaymentID,
				Items:      tt.fields.Items,
				Status:     tt.fields.Status,
			}
			if tt.on != nil {
				tt.on(aggregate)
			}

			got, err := b.Cancel()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cancel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				assert.Equal(t, tt.want.EventName(), got.EventName())
				assert.IsType(t, tt.want.Payload(), got.Payload())
				assert.Equal(t, tt.want.Metadata(), got.Metadata())
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func TestBasket_Checkout(t *testing.T) {
	store := &Store{
		ID:   "store-id",
		Name: "store-name",
	}
	product := &Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}
	item := Item{
		StoreID:      store.ID,
		ProductID:    product.ID,
		StoreName:    store.Name,
		ProductName:  product.Name,
		ProductPrice: product.Price,
		Quantity:     1,
	}

	type fields struct {
		CustomerID string
		PaymentID  string
		Items      map[string]Item
		Status     BasketStatus
	}
	type args struct {
		paymentID string
	}
	tests := map[string]struct {
		fields  fields
		args    args
		on      func(a *es.MockAggregate)
		want    ddd.Event
		wantErr bool
	}{
		"OpenBasket": {
			fields: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: item,
				},
				Status: BasketIsOpen,
			},
			args: args{paymentID: "payment-id"},
			on: func(a *es.MockAggregate) {
				a.On("AddEvent", BasketCheckedOutEvent, &BasketCheckedOut{
					PaymentID: "payment-id",
				})
			},
			want: ddd.NewEvent(BasketCheckedOutEvent, &Basket{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsCanceled,
			}),
		},
		"OpenBasket.NoItems": {
			fields: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsOpen,
			},
			args:    args{paymentID: "payment-id"},
			wantErr: true,
		},
		"CheckedOutBasket": {
			fields: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsCheckedOut,
			},
			args:    args{paymentID: "payment-id"},
			wantErr: true,
		},
		"CanceledBasket": {
			fields: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsCanceled,
			},
			args:    args{paymentID: "payment-id"},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			aggregate := es.NewMockAggregate(t)
			b := &Basket{
				Aggregate:  aggregate,
				CustomerID: tt.fields.CustomerID,
				PaymentID:  tt.fields.PaymentID,
				Items:      tt.fields.Items,
				Status:     tt.fields.Status,
			}
			if tt.on != nil {
				tt.on(aggregate)
			}

			// Act
			got, err := b.Checkout(tt.args.paymentID)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("Checkout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				assert.Equal(t, tt.want.EventName(), got.EventName())
				assert.IsType(t, tt.want.Payload(), got.Payload())
				assert.Equal(t, tt.want.Metadata(), got.Metadata())
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func TestBasket_RemoveItem(t *testing.T) {
	store := &Store{
		ID:   "store-id",
		Name: "store-name",
	}
	product := &Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}
	item := Item{
		StoreID:      store.ID,
		ProductID:    product.ID,
		StoreName:    store.Name,
		ProductName:  product.Name,
		ProductPrice: product.Price,
		Quantity:     10,
	}

	type fields struct {
		CustomerID string
		PaymentID  string
		Items      map[string]Item
		Status     BasketStatus
	}
	type args struct {
		product  *Product
		quantity int
	}
	tests := map[string]struct {
		fields  fields
		args    args
		on      func(a *es.MockAggregate)
		wantErr bool
	}{
		"OpenBasket": {
			fields: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: item,
				},
				Status: BasketIsOpen,
			},
			args: args{
				product:  product,
				quantity: 1,
			},
			on: func(a *es.MockAggregate) {
				a.On("AddEvent", BasketItemRemovedEvent, &BasketItemRemoved{
					ProductID: product.ID,
					Quantity:  1,
				})
			},
		},
		"OpenBasket.NoItems": {
			fields: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsOpen,
			},
			args: args{
				product:  product,
				quantity: 1,
			},
		},
		"CheckedOutBasket": {
			fields: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsCheckedOut,
			},
			args: args{
				product:  product,
				quantity: 1,
			},
			wantErr: true,
		},
		"CanceledBasket": {
			fields: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsCanceled,
			},
			args: args{
				product:  product,
				quantity: 1,
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			aggregate := es.NewMockAggregate(t)
			b := &Basket{
				Aggregate:  aggregate,
				CustomerID: tt.fields.CustomerID,
				PaymentID:  tt.fields.PaymentID,
				Items:      tt.fields.Items,
				Status:     tt.fields.Status,
			}
			if tt.on != nil {
				tt.on(aggregate)
			}

			if err := b.RemoveItem(tt.args.product, tt.args.quantity); (err != nil) != tt.wantErr {
				t.Errorf("RemoveItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBasket_Start(t *testing.T) {
	type fields struct {
		CustomerID string
		PaymentID  string
		Items      map[string]Item
		Status     BasketStatus
	}
	type args struct {
		customerID string
	}
	tests := map[string]struct {
		fields  fields
		args    args
		on      func(a *es.MockAggregate)
		want    ddd.Event
		wantErr bool
	}{
		"New": {
			fields: fields{},
			args:   args{customerID: "customer-id"},
			on: func(a *es.MockAggregate) {
				a.On("AddEvent", BasketStartedEvent, &BasketStarted{
					CustomerID: "customer-id",
				})
			},
			want: ddd.NewEvent(BasketStartedEvent, &Basket{
				CustomerID: "customer-id",
				PaymentID:  "",
				Items:      make(map[string]Item),
				Status:     BasketIsOpen,
			}),
		},
		"OpenBasket": {
			fields: fields{
				Status: BasketIsOpen,
			},
			args:    args{customerID: "customer-id"},
			wantErr: true,
		},
		"CheckedOutBasket": {
			fields: fields{
				Status: BasketIsCheckedOut,
			},
			args:    args{customerID: "customer-id"},
			wantErr: true,
		},
		"CanceledBasket": {
			fields: fields{
				Status: BasketIsCanceled,
			},
			args:    args{customerID: "customer-id"},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			aggregate := es.NewMockAggregate(t)
			b := &Basket{
				Aggregate:  aggregate,
				CustomerID: tt.fields.CustomerID,
				PaymentID:  tt.fields.PaymentID,
				Items:      tt.fields.Items,
				Status:     tt.fields.Status,
			}
			if tt.on != nil {
				tt.on(aggregate)
			}

			got, err := b.Start(tt.args.customerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				assert.Equal(t, tt.want.EventName(), got.EventName())
				assert.IsType(t, tt.want.Payload(), got.Payload())
				assert.Equal(t, tt.want.Metadata(), got.Metadata())
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func TestBasket_ToSnapshot(t *testing.T) {
	store := &Store{
		ID:   "store-id",
		Name: "store-name",
	}
	product := &Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}
	item := Item{
		StoreID:      store.ID,
		ProductID:    product.ID,
		StoreName:    store.Name,
		ProductName:  product.Name,
		ProductPrice: product.Price,
		Quantity:     10,
	}

	type fields struct {
		CustomerID string
		PaymentID  string
		Items      map[string]Item
		Status     BasketStatus
	}
	tests := map[string]struct {
		fields fields
		want   es.Snapshot
	}{
		"V1": {
			fields: fields{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: item,
				},
				Status: BasketIsOpen,
			},
			want: &BasketV1{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: item,
				},
				Status: BasketIsOpen,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b := &Basket{
				Aggregate:  es.NewMockAggregate(t),
				CustomerID: tt.fields.CustomerID,
				PaymentID:  tt.fields.PaymentID,
				Items:      tt.fields.Items,
				Status:     tt.fields.Status,
			}

			if got := b.ToSnapshot(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToSnapshot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBasket(t *testing.T) {
	type args struct {
		id string
	}
	tests := map[string]struct {
		args args
		want *Basket
	}{
		"Basket": {
			args: args{id: "basket-id"},
			want: &Basket{
				Aggregate: es.NewAggregate("basket-id", BasketAggregate),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := NewBasket(tt.args.id)

			assert.Equal(t, tt.want.ID(), got.ID())
			assert.Equal(t, tt.want.AggregateName(), got.AggregateName())
		})
	}
}
