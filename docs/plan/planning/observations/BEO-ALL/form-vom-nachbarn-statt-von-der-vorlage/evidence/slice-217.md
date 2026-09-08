**Vorgang:** slice-217
**Fund:** **Der Slice führte die Regel im Titel und verfehlte sie im eigenen
Text.** §1 lautete *„Die Form-Frage ist an der Vorlage geklärt, nicht am
Nachbarn"* — und schrieb der Baseline-Vorlage den Trenner `| --- |` zu.
Gemessen führt
[`harness/README.template.md`](../../../../../../../.harness/baseline/v6.5.0/templates/harness/README.template.md)
`|---|---|---|` **ohne Leerzeichen**. Der Trenner stammte von den drei übrigen
Tabellen derselben Datei — also vom Nachbarn.

**Die Vorlage klärte tatsächlich etwas, nur nicht das.** Sie beantwortet die
**Padding**-Frage (39 Tabellenzeilen, null Padding), und diese Antwort trägt
den Slice. Der **Trenner-Abstand** ist eine zweite Form-Frage, die im selben
Satz mitgenommen und stillschweigend derselben Quelle zugeschrieben wurde.

**Die Wahl selbst war richtig** und bleibt: Die verkörperte Form führt
gegenüber der Referenz-Form, und zwei Trenner-Stile in einer Datei hätten
gegen den Zweck des Slice gearbeitet. Falsch war ausschließlich die
**Herkunfts-Behauptung** — und genau die macht den Fehler unsichtbar, weil der
Satz die richtige Autorität nennt.

**Gefunden hat es der unabhängige Review**, nicht der Autor — obwohl der Autor
die Vorlage vorher selbst gemessen und `|---|---|---|` mit eigenen Augen in
der Ausgabe stehen hatte.
