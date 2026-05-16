# 04 — Write-Ahead Log (WAL)

> Maps to **Module 04: Architecture Styles** — specifically the durability guarantees that underpin databases, message queues, and any system that cannot lose data.

## What It Shows

A minimal Write-Ahead Log (WAL): an append-only log file that records every mutation *before* applying it to the in-memory state machine. On startup, the system replays the log to reconstruct state. Crash recovery happens automatically.

This is how PostgreSQL, Kafka, etcd, and virtually every durable system works internally.

## The Core Idea

```
Without WAL:
  write to memory → crash → state is gone

With WAL:
  write to log (disk) → apply to memory → crash → replay log on startup → state restored
```

The WAL is the only thing that makes "I wrote it" a meaningful statement.

## What to Observe

```bash
$ make demo

=== Write-Ahead Log Demo ===

[1] Writing keys to store...
  SET name=Alice
  SET city=Singapore
  SET role=engineer

[2] Simulating crash (close without cleanup)...

[3] Recovering from WAL...
  Replayed 3 entries

[4] State after recovery:
  name  = Alice
  city  = Singapore
  role  = engineer

Recovery complete. No data lost.

[5] Compaction: 3 entries → 3 entries (already minimal)
```

Then try modifying `cmd/demo/main.go` to crash mid-write and observe partial recovery.

## How to Run

```bash
make test    # unit tests for WAL and store
make demo    # automated crash-and-recover demonstration
```

**Requirements**: Go 1.22+. No external dependencies.

## What to Look At

1. `wal/wal.go` — `Append` writes a length-prefixed, checksummed record to disk and calls `Sync` before returning. This is the durability guarantee.

2. `wal/wal.go` — `Recover` reads entries from the log file, validates checksums, and replays them. Truncated or corrupt entries (from a crash mid-write) are detected and skipped.

3. `store/store.go` — `Set` calls `wal.Append` before updating the in-memory map. The order matters: log first, memory second.

4. `wal/wal_test.go` — `TestCrashRecovery` simulates a crash by closing the WAL without flushing, then reopening and verifying the recovered state.

## The Durability Trade-off

`Sync` (fdatasync) ensures the OS has flushed the write to physical media. It's expensive: ~1ms on SSD, ~10ms on spinning disk.

Databases offer a knob: `fsync=on` (PostgreSQL default) → sync every commit → durable, slower. `fsync=off` → no sync → fast, loses data on crash.

Kafka has the same knob: `log.flush.interval.messages`. Most Kafka deployments rely on replication, not fsync, for durability — a different trade-off.

## Project Structure

```
04-wal/
├── README.md
├── Makefile
├── go.mod
├── wal/
│   ├── wal.go         # WAL implementation: append, recover, compact
│   └── wal_test.go    # crash recovery tests
├── store/
│   ├── store.go       # WAL-backed KV store
│   └── store_test.go  # integration tests
└── cmd/
    └── demo/
        └── main.go    # crash-and-recover demonstration
```
