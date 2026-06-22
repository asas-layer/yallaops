---
description: Plan and spec a new feature before writing any code
---

## Task

I want to build: $ARGUMENTS

Before writing any code, do the following:

1. **Ask me clarifying questions** about anything ambiguous — edge cases, error states, API shape, where it fits in the existing architecture. Check @ARCHITECTURE.md and the current phase in @CLAUDE.md first so you don't ask things already decided.

2. **Write a spec** covering:
   - What this feature does (one paragraph)
   - Proto changes needed (new messages, new RPC methods)
   - Database changes needed (new tables or columns, migration required?)
   - Go packages affected (`core/internal/release/`, `core/internal/api/`, etc.)
   - Acceptance criteria — a list of things that must be true for this to be "done"
   - Out of scope — what this feature explicitly does NOT do

3. **Propose the branch name** following our naming convention.

4. **Ask for my approval** before writing any code.

Only start coding after I say "looks good" or "go ahead".
