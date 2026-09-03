# Slice slice-151: Die urteilsfreie Hälfte, so weit wie der Kanon sie benennt

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos.** Geschnitten vom Delta-Audit in
[slice-149](welle-85/slice-149-baseline-v5120-delta-audit.md) als Etappe C der
Baseline-Migration, bei deren Closure aber **herausgelöst**: der Gegenstand
berührt die Pin-Hebung nicht, und die Closure-Bedingung geht nicht über die
eigene DoD hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine
Welle braucht).

**Bezug:** [slice-143](../done/slice-143-structure-abschnitts-skopus.md) (die
heutige Deckung), [ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md),
[`BEO-015`](../observations.md),
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in).

**Berührte Spec-Stellen:** `spec/lastenheft.md` — **falls** die Messung ein
Produkt-Delta ergibt; Bump und Historie dann nach
[`MR-032`](../../../../harness/conventions.md#mr-032).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-26.

---

## 1. Ziel

Das Regelwerk benennt seit `v5.12.0` die urteilsfreie Hälfte der
Drei-Ausgänge-Regel ausdrücklich: urteilsfrei ist, **dass** zu jedem notierten
Risiko ein Ausgang dasteht **und welcher der drei** es ist — die drei sind eine
geschlossene Menge, kein Freitext. Es schließt mit: *„Welches Werkzeug die
urteilsfreie Hälfte prüft, ist Repo-Entscheidung; dass sie eine hat, ist es
nicht."*

Wir haben eine — sie prüft den **häufigsten Auslöser**, den stehengebliebenen
Vorlagen-Platzhalter. Zwei Fälle deckt sie nicht: ein Risiko **ganz ohne**
Ausgang, und ein Ausgang als **Freitext** statt einer der drei Formen. Der
zweite ist genau die Gestalt, in der
[`BEO-015`](../observations.md) auftrat — der erfundene vierte Ausgang.

**Die erste Frage ist eine Messung, keine Konstruktion:** Trägt `structure` die
Aussage *„jedes Risiko in §5 hat einen Ausgang"* überhaupt? Die Bedingungen
wirken auf den **Abschnitts-Text**, nicht je Listen-Eintrag — eine Korrelation
Risiko ↔ Ausgang ist damit womöglich nicht ausdrückbar.

## 2. Vorgehen

1. **Messen, was ausdrückbar ist**, bevor irgendetwas entschieden wird:
   `require-pattern`, `require-all`, `forbid-pattern` gegen die drei Formen —
   und ehrlich benennen, was davon *je Abschnitt* statt *je Risiko* wirkt.
2. **Am Bestand messen**, was eine kandidierende Regel im heutigen `done/`
   melden würde. Die Prüfmenge ist der Bestand von `done/` **zum Zeitpunkt der
   Messung** — beim Anlegen dieses Slice 142 Slice-Dateien; die Zahl ist zu
   messen, nicht aus einem älteren Slice zu übernehmen.
3. Reicht `structure` nicht, ist die Frage ein **Produkt-Delta** — und dann eine
   eigene Anforderung mit ADR, nicht ein Anhängsel.
4. Bewusstes Brechen je gedeckter Form, **Ursache gelesen** — nicht nur der
   Exit (Regelwerk `modul-13`, Schritt 6, seit `v5.12.0` in beiden Hälften).
5. `make gates`, `make fullbuild`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Urteils-Prüfung.** Ob ein eingetragener Ausgang inhaltlich **trägt**,
  bleibt Urteil — der Kanon sagt das ausdrücklich, und
  [ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md) hat die
  Grenze schon gezogen.
- **Kein Konformitätsdruck.** Der Kurs hat die Ausweitung ausdrücklich als
  unsere Entscheidung und **kein** Konformitätsthema bezeichnet. Ein „wir
  müssen" wäre hier falsch.
- **Keine zweite Mechanik.** Was entsteht, ersetzt oder erweitert die
  bestehende Regel; es tritt nicht daneben.

## 4. Definition of Done

- [x] Die Ausdrückbarkeits-Frage ist **gemessen** beantwortet, nicht
      angenommen — mit der Grenze *je Abschnitt vs. je Risiko*. *(Und der erste
      Anlauf hat sie falsch beantwortet: siehe Closure-Notiz.)*
- [x] Der Bestand ist gemessen; jede Fundstelle geräumt oder ausgewiesen.
      *(Ausgewiesen, nicht geräumt — die Begründung steht in der Notiz.)*
- [x] Bei Umsetzung: je gedeckter Form ein konstruierter Verstoß mit **gelesener
      Ursache**, Rückbau je grün. *(Beide Befund-Formen: `section-forbidden`
      und `section-missing`.)*
- [x] ~~Bei Nicht-Umsetzung: die Entscheidung steht **mit ihrer Messung** in der
      Closure-Notiz, und [`BEO-015`](../observations.md) trägt sie.~~
      **Entfällt — umgesetzt.**
- [x] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Eine Regel je Abschnitt kann eine Aussage je Risiko nicht treffen.** Fällt
  die Messung so aus, ist die ehrliche Antwort ein Produkt-Delta oder ein
  Verzicht — nicht eine Regel, die weniger prüft und mehr verspricht. —
  **Ausgang: eingetreten, aber nicht so, wie der Punkt es meinte.** Die Messung
  fiel zunächst genau so aus — und der Schluss daraus war **falsch**. Eine
  Bedingung je Abschnitt kann die Aussage je Risiko sehr wohl treffen, wenn sie
  ein `forbid-pattern` ist: das ist über **jedes Vorkommen** quantifiziert. Die
  ausgelieferte Regel prüft deshalb je Risiko, nicht je Abschnitt, und die
  Alternative „Produkt-Delta oder Verzicht" war nie nötig. Was der Punkt
  richtig vorhergesehen hat, ist die Gefahr: der erste Anlauf **hat** weniger
  geprüft und mehr versprochen, und ohne den Review wäre er so geschlossen
  worden.
- **Der Bestand ist gewachsen und uneinheitlich.** 139 Slices tragen ihre
  Ausgänge in Prosa; eine Form-Prüfung könnte breit rot laufen und damit einen
  Retrofit erzwingen, den niemand beschlossen hat. — **Ausgang: entfallen.** Die
  Ausnahme auf den abgeschlossenen Altbestand wendet es ab: ohne sie 122
  Befunde, mit ihr null, und kein künftiger Slice fällt heraus. Die Planzahl
  „139" ist überholt — gemessen sind es **150** Dateien, davon **51** mit
  Ausgangs-Markern.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Messung ein Produkt-Delta
verlangt, das eine eigene Anforderung braucht.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Module (GF), Harness-Regeltext (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-26):
  [`BEO-015`](../observations.md) ist der Anlass;
  [`BEO-011`](../observations.md) für jede Aussage darüber, was die Regel
  „vollständig" abdecke.

Slice-ID: slice-151. Betroffene IDs:
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in).
Module: `structure`, `planning`. Gates: `make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Messung vor Konstruktion an bestehender
Modul-Mechanik.

## 9. Closure-Notiz (nach `done/`)

**Die urteilsfreie Hälfte hat jetzt zwei Regeln statt einer.** Die vorhandene
fängt den **vergessenen** Ausgang (Vorlagen-Platzhalter), die neue den
**erfundenen** — die Gestalt aus [`BEO-015`](../observations.md). Der
Wortschatz ist geschlossen und umfasst alle drei Kanon-Ausgänge;
[`MR-049`](../../../../harness/conventions.md#mr-049) trägt sie.

**Der erste Anlauf war falsch, und zwar an der Stelle, die er selbst als
gemessen ausgab.** Er schloss: *„Die Je-Risiko-Korrelation ist **belegt** nicht
ausdrückbar."* Gemessen war nur, dass RE2 `(?!` abweist. Daraus folgt nichts
über die Frage — `forbid-pattern` ist über **jedes Vorkommen** quantifiziert
und trifft damit je Risiko, und das Komplement einer endlichen Wortmenge ist in
RE2 ohne Lookahead darstellbar. Der Review hat das nicht behauptet, sondern
**vorgeführt**. Die Lehre ist nicht „RE2 kann mehr als gedacht", sondern: **ein
gescheiterter Weg ist kein Beweis, dass es keinen gibt** — und das Wort
„belegt" stand vor der Messung, nicht dahinter.

**Der zweite Fehler war teurer, weil er still gewesen wäre.** Das erste Muster
erlaubte zwischen Marker und Ausgangswort höchstens ein Leerzeichen und zwei
Sterne. Es meldete damit **kanonische** Ausgänge als Verstoß, sobald der
80-Spalten-Umbruch dieses Repos vor das Wort fiel oder der Autor es fett setzte
— sieben der vierzehn gemeldeten Dateien waren in Wahrheit korrekt. Ein
Closure-Gate, das richtige Arbeit ablehnt, erzieht zum Ausweichen auf die
Formulierung, die durchkommt.

**Und `weiter offen` war nicht ungeprüft, sondern abgelehnt.** Der dritte
Kanon-Ausgang machte das Gate rot. Meine Begründung trug den Verzicht auf eine
*Zusatzprüfung*, nicht den Ausschluss aus der *akzeptierten Menge* — zwei
Sätze, die sich gleich lesen und Verschiedenes sagen.

**Eine eigene Warnung habe ich in derselben Botschaft ausgesprochen und
gebrochen.** Der verworfene Glob-Kandidat `slice-1[4-9][0-9]` hörte bei
`slice-200` still auf zu greifen — das stand als Messung 3 in der Botschaft.
Die stattdessen gewählte Ausnahme `slice-1[0-3]*` nimmt `slice-1000` bis
`slice-1399` heraus: derselbe Fehler mit weiterem Horizont, weil `matchGlob`
segmentweise matcht und `*` den Rest frisst. Die Ziffernzahl ist jetzt
festgenagelt.

**Der Bestand ist ausgewiesen, nicht geräumt — und die erste Auszählung war
falsch etikettiert.** „107 tragen gar keinen §5-Abschnitt" stimmt nicht:
**alle 150** tragen eine `## 5.`-Überschrift, 107 davon unter einem anderen
Titel (`## 5. Trigger` 43 · `## 5. Closure-Trigger` 31 ·
`## 5. Risiken / offene Punkte` 19 · weitere). Geräumt wird nichts: die
Ausnahme trägt den Altbestand, und ein Retrofit über 122 Dateien hat niemand
beschlossen — der Kurs hat die Ausweitung ausdrücklich als Repo-Entscheidung
und kein Konformitätsthema bezeichnet.

**Was die Sammelzahl verdeckte, ist der schwerste Bestands-Befund.** Drei
Dateien — `slice-106`, `slice-110`, `slice-111` — führen ein §5 mit Risiken und
**keinen einzigen** Ausgang. Das ist der Kanon-Kernsatz *„Ein Slice geht nicht
nach `done/`, während ein Risiko ohne Ausgang dasteht"*, dreimal verletzt, und
**keine** der beiden Regeln erreicht ihn: die Platzhalter-Regel braucht einen
Platzhalter, die neue braucht einen Marker. Ausgewiesen in
[`MR-049`](../../../../harness/conventions.md#mr-049) und hier; die Lücke ist
der Auflösungs-Trigger jenes Eintrags.

**Und eine Gegenprobe ist mir zunächst danebengegangen, genau wie das Register
es beschreibt.** Der erste Versuch hängte die Verstoß-Zeile ans **Dateiende** —
also hinter `## 6.`, außerhalb des geprüften Abschnitts — und meldete nichts.
[`BEO-017`](../observations.md) führt diese Gestalt wörtlich: *„eine Probe wird
ans Dateiende angehängt und landet im einzigen Abschnitt, den die Prüfung
ausnimmt"*. Sie zu kennen hat nicht geholfen; das Lesen der Ausgabe hat.

**Sensors:** `make gates` (Exit 0, zehn Glieder, 533 Dateien, 0 Befunde),
`make verify-closure-notes` (Exit 0, 479 Dateien, 0 Befunde). Gegenproben je
Befund-Form (`section-forbidden`, `section-missing`) mit gelesener Meldung; vier
kanonische Schreibweisen, die der erste Anlauf falsch rot meldete, laufen grün.
Ein unabhängiger Review ist gelaufen; seine zwei HIGH, fünf MEDIUM und drei LOW
sind in `2a2849f` eingearbeitet.
