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

### Changed

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
