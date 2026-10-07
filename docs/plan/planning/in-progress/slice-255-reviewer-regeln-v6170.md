# slice-255: Reviewer-Regeln aus `v6.17.0` in die Reviewer-Skills

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [`MR-074`](../../../../harness/conventions.md#mr-074) (Bewegung 3:
Reviewer-Regeln), Baseline `v6.17.0` · `regelwerk/modul-10-review-harness.md`
§Ziel-Form: Reviewer-Skill und die Vorlagen
`templates/.harness/skills/reviewer.template.md` /
`closure-note-reviewer.template.md`; Auftraggeber-Freigabe 2026-10-07.

**Berührte Spec-Stellen:** — *(Harness-Skills, keine Spec-Aussage)*.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Die Reviewer-Skills folgen den drei Reviewer-Regeln der Baseline
`v6.17.0`: **LOW nur mit Konventions-Anker** (ADR, Hard Rule, Linter-Regel,
Eintrag im Skill); die **Failure-Szenario-Pflicht gilt HIGH und MEDIUM**; das
Feld **`pfad`** ankert die Fundstelle als wörtliches, in der Datei eindeutig
auffindbares Kurzzitat, die Zeile ist Lesehilfe. Die vierte Regel der
Bewegung, „Kein Stil-Polizist", trug der Skill bereits. *(Plan-Änderung nach
R1: die Kontext-Eskalation, das `quelle`-Feld und die zweite
`pfad`-Definition unter §Ablage werden mitgezogen, weil die neuen Regeln
sonst mit ihnen kollidieren; dazu die veraltete Zahl der Prüffragen —
achtzehn statt sechzehn — in Skill und Agent-Spiegel.)*

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die übrigen Bewegungen aus dem Hebung-Eintrag** (Spezifikation,
  Werkzeug-Teil) — eigene Slices.
- **Nachträgliche Umformung bestehender Review-Reports** — sie sind Lauf-Belege
  und frozen; die Regel gilt ab dem nächsten Report.
- **Ein Sensor auf die Report-Form** — `pfad` als Kurzzitat ist ein Urteil über
  Eindeutigkeit, kein Formmerkmal, das ein Gate ohne Fehlalarme prüft.

**Zur Lockerung der Failure-Szenario-Pflicht:** der Skill verbot bisher *jedes*
Finding ohne Failure-Szenario, strenger als die Baseline. Die Baseline lässt
LOW (mit Konventions-Anker) und INFO ohne erzähltes Versagen zu — ein
Konventionsverstoß ist ein Befund, auch wenn er heute nichts bricht, und INFO
trägt Beobachtungen außerhalb des Auftrags. Ob die strengere Form in der
Praxis etwas verhindert hat, was die Baseline-Form durchließe, ist nicht
gemessen; der Tausch folgt der Baseline-Regel. Kein Gate ist betroffen
(`AGENTS.md` §3.6 nicht berührt).

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] `.harness/skills/reviewer.md`: LOW-Anker, Failure-Szenario nur für
      HIGH/MEDIUM, `pfad` als Kurzzitat; Version und Datum gehoben; jede
      übernommene Regel wortnah zur Baseline-Quelle, mit `d-check:cite`, wo
      wörtlich zitiert wird.
- [ ] `.harness/skills/closure-note-reviewer.md`: `pfad` als Kurzzitat;
      Version gehoben.
- [ ] Spiegel geprüft (`.claude/agents/reviewer.md`, Agent-Prompts in
      `.claude/`); `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor —
      Rollenwechsel nach Schritt 8, kein Self-Review.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft;
      [MR-074](../../../../harness/conventions.md#mr-074) Bewegung 3 als eingelöst vermerkt.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.harness/skills/reviewer.md` | update | drei Regeln, Version |
| `.harness/skills/closure-note-reviewer.md` | update | `pfad`, Version |
| Hebung-Eintrag (MR-Datei) | update | Einlösung vermerkt (Closure) |

**Spiegel vor dem Editieren** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
`.claude/agents/reviewer.md` (Output-Beschreibung), die Agent-Prompts, die
`pfad` als `Datei:Zeile` beschreiben, `harness/README.md` §Guides (nennt die
Skills nur), die vorhandenen Cite-Direktiven im Reviewer-Skill.

## 4. Trigger

**Start** (`next` → `in-progress`): Implementer übernimmt; `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): die Regeln verlangen eine Umstellung des
  Prüffragen-Katalogs über die drei Stellen hinaus.
- `in-progress` → `open` (blockiert): keiner bekannt.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- Kurzzitate als Anker sind bei sehr kurzen oder sich wiederholenden Zeilen
  nicht eindeutig. — **Ausgang:** *(offen)*

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

**Vorgelagert — Sub-Area-Wahl prüfen:** eine berührte Sub-Area: der Harness
des Repos (`*`, `ALL`); deklariert.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:**
[`BEO-ALL/citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/state.md)
— wörtliche Übernahmen aus der Baseline werden als `d-check:cite` geführt,
damit ein Zitat nicht mehr trägt, als es sagt.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-07 — beide Läufe grün.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
