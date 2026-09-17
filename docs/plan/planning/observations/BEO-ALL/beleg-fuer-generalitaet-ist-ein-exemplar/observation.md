# Eine Aussage über "jeden Fall einer Klasse" ist an einem einzigen Exemplar geeicht

**Sub-Area:** `*`

Ein Fix wird als Lösung für eine ganze Klasse formuliert („jedes Präfix",
„jede Fremdbibliothek dieser Art"), belegt ist er aber nur an dem einen
Exemplar, das den Anlass gab — hier: `git maintenance
run --task=loose-objects` mit dem Präfix `loose-`. Anders als
[`BEO-ALL/rule-drawn-from-occasion-not-inventory`](../rule-drawn-from-occasion-not-inventory/observation.md)
ist das hier **keine fehlerhafte** Verallgemeinerung — die Mechanik
(Hash-Suffix + passender Index) ist präfix-unabhängig konstruiert und deckt
jedes Exemplar der Klasse strukturell ab, nicht nur das eine, das getestet
wurde. Der Unterschied ist die **Evidenzbasis**: Ein struktureller
Mechanismus, der nur an einem Exemplar beobachtet wurde, ist etwas anderes
als einer, der an mehreren beobachtet wurde — auch wenn beide korrekt sein
können. [ADR-0072](../../../../adr/0072-workflows-modul.md) benennt exakt
dieselbe Lücke für die yaml-Kapsel-Erweiterung und akzeptiert sie
ausdrücklich als Grenze, nicht als Mangel.
