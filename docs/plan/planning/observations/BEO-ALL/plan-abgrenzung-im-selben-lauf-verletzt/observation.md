# Der Implementer weitet die eigene, im Plan geschriebene Abgrenzung im selben Lauf aus

**Sub-Area:** `*`

`AGENTS.md` §6 Schritt 4 verlangt: „Er darf die Abgrenzung nicht ausweiten,
weder still noch begründet. Wer im Lauf etwas mitnimmt, das der Plan
ausschließt, hat den Plan geändert — und das gehört vor den Code, nicht in
den Bericht danach." Der Implementer-Kontext, der den Plan geschrieben hat,
ist derselbe, der ihn im Lauf verletzt — er sieht den eigenen Bruch nicht,
weil er die eigene Abgrenzung nicht als Prüf-Checkliste gegen den entstehenden
Diff hält, sondern nur beim Schreiben des Plans einmal bedacht hat.

Ableiter, ein Schritt: **vor** Schritt 8 (Pre-completion) den eigenen Diff
gegen §1 „Ausdrücklich NICHT" durchgehen, Datei für Datei — nicht aus dem
Gedächtnis, sondern durch Lesen der Abgrenzungs-Liste gegen `git diff
--stat`.
