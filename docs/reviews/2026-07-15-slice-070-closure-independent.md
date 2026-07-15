# Review — slice-070 unabhängiger Closure-Review

**Datum:** 2026-07-15
**Review-Art:** unabhängiger Closure-Review (kontext-getrennt; frisches
Kontextfenster, nicht der Implementierer-Kontext der Self-Reviews R1–R7)
**Gegenstand:**
[`slice-070`](../plan/planning/done/slice-070-trace-tabellenquellen-nullmengen-guard.md) —
Trace-Tabellenquellen und Nullmengen-Guard (Feat-Commit `5b6b284` + uncommittete
R4-/R6-Fixes im Arbeitsbaum)
**Reviewer:** Claude (personell/kontext-unabhängiger Lauf)
**Skill:** `.harness/skills/reviewer.md` v1.2.0
**Modell:** Opus 4.8 (1M context)

## Eingangs-Kontext

- Slice-Vertrag §2/§3/§4/§7
  [`slice-070`](../plan/planning/done/slice-070-trace-tabellenquellen-nullmengen-guard.md)
- Lastenheft
  [`DC-FA-REQ-001`](../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen),
  [`DC-FA-CLI-009`](../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix),
  [`DC-FA-MOD-001`](../../spec/lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in)
- Spezifikation
  [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
- [ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md) (Proposed)
- Hard Rules [`AGENTS.md`](../../AGENTS.md) §3, `harness/conventions.md`
- Vor-Findings (Kontext, nicht Autorität): R4/R5/R6/R7
- Prüfgegenstand: `internal/hexagon/core/app/trace_table.go`,
  `internal/hexagon/core/app/trace.go`,
  `internal/hexagon/core/model/config.go`,
  `internal/adapter/driven/configyaml/configyaml.go`,
  `internal/adapter/driving/cli/cli.go`,
  `internal/adapter/driving/cli/config_template.go`,
  `internal/adapter/driving/cli/cli_acceptance_test.go`

Gemäß Reviewer-Skill war die DoD-Abhakung nicht Grundlage dieses Reviews
(Verifier-Kontext, getrennt).

## Findings

### F-1 · LOW · offen

- **kategorie:** LOW
- **quelle:** Maintainability
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go:209`
- **befund:** Der uncommittete R6-Fix hat `text-column`/`text-columns` auf
  Zeiger-Typen (`*string` / `*[]string`) umgestellt, den Struct-Tag-Block aber
  nicht gofmt-konform ausgerichtet: die Tag-Spalte von
  `TextColumns *[]string` beginnt genau eine Spalte rechts der vier
  Geschwisterfelder (Backtick-Spalte 26 vs. 25). Der `rawTraceTable`-Block ist
  damit nicht gofmt-clean. Das aktive Linter-Profil (`.golangci.yml`,
  `default: none`) aktiviert weder `gofmt`/`gofumpt` noch `goimports`, daher
  bleiben `make lint` und `make test` grün und der Befund ist nicht gate-gefangen.
- **verifizierbar:** ja — `gofmt -l internal/adapter/driven/configyaml/configyaml.go`
  listet die Datei (gofmt-Umbruch), bzw. `sed -n '207,211p' … | cat -A` zeigt die
  ungleichen Tag-Spalten. Kein bestehender Gate-Lauf bestätigt ihn.

### F-2 · INFO · offen

- **kategorie:** INFO
- **quelle:**
  [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
  (Fehlerpräzedenz „… → Tabellenstruktur/Header → Duplicate-ID → …")
- **pfad:** `internal/hexagon/core/app/trace_table.go:43` (unbenutzter
  Text-Header-Guard nach dem Scan) vs. `:124` (Duplicate-ID-Fehler während des
  Scans)
- **befund:** Der „jede `text-columns`-Alternative muss mindestens einmal
  gebunden werden"-Guard (laut §4 eine Header-Phasen-Prüfung) läuft erst nach
  dem vollständigen Zeilen-Scan. Der Duplicate-ID-Fehler unter Politik `error`
  wird dagegen mitten im Scan geworfen. Treffen beide Fehlerbedingungen in
  derselben Quelle zusammen (eine deklarierte, aber nirgends gebundene
  Text-Alternative UND eine doppelte ID in einer gebundenen Tabelle), gewinnt die
  Duplicate-ID-Meldung — entgegen der in der Spezifikation gelisteten Reihenfolge
  „Header → Duplicate-ID". Exit-Code (2) und Determinismus bleiben unberührt; nur
  die Diagnose-Meldung weicht in diesem konstruierten Doppelfehler-Fall ab.
- **verifizierbar:** ja — eine Fixture mit `text-columns: [Text, Tippfehler]`,
  einer Tabelle, die nur `Text` bindet und darin dieselbe ID doppelt trägt
  (`duplicate-ids: error`), liefert die Meldung `doppelte Anforderungs-ID` statt
  `Tippfehler … kommt in keiner Tabelle … vor`. Beide Wege bleiben Exit 2.

## Negativbefunde (geprüft, ohne Befund)

- **Zwei explizite Formate & strikte Header-Bindung:** `format: headings`
  (Default) / `table`; `resolveTrace`/`loadTraceRequirements`
  (`trace.go:94–200`) und `bindTableColumns` (`trace_table.go:152–191`) binden
  Spalten ausschließlich über exakt, case-sensitiv verglichene Header-Namen; genau
  eine von `text-column`/`text-columns` ist Pflicht; keine Positions-/
  Synonym-Heuristik. Geprüft, ohne Befund.
- **Ganzzellen-ID & Modalitätsquelle:** `addTableRequirement`
  (`trace_table.go:119–141`) definiert eine Anforderung nur bei
  `isFullReqID(reqPat, id)` (`FindString==tok`); Modalität kommt exklusiv aus
  `modality-column`, sonst der gebundenen Textzelle. Boundary-Test bestätigt, dass
  der widersprechende Textinhalt die Klassifikation nicht beeinflusst. Geprüft,
  ohne Befund.
- **Duplikatpolitik & Nullmengen-Guard (fail-closed):** `error`(Default)/`first`/
  `last` deterministisch (Sortierung am Ende); expliziter Source ODER
  Tabellenmodus setzt `strictReqs` (`trace.go:102–115`); fehlende/leere Quelle
  und null erkannte Anforderungen ⇒ Exit 2 vor dem Reporter; `--require-complete`
  liefert bei `BuildTraceMatrix`-Fehler Exit 2 statt Exit 1 (`cli.go:181–194`);
  `source: ""` aktiviert keinen Guard. Geprüft, ohne Befund.
- **Fehlerpräzedenz & Tabellen-Lexik:** Config-Schema (`applyTrace…`) vor
  Dateilesen; Tabellenstruktur (fehlende/doppelte Header, Zeilenbreite) vor
  Nullmenge; `\|`, Rand-Pipes und einzeilige, korrekt geschlossene Backtick-Spans
  teilen keine Zelle; ungeschlossener Backtick bleibt literal (R4-F-1-Fix,
  `hasClosingBacktickRun`, `trace_table.go:330–346`, Test `unclosed code span`);
  Fenced-Code wird über `markdownTableLines` ausgeklammert. Außer F-2 (INFO)
  geprüft, ohne Befund.
- **Kompatibilität:** ohne `trace`-Block bleibt `tc` der Nullwert ⇒ Heading-Pfad
  mit First-Wins-Deduplizierung, keine Guards, leere Nullmenge Exit 0
  (`traceHeadingRequirements` strict=false, `trace.go:247–280`); der
  Strict-Zweig (`strict=true`) greift nur bei nichtleer explizitem Source.
  Geprüft, ohne Befund.
- **Config-Schema (R6-Fix-Inhalt):** Exklusivität `text-column`/`text-columns`
  jetzt präsenz- statt inhaltsbasiert (`*string`/`*[]string`); leerer
  Einzelwert, leere/duplizierte/leere-Element-Liste sowie unbekannte
  `duplicate-ids` sind fail-closed (`configyaml.go:450–517`); vier neue
  CLI-Negativfälle verriegeln die Grenzen. Außer F-1 (LOW, reine Formatierung)
  geprüft, ohne Befund.
- **Hard Rules / Architektur:** keine neue Dependency (nur stdlib + `core/model`
  + `port/driven`; `go.mod` unberührt), keine `//nolint`-Suppression, keine
  Gate-/Schwellen-Lockerung, kein Host-Go, Import-Richtung
  ([ADR-0005](../plan/adr/0005-modul-layout-hexagon-ordner.md)) eingehalten.
  Geprüft, ohne Befund.
- **Determinismus / Read-only / Netzlos (DC-QA-02/03):** ID-Reihenfolge via
  `sort.Strings`, Refs/Coverage sortiert; nur `fsys.ReadFile`-Lesepfade, keine
  Schreib-/Netz-Operation. Geprüft, ohne Befund.
- **Dogfooding-Rückwirkung:** d-checks eigene `.d-check.yml` trägt keinen
  `trace`-Block ⇒ Default-Heading-Pfad, `strictReqs=false` ⇒ der
  Completeness-/RTM-Dogfood bleibt unverändert. Geprüft, ohne Befund.
- **Sensor-Verifikation:** `make test` auf dem aktuellen Arbeitsbaum grün (alle
  Pakete `ok`, inkl. `cli`/`core/app`; Image `sha256:e30ee9eb…`). Geprüft, ohne
  Befund.

**Scope-Hinweis (Verifier-Domäne):** Die DoD-Vollständigkeit (Release,
Digest-Backfill, Realdatenbeleg, `make gates`/`ci`/`fullbuild`-Abhakung) ist
Verifier-Sache (Modul 11) und wurde nicht als Review-Kriterium bewertet; es fiel
dabei keine DoD-Lücke auf.

## Kategorie-Summary

| Kategorie | gefunden | blockierend |
| --------- | -------: | ----------: |
| HIGH      |        0 |           0 |
| MEDIUM    |        0 |           0 |
| LOW       |        1 |           0 |
| INFO      |        1 |           0 |

## Verdikt

**ACCEPT-WITH-NITS.** Der öffentliche Tabellen-/Nullmengen-Vertrag
(DC-FA-REQ-001 / DC-FA-REQ-001.a / ADR-0037) ist korrekt, fail-closed und
determinismus-/read-only-treu umgesetzt; die R4-Backtick- und R6-Config-Fixes
halten und sind durch Akzeptanztests verriegelt; die Default-Kompatibilität
bleibt gewahrt. Kein HIGH-/MEDIUM-Befund blockiert die Closure. Verbleiben zwei
nicht-blockierende Nits: F-1 (LOW, nicht gofmt-clean im R6-Fix, nicht
gate-gefangen) und F-2 (INFO, Präzedenz-Nuance zwischen unbenutztem
Text-Header-Guard und Duplicate-ID-Fehler, Exit-Code/Determinismus unberührt).
Beide können unabhängig vom Lifecycle-Move behoben oder als Won't-Fix notiert
werden.
