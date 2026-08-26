# Slice slice-153: Sagt das Lastenheft noch, was `matrix` tut?

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix);
[`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix);
[`MR-032`](../../../../harness/conventions.md#mr-032) (Bump und Historie);
[slice-144](../done/slice-144-commit-hash-muster.md) (der Anlass).

**Berührte Spec-Stellen:** `spec/lastenheft.md` — die Beschreibung und die
Akzeptanzkriterien von §[`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix) sowie der Glossar-Eintrag
*Dokumentklasse*; §[`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
bleibt unverändert — **falls** die
Antwort so ausfällt; Bump und Historie dann nach
[`MR-032`](../../../../harness/conventions.md#mr-032).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-26.

---

## 1. Ziel

[`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
beschreibt: *„Die Konfiguration deklariert Dokumentklassen über **Pfad-Muster**"*.
Seit [slice-144](../done/slice-144-commit-hash-muster.md) fährt dieses Repo eine
Klasse **ohne** Pfade — ihr Gegenstand ist eine Zeichenkette, kein Dokument. Die
**technische** Spezifikation ist nachgezogen; die **vertragliche** Beschreibung
sagt weiterhin etwas Engeres.

**Das ist eine Vertrags-Frage, keine Implementierungs-Frage** — und darum ein
eigener Slice: Ändert man einen Lastenheft-Satz, weil die Umsetzung ihn
überholt hat, oder war die Umsetzung dann zu weit?

## 2. Vorgehen

1. **Zuerst die Frage richtig stellen:** Beschreibt der Satz eine *Zusage* oder
   ein *Beispiel*? Er steht unter **Beschreibung**, nicht unter
   Akzeptanzkriterien — das ist ein Unterschied, und er gehört gelesen, bevor
   irgendetwas geändert wird.
2. **Den Bestand prüfen:** Gibt es weitere Stellen im Lastenheft, die Klassen
   über Pfade definieren? Eine einzelne Zeile zu ändern und drei stehenzulassen
   wäre die Spiegel-Falle aus [`MR-025`](../../../../harness/conventions.md#mr-025).
3. Fällt die Antwort für eine Änderung: Bump und Historie nach
   [`MR-032`](../../../../harness/conventions.md#mr-032), Formulierung so eng
   wie möglich.
4. Fällt sie dagegen: die Token-Ziel-Klasse bleibt eine **benannte Ausnahme**
   der technischen Spezifikation, und das gehört dort ausgeschrieben.
5. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Rücknahme der Klasse.** Sie ist gemessen, sie greift, und sie
  verschärft — die Frage ist die Beschreibung, nicht das Verhalten.
- **Keine Ausweitung auf andere `DC-*`-Beschreibungen** in diesem Zug.

## 4. Definition of Done

- [x] Die Frage ist am Text beantwortet — **und zwei meiner drei Belege
      trugen nicht** (§9). Getragen wird die Antwort allein von der
      Glossar-Definition, und die ist eine Zusage.
- [x] Der Bestand ist geprüft — **im zweiten Anlauf**. Der erste zählte vier
      Spiegel, weil er nach dem *neuen* Vokabular suchte;
      [`MR-025`](../../../../harness/conventions.md#mr-025) verlangt das `grep`
      nach dem **alten**. Danach sind es vierzehn.
- [x] Lastenheft-Änderung mit Bump (`0.65.4` → `0.66.0`) und Historie-Zeile
      nach [`MR-032`](../../../../harness/conventions.md#mr-032); **zwei neue
      Akzeptanzkriterien**, nicht nur Prosa. Kein CR — `Status: Draft`.
- [x] `make gates` Exit 0 (zehn Glieder), `make fullbuild` Exit 0; unabhängiger
      Review ([Report](../../../reviews/2026-08-26-slice-153-lastenheft-token-ziel-klasse-review.md)),
      blockierend mit vier MEDIUM — alle zehn eingearbeitet, und die
      **Gegenthese des Reviewers hat den Entwurf ersetzt**.

## 5. Abnahme-Punkte / Risiken

- **Die bequeme Antwort ist, den Vertrag der Umsetzung anzupassen.** Genau
  dieselbe Lage wie in [slice-150](slice-150-pin-gebundene-zitate.md),
  wo sie sich als falsch erwiesen hat. Die Begründung muss aus dem Text kommen.
  — **Ausgang:** *eingetreten, halb.* Der Vertrag ist der Umsetzung angepasst
  worden — das war die bequeme Richtung. Sie ist diesmal die richtige, weil die
  Glossar-Definition tatsächlich enger war als die Fähigkeit. Aber ich habe
  **drei** Belege genannt, wo **einer** trägt: das *„z. B."* hedged die
  Aufzählung, nicht den Mechanismus, und das *„zusätzlich"* kontrastiert Token
  gegen **Link**, nicht gegen Pfade. Die Begründung kam aus dem Text — nur zu
  einem Drittel.
- **Ein Lastenheft-Satz hat Spiegel.** Wer einen ändert, ohne die anderen zu
  zählen, hinterlässt einen Rand, der eine Fassung referiert, die es nicht mehr
  gibt. — **Ausgang:** *eingetreten, vollständig, und mit angesagtem Ausgang.*
  Ich habe gezählt und das Falsche gezählt: `grep` nach dem **neuen** Vokabular
  statt nach dem alten. Vier statt vierzehn — zwei davon in der Datei, die ich
  gerade bearbeitete. Das Risiko stand wörtlich in diesem Abschnitt.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Antwort eine
Auftraggeber-Entscheidung verlangt, die nicht vorliegt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Spec-Straten (GF), Konfigurations-Profil (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-26):
  [`BEO-002`](../observations.md) für die Spiegel einer Semantik-Änderung;
  [`BEO-012`](../observations.md) für jede Aussage darüber, was der
  Lastenheft-Satz zusagt.

Slice-ID: slice-153. Betroffene IDs:
[`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)
(die geänderte Anforderung),
[`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
(die geprüfte, unverändert gebliebene).
Module: `matrix`, Spec-Straten. Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Textarbeit am eigenen Vertrag.

## 9. Closure-Notiz (nach `done/`)

Geliefert: Das Lastenheft sagt jetzt, was `matrix` tut — die Definition der
**Dokumentklasse** ist **geweitet** (Regelfall Pfad-Muster, dazu das
**Token-Ziel**), [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix) trägt die Zusage samt **zwei
Akzeptanzkriterien**, Version `0.66.0` mit Historie-Zeile.

**Die Antwort auf die Kernfrage steht, aber auf einem Bein statt auf dreien.**
Ich hatte drei Textbelege genannt. Zwei tragen nicht: das *„z. B."* in
[`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix) hedged die **Aufzählung** der Klassen, nicht den Mechanismus
*„über Pfad-Muster"*; und das *„zusätzlich"* in der Beschreibung derselben Anforderung kontrastiert
Token gegen **Link** (*„bisher nur als Links … daher"*), nicht gegen Pfade.
Getragen wird die Antwort allein von der **Glossar-Definition** — und die
reicht, weil sie ein Definiens ist und kein Beispiel.

**Der Entwurf war falsch, und die Gegenthese des Reviewers war die Lösung.**
Ich hatte der Dokumentklasse einen **Gegenbegriff** danebengestellt: eine
„Token-Ziel-Klasse", die *„damit **keine** Dokumentklasse"* sei. Das erzeugte
einen Widerspruch zu **vierzehn** lebenden Stellen, die weiter von
Dokumentklassen sprechen — darunter die Glossar-Zeile *Referenzmatrix* **eine
Zeile unter meinem eigenen Eintrag**. Die geweitete Definition erzeugt keinen
einzigen. Der Preis meiner Abgrenzung war nirgends abgewogen; der Reviewer hat
ihn ausgerechnet.

**Und die Spiegel habe ich falsch gezählt, während ich behauptete, sie nach
[`MR-025`](../../../../harness/conventions.md#mr-025) gezählt zu haben.** Der
Eintrag sagt ausdrücklich: *der Ableiter ist das `grep` nach dem **alten**
Wortlaut*. Ich habe nach dem **neuen** gesucht — also nach dem, was ich gerade
geschrieben hatte, und damit genau die Stellen nicht gefunden, um die es geht.
Vier statt vierzehn. Eine davon, der Kern-Kommentar `MatrixClass ist eine über
Pfad-Globs deklarierte Dokumentklasse`, hatte mir der auslösende Review mit
Datei und Zeile übergeben.

**Ein Befund traf die Textsorte selbst.** Meine neue Zusage stand nur in der
*Beschreibung* — in genau der Textsorte, die ich einen Absatz zuvor als
nicht-zusagend eingeordnet hatte, um die Frage überhaupt zu beantworten. Der
Minor-Bump stützte sich damit auf ein Wachstum, das kein Kriterium prüft.
[`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix) trägt jetzt beide Richtungen als Kriterium: das Token-Ziel als
**Ziel** (Befund) und als **Quelle** (kein Befund — es hat keine Mitglieder).

**Register:** [`BEO-002`](../observations.md) auf Zähler **4**.
