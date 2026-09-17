# Ein Mechanismuswechsel zieht die Spec nach, aber eine Guide-Doku außerhalb der Spec nicht automatisch mit

**Sub-Area:** `*`

Ein Slice, der einen internen Mechanismus austauscht (nicht nur eine
Zusage erweitert), nennt im Kopf die berührten **Spec-Stellen** — und die
werden nachgezogen, geprüft, im Review verglichen. Eine **Guide-Doku** wie
`harness/sensors/*.md`, die denselben Sachverhalt in eigener Prosa
beschreibt (Fehlerbilder, Testnamen, den Zustand eines referenzierten
Carveouts), trägt aber **kein** Feld, das ihre Berührung erzwingt oder
sichtbar macht — sie ist kein „öffentlicher Vertrag" im Sinn der Spec-
Straten, sondern eine zusätzliche, freiwillige Erklärung desselben Vertrags
für Menschen.

**Gemessen:** `slice-220` (vcs-Kandidaten-Bestimmung, Diff → Mengen-
Differenz) nannte `spec/spezifikation.md` korrekt als berührte Spec-Stelle,
wurde dort nachgezogen und in **zwei** unabhängigen Review-Runden geprüft.
`harness/sensors/adr-check.md` und `trace-check.md` beschrieben denselben
Mechanismus in eigener Prosa — mit veralteten Testnamen
(`TestVCSAddedMeldetUnlesbareBasis`, seit slice-220 umbenannt), einer
überholten Drei-statt-Vier-Ausprägungen-Tabelle und `CO-001` als
„weiterhin offen" — und blieben unangetastet, bis die CO-001-Auflösung
(`slice-219`) sie zufällig traf, weil ihr eigener Doc-Check-Lauf die
gebrochene Link-Tiefe nach dem Carveout-Move meldete. Ohne den Move wäre
die Staleness vermutlich unentdeckt geblieben — kein Sensor liest
Guide-Prosa gegen Code.

**Der Ableiter ist eine Frage an den Implementer-Schritt 7
(„Doku/Indizes aktualisieren, falls ein öffentlicher Vertrag berührt"),
nicht an die Reviewer-Fragenliste:** Schritt 7 nennt „öffentlicher
Vertrag", ohne zu sagen, dass `harness/README.md`s Guide-Tabelle
(`harness/sensors/*.md`) denselben Rang trägt wie Spec/Handbuch/README.
Ein Implementer, der nur die im Slice-Kopf genannte Spec-Stelle prüft,
lässt die Guide-Doku unbemerkt zurück.
