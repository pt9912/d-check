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
# Verfügbar: links, anchors, ids, matrix, codepaths, spans, hostpaths, diagrams, versions, pins, immutable, vcs, commits, planning, tracked, targets, citations, structure, external, sources
# (external und sources sind die einzigen Netzwerk-Türen — beide strikt opt-in
#  (external prüft http(s)-Erreichbarkeit, sources Upstream-Content-Drift); vcs und commits sind
#  git-basiert und brauchen .git + eine Commit-Range — strikt opt-in; tracked ist
#  git-basiert und braucht nur .git (Index, ohne Range) — strikt opt-in; planning
#  und targets sind hermetisch (Roadmap-↔-in-progress- bzw. Doku-↔-Makefile-
#  Konsistenz, kein git) — strikt opt-in.)

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
#   # closure:                          # zweite Fähigkeit: Struktur der Closure-Notizen (opt-in im opt-in)
#   #   dir: docs/plan/planning/done    # Aktivierungs-Schalter; leer ⇒ inert (keine Slice-Datei wird geöffnet)
#   #   glob: '*.md'                    # EIGENER Kandidaten-Filter; WEGLASSEN ⇒ es gilt slice-glob (Verweis, kein Literal)
#   #   heading-pattern: '^#{2,3} .*Closure-Notiz'  # RE2 (Default); kein Treffer ⇒ closure-note-missing
#   #   min-sentences: 4                # Satzenden ausserhalb Fenced- UND Inline-Code, nur vor Whitespace; < 1 ⇒ Exit 2
#   #   boilerplate: []                 # Floskeln, case-insensitiv an WORTGRENZEN; Default LEER (keine Sprach-Annahme)
#   #   placeholder: false              # unausgefuellte Vorlagen-Platzhalter (<feld>); Default AUS, ignoriert Inline-Code

# --- tracked: Getrackt-Status aufgelöster Link-/Bild-Ziele — git-Index (ohne Range), opt-in ---
#   (--enable tracked; fail-closed ohne lesbares .git. Ein existierendes, aber
#    untracked/gitignoriertes Ziel wäre auf jedem frischen Klon target-missing.)
# tracked:
#   exempt-targets: []   # Globs über den AUFGELÖSTEN Zielpfad — absichtlich untrackte Ziele (referenz-weit)

# --- structure: Struktur-Invarianten INNERHALB eines Dokuments — hermetisch, opt-in ---
#   (--enable structure; Post-Pass. Jede Regel benennt ihre Dateien SELBST über
#    eigene Globs — unabhängig von scan.roots/scan.ignore, daher kein scope.
#    Leere Liste ⇒ Modul inert. Null Kandidaten ⇒ section-missing, auch wenn erst
#    exempt-paths die Menge geleert hat.)
# structure:
#   - files: "docs/plan/planning/done/slice-*.md"   # Pflicht; Glob über Wurzel-relative Pfade
#     section: "## 9. Closure-Notiz"                # Klartext, EXAKT inkl. #-Folge …
#     # section-pattern: '^#{2,3} .*Closure'        # … oder RE2; genau eines von beiden
#     # sections: one                               # one (Default) | each
#     non-empty: true                               # sonst section-empty
#     # min-sentences: 4                            # sonst section-thin; abwesend = Bedingung aus
#     # max-tasks: 0                                # sonst section-oversized (Task-Items IM Abschnitt)
#     # forbid-pattern: 'TODO'                      # Treffer ⇒ section-forbidden
#     # require-pattern: 'Beleg'                    # kein Treffer ⇒ section-pattern-missing
#     # require-all: ["Beleg", "Lernsignal"]        # fehlende Marke ⇒ section-marker-missing
#     # exempt-paths: []                            # Globs; Treffer werden von DIESER Regel nicht geprüft

# --- targets: Deklarations-Konsistenz Doku ↔ Build-Targets — hermetisch (kein git), opt-in ---
#   (Aufruf über das make-Target gate-consistency bzw. --enable targets. NICHT in modules: oben.)
# targets:
#   makefiles: [Makefile]                        # Regelnamen-Quelle(n)
#   doc-tables: [AGENTS.md, harness/README.md]   # Dateien mit make-X-Tabellen (Richtung 1 ⇒ gate-phantom)
#   authority: AGENTS.md                         # Vollständigkeits-Quelle (Richtung 2 ⇒ gate-undocumented)
#   exempt-targets: []                           # Regelnamen EXAKT (kein Glob, anders als tracked) — Utility-Targets ohne Doku-Pflicht

# --- external: Erreichbarkeit von http(s)-Links — NETZZUGRIFF, opt-in ---
# external:
#   timeout-seconds: 10
#   parallel: 4

# --- sources: Upstream-Content-Drift externer Quellen — NETZZUGRIFF, opt-in ---
#   (zweite Netz-Tür neben external. Holt eine auf sha256 gepinnte http(s)-Quelle,
#    hasht sie und meldet source-drift (voller Ist-Hash) bzw. source-unreachable.
#    Pin per Marker am Link — [text](URL) <!-- source-pin: [zip] sha256:… --> — ODER
#    per Liste hier. unpack: zip hasht das Content-Manifest statt der Roh-Bytes.)
# sources:
#   - url: https://example.org/regelwerk.md   # Einzeldatei: Hash der Roh-Bytes
#     sha256: 0000000000000000000000000000000000000000000000000000000000000000
#   - url: https://example.org/bundle.zip     # Archiv: Hash des Content-Manifests
#     sha256: 0000000000000000000000000000000000000000000000000000000000000000
#     unpack: zip

# --- trace: konfigurierbare Quellen der Requirements Traceability Matrix (--trace) ---
#   (KEIN Modul — steuert nur --trace / --require-complete. Jedes Feld optional;
#    abwesend ⇒ d-checks Konventions-Default ⇒ RTM byte-identisch. Eine gesetzte
#    file-pattern braucht eine Capture-Gruppe (…) = Owner-Kennung, sonst Exit 2.)
# trace:
#   requirements:
#     source: spec/lastenheft.md   # nichtleer explizit => 0 Anforderungen ist Exit 2
#     id-pattern: '[A-Z][A-Z0-9]*-(?:FA-[A-Z]+|QA)-\d+[A-Za-z]?'   # Kennung als Ganz-Token/-Zelle UND als Referenz
#     format: headings             # headings (Default) oder table
#     # table:                     # Pflicht bei format: table; exakte Header-Namen
#     #   id-column: Kennung
#     #   text-column: Anforderung
#     #   # text-columns: [Anforderung, Akzeptanzkriterium]  # alternativ zu text-column
#     #   modality-column: Prioritaet  # optional; sonst klassifiziert modality die Textspalte
#     #   duplicate-ids: error      # error (Default), first oder last
#     modality:                    # RFC-2119-Klassifikation je Anforderung (eigene Spalte); {} = Built-in DE+EN-Defaults
#       # levels:                  # Stufe -> Modal-Verb-Keywords (überschreibt die Defaults; kein Keyword in 2 Stufen)
#       #   must:   [MUSS, MUESSEN, "DARF NICHT", MUST, SHALL]
#       #   should: [SOLLTE, SHOULD]
#       #   may:    [KANN, "MUSS NICHT", MAY]
#       require-levels: [must]     # nur diese Stufen gaten --require-complete (unknown nur wenn gelistet)
#   adrs:
#     dir: docs/plan/adr
#     file-pattern: '^(\d{4})-.*\.md$'          # Capture 1 = Owner-Kennung
#     id-prefix: 'ADR-'
#   slices:
#     dir: docs/plan/planning
#     file-pattern: '^(slice-\d+)-.*\.md$'      # z. B. '^(\d+)-.*\.md$' für NNN-titel.md
#     id-prefix: ''                             # der Owner-Kennung vorangestellt (Default leer)
#   coverage:                                   # kuratierte Coverage-Quellen (Liste; eigene RTM-Spalte)
#     - files: [docs/plan/traceability.md]      # EXPLIZITE Dateien (keine dir/pattern) — gegen ADR-Kontamination
#       label: 'Trace'                          # Owner-Kennung in der Coverage-Spalte
#       ranges: true                            # GG-QA-001..006 -> alle sechs (+ /-Aufzählung); gegen id-pattern validiert
#       # sections: ["27.1 Anforderung zu Design"]              # Whitelist: NUR diese Abschnitte (voller Heading-Text!)
#       # exclude-sections: ["27.1.1 Anforderungen ohne Design-Artefakt"]  # Blacklist (voller Heading-Text)
#   # cross-consistency:                          # Mengenabgleich zweier Traceability-Sichten (advisory; gatet via --require-complete)
#   #   forward:                                  # Vorwärts-Sicht: Anforderung -> Design-Menge (der kuratierte Spiegel, der driftet)
#   #     file: docs/plan/traceability.md
#   #     sections: ["27.1 Anforderung zu Design"]   # optional; Whitelist (voller Heading-Text!)
#   #     # exclude-sections: ["27.1.1 Anforderungen ohne Design-Artefakt"]
#   #     req-column: Anforderung                 # header-gebunden; je relevanter Tabelle genau einmal
#   #     design-column: Design-Artefakte         # header-gebunden
#   #     design-pattern: 'GG-AR-[A-Z-]+\d*'      # extrahiert die Artefakt-IDs; GETEILT mit backward (gemeinsamer Namensraum!)
#   #     req-pattern: 'GG-[A-Z]+-\d{3}'          # erkennt die Anforderungs-IDs der req-column; Default: requirements.id-pattern
#   #                                             # ACHTUNG: der Vergleichs-Scope ist NICHT die RTM-Anforderungsmenge — scopt Ihre
#   #                                             # RTM eine Familie aus, setzen Sie ihn hier explizit, sonst ist die Vorwärts-Sicht leer
#   #     ranges: true                            # GG-SIM-001..009 der ID-Spalte expandieren
#   #   backward:                                 # Rück-Kanten: Design -> Anforderung (Quelle der Wahrheit)
#   #     file: spec/architecture.md
#   #     sections: ["4 Komponenten", "5 Ports"]  # optional; Whitelist
#   #     artifact-id-column: first               # 'first' = erste Spalte (heterogene ID-Header) ODER ein Header-Name
#   #     edge-column: Bezug                      # header-gebunden; macht die Tabelle relevant
#   #     req-pattern: 'GG-[A-Z]+-\d+'            # erkennt die Anforderungs-IDs in der Bezug-Zelle
#   #     ranges: true
#   #   mode: equal                               # equal (Default: F\B und B\F gaten) oder superset (nur B\F)
#   #   exclude-req: '^GG-SPEC-'                  # Ableitungssprünge (Mittelschicht-IDs) aus dem Abgleich nehmen
`
