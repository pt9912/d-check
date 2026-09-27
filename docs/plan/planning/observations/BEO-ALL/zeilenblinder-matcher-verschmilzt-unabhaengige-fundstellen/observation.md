# Ein zeilenblinder Klammer-Matcher über mehrzeilig gefaltetem Text verschmilzt unabhängige Fundstellen

**Sub-Area:** `*`

Wird ein String-Matcher, der keine Zeilengrenzen kennt (`matchBracket` zählt
`[`/`]`- bzw. `(`/`)`-Tiefe ohne Rücksicht auf Zeilenumbrüche), über mehrere
mit `\n` zusammengefügte Zeilen laufen gelassen, um eine zeilenübergreifende
Form zu erkennen, kann ein unbalanciertes Öffnungszeichen aus einer völlig
unabhängigen, weiter oben stehenden Zeile mit einer späteren, ebenfalls
unabhängigen Schließsequenz zu einem erfundenen Fund verschmelzen — ein
Befund, der ohne die Erweiterung nicht existierte, aus zwei inhaltlich
fremden Textstellen. Der sichere Entwurf begrenzt den Lookahead auf **eine**
Nachbarzeile und **nur** auf die Klammer, die die neue Form tatsächlich
verlangt — nie auf beide Klammern eines mehrteiligen Tokens gemeinsam.
