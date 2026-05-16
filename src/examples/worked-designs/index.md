# Worked Designs

Full architecture documents for systems you'd be asked to design in interviews or in real work. Each one applies the patterns from Modules 01–07 end-to-end: requirements → estimates → high-level → deep-dives → ADRs.

These are **gold standard examples**. Use them as:

- **Reference** when writing your own architecture documents
- **Critique target** — try to find what they got wrong or could improve
- **Interview prep** — read the structure, then try a similar prompt under time pressure

## Available

| Design                                                    | Scale                               | Maps to                                      |
| --------------------------------------------------------- | ----------------------------------- | -------------------------------------------- |
| [PropertyHub Real Estate Search](./real-estate-search.md) | 5M listings, 5K QPS peak            | Module 03, 04, 05, 06                        |
| [URL Shortener](./url-shortener.md)                       | 50M URLs/month, 100:1 R:W           | Module 01, 03 — the "small system" case      |
| [Real-Time Chat (WhatsApp-style)](./chat-system.md)       | 1B DAU, 5M msg/sec peak             | Module 02, 03, 05, 06, 07                    |
| [Payment Gateway](./payment-gateway.md)                   | 50M tx/day, multi-acquirer          | Module 05, 06, 08 — correctness over latency |
| [Social Media Feed](./social-media-feed.md)               | 500M DAU, 5M posts/day              | Module 03, 05, 07 — fanout strategy          |
| [Video Streaming Platform](./video-streaming.md)          | 500 hrs uploaded/min, 1B views/day  | Module 03, 04, 06 — pipeline design          |
| [Ride-Sharing Platform](./ride-sharing.md)                | 5M rides/day, sub-5s match          | Module 03, 06 — geospatial at scale          |
| [Notification Service](./notification-service.md)         | 50M notifications/day               | Module 05, 06 — idempotency and reliability  |
| [Distributed Rate Limiter](./rate-limiter.md)             | 100K req/sec, multi-region          | Module 06, 07 — approximate counting         |

## Reading Order

If you're going through these in sequence:

1. **URL Shortener** first — small scale, teaches the discipline of estimating before designing. Beware over-engineering.
2. **Real Estate Search** — medium scale, multi-storage system (PostgreSQL + OpenSearch + Redis). Polyglot persistence in action.
3. **Chat System** — large scale (1B DAU). Stateful gateways, per-conversation ordering, the trade-offs that emerge at scale.
4. **Payment Gateway** — small/medium scale, but correctness-critical. Different quality attribute ranking entirely.
5. **Social Media Feed** — 500M DAU. Fanout strategy is the core decision; push vs pull vs hybrid.
6. **Video Streaming** — pipeline design at scale. Upload → transcode → CDN delivery.
7. **Ride-Sharing** — real-time geospatial matching. Location tracking and race conditions.
8. **Notification Service** — idempotency and multi-channel fan-out. Deceptively hard to do correctly.
9. **Rate Limiter** — approximate counting. Exact consistency is impossible; learn where to accept the trade-off.

The contrast between all nine is the lesson: **scale alone doesn't determine complexity. The ranked quality attributes do.**

## What "Worked" Means

Each design includes:

1. **Executive summary** (1 paragraph)
2. **Requirements** — functional, non-functional, what's NOT in scope
3. **Capacity estimates** with traceable assumptions
4. **C4 Context + Container diagrams** (Mermaid)
5. **Component drill-downs** for the 2–3 most critical pieces
6. **Sequence diagrams** for important flows
7. **Trade-off analysis** with at least 3 explicit decisions
8. **Failure modes** and how the design degrades
9. **Migration / rollout strategy** if applicable
10. **Open questions** — what's deliberately deferred
11. **ADR set** — links to the supporting decisions

## Reading Strategy

First pass: skim the structure, note the headers.

Second pass: focus on Section 7 (trade-offs) and Section 11 (ADRs). These are where the architecture _thinking_ lives.

Third pass: critique. Where would you push back? What's missing? What assumed scale that may not hold?

If after three passes you couldn't have produced something like this, you have more practice ahead. That's fine. The structure repeats; the muscle builds.
