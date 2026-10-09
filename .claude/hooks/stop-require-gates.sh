#!/usr/bin/env bash
# stop-require-gates — gibt den Stop nur frei, wenn der aktuelle
# Repo-INHALT durch einen erfolgreichen `make gates`-Lauf abgedeckt ist.
# Nutzt dieselbe inhaltsbasierte Hash-Funktion wie record-gates (keine
# Logik-Dopplung; harness/conventions.md MR-004/MR-005).
#
# Da der Hash inhaltsbasiert ist, gilt der Nachweis über Commits hinweg
# (gleicher Inhalt = gleicher Hash) — ein Commit ohne Gate-Lauf macht
# den Stop dagegen NICHT mehr grün. Restlücke: frischer Klon bzw.
# gelöschter .harness-State mit cleanem Tree wird freigegeben
# (kein Nachweis prüfbar); CI ist dort das Netz.
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

if [ ! -f "$state_file" ]; then
  # Ein `git status`, das scheitert (kein Repository, kaputter Index), liest
  # den Zustand nicht — leer wäre dann kein sauberer Baum.
  if ! status="$(git status --porcelain=v1)"; then
    cat <<'JSON'
{
  "decision": "block",
  "reason": "The working tree state could not be read (git status failed). Fix the cause, then run `make gates`."
}
JSON
    exit 0
  fi
  if [ -z "$status" ]; then
    # Frischer Klon ohne lokale Änderungen: kein Nachweis prüfbar.
    cat <<'JSON'
{"decision":"approve"}
JSON
    exit 0
  fi
  cat <<'JSON'
{
  "decision": "block",
  "reason": "There are working tree changes, but no recorded successful make gates run. Run `make gates`."
}
JSON
  exit 0
fi

# Ein Hash, der nicht entsteht, ist kein Nachweis: der Hook blockt mit Grund,
# statt unter `set -e` ohne Antwort zu enden — das ließe den Stop durch.
if ! current="$(bash tools/harness/working-tree-hash.sh)" \
    || ! recorded="$(cat "$state_file")"; then
  cat <<'JSON'
{
  "decision": "block",
  "reason": "The content hash of the working tree could not be computed or the recorded one could not be read. Fix the cause, then run `make gates`."
}
JSON
  exit 0
fi

if [ "$current" != "$recorded" ]; then
  cat <<'JSON'
{
  "decision": "block",
  "reason": "Repo content changed after the last recorded gates run (commits without a gates run count too). Run `make gates` again."
}
JSON
  exit 0
fi

cat <<'JSON'
{"decision":"approve"}
JSON
