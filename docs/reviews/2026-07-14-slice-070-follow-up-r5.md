# Review R5 — slice-070 R4-Fixverifikation

**Datum:** 2026-07-14  
**Review-Art:** technischer Folgereview der R4-Befunde  
**Gegenstand:** Arbeitsbaum-Diff nach
[`Review R4`](2026-07-14-slice-070-follow-up-r4.md) zu
[`slice-070`](../plan/planning/done/slice-070-trace-tabellenquellen-nullmengen-guard.md)  
**Reviewer:** Codex (Self-Folgereview; kein personell unabhängiger Review)  
**Skill:** `.harness/skills/reviewer.md` v1.2.0  
**Modell:** GPT-5 Codex

## Findings

### R4-F-1 · MEDIUM · behoben

- **kategorie:** MEDIUM
- **quelle:**
  [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen),
  Tabellen-Lexik und Datenzeilen
- **pfad:** `internal/hexagon/core/app/trace_table.go:309`
- **befund:** Der Splitter öffnet einen Backtick-Code-Span nur noch, wenn auf
  derselben Zeile eine exakt gleich lange schließende Folge existiert. Ein
  ungeschlossener Backtick bleibt literal, sodass die folgende Pipe als
  Zelltrenner wirkt und eine falsche Zeilenbreite Exit 2 auslöst.
- **verifizierbar:** ja — der neue Akzeptanzfall `unclosed code span` in
  `TestCLI070_TraceTable_Negative` ist mit `make test` grün.

### R4-F-2 · LOW · behoben

- **kategorie:** LOW
- **quelle:** Planning-/Lifecycle-Wahrheit
- **pfad:**
  `docs/plan/planning/in-progress/slice-070-trace-tabellenquellen-nullmengen-guard.md:15`
- **befund:** Der Kopf dokumentiert v0.43.0 nun als veröffentlicht; der
  Lifecycle-Abschnitt beschreibt den tatsächlichen Zustand nach R4 und hält den
  Slice bis zum personell unabhängigen Closure-Review in `in-progress/`.
- **verifizierbar:** ja — Kopf, Lifecycle-Abschnitt, DoD und Closure-Notiz sind
  widerspruchsfrei; `make planning-check` bleibt grün.

## Negativbefunde

- **Code-Span-Grenzen:** passend begrenzte Einzel- und Mehrfach-Backtick-Spans
  schützen ihre inneren Pipes weiterhin; ohne weiteren Befund.
- **Pipe-/Breitensemantik:** escaped Pipes, optionale Rand-Pipes und bestehende
  Zeilenbreitenfehler bleiben unverändert; ohne Befund.
- **Trace-Vertrag:** Headerbindung, Ganzzellen-ID, Duplikatpolitik, Modalität und
  Nullmengen-Guard werden vom Fix nicht verändert; ohne Befund.
- **Planning-Wahrheit:** Release-, Status- und Closure-Aussagen wurden
  gegeneinander geprüft; ohne weiteren Befund.
- **Hard Rules:** keine neue Dependency, keine Suppression, keine
  Gate-Lockerung und kein Host-Go; ohne Befund.

## Kategorie-Summary

| Kategorie | übernommen | offen |
|---|---:|---:|
| HIGH | 0 | 0 |
| MEDIUM | 1 | 0 |
| LOW | 1 | 0 |
| INFO | 0 | 0 |

## Sensoren

- `make test` — grün
- `make doc-check` — grün, 0 Befunde

## Verdikt

**ACCEPT.** Beide R4-Befunde sind im Arbeitsbaum behoben und durch Test bzw.
Planning-/Doku-Sensor verifizierbar. Der für die Slice-Closure geforderte
personell unabhängige Review bleibt davon ausdrücklich unberührt.
