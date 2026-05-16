# Capacity Estimator

A back-of-envelope estimator for system design. Convert DAU + behavioral inputs into QPS, storage, and bandwidth — the way you should at the start of every design.

## Why This Exists

Module 01 says: _numbers anchor every decision_. Most engineers do this math sloppily in interviews — or skip it entirely. This tool makes it precise, repeatable, and version-controllable.

You can:

- Run scenarios from the CLI: `./capest estimate --preset twitter`
- Save scenarios as YAML
- Get markdown-formatted reports for design docs

## Quick Start

```bash
go install ./cmd/capest
capest estimate --preset twitter
```

Output (abbreviated):

```
Twitter-scale estimate
=====================
Inputs:
  DAU:              200M
  Actions/user/day: 50 (post=2, read=48)
  Avg post size:    280 bytes
  Read amplification: 100 (each post seen by avg 100 followers)
  Peak multiplier:  3x

Derived:
  Writes/sec avg:   ~4.6 K/sec
  Writes/sec peak:  ~14 K/sec
  Reads/sec avg:    ~111 K/sec
  Reads/sec peak:   ~333 K/sec
  Storage/day:      ~115 GB
  Storage/year:     ~42 TB
  Bandwidth (read): ~31 MB/sec peak
```

## Presets

Built-in presets you can use as starting points:

| Preset        | Description                                |
| ------------- | ------------------------------------------ |
| `twitter`     | Twitter-scale microblog                    |
| `instagram`   | Photo-share, fan-out heavy                 |
| `whatsapp`    | Messaging, no fan-out                      |
| `uber`        | Ride share, location-heavy                 |
| `propertyhub` | The author's domain (real estate platform) |

## Custom Scenarios

Save to YAML:

```yaml
# my-system.yaml
name: "My System"
dau: 5_000_000
actions:
  writes_per_user_per_day: 10
  reads_per_user_per_day: 200
avg_write_bytes: 1024
read_amplification: 1
peak_multiplier: 4
replication_factor: 3
retention_days: 365
```

Run:

```bash
capest estimate --file my-system.yaml
```

## Output Formats

- Default: human-readable text
- `--format markdown`: insert directly into design docs
- `--format json`: machine-readable for scripting

## What This Project Demonstrates

This is the smallest "real" Go project — it gives you a tour of:

- CLI structure with `cobra` (production-typical)
- YAML configuration with validation
- Testable pure functions (the math is unit-tested)
- Multiple output formats from one core function

## Build and Test

```bash
go test ./...     # run tests
go build ./...    # build
make run          # run example
```
