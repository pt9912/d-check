#!/usr/bin/env bash
# gate-consistency.sh — Meta-Gate gegen Harness-Lügen (Kurs-Modul 13):
#
#   (1) Jedes in AGENTS.md §4 bzw. harness/README.md §Sensors als
#       Tabellenzeile dokumentierte `make`-Target existiert im Makefile
#       (kein halluziniertes Gate).
#   (2) Jedes Makefile-Target ist in AGENTS.md §4 gelistet — AGENTS'
#       eigene Zusage: "Nur hier gelistete Targets existieren im
#       Makefile" (kein undokumentiertes Gate).
#   (3) DC-QA-03-Zusage des Netzlos-Gates: die modules-Liste der
#       .d-check.yml enthält alle Module außer `external` — sonst
#       verliert der `--network none`-Lauf still seine Beweis-Aussage
#       (Review R1 zu slice-008).
#
# Negativ-Test: vor der echten Prüfung läuft ein Selbsttest, in dem
# ein absichtlich dokumentiertes Phantom-Target das Gate nachweislich
# feuern lässt (analog verify-depguard-Idee).
set -euo pipefail

cd "$(dirname "$0")/.."

# Dokumentierte Targets: alle `make <name>`-Tokens in Tabellenzeilen
# (auch kombinierte Zellen wie "`make build` / `make run`"; eine
# Erwähnung in einer Tabellenzeile ist eine Existenz-Behauptung).
doc_targets() {
  grep -E '^\|' "$1" | grep -oE '`make [a-z][a-z0-9_-]*`' \
    | sed -E 's/`make ([a-z0-9_-]+)`/\1/' | sort -u
}

# Makefile-Targets: Regelzeilen `<name>:` am Zeilenanfang.
makefile_targets() {
  grep -oE '^[a-zA-Z][a-zA-Z0-9_-]*:' "$1" | tr -d ':' | sort -u
}

# Richtung 1: dokumentierte Targets müssen im Makefile existieren.
check_documented_exist() { # $1 Makefile, $2.. Doku-Dateien
  local mk="$1" fail=0 doc t
  shift
  local mk_targets
  mk_targets="$(makefile_targets "$mk")"
  for doc in "$@"; do
    while IFS= read -r t; do
      [ -z "$t" ] && continue
      if ! grep -qx "$t" <<<"$mk_targets"; then
        echo "gate-consistency: FAIL — $doc dokumentiert 'make $t', das Makefile kennt es nicht" >&2
        fail=1
      fi
    done <<<"$(doc_targets "$doc")"
  done
  return "$fail"
}

self_test() {
  local tmp
  tmp="$(mktemp -d)"
  printf '| `make phantom-target` | x |\n' > "$tmp/doc.md"
  printf 'echtes-target:\n\ttrue\n' > "$tmp/Makefile"
  if check_documented_exist "$tmp/Makefile" "$tmp/doc.md" 2>/dev/null; then
    echo "gate-consistency: Selbsttest FEHLGESCHLAGEN — Phantom-Target nicht erkannt" >&2
    rm -rf "$tmp"
    exit 2
  fi
  rm -rf "$tmp"
}

self_test
fail=0

# (1) Doku → Makefile
check_documented_exist Makefile AGENTS.md harness/README.md || fail=1

# (2) Makefile → AGENTS.md §4 (AGENTS' "Nur hier gelistete"-Zusage)
agents_targets="$(doc_targets AGENTS.md)"
while IFS= read -r t; do
  [ -z "$t" ] && continue
  if ! grep -qx "$t" <<<"$agents_targets"; then
    echo "gate-consistency: FAIL — Makefile-Target '$t' fehlt in AGENTS.md §4" >&2
    fail=1
  fi
done <<<"$(makefile_targets Makefile)"

# (3) DC-QA-03-Modulliste des Netzlos-Gates (.d-check.yml)
modules_line="$(grep -E '^modules:' .d-check.yml || true)"
for m in links anchors ids matrix; do
  if [[ "$modules_line" != *"$m"* ]]; then
    echo "gate-consistency: FAIL — .d-check.yml modules ohne '$m'; der Netzlos-Lauf beweist DC-QA-03 nur mit allen Modulen außer external" >&2
    fail=1
  fi
done
if [[ "$modules_line" == *external* ]]; then
  echo "gate-consistency: FAIL — .d-check.yml aktiviert external; das Netzlos-Gate darf kein Netz-Modul tragen (DC-QA-03)" >&2
  fail=1
fi

if [ "$fail" -ne 0 ]; then
  exit 1
fi
echo "gate-consistency ok: Doku ↔ Makefile konsistent, QA-03-Modulliste intakt (Selbsttest gefeuert)."
