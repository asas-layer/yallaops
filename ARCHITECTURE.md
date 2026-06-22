# YallaOps — Architecture

This document records every major technical decision made for YallaOps, the reasoning behind each choice, and the constraints that future contributors should be aware of.

---

## 1. What YallaOps is

YallaOps is a **release lifecycle and deployment orchestration platform**. It manages software delivery from build → release → approval → environment promotion → deployment across Kubernetes, Docker Compose, and Docker Swarm.

It is not a GitOps sync tool (that's ArgoCD/Flux). It is not a CI pipeline runner (that's Tekton/GitHub Actions). It sits between those layers and manages the **release object** — a versioned, auditable unit of work that moves through environments with human approval gates.

---

## 2. Core concept: the release object

Everything in YallaOps revolves around a release. A release is created once and promoted through environments — it is never re-created per environment.

```
Release {
  id:          "release-1.4.2"
  service:     "payment-api"
  version:     "1.4.2"
  source:      { type: "image", location: "ghcr.io/org/payment-api:1.4.2" }
  status:      draft | running | failed | deployed
  environments: {
    dev:     deployed
    staging: pending
    prod:    blocked
  }
  created_at:  timestamp
}
```

You do not deploy directly. You promote releases.

---

## 3. Repository structure

YallaOps is a **monorepo**. Three runtimes (Go, Python, Rust) share the same repository because they are tightly coupled at the API boundary. The `proto/` directory is the contract between all of them.

```
yallaops/
├── core/              # Go — control plane (the brain)
├── agent/             # Rust — runtime agent (Phase 4)
├── cli/               # Python — CLI tool
├── ai-generator/      # Python — LLM deployment generator (Phase 5)
├── proto/             # Protobuf definitions — source of truth for all APIs
├── infra/             # K8s manifests, Helm charts for deploying YallaOps itself
├── docs/              # ADRs and specs
├── docker-compose.yml # Local dev infrastructure
├── Justfile           # Task runner
└── buf.yaml           # Protobuf toolchain config
```

---

## 4. Technology decisions

### 4.1 Go — control plane

**Decision:** Go for the core control plane.

**Reasons:**
- Native concurrency model (goroutines) fits the event-driven release engine
- `client-go` is the standard K8s API library — written in Go, best supported in Go
- `helm/v3` embeds directly as a Go library — no shelling out to the Helm binary
- Same stack as ArgoCD, Flux, Tekton — the proven choice for this problem domain
- Single compiled binary — easy to ship and deploy

**Key libraries:**
- `google.golang.org/grpc` — gRPC server
- `github.com/jackc/pgx/v5` — Postgres driver
- `github.com/sqlc-dev/sqlc` — SQL → type-safe Go code generation
- `k8s.io/client-go` — Kubernetes API client
- `helm.sh/helm/v3` — Helm chart rendering
- `github.com/redis/go-redis/v9` — Redis client
- `log/slog` — structured logging (stdlib, no external dep)

---

### 4.2 Python — CLI and AI generator

**Decision:** Python for the CLI tool and AI deployment generator.

**Reasons:**
- `typer` gives a clean CLI interface with minimal boilerplate
- `grpcio` + `grpcio-tools` generate client stubs from the same `.proto` files as Go — always in sync
- LLM integration (for the AI generator) is most mature in Python
- Fast to iterate on — CLI commands change frequently during early dev

**Key libraries:**
- `typer` — CLI framework
- `httpx` — HTTP client (for health checks and webhooks)
- `grpcio` + `grpcio-tools` — gRPC client
- `rich` — terminal output formatting
- `pydantic` — config validation

---

### 4.3 Rust — runtime agent (Phase 4)

**Decision:** Rust for the optional runtime agent that runs inside clusters/nodes.

**Reasons:**
- Low resource footprint — the agent runs on every node, must be lightweight
- Memory safety without GC pauses — critical for a metrics/log collection process
- Excellent async networking with `tokio`

**Deferred to Phase 4.** Do not start this until the Go core and CLI are stable.

---

### 4.4 gRPC — API layer

**Decision:** gRPC with protobuf from day one for all API communication.

**Reasons:**
- Proto files in `proto/` are the single source of truth — Go server and Python CLI both generate from the same definitions, so they are always in sync by construction
- Strongly typed — no stringly-typed REST JSON to maintain
- Bidirectional streaming — useful for log tailing and real-time deployment status
- Future web UI can use gRPC-Web or a lightweight REST gateway

**What this means in practice:**
- Every new API endpoint starts as a `.proto` definition, not a Go handler
- Run `buf generate` to regenerate Go and Python stubs after any proto change
- Never bypass gRPC with direct function calls between CLI and core — keep the boundary clean

**No REST API.** If a REST endpoint is needed for webhooks or health checks, add a thin HTTP handler in `core/cmd/server/` that does not duplicate business logic.

---

### 4.5 PostgreSQL — database

**Decision:** PostgreSQL everywhere, via `pgx` + `sqlc`.

**Reasons:**
- YallaOps data is relational — releases belong to services, environments belong to releases, approvals belong to environments. SQL joins are the natural query model.
- `sqlc` generates type-safe Go from raw SQL — no ORM magic, full control over queries
- Runs anywhere: `docker run` locally, a pod in K8s, or managed (Neon, RDS, Supabase)
- ACID transactions — updating release status + writing an audit log entry must be atomic

**Why not etcd:** etcd is a key-value coordination store, not an application database. It has no SQL, no joins, a 1MB per-value limit, and is tightly coupled to Kubernetes. ArgoCD and Flux use it because they store their state as K8s CRDs (which live in etcd for free). YallaOps manages multiple runtimes — it cannot be K8s-only, so it needs a standalone database.

**Why not MongoDB:** YallaOps data has clear relations. MongoDB's flexibility is not needed and would require manual joins in Go application code.

**Why not SQLite:** No concurrent writes. Requires CGo in Go. Would need replacing before production — migration pain is not worth the zero-install convenience.

**Local dev:**
```bash
docker run -e POSTGRES_PASSWORD=yalla -e POSTGRES_DB=yallaops -p 5432:5432 postgres:16-alpine
```

**Migrations:** managed with `golang-migrate`. All migrations live in `core/db/migrations/`. Never edit a migration file after it has been committed — add a new one.

**Queries:** written as raw SQL in `core/db/queries/`. Run `sqlc generate` to regenerate Go code. Never write raw SQL strings in Go application code — always go through sqlc-generated functions.

---

### 4.6 Redis — event bus

**Decision:** Redis pub/sub for the internal notification center event bus.

**Reasons:**
- Already in the stack for caching — zero additional infrastructure
- Pub/sub is sufficient for YallaOps event volume — this is not a high-throughput message queue
- Simple to reason about — publish an event, subscribers receive it, send notifications

**Event flow:**
```
Release engine publishes event → Redis channel
        ↓
Notification center subscribes → fans out to:
  - Slack webhook
  - Email
  - Web dashboard (SSE or gRPC stream)
  - External webhooks
```

**Why not NATS:** NATS is excellent but adds an extra service to run. Redis is already required — use it.

**Why not RabbitMQ/Kafka:** Overkill for this use case. If event volume grows significantly, Redis Streams can replace pub/sub with zero infrastructure change.

---

### 4.7 Helm + Kustomize — deployment sources

**Decision:** Support Helm charts and Kustomize overlays as deployment sources, embedded as Go libraries.

**Why embedded:**
- `helm.sh/helm/v3` is a Go library — import it, call `RenderChart()`, get back plain K8s manifests, apply them with `client-go`. No shelling out to the Helm binary.
- `sigs.k8s.io/kustomize/api` same pattern — render overlays in-process.

**This means:** users can point a YallaOps release at a Helm chart or a Kustomize directory, not just raw YAML. Phase 3 adds Helm support, Phase 4 adds Kustomize.

---

## 5. What YallaOps does NOT do

- **Not a CI runner** — YallaOps does not build images or run tests. It consumes artifacts that CI has already produced.
- **Not a secret manager** — use Sealed Secrets, External Secrets, or Vault. YallaOps passes through whatever secrets configuration exists.
- **Not a service mesh** — networking between services is not YallaOps' concern.
- **Not GitOps** — YallaOps does not continuously sync Git state to a cluster. It manages discrete release promotion events.

---

## 6. Local development environment

Two services required locally: Postgres and Redis. Both run via Docker Compose.

```yaml
# docker-compose.yml (root of repo)
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_PASSWORD: yalla
      POSTGRES_DB: yallaops
    ports: ["5432:5432"]
    volumes: ["postgres_data:/var/lib/postgresql/data"]

  redis:
    image: redis:7-alpine
    ports: ["6379:6379"]

volumes:
  postgres_data:
```

Start with: `docker compose up -d`

---

## 7. Environment promotion rules

```
dev     → auto-deploy on release creation (no approval)
staging → requires at least one approval
prod    → requires approval + policy validation pass
```

Policy validation examples:
- Production must have ≥ 2 replicas
- Image digest in prod must match image digest that passed staging
- All environment variables present in staging must be present in prod
- Health checks must be defined

---

## 8. RBAC model

| Role | Permissions |
|---|---|
| developer | create releases, view status |
| reviewer | approve promotions |
| admin | all permissions + manage policies + rollback |

RBAC is enforced at the gRPC interceptor layer in `core/internal/api/`.

---

## 9. Decisions still open

| Decision | Status | Notes |
|---|---|---|
| Web UI framework | Not started | React or Vue, Phase 3+ |
| Auth provider | Not started | Likely OIDC (Dex or Auth0) |
| Managed Postgres for K8s | Not started | CloudNativePG or Neon |
| Observability | Not started | Prometheus metrics from Go core |

---

## 10. ADR index

Architecture Decision Records live in `docs/adr/`. Each significant decision gets its own file.

| ADR | Decision |
|---|---|
| ADR-001 | Monorepo structure |
| ADR-002 | gRPC over REST |
| ADR-003 | PostgreSQL over etcd/MongoDB/SQLite |
| ADR-004 | Redis pub/sub for event bus |
| ADR-005 | Helm embedded as Go library |
