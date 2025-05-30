package server

import (
	"context"
	"errors"
	"fmt"
	"log"

	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var tracer = otel.Tracer("FOXDEN")

func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Start a span for the incoming request
		spanName := c.FullPath() // or c.Request.URL.Path
		ctx, span := tracer.Start(ctx, spanName)
		defer span.End()

		// Add useful attributes
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("gin.route", c.FullPath()),
		)

		// Inject the context with span into the request
		c.Request = c.Request.WithContext(ctx)

		// Continue to the next middleware/handler
		c.Next()

		// Optionally record the status
		statusCode := c.Writer.Status()
		span.SetAttributes(attribute.Int("http.status_code", statusCode))
		if statusCode >= 500 {
			span.RecordError(c.Errors.Last())
		}
	}
}

func InitTracer() (*sdktrace.TracerProvider, error) {
	var spanProcessors []sdktrace.SpanProcessor

	// Resource (e.g., service name)
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(srvConfig.Config.OpenTelemetry.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	if srvConfig.Config.OpenTelemetry.JaegerEndpoint != "" {
		exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(srvConfig.Config.OpenTelemetry.JaegerEndpoint)))
		if err != nil {
			return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
		}
		spanProcessors = append(spanProcessors, sdktrace.NewBatchSpanProcessor(exp))
		log.Println("Tracing: Jaeger enabled")
	}

	if srvConfig.Config.OpenTelemetry.OTLPEndpoint != "" {
		exp, err := otlptracehttp.New(context.Background(),
			otlptracehttp.WithEndpoint(srvConfig.Config.OpenTelemetry.OTLPEndpoint),
			otlptracehttp.WithInsecure(), // for local/dev testing
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
		}
		spanProcessors = append(spanProcessors, sdktrace.NewBatchSpanProcessor(exp))
		log.Println("Tracing: OTLP enabled")
	}

	if srvConfig.Config.OpenTelemetry.EnableStdout {
		exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout exporter: %w", err)
		}
		spanProcessors = append(spanProcessors, sdktrace.NewSimpleSpanProcessor(exp))
		log.Println("Tracing: Stdout enabled")
	}

	if len(spanProcessors) == 0 {
		msg := "Tracing: No exporters configured, tracing is disabled"
		return nil, errors.New(msg)
	}

	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(res),
	}

	for _, sp := range spanProcessors {
		opts = append(opts, sdktrace.WithSpanProcessor(sp))
	}

	tp := sdktrace.NewTracerProvider(opts...)

	otel.SetTracerProvider(tp)
	log.Println("Tracing: TracerProvider initialized")

	return tp, nil
}
