**Vorgang:** slice-279
**Fund:** Die Default-Erkennung des Trace macht jede Datei `slice-<name>.md` zum
Slice, ohne zu prüfen, ob sie ein Slice-Plan ist; eine Notiz mit `slice-`-Präfix
im Planning-Baum deckte eine Anforderung. Am Dateinamen nicht entscheidbar — die
Grenze steht in [ADR-0110](../../../../../adr/0110-trace-default-erkennt-namens-slices.md), die Abhilfe ist ein gesetztes
`trace.slices.file-pattern`. Eigenbestand: kein solcher Fall (549 Dateien, 279
Ziffern-Slices, 270 Belege unter `observations/`). Dieselbe Klasse wie der
Pack-Name als Beleg für den Hash, hier mit der Sub-Area `*` statt des
Git-Adapters.
