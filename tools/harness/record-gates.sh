#!/usr/bin/env bash
# record-gates — Nachweis schreiben, dass `make gates` den aktuellen
# Arbeitsbaum-Zustand abgedeckt hat. `make gates` ruft es im eigenen Rezept,
# nach allen grünen Gliedern (MR-076). Der Stop-Hook vergleicht denselben Hash.
# Übernommen aus b-cad (harness/conventions.md MR-004).
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

mkdir -p .harness/state
bash tools/harness/working-tree-hash.sh > .harness/state/gates-passed.diffsha
