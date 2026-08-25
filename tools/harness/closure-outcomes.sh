#!/usr/bin/env bash
# closure-outcomes — Waechter zur Drei-Ausgaenge-Regel des Baseline-Regelwerks
# (modul-05 §Offene Risiken werden bei Closure aufgeloest): "Ein Slice geht
# nicht nach done/, waehrend ein Risiko ohne Ausgang dasteht."
#
# WAS ER ZUSICHERT, GENAU: kein Slice in done/ traegt einen UNAUFGELOESTEN
# Vorlagen-Platzhalter. Das ist die urteilsfreie Haelfte der Regel -- OB der
# eingetragene Ausgang inhaltlich traegt, bleibt Urteil und ist ausdruecklich
# nicht Gegenstand.
#
# INLINE-CODE WIRD UEBERSPRUNGEN, und das ist tragend: ein Slice, der ueber die
# Platzhalter SCHREIBT, nennt sie zwangslaeufig. Ohne diese Trennung koennte
# genau der Slice, der diesen Waechter einfuehrt, nie schliessen -- er meldete
# seine eigene Dokumentation. Entfernt werden `…`-Spannen paarweise; eine
# ungerade Backtick-Zahl laesst den Rest stehen und damit sichtbar (fail-loud).
#
# FAIL-CLOSED BEI LEERER MENGE: findet er keinen done/-Slice, bricht er ROT ab.
# "Nichts gefunden" und "nichts zu pruefen" duerfen im Exit nicht gleich
# aussehen -- dieselbe Klasse, die in dieser Arbeit schon dreimal aufgefallen ist.
#
# GRENZE, benannt: die Platzhalter-Liste ist eine LISTE. Aendert die Vorlage
# ihre Form, greift der Waechter nicht mehr und schweigt. Sie gehoert beim
# naechsten Vorlagen-Bump mitgeprueft.
#
# Exit: 0 = sauber, 1 = mindestens ein Befund oder leere Pruefmenge.
# bash + coreutils + grep + sed.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

dir="docs/plan/planning/done"
[ -d "$dir" ] || { echo "closure-outcomes: ${dir} fehlt" >&2; exit 1; }

mapfile -t files < <(find "$dir" -maxdepth 1 -type f -name 'slice-*.md' | LC_ALL=C sort)
if [ "${#files[@]}" -eq 0 ]; then
  echo "closure-outcomes: kein Slice in ${dir} — leere Prüfmenge, fail-closed" >&2
  exit 1
fi

# Die drei Formen: zwei repo-eigene und der Platzhalter der Kanon-Vorlage.
PATTERNS=(
  '(bei Closure)'
  'wird mit dem Closure-Body gefüllt'
  '<…>'
)

findings=0
for f in "${files[@]}"; do
  # Inline-Code paarweise entfernen, dann Zeilennummern erhalten.
  while IFS= read -r hit; do
    line="${hit%%:*}"
    text="${hit#*:}"
    for pat in "${PATTERNS[@]}"; do
      case "$text" in
        *"$pat"*)
          echo "${f}:${line}	closure-outcome-open	unaufgelöster Vorlagen-Platzhalter: ${pat}"
          findings=$((findings + 1))
          ;;
      esac
    done
  done < <(sed 's/`[^`]*`//g' "$f" | grep -n . || true)
done

if [ "$findings" -gt 0 ]; then
  echo "closure-outcomes: ${findings} Befund(e) über ${#files[@]} Slices in done/ — ein Risiko ohne Ausgang (modul-05)" >&2
  exit 1
fi

echo "closure-outcomes: ok (${#files[@]} Slices in done/, kein offener Platzhalter)"
