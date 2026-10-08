#!/usr/bin/env bash
# coverage-gate.sh — Go-Coverage-Gate. Schwelle, Messbasis und Randformen
# legt spec/spezifikation.md §7 (SPEC-089) fest; Verfehlung ⇒
# Carveout-Pflicht, Senkung nur per ADR. Muster: u-boot
# scripts/coverage-gate.sh (gleiche Build-Familie).
#
# Aufruf:
#   coverage-gate.sh <coverage-func.txt> <threshold>
#
# Liest die Ausgabe von `go tool cover -func=<profile>` und erzwingt
# die Gesamt-Coverage gegen die Schwelle (Prozent, ganz- oder
# gleitkommazahlig). Läuft als Dockerfile-Stage (make coverage-gate).
set -euo pipefail
# Zahlen mit Punkt: unter einer Locale mit Dezimalkomma liest printf "93.0"
# nicht als Zahl und endet mit 1 — auch im bestandenen Fall.
export LC_ALL=C

if [[ $# -ne 2 ]]; then
  echo "usage: $0 <coverage-func.txt> <threshold>" >&2
  exit 2
fi

func_file="$1"
threshold="$2"

# Die Schwelle ist ein Build-Parameter. Leer, nicht numerisch oder negativ
# verglich awk sie bisher als 0 — jede Coverage bestand still. Nur eine
# nicht negative Zahl ist eine Schwelle; alles andere ist gescheitert.
if [[ ! "$threshold" =~ ^[0-9]+(\.[0-9]+)?$ ]]; then
  echo "coverage-gate: Schwelle '$threshold' ist keine nicht negative Zahl" >&2
  exit 2
fi

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

# `|| true`: ohne Treffer bricht grep die Zuweisung unter `set -e pipefail`
# sonst mit Exit 1 ab, und der Exit-2-Zweig darunter wird nie erreicht.
total_pct="$(echo "$total_line" | grep -oE '[0-9]+\.[0-9]+%?$' | tr -d '%' || true)"
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
