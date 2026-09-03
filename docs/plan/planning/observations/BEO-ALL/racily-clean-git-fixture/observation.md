# Ein Test-Fixture ändert Inhalt, ohne die Größe zu ändern, und fällt darum sporadisch aus

**Sub-Area:** `internal/adapter/driven/git`

Ein Fixture schrieb dieselbe Datei mit gleicher Byte-Länge in derselben
Sekunde neu; eine stat-basierte Änderungserkennung führte sie dann als
unverändert (*racily clean*). Jedes git-Fixture, das Änderung über gleich
lange Inhalte in einem Zug herstellt, trägt dieselbe Klasse. Bereits
informell adressiert durch `coretest.GitFixtureRewriteHazard`
([`BEO-ALL/gruene-wettlauf-probe-beweist-nichts`](../gruene-wettlauf-probe-beweist-nichts/observation.md)),
gemessener Bestand: zehn Fundstellen in acht Testfunktionen, zwei Paketen.
