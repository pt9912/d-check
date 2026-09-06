# MR-007 — Auflösung von MR-003: doc-check als Dogfooding

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** [`modul-13-quality-gates.md` §Hard Rule (Doku-Disziplin)](../../.harness/baseline/v6.3.1/regelwerk/modul-13-quality-gates.md#hard-rule-doku-disziplin)
- **Datum:** 2026-06-10
- **Geltungsbereich:** `make doc-check`, [`.d-check.yml`](../../.d-check.yml)
- **Adaption:** Der Auflösungs-Trigger von
  [`MR-003`](../conventions.md#mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh)
  ist eingetreten: `make doc-check` läuft über `d-check` selbst
  (Runtime-Image, read-only-Mount; Module `links` + `anchors` über
  die gesamte Repo-Wurzel via `scan.roots: ["."]`). Das vendorte
  Skript `tools/verify-doc-refs.sh` ist gelöscht; der <!-- d-check:ignore (historisch: gelöscht) -->
  Geltungsbereich-Link in MR-003 wurde dafür auf einen Code-Span
  umgestellt (Form-, keine Inhaltsänderung). Vergleichslauf
  (erster Datenpunkt für
  [`DC-QA-04`](../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)):
  Alt-Skript 0 broken links, `d-check` 0 Befunde bei 23 Dateien —
  bei strikt größerer Abdeckung (zusätzlich Anker-Validierung und
  Bildreferenzen).
- **Begründung:** Dogfooding-Ziel von slice-004; die BF-Sub-Area aus
  der Modus-Tabelle ist damit graduiert (gelöscht).
- **Auflösungs-Trigger:** permanent (Dogfooding ist der Zielzustand).
