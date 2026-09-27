**Vorgang:** slice-237
**Fund:** Beim Einführen von `structure[].max-lines`
([ADR-0089](../../../../../adr/0089-structure-max-lines-zwoelfte-bedingung.md))
als Erweiterung neben dem bereits bestehenden `file[].max-lines`
([ADR-0088](../../../../../adr/0088-file-modul-groessengrenzen.md), slice-236)
wurde die Namens-Kollision erkannt und im Vertrag beidseitig abgegrenzt
(Spezifikation §2-Schema-Zeilen,
[ADR-0089](../../../../../adr/0089-structure-max-lines-zwoelfte-bedingung.md)
§Konsequenzen, Re-Evaluierungs-Trigger 2). Keine reale Verwechslung ist bislang
aufgetreten; das strukturelle Risiko bleibt, solange beide Schlüssel den
Namen teilen.
