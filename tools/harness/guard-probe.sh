#!/usr/bin/env bash
# Proben gegen den Tool-Call-Wächter. Liegt als DATEI vor: die Proben-Daten
# enthalten die blockierten Wörter, und die quote-blinde Segmentierung des
# Wächters würde den Aufruf sonst selbst blocken.
set -uo pipefail
cd "$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
G=.claude/hooks/pretooluse-command-guard.sh
fails=0

out=$(mktemp)
trap 'rm -f "$out"' EXIT

# Der Exit des Wächters wird MITGENOMMEN, nicht hinter einer Pipe verworfen:
# ein Wächter, der abstürzt, gibt nichts aus — und wäre sonst von „ausdrücklich
# erlaubt" nicht zu unterscheiden. Darum das dritte Verdikt `crash`.
verdict() {  # verdict <roh-json> -> setzt GOT
  local rc
  set +e
  printf '%s' "$1" | bash "$G" >"$out" 2>/dev/null
  rc=$?
  set -e
  if grep -q '"permissionDecision": "deny"' "$out"; then
    # Beide Kanäle gehören zusammen: fehlt der Exit-Riegel, hängt der Block
    # allein an der Antwortform. Das ist ein eigenes Verdikt, kein `block`.
    if [ "$rc" -eq 0 ]; then GOT=halb; else GOT=block; fi
  elif [ "$rc" -ne 0 ]; then GOT=crash
  else GOT=pass
  fi
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
guard_ohne_pfad() {
  local rc
  set +e
  printf '{"tool_input":{"command":"make gates"}}' \
    | env -i PATH=/nonexistent /bin/bash "$PWD/$G" >"$out" 2>/dev/null
  rc=$?
  set -e
  if grep -q '"permissionDecision": "deny"' "$out"; then
    # Beide Kanäle gehören zusammen: fehlt der Exit-Riegel, hängt der Block
    # allein an der Antwortform. Das ist ein eigenes Verdikt, kein `block`.
    if [ "$rc" -eq 0 ]; then GOT=halb; else GOT=block; fi
  elif [ "$rc" -ne 0 ]; then GOT=crash
  else GOT=pass
  fi
}
guard_ohne_pfad
report block 'leeres PATH: awk fehlt'

echo "== Quote-blinde Falsch-Positiv-Klasse: benannte Grenze, nicht behoben"
raw block 'Interpreter als Daten in Klammern' '{"tool_input":{"command":"sed -i /Bash(python3 -)/d f"}}'
raw pass  'dieselben Daten ohne Leerzeichen'  '{"tool_input":{"command":"sed -i /Bash(python3)/d f"}}'
raw block 'blockiertes Wort nach Trenner in Text' '{"tool_input":{"command":"echo a; pip b"}}'

printf '\n== Fehlschläge: %d\n' "$fails"
