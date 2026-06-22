# Claude Code Commands

YallaOps ships custom slash commands for [Claude Code](https://claude.ai/code) in `.claude/commands/`. Each command is a markdown file with instructions and allowed tools. Invoke them by typing `/command-name` in the Claude Code prompt.

---

## `/new-branch`

**What it does:** Creates a correctly named branch from the latest `main`.

Picks the right branch type (`feat/`, `fix/`, `chore/`), generates a short lowercase name from your description, checks for uncommitted changes, fetches and pulls `main`, and confirms what the first task should be based on the current phase in `CLAUDE.md`.

**When to use:** At the start of any new task.

**Example:**
```
/new-branch add release state machine
```

---

## `/new-feature`

**What it does:** Plans and specs a new feature before writing any code.

Reads `ARCHITECTURE.md` and the current phase, asks clarifying questions about anything ambiguous, then writes a spec covering: what the feature does, proto changes, DB changes, Go packages affected, acceptance criteria, and out-of-scope items. Proposes a branch name and waits for your approval before touching any code.

**When to use:** Before starting any significant new capability.

**Example:**
```
/new-feature approval gate for staging promotions
```

---

## `/commit`

**What it does:** Generates a conventional commit message for your staged changes.

Reads `git diff --cached`, produces a message in `<type>(<scope>): <description>` format with an optional body. Output is ready to paste directly into `git commit -m`.

**When to use:** After staging your changes with `git add`.

**Example:**
```
/commit
```

---

## `/review`

**What it does:** Deep code review of the current branch diff vs `main`.

Checks for: logic errors and race conditions, Go convention violations (from `.claude/rules/go.md`), security issues (hardcoded secrets, missing auth), database problems (missing transactions, N+1 queries), test coverage gaps, and proto breaking changes. Returns findings as **Must fix**, **Should fix**, and **Suggestions**.

**When to use:** Before opening a PR, or when you want a second pass on your own work.

**Example:**
```
/review
```

---

## `/pr-ready`

**What it does:** Full pre-PR checklist.

Runs Go tests and lint, checks commit message format, scans the diff for secrets and hardcoded URLs, validates that proto stubs were regenerated if `.proto` files changed, checks that migration `.down.sql` files exist, and drafts a PR title and description you can copy-paste into GitHub.

**When to use:** Right before pushing and opening a PR.

**Example:**
```
/pr-ready
```

---

## `/debug`

**What it does:** Structured debugging — reproduce, locate, explain, fix, verify.

Walks through: (1) write the minimal command or test that triggers the issue, (2) locate the exact file and function, (3) explain what's happening vs what should happen, (4) implement a minimal fix, (5) verify the fix, (6) scan for other locations with the same pattern.

**When to use:** When you have a specific bug or error to investigate.

**Example:**
```
/debug nil pointer in promotion handler when release has no environments
```

---

## `/sync-phase`

**What it does:** Reviews progress and updates the current phase in `CLAUDE.md`.

Reads recent commits, scans `core/internal/`, `proto/`, and `core/db/migrations/`, runs tests, then tells you: what's been completed, what's remaining, whether you're ready to move to the next phase, and what the next 2–3 tasks should be. Updates the "Current phase" section in `CLAUDE.md`.

**When to use:** At the end of a work session, or when you're not sure what to pick up next.

**Example:**
```
/sync-phase
```
