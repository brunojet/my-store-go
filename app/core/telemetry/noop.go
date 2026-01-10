package telemetry

import (
	"context"
	"time"
)

type NoopProvider struct{}

type noopTracer struct{}

type noopSpan struct{}

type noopMetrics struct{}

type noopLogger struct{}

func (NoopProvider) Tracer() Tracer   { return noopTracer{} }
func (NoopProvider) Metrics() Metrics { return noopMetrics{} }
func (NoopProvider) Logger() Logger   { return noopLogger{} }

func (noopTracer) Start(ctx context.Context, _ string, _ ...Field) (context.Context, Span) {
	return ctx, noopSpan{}
}

func (noopSpan) End(_ error, _ ...Field) {}

func (noopMetrics) Inc(_ string, _ int64, _ ...Field)                     {}
func (noopMetrics) ObserveDuration(_ string, _ time.Duration, _ ...Field) {}

func (noopLogger) Info(_ context.Context, _ string, _ ...Field)           {}
func (noopLogger) Error(_ context.Context, _ string, _ error, _ ...Field) {}
