# Slice slice-208: Regel- und Template-Adoption des `v6.5.0`-Deltas

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [slice-207](../open/slice-207-baseline-v650-bump.md) (misst den Delta und
hebt den Pin; dieser Slice **urteilt** über ihn),
[Antwort auf den ausgehenden CR](../../cr/2026-09-06-cr-ai-harness-course-slice-formluecken.md)
(die `v6.4.0`-Hälfte des Deltas ist die Umsetzung unserer beiden angenommenen
Bitten), [`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage)
(Adoptions-Erklärung).

**Berührte Spec-Stellen:** — *(voraussichtlich keine; sollte der Delta eine
Spec-Stelle berühren, wird dieser Kopf bei der Beanspruchung nachgezogen)*

**Verantwortlich:** — · **Autor:** pt9912. **Datum:** 2026-09-06.

---

## 1. Ziel

Je Regel des gemessenen Deltas **eine Antwort**: übernommen · nicht anwendbar
mit Begründung · abweichend als deklarierte Adaption. Kern der bekannten Hälfte
ist die **Auflösung der Slice-Haus-Form**, die der Kanon mit `v6.4.0`
entbehrlich macht.

## 2. Vorgehen

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Delta-Audit | neu | je Regel eine Antwort, im Slice festgehalten — die Form, die der vorige Adoptions-Slice gelebt hat |
| Slice-Vorlage und Haus-Form | update | `v6.4.0` führt den Ausschluss als zweite Hälfte von §1 und §8 mit unbedingtem Kopf; unsere §1/§3/§5-Form wird damit entbehrlich |
| [`.d-check.closure.yml`](../../../../.d-check.closure.yml) | update | zwei Regeln keilen heute auf Haus-Form-Titel; eine bindet `## 5. Abnahme-Punkte / Risiken` **wörtlich** |
| [`harness/conventions/`](../../../../harness/conventions/) | update | was abweicht, wird deklarierte Adaption; was der Kanon jetzt selbst sagt, löst seinen Eintrag auf |

**Der Bestand ist gemessen und klein — aber er ist eingefroren.** **Sechs**
`done/`-Slices tragen die Haus-Form (`## 3. Ausdrücklich NICHT`,
`## 5. Abnahme-Punkte / Risiken`). Sie sind Lauf-Belege und werden **nicht**
nachgezogen; ein umgeschriebener DoD-Punkt fälschte einen Beleg. Die
Gate-Regeln müssen deshalb **beide** Formen tragen oder sauber nach
Zeitpunkt geschieden werden — das ist die eigentliche Arbeit, nicht das
Umbenennen der Abschnitte.

**Was der Kanon selbst sagt, braucht keinen Eintrag mehr.** Die Antwort auf
unseren CR nennt es ausdrücklich: „euren lokalen Fork könnt ihr dann
auflösen". Ob daraus die **Auflösung** eines Konventions-Eintrags folgt oder
nur seine Umformulierung, ist je Eintrag zu entscheiden.

## 3. Ausdrücklich NICHT in diesem Slice

- **Die Pin-Hebung selbst.** Sie liegt in
  [slice-207](../open/slice-207-baseline-v650-bump.md); ohne den dortigen `diff -I`-Beleg
  hat dieser Slice keinen Gegenstand.
- **Ein Retrofit des `done/`-Bestands.** Die sechs Slices in Haus-Form bleiben,
  wie sie sind — eingefrorene Lauf-Belege.
- **Jede Regel, die der Delta nicht berührt.** Der Slice adoptiert, was
  `v6.4.0`/`v6.5.0` ändern, und benutzt die Gelegenheit nicht, um Nachbarregeln
  mitzunehmen.

## 4. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** Zu **jeder** Regel des gemessenen Deltas steht eine Antwort im
      Slice: übernommen (mit Träger) · nicht anwendbar (mit Begründung) ·
      abweichend (mit Adaptions-Eintrag). Eine Regel ohne Antwort ist ein
      offener Punkt, kein stilles Übergehen.
- [ ] **(2)** Die Slice-Haus-Form ist aufgelöst: die Vorlage folgt der
      Baseline-Form, und die Regeln in
      [`.d-check.closure.yml`](../../../../.d-check.closure.yml) tragen den
      Bestand **und** die neue Form — mit einem **Bruch-Test je Richtung**, der
      belegt, dass beide noch gefangen werden.
- [ ] **(3)** Die Konventions-Einträge sind nachgezogen: was der Kanon jetzt
      selbst sagt, ist aufgelöst; was abweicht, ist deklariert.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §5 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 5. Abnahme-Punkte / Risiken

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Dieser Slice wird wahrscheinlich die Ein-Sitzungs-Review-Grenze
  überschreiten**, und das ist **vorab** bekannt: Sein Vorgänger — dieselbe
  Arbeit über einen kleineren Delta — wuchs von drei auf 23 Träger. Wird er
  nicht zurückgeführt, verlangt
  [`MR-066`](../../../../harness/conventions.md#mr-066) **zweierlei**: den
  Grund **und** die **Ersatz-Form der Prüfung**, benannt im Plan und vollzogen
  im Report. Beides gehört bei der Beanspruchung in §6, nicht erst in den
  Review. — **Ausgang:** \<offen\>
- **Der Umfang steht erst nach slice-207 fest.** `v6.4.0` ist angekündigt und
  bekannt, `v6.5.0` nicht. Trägt es eine eigene Regel-Änderung, ist dieser
  Slice **vor** der Beanspruchung neu zu schneiden — eine Schätzung jetzt wäre
  aus dem Anlass gezogen und nicht aus dem Bestand
  ([`BEO-ALL/rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md),
  7×). — **Ausgang:** \<offen\>
- **Die Gate-Regeln sind die stille Stelle.** Eine Regel, die auf einen
  Haus-Form-Titel keilt, wird nach der Umbenennung **grün, ohne noch etwas zu
  prüfen** — dieselbe Klasse wie ein Sensor, der auf `done/*.md` keilt und die
  archivierten Stubs nicht mehr sieht. Deshalb der Bruch-Test je Richtung in
  DoD (2) und nicht bloß ein grüner Lauf. — **Ausgang:** \<offen\>
- **Der `done/`-Bestand und die neue Form leben nebeneinander.** Solange beide
  existieren, ist „die Form eines Slice" zweideutig, und jede Regel darüber
  braucht eine Zeit- oder Verzeichnis-Grenze. Ob die sechs Bestands-Slices als
  feste Liste oder über eine Ziffern-Schwelle ausgenommen werden, ist ein
  Entscheid — die Ziffern-Schwelle ist die gelebte Form
  ([`MR-056`](../../../../harness/conventions.md#mr-056)). — **Ausgang:** \<offen\>

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-207](../open/slice-207-baseline-v650-bump.md)
liegt in `done/`, der Delta ist mit `diff -I` gemessen und als Liste
festgehalten. WIP-Limit frei.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Erweist sich die Haus-Form-Auflösung als
  eigene Arbeit neben dem übrigen Delta — insbesondere wenn `v6.5.0` eine
  zweite Regel-Änderung trägt —, wird sie ein eigener Slice. **Die
  Rückführung ist hier der Regelfall, nicht die Ausnahme**; wer sie nicht
  zieht, schuldet [`MR-066`](../../../../harness/conventions.md#mr-066).
- `in-progress` → `open` (blockiert): Zeigt sich, dass eine Kanon-Regel dem
  gelebten Bestand widerspricht, ohne dass eine Adaption sie trägt, ruht der
  Slice bis zum Entscheid — Adoption ist keine Erlaubnis, den Widerspruch
  stillschweigend nach einer Seite aufzulösen
  ([`AGENTS.md`](../../../../AGENTS.md) §1).

**Closure-Trigger.** Zwei beobachtbare Kriterien und ein Lerneintrag: (a) zu
jeder Delta-Regel steht eine Antwort und `make gates` ist grün; (b) die
Bruch-Tests aus DoD (2) sind gefahren und ihre Ausgabe steht im Report.

## 7. Vorgelagert (vor der Modus-Begründung)

\<entsteht spätestens bei der Beanspruchung — ein Plan in `open/` trägt die drei
Vorprüfungen noch nicht\>

## 8. Sub-Area-Modus-Begründung

\<entsteht mit den Vorprüfungen bei der Beanspruchung\>

## 9. Closure-Notiz (nach `done/`)

\<wird vor dem `git mv` nach `done/` gefüllt\>
