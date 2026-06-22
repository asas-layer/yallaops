---
allowed-tools: Bash(git status:*), Bash(git branch:*), Bash(git checkout:*), Bash(git fetch:*), Bash(git log:*)
description: Create a correctly named branch from latest main for a new task
---

## Context

- Current branch: !`git branch --show-current`
- Current status: !`git status --short`
- Latest branches: !`git branch -a | head -20`

## Task

I want to work on: $ARGUMENTS

Do the following:

1. **Decide the branch type** based on what I described:
   - `feat/` — new feature or capability
   - `fix/` — bug fix
   - `chore/` — tooling, CI, config, docs, refactoring with no behavior change

2. **Generate a branch name** following these rules:
   - Format: `<type>/<short-name>`
   - Short name: lowercase, hyphens only, no underscores, max 4 words
   - Must describe what the branch does, not what phase it's in
   - Good: `feat/release-state-machine`, `fix/nil-context-promotion`, `chore/sqlc-setup`
   - Bad: `feat/phase1`, `fix/bug`, `chore/stuff`

3. **Check for uncommitted changes** — if `git status` shows dirty state, warn me and ask if I want to stash before switching.

4. **Create the branch from latest main**:
   ```bash
   git fetch origin
   git checkout main
   git pull origin main
   git checkout -b <branch-name>
   ```

5. **Confirm** by showing me:
   - The branch name chosen and why
   - The current branch after switching
   - What the first task should be based on @CLAUDE.md and @ARCHITECTURE.md current phase

Do not start writing any code. Just create the branch and orient me on what to do first.
