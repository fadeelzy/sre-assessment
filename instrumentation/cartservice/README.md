Notes:
- Use OpenTelemetry .NET packages:
  OpenTelemetry
  OpenTelemetry.Exporter.OpenTelemetryProtocol
  OpenTelemetry.Extensions.Hosting
  OpenTelemetry.Instrumentation.AspNetCore
- Use env OTEL_EXPORTER_OTLP_ENDPOINT pointing to the gateway. For Kubernetes, set OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-gateway-opentelemetry-collector.observability.svc.cluster.local:4317
- Keep ELASTIC_APM_TOKEN as k8s secret for collector -> Elastic Cloud; services export to gateway.