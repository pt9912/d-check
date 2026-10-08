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
# Das Repo als Kopie seiner nicht ignorierten Dateien, wie die Fixtures unter
# $WORK/fx und gleich lesbar gemacht: EINE Mount-Quelle fuer alle Laeufe, die
# der Kanarienlauf mitbelegt. Ein direkter Mount des Arbeitsbaums waere eine
# zweite Quelle mit eigenen Rechten (umask), die keine Kanarie prueft.
mkdir -p "$WORK/fx/repo"
# Getrackte, aber im Arbeitsbaum geloeschte Dateien fallen heraus — kopiert
# wird, was der Arbeitsbaum traegt, nicht was der Index nennt.
git ls-files -z -co --exclude-standard \
  | while IFS= read -r -d '' f; do { [ -e "$f" ] || [ -L "$f" ]; } && printf '%s\0' "$f"; done \
  | tar --null -T - -cf - | tar -x -C "$WORK/fx/repo" \
  || fail "Kopie des Arbeitsbaums gescheitert"
chmod -R a+rX "$WORK/fx"
fixtures+=("repo=$WORK/fx/repo")

# --- Laeufe und Vergleich -------------------------------------------
lauf() { # image dir form outprefix
  local rc=0
  if [ "$3" = "-" ]; then
    docker run --rm --network none -v "$2":/repo:ro "$1" > "$4.out" 2> "$4.err" || rc=$?
  else
    docker run --rm --network none -v "$2":/repo:ro "$1" "$3" > "$4.out" 2> "$4.err" || rc=$?
  fi
  echo "$rc" > "$4.rc"
  # d-check endet nur mit 0, 1 oder 2 (auch ein Absturz endet mit 2); jeder
  # andere Exit heisst, der Container lief nicht oder wurde beendet
  # (125/126/127 Start, 137 OOM) — kein Lauf, den man vergleichen koennte.
  # 0, 1 und 2 sind Verhalten des Werkzeugs und werden verglichen.
  case "$rc" in
    0|1|2) ;;
    *) fail "${1} auf ${2} (${3}): Exit ${rc}, kein Lauf des Werkzeugs — $(tail -n 1 "$4.err")" ;;
  esac
}

# Kanarienlauf vor und nach allen Vergleichen, auf beiden Images, ueber den
# Exit-Vertrag allein: `sauber` muss mit 0 enden, `links` mit 1. Er belegt
# Daemon, die Mount-Quelle aller Laeufe ($WORK/fx) und den Inhalt, den der
# Container sieht — ein leerer Mount endet mit 2, ein Mount ohne Markdown
# bei `links` mit 0, ein Daemon-Ausfall bei `sauber` mit 1. Kein Wortlaut
# einer Meldung: ein Vorgang, der die Ausgabe aendert, bricht ihn nicht.
kanarie() { # wann
  local img fx want ursache
  for img in "$VORHER" "$NACHHER"; do
    for fx in sauber:0 links:1; do
      want="${fx#*:}"; fx="${fx%%:*}"
      [ -d "$WORK/fx/$fx" ] || fail "Kanarien-Fixture ${FIXTURES_DIR}/${fx} fehlt"
      lauf "$img" "$WORK/fx/$fx" - "$WORK/kanarie"
      [ "$(cat "$WORK/kanarie.rc")" = "$want" ] && continue
      # Besteht das Vorher-Image denselben Kanarienlauf und nur das Nachher
      # nicht, liegt eher eine Aenderung am Exit-Vertrag vor als eine
      # Stoerung der Umgebung — die Meldung sagt das, statt zu raten.
      ursache="die Umgebung traegt keinen Vergleich (Daemon, Mount oder gesehener Inhalt)"
      [ "$img" = "$NACHHER" ] \
        && ursache="das Vorher-Image bestand denselben Lauf — vermutlich eine Aenderung am Exit-Vertrag des Nachher-Stands, keine Umgebungsstoerung"
      fail "Kanarienlauf ${1} auf ${img} (${fx}): Exit $(cat "$WORK/kanarie.rc") statt ${want} — $(tail -n 1 "$WORK/kanarie.err"); ${ursache}"
    done
  done
}

kanarie "vorher"
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
kanarie "nachher"

if [ "$abweichungen" -gt 0 ]; then
  echo "blackbox-probe: ${abweichungen} von ${vergleiche} Vergleichen weichen ab (Vorher $(git rev-parse --short "$REF"))."
  exit 1
fi
echo "blackbox-probe: byte-identisch über ${vergleiche} Vergleiche (Vorher $(git rev-parse --short "$REF"), stdout/stderr/Exit getrennt, Kanarienlauf vorher und nachher)."
