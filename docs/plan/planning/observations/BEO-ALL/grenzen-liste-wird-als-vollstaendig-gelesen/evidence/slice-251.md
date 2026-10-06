**Vorgang:** slice-251
**Fund:** Drittes Auftreten in Folge am Rand des Moduls `targets` (nach
slice-249 und slice-250), diesmal im Dekodier-Verhalten der YAML-Bibliothek.
Lastenheft, Spezifikation und ADR sagten „leerer Eintrag ⇒ Exit 2" und
Byte-Identität der String-Form zu; der Code dekodierte Listen in
`[]string`, wobei die Bibliothek Null-Elemente still verwirft (Richtung 2
entfiel unbemerkt), und lehnte Glob-Zeichen ab, was einen wörtlichen Pfad mit
eckiger Klammer brach. Die Unit-Tests trugen nur die beschriebenen Fälle;
gefunden hat beides der unabhängige Review per Black-Box gegen ein
Vorher-Image (R1-H1, M1).
