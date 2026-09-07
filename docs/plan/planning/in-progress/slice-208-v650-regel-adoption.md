# Slice slice-208: Regel- und Template-Adoption des `v6.5.0`-Deltas

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [slice-207](../done/slice-207-baseline-v650-bump.md) (misst den Delta und
hebt den Pin; dieser Slice **urteilt** über ihn),
[Antwort auf den ausgehenden CR](../../cr/2026-09-06-cr-ai-harness-course-slice-formluecken.md)
(die `v6.4.0`-Hälfte des Deltas ist die Umsetzung unserer beiden angenommenen
Bitten), [`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage)
(Adoptions-Erklärung).

**Berührte Spec-Stellen:** — *(voraussichtlich keine; sollte der Delta eine
Spec-Stelle berühren, wird dieser Kopf bei der Beanspruchung nachgezogen)*

**Verantwortlich:** pt9912 · **Autor:** pt9912. **Datum:** 2026-09-06.

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

### Der Delta-Audit (DoD 1)

**Zuerst die Form des Gegenstands, wie der Register-Eintrag es verlangt.** Der
Vorgänger hat **zwölf Dateien** mit Delta gemessen. Eine Datei ist aber keine
Regel: Dieselbe Regel steht oft in mehreren Trägern, und ein Träger kann eine
Regel nur nachziehen. Gezählt wird deshalb, was ein Implementer **tun oder
lassen** muss. Ergebnis: **sieben Regeln in zwölf Trägern.**

| # | Regel | Träger | Antwort |
|---|---|---|---|
| R1 | §1 heißt *Ziel und Abgrenzung* — vier Klassen, Begründung je Punkt, keine Mindestzahl, kein Sensor | `modul-05`, `slice.template`, `templates/README`, `modul-06`, `modul-09` | **übernommen** |
| R2 | *„Die Adresse muss die Sendung annehmen"* — ein Folge-Slice, der den verwiesenen Punkt selbst ausschließt oder **vor** dem verweisenden schließt, ist keine Adresse | `modul-05` | **übernommen** |
| R3 | §8 heißt *Sub-Area-Prüfungen und Modus-Begründung*; die zwei *Vorgelagert*-Blöcke sind unbedingter Kopf | `slice.template`, `templates/README` | **übernommen** |
| R4 | Der Lauf weitet die Abgrenzung nicht stillschweigend — Mitnahme ist eine **Plan-Änderung** und gehört vor den Code | `modul-09` | **übernommen** |
| R5 | Die **RTM** wird *erzeugt, nicht gepflegt*; was eine Anforderung **entlastet**, ist eine Setzung und gehört aufgeschrieben | `grundlagen-traceability`, `grundlagen-begriffe` | **übernommen, mit Handlung** |
| R6 | **Einfrierend gegen lebend**: ein lebendes Artefakt verlinkt `harness/sensors/<target>.md`, ein einfrierendes nennt `make <target>` als Token | `grundlagen-harness-dateien` | **übernommen, Bestand konform** |
| R7 | **Zitier-Form** einfrierender Artefakte: Kennung statt Adresse, Baseline-Stelle als Tag + Pfad in Inline-Code statt als Link | `review-report`, `archiv-stub-slice`, `archiv-stub-welle`, `welle-results` | **übernommen, template-forward** |

**R1–R4 sind die CR-Umsetzung** und lösen d-checks Haus-Form auf. §1 und §3
fallen zusammen, §7 und §8 ebenso, §6 spaltet sich in Trigger und
Closure-Trigger; neun Abschnitte werden acht. Die Zuordnung ist **keine
Bijektion**, und deshalb ist es keine Umbenennung, sondern eine Migration.

**R2 ist neu und wurde nicht erbeten.** Sie schärft die Klasse *„ein
Folge-Slice übernimmt es"*, die dieses Repo selbst geliefert hat — und sie
trifft eine reale Praxis: Von den sechs gelebten Ausschluss-Abschnitten nennen
vier einen Folge-Slice, und ob der die Sendung annimmt, hat bisher niemand
geprüft.

**R5 verlangt als einzige eine Handlung am Bestand, und sie ist klein.**
d-checks `--trace` **erzeugt** die RTM aus den Ankern — die erste Hälfte ist
also erfüllt, ohne dass etwas zu tun wäre. Die zweite ist offen: Der
`trace`-Block deklariert `adrs` und `slices` als Quellen, die eine Anforderung
entlasten, **und nennt nirgends den Grund**. Genau davor warnt der Kanon: Das
ist eine Konfigurationsentscheidung, keine Naturgesetzlichkeit — dieselbe ADR
entlastet oder entlastet nicht, je nachdem, was das Repo gesetzt hat.

**R6 ist gemessen und der Bestand trägt sie schon.** 26 Verweise aus lebenden
Artefakten verlinken die Sensor-Datei — richtig. Aus dem eingefrorenen Bestand
tut es **genau einer**, und der steht in einem Review-Report, der nicht mehr
angefasst wird. Kein Handlungsbedarf am Bestand; die Regel greift vorwärts.

**R7 ist die Regel, die der Pin-Hebung ihre Ausnahme nimmt.** Die letzte
Hebung brauchte einen Tombstone, weil eine `Accepted`-ADR einen **Link** in den
vendorten Baum trug. Unter R7 hätte sie ihn nie getragen. Träger bei uns ist
nicht die Baseline-Vorlage — dieses Repo führt eine eigene Review-Form —,
sondern der **Reviewer-Skill**.

**Keine Regel des Deltas ist *nicht anwendbar*, und keine wird *abweichend*
adoptiert.** Das ist ein Befund und keine Selbstverständlichkeit: Die beiden
Vorgänger-Adoptionen trugen je mehrere Nicht-anwendbar-Antworten (Wellen-Betrieb,
Mehr-Schreiber-Teile). Dieser Delta trifft ausschließlich Bereiche, in denen
dieses Repo aktiv arbeitet.

## 3. Ausdrücklich NICHT in diesem Slice

- **Die Pin-Hebung selbst.** Sie liegt in
  [slice-207](../done/slice-207-baseline-v650-bump.md); ohne den dortigen `diff -I`-Beleg
  hat dieser Slice keinen Gegenstand.
- **Ein Retrofit des `done/`-Bestands.** Die sechs Slices in Haus-Form bleiben,
  wie sie sind — eingefrorene Lauf-Belege.
- **Jede Regel, die der Delta nicht berührt.** Der Slice adoptiert, was
  `v6.4.0`/`v6.5.0` ändern, und benutzt die Gelegenheit nicht, um Nachbarregeln
  mitzunehmen.

## 4. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [x] **(1)** Zu **jeder** Regel des gemessenen Deltas steht eine Antwort im
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

- **Zwei `v6.5.0`-Regeln gehen über eine Template-Adoption hinaus** — beim
  Anlegen dieses Plans war das unbekannt, die Delta-Messung des Vorgängers hat
  es gezeigt. **(a) Die RTM ist Kanon-Begriff geworden**
  (`grundlagen-traceability.md` §Die zweite Richtung), mit zwei Setzungen:
  *„sie wird **erzeugt**, nicht gepflegt"* und *„was eine Anforderung
  **entlastet**, ist eine Konfigurationsentscheidung und gehört
  aufgeschrieben"*. Das beschreibt `--trace` — **unser Produkt** —, und ob
  daraus eine Spec-Frage folgt oder nur eine Bestätigung, ist offen. **(b) Die
  Unterscheidung *einfrierend / lebend*** verlangt für einfrierende Artefakte
  die **Kennung statt der Adresse**; sie berührt Bestand (Review-Reports,
  `done/`-Slices) und macht den ADR-Tombstone der letzten Hebung künftig
  entbehrlich. Beide sind **Urteile**, keine Umbenennungen, und beide können
  den Zuschnitt sprengen. — **Ausgang:** \<offen\>
- **Die Umnummerierung bewegt JEDEN Abschnitt** — gemessen: neun Haus-Form-
  Abschnitte gegen acht der Baseline, und die Zuordnung ist keine Bijektion
  (§1+§3 fallen zusammen, §6 spaltet sich in §4+§5). **Gate-seitig ist es
  weniger, als es aussieht:** Von drei Regeln im Closure-Profil, die einen
  Abschnitt adressieren, keilt genau **eine** auf einen wörtlichen Titel
  (`## 5. Abnahme-Punkte / Risiken`); die beiden anderen sind Muster über
  `Definition of Done` und überleben. Die Zahl ist gemessen, nicht geschätzt —
  aber sie ist eine Zählung, und der Register-Eintrag zu Zählungen ist
  gesichtet. — **Ausgang:** \<offen\>
## 6. Trigger

**Start** (`open` → `in-progress`): [slice-207](../done/slice-207-baseline-v650-bump.md)
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

**[`MR-066`](../../../../harness/conventions.md#mr-066) — vorab, nicht im
Nachhinein.** Dieser Slice überschreitet die Ein-Sitzungs-Review-Grenze
voraussichtlich: Sein Vorgänger derselben Art wuchs von drei auf 23 Träger,
und hier kommen zwei Urteils-Fragen (§5) zu einer Form-Migration hinzu, die
jeden Abschnitt bewegt. **Der Grund, ihn nicht zurückzuführen:** Eine Teilung
zerrisse den Delta-Audit — die Frage *„trägt dieser Eintrag noch?"* ist je
Regel nur beantwortbar, wenn der ganze Delta danebenliegt, und die
Form-Migration hängt an derselben Liste.

**Die Ersatz-Form der Prüfung, benannt:** **Zwei Review-Runden gegen je einen
abgeschlossenen Stand** — Runde 1 gegen den **Delta-Audit** (DoD 1), sobald zu
jeder Regel eine Antwort steht und **bevor** die Form angefasst wird; Runde 2
gegen die **Form-Migration** (DoD 2 und 3). Das teilt die Last nach
Gegenstand, nicht nach Umfang, und jede Runde prüft einen Stand, der für sich
vollständig ist. **Nicht** dasselbe wie zwei Runden über denselben Bereich:
Runde 1 sieht die Migration nicht, Runde 2 den Audit als gegeben.

Der Vollzug gehört in die **Reports**, nicht in die Closure-Notiz — beide
Runden liegen unter [`docs/reviews/`](../../../reviews/), und ihre
Gegenstands-Zeile nennt den jeweils geprüften Stand.

**Closure-Trigger.** Zwei beobachtbare Kriterien und ein Lerneintrag: (a) zu
jeder Delta-Regel steht eine Antwort und `make gates` ist grün; (b) die
Bruch-Tests aus DoD (2) sind gefahren und ihre Ausgabe steht im Report.

## 7. Vorgelagert (vor der Modus-Begründung)

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:268-269 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area: `*` (Repo-Default). Der Slice ändert Konventionen, die
Slice-Vorlage und ein Prüf-Profil — alles unter dem Default.
`tools/harness/` ist **nicht** berührt: Der Delta hat kein Werkzeug angefasst,
und die Hebung selbst liegt hinter uns.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:274-274 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, **36** Verzeichnisse — über **beide**
Kürzel gezählt, `BEO-ALL` und `BEO-HARN`; die Zählung nur über `BEO-ALL` war
der Fehler des Vorgänger-Slice). Vier Einträge sind einschlägig:

- [`semantic-change-body-only-edges-stale`](../observations/BEO-ALL/semantic-change-body-only-edges-stale/observation.md)
  (11×) — **der zentrale.** Dieser Slice ändert eine Semantik, die an vielen
  Stellen gespiegelt ist: Die Abschnitts-Nummern des Slice-Plans stehen in der
  Vorlage, im Prüf-Profil, in Konventions-Einträgen und in der Prosa, die auf
  sie verweist. Sein Ableiter ist hier wörtlich anzuwenden — Spiegel **vor**
  dem Editieren auflisten —, und die Präzisierung aus slice-206 ebenso: Ein
  Spiegel, der eine **abgeleitete Aussage** ist, lässt sich nicht ergreppen.
- [`mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
  (3×, Ausgang *geplant*) — unmittelbar: Eine Abschnitts-Umnummerierung ist
  genau die mechanische Ersetzung, die dreimal zu weit gegriffen hat. **Die
  sechs `done/`-Slices in Haus-Form sind eingefrorene Lauf-Belege und werden
  nicht angefasst.** Der Folge-Slice
  [slice-209](../open/slice-209-frozen-klassen-vor-mechanischer-ersetzung.md)
  schreibt die Regel dazu und wartet auf **diesen** Slice — hier ist sie also
  noch Disziplin, nicht Konvention.
- [`zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
  (2×) — die Delta-Liste aus dem Vorgänger ist eine **Zählung**, und dieser
  Slice urteilt auf ihr. Vor jeder Aussage „so viele Regeln sind betroffen"
  gehört die Form des Gegenstands ausgeschrieben.
- [`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md)
  (7×) — für die Gegenrichtung: Wo der Delta eine Regel **entbehrlich** macht,
  ist der vorhandene Konventions-Eintrag aufzulösen, nicht umzuformulieren. Ein
  Eintrag, der nur noch wiederholt, was der Kanon selbst sagt, ist eine zweite
  Quelle.

**Geprüft und ausgeschlossen:**
[`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
(15×) — der Slice zitiert den neuen Kanon, aber die Zitat-Spannen sind mit der
Hebung bereits neu geankert und von `citations` geprüft. Keiner der vier
erreicht mit diesem Slice die Schwelle erstmalig.

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-07 gelesen. `image-scan.yml` **grün**.
`upstream-drift.yml` meldet weiterhin **ROT** — und hier greift die benannte
Grenze des Targets: Es liest den **jüngsten** Lauf, nicht sein **Alter**. Der
gemeldete Lauf datiert von vor der Pin-Hebung; lokal meldet
`make baseline-freshness` **beide** Achsen grün. Der nächste Nachtlauf löst die
Meldung auf, ohne dass hier etwas zu tun wäre.

## 8. Sub-Area-Modus-Begründung

**Modus:** `*` ist **GF** (Greenfield, Repo-Default).

- **Konventions-Dichte:** hoch, aber **an der falschen Stelle**. Die
  Slice-Form ist dicht geregelt — in der Vorlage, in
  [`MR-049`](../../../../harness/conventions.md#mr-049) und
  [`MR-056`](../../../../harness/conventions.md#mr-056), im Closure-Profil.
  Genau diese Dichte ist der Aufwand: Jede Regel, die einen Abschnitts**titel**
  nennt, muss die Umbenennung überleben.
- **Phase-Reife:** Phase 5 für den Vorgang (Adoption ist zweimal gelebt,
  slice-107 und slice-203), Phase 3 für den **Gegenstand** — die Auflösung der
  Haus-Form ist neu und hat keinen Präzedenzfall.
- **Evidenz-/Diskrepanz-Risiko:** **hoch**, und die drei gesichteten Einträge
  benennen es je einzeln: eine gespiegelte Semantik
  (`semantic-change-body-only-edges-stale`, 11×), eine mechanische Ersetzung
  über eingefrorenen Bestand (`mechanical-id-rewrite-misses-frozen-classes`,
  3×) und eine Zählung als Urteilsgrundlage
  (`zaehlmethode-misst-proxy-statt-gegenstand`, 2×). Das ist die höchste
  Risiko-Dichte, die ein Slice dieses Repos bisher vorab getragen hat.
- **Reconciliation-Aufwand:** keiner (GF). Graduation entfällt.

## 9. Closure-Notiz (nach `done/`)

\<wird vor dem `git mv` nach `done/` gefüllt\>
