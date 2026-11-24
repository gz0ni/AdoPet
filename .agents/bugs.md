# Bugs

Known bugs found in this project and how they were fixed, so they do not reappear. Logged at the end of every session. Status: `open` / `fixed` / `won't fix`.

## Format

```
### <Short bug title> (YYYY-MM-DD) — status

- Symptom
- Root cause
- Fix applied
- How to prevent recurrence
```

---

<!-- New entries go here, newest first. -->

### pet-core uses ~130 MB on macOS, target 40 MB (2026-11-24) — fixed

- Symptom: on Apple Silicon RSS ≈ 150 MB, physical footprint ≈ 130 MB (peak 179 MB) for a 96×104 pt window.
- Root cause: Ebiten Metal backend; `vmmap` shows ~38 MB `IOAccelerator (graphics)` plus Metal/AppKit malloc.
- Fix applied: own Cocoa `NSPanel` + `CALayer` backend; footprint ≈ 32 MB (peak 49 MB at startup).
- How to prevent recurrence: measure footprint with `vmmap --summary` after any macOS rendering change.

### pet-core shows a Dock icon on macOS (2026-11-24) — fixed

- Symptom: Dock icon appears; `lsappinfo` reports `type="Foreground"`.
- Root cause: Ebiten `RunGameOptions.SkipTaskbar` only works on Windows; activation policy stays Regular.
- Fix applied: `petInit` sets `NSApplicationActivationPolicyAccessory`; `lsappinfo` reports `UIElement`.
- How to prevent recurrence: set `NSApplicationActivationPolicyAccessory` (or `LSUIElement` in Info.plist for the bundle).

### Ebiten build saved a minimized window position (2026-11-24) — fixed

- Symptom: `config.json` had `position.y = -32817`; next start put the pet in the top-right corner.
- Root cause: on Windows a minimized window sits near (-32000, -32000); the Ebiten drag-end handler saved `WindowPosition()` without checking that.
- Fix applied: Windows no longer uses Ebiten; the Win32 backend saves only on `WM_EXITSIZEMOVE`. `place` clamps any loaded position on-screen.
- How to prevent recurrence: never persist a position read while the window is minimized or hidden.

### pet-core uses ~108 MB on Windows, target 40 MB (2026-11-24) — fixed

- Symptom: idle WorkingSet ≈ 108 MB, private ≈ 140 MB on RTX 4060; Go heap is only ~6 MB (`gctrace`).
- Root cause: Ebiten transparent windows require OpenGL on Windows; the NVIDIA GL driver dominates the footprint.
- Fix applied: own Win32 layered window with a CPU-filled DIB; WS ≈ 18 MB, private ≈ 21 MB.
- How to prevent recurrence: measure WS after any rendering change; do not bring a GPU API into the Windows window.
