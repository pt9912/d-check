# Slice slice-141: Ein Fixture, das Änderung über gleich lange Inhalte herstellt

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`BEO-014`](../observations.md); `internal/adapter/driven/git/git_test.go`.

**Berührte Spec-Stellen:** — (Testdaten; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-25.

---

## 1. Ziel

`TestRangeAndFileAt` schreibt dieselbe Datei von `"v1\n"` auf `"v2\n"` und
committet sofort. Beide Fassungen sind **drei Byte** lang und entstehen in
**derselben Sekunde**; eine stat-basierte Änderungserkennung darf die Datei dann
als unverändert führen. Der Test ist damit zeitabhängig — einmal beobachtet:
`make gates` rot, der `make fullbuild` unmittelbar danach auf demselben
Arbeitsbaum grün.

**Warum das mehr ist als ein Testfehler:** der Test hängt in `make gates`. Ein
Gate, das gelegentlich ohne Grund rot wird, erodiert die Zusage „grün heißt
geprüft" schneller als eines, das gar nicht existiert.

## 2. Vorgehen

1. Die Klasse **im Bestand suchen**, nicht nur den einen Fall: welche Fixtures
   stellen Änderung über gleich lange Inhalte in einem Zug her?
2. Inhalte auf **unterschiedliche Länge** bringen — das ist die kleinste
   Änderung, die die Bedingung entfernt.
3. **Belegen, dass es die Ursache war**, nicht nur dass es jetzt grün ist: der
   reparierte Test muss bei künstlich gleichgehaltener Länge weiterhin
   sporadisch scheitern können.
4. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Retry-Schleife.** Einen zeitabhängigen Test zu wiederholen, bis er
  grün wird, macht ihn nicht deterministisch, sondern unsichtbar.
- **Kein Umstieg der Bibliothek.** Die Ursache liegt im Fixture, nicht im
  Adapter.

## 4. Definition of Done

- [x] Die Klasse ist im Bestand **gemessen**; jede Fundstelle ist benannt.
      *(Sechs gleich lange Neuschreibungen derselben Datei in vier
      Testfunktionen — nicht eine.)*
- [x] Die Inhalte unterscheiden sich in der Länge; die Begründung steht im Test.
      *(Und sie steht dort, wo sie greift: der Fixture-Helfer lehnt eine gleich
      lange Neuschreibung ab.)*
- [x] Der Ursachen-Beleg ist geführt, nicht nur der grüne Lauf. *(Gelesene
      Meldung: `cannot create empty commit: clean working tree`.)*
- [x] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein einmal beobachteter Flake ist schwer zu reproduzieren.** Wer ihn nicht
  herstellen kann, kann auch nicht belegen, dass er weg ist — die Reparatur
  bliebe eine Vermutung. — **Ausgang: entfallen, aber knapp.** Die Bedingung
  ließ sich herstellen: gleiche Länge **und**, per `Chtimes`, dieselbe
  Änderungszeit auf die Nanosekunde. Sie reproduziert — allerdings **nicht
  zuverlässig**, und genau daran wäre der Slice fast gescheitert: die ersten
  zwei Läufe der erzwungenen Probe waren grün und wurden als Widerlegung
  gelesen. Der dritte meldete `cannot create empty commit: clean working tree`.
  Belegt ist die Ursache damit über die **Meldung**, nicht über einen grünen
  Lauf; die Klasse hinter dem Beinahe-Fehlschluss trägt
  [`BEO-019`](../observations.md).

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls sich zeigt, dass die Ursache im Adapter liegt und nicht im Fixture.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Testdaten (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-25): [`BEO-014`](../observations.md) ist der Anlass; [`BEO-011`](../observations.md) für die Aussage, wie viele Fixtures die Klasse tragen.

Slice-ID: slice-141. Betroffene IDs: — (Testdaten; keine Anforderung). Module: Test-Fixtures.
Gates: `make test`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Reparatur an eigenen Testdaten.

## 9. Closure-Notiz (nach `done/`)

**Die Ursache ist bestätigt, und der Weg dorthin ist die eigentliche Lehre.**
Die Änderungserkennung von `go-git` ist stat-basiert: schreibt ein Test dieselbe
Datei zweimal mit gleich langem Inhalt in einem Zug, kann der zweite Commit den
Baum als sauber sehen. Gelesene Meldung:
`cannot create empty commit: clean working tree`.

**Die Klasse ist breiter als der Plan.** Nicht ein Fall, sondern **sechs** —
`keep.md` in `TestRangeAndFileAt`, `adr.md` in `TestStaged` und je zwei
`f.md`-Paare in `TestCommitMessages` und `TestCommitMessagesSkipsMerges`. Alle
sechs tragen jetzt unterschiedlich lange Inhalte.

**Beinahe wäre das Gegenteil geliefert worden.** Die erzwungene Probe — gleiche
Länge **und**, per `Chtimes`, dieselbe Änderungszeit auf die Nanosekunde — lief
**zweimal grün**. Daraus wurde geschlossen, die Deutung sei widerlegt und die
Reparatur bliebe eine Vermutung; die Closure-Notiz stand bereits so da. Der
**dritte** Lauf desselben unveränderten Codes meldete den Fehler. Zwei grüne
Läufe eines nichtdeterministischen Prüflings sind kein Ergebnis — nur das Rot
trug eine Aussage. Als Klasse geführt in [`BEO-019`](../observations.md).

**Der Wächter steht dort, wo er nicht auf Zeit angewiesen ist.** Ein Test, der
auf die Kollision wartet, meldet die Klasse nur manchmal — er wäre selbst der
Fehler, den er sucht. Stattdessen lehnt der Fixture-Helfer `put` eine gleich
lange Neuschreibung derselben Datei **beim Schreiben** ab, mit Begründung in der
Meldung. Das ist deterministisch, greift beim Autor statt beim Lauf, und es
verhindert die Rückkehr der Klasse in jedem künftigen Fixture.

**Belegt statt behauptet:** die Gegenprobe setzt eine der sechs Stellen auf
gleiche Länge zurück und meldet
*„Fixture `keep.md` wird mit gleich langem Inhalt neu geschrieben (3 Byte)"*;
der Bestand ist danach unverändert. Und `make test` lief **fünfmal** ohne
Fehlschlag — nicht als Beweis der Stabilität (siehe `BEO-019`), sondern als
Gegenprobe auf dem reparierten Bestand.

**Der Rückführungs-Trigger greift nicht.** §6 sieht `next` vor, *„falls sich
zeigt, dass die Ursache im Adapter liegt und nicht im Fixture"*. Sie liegt im
Zusammenspiel: der Adapter verhält sich wie dokumentiert, das Fixture verlässt
sich auf etwas anderes. Repariert ist das Fixture — und die Annahme trägt jetzt
einen Wächter.

**Sensors:** `make gates` (Exit 0, zehn Glieder), `make test` (fünf Läufe, Paket
`driven/git` je grün), Gegenprobe mit gelesener Meldung und
byte-identischem Rückbau.
