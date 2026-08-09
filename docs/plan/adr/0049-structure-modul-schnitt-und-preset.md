# ADR-0049 — Modul `structure`: Schnitt, Grenze und Preset-Kopplung statt Supersede

**Status:** Proposed
**Datum:** 2026-08-09
**Autor:** pt9912
**Bezug:** [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(neu),
[`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(Mit-Modifikation: Preset-Ausweisung + `closure-note-ambiguous`);
Schnitt-Kriterium [ADR-0044](0044-geteiltes-referenz-ventil-quell-skopus.md);
Vorläufer der Fähigkeit [ADR-0048](0048-closure-note-struktur-im-planning-modul.md)
(**nicht** superseded); Modul-Fundament
[ADR-0005](0005-modul-layout-hexagon-ordner.md),
[ADR-0012](0012-kern-paketschnitt-model-rules-app.md); Nullmengen-Logik
[ADR-0037](0037-trace-tabellenquellen-nullmengen-guard.md).
**Schärft:** die neue Algorithmus-Sektion
[`spec/spezifikation.md` §DC-FA-STRUCT-001.a](../../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
und die Preset-Kopplung in
[§DC-FA-PLAN-001.a](../../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning).

## Kontext

Ein Adopter hat seine vier handgeschriebenen Prüfskripte gegen d-checks Modulsatz
vermessen und dabei eine Grenze benannt, die nie ausgesprochen war: der Modulsatz
deckt **Referenz**-Invarianten lückenlos ab (Ziel existiert, Kennung verlinkt,
Richtung erlaubt, Target deklariert, Core unverändert) — aber **keine einzige
Aussage über die Form eines Dokuments selbst**. Die Module sind entlang „zeigt
dieses Dokument korrekt auf andere?" gewachsen, nie entlang „ist dieses Dokument
selbst richtig gebaut?".

Die Ursache, warum das erst jetzt auffiel, ist selbst ein Befund: jeder Adopter
füllt die Lücke lokal mit einem Skript, und ein lokales Skript sieht wie
Repo-Eigenheit aus, nicht wie fehlende Werkzeug-Fähigkeit. **Kein Gate kann
melden, dass ein anderes Gate hätte existieren können.**

Eigene Messung (2026-08-09) gegen die drei einschlägigen Skripte des Adopters
(480 Zeilen Shell): sie tragen **elf** Prüfungen — 2 heute gedeckt, 2 nach
Kalibrierung, **6 ungedeckt**, 1 außerhalb, 1 als Ventil ausdrückbar.

Erschwerend: d-check hat wenige Tage zuvor mit
[ADR-0048](0048-closure-note-struktur-im-planning-modul.md) eine
planning-lokale Closure-Note-Struktur ausgeliefert (v0.52.0), die der
**Spezialfall** genau dieser Abstraktion ist.

## Entscheidung

1. **Neues Modul `structure`, Bereichskürzel `STRUCT`** — als **Liste** von
   Regeln über Datei-Globs, nicht als Ausbau von `planning.closure`.
   Das Kriterium ist nicht neu gewählt, sondern
   [ADR-0044](0044-geteiltes-referenz-ventil-quell-skopus.md) entnommen:
   querschnittlich ⇒ eigenes Kürzel, Einzelmodul ⇒ bestehende Anforderung
   ändern. Entschieden hat es eine einzige Messzeile: die Pflicht-Bausteine
   einer **Anforderung** im Lastenheft sind dieselbe Frageform wie die Substanz
   einer Closure-Notiz, aber eine **andere Dokumentklasse**.
   Die Alternative reicht **technisch** (der Kandidaten-Filter ließe sich
   erweitern) und scheitert am **Namen**: ein Modul namens `planning`, das das
   Lastenheft prüft, täuscht über seinen Gegenstand.

2. **Die Modul-Grenze ist „innerhalb eines Dokuments" — und sie steht im
   Vertrag.** Dieselbe Skript-Familie prüft auch, ob ein Dateiname eine Kennung
   hergibt. Das ist eine Aussage über den **Ort**, nicht über den Text, und
   damit ausdrückliches Nicht-Ziel. Ohne diese Grenze wird „Struktur-Invarianten"
   eine Kategorie statt einer Prüfung, und das Modul zum Sammelbecken.

3. **Nichts wird superseded — die Closure-Fähigkeit wird zum Preset.** Die
   naheliegende Bauform (alte Anforderung superseden, Config-Schlüssel als Alias)
   ist **ausgeschlossen**, nicht abgewogen: die Spezifikation führt die
   Grund-Codes als „stabil, maschinenlesbar" — im ausdrücklichen Unterschied zur
   `message`, die nicht stabilitätsgarantiert ist. `closure-note-*` gegen
   `section-*` zu tauschen bräche eine zugesagte Fläche; jede Konsumenten-CI, die
   auf den Code filtert, bricht mit.
   Also: `structure` entsteht **neben** der Fähigkeit mit eigenen Codes,
   [`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
   und [ADR-0048](0048-closure-note-struktur-im-planning-modul.md) bleiben
   stehen. Zusammengeführt wird die **Semantik**: die Closure-Fähigkeit ist
   spezifiziert als Preset derselben Schritte (Abschnitts-Bestimmung,
   Fence-Bereinigung, Zählung). Doppelt ist die Config-Oberfläche, nicht die
   Mechanik.
   Das ist die Form aus [ADR-0044](0044-geteiltes-referenz-ventil-quell-skopus.md)
   aus dem **umgekehrten** Grund: dort blieb der alte Schlüssel Alias, *weil* die
   Codes gleich waren; hier bleibt die alte Anforderung ganz stehen, *weil* sie
   verschieden sind.

4. **Mehrdeutigkeit ist ein eigener Befund, in beiden Oberflächen**
   (`section-ambiguous`, additiv `closure-note-ambiguous`). Der
   Aktiv-Status-Guard desselben Moduls ist an dieser Stelle längst fail-closed;
   dass die Closure-Fähigkeit still den **ersten** Treffer nahm, war ein
   Versäumnis — und ein stilles, denn ein zweiter, stehengebliebener Abschnitt
   ist der typische Rest einer Vorlage. Ohne eindeutigen Abschnitt wird **nicht
   gemessen**: eine Satzzahl über dem falschen Abschnitt ist keine Aussage.
   Eigener Code statt Sammelbefund, nach der Begründung aus
   [ADR-0048](0048-closure-note-struktur-im-planning-modul.md): die Reparatur ist
   „den überzähligen Abschnitt entfernen", nicht „den fehlenden schreiben".

5. **Zwei Marken-Formen, nicht eine — und zeilenverankert.** Der Antrag nennt
   eine Liste, die vollständig vorkommen muss. Die Messung zeigt beide Bedarfe:
   eine Anforderung muss **alle** Akzeptanz-Bausteine tragen, ein Lerneintrag
   genügt mit **einer** von mehreren zulässigen Formen. Eine Liste deckt das
   nicht ab, also `require-all` **und** `require-any`.
   Und die Marke ist eine **hervorgehobene Zeilen-Anfangs-Marke**, kein
   Teilstring: ein `Boundary` mitten im Satz erfüllt die Zusage nicht. Ein
   Teilstring-Vergleich erzeugte ein Falsch-Grün — teurer als ein Falsch-Rot,
   weil es unbemerkt bleibt.

6. **Keine Stichtags-Mechanik; Pfad-Ausnahmen genügen.** Der gelebte
   Grandfathering-Fall des Adopters ist eine Zahlen-Schwelle im Dateinamen. Sie
   zu lernen hieße, dessen Kennungs-Konvention zu interpretieren — d-check
   müsste aus einem Dateinamen eine Ordnung lesen. Dieselbe Grenze wie die
   Index-Wahrheit in
   [`DC-FA-TRK-001`](../../../spec/lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in):
   **kein zweiter Regel-Interpreter im Werkzeug.** `exempt-paths` je Regel deckt
   den Pfad-Fall; die Grenze gehört benannt, damit der nächste Antrag sie nicht
   neu verhandelt.

## Verglichene Alternativen

| Alternative | Warum verworfen |
|---|---|
| `planning.closure` um einen Datei-Glob erweitern | Reicht technisch, täuscht aber über den Gegenstand: `planning` prüft den Planning-Lifecycle, nicht das Lastenheft |
| Alte Anforderung superseden, Schlüssel als Alias | Bräche die stabil zugesagten Grund-Codes — ausgeschlossen, nicht abgewogen |
| `structure` emittiert je nach Config-Pfad mal `section-*`, mal `closure-note-*` | Ein Modul mit zwei Code-Familien je nach Herkunft der Konfiguration; die Meldung wäre nicht mehr aus dem Modul erklärbar |
| Mehrdeutigkeit als `section-constraint` mitführen | Andere Reparatur, andere Klasse; Sammelbefunde waren schon bei den drei Closure-Codes die verworfene Bauform |
| `require-strong` als einzelne Liste (wie beantragt) | Deckt „mindestens eine von" nicht ab; der Adopter behielte ein Skript für genau einen Fall |
| Marken als Teilstring prüfen | Falsch-Grün: ein Wort im Fließtext erfüllte eine Gliederungs-Zusage |
| Stichtags-Schwelle als eigener Mechanismus | Zweiter Regel-Interpreter im Werkzeug; dieselbe Grenze, die `tracked` schon gezogen hat |

## Konsequenzen

- **Die Identität des Werkzeugs wächst**, und zwar erklärt: `structure` ist das
  erste Modul ohne Referenz-Invariante. §1 des Lastenhefts nimmt die zweite
  Frage ausdrücklich auf. Ohne diesen Satz wäre die Erweiterung still — genau
  die Klasse, die den Antrag ausgelöst hat.
- **Zwei Config-Wege für verwandte Fragen.** Der Preis der Code-Stabilität. Die
  Spezifikation koppelt sie, damit sie nicht driften: eine Änderung an einer der
  beiden Stellen ohne die andere ist ein Spezifikations-Bug.
- **`closure-note-ambiguous` ist additiv** ⇒ SemVer-Minor. Ein Repo mit zwei
  Closure-Abschnitten wird danach rot, wo es vorher still den ersten las. Das
  gehört in die Release-Notiz.
- **Der Paritäts-Anspruch wird erreichbar:** mit dem Mehrdeutigkeits-Code kann
  die bewusst verengte Zusage des Zähl-Slice wieder auf volle Deckung gehen.
- **Die Fixtures liegen im Adopter-Repo**, nicht hier. Sie werden beigezogen,
  nicht nachgebaut — ein nachgebautes Fixture belegt die Parität nicht.

## Fitness Function

- **Preset-Kopplung:** ein Test fährt dieselbe Eingabe durch beide Oberflächen
  und vergleicht die Befund-**Positionen**; divergieren sie, sind die Semantiken
  auseinandergelaufen.
- **Mehrdeutigkeit schlägt Messung:** ein Dokument mit zwei passenden
  Abschnitten und einem zu dünnen ersten meldet **nur** den Mehrdeutigkeits-Code.
- **Marken-Verankerung:** ein Dokument, das eine Marke nur als Fließtext-Wort
  trägt, meldet die Bedingung als verletzt.
- **Leerlauf:** eine Regel, deren Glob keine Datei trifft, meldet — statt grün zu
  sein.
- **Byte-Identität:** ohne aktives `structure` bzw. ohne Regeln ist der
  Befundsatz unverändert.

## Re-Evaluierungs-Trigger

- Wenn eine dritte Oberfläche dieselbe Struktur-Semantik brauchte, ist die
  Preset-Kopplung an ihrer Grenze — dann ist die Frage „ein Mechanismus, zwei
  Namen" neu zu stellen.
- Wenn Adopter wiederholt Stichtags-Ausnahmen als Glob-Ketten schreiben müssen,
  ist Entscheidung 6 gegen die gelebte Praxis zu prüfen — die Antwort wäre dann
  eine eigene Anforderung, nicht eine Ausweitung dieser.
- Wenn `structure`-Regeln beginnen, Aussagen über Datei**orte** zu schmuggeln
  (etwa über `files`-Globs, die eine Konvention erzwingen sollen), ist die
  Grenze aus Entscheidung 2 zu schärfen.

## Geschichte

- 2026-08-09: Proposed (doc-first, `slice-096`).
