# Git rules

- Never commit directly to `main`. Always work on a branch.
- Branch naming: `feat/<name>`, `fix/<name>`, `chore/<name>`.
- Branch lifetime: max 1-2 days. If it lives longer, it's too big — split it.
- One logical change per PR. No "misc fixes" PRs.
- Commit messages: conventional commits format — `feat:`, `fix:`, `chore:`, `docs:`, `test:`, `refactor:`.
- Squash merge only into main. Linear history.
- Never commit secrets, API keys, or internal URLs. This repo is public.
- Run `just lint` and `just test` before pushing.
