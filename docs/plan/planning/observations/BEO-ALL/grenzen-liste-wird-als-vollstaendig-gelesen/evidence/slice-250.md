**Vorgang:** slice-250
**Fund:** Die Grenzen der neuen Glob-Expansion wurden gegen die eigene
Beschreibung geschrieben, nicht gegen den Code. Lastenheft, Spezifikation,
ADR und Commit-Botschaft sagten „keine Symlinks" und übersprungene
Verzeichnisse „wie beim Modul `file`" zu. Der Code prüfte per `Kind` nur die
letzte Präfix-Komponente (ein Symlink davor wurde verfolgt, auch aus der
Wurzel hinaus), überging einen passenden Symlink darunter still, und
behandelte ausdrücklich genannte übersprungene Namen anders als `file`.
Keine der Grenzen trug einen Negativtest. Gefunden vom unabhängigen Review
(R1-M1, M2, M3, L2) per Black-Box gegen den echten Dateisystem-Adapter, obwohl
die Regel (`AGENTS.md` §5 Regel 13) stand.
