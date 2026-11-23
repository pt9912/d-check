# Der YAML-Dekoder kürzt eine Kommazahl still auf eine ganze Zahl

**Sub-Area:** `*`

Ein Zahl-Schlüssel der Konfiguration, der als Ganzzahl dekodiert wird, nimmt
`80.5` still als `80` an, und ein Wahrheitswert nimmt `yes` oder `on` als
`true`. Die Konfiguration ist dann gültig, obwohl sie nicht sagt, was sie meint.
Ein Modul, das seine Werte am YAML-Typ prüft, schließt die Lücke für sich; die
Nachbarmodule tragen sie weiter.
