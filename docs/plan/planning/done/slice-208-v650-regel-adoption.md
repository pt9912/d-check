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

**Der Bestand ist gemessen und klein — aber er ist eingefroren.** **Acht**
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
| R5 | Die **RTM** — vier Setzungen (s. u.) | `grundlagen-traceability` | `grundlagen-begriffe` | **vier erfüllt, keine Abweichung** |
| R6 | **Kennung statt Adresse** für einfrierende Artefakte, in **drei Formen**, für eine **erweiterte** Klasse | `grundlagen-harness-dateien` | `review-report`, `archiv-stub-slice`, `archiv-stub-welle`, `welle-results` | **übernommen, mit Handlung** (Skill) |
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

**Die vierte war eine Fehlmessung, und Review-Runde 2 hat sie gefunden.** Sie
lautete zunächst *„offen"*: Der Kanon verlangt, einen **anderen** Schnitt zu
deklarieren, *„wie jede Abweichung von der Baseline"*, und dieser Slice trug
dafür kurz einen eigenen Konventions-Eintrag. Er ist **zurückgezogen**. Der
Kanon-Auslöser gilt einem Repo, *„das das anders schneidet"* — die
`.d-check.yml` dieses Repos schneidet nicht anders: ihr `trace:`-Block führt
`requirements`, `adrs` und `slices` und **kein** `coverage`. Was abweicht, ist
eine **Produkt-Fähigkeit** ([`DC-FA-COV-001`](../../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in),
strikt opt-in, default-aus byte-identisch) — die eigene Konfiguration und der
Funktionsumfang des eigenen Werkzeugs sind zwei Mengen, und die Regel gilt
der ersten. Ein Eintrag, der eine nicht gelebte Abweichung deklariert, ist
schlechter als keiner.

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
undeklarierte Abweichung. Der Skill **trug die Zitier-Form nicht** — sie
ist mit diesem Slice hineingeschrieben: [`.harness/skills/reviewer.md`](../../../../.harness/skills/reviewer.md)
§Zitier-Form, Version 1.14.0. Damit sind es **drei** Handlungen an diesem
Delta, nicht zwei; die Antwort-Spalte der R6-Zeile sagte das zuerst nicht.

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
- **Ein Retrofit des `done/`-Bestands.** Die **acht** `done/`-Slices in Haus-Form bleiben,
  wie sie sind — eingefrorene Lauf-Belege.
- **Eine Umschrift der beiden lebenden Pläne.** Dieser Plan selbst und
  [slice-209](../in-progress/slice-209-frozen-klassen-vor-mechanischer-ersetzung.md)
  sind vor der Adoption geschrieben und tragen die Haus-Form. Sie werden
  **nicht** umnummeriert: Dieser hier trägt `d-check:cite`-Spannen und
  §-Verweise, die in Commit-Botschaften und in zwei Review-Reports zitiert
  sind — eine Umschrift mitten im Lauf machte jeden dieser Verweise falsch.
  Die Adoption ist **template-forward**: Sie gilt jedem Plan, der nach ihr
  entsteht, und beide Formen sind seit diesem Slice in den Deklarations-Trägern
  und im Prüf-Profil benannt. **Nachgetragen bei der Closure**, weil die
  Abgrenzung beim Anlegen implizit blieb — die Regel, die dieser Slice mit R4
  adoptiert, verlangt sie ausgeschrieben.
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
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [x] Jedes Risiko aus §5 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

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
  Review. — **Ausgang:** eingetreten — und zwar wie vorhergesehen. Beide
  Pflichten aus [`MR-066`](../../../../harness/conventions.md#mr-066) stehen in
  §6, geschrieben **vor** der ersten Zeile Arbeit: der Grund (eine Teilung
  zerrisse den Delta-Audit) und die Ersatz-Form (zwei Runden gegen je einen
  abgeschlossenen Stand). Vollzogen sind sie in den beiden Reports unter
  [`docs/reviews/`](../../../reviews/) — Runde 1 gegen den Delta-Audit, ohne
  die Migration zu sehen, Runde 2 gegen die Migration, den Audit als gegeben.
  **Damit ist dies die vierte Instanz**, an der sich der Eintrag messen lassen
  wollte, und die erste mit Vorab-Deklaration. Sie trägt: Runde 1 fand vier
  Lücken im Audit, die eine zweite Runde über denselben Bereich nicht gefunden
  hätte, weil sie ihn schon als geprüft gelesen hätte.
- **Der Umfang steht erst nach slice-207 fest.** `v6.4.0` ist angekündigt und
  bekannt, `v6.5.0` nicht. Trägt es eine eigene Regel-Änderung, ist dieser
  Slice **vor** der Beanspruchung neu zu schneiden — eine Schätzung jetzt wäre
  aus dem Anlass gezogen und nicht aus dem Bestand
  ([`BEO-ALL/rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md),
  7×). — **Ausgang:** eingetreten — `v6.5.0` trug **zwei** eigene
  Regel-Änderungen (RTM als Kanon-Begriff, Zitier-Form für einfrierende
  Artefakte), und der Slice wurde **nicht** neu geschnitten. Das ist der Fall,
  den der Risiko-Text als Alternative benennt, und er ist bewusst gegen den
  Schnitt entschieden: Beide Regeln hängen an derselben Delta-Liste wie die
  Form-Migration, ein Schnitt hätte sie zerrissen. Die Deckung dafür liefert
  der Ausgang darüber — [`MR-066`](../../../../harness/conventions.md#mr-066)
  ist genau für diese Wahl geschrieben. Kein Folge-Slice; die beiden Regeln
  sind in diesem Slice beantwortet (R5, R6/R7).
- **Die Gate-Regeln sind die stille Stelle.** Eine Regel, die auf einen
  Haus-Form-Titel keilt, wird nach der Umbenennung **grün, ohne noch etwas zu
  prüfen** — dieselbe Klasse wie ein Sensor, der auf `done/*.md` keilt und die
  archivierten Stubs nicht mehr sieht. Deshalb der Bruch-Test je Richtung in
  DoD (2) und nicht bloß ein grüner Lauf. — **Ausgang:** eingetreten — die
  stille Stelle war real und wurde gefunden: Genau **eine** Regel keilte auf
  einen wörtlichen Titel und wäre nach der Umbenennung grün geblieben, ohne
  noch etwas zu prüfen. Der Bruch-Test je Richtung ist gefahren und im Report
  von Runde 2 unabhängig reproduziert — beide Titel melden `section-forbidden`
  auf derselben Zeile, ein dritter Titel meldet weiterhin `section-missing`.
  **Der Review hat die Gegenrichtung mitgemessen und eine neue Fläche
  benannt** (F-15): Trägt eine Datei beide Titel, ist der Selektor mehrdeutig
  und die Prüfung entfällt für sie — laut, nicht still, und jetzt im
  Kommentar über der Regel.
- **Der `done/`-Bestand und die neue Form leben nebeneinander.** Solange beide
  existieren, ist „die Form eines Slice" zweideutig, und jede Regel darüber
  braucht eine Zeit- oder Verzeichnis-Grenze. Ob die acht Bestands-Slices als
  feste Liste oder über eine Ziffern-Schwelle ausgenommen werden, ist ein
  Entscheid — die Ziffern-Schwelle ist die gelebte Form
  ([`MR-056`](../../../../harness/conventions.md#mr-056)). — **Ausgang:**
  eingetreten, und der Entscheid ist gefallen: **beides, an getrennten
  Achsen**. Die Ziffern-Schwelle (`exempt-paths` bis `slice-13?`) trennt den
  Altbestand, der den Abschnitt gar nicht führt; das **Doppel-Muster** im
  Selektor trägt die acht Slices, die ihn in Haus-Form führen, neben der neuen
  Form. Eine feste Dateiliste wäre die dritte Möglichkeit gewesen und ist
  verworfen: Sie müsste bei jedem geschlossenen Slice gepflegt werden, und
  genau das ist die Klasse, die dieses Repo als zweite Quelle meidet.

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
  den Zuschnitt sprengen. — **Ausgang:** eingetreten — beide sind Urteile
  geworden, und beide haben den Zuschnitt gehalten. **(a)** Die RTM: vier
  Setzungen, alle vier ohne Zutun erfüllt. Der kurzzeitig geschriebene
  Konventions-Eintrag dazu ist **zurückgezogen** — er deklarierte eine
  Abweichung, die dieses Repo nicht praktiziert. **(b)** Die Unterscheidung
  *einfrierend / lebend*: übernommen in zwei Trägern — die Zitier-Form steht
  im Reviewer-Skill (1.14.0), und das `ignore-refs`-Ventil ist als
  Gate-Senkung deklariert
  ([`MR-069`](../../../../harness/conventions.md#mr-069)). Was sie **nicht**
  entbehrlich macht, ist der ADR-Tombstone: Die ADR-Vorlage trägt die
  Zitier-Form nicht, und für die 18 `Accepted`-ADRs gibt es damit keinen
  Träger — als benannte Grenze in [`MR-069`](../../../../harness/conventions.md#mr-069)
  festgehalten, nicht stillschweigend
  weggelassen.
- **Die Umnummerierung bewegt JEDEN Abschnitt** — gemessen: neun Haus-Form-
  Abschnitte gegen acht der Baseline, und die Zuordnung ist keine Bijektion
  (§1+§3 fallen zusammen, §6 spaltet sich in §4+§5). **Gate-seitig ist es
  weniger, als es aussieht:** **Fünf** `structure`-Regeln des Closure-Profils
  adressieren einen Abschnitt der `done/`-Slices; genau **eine** keilte auf
  einen wörtlichen Titel (`## 5. Abnahme-Punkte / Risiken`). Die vier übrigen
  sind form-agnostisch — zweimal `Definition of Done` mit offener Ziffer, dazu
  `Closure-Notiz` und jedes H1 — und überleben die Umbenennung unverändert.
  **Die erste Fassung zählte drei**, und der unabhängige Review hat
  nachgemessen: Sie zählte die Regeln, die beim Schreiben vor Augen standen,
  und sagte über alle aus. Der Schluss bleibt richtig, die Zahl unter ihm war
  falsch — genau die Klasse, die dieser Slice in §7 als gesichtet führt.
  — **Ausgang:** eingetreten — und zwar in der Form, die der Risiko-Text
  vorwegnimmt: *„sie ist eine Zählung, und der Register-Eintrag zu Zählungen
  ist gesichtet"*. Die Zählung war zweimal falsch — „drei Regeln" statt fünf,
  und beim Korrigieren „16 Ventil-Einträge" statt sieben von 25 —, beide Male
  weil über eine andere Menge geredet als gezählt wurde. Eingetragen als
  achter Beleg bei
  [`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md),
  die Proxy-Variante als dritter bei
  [`zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md).
  Das Sichten hat die Klasse benannt und ihr Eintreten nicht verhindert; der
  unabhängige Review hat sie gefangen.

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
  **acht** `done/`-Slices in Haus-Form sind eingefrorene Lauf-Belege und werden
  nicht angefasst.** Der Folge-Slice
  [slice-209](../in-progress/slice-209-frozen-klassen-vor-mechanischer-ersetzung.md)
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

**Geliefert.** Ein Delta-Audit über **sieben** Regeln in zwölf Trägern, je
Regel eine Antwort; die Auflösung der Slice-Haus-Form (Vorlage, Prüf-Profil,
Deklarations-Träger) mit einem Bruch-Test je Richtung; drei Handlungen am
Delta — der Out-of-Scope-Schritt in beiden Workflow-Kopien, die Zitier-Form im
Reviewer-Skill (1.14.0), das `ignore-refs`-Ventil als deklarierte Gate-Senkung
([`MR-069`](../../../../harness/conventions.md#mr-069)). Zwei Review-Runden
gegen je einen abgeschlossenen Stand, beide blockierend, beide eingearbeitet.
`make gates` grün (zehn Gates, 693 Dateien), `make verify-closure-notes` grün
(593 Dateien).

**Was funktioniert hat: die vorab deklarierte Ersatz-Form.** Dies ist die
vierte Instanz von [`MR-066`](../../../../harness/conventions.md#mr-066) und
die erste, in der Grund und Ersatz-Form **vor** der Arbeit standen. Der Gewinn
ist messbar und nicht rhetorisch: Runde 1 sah nur den Delta-Audit und fand
darin vier Lücken — eine Regel ganz ohne Antwort, zwei, die in Wahrheit eine
waren, vier fehlende Träger. Eine zweite Runde über denselben Gesamt-Bereich
hätte den Audit als geprüft gelesen. Die Teilung **nach Gegenstand** statt nach
Umfang ist der Grund, dass die Lücke gefunden wurde.

**Was Friktion war: die Deklaration hinkt der Konfiguration hinterher.** Die
Gate-Regel wurde auf beide Abschnitts-Formen umgestellt und mit einem
Bruch-Test belegt — und sechs lebende Stellen, die über genau diese Regel
**reden**, nannten weiter die alte Nummer. Der Slice hatte
[`semantic-change-body-only-edges-stale`](../observations/BEO-ALL/semantic-change-body-only-edges-stale/observation.md)
in §7 selbst als *„den zentralen"* Eintrag geführt und sich dessen Ableiter
verordnet; angewandt wurde er auf die Ausführung, nicht auf die Rede darüber.
Ein Spiegel ist nicht nur, wer die Regel **ausführt**, sondern auch, wer sie
**zitiert**.

**Steering-Loop-Lerneintrag — neu benannt:** *Eine Review-Korrektur wird nur an
der zitierten Fundstelle eingearbeitet*
([`review-fix-applied-only-at-cited-site`](../observations/BEO-ALL/review-fix-applied-only-at-cited-site/observation.md),
1×). Review-Runde 2 hat das Muster selbst benannt: drei ihrer Befunde sind
Stellen, an denen eine Korrektur aus Runde 1 an genau der zitierten Fundstelle
landete und ihre Geschwister stehen blieben — im deutlichsten Fall blieb die
widerlegte Zahl in §3 stehen, also dort, wo sie den Umfang **normativ**
festlegt. Das `pfad`-Feld eines Findings sagt, wo der Befund **gefunden**
wurde, nicht wo er **steht**; ein Reviewer belegt, er inventarisiert nicht.
Der Eintrag steht bei 1× und wartet, wie es die Regel vorsieht.

**Drei weitere Belege**, alle drei in §7 vorab gesichtet und trotzdem
eingetreten:
[`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
(dreimal, dazu derselbe Fehler beim Korrigieren),
[`zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
(Slice-Kennungen gezählt, über Folge-Slice-Verweise ausgesagt) und
[`semantic-change-body-only-edges-stale`](../observations/BEO-ALL/semantic-change-body-only-edges-stale/observation.md)
(oben). **Das Sichten hat sie benannt und ihr Eintreten nicht verhindert** —
gefangen hat sie der unabhängige Review, jedes Mal. Wer daraus schließt, das
Sichten sei wirkungslos, liest zu viel hinein: Ohne die Sichtung hätte der
Slice die Klassen nicht benennen können, in die seine Befunde fallen. Aber als
*Vermeidungs*-Werkzeug hat es in diesem Lauf nicht getragen, und das gehört
so notiert.

**Ein zurückgezogener Konventions-Eintrag.** `MR-068` <!-- d-check:ignore (zurückgezogen, es gibt kein Ziel mehr) --> deklarierte eine
Abweichung über `trace.coverage` — einen Schlüssel, den die `.d-check.yml`
dieses Repos nicht führt. Er ist gelöscht, nicht umformuliert: Ein Eintrag, der
eine nicht gelebte Abweichung behauptet, ist schlechter als keiner. Die
Nummer bleibt vergeben und wird nicht nachbelegt
([`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage)).

**Die drei Paarungen, gemessen.** **(a) Anker** — vakant: Der
Steering-Loop-Eintrag dieses Slice steht bei 1× und ist damit *gezählt, nicht
verkörpert*; er trägt kein `liegt in`-Pflichtfeld und ist kein Gegenstand der
Paarung. Was dieser Slice verkörpert ([`MR-069`](../../../../harness/conventions.md#mr-069),
die Zitier-Form im Skill,
Schritt 4), stammt aus dem **Kanon-Delta**, nicht aus der 3×-Schwelle, und
braucht deshalb keinen Herkunfts-Anker. **(b) Folge-Slice** — der einzige
genannte ist
[slice-209](../in-progress/slice-209-frozen-klassen-vor-mechanischer-ersetzung.md),
und er liegt im Lifecycle (`open/`). **(c) Register** — alle in diesem Plan
zitierten Beobachtungs-Pfade lösen auf; von den jetzt **37** Verzeichnissen
trägt genau eines kein `evidence/`:
[`kanal-kennung-als-inhalt-gelesen`](../observations/BEO-ALL/kanal-kennung-als-inhalt-gelesen/observation.md).
Das ist **kein** Defekt dieses Slice, sondern ein bewusst als *Wachposten*
geführter Eintrag mit Zähler 0 — der Kanon lässt ein Vorkommen ohne
abgeschlossenen Vorgang ausdrücklich zu („benannt, nicht gezählt"). Die
zweite Hälfte der Register-Paarung („jede Zeile trägt mindestens einen
Beleg") und diese Form widersprechen einander; hier nur **gemessen und
benannt**, nicht aufgelöst — die Auflösung wäre eine Änderung an der
Paarungs-Regel und gehört nicht in einen Adoptions-Slice.

**Was offen bleibt, benannt statt verschwiegen.**
[`MR-069`](../../../../harness/conventions.md#mr-069) trägt drei Grenzen,
und die dritte ist die unbequeme: Die Auflösung nach vorn — Kennung statt
Adresse — hat für Review-Reports einen Träger (den Skill), für die **18
`Accepted`-ADRs** keinen. Eine ADR entsteht aus einer Vorlage, und die
adoptierte Vorlage führt die Zitier-Form nicht. Solange das so ist, erzeugt die
gelebte ADR-Praxis weiter Ventil-Einträge. Der Auflösungs-Trigger dieses Eintrags
verlangt deshalb **beide** Bedingungen und kann heute nicht feuern.
