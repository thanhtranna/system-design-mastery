# 03 — Consistent Hashing

> Maps to **Module 03: Data at Scale**. Makes consistent hashing tangible: add a node, watch minimal keys move.

## What It Shows

A consistent hash ring with virtual nodes. The core question in distributed data systems: *when you add or remove a node, how many keys have to move?*

With naive sharding (`key % N`): adding one node remaps ~50% of all keys. With consistent hashing: only `1/N` of keys move — the minimum theoretically possible.

This implementation lets you observe that directly:

```
$ make demo

Ring with 3 nodes, 150 virtual nodes each.
Distributing 10000 keys...
  node-A: 3312 keys (33.1%)
  node-B: 3389 keys (33.9%)
  node-C: 3299 keys (33.0%)

Adding node-D...
  node-A: 2498 keys  (-814, -24.6%)
  node-B: 2501 keys  (-888, -26.2%)
  node-C: 2498 keys  (-801, -24.3%)
  node-D: 2503 keys  (+2503)
  Keys remapped: 2503 / 10000 (25.0%)  ← close to ideal 1/4

Removing node-B...
  node-A: 3332 keys  (+834)
  node-C: 3329 keys  (+831)
  node-D: 3339 keys  (+836)
  Keys remapped: 2501 / 10000 (25.0%)  ← again, ~1/N
```

## Why This Matters in Interviews

The interviewer asks "how do you shard your database?" or "how does your cache cluster handle node failure?" Consistent hashing is the canonical answer. Being able to explain virtual nodes and why they improve key distribution is the difference between a mid-level and senior-level answer.

The code makes the mental model concrete: you'll have run it and watched the numbers.

## How to Run

```bash
make test    # run unit tests
make run     # interactive CLI: add/remove nodes, look up keys
make demo    # automated demo showing remapping percentages
```

**Requirements**: Go 1.22+. No external dependencies (pure stdlib).

## What to Look At

1. `ring/ring.go` — the ring implementation. Key insight: virtual nodes are just multiple hash positions for the same physical node. `AddNode` hashes N virtual node IDs and inserts them into a sorted slice.

2. The `lookup` function: binary search on the sorted slice, wrap around at the end. 10 lines of code.

3. `ring/ring_test.go` — the `TestRemapping` test measures actual remapping percentage. Run it and see the numbers.

## The Virtual Node Insight

Without virtual nodes, 3 physical nodes → 3 points on the ring → key distribution is highly uneven (depends entirely on where the 3 points land).

With 150 virtual nodes per physical node: 3 × 150 = 450 points distributed around the ring. By the law of large numbers, each physical node ends up owning ~1/3 of the ring regardless of the hash function's behaviour on the node names.

The trade-off: more virtual nodes = better distribution but more memory (the ring is a sorted slice of all virtual node positions). 100-200 virtual nodes per physical node is the typical production setting.

## Project Structure

```
03-consistent-hashing/
├── README.md
├── Makefile
├── go.mod
├── ring/
│   ├── ring.go        # consistent hash ring
│   └── ring_test.go   # unit tests + remapping benchmark
└── cmd/
    └── demo/
        └── main.go    # CLI demo
```
