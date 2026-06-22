# Proto rules

- All .proto files live in `proto/`. Package: `yallaops.v1`.
- Every new API feature starts as a proto definition — never write a Go handler first.
- After any proto change: run `just proto` before touching Go or Python code.
- Field names: snake_case. Enum values: SCREAMING_SNAKE_CASE prefixed with enum name (e.g. `RELEASE_STATUS_DRAFT`).
- Never delete or renumber fields — add new ones, mark old ones deprecated with a comment.
- Run `buf lint` before committing any proto change.
