**Vorgang:** slice-222
**Fund:** **Der Eintrag hat diesmal gewirkt, bevor der Fehler entstand — und
das ist sein zweiter Beleg, nicht sein erster Rückfall.**

Die Lage war dieselbe wie bei seiner Anlage: Eine Form-Frage stand an
(`AGENTS.md` §4 — Index-Zellen kürzen oder Tabelle streichen?), und der
Bestand hätte eine plausible Antwort gegeben. [slice-221](../../../../next/slice-221-agents-md-tabellenzellen.md)
hatte sie schon abgeleitet — kürzere Zellen plus `cell-max-chars`-Schwelle,
grün und bruchgetestet.

**Die Vorlage sagte etwas anderes:** `AGENTS.template.md` in `v6.6.0` führt die
Tabelle **gar nicht mehr** — *„Der Gate-Index steht einmal, in
`harness/README.md` §Sensors … Diese Datei führt die Liste nicht."* Die
Bestands-Antwort war also nicht bloß ungenauer, sondern die **falsche
Operation**.

**Ausgelöst hat die Prüfung nicht der Eintrag, sondern der Auftraggeber**
(*„in v6.6.0 gibt es weniger Tabellen in AGENTS.md"*). Das ist die ehrliche
Zurechnung: Der Eintrag stand seit einem Tag da, und der Lauf hat trotzdem
weitergearbeitet, ohne die Vorlage zu lesen — sie war zu dem Zeitpunkt auch
noch nicht vendored. **Genau darin liegt die Verschärfung:** Wo die Vorlage
erst nach einem Bump lesbar wird, ist „an der Vorlage nachschlagen" keine
Handlung, die der Lauf jederzeit ausführen kann — sie hat eine
**Reihenfolge**, und die gehört in den Zuschnitt.

**Folge:** slice-221 ging nach `next/` zurück, statt eine Form zu zementieren,
die die Vorlage abschafft.
