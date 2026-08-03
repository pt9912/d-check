# Slice slice-090: Regelwerk-Migration Etappe D — Review-Infrastruktur (reviewer.md-Currency · Report-Kopffelder)

**Status:** In Arbeit (welle-67).

**Welle:** welle-67-baseline-v500-migration (Etappe D, dritter „Mini-Welle"-Slice, nach
[slice-089](../done/slice-089-doc-form-agents.md)).

**Bezug:** Etappe D (Form-Konformität) aus
[slice-085](../done/slice-085-etappe-b-modul-delta.md) §3.2 — **D-10** (reviewer.md-Kopf-
Currency), **D-6** (Review-Report-Kopffelder + Finding-`klasse`) und **D-7**
(`closure-note-reviewer`-Skill + `verify-closure-notes`-Gate). Mini-Welle: slice-088
Planning-Layer (done) → slice-089 Doc-Form (done) → **slice-090 Review-Infrastruktur** →
slice-091 Slice-`Status:`-Feld (D-5). **Kein Release** (D-6/D-10); D-7 ist gesondert
(siehe §3).

**Autor:** pt9912. **Datum:** 2026-08-02.

---

## 1. Ziel

Die d-check-Review-Infrastruktur an die v5.0.0-Baseline-Form angleichen: den
`.harness/skills/reviewer.md`-**Kopf** aktualisieren (D-10) und das **Output-/Report-
Schema** um die Baseline-Kopffelder (Review-Art · Skill-Version · Modell-ID) und das
Finding-Feld `klasse` ergänzen (D-6), abgestimmt auf die vendored
`review-report.template.md`. **D-7** (der dedizierte Closure-Note-Reviewer **samt Gate**)
ist eine From-Scratch-**Code**-Adoption → Abnahme-Punkt (§3).

## 2. Vorgehen

1. **D-10 reviewer.md-Kopf-Currency.** Der Kopf zitiert die **retirete**
   `grundlagen-konventionen.md` (→ `grundlagen-referenz-richtung.md`) und „Agents-Digest
   Kurs-Welle 18 §8" (→ v5.0.0-`modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill).
   Version bumpen (1.2.0 → 1.3.0), Datum aktualisieren.
2. **D-6 Output-/Report-Schema.** In `reviewer.md` §Output-Schema das Feld **`klasse`**
   (stabile Fehlermuster-Bezeichnung) ergänzen; §Ablage um die **Kopf-Metadaten**
   präzisieren, die die vendored `review-report.template.md` führt: **Review-Art**
   (Plan/Design/Code), **Skill** (`@ Version/Commit`), **Modell-ID**, Gegenstand,
   Eingangs-Kontext. Zeiger auf die Template-Form, nicht Duplikat.
3. **D-7 (Abnahme-Punkt §3).** Der `closure-note-reviewer`-Skill (semantische Schicht)
   **über** dem `verify-closure-notes`-Gate (Struktur-Prüfung, Python + make-Target)
   fehlt in d-check komplett; die Closure-Note-Pflicht-ADR des Kurs-Templates hat in
   d-check **keine Entsprechung** (die 0011er-Nummer trägt dort die Digest-Pins-
   Entscheidung). Adoption = neues Skript + Gate + Skill + **eigene ADR** → **nicht** in
   diesem Doc-Form-Slice.
4. **Gate.** `make gates` + `make adr-check` grün; unabhängiger Frischkontext-Review.

## 3. Abnahme-Punkte

1. **D-7 Closure-Note-Reviewer + Gate.** Der vendored Skill ist die *inferentielle*
   Schicht über `make verify-closure-notes` (`check_closure_notes.py`: Heading +
   Satzzahl + Floskel-Liste). In d-check fehlen **alle drei** Teile; es wäre eine
   From-Scratch-Adoption (Go-/Python-Gate + `Makefile`-Target + Skill + ADR, plus die
   Frage, ob es als `targets`/`planning`-Modul-Fähigkeit mechanisiert statt als Skript
   lebt). **Entscheid:** (a) in einen **Folge-Produkt-Slice** herausschneiden (empfohlen,
   analog C-3 — Code+ADR+ggf. Release gehört nicht in die Doc-Form-Mini-Welle); (b) jetzt
   in slice-090 mit-adoptieren (größer, verlässt den Doc-Form-Rahmen); (c) als deklarierte
   Adaption führen (d-check reviewt Closure-Notes über den allgemeinen `reviewer.md` bei
   jeder Slice-Closure — kein dediziertes Gate).

## 4. Definition of Done

- [ ] `reviewer.md`-Kopf v5.0.0-konform (D-10: `grundlagen-referenz-richtung`, `modul-10`;
  Version/Datum gebumpt).
- [ ] `reviewer.md` §Output-Schema trägt `klasse`; §Ablage nennt die Kopf-Metadaten
  (Review-Art/Skill-Version/Modell-ID) abgestimmt auf `review-report.template.md` (D-6).
- [ ] D-7-Entscheid (Abnahme-Punkt 1) festgehalten.
- [ ] `make gates` + `make adr-check` grün; unabhängiger Frischkontext-Review.

## 5. Risiken / offene Punkte

- **reviewer.md ist der Skill hinter jedem Review** — Schema-Änderungen betreffen alle
  Folge-Reviews; abwärtskompatibel halten (Felder ergänzen, nicht umbenennen).
- **D-7-Scope-Disziplin:** die Gate-Adoption nicht schleichend in diesen Slice ziehen.

## 6. Trigger

Abschluss von [slice-089](../done/slice-089-doc-form-agents.md) (Doc-Form).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** berührt *Harness/Prozess*-Doku (`.harness/skills/reviewer.md`) — GF.
- **Offene Beobachtungen sichten:** `observations.md` = `— keine —`; nichts zu
  berücksichtigen.

## 8. Sub-Area-Modus-Begründung

GF (Repo-Default): Doc/Prozess führt. Berührt den Reviewer-Skill (Harness/Prozess);
greenfield-Currency-/Form-Angleich an die adoptierte Baseline.

## 9. Closure-Notiz (nach `done/`)

_Ausstehend._
