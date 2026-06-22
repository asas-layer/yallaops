---
allowed-tools: Bash(git diff:*), Bash(git log:*)
description: Deep code review of the current branch changes
---

## Context

- Diff vs main: !`git diff main`
- Files changed: !`git diff main --name-only`

## Task

Do a thorough code review of the changes above. Focus on:

### Correctness
- Logic errors, off-by-one errors, incorrect state transitions
- Race conditions or missing mutex locks in Go concurrent code
- Unhandled error cases

### YallaOps conventions (from @.claude/rules/go.md)
- Is context passed correctly to all DB and network calls?
- Are errors wrapped with context?
- Is slog used instead of fmt.Println?
- Are there any raw SQL strings that should go through sqlc?
- Any global state introduced?

### Security
- Hardcoded secrets or credentials
- Missing input validation on gRPC handlers
- Any auth bypass risk

### Database
- Missing transactions on multi-table writes
- N+1 query patterns
- Missing indexes on columns that will be queried frequently

### Tests
- Are new functions covered by tests?
- Are tests table-driven?
- Any test that only tests the happy path with no error cases?

### Proto
- If proto files changed, were stubs regenerated?
- Any field deletions or renumbering (breaking changes)?

Format your review as:
- **Must fix** — blockers, cannot merge
- **Should fix** — not blockers but important
- **Suggestions** — optional improvements

Be direct and specific. Reference line numbers or function names where possible.
