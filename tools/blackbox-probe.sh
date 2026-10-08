#!/usr/bin/env bash
# blackbox-probe.sh — vergleicht das Verhalten des Werkzeugs VOR und NACH einer
# Aenderung, von aussen: dieselben Fixtures, dieselben Ausgabeformen, stdout,
# stderr und Exit GETRENNT (zusammengefuehrte Streams verschraenken sich und
# zeigen Abweichungen, die keine sind).
#
# ZUSAGE: meldet jede Abweichung zwischen dem Vorher-Image (gebaut aus
# `git archive $REF`) und dem Nachher-Image ($IMAGE:latest, `make build`) ueber
# die Fixtures unter tools/blackbox-probe/fixtures/ und das Repo selbst.
# ABGRENZUNG: kein Gate — ob eine Abweichung gewollt ist, entscheidet der
# Vorgang, der sie erzeugt; das Werkzeug urteilt nicht ueber den Repo-Zustand.
# GRENZE: es findet einen Unterschied, den niemand erwartet; dass eine
# erwartete Gleichheit die falsche Zusage ist, findet es nicht. Es vergleicht
# nur die Faelle, die Fixtures und Formen abdecken.
#
# Eingaben: REF (git-Referenz des Vorher-Stands, Pflicht), IMAGE (Default
# d-check), PROBE_FORMS (Default: Standard, --json, --yaml, --doctor).
# Exit: 0 = byte-identisch, 1 = Abweichung(en), 2 = Lauf gescheitert.
set -uo pipefail

REF="${REF:-}"
IMAGE="${IMAGE:-d-check}"
NACHHER="${IMAGE}:latest"
VORHER="${IMAGE}:probe-vorher"
FIXTURES_DIR="tools/blackbox-probe/fixtures"
# `-` steht fuer den Lauf ohne Zusatz-Argument; eine leere Form waere in einer
# Wortliste nicht darstellbar.
PROBE_FORMS="${PROBE_FORMS:-- --json --yaml --doctor}"

fail() {
  echo "blackbox-probe: GESCHEITERT — $1" >&2
  exit 2
}

[ -n "$REF" ] || fail "REF fehlt (z. B. make blackbox-probe REF=HEAD~1)"
git rev-parse --verify --quiet "${REF}^{commit}" > /dev/null \
  || fail "REF '${REF}' ist kein Commit dieses Repos"
[ -n "$(printf '%s' "$PROBE_FORMS" | tr -d '[:space:]')" ] \
  || fail "PROBE_FORMS ist leer — nichts zu vergleichen ist kein gleicher Stand"
docker image inspect "$NACHHER" > /dev/null 2>&1 \
  || fail "Nachher-Image ${NACHHER} fehlt (make build)"

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

# --- Vorher-Image aus dem Stand von REF -----------------------------
# Gebaut mit dem Dockerfile und den Build-Arg-Defaults DIESES Stands; VERSION
# gleich dem Nachher-Default, damit die eingebettete Version keinen
# Unterschied erzeugt, der keine Verhaltensaenderung ist.
mkdir -p "$WORK/src"
git archive "$REF" | tar -x -C "$WORK/src" || fail "git archive ${REF} gescheitert"
echo "blackbox-probe: baue Vorher-Image aus $(git rev-parse --short "$REF") …"
docker build -q --build-arg VERSION=0.0.0-dev --target runtime -t "$VORHER" "$WORK/src" > /dev/null \
  || fail "Vorher-Image aus ${REF} nicht baubar"

# --- Fixtures --------------------------------------------------------
# Kopie mit Leserechten fuer den nonroot-Nutzer des Images (65532); der
# Arbeitsbaum selbst bleibt unberuehrt.
mkdir -p "$WORK/fx"
fixtures=()
for d in "$FIXTURES_DIR"/*/; do
  [ -d "$d" ] || continue
  name="$(basename "$d")"
  cp -r "$d" "$WORK/fx/$name"
  fixtures+=("$name=$WORK/fx/$name")
done
[ "${#fixtures[@]}" -gt 0 ] || fail "keine Fixtures unter ${FIXTURES_DIR}/"
chmod -R a+rX "$WORK/fx"
fixtures+=("repo=$PWD")

# --- Laeufe und Vergleich -------------------------------------------
lauf() { # image dir form outprefix
  local rc=0
  if [ "$3" = "-" ]; then
    docker run --rm --network none -v "$2":/repo:ro "$1" > "$4.out" 2> "$4.err" || rc=$?
  else
    docker run --rm --network none -v "$2":/repo:ro "$1" "$3" > "$4.out" 2> "$4.err" || rc=$?
  fi
  echo "$rc" > "$4.rc"
}

vergleiche=0
abweichungen=0
for fx in "${fixtures[@]}"; do
  name="${fx%%=*}" dir="${fx#*=}"
  for form in $PROBE_FORMS; do
    tag="${name}/${form}"
    p="$WORK/run.${name}.${form//[^a-z-]/_}"
    lauf "$VORHER" "$dir" "$form" "$p.v"
    lauf "$NACHHER" "$dir" "$form" "$p.n"
    vergleiche=$((vergleiche + 1))
    diffs=""
    cmp -s "$p.v.out" "$p.n.out" || diffs="${diffs} stdout"
    cmp -s "$p.v.err" "$p.n.err" || diffs="${diffs} stderr"
    [ "$(cat "$p.v.rc")" = "$(cat "$p.n.rc")" ] \
      || diffs="${diffs} exit($(cat "$p.v.rc")→$(cat "$p.n.rc"))"
    if [ -z "$diffs" ]; then
      echo "  gleich       ${tag} (exit $(cat "$p.n.rc"))"
    else
      abweichungen=$((abweichungen + 1))
      echo "  ABWEICHUNG   ${tag}:${diffs}"
      for s in out err; do
        cmp -s "$p.v.$s" "$p.n.$s" && continue
        diff -u --label "vorher.std${s}" --label "nachher.std${s}" "$p.v.$s" "$p.n.$s" \
          | head -n 20 | sed 's/^/      /'
      done
    fi
  done
done

if [ "$abweichungen" -gt 0 ]; then
  echo "blackbox-probe: ${abweichungen} von ${vergleiche} Vergleichen weichen ab (Vorher $(git rev-parse --short "$REF"))."
  exit 1
fi
echo "blackbox-probe: byte-identisch über ${vergleiche} Vergleiche (Vorher $(git rev-parse --short "$REF"), stdout/stderr/Exit getrennt)."
