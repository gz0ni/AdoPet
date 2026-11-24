# Gotchas

Traps and pitfalls the agent has stepped on in this project. Log a gotcha immediately whenever something does not build, run, or pass — right after the incident, while context is fresh, do not wait for session end. Check this file before repeating similar work. One entry per trap.

## Format

```
### <Short trap title> (YYYY-MM-DD)

- What happened / what was wrong
- Root cause
- How to avoid next time
```

---

<!-- New entries go here, newest first. -->

### Very long multi-heredoc Bash commands fail to parse (2026-11-24)

- One Bash call writing ~14 files via heredocs failed with `unexpected EOF while looking for matching '`; nothing was written.
- Root cause: shell parsing of one huge command with many heredocs and quotes.
- Avoid: write many files with the Write tool, or split into small Bash calls.

### golangci-lint is not installed locally (2026-11-24)

- `golangci-lint run` is unavailable on the dev machine; only `go vet` was run.
- Root cause: tool not on PATH.
- Avoid: install it (`go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`) or rely on `go vet` and say so.
