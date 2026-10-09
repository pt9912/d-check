# slice-270: Die RTM zeigt die Test-Nachweise

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** Auftraggeber-Frage 2026-10-09: `make trace` zeigt nur ADR- und
Slice-Spalten, keine Tests; die Schwester-Repos pg-change-feed und
pgwire-recorder binden ihre Test-Nachweise über `trace.coverage` ein.
Auftraggeber-Entscheid: erzeugt und gewächtert, nicht von Hand gepflegt.

**Berührte Spec-Stellen:** —

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Die RTM dieses Repos (`make trace`) zeigt je Anforderung, welche
Tests sie belegen — in zwei Spalten: **Tests** (die Go-Suite von `make test`)
und **E2E** (`make image-test` gegen das gebaute Image). Je Spalte eine
Abdeckungs-Datei mit der Tabelle `Kennung → Test → Datei:Zeile`, eingebunden
über `trace.coverage` in der `.d-check.yml`. Die Tabelle wird aus den
Testquellen **abgeleitet**: ein Test in `make test` erzeugt sie und vergleicht
sie mit der committeten Datei; weicht sie ab, ist `make test` rot.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine neue Produkt-Fähigkeit** — `trace.coverage` gibt es; der Slice ist
  Konfiguration und ein Test dieses Repos über sich selbst.
- **Das Urteil, ob ein Test seine Anforderung wirklich belegt** — die Tabelle
  zeigt die Deklaration; ob sie trägt, bleibt Review und Verifikation.
- **`make bench` und die Bestandsproben** (`blackbox-probe`) als eigene
  Spalten — kein Anlass; ein eigener Schnitt, wenn die RTM sie braucht.
- **Waisen schließen** — der Slice macht Belege sichtbar, er schreibt keine
  neuen Tests.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] Zwei Abdeckungs-Dateien (Tests, E2E), aus den Testquellen abgeleitet,
      über `trace.coverage` eingebunden; `make trace` zeigt beide Spalten.
- [ ] Ein Test in `make test` hält jede Datei gegen ihre Ableitung — eine
      neue oder entfernte Deklaration ohne nachgezogene Datei ist rot
      (bewusstes Brechen belegt); `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Abdeckungs-Datei Tests, Abdeckungs-Datei E2E | neu | die zwei Spalten |
| Test über die Ableitung (Go, in `make test`) | neu | Wächter gegen Drift |
| `tools/image-test.sh` | update | jede Phase deklariert ihre Kennungen an Ort und Stelle |
| `.d-check.yml` (`trace.coverage`) | update | Einbindung |
| `harness/sensors/test.md` | update | dritte Zusage des Repos über sich selbst |

**Beim Beanspruchen zu entscheiden und hier einzutragen:**

- **Die Deklarations-Form im Go-Test** — heute nennen 57 Testdateien
  56 verschiedene Kennungen an beliebiger Stelle (Kommentar, Testname,
  Fehlermeldung). Gezählt wird nur eine **feste** Form (etwa die Kennung im
  Doc-Kommentar direkt über `func Test…`); Messung, wie viele Tests sie schon
  tragen, und Negativliste vor dem Code (Workflow-Skelett Schritt 19).
- **Der Ort der Dateien** — die Schwester-Repos legen sie unter `docs/user/`;
  hier ist Maintainer-Doku seit [`MR-077`](../../../../harness/conventions.md#mr-077)
  unter `docs/maintainer/`, dessen Geltungsbereich die Releasing-Doku nennt.
- **Die Tabelle bindet an `Datei:Zeile`** — jede Zeilenverschiebung über einem
  Test ändert die Datei; die Schwester-Repos nehmen das in Kauf. Abwägen gegen
  eine Bindung nur an den Testnamen.

## 4. Trigger

**Start** (`next` → `in-progress`): `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): die feste Deklarations-Form verlangt, dass
  mehr als eine Handvoll Tests umgeschrieben werden — dann zuerst die Form,
  die Spalte danach.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- **Die Spalte behauptet mehr, als der Test prüft** — eine Kennung im
  Kommentar ist eine Deklaration, kein Beleg; die RTM liest sich danach wie
  ein Nachweis. — **Ausgang:** *(offen)*

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

**Vorgelagert — Sub-Area-Wahl prüfen:** beim Beanspruchen.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** beim Beanspruchen.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
beim Beanspruchen.

**Modus-Begründungsblock:** beim Beanspruchen.
