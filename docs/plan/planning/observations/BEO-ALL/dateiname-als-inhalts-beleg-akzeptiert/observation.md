# Ein Dateiname wird als Beleg für seinen Inhalt akzeptiert, ohne den Inhalt zu prüfen

**Sub-Area:** `internal/adapter/driven/git`

Eine Erkennungsregel liest eine Eigenschaft (hier: den Content-Hash eines
Packs) aus dem **Dateinamen**, statt sie am tatsächlichen Inhalt zu
verifizieren. Solange die Datei von vertrauenswürdiger Seite entsteht
(`git`/`git maintenance` benennt einen Pack korrekt nach seinem echten
Hash), trifft Name und Inhalt zusammen — die Regel ist dann eine Heuristik,
kein Beweis. Ein Name, der zufällig oder absichtlich einen validen
Hash-Suffix trägt, ohne der echte Content-Hash zu sein, würde fälschlich
akzeptiert. Ob eine nachgelagerte Schicht (hier: go-gits eigener
Checksummen-Vergleich beim Öffnen der Idx-Datei) den Fehlgriff auffängt, ist
eine **Folge dieser Schicht**, keine Garantie der eigenen Regel — und bleibt
deshalb ein eigenes Risiko, auch wenn der praktische Schaden meist an dieser
zweiten Schicht hängen bleibt.
