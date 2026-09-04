# ADR-0061 — Der Ventil-Marker bekommt die geteilte Prosa-Antwort, wo seine Eingabe Prosa ist — und nur dort

**Status:** Accepted
**Datum:** 2026-08-27
**Autor:** pt9912
**Schärft:** [`DC-FA-CODE-001.a`](../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code)
und [`DC-FA-ID-001.a`](../../../spec/spezifikation.md#dc-fa-id-001a--kennungs-prüfung)
(je der Zeilen-Marker als Ventil)
**Bezug:**
[ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) (Entscheidung 1 und
ihr Re-Evaluierungs-Trigger — dieser Eintrag ist seine Antwort),
[ADR-0060](0060-citations-marker-scan-geteilte-prosa-antwort.md) (dieselbe
Frage für die Zitat-Direktive; Entscheidung 7 hat diesen Fall gescopt),
[ADR-0025](0025-codepaths-ignore-refs.md) (das Ventil selbst),
[ADR-0019](0019-versions-pin-fence-ausnahme.md) (die Roh-Lesung von `versions`)

## Kontext

Seit [ADR-0060](0060-citations-marker-scan-geteilte-prosa-antwort.md)
beantwortet das Produkt die Frage *„ist diese Zeile eine Direktive"*
**zweifach**: die Zitat-Direktive wird auf dem inline-code-gestrippten Text
erkannt, der Ventil-Marker weiter auf der **rohen** Zeile. Das ist der Defekt
aus [ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 1 —
*zwei Antworten auf dieselbe Frage in einem Lauf sind ein stiller Grün-Pfad,
den kein Gate sieht* — und zugleich ihr eigener Re-Evaluierungs-Trigger
(*„eine vierte Stelle beantwortet eine Lexik-Frage selbst"*).

**Die vier Konsumenten sind nicht gleich.** Zwei lesen Prosa-Zeilen, zwei
nicht — und das entscheidet, nicht der Aufwand:

| Konsument | Eingabe | geteilte Prosa-Antwort? |
|---|---|---|
| `codepaths` | Prosa-Zeilen (Fences entfallen) | **ja** |
| `ids` | Prosa-Zeilen | **ja** |
| `versions` | **alle** Zeilen, einschließlich Fenced-Code ([ADR-0019](0019-versions-pin-fence-ausnahme.md)) | nein |
| `diagrams` | die Zeilen **innerhalb** eines Fence | nein |

Innerhalb eines Fenced-Blocks ist ein Backtick **literaler Inhalt**; es gibt
dort kein Inline-Code, das zu strippen wäre. Für `versions` und `diagrams`
existiert die geteilte Antwort also gar nicht — sie beantworten eine andere
Frage, genau im Sinne von
[ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 2.

**Die Folge der Roh-Lesung ist gemessen, nicht geschätzt.** Über 544 getrackte
Markdown-Dateien tragen **250** Prosa-Zeilen den Marker; **160** davon
ausschließlich in Inline-Code. Der Angleich der zwei Prosa-Konsumenten legt
**fünf** Befunde frei — und alle fünf sind **echt**: es sind Zeilen des
Lastenhefts, die das Ventil **beschreiben** und sich durch ihre eigene
Beschreibung ausgenommen haben. Die Ausnahme war nie beabsichtigt.

Die im Vorfeld genannte Obergrenze von 58 Befunden galt einem **ganz**
abgeschalteten Ventil, nicht diesem Angleich; sie ist um den Faktor zehn zu
groß für die Frage, die hier entschieden wird.

**Kein Test hat das alte Verhalten je behauptet.** Die vollständige Suite läuft
gegen die angeglichene Fassung grün — die Roh-Lesung war eine
Implementierungs-Eigenschaft, keine zugesagte.

## Entscheidung

1. **Wo die Eingabe Prosa ist, gilt die geteilte Antwort.** `codepaths` und
   `ids` erkennen den Marker auf dem fence-bewussten, inline-code-gestrippten
   Text. Ein Marker in Backticks ist eine **Erwähnung**, keine Direktive.

2. **Wo die Eingabe keine Prosa ist, bleibt die Roh-Lesung — als Skopierung,
   nicht als Rest.** `versions` und `diagrams` sind unverändert. Der Grund ist
   strukturell: in ihrer Eingabe gibt es kein Inline-Code. Eine „Angleichung"
   wäre dort keine, sondern eine neue, dritte Antwort.

3. **Ein Test hält die Klasse, nicht ein Review.** Das ist die Antwort auf
   [ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md)s
   Re-Evaluierungs-Trigger *„kann ein Gate die Klasse prüfen statt eines
   Reviews"*: ein Wächter über den **Quelltext** des Regel-Pakets meldet jeden
   lesenden Zugriff auf den Marker über eine rohe Zeile, sofern die Datei nicht
   mit **Grund** in einer Liste steht — und meldet umgekehrt einen Eintrag,
   den niemand mehr braucht. **Seine Grenze steht bei ihm:** er liest
   Schreibweisen, nicht Verhalten, und fängt keine Umgehung über eine
   Zwischenvariable oder ein anderes Paket. Er ist ein Stolperdraht für den
   nächsten Konsumenten.

4. **Die fünf freigelegten Zeilen bekommen einen ausdrücklichen Marker, keine
   Ausnahme-Liste.** Sie brauchen die Ausnahme wirklich — ihre Kennungen und
   Pfade sind Beispiele, keine Verweise. Der Unterschied ist, dass die Ausnahme
   jetzt **gesetzt** ist statt aus der eigenen Prosa zu folgen.

5. **Die Form des Markers ist hier nicht Gegenstand.** Die Spezifikation nennt
   ihn einen **HTML-Kommentar**; die Implementierung akzeptiert jedes Vorkommen
   der Zeichenkette, und **160 von 250** Zeilen im Bestand tragen die bare Form.
   Das ist eine zweite, größere Lücke zwischen Vertrag und Code — eine **andere
   Frage** als die Inline-Code-Blindheit, und sie wird hier weder entschieden
   noch stillschweigend mitgezogen.

## Alternativen

- **Alle vier Konsumenten angleichen.** Verworfen nach Entscheidung 2: für zwei
  von ihnen gibt es keine Prosa-Antwort zu übernehmen. Der Angleich wäre dort
  eine Erfindung, keine Übernahme — und genau die Sorte zweiter Antwort, gegen
  die [ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) geschrieben ist.
- **Gar nicht angleichen und die Zweifach-Antwort als Grenze führen.**
  Verworfen: die Messung zeigt fünf **echte** Befunde, also einen wirksamen
  stillen Grün-Pfad. Eine benannte Grenze verwaltet ihn, sie behebt ihn nicht.
- **Die fünf Zeilen umschreiben, statt sie zu markieren.** Verworfen: die
  Kennungen sind Teil der Akzeptanzkriterien und des reproduzierten
  Fremd-Befundes; sie umzuschreiben verfälschte den Vertragstext, um ein Gate
  ruhigzustellen.
- **Zugleich die Marker-Form auf den HTML-Kommentar verengen.** Verworfen nach
  Entscheidung 5 — 160 Zeilen sind eine eigene Entscheidung mit eigener
  Messung, und sie an dieser hier vorbeizuschieben hieße, die zweite an der
  Grünheit der ersten zu messen.

## Konsequenzen

- Eine Zeile, die das Ventil **beschreibt**, nimmt sich nicht mehr selbst aus.
  Das ist der Zweck; es ist zugleich der Preis für jeden, der den Marker künftig
  in Prosa erwähnt und die Ausnahme dann **setzen** muss.
- **Die zwei Antworten bestehen fort**, nur nicht mehr unbegründet: sie folgen
  jetzt der Eingabe-Art und stehen als Tabelle im Vertrag. Wer eine dritte
  Antwort einführt, stolpert über den Wächter aus Entscheidung 3.
- Der Wächter kostet einen Eintrag pro legitimer Roh-Lesung. Wächst die Liste
  ohne Grund, ist das sichtbar — das ist seine zweite Richtung.
- **Die Marker-Form bleibt offen**, mit gemessener Größe (160 von 250). Sie ist
  benannt, nicht verdeckt.

## Re-Evaluierungs-Trigger

**Permanent** für Entscheidung 2 — sie folgt aus der Eingabe-Art, nicht aus
einer Lage, die sich ändern kann. Ändert ein Modul seine Eingabe (etwa
`versions` auf Prosa-Zeilen), ist es kein Trigger dieser ADR, sondern eine
Änderung ihrer Voraussetzung.

**Wiedervorlage für Entscheidung 5**, sobald die Marker-Form entschieden wird:
verengt sie sich auf den HTML-Kommentar, ist der Wächter aus Entscheidung 3
auf die neue Form zu ziehen.

**Wiedervorlage für Entscheidung 3**, wenn der Wächter eine Umgehung durchlässt,
die im Bestand eintritt — dann ist die Frage, ob die Klasse ein Werkzeug
braucht, das Verhalten statt Schreibweisen liest.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-08-27 | Accepted | [slice-162](../planning/done/slice-162-ignore-marker-geteilte-antwort.md) |
| 2026-08-27 | Herleitung, Zahlen und Grenzen abgelöst — das Produktverhalten bleibt | [ADR-0062](0062-ventil-marker-versions-ist-eine-benannte-grenze.md) |
