**Vorgang:** slice-218
**Fund:** **Dreimal in einem einzigen Slice, jedes Mal vom unabhängigen Review
und nie vom Autor.** Derselbe Defekt — ein unlesbares git-Objekt lässt eine
echte `core-drift-vcs`-Verletzung verschwinden — trat in drei Formen auf, und
jeder Fix deckte genau die, an der er gefunden wurde:

| Form | Fix-Schicht | Ergebnis |
|---|---|---|
| BASE-**Blob** unsichtbar | Adapter (`FileAt`) | zu |
| BASE-**Tree** unsichtbar, mit Pendant | Regel (`A`-Zweig) | zu |
| BASE-**Tree** unsichtbar, **ohne** Pendant | — | offen |

**Jeder Fix sah beim Schreiben vollständig aus**, hatte einen Bruch-Test in
beide Richtungen und wurde mit einer Kontroll-Matrix belegt. Die jeweils
nächste Form lag **eine Schicht davor**: hinter dem Adapter die Regel, hinter
der Regel die Diff-Erzeugung, in der der Tree-Walker aus einem nicht ladbaren
Unterbaum ein `io.EOF` macht.

**Der Ableiter ist gemessen, nicht abgeleitet:** Nach der dritten Form wurde
der naheliegende Wachposten eine Schicht tiefer probiert (`tree.Files()`) —
er benutzt denselben Walker und schweigt ebenso. Genau diese Probe wäre nach
dem **ersten** Fix fällig gewesen und hätte die Schichtfrage sofort gestellt.

**Der Slice hat die Ein-Sitzungs-Review-Grenze deshalb zweimal überschritten**
([`MR-066`](../../../../../../../harness/conventions.md#mr-066)) und drei
Review-Runden gebraucht — die Kosten der flachen Reparatur sind an ihm
ablesbar.
