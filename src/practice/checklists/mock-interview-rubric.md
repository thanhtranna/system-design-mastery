# Mock Interview Rubric

Use this after every 45–60 minute system design mock. Score honestly. Track trends over 6 months.

## Scoring Scale

| Score | Label     | Description                                        |
| ----- | --------- | -------------------------------------------------- |
| **1** | Junior    | Skipped this dimension entirely or got it wrong    |
| **2** | Mid       | Mentioned but shallow                              |
| **3** | Senior    | Solid execution at expected level                  |
| **4** | Staff     | Strong, with original insight or nuance            |
| **5** | Architect | Exceptional. Could lead the design at a senior org |

**Target by end of Module 08**: average 4+ across all dimensions. Occasional 5s.

---

## The Rubric (10 dimensions)

### 1. Clarifying Requirements

| Score | Behavior                                                                          |
| ----- | --------------------------------------------------------------------------------- |
| 1     | Dives in drawing boxes within 30 seconds                                          |
| 2     | Asks "what's the scale?" then proceeds                                            |
| 3     | Asks scale, read:write ratio, latency target                                      |
| 4     | Asks scale, SLOs, functional scope, **what's NOT in scope**                       |
| 5     | Explicitly frames the system's purpose, ranks quality attributes before designing |

### 2. Back-of-Envelope Estimation

| Score | Behavior                                                                               |
| ----- | -------------------------------------------------------------------------------------- |
| 1     | "It'll be a lot of traffic"                                                            |
| 2     | Picks a DAU number, no derivation                                                      |
| 3     | DAU → QPS with stated assumptions                                                      |
| 4     | QPS, storage, bandwidth, all derived from a few inputs; identifies dominant constraint |
| 5     | Numbers anchor every subsequent decision; estimates revised as design evolves          |

### 3. Top-Down Structure

| Score | Behavior                                                                                    |
| ----- | ------------------------------------------------------------------------------------------- |
| 1     | Random boxes appear; no clear order                                                         |
| 2     | Eventually settles into layers, with backtracking                                           |
| 3     | Linear: clients → API → services → storage                                                  |
| 4     | Clear C4-style layering; abstracts before details                                           |
| 5     | Each layer earns its complexity; explicitly justifies why something is a separate component |

### 4. Trade-off Articulation

| Score | Behavior                                                                                                         |
| ----- | ---------------------------------------------------------------------------------------------------------------- |
| 1     | "Best practice is X" with no comparison                                                                          |
| 2     | Acknowledges one alternative ("could also use Y but...")                                                         |
| 3     | Names at least 2 explicit trade-offs (latency vs consistency, cost vs availability)                              |
| 4     | Surfaces 3+ trade-offs, each grounded in a quality attribute                                                     |
| 5     | Frames every key decision as a trade-off with explicit rationale; pushes back on interviewer's leading questions |

### 5. Storage / Data Modeling

| Score | Behavior                                                                                                |
| ----- | ------------------------------------------------------------------------------------------------------- |
| 1     | "We'll use a database"                                                                                  |
| 2     | Picks SQL or NoSQL with one-word justification                                                          |
| 3     | Picks specific store with access-pattern reasoning                                                      |
| 4     | Discusses partitioning, replication, indexes appropriately                                              |
| 5     | Discusses isolation level, hot keys, capacity planning, evolution; chooses polyglot only when justified |

### 6. Reliability / Failure Modes

| Score | Behavior                                                                                     |
| ----- | -------------------------------------------------------------------------------------------- |
| 1     | Happy path only                                                                              |
| 2     | Mentions retries                                                                             |
| 3     | Mentions retries, timeouts, replication                                                      |
| 4     | Identifies blast radii; applies circuit breakers, bulkheads, rate limiting where appropriate |
| 5     | Predicts which component fails first at 10× scale; discusses chaos engineering and recovery  |

### 7. Consistency Model

| Score | Behavior                                                                            |
| ----- | ----------------------------------------------------------------------------------- |
| 1     | "It'll be consistent"                                                               |
| 2     | Mentions CAP but as "pick 2 of 3"                                                   |
| 3     | Picks AP or CP and justifies                                                        |
| 4     | Specifies consistency model (eventual, causal, strong) per use case                 |
| 5     | Different parts of system have different models, justified by business requirements |

### 8. Scalability Discussion

| Score | Behavior                                                                              |
| ----- | ------------------------------------------------------------------------------------- |
| 1     | "We'll scale horizontally"                                                            |
| 2     | Mentions sharding                                                                     |
| 3     | Discusses partitioning strategy with hot-key awareness                                |
| 4     | Discusses scaling reads vs writes independently; caching strategy                     |
| 5     | Identifies bottlenecks at multiple scales (1x, 10x, 100x); plans graceful degradation |

### 9. Out-of-Scope Discipline

| Score | Behavior                                                                    |
| ----- | --------------------------------------------------------------------------- |
| 1     | Tries to design everything                                                  |
| 2     | Skips topics when running low on time                                       |
| 3     | Says "I'll skip X for now"                                                  |
| 4     | Explicitly lists what's NOT in v1 with reasons                              |
| 5     | Pushes back on scope creep from interviewer; recommends incremental rollout |

### 10. Communication & Time Management

| Score | Behavior                                                                  |
| ----- | ------------------------------------------------------------------------- |
| 1     | Rambling, jumps around, runs out of time                                  |
| 2     | Reasonably linear but uneven time allocation                              |
| 3     | Follows the 5-phase framework; finishes on time                           |
| 4     | Adapts to interviewer cues (when to go deeper, when to move on)           |
| 5     | Comes across like a senior peer collaborating, not a candidate performing |

---

## Score Sheet

Print this. Fill out after every mock.

```
Date: __________   Prompt: __________________________
Duration: ____ min   Interviewer: __________

Dimension                       Score (1-5)    Notes
1. Clarifying requirements      [ ]
2. Back-of-envelope estimation  [ ]
3. Top-down structure           [ ]
4. Trade-off articulation       [ ]
5. Storage / data modeling      [ ]
6. Reliability / failure modes  [ ]
7. Consistency model            [ ]
8. Scalability discussion       [ ]
9. Out-of-scope discipline      [ ]
10. Communication / timing      [ ]

Average: ____ / 5.0

Top 2 strengths:
-
-

Top 2 things to work on:
-
-

One specific thing I'll change next time:

```

---

## Tracking Over Time

Keep a spreadsheet. Plot average score over weeks. The graph should bend upward. If it plateaus, your gap is one specific dimension — drill that.

| Week | Average | Weakest Dim | Action                             |
| ---- | ------- | ----------- | ---------------------------------- |
| 1    | 2.4     | Estimation  | Practice 5 estimate-only prompts   |
| 2    | 2.8     | Trade-offs  | Re-read Module 04 trade-off tables |
| ...  | ...     | ...         | ...                                |

By Module 08, you should be consistently 4+ across the board. That's an "architect-level" mock — comparable to what you'd see in a real principal/staff interview.
