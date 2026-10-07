# Die Plan-Änderung wird erst im Fix-Commit nachgetragen, nicht vor dem Code

**Sub-Area:** `*`

`AGENTS.md` §6 Schritt 4: wer im Lauf etwas mitnimmt, das der Plan nicht
führt, hat den Plan geändert — „und das gehört vor den Code, nicht in den
Bericht danach". Beobachtet wird die schwächere Form des Bruchs: die
Mitnahme ist nicht ausgeschlossen, nur ungeplant, und der Vermerk im Plan
entsteht erst, wenn ein Review die Mitnahme sichtbar macht — im selben
Commit wie die Korrektur oder später, nicht davor. Unterschied zu
[`BEO-ALL/plan-abgrenzung-im-selben-lauf-verletzt`](../plan-abgrenzung-im-selben-lauf-verletzt/observation.md):
dort wird eine ausdrückliche Abgrenzung verletzt, hier fehlt der Plan-Eintrag
für eine nicht ausgeschlossene Mitnahme. Ableiter: eine Datei, die der Plan
in §3 nicht nennt, löst den Vermerk aus, bevor sie committet wird.
