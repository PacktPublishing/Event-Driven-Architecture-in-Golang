package errorsotel

import (
	"github.com/stackus/errors"
	"go.opentelemetry.io/otel/attribute"
)

func ErrAttrs(err error) []attribute.KeyValue {
	kvs := []attribute.KeyValue{
		attribute.String("error.message", err.Error()),
	}

	var typeCoder errors.TypeCoder
	if errors.As(err, &typeCoder) {
		kvs = append(kvs, attribute.String("error.type", typeCoder.TypeCode()))
	}

	var httpCoder errors.HTTPCoder
	if errors.As(err, &httpCoder) {
		kvs = append(kvs, attribute.Int64("error.http_code", int64(httpCoder.HTTPCode())))
	}

	var statusCoder errors.GRPCCoder
	if errors.As(err, &statusCoder) {
		kvs = append(kvs, attribute.Int64("error.grpc_code", int64(statusCoder.GRPCCode())))
	}

	return kvs
}
