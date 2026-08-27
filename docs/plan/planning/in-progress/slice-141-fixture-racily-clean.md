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

- [ ] Die Klasse ist im Bestand **gemessen**; jede Fundstelle ist benannt.
- [ ] Die Inhalte unterscheiden sich in der Länge; die Begründung steht im Test.
- [ ] Der Ursachen-Beleg ist geführt, nicht nur der grüne Lauf.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein einmal beobachteter Flake ist schwer zu reproduzieren.** Wer ihn nicht
  herstellen kann, kann auch nicht belegen, dass er weg ist — die Reparatur
  bliebe eine Vermutung. — **Ausgang:** *(bei Closure)*

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

*(wird mit dem Closure-Body gefüllt)*
