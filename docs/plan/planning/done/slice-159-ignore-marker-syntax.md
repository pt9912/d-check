# Slice slice-159: Der `d-check:ignore`-Marker hat keine Syntax

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:**
[`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
(der Marker in seiner ersten Heimat);
[ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md) (das
geteilte Ventil); [ADR-0019](../../adr/0019-versions-pin-fence-ausnahme.md);
[`BEO-013`](../observations.md) (der Anlass);
[slice-146](../done/slice-146-ignore-marker-wirkung.md) (die Messung).

**Berührte Spec-Stellen:**
[`DC-FA-CODE-001.a`](../../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code),
[`DC-FA-ID-001.a`](../../../../spec/spezifikation.md#dc-fa-id-001a--kennungs-prüfung),
[`DC-FA-VER-001.a`](../../../../spec/spezifikation.md#dc-fa-ver-001a--versions-pin-konsistenz-versions),
[`DC-FA-DIAG-001.a`](../../../../spec/spezifikation.md#dc-fa-diag-001a--kennungs-konsistenz-in-diagramm-fences-diagrams).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-27.

---

## 1. Ziel

Der Unterdrückungs-Marker wird als **blanke Teilkette der Rohzeile** gematcht:
`strings.Contains(raw, "d-check:ignore")`. Damit unterdrückt **jede Erwähnung**
des Markers — in Prosa, in Inline-Code, in einer Tabellenzelle — genauso wie ein
gesetzter Marker. Die Dokumentation des Ventils schaltet die Prüfung ab, über
die sie schreibt.

**Gemessen** ([slice-146](../done/slice-146-ignore-marker-wirkung.md),
Stand vor jenem Slice): **233** Prosa-Zeilen tragen die Zeichenkette und
unterdrücken damit für `codepaths` und `ids`. Nur **62** sind gesetzte Marker;
**171** sind Erwähnungen — 146 blank, **25** zitieren die Kommentar-Form in
Inline-Code. Wird die Konstante ins Leere gelenkt, treten **58** Befunde auf
**48** Zeilen hervor (38 `id-unlinked`, 18 `codepath-missing`, 2 `repo-escape`);
die übrigen **185** Marker-Zeilen unterdrücken nichts.

**Die Mengen sind je Konsument verschieden**, und wer das übersieht, baut die
falsche Regel: `versions` liest **alle** Zeilen, in Fences wie außerhalb —
seine Menge ist **236**. `diagrams` liest **nur** Diagramm-Fence-Zeilen und die
Öffnungszeile; seine Menge ist zu den 233 **disjunkt**, und auf der
Öffnungszeile nimmt der Marker den **ganzen Block** aus.

**WIDERSPRUCH, gemeldet statt aufgelöst** ([`AGENTS.md`](../../../../AGENTS.md)
§1): Die unten naheliegende Kommentar-Form widerspricht einer **ausdrücklichen**
Festlegung der Spezifikation für `diagrams` — dort ist der Marker *„ein
**Token**, kein HTML-Kommentar"*, mit der Begründung, eine Kommentar-Lexik je
Fence-Sprache wäre ein Grammatik-Parser. Das Benutzerhandbuch sagt es als
Nutzer-Zusage weiter, und der Bestand trägt eine gesetzte Fundstelle in genau
dieser Form: eine `mermaid`-Öffnungszeile mit blankem Token im Infostring. Eine
repo-weit einheitliche Kommentar-Form würde sie **verlieren**. Der Widerspruch
ist damit vor jedem Code zu entscheiden — die Spec sticht den Plan, und wer sie
ablösen will, braucht eine ADR.

**Das ist dieselbe Klasse wie
[slice-158](../done/slice-158-citations-inline-code.md)**, nur mit umgekehrtem
Vorzeichen: dort bricht ein Modul an seiner eigenen Doku laut ab, hier schweigt
es lautlos. Die laute Variante fällt beim ersten Lauf auf; diese nicht.

**Nachtrag bei der Beanspruchung — dieser Abschnitt ist zur Hälfte überholt,
und zwar durch die eigene Diagnose.** §2 Schritt 2 unten benennt die tragende
Achse selbst: *„nicht die Kommentar-Form, sondern **Marker gegen zitierten
Marker**: die Erkennung muss Inline-Code auslassen."* Genau die hat
[slice-162](../done/slice-162-ignore-marker-geteilte-antwort.md) gebaut —
`codepaths` und `ids` lesen den Marker seither auf dem gestrippten Text.

**Was davon bleibt, ist die FORM** — und der Widerspruch unten, den dieser Plan
richtig gemeldet statt aufgelöst hat.

**Die Zahlen dieses Abschnitts sind Stand vor
[slice-146](../done/slice-146-ignore-marker-wirkung.md)** und nicht mehr die
Menge, über die hier entschieden wird. Nachgemessen mit der Produkt-Lexik zum
Zeitpunkt dieses Slice: **558** getrackte Dateien, **259** Marker-Prosa-Zeilen,
**66** davon wirksam, **193** nur in Inline-Code (die sind seit slice-162
entschärft). Die Verengung auf die Kommentar-Form kostet im Bestand deshalb
**null** Befunde —
und der Grund dafür ist eng: **65 der 66** wirksamen Marker tragen die Form
bereits, der 66. ist eine Erwähnung in Backticks, die nur durch das
Paritäts-Leck wirkt (benannte Grenze von
[ADR-0062](../../adr/0062-ventil-marker-versions-ist-eine-benannte-grenze.md))
und in einem für beide Konsumenten datei-weit ausgenommenen Verzeichnis liegt.
Einen **baren** wirksamen Marker gibt es im Bestand nicht.

**Auch §2 und §3 tragen denselben überholten Stand** — die Zahlen dort (87, 25,
48, drei, fünf, 58, 185) stammen aus derselben Messung wie die oben. Der
tatsächliche Preis dieses Slice ist **eine** Zeile, nicht drei.

## 2. Vorgehen

1. **Den Widerspruch oben zuerst entscheiden.** Solange die Spezifikation für
   `diagrams` den Token ausdrücklich festlegt, ist „eine Form für alle vier"
   keine Wahl, sondern eine Spec-Änderung. Möglich sind: die Spec per ADR
   ablösen, oder die Form **je Konsument** festlegen und die Uneinheitlichkeit
   begründen.
2. **Die Kommentar-Form allein reicht nicht — gemessen.** Von den 87 Zeilen in
   Kommentar-Form stehen **25 in Inline-Code**; eine Verengung auf die rohe
   Teilkette `<!-- d-check:ignore` würde sie **weiter** unterdrücken. Der
   diagnostizierte Defekt bliebe damit zu rund einem Drittel bestehen. Die
   tragende Achse ist deshalb nicht die Kommentar-Form, sondern **Marker gegen
   zitierten Marker**: die Erkennung muss Inline-Code auslassen, wie jedes
   prosa-lesende Modul es tut.
3. **Den Preis zählen — er ist bereits erhoben.** Von den 48 wirksamen Zeilen
   hängen **drei** an einer blanken Erwähnung statt an einem Marker
   (`spec/lastenheft.md`, Akzeptanzkriterien- und Historien-Prosa), und sie
   tragen **fünf** der 58 Befunde. Das ist der Preis einer Verengung — und alle
   drei sind selbst Instanzen des Defekts, den dieser Slice behebt.
4. Trägt die Änderung: Spezifikation (die Marker-Definition je Konsument), eine
   ADR mit `Schärft:`-Feld, dann Code und Tests. Ob das Lastenheft berührt ist,
   entscheidet sich daran, ob die Marker-Form dort zugesagt ist.
5. **Die vier Konsumenten gemeinsam — soweit Schritt 1 es zulässt.** Sie teilen
   heute eine Konstante, und eine Verengung in nur einem wäre die zweite Antwort
   auf dieselbe Frage
   ([ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)).
   Fällt Schritt 1 dagegen zugunsten der Spec aus, ist die Uneinheitlichkeit
   **keine** zweite Antwort, sondern eine **andere Frage** — dann gehört sie
   nach demselben ADR ausdrücklich gescopt und begründet.
6. Bewusstes Brechen je Konsument, **Ursache gelesen**; Rückbau grün.
7. `make gates`, `make fullbuild`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Regel „dieser Marker unterdrückt nichts".** Sie ist erst sinnvoll,
  wenn ein Marker ein Marker ist — vorher meldete sie 185 Zeilen Prosa. Das ist
  [slice-146](../done/slice-146-ignore-marker-wirkung.md)s Frage, und sie wartet
  auf dieses Ergebnis.
- **Keine Räumung eingefrorener Dokumente.** `done/` und `docs/reviews/`
  bleiben, wie sie sind; was dort nach der Verengung rot würde, braucht ein
  Ventil oder eine benannte Grenze.
- **Keine Änderung an den Ventil-Semantiken selbst** — welche Befunde ein
  gesetzter Marker unterdrückt, bleibt unangetastet.

## 4. Definition of Done

- [x] Der Widerspruch zur `diagrams`-Festlegung ist **entschieden**, nicht
      umgangen: die Spec bleibt, und die Form ist **je Eingabe** gescopt
      ([ADR-0063](../../adr/0063-marker-form-folgt-der-kommentar-lexik-der-eingabe.md)
      Entscheidung 1).
- [x] Die Marker-Form ist festgelegt und gegen den Bestand geprüft. **Es geht
      keine gesetzte Unterdrückung verloren** — die eine bekannte Kandidatin
      (`mermaid`-Öffnungszeile mit Token) ist belegt sicher: `proseLines` sieht
      sie gar nicht, die Form-Bedingung erreicht sie nie. Berührt wird **eine**
      Zeile, und sie war kein gesetzter Marker.
- [x] Die Erkennung lässt **Inline-Code** aus — geleistet von
      [slice-162](../done/slice-162-ignore-marker-geteilte-antwort.md); dieser
      Slice hat es beim Anspruch als überholt ausgewiesen statt es sich
      zuzuschreiben.
- [x] Spezifikation, ADR, Code und Tests hängen zusammen; die **verschiedenen**
      Antworten der Konsumenten sind als Tabelle gescopt und begründet.
- [x] Ein konstruierter Verstoß je Konsument mit **gelesener Ursache** — für
      `diagrams` erst in der Nacharbeit, dort fehlte die Assertion ganz.
- [x] `make gates` grün (Exit explizit), `make fullbuild` grün; unabhängiger
      Review.

## 5. Abnahme-Punkte / Risiken

- **Die Verengung deckt Befunde auf, die niemand bestellt hat.** Was heute eine
  Erwähnung deckt, wird danach rot — in lebenden Dokumenten ist das ein Gewinn,
  in eingefrorenen ein Problem ohne Adressaten. Die Menge gehört gezählt und
  entschieden, nicht erlitten. — **Ausgang: entfallen, gemessen.** Aufgedeckt
  wird **nichts**: 65 der 66 wirksamen Marker tragen die Form bereits, und die
  eine übrige Zeile liegt in `docs/reviews/`, das beide Konsumenten datei-weit
  ausnehmen. Der Preis, den §2 Schritt 3 auf drei Zeilen und fünf Befunde
  schätzte, ist **eine** Zeile und **null** Befunde.
- **Eine geteilte Konstante zu verengen, ändert vier Module auf einmal.** Die
  Gegenprobe muss je Konsument zeigen, dass ein **gesetzter** Marker weiter
  wirkt — sonst wird aus einem stillen Grün ein stilles Rot. —
  **Ausgang: eingetreten, behoben — und erst der Review hat es gesehen.** Die
  Verengung trifft nur zwei der vier; für `versions` gab es eine Probe, für
  `diagrams` **keine**, obwohl
  [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 4 je Konsument eine verlangt. Nachgeholt in der Fence-eigenen
  Konstellation; ein stilles Rot ist damit ausgeschlossen, nicht gehofft.
- **Ungeplant, vom Review aufgedeckt: die Form-Bedingung hat zwei Grenzen, und
  eine zeigt in die leise Richtung.** — **Ausgang: eingetreten, benannt.**
  Geprüft wird **zeilenweise** — ein Marker, dessen HTML-Kommentar auf einer
  früheren Zeile öffnet, gilt nicht (Falsch-Rot, laut). Und gefordert ist nur
  der **Öffner**, nicht der Abschluss — ein nie geschlossener Kommentar wirkt
  (stilles Grün, leise). Die Spezifikation sagte *„Klammer gefordert"* und lag
  damit über dem Code; sie ist auf die tatsächliche Bedingung gezogen, und beide
  Grenzen sind Proben.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls der gezählte Preis einen
Auftraggeber-Entscheid verlangt — dann bleibt die Lücke benannt, und
[`BEO-013`](../observations.md) trägt sie weiter.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Module (GF), Spec-Straten (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-27):
  [`BEO-013`](../observations.md) ist der Anlass;
  [`BEO-011`](../observations.md) — die Form gehört aus dem Bestand, nicht aus
  der einen Fundstelle, die aufgefallen ist.

Slice-ID: slice-159. Betroffene IDs:
[`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in).
Module: `codepaths`, `ids`, `versions`, `diagrams`. Gates: `make gates`,
`make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Verengung einer vorhandenen Erkennungs-Form
an vorhandenen Modulen.

## 9. Closure-Notiz (nach `done/`)

**Die erste Arbeit war, nicht loszulegen: die Hälfte dieses Plans war erledigt,
bevor er beansprucht wurde — durch seine eigene Diagnose.**

**§2 Schritt 2 benennt die tragende Achse selbst:** *„nicht die Kommentar-Form,
sondern **Marker gegen zitierten Marker**: die Erkennung muss Inline-Code
auslassen."* Genau die hat
[slice-162](../done/slice-162-ignore-marker-geteilte-antwort.md) gebaut. Der
Plan war also nicht falsch, sondern **überholt**, und zwar von einem Slice, den
seine eigene Analyse angestoßen hat. Ausgewiesen als Nachtrag in §1, nicht
stillschweigend weiterbenutzt — die Zahlen dort stammten von vor
[slice-146](../done/slice-146-ignore-marker-wirkung.md).

**Was blieb, war die Form — und ein Widerspruch, den der Plan richtig gemeldet
statt aufgelöst hat.** Die Spezifikation legt für `diagrams` den **Token**
ausdrücklich fest, weil im `mermaid`-Fence die Diagramm-Sprache den Kommentar
bildet und eine Lexik je Fence-Sprache ein Grammatik-Parser wäre. Eine
repo-weit einheitliche Form war damit nie verfügbar. Die Auflösung ist
dieselbe Regel auf verschiedene Eingaben: Markdown-Prosa bekommt den
HTML-Kommentar, Fence-Inneres und sprachgemischte Zeilen den Token.

**Der Ertrag im Bestand ist null Befunde — und der Grund dafür ist enger, als
die erste Fassung sagte.** Von **66** wirksamen Markern über 558 Dateien tragen
**65** die Form bereits. Der 66. liegt in `docs/reviews/`, das beide Konsumenten
datei-weit ausnehmen; der Repo-Lauf **konnte** nichts finden. Der einzige
Zähne-Beleg ist die konstruierte Gegenprobe: ein blanker Marker neben einem
toten Pfad wird mit der Verengung gemeldet und ohne sie nicht.

**Der schwerste Befund des Reviews galt meiner Messung, nicht dem Code.** Ich
hatte jene 66. Zeile als **bare Form** klassifiziert. Sie trägt den Marker in
**Backticks** und wirkt nur, weil ihr Absatz **ungerade Backtick-Parität** hat
— 221 Backticks, selbst nachgezählt. Sie ist also keine bare Form, sondern eine
**eingetretene Instanz** der Grenze *„Erzeugung ⇒ stilles Grün"*, die
[ADR-0062](../../adr/0062-ventil-marker-versions-ist-eine-benannte-grenze.md)
selbst benannt und dann für **nicht eingetreten** erklärt hatte. Einen baren
wirksamen Marker gibt es im Bestand nicht.

**Zwei Folgen daraus, beide festgehalten.** Erstens war
[ADR-0062](../../adr/0062-ventil-marker-versions-ist-eine-benannte-grenze.md)s
Trigger erfüllt, und dieser Slice hat die Instanz beseitigt, ohne dass jemand
entschieden hätte — als `## Geschichte`-Zeile dort vermerkt, samt der
Feststellung, dass die **Grenze fortbesteht** und nur ihre Fundstelle weg ist.
Zweitens waren die Zahlen wieder eine andere Grundgesamtheit
([`BEO-020`](../observations.md), vierte Instanz in dieser Kette): 249/183
gehören zum Stand jener ADR, nicht zu diesem. Nachgemessen über die
**Produkt-Scan-Menge**: 558 / 259 / 66 / 193.

**Ein Overclaim ist zurückgenommen.** *„Damit ist jede Erwähnungs-Form
entschärft"* — entschärft ist die **blanke**. Weiter wirksam sind eine escapte
Erwähnung, eine in einem HTML-Attribut und eine in einem eingerückten
Code-Block. Die haltbare Aussage über diese Änderung ist eine andere und
**stärker**: die Bedingung ist **monoton verengend** und kann keinen neuen
stillen Grün-Pfad erzeugen.

**Und `diagrams` hatte keine Assertion** — obwohl
[ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
Entscheidung 4 je Konsument eine verlangt und §4 dieses Slice sie fordert. Das
Akzeptanzkriterium war zudem in der Prosa-Konstellation formuliert, die es bei
`diagrams` nie gibt. Beides behoben; die Bestands-Fundstelle auf der
`mermaid`-Öffnungszeile ist damit **belegt** sicher statt behauptet.

**Sensors:** `make gates` (Exit 0, zehn Glieder, 559 Dateien, 0 Befunde,
Coverage 94,9 %), `make fullbuild` (Exit 0, 48 Anforderungen / 0 Waisen, bench
Median 710 ms), `make test` (Exit 0, **elf** Grenz-Proben — sechs zur Form, drei
zu `diagrams`, zwei zu den Rändern), `make adr-check` über die Range (0 Befunde;
beide ADR-Anhänge sind `## Geschichte`, kein Kern-Eingriff). **Proben, Ursache
je gelesen:** blanker Marker gegen Kommentar-Marker je Modul; die `>`-Kante; der
mehrzeilige Kommentar; der nie geschlossene; und drei `diagrams`-Fälle in ihrer
eigenen Konstellation. Ein unabhängiger Review ist gelaufen; sein Urteil war
*„schließbar nach Nacharbeit"*, seine vierzehn Befunde sind eingearbeitet, und
seine zwei HIGH sind eigens nachgemessen statt übernommen.
