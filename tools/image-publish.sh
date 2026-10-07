#!/usr/bin/env bash
# image-publish.sh — veröffentlicht den Multi-Plattform-Index
# (DC-FA-DIST-001, ADR-0102) in drei Schritten, damit unter keinem Tag eine
# ungeprüfte Variante liegt:
#
#   (1) Build beider Plattformen und Push OHNE Tag (push-by-digest) — der
#       Index liegt nur unter seinem Digest in der Registry.
#   (2) Gegenprobe am gepushten Index (image-verify-published.sh):
#       Plattformen, Labels je Plattform, Binary je Plattform gleich dem
#       geprüften.
#   (3) Erst danach die Tags v$VERSION (und :latest bei PUBLISH_LATEST=true,
#       ADR-0014) auf genau diesen Digest — eine Index-Kopie innerhalb der
#       Registry, kein zweiter Bau.
#
# Grenze: fällt (2), liegt der ungetaggte Index unter seinem Digest in der
# Registry; erreichbar ist er nur über diesen Digest, den die Meldung nennt.
# Bricht (3) zwischen den Tags ab, kann v$VERSION schon gesetzt sein und
# :latest noch nicht — jeder gesetzte Tag zeigt dann auf den geprüften
# Digest, nie auf einen anderen.
#
# Eingaben: PUBLISH_REPO, VERSION, PUBLISH_LATEST, GO_VERSION,
# GOLANGCI_LINT_VERSION, TESTED_AMD64, TESTED_ARM64, optional PROGRESS_FLAG.
# Letzte Ausgabezeile bei Erfolg: `image-publish: <repo>@<index-digest>`.
set -euo pipefail

: "${PUBLISH_REPO:?PUBLISH_REPO fehlt (z. B. ghcr.io/pt9912/d-check)}"
: "${VERSION:?VERSION fehlt}" "${GO_VERSION:?GO_VERSION fehlt}"
: "${GOLANGCI_LINT_VERSION:?GOLANGCI_LINT_VERSION fehlt}"
: "${TESTED_AMD64:?TESTED_AMD64 fehlt}" "${TESTED_ARM64:?TESTED_ARM64 fehlt}"
PUBLISH_LATEST="${PUBLISH_LATEST:-false}"

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

fail() {
  echo "image-publish: FAIL — $1" >&2
  exit 1
}

# Abbrüche über set -e (ein Docker-Aufruf scheitert) nennen ihren Schritt und
# den Stand, den er hinterlässt.
stage="(1) Build/Push ohne Tag — kein Tag gesetzt"
trap 'echo "image-publish: FAIL — Schritt $stage" >&2' ERR

# --- (1) Build und Push ohne Tag ------------------------------------
docker buildx build ${PROGRESS_FLAG:+"$PROGRESS_FLAG"} --platform linux/amd64,linux/arm64 \
  --provenance=false --sbom=false \
  --build-arg GO_VERSION="$GO_VERSION" \
  --build-arg GOLANGCI_LINT_VERSION="$GOLANGCI_LINT_VERSION" \
  --build-arg VERSION="$VERSION" --target runtime \
  --output "type=image,name=$PUBLISH_REPO,push-by-digest=true,name-canonical=true,push=true" \
  --metadata-file "$WORK/meta.json" .

digest="$(grep -o '"containerimage.digest": *"sha256:[0-9a-f]\{64\}"' "$WORK/meta.json" \
  | grep -o 'sha256:[0-9a-f]\{64\}' || true)"
[ -n "$digest" ] || fail "Index-Digest nicht aus den Build-Metadaten lesbar — nichts getaggt"
ref="$PUBLISH_REPO@$digest"
echo "image-publish: (1) Index ohne Tag gepusht — $ref"

# --- (2) Gegenprobe am gepushten Index ------------------------------
stage="(2) Gegenprobe — kein Tag gesetzt; der ungetaggte Index liegt unter $ref"
REF="$ref" VERSION="$VERSION" TESTED_AMD64="$TESTED_AMD64" TESTED_ARM64="$TESTED_ARM64" \
  bash tools/image-verify-published.sh \
  || fail "Gegenprobe rot — KEIN Tag gesetzt; der ungetaggte Index liegt unter $ref"
echo "image-publish: (2) Gegenprobe grün"

# --- (3) Tags auf den geprüften Digest ------------------------------
stage="(3) Tags — v$VERSION kann gesetzt sein, :latest noch nicht; gesetzte Tags zeigen auf $digest"
tags=(-t "$PUBLISH_REPO:v$VERSION")
[ "$PUBLISH_LATEST" != true ] || tags+=(-t "$PUBLISH_REPO:latest")
docker buildx imagetools create "${tags[@]}" "$ref"
check_tags=("v$VERSION")
[ "$PUBLISH_LATEST" != true ] || check_tags+=(latest)
for t in "${check_tags[@]}"; do
  tagged="$(docker buildx imagetools inspect --format '{{.Manifest.Digest}}' "$PUBLISH_REPO:$t" 2>/dev/null || true)"
  [ "$tagged" = "$digest" ] \
    || fail "$PUBLISH_REPO:$t zeigt auf [${tagged:-<leer>}], nicht auf den geprüften $digest"
done
echo "image-publish: (3) getaggt — v$VERSION$([ "$PUBLISH_LATEST" != true ] || echo ' + latest')"
echo "image-publish: $ref"
