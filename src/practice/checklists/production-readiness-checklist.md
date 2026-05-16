# Production Readiness Checklist

Use this checklist **before launching any service to production**. A service that clears all items here has a fighting chance of surviving contact with real traffic. A service that skips sections will eventually fail in ways that are embarrassing and avoidable.

This is intentionally exhaustive. If your service is genuinely small (< 100 users, non-critical), mark items N/A with a reason rather than silently skipping.

---

## Observability

### Metrics

- [ ] Service exports RED metrics: **R**ate (requests/sec), **E**rrors (error rate), **D**uration (latency percentiles: p50, p95, p99)
- [ ] Downstream dependency calls are instrumented with the same RED metrics
- [ ] Business-level metrics are defined and exported (e.g. orders/min, payments processed/hr)
- [ ] Dashboards exist and are linked from the service README
- [ ] Dashboards show at least 7 days of history

### Logs

- [ ] Logs are structured (JSON) — not free-form text
- [ ] Each log entry includes: timestamp, severity, service name, trace_id, request_id
- [ ] Logs are shipped to a centralised log store (not just local disk)
- [ ] Error logs include enough context to diagnose without additional lookups
- [ ] PII is not logged (names, emails, tokens, card numbers)

### Tracing

- [ ] Distributed traces are emitted for all inbound requests
- [ ] Trace context is propagated to all downstream calls (HTTP headers, Kafka message headers)
- [ ] Traces are sampled appropriately (100% for errors, 1-10% for success paths)
- [ ] A trace for a representative request flow has been reviewed and makes sense

---

## Alerting

- [ ] SLO is defined: target availability and latency (e.g., "99.9% of requests succeed in < 500ms over 30 days")
- [ ] SLO-based alert is configured and tested (alert fires when error budget is burning fast)
- [ ] At least one on-call engineer is paged for SLO breach — not just email/Slack
- [ ] Alert has a linked runbook (not "investigate in Datadog")
- [ ] Runbook covers: symptoms, likely causes, diagnostic steps, remediation steps, escalation path
- [ ] Alert has been fired and resolved in staging to verify it works end-to-end
- [ ] No alert fires in staging that is not actionable (alert fatigue is a reliability risk)

---

## Reliability

### Failure Handling

- [ ] All external HTTP/gRPC calls have **timeouts** configured (not relying on default/infinite)
- [ ] Retries are implemented with **exponential backoff and jitter** (not fixed interval)
- [ ] Non-idempotent operations are NOT automatically retried without idempotency key
- [ ] **Circuit breaker** is in place for calls to dependencies with known reliability issues
- [ ] Service has a defined behaviour for each dependency failure (degrade gracefully vs fail fast)
- [ ] Graceful shutdown is implemented: service drains in-flight requests before exiting

### Resilience

- [ ] Service has been tested with its primary dependency unavailable (what happens?)
- [ ] Service has been tested at 2× expected peak load (does it degrade gracefully or crash?)
- [ ] Rate limiting is in place for inbound requests (prevent a single caller from overloading the service)
- [ ] Health check endpoint exists and is correctly wired to load balancer / orchestrator

---

## Performance

- [ ] p99 latency has been measured under realistic load (not just happy-path)
- [ ] Capacity headroom is understood: at what traffic level does the service need to scale out?
- [ ] Autoscaling is configured with validated thresholds (scale out before, not after, saturation)
- [ ] Database queries are reviewed: no N+1 queries, relevant indexes exist, slow query log reviewed
- [ ] Any caches have defined TTLs and cache invalidation strategy
- [ ] Memory and CPU usage have been profiled under load; no obvious leaks

---

## Security

- [ ] All inbound requests are authenticated (API key, JWT, mTLS — appropriate for the threat model)
- [ ] Authorisation is enforced at the service level (not assumed from the caller)
- [ ] All user-supplied inputs are validated before use (types, lengths, allowed values)
- [ ] SQL queries use parameterised statements or ORM — no string interpolation into queries
- [ ] Secrets are in a secrets manager (Vault, AWS Secrets Manager) — not in environment variables, not in code
- [ ] Dependencies are pinned to specific versions and scanned for known vulnerabilities (Snyk, Dependabot, govulncheck)
- [ ] TLS is enforced for all external communication; certificates are auto-rotating
- [ ] CORS policy is intentional and minimal (not `*` on a service with sensitive data)

---

## Data

- [ ] Database backups are enabled and tested (restore has been verified, not just assumed)
- [ ] If schema migration is needed: migration has been rehearsed, rollback plan exists
- [ ] Data retention policy is defined and enforced (how long is data kept, who deletes it)
- [ ] GDPR / data residency requirements have been reviewed and are satisfied
- [ ] Sensitive data at rest is encrypted

---

## Operations

### Deployment

- [ ] Deployment is automated (CI/CD pipeline — not manual SSH and `git pull`)
- [ ] Deployment is zero-downtime (rolling update, blue/green, or canary)
- [ ] Rollback procedure is documented and has been tested
- [ ] Config changes (feature flags, environment variables) do not require a full redeploy
- [ ] Deployment produces an artifact that is tagged with the git commit SHA

### Runbook

- [ ] Runbook is written and linked from the service repository
- [ ] Runbook covers: how to scale up/down, how to roll back, how to check health, how to access logs
- [ ] On-call team has been briefed on the service (not "read the docs on your first page")

---

## Launch

- [ ] Feature flag or traffic ramp plan exists for the initial rollout
- [ ] Rollout plan specifies: who monitors, what metrics indicate success, what triggers a rollback
- [ ] Stakeholders are notified of the launch plan (timing, scope, rollback decision criteria)
- [ ] Post-launch monitoring window is scheduled (at minimum: first hour, first 24 hours)

---

## Go / No-Go Decision

| Area | Status | Notes |
|---|---|---|
| Observability | ✅ / ⚠️ / ❌ | |
| Alerting | ✅ / ⚠️ / ❌ | |
| Reliability | ✅ / ⚠️ / ❌ | |
| Performance | ✅ / ⚠️ / ❌ | |
| Security | ✅ / ⚠️ / ❌ | |
| Data | ✅ / ⚠️ / ❌ | |
| Operations | ✅ / ⚠️ / ❌ | |
| Launch plan | ✅ / ⚠️ / ❌ | |

**Go criteria**: all areas ✅, or all ⚠️ with documented accepted risk. Any ❌ = No-Go.

---

## Common Failure Modes

| Skipped item | What goes wrong |
|---|---|
| No timeouts on external calls | One slow dependency causes request handler threads to exhaust; service becomes unresponsive |
| No structured logs | Incident happens, nobody can query logs; engineers read raw text files under pressure |
| No SLO-based alerting | Problems detected by users before engineering; trust erodes |
| Secrets in env vars | Secrets appear in crash dumps, log files, `ps aux` output; rotate and audit |
| Untested rollback | Rollback attempted under incident pressure; fails; incident is now also a deployment incident |
| No load test | Service works fine in staging with 10 users; crashes in production with 10,000 |
| Circuit breakers absent | One degraded dependency cascades to total service failure |
