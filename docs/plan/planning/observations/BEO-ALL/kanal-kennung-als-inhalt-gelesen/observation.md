# Eine Kennung, die einen Kanal ausweist, wird als Bezug auf den Inhalt gelesen

**Sub-Area:** `*`

Die Dependabot-Commits tragen `[ADR-0067]` im Präfix, <!-- d-check:ignore (literaler Praefix-Wert, kein Verweis) --> weil das
Traceability-Gate von jeder Botschaft eine Kennung verlangt. Die Kennung
sagt „dieser Commit entsteht aus dem Kanal, den [ADR-0067](../../../../adr/0067-dependabot-als-hebender-kanal.md) erlaubt" — sie
sagt **nicht**, dass die gehobene Bibliothek etwas mit jener Entscheidung zu
tun hat. Die Klasse ist allgemeiner als der Anlass: überall, wo eine ID ein
Verfahren statt eines Gegenstands ausweist, lädt sie zur Fehl-Lesung ein.

## Benannt, nicht gezählt

Noch nicht eingetreten, als Wachposten geführt (Zähler 0) — Eintritts-Bedingung:
ein Commit, ein Review-Befund oder eine Doku-Zeile begründet etwas mit
`[ADR-0067]`, das über „stammt aus dem Dependabot-Kanal" hinausgeht. <!-- d-check:ignore (literaler Praefix-Wert, kein Verweis) --> Zu
streichen, sobald die Dependabot-Commits einen anderen Traceability-Weg
haben.
