#!/usr/bin/env bash
# fetch-baseline-cache — materialisiert die lokale Lese-Form der adoptierten
# Baseline (AGENTS.md §1) im gitignorierten Cache: lädt die Release-Assets
# `lab-regelwerk.zip` + `lab-templates.zip` und entpackt sie nach
# `.harness/cache/<tag>/{regelwerk,templates}/`.
#
# Tag-Quelle: ohne Argument die §Baseline-Stand-Zeile in
# harness/conventions.md (Single Source of Truth — kein Drift; der nächste
# Baseline-Bump zieht automatisch die neue Version); mit Argument ein
# expliziter Tag (z. B. `v1.4.0`). Der Cache ist ephemer/gitignored und aus
# dem Dogfooding-Selbst-Scan ausgenommen (harness/conventions.md MR-017).
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

repo="pt9912/ai-harness-course"
conventions="harness/conventions.md"

tag="${1:-}"
if [ -z "$tag" ]; then
  tag="$(grep -m1 '\*\*Stand:\*\*' "$conventions" \
    | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' | head -1 || true)"
fi
if ! [[ "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "fetch-baseline-cache: ungültiger/leerer Tag '${tag}' — Argument vMAJOR.MINOR.PATCH angeben oder §Baseline in ${conventions} prüfen" >&2
  exit 1
fi

for cmd in curl unzip; do
  command -v "$cmd" >/dev/null 2>&1 \
    || { echo "fetch-baseline-cache: '${cmd}' nicht gefunden (Host-Werkzeug)" >&2; exit 1; }
done

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

fetch() {
  # $1 = Release-Asset, $2 = Ziel-Unterverzeichnis im Cache
  local asset="$1" subdir="$2"
  local url="https://github.com/${repo}/releases/download/${tag}/${asset}"
  local dest=".harness/cache/${tag}/${subdir}"
  echo "fetch-baseline-cache: ${tag}/${asset} -> ${dest}/"
  curl -fsSL -o "${tmp}/${asset}" "$url"
  rm -rf "$dest"
  mkdir -p "$dest"
  unzip -oq "${tmp}/${asset}" -d "$dest"
}

fetch lab-regelwerk.zip regelwerk
fetch lab-templates.zip templates

echo "fetch-baseline-cache: fertig — .harness/cache/${tag}/{regelwerk,templates}/"
