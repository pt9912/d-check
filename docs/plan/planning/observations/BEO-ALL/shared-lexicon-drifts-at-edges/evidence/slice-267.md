**Vorgang:** slice-267
**Fund:** Das Modul `vcs` beantwortete „was ist ein Link" mit eigenen
zeilenweisen Mustern statt mit der Link-Erkennung des Moduls `links`; vier
Runden flickten die Muster, und ein Fix (Code-Span-Überlappung) machte den
Anlassfall selbst zur Drift (R4 H-1, I-3). Erst die geteilte Erkennung
(`PreprocessMarkdown`, `ExtractLinkSpans`, `definitionRe`) trug.
