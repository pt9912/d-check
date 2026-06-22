#!/usr/bin/env bash
# completeness-check.sh — Requirements-Completeness-Gate (slice-042, ADR-0017).
#
# Closure-Invariante: jede DC-*-Anforderung muss von >=1 Slice referenziert
# sein. Quelle ist die RTM von `d-check --trace --json` (DC-FA-CLI-009):
# das top-level-Feld `orphans` (int). orphans>0 => FAIL (mit Waisen-IDs).
#
# `--trace` selbst bleibt advisory (Exit 0 bei Waisen, spec-fixiert); die
# Durchsetzung lebt hier. Bindepunkt: Closure (an `make fullbuild`), bewusst
# NICHT in `make gates`/`ci` — Greenfield erlaubt transiente Waisen, bis der
# umsetzende Slice landet.
#
# fail-closed: leere/feldlose/nicht-numerische/kaputte RTM => FAIL, nie
# stilles "0". Parsing rein bash/grep/awk (kein jq/python — keine
# Host-Binary-Abhaengigkeit). Negativ-Selbsttest in beide Richtungen bei
# jedem Lauf (analog tools/gate-consistency.sh, tools/adr-immutable-check.sh).
set -euo pipefail

cd "$(dirname "$0")/.."

IMAGE="${IMAGE:-d-check}"

# extract_orphans <json> -> druckt die top-level orphans-Ganzzahl, Exit 0;
# fehlendes/nicht-numerisches Feld => Exit 1 (fail-closed). Das top-level
# "orphans" (Plural) ist eindeutig von den per-Requirement "orphan"
# (Singular) unterscheidbar.
extract_orphans() {
  local val
  val="$(printf '%s' "$1" \
    | grep -oE '"orphans"[[:space:]]*:[[:space:]]*[0-9]+' \
    | head -n1 | grep -oE '[0-9]+$' || true)"
  [ -n "$val" ] || return 1
  printf '%s' "$val"
}

# decide <json> -> Exit 0 bei genau 0 Waisen, sonst Exit 1 (auch bei
# nicht-extrahierbarem orphans => fail-closed).
decide() {
  local n
  n="$(extract_orphans "$1")" || return 1
  [ "$n" -eq 0 ]
}

self_test() {
  local fail=0
  # Richtung A (muss PASS): 0 Waisen.
  decide '{"total":3,"orphans":0,"requirements":[]}' \
    || { echo "self-test: orphans:0 sollte PASS" >&2; fail=1; }
  # Richtung B (muss FAIL): >0 Waisen.
  if decide '{"total":3,"orphans":2,"requirements":[]}'; then
    echo "self-test: orphans:2 sollte FAIL" >&2; fail=1; fi
  # Stilles-Gruen-Vektoren — alle muessen FAIL (fail-closed):
  if decide '';                          then echo "self-test: leere RTM sollte FAIL" >&2; fail=1; fi
  if decide '{"total":3,"requirements":[]}'; then echo "self-test: fehlendes orphans-Feld sollte FAIL" >&2; fail=1; fi
  if decide '{"orphans":"x"}';           then echo "self-test: nicht-numerisches orphans sollte FAIL" >&2; fail=1; fi
  if decide '{kaputt';                   then echo "self-test: kaputte RTM sollte FAIL" >&2; fail=1; fi
  if [ "$fail" -ne 0 ]; then
    echo "completeness-check: Selbsttest FEHLGESCHLAGEN" >&2
    exit 2
  fi
}

self_test

# Echter Lauf: RTM frisch aus dem read-only/netzlosen --trace --json.
json="$(docker run --rm --network none -v "$(pwd)":/repo:ro "$IMAGE":latest --trace --json 2>/dev/null || true)"

if [ -z "$json" ]; then
  echo "completeness-check: FAIL — leere RTM-Ausgabe (Image fehlt/Lauf-Fehler); fail-closed" >&2
  exit 1
fi

orphans="$(extract_orphans "$json")" || {
  echo "completeness-check: FAIL — top-level 'orphans' nicht als Ganzzahl extrahierbar (RTM-Schema-Drift?); fail-closed" >&2
  exit 1
}

if [ "$orphans" -gt 0 ]; then
  echo "completeness-check: FAIL — $orphans Requirements-Waise(n) (Anforderung ohne referenzierenden Slice):" >&2
  # Waisen-IDs: id steht je Requirement vor dem orphan-Flag; -F'"' ist
  # portabel (gawk/mawk), kein gensub/3-arg-match noetig.
  printf '%s\n' "$json" | awk -F'"' '
    /"id"[[:space:]]*:/ { id=$4 }
    /"orphan"[[:space:]]*:[[:space:]]*true/ { print "  - " id }
  ' >&2
  echo "  (Details: make trace)" >&2
  exit 1
fi

total="$(printf '%s' "$json" | grep -oE '"total"[[:space:]]*:[[:space:]]*[0-9]+' | head -n1 | grep -oE '[0-9]+$' || true)"
echo "completeness-check: OK — 0 Requirements-Waisen, ${total:-?} Anforderung(en) von >=1 Slice abgedeckt (Selbsttest gefeuert)."
