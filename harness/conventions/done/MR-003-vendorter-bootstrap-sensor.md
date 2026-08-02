# MR-003 — Vendorter Bootstrap-Sensor `tools/verify-doc-refs.sh`

- **Status:** Accepted
- **Aufgelöst durch:** MR-007 (doc-check als Dogfooding)
- **Datum:** 2026-06-10
- **Geltungsbereich:** `tools/verify-doc-refs.sh` (gelöscht mit slice-004, siehe [`MR-007`](../../conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)), `make doc-check` <!-- d-check:ignore (historisch: gelöscht) -->
- **Adaption:** Bis `d-check` sich selbst prüfen kann, läuft
  `make doc-check` über ein aus `d-migrate` vendortes Shell-Skript
  (Markdown-Linkziel-Prüfung). Das ist Fremd-Code ohne eigene Spec in
  diesem Repo (Sub-Area in BF, siehe Modus-Tabelle).
- **Begründung:** Ein Doku-Repo ohne Doku-Sensor wäre ein blinder
  Bootstrap; das vendorte Skript ist dependency-frei (bash/awk) und
  deckt den Kern von
  [`DC-FA-LINK-001`](../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
  ab.
- **Auflösungs-Trigger:** slice-004 — `make doc-check` läuft über
  `d-check` selbst (Dogfooding), das Skript wird gelöscht
  ([Slice-Plan](../../../docs/plan/planning/done/slice-004-anchors-modul-und-dogfooding.md)).
