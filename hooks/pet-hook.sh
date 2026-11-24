#!/bin/sh
# Forwards a Claude Code hook payload to pet-core. Never prints, never fails.
EVENT="$1"
case "$(uname)" in
  Darwin) DIR="$HOME/Library/Application Support/ClaudePet" ;;
  *)      DIR="${XDG_CONFIG_HOME:-$HOME/.config}/ClaudePet" ;;
esac
PORT=$(cat "$DIR/port" 2>/dev/null) || exit 0
TOKEN=$(cat "$DIR/token" 2>/dev/null) || exit 0
curl -s -m 1 -X POST "http://127.0.0.1:$PORT/event?event=$EVENT" \
  -H "X-Pet-Token: $TOKEN" -H "Content-Type: application/json" \
  --data-binary @- >/dev/null 2>&1
exit 0
