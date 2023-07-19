package otel

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// SpanFromContext takes a context and builds a span from it.
// If the context already contains a span, it will be used.
// If the context does not contain a span, a new one will be created.
func SpanFromContext(ctx context.Context, instrumentationName, spanName string) (context.Context, trace.Span) {
	return trace.SpanFromContext(ctx).TracerProvider().Tracer(instrumentationName).Start(ctx, spanName)
}
