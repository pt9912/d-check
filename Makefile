# d-check — Harness-Gates (Bootstrap-Stand, welle-01).
#
# `make gates` aggregiert nur real existierende Targets (Kurs-Modul 13).
# Implementierungs-Gates (lint/test/arch-check/coverage-gate) entstehen
# ab slice-003 — bis dahin werden sie hier bewusst NICHT behauptet.

.PHONY: help doc-check record-gates gates

# Der gates-Nachweis (record-gates) darf erst nach grünen Gates
# entstehen — unter `make -j` liefen Prerequisites parallel und der
# Nachweis entstünde trotz roter Gates (Review R1).
.NOTPARALLEL:

help: ## Targets anzeigen
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'

doc-check: ## interne Markdown-Linkziele existieren (Bootstrap-Sensor; Ablösung: slice-004)
	@bash tools/verify-doc-refs.sh

record-gates: ## Nachweis schreiben: Working-Tree-Hash (für den Stop-Hook)
	@bash tools/harness/record-gates.sh

# record-gates läuft als LETZTER Prerequisite — der Nachweis entsteht
# nur, wenn alle Gates grün sind (sonst bricht make vorher ab).
gates: doc-check record-gates ## alle inneren Gates (mandatory vor Handoff)
