# MR-066 — Was „nicht still weiterschieben" verlangt, wenn die Rückführung nicht gezogen wird

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine — der Eintrag **präzisiert** eine
  Kanon-Regel, statt von ihr abzuweichen.
  [`modul-05-planning-harness.md` §Ziel-Form: Slice](../../.harness/baseline/v6.3.1/regelwerk/modul-05-planning-harness.md#ziel-form-slice)
  sagt für den zu groß gewordenen Slice: *„Dann zurück zum Schneiden
  (`in-progress→next`), nicht still weiterschieben."* Der Satz verbietet das
  **stille** Weiterschieben — er sagt nicht, was an die Stelle der Stille
  tritt, wenn die Rückführung aus einem Sachgrund unterbleibt. Genau diese
  Lücke füllt der Eintrag. Die Rückführung bleibt der Regelfall; sie wird
  weder abgeschafft noch aufgeweicht.
- **Datum:** 2026-09-06 · **Herkunft:** seit slice-204 (Steering Loop,
  `BEO-ALL/large-migration-exceeds-session-review-limit` 3×)
- **Geltungsbereich:** jeder Slice, dessen Umfang die Ein-Sitzungs-Review-Grenze
  überschreitet und der **nicht** zurückgeführt wird. **Nicht** erfasst: ein
  Slice, der die Grenze hält, und einer, der regulär zurückgeführt wird — für
  beide gilt der Kanon unverändert.
- **Adaption:** Zwei Pflichten, und die zweite ist der Gehalt:

  1. **Der Grund steht im Slice-Plan.** Vorab, wo er vorab bekannt ist;
     **beim Überschreiten**, wo er erst während der Arbeit entsteht.
  2. **Die Ersatz-Form der Prüfung steht ebenfalls im Plan — und ihr Vollzug
     im Review-Report.** Nicht der Slice wird geteilt, sondern die Prüfung:
     mehrere Runden gegen je einen abgeschlossenen Stand · ein deklarierter
     Stichproben-Fokus · eine mechanische Vollprüfung entlang einer Achse
     (Byte-Diff, Feld-Vergleich) · eine andere, die der Plan benennt.

  **Warum die Deklaration und nicht die Praxis der Gegenstand ist.** Der
  Befund aus den drei Instanzen ist nicht, dass die Prüfung ausblieb — sie
  fand jedes Mal statt, mit sehr unterschiedlicher Gründlichkeit. Der Befund
  ist, dass sie **nirgends vorher stand**: Wie gründlich geprüft wurde, erfuhr
  man erst hinterher aus dem Report, und damit war die Prüftiefe weder planbar
  noch einforderbar. Eine Regel, die nur den *Grund* verlangt, ändert daran
  nichts.

  | Vorgang | Warum nicht zurückgeführt | Was tatsächlich geprüft wurde | Vorab benannt? |
  |---|---|---|---|
  | slice-195 | eine Teilung hätte den Zähler-Diff-Beleg zerrissen | Stichprobe, 8 von 29 Einträgen | nein |
  | slice-197 | eine Teilung hätte die Werkzeug-Korrektur vervielfacht | Byte-Diff, Titel- und Feld-Vergleich über **alle 45**, dazu drei Dateien Zeile für Zeile — fand so ein HIGH | nein |
  | slice-203 | kein einzelner Nachsteuerungs-Schritt sprengte die Grenze | zwei Review-Runden, die zweite über den gesamten Bereich | nein |

  **Die mittlere Zeile ist der Grund für die Form von Pflicht 2.** Ausgerechnet
  dort war der Prüfumfang am größten — und der Slice-Plan sagte darüber nichts.
  Wer nur die Reports vergleicht, hielte slice-197 für den ungedeckten Fall;
  wer die Praxis misst, findet das Gegenteil. Deklariert wird deshalb die
  **Ersatz-Form**, nicht ihre Gründlichkeit: Die ist ein Urteil, ihre
  Benennung nicht.

  **Zwei Lagen, und nur die erste ist eine Abwägung.** In slice-195 und
  slice-197 war die Teilung *möglich* und hätte die Prüfbarkeit
  **verschlechtert** — dort ist die Nicht-Rückführung eine begründete Wahl.
  In slice-203 stellte sich die Frage nie: Das Wachstum lief schrittweise, kein
  einzelner Schritt sprengte die Grenze, nur ihre Summe. **Für diese zweite
  Lage ist die erste Pflicht das Bemerken** — und danach steht die Wahl
  zwischen Rückführung und Pflicht 2 offen wie in der ersten Lage. Wer die
  zweite Lage mit der ersten begründet, hat sie nicht verstanden.

  **Was diese Regel nicht leistet, ausgeschrieben:**

  - Sie rettet die **Prüfbarkeit**, nicht die Größe. Der Diff bleibt groß, der
    Merge-Konflikt-Raum auch, und ein langlebiger Branch bleibt ein Risiko.
  - Sie ist **kein Sensor.** Ob ein Slice die Ein-Sitzungs-Grenze überschreitet,
    ist ein Urteil; eine Zeichen- oder Dateizahl daraus zu machen, tauschte ein
    ehrliches Urteil gegen falsche Genauigkeit. Kein Gate prüft diesen Eintrag.
  - Sie prüft **nicht die Gründlichkeit** der Ersatz-Form, nur ihre Benennung
    und ihren Vollzug. Eine deklarierte Stichprobe von zwei Zeilen erfüllt sie
    formal — dass das zu wenig ist, bleibt Review-Urteil.
  - Sie ist aus **drei** Anlässen gezogen, nicht aus einer Inventur
    (`BEO-ALL/rule-drawn-from-occasion-not-inventory`). Ob sie trägt,
    entscheidet die vierte Instanz — bis dahin ist sie belegt, nicht bewiesen.

- **Begründung:** Dreimal stand dieselbe Spannung, dreimal wurde sie einzeln
  aufgelöst, und in keinem der drei Fälle stand vorher, wie die verlorene
  Prüftiefe aufgefangen wird. Ohne geschriebene Regel entscheidet jeder Lauf
  neu und im Nachhinein; der bequeme Ausgang ist, die Rückführung zu
  unterlassen und die Frage gar nicht zu stellen.
- **Auflösungs-Trigger:** der Kanon sagt selbst, was an die Stelle der Stille
  tritt — dann gilt seiner. Bis dahin permanent.
