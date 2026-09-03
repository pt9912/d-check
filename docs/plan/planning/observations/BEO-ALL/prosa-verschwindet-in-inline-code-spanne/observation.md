# Prosa verschwindet in einer absatzweiten Inline-Code-Spanne, und kein Gate sieht es

**Sub-Area:** `*`

Die geteilte Lexik paart Backticks je Absatz, nicht je Zeile; zwei Backticks
in verschiedenen Zeilen desselben Absatzes ziehen eine Spanne über alles
dazwischen und leeren es positionserhaltend für jedes prosa-lesende Modul.
Als Markdown-Lesart korrekt — der Schaden entsteht, wenn die Spanne
unbeabsichtigt ist. `spans` deckt es nicht: sein `span-unclosed` meldet
ungerade Parität, dieser Fall hat gerade. Die Klasse ist breiter als
gedacht: derselbe stille Ausfall entsteht auch über einen vergessenen
Schluss-Fence.
