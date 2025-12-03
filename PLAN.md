# Decision: Fork Wakapi for OxyWaka

**Date:** 2025-12-02
**Decision:** Fork Wakapi and extend it with AI insights, rather than build from scratch

---

## Context

OxyWaka aims to provide AI-powered productivity insights by tracking:
- Coding activity (via WakaTime plugins)
- Browser activity (via WakaTime browser extensions)
- Desktop activity (via ActivityWatch - Phase 2)

**Core question:** Build heartbeat ingestion from scratch or fork existing solution?

---

## The Decision

**Fork Wakapi and extend it with insights module**

### Rationale

1. **Time to Market**
   - Fork: 4 weeks to working insights
   - Scratch: 10 weeks to working insights
   - **5.5 weeks saved (40% faster)**

2. **Learning Approach**
   - Reading 20k lines of production Go code teaches more than writing 5k lines of beginner code
   - Learn battle-tested patterns, not trial-and-error basics
   - Understand real-world trade-offs and decisions

3. **Focus on Unique Value**
   - OxyWaka's value = AI insights + activity correlation
   - NOT heartbeat ingestion (solved problem)
   - Forking lets us focus on insights immediately

4. **Production-Ready Foundation**
   - Authentication, API keys, rate limiting: ✅ Done
   - Database schema, migrations: ✅ Done
   - Heartbeat validation, deduplication: ✅ Done
   - Aggregation queries, time handling: ✅ Done
   - Bug fixes from real-world usage: ✅ Done

5. **Low Risk**
   - Evaluation cost: 2-3 days to understand codebase
   - If unsuitable: Build from scratch with learned patterns
   - Asymmetric risk/reward in favor of forking

---

## What We Get from Wakapi

### Keep (Production-Ready Components)
- ✅ Heartbeat ingestion & validation
- ✅ Authentication & API keys
- ✅ Database schema & migrations
- ✅ REST API endpoints (WakaTime-compatible)
- ✅ Aggregation logic (daily/weekly summaries)
- ✅ Multi-user support (even if single-user initially)

### Strip (If Needed)
- Public leaderboards (if present)
- Complex org/team features (if too complex)
- Email notifications (rebuild differently)

### Add (Our Innovation)
- ✅ Insights module (`internal/insights/`)
  - Activity categorization (productive/learning/distracted)
  - Pattern detection (deep work, context switching)
  - Productivity scoring
  - AI-powered recommendations
- ✅ Enhanced dashboard with insights visualization
- ✅ ActivityWatch integration (Phase 2)

---

## Phased Roadmap

### Phase 1: Wakapi + Insights (4 weeks)
**Goal:** AI insights on existing coding + browser data

- Week 1: Fork, evaluate, understand codebase
- Week 2: Design & implement insights module
- Week 3: Build categorization & pattern detection
- Week 4: Dashboard UI for insights

**Deliverable:** OxyWaka V1 with AI insights on WakaTime data

### Phase 2: ActivityWatch Integration (2 weeks)
**Goal:** Add desktop app tracking for complete picture

- Week 5: ActivityWatch API client
- Week 6: Data correlation & unified timeline

**Deliverable:** OxyWaka V2 with full desktop tracking

### Phase 3: Advanced Features (ongoing)
- Mobile tracking
- Advanced AI insights
- Team features (if validated)

---

## Data Sources (Phase 1)

**WakaTime Coverage:**
- ✅ Code editors (VS Code, JetBrains, etc.)
- ✅ Browser activity (Chrome, Firefox, Edge extensions)
- ✅ Project/language/file tracking
- ✅ Git branch tracking

**What This Enables:**
- Context switching quantification
- Productive vs learning vs distracted time
- Deep work session detection
- Project time breakdown
- Research patterns (docs → coding correlation)

**Estimated Coverage:** 80% of developer activity (sufficient for V1)

---

## Alternatives Considered

### Alternative 1: Build from Scratch
- **Time:** 10 weeks to insights
- **Learning:** Maximum control, maximum effort
- **Verdict:** Reinvents solved problems, delays unique value

### Alternative 2: Use Wakapi As-Is (Separate Service)
- **Time:** 4 weeks to insights
- **Architecture:** Two services to maintain
- **Verdict:** Less integrated, API overhead

### Alternative 3: Skip Wakapi, Use ActivityWatch Only
- **Coverage:** Desktop apps but no coding detail
- **Integration:** WakaTime plugins wouldn't work
- **Verdict:** Loses coding-specific insights

---

## Success Criteria for Fork Evaluation

**Spend 2-3 days evaluating Wakapi codebase:**

### Must Answer:
1. ✅ Is code readable and well-structured?
2. ✅ Does architecture align with OxyWaka vision?
3. ✅ Can insights module integrate cleanly?
4. ✅ Is database schema suitable for our needs?
5. ✅ Are dependencies manageable (Go version, libraries)?

### Decision Point:
- **If YES to 4+**: Proceed with fork
- **If NO to majority**: Build from scratch with learned patterns

**Worst case:** 2 days spent learning from production code before building own solution
**Best case:** 5.5 weeks of development time saved

---

## Technical Architecture (Post-Fork)

```
oxywaka/ (forked from Wakapi)
├── cmd/
│   └── server/              # Wakapi's server (keep)
├── internal/
│   ├── api/                 # Wakapi's API handlers (keep/extend)
│   ├── models/              # Wakapi's data models (keep)
│   ├── storage/             # Wakapi's database layer (keep)
│   ├── insights/            # NEW: Our insights module
│   │   ├── categorizer.go   # Activity categorization
│   │   ├── patterns.go      # Deep work, context switching
│   │   ├── scorer.go        # Productivity scoring
│   │   └── recommendations.go
│   └── clients/             # NEW: External integrations
│       └── activitywatch.go # Phase 2
├── web/                     # Dashboard (enhance Wakapi's UI)
└── migrations/              # Database (add insights tables)
```

---

## Key Principles

1. **Ship Insights Fast**
   - Infrastructure is a means, not the goal
   - Focus on unique value from day 1

2. **Learn by Reading**
   - Production code teaches real-world patterns
   - Understand trade-offs experts made

3. **Incremental Complexity**
   - V1: WakaTime only (code + browser)
   - V2: Add ActivityWatch (desktop apps)
   - V3: Advanced features

4. **Validate Early**
   - Are AI insights actually useful?
   - Is unified timeline valuable?
   - Answer these in 4 weeks, not 10

---

## Next Steps

### Immediate (This Week)
1. Fork Wakapi repository on GitHub
2. Clone locally: `git clone https://github.com/YOUR_USERNAME/wakapi.git oxywaka`
3. Run and explore codebase
4. Install WakaTime browser extension
5. Collect sample data (code + browse for 2 days)

### Week 1
- Read key files: main.go, database schema, heartbeat handler
- Understand project structure and patterns
- Make fork evaluation decision
- Design insights module architecture

### Week 2+
- Implement insights module
- Build AI categorization logic
- Create dashboard for insights

---

## References

- **Wakapi:** https://github.com/muety/wakapi
- **WakaTime API:** https://wakatime.com/developers
- **Project Docs:** See CLAUDE.md, ROADMAP.md in this repo

---

**Decision made:** 2025-12-02
**Status:** Approved - Proceeding with fork
**Next review:** After Week 1 evaluation phase
