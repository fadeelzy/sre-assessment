Quick notes:
- Place otel.go into package `instrumentation` and import it from main.
- Ensure envs:
  OTEL_EXPORTER_OTLP_ENDPOINT (default set to gateway FQDN)
  OTEL_SERVICE_NAME=frontend
  SERVICE_VERSION=1.0.0
  DEPLOY_ENV=staging
  ELASTIC_APM_TOKEN secret in k8s
- Create counters in main after Init:
  meter := otel.GetMeterProvider().Meter("frontend")
  cartAddCounter, _ := meter.Int64Counter("cart_additions_total")
- Wrap handlers with otelhttp.NewHandler and start spans with otel.Tracer("frontend").