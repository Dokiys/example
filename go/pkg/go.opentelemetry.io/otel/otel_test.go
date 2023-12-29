package otel

import (
	"context"
	"fmt"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func fibonacci(ctx context.Context, n int) int {
	// 设置全局 Tracer Provider
	provider := otel.GetTracerProvider()
	// 创建一个全局 tracer
	tracer := provider.Tracer("fibonacci-tracer")
	// 创建一个 span
	ctx, span := tracer.Start(ctx, "Fibonacci")
	defer span.End()
	// 设置 span 的属性
	span.SetAttributes(
		attribute.Int("input", n),
	)

	if n <= 1 {
		return n
	}

	return fibonacci(ctx, n-1) + fibonacci(ctx, n-2)
}

func TestOtelTrace(t *testing.T) {
	var ctx = context.Background()

	// init
	exporter, err := newExporter()
	if err != nil {
		panic("failed to initialize exporter: " + err.Error())
	}
	tp := newTraceProvider(exporter)
	defer func() { _ = tp.Shutdown(ctx) }()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		// https://opentelemetry.io/docs/instrumentation/go/manual/#propagators-and-context
		propagation.TraceContext{},
		// https://opentelemetry.io/docs/concepts/signals/baggage/
		propagation.Baggage{},
	))

	var n = 3
	fmt.Printf("Fibonacci(%d) = %d\n", n, fibonacci(ctx, n))
}

// https://opentelemetry.io/docs/instrumentation/go/exporters/
func newExporter() (sdktrace.SpanExporter, error) {
	return stdouttrace.New()
}

// https://opentelemetry.io/docs/instrumentation/go/manual/
func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("ExampleService"),
		),
	)
	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}
