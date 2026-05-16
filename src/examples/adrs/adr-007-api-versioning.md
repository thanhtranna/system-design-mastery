# ADR-007: URL Path Versioning for Public API

| Field | Value |
|---|---|
| **Status** | Accepted |
| **Date** | 2026-05-16 |
| **Author** | Thanh Tran |
| **Reviewers** | API Working Group, Mobile Team |
| **Review date** | 2027-05-16 |

---

## Context

We operate a public REST API consumed by third-party integrators, our mobile apps (iOS/Android), and web clients. We have 47 known external integrators as of this writing. Our API is served through Kong gateway in front of a Go backend.

We need to introduce breaking changes (v2) to three endpoints as part of our data model redesign. This requires a versioning strategy that:
- Allows v1 and v2 to coexist during the migration window (~6 months)
- Is understandable to external developers without reading extensive documentation
- Does not require changes to Kong routing logic for every version increment
- Works correctly with our existing API docs tooling (OpenAPI / Swagger)

**Breaking change scope**: field renames in response bodies, removal of deprecated fields, new required request fields.

---

## Decision

We will use **URL path versioning**: all API routes are prefixed with the major version number.

```
/v1/orders/{id}      → current stable API
/v2/orders/{id}      → new API with breaking changes
/v1/products         → unchanged; v1 and v2 may serve identical responses for some endpoints
```

Kong routing rule:
```
/v1/* → order-service:v1
/v2/* → order-service:v2 (or same service with version header passed downstream)
```

Version deprecation policy:
- v1 is supported for 12 months after v2 GA
- Deprecated versions return `Sunset` header indicating end-of-life date (RFC 8594)
- 90 days before end-of-life, v1 responses include `Deprecation: true` header

---

## Consequences

### Positive

- **Immediately visible in URLs**: developers, log files, browser DevTools — everyone can see which version is being called. No guessing which version a request is using.
- **Easy to test**: `curl https://api.example.com/v1/orders/123` vs `/v2/orders/123`. No custom header required.
- **OpenAPI tooling works natively**: OpenAPI specs per version are straightforward. Swagger UI handles `/v1` and `/v2` independently.
- **Caching is unambiguous**: CDN and browser caches treat `/v1/orders/123` and `/v2/orders/123` as distinct resources. No cache pollution between versions.

### Negative

- **URLs are not "pure" REST**: REST purists argue that a resource's URL should be stable across versions. Accepted trade-off — pragmatism wins over purity.
- **Clients must update base URL**: third-party integrators must change their base URL from `v1` to `v2`. This is the same cost as any versioning scheme — we cannot change breaking semantics without client changes.
- **Version number in URL can become stale**: if a client hardcodes `/v1/`, they may not notice the deprecation warning and miss the migration deadline. Mitigated by `Deprecation` and `Sunset` headers.

### Neutral

- **Internal services use v2 internally**: internal service-to-service calls also go through the versioned path. This is consistent — internal callers are treated like external callers.

---

## Alternatives Considered

### Header-based versioning (`Accept: application/vnd.api.v2+json`)

Version specified in the `Accept` header. URL is stable across versions.

**Rejected because**:
- Not visible in logs, proxies, or browser DevTools — debugging a version mismatch requires header inspection
- Curl/Postman/testing requires custom header setup for every request
- Our 47 external integrators range widely in technical sophistication; header-based versioning is consistently reported as confusing in API developer surveys
- Kong routing on header values is possible but more complex than path-based routing

### Content negotiation (custom `API-Version` header)

Similar to above but with a simpler custom header instead of `Accept` vendor media type.

**Rejected because**: same visibility problems as header-based. Additionally, custom headers are not cached correctly by some proxies. The header approach requires documentation to be clear about header name, format, and fallback behaviour (what happens if the header is absent?).

### No versioning (semver, breaking changes via deprecation warnings only)

Ship breaking changes as part of the same endpoint. Deprecate old fields with a warning response field.

**Rejected because**: with 47 external integrators on different release schedules, we cannot coordinate all clients to update simultaneously. A breaking change to `/orders/{id}` would break some integrators immediately. The 6-month migration window requires running v1 and v2 in parallel, which requires some form of versioning.

### GraphQL (eliminate REST versioning entirely)

GraphQL allows additive schema changes without versioning. Clients request exactly the fields they need.

**Rejected because**: our API surface is primarily CRUD, not graph-shaped data. GraphQL adds query parsing, schema introspection overhead, and a learning curve for our existing integrators. Out of scope for this decision cycle; may be revisited for v3.

---

## Open Questions

- What is the process for minor, non-breaking changes within a version? (Decision: non-breaking additions — new optional fields, new endpoints — are deployed to the current version without a version bump. Only breaking changes require a new version.)
- How do we notify integrators of deprecation? (Plan: email to registered developer accounts + `Deprecation` + `Sunset` headers + changelog entry. Final 30 days: weekly reminder emails.)
- What version does a request get if no version prefix is present? (Decision: 400 Bad Request with an error message directing the client to include the version prefix. No silent fallback to avoid version ambiguity.)

---

## References

- RFC 8594: The Sunset HTTP Header Field (deprecation signalling)
- Stripe API versioning documentation (reference for version sunset communication)
- ADR-003 (Modular Monolith) — the backend service structure this versioning strategy routes to
