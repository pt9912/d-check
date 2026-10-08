# Der Fix-Commit ist größer als seine Befunde und liegt außerhalb jeder Review-Range

**Sub-Area:** `*`

Ein Review benennt Befunde, die Einarbeitung schreibt dafür neuen Code — eine
neue Mechanik, nicht nur eine Korrektur am Befund. Dieser Code entsteht nach
dem Review und wird von ihm nicht gesehen; erst die Verifikation bemerkt, dass
der größte Code-Anteil des Vorgangs ungelesen ist. Ableiter: ein Fix, der mehr
schreibt, als seine Befunde verlangen, bekommt eine eigene Review-Runde über
genau seinen Commit, bevor die Verifikation läuft.
