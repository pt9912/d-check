# Slice slice-140: Konsumenten-CR an den Kurs — vier Punkte, jeder belegt

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

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-25.

---

## 1. Ziel

Vier Beobachtungen aus dem eigenen Betrieb treffen **jedes** Adopter-Repo, nicht
nur dieses. Jede ist belegt, keine ist ein Vorschlag ins Blaue.

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

**Punkt 3 — „bewusst gebrochen" sagt nicht, dass das Rot von der Regel kommt.**
`modul-13` beschreibt als sechsten Schritt der Fitness-Function-Übersetzung das
**Bewusste Brechen**: *„das Gate läuft rot mit `ADR-<NNNN> violated`. Genau der
Effekt, der eine ADR von einer Absichtserklärung trennt."* Was fehlt, ist ein
halber Satz: **das Rot muss von der gebrochenen Regel kommen, und seine Ursache
gehört gelesen.** Ohne ihn ist die Bedingung erfüllt, sobald *irgendetwas* rot
wird. Drei Instanzen aus zwei aufeinanderfolgenden Arbeitstagen, alle in einem
Repo, das diese Disziplin ernst nimmt: eine Suppression-Probe, die den falschen Linter
nannte (rot — aber vom Nachbarlinter); eine Direktiv-Form, die als *wohlgeformt*
galt und nie gefahren worden war (sie wird gemeldet); und eine Probe, die per
Anhängen im ausgenommenen Abschnitt landete (still — beinahe als *„der Wächter
greift nicht"* verbucht). Der Punkt verlangt **kein** Gate; er schärft einen
Satz, den der Kanon schon hat.

**Punkt 4 — die Reichweite eines Zitats ist nur im Einzelfall geregelt.** Der
Kanon ordnet Quellen (`grundlagen-source-precedence.md` §Source Precedence) und
beantwortet die Reichweitenfrage dort zweimal — für Adaptions-Einträge und für
seinen eigenen Rangordnungs-Satz. Was fehlt, ist der Schritt von diesen
Einzelfällen zur Frage als solcher: **wie weit trägt eine einzelne Aussage aus
einer gerankten Quelle?** Genau dort entsteht eine eigene Fehlerklasse: der Text stimmt, die in
Anspruch genommene Reichweite nicht — ein Satz, der nur für einen benannten Fall
gilt, wird universal geführt; ein Adaptions-Eintrag nach seinem **Titel**
zitiert statt nach seinem Feld `Geltungsbereich`; eine ADR-Entscheidung, die
einen einmaligen Akt beschreibt, als stehendes Verbot gelesen. Dieses Repo führt
sie als [`BEO-012`](../observations.md) mit erreichter Schwelle. Weil ein Zitat
wie ein Beleg **aussieht**, ist es schwerer zu bemerken als eine unbelegte
Behauptung — und weil der Kanon selbst die zitierte Quelle ist, trifft es jeden
Adopter.

## 2. Vorgehen

0. **Nur Punkt 1 hat eine Vorbedingung.** Die Punkte 2 bis 4 hängen an nichts;
   sie können mit dem Dokument entstehen.
1. **Reihenfolge:** Punkt 1 erst schreiben, wenn
   [slice-139](../done/slice-139-closure-ausgang-waechter.md) in `done/` liegt. Dann
   trägt der CR eine **gebaute** Form und gemessene Zahlen statt eines
   Vorschlags — genau die Bedingung, an der ein früherer CR-Punkt dieses Repos
   zu Recht gescheitert ist (*„ohne Baubarkeit wäre sie ein behauptetes Gate"*).
2. **Je Punkt: Beleg vor Bitte.** Was steht wo, was fehlt, was folgt daraus für
   einen Adopter — und ausdrücklich, was der CR **nicht** verlangt.
3. **Getrennt halten.** Die vier Punkte haben nichts miteinander zu tun; sie
   zusammenzubinden schwächt alle. Ein CR-Dokument, vier nummerierte Punkte,
   jeder für sich annehmbar oder ablehnbar.
4. **Vorlegen, nicht senden.** Ob und wann der CR beim Kurs eingeht, ist eine
   Auftraggeber-Entscheidung.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Änderung am vendorten Baum.** Der CR ist eine Bitte an die Quelle;
  bis der Kurs entscheidet, gilt hier die Baseline unverändert.
- **Keine fünfte Bitte.** Die Grenze stand zuerst bei zwei Punkten, mit der
  Bedingung *„was sonst noch auffiel, wartet auf eigene Belege"*. Für zwei
  weitere liegen sie inzwischen vor — je drei gemessene Instanzen —, und der
  Auftraggeber hat die Erweiterung entschieden. Sie ist damit verschoben, nicht
  gefallen: was **keine** drei Belege hat, wartet weiter. Ausdrücklich draußen
  bleiben die repo-eigenen Funde — das Fixture mit gleicher Datei-Größe und die
  Eigenheiten unserer Messskripte.
- **Keine Vorwegnahme der Antwort.** Lehnt der Kurs Punkt 2 ab und meint
  Code-Pfade, wird die dortige Adaption gegenstandslos — das ist ein Ergebnis,
  kein Verlust.

## 4. Definition of Done

- [ ] **Alle vier** Punkte stehen mit **Zitat und Fundstelle**, nicht mit
      Auslegung.
- [ ] Punkt 1 nennt die gebaute Form und ihre gemessenen Zahlen; Punkt 3 und 4
      nennen je **drei** Instanzen mit Fundstelle.
- [ ] Je Punkt steht da, was der CR **nicht** verlangt — bei Punkt 3 und 4
      ausdrücklich: **kein Gate**, nur ein geschärfter Satz.
- [ ] Das Dokument liegt vor; die Entscheidung über das Absenden ist
      ausdrücklich offen gelassen.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein CR aus dem eigenen Schmerz heraus ist selten allgemein.** Beide Punkte
  müssen zeigen, dass sie **jedes** Adopter-Repo treffen — sonst sind sie eine
  repo-lokale Adaption, und dafür gibt es den Konventionsspeicher. —
  **Ausgang:** *(bei Closure)*
- **Vier Punkte in einem CR laden zur Teil-Annahme ein.** Beim letzten CR
  dieses Repos wurden zwei von mehreren Punkten abgelehnt — das ist gesund,
  solange jeder Punkt für sich steht. Kippt aber die Nummerierung zur
  Wunschliste, sinkt die Chance für alle. Zu prüfen ist deshalb je Punkt, ob er
  **ohne** die anderen trägt. — **Ausgang:** *(bei Closure)*
- **Punkt 2 könnte sich als Lesefehler erweisen.** Wenn der Kanon anderswo
  klärt, was *Modul-Pfad* heißt, ist der CR gegenstandslos und jene Adaption
  falsch. Vor dem Schreiben ist der ganze Kanon danach zu durchsuchen, nicht nur
  die zwei bekannten Stellen. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-139](../done/slice-139-closure-ausgang-waechter.md)
in `done/` — für Punkt 1 ist die gebaute Form die Eintrittskarte.

**Rückführungen:** `in-progress` → `next`, falls die Kanon-Suche zu Punkt 2 die
Mehrdeutigkeit auflöst; dann bleibt nur Punkt 1, und das ist eine
Auftraggeber-Frage.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Regeltext (GF), Konventionsspeicher (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-25):
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
