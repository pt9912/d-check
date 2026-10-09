# Ein leeres Verzeichnis bleibt lokal stehen und verdeckt einen Befund, den die CI sieht

**Sub-Area:** `*`

`git mv` entfernt die Datei, nicht das Verzeichnis, das sie enthielt. Lokal
bleibt es leer stehen; ein frischer Klon kennt es nicht, weil git keine
leeren Verzeichnisse führt. Eine Prüfung im lokalen Arbeitsbaum, die einen
Pfad auf dieses Verzeichnis auflöst, ist grün, wo die CI rot ist. Ableiter:
nach einem Umzug das leere Verzeichnis entfernen und im frischen Klon messen.
