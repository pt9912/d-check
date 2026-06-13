# Slice slice-002: Architektur- und Spezifikations-Stratum

**Status:** done.

**Welle:** welle-01-fundament.

**Bezug:** [`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate),
[`DC-FA-ANCH-001`](../../../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors),
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus); [ADR-0001](../../adr/0001-implementierungssprache.md)–[ADR-0004](../../adr/0004-architektur-pattern-hexagonal.md) (aus slice-001).

**Autor:** pt9912. **Datum:** 2026-06-10.

---

## 1. Ziel

Die Spec-Straten 2 und 3 existieren: `spec/spezifikation.md`
(technische Details) und `spec/architecture.md` (Komponenten- und
Schichtensicht), konsistent mit dem Lastenheft und den
Fundament-ADRs.

## 2. Definition of Done

- [x] [`spec/spezifikation.md`](../../../../spec/spezifikation.md) existiert; enthält mindestens: Befund-Format (Text + JSON-Schema, [`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)), GitHub-Slug-Algorithmus ([`DC-FA-ANCH-001`](../../../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors)), Default-Konfiguration und `.d-check.yml`-Schema ([`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)), Determinismus-Festlegungen (Sortierung, [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- [x] [`spec/architecture.md`](../../../../spec/architecture.md) existiert; hexagonaler Schnitt gemäß [ADR-0004](../../adr/0004-architektur-pattern-hexagonal.md) (Kern: Scanner, Regelmodule, Befund-Modell; Ports: Filesystem, HTTP, Config, Reporter; driving: CLI), Schichten-Constraints als Grundlage der `arch-check`-Fitness-Function; sprach- und meilensteinfrei ([`AGENTS.md`](../../../../AGENTS.md) §3.4).
- [x] Source-Precedence-Tabellen in [`AGENTS.md`](../../../../AGENTS.md) und [`harness/README.md`](../../../../harness/README.md) von „geplant" auf echte Links umgestellt (Auflösung des Geplant-Zustands aus [`MR-001`](../../../../harness/conventions.md#mr-001--source-precedence-mit-eigener-spezifikations-schicht)).
- [x] `make gates` grün (`doc-check` prüft die neuen Querverweise mit).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`spec/spezifikation.md`](../../../../spec/spezifikation.md) | neu | Rang-2-Stratum, von ADRs schärfbar |
| [`spec/architecture.md`](../../../../spec/architecture.md) | neu | Rang-3-Stratum, Zielbild für Implementierungs-Slices |
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

**Umsetzung:** Commit `e13b60b` (spezifikation.md, architecture.md,
Precedence-/Guides-Umstellung in `AGENTS.md` + `harness/README.md`).

- **Was hat funktioniert:** Die ADRs aus slice-001 haben beide Straten
  fast vollständig determiniert — die Spezifikation musste nur noch
  normieren (Slug-Schritte, Grund-Codes, Schema-Constraints), nicht
  entscheiden. Die Import-Constraint-Tabelle der Architektur ist
  direkt als `arch-check`-Quelle formuliert (slice-003 kann sie
  mechanisch übersetzen).
- **Anders als geplant:** `AGENTS.md` §3.4 wurde von „sprach- und
  meilensteinfrei" auf „meilensteinfrei" umbenannt — die Architektur
  referenziert sprachkonkrete Modul-Pfade (per [`ADR-0001`](../../adr/0001-implementierungssprache.md)/0004 fixiert);
  das alte Wording wäre eine Harness-Lüge gewesen. Setext-Headings
  sind in der Slug-Spezifikation bewusst als nicht unterstützt
  deklariert (Quell-Tools nutzten nur ATX; fortschreibbar).
- **Steering-Loop-Lerneintrag (aus slice-001 übernommen, geprüft):**
  Die `Schärft:`-Felder der ADRs bleiben unverändert — ADRs sind nach
  `Accepted` immutable (`AGENTS.md` §3.5); die Formulierung „entstehen
  mit slice-002" ist als historische Aussage weiterhin korrekt.
  Künftige ADRs können auf konkrete Abschnitte der nun existierenden
  Straten zeigen.
- **Folge-Slices:** keine neuen. welle-01-fundament ist mit diesem
  Slice abgeschlossen (Closure-Trigger erfüllt); welle-02-mvp
  (slice-003, slice-004) ist entsperrt.

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (reine Doku-Arbeit; siehe Kurs Modul 5
§Worked Mini-Example).
