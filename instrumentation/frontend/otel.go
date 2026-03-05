package instrumentation

import (
  "context"
  "fmt"
  "os"
  "time"

  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
  "go.opentelemetry.io/otel/metric"
  "go.opentelemetry.io/otel/sdk/resource"
  sdktrace "go.opentelemetry.io/otel/sdk/trace"
  semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
  "google.golang.org/grpc"
)

var (
  Meter metric.Meter
)

func Init(ctx context.Context) (func(context.Context) error, error) {
  endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
  if endpoint == "" {
    endpoint = "otel-gateway-opentelemetry-collector.observability.svc.cluster.local:4317"
  }

  traceExp, err := otlptracegrpc.New(ctx,
    otlptracegrpc.WithEndpoint(endpoint),
    otlptracegrpc.WithDialOption(grpc.WithInsecure()),
  )
  if err != nil {
    return nil, fmt.Errorf("trace exporter: %w", err)
  }

  res, err := resource.New(ctx,
    resource.WithAttributes(
      semconv.ServiceNameKey.String(getEnvDefault("OTEL_SERVICE_NAME", "frontend")),
      semconv.ServiceVersionKey.String(getEnvDefault("SERVICE_VERSION", "1.0.0")),
      attribute.String("deployment.environment", getEnvDefault("DEPLOY_ENV", "staging")),
    ),
  )
  if err != nil {
    return nil, fmt.Errorf("resource: %w", err)
  }

  tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(traceExp), sdktrace.WithResource(res))
  otel.SetTracerProvider(tp)

  Meter = otel.Meter("frontend-meter")

  return func(ctx context.Context) error {
    ctxShutdown, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    if err := tp.Shutdown(ctxShutdown); err != nil {
      return err
    }
    return nil
  }, nil
}

func getEnvDefault(k, d string) string {
  if v := os.Getenv(k); v != "" {
    return v
  }
  return d
}