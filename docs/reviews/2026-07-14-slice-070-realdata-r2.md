# Review R2 — slice-070 Realdaten-Schnitt

**Datum:** 2026-07-14
**Review-Art:** Realdaten-Review während der Implementierung
**Gegenstand:** Tabellenvertrag von
[`slice-070`](../plan/planning/in-progress/slice-070-trace-tabellenquellen-nullmengen-guard.md)
gegen `m-trace/spec/lastenheft.md`
**Reviewer:** Codex (Self-Review; kein unabhängiger Review)
**Skill:** `.harness/skills/reviewer.md` v1.2.0
**Modell:** GPT-5 Codex

## Eingangs-Kontext

- m-trace-Arbeitsstand vom 2026-07-14, nur lesend untersucht; Lauf gegen eine
  Kopie unter `/tmp/m-trace-slice070`
- lokales `d-check:latest` mit dem ersten Tabellenparser-Stand
- Konfiguration mit `id-column: Kennung`, `text-column: Anforderung` und
  `modality-column: Prioritaet`
- beobachteter Lauf: `236 Anforderung(en), 98 Waise(n).`

## Findings

### F-1 · HIGH

- **kategorie:** HIGH
- **quelle:** Realdaten-DoD von
  [`slice-070`](../plan/planning/in-progress/slice-070-trace-tabellenquellen-nullmengen-guard.md)
- **pfad:** `m-trace/spec/lastenheft.md:502,1898`
- **befund:** Der Slice behauptet eine einheitliche Textspalte `Anforderung`.
  Tatsächlich tragen 236 Requirement-Zeilen diesen Header, weitere 136
  Requirement-Zeilen stehen in Tabellen mit dem Header
  `Akzeptanzkriterium`. Ein einzelnes exaktes `text-column` kann daher nie alle
  371 eindeutigen IDs erfassen; der reale Lauf erkennt nur 236.
- **verifizierbar:** ja; Tabellenzeilen nach Headerklasse zählen und den
  `--trace`-Summenwert prüfen.

### F-2 · HIGH

- **kategorie:** HIGH
- **quelle:** Realdaten-DoD und Duplikatvertrag von
  [`slice-070`](../plan/planning/in-progress/slice-070-trace-tabellenquellen-nullmengen-guard.md)
- **pfad:** `m-trace/spec/lastenheft.md:2045,2060`
- **befund:** m-trace enthält 372 passende Tabellenzeilen, aber 371 eindeutige
  IDs, weil `RAK-51` als historische Kann-Aussage und spätere Muss-Hochstufung
  zweimal definiert ist. Sobald beide Text-Header gelesen werden, verhindert
  der zwingende Duplikatfehler den geforderten Realdaten-Happy-Path.
- **verifizierbar:** ja; erste Tabellenspalte zählen (`372` Zeilen, `371`
  eindeutige Werte, `RAK-51` zweimal).

## Negativbefunde

- Der Tabellenlexer selbst erklärt die Differenz nicht: Alle 372 ID-Zeilen
  folgen der festgelegten Pipe-Tabellengrammatik.
- ID-Regex und Modalitätsspalte sind für beide Headerklassen identisch; hierfür
  ist keine weitere Konfigurationsachse erforderlich.
- Das m-trace-Original wurde nicht verändert; alle Ausführungen nutzten eine
  read-only gemountete Kopie.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---:|
| HIGH | 2 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

## Verdikt

**NACHBESSERN.** Der Tabellenvertrag braucht explizite alternative Text-Header
und eine deterministische, opt-in Duplikatpolitik. Der sichere Default bleibt
ein Fehler; m-trace kann die historische Hochstufung gezielt mit `last`
auflösen. Erst danach ist der 371er-Realdatenbeleg aussagekräftig.
