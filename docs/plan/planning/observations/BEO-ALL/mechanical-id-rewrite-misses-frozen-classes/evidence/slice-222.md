**Vorgang:** slice-222
**Fund:** **Die Frozen-Liste wurde vor der Ersetzung erstellt — und fiel beim
ersten Versuch durch.** [`MR-070`](../../../../../../../harness/conventions.md#mr-070) verlangt die Einteilung über die
**Eigenschaft**, nicht über Verzeichnisse. Der erste Anlauf tat trotzdem
genau das: `case`-Muster auf Pfad-Präfixe. Ergebnis — **alle 108 Dateien**
landeten in der Klasse „LEBEND", auch der vendorte Baum und die
`done/`-Slices.

**Der Fehler war stumm.** Die Ausgabe sah plausibel aus: eine lange, sortierte
Liste mit Zahlen. Nur wer sie **liest**, sieht, dass 28 Regelwerk-Dateien als
„lebend" geführt sind. Ein Lauf, der die Zahl genommen und weitergemacht
hätte, hätte den vendorten Baum retargetet und die eingefrorenen Belege
mitgehoben.

**Nach Eigenschaft ging es**: vendorter Baum (wird ersetzt) · Lauf-Beleg
(`done/`-Slice, Review-Report) · Versions-Historie · Register-Beleg · lebend.
Fünf Klassen, 108 Dateien, 299 Vorkommen — und die Zahl, die den Rest des
Slice trug, war die Aufteilung **77 lebende Pfad-Verweise gegen 33
eingefrorene**.

**Der Eintrag ist damit zum dritten Mal an derselben Stelle belegt** und seine
Verkörperung als [`MR-070`](../../../../../../../harness/conventions.md#mr-070) bestätigt: Die Regel stand, wurde befolgt — und die
*Ausführung* der Einteilung ist trotzdem der Ort, an dem es schiefgeht.
