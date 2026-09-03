# archive-wave — eigenständiges Makefile, damit dieses Verzeichnis
# unverändert in ein Fremd-Repo kopiert werden kann und dort dieselbe
# Docker/make-only-Disziplin trägt wie das aufrufende Repo (AGENTS.md
# §3.1, falls vorhanden): kein roher `docker`-Aufruf in der Doku, jede
# Operation hinter einem Target.
#
# Aufrufbar direkt (aus diesem Verzeichnis) oder per `make -C
# tools/archive-wave <target>` aus dem Wurzel-Makefile eines
# aufrufenden Repos — beide Wege landen auf denselben Targets.

GO_VERSION    ?= 1.27.0
IMAGE         ?= archive-wave
PROGRESS_FLAG ?=
# abspath: `docker run -v` verlangt einen absoluten Host-Pfad, ein
# relativer würde als Named-Volume-Name fehlinterpretiert.
ROOT          ?= $(abspath ../..)

.PHONY: help test build run

help: ## Targets anzeigen
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | sort

test: ## Testsuite fahren (Docker, kein Host-Go)
	docker build $(PROGRESS_FLAG) --build-arg GO_VERSION=$(GO_VERSION) \
	    --target test -f Dockerfile -t $(IMAGE):test .

build: ## Binary-Image bauen
	docker build $(PROGRESS_FLAG) --build-arg GO_VERSION=$(GO_VERSION) \
	    -f Dockerfile -t $(IMAGE):latest .

run: build ## Welle archivieren: make run WELLE=welle-NN [ROOT=/pfad] [APPLY=1].
	@if [ -z "$(WELLE)" ]; then \
	    echo "run: WELLE=welle-NN ist Pflicht" >&2; exit 2; \
	fi
	docker run --rm -u "$$(id -u):$$(id -g)" \
	    -v "$(ROOT)":/repo$(if $(APPLY),,:ro) \
	    $(IMAGE):latest -root=/repo -welle=$(WELLE) $(if $(APPLY),-apply,)
