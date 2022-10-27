package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"eda-in-golang/stores/storespb"

	"eda-in-golang/stores/internal/application"
	"eda-in-golang/stores/internal/application/commands"
	"eda-in-golang/stores/internal/application/queries"
	"eda-in-golang/stores/internal/domain"
)

type server struct {
	app application.App
	storespb.UnimplementedStoresServiceServer
}

var _ storespb.StoresServiceServer = (*server)(nil)

func RegisterServer(_ context.Context, app application.App, registrar grpc.ServiceRegistrar) error {
	storespb.RegisterStoresServiceServer(registrar, server{app: app})
	return nil
}

func (s server) CreateStore(ctx context.Context, request *storespb.CreateStoreRequest) (*storespb.CreateStoreResponse, error) {
	storeID := uuid.New().String()

	err := s.app.CreateStore(ctx, commands.CreateStore{
		ID:       storeID,
		Name:     request.GetName(),
		Location: request.GetLocation(),
	})
	if err != nil {
		return nil, err
	}

	return &storespb.CreateStoreResponse{
		Id: storeID,
	}, nil
}

func (s server) EnableParticipation(ctx context.Context, request *storespb.EnableParticipationRequest) (*storespb.EnableParticipationResponse, error) {
	err := s.app.EnableParticipation(ctx, commands.EnableParticipation{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &storespb.EnableParticipationResponse{}, nil
}

func (s server) DisableParticipation(ctx context.Context, request *storespb.DisableParticipationRequest) (*storespb.DisableParticipationResponse, error) {
	err := s.app.DisableParticipation(ctx, commands.DisableParticipation{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &storespb.DisableParticipationResponse{}, nil
}

func (s server) RebrandStore(ctx context.Context, request *storespb.RebrandStoreRequest) (*storespb.RebrandStoreResponse, error) {
	err := s.app.RebrandStore(ctx, commands.RebrandStore{
		ID:   request.GetId(),
		Name: request.GetName(),
	})

	return &storespb.RebrandStoreResponse{}, err
}

func (s server) GetStore(ctx context.Context, request *storespb.GetStoreRequest) (*storespb.GetStoreResponse, error) {
	store, err := s.app.GetStore(ctx, queries.GetStore{ID: request.GetId()})
	if err != nil {
		return nil, err
	}

	return &storespb.GetStoreResponse{Store: s.storeFromDomain(store)}, nil
}

func (s server) GetStores(ctx context.Context, request *storespb.GetStoresRequest) (*storespb.GetStoresResponse, error) {
	stores, err := s.app.GetStores(ctx, queries.GetStores{})
	if err != nil {
		return nil, err
	}

	protoStores := []*storespb.Store{}
	for _, store := range stores {
		protoStores = append(protoStores, s.storeFromDomain(store))
	}

	return &storespb.GetStoresResponse{
		Stores: protoStores,
	}, nil
}

func (s server) GetParticipatingStores(ctx context.Context, request *storespb.GetParticipatingStoresRequest) (*storespb.GetParticipatingStoresResponse, error) {
	stores, err := s.app.GetParticipatingStores(ctx, queries.GetParticipatingStores{})
	if err != nil {
		return nil, err
	}

	protoStores := []*storespb.Store{}
	for _, store := range stores {
		protoStores = append(protoStores, s.storeFromDomain(store))
	}

	return &storespb.GetParticipatingStoresResponse{
		Stores: protoStores,
	}, nil
}

func (s server) AddProduct(ctx context.Context, request *storespb.AddProductRequest) (*storespb.AddProductResponse, error) {
	id := uuid.New().String()
	err := s.app.AddProduct(ctx, commands.AddProduct{
		ID:          id,
		StoreID:     request.GetStoreId(),
		Name:        request.GetName(),
		Description: request.GetDescription(),
		SKU:         request.GetSku(),
		Price:       request.GetPrice(),
	})
	if err != nil {
		return nil, err
	}

	return &storespb.AddProductResponse{Id: id}, nil
}

func (s server) RebrandProduct(ctx context.Context, request *storespb.RebrandProductRequest) (*storespb.RebrandProductResponse, error) {
	err := s.app.RebrandProduct(ctx, commands.RebrandProduct{
		ID:          request.GetId(),
		Name:        request.GetName(),
		Description: request.GetDescription(),
	})
	return &storespb.RebrandProductResponse{}, err
}

func (s server) IncreaseProductPrice(ctx context.Context, request *storespb.IncreaseProductPriceRequest) (*storespb.IncreaseProductPriceResponse, error) {
	err := s.app.IncreaseProductPrice(ctx, commands.IncreaseProductPrice{
		ID:    request.GetId(),
		Price: request.GetPrice(),
	})
	return &storespb.IncreaseProductPriceResponse{}, err
}

func (s server) DecreaseProductPrice(ctx context.Context, request *storespb.DecreaseProductPriceRequest) (*storespb.DecreaseProductPriceResponse, error) {
	err := s.app.DecreaseProductPrice(ctx, commands.DecreaseProductPrice{
		ID:    request.GetId(),
		Price: request.GetPrice(),
	})
	return &storespb.DecreaseProductPriceResponse{}, err
}

func (s server) RemoveProduct(ctx context.Context, request *storespb.RemoveProductRequest) (*storespb.RemoveProductResponse, error) {
	err := s.app.RemoveProduct(ctx, commands.RemoveProduct{
		ID: request.GetId(),
	})

	return &storespb.RemoveProductResponse{}, err
}

func (s server) GetProduct(ctx context.Context, request *storespb.GetProductRequest) (*storespb.GetProductResponse, error) {
	product, err := s.app.GetProduct(ctx, queries.GetProduct{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &storespb.GetProductResponse{Product: s.productFromDomain(product)}, nil
}

func (s server) GetCatalog(ctx context.Context, request *storespb.GetCatalogRequest) (*storespb.GetCatalogResponse, error) {
	products, err := s.app.GetCatalog(ctx, queries.GetCatalog{StoreID: request.GetStoreId()})
	if err != nil {
		return nil, err
	}

	protoProducts := make([]*storespb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = s.productFromDomain(product)
	}

	return &storespb.GetCatalogResponse{
		Products: protoProducts,
	}, nil
}

func (s server) storeFromDomain(store *domain.MallStore) *storespb.Store {
	return &storespb.Store{
		Id:            store.ID,
		Name:          store.Name,
		Location:      store.Location,
		Participating: store.Participating,
	}
}

func (s server) productFromDomain(product *domain.CatalogProduct) *storespb.Product {
	return &storespb.Product{
		Id:          product.ID,
		StoreId:     product.StoreID,
		Name:        product.Name,
		Description: product.Description,
		Sku:         product.SKU,
		Price:       product.Price,
	}
}
