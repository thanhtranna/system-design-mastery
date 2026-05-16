# ADR-003: Adopt Modular Monolith architecture for the PropertyHub 2.0 rebuild — explicitly NOT microservices

## Status

**Accepted** — 2026-01-22
**Authors**: Thanh Tran, Linh Pham (Engineering Manager)
**Reviewers**: CTO, all backend leads, Platform team
**Supersedes**: discussion document `2025-12-spike-microservices.md`
**Review date**: 2027-01-22

## Context

In Q4 2025, the PropertyHub engineering team conducted a 4-week design spike on whether to rebuild the aging PHP monolith as:

- **(A)** A new monolith (one deployable, modern stack)
- **(B)** A modular monolith (one deployable, but with strict internal boundaries)
- **(C)** Microservices (multiple deployables, independent teams)

The discussion was heavily influenced by:

- Two engineers had recently joined from companies running microservices and brought enthusiasm
- A widely-shared blog post from a competitor describing their microservices migration
- Leadership pressure to "modernize" before fundraising

We collected hard data on the current state:

| Metric                            | Current (PHP monolith)                                          |
| --------------------------------- | --------------------------------------------------------------- |
| Engineering team size             | 25 (planned to grow to 40 in 12 months)                         |
| Number of business domains        | 6 (Listings, Agents, Search, Notifications, Reporting, Billing) |
| Current deploy frequency          | 1–2 per day, with rollbacks ~10% of the time                    |
| Average build+test time           | 18 minutes                                                      |
| Active oncall rotations           | 1 (full backend team shares load)                               |
| Microservices experience on team  | 3 of 25 engineers                                               |
| Kubernetes operational experience | 1 engineer (recently hired)                                     |
| DevOps platform team size         | 2                                                               |

Constraints:

- **Hard**: must ship MVP of rebuild within 9 months. Cannot stall product roadmap.
- **Soft**: avoid adding new on-call rotations without dedicated platform team funding (not on roadmap).
- **Self-imposed**: at most one paradigm shift at a time. We are already changing language (PHP → Go), deployment (VMs → containers), and storage (MySQL → PostgreSQL).

Conway's Law analysis: our 25 engineers are currently organized as a single backend team that collaboratively owns all domains. Splitting into autonomous teams requires explicit reorg + headcount investment that is not funded.

## Decision

**We will rebuild PropertyHub as a modular monolith in Go.** This means:

- One deployable artifact, one binary
- Strict internal module boundaries enforced by package structure and import linters
- Each module owns its data tables; cross-module access only through defined interfaces
- One database (PostgreSQL) but logical schema separation per module
- Single deploy pipeline, single oncall rotation, single observability stack

We will **not** split into microservices in 2026.

## Consequences

### Positive

- **Aligned with team capability.** 25 engineers, mostly without microservices experience, can iterate fast on a single codebase.
- **Operational simplicity.** One deploy, one oncall, one observability stack. Our 2-person DevOps team can handle this comfortably.
- **Strong consistency by default.** Cross-module operations use ACID transactions. No saga complexity for v1.
- **Fast development loop.** Local dev: spin up one binary + Postgres. No service mesh, no distributed tracing required for basic debugging.
- **Refactoring across module boundaries is cheap.** Compile-time checks, single PR. We'll learn the right boundaries before paying the cost of network boundaries.
- **Faster to MVP.** Estimated 9 months for monolith vs 15+ months for microservices (per spike).
- **Future-flexible.** If we later need to extract a service (e.g., search at very high scale), modular boundaries make extraction cheap. We've front-loaded the boundary discipline without paying the deployment tax.

### Negative / Costs

- **Single point of failure for deploys.** A bug in one module can block deploys for all. Mitigated by: feature flags, canary deploys, automated rollback.
- **Independent scaling not possible.** If Search becomes CPU-bound, we have to scale the whole monolith. Acceptable trade-off until we have data showing this is the bottleneck (we don't).
- **Long-term ceiling.** If we grow to 100+ engineers, modular monolith will likely need to split. We accept that this ADR has a ~3-year horizon, not 10.
- **Cultural risk**: some engineers were excited about microservices; we may need to explicitly recognize that the _concepts_ of microservices (strict boundaries, clear contracts) are valuable and being applied — just not the deployment topology.
- **Discipline cost**: enforcing module boundaries requires tooling and code review attention. If discipline slips, we'll regress to a tangled monolith.

### Neutral / Trade-offs accepted

- **No technology heterogeneity.** Everything is Go. We see this as a feature (consistency, hiring) not a bug.
- **No team autonomy via service ownership.** Team coordination remains a thing. Will partly address through clear module ownership.

## Alternatives Considered

### Option C: Microservices from day one

What the loud minority of the team wanted.

**Why rejected**:

| Concern                | Microservices reality                | Our situation                                             |
| ---------------------- | ------------------------------------ | --------------------------------------------------------- |
| Independent deploys    | Yes, but requires CI/CD per service  | We have 2 DevOps engineers; can't run 6 pipelines well    |
| Independent scaling    | Yes, but rarely the bottleneck early | Our peak load is modest (~5K QPS)                         |
| Team autonomy          | Yes, with autonomous teams           | Our 25 engineers aren't organized into autonomous teams   |
| Technology flexibility | Yes, polyglot                        | We're already standardizing on Go for the rebuild         |
| Cascade resilience     | Sometimes, with correct patterns     | We'd be implementing patterns we don't have time to learn |
| Operational cost       | High; requires platform investment   | We don't have the budget for a platform team              |

The honest summary: **microservices solve organizational problems we don't have**, at the cost of operational problems we can't afford to take on.

### Option A: Plain new monolith (no module discipline)

The "just rewrite it in Go" approach.

**Why rejected**:

- We've seen where this leads (our current PHP monolith). After 4 years it became hard to change.
- Modular boundaries are cheap to add up-front, expensive to retrofit later.
- The bar of "module-per-package + import linter" costs little.

### Option D: Hybrid — start as modular monolith, extract specific high-load services

This is what we'll likely do _eventually_, but not now. Naming it as a future option rather than current decision.

## What "Modular Monolith" Means in Practice

To prevent the cultural drift back toward a "monolith without modules," we commit to:

1. **Module-per-package** in Go. Each module is `internal/<domain>/`.
2. **Public API per module**: each module exposes a typed interface; other modules import only this interface.
3. **Data ownership**: each module's tables prefixed `<module>_*` (e.g., `listings_*`, `agents_*`). Cross-module reads go through the owning module's interface, not direct SQL.
4. **Import linter (`go-arch-lint`)**: CI fails if a module's internal types are imported elsewhere.
5. **Per-module owners**: each module has 1–3 named owners; cross-module changes require their review.

## Open Questions

- **At what signal do we start extracting services?** Tentative criteria: a module exceeds 30 active engineers contributing to it monthly, OR has independent scaling characteristics severely mismatched from the rest, OR has a strong availability SLO mismatch.
- **How do we handle cross-module data joins for reporting?** Likely: a separate read replica + a read-only "reporting" module that can cross-cut. Will detail when reporting becomes the priority.
- **What's our test boundary strategy?** Each module needs unit tests; do we need module-level integration tests? Decision: yes, with a shared test harness.

## References

- Spike document: `2025-12-spike-microservices.md` (40 pages with team interviews and benchmarks)
- Simon Brown, "Modular Monoliths" (https://www.codingthearchitect.com/2018/07/22/modular_monoliths.html)
- Sam Newman, "Monolith to Microservices" (book) — especially Chapters 1–2 on when NOT to do microservices
- Course material: Module 04 — Architecture Styles
- Shopify engineering: "Deconstructing the Monolith" (parallel internal blog post)

## Decision Record History

| Date       | Status   | Notes                                                                   |
| ---------- | -------- | ----------------------------------------------------------------------- |
| 2025-12-15 | Proposed | First draft after spike completion                                      |
| 2026-01-12 | Revised  | Added explicit "what modular means" section per CTO concern about drift |
| 2026-01-22 | Accepted | Approved at all-hands engineering meeting                               |

## Postscript (added 2026-04-22, 3 months in)

After 3 months of building under this ADR:

- Module boundaries are holding. `go-arch-lint` has caught 14 cross-module violations in PRs.
- Two engineers initially skeptical of "monolith" decision now advocate for the modular approach internally.
- Build times remain under 4 minutes despite 3× LOC growth.
- One emerging concern: Notifications module is doing heavy outbound HTTP and slowing test runs. May need to mock more aggressively or extract Notifications first. Tracking in `TICKET-2026-MQ-014`.
