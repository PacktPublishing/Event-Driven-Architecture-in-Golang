package grpc

import (
	"context"

	"google.golang.org/grpc"

	"eda-in-golang/search/internal/application"
	"eda-in-golang/search/searchpb"
)

type server struct {
	app application.Application
	searchpb.UnimplementedSearchServiceServer
}

func RegisterServer(_ context.Context, app application.Application, registrar grpc.ServiceRegistrar) error {
	searchpb.RegisterSearchServiceServer(registrar, server{app: app})
	return nil
}

func (s server) SearchOrders(ctx context.Context, request *searchpb.SearchOrdersRequest) (*searchpb.SearchOrdersResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s server) GetOrder(ctx context.Context, request *searchpb.GetOrderRequest) (*searchpb.GetOrderResponse, error) {
	// TODO implement me
	panic("implement me")
}
