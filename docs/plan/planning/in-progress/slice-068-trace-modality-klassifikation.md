# Slice slice-068: Modalitäts-Klassifikation der RTM-Anforderungen (`trace.requirements.modality`)

**Status:** in-progress (welle-57-trace-modality-klassifikation). Doc-first-
Fundament gelegt (Lastenheft-CR, Spezifikation, [ADR-0036](../../adr/0036-trace-modality-klassifikation.md)
Proposed); Implementierung + Review + Release folgen. Lifecycle `in-progress`→`done`
mit Roadmap-Flip bei Closure
([`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)).

**Welle:** welle-57-trace-modality-klassifikation.

**Bezug:** [ADR-0036](../../adr/0036-trace-modality-klassifikation.md)
(konfigurierbare Keywords + Defaults, erster/längster Treffer, `unknown`, Gating) +
[`DC-FA-MOD-001`](../../../../spec/lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in)
(neue Anforderung; Spezifikation
[§`DC-FA-MOD-001.a`](../../../../spec/spezifikation.md#dc-fa-mod-001a--modalitäts-klassifikation-tracerequirementsmodality)).
Mit-modifiziert [`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
(Modality-Spalte) und [`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
(gatende Stufen). **Lastenheft-CR** — Bereich `MOD`, **kein neues Modul**, **kein
neuer Grund-Code**. **Release** geplant (v0.42.0).

**Autor:** pt9912. **Datum:** 2026-07-11.

---

## 1. Ziel

Die RTM behandelt jede Anforderung gleich, obwohl das Lastenheft
**RFC-2119-Modalität** trägt (grid-gym: 112 MUSS · 50 SOLLTE · 9 KANN · 3 DARF
NICHT). Eine reine **KANN**-Anforderung ohne Slice/Coverage erscheint wie eine
unabgedeckte **MUSS**-Pflicht — die 10 Coverage-Rest-Waisen von grid-gym sind
5× KANN (Future) + 4× Nicht-Ziele + 1× DARF NICHT.

Ziel: eine opt-in **Modalitäts-Klassifikation** — konfigurierbare Modal-Verb-
Keywords (mit DE+EN-Defaults) klassifizieren jede Anforderung; eigene
Modality-Spalte; `--require-complete` bricht **nur** auf verpflichtende Stufen;
default-aus byte-identisch.

## 2. Entscheidungen (aus [ADR-0036](../../adr/0036-trace-modality-klassifikation.md))

- **`trace.requirements.modality`** opt-in: `levels` (Stufe→Keywords, Built-in
  DE+EN-RFC-2119-Default inkl. Negationen `DARF NICHT`→must / `MUSS NICHT`→may) +
  `require-levels` (welche Stufen gaten, Default `[must]`).
- **Klassifikation:** erster/**längster** Keyword-Treffer im Body-Abschnitt
  (Span-Mechanik von `trace.coverage`), case-insensitiv, wortgrenzen-genau; eine
  Stufe je Anforderung; kein Treffer ⇒ **`unknown`** (sichtbar).
- **Gating:** Waise gatet nur, wenn Stufe ∈ `require-levels`; ohne `modality`
  gaten **alle** Waisen (byte-identisch). `unknown` gatet nur, wenn explizit
  gelistet (bewusste, sichtbare fail-open-Grenze).
- **Konditionale Modality-Spalte** + `modality`-json/yaml-Feld (omitempty) ⇒ ohne
  Block byte-identisch. **Fail-closed:** leerer Stufen-Name/Keyword, `require-levels`
  weder Stufe noch `unknown` ⇒ Exit 2.

## 3. Definition of Done

- [x] **Lastenheft-CR**: die neue Anforderung (Bereich `MOD`) + die beiden
  CLI-Modifikationen (Modality-Spalte + gatende Stufen); **v0.42.0** + Historie.
- [x] **Spezifikation** [§`DC-FA-MOD-001.a`](../../../../spec/spezifikation.md#dc-fa-mod-001a--modalitäts-klassifikation-tracerequirementsmodality)
  (Klassifikation 1–5) + §2-Schema (`modality.levels`/`require-levels`) + Historie.
- [x] **[ADR-0036](../../adr/0036-trace-modality-klassifikation.md)** (Proposed) + ADR-Index.
- [ ] **Modell** [`config.go`](../../../../internal/hexagon/core/model/config.go):
  `TraceModality` (Levels/RequireLevels) + `TraceConfig.Modality`; `TraceRow.Modality`
  (omitempty); `TraceMatrix.ModalityActive`.
- [ ] **Config-Decode** [`configyaml.go`](../../../../internal/adapter/driven/configyaml/configyaml.go):
  `modality`-Block; Default-Keywords bei leerem `levels`; `require-levels` gegen
  Stufen+`unknown` validiert; fail-closed.
- [ ] **Klassifikator** [`trace.go`](../../../../internal/hexagon/core/app/trace.go):
  Body-Abschnitt (Span-Wiederverwendung), erster/längster Treffer, `unknown`;
  Gating-Filter für `--require-complete`.
- [ ] **Reporter** [`report.go`](../../../../internal/adapter/driven/report/report.go):
  konditionale Modality-Spalte (nach Coverage, vor Status); `modality` in json/yaml.
- [ ] **`--print-config`** [`config_template.go`](../../../../internal/adapter/driving/cli/config_template.go):
  kommentierter `modality`-Block.
- [ ] **Tests**: Klassifikation (must/should/may aus Body); **Längster-Treffer**
  (`MUSS NICHT`→may, `DARF NICHT`→must); `unknown`; Gating nach `require-levels`
  (KANN advisory, MUSS gatet); Negative-Config (leeres Keyword / ungültiges
  `require-levels` ⇒ Exit 2); Default-aus byte-identisch; Klassifikator-Unit-Tests;
  mutations-hart.
- [ ] **Release-Prep**: Handbuch §4.12 (Modality) + §5 (Schema), `CHANGELOG.md`;
  slice-061/062-Harnesse grün; bare-Tag-Sweep + `version.md`.
- [ ] `make gates` / `make ci` grün; **unabhängige Reviews** (Doc-first + Impl);
  Verifikation grid-gym (KANN/Nicht-Ziel advisory, GG-MVP-004 gatet); Closure +
  **Lerneintrag**. **Release v0.42.0**.

## 4. Trigger

Nutzer-Frage (2026-07-11) „unterscheidet Traceability MUSS/SOLLTE?" + Hinweis, dass
das Lastenheft MUSS/SOLLTE/DARF NICHT trägt. Meine erste Analyse („Modalität hülfe
nicht") war durch einen Body-Scan-Bug falsch — korrigiert: die 10 Rest-Waisen
sind 5 KANN + 4 Nicht-Ziele + 1 DARF NICHT. Nutzer bestätigte: konfigurierbare
Modal-Verben mit Defaults + opt-in, voller Slice.

## 5. Offene Punkte / Risiken

- **Fail-open bei `unknown`:** ein echtes MUSS mit unaufgeführtem Verb fiele auf
  `unknown` (advisory) — Gegenmittel: umfassende Defaults, **sichtbare** Spalte,
  `require-levels: [must, unknown]`-Strikt-Modus. In Tests + Handbuch benennen.
- **Negations-Kanten:** `MUSS NICHT` (may) ≠ `DARF NICHT` (must) — **Längster-
  Treffer-zuerst** ist Pflicht, sonst verschluckt `MUSS` die Negation. Der stärkste
  Test.
- **Wortgrenzen:** `MUSS` darf nicht `musste`/`Mussten` matchen; Keyword-Match
  case-insensitiv + wortgrenzen-genau (Unicode-Wortgrenzen in RE2 beachten).
- **Byte-Identität:** RTM/`--require-complete` ohne `modality` == v0.41.0
  (kein Spalten-/Feld-/Gate-Diff); die slice-066/067-Byte-Identitäts-Tests +
  Handbuch-E2E müssen grün bleiben.
- **Default-Keywords als Config-Rand:** die eingebaute Menge lebt im Kern (nicht
  in der Doku-`yaml`) — Doku-Beispiel + Kern-Default konsistent halten.
