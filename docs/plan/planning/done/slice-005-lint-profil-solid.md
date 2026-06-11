# Slice slice-005: SOLID-nahes Lint-Profil (u-boot-Parität)

**Status:** done.

**Welle:** welle-03-regelmodule (Gate-Ausbau).

**Bezug:** [ADR-0001](../../adr/0001-implementierungssprache.md)
(Lint-Profil als Konsequenz), `ADR-0006` (entsteht in diesem Slice);
[`AGENTS.md`](../../../../AGENTS.md) §3.2 (Suppression-Verbot). Kein
direkter Lastenheft-Bezug — Qualitätsinfrastruktur (Kurs-Modul 13).

**Autor:** pt9912. **Datum:** 2026-06-10.

---

## 1. Ziel

Das golangci-lint-Profil erreicht u-boot-Parität (5 Default- + 24
SOLID-nahe Linter inkl. Kalibrierung und Why-kommentierten
Ausnahmen), beschlossen als `ADR-0006`; der Code ist ohne
Inline-Suppressions lint-clean.

## 2. Definition of Done

- [x] `ADR-0006` `Accepted` (min. 3 Alternativen; Abweichung vom
  Vorbild begründet: kein `depguard` — Architektur-Regeln bleiben in
  der `ADR-0005`-Fitness-Function `tools/arch-check.sh`).
- [x] `.golangci.yml` trägt das Profil; alle Ausnahmen zentral mit
  `Why:`-Kommentar (keine `//nolint`-Direktiven im Code).
- [x] Code lint-clean; Globals außerhalb `cmd/` eliminiert statt
  ausgenommen (gochecknoglobals ohne Breitband-Carveout).
- [x] Sensors-Tabelle ([`harness/README.md`](../../../../harness/README.md))
  und [`CHANGELOG.md`](../../../../CHANGELOG.md) aktualisiert.
- [x] `make gates` grün.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/0006-lint-profil-solid.md` | neu | Profil-Entscheidung (Vorbild: u-boot-`ADR-0003`) |
| [`.golangci.yml`](../../../../.golangci.yml) | update | Profil + Kalibrierung + Ausnahmen |
| `internal/hexagon/core` (Globals → Funktionen) | refactor | gochecknoglobals ehrlich erfüllen |
| `internal/adapter/driven/configyaml` | refactor | Regex lokal; Test als Black-Box-Paket |
| [`harness/README.md`](../../../../harness/README.md), [`docs/plan/adr/README.md`](../../adr/README.md) | update | Sensors-Vertrag, ADR-Index |

## 4. Trigger

Wellen-Start welle-03 (Gate-Ausbau ist Teil des Wellen-Umfangs).

## 5. Closure-Trigger

DoD vollständig + Commit(s) auf `main` + Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Die strengen Komplexitäts-Linter (funlen/gocognit/maintidx) können
  Refactorings in `cli.Run` erzwingen — gewollt, aber Verhalten muss
  durch die Akzeptanztests gedeckt bleiben.
- `testpackage` kollidiert mit den White-Box-Kern-Tests — Ausnahme
  nur für `internal/hexagon/core/` mit Why (Black-Box-Abdeckung
  liegt in den CLI-Akzeptanztests).

## 7. Closure-Notiz (nach `done/`)

**Umsetzung:** Commit `33e2cd5` (`ADR-0006`, Profil, Refactorings,
Doku-Nachzug).

- **Was hat funktioniert:** Die u-boot-Kalibrierung passte
  unverändert; das Profil fand sofort 12 echte Verstöße (6×
  Komplexität, Builtin-Shadowing, fehlende Export-Kommentare) — alle
  durch Refactoring statt Ausnahmen gelöst (Globals → Funktionen,
  Helfer-Extraktion). Endstand: 0 Issues ohne ein einziges `//nolint`.
- **Anders als geplant:** nichts Wesentliches; der Plan hatte die
  cli/configyaml-Splits bereits als Risiko antizipiert.
- **Steering-Loop-Lerneintrag:** Die `testpackage`-Ausnahme für die
  White-Box-Kern-Tests ist die einzige Pfad-Ausnahme im Profil —
  beim Schneiden der ids/matrix/external-Slices darauf achten, neue
  Kern-Tests möglichst black-box zu formulieren, damit die Ausnahme
  nicht wächst.
- **Folge-Slices:** keine neuen; die restlichen welle-03-Slices
  (ids, matrix, external, coverage-gate, gate-consistency) stehen
  noch zum Schneiden aus.

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Konfigurations-/Refactor-Arbeit am
spec-geführten Code; siehe Kurs Modul 5 §Worked Mini-Example).
