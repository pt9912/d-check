# ADR-0063 — Die Form des Ventil-Markers folgt der Kommentar-Lexik seiner Eingabe

**Status:** Accepted
**Datum:** 2026-08-27
**Autor:** pt9912
**Schärft:** [`DC-FA-CODE-001.a`](../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code)
Schritt 1 und [`DC-FA-ID-001.a`](../../../spec/spezifikation.md#dc-fa-id-001a--kennungs-prüfung)
Bedingung 4 (je die **Form** des Zeilen-Markers)
**Bezug:**
[ADR-0062](0062-ventil-marker-versions-ist-eine-benannte-grenze.md) (dieselbe
Frage für die **Lage** des Markers — diese ADR ist ihre zweite Hälfte),
[ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) (die Trennlinie
*andere Antwort* gegen *andere Frage*),
[`DC-FA-DIAG-001.a`](../../../spec/spezifikation.md#dc-fa-diag-001a--kennungs-konsistenz-in-diagramm-fences-diagrams)
(die **festgelegte** Token-Form für `diagrams`)

## Kontext

[ADR-0062](0062-ventil-marker-versions-ist-eine-benannte-grenze.md) hat die
**Lage** des Markers geregelt — außerhalb von Inline-Code, wo die Eingabe Prosa
ist. Offen blieb seine **Form**: die Spezifikation nennt ihn an zwei Stellen
einen **HTML-Kommentar**, die Erkennung akzeptierte jedes Vorkommen der
Zeichenkette. Eine blanke Erwähnung in Prosa — *„der Marker `d-check:ignore`
nimmt die Zeile aus"*, ohne Backticks — wirkte damit wie ein gesetzter Marker.

**Eine repo-weit einheitliche Kommentar-Form ist nicht verfügbar, und das ist
kein Aufwands-Argument.** Die Spezifikation legt für `diagrams` **ausdrücklich
und begründet** den Token fest:

> Der Marker ist ein **Token**, kein HTML-Kommentar. Das Modul sucht die
> Zeichenfolge auf der Roh-Zeile; in einem `mermaid`-Fence versteckt ihn die
> Diagramm-Sprache (`%% d-check:ignore`), nicht Markdown. Eine Kommentar-Lexik
> je Fence-Sprache wäre ein Grammatik-Parser.

Der Slice, der diese ADR trägt, hat den Widerspruch **gemeldet statt aufgelöst**
([`AGENTS.md`](../../../AGENTS.md) §1) — richtig so: die Spec sticht den Plan.
Dieselbe Überlegung trifft `versions`: es liest **alle** Zeilen, auch die in
Fences fremder Sprachen. Es gibt dort keine **eine** Kommentar-Lexik, die zu
fordern wäre.

**Für `codepaths` und `ids` gibt es sie.** Ihre Eingabe sind Markdown-Prosa-
Zeilen, und Markdowns Kommentar-Lexik ist der HTML-Kommentar. Die Verengung ist
dort verfügbar — **und im Bestand kostenlos**: gemessen über 557 Dateien meldet
der Repo-Lauf mit geforderter Kommentar-Form **null** zusätzliche Befunde. Jeder
heute wirksame Marker dieser zwei Konsumenten trägt die Form bereits.

**Dass die Verengung trotzdem Zähne hat, ist eigens belegt** — Grün allein wäre
kein Ergebnis. Gegenprobe: eine Zeile mit blankem Marker und nicht auflösbarem
Pfad wird mit der Verengung gemeldet und ohne sie nicht; die Zeile mit
Kommentar-Marker bleibt in beiden Fassungen still.

## Entscheidung

1. **Die Form folgt der Kommentar-Lexik der Eingabe.** Sie ist damit **je
   Konsument** verschieden, und das ist keine zweite Antwort auf dieselbe Frage,
   sondern die Anwendung derselben Regel auf verschiedene Eingaben:

   | Konsument | Eingabe | Kommentar-Lexik | geforderte Form |
   |---|---|---|---|
   | `codepaths`, `ids` | Markdown-Prosa | HTML-Kommentar | **`<!-- … -->`** |
   | `diagrams` | Fence-Inneres + Öffnungszeile | die der Fence-Sprache | **Token** (festgelegt) |
   | `versions` | alle Zeilen, sprachgemischt | keine einheitliche | **Token** |

2. **Bei `codepaths` und `ids` muss der Marker in einem HTML-Kommentar
   stehen.** Geprüft wird, dass er nach `<!--` und vor dessen Ende steht.

3. **Die Bedingung ist konservativ formuliert.** Ein `>` im Kommentar vor dem
   Marker lässt ihn **nicht** gelten. Ein verpasster Marker ist Falsch-Rot —
   laut und sichtbar; ein erfundener wäre stilles Grün. Die Richtung des
   Irrtums ist gewählt, nicht zufällig.

4. **`diagrams` und `versions` bleiben unangetastet.** Für `diagrams` ist das
   die Spec-Festlegung, für `versions` folgt es aus derselben Überlegung. Beide
   behalten damit auch die **Lage**-Regel aus
   [ADR-0062](0062-ventil-marker-versions-ist-eine-benannte-grenze.md), also die
   Roh-Lesung.

5. **Die Uneinheitlichkeit ist der Preis und wird ausgewiesen.** Wer den Marker
   setzt, muss wissen, für welches Modul er ihn setzt. Das stand vorher nicht
   im Vertrag; jetzt steht es dort als Tabelle.

## Alternativen

- **Eine Form für alle vier.** Verworfen: nicht verfügbar. Sie hieße, die
  ausdrückliche Spec-Festlegung für `diagrams` per ADR abzulösen und dort einen
  Grammatik-Parser je Fence-Sprache zu verlangen — oder umgekehrt, für
  `codepaths`/`ids` beim Token zu bleiben und die Spec-Aussage *„HTML-Kommentar"*
  an zwei Stellen zu streichen. Die erste ist teuer und falsch, die zweite gäbe
  die einzige verfügbare Verengung auf.
- **Gar nichts verengen und die Form als ungeregelt führen.** Verworfen: der
  Vertrag sagt an zwei Stellen etwas anderes als der Code, und die Verengung ist
  im Bestand kostenlos. Eine benannte Grenze verwaltet den Widerspruch, sie
  behebt ihn nicht.
- **Die Bedingung großzügiger fassen** (Marker irgendwo auf einer Zeile, die
  auch einen Kommentar trägt). Verworfen nach Entscheidung 3 — das ist die
  Richtung zum stillen Grün.

## Konsequenzen

- Eine **blanke** Erwähnung des Markers in Prosa wirkt bei `codepaths` und `ids`
  nicht mehr. Zusammen mit
  [ADR-0062](0062-ventil-marker-versions-ist-eine-benannte-grenze.md) ist damit
  jede Erwähnungs-Form dieser zwei Konsumenten entschärft — die in Backticks
  über die **Lage**, die blanke über die **Form**.
- **Die Form ist jetzt je Modul verschieden.** Wer einen Marker setzt, ohne zu
  wissen, welches Modul er stumm schalten will, setzt ihn womöglich in der
  falschen Form. Das ist der ausgewiesene Preis.
- Die zwei Ränder des Strippens aus
  [ADR-0062](0062-ventil-marker-versions-ist-eine-benannte-grenze.md) gelten
  **unverändert** weiter — die Form-Bedingung liegt hinter der Lage-Bedingung
  und ändert an ihr nichts.
- `versions` behält beide Roh-Eigenschaften und bleibt damit die **benannte
  Grenze** aus
  [ADR-0062](0062-ventil-marker-versions-ist-eine-benannte-grenze.md)
  Entscheidung 3.

## Re-Evaluierungs-Trigger

**Permanent** für Entscheidung 1 — sie folgt aus der Eingabe-Art.

**Wiedervorlage für Entscheidung 3**, wenn ein legitimer Marker an der
`>`-Bedingung scheitert: dann ist die konservative Form zu weit gefasst und
braucht die vollständige Kommentar-Grammatik.

**Wiedervorlage für Entscheidung 4** gemeinsam mit
[ADR-0062](0062-ventil-marker-versions-ist-eine-benannte-grenze.md)s
`versions`-Grenze — Form und Lage gehören dort in eine Entscheidung.
