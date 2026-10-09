# slice-263: Die Closure-Prüfung sieht die Slices in den Unterverzeichnissen von `done/`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** Befund F-1 (HIGH) aus dem Review von slice-260;
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(Closure-Fähigkeit von `planning`),
[ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md);
Auftraggeber-Entscheid 2026-10-09 (eigener Slice statt Mitnahme in slice-260).

**Berührte Spec-Stellen:** [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) (Closure-Kandidaten,
`planning.closure.dir`), `spec/spezifikation.md` §2 (Schlüssel der
`planning`-Konfiguration).

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** `make verify-closure-notes` prüft jeden geschlossenen Slice im
Volltext, gleich ob er direkt unter `done/`, unter `done/wellenlos/` oder unter
einem Wellen-Verzeichnis liegt — die Closure-Notiz (Modul `planning`) ebenso wie
die Abschnitts-Invarianten (Modul `structure`). Archivierte Stubs bleiben
ausgenommen, an ihrer Form erkannt, nicht an ihrem Verzeichnis. Gemessen beim
Schnitt: seit 2026-09-29 liegen 17 Volltexte unter `done/wellenlos/`, keiner
davon wurde geprüft; vier tragen einen Risiko-Ausgang außerhalb des
geschlossenen Wortschatzes (slice-240 bis slice-243).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Auslöser im Hook und in der CI** — slice-260 hat ihn bereits auf jeden
  `slice-*.md` unter `done/` erweitert.
- **Andere Module mit einem `done/`-Verzeichnis** (`reviews.done-dir`) — beim
  Beanspruchen messen; trifft dieselbe Blindheit zu, ist das ein eigener
  Befund mit eigenem Slice.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] `planning` liest die Closure-Kandidaten auch aus Unterverzeichnissen
      (rekursiv oder als Liste — Entscheidung im Slice, mit Anforderung oder
      Verfeinerung); Test, der ohne die Änderung aus dem richtigen Grund rot ist.
- [ ] `.d-check.closure.yml` trifft die Volltexte in den Unterverzeichnissen
      und nimmt die Stubs an ihrer Form aus; die vier Altverstöße sind
      behoben; `make verify-closure-notes` grün.
- [ ] Ohne die neue Konfiguration ist die Ausgabe unverändert
      (`make blackbox-probe`); `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Modul `planning` (Closure-Kandidaten) | update | Unterverzeichnisse |
| `spec/lastenheft.md` / `spec/spezifikation.md` | update | Anforderung bzw. Verfeinerung |
| `.d-check.closure.yml` | update | Globs und Stub-Ausnahme |
| `done/wellenlos/slice-240` bis `slice-243` | update | Altverstöße |
| [`SPEC-095`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge) (Satz zum ausgelösten Lauf), `harness/sensors/hooks.md` (Vertrag, Grenze 4), `harness/sensors/verify-closure-notes.md` (Vertrag „direkt", Grenze 8, Bindung), GRENZE-Kommentar in `.githooks/pre-commit` | update | die Aussagen, dass der Lauf die Unterverzeichnisse nicht prüft, werden mit der Behebung zurückgenommen |

## 4. Trigger

**Start** (`next` → `in-progress`): `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): die Stub-Erkennung braucht ein eigenes
  Kriterium im Produkt — dann Produkt und Config trennen.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft. Produkt-Verhalten — geht mit dem nächsten Release
hinaus.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** das Produkt und die
Closure-Konfiguration unter dem Default `*` (`ALL`); beim Beanspruchen neu
prüfen.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** beim Beanspruchen neu lesen.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
beim Beanspruchen aus dem jüngsten Lauf lesen.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
