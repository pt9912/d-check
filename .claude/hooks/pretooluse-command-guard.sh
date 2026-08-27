#!/usr/bin/env bash
# pretooluse-command-guard — verbietet Host-Paketmanager, die Host-Go-Toolchain
# und Host-Skript-Interpreter (d-check ist make/Docker-only, AGENTS.md §3.1,
# ADR-0001).
#
# Reines bash + awk. Die JSON-Extraktion liegt in tools/harness/extract-command.awk
# und meldet Parse-Zweifel per Exit 3; dieser Guard blockt dann (fail-closed).
#
# Geprüft wird die Befehlsposition jedes Kommando-Segments (Trennung an
# ; & && || | $( ` ( und Zeilenenden). `/usr/bin/pip` und `sudo pip` werden
# erkannt; Zuweisungs-, Wrapper- und Brace-Group-Präfixe werden übersprungen —
# ohne Letztere entkäme das Werkzeug in `{ go build; }` als zweites Token.
# Sub-Shell-Strings (`bash -c "…"`, auch in Flag-Bündeln wie -lc/-ec/-cx)
# werden rekursiv geprüft, Tiefe ≤ 3, darüber fail-closed (MR-005).
#
# SKRIPT-INTERPRETER: `python` samt Versions-Suffix, `perl`, `ruby`, `node`,
# `deno`, `bun` sind blockiert (MR-040, MR-042). Nicht die Sprache ist der
# Grund, sondern der Skopus: ein General-Purpose-Interpreter kann alles, was die
# in §3.1 genannten Werkzeuge können, ohne deren Grenze zu erben.
#
# GRENZE, quote-blind: ein Trenner INNERHALB eines Arguments startet ein neues
# Segment. Ein Kommando, das ein blockiertes Wort als DATEN trägt, wird deshalb
# geblockt (Falsch-Positiv). Abhilfe ist die Repo-Praxis ohnehin: solchen Inhalt
# in eine Datei legen — `git commit -F <datei>`, Proben als Skript.
#
# GRENZE, Umfang — gemessen, nicht geschätzt. Ungeprüft bleiben:
#   1. Segment-Köpfe, die Shell-Schlüsselwörter sind: `if true; then pip …; fi`,
#      `for f in a; do go build; done`, `while …; do npm i; done`, `! pip …`.
#   2. Wrapper außerhalb von PREFIXES: `nohup`, `timeout 5`, `stdbuf -o0`,
#      `setsid`. Die Liste ist eine Liste und damit unvollständig.
#   3. Wort-interne Splices: `p"i"p`, `g\o` — bash setzt sie zusammen, die
#      Token-Sicht hier nicht.
#   4. Verschachtelung mit escapten Quotes: bei `bash -c "bash -c \"…\""`
#      greift strip_quotes nicht, die Rekursion endet.
# Dazu `find -exec`, `awk`-Programme und jeder Interpreter, den die Liste nicht
# kennt. Stolperdraht, keine Sandbox.
#
# Im Pass-Fall: KEINE Ausgabe — "approve" würde das Permission-System
# überspringen; ohne Ausgabe läuft die normale Permission-Entscheidung.
set -euo pipefail

# Pfad und Eingabe ohne externes Programm: `dirname` und `cat` wären zwei
# weitere Host-Abhängigkeiten, und ein fehlendes PATH ließe den Guard mit 127
# enden — ohne Ausgabe, also OHNE Block. Fail-closed hängt sonst an Werkzeugen,
# die dieser Guard nicht selbst prüft.
here=${BASH_SOURCE[0]%/*}
extractor="$here/../../tools/harness/extract-command.awk"

# ZWEI KANÄLE, absichtlich beide. Die JSON-Antwort trägt den Grund, den der
# Aufrufer wörtlich sieht; der Nicht-Null-Exit blockt unabhängig von jeder
# Antwortform. Sie widersprechen sich nicht — der Grund erscheint, der Exit
# liegt als Riegel darunter. Ohne den zweiten Kanal hinge JEDER Block an einem
# Format, dessen Auslegung dieser Guard nicht kontrolliert (MR-044).
emit_block() {
  printf '%s\n' '{' \
    '  "hookSpecificOutput": {' \
    '    "hookEventName": "PreToolUse",' \
    '    "permissionDecision": "deny",' \
    '    "permissionDecisionReason": "d-check is make/Docker-only (AGENTS.md §3.1, ADR-0001). Use make targets and the POSIX host tools the gate scripts use (grep/sed/awk/find); do not run host package managers, host go, or host script interpreters (apt/brew/pip/npm/cargo/go/python/perl/ruby/node). On parse doubt the guard fails closed."' \
    '  }' \
    '}'
  exit 2
}

# Host-Go: ADR-0001 + AGENTS §3.1. Paketmanager: AGENTS §3.1.
# Skript-Interpreter: MR-040/MR-042.
BLOCKED="apt apt-get brew pip pip3 pipx npm pnpm yarn npx corepack cargo rustup
gem conda go gofmt golangci-lint staticcheck perl ruby node deno bun"
# python, python3, python3.12 … — die Versions-Suffixe machen eine Liste
# unvollständig, darum als Muster.
BLOCKED_RE='^python[0-9]*(\.[0-9]+)*$'
PREFIXES="sudo env command exec nice time xargs eval"
SHELLS="bash sh zsh dash ksh"

in_set() {  # in_set <whitespace-getrennte-menge> <wort>
  local w
  for w in $1; do [ "$w" = "$2" ] && return 0; done
  return 1
}

# Ergebnis in der globalen STRIPPED: kein Subshell-Fork je Token — der Guard
# läuft vor JEDEM Bash-Call, Latenz zählt.
strip_quotes() {
  local s=$1
  while [ -n "$s" ]; do case $s in \"*|\'*) s=${s#?};; *) break;; esac; done
  while [ -n "$s" ]; do case $s in *\"|*\') s=${s%?};; *) break;; esac; done
  STRIPPED=$s
}

scan() {  # scan <cmd> <tiefe>; return 0 = BLOCK, 1 = ok
  local cmd=$1 depth=$2
  [ "$depth" -gt 3 ] && return 0          # zu tief verschachtelt -> fail-closed
  local s=$cmd
  s=${s//'&&'/$'\n'}; s=${s//'&'/$'\n'}; s=${s//'||'/$'\n'}; s=${s//'|'/$'\n'}
  s=${s//';'/$'\n'};  s=${s//\$\(/$'\n'}; s=${s//'`'/$'\n'}
  s=${s//'('/$'\n'};  s=${s//$'\r'/$'\n'}
  local seg head i j rest x
  local -a toks stoks
  while IFS= read -r seg; do
    read -ra toks <<< "$seg"
    [ "${#toks[@]}" -eq 0 ] && continue
    stoks=()
    for x in "${toks[@]}"; do strip_quotes "$x"; stoks+=("$STRIPPED"); done
    i=0
    while [ "$i" -lt "${#stoks[@]}" ]; do
      if [[ "${stoks[$i]}" =~ ^[A-Za-z_][A-Za-z0-9_]*= ]]; then i=$((i+1)); continue; fi
      in_set "$PREFIXES" "${stoks[$i]}" && { i=$((i+1)); continue; }
      case "${stoks[$i]}" in "{"|"}") i=$((i+1)); continue;; esac
      break
    done
    [ "$i" -ge "${#stoks[@]}" ] && continue
    head=${stoks[$i]}; head=${head##*/}   # /usr/bin/pip -> pip
    in_set "$BLOCKED" "$head" && return 0
    [[ "$head" =~ $BLOCKED_RE ]] && return 0
    if in_set "$SHELLS" "$head"; then
      # -c auch in Flag-Bündeln (-lc, -ec, -cx, …): bei sh/bash ist c das
      # einzige Single-Letter-Flag mit Kommando-String-Semantik.
      j=$((i+1))
      while [ "$j" -lt "${#stoks[@]}" ]; do
        if [[ "${stoks[$j]}" =~ ^-[a-z]*c[a-z]*$ ]]; then
          rest="${stoks[*]:$((j+1))}"
          scan "$rest" "$((depth+1))" && return 0
          break
        fi
        j=$((j+1))
      done
    fi
  done <<< "$s"
  return 1
}

# Eingabe per read-Builtin, nicht per `cat` (zweite Host-Abhängigkeit) und
# nicht per `$(</dev/stdin)`: Letzteres schlägt fehl, wo stdin keine
# nachlesbare Datei ist — unter `set -e` endet der Guard dann ohne Ausgabe,
# also fail-OPEN. `read -d ''` liest bis EOF und meldet dabei 1; das ist der
# Normalfall, kein Fehler.
# GRENZE des gewählten Wegs: `read -d ''` hält an einem NUL-Byte. Ein NUL im
# Kommando-Wert schneidet das JSON ab -> Extraktor-Zweifel -> Block; ein NUL
# zwischen zwei vollständigen Objekten verschluckt das zweite. JSON kennt kein
# rohes NUL, die Eingabe ist maschinell erzeugt — erreichbar ist das nicht.
IFS= read -r -d '' input || true

# Ohne awk keine Prüfung -> fail-closed (awk ist POSIX-Basis, AGENTS.md §3.1).
command -v awk >/dev/null 2>&1 || emit_block
[ -f "$extractor" ] || emit_block

set +e
cmd="$(printf '%s' "$input" | awk -f "$extractor")"
rc=$?
set -e
[ "$rc" -ne 0 ] && emit_block                # Parse-Zweifel -> fail-closed

scan "$cmd" 0 && emit_block
# Pass-Fall: keine Ausgabe — normale Permission-Prüfung übernimmt.
exit 0
