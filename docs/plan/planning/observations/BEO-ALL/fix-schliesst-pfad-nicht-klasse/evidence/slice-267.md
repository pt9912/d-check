**Vorgang:** slice-267
**Fund:** Der Schlüssel `vcs.ignore-link-targets` leerte Link-Ziele über eigene
Muster. Sechs Review-Runden fanden je neue Markdown-Formen, in denen
Inhaltstext als Ziel geleert wurde — Fußnote, Prosa in Referenz-Form, `](`
ohne Link (R1), Code und Escapes (R2, R3), Code-Span im Linktext als Gegenfehler
(R4), Klammertext, Zitate, Container-Code (R5, R6). Jeder Fix schloss die
gemeldete Form. Geschlossen wurde die Klasse erst eine Schicht tiefer: ein Ziel
wird nur normiert, wenn es als Datei auflöst — Inhaltstext löst nicht auf.
