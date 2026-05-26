# SATOS Overview

**S**warm **A**gentic **T**eam **O**perating **S**ystem

## What SATOS Is

A swarm of specialized AI agents that discuss, debate, research, and THEN code. The goal is discussions first, coding second.

**Core thesis:** Agents need proper guidance + dedicated roles beat generalist agents. Copilot CLI fails without guidance. SATOS builds directed agent teams with proper prompts and task definitions.

## The Workflow

```
Discuss → Research → Debate → Code
```

Teams don't just output — they explore all solutions, challenge each other, find the optimal path. Coding is the last step, not the first.

## Work Type

**High-stakes, get-it-right problems.** Quality > speed. Getting it right once through thorough debate beats fast-then-wrong-then-retry cycles.

## Key Documents

| Document | Description | Produced By |
|----------|-------------|-------------|
| PRD | Product Requirements Document | Researcher Team |
| TRD | Technical Requirements Document | Backend/Database Teams |

> These are documents, NOT teams. Teams produce these as part of their work.

## Directory Structure

```
SATOS/
  .satos/              # SATOS state (per-project)
    tasks/
      queue/           # Pending tasks
      done/            # Completed tasks
      failed/          # Failed tasks
    state.db           # Task state database
    director.log       # Audit log
    verifiers/         # Test verification scripts

  docs/                # All documentation
    satos.md           # This file
    architecture/      # Architecture docs
    teams/             # Team descriptions
    workflow/          # Workflow docs
    pr-standards/      # PR standards

  src/                 # Project source
  tests/              # Project tests
  evidence/           # Video/image evidence
```

## Quick Links

- [Architecture](architecture/README.md)
- [Teams](teams/README.md)
- [Workflow](workflow/README.md)
- [PR Standards](pr-standards/README.md)