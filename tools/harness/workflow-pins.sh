#!/usr/bin/env bash
# workflow-pins — Wächter zu AGENTS.md §3.9: jeder `uses:`-Eintrag in
# .github/workflows/ nennt einen VOLLEN Commit-SHA mit Tag-Kommentar dahinter.
#
# WAS ER ZUSICHERT, GENAU: die FORM des Pins, nicht seine Gültigkeit. Ob der
# SHA existiert und ob er den Commit bezeichnet, den der Tag-Kommentar
# behauptet, ist eine Netz-Frage und ausdruecklich NICHT Teil der Zusage --
# dieser Wächter ist netzlos und laeuft in `make gates`.
#
# ER LIEST NUR ECHTE YAML-SCHLUESSEL. Die Workflow-Koepfe erklaeren die
# Pin-Konvention in Prosa und nennen `uses:` dabei mehrfach; ein Wächter, der
# auf die Zeichenkette trifft statt auf den Schluessel, meldete seine eigene
# Dokumentation. Gefiltert wird deshalb auf `^\s*(- )?uses:` -- Kommentarzeilen
# beginnen nach Leerraum mit `#` und fallen heraus.
#
# FAIL-CLOSED BEI LEERER MENGE: findet er keine Workflow-Datei oder keinen
# einzigen `uses:`-Schluessel, bricht er ROT ab statt gruen zu melden. Ein
# Wächter, der nichts gefunden hat, und einer, der nichts zu pruefen hatte,
# sehen im Exit-Code sonst identisch aus.
#
# BEIDE ENDUNGEN: GitHub liest `.yml` UND `.yaml`. Ein Glob auf nur eine von
# beiden waere derselbe stille Gruen-Pfad eine Ebene tiefer -- die Regel bindet
# an das VERZEICHNIS, nicht an eine Endung.
#
# BENANNTE GRENZE (kein Befund, sondern Falsch-Positiv-Richtung): als Pin gilt
# hier ein 40-stelliger git-SHA. Ein Digest-Verweis (`docker://image@sha256:…`)
# oder ein in Anfuehrungszeichen gesetzter Wert wuerde gemeldet, obwohl er
# gepinnt ist. Im Bestand existiert keine solche Form; taucht eine auf, ist der
# Wächter zu erweitern -- nicht die Meldung wegzunehmen.
#
# Exit: 0 = alle Pins formgerecht, 1 = mindestens ein Befund oder leere Menge.
# bash + coreutils + grep + find.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

dir=".github/workflows"
[ -d "$dir" ] || { echo "workflow-pins: ${dir} fehlt" >&2; exit 1; }

# Dateiliste getrennt bilden: `grep -H` erzwingt den Dateinamen-Praefix auch
# bei GENAU EINER Datei -- ohne ihn laesst GNU-grep ihn dann weg, und die
# Zerlegung in file/line liefe still auf falsche Felder.
mapfile -t files < <(find "$dir" -maxdepth 1 -type f \( -name '*.yml' -o -name '*.yaml' \) | LC_ALL=C sort)
if [ "${#files[@]}" -eq 0 ]; then
  echo "workflow-pins: keine Workflow-Datei in ${dir} — leere Prüfmenge, fail-closed" >&2
  exit 1
fi

# WARUM ZWEI awk-PROGRAMME UND KEIN YAML-PARSER: die Block-Form der
# `permissions:`-Bloecke traegt der Bestand, und eine vierte Toolchain waere ein
# Entscheid, kein Werkzeug nebenbei (harness/conventions.md MR-046). Was die
# Zerlegung NICHT sicher liest, gibt sie als `?<text>` aus -- der Aufrufer macht
# daraus einen Befund. Ein Waechter, der bei unbekannter Schreibweise schweigt,
# waere derselbe stille Gruen-Pfad, den diese Pruefung schliesst.
#
# NICHT GEDECKT, benannt: YAML-Anker, Flow-Style-Mappings mit Inhalt,
# Mehrfach-Dokumente und eine Job-Einrueckung ausserhalb der 2/4-Form. Alle
# vier melden sich als unlesbar, statt zu passieren.
AWK_WANTS='
{ line = $0; sub(/[[:space:]]+#.*$/, "", line); sub(/[[:space:]]+$/, "", line) }
line ~ /^[[:space:]]*#/ { next }
line == "" { next }
{
  ind = match(line, /[^ ]/) - 1
  if (line ~ /^[[:space:]]*permissions:[[:space:]]*$/)                  { inperm = 1; pind = ind; next }
  if (line ~ /^[[:space:]]*permissions:[[:space:]]*\{[[:space:]]*\}$/)  { inperm = 0; next }
  if (line ~ /^[[:space:]]*permissions:/) {
    v = line; sub(/^[[:space:]]*permissions:[[:space:]]*/, "", v); print "?" v; inperm = 0; next
  }
  if (inperm) {
    if (ind <= pind) { inperm = 0 }
    else if (line ~ /^[[:space:]]*[a-z-]+:[[:space:]]*(read|write|none)$/) {
      k = line; sub(/^[[:space:]]*/, "", k); sub(/:[[:space:]]*/, " ", k); print k; next
    }
    else { s = line; sub(/^[[:space:]]*/, "", s); print "?" s; next }
  }
}'

AWK_CALLER='
{ line = $0; sub(/[[:space:]]+#.*$/, "", line); sub(/[[:space:]]+$/, "", line) }
{
  if (line ~ /^  [A-Za-z0-9_.-]+:[[:space:]]*$/) {
    job = line; sub(/^  /, "", job); sub(/:$/, "", job); inperm = 0
  }
  if (NR == TARGET) hit = job
  if (job == "") next
  ind = match(line, /[^ ]/) - 1
  if (line ~ /^    permissions:[[:space:]]*$/)                 { inperm = 1; decl[job] = 1; next }
  if (line ~ /^    permissions:[[:space:]]*\{[[:space:]]*\}$/) { decl[job] = 1; inperm = 0; next }
  if (line ~ /^    permissions:/) {
    v = line; sub(/^    permissions:[[:space:]]*/, "", v); p[job] = p[job] "?" v "\n"; decl[job] = 1; inperm = 0; next
  }
  if (inperm) {
    if (line == "") next
    if (line !~ /^      /) { inperm = 0; next }
    if (line ~ /^      [a-z-]+:[[:space:]]*(read|write|none)$/) {
      k = line; sub(/^[[:space:]]*/, "", k); sub(/:[[:space:]]*/, " ", k); p[job] = p[job] k "\n"; next
    }
    s = line; sub(/^[[:space:]]*/, "", s); p[job] = p[job] "?" s "\n"; next
  }
}
END {
  if (hit == "") { print "!nojob"; exit }
  if (!(hit in decl)) { print "!undeclared"; exit }
  printf "%s", p[hit]
}'

# perms_level ordnet die drei GitHub-Stufen; ein Scope, den der Aufrufer nicht
# nennt, ist `none` -- das ist GitHubs eigene Semantik, kein Zugestaendnis.
perms_level() {
  case "$1" in none) echo 0 ;; read) echo 1 ;; write) echo 2 ;; *) echo -1 ;; esac
}

# perms_findings prueft EINE lokale Referenz und schreibt je Befund eine Zeile.
perms_findings() {
  local caller="$1" ln="$2" tgt="$3" wants have scope lvl hscope hlvl found
  wants="$(awk "$AWK_WANTS" "$tgt" || true)"
  # Verlangt das Ziel nichts, gibt es nichts zu vergleichen.
  [ -n "$wants" ] || return 0
  if printf '%s\n' "$wants" | grep -q '^?'; then
    printf '%s:%s\tuses-local-perms-unreadable\tRechte des Ziels ./%s nicht sicher lesbar: %s\n' \
      "$caller" "$ln" "$tgt" "$(printf '%s\n' "$wants" | grep '^?' | head -1 | cut -c2-)"
    return 0
  fi
  have="$(awk -v TARGET="$ln" "$AWK_CALLER" "$caller" || true)"
  case "$have" in
    '!nojob')
      printf '%s:%s\tuses-local-perms-unreadable\tJob-Block der Referenz nicht bestimmbar (erwartet: Job-Kopf mit zwei Leerzeichen)\n' \
        "$caller" "$ln"; return 0 ;;
    '!undeclared')
      printf '%s:%s\tuses-local-perms-undeclared\tJob traegt kein eigenes permissions:, das Ziel ./%s verlangt aber %s — der Job erbt sonst den Workflow-Kopf, und GitHub bricht den ganzen Lauf ab\n' \
        "$caller" "$ln" "$tgt" "$(printf '%s\n' "$wants" | tr '\n' ',' | sed 's/,$//')"; return 0 ;;
  esac
  if printf '%s\n' "$have" | grep -q '^?'; then
    printf '%s:%s\tuses-local-perms-unreadable\tRechte des aufrufenden Jobs nicht sicher lesbar: %s\n' \
      "$caller" "$ln" "$(printf '%s\n' "$have" | grep '^?' | head -1 | cut -c2-)"
    return 0
  fi
  while read -r scope lvl; do
    [ -n "$scope" ] || continue
    found=0
    while read -r hscope hlvl; do
      [ "$hscope" = "$scope" ] || continue
      found=1
      if [ "$(perms_level "$hlvl")" -lt "$(perms_level "$lvl")" ]; then
        printf '%s:%s\tuses-local-perms-narrow\tScope %s: Aufrufer fuehrt %s, Ziel ./%s verlangt %s\n' \
          "$caller" "$ln" "$scope" "$hlvl" "$tgt" "$lvl"
      fi
    done <<< "$have"
    if [ "$found" -eq 0 ] && [ "$(perms_level "$lvl")" -gt 0 ]; then
      printf '%s:%s\tuses-local-perms-narrow\tScope %s: Aufrufer nennt ihn nicht (= none), Ziel ./%s verlangt %s\n' \
        "$caller" "$ln" "$scope" "$tgt" "$lvl"
    fi
  done <<< "$wants"
}

# --selftest faehrt die Rechte-Pruefung gegen konstruierte Proben. Ohne sie
# waere ihre Zusage eine Erinnerung: der Anlassfall ist EIN Exemplar, und eine
# Regel, die nur an ihrem Anlass geeicht ist, traegt nicht ueber ihn hinaus
# (docs/plan/planning/observations.md BEO-011). Netzlos, kein git.
if [ "${1:-}" = "--selftest" ]; then
  tmp="$(mktemp -d)"
  trap 'rm -rf "$tmp"' EXIT
  probe_fail=0
  probe() { # name, erwartetes-Muster (leer = kein Befund), caller-yaml, ziel-yaml
    local name="$1" want="$2" cy="$3" ty="$4" got ln
    printf '%s' "$cy" > "$tmp/caller.yml"; printf '%s' "$ty" > "$tmp/target.yml"
    ln="$(grep -n 'uses: \./' "$tmp/caller.yml" | head -1 | cut -d: -f1)"
    got="$(perms_findings "$tmp/caller.yml" "$ln" "$tmp/target.yml" || true)"
    if { [ -z "$want" ] && [ -z "$got" ]; } || { [ -n "$want" ] && printf '%s' "$got" | grep -q "$want"; }; then
      echo "  ok    ${name}"
    else
      echo "  FEHL  ${name}: erwartet ${want:-<kein Befund>}, bekommen: ${got:-<kein Befund>}"
      probe_fail=$((probe_fail + 1))
    fi
  }
  C_UNDECL='permissions: {}
jobs:
  a:
    uses: ./target.yml
'
  C_READ='permissions: {}
jobs:
  a:
    permissions:
      contents: read
    uses: ./target.yml
'
  T_NONE='jobs:
  s:
    runs-on: ubuntu-latest
'
  T_READ='permissions: {}
jobs:
  s:
    permissions:
      contents: read
'
  T_WRITE='permissions: {}
jobs:
  s:
    permissions:
      contents: write
'
  T_PKG='permissions: {}
jobs:
  s:
    permissions:
      packages: write
'
  T_ALL='permissions: read-all
jobs:
  s:
    runs-on: ubuntu-latest
'
  echo "workflow-pins --selftest:"
  # Der reale Ausfall von v0.66.0, nachgebildet.
  probe "stille Vererbung (der v0.66.0-Fall)" "uses-local-perms-undeclared" "$C_UNDECL" "$T_READ"
  probe "deklariert und ausreichend"          ""                            "$C_READ"   "$T_READ"
  probe "zu eng: read gegen write"            "uses-local-perms-narrow"     "$C_READ"   "$T_WRITE"
  probe "Scope fehlt beim Aufrufer"           "uses-local-perms-narrow"     "$C_READ"   "$T_PKG"
  probe "Ziel verlangt nichts"                ""                            "$C_UNDECL" "$T_NONE"
  probe "unlesbare Form (read-all)"           "uses-local-perms-unreadable" "$C_READ"   "$T_ALL"
  # Ein Ziel ohne Forderung macht auch die fehlende Deklaration harmlos.
  probe "keine Forderung, keine Deklaration"  ""                            "$C_UNDECL" "$T_NONE"
  if [ "$probe_fail" -gt 0 ]; then
    echo "workflow-pins --selftest: ${probe_fail} Probe(n) fehlgeschlagen" >&2
    exit 1
  fi
  echo "workflow-pins --selftest: 7 Probe(n) ok"
  exit 0
fi

findings=0
checked=0
local_refs=0

while IFS= read -r hit; do
  file="${hit%%:*}"
  rest="${hit#*:}"
  line="${rest%%:*}"
  text="${rest#*:}"
  checked=$((checked + 1))

  # LOKALE Workflow-Referenz (`uses: ./.github/workflows/x.yml`): sie kann
  # keinen SHA tragen und BRAUCHT keinen. Sie loest auf denselben Commit auf
  # wie der aufrufende Workflow und ist damit staerker gebunden als ein
  # SHA-Pin -- sie kann per Konstruktion nicht driften. Der Zweck von §3.9
  # (keine beweglichen Referenzen) ist hier ohne Pin erfuellt.
  #
  # STATT DES PIN-CHECKS die Frage, die hier Sinn ergibt: EXISTIERT das Ziel?
  # Ein vertippter lokaler Verweis traegt keinen Pin und faellt sonst erst zur
  # Laufzeit auf. Die Ausnahme nimmt damit eine Pruefung weg und setzt eine
  # andere an ihre Stelle -- sie laesst keinen Eintrag ungeprueft.
  #
  # ABGRENZUNG: das gilt NUR fuer den `./`-Praefix. Eine Referenz auf einen
  # anderen Repo-Pfad (`owner/repo/.github/workflows/x.yml@ref`) ist fremd und
  # faellt unter die Regel wie jede Action.
  if printf '%s' "$text" | grep -qE 'uses:[[:space:]]*\./'; then
    local_refs=$((local_refs + 1))
    target="$(printf '%s' "$text" | sed -nE 's|.*uses:[[:space:]]*\./([^[:space:]]+).*|\1|p')"
    if [ -z "$target" ] || [ ! -f "$target" ]; then
      echo "${file}:${line}	uses-local-missing	lokale Workflow-Referenz zeigt auf keine Datei: ./${target}"
      findings=$((findings + 1))
      continue
    fi
    # DIE ZWEITE FRAGE AN DIESELBE REFERENZ (slice-169): hat der aufrufende Job
    # die Rechte, die das Ziel verlangt? Ein aufgerufener Workflow bekommt nur,
    # was der Aufrufer selbst fuehrt; verlangt er mehr, lehnt GitHub den GANZEN
    # Lauf vor dem ersten Job ab (startup_failure, kein Log). Bis hierher las
    # dieser Waechter die `uses:`-ZEILE und oeffnete die Zieldatei, ohne sie
    # anzusehen -- genau die Luecke, die AGENTS.md §3.8 als Frage stellt.
    n="$(perms_findings "$file" "$line" "$target")"
    if [ -n "$n" ]; then
      printf '%s\n' "$n"
      findings=$((findings + $(printf '%s\n' "$n" | grep -c '')))
    fi
    continue
  fi

  # Form: <ref>@<40 Hex> gefolgt von einem Kommentar, der den Tag nennt.
  if ! printf '%s' "$text" | grep -qE '@[0-9a-f]{40}([[:space:]]|$)'; then
    echo "${file}:${line}	uses-pin-missing	kein voller 40-stelliger Commit-SHA"
    findings=$((findings + 1))
    continue
  fi
  if ! printf '%s' "$text" | grep -qE '@[0-9a-f]{40}[[:space:]]+#[[:space:]]*\S'; then
    echo "${file}:${line}	uses-pin-untagged	SHA ohne Tag-Kommentar dahinter"
    findings=$((findings + 1))
  fi
done < <(printf '%s\n' "${files[@]}" | xargs -r grep -HnE '^[[:space:]]*(- )?uses:[[:space:]]' || true)

if [ "$checked" -eq 0 ]; then
  echo "workflow-pins: kein einziger uses:-Schlüssel gefunden — leere Prüfmenge, fail-closed" >&2
  exit 1
fi

if [ "$findings" -gt 0 ]; then
  echo "workflow-pins: ${findings} Befund(e) über ${checked} uses:-Einträge (AGENTS.md §3.9)" >&2
  exit 1
fi

echo "workflow-pins: ok (${checked} uses:-Einträge geprüft, davon ${local_refs} lokale Referenz(en) ohne Pin-Pflicht; alle übrigen SHA-gepinnt mit Tag-Kommentar)"
