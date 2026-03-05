# SRE Assessment - Instrumentation Template

This repo is a template to instrument frontend (Go), cartservice (.NET), and paymentservice (Node.js) and deploy OpenTelemetry Collector Gateway + Agents that forward to Elastic Cloud.

Important: DO NOT commit secrets. Store your Elastic ApiKey in a Kubernetes Secret named `elastic-apm-token` in namespace `observability`.

Quick start (WSL):

1) Create namespace & secret
   kubectl create namespace observability || true
   kubectl create secret generic elastic-apm-token --from-literal=token="REPLACE_WITH_API_KEY" -n observability

2) Deploy OTel Collector (Gateway + Agent) via Helm
   helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
   helm repo update

   helm install otel-gateway open-telemetry/opentelemetry-collector -f otel-collector/values-gateway.yaml -n observability
   helm install otel-agent open-telemetry/opentelemetry-collector -f otel-collector/values-agent.yaml -n observability

3) Instrument services
   - frontend (Go): copy instrumentation/frontend/otel.go into your frontend project and call Init(ctx). Create metrics counters and wrap handlers with otelhttp.NewHandler. (See instrumentation/frontend/README.md)
   - cartservice (.NET): use the stub in instrumentation/cartservice, install OpenTelemetry NuGet packages and configure OTLP exporter to gateway.
   - paymentservice (Node): copy instrumentation/paymentservice/* files, run `npm install`, require('./otel') at app startup.

4) Build & Deploy service images (optional local Docker Desktop)
   - Build: docker build -t local/frontend:latest ./frontend
            docker build -t local/cartservice:latest ./cartservice
            docker build -t local/paymentservice:latest ./paymentservice
   - Deploy minimal k8s deployments (see your own manifests); ensure env:
       OTEL_EXPORTER_OTLP_ENDPOINT=otel-gateway-opentelemetry-collector.observability.svc.cluster.local:4317

5) Generate traffic & verify traces & metrics in Elastic APM (Elastic Cloud)
   Example:
     curl -X POST http://<frontend-svc>:8080/cart -H "X-User-ID: user123" -d '{"product_id":1}'