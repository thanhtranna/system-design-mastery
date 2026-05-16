# RFC Template

An RFC (Request for Comment) is a **proposal document** written before a significant technical decision is made. It is the pre-cursor to an ADR: where an ADR records a decision already taken, an RFC invites feedback on a proposed direction.

Use an RFC when the decision is non-trivial, affects multiple teams or systems, or when you want structured feedback before committing to an approach.

---

## When to Write an RFC vs an ADR

| Write an RFC when... | Write an ADR when... |
|---|---|
| The decision is not yet made | The decision is final (or already implemented) |
| You want structured feedback from reviewers | You want to record context for future readers |
| The change affects multiple teams | The change is within your team's scope |
| You're proposing a new pattern or architecture | You're applying an existing pattern |
| The trade-offs are genuinely unclear | The trade-offs are understood; you're documenting the choice |

**RFC → ADR lifecycle**: an accepted RFC typically becomes an ADR. The RFC captures the deliberation; the ADR records the decision. Link them together.

---

## RFC Template

Copy the section below. Delete this header before submitting.

---

```markdown
# RFC-NNN: [Short descriptive title]

**Status**: Draft | In Review | Accepted | Rejected | Withdrawn
**Author(s)**: [Your name(s)]
**Reviewers**: [Who you're asking for feedback]
**Created**: [YYYY-MM-DD]
**Last updated**: [YYYY-MM-DD]
**Review deadline**: [YYYY-MM-DD — give reviewers a clear window]

---

## Problem Statement

[2–4 paragraphs. What is the problem? Why does it matter now? Who is affected?

Be concrete: name the systems, teams, users, and scale. Avoid "we have a problem with X" — say "we receive Y requests/sec to service Z, which causes W every N days."]

## Proposed Solution

[Describe what you are proposing to do. Be specific:
- What changes?
- What does NOT change?
- What is the implementation order?

Include diagrams where they clarify. Mermaid works in mdBook.]

## Alternatives Considered

[For each alternative, state what it is and WHY you rejected it. "Less suitable" is not a reason.

Minimum 2 alternatives. "Do nothing" is always a valid alternative and should be listed.]

### Alternative A: [Name]

[Description]

**Rejected because**: [concrete reason — latency, complexity, cost, team expertise, etc.]

### Alternative B: [Name]

[Description]

**Rejected because**: [concrete reason]

### Alternative C: Do nothing

[What happens if we don't make this change?]

**Rejected because**: [why the status quo is unacceptable]

## Impact Analysis

### Systems affected

| System | Change required | Owner | Effort estimate |
|---|---|---|---|
| [System A] | [What changes] | [Team] | [S/M/L/XL] |

### Breaking changes

[Are there API changes? Data format changes? Will existing clients break?
If yes: migration path, backwards compatibility window, deprecation plan.]

### Performance impact

[Expected change in latency, throughput, resource usage. If unknown, say so and describe how you'll measure.]

### Security impact

[New attack surface? Changes to auth model? Data exposure? If none, say "No security impact."]

### Cost impact

[Infrastructure cost change, if any. Order of magnitude is fine if exact is unknown.]

## Open Questions

[List what is NOT decided and needs input from reviewers. This section is why you're writing an RFC.

Format: "Q: [question]" — reviewers can answer inline.]

Q: [First open question]

Q: [Second open question]

## Success Criteria

[How will you know this worked? Specific, measurable.

Examples:
- p99 latency of endpoint X drops from Yms to Zms
- Incident rate for failure mode W drops to near zero
- Team B can deploy changes to system A without coordinating with team C]

## Timeline

| Milestone | Target date | Notes |
|---|---|---|
| RFC accepted | [date] | |
| Implementation starts | [date] | |
| First deployment (staging) | [date] | |
| Production rollout begins | [date] | |
| Full rollout / feature flag removed | [date] | |

## Reviewer Guidance

[Tell reviewers what kind of feedback is most valuable. Are you looking for:
- Feedback on the overall approach?
- Specific expertise on technology X?
- Validation that the impact analysis is complete?
- Challenge to assumptions?]

Please comment inline on sections you have concerns with. Add a summary comment with your overall stance: **Approve / Approve with changes / Request changes / Block**.

**Block** should be used sparingly and must include a specific, addressable objection.
```

---

## Example: Filled RFC Snippet

```markdown
# RFC-007: Replace synchronous user lookup with async pre-fetch cache

**Status**: In Review
**Author(s)**: Thanh Tran
**Reviewers**: Backend Guild, Platform Team
**Created**: 2026-04-10
**Review deadline**: 2026-04-17

---

## Problem Statement

The `/api/feed` endpoint performs a synchronous user lookup (database call) for every
request to resolve the requesting user's permissions. At 5,000 req/sec, this generates
5,000 DB reads/sec against the users table — 40% of our DB query budget for one table.

The user profile data (permissions, preferences) changes at most once per hour per user.
We are paying per-request DB cost for data that is effectively static on a minute-by-minute basis.

This is causing p99 latency on `/api/feed` to be 340ms (SLO: 200ms) because the users
table query is the bottleneck (p99 = 290ms under current load).

## Proposed Solution

Pre-fetch user data into Redis at login time and on user profile change events.
The feed endpoint reads from Redis (p99 ~1ms) instead of PostgreSQL.

Cache TTL: 15 minutes. On profile change: invalidate via Pub/Sub message.

[diagram of new flow]

## Open Questions

Q: What is the acceptable staleness window? 15 minutes was chosen; is that acceptable
   for permission changes? (If a user's account is suspended, they could still have
   cache hits for up to 15 minutes.)

Q: Should we use Redis Cluster or a single Redis instance? The dataset is ~10M users ×
   500 bytes = 5 GB. Fits on a single instance, but no HA.
```

---

## After the RFC is Accepted

- [ ] Status changed to "Accepted"
- [ ] RFC linked from the team wiki / decision log
- [ ] Key decisions from the RFC written up as an ADR (the ADR is the permanent record)
- [ ] Open questions resolved and documented in the RFC before it is archived
- [ ] If the RFC was rejected: status set to "Rejected", reason documented, future readers can reference it to understand why this was considered and not pursued

## After the RFC is Rejected

Rejected RFCs are valuable. They prevent re-litigating the same proposals. Document:
- Why it was rejected
- What would need to change for it to be worth reconsidering
- Who made the decision and when

A rejected RFC that is discoverable is better than a rejected RFC that is deleted.
