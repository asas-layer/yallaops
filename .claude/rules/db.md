# Database rules

- PostgreSQL only. No SQLite, no MongoDB, no etcd for application state.
- Never edit a migration file after it has been committed — always add a new one.
- Migration files: `core/db/migrations/NNNN_description.up.sql` and `.down.sql`.
- Query files: `core/db/queries/*.sql`. Run `just sqlc` after any change to regenerate Go code.
- Multi-table writes must use transactions — updating release status + audit log = one transaction.
- Never write raw SQL strings in Go. All queries go through sqlc-generated functions.
- Connection string from env var `DATABASE_URL` only. Never hardcode.
