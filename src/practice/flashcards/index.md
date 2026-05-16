# Flashcard Deck

A curated set of ~150 flashcards covering key concepts from all 8 modules. Designed for **spaced repetition** — review daily, 10 minutes, for sustained retention.

## How To Use

### Option A: Anki (recommended)

1. Install [Anki](https://apps.ankiweb.net/) (free, all platforms)
2. Download [`flashcards.csv`](./flashcards.csv) from this repo
3. In Anki: File → Import → select the CSV
4. Choose deck name "System Design Mastery"
5. Map: Field 1 → Front, Field 2 → Back, Field 3 → Tags
6. Review daily

### Option B: Plain reading

The cards are also rendered below as a list. Less effective without spaced repetition, but workable.

### Option C: Build your own

The plain CSV makes it easy to convert to Quizlet, Mochi, Memorize, or any other SRS tool.

## Recommended Schedule

- **Week 1–4**: focus on Phase I cards (modules 01–03). New cards: 5/day, reviews: ~15/day.
- **Week 5–8**: add Phase II cards. New cards: 5/day, reviews: ~20/day.
- **Week 9+**: add Phase III cards. Maintenance: ~10 min/day.

After 90 days of consistent review, recall should be reliable for ~85% of the deck. That's the goal.

## Tags

Cards are tagged by module and concept type:

- `m01`, `m02`, ..., `m08` — module number
- `concept` — definition or principle
- `tradeoff` — comparison or trade-off
- `pattern` — named pattern
- `failure-mode` — what breaks and why
- `number` — order-of-magnitude facts

Use tags to filter in Anki for focused review.

## Sample Cards

### Foundations (Module 01)

> **Front**: What does PACELC add to CAP?
>
> **Back**: PACELC says: if Partitioned, choose between **A**vailability and **C**onsistency. **E**lse (no partition), choose between **L**atency and **C**onsistency. CAP only covers partition behavior; PACELC includes the normal case, where you still trade latency for consistency.

> **Front**: What are the three types of constraints, with one example of each?
>
> **Back**:
>
> - **Hard**: regulatory, physical, contractual (e.g., GDPR data residency)
> - **Soft**: organizational, budgetary (e.g., team lacks Rust experience)
> - **Self-imposed**: prior decisions / preferences (e.g., "we'll stay on AWS")

### Distributed Systems (Module 02)

> **Front**: What does the FLP impossibility result state?
>
> **Back**: In an _asynchronous_ distributed system with even one faulty node, no deterministic consensus algorithm can guarantee termination. Real-world protocols (Raft, Paxos) work _in practice_ by using timeouts, accepting that they may not terminate in adversarial scenarios.

> **Front**: For N=5 replicas, what R and W guarantee strong consistency?
>
> **Back**: R + W > N. So with N=5, configurations like R=3, W=3 or R=4, W=2 give strong consistency. R=3, W=3 also tolerates 2 failures for both reads and writes. R=1, W=5 gives the cheapest reads but no write availability under any single failure.

### Data at Scale (Module 03)

> **Front**: Why does an LSM-tree have higher write amplification than a B-tree?
>
> **Back**: B-trees update pages in place (write amplification ~1×). LSM-trees write to memtable, flush to SSTable, then _compact_ SSTables in background — each compaction rewrites data. Total write amplification typically 10-30× over the data's lifetime. The trade-off: LSM writes are sequential (fast), B-tree writes are random (slow).

### Architecture Styles (Module 04)

> **Front**: What's the difference between a "modular monolith" and a "tangled monolith"?
>
> **Back**: A modular monolith deploys as one unit but enforces strict internal module boundaries (typically via package structure + import linters). Modules can only talk through defined interfaces; each module owns its data. A tangled monolith has none of these boundaries — modules access each other's internals freely. The modular approach captures ~80% of microservices' benefit at ~20% of the operational cost.

### Event-Driven (Module 05)

> **Front**: What's the "dual-write problem" and what pattern solves it?
>
> **Back**: When an app writes to a database AND publishes to Kafka in separate operations, they can diverge: DB succeeds but Kafka fails (or vice versa) leaves the system inconsistent. Retries don't fix it (Two-Generals problem). The **Transactional Outbox pattern** solves it: write business data AND the event to an `outbox` table in the same database transaction; a relay process polls the outbox and publishes to Kafka.

### Reliability (Module 06)

> **Front**: When should you NOT use exponential backoff with retries?
>
> **Back**:
>
> - For non-idempotent operations (you'd risk double-writes)
> - For 4xx errors (it's your bug, retrying won't help)
> - When the circuit breaker is already open (don't retry at all — let it cool)
> - When the operation's latency is already part of someone else's timeout budget

> **Front**: What's a "bulkhead" in software architecture?
>
> **Back**: Separate resource pools (connection pools, thread pools) per downstream dependency, so one downstream's slowness can't exhaust shared resources and starve calls to healthy downstreams. Named after ship compartments that contain flooding to one section.

### Design at Scale (Module 07)

> **Front**: Twitter's timeline: when do you use fan-out-on-write vs fan-out-on-read?
>
> **Back**:
>
> - **Fan-out on write** (push) for most users: when they post, write to every follower's timeline cache. O(followers) write, O(1) read.
> - **Fan-out on read** (pull) for celebrities (>10K followers): just store the tweet; followers fetch + merge at read time.
> - **Hybrid in production**: combine both. Followers' reads merge pushed tweets + recently-fetched celebrity tweets.

### Architect's Craft (Module 08)

> **Front**: What's the test for whether to write an ADR?
>
> **Back**: Would a future engineer reading this in 18 months benefit from knowing _why_ the choice was made? If yes → write an ADR. If they'd just say "yeah, makes sense, moving on" → skip it. Don't ADR every variable name; do ADR every Type 1 (irreversible) decision.

---

## Full Deck Download

Download the complete deck: [`flashcards.csv`](./flashcards.csv)

Format: `Front,Back,Tags` (CSV, UTF-8). Compatible with Anki, Quizlet, Mochi, and any SRS tool that imports CSV.
