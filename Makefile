# d-check — Doc-Referenz-Checker.
#
# Docker-only-Workflow (AGENTS.md §3.1; Vorbild:
# https://github.com/pt9912/u-boot): Build/Lint/Test
# laufen über `docker build --target <stage>` in Containern; arch-check
# läuft über das digest-gepinnte a-check-Image (a-check.mk, ADR-0029).
# Der Host braucht nur Docker, GNU make, bash und git.
#
# `make gates` aggregiert nur real existierende Targets (Kurs-Modul 13);
# ci/fullbuild bauen darauf auf (harness/README.md §Sensors).

IMAGE                 ?= d-check
GO_VERSION            ?= 1.27.1
GOLANGCI_LINT_VERSION ?= v2.13.2

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
NO_CACHE_FILTER_COV  := --no-cache-filter coverage

# Architektur-Gate via Schwester-Tool a-check (ADR-0029): das Fragment
# liefert A_CHECK_IMAGE (digest-gepinnt) + das Basis-Target; das
# arch-check-Target unten delegiert dorthin und bleibt in DIESER Datei
# (gate-consistency parst nur das Makefile, keine includes).
include a-check.mk

# Kalibrierungs-Bindung (harness/README.md §Sensors): 93 % seit
# 2026-06-11 (Kalibrierung nach Test-Ausbau, Ist 95,1 %; zuvor Ramp
# 85 → 90 bei welle-03 done). Override: `make coverage-gate
# THRESHOLD=…`; Senkung nur per ADR (AGENTS.md §3.6).
THRESHOLD ?= 93

DOCKER_BUILD := docker build $(PROGRESS_FLAG) \
    --build-arg GO_VERSION=$(GO_VERSION) \
    --build-arg GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION)

.DEFAULT_GOAL := help

.PHONY: nightly-state freshness-semgrep semgrep-digest freshness-a-check a-check-digest help deps compile lint test arch-check baseline-verify baseline-freshness workflow-pins freshness-go freshness-golangci runtime-base-digest go-base-digest lint-base-digest checkout-pin-freshness login-pin-freshness coverage-gate gate-consistency planning-check verify-closure-notes bench image-test semgrep versions build run doc-check trace record-gates guard-probe gates ci fullbuild completeness-check trace-check adr-check hooks clean tidy image-scan freshness-trivy trivy-digest archive-wave-test archive-wave

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

arch-check: a-check ## Import-Regeln R1–R6 (ADR-0005/ADR-0012) via a-check-Image (.a-check.yml; netzlos, read-only — DC-QA-03). ADR-0029 (löst tools/arch-check.sh ab).

coverage-gate: ## Coverage-Schwelle (Kalibrierungs-Bindung: 93 %, Historie in harness/README §Sensors).
	$(DOCKER_BUILD) $(NO_CACHE_FILTER_COV) \
	    --build-arg COVERAGE_THRESHOLD=$(THRESHOLD) \
	    --target coverage -t $(IMAGE):coverage .

gate-consistency: build ## Meta-Gate: Deklarations-Konsistenz Doku↔Makefile via Modul targets (Image, dogfood). ADR-0031/ADR-0032 (Skript voll abgelöst; DC-QA-03-Modulliste jetzt Go-Test in `make test`).
	$(DCHECK_RUN) --enable targets $(FOCUS_DISABLE)

planning-check: build ## Meta-Gate via Modul planning (Image, dogfood): Roadmap §Offene Wellen (Ruhe-Marker) ↔ in-progress/slice-* (Planning-Drift-Schutz, in gates). ADR-0028 (löst die Skript-Mechanik von slice-040 ab).
	$(DCHECK_RUN) --enable planning $(FOCUS_DISABLE)

bench: build ## DC-QA-01-Benchmark: generiertes Fixture, N=3 Läufe, Median < 5 s (Spez §DC-QA-01.a).
	@bash tools/bench-fixture.sh

image-test: build ## DC-FA-DIST-001-Akzeptanzkriterien gegen das lokale Image (nativ vs. Container, :ro, Mount-Hinweis).
	@bash tools/image-test.sh

image-scan: ## CVE-Scan gegen die PUBLIZIERTEN Images (Netz, NICHT in gates, Trivy digest-gepinnt; ADR-0066). Exit 1 = behebbare CRITICAL/HIGH, 2 = Scan gescheitert.
	@bash tools/image-scan.sh

semgrep: ## Security-/Static-Analysis-Gate: gepinntes semgrep-Image + gepinntes, lokal gecachtes go/lang/security-Regelset, netzloser Scan (Bestandteil von gates; ADR-0010).
	@bash tools/semgrep.sh

versions: ## Reproduzierbarkeits-Pins ausgeben (Go, Lint, Basis-Image-Digests, semgrep, a-check, Runtime-Image-ID).
	@echo "GO_VERSION=$(GO_VERSION)"
	@echo "GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION)"
	@grep -E '^FROM ' Dockerfile | grep -v '^FROM deps' | sort -u
	@echo "semgrep-image=semgrep/semgrep:$$(sed -nE 's/.*SEMGREP_VERSION:-([^}]+)\}.*/\1/p' tools/semgrep.sh | head -1)@$$(sed -nE 's/.*SEMGREP_DIGEST:-([^}]+)\}.*/\1/p' tools/semgrep.sh | head -1)"
	@echo "a-check-image=$(A_CHECK_IMAGE)"
	@echo "trivy-image=aquasec/trivy:$$(sed -nE 's/.*TRIVY_VERSION:-([^}]+)\}.*/\1/p' tools/image-scan.sh | head -1)@$$(sed -nE 's/.*TRIVY_DIGEST:-([^}]+)\}.*/\1/p' tools/image-scan.sh | head -1)"
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
# Wie DCHECK_RUN, aber mit -i (stdin offen) — für den commit-msg-Hook, der die
# Pending-Message über stdin an `--commit-msg -` pipet (Modul commits, ADR-0027).
DCHECK_RUN_I = docker run --rm -i --network none -v "$(CURDIR)":/repo:ro $(IMAGE):latest

# Vollständigkeits-Flag-Satz — EINE Quelle, geteilt von `doc-complete`
# (Konsumenten-/print-mk-Name) und `completeness-check` (Closure-Gate); so können
# die beiden Recipes nicht still divergieren (seit slice-055). DC-FA-CLI-011.
COMPLETE_FLAGS = --trace --require-complete

# Dogfooding (MR-007, Selbstkonfiguration slice-007): d-check prüft die
# eigene Doku — die Modulliste führt die .d-check.yml, der Scan-Bereich ihr
# scan-Block (Repo-Wurzel).
# Zugleich die automatisierte DC-QA-03-Messmethode (slice-008):
# read-only-Mount + --network none — alle Module außer external aktiv,
# der Lauf beweist Seiteneffektfreiheit und Netzlosigkeit.
doc-check: build ## Links, Anker, Kennungs-Linkpflicht, Referenzmatrix, Inline-Code-Pfade, Spans, Host-Pfade, Versions-Pins, Abschnitts-Invarianten und Diagramm-Kennungen via d-check selbst (Dogfooding; netzlos: DC-QA-03).
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

# Der Tool-Call-Wächter ist werkzeug-lokal und hat kein Gate (AGENTS.md §3.1).
# Ohne wiederholbare Proben wäre seine Zusage eine Erinnerung; dieses Target
# macht sie nachfahrbar. Bewusst NICHT in gates: der Wächter ist keine
# Repo-Invariante, sondern eine Werkzeug-Einstellung.
guard-probe: ## Tool-Call-Wächter gegen seine Proben fahren (werkzeug-lokal, NICHT in gates).
	@bash tools/harness/guard-probe.sh

# Der vendorte Baseline-Bestand ist die netzlose Leseform des adoptierten
# Regelwerks (MR-011-Kette). Zwei Fragen, zwei Targets, zwei Fehlerpolitiken —
# die Trennung ist tragend und steht deshalb im Namen:
#   baseline-verify    "ist der committete Bestand unversehrt?" — netzlos,
#                      fail-closed, IN gates. DREI Fragen: sha256sum -c
#                      erkennt geänderte/gelöschte Dateien, die Manifest-Deckung
#                      zusätzlich EINGELEGTE (eine untermengige, in sich
#                      konsistente SHA256SUMS passierte sonst grün), und die
#                      Aufloesung der Aliase unter .claude/rules/ — ein
#                      Symlink bindet denselben Pin, steht aber in keiner
#                      Manifest-Zeile (MR-055).
#   baseline-freshness "ist der gepinnte Stand noch aktuell und authentisch?" —
#                      Netz, fail-open (Ausfall ⇒ SKIP je Teil), NICHT in gates.
#                      Der netzlose innere Lauf ist eine Eigenschaft DIESES
#                      REPOS; DC-QA-03 gilt dem Produkt und ist hier nur über
#                      doc-check berührt, dessen Container-Lauf ihre Messmethode
#                      ist. Meldet nur; die Pin-Hebung bleibt ein bewusster Akt
#                      der MR-011-Kette.
# Wächter zu AGENTS.md §3.9. Er prüft die FORM des Pins (voller Commit-SHA plus
# Tag-Kommentar), nicht seine Gültigkeit — ob der SHA existiert und zum Tag
# passt, wäre Netz und gehört zur Freshness-Familie. Netzlos, fail-closed auch
# bei leerer Prüfmenge, deshalb IN gates.
workflow-pins: build ## uses:-Einträge der Workflows via Modul workflows (Image, dogfood): voller SHA + Tag-Kommentar; lokale Referenz existiert und bekommt die Rechte ihres Jobs (netzlos, in gates). AGENTS.md §3.9, ADR-0072.
	$(DCHECK_RUN) --enable workflows $(FOCUS_DISABLE)

# Netzlos, fail-closed auch bei leerer Kandidatenmenge -- aber bewusst NOCH
# NICHT in gates: eine neue Modul-Klasse startet als eigenstaendiger Fokus-Lauf
# (dieselbe Vorsicht wie bei trace-check/commits), Aufnahme in gates ist eine
# spaetere, eigene Entscheidung.
review-coverage: build ## Review-Report-Deckung via Modul reviews (Image, dogfood): jede DoD-Zusage "unabhängiger Review" braucht einen passenden Report unter docs/reviews/ (netzlos, NICHT in gates). ADR-0081.
	$(DCHECK_RUN) --enable reviews $(FOCUS_DISABLE)

baseline-probe: ## Faehrt die Alias-Aufloesung von baseline-verify gegen ihre Proben (neun Faelle, netzlos, NICHT in gates). MR-055.
	@bash tools/harness/fetch-baseline-cache.sh --selftest

baseline-verify: ## Vendorte Baseline gegen SHA256SUMS: Integrität + Manifest-Deckung (netzlos, in gates). MR-011-Kette.
	@bash tools/harness/fetch-baseline-cache.sh --verify

# Zwei Toolchain-Achsen, ein Sensor. Netz, fail-open, NICHT in gates -- wie
# baseline-freshness und aus demselben Grund. Sie MELDEN; die Hebung bleibt ein
# bewusster Akt, und beim Go-Bump zieht das golangci-Pendant mit.
freshness-go: ## Neuere stabile Go-Version als GO_VERSION melden (Netz, Quelle go.dev, NICHT in gates, fail-open).
	@NAME='go' PINNED='$(GO_VERSION)' \
	  ADVICE='GO_VERSION (Makefile) heben, Dockerfile-Digest nachziehen; das golangci-Pendant zieht mit.' \
	  bash tools/harness/pin-freshness.sh --godev

freshness-golangci: ## Neueren golangci-lint-Release als GOLANGCI_LINT_VERSION melden (Netz, NICHT in gates, fail-open).
	@NAME='golangci-lint' PINNED='$(GOLANGCI_LINT_VERSION)' \
	  ADVICE='GOLANGCI_LINT_VERSION (Makefile) heben und den Dockerfile-lint-Digest nachziehen.' \
	  bash tools/harness/pin-freshness.sh --github golangci/golangci-lint

# Die dritte Achsen-FORM beantwortet eine ANDERE Frage als die zwei oben: nicht
# „gibt es einen neueren Tag", sondern „traegt derselbe Tag inzwischen einen
# anderen Digest". Sie gilt ALLEN drei Basis-Images: die Tag-Achsen wachen die
# VERSION, nicht den Bau — `freshness-go` meldet `ok`, waehrend derselbe
# golang-Tag laengst neu gebaut ist. Fuer das Runtime-Image ist sie ausserdem
# die EINZIGE Handhabe, weil sein Tag `nonroot` keine Version fuehrt.
#
# Referenz UND Pin kommen aus DERSELBEN FROM-Zeile — zwei Quellen waeren zwei
# Spiegel, und ein Wechsel des Images liesse den einen stehen.
image-digest-axis = @set -- $$(awk -v p='$(1)' '$$1=="FROM" && index($$2,p)==1 { d=$$2; sub(/.*@/,"",d); sub(/@.*/,"",$$2); print $$2, d; exit }' Dockerfile \
	  | sed -e 's|$${GO_VERSION}|$(GO_VERSION)|' -e 's|$${GOLANGCI_LINT_VERSION}|$(GOLANGCI_LINT_VERSION)|'); \
	NAME="$$1" PINNED="$$2" \
	  ADVICE='Dockerfile-Digest der Stage nachziehen (ADR-0011); make versions zeigt den Stand.' \
	  bash tools/harness/pin-freshness.sh --digest "$$1"

runtime-base-digest: ## Neueren Digest fuer denselben Runtime-Basis-Tag melden (Netz, NICHT in gates, fail-open).
	$(call image-digest-axis,gcr.io/distroless)

go-base-digest: ## Neueren Digest fuer denselben Go-Basis-Tag melden (Netz, NICHT in gates, fail-open).
	$(call image-digest-axis,golang:)

lint-base-digest: ## Neueren Digest fuer denselben Lint-Basis-Tag melden (Netz, NICHT in gates, fail-open).
	$(call image-digest-axis,golangci/)

# Die Action-Pins brauchen KEINE neue Quellen-Form: der Pin ist ein SHA, aber
# der Tag-Kommentar daneben traegt den Release-Tag — genau die Groesse, die
# `--github` vergleicht. Ein alter Action-Pin ist zudem sehr wohl ein
# Sicherheitsthema: ein SHA-Pin schliesst das Umhaengen eines Tags aus und macht
# zugleich blind fuer die Behebung.
action-pin-axis = @NAME='$(1)' \
	PINNED="$$(grep -h 'uses: $(1)@' .github/workflows/*.yml | head -1 | awk '{print $$NF}')" \
	ADVICE='SHA + Tag-Kommentar in .github/workflows/ heben (AGENTS.md §3.9).' \
	bash tools/harness/pin-freshness.sh --github $(1)

checkout-pin-freshness: ## Neueren actions/checkout-Release als den Tag-Kommentar melden (Netz, NICHT in gates, fail-open).
	$(call action-pin-axis,actions/checkout)

login-pin-freshness: ## Neueren docker/login-action-Release als den Tag-Kommentar melden (Netz, NICHT in gates, fail-open).
	$(call action-pin-axis,docker/login-action)

hubdesc-pin-freshness: ## Neueren peter-evans/dockerhub-description-Release als den Tag-Kommentar melden (Netz, NICHT in gates, fail-open).
	$(call action-pin-axis,peter-evans/dockerhub-description)

# Die zwei uebrigen gepinnten Fremd-Images. Sie stehen NICHT im Dockerfile —
# semgrep im Gate-Skript, a-check im include-Fragment —, deshalb je ein eigener
# Extraktor statt des FROM-Musters oben. Beide tragen Tag UND Digest, also
# stehen ihnen beide Fragen offen: der Tag fragt nach der Version, der Digest
# nach dem Bau desselben Tags.
freshness-semgrep: ## Neueren semgrep-Release als SEMGREP_VERSION melden (Netz, NICHT in gates, fail-open).
	@NAME='semgrep' \
	  PINNED="$$(sed -nE 's/.*SEMGREP_VERSION:-([^}]+)\}.*/\1/p' tools/semgrep.sh | head -1)" \
	  ADVICE='SEMGREP_VERSION + SEMGREP_DIGEST in tools/semgrep.sh heben (ADR-0010/ADR-0011).' \
	  bash tools/harness/pin-freshness.sh --github semgrep/semgrep

freshness-trivy: ## Neueren Trivy-Release als TRIVY_VERSION melden (Netz, NICHT in gates, fail-open).
	@NAME=aquasec/trivy \
	  PINNED="$$(sed -nE 's/.*TRIVY_VERSION:-([^}]+)\}.*/\1/p' tools/image-scan.sh | head -1)" \
	  ADVICE='TRIVY_VERSION + TRIVY_DIGEST in tools/image-scan.sh heben (ADR-0011/ADR-0066).' \
	  bash tools/harness/pin-freshness.sh --github aquasecurity/trivy

trivy-digest: ## Neuen Bau desselben Trivy-Tags melden (Netz, NICHT in gates, fail-open).
	@NAME="aquasec/trivy:$$(sed -nE 's/.*TRIVY_VERSION:-([^}]+)\}.*/\1/p' tools/image-scan.sh | head -1)" \
	  PINNED="$$(sed -nE 's/.*TRIVY_DIGEST:-([^}]+)\}.*/\1/p' tools/image-scan.sh | head -1)" \
	  ADVICE='TRIVY_DIGEST in tools/image-scan.sh nachziehen (ADR-0011).' \
	  bash tools/harness/pin-freshness.sh --digest "aquasec/trivy:$$(sed -nE 's/.*TRIVY_VERSION:-([^}]+)\}.*/\1/p' tools/image-scan.sh | head -1)"
semgrep-digest: ## Neueren Digest fuer denselben semgrep-Tag melden (Netz, NICHT in gates, fail-open).
	@NAME="semgrep/semgrep:$$(sed -nE 's/.*SEMGREP_VERSION:-([^}]+)\}.*/\1/p' tools/semgrep.sh | head -1)" \
	  PINNED="$$(sed -nE 's/.*SEMGREP_DIGEST:-([^}]+)\}.*/\1/p' tools/semgrep.sh | head -1)" \
	  ADVICE='SEMGREP_DIGEST in tools/semgrep.sh nachziehen (ADR-0011).' \
	  bash tools/harness/pin-freshness.sh --digest "semgrep/semgrep:$$(sed -nE 's/.*SEMGREP_VERSION:-([^}]+)\}.*/\1/p' tools/semgrep.sh | head -1)"

freshness-a-check: ## Neueren a-check-Release als A_CHECK_VERSION melden (Netz, NICHT in gates, fail-open).
	@NAME='a-check' PINNED='$(A_CHECK_VERSION)' \
	  ADVICE='A_CHECK_VERSION + Digest in a-check.mk heben; Fragment per --print-mk neu erzeugen (ADR-0029).' \
	  bash tools/harness/pin-freshness.sh --github pt9912/a-check

a-check-digest: ## Neueren Digest fuer denselben a-check-Tag melden (Netz, NICHT in gates, fail-open).
	@NAME='ghcr.io/pt9912/a-check:$(A_CHECK_VERSION)' \
	  PINNED="$$(printf '%s' '$(A_CHECK_IMAGE)' | sed 's/.*@//')" \
	  ADVICE='Digest in a-check.mk nachziehen (ADR-0011).' \
	  bash tools/harness/pin-freshness.sh --digest 'ghcr.io/pt9912/a-check:$(A_CHECK_VERSION)'

# Lese-Schritt, kein Gate: er sagt, ob der Nachtlauf gelesen werden muss, und
# haengt an dem Moment, an dem ohnehin jemand hinsieht (Slice-Planung, dritte
# Vorpruefung nach MR-053). Netz, fail-open, immer Exit 0 — der Ausgang steht
# in der AUSGABE, damit ein Exit-Code ihn nicht verdecken kann.
nightly-state: ## Ausgang des juengsten Nachtlaufs lesen (Netz, fail-open, NICHT in gates; Vorpruefung der Slice-Planung).
	@bash tools/harness/nightly-state.sh

baseline-freshness: ## Upstream-Audit des Baseline-Pins: neuerer Release-Tag (Currency) + Content-Drift am gepinnten Tag (Netz, NICHT in gates, fail-open). MR-011-Kette.
	@bash tools/harness/fetch-baseline-cache.sh --check-latest

# record-gates läuft als LETZTER Prerequisite — der Nachweis entsteht
# nur, wenn alle Gates grün sind (sonst bricht make vorher ab).
gates: baseline-verify workflow-pins doc-check lint test arch-check coverage-gate semgrep gate-consistency planning-check record-gates ## alle inneren Gates (mandatory vor Handoff).
	@echo "[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green"

# ci = gates + Image-Integrationstests — das Target, das die
# Release-Pipeline (slice-011) fährt. fullbuild = volle Closure vor
# Welle-Merge/Release inkl. Benchmark; bewusst NICHT Teil von gates
# (inner loop bleibt schnell).
ci: gates image-test ## CI-äquivalenter Lauf: gates + image-test (DC-FA-DIST-001).
	@echo "[ci] gates + image-test green"

fullbuild: ci bench completeness-check verify-closure-notes ## volle Closure: ci + bench + completeness-check + verify-closure-notes; schließt mit dem Image-Hash (Reproduzierbarkeits-Bindung).
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

# Zweiter Closure-Bindepunkt neben completeness-check: die STRUKTUR der
# Closure-Notizen im done/-Bestand (Modul planning, DC-FA-PLAN-001). Fährt ein
# eigenes Prüf-Profil über --config (DC-FA-CLI-012) — stünde der closure-Block
# in der konventionellen .d-check.yml, liefe er in `gates` mit, und eine
# Closure-Frage gehört nicht in den Inner-Loop (ADR-0048).
#
# DAS DRITTE MODUL MACHT DEN BINDEPUNKT SELBSTGENUEGSAM (ADR-0077):
# planning und structure lesen den BEREINIGTEN Abschnitts-Text; ein vergessener
# Schluss-Fence verschluckt alles dahinter, und ihre Zusagen werden STILL wahr.
# spans meldet das — findet dabei aber NICHTS, was doc-check nicht schon beim
# Commit faende (gemessen; die Scan-Menge hier ist eine Teilmenge). Gekauft ist
# die Unabhaengigkeit von einem fremden Profil, nicht neue Deckung.
verify-closure-notes: build ## Struktur des done/-Bestands: Closure-Notizen + Register-Deckung (Modul planning) UND Abschnitts-Invarianten (Modul structure), plus die Span-/Fence-Artefakte (Modul spans) UND Review-Report-Deckung (Modul reviews) -- via eigenes --config-Profil; Closure-Gate (in fullbuild, NICHT gates/ci). ADR-0048/ADR-0049/ADR-0077/ADR-0081/ADR-0082.
	$(DCHECK_RUN) --config .d-check.closure.yml --enable planning --enable structure --enable spans --enable reviews

# ---- traceability ------------------------------------------------------------

# Traceability-Gate (ADR-0013 Policy; Mechanik seit slice-056/ADR-0027 in-Produkt
# über das Modul commits statt tools/trace-check.sh). Jede Commit-Message nennt
# eine DC-/ADR-/MR-/slice-ID. Bewusst NICHT Teil von `gates`/`ci` — anderer
# Bindepunkt (Commit-Zeit, nicht Arbeitsbaum-Inhalt). Zwei Modi, eine Wahrheit
# (das Modul commits im Image):
#   MSGFILE=<datei>  commit-msg-Hook: Pending-Message via stdin an --commit-msg -
#   sonst            Range-Modus: --enable commits (fokussiert) --range (Default HEAD~1..HEAD)
# Der CI-Workflow (.github/workflows/ci.yml) ruft den Range-Modus, der Hook den
# Message-Modus. Fokus-Disable wie adr-check (nur commits läuft).
trace-check: build ## Traceability-Gate via Modul commits (Image, dogfood): DC-/ADR-/MR-/slice-ID in Commit-Messages (RANGE=a..b für CI, MSGFILE=… für den Hook, sonst HEAD~1..HEAD). ADR-0027 (löst die Skript-Mechanik von ADR-0013 ab).
	@$(if $(MSGFILE),$(DCHECK_RUN_I) --commit-msg - < $(MSGFILE),$(DCHECK_RUN) --enable commits $(FOCUS_DISABLE) --range $(if $(RANGE),$(RANGE),HEAD~1..HEAD))

# FOCUS_DISABLE wählt ALLE .d-check.yml-modules ab, sodass ein fokussiertes Gate
# nur sein eines opt-in-Modul laufen lässt (adr-check nur vcs, trace-check nur
# commits) — sonst über-feuerten die Datei-Module auf den Arbeitsbaum-Inhalt (im
# STAGED-Hook auf ungestaged WIP), entgegen ADR-0024 „grün, sofern keine
# Accepted-ADR berührt". Spiegelt die .d-check.yml-modules-Liste; wächst die dort,
# hier nachziehen (ein neues Default-Modul liefe sonst mit — kein Silent-Grün,
# aber Über-Feuern).
FOCUS_DISABLE := --disable links --disable anchors --disable ids --disable matrix \
    --disable codepaths --disable spans --disable hostpaths --disable versions \
    --disable structure --disable diagrams --disable citations
adr-check: build ## ADR-Immutable-Gate via Modul vcs (Image, dogfood, nur vcs): Accepted-ADRs nicht inhaltlich ändern (RANGE=a..b für CI, STAGED=1 für den Hook, sonst HEAD~1..HEAD). ADR-0024 (löst die Skript-Mechanik von ADR-0016 ab); ADR-0025 entfernt das Alt-Skript.
	$(DCHECK_RUN) --enable vcs $(FOCUS_DISABLE) $(if $(STAGED),--staged,--range $(if $(RANGE),$(RANGE),HEAD~1..HEAD))

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
	    $(IMAGE):lint $(IMAGE):test \
	    $(IMAGE):coverage 2>/dev/null || true
	@echo "[clean] images removed"

# archive-wave (Baseline-Regelwerk modul-06-roadmap.md §Wellen-Closure-
# Prozedur, Schritt 4): eigenstaendiges Werkzeug unter tools/archive-wave/,
# eigenes go.mod, eigenes Dockerfile, eigenes Makefile -- portabel fuer
# jedes Repo mit demselben Planning-Layout. Delegiert an das
# lokale Makefile statt den Docker-Aufruf hier zu duplizieren -- eine
# Quelle fuer den Mount/UID-Umgang, kein Drift-Risiko zwischen zwei
# Kopien. Sicherer Default: ohne APPLY=1 wird NICHTS geschrieben und der
# Mount ist read-only (das lokale Makefile schaltet ihn erst bei
# APPLY=1 beschreibbar).
archive-wave-test: ## archive-wave-Testsuite (eigenes go.mod, nicht Teil von `make test`).
	$(MAKE) -C tools/archive-wave test GO_VERSION=$(GO_VERSION) PROGRESS_FLAG='$(PROGRESS_FLAG)'

archive-wave: ## Welle oder wellenlosen Slice archivieren: make archive-wave WELLE=welle-NN|SLICE=slice-NNN [APPLY=1].
	$(MAKE) -C tools/archive-wave run WELLE=$(WELLE) SLICE=$(SLICE) APPLY=$(APPLY) ROOT=$(CURDIR) \
	    GO_VERSION=$(GO_VERSION) PROGRESS_FLAG='$(PROGRESS_FLAG)'
