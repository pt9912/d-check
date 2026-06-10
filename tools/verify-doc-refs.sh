#!/usr/bin/env bash
# verify-doc-refs — Bootstrap-Sensor: lokale Markdown-Linkziele
# ([text](pfad)) in der Doku müssen existieren. Externe Links und
# reine Anker (#…) werden ignoriert.
#
# VENDORED aus d-migrate/scripts/verify-doc-refs.sh (Stand 2026-06-10),
# angepasst: zusätzliche Scan-Wurzel harness/ und Top-Level AGENTS.md,
# CLAUDE.md; Fenced-Code-Blöcke (``` / ~~~) werden übersprungen.
# Deklariert in harness/conventions.md MR-003 —
# Auflösungs-Trigger: slice-004 ersetzt dieses Skript durch d-check
# selbst (Dogfooding). Bis dahin: bei Fixes an der Quelle nachziehen.
#
# Usage:
#   tools/verify-doc-refs.sh [root-dir]
#
# Exit codes:
#   0  passed
#   1  broken local link target detected
#   2  environment error
set -euo pipefail

script_dir="$(cd "$(dirname "$0")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"
root="${1:-$repo_root}"

if [[ ! -d "$root" ]]; then
    echo "ERROR: root directory not found: $root" >&2
    exit 2
fi

extract_local_markdown_links() {
    awk '
        # Fenced-Code-Blöcke überspringen — Markdown-Beispiele mit
        # relativen Links darin sind keine echten Referenzen.
        /^[[:space:]]*(```|~~~)/ { in_fence = !in_fence; next }
        in_fence { next }
        {
            line = $0
            # Inline-Code (einfache Backticks) vor dem Link-Scan
            # entfernen — sonst werden Regex-/Code-Beispiele mit "]("
            # in Backticks als Markdown-Links fehlinterpretiert.
            gsub(/`[^`]*`/, "", line)
            while (match(line, /!?\[[^]]*\]\([^)]*\)/)) {
                link = substr(line, RSTART, RLENGTH)
                line = substr(line, RSTART + RLENGTH)

                if (substr(link, 1, 1) == "!") {
                    continue
                }

                sub(/^!?\[[^]]*\]\(/, "", link)
                sub(/\)$/, "", link)

                if (link ~ /^</) {
                    sub(/^</, "", link)
                    sub(/>.*/, "", link)
                } else {
                    sub(/[[:space:]].*/, "", link)
                }

                sub(/#.*/, "", link)

                if (link == "" ||
                    link ~ /^[a-zA-Z][a-zA-Z0-9+.-]*:/) {
                    continue
                }

                print link
            }
        }
    ' "$1" | sort -u
}

broken=0

while IFS= read -r md; do
    rel="${md#"$root"/}"
    while IFS= read -r target; do
        if [[ "$target" == /* ]]; then
            resolved="$target"
        else
            resolved="$(dirname "$md")/$target"
        fi
        if [[ ! -e "$resolved" ]]; then
            echo "BROKEN: $rel -> $target"
            ((++broken))
        fi
    done < <(extract_local_markdown_links "$md")
done < <(
    {
        for docs_dir in "$root/docs" "$root/spec" "$root/harness"; do
            if [[ -d "$docs_dir" ]]; then
                find "$docs_dir" -name '*.md' -type f
            fi
        done
        for top_level_doc in "$root/README.md" "$root/CHANGELOG.md" \
                             "$root/AGENTS.md" "$root/CLAUDE.md"; do
            if [[ -f "$top_level_doc" ]]; then
                printf '%s\n' "$top_level_doc"
            fi
        done
    } | sort
)

if [[ "$broken" -gt 0 ]]; then
    echo "$broken broken documentation link(s)"
    exit 1
fi
echo "All documentation links OK."
