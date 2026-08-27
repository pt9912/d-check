# Slice slice-162: `d-check:ignore` beantwortet dieselbe Frage anders als `d-check:cite`

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
(Entscheidung 1: ein Konsument, der eine Lexik-Frage selbst beantwortet, ist ein
Defekt — und der Re-Evaluierungs-Trigger *„eine vierte Stelle"*);
[ADR-0060](../../adr/0060-citations-marker-scan-geteilte-prosa-antwort.md) (die
Skopierung, die diesen Slice schneidet);
[ADR-0025](../../adr/0025-codepaths-ignore-refs.md) (das Ventil selbst);
[slice-158](../done/slice-158-citations-inline-code.md) (der Anlass).

**Berührte Spec-Stellen:** [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in),
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids),
[`DC-FA-VER-001`](../../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
[`DC-FA-DIAG-001`](../../../../spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in)
— je die Ventil-Zusage.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-27.

---

## 1. Ziel

`d-check` kennt zwei Direktiven, und seit
[slice-158](../done/slice-158-citations-inline-code.md) beantworten sie
**dieselbe** Frage verschieden. *„Ist diese Zeile eine Direktive"* liest
`citations` auf dem inline-code-gestrippten Text; die vier Konsumenten des
Ventils lesen weiter roh:
[`codepaths.go`](../../../../internal/hexagon/core/rules/codepaths.go),
[`ids.go`](../../../../internal/hexagon/core/rules/ids.go),
[`versions.go`](../../../../internal/hexagon/core/rules/versions.go),
[`diagrams.go`](../../../../internal/hexagon/core/rules/diagrams.go). Eine
Erwähnung von `d-check:ignore` in Backticks **wirkt** dort als Ventil.

Das ist die Klasse aus
[ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
Entscheidung 1 — *zwei Antworten auf dieselbe Frage in einem Lauf sind ein
stiller Grün-Pfad, den kein Gate sieht* — und zugleich ihr eigener
Re-Evaluierungs-Trigger (*„eine vierte Stelle beantwortet eine Lexik-Frage
selbst"*).

**Gemessen** (Produkt-Lexik: Fence-Automat, absatzweise Spannen gleicher
Backtick-Länge, über 544 getrackte Markdown-Dateien):

| Lage von `d-check:ignore` | Zeilen |
|---|---|
| frei in Prosa (wirkt so oder so) | **63** |
| **nur** in Inline-Code (wirkt heute, würde nach Angleich nicht mehr) | **173** |

**Obergrenze der Sprengweite, gemessen:** wird das Ventil **ganz** abgeschaltet,
meldet der Repo-Lauf **58** Befunde (21 davon in
[`spec/lastenheft.md`](../../../../spec/lastenheft.md)). Der Angleich beträfe
nur die 173 Zeilen, also **höchstens** 58 — wie viele davon, ist die erste
Messung dieses Slice und **nicht** vorweggenommen.

## 2. Vorgehen

1. **Die Zahl unter der Obergrenze messen**, bevor entschieden wird: welche der
   58 Befunde hängen an einer Erwähnung in Inline-Code, welche an einem freien
   Marker? Ein Lauf mit gestripptem Ventil gegen einen mit rohem, Befundsätze
   verglichen.
2. **Die Richtung ist die unangenehme.** `citations` wurde **stiller**, das
   Ventil würde **lauter**: nach dem Angleich melden Zeilen, die heute
   schweigen. Jeder dieser Befunde ist einzeln zu beurteilen — echt (das Ventil
   stand nie legitim dort) oder Falsch-Rot (die Zeile braucht ein Ventil, nur
   ein anderes).
3. **Die Vertrags-Frage.** Vier Anforderungen tragen die Ventil-Zusage. Ist der
   Angleich eine Schärfung je Anforderung oder eine querschnittliche? Und trägt
   [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)s
   Trigger-Frage — *„kann ein **Gate** die Klasse prüfen statt eines Reviews"* —
   hier eine Antwort? Ein Wächter *„kein Konsument liest roh, wo die geteilte
   Antwort existiert"* wäre die strukturelle Reparatur statt der punktuellen.
4. Nur bauen, was die Messung trägt; die Entscheidung gegen einen Angleich wäre
   ebenso auszuweisen wie einer für ihn.
5. Bewusstes Brechen je gedeckter Form, **Ursache gelesen**; Rückbau grün.
6. `make gates`, `make fullbuild`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Änderung an `citations`.** Dessen Antwort ist entschieden
  ([ADR-0060](../../adr/0060-citations-marker-scan-geteilte-prosa-antwort.md)).
- **Keine Änderung an den drei gescopten Roh-Lesungen** aus
  [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) —
  `versions`-Pins, `pins`-Span, `immutable`-Core beantworten andere Fragen.
- **Kein Umschreiben der 173 Fundstellen.** Wenn ein Befund echt ist, wird die
  Zeile repariert; wenn nicht, das Ventil richtig gesetzt — keine
  Massen-Umformatierung, um ein Gate ruhigzustellen.

## 4. Definition of Done

- [x] Die Zahl unter der Obergrenze ist **gemessen**, nicht geschätzt: **fünf**
      Befunde, nicht 58. Die Obergrenze galt einem **ganz** abgeschalteten
      Ventil und war damit Faktor zehn zu groß für diese Frage.
- [x] Je aufgedecktem Befund eine Beurteilung — und sie fiel anders aus als
      erwartet: **keiner** ist ein Doku-Defekt, alle fünf brauchen die Ausnahme
      wirklich. Was sich ändert, ist dass sie **gesetzt** ist statt aus der
      eigenen Prosa zu folgen.
- [x] Die Vertrags-Frage ist entschieden: `codepaths` und `ids` angeglichen,
      `diagrams` strukturell ausgenommen, `versions` als **benannte Grenze** —
      die Straten hängen zusammen (Lastenheft 0.68.0, Spezifikation §Achsen-
      Präzedenz **und** beide Algorithmus-Stellen,
      [ADR-0062](../../adr/0062-ventil-marker-versions-ist-eine-benannte-grenze.md)).
- [x] Die Gate-Frage aus
      [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)s
      Trigger ist **beantwortet**: ein Quelltext-Wächter hält die Klasse statt
      eines Reviews — mit gemessener Reichweite (zwei von sechs Umgehungen,
      datei-granular), die als seine Grenze bei ihm steht.
- [x] Ein konstruierter Verstoß je gedeckter Form, **Ursache gelesen**: die
      Erwähnung in Backticks je Modul, die Gegenprobe mit freiem Marker, die
      Divergenz auf derselben Prosa-Zeile, das Paritäts-Loch und die zwei neu
      gefangenen Umgehungs-Formen.
- [x] Doku-Currency: Handbuch, beide `README`-Fassungen, `CHANGELOG`.
- [x] `make gates` grün (Exit explizit), `make fullbuild` grün; unabhängiger
      Review.

## 5. Abnahme-Punkte / Risiken

- **Der Angleich macht ein Gate lauter, nicht leiser.** Bis zu 58 heute stille
  Zeilen melden danach — und ein Gate, das auf einen Schlag viel meldet, wird
  abgeschaltet statt gelesen. — **Ausgang: entfallen, gemessen.** Es sind
  **fünf**, nicht 58; die Obergrenze stammte aus einem Lauf mit **ganz**
  abgeschaltetem Ventil und beantwortete damit eine andere Frage.
- **Die Reparatur könnte punktuell statt strukturell ausfallen.** Vier
  Konsumenten einzeln umzustellen löst diesen Fall und nicht die Klasse; der
  fünfte Konsument entstünde morgen wieder roh. — **Ausgang: eingetreten,
  teilweise behoben.** Der Quelltext-Wächter ist die strukturelle Antwort und
  hält den fünften Konsumenten auf — aber nur in **zwei von sechs**
  Schreibweisen, und seine Erlaubnis ist **datei**-granular: eine zusätzliche
  Roh-Lesung in einer bereits gelisteten Datei bliebe unsichtbar. Beides steht
  als seine Grenze in
  [ADR-0062](../../adr/0062-ventil-marker-versions-ist-eine-benannte-grenze.md)
  Entscheidung 6, mit Wiedervorlage.
- **Ein Ventil, das enger wird, kann legitime Ausnahmen entwerten.** Wer
  `d-check:ignore` bisher in Backticks setzte, tat es womöglich absichtlich und
  in gutem Glauben. — **Ausgang: eingetreten, behandelt.** Genau das war der
  gesamte Ertrag: fünf Zeilen, die die Ausnahme **wirklich** brauchen und sie
  bisher aus ihrer eigenen Beschreibung bezogen. Sie tragen jetzt einen
  gesetzten Marker mit Begründung; entwertet wurde keine.
- **Ungeplant, vom Review aufgedeckt: die geteilte Antwort bringt zwei Ränder
  mit, und einer davon ist ein stiller Grün-Pfad.** — **Ausgang: eingetreten,
  benannt statt behoben.** Eine Code-Spanne desselben Absatzes verschluckt auch
  einen **gesetzten** Marker (Falsch-Rot); ein **unpaariger** Backtick weiter
  oben kippt die Parität, sodass eine Erwähnung doch als Direktive wirkt
  (stilles Grün). Beide gemessen, beide im Vertrag, beide mit Wiedervorlage —
  sie zu schließen bräuchte eine eigene Lexik neben der geteilten, also genau
  die zweite Antwort, gegen die
  [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  geschrieben ist.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Messung zeigt, dass der
Angleich eine Auftraggeber-Entscheidung über den Ventil-Vertrag verlangt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Module (GF), Spec-Straten (GF), Doku (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-27):
  [`BEO-020`](../observations.md) — die gemessene Menge muss die sein, über die
  geredet wird (die 58 sind eine **Obergrenze**, keine Antwort);
  [`BEO-011`](../observations.md) — die Regel gehört aus dem Bestand, nicht aus
  dem Anlass; [`BEO-017`](../observations.md) — ein rotes Gate muss vom
  geprüften Grund kommen.

Slice-ID: slice-162. Betroffene IDs:
[`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in),
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids),
[`DC-FA-VER-001`](../../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
[`DC-FA-DIAG-001`](../../../../spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in).
Module: `codepaths`, `ids`, `versions`, `diagrams`.
Gates: `make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Angleich einer vorhandenen Antwort an eine
vorhandene geteilte Antwort.

## 9. Closure-Notiz (nach `done/`)

**Der Angleich kostet fünf Zeilen Code und hat null Doku-Defekte gefunden. Sein
Wert liegt woanders — und drei meiner vier tragenden Aussagen darüber waren
falsch.**

**Die Sache selbst.** Seit
[ADR-0060](../../adr/0060-citations-marker-scan-geteilte-prosa-antwort.md)
beantwortete das Produkt *„ist diese Zeile eine Direktive"* zweifach: die
Zitat-Direktive gestrippt, der Ventil-Marker roh. Angeglichen sind jetzt die
zwei Konsumenten mit Prosa-Eingabe, `codepaths` und `ids`. Der Marker in
Backticks ist dort eine **Erwähnung**.

**Die Messung, die der Slice-Plan verlangte — und wie weit die Obergrenze
danebenlag.** §1 nannte **58** als Obergrenze. Das war die Zahl aus einem Lauf
mit **ganz** abgeschaltetem Ventil und beantwortete damit eine andere Frage. Der
Angleich legt **fünf** Befunde frei, Faktor zehn darunter.

**Und alle fünf sind Falsch-Rot im Wortsinn des Plans.** §2 unterschied *echt*
(das Ventil stand nie legitim dort) von *Falsch-Rot* (die Zeile braucht ein
Ventil, nur ein anderes). Die fünf brauchen die Ausnahme **wirklich** — ihre
Kennungen und Pfade sind Beispiele eines Akzeptanzkriteriums und ein
reproduzierter Fremd-Befund. **Gefundene Doku-Defekte: null.** Was sich ändert,
ist dass die Ausnahme **gesetzt** ist statt aus der eigenen Beschreibung zu
folgen. Die erste Fassung nannte sie „echt" und widersprach damit sowohl dem
Plan als auch der eigenen ADR zwei Absätze weiter. **Die Begründung des
Angleichs liegt damit im künftigen Fall, nicht im gemessenen Bestand** — das ist
eine schwächere Begründung, und sie steht jetzt als solche da.

**Drei weitere Aussagen waren am Bestand widerlegt, alle drei vom Review.**

1. **„`versions` liest keine Prosa" ist falsch.** Es liest **alle** Zeilen, also
   eine **Obermenge**, die die Prosa-Zeilen einschließt. Eigene Gegenprobe: auf
   **einer** Zeile mit dem Marker in Backticks meldet `codepaths` und schweigt
   `versions`; die Kontrollzeile ohne Marker meldet beides. Die Zweifach-Antwort
   ist dort **verschoben, nicht behoben**. Die Entscheidung hält — sie ist jetzt
   eine **benannte Grenze** statt einer „anderen Frage".
2. **Die Zahlen waren die falsche Grundgesamtheit.** *„160 Zeilen nur in
   Inline-Code"* war die Zahl der **baren Marker-Form**. Ich hatte die
   Beschriftung der einen Messung für die Zahl der anderen übernommen —
   [`BEO-020`](../observations.md), dritte Instanz in dieser Kette, und diesmal
   habe ich sie beim Formulieren des Review-Auftrags selbst bemerkt, aber erst
   danach. Nachgemessen: **553** Dateien, **249** Marker-Prosa-Zeilen, **66**
   frei, **183** nur in Inline-Code; unter den 66 wirksamen trägt genau **einer**
   die bare Form.
3. **Zwei Quellen über ihren Geltungsbereich hinaus zitiert.**
   [ADR-0019](../../adr/0019-versions-pin-fence-ausnahme.md) und
   [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) §2
   skopieren den **Pin**-Scan, nicht die Marker-Erkennung. Der Slice-Plan hatte
   es noch richtig (*„die `versions`-**Pins**"*); zwischen Plan und ADR wurde
   daraus „die Roh-Lesung von `versions`". [`BEO-012`](../observations.md).

**Was der Review zusätzlich fand und was gebaut ist.** Die geteilte Antwort
bringt **zwei Ränder** mit, die
[ADR-0060](../../adr/0060-citations-marker-scan-geteilte-prosa-antwort.md) für
die Zitat-Direktive ausdrücklich führt und die ich für den Ventil-Marker
übernommen, aber nicht benannt hatte: Verschluckung (Falsch-Rot) und Erzeugung
(**stilles Grün** — ein unpaariger Backtick kippt die Parität, und die Erwähnung
wirkt doch). Beide stehen jetzt im Vertrag, mit Richtung. Der Wächter fängt
**zwei von sechs** konstruierten Umgehungen; zwei billige sind geschlossen, die
Datei-Granularität bleibt seine Grenze. Und die Behauptung, der Compiler fange
die naheliegende Umgehung, ist **zurückgenommen**: sie gilt an genau einer
Stelle.

**Zur Form: eine ablösende ADR statt eines Edits.**
[ADR-0061](../../adr/0061-ventil-marker-geteilte-antwort-wo-die-eingabe-prosa-ist.md)
war beim Review-Ergebnis bereits gepusht und `Accepted`; §3.5 schützt ihren
Kern, und die widerlegten Aussagen standen genau dort. Eine
`## Geschichte`-Zeile hält fest, was sich geändert hat — sie kann keine falsche
**Herleitung** richtigstellen. Also
[ADR-0062](../../adr/0062-ventil-marker-versions-ist-eine-benannte-grenze.md),
die Herleitung, Zahlen und Grenzen ersetzt und das Verhalten unangetastet lässt.
**Die Lehre daran ist prozessual:** zum zweiten Mal in dieser Kette hat eine ADR
den Status `Accepted` getragen, bevor der Review sie gesehen hat. Beim ersten
Mal ließ sich der Commit noch neu aufbauen, hier nicht mehr.

**Sensors:** `make gates` (Exit 0, zehn Glieder, 557 Dateien, 0 Befunde,
Coverage 94,9 %), `make fullbuild` (Exit 0, 48 Anforderungen / 0 Waisen, bench
Median 714 ms), `make test` (Exit 0, sechs neue Grenz-Tests plus der Wächter).
**Proben, Ursache je gelesen:** Erwähnung in Backticks je Modul mit Gegenprobe
bei freiem Marker; die Divergenz `codepaths` gegen `versions` auf derselben
Zeile; das Paritäts-Loch; und der Wächter gegen zwei wiedereingebaute
Umgehungen — je gemeldet mit Datei und Zeile, Rückbau grün. Ein unabhängiger
Review ist gelaufen; sein Urteil war *„schließbar nach Nacharbeit"*, seine elf
Befunde sind eingearbeitet, und er hat seine Messmethode gegen
[ADR-0060](../../adr/0060-citations-marker-scan-geteilte-prosa-antwort.md)s
Zahlen **validiert**, bevor er meine widerlegt hat.
