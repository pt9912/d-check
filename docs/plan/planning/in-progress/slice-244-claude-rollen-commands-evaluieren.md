# slice-244: `.claude`-Rollen und Commands aus ai-harness-init evaluieren

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle as State Machine.

**Welle:** welle-91.

**Bezug:** [`welle-91`](../welle-91-adoption-ai-harness-init.md),
[`MR-073`](../../../../harness/conventions.md#mr-073)
(Baseline-Stand, der dem Schwester-Bootstrap zugrunde liegt).
Keine `DC-*` — Rollen- und Command-Doku berührt kein
Produkt-Anforderungs-Delta.

**Berührte Spec-Stellen:** — *(Agent-/Command-Doku ist kein Spec-Stratum)*.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-29.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die sechs `.claude/agents/`-Rollen der Schwester (architect,
implementer, planner, reviewer, validator, verifier) und die drei
`.claude/commands/` (plan-welle, implement-slice, close-welle) sind je
evaluiert — adoptiert (an d-check angepasst: d-check führt reviewer und
verifier bereits als Agents, implement-slice als Command unter
`.claude/commands/`) oder abgelehnt (je mit Begründung gegen Modul 8 und
d-checks Betrieb). Die **Implementer-Rolle ist abgelehnt**: ihre Arbeit läuft
im Hauptlauf nach AGENTS.md §6 bzw. über den implement-slice-Command — ein
Agent-Duplikat trüge dieselbe Anweisung an zweiter Stelle.

**Ausdrücklich NICHT in diesem Slice:**

- Die **Erfassungsschicht** (`span-emit.sh`, `erfassung.mk`) — Out-of-Scope
  von welle-91 (Welle §6).
- Umbauten an den bestehenden Rollen reviewer/verifier — nur Neu-Einführung
  oder Nein.
- Kein Verpflichten auf Rollen-Besetzung: adoptierte Rollen-Doku beschreibt
  die Rolle; wer sie füllt, entscheidet sich pro Vorgang.

## 2. Definition of Done

- [ ] Je Rolle (6) und Command (3) eine belegte Entscheidung:
      adoptiert (Datei in `.claude/`, an d-check angepasst) oder abgelehnt
      (Begründung, gegen Modul 8 und den d-check-Betrieb).
- [ ] Adoptierte Rollen-Dateien tragen die Kontext-Trennung der Vorlage
      („Was du NICHT bist") und verweisen auf die d-check-Pendants.
- [ ] `make gates` grün.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.claude/agents/*.md` (nach Evaluierung) | neu | adoptierte Rollen |
| `.claude/commands/*.md` (nach Evaluierung) | neu | adoptierte Commands |
| dieser Plan | update | etwaige Anker-Korrekturen |

## 4. Trigger

**Start** (`open` → `in-progress`): direkt beansprucht — die
Welle-Eröffnung (welle-91, Auftraggeber 2026-09-29) ist die Beanspruchung;
der Nachtlauf-Stand wird bei der Beanspruchung gelesen
([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Rückführungen:** `in-progress` → `next` (zu groß): fällt die Evaluierung
je Rolle verschieden scharf aus (Adoption vs. Diskussion), wird pro Rolle
geschnitten. `in-progress` → `open`: keiner bekannt.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag (drei Paarungen
wellenlos hier geprüft).

## 6. Risiken und offene Punkte

- Adoptierte Rollen bleiben unbesetzt (Prosa ohne Betrieb) — der
  Rollen-Zweck wird je dokumentiert; Betrieb-Erfahrungen sammelt die
  Welle-Closure. — **Ausgang:** *(offen)*
- planner/plan-welle/close-welle setzen Wellen-Betrieb voraus, der noch
  nicht etabliert ist — die Commands stützen den künftigen Betrieb
  (Auftraggeber: „Wir werden auch in Zukunft Wellen erstellen"). —
  **Ausgang:** *(offen)*

## 7. Closure-Notiz

*(gefüllt vor dem `git mv` nach `done/`)*

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —
- **Nachtlauf-Stand bei der Beanspruchung** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
  beide Nachtläufe grün (upstream-drift 2026-09-29, image-scan 2026-09-28).

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: das
Harness-Werkzeug selbst (`*`, Kürzel `ALL`); bereits deklariert.

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** keine berührte Sub-Area
mit Treffern in offenen Beobachtungen.

**Vorgelagert — Nachtlauf-Stand lesen** (bei der Beanspruchung):
`make nightly-state` wird bei der Beanspruchung gelesen und der Stand in
§7 notiert.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.

### Sub-Area: `*` (Harness-Werkzeug)

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — die Vorlage ist die der Schwester, die
  Kontext-Trennungen sind Modul-8-fest.
- **Phase-Reife:** Phase 3.
- **Evidenz-/Diskrepanz-Risiko:** Niedrig — Doku-Adoption mit je Rolle
  begründeter Entscheidung.
- **Reconciliation-Aufwand:** Keiner.
