# ADR-0060 — Der `citations`-Marker-Scan bekommt die geteilte Prosa-Antwort; der Zitattext bleibt roh

**Status:** Accepted
**Datum:** 2026-08-27
**Autor:** pt9912
**Schärft:** [`DC-FA-CITE-001.a`](../../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations)
(Kopfsatz und Schritt 1, Marker-Erkennung)
**Bezug:**
[`DC-FA-CITE-001`](../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in)
(die Anforderung),
[ADR-0045](0045-zitat-verifikation-codepaths-erweiterung-und-citations-modul.md)
(das Modul und seine Fail-closed-Entscheidung),
[ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) (die Trennlinie
*andere Antwort* gegen *andere Frage*),
[ADR-0050](0050-fence-unclosed-in-spans.md) (die Fence-Lexik selbst),
[ADR-0025](0025-codepaths-ignore-refs.md) (die Schwester-Direktive `d-check:ignore`)

## Kontext

[`DC-FA-CITE-001.a`](../../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations)
sagt über den Marker-Scan zu: *„Arbeitet auf den rohen Zeilen (fence-aware wie
die übrigen Module)."* Die Implementierung löste die zweite Hälfte des Satzes
ein und die erste wörtlich: fence-bewusst ja, aber **inline-code-blind** —
anders als jedes andere prosa-lesende Modul, deren gemeinsame Vorverarbeitung
Backtick-Spannen positionserhaltend leert.

**Die Folge trägt der Bestand.** Gemessen mit der **Produkt-Lexik**
(Fence-Automat, absatzweise Spannen gleicher Backtick-Länge) über alle 544
getrackten Markdown-Dateien außerhalb des vendorten Baums — der Gegenstand ist
der **Marker**, also ein geöffneter HTML-Kommentar, nicht die bloße Erwähnung
der Zeichenkette:

| Lage | Marker-Zeilen | Form |
|---|---|---|
| Inline-Code | **24** | 20 malformt, 4 wohlgeformt |
| Fenced-Block | **1** | malformt, schon vorher übersprungen |
| frei in Prosa | **0** | — |
| **gesamt** | **25** in 13 Dateien | |

**Null** davon ist eine produktive Direktive: die vier wohlgeformten nennen
zweimal denselben Fixture-Pfad eines Review-Reports, einen Pfad aus einer nicht
mehr vendorten Baseline-Fassung und das Beispiel aus dem Lastenheft selbst. Und
weil **kein einziger** Marker frei in Prosa steht, liest das Modul nach dieser
Entscheidung **gar keinen** mehr — der Bestand, der fail-closed auslösen kann,
fällt von **25 auf 0**. Gemessen: vorher endet der Lauf an
`CHANGELOG.md:592`, nachher meldet er 546 Dateien, 0 Befunde, Exit 0.

**Die Zeichenketten-Erwähnung ist eine andere, größere Menge** und hier
ausdrücklich **nicht** der Gegenstand: der Direktiven-Name kommt als bloßes
Token **78**-mal auf **74** Zeilen in **22** Dateien vor. Die Differenz sind
Sätze, die ihn nennen, ohne einen Kommentar zu öffnen; sie können nie eine
Direktive sein. Wer beide Zahlen mischt, misst richtig und redet über etwas
anderes.

**Ein Ventil gibt es, aber kein passendes.** `citations` trägt **keine feine**
Achse — kein `exempt-paths`, kein `ignore-refs`, keinen Zeilen-Marker. Grob
wirken `scan.ignore` und `citations.scope` (Wurzeln plus `ignore`). Beide nehmen
die **ganze Datei** aus dem Modul, und genau das taugt hier nicht: die 25 Marker
liegen in `CHANGELOG.md`, beiden `README`-Fassungen, beiden Spec-Straten, dem
Benutzerhandbuch und den ADRs — also in den Dokumenten, die **echte** Zitate
tragen oder tragen werden. Wer sie ausschließt, schaltet das Modul dort ab, wo
es wirken soll.

**Eine Doku-Konvention *„Syntax nur noch in Fences"* ist nicht teuer, sondern
unmöglich:** **12** der 25 Marker liegen in **fünf** unantastbaren Dateien —
vier eingefrorenen Review-Reporten und
[ADR-0045](0045-zitat-verifikation-codepaths-erweiterung-und-citations-modul.md),
die `Accepted` und damit von §3.5 geschützt ist.

## Entscheidung

1. **Die Frage „ist diese Zeile eine Direktive" ist eine Prosa-Frage und
   bekommt die geteilte Antwort.** [ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md)
   Entscheidung 2 zieht die Trennlinie zwischen *andere Antwort* (Defekt) und
   *andere Frage* (per ADR gescopt) und nennt als Beispiele der geteilten Seite
   *„ist das eine Überschrift", „ist das ein Anker", „ist das derselbe
   Absatz"*. *„Ist das eine Direktive"* steht in derselben Reihe: sie fragt
   nach der **Rolle** einer Zeile im Dokument, nicht nach ihren Bytes. Marker-
   Suche **und** Direktiven-Parse laufen deshalb auf dem fence-bewussten,
   inline-code-gestrippten Text. Eine in Inline-Code geschriebene
   Direktiv-Syntax ist eine **Erwähnung**.

2. **Der Zitattext bleibt roh.** Er ist die andere Frage: dort sind die Bytes
   der Gegenstand des Vergleichs, nicht ihre Prosa-Rolle — dieselbe Bauform,
   mit der `pins` den gehashten Span roh liest
   ([ADR-0020](0020-content-pin-fence-ausnahme.md)) und `versions` die Pins auch
   in Fences erkennt ([ADR-0019](0019-versions-pin-fence-ausnahme.md)). Ein
   Zitat, das Backticks enthält, muss zeichengenau vergleichbar bleiben; ein
   Stripping machte den Vergleich stillschweigend großzügig.

3. **Auch der Parse läuft auf dem gestrippten Text, nicht nur die Suche.**
   Trägt eine Zeile beides — eine Erwähnung in Backticks und daneben eine echte
   Direktive —, träfe der Regex auf der rohen Zeile die **Erwähnung** zuerst
   und zitierte gegen deren Pfad. Der Unterschied ist kein Randfall der
   Ästhetik: er entschiede über die geprüfte Quelle.

4. **Der Pfad im Muster beginnt mit einem Nicht-Whitespace-Zeichen.** Das ist
   die Folge von Entscheidung 3, keine Kosmetik: steht der **Pfad** in
   Backticks, ersetzt das positionserhaltende Strippen ihn durch Leerzeichen.
   Ohne diese Forderung wäre das ein gültiger Parse mit **leerem** Ziel, und der
   Befund nennte die nackten Zeilennummern statt eines Pfades. Ein fehlender
   Pfad ist ein malformter Span und damit fail-closed, wie Schritt 1 es
   vorsieht.

5. **Fail-closed bleibt, und die Messung ist jetzt sein Argument.** Nach
   Entscheidung 1 fällt der auslösende Bestand von 25 Markern auf **null**; was
   übrig bleibt, ist eine frei in Prosa stehende, malformte Direktive — ein
   Autoren-Fehler, kein Doku-Nebeneffekt. Die Härte aus
   [ADR-0045](0045-zitat-verifikation-codepaths-erweiterung-und-citations-modul.md)
   trifft damit genau die Menge, für die sie gedacht war, und wird nicht
   gelockert. Die **Regel** ändert sich nicht, nur die Menge, die sie trifft —
   und das ist der Zweck der Änderung.

6. **Keine feine Ventil-Achse.** Weder `exempt-paths` noch `ignore-refs` noch
   ein Zeilen-Marker wird ergänzt, und die zwei groben Achsen bleiben, wie sie
   sind. Das Problem, das ein Ventil lösen sollte, verschwindet mit
   Entscheidung 1; eine Konfigurations-Fläche, die niemand braucht, ist eine,
   die jemand falsch benutzt
   ([ADR-0058](0058-konfigurations-flaechen-additiv-weiten.md)).

7. **Diese Entscheidung gilt der Zitat-Direktive, nicht der Ventil-Direktive —
   und das ist eine Skopierung, keine Auslassung.** Nach Entscheidung 1
   beantwortet das Produkt *„ist das eine Direktive"* ab jetzt **zweifach**: die
   vier Konsumenten des Ventils lesen weiter roh, eine Erwähnung des
   Ventil-Markers in Backticks **wirkt** dort. Das ist genau der Defekt aus
   [ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 1 und
   wird hier **nicht** behoben, weil die Richtung entgegengesetzt und die Folge
   unvermessen ist: `citations` wird **stiller**, das Ventil würde **lauter**.
   Gemessen sind die Ränder — **173** Prosa-Zeilen tragen den Ventil-Marker
   ausschließlich in Inline-Code (63 tragen ihn frei), und ein Lauf mit **ganz**
   abgeschaltetem Ventil meldet **58** Befunde. Der Angleich beträfe höchstens
   diese 58; wie viele wirklich, ist nicht gemessen und wird hier nicht
   behauptet. Der Angleich ist als eigener Slice geschnitten.

## Alternativen

- **Doku-Konvention statt Produkt-Änderung** — die Syntax nur noch in Fences
  schreiben. Verworfen, weil nicht durchführbar: 12 der 25 Marker liegen in
  vier eingefrorenen Review-Reporten und einer `Accepted`-ADR.
- **Eine feine Ventil-Achse für `citations`** (`exempt-paths`). Verworfen nach
  Entscheidung 6: sie benennt die betroffenen Dateien statt die Regel, und sie
  wüchse mit jedem neuen Review-Report weiter. Dasselbe gilt für die
  vorhandenen groben Achsen (`scan.ignore`, `citations.scope`): sie nehmen die
  ganze Datei aus dem Modul und schalteten es genau in den Dokumenten ab, die
  echte Zitate tragen.
- **Fail-closed zu einem Befund abschwächen.** Verworfen: das löste den
  Abbruch, nicht die Ursache — eine Erwähnung bliebe ein Befund, nur ein
  leiserer, und der Bestand trüge 24 davon.
- **Nur die Suche strippen, den Parse roh lassen.** Verworfen nach
  Entscheidung 3 — die Zeile mit beidem entschiede sonst gegen die falsche
  Quelle.
- **Beide Direktiven in einem Zug angleichen.** Verworfen nach Entscheidung 7:
  die eine Änderung verengt die geprüfte Menge und ist im Bestand folgenlos,
  die andere weitet sie um bis zu 58 neue Befunde. Sie in einen Slice zu legen
  hieße, die zweite an der Grünheit der ersten zu messen.

## Konsequenzen

- `citations` ist **aktivierbar**: der Lauf über den Bestand meldet 546
  Dateien, 0 Befunde, Exit 0. Das Scharfschalten selbst ist ein eigener Schritt.
- **Das Strippen kann eine Direktive auch erzeugen**, nicht nur verschlucken.
  Steht eine Code-Spanne zwischen Kommentar-Öffner und Marker, verschwindet sie
  und macht aus einer Nicht-Direktive eine — gemessen an einer Probe, die
  vorher 0 Befunde/Exit 0 lieferte und jetzt einen Befund/Exit 1 liefert. Das
  Verhalten ist als Test festgenagelt, nicht gewünscht; es folgt aus der
  Positionserhaltung und ließe sich nur mit einer zweiten, eigenen Lexik
  verhindern.
- **Eine echte Direktive kann still übersprungen werden — und die Klasse ist
  weiter als „steht in Backticks".** Es genügt, dass eine Code-Spanne
  **desselben Absatzes** sie umschließt: die Spanne öffnet eine Zeile vorher
  und schließt eine Zeile später, die Direktive selbst steht unverklammert.
  Gemessen an einem Dateipaar, das sich nur um zwei Backticks unterscheidet.
  Es ist derselbe Preis, den jedes prosa-lesende Modul zahlt — ein Link in
  Backticks wird ebenso wenig geprüft —, aber der Satz muss die
  Absatz-Reichweite nennen. Im Bestand betrifft er null Fundstellen.
- **Die Ziel-Seite bleibt außerhalb jeder Prosa-Zusage.** Das Modul liest die
  zitierte Quell-Spanne **roh und typunabhängig**: kein Fence-Bewusstsein, kein
  Strippen, und die Zieldatei muss weder Markdown sein noch in der Scan-Menge
  liegen — `scan.ignore` gilt der prüfenden Datei, nicht dem Ziel (§3.8: ein
  Modul verspricht nur über das, was es scannt). Das ist Absicht, weil dort
  Bytes verglichen werden, und gehört als Grenze in den Vertrag.
- Die Lexik-Änderung wirkt **nur** in diesem Modul: die geteilte Vorverarbeitung
  bleibt unverändert, und die drei gescopten Roh-Lesungen aus
  [ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) — `versions`,
  `pins`, `immutable` — sind nicht berührt.
- **[ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md)s
  Re-Evaluierungs-Trigger ist eingetreten** (*„eine vierte Stelle beantwortet
  eine Lexik-Frage selbst"* — hier derselbe Konsument mit einer anderen
  Lexik-Frage). Seine Folgefrage — ob ein **Gate** die Klasse prüfen kann statt
  eines Reviews — wird hier **nicht** beantwortet; sie ist eine DoD-Position
  von jenem Slice.

## Re-Evaluierungs-Trigger

**Permanent** für Entscheidung 2 (der Zitattext bleibt roh) — sie folgt aus dem
Gegenstand des Vergleichs, nicht aus einer Lage, die sich ändern kann.

**Wiedervorlage für Entscheidung 1**, wenn eine **produktive** Direktive von
einer Code-Spanne verschluckt wird: dann ist das stille Überspringen kein
Nullposten mehr, und die Frage nach einer Warnung — nicht nach einer Rücknahme
— stellt sich neu.

**Wiedervorlage für Entscheidung 7** mit der Closure jenes Slice: entweder der
Angleich fällt, oder die Zweifach-Antwort wird zur dauerhaft benannten Grenze.
Beides ist eine Entscheidung; das Weiterbestehen ohne Entscheidung ist keine.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-08-27 | Accepted | [slice-158](../planning/done/slice-158-citations-inline-code.md) |
