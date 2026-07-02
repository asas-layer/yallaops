# YallaOps

**A multi-runtime release orchestration platform.** YallaOps manages software delivery through structured releases, environment promotion workflows, policy validation, notifications, and AI-generated deployment configurations.

> يلا — Arabic for "let's go". Because shipping should be fast, not painful.

---

## What it does

Traditional tools either sync Git to Kubernetes (ArgoCD, Flux) or run pipelines (GitHub Actions, Tekton). YallaOps sits in between: it manages the **release lifecycle** — from build to approval to environment promotion to deployment — across Kubernetes, Docker Compose, and Docker Swarm.

```
Source (Git / Image / S3 / AI input)
        ↓
Release created
        ↓
Auto-deployed to dev
        ↓
Approval → staging
        ↓
Policy validation → production
        ↓
Monitoring + notifications + audit
```

You do not deploy directly. You promote **releases**.

---

## Key features

- **Release lifecycle management** — every deployment is a versioned, trackable release object
- **Environment promotion** — dev → staging → prod with configurable approval gates
- **Multi-runtime** — Kubernetes, Docker Compose, and Docker Swarm from one control plane
- **Policy engine** — block promotions that violate rules (replica count, health checks, image digest match)
- **Notification center** — Slack, email, webhooks, and web dashboard for every release event
- **AI deployment generator** — describe your service in plain text, get production-ready K8s YAML or Compose config
- **Rust agent** — optional lightweight agent for metrics, logs, and secure execution inside clusters
- **Full audit log** — every action, every approval, every promotion is recorded

---

## Architecture

```
CLI (Go)      ──gRPC──▶  Go Control Plane
                               │
              ┌────────────────┼────────────────┐
              │                │                │
       Release Engine   Notification      Policy Engine
              │           Center (Redis)        │
              ▼                                 │
     Deployment Orchestrator ◀──────────────────┘
              │
    ┌─────────┼─────────┐
    │         │         │
Kubernetes  Compose   Swarm
              │
         Rust Agent (optional)
```

**Stack:**
| Layer | Technology |
|---|---|
| Core | Go |
| CLI | Go |
| AI generator | Python |
| Runtime agent | Rust (Phase 4) |
| API | gRPC + protobuf |
| Database | PostgreSQL (sqlc + pgx) |
| Events | Redis pub/sub |
| K8s integration | client-go + Helm library |
| Local dev | Docker Compose |

---

## Project status

YallaOps is in active early development. Current phase: **Phase 1 — core release engine**.

| Phase | Status | Description |
|---|---|---|
| 0 — Repo + tooling | ✅ Done | Monorepo structure, CI, proto setup |
| 1 — Release engine | 🔄 In progress | Go core, Postgres, gRPC API |
| 2 — CLI | ⏳ Planned | Go CLI, create/promote/status |
| 3 — K8s deployment | ⏳ Planned | client-go, Helm, rollout tracking |
| 4 — Approvals + notifications | ⏳ Planned | Approval workflow, Slack, Redis events |
| 5 — Docker Compose + Swarm | ⏳ Planned | Multi-runtime support |
| 6 — Rust agent | ⏳ Planned | Metrics, logs, secure exec |
| 7 — AI generator | ⏳ Planned | LLM-powered deployment config gen |

---

## Getting started

### Prerequisites

- Go 1.23+
- Python 3.12+ (for the AI generator, Phase 5)
- Docker + Docker Compose
- `protoc` + Go protobuf plugins
- `buf` (protobuf toolchain)
- `sqlc` (SQL → Go code generation)

### Run locally

```bash
# Clone the repo
git clone https://github.com/your-org/yallaops
cd yallaops

# Start Postgres + Redis
docker compose up -d

# Run database migrations
just migrate

# Start the Go control plane
just dev-core

# In another terminal, build and use the CLI
cd cli
go build -o yallaops ./cmd/yallaops
./yallaops status
```

### Run with Docker

```bash
docker compose --profile full up
```

---

## Repository structure

```
yallaops/
├── core/                  # Go — control plane
│   ├── cmd/server/        # gRPC server entrypoint
│   ├── internal/
│   │   ├── release/       # release engine + state machine
│   │   ├── deploy/        # K8s, Compose, Swarm orchestrators
│   │   ├── policy/        # promotion policy engine
│   │   ├── notify/        # event bus + notification center
│   │   └── api/           # gRPC handlers
│   ├── db/
│   │   ├── migrations/    # SQL migrations
│   │   └── queries/       # sqlc query files
│   └── go.mod
│
├── agent/                 # Rust — runtime agent
│   └── src/
│
├── cli/                   # Go — CLI tool
│   ├── cmd/yallaops/      # CLI entrypoint (Cobra root command)
│   ├── internal/
│   │   ├── commands/      # create, promote, status, dashboard
│   │   ├── client/        # gRPC client wrapper
│   │   └── config/        # ~/.yallaops/config.yaml handling
│   └── go.mod
│
├── ai-generator/          # Python — LLM deployment config generator
│
├── proto/                 # Protobuf definitions (source of truth)
│   ├── release.proto
│   ├── environment.proto
│   └── agent.proto
│
├── infra/                 # K8s manifests, Helm charts
├── docs/                  # Architecture decisions (ADRs)
├── docker-compose.yml     # Local dev: Postgres + Redis
├── Justfile               # Task runner
└── buf.yaml               # Protobuf toolchain config
```

---

## Contributing

YallaOps is open source and welcomes contributions. Please read [CONTRIBUTING.md](./CONTRIBUTING.md) before opening a PR.

**Branching:**
- `main` — always deployable, protected
- `feat/<name>` — new features, max 1-2 days old
- `fix/<name>` — bug fixes
- `chore/<name>` — tooling, CI, docs

**PR rules:**
- Squash merge only
- Must pass CI (Go tests, proto lint)
- One logical change per PR

---

## License

MIT — see [LICENSE](./LICENSE).

---

## Roadmap

See [docs/ROADMAP.md](./docs/ROADMAP.md) for the detailed phase plan.

Community-requested features are tracked in [GitHub Issues](https://github.com/your-org/yallaops/issues). If you want something built, open an issue.
