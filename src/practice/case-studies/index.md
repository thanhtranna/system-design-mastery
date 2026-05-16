# Failure Case Studies

Real production outages. Real lessons. **More learning per page than any textbook.**

Each case study includes:

- Timeline of what happened
- Root cause analysis
- Architecture before and after (where public)
- Lessons mapped to course modules
- Discussion questions for self-study or group

## Reading Schedule

One a week, alongside your current module. Each takes ~30 minutes to read + reflect.

| Case                                                                  | Year | Primary Module | Why read it                                      |
| --------------------------------------------------------------------- | ---- | -------------- | ------------------------------------------------ |
| [Knight Capital — $440M in 45 minutes](./knight-capital-2012.md)      | 2012 | M01, M04       | Quality attributes you don't name will kill you  |
| [GitHub — 24-hour MySQL outage](./github-2018.md)                     | 2018 | M02, M06       | Failover automation and split-brain              |
| [Roblox — 73 hours down](./roblox-2021.md)                            | 2021 | M02, M06       | Consensus systems have capacity limits           |
| [AWS S3 — typo takes down half the internet](./aws-s3-2017.md)        | 2017 | M06            | Blast radius and tooling safeguards              |
| [Discord → ScyllaDB migration](./discord-2022.md)                     | 2022 | M03            | When to migrate; how to migrate at scale         |
| [Facebook BGP — 6-hour global outage](./facebook-bgp-2021.md)         | 2021 | M02, M06, M08  | Management plane must survive data plane failure |
| [Cloudflare — regex that stopped the internet](./cloudflare-2019.md)  | 2019 | M06, M08       | CPU exhaustion is an outage; test hot-path code  |
| [Slack — thundering herd after maintenance](./slack-2022.md)          | 2022 | M06, M07       | Reconnection storms; jitter and gradual restore  |

## How To Read a Case Study

Read **with a notebook open**. After each one, write down:

1. **What was the immediate cause?** (the thing that broke)
2. **What was the root cause?** (the underlying design or decision)
3. **What quality attribute was implicitly deprioritized?**
4. **In my current system, what's the equivalent risk?**
5. **One specific change I'd make as a result of reading this.**

The fourth question is the highest-leverage one. Translation from "their outage" to "my system" is where lessons stick.

## The Larger Pattern

Read 5+ case studies and a pattern emerges: **most major outages are not bugs. They're architecture.** Specifically:

- **Blast radius too large** (one component fail → cascade)
- **Quality attribute never named** (e.g., "operability" of feature flags)
- **Untested failure mode** (recovery procedure never rehearsed)
- **Defaults wrong for context** (acceptable in one environment, dangerous in another)
- **Human in the loop where they shouldn't be** (a 30-second typo → 4-hour outage)

These are _architectural_ failures, addressable at design time. That's what this course exists to develop the muscle for.

## Further Reading

Public post-mortems worth your time:

- **Cloudflare engineering blog** — exemplary post-mortems, especially the BGP and routing ones
- **AWS status page archive** — gold standard for transparency
- **Stripe engineering blog** — payments-domain outages are uniquely instructive
- **HackerNews "post-mortem" tag** — community-curated
- **The Morning Paper** (now archived) — academic + industry summaries
- **John Allspaw's writings** — the foundational human-factors-of-outages thinker
