#!/usr/bin/env bash
# record-gates — Nachweis schreiben, dass `make gates` den aktuellen
# Arbeitsbaum-Zustand abgedeckt hat. `make gates` ruft es im eigenen Rezept,
# nach allen grünen Gliedern (MR-076). Der Stop-Hook vergleicht denselben Hash.
# Übernommen aus b-cad (harness/conventions.md MR-004).
#
# HEAD-STEMPEL: neben dem Hash schreibt dieses Skript die aufgelöste Commit-SHA von
# HEAD nach .harness/state/gates-passed.head. Den Stempel liest der Stop-Hook nur im
# lockeren Modus (nur mit .harness/stop-gate-locker, MR-079); im strengen Modus dieses
# Repos zählt allein der Hash. head_wert ermittelt den Wert genauso wie
# .claude/hooks/stop-require-gates.sh — Kopplung: beide Fassungen gleich halten. Ist
# HEAD nicht auflösbar und das Repo nicht eng als commitlos erkannt, endet das Skript
# rot, bevor es einen Stempel schreibt.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

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

if ! head="$(head_wert)"; then
  echo "record-gates: FEHLER — HEAD ist nicht aufloesbar, und das Repo ist nicht als commitlos erkannt (ungeborener Zweig ohne einzigen Commit) — kein Stempel geschrieben." >&2
  exit 1
fi

mkdir -p .harness/state
bash tools/harness/working-tree-hash.sh > .harness/state/gates-passed.diffsha
printf '%s\n' "$head" > .harness/state/gates-passed.head
