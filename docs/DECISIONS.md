Date: 2026-03-05  
Authors: SRE / Infra Automation (repo: sre-assessment/infrastructure)  
Scope: Observability / metrics, logs, and alerting for services (Postgres, Redis, NGINX, Kibana) using Metricbeat/Elastic Agent and Kibana Rules.ADR-001: Observability PlatformDecision: Use Elastic Stack (Elasticsearch + Kibana) for metrics, logs, and alerting.Status: AcceptedRationale: Single platform for metrics, logs, APM; built-in integrations and alerting; works with Elastic Cloud Serverless.Consequences:
Use Elasticsearch endpoint + API keys for ingestion (Serverless restriction: no cloud.id/cloud.auth/username).
Centralized dashboards and rule engine; dependency on Elastic licensing/features.

Alternatives considered: Prometheus+Grafana (+Loki), Datadog. Rejected due to consolidation preference and existing Elastic usage.ADR-002: Data Collection AgentsDecision: Use Metricbeat (or Elastic Agent via Fleet) for metrics collection; Filebeat (or Elastic Agent) for logs where needed.Status: AcceptedRationale: Native Metricbeat modules for Postgres, Redis, NGINX produce standard fields (metricset.module) and integrate with Metricbeat dashboards; Elastic Agent simplifies Fleet management.Consequences:
Use modules.d/*.yml for Metricbeat-native or Fleet policy variables for Elastic Agent.
Ensure agents run where they can reach services' metrics endpoints and logs.

Alternatives: Custom exporters (Prometheus exporters) + Beats; rejected for extra orchestration and mapping.ADR-003: Integrations ImplementedDecision: Provide standardized config files under infrastructure/ for:
Postgres: JDBC/Metricbeat postgresql module (metrics: database, bgwriter, activity)
Redis: Metricbeat redis module (info, keyspace)
NGINX: Metricbeat nginx module (stubstatus) + Filebeat-style log inputs for access/error logs

Status: Implemented (files added)Rationale: Cover typical infra telemetry and expose fields useful for dashboards and alerts.Consequences:
Credentials must be injected securely (see ADR-005).
Path/host vars are parameterized for portability.

ADR-004: Alerting StrategyDecision: Implement rule-based alerts in Kibana for:
Service error rate threshold (apm.error_rate)
High P95 transaction duration (apm.transaction_duration)
Host CPU high (metrics.threshold on Metricbeat indices)

Status: Implemented (NDJSON saved-objects under alerting-rules/alerts.ndjson)Rationale: Use built-in rule types for APM and metrics for straightforward detection; schedule run every 1m with short lookbacks (5m) to surface regressions quickly.Consequences:
Rules reference connectors (actions) for notification. NDJSON contains a placeholder connector id (synthetic).
Environment-scoped: rules configured for staging. To monitor other environments, duplicate/update rules or use wildcard.

Alternatives: External alerting pipelines (Prometheus Alertmanager). Rejected to keep alerts inside Kibana and near data.ADR-005: Secrets & AuthenticationDecision: Never commit secrets. Use environment variables, Fleet policy variables, Kubernetes Secrets, or a secrets manager. Use Elasticsearch API keys for output authentication to Serverless.Status: Accepted and applied in configs (placeholders like ${POSTGRES_PASSWORD}).Rationale: Security best practices and Serverless requirements (API keys only).Consequences:
Deployment manifests must wire secrets into agent configs.
For rule creation via API, authenticate calls with a valid Kibana session or API key.

ADR-006: Connector / Notification HandlingDecision: Use Kibana Connectors (Slack, Email, Webhook) for alert actions; saved object NDJSON uses a generated UUID placeholder 7f3b2c9e-4a8d-4d1f-9b2e-1c6a8f0b2d5e which must be replaced by the real connector id in each environment before import.Status: Implemented (placeholder inserted)Rationale: Leverage Kibana Actions for flexible notifications and integrations.Consequences:
Users must create connectors in the target Kibana and update NDJSON before import, or create rules programmatically via the Rules API referencing real connector ids.

ADR-007: Environment Tagging & ScopeDecision: All integration configs add an environment: staging tag and alert rules target environment = "staging".Status: ImplementedRationale: Repository defaults and the user's confirmation that environment is staging.Consequences:
Alerts will only fire for telemetry tagged with staging. To cover multiple environments, use ENVIRONMENT_ALL/* or create environment-specific rules.

ADR-008: Deployment ModelDecision: Provide configs that support two deployment models:
Fleet-managed Elastic Agent (recommended for centralized management)
Standalone Metricbeat/Filebeat agents (for simple or constrained environments)

Status: Implemented (files and notes)Rationale: Support both managed (Fleet) and unmanaged setups to ease adoption.Consequences:
Fleet simplifies secret injection and policy updates; standalone requires manual config and secure secret handling.

ADR-009: Verification & Observability HealthchecksDecision: Add verification steps and test commands (module list/test, curl to Elasticsearch indices, sample searches) in repo documentation to confirm data ingestion and alerting.Status: Included in docs/instructionsRationale: Reduce time-to-detect misconfiguration; provide standard debug workflow.Consequences:
Operators must run these checks after deploy to ensure indices (metricbeat-, apm-) and dashboards/rules are present.

ADR-010: Saved-Object Format & CompatibilityDecision: Provide alerting rules as NDJSON saved-object lines, with note about version compatibility; recommend creating rules via the Rules API for Serverless/Cloud environments when possible.Status: ImplementedRationale: NDJSON import is quick for local testing; APIs are more robust for CI/CD.Consequences:
Before importing NDJSON, verify Kibana version compatibility. For automated provisioning prefer scripted API calls that authenticate and create connectors + rules atomically.

Action Items / Next StepsReplace connector placeholder uuid with real connector id(s) created in Kibana.Create Fleet policies (or Metricbeat DaemonSet) and wire secrets (Postgres/Redis passwords, Elasticsearch API key).Import NDJSON (after replacing connector ids) or create rules via the Rules API using CI/CD.Create dashboards or confirm Metricbeat dashboard imports via metricbeat setup.Add Kubernetes manifests or Fleet policy JSON to repo if desired for automated deployment.Add a runbook that describes how to respond to each alert (escalation, playbooks, runbook links).