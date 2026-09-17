# Ein korrekter Fix ändert nebenbei einen zweiten, nicht in Root-Cause/DoD benannten Aufruf-Pfad mit

**Sub-Area:** `*`

Ein Fix wirkt auf der richtigen Schicht (anders als bei
[`fix-schliesst-pfad-nicht-klasse`](../fix-schliesst-pfad-nicht-klasse/observation.md),
wo die Reparatur zu flach ansetzt und ein zweiter Weg zur selben Ursache
offen bleibt) — hier schließt der Fix die ganze Klasse bereits korrekt,
weil er nicht pfad-spezifisch, sondern an einer gemeinsamen Funktion
ansetzt. Genau das lässt einen zweiten, real existierenden Aufruf-Pfad
derselben Funktion **unbemerkt mit durchlaufen**: Die Root-Cause-Analyse
und die Tests wurden nur gegen den im Anlass genannten Pfad (`--range`)
geschrieben, ein struktogleicher zweiter Pfad (`--staged`) teilt sich
denselben Code und wechselt sein Verhalten ebenso — ohne dass DoD, Plan
oder ein Test das je benennt. Der Nettoeffekt kann sogar korrekt und
erwünscht sein; das Risiko liegt darin, dass die einzige Stelle, an der das
Verhalten je beobachtet wurde, mit dem Review-Report verschwindet, wenn sie
nicht auch als Test verankert wird.
