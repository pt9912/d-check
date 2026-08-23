#!/usr/bin/env bash
# workflow-pins — Wächter zu AGENTS.md §3.9: jeder `uses:`-Eintrag in
# .github/workflows/ nennt einen VOLLEN Commit-SHA mit Tag-Kommentar dahinter.
#
# WAS ER ZUSICHERT, GENAU: die FORM des Pins, nicht seine Gültigkeit. Ob der
# SHA existiert und ob er den Commit bezeichnet, den der Tag-Kommentar
# behauptet, ist eine Netz-Frage und ausdruecklich NICHT Teil der Zusage --
# dieser Wächter ist netzlos und laeuft in `make gates`.
#
# ER LIEST NUR ECHTE YAML-SCHLUESSEL. Die Workflow-Koepfe erklaeren die
# Pin-Konvention in Prosa und nennen `uses:` dabei mehrfach; ein Wächter, der
# auf die Zeichenkette trifft statt auf den Schluessel, meldete seine eigene
# Dokumentation. Gefiltert wird deshalb auf `^\s*(- )?uses:` -- Kommentarzeilen
# beginnen nach Leerraum mit `#` und fallen heraus.
#
# FAIL-CLOSED BEI LEERER MENGE: findet er keine Workflow-Datei oder keinen
# einzigen `uses:`-Schluessel, bricht er ROT ab statt gruen zu melden. Ein
# Wächter, der nichts gefunden hat, und einer, der nichts zu pruefen hatte,
# sehen im Exit-Code sonst identisch aus.
#
# Exit: 0 = alle Pins formgerecht, 1 = mindestens ein Befund oder leere Menge.
# bash + coreutils + grep.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

dir=".github/workflows"
[ -d "$dir" ] || { echo "workflow-pins: ${dir} fehlt" >&2; exit 1; }

findings=0
checked=0

while IFS= read -r hit; do
  file="${hit%%:*}"
  rest="${hit#*:}"
  line="${rest%%:*}"
  text="${rest#*:}"
  checked=$((checked + 1))

  # Form: <ref>@<40 Hex> gefolgt von einem Kommentar, der den Tag nennt.
  if ! printf '%s' "$text" | grep -qE '@[0-9a-f]{40}([[:space:]]|$)'; then
    echo "${file}:${line}	uses-pin-missing	kein voller 40-stelliger Commit-SHA"
    findings=$((findings + 1))
    continue
  fi
  if ! printf '%s' "$text" | grep -qE '@[0-9a-f]{40}[[:space:]]+#[[:space:]]*\S'; then
    echo "${file}:${line}	uses-pin-untagged	SHA ohne Tag-Kommentar dahinter"
    findings=$((findings + 1))
  fi
done < <(grep -nE '^[[:space:]]*(- )?uses:[[:space:]]' "${dir}"/*.yml 2>/dev/null | sed "s|^|${dir}/|; s|^${dir}/${dir}/|${dir}/|")

if [ "$checked" -eq 0 ]; then
  echo "workflow-pins: kein einziger uses:-Schlüssel gefunden — leere Prüfmenge, fail-closed" >&2
  exit 1
fi

if [ "$findings" -gt 0 ]; then
  echo "workflow-pins: ${findings} Befund(e) über ${checked} uses:-Einträge (AGENTS.md §3.9)" >&2
  exit 1
fi

echo "workflow-pins: ok (${checked} uses:-Einträge, alle SHA-gepinnt mit Tag-Kommentar)"
