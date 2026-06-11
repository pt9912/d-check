#!/usr/bin/env bash
# coverage-gate.sh — Go-Coverage-Gate (Kalibrierungs-Bindung: aktuelle
# Schwelle und Historie in harness/README.md §Sensors; Verfehlung ⇒
# Carveout-Pflicht). Muster: u-boot scripts/coverage-gate.sh (gleiche
# Build-Familie).
#
# Aufruf:
#   coverage-gate.sh <coverage-func.txt> <threshold>
#
# Liest die Ausgabe von `go tool cover -func=<profile>` und erzwingt
# die Gesamt-Coverage gegen die Schwelle (Prozent, ganz- oder
# gleitkommazahlig). Läuft als Dockerfile-Stage (make coverage-gate).
set -euo pipefail

if [[ $# -ne 2 ]]; then
  echo "usage: $0 <coverage-func.txt> <threshold>" >&2
  exit 2
fi

func_file="$1"
threshold="$2"

if [[ ! -s "$func_file" ]]; then
  echo "coverage-gate: Coverage-Eingabe fehlt oder ist leer: $func_file" >&2
  echo "Hinweis: ist 'go test -coverprofile' vorher fehlgeschlagen?" >&2
  exit 2
fi

# Letzte Zeile von `go tool cover -func`: "total:\t(statements)\tXX.X%"
total_line="$(grep -E '^total:' "$func_file" || true)"
if [[ -z "$total_line" ]]; then
  echo "coverage-gate: keine 'total:'-Zeile in $func_file" >&2
  exit 2
fi

total_pct="$(echo "$total_line" | grep -oE '[0-9]+\.[0-9]+%?$' | tr -d '%')"
if [[ -z "$total_pct" ]]; then
  echo "coverage-gate: Coverage-Prozent nicht parsbar: $total_line" >&2
  exit 2
fi

pass="$(awk -v p="$total_pct" -v t="$threshold" 'BEGIN { print (p+0 >= t+0) ? 1 : 0 }')"
if [[ "$pass" != "1" ]]; then
  printf "coverage-gate: FAIL — Coverage %.2f%% unter Schwelle %s%%\n" "$total_pct" "$threshold" >&2
  exit 1
fi

printf "coverage-gate: OK — Coverage %.2f%% erfüllt Schwelle %s%%\n" "$total_pct" "$threshold"
