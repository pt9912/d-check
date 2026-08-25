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

- [x] **Alle vier** Punkte stehen mit **Zitat und Fundstelle** — wahr erst seit
      dem Review: Punkt 4 trug beides nicht (§9). Acht Blockzitate, alle acht
      maschinell **wörtlich** gegen die vendorte Quelle gehalten.
- [x] Punkt 1 nennt die gebaute Form und ihre Zahlen: 137 Slices, null Treffer,
      **vier** konstruierte Platzhalter-Verstöße — nicht sechs, die Zahl mischte
      zwei Proben-Klassen — dazu zwei Fail-Closed-Proben. Punkt 3 und 4 nennen
      je **drei** Instanzen mit Fundstelle.
- [x] Je Punkt steht da, was der CR **nicht** verlangt; bei Punkt 3 und 4
      ausdrücklich **kein Gate**. Vom Review je Punkt gegengeprüft.
- [x] Das Dokument liegt unter
      [`docs/plan/cr/`](../../cr/2026-08-25-cr-regelwerk-v5110.md); das
      Absenden bleibt ausdrücklich offen — Auftraggeber-Entscheidung.
- [x] `make gates` Exit 0 (zehn Glieder); unabhängiger Review
      ([Report](../../../reviews/2026-08-25-slice-140-konsumenten-cr-review.md)),
      blockierend mit **drei MEDIUM**, alle sechs Befunde eingearbeitet — dazu
      ein siebter aus der eigenen Gegenprobe.

## 5. Abnahme-Punkte / Risiken

- **Ein CR aus dem eigenen Schmerz heraus ist selten allgemein.** Beide Punkte
  müssen zeigen, dass sie **jedes** Adopter-Repo treffen — sonst sind sie eine
  repo-lokale Adaption, und dafür gibt es den Konventionsspeicher. —
  **Ausgang:** *eingetreten — in Punkt 4, und zwar genau als der Fehler, den
  Punkt 4 beschreibt.* Er erklärte die Reichweitenfrage pauschal für ungeregelt,
  während die von ihm selbst genannte Datei sie zweimal regelt. Der Review fand
  es, die Bitte ist auf das Belegte zurückgeschnitten. Kein Carveout und kein
  Folge-Slice: der Rest ist null und gemessen — acht von acht Blockzitaten
  wörtlich belegt.
- **Vier Punkte in einem CR laden zur Teil-Annahme ein.** Beim letzten CR
  dieses Repos wurden zwei von mehreren Punkten abgelehnt — das ist gesund,
  solange jeder Punkt für sich steht. Kippt aber die Nummerierung zur
  Wunschliste, sinkt die Chance für alle. Zu prüfen ist deshalb je Punkt, ob er
  **ohne** die anderen trägt. — **Ausgang:** *entfallen — die Prüfung ist
  gelaufen.* Der Review hat je Punkt geprüft, ob er allein trägt; Ergebnis: alle
  vier. Teil-Annahme bleibt möglich und ist gesund — sie war nie das Risiko;
  das Risiko war die Kippe zur Wunschliste, und die ist geprüft.
- **Punkt 2 könnte sich als Lesefehler erweisen.** Wenn der Kanon anderswo
  klärt, was *Modul-Pfad* heißt, ist der CR gegenstandslos und jene Adaption
  falsch. Vor dem Schreiben ist der ganze Kanon danach zu durchsuchen, nicht nur
  die zwei bekannten Stellen. — **Ausgang:** *entfallen für Punkt 2 —
  widerlegt, doppelt.* Die Baum-Suche fand **genau zwei** Fundstellen, keine
  dritte, keine andere Schreibweise; der Reviewer hat das unabhängig
  reproduziert, Schreibvarianten eingeschlossen. Punkt 2 ist damit die stärkste
  Stelle des CR. **Der Preis steht in §9:** das Risiko war auf Punkt 2
  geschnitten statt auf seine Klasse — und für Punkt 4 hat darum niemand
  gesucht.

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

Geliefert: vier Punkte gegen Regelwerk `v5.11.0`, jeder mit Zitat, Fundstelle
und der ausdrücklichen Angabe, was er **nicht** verlangt. Keiner bittet um ein
Gate. Das Dokument liegt vor; ob es beim Kurs eingeht, ist nicht entschieden.

**Der CR über Zitat-Reichweite hat selbst überzogen.** Drei der sechs
Review-Befunde sind [`BEO-012`](../observations.md)-Klasse — die Klasse, die
dieser CR zu seinem Punkt 4 macht. Der teuerste sitzt in Punkt 4 selbst: er
erklärte die Reichweitenfrage für ungeregelt, während
`grundlagen-source-precedence.md` — die Datei, die er als Quelle nennt — sie
zweimal beantwortet, für Adaptions-Einträge und für den eigenen
Rangordnungs-Satz. Die Bitte ist jetzt die schmalere und die belegte: nicht
*„das fehlt"*, sondern *„das steht zweimal für den Einzelfall und nie als
Frage"*.

**Warum es niemand suchte, steht in §5.** Das Risiko *„könnte sich als
Lesefehler erweisen"* war auf Punkt 2 geschnitten — auf die Stelle, an der der
Zweifel entstand, statt auf seine Klasse: **jeder** Punkt, der eine Kanon-Lücke
behauptet. Für Punkt 2 lief die Baum-Suche und entkräftete das Risiko doppelt;
für Punkt 4 lief sie nie. Das ist
[`BEO-011`](../observations.md) Ausprägung (c), begangen im Risiko-Abschnitt
selbst.

**Der siebte Befund kam nicht vom Reviewer.** Nach dem Einarbeiten habe ich alle
acht Blockzitate maschinell gegen den Baseline-Baum gehalten. Das
**Eröffnungszitat** von Punkt 1 — der meistgelesene Satz des Dokuments — tilgte
drei Einschübe ohne Auslassungszeichen, darunter ausgerechnet *„(ohne sie ist es
stilles Vergessen)"*. Es stand seit der ersten Fassung da und ist keinem der
beiden Leser aufgefallen.

**Die Gegenprobe selbst war beim ersten Lauf falsch** und meldete vier
Fehlschläge, weil sie die Auszeichnung nur auf **einer** Seite entfernte. Erst
beidseitig normalisiert blieb der eine echte Fund übrig. Dieselbe Sorte
Messfehler wie beim Abschnitts-Skopus zwei Slices zuvor — ein Messskript, das
weniger sieht als das, was es prüfen soll, ist kein Beleg, sondern ein zweiter
ungewachter Spiegel.

**Auftraggeber-Vorgabe während der Einarbeitung: keine Forensik im CR.** Ich
hatte einen Absatz eingefügt, der dem Kurs die Entstehungsgeschichte unseres
eigenen Punktes erzählt. Er ist wieder raus. Das Dokument trägt Bitte und Beleg;
was der CR über sich selbst zu sagen hätte, steht hier.

**Was hielt.** Punkt 2 ist die stärkste Stelle: die Erschöpfungs-Behauptung
*„genau zwei Fundstellen"* hat der Reviewer unabhängig über den ganzen Baum
reproduziert. Punkt 3 trägt drei belegte Instanzen — nach der Berichtigung aus
**zwei** Arbeitstagen statt aus einem; die Dichte-Behauptung stammte aus der
Datums-Verwechslung, die am Kopf derselben Datei längst berichtigt war.

**Register:** [`BEO-012`](../observations.md) auf Zähler **4**,
[`BEO-011`](../observations.md) auf Zähler **4** (Ausprägung (c), erstmals aus
einem Risiko-Abschnitt). Für den fehlenden Reviewer-Anker zu `BEO-012` liegt
jetzt [slice-147](../open/slice-147-reviewer-anker-reichweite.md) in `open/` —
die Klasse hat viermal zugeschlagen und wird jedes Mal vom zweiten Leser
gefunden, nie vom Schreibenden. Das ist eine Feedforward-Lücke, und die gehört
nicht in eine Closure-Notiz, sondern in eine Datei.
