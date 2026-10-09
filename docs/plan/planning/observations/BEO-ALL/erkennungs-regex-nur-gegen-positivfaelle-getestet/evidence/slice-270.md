**Vorgang:** slice-270
**Fund:** Die Frage, ob `go test` eine Testdatei baut, beantworteten zuerst
eigene Regeln (Testname, Signatur, `//go:build`-Auswertung); R2 fand drei
Formen, die sie falsch zählten (GOOS/GOARCH-Suffix im Dateinamen, `_`/`.` am
Anfang, `// +build`). Erst `go/build` (`Context.MatchFile`) gab die Antwort
des Werkzeugs — obwohl Schritt 19 seit slice-267 genau das verlangt.
