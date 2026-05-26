# Architecture

## Overview

```
Director (top level - coordinates TLs)
  └── TL (Team Leader)
        └── Team Members (specialists)
              ├── UI Team
              ├── Backend Team
              ├── Researcher Team
              └── Database Team
```

## Core Components

### Director

**Role:** Top-level coordinator. Stateless, per-task invocations.

**Responsibilities:**
- Task assignment to TLs
- Cross-team coordination
- Final arbitration when TLs disagree
- Maintains decision audit trail

**Status:** Stateless per-task (pending implementation validation)

### Team Leader (TL)

**Role:** Judge and coordinator for a team.

**Responsibilities:**
- Direct team members
- Collect and evaluate team proposals
- Run internal team debates
- Pick winner within team
- Report to Director

**TL as Judge:** Not just a coordinator — evaluates multiple solutions, weighs trade-offs, selects optimal path with documented rationale.

### Team Members

**Role:** Specialists in their domain.

**Teams:**
- **UI Team** — frontend, design, UX
- **Backend Team** — server-side, APIs
- **Researcher Team** — business research, market analysis, product thinking
- **Database Team** — data modeling, queries, data layer

**Core Behavior:** Explore all solutions, challenge each other, debate trade-offs — THEN code.

## Communication Flow

### Intra-Team (Same TL)
```
Team Member A ←→ Team Member B (debate directly)
        ↓
       TL (judges winner)
```

### Inter-Team (Different TLs)
```
Researcher Team Member ←→ TL-Researcher
        ↓
     Director (arbitrates cross-team conflicts)
        ↓
     TL-Backend
        ↓
Backend Team Member
```

> **Note:** Cross-team debates route through Director. Direct TL-to-TL bypasses Director oversight.

## Testing Layer

**Placement:** Independent verifier layer, reports to Director.

**Why:** If testing lives under a TL, that TL can suppress inconvenient test results. Testing needs independence to function as objective judge.

```
Director
  ├── TL-UI
  ├── TL-Backend
  ├── TL-Researcher
  └── TL-Database
  └── Testing Layer (independent)
```

## Open Architecture Decisions

| Decision | Status | Notes |
|----------|--------|-------|
| Director persistence | OPEN | Stateless (current) vs Persistent (needs validation) |
| Cross-team debate path | OPEN | Current: 3 hops via Director. Could be optimized. |
| Debate termination trigger | OPEN | No "debate ends here" mechanism defined |
| Shared reasoning substrate | OPEN | Currently outputs-only. Reasoning chains not shared. |
| TL conflict resolution | OPEN | Who judges when judges (TLs) disagree? |

## Key Gaps (From Party Mode Review)

1. **Decision audit trail** — Where does "why TL picked winner" live 6 months later?
2. **Debate termination** — What prevents a 10-minute task going 3 hours in debate?
3. **Testing independence** — Must report to Director, not under a TL.
4. **Shared reasoning** — Agents exchange outputs (conclusions) not reasoning chains.

## Related

- [Teams](docs/teams/README.md)
- [Workflow](docs/workflow/README.md)