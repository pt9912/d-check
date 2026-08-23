# Slice slice-136: §3.4 gegen den Kanon halten — Doppelung, Verschärfung oder keins von beidem?

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-84-durchsetzung](welle-84-durchsetzung.md).

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §3.4; Baseline-Regelwerk
[`modul-03-spec.md` §Ziel-Form: Architektur-Sicht](../../../../.harness/baseline/v5.11.0/regelwerk/modul-03-spec.md)
und
[`grundlagen-referenz-richtung.md`](../../../../.harness/baseline/v5.11.0/regelwerk/grundlagen-referenz-richtung.md);
die vendorte `architecture.template.md`; der Zensus in
[slice-132](../done/slice-132-hard-rule-zensus.md).

**Berührte Spec-Stellen:** — (Harness-Regeltext; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

Der Zensus hat §3.4 als *teilgedeckt* eingestuft und die ungedeckte Hälfte —
die **Sprachfreiheit der Architektur-Sicht** — mit *Auflösungs-Trigger:
permanent* versehen, also für grundsätzlich unmechanisierbar erklärt. Auf
Auftraggeber-Nachfrage stellte sich heraus: **das ist zu grob.** Zwei Fragen
hängen daran, und keine war beantwortet.

**Erstens: Ist §3.4 eine Doppelung?** Seine zweite Aussage — *kein Spec-Stratum
referenziert ADRs, Wellen, Slices, Commit-Hashes oder Closure-Daten* — steht
im Kanon **wörtlich**, in `modul-03` und in `grundlagen-referenz-richtung.md`.

**Zweitens: Ist §3.4 eine Verschärfung?** Wir schreiben *keine
Sprach-/Modul-Pfade*; `modul-03` schreibt, die Sicht *„referenziert
Modul-Pfade"*. Meint das dasselbe? Wenn ja, verschärfen wir — und
[`MR-031`](../../../../harness/conventions.md#mr-031) verlangt dafür einen
Eintrag im Konventionsspeicher, den es nicht gibt.

## 2. Vorgehen

1. **Die Doppelungs-Frage entscheiden** — an `modul-09`, nicht am Gefühl: eine
   Hard Rule soll in **zwei** Quadranten liegen, und der Feedforward-Quadrant
   *ist* der `AGENTS.md`-Eintrag. Prüfen, ob daraus folgt, dass eine gerankte
   Wiederholung des Kanons hier **erwünscht** statt verboten ist.
2. **Die Verschärfungs-Frage messen, nicht auslegen.** Der Bullet-Text ist
   mehrdeutig; die vendorte `architecture.template.md` ist es nicht — sie
   **praktiziert** eine Lesart. Zählen, welche Pfade sie führt und welche
   Kennungen ihre Komponenten-Tabelle trägt.
3. **Je nach Ergebnis:** `MR-`Eintrag (bei Verschärfung), Kürzung (bei
   unnötiger Doppelung) — oder eine berichtigte Ausweisung, wenn beides
   entfällt.
4. **Die Ausweisung berichtigen.** *Permanent* gilt allenfalls für die
   Technologie-Hälfte; die Pfad-Hälfte ist ein **detektierbarer Zustand**
   (heute gemessen: `spec/architecture.md` null, `spec/spezifikation.md` fünf).
   Was davon baubar ist und womit, gehört benannt statt pauschal verneint.
5. **CR-Kandidat festhalten**, falls die Mehrdeutigkeit im Kanon bleibt: sie
   trifft jedes Adopter-Repo, nicht nur uns.
6. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Sensor.** Selbst wenn die Pfad-Hälfte baubar ist — sie zu bauen ist
  ein eigener Slice mit eigener Ortswahl.
- **Kein Konsumenten-CR.** Ein CR-**Kandidat** wird benannt; ob und wann er
  geschrieben wird, ist eine Auftraggeber-Entscheidung.
- **Keine Änderung an den Spec-Straten selbst.** Der Slice prüft eine Regel,
  nicht den Bestand, den sie regelt.

## 4. Definition of Done

- [x] Beide Fragen mit **Zitat** beantwortet, nicht mit Auslegung — **und die
      zweite Antwort ist die umgekehrte der ersten Fassung** (§9).
- [x] Die Folge ist gezogen, und es war die unbequeme:
      [`MR-033`](../../../../harness/conventions.md#mr-033) statt einer
      Konformitäts-Erklärung.
- [x] §3.4 trägt keine unbelegte Unmöglichkeits-Behauptung mehr — und darüber
      hinaus keinen unbelegten Deckungs-Anspruch: die Abwärts-Sperre steht jetzt
      auf **zwei von fünf** Kategorien.
- [x] `make gates` Exit 0 (zehn Glieder, 479 Dateien, 0 Befunde); unabhängiger
      Review ([Report](../../../reviews/2026-08-23-slice-136-agents-34-klaerung-review.md)),
      blockierend mit einem HIGH, alle fünf Befunde eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **Die Versuchung ist, die bequeme Lesart zu wählen.** — **Ausgang:**
  *eingetreten. Ich bin ihr erlegen.* Der erste Anlauf schloss auf „keine
  Verschärfung" und stützte sich dabei auf die **Praxis** einer Vorlage, während
  die direkteste Quelle ungelesen blieb. Der Risiko-Satz stand geschrieben,
  bevor er zutraf — und hat mich nicht davor bewahrt.
- **Eine Vorlage ist kein Regeltext.** — **Ausgang:** *eingetreten, und zwar
  genau an der Stelle, an der es zählte.* Die Grenze war benannt; sie zu nennen
  hat den Fehlschluss nicht verhindert. Was ihn verhindert hätte, war eine
  einzige Suche in der richtigen Datei.

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-135](../done/slice-135-uses-pin-sensor.md)
in `done/`, WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Messung eine echte
Verschärfung zeigt und die Frage, ob wir sie behalten wollen, eine
Auftraggeber-Entscheidung ist.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Regeltext (GF), Konventionsspeicher (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23):
  [`BEO-012`](../observations.md) ist **zentral** — dieser Slice besteht ganz
  aus der Frage, wie weit ein Kanon-Satz reicht, und genau daran ist die Welle
  schon fünfmal gescheitert. [`BEO-011`](../observations.md) für jede Aussage
  über den Bestand der Spec-Straten.

Slice-ID: slice-136. Betroffene IDs: — (Harness-Regeltext; keine Anforderung).
Module: Harness-Regeltext, Konventionsspeicher. Gates: `make doc-check`,
`make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Prüfung eigener Regeln gegen die Baseline.

## 9. Closure-Notiz (nach `done/`)

**Die Antwort dieses Slice ist die umgekehrte seiner ersten Fassung**, und das
ist sein Ergebnis.

**Doppelung? Nein — die Frage war falsch gestellt.** §3.4s Abwärts-Sperre steht
im Kanon **und** ist bei uns gate-gedeckt. Für genau diese Lage verlangt
`modul-09` **beides**: *„Jede Hard Rule liegt in zwei Quadranten … nur als
Fitness Function ohne `AGENTS.md`-Eintrag versteht der Agent das Warum nicht"*
— und gegen die Fehlannahme direkt: *„Beides ist Pflicht."* Der Eintrag ist die
**vorgeschriebene** Feedforward-Hälfte, nicht die überflüssige Kopie.

**Verschärfung? Ja — und ich hatte zuerst das Gegenteil geschlossen.** Der erste
Anlauf stützte sich auf die *Praxis* der Architektur-Vorlage: null
Code-Pfad-Token, Rollen statt Pfade. Die direkteste Quelle blieb ungelesen —
[`AGENTS.template.md`](../../../../.harness/baseline/v5.11.0/templates/AGENTS.template.md)
§3.4, das Baseline-Pendant zu unserer eigenen Datei, sagt wörtlich:

> `spec/architecture.md` **referenziert Modul-Pfade**, aber keine Wellen,
> Slices, Commit-Hashes oder Closure-Daten.

Der Kanon **erlaubt** der Sicht also, was wir ihr **verbieten**. Das ist eine
Verschärfung; [`MR-031`](../../../../harness/conventions.md#mr-031) verlangt
dafür einen Eintrag, und
[`MR-033`](../../../../harness/conventions.md#mr-033) trägt ihn nach.

**Der Slice hat sein eigenes Risiko wörtlich erfüllt.** §5 stand geschrieben,
bevor es zutraf: *„die bequeme Lesart macht unsere Regel konform … wer sie
wählt, weil sie weniger Arbeit macht, hat `BEO-012` begangen."* Ich habe
indirekte Evidenz genommen, die mein Ergebnis stützte, und die direkte nicht
gesucht. **Ein Risiko zu benennen ersetzt nicht, es zu prüfen** — was hier
geholfen hätte, war keine Warnung, sondern eine Suche in der richtigen Datei.

**Und ein zweiter Fund, der eine Wiederholung ist.** §3.4s Abwärts-Sperre nennt
**fünf** Kategorien — ADRs, Wellen, Slices, Commit-Hashes, Closure-Daten.
`matrix` trägt Klassen für **zwei** davon; für Wellen, Commit-Hashes und
Closure-Daten gibt es weder Klasse noch Muster. Meine Zensus-Probe prüfte genau
eine (`slice-999`). Der Review zu
[slice-132](slice-132-hard-rule-zensus.md) hatte **dieselbe Form** einen Slice
zuvor als HIGH gemeldet, im **selben Abschnitt**. §3.4 sagt jetzt *zwei von
fünf* und trägt einen auflösenden Trigger für die drei fehlenden.

**Zwei Zahlen sind aus dem Regeltext verschwunden**, weil sie suchmuster-
abhängig waren und die Aussage nicht trugen, auf die ich sie gestützt hatte: die
*„null Code-Pfad-Token"* der Vorlage übersehen ihren Bedienhinweis-Block, und
die *„fünf"* in der Spezifikation sind ausschließlich `tools/*.sh`-Tombstones,
die `codepaths.ignore-refs` ohnehin ausnimmt.

**Offen und benannt:** Die Mehrdeutigkeit des Kanon-Bullets bleibt bestehen —
*„referenziert Modul-Pfade"* sagt nicht, ob Code- oder Dokument-Pfade gemeint
sind. Für uns ist die Frage durch [`MR-033`](../../../../harness/conventions.md#mr-033) entschieden; für den Kanon ist sie
ein **CR-Kandidat**, weil sie jedes Adopter-Repo trifft. Ob und wann er
geschrieben wird, ist eine Auftraggeber-Entscheidung und ausdrücklich nicht Teil
dieses Slice.
