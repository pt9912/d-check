#!/usr/bin/env bash
# stop-require-gates — gibt den Stop nur frei, wenn der aktuelle Repo-INHALT durch
# einen erfolgreichen `make gates`-Lauf abgedeckt ist. Nutzt dieselbe
# inhaltsbasierte Hash-Funktion wie record-gates (keine Logik-Dopplung;
# harness/conventions.md MR-004/MR-005).
#
# STRENG (Default, ohne jeden Schalter). Jedes Turn-Ende mit ungedecktem Inhalt blockt,
# auch ohne neuen HEAD. Ein fehlender, gelöschter oder unlesbarer Schalter lockert
# nichts: gelockert wird nur durch eine Datei, die da ist.
#
# LOCKER (nur mit `.harness/stop-gate-locker`). Der Hook blockt dann nur bei neuem HEAD
# UND abweichendem Hash: ein Turn-Ende, dessen HEAD dem von record-gates gestempelten
# (.harness/state/gates-passed.head) gleicht, geht auch mit ungedeckter Änderung
# frei. Das ist eine Lockerung des Handoff-Gates und braucht eine ADR
# (AGENTS.md §3.6); dieses Repo legt die Datei nicht an (MR-079). STOP_GATE_STRENG
# genau `1` hebt die Datei auf; jeder andere Wert der Variablen sagt nichts. Der Hook
# liest die Datei, er schreibt und löscht sie nie.
#
# Da der Hash inhaltsbasiert ist, gilt der Nachweis über Commits hinweg (gleicher
# Inhalt = gleicher Hash) — ein Commit ohne Gate-Lauf macht den Stop dagegen NICHT
# mehr grün. Restlücke: ein frischer Klon bzw. gelöschter .harness-State mit
# cleanem Tree wird freigegeben (kein Nachweis prüfbar); CI ist dort das Netz.
#
# FAIL-CLOSED. Was der Hook nicht lesen kann (Hash, Nachweis, `git status`, HEAD),
# beantwortet er mit `decision: block` und Grund (MR-076). Jeder darüber hinaus
# unerwartete Fehler endet mit Exit 2 (blockierend), nie mit einem anderen Code
# (Trap unten) — Exit 1 ließe den Stop durch. head_wert ermittelt den HEAD-Wert
# genauso wie tools/harness/record-gates.sh — Kopplung: beide Fassungen gleich halten.
set -euo pipefail
rc=0
trap 'rc=$?; [ "$rc" -eq 0 ] || exit 2' EXIT
cd "$(git rev-parse --show-toplevel)"

state_file=".harness/state/gates-passed.diffsha"
head_file=".harness/state/gates-passed.head"

approve() {
  printf '%s\n' '{"decision":"approve"}'
  exit 0
}

block() {
  printf '{\n  "decision": "block",\n  "reason": "%s"\n}\n' "$1"
  exit 0
}

# head_wert — die aufgelöste SHA von HEAD; `kein-commit` nur auf einem ungeborenen
# Zweig in einem Repo ohne einzigen Commit; jede andere Lage ist ein Git-Fehler.
head_wert() {
  local sha alle
  if sha="$(git rev-parse --verify -q HEAD)"; then
    printf '%s\n' "$sha"
    return 0
  fi
  git symbolic-ref -q HEAD >/dev/null || return 1
  alle="$(git rev-list -n1 --all)" || return 1
  [ -z "$alle" ] || return 1
  printf '%s\n' kein-commit
}

# Schleifen-Schutz: Hat dieser Hook den Stop bereits einmal blockiert
# (stop_hook_active), nicht erneut blockieren — sonst Endlosschleife
# bei dauerhaft rotem Gate.
input="$(cat || true)"
if printf '%s' "$input" | grep -q '"stop_hook_active"[[:space:]]*:[[:space:]]*true'; then
  approve
fi

locker=0
if [ -e .harness/stop-gate-locker ]; then locker=1; fi
if [ "${STOP_GATE_STRENG:-}" = 1 ]; then locker=0; fi

if [ "$locker" -eq 1 ] && [ -e "$head_file" ]; then
  if ! current_head="$(head_wert)" || ! recorded_head="$(cat "$head_file")"; then
    block "HEAD or the recorded HEAD stamp could not be read. Fix the cause, then run \`make gates\`."
  fi
  if [ "$current_head" = "$recorded_head" ]; then approve; fi
fi

if [ ! -f "$state_file" ]; then
  # Ein `git status`, das scheitert (kein Repository, kaputter Index), liest
  # den Zustand nicht — leer wäre dann kein sauberer Baum.
  if ! status="$(git status --porcelain=v1)"; then
    block "The working tree state could not be read (git status failed). Fix the cause, then run \`make gates\`."
  fi
  # Frischer Klon ohne lokale Änderungen: kein Nachweis prüfbar.
  if [ -z "$status" ]; then approve; fi
  block "There are working tree changes, but no recorded successful make gates run. Run \`make gates\`."
fi

# Ein Hash, der nicht entsteht, ist kein Nachweis: der Hook blockt mit Grund,
# statt unter `set -e` ohne Antwort zu enden — das ließe den Stop durch.
if ! current="$(bash tools/harness/working-tree-hash.sh)" \
    || ! recorded="$(cat "$state_file")"; then
  block "The content hash of the working tree could not be computed or the recorded one could not be read. Fix the cause, then run \`make gates\`."
fi

if [ "$current" != "$recorded" ]; then
  block "Repo content changed after the last recorded gates run (commits without a gates run count too). Run \`make gates\` again."
fi

approve
