#!/usr/bin/env bash
# working-tree-hash — inhaltsbasierter Hash über alle getrackten und
# untracked Dateien (exclude-standard). Gemeinsame Quelle der Wahrheit
# für `record-gates` und den Stop-Hook.
#
# INHALTSBASIERT statt diff-basiert (Abweichung vom b-cad-Vorbild,
# harness/conventions.md MR-005): Ein Commit ohne Inhaltsänderung
# lässt den Hash unverändert — der Gate-Nachweis bleibt gültig. Ein
# Commit OHNE vorherigen Gate-Lauf erzeugt dagegen keinen passenden
# Nachweis und wird vom Stop-Hook nicht mehr durchgewunken.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

git ls-files -z --cached --others --exclude-standard \
  | sort -zu \
  | while IFS= read -r -d '' f; do
      if [ -L "$f" ]; then
        printf 'LINK %s -> %s\n' "$f" "$(readlink -- "$f")"
      elif [ -f "$f" ]; then
        sha256sum -- "$f"
      else
        printf 'GONE %s\n' "$f"
      fi
    done | sha256sum | awk '{print $1}'
