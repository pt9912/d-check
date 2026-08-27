#!/usr/bin/env bash
# Proben gegen den Tool-Call-Wächter. Liegt als DATEI vor: die Proben-Daten
# enthalten die blockierten Wörter, und die quote-blinde Segmentierung des
# Wächters würde den Aufruf sonst selbst blocken.
set -uo pipefail
cd "$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
G=.claude/hooks/pretooluse-command-guard.sh
fails=0

probe() {  # probe <erwartung: block|pass> <kommando>
  local want="$1" cmd="$2" n got mark
  n=$(printf '{"tool_input":{"command":"%s"}}' "$cmd" | bash "$G" | grep -c 'decision' || true)
  got=pass; [ "$n" != "0" ] && got=block
  mark="  ok "; [ "$got" != "$want" ] && { mark="FAIL"; fails=$((fails+1)); }
  printf '%s  %-6s %-44s (erwartet %s)\n' "$mark" "$got" "${cmd:0:44}" "$want"
}

raw() {  # raw <erwartung> <label> <roh-json>
  local want="$1" label="$2" json="$3" n got mark
  n=$(printf '%s' "$json" | bash "$G" | grep -c 'decision' || true)
  got=pass; [ "$n" != "0" ] && got=block
  mark="  ok "; [ "$got" != "$want" ] && { mark="FAIL"; fails=$((fails+1)); }
  printf '%s  %-6s %-44s (erwartet %s)\n' "$mark" "$got" "$label" "$want"
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

echo "== Quote-blinde Falsch-Positiv-Klasse: benannte Grenze, nicht behoben"
raw block 'Interpreter als Daten in Klammern' '{"tool_input":{"command":"sed -i /Bash(python3 -)/d f"}}'
raw pass  'dieselben Daten ohne Leerzeichen'  '{"tool_input":{"command":"sed -i /Bash(python3)/d f"}}'
raw block 'blockiertes Wort nach Trenner in Text' '{"tool_input":{"command":"echo a; pip b"}}'

printf '\n== Fehlschläge: %d\n' "$fails"
