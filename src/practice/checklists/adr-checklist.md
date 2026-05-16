# ADR Submission Checklist

Use this **before submitting** any ADR for review. ADRs that hit these bars get accepted faster, age better, and survive team changes.

---

## Structure

- [ ] **Title** is the decision, not a question
  - ✅ "Use PostgreSQL 16 for OLTP store"
  - ❌ "PostgreSQL vs MySQL?"
  - ❌ "Database decision"
- [ ] **Status** field set (Proposed / Accepted / Deprecated / Superseded)
- [ ] **Author** and **reviewers** named
- [ ] **Review date** set (most decisions should be revisited)

## Context

- [ ] Sets the stage in 2–4 paragraphs (no longer)
- [ ] Names the **forces** acting on the decision: scale, deadline, team competence, business pressure
- [ ] Separates hard constraints from soft and self-imposed
- [ ] **Quantified** where possible (X requests/sec, Y engineers, $Z budget)
- [ ] Names stakeholders affected
- [ ] No code in this section

## Decision

- [ ] One short sentence, active voice
  - ✅ "We will use Amazon RDS for PostgreSQL 16 as the primary OLTP store."
  - ❌ "After consideration, the team has decided to perhaps adopt PostgreSQL."
- [ ] Implementation details kept brief — major bullets only

## Consequences

- [ ] **Positive** section: concrete benefits, not platitudes
  - ✅ "Reduces per-request DB latency from ~80ms to ~12ms based on benchmark"
  - ❌ "Better performance"
- [ ] **Negative / Costs** section: honest, not buried
- [ ] **Neutral / Trade-offs accepted** section: things that aren't strictly bad but are deliberate choices
- [ ] Capacity / cost numbers where applicable

## Alternatives Considered

- [ ] At least 2 alternatives named (even if one is "do nothing")
- [ ] Each alternative has a **concrete** rejection reason
  - ✅ "p99 latency 1.8× higher in our benchmark with realistic data volume"
  - ❌ "Less suitable"
- [ ] For unconventional choices: what _would_ have made the alternative right? When would we revisit?

## Open Questions

- [ ] What's NOT decided yet
- [ ] When/how will those open questions get resolved?
- [ ] Naming uncertainty is a strength, not a weakness

## References

- [ ] Benchmarks, prior art, vendor docs
- [ ] Related ADRs (especially if superseding)
- [ ] Incident reports (if decision was incident-driven)
- [ ] External articles / papers where relevant

---

## Quality Tests

### The "New Hire" Test

Hand the ADR to a senior engineer who joined yesterday. Can they:

- Read it in under 5 minutes?
- Explain the decision back to you in 2 minutes?
- Name one alternative and why it was rejected?

If no to any → **revise before submitting**.

### The "Why Didn't We" Test

After reading, would a reasonable reader still ask "but why didn't we just do X?" If yes → that "X" should be in Alternatives Considered.

### The "Future Self" Test

In 18 months, you (or someone else) will read this ADR while debugging a problem caused by this decision. Will the ADR explain _why_ the trade-off was acceptable at the time?

If no → add context.

### The "Bus Factor" Test

If everyone who authored this ADR left the company, could a new hire understand the system and continue making good decisions in this area?

If no → the _context_ is undocumented. Fix it.

---

## Common Failure Modes

| Pattern                                               | Fix                                                   |
| ----------------------------------------------------- | ----------------------------------------------------- |
| Title is the problem, not the decision                | Rewrite as active-voice statement                     |
| All upsides, no downsides                             | Force at least 2 honest downsides                     |
| Alternatives section is "we considered other options" | Name each, with concrete rejection reason             |
| ADR is 12 pages                                       | Cut to 3. Anything longer should be an RFC            |
| Open questions section absent                         | Either it's perfect (unlikely) or you're hiding doubt |
| No review date                                        | Add one. Decisions decay                              |
| Buried conclusion                                     | Move Decision to top after Context                    |

---

## When NOT to Write an ADR

Not every choice needs an ADR. Skip if:

- It's reversible in under a day
- The team doesn't care which way you decide
- You're following an existing pattern (link to the originating ADR)
- It's an implementation detail (variable name, log format)

**Rule**: if a future engineer would benefit from knowing _why_, write an ADR. If they'd just say "yeah, makes sense, moving on" — skip.

---

## After Acceptance

- [ ] Status changed to "Accepted"
- [ ] Linked from team docs, runbook, or onboarding guide
- [ ] Calendar reminder set for review date
- [ ] If supersedes an earlier ADR: that ADR's status updated to "Superseded by X"
