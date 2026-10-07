**Vorgang:** slice-252
**Fund:** Die R1-Behebung (nur die erste Zelle führt ein Target) setzte im
geteilten Extraktor der Doku-Tabellen an — gemeinsam für die neue
Disjunktheit **und** für Richtung 1 und 2. Der Zugriff `tableCells(line)[0]`
ohne Längenprüfung ließ eine Tabellenzeile, die nur aus `|` besteht, in
**jedem** `targets`-Lauf mit Panic abbrechen, auch ohne Schalter; der Fix
war für die Disjunktheit gedacht und änderte unbenannt den Pfad der beiden
Bestands-Richtungen mit. Gefunden vom Verifier per Black-Box (V1, HIGH).
