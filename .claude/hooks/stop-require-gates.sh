#!/usr/bin/env bash
# stop-require-gates — gibt den Stop nur frei, wenn der aktuelle
# Arbeitsbaum durch einen erfolgreichen `make gates`-Lauf abgedeckt ist.
# Nutzt dieselbe Hash-Funktion wie record-gates (keine Logik-Dopplung).
# Übernommen aus b-cad (harness/conventions.md MR-004).
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

state_file=".harness/state/gates-passed.diffsha"

# Schleifen-Schutz: Hat dieser Hook den Stop bereits einmal blockiert
# (stop_hook_active), nicht erneut blockieren — sonst Endlosschleife
# bei dauerhaft rotem Gate.
input="$(cat || true)"
if printf '%s' "$input" | grep -q '"stop_hook_active"[[:space:]]*:[[:space:]]*true'; then
  cat <<'JSON'
{"decision":"approve"}
JSON
  exit 0
fi

# Kein Diff und keine untracked files → nichts zu blockieren.
# (Bewusste Lücke: ein Commit ohne vorherigen Gate-Lauf macht den Tree
# clean — Konvention "gates vor Commit" + CI fangen das ab; Review R1.)
if [ -z "$(git status --porcelain=v1)" ]; then
  cat <<'JSON'
{"decision":"approve"}
JSON
  exit 0
fi

if [ ! -f "$state_file" ]; then
  cat <<'JSON'
{
  "decision": "block",
  "reason": "There are working tree changes, but no recorded successful make gates run. Run `make gates`."
}
JSON
  exit 0
fi

current="$(bash tools/harness/working-tree-hash.sh)"
recorded="$(cat "$state_file")"

if [ "$current" != "$recorded" ]; then
  cat <<'JSON'
{
  "decision": "block",
  "reason": "The working tree changed after the last recorded gates run. Run `make gates` again."
}
JSON
  exit 0
fi

cat <<'JSON'
{"decision":"approve"}
JSON
