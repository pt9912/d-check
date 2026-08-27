#!/usr/bin/env bash
# pin-freshness — read-only Sensor auf einen gepinnten Fremd-Bestandteil: meldet,
# wenn upstream einen NEUEREN Stand fuehrt als unser Pin.
#
# DREI QUELLEN-FORMEN. Die ersten zwei fragen „gibt es einen neueren Tag" und
# unterscheiden sich nur in der Sonderquelle; die dritte fragt etwas anderes —
# „traegt derselbe Tag einen anderen Digest":
#   --github <owner/repo>  dem Redirect von .../releases/latest folgen und die
#                          effektive URL lesen; sie endet auf /releases/tag/<x>.
#                          Kein jq, keine API, kein Token (DC-QA-03-Sparsamkeit).
#   --godev                Fuer golang/go endet releases/latest auf
#                          .../releases -- OHNE /tag/. Der GitHub-Zweig SKIPpte
#                          hier also, statt falsch zu vergleichen. Die stabile
#                          Version kommt als PLAINTEXT von go.dev/VERSION?m=text
#                          (erste Zeile, z. B. `go1.27.0`) und wird auf ihre
#                          Form geprueft, bevor sie als Stand gilt.
#   --digest <ref>         Der Digest, den <ref> heute traegt, via
#                          `docker buildx imagetools inspect`. Fuer einen Tag
#                          OHNE Version die einzige Handhabe.
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
# Exit: 0 = aktuell ODER SKIP, 3 = VERALTET. bash + coreutils; je Zweig curl
# bzw. docker — der Werkzeug-Riegel steht deshalb IM Zweig, nicht davor.
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

# Der Werkzeug-Riegel steht JE ZWEIG, nicht davor: die Digest-Form braucht kein
# curl, und ein SKIP wegen eines Werkzeugs, das der Zweig gar nicht ruft, wäre
# ein stilles Abschalten aus dem falschen Grund.
brauche_curl() {
  command -v curl >/dev/null 2>&1 \
    || { echo "pin-freshness: 'curl' nicht gefunden — SKIP" >&2; exit 0; }
}

mode="${1:-}"; shift || true
name="${NAME:-?}"
pinned="${PINNED:-}"
[ -n "$pinned" ] || { echo "pin-freshness: PINNED ist leer — SKIP" >&2; exit 0; }

upstream=""
case "$mode" in
  --github)
    brauche_curl
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
    brauche_curl
    # Die Antwort wird auf ihre FORM geprueft, bevor sie als Stand gilt. Ohne
    # das waere eine HTTP-200-Nicht-Versions-Antwort (Fehlerseite, Wartungstext)
    # ein gueltiger "upstream" -- und ergaebe ein falsches VERALTET statt des
    # zugesagten SKIP. Ein Parse-Ausfall ist ein Ausfall, kein Befund.
    raw="$(curl -fsSL --connect-timeout "$CT" --max-time "$MT" \
             'https://go.dev/VERSION?m=text' 2>/dev/null \
           | head -1 | tr -d '\r' || true)"
    raw="${raw#go}"
    if printf '%s' "$raw" | grep -qE '^[0-9]+\.[0-9]+(\.[0-9]+)?([a-z]+[0-9]+)?$'; then
      upstream="$raw"
    else
      upstream=""
    fi
    # BEIDE Seiten normalisieren. Traegt der Pin eines Tages `go1.27.0`,
    # verglichen wir sonst dauerhaft ungleich und meldeten fuer immer VERALTET.
    pinned="${pinned#go}"
    ;;
  --digest)
    # DRITTE Quellen-Form, und sie beantwortet eine ANDERE Frage als die zwei
    # oben: nicht „gibt es einen neueren Tag", sondern „traegt derselbe Tag
    # inzwischen einen anderen Digest". Fuer ein Basis-Image ohne Version im
    # Tag — `distroless/static-debian12:nonroot` — ist das die einzige
    # Handhabe: eine Tag-Frische-Achse gibt es dort nicht, weil es keinen Tag
    # gibt, der sich bewegt.
    #
    # Die Quelle ist `docker buildx imagetools inspect`, nicht curl: eine
    # Registry-Abfrage von Hand braeuchte Token-Austausch und
    # Manifest-Listen-Aufloesung — beides ein Parser durch die Hintertuer.
    # Docker steht ohnehin in der Host-Klasse (AGENTS.md 3.1).
    #
    # GRENZE: verglichen wird der Digest der ADRESSIERTEN Referenz. Bei einer
    # Multi-Plattform-Liste ist das der Listen-Digest — dieselbe Groesse, die
    # im Dockerfile steht. Ein Plattform-Digest darunter ist eine andere Zahl
    # und hier nicht gemeint.
    ref="${1:-}"
    command -v docker >/dev/null 2>&1 \
      || { echo "pin-freshness: 'docker' nicht gefunden — SKIP" >&2; exit 0; }
    raw="$(docker buildx imagetools inspect "$ref" --format '{{.Manifest.Digest}}' 2>/dev/null || true)"
    # Form pruefen, bevor der Wert als Stand gilt — dieselbe Begruendung wie
    # beim godev-Zweig: eine Nicht-Digest-Antwort waere sonst ein falsches
    # VERALTET statt des zugesagten SKIP.
    if printf '%s' "$raw" | grep -qE '^sha256:[0-9a-f]{64}$'; then
      upstream="$raw"
    else
      upstream=""
    fi
    ;;
  *)
    echo "pin-freshness: Modus fehlt (--github <owner/repo> | --godev | --digest <ref> | --compare)" >&2
    exit 1
    ;;
esac

compare "$name" "$pinned" "$upstream"
