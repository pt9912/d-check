#!/usr/bin/env bash
# working-tree-hash — ein Hash über den GESAMTEN Arbeitsbaum-Zustand.
# Gemeinsame Quelle der Wahrheit für `record-gates` und den Stop-Hook,
# damit die Logik nicht zwischen Makefile und Hook driftet.
# Übernommen aus b-cad (harness/conventions.md MR-004).
#
# Erfasst: tracked-modified, staged, untracked, gelöschte, umbenannte
# Dateien (anders als `git diff`, das untracked NICHT sieht).
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

{
  git status --porcelain=v1
  git diff --binary
  git diff --cached --binary

  git ls-files --others --exclude-standard -z \
    | sort -z \
    | while IFS= read -r -d '' file; do
        printf 'UNTRACKED %s\n' "$file"
        sha256sum "$file"
      done
} | sha256sum | awk '{print $1}'
