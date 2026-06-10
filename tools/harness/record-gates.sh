#!/usr/bin/env bash
# record-gates — Nachweis schreiben, dass `make gates` den aktuellen
# Arbeitsbaum-Zustand abgedeckt hat. Läuft als letzter gates-Prerequisite
# (nur bei grünen Gates). Der Stop-Hook vergleicht denselben Hash.
# Übernommen aus b-cad (harness/conventions.md MR-004).
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

mkdir -p .harness/state
bash tools/harness/working-tree-hash.sh > .harness/state/gates-passed.diffsha
