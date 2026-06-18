#!/usr/bin/env bash
# image-test.sh — DC-FA-DIST-001-Akzeptanzkriterien gegen das lokal
# gebaute Runtime-Image (slice-010):
#
#   (1) Happy:    Repo mit kaputtem Link → Befund-Ausgabe und
#                 Exit-Code des Containers byte-identisch zur nativen
#                 Ausführung desselben Binaries (DC-QA-02).
#   (2) Boundary: read-only-Mount (:ro) → vollständige Prüfung ohne
#                 Schreibfehler (sauberes Fixture → Exit 0).
#   (3) Negative: fehlender /repo-Mount → Exit 2 mit Mount-Hinweis.
#   (4) Modi:     --doctor und --repair → Ausgabe nativ vs. Container
#                 byte-identisch (DC-QA-02; CLI-Optionen als Container-
#                 Argumente identisch zur nativen Ausführung,
#                 DC-FA-DIST-001 / DC-FA-CLI-007 / DC-FA-CLI-008).
#
# „Nativ" in einem Docker-only-Repo: das statische Binary wird aus dem
# Runtime-Image extrahiert (docker cp) und direkt ausgeführt — kein
# Host-Go (AGENTS.md §3.1). Annahme: Host-Architektur = Image-
# Architektur (amd64); auf abweichenden Hosts (arm64) bricht der
# Nativ-Lauf laut ab — dann Binary via qemu/binfmt ausführen oder den
# Vergleich auf zwei Container-Varianten umstellen.
set -euo pipefail

IMAGE="${IMAGE:-d-check}"
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

fail() {
  echo "image-test: FAIL — $1" >&2
  exit 1
}

# Binary aus dem Runtime-Image extrahieren (identisches Artefakt).
cid="$(docker create "$IMAGE":latest)"
docker cp -q "$cid":/d-check "$WORK/d-check"
docker rm "$cid" > /dev/null
chmod +x "$WORK/d-check"

# Fixture: ein gültiger Link, ein kaputter Link.
mkdir -p "$WORK/fixture/docs"
printf '# A\n\n[ok](b.md)\n[kaputt](fehlt.md)\n' > "$WORK/fixture/docs/a.md"
printf 'ziel\n' > "$WORK/fixture/docs/b.md"

# --- (1) Happy: nativ vs. Container byte-identisch ------------------
native_exit=0
"$WORK/d-check" "$WORK/fixture" > "$WORK/native.out" 2> "$WORK/native.err" || native_exit=$?
container_exit=0
docker run --rm --network none -v "$WORK/fixture":/repo:ro "$IMAGE":latest \
  > "$WORK/container.out" 2> "$WORK/container.err" || container_exit=$?

[ "$native_exit" -eq 1 ] || fail "nativer Lauf: Exit $native_exit, want 1"
[ "$container_exit" -eq 1 ] || fail "Container-Lauf: Exit $container_exit, want 1"
cmp -s "$WORK/native.out" "$WORK/container.out" \
  || fail "stdout nativ vs. Container nicht byte-identisch (DC-QA-02)"
cmp -s "$WORK/native.err" "$WORK/container.err" \
  || fail "stderr nativ vs. Container nicht byte-identisch"
grep -q 'fehlt.md' "$WORK/container.out" || fail "Befund fehlt in der Ausgabe"
echo "image-test: (1) Happy — nativ und Container byte-identisch, Exit 1"

# --- (2) Boundary: read-only-Mount, sauberes Fixture → Exit 0 -------
rm "$WORK/fixture/docs/a.md"
printf '# A\n\n[ok](b.md)\n' > "$WORK/fixture/docs/a.md"
ro_exit=0
docker run --rm --network none -v "$WORK/fixture":/repo:ro "$IMAGE":latest \
  > /dev/null 2> "$WORK/ro.err" || ro_exit=$?
[ "$ro_exit" -eq 0 ] || fail "read-only-Lauf: Exit $ro_exit, want 0 (stderr: $(cat "$WORK/ro.err"))"
echo "image-test: (2) Boundary — read-only-Mount, vollständige Prüfung, Exit 0"

# --- (3) Negative: kein Mount → Exit 2 + Mount-Hinweis --------------
nomount_exit=0
docker run --rm --network none "$IMAGE":latest \
  > /dev/null 2> "$WORK/nomount.err" || nomount_exit=$?
[ "$nomount_exit" -eq 2 ] || fail "Lauf ohne Mount: Exit $nomount_exit, want 2"
grep -q '/repo gemountet' "$WORK/nomount.err" \
  || fail "Mount-Hinweis fehlt auf stderr: $(cat "$WORK/nomount.err")"
echo "image-test: (3) Negative — kein Mount, Exit 2 mit Mount-Hinweis"

# --- (4) Modi: --doctor und --repair nativ vs. Container ------------
# Eigenes Fixture mit einer nackten Kennung (id-unlinked), damit --repair
# einen nicht-leeren Patch liefert (ids-Modul + existierendes Target).
mkdir -p "$WORK/idsfix/docs/plan/adr"
cat > "$WORK/idsfix/.d-check.yml" <<'YML'
modules: [ids]
ids:
  patterns:
    - regex: 'ADR-\d{4}'
      target: docs/plan/adr/
YML
printf '# ADR\n' > "$WORK/idsfix/docs/plan/adr/0042-x.md"
printf '# Doc\n\nnacktes ADR-0042 im Text\n' > "$WORK/idsfix/docs/a.md"

for mode in doctor repair; do
  nx=0
  "$WORK/d-check" "--$mode" "$WORK/idsfix" > "$WORK/n.$mode.out" 2> "$WORK/n.$mode.err" || nx=$?
  cx=0
  docker run --rm --network none -v "$WORK/idsfix":/repo:ro "$IMAGE":latest "--$mode" \
    > "$WORK/c.$mode.out" 2> "$WORK/c.$mode.err" || cx=$?
  [ "$nx" -eq "$cx" ] || fail "--$mode: Exit nativ $nx != Container $cx"
  [ "$nx" -eq 1 ] || fail "--$mode: Exit $nx, want 1"
  cmp -s "$WORK/n.$mode.out" "$WORK/c.$mode.out" \
    || fail "--$mode: stdout nativ vs. Container nicht byte-identisch (DC-QA-02)"
  cmp -s "$WORK/n.$mode.err" "$WORK/c.$mode.err" \
    || fail "--$mode: stderr nativ vs. Container nicht byte-identisch"
done
grep -q 'Diagnose' "$WORK/c.doctor.out" || fail "--doctor: Diagnose-Ausgabe fehlt"
grep -q '+nacktes \[`ADR-0042`\]' "$WORK/c.repair.out" || fail "--repair: erwarteter Hunk fehlt"
echo "image-test: (4) Modi — --doctor/--repair nativ == Container, Exit 1"

echo "image-test: OK — DC-FA-DIST-001-Akzeptanzkriterien erfüllt"
