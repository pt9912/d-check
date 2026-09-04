# Ein Review-Report ohne `slice-<NNN>` im Dateinamen ist für beide Archivierungs-Modi strukturell unsichtbar

**Sub-Area:** `*`

`tools/archive-wave`s Sammel-Logik (beide Modi, Wellen- und Einzel-Slice)
findet Review-Reports ausschließlich über ein `slice-<NNN>`-Muster im
Dateinamen. Ein Report ohne dieses Muster — reine CR-/Baseline-/Release-Prep-
Reviews — bleibt dauerhaft unarchiviert, unabhängig davon, wie viele Wellen
oder Einzel-Slices archiviert werden. Das ist keine einmalige Lücke, sondern
eine strukturelle Eigenschaft der Sammel-Logik selbst.
