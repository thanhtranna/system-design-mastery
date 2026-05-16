# ADR-008: Caching Strategy — Cache-Aside for Catalogue, Write-Through for Cart

| Field | Value |
|---|---|
| **Status** | Accepted |
| **Date** | 2026-05-16 |
| **Author** | Thanh Tran |
| **Reviewers** | Backend Guild, Platform Team |
| **Review date** | 2027-05-16 |

---

## Context

Our e-commerce platform has two data access patterns with very different consistency requirements:

**Product catalogue** (10M products): reads are highly cacheable — products change infrequently (price updates ~100/sec, product description updates ~10/sec). Under peak load, we receive 50K catalogue reads/sec. Without caching, our PostgreSQL catalogue DB would need to handle all 50K reads/sec — well beyond its capacity (currently sized for 5K reads/sec writes + metadata queries).

**Shopping cart** (1M active carts): each cart change (add item, remove item, update quantity) must be immediately visible to the user on their next request, possibly from a different server. Cart data must not be lost. Carts are temporary (TTL 30 days) and small (~200 bytes per cart).

Redis is already deployed in the infrastructure (used for session management). We are choosing caching patterns for these two data types.

**Constraint**: team of 8 engineers; complexity budget is limited. We want the simplest pattern that satisfies correctness requirements for each data type.

---

## Decision

We will use **different caching strategies per data type**, based on consistency requirements:

### Product Catalogue — Cache-Aside (Lazy Loading) with TTL

```
Read path:
  1. Check Redis: GET product:{id}
  2. Cache hit → return (p99 ~1ms)
  3. Cache miss → query PostgreSQL → store in Redis (TTL 5 min) → return

Write path (price/product update):
  1. Write to PostgreSQL
  2. DEL product:{id} from Redis (invalidate, do not update cache)
  → next read repopulates cache from DB
```

**TTL: 5 minutes**. Rationale: product data can be 5 minutes stale with no user-visible impact. A price change may display the old price for up to 5 minutes — acceptable, and consistent with the "price is confirmed at checkout" user expectation.

**Invalidation on write**: when a product is updated, we delete the cache key (not update it). This avoids the race condition where a stale write overwrites a fresher one. The next read repopulates from the DB.

### Shopping Cart — Write-Through

```
Read path:
  1. Check Redis: HGETALL cart:{user_id}
  2. Cart is always in Redis for active sessions (TTL 30 days)
  3. Cache miss (cold start or TTL expired) → read from PostgreSQL → load into Redis

Write path (add/remove/update item):
  1. Write to Redis atomically (HSET / HINCRBY)
  2. Write to PostgreSQL asynchronously (via background job or immediate sync write)
  → Redis is the primary read store; PostgreSQL is the durable backup
```

**Why write-through for cart**: cart changes must be immediately consistent — the user who just added an item must see it on their next page load. Cache-aside for writes would require a cache invalidation + repopulation on every cart change, which is identical in effect to write-through but more code. Write-through simplifies the read path: the cart is always in Redis.

**Durability**: cart data is written to PostgreSQL synchronously on each change (not just async). Redis is the fast read layer; PostgreSQL is the source of truth. If Redis is flushed, carts are rebuilt from PostgreSQL.

---

## Consequences

### Positive

- **Catalogue reads at 50K req/sec from Redis**: p99 ~1ms vs ~20ms from PostgreSQL. Within our 200ms API SLO with substantial headroom.
- **Catalogue DB load reduced by ~95%**: from 50K reads/sec to ~2.5K reads/sec (cache miss rate ~5% at 5min TTL with warm cache).
- **Cart consistency guaranteed**: write-through ensures cart reads always reflect the most recent write, regardless of which server handles the request.
- **Simple implementation**: both patterns are well-understood, library-supported, and do not require complex coordination logic.

### Negative

- **Cache-aside allows stale catalogue data**: up to 5 minutes stale. Acceptable for product descriptions and images; requires a UX note at checkout ("price confirmed at checkout").
- **Write-through doubles write latency for cart operations**: every cart change writes to both Redis and PostgreSQL. At current scale (~10K cart writes/sec), both writes complete in < 5ms — acceptable.
- **Cache stampede risk for popular products**: when a high-traffic product's cache key expires, many concurrent requests may simultaneously miss and query PostgreSQL. Mitigation: probabilistic early expiration (refresh cache at TTL × random(0.8, 1.0)) for the top 1000 products by view count.

### Neutral

- **Two different patterns in the same Redis instance**: complexity is contained — each data type has its own key namespace (`product:*` vs `cart:*`). The Redis team manages both.
- **TTL-based invalidation is approximate**: products may be served stale for up to 5 minutes. This is a deliberate choice. Event-driven invalidation (exact) is available in v2 if staleness tolerance decreases.

---

## Alternatives Considered

### Write-Through for Catalogue

Update the cache on every product write (price change, description update).

**Rejected because**: the write path for catalogue updates goes through a separate admin service, batch import jobs, and vendor feeds — multiple writers to the same data. Keeping all writers in sync with cache updates adds coupling and failure modes. Cache-aside with TTL is simpler when there are multiple writers.

### Event-Driven Cache Invalidation (Pub/Sub) for Catalogue

On product update, publish a `product_updated` Kafka event. Cache invalidation workers consume the event and delete Redis keys.

**Rejected for v1 because**: adds Kafka consumer complexity, a new service, and failure modes (what if the invalidation worker is lagging?). The 5-minute TTL provides acceptable staleness today. If catalogue update frequency increases or staleness tolerance decreases, this is the v2 approach.

### Read-Through (Redis as the primary, DB as backup, managed by Redis module)

Use RedisJSON or a custom Redis module to manage the read-through logic server-side.

**Rejected because**: Redis modules add operational complexity and are not natively supported by all managed Redis providers (AWS ElastiCache has limited module support). Standard Redis commands with application-layer logic are more portable and debuggable.

### Cache-Aside for Cart

Treat cart data the same as catalogue — read from cache, invalidate on write, repopulate on miss.

**Rejected because**: the cache miss on a cart read would require a DB query on every cart page load for users with an expired or evicted cart. Cart operations are user-interactive (< 200ms latency expectation). Write-through eliminates the per-read DB dependency for active carts.

---

## Open Questions

- Should we implement probabilistic early expiration for hot product pages? (Current decision: yes, for the top 1000 products by traffic; implemented in the next sprint.)
- Cart TTL of 30 days — is this the right balance between abandoned cart recovery and Redis memory? (At 1M active carts × 200 bytes = 200 MB in Redis. 30 days is fine at this scale.)
- When write-through cart writes conflict (two devices writing simultaneously), which wins? (Last-write-wins via Redis HSET. This is acceptable — simultaneous cart edits from two devices are rare and the last write is the most recent user intent.)

---

## References

- ADR-001: PostgreSQL as OLTP store — the source of truth for both catalogue and cart data
- ADR-004: Meilisearch — catalogue search index has a separate invalidation path (not this ADR)
- AWS ElastiCache Redis documentation — cache-aside pattern
- "Database Reliability Engineering" by Laine Campbell — Chapter 7 (caching strategies)
