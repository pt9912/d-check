#!/usr/bin/env bash
# nightly-state — liest den Ausgang des juengsten Nachtlaufs und sagt, ob er
# gelesen werden muss. Er ERSETZT den Nachtlauf nicht; er macht seinen Stand an
# dem Moment verfuegbar, an dem jemand ohnehin hinsieht: der Slice-Planung
# (dritte Vorpruefung, MR-053).
#
# WARUM NICHT `gh`: der Host dieses Repos ist auf git, make, bash, Docker und
# die POSIX-Werkzeuge festgelegt; die Netz-Skripte erwarten zusaetzlich curl
# (AGENTS.md §3.1 — dieses Skript ist dort als drittes genannt). `gh` waere eine
# neue Erwartung fuer einen Aufruf, den curl ebenso beantwortet.
#
# KEIN jq: das Werkzeug steht nicht in der Host-Klasse (§3.1). Die Antwort wird
# feldweise gelesen, und JEDER gelesene Wert wird gegen sein zulaessiges
# Vokabular geprueft, bevor er als Ausgang gilt — sonst waere eine
# HTTP-200-Nicht-Lauf-Antwort ein gueltiges Urteil statt des zugesagten SKIP.
#
# FAIL-OPEN: Netz-, Werkzeug- oder Parse-Ausfall meldet SKIP. IMMER Exit 0 —
# nicht weil ein Exit-Code etwas verdecken koennte (beide werden emittiert),
# sondern damit dieser reine Lese-Schritt keine Kette bricht und niemanden
# aufhaelt. PREIS, benannt: der Ausgang ist damit maschinell nicht mehr
# unterscheidbar (ROT und SKIP sehen fuer ein `make`-Glied gleich aus) — anders
# als bei pin-freshness.sh, das Exit 3 traegt. Deshalb liegen die Urteile auf
# STDERR, wie dort: ein `>/dev/null` loescht sonst das einzige Signal.
#
# NETZLOS PRUEFBAR: `--parse <datei>` fuehrt nur die Feld-Extraktion und das
# Urteil; `--selftest` faehrt die eingebauten Proben. Ohne diesen Einstieg waere
# die Semantik nur mit Netz zu pruefen, und damit gar nicht.
#
# GRENZEN, benannt: (1) Es liest den JUENGSTEN Lauf, nicht sein Alter — ein
# abgeschalteter Nachtlauf meldet weiter "gruen" (das Datum steht in der Zeile,
# das Urteilswort sagt es nicht). (2) Bei privatem Repo oder umbenanntem
# Workflow scheitert curl und meldet dauerhaft SKIP, ununterscheidbar von einer
# Netzstoerung. (3) Der Repo-Slug ist ein Default, kein Fund: in einem Fork
# meldet es ohne NIGHTLY_REPO den Nachtlauf des Originals.
set -uo pipefail

CT=10
MT=30
REPO="${NIGHTLY_REPO:-pt9912/d-check}"
WF="${NIGHTLY_WORKFLOW:-upstream-drift.yml}"

# Ein Feld je Zeile greifen. Das Muster ist an BEIDE Anfuehrungszeichen
# verankert; in einem JSON-String ist ein literales " als \" maskiert, die
# Teilkette "conclusion" kann in einem Wert also nicht entstehen.
feld() {
  printf '%s' "$1" | tr ',' '\n' | grep -m1 "\"$2\"" \
    | sed 's/^[^:]*: *//; s/^"//; s/[",}]*[[:space:]]*$//'
}

urteil() {
  local json="$1" conclusion created url
  conclusion="$(feld "$json" conclusion)"
  created="$(feld "$json" created_at)"
  url="$(printf '%s' "$json" | tr ',' '\n' | grep -m1 '"html_url".*actions/runs' \
         | sed 's/^[^:]*: *//; s/^"//; s/[",}]*[[:space:]]*$//')"

  # Form pruefen, bevor der Wert als Stand gilt: ein Zeitstempel, der keiner
  # ist, macht aus einem Parse-Ausfall ein Urteil.
  printf '%s' "$created" | grep -qE '^[0-9]{4}-[0-9]{2}-[0-9]{2}T' || {
    echo "nightly-state: SKIP — keine Lauf-Antwort (Form nicht erkannt)" >&2
    return 0
  }

  case "$conclusion" in
    success)
      echo "nightly-state: gruen — juengster Lauf ${created}"
      ;;
    null|"")
      # Die API liefert bei laufendem Job JSON-null, nicht ein fehlendes Feld.
      echo "nightly-state: SKIP — der juengste Lauf (${created}) hat noch keinen Ausgang" >&2
      ;;
    failure)
      echo "nightly-state: ROT (failure) — juengster Lauf ${created}" >&2
      echo "nightly-state: ${url}" >&2
      echo "nightly-state: LESEN, nicht wegklicken. Eine PLANMAESSIGE Meldung" >&2
      echo "nightly-state: (Fremd-Release; Zitat-Spanne nach einem Bump, MR-051)" >&2
      echo "nightly-state: wird anders behandelt als eine unerwartete — der" >&2
      echo "nightly-state: Unterschied steht in der AUSGABE des Laufs, nicht in" >&2
      echo "nightly-state: seiner Farbe." >&2
      ;;
    *)
      # cancelled, timed_out, startup_failure, action_required, neutral,
      # skipped, stale: kein Achsen-Urteil, sondern eine Lauf-Stoerung. Der
      # Deutungsrahmen oben waere hier schlicht falsch.
      echo "nightly-state: ROT (${conclusion}) — juengster Lauf ${created}" >&2
      echo "nightly-state: ${url}" >&2
      echo "nightly-state: Lauf-Stoerung, KEIN Achsen-Urteil — die Achsen haben" >&2
      echo "nightly-state: nichts gemeldet. Neu ausloesen, dann lesen." >&2
      ;;
  esac
}

if [ "${1:-}" = "--parse" ]; then
  [ -r "${2:-}" ] || { echo "nightly-state: --parse braucht eine lesbare Datei" >&2; exit 0; }
  urteil "$(cat "$2")"
  exit 0
fi

if [ "${1:-}" = "--selftest" ]; then
  tmp="${TMPDIR:-/tmp}/nightly-selftest.$$"
  fails=0
  probe() {
    local name="$1" json="$2" erwartet="$3" got
    printf '%s' "$json" > "$tmp"
    got="$(urteil "$(cat "$tmp")" 2>&1 | head -1)"
    case "$got" in
      *"$erwartet"*) printf '  ok   %-34s %s\n' "$name" "$erwartet" ;;
      *) printf '  FAIL %-34s erwartet %s, war: %s\n' "$name" "$erwartet" "$got"; fails=$((fails + 1)) ;;
    esac
  }
  probe "abgeschlossen, gruen"   '{"created_at": "2026-08-27T01:00:00Z", "conclusion": "success"}'          'gruen'
  probe "abgeschlossen, rot"     '{"created_at": "2026-08-27T01:00:00Z", "conclusion": "failure"}'          'ROT (failure)'
  probe "laeuft gerade (null)"   '{"created_at": "2026-08-27T01:00:00Z", "conclusion": null}'               'noch keinen Ausgang'
  probe "abgebrochen"            '{"created_at": "2026-08-27T01:00:00Z", "conclusion": "cancelled"}'        'ROT (cancelled)'
  probe "Nicht-Lauf-Antwort"     '{"message": "Not Found"}'                                                 'Form nicht erkannt'
  probe "leere Lauf-Liste"       '{"total_count": 0, "workflow_runs": []}'                                  'Form nicht erkannt'
  rm -f "$tmp"
  echo
  echo "== Fehlschlaege: $fails"
  [ "$fails" -eq 0 ]
  exit $?
fi

command -v curl >/dev/null 2>&1 || {
  echo "nightly-state: 'curl' nicht gefunden — SKIP" >&2
  exit 0
}

json="$(curl -fsSL --connect-timeout "$CT" --max-time "$MT" \
  "https://api.github.com/repos/${REPO}/actions/workflows/${WF}/runs?per_page=1" \
  2>/dev/null || true)"

if [ -z "$json" ]; then
  echo "nightly-state: SKIP — Lauf-Stand nicht abrufbar (Netz, privates Repo oder umbenannter Workflow)" >&2
  exit 0
fi

urteil "$json"
exit 0
