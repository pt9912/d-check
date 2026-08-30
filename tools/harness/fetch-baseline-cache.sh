#!/usr/bin/env bash
# fetch-baseline-cache — materialisiert die lokale, netzlose Lese-Form der
# adoptierten Baseline (AI-Harness-Kurs) als committet-vendored Baseline.
# Seit dem Kurs-v5.0.0-Release liefert EIN self-contained Bundle
# (lab-regelwerk.zip) BEIDE Bäume; lab-templates.zip existiert nicht mehr:
#
#   regelwerk → .harness/baseline/<tag>/regelwerk/   (COMMITTET, vendored)
#   templates → .harness/baseline/<tag>/templates/   (COMMITTET, vendored)
#               + .harness/baseline/<tag>/SHA256SUMS  (Manifest über BEIDE Bäume)
#
# Beide Bäume vendored: das Regelwerk verweist mit ../templates/… auf die
# Templates als „Ziel-Form" — netzlos nur auflösbar, wenn templates/ parallel zu
# regelwerk/ vendored ist. Der frühere Cache-Zweig (lab-templates.zip nach
# .harness/cache/) entfällt mit dem self-contained Bundle (Baseline-Migration
# v1.4.0 → v5.0.0, slice-084/Etappe A).
#
# Modi:
#   (default)       re-vendor: zieht lab-regelwerk.zip, entpackt beide Bäume
#                   TOLERANT (regelwerk am modul-00-Marker, templates als
#                   Geschwister) in den committeten Vendor-Pfad, prüft die
#                   Under-Copy-Barriere, (re)generiert SHA256SUMS über beide
#                   Bäume und verifiziert. Netz nötig — Anlass: Baseline-Bump.
#   --verify        Integritätsprüfung des committeten Bestands: sha256sum -c,
#                   Manifest-Deckung UND Aufloesung der Aliase unter
#                   .claude/rules/ (MR-055). Offline, kein Netz — CI/Audit.
#   --selftest      Proben der Alias-Aufloesung (neun Faelle, netzlos, ohne
#                   Wirkung auf das Repo) — `make baseline-probe`.
#   --check-latest  Upstream-Audit (Netz, informativ, KEIN Gate, KEIN
#                   fail-closed). Zwei Teile:
#                   (A) Currency: die Release-LISTE des Kurs-Repos gegen den Pin
#                       (exit 3 = neuerer Release-Tag; kein Auto-Update).
#                   (B) Content-Drift am GEPINNTEN Tag: die Bytes BEIDER Bäume des
#                       Release-Assets gegen das committete SHA256SUMS (exit 4 =
#                       Tag verschoben / Asset neu). Netz-/Werkzeug-/Manifest-
#                       Ausfall → SKIP je Teil. Exit = schlimmster Fall (4>3>0).
#
# Integrität ist nicht Aktualität: --verify beantwortet „ist der vendorte Bestand
# unversehrt?", --check-latest „ist der gepinnte Stand noch aktuell & authentisch?".
#
# Tag-Quelle: ohne Argument die §Baseline-`**Stand:**`-Zeile in
# harness/conventions.md; mit Argument ein expliziter Tag (z. B. v5.0.0) — der
# erste Vendor-Lauf einer Migration nutzt das explizite Argument, weil der Pin in
# conventions.md erst im selben Bogen umgestellt wird.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
root="$(pwd -P)"

repo="pt9912/ai-harness-course"
conventions="harness/conventions.md"

mode="vendor"
case "${1:-}" in
  --verify)       mode="verify"; shift ;;
  --selftest)     mode="selftest"; shift ;;
  --check-latest) mode="check-latest"; shift ;;
esac

tag="${1:-}"
if [ -z "$tag" ]; then
  tag="$(grep -m1 '\*\*Stand:\*\*' "$conventions" 2>/dev/null \
    | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' | head -1 || true)"
fi
if ! [[ "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "fetch-baseline-cache: ungültiger/leerer Tag '${tag}' — Argument vMAJOR.MINOR.PATCH angeben oder §Baseline in ${conventions} prüfen" >&2
  exit 1
fi

baseline=".harness/baseline/${tag}"
sums="${baseline}/SHA256SUMS"

# extract_both_trees <zip> <dest> — entpackt das self-contained Bundle TOLERANT
# nach <dest>/{regelwerk,templates}/ (regelwerk am modul-00-Marker erkannt, egal
# ob flach oder verschachtelt; templates als Geschwister gefordert) und setzt die
# globale src_n = Quell-Dateizahl beider Bäume (für die Under-Copy-Barriere).
extract_both_trees() {
  local zip="$1" dest="$2" stage
  stage="$(mktemp -d)"
  unzip -oq "$zip" -d "$stage"
  local modul0; modul0="$(find "$stage" -type f -name 'modul-00-*.md' | head -1)"
  [ -n "$modul0" ] \
    || { echo "fetch-baseline-cache: kein Regelwerk (modul-00-*.md) im ZIP gefunden" >&2; rm -rf "$stage"; return 1; }
  local src_regelwerk archive_root src_templates
  src_regelwerk="$(dirname "$modul0")"
  archive_root="$(dirname "$src_regelwerk")"
  src_templates="${archive_root}/templates"
  [ -d "$src_templates" ] \
    || { echo "fetch-baseline-cache: kein templates/-Geschwisterbaum neben regelwerk/ gefunden (das self-contained Bundle liefert beide als Geschwister)" >&2; rm -rf "$stage"; return 1; }
  # Ziel absolut auflösen: ein repo-relatives dest (Vendor-Modus) wird über
  # $root verankert, weil die Kopier-Subshell das Verzeichnis wechselt; ein
  # bereits absolutes dest (check-latest: mktemp) bleibt unverändert — sonst
  # landet die Kopie unter <repo>/<mktemp-Pfad> und der Vergleich läuft leer.
  local tgt
  case "$dest" in /*) tgt="$dest" ;; *) tgt="${root}/${dest}" ;; esac
  rm -rf "${tgt}/regelwerk" "${tgt}/templates"
  mkdir -p "${tgt}/regelwerk" "${tgt}/templates"
  # Regelwerk-Markdown (README + grundlagen-* + modul-*) REKURSIV, Struktur
  # erhalten (nicht maxdepth-1: ein künftiger Reorg mit Unterverzeichnissen darf
  # nicht still unter-kopiert werden).
  ( cd "$src_regelwerk" && find . -type f -name '*.md' -print0 | while IFS= read -r -d '' f; do
      d="${tgt}/regelwerk/${f#./}"; mkdir -p "$(dirname "$d")"; cp "$f" "$d"; done )
  ( cd "$src_templates" && find . -type f -print0 | while IFS= read -r -d '' f; do
      d="${tgt}/templates/${f#./}"; mkdir -p "$(dirname "$d")"; cp "$f" "$d"; done )
  src_n="$(( $(find "$src_regelwerk" -type f -name '*.md' | wc -l) + $(find "$src_templates" -type f | wc -l) ))"
  rm -rf "$stage"
}

# check_aliases ist die DRITTE Frage von verify(): Aliase IN den gepinnten Baum.
# Ein Symlink unter .claude/rules/ bindet denselben Pin — er wird aber von keinem
# Modul gescannt und steht in keiner Manifest-Zeile. Beim Bump braeche er STILL,
# und weder sha256sum -c noch die Deckungszaehlung saehe es (MR-055).
#
# REKURSIV und DOTFILE-BEWUSST: die Zusage lautet "jeder Symlink UNTERHALB von
# .claude/rules/". Ein flacher Glob liesse ein Unterverzeichnis und jeden
# Punkt-Namen still passieren — genau die Klasse, gegen die diese Frage steht.
#
# GRENZEN, benannt: geprueft wird die AUFLOESUNG, nicht das Ziel (ein Alias auf
# eine Datei ausserhalb des gepinnten Baums passiert); ein FEHLENDES oder leeres
# .claude/rules/ ist von "hier gibt es keine Aliase" nicht unterscheidbar und
# meldet nicht.
check_aliases() {
  local root="${1:-.claude/rules}" dangling=0 l tgt
  [ -d "$root" ] || return 0
  while IFS= read -r l; do
    [ -n "$l" ] || continue
    if tgt="$(readlink -e "$l" 2>/dev/null)" && [ -n "$tgt" ]; then continue; fi
    if [ -e "$l" ]; then
      echo "fetch-baseline-cache: unaufloesbarer Symlink ${l} — Zyklus oder zu tiefe Kette" >&2
    else
      echo "fetch-baseline-cache: toter Symlink ${l} — Ziel fehlt (Baseline-Bump?)" >&2
    fi
    dangling=1
  done < <(find "$root" -type l 2>/dev/null | LC_ALL=C sort)
  [ "$dangling" = 0 ]
}

# selftest faehrt check_aliases gegen einen eigenen Baum unter TMPDIR: je Fall
# Erwartung und Ergebnis, Fehlschlag-Zaehler am Ende. Netzlos, ohne Wirkung auf
# das Repo.
#
# WARUM: eine Zusage ohne wiederholbare Probe ist eine Erinnerung — dieselbe
# Begruendung, die `make guard-probe` traegt. Die Faelle sind die, an denen der
# flache Glob der ersten Fassung still passierte (MR-055).
selftest() {
  local d fails=0
  d="$(mktemp -d "${TMPDIR:-/tmp}/fbc-selftest.XXXXXX")" || exit 1
  # shellcheck disable=SC2064
  trap "rm -rf '$d'" EXIT

  probe() { # name erwartung(ok|rot) verzeichnis
    local name="$1" want="$2" dir="$3" got
    if check_aliases "$dir" >/dev/null 2>&1; then got="ok"; else got="rot"; fi
    if [ "$got" = "$want" ]; then
      printf 'fetch-baseline-cache: selftest OK   %-34s erwartet=%s\n' "$name" "$want"
    else
      printf 'fetch-baseline-cache: selftest FEHL %-34s erwartet=%s ergebnis=%s\n' "$name" "$want" "$got"
      fails=$((fails + 1))
    fi
  }

  mkdir -p "$d/ziel" "$d/leer" "$d/flach" "$d/tief/sub" "$d/punkt" "$d/dir" "$d/zyklus" "$d/echt"
  : > "$d/ziel/da.md"

  ln -s ../ziel/da.md            "$d/flach/gut.md"
  ln -s ../ziel/weg.md           "$d/tief/tot.md"
  ln -s ../../ziel/weg.md        "$d/tief/sub/tot.md"
  ln -s ../ziel/weg.md           "$d/punkt/.versteckt.md"
  ln -s ../ziel                  "$d/dir/als-verzeichnis"
  ln -s a.md                     "$d/zyklus/b.md"
  ln -s b.md                     "$d/zyklus/a.md"
  printf 'nur text\n'          > "$d/echt/keine-datei.md"

  probe "gesunder Alias"              ok  "$d/flach"
  probe "toter Alias, flach"          rot "$d/tief"
  probe "toter Alias im Unterbaum"    rot "$d/tief/sub"
  probe "toter Alias als Punkt-Name"  rot "$d/punkt"
  probe "Alias auf Verzeichnis"       ok  "$d/dir"
  probe "Symlink-Zyklus"              rot "$d/zyklus"
  probe "echte Datei, kein Alias"     ok  "$d/echt"
  probe "leeres Verzeichnis"          ok  "$d/leer"
  probe "Verzeichnis fehlt"           ok  "$d/gibt-es-nicht"

  if [ "$fails" -eq 0 ]; then
    echo "fetch-baseline-cache: selftest ok (9 Proben)"
    return 0
  fi
  echo "fetch-baseline-cache: selftest FEHLGESCHLAGEN (${fails} von 9)" >&2
  return 1
}

verify() {
  for c in sha256sum find readlink; do
    command -v "$c" >/dev/null 2>&1 \
      || { echo "fetch-baseline-cache: '$c' nicht gefunden (Host-Werkzeug)" >&2; exit 1; }
  done
  [ -f "$sums" ] \
    || { echo "fetch-baseline-cache: ${sums} fehlt — erst re-vendor (ohne --verify) laufen" >&2; exit 1; }
  echo "fetch-baseline-cache: verify ${baseline} gegen SHA256SUMS"
  ( cd "$baseline" && sha256sum -c SHA256SUMS )
  # Manifest-Deckung: Datei-Anzahl auf Platte == Manifest-Zeilen. Das faengt
  # Post-Vendor-Drift (gelöschte/zusätzliche Datei, unvollständiger Checkout) ab —
  # `sha256sum -c` allein passierte eine untermengige, in sich konsistente
  # SHA256SUMS grün. Die Under-Copy-Barriere (Quelle vs. vendored) sitzt im
  # re-vendor-Pfad.
  # Gezählt wird das GANZE Tag-Verzeichnis ohne das Manifest selbst, nicht nur
  # die zwei Bäume: eine Datei, die als Geschwister von regelwerk/ und templates/
  # abgelegt wird, liegt in keinem der beiden und stünde auch in keiner
  # Manifest-Zeile — beide Zählungen blieben gleich, und der Gate meldete grün.
  local on_disk manifest
  on_disk="$(find "${baseline}" -type f ! -path "$sums" 2>/dev/null | wc -l)"
  manifest="$(grep -c . "$sums" || true)"
  [ "$on_disk" -gt 0 ] \
    || { echo "fetch-baseline-cache: 0 Dateien — leeres/kaputtes Vendoring" >&2; exit 1; }
  [ "$on_disk" = "$manifest" ] \
    || { echo "fetch-baseline-cache: Manifest (${manifest} Zeilen) != Dateien auf Platte (${on_disk}) — unvollständig" >&2; exit 1; }
  check_aliases || exit 1
  echo "fetch-baseline-cache: verify ok (${manifest} Dateien, vollständig)"
}

check_latest() {
  # Upstream-Audit: (A) Currency (Release-Liste) und (B) Content-Drift am
  # gepinnten Tag (beide Bäume). Netz, informativ, KEIN Gate, KEIN fail-closed
  # (Ausfall → SKIP je Teil). Läuft unter deaktiviertem errexit (Dispatch).
  command -v curl >/dev/null 2>&1 \
    || { echo "fetch-baseline-cache: 'curl' nicht gefunden (Host-Werkzeug)" >&2; exit 1; }

  # Zeitgrenzen sind hier KEIN Komfort, sondern die Bedingung des fail-open:
  # ein Fehlschlag wird zu SKIP, eine haengende Verbindung wuerde ohne sie erst
  # von der Job-Decke abgeraeumt — und faerbte den Nachtlauf rot, statt still zu
  # ueberspringen. Grosszuegig gewaehlt: ein langsames Netz soll melden duerfen.
  local CT=10 MT=60

  # --- (A) Currency: Release-LISTE gegen den Pin (nicht /releases/latest, der
  #         Prereleases überspringt und einen zurückgezogenen Pin verbirgt) ---
  local api tags newest newer currency="skip"
  api="https://api.github.com/repos/${repo}/releases?per_page=100"
  tags="$(curl -fsSL --connect-timeout "$CT" --max-time "$MT" -H 'Accept: application/vnd.github+json' "$api" 2>/dev/null \
    | grep -o '"tag_name"[[:space:]]*:[[:space:]]*"[^"]*"' \
    | sed 's/.*"\([^"]*\)"$/\1/' \
    | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | LC_ALL=C sort -V -u || true)"
  if [ -n "$tags" ]; then
    newest="$(printf '%s\n' "$tags" | tail -1)"
    newer="$(printf '%s\n' "$tags" | LC_ALL=C sort -V | sed -n "/^${tag}\$/,\$p" | tail -n +2)"
    if [ "$newest" = "$tag" ]; then currency="current"
    elif [ -n "$newer" ];   then currency="newer"
    else                         currency="ahead"; fi
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
    if curl -fsSL --connect-timeout "$CT" --max-time "$MT" -o "${td}/lab-regelwerk.zip" "$url" 2>/dev/null \
       && extract_both_trees "${td}/lab-regelwerk.zip" "${td}/vendored" 2>/dev/null; then
      local up ve
      # xargs -r: eine leere Datei-Liste darf sha256sum NICHT ohne Argumente
      # starten — das hashte stdin und meldete einen falschen Drift.
      up="$( { cd "${td}/vendored" && find regelwerk templates -type f | LC_ALL=C sort | xargs -r sha256sum; } 2>/dev/null | LC_ALL=C sort )"
      ve="$( LC_ALL=C sort < "$sums" )"
      if [ -n "$up" ]; then
        [ "$up" = "$ve" ] && authenticity="ok" || authenticity="drift"
      else
        note="Upstream-Extraktion lieferte keine Dateien (Bundle-Layout?)"
      fi
    else
      note="Asset ${tag}/lab-regelwerk.zip nicht ladbar/entpackbar (Netz/Layout)"
    fi
    rm -rf "$td"
  fi

  local rc=0
  case "$currency" in
    current) echo "fetch-baseline-cache: check-latest OK (Currency) — Pin ${tag} ist der neueste Release-Tag." ;;
    newer)   echo "fetch-baseline-cache: check-latest — NEUER RELEASE verfügbar (Pin ${tag}):" >&2
             printf '%s\n' "$newer" | sed 's/^/  /' >&2
             echo "  -> Re-Adopt erwägen: fetch-baseline-cache.sh <neuer-tag> (re-vendor), §Baseline-Pin + Pointer nach MR-Bump-Prozedur." >&2
             rc=3 ;;
    ahead)   echo "fetch-baseline-cache: check-latest (Currency) — Pin ${tag} nicht in der Release-Liste (zurückgezogen/Fenster?); manuell prüfen." >&2 ;;
    skip)    echo "fetch-baseline-cache: check-latest SKIP (Currency) — Release-Liste nicht lesbar (Netz/API/Rate-Limit). Pin: ${tag}." ;;
  esac
  case "$authenticity" in
    ok)    echo "fetch-baseline-cache: check-latest OK (Content) — gepinnter Tag ${tag} upstream unverändert (Bytes == vendored SHA256SUMS)." ;;
    drift) echo "fetch-baseline-cache: check-latest — UPSTREAM-CONTENT-DRIFT: die Bytes von ${tag}/lab-regelwerk.zip weichen vom vendored SHA256SUMS ab (Tag verschoben / Asset neu?)." >&2
           echo "  -> Provenienz prüfen: das gepinnte Release-Asset hat sich geändert. NICHT blind re-vendoren." >&2
           rc=4 ;;
    skip)  echo "fetch-baseline-cache: check-latest SKIP (Content) — ${note}. Pin: ${tag}." ;;
  esac
  return $rc
}

if [ "$mode" = "selftest" ]; then selftest; exit $?; fi
if [ "$mode" = "verify" ]; then verify; exit 0; fi
if [ "$mode" = "check-latest" ]; then set +e; check_latest; rc=$?; set -e; exit "$rc"; fi

# --- re-vendor (Netz) ---
for cmd in curl unzip sha256sum find; do
  command -v "$cmd" >/dev/null 2>&1 \
    || { echo "fetch-baseline-cache: '${cmd}' nicht gefunden (Host-Werkzeug)" >&2; exit 1; }
done

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

url="https://github.com/${repo}/releases/download/${tag}/lab-regelwerk.zip"
echo "fetch-baseline-cache: ${tag}/lab-regelwerk.zip -> ${baseline}/{regelwerk,templates}/"
curl -fsSL -o "${tmp}/lab-regelwerk.zip" "$url"

src_n=0
extract_both_trees "${tmp}/lab-regelwerk.zip" "$baseline"

# Under-Copy-Barriere (die eigentliche Vollständigkeitsprüfung): der vendorte
# Bestand muss exakt so viele Dateien tragen wie die Quelle im ZIP. Ein reiner
# Post-Copy-vs-Manifest-Vergleich (verify) sähe einen übergangenen Quell-Zweig
# nicht, weil Manifest und Platte beide post-copy sind.
dst_n="$(find "${baseline}/regelwerk" "${baseline}/templates" -type f | wc -l)"
[ "$src_n" = "$dst_n" ] \
  || { echo "fetch-baseline-cache: Quelle (${src_n}) != vendored (${dst_n}) — Kopierschritt hat Dateien übergangen" >&2; exit 1; }

# Manifest über den TATSÄCHLICHEN Dateibestand BEIDER Bäume (find, nicht
# Top-Level-Glob).
( cd "$baseline" && find regelwerk templates -type f | LC_ALL=C sort | xargs sha256sum > SHA256SUMS )
verify

echo "fetch-baseline-cache: fertig — vendored ${baseline}/{regelwerk,templates} (+SHA256SUMS)"
