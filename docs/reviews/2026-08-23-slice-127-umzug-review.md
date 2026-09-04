# Review-Report: slice-127 (zweiter Anlauf) — zwei Waisen nach AGENTS.md umgezogen

**Review-Art:** Code- und Plan-Review (Slice-Plan · Hard Rules · gepinnter Kanon gegen den Diff), unabhängiger Reviewer ohne Anteil an der Arbeit
**Gegenstand:** `git diff 84b588e..b806b03` — die zwei Slice-Commits `753fb46` (beansprucht, Lifecycle-Move `next/` → `in-progress/`) und `b806b03` (zwei Zuzüge nach AGENTS.md, CLAUDE.md 23 → 4 Zeilen)
**Skill:** `.harness/skills/reviewer.md` @ 1.9.0 · **Modell-ID:** `claude-opus-5[1m]` · **Datum:** 2026-08-23

**Bewegliches Ziel — ausdrücklich festgehalten.** Zu Review-Beginn war `b806b03` HEAD und der Baum sauber. Während des Reviews sind ein dritter Commit `b3a74e4` („Herkunfts-Prosa aus beiden Zuzügen gestrichen") und weitere, nicht committete Änderungen an `AGENTS.md` und `harness/README.md` dazugekommen. Geprüft ist der beauftragte Stand `84b588e..b806b03`; wo ein Nachzügler einen Befund bereits erledigt, steht das beim Befund.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan [`slice-127-claude-md-pointer.md`](../plan/planning/done/slice-127-claude-md-pointer.md) §1, §2, §3, §4, §5, §7
- [`CLAUDE.md`](../../CLAUDE.md) vorher (`git show 84b588e:CLAUDE.md`) und nachher
- [`AGENTS.md`](../../AGENTS.md) §1, §2, §3, §5, §6 (Zeilen-Angaben gegen `git show b806b03:AGENTS.md`)
- Gepinnter Kanon: [`grundlagen-source-precedence.md`](../../.harness/baseline/v5.11.0/regelwerk/grundlagen-source-precedence.md) §Source Precedence/§Vollständigkeit, [`grundlagen-durchsetzungsschicht.md`](../../.harness/baseline/v5.11.0/regelwerk/grundlagen-durchsetzungsschicht.md) §Referenz-Implementierung, [`grundlagen-traceability.md`](../../.harness/baseline/v5.11.0/regelwerk/grundlagen-traceability.md) §Herkunfts-Anker, [`modul-09-implementierung.md`](../../.harness/baseline/v5.11.0/regelwerk/modul-09-implementierung.md) §Minimal Agent Workflow/§AGENTS.md-Regeln, [`modul-02-harness-bootstrap.md`](../../.harness/baseline/v5.11.0/regelwerk/modul-02-harness-bootstrap.md) §Freshness-Audit
- [`harness/README.md`](../../harness/README.md) §Leseordnung, §Source precedence · [`harness/conventions.md`](../../harness/conventions.md) §Modus-Deklaration, MR-000, MR-013, MR-015
- `.claude/commands/implement-slice.md`, `.claude/settings.json`, `.claude/hooks/*`
- Delta-Audit [`slice-129-baseline-v5110-delta-audit.md`](../plan/planning/done/welle-83/slice-129-baseline-v5110-delta-audit.md) §4a (Fünf-Fundorte-Zensus)
- Report des ersten Anlaufs [`2026-08-23-slice-127-claude-md-pointer-review.md`](2026-08-23-slice-127-claude-md-pointer-review.md)
- [`observations.md`](../plan/planning/observations.md) BEO-002, BEO-008, BEO-009, BEO-011

**Vom Reviewer selbst gefahren** (Ausgabe in Datei umgeleitet, Exit direkt aus `$?` — BEO-007):
`make gates` → **Exit 0**; `d-check: 461 Datei(en) geprüft, 0 Befund(e)` (dreimal: doc-check, gate-consistency, planning-check), `coverage-gate: OK — Coverage 94.80% erfüllt Schwelle 93%`, `Makefile`:158 listet **acht** Gates plus `record-gates`. Der Lauf lag auf dem Arbeitsbaum, der zu dem Zeitpunkt gegenüber `b806b03` nur die Anker-Klammern gekürzt hatte — Dateizahl und Modulmenge sind davon unberührt.

---

## Findings

### F-1 · MEDIUM — der Herkunfts-Anker `(seit slice-127)` datiert einen Umzug, nicht einen Ursprung

- `kategorie`: MEDIUM
- `quelle`: `grundlagen-traceability.md` §Herkunfts-Anker (Geltungsbereich eng; *Ab Einführung, kein Nachrüsten*); `modul-09` §AGENTS.md-Regeln
- `pfad`: `AGENTS.md`:19 und :306 (`git show b806b03:AGENTS.md`) · `git show ec11ed0:CLAUDE.md`:19-23
- `befund`: Beide Regeln stehen seit dem Bootstrap-Commit `ec11ed0` (2026-06-10) unverändert in `CLAUDE.md`; slice-127 verschiebt sie, er verkörpert sie nicht. Der Kanon reserviert den Anker für Regeln, die *aus dem Steering Loop entstanden* und die 3×-Schwelle erreicht haben, und verbietet das Nachrüsten für Bestandsregeln ausdrücklich („`seit unbekannt` wäre eine Harness-Lüge. Der leere Zustand *ist* die ehrliche Information"). Keine der beiden Regeln hat einen Eintrag im Beobachtungs-Register, und `git log -S` über beide Wortlaute zeigt vor slice-127 keinen anderen Träger als `CLAUDE.md`. Wer den Anker auflöst, sucht in slice-127 nach einer Schwellen-Überschreitung, die es nicht gab.
- `verifizierbar`: nein — kein Gate liest Anker-Semantik; reproduzierbar per `git log -S "Konflikt melden und der höherrangigen Quelle folgen" -- CLAUDE.md` (ältester Treffer `ec11ed0`).
- `klasse`: `herkunfts-anker-datiert-umzug-statt-ursprung`

### F-2 · MEDIUM — „genau die zwei, die der Review des ersten Anlaufs gefunden hat" — der Review fand drei

- `kategorie`: MEDIUM
- `quelle`: BEO-009 Richtung (a); Report des ersten Anlaufs F-3
- `pfad`: Botschaft `b806b03` („Sechs waren gedeckt, zwei nicht — genau die zwei …"), Zuordnungstabelle Zeile 7 · `AGENTS.md`:312 · `.claude/commands/implement-slice.md`:26
- `befund`: Der erste Review beanstandete unter F-3 die Zeile „Kein Erfolg ohne echte Gate-**Ausgabe**" als *nicht* durch §6.8 gedeckt („keine Erfolgsmeldung ohne Gate-**Ausführung**" — Ausführung ist der Lauf, Ausgabe der vorzuzeigende Beleg) und wies den wörtlichen Zwilling in `.claude/commands/implement-slice.md`:26 nach, also außerhalb jeder gerankten Quelle. Die neue Tabelle führt exakt dieselbe Zuordnung wieder, ohne Vorbehalt und ohne die `[ZUGEZOGEN]`-Marke, die sie den zwei anerkannten Waisen gibt; die Zeilen 21 und 22 der Vorfassung bekommen dabei ungleiche Behandlung — Zeile 21 trägt den Vorbehalt „(dort präziser: vor Handoff)", Zeile 22 keinen. `grep -n "Gate-Ausgabe" AGENTS.md` bleibt ohne Treffer, `grep` über die gerankten Quellen ebenfalls.
- `verifizierbar`: nein — reproduzierbar per `grep -n "Gate-Ausgabe\|Gate-Ausführung" AGENTS.md harness/README.md README.md spec/*.md docs/user/*.md` (ein Treffer, §6.8, mit „Ausführung").
- `klasse`: `woertlich-behauptet-sinngemaess-belegt`

### F-3 · MEDIUM — „der Widerspruch gehört benannt" wird dem Kanon als Teil der Konflikt-Hard-Rule zugeschrieben; dort trägt der Satz einen anderen Fall

- `kategorie`: MEDIUM
- `quelle`: `grundlagen-source-precedence.md`:126-130 und :132-135; `modul-02-harness-bootstrap.md`:243-245; BEO-011 Ausprägung (c)
- `pfad`: Botschaft `b806b03` („Der Kanon ist strenger als beide: … und der Widerspruch GEHÖRT BENANNT") · Slice-Plan §1:47-49 · `AGENTS.md`:18-21
- `befund`: Im gepinnten Kanon kommt „der Widerspruch gehört (aber) benannt" genau zweimal vor, beide Male im **Freshness-Audit-Ausgang** *eine `MR-<NNN>` widerspricht der neuen Baseline-Fassung* — dem einen Fall, in dem die niedriger rangierte Quelle gerade **nicht** angepasst wird, sondern in ihrem Geltungsbereich weitergilt; der Satz sitzt dort als Preis für das Weitergelten. Als **universal (Hard Rule)** bezeichnet der Kanon ausschließlich die *Anpassung* (`:132-135`). Botschaft und Plan schreiben beide Hälften dem Kanon zu und stützen darauf die Aussage, der Umzug sei eine „Angleichung" statt einer Verschiebung. Der Text in `AGENTS.md`:19-21 selbst verweist demgegenüber nur für die Anpassung auf den Kanon — Botschaft und Plan reichen weiter als das Artefakt, das sie beschreiben.
- `verifizierbar`: nein — reproduzierbar per `grep -rn "gehört .*benannt" .harness/baseline/v5.11.0/regelwerk/` (zwei Sach-Treffer, beide Freshness-Audit).
- `klasse`: `kanon-zitat-aus-fremdem-fall`

### F-4 · MEDIUM — `welle-83`:87: die Zeile wurde angefasst und behauptet danach zwei Lifecycle-Positionen gleichzeitig

- `kategorie`: MEDIUM
- `quelle`: BEO-008 Richtung 1 (Link-Ziel neu, Prosa alt); BEO-002
- `pfad`: [`welle-83-baseline-v5110-migration.md`](../plan/planning/done/welle-83/welle-83-baseline-v5110-migration.md):87 und :90-91 · Botschaft `753fb46`
- `befund`: `753fb46` hebt in derselben Zeile das Link-Ziel von `next/` auf `in-progress/` und lässt den Prosa-Teil „liegt in `next/` und ist startbar" stehen; die Zeile sagt danach in einem Satz zweierlei, und die stehen gebliebene Hälfte widerspricht dem Verzeichnis. Ebenso bleibt §5 „**Blockiert:** slice-127 — er folgt einer Regel, die erst mit diesem Bump gepinnt ist" unverändert, obwohl derselbe Commit den Slice beansprucht und der Wartegrund laut seiner eigenen Botschaft erledigt ist. Dies ist die Klasse, die der Review zu [slice-128](../plan/planning/done/welle-83/slice-128-baseline-v5110-vendoring.md) vier Tage zuvor an derselben Datei-Familie gemeldet hat.
- `verifizierbar`: nein — `make planning-check` Exit 0 selbst gefahren; das Gate misst Ruhe-Marker ↔ `in-progress/`, nicht Prosa-Lifecycle-Aussagen.
- `klasse`: `prosa-pin-nicht-gehoben`

### F-5 · MEDIUM — eine Verschärfung eines **benannten** Baseline-Schritts ohne Adaptions-Eintrag

- `kategorie`: MEDIUM
- `quelle`: `grundlagen-source-precedence.md`:74-80 (Definition der Adaption) und :116-124; `modul-02` §Freshness-Audit; MR-000 (Baseline-Aussage, Geltungsbereich gesamtes Repo)
- `pfad`: `AGENTS.md`:304-307 · Slice-Plan §2:95-102
- `befund`: §6 Schritt 3 weicht nach dem Commit vom kanonischen Schritt 3 des Minimal Agent Workflow ab: aus „Betroffene Requirement-/ADR-IDs **identifizieren**" wird eine Nenn-Pflicht über fünf Positionen. Der Plan lehnt einen `MR`-Eintrag mit dem Präzedenzfall §3.2 (Suppression-Verbot) ab — §3.2 ändert jedoch keine Baseline-Regel, sondern setzt eine zusätzliche; die hier vorliegende Abweichung ersetzt eine benannte. Der Freshness-Audit läuft über die `MR`-Einträge; eine nur in AGENTS.md-Prosa deklarierte Abweichung hat dort keine Zeile, die beim nächsten Bump gegen die neue Fassung gehalten würde, und der mitgelieferte Auflösungs-Trigger („die Baseline verlangt die Nennung selbst") hat damit keinen Konsumenten, der ihn abfragt.
- `verifizierbar`: nein — kein Gate vergleicht `AGENTS.md` §6 gegen den vendorten `modul-09`-Schritt.
- `klasse`: `baseline-abweichung-ohne-adaptionseintrag`

### F-6 · MEDIUM — zwei Regeln nennen sich „Hard Rule" und liegen außerhalb von §3

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3 („Harte Regeln"); `.harness/skills/reviewer.md` §Eingangs-Kontext („die Hard Rules (`AGENTS.md` §3)")
- `pfad`: `AGENTS.md`:19 (§1) und :306 (§6 Schritt 3)
- `befund`: Beide Zuzüge tragen die Selbstbezeichnung „Hard Rule seit slice-127", liegen aber im Einleitungsabschnitt §1 bzw. als Klammer in einem Schritt der §6-Workflow-Liste. Die §3-Kette `3.1`–`3.7` ist der Ort, den der Reviewer-Skill als Eingangs-Kontext benennt und den §3.7 selbst als Träger von Hard Rules mit Herkunfts-Anker vorführt. Ein Reviewer, der seinem Skill folgt und §3 liest, bekommt die zwei neuen Hard Rules nicht zu sehen; dasselbe gilt für jeden Agenten, der auf „die harten Regeln" verwiesen wird.
- `verifizierbar`: nein — kein Gate prüft, in welchem Abschnitt eine Hard Rule steht (`structure` fährt keine Abschnitts-Invariante über `AGENTS.md` §3).
- `klasse`: `hard-rule-ausserhalb-der-hard-rule-sektion`

### F-7 · LOW — „AGENTS.md §2 (dort steht die Tabelle)": §2 ist eine nummerierte Liste

- `kategorie`: LOW
- `quelle`: BEO-009 Richtung (a); Report des ersten Anlaufs F-6
- `pfad`: Botschaft `b806b03`, Zuordnungstabelle Zeile 3 · `AGENTS.md`:64-76 · `harness/README.md`:35 ff.
- `befund`: Wortgleiche Wiederholung des im ersten Review als F-6 gemeldeten Fehlers: `AGENTS.md` §2 führt neun nummerierte Zeilen, die Tabelle steht in `harness/README.md` §Source precedence. Die Deckung selbst ist tragfähig; nur der Klammerzusatz beschreibt eine Form, die die genannte Stelle nicht hat.
- `verifizierbar`: nein — reproduzierbar durch Lesen von `AGENTS.md`:64-76.
- `klasse`: `zuordnung-nennt-falsche-form`

### F-8 · LOW — die Herkunfts-Anker tragen Prosa statt eines Feldes (durch `b3a74e4` behoben)

- `kategorie`: LOW
- `quelle`: `grundlagen-traceability.md`:33-38 („**Form** — ein Feld, kein Konstrukt"); `modul-09` §Was der Agent in den Kommentar schreibt
- `pfad`: `AGENTS.md`:19-21 und :304-307 (`git show b806b03:AGENTS.md`)
- `befund`: Im geprüften Commit trägt jede der beiden Klammern drei bis vier Sätze Herkunfts- und Begründungsprosa („zugezogen aus `CLAUDE.md`, wo sie als einzige Fundstelle stand"; „Die volle Form buchstabiert das Workflow-Skelett … aus, das sie bis dahin als einzige Quelle trug") — die kanonische Form ist ein Feld nach dem Muster `(seit welle-3)`, die Vorgeschichte gehört nach `modul-09` in §3/§6 des Slice. **Stand bei Review-Ende:** der während des Reviews entstandene Commit `b3a74e4` hat beide Klammern auf Anker, Rang-Zeiger und Auflösungs-Trigger gekürzt; der Befund gilt dem beauftragten Stand und ist am Kopf erledigt.
- `verifizierbar`: nein — kein Gate misst Klammer-Inhalte.
- `klasse`: `herkunfts-prosa-statt-feld`

### F-9 · LOW — die Lesepflicht für `harness/conventions.md` ist beim Umzug enger geworden

- `kategorie`: LOW
- `quelle`: `grundlagen-source-precedence.md` §Vollständigkeit; Botschaft `b806b03`, Zuordnungstabelle Zeile 2
- `pfad`: `git show 84b588e:CLAUDE.md`:5,9 · `harness/README.md`:25-26 · `AGENTS.md`:302-303
- `befund`: Die Vorfassung band das Lesen von `harness/conventions.md` an „**jede** Änderung an Code **oder** Dokumentation". Die als Deckung genannte Stelle — `harness/README.md` §Leseordnung Punkt 2 — trägt den engeren Trigger „vor jeder **Doku-/Konventions**-Änderung", und `AGENTS.md` §6 nennt `conventions.md` in keinem seiner acht Schritte (Schritt 2 verlangt die *kanonische Quelle*; der Konventionsspeicher steht laut Kanon bewusst außerhalb der Rangliste). Für eine reine Code-Änderung existiert die Lesepflicht danach in keiner Quelle. Die Tabellenzeile fasst sechs Leselisten-Punkte zu einer „Aussage" zusammen und macht die Verengung damit unsichtbar.
- `verifizierbar`: nein — reproduzierbar durch Lesen beider Abschnitte.
- `klasse`: `zuordnung-deckt-teilmenge`

### F-10 · LOW — `MR-015` wird als allgemeine Regel zitiert; sein Körper entscheidet eine andere Frage

- `kategorie`: LOW
- `quelle`: MR-015; Maintainability
- `pfad`: Trailer beider Botschaften (`MR-015, slice-127`) · Slice-Plan §Bezug:15 und §5:134-138 · `harness/conventions/MR-015-agents-md-routet.md`:6-35
- `befund`: Plan und Botschaften führen `MR-015` als die Regel „AGENTS.md **routet**, spiegelt nicht" und leiten daraus das zentrale Risiko des Slice ab. Der Eintrag selbst hat den Geltungsbereich `AGENTS.md` §1/`MR-011`/`MR-012`/§Adoptierte Konventions-Quellen und entscheidet genau eine Frage: dass `AGENTS.md` §1 für Quelldatei und Stand der Baseline auf `conventions.md` routet, statt die `agents-regelwerk.md`-Raw-URL zu spiegeln (Auflösungs-Trigger: permanent, *Provenienz*). Die allgemeine Nicht-Duplikations-Regel steht in `AGENTS.md`:10-12 und im Kanon (`modul-09` §Ziel-Form: AGENTS.md), nicht in `MR-015`. Der Trailer bleibt formal korrekt (die Änderung liegt in §1, dem deklarierten Geltungsbereich); die Begründungslast trägt der Eintrag nicht.
- `verifizierbar`: nein — `make trace-check` prüft die Kennungs-*Form*, nicht die Passung.
- `klasse`: `mr-als-allgemeine-regel-zitiert`

### F-11 · INFO — „sechs Verweise" sind zwölf Verweis-Zeilen in sechs Dateien

- `kategorie`: INFO
- `quelle`: BEO-009 Richtung (a)
- `pfad`: Botschaft `753fb46` („gebündelt nach MR-013 mit dem Roadmap-Flip und den sechs Verweisen")
- `befund`: `git show 753fb46` ändert **zwölf** Verweis-Zeilen, verteilt auf **sechs** Dateien (`slice-128` 1, `slice-129` 5, `slice-131` 1, `welle-83` 3, zwei Review-Reports je 1), plus Roadmap-Flip und den Move selbst. Die Zahl stimmt für Dateien, nicht für Verweise; sie ist die einzige nachzählbare Größe der Botschaft, deren Bezugsmenge nicht mitgenannt ist.
- `verifizierbar`: ja — `git show 753fb46 -- <die sechs Dateien> | grep -c "^+[^+]"` = 12.
- `klasse`: `zahl-ohne-bezugsmenge`

### F-12 · INFO — der erste `(seit slice-<NNN>)`-Anker des Repos hat kein auflösendes §7-Feld

- `kategorie`: INFO
- `quelle`: `grundlagen-traceability.md`:39-46 und :52-60 (Sensor 1: Pflichtfeld `liegt in <Zielort>` in §7 der `done/`-Slice-Datei)
- `pfad`: `AGENTS.md`:19, :306 · Slice-Plan §7:179-193 · Slice-Plan §4 (DoD)
- `befund`: Der Kanon lässt `seit slice-<NNN>` „über `done/slice-<NNN>-*.md` §7" auflösen und macht das Pflichtfeld `liegt in <Zielort>` zum Auslöser der Anker-Paarung. `grep -rn "liegt in \`" docs/plan/planning/` findet im ganzen Repo keine §7-Instanz dieses Feldes; slice-127 schreibt als erster Slice überhaupt einen `(seit slice-<NNN>)`-Anker in `AGENTS.md`, und weder sein §7 noch seine DoD sehen das Feld vor. Der Retirement-Check („Regel seit slice-127 — ist die Beobachtung seither wieder aufgetreten?") hätte kein Ziel. Repo-weite Lage; hier zum ersten Mal mit Konsequenz.
- `verifizierbar`: nein — kein Gate im Repo prüft die Anker-Paarung.
- `klasse`: `anker-ohne-aufloesendes-feld`

---

## Negativbefunde (geprüft, ohne Befund)

1. **Zeilenweiser Abgleich alt → neu, alle 23 Zeilen.** `git show 84b588e:CLAUDE.md` trägt Titel (Z. 1), Modus-Satz (Z. 3), Lese-Trigger + sechs Punkte (Z. 5-12), Überschrift „Regeln:" (Z. 14) und sechs Regelzeilen (Z. 16-23). Die acht Tabellenzeilen decken jede Sachzeile ab; eine **neunte, gar nicht genannte** Aussage habe ich nicht gefunden. Die Ungenauigkeiten liegen innerhalb der Zeilen (F-2, F-9) und in Zeile 3, die zwei Aussagen trägt („folgt dem AI-Harness-Prozess" + „Greenfield: Doc führt") und nur mit der zweiten in der Tabelle steht — die erste ist über `AGENTS.md` §1 und MR-000 gedeckt.
2. **Zuordnung Zeile 1 (Greenfield).** `harness/conventions.md`:158 trägt „`*` (Default für gesamtes Repo) | Greenfield | Projekt startet spec-first; **Doc führt, Code folgt**" — die Aussage steht dort wörtlich. Anzumerken bleibt, dass `conventions.md` der **Konventionsspeicher** ist und nicht in der Rangliste steht; der Kanon lässt ihn als Fundort ausdrücklich zu (§Vollständigkeit nennt drei Orte), die DoD des Slice sagt enger „in einer **gerankten** Quelle". Kein Verstoß gegen den Kanon, deshalb kein Finding — die Prüfung der DoD-Formulierung gehört der Verifikation.
3. **Zuordnung Zeile 4 (make-only).** `AGENTS.md` §3.1 deckt „Nur `make`-Targets für Checks und Gates; keine Host-Paketmanager oder -Toolchains" vollständig und breiter (Host-`go` zusätzlich genannt).
4. **Zuordnung Zeile 6 (`make gates`).** `AGENTS.md` §6.6 „Repo-weiten Gate-Lauf vor Handoff" plus §4 („mandatory vor Handoff") tragen die Aussage; der Klammerzusatz „dort präziser: vor Handoff" beschreibt einen früheren, nicht einen engeren Zeitpunkt, was die Pflicht nicht schwächt.
5. **Zuordnung Zeilen 2 und 5.** Von den sechs Leselisten-Punkten haben fünf eine tragfähige Fundstelle (`harness/README.md` §Leseordnung 1-3, `AGENTS.md` §6.1 für `harness/README.md`, §6.2 + §2 für ADRs und `spec/`); der sechste trägt F-9. Die Benenn-Pflicht ist mit §6.3 real zugezogen und vollständig (alle fünf Positionen).
6. **Kanon-Treue der Rolle von `CLAUDE.md`.** Der Kanon verlangt: „Sie **verweist** … und legt nichts fest" und „bringt `AGENTS.md` in den Lauf-Kontext, wo Modul 9 es für jeden Lauf verlangt". Die neue Zeile tut genau das: „Vor jeder Änderung an Code oder Dokumentation zuerst `AGENTS.md` lesen und befolgen." Der Imperativ ist **keine** Festlegung, sondern das Ausbuchstabieren einer kanonischen Forderung (`modul-09`: „gehört in jeden Lauf-Kontext"); der Trigger „Code oder Dokumentation" ist deckungsgleich mit `AGENTS.md` §1. Kein Eigeninhalt bleibt zurück.
7. **Wahrheit des Pointers.** „sie trägt die Hard Rules und routet weiter" — §3 trägt sieben Hard Rules, §2 routet auf die kanonischen Quellen, §1 auf `harness/conventions.md` und die vendorte Baseline, §4 auf die Gates, §6.1 auf `harness/README.md`. Beide Zusagen sind eingelöst; das Wort „Leseordnung", an dem der erste Anlauf scheiterte, kommt nicht mehr vor.
8. **Der Widerspruchs-Fall „zwei kanonische Quellen".** Vom Kanon gedeckt: `:132-135` erklärt „dass bei Konflikt die niedriger rangierte Quelle angepasst wird" als **universal (Hard Rule)** — das trägt auch zwei gerankte Quellen gegeneinander. Die neue Formulierung greift zugleich **nicht** in die zwei Fälle über, die der Kanon anders regelt (Konventionsspeicher = „kein Konflikt, sondern eine Zuständigkeit"; `MR-<NNN>` gilt in ihrem Geltungsbereich vor der Baseline) — sie nennt ihre beiden Fälle ausdrücklich.
9. **Scope-Treue gegen §3, alle vier Ausschlüsse.** `git diff --stat 84b588e..b806b03`: kein `.claude/`-Pfad berührt (Hooks, `settings.json`, Workflow-Skelett unverändert); `.d-check.yml` unverändert (kein `pins`, kein `dpin`, kein `FOCUS_DISABLE`-Nachzug); kein neues Gate, kein Makefile-Diff. „Keine neue Hard Rule" hält im Sinn des Plans: beide Regeln existierten, der Wechsel von „folgen" auf „angepasst" ist die vom Plan §2 Schritt 1 ausdrücklich beauftragte Angleichung an den Kanon, keine Erfindung. (Der Instrumentierungs-Punkt dazu steht als F-5, nicht als Scope-Bruch.)
10. **Move-Commit-Form.** `git show --stat -M 753fb46` zeigt `{next => in-progress}` mit **0** geänderten Zeilen an der Slice-Datei (R100); Roadmap-Flip und Verweis-Nachzug sind gebündelt, wie es `make planning-check` erzwingt. Selbst gefahren: Exit 0.
11. **Ruhe-Marker.** `in-progress/` trägt seit `753fb46` einen Slice, also muss „Nichts in Arbeit." weichen; der Wächter misst beide Richtungen, `planning-check` Exit 0 bestätigt den Zustand. Die verbleibenden drei aufeinanderfolgenden Leerzeilen in `roadmap.md`:33-36 haben keinen Konventions-Anker (kein markdownlint in `Makefile` oder CI) — nach dem Anti-Pattern *Kein Stil-Polizist* kein Finding, nur festgehalten.
12. **Zahlen der Botschaften.** „23 auf 4 Zeilen": `git show 84b588e:CLAUDE.md` = 23 Zeilen, `CLAUDE.md` = 4 — bestätigt. „461 Dateien, 0 Befunde": selbst gefahren, exakt diese Zeile dreimal. „acht Gates": `Makefile`:158 listet acht plus `record-gates` (Nachweis, kein Gate). „`make gates` Exit 0": selbst gefahren, `EXIT=0` aus `$?`.
13. **Eingehende Verweise auf `CLAUDE.md`.** `grep -rn "CLAUDE" .` über `.claude/`, `.github/`, `Makefile`, Skripte, YAML und Markdown: Sach-Treffer sind Review-Reports und `.claude/commands/implement-slice.md`:7 („Read `CLAUDE.md`"); der einzige Nicht-Doku-Treffer ist `CLAUDE_PROJECT_DIR`. Der Kommando-Ablauf bleibt funktionsfähig — Schritt 1 liest jetzt den Pointer, Schritt 3 liest `AGENTS.md` ohnehin als eigenen Schritt.
14. **Delta-Audit-Konsistenz.** Die zwei Zuzüge decken sich genau mit den zwei Waisen, die der Fünf-Fundorte-Zensus (`slice-129` §4a) für `CLAUDE.md` führt („Meldepflicht bei Quellen-Konflikt, Benenn-Pflicht vor der Implementierung"); der dort als „Grenzfall / im Slice zu entscheiden" geführte Workflow-Skelett-Punkt ist durch §6.3 mit-entschieden. Die übrigen drei Fundorte bleiben korrekt unangetastet (`slice-131` und Folgepunkte).
15. **Zustandsfeld-Hygiene (`AGENTS.md` §3.7).** Der Diff fügt kein Code-, Konfigurations- oder Skript-Kommentar und kein Zustandsfeld hinzu; der Slice-Plan trägt `**Lifecycle:**` statt `**Status:**`.
16. **Traceability-Form.** Beide Botschaften nennen `MR-015` bzw. `MR-013` und `slice-127`; `make trace-check` ist bewusst nicht Teil von `gates`, die Form ist erfüllt.
17. **Hooks.** Nicht berührt und im Plan §3 ausdrücklich ausgeschlossen; der Plan macht — anders als der erste Anlauf — keine Aussage mehr darüber, was sie abdecken. Die im ersten Review als F-5 gemeldete Überdehnung ist damit gegenstandslos.
18. **`make gates` gegen den geprüften Stand.** Exit 0, acht Gates grün, 0 Befunde. Der Lauf ist kein DoD-Haken (das ist Verifier-Arbeit), sondern der Beleg, dass keiner der obigen Befunde von einem Gate gefangen wird.

---

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
| --------- | ------ | -------- |
| HIGH      | 0      | — |
| MEDIUM    | 6      | F-1, F-2, F-3, F-4, F-5, F-6 |
| LOW       | 4      | F-7, F-8 (durch `b3a74e4` behoben), F-9, F-10 |
| INFO      | 2      | F-11, F-12 |

---

## Verdikt

**Blockierend** — kein HIGH, sechs MEDIUM, vier LOW (davon eines am Kopf bereits erledigt), zwei INFO.

Der Neuschnitt trägt: die Reihenfolge des Kanons (belegen → umziehen → kürzen) ist eingehalten, die zwei anerkannten Waisen sind real umgezogen statt gelöscht, die Pointer-Zeile verspricht nur noch, was `AGENTS.md` einlöst, alle vier Out-of-Scope-Zusagen halten, der Move ist R100, und jede nachzählbare Zahl außer einer stimmt. Die Kernfrage des Slice — *ist die Vollständigkeits-Aussage diesmal am Bestand statt am Anlass gebildet?* — ist zu etwa vier Fünfteln mit Ja zu beantworten: die Tabelle ist zeilenweise geführt und hält für sechs der acht Zeilen.

Sie kippt aber an derselben Stelle wie beim ersten Mal, nur kleiner. Der erste Review hat **drei** Zuordnungen beanstandet, nicht zwei; die dritte (F-2) steht unverändert wieder in der Tabelle, ohne den Vorbehalt, den die Nachbarzeile bekommt — und die Botschaft erklärt die zwei behandelten zu „genau den zwei, die der Review gefunden hat". Dazu kommen zwei Aussagen über den Kanon, die er so nicht trägt (F-3) bzw. deren Instrumentierung er anders regelt (F-5), und ein Herkunfts-Anker, der einen Umzug als Ursprung datiert, obwohl der Kanon das Nachrüsten ausdrücklich verbietet (F-1). F-6 ist der leiseste und zugleich der folgenreichste: zwei Regeln nennen sich Hard Rule und liegen außerhalb des Abschnitts, auf den jeder Reviewer-Auftrag zeigt.

F-4 ist unabhängig von der Kürzung und trifft den Move-Commit: eine angefasste Zeile, deren Prosa-Hälfte danach dem Verzeichnis widerspricht — dieselbe Klasse, die der Review zu `slice-128` an derselben Datei-Familie gemeldet hat. Das ist die dritte Wiederholung dieser Klasse in kurzer Folge und damit nach der Kontext-Eskalations-Regel ein Steering-Loop-Signal, kein Einzelbefund.

F-7 bis F-12 blockieren nicht. F-8 ist bereits behoben und steht nur zur Vollständigkeit des geprüften Stands; F-12 ist eine repo-weite Lage, die dieser Slice als erster mit Konsequenz belegt.
