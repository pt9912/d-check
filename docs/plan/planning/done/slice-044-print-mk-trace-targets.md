# Slice slice-044: `doc-trace`/`doc-complete`-Fragment-Targets + opt-in `--require-complete`

**Status:** done (abgeschlossen, welle-33-print-mk-trace).

**Welle:** welle-33-print-mk-trace (Trigger: a-check-Bootstrap — Konsumenten
wollen die RTM/Vollständigkeits-Invariante als Makefile-Gate binden, ohne die
`completeness-check.sh`-Parsing-Logik zu kopieren).

**Bezug:** Neue Anforderung
[`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
(opt-in `--trace --require-complete`) + Change Request an
[`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
(Fragment um `doc-trace`/`doc-complete` erweitert). Baut auf
[`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
(RTM, slice-036) auf. Verwandt mit
[ADR-0017](../../adr/0017-requirements-completeness-gate.md) (Repo-eigenes
`make completeness-check`): der opt-in Strict-Exit ist die **Produkt-seitige**
Schwester der dortigen Wrapper-Durchsetzung — der **Default**-`--trace` bleibt
advisory, die ADR-Entscheidung unberührt. **Kein ADR** (additiv; neue CLI-Fähigkeit +
Fragment-Erweiterung, Präzedenz slice-036/slice-038 ohne ADR).

**Autor:** pt9912. **Datum:** 2026-06-23.

---

## 1. Ziel

Konsumenten von `d-check` (Anlass: a-check-Bootstrap) sollen die
RTM-Vollständigkeit als **Makefile-Gate** einbinden können, ohne ein Skript zu
kopieren. Zwei zusammengehörige Bausteine:

1. **`--trace --require-complete`** (neu,
   [`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)):
   derselbe read-only RTM-Lauf, aber `orphans > 0` ⇒ **Exit 1** statt 0. Default
   `--trace` bleibt advisory (Exit 0) — die Durchsetzung ist strikt opt-in.
2. **`--print-mk`-Fragment** (Erweiterung,
   [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)):
   zusätzlich die Targets `doc-trace` (advisory RTM) und `doc-complete`
   (`--trace --require-complete`, das Gate) plus eine überschreibbare
   `TRACE_FLAGS`-Variable.

## 2. Entscheidungen

- **Warum ein opt-in Flag statt Recipe-Logik:** Die Gate-Form darf weder den
  Default-`--trace` auf Exit ≠ 0 umbiegen (verboten durch
  [`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
  + [ADR-0017](../../adr/0017-requirements-completeness-gate.md), dort als
  Alternative bewusst verworfen — bräche jeden GF-Zwischenstand mit transienter
  Waise) noch die `completeness-check.sh`-Parsing-Logik ins Konsumenten-Recipe
  kopieren (verboten durch
  [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben):
  „keine Recipe-/Skript-Kopie"). Einziger sauberer Weg: die Fail-Logik lebt **in
  d-check** als opt-in Modifikator. Der Default bleibt advisory → die
  ADR-Invariante gewahrt.
- **Code (`internal/adapter/driving/cli/cli.go`):** `options.requireComplete`;
  Flag `--require-complete`; `comboError` weist `--require-complete` ohne
  `--trace` als Nutzungsfehler ab (Exit 2); `runTrace` liefert nach dem
  unveränderten Rendern bei `matrix.Orphans > 0` Exit 1 (Zähl-Zeile auf stderr,
  RTM bleibt rein auf stdout). `internal/hexagon/core/app/trace.go` liefert
  `matrix.Orphans` bereits — keine Kern-Änderung.
- **Fragment (`internal/adapter/driving/cli/print_mk.go`):** `mkTemplate` um
  `TRACE_FLAGS ?=` + `doc-trace`-Target (`--trace $(TRACE_FLAGS)`) +
  `doc-complete`-Target (`--trace --require-complete $(TRACE_FLAGS)`) erweitert;
  Determinismus unberührt (hängt nur an der eingebetteten Version).
- **Spec:** neue
  [`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
  + CR an
  [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
  (Lastenheft **0.24.0**, Historie); `spec/spezifikation.md`: die `.a`-Sektionen
  zu `--require-complete` (Exit-Semantik) und `--print-mk` (Fragment-Targets). Keine Abwärts-Links
  auf ADR/Slice ([`MR-006`](../../../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)).
- **Repo-eigenes `make completeness-check` bleibt unberührt:**
  [ADR-0017](../../adr/0017-requirements-completeness-gate.md) §2 fixiert dessen
  bash/grep-Wrapper-Mechanik; ein Umbau auf das neue Flag wäre ein ADR-Re-Eval,
  **nicht** dieser Slice (bewusst out of scope).

## 3. Definition of Done

- [ ] **Lastenheft** (0.24.0): neue
  [`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
  mit drei Akzeptanzkriterien + Out-of-Scope; CR an
  [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
  (Beschreibung/AK/Out-of-Scope) + Historie-Zeile.
- [ ] `spec/spezifikation.md`: die `.a`-Sektionen zu `--require-complete` (Exit-Semantik) und `--print-mk` (Fragment-Targets).
- [ ] Code: `--require-complete`-Flag + `comboError` + `runTrace`-Strict-Exit;
  `mkTemplate` um `TRACE_FLAGS`/`doc-trace`/`doc-complete` erweitert.
- [ ] Tests (`internal/adapter/driving/cli/cli_acceptance_test.go`):
  Happy (0 Waisen ⇒ 0), Boundary (Waise ⇒ 1, RTM bleibt, read-only),
  JSON-Gate-Pfad (⇒ 1), Negative (`--require-complete` ohne `--trace` ⇒ 2);
  Fragment-Test (`doc-trace`/`doc-complete`/`TRACE_FLAGS`).
- [ ] Doku: [`docs/user/operations.md`](../../../user/operations.md) (Fragment-
  Targets + `--require-complete`), [`docs/user/benutzerhandbuch.md`](../../../user/benutzerhandbuch.md)
  §4.13, [`CHANGELOG.md`](../../../../CHANGELOG.md).
- [ ] `make gates` grün (Coverage-Schwelle gehalten); unabhängiges Review R1;
  Closure (Slice → `done/` mit Roadmap-Flip,
  [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)).

## 4. Risiken / offene Punkte

- **Verhältnis zu [ADR-0017](../../adr/0017-requirements-completeness-gate.md):**
  Der opt-in `--require-complete` widerspricht der ADR **nicht** — die dort
  verworfene Alternative war „den **Default**-`--trace` auf Exit ≠ 0", was jeden
  GF-Zwischenstand bräche. Hier bleibt der Default advisory; nur der explizit
  angeforderte Modus failt. Provenienz als Geschichte-Zeile in der ADR
  vermerkt.
- **Determinismus** ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
  `--require-complete` ändert **nur** den Exit-Code, die RTM-Ausgabe bleibt
  byte-identisch — Boundary-Test prüft, dass die RTM trotz Exit 1 vollständig
  auf stdout steht.
- **Semantik-Grenze (ehrlich benannt):** „0 Waisen" = jede Anforderung von ≥1
  Slice **beansprucht**, nicht *done* — identisch zur dortigen Semantik; im
  Out-of-Scope von
  [`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
  benannt.
- **Fragment-Drift:** `doc-trace`/`doc-complete` müssen mit dem Image-Pin
  konsistent bleiben; das Fragment ist deterministisch aus `mkTemplate` erzeugt
  (ein Wahrheitsort).

## 5. Trigger

a-check-Bootstrap (2026-06-20 ff.): das `d-check.mk` wurde dort interim von Hand
gepflegt; slice-038 verlagerte den `doc-check`-Pin nach d-check, ließ aber die
RTM-/Vollständigkeits-Gates aus. Dieser Slice zieht `doc-trace`/`doc-complete`
nach und liefert das opt-in Gate-Flag, das `doc-complete` trägt.

## 6. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Produkt-Code `cli`/`print_mk` + Spec; Greenfield-
Default „Doc führt, Code folgt"). Keine BF-Sub-Area berührt.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung.** Neue Anforderung
[`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code):
`--trace --require-complete` bindet die RTM-Waisen-Markierung an den Exit-Code
(`matrix.Orphans > 0` ⇒ Exit 1, Zähl-Zeile auf stderr, RTM rein auf stdout;
sonst Exit 0); `comboError` weist das Flag ohne `--trace` als Nutzungsfehler ab
(Exit 2). Der Default-`--trace` bleibt advisory (Exit 0) — die Durchsetzung ist
strikt opt-in, der Vertrag von
[`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
und die Wrapper-Mechanik von
[ADR-0017](../../adr/0017-requirements-completeness-gate.md) bleiben unangetastet
(Provenienz-Zeile dort angehängt). CR an
[`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben):
das `--print-mk`-Fragment trägt zusätzlich `doc-trace` (advisory RTM) und
`doc-complete` (`--trace --require-complete`) plus die `TRACE_FLAGS`-Variable.
Kein Kern-Eingriff — `matrix.Orphans` existierte bereits
(`internal/hexagon/core/app/trace.go`); der Strict-Exit lebt im CLI-Adapter
(`internal/adapter/driving/cli/cli.go`), das Fragment in
`internal/adapter/driving/cli/print_mk.go`. **Kein ADR** (additiv; Präzedenz
slice-036/slice-038).

**Belege.** `make gates` grün (lint, test, arch-check, coverage 93,90 %, semgrep
0, doc-check, gate-consistency, planning-check); `make completeness-check` grün
(0 Waisen, 28 Anforderungen abgedeckt). Fünf neue Akzeptanztests `TestCLI044_*`
(Happy 0 Waisen⇒0, Boundary Waise⇒1 mit RTM rein auf stdout, JSON-Gate-Pfad⇒1,
Negative ohne `--trace`⇒2, Fragment-Targets). Minor-Release **v0.24.0** auf GHCR
(Run `28008942708` grün in 2m11s, Tags `v0.24.0`+`latest`), Digest-Pin
`ghcr.io/pt9912/d-check@sha256:1c28a2b7e0e624763577ecba75b027f384692ecaa8a78a6e353a1a0c1889a4f8`.

**Reviews.** R1 (`docs/reviews/2026-06-23-slice-044-r1.md`): unabhängig, **ACCEPT**
(0 HIGH/0 MEDIUM/2 LOW behoben/1 INFO won't-fix). R2
(`docs/reviews/2026-06-23-slice-044-r2.md`): Slice-Doc gegen Implementierung,
**ACCEPT**, fand F-A (Handbuch-Versionsdrift) — release-prep-gebunden aufgelöst;
der v0.24.0-Stempel (Handbuch 1.6, CHANGELOG-Datum) lebt im separaten
Release-Prep-Commit (Präzedenz `7716be3`).

**Lerneintrag.** Eine Gate-Form für Konsumenten darf weder den advisory-Default
umbiegen (die ADR-Invariante) noch Skript-Logik ins Recipe kopieren (das Verbot
aus [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)) —
der einzige saubere Schnitt ist ein **opt-in Produkt-Exit-Flag**, dessen Default
nichts ändert. Und: `link-policy: always` zählt **jede** ID-Wiederholung im
Fließtext (R2-Vorlauf: neun `id-unlinked` durch bare Wiederholungen im
Slice-Entwurf) — Wiederholungen als „die Anforderung"/„die ADR" formulieren,
nicht erneut als nackte Kennung.
