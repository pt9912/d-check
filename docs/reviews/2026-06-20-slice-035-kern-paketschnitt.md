# Review R1 — slice-035: Kern-Paketschnitt `model`/`rules`/`app`

- **Datum:** 2026-06-20
- **Reviewer:** Claude (Agent)
- **Modell:** `claude-opus-4-8[1m]`
- **Gegenstand:** slice-035, Commit-Range `d0d7899..HEAD`
  (`5e78237`, `bf567d4`, `72f38e8`, `4e7d26d`, `81d6bcd`, `3fdbcd4`,
  `daec49a`)
- **Rolle:** unabhängiger Code-Reviewer (kein Verifier — Gates **nicht**
  als grün angenommen; Plan-/DoD-Konformität ist nicht Gegenstand).

## Eingangs-Kontext

- ADR: [`docs/plan/adr/0012-kern-paketschnitt-model-rules-app.md`](../plan/adr/0012-kern-paketschnitt-model-rules-app.md)
  (Accepted) — schärft [`spec/architecture.md`](../../spec/architecture.md)
  §Kern, ergänzt [ADR-0005](../plan/adr/0005-modul-layout-hexagon-ordner.md)
  (arch-check). Bindung [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)
  (kein Befund-Delta), [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
  (I/O-Freiheit).
- Slice + Closure: [`docs/plan/planning/done/slice-035-kern-paketschnitt.md`](../plan/planning/done/slice-035-kern-paketschnitt.md).
- Hard Rules: [`AGENTS.md`](../../AGENTS.md) §3.
- Reviewer-Skill: [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) v1.0.0.

## Findings

### F1 — Über-Export ohne Grenz-Überschreiter (`MatchGlob`)

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** `internal/hexagon/core/rules/paths.go:96`
- **befund:** `MatchGlob` wurde in `72f38e8` exportiert mit der
  Commit-Begründung „Symbole, die künftig die rules→app- bzw.
  →model-Grenze kreuzen"; die einzigen Aufrufer sind jedoch
  `paths.go`, `scan.go`, `matrix.go` und `paths_test.go` — alle im Paket
  `rules` selbst, kein Konsument in `app`/`adapter`/`cmd`. Der Export
  überschreitet keine Paketgrenze; ein unexportiertes `matchGlob`
  hätte alle aktuellen Aufrufe bedient. Kein Verhaltens-Effekt.
- **verifizierbar:** ja — `make arch-check` bleibt grün (Export verletzt
  keine Importregel); ein paketweiter Symbol-Such-Lauf über
  `internal/hexagon/core/{app}` + `internal/adapter` + `cmd` nach
  `MatchGlob` liefert null Treffer.

### F2 — Über-Export ohne Grenz-Überschreiter (`Ignored`)

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** `internal/hexagon/core/rules/paths.go:123`
- **befund:** `Ignored` wurde exportiert (und der Helfer von `scan.go`
  nach `paths.go` verschoben), begründet damit, dass „ids ihn nutzt"
  und so `rules → app` vermieden wird (ADR-0012 Konsequenzen). Da `ids`
  nach dem Schnitt in `rules` liegt, sind alle Aufrufer
  (`ids.go:30`, `scan.go:94`, `scan.go:120`) paket-intern; ein
  cross-package-Konsument existiert nicht. Der Export ist breiter als
  nötig; kein Verhaltens-Effekt.
- **verifizierbar:** ja — Symbol-Such-Lauf nach `Ignored` außerhalb von
  `internal/hexagon/core/rules` liefert null Treffer; `make arch-check`
  bleibt grün.

### F3 — Stale Rationale-Kommentar in `.golangci.yml` (umbenanntes Primitiv)

- **kategorie:** LOW
- **quelle:** Maintainability
- **pfad:** `.golangci.yml:178`
- **befund:** Der Kommentar begründet die `testpackage`-Ausnahme mit
  „die Kern-Tests sind bewusst White-Box (unexportierte
  Analyse-Primitiven wie stripInlineCode/matchGlob …)". `matchGlob`
  wurde in diesem Slice (`72f38e8`) zu `MatchGlob` exportiert (ist also
  nicht mehr unexportiert), und `stripInlineCode` existiert im Kern
  nicht (mehr). Die zitierten Beispiele tragen die angegebene
  Begründung nicht mehr; `.golangci.yml` wurde im Slice-Diff nicht
  mitgezogen.
- **verifizierbar:** ja — `make lint` liefert keine Meldung (reine
  Kommentar-Drift, kein Regel-Effekt); Beleg ist der Symbol-Such-Lauf:
  `MatchGlob` ist exportiert (`paths.go:96`), `stripInlineCode` fehlt.

### F4 — architecture.md §Kern nennt Paketnamen und Flags

- **kategorie:** INFO
- **quelle:** Hard-Rule AGENTS.md §3.4
- **pfad:** `spec/architecture.md:46`
- **befund:** Die geschärfte §Kern-Zeile führt die konkreten Paketnamen
  `model`/`rules`/`app` sowie die Flags `--doctor`/`--repair`/
  `--suggest-config` in die Architektur-Spec ein; die Flag-Namen
  standen vor diesem Slice nicht in `architecture.md`. AGENTS.md §3.4
  verlangt „Schichten und Rollen statt Technologie — keine Sprach-/
  Modul-Pfade". ADR-0012 (`Schärft:` §Kern) sanktioniert die Bearbeitung
  ausdrücklich, und die drei Namen sind als Rollen (Daten/Engine/Modi)
  glossiert — daher Spannung statt klarer Verstoß; dokumentiert zur
  Bewertung, ob Paket-/Flag-Namen in Rang-3-Spec gehören.
- **verifizierbar:** ja — `make doc-check` bleibt grün (kein
  Anker-/Link-Bruch); der Befund ist Konventions-Konformität, kein
  Gate-Versagen.

## Negativbefunde (geprüft, ohne Befund)

- **R6 ist kein No-op (arch-check):** `tools/arch-check.sh` Zeilen 75–88 —
  die `case`-Arme matchen die realen Importpfade
  `github.com/pt9912/d-check/internal/hexagon/core/{rules,app}` (Modul =
  `github.com/pt9912/d-check`, `go.mod:1`), exakt wie `app` sie
  importiert (`internal/hexagon/core/app/diagnose.go` Zeilen 8/10). Ein
  `model→rules`- oder `rules→app`-Import setzt `fail=1`; die Logik
  feuert nachweislich.
- **Verhaltens-Delta DC-QA-02 (sed-Qualifizierung):** Die Diffs der
  hoch-geänderten Module (`ids.go`, `anchors.go`, `external.go`,
  `codepaths.go`, `matrix.go`, `report.go`) zeigen ausschließlich
  Bezeichner-Umbenennungen (klein→groß) und `model.`/`rules.`-
  Qualifizierungen; jede logiktragende Zeile (`findings = append`,
  `Reason: …`, Struct-Felder, Argument-Reihenfolge) bleibt identisch.
  Keine Über-/Unter-Qualifizierung gefunden.
- **`Line`-Feld-Kollision (zurückgenommen):** Der Typ `rules.Line`
  (`markdown.go:8`, Feld `No`) und das Feld `Finding.Line` (int,
  `model/finding.go:25`) koexistieren konfliktfrei — Struct-Literale
  setzen `Line: ln.No` (`rules/ids.go:143`) bzw. `Line: pl.no`. Keine
  Restspur einer `rules.Line`-Fehlqualifizierung im Code.
- **String-/Kommentar-Korruption:** Kein `model.`/`rules.`-Präfix in
  produktiven String-Literalen; einziger Treffer ist ein Test-
  Fehlermeldungs-String (`rules/ids_test.go:83`, „model.Config-
  Constraint"), der den Typ benennt, kein Verhalten ändert.
- **`model` ist Blatt (DC-QA-03):** `model/finding.go` importiert nur
  `sort`, `model/config.go` nur `fmt`/`regexp`/`sort`/`strings` — keine
  Kern-Pakete, keine I/O-APIs.
- **`rules ↛ app`:** Kein `rules`-Nichttest-File importiert
  `internal/hexagon/core/app`.
- **`coretest` nur Port:** `coretest/memfs.go:13` importiert allein
  `internal/hexagon/port/driven` (+ stdlib).
- **Test-only-Exporte (4 erwartete Interna unexportiert):**
  `htmlAnchors` (`rules/anchors.go:120`), `normalizeCodepath`
  (`rules/codepaths.go:62`), `classifyCodepath` (`rules/codepaths.go:94`),
  `checkURLs` (`rules/external.go:98`) sind alle weiterhin
  klein-geschrieben (unexportiert).
- **`Check*` sind nicht test-only:** Die acht exportierten
  `Check*`-Funktionen werden vom paket-internen Orchestrator
  `checkFile` aufgerufen (`rules/run.go:94-115`), nicht nur aus Tests —
  der Export trägt die in ADR-0012 begründete White-box-Test-Kopplung.
- **Vor-Refactor-Exporte nicht diesem Slice anzulasten:** `Slugify`,
  `HeadingSlugs`, `ExtractLinks` waren bereits bei `d0d7899` exportiert
  (vor slice-035) — keine neu eingeführten Test-Exporte.
- **Adapter-/CLI-Qualifizierung sauber:** `report.go` routet
  Reporter-/Config-Typen nach `model.*` und Diagnose-/Repair-Helfer
  (`ReasonText`, `FixCandidateFor`, `RepairEdit`) nach `app.*`,
  konsistent mit ADR-0012 (Adapter importieren `…/core/model` bzw.
  `…/core/app`); keine Logikänderung.
- **Deviation run/scan→rules (Begründung stichhaltig):** `run.go`/
  `scan.go` koppeln im White-box-Test `Run` mit Modul-Interna
  (`runState`, `checkFile`, `discoverScopes`); in `rules` bleiben sie
  ohne Interna-Export testbar. „Engine vs. Modi" ist die vom Code
  getragene Naht; verwässert die rules-Semantik nicht problematisch.
- **Doku-Kohärenz / MR-006:** `spec/architecture.md` verweist nirgends
  abwärts auf ADRs (MR-006 gewahrt); ADR-0012, §Kern und Closure-Notiz
  beschreiben übereinstimmend `run`/`scan` in `rules`; ADR-Index
  (`docs/plan/adr/README.md`) korrekt auf Accepted gebumpt.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 (F3) |
| INFO | 3 (F1, F2, F4) |

## Verdikt

**Merge-fähig.** Keine HIGH/MEDIUM-Befunde: Importrichtung maschinell
erzwungen (R6 feuert real), Paketgrenzen sauber, kein Spur eines
Verhaltens-Deltas aus der sed-Qualifizierung, die `Line`-Feld-Kollision
nachweislich zurückgenommen. Verbleibend ein LOW (stale `.golangci.yml`-
Kommentar) und drei INFO (zwei breiter-als-nötige Exporte; Paket-/
Flag-Namen in der Architektur-Spec) — keiner blockiert. Hinweis an die
Verifikation (getrennter Kontext): die Gate-Grün-Behauptung inkl.
`arch-check` R1–R6 und Coverage 93 % wurde hier **nicht** nachgefahren.
