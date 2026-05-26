# Workflow

## Core Principle

**Discuss → Research → Debate → Code**

More discussions, more debates, THEN coding. Coding is the last step, not the first.

## High-Level Flow

```
Problem → Teams Explore Solutions → Teams Debate → Testing Validates → Winner Picked → Code
```

## Detailed Workflow

### Phase 1: Problem Definition

1. **Director** receives the problem/task
2. **Director** assigns to appropriate TLs
3. **TLs** brief their team members

### Phase 2: Solution Exploration

1. **Team Members** independently explore different solutions
2. Multiple proposals created (not one solution)
3. Each proposal documents reasoning, trade-offs, alternatives considered

### Phase 3: Debate

1. **Team Members** present solutions
2. Cross-examination happens:
   - Challenge assumptions
   - Stress-test conclusions
   - Explore failure modes
3. **TL** facilitates debate (does not judge yet)
4. Debate continues until:
   - Convergence on optimal path, OR
   - Clear trade-offs identified, OR
   - Termination trigger fires (TBD)

### Phase 4: Judgment

1. **TL** reviews debate outcomes
2. **TL** evaluates alternatives with documented rationale
3. **TL** picks winner with explanation
4. Decision logged to audit trail

### Phase 5: Testing

1. **Testing Team** receives solution
2. **Testing Team** validates:
   - Functional correctness
   - Performance benchmarks
   - Load testing (k6)
   - Video evidence for UI changes
3. **Test Report** generated with all evidence
4. Failures route back to Phase 2

### Phase 6: Implementation

1. **Winning solution** implemented by team
2. **PR** created with proper standards
3. **Evidence** attached (video, tests, benchmarks)

## Retry Cycle

If Testing fails:
```
Testing → Reports failure → Team addresses → Testing validates again
```

Retry cycle is intentional. Getting it right after 2 iterations beats shipping wrong on iteration 1.

## Document Flow

| Phase | Documents Produced |
|-------|-------------------|
| Problem Definition | PRD (by Researcher Team) |
| Technical Design | TRD (by Backend/Database Teams) |
| Testing | Test Report |
| Implementation | PR with evidence |

## State Management

Each project has `.satos/` directory:
```
.satos/
  tasks/
    queue/     # Pending tasks
    done/      # Completed tasks
    failed/    # Failed tasks
  state.db     # SQLite: task state, team assignments
  director.log # Decision audit trail
```

## Open Workflow Decisions

| Decision | Status | Notes |
|----------|--------|-------|
| Debate termination trigger | OPEN | How to know when debate has produced enough signal? |
| Shared reasoning substrate | OPEN | Currently outputs-only. Reasoning chains not shared. |
| TTL calibration | OPEN | Adaptive vs fixed timeout per task |

## Related

- [Architecture](docs/architecture/README.md)
- [Teams](docs/teams/README.md)
- [PR Standards](docs/pr-standards/README.md)