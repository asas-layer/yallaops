# YallaOps

Open source multi-runtime release orchestration platform. Manages releases through dev → staging → prod with approval gates, across Kubernetes, Docker Compose, and Swarm.

Full decisions: @ARCHITECTURE.md
Full spec: @docs/SPEC.md

## Commands

```bash
docker compose up -d          # start postgres + redis
just migrate                  # run db migrations
just dev-core                 # start gRPC server
just proto                    # buf generate (regenerate all stubs)
just test                     # run all tests
just lint                     # golangci-lint + ruff + buf lint
```

## Rules

@.claude/rules/go.md
@.claude/rules/proto.md
@.claude/rules/git.md
@.claude/rules/db.md

## Current phase

Phase 1 — Go core release engine. Focus: proto definitions → postgres schema → release store → state machine → gRPC handlers. Do not touch agent/ or ai-generator/ yet.
