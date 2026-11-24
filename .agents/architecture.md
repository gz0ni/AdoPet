# Architecture

## Entry Point

- `core/cmd/pet-core/main.go`: sets a 48 MiB Go memory limit, loads config (defaults on missing file; invalid file → defaults and no saving), parses `pet.json`, decodes the sheet, runs `window.Run`, saves the position on drag end. Will wire server, state, tray, hotkeys.
- `settings/lib/main.dart`: initializes `window_manager` (900×620), runs `PetMacosApp` on macOS, `PetFluentApp` elsewhere.

## Core Layer

### core/internal/paths

- Per-OS data dir (`os.UserConfigDir()/ClaudePet`) and `config.json` path.
- Files: `paths.go`.

### core/assets

- `embed.go` exposes `Spritesheet` and `PetJSON` bytes.
- `pet.json`: cell 192×208, animations by row; `look` uses rows 9–10 in direction mode.

### core/internal/config

- `config.go`: `Config` subset (`version`, `pet`, `scale`, `position`, `alwaysOnTop`), `Default`, `Load` (fills missing keys from defaults), `Save` (deep-merges known fields over the raw JSON so unknown keys survive, `UseNumber`, temp file + rename).

### core/internal/sprite

- `sprite.go`: `ParsePet`, `FrameRect` (single row or `rows` direction mode), `FrameAt` (loop / one-shot timing), `ScaleSheet` (per-cell CatmullRom into `*image.RGBA`, no bleeding between cells).

### core/internal/window

- `options.go`: `Options` shared by both backends.
- `window_windows.go` + `win32_windows.go`: Windows backend without GPU. `WS_POPUP` with `WS_EX_LAYERED | WS_EX_TOOLWINDOW | WS_EX_NOACTIVATE` (+ `WS_EX_TOPMOST`), per-monitor DPI v2, 32-bit top-down DIB section updated by `UpdateLayeredWindow(ULW_ALPHA)`. `SetTimer` at the animation FPS redraws only on frame change. Drag through `WM_NCHITTEST → HTCAPTION`; fully transparent pixels pass clicks through. Position saved on `WM_EXITSIZEMOVE`. Monitors from `EnumDisplayMonitors`, primary first; coordinates relative to the work area (taskbar excluded), config offsets in DIP.
- `pixels.go`: `blitBGRA` (premultiplied RGBA → BGRA for the DIB).
- `window_other.go` (`!windows`): Ebiten game. Pre-scales the sheet to `scale × deviceScaleFactor`, `Layout` returns physical cell size (crisp on HiDPI), window size is in DIP. Transparent, undecorated, floating, `SkipTaskbar`, `InitUnfocused`. `FPSModeVsyncOffMinimum` + `ScheduleFrame` ticker at the animation FPS; draws only on frame change. Left-drag moves the window; on release reports the anchor-relative position and monitor index.
- `placement.go`: `place` (anchor offset → monitor coords, clamped on-screen) and `relative` (inverse).

### core/internal/{state,server,tray,hotkeys,autostart}

- Stubs (`doc.go` only). Responsibilities as in [plan.md](plan.md) §3.

### settings/lib/core

- `core_client.dart`: `CoreClient` for `/health`, `/state`, `/preview`, `/reload` with `X-Pet-Token`.

### settings/lib/pages

- `pages.dart`: `SettingsPage` enum (general, animations, events, behavior, hotkeys, about) with titles.

## Feature Layer

- `settings/lib/ui_fluent/fluent_app.dart`: `NavigationView` + `NavigationPane`, empty page bodies.
- `settings/lib/ui_macos/macos_app.dart`: `MacosWindow` + `Sidebar`, empty content.
- `hooks/pet-hook.sh`, `hooks/pet-hook.ps1`: read port/token from the data dir, POST stdin to `/event?event=<Name>`, silent, exit 0. Still on port/token; M2 moves them to `pet.sock`.
- `hooks/settings.snippet.json` (macOS), `hooks/settings.snippet.windows.json` (Windows).

## Known Gaps

- No server, state machine, tray, hotkeys, autostart yet; only `idle` is played.
- Windows: `WM_DPICHANGED` is not handled; the pet keeps its pixel size on a monitor with another DPI.
- Device scale factor is read once at startup; moving to a monitor with another DPI is not handled (M7).
- Click-through and fsnotify config watch not implemented.
- macOS (Ebiten backend) only checked with `GOOS=darwin go vet`; not run, RSS unknown.
- Settings pages have no content; no `ConfigRepository`; no «Pet is not running» banner.
- `CoreClient` uses port + `X-Pet-Token`; M2 switches it to the unix socket.
- `tools/slice_sprites.py` is a placeholder; `packaging/` is empty.

## Roadmap

1. M1: run on macOS.
2. M2: HTTP server on a unix socket, state machine.
3. M3: tray, hotkeys, autostart.
4. M4–M5: settings pages and hook installer.
5. M6–M7: packaging, polish.
