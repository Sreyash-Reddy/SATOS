# SATOS

**S**warm **A**gentic **T**eam **O**perating **S**ystem

A swarm of specialized AI agents that discuss, debate, research, and THEN code. Built for high-stakes, get-it-right problem solving.

## Quick Start

```bash
# SATOS manages projects via .satos/ directory in each project
# See docs/ for full documentation
```

## The Workflow

```
Discuss → Research → Debate → Code
```

More discussions first. Coding is the last step, not the first.

## Architecture

```
Director
  └── Team Leader (TL)
        └── Teams
              ├── UI Team
              ├── Backend Team
              ├── Researcher Team
              ├── Database Team
              └── Testing Team (independent layer)
```

Teams explore multiple solutions, debate trade-offs, and the TL picks the optimal path with documented rationale.

## Key Principles

- **Discuss before code** — Multiple solutions explored, not just one
- **Testing as independence** — Reports to Director, not under a TL
- **Evidence-backed** — Every PR has video/test evidence
- **High-stakes** — Get it right once, not wrong five times

## Directory Structure

```
.satos/           # SATOS state (per-project)
docs/             # Full documentation
  satos.md        # Overview
  architecture/   # System architecture
  teams/          # Team descriptions
  workflow/       # Workflow process
  pr-standards/   # PR standards
evidence/         # Video/image evidence
```

## Documentation

- [Overview](docs/satos.md) — What is SATOS
- [Architecture](docs/architecture/README.md) — System design
- [Teams](docs/teams/README.md) — Team structure and roles
- [Workflow](docs/workflow/README.md) — Discuss → Research → Debate → Code
- [PR Standards](docs/pr-standards/README.md) — Evidence requirements

## Status

🚧 **V1 Bootstrap** — Initial setup in progress

See [milestones](https://github.com/Sreyash-Reddy/SATOS/milestones) for roadmap.