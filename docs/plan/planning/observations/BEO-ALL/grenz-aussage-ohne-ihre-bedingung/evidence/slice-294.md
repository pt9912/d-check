**Vorgang:** slice-294
**Fund:** Die Kommentare in `linelength.go` (Dateikopf) und `markdown.go` (am Helfer) begründen die gemessenen Fälle Tab, U+00A0 und `>` mit „die Lexik hat keinen Listenkontext“; die Begründung trägt nur für vier und mehr Leerzeichen, `> |` bleibt wegen einer Entscheidung der ADR gemessen und U+00A0 ist keine Einrückung (Review R1 F-1).
