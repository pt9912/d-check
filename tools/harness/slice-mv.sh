#!/usr/bin/env bash
# slice-mv.sh — Lifecycle-Wechsel eines Slice UND der Verweise, die er bricht,
# emittiert von ai-harness-init.
#
# ZUSAGE. `make slice-mv SLICE=slice-<Kennung> TO=<open|next|in-progress|done>`
# bewegt den Slice per `git mv` und committet den reinen Move SOFORT als eigenen
# Commit (ein Move und eine Inhaltsaenderung in einem Commit verlieren die
# Rename-Erkennung). Danach zieht es reale Verweise nach — EINGEHEND (jede
# Praefix-Form auf die bewegte Datei, repo-weit, dazu die praefixlose
# Link-Form aus den Geschwistern im Ausgangsverzeichnis) UND AUSGEHEND (praefixlose
# Ziele INNERHALB der bewegten Datei, die nach dem Wechsel ins falsche
# Verzeichnis zeigen; getroffen wird eine nummerierte wie eine benannte
# Slice-Kennung gleichermassen) — und committet diese Inhaltsaenderung, falls
# welche anfiel, als ZWEITEN, vom Move getrennten Commit; fiel keine an, bleibt
# es beim einen Move-Commit.
#
# WAS DIESE COMMITS TRAGEN. Beide Messages nennen den bewegten Slice mit seinem
# Dateinamen, und der Dateiname IST die Kennung eines Slice: die Messages tragen
# sie damit verbatim. Ob eine Kennungs-Menge dieses Repos diese Form trifft, ist
# eine andere Achse — dieses Werkzeug fuehrt keine Kennungs-Menge und setzt keine
# voraus, und es prueft die Message nicht. Setzt das Repo einen Waechter ein, der
# einen anderen Zuschnitt prueft, fallen die zwei Commits NACH dem `git mv` durch
# ihn; was dann liegen bleibt, ist ein gestagter Rename.
#
# VORAUSSETZUNG. Weil dieses Skript selbst committet, verlangt es einen sauberen
# Arbeitsbaum (keine gestagten oder ungestagten Aenderungen an getrackten
# Dateien), BEVOR es startet — sonst landet ein fremder, zufaellig anwesender
# Diff in einem der beiden Commits. Ein Verstoss bricht den Aufruf vor dem
# ersten `git mv` ab (main(), erste Pruefung).
#
# REPO-POLITIK statt Mechanik: die Pfade, die der EINGEHEND-Nachzug ausnimmt
# (Variable SLICE_MV_AUSGENOMMENE_PFADE unten). Die zwei Vorgaben sind
# Ableitungen aus dem mitemittierten Regelwerk, keine Vorliebe des Werkzeugs:
# `.harness/baseline` ist unveraenderter Fremdtext, und eine `Accepted`-ADR
# wird nach der Hard Rule fuer Accepted-ADRs nicht inhaltlich ueberschrieben.
# Diese Datei wird bei jedem Bootstrap kanonisch neu geschrieben, eine Aenderung
# an ihr ueberlebt den naechsten Lauf nicht; gesetzt wird die Variable darum von
# aussen (dieselbe im mitemittierten Make-Fragment, s. dort) — der Ort dafuer ist
# jede Make-Quelle oder die Umgebung des Aufrufs.
#
# GRENZEN (vier, jede mit ihrer Ursache):
# (1) Das Werkzeug zieht PFADE nach, keine ZUSTANDSSAETZE. Eine Zeile
#     "In Arbeit: <slice>" bleibt nach dem Wechsel stehen; ihr Verweis wird
#     richtig, ihre Aussage falsch. Welcher Satz einen Zustand behauptet, ist
#     Urteil, kein Match.
# (2) Ein WELLE-Plan wechselt beim Closure-Move die Verzeichnis-TIEFE (flach ->
#     done/), nicht nur das Verzeichnis — eine andere Ersetzung als der Tausch
#     auf gleicher Ebene. Dieses Werkzeug bewegt nur SLICE-Dateien
#     (SLICE=slice-<Kennung>) und ersetzt in der Ausgehend-Richtung darum auch
#     nur "slice-"-Ziele; ein praefixloses "welle-"-Ziel bleibt unberuehrt.
# (3) Die praefixlose EINGEHEND-Ersetzung (rewrite_incoming_bare_in_file)
#     erkennt einen Verweis ohne Verzeichnis-Segment nur als Markdown-Link
#     "](<datei>)" oder "](<datei>#…)" und nur in den getrackten Dateien, die
#     flach im $from-Verzeichnis liegen — dort loest der blanke Name gegen das
#     Verzeichnis auf, das die Datei verlassen hat. Eine andere Schreibweise
#     desselben Verweises ("](./<datei>)", "](<<datei>>)", eine
#     Referenz-Definition "[x]: <datei>") bleibt stehen; was danach tot ist,
#     meldet das Doku-Gate des Repos als toten Link. Markdown liest die
#     Ersetzung nicht: steht die Link-Syntax selbst mit genau diesem Namen in
#     einem Code-Span oder Code-Block, wird sie mitersetzt.
# (4) Die AUSGEHEND-Ersetzung trifft die lowercase-Kebab-Form einer benannten
#     Kennung (Zeichenklasse "[0-9a-z]"). Eine Slice-Kennung, die das Praefix
#     eines vorhandenen Ankers traegt (LH-*, ADR-*, CO-*) und darum
#     Grossbuchstaben fuehrt, bleibt unerkannt.
set -euo pipefail

PLANNING="docs/plan/planning"
LIFECYCLE="open next in-progress done"

# psed_i — portables `sed -i`: BSD-sed (macOS) verlangt nach `-i` zwingend eine eigene
# Backup-Extension als naechstes Token (auch leer) und verschluckt sonst das naechste
# Argument dafuer — ein blosses `sed -i -E SCRIPT FILE` (Extension = "-E", -E greift
# nicht) scheitert dort mit "\1 not defined in the RE", auf GNU-sed nicht. Kein -i:
# Ausgabe in eine temporaere Datei, dann in die Zieldatei GESCHRIEBEN statt ueber sie
# verschoben — `mktemp` legt die temporaere Datei mit 0600 an, und `mv` traegt diesen
# Modus auf das Ziel; ein `d-check`-Container liest als Nicht-Root, und ein derart auf
# 0600 gefallenes Ziel wird dort unlesbar. `cat >` in die bestehende Zieldatei behaelt
# deren Inode und damit ihren Modus. Aufruf wie `sed -i`: optionale Flags, dann SCRIPT,
# dann FILE als letztes Argument.
psed_i() {
  local tmp ziel
  tmp="$(mktemp -p "${TMPDIR:-/tmp}")"
  ziel="${!#}"
  sed "$@" >"$tmp"
  cat "$tmp" >"$ziel"
  rm -f "$tmp"
}

# REPO-POLITIK (s. Kopf): die `git grep`-Pathspecs, die der EINGEHEND-Nachzug
# NICHT durchsucht — whitespace-getrennt, eine Angabe je Eintrag. Die Vorgabe
# ist aus dem mitemittierten Regelwerk abgeleitet (Fremdtext und die Hard Rule
# fuer Accepted-ADRs). Ein Repo, das weitere Baeume ausnimmt oder Zeitdokumente
# anders behandelt, setzt diese Variable in der Umgebung des Aufrufs oder in
# einer Make-Quelle — das mitemittierte Fragment reicht seinen Wert als Umgebung
# an dieses Skript durch. Der Vorgabewert hier greift nur, wenn von aussen keiner
# gesetzt ist.
SLICE_MV_AUSGENOMMENE_PFADE="${SLICE_MV_AUSGENOMMENE_PFADE:-:!.harness/baseline :!docs/plan/adr}"

# d-check-ANPASSUNG · seit slice-245: done/ traegt Lifecycle-Unterordner (wellenlos
# fuer wellenlose Slices, welle-NN fuer Wellen-Slices). Die Variable
# SLICE_MV_DONE_UNTERORDNER mappt TO=done auf den Unterordner — gesetzt in der
# Make-Quelle dieses Repos. Ohne sie gilt die Schwester-Form (flaches done/).
SLICE_MV_DONE_UNTERORDNER="${SLICE_MV_DONE_UNTERORDNER:-}"

usage() {
  cat >&2 <<'USAGE'
Aufruf: make slice-mv SLICE=slice-<Kennung>[-kurztitel[.md]] TO=<open|next|in-progress|done>

  Bewegt den Slice per `git mv`, committet den reinen Move sofort, und zieht
  danach die Verweise nach — repo-weit eingehend und innerhalb der Datei selbst
  ausgehend; fielen Verweise an, committet es sie getrennt vom Move. Verlangt
  einen sauberen Arbeitsbaum. Grenzen: siehe Skriptkopf.
USAGE
}

# Erweiterte Regex-Metazeichen im Dateinamen entschaerfen — ein Slice-Titel
# traegt mindestens einen Punkt (die Endung ".md"), der in ERE sonst "ein
# beliebiges Zeichen" bedeutet statt sich selbst.
re_escape() {
  printf '%s' "$1" | sed -e 's/[][\.^$*+?(){}|\\]/\\&/g'
}

# Die Ausnahmeliste als reine Funktion — main() liest sie hier aus, und ein
# Test kann sie rufen, ohne ein Repo zu bewegen. Zerlegt wird an Leerraum, weil
# die Variable als eine Zeile mit mehreren Eintraegen gesetzt wird.
eingehend_ausgenommene_pfade() {
  local p
  for p in $SLICE_MV_AUSGENOMMENE_PFADE; do printf '%s\n' "$p"; done
}

# EINGEHEND: jedes Vorkommen von "$from/$base" in $file wird zu "$to/$base" —
# an einer Wortgrenze (Zeilenanfang oder ein Zeichen davor, das kein Bestandteil
# eines Wortes/Verzeichnisnamens ist — Buchstabe, Ziffer, Unterstrich UND
# Bindestrich zaehlen als Wortzeichen, weil "in-progress" selbst einen
# Bindestrich traegt und ein glued Praefix wie "sibling-open/" sonst faelschlich
# traefe), nicht an einer festen Praefix-Liste. Das deckt jede Praefix-Tiefe
# (vom nackten Verzeichnisnamen bis zum doppelten Aufstieg ueber
# "docs/plan/planning/") mit einer Regel statt einer Liste, die driftet.
rewrite_incoming_in_file() {  # $1=datei $2=base $3=from $4=to
  local file="$1" base="$2" from="$3" to="$4" esc_base
  esc_base="$(re_escape "$base")"
  psed_i -E "s#(^|[^A-Za-z0-9_-])$from/$esc_base#\\1$to/$base#g" "$file"
}

# EINGEHEND, PRAEFIXLOS: jeder Markdown-Link "](<base>)" oder "](<base>#…)" in
# $file wird zu "](../<to>/<base>…)". main() ruft das fuer die Geschwister im
# Ausgangsverzeichnis der bewegten Datei auf — dort loest der blanke Name gegen
# das Verzeichnis auf, das die Datei gerade verlassen hat. Die Regel ankert an
# der Link-Klammer "](" und am Ende des Namens (")" oder "#"): ein Code-Span mit
# dem blossen Namen, ein Tree-Operand "<sha>:<base>", ein Verweis mit
# Verzeichnis-Segment und ein laengerer Name mit demselben Anfang bleiben
# stehen. Die Regel liest kein Markdown: steht die Link-Syntax selbst mit
# genau diesem Namen in einem Code-Span oder Code-Block, wird sie mitersetzt
# (Grenze 3 im Skriptkopf). Gibt die Anzahl ersetzter Links auf stdout aus.
rewrite_incoming_bare_in_file() {  # $1=datei $2=base $3=to
  local file="$1" base="$2" to="$3" esc_base count
  esc_base="$(re_escape "$base")"
  count="$( { grep -oE "[]]\\(${esc_base}[)#]" "$file" 2>/dev/null || true; } | wc -l)"
  psed_i -E "s|[]]\\($esc_base([)#])|](../$to/$base\\1|g" "$file"
  printf '%d\n' "$((count))"
}

# AUSGEHEND: praefixlose "](slice-…)"-Ziele INNERHALB von $file, deren Datei im
# $from-Verzeichnis liegen geblieben ist, bekommen die Stufen-Kette aus $3
# ("../" je Segment) plus "$from/" vorangestellt — sonst zeigt der Verweis nach
# dem Wechsel ins neue (falsche) Verzeichnis. Die Stufen-Kette traegt die Tiefe
# des neuen Orts: aus "done/welle-91" heraus steigt ein Verweis mit zwei Stufen
# ("../../in-progress/X") auf das planning-Niveau, aus "next" heraus reicht
# eine ("../in-progress/X").
# Nur "slice-"-Ziele (Grenze 2 im Skriptkopf); ein Ziel, das nicht (mehr) unter
# $from liegt, bleibt unberuehrt (kein Rateversuch, welches Verzeichnis stimmt).
# Das Fundmuster trifft eine nummerierte Kennung (slice-NNN…) ebenso wie eine
# benannte (slice-<slug>, lowercase Kebab-Case ohne Ziffern-Praefix).
# Gibt die Anzahl umgehaengter Ziele auf stdout aus — main() liest sie per
# Kommando-Substitution, statt Vorher/Nachher getrennt zu zaehlen.
rewrite_outgoing_bare_in_file() {  # $1=datei $2=from $3=stufen
  local file="$1" from="$2" stufen="$3" t esc_t count=0
  while IFS= read -r t; do
    [ -n "$t" ] || continue
    [ -f "$PLANNING/$from/$t" ] || continue
    esc_t="$(re_escape "$t")"
    psed_i -E "s#\\]\\($esc_t\\)#]($stufen$from/$t)#g" "$file"
    count=$((count + 1))
  done < <(grep -ohE '\]\(slice-[0-9a-z][^)/]*\)' "$file" 2>/dev/null \
             | sed -E 's/^\]\(//; s/\)$//' | sort -u)
  printf '%d\n' "$count"
}

main() {
  local SLICE="${1:-}" TO="${2:-}"
  [ -n "$SLICE" ] && [ -n "$TO" ] || { usage; exit 2; }

  cd "$(dirname "$0")/../.."

  # Sauberer Arbeitsbaum (VORAUSSETZUNG im Skriptkopf) — sonst landet ein
  # fremder Diff in einem der beiden automatischen Commits weiter unten.
  if ! git diff --quiet || ! git diff --cached --quiet; then
    echo "slice-mv: Arbeitsbaum nicht sauber — erst committen oder stashen (das Skript committet selbst, siehe Skriptkopf VORAUSSETZUNG)" >&2
    exit 2
  fi

  case " $LIFECYCLE " in
    *" $TO "*) ;;
    *) echo "slice-mv: '$TO' ist kein Lifecycle-Verzeichnis ($LIFECYCLE)" >&2; exit 2 ;;
  esac

  # d-check-Anpassung: TO=done auf den Lifecycle-Unterordner mappen (Kopf).
  local ziel="$TO"
  if [ "$TO" = "done" ] && [ -n "${SLICE_MV_DONE_UNTERORDNER:-}" ]; then
    ziel="$TO/$SLICE_MV_DONE_UNTERORDNER"
  fi

  # Quelle finden: Praefix oder voller Dateiname, in genau EINEM Verzeichnis —
  # zwei Treffer (auch ueber Verzeichnisse hinweg) sind mehrdeutig und brechen
  # ab, statt zu raten.
  local found="" d f
  for d in $LIFECYCLE; do
    for f in "$PLANNING/$d/${SLICE%.md}"*.md; do
      [ -e "$f" ] || continue
      if [ -n "$found" ]; then
        echo "slice-mv: '$SLICE' ist mehrdeutig — $found und $f" >&2
        exit 2
      fi
      found="$f"
    done
  done
  [ -n "$found" ] || { echo "slice-mv: kein Slice '$SLICE' unter $PLANNING/" >&2; exit 2; }

  local base from
  base="$(basename "$found")"
  from="$(basename "$(dirname "$found")")"
  if [ "$from" = "$TO" ]; then
    echo "slice-mv: '$base' liegt bereits in $from/" >&2
    exit 2
  fi

  mkdir -p "$PLANNING/$ziel"
  git mv "$found" "$PLANNING/$ziel/$base"

  # Identitaet-Fallback: ein frischer Klon traegt keine git-Identitaet — die
  # Commits des Skripts bringen sie selbst mit (dieselbe Form wie die
  # Selbstpruefung), ohne eine konfigurierte Identitaet zu ueberschreiben.
  # Die Expansion ${ident[@]+…} haelt das leere Array unter bash < 4.4 fuer
  # `set -u` unschaedlich (dort ist die Expansion eines leeren Arrays ein
  # unbound-variable-Fehler), auf jungen bash ist sie ohne Wirkung.
  local -a ident=()
  if [ -z "$(git config user.email || true)" ]; then
    ident=(-c user.email=d-check@local -c user.name="d-check slice-mv")
  fi

  # Commit 1 — reiner Move, kein Byte Inhalt veraendert: der Arbeitsbaum war
  # laut Vorpruefung sauber, `git mv` ist die einzige gestagte Aenderung, also
  # committet dieser Aufruf genau sie.
  git ${ident[@]+"${ident[@]}"} commit -q -m "slice-mv: $base  $from/ -> $ziel/ (reiner Move)"

  # EINGEHEND, repo-weit — ausser der Liste aus eingehend_ausgenommene_pfade().
  # Zeitdokumente sind NICHT ausgenommen: `done/`, `docs/reviews/` und die
  # Belege des Registers sind reale Verweisziele, auf die das Doku-Gate sieht;
  # ein Verweis auf die bewegte Datei bricht dort genauso wie ueberall sonst.
  # Nur der Pfad aendert sich, die umgebende Aussage bleibt stehen (Grenze 1).
  local -a in_pathspec=()
  while IFS= read -r p; do in_pathspec+=("$p"); done < <(eingehend_ausgenommene_pfade)

  local in_count=0 rf
  local -a touched=()
  while IFS= read -r rf; do
    [ -n "$rf" ] || continue
    rewrite_incoming_in_file "$rf" "$base" "$from" "$ziel"
    touched+=("$rf")
    in_count=$((in_count + 1))
  done < <(git grep -l -F -e "$from/$base" -- "${in_pathspec[@]}" 2>/dev/null || true)

  # EINGEHEND, praefixlos — in den getrackten Geschwistern, die flach im
  # Ausgangsverzeichnis liegen (":(glob)" haelt "*" innerhalb eines Segments),
  # unter derselben Ausnahmeliste. Eine Datei, die die Praefix-Ersetzung schon
  # getroffen hat, zaehlt in $in_count nicht doppelt.
  local bare_count=0 n sf
  while IFS= read -r sf; do
    [ -n "$sf" ] || continue
    n="$(rewrite_incoming_bare_in_file "$sf" "$base" "$ziel")"
    [ "$n" -gt 0 ] || continue
    bare_count=$((bare_count + n))
    case " ${touched[*]-} " in
      *" $sf "*) ;;
      *) touched+=("$sf"); in_count=$((in_count + 1)) ;;
    esac
  done < <(git grep -l -F -e "]($base" -- ":(glob)$PLANNING/$from/*.md" "${in_pathspec[@]}" 2>/dev/null || true)

  # AUSGEHEND — nur in der bewegten Datei selbst, an ihrem NEUEN Ort. Die
  # Stufen-Kette ("../" je Segment von $ziel) haelt die Ersetzung relativ zum
  # neuen Ort richtig — auch dann, wenn $ziel ein Unterordner ist
  # (done/welle-91: zwei Stufen, ../../$from/).
  local out_count stufen="" rest="$ziel"
  while [ -n "$rest" ]; do
    stufen="$stufen../"
    case "$rest" in */*) rest="${rest#*/}" ;; *) rest="" ;; esac
  done
  out_count="$(rewrite_outgoing_bare_in_file "$PLANNING/$ziel/$base" "$from" "$stufen")"
  [ "$out_count" -gt 0 ] && touched+=("$PLANNING/$ziel/$base")

  # Commit 2 — Inhaltsaenderung, GETRENNT vom Move, nur wenn ueberhaupt ein
  # Verweis anfiel; explizite Pfade statt `git add -A`, damit kein anderer
  # (eigentlich schon per VORAUSSETZUNG ausgeschlossener) Diff mitgenommen wird.
  if [ "${#touched[@]}" -gt 0 ]; then
    git add -- "${touched[@]}"
    git ${ident[@]+"${ident[@]}"} commit -q -m "slice-mv: Verweise auf $base nach $ziel/ nachgezogen ($in_count eingehend, $out_count ausgehend, $bare_count praefixlos aus $from/)"
  fi

  echo "slice-mv ok: $base  $from/ -> $ziel/"
  echo "  Commit 1 (reiner Move): $from/$base -> $ziel/$base"
  echo "  eingehend: $in_count Datei(en) mit Verweisen nachgezogen, darin $bare_count praefixlose(r) Link(s) aus Geschwistern unter $from/"
  echo "  ausgehend: $out_count praefixloses Ziel(e) in der bewegten Datei auf ../$from/ umgehaengt"
  if [ "${#touched[@]}" -gt 0 ]; then
    echo "  Commit 2 (Inhalt, getrennt vom Move): $in_count eingehend, $out_count ausgehend, $bare_count praefixlos aus $from/"
  else
    echo "  Kein Verweis zu ziehen — kein zweiter Commit noetig."
  fi
}

# BASH_SOURCE-Waechter: ein Test sourced dieses Skript, um
# die Ersetzungs-Funktionen direkt zu pruefen,
# ohne main() (und damit git mv) auszuloesen — sonst misst der Test sich selbst
# statt der Ersetzung.
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
  main "$@"
fi
