# Slice slice-146: Ein `d-check:ignore`, der nichts mehr unterdrückt

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`BEO-013`](../observations.md); [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix); `nolintlint` mit `allow-unused` als Vorbild aus [slice-134](../done/slice-134-nolintlint.md).

**Berührte Spec-Stellen:** `spec/lastenheft.md` — falls das Produkt die Information nicht ohnehin hält, wächst eine Anforderung; Bump und Historie nach [`MR-032`](../../../../harness/conventions.md#mr-032).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-25.

---

## 1. Ziel

Ein Unterdrückungs-Marker wird gegen **einen** konkreten Befund gesetzt.
Verschwindet dieser, bleibt der Marker — und schweigt von da an über etwas, das
niemand stummschalten wollte. Für Go ist die Frage seit `nolintlint`
(`allow-unused: false`) gemessen; für `<!-- d-check:ignore -->` gibt es kein
Gegenstück. Bestand: **zwölf** aktive Marker, elf eingefroren in `done/`,
**einer** in einem lebenden Dokument.

**Die erste Frage ist eine Messung, keine Konstruktion:** hält das Produkt an
einer markierten Zeile überhaupt die Information, ob ohne den Marker ein Befund
entstünde?

## 2. Vorgehen

1. **Messen**, ob die Auswertungsreihenfolge diese Information hergibt. Wenn
   nein, ist das die Antwort — und der Slice endet mit einer Ausweisung.
2. Wenn ja: Befund-Code und Ventil entwerfen, am Bestand messen (zwölf Marker,
   erwartet null Treffer).
3. Konstruierter Verstoß: ein Marker an einer Zeile ohne Befund ⇒ rot.
4. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Räumung der elf eingefrorenen Marker.** Sie stehen in `done/` und
  können nicht altern.
- **Kein neues Ventil ohne Bedarf.** Zwölf Marker sind ein kleiner Bestand; der
  Ertrag gehört gegen den Aufwand gestellt.

## 4. Definition of Done

- [ ] Die Vorfrage ist **gemessen** beantwortet, nicht angenommen.
- [ ] Bei Umsetzung: Befund-Code, Bestandsmessung, konstruierter Verstoß mit
      gelesener Ursache.
- [ ] Bei Nicht-Umsetzung: [`BEO-013`](../observations.md) trägt die Messung als
      Begründung.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Der Ertrag ist klein.** Zwölf Marker, elf davon unbeweglich — der Slice
  könnte mehr kosten als er bringt. Das gehört vor der Umsetzung entschieden,
  nicht danach bedauert. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Vorfrage ein Produkt-Delta verlangt, das den Aufwand nicht trägt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Module (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-25): [`BEO-013`](../observations.md) ist der Anlass; [`BEO-011`](../observations.md) für die Aussage über den Marker-Bestand.

Slice-ID: slice-146. Betroffene IDs: [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix). Module: Produkt-Module.
Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Messung vor Konstruktion an bestehender Modul-Mechanik.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
