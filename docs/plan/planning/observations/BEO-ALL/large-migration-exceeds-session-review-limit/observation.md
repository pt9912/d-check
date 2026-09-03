# Ein Formatmigrations-Slice sprengt seinen eigenen Größen-Trigger und wird nicht rückgeführt

**Sub-Area:** `*`

Der Umfang einer einmaligen Datenformat-Migration lässt sich vor der
Ausführung unterschätzen (geschätzt ~27 Einträge, real 28 plus rund 180
geänderte Dateien) und übersteigt dann die Ein-Sitzungs-Review-Grenze aus
Modul 5. Der Slice-Plan hatte für genau diesen Fall vorab einen Trigger
(`in-progress` → `next`, Aufteilung nach Sub-Area) benannt — der Trigger
griff, weil eine Aufteilung den Zähler-Diff-Beleg zwischen alter und neuer
Form zerrissen hätte, ohne den die stille Bedeutungsverschiebung nicht
ausgeschlossen werden kann. Diese Spannung — „nicht mitten in der
Zähler-Verifikation teilen" gegen „einen Slice nicht über die
Review-Sitzungs-Grenze wachsen lassen" — hat keine etablierte Auflösung.
