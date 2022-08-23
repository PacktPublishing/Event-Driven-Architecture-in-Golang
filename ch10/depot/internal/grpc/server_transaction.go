package grpc

import (
	"context"
	"database/sql"

	"google.golang.org/grpc"

	"eda-in-golang/depot/depotpb"
	"eda-in-golang/depot/internal/application"
	"eda-in-golang/internal/di"
)

type serverTx struct {
	c di.Container
	depotpb.UnimplementedDepotServiceServer
}

var _ depotpb.DepotServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	depotpb.RegisterDepotServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) CreateShoppingList(ctx context.Context, request *depotpb.CreateShoppingListRequest) (resp *depotpb.CreateShoppingListResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.CreateShoppingList(ctx, request)
}

func (s serverTx) CancelShoppingList(ctx context.Context, request *depotpb.CancelShoppingListRequest) (resp *depotpb.CancelShoppingListResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.CancelShoppingList(ctx, request)
}

func (s serverTx) AssignShoppingList(ctx context.Context, request *depotpb.AssignShoppingListRequest) (resp *depotpb.AssignShoppingListResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.AssignShoppingList(ctx, request)
}

func (s serverTx) CompleteShoppingList(ctx context.Context, request *depotpb.CompleteShoppingListRequest) (resp *depotpb.CompleteShoppingListResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.CompleteShoppingList(ctx, request)
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
