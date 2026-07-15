# Review R3 — slice-070 Implementierung

**Datum:** 2026-07-14
**Review-Art:** Implementierungsreview nach Realdatenkorrektur
**Gegenstand:** uncommitteter Arbeitsbaum zu
[`slice-070`](../plan/planning/done/slice-070-trace-tabellenquellen-nullmengen-guard.md)
**Reviewer:** Codex (Self-Review; kein unabhängiger Review)
**Skill:** `.harness/skills/reviewer.md` v1.2.0
**Modell:** GPT-5 Codex

## Eingangs-Kontext

- Lastenheft 0.43.0, technische Spezifikation und Proposed
  [ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md)
- Config-/Modelländerungen und Tabellenextraktor unter
  `internal/adapter/driven/configyaml/` bzw. `internal/hexagon/core/app/`
- CLI-Akzeptanztests einschließlich synthetischem 372-Zeilen-/371-ID-Fall
- Realdatenbeleg aus
  [`slice-070` §7](../plan/planning/done/slice-070-trace-tabellenquellen-nullmengen-guard.md#7-realdatenbeleg)
- Nutzervertrag in Handbuch, Operations und Changelog

## Findings

### F-1 · HIGH · behoben

- **kategorie:** HIGH
- **quelle:**
  [`DC-FA-REQ-001`](../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen)
- **pfad:** `internal/hexagon/core/app/trace_table.go`
- **befund:** Bei `text-columns: [Anforderung, Akzeptanzkriteriuum]` akzeptierte
  der erste Stand die 236 Tabellen mit dem gültigen ersten Header und ignorierte
  den vertippten zweiten Header. Der Nullmengen-Guard griff nicht; eine
  unvollständige RTM wäre mit Exit 0 möglich gewesen.
- **fix/verifikation:** Der Extraktor protokolliert jeden tatsächlich gebundenen
  Text-Header und verlangt nach dem Scan mindestens einen Treffer für jeden
  konfigurierten Alias. Akzeptanztest `unused text alternative`; m-trace-
  Gegenprobe endet mit Exit 2 und nennt `Akzeptanzkriteriuum`.

### F-2 · MEDIUM · behoben

- **kategorie:** MEDIUM
- **quelle:**
  [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
- **pfad:** `internal/hexagon/core/app/trace_table.go`
- **befund:** Die erste Fence-Erkennung toggelte bei jedem Präfix aus drei
  Backticks oder Tilden. Ein `~~~` innerhalb eines mit Backticks geöffneten
  Fences konnte den Bereich daher fälschlich verlassen und Beispieltabellen als
  Anforderungen lesen.
- **fix/verifikation:** Öffnungszeichen und Mindestlänge werden gespeichert;
  nur dasselbe Zeichen mit mindestens derselben Länge schließt. Der
  Fence-Akzeptanztest enthält nun absichtlich `~~~` innerhalb eines
  Backtick-Fences.

### F-3 · LOW · behoben

- **kategorie:** LOW
- **quelle:** Benutzerhandbuch-Currency
- **pfad:** `docs/user/benutzerhandbuch.md`
- **befund:** Die bereits vorhandene Historienzeile 1.27 war noch nicht in den
  Header-Stempel übernommen; Stand und Handbuch-Version blieben auf 1.26 bzw.
  2026-07-04.
- **fix/verifikation:** Header auf Handbuch 1.27 und Stand 2026-07-14 korrigiert.

## Negativbefunde

- Default-Heading-Pfad geprüft: `source: ""`, Nullmenge und First-Wins-
  Deduplizierung bleiben kompatibel; die neuen Strict-Fehler greifen nur bei
  nichtleer expliziter Quelle oder Tabellenmodus.
- Modalität geprüft: Eine konfigurierte Modalitätsspalte ist exklusiv; ohne sie
  wird genau die pro Tabelle gebundene Textspalte klassifiziert.
- Tabellenlexik geprüft: Ganzzellen-ID, relevante Zeilenbreite, escaped Pipe,
  Code-Span-Pipe, alternative Header, mehrere Tabellen und Duplikatpolitiken
  `error`/`first`/`last` sind beobachtbar getestet.
- Architektur/Hard Rules geprüft: gemeinsames RTM-Modell, kein zweiter
  Reporter, keine neue Dependency, keine Suppression, keine Gate-Lockerung und
  kein Host-Go.
- Doku-/Referenzrichtung geprüft: `make doc-check` meldet 0 Befunde.
- Lifecycle geprüft: `slice-070` bleibt `open`, solange `slice-069` die einzige
  aktive Welle ist; kein vorzeitiger Move oder Release-Register-Bump.

## Kategorie-Summary

| Kategorie | gefunden | offen |
|---|---:|---:|
| HIGH | 1 | 0 |
| MEDIUM | 1 | 0 |
| LOW | 1 | 0 |
| INFO | 0 | 0 |

## Sensoren

- `make test` — grün
- `make lint` — grün, 0 Issues
- `make doc-check` — grün, 0 Befunde
- `make gates` — grün; Coverage 93,70 %, Architektur/Semgrep/Planning/
  Gate-Konsistenz ohne Befund
- `make ci` — grün; Gates plus Image-Test
- m-trace — 371 Anforderungen; falscher Alias und unpassende Regex jeweils
  Exit 2

## Verdikt

**ACCEPT.** Alle im Self-Review gefundenen Befunde sind behoben und durch
Akzeptanz- oder Realdatenproben verriegelt. Die organisatorische DoD eines
**unabhängigen** Reviews bleibt davon ausdrücklich unberührt und offen.
