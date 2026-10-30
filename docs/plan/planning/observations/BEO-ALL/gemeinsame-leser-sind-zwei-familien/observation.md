# „Dieselben Leser wie die übrigen Module“ stimmt nur für eine von zwei Familien

**Sub-Area:** `internal/hexagon/core/rules`

Eine Spezifikation, die Zäune und Tabellenzeilen „wie die übrigen Module“ erkennen lässt, trifft zwei Automaten: die einen schalten bei jeder Fence-Zeile um (`links`, `anchors`, `ids`), die anderen schließen längengleich mit dem Schluss-Test (RTM-Tabellen-Leser, `spans`). Der gemeinsame Öffner-Test verdeckt, dass das Schließen nicht gemeinsam ist.
