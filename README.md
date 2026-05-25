# System Design Mastery

> _"An architect's job isn't to draw diagrams. It's to make decisions reversible, or make irreversible decisions correctly."_

A 6-month, production-grade curriculum for senior engineers transitioning into software architecture roles — built as an [mdBook](https://rust-lang.github.io/mdBook/).

**Author:** Thanh Tran

---

## What This Is

A structured self-study course covering the full arc from senior engineer to software architect. Not another "how to pass a FAANG interview" guide — this is about becoming someone who can define problems, navigate trade-offs, and lead technical decisions.

**Time commitment:** 8–12 hours/week for ~24 weeks.

---

## Who It's For

- Senior engineers (5+ years) targeting **tech lead → architect → principal** paths
- Engineers preparing for **L6/Staff/Principal** system design interviews
- Anyone who's read DDIA and wants to know what's next

---

## Curriculum

### Phase I — Foundations (Weeks 1–8)

| #   | Module                     | Focus                                                         |
| --- | -------------------------- | ------------------------------------------------------------- |
| 01  | Thinking in Systems        | Quality attributes, constraints, trade-offs, back-of-envelope |
| 02  | Distributed Systems Theory | Consistency, consensus, time, failure                         |
| 03  | Data at Scale              | Storage engines, indexing, sharding, replication              |

### Phase II — Patterns (Weeks 9–16)

| #   | Module               | Focus                                             |
| --- | -------------------- | ------------------------------------------------- |
| 04  | Architecture Styles  | Monolith → microservices, hexagonal, event-driven |
| 05  | Event-Driven & CQRS  | Outbox, sagas, event sourcing, Kafka internals    |
| 06  | Reliability Patterns | Circuit breakers, bulkheads, rate limiting, chaos |

### Phase III — Craft (Weeks 17–24)

| #   | Module                | Focus                                                 |
| --- | --------------------- | ----------------------------------------------------- |
| 07  | Design at Scale       | Chat, feed, search, payments — the interview classics |
| 08  | The Architect's Craft | ADRs, RFCs, C4, communication, leadership             |

---

## What's Included

- **8 core modules** with concept deep-dives, Go code examples, and trade-off tables
- **9 worked system designs** — Real Estate Search, URL Shortener, Chat, Payment Gateway, Social Feed, Video Streaming, Ride-Sharing, Notification Service, Rate Limiter
- **8 ADR examples** covering PostgreSQL, Outbox Pattern, Modular Monolith, Meilisearch, ScyllaDB, Saga Pattern, API Versioning, and Caching Strategy
- **6 runnable Go code examples** — Capacity Estimator, Vector Clock KV, Consistent Hashing, Write-Ahead Log, Outbox Pattern, Circuit Breaker
- **Practice resources** — Anki flashcards, mock interview rubrics, architecture review checklists, production readiness checklist, RFC template
- **7 real failure case studies** — Knight Capital, GitHub, Roblox, AWS S3, Discord, Facebook BGP, Cloudflare, Slack
- **Interactive visual roadmap**

---

## Project Structure

```bash
sd-mastery/
├── README.md                   # This file
├── book.toml                   # mdBook configuration
├── src/
│   ├── SUMMARY.md              # Table of contents
│   ├── introduction.md
│   ├── roadmap.md
│   ├── modules/                # 8 core curriculum modules
│   ├── examples/
│   │   ├── worked-designs/     # 9 system design walkthroughs
│   │   ├── adrs/               # 8 Architecture Decision Records
│   │   └── code/               # 6 runnable Go examples
│   └── practice/
│       ├── flashcards/
│       ├── checklists/         # Rubrics, checklists, RFC template
│       └── case-studies/       # Real-world failure post-mortems
└── book/                       # Generated HTML output
```

---

## Getting Started

### Prerequisites

- [Rust](https://www.rust-lang.org/tools/install)
- [mdBook](https://rust-lang.github.io/mdBook/guide/installation.html)
- [mdbook-mermaid](https://github.com/badboy/mdbook-mermaid) (for diagrams)

```bash
cargo install mdbook
cargo install mdbook-mermaid
```

### Build and Serve

```bash
cd sd-mastery
mdbook serve --open
```

The book will open at `http://localhost:3000`.

---

## License

[![License: CC BY 4.0](https://img.shields.io/badge/License-CC_BY_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by/4.0/)

Free to use with attribution when sharing or redistributing.
