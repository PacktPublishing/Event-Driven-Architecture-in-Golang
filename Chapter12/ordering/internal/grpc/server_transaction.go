package grpc

import (
	"context"
	"database/sql"

	"google.golang.org/grpc"

	"eda-in-golang/internal/di"
	"eda-in-golang/ordering/internal/application"
	"eda-in-golang/ordering/internal/constants"
	"eda-in-golang/ordering/orderingpb"
)

type serverTx struct {
	c di.Container
	orderingpb.UnimplementedOrderingServiceServer
}

var _ orderingpb.OrderingServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	orderingpb.RegisterOrderingServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) CreateOrder(ctx context.Context, request *orderingpb.CreateOrderRequest) (resp *orderingpb.CreateOrderResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CreateOrder(ctx, request)
}

func (s serverTx) GetOrder(ctx context.Context, request *orderingpb.GetOrderRequest) (resp *orderingpb.GetOrderResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetOrder(ctx, request)
}

func (s serverTx) CancelOrder(ctx context.Context, request *orderingpb.CancelOrderRequest) (resp *orderingpb.CancelOrderResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CancelOrder(ctx, request)
}

func (s serverTx) ReadyOrder(ctx context.Context, request *orderingpb.ReadyOrderRequest) (resp *orderingpb.ReadyOrderResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ReadyOrder(ctx, request)
}

func (s serverTx) CompleteOrder(ctx context.Context, request *orderingpb.CompleteOrderRequest) (resp *orderingpb.CompleteOrderResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CompleteOrder(ctx, request)
}

func (s serverTx) closeTx(tx *sql.Tx, err error) error {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if err != nil {
		_ = tx.Rollback()
		return err
	} else {
		return tx.Commit()
	}
}
