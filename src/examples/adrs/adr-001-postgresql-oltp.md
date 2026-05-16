# ADR-001: Use PostgreSQL 16 as primary OLTP store for the Listings service

## Status

**Accepted** — 2026-02-14
**Authors**: Thanh Tran (Tech Lead, Listings Platform)
**Reviewers**: Platform team, DBA group, CTO
**Supersedes**: none
**Review date**: 2027-02-14 (annual)

## Context

The Listings service is the new core of Real Estate Company, replacing a 4-year-old PHP+MySQL monolith. It will own:

- ~5M active property listings, growing ~100K/month
- ~500K updates per day (price changes, status flips, photo additions)
- ~50M search/read requests per day with peak ~5K QPS
- Strict transactional needs for agent ownership transfers and listing status changes (cannot have two listings in `sold` state for the same property in any race condition)
- Geographic queries (within radius, bounding box)
- Some semi-structured data: custom listing attributes per property type

Stakeholders:

- **Backend team**: 8 engineers, deep Go experience, moderate SQL, no NoSQL production experience
- **DBA group**: small (2 people), highly skilled with MySQL/PostgreSQL, currently overloaded with the monolith
- **Business**: needs the new service live in Q3 2026; 5-year horizon for the chosen store
- **Compliance**: data residency required for SG and MY tenants

Constraints:

- **Hard**: data must stay within AWS Singapore + Malaysia regions; vendor-managed service preferred (small DBA group)
- **Soft**: total cost of ownership target <$8K/month at year-1 scale
- **Self-imposed**: avoid introducing a brand-new database technology in this service (we're already changing the runtime to Go and the deployment to Kubernetes — one bet at a time)

## Decision

**We will use Amazon RDS for PostgreSQL 16 as the primary OLTP store for the Listings service.** Geo queries will use PostGIS extension. Vector/embedding features (planned for semantic search in 2027) will use pgvector on the same instance until we hit its limits.

## Consequences

### Positive

- **Strong consistency by default**: serializable + repeatable read available; covers the agent-ownership race natively without application-side coordination.
- **Mature ecosystem**: ORMs, migration tools, observability, replication are battle-tested. Onboarding cost near zero for our SQL-fluent team.
- **PostGIS** handles geo queries (within radius, bounding box) with R-tree indexes; we don't need a separate geo store.
- **Mixed workload friendly**: JSONB for custom listing attributes lets us avoid premature schema decisions for new property types.
- **pgvector path**: when we add semantic search next year, no new datastore is needed for <10M vectors.
- **DBA familiarity**: existing skills transfer directly; minimal training.
- **RDS Multi-AZ**: failover-ready out of the box; satisfies our 99.95% availability target.
- **Cost**: ~$3,200/month for db.r6g.2xlarge Multi-AZ + storage; well under budget.

### Negative / Costs

- **Write scaling ceiling**: single primary writer caps us at ~5–10K writes/sec on this instance class. Our peak is ~500/sec today, so we have ~10–20× headroom. **Beyond that, we will need to shard application-side or migrate to a horizontally-scalable store.** We accept this; current data trajectory gives us 3–4 years before it becomes pressing.
- **Vacuum / bloat operations**: PostgreSQL's MVCC requires periodic VACUUM; can cause latency spikes on heavily-updated tables. Mitigation: partition the listings table by month for historical data; tune autovacuum aggressively for hot tables.
- **No native global distribution**: cross-region reads require logical replication, which we will set up read-only for analytics but not for serving.

### Neutral / Trade-offs accepted

- **Schema migrations require discipline**: PostgreSQL DDL is transactional but `ALTER TABLE` can lock; we adopt Strong Migrations conventions (concurrent index builds, no large in-place rewrites).
- **JSONB is not free**: indexed JSONB queries are slower than typed columns. We'll only use JSONB for truly variable attributes, not as a "schemaless" escape hatch.
- **pgvector at scale**: known to slow down past ~10M vectors. Acceptable for the planned semantic search v1; we'll re-evaluate.

## Alternatives Considered

### Option B: MySQL 8.4 (incumbent)

The PHP monolith already uses MySQL. Migrating to PostgreSQL adds a second flavor for the DBA team.

**Why rejected**:

- PostGIS is substantially more capable than MySQL's spatial features (our geo queries are non-trivial).
- pgvector has no MySQL equivalent of similar maturity; we'd need a second store for embeddings.
- PostgreSQL JSONB indexing (GIN) is more flexible than MySQL JSON for our custom-attributes use case.
- Migration cost is similar either way — we're rewriting the service.

### Option C: CockroachDB

Distributed SQL, would address the write-scaling concern proactively.

**Why rejected**:

- Solves a problem we don't have yet (we're at 500 writes/sec; the service can grow 10–20× before we hit PostgreSQL limits).
- Operational complexity is higher; we don't have CockroachDB experience on the team.
- Cost is meaningfully higher at our scale.
- Compliance: at the time of decision, CockroachDB Cloud's SG region was limited; AWS managed had no offering.
- "When you need it, you'll know." Until then, distributed-SQL premium is wasted.

### Option D: DynamoDB

AWS-native, infinite scale, fully managed.

**Why rejected**:

- Our access patterns include ad-hoc queries (admin tools, reports) that don't fit DynamoDB's key-based model.
- Strong consistency only for single-item reads; multi-item transactions are limited and expensive.
- Geo queries on DynamoDB require additional indexing (geohash) and are not first-class.
- Team has zero DynamoDB experience; the Go SDK + data modeling shift is non-trivial.
- Cost at our access pattern (high read fanout, complex filters) projects 3–4× higher than PostgreSQL.

### Option E: MongoDB

Document model, would fit our "custom attributes per listing type" use case.

**Why rejected**:

- Lost the historical reputation issue, but lacks our consistency confidence for ownership transfers (transactions exist but are more expensive than PG).
- No PostGIS-quality geo + no pgvector equivalent.
- We don't actually need document model — JSONB on top of typed columns is sufficient.
- Yet another technology in the stack.

## Open Questions

- **Replication topology for read scaling**: do we run a read replica from day 1, or only when CPU on primary trends to 70%? _Decision deferred to operational review at month 3._
- **Partitioning strategy for `listings.history`**: monthly partitions seem right but we lack production data to confirm. _Will revisit after 3 months of telemetry._
- **pgvector capacity at our growth rate**: when do we cross the threshold where a dedicated vector store (Pinecone/Milvus) becomes justified? _Build the embedding pipeline first; revisit when index size > 5GB or query latency exceeds 200ms p95._

## References

- AWS RDS PostgreSQL pricing calculator (SG region, db.r6g.2xlarge Multi-AZ): https://calculator.aws/
- PostGIS performance comparison report (internal benchmark, attached: `bench-2026-01-31.pdf`)
- pgvector benchmarks at our embedding dimensions (1536, OpenAI): https://github.com/pgvector/pgvector#performance
- Strong Migrations guide: https://github.com/ankane/strong_migrations
- Prior incident: "MongoDB ownership-race incident" (2024-08-12) — historical context on why strong consistency matters here.

## Decision Record History

| Date       | Status   | Notes                                                             |
| ---------- | -------- | ----------------------------------------------------------------- |
| 2026-01-15 | Proposed | Initial draft, circulated to Platform team                        |
| 2026-01-28 | Revised  | Added cost analysis, expanded Option C reasoning per CTO feedback |
| 2026-02-14 | Accepted | Approved at architecture review meeting                           |
