**Vorgang:** slice-233
**Fund:** Zweifacher Auftritt im selben Slice. R1-H1: die Erstfassung der
Definitions-Erkennung validierte den Titel-Rest nicht und las gewöhnliche
Prosa (`[TERM]: First In, First Out`) als Definition mit erfundenem Ziel.
Nach der Korrektur (Titel-Delimiter-Pflicht) fand R2-M1: ein whitespace-
freies, klammerartiges Ziel-Token (`[TODO]: (spaeter)`) wird weiterhin ohne
Wohlgeformtheits-Prüfung als Ziel akzeptiert — nicht blockierend (R2-Verdikt),
aber derselbe Mechanismus (Ziel-Token wird nicht validiert) eine Ebene
tiefer als der zuerst behobene Fund.
