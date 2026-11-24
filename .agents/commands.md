# Commands

Verified on Windows 11 on 2026-11-24 unless marked otherwise.

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
cd core && golangci-lint run              # not installed locally yet
cd settings && flutter analyze && flutter test
```

## Build / Package

```bash
cd settings && flutter build windows --debug     # verified
cd settings && flutter build windows --release
cd core && go build -ldflags "-H=windowsgui -s -w" -o pet-core.exe ./cmd/pet-core   # not verified
```

Packaging (Inno Setup, dmg, `lipo` universal binary) is not set up yet, see [plan.md](plan.md) §6.

## CI

`.github/workflows/build.yml`: on push to `main`, tags `v*`, and PRs. Matrix `windows-latest`, `macos-latest`: Go vet/test/build in `core/`, then `flutter pub get`, `flutter analyze`, `flutter test` in `settings/`. No artifacts yet.

## Tests

- Go: none yet. Planned: state machine, config, sprite slicing, server auth (`httptest`).
- Flutter: `settings/test/pages_test.dart` (page titles). Planned: config repo round-trip, event-map editor, hook-settings merge.
