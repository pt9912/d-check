# Ein Negativtest für eine symmetrische Regel bricht immer dieselbe Seite

**Sub-Area:** `*`

Eine Validierungsfunktion prüft zwei (oder mehr) gleichrangige Seiten
derselben Regel in einer Schleife und kehrt beim ersten Fehler zurück. Ein
Negativtest mit mehreren Unterfällen variiert dabei unbeabsichtigt nur die
**erste** Seite — die Schleife verlässt jeden Unterfall, bevor die zweite
Seite je geprüft wird. Der Code ist korrekt, aber ein im Anforderungstext
wörtlich genannter Fall (die zweite Seite) hat keinen dedizierten Beleg.
