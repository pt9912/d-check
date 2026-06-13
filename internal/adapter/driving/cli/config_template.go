package cli

// configTemplate ist das statische `.d-check.yml`-Startgerüst, das
// `--print-config` unverändert auf stdout schreibt (DC-FA-CLI-005).
// Es dekodiert über den eigenen Konfigurations-Parser fehlerfrei und
// führt die aktiven Defaults plus die übrigen Module mit ihren Optionen
// als kommentierte Beispiele. Statisch — nicht aus Repo-Inhalt
// abgeleitet (Determinismus DC-QA-02; kein Repo-Zugriff DC-QA-03).
const configTemplate = `# .d-check.yml — Startgerüst, erzeugt von: d-check --print-config
# Umleiten zum Anlegen:  d-check --print-config > .d-check.yml
#
# Aktiv ohne Bearbeitung: die Module links + anchors über die ganze
# Repo-Wurzel. Die auskommentierten Blöcke zeigen die weiteren Module
# und ihre Optionen — einkommentieren und an das Repo anpassen.

scan:
  roots: ["."]            # zu prüfende Wurzeln; "." = ganze Repo-Wurzel
  # ignore: ["pfad/**"]   # Glob, relativ zur Wurzel (prunt den Abstieg)

modules: [links, anchors]
# Verfügbar: links, anchors, ids, matrix, codepaths, spans, hostpaths, external
# (external ist die einzige Netzwerk-Tür — strikt opt-in.)

# --- ids: Linkpflicht für Kennungen ---
# ids:
#   patterns:                      # Reihenfolge = Präzedenz
#     - regex: 'ADR-\d{4}'
#       target: docs/plan/adr/     # wo die Kennung definiert ist (Datei oder Verzeichnis)
#       link-policy: prose         # prose (Default) | always (auch Inline-Code linkpflichtig)
#       # exempt-paths: [CHANGELOG.md, "docs/reviews/**"]   # nur bei always
#   # scope:                       # optional: ersetzt den globalen Scan nur für ids
#   #   roots: [spec, docs]

# --- matrix: erlaubte Referenzen zwischen Dokumentklassen ---
# matrix:
#   classes:
#     - {name: contract, paths: [spec/lastenheft.md]}
#     - {name: adr, paths: ["docs/plan/adr/[0-9]*.md"]}
#   rules:
#     - {from: contract, to: adr, allow: false}
#   status: {forbidden: [superseded, deprecated]}
#   exclude-sections: [Historie]

# --- codepaths: explizite Pfade in Inline-Code ---
# codepaths:
#   roots: [docs, spec]            # Wurzel-Präfixe; ./ und ../ werden immer erkannt

# --- hostpaths: host-lokale absolute Pfade (Maschinen-Layout-Leaks) ---
# hostpaths:
#   prefixes: [home, mnt]          # ersetzt die Default-Präfixliste

# --- spans: Markdown-Span-Artefakte (ungeschlossene Code-Spans u. a.) ---
#   (keine eigenen Optionen; über modules aktivieren)

# --- external: Erreichbarkeit von http(s)-Links — NETZZUGRIFF, opt-in ---
# external:
#   timeout-seconds: 10
#   parallel: 4
`
