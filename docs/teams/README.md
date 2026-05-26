# Teams

## Overview

Each team has complete knowledge in their domain. They discuss, explore options, debate trade-offs — then code.

**Motto:** Discuss → Research → Debate → Code

## Team Structure

```
TL (Team Leader)
  └── Team Members (2-5 specialists)
```

## Teams

### UI Team

**Domain:** Frontend, design, UX

**Responsibilities:**
- User interface design and implementation
- Design system maintenance
- Accessibility compliance
- Cross-browser/cross-device testing

**Key Discussions:**
- Component architecture choices
- Design patterns vs flexibility
- User experience trade-offs

### Backend Team

**Domain:** Server-side, APIs, business logic

**Responsibilities:**
- API design and implementation
- Business logic layer
- Performance optimization
- Security implementation

**Key Discussions:**
- API design (REST vs GraphQL vs other)
- Data flow architecture
- Caching strategies
- Error handling approaches

### Researcher Team

**Domain:** Business research, market analysis, product thinking

**Responsibilities:**
- Problem space exploration
- Market analysis
- Product requirements definition
- PRD creation
- Business model validation

**Key Discussions:**
- What problem are we solving?
- Is this worth building?
- Who is the user?
- What does success look like?

### Database Team

**Domain:** Data modeling, queries, data layer

**Responsibilities:**
- Data model design
- Query optimization
- Data migration strategies
- Data integrity

**Key Discussions:**
- Schema design trade-offs
- Normalization vs performance
- Query strategies
- Migration approaches

## Testing Team

**Placement:** Independent layer, reports to Director (not under any TL)

**Why:** Testing needs independence to avoid suppression of inconvenient results.

**Responsibilities:**
- Functional testing
- Load testing (calls out to k6)
- Performance benchmarking
- Video evidence for UI bug fixes
- Test Report generation

## Key Principle: Debate Before Code

Teams don't just output solutions. Before coding begins:

1. **Explore** — Multiple team members propose different approaches
2. **Challenge** — Each solution is questioned and stress-tested
3. **Debate** — Trade-offs are discussed openly
4. **Converge** — TL picks optimal solution with documented rationale

**TL as Judge:** Evaluates multiple solutions, not just collects outputs. The goal is adversarial friction that produces genuinely better solutions, not polished presentations.

## Team Communication

### Within Team
- Team members debate directly with each other
- TL judges and picks winner
- All debates documented in team log

### Across Teams
- Cross-team debates route through Director
- Director arbitrates conflicts
- Solutions shared via shared context files

## Related

- [Architecture](docs/architecture/README.md)
- [Workflow](docs/workflow/README.md)