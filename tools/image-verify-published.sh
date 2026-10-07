#!/usr/bin/env bash
# image-verify-published.sh — prüft einen gepushten Multi-Plattform-Index
# gegen die vorher geprüften lokalen Bilder (DC-FA-DIST-001, ADR-0102).
#
#   (1) Plattformen: der Index trägt genau linux/amd64 und linux/arm64 —
#       keine fehlende Variante, kein Attestations-Eintrag.
#   (2) Labels: je Plattform sind die fünf OCI-Labels gesetzt, und
#       org.opencontainers.image.version entspricht VERSION.
#   (3) Binary: je Plattform ist /d-check im gepushten Index sha256-gleich
#       zu dem im geprüften lokalen Bild — Prüfung und Push sind zwei
#       buildx-Aufrufe, und erst dieser Vergleich macht das veröffentlichte
#       Binary zum geprüften.
#
# Gelesen wird aus der Registry, nicht aus dem lokalen Daemon (ADR-0102):
# zwei aus demselben lokalen Bild abgeleitete Werte wären trivial gleich.
# Jeder Abweichung folgt Exit 1; eine
# Antwort, die sich nicht lesen lässt, ist ebenfalls Exit 1 — ungeprüft ist
# nicht bestätigt.
#
# Eingaben: REF (gepushte Referenz, Tag oder `<repo>@<digest>`),
# VERSION, TESTED_AMD64 und TESTED_ARM64 (lokale, geprüfte Bilder).
set -euo pipefail

: "${REF:?REF fehlt}" "${VERSION:?VERSION fehlt}"
: "${TESTED_AMD64:?TESTED_AMD64 fehlt}" "${TESTED_ARM64:?TESTED_ARM64 fehlt}"

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

fail() {
  echo "image-verify-published: FAIL — $1" >&2
  exit 1
}

# --- (1) Plattformen ------------------------------------------------
platforms="$(docker buildx imagetools inspect "$REF" \
  --format '{{range .Manifest.Manifests}}{{.Platform.OS}}/{{.Platform.Architecture}}{{"\n"}}{{end}}' \
  | sed '/^$/d' | sort | tr '\n' ' ')" || fail "$REF nicht lesbar oder kein Multi-Plattform-Index"
[ "$platforms" = "linux/amd64 linux/arm64 " ] \
  || fail "Index $REF trägt [${platforms% }], verlangt [linux/amd64 linux/arm64]"
echo "image-verify-published: (1) Plattformen — linux/amd64 linux/arm64"

# --- (2) Labels je Plattform ----------------------------------------
for label in \
  org.opencontainers.image.source \
  org.opencontainers.image.description \
  org.opencontainers.image.licenses \
  org.opencontainers.image.title \
  org.opencontainers.image.vendor \
  org.opencontainers.image.version
do
  docker buildx imagetools inspect "$REF" \
    --format "{{range \$p, \$i := .Image}}{{\$p}}={{index \$i.Config.Labels \"$label\"}}{{\"\n\"}}{{end}}" \
    | sed '/^$/d' > "$WORK/labels" || fail "Labels von $REF nicht lesbar"
  [ "$(wc -l < "$WORK/labels")" -eq 2 ] || fail "$label: nicht je Plattform lesbar ($(tr '\n' ' ' < "$WORK/labels"))"
  while IFS='=' read -r plat value; do
    [ -n "$value" ] || fail "$label leer auf $plat"
    if [ "$label" = org.opencontainers.image.version ] && [ "$value" != "$VERSION" ]; then
      fail "$label=[$value] auf $plat entspricht nicht der Tag-Version [$VERSION]"
    fi
  done < "$WORK/labels"
done
echo "image-verify-published: (2) Labels — je Plattform gesetzt, version=$VERSION"

# --- (3) Binary je Plattform: gepusht == geprüft --------------------
binary_sha() {
  local plat="$1" ref="$2" cid
  cid="$(docker create --platform "$plat" "$ref")" || return 1
  docker cp -q "$cid":/d-check "$WORK/bin" || { docker rm "$cid" > /dev/null; return 1; }
  docker rm "$cid" > /dev/null
  sha256sum "$WORK/bin" | cut -d' ' -f1
}
# Gezogen wird je Plattform über ihren Manifest-Digest aus dem Index: der
# Daemon bindet eine Digest-Referenz an genau ein lokales Bild, ein zweites
# Ziehen derselben Digest-Referenz für die andere Plattform scheitert; über den
# Manifest-Digest je Plattform arbeiten Tag- und Digest-Referenz gleich.
case "$REF" in *@*) repo="${REF%%@*}" ;; *) repo="${REF%:*}" ;; esac
docker buildx imagetools inspect "$REF" \
  --format '{{range .Manifest.Manifests}}{{.Platform.OS}}/{{.Platform.Architecture}} {{.Digest}}{{"\n"}}{{end}}' \
  > "$WORK/manifests" || fail "Manifeste von $REF nicht lesbar"
for pair in "linux/amd64=$TESTED_AMD64" "linux/arm64=$TESTED_ARM64"; do
  plat="${pair%%=*}" tested="${pair#*=}"
  mdigest="$(awk -v p="$plat" '$1==p { print $2 }' "$WORK/manifests")"
  [ -n "$mdigest" ] || fail "kein Manifest für $plat in $REF"
  pushed="$repo@$mdigest"
  docker pull -q --platform "$plat" "$pushed" > /dev/null || fail "$pushed ($plat) nicht ziehbar"
  pushed_sha="$(binary_sha "$plat" "$pushed")" || fail "Binary aus $pushed ($plat) nicht lesbar"
  tested_sha="$(binary_sha "$plat" "$tested")" || fail "Binary aus $tested ($plat) nicht lesbar"
  [ "$pushed_sha" = "$tested_sha" ] \
    || fail "$plat: Binary im gepushten Index ($pushed_sha) ist nicht das geprüfte aus $tested ($tested_sha)"
  echo "image-verify-published: (3) $plat — Binary gleich dem geprüften (${pushed_sha:0:12})"
done

digest="$(docker buildx imagetools inspect "$REF" --format '{{.Manifest.Digest}}')" \
  || fail "Index-Digest von $REF nicht lesbar"
echo "image-verify-published: OK — $REF (Index $digest)"
