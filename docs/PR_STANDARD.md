# PR Standards

Every PR in this project follows a strict standard. No "normal" PRs.

## Why Standards?

SATOS is a high-stakes, get-it-right system. PRs are the delivery mechanism. Every PR must be:
- Self-contained
- Evidence-backed
- Properly labeled
- Fully documented

## PR Title Format

```
{type}: {short description}

Types: feat, fix, docs, refactor, test, chore
```

**Examples:**
```
feat: add user authentication
fix: resolve memory leak in cache
docs: add API documentation
test: add load testing for auth endpoint
```

## PR Description

Every PR must include:

### 1. Summary
What does this PR do? One paragraph max.

### 2. Type
One of: `bug-fix`, `feature`, `docs`, `refactor`, `test`, `chore`

### 3. Evidence
Required evidence varies by type:

| Type | Evidence |
|------|----------|
| Bug Fix | Video/image of bug + video/image of fix |
| Feature | Video/image of feature working |
| Docs | Justification + scope of changes |
| Backend | Test results + evidence |
| Refactor | Before/after comparison |

### 4. Testing
- How was this tested?
- What tests were run?
- Results?

### 5. Tags/Labels
Required labels:
- `documentation` — for docs changes
- `bug-fix` — for bug fixes
- `feature` — for new features
- `backend` — for backend changes
- `frontend` — for frontend changes
- `test` — for test changes
- `refactor` — for refactors

## Evidence Requirements

### Video Evidence
For UI changes and bug fixes, video evidence is required:
- Named: `{issue-name}-{before|after}.webm`
- Stored in: `evidence/` directory
- Shows: bug existing → bug fixed

### Image Evidence
For simple fixes or features:
- Named: `{feature-name}.png`
- Stored in: `evidence/` directory
- Shows: working state

### Test Results
For backend/code changes:
- Include test output in description
- Attach test report if available
- Include performance benchmarks if applicable

## Checklist

Before merging, verify:
- [ ] Title follows format
- [ ] Description complete
- [ ] Evidence attached (or justification for why not)
- [ ] Tests pass
- [ ] Labels applied
- [ ] No merge conflicts
- [ ] Docs updated (if applicable)

## Example PR

### Title
```
feat: add project creation workflow
```

### Description
```
## Summary
Adds a complete project creation workflow including form validation,
API integration, and success/error states.

## Type
feature

## Evidence
[Video showing workflow end-to-end](evidence/project-creation.webm)

## Testing
- Unit tests: 12 passing
- Integration tests: 3 passing
- Manual testing: Complete

## Tags
feature, backend, frontend
```

## Review Process

1. PR created → automated checks run
2. Reviewer checks: evidence, tests, standards
3. Changes requested (if needed)
4. Approved → merged

No PR merges without proper evidence.
