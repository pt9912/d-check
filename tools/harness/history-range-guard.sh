#!/usr/bin/env bash
# history-range-guard — Vorlauf-Waechter fuer die zwei history-lesenden d-check-Targets
# (`doc-immutable` ueber das Modul vcs, `doc-commits` ueber commits). Er prueft VOR dem
# Modul-Lauf, dass die angeforderte Range git-seitig AUFLOESBAR und NICHT LEER ist.
#
# ANLASS: ein Klon der Tiefe 1 liefert fuer eine AUFLOESBARE, aber dort LEERE Range
# (z. B. `HEAD..HEAD`) "0 Befund(e)", Exit 0 — das Modul unterscheidet das nicht von
# "wirklich nichts zu melden". Das Ergebnis ist gruen ueber leerem Pruefbereich.
#
# GRENZE — was dieser Waechter NICHT deckt: eine UNAUFLOESBARE Basis (z. B. `HEAD~1` im
# Klon der Tiefe 1) bricht schon ohne ihn ab (git kennt den Commit dort nicht). Ein
# `--staged`-Lauf vergleicht den Index gegen HEAD, braucht keine Tiefe > 1 und loest den
# Range-Check darum nicht aus; den Leerfall meldet er trotzdem.
#
# Kein Docker, kein Netz: bash + git. Das Doc-Gate-Fragment haengt ihn darum als
# Vorbedingung vor die zwei Targets — er laeuft, bevor der Container startet.
set -euo pipefail

# depth_info liefert die Zahl der Shallow-Grenzen fuer die Meldung: die Zeilenzahl in
# .git/shallow (eine Zeile je Shallow-Grenze), oder "voll", wenn die Datei fehlt. Nur
# fuer die Meldung — die Entscheidung haengt an der Commit-Zahl der Range, nicht an
# diesem Wert: `decide()` bewertet ausschliesslich `count`.
depth_info() {
  if [ -f .git/shallow ]; then
    wc -l <.git/shallow | tr -d ' '
  else
    echo "voll (kein Shallow-Klon)"
  fi
}

# decide <range> <count> — REIN: kein git-Aufruf, nur die schon ermittelte Commit-Zahl
# wird bewertet. Exit 0: nicht leer, OK. Exit 1: leer (0 Commits) — der Fall, den dieser
# Waechter faengt. Exit 2: `count` ist keine Ganzzahl — fail-closed statt in den OK-Zweig
# zu fallen; ueber den produktiven Pfad nicht erreichbar (`git rev-list --count` liefert
# bei Erfolg stets eine Ganzzahl).
decide() {
  local range="$1" count="$2"
  case "$count" in
    '' | *[!0-9]*)
      echo "history-range-guard: Range '$range' liefert keine gueltige Commit-Zahl ('$count')." >&2
      return 2
      ;;
  esac
  if [ "$count" -eq 0 ]; then
    echo "history-range-guard: Range '$range' ist aufloesbar, aber LEER (0 Commits)." >&2
    echo "  Shallow-Grenzen: $(depth_info)" >&2
    echo "  Angeforderte Range: $range" >&2
    echo "  -> Checkout braucht 'fetch-depth: 0' (oder ausreichende Tiefe) fuer diesen Job." >&2
    return 1
  fi
  echo "history-range-guard: Range '$range' aufgeloest, $count Commit(s) — OK."
}

# decide_staged <has_staged> — der Leerfall des `--staged`-Zweigs wird GEMELDET statt
# schweigend durchgereicht (dieselbe Klasse wie bei der leeren Range); Exit 0 in beiden
# gueltigen Faellen — `--staged` loest den Range-Check nicht aus. Exit 2 bei einem Wert
# ausserhalb `0|1`, derselbe fail-closed-Grundsatz wie in `decide()`; ueber den
# produktiven Pfad nicht erreichbar (der `--staged`-Zweig uebergibt die Literale 0/1).
decide_staged() {
  local has_staged="$1"
  case "$has_staged" in
    0 | 1) ;;
    *)
      echo "history-range-guard: --staged erwartet 0 oder 1, nicht '$has_staged'." >&2
      return 2
      ;;
  esac
  if [ "$has_staged" -eq 0 ]; then
    echo "history-range-guard: --staged ohne gestagte Aenderung — nichts zu pruefen." >&2
  fi
}

range="${1:?Usage: history-range-guard.sh <base>..<head> | --staged}"

# --staged vergleicht den Index gegen HEAD und braucht keine Tiefe > 1 — der Range-Check
# wird nicht ausgeloest.
if [ "$range" = "--staged" ]; then
  cd "$(git rev-parse --show-toplevel)"
  if git diff --cached --quiet; then
    decide_staged 0
  else
    decide_staged 1
  fi
  exit 0
fi

cd "$(git rev-parse --show-toplevel)"

count="$(git rev-list --count "$range" 2>/dev/null)" || {
  echo "history-range-guard: Range '$range' ist NICHT aufloesbar (Basis fehlt im Klon?)." >&2
  echo "  Shallow-Grenzen: $(depth_info)" >&2
  exit 2
}

rc=0
decide "$range" "$count" || rc=$?
exit "$rc"
