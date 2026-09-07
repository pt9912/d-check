# Eine Review-Korrektur wird nur an der zitierten Fundstelle eingearbeitet

**Sub-Area:** `*`

Das `pfad`-Feld eines Findings sagt, **wo der Befund gefunden wurde** — nicht,
wo er **steht**. Wer es als Umfang liest, korrigiert eine Fundstelle und lässt
ihre Geschwister stehen; der nächste Review findet dieselbe Aussage an der
nächsten Stelle und liest sie zu Recht als neuen Befund. Das Feld ist dabei
kein Fehler des Reports: Ein Reviewer belegt, er inventarisiert nicht.
Unterschied zu
[`BEO-ALL/semantic-change-body-only-edges-stale`](../semantic-change-body-only-edges-stale/observation.md):
dort ist der Auslöser die **eigene** Änderung und die Spiegel sind
Rand-Artefakte; hier ist der Auslöser ein **fremder Befund**, und die
Geschwister sind gleichrangige Vorkommen derselben Aussage. Ableiter: vor dem
Abhaken eines Findings die **Formulierung** des Befunds repo-weit suchen, nicht
den genannten Pfad öffnen.
