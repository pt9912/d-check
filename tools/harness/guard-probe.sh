#!/usr/bin/env bash
# Proben gegen den Tool-Call-Wächter. Liegt als DATEI vor: die Proben-Daten
# enthalten die blockierten Wörter, und die quote-blinde Segmentierung des
# Wächters würde den Aufruf sonst selbst blocken.
set -uo pipefail
cd "$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
G=.claude/hooks/pretooluse-command-guard.sh
fails=0

out=$(mktemp)
tmpin=$(mktemp)
trap 'rm -f "$out" "$tmpin"' EXIT

# Der Exit des Wächters wird MITGENOMMEN, nicht hinter einer Pipe verworfen:
# ein Wächter, der abstürzt, gibt nichts aus — und wäre sonst von „ausdrücklich
# erlaubt" nicht zu unterscheiden. Vier Verdikte, eine Auswertung: `block`
# (abgelehnt über beide Kanäle), `halb` (abgelehnt, aber ohne Exit-Riegel),
# `crash` (keine Ausgabe), `pass` (durchgelassen).
klassifiziere() {  # klassifiziere <rc> -> setzt GOT
  if grep -q '"permissionDecision": "deny"' "$out"; then
    if [ "$1" -eq 0 ]; then GOT=halb; else GOT=block; fi
  elif [ "$1" -ne 0 ]; then GOT=crash
  else GOT=pass
  fi
}

verdict() {  # verdict <roh-json> -> setzt GOT
  local rc
  set +e
  printf '%s' "$1" | bash "$G" >"$out" 2>/dev/null
  rc=$?
  set -e
  klassifiziere "$rc"
}

report() {  # report <erwartung> <label>
  local mark="  ok "
  [ "$GOT" != "$1" ] && { mark="FAIL"; fails=$((fails+1)); }
  printf '%s  %-6s %-44s (erwartet %s)\n' "$mark" "$GOT" "${2:0:44}" "$1"
}

probe() {  # probe <erwartung: block|pass|crash|halb> <kommando>
  verdict "$(printf '{"tool_input":{"command":"%s"}}' "$2")"
  report "$1" "$2"
}

raw() {  # raw <erwartung> <label> <roh-json>
  verdict "$3"
  report "$1" "$2"
}

echo "== Segmentierung: Brace-Group und einzelnes kaufmännisches Und"
probe block '{ go build ./... ; }'
probe block 'echo x & pip install y'

echo "== Interpreter-Sperre (MR-040, MR-042)"
probe block 'python3 x.py'
probe block '/usr/bin/python3.12 x.py'
probe block 'node -e 1'
probe block 'perl -e 1'
probe block 'ruby -e 1'
probe block 'uv run x.py'
probe block 'uvx ruff check'

echo "== Paketmanager und Host-Go (MR-005)"
probe block 'pip install foo'
probe block 'sudo apt-get install foo'
probe block '/usr/bin/go test ./...'
raw   block 'bash -lc mit Sub-Shell' '{"tool_input":{"command":"bash -lc \"go build\""}}'
probe block 'FOO=1 cargo build'

echo "== Gegenkontrollen, müssen durchlaufen"
probe pass 'make gates'
probe pass 'git commit -F msg.txt'
probe pass 'docker run img npm test'
probe pass 'grep -rn pip AGENTS.md'
probe pass 'sed -n 1,5p Makefile'

echo "== Fail-closed am Extraktor: Parse-Zweifel blockiert"
raw block 'malformed JSON'        '{"tool_input":{"command":'
raw block 'abgeschnitten'         '{"tool_input":{"command":"make gates'
raw block 'kein JSON-Objekt'      'einfach text'
raw block 'u-Escape im Befehl'    '{"tool_input":{"command":"make \u0067ates"}}'
raw block 'u-Escape zu kurz'      '{"tool_input":{"command":"make \u06"}}'
raw block 'Müll ausserhalb String' '{"tool_input":{"command":"a" garbage }}'
raw pass  'kein command-Feld'     '{"tool_input":{"other":"x"}}'
raw pass  'command als VALUE-Wort' '{"tool_input":{"command":"echo command pip"}}'
raw block 'zwei Strings ohne Trenner'  '{"tool_input":{"command":"pip x""make gates"}}'
raw block 'zwei Werte mit Leerzeichen' '{"tool_input":{"command":"pip x" "make gates"}}'
raw block 'u-Escape im Schluessel'     '{"tool_input":{"\u0063ommand":"pip x"}}'
raw block 'u-Escape im Eltern-Schluessel' '{"\u0074ool_input":{"command":"pip x"}}'

echo "== Der Wächter selbst: ohne Host-PATH bleibt er fail-closed"
ohne_pfad() {
  local rc
  set +e
  printf '{"tool_input":{"command":"make gates"}}' \
    | env -i PATH=/nonexistent /bin/bash "$PWD/$G" >"$out" 2>/dev/null
  rc=$?
  set -e
  klassifiziere "$rc"
}
ohne_pfad
report block 'leeres PATH: awk fehlt'

# Der Ausfall, der diesen Slice ausgelöst hat, lag NICHT im Urteil, sondern im
# LESEWEG: eine stdin-Form, die der Guard nicht lesen konnte, beendete ihn ohne
# Ausgabe — und eine Probenmenge, die stdin nur als Pipe bespielt, sieht genau
# das nicht. Jede Form, die ein Aufrufer liefern kann, gehört gefahren.
# Erwartung ist überall `block`: entweder wegen des Kommandos oder, wo gar keine
# Eingabe ankommt, wegen fehlender Parse-Grundlage.
echo "== Leseweg: jede stdin-Form, die ein Aufrufer liefern kann"
stdin_form() {  # stdin_form <erwartung> <label> <pipe|datei|herestring|devnull|zu>
  local rc j='{"tool_input":{"command":"pip --version"}}'
  set +e
  case "$3" in
    pipe)       printf '%s' "$j" | bash "$G" >"$out" 2>/dev/null ;;
    datei)      printf '%s' "$j" >"$tmpin"; bash "$G" >"$out" 2>/dev/null <"$tmpin" ;;
    herestring) bash "$G" >"$out" 2>/dev/null <<<"$j" ;;
    devnull)    bash "$G" >"$out" 2>/dev/null </dev/null ;;
    zu)         bash "$G" >"$out" 2>/dev/null <&- ;;
  esac
  rc=$?
  set -e
  klassifiziere "$rc"
  report "$1" "$2"
}
stdin_form block 'stdin als Pipe'          pipe
stdin_form block 'stdin als Datei'         datei
stdin_form block 'stdin als Here-String'   herestring
stdin_form block 'stdin auf /dev/null'     devnull
stdin_form block 'stdin GESCHLOSSEN'       zu

echo "== Quote-blinde Falsch-Positiv-Klasse: benannte Grenze, nicht behoben"
raw block 'Interpreter als Daten in Klammern' '{"tool_input":{"command":"sed -i /Bash(python3 -)/d f"}}'
raw pass  'dieselben Daten ohne Leerzeichen'  '{"tool_input":{"command":"sed -i /Bash(python3)/d f"}}'
raw block 'blockiertes Wort nach Trenner in Text' '{"tool_input":{"command":"echo a; pip b"}}'

printf '\n== Fehlschläge: %d\n' "$fails"
# Das Urteil gehört in den Exit, nicht nur in die Prosa: ein Aufrufer, der die
# Ausgabe nicht liest, muss den Fehlschlag trotzdem sehen.
[ "$fails" -eq 0 ]
