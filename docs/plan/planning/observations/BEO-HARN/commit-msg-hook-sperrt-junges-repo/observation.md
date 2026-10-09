# Der `commit-msg`-Hook sperrt ein Repo ohne oder mit einem Commit

**Sub-Area:** `HARN`

`make trace-check` fährt den Vorlauf `history-range-guard` mit der Default-Range
`HEAD~1..HEAD` auch im Modus des `commit-msg`-Hooks, der keine Historie liest.
In einem Repo ohne oder mit einem Commit und in einem Klon der Tiefe 1 ist die
Range nicht auflösbar; der installierte Hook bricht dann jeden Commit ab,
auch einen mit Kennung. Benannt in
[`SPEC-102`](../../../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge);
die Behebung — den Vorlauf im Botschafts-Modus auslassen — ist eine Änderung
am Werkzeug.
