package rpc

import (
	"context"

	"github.com/stackus/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func clientErrorUnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return errors.ReceiveGRPCError(invoker(ctx, method, req, reply, cc, opts...))
	}
}

func Dial(ctx context.Context, endpoint string) (conn *grpc.ClientConn, err error) {
	return grpc.DialContext(ctx, endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(),
			clientErrorUnaryInterceptor(),
		),
		// If there are streaming endpoints also add
		// grpc.WithStreamInterceptor(
		// 	otelgrpc.StreamClientInterceptor(),
		// ),
	)
}
