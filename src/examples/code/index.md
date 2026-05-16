# Runnable Code Examples

Four production-shaped Go projects that make the patterns from the course **tangible and executable**. You can `git clone`, `make run`, and see behavior.

These are not snippets glued into a slide deck. They're complete, runnable, tested. Each has Docker Compose so you can spin up the full stack including dependencies.

## Projects

| Project                                             | Maps to module | What it shows                                                                                          |
| --------------------------------------------------- | -------------- | ------------------------------------------------------------------------------------------------------ |
| [01 — Capacity Estimator](./01-capacity-estimator/) | Module 01      | A CLI tool you'll actually use in interviews. Take DAU + behavior → QPS/storage/bandwidth              |
| [02 — Vector Clock KV](./02-vector-clock-kv/)       | Module 02      | 3-node distributed KV store. Concurrent writes detected via vector clocks; client reconciles           |
| [05 — Outbox Pattern](./05-outbox/)                 | Module 05      | Postgres + Kafka with the Transactional Outbox pattern. Idempotent consumer. Network failure resilient |
| [06 — Circuit Breaker](./06-circuit-breaker/)       | Module 06      | Production-grade circuit breaker library with rolling window, half-open probes, and a chaos test       |

## Why These Exist

Reading about the Outbox pattern doesn't teach you the pattern. **Writing the relay** does. Reading about vector clocks doesn't teach them. **Watching two nodes produce divergent state** does.

These projects make the abstract concrete. They're also good portfolio code — well-structured, tested, properly Dockerized. You can put them on your CV.

## Suggested Order

If you're working through the modules in sequence:

1. **Capacity Estimator** (with Module 01) — tiny, accessible, immediate utility.
2. **Vector Clock KV** (with Module 02) — see distributed-systems theory operate.
3. **Outbox** (with Module 05) — the event-driven cornerstone pattern.
4. **Circuit Breaker** (with Module 06) — production-grade reliability building block.

## Each Project's Structure

A consistent layout across all four:

```
project/
├── README.md           # What it does, how to run, what to look at
├── Makefile            # `make test`, `make run`, `make demo`
├── go.mod
├── docker-compose.yml  # (where dependencies are needed)
├── Dockerfile          # (where containerized)
├── cmd/                # binaries
└── <library packages>  # the interesting code
```

## What You'll Take Away

After working through these four projects, you'll have:

- A capacity estimation tool you'll keep using
- A working mental model of vector clocks and concurrent writes
- The Outbox pattern in your bones — you'll never propose a dual-write design again
- A circuit breaker library you could lift into a real project tomorrow

That's all four reliability/correctness building blocks, runnable on your laptop.
