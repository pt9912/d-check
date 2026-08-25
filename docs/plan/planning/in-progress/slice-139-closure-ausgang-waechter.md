# Slice slice-139: Ein Risiko ohne Ausgang darf nicht in `done/` liegen

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos** (Baseline-Regelwerk
[`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht](../../../../.harness/baseline/v5.11.0/regelwerk/modul-06-roadmap.md)):
seine Closure-Bedingung wäre seine eigene DoD.

**Bezug:** Baseline-Regelwerk
[`modul-05-planning-harness.md` §Offene Risiken werden bei Closure aufgelöst](../../../../.harness/baseline/v5.11.0/regelwerk/modul-05-planning-harness.md)
— *„Ein Slice geht nicht nach `done/`, während ein Risiko ohne Ausgang
dasteht."* Dazu [`BEO-015`](../observations.md);
[ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md) (die
Closure-Notiz-Struktur im Modul `planning`);
[`AGENTS.md`](../../../../AGENTS.md) §4.

**Berührte Spec-Stellen:** — (Harness-Gate; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

Der Kanon verlangt, dass **jeder** offene Punkt beim Übergang nach `done/` genau
einen von drei Ausgängen bekommt, und schließt mit einem harten Satz: *„Ein
Slice geht nicht nach `done/`, während ein Risiko ohne Ausgang dasteht."*

**Diese Regel ist heute vollständig ungewacht** — obwohl ihr Verstoß die
maschinenlesbarste Form hat, die es gibt: ein wörtlicher Vorlagen-Platzhalter in
einer Datei an einem bekannten Ort. Gemessen:

- Ein `done/`-Slice mit unaufgelöstem `*(bei Closure)*` in §5 läuft durch
  `make verify-closure-notes` mit **0 Befunden**.
- Auch nachdem `"(bei Closure)"` in die `boilerplate`-Liste aufgenommen wurde:
  weiterhin **0** — die Liste greift nur in der Closure-Notiz, nicht im
  Risiko-Abschnitt.
- `structure` kennt `non-empty`, `table-order`, `headings-match` — **kein**
  „darf Muster X nicht enthalten".

**Der Bestand ist sauber**, dreifach gemessen: null `*(bei Closure)*`, null
`*(wird mit dem Closure-Body gefüllt)*`, null `<…>` in `done/`-Slices. Der
Wächter startet also grün und wirkt ab dem ersten Verstoß.

## 2. Vorgehen

1. **Ein Skript**, das `done/slice-*.md` gegen die drei Platzhalter-Formen hält
   — die zwei repo-lokalen und die `<…>`-Form der Kanon-Vorlage.
2. **Fail-closed bei leerer Prüfmenge**: findet es keine `done/`-Slices, bricht
   es rot ab. „Nichts gefunden" und „nichts zu prüfen" dürfen im Exit nicht
   gleich aussehen.
3. **Ort:** ein eigenes `make`-Target, gehängt an **`fullbuild`** neben
   `verify-closure-notes` — nicht an `gates`. Begründung: die Regel gilt dem
   **Übergang nach `done/`**, und das ist der Closure-Bindepunkt; dieselbe
   Einordnung, die `completeness-check` und `verify-closure-notes` schon tragen.
4. **Bewusstes Brechen** je Platzhalter-Form: gesetzt ⇒ rot mit **gelesener
   Fundstelle**; Rückbau ⇒ grün. Plus die leere Prüfmenge.
5. `AGENTS.md` §4 und die Sensors-Tabelle nachziehen — sonst
   `gate-consistency`-rot.
6. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Modul-Fähigkeit.** Die abschnittsgenaue Form (`structure` mit einem
  `forbid-match`) ist der spätere, sauberere Weg — Produkt-Delta mit ADR, und
  sie gehört **hinter** einen Zensus über die Closure-Prozedur.
- **Keine Ausweitung auf `boilerplate`.** Den ganzen Slice gegen die
  Floskel-Liste zu halten erzeugte Falsch-Positive: ein Slice-Plan darf
  „fertig" in Prosa enthalten.
- **Keine Prüfung, ob der Ausgang inhaltlich trägt.** Ob *„eingetreten"* die
  richtige Antwort war, ist Urteil. Geprüft wird, dass **überhaupt** einer
  dasteht.

## 4. Definition of Done

- [ ] Das Target hält `done/slice-*.md` gegen **alle drei** Platzhalter-Formen
      und ist **fail-closed bei leerer Prüfmenge**.
- [ ] Es hängt an `fullbuild`, **nicht** an `gates` — mit der Begründung im
      Target-Kommentar.
- [ ] **Vier** konstruierte Verstöße rot gesehen, jeder mit gelesener
      Fundstelle: drei Platzhalter-Formen und die leere Prüfmenge; Rückbau grün.
- [ ] `AGENTS.md` §4 und die Sensors-Tabelle tragen das Target;
      `gate-consistency` grün.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein Wächter auf eine Zeichenkette ist so gut wie seine Liste.** Ändert
  jemand die Platzhalter-Form im Slice-Kopf, greift er nicht mehr — und
  schweigt. Die Liste gehört an die Vorlage gebunden, nicht geraten. —
  **Ausgang:** *(bei Closure)*
- **`fullbuild` statt `gates` heißt: der Verstoß fällt später auf.** Ein
  falsch geschlossener Slice könnte einen Commit lang unbemerkt bleiben. Das ist
  die bewusste Kehrseite der Einordnung als Closure-Bindepunkt und gehört
  benannt, nicht wegargumentiert. — **Ausgang:** *(bei Closure)*
- **Der Wächter prüft nur `done/`.** Ein Slice, der nie dorthin wandert, ist
  ihm gleichgültig — und das ist richtig so, aber es heißt auch: er sagt nichts
  über den Bestand in `open/`. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls sich zeigt, dass die
Platzhalter-Formen im Bestand uneinheitlich sind — dann ist ihre Vereinheitlichung
ein eigener Slice und der Wächter hinge an einer geratenen Liste.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Werkzeuge (GF), Gate-Landschaft (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23):
  [`BEO-015`](../observations.md) ist der **Anlass** dieses Slice — er schließt
  die Feedback-Hälfte zu genau der Regel, deren vierten Ausgang der Eintrag
  benennt. [`BEO-007`](../observations.md) für jeden Beleg-Lauf.
  [`BEO-010`](../observations.md), weil ein neues Target in drei Doku-Flächen
  erscheint.

Slice-ID: slice-139. Betroffene IDs: — (Harness-Gate; keine Anforderung).
Module: Harness-Werkzeuge, Gate-Landschaft. Gates: `make gate-consistency`,
`make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — neuer Wächter auf eigenem Bestand.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
