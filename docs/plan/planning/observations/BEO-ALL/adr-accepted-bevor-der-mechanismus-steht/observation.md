# Eine ADR wird `Accepted`, bevor der Mechanismus steht, den sie beschreibt

**Sub-Area:** `*`

Eine ADR, die mit ihrem ersten Commit `Accepted` ist, friert den Wortlaut des
ersten Entwurfs ein. Ändert das Review danach den Mechanismus, kann ihr Körper
nicht folgen; jede Präzisierung landet als Anhang in `## Geschichte`, und der
Körper beschreibt einen Stand, den es nie gab. Ableiter: eine ADR, deren
Mechanismus ein Review noch ändern kann, bleibt `Proposed` bis zur Closure.
