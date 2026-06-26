#!/usr/bin/env bash
# fetch-baseline-cache — materialisiert die lokale Lese-Form der adoptierten
# Baseline (AGENTS.md §1). Regelwerk und Templates werden unterschiedlich
# verkörpert (MR-019 löst MR-017 ab; MR-018):
#
#   regelwerk → .harness/baseline/<tag>/regelwerk/   (COMMITTET, vendored)
#               + .harness/baseline/<tag>/SHA256SUMS  (Integritäts-/Provenienz-
#               Manifest über die vendorten Dateien)
#   templates → .harness/cache/<tag>/templates/       (gitignored, ephemer;
#               nur Adoptions-/Drift-Audit-Staging, MR-018)
#
# Modi:
#   (default)  re-vendor: zieht die Release-Assets, entpackt das Regelwerk in
#              den committeten Vendor-Pfad, (re)generiert SHA256SUMS und
#              verifiziert; staged die Templates im ephemeren Cache. Netz nötig
#              (Release-Download) — Anlass: Baseline-Pin-Bump.
#   --verify   nur Integritätsprüfung des committeten Regelwerks gegen
#              SHA256SUMS. Offline, kein Netz — für CI/Audit/frischen Checkout.
#
# Tag-Quelle: ohne Argument die §Baseline-Stand-Zeile in
# harness/conventions.md (Single Source of Truth — kein Drift; der nächste
# Baseline-Bump zieht automatisch die neue Version); mit Argument ein
# expliziter Tag (z. B. `v1.4.0`).
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

repo="pt9912/ai-harness-course"
conventions="harness/conventions.md"

mode="vendor"
if [ "${1:-}" = "--verify" ]; then mode="verify"; shift; fi

tag="${1:-}"
if [ -z "$tag" ]; then
  tag="$(grep -m1 '\*\*Stand:\*\*' "$conventions" \
    | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' | head -1 || true)"
fi
if ! [[ "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "fetch-baseline-cache: ungültiger/leerer Tag '${tag}' — Argument vMAJOR.MINOR.PATCH angeben oder §Baseline in ${conventions} prüfen" >&2
  exit 1
fi

baseline=".harness/baseline/${tag}"
sums="${baseline}/SHA256SUMS"

verify() {
  # Integritätsprüfung des committeten Regelwerks gegen SHA256SUMS (offline).
  command -v sha256sum >/dev/null 2>&1 \
    || { echo "fetch-baseline-cache: 'sha256sum' nicht gefunden (Host-Werkzeug)" >&2; exit 1; }
  [ -f "$sums" ] \
    || { echo "fetch-baseline-cache: ${sums} fehlt — erst re-vendor (ohne --verify) laufen" >&2; exit 1; }
  echo "fetch-baseline-cache: verify ${baseline}/regelwerk gegen SHA256SUMS"
  ( cd "$baseline" && sha256sum -c SHA256SUMS )
  echo "fetch-baseline-cache: verify ok"
}

if [ "$mode" = "verify" ]; then
  verify
  exit 0
fi

# --- re-vendor (Netz) ---
for cmd in curl unzip sha256sum; do
  command -v "$cmd" >/dev/null 2>&1 \
    || { echo "fetch-baseline-cache: '${cmd}' nicht gefunden (Host-Werkzeug)" >&2; exit 1; }
done

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

unpack() {
  # $1 = Release-Asset, $2 = Zielverzeichnis
  local asset="$1" dest="$2"
  local url="https://github.com/${repo}/releases/download/${tag}/${asset}"
  echo "fetch-baseline-cache: ${tag}/${asset} -> ${dest}/"
  curl -fsSL -o "${tmp}/${asset}" "$url"
  rm -rf "$dest"
  mkdir -p "$dest"
  unzip -oq "${tmp}/${asset}" -d "$dest"
}

# Regelwerk: committeter Vendor-Pfad + Manifest über die vendorten Dateien.
unpack lab-regelwerk.zip "${baseline}/regelwerk"
( cd "$baseline" && sha256sum regelwerk/*.md > SHA256SUMS )
verify

# Templates: ephemerer Cache, nur Adoptions-/Drift-Audit-Staging (MR-018).
unpack lab-templates.zip ".harness/cache/${tag}/templates"

echo "fetch-baseline-cache: fertig — vendored ${baseline}/regelwerk (+SHA256SUMS), staged .harness/cache/${tag}/templates"
