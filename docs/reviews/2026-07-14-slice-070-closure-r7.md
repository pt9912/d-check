# Review R7 — slice-070 R6-Fixverifikation

**Datum:** 2026-07-14  
**Review-Art:** technischer Folgereview des R6-Befunds  
**Gegenstand:** Arbeitsbaum-Diff nach
[`Review R6`](2026-07-14-slice-070-closure-r6.md) zu
[`slice-070`](../plan/planning/in-progress/slice-070-trace-tabellenquellen-nullmengen-guard.md)  
**Reviewer:** Codex (Self-Folgereview; kein personell unabhängiger Review)  
**Skill:** `.harness/skills/reviewer.md` v1.2.0  
**Modell:** GPT-5 Codex

## Findings

### R6-F-1 · MEDIUM · behoben

- **kategorie:** MEDIUM
- **quelle:**
  [`DC-FA-REQ-001`](../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen),
  [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go:196`
- **befund:** Das Raw-Schema bewahrt die Präsenz von `text-column` und
  `text-columns` nun getrennt von ihrem Inhalt. Sobald beide Formen gesetzt
  sind, folgt Exit 2 — auch bei `text-column: ""` oder `text-columns: []`;
  die jeweils allein gesetzte Form wird weiterhin auf nichtleere Header
  validiert.
- **verifizierbar:** ja — die neuen CLI-Negativfälle `empty single plus list`
  und `single plus empty list` sind mit `make test` grün.

## Negativbefunde

- **Einzelform:** Ein allein gesetzter nichtleerer `text-column`-Wert wird
  weiterhin übernommen; ohne Befund.
- **Listenform:** Eine allein gesetzte, nichtleere `text-columns`-Liste bleibt
  zulässig; leere, doppelte oder unbenutzte Header bleiben fail-closed; ohne
  Befund.
- **Übriger Tabellenvertrag:** Headerbindung, Duplikatpolitik, Modalität,
  Nullmengen-Guard und Extraktion sind vom Config-Fix unberührt; ohne Befund.
- **Hard Rules:** keine neue Dependency, keine Suppression, keine
  Gate-Lockerung und kein Host-Go; ohne Befund.

## Kategorie-Summary

| Kategorie | übernommen | offen |
|---|---:|---:|
| HIGH | 0 | 0 |
| MEDIUM | 1 | 0 |
| LOW | 0 | 0 |
| INFO | 0 | 0 |

## Sensoren

- `make test` — grün
- `make lint` — grün, 0 Issues

## Verdikt

**ACCEPT.** Der R6-Befund ist im Arbeitsbaum behoben und durch zwei
CLI-Negativfälle verriegelt. Der für die Slice-Closure geforderte personell
unabhängige Review bleibt davon ausdrücklich unberührt.
