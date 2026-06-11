# Slice slice-009: coverage-gate, gate-consistency + DC-QA-01-Benchmark

**Status:** open.

**Welle:** welle-03-regelmodule (Gate-Ausbau, Abschluss).

**Bezug:** [`DC-QA-01`](../../../../spec/lastenheft.md#dc-qa-01--performance)
(inkl. Einlösung der Spez-Schuld „Definition folgt"),
[ADR-0006](../../adr/0006-lint-profil-solid.md) (Gate-Disziplin);
Kurs-Modul 13 (Kalibrierungs-Bindung, Meta-Gate).

**Autor:** pt9912. **Datum:** 2026-06-10.

---

## 1. Ziel

Die beiden ausstehenden Gates existieren (`coverage-gate`
bootstrap-aware mit Ramp, `gate-consistency` als
Doku↔Makefile-Meta-Gate), und die `DC-QA-01`-Performance-Anforderung
hat ihre Benchmark-Definition in der Spezifikation plus eine
dokumentierte Messung.

## 2. Definition of Done

- [ ] `make coverage-gate`: Dockerfile-Stage + Schwellen-Skript
  (Variable `THRESHOLD ?= 85`); Sensors-Bindung als
  Kalibrierungs-Bindung „Schwelle 85 %, welle-03 done → 90 %";
  Verfehlung nach Trigger ⇒ Carveout-Pflicht (Kurs-Hard-Rule)
  dokumentiert.
- [ ] `make gate-consistency`: neues Target
  (`tools/gate-consistency.sh`), das alle in
  [`harness/README.md`](../../../../harness/README.md) §Sensors und
  [`AGENTS.md`](../../../../AGENTS.md) §4 als real dokumentierten
  Make-Targets gegen das Makefile abgleicht; jedes fehlende Target →
  Exit 1 mit Auflistung (Meta-Gate gegen Harness-Lügen); mit
  Negativ-Test (absichtlich dokumentiertes Phantom-Target lässt das
  Gate nachweislich feuern, analog `verify-depguard`-Idee).
  Zusätzlich prüft das Meta-Gate die `DC-QA-03`-Zusage des
  Netzlos-Gates: die `modules`-Liste der
  [`.d-check.yml`](../../../../.d-check.yml) muss alle Module außer
  `external` enthalten — sonst verliert der `--network none`-Lauf
  still seine Beweis-Aussage (Review-R1-Finding zu slice-008:
  Config-Kopplung des QA-03-Gates).
- [ ] Spez-Schuld eingelöst: `spec/spezifikation.md` erhält einen
  Abschnitt `DC-QA-01.a — Benchmark` mit (1) Fixture-Spezifikation
  (generiert: 1.000 Markdown-Dateien, ≤ 20 MB, definierter
  Link-/Heading-Mix), (2) Messprotokoll (Default-Module, N ≥ 3 Läufe
  im Container, Median zählt), (3) Pass-Kriterium (< 5 s) — das
  „(folgt)" im Lastenheft ist damit erfüllt; eine durchgeführte
  Messung ist in der Closure-Notiz dokumentiert.
- [ ] Beide Gates in `make gates` aggregiert; „Nicht behauptet"-Listen
  in AGENTS/harness entsprechend verkürzt.
- [ ] `make gates` grün; [`CHANGELOG.md`](../../../../CHANGELOG.md);
  Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`Dockerfile`](../../../../Dockerfile) (coverage-Stage), `tools/coverage-gate.sh` | neu | bootstrap-aware Schwelle (u-boot-Muster) |
| `tools/gate-consistency.sh` | neu | dokumentierte Targets ↔ Makefile |
| [`Makefile`](../../../../Makefile) | update | Targets + gates-Aggregation |
| [`spec/spezifikation.md`](../../../../spec/spezifikation.md) | update | `DC-QA-01`-Benchmark-Definition (+ Historie) |
| `tools/bench-fixture.sh` o. ä. | neu | generiertes Fixture-Repo für die Messung |

## 4. Trigger

slice-006, slice-007 und slice-008 done (Coverage-Endausbau misst
alle fünf Module; gate-consistency friert den finalen Sensor-Satz
der Welle ein).

## 5. Closure-Trigger

DoD vollständig + Commit(s) auf `main` + Closure-Notiz geschrieben —
damit ist zugleich der welle-03-Closure-Trigger erfüllt.

## 6. Risiken und offene Punkte

- Coverage im Docker-Stage muss `-coverpkg ./internal/...` über die
  Paketgrenzen messen (u-boot-Muster), sonst zählt nur
  paket-lokale Abdeckung.
- Der Benchmark misst im Container — Schwankungen der Host-Last;
  Schwelle (< 5 s) hat Puffer, Messung mehrfach ausführen.
- `make doc-check` aktiviert `external` bewusst **nicht** (kein Netz
  im Gate, `DC-QA-03`); die Coverage misst das Modul über seine
  Tests, nicht über das Dogfooding.

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Gate-/Spez-Arbeit; siehe Kurs Modul 5
§Worked Mini-Example).
