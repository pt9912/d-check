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

.PHONY: help deps compile lint test arch-check coverage-gate gate-consistency planning-check bench image-test semgrep versions build run doc-check trace record-gates gates ci fullbuild completeness-check trace-check adr-check hooks clean tidy

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

planning-check: ## Meta-Gate: Roadmap §Aktuelle Welle ↔ in-progress/slice-* (Planning-Drift-Schutz; slice-040).
	@bash tools/planning-consistency.sh

bench: build ## DC-QA-01-Benchmark: generiertes Fixture, N=3 Läufe, Median < 5 s (Spez §DC-QA-01.a).
	@bash tools/bench-fixture.sh

image-test: build ## DC-FA-DIST-001-Akzeptanzkriterien gegen das lokale Image (nativ vs. Container, :ro, Mount-Hinweis).
	@bash tools/image-test.sh

semgrep: ## Security-/Static-Analysis-Gate: gepinntes semgrep-Image + gepinntes, lokal gecachtes go/lang/security-Regelset, netzloser Scan (Bestandteil von gates; ADR-0010).
	@bash tools/semgrep.sh

versions: ## Reproduzierbarkeits-Pins ausgeben (Go, Lint, Basis-Image-Digests, semgrep, Runtime-Image-ID).
	@echo "GO_VERSION=$(GO_VERSION)"
	@echo "GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION)"
	@grep -E '^FROM ' Dockerfile | grep -v '^FROM deps' | sort -u
	@echo "semgrep-image=semgrep/semgrep:$$(sed -nE 's/.*SEMGREP_VERSION:-([^}]+)\}.*/\1/p' tools/semgrep.sh | head -1)@$$(sed -nE 's/.*SEMGREP_DIGEST:-([^}]+)\}.*/\1/p' tools/semgrep.sh | head -1)"
	@docker image inspect $(IMAGE):latest --format 'runtime-image={{.Id}}' 2>/dev/null \
	    || echo "runtime-image=(nicht gebaut — make build)"

# VERSION fließt ins OCI-Label org.opencontainers.image.version; die
# Release-Pipeline setzt sie aus dem Git-Tag (make ci VERSION=…).
VERSION ?= 0.0.0-dev

build: ## Runtime-Image bauen (distroless static, nonroot — ADR-0002).
	$(DOCKER_BUILD) --build-arg VERSION=$(VERSION) --target runtime -t $(IMAGE):latest .

run: build ## Smoke-Test: d-check prüft das eigene Repo (read-only).
	docker run --rm -v "$(CURDIR)":/repo:ro $(IMAGE):latest

# ---- docs gates --------------------------------------------------------------

# Gemeinsamer Dogfooding-Lauf: das LOKAL gebaute Image ($(IMAGE):latest, NICHT
# ein Release-Pin), read-only-Mount + --network none (DC-QA-03). Basis für
# doc-check/trace/doc-complete — der Pin existiert hier nicht (Produzent baut
# selbst), daher inline statt --print-mk-Fragment (das ist konsumenten-seitig).
DCHECK_RUN = docker run --rm --network none -v "$(CURDIR)":/repo:ro $(IMAGE):latest

# Vollständigkeits-Flag-Satz — EINE Quelle, geteilt von `doc-complete`
# (Konsumenten-/print-mk-Name) und `completeness-check` (Closure-Gate); so können
# die beiden Recipes nicht still divergieren (slice-055-R2-LOW). DC-FA-CLI-011.
COMPLETE_FLAGS = --trace --require-complete

# Dogfooding (MR-007, Selbstkonfiguration slice-007): d-check prüft die
# eigene Doku — Module links + anchors + ids + matrix über die gesamte
# Repo-Wurzel (.d-check.yml, scan.roots ".").
# Zugleich die automatisierte DC-QA-03-Messmethode (slice-008):
# read-only-Mount + --network none — alle Module außer external aktiv,
# der Lauf beweist Seiteneffektfreiheit und Netzlosigkeit.
doc-check: build ## Doku-Links, Anker, ID-Linkpflicht + Referenzmatrix via d-check selbst (Dogfooding, DC-FA-LINK/ANCH/ID/MTX; netzlos: DC-QA-03).
	$(DCHECK_RUN)

# Dogfooding-Render (kein Gate): die Requirements Traceability Matrix
# (--trace, DC-FA-CLI-009) über Lastenheft/ADR/Planning auf stdout — Komfort
# über das released Feature, read-only + --network none wie doc-check. Bewusst
# NICHT in `gates` (rein informativ, kein Pass/Fail).
trace: build ## Requirements Traceability Matrix via d-check selbst (Dogfooding, --trace; netzlos: DC-QA-03). DC-FA-CLI-009.
	$(DCHECK_RUN) --trace

# Dogfooding des opt-in Vollständigkeits-Modus (--require-complete, DC-FA-CLI-011):
# Exit 1 bei Requirements-Waisen, sonst 0 — gegen das lokal gebaute Image. Dies ist
# der konsumenten-/print-mk-seitige Name; den Closure-Bindepunkt bildet
# `make completeness-check` (dieselbe Mechanik, an fullbuild gehängt; ADR-0026).
doc-complete: build ## Vollständigkeits-Dogfood via d-check selbst (--trace --require-complete, Waise⇒Exit1; netzlos: DC-QA-03). DC-FA-CLI-011.
	$(DCHECK_RUN) $(COMPLETE_FLAGS)

# ---- harness -----------------------------------------------------------------

record-gates: ## Nachweis schreiben: Working-Tree-Hash (für den Stop-Hook).
	@bash tools/harness/record-gates.sh

# record-gates läuft als LETZTER Prerequisite — der Nachweis entsteht
# nur, wenn alle Gates grün sind (sonst bricht make vorher ab).
gates: doc-check lint test arch-check coverage-gate semgrep gate-consistency planning-check record-gates ## alle inneren Gates (mandatory vor Handoff).
	@echo "[gates] doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green"

# ci = gates + Image-Integrationstests — das Target, das die
# Release-Pipeline (slice-011) fährt. fullbuild = volle Closure vor
# Welle-Merge/Release inkl. Benchmark; bewusst NICHT Teil von gates
# (inner loop bleibt schnell).
ci: gates image-test ## CI-äquivalenter Lauf: gates + image-test (DC-FA-DIST-001).
	@echo "[ci] gates + image-test green"

fullbuild: ci bench completeness-check ## volle Closure: ci + bench + completeness-check; schließt mit dem Image-Hash (Reproduzierbarkeits-Bindung).
	@docker image inspect $(IMAGE):latest --format '[fullbuild] green — image-hash {{.Id}}'

# Requirements-Completeness-Gate (Policy slice-042/ADR-0017; Mechanik seit
# slice-055/ADR-0026 in-Produkt): failt bei Requirements-Waisen über den Flag
# `--trace --require-complete` (DC-FA-CLI-011) statt des abgelösten Skripts —
# die Waisen-IDs erscheinen als `WAISE`-Zeilen der --trace-Tabelle. Closure-
# Bindepunkt, an `fullbuild` gehängt, bewusst NICHT in `gates`/`ci` (GF erlaubt
# transiente Waisen). Dieselbe Mechanik wie `doc-complete`; completeness-check
# ist die Gate-/Closure-Rolle (d-check isst sein eigenes verteiltes Futter).
completeness-check: build ## Requirements-Completeness via in-Produkt-Flag (--trace --require-complete, Waise⇒Exit1); Closure-Gate (in fullbuild, NICHT gates/ci). ADR-0026 (löst ADR-0017-Skript ab).
	$(DCHECK_RUN) $(COMPLETE_FLAGS)

# ---- traceability ------------------------------------------------------------

# Traceability-Gate (ADR-0013): jede Commit-Message nennt eine
# DC-/ADR-/slice-ID. Bewusst NICHT Teil von `gates`/`ci` — anderer
# Bindepunkt (Commit-Zeit, nicht Arbeitsbaum-Inhalt): der CI-Workflow
# (.github/workflows/ci.yml) ruft es getrennt über den Commit-Range, der
# commit-msg-Hook prüft lokal. Eine Skript-Wahrheit (tools/trace-check.sh).
trace-check: ## Traceability-Gate: DC-/ADR-/slice-ID in Commits (Selbsttest + HEAD; RANGE=a..b für CI). ADR-0013.
	@bash tools/trace-check.sh $(if $(RANGE),--range $(RANGE),)

# adr-check fokussiert auf das Modul vcs: --enable vcs aktiviert es, und ALLE
# übrigen .d-check.yml-modules werden abgewählt — sonst über-feuerte das
# ADR-Immutable-Gate auf Nicht-ADR-Befunde des Arbeitsbaums (im STAGED-Hook auf
# ungestaged WIP), entgegen ADR-0024 „grün, sofern keine Accepted-ADR berührt".
# Die VCS_DISABLE-Liste spiegelt die .d-check.yml-modules; wächst die dort, hier
# nachziehen (das neue Modul liefe sonst in adr-check mit — kein Silent-Grün,
# aber Über-Feuern).
VCS_DISABLE := --disable links --disable anchors --disable ids --disable matrix \
    --disable codepaths --disable spans --disable hostpaths --disable versions
adr-check: build ## ADR-Immutable-Gate via Modul vcs (Image, dogfood, nur vcs): Accepted-ADRs nicht inhaltlich ändern (RANGE=a..b für CI, STAGED=1 für den Hook, sonst HEAD~1..HEAD). ADR-0024 (löst die Skript-Mechanik von ADR-0016 ab); ADR-0025 entfernt das Alt-Skript.
	$(DCHECK_RUN) --enable vcs $(VCS_DISABLE) $(if $(STAGED),--staged,--range $(if $(RANGE),$(RANGE),HEAD~1..HEAD))

hooks: ## git-Hooks installieren (core.hooksPath -> .githooks; commit-msg Traceability + pre-commit ADR-Immutable). ADR-0013/0016.
	@git config core.hooksPath .githooks
	@echo "[hooks] core.hooksPath=.githooks — commit-msg Traceability + pre-commit ADR-Immutable aktiv"

# ---- maintenance -------------------------------------------------------------

# go.mod/go.sum pflegen: die Go-Toolchain läuft in Docker (kein Host-Go,
# §3.1), schreibt als Host-User in ephemere Caches; `go mod tidy` nimmt die
# importierten Module auf (z. B. go-git für das Modul vcs, ADR-0024) und
# erneuert go.sum. Bewusster Akt am Dependency-Stand, kein Routine-Gate —
# go.sum ist der Reproduzierbarkeits-Anker, die deps-Stage prüft ihn beim
# Build (`-mod=readonly`).
tidy: ## go.mod/go.sum aufräumen (go mod tidy in Docker; Dependency-Pflege).
	docker run --rm -u "$$(id -u):$$(id -g)" \
	    -e HOME=/tmp -e GOCACHE=/tmp/gc -e GOMODCACHE=/tmp/gm \
	    -e GOTOOLCHAIN=local -e GOFLAGS=-mod=mod \
	    -v "$(CURDIR)":/src -w /src golang:$(GO_VERSION) \
	    go mod tidy

clean: ## Lokale Images entfernen.
	@-docker image rm \
	    $(IMAGE):latest $(IMAGE):deps $(IMAGE):compile \
	    $(IMAGE):lint $(IMAGE):test $(IMAGE):arch-check \
	    $(IMAGE):coverage 2>/dev/null || true
	@echo "[clean] images removed"
