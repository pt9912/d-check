#!/usr/bin/env bash
# planning-consistency.sh — Meta-Gate gegen Planning-Drift (slice-040):
#
#   Die Roadmap darf "Keine aktive Welle" GENAU DANN behaupten, wenn kein
#   `slice-*` in docs/plan/planning/in-progress/ liegt — und umgekehrt.
#   Die Lifecycle-/Roadmap-Konsistenz ist in AGENTS.md §3.3 dokumentiert,
#   war aber bis slice-040 nicht maschinell erzwungen (Nutzer-Audit
#   2026-06-21: roadmap.md sagte "Keine aktive Welle", während ein Slice
#   in in-progress/ lag).
#
#   Regel:
#     hasSlices := es existiert >=1 in-progress/slice-*.md
#     hasActive := der Abschnitt "## Aktuelle Welle" der roadmap.md enthält
#                  NICHT den Marker "Keine aktive Welle"
#     konsistent <=> hasActive == hasSlices   (sonst Exit 1, fail-closed)
#
# Negativ-Test: vor der echten Prüfung läuft ein Selbsttest, der beide
# Inkonsistenz-Richtungen nachweislich feuern lässt (analog
# tools/gate-consistency.sh).
set -euo pipefail

cd "$(dirname "$0")/.."

PROGRESS_DIR="docs/plan/planning/in-progress"
ROADMAP="$PROGRESS_DIR/roadmap.md"
ACTIVE_MARKER="Keine aktive Welle"

# hasSlices: >=1 Datei in-progress/slice-*.md? (Glob ohne Treffer bleibt
# wörtlich → -e ist dann false; nullglob brauchen wir nicht.)
has_slices() { # $1 progress-dir
  local d="$1" f
  for f in "$d"/slice-*.md; do
    [ -e "$f" ] && return 0
  done
  return 1
}

# Kanonische Überschrift des Aktiv-Status-Abschnitts. EINE Wahrheit für
# heading_present (grep) und active_section (awk) — beide Muster müssen
# wörtlich gleich bleiben, sonst entsteht zwischen Guard und Extraktion
# eine Lücke (R1 zu slice-040).
HEADING_RE='^## Aktuelle Welle[[:space:]]*$'

# heading_present: die kanonische H2 existiert EXAKT. Fehlt sie (Tippfehler,
# Umbenennung, Zusatztext wie "## Aktuelle Welle (Stand …)"), kann der
# Aktiv-Status nicht bestimmt werden → der Aufrufer schlägt fehl
# (fail-closed). Ohne diesen Guard liefe active_section leer und has_active
# meldete still "aktiv" — ein silent-green, wenn zugleich ein Slice in
# in-progress/ läge (R1-MEDIUM zu slice-040).
heading_present() { # $1 roadmap
  grep -qE "$HEADING_RE" "$1"
}

# active-Block: der Abschnitt "## Aktuelle Welle" bis zur nächsten
# H2-Überschrift. Der Marker wird NUR in diesem Block gesucht, damit ein
# erklärender Verweis anderswo den Status nicht verfälscht.
active_section() { # $1 roadmap
  awk '
    /^## Aktuelle Welle[[:space:]]*$/ { inblk=1; next }
    /^## / { if (inblk) exit }
    inblk { print }
  ' "$1"
}

# hasActive: Marker NICHT im Aktuelle-Welle-Block → aktive Welle.
has_active() { # $1 roadmap
  local sec
  sec="$(active_section "$1")"
  if grep -qF "$ACTIVE_MARKER" <<<"$sec"; then
    return 1
  fi
  return 0
}

# Kernregel auf einem Paar (progress-dir, roadmap). Gibt 0 bei Konsistenz,
# 1 bei Drift. Schreibt die Befund-Meldung nach stderr.
check_pair() { # $1 progress-dir, $2 roadmap
  local d="$1" rm="$2" slices active
  if ! heading_present "$rm"; then
    echo "planning-consistency: FAIL — §\"## Aktuelle Welle\" nicht in $rm gefunden; ohne die kanonische Überschrift ist der Aktiv-Status nicht bestimmbar (fail-closed)" >&2
    return 1
  fi
  if has_slices "$d"; then slices=1; else slices=0; fi
  if has_active "$rm"; then active=1; else active=0; fi
  if [ "$slices" -eq "$active" ]; then
    return 0
  fi
  if [ "$slices" -eq 1 ]; then
    echo "planning-consistency: FAIL — $d enthält slice-*, aber $rm §Aktuelle Welle sagt \"$ACTIVE_MARKER\" (Roadmap muss die aktive Welle benennen)" >&2
  else
    echo "planning-consistency: FAIL — kein slice-* in $d, aber $rm §Aktuelle Welle benennt eine aktive Welle (Roadmap muss \"$ACTIVE_MARKER\" tragen)" >&2
  fi
  return 1
}

self_test() {
  local tmp
  tmp="$(mktemp -d)"
  local act="## Aktuelle Welle"$'\n\n'"welle-x in Arbeit."
  local idle="## Aktuelle Welle"$'\n\n'"$ACTIVE_MARKER — wartet auf Trigger."

  # Richtung A: Slice da, Roadmap idle → muss feuern.
  mkdir -p "$tmp/a"; : > "$tmp/a/slice-001-x.md"; printf '%s\n' "$idle" > "$tmp/ra.md"
  if check_pair "$tmp/a" "$tmp/ra.md" 2>/dev/null; then
    echo "planning-consistency: Selbsttest FEHLGESCHLAGEN — Slice+idle nicht erkannt" >&2
    rm -rf "$tmp"; exit 2
  fi
  # Richtung B: kein Slice, Roadmap aktiv → muss feuern.
  mkdir -p "$tmp/b"; printf '%s\n' "$act" > "$tmp/rb.md"
  if check_pair "$tmp/b" "$tmp/rb.md" 2>/dev/null; then
    echo "planning-consistency: Selbsttest FEHLGESCHLAGEN — kein-Slice+aktiv nicht erkannt" >&2
    rm -rf "$tmp"; exit 2
  fi
  # Richtung C (Heading-Guard, R1-MEDIUM): kaputte Überschrift trägt den
  # idle-Marker, ein Slice liegt vor → ohne Guard wäre das silent-green;
  # mit Guard muss es feuern.
  local broken="## Aktuelle Wellen"$'\n\n'"$ACTIVE_MARKER — wartet."
  printf '%s\n' "$broken" > "$tmp/rc.md"
  if check_pair "$tmp/a" "$tmp/rc.md" 2>/dev/null; then
    echo "planning-consistency: Selbsttest FEHLGESCHLAGEN — kaputte §Aktuelle-Welle-Überschrift nicht als fail-closed erkannt" >&2
    rm -rf "$tmp"; exit 2
  fi
  # Konsistente Gegenproben dürfen NICHT feuern.
  if ! check_pair "$tmp/a" "$tmp/rb.md" 2>/dev/null; then  # Slice+aktiv
    echo "planning-consistency: Selbsttest FEHLGESCHLAGEN — Slice+aktiv fälschlich gefeuert" >&2
    rm -rf "$tmp"; exit 2
  fi
  if ! check_pair "$tmp/b" "$tmp/ra.md" 2>/dev/null; then  # kein-Slice+idle
    echo "planning-consistency: Selbsttest FEHLGESCHLAGEN — kein-Slice+idle fälschlich gefeuert" >&2
    rm -rf "$tmp"; exit 2
  fi
  rm -rf "$tmp"
}

self_test

if [ ! -f "$ROADMAP" ]; then
  echo "planning-consistency: FAIL — $ROADMAP fehlt (fail-closed)" >&2
  exit 1
fi

if ! check_pair "$PROGRESS_DIR" "$ROADMAP"; then
  exit 1
fi

echo "planning-consistency ok: Roadmap §Aktuelle Welle ↔ in-progress/slice-* konsistent (Selbsttest gefeuert)."
