# ADR-005: Migrate audit-log storage from Apache Cassandra to ScyllaDB

## Status

**Accepted** — 2026-04-30
**Authors**: Thanh Tran
**Reviewers**: Platform team, DBA group, Compliance
**Supersedes**: **ADR-2022-007 — "Use Apache Cassandra for audit log storage"**
**Review date**: 2027-04-30

## Context

In 2022, the team chose Apache Cassandra for the audit log store. The 2022 decision was sound at the time:

- Write-heavy (~5K writes/sec sustained at decision time)
- Time-series-like access patterns (audit logs queried by user_id + time range)
- Required retention: 7 years (regulatory)
- AP consistency model acceptable for an append-only audit log
- 3-node cluster was sufficient for projected 2-year load

Four years later, the system has grown:

| Metric                                   | At 2022 decision    | Today                                        |
| ---------------------------------------- | ------------------- | -------------------------------------------- |
| Sustained writes/sec                     | 5K                  | 38K                                          |
| Cluster size                             | 3 nodes (i3.xlarge) | 14 nodes (i3.2xlarge)                        |
| Total storage                            | 1.2 TB              | 41 TB                                        |
| Avg p99 write latency                    | 8ms                 | 47ms                                         |
| Frequency of GC pauses > 200ms           | Rare                | 3–5/hour per node                            |
| Operational incidents/month              | <1                  | 4–6 (mostly GC-related and node-replacement) |
| Engineering-hours/month on Cassandra ops | ~10                 | ~80                                          |

The team has invested significant effort in Cassandra tuning over the past year (JVM heap, GC algorithm, compaction strategy, off-heap caches). Improvements were marginal. The fundamental issue is that **JVM GC at our throughput is the bottleneck**, and no amount of tuning fixes it without ever-larger machines.

A spike in Q1 2026 evaluated alternatives. ScyllaDB emerged as the most promising:

- API-compatible with Cassandra (same CQL); migration is realistic
- C++ implementation, no GC pauses
- Public benchmarks and the well-known Discord migration (2022) showed dramatic improvements at comparable scale
- ScyllaDB Cloud + self-hosted options both available

We ran a 30-day production-parallel test: ScyllaDB cluster received a mirrored write stream of the audit events.

## Decision

**We will migrate the audit log store from Apache Cassandra to ScyllaDB Open Source 5.4.**

Migration approach:

1. **Phase 1** (month 0–1): Deploy 4-node ScyllaDB cluster in production. Begin dual-writing audit events to both Cassandra and ScyllaDB via an extended outbox relay.
2. **Phase 2** (month 1–2): Replay historical Cassandra data into ScyllaDB using `sstableloader`. Reconcile.
3. **Phase 3** (month 2–3): Shift reads to ScyllaDB. Cassandra remains receiving writes as a safety net.
4. **Phase 4** (month 3): Stop dual-writes. Snapshot Cassandra for backup; decommission cluster.

Total project budget: 3 months, 1.5 engineers.

## Consequences

### Positive

- **Dramatic operational improvement.** Production-parallel test showed:
  - p99 write latency: 47ms → 3.2ms (15× improvement)
  - GC pauses eliminated entirely (no GC in ScyllaDB)
  - Node count projected: 14 → 4 (3.5× reduction)
- **Significant cost reduction.** EC2 instance cost projected to drop ~65% ($28K/month → $10K/month).
- **CQL compatibility.** Application code requires zero changes (CQL queries, driver compatible).
- **Operational simplicity.** No JVM tuning. Single binary. Auto-tuning shards itself to CPU cores.
- **Headroom for growth.** ScyllaDB benchmarks suggest comfortable handling of 5–10× our current write rate on the same cluster size.

### Negative / Costs

- **Migration complexity and risk.** 7 years of audit data, regulatory implications if data is lost or corrupted in transit.
- **Smaller community than Cassandra.** Vendor risk (ScyllaDB Inc. as commercial backer). Mitigated: ScyllaDB Open Source is Apache 2.0; we self-host; vendor failure doesn't strand us.
- **Some edge-case feature gaps.** ScyllaDB doesn't support every Cassandra feature (e.g., some materialized view edge cases, secondary indexes have different performance profile). We validated all our queries during the spike; none affected.
- **Team learning curve.** New monitoring tools, new operational playbook. Estimate 1 month for the on-call team to be fully comfortable.
- **Migration takes engineering time.** ~3 person-months from a small team, but ROI is large (recovers ~10 engineer-hours/week of ops afterward).

### Neutral / Trade-offs accepted

- **Self-hosted, not managed.** We considered ScyllaDB Cloud but went self-hosted for cost and data residency. Will re-evaluate if ops burden grows.

## Why This Supersedes ADR-2022-007

The 2022 ADR was correct given 2022 context. **Architectural decisions decay** — the world changes. Specifically:

- Scale grew 7.6× beyond the projected curve
- JVM GC characteristics that were acceptable at low throughput became unacceptable at high throughput
- The market for Cassandra-compatible stores has matured (ScyllaDB has stabilized significantly since 2022)
- The cost of operating Cassandra at our scale exceeded the cost of migration

**This is not a criticism of the 2022 decision.** It is the normal lifecycle of architectural choices. The 2022 ADR even included a review date and the open question "How do we handle GC pause issues if they emerge?" — the answer turned out to be "migrate."

The lesson encoded by this supersession: **revisit ADRs on their review dates, with data.**

## Alternatives Considered

### Option B: Stay on Cassandra, upgrade hardware

Buy more / bigger machines.

**Why rejected**:

- Linear cost growth at sub-linear performance gains.
- GC is a software-level issue; bigger machines don't fix it.
- Even doubled instances we estimated would push p99 from 47ms to ~30ms — still 10× worse than ScyllaDB.

### Option C: Stay on Cassandra, migrate to Cassandra 5.0 with TPC architecture

Cassandra 5.0 brings significant performance improvements (Thread-Per-Core).

**Why rejected**:

- 5.0 was released only 3 months before this decision; not yet production-hardened.
- Improvements significant but still GC-bound at our scale (per public benchmarks).
- We'd be the early adopter and bear the migration risk anyway.
- ScyllaDB has been production-stable for years.

### Option D: Migrate to a different data model (e.g., TimescaleDB)

Audit logs are time-series-shaped; TimescaleDB is built for this.

**Why rejected**:

- Single-master architecture limits write throughput vs Cassandra/Scylla's masterless model.
- Migration cost is substantial (CQL → SQL, application changes).
- Solving a different problem; we're write-throughput-bound, not query-pattern-bound.

### Option E: Managed Cassandra (AWS Keyspaces, DataStax Astra)

Outsource the ops.

**Why rejected**:

- AWS Keyspaces has incomplete Cassandra API compatibility; some of our queries don't work.
- Astra cost is 4–6× our self-hosted estimate.
- Doesn't solve the fundamental architecture issue, just hides it.

## Open Questions

- **Backup/snapshot frequency for ScyllaDB**: matching Cassandra's nightly + WAL or different cadence? Decision pending operational review at month 2.
- **Compaction strategy fine-tuning**: time-window compaction is the obvious default for audit logs; we'll validate against production traffic for 30 days post-migration.
- **Monitoring**: integrating ScyllaDB's Prometheus exporters with our existing Grafana dashboards. Some metrics differ from Cassandra; will rewrite alerting rules during Phase 1.
- **Long-term retention**: should we tier old audit data (>2 years) to S3 with Athena query? Out of scope for this migration but a likely follow-up.

## References

- Production-parallel test results: `bench/2026-q1-scylladb-eval/` (this repo)
- Original ADR-2022-007: archived in `decisions/superseded/`
- Discord's migration write-up (https://discord.com/blog/how-discord-stores-trillions-of-messages) — reference architecture
- ScyllaDB documentation: https://docs.scylladb.com/
- Course material: Module 03 — Data at Scale
- Incident reports referenced: INC-2025-11-07-cass-gc-cascade.md, INC-2026-01-14-node-replacement-saga.md

## Decision Record History

| Date       | Status   | Notes                                        |
| ---------- | -------- | -------------------------------------------- |
| 2026-04-08 | Proposed | First draft after production-parallel test   |
| 2026-04-22 | Revised  | Migration plan refined per DBA team feedback |
| 2026-04-30 | Accepted | Approved; migration kickoff scheduled        |
