**Vorgang:** slice-217
**Fund:** **Eine Definition liest sich als Menge und war eine Auswahl** — hier
nicht in einer Grenzen-Liste, sondern in der **Klassen-Definition** selbst,
die der Slice zu seinem DoD (1) gemacht hatte.

§3 definierte *überflüssiges Leerzeichen* als Lauf von zwei oder mehr
Leerzeichen **vor** einem Pipe. Die Klasse ist damit **einseitig**. DoD (2)
behauptete daneben, danach stehe *„genau ein Leerzeichen **beidseits** jedes
Zellinhalts"* — eine Aussage, die aus der Klasse nicht folgt. Sie traf zu,
aber aus einem anderen Grund: Links-Padding gab es im Ausgangsstand nicht.
**Nicht gedeckt** waren außerdem Zeilenende-Whitespace, Tabs und `&nbsp;` —
alle vier nachgemessen und bei null, das Ergebnis steht also; die Deckung
nicht.

**Und dieselbe Klasse in der Kanten-Prüfung:** DoD (3) nannte *„die zwei
Konfigurations-Kanten nachweislich intakt"*. Eine davon — die
`structure`-Regel mit `cell-min-chars` — bindet §Sensors, einen Abschnitt, den
der Slice **nicht anfasst**. Ihr grüner Lauf ist über diese Änderung eine
Tautologie; die Zwei war eine Aufzählung, keine Deckung.

**Der gemessene Mechanismus des Eintrags trifft wieder:** Die fehlende Zeile
stand bereits in derselben Datei, in der anderen Rahmung — §3 nannte die
Klasse einseitig, §2 las sie zweiseitig, zwei Abschnitte voneinander entfernt.
