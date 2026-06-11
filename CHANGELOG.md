# Changelog

Alle nennenswerten Änderungen an diesem Projekt werden in dieser Datei
dokumentiert. Das Format folgt [Keep a Changelog](https://keepachangelog.com/de/1.1.0/),
die Versionierung folgt [SemVer](https://semver.org/lang/de/).

## [Unreleased]

### Added

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
  Text-/JSON-Reporter; Layout nach ADR-0005 (hexagon-/adapter-Ordner,
  u-boot-Konvention); Dockerfile-Stages + Make-Gates `lint`, `test`,
  `arch-check` (Fitness Function R1–R5), Runtime-Image
  distroless/static mit Selbst-Smoke-Test (`make run`).

- slice-004 — Modul `anchors` (GitHub-Slug-Verfahren inkl.
  Duplikat-Suffixen, Fragment-Dekodierung, Schweigen bei fehlender
  Zieldatei) und **Dogfooding**: `make doc-check` läuft über `d-check`
  selbst (`scan.roots: ["."]`, Module links+anchors — erstmals mit
  Anker-Validierung); vendorter Bootstrap-Sensor gelöscht
  (MR-003 → MR-007 aufgelöst); Vergleichslauf als erster
  DC-QA-04-Datenpunkt.

- slice-005 — SOLID-nahes Lint-Profil (ADR-0006, u-boot-Parität ohne
  depguard): 5 Default- + 23 Linter mit u-boot-Kalibrierung,
  gomodguard-Anti-Module, Why-kommentierte Ausnahmen; Code-Refactoring
  statt Carveouts (Globals → Funktionen, Komplexitäts-Splits in
  cli/configyaml/core) — lint-clean ohne //nolint.

- slice-006 — Modul `ids` (DC-FA-ID-001): Linkpflicht für Kennungen
  nach konfigurierten Regex-Mustern (Reihenfolge = Präzedenz, erstes
  Match gewinnt pro Vorkommen); „verlinkt" = Vorkommen im Linktext
  eines Markdown-Links, Ziel-Klammern und Bildreferenzen sind
  linkpflichtfrei (kein Fließtext); Grund-Code `id-unlinked`;
  Config-Constraint `ids.patterns[].target` muss existieren (Exit 2).

### Changed

- Review R1 zu slice-006: Config-Pfade (`scan.roots`,
  `ids.patterns[].target`) dürfen die Repo-Wurzel nicht verlassen
  (Exit 2 statt stillem Escape); Leerstring-matchende ids-Regexe sind
  Konfigurationsfehler; Inline-Code-Stripping positionserhaltend
  (keine Phantom-Kennungen durch Text-Verschmelzung); zeilenbasierte
  Link-Extraktion als normative Grenze dokumentiert.
- Harness-Hooks gehärtet (MR-005): Gate-Nachweis inhaltsbasiert
  (Commit ohne Gate-Lauf wird vom Stop-Hook nicht mehr freigegeben),
  PreToolUse-Guard prüft `bash/sh -c`-Sub-Shell-Strings rekursiv.
- Referenzrichtungs-Korrektur (MR-006): ADR-Abwärtsverweise aus
  `spec/spezifikation.md` und `spec/architecture.md` entfernt
  (Kurs-Template-Fehler; Spec-Straten verweisen nie abwärts,
  Traceability über die `Schärft:`-Felder der ADRs).
- `spec/architecture.md` sprachneutral umformuliert (Schichten/Rollen
  statt Modul-Pfade und Imports; sprachkonkrete Übersetzung lebt in
  ADR-0004) — Template-Hard-Rule „sprach- und meilensteinfrei" wieder
  voll erfüllt.
- Lastenheft 0.2.1 (redaktionell): Beispiel-Kennungen auf fiktive
  Nummern (`ADR-0042`, `ADR-0099`) — keine Kollision mit real
  entstandenen eigenen ADRs.
