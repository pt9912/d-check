# Review R4 — slice-070 Follow-up

**Datum:** 2026-07-14  
**Review-Art:** technisches Implementierungs- und Plan-Folgereview  
**Gegenstand:** Commit-Range `4fc4d81..f047bc3` zu
[`slice-070`](../plan/planning/done/slice-070-trace-tabellenquellen-nullmengen-guard.md)  
**Reviewer:** Codex (kein personell unabhängiger Review)  
**Skill:** `.harness/skills/reviewer.md` v1.2.0  
**Modell:** GPT-5 Codex

## Eingangs-Kontext

- Lastenheft-Anforderung
  [`DC-FA-REQ-001`](../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen)
  und technischer Vertrag
  [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
- Proposed
  [ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md)
- Tabellenextraktor, Trace-Auflösung, Config-Adapter und CLI-Akzeptanztests
- Nutzervertrag in Handbuch, Operations und Changelog
- Hard Rules aus [`AGENTS.md`](../../AGENTS.md) §3

Die DoD-Abhakung des Slice war gemäß Reviewer-Skill nicht Grundlage dieses
Reviews.

## Findings

### F-1 · MEDIUM · offen

- **kategorie:** MEDIUM
- **quelle:**
  [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen),
  Tabellen-Lexik und Datenzeilen
- **pfad:** `internal/hexagon/core/app/trace_table.go:291`
- **befund:** Jeder Backtick-Run öffnet im Splitter unmittelbar einen Code-Span,
  auch wenn auf derselben Zeile kein gleich langer Abschluss folgt. Dadurch kann
  eine Pipe hinter einem ungeschlossenen, laut Vertrag literalen Backtick ihre
  Funktion als Zelltrenner verlieren und eine relevante Zeile mit falscher
  Zellenzahl akzeptiert werden, statt Exit 2 auszulösen.
- **verifizierbar:** ja — eine Tabellen-Fixture mit Header `ID | Text` und
  Datenzeile `| R-1 | \`abc | def` muss wegen drei statt zwei Zellen Exit 2
  liefern; der aktuelle Parser übernimmt sie als zwei Zellen.

### F-2 · LOW · offen

- **kategorie:** LOW
- **quelle:** Planning-/Lifecycle-Wahrheit
- **pfad:**
  `docs/plan/planning/in-progress/slice-070-trace-tabellenquellen-nullmengen-guard.md:15`
- **befund:** Der Slice bezeichnet Release v0.43.0 noch als geplant und erklärt
  später, er könne nun nach `next/` und `in-progress/` aufgenommen werden. Das
  widerspricht dem eigenen Status `in-progress` sowie den Abschnitten, die
  Release, Digest-Backfill und den verbleibenden Closure-Schritt als
  abgeschlossen beziehungsweise offen dokumentieren.
- **verifizierbar:** ja — Zeilen 15–16 und 115–116 gegen Statuszeile 3 sowie
  Zeilen 95–98 und 153–158 vergleichen; `make planning-check` erkennt diese
  semantische Prosa-Drift nicht.

## Negativbefunde

- **Config und Strict-Grenze:** Format-/Block-Konsistenz, Pflichtspalten,
  alternative Text-Header, `source: ""` und der Nullmengen-Guard wurden geprüft;
  außer F-1 kein Befund.
- **Tabellenaggregation:** Headerbindung, Ganzzellen-ID, mehrere Tabellen und die
  Duplikatpolitiken `error`/`first`/`last` wurden geprüft; ohne weiteren Befund.
- **Modalität und RTM-Modell:** exklusive Modalitätsspalte, Text-Fallback sowie
  unveränderte ADR-/Slice-/Coverage-/Waisensemantik wurden geprüft; ohne Befund.
- **Kompatibilität:** Der unkonfigurierte Heading-Pfad, leere Defaultquelle und
  First-Wins-Deduplizierung wurden geprüft; ohne Befund.
- **Dokumentation:** Handbuch-Vertrag, Brownfield-Anleitung, Waisen-Semantik und
  Referenzscan-Details wurden geprüft; außer F-2 ohne Befund.
- **Architektur und Hard Rules:** gemeinsames RTM-Modell, keine neue Dependency,
  keine Inline-Suppression, keine Gate-Lockerung und kein Host-Go; ohne Befund.

## Kategorie-Summary

| Kategorie | gefunden | offen |
|---|---:|---:|
| HIGH | 0 | 0 |
| MEDIUM | 1 | 1 |
| LOW | 1 | 1 |
| INFO | 0 | 0 |

## Sensor

- `make test` — grün; der F-1-Grenzfall ist im bestehenden Testsatz nicht
  enthalten.

## Verdikt

**REQUEST CHANGES.** F-1 ist eine offene Spec-Treue-Lücke im neuen öffentlichen
Tabellenvertrag und blockiert den Abschluss. F-2 ist nicht blockierend, soll
aber vor dem Lifecycle-Move die Planning-Wahrheit wiederherstellen. Dieser Lauf
ist technisch eigenständig, erfüllt jedoch nicht die im Slice verlangte
personelle Unabhängigkeit.
