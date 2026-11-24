# TODO

Task checklist. The agent maintains this file: creates items, marks them done, keeps it honest. Work from here — do not lose tasks between sessions.

## Format

- `[ ]` pending task — add due date or note if needed
- `[x]` done task — keep for history, do not delete immediately
- `[~]` in progress — who/what is working on it

## Rules

- Add a task the moment it is discovered, even if vague.
- Mark `[x]` only when actually verified done (checks passed), not by intent.
- Review this file at the end of every session: close finished items, add follow-ups.

---

## Open

- [~] M1: pet on screen — Windows and macOS done and checked by the user; branch `m1-pet-window` not merged yet
- [ ] Handle `WM_DPICHANGED` in the Win32 backend (M7)
- [x] M1 on macOS: build, transparency, idle loop, drag, saved position checked on Apple Silicon, macOS 27 (2026-11-24)
- [ ] Run the Windows build on the Windows machine after Ebiten removal (`go mod tidy` changed `go.sum`)
- [ ] Handle `backingScaleFactor` changes on macOS (M7)
- [ ] M2: events — HTTP over `<dataDir>/pet.sock` (no port/token), stale-socket check + single instance, state machine with all animations, hooks and `CoreClient` on the socket, `curl` test, then a real Claude Code session (server reads the event name from `?event=`); verify AF_UNIX on macOS
- [ ] M3: tray + hotkeys + autostart (confirm systray `RunWithExternalLoop` with the Cocoa `[NSApp run]` loop on macOS early)
- [ ] M4: settings app — config repo, General + Animations pages on both UI kits, `/preview`, `/reload`
- [ ] M5: Claude events page + hook installer (merge into `~/.claude/settings.json` with backup, uninstall)
- [ ] M6: packaging — Inno Setup, dmg, CI artifacts on tag
- [ ] M7: polish — wander/follow-cursor, sounds, DND, multi-monitor, HiDPI, custom pet import
- [ ] Verify Claude Code hook event names and payload fields against current docs before M2
- [ ] Choose Flutter state management (ChangeNotifier vs Riverpod) and record it
- [ ] Translate settings UI page titles to Russian (currently English)
- [ ] Implement `tools/slice_sprites.py` (optional)
- [ ] First commit of the skeleton (on user request)

## In Progress

## Done

- [x] Cocoa `NSPanel` + `CALayer` backend on macOS replaces Ebiten; footprint 130 → 32 MB, CPU ≈ 0.1%, no Dock icon; Ebiten removed from `go.mod` (2026-11-24)
- [x] golangci-lint v2.15.0 installed on the Mac, 0 issues (2026-11-24)
- [x] Win32 layered-window backend for Windows replaces Ebiten there; memory 108 → 18 MB (2026-11-24)

- [x] Create `.agents/` structure and `AGENTS.md` (2026-11-24)
- [x] Install `golangci-lint` v2.15.0 locally; `golangci-lint run` reports 0 issues (2026-11-24)
- [x] M0: skeleton — layout, Go module, Flutter project, hooks, CI, README; `go vet`/`go build`, `flutter analyze`/`flutter test`, `flutter build windows --debug` pass (2026-11-24)
