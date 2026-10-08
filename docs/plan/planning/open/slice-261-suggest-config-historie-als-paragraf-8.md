# slice-261: `--suggest-config` kennt die Spezifikations-Historie als §8

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** Befund F-6 aus dem Review von slice-259; Baseline `v6.17.0` ·
`templates/spec/spezifikation.template.md` (Historie als §8).

**Berührte Spec-Stellen:** `spec/spezifikation.md` §2 (Beispiel der
vorgeschlagenen `matrix`-Konfiguration).

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-08.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Die von `--suggest-config` vorgeschlagene
`matrix.exclude-sections`-Liste nimmt neben `7. Historie` auch
`8. Historie` auf — die Überschrift, unter der eine Spezifikation nach der
Vorlage seit Baseline `v6.17.0` ihre Historie führt. Ein Adopter, der die
neue Vorlage nutzt, bekäme sonst eine Ausnahme vorgeschlagen, die seine
Spezifikations-Historie nicht trifft.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Andere Vorschläge von `--suggest-config`** — nur diese Überschrift hat
  sich mit der Baseline-Vorlage verschoben.
- **Die Konfiguration dieses Repos** — slice-259 hat sie bereits
  nachgezogen.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] `--suggest-config` schlägt `8. Historie` mit vor; Test, der ohne den
      Fix aus dem richtigen Grund rot ist; Beispiel in der Spezifikation §2
      nachgezogen.
- [ ] Ohne die Option ist die Ausgabe unverändert (Black-Box-Probe gegen den
      Vorher-Stand); `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Vorschlags-Erzeugung von `--suggest-config` | update | `8. Historie` ergänzen |
| `spec/spezifikation.md` §2 | update | Beispiel |

## 4. Trigger

**Start** (`next` → `in-progress`): Implementer übernimmt; `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `open` (blockiert): keiner bekannt.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft. Die Änderung ist Produkt-Verhalten und geht mit dem
nächsten Release hinaus.

## 6. Risiken und offene Punkte

- Keine bekannt.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** eine berührte Sub-Area: das Produkt
unter dem Default `*` (`ALL`); deklariert.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** beim Beanspruchen neu lesen.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
beim Beanspruchen aus dem jüngsten Lauf lesen.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
