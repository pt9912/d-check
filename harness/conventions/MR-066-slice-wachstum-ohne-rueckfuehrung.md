# MR-066 — Wird die Rückführung nicht gezogen, wird die Review-Last geteilt statt des Slice

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon kennt für den zu groß gewordenen
  Slice **eine** Antwort — die Rückführung `in-progress` → `next`
  ([`modul-05-planning-harness.md` §Trigger je Lifecycle-Übergang](../../.harness/baseline/v6.3.1/regelwerk/modul-05-planning-harness.md#trigger-je-lifecycle-übergang-und-wip-limit-modul-5)).
  Er sagt nichts über den Fall, dass ihr etwas **entgegensteht**. Dieser
  Eintrag ergänzt ihn dort und weicht nicht von ihm ab: Die Rückführung bleibt
  der Regelfall.
- **Datum:** 2026-09-06 · **Herkunft:** seit slice-204 (Steering Loop,
  `BEO-ALL/large-migration-exceeds-session-review-limit` 3×)
- **Geltungsbereich:** jeder Slice, dessen Umfang die Ein-Sitzungs-Review-Grenze
  überschreitet und der **nicht** zurückgeführt wird. **Nicht** erfasst: ein
  Slice, der die Grenze hält, und einer, der regulär zurückgeführt wird — für
  beide gilt der Kanon unverändert.
- **Adaption:** Wird die Rückführung nicht gezogen, sind **zwei** Dinge fällig,
  und das zweite ist der eigentliche Gehalt dieser Regel:

  1. **Der Grund steht geschrieben** — im Slice-Plan. Vorab, wo er vorab
     bekannt ist; **beim Überschreiten**, wo er erst während der Arbeit
     entsteht. Der zweite Fall ist der unauffälligere und verlangt, dass das
     Wachstum überhaupt **bemerkt** wird: Wächst ein Slice schrittweise, sprengt
     kein einzelner Schritt die Grenze — nur ihre Summe, und die sieht niemand,
     der nur auf den nächsten Schritt schaut.
  2. **Die Review-Last bekommt eine benannte Ersatz-Form.** Nicht der Slice
     wird geteilt, sondern die Prüfung: mehrere Review-Runden gegen je einen
     abgeschlossenen Zwischenstand · ein gezielter Stichproben-Fokus statt
     vollständiger Zeile-für-Zeile-Prüfung · eine andere, die der Plan benennt.
     **Ohne diesen zweiten Teil ist es kein Ausnahmefall, sondern nur ein zu
     großer Slice mit einer Begründung davor.**

  **Warum das die Größenregel nicht aufweicht.** Sie existiert aus **einem**
  Grund: dass ein Reviewer den Diff prüfen kann. Wo die Teilung diesen Grund
  besser bedient, wird geteilt — das ist der Regelfall und bleibt es. Wo sie
  ihn *verschlechtert*, weil sie einen Beleg zerreißt, der nur im Zusammenhang
  trägt, ist die Teilung das falsche Mittel für das richtige Ziel; dann wird
  das Ziel direkt verfolgt.

  **Gezogen aus drei Instanzen, und die dritte ist die, an der die Regel sich
  bewähren muss** (`BEO-ALL/large-migration-exceeds-session-review-limit`, 3×):

  | Vorgang | Warum nicht zurückgeführt | Ersatz-Form |
  |---|---|---|
  | slice-195 | eine Teilung hätte den Zähler-Diff-Beleg zerrissen | gezielter Stichproben-Fokus |
  | slice-197 | eine Teilung hätte die Werkzeug-Korrektur vervielfacht | **keine benannt** |
  | slice-204s Anlass (slice-203) | kein einzelner Nachsteuerungs-Schritt sprengte die Grenze | zwei Review-Runden gegen je einen abgeschlossenen Stand |

  **Die mittlere Zeile ist der Grund für Teil 2.** In slice-197 wurde die
  Teilung mit gutem Grund unterlassen und **nichts** an ihre Stelle gesetzt;
  die Grenze war damit einfach weg. Wer nur Teil 1 verlangt, hätte diesen Fall
  durchgewinkt.

  **Was diese Regel nicht leistet, ausgeschrieben:**

  - Sie rettet die **Prüfbarkeit**, nicht die Größe. Der Diff bleibt groß, der
    Merge-Konflikt-Raum auch, und ein langlebiger Branch bleibt ein Risiko.
  - Sie ist **kein Sensor.** Ob ein Slice die Ein-Sitzungs-Grenze überschreitet,
    ist ein Urteil; eine Zeichen- oder Dateizahl daraus zu machen, tauschte ein
    ehrliches Urteil gegen falsche Genauigkeit. Kein Gate prüft diesen Eintrag.
  - Sie sagt **nicht**, wie viele Runden genügen. Zwei waren es einmal; die
    Zahl folgt dem Gegenstand, nicht der Regel.
  - Sie ist aus **drei** Anlässen gezogen, nicht aus einer Inventur
    (`BEO-ALL/rule-drawn-from-occasion-not-inventory`). Ob sie trägt,
    entscheidet die vierte Instanz — bis dahin ist sie belegt, nicht bewiesen.

- **Begründung:** Dreimal stand dieselbe Spannung, dreimal wurde sie einzeln
  aufgelöst, und beim dritten Mal fiel auf, dass die Auflösungen sich
  unterscheiden — einmal mit Ersatz, einmal ohne. Ohne geschriebene Regel
  entscheidet jeder Lauf neu, und der bequeme Ausgang ist immer, die
  Rückführung zu unterlassen und nichts an ihre Stelle zu setzen.
- **Ausgelöst durch Baseline-Stand:** v6.3.1
- **Auflösungs-Trigger:** der Kanon nimmt den Fall selbst auf — dann gilt
  seiner. Bis dahin permanent.
