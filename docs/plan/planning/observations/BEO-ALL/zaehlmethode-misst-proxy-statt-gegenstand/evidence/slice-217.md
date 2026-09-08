**Vorgang:** slice-217
**Fund:** **Dreimal in einem Slice, dessen ganzer Gegenstand eine
Zeichen-Klasse ist.** Der Plan machte die Form des Gegenstands zu DoD (1) —
*ausschreiben, bevor gezählt wird* — und verfehlte die Messung trotzdem drei
Mal.

**Zweimal in der Voruntersuchung, beide vom Autor bemerkt:** Das Muster
`[^ ]  *|` traf auch **ein** Leerzeichen (`*` ist null oder mehr) und meldete
75 von 75 Tabellenzeilen statt 20. Das Trennzeilen-Muster `^\|[ -]+\|$`
vergaß die **inneren** Pipes im Zeichensatz und fand null Zeilen, wo zwei
stehen. Beide Male war die falsche Zahl **plausibel** — 75 sah nach „alle
Tabellenzeilen tragen Padding" aus, 0 nach „es gibt keine überlangen
Trennzeilen".

**Einmal im Ergebnis, vom Review gefunden:** Die Gesamtzahl 4737 zählte die
**ganzen** Leerzeichen-Läufe, die Teilzahlen 3394 und 1297 dieselben Läufe je
**um ein Zeichen gekürzt**. Beide richtig gemessen — aber zwei verschiedene
Größen, und die Summe ging um genau die Zahl der Läufe (46) nicht auf. Nach
der in §3 selbst deklarierten Klasse ist **4691** die richtige Gesamtzahl.

**Der Ableiter greift eine Stufe tiefer, als der Eintrag ihn bisher führt:**
Es genügt nicht, die Form vor der Messung auszuschreiben — sie muss **bis in
die Arithmetik** durchgehalten werden. Eine Gesamtzahl und ihre Teilzahlen
gegeneinander zu prüfen ist dabei die billigste Probe, die es gibt, und sie
hätte alle drei Fälle nicht gefangen, aber den dritten schon.
