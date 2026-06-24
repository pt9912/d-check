# Changelog

Alle nennenswerten Änderungen an diesem Projekt werden in dieser Datei
dokumentiert. Das Format folgt [Keep a Changelog](https://keepachangelog.com/de/1.1.0/),
die Versionierung folgt [SemVer](https://semver.org/lang/de/).

## [0.28.0] — 2026-06-24

### Added

- slice-048 — neues opt-in Regelmodul `versions` (zehntes Modul): prüft, dass alle
  gepinnten `ghcr`-Image-Verweise die aktuelle Version aus `version.md#aktuell`
  tragen, sonst Befund `version-stale`; liest die Pins **auch in Fenced-Code**
  (gescopte Fence-Ausnahme), Ventile `exempt-paths`/`d-check:ignore`, fail-closed
  bei unauflösbarer Quelle, diagnose-only
  ([`DC-FA-VER-001`](spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
  [ADR-0019](docs/plan/adr/0019-versions-pin-fence-ausnahme.md), Lastenheft 0.28.0).
- Release-Register `version.md` (only-current-anchor): kanonische Quelle der
  aktuellen Version; `--print-config` führt den `versions:`-Block.

### Changed

- Dogfooding: `.d-check.yml` aktiviert `versions` — die `ghcr`-Image-Pins in
  README und Benutzerhandbuch sind ab jetzt gateguarded.

## [0.27.0] — 2026-06-23

### Added

- slice-047 — `--print-mk`-Fragment um drei Targets + eine Variable erweitert:
  `doc-doctor` (`--doctor`-Diagnose), `doc-repair` (`--repair`-Patch, Recipe-Echo
  unterdrückt für `git apply`-reine stdout), `doc-help` (namespaced Self-Doku der
  `doc-*`-Targets via `##`-Annotationen) sowie `DCHECK_DIGEST` (Digest-Override per
  `ifeq`, sticht den Tag von `DCHECK_IMAGE`). Alle Targets `##`-annotiert
  ([`DC-FA-CLI-010`](spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
  Change Request, Lastenheft 0.27.0).

## [0.26.0] — 2026-06-23

### Changed

- slice-046 — `--suggest-config ai-harness[-init]`: die Ausgabe nennt die nicht
  aktivierten situativen opt-in-Module (`external`, `spans`, `hostpaths`,
  `diagrams`) jetzt in einem Kommentar mit Verweis auf `d-check --print-config`
  (Auffindbarkeit ohne Aktivieren eines inerten Moduls — `diagrams` braucht
  repo-spezifische `patterns`/`defined-in`). Schärfung
  ([`DC-FA-CLI-006`](spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten),
  Lastenheft 0.26.0).

## [0.25.0] — 2026-06-23

### Added

- slice-045 — Modul `diagrams` (opt-in): öffnet gezielt benannte
  Diagramm-Fences (Default `mermaid`) und prüft die darin gefundenen Kennungen
  auf **Existenz** in ihrer `defined-in`-Quelle (Befund `diagram-id-undefined`).
  Reine Token-Extraktion ohne Mermaid-Parser, **Existenz statt Link-Policy** (in
  Fences kein Markdown-Link möglich), read-only/netzlos (`DC-QA-03`),
  deterministisch (`DC-QA-02`), Default aus (byte-identisch). Fängt Drift/Typos
  in Diagramm-Kennungen (z. B. Architektur-IDs in `mermaid`), die bei opaken
  Fences heute unsichtbar bleiben
  ([`DC-FA-DIAG-001`](spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in),
  [ADR-0018](docs/plan/adr/0018-diagram-fence-ausnahme.md)).

## [0.24.0] — 2026-06-23

### Added

- slice-044 — Option `--require-complete` (nur mit `--trace`): bindet die
  RTM-Waisen-Markierung an den Exit-Code — ≥1 Requirements-Waise ⇒ **Exit 1**
  statt 0, sonst 0; die RTM bleibt unverändert auf stdout, der Default-`--trace`
  bleibt advisory (Exit 0). Erlaubt Konsumenten ein Vollständigkeits-Gate im
  eigenen Makefile, ohne Parsing-Logik zu kopieren. read-only (`DC-QA-03`),
  deterministisch (`DC-QA-02`)
  ([`DC-FA-CLI-011`](spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)).

### Changed

- slice-044 — `--print-mk`-Fragment: zusätzlich die Targets `doc-trace`
  (advisory RTM) und `doc-complete` (`--trace --require-complete`, das
  Vollständigkeits-Gate) plus eine überschreibbare `TRACE_FLAGS`-Variable
  ([`DC-FA-CLI-010`](spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
  Change Request 0.24.0).

## [0.23.0] — 2026-06-22

### Changed

- slice-043 — Modul `codepaths`: neues Ventil `exempt-paths` (Glob-Liste,
  Syntax wie `scan.ignore`) nimmt **ganze Dateien** von der codepath-Prüfung
  aus — Parität zum gleichnamigen `ids`-Ventil; datei-weit, komplementär zum
  `d-check:ignore`-Marker. Abwärtskompatibel: ohne `exempt-paths`
  byte-identisch ([`DC-QA-02`](spec/lastenheft.md#dc-qa-02--determinismus)).
  Dogfooding: die eigene `.d-check.yml` nimmt `docs/reviews/**` aus
  ([`DC-FA-CODE-001`](spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
  Change Request 0.23.0).

## [0.22.0] — 2026-06-22

### Added

- slice-038 — Option `--print-mk`: gibt ein include-bares `d-check.mk`
  (überschreibbare `DCHECK_IMAGE`-Variable mit version-gepinntem Image +
  `doc-check`-Target) auf stdout aus — Konsumenten `include`-n statt
  Recipe/Skript zu kopieren; der Image-Ref ist die ins Binary eingebettete
  Release-Version (Digest via `DCHECK_IMAGE`-Override). read-only
  (`DC-QA-03`), deterministisch (`DC-QA-02`)
  ([`DC-FA-CLI-010`](spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)).
- slice-036 — Option `--trace`: gibt eine **Requirements Traceability
  Matrix** (je Anforderung die referenzierenden ADRs/Slices + Waisen-
  Markierung) auf stdout aus — Default Markdown-Tabelle, mit `--json`/`--yaml`
  maschinenlesbar; read-only (`DC-QA-03`), deterministisch (`DC-QA-02`),
  kein Dokument erzeugt; Doku-Domäne (Lastenheft/ADR/Planning)
  ([`DC-FA-CLI-009`](spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)).

### Changed

- slice-037 — `--suggest-config ai-harness[-init]`: das Anforderungs-`ids`-
  Muster ist nicht mehr fix `DC-`. Neues Flag `--id-prefix <PREFIX>` setzt
  das Präfix explizit; der Modus `ai-harness` leitet es aus
  `spec/lastenheft.md` ab (mehrere verschiedene Präfixe ⇒ Fehler).
  **Breaking:** ohne Präfix/Ableitung (typisch `ai-harness-init` im leeren
  Repo) erscheint ein markierter Platzhalter `<PREFIX>` + `# TODO` statt
  eines stillen `DC-`
  ([`DC-FA-CLI-006`](spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten),
  `ADR-0015`).
- slice-034 — Distribution: der `:latest`-Tag zeigt explizit auf das
  **neueste stabile** Release; Vorabversionen (Prereleases) erhalten kein
  `:latest`. `:latest` ist Komfort-Einstieg — für reproduzierbare Läufe
  weiterhin auf eine feste Version oder den `@sha256:`-Digest pinnen
  (ratifiziert `ADR-0002` §4; `ADR-0014`).

## [0.19.0] — 2026-06-20

### Added

- slice-031 — YAML-Ausgabeformat `--yaml`: gibt die Befunde strukturgleich
  zu `--json` als YAML auf stdout aus (`findings`/`summary`/`exitCode`);
  `--doctor --yaml` analog `--doctor --json`. Deterministisch (`DC-QA-02`),
  read-only (`DC-QA-03`)
  ([`DC-FA-CLI-004`](spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)).

### Changed

- slice-032 — semgrep als hermetisches Security-/Static-Analysis-Gate in
  `make gates`: gepinntes Image (`semgrep/semgrep:1.167.0`) + gepinntes,
  lokal außerhalb des Repos gecachtes `go/lang/security`-Regelset, netzloser
  Scan (`--network none`); `--error` bricht das Gate bei Befund. Ergänzt
  golangci-lint sprachübergreifend; reproduzierbar (`DC-QA-02`), netzlos
  (`DC-QA-03`) (`ADR-0010`).
- slice-033 — alle Build- und Gate-Images per `@sha256:`-Digest gepinnt
  (alle Dockerfile-`FROM` — golang, golangci-lint, distroless — und das
  semgrep-Image; Manifest-Listen-Digest amd64+arm64, inline neben dem Tag);
  `make versions` belegt die Pins. Schließt die `ADR-0002`-§1-Digest-Drift
  und vereinheitlicht die Image-Pin-Politik (`ADR-0011`).

## [0.18.0] — 2026-06-19

### Added

- slice-030 — `--suggest-config ai-harness` / `ai-harness-init`: schlägt
  ein an die ai-harness-course-Konvention (Baseline v1.3.0) angelehntes
  `.d-check.yml` vor — kanonische `ids`-Muster, `matrix`-Klassen samt
  Referenzrichtung, Standard-Modulset und Scan-Scope. **Zwei Modi:**
  `ai-harness-init` (Voll-Kanon, alle Blöcke aktiv — Zielbild fürs leere
  Repo, läuft nach Struktur-Anlage) und `ai-harness` (repo-bewusst — nur
  vorhandene Pfade aktiv, fehlende auskommentiert mit Hinweis). Read-only,
  advisory, deterministisch (`DC-QA-02`); mit echten Quellen kombinierbar
  ([`DC-FA-CLI-006`](spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)).

### Changed

- Spezifikation/Dokumentation präzisiert: `--repair-broad` löst nur
  **Verschiebungen** auf (gleicher Name), keine Umbenennungen; das
  Out-of-Scope von `DC-FA-CLI-008` benennt die nicht-reparierbaren
  Befundarten und schließt git-historienbasierte Move-/Rename-Erkennung
  aus; die Reparatur-Ableitbarkeit ist als Entscheidung in `ADR-0008`
  festgehalten (Handbuch §4.10).

## [0.17.0] — 2026-06-19

### Added

- slice-029 — Maschinenlesbare Diagnose `--doctor --json`: dieselbe
  Diagnose wie `--doctor`, aber als JSON-Dokument auf stdout. Die
  `findings` tragen je Eintrag zusätzlich `reasonText` (Grund-Klartext)
  und `fixCandidate` (`{original, replacement, note}` oder explizit
  `null`, wo kein eindeutiger Fix existiert); `summary`/`exitCode` wie
  bei `--json`. Ein drittes Rendering desselben Fix-Kandidaten-Modells
  neben Prosa (`--doctor`) und Patch (`--repair`); deterministisch
  (`DC-QA-02`), read-only (`DC-QA-03`).

### Changed

- `--doctor` ist nun **mit `--json` kombinierbar** (zuvor Nutzungsfehler).
  Nutzungsfehler bleiben nur `--repair`+`--json` und `--doctor`+`--repair`
  ([`DC-FA-CLI-007`](spec/lastenheft.md#dc-fa-cli-007--diagnose-modus),
  [`DC-FA-CLI-004`](spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)).

## [0.12.0] — 2026-06-18

### Added

- slice-025 — Diagnose-Modus `--doctor`: erklärende, nach Datei
  gruppierte Klartext-Diagnose auf stdout statt der knappen Befund-Zeilen,
  mit Fix-Kandidaten wo eindeutig ableitbar (in dieser Version
  `id-unlinked` → Markdown-Link auf das ids-Definitions-`target`). Grund-
  Klartext für alle 14 Grund-Codes, abgesichert durch eine
  Vollständigkeits-Prüfung gegen die Reason-Konstanten. Read-only,
  stdout-only; `--doctor` ist nicht mit `--json` kombinierbar
  (Nutzungsfehler, Exit 2). Das Fix-Kandidaten-Modell ist die
  wiederverwendbare Eingabe für den folgenden Patch-Modus `--repair`
  (slice-026)
  ([`DC-FA-CLI-007`](spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)).
- slice-026 — Reparatur-Patch `--repair`: gibt einen unified diff auf
  stdout aus (`git apply`-kompatibel), der ableitbare Befunde behebt;
  schreibt selbst nichts. **Konservativ** (Default) nur eindeutige Fixes
  (`id-unlinked` → Definitions-Link, nur nackte Prosa-Vorkommen — keine
  Inline-Code- oder Mehrdeutigkeits-Reparatur); **breit** (`--repair-broad`,
  opt-in) zusätzlich Best-Guess (`target-missing` → Datei eindeutig
  gleichen Basisnamens), review-pflichtig mit Marker auf stderr, sodass
  der Patch `git apply`-rein bleibt. Nicht mit `--json`/`--doctor`
  kombinierbar; deterministisch (`DC-QA-02`), read-only (`DC-QA-03`).
  Wiederverwendung des Fix-Kandidaten-Modells aus slice-025
  ([`DC-FA-CLI-008`](spec/lastenheft.md#dc-fa-cli-008--reparatur-patch)).

### Changed

- slice-027 — `make image-test` deckt nun auch `--doctor`/`--repair` ab
  (nativ == Container byte-identisch, stdout + stderr + Exit-Code); E2E-
  Härtung des
  [`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image)-
  Vertrags für die neuen Ausgabe-Modi (`DC-QA-02`-Parität).

## [0.11.0] — 2026-06-17

### Added

- slice-024 — Modul `matrix`: opt-in `allow-supersede-lineage` (mit
  `supersede-fields`) nimmt die Supersede-Lineage-Kante von der
  Status-Prüfung aus — eine ablösende Datei darf auf das von ihr
  abgelöste (inaktive) Dokument verweisen, ohne `matrix-inactive` zu
  erzeugen, sofern sie die Ablösung über ein deklariertes Feld benennt
  (Match über Linktext bzw. Zielpfad der Referenz)
  ([`DC-FA-MTX-001`](spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
  Change Request 0.14.0). Wirkt nur auf `matrix-inactive`;
  `matrix-forbidden` (Klassen-Regeln) bleibt unberührt. `matrix` trägt
  bewusst **keinen** `d-check:ignore`-Marker — legitime Ausnahmen sind
  strukturell (`exclude-sections`, `allow-supersede-lineage`).
  Abwärtskompatibel: Default aus ⇒ Befundsatz byte-identisch.

## [0.10.0] — 2026-06-16

### Changed

- slice-023 — die `ids`-Ventile `exempt-paths` und `d-check:ignore`
  gelten jetzt für **alle** Vorkommen eines Musters — nackt im Fließtext
  **wie** in Inline-Code — und unabhängig von der `link-policy`
  ([`DC-FA-ID-001`](spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
  Change Request 0.13.0). Bisher griffen beide nur auf die
  `always`-Inline-Code-Vorkommen; eine nackte Kennung in einer
  `exempt-paths`-Datei (oder auf einer `d-check:ignore`-Zeile) wurde
  weiterhin als `id-unlinked` gemeldet. Jetzt ein Ganzdatei- bzw.
  Ganzzeilen-Carve-out. Abwärtskompatibel: Configs ohne gesetzte Ventile
  sind byte-identisch; die Wirkung geht nur in Richtung *weniger* Befunde
  in explizit ausgenommenen Dateien/Zeilen.

## [0.9.0] — 2026-06-15

### Added

- slice-022 — Inline-HTML-Anker als gültige Anker-Menge
  ([`DC-FA-ANCH-001`](spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors)
  Schärfung 0.12.0): das Modul `anchors` (und mittelbar `codepaths`)
  akzeptiert zusätzlich zu Heading-Slugs die Inline-HTML-Anker der
  Zieldatei — `id` an beliebigem Element und `name` an `<a>`
  (GitHub-Parität, wörtlicher/case-sensitiver Vergleich). HTML in
  Code-Auszeichnung (Fenced-Block oder Inline-Code) zählt nicht.
  Abwärtskompatibel: reduziert Falsch-Befunde `anchor-missing`, erzeugt
  nie neue.

## [0.8.0] — 2026-06-13

Reichhaltige `--help` (Schärfung `DC-FA-CLI-001`, slice-021): die Hilfe
nennt Synopsis und das Pfad-Argument und verweist fürs Config-Format
auf `--print-config`.

### Changed

- slice-021 — reichhaltige `--help`/`-h`
  ([`DC-FA-CLI-001`](spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel)
  Schärfung 0.11.0): die Hilfe nennt jetzt eine Kurzbeschreibung, die
  Synopsis `d-check [optionen] [pfad]` und beschreibt das bislang
  verschwiegene Pfad-Argument (Scan-Wurzel, Default cwd); für das
  Config-Format verweist sie auf `--print-config`/`--suggest-config`
  (kein Format-Duplikat). Exit 0 / stderr unverändert.

## [0.7.0] — 2026-06-13

Konfigurations-Vorschlag aus Autoritäts-Dokumenten (`--suggest-config`,
Change Request 0.10.0, slice-020) — inkl. Review R1.

### Added

- slice-020 — Option `--suggest-config`
  ([`DC-FA-CLI-006`](spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten),
  Change Request 0.10.0): liest benannte Autoritäts-Quellen und schlägt
  ein `.d-check.yml` vor — leitet je Quelle ein `ids`-Muster aus den in
  Überschriften **definierten** Kennungen ab (Präfix-Alternation,
  Round-Trip-Garantie; Quell-Kennungen als Kommentar) und schlägt
  opt-in-Module nach Signal vor. **Liest, schreibt nie** (read-only-
  Vertrag; Umleiten macht der Aufrufer). Bewusste Grenze: erkennt nur
  großgeschriebene Heading-Token-IDs — Scaffold, kein Orakel, der
  Mensch verengt/ergänzt. Dazu Schärfung von
  [`DC-FA-ID-001`](spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids):
  Muster-Ableitung bleibt für die **Prüfung** Out-of-Scope, nur der
  advisory Scaffold-Modus leitet ab. Korpora-Gegentest dokumentiert
  (d-check round-trippt, b-trace zeigt die Heading-Grenze).

### Fixed

- Review R1 zu slice-020 (`/code-review`): das `--suggest-config`-Gerüst
  nimmt jetzt `ids` in die Modul-Liste auf (sonst waren die abgeleiteten
  Muster im erzeugten Config inaktiv — gültiges YAML, semantisch
  wirkungslos); Modul-Probelauf nutzt denselben Scope (`roots: ["."]`)
  wie das Gerüst; `target` wird gequotet (Quellpfade mit `:`/`#` brechen
  das YAML nicht mehr); Heading-Token-Extraktion strippt Links und
  Satzzeichen (`ADR-0001:` wird erkannt); leere Quellenliste ist ein
  Nutzungsfehler. Report unter `docs/reviews/`.

## [0.6.0] — 2026-06-13

Konfigurations-Startgerüst (`--print-config`, Change Request 0.9.0,
slice-019): neue Repos kommen ohne Handarbeit zu einer `.d-check.yml` —
das Werkzeug gibt aus, der Aufrufer leitet um; read-only bleibt.

### Added

- slice-019 — Option `--print-config`
  ([`DC-FA-CLI-005`](spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben),
  Change Request 0.9.0): gibt ein kommentiertes `.d-check.yml`-Startgerüst
  auf stdout aus und endet mit Exit 0 — **kein Scan, schreibt nichts**
  (read-only-Vertrag bleibt; Anlegen via `d-check --print-config >
  .d-check.yml`). Das Gerüst ist statisch, deterministisch und
  dekodiert über den eigenen Parser; es macht die verfügbaren Module
  und Optionen als Kommentare sichtbar. Senkt die Adoptions-Reibung in
  neuen Repos ohne Konfiguration.

## [0.5.0] — 2026-06-13

Konfigurierbare Link-Politik für das Modul `ids` (Change Request 0.8.0,
slice-018): „gut verlinkte Dokumente" wird ein im `.d-check.yml`
konfigurierbares, gemessenes Property.

### Added

- slice-018 — konfigurierbare Link-Politik `ids.patterns[].link-policy`
  ([`DC-FA-ID-001`](spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids),
  Change Request 0.8.0): `link-policy: always` macht auch
  Inline-Code-Vorkommen einer Kennung linkpflichtig — „gut verlinkt"
  als gemessenes, konfigurierbares Property statt menschlicher
  Aufmerksamkeit. Default `prose` (byte-identisch, opt-in fürs Gating).
  Zwei Ventile: `exempt-paths` (Glob-Liste je Muster) und der
  Zeilen-Marker `d-check:ignore` (Geltungsbereich von `codepaths` auf
  `ids` erweitert — illustrative Beispiel-IDs). Kalibriert über die
  drei ids-Repos (d-check, u-boot, b-trace); Dogfooding aktiv (d-check
  setzt `always` und verlinkte den eigenen Befundsatz). Nutzersichtbar
  dokumentiert in [`docs/user/operations.md`](docs/user/operations.md).

### Changed

- Dogfooding-Sweep: d-checks eigene Doku auf `link-policy: always`
  umgestellt; Inline-Code-Kennungen in Slices, ADRs, AGENTS, harness
  und Spezifikation als Links ausgeführt (Sektions-Referenzen `.a` auf
  ihre Spez-Anker), `exempt-paths` für CHANGELOG + Reviews,
  Beispiel-IDs mit `d-check:ignore`.

## [0.4.0] — 2026-06-13

Welle-06-sensorik: zwei opt-in-Sensormodule — `spans` (`DC-FA-SPAN-001`,
Markdown-Span-Artefakte) und `hostpaths` (`DC-FA-HOST-001`, host-lokale
absolute Pfade), je über 14 Korpora kalibriert und gegen die Alt-Sensoren
gegengeprüft; schließt welle-06 (slice-015, slice-016).

### Added

- slice-016 — Modul `hostpaths`
  ([`DC-FA-HOST-001`](spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in),
  Change Request 0.6.0, opt-in): meldet host-lokale absolute Pfade
  (Maschinen-Layout-Leaks) in Prosa **und Inline-Code**
  (`hostpath-forbidden`); Unix-Präfixliste konfigurierbar via
  `hostpaths.prefixes` (Default ohne tmp — Lastenheft 0.7.2 aus dem
  Kalibrierungs-Befund: Laufzeit-Doku ist legitim),
  Windows-/UNC-Muster fest (UNC-Servername alphanumerisch — Schutz
  vor Regex-Beispiel-Treffern), Fences ausgenommen, kein
  Opt-out-Marker. Paritäts-Gegentest gegen den
  bess-ems-Rest-Sensor: identische Befunde auf identischen Zeilen;
  Kalibrierung über 14 Korpora trennte echte Workspace-Leaks
  (k-deskflight-Spec gefixt) von gewollter
  Windows-/WSL-Plattform-Doku (Opt-in-Entscheidung der Repos).

- slice-015 — Modul `spans`
  ([`DC-FA-SPAN-001`](spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in),
  Change Request 0.5.0, opt-in): meldet ungeschlossene Code-Spans,
  die an Nicht-Whitespace kleben (`span-unclosed`, absatzweise —
  alleinstehende Backticks bleiben literal) und Link-Syntax im
  Linktext (`span-nested-link`; Badge-Muster `[![…](…)](…)` ist
  legal — Lastenheft 0.7.1 aus dem Kalibrierungs-Befund). Dogfooding
  aktiv; Kalibrierung über 14 Korpora fand 17 echte Artefakte (in
  den Ziel-Repos gefixt), historischer Gegentest: 14 Befunde auf dem
  u-boot-Stand vor den slice-014-Reparaturen, 0 danach.

## [0.3.0] — 2026-06-12

Modul-lokaler Scan-Scope (`<modul>.scope`, Change Request des
Erst-Bedarfsträgers grid-gym) — dazu der dokumentierte Abschluss des
13/13-Migrations-Rollouts.

### Added

- slice-017 — Modul-lokaler Scan-Scope
  ([`DC-FA-CONF-002`](spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope),
  Change Request des Erst-Bedarfsträgers grid-gym): optionaler
  Schlüssel `<modul>.scope` (`roots` Pflicht, `ignore` optional)
  ersetzt für genau dieses Modul den globalen Scan-Scope — eigener
  Discover-Lauf mit den bekannten Scan-Regeln, Lauf über die
  Vereinigungsmenge mit Einmal-Lese-Garantie; ohne `scope`
  byte-identisches Verhalten (belegt gegen v0.2.1 auf zwei Korpora).
  Konsumenten-Abnahme grid-gym: `ids` kuratiert auf `spec/` +
  `docs/user/` → 311 statt 2378 Befunde, `links`/`anchors`
  unverändert global.

- slice-014 — Rollout abgeschlossen
  ([`DC-QA-04`](spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)
  vollständig, 13/13): alle verbleibenden neun Alt-Tool-Vorkommen
  migriert (Shell: b-trace, m-trace, cmake-xray, k-deskflight;
  Python: c-hsm-doc, grid-gym; JS: euler-fourier-hilbert, b-cad;
  eigenständige Linie: bess-ems — Inventur-Nachtrag Lastenheft
  0.4.0). 16 echte Mehr-Befunde in den Ziel-Repos gefixt;
  Rest-Sensoren für Math-Validierung, Host-Pfad-Prosa und
  Modul-Nummern verbleiben dort. Zusatz: Neu-Adoption pkcs11-course
  (Auslöser der v0.2.1-Scan-Härtung). Vergleichstabellen in der
  slice-014-Closure-Notiz; schließt welle-05.

## [0.2.1] — 2026-06-12

Scan-Härtung aus der pkcs11-course-Adoption (slice-014) plus der
dokumentierte Rollout-Stand.

### Fixed

- `scan.ignore`-Muster prunen jetzt den Verzeichnis-Abstieg: ein
  vollständig ignorierter Teilbaum (`pfad/**` oder direkt matchendes
  Muster) wird nicht betreten — unlesbare ignorierte Verzeichnisse
  (z. B. root-eigene Laufzeit-Residuen wie SoftHSM-Tokens) brechen
  den Lauf nicht mehr mit Exit 2 ab
  ([`DC-FA-SCAN-001`](spec/lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln)).
- `SKIP_DIRS` um `.gradle` ergänzt (Parität zur JS-Alt-Familie);
  Spezifikation §3 inkl. Querverweis aus dem Config-Schema.

### Added

- slice-012 — Pilot-Migrationen
  ([`DC-QA-04`](spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)):
  drei Schwester-Repos prüfen ihre Doku jetzt digest-gepinnt über das
  GHCR-Image — d-migrate (Shell-Familie, v0.1.0-Digest, Alt-Skript
  gelöscht), ai-harness-course (JS-Familie, v0.2.0-Digest,
  `docs-check.js` auf den Modul-Nummern-Rest-Sensor geschrumpft) und
  u-boot (Python-Vollausbau, v0.2.0-Digest, `check_refs.py`
  deprecated). Vergleichsläufe und Triage in der Closure-Notiz des
  Slices; schließt welle-04 und Meilenstein M3.

## [0.2.0] — 2026-06-12

Modul `codepaths` (Change Request 0.3.0) und der absatzweise
Inline-Code-Parser aus dem `DC-QA-04`-Gegentest — damit enthält das
Image alle sechs Regelmodule (slice-013, slice-012-Vorlauf).

### Added

- slice-013 — Modul `codepaths` (`DC-FA-CODE-001`, opt-in): explizite
  Pfade in Inline-Code (`./`, `../`, konfigurierte Wurzel-Präfixe via <!-- d-check:ignore (Syntax-Beispiel) -->
  `codepaths.roots`) werden auf Existenz, Repo-Escape und — bei
  Markdown-Zielen mit Fragment — Anker geprüft; Wert-Normalisierung
  (Anführungszeichen, Satzzeichen, `Datei:Zeile`-Suffix), Headings
  ausgenommen, Zeilen-Opt-out `<!-- d-check:ignore (Begründung) -->`
  wirkt nur auf dieses Modul. Dogfooding aktiv (eigene Doku
  befundfrei; 16 begründete Marker an historischen/Beispiel-Pfaden).
- Review-Infrastruktur nach Digest-Welle-18-Konvention:
  Reviewer-Skill (`.harness/skills/reviewer.md` — Kategorien-Anker,
  Output-Schema, Negativbefund-Pflicht) und Report-Ablage
  `docs/reviews/` (ein Report pro Lauf); `.gitignore` auf
  `.harness/state/` verengt (Skills sind versionierte
  Harness-Mechanik).

### Changed

- Lastenheft 0.3.0/0.3.1 — Change Request `DC-FA-CODE-001`: neues
  opt-in-Modul `codepaths` (explizite Pfade in Inline-Code,
  Zeilen-Opt-out `d-check:ignore` nur für dieses Modul) inkl.
  Review-R1-Präzisierungen (Wert-Normalisierung, Anker-Prüfung bei
  Markdown-Zielen); Umsetzung folgt mit slice-013.

### Fixed

- Inline-Code-Erkennung absatzweise statt zeilenweise (CommonMark:
  Spans dürfen Zeilenumbrüche enthalten; Absatzgrenzen sind
  Leerzeile/Fence, ungeschlossene Backtick-Folgen sind literal und
  brechen den Scan nicht mehr ab). Behebt False-Positive-
  `id-unlinked`-Befunde auf korrekt verlinkten Kennungen nach
  Span-Fortsetzungszeilen — gefunden im
  [`DC-QA-04`](spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Gegentest
  gegen den u-boot-Vollausbau (slice-012); Spezifikation
  §`DC-FA-LINK-001.a` Schritt 2 fortgeschrieben.

## [0.1.0] — 2026-06-11

Erster Release: alle fünf Regelmodule, Gate-Vollausbau, Distribution
via `ghcr.io/pt9912/d-check` (slice-011 — `DC-FA-DIST-001`).

### Added

- GHCR-Release-Pipeline
  ([`.github/workflows/release.yml`](.github/workflows/release.yml)):
  Tag-Push `v*` → SemVer-Validate → `make ci` → OCI-Label-Pin
  (`org.opencontainers.image.version` muss der Tag-Version
  entsprechen) → Push mit Semver-Tag (+ `latest` nur für stabile
  Releases) → Digest-Pin im Job-Summary und GitHub-Release;
  Betriebs-/Release-Doku unter `docs/user/` (löst `MR-009`:
  `docs/user` ist jetzt Rang 6 der Source Precedence).
- MIT-Lizenz als Repository-Lizenz ergänzt ([`LICENSE`](LICENSE)).
- Lastenheft 0.1.0 ([`spec/lastenheft.md`](spec/lastenheft.md)):
  Konsolidierung von 12 Quell-Tools, Modul-Schnitt (`links`, `anchors`,
  `ids`, `matrix`, `external`), Docker-Distribution.
- Greenfield-Harness-Bootstrap nach AI-Harness-Kurs: `AGENTS.md`,
  `harness/`, Planning-Struktur mit Roadmap und Slices 001–004,
  Makefile-Gates (`doc-check`, `gates`), `.claude`-Hooks.
- Fundament-ADRs 0001–0004 (slice-001): Go, GHCR-Image auf
  distroless/static (Binary-Distribution vertagt mit Trigger),
  striktes YAML via yaml.v3, Architektur Hexagon light.
- Spec-Straten 2+3 (slice-002): `spec/spezifikation.md` (Prüflauf-,
  Slug-, Modul-Algorithmen; `--json`- und `.d-check.yml`-Schema;
  Defaults; Grund-Codes) und `spec/architecture.md` (Hexagon-Schnitt,
  Import-Constraints als arch-check-Grundlage, Sequenzen).

- slice-003 — erster Go-Code: CLI-Kern (`d-check [pfad]`,
  `--enable/--disable/--json`), Scanner mit Default-Wurzeln/Ignores,
  Modul `links` (Linkziele, Repo-Escape, Symlink-Vorrang,
  RFC-3986-Dekodierung), strikte `.d-check.yml`-Validierung,
  Text-/JSON-Reporter; Layout nach `ADR-0005` (hexagon-/adapter-Ordner,
  u-boot-Konvention); Dockerfile-Stages + Make-Gates `lint`, `test`,
  `arch-check` (Fitness Function R1–R5), Runtime-Image
  distroless/static mit Selbst-Smoke-Test (`make run`).

- slice-004 — Modul `anchors` (GitHub-Slug-Verfahren inkl.
  Duplikat-Suffixen, Fragment-Dekodierung, Schweigen bei fehlender
  Zieldatei) und **Dogfooding**: `make doc-check` läuft über `d-check`
  selbst (`scan.roots: ["."]`, Module links+anchors — erstmals mit
  Anker-Validierung); vendorter Bootstrap-Sensor gelöscht
  (`MR-003` → `MR-007` aufgelöst); Vergleichslauf als erster
  `DC-QA-04`-Datenpunkt.

- slice-005 — SOLID-nahes Lint-Profil (`ADR-0006`, u-boot-Parität ohne
  depguard): 5 Default- + 23 Linter mit u-boot-Kalibrierung,
  gomodguard-Anti-Module, Why-kommentierte Ausnahmen; Code-Refactoring
  statt Carveouts (Globals → Funktionen, Komplexitäts-Splits in
  cli/configyaml/core) — lint-clean ohne //nolint.

- slice-006 — Modul `ids` (`DC-FA-ID-001`): Linkpflicht für Kennungen
  nach konfigurierten Regex-Mustern (Reihenfolge = Präzedenz, erstes
  Match gewinnt pro Vorkommen); „verlinkt" = Vorkommen im Linktext
  eines Markdown-Links, Ziel-Klammern und Bildreferenzen sind
  linkpflichtfrei (kein Fließtext); Grund-Code `id-unlinked`;
  Config-Constraint `ids.patterns[].target` muss existieren (Exit 2).

- slice-007 — Modul `matrix` (`DC-FA-MTX-001`): Dokumentklassen per
  Glob (Reihenfolge = Präzedenz), Referenzregeln pro Klassen-Paar,
  Status-Bedingungen (`**Status:**`-Zeile vor Status-Heading,
  Präfix-Match case-insensitiv, ohne Status aktiv) und
  `exclude-sections` (Provenance-Ausnahme); Grund-Codes
  `matrix-forbidden`/`matrix-inactive`. **Dogfooding-
  Selbstkonfiguration:** die eigene `.d-check.yml` aktiviert
  `ids` + `matrix` (Muster `ADR-*`/`MR-*`/`DC-*`; `MR-006`-
  Referenzrichtung maschinell kodiert); ids-Fortschreibung aus dem
  Selbstlauf: Headings und Definitions-Ort des Musters sind
  linkpflichtfrei; ~50 nackte Kennungen der eigenen Doku verlinkt
  bzw. als Code-Span fixiert.

- slice-008 — Modul `external` (`DC-FA-EXT-001`, opt-in): HTTP-Port im
  Hexagon + `httpcheck`-Adapter (HEAD mit GET-Fallback bei 405/501,
  Timeout konfigurierbar, Redirect-Limit 5, begrenzte Parallelität,
  eine Prüfung pro URL); Grund-Codes
  `external-status`/`external-timeout`/`external-redirects`;
  `make doc-check` läuft jetzt mit `--network none` und ist damit die
  automatisierte `DC-QA-03`-Messmethode (netzloser Lauf aller Module
  außer `external`); Interim-Mechanismus
  `isImplemented`/`SkippedModules` entfernt — alle fünf
  Vertragsmodule sind lauffähig.

- slice-010 — Image-Integrationstests + Reproduzierbarkeits-Belege:
  `make image-test` automatisiert die `DC-FA-DIST-001`-Akzeptanzkriterien
  gegen das lokal gebaute Image (Befund-Ausgabe und Exit-Code nativ
  vs. Container byte-identisch, read-only-Mount, fehlender Mount →
  Exit 2 mit Hinweis); `make versions` (Pins + Runtime-Image-ID),
  `make ci` (gates + image-test — Target der Release-Pipeline) und
  `make fullbuild` (volle Closure inkl. Benchmark, schließt mit dem
  Image-Hash); die „Nicht behauptet"-Listen sind leer.
- slice-009 — Gate-Endausbau (welle-03-Abschluss): `make coverage-gate`
  (Coverage über `./internal/...` per `-coverpkg`, Kalibrierungs-Bindung
  85 % → 90 % bei welle-03 done; Ist-Stand 92,9 %),
  `make gate-consistency` (Meta-Gate gegen Harness-Lügen: dokumentierte
  Targets ↔ Makefile in beide Richtungen, `DC-QA-03`-Modulliste der
  `.d-check.yml`, Selbsttest mit Phantom-Target bei jedem Lauf) — beide
  in `make gates` aggregiert; `make bench` mit `DC-QA-01`-Benchmark
  (Spez §`DC-QA-01.a` eingelöst: 1.000-Dateien-Fixture, Median aus 3
  Container-Läufen; gemessen: 551 ms ≪ 5 s).

### Changed

- Review R1 zu slice-009/010 (Gate-Infrastruktur): Benchmark-Median
  aus `RUNS` abgeleitet statt hart verdrahtet (latente
  Stilles-Grün-Falle); `--cpus 2` im Benchmark-Lauf + Spez-Präzisierung
  (2-vCPU-Normierung aus `DC-QA-01`); Meta-Gate-Parser erkennt
  Mehrfach-Target-Zeilen und schließt Variablen-Zuweisungen aus (mit
  Parser-Selbsttest); `fullbuild: ci bench` statt Kettenduplikat;
  `make versions` ohne Stage-FROM-Rauschen; drei dokumentierte
  Annahmen (QA-03-Listenformat fail-closed, amd64-Kopplung des
  image-test, Bench-Fixture bleibt zur Inspektion liegen);
  93-%-Kalibrierung als Nachtrag in der slice-009-Closure-Notiz.

- Review R1 zu slice-008: Fragmente werden vor Prüfung/Dedupe entfernt
  (eine Prüfung pro Ressource, Befund am Original-Linkziel);
  Schema-Vergleich case-insensitiv (kein stiller Gap zwischen `links`
  und `external` bei `HTTP://`); explizit gesetzte 0 in
  `external.timeout-seconds`/`parallel` ist Konfigurationsfehler statt
  stillem Default; GET-Fallback drained den Body begrenzt (64 KB);
  HTTP-Adapter wird nur noch bei aktivem `external` verdrahtet
  (strukturelle Opt-in-Absicherung); Timeout-Semantik (pro Request)
  spezifiziert; QA-03-Config-Kopplung als gate-consistency-Auftrag in
  slice-009 eingetragen.
- Review R1 zu slice-007: Status-Extraktion fence-aware (Fence-Inhalt
  ist kein Statuswert) und nur für Markdown-Ziele (kein Voll-Read von
  Binärdateien); Doppel-Befund forbidden+inactive als unabhängige
  Verletzungen spezifiziert; gemeinsamer Fence-Scanner (`proseLines`)
  und gemeinsame Ziel-Auflösung (`localTarget`) statt Drittkopien;
  `exclude-sections` der Selbstkonfiguration um die realen
  nummerierten Headings („7. Historie") ergänzt.
- Review R1 zu slice-006: Config-Pfade (`scan.roots`,
  `ids.patterns[].target`) dürfen die Repo-Wurzel nicht verlassen
  (Exit 2 statt stillem Escape); Leerstring-matchende ids-Regexe sind
  Konfigurationsfehler; Inline-Code-Stripping positionserhaltend
  (keine Phantom-Kennungen durch Text-Verschmelzung); zeilenbasierte
  Link-Extraktion als normative Grenze dokumentiert.
- Harness-Hooks gehärtet (`MR-005`): Gate-Nachweis inhaltsbasiert
  (Commit ohne Gate-Lauf wird vom Stop-Hook nicht mehr freigegeben),
  PreToolUse-Guard prüft `bash/sh -c`-Sub-Shell-Strings rekursiv.
- Referenzrichtungs-Korrektur (`MR-006`): ADR-Abwärtsverweise aus
  `spec/spezifikation.md` und `spec/architecture.md` entfernt
  (Kurs-Template-Fehler; Spec-Straten verweisen nie abwärts,
  Traceability über die `Schärft:`-Felder der ADRs).
- `spec/architecture.md` sprachneutral umformuliert (Schichten/Rollen
  statt Modul-Pfade und Imports; sprachkonkrete Übersetzung lebt in
  `ADR-0004`) — Template-Hard-Rule „sprach- und meilensteinfrei" wieder
  voll erfüllt.
- Lastenheft 0.2.1 (redaktionell): Beispiel-Kennungen auf fiktive
  Nummern (`ADR-0042`, `ADR-0099`) — keine Kollision mit real
  entstandenen eigenen ADRs.
