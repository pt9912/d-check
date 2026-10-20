#!/usr/bin/env bash
# ci-range-probe.sh — fährt ci-range.sh gegen ein Wegwerf-Repo, je Fall eine Probe:
# Pull-Request, neuer Branch, normaler Push, Force-Push auf einen Feature-Branch
# (Basis fehlt), Force-Push auf den Default-Branch (laut). EIN KOMMANDO, KEIN GATE
# (ADR-0109). Kein Netz, kein Docker: nur git und bash.
#
# GRENZE. Geprüft ist die Range, die das Skript liefert, nicht ihre Verwendung in
# der CI; ob die Workflow-Datei das Skript ruft, liest `make doc-check` nicht.
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
script="$root/tools/harness/ci-range.sh"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

git init -q --bare "$tmp/origin.git"
git init -q -b main "$tmp/w"
cd "$tmp/w"
git config user.email probe@example.invalid
git config user.name probe
git remote add origin "$tmp/origin.git"
commit() { git commit -q --allow-empty -m "$1"; git rev-parse HEAD; }
c1=$(commit "c1"); git push -q origin main
git checkout -q -b feature
f1=$(commit "f1")
f2=$(commit "f2")
git push -q origin feature

fail=0
expect() { # name expected-stdout [expected-exit] -- env...
  local name="$1" want="$2" wantexit="$3"; shift 3
  local out rc=0
  out=$(env -i PATH="$PATH" HOME="$tmp" "$@" bash "$script" 2>/dev/null) || rc=$?
  if [[ "$rc" != "$wantexit" || "$out" != "$want" ]]; then
    echo "ci-range-probe: FAIL $name -- Exit $rc (erwartet $wantexit), Ausgabe '$out' (erwartet '$want')" >&2
    fail=1
  else
    echo "ci-range-probe: ok   $name"
  fi
}

zero=0000000000000000000000000000000000000000
gone=1111111111111111111111111111111111111111
expect "pull-request" "$c1..$f2" 0 EVENT=pull_request PR_BASE="$c1" PUSH_BEFORE="" HEAD_SHA="$f2" DEFAULT_BRANCH=main REF_NAME=feature
expect "neuer-branch" "origin/main..$f2" 0 EVENT=push PUSH_BEFORE="$zero" HEAD_SHA="$f2" DEFAULT_BRANCH=main REF_NAME=feature
expect "normaler-push" "$f1..$f2" 0 EVENT=push PUSH_BEFORE="$f1" HEAD_SHA="$f2" DEFAULT_BRANCH=main REF_NAME=feature
expect "force-push-feature" "origin/main..$f2" 0 EVENT=push PUSH_BEFORE="$gone" HEAD_SHA="$f2" DEFAULT_BRANCH=main REF_NAME=feature
expect "force-push-default" "" 2 EVENT=push PUSH_BEFORE="$gone" HEAD_SHA="$f2" DEFAULT_BRANCH=main REF_NAME=main

# Ohne origin/<default> kein Rückfall: laut statt leerer Range.
git remote remove origin
git update-ref -d refs/remotes/origin/main
expect "force-push-ohne-origin" "" 2 EVENT=push PUSH_BEFORE="$gone" HEAD_SHA="$f2" DEFAULT_BRANCH=main REF_NAME=feature

exit "$fail"
