# Changelog

Alle nennenswerten Änderungen an diesem Projekt werden in dieser Datei
dokumentiert. Das Format folgt [Keep a Changelog](https://keepachangelog.com/de/1.1.0/),
die Versionierung folgt [SemVer](https://semver.org/lang/de/).

## [Unreleased]

### Added

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
