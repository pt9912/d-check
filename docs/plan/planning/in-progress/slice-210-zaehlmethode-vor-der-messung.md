# Slice slice-210: Die Form des Gegenstands steht vor seiner Zählung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
(3× erreicht, Ausgang *geplant*),
[`BEO-ALL/eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
(der Geschwister-Eintrag: falsche **Menge**, richtige Kategorie).

**Berührte Spec-Stellen:** — *(keine; der Slice verkörpert eine
Steering-Loop-Regel, er ändert keine Anforderung)*

**Verantwortlich:** pt9912 · **Autor:** pt9912. **Datum:** 2026-09-07.

**Form-Hinweis.** Erster Plan dieses Repos in der **Baseline-Form**
(`v6.5.0`, acht Abschnitte). Die Haus-Form ist mit
[slice-208](../done/slice-208-v650-regel-adoption.md) aufgelöst —
template-forward, also ab diesem Plan.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Die Regel schreiben, die drei Vorgänge verlangt haben: **Vor einer
Messung steht die Form ihres Gegenstands** — was macht eine Kennung zu einer
Kennung, einen Verweis zu einem Folge-Verweis, eine Regel zu einer Regel? —
**und die Trefferliste wird stichprobenweise gegen diese Form gehalten, nicht
nur ihre Zahl gelesen.**

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Geschwister-Eintrag `eigene-menge-gemessen-fremde-behauptet`.** Er
  steht bei 10×, sein Stand ist *gemischt* (differenziert nach Instanz, zwei
  davon ohne formgültigen Ausgang); seine Frage ist die **Menge**
  (*wer ändert, was ich zähle?*), diese hier ist die **Kategorie** (*ist das
  Gezählte überhaupt der Gegenstand?*). Zwei Fragen, zwei Regeln — eine
  gemeinsame Formulierung machte beide unschärfer.
- **Ein Sensor.** Ob jemand vor einer Messung ihre Form ausgeschrieben hat, ist
  eine Aussage über einen **Akt**, nicht über einen ruhenden Zustand — dieselbe
  Lage wie bei [`AGENTS.md`](../../../../AGENTS.md) §3.6. Es gäbe nichts zu
  bauen, und ein Gate zu behaupten wäre schlimmer als keines.
- **Ein Retrofit der drei Vorgänge.** slice-205, slice-207 und slice-208 liegen
  in `done/` und sind eingefrorene Lauf-Belege; ihre Fehlmessungen stehen in den
  Evidence-Dateien und bleiben dort.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [x] **(1)** Der **Zielort ist entschieden und begründet**: trägt die Regel
      [`AGENTS.md`](../../../../AGENTS.md) §5 (wie die beiden
      Geschwister-Klassen), der Reviewer-Skill, oder beide? Ein Ort, der aus
      der Gewohnheit gewählt ist, ist keiner.
- [x] **(2)** Die Regel steht am entschiedenen Ort, mit Herkunfts-Anker
      `seit slice-210` und mit der **Grenze**, was sie nicht leistet.
- [x] **(3)** Der Registereintrag trägt den Ausgang *verkörpert* mit
      auflösbarem Zielort.
- [x] `make gates` grün.
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| der in DoD (1) entschiedene Zielort | update | die Regel selbst |
| [`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/state.md) | update | Ausgang von *geplant* auf *verkörpert*, mit Zielort und Anker |

**Die drei Vorgänge, und was sie gemeinsam haben.** In slice-205 wurden
`###`-Überschriften als Anforderungs-Kennungen gezählt — deutsche Komposita
zählten mit, und daraus entstand ein fünftes ID-Schema, das es nicht gibt. In
slice-207 wurden geänderte **Zeilen** gezählt und über geänderte **Regeln**
ausgesagt. In slice-208 wurden **Slice-Kennungen** gezählt und über
**Folge-Slice-Verweise** ausgesagt. Dreimal ähnelte das Muster dem Gegenstand
genug, um eine plausible Zahl zu liefern, und dreimal wurde aus den
Fehltreffern eine Aussage abgeleitet, die es nie gab.

**Dreimal von vier fand es der unabhängige Review, nicht der Zählende** — das
ist die Eigenschaft, die den Zielort mitbestimmt: Eine Regel, die nur im
Briefing steht, adressiert genau den, der sie beim Zählen nicht liest.

### Der Zielort, entschieden (DoD 1)

**Beide — [`AGENTS.md`](../../../../AGENTS.md) §5 und der Reviewer-Skill.** Die
Entscheidung hat zwei Teile, und der erste ist die Frage, ob es überhaupt eine
Lücke gibt.

**Teil 1: Deckt `AGENTS.md` §5 die Klasse schon?** Nein — und das ist
**nicht** meine Lesart, sondern die Selbstauskunft des Trägers. Der §5-Absatz
sagt, ein Schluss reiche *„nicht weiter als die gemessene Menge"*; der
zugehörige Anker 8 des Reviewer-Skills schreibt dazu wörtlich:

> diese ist MEDIUM, weil die **Messung stimmt** und nur ihre Reichweite
> überdehnt ist.

**Beide setzen eine korrekte Messung voraus.** Der Fall, in dem die Messung
selbst den falschen Gegenstand zählt, liegt außerhalb — nach ihren eigenen
Worten, nicht nach meiner Auslegung. Das ist die Vorsichtsmaßnahme gegen
[`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
(16×, zweimal in Folge eingetreten): Ich lese §5 nicht enger, um meine Regel zu
retten — der Träger von §5 nennt seinen Geltungsbereich selbst.

**Teil 2: Warum beide Orte, und warum das keine Gewohnheit ist.** Die
Geschwister-Einträge stehen ebenfalls an beiden Orten; genau deshalb ist
*„beide"* die verdächtige Antwort und braucht ein eigenes Argument.

- **Zwei verschiedene Handlungen zu zwei verschiedenen Zeitpunkten.**
  `AGENTS.md` adressiert den **Messenden**, *bevor* er zählt: die Form des
  Gegenstands ausschreiben. Der Skill adressiert den **Prüfenden**, *nachdem*
  gezählt wurde: die Trefferliste gegen diese Form halten. Keine der beiden
  Anweisungen lässt sich am anderen Ort ausführen.
- **Gemessen, wo die Klasse tatsächlich auffällt:** Von den **vier** Instanzen
  in drei Vorgängen fand sie **dreimal** der unabhängige Review und **einmal**
  der Zählende selbst. Eine Regel nur in `AGENTS.md` stünde damit an der
  Stelle, an der sie nachweislich seltener greift.

**Und der Test aus
[`begruendung-traegt-entscheidung-nicht`](../observations/BEO-ALL/begruendung-traegt-entscheidung-nicht/observation.md)
(2×), angewandt:** *Bliebe die Entscheidung richtig, wenn die Begründung falsch
wäre?* **Ja.** Das zweite Argument ist eine Zahl, und sie stärkt die
Skill-Hälfte, sie erzeugt sie nicht — die Entscheidung ruht auf dem ersten
Argument, den zwei Handlungen. Das ist hier keine Formalie: Die Zahl **war**
zuerst falsch. Die `state.md` dieses Eintrags nannte *„viermal von fünf"*;
nachgezählt sind es **dreimal von vier**, und die fünfte Instanz gehörte nie
hierher, sondern zum Geschwister-Eintrag. **Der Slice über Fehlzählungen hat
sich beim Zählen seiner eigenen Belege vertan** — korrigiert, und als Beleg für
genau diese Klasse zu behandeln.

## 4. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei. Kein inhaltlicher
Vorgänger — die Schwelle ist mit slice-208 erreicht, der Eintrag trägt den
Ausgang *geplant* mit dieser Kennung.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Zeigt DoD (1), dass die Regel an **zwei**
  Orten stehen muss und der zweite eine eigene Form verlangt (etwa ein
  Skill-Anker mit eigener Kategorie), wird der zweite Ort ein eigener Slice.
- `in-progress` → `open` (blockiert): Ergibt DoD (1), dass die Regel von
  [`AGENTS.md`](../../../../AGENTS.md) §5 bereits mitgetragen wird — der
  Absatz über die *gemessene Menge* liegt nah —, ruht der Slice bis zum
  Entscheid, ob eine Schärfung dort die Klasse deckt oder verwässert.

## 5. Closure-Trigger

Zwei beobachtbare Kriterien und ein Lerneintrag: (a) die Regel steht am
entschiedenen Ort und `make gates` ist grün; (b) der Registereintrag trägt
*verkörpert* mit einem Zielort, der auflöst, und dem Anker `seit slice-210`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Die Regel könnte in `AGENTS.md` §5 schon stecken.** Dort steht, der Schluss
  reiche *„nicht weiter als die gemessene Menge"*. Das ist nah, aber es ist die
  Aussage über den **Schluss**, nicht über die **Messmethode** — der Unterschied
  ist genau der zwischen `commit-message-overclaims-work` und diesem Eintrag.
  Ob er trägt, ist DoD (1). — **Ausgang:** entfallen — DoD (1) hat es
  entschieden, und zwar nicht durch Auslegung: Anker 8 des Reviewer-Skills
  schreibt die korrekte Messung wörtlich als seine **Voraussetzung** fest
  (*„weil die Messung stimmt"*). Der Fall, in dem die Messung selbst den
  falschen Gegenstand zählt, liegt damit nach der Selbstauskunft des Trägers
  außerhalb. **Der unabhängige Review hat genau diese Frage nachgeprüft** — sie
  stand als erste in seinem Auftrag — und die Lücke bestätigt. Das Risiko kann
  nicht mehr eintreten.
- **Eine Regel aus drei Vorgängen ist keine Inventur**
  ([`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md),
  8×, zuletzt slice-209). Alle drei Vorgänge sind Doku-/Planning-Messungen; ob die Regel für eine
  Code-Messung trägt, ist unbelegt und gehört als Grenze in die Formulierung.
  — **Ausgang:** weiter offen — eingetragen als neunter Beleg bei
  [`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md),
  zweiter in Folge nach slice-209. Die Grenze steht in
  [`AGENTS.md`](../../../../AGENTS.md) §5, wie der Plan es verlangt hat: belegt
  ist die Klasse nur an Doku-/Planning-Messungen, für eine Code-Messung ist sie
  unbelegt. Zweimal hintereinander eine Regel aus dem Anlass statt aus dem
  Bestand zu ziehen, ist selbst ein Muster — es steht in der Closure-Notiz und
  nicht nur hier.
- **Kein Gate, und das ist eine Aussage über die Wirkung.** Die Regel verschiebt
  einen Fehler von *unsichtbar* nach *vermeidbar*, nicht nach *unmöglich*; drei
  von vier Instanzen fand ein Review, und daran ändert eine geschriebene Regel
  zunächst nichts. Das gehört ausgeschrieben, sonst liest sie sich stärker, als
  sie ist. — **Ausgang:** entfallen — die Grenze ist ausgeschrieben (*„gilt der
  Methode, nicht der Sorgfalt — sie verhindert keinen Zählfehler, sie macht ihn
  auffindbar"*), und der Skill-Anker nennt zusätzlich, **wo** die Regel greift:
  in drei von vier belegten Instanzen fand die Klasse der Review, nicht der
  Zählende. **Der Review hat die Grenzen-Liste trotzdem als unvollständig
  erwiesen** (F-4): Die Regel nannte keinen Ort, an dem die Form stehen muss,
  und war damit nicht befolgbar — anders als beide Nachbar-Regeln. Behoben; die
  Form steht jetzt dort, wo die Zahl steht. **Zweiter Slice in Folge, dessen
  ausgeschriebene Grenzen die eigene Lücke nicht enthielten.**

## 7. Closure-Notiz

**Geliefert.** Die Regel *„Vor einer Messung steht die Form ihres
Gegenstands"* an zwei Orten — [`AGENTS.md`](../../../../AGENTS.md) §5 für den
Messenden, [`.harness/skills/reviewer.md`](../../../../.harness/skills/reviewer.md)
Anker 17 für den Prüfenden (Version 1.15.0) —, beide mit Herkunfts-Anker
`seit slice-210` und ausgeschriebener Grenze. Registereintrag auf
*verkörpert*. Ein unabhängiger Review, blockierend, vier MEDIUM plus zwei LOW
und ein INFO — alle eingearbeitet bis auf das INFO, das benannt bleibt.
`make gates` grün (zehn Gates, 711 Dateien).

**Was funktioniert hat: die Lücke wurde nicht behauptet, sondern belegt.** DoD
(1) hätte sich bequem mit *„§5 deckt das nicht"* beantworten lassen. Statt
dessen steht dort das Zitat aus Anker 8 — *„weil die **Messung stimmt** und nur
ihre Reichweite überdehnt ist"* —, also die Selbstauskunft des Trägers über
seinen eigenen Geltungsbereich. Das war die Vorsichtsmaßnahme gegen
[`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md),
und der Review hat sie eigens nachgeprüft: Die Lücke ist echt, §5 wurde nicht
enger gelesen, um die Regel zu retten. **Ein Beleg aus dem Zielartefakt selbst
ist stärker als jede Auslegung** — und er war billiger zu haben, als die
Diskussion darüber gewesen wäre.

**Was Friktion war: der Slice über Fehlzählungen hat dreimal falsch gezählt.**
Einmal selbst gefunden (*„viermal von fünf"* statt *„dreimal von vier"*, plus
eine Instanz, die dem Geschwister-Eintrag gehört), zweimal vom Review
(`10×, verkörpert` für einen Eintrag, dessen Stand *gemischt* ist und der zwei
Instanzen ohne formgültigen Ausgang führt; `7×` gegen `8×` für denselben
Eintrag im selben Plan). **Alle drei betreffen Zahlen über das
Beobachtungs-Register** — also über das Artefakt, das dieser Slice am
gründlichsten gelesen hat.

**Der Steering-Loop-Lerneintrag ist deshalb nicht die Regel, sondern ihre
Reichweite.** Die geschriebene Regel hätte keinen der drei Fehler verhindert:
Sie verlangt die **Form** des Gezählten, und die war in allen drei Fällen klar
— falsch war die **Menge** und der abgelesene **Zustand**. Das ist die
Arbeitsteilung mit
[`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
(jetzt 11×) und
[`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
(17×), und sie ist mit diesem Lauf gemessen statt vermutet: **Drei
Register-Klassen liegen nebeneinander, und die neue deckt die schmalste von
ihnen.** Wer sie für die allgemeine Zahlen-Disziplin hält, liest zu viel
hinein — der Satz gehört in die Regel und steht dort.

**Zwei Muster über Slice-Grenzen hinweg, beide zum zweiten Mal in Folge.**
**(1)** Eine Regel aus dem **Anlass** statt aus dem Bestand: slice-209 und
slice-210 haben beide die 3×-Schwelle bedient, beide die Inventur nicht
gemacht, beide die Grenze vorab benannt
([`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md),
neunter Beleg). **(2)** Eine ausgeschriebene **Grenzen-Liste**, die die eigene
Lücke nicht enthielt: In slice-209 fand der Review eine sechste Grenze, hier
(F-4) die fehlende Ortsangabe, ohne die die Regel nicht befolgbar war. Eine
Grenzen-Liste macht eine Regel ehrlicher, nicht vollständig — und sie braucht
denselben fremden Leser wie die Regel selbst.

**Eine Kanon-Regel wurde nicht befolgt, und das ist kein Steering-Loop-Signal.**
Die `state.md` speicherte *„Drei Vorgänge, vier Instanzen"*, obwohl der Kanon
sagt, der Zähler werde **abgeleitet**: *„Es gibt kein Feld, in das man ihn
schreibt, und deshalb keines, das falsch stehen kann."* Der Review hat die
Drift **gemessen** — dieselbe Evidence-Basis trug nacheinander drei
verschiedene Zahlen. Das gehört nicht ins Register: Die Regel existiert, sie
wurde übergangen. Beide Zahlen sind entfernt.

**Die drei Paarungen, gemessen.** **(a) Anker** — beide Zielorte tragen
`seit slice-210`, beide lösen auf, der Registereintrag nennt beide.
**(b) Folge-Slice** — dieser Slice nennt keinen; nichts wurde auf später
verwiesen. **(c) Register** — alle zitierten Pfade lösen auf, die vier neuen
Belege liegen als `evidence/slice-210.md` in ihren Verzeichnissen. Der
Wachposten
[`kanal-kennung-als-inhalt-gelesen`](../observations/BEO-ALL/kanal-kennung-als-inhalt-gelesen/observation.md)
trägt weiterhin kein `evidence/` — unverändert die benannte Spannung aus
slice-208.
## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt selbst entfällt nie** —
die vorgelagerten Prüfungen hängen weder am Modus noch am Slice-Typ; bedingt
ist allein der Modus-Block am Ende. Dieses Repo führt **drei** Prüfungen: die
zwei kanonischen und, als Adaption, den Nachtlauf-Stand
([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:268-269 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area: `*` (Repo-Default). Der Slice schreibt eine Regel über das
Messen und ändern höchstens `AGENTS.md`, den Reviewer-Skill und eine
Register-Datei. `tools/harness/` ist **nicht** berührt — §1 schließt einen
Sensor ausdrücklich aus, und ohne Sensor gibt es dort nichts anzufassen.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:274-274 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, **37** Verzeichnisse über beide
Kürzel). **Sechs** Einträge sind einschlägig, und drei davon entscheiden über
den Zuschnitt:

- [`zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
  (3×, Ausgang *geplant*) — **der Gegenstand selbst.** Der Eintrag nennt
  diesen Slice namentlich; sein Ausgang wird hier zu *verkörpert*.
- [`commit-message-overclaims-work`](../observations/BEO-ALL/commit-message-overclaims-work/observation.md)
  (9×, verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §5) — **der Eintrag,
  der DoD (1) entscheidet.** Seine Regel sagt, ein Schluss reiche *„nicht
  weiter als die gemessene Menge"*. Das ist die Aussage über den **Schluss**;
  hier geht es um die **Messmethode**. Ob die Differenz trägt oder ob eine
  Schärfung dort genügt, ist genau die Frage von DoD (1) — und sie ist
  **nicht** vorentschieden.
- [`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
  (16×) — **zweimal in Folge eingetreten**, zuletzt in slice-209, und beide
  Male in derselben Bauform: eine vorhandene Regel wurde weiter gelesen, als
  ihr Geltungsbereich reicht. DoD (1) stellt genau diese Frage an
  `AGENTS.md` §5. Der Ableiter gilt hier ungekürzt — den **Absatz** lesen,
  nicht den Titel; und ein **Ja** („§5 deckt es schon"), das die dortige Regel
  weiter zieht, als sie reicht, wäre der teuerste Fehler dieses Slice.
- [`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
  (10×, Stand *gemischt* — differenziert nach Instanz, zwei davon ohne
  formgültigen Ausgang) — der **Geschwister-Eintrag**, den §1 ausdrücklich
  ausschließt. Die Abgrenzung ist beim Schreiben nachzuhalten: dort ist die
  **Menge** falsch und die Kategorie richtig, hier ist die **Kategorie**
  erfunden. Eine gemeinsame Formulierung machte beide unschärfer.
- [`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md)
  (8×, zuletzt slice-209) — §6 führt es als Risiko: drei Vorgänge, alle drei
  Doku-/Planning-Messungen. Die Grenze gehört in die Formulierung.
- [`begruendung-traegt-entscheidung-nicht`](../observations/BEO-ALL/begruendung-traegt-entscheidung-nicht/observation.md)
  (2×, frisch aus slice-209) — für DoD (1): Der Zielort kann **richtig**
  gewählt und **falsch** begründet sein, und dann leitet die nächste Änderung
  ihre Reichweite aus dem falschen Satz ab. Der Test des Eintrags gehört
  angewandt: *Bliebe die Entscheidung richtig, wenn die Begründung falsch
  wäre?*

**Keiner der sechs erreicht mit diesem Slice die Schwelle erstmalig.**

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-07 gelesen: **beide Nachtläufe grün** —
`upstream-drift.yml` (jüngster Lauf 2026-09-07T05:33:45Z) und `image-scan.yml`
(2026-09-07T08:21:32Z). Nichts zu tun.

**Modus-Begründungsblock.** Alle berührten Sub-Areas GF — ein Block genügt.

### Sub-Area: `*`

- **Modus:** GF (Repo-Default).
- **Konventions-Dichte:** hoch für die **Träger**, null für den
  **Gegenstand**. Wie eine Regel in [`AGENTS.md`](../../../../AGENTS.md) §5
  oder im Reviewer-Skill auszusehen hat, ist dicht geregelt; *dass* eine
  Messmethode vor ihrer Zahl steht, sagt heute kein Artefakt. Genau diese
  Asymmetrie ist der Slice.
- **Phase-Reife:** Phase 5 für den Vorgang — eine Steering-Loop-Regel aus
  einem 3×-Stand ist der eingespielteste Ablauf dieses Repos, zuletzt in
  slice-209. Phase 4 für den Gegenstand: Die Klasse ist benannt und dreimal
  belegt, aber noch nie in eine Formulierung gebracht worden.
- **Evidenz-/Diskrepanz-Risiko:** **niedrig für den Bestand, hoch für DoD (1).**
  Am Bestand ist nichts zu inventarisieren — die drei Vorgänge sind in den
  Evidence-Dateien belegt und werden nicht angefasst. Das Risiko sitzt allein
  im Urteil über `AGENTS.md` §5, und es ist der 16×-Eintrag
  `citation-stretched-beyond-scope`, der es benennt — zweimal in Folge
  eingetreten.
- **Reconciliation-Aufwand:** keiner (GF). Graduation entfällt.
