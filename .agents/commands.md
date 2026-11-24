# Commands

Verified on Windows 11 on 2026-11-24 unless marked otherwise. `go vet/test/build` and `go run ./cmd/pet-core` also verified on macOS 27 (Apple Silicon) on 2026-11-24.

## Development

```bash
cd core && go run ./cmd/pet-core
cd settings && flutter pub get
cd settings && flutter run -d windows    # macOS: -d macos
```

## Verification

Run before claiming work is done:

```bash
cd core && go vet ./... && go test ./... && go build ./...
cd core && golangci-lint run              # v2.15.0, in $(go env GOPATH)/bin on both machines
cd settings && flutter analyze && flutter test
```

## Build / Package

```bash
cd settings && flutter build windows --debug     # verified
cd settings && flutter build windows --release
cd core && go build -ldflags "-H=windowsgui -s -w" -o pet-core.exe ./cmd/pet-core   # verified, ~5.3 MB
```

Packaging (Inno Setup, dmg, `lipo` universal binary) is not set up yet, see [plan.md](plan.md) §6.

## Measuring pet-core

```powershell
$p = Start-Process .\pet-core.exe -PassThru; Start-Sleep 5; $p.Refresh(); $p.WorkingSet64/1MB
```

```bash
ps -o rss=,%cpu= -p $(pgrep -x pet-core); vmmap --summary $(pgrep -x pet-core) | grep "Physical footprint"   # macOS
```

`GODEBUG=gctrace=1` shows the Go heap. On macOS judge memory by physical footprint, not RSS. Cross-check the Windows backend from a Mac: `GOOS=windows GOARCH=amd64 go vet ./...`. The macOS backend uses cgo and cannot be cross-checked from Windows.

## CI

`.github/workflows/build.yml`: on push to `main`, tags `v*`, and PRs. Matrix `windows-latest`, `macos-latest`: Go vet/test/build in `core/`, then `flutter pub get`, `flutter analyze`, `flutter test` in `settings/`. No artifacts yet.

## Tests

- Go: `internal/config` (defaults, merge, unknown keys, atomic save), `internal/sprite` (pet.json, frame rects, frame timing, scaling), `internal/window` (anchor placement, BGRA/RGBA blit, Win32 screen ↔ DIP and Cocoa ↔ config conversion). Planned: state machine, server auth (`httptest`).
- Flutter: `settings/test/pages_test.dart` (page titles). Planned: config repo round-trip, event-map editor, hook-settings merge.
