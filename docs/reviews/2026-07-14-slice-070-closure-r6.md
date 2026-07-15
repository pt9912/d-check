# Review R6 — slice-070 Closure

**Datum:** 2026-07-14  
**Review-Art:** technischer Code-/Closure-Review  
**Gegenstand:** Commit-Range `4fc4d81..f047bc3` einschließlich des aktuellen
Arbeitsbaum-Diffs zu
[`slice-070`](../plan/planning/in-progress/slice-070-trace-tabellenquellen-nullmengen-guard.md)  
**Reviewer:** Codex (technisch eigenständiger Lauf; personelle Unabhängigkeit
nicht belegbar)  
**Skill:** `.harness/skills/reviewer.md` v1.2.0  
**Modell:** GPT-5 Codex

## Eingangs-Kontext

- Lastenheft-Anforderungen
  [`DC-FA-REQ-001`](../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen),
  [`DC-FA-CLI-009`](../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
  und
  [`DC-FA-MOD-001`](../../spec/lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in)
- Technischer Vertrag
  [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
  und Proposed
  [ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md)
- Config-Adapter, Tabellenextraktor, Trace-Auflösung, CLI-Akzeptanztests und
  Nutzervertrag
- Hard Rules aus [`AGENTS.md`](../../AGENTS.md) §3
- R4-/R5-Folgeläufe als Historie der bereits bearbeiteten Befunde; die
  DoD-Abhakung des Slice war gemäß Reviewer-Skill keine Review-Grundlage

## Findings

### F-1 · MEDIUM · offen

- **kategorie:** MEDIUM
- **quelle:**
  [`DC-FA-REQ-001`](../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen),
  [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go:484`
- **befund:** Die Exklusivität von `text-column` und `text-columns` wird erst
  nach Trimmen beziehungsweise über `len(text-columns) > 0` entschieden. Damit
  werden `text-column: ""` plus gültiger `text-columns`-Liste sowie ein gültiger
  Einzelwert plus explizitem `text-columns: []` akzeptiert, obwohl der
  öffentliche Schema-Vertrag genau eine Form und nichtleere Header verlangt und
  dafür Exit 2 zusagt.
- **verifizierbar:** ja — je eine CLI-Negativprobe für beide Konfigurationen
  zeigt aktuell Exit 0 statt Exit 2; der zugehörige Akzeptanztest würde über
  `make test` laufen.

## Negativbefunde

- **R4-Nachbesserungen:** Ungeschlossene Backtick-Folgen bleiben literal und
  nachfolgende Pipes Zelltrenner; Release-/Lifecycle-Prosa ist wieder
  widerspruchsfrei. Außer F-1 ohne Befund.
- **Config und Strict-Grenze:** Format-/Block-Konsistenz, ID-/Modalitätsspalten,
  alternative Text-Header, Duplikatpolitik, nichtleere explizite Quelle,
  `source: ""` und Nullmengen-Guard geprüft; außer F-1 ohne Befund.
- **Tabellenextraktion:** Headerbindung, Ganzzellen-ID, Zeilenbreite, Escapes,
  geschlossene Code-Spans, Fences, mehrere Tabellen und `error`/`first`/`last`
  geprüft; ohne weiteren Befund.
- **Modalität und gemeinsames RTM-Modell:** exklusive Modalitätsspalte,
  Text-Fallback sowie unveränderte Referenz-/Coverage-/Waisen-/Reporter-Semantik
  geprüft; ohne Befund.
- **Default-Kompatibilität:** unkonfigurierter Heading-Pfad, leere Defaultquelle
  und First-Wins-Deduplizierung geprüft; ohne Befund.
- **Nutzerdokumentation und Referenzrichtung:** Handbuch, Changelog,
  Release-/Versionseinträge sowie Provenance-Marker-Ehrlichkeit geprüft; ohne
  Befund.
- **Architektur und Hard Rules:** keine neue Dependency, keine
  Inline-Suppression, keine Gate-Lockerung, kein Host-Go und kein zusätzlicher
  Netz-/Schreibpfad; ohne Befund.

## Kategorie-Summary

| Kategorie | gefunden | offen |
| --------- | -------: | ----: |
| HIGH      |        0 |     0 |
| MEDIUM    |        1 |     1 |
| LOW       |        0 |     0 |
| INFO      |        0 |     0 |

## Sensoren

- `make test` — grün
- `make gates` — grün; Coverage 93,70 %, Architektur, Semgrep, Doku,
  Planning und Gate-Konsistenz ohne Befund

## Verdikt

**REQUEST CHANGES.** F-1 ist eine offene Spec-Treue-Lücke im neuen
öffentlichen Config-Vertrag und blockiert die technische Closure.
