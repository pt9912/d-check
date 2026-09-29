# slice-241: Adoption „Trigger-Audit der Welle" (modul-06, v6.13.0) in den wellenlosen Closure-Vollzug

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Adoption einer Closure-Praxis ist kein
repo-weites Mehr (Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).

**Bezug:** die Hebung
[`MR-073`](../../../../harness/conventions.md#mr-073--baseline-pin-hebung-auf-v6130-fünfzehnter-nachtrag-zu-mr-011-nachtrag-zu-mr-023)
liefert das Delta (Baseline-Regelwerk `modul-06-roadmap.md` §Closure,
Schritt 2 „Trigger-Audit der Welle"), [`MR-051`](../../../../harness/conventions.md#mr-051)
(Re-Ankern), [`MR-054`](../../../../harness/conventions.md#mr-054),
[`MR-053`](../../../../harness/conventions.md#mr-053). Keine `DC-*` —
keine Produkt-Anforderung berührt.

**Berührte Spec-Stellen:** — *(die Closure-Praxis ist kein Spec-Stratum;
`spec/` bleibt unverändert)*.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-29.

---

## 1. Ziel und Abgrenzung

**Ziel:** Den neuen Welle-Closure-Schritt „Trigger-Audit der Welle"
(`modul-06-roadmap.md` v6.13.0, §Closure Schritt 2: vier Artefaktklassen
tragen einen Trigger — **Carveout** · **bootstrap-aware Gate** · **ADR** ·
**Hard Rule** — „ein Trigger ohne Wächter ist eine Absichtserklärung mit
Verfallsdatum") in den wellenlosen Betrieb übernehmen: d-check schließt
wellenlos, dort löst die Slice-Closure die Register-Lese-Schritte selbst
aus — der Audit wird als fester Schritt in diese Closure-Praxis
aufgenommen, dort dokumentiert, wo die Closure-Schritte beschrieben sind,
und bei der Closure dieses eigenen Slices erstmals vollzogen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Keine Mechanisierung** (Sensor/Config-Schlüssel für den Audit) —
  „verkörpert heißt nicht zwangsläufig automatisiert" (modul-06 v6.13.0
  ebenda); ein Sensor wäre eine Produkt-Entscheidung mit eigenem
  Anforderungs-Delta, nicht eine Konventions-Adoption.
- **Keine Retro-Audit über den Bestand** — die `done/`-Slices gelten als
  geprüft durch ihre Reviews; der Audit gilt ab seiner Einführung für
  neue Closures.
- **Keine Änderung der Closure-Notiz-Schema-Form** ([`MR-056`](../../../../harness/conventions.md#mr-056)
  bleibt, wie sie ist) — der Audit ist ein Prüfschritt, kein neues
  Pflichtfeld, solange der Vollzug nichts anderes zwingend macht.

## 2. Definition of Done

- [ ] Der Trigger-Audit ist als Closure-Schritt dokumentiert — an dem Ort,
      an dem die wellenlosen Closure-Schritte beschrieben sind
      (Entscheidung README-Abschnitt vs. Reviewer-Skill im Vollzug, mit
      Begründung) — und benennt die vier Klassen (Carveout ·
      bootstrap-aware Gate · ADR · Hard Rule) mit je Trigger und Wächter.
- [ ] Der erste Audit-Vollzug ist belegt: an der eigenen Closure dieses
      Slices sind die vier Klassen je geprüft (Beleg in der
      Closure-Notiz).
- [ ] `make gates` grün.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/planning/observations/README.md` **oder** `.harness/skills/reviewer.md` | update | der Ort des Closure-Schritts — Entscheidung im Vollzug je Gegenstand: README beschreibt den Ablauf, der Skill prüft den Report |
| dieser Plan | update | etwaige Anker-Korrekturen, falls der Vollzug Zeilen verschiebt |

## 4. Trigger

**Start** (`open` → `in-progress`): direkt beansprucht — die
Auftraggeber-Freigabe „ja bitte" zu den vier Deltas der v6.13.0-Hebung
(2026-09-29) ist die Beanspruchung selbst; der Nachtlauf-Stand wird bei
der Beanspruchung gelesen ([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): die Ort-Frage (README vs. Skill)
  zwingt zu zwei Dokumentations-Hälften, die sich nicht je Fundstelle
  trennen lassen — dann ist der Nachzug ein eigener, feinerer Slice.
- `in-progress` → `open` (Blocker): keiner bekannt; der Delta-Wortlaut ist
  vendiert und gelesen.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag (drei Paarungen
wellenlos hier geprüft).

## 6. Risiken und offene Punkte

- Eine Doppel-Dokumentation (README **und** Skill) beschreiben denselben
  Schritt verschieden — eine Quelle gewinnt, die andere bleibt Zeiger. —
  **Ausgang:** *(offen)*
- Der Audit bleibt Prosa ohne Vollzugs-Beleg — DoD 2 deckt die erste
  Ausführung ab. — **Ausgang:** *(offen)*
- Die Vier-Klassen-Terminologie (bootstrap-aware Gate) passt nicht 1:1 auf
  d-checks Gate-Landschaft (kein Stufen-Gate) — die Übertragung wird im
  Vollzug je Klasse belegt oder als n.a. begründet. — **Ausgang:** *(offen)*

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

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: das
Harness-Werkzeug selbst (`*`, Kürzel `ALL`); bereits deklariert.

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:**
[`BEO-ALL/registerzeile-ohne-ausgang-nach-schwelle`](../observations/BEO-ALL/registerzeile-ohne-ausgang-nach-schwelle/observation.md)
ist benachbart offen (Trigger mit Verfallsdatum ohne Ausgang — derselbe
Gegenstandsbereich Register/Steering-Loop); keine weitere berührte
Sub-Area mit Treffern; der Ausgang bleibt ausgeschildert.

**Vorgelagert — Nachtlauf-Stand lesen** (bei der Beanspruchung):
`make nightly-state` wird bei der Beanspruchung gelesen und der Stand in
§7 notiert.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.

### Sub-Area: `*` (Harness-Werkzeug)

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — die Prozedur ist verkörpert (MR-Kette,
  Reviewer-Skill, Closure-Profile), das Delta ist vendiert und gelesen.
- **Phase-Reife:** Phase 3.
- **Evidenz-/Diskrepanz-Risiko:** Niedrig — Doku-Adoption mit Beleg-Vollzug.
- **Reconciliation-Aufwand:** Keiner — Graduation-Trigger bleibt
  [`BEO-ALL/registerzeile-ohne-ausgang-nach-schwelle`](../observations/BEO-ALL/registerzeile-ohne-ausgang-nach-schwelle/observation.md).
