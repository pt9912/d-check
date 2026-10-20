**Vorgang:** slice-277
**Fund:** Im MR-Eintrag standen drei Zahlen, die nicht dem Gegenstand entsprachen:
„38 aktive MR-*" (35), „19 lebende d-check:cite-Direktiven" (15) und die daraus folgende
„übrigen 18" (14). `git grep -h … | wc -l` zählte **Zeilen** über einen Pfad-Ausschluss,
der `conventions/done/` (eingefroren, nicht im `citations`-Scope) und den Plan des Slice
mitnahm; die Summe der Dateien (49) trug die 38 nicht. Reviewer und Verifier zählten am
Vorzustand nach und fanden es; die Aussage über den Gegenstand (genau eine Spanne
verschoben) stimmte.
