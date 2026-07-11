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

## 6. Review-Nachtrag (Doc-first R1 + Impl R2)

- **R1 Doc-first — NACHBESSERN → behoben** (inline eingearbeitet, keine separate
  Report-Datei). MEDIUM-1: der Body-Span trägt rohen Markdown-Text; ein
  umbrochenes/emphasiertes `**MUSS** NICHT` bzw. `MUSS\nNICHT` würde als `must`
  fehlklassifiziert ⇒ **Whitespace-/Markup-Normalisierung** (`normalizeBody`)
  ergänzt. MEDIUM-2: dasselbe Keyword in zwei Stufen ist nondeterministisch ⇒
  **fail-closed Exit 2**. Dazu LOW/INFO: Aktivierung = **Schlüssel-Präsenz**
  (nicht `len(levels)>0`), gatende Zahl in der `--require-complete`-Meldung, das
  RE2-ASCII-`\b`-Caveat 1:1 im Vertrag, kanonische Defaults in der Spezifikation.
- **R2 Impl — ACCEPT WITH NITS → behoben** (Report:
  [`docs/reviews/2026-07-11-slice-068-trace-modality-impl-r2.md`](../../../reviews/2026-07-11-slice-068-trace-modality-impl-r2.md)).
  MEDIUM M1: „leerer Stufen-Name ⇒ Exit 2" war ungetestet ⇒ Fall in
  `TestCLI068_Modality_NegativeConfig` (mutations-verriegelt: nur der
  `leerer Stufen-Name`-Guard löst Exit 2 aus). LOW L2: positive `unknown`-Gating-
  Richtung (`require-levels: [must, unknown]`) ⇒ `TestCLI068_Modality_UnknownGating`.
  INFO I1: [ADR-0036](../../adr/0036-trace-modality-klassifikation.md)-Default-Keyword-Beispiel mit der Impl synchronisiert
  (`SHALL NOT`/`SOLLTEN NICHT`). I4: negativer Index-Panic im Test-`Fatalf`
  entschärft. **Bewusst nicht behoben:** L1 (bare `modality:` YAML-null inaktiv —
  im Handbuch/Template durchgängig `modality: {}` gezeigt), I2 (Quelle bei aktiver
  Modalität zweimal read-only gelesen — kein Korrektheits-Impact), I3
  (`normalizeBody` strippt `*`/Backtick, nicht `_` — spec-konsistent).

Explizit sauber verifiziert (R2): Spec-Konformität, Klassifikator (längster/
frühester Treffer, Wortgrenze `MUSS`≠`musste`, Markup-Normalisierung),
Byte-Identität (`omitempty` json+yaml, `ModalityActive json:"-"`, 5-Spalten-
Fallback), Fail-closed bis Exit 2 config-zeitig, Gating-Semantik (`GatingOrphans`
= Waisen deren Stufe in `require-levels`), Hexagon-Grenzen (kein git/IO im Kern).

## 7. Closure-Notiz (nach done)

**Umgesetzt:** Die RTM (`--trace`,
[`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix))
trägt eine opt-in **Modalitäts-Klassifikation**
[`trace.requirements.modality`](../../../../spec/lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in):
aus **konfigurierbaren** Modal-Verb-Keywords (Built-in DE+EN-RFC-2119-Defaults,
`levels`/`require-levels`) klassifiziert sie jede Anforderung nach RFC-2119-Stufe
(MUSS/SOLLTE/KANN) in einer **eigenen Modality-Spalte** — **längster/frühester
Treffer** im **markup-normalisierten** Body (`**MUSS** NICHT`/`MUSS\nNICHT` ⇒
`may`), **wortgrenzen-genau** (`MUSS` ≠ `musste`), `unknown`-Fallback sichtbar.
`--require-complete` bricht **nur** auf `require-levels`
([`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
angepasst; Default `[must]`) — SOLLTE/KANN/`unknown` advisory. Fail-closed (leerer
Stufen-Name / reserviertes `unknown` / leeres Keyword / Keyword in zwei Stufen /
ungültiges `require-levels` ⇒ Exit 2). **Ohne `modality` byte-identisch** (keine
Spalte, kein Feld, unverändertes Gating;
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).

**Belege:** `make gates` grün (doc-check + lint + test + arch-check + coverage-gate
+ semgrep + gate-consistency + planning-check). **End-to-End gegen grid-gyms echtes
Lastenheft** verifiziert: die **10** Coverage-Rest-Waisen (slice-067) gaten unter
`modality: {}` + Default `require-levels: [must]` nur noch **2** — GG-MVP-004
(`DARF NICHT` ⇒ must) und GG-NONGOAL-005 (Klausel `… muessen offen austauschbar …`
⇒ must); die übrigen acht (5 KANN + 3 modalitätslose Nicht-Ziele) werden advisory.
Zwei unabhängige Reviews (Doc-first R1 NACHBESSERN + Impl R2 ACCEPT-WITH-NITS, alle
eingearbeitet, §6).

**Commit-Kette:** `e85f9b3` (doc-first) · `4643425` (feat) · `c759590`
(release-prep v0.42.0) · `dcb8f46` (Review R1+R2) · `f8c0be7` (Closure-Move) ·
Closure-Body (dieser, = Tag `v0.42.0`) · digest-backfill (folgt). **Release
v0.42.0** auf GHCR — Digest-Pin folgt bei Backfill.

**Lehre:** (i) Modalität ist eine **andere Dimension** als Deckung/Slice — eine
**eigene Spalte** + ein **eigenes Gate** (`require-levels`) halten „was ist
gefordert" und „wie stark ist es gefordert" getrennt; MUSS und KANN als gleich
gewichtete Waisen zu zählen war das eigentliche Rauschen. (ii) Eine am **rohen
Body** ansetzende Klassifikation muss **vor** dem Match normalisieren (Whitespace +
Emphasis/Code-Markup) — sonst verschluckt ein umbrochenes/emphasiertes `**MUSS**
NICHT` die Negation und fällt still auf `must` zurück (Doc-first-R1). (iii)
Überlappende Mehr-Wort-Keywords (`MUSS` ⊂ `MUSS NICHT`) verlangen **längster-
Treffer-zuerst**; bei Byte-Sortierung ist zu prüfen, dass die Umlaut-Bytes die
Längen-Ordnung nicht verdrehen (nur ASCII-Suffixe ` NICHT`/` NOT` angehängt). (iv)
Ein **konfigurierbarer** Keyword-Satz gehört als Built-in-Default in den **Kern**
(nicht in die Doku-`yaml`) — Kern-Default, Spec-Beispiel und ADR-Beispiel
konsistent halten (Impl-R2 I1: `SHALL NOT`/`SOLLTEN NICHT` nachgezogen).
