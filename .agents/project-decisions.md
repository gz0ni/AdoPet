# Project Decisions

Recorded decisions and the reasoning behind them, so future work does not silently contradict them. Add an entry whenever a non-trivial choice is made. Latest first.

## Format

```
### <Decision title> (YYYY-MM-DD)

- Decision: <what was chosen, one sentence>
- Reason: <why this and not the alternatives>
- Alternatives considered: <what was rejected and why>
- Consequences: <what this decision affects>
```

---

<!-- New entries go here, newest first. -->

### Own Cocoa window on macOS, Ebiten removed (2026-11-24)

- Decision: macOS backend is a cgo/Objective-C borderless non-activating `NSPanel` whose layer shows one pre-built `CGImage` per frame; Ebiten is dropped from the module.
- Reason: Ebiten (Metal) used ~130 MB physical footprint and showed a Dock icon; the Cocoa window runs at ~32 MB, ≈ 0.1% CPU, no Dock icon, and the binary shrank from 11 to 5.3 MB. With both OS on native backends Ebiten had no users left.
- Alternatives considered: keep Ebiten and relax the 40 MB target; `NSView drawRect` with a CPU backing store (more code, no gain over layer contents).
- Consequences: macOS builds need cgo and a Mac or the macOS CI runner; `GOOS=darwin` cross-vet from Windows no longer covers the window package; tray in M3 must share the `[NSApp run]` main-thread loop; the "redraw on demand with minimum FPS mode" decision is obsolete.

### Own Win32 layered window on Windows, Ebiten only on macOS (2026-11-24)

- Decision: `internal/window` has two backends behind one `Run(Options)`: Win32 (`UpdateLayeredWindow` + DIB, no GPU) on Windows, Ebiten elsewhere.
- Reason: Ebiten transparent windows force OpenGL on Windows and the GL driver took ~100 MB; the layered window runs at ~18 MB and ≈ 0% CPU, and free per-pixel hit-testing and system drag come with it.
- Alternatives considered: keep Ebiten and relax the 40 MB target; rewrite in Rust (same driver cost with a GPU API, only ~5–10 MB gain without one).
- Consequences: Windows window code is pure `x/sys/windows` syscalls; sprite, config, and placement stay shared; tray and hotkeys on Windows can share the Win32 message loop; DPI changes need `WM_DPICHANGED` handling later.

### Local API over a unix socket instead of a TCP port (2026-11-24)

- Decision: pet-core serves HTTP on `<dataDir>/pet.sock` (AF_UNIX on both OS); no port, no `port` file, no token.
- Reason: the OS enforces access (per-user directory and file permissions), nothing listens on the network stack, browsers cannot reach it, no port conflicts. A localhost port needs a token to compensate for being reachable by any local process or page.
- Alternatives considered: `127.0.0.1` port + `X-Pet-Token` (simpler tooling, weaker by design); named pipes on Windows (would need a different transport per OS; AF_UNIX works on Windows 10 1803+, checked with `curl.exe` and Dart).
- Consequences: Windows minimum is 10 1803; hook scripts and `CoreClient` change in M2; `server.port` leaves the config schema; socket path must stay under 104 bytes.

### Redraw on demand with the deprecated minimum FPS mode (2026-11-24) — superseded

- Decision: `ebiten.SetFPSMode(FPSModeVsyncOffMinimum)` plus a ticker calling `ScheduleFrame` at the animation FPS.
- Reason: drops idle CPU from ~4.3% to ~0% of a core; input still wakes the loop for dragging.
- Alternatives considered: vsync on + `SetScreenClearedEveryFrame(false)` (still ~4%), vsync off (unbounded FPS).
- Consequences: relies on a deprecated API; revisit on Ebiten upgrades.

### Pre-scale the sheet for HiDPI at startup (2026-11-24)

- Decision: scale cells by `scale × deviceScaleFactor` once, window size in DIP, `Layout` in physical pixels.
- Reason: crisp pet on 125–200% displays without keeping the full 14 MB sheet.
- Alternatives considered: scaling in `Draw` (blurry, keeps the original).
- Consequences: DPI changes after start are not handled until M7.

### Config holds only the fields in use (2026-11-24)

- Decision: `config.Config` covers `version`, `pet`, `scale`, `position`, `alwaysOnTop`; other schema keys are added in their milestones and survive saves as unknown keys.
- Reason: avoid inventing defaults for `eventMap`, hotkeys, etc. before they are designed.
- Alternatives considered: full plan schema now.
- Consequences: later milestones extend the struct and `Default`.

### Position stored relative to an anchor corner (2026-11-24)

- Decision: `position.x/y` are offsets from the `anchor` corner of the monitor (pet corner aligned), clamped on-screen; drag keeps the existing anchor.
- Reason: matches the plan default `{40, -40, bottom-left}` and survives resolution changes.
- Alternatives considered: absolute coordinates.
- Consequences: settings «reset position» writes the default anchor offsets.

### Hook event name passed in the query string (2026-11-24)

- Decision: hook scripts POST the raw hook JSON as the body and the event name as `/event?event=<Name>`.
- Reason: scripts stay trivial — no JSON parsing or rewriting in sh/PowerShell.
- Alternatives considered: injecting `event` into the JSON body (needs jq or JSON handling in the hook).
- Consequences: the M2 `/event` handler must read `event` from the query; the payload's `hook_event_name` can serve as a fallback.

### Embed only spritesheet.png and pet.json (2026-11-24)

- Decision: `core/assets/embed.go` embeds just `spritesheet.png` and `pet.json`; the other asset folders stay as source material.
- Reason: keep the binary small; work folders are not needed at runtime.
- Alternatives considered: embedding the whole `assets/` directory.
- Consequences: replacing the default pet means updating these two files.

### Two processes: Go core + Flutter settings (2026-11-24)

- Decision: always-running pet in Go (native window per OS), settings UI in Flutter launched on demand.
- Reason: small idle footprint for the pet; native-looking settings on both OS.
- Alternatives considered: see [plan.md](plan.md).
- Consequences: they communicate only through config.json and the local HTTP API.
