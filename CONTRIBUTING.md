# Contributing to YallaOps

First off — thanks for taking the time to contribute. YallaOps is early-stage and every contribution matters.

---

## Before you start

- Check [open issues](https://github.com/asas-layer/yallaops/issues) to see if your idea or bug is already tracked
- For anything significant (new feature, architectural change), open an issue first and discuss it before writing code — saves everyone time
- Small fixes (typos, docs, obvious bugs) can go straight to a PR

---

## Setup

### Prerequisites

- Go 1.23+
- Python 3.12+
- Docker + Docker Compose
- `buf` — protobuf toolchain (`brew install bufbuild/buf/buf` or see [buf.build](https://buf.build/docs/installation))
- `sqlc` — SQL code generation (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)
- `golangci-lint` — Go linter (`brew install golangci-lint` or see [golangci-lint.run](https://golangci-lint.run/usage/install/))
- `just` — task runner (`brew install just` or see [just.systems](https://just.systems))

### Local environment

```bash
# 1. Fork and clone
git clone https://github.com/<your-username>/yallaops
cd yallaops

# 2. Start Postgres + Redis
docker compose up -d

# 3. Run migrations
just migrate

# 4. Start the gRPC server
just dev-core

# 5. Install the CLI (in another terminal)
cd cli && pip install -e .
yallaops status
```

If `yallaops status` returns a response, you're good to go.

---

## Workflow

### 1. Create a branch

Never work directly on `main`. Always branch off from it:

```bash
git checkout main
git pull origin main
git checkout -b feat/your-feature-name
```

Branch naming:

| Prefix | When to use |
|---|---|
| `feat/` | New feature |
| `fix/` | Bug fix |
| `chore/` | Tooling, CI, config, docs |

Keep branch names short and lowercase: `feat/release-state-machine`, not `feat/AddTheReleaseStateMachineForTheGoCore`.

### 2. Make your changes

A few hard rules:

- **Proto first** — if your change touches the API, define it in `proto/` before writing Go or Python code. Run `buf generate` to regenerate stubs.
- **No ORMs** — all database access goes through sqlc-generated functions. No raw SQL strings in Go code.
- **No REST** — the API is gRPC only (except `GET /healthz`).
- **Tests required** — new Go code needs table-driven tests. PRs without tests will be asked to add them.
- **No secrets** — this repo is public. Never commit API keys, passwords, or internal URLs.

See `ARCHITECTURE.md` for the full list of technical decisions and constraints.

### 3. Run checks locally

Before pushing, make sure everything passes:

```bash
just lint      # golangci-lint + ruff + buf lint
just test      # go test ./... + python tests
just proto     # buf generate (if you changed any .proto files)
just sqlc      # sqlc generate (if you changed any query files)
```

All of these must pass before opening a PR. CI will run the same checks and a failing CI blocks merge.

### 4. Commit messages

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <short description>

[optional body]
```

Types: `feat`, `fix`, `chore`, `docs`, `test`, `refactor`, `perf`

Scope: the package or area — `release`, `deploy`, `proto`, `cli`, `db`, `ci`

Examples:

```
feat(release): add state machine for release lifecycle
fix(api): handle nil context in promotion handler
chore(ci): add path filters for Go and Python workflows
docs(proto): document release status enum values
test(release): add table-driven tests for promotion logic
```

Rules:
- Lowercase, no period at the end
- Imperative mood — "add" not "added", "fix" not "fixed"
- Max 72 characters on the first line
- If it's a breaking change, add `BREAKING CHANGE:` in the footer

### 5. Open a PR

Push your branch and open a PR against `main`:

```bash
git push origin feat/your-feature-name
```

PR description should cover:

- **What** — one sentence describing what this PR does
- **Why** — one sentence on why it's needed
- **Changes** — bullet list of key changes
- **Testing** — how to verify it works

Keep PRs small and focused. One logical change per PR. If you find yourself writing "and also..." in the description, split it into two PRs.

### 6. Review process

- All PRs require at least one approval before merge
- CI must be green (lint + tests)
- We squash merge — your branch commits collapse into one commit on `main`
- Main is always deployable — don't merge broken code

---

## Claude Code commands

If you use [Claude Code](https://claude.ai/code), the project ships custom slash commands in `.claude/commands/`. They automate the most common workflows. Invoke them by typing `/command-name` in the Claude Code prompt.

| Command | What it does | Example |
|---------|-------------|---------|
| `/new-branch` | Creates a correctly named branch from latest `main` — picks the right prefix, pulls latest, and orients you on what to do first | `/new-branch add release state machine` |
| `/new-feature` | Plans and specs a new feature before writing any code — asks clarifying questions, writes a proto/DB/Go breakdown, and proposes a branch name | `/new-feature approval gate for staging promotions` |
| `/commit` | Reads staged changes and generates a conventional commit message ready to copy-paste | `/commit` |
| `/review` | Deep code review of the current branch diff — checks correctness, Go conventions, security, database patterns, and test coverage | `/review` |
| `/pr-ready` | Runs lint and tests, checks for secrets, validates commit format, and drafts a PR description | `/pr-ready` |
| `/debug` | Walks through reproduce → locate → explain → fix → verify for a specific issue | `/debug nil pointer in promotion handler` |
| `/sync-phase` | Reads recent commits and the codebase state, tells you what's done and what's left in the current phase, and updates `CLAUDE.md` | `/sync-phase` |

See [`docs/commands.md`](./docs/commands.md) for full details on each command.

---

## What to work on

Good first issues are tagged [`good first issue`](https://github.com/asas-layer/yallaops/issues?q=is%3Aissue+label%3A%22good+first+issue%22).

The project follows a phased roadmap — see `ARCHITECTURE.md` section 4 and `docs/ROADMAP.md`. Work within the current phase unless you've discussed otherwise in an issue. Opening a PR for Phase 5 work when we're in Phase 1 will be deferred.

---

## Project structure

```
core/          Go — gRPC control plane, release engine, Postgres, Redis
agent/         Rust — runtime agent (Phase 4, not started)
cli/           Python — CLI tool
ai-generator/  Python — LLM config generator (Phase 5, not started)
proto/         Protobuf definitions — source of truth for all APIs
infra/         K8s manifests and Helm charts
docs/          Architecture decisions and specs
```

---

## Questions

Open a [GitHub Discussion](https://github.com/asas-layer/yallaops/discussions) for anything that isn't a bug or feature request. We're happy to help you get oriented.

---

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](./LICENSE).
