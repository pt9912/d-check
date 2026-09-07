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

**Zuerst die Form des Gegenstands.** Der Vorgänger hat **zwölf Dateien** mit
Delta gemessen. Eine Datei ist keine Regel: Dieselbe steht oft in mehreren
Trägern, und ein Träger kann sie **normativ** führen oder nur **nachziehen**.
Gezählt wird, was ein Implementer tun oder lassen muss.

**Die erste Fassung dieses Audits zählte trotzdem datei-weise** — jede Datei
bekam die Regel, die als ihre Hauptaussage gelesen wurde. Alle zwölf waren
abgedeckt, aber nicht alle Regeln *in* ihnen: Eine Regel hatte gar keine
Antwort, zwei waren in Wahrheit eine, und vier Träger fehlten. Der unabhängige
Review hat das gemessen. Die Fassung unten ist regel-weise gebildet — jede
Regel gegen **jede** Datei gehalten.

| # | Regel | Normativ in | Nachgezogen in | Antwort |
|---|---|---|---|---|
| R1 | §1 heißt *Ziel und Abgrenzung* — vier Klassen, Begründung je Punkt, keine Mindestzahl, kein Sensor | `modul-05` | `slice.template`, `templates/README`, `modul-06`, `modul-09` | **übernommen** (DoD 2) |
| R2 | *„Die Adresse muss die Sendung annehmen"* — ein Folge-Slice, der den verwiesenen Punkt selbst ausschließt oder **vor** dem verweisenden schließt, ist keine Adresse | `modul-05` | `slice.template` | **übernommen, Bestand geprüft** |
| R3 | §8 heißt *Sub-Area-Prüfungen und Modus-Begründung*; die zwei *Vorgelagert*-Blöcke sind unbedingter Kopf, der Modus-Block bedingter Rumpf | `modul-05` | `slice.template`, `templates/README` | **übernommen** (DoD 2) |
| R4 | Der Lauf weitet die Abgrenzung nicht — Mitnahme ist eine **Plan-Änderung** und gehört vor den Code | `modul-09` | `modul-05`, `slice.template` | **übernommen, mit Handlung** |
| R5 | Die **RTM** — vier Setzungen (s. u.) | `grundlagen-traceability` | `grundlagen-begriffe` | **drei erfüllt, eine deklarationspflichtig** |
| R6 | **Kennung statt Adresse** für einfrierende Artefakte, in **drei Formen**, für eine **erweiterte** Klasse | `grundlagen-harness-dateien` | `review-report`, `archiv-stub-slice`, `archiv-stub-welle`, `welle-results` | **übernommen, template-forward** |
| R7 | Ein **Ausnahme-Ventil** im Prüfbereich ist *„eine Gate-Senkung mit eigener Begründungslast"* | `grundlagen-harness-dateien` | — | **übernommen, mit Handlung** |

**R1–R4 sind die CR-Umsetzung.** §1 und §3 fallen zusammen, §7 und §8 ebenso,
§6 spaltet sich; neun Abschnitte werden acht. Keine Bijektion, also eine
Migration und keine Umbenennung.

**R2, am Bestand geprüft — und die erste Messung war ein Proxy.** Die erste
Fassung schrieb *„von sechs Ausschluss-Abschnitten nennen vier einen
Folge-Slice"*. Gezählt hatte sie **Slice-Kennungen** im Abschnitt, nicht
**Folge-Slice-Verweise**. Richtig: **acht** geschlossene Slices führen den
Abschnitt, und **zwei** nennen darin einen echten Folge-Slice
(slice-202 → slice-203, slice-207 → slice-208). Die beiden anderen Treffer
waren Rückverweise auf bereits geschlossene Slices. **Beide echten Adressen
nehmen die Sendung an** — slice-203 hat den verwiesenen Punkt geliefert,
slice-208 führt ihn als DoD (1). Kein Handlungsbedarf, und die Regel ist am
Bestand belegt statt behauptet.

**R4 verlangt eine Handlung, die die erste Fassung übersehen hat.** Der Kanon
bindet die Regel an **Schritt 4** des Acht-Schritt-Workflows. Dieses Repo führt
den Workflow als adoptierte Kopie in [`AGENTS.md`](../../../../AGENTS.md) §6
und in [`harness/README.md`](../../../../harness/README.md) — und sein
Schritt 4 lautet vollständig *„Kleinste sinnvolle Änderung planen."*, ohne
Out-of-Scope und ohne die Plan-Änderungs-Pflicht. Beide Träger gehören
nachgezogen.

**R5 — vier Setzungen, drei ohne Zutun erfüllt.** Der Kanon sagt: (a) die RTM
*wird erzeugt, nicht gepflegt* — `--trace` tut genau das; (b) *Bericht und Gate
sind derselbe Lauf* — `--trace` gegen `--trace --require-complete`, dieselbe
Mechanik; (c) *der Vorschlag des Kurses: der **Slice** entlastet, die ADR steht
als Spalte* — **genau so gesetzt**, und die erste Fassung dieses Audits
behauptete das Gegenteil: sie las `adrs:` als entlastende Quelle. Das
Lastenheft sagt wörtlich, *„eine bloße ADR-Referenz ohne Slice/Coverage deckt
weiterhin **nicht** ab"*
([`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)).
Die abgeleitete „Handlung" zielte auf eine Konfiguration, die es nicht gibt.

**Offen ist die vierte:** Der Kanon verlangt, einen **anderen** Schnitt zu
deklarieren, *„wie jede Abweichung von der Baseline"* — und nennt als Beispiel
genau unseren Fall, eine kuratierte Nachweis-Datei als entlastende Quelle. Das
ist [`DC-FA-COV-001`](../../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
(`trace.coverage`), eine **dritte, opt-in** Referenzklasse. Sie ist im
Lastenheft beschrieben, aber **nicht als Abweichung vom Kanon-Vorschlag
deklariert** — der Kanon kannte diesen Vorschlag bis `v6.5.0` nicht.

**R6 ist eine Regel in drei Formen, nicht zwei Regeln.** Der Kanon schreibt sie
so: Gate-Token statt Sensor-Link · `slice-NNN` statt Lifecycle-Pfad ·
Baseline-Stelle als Tag + Pfad in Inline-Code. Und er **erweitert die Klasse**:
einfrierend sind jetzt auch Archiv-Stub, `Accepted`-ADR und geschlossener
Slice. Am Bestand gemessen: **26** Verweise aus lebenden Artefakten verlinken
die Sensor-Datei (richtig), aus dem eingefrorenen Bestand tut es **einer** —
und der ist ein **Zitat** aus einem Review-Report, keine eigene Referenz. Kein
Handlungsbedarf am Bestand.

**Der Träger-Satz der ersten Fassung war falsch.** Sie schrieb, dieses Repo
führe *„eine eigene Review-Form"* und der Reviewer-Skill sei der Träger. Beides
hält nicht: Der Skill nennt die **Baseline**-Vorlage als Ziel-Form, es gibt
keine lokale Kopie, und kein Konventions-Eintrag deklariert eine Abweichung —
eine „eigene Form" wäre nach
[`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage) eine
undeklarierte Abweichung. Der Skill **trägt die Zitier-Form nicht**; sie
gehört hinein.

**R7 ist die Regel, die die erste Fassung ganz übersehen hat — und sie trifft
uns am härtesten.** Der Kanon schließt: *„Steht die Adresse erst im
eingefrorenen Artefakt, bleiben zwei Wege: es doch anfassen — dann ist es kein
Zeitdokument mehr — oder ein Ausnahme-Ventil im Prüfbereich, also eine
**Gate-Senkung mit eigener Begründungslast**."* Genau dieses Ventil betreibt
dieses Repo: `ignore-refs` in [`.d-check.yml`](../../../../.d-check.yml) trägt
**25** Tombstone-Einträge über **zehn** entfernte Baseline-Bäume, zuletzt für
`v6.3.1` mit [`MR-067`](../../../../harness/conventions.md#mr-067). Dahinter
stehen **28** Dateien mit Links in Bäume, die es nicht mehr gibt — 18
`Accepted`-ADRs, 6 aufgelöste Konventions-Einträge, 3 `done/`-Slices, ein CR.

**Das Ventil ist nicht falsch, aber es wächst mit jedem Bump**, und der Kanon
nennt es jetzt eine Gate-Senkung. Die Begründungslast ist damit fällig — und
R6 ist ihre Auflösung nach vorn: Wer die Kennung statt der Adresse schreibt,
braucht beim nächsten Bump keinen neuen Eintrag.

**Keine Regel ist *nicht anwendbar*, keine wird *abweichend* adoptiert.** Zum
Vergleich mit den Vorgängern, gemessen statt behauptet: `slice-107` führte
einen Stufen-Audit über sechs Stufen **mit** mehreren Nicht-anwendbar-Antworten;
`slice-203` führte **gar keinen** Regel-Audit, sondern übernahm Template-Deltas
direkt. Der Vergleich der ersten Fassung („beide Vorgänger") traf also nur auf
einen zu.

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
- [x] **(2)** Die Slice-Haus-Form ist aufgelöst: die Vorlage folgt der
      Baseline-Form, und die Regeln in
      [`.d-check.closure.yml`](../../../../.d-check.closure.yml) tragen den
      Bestand **und** die neue Form — mit einem **Bruch-Test je Richtung**, der
      belegt, dass beide noch gefangen werden.
- [x] **(3)** Die Konventions-Einträge sind nachgezogen: was der Kanon jetzt
      selbst sagt, ist aufgelöst; was abweicht, ist deklariert.
- [x] `make gates` grün.
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
  (1× — der in slice-207 angekündigte zweite Beleg wurde dort nie geschrieben; er kommt mit diesem Slice) — die Delta-Liste aus dem Vorgänger ist eine **Zählung**, und dieser
  Slice urteilt auf ihr. Vor jeder Aussage „so viele Regeln sind betroffen"
  gehört die Form des Gegenstands ausgeschrieben.
- [`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md)
  (7×) — für die Gegenrichtung: Wo der Delta eine Regel **entbehrlich** macht,
  ist der vorhandene Konventions-Eintrag aufzulösen, nicht umzuformulieren. Ein
  Eintrag, der nur noch wiederholt, was der Kanon selbst sagt, ist eine zweite
  Quelle.

- [`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
  (9×) — **nachgetragen nach Review-Runde 1**, und er war der einschlägigste
  von allen: Drei Befunde dieser Runde fallen in seine Klasse. Der Audit
  zählte Slice-**Kennungen** und sagte über Folge-Slice-**Verweise** aus; er
  las `adrs:` als entlastende Quelle und sagte über die
  Waisen-Definition aus; und er verglich mit „beiden Vorgängern", von denen
  einer gar keinen Audit führt. Der Test des Eintrags — *wer ändert die Menge,
  die ich zähle, und wer die, über die ich rede?* — hätte alle drei gefangen.

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
  (`zaehlmethode-misst-proxy-statt-gegenstand`, 1×). Das ist die höchste
  Risiko-Dichte, die ein Slice dieses Repos bisher vorab getragen hat.
- **Reconciliation-Aufwand:** keiner (GF). Graduation entfällt.

## 9. Closure-Notiz (nach `done/`)

\<wird vor dem `git mv` nach `done/` gefüllt\>
