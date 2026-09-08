# Slice slice-222: Baseline-Pin auf `v6.6.0` — mechanisch, ohne Urteil über den Delta

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`MR-011`](../../../../harness/conventions.md#mr-011)-Pin-Serie
(Vorgänger-Eintrag: [`MR-067`](../../../../harness/conventions.md#mr-067)),
[`MR-021`](../../../../harness/conventions.md#mr-021) (pin-gebundene
Verweise), [`MR-051`](../../../../harness/conventions.md#mr-051)
(`d-check:cite`-Spannen neu ankern),
[`MR-055`](../../../../harness/conventions.md#mr-055) (Symlink als Träger),
[`MR-069`](../../../../harness/conventions.md#mr-069) (`ignore-refs` als
deklarierte Gate-Senkung),
[`MR-070`](../../../../harness/conventions.md#mr-070) (Frozen-Klassen vor
mechanischer Ersetzung). Der neue Eintrag der Serie wird **`MR-071`** <!-- d-check:ignore (entsteht erst mit diesem Slice) -->.

**Berührte Spec-Stellen:** — (der Slice berührt keine Spec-Stelle).

**Verantwortlich:** pt9912 (Implementer-Rolle).

**Autor:** pt9912.

## 1. Ziel und Abgrenzung

**Ziel.** Den vendorten Baseline-Bestand von `v6.5.0` auf `v6.6.0` heben und
alle pin-gebundenen Verweise nachziehen — **mechanisch, ohne den Regel-Delta
zu beurteilen**. Was der Delta inhaltlich verlangt, entscheidet der
Folge-Slice; dieselbe Zerlegung wie beim Vorgänger-Paar
[slice-207](../done/slice-207-baseline-v650-bump.md) /
[slice-208](../done/slice-208-v650-regel-adoption.md).

**Der Anlass ist am Sensor bestätigt, nicht übernommen:**
`make baseline-freshness` meldet *„NEUER RELEASE verfügbar (Pin v6.5.0):
v6.6.0"* und zugleich, dass der **gepinnte** Tag upstream inhaltlich
unverändert ist (Bytes == vendored `SHA256SUMS`).

**Abgrenzung — vier Punkte, jeder mit Grund:**

1. **Kein Urteil über den Regel-Delta.** Was `v6.6.0` an Regeln ändert, wird
   **gemessen und gelistet**, aber nicht beantwortet — *ein Folge-Slice
   übernimmt es*, und zwar mit einer Antwort je Regel (übernommen · nicht
   anwendbar mit Begründung · abweichend als Adaption).
2. ~~**Keine Template-Adoption.**~~ **Aufgehoben — siehe Plan-Änderung
   unten.**
3. **slice-221 wird nicht angefasst.** Weder beansprucht noch nachgezogen —
   *Schicht-Abgrenzung*: Dieser Slice hebt einen Pin, er räumt keine Tabelle
   auf. **Das gilt weiter**, auch nach der Plan-Änderung: Die §4-Tabelle
   verschwindet **ganz**, statt gekürzt zu werden — damit erledigt sich
   slice-221s Gegenstand, aber das zu entscheiden ist Sache jenes Slice.
4. **Kein Aufräumen der `ignore-refs`-Einträge.** Die Hebung **fügt** einen
   hinzu (der `v6.5.0`-Baum verschwindet, eingefrorene Artefakte zitieren ihn
   weiter); dass die Liste damit wächst, ist als deklarierte Gate-Senkung
   bereits geführt ([`MR-069`](../../../../harness/conventions.md#mr-069)) —
   *Bestand bleibt bewusst stehen*.

**Plan-Änderung (2026-09-08) — die Template-Adoption kommt hinzu.**
Abgrenzung 2 schloss sie aus, mit der Begründung, der Bump sei mechanisch und
die Adoption ein Urteil. **Der Auftraggeber hat widersprochen, und das
Argument trägt:** *„Wenn wir AGENTS.md nicht anpassen würden, bräuchten wir
nicht auf v6.6.0 umstellen."* Der gemessene Delta ist klein (sechs
Regelwerk-Dateien), und seine **Schlagzeile ist genau diese Formänderung** —
ein Pin ohne sie wäre eine Versionsnummer. Schwerer wiegt: `AGENTS.md`
widerspräche der vendorten Vorlage, die **jeder Lauf** mitlädt.

**Was übernommen wird**, wörtlich aus der neuen `AGENTS.template.md`: *„Der
Gate-Index steht **einmal**, in `harness/README.md` §Sensors … Diese Datei
führt die Liste nicht."* Dazu die Template-`.d-check.yml`, die für das Modul
`targets` **beide** Schlüssel auf `harness/README.md` setzt.

**Machbarkeit vorab gemessen, nicht angenommen:** 54 Makefile-Regeln, 54 in
`harness/README.md` — keine Lücke in beide Richtungen. `exempt-targets` bleibt
deshalb leer; die Autorität wechselt die Datei, nicht die Strenge.

**Der Slice überschreitet damit die Ein-Sitzungs-Review-Grenze**
([`MR-066`](../../../../harness/conventions.md#mr-066)). **Grund für die
Nicht-Rückführung:** Bump und Adoption getrennt zu schneiden hieße, den
Widerspruch zwischen Vorlage und `AGENTS.md` für die Dauer eines Slice stehen
zu lassen — in der Datei, die jeder Lauf lädt. **Ersatz-Form der Prüfung,
vorab benannt:** ein **Bruch-Test in beide Richtungen** am Autoritäts-Wechsel
(ein erfundenes `make`-Target im neuen Index muss `gate-phantom` melden; eine
Makefile-Regel ohne Index-Eintrag muss `gate-undocumented` melden), dazu der
unabhängige Review über den Gesamt-Stand.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [x] **(1)** `.harness/baseline/v6.6.0/` ist materialisiert (`regelwerk/`,
      `templates/`, `SHA256SUMS`), der `v6.5.0`-Baum entfernt, §Baseline in
      [`harness/conventions.md`](../../../../harness/conventions.md) zeigt auf
      den neuen Tag, und **`MR-071`** <!-- d-check:ignore (entsteht erst mit diesem Slice) --> trägt die Hebung als nächster Eintrag der
      [`MR-011`](../../../../harness/conventions.md#mr-011)-Serie.
      `make baseline-verify` grün.
- [x] **(2)** **Alle vier Spiegel-Klassen sind nachgezogen, nicht nur die
      grep-bare.**
      [`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
      (5×) nennt sie: Pfad-Verweise (gate-gedeckt) · Release-/Tree-**URLs** mit
      dem Tag · **Prosa-/Ellipsen-Pins** · der **zitierende Verweis**, dessen
      Wortlaut am neuen Ziel nicht mehr stehen muss
      ([`MR-051`](../../../../harness/conventions.md#mr-051)). Je Klasse steht
      im Slice, **wie** sie gesucht wurde — drei davon deckt kein Gate.
- [x] **(3)** **Die Frozen-Klassen sind VOR der mechanischen Ersetzung
      aufgelistet** ([`MR-070`](../../../../harness/conventions.md#mr-070),
      Geltungsbereich trifft hier zu: eine mechanische Ersetzung über mehr als
      eine Datei) — über die **Eigenschaft**, nicht über Verzeichnisse. Was
      eingefroren ist, wird **nicht** retargetet, sondern über `ignore-refs`
      abgefangen.
- [x] `make gates` grün.
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

**Die Reihenfolge ist nicht beliebig, und der Grund steht in
[`MR-070`](../../../../harness/conventions.md#mr-070):** Die Frozen-Liste
entsteht **vor** der ersten Ersetzung, sonst ist sie eine Rechtfertigung
hinterher.

1. **Frozen-Klassen auflisten** — über die Eigenschaft *„würde ein
   korrigierter Wert verfälschen, was dieses Artefakt zu seinem Datum
   festgehalten hat?"*, nicht über Verzeichnisnamen.

**Ausgeführt am 2026-09-08, vor der ersten Ersetzung** — die Klassifikation
stand vorher, die **Zahlen** darunter sind im Review als falsch aufgefallen und
am Vorzustand (`git archive c0730c13`) neu gemessen.

**Die Form der Messung, ausgeschrieben, weil ohne sie jede Zahl beliebig ist:**
*Datei* = enthält die Zeichenkette `v6.5.0` mindestens einmal (`grep -rl`).
*Vorkommen* = **jedes einzelne Auftreten** (`grep -ro | wc -l`), nicht die
Zeile — 315 Vorkommen verteilen sich auf 302 Zeilen. *Pfad-Verweis* = ein
Vorkommen, dem `baseline/` unmittelbar vorausgeht; alles andere ist
*Versionsnennung*. **Die weite Form ist Absicht:** `.harness/baseline/v6.5.0`
als Muster übersieht die fünf **relativen** Schreibweisen `../baseline/v6.5.0` <!-- d-check:ignore (der Baum ist mit diesem Slice entfernt) -->,
und genau eine davon ist im Lauf dann auch übersehen worden (§7).

108 Dateien nennen `v6.5.0`, mit **315** Vorkommen:

| Klasse | Dateien | Pfad-Verweise | nur Versionsnennung | Behandlung |
|---|---|---|---|---|
| vendorter Baum | 28 | 0 | 29 | wird **ersetzt** |
| `done/`-Slice (Lauf-Beleg) | 12 | 27 | 28 | **eingefroren** |
| Review-Report (Lauf-Beleg) | 14 | 7 | 105 | **eingefroren** |
| `CHANGELOG.md` (Versions-Historie) | 1 | 0 | 1 | **eingefroren** |
| Register-Beleg (ab Merge unveränderlich) | 1 | 1 | 0 | **eingefroren** |
| **lebend** | 52 | **87** | **30** | retargeten |
| **Summe** | **108** | **122** | **193** | — |

**Die Einteilung ist die Eigenschaft, nicht der Pfad** — und der erste Versuch
über Verzeichnismuster fiel durch, er warf alles in eine Klasse. Das ist der
Grund für [`MR-070`](../../../../harness/conventions.md#mr-070) und hier
gleich wieder belegt; beim Nachmessen im Review ist derselbe Fehler ein zweites
Mal passiert, weil `grep -rl … .` die Pfade **ohne** `./`-Präfix ausgibt und die
`case`-Muster eines annahmen.

**Die 35 eingefrorenen Pfad-Verweise sagen nicht, wie viele Befunde entstehen.**
Der ursprüngliche Plan schloss von ihnen auf das Ventil — das ist ein Schluss
über die falsche Menge: **ein Verweis feuert, wenn ein Modul ihn auflöst**, und
24 der 27 in `done/` sind `d-check:cite`-Direktiven (deren Verzeichnisse
`citations.scope` ausnimmt), die sieben in `docs/reviews/` stehen in Inline-Code
und Fences. **Gemessen** — ein Lauf ohne den neuen Eintrag — sind es **4**
Befunde, alle vier `target-missing` auf dieselbe Datei. Der elfte
`ignore-refs`-Eintrag deckt genau sie; `.d-check.yml` führte bereits zehn, von
`v1.4.0` bis `v6.3.1` ([`MR-069`](../../../../harness/conventions.md#mr-069)).

Die **30** lebenden **Versionsnennungen ohne Pfad** sind die gate-blinde
Prosa-Klasse: Jede einzelne ist ein Urteil — Gegenwarts-Aussage (heben) oder
Vergangenheits-Aussage (stehen lassen). **Keine ADR nennt `v6.5.0`**, die
immutable Klasse ist also gar nicht betroffen.
2. **Delta messen, nicht beurteilen:** `diff -I '<!-- Quelle:'` gegen den
   alten Baum — die Herkunftszeile in Zeile 3 jeder Regelwerk-Datei trägt den
   Tag und meldete sonst *jede* Datei als geändert. Die Liste wandert in den
   Slice und ist die Eingabe des Folge-Slice.
3. **Re-vendoren:** `bash tools/harness/fetch-baseline-cache.sh v6.6.0`, alten
   Baum entfernen, `make baseline-verify`.
4. **Vier Spiegel-Klassen nachziehen**, je mit benannter Suchform.
5. **`d-check:cite`-Spannen neu ankern** — der Bump verschiebt Zeilennummern,
   und `citations` ist fail-closed im inneren Loop
   ([`MR-051`](../../../../harness/conventions.md#mr-051)).
6. `MR-071` <!-- d-check:ignore (entsteht erst mit diesem Slice) --> schreiben, §Baseline umstellen, `make gates`, Handoff.

## 4. Trigger

**Beanspruchung:** WIP-Limit frei (`in-progress/` ist leer, gemessen nach der
Rückführung von slice-221), `make baseline-freshness` meldet den neuen
Release, und der Auftraggeber hat den Vorgang am 2026-09-08 beauftragt.

**Rückführung nach `next/`** (`in-progress→next`): wenn der gemessene Delta
so groß ist, dass das **Nachziehen der Spiegel** selbst mehrere Sitzungen
braucht — dann trennt sich der Bump in Vendoring und Retargeting.

**Rückführung nach `open/`** (`in-progress→open`): wenn `v6.6.0` eine
**Struktur**-Änderung am vendorten Baum mitbringt (andere Verzeichnisnamen,
anderes Bundle-Layout), die die Pfad-Verweise nicht mechanisch abbildbar
macht. Dann ist vorher eine Entscheidung über die Verweis-Form fällig, wie sie
[`MR-023`](../../../../harness/conventions.md#mr-023) beim
Bundle-Layout-Wechsel gebraucht hat.

## 5. Closure-Trigger

DoD (1) bis (3) abgehakt, `make gates` und `make baseline-verify` grün mit
echter Ausgabe, unabhängiger Review durchgeführt und eingearbeitet,
Closure-Notiz geschrieben, Register fortgeschrieben, jedes Risiko aus §6 mit
einem der drei Ausgänge.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Drei der vier Spiegel-Klassen deckt kein Gate — und der Eintrag dazu steht
  bei 5×, ohne formgültigen Ausgang.** Ein grüner `make gates`-Lauf nach dem
  Bump sagt über Release-URLs, Prosa-Pins und zitierende Verweise **nichts**.
  Die einzige Gegenmaßnahme ist, je Klasse die Suchform aufzuschreiben und die
  Trefferliste anzusehen — DoD (2) verlangt genau das, und es ist eine
  Disziplin, kein Wächter. — **Ausgang:** eingetreten, und die Disziplin
  hat nur zur Hälfte getragen. Die Suchform je Klasse steht im Slice, aber
  **zwei der drei gate-blinden Klassen fielen erst nachträglich auf**: die
  Über-Hebung von [`MR-067`](../../../../harness/conventions.md#mr-067) beim Lesen des Diffs, die relative Schreibweise
  in `reviewer.md` erst, als der entfernte Baum sie zu `target-missing`
  machte — also durch einen Zufall, nicht durch die Suche. Was gewirkt hat,
  war die **Änderung der Suchform**: nach der Version statt nach dem Pfad,
  mit Gruppierung nach Präfix. Sie machte fünf Schreibweisen sichtbar.
  Beleg bei [`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md) (6×).
- **Der `citations`-Bruch ist die planmäßige Rot-Quelle, nicht ein Unfall.**
  Der Bump verschiebt Zeilenspannen; `citations` läuft fail-closed im inneren
  Loop und nimmt den `pre-commit`-Hook mit. Wer das nicht erwartet, hält es für
  einen Defekt und sucht an der falschen Stelle
  ([`MR-051`](../../../../harness/conventions.md#mr-051)). — **Ausgang:**
  eingetreten, aber **nicht** bei den Baseline-Zitaten. Die vier
  `d-check:cite`-Direktiven, die in geänderte Regelwerk-Dateien zeigen,
  blieben grün — die Spannen haben sich nicht verschoben, gemessen und nicht
  vorausgesetzt. Rot wurde `citations` **zweimal an einer internen Spanne**:
  `reviewer.md` zitiert `AGENTS.md` zeilengenau, und erst die
  Chronik-Löschung, dann die Tabellen-Entfernung haben sie verschoben. Die
  Rot-Quelle war also die richtige Klasse am unerwarteten Ort.
- **Die Über-Hebung ist so teuer wie die vergessene.** Eine
  Vergangenheits-Aussage („der Stand v6.5.0 führte …") mitzuheben macht sie
  falsch, und kein Gate meldet es — der Registereintrag nennt die Blindheit
  ausdrücklich **in beide Richtungen**. Die Frozen-Liste aus DoD (3) ist die
  Antwort darauf, und ob sie vollständig ist, bleibt Urteil. — **Ausgang:**
  eingetreten, und die Frozen-Liste hat sie **nicht** verhindert. Sie war
  korrekt und vor der Ersetzung erstellt — trotzdem hob das `sed` den
  `Geltungsbereich` von [`MR-067`](../../../../harness/conventions.md#mr-067)
  auf `v6.6.0`, obwohl der Eintrag die Hebung **auf v6.5.0 ist**. Der Grund
  ist eine Lücke in der Klassifikation, nicht in ihrer Ausführung: [`MR-067`](../../../../harness/conventions.md#mr-067)
  liegt in `harness/conventions/` und ist damit **lebend** — die
  Vergangenheits-Aussage steckt in seinem *Inhalt*, nicht in seiner Lage.
  **Eine Frozen-Liste über Datei-Eigenschaften kann das prinzipiell nicht
  fangen**; gefunden hat es erst das Lesen des Diffs. Gefunden **und** korrekt
  stehen gelassen wurden dagegen fünf weitere Vergangenheits-Aussagen — der
  `ignore-refs`-Tombstone, das Fremdzitat im CR und dreimal [`MR-067`](../../../../harness/conventions.md#mr-067) selbst.
  Beleg bei [`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md) (6×).

## 7. Closure-Notiz

**Geliefert.** Der Baseline-Pin steht auf `v6.6.0` (54 Dateien, `verify ok`),
[`MR-071`](../../../../harness/conventions.md#mr-071) trägt die Hebung als
dreizehnter Nachtrag der Serie, [`MR-067`](../../../../harness/conventions.md#mr-067) liegt in `harness/conventions/done/`. Dazu —
nach ausdrücklicher Plan-Änderung — die **Template-Adoption**: `AGENTS.md`
führt **keine Tabellenzeile** mehr (43 → 0) und ist von 59 826 auf
40 320 Bytes geschrumpft (−32,6 %), die Gate-Autorität liegt bei
`harness/README.md`. `make gates` grün. **Die Prozentzahl ist gegen den
Eingangs-Stand dieses Slice gemessen, nicht gegen den vor slice-221** — jener
hatte bereits 9 785 Bytes an Zell-Kürzung geliefert; sie hier mitzuzählen wäre
fremdes Ergebnis.

**Was funktioniert hat: die Machbarkeit vor der Änderung messen.** Bevor die
§4-Tabelle fiel, stand fest: 54 Makefile-Regeln, 54 in `harness/README.md`,
keine Lücke in beide Richtungen. Ohne diese Zahl wäre der Autoritäts-Wechsel
ein Sprung gewesen; mit ihr war er eine Buchung. Der Bruch-Test bestätigte
beide Richtungen — `gate-phantom` für ein erfundenes Target, `gate-undocumented`
für eine Regel ohne Eintrag, und der Befund nennt die neue Autoritäts-Doku
beim Namen.

**Was Friktion war: die mechanische Ersetzung ging zweimal daneben, und beide
Male sagte der Registereintrag es vorher.**
[`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
(6×) nennt vier Spiegel-Klassen und ihre Blindheit **in beide Richtungen** —
genau so trat es ein. Eine **Über-Hebung**: der Eintrag [`MR-067`](../../../../harness/conventions.md#mr-067) trug nach der Ersetzung einen Geltungsbereich,
der `v6.6.0` sagte — obwohl der Eintrag die v6.5.0-Hebung *ist*. Und eine
**vergessene**: `.harness/skills/reviewer.md` schreibt `../baseline/v6.5.0/` <!-- d-check:ignore (der Baum ist mit diesem Slice entfernt) -->
relativ, was das Muster `.harness/baseline/v6.5.0/` nicht traf. Gefunden hat
sie der Zufall — der entfernte Baum machte sie zu `target-missing`.

**Steering-Loop-Lerneintrag: die Suchform entscheidet, nicht die Sorgfalt.**
Erst `grep -rhoE '[^ ("]*baseline/v6\.5\.0'` **mit Gruppierung nach Präfix**
machte sichtbar, dass fünf Schreibweisen existieren — vier Pfad-Formen
unterschiedlicher Tiefe plus die relative. Nach der **Version** zu suchen
statt nach dem **Pfad** ist der Ableiter; er kostet nichts und hätte beide
Fehler vorab gefangen.

**Zweiter Lerneintrag, und er begrenzt eine Regel, die dieses Repo für stark
hielt.** Die Frozen-Liste nach
[`MR-070`](../../../../harness/conventions.md#mr-070) war korrekt und stand
**vor** der Ersetzung — sie hat die Über-Hebung trotzdem nicht verhindert.
[`MR-067`](../../../../harness/conventions.md#mr-067) liegt in `harness/conventions/` und ist damit **lebend**; seine
Vergangenheits-Aussage steckt im *Inhalt*, nicht in der *Lage*. **Eine
Frozen-Liste über Datei-Eigenschaften kann das prinzipiell nicht fangen** —
das ist keine Ausführungslücke, sondern die Grenze der Regel, und sie gehört
zu ihr gesagt. Beleg bei
[`mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
(4×), wo zusätzlich steht, dass der **erste** Einteilungsversuch nach
Verzeichnismustern alle 108 Dateien in eine Klasse warf — stumm und plausibel
aussehend.

**Dritter: der Auftraggeber-Einwand hat den Zuschnitt gerettet.**
[`form-vom-nachbarn-statt-von-der-vorlage`](../observations/BEO-ALL/form-vom-nachbarn-statt-von-der-vorlage/observation.md)
(2×) — slice-221 hatte die Zielform aus dem **Bestand** abgeleitet (kürzere
Zellen, Schwellen-Regel, grün und bruchgetestet), während die neue Vorlage die
Tabelle **ganz** streicht. Die Bestands-Antwort war nicht ungenauer, sondern
die **falsche Operation**. **Die Verschärfung, die daraus folgt:** Wo die
Vorlage erst nach einem Bump lesbar wird, ist *„an der Vorlage nachschlagen"*
keine jederzeit ausführbare Handlung — sie hat eine **Reihenfolge**, und die
gehört in den Zuschnitt.

**Vierter, und er wurde nicht vom Review gefunden, sondern beim Nachmessen einer
Zahl.** Die Kern-Lieferung dieses Slice — die Template-Adoption in `AGENTS.md`,
die Autoritäts-Umstellung in [`.d-check.yml`](../../../../.d-check.yml) und der
Zitat-Neuanker — lag in einem Commit, dessen Botschaft **ausschließlich
slice-223 nennt**. `git log --grep=slice-222` fand die Streichung nicht.
**Kein Gate sieht das:** `make trace-check` prüft, *dass* eine Kennung dasteht,
nicht *welche* — und `slice-223` ist eine gültige. Aufgefallen ist es nur, weil
die Größenangabe für diese Notiz gegen `git show` geprüft wurde statt aus der
Erinnerung geschrieben. Behoben, weil noch nichts gepusht war: der Misch-Commit
ist in zwei zerlegt, der Baum bit-identisch. Beleg bei
[`liefer-punkt-in-fremdem-commit`](../observations/BEO-ALL/liefer-punkt-in-fremdem-commit/observation.md)
(2×) — **derselbe Anlass wie beim Erstauftreten**: ein Kontext-Wechsel mitten in
der Arbeit.

**Review-Runde 1: 0 HIGH, 5 MEDIUM, 3 LOW, 3 INFO** — der Report liegt unter
[`docs/reviews/2026-09-08-slice-222-baseline-v660-review-r1.md`](../../../reviews/2026-09-08-slice-222-baseline-v660-review-r1.md).
Jeder Befund ist vor der Übernahme selbst nachgefahren; **alle acht
handlungsrelevanten sind behoben**, zwei INFO haben stattdessen eine Adresse
bekommen (unten). Zwei davon sind eigene Klassen, die dieser Slice **schon
kannte** und trotzdem produziert hat:

**Fünfter Lerneintrag — die Spiegelliste, die nicht geschrieben wurde.** Die
Autoritäts-Umschaltung ist eine Semantik-Änderung; für die verlangt
[`MR-025`](../../../../harness/conventions.md#mr-025), die Spiegel **vor** dem
Editieren per `grep` nach dem **alten** Wortlaut aufzulisten. Das ist nicht
geschehen — und der Slice hatte für einen anderen Gegenstand am selben Tag
genau so eine Liste erstellt. Stehen blieben acht Stellen: die beiden
**Beschreibungen des Gates selbst** (Blockkommentar der
[`.d-check.yml`](../../../../.d-check.yml), Vertrags-Sektion von
[`harness/sensors/gate-consistency.md`](../../../../harness/sensors/gate-consistency.md))
und sechs Zeiger auf einen entleerten Abschnitt. **Kein Gate sieht das**, und
zwar aus einem Grund, der die Klasse definiert: Der Link *löst auf*; leer ist
nur seine Aussage. Beleg bei
[`semantic-change-body-only-edges-stale`](../observations/BEO-ALL/semantic-change-body-only-edges-stale/observation.md)
(17×).

**Sechster — zweimal ein Proxy statt des Gegenstands, im selben Slice.** Die
Frozen-Tabelle mischte Zeilen und Vorkommen (299 behauptet, 315 gemessen, die
Tabelle selbst summierte 300), und aus „33 eingefrorene Pfad-Verweise" schloss
der Plan auf den Umfang des Ventils — **gemessen sind es 4 Befunde**, weil ein
Verweis feuert, wenn ein Modul ihn *auflöst*, und die übrigen
`d-check:cite`-Direktiven in ausgenommenen Verzeichnissen sind. Beide Zahlen
sind korrigiert, in §3 **und** im lebenden
[`MR-071`](../../../../harness/conventions.md#mr-071), jeweils mit der
**Form der Messung daneben**. Beleg bei
[`zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
(6×) — dort steht auch, dass die Klassifikation beim Nachmessen ein **zweites
Mal** still scheiterte, weil `grep -rl … .` die Pfade ohne `./`-Präfix ausgibt.

**Siebter — die Ausnahme war richtig, ihr Grund nicht.** Der elfte
`ignore-refs`-Eintrag deckt genau die vier Befunde; begründet war er mit der
fehlenden §4-Tabelle der `AGENTS`-Vorlage, während beide ausgenommenen Belege
ausschließlich die `README`-Vorlage zitieren (39 → 40 Tabellenzeilen). Grün,
eng, richtig — und der Satz daneben ist die einzige Fassung, aus der der
nächste Bump seine Reichweite ableitet. Beleg bei
[`begruendung-traegt-entscheidung-nicht`](../observations/BEO-ALL/begruendung-traegt-entscheidung-nicht/observation.md)
(3×).

**Ein Befund traf eine Regel, die dieses Repo schon verkörpert hat.** §4 sagte
nach der Umstellung *„Kein Target nennen … auch nicht in Prosa"* und im selben
Atemzug *„maschinell gehalten"* — die Prosa-Hälfte trägt kein Mechanismus, und
seit der Umschaltung ist `AGENTS.md` weder Scan-Ziel noch Autorität. Das ist
[§3.8](../../../../AGENTS.md#38-ein-modul-verspricht-nur-über-das-was-es-scannt)
auf diesen Sensor angewandt; die Grenze steht jetzt bei der Zusage. Kein neuer
Registereintrag — die Regel gibt es, sie wurde verletzt.

**Was offen bleibt — drei Punkte, jeder mit Adresse oder benannter
Adresslosigkeit.** **(1)** `make baseline-verify` prüft die **Auflösung** der
Symlinks, nicht ihren **adoptierten Stand** — das Schwester-Repo `a-check` hat
dafür einen eigenen Sensor (`symlink-check`), der beides hält. Hier wäre ein
vergessener Symlink nur zufällig aufgefallen: weil der alte Baum verschwand,
wäre er tot gewesen. **Im Fenster, in dem beide Bäume existieren, meldete er
grün.** Der Review hat denselben Punkt an seiner dauerhaften Stelle gefunden:
[`harness/sensors/baseline-verify.md`](../../../../harness/sensors/baseline-verify.md)
§Grenze nennt den Fall nicht, der beim Bump trägt. **Bewusst nicht hier
mitgenommen** — der Slice fasst keine Sensor-Zusage an, und die Ergänzung ohne
den Sensor wäre die Hälfte, die schon jetzt in der Closure-Notiz steht. Eigener
Vorgang, weiterhin ohne Kennung; das ist die benannte Adresslosigkeit, keine
Erledigung. **(2)** Der Regel-Delta ist **gemessen und gelistet, nicht
beantwortet** — sechs Regelwerk-Dateien, fünf Templates. Die Antwort je Regel
ist Sache eines Adoptions-Slice; slice-221 wartet darauf und
[slice-223](../open/slice-223-commit-zerlegung-ausnahmen-aufloesen.md) nimmt
den ersten Teil davon vorweg. **(3)** Dazu gehört ein Spiegel, den erst die
Umschaltung sichtbar macht: die **Produkt**-Oberfläche zeigt in
[`docs/user/benutzerhandbuch.md`](../../../user/benutzerhandbuch.md) und im
emittierten Config-Template weiter `authority: AGENTS.md`, während die
adoptierte Vorlage `harness/README.md` zeigt. Das fällt unter Abgrenzung 1 —
**kein Verstoß, aber auch kein Zufall**: Es ist derselbe Adoptions-Slice, der
ihn zu beantworten hat.

**Die drei Paarungen, gemessen.** **(a) Anker** — vakant: Der Slice verkörpert
keine neue Steering-Loop-Regel; seine sieben Lerneinträge liegen bei bestehenden
Registereinträgen. **(b) Folge-Slice** —
[slice-223](../open/slice-223-commit-zerlegung-ausnahmen-aufloesen.md) ist
genannt und liegt als Datei in `open/`; die `symlink-check`-Lücke ist bewusst
**ohne** Kennung gelassen. **(c) Register** — alle zitierten Pfade lösen auf,
die sieben Belege liegen als `evidence/slice-222.md` in ihren Verzeichnissen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende.

**Die drei Vorprüfungen sind bei der Beanspruchung am 2026-09-08 bestätigt**
— Nachtlauf und Register unverändert gegenüber dem Anlege-Stand desselben
Tages; die Sub-Area-Wahl trägt der Block unten. Ursprünglich notiert war:
([`AGENTS.md`](../../../../AGENTS.md) §5). **Zwei sind schon gelaufen**, weil
sie den Zuschnitt tragen, und werden bei der Beanspruchung gegen den dann
gültigen Stand wiederholt:

- **Nachtlauf-Stand** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
  am 2026-09-08 beide grün. **Für diesen Slice ausnahmsweise nicht bezugslos:**
  `make baseline-freshness` ist Teil desselben Nachtlaufs, und sein Exit 3 ist
  der Auslöser dieses Vorgangs.
- **Register** (40 Verzeichnisse): **drei** Einträge sind einschlägig.
  [`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
  (5×) ist **der tragende** und formt DoD (2) — er nennt die vier Klassen und
  benennt, dass drei davon gate-blind sind, **in beide Richtungen**. Sein Stand
  ist *kein formgültiger Ausgang*, geführt bei
  [`registerzeile-ohne-ausgang-nach-schwelle`](../observations/BEO-ALL/registerzeile-ohne-ausgang-nach-schwelle/observation.md);
  dieser Slice ändert daran nichts und erfindet auch keinen.
  [`mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
  (3×, verkörpert als [`MR-070`](../../../../harness/conventions.md#mr-070))
  ist hier **anwendbar** — anders als in slice-221: Der Geltungsbereich nennt
  *„jede mechanische Ersetzung über mehr als eine Datei"*, und genau das ist
  das Retargeting. DoD (3) löst die Pflicht ein.
  [`form-vom-nachbarn-statt-von-der-vorlage`](../observations/BEO-ALL/form-vom-nachbarn-statt-von-der-vorlage/observation.md)
  (1×) ist der Grund, warum slice-221 wartet — und die Mahnung, in diesem Slice
  die **Form** der neuen Vorlage nicht nebenbei zu übernehmen (Abgrenzung 2).

**Modus-Begründungsblock.** Alle berührten Sub-Areas GF — ein Block genügt.

### Sub-Area: `*`

- **Modus:** GF (Repo-Default).
- **Konventions-Dichte:** **sehr hoch.** Die Pin-Serie ist zwölfmal gelaufen
  ([`MR-011`](../../../../harness/conventions.md#mr-011) bis
  [`MR-067`](../../../../harness/conventions.md#mr-067)), das Vorgehen ist in
  fünf `MR`-Einträgen verankert, und der Vorgänger-Slice liegt als Vorlage vor.
- **Phase-Reife:** Phase 5.
- **Evidenz-/Diskrepanz-Risiko:** **niedrig für den vendorten Baum**
  (`baseline-verify` prüft ihn hart), **mittel für die Spiegel** — drei der
  vier Klassen sind gate-blind, und der Registereintrag dazu steht bei 5×.
- **Reconciliation-Aufwand:** keiner (GF).
