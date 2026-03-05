
<script src="https://cdn.jsdelivr.net/npm/@elastic/apm-rum@5/dist/bundles/elastic-apm-rum.umd.min.js"></script>
<script>
  var apm = elasticApm.init({
    serviceName: 'frontend',
    serverUrl: 'https://my-observability-project-d46b5c.ingest.us-central1.gcp.elastic.cloud:443',
    serviceVersion: '1.0.0',
    environment: 'staging'
    // For public token handling, serve a token or configure APM Server CORS.
  });
</script>

