package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"eda-in-golang/internal/waiter"
)

var clients = flag.Int("clients", 5, "Number of clients to have running concurrently [min:1, max:25]")
var hostAddr = flag.String("host", "localhost:8080", "Sets the host address of the mallbots application")
var otlpAddr = flag.String("otlp", "http://collector:4317", "Sets the host address of the OpenTelemetry Collector")

func main() {
	log.SetFlags(log.Ltime)
	if err := run(); err != nil {
		log.Println(err.Error())
	}
	log.Println("busywork shutdown")
}

func run() error {
	flag.Parse()

	tp, err := initOpenTelemetry()
	if err != nil {
		return err
	}
	defer tp.Shutdown(context.Background())

	if *clients < 0 || *clients > 25 {
		*clients = 5
	}

	wait := waiter.New(waiter.CatchSignals())

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < *clients; i++ {
		jitter := time.Duration(rand.Int63n(int64(3 * time.Second)))
		interval := 8*time.Second + jitter
		wait.Add(newBusyworkClient(fmt.Sprintf("Client %d", i+1), interval).run)
	}

	return wait.Wait()
}

func initOpenTelemetry() (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithEndpoint(*otlpAddr))
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
