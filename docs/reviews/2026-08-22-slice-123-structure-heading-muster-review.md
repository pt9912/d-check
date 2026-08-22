# Review-Report: slice-123 — `structure`: jede Überschrift des Abschnitts, positiv geprüft

**Datum:** 2026-08-22 · **Review-Art:** Code- und Design-Review (Lastenheft/Spezifikation/ADR gegen Implementierung), unabhängiger Reviewer ohne Anteil an der Arbeit
**Gegenstand:** die vier Slice-Commits der Kette (Beanspruchung · Lastenheft-CR 0.64.0 · ADR-0058 Entscheidung 2 · Implementierung samt Spezifikation und Profil-Umstellung); Prüfbaum per `git archive` außerhalb des Repos, am Repo hat der Reviewer nichts geändert
**Skill:** `.harness/skills/reviewer.md` @ 1.8.0 · **Modell-ID:** `claude-opus-5[1m]`
**Vom Reviewer selbst gefahren (Exit explizit gelesen):** `make build` 0 · `make test` 0 · `make doc-check` 0 · `make gates` 0 · `make adr-check` 0 · `make trace-check` 0

**Verdikt des Reviews: blockierend** — ein HIGH, sechs MEDIUM, vier LOW, vier INFO. Alle eingearbeitet.

## Befunde und Einarbeitung

### F-1 · HIGH — der neue Config-Kommentar trug Chronik

`.d-check.yml` — der in diesem Slice **neu geschriebene** Kommentar erzählte,
was einmal war und was nicht mehr da ist: „*die abgeloeste Praefix-Negation
musste das und war dabei einmal zu eng*" (Review-Historie über eine entfernte
Konstruktion) und „*Grenze 1 entfaellt mit dem eigenen Schluessel*" (ein
Grabstein — der Kommentar beschreibt eine Grenze, die es nicht mehr gibt). Die
Bestands-Ausnahme greift nicht: die Zeilen sind neu.

**Eingearbeitet:** auf die Zusage reduziert — geprüft wird der
Überschriften-**Text** jeder Sektion, die Ebene ist der Default. Die
weggefallene Grenze ist ersatzlos gestrichen, die verbleibenden sind
neu durchnummeriert.

### F-2 · MEDIUM — das Lastenheft widersprach sich an drei Stellen

Rang 1 sagte dreimal „**als einzige**": einmal über die neue Bedingung, einmal
(bestehend) über die Chronologie-Bedingung, einmal in der §7-Historie. Zwei
„als einzige"-Aussagen über dieselbe Eigenschaft können nicht beide wahr sein
— und nach Source Precedence hätte die falsche gewonnen. Spezifikation und ADR
sagten korrekt „die zweite".

**Eingearbeitet:** alle drei Stellen auf „zwei benannte Ausnahmen" gezogen,
einschließlich der bestehenden Chronologie-Formulierung.

### F-3 · MEDIUM — dieselbe Aussage im selben Dokument vergessen

Schritt 5 der Spezifikation sagte weiter „mit **einer** benannten Ausnahme",
während Schritt 6 desselben Dokuments „die zweite" sagte. Die ADR behauptete,
die Aussage sei „an **beiden** Ausnahmen benannt" — sie war es an einer.

**Eingearbeitet:** Schritt 5 nennt jetzt beide.

### F-4 · MEDIUM — `Schärft:` nannte die Ziele der zweiten Entscheidung nicht

Das Feld führte nur `§DC-FA-VER-001.a` und `SPEC-005` aus Entscheidung 1.
Entscheidung 2 schärft `§DC-FA-STRUCT-001.a` und **erzeugt** `SPEC-067`.
Nach dem Übergang auf `Accepted` schließt `make adr-check` das Feld ab — die
Frist lief.

**Eingearbeitet:** beide Ziele ergänzt, solange die ADR `Proposed` ist.

### F-5 · MEDIUM — ein Schlüsselname mit zwei Rollen

`heading-pattern` existierte nach diesem Slice **zweimal** mit gegensätzlicher
Bedeutung: unter `planning.closure` als **Selektor** („welche Überschrift
eröffnet den geprüften Abschnitt"), unter `structure[]` als **Bedingung**
(„welche Form müssen die Überschriften haben"). Beide Blöcke leben in diesem
Repo im **selben** Profil, und die Spezifikation weist die Closure-Fähigkeit
als *Preset derselben Semantik* aus — Namensgleichheit läse sich dort als
Bedeutungsgleichheit. Wer den Selektor-Wert überträgt, prüft `'^#{2,3} …'`
gegen den Überschriften-**Text**, der nie ein `#` enthält: jede
Unterüberschrift rot, aus einem Grund, den die Meldung nicht erklärt.

**Eingearbeitet — umbenannt, nicht dokumentiert:** der Schlüssel heißt
`headings-match`/`headings-level`. Der Plural sagt die Semantik (**alle**
Überschriften, nicht eine), und der Name kollidiert nicht mehr. Nichts ist
released — die Umbenennung kostet hier genau nichts und ist danach nie wieder
so billig. Die Abgrenzung steht zusätzlich in der Schema-Zeile und in der ADR.

### F-6 · MEDIUM — die Botschaft verallgemeinerte über die Messung hinaus

Die acht behaupteten Formen hat der Reviewer **nachgemessen und bestätigt**.
Nicht haltbar war der Schluss daraus: „*Die Umstellung ist damit
verhaltenserhaltend, nicht heilend.*" Der Reviewer fand die neunte Form —
eine Überschrift **innerhalb eines mehrzeiligen Inline-Code-Spans**. Die
Bereinigung räumt diese Zeile leer, die abgelöste Negation sah dort nichts;
die Überschriften-Erkennung sieht die Überschrift, wie `anchors` auch, das ihr
einen Slug gibt.

Selbst nachgemessen, isoliertes Fixture, beide Konstruktionen im selben Image:

```
ALT (Negation):        1 Datei, 0 Befunde, Exit 0
NEU (headings-match):  docs/a.md:6  section-heading-mismatch, Exit 1
```

**Die Umstellung heilt also eine reale stille Klasse.** Der Commit ist
geschrieben und wird nicht umgeschrieben; die Korrektur steht hier, in der
Closure-Notiz und im Register.

**Register-Folge:** das ist die dritte Instanz von `BEO-009`. Die Klasse ist um
diese zweite Richtung erweitert (Botschaft verallgemeinert über die Messung
hinaus), der Zähler steht auf **3**, und weil die benannte mechanische Form nur
die erste Richtung deckt, ist die zweite im Reviewer-Skill als MEDIUM-Anker
verkörpert (1.9.0).

### F-7 · MEDIUM — nichtdeterministische Fehlermeldung

Die Musterprüfung des Config-Randes iterierte über ein **Map-Literal**. Bei
zwei ungültigen Mustern in derselben Regel entschied die Iterationsreihenfolge,
welche Meldung der Nutzer sieht — gemessen 6× / 2× über acht Läufe derselben
Konfiguration. `DC-QA-02` sagt das unqualifiziert zu. Die Klasse war
vorbestehend (drei Einträge), dieser Slice verbreiterte sie auf vier.

**Eingearbeitet:** geordneter Slice statt Map, mit einem Test, der zwanzig
Läufe auf Wortgleichheit prüft.

### F-8 · LOW — zwei Kommentare behaupteten dieselbe überholte Invariante

`structure.go` („die siebte Bedingung liest **als einzige** die rohen Zeilen")
und `sections.go` („**jede** Bedingung beider Konsumenten arbeitet auf diesem
Text") — beide seit ADR-0057 falsch, jetzt doppelt. **Eingearbeitet:** beide
auf zwei benannte Ausnahmen gezogen.

### F-9 · LOW — der Slice-Plan trug eine widerlegte Annahme

§1, §2.5 und der DoD-Haken behaupteten, die Negation sei **heute** still bei
eingerückten Sektionen. Der Reviewer hat die Gegenbehauptung des Autors belegt:
`db6a10c` (slice-114-Review) hatte das Muster auf die Modul-Lexik korrigiert.
**Eingearbeitet:** der Closure-Body weist §2.5/§4 als widerlegt aus — und
erfüllt den DoD-Sinn trotzdem, mit der Gegenprobe aus F-6.

### F-10 · LOW — falsche Schritt-Nummer

Der Bereich der Überschriften verwies auf „Schritt 4" (Kardinalität); die
Abschnitts-Grenze definiert Schritt 5. **Korrigiert.**

### F-11 · LOW — eine Grenzen-Begründung, die nicht mehr trägt

Der Config-Kommentar begründete die Tabellen-Grenze mit „die Bereinigung leert
Inline-Code" — für diese Regel gilt das nicht mehr, sie liest den bereinigten
Text gar nicht. **Eingearbeitet:** der wahre Grund steht da (sie sieht nur
Überschriften).

### F-12 · LOW — Herkunfts-Prosa im Code-Kommentar

„*eine nachgebaute Heading-Lexik war der Anlass dieser Bedingung*" — Herkunft
ohne auflösbares Feld, dieselbe Klasse wie F-1, schwächer. **Entfernt**, die
Kopplungs-Aussage davor trägt.

### INFO — festgehalten

- **INFO-1:** `SPEC-067` steht in §4 zwischen `SPEC-058` und `SPEC-059`
  (Modul-Gruppierung als Ordnungsprinzip); die Kennungs-Spalte ist damit die
  einzige Bruchstelle einer sonst aufsteigenden Spalte. Bewusst so, kein Gate
  sieht es.
- **INFO-2:** das Handbuch ist jetzt an drei Stellen widersprüchlich („die
  **eine** Ausnahme", „bis zu **sieben** Bedingungen", Config-Beispiel ohne die
  neuen Schlüssel). Ausdrücklich out-of-scope — benannt für den
  Release-Prep-Slice.
- **INFO-3:** zwei Ebenen in **einer** Regel sind nicht bloß unvorgesehen,
  sondern durch die Regel-Identität **gesperrt** (Exit 2, gemessen).
  **Eingearbeitet:** als dritte Grenze in die ADR.
- **INFO-4:** doppelte Leerzeile im Test. **Entfernt.**

## Negativbefunde des Reviews (geprüft, ohne Befund)

- **Alle sieben Akzeptanzkriterien** einzeln gegen Code **und** Test geprüft —
  jedes hat einen Test, und jeder Test misst, was er behauptet.
- **Vier Off-by-one-Proben** an den Bereichsgrenzen (Überschrift direkt nach
  dem Kopf, als letzte Datei-Zeile, direkt vor dem Terminator, der Terminator
  selbst) — alle korrekt.
- **Der Fence-Verdacht ist widerlegt:** der Automat startet bei `headingNo`
  mit `inFence=false`, und das ist beweisbar richtig, weil der Abschnittskopf
  aus einer Erkennung stammt, die Überschriften **innerhalb** einer Fence
  verwirft. Gegenproben gemessen (unbalancierte Fence vor und im Abschnitt).
- **Ebenen-Semantik bei `section-pattern`/`sections: each`:** die Ebene wird
  **je Abschnitt** ausgewertet, nicht je Regel — mit einem Muster gemessen,
  das Ebene 2 und 3 trifft.
- **`DC-QA-02` gemessen, nicht geglaubt:** ganzer Baum mit alter Konfiguration
  gegen das gepinnte Vorgänger-Image byte-identisch; roter Fixture ohne den
  neuen Schlüssel byte-identisch; Profil-Umstellung byte-identisch;
  Gegenprobe `field heading-pattern not found`, Exit 2.
- **Die kritische Richtung der Deckung ist zu:** über **22** konstruierte
  Formen fand der Reviewer **keinen** Fall, in dem der neue Schlüssel still
  bleibt und die Negation eine **echte** Überschrift gemeldet hätte — mit
  Strukturargument: die Bereinigung ersetzt Zeichen nur durch Leerzeichen, sie
  kann ein `SPEC-NNN␣`-Präfix nie erzeugen.
- **Spiegel (MR-025)** einzeln abgehakt, zwölf Flächen; fehlend waren nur die
  in F-2/F-3/F-8 genannten.
- **Fail-closed vollständig**, `Re-Evaluierungs-Trigger` vorhanden (die
  slice-122-Lehre ist eingearbeitet), Lifecycle-Commit korrekt nach MR-013 —
  und die BEO-006-Falle des Vorgängers vermieden.

## Nachmessung nach der Einarbeitung

- `make gates` grün (acht Gates, Exit 0 explizit gelesen).
- Der Heilungs-Fall aus F-6 **nach** der Umbenennung erneut gemessen: alte
  Konstruktion Exit 0, neue Exit 1 auf der Zeile der Überschrift.
- Die Profil-Umstellung erneut gegen den Stand **vor** ihr: byte-identisch.
