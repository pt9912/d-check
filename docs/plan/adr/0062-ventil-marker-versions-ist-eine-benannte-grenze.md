# ADR-0062 — `versions` bleibt roh als **benannte Grenze**, nicht als andere Frage (ersetzt ADR-0061)

**Status:** Accepted
**Datum:** 2026-08-27
**Supersedes:** [ADR-0061](0061-ventil-marker-geteilte-antwort-wo-die-eingabe-prosa-ist.md)
— die **Entscheidung** bleibt unverändert (`codepaths` und `ids` angeglichen,
`versions` und `diagrams` roh); ihre **Herleitung**, ihre **Zahlen** und ihre
**Grenzen** waren falsch bzw. unvollständig. Der Kern einer `Accepted`-ADR wird
nicht überschrieben (AGENTS.md §3.5), deshalb dieser Nachfolger.
**Autor:** pt9912
**Schärft:** [`DC-FA-CODE-001.a`](../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code)
Schritt 1 und [`DC-FA-ID-001.a`](../../../spec/spezifikation.md#dc-fa-id-001a--kennungs-prüfung)
Bedingung 4 (je die Erkennung des Zeilen-Markers)
**Bezug:**
[ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) (Entscheidung 1 und
ihr Re-Evaluierungs-Trigger),
[ADR-0060](0060-citations-marker-scan-geteilte-prosa-antwort.md) (dieselbe
Frage für die Zitat-Direktive — und die **drei** Grenzen, die dort benannt sind),
[ADR-0019](0019-versions-pin-fence-ausnahme.md) (die Fence-Ausnahme des
**Pin-Scans** von `versions` — **nicht** seiner Marker-Erkennung),
[ADR-0025](0025-codepaths-ignore-refs.md) (das Ventil selbst)

## Kontext

[ADR-0061](0061-ventil-marker-geteilte-antwort-wo-die-eingabe-prosa-ist.md) hat
`codepaths` und `ids` auf die geteilte Prosa-Antwort umgestellt und `versions`
und `diagrams` roh gelassen. Der unabhängige Review hat die **Umstellung**
bestätigt und **vier** tragende Aussagen widerlegt. Diese ADR ersetzt sie; das
Verhalten des Produkts ändert sich dadurch **nicht**.

**1. „`versions` liest keine Prosa" ist falsch.** Es liest **alle** Zeilen —
eine **Obermenge**, die die Prosa-Zeilen einschließt. Auf einer Prosa-Zeile gibt
es sehr wohl Inline-Code, und dort antwortet das Produkt seit
[ADR-0061](0061-ventil-marker-geteilte-antwort-wo-die-eingabe-prosa-ist.md)
**zweifach**. Gemessen an einer Zeile mit `` `d-check:ignore` `` in Backticks,
einem veralteten Pin und einem nicht auflösbaren Pfad:

| Modul | Verhalten auf **derselben** Zeile |
|---|---|
| `codepaths` | meldet — der Marker ist eine **Erwähnung** |
| `versions` | schweigt — der Marker ist eine **Direktive** |

Die Kontrollzeile ohne Marker meldet beides. Die Divergenz ist also **verschoben,
nicht behoben** — für `diagrams` trägt die Begründung, für `versions` nicht.

**2. Die Zahlen waren die falsche Grundgesamtheit.** *„160 Zeilen ausschließlich
in Inline-Code"* war die Zahl der Zeilen mit **barer Marker-Form** — eine andere
Frage. Nachgemessen mit der Produkt-Lexik über den Stand dieser ADR:

| Größe | Wert |
|---|---|
| getrackte Markdown-Dateien außerhalb des vendorten Baums | **553** |
| Prosa-Zeilen mit dem Marker | **249** |
| davon **frei** (wirksam) | **66** |
| davon **nur in Inline-Code** | **183** |
| unter den 66 wirksamen: bare Form ohne Kommentar-Klammer | **1** |

Die 173 aus [ADR-0060](0060-citations-marker-scan-geteilte-prosa-antwort.md)
waren für **ihren** Stand richtig; die Menge ist seither gewachsen. Die
Fehlklassifikation war, die Beschriftung der einen Messung für die Zahl der
anderen zu übernehmen.

**3. Zwei Grenzen des Strippens wurden mitgeerbt und nicht benannt.**
[ADR-0060](0060-citations-marker-scan-geteilte-prosa-antwort.md) hat sie für die
Zitat-Direktive ausdrücklich geführt; für den Ventil-Marker gelten sie ebenso,
mit **entgegengesetzter** Richtung — beide gemessen:

- **Verschluckung ⇒ Falsch-Rot.** Ein **gesetzter**, freier Marker wirkt
  **nicht**, wenn eine Code-Spanne desselben Absatzes ihn umschließt.
- **Erzeugung ⇒ stilles Grün.** Steht ein **unpaariger** Backtick weiter oben im
  Absatz, kippt die Parität, und die Erwähnung in Backticks **wirkt weiterhin
  als Direktive** — gemessen: der Pfad auf jener Zeile wird nicht gemeldet, die
  Kontrollzeile schon.

Die Zusage aus
[ADR-0061](0061-ventil-marker-geteilte-antwort-wo-die-eingabe-prosa-ist.md)
§Konsequenzen — *„Eine Zeile, die das Ventil beschreibt, nimmt sich nicht mehr
selbst aus"* — gilt damit **nur für Absätze mit gerader Backtick-Parität**.

**4. Zwei Quellen wurden über ihren Geltungsbereich hinaus zitiert.**
[ADR-0019](0019-versions-pin-fence-ausnahme.md) skopiert ausdrücklich den
**Pin-Scan** (*„Der Fence-Scan greift ausschließlich für das konfigurierte
`pin-pattern`"*); über die Marker-Erkennung sagt sie nichts. Und
[ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 2 nennt
als *andere Frage* die **Pins** in Fences, nicht den Marker. *„Ist das eine
Direktive"* ist bei `versions` **dieselbe** Frage wie bei `codepaths` — sie wird
dort nur anders beantwortet.

## Entscheidung

1. **`codepaths` und `ids` bleiben angeglichen** — unverändert gegenüber
   [ADR-0061](0061-ventil-marker-geteilte-antwort-wo-die-eingabe-prosa-ist.md).
   Die Frage *„ist diese Zeile eine Direktive"* bekommt dort die geteilte
   Prosa-Antwort.

2. **`diagrams` bleibt roh, und das ist strukturell.** Es liest die Zeilen
   **innerhalb** eines Fence und zusätzlich die **Öffnungszeile** — beides keine
   Prosa, dort ist ein Backtick literaler Inhalt. Es gibt keine geteilte Antwort
   zu übernehmen. (Die Öffnungszeile fehlte in der Tabelle des Vorgängers; sie
   ist eine dritte Eingabe-Klasse und trägt bei einer `~~~`-Fence Backticks.)

3. **`versions` bleibt roh als BENANNTE GRENZE, nicht als andere Frage.** Seine
   Eingabe **enthält** Prosa-Zeilen; auf ihnen wäre die geteilte Antwort
   möglich und wird bewusst nicht gegeben. Der Grund ist nicht struktureller,
   sondern **abgrenzender** Art: der Marker nimmt dort die Zeile **allen**
   Muster-Quellen-Paaren aus, und die Umstellung wäre eine eigene Messung mit
   eigener Sprengweite — sie an dieser hier vorbeizuschieben hieße, sie an
   deren Grünheit zu messen. **Die Zweifach-Antwort besteht dort fort und ist
   ab hier ausgewiesen**, nicht wegerklärt.

4. **Die zwei Grenzen des Strippens stehen im Vertrag**, mit Richtung:
   Verschluckung ist Falsch-Rot, Erzeugung ist stilles Grün. Keine wird
   behoben — beide folgen aus der geteilten Lexik, und eine eigene Lexik dafür
   wäre genau die zweite Antwort, gegen die
   [ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) geschrieben ist.

5. **Der Ertrag im Bestand ist null gefundene Doku-Defekte.** Die fünf
   freigelegten Befunde brauchten die Ausnahme **wirklich** — ihre Kennungen und
   Pfade sind Beispiele. Was sich ändert, ist dass die Ausnahme **gesetzt** ist
   statt aus der eigenen Prosa zu folgen. Der Vorgänger nannte sie „echt" und
   widersprach damit sowohl der Definition des Slice-Plans als auch seinem
   eigenen §Entscheidung 4. **Die Begründung des Angleichs liegt im künftigen
   Fall, nicht im gemessenen Bestand** — und das ist eine schwächere, aber
   haltbare Begründung.

6. **Der Quelltext-Wächter bleibt, mit gemessener Reichweite.** Er ist die
   Antwort auf [ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md)s
   Gate-Frage. Von sechs konstruierten Umgehungen fängt er **zwei**; die
   Erlaubnis-Liste ist **datei-granular**, eine *zusätzliche* Roh-Lesung in
   einer bereits gelisteten Datei bliebe unsichtbar. Beides ist seine Grenze,
   nicht sein Fehler — er ist ein Stolperdraht für den nächsten Konsumenten,
   keine Beweisführung. **Zurückgenommen wird die Behauptung, der Compiler fange
   die naheliegende Umgehung**: das gilt an genau einer Stelle, wo eine Variable
   unbenutzt würde, nicht allgemein.

## Alternativen

- **`versions` mit angleichen.** Verworfen nach Entscheidung 3 — eigene
  Sprengweite, eigene Messung, und der Marker wirkt dort über alle
  Muster-Quellen-Paare. Als Kandidat benannt, nicht als Rest.
- **Die zwei Strip-Grenzen schließen.** Verworfen nach Entscheidung 4: dafür
  bräuchte der Marker eine eigene Lexik neben der geteilten.
- **Den Wächter auf alle sechs Umgehungen weiten.** Teilweise umgesetzt (das
  Literal und der `Index`-Aufruf sind ergänzt); die restlichen drei brauchen
  eine Quelltext-Analyse statt einer Zeilen-Suche und wären ein eigenes
  Werkzeug.
- **ADR-0061 mit einer `## Geschichte`-Zeile korrigieren.** Verworfen: eine
  Geschichte-Zeile hält fest, was sich geändert hat — sie kann keine falsche
  **Herleitung** im Kern richtigstellen. §3.5 sieht dafür den Nachfolger vor.

## Konsequenzen

- Das Produktverhalten ist **unverändert** gegenüber
  [ADR-0061](0061-ventil-marker-geteilte-antwort-wo-die-eingabe-prosa-ist.md).
  Diese ADR korrigiert Herleitung, Zahlen und Grenzen.
- **Die Zweifach-Antwort besteht bei `versions` fort** und ist ab hier eine
  ausgewiesene Grenze. Wer sie schließt, misst ihre Sprengweite eigens.
- Die Zusage *„eine Zeile, die das Ventil beschreibt, nimmt sich nicht mehr
  selbst aus"* gilt **nur bei gerader Backtick-Parität** des Absatzes.
- **Die Marker-Form bleibt ungeregelt.** Der Vertrag nennt an zwei Stellen einen
  HTML-Kommentar; die Erkennung akzeptiert jede Zeichenkette außerhalb von Code.
  Unter den **66** heute wirksamen Markern trägt genau **einer** die bare Form —
  die Exposition ist damit klein und benannt.

## Re-Evaluierungs-Trigger

**Permanent** für Entscheidung 2 — sie folgt aus der Eingabe-Art.

**Wiedervorlage für Entscheidung 3**, sobald die `versions`-Sprengweite gemessen
ist: dann ist zu entscheiden, ob die Grenze bleibt oder fällt.

**Wiedervorlage für Entscheidung 4**, wenn eine der zwei Grenzen im Bestand
eintritt — heute tut es keine.

**Wiedervorlage für Entscheidung 6**, wenn eine der vier ungefangenen Umgehungen
im Bestand auftaucht; dann ist die Frage, ob die Klasse ein Werkzeug braucht,
das Quelltext analysiert statt Zeilen zu lesen.
