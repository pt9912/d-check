# Slice slice-140: Konsumenten-CR an den Kurs — zwei Punkte, beide belegt

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos.**

**Bezug:** Baseline-Regelwerk
[`modul-05-planning-harness.md` §Offene Risiken werden bei Closure aufgelöst](../../../../.harness/baseline/v5.11.0/regelwerk/modul-05-planning-harness.md)
und
[`modul-06-roadmap.md` §Das Beobachtungs-Register](../../../../.harness/baseline/v5.11.0/regelwerk/modul-06-roadmap.md);
[`modul-03-spec.md` §Ziel-Form: Architektur-Sicht](../../../../.harness/baseline/v5.11.0/regelwerk/modul-03-spec.md)
und `AGENTS.template.md` §3.4;
[`MR-033`](../../../../harness/conventions.md#mr-033);
[`BEO-015`](../observations.md).

**Berührte Spec-Stellen:** — (Fremd-Repo; keine Anforderung dieses Repos).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

Zwei Beobachtungen aus dem eigenen Betrieb treffen **jedes** Adopter-Repo, nicht
nur dieses. Beide sind belegt, keine ist ein Vorschlag ins Blaue.

**Punkt 1 — die urteilsfreie Hälfte der Drei-Ausgänge-Regel ist nicht benannt.**
`modul-05` verlangt, dass jeder offene Punkt beim Übergang nach `done/` genau
einen von drei Ausgängen bekommt, und schließt mit *„Ein Slice geht nicht nach
`done/`, während ein Risiko ohne Ausgang dasteht."* Im **Register**-Abschnitt
benennt der Kanon für eine vergleichbare Regel ausdrücklich, was maschinell
entscheidbar ist — *„drei Prüfungen ohne Urteil: Form · Anzahl · Lage"*, und
*„welches Werkzeug, ist Repo-Entscheidung"*. Für die Drei-Ausgänge-Regel fehlt
diese Unterscheidung, obwohl sie eine hat: **ein Slice in `done/` trägt keinen
unaufgelösten Vorlagen-Platzhalter.** Der CR bittet nicht um ein Gate, sondern
um denselben Satz an einer zweiten Stelle.

**Punkt 2 — *„referenziert Modul-Pfade"* ist mehrdeutig.** `modul-03` und
`AGENTS.template.md` §3.4 erlauben der Architektur-Sicht, *Modul-Pfade* zu
referenzieren. Ob damit **Code**-Pfade oder **Dokument**-Pfade gemeint sind,
sagt keine der beiden Stellen — und die Antwort entscheidet, ob ein Repo, das
Code-Pfade in der Sicht verbietet, konform ist oder verschärft. Dieses Repo hat
sich entschieden und die Verschärfung deklariert
([`MR-033`](../../../../harness/conventions.md#mr-033)); jedes andere Repo steht
vor derselben ungeklärten Frage.

## 2. Vorgehen

1. **Reihenfolge:** Punkt 1 erst schreiben, wenn
   [slice-139](../in-progress/slice-139-closure-ausgang-waechter.md) in `done/` liegt. Dann
   trägt der CR eine **gebaute** Form und gemessene Zahlen statt eines
   Vorschlags — genau die Bedingung, an der ein früherer CR-Punkt dieses Repos
   zu Recht gescheitert ist (*„ohne Baubarkeit wäre sie ein behauptetes Gate"*).
2. **Je Punkt: Beleg vor Bitte.** Was steht wo, was fehlt, was folgt daraus für
   einen Adopter — und ausdrücklich, was der CR **nicht** verlangt.
3. **Getrennt halten.** Die zwei Punkte haben nichts miteinander zu tun; sie
   zusammenzubinden schwächt beide. Ein CR-Dokument, zwei nummerierte Punkte,
   jeder für sich annehmbar oder ablehnbar.
4. **Vorlegen, nicht senden.** Ob und wann der CR beim Kurs eingeht, ist eine
   Auftraggeber-Entscheidung.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Änderung am vendorten Baum.** Der CR ist eine Bitte an die Quelle;
  bis der Kurs entscheidet, gilt hier die Baseline unverändert.
- **Keine dritte Bitte.** Was sonst noch auffiel, wartet auf eigene Belege.
- **Keine Vorwegnahme der Antwort.** Lehnt der Kurs Punkt 2 ab und meint
  Code-Pfade, wird die dortige Adaption gegenstandslos — das ist ein Ergebnis,
  kein Verlust.

## 4. Definition of Done

- [ ] Beide Punkte stehen mit **Zitat und Fundstelle**, nicht mit Auslegung.
- [ ] Punkt 1 nennt die gebaute Form und ihre gemessenen Zahlen.
- [ ] Je Punkt steht da, was der CR **nicht** verlangt.
- [ ] Das Dokument liegt vor; die Entscheidung über das Absenden ist
      ausdrücklich offen gelassen.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein CR aus dem eigenen Schmerz heraus ist selten allgemein.** Beide Punkte
  müssen zeigen, dass sie **jedes** Adopter-Repo treffen — sonst sind sie eine
  repo-lokale Adaption, und dafür gibt es den Konventionsspeicher. —
  **Ausgang:** *(bei Closure)*
- **Punkt 2 könnte sich als Lesefehler erweisen.** Wenn der Kanon anderswo
  klärt, was *Modul-Pfad* heißt, ist der CR gegenstandslos und jene Adaption
  falsch. Vor dem Schreiben ist der ganze Kanon danach zu durchsuchen, nicht nur
  die zwei bekannten Stellen. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-139](../in-progress/slice-139-closure-ausgang-waechter.md)
in `done/` — für Punkt 1 ist die gebaute Form die Eintrittskarte.

**Rückführungen:** `in-progress` → `next`, falls die Kanon-Suche zu Punkt 2 die
Mehrdeutigkeit auflöst; dann bleibt nur Punkt 1, und das ist eine
Auftraggeber-Frage.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Regeltext (GF), Konventionsspeicher (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23):
  [`BEO-012`](../observations.md) ist **zentral** — ein CR behauptet, was eine
  fremde Quelle sagt und nicht sagt; genau daran ist diese Arbeit mehrfach
  gescheitert. [`BEO-011`](../observations.md) für jede Aussage darüber, dass
  ein Punkt „jedes Repo" trifft.

Slice-ID: slice-140. Betroffene IDs:
[`MR-033`](../../../../harness/conventions.md#mr-033). Module:
Harness-Regeltext. Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Beleg-Arbeit an fremdem Regeltext.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
