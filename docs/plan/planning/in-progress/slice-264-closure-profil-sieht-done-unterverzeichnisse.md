# slice-264: Das Closure-Profil prüft die Slices in den Unterverzeichnissen von `done/`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** geteilt aus slice-263, der die Produkt-Schlüssel liefert; Befund
F-1 (HIGH) aus dem Review von slice-260;
[ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md).

**Berührte Spec-Stellen:** [`SPEC-095`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)
(der Satz zum ausgelösten Lauf).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** `make verify-closure-notes` prüft jeden geschlossenen Slice im
Volltext, gleich ob er direkt unter `done/`, unter `done/wellenlos/` oder unter
einem Wellen-Verzeichnis liegt, und lässt die Stubs aus — mit den Schlüsseln
aus slice-263 in `.d-check.closure.yml`. Die vier Altverstöße, die beim
Schnitt von slice-263 gemessen wurden (Risiko-Ausgang außerhalb des
Wortschatzes in slice-240 bis slice-243), sind behoben, und die Aussagen, dass
der Lauf die Unterverzeichnisse nicht prüft, sind zurückgenommen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Produkt-Schlüssel** — slice-263.
- **Die Archivierung der 23 Volltexte unter `done/wellenlos/`** — ein eigener
  Vorgang am Bestand, kein Teil der Prüfung.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] `.d-check.closure.yml` trifft die Volltexte in den Unterverzeichnissen
      und nimmt die Stubs an ihrem Inhalt aus; bewusst gebrochen: ein
      kaputter Volltext unter `done/wellenlos/` macht den Lauf rot, ein Stub
      nicht.
- [ ] Die Altverstöße sind behoben; `make verify-closure-notes` und
      `make gates` grün.
- [ ] Zurückgenommen: der Satz in [`SPEC-095`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge), Grenze 4 und Vertrag von
      `harness/sensors/hooks.md`, Vertrag „direkt", Grenze 8 und Bindung von
      `harness/sensors/verify-closure-notes.md`, der GRENZE-Kommentar in
      `.githooks/pre-commit`.
- [ ] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.d-check.closure.yml` | update | Rekursion, Globs, Stub-Ausnahme |
| `done/wellenlos/slice-240` bis `slice-243` | update | Altverstöße |
| `spec/spezifikation.md` ([`SPEC-095`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)), `harness/sensors/{hooks,verify-closure-notes}.md`, `.githooks/pre-commit` | update | Rücknahme der Grenz-Aussagen |

## 4. Trigger

**Start** (`next` → `in-progress`): slice-263 in `done/`; `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `open` (blockiert): die Prüfung der Volltexte zeigt mehr
  Altverstöße, als ein Slice beheben sollte — dann die Behebung in einen
  eigenen Slice, mit befristeter Ausnahme.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- **Weitere Altverstöße** — die Messung beim Schnitt galt nur dem Modul
  `structure`; die `planning`-Hälfte (Platzhalter, Floskeln, dünne Notiz) kann
  weitere zeigen. — **Ausgang:** *(offen)*

## 7. Closure-Notiz

*(gefüllt vor dem `git mv` nach `done/`)*

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind die Closure-Konfiguration,
vier `done/`-Slices und die Harness-Doku samt Spezifikation §7 — alle unter dem
Default `*` (`ALL`). Eine eigene Sub-Area erfüllt keine davon: Keine trägt eine
eigene Konvention, einen eigenen Modus oder eine eigene Inventur-Linie.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Zwei berühren den Slice.
[`BEO-ALL/messmethode-scope-enger-als-dokumentierte-ziel-form`](../observations/BEO-ALL/messmethode-scope-enger-als-dokumentierte-ziel-form/observation.md)
(1×) beschreibt genau die Lage, die dieser Slice schließt — ein Scan-Bereich,
der enger ist als das, wofür sein Grün gelesen wird; der Slice ist ihre
Behebung, kein zweites Auftreten.
[`BEO-ALL/leeres-verzeichnis-lokal-verdeckt-ci-befund`](../observations/BEO-ALL/leeres-verzeichnis-lokal-verdeckt-ci-befund/observation.md)
(1×) betrifft das bewusste Brechen: Es wird auch im frischen Klon gemessen.
Keine erreicht mit diesem Slice 3×.

**Messung beim Beanspruchen:** Mit Rekursion, `**/`-Globs und der
Stub-Ausnahme `(?m)^> \*\*ARCHIVIERT` in einer Wegwerf-Kopie des Profils
meldet der Lauf über 974 Dateien genau **vier** Befunde — den Risiko-Ausgang
`*(offen)*` in §6 von slice-240 bis slice-243. Die `planning`-Hälfte meldet
nichts, die `reviews`-Hälfte ebenfalls nichts. Unter `done/wellenlos/` liegen
23 Volltexte; die Wellen-Verzeichnisse tragen nur Stubs.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`upstream-drift` grün (2026-10-09T07:17Z). `image-scan` rot (2026-10-09T10:37Z)
— der Lauf liegt vor dem Release v0.85.0, das die gemeldeten CVEs behebt;
der Slice berührt das Image nicht.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
