# Slice slice-152: `citations` scharfschalten — und vorher die eigene Doku entschärfen

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos.** Geschnitten von
[slice-150](welle-85/slice-150-pin-gebundene-zitate.md) als Etappe C der
Baseline-Migration, bei deren Closure aber **herausgelöst**: der Blocker ist
älter als die Welle, und die Closure-Bedingung geht nicht über die eigene DoD
hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in);
[ADR-0045](../../adr/0045-zitat-verifikation-codepaths-erweiterung-und-citations-modul.md)
(das Modul und seine Fail-closed-Entscheidung);
[`MR-038`](../../../../harness/conventions.md#mr-038) (was ein Zitat beim Bump
tut); [`BEO-008`](../observations.md) (vierte Spiegel-Klasse).

**Berührte Spec-Stellen:** `spec/lastenheft.md` / `spec/spezifikation.md` — nur
falls die Lösung ein **Produkt**-Delta ist (Inline-Code-Bewusstsein des
Direktiven-Scans); Bump und Historie dann nach
[`MR-032`](../../../../harness/conventions.md#mr-032).

---

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-26.

## 1. Ziel

Die vierte Spiegel-Klasse hat eine mechanische Form, und sie liegt seit
`v0.50.0` im Produkt: das Modul `citations` prüft ein per
`d-check:cite`-Direktive ausgezeichnetes Zitat **gegen die von ihm zitierte
Quelle**, whitespace-normalisiert. Das ist genau die Prüfung, an der ein
Korpus-Test scheitert.

**Sie ist heute nicht aktivierbar, und das ist gemessen.** Ein Probelauf über
den Bestand bricht fail-closed an der ersten Fundstelle ab:

```text
d-check: error: CHANGELOG.md:592: malformte d-check:cite-Direktive — erwartet
<!-- d-check:cite <pfad>:<von>-<bis> --> (DC-FA-CITE-001.a Schritt 1, fail-closed)
```

Die Fundstelle ist die **Dokumentation der Direktive selbst** — der Kopfteil des
Musters steht dort in Inline-Code. Der Scan ist **fence**-bewusst, aber nicht
inline-code-bewusst; eine Fenced-Darstellung wäre immun, die gewählte ist es
nicht.

**Der Bestand ist gezählt, nicht geschätzt** (Marker außerhalb von Fences, über
`git ls-files`, Markdown-Dateien, ohne den vendorten Baum): **zehn** Dateien
tragen ihn — `CHANGELOG.md`, `README.md`, `README.de.md`, `spec/lastenheft.md`,
`spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`,
`docs/plan/adr/0045-…md` und drei Reporte aus `docs/reviews/2026-07-18-…`.
Der Lauf bricht dabei an **zwei** verschiedenen Stellen des Algorithmus:
**15** Marker sind malformt (Schritt 1), und **zwei** wohlgeformte Direktiven
tragen kein folgendes Zitat (Schritt 2). Vier der zehn Dateien sind
eingefroren (`docs/reviews/`) und dürfen nicht editiert werden — für sie taugt
nur ein Ventil oder ein Produkt-Delta.

**Neu ist das nicht.** Der Design-Review des Moduls hält denselben Blocker seit
**2026-07-18** als INFO-Befund fest, mit derselben Ursache und derselben
Datei-Klasse — und mit der damals gewählten Einordnung als *bewusste,
dokumentierte Fail-closed-Semantik*. Dieser Slice bringt keine neue
Entdeckung, sondern die Frage, ob diese Einordnung noch trägt, wenn das Modul
scharfgeschaltet werden soll.

**Stand nach dem ersten Anlauf (2026-08-27): zurückgeführt, weil die Wegwahl
ein Produkt-Delta verlangt.** Schritt 1 und 2 sind gefahren, das Ergebnis
schließt beide Wege bis auf einen — und der braucht eine eigene Anforderung.
Die Zahlen dieses Abschnitts oben sind dabei **überholt**: gemessen sind
**72** Vorkommen in **20** getrackten Dateien (nicht zehn), davon **70**
außerhalb eines Fenced-Blocks.

- **Weg A (Doku-Konvention) ist nicht teuer, sondern unmöglich.** Die 70
  Vorkommen außerhalb von Fences verteilen sich auf **neun** eingefrorene
  Review-Reporte, einen `done/`-Slice und zwei `Accepted`-ADRs
  ([ADR-0045](../../adr/0045-zitat-verifikation-codepaths-erweiterung-und-citations-modul.md),
  [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)). §3
  dieses Slice verbietet genau deren Bearbeitung, §3.5 die der ADRs.
- **Ein Ventil gibt es, aber kein passendes.** `citations` trägt **keine feine**
  Achse — kein `exempt-paths`, kein `ignore-refs`, keinen Zeilen-Marker. Grob
  wirken `scan.ignore` und `citations.scope`; beide nehmen die **ganze Datei**
  aus dem Modul und schalteten es damit gerade in `CHANGELOG.md`, den
  `README`-Fassungen, den Spec-Straten und dem Handbuch ab — also dort, wo echte
  Zitate stehen oder stehen werden. Ein Scharfschalten mit benannter Ausnahme
  ist damit nicht sinnvoll möglich.
- **Weg B ist ein Vertrags-Delta, kein Bugfix.** Die Spezifikation sagt
  ausdrücklich zu: *„Arbeitet auf den rohen Zeilen (fence-aware wie die übrigen
  Module)."* Das Verhalten ist spezifiziert, nicht abweichend — die Änderung
  braucht Lastenheft, Spezifikation und ADR.

Ausgetragen als [slice-158](../done/slice-158-citations-inline-code.md). Dieser
Slice wartet auf dessen Ergebnis; Schritt 3 bis 5 sind unberührt.

**Nachtrag bei der Closure — beide Mess-Generationen dieses Abschnitts sind
überholt.** Die zehn Dateien oben und die 72/20/70 hier zählten das
Direktiv-**Token**, also jede Erwähnung der Zeichenkette. Der Gegenstand ist
der **Marker**, also ein geöffneter HTML-Kommentar. Mit der Produkt-Lexik
nachgemessen ([slice-158](../done/slice-158-citations-inline-code.md)):
**25** Marker-Zeilen in 13 Dateien, davon 24 in Inline-Code, eine im Fence,
**keine** frei. Die Wegwahl bleibt davon unberührt — die Klasse der betroffenen
Dateien ist dieselbe —, aber die Zahlen sind es nicht.

## 2. Vorgehen

1. **Den Bestand zählen**, bevor entschieden wird: wie viele Stellen schreiben
   die Direktiv-Syntax in Prosa, und wie viele davon in Inline-Code gegen
   Fenced-Block?
2. **Die Wegwahl treffen und begründen** — Doku-Konvention (Syntax nur noch in
   Fenced-Blöcken) oder Produkt-Delta (der Scan überspringt Inline-Code wie die
   übrigen Module). Der zweite Weg ist die geteilte Lexik, der erste eine Regel
   für Menschen; beide haben einen Preis, und er gehört benannt.
3. Erst danach `citations` in den aktiven Modul-Satz — **mit** Bestandsmessung
   vor dem Scharfschalten.
4. Die Zitate aktiver `MR-*`-Einträge auszeichnen, so weit sie den **gepinnten**
   Stand zitieren. Historisch gestempelte Fassungen nach
   [`MR-038`](../../../../harness/conventions.md#mr-038) bleiben **ohne**
   Direktive — ihre Quelle existiert nicht mehr, und ein Wächter darauf wäre
   dauerhaft rot.
5. Bewusstes Brechen: ein verfälschtes Zitat ⇒ `citation-mismatch`, **Ursache
   gelesen**; Rückbau grün.
6. `make gates`, `make fullbuild`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Auszeichnung eingefrorener Dokumente.** `done/` und Review-Reporte
  zitieren den Stand ihrer Zeit.
- **Keine Direktive auf eine verschwundene Quelle.** Das wäre ein Wächter, der
  konstruktionsbedingt rot bleibt.
- **Keine Ausweitung auf Zitate außerhalb des Konventionsspeichers** in diesem
  Zug — der Bestand dort ist gemessen und klein.

## 4. Definition of Done

- [x] Der Bestand der Direktiv-Erwähnungen ist **gezählt**, getrennt nach
      Inline-Code und Fenced-Block.
- [x] Die Wegwahl ist begründet, mit dem benannten Preis des verworfenen Wegs.
- [x] `citations` läuft über den Bestand — Exit und Befundzahl genannt.
- [x] Die Zitate aktiver Einträge sind ausgezeichnet, soweit sie den gepinnten
      Stand zitieren; die Ausnahmen sind **benannt**, nicht übergangen.
- [x] Ein konstruierter Verstoß meldet `citation-mismatch` mit gelesener
      Ursache; Rückbau grün.
- [x] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Die Direktive trägt Zeilennummern in den vendorten Baum.** Jeder Bump
  verschiebt sie, und dann bricht die Prüfung — laut, aber sie bricht. Ob das
  Wartungs-Last oder gewollter Alarm ist, gehört entschieden statt erlitten. —
  **Ausgang: eingetreten, entschieden.** Gewollter Alarm: er ist genau die
  vierte Spiegel-Klasse aus [`BEO-008`](../observations.md), die sonst still
  alterte. Der Preis ist benannt — die Direktive meldet auch dann, wenn nur
  Zeilennummern gewandert sind —, und der Schritt liegt jetzt in der Prozedur,
  die den Bump ohnehin beschreibt: [`MR-051`](../../../../harness/conventions.md#mr-051),
  mit **drei** Fällen, weil der Grund-Code sie unterscheidet.
- **Fail-closed heißt: ein Fehler in EINER Direktive nimmt den ganzen Lauf
  mit.** Bei einem Gate im inneren Loop ist das eine andere Zumutung als bei
  einem Closure-Gate; die Bindepunkt-Frage gehört mitentschieden. —
  **Ausgang: eingetreten, entschieden.** Bindepunkt ist der **innere** Loop
  (`.d-check.yml`, damit `gates` und `pre-commit`). Tragbar, weil
  [ADR-0060](../../adr/0060-citations-marker-scan-geteilte-prosa-antwort.md)
  den Doku-Nebeneffekt entfernt hat. **Die erste Fassung dieser Begründung war
  zu weit** und behauptete, was übrig bleibe, sei ein Autoren-Fehler an einer
  gerade geschriebenen Zeile — drei Gegenwege sind belegt: der **Bump** meldet
  planmäßig, ein **Merge** trägt eine fremde Direktive am Hook vorbei, und
  `d-check:ignore` greift hier **nicht**. Alle drei stehen jetzt in
  [`AGENTS.md`](../../../../AGENTS.md) §4, in
  [`harness/README.md`](../../../../harness/README.md) §Sensors und am
  `modules`-Eintrag der [`.d-check.yml`](../../../../.d-check.yml).
- **Die Spiegel der Modulliste sind gate-blind.** Wer ein Modul aufnimmt, ändert
  still drei weitere Stellen mit. — **Ausgang: eingetreten.** Zwei Spiegel waren
  gebrochen, und keiner der beiden fiel auf: `FOCUS_DISABLE` ließ `citations` in
  vier fokussierten Gates mitlaufen, und der Netzlos-Guard blieb grün, weil er
  eine **Teilmenge** prüft. Als zweiter Beleg von [`BEO-010`](../observations.md)
  geführt, samt der Schärfung, dass dort nur **eine Richtung einer Hälfte**
  gebunden ist.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Wegwahl ein Produkt-Delta
verlangt, das eine eigene Anforderung braucht.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Module (GF), Konventionsspeicher (GF), Doku (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-26):
  [`BEO-008`](../observations.md) ist der Anlass;
  [`BEO-011`](../observations.md) für jede Aussage darüber, dass der Bestand
  „vollständig" ausgezeichnet sei.

Slice-ID: slice-152. Betroffene IDs:
[`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in).
Module: `citations`. Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Scharfschalten einer vorhandenen,
opt-in-Fähigkeit.

## 9. Closure-Notiz (nach `done/`)

**Das elfte Modul läuft. Der Ertrag ist nicht das Scharfschalten, sondern was
die Prüfung über den eigenen Bestand sagt: von 38 Zitaten in den aktiven
Konventions-Einträgen waren fünf nicht wörtlich.**

**Die Messung, mit der Produkt-Lexik** (Fence-Automat, absatzweise Spannen
gleicher Backtick-Länge, Zitate ab 16 Zeichen):

| Klasse | Anzahl |
|---|---|
| wörtliche Teilstrings des gepinnten Stands | **10** |
| Baseline-Zitate mit abweichender Auszeichnung oder Zeichensetzung | **5** |
| gegen den gepinnten Stand nicht prüfbar (Selbstzitate, Werkzeug-Meldungen, historische Pins, elidierte Zitate) | **23** |

**Ausgezeichnet sind neun** — acht der zehn wörtlichen plus einer nach
Korrektur. Die zwei übrigen wörtlichen stehen in der Delta-Tabelle von
[`MR-039`](../../../../harness/conventions.md#mr-039), deren Zeile den **alten
und den neuen** Wortlaut trägt: die Direktive paart mit dem **ersten** Zitat des
Absatzes, also mit dem historischen. Mit der heutigen Paarungsregel nicht
adressierbar; benannt statt weggelassen.

**Die fünf Abweichungen zerfallen in zwei Klassen, und die Grenze ist die
Richtung — nicht das Datum.** Was die Quelle **hat** und das Zitat weglässt
(eine Fettung, eine Kursivierung) oder was das Zitat **hinzufügt**, wo die
Quelle weiterläuft (ein Schluss-Punkt), ist ein **Transkriptions-Fehler**: drei
Einträge, korrigiert und ausgezeichnet. Was das Zitat an **Hervorhebung
hinzufügt**, ist ein **Autoren-Akt**: zwei Zitate in einem Eintrag, jetzt
deklariert statt stillschweigend, und ohne Direktive.

**Die erste Fassung begründete das mit dem Datum — und das war zweifach
falsch.** Sie sagte, für die vor dem Bump geschriebenen Einträge sei *„am Repo
nicht entscheidbar"*, ob ungenau transkribiert oder seither gedriftet, weil der
alte Pin nicht mehr vendored ist. Erstens führt der Bump-Eintrag selbst die
Delta-Tabelle, und eine der beiden Quelldateien steht **nicht** darunter.
Zweitens — und das ist der eigentliche Punkt — liegt der **alte Baum in der
git-Historie**: der Bump entfernt ihn aus dem Arbeitsbaum, nicht aus dem Repo.
Der Diff zeigt in beiden Quelldateien nur die geänderte Provenienz-URL plus
eine reine Ergänzung weiter unten. Die zitierten Zeilen sind unverändert. *„Am
Repo nicht entscheidbar"* war eine Aussage über den **Arbeitsbaum**, ausgegeben
als Aussage über das **Repo**.

**Zwei Entscheidungen sind getroffen, nicht erlitten** — §5 führt sie mit ihren
Ausgängen. Kurz: der Bump-Alarm ist gewollt, sein Preis benannt, und der
Neuanker liegt in [`MR-051`](../../../../harness/conventions.md#mr-051) mit drei
Fällen; der Bindepunkt ist der innere Loop, und die drei Wege, auf denen dort
ohne Autoren-Fehler Rot entsteht, stehen im Vertrag statt in der Erinnerung.

**Zwei Spiegel der Modulliste waren gebrochen, und keiner fiel auf.**
`FOCUS_DISABLE` ließ das neue Modul in vier fokussierten Gates mitlaufen —
fail-closed, über den ganzen Baum, `adr-check` im `pre-commit`-Hook —, obwohl
der Kommentar direkt darüber die Kopplung wörtlich benennt. Und der
Netzlos-Guard im Go-Test blieb **grün**, weil er eine **Teilmenge** prüft: er
fängt ein entferntes Modul, kein hinzugefügtes ohne Spiegel. Beide nachgezogen;
[`BEO-010`](../observations.md) steht auf zwei und trägt die Schärfung, dass
dort nur eine **Richtung** einer Hälfte gebunden ist.

**Eine Form-Sache, die ohne Renderer unsichtbar geblieben wäre.** Die
Direktiven standen zunächst allein auf einer Zeile. Eine Zeile, die mit dem
Kommentar-Öffner beginnt, ist nach CommonMark ein HTML-**Block** und
**unterbricht den Absatz** — mitten im Satz. Sie hängen jetzt am **Ende der
Vorzeile**; die Paarung bleibt, weil das Modul erst ab der Folgezeile sucht, und
das ist nachgeprüft.

**Was offen bleibt, benannt.** [`MR-039`](../../../../harness/conventions.md#mr-039)
§Geltungsbereich nennt auch [`AGENTS.md`](../../../../AGENTS.md),
[`harness/README.md`](../../../../harness/README.md) und die Skills; §3 dieses
Slice hat sie ausgeklammert. Dort stehen mindestens zwei wörtliche,
auszeichenbare Zitate — geschnitten als
[slice-163](../done/slice-163-zitate-ausserhalb-des-speichers.md), zusammen mit
dem vorbestehenden Befund, dass die historische Spalte der Delta-Tabelle ihr
eigenes Zitat nicht wörtlich wiedergibt.

**Sensors:** `make gates` (Exit 0, zehn Glieder, 552 Dateien, 0 Befunde),
`make test` (Exit 0, Netzlos-Guard auf elf Module), `make adr-check` (Exit 0,
wieder ohne `citations`), `make fullbuild` (Exit 0, 48 Anforderungen / 0
Waisen). **Bewusstes Brechen, Ursache gelesen:** ein Wort im ausgezeichneten
Kanon-Kernsatz verfälscht ⇒ `citation-mismatch` an der Direktiven-Zeile mit dem
richtigen Ziel und Grund-Code; Rückbau grün. Ein unabhängiger Review ist
gelaufen; sein Urteil war *„schließbar nach Nacharbeit"*, und seine vierzehn
Befunde sind eingearbeitet — die zwei HIGH sind die gebrochene Spiegel-Kopplung
und eine Zahl, die über zwei verschiedene Mengen sprach.
