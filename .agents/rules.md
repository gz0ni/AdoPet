# Rules

## Code Style

- Go: `gofmt`, `go vet`, golangci-lint v2 with `linters.default: standard` (`core/.golangci.yml`). Packages live under `core/internal/<area>`.
- Dart: `flutter_lints` via `settings/analysis_options.yaml`; `flutter analyze` must report no issues.
- Shared settings logic goes in `settings/lib/core` and `settings/lib/pages`; widgets only in `ui_fluent/` or `ui_macos/`.

## Languages

- User-facing UI text: Russian.
- Identifiers, file names, commits, docs: English.

## Design

- Follow [plan.md](plan.md); do not invent new architecture without recording it in `project-decisions.md`.
- Windows UI uses only `fluent_ui` widgets, macOS UI only `macos_ui` widgets; same pages on both.
- pet-core must stay small: idle RSS ≤ 40 MB, CPU ≈ 0% idle, < 2% animating.

## State and Persistence

- Config: `config.json` in `%APPDATA%\ClaudePet` / `~/Library/Application Support/ClaudePet` (`core/internal/paths`). Atomic writes, unknown keys preserved, missing keys filled from defaults.
- `port` and `token` files in the same data dir, read by hooks and settings.
- Flutter state: ChangeNotifier or Riverpod — not chosen yet; pick one and record it.

## Testing

- Go: standard `testing` + `httptest`.
- Flutter: `flutter_test`.
- Hook scripts must never print to stdout and must always exit 0.

## Generated and External Files

- Do not edit `settings/windows/` and `settings/macos/` runner code unless a platform feature requires it.
- `core/assets/` work folders (`decoded/`, `frames/`, `final/`, `prompts/`, `qa/`, …) are sprite source material; only `spritesheet.png` and `pet.json` are embedded.
- `settings/build/` and `.dart_tool/` are build output.

## Git

- Do not commit, push, create PRs, or run destructive git operations without explicit user request.
- Before any commit (if requested): check `git status`, `git diff`, `git log --oneline -10`; stage only intended files.
