# Architecture

## Entry Point

- `core/cmd/pet-core/main.go`: resolves the data dir and logs startup. Will wire window, server, state, tray, hotkeys.
- `settings/lib/main.dart`: initializes `window_manager` (900×620), runs `PetMacosApp` on macOS, `PetFluentApp` elsewhere.

## Core Layer

### core/internal/paths

- Per-OS data dir (`os.UserConfigDir()/ClaudePet`) and `config.json` path.
- Files: `paths.go`.

### core/assets

- `embed.go` exposes `Spritesheet` and `PetJSON` bytes.
- `pet.json`: cell 192×208, animations by row; `look` uses rows 9–10 in direction mode.

### core/internal/{config,sprite,state,server,window,tray,hotkeys,autostart}

- Stubs (`doc.go` only). Responsibilities as in [plan.md](plan.md) §3.

### settings/lib/core

- `core_client.dart`: `CoreClient` for `/health`, `/state`, `/preview`, `/reload` with `X-Pet-Token`.

### settings/lib/pages

- `pages.dart`: `SettingsPage` enum (general, animations, events, behavior, hotkeys, about) with titles.

## Feature Layer

- `settings/lib/ui_fluent/fluent_app.dart`: `NavigationView` + `NavigationPane`, empty page bodies.
- `settings/lib/ui_macos/macos_app.dart`: `MacosWindow` + `Sidebar`, empty content.
- `hooks/pet-hook.sh`, `hooks/pet-hook.ps1`: read port/token from the data dir, POST stdin to `/event?event=<Name>`, silent, exit 0.
- `hooks/settings.snippet.json` (macOS), `hooks/settings.snippet.windows.json` (Windows).

## Known Gaps

- No window, server, state machine, tray, hotkeys, autostart, or config code yet.
- Go dependencies not added.
- Settings pages have no content; no `ConfigRepository`; no «Pet is not running» banner.
- `CoreClient` has no port/token discovery.
- `tools/slice_sprites.py` is a placeholder; `packaging/` is empty.

## Roadmap

1. M1: Ebiten window with `idle` animation.
2. M2: HTTP server, token, state machine.
3. M3: tray, hotkeys, autostart.
4. M4–M5: settings pages and hook installer.
5. M6–M7: packaging, polish.
