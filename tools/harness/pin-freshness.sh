#!/usr/bin/env bash
# pin-freshness — read-only Sensor auf einen gepinnten Fremd-Bestandteil: meldet,
# wenn upstream einen NEUEREN Stand fuehrt als unser Pin (slice-137).
#
# ZWEI QUELLEN-FORMEN, weil eine davon eine Sonderquelle ist:
#   --github <owner/repo>  dem Redirect von .../releases/latest folgen und die
#                          effektive URL lesen; sie endet auf /releases/tag/<x>.
#                          Kein jq, keine API, kein Token (DC-QA-03-Sparsamkeit).
#   --godev                golang/go publiziert KEINE Release-Objekte -- der
#                          releases/latest-Pfad redirected ins Leere. Die stabile
#                          Version kommt als PLAINTEXT von go.dev/VERSION?m=text
#                          (erste Zeile, z. B. `go1.27.0`).
#
# NORMALISIERUNG: go.dev sagt `go1.27.0`, unser Pin ist bar `1.27.0`; GitHub sagt
# `v2.13.1`, unser Pin ebenso. Auf EIN Format bringen heisst hier: fuehrendes
# `go` strippen, sonst nichts. Ein `v` wird NICHT gestrippt -- die golangci-Achse
# fuehrt es auf beiden Seiten, und es zu entfernen machte den Vergleich
# grosszuegiger, nicht genauer.
#
# VERGLEICH IST GLEICH/UNGLEICH, kein Semver-Sort. Beide Reihen sind monoton;
# ein "neuer, aber aelter" existiert dort nicht. Waere eine Quelle es nicht,
# muesste das eigens entschieden werden -- nicht hier stillschweigend.
#
# FAIL-OPEN, und das ist die tragende Entscheidung: Netz-, Werkzeug- oder
# Parse-Ausfall meldet SKIP und Exit 0. Ein Sensor, der bei fremder Stoerung rot
# wird, wird abgeschaltet -- und ein abgeschalteter Waechter ist schlechter als
# ein loechriger. Zeitgrenzen gehoeren dazu: ohne sie waere eine HAENGENDE
# Verbindung kein SKIP, sondern ein Job-Timeout.
#
# NETZLOS PRUEFBAR: `--compare <name> <gepinnt> <upstream>` ruft NUR den
# Vergleicher. Ohne diesen Einstieg waere die Semantik nur mit Netz zu pruefen,
# und damit gar nicht.
#
# Exit: 0 = aktuell ODER SKIP, 3 = VERALTET. bash + coreutils + curl.
set -euo pipefail

CT=10   # connect-timeout
MT=60   # max-time

compare() {
  local name="$1" pinned="$2" upstream="$3"
  if [ -z "$upstream" ]; then
    echo "pin-freshness: ${name} SKIP — kein Upstream-Stand ermittelbar (Pin ${pinned})"
    return 0
  fi
  if [ "$pinned" = "$upstream" ]; then
    echo "pin-freshness: ${name} ok — Pin ${pinned} ist der neueste Stand"
    return 0
  fi
  echo "pin-freshness: ${name} VERALTET — Pin ${pinned}, upstream ${upstream}" >&2
  [ -n "${ADVICE:-}" ] && echo "pin-freshness: ${name} — ${ADVICE}" >&2
  return 3
}

# Reiner Vergleicher, ohne Netz: der Testeinstieg.
if [ "${1:-}" = "--compare" ]; then
  shift
  compare "${1:-}" "${2:-}" "${3:-}"
  exit $?
fi

command -v curl >/dev/null 2>&1 \
  || { echo "pin-freshness: 'curl' nicht gefunden — SKIP" >&2; exit 0; }

mode="${1:-}"; shift || true
name="${NAME:-?}"
pinned="${PINNED:-}"
[ -n "$pinned" ] || { echo "pin-freshness: PINNED ist leer — SKIP" >&2; exit 0; }

upstream=""
case "$mode" in
  --github)
    repo="${1:-}"
    eff="$(curl -fsSLo /dev/null -w '%{url_effective}' \
             --connect-timeout "$CT" --max-time "$MT" \
             "https://github.com/${repo}/releases/latest" 2>/dev/null || true)"
    case "$eff" in
      */releases/tag/*) upstream="${eff##*/releases/tag/}" ;;
      *)                upstream="" ;;   # kein Tag in der Endstation ⇒ SKIP
    esac
    ;;
  --godev)
    upstream="$(curl -fsSL --connect-timeout "$CT" --max-time "$MT" \
                  'https://go.dev/VERSION?m=text' 2>/dev/null \
                | head -1 | tr -d '\r' | sed 's/^go//' || true)"
    ;;
  *)
    echo "pin-freshness: Modus fehlt (--github <owner/repo> | --godev | --compare)" >&2
    exit 1
    ;;
esac

compare "$name" "$pinned" "$upstream"
