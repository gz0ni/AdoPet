# Project Decisions

Recorded decisions and the reasoning behind them, so future work does not silently contradict them. Add an entry whenever a non-trivial choice is made. Latest first.

## Format

```
### <Decision title> (YYYY-MM-DD)

- Decision: <what was chosen, one sentence>
- Reason: <why this and not the alternatives>
- Alternatives considered: <what was rejected and why>
- Consequences: <what this decision affects>
```

---

<!-- New entries go here, newest first. -->

### Hook event name passed in the query string (2026-11-24)

- Decision: hook scripts POST the raw hook JSON as the body and the event name as `/event?event=<Name>`.
- Reason: scripts stay trivial — no JSON parsing or rewriting in sh/PowerShell.
- Alternatives considered: injecting `event` into the JSON body (needs jq or JSON handling in the hook).
- Consequences: the M2 `/event` handler must read `event` from the query; the payload's `hook_event_name` can serve as a fallback.

### Embed only spritesheet.png and pet.json (2026-11-24)

- Decision: `core/assets/embed.go` embeds just `spritesheet.png` and `pet.json`; the other asset folders stay as source material.
- Reason: keep the binary small; work folders are not needed at runtime.
- Alternatives considered: embedding the whole `assets/` directory.
- Consequences: replacing the default pet means updating these two files.

### Two processes: Go core + Flutter settings (2026-11-24)

- Decision: always-running pet in Go (Ebiten), settings UI in Flutter launched on demand.
- Reason: small idle footprint for the pet; native-looking settings on both OS.
- Alternatives considered: see [plan.md](plan.md).
- Consequences: they communicate only through config.json and the local HTTP API.
