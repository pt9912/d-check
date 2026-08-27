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

- [x] Die Vorfrage ist **gemessen** beantwortet, nicht angenommen.
- [x] ~~Bei Umsetzung: Befund-Code, Bestandsmessung, konstruierter Verstoß mit
      gelesener Ursache.~~ **Entfällt — nicht umgesetzt.**
- [x] Bei Nicht-Umsetzung: [`BEO-013`](../observations.md) trägt die Messung als
      Begründung.
- [x] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Der Ertrag ist klein.** Zwölf Marker, elf davon unbeweglich — der Slice
  könnte mehr kosten als er bringt. Das gehört vor der Umsetzung entschieden,
  nicht danach bedauert. — **Ausgang: entfallen, weil die Prämisse falsch war.**
  Es sind nicht zwölf Marker, sondern **233** unterdrückende Zeilen; der Ertrag
  wäre also groß. Gebaut wird trotzdem nicht — aus dem umgekehrten Grund:
  **185** davon unterdrücken nichts, und die Regel meldete sie **zu Recht und
  nutzlos**, weil **146** blanke Erwähnungen in Doku-Prosa sind. Die Abwägung
  Aufwand-gegen-Ertrag wurde damit nie gebraucht; entschieden hat ein
  tieferliegender Defekt, ausgetragen als
  [slice-159](../open/slice-159-ignore-marker-syntax.md).

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

**Die Vorfrage ist beantwortet, und die Antwort ist zweiteilig.** Das Produkt
hält die Information **nicht**: in allen vier Konsumenten — `codepaths`, `ids`,
`versions`, `diagrams` — ist der Marker ein frühes `continue` **vor** der
Prüfung, die markierte Zeile wird also nie ausgewertet. Erreichbar wäre sie
durch Umstellung auf *prüfen, dann unterdrücken*; unmöglich ist sie nicht.

**Gebaut wird sie trotzdem nicht — und der Grund ist nicht der, den §5
erwartet hat.** Der Abnahme-Punkt fürchtete einen zu kleinen Ertrag bei zwölf
Markern. Die Prämisse ist falsch: gemessen sind **233** Zeilen außerhalb von
Fences, die die Zeichenkette tragen und damit unterdrücken. Der Ertrag wäre
also groß.

**Er wäre nur wertlos.** Die Erkennung ist eine blanke Teilketten-Suche
(`strings.Contains`), und deshalb wirkt **jede Erwähnung** des Markers wie ein
gesetzter Marker: **146** der 233 Zeilen sind Doku-Prosa, die über das Ventil
schreibt. Eine Regel „dieser Marker unterdrückt nichts" meldete davon **185** —
richtig und nutzlos. Erst muss ein Marker ein Marker sein.

**Die eigentliche Zahl steht damit fest, und sie ist unangenehm.** Lenkt man die
Marker-Konstante ins Leere, treten **58** Befunde auf **48** Zeilen hervor:
38 `id-unlinked`, 18 `codepath-missing`, 2 `repo-escape`. Diese 48 Zeilen sind
die einzigen, an denen der Marker heute wirklich etwas tut. Die übrigen 185
schweigen über nichts — oder über etwas, das niemand geprüft hat.

**Das ist dieselbe Klasse wie
[slice-158](../open/slice-158-citations-inline-code.md), mit umgekehrtem
Vorzeichen:** ein Modul, das seine eigene Dokumentation liest. Dort bricht es
laut ab und fällt beim ersten Lauf auf; hier schweigt es lautlos und ist seit
der Einführung des Ventils unbemerkt. Ausgetragen als
[slice-159](../open/slice-159-ignore-marker-syntax.md);
[`BEO-013`](../observations.md) trägt die Messung und wartet darauf.

**Was dieser Slice nicht geprüft hat, und es gehört gesagt:** ob die 48 wirkenden
Unterdrückungen berechtigt sind. Der Slice fragt nach Markern, die **nichts**
unterdrücken — nicht nach solchen, die **zu Unrecht** unterdrücken. Die zweite
Frage ist ein Urteil und bleibt eines.

**Sensors:** `make gates` (Exit 0, zehn Glieder), `make doc-check` (535 Dateien,
0 Befunde). Die Messungen liefen mit dem Produkt gegen eine temporär geänderte
Marker-Konstante; der Quellstand ist danach byte-identisch (`git diff` leer).
