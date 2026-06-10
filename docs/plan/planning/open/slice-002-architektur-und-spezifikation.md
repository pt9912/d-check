# Slice slice-002: Architektur- und Spezifikations-Stratum

**Status:** open.

**Welle:** welle-01-fundament.

**Bezug:** [`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate),
[`DC-FA-ANCH-001`](../../../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors),
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus); ADR-0001–0003 (aus slice-001).

**Autor:** pt9912. **Datum:** 2026-06-10.

---

## 1. Ziel

Die Spec-Straten 2 und 3 existieren: `spec/spezifikation.md`
(technische Details) und `spec/architecture.md` (Komponenten- und
Schichtensicht), konsistent mit dem Lastenheft und den
Fundament-ADRs.

## 2. Definition of Done

- [ ] `spec/spezifikation.md` existiert; enthält mindestens: Befund-Format (Text + JSON-Schema, [`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)), GitHub-Slug-Algorithmus ([`DC-FA-ANCH-001`](../../../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors)), Default-Konfiguration und `.d-check.yml`-Schema ([`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)), Determinismus-Festlegungen (Sortierung, [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- [ ] `spec/architecture.md` existiert; Komponenten-Übersicht (Scanner, Regelmodule, Reporter, Config), Schichten-Constraints; sprach- und meilensteinfrei ([`AGENTS.md`](../../../../AGENTS.md) §3.4).
- [ ] Source-Precedence-Tabellen in [`AGENTS.md`](../../../../AGENTS.md) und [`harness/README.md`](../../../../harness/README.md) von „geplant" auf echte Links umgestellt (Auflösung des Geplant-Zustands aus [`MR-001`](../../../../harness/conventions.md#mr-001--source-precedence-mit-eigener-spezifikations-schicht)).
- [ ] `make gates` grün (`doc-check` prüft die neuen Querverweise mit).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` | neu | Rang-2-Stratum, von ADRs schärfbar |
| `spec/architecture.md` | neu | Rang-3-Stratum, Zielbild für Implementierungs-Slices |
| [`AGENTS.md`](../../../../AGENTS.md) | update | Precedence-Tabelle: Links statt „geplant" |
| [`harness/README.md`](../../../../harness/README.md) | update | Precedence-/Guides-Tabelle: Links statt „geplant" |

## 4. Trigger

slice-001 done (ADRs bestimmen Architektur- und Spezifikationsinhalte).

## 5. Closure-Trigger

DoD vollständig + Commit(s) auf `main` + Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Versuchung, Implementierungsdetails ins Lastenheft zurückzuschreiben —
  Schärfungen gehören in `spezifikation.md` (Spec-Stratifizierung).
- Das JSON-Schema des `--json`-Outputs muss zum Lastenheft-Mindestumfang
  passen, darf ihn aber erweitern.

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (reine Doku-Arbeit; siehe Kurs Modul 5
§Worked Mini-Example).
