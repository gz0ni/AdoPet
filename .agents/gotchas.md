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

### macOS RSS overstates pet-core memory (2026-11-24)

- `ps` RSS stays ~115 MB for the Cocoa backend while `vmmap` physical footprint is ~32 MB.
- Root cause: RSS counts shared, clean pages of system frameworks (AppKit, Core Animation).
- Avoid: judge the 40 MB budget on macOS by `vmmap --summary` physical footprint (what Activity Monitor shows as Memory).

### Agent cannot screenshot a bare pet-core binary on macOS (2026-11-24)

- `screencapture -l <id>` failed (`could not create image from window`); computer-use cannot grant pet-core because it has no bundle.
- Root cause: the shell lacks Screen Recording permission; computer-use works only with bundled apps.
- Avoid: check the window with a CGWindowList script (bounds, layer, onscreen) and ask the user for visual checks, or wrap pet-core in a `.app` bundle.

### macOS CPU at startup looks high (2026-11-24)

- `ps` right after launch showed ~10% CPU; after ~10 s it settles to 0.3–1.6%.
- Avoid: measure CPU and memory after at least 10 s, average several samples.

### Unix socket bind fails with a long path (2026-11-24)

- `net.Listen("unix", <scratchpad>\pet.sock)` on Windows failed with `bind: invalid argument`.
- Root cause: `sun_path` limit (108 bytes on Windows, 104 on macOS); the scratchpad path was longer.
- Avoid: keep socket paths short (`<dataDir>/pet.sock` is ~50 bytes); test sockets under `%TEMP%\<short>`.

### Ebiten DirectX backend refuses transparent windows (2026-11-24)

- `EBITENGINE_GRAPHICS_LIBRARY=directx` with `ScreenTransparent: true` exits with `ui: DirectX is not available with a transparent window`.
- Root cause: Ebiten on Windows supports transparency only through OpenGL.
- Avoid: do not try DirectX to cut memory; the GL driver footprint (~100 MB on NVIDIA) is the cost of an Ebiten transparent window.

### Vsync-on Ebiten loop costs ~4% of a core while idle (2026-11-24)

- With vsync on and `SetScreenClearedEveryFrame(false)` CPU stayed at ~4.3% of one core even with skipped draws.
- Root cause: the loop still presents at the monitor refresh rate.
- Avoid: use `FPSModeVsyncOffMinimum` + `ScheduleFrame` at the animation FPS (→ ~0%). It is deprecated, so keep the `//nolint:staticcheck` directives.

### golangci-lint errcheck flags cleanup calls (2026-11-24)

- `f.Close()` / `os.Remove(tmp)` on error paths failed errcheck.
- Avoid: assign to `_` or use a deferred cleanup keyed on the named error.

### Redirecting run logs with relative paths escapes the repo (2026-11-24)

- `../../pet-run.log` from `core/` landed in `D:\Work\Project`, outside the repo.
- Avoid: write run logs to the session scratchpad with an absolute path.

### Very long multi-heredoc Bash commands fail to parse (2026-11-24)

- One Bash call writing ~14 files via heredocs failed with `unexpected EOF while looking for matching '`; nothing was written.
- Root cause: shell parsing of one huge command with many heredocs and quotes.
- Avoid: write many files with the Write tool, or split into small Bash calls.

### golangci-lint is not on PATH (2026-11-24)

- Installed with `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest` (v2.15.0), but `$(go env GOPATH)/bin` is not on PATH.
- Avoid: call it as `"$(go env GOPATH)/bin/golangci-lint" run ./...`.
