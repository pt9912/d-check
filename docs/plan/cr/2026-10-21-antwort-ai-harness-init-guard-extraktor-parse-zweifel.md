# Antwort von `ai-harness-init` — Der Extraktor des Tool-Call-Wächters lässt zwei Parse-Zweifel durch

**Absender:** `ai-harness-init` · **Eingegangen:** 2026-10-21
**Richtung:** eingehende Antwort auf den ausgehenden
[CR](2026-10-21-cr-ai-harness-init-guard-extraktor-parse-zweifel.md)
([`MR-036`](../../../harness/conventions.md#mr-036))
**Entscheidung:** angenommen und umgesetzt, auf `main` des Werkzeug-Repos
geschlossen; noch nicht ausgeliefert — ein Ziel-Repo bekommt den Fix mit dem
nächsten Release.

---

## Wortlaut

> Antwort an d-check zum Change Request 2026-10-21-cr-ai-harness-init-guard-extraktor-parse-zweifel
>
> Von: ai-harness-init · Datum: 2026-10-21
> Stand: umgesetzt und geschlossen auf main bei v0.10.0-172-gd8312c3b. Noch nicht ausgeliefert: Die emittierte Fassung steckt im Binary, und v0.10.0 bleibt unverändert. Ein Ziel-Repo bekommt den Fix erst mit dem nächsten Release.
>
> Befund bestätigt. Beide Eingaben aus dem Change Request lieferten bei uns Exit 0, im busybox-awk 1.36.1 des gepinnten BATS_IMAGE (dort läuft make test) und im BWK-awk 20200816 auf dem macOS-Host:
> - {"tool_input":{"command":"pip x""make gates"}} ergab make gates.
> - {"tool_input":{"command":"a" garbage }} ergab a.
>
> Ursache. Nach einem String-Ende prüfte der Scanner nicht, was folgen darf. Ein zweites " öffnete still einen neuen String, und Buchstaben außerhalb von Strings wurden ignoriert.
>
> Umsetzung. Eine Zustandsmarke afterstr wird am String-Ende gesetzt. Danach wird Leerraum übersprungen, und jedes Zeichen außer : , } ] beendet den Extraktor mit Exit 3, auch ein weiteres ". Eines dieser vier Zeichen setzt die Marke zurück. Die Regel haben wir selbst geschrieben und von Ihrem Extraktor nur die Parse-Regel gelesen, keinen Code übernommen. Der Fix steht in beiden Fassungen, harness/tools/extract-command.awk und der emittierten Vorlage. Die Code-Zeilen beider sind gleich, nur die Kommentare unterscheiden sich. Die Gleichheit ist eine einmal gemessene Handmessung und kein Sensor.
>
> Gemessen (im gepinnten BATS_IMAGE, --network none):
> - Vorher: beide Eingaben Exit 0.
> - Nachher, beide Fassungen: beide Eingaben Exit 3 und leere Ausgabe. Dazu "a" "b" und "a"x, ebenfalls Exit 3.
> - Falsch-Positive: Ein Gegenlauf mit 20 wohlgeformten Eingaben in drei Fassungen (lokal, Vorlage, Vorher-Stand) ergab keinen Falsch-Positiv. Geprüft waren Escapes (\", \\, auch \" direkt vor dem schließenden "), Folgefelder wie description und timeout, Literale, Zahlen, verschachtelte Objekte und Arrays, Leerraum, Zeilenumbruch, CRLF, BOM und ein String am Dateiende.
>
> Zu Bitte 2 (Probenmenge). test/guard.bats führt die Fälle jetzt in einer Schleife über beide Fassungen, dazu eine Regressionsprobe für wohlgeformtes JSON. Wir haben zusätzlich vier Mutations-Fälle (780 bis 783). Zwei heben die neue Prüfung auf (Lücke), zwei die Rücksetz-Zeile der Marke (Falsch-Positiv-Richtung). Alle vier färben ihren gepinnten Test rot. Die Proben stammen von Ihrem Abschnitt „Fail-closed am Extraktor" als Vorlage, die Fälle sind unsere.
>
> Was nicht gemessen ist:
> - gawk und mawk. Kein gepinntes Image führt sie, und ein Host-Lauf wäre unter unserer Norm kein Beleg. Der Extraktor-Kopf sagt weiter „POSIX-awk (busybox/gawk/BSD)". Gemessen sind nur busybox und BWK.
> - Ein realer Guard-Lauf im emittierten Ziel mit einer Parse-Zweifel-Eingabe. Unser full-smoke hat dafür keine Stufe. Die emittierte Fassung ist nur über den Extraktor der Vorlage gedeckt (Schleife und Fälle 781 und 783).
> - Resttext hinter dem Wurzelobjekt ({…}x) und Syntaxfehler ohne vorangehenden String ({"a":1 2}) bleiben Exit 0. Das war nicht Gegenstand.
>
> Zu „Nicht gefordert". Wir haben nichts davon angefasst: Interpreter-Sperre, Brace-Group, Sub-Shell-Rekursion, \u-Escape, den Exit-Riegel und das leere PATH. Zu Ihrem Hinweis auf [MR-044](../../../harness/conventions.md#mr-044): Bei uns trägt diese Nummer eine andere Regel („Das Technik-Stratum trägt die ID-Spalte der Ziel-Form"). Der Zweikanal-Wächter ist dort nicht abgebildet. Bei der Zählung zwischen den Repos bitte nicht gleichsetzen.
>
> Commits (auf main): 2afb5351 Extraktor, Proben und Fälle 780/781 · 7123b413 Kopfkommentare der beiden Guard-Fassungen · 507ec7a2 Grenze der emittierten Fassung am Ort · 91833454 Fälle 782/783.
>
> Offen auf unserer Seite: Dieses Repo hat keine Eingangsform für Change Requests von Adoptern. Wir haben Ihren Change Request als wellenlosen Slice behandelt und halten das so fest, bis ein zweiter Fall eine Form begründet.
