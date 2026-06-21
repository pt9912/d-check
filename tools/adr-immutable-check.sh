#!/usr/bin/env bash
# adr-immutable-check.sh — ADR-Immutable-Gate (ADR-0016, slice-041).
#
# Erzwingt AGENTS.md §3.5: eine Accepted-ADR wird nicht inhaltlich
# ueberschrieben. Erlaubt bleiben nur Anhaenge unter `## Geschichte` und der
# `**Status:**`-Uebergang. EINE Wahrheit fuer drei Aufrufer:
#   - lokaler pre-commit-Hook:  adr-immutable-check.sh --staged
#   - PR-/Push-CI:              adr-immutable-check.sh --range <BASE>..<HEAD>
#   - lokal (make adr-check):   adr-immutable-check.sh            (Selbsttest + HEAD~1..HEAD)
#                               adr-immutable-check.sh --self-test (nur Selbsttest)
#
# Methode (robust statt roher Zeilen-Diff): "immutable core" einer ADR-Datei
# = Inhalt OHNE den ## Geschichte-Abschnitt (bis zur naechsten ## -H2 bzw.
# EOF) UND OHNE die **Status:**-Zeile. Geprueft wird je MODIFIZIERTE
# docs/plan/adr/NNNN-*.md, deren BASE-Version `Status: Accepted` traegt:
# core(BASE) != core(HEAD) -> FAIL. Zusaetzlich muss die HEAD-Status-Zeile
# zulaessig bleiben (Accepted… | Superseded by ADR-NNNN). Geloeschte/
# umbenannte Accepted-ADRs -> FAIL.
#
# fail-closed: unbekannter Modus / fehlende Argumente / kaputter Range -> 2.
set -euo pipefail

cd "$(git rev-parse --show-toplevel 2>/dev/null || echo .)"

STATUS_ACCEPTED_RE='^\*\*Status:\*\* Accepted'
STATUS_HEAD_OK_RE='^\*\*Status:\*\* (Accepted|Superseded by ADR-[0-9]{4})'

# Nur echte ADR-Dateien: docs/plan/adr/NNNN-*.md (Index README.md ausgenommen).
is_adr_file() { [[ "$1" =~ ^docs/plan/adr/[0-9]{4}-.*\.md$ ]]; }

# immutable core: Geschichte-Abschnitt + die Metadaten-Status-Zeile entfernen.
# Liest stdin. Die **Status:**-Zeile wird NUR im Kopf-Block (vor der ersten
# `## `-H2) gestrippt — eine gleichlautende Zeile im Koerper bleibt Teil des
# core (sonst False-Negative: ein Edit daran wuerde durchrutschen; R1-MEDIUM).
core() {
  awk '
    /^## Geschichte[[:space:]]*$/ { skip=1; next }
    skip && /^## /                { skip=0; body=1; print; next }
    skip                          { next }
    /^## /                        { body=1; print; next }
    !body && /^\*\*Status:\*\*/   { next }
    { print }
  '
}

# Vergleicht zwei Inhalts-Dateien (BASE, HEAD) einer ADR. 0 = ok, 1 = Verstoss.
classify() { # $1 base-file, $2 head-file, $3 pfad-label
  local base="$1" head="$2" path="$3" fail=0
  # Nur die Metadaten-Status-Zeile (erste **Status:**-Zeile) zaehlt — eine
  # gleichlautende Zeile im Koerper darf den Status-Befund nicht verfaelschen.
  local base_status head_status
  base_status="$(grep -m1 -E '^\*\*Status:\*\*' "$base" || true)"
  head_status="$(grep -m1 -E '^\*\*Status:\*\*' "$head" || true)"
  # Nur Accepted-ADRs sind immutabel; Proposed (BASE) ist frei.
  grep -qE "$STATUS_ACCEPTED_RE" <<<"$base_status" || return 0
  if ! diff -q <(core <"$base") <(core <"$head") >/dev/null 2>&1; then
    echo "adr-immutable: FAIL — $path aendert den Koerper einer Accepted-ADR (nur ## Geschichte-Anhang + Status-Uebergang erlaubt; AGENTS §3.5)" >&2
    fail=1
  fi
  if ! grep -qE "$STATUS_HEAD_OK_RE" <<<"$head_status"; then
    echo "adr-immutable: FAIL — $path setzt den Status einer Accepted-ADR auf einen unzulaessigen Wert (erlaubt: 'Accepted…' oder 'Superseded by ADR-NNNN')" >&2
    fail=1
  fi
  return "$fail"
}

# Ein Diff-Eintrag (status, p1, p2) gegen zwei git-Refs pruefen.
# $4 = wie BASE-Inhalt zu holen ist (ref oder ":"=index), $5 = HEAD-Inhalt.
check_entry() { # $1 st $2 p1 $3 p2 $4 base-ref $5 head-ref
  local st="$1" p1="$2" p2="$3" base_ref="$4" head_ref="$5"
  local bt ht rc=0
  case "$st" in
    A*) return 0 ;;  # neue ADR — noch nicht immutabel
    D*)
      is_adr_file "$p1" || return 0
      if show_ref "$base_ref" "$p1" | grep -qE "$STATUS_ACCEPTED_RE"; then
        echo "adr-immutable: FAIL — $p1 (Accepted) geloescht; AGENTS §3.5" >&2
        return 1
      fi ;;
    R*)
      # Rename: p1=alt, p2=neu. Eine Accepted-ADR behaelt ihren NNNN-Pfad.
      is_adr_file "$p1" || return 0
      if show_ref "$base_ref" "$p1" | grep -qE "$STATUS_ACCEPTED_RE"; then
        echo "adr-immutable: FAIL — $p1 (Accepted) umbenannt nach ${p2:-?}; ADR-Pfad ist stabil; AGENTS §3.5" >&2
        return 1
      fi ;;
    M*|T*)
      is_adr_file "$p1" || return 0
      bt="$(mktemp)"; ht="$(mktemp)"
      show_ref "$base_ref" "$p1" > "$bt" || true
      show_ref "$head_ref" "$p1" > "$ht" || true
      classify "$bt" "$ht" "$p1" || rc=1
      rm -f "$bt" "$ht"
      return "$rc" ;;
    *)
      is_adr_file "$p1" || return 0
      echo "adr-immutable: FAIL — $p1: unerwarteter Diff-Status '$st' (fail-closed)" >&2
      return 1 ;;
  esac
  return 0
}

# Holt Datei-Inhalt aus einem Ref. ref=":" -> staged (index): git show :path.
show_ref() { # $1 ref, $2 path
  if [ "$1" = ":" ]; then
    git show ":$2" 2>/dev/null || true
  else
    git show "$1:$2" 2>/dev/null || true
  fi
}

# Treiber: jede geaenderte ADR im Diff base..head bzw. staged pruefen.
process_diff() { # $1 base-ref $2 head-ref $3 diff-cmd-tag (range|staged)
  local base_ref="$1" head_ref="$2" tag="$3" fail=0 n=0
  while IFS=$'\t' read -r st p1 p2; do
    [ -z "$st" ] && continue
    n=$((n + 1))
    check_entry "$st" "$p1" "${p2:-}" "$base_ref" "$head_ref" || fail=1
  done < <(adr_diff "$tag" "$base_ref" "$head_ref")
  # Hook-Hygiene: im staged-Modus bei Erfolg still (wie trace-check --message).
  if [ "$fail" -eq 0 ] && [ "$tag" != staged ]; then
    echo "adr-immutable: $n ADR-Diff-Eintrag/e geprueft, keine Accepted-Verletzung (Selbsttest gefeuert)."
  fi
  return "$fail"
}

adr_diff() { # $1 tag, $2 base-ref, $3 head-ref
  case "$1" in
    range)  git diff --name-status -M "$2" "$3" -- docs/plan/adr/ ;;
    staged) git diff --cached --name-status -M -- docs/plan/adr/ ;;
  esac
}

# Negativ-Selbsttest auf der classify-Kernlogik (ohne git): beweist bei jedem
# Lauf, dass Koerper-Edit/Status-Rueckfall feuern und legitime Aenderungen
# (Geschichte-Anhang, Superseded-Uebergang, Proposed-BASE) NICHT feuern.
self_test() {
  local tmp; tmp="$(mktemp -d)"
  local body='# ADR-0099 — X

**Status:** Accepted
**Datum:** 2026-01-01

## Entscheidung

Tue A.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-01-01 | Proposed → Accepted |'

  printf '%s\n' "$body" > "$tmp/base"

  # (1) Geschichte-Anhang -> darf NICHT feuern.
  { printf '%s\n' "$body"; printf '| 2026-02-02 | Notiz |\n'; } > "$tmp/hist"
  classify "$tmp/base" "$tmp/hist" "selftest-hist" >/dev/null 2>&1 \
    || { echo "adr-immutable: Selbsttest FEHLGESCHLAGEN — Geschichte-Anhang faelschlich gefeuert" >&2; rm -rf "$tmp"; exit 2; }

  # (2) Status -> Superseded (+ Geschichte) -> darf NICHT feuern.
  printf '%s\n' "${body/Accepted/Superseded by ADR-0100}" > "$tmp/sup"
  classify "$tmp/base" "$tmp/sup" "selftest-sup" >/dev/null 2>&1 \
    || { echo "adr-immutable: Selbsttest FEHLGESCHLAGEN — Superseded-Uebergang faelschlich gefeuert" >&2; rm -rf "$tmp"; exit 2; }

  # (3) Koerper-Edit -> MUSS feuern.
  printf '%s\n' "${body/Tue A./Tue B.}" > "$tmp/body"
  if classify "$tmp/base" "$tmp/body" "selftest-body" >/dev/null 2>&1; then
    echo "adr-immutable: Selbsttest FEHLGESCHLAGEN — Koerper-Edit nicht erkannt" >&2; rm -rf "$tmp"; exit 2
  fi

  # (4) Status-Rueckfall Accepted -> Proposed -> MUSS feuern.
  printf '%s\n' "${body/Accepted/Proposed}" > "$tmp/regress"
  if classify "$tmp/base" "$tmp/regress" "selftest-regress" >/dev/null 2>&1; then
    echo "adr-immutable: Selbsttest FEHLGESCHLAGEN — Status-Rueckfall nicht erkannt" >&2; rm -rf "$tmp"; exit 2
  fi

  # (5) BASE Proposed -> frei (Koerper-Edit darf NICHT feuern).
  printf '%s\n' "${body/Accepted/Proposed}" > "$tmp/pbase"
  printf '%s\n' "${body/Accepted/Proposed}" | sed 's/Tue A./Tue B./' > "$tmp/pedit"
  classify "$tmp/pbase" "$tmp/pedit" "selftest-proposed" >/dev/null 2>&1 \
    || { echo "adr-immutable: Selbsttest FEHLGESCHLAGEN — Proposed-BASE faelschlich gefeuert" >&2; rm -rf "$tmp"; exit 2; }

  # (6) Koerper-Abschnitt NACH ## Geschichte -> Edit dort MUSS feuern
  # (core-resume an der naechsten ## -H2; slice-040-Klasse Boundary, R1-LOW).
  local g6='# ADR-0099 — X

**Status:** Accepted

## Entscheidung

Tue A.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-01-01 | Proposed → Accepted |

## Anhang

Detail A.'
  printf '%s\n' "$g6" > "$tmp/g6base"
  printf '%s\n' "${g6/Detail A./Detail B.}" > "$tmp/g6head"
  if classify "$tmp/g6base" "$tmp/g6head" "selftest-resume" >/dev/null 2>&1; then
    echo "adr-immutable: Selbsttest FEHLGESCHLAGEN — Edit im Abschnitt nach ## Geschichte nicht erkannt" >&2; rm -rf "$tmp"; exit 2
  fi

  # (7) Koerper-Zeile, die mit **Status:** beginnt -> Edit MUSS feuern
  # (nur die Metadaten-Status-Zeile darf gestrippt werden; R1-MEDIUM-Regress).
  local s7='# ADR-0099 — X

**Status:** Accepted

## Konsequenzen

**Status:** der Migration bleibt offen.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-01-01 | Proposed → Accepted |'
  printf '%s\n' "$s7" > "$tmp/s7base"
  printf '%s\n' "${s7/bleibt offen./ist geklaert.}" > "$tmp/s7head"
  if classify "$tmp/s7base" "$tmp/s7head" "selftest-bodystatus" >/dev/null 2>&1; then
    echo "adr-immutable: Selbsttest FEHLGESCHLAGEN — Edit an Koerper-**Status:**-Zeile nicht erkannt (core strippt zu breit)" >&2; rm -rf "$tmp"; exit 2
  fi

  rm -rf "$tmp"
}

mode="${1:-}"
case "$mode" in
  --range)
    range="${2:-}"
    # `..`-Separator verlangen (sonst base==head -> leerer Diff, still gruen; R1-LOW).
    [[ -n "$range" && "$range" == *..* ]] || { echo "adr-immutable: --range braucht <base>..<head>" >&2; exit 2; }
    self_test
    base="${range%%..*}"; head="${range##*..}"
    [ -n "$base" ] && [ -n "$head" ] || { echo "adr-immutable: --range braucht <base>..<head>" >&2; exit 2; }
    # fail-closed: nicht aufloesbare Basis blockiert laut (analog trace-check).
    if [[ "$base" =~ ^0*$ ]] || ! git rev-parse -q --verify "${base}^{commit}" >/dev/null 2>&1; then
      echo "adr-immutable: FAIL — Range-Basis '$base' nicht aufloesbar; der CI-Workflow muss eine gueltige Basis liefern." >&2
      exit 2
    fi
    process_diff "$base" "$head" range
    ;;
  --staged)
    self_test
    # BASE = HEAD-Version, HEAD = staged (index). Beim allerersten Commit gibt
    # es kein HEAD -> nichts zu schuetzen.
    if git rev-parse -q --verify HEAD >/dev/null 2>&1; then
      process_diff HEAD ":" staged
    fi
    ;;
  --self-test)
    self_test
    echo "adr-immutable: Selbsttest gruen."
    ;;
  "")
    self_test
    if git rev-parse -q --verify HEAD~1 >/dev/null 2>&1; then
      process_diff HEAD~1 HEAD range
    else
      echo "adr-immutable: kein HEAD~1 — nur Selbsttest (gruen)."
    fi
    ;;
  *)
    echo "adr-immutable: unbekannter Modus '$mode' (erwartet --range|--staged|--self-test|<leer>)" >&2
    exit 2
    ;;
esac
