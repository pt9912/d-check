#!/usr/bin/env bash
# ci-range.sh — löst die Commit-Range des CI-Schritts „Traceability + ADR-Immutable
# + Slice-Closure-Übergang" auf und schreibt sie auf stdout (ADR-0013, ADR-0109).
#
# ZUSAGE. Ein Push-Lauf, dessen `before`-Commit der Klon nicht mehr kennt (Force-
# Push, Rebase), prüft `origin/<default>..HEAD` statt an der unauflösbaren Basis
# zu scheitern: mehr Commits, nie weniger. Ein neuer Branch, ein normaler Push und
# ein Pull-Request behalten ihre Range.
#
# GRENZE. Auf dem Default-Branch selbst bleibt eine nicht auflösbare Basis ein
# Fehler (Exit 2) — dort gäbe es keinen älteren Stand, gegen den man prüfen
# könnte, und `origin/<default>..HEAD` wäre leer (stilles Grün). Geprüft wird die
# AUFLÖSBARKEIT des `before`-Commits, nicht seine Eignung als Basis. Ein Repo
# ohne Remote `origin` gibt für den Rückfall keine Range her (Exit 2).
#
# EINGABE (Umgebung): EVENT (`pull_request` | `push`), PR_BASE, PUSH_BEFORE,
# HEAD_SHA, DEFAULT_BRANCH, REF_NAME (Branch des Laufs). Hinweise auf stderr.
set -euo pipefail

: "${EVENT:?EVENT fehlt}" "${HEAD_SHA:?HEAD_SHA fehlt}" "${DEFAULT_BRANCH:?DEFAULT_BRANCH fehlt}"

fetch_default() {
  git fetch --no-tags origin "+refs/heads/$DEFAULT_BRANCH:refs/remotes/origin/$DEFAULT_BRANCH" 2>/dev/null || true
}

if [[ "$EVENT" == "pull_request" ]]; then
  echo "$PR_BASE..$HEAD_SHA"
  exit 0
fi

if [[ "${PUSH_BEFORE:-}" =~ ^0+$ ]]; then
  # Neuer Branch (kein before-Commit): gegen den Default-Branch prüfen statt nur HEAD.
  fetch_default
  echo "origin/$DEFAULT_BRANCH..$HEAD_SHA"
  exit 0
fi

if git cat-file -e "${PUSH_BEFORE}^{commit}" 2>/dev/null; then
  echo "$PUSH_BEFORE..$HEAD_SHA"
  exit 0
fi

if [[ "${REF_NAME:-}" == "$DEFAULT_BRANCH" ]]; then
  echo "ci-range: Basis $PUSH_BEFORE auf dem Default-Branch '$DEFAULT_BRANCH' nicht aufloesbar -- kein Rueckfall (stilles Gruen ueber leerer Range)." >&2
  exit 2
fi

fetch_default
if ! git rev-parse --verify --quiet "origin/$DEFAULT_BRANCH^{commit}" >/dev/null; then
  echo "ci-range: Basis $PUSH_BEFORE nicht aufloesbar und origin/$DEFAULT_BRANCH fehlt -- keine Range." >&2
  exit 2
fi
echo "ci-range: Basis $PUSH_BEFORE nicht aufloesbar (Force-Push?) -- pruefe gegen origin/$DEFAULT_BRANCH." >&2
echo "origin/$DEFAULT_BRANCH..$HEAD_SHA"
