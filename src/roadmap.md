# Interactive Roadmap

> **System Design Mastery —** _think in systems, not just code._
>
> A 6-month curriculum for the architect-in-waiting.

<p style="margin: 1.5rem 0; padding: 1rem; background: rgba(212,165,116,0.08); border-left: 3px solid #d4a574; border-radius: 4px;">
  Want the visual, interactive version with progress tracking?
  <strong><a href="./static/roadmap.html" target="_blank">Open the interactive roadmap ↗</a></strong>
  (saves your progress locally in this browser)
</p>

---

## At a Glance

|                 |                       |
| --------------- | --------------------- |
| **Modules**     | 8                     |
| **Topics**      | 42                    |
| **Capstones**   | 8                     |
| **Duration**    | ~24 weeks (~6 months) |
| **Time / week** | 8–12 hours            |

---

## The Curriculum

### Phase I — Foundations

_Think before you build._ · Weeks 1–8

#### Module 01 · Thinking in Systems

**Foundations · Week 1–3 · 3 weeks · 1 capstone**

The architect's mental model. Quality attributes, constraints, trade-offs.
Why "it works" is not "it's done." How to read a system before changing it.

`Quality Attributes` `CAP / PACELC` `Latency vs Throughput` `SLI / SLO / SLA` `Constraints` `Back-of-envelope`

→ [Open module](modules/01-thinking-in-systems.md)

---

#### Module 02 · Distributed Systems Theory

**Foundations · Week 4–6 · 3 weeks · 1 capstone**

The hard truths. Consistency models, consensus, time, failure detection.
Why everything you know about local programs is a lie at scale.

`Consistency Models` `Raft / Paxos` `Lamport Clocks` `CRDTs` `Failure Modes` `FLP / 2 Generals`

→ [Open module](modules/02-distributed-systems-theory.md)

---

#### Module 03 · Data at Scale

**Foundations · Week 7–8 · 2 weeks · 1 capstone**

Storage engines, indexing, partitioning, replication. SQL vs NoSQL trade-offs
informed by access patterns, not marketing. ACID, BASE, and everything between.

`LSM vs B-Tree` `Sharding` `Replication` `OLTP vs OLAP` `pgvector / Vector DBs` `CDC`

→ [Open module](modules/03-data-at-scale.md)

---

### Phase II — Patterns

_The language of design._ · Weeks 9–16

#### Module 04 · Architecture Styles

**Patterns · Week 9–11 · 3 weeks · 1 capstone**

Monoliths, microservices, modular monoliths, serverless. The honest comparison
most people skip. When to split, when to consolidate, when each style is right.

`Modular Monolith` `Microservices` `Hexagonal` `Clean / Onion` `Event-Driven` `Service Mesh`

→ [Open module](modules/04-architecture-styles.md)

---

#### Module 05 · Event-Driven & CQRS

**Patterns · Week 12–14 · 3 weeks · 1 capstone**

Decoupling through events. Outbox, sagas, idempotency, ordering. Why
eventual consistency is a feature, not a bug — and how to communicate that.

`Event Sourcing` `CQRS` `Outbox Pattern` `Sagas` `Kafka Internals` `Idempotency`

→ [Open module](modules/05-event-driven-cqrs.md)

---

#### Module 06 · Reliability Patterns

**Patterns · Week 15–16 · 2 weeks · 1 capstone**

Building systems that break gracefully. Circuit breakers, bulkheads, retries,
rate limiting, load shedding. Capacity planning without spreadsheets-of-doom.

`Circuit Breaker` `Bulkhead` `Backpressure` `Retry / Hedging` `Rate Limiting` `Chaos Eng.`

→ [Open module](modules/06-reliability-patterns.md)

---

### Phase III — Craft

_The architect's actual job._ · Weeks 17–24

#### Module 07 · Design at Scale

**Craft · Week 17–20 · 4 weeks · 1 capstone**

The classic interview prompts done right. Chat, feeds, search, payments, ride-share.
Not memorizing answers — building the muscle to design anything in 45 minutes.

`Chat (WhatsApp)` `Feed (Twitter)` `Search (Google)` `Payments` `Ride-share` `URL Shortener`

→ [Open module](modules/07-design-at-scale.md)

---

#### Module 08 · The Architect's Craft

**Craft · Week 21–24 · 4 weeks · Final capstone**

The non-technical 50%. ADRs, RFCs, technical leadership, influence without authority,
stakeholder communication. The skills you can't learn from a book on databases.

`ADRs / RFCs` `C4 Diagrams` `Stakeholders` `Tech Strategy` `Influence` `Career Path`

→ [Open module](modules/08-architects-craft.md)

---

## Three Mindset Shifts

The senior → architect transition is mostly mental. These three shifts run through every module.

### Shift 01 — From / To

> ~~Solve problems~~ → **Define constraints**

Seniors are paid to solve. Architects are paid to define the problem with such
clarity that the solution is obvious. The right question is harder than the right answer.

### Shift 02 — From / To

> ~~Functional first~~ → **Quality attributes first**

Features are table stakes. Architecture is the system's response to non-functional
forces: latency, availability, security, evolvability. These are what break in production.

### Shift 03 — From / To

> ~~Code is the artifact~~ → **Communication is**

ADRs, RFCs, diagrams, conversations. The architect's deliverable is alignment.
A great design that no one implements correctly is a failed design.

---

## What You'll Walk Away With

By week 24:

1. A portfolio of **8 production-grade designs** with full ADRs, C4 diagrams, and trade-off analysis.
2. A working capstone — a **distributed system you designed and built end-to-end**, deployable on AWS or k8s.
3. Fluency in the language architects speak: **quality attributes, CAP, consistency models, patterns** — not as buzzwords but as tools.
4. A repeatable framework for tackling **any system design problem** in 45 minutes or 6 months.
5. The non-technical craft: **writing ADRs that get read**, RFCs that get accepted, and presentations that move stakeholders.

---

## Progress Tracking

Use the [interactive roadmap](./static/roadmap.html) for clickable progress tracking — it saves to your browser's localStorage.

Or track here manually:

- [ ] **Module 01** — Thinking in Systems
- [ ] **Module 02** — Distributed Systems Theory
- [ ] **Module 03** — Data at Scale
- [ ] **Module 04** — Architecture Styles
- [ ] **Module 05** — Event-Driven & CQRS
- [ ] **Module 06** — Reliability Patterns
- [ ] **Module 07** — Design at Scale
- [ ] **Module 08** — The Architect's Craft

---

> _"An architect's job isn't to draw diagrams.
> It's to **make decisions reversible**, or **make irreversible decisions correctly.**"_
