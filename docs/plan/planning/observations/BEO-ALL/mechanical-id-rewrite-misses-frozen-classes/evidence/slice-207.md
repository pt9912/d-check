**Vorgang:** slice-207
**Fund:** Dritte Instanz, und diesmal ohne jede Ausnahme-Liste. Eine
Pin-Hebung ersetzte `.harness/baseline/<alt>/` durch `<neu>` über den ganzen
Baum außer dem vendorten Verzeichnis selbst. Gehoben wurden dabei auch die
`d-check:cite`-Direktiven **eingefrorener** `done/`-Slices — zehn Stück in
fünf Dateien —, deren Zitat gegen den **alten** Tag geschrieben war; die
Zeilen-Spannen wanderten nicht mit, und die Direktiven zeigten danach auf
Zeilen, an denen der zitierte Satz nicht mehr steht. Dazu zwei Prosa-Pfade in
denselben Slices und der Kern einer `Accepted`-ADR (den `adr-check` sofort
meldete). **Der Bestand trug die Regel die ganze Zeit sichtbar:** Zwei ältere
`done/`-Slices stehen seit der vorigen Hebung auf dem damaligen Tag — genau
deshalb. Gefangen hat es kein Gate, sondern der unabhängige Review;
`citations` sieht diese Verzeichnisse nicht, weil `citations.scope` sie
ausnimmt.
