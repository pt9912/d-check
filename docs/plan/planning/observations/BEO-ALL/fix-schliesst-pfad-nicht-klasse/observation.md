# Ein Fix schließt den gefundenen Pfad, nicht die Klasse

**Sub-Area:** `*`

Ein Defekt wird an der Stelle behoben, an der er **gefunden** wurde. Der Fix
ist korrekt, sein Bruch-Test grün, die Gegenrichtung geprüft — und trotzdem
ist der Defekt nicht weg, weil dieselbe Ursache über einen **zweiten Weg**
denselben Ausgang erzeugt. Die Reparatur wirkt auf der Schicht des Fundes,
nicht auf der Schicht der Ursache.

**Warum das anders ist als eine übersehene N+1-te Form:** Bei
[`commit-message-overclaims-work`](../commit-message-overclaims-work/observation.md)
reicht ein **Schluss** weiter als die Messung — dort hilft, mehr Formen zu
messen. Hier ist die **Reparatur** zu flach angesetzt, und mehr Proben auf
derselben Schicht finden den Rest nicht: Sie liegen davor. Wer nach einem Fix
nur die gefundene Form nachmisst, bestätigt genau die Grenze, die er gezogen
hat.

**Der Ableiter ist eine Frage an den Fix, nicht an die Proben:** *Auf welcher
Schicht wirkt er — und was liegt davor?* Beim Adapter-Fix lag die
Regel-Schicht davor, bei der Regel-Schicht die Diff-Erzeugung, und dort war
die dritte Form. Die Frage ist billig und wäre dreimal fällig gewesen.

**Gegenprobe, die die Klasse von bloßer Unvollständigkeit trennt:** Lässt sich
ein Wachposten **eine Schicht tiefer** setzen, der alle bekannten Formen auf
einmal abdeckt? Trägt er nicht — gemessen, nicht vermutet —, ist der flache
Fix legitim und seine Grenze gehört benannt.
