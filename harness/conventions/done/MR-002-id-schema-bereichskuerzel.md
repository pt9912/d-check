# MR-002 — ID-Schema mit Bereichskürzeln ab initialer Fassung

- **Status:** Accepted
- **Datum:** 2026-06-10
- **Geltungsbereich:** [`spec/lastenheft.md`](../../../spec/lastenheft.md), alle Traceability-Verweise
- **Adaption:** Funktionale Anforderungen verwenden von Beginn an
  Bereichskürzel: `DC-FA-<BEREICH>-<NNN>` (z. B. [`DC-FA-LINK-001`](../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links))
  statt des zweistelligen Kurs-Defaults `<PREFIX>-FA-<NN>`.
  Nichtfunktionale Anforderungen bleiben beim Kurs-Default
  (`DC-QA-<NN>`).
- **Begründung:** Das Lastenheft konsolidiert zwölf Quell-Tools und hat
  dadurch von Anfang an viele Funktionsbereiche; das Kurs-Beispiel
  (DocSearch) zeigt, dass eine spätere Schema-Migration teurer ist als
  ein Bereichsschema ab Welle 1.
- **Auflösungs-Trigger:** permanent.
