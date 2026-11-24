# Claude Desktop Pet — Implementation Plan

A lightweight desktop pet (Chibi Ado) that reacts to Claude Code activity.
Two processes:

- **pet-core** (Go): always running, small footprint. Draws the pet, owns the tray/menu-bar icon and global hotkeys, receives events from Claude Code hooks.
- **settings** (Flutter): started on demand from the tray and exits when closed. Uses `fluent_ui` on Windows and `macos_ui` on macOS.

Targets: Windows 10/11 and macOS 12+.

Scope: Claude Code in the terminal and the Code tab of Claude Desktop (both read hooks from `~/.claude/settings.json`). Claude chat (claude.ai and the Desktop Chat tab) exposes no events. There the pet is decorative only.

---

## 1. Repository layout

```
claude-pet/
├── core/                    # Go module: pet-core
│   ├── cmd/pet-core/main.go
│   ├── internal/
│   │   ├── config/          # load/save/watch config.json, defaults, schema version
│   │   ├── sprite/          # spritesheet loading, frame slicing, animation player
│   │   ├── state/           # event → pet state machine
│   │   ├── server/          # HTTP API over a unix socket
│   │   ├── window/          # Ebiten game loop, transparent window, drag, click-through
│   │   ├── tray/            # tray / menu-bar icon and menu
│   │   ├── hotkeys/         # global hotkeys
│   │   ├── autostart/       # launch at login (registry / LaunchAgent)
│   │   └── paths/           # per-OS config & data dirs
│   └── assets/              # embedded spritesheet + pet.json (go:embed)
├── settings/                # Flutter app
│   ├── lib/
│   │   ├── main.dart        # picks FluentApp or MacosApp by Platform
│   │   ├── core/            # models, config repo, core API client
│   │   ├── pages/           # shared page logic (state/controllers)
│   │   ├── ui_fluent/       # Windows widgets
│   │   └── ui_macos/        # macOS widgets
│   └── pubspec.yaml
├── hooks/
│   ├── pet-hook.sh          # macOS/Linux
│   ├── pet-hook.ps1         # Windows
│   └── settings.snippet.json
├── tools/
│   └── slice_sprites.py     # optional: chroma check, per-row frame export, previews
├── packaging/
│   ├── windows/             # Inno Setup script (or MSIX manifest)
│   └── macos/               # Info.plist, entitlements, dmg script
├── .github/workflows/       # build matrix: windows-latest, macos-latest
└── README.md
```

---

## 2. Sprite assets

Source: `Chibi_Ado_Work.zip` → `ado-run/final/spritesheet-extended.png`.

- 1536×2288 RGBA, **8 columns × 11 rows**, cell **192×208**, already transparent (chroma `#FF00FF` removed, 0 residue).
- Unused cells are empty; frame counts per row:

| Row | State | Frames | Suggested use |
|---|---|---|---|
| 0 | `idle` | 6 | default |
| 1 | `running-right` | 8 | walk right |
| 2 | `running-left` | 8 | walk left |
| 3 | `waving` | 4 | needs attention / greeting |
| 4 | `jumping` | 5 | task done |
| 5 | `failed` | 8 | tool error / failed tests |
| 6 | `waiting` | 6 | waiting for user (permission) |
| 7 | `running` | 6 | working (tool use) |
| 8 | `review` | 6 | reading / thinking |
| 9 | `look-000-to-157.5` | 8 | follow cursor (directions 0°–157.5°) |
| 10 | `look-180-to-337.5` | 8 | follow cursor (directions 180°–337.5°) |

Ship `assets/pet.json` describing the sheet so another pet can be dropped in later:

```json
{
  "name": "Chibi Ado",
  "sheet": "spritesheet.png",
  "cell": { "w": 192, "h": 208 },
  "animations": {
    "idle":          { "row": 0, "frames": 6, "fps": 6,  "loop": true },
    "running-right": { "row": 1, "frames": 8, "fps": 10, "loop": true },
    "running-left":  { "row": 2, "frames": 8, "fps": 10, "loop": true },
    "waving":        { "row": 3, "frames": 4, "fps": 6,  "loop": true },
    "jumping":       { "row": 4, "frames": 5, "fps": 10, "loop": false, "next": "idle" },
    "failed":        { "row": 5, "frames": 8, "fps": 8,  "loop": false, "next": "idle" },
    "waiting":       { "row": 6, "frames": 6, "fps": 6,  "loop": true },
    "running":       { "row": 7, "frames": 6, "fps": 10, "loop": true },
    "review":        { "row": 8, "frames": 6, "fps": 6,  "loop": true },
    "look":          { "rows": [9, 10], "frames": 16, "mode": "direction" }
  }
}
```

Memory: decode the sheet once (≈14 MB RGBA at full size). Pre-scale to the configured display size on load (e.g. 50% → ≈3.5 MB) and drop the original.

---

## 3. pet-core (Go)

### Dependencies
- `github.com/hajimehoshi/ebiten/v2`: window + rendering (transparent, undecorated, floating, mouse passthrough).
- `fyne.io/systray`: tray / menu bar. Use `systray.RunWithExternalLoop` so it coexists with Ebiten's main-thread loop.
- `golang.design/x/hotkey`: global hotkeys (on macOS must be registered on the main thread; use `mainthread`).
- `github.com/fsnotify/fsnotify`: watch config.json.
- stdlib `net/http`, `encoding/json`, `embed`, `log/slog`.

### Window
- Windows (decided 2026-11-24): own Win32 layered window (`UpdateLayeredWindow` + DIB, no GPU), see `project-decisions.md`. The Ebiten notes below apply to macOS.
- Ebiten options: `ScreenTransparent: true`, `WindowDecorated(false)`, `WindowFloating(true)`, `SetRunnableOnUnfocused(true)`, `SetTPS` low (e.g. 30), `SetVsyncEnabled(false)`.
- Skip redraw when the frame hasn't changed (`ebiten.SetScreenClearedEveryFrame(false)` + dirty flag) to keep CPU near 0.
- Drag to move (left button); position saved to config.
- Optional click-through mode (`SetWindowMousePassthrough`), toggled by a hotkey.
- Hide from taskbar/Dock: macOS `LSUIElement=1` in Info.plist; Windows: tool-window style (small platform shim via `golang.org/x/sys/windows` if Ebiten doesn't expose it).

### State machine (`internal/state`)
Priority-based; higher priority wins, transient states return to the base state.

| Hook event | Condition | Pet state | Priority |
|---|---|---|---|
| `SessionStart` | — | `waving` (2 s) → `idle` | low |
| `UserPromptSubmit` | — | `review` | mid |
| `PreToolUse` | tool ∈ Read/Grep/Glob/WebFetch/WebSearch | `review` | mid |
| `PreToolUse` | tool ∈ Edit/Write/MultiEdit/Bash/… | `running` | mid |
| `PostToolUse` | error / non-zero exit | `failed` (one shot) | high |
| `Notification` | permission / idle prompt | `waiting` (+ optional sound) | highest |
| `Stop` / `SubagentStop` | — | `jumping` (one shot) → `idle` | high |
| `SessionEnd` | — | `idle` | low |
| no events for N min | — | `idle`, optional random walk | — |

- Track sessions by `session_id` so several Claude Code sessions don't fight; show the most urgent state across sessions.
- Every mapping is overridable from config (`eventMap`), and disabled animations fall back to `idle`.

### HTTP API (`internal/server`) over a unix domain socket
| Method | Path | Purpose |
|---|---|---|
| `POST` | `/event` | hook payload (raw Claude Code hook JSON + `event` name) |
| `POST` | `/preview` | `{ "animation": "failed" }`: play once (from settings) |
| `POST` | `/reload` | re-read config |
| `GET` | `/state` | current state, sessions, version (for settings UI) |
| `GET` | `/health` | liveness |

- Transport: HTTP over an AF_UNIX socket at `<dataDir>/pet.sock` on both OS (`net.Listen("unix", …)` + `http.Serve`). Windows 10 1803+ supports AF_UNIX. No TCP port, no `port` file, no token.
- Access control comes from the OS: the data dir is per-user (`0700` on macOS, user-only ACL of `%APPDATA%` on Windows), the socket file is `0600`. Other users, other machines, and browser pages cannot reach it.
- Startup: if `pet.sock` exists, try to connect. Success → another pet-core is running, exit (single instance for free). Failure → stale file, remove it and listen. Remove the socket on exit.
- Path limit: `sun_path` is 108 bytes on Windows, 104 on macOS. Check the length at startup and log a clear error; a long user profile path is the realistic failure.
- Clients: hooks use `curl --unix-socket` (curl in Windows has `UnixSockets` built in, macOS curl too); settings use Dart `HttpClient` with `connectionFactory` → `Socket.startConnect(InternetAddress(path, type: InternetAddressType.unix), 0)`. The URL host is a dummy `localhost`.
- Verified on Windows 11 (2026-11-24): Go server on AF_UNIX, `curl.exe --unix-socket` POST, Dart raw `Socket` and `HttpClient` both got `200`. macOS still to check in M2.
- Hook must never block Claude: handler returns `204` immediately and processes the event asynchronously.

### Tray / menu
Show/Hide pet · Pause reactions · Do not disturb · Settings… (spawns the Flutter app; one instance only) · Start at login ✓ · Quit.

### Hotkeys (defaults, configurable)
- `Ctrl/Cmd+Alt+P`: show/hide pet
- `Ctrl/Cmd+Alt+D`: do not disturb
- `Ctrl/Cmd+Alt+T`: click-through on/off
- `Ctrl/Cmd+Alt+,`: open settings

### Config (`internal/config`)
Location: Windows `%APPDATA%\ClaudePet\config.json`, macOS `~/Library/Application Support/ClaudePet/config.json`.

```json
{
  "version": 1,
  "pet": "chibi-ado",
  "scale": 0.5,
  "position": { "x": 40, "y": -40, "anchor": "bottom-left", "display": 0 },
  "alwaysOnTop": true,
  "clickThrough": false,
  "animations": { "idle": true, "running-left": true, "running-right": true, "look": true, "...": true },
  "eventMap": { "Notification": "waiting", "Stop": "jumping" },
  "behavior": { "wander": false, "followCursor": true, "idleAfterSec": 120 },
  "sound": { "onWaiting": true, "volume": 0.6 },
  "hotkeys": { "toggle": "Ctrl+Alt+P", "dnd": "Ctrl+Alt+D", "clickThrough": "Ctrl+Alt+T", "settings": "Ctrl+Alt+Comma" },
  "autostart": true
}
```
- Unknown keys preserved; missing keys filled from defaults; atomic write (temp file + rename).
- fsnotify + `/reload` both trigger re-apply.

### Autostart
- Windows: `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`.
- macOS: LaunchAgent plist in `~/Library/LaunchAgents/` (or `SMAppService` later when signed).

### Memory/CPU targets
- Idle RSS ≤ 40 MB, CPU ≈ 0% when idle, < 2% while animating.
- Measure with Task Manager / Activity Monitor and `GODEBUG=gctrace=1`; set `debug.SetMemoryLimit` (e.g. 48 MiB) and `GOGC=50` if needed.

---

## 4. Claude Code hooks

`hooks/pet-hook.sh`:
- Reads hook JSON from stdin, `curl -s -m 1 --unix-socket "<dataDir>/pet.sock" --data-binary @- "http://localhost/event?event=<Name>"`, always `exit 0`.
- Never prints to stdout (stdout from some hooks is fed back to Claude).

`hooks/pet-hook.ps1`: same through `curl.exe --unix-socket` (`Invoke-RestMethod` has no unix socket support), wrapped in `try {} catch {}`.

`settings.snippet.json` (merged into `~/.claude/settings.json`, user-level so terminal and Desktop Code tab share it):

```json
{
  "hooks": {
    "SessionStart":     [{ "hooks": [{ "type": "command", "command": "~/.claude-pet/pet-hook.sh SessionStart" }] }],
    "UserPromptSubmit": [{ "hooks": [{ "type": "command", "command": "~/.claude-pet/pet-hook.sh UserPromptSubmit" }] }],
    "PreToolUse":       [{ "matcher": "*", "hooks": [{ "type": "command", "command": "~/.claude-pet/pet-hook.sh PreToolUse" }] }],
    "PostToolUse":      [{ "matcher": "*", "hooks": [{ "type": "command", "command": "~/.claude-pet/pet-hook.sh PostToolUse" }] }],
    "Notification":     [{ "hooks": [{ "type": "command", "command": "~/.claude-pet/pet-hook.sh Notification" }] }],
    "Stop":             [{ "hooks": [{ "type": "command", "command": "~/.claude-pet/pet-hook.sh Stop" }] }],
    "SubagentStop":     [{ "hooks": [{ "type": "command", "command": "~/.claude-pet/pet-hook.sh SubagentStop" }] }],
    "SessionEnd":       [{ "hooks": [{ "type": "command", "command": "~/.claude-pet/pet-hook.sh SessionEnd" }] }]
  }
}
```
Windows variant uses `powershell -NoProfile -ExecutionPolicy Bypass -File "%USERPROFILE%\\.claude-pet\\pet-hook.ps1" <Event>`.

Verify the current hook event names and payload fields against the Claude Code docs before implementing (they may have changed).

The settings app gets an **"Install Claude Code hooks"** button. It merges the snippet into `~/.claude/settings.json` without overwriting existing hooks, makes a backup first, and offers an uninstall.

---

## 5. Settings app (Flutter)

### Packages
- `fluent_ui` (Windows), `macos_ui` (macOS)
- `window_manager`: size, title bar, single instance focus
- `http`: talk to pet-core
- `path_provider`: config dir
- `hotkey_manager` (recorder widget only; registration stays in core)
- `file_picker`: import another pet (sheet + pet.json), later

### Architecture
- Shared layer (`lib/core`, `lib/pages`): `ConfigRepository` (read/write JSON atomically), `CoreClient` (`/state`, `/preview`, `/reload`), page controllers (`ChangeNotifier` or Riverpod; pick one and stick to it).
- Platform layer: `ui_fluent/` with `NavigationView` + `NavigationPane`, `ui_macos/` with `MacosWindow` + `Sidebar`. Same pages, different widgets.
- `main.dart`: `Platform.isMacOS ? MacosApp(...) : FluentApp(...)`.
- On launch: if core isn't reachable, show a banner "Pet is not running" + Start button.

### Pages
1. **General**: pet on/off, scale slider, always on top, click-through, start at login, reset position.
2. **Animations**: grid of all animations with live preview (render frames from the sheet in Flutter with `CustomPainter`), toggle each, FPS override, "Play on pet" button (`/preview`).
3. **Claude events**: table "event → animation" with dropdowns, sound on waiting, "Install / Remove Claude Code hooks", hook status (last event received, from `/state`).
4. **Behavior**: wander, follow cursor, idle timeout, do-not-disturb schedule (optional).
5. **Hotkeys**: recorder fields + conflict warnings.
6. **About**: version, links, open config folder, logs.

Save on change (debounced 300 ms) → write config → `POST /reload`.

---

## 6. Packaging

Both binaries ship together; core finds settings next to itself.

**Windows**
- `go build -ldflags "-H=windowsgui -s -w"` → `pet-core.exe` (no console window), embed icon + manifest via `go-winres`.
- `flutter build windows --release` → `settings/` folder.
- Inno Setup → `ClaudePet-Setup.exe` (installs to `%LOCALAPPDATA%\Programs\ClaudePet`, Start-menu shortcut, optional autostart, copies hook script to `~\.claude-pet\`).

**macOS**
- `ClaudePet.app` = core as the main executable (`LSUIElement=1`), Flutter settings app embedded at `Contents/Helpers/Settings.app`.
- Universal binary: `GOARCH=amd64` + `arm64` → `lipo`; `flutter build macos --release`.
- Ebiten on macOS needs cgo, so build on a Mac (or the macOS CI runner).
- `create-dmg` → `ClaudePet.dmg`. Unsigned builds need right-click → Open the first time; signing + notarization can come later with an Apple Developer account.

**CI**: GitHub Actions matrix (`windows-latest`, `macos-latest`) producing the installer and dmg as artifacts on tag.

---

## 7. Milestones

1. **M0: skeleton.** Repo layout, Go module, Flutter project, README, lint (`golangci-lint`, `flutter analyze`).
2. **M1: pet on screen.** Embedded sheet, `idle` loop, transparent undecorated floating window, drag, saved position. Check RSS.
3. **M2: events.** HTTP server on `<dataDir>/pet.sock` (no port, no token), stale-socket check and single instance, state machine with all animations, hook scripts on `curl --unix-socket`, `CoreClient` over the socket, manual test with `curl`, then a real Claude Code session. Check AF_UNIX on macOS and socket removal on exit on Windows.
4. **M3: tray + hotkeys + autostart.**
5. **M4: settings app.** Config repo, General + Animations pages on both UI kits, `/preview`, `/reload`.
6. **M5: Claude events page + hook installer.**
7. **M6: packaging.** Installer, dmg, CI.
8. **M7: polish.** Wander/follow-cursor, sounds, DND, multi-monitor, HiDPI, custom pet import.

---

## 8. Testing

- **Go unit tests:** state machine (event sequences → expected states, priorities, multiple sessions), config defaults/merge/atomic write, sprite slicing (frame rects from pet.json), server auth and async handling (`httptest`).
- **Flutter tests:** config repository round-trip, event-map editor logic, hook-settings merge (does not drop existing user hooks).
- **Manual checklist per OS:** transparency, no taskbar/Dock icon, drag, multi-monitor, HiDPI scaling, sleep/wake, hotkeys, autostart, idle RAM/CPU numbers.

---

## 9. Risks and open questions

- Ebiten transparency/floating behavior on each OS: prototype M1 on both machines first. If Windows tool-window style isn't reachable, add a small `x/sys/windows` shim.
- Ebiten and systray both want the main thread on macOS: confirm `RunWithExternalLoop` works with Ebiten's loop early in M3.
- Hook payload format can change between Claude Code versions: parse defensively, ignore unknown fields.
- Claude chat (non-Code) gives no events. Out of scope, so the pet stays idle there.
