**Vorgang:** slice-206
**Fund:** Drei Tests eines neuen Moduls konnten nicht fallen. Der
Determinismus-Test verglich fünf Läufe gegen eine feste Reihenfolge — die
Fixture konnte sie nicht verletzen, weil das Speicher-Dateisystem selbst
sortiert; mit entfernter Sortierung blieb er grün. Der einzige Test der
`scan.ignore`-Auswertung benutzte ein Verzeichnis, das schon die feste
Skip-Liste prunt; mit vollständig entfernter Auswertung blieb er grün. Und die
Runen-Basiertheit der rechten Grenze hielt gar kein Test, obwohl ihre
Einführung als Grund für einen Fassungswechsel genannt war. Alle drei fielen
erst auf, nachdem ein Reviewer sie am Produktivcode mutierte — die Zusage im
Testnamen ist von einem echten Wächter nicht zu unterscheiden.
