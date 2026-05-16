# ADR-004: Use Meilisearch as the search engine for PropertyHub listings

## Status

**Accepted** — 2026-03-18
**Authors**: Thanh Tran
**Reviewers**: Search team, Listings team, Platform team
**Supersedes**: none
**Review date**: 2027-03-18

## Context

PropertyHub needs full-text + faceted search over property listings:

- ~5M active listings (growing ~100K/month)
- Search fields: title, description, location (text + geo coordinates), agent name
- Facets: property type, price range (bucketed), bedrooms, listing status, neighborhood
- Typo tolerance, prefix matching, instant-search (search-as-you-type) on mobile and web
- Geo queries (within radius)
- Multi-language: English, Vietnamese (significant), Chinese (smaller volume)

Performance requirements:

- p95 query latency < 50ms (instant-search on mobile)
- p99 query latency < 150ms
- Index updates within 5 seconds of source change
- Peak query rate: ~5K QPS

Operational constraints:

- Search team: 2 backend engineers, no prior Elasticsearch experience
- DevOps team: 2 engineers, limited bandwidth
- Self-hosted preferred for cost reasons (managed search engines are expensive at our scale)

## Decision

**We will use Meilisearch v1.7 (self-hosted on AWS) as the search engine for product listings.** Indexing will be driven by the Outbox pattern (see ADR-002): listing events flow to Kafka, a small Go consumer transforms them into Meilisearch documents and pushes them via Meilisearch HTTP API.

## Consequences

### Positive

- **Excellent out-of-the-box developer experience.** Index creation, schema, ranking all configured via REST or simple TOML. Steeper learning curve avoided.
- **Built-in typo tolerance** and prefix search. No custom analyzer chains needed.
- **Latency profile matches our SLO.** Benchmarks (below) show p99 < 80ms on our dataset.
- **Resource-efficient.** Single-server deployments handle our scale comfortably. Rust-based, predictable memory usage.
- **Vietnamese tokenization** support is acceptable for our needs (verified during benchmark).
- **Active development** and growing ecosystem (admin UI, official SDKs for Go, JS).
- **Lower operational cost** vs OpenSearch: simpler deployment, less tuning, smaller instance footprint.

### Negative / Costs

- **Less mature than OpenSearch/Elasticsearch.** Some advanced features missing or limited:
  - No native query DSL as expressive as ES (we don't need this; our queries are straightforward)
  - Less sophisticated relevance scoring tunables
  - Aggregations and analytics support is limited (we'll use the OLTP database for analytics, not search engine)
- **Horizontal scaling is limited.** Meilisearch is primarily single-node. Cluster mode is newer; we'll need to plan migration if our index exceeds single-node limits (~20–30M documents based on community reports).
- **Smaller ecosystem.** Fewer community resources, plugins, third-party integrations vs ES.
- **Single-node availability**: must operate with hot-standby and snapshot/restore for HA. We accept this; downtime budget allows for it.

### Neutral / Trade-offs accepted

- **No mature security plugin ecosystem.** We rely on network-level isolation (VPC, security groups) rather than fine-grained access control inside Meilisearch.
- **Limited custom analyzers.** We can't write Lucene-style analyzers. For us, this is fine — we don't have custom tokenization needs.

## Benchmark Comparison

A 2-week spike compared three options against a representative dataset (1M listings, 10K queries from production logs replayed):

| Metric                                     | OpenSearch 2.13                         | Meilisearch 1.7       | PostgreSQL FTS (tsvector)          |
| ------------------------------------------ | --------------------------------------- | --------------------- | ---------------------------------- |
| Indexing time (1M docs)                    | 28 min                                  | 6 min                 | 12 min (with GIN index)            |
| p50 query latency                          | 18ms                                    | 8ms                   | 35ms                               |
| p95 query latency                          | 62ms                                    | 24ms                  | 110ms                              |
| p99 query latency                          | 140ms                                   | 75ms                  | 280ms                              |
| Typo tolerance (manual eval, 30 queries)   | Excellent (with edge n-grams)           | Excellent (built-in)  | Poor (requires extensions)         |
| Multi-language (VI tokenization)           | Acceptable (with custom analyzer setup) | Acceptable (built-in) | Acceptable (Vietnamese dictionary) |
| RAM at idle                                | 4.2 GB                                  | 1.1 GB                | 800 MB (shared with OLTP)          |
| RAM under 5K QPS                           | 8.5 GB                                  | 2.6 GB                | 1.4 GB                             |
| Setup complexity (engineer-hours)          | ~40 (cluster, tuning)                   | ~6                    | ~3 (already running PG)            |
| Re-indexing operational cost               | High (re-shard, takes hours)            | Low (snapshot + load) | Low                                |
| Cost at 5M docs, 2x HA (AWS, monthly est.) | ~$780                                   | ~$230                 | ~$0 incremental                    |

Benchmark code and full results: `bench/2026-03-search-comparison/`

Why not just PostgreSQL FTS:

- p95 query latency 110ms violates our 50ms SLO on read latency for instant-search.
- Typo tolerance and prefix matching require additional extensions (pg_trgm) and are slower.
- Faceted aggregations on PostgreSQL FTS are slow without materialized views, adding complexity.

Why not OpenSearch:

- 4× the operational complexity for capabilities we don't need.
- 3× the cost.
- Team has no operational experience; staffing risk.

## Alternatives Considered

### Option B: OpenSearch (rejected per benchmark above)

The default "industry choice" for full-text search.

**Why rejected**:

- Operational complexity disproportionate to our needs.
- Cost premium ~3× without proportionate benefit at our scale.
- Team would need significant ramp-up time.
- Most of its advanced features (custom analyzers, sophisticated ranking, security plugins, analytics) we don't currently need.

**When we would revisit**: if index exceeds 20M documents OR if we need sophisticated query DSL OR if compliance requires fine-grained access control inside the search engine.

### Option C: PostgreSQL Full-Text Search (rejected per benchmark)

Stay within the database.

**Why rejected**:

- Latency profile (p95 110ms) violates SLO.
- Typo tolerance and prefix search require extensions; quality is lower than dedicated engines.
- Co-locating heavy search workload with OLTP on same instance creates resource contention.

**When this would have been right**: if our scale were 10× smaller. For ~500K listings with < 50 QPS, PostgreSQL FTS is enough and the simplicity wins.

### Option D: Algolia (managed)

SaaS search-as-a-service.

**Why rejected**:

- Cost at our scale: estimated $3K+/month, multiples of our self-hosted Meilisearch projection.
- Data residency: must verify SG/MY region compliance; adds procurement overhead.
- Lock-in to vendor.

### Option E: Typesense

Similar profile to Meilisearch (Rust, simple).

**Why rejected**:

- Smaller community, less mature.
- Vietnamese tokenization less well-validated.
- Functionally similar to Meilisearch for our needs; we picked Meilisearch on slight ecosystem advantage.

## Open Questions

- **HA topology**: Meilisearch's clustering is newer than OpenSearch's. We'll start with active + hot-standby with periodic snapshot replication; will revisit if we need true clustering at month 6.
- **Index versioning during schema changes**: dual-write to old and new indices during migration. Will codify the playbook after first migration.
- **Re-ranking with ML signals (planned for 2027)**: Meilisearch supports custom ranking rules but not ML-trained re-rankers. We'll likely add a thin re-rank service in front when this becomes a priority.

## References

- Benchmark workspace: `bench/2026-03-search-comparison/` (this repo)
- Meilisearch documentation: https://www.meilisearch.com/docs
- "Choosing a search engine in 2024" — internal write-up by Search team
- Course material: Module 03 — Data at Scale
- ADR-002 (Outbox Pattern) — indexing pipeline depends on this

## Decision Record History

| Date       | Status   | Notes                                                       |
| ---------- | -------- | ----------------------------------------------------------- |
| 2026-03-03 | Proposed | First draft after benchmark spike                           |
| 2026-03-12 | Revised  | Added Algolia and Typesense as alternatives per CTO request |
| 2026-03-18 | Accepted | Approved; rollout planned end of Q2                         |
