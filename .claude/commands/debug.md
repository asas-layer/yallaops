---
allowed-tools: Bash(git diff:*), Bash(go test:*), Bash(cat:*)
description: Debug a specific issue — $ARGUMENTS
---

## Task

Debug this issue: $ARGUMENTS

Do the following in order:

1. **Reproduce** — write the minimal command or test that triggers the issue
2. **Locate** — identify the exact file and function where the problem occurs
3. **Explain** — explain what is happening vs what should happen
4. **Fix** — implement the fix, keeping it minimal and focused
5. **Verify** — run the relevant test or command to confirm it's fixed
6. **Check for recurrence** — are there other places in the codebase with the same pattern that could have the same bug?

Do not guess. If you need to read a file to understand the issue, read it first.
