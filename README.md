# Claude Pet

Desktop pet (Chibi Ado) that reacts to Claude Code activity via hooks.

- `core/`: **pet-core** (Go, native Win32 / Cocoa window). Pet window, tray, hotkeys, local HTTP API.
- `settings/`: **settings** (Flutter, `fluent_ui` on Windows / `macos_ui` on macOS). Launched on demand.
- `hooks/`: Claude Code hook scripts and `settings.json` snippets.
- `packaging/`, `.github/workflows/`: installers and CI.

## Dev

```sh
cd core && go run ./cmd/pet-core
cd settings && flutter run -d windows   # or macos
```

Lint: `golangci-lint run` (core), `flutter analyze` (settings).

See the implementation plan for milestones (M0–M7).
