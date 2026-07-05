#!/usr/bin/env bash
# gate-consistency.sh — Meta-Gate gegen Harness-Lügen, **Rest-Prüfung (3)**:
# die DC-QA-03-Modulliste des netzlosen doc-check-Gates.
#
# Der cross-repo-driftende **Kern** (Prüfung 1+2: dokumentierte `make X` ↔
# Makefile-Regeln, in beide Richtungen) ist seit ADR-0031/slice-063 das opt-in
# Modul `targets` — `make gate-consistency` dogfoodet es via Image (verteilbar,
# kein kopiertes Skript mehr; der phantom-Target-Selbsttest lebt als
# Modul-Akzeptanztest in rules/targets_test.go). Hier bleibt nur die
# **repo-spezifische** Modul-Listen-Selbstkonsistenz (kein cross-repo-Kopie-
# Drift, daher nicht d-check-mechanisiert):
#
#   (3) DC-QA-03: die modules-Liste der .d-check.yml führt alle netzlosen
#       Doku-Module (links/anchors/ids/matrix/codepaths) und **nicht** external
#       — sonst verliert der `--network none`-Lauf still seine Beweis-Aussage.
set -euo pipefail
cd "$(dirname "$0")/.."

fail=0

# Format-Annahme: einzeilige Flow-Style-Liste (`modules: [a, b]`). Bei Umstellung
# auf YAML-Listenform wird dieser Check laut rot (fail-closed) — dann den Parser
# hier mitziehen.
modules_line="$(grep -E '^modules:' .d-check.yml || true)"
for m in links anchors ids matrix codepaths; do
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
echo "gate-consistency ok: QA-03-Modulliste intakt (Doku↔Makefile-Kern via Modul targets, ADR-0031)."
