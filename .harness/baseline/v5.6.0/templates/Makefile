# Gates. Doc-Gate via d-check: das Fragment `d-check.mk` ist tool-generiert
# (`d-check --print-mk`) und wird included — daher `--network none` + DCHECK_DIGEST-
# Re-Pin, kein handgepflegtes Recipe. Nur `docs-check` ist ein *behaupteter* Gate
# (die doc-*-Targets sind verfügbar, nicht behauptet, Modul 13); nur behauptete
# Targets in AGENTS.md / harness/README.md eintragen (keine halluzinierten Gates).
# Code-Gates (lint/test/build) repo-spezifisch — Vorbild: lab/example/<lang> im Kurs.
# `d-check.mk` wird NICHT mitgeliefert (tool-/versionsspezifisch) — beim Bootstrap
# aus deiner d-check erzeugen (Modul 2) und bei jedem d-check-Bump neu:
# `d-check --print-mk > d-check.mk`. Den Platzhalter unten durch den Digest deiner
# gepinnten d-check ersetzen (`sha256:…`); leer lassen nutzt den Tag aus d-check.mk
# (weniger reproduzierbar). `-include` bleibt still, bis das Fragment existiert.
DCHECK_DIGEST ?= <sha256-digest-deiner-gepinnten-d-check>
-include d-check.mk

.PHONY: help gates docs-check
help: ## Targets anzeigen
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-14s %s\n", $$1, $$2}'

# docs-check brückt das tool-generierte `doc-check` (reiner d-check) auf den
# kurs-kanonischen Namen — kein Eingriff ins Generat.
docs-check: doc-check ## Doku-Referenzen prüfen (links/anchors; ids/codepaths laut .d-check.yml)

gates: docs-check ## alle aktuell lauffähigen Gates (Code-Gates ergänzen)
