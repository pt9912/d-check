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
#   --check-latest  Upstream-Audit (Netz, informativ, KEIN Gate, KEIN
#              fail-closed — anders als das netzlose fail-closed --verify).
#              Zwei Prüfungen gegen Upstream:
#              (A) Currency: gepinnter Tag vs. neuestes STABILES Release
#                  (GitHub-API /releases/latest; Prereleases/Drafts raus).
#              (B) Content-Drift am GEPINNTEN Tag: lädt dessen lab-regelwerk.zip
#                  und vergleicht die Bytes (sha256sum regelwerk/*.md) gegen das
#                  committete SHA256SUMS — fällt auf, wenn der Tag verschoben /
#                  das Asset neu hochgeladen wurde (VERIFIZIERT die
#                  Tag-Immutabilität aus MR-011, statt ihr nur zu vertrauen).
#              Exit (schlimmster Fall): 0 = aktuell & authentisch; 3 = neuerer
#              Tag verfügbar (Signal); 4 = Content-Drift am gepinnten Tag
#              (Provenienz!). Nicht erreichbare Teile (Netz/API/Werkzeug/Manifest)
#              -> SKIP je Teil (exit 0, sofern der andere Teil nicht 3/4 meldet).
#              Nudge zum bewussten Re-Adopt (MR-019/MR-022), nie automatischer
#              Bump. Übernommen aus dem Kurs-Beispiel check_regelwerk_drift.py
#              (Content-Hash-Drift), auf d-checks Tag-Pin-Modell übersetzt.
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
case "${1:-}" in
  --verify)       mode="verify"; shift ;;
  --check-latest) mode="check-latest"; shift ;;
esac

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

check_latest() {
  # Upstream-Audit: (A) Currency (neuerer Tag?) und (B) Content-Drift am
  # gepinnten Tag (Bytes des Release-Assets vs. committetes SHA256SUMS — Tag
  # verschoben/Asset neu?). Beides Netz, informativ, KEIN Gate, KEIN fail-closed
  # (Netz-/Werkzeug-Ausfall -> SKIP je Teil). (B) VERIFIZIERT die
  # Tag-Immutabilität (MR-011); automatisiert nichts (MR-019/MR-022).
  # Läuft unter deaktiviertem errexit (Dispatch: set +e), daher kein `|| true`.
  command -v curl >/dev/null 2>&1 \
    || { echo "fetch-baseline-cache: 'curl' nicht gefunden (Host-Werkzeug)" >&2; exit 1; }

  # --- (A) Currency: neuestes stabiles Release ---
  local api="https://api.github.com/repos/${repo}/releases/latest"
  local latest currency="skip"
  latest="$(curl -fsSL -H 'Accept: application/vnd.github+json' "$api" 2>/dev/null \
    | grep -oE '"tag_name"[[:space:]]*:[[:space:]]*"v[0-9]+\.[0-9]+\.[0-9]+"' \
    | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' | head -1)"
  if [ -n "$latest" ]; then
    if [ "$latest" = "$tag" ]; then
      currency="current"
    else
      local newest; newest="$(printf '%s\n%s\n' "$tag" "$latest" | sort -V | tail -1)"
      [ "$newest" = "$latest" ] && currency="newer" || currency="ahead"
    fi
  fi

  # --- (B) Content-Drift am gepinnten Tag (verifiziert Tag-Immutabilität) ---
  local authenticity="skip" note=""
  if ! command -v unzip >/dev/null 2>&1 || ! command -v sha256sum >/dev/null 2>&1; then
    note="unzip/sha256sum fehlt (Host-Werkzeug)"
  elif [ ! -f "$sums" ]; then
    note="kein vendored Manifest ${sums}"
  else
    local td; td="$(mktemp -d)"
    local url="https://github.com/${repo}/releases/download/${tag}/lab-regelwerk.zip"
    if curl -fsSL -o "${td}/lab-regelwerk.zip" "$url" 2>/dev/null \
       && unzip -oq "${td}/lab-regelwerk.zip" -d "${td}/regelwerk" 2>/dev/null; then
      local up ve
      up="$( { cd "$td" && sha256sum regelwerk/*.md 2>/dev/null; } | LC_ALL=C sort )"
      ve="$( LC_ALL=C sort < "$sums" )"
      [ "$up" = "$ve" ] && authenticity="ok" || authenticity="drift"
    else
      note="Asset ${tag}/lab-regelwerk.zip nicht ladbar (Netz/API)"
    fi
    rm -rf "$td"
  fi

  # --- Ausgabe beider Teile; Exit = schlimmster Fall (4 > 3 > 0) ---
  local rc=0
  case "$currency" in
    current) echo "fetch-baseline-cache: check-latest OK (Currency) — Baseline aktuell (gepinnt == latest == ${tag})." ;;
    newer)   echo "fetch-baseline-cache: check-latest — NEUER RELEASE ${latest} verfügbar (gepinnt ${tag})." >&2
             echo "  -> Re-Adopt erwägen: fetch-baseline-cache.sh ${latest} (re-vendor), §Baseline-Stand+Pin bumpen, Adaptionen prüfen (MR-019/MR-020/MR-021)." >&2
             rc=3 ;;
    ahead)   echo "fetch-baseline-cache: check-latest (Currency) — latest ${latest} ist nicht neuer als gepinnt ${tag} (Pin voraus?); nichts zu tun." ;;
    skip)    echo "fetch-baseline-cache: check-latest SKIP (Currency) — neuestes Release nicht ermittelbar (Netz/API/Rate-Limit). Gepinnt: ${tag}." ;;
  esac
  case "$authenticity" in
    ok)    echo "fetch-baseline-cache: check-latest OK (Content) — gepinnter Tag ${tag} upstream unverändert (Bytes == vendored SHA256SUMS)." ;;
    drift) echo "fetch-baseline-cache: check-latest — UPSTREAM-CONTENT-DRIFT: die Bytes von ${tag}/lab-regelwerk.zip weichen vom vendored SHA256SUMS ab (Tag verschoben / Asset neu?)." >&2
           echo "  -> Provenienz prüfen: das gepinnte Release-Asset hat sich geändert. NICHT blind re-vendoren; upstream verifizieren (MR-011-Immutabilität verletzt)." >&2
           rc=4 ;;
    skip)  echo "fetch-baseline-cache: check-latest SKIP (Content) — ${note}. Gepinnt: ${tag}." ;;
  esac
  return $rc
}

if [ "$mode" = "verify" ]; then
  verify
  exit 0
fi

if [ "$mode" = "check-latest" ]; then
  set +e; check_latest; rc=$?; set -e
  exit "$rc"
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
