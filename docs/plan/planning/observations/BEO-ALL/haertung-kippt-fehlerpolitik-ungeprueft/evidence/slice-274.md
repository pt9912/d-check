**Vorgang:** slice-274
**Fund:** Ein neuer Leseweg — die Strikt-Prüfung der Einträge von
`trace.requirements` über Neu-Serialisierung des Teilbaums — wurde nur gegen den
Fall gefahren, den er beheben sollte (unbekannter Schlüssel in der Liste).
Aliase und Merge-Keys, die Zeilennummern bei Kommentaren und die Objekt-Form
blieben ungeprüft, bis der Reviewer über genau den Code-Commit sie fuhr (Runde 1,
F-1, F-3). Der fremde Leser (Schritt 14) fing es, noch vor einem Release.
