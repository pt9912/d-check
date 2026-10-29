**Vorgang:** slice-279
**Fund:** Die Negativliste des Plans war vollständig umgesetzt und getestet, und
trotzdem stand der `observations`-Ausschluss im Code vor dem Ziffern-Fall: ein
Beleg `evidence/slice-007-foo.md`, bisher ein erkannter Slice, wurde keiner mehr
— und der Negativ-Test pinnte genau das. Die Liste prüfte, was die neue Regel
ausschließt, nicht was sie dem **bisher erkannten** Bestand wegnimmt. Der Review
fand es (F-1, MEDIUM); der Vergleich der RTM gegen `v0.87.1` am Eigenbestand hatte
es nicht gezeigt, weil dessen Belege titellos heißen.
