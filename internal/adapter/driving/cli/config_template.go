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
# Verfügbar: links, anchors, ids, matrix, codepaths, spans, hostpaths, diagrams, versions, pins, immutable, vcs, commits, planning, tracked, external
# (external ist die einzige Netzwerk-Tür — strikt opt-in; vcs und commits sind
#  git-basiert und brauchen .git + eine Commit-Range — strikt opt-in; tracked ist
#  git-basiert und braucht nur .git (Index, ohne Range) — strikt opt-in; planning
#  ist hermetisch (Roadmap-↔-in-progress-Konsistenz, kein git) — strikt opt-in.)

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
#     # order/direction prüft die Richtung INNERHALB einer geschichteten Klasse
#     # (autoritativste Schicht zuerst; Abwärtsverweis ⇒ matrix-downward):
#     # - name: spec
#     #   paths: [spec/lastenheft.md, spec/spezifikation.md, spec/architecture.md]
#     #   order: [spec/lastenheft.md, spec/spezifikation.md, spec/architecture.md]
#     #   direction: no-downward
#     # token erkennt Referenzen auch als bare ID-Token im Körper (nicht nur Link):
#     # - {name: slice, paths: ["docs/plan/planning/**/slice-*.md"], token: 'slice-\d{3}'}
#   rules:
#     - {from: contract, to: adr, allow: false}
#     # - {from: adr, to: slice, allow: false}   # ADR nennt Slice nur als deklarierte
#     #   # Provenance: <!-- d-check:status-provenance --> auf der Zeile nimmt aus
#   status: {forbidden: [superseded, deprecated]}
#     # opt-in: ablösende Datei darf auf ihr abgelöstes (inaktives) Ziel zeigen
#     # allow-supersede-lineage: true
#     # supersede-fields: [Supersedes, Aenderungstyp]
#   exclude-sections: [Historie]
#   # exempt-paths: ["docs/plan/adr/0001-*.md"]   # immutable Alt-ADRs grandfathern

# --- codepaths: explizite Pfade in Inline-Code ---
# codepaths:
#   roots: [docs, spec]            # Wurzel-Präfixe; ./ und ../ werden immer erkannt
#   exempt-paths: [CHANGELOG.md, "docs/reviews/**"]   # Globs: Dateien ohne codepath-Prüfung (datei-weit, wie ids)
#   ignore-refs: ["tools/altes-skript.sh"]            # Globs: Ziel-Pfade ohne Existenz-Prüfung (referenz-weit; Tombstones entfernter Artefakte)

# --- hostpaths: host-lokale absolute Pfade (Maschinen-Layout-Leaks) ---
# hostpaths:
#   prefixes: [home, mnt]          # ersetzt die Default-Präfixliste

# --- spans: Markdown-Span-Artefakte (ungeschlossene Code-Spans u. a.) ---
#   (keine eigenen Optionen; über modules aktivieren)

# --- diagrams: Kennungs-Existenz in Diagramm-Fences (z. B. mermaid) ---
# diagrams:
#   fences: [mermaid]              # zu öffnende Diagramm-Sprachen (Default mermaid)
#   patterns:                      # wie ids; geprüft wird Existenz in defined-in (nicht Linkpflicht)
#     - regex: 'ARC-\d{2}'
#       defined-in: spec/architecture.md

# --- versions: Versions-Pin-Konsistenz (alle Pins == aktuelle Version) ---
# versions:
#   pin-pattern: 'ghcr\.io/[^\s:]+:(v[0-9]+\.[0-9]+\.[0-9]+)'   # Version in Capture-Gruppe 1
#   current-from: version.md#aktuell   # Datei#Anker (Span) mit der aktuellen Version
#   exempt-paths: [CHANGELOG.md, "docs/plan/planning/done/**"]  # historische Pins (datei-weit)

# --- pins: Content-Pin gegen inhaltlichen Drift verlinkter Ziele ---
#   (keine eigenen Optionen; über modules aktivieren. Ein Link mit
#    <!-- dpin: sha256:… --> wird gegen den Hash seines rohen Ziel-Spans geprüft.)

# --- immutable: Immutabilitäts-Pin gegen Core-Drift — hermetisch (kein git) ---
# immutable:
#   exclude-sections: [Geschichte]   # Abschnitte, die nicht zum gehashten Core zählen
#   # Eine Datei mit <!-- immutable: sha256:… --> wird gegen ihren normalisierten
#   # Core (ohne Marker-Zeile + exclude-sections) gehasht; Abweichung → core-drift.

# --- vcs: git-Diff-Immutabilität des Core über eine Commit-Range — git, opt-in ---
#   (braucht .git + eine Range; Aufruf über das print-mk-Target doc-immutable
#    bzw. --range <base>..<head> / --staged. NICHT in modules: oben aufnehmen.)
# vcs:
#   paths: ["docs/plan/adr/[0-9]*.md"]                 # geschützte Datei-Klasse (Glob)
#   immutable-when: '^\*\*Status:\*\* Accepted'        # BASE ab dieser Zeile immutabel
#   exclude-sections: [Geschichte]                     # nicht zum Core zählende Abschnitte
#   status-line: '^\*\*Status:\*\*'                    # Kopf-Status-Zeile (aus dem Core gestrippt)
#   head-allow: '^\*\*Status:\*\* (Accepted|Superseded by ADR-[0-9]{4})'  # erlaubter Status-Übergang

# --- commits: Traceability-Kennung in Commit-Messages über eine Range — git, opt-in ---
#   (braucht .git + eine Range bzw. --commit-msg; Aufruf über das make-Target
#    trace-check bzw. --range <base>..<head> / --commit-msg <datei|->. NICHT in modules: oben.)
# commits:
#   id-patterns:                                       # Regex-Liste gültiger Kennungen
#     - 'ADR-\d{4}'
#     - 'DC-(FA-[A-Z]+|QA)-\d+'
#     - 'slice-\d+'
#   exempt-pattern: '^(Merge |Revert )'                # kennungsfreier Betreff (Merge/Revert)

# --- planning: Roadmap-↔-in-progress-Lifecycle-Konsistenz — hermetisch (kein git), opt-in ---
#   (Aufruf über das make-Target planning-check bzw. --enable planning. NICHT in modules: oben.)
# planning:
#   roadmap: docs/plan/planning/in-progress/roadmap.md   # Roadmap-Datei; ihr Verzeichnis = Slice-Verzeichnis
#   # heading: "## Aktuelle Welle"      # kanonische H2 (Default); fehlt sie ⇒ planning-drift (fail-closed)
#   # marker: "Keine aktive Welle"      # literaler Ruhe-Marker (Default)
#   # slice-glob: "slice-*.md"          # Basisnamen-Glob der Slice-Dateien (Default)

# --- tracked: Getrackt-Status aufgelöster Link-/Bild-Ziele — git-Index (ohne Range), opt-in ---
#   (--enable tracked; fail-closed ohne lesbares .git. Ein existierendes, aber
#    untracked/gitignoriertes Ziel wäre auf jedem frischen Klon target-missing.)
# tracked:
#   exempt-targets: []   # Globs über den AUFGELÖSTEN Zielpfad — absichtlich untrackte Ziele (referenz-weit)

# --- external: Erreichbarkeit von http(s)-Links — NETZZUGRIFF, opt-in ---
# external:
#   timeout-seconds: 10
#   parallel: 4
`
