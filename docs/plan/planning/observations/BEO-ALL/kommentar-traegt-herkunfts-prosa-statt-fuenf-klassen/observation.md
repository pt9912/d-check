# Ein neu geschriebener Kommentar trägt eine externe Vorgangs-Referenz statt einer der fünf zulässigen Klassen

**Sub-Area:** `*`

`AGENTS.md` §3.7 verlangt, dass ein Kommentar eine der fünf Klassen trägt
(Zusage · Kopplung · Abgrenzung · Rang-Zeiger · Grenze) und Herkunft nur als
**ein** auflösbares Feld nach dem Baseline-Schema (`DC-*`/`ADR-*`/`MR-*`/
`seit welle-<NN>`). Eine externe Vorgangs-Kennung — hier eine
GitHub-Issue-Nummer in Klammern, als Begründung, warum ein Test existiert —
ist dieselbe Klasse von Herkunfts-Prosa wie eine Slice- oder
Befund-Nummer, fällt aber leichter durch die eigene Aufmerksamkeit, weil sie
sich wie ein legitimer Beleg liest (die Nummer existiert ja tatsächlich,
extern nachprüfbar) statt wie ein offensichtlicher Zeitschnitt-Verweis. Die
Bestandsgrenze in §3.7 gilt ausdrücklich nicht rückwirkend für
**Neuzugänge** — ein frisch geschriebener Testkommentar ist davon nicht
ausgenommen, nur weil er in einer Testdatei statt in Produktionscode steht.
