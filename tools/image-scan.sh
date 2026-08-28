#!/usr/bin/env bash
# Trivy gegen die PUBLIZIERTEN Container-Images (ADR-0066).
#
# ZUSAGE: meldet bekannte Schwachstellen in dem Bild, das Anwender ziehen —
# nicht im Arbeitsbaum. Zwischen zwei Releases altert das publizierte Image,
# ohne dass sich das Repo aendert; ein push-getriggertes Gate ist gegen diesen
# Fall prinzipiell blind.
#
# ABGRENZUNG: NICHT in `gates`. Der Scan braucht Netz fuer die Vuln-DB, und das
# ist hier der ZWECK, nicht ein Zugestaendnis — eine gepinnte DB faende nur die
# CVEs von gestern. Damit steht er in derselben Klasse wie die Frische-Achsen,
# nicht bei `semgrep` (ADR-0010, netzlos und hermetisch).
#
# KOPPLUNG: BEIDE Trivy-Laeufe fahren `--exit-code 0`, und das ist der Kern der
# Fehlerbehandlung. GEMESSEN: ein nicht existierendes Image quittiert Trivy mit
# `--exit-code 1` ebenfalls mit 1 -- Fehler und Befund waeren nicht zu
# unterscheiden, und das Gate meldete "behebbare CRITICAL/HIGH", wo gar nicht
# geprueft wurde. Mit `--exit-code 0` heisst ein Nicht-Null-Exit von Trivy
# eindeutig "Scan gescheitert"; ueber Befunde entscheidet die AUSWERTUNG.
#
# GRENZE: kein Docker-Socket. Trivy liest die Images aus der Registry; fuer
# publizierte Bilder braucht es ihn nicht, und ein gemounteter Socket waere ein
# Host-Root-Pfad fuer ein Werkzeug, das ihn nicht noetig hat.
#
# GRENZE: das Runtime-Image ist `distroless/static` plus statisches Go-Binary --
# es gibt praktisch keine OS-Paket-Flaeche. Der Fund-Raum ist die Go-Build-Info
# des Binaries. Ein gruener Lauf sagt "nichts Bekanntes in diesem Raum", nicht
# "das Image ist sicher".
#
# Exit-Codes: 0 = keine behebbaren CRITICAL/HIGH, 1 = solche gefunden,
#             2 = Scan gescheitert (Registry, DB, Image).
set -uo pipefail

# Digest-Pin (ADR-0011): Tag bleibt lesbar, der @sha256:-Digest ist die
# Wahrheit. Ein Scanner, der sich unter der Hand aendert, macht
# Befund-Vergleiche ueber die Zeit wertlos.
TRIVY_VERSION="${TRIVY_VERSION:-0.74.0}"
TRIVY_DIGEST="${TRIVY_DIGEST:-sha256:62b1e65e8869bc4b4c6aa4fa2b21595256c7c2f6018a9d9ad61caf87187c1969}"

# Der Docker-Hub-Spiegel gehoert hier hinein, SOBALD er ein Bild traegt
# (DC-FA-DIST-002). Ein Ref ohne Bild machte den Nachtlauf ab dem ersten Tag rot.
IMAGE_SCAN_REFS="${IMAGE_SCAN_REFS:-ghcr.io/pt9912/d-check:latest}"

# Cache ausserhalb des Repos, wie das Regelset von `semgrep` (ADR-0010): der
# Arbeitsbaum bleibt sauber, und `git status` meldet keine Werkzeug-Artefakte.
CACHE="${TRIVY_CACHE:-$HOME/.cache/d-check/trivy}"
mkdir -p "$CACHE"

# Eine Zeile je Befund plus ein zaehlbarer Marker -- Trivys eigenes
# Template-Format statt eines JSON-Parsers. Ein Fremd-Interpreter waere eine
# vierte Toolchain (AGENTS.md §3.1, MR-040).
TPL='{{ range . }}{{ range .Vulnerabilities }}FINDING {{ .Severity }} {{ .PkgName }} {{ .InstalledVersion }} -> fix {{ .FixedVersion }} {{ .VulnerabilityID }}
{{ end }}{{ end }}'

trivy() {
  docker run --rm \
    -v "${CACHE}:/root/.cache/trivy" \
    "aquasec/trivy:${TRIVY_VERSION}@${TRIVY_DIGEST}" \
    image --no-progress --scanners vuln --exit-code 0 "$@"
}

findings=0
errored=0

for ref in ${IMAGE_SCAN_REFS}; do
  echo "=============================================================="
  echo "== Vollbericht (alle Schweregrade): ${ref}"
  echo "=============================================================="
  # Faellt nie an Befunden -- beantwortet "was steckt gerade drin", auch wenn
  # nichts davon behebbar ist.
  if ! trivy --severity CRITICAL,HIGH,MEDIUM,LOW,UNKNOWN --format table "${ref}"; then
    echo "image-scan: Scan von ${ref} ist GESCHEITERT (nicht: Befunde gefunden)."
    errored=1
    continue
  fi

  echo
  echo "--------------------------------------------------------------"
  echo "-- Handlungspflichtig (CRITICAL/HIGH mit verfuegbarem Fix): ${ref}"
  echo "--------------------------------------------------------------"
  # Nur DIESER Lauf entscheidet ueber rot. Ein Nachtlauf, der an nicht
  # behebbaren Basis-Image-CVEs rot wird, ist in zwei Wochen ein weggeklicktes
  # Abzeichen und dann schlechter als nichts.
  if ! out="$(trivy --severity CRITICAL,HIGH --ignore-unfixed \
               --format template --template "${TPL}" "${ref}")"; then
    echo "image-scan: Entscheidungslauf fuer ${ref} ist GESCHEITERT."
    errored=1
    continue
  fi

  # `grep -c` liefert 1, wenn nichts passt -- deshalb `|| true`, sonst risse
  # der Zaehl-Pfad den Lauf ab und ein SAUBERES Image saehe aus wie ein Fehler.
  count="$(printf '%s' "${out}" | grep -c '^FINDING ' || true)"
  if [ "${count}" = "0" ]; then
    echo "OK — keine behebbaren CRITICAL/HIGH in ${ref}."
  else
    printf '%s\n' "${out}" | grep '^FINDING ' | sed 's/^FINDING /  /'
    echo "image-scan: ${ref}: ${count} behebbare CRITICAL/HIGH-Befunde."
    findings=1
  fi
  echo
done

if [ "${errored}" = "1" ]; then
  echo "image-scan: mindestens ein Scan ist GESCHEITERT — der Befundstand ist UNBEKANNT, nicht gruen."
  exit 2
fi
if [ "${findings}" = "1" ]; then
  exit 1
fi
echo "image-scan: keine behebbaren CRITICAL/HIGH in: ${IMAGE_SCAN_REFS}"
exit 0
