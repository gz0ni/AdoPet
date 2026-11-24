# Project Context

Claude Pet (AdoPet): a lightweight desktop pet (Chibi Ado) that reacts to Claude Code activity. Two processes: **pet-core** (Go) is always running, draws the pet, owns the tray icon and global hotkeys, and receives events from Claude Code hooks over a local HTTP API; **settings** (Flutter) is started on demand from the tray and exits when closed. Scope: Claude Code in the terminal and the Desktop Code tab (hooks from `~/.claude/settings.json`); Claude chat exposes no events, so the pet is decorative there. Full plan: [plan.md](plan.md).

## Platforms

Windows 10/11 and macOS 12+.

## Toolchain

- core: Go 1.26.5 (`core/go.mod`), Ebiten v2 (planned); cgo required on macOS.
- settings: Flutter 3.44.6 stable, Dart SDK ^3.12.2, pub.
- Hooks: POSIX sh + curl (macOS), PowerShell (Windows).
- CI: GitHub Actions, matrix `windows-latest` / `macos-latest`.
- Lint: `go vet`, `golangci-lint` (v2 config), `flutter analyze`.

## Package / Repository

- `name:` core module `github.com/gz0ni/claude-pet/core`; Flutter package `claude_pet_settings`
- `version:` core `0.0.0-dev` (ldflag `main.version`); settings `1.0.0+1`
- `org:` / owner: gz0ni; Flutter org `dev.claudepet`

## Dependencies

core (planned, not yet in go.mod):
- `github.com/hajimehoshi/ebiten/v2`: transparent floating window, rendering.
- `fyne.io/systray`: tray / menu bar (`RunWithExternalLoop`).
- `golang.design/x/hotkey`: global hotkeys.
- `github.com/fsnotify/fsnotify`: watch config.json.
- stdlib `net/http`, `encoding/json`, `embed`, `log/slog`.

settings (added):
- `fluent_ui` (Windows UI), `macos_ui` (macOS UI), `window_manager`, `http` (core API), `path_provider`.
- Planned: `hotkey_manager` (recorder only), `file_picker` (pet import).

## Design Basis

- Implementation plan: [plan.md](plan.md).
- Sprites: `core/assets/spritesheet.png` (copy of `core/assets/final/spritesheet-extended.png`), 1536×2288, 8×11 grid, cell 192×208, described by `core/assets/pet.json`.
- Settings UI: native look per OS — Fluent (`NavigationView`) on Windows, `MacosWindow` + `Sidebar` on macOS.

## Status

M0 (skeleton) done: repo layout, Go module with stub packages, embedded assets, hook scripts and snippets, Flutter app with platform-split navigation, CI workflow. M1 (pet on screen) not started.
