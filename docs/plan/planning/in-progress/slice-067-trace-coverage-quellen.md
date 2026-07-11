# Slice slice-067: Kuratierte Coverage-Quellen der RTM (`trace.coverage`, range-aware)

**Status:** in-progress (welle-56-trace-coverage-quellen). Doc-first-Fundament
gelegt (Lastenheft-CR, Spezifikation, [ADR-0035](../../adr/0035-trace-coverage-quellen.md)
Proposed); Implementierung + Review + Release folgen. Lifecycle
`in-progress`→`done` mit Roadmap-Flip bei Closure
([`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)).

**Welle:** welle-56-trace-coverage-quellen.

**Bezug:** [ADR-0035](../../adr/0035-trace-coverage-quellen.md)
(Coverage-Klasse, Range-Parser, Sektions-Semantik) +
[`DC-FA-COV-001`](../../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
(neue Anforderung; Spezifikation
[§`DC-FA-COV-001.a`](../../../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)).
Mit-modifiziert [`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
(Coverage-Spalte + json/yaml-Feld) und [`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
(Waise = ¬slice ∧ ¬coverage). **Lastenheft-CR** (neue Anforderung + nutzersichtbarer
Output) — Bereich `COV`, **kein neues Modul**, **kein neuer Grund-Code** (`--trace`
bleibt advisory). **Release** geplant (v0.41.0).

**Autor:** pt9912. **Datum:** 2026-07-11.

---

## 1. Ziel

Die RTM ([ADR-0034](../../adr/0034-trace-konfigurierbare-quellen.md)) leitet Abdeckung aus ADR-/Slice-Referenz-Scans ab.
Anforderungen, deren Deckung in einer **kuratierten Matrix** (`traceability.md`)
liegt, erscheinen als **falsche Waisen** — die Matrix ist weder ADR noch Slice und
nutzt Bereichs-Notation (`GG-QA-001..006`), die der Scanner nur als erste ID
erkennt. Der Konsument grid-gym: 171 „Waisen", davon ≥122 anderswo belegt.

Ziel: eine **dritte, opt-in Referenzklasse `trace.coverage`** — eine Liste
benannter kuratierter Quellen, die **range-aware** und **abschnitts-gescopt** als
Coverage einlesen, **ohne** `adrs`/`slices` zu berühren; eigene RTM-Spalte; Waise
= ¬slice ∧ ¬coverage; Default-aus byte-identisch.

## 2. Entscheidungen (aus [ADR-0035](../../adr/0035-trace-coverage-quellen.md))

- **`trace.coverage` als Liste** benannter Quellen: je Quelle `files` (explizite
  Pfade — **keine** `dir`/`file-pattern`, gegen ADR-Kontamination), `label` (fester
  Owner in eigener Coverage-Spalte), `ranges` (bool, Default true),
  `sections`/`exclude-sections`.
- **Range-Parser** (isolierte, unit-getestete Kernfunktion, parametrisiert über
  `id-pattern`): `<FAM>-AAA..BBB` inklusiv **breiten-erhaltend** + `<FAM>-AAA/BBB/CCC`;
  jede expandierte ID gegen `id-pattern` validiert (Nicht-Treffer verworfen);
  `AAA>BBB`/Breiten-Mismatch ⇒ Exit 2. Nur Kurzform-Enden.
- **Abschnitts-Scoping über die bestehende Span-Semantik** (`rules`-`excludedRanges`,
  wie `matrix.exclude-sections`): Whitelist `sections` + Blacklist `exclude-sections`,
  erst Whitelist dann Blacklist ⇒ „§27.1 ohne §27.1.1". **Kein** neuer Span-Typ.
- **Waise = ¬slice ∧ ¬coverage** (Änderung der Vollständigkeits-Prüfung); ADR-Referenz deckt
  weiter nicht. **Konditionale Coverage-Spalte** + `coverage`-json/yaml-Feld
  (omitempty) ⇒ ohne `trace.coverage` byte-identisch.
- **Fail-closed** (Exit 2): fehlende `files`-Datei, leeres `label`, ungültige Range.

## 3. Definition of Done

- [x] **Lastenheft-CR**: die neue Anforderung (Bereich `COV` in §3) + die beiden
  CLI-Modifikationen (RTM-Spalte + Waisen-Definition); **v0.41.0** + Historie.
- [x] **Spezifikation** [§`DC-FA-COV-001.a`](../../../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)
  (Verrechnung + Range-Parser + Sektionen) + §2-Schema (`trace.coverage[].*`) + Historie.
- [x] **[ADR-0035](../../adr/0035-trace-coverage-quellen.md)** (Proposed) + ADR-Index.
- [ ] **Modell** [`config.go`](../../../../internal/hexagon/core/model/config.go):
  `TraceCoverage` + `TraceConfig.Coverage`; `TraceRow.Coverage` (omitempty);
  Coverage-Aktiv-Flag der `TraceMatrix`.
- [ ] **Config-Decode** [`configyaml.go`](../../../../internal/adapter/driven/configyaml/configyaml.go):
  `rawTrace.Coverage`; `files`/`label` nicht-leer, `ranges`-Default-true (Pointer),
  fail-closed.
- [ ] **Range-Parser** ([`trace.go`](../../../../internal/hexagon/core/app/trace.go)):
  isolierte Funktion + `rules`-Section-Span-Wiederverwendung (ggf. `excludedRanges`
  exportieren/`SelectSections`-Helfer); Coverage-Scan; Waise-Neubestimmung.
- [ ] **Reporter** [`report.go`](../../../../internal/adapter/driven/report/report.go):
  konditionale Coverage-Spalte (Markdown), `coverage` in json/yaml.
- [ ] **`--print-config`** [`config_template.go`](../../../../internal/adapter/driving/cli/config_template.go):
  kommentierter `coverage`-Block.
- [ ] **Tests**: Coverage-Klasse (Label, kein Waise); Range (`001..006` ⇒ alle
  sechs) + `/`-Enum; Sektionen (exclude §27.1.1 nicht kreditiert, include-Whitelist);
  keine ADR-Kontamination; Negative (fehlende Datei / `AAA>BBB` / Breiten-Mismatch
  ⇒ Exit 2); `--require-complete` mit Coverage; **Default-aus byte-identisch**
  (kein Spalten-/Feld-Diff); Range-Parser-Unit-Tests; mutations-hart.
- [ ] **Release-Prep**: Handbuch §4.12 (Coverage) + §5 (Schema), `CHANGELOG.md`;
  slice-061/062-Harnesse grün; bare-Tag-Sweep + `version.md` nach
  [`releasing.md` §4](../../../../docs/user).
- [ ] `make gates` / `make ci` grün; **unabhängiger Impl-Review**; Verifikation an
  grid-gym-Realdaten (Waisen-Rückgang); Closure-Move + Body + **Lerneintrag**.
  **Release v0.41.0**.

## 4. Trigger

Konsumenten-Analyse grid-gym (2026-07-11): 171 „Waisen" der v0.40.0-RTM waren zu
≥122 anderswo belegt (ADR/traceability.md/Wellen); die slice-zentrische RTM
verfehlt die kuratierte Matrix. Nutzer-Vorschlag „Feature-Spec `trace.coverage`"
(9 funktionale Punkte) → auf d-checks ein-Requirement-Konvention geschnitten
(eine neue Anforderung + zwei CLI-Modifikationen), beide Sektions-Modi, voller Slice.

## 5. Offene Punkte / Risiken

- **Range-Parser = die neue Kernlogik.** Breiten-Erhaltung (`001..007`), Familie =
  Kennung ohne Trailing-Ziffern, Validierung gegen `id-pattern`, fail-closed
  (`AAA>BBB`/Breite). Gegen die echten grid-gym-Formen (`GG-ACCEPT-001..003`,
  `GG-RT-004/005`) unit-testen — die einzige Kanten-reiche Stelle.
- **Sektions-Span:** die Whitelist+Blacklist-Kombination braucht **nur** die
  bestehende `matrix`-Span-Semantik (`excludedRanges`) — beim Export/Refactor die
  Section-Extraktion nicht duplizieren; Selbst-Test gegen §27.1.1.
- **Byte-Identität:** stärkster Regressionsschutz — RTM **ohne** `trace.coverage`
  == v0.40.0 (kein Spalten-/Feld-Diff; `omitempty` auf `coverage`); die
  slice-066-Byte-Identitäts-Tests + Handbuch-E2E (5-Spalten) müssen grün bleiben.
- **`ranges: false` vs. Default true:** Pointer im Decode (nil ⇒ true), sonst
  kippt Default-Semantik.
- **Fehlende `files`-Datei ⇒ Exit 2** (anders als `adrs.dir`/`slices.dir`, wo
  Fehlen = Skip): Coverage-`files` sind **explizit benannt**, fail-closed.
