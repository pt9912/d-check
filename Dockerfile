# syntax=docker/dockerfile:1.7

# ---------------------------------------------------------------------------
# d-check — Doc-Referenz-Checker.
#
# Docker-only-Workflow (AGENTS.md §3.1; Vorbild:
# https://github.com/pt9912/u-boot): das Repo hat
# keine Host-Go-Anforderung. Build/Lint/Test laufen über
# `docker build --target <stage>` und sind im Makefile gewrappt;
# das Architektur-Gate (make arch-check) läuft seit ADR-0029 über das
# digest-gepinnte a-check-Image (a-check.mk + .a-check.yml).
#
# Stages:
#   deps       — Go-Modul-Auflösung (Cache-Layer).
#   compile    — schnelles Compile-Feedback ohne Tests/Lint.
#   lint       — golangci-lint mit dem Projekt-Profil.
#   test       — `go test ./...`.
#   coverage   — `go test -coverpkg` + tools/coverage-gate.sh
#                (Kalibrierungs-Bindung, harness/README §Sensors).
#   build      — statisch gelinktes Binary (CGO=0, -ldflags "-s -w").
#   runtime    — distroless/static:nonroot (ADR-0002).
#
# Pin-Politik (u-boot-Konvention): GO_VERSION und
# GOLANGCI_LINT_VERSION sind Routine-Pins (Tag); zusätzlich ist jede
# FROM-Zeile per @sha256:-Digest gepinnt (ADR-0011, ADR-0002 §1) — der
# Digest ist die Wahrheit, der Tag bleibt lesbar daneben. Eine Hebung
# (Version UND Digest gemeinsam) ist ein bewusster Commit, Begründung
# im Body.
# ---------------------------------------------------------------------------

ARG GO_VERSION=1.27.1
ARG GOLANGCI_LINT_VERSION=v2.13.1

# ---- deps ------------------------------------------------------------------
FROM golang:${GO_VERSION}@sha256:512690a5660563b57d37ecc31129e7f136e831db2aed24a1dbeb8ad7380dc0fa AS deps

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
FROM golangci/golangci-lint:${GOLANGCI_LINT_VERSION}@sha256:ba07dffad130794ae79ebaa0056809d18c0168f3f846480ffd3eb6c04578b83d AS lint

WORKDIR /src
COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY . .
RUN golangci-lint run ./...

# ---- test ------------------------------------------------------------------
FROM deps AS test

COPY . .
RUN CGO_ENABLED=0 go test ./...

# ---- coverage --------------------------------------------------------------
# Kalibrierungs-Bindung (harness/README.md §Sensors): Schwelle 93 %
# (Kalibrierung 2026-06-11; zuvor Ramp 85 → 90 bei welle-03 done);
# Verfehlung ⇒ Carveout-Pflicht.
# `-coverpkg` misst über die Paketgrenzen von ./internal/... (u-boot-
# Muster) — sonst zählt nur paket-lokale Abdeckung.
# `pipefail` via SHELL, damit `go test … | tee` den Exit-Code nicht
# maskiert.
FROM deps AS coverage

SHELL ["/bin/bash", "-eo", "pipefail", "-c"]

ARG COVERAGE_THRESHOLD=93
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

# VERSION (Git-Tag) ins Binary einbetten — Quelle des Image-Refs in
# `--print-mk` (DC-FA-CLI-010, slice-038). Default für Dev-/Gate-Builds;
# die Release-Pipeline setzt den Tag (make ci VERSION=…).
ARG VERSION=0.0.0-dev
COPY . .
RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w -X 'github.com/pt9912/d-check/internal/adapter/driving/cli.version=${VERSION}'" \
    -o /out/d-check \
    ./cmd/d-check

# ---- runtime ---------------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab AS runtime

# VERSION wird von der Release-Pipeline aus dem Git-Tag durchgereicht
# (make ci VERSION=…); der Workflow pinnt das Label gegen den Tag —
# ein Build mit Version-Drift darf nicht shippen (release.yml).
ARG VERSION=0.0.0-dev

LABEL org.opencontainers.image.source="https://github.com/pt9912/d-check" \
      org.opencontainers.image.description="d-check — Doc-Referenz-Checker für Markdown-Dokumentation." \
      org.opencontainers.image.title="d-check" \
      org.opencontainers.image.vendor="pt9912" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.version="${VERSION}"

COPY --from=build /out/d-check /d-check

USER 65532:65532
# Default-Befehl: Prüfung von /repo; CLI-Optionen werden als
# Container-Argumente angehängt (DC-FA-DIST-001, ADR-0002).
ENTRYPOINT ["/d-check", "/repo"]
