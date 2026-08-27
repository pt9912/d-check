# Slice slice-160: Der Reviewer-Skill trägt sechzehn Anker

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`.harness/skills/reviewer.md`](../../../../.harness/skills/reviewer.md);
Baseline
[`modul-10-review-harness.md`](../../../../.harness/baseline/v5.12.0/regelwerk/modul-10-review-harness.md)
§Ziel-Form: Reviewer-Skill;
[slice-131](../done/slice-131-reviewer-skill-waisen.md) (der Waisen-Zensus);
[slice-147](../done/slice-147-reviewer-anker-reichweite.md) (der Anlass).

**Berührte Spec-Stellen:** — (Rollen-Skill; das Produkt bleibt unberührt).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-27.

---

## 1. Ziel

Der Reviewer-Skill trägt **16** benannte HIGH-/MEDIUM-Anker — 7 HIGH, 2 einzeln
stehende MEDIUM, 7 MEDIUM in einem Sammelpunkt. Jeder einzelne ist gegen einen
belegten Fall geschrieben; zusammen sind sie eine Leseliste, die vor jedem
Review vollständig im Kopf sein müsste.

**Die Sorge ist nicht neu, sie ist nur nie gemessen worden.** Mehrere Slices
tragen den Risikosatz *„der fünfzehnte wird nicht mehr gelesen"* als
Abnahme-Punkt und beantworten dann die Nebenfrage („schärft dieser Anker einen
bestehenden?") statt der gestellten. Die Schwelle war beim sechzehnten Anker
längst überschritten.

**Die erste Frage ist eine Messung, keine Umstrukturierung:** greifen die Anker
überhaupt? Über die abgelegten Reports lässt sich zählen, welcher Anker wie oft
zu einem Befund geführt hat — und welcher seit seiner Einführung **nie**.

## 2. Vorgehen

1. **Zählen, bevor umgebaut wird.** Je Anker: wie viele Befunde in
   `docs/reviews/` lassen sich ihm zuordnen, und seit wann steht er? Ein Anker
   ohne Treffer ist nicht automatisch wertlos — er kann die Klasse verhindert
   haben —, aber die Zahl gehört auf den Tisch, bevor jemand kürzt.
2. **Die Kategorien-Semantik der Baseline gegenprüfen.** `modul-10` gibt
   HIGH/MEDIUM/LOW/INFO vor; ob 16 repo-eigene Anker darunter die Ziel-Form
   noch treffen oder sie überschreiben, ist zu lesen, nicht anzunehmen.
3. **Erst dann die Form entscheiden.** Kandidaten: Zusammenzug verwandter Anker;
   eine Zwei-Ebenen-Form (kurze Prüffragen-Liste vorn, Begründungen hinten);
   Auslagerung der Begründungen in die Registerzeilen, auf die sie ohnehin
   zeigen. Jede Form hat einen Preis — der gehört benannt.
4. **Kein Anker wird gestrichen, ohne dass seine Klasse einen anderen Ort hat.**
   Sonst entsteht die Waise in der Gegenrichtung: die Regel steht in
   `AGENTS.md`, und niemand prüft sie mehr.
5. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine neue Anker-Klasse.** Der Slice räumt, er ergänzt nicht.
- **Keine Änderung an den Hard Rules selbst.** Was in `AGENTS.md` §3/§5 steht,
  bleibt; hier geht es um die Feedback-Hälfte.
- **Kein Zusammenzug ohne Messung.** „Wirkt verwandt" ist kein Grund; zwei
  Anker mit unvereinbaren Arbeitsanweisungen gehören auseinander, auch wenn sie
  dieselbe Kategorie tragen.

## 4. Definition of Done

- [ ] Je Anker eine **Trefferzahl** aus `docs/reviews/`, mit dem Zeitraum seit
      seiner Einführung.
- [ ] Die Form ist entschieden und ihr Preis benannt.
- [ ] Kein gestrichener Anker ohne benannten Ersatz-Ort für seine Klasse.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Die Trefferzahl misst die Vergangenheit, nicht den Wert.** Ein Anker ohne
  Treffer kann seine Klasse verhindert haben; wer nach Zahlen kürzt, entfernt
  womöglich genau den, der wirkt. Die Zahl ist ein Eingang in die Entscheidung,
  nicht die Entscheidung. — **Ausgang:** *(bei Closure)*
- **Kürzen ist die Bewegung, die sich gut anfühlt.** Ein kürzerer Skill liest
  sich besser und prüft schlechter; der Slice darf nicht am Umfang gemessen
  werden, sondern an der Frage, ob eine Klasse ihren Ort behält. —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Messung ergibt, dass die
Anker greifen und die Last tragbar ist — dann ist die ehrliche Lieferung eine
gestrichene Sorge, keine Umstrukturierung.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Rollen-Skills (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-27):
  [`BEO-011`](../observations.md) — die Form gehört aus dem Bestand, nicht aus
  dem einen Anker, der zuletzt dazukam;
  [`BEO-012`](../observations.md) — vor jedem Zitat der Baseline deren
  Geltungsbereich lesen.
- **Nachtlauf-Stand lesen** (`make nightly-state`, dritte Vorprüfung nach
  [`MR-053`](../../../../harness/conventions.md#mr-053)): **ROT**, jüngster
  Lauf `2026-08-27T10:49:23Z`, `head_sha 48cf132`. **Gelesen:** derselbe Lauf,
  den [slice-164](../done/slice-164-nachtlauf-kadenz.md) §7 bereits eingeordnet
  hat — er lief vor den sechs Pin-Hebungen aus
  [slice-161](../done/slice-161-sechs-pins-heben.md), die zum Dispatch noch
  nicht auf `origin/main` lagen. Keine neue Meldung; der nächste Lauf ist die
  Probe darauf, und die steht noch aus.

Slice-ID: slice-160. Betroffene IDs: — (kein `DC-`-Bezug; Rollen-Skill).
Module: — . Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Aufräumen an einem vorhandenen Artefakt der
Review-Schicht.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
