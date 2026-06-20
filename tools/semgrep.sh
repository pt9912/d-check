#!/usr/bin/env bash
# semgrep.sh — hermetisches Security-/Static-Analysis-Gate über ein
# GEPINNTES semgrep-Image und ein GEPINNTES, lokal gecachtes Regelset
# (Docker/make-only, AGENTS.md §3.1; ADR-0010).
#
# Hermetik (DC-QA-02 Determinismus, DC-QA-03 Netzlosigkeit):
#   - Regelset NICHT ins Repo vendort (Semgrep Rules License v1.0), sondern
#     `semgrep/semgrep-rules` an einem festen Commit-Pin EINMALIG in einen
#     Cache AUSSERHALB des Repos geholt (wie `go mod`/Image-Pull — Setup-Netz,
#     nicht Teil der Analyse). Außerhalb des Repos, damit die `.go`-Fixtures
#     der Regeln weder `go list ./...` noch den d-check-Selbstscan stören.
#   - Der Scan selbst läuft `--network none` mit gepinntem Image + lokalen
#     Regeln. Gepinnter Commit + gepinntes Image ⇒ identische Befunde.
#
# Umfang `go/lang/security` (kuratiert, hoch-Signal für die Go-Codebasis;
# bash/dockerfile bewusst ausgelassen — siehe ADR-0010). 0 Befunde auf
# d-check ⇒ keine zentrale `--exclude-rule`-Ausnahme nötig.
set -euo pipefail

SEMGREP_VERSION="${SEMGREP_VERSION:-1.167.0}"
RULES_COMMIT="${SEMGREP_RULES_COMMIT:-d41fb34cf74466e2878af5f268ebf54466a04541}"
RULES_SUBSET="go/lang/security"
RULES_REMOTE="https://github.com/semgrep/semgrep-rules.git"

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CACHE_ROOT="${SEMGREP_RULES_CACHE:-${XDG_CACHE_HOME:-$HOME/.cache}/d-check/semgrep-rules}"
RULES_DIR="$CACHE_ROOT/$RULES_COMMIT"

# --- Setup: Regel-Cache am Pin holen (einmalig, Netz) ---------------
# Idempotent über den Commit-Schlüssel: ist der Umfang am Pin schon da,
# wird nichts geholt. Atomar via .tmp + mv, damit ein Abbruch keinen
# halben Cache als „fertig" zurücklässt.
if [ ! -d "$RULES_DIR/$RULES_SUBSET" ]; then
  echo "semgrep: hole Regel-Cache am Pin ${RULES_COMMIT} (einmalig, Netz; wie go mod/Image-Pull) ..." >&2
  rm -rf "${RULES_DIR}.tmp"
  mkdir -p "${RULES_DIR}.tmp"
  git -C "${RULES_DIR}.tmp" init -q
  git -C "${RULES_DIR}.tmp" remote add origin "$RULES_REMOTE"
  git -C "${RULES_DIR}.tmp" fetch -q --depth 1 origin "$RULES_COMMIT"
  git -C "${RULES_DIR}.tmp" checkout -q FETCH_HEAD
  mv "${RULES_DIR}.tmp" "$RULES_DIR"
fi

# --- Scan: netzlos, gepinntes Image + lokale Regeln -----------------
# --disable-version-check: ohne ihn liefe semgreps Versions-Ping unter
# `--network none` in einen ~2-min-Timeout (statt ~3 s Scan).
out="$(mktemp)"
trap 'rm -f "$out"' EXIT

set +e
docker run --rm --init --network none \
  -v "$ROOT:/src:ro" \
  -v "$RULES_DIR:/rules:ro" \
  "semgrep/semgrep:${SEMGREP_VERSION}" \
  semgrep scan --error --metrics off --disable-version-check \
    --config "/rules/${RULES_SUBSET}" \
    /src 2>&1 | tee "$out"
status="${PIPESTATUS[0]}"
set -e

# Schutz gegen stilles Grün (Review R1, HIGH-1): `--error` setzt den
# Exit-Code nur BEI Befunden, nicht bei 0 GELADENEN Regeln. Ein
# regel-leeres oder upstream umbenanntes Subset (z. B. nach Pin-Hebung)
# ergäbe sonst „0 findings" -> Exit 0 -> grünes Gate ohne Scan. Wir
# verlangen daher die positive „Ran N rules"-Zeile (N>=1) als Nachweis,
# dass tatsächlich Regeln liefen — sonst lautes Rot statt stillem Grün.
if ! grep -qE 'Ran [1-9][0-9]* rules' "$out"; then
  echo "semgrep: FEHLER — 0 Regeln geladen (Umfang ${RULES_SUBSET} am Pin ${RULES_COMMIT} leer/umbenannt?); breche ab statt still grün." >&2
  exit 2
fi
exit "$status"
