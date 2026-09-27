# Zwei Regelmodule tragen einen gleich benannten Schlüssel `max-lines` mit unterschiedlicher Semantik

**Sub-Area:** `*`

`structure[].max-lines` ([ADR-0089](../../../../adr/0089-structure-max-lines-zwoelfte-bedingung.md))
zählt die Zeilenumbrüche des **bereinigten** Abschnittstextes einer
Markdown-Datei; `file[].max-lines`
([ADR-0088](../../../../adr/0088-file-modul-groessengrenzen.md)) zählt die
**rohen** Zeilen einer **ganzen** Datei jeder Art. Beide leben unter
verschiedenen Top-Level-Blöcken und kollidieren syntaktisch nicht, aber der
gleiche Schlüsselname bei unterschiedlicher Grundmenge ist ein
Verwechslungsrisiko in Konfiguration, Support und Doku. Die Abgrenzung steht
im Vertrag (Spezifikation, beide Schema-Zeilen;
[ADR-0089](../../../../adr/0089-structure-max-lines-zwoelfte-bedingung.md)
§Konsequenzen), ist aber keine Sensor-Zusage — kein Gate verhindert eine Fehlkonfiguration, die
das falsche Modul für den gemeinten Zweck wählt.
