package trace

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func NewExporter() (sdktrace.SpanExporter, error) {
	return otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithInsecure(),
	)
}
