# d-check — Doc-Referenz-Checker.
#
# Docker-only-Workflow (AGENTS.md §3.1; Vorbild:
# https://github.com/pt9912/u-boot): Build/Lint/
# Test/arch-check laufen über `docker build --target <stage>` in
# Containern. Der Host braucht nur Docker, GNU make, bash und git.
#
# `make gates` aggregiert nur real existierende Targets (Kurs-Modul 13);
# ci/fullbuild bauen darauf auf (harness/README.md §Sensors).

IMAGE                 ?= d-check
GO_VERSION            ?= 1.26.4
GOLANGCI_LINT_VERSION ?= v2.12.2

# `--progress=plain` für CI-taugliche BuildKit-Logs (u-boot-Konvention).
PROGRESS_FLAG :=
ifeq ($(CI),1)
PROGRESS_FLAG := --progress=plain
endif

# `--no-cache-filter <stage>`: erzwingt die Neu-Auswertung der
# Gate-Stage, ohne den deps-Cache zu verlieren — ein stale Layer-Hash
# darf keine rote Stage maskieren.
NO_CACHE_FILTER_TEST := --no-cache-filter test
NO_CACHE_FILTER_LINT := --no-cache-filter lint
NO_CACHE_FILTER_ARCH := --no-cache-filter arch-check
NO_CACHE_FILTER_COV  := --no-cache-filter coverage

# Kalibrierungs-Bindung (harness/README.md §Sensors): 93 % seit
# 2026-06-11 (Kalibrierung nach Test-Ausbau, Ist 95,1 %; zuvor Ramp
# 85 → 90 bei welle-03 done). Override: `make coverage-gate
# THRESHOLD=…`; Senkung nur per ADR (AGENTS.md §3.6).
THRESHOLD ?= 93

DOCKER_BUILD := docker build $(PROGRESS_FLAG) \
    --build-arg GO_VERSION=$(GO_VERSION) \
    --build-arg GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION)

.DEFAULT_GOAL := help

.PHONY: help deps compile lint test arch-check coverage-gate gate-consistency bench image-test versions build run doc-check record-gates gates ci fullbuild clean

# Der gates-Nachweis (record-gates) darf erst nach grünen Gates
# entstehen — unter `make -j` liefen Prerequisites parallel und der
# Nachweis entstünde trotz roter Gates (MR-005).
.NOTPARALLEL:

help: ## Targets anzeigen.
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# ---- inner-loop ------------------------------------------------------------

deps: ## Go-Modul-Auflösung (deps-Cache-Layer).
	$(DOCKER_BUILD) --target deps -t $(IMAGE):deps .

compile: ## Schnelles Compile-Feedback (ohne Tests/Lint).
	$(DOCKER_BUILD) --target compile -t $(IMAGE):compile .

lint: ## golangci-lint mit dem Projekt-Profil (AGENTS.md §3.2).
	$(DOCKER_BUILD) $(NO_CACHE_FILTER_LINT) --target lint -t $(IMAGE):lint .

test: ## `go test ./...` in Docker (Akzeptanzkriterien DC-FA-*).
	$(DOCKER_BUILD) $(NO_CACHE_FILTER_TEST) --target test -t $(IMAGE):test .

arch-check: ## ADR-0005 — Import-Regeln des Hexagon-Schnitts (DC-QA-03).
	$(DOCKER_BUILD) $(NO_CACHE_FILTER_ARCH) --target arch-check -t $(IMAGE):arch-check .

coverage-gate: ## Coverage-Schwelle (Kalibrierungs-Bindung: 93 %, Historie in harness/README §Sensors).
	$(DOCKER_BUILD) $(NO_CACHE_FILTER_COV) \
	    --build-arg COVERAGE_THRESHOLD=$(THRESHOLD) \
	    --target coverage -t $(IMAGE):coverage .

gate-consistency: ## Meta-Gate: dokumentierte Targets ↔ Makefile, QA-03-Modulliste (Harness-Lügen-Schutz).
	@bash tools/gate-consistency.sh

bench: build ## DC-QA-01-Benchmark: generiertes Fixture, N=3 Läufe, Median < 5 s (Spez §DC-QA-01.a).
	@bash tools/bench-fixture.sh

image-test: build ## DC-FA-DIST-001-Akzeptanzkriterien gegen das lokale Image (nativ vs. Container, :ro, Mount-Hinweis).
	@bash tools/image-test.sh

versions: ## Reproduzierbarkeits-Pins ausgeben (Go, Lint, Basis-Images, Runtime-Image-ID).
	@echo "GO_VERSION=$(GO_VERSION)"
	@echo "GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION)"
	@grep -E '^FROM ' Dockerfile | grep -v '^FROM deps' | sort -u
	@docker image inspect $(IMAGE):latest --format 'runtime-image={{.Id}}' 2>/dev/null \
	    || echo "runtime-image=(nicht gebaut — make build)"

build: ## Runtime-Image bauen (distroless static, nonroot — ADR-0002).
	$(DOCKER_BUILD) --target runtime -t $(IMAGE):latest .

run: build ## Smoke-Test: d-check prüft das eigene Repo (read-only).
	docker run --rm -v "$(CURDIR)":/repo:ro $(IMAGE):latest

# ---- docs gates --------------------------------------------------------------

# Dogfooding (MR-007, Selbstkonfiguration slice-007): d-check prüft die
# eigene Doku — Module links + anchors + ids + matrix über die gesamte
# Repo-Wurzel (.d-check.yml, scan.roots ".").
# Zugleich die automatisierte DC-QA-03-Messmethode (slice-008):
# read-only-Mount + --network none — alle Module außer external aktiv,
# der Lauf beweist Seiteneffektfreiheit und Netzlosigkeit.
doc-check: build ## Doku-Links, Anker, ID-Linkpflicht + Referenzmatrix via d-check selbst (Dogfooding, DC-FA-LINK/ANCH/ID/MTX; netzlos: DC-QA-03).
	docker run --rm --network none -v "$(CURDIR)":/repo:ro $(IMAGE):latest

# ---- harness -----------------------------------------------------------------

record-gates: ## Nachweis schreiben: Working-Tree-Hash (für den Stop-Hook).
	@bash tools/harness/record-gates.sh

# record-gates läuft als LETZTER Prerequisite — der Nachweis entsteht
# nur, wenn alle Gates grün sind (sonst bricht make vorher ab).
gates: doc-check lint test arch-check coverage-gate gate-consistency record-gates ## alle inneren Gates (mandatory vor Handoff).
	@echo "[gates] doc-check + lint + test + arch-check + coverage-gate + gate-consistency green"

# ci = gates + Image-Integrationstests — das Target, das die
# Release-Pipeline (slice-011) fährt. fullbuild = volle Closure vor
# Welle-Merge/Release inkl. Benchmark; bewusst NICHT Teil von gates
# (inner loop bleibt schnell).
ci: gates image-test ## CI-äquivalenter Lauf: gates + image-test (DC-FA-DIST-001).
	@echo "[ci] gates + image-test green"

fullbuild: ci bench ## volle Closure: ci + bench; schließt mit dem Image-Hash (Reproduzierbarkeits-Bindung).
	@docker image inspect $(IMAGE):latest --format '[fullbuild] green — image-hash {{.Id}}'

# ---- maintenance -------------------------------------------------------------

clean: ## Lokale Images entfernen.
	@-docker image rm \
	    $(IMAGE):latest $(IMAGE):deps $(IMAGE):compile \
	    $(IMAGE):lint $(IMAGE):test $(IMAGE):arch-check \
	    $(IMAGE):coverage 2>/dev/null || true
	@echo "[clean] images removed"
