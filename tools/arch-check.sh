#!/usr/bin/env bash
# arch-check — Fitness Function zu ADR-0005 (Hexagon light,
# u-boot-Ordnerkonvention). Prüft die Import-Regeln über `go list`;
# Verstöße → Exit 1. Läuft als Dockerfile-Stage (make arch-check).
#
# Regeln:
#   R1  hexagon/* importiert weder I/O-APIs (os, syscall, io/fs,
#       net-Sockets, net/http) noch Adapter noch die YAML-Bibliothek.
#       Reine Parser ohne I/O (net/url) sind erlaubt.
#   R2  net/http ausschließlich in internal/adapter/driven/httpcheck.
#   R3  gopkg.in/yaml.v3 ausschließlich in
#       internal/adapter/driven/{configyaml,report}
#       (ADR-0009: Encode der YAML-Ausgabe im report-Adapter).
#   R4  os ausschließlich in internal/adapter/driven/fs,
#       internal/adapter/driving/cli und cmd/* (Composition Root).
#   R5  driven Adapter importieren einander nicht.
set -euo pipefail

MODULE="$(go list -m)"
fail=0

violation() {
  echo "ARCH-CHECK FAIL (ADR-0005, $1): $2" >&2
  fail=1
}

while IFS='|' read -r pkg imports; do
  rel="${pkg#"$MODULE"}"
  rel="${rel#/}"
  for imp in $imports; do
    case "$rel" in
      internal/hexagon/*)
        case "$imp" in
          net/url) ;; # reiner Parser, kein I/O — erlaubt
          os|os/*|net|net/*|syscall|io/fs)
            violation R1 "$rel importiert $imp" ;;
          "$MODULE"/internal/adapter/*)
            violation R1 "$rel importiert Adapter $imp" ;;
          gopkg.in/yaml.v3)
            violation R1 "$rel importiert yaml" ;;
        esac
        ;;
    esac
    if [ "$imp" = "net/http" ] && [ "$rel" != "internal/adapter/driven/httpcheck" ]; then
      violation R2 "$rel importiert net/http"
    fi
    if [ "$imp" = "gopkg.in/yaml.v3" ]; then
      case "$rel" in
        internal/adapter/driven/configyaml|internal/adapter/driven/report) ;;
        *) violation R3 "$rel importiert gopkg.in/yaml.v3" ;;
      esac
    fi
    if [ "$imp" = "os" ]; then
      case "$rel" in
        internal/adapter/driven/fs|internal/adapter/driving/cli|cmd/*) ;;
        *) violation R4 "$rel importiert os" ;;
      esac
    fi
    case "$rel" in
      internal/adapter/driven/*)
        case "$imp" in
          "$MODULE"/internal/adapter/driven/*)
            from="${rel#internal/adapter/driven/}"
            to="${imp#"$MODULE"/internal/adapter/driven/}"
            if [ "${from%%/*}" != "${to%%/*}" ]; then
              violation R5 "$rel importiert $imp"
            fi
            ;;
        esac
        ;;
    esac
  done
done < <(go list -f '{{.ImportPath}}|{{join .Imports " "}}' ./...)

if [ "$fail" -ne 0 ]; then
  exit 1
fi
echo "arch-check ok: Import-Regeln R1–R5 (ADR-0005) eingehalten."
