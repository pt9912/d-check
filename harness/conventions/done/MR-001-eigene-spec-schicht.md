# MR-001 — Source Precedence mit eigener Spezifikations-Schicht

- **Status:** Accepted
- **Datum:** 2026-06-10
- **Geltungsbereich:** [`harness/README.md` §Source precedence](../../README.md#source-precedence)
- **Adaption:** Die Source-Precedence-Tabelle führt
  `spec/spezifikation.md` als eigenen **Rang 2** zwischen Lastenheft
  (Rang 1) und Architektur (Rang 3). Der Kurs-Default setzt zwei
  Spec-Ränge; dieses Repo nutzt drei. Beide Dateien der Ränge 2–3
  entstehen mit slice-002; bis dahin sind sie in den Tabellen als
  „geplant" markiert und nicht verlinkt.
- **Begründung:** Spec-Stratifizierung mit drei Spec-Dateien; die
  ADR-Schärfungs-Regel („ADR darf Spezifikation schärfen, nicht
  Lastenheft") soll strukturell sichtbar sein.
- **Auflösungs-Trigger:** permanent.
