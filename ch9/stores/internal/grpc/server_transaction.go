package grpc

import (
	"context"
	"database/sql"

	"google.golang.org/grpc"

	"eda-in-golang/internal/di"
	"eda-in-golang/stores/internal/application"
	"eda-in-golang/stores/storespb"
)

type serverTx struct {
	c di.Container
	storespb.UnimplementedStoresServiceServer
}

var _ storespb.StoresServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	storespb.RegisterStoresServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) CreateStore(ctx context.Context, request *storespb.CreateStoreRequest) (resp *storespb.CreateStoreResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.CreateStore(ctx, request)
}

func (s serverTx) EnableParticipation(ctx context.Context, request *storespb.EnableParticipationRequest) (resp *storespb.EnableParticipationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.EnableParticipation(ctx, request)
}

func (s serverTx) DisableParticipation(ctx context.Context, request *storespb.DisableParticipationRequest) (resp *storespb.DisableParticipationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.DisableParticipation(ctx, request)
}

func (s serverTx) RebrandStore(ctx context.Context, request *storespb.RebrandStoreRequest) (resp *storespb.RebrandStoreResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.RebrandStore(ctx, request)
}

func (s serverTx) GetStore(ctx context.Context, request *storespb.GetStoreRequest) (resp *storespb.GetStoreResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetStore(ctx, request)
}

func (s serverTx) GetStores(ctx context.Context, request *storespb.GetStoresRequest) (resp *storespb.GetStoresResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetStores(ctx, request)
}

func (s serverTx) GetParticipatingStores(ctx context.Context, request *storespb.GetParticipatingStoresRequest) (resp *storespb.GetParticipatingStoresResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetParticipatingStores(ctx, request)
}

func (s serverTx) AddProduct(ctx context.Context, request *storespb.AddProductRequest) (resp *storespb.AddProductResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.AddProduct(ctx, request)
}

func (s serverTx) RebrandProduct(ctx context.Context, request *storespb.RebrandProductRequest) (resp *storespb.RebrandProductResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.RebrandProduct(ctx, request)
}

func (s serverTx) IncreaseProductPrice(ctx context.Context, request *storespb.IncreaseProductPriceRequest) (resp *storespb.IncreaseProductPriceResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.IncreaseProductPrice(ctx, request)
}

func (s serverTx) DecreaseProductPrice(ctx context.Context, request *storespb.DecreaseProductPriceRequest) (resp *storespb.DecreaseProductPriceResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.DecreaseProductPrice(ctx, request)
}

func (s serverTx) RemoveProduct(ctx context.Context, request *storespb.RemoveProductRequest) (resp *storespb.RemoveProductResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.RemoveProduct(ctx, request)
}

func (s serverTx) GetProduct(ctx context.Context, request *storespb.GetProductRequest) (resp *storespb.GetProductResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetProduct(ctx, request)
}

func (s serverTx) GetCatalog(ctx context.Context, request *storespb.GetCatalogRequest) (resp *storespb.GetCatalogResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetCatalog(ctx, request)
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
