---
allowed-tools: Bash(git diff:*), Bash(git status:*), Bash(git branch:*)
description: Generate a conventional commit message for staged changes
---

## Context

- Staged changes: !`git diff --cached`
- Unstaged changes: !`git diff`
- Current branch: !`git branch --show-current`

## Task

Generate a conventional commit message for the staged changes above.

Rules:
- Format: `<type>(<scope>): <short description>`
- Types: `feat`, `fix`, `chore`, `docs`, `test`, `refactor`, `perf`
- Scope: the package or area changed — e.g. `release`, `deploy`, `proto`, `cli`, `db`, `ci`
- Short description: lowercase, no period, max 72 chars, imperative mood ("add" not "added")
- If the change is significant, add a body after a blank line (what and why, not how)
- If there are breaking changes, add `BREAKING CHANGE:` footer
- Use + for bullet points in the body, not *

Output only the commit message, nothing else. No explanation. No markdown fences.

Output the commit message inside a plain code block so it can be copied exactly as-is.

Example output with a list body:
```
chore(ci): add path-filtered workflows for Go and Python

+ Go tests run only on changes under core/
+ Python tests run only on changes under cli/
+ buf lint runs only on changes under proto/
```

Example output with a prose body:
```
feat(release): add state machine for release lifecycle

Implements draft → running → deployed/failed transitions.
State is validated before any promotion attempt.
```