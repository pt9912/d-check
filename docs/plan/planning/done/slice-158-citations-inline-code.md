# Slice slice-158: Der `citations`-Scan sieht Inline-Code nicht

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:**
[`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in)
(die Anforderung); [ADR-0045](../../adr/0045-zitat-verifikation-codepaths-erweiterung-und-citations-modul.md)
(das Modul); [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
(die geteilte Lexik und ihre gescopten Ausnahmen);
[slice-152](../done/slice-152-citations-scharfschalten.md) (der Anlass).

**Berührte Spec-Stellen:**
[`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in)
und [`DC-FA-CITE-001.a`](../../../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations) Schritt 1.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-27.

---

## 1. Ziel

[`DC-FA-CITE-001.a`](../../../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations) sagt über den Marker-Scan: *„Arbeitet auf den rohen Zeilen
(fence-aware wie die übrigen Module)."* Das ist eine **Zusage**, kein Versehen —
und sie hat eine Folge, die erst beim Scharfschalten sichtbar wird: **die
Dokumentation der Direktive ist selbst ein Fund**. Wer die Syntax in
Inline-Code schreibt — `<!-- d-check:cite <pfad>:<von>-<bis> -->` —, erzeugt
einen malformten Marker, und das Modul bricht fail-closed über den ganzen Lauf.

**Gemessen** ([slice-152](../done/slice-152-citations-scharfschalten.md)):
**72** Vorkommen des Markers in **20** getrackten Dateien, davon **70
außerhalb** eines Fenced-Blocks. Neun der zwanzig sind eingefrorene
Review-Reporte, dazu ein `done/`-Slice und zwei `Accepted`-ADRs — alle
unantastbar. Eine Doku-Konvention „Syntax nur noch in Fences" ist damit keine
teure Option, sondern unmöglich, und ein passendes Ventil hat das Modul nicht:
es trägt **keine feine** Achse (kein `exempt-paths`, kein `ignore-refs`, keinen
Zeilen-Marker), und die zwei groben — `scan.ignore` und `citations.scope` —
nähmen die **ganze Datei** aus dem Modul, gerade dort, wo echte Zitate stehen.

**Die Frage dieses Slice ist deshalb eine Vertrags-Frage**, keine
Implementierungs-Frage: soll der Marker-Scan Inline-Code überspringen wie jedes
andere prosa-lesende Modul — und wenn ja, ist das eine Schärfung der
vorhandenen Anforderung oder eine neue?

## 2. Vorgehen

1. **Die Vertrags-Frage zuerst.** [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
   trennt *„andere Antwort"* (Defekt) von *„andere Frage"* (per ADR gescopt) und
   führt `versions`, `pins`, `immutable` als gescopte Ausnahmen. `citations`
   steht dort **nicht** — die Rohzeilen-Lesart ist in der Spec zugesagt, aber
   nirgends als *andere Frage* begründet. Ob das eine Lücke oder eine
   ungeschriebene Absicht ist, gehört entschieden, bevor Code entsteht.
2. **Den Preis beider Antworten messen.** Überspringt der Scan Inline-Code,
   verschwindet die Selbst-Fundstelle — aber auch jede **echte** Direktive, die
   jemand in Inline-Code setzt. Ob es solche gibt, ist zu zählen, nicht
   anzunehmen.
3. Trägt die Änderung: [`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in) in
   [`spec/lastenheft.md`](../../../../spec/lastenheft.md) (Akzeptanzkriterium,
   Versions-Bump, Historie), [`DC-FA-CITE-001.a`](../../../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations) Schritt 1 in der Spezifikation,
   eine ADR mit `Schärft:`-Feld, dann Code und Tests.
4. **Die Fail-closed-Frage mitentscheiden**, weil sie an derselben Stelle
   hängt: ein malformter Marker nimmt heute den **ganzen Lauf** mit. Bei einem
   Modul im inneren Loop ist das eine andere Zumutung als bei einem
   Closure-Gate. Entweder bleibt es so — dann steht die Begründung im Vertrag —
   oder es wird ein Befund wie jeder andere.
5. Bewusstes Brechen je gedeckter Form, **Ursache gelesen**; Rückbau grün.
6. `make gates`, `make fullbuild`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Scharfschalten.** Das ist [slice-152](../done/slice-152-citations-scharfschalten.md),
  und es wartet auf dieses Ergebnis.
- **Keine Auszeichnung von Zitaten.** Auch das gehört zu slice-152.
- **Keine Änderung an den drei gescopten Ausnahmen** aus
  [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) —
  `versions`, `pins`, `immutable` beantworten andere Fragen und bleiben, wie
  sie sind.

## 4. Definition of Done

- [x] Die Vertrags-Frage ist entschieden: Schärfung, neue Anforderung oder
      bewusster Verzicht — mit Begründung gegen
      [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md).
- [x] Bei Änderung: Lastenheft, Spezifikation, ADR, Code und Tests hängen
      zusammen; kein Stratum verweist abwärts.
- [x] Die Zahl der **echten** Direktiven in Inline-Code ist gezählt, nicht
      angenommen.
- [x] Die Fail-closed-Entscheidung steht im Vertrag, nicht nur im Code.
- [x] Ein konstruierter Verstoß je gedeckter Form, **Ursache gelesen**.
- [x] **Doku-Currency:** `CHANGELOG`, Benutzerhandbuch und beide
      `README`-Fassungen tragen die neue Regel. *(Nachgetragen — die DoD führte
      diese Position zuerst nicht, obwohl §7 des Workflows sie verlangt: ein
      öffentlicher Vertrag ist berührt.)*
- [x] `make gates` grün (Exit explizit), `make fullbuild` grün; unabhängiger
      Review.

## 5. Abnahme-Punkte / Risiken

- **Eine Lexik-Änderung wirkt über das Modul hinaus.** Wer den Marker-Scan auf
  die geteilte Antwort umstellt, ändert eine Zusage, die andere Module teilen —
  und die Gegenprobe muss zeigen, dass sich für sie **nichts** ändert. —
  **Ausgang: entfallen, gemessen.** Geändert ist genau eine Datei
  ([`citations.go`](../../../../internal/hexagon/core/rules/citations.go)); die
  geteilte Vorverarbeitung ist unangetastet, und die drei gescopten
  Roh-Lesungen (`versions`, `pins`, `immutable`) sind nicht berührt. **Aber die
  Umkehrung ist eingetreten** und steht als eigener Punkt darunter: nicht das
  Modul wirkt über sich hinaus, sondern die **Frage** — sie hat jetzt zwei
  Antworten im Produkt.
- **Die Selbst-Fundstelle ist ein Sonderfall, der wie ein Allgemeinfall
  aussieht.** Dass ausgerechnet die Doku der Direktive stolpert, verführt zu
  einer Regel, die nur diesen Fall löst. Der Bestand entscheidet, nicht der
  Anlass ([`BEO-011`](../observations.md)). — **Ausgang: eingetreten, in beiden
  Ausprägungen von [`BEO-011`](../observations.md), und beide erst im Review
  gefunden.** (a) Eine **Exklusivitäts-Aussage**: *„`citations` führt keinen
  einzigen Konfigurations-Schlüssel"* — das Modul trägt sehr wohl
  `citations.scope`, und das Benutzerhandbuch sagte es die ganze Zeit richtig.
  (b) Eine **Zahl aus dem Anlass statt dem Bestand**: gezählt wurde die
  Zeichenketten-Erwähnung, geredet wurde über Marker. Beide sind ersetzt, und
  die Regel steht auf der Marker-Menge.
- **Die Entscheidung könnte an einem von zwei Konsumenten hängenbleiben.** Das
  Produkt kennt zwei Direktiven; wird nur eine umgestellt, beantwortet es
  dieselbe Frage zweifach. — **Ausgang: eingetreten, bewusst und
  ausgewiesen.** Die Ventil-Direktive bleibt roh; das ist der Defekt aus
  [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 1 und in
  [ADR-0060](../../adr/0060-citations-marker-scan-geteilte-prosa-antwort.md)
  Entscheidung 7 skopiert statt übergangen — mit gemessenen Rändern (173 Zeilen
  nur in Inline-Code; höchstens 58 Befunde) und einem eigenen Slice als Ausgang:
  [slice-162](../in-progress/slice-162-ignore-marker-geteilte-antwort.md).

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Vertrags-Frage einen
Auftraggeber-Entscheid verlangt — dann bleibt die Lücke benannt, und
[slice-152](../done/slice-152-citations-scharfschalten.md) wartet weiter.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Module (GF), Spec-Straten (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-27):
  [`BEO-011`](../observations.md) — die Regel gehört aus dem Bestand, nicht aus
  dem Anlass; [`BEO-017`](../observations.md) — ein rotes Gate muss vom
  geprüften Grund kommen.

Slice-ID: slice-158. Betroffene IDs:
[`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in).
Module: `citations`. Gates: `make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Schärfung eines vorhandenen Vertrags an
einem vorhandenen Modul.

## 9. Closure-Notiz (nach `done/`)

**Die Änderung ist fünf Zeilen Code; die Arbeit steckt darin, dass die erste
Fassung dreimal mehr behauptete, als sie gemessen hatte.**

**Die Sache selbst.** Die Spezifikation sagte über den Marker-Scan *„Arbeitet
auf den rohen Zeilen (fence-aware wie die übrigen Module)"* — die
Implementierung löste die zweite Hälfte des Satzes ein und die erste wörtlich:
fence-bewusst ja, inline-code-blind. Damit war die **Dokumentation der
Direktive** selbst ein malformter Marker, und weil das fail-closed ist, brach
der ganze Lauf ab. Gemessen: vorher endet er an `CHANGELOG.md:592`, nachher
meldet er 546 Dateien, 0 Befunde, Exit 0. Beide Seiten sind gefahren, die
Vorzustands-Seite gegen ein eigens aus dem Vorzustand gebautes Image.

**Die Vertrags-Frage, wie §2 sie verlangt.**
[ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
Entscheidung 2 trennt *andere Antwort* (Defekt) von *andere Frage* (per ADR
gescopt) und nennt für die geteilte Seite *„ist das eine Überschrift", „ist das
ein Anker", „ist das derselbe Absatz"*. *„Ist das eine Direktive"* fragt
ebenso nach der **Rolle** einer Zeile und bekommt deshalb die geteilte Antwort;
der **Zitattext** ist die andere Frage — dort sind Bytes der Gegenstand — und
bleibt roh, dieselbe Bauform wie beim gehashten Span von `pins`.

**Drei eigene Aussagen waren am Bestand widerlegt, und alle drei fand erst der
Review.**

1. **„`citations` führt keinen einzigen Konfigurations-Schlüssel."** Falsch —
   das Modul trägt `citations.scope`, und das Benutzerhandbuch sagte es die
   ganze Zeit genauer: keine **feine** Achse. Die Entscheidung hält trotzdem,
   nur steht jetzt der Grund statt der Behauptung: beide groben Achsen nehmen
   die **ganze Datei** aus dem Modul und schalteten es genau in den Dokumenten
   ab, die echte Zitate tragen.
2. **Die Messtabelle zählte die falsche Menge.** Gezählt war die
   Zeichenketten-**Erwähnung** (78 Vorkommen auf 74 Zeilen in 22 Dateien),
   geredet wurde über **Marker**. Mit der Produkt-Lexik nachgemessen sind es
   **25** Marker-Zeilen in 13 Dateien: 24 in Inline-Code (20 malformt, 4
   wohlgeformt), eine im Fence, **keine** frei. Der Satz *„von 74 auf null"* war
   in Wahrheit *„von 25 auf null"* — und dieselbe Wendung *„Marker-Vorkommen"*
   steht in [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
   für die andere Basis, sodass der Leser eine Vervierfachung sah, wo 18 → 25
   steht. Das ist [`BEO-020`](../observations.md), zwei Slices nach seiner
   Anlage.
3. **Eine Konsequenz war schlicht widerlegt.** Die ADR sagte, ein Pfad in
   Backticks falle fail-closed. Gemessen fiel er in einen Befund mit **leerem**
   Ziel, weil `(.+?)` auf die eingesetzten Leerzeichen matchte. Statt den Satz
   zu streichen ist das **Muster** gehärtet — der Pfad beginnt jetzt mit einem
   Nicht-Whitespace-Zeichen —, womit die Aussage stimmt und der krumme Befund
   verschwindet.

**Zwei Verhaltensweisen waren gar nicht benannt, beide gemessen.** Das Strippen
kann eine Direktive auch **erzeugen** (eine Code-Spanne zwischen `<!--` und dem
Marker verschwindet), und der Preis reicht weiter als *„steht in Backticks"*:
eine **freie**, unverklammerte Direktive wird ebenso verschluckt, wenn eine
Spanne **desselben Absatzes** sie umschließt. Beide stehen jetzt im Vertrag und
sind als Tests festgenagelt — nicht als Wunsch, sondern als bekanntes
Verhalten.

**Der schwerste Befund war keiner am Code.** Nach dieser Änderung beantwortet
das Produkt *„ist das eine Direktive"* **zweifach**: die Ventil-Direktive wird
weiter roh erkannt und wirkt aus Backticks heraus. Das ist genau der Defekt aus
[ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
Entscheidung 1 — und zugleich deren eigener Re-Evaluierungs-Trigger. Er wird
hier **nicht** behoben, weil die Richtung entgegengesetzt ist: dieses Modul wird
**stiller**, das Ventil würde **lauter**. Die Ränder sind gemessen statt
geschätzt — **173** Prosa-Zeilen tragen den Ventil-Marker ausschließlich in
Inline-Code, 63 frei, und ein Lauf mit **ganz** abgeschaltetem Ventil meldet
**58** Befunde. Der Angleich beträfe höchstens diese 58; wie viele wirklich, ist
die erste Messung von
[slice-162](../in-progress/slice-162-ignore-marker-geteilte-antwort.md).

**Was §3.8 hier zu benennen gab.** Das Modul liest die zitierte Quell-Spanne
**roh und typunabhängig** — kein Fence-Bewusstsein, kein Strippen, und die
Zieldatei muss weder Markdown sein noch in der Scan-Menge liegen; `scan.ignore`
skopiert die **prüfende** Datei, nicht das Ziel. Diese Eingabe hat das Modul nie
gescannt, und keine der neuen Zusagen gilt dort. Steht jetzt in beiden
Spec-Straten.

**Zur Form: der Feat-Commit ist vor der Veröffentlichung neu gebaut worden.**
[ADR-0060](../../adr/0060-citations-marker-scan-geteilte-prosa-antwort.md) wurde
als `Accepted` angelegt, und §3.5 schützt ihren Kern — die drei widerlegten
Aussagen standen genau dort. Der `pre-commit`-Hook wies den Kern-Edit korrekt ab
(`adr-check` im `STAGED=`-Modus vergleicht gegen `HEAD`); ihn zu übergehen wäre
ein umgangenes Gate gewesen, ein `Supersedes`-Nachfolger für eine ADR, die
diesen Zweig nie verlassen hat, eine leere Form. Weil der Commit **nicht
gepusht** war, ist er stattdessen zurückgenommen und neu aufgebaut: die ADR geht
**einmal** in die Historie, und zwar richtig. Der Beleg, dass drei ihrer
Aussagen erst falsch waren, ist diese Notiz — nicht eine Commit-Kette, die die
falsche Fassung konserviert.

**Sensors:** `make gates` (Exit 0, zehn Glieder, 548 Dateien, 0 Befunde),
`make fullbuild` (Exit 0, 48 Anforderungen / 0 Waisen), `make test` (Exit 0),
`make adr-check RANGE=origin/main..HEAD` (0 Befunde),
`make trace-check RANGE=origin/main..HEAD` (0 Befunde). **Proben, Ursache je
gelesen:** frei + malformt ⇒ fail-closed an der gesetzten Zeile; frei +
wohlgeformt + verfälscht ⇒ `citation-mismatch`; korrektes Zitat grün, während
dieselbe Direktive in Backticks mit einem Bereich weit hinter dem Datei-Ende
**stumm** bleibt — sie müsste sonst `citation-out-of-range` melden, die Stille
kommt also vom Strippen und nicht von Toleranz; Rückbau grün. Dazu sieben
Unit-Tests, darunter die absatzweite Spanne, die Zeile mit Erwähnung **und**
Direktive, der Backtick-Pfad und die erzeugte Direktive. Ein unabhängiger
Review ist gelaufen; sein Urteil war *„schließbar nach Nacharbeit"*, und seine
dreizehn Befunde sind eingearbeitet — die vier HIGH sind die drei widerlegten
Aussagen oben und die Zweifach-Antwort.
