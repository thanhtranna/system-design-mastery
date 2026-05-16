# System Design Mastery — From Senior Engineer to Software Architect

> _"An architect's job isn't to draw diagrams. It's to make decisions reversible, or make irreversible decisions correctly."_

A 6-month, production-grade curriculum for senior engineers transitioning into software architecture roles. Designed for engineers who've built systems and want to design them.

---

## Who This Is For

- Senior engineers (5+ years) targeting **tech lead → architect → principal** paths
- Engineers who can ship features but feel underwater in design reviews
- Anyone preparing for **L6/Staff/Principal** system design interviews
- Engineers who've read DDIA and want the "what next"

If you can read this and think "I know the words but not the system" — this is for you.

---

## How To Use This Course

**Time commitment**: 8–12 hours/week for ~24 weeks (~6 months). Adjustable.

**Per module, expect**:

1. **Read** (2–3h) — concept deep-dive with Go code examples
2. **Practice** (3–4h) — exercises + design challenges
3. **Build** (4–6h) — capstone project
4. **Communicate** (1–2h) — write an ADR, do a mock interview

**Don't skip the writing.** The ADR and mock-interview parts are 50% of the value. Architects who can't communicate are senior engineers with a fancier title.

---

## Curriculum Map

### Phase I — Foundations (Weeks 1–8)

_Think before you build._

| #   | Module                         | Focus                                                         | Weeks |
| --- | ------------------------------ | ------------------------------------------------------------- | ----- |
| 01  | **Thinking in Systems**        | Quality attributes, constraints, trade-offs, back-of-envelope | 1–3   |
| 02  | **Distributed Systems Theory** | Consistency, consensus, time, failure                         | 4–6   |
| 03  | **Data at Scale**              | Storage engines, indexing, sharding, replication              | 7–8   |

### Phase II — Patterns (Weeks 9–16)

_The language of design._

| #   | Module                   | Focus                                             | Weeks |
| --- | ------------------------ | ------------------------------------------------- | ----- |
| 04  | **Architecture Styles**  | Monolith → microservices, hexagonal, event-driven | 9–11  |
| 05  | **Event-Driven & CQRS**  | Outbox, sagas, event sourcing, Kafka internals    | 12–14 |
| 06  | **Reliability Patterns** | Circuit breakers, bulkheads, rate limiting, chaos | 15–16 |

### Phase III — Craft (Weeks 17–24)

_The architect's actual job._

| #   | Module                    | Focus                                                | Weeks |
| --- | ------------------------- | ---------------------------------------------------- | ----- |
| 07  | **Design at Scale**       | The interview classics: chat, feed, search, payments | 17–20 |
| 08  | **The Architect's Craft** | ADRs, RFCs, C4, communication, leadership            | 21–24 |

---

## What Makes This Course Different

**1. Trade-offs over recipes.** Every pattern includes when _not_ to use it.

**2. Go-flavored code, language-agnostic ideas.** Examples in Go because runtime models matter, but the concepts transfer to any stack.

**3. Production scars, not just textbooks.** Each module includes real failure modes — the ones that get you paged at 3 AM.

**4. ADR practice baked in.** You'll write 8 ADRs by the end. This is how you actually become an architect.

**5. Mock interviews per module.** With a rubric, not just "good answer."

---

## File Structure

```
system-design-mastery/
├── README.md                    # This file
├── roadmap.html                 # Interactive visual roadmap
└── modules/
    ├── 01-thinking-in-systems.md
    ├── 02-distributed-systems-theory.md
    ├── 03-data-at-scale.md
    ├── 04-architecture-styles.md
    ├── 05-event-driven-cqrs.md
    ├── 06-reliability-patterns.md
    ├── 07-design-at-scale.md
    └── 08-architects-craft.md
```

Each module follows this structure:

1. **Mindset** — the mental shift this module forces
2. **Core Concepts** — theory with concrete examples
3. **Patterns** — the named techniques you'll use
4. **Go Implementation** — production-grade code
5. **Trade-offs Table** — when to use what
6. **Real-World Failures** — what breaks in production
7. **Design Challenges** — practice problems
8. **Capstone Project** — build something real
9. **ADR Practice** — write a decision record
10. **Mock Interview** — practice prompt + rubric
11. **Further Reading** — curated, not exhaustive

---

## Recommended Prerequisites

Already mastered:

- A backend language (Go, Java, Python, etc.) at a senior level
- SQL fundamentals, basic NoSQL exposure
- HTTP, REST, basic networking
- Some experience with cloud (AWS preferred for examples)
- Has shipped production systems

If you don't have these — finish a "backend engineer" path first. This course assumes you've earned scars.

---

## The Three Mindset Shifts

These run through every module. Internalize them.

### Shift 1: _From solving to defining_

Seniors are paid to solve problems. Architects are paid to _define_ problems with such clarity that the solution becomes obvious. The right question is harder than the right answer.

When stuck, the architect asks: **"What are we actually optimizing for?"**

### Shift 2: _From functional to quality attributes_

Features are table stakes. Architecture is the system's response to **non-functional forces**: latency, availability, security, evolvability. These are what break in production, not "the button doesn't work."

When designing, the architect asks: **"Which quality attributes are we trading off, and at what cost?"**

### Shift 3: _From code to communication_

Code is one artifact. The architect's real deliverable is **alignment**: ADRs, RFCs, diagrams, stakeholder conversations, getting 12 teams pointing the same direction. A great design no one implements correctly is a failed design.

When working, the architect asks: **"Will the team understand and accept this in 6 months when I'm not in the room?"**

---

## License

[![License: CC BY 4.0](https://img.shields.io/badge/License-CC_BY_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by/4.0/)

Use freely, attribution required when sharing or redistributing.

---

## Next Step

→ Open `roadmap.html` in a browser for the visual map
→ Then start with `modules/01-thinking-in-systems.md`

Don't speedrun. The point isn't completion — it's transformation.
