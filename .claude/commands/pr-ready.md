---
allowed-tools: Bash(git diff:*), Bash(git log:*), Bash(git status:*), Bash(git branch:*)
description: Check if the current branch is ready to open a PR
---

## Context

- Current branch: !`git branch --show-current`
- Changed files: !`git diff main --name-only`
- Full diff: !`git diff main`
- Commits on this branch: !`git log main..HEAD --oneline`
- Test output: !`cd core && go test ./... 2>&1 | tail -30`
- Lint output: !`cd core && golangci-lint run ./... 2>&1 | tail -20`

## Task

Review the branch and tell me if it is ready to open a PR. Check:

1. **Tests** — do all tests pass? If not, list the failures.
2. **Lint** — any lint errors? List them.
3. **Scope** — does this branch do one logical thing, or is it mixing concerns? If mixed, suggest how to split.
4. **Commit messages** — are they conventional commit format (`feat:`, `fix:`, `chore:` etc.)? List any that aren't.
5. **Secrets check** — does the diff contain any API keys, passwords, or hardcoded URLs? Flag immediately if so.
6. **Proto** — if any `.proto` files changed, were the generated stubs (`buf generate`) also committed?
7. **Migration check** — if any migration files changed, is there a corresponding `.down.sql`?
8. **PR description** — draft a PR title and description I can copy-paste, following this format:

```
## What
<one sentence — what does this PR do>

## Why
<one sentence — why is this needed>

## Changes
- <bullet list of key changes>

## Testing
<how to verify this works>
```

Be direct. If it's not ready, tell me exactly what to fix.
