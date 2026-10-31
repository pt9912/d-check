#!/usr/bin/env bash
# image-test.sh — Akzeptanzkriterien der Distribution (DC-FA-DIST-001) und
# der Modi, die an ihrer Grenze gelten, gegen das lokal gebaute Runtime-Image:
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
#   (5) Doku:     --manual ohne Netz und ohne Mount → Abschnitte aus
#                 Handbuch und Spezifikation, nativ vs. Container
#                 byte-identisch (DC-FA-CLI-013).
#   (6) metrics:  Messung aus git im Container unter :ro, mit fremdem Benutzer
#                 und ohne .git (DC-FA-MET-001).
#   (7) metrics:  Schätzung aus der Tabelle des Plans (DC-FA-MET-001).
#   (8) pinsync:  Wertvergleich unter :ro, mit fremdem Benutzer, ohne Netz
#                 (DC-FA-PSYNC-001).
#
# „Nativ" in einem Docker-only-Repo: das statische Binary wird aus dem
# Runtime-Image extrahiert (docker cp) und direkt ausgeführt — kein
# Host-Go (AGENTS.md §3.1).
#
# Plattform: IMAGE_REF nennt das Bild (Default `$IMAGE:latest`), PLATFORM
# optional die Plattform (`linux/amd64`, `linux/arm64`; leer = die des
# Hosts). Container und extrahiertes Binary laufen auf derselben Plattform;
# die ELF-Maschine des Binaries wird gegen sie geprüft, sonst prüfte ein
# Bild, das still die Host-Variante liefert, die falsche (DC-FA-DIST-001,
# ADR-0102). Grenze: eine Plattform, die nicht die des Hosts ist, braucht
# binfmt/QEMU auf dem Host — fehlt es, bricht der Lauf mit Hinweis ab. Der Host
# selbst muss Linux sein: ein Linux-Binary läuft auf macOS nicht nativ.
set -euo pipefail

IMAGE="${IMAGE:-d-check}"
REF="${IMAGE_REF:-$IMAGE:latest}"
PLATFORM="${PLATFORM:-}"
PLAT=()
[ -z "$PLATFORM" ] || PLAT=(--platform "$PLATFORM")
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

fail() {
  echo "image-test: FAIL — $1" >&2
  exit 1
}

# Binary aus dem Runtime-Image extrahieren (identisches Artefakt).
cid="$(docker create ${PLAT[@]+"${PLAT[@]}"} "$REF")"
docker cp -q "$cid":/d-check "$WORK/d-check"
docker rm "$cid" > /dev/null
chmod +x "$WORK/d-check"

# ELF-Maschine (Byte 18) gegen die verlangte Plattform: 3e = x86-64,
# b7 = AArch64.
want_arch="${PLATFORM#linux/}"
if [ -z "$want_arch" ]; then
  case "$(uname -m)" in x86_64) want_arch=amd64 ;; aarch64|arm64) want_arch=arm64 ;; *) want_arch="$(uname -m)" ;; esac
fi
case "$want_arch" in amd64) want_mach=3e ;; arm64) want_mach=b7 ;; *) fail "Plattform $want_arch nicht unterstützt" ;; esac
got_mach="$(od -An -tx1 -j18 -N1 "$WORK/d-check" | tr -d ' \n')"
[ "$got_mach" = "$want_mach" ] \
  || fail "Binary in $REF ist ELF-Maschine 0x$got_mach, verlangt $want_arch (0x$want_mach)"
echo "image-test: Plattform $want_arch — $REF"

# Fixture: ein gültiger Link, ein kaputter Link.
mkdir -p "$WORK/fixture/docs"
printf '# A\n\n[ok](b.md)\n[kaputt](fehlt.md)\n' > "$WORK/fixture/docs/a.md"
printf 'ziel\n' > "$WORK/fixture/docs/b.md"

# --- (1) Happy: nativ vs. Container byte-identisch ------------------
# abdeckung: DC-FA-DIST-001, DC-QA-02
native_exit=0
"$WORK/d-check" "$WORK/fixture" > "$WORK/native.out" 2> "$WORK/native.err" || native_exit=$?
if [ "$native_exit" -eq 126 ]; then
  [ "$(uname -s)" = Linux ] \
    || fail "Host ist $(uname -s), kein Linux — das Linux-Binary läuft hier nicht nativ (Wrapper: docs/maintainer/releasing.md)"
  fail "Binary für $want_arch auf diesem Host nicht ausführbar — binfmt/QEMU für $want_arch fehlt"
fi
container_exit=0
docker run --rm ${PLAT[@]+"${PLAT[@]}"} --network none -v "$WORK/fixture":/repo:ro "$REF" \
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
# abdeckung: DC-FA-DIST-001, DC-QA-03
rm "$WORK/fixture/docs/a.md"
printf '# A\n\n[ok](b.md)\n' > "$WORK/fixture/docs/a.md"
ro_exit=0
docker run --rm ${PLAT[@]+"${PLAT[@]}"} --network none -v "$WORK/fixture":/repo:ro "$REF" \
  > /dev/null 2> "$WORK/ro.err" || ro_exit=$?
[ "$ro_exit" -eq 0 ] || fail "read-only-Lauf: Exit $ro_exit, want 0 (stderr: $(cat "$WORK/ro.err"))"
echo "image-test: (2) Boundary — read-only-Mount, vollständige Prüfung, Exit 0"

# --- (3) Negative: kein Mount → Exit 2 + Mount-Hinweis --------------
# abdeckung: DC-FA-DIST-001
nomount_exit=0
docker run --rm ${PLAT[@]+"${PLAT[@]}"} --network none "$REF" \
  > /dev/null 2> "$WORK/nomount.err" || nomount_exit=$?
[ "$nomount_exit" -eq 2 ] || fail "Lauf ohne Mount: Exit $nomount_exit, want 2"
grep -q '/repo gemountet' "$WORK/nomount.err" \
  || fail "Mount-Hinweis fehlt auf stderr: $(cat "$WORK/nomount.err")"
echo "image-test: (3) Negative — kein Mount, Exit 2 mit Mount-Hinweis"

# --- (4) Modi: --doctor und --repair nativ vs. Container ------------
# abdeckung: DC-FA-CLI-007, DC-FA-CLI-008, DC-QA-02
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
  docker run --rm ${PLAT[@]+"${PLAT[@]}"} --network none -v "$WORK/idsfix":/repo:ro "$REF" "--$mode" \
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

# --- (5) Doku aus dem Image: --manual netzlos, nativ vs. Container ---
# abdeckung: DC-FA-CLI-013, DC-QA-02
nx=0
"$WORK/d-check" --manual planning > "$WORK/n.manual.out" 2> "$WORK/n.manual.err" || nx=$?
cx=0
docker run --rm ${PLAT[@]+"${PLAT[@]}"} --network none "$REF" --manual planning \
  > "$WORK/c.manual.out" 2> "$WORK/c.manual.err" || cx=$?
[ "$cx" -eq 0 ] || fail "--manual: Exit $cx im Container ohne Netz und ohne Mount, want 0"
[ "$nx" -eq "$cx" ] || fail "--manual: Exit nativ $nx != Container $cx"
cmp -s "$WORK/n.manual.out" "$WORK/c.manual.out" \
  || fail "--manual: stdout nativ vs. Container nicht byte-identisch (DC-QA-02)"
grep -q '^==> docs/user/benutzerhandbuch.md:' "$WORK/c.manual.out" \
  || fail "--manual: kein Abschnitt aus dem Handbuch"
grep -q '^==> spec/spezifikation.md:' "$WORK/c.manual.out" \
  || fail "--manual: kein Abschnitt aus der Spezifikation"
echo "image-test: (5) --manual — netzlos im Container, nativ == Container, Exit 0"

# --- (6) metrics: Messung aus git unter :ro, fremdem Eigentümer, ohne .git ---
# abdeckung: DC-FA-MET-001
# Das Fixture ist ein echtes git-Repo (Basis-Commit, danach Änderungen über der
# Marke). Der Container läuft als anderer Benutzer als der Eigentümer des Repos
# und mit :ro; ohne `.git` endet der Lauf mit Exit 2 (nicht „0 Befunde").
mfx="$WORK/metricsfix"
mkdir -p "$mfx/src" "$mfx/plans"
cat > "$mfx/.d-check.yml" <<'YML'
modules: [links]
metrics:
  plans-dir: plans
  block: groessenmarke
  surcharge: {numerator: 135, denominator: 10000}
  rounding: 100
  table: {header: ["Nr.", "Klasse"], key-column: "Nr.", class-column: "Klasse", classes: {x: 1000}}
YML
printf 'a\nb\n' > "$mfx/src/x.txt"
git -C "$mfx" init -q .
git -C "$mfx" add -A
git -C "$mfx" -c user.name=t -c user.email=t@example.invalid commit -q -m base
base="$(git -C "$mfx" rev-parse --short=7 HEAD)"
seq 1 1600 > "$mfx/src/x.txt"
printf '# Plan\n\n```groessenmarke\npath: src/\nbase: %s\nestimate: 1000\nlimit: 1400\n```\n\n| Nr. | Klasse |\n|---|---|\n| 1 | x |\n' "$base" > "$mfx/plans/slice-1.md"
git -C "$mfx" add -A
git -C "$mfx" -c user.name=t -c user.email=t@example.invalid commit -q -m head
chmod -R a+rX "$mfx"
mx=0
docker run --rm ${PLAT[@]+"${PLAT[@]}"} --network none --user 12345:12345 -v "$mfx":/repo:ro "$REF" --enable metrics \
  > "$WORK/m.out" 2> "$WORK/m.err" || mx=$?
[ "$mx" -eq 1 ] || fail "metrics: Exit $mx unter :ro und fremdem Benutzer, want 1 ($(cat "$WORK/m.err"))"
grep -q 'metrics-over-limit' "$WORK/m.out" || fail "metrics: metrics-over-limit fehlt: $(cat "$WORK/m.out")"
rm -rf "$mfx/.git"
mx=0
docker run --rm ${PLAT[@]+"${PLAT[@]}"} --network none -v "$mfx":/repo:ro "$REF" --enable metrics \
  > "$WORK/m2.out" 2> "$WORK/m2.err" || mx=$?
[ "$mx" -eq 2 ] || fail "metrics ohne .git: Exit $mx, want 2"
echo "image-test: (6) metrics — Messung unter :ro und fremdem Benutzer, Exit 1; ohne .git Exit 2"

# --- (7) metrics: Schätzung aus der Tabelle — Befund bei Messung 0, Hinweis danach ---
# abdeckung: DC-FA-MET-001
tfx="$WORK/tablefix"
mkdir -p "$tfx/src" "$tfx/plans"
cp "$WORK/metricsfix/.d-check.yml" "$tfx/.d-check.yml"
printf 'a\n' > "$tfx/src/x.txt"
git -C "$tfx" init -q .
git -C "$tfx" add -A
git -C "$tfx" -c user.name=t -c user.email=t@example.invalid commit -q -m base
tbase="$(git -C "$tfx" rev-parse --short=7 HEAD)"
printf '# P1\n\n```groessenmarke\npath: src/\nbase: %s\nestimate: 999\nlimit: 1400\n```\n\n| Nr. | Klasse |\n|---|---|\n| 1 | x |\n' "$tbase" > "$tfx/plans/slice-1.md"
printf '# P2\n\n```groessenmarke\npath: src/\n```\n\n| Nr. | Klasse |\n|---|---|\n| 1 | x |\n' > "$tfx/plans/slice-2.md"
git -C "$tfx" add -A
git -C "$tfx" -c user.name=t -c user.email=t@example.invalid commit -q -m plans
tx=0
docker run --rm ${PLAT[@]+"${PLAT[@]}"} --network none -v "$tfx":/repo:ro "$REF" --enable metrics \
  > "$WORK/t.out" 2> "$WORK/t.err" || tx=$?
[ "$tx" -eq 1 ] || fail "metrics Tabelle: Exit $tx bei Messung 0, want 1 ($(cat "$WORK/t.err"))"
grep -q 'metrics-estimate-mismatch' "$WORK/t.out" || fail "metrics Tabelle: metrics-estimate-mismatch fehlt: $(cat "$WORK/t.out")"
grep -q 'Marke ausstehend (estimate 1000, limit 1400)' "$WORK/t.err" || fail "metrics Tabelle: Hinweis mit Zahlen fehlt: $(cat "$WORK/t.err")"
printf 'a\nb\n' > "$tfx/src/x.txt"
git -C "$tfx" add -A
git -C "$tfx" -c user.name=t -c user.email=t@example.invalid commit -q -m code
tx=0
docker run --rm ${PLAT[@]+"${PLAT[@]}"} --network none -v "$tfx":/repo:ro "$REF" --enable metrics \
  > "$WORK/t2.out" 2> "$WORK/t2.err" || tx=$?
[ "$tx" -eq 0 ] || fail "metrics Tabelle: Exit $tx bei Messung >= 1, want 0 ($(cat "$WORK/t2.out"))"
grep -q 'weicht von der Tabelle ab' "$WORK/t2.err" || fail "metrics Tabelle: Hinweis bei Messung >= 1 fehlt: $(cat "$WORK/t2.err")"
echo "image-test: (7) metrics — Schätzung aus der Tabelle: Befund bei Messung 0, Hinweis bei Messung >= 1"

# --- (8) pinsync: Wertvergleich unter :ro, fremdem Eigentümer, ohne Netz ---
# abdeckung: DC-FA-PSYNC-001
# Dasselbe Ergebnis gegen ein Repo, das einem anderen Benutzer gehört und nur
# lesbar eingehängt ist, wie gegen ein eigenes beschreibbares: gleiche Werte
# Exit 0, ungleiche Exit 1 mit pinsync-mismatch, eine Symlink-Quelle Exit 2.
pfx="$WORK/pinsyncfix"
mkdir -p "$pfx"
cat > "$pfx/.d-check.yml" <<'YML'
modules: [links]
pinsync:
  pins:
    - name: go
      sources:
        - {file: Dockerfile, pattern: '^ARG GO_VERSION=v?(\S+)$'}
        - {file: Makefile, pattern: '^GO_VERSION\s*\?=\s*v?(\S+)$'}
YML
printf 'ARG GO_VERSION=1.27.3\n' > "$pfx/Dockerfile"
printf 'GO_VERSION ?= v1.27.3\n' > "$pfx/Makefile"
chmod -R a+rX "$pfx"
px=0
docker run --rm ${PLAT[@]+"${PLAT[@]}"} --network none --user 12345:12345 -v "$pfx":/repo:ro "$REF" --enable pinsync \
  > "$WORK/p.out" 2> "$WORK/p.err" || px=$?
[ "$px" -eq 0 ] || fail "pinsync: Exit $px bei gleichen Werten unter :ro und fremdem Benutzer, want 0 ($(cat "$WORK/p.out") $(cat "$WORK/p.err"))"
grep -q 'pinsync: 1 Pin(s) geprüft, 2 Quelle(n) gelesen' "$WORK/p.err" || fail "pinsync: Bezugsmengen-Zeile fehlt: $(cat "$WORK/p.err")"
printf 'GO_VERSION ?= v1.27.4\n' > "$pfx/Makefile"
px=0
docker run --rm ${PLAT[@]+"${PLAT[@]}"} --network none --user 12345:12345 -v "$pfx":/repo:ro "$REF" --enable pinsync \
  > "$WORK/p1.out" 2> "$WORK/p1.err" || px=$?
[ "$px" -eq 1 ] || fail "pinsync: Exit $px bei ungleichen Werten unter :ro und fremdem Benutzer, want 1 ($(cat "$WORK/p1.err"))"
grep -q 'pinsync-mismatch' "$WORK/p1.out" || fail "pinsync: pinsync-mismatch fehlt: $(cat "$WORK/p1.out")"
px=0
docker run --rm ${PLAT[@]+"${PLAT[@]}"} --network none -v "$pfx":/repo "$REF" --enable pinsync \
  > "$WORK/p2.out" 2> "$WORK/p2.err" || px=$?
[ "$px" -eq 1 ] || fail "pinsync: Exit $px gegen ein beschreibbares Repo, want 1"
cmp -s "$WORK/p1.out" "$WORK/p2.out" || fail "pinsync: Ausgabe unter :ro und fremdem Benutzer weicht vom beschreibbaren Lauf ab"
rm "$pfx/Makefile"
printf 'GO_VERSION ?= 1.27.3\n' > "$pfx/echt.mk"
ln -s echt.mk "$pfx/Makefile"
chmod -R a+rX "$pfx"
px=0
docker run --rm ${PLAT[@]+"${PLAT[@]}"} --network none --user 12345:12345 -v "$pfx":/repo:ro "$REF" --enable pinsync \
  > "$WORK/p3.out" 2> "$WORK/p3.err" || px=$?
[ "$px" -eq 2 ] || fail "pinsync: Exit $px bei Symlink-Quelle, want 2 ($(cat "$WORK/p3.err"))"
echo "image-test: (8) pinsync — gleiche Werte Exit 0, ungleiche Exit 1 (:ro, fremder Benutzer == beschreibbar), Symlink-Quelle Exit 2"

echo "image-test: OK — Akzeptanzkriterien erfüllt (DC-FA-DIST-001, DC-FA-CLI-007, DC-FA-CLI-008, DC-FA-CLI-013, DC-FA-MET-001, DC-FA-PSYNC-001)"
