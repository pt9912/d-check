# Eine Deklaration entlastet eine Anforderung, ohne sie zu belegen

**Sub-Area:** `*`

Eine Kennung im Doc-Kommentar eines Tests macht eine Anforderung in der RTM
waisenfrei (`trace.coverage`, [ADR-0104](../../../../adr/0104-test-nachweise-entlasten-in-der-rtm.md)) — auch wenn der Test sie gar nicht
prüft, etwa weil er die Konfiguration des Repos statt des Produkts testet.
Umgekehrt zählt ein Test, der die Anforderung prüft, aber die Kennung nur im
Datei-Kommentar trägt, nicht. Ableiter: bei jeder neuen Kennung im
Doc-Kommentar eines Tests fragen, ob der Test die Anforderung über das
Produkt prüft.
