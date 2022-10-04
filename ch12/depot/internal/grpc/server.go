package grpc

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

	"eda-in-golang/depot/depotpb"
	"eda-in-golang/depot/internal/application"
	"eda-in-golang/depot/internal/application/commands"
	"eda-in-golang/internal/errorsotel"
)

type server struct {
	app application.App
	depotpb.UnimplementedDepotServiceServer
}

var _ depotpb.DepotServiceServer = (*server)(nil)

func Register(app application.App, registrar grpc.ServiceRegistrar) error {
	depotpb.RegisterDepotServiceServer(registrar, server{app: app})
	return nil
}

func (s server) CreateShoppingList(ctx context.Context, request *depotpb.CreateShoppingListRequest) (*depotpb.CreateShoppingListResponse, error) {
	span := trace.SpanFromContext(ctx)

	id := uuid.New().String()

	span.SetAttributes(
		attribute.String("ShoppingListID", id),
		attribute.String("OrderID", request.GetOrderId()),
	)

	items := make([]commands.OrderItem, 0, len(request.GetItems()))
	for _, item := range request.GetItems() {
		items = append(items, s.itemToDomain(item))
	}

	err := s.app.CreateShoppingList(ctx, commands.CreateShoppingList{
		ID:      id,
		OrderID: request.GetOrderId(),
		Items:   items,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &depotpb.CreateShoppingListResponse{Id: id}, err
}

func (s server) CancelShoppingList(ctx context.Context, request *depotpb.CancelShoppingListRequest) (*depotpb.CancelShoppingListResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("ShoppingListID", request.GetId()),
	)

	err := s.app.CancelShoppingList(ctx, commands.CancelShoppingList{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &depotpb.CancelShoppingListResponse{}, err
}

func (s server) AssignShoppingList(ctx context.Context, request *depotpb.AssignShoppingListRequest) (*depotpb.AssignShoppingListResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("ShoppingListID", request.GetId()),
	)

	err := s.app.AssignShoppingList(ctx, commands.AssignShoppingList{
		ID:    request.GetId(),
		BotID: request.GetBotId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &depotpb.AssignShoppingListResponse{}, err
}

func (s server) CompleteShoppingList(ctx context.Context, request *depotpb.CompleteShoppingListRequest) (*depotpb.CompleteShoppingListResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("ShoppingListID", request.GetId()),
	)

	err := s.app.CompleteShoppingList(ctx, commands.CompleteShoppingList{ID: request.GetId()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &depotpb.CompleteShoppingListResponse{}, err
}

func (s server) itemToDomain(item *depotpb.OrderItem) commands.OrderItem {
	return commands.OrderItem{
		StoreID:   item.GetStoreId(),
		ProductID: item.GetProductId(),
		Quantity:  int(item.GetQuantity()),
	}
}
