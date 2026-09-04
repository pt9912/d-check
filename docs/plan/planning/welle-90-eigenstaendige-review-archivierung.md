# Welle welle-90: Eigenständige Review-Archivierung

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach in `docs/plan/planning/`; bei der Closure wandert sie nach `done/` und
bekommt ihre Ergebnisnotiz `welle-90-results.md` daneben.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-04.

---

## 1. Welle-Ziel

**`docs/reviews/` trägt nach welle-89 noch 11 Review-Reports ohne
`slice-<NNN>` im Dateinamen** — CR-, ADR-, Baseline-/MR-, Backlog-, Wellen-
und Release-Prep-Reviews. `tools/archive-wave`s Sammel-Logik (beide bisherige
Modi) findet Review-Reports ausschließlich über ein solches Dateinamens-
Muster; diese 11 sind für sie strukturell unsichtbar. Als
[`BEO-ALL/review-collection-misses-non-slice-filenames`](observations/BEO-ALL/review-collection-misses-non-slice-filenames/observation.md)
registriert (Zähler 1, unter der Schwelle) — Auftraggeber-Entscheid: diesen
Fund jetzt beheben, statt auf die 3×-Schwelle zu warten.

**Design-Entscheidung, die den Umfang trägt:** Ein Review-Report zu einem
Slice bekommt beim Archivieren keinen Stub — seine Identität kommt vom
Slice, der seinerseits einen Stub trägt (Baseline-Regelwerk, Kanon-Begründung
für Review-Reports). Ein **eigenständiger** Review ohne Slice-Partner *ist*
selbst der abgeschlossene Vorgang — er bekommt deshalb, wie ein Slice oder
eine Welle, einen eigenen Stub an seiner Stelle, statt spurlos zu
verschwinden. Das ist eine neue Adaption (kein Kanon-Fall), analog zu
[ADR-0083](../adr/0083-beobachtungsregister-verzeichnis-modus.md)s additiver
Erweiterung.

**Zwei Slices, nach dem welle-87/89-Muster (Werkzeug bauen, dann
anwenden):**

- **slice-198** erweitert `tools/archive-wave` um einen dritten Modus
  (`-review=<dateiname>`, mutually exclusive zu `-welle`/`-slice`): ein
  eigenständiger Review-Report ohne `slice-<NNN>` im Dateinamen wird nach
  `docs/reviews/archiv/<basisname>-archiv.zip` archiviert, gekürzter Stub
  im selben Verzeichnis. An einem konstruierten Fixture bewiesen.
- **slice-199** wendet den neuen Modus auf alle 11 eigenständigen Reviews
  an, die nach welle-89 noch in `docs/reviews/` liegen.

## 2. Trigger (Welle startet)

Kein Vorwellen-Trigger nötig — der Anlass ist die eigene, bereits gemeldete
Beobachtung. `in-progress/` ist leer, das WIP-Limit ist frei.

## 3. Closure-Trigger (Welle schließt)

- slice-198, slice-199 beide in `done/`.
- `docs/reviews/` enthält (außerhalb `docs/reviews/archiv/`) <!-- d-check:ignore (Ziel-Form, entsteht erst mit slice-198) --> keinen Report
  mehr ohne `slice-<NNN>` im Dateinamen.
- `make gates` und `make fullbuild` grün auf dem Endstand.
- `make archive-wave-test` grün.
- Closure-Notiz `welle-90-results.md` geschrieben.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| slice-198 | `tools/archive-wave`: Modus für eigenständige Review-Archivierung | `BEO-ALL/review-collection-misses-non-slice-filenames` |
| slice-199 | Alle 11 eigenständigen Reviews archivieren | slice-198 |

## 5. Abhängigkeiten

- slice-199 wird blockiert von slice-198 (braucht den neuen Werkzeug-Modus).

## 6. Out-of-Scope für diese Welle

- Eine Änderung der beiden bestehenden Modi (`-welle`, `-slice`).
- Reviews mit `slice-<NNN>` im Dateinamen, deren Slice noch nicht archiviert
  ist — bleiben unberührt, gehören in den Slice-Modus, sobald ihr Slice
  archivierbar wird.
- Eine rückwirkende Umbenennung der 11 Reviews — sie behalten ihren
  historischen Dateinamen im Archiv.

## 7. Closure-Notiz

<!-- wird erst nach Welle-Abschluss gefüllt -->

Ergebnis: `done/welle-90-results.md`
