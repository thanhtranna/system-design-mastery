# ADR Examples

A working ADR (Architecture Decision Record) is worth more than a thousand abstract templates. These five examples show what real ADRs look like — context, decision, honest trade-offs, alternatives, and even a superseded decision.

Use them as **reference**: copy the structure, the voice, the level of rigor. Adapt them to your context.

## The Eight Examples

| #   | Title                                                                              | Module | Notable for                                             |
| --- | ---------------------------------------------------------------------------------- | ------ | ------------------------------------------------------- |
| 001 | [PostgreSQL as primary OLTP store](./adr-001-postgresql-oltp.md)                   | Mod 03 | Classic boring-technology decision; honest alternatives |
| 002 | [Adopt Outbox Pattern for cross-service events](./adr-002-outbox-pattern.md)       | Mod 05 | Incident-driven; specific failure mode named            |
| 003 | [Modular Monolith, not microservices](./adr-003-modular-monolith.md)               | Mod 04 | Contrarian decision against fashion                     |
| 004 | [Meilisearch over OpenSearch for product search](./adr-004-meilisearch.md)         | Mod 03 | Benchmark-driven with concrete numbers                  |
| 005 | [Migrate from Cassandra to ScyllaDB](./adr-005-scylladb.md)                        | Mod 03 | Supersedes a 4-year-old decision                        |
| 006 | [Orchestration-Based Saga for distributed txns](./adr-006-saga-pattern.md)         | Mod 05 | Orchestration vs choreography; clear rationale          |
| 007 | [URL Path Versioning for public API](./adr-007-api-versioning.md)                  | Mod 08 | Pragmatism over REST purity; concrete alternatives      |
| 008 | [Cache-aside for catalogue, write-through for cart](./adr-008-caching-strategy.md) | Mod 03 | Different strategies for different data types           |

## What Makes a Good ADR

Read these examples and notice what they share:

1. **Title is the decision, not the question.** "Use PostgreSQL" not "PostgreSQL vs MySQL?"
2. **Context names the forces**: scale numbers, team competence, deadline.
3. **Decision is one sentence, active voice.**
4. **Consequences split into positive / negative / neutral.** Honesty about cost.
5. **Alternatives have concrete rejection reasons** — not "we didn't like it."
6. **Open questions are named.** No ADR is omniscient.
7. **Review date** is set for every decision that might age.

## What Makes an ADR Get Read

- **1–3 pages**. Not 12.
- **Skim-able**: headers, bullets, no walls of text.
- **The conclusion is at the top.** Reader knows decision in 30 seconds.
- **Written for someone who joined the company yesterday.**

If the reader has to ask "but why didn't we just do X?" — your Alternatives section failed.
