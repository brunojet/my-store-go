package telemetry

import (
	"context"
	"time"
)

type Field struct {
	Key   string
	Value any
}

type Span interface {
	End(err error, fields ...Field)
}

type Tracer interface {
	Start(ctx context.Context, name string, fields ...Field) (context.Context, Span)
}

type Metrics interface {
	Inc(name string, value int64, fields ...Field)
	ObserveDuration(name string, d time.Duration, fields ...Field)
}

type Logger interface {
	Info(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, err error, fields ...Field)
}

type Provider interface {
	Tracer() Tracer
	Metrics() Metrics
	Logger() Logger
}
