# Slice slice-009: coverage-gate, gate-consistency + DC-QA-01-Benchmark

**Status:** done.

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
Doku↔Makefile-Meta-Gate), und die [`DC-QA-01`](../../../../spec/lastenheft.md#dc-qa-01--performance)-Performance-Anforderung
hat ihre Benchmark-Definition in der Spezifikation plus eine
dokumentierte Messung.

## 2. Definition of Done

- [x] `make coverage-gate`: Dockerfile-Stage + Schwellen-Skript
  (Variable `THRESHOLD ?= 85`); Sensors-Bindung als
  Kalibrierungs-Bindung „Schwelle 85 %, welle-03 done → 90 %";
  Verfehlung nach Trigger ⇒ Carveout-Pflicht (Kurs-Hard-Rule)
  dokumentiert. *(Ramp beim Wellen-Closure vollzogen: Default jetzt
  90 %, Ist 92,9 % — siehe Closure-Notiz.)*
- [x] `make gate-consistency`: neues Target
  (`tools/gate-consistency.sh`), das alle in
  [`harness/README.md`](../../../../harness/README.md) §Sensors und
  [`AGENTS.md`](../../../../AGENTS.md) §4 als real dokumentierten
  Make-Targets gegen das Makefile abgleicht; jedes fehlende Target →
  Exit 1 mit Auflistung (Meta-Gate gegen Harness-Lügen); mit
  Negativ-Test (absichtlich dokumentiertes Phantom-Target lässt das
  Gate nachweislich feuern, analog `verify-depguard`-Idee).
  Zusätzlich prüft das Meta-Gate die [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Zusage des
  Netzlos-Gates: die `modules`-Liste der
  [`.d-check.yml`](../../../../.d-check.yml) muss alle Module außer
  `external` enthalten — sonst verliert der `--network none`-Lauf
  still seine Beweis-Aussage (Review-R1-Finding zu slice-008:
  Config-Kopplung des QA-03-Gates). *(Umgesetzt inkl. Gegenrichtung:
  Makefile-Targets müssen in `AGENTS.md` §4 gelistet sein.)*
- [x] Spez-Schuld eingelöst: `spec/spezifikation.md` erhält einen
  Abschnitt `DC-QA-01.a — Benchmark` mit (1) Fixture-Spezifikation <!-- d-check:ignore (Sektions-Titel-Zitat, kein Requirement-Link) -->
  (generiert: 1.000 Markdown-Dateien, ≤ 20 MB, definierter
  Link-/Heading-Mix), (2) Messprotokoll (Default-Module, N ≥ 3 Läufe
  im Container, Median zählt), (3) Pass-Kriterium (< 5 s) — das
  „(folgt)" im Lastenheft ist damit erfüllt (0.2.3, redaktionell);
  die durchgeführte Messung ist in der Closure-Notiz dokumentiert.
- [x] Beide Gates in `make gates` aggregiert; „Nicht behauptet"-Listen
  in AGENTS/harness entsprechend verkürzt.
- [x] `make gates` grün; [`CHANGELOG.md`](../../../../CHANGELOG.md);
  Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`Dockerfile`](../../../../Dockerfile) (coverage-Stage), `tools/coverage-gate.sh` | neu | bootstrap-aware Schwelle (u-boot-Muster) |
| `tools/gate-consistency.sh` | neu | dokumentierte Targets ↔ Makefile |
| [`Makefile`](../../../../Makefile) | update | Targets + gates-Aggregation |
| [`spec/spezifikation.md`](../../../../spec/spezifikation.md) | update | [`DC-QA-01`](../../../../spec/lastenheft.md#dc-qa-01--performance)-Benchmark-Definition (+ Historie) |
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
  im Gate, [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)); die Coverage misst das Modul über seine
  Tests, nicht über das Dogfooding.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung:** Commit `78716b2` (Gates, Benchmark, Spez-Schuld) plus
Kalibrierungs-Schaltung im Closure-Commit. **Messungen:** Coverage
**92,9 %** (Erstlauf gegen Schwelle 85, Ramp auf 90 vollzogen —
Erhöhung, kein ADR nötig, `AGENTS.md` §3.6 betrifft Senkungen);
[`DC-QA-01`](../../../../spec/lastenheft.md#dc-qa-01--performance)-Benchmark **Median 551 ms** (Läufe 559/517/551 ms, 1.000
Dateien / 8 MB, inkl. Container-Start) — Faktor ~9 unter dem
5-s-Budget.

- **Was hat funktioniert:** Das u-boot-Coverage-Muster (Stage +
  Schwellen-Skript, `-coverpkg` über Paketgrenzen) ließ sich
  unverändert übernehmen; der [`DC-QA-01`](../../../../spec/lastenheft.md#dc-qa-01--performance)-Benchmark brauchte kein
  eigenes Framework — deterministisches Shell-Fixture plus drei
  getimte Container-Läufe genügen dem Messprotokoll.
- **Anders als geplant:** (a) Das Meta-Gate prüft zusätzlich die
  *Gegenrichtung* (jedes Makefile-Target muss in `AGENTS.md` §4
  stehen) — die AGENTS-eigene „Nur hier gelistete"-Zusage war sonst
  unbewacht. (b) Der Ramp 85 → 90 wurde direkt beim Wellen-Closure
  vollzogen statt offen gelassen (Ist 92,9 % > 90).
- **Steering-Loop-Lerneintrag:** Der *Erstlauf* des Meta-Gates fand
  sofort sechs Diskrepanzen — drei frische Targets waren noch nicht
  dokumentiert, und drei Bestands-Targets (`run`, `compile`, `clean`)
  steckten in kombinierten Tabellen-Zellen, die der erste Parser
  übersah. Beleg für zwei Kurs-Thesen: (1) ein Meta-Gate gegen
  Harness-Lügen zahlt sich beim ersten Lauf aus; (2) der Sensor
  selbst braucht den Steering-Loop (Parser-Härtung auf alle
  `make`-Tokens pro Tabellenzeile). Außerdem: den Negativ-Test als
  Selbsttest in den Gate-Lauf einzubetten ist billiger als ein
  separater Test und beweist die Wirksamkeit bei jedem `make gates`.
- **Folge-Slices:** keine — welle-03 ist mit diesem Slice
  abgeschlossen (Wellen-Closure-Trigger erfüllt: alle fünf Module
  implementiert und getestet, Selbstkonfiguration aktiv,
  coverage-gate + gate-consistency in `gates` aggregiert). welle-04
  (Distribution + Migration) übernimmt als aktive Welle.

**Nachkalibrierung (2026-06-11, nach Closure):** Auf User-Wunsch
Schwelle 93 % — gezielter Test-Ausbau hob das Ist von 92,9 % auf
95,1 % (Commit `11731c7`), erst danach wurde geschaltet (statt das
Gate rot zu stellen). Schwellen-Historie: 85 (DoD) → 90 (Ramp,
welle-03 done) → 93 (Kalibrierung); maßgeblich dokumentiert in der
Sensors-Bindung (`harness/README.md` §Sensors). **Review R1
(gebündelt mit slice-010, nach Closure):** Median-Position aus `RUNS`
abgeleitet statt hart `sed -n 2p`; `--cpus 2` im Benchmark +
Spez-Präzisierung (2-vCPU-Normierung aus dem Lastenheft);
Meta-Gate-Parser um Mehrfach-Target-Zeilen erweitert (mit
Parser-Selbsttest).

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Gate-/Spez-Arbeit; siehe Kurs Modul 5
§Worked Mini-Example).
