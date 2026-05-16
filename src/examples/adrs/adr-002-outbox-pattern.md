# ADR-002: Adopt the Transactional Outbox Pattern for cross-service events

## Status

**Accepted** — 2026-03-04
**Authors**: Thanh Tran
**Reviewers**: Platform team, Notifications team, Search team
**Supersedes**: none
**Review date**: 2027-03-04

## Context

The Listings service publishes events consumed by:

- **Search service** — must reindex Meilisearch within seconds of a listing change
- **Notifications service** — agents are alerted when their listings are viewed, favorited, or status-changed
- **Analytics pipeline** — Kafka → S3 → Athena for business reporting

The current implementation in the monolith uses _dual writes_:

```php
$db->save($listing);
$kafka->publish('ListingUpdated', $listing);
```

On 2026-02-22, an incident exposed this pattern's failure mode:

- 14:32:08 — Kafka broker in AZ-1a became unresponsive (network issue)
- 14:32:08–14:34:12 — Listings service committed ~340 listings to PostgreSQL, but `kafka->publish` failed silently (Kafka client retried then dropped after 30s timeout)
- 14:34:12 — Kafka recovered; new publishes succeeded
- **Result**: 340 listings updated in DB but never propagated to search index. Customer-facing impact: stale prices on search results for 4 hours until manual reconciliation script ran.

Post-incident review identified this as a textbook **dual-write problem**: when two systems must be updated atomically, doing them sequentially with no coordination cannot guarantee consistency across failure modes.

This is a structural issue, not a Kafka issue. Retrying does not solve it: an app crash _between_ the DB commit and the Kafka publish leaves the DB ahead with no event ever published.

## Decision

**We will adopt the Transactional Outbox Pattern for all cross-service events emitted from the Listings service.**

Implementation:

1. An `outbox` table in the same PostgreSQL database as the business tables.
2. Domain operations write business data **and** the outbox event in the **same database transaction**.
3. A separate relay process polls the `outbox` table and publishes pending events to Kafka, marking them as published after broker acknowledgment.
4. Consumers must be idempotent (handle duplicate delivery).

We will use `SELECT ... FOR UPDATE SKIP LOCKED` to allow horizontal scaling of relay workers without coordination.

## Consequences

### Positive

- **At-least-once delivery is guaranteed.** The DB transaction is the boundary; if it commits, the event is durable. The relay will publish it eventually.
- **No more dual-write incidents.** The 340-listing scenario from Feb 22 cannot recur with this pattern.
- **Audit trail.** Every emitted event is in the outbox table with timestamps. Easy to replay, debug, or recover.
- **No new infrastructure.** Uses the database we already operate; no Saga coordinator service needed.
- **Backpressure flows naturally**: if Kafka is slow, the outbox grows but doesn't fail. If the outbox grows unbounded, we'll see it in monitoring before it becomes critical.

### Negative / Costs

- **Added latency**: events are published asynchronously, not synchronously. Search index lag will go from "0 seconds (mostly)" to "consistently 1–3 seconds." This is acceptable for our use case but must be communicated.
- **All consumers must be idempotent.** Some current consumers are not. We will need to audit and fix:
  - Notifications service: needs `(event_id, user_id)` idempotency key
  - Analytics pipeline: idempotent by design (upsert by event_id)
  - Search service: already idempotent (replays full listing state)
- **Operational concern**: outbox table grows. We will partition by month and archive partitions older than 30 days.
- **Two-phase mental model**: developers must remember that "publishing an event" means "inserting to outbox," not "calling Kafka." A small training cost.

### Neutral / Trade-offs accepted

- **At-least-once, not exactly-once.** This is fundamental to distributed systems, not a limitation of the pattern. Consumers handle it.
- **Cannot order events strictly across partitions.** We will use `aggregate_id` as the Kafka partition key, which guarantees per-aggregate ordering. Cross-aggregate ordering is not guaranteed and not required by current consumers.

## Alternatives Considered

### Option B: Change Data Capture (CDC) via Debezium

Read the PostgreSQL WAL directly; emit changes to Kafka.

**Why rejected**:

- Couples Kafka consumers to the database schema. Renaming a column in `listings` table would break downstream.
- Less control over event shape: we want clean domain events (`ListingPriceChanged`), not raw row diffs.
- Operational complexity: Debezium runs as a Kafka Connect cluster, separate failure domain.
- The Outbox pattern gives us 95% of CDC's benefit at 20% of the operational cost for our scale.

**When we would revisit**: if we add many more services needing many different event shapes, the centralized table approach becomes a bottleneck. Re-evaluate at 10+ consumer services.

### Option C: Two-Phase Commit (XA transactions)

Synchronously commit to DB and Kafka together.

**Why rejected**:

- Kafka does not natively support XA. We'd need an external coordinator.
- Blocking: if the coordinator fails mid-transaction, participants block.
- Worse availability characteristics than our current dual-write.
- Industry consensus: 2PC across heterogeneous systems is to be avoided.

### Option D: Application-level retry queue

App writes to DB, then publishes to Kafka; on publish failure, queues for retry in a local queue.

**Why rejected**:

- App crashes lose the in-flight queue.
- Doesn't solve the "DB committed, app died before queueing" case.
- Half-solutions for distributed transactions are usually worse than full solutions.

### Option E: Saga orchestrator (Temporal / Conductor)

Use a workflow engine to coordinate cross-service operations.

**Why rejected**:

- Overkill for our current needs. Outbox solves the event-publishing problem. Sagas solve a different problem (multi-service workflows with compensation).
- We have one cross-service workflow today (listing-published-triggers-notification); not a saga-shaped need.
- Will reconsider when we have 3+ services collaborating in a workflow with rollback needs.

## Open Questions

- **Outbox table partitioning**: monthly or weekly? Will decide based on row count after 30 days of production data.
- **Relay scaling**: start with 2 instances behind `FOR UPDATE SKIP LOCKED`; revisit if backlog ever exceeds 10K events.
- **Schema for the `payload` JSONB column**: do we use Avro/Protobuf binary or JSON? Decision: JSON for v1 (debuggability); reconsider when payload sizes hit a meaningful storage cost.
- **Cleanup policy**: archive published rows after 7 days (keep last week for debug) or 30 days (defensive)? Will set to 30 initially, monitor disk growth.

## References

- Original Outbox pattern write-up: Gunnar Morling, "Reliable Microservices Data Exchange With the Outbox Pattern" (https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/)
- Course material: Module 05 — Event-Driven & CQRS
- Incident report: `INC-2026-02-22-search-stale.md` (internal)
- Implementation reference: `examples/05-outbox/` in this repo

## Decision Record History

| Date       | Status   | Notes                                         |
| ---------- | -------- | --------------------------------------------- |
| 2026-02-23 | Proposed | Drafted in response to Feb 22 incident        |
| 2026-02-28 | Revised  | Added Option B (CDC) per Search team feedback |
| 2026-03-04 | Accepted | Approved; rollout target end of Q1            |
