# ADR-0079: Die Register-Deckung zählt Linktext, nicht nur Prosa

**Status:** Accepted

**Datum:** 2026-08-31

**Autor:** pt9912

**Bezug:**
[`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(die erweiterte Anforderung),
[ADR-0028](0028-planning-lifecycle-modul.md) (das Modul und seine Bauform
„opt-in innerhalb des opt-in Moduls"),
[`DC-FA-ID-001`](../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
(dessen Linktext-Begriff hier wiederverwendet statt nachgebaut wird),
[ADR-0060](0060-citations-marker-scan-geteilte-prosa-antwort.md) (der
Präzedenzfall „Inline-Code zählt nicht", der hier **nicht** trägt)

**Schärft:**
[`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)

---

## Kontext

Das Baseline-Regelwerk nennt die Register-Paarung selbst als **maschinell
entscheidbar**: eine zitierte Beobachtungs-Kennung hat eine Zeile im Register.
Gebaut hatte sie niemand — die Deckung war gelebte Praxis ohne Sensor, und ein
erfundenes `BEO-999` in einer Closure-Notiz wäre niemandem aufgefallen.

Beim Bau stellte sich heraus, dass die naheliegende Erkennungs-Regel dieses
Repos hier **das Gegenteil** des Gewollten tut. `citations` und die
Closure-Platzhalter-Bedingung zählen Vorkommen nur **außerhalb** von
Inline-Code — mit guter Begründung: dort stehen die Illustrationen.

Für Kennungen ist die Lage umgekehrt, und das ist gemessen, nicht vermutet.
Im Planning-Baum stehen:

| Form | Anzahl |
|---|---|
| Link **mit** Backticks im Linktext — die verbreitete Zitier-Form | 366 |
| Link ohne Backticks | 5 |
| Prosa ohne Backtick | 293 |
| reines Inline-Code-Span, kein Link | 112 |

Die beiden einzigen Kennungen ohne Registerzeile — `BEO-999` in einem
Slice-Plan über diese Regel, `BEO-099` in einem Review-Report — sind **reine
Code-Spans**. Eine Regel „Inline-Code zählt nicht" hätte also die 366 echten
Zitate übersehen und ausschließlich die Beispiele gesehen.

## Entscheidung

**Gezählt wird ein Vorkommen in Prosa und ein Vorkommen im Linktext; ein
reines Inline-Code-Span zählt nicht.** Die Trennlinie verläuft zwischen
**Behauptung** und **Beispiel**, nicht zwischen Prosa und Code.

Die Linktext-Erkennung wird **wiederverwendet, nicht nachgebaut**: `ids` trifft
dieselbe Unterscheidung für seine Linkpflicht (*ein Code-Span-Vorkommen ist
linkpflichtfrei, wenn es im Linktext liegt*). Zwei Antworten auf dieselbe Frage
wären ein Defekt.

**Deklarierend ist nur die erste Zelle einer Tabellenzeile**, und sie muss
**ganz** aus der Kennung bestehen. Andernfalls deklarierte sich jede
Quer-Referenz im Fließtext einer Registerzelle selbst — der Wächter fände nie
etwas.

**Träger ist eine vierte Fähigkeit im vorhandenen Modul `planning`, kein
eigenes Modul.** `planning` besitzt Layout, hermetischen Vertrag und die
Bauform bereits.

**Die Scan-Menge ist konfigurierbar, die Richtung nicht.** Geprüft wird
ausschließlich *Zitat ⇒ Zeile*.

## Verglichene Alternativen

**Kennung als `ids`-Linkpflicht auf einen Anker je Registerzeile** — die Form,
die `version.md` für Versionen fährt: jede Zeile bekommt `<a id="beo-NNN">`,
jede Nennung muss dorthin verlinken, und `anchors` fängt die erfundene Kennung
gratis. **Gemessen verworfen:** das bräuchte Anker in 22 Zeilen (klein),
`#beo-NNN` in 366 bestehende Links (groß) und Links für 293 Prosa-Nennungen
allein im Planning-Baum (groß). Repo-weit sind es 1450 Nennungen, davon 1042
nackt, allein 708 in `done/` und 626 in `docs/reviews/` — das Retrofit schriebe
**eingefrorene Lauf-Belege** um.

**Ein eigenes Modul `registry`** — die Form ist im Register selbst
vorgezeichnet (*ein Verzeichnis-Muster je Register plus Autoritäts-Datei, eine
Richtung, ein Grund-Code*) und deckte später den ADR-Index und den
Konventionsspeicher-Index mit. **Verworfen, weil die Rechtfertigung aus
künftigem Bedarf stammt** und nicht aus Bestand — genau die Klasse, die das
Register als wiederkehrend führt. Fällt die dritte Anwendung an, ist die
Verallgemeinerung dann mit Bestand zu begründen.

**Auch die Beleg-Form prüfen (Form · Anzahl · Lage)** — sie steht in derselben
Kanon-Stelle. **Verworfen für diesen Schnitt:** *Lage* deckt bereits `links`
(alle Beleg-Zellen sind `done/`-Links), *Form* ist erfüllt, und *Anzahl* ginge
auf dem Bestand rot — zwei Zeilen tragen eine **in der Zeile begründete**
Abweichung, und die Quelle des Kanons räumt ein, zwei Regeln zu schulden. Der
Weg dorthin ist eine benannte Bestands-Ausnahme, kein Carveout.

## Konsequenzen

**Der Wächter ist von Anfang an grün** — am eigenen Bestand 0 Befunde. Das ist
der Zweck: er hält eine Praxis, die heute gelebt wird und die niemand prüft.
Zugleich ist es die Lage, in der ein Sensor später für überflüssig gehalten
wird; die Anforderung sagt deshalb, **wovor** er schützt.

**Eine Kennung außerhalb der konfigurierten Verzeichnisse wird nicht geprüft.**
Das Modul verspricht nur über seine Scan-Menge; die Grenze steht in der
Anforderung, nicht nur im Code.

**Ein vertipptes Zitier-Verzeichnis meldet, statt still inert zu gehen.** Der
Rand entstand aus einem fehlgeschlagenen eigenen Test: der In-Memory-Adapter
liefert für ein fehlendes Verzeichnis eine **leere Liste** statt eines Fehlers.
Geprüft wird deshalb die Art des Pfades, nicht das Verhalten des Listings.

**Das Urteil bleibt beim Menschen.** Ob zwei Beobachtungen dieselbe sind,
entscheidet niemand maschinell; die Umkehrung der Paarung bleibt ungeprüft.

## Fitness Function (falls maschinell prüfbar)

Die Entscheidung prüft sich selbst: `TestObservationsCodeSpanIsExampleButLinkTextCounts`
hält alle vier Formen (Code-Span, Linktext mit und ohne Backticks, Prosa)
gegeneinander; `TestObservationsOnlyFirstCellDeclares` hält die Deklarations-Regel;
`TestObservationsFailClosedAndInert` die drei Ränder. Am Produkt gemessen:
eigener Bestand 0 Befunde, konstruierte Kennung als Link 1, dieselbe als
Code-Span 0.

## Re-Evaluierungs-Trigger

**Wenn eine dritte Anwendung derselben Bauform anfällt** — ADR-Index oder
Konventionsspeicher-Index brauchen dieselbe Deckung —, ist die
Modul-Verallgemeinerung neu zu bewerten, dann mit Bestand statt mit Vermutung.

**Wenn die Quelle des Kanons die zwei fehlenden Beleg-Regeln setzt** (Vorkommen
außerhalb einer Slice-Closure; zweites Vorkommen derselben Klasse im selben
Slice), ist die Anzahl-Achse fällig und diese ADR um sie zu ergänzen oder
abzulösen.

## Geschichte
