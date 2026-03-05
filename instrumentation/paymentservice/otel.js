"use strict";

const { NodeSDK } = require('@opentelemetry/sdk-node');
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-grpc');
const { OTLPMetricExporter } = require('@opentelemetry/exporter-metrics-otlp-grpc');
const api = require('@opentelemetry/api');

const endpoint = process.env.OTEL_EXPORTER_OTLP_ENDPOINT || 'otel-gateway-opentelemetry-collector.observability.svc.cluster.local:4317';

const traceExporter = new OTLPTraceExporter({ url: endpoint });
const metricExporter = new OTLPMetricExporter({ url: endpoint });

const sdk = new NodeSDK({
  traceExporter,
  metricExporter,
  instrumentations: [getNodeAutoInstrumentations()],
});

sdk.start().then(() => {
  console.log('OTEL SDK started');
}).catch((err) => {
  console.error('Error starting OTEL SDK', err);
});

// Example manual span usage in your route handler
function startPaymentSpan(ctx, amount, userId) {
  const tracer = api.trace.getTracer('paymentservice');
  const span = tracer.startSpan('validate-payment', {
    attributes: { 'order.total': amount, 'user.id': userId },
  });
  return span;
}

module.exports = { startPaymentSpan };