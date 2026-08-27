#!/usr/bin/env bash
# nightly-state — liest den Ausgang des juengsten Nachtlaufs und sagt, ob er
# gelesen werden muss. Er ERSETZT den Nachtlauf nicht; er macht seinen Stand an
# dem Moment verfuegbar, an dem jemand ohnehin hinsieht: der Slice-Planung.
#
# WARUM NICHT `gh`: der Host dieses Repos ist auf git, make, bash, Docker und
# die POSIX-Werkzeuge festgelegt; die Netz-Targets erwarten zusaetzlich curl
# (AGENTS.md §3.1). `gh` waere eine neue Erwartung fuer einen Aufruf, den curl
# ebenso beantwortet — die GitHub-API liefert den Lauf-Ausgang eines
# oeffentlichen Repos ohne Token.
#
# FAIL-OPEN wie die Frische-Achsen: Netz-, Werkzeug- oder Parse-Ausfall meldet
# SKIP und Exit 0. Diese Vorpruefung darf niemanden aufhalten — sie ist ein
# Lese-Schritt, kein Gate. Aus demselben Grund ist sie NICHT in `gates`.
#
# EXIT: immer 0. Der Ausgang steht in der AUSGABE, nicht im Code — wer sie
# nicht liest, hat den Schritt nicht getan, und das soll ein Exit-Code nicht
# verdecken koennen.
set -uo pipefail

CT=10
MT=30
REPO="${NIGHTLY_REPO:-pt9912/d-check}"
WF="${NIGHTLY_WORKFLOW:-upstream-drift.yml}"

command -v curl >/dev/null 2>&1 || {
  echo "nightly-state: 'curl' nicht gefunden — SKIP"
  exit 0
}

json="$(curl -fsSL --connect-timeout "$CT" --max-time "$MT" \
  "https://api.github.com/repos/${REPO}/actions/workflows/${WF}/runs?per_page=1" \
  2>/dev/null || true)"

if [ -z "$json" ]; then
  echo "nightly-state: SKIP — Lauf-Stand nicht abrufbar (Netz oder API)"
  exit 0
fi

# Ein Feld je Zeile, dann die drei gesuchten greifen. Kein jq: dieselbe
# Sparsamkeit wie bei den Frische-Achsen (DC-QA-03).
feld() { printf '%s' "$json" | tr ',' '\n' | grep -m1 "\"$1\"" | sed 's/^[^:]*: *//; s/^"//; s/"$//'; }

conclusion="$(feld conclusion)"
created="$(feld created_at)"
url="$(printf '%s' "$json" | tr ',' '\n' | grep -m1 '"html_url".*actions/runs' | sed 's/^[^:]*: *//; s/^"//; s/"$//')"

case "$conclusion" in
  success)
    echo "nightly-state: gruen — juengster Lauf ${created}"
    ;;
  "")
    echo "nightly-state: SKIP — kein Ausgang im Lauf-Objekt (laeuft er gerade?)"
    ;;
  *)
    echo "nightly-state: ROT (${conclusion}) — juengster Lauf ${created}"
    echo "nightly-state: ${url}"
    echo "nightly-state: LESEN, nicht wegklicken. Eine planmaessige Meldung"
    echo "nightly-state: (Fremd-Release; Zitat-Spanne nach einem Bump, MR-051)"
    echo "nightly-state: wird anders behandelt als eine unerwartete — der"
    echo "nightly-state: Unterschied steht in der AUSGABE des Laufs, nicht in"
    echo "nightly-state: seiner Farbe."
    ;;
esac
exit 0
