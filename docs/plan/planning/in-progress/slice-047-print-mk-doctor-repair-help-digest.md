# Slice slice-047: `--print-mk` — `doc-doctor`/`doc-repair`/`doc-help` + `DCHECK_DIGEST`

**Status:** in-progress (welle-36-print-mk-erweiterung).

**Welle:** welle-36-print-mk-erweiterung (Trigger: Auftraggeber-Wunsch, das
`--print-mk`-Fragment um Diagnose-/Reparatur-Targets, Self-Doku und einen
Digest-Komfort zu erweitern).

**Bezug:** Change Request an
[`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
(Lastenheft **0.27.0**). Exponiert die bestehenden Modi
[`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)
(`--doctor`) und
[`DC-FA-CLI-008`](../../../../spec/lastenheft.md#dc-fa-cli-008--reparatur-patch)
(`--repair`) als Fragment-Targets. **Kein ADR** (additiv, Fragment-Erweiterung wie
slice-044).

**Autor:** pt9912. **Datum:** 2026-06-23.

---

## 1. Ziel

Das `--print-mk`-Fragment bekommt drei Targets + eine Variable: `doc-doctor`
(`--doctor`), `doc-repair` (`--repair`, `git apply`-rein), `doc-help` (Self-Doku
der `doc-*`-Targets) sowie `DCHECK_DIGEST` (Digest-Override per `ifeq`, sticht den
Tag). Alle Targets werden `##`-annotiert (das `help` des Konsumenten greift sie
auf).

## 2. Entscheidungen

- **`doc-help`, nicht `help`** (Namens-Kollision mit dem Konsumenten-Makefile):
  namespaced Target + `##`-Annotationen auf allen Targets; `doc-help` listet via
  `grep … $(MAKEFILE_LIST) | sed` über die Annotationen.
- **`doc-repair` mit `@`** (Recipe-Echo unterdrückt): `make` druckt die
  Recipe-Zeile auf stdout; ohne `@` verunreinigte „docker run …" den Patch und
  bräche `git apply`. Nur `doc-repair` braucht das (Patch-Konsum).
- **`DCHECK_DIGEST` per `ifeq`:** gesetzt ⇒ `DCHECK_REF := …/d-check@$(DCHECK_DIGEST)`,
  sonst der Tag von `DCHECK_IMAGE`. Komfort gegenüber dem vollen
  `DCHECK_IMAGE`-Override.
- **Renderer-Falle:** das `mkTemplate` trägt außer dem Versions-`%s` **kein** `%`
  (das `doc-help`-Recipe nutzt `sed`, nicht `awk`-`printf`) — sonst bräche
  `fmt.Sprintf`. Im Code-Kommentar festgehalten.
- **Kein `diagrams`-Target:** das Modul läuft (sofern konfiguriert) als Teil von
  `doc-check`; ein eigenes Target wäre Redundanz (bewusst ausgelassen).
- **Spec:** CR an
  [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
  (Beschreibung/AK/Out-of-Scope, 0.27.0) + erweiterte `.a`-Spezifikationssektion.

## 3. Definition of Done

- [x] **Lastenheft-CR** (0.27.0): Beschreibung (sechs Targets + zwei Variablen),
  Happy-/Boundary-AK, Out-of-Scope (Targets + `help`-Kollision + Digest-Form), §7.
- [x] `spec/spezifikation.md` `.a`-Sektion (DCHECK_DIGEST/`ifeq`, sechs
  `##`-Targets, `@`- und Namespacing-Begründung).
- [x] Code `print_mk.go`: `mkTemplate` neu (DCHECK_DIGEST/`ifeq`, sechs Targets,
  `##`, `doc-repair` `@`, `doc-help` `sed`); `%`-frei außer `%s`.
- [x] Tests: `TestCLI038/044` an `##`-Format angepasst; `TestCLI047_PrintMK_NeueTargets`
  (neue Targets + `DCHECK_DIGEST` + `ifeq` + `@`). Fragment-Parse-Check (`make -n`)
  + Digest-Override manuell belegt.
- [ ] `--print-config`-/Handbuch-Parität prüfen (§4.13). [`CHANGELOG.md`](../../../../CHANGELOG.md).
- [ ] `make gates` grün; unabhängiges Review R1; Closure (Move nach `done/` +
  Roadmap-Flip, [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)).

## 4. Risiken / offene Punkte

- **Konsumenten-Kollision:** `doc-help` ist namespaced; `DCHECK_DIGEST`/`DCHECK_REF`
  sind spezifische Namen — Restkollisions-Risiko minimal, aber benannt.
- **`fmt.Sprintf`-`%`-Falle:** ein künftiges `%` im Template (z. B. `awk printf`)
  bräche den Renderer still — Code-Kommentar + `TestCLI047` (Render-Erfolg) schützen.
- **Determinismus** ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
  statisches Template, hängt nur an der Version.

## 5. Trigger

Auftraggeber: „Können wir `--print-mk` erweitern?" → vier Ideen (`doc-repair`,
`doc-doctor`, Self-Doc, Digest-Komfort), alle bestätigt (2026-06-23).

## 6. Sub-Area-Modus-Begründung

GF (Produkt-Code `cli` + Spec; „Doc führt, Code folgt"). Keine BF-Sub-Area.

## 7. Closure-Notiz (nach `done/`)

_(Wird beim Move nach `done/` gefüllt — Belege: `make gates`-Ausgabe,
Fragment-Parse-Check, unabhängiges Review R1, ggf. Release.)_
