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

- [ ] M1: pet on screen — embedded sheet, `idle` loop, transparent undecorated floating window, drag, saved position, check RSS (prototype on Windows and macOS first)
- [ ] M2: events — HTTP server + token, state machine, hook scripts, `curl` test, then a real Claude Code session (server reads the event name from `?event=`)
- [ ] M3: tray + hotkeys + autostart (confirm systray `RunWithExternalLoop` with Ebiten on macOS early)
- [ ] M4: settings app — config repo, General + Animations pages on both UI kits, `/preview`, `/reload`
- [ ] M5: Claude events page + hook installer (merge into `~/.claude/settings.json` with backup, uninstall)
- [ ] M6: packaging — Inno Setup, dmg, CI artifacts on tag
- [ ] M7: polish — wander/follow-cursor, sounds, DND, multi-monitor, HiDPI, custom pet import
- [ ] Verify Claude Code hook event names and payload fields against current docs before M2
- [ ] Install `golangci-lint` locally and run it
- [ ] Choose Flutter state management (ChangeNotifier vs Riverpod) and record it
- [ ] Translate settings UI page titles to Russian (currently English)
- [ ] Implement `tools/slice_sprites.py` (optional)
- [ ] First commit of the skeleton (on user request)

## In Progress

## Done

- [x] Create `.agents/` structure and `AGENTS.md` (2026-11-24)
- [x] M0: skeleton — layout, Go module, Flutter project, hooks, CI, README; `go vet`/`go build`, `flutter analyze`/`flutter test`, `flutter build windows --debug` pass (2026-11-24)
