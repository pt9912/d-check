# `make blackbox-probe` — vergleicht das Verhalten vor und nach einer Änderung von außen

## Vertrag

Vergleicht ein **Vorher**-Image aus einer Git-Referenz (`REF=<ref>`) mit dem
**Nachher**-Image (`make build`) über dieselben Fixtures und Ausgabeformen,
stdout, stderr und Exit getrennt — zusammengeführte Streams verschränken sich
und zeigen Abweichungen, die keine sind. Was verglichen wird, wann der
Kanarienlauf den Vergleich trägt und in welchen Fällen der Lauf scheitert,
legt [`SPEC-096`](../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)
fest. Die Fixtures sind je ein kleines Repository im Default-Zustand eines
Moduls; eine gelöschte, noch getrackte Datei fällt aus der Kopie des Repos
heraus.

**Nebenwirkung:** das Target baut `$(IMAGE):latest` neu, mit
`VERSION=0.0.0-dev` — dasselbe Bild wie ein `make build` ohne Version; ein
zuvor mit Release-Version gebautes `:latest` ist danach ersetzt — auch, wenn
`VERSION` auf der Kommandozeile mitgegeben wird.

Gedacht für die Zusage „ohne Schalter unverändert": vor dem Review laufen
lassen, mit `REF` auf den Stand vor der Änderung.

## Grenze — was das Grün nicht abdeckt

1. **Es vergleicht nur, was Fixtures und Formen abdecken.** Ein Modul ohne
   Fixture ist nur über den Lauf auf dem Repo selbst erfasst; was dessen
   Konfiguration nicht einschaltet, bleibt unverglichen. `ls
   tools/blackbox-probe/fixtures/` zeigt den Ausschnitt.
2. **Es findet einen Unterschied, den niemand erwartet — nicht die falsche
   Erwartung.** Ob „unverändert" die richtige Zusage ist, bleibt beim fremden
   Leser.
3. **Verglichen werden zwei Binaries, nicht zwei Stände.** Beide Images lesen
   dieselben Eingaben — die heutigen Fixtures und den heutigen Arbeitsbaum
   samt `.d-check.yml`, nicht die von `REF`. Eine Konfiguration, die nur der
   neue Stand versteht, zeigt sich als Abweichung (das alte Binary endet mit
   Exit 2) — erwartbar, aber kein Verhaltensvergleich.
4. **Die Umgebung prüft ein Kanarienlauf, nicht der Vergleich.** Alle Läufe
   mounten aus **einer** Quelle: Kopien der Fixtures und des Repos (seine
   nicht ignorierten Dateien) unter einem Arbeitsverzeichnis, lesbar
   gemacht — eine `umask 077` im Klon ändert daran nichts. Vor und nach
   allen Läufen fährt das Target beide Images über `sauber` und `links` und
   verlangt den Exit-Vertrag allein: 0 und 1. Ein leerer Mount endet mit 2,
   ein Mount ohne Markdown bei `links` mit 0, ein Daemon-Ausfall bei
   `sauber` mit 1 — die Probe bricht dann mit Exit 2 ab, ohne vom Wortlaut
   einer Meldung abzuhängen. Jeder Exit außer 0, 1 und 2 bricht ebenfalls ab
   (Start 125/126/127, OOM 137); 0, 1 und 2 sind Verhalten des Werkzeugs und
   werden verglichen, auch ein Absturz. **Restgrenze:** ein Ausfall, der
   zwischen den Kanarienläufen nur einzelne Läufe trifft, mit 0, 1 oder 2
   endet und auf beiden Seiten gleich aussieht, bliebe unerkannt. Ändert ein
   Vorgang den Exit von `sauber` oder `links`, scheitert der Kanarienlauf —
   das ist eine Änderung am Exit-Vertrag und gehört ohnehin benannt.
5. **Kein Gate.** Eine gewollte Änderung erzeugt Abweichungen; ob eine
   Abweichung gewollt ist, entscheidet der Vorgang, der sie erzeugt.

## Ausgabe lesen

Die Exit-Codes legt
[`SPEC-096`](../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)
fest; über `make` endet jedes Scheitern mit 2, und 1 und 2 trennt die
**Ausgabe**. Eine Abweichung nennt je Fall die Ströme und einen Diff-Auszug.
Scheitert der Kanarienlauf nur auf dem Nachher-Image, sagt die Meldung, dass
eher der Exit-Vertrag des neuen Stands sich geändert hat als die Umgebung.

## Bindung

Werkzeug, kein Gate: nicht in `gates`, `ci` oder einem Nachtlauf.
