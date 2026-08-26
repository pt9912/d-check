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

**Berührte Spec-Stellen:** `spec/lastenheft.md`, §[`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix) — **falls** die
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

- [ ] Die Frage *Zusage oder Beispiel* ist am Text beantwortet, nicht am Wunsch.
- [ ] Der Bestand ist geprüft: alle Stellen, die Klassen über Pfade definieren.
- [ ] Entweder Lastenheft-Änderung mit Bump und Historie — oder die benannte
      Ausnahme steht in der technischen Spezifikation.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Die bequeme Antwort ist, den Vertrag der Umsetzung anzupassen.** Genau
  dieselbe Lage wie in [slice-150](../done/slice-150-pin-gebundene-zitate.md),
  wo sie sich als falsch erwiesen hat. Die Begründung muss aus dem Text kommen.
  — **Ausgang:** *(bei Closure)*
- **Ein Lastenheft-Satz hat Spiegel.** Wer einen ändert, ohne die anderen zu
  zählen, hinterlässt einen Rand, der eine Fassung referiert, die es nicht mehr
  gibt. — **Ausgang:** *(bei Closure)*

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
[`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix).
Module: `matrix`, Spec-Straten. Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Textarbeit am eigenen Vertrag.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
