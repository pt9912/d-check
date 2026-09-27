**Vorgang:** slice-232
**Fund:** Die Erstfassung der Entscheidung zur Zeilenumbruch-Erkennung
(mittlerweile Superseded) faltete einen ganzen Markdown-Absatz zu einem
String und ließ `matchBracket` (zeilenblind) sowohl über die
Linktext-Klammer `[…]` als auch die Adress-Klammer `(…)` darüber laufen. Der
unabhängige Review (R1-H1, HIGH) konstruierte ein unbalanciertes `[` in
gewöhnlicher Prosa, das mit einer späteren, unabhängigen `](…)`-Sequenz im
selben Absatz zu einem erfundenen Link verschmolz. Die Korrektur
([ADR-0092](../../../../../adr/0092-links-zeilenumbruch-begrenzter-lookahead.md))
begrenzt den Lookahead auf eine Zeile und ausschließlich die Klammer, die
den Change Request tatsächlich trägt (die Adresse) — die Linktext-Klammer
bleibt strikt zeilenlokal, unabhängig vom Absatz-Inhalt.
