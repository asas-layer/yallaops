# Go rules

- No ORMs. All DB access via sqlc-generated functions in `core/internal/store/`. Never write raw SQL strings in Go.
- No global state. Inject all dependencies via struct constructors.
- Error wrapping always: `fmt.Errorf("release store: get: %w", err)` — never discard errors.
- No panics in library code. Only in `main()` for fatal startup.
- Logging: `slog` only. Never `fmt.Println` for operational output. Always structured key-value pairs.
- Context first: every function touching DB, network, or K8s takes `context.Context` as first arg.
- No logic in gRPC handlers — handlers call internal packages only.
- Tests: table-driven only. No single-assertion test functions.
- No shelling out to `helm` or `kubectl` — use embedded Go libraries.
