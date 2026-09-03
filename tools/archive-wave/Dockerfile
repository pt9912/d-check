# syntax=docker/dockerfile:1.7

# ---------------------------------------------------------------------------
# archive-wave — Wellen-Archivierungs-Werkzeug (Baseline-Regelwerk
# modul-06-roadmap.md §Wellen-Closure-Prozedur, Schritt 4).
#
# Eigenstaendiges Modul, portabel: kein Import aus d-checks internen
# Paketen. Baubar mit demselben Docker/make-only-Workflow wie d-check
# selbst (AGENTS.md §3.1) -- jedes Repo mit demselben Planning-Layout kann
# dieses Verzeichnis unveraendert uebernehmen und mit `docker build` bauen.
#
# Build-Kontext ist DIESES Verzeichnis (tools/archive-wave/), nicht die
# Repo-Wurzel -- siehe das `archive-wave`-Target im Wurzel-Makefile.
# ---------------------------------------------------------------------------

ARG GO_VERSION=1.27.0

# ---- deps -------------------------------------------------------------------
FROM golang:${GO_VERSION}@sha256:4013ae0f9e7994f8535c58c811f8f863fbed38b72e0d51e6592156f758d66146 AS deps

WORKDIR /src
ENV GOFLAGS="-mod=readonly -buildvcs=false" \
    GOMODCACHE=/go/pkg/mod \
    GOCACHE=/root/.cache/go-build

COPY go.mod ./
COPY go.su[m] ./
RUN mkdir -p "$GOMODCACHE" && go mod download

# ---- test ---------------------------------------------------------------
# Eigenes Ziel, weil test/build getrennte Stages im Docker/make-only-Workflow
# sind (dasselbe Muster wie d-checks eigenes Dockerfile) -- die Fixture-
# Verifikation und die drei Umkehr-Proben (BEO-023) laufen hier, nicht nur
# einmalig von Hand.
FROM deps AS test

COPY . .
RUN CGO_ENABLED=0 go test ./...

# ---- build --------------------------------------------------------------
FROM deps AS build

COPY . .
RUN CGO_ENABLED=0 go build -o /out/archive-wave .

# ---- runtime ---------------------------------------------------------------
# Dasselbe Basis-Image + Digest wie d-checks eigenes Runtime-Image (ADR-0002,
# ADR-0011) -- der im Image erklaerte "nonroot"-User ist nur ein Default;
# `docker run -u $(id -u):$(id -g)` (im Makefile-Target) ueberschreibt ihn
# zur Laufzeit, damit geschriebene Dateien dem aufrufenden Bediener gehoeren,
# nicht einer festen Image-UID -- dieses Werkzeug schreibt bewusst in das
# gemountete Repo, anders als d-checks eigenes Read-Only-Image.
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab AS runtime

COPY --from=build /out/archive-wave /archive-wave

ENTRYPOINT ["/archive-wave"]
