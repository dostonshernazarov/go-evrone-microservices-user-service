package otlp

import (
	"context"

	otelpkg "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Span interface {
	trace.Span
	EndError(err error, options ...trace.SpanEndOption)
	Error(err error)
}

func Start(ctx context.Context, name, spanName string) (context.Context, Span) {
	ctx, _span := otelpkg.Tracer(name).Start(ctx, spanName)
	return ctx, &span{span: _span}
}

type span struct {
	span trace.Span
}

func (s *span) End(options ...trace.SpanEndOption) {
	s.span.End(options...)
}

func (s *span) EndError(err error, options ...trace.SpanEndOption) {
	s.Error(err)
	s.span.End(options...)
}

func (s *span) AddEvent(name string, options ...trace.EventOption) {
	s.span.AddEvent(name, options...)
}

func (s *span) IsRecording() bool {
	return s.span.IsRecording()
}

func (s *span) RecordError(err error, options ...trace.EventOption) {
	s.span.RecordError(err, options...)
}

func (s *span) SpanContext() trace.SpanContext {
	return s.span.SpanContext()
}

func (s *span) SetStatus(code codes.Code, description string) {
	s.span.SetStatus(code, description)
}

func (s *span) SetName(name string) {
	s.span.SetName(name)
}

func (s *span) SetAttributes(kv ...attribute.KeyValue) {
	s.span.SetAttributes(kv...)
}

func (s *span) TracerProvider() trace.TracerProvider {
	return s.span.TracerProvider()
}

func (s *span) Error(err error) {
	if err != nil {
		s.span.SetStatus(codes.Error, err.Error())
	}
}

// RestoreTraceContext function forms context and span from trace_id and span_id
//
// span_id and trace_id should both be strings in hex format.
//
// Although this function returns both context and span it is required to call Start method to start tracing
// WARNING: if error IS NOT NIL, then context and span are BOTH NIL.
func RestoreTraceContext(traceIdStr, spanIdStr string) (context.Context, trace.Span, error) {
	spanId, err := trace.SpanIDFromHex(spanIdStr)
	if err != nil {
		return nil, nil, err
	}

	traceId, err := trace.TraceIDFromHex(traceIdStr)
	if err != nil {
		return nil, nil, err
	}

	ctx := trace.ContextWithRemoteSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceId,
		SpanID:  spanId,
	}))

	return ctx, trace.SpanFromContext(ctx), nil
}
