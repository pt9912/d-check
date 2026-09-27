# Eine Erkennungs-Regex für eine Markdown-Mikrosyntax bleibt gegen Freitext-Negativfälle ungeprüft, bis ein Review adversarial dagegen testet

**Sub-Area:** `*`

Ein Change Request nennt typischerweise nur **Positivfälle** (die Form, die
erkannt werden soll) und **einzelne** benannte Negativfälle. Eine daraus
abgeleitete Regex wird gegen genau diese Fälle getestet und ist grün — bleibt
aber gegenüber der viel größeren Menge an **Freitext, der zufällig ähnlich
aussieht**, ungeprüft, bis ein unabhängiger Review gezielt adversarielle
Gegenproben konstruiert. Dieselbe Grund-Ursache trat in slice-233 zweimal
auf: Erst ein naiver Whitespace-Schnitt las gewöhnliche Prosa (`[TERM]: First
In, First Out`) als Definition mit erfundenem Ziel (R1-H1); die Korrektur
validierte den Titel-Delimiter, ließ aber ein whitespace-freies,
klammerartiges Ziel-Token (`[TODO]: (spaeter)`) unvalidiert durch (R2-M1) —
derselbe Mechanismus (Ziel-Token-Wohlgeformtheit wird nicht geprüft), eine
Ebene tiefer.

## Benannt, nicht gezählt

- R1-L1/R2-L2 (`matrix`-Lineage-Match über das Definitions-Label
  undokumentiert) und R2-L1 (Backslash-escapte Anführungszeichen im Titel
  erzeugen ein Falsch-Negativ statt eines Falsch-Positivs) sind LOW-Funde
  derselben Review-Runden, aber eine andere Fehlerrichtung (fehlende Doku
  bzw. zu strenge statt zu lockere Erkennung) — hier benannt, nicht als
  eigener Beleg gezählt.
