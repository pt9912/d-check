# MR-0098 — Die Matrix nimmt die `7. Historie` beider Spec-Straten aus der Referenz-Richtung aus

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** [`grundlagen-referenz-richtung.md`](../../.harness/baseline/v6.17.0/regelwerk/grundlagen-referenz-richtung.md)
  §Referenz-Richtung (SDP), Regel 5 — „in keinem Abschnitt, auch nicht in
  seiner Historie". Die Ausnahme hebt genau diese Hälfte für die
  `7. Historie` beider Spec-Straten aus der Matrix-Prüfung heraus; das
  lebende Dokument bleibt voll geprüft. Begründung und Messung in
  [ADR-0097](../../docs/plan/adr/0097-matrix-aussen-adaptionsblock-historie-status-ausnahmen.md)
  (40 gemessene Befunde, davon 18 in den Historie-Zeiten).
- **Datum:** 2026-09-29 · **Herkunft:** seit slice-248 (Adoption aus
  ai-harness-init; die Schwester fährt dieselbe Ausnahme mit drei
  `exempt-paths`-Einträgen)
- **Geltungsbereich:** `matrix.exclude-sections` — der heading-genau
  Abschnitt `7. Historie` in `spec/lastenheft.md` und
  `spec/spezifikation.md`. Die ADR-Geschichte bleibt separat ausgenommen
  (`[Geschichte]`, [ADR-0047](../../docs/plan/adr/0047-matrix-spec-historie-nicht-provenance-exempt.md)-Form); die Referenz-Richtung des lebenden
  Dokumentkörpers ist von der Ausnahme nicht berührt.
- **Adaption:** die `7. Historie` ist ein Zeitdokument — ihre Zeilen
  frieren Provenance zum jeweiligen Datum (MR-Verweise, CR-Verweise,
  Agenten-Briefing-Referenzen). Die Referenz-Richtung gilt im lebenden
  Dokumentkörper unverändert; der frozen Bestand wird von der Matrix
  nicht auf Abwärtsverweise geprüft.
- **Begründung:** der Kanon lässt in der Lastenheft-Historie den externen
  Change Request ausdrücklich zu (v5.11.0) — die Verweis-Spalte ist
  canonisches Provenanz-Feld. Die aussen/adaptionsblock-Regeln würden
  18 frozen Verweise flaggen (gemessen in slice-246/248), ohne dass
  ein lebender Text sie verletzt; die Ausnahme hält die Regel am lebenden
  Text scharf, statt frozen Zeilen still durch die allgemeine Prüfung zu
  lassen.
- **Auflösungs-Trigger:** ändert die Baseline die Historie-Semantik
  (Referenz-Richtung gilt doch in frozen Abschnitten), wird der Eintrag
  aufgelöst und die Ausnahme gezogen.
