**Vorgang:** slice-218
**Fund:** **Dieselbe Aussage an vier Orten, und jede Review-Runde fand einen,
der hinterherhing.** Der Slice schrieb eine Grenze in zwei Sensor-Dateien, das
Benutzerhandbuch und später die Spezifikation. §6 hatte die Drift-Gefahr
**vorab als Risiko benannt** — und sie trat trotzdem dreimal ein:

- Runde 1: die fail-closed-Zusage war an allen Orten falsch.
- Runde 2: alle Orte sagten den Abschluss über die ganze Klasse zu, während
  nur eine Hälfte gemessen war.
- Runde 3: alle Orte führten für **beide** Ausprägungen dieselbe Ursache an,
  obwohl für die zweite ein anderer Mechanismus wirkt; und drei Orte
  beschränkten `CO-001` weiter auf den `RANGE=`-Modus, obwohl der Carveout
  längst beide führte.

**Die Gegenmaßnahme ist Struktur, nicht Sorgfalt** — das ist der Ertrag: Die
vollständige Messung steht jetzt an **einem** Ort
(`harness/sensors/adr-check.md`), die übrigen verweisen darauf, statt sie zu
wiederholen. **Ein Risiko vorab zu benennen hat hier nichts verhindert**; es
hat nur erklärt, was passierte.
