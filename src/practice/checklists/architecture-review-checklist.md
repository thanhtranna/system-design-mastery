# Architecture Review Checklist

Use this **before** any architecture review you're presenting or running. Most review meetings fail because the doc wasn't ready, not because the design was wrong. This checklist makes the doc review-ready.

If you can't check 90%+ of these, the meeting will be a fishing expedition. Postpone the meeting and use this list as your TODO.

---

## Before Sharing the Doc

### Problem Framing

- [ ] One sentence describing what we're building, in business language (no jargon)
- [ ] One paragraph describing why this matters now (the trigger, the pain)
- [ ] Explicit list of what we're **NOT** building (non-goals)
- [ ] At least one alternative considered and explicitly rejected

### Requirements

- [ ] Functional requirements listed as bullets
- [ ] Non-functional / quality attributes ranked (top 3)
- [ ] Hard constraints separated from soft constraints
- [ ] Capacity estimates with assumptions stated (not "it'll scale")

### Design

- [ ] At least one diagram, at the right C4 level for the audience
- [ ] Every box labeled with what it IS and language/runtime
- [ ] Every arrow labeled with protocol/direction
- [ ] At least one sequence diagram for the most important flow

### Trade-offs

- [ ] At least 3 explicit trade-offs surfaced
- [ ] For each: option A, option B, why we chose
- [ ] At least one decision called "Type 1" (irreversible)

### Reliability

- [ ] Failure modes table: what fails → blast radius → mitigation
- [ ] Recovery plan: how do we degrade vs. fail
- [ ] Capacity headroom: how many ×s of growth before redesign

### Migration / Rollout (if applicable)

- [ ] Phases with concrete checkpoints
- [ ] Rollback plan for each phase
- [ ] Success metrics per phase

### Open Questions

- [ ] Explicit list of what we don't know yet
- [ ] What needs to be decided before code
- [ ] What can be decided as we build

### Communication

- [ ] Title is a decision, not a question ("Use X" not "X or Y?")
- [ ] Conclusion at the top — reader knows the recommendation in 30 seconds
- [ ] Length: 2-10 pages depending on scope; not 25
- [ ] Skim-able: headers, bullets, no walls of text

---

## Before the Meeting

- [ ] Doc sent at least 48 hours in advance
- [ ] Reviewers given async comment window
- [ ] Pre-wired 1:1 with at least 2 key stakeholders (no group meeting surprises)
- [ ] Meeting agenda is clear: deciding, exploring, or reviewing?
- [ ] Open questions surfaced upfront — not buried

---

## During the Meeting

- [ ] Start with the decision being requested
- [ ] Confirm everyone read it (if not, **postpone**)
- [ ] Focus on disagreement and open questions, not walkthrough
- [ ] Time-box each discussion ("5 more minutes on this, then move on")
- [ ] Capture decisions live in the doc
- [ ] End with concrete next steps and owners

---

## After the Meeting

- [ ] Decision documented in the doc within 24 hours
- [ ] ADR written if decision is Type 1
- [ ] Action items distributed with owners and dates
- [ ] Folks who couldn't attend get an async summary
- [ ] Calendar invite for follow-up if decision wasn't reached

---

## Common Failure Modes (avoid these)

| Anti-pattern                          | How to fix                                    |
| ------------------------------------- | --------------------------------------------- |
| "Let's discuss further" (no decision) | Force a decision criterion before the meeting |
| Walking through the entire doc        | Send async-first; assume people read          |
| Architect monopolizes airtime         | Plan questions for the room; pause for input  |
| Decisions made by loudest voice       | Use a structured technique (dot voting, RACI) |
| Same questions every review           | Document FAQ in the doc itself                |
| No follow-through                     | Decisions without owners = no decisions       |

---

## The 30-Second Test

Hand this doc to a senior engineer who joined the company yesterday. Set a 30-second timer. Ask: "What is this proposing and why?"

If they can answer correctly, the doc is review-ready. If they can't — the executive summary needs work. **Revise before reviewing.**
