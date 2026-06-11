# syntax=docker/dockerfile:1.7

# ---------------------------------------------------------------------------
# d-check — Doc-Referenz-Checker.
#
# Docker-only-Workflow (AGENTS.md §3.1; Vorbild:
# https://github.com/pt9912/u-boot): das Repo hat
# keine Host-Go-Anforderung. Build/Lint/Test/arch-check laufen über
# `docker build --target <stage>` und sind im Makefile gewrappt.
#
# Stages:
#   deps       — Go-Modul-Auflösung (Cache-Layer).
#   compile    — schnelles Compile-Feedback ohne Tests/Lint.
#   lint       — golangci-lint mit dem Projekt-Profil.
#   test       — `go test ./...`.
#   arch-check — Fitness Function zu ADR-0005 (Import-Regeln).
#   coverage   — `go test -coverpkg` + tools/coverage-gate.sh
#                (Kalibrierungs-Bindung 85 % → 90 %).
#   build      — statisch gelinktes Binary (CGO=0, -ldflags "-s -w").
#   runtime    — distroless/static:nonroot (ADR-0002).
#
# Pin-Politik (u-boot-Konvention): GO_VERSION und
# GOLANGCI_LINT_VERSION sind Routine-Pins; Hebung ohne eigene ADR,
# Begründung im Commit-Body.
# ---------------------------------------------------------------------------

ARG GO_VERSION=1.26.4
ARG GOLANGCI_LINT_VERSION=v2.12.2

# ---- deps ------------------------------------------------------------------
FROM golang:${GO_VERSION} AS deps

WORKDIR /src
ENV GOFLAGS="-mod=readonly -buildvcs=false" \
    GOMODCACHE=/go/pkg/mod \
    GOCACHE=/root/.cache/go-build

COPY go.mod ./
# go.su[m]-Trick (u-boot/k-deskflight): matcht go.sum, falls vorhanden,
# und still nichts im Bootstrap-Zustand vor `go mod tidy`.
COPY go.su[m] ./

RUN mkdir -p "$GOMODCACHE" && go mod download

# ---- compile ---------------------------------------------------------------
FROM deps AS compile

COPY . .
RUN CGO_ENABLED=0 go build -o /tmp/d-check ./cmd/d-check

# ---- lint ------------------------------------------------------------------
FROM golangci/golangci-lint:${GOLANGCI_LINT_VERSION} AS lint

WORKDIR /src
COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY . .
RUN golangci-lint run ./...

# ---- test ------------------------------------------------------------------
FROM deps AS test

COPY . .
RUN CGO_ENABLED=0 go test ./...

# ---- arch-check ------------------------------------------------------------
# Fitness Function zu ADR-0005 (Hexagon light, u-boot-Ordnerkonvention):
# Import-Regeln werden über `go list` geprüft — Verstöße brechen den
# Build (strukturelle Durchsetzung von DC-QA-03).
FROM deps AS arch-check

COPY . .
RUN bash tools/arch-check.sh

# ---- coverage --------------------------------------------------------------
# Kalibrierungs-Bindung (harness/README.md §Sensors): Schwelle 85 %,
# welle-03 done → 90 %; Verfehlung nach Trigger ⇒ Carveout-Pflicht.
# `-coverpkg` misst über die Paketgrenzen von ./internal/... (u-boot-
# Muster) — sonst zählt nur paket-lokale Abdeckung.
# `pipefail` via SHELL, damit `go test … | tee` den Exit-Code nicht
# maskiert.
FROM deps AS coverage

SHELL ["/bin/bash", "-eo", "pipefail", "-c"]

ARG COVERAGE_THRESHOLD=85
ENV COVERAGE_THRESHOLD=${COVERAGE_THRESHOLD}

COPY . .
RUN mkdir -p /out && \
    COVERPKG=$(go list ./internal/... | tr '\n' ',' | sed 's/,$//') && \
    CGO_ENABLED=0 go test \
        -coverpkg="$COVERPKG" \
        -coverprofile=/out/coverage.out \
        -covermode=atomic \
        ./... && \
    go tool cover -func=/out/coverage.out | tee /out/coverage-func.txt && \
    bash tools/coverage-gate.sh /out/coverage-func.txt "$COVERAGE_THRESHOLD"

# ---- build -----------------------------------------------------------------
FROM deps AS build

COPY . .
RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -o /out/d-check \
    ./cmd/d-check

# ---- runtime ---------------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

LABEL org.opencontainers.image.source="https://github.com/pt9912/d-check" \
      org.opencontainers.image.description="d-check — Doc-Referenz-Checker für Markdown-Dokumentation." \
      org.opencontainers.image.title="d-check" \
      org.opencontainers.image.vendor="pt9912"

COPY --from=build /out/d-check /d-check

USER 65532:65532
# Default-Befehl: Prüfung von /repo; CLI-Optionen werden als
# Container-Argumente angehängt (DC-FA-DIST-001, ADR-0002).
ENTRYPOINT ["/d-check", "/repo"]
