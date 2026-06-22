# Slice slice-042: Requirements-Completeness-Gate (Waisen als Closure-Invariante)

**Status:** in-progress (Planung; welle-31-requirements-completeness).

**Welle:** welle-31-requirements-completeness (Trigger: Nutzer 2026-06-22 —
„damit kann man prüfen, ob die Arbeit abgeschlossen ist"; die RTM
(`--trace`, slice-036) markiert Waisen, **erzwingt** sie aber nicht).

**Bezug:** [`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
(RTM-Modus `--trace`: je Anforderung referenzierende ADRs/Slices +
Waisen-Markierung; advisory, Exit 0). Mechanik-Präzedenz: die Meta-Gates
[slice-040](../done/slice-040-planning-consistency-gate.md) (`planning-check`)
und [slice-041](../done/slice-041-adr-immutable-gate.md) (`adr-check`) —
`tools/`-Skript + Negativ-Selbsttest + Doku-Kopplung. Bindepunkt-Policy in
eigener Prozess-ADR ([ADR-0017](../../adr/0017-requirements-completeness-gate.md),
entsteht mit diesem Slice).

**Autor:** pt9912. **Datum:** 2026-06-22.

---

## 1. Ziel

Ein Wächter `make completeness-check`, der bei **Requirements-Waisen**
(Anforderung ohne referenzierenden Slice) fehlschlägt — die *Abdeckungs*-Hälfte
von „fertig" maschinell erzwingen. `--trace` markiert Waisen heute nur
(advisory, Exit 0, spec-fixiert in
[`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix));
dieser Slice macht „0 Waisen" zu einer **Closure-Invariante**, ohne den
Advisory-Vertrag von `--trace` anzutasten.

## 2. Entscheidungen

- **Bindepunkt: Closure, nicht per Commit.** Gebunden an `make fullbuild`
  (volle Closure vor Welle-Merge/Release), bewusst **NICHT** in
  `make gates`/`make ci`. Begründung: Greenfield „Doc führt, Code folgt" —
  eine frische Anforderung ist legitim **transient Waise**, bis ihr Slice
  landet; ein per-Commit-Gate bräche den GF-Ablauf. Schwester-Logik zu
  `trace-check`/`adr-check`, die bewusst eigene Bindepunkte außerhalb von
  `gates` haben.
- **Mechanik: dünner Wrapper über `--trace --json`, eine Skript-Wahrheit.**
  Das Gate ruft das Runtime-Image `d-check --trace --json` (`--network none`,
  ro-Mount, wie `make trace`/`doc-check`) und liest das Feld `orphans` (int)
  des RTM-JSON: `orphans > 0` ⇒ FAIL mit der Liste der Waisen-IDs
  (`requirements[].orphan == true`). `--trace` bleibt unangetastet Exit 0
  (kein Spec-Change an
  [`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)).
  **Parsing rein mit bash/grep** (kein `jq`/`python` — keine
  Host-Binary-Abhängigkeit, konsistent mit den bestehenden Gate-Skripten unter
  `tools/`); ein fehlendes oder nicht als Ganzzahl matchbares `orphans` ⇒ FAIL,
  nicht stilles „0" (ein `jq '.orphans // 0'`-Notnagel ist verboten — er macht
  aus fehlendem Feld ein grünes „0"). R1-F-2.
- **ADR: ja** ([ADR-0017](../../adr/0017-requirements-completeness-gate.md),
  Proposed → Accepted nach Review) — pinnt die Bindepunkt-Policy (Closure-only)
  + die Advisory/Gate-Trennung; schärft keine Spec-Stelle (Prozess-ADR, analog
  [ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md)/[ADR-0016](../../adr/0016-adr-immutable-gate.md)).
- **Negativ-Selbsttest (fail-closed, beide Richtungen — R1-F-1):** auf der
  Parser-Kernlogik (synthetisches JSON, ohne Container) — testet ausdrücklich
  die Stilles-Grün-Vektoren, nicht nur den Happy-Pfad: `orphans: 1` feuert ·
  `orphans: 0` feuert nicht · **leeres stdout** (Image-Exit ≠ 0 / leerer Lauf)
  ⇒ FAIL · **JSON ohne `orphans`-Feld** oder nicht-numerischer Wert ⇒ FAIL ·
  kaputtes JSON ⇒ FAIL. Läuft bei jedem Gate-Lauf — wie
  `gate-consistency`/`adr-check`. Hintergrund: dieselbe Falle wie R1 zu
  slice-040 (Heading-Guard) / slice-041 (Status-Strip) — der Wächter muss seine
  Versagensrichtung selbst beweisen, sonst silent-green.

## 3. Definition of Done

Geplante Artefakte (Pfad gefenced, da noch nicht existent):

```text
tools/completeness-check.sh   # Waechter: --trace --json lesen, orphans>0 => FAIL; + Negativ-Selbsttest
```

- [ ] `completeness-check.sh` (unter `tools/`): ruft das Image `--trace --json`,
  parst `orphans` **mit bash/grep** (kein `jq`/`python`), listet die Waisen-IDs
  bei Fehlschlag; **Negativ-Selbsttest in beide Richtungen** — `orphans: 1`
  feuert, `orphans: 0` nicht, **und** leeres stdout / fehlendes `orphans`-Feld /
  nicht-numerischer Wert / kaputtes JSON ⇒ FAIL (kein stilles „0"); Selbsttest
  bei jedem Lauf, fail-closed.
- [ ] `make completeness-check` (`build`-Prerequisite; ro-Mount +
  `--network none`); **nicht** in `gates`/`ci`; eingehängt in `make fullbuild`;
  `.PHONY` + help.
- [ ] Doku-Sync: [`AGENTS.md` §4](../../../../AGENTS.md#4-quality-gates) **und**
  [`harness/README.md` §Sensors](../../../../harness/README.md#sensors-feedback-gates)
  (sonst rotes `make gate-consistency`). In der Gate-Taxonomie-Tabelle als
  **Meta-/Governance-Gate** mit **neuem dritten Bindepunkt „Closure"** (neben
  „in `gates`" und „Commit-/Diff-Bindepunkt") explizit eintragen — diese
  Bindepunkt-Klasse ist eine Form-Setzung, die `gate-consistency` **nicht**
  prüft (R1-F-3), daher bewusst benannt.
- [ ] Prozess-ADR ([ADR-0017](../../adr/0017-requirements-completeness-gate.md)
  Proposed → Accepted nach Review) + ADR-Index.
- [ ] `make gates` grün; `make completeness-check` grün (Ist: 0 Waisen);
  Akzeptanz: Wegwerf-Fixture mit Waise → Exit ≠ 0. Unabhängiges Review R1;
  Closure (ADR → Accepted, Slice → `done/` mit Roadmap-Flip,
  [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)).

## 4. Risiken / offene Punkte

- **JSON-Form-Kopplung:** das Gate parst `orphans` (int) aus `--trace --json`
  ([`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix));
  ändert sich das RTM-Schema, muss das Skript mit (Re-Eval-Trigger; fail-closed
  bei Parse-Fehler).
- **Semantik-Grenze (ehrlich):** „0 Waisen" = jede Anforderung von ≥1 Slice
  *beansprucht* — **nicht**, dass der Slice *done* ist (Slice-Status
  out-of-scope in
  [`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix))
  noch dass die Implementierung die Anforderung *erfüllt* (das prüfen
  Tests/Verification). Das Gate ist die Abdeckungs-, nicht die
  Fertigstellungs-Garantie — dokumentieren, nicht überverkaufen.
- **fullbuild-Härte:** eine bewusst vorgezogene Zukunfts-Anforderung
  (Lastenheft vor Slice) failt dann `make fullbuild`/Release — gewollt (kein
  Release mit unimplementierter Vertrags-Anforderung), als Konsequenz benannt.

## 5. Trigger

Nutzer 2026-06-22: „damit kann man prüfen, ob die Arbeit abgeschlossen ist" —
die RTM-Waisen-Markierung (slice-036) soll als maschinelle Abschluss-Invariante
erzwingbar werden. Schwester-Gate zu slice-040/041 (dokumentierte → erzwungene
Harness-Regel).

## 6. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Harness-Mechanik/Doku; Greenfield-Default;
`tools/`-nahe Konvention).

## 7. Closure-Notiz (nach `done/`)

_(folgt mit der Umsetzung.)_
