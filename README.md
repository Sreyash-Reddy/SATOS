# SATOS

**S**elf-**A**gent **T**eam **O**perating **S**ystem

An agentic team operating system for personal use. Manages multiple AI agent teams (research, development, testing) that collaborate autonomously to solve high-stakes problems correctly the first time.

## Overview

SATOS is a personal AI orchestration layer that:
- Runs multiple specialized AI agents in parallel
- Teams debate and explore multiple solutions
- Testing validates with evidence (video, tests, benchmarks)
- A Team Leader reviews all options and picks the optimal solution
- Every output is structured, documented, and evidenced

## Architecture

```
User (you)
  └── Director (stateless, per-task)
        └── Team Leader (coordinates teams)
              ├── PRD Team (what to build, why)
              └── TRD Team (how, testing, validation)
```

## Getting Started

> SATOS setup instructions coming soon.

## Workflow

1. **Define** — PRD Team produces the Product Requirements Document
2. **Design** — TRD Team produces the Technical Requirements Document
3. **Build** — Implementation agents work from TRD
4. **Test** — Testing team validates with evidence
5. **Review** — Team Leader picks optimal solution
6. **Ship** — Merged via structured PR

## Standards

Every PR follows a strict standard:
- Structured title and description
- Evidence (video/image) for all changes
- Tags: `docs`, `bug-fix`, `feature`, `backend`, `test`
- Test results included

See [docs/PR_STANDARD.md](docs/PR_STANDARD.md) for full PR requirements.

## Directory Structure

```
.satos/           # SATOS state (per-project)
  tasks/
    queue/        # Pending tasks
    done/         # Completed tasks
    failed/       # Failed tasks
  verifiers/      # Test verification scripts
  state.db        # Task state database
  director.log    # Audit log

.github/
  PULL_REQUEST_TEMPLATE/  # PR template

docs/             # Project documentation
```

## Current Status

🚧 **Under Development** — V1 bootstrap in progress
