# Review-Report: slice-140 — Konsumenten-CR an den Kurs — 2026-08-25

**Review-Art:** Plan-/Doku-Review (Zitate, Fundstellen und Reichweiten gegen
die vendorte Baseline und gegen den eigenen Bestand; Modul 10 §Drei
Review-Arten) · **Gegenstand:** Commit `9b7c932` (`feat(plan): Konsumenten-CR
an den Kurs geschrieben — vier Punkte, Ablage deklariert (MR-035,
slice-140)`) plus die Datums-Berichtigung `070be0e`, soweit sie dieselben
Dateien berührt. 3 Dateien laut `git show --stat 9b7c932` (neu
`docs/plan/cr/2026-08-23-cr-regelwerk-v5110.md` +174, `harness/conventions.md`
+1, neu `harness/conventions/MR-035-cr-ablage.md` +29), 204 Einfügungen /
0 Löschungen; aus `070be0e` dazu die Umbenennung nach
`docs/plan/cr/2026-08-25-cr-regelwerk-v5110.md` (+1/-1, Kopf-Datum) und
`harness/conventions/MR-035-cr-ablage.md` (+1/-1, Feld `Datum`).

**Skill:** `.harness/skills/reviewer.md` @ 1.10.0 · **Modell-ID:**
`claude-opus-5[1m]` · **Datum:** 2026-08-25

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-140-konsumenten-cr.md`
  vollständig (§1 Ziel, §3 NICHT, §4 DoD, §5 Risiken, §7 Vorprüfungen)
- Das zu prüfende Dokument
  `docs/plan/cr/2026-08-25-cr-regelwerk-v5110.md` vollständig (174 Zeilen)
- `harness/conventions/MR-035-cr-ablage.md` vollständig und die neue
  Index-Zeile in `harness/conventions.md` §Aktive Adaptionen
- `AGENTS.md` §3 (Hard Rules, insbesondere §3.4 und §3.7) und §5
- `harness/conventions.md` → `MR-000` (Baseline-Aussage), `MR-033`, `MR-035`
- `docs/plan/planning/observations.md` → `BEO-011`, `BEO-012`, `BEO-015`
- Vendorte Baseline `.harness/baseline/v5.11.0/` — für die Zitat-Prüfung im
  Volltext: `regelwerk/modul-03-spec.md` §Ziel-Form: Architektur-Sicht,
  `regelwerk/modul-05-planning-harness.md` §Offene Risiken werden bei Closure
  aufgelöst, `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register,
  `regelwerk/modul-13-quality-gates.md` §Fitness Function aus einem ADR-Satz,
  `regelwerk/grundlagen-source-precedence.md` §Source Precedence,
  `regelwerk/grundlagen-harness-dateien.md` §Konventionsspeicher,
  `regelwerk/modul-02-harness-bootstrap.md` §Freshness-Audit,
  `templates/AGENTS.template.md` §3.4,
  `templates/spec/architecture.template.md` vollständig,
  `templates/docs/plan/planning/slice.template.md` §6
- **Vorherige Findings am gleichen Gegenstand:**
  `docs/reviews/2026-08-23-slice-136-agents-34-klaerung-review.md` (F-1
  „gedeckt-Verdikt deckt nur eine von fünf Kategorien", F-3
  „vorlagen-praxis-statt-direkterer-textbeleg" — derselbe `modul-03`-Satz wie
  Punkt 2 des CR, F-4 „suchmuster-abhaengige-kennzahl"),
  `docs/reviews/2026-08-25-slice-137-toolchain-freshness-review.md` (F-4
  „unbelegte-uebernommene-quellenaussage"),
  `docs/reviews/2026-08-25-slice-138-matrix-wellen-klasse-review.md` (F-3
  „bestandsaussage-ueber-modul-faehigkeit-unbelegt"),
  `docs/reviews/2026-08-25-slice-139-closure-ausgang-waechter-review.md`
  (F-1/F-2 — die zwei Fallen, die Punkt 1 des CR als Beleg führt)
- `tools/harness/closure-outcomes.sh` vollständig (92 Zeilen) und
  `internal/hexagon/core/rules/planning.go` (`placeholderRE`,
  `checkClosurePlaceholder`) — zur Prüfung der in Punkt 1 genannten Zahlen und
  der dritten „Falle"
- `.d-check.yml` (`scan.roots`/`ignore`, `ids.exempt-paths`,
  `links.resolve-from`) — zur Frage, welche Zusage für den neuen Pfad
  `docs/plan/cr/` überhaupt gilt

**Nicht erhalten:** die DoD-Abhakung (Verifikations-Rolle, getrennter
Kontext, anderes Prüf-Artefakt).

**Vom Reviewer selbst gefahren** (ausschließlich Lesekommandos; kein `make`,
kein Netzzugriff, keine Änderung an einer Repo-Datei außer diesem Report):
`git show`/`git log --format=%B`/`git ls-tree`/`git log --follow` auf
`9b7c932`, `070be0e`, `612a619`, `2378cb4`, `cf41502` und auf die Slice-Dateien
von slice-132/134/138/139; `grep -rn` über den **gesamten** Baum
`.harness/baseline/v5.11.0/` für jedes im CR zitierte Satzfragment sowie für
`modul[- ]?pfad` (case-insensitiv), `Platzhalter`, `Change Request`,
`Konsument`, `Geltungsbereich`, `Reichweite`; `find … -name 'slice-*.md' | wc -l`
über `docs/plan/planning/done` (Ist-Stand und Stand je Commit über
`git ls-tree`); Lesen von `planning.go`, um die in Punkt 1 behauptete
Whitespace-Grenze des produkteigenen Platzhalter-Erkenners zu prüfen.

**Verdikt: blockierend** — kein HIGH, drei MEDIUM, drei LOW.

---

## Findings

### F-1 — „Belege aus einem einzigen Arbeitstag" ist durch die eigene git-Historie widerlegt: zwei der drei Instanzen stammen vom 23., die dritte vom 25. August

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §5 („Eine Commit-Botschaft oder Closure-Notiz
  behauptet **nicht mehr, als die Arbeit trägt** … ihr Schluss reicht nicht
  weiter als die gemessene Menge", `BEO-009`) ·
  [`BEO-012`](../plan/planning/observations.md) (Reichweite einer belegten
  Aussage)
- **pfad:** `docs/plan/cr/2026-08-25-cr-regelwerk-v5110.md:120` („**Belege aus
  einem einzigen Arbeitstag**, alle in einem Repo, das diese Disziplin ernst
  nimmt"); gespiegelt in
  `docs/plan/planning/in-progress/slice-140-konsumenten-cr.md:54` und in der
  Commit-Botschaft von `9b7c932` („Drei Instanzen aus einem Arbeitstag, jede
  mit ihrer Form")
- **befund:** Die drei Instanzen von Punkt 3 liegen nachweislich auf **zwei**
  Kalendertagen: Instanz 1 (Suppression-Probe nennt den falschen Linter) steht
  in `docs/plan/planning/done/slice-132-hard-rule-zensus.md:167` und Instanz 2
  (Direktiv-Form galt als „wohlgeformt", nie gefahren) in
  `docs/plan/planning/done/slice-134-nolintlint.md:145-146` — beide zu
  welle-84, deren Commits durchgehend das Autor-Datum **2026-08-23** tragen
  (`1cbd13b`, `2539409`, `54d102d`); Instanz 3 (Probe ans Dateiende angehängt,
  landet im ausgenommenen Abschnitt) steht in
  `docs/plan/planning/done/slice-138-matrix-wellen-klasse.md:177-181`, dessen
  Commits das Autor-Datum **2026-08-25** tragen (`d30d2b1`, `a375222`,
  `7e0d571`, `42b7651`). Die Aussage entstand unter der in `070be0e` selbst
  berichtigten Datums-Verwechslung; dieselbe Berichtigung hat das Kopf-Datum
  des CR auf 2026-08-25 gezogen (`070be0e`, `+1/-1`), den Satz im Rumpf aber
  nicht mehr berührt — seither ist er an derselben Datei prüfbar falsch.
- **verifizierbar:** ja — `git log --format='%h %ad' --date=short -- docs/plan/planning/done/slice-134-nolintlint.md`
  gegen dieselbe Abfrage auf `slice-138-matrix-wellen-klasse.md` zeigt die zwei
  Tage; kein Gate deckt Prosa-Behauptungen über Zeitpunkte.
- **klasse:** dichte-behauptung-ueber-einen-tag-durch-eigene-historie-widerlegt

### F-2 — Punkt 4 erklärt die Reichweitenfrage im Kanon für ungeregelt, während die von ihm genannte Datei sie für genau eine seiner drei Formen ausdrücklich regelt

- **kategorie:** MEDIUM
- **quelle:** [`BEO-012`](../plan/planning/observations.md) („Eine Quelle wird
  über ihren Geltungsbereich hinaus zitiert" — hier in der Umkehrung: eine
  Quelle wird für *weniger* in Anspruch genommen, als sie sagt) ·
  slice-140 §5 Risiko 3 („Vor dem Schreiben ist der ganze Kanon danach zu
  durchsuchen, nicht nur die zwei bekannten Stellen")
- **pfad:** `docs/plan/cr/2026-08-25-cr-regelwerk-v5110.md:145-160` („**Was
  dasteht.** `grundlagen-source-precedence.md` ordnet Quellen und regelt den
  Konfliktfall … **Was fehlt.** Eine Aussage darüber, **wie weit eine einzelne
  Aussage aus einer gerankten Quelle trägt**" mit der zweiten Form „Ein
  Adaptions-Eintrag wird nach seinem **Titel** zitiert statt nach seinem Feld
  `Geltungsbereich`") gegen
  `.harness/baseline/v5.11.0/regelwerk/grundlagen-source-precedence.md:76-77`,
  `:111-113` und `:132-135`
- **befund:** Dieselbe Datei, die der CR als „ordnet Quellen und regelt den
  Konfliktfall" zusammenfasst, trägt drei Aussagen über **Reichweite**: (a)
  Zeile 76–77 setzt sie für Adaptions-Einträge wörtlich — *„Eine `MR-<NNN>`
  gilt innerhalb ihres deklarierten Geltungsbereichs vor der Baseline;
  außerhalb davon gilt die Baseline unverändert"* —, (b) Zeile 111–113 trennt
  dafür zwei Felder (`Geltungsbereich` = Repo-Ausschnitt,
  `Ersetzt-Baseline-Regel` = Baseline-Ausschnitt), und (c) Zeile 132–135
  markiert unter *Universal vs projektabhängig* ausdrücklich, welcher Teil der
  eigenen Kernaussage universal gilt und welcher projektspezifisch ist. Die
  zweite der drei im CR aufgezählten Formen ist damit keine Kanon-Lücke,
  sondern der Verstoß gegen eine bestehende, in derselben Datei stehende
  Regel — was der eigene Registereintrag `BEO-012` in seiner Prozedur-Spalte
  auch so führt („bei einer `MR-<NNN>` das Feld `Geltungsbereich` **und**
  `Ersetzt-Baseline-Regel`"). Das §5-Risiko „könnte sich als Lesefehler
  erweisen" wurde ausweislich der Commit-Botschaft nur für Punkt 2 durch eine
  Baum-Suche entkräftet, für Punkt 4 nicht.
- **verifizierbar:** ja — `grep -n "gilt innerhalb ihres" .harness/baseline/v5.11.0/regelwerk/grundlagen-source-precedence.md`
  liefert Zeile 76; kein Gate prüft die Vollständigkeit einer Lücken-Behauptung
  über eine Fremdquelle.
- **klasse:** behauptete-kanon-luecke-in-der-zitierten-datei-bereits-geschlossen

### F-3 — Das Kronzeugen-Zitat für Punkt 1 gibt Reihenfolge und Wortlaut der Quelle nicht wieder: die zwei Sätze stehen in umgekehrter Folge und in zwei verschiedenen Aufzählungspunkten

- **kategorie:** MEDIUM
- **quelle:** [`BEO-012`](../plan/planning/observations.md) ·
  `AGENTS.md` §5 · Baseline
  `.harness/baseline/v5.11.0/regelwerk/modul-06-roadmap.md:103-104`
- **pfad:** `docs/plan/cr/2026-08-25-cr-regelwerk-v5110.md:29-34` („benennt
  dann, was maschinell entscheidbar ist: *„drei Prüfungen ohne Urteil:
  **Form** · **Anzahl** · **Lage**"*, gefolgt von *„Muster: schreiben →
  committen → Gate prüft. Welches Werkzeug, ist Repo-Entscheidung."*")
- **befund:** In der Quelle steht der Satz *„Muster: schreiben → committen →
  Gate prüft. Welches Werkzeug, ist Repo-Entscheidung."* am **Ende** des
  Aufzählungspunkts *Mensch urteilt, Maschine prüft Deckung* (Zeile 103); die
  drei Prüfungen stehen im **folgenden**, eigenen Punkt *Der Beleg ist
  formgebunden* (Zeile 104). Das „gefolgt von" dreht diese Folge um. Zusätzlich
  ist die als Zitat ausgewiesene Zeichenkette nicht der Wortlaut: die Quelle
  schreibt *„**Form** (Slice-Kennung `slice-<NNN>`, kein Freitext) ·
  **Anzahl** (so viele wie der Zähler) · **Lage** (führt das Repo die
  Slice-Datei, liegt sie in `done/`)"*, und genau diese drei ohne Auslassungs­
  zeichen getilgten Klammern binden die drei Prüfungen an den **Beleg einer
  Registerzeile** — nicht an ein allgemein übertragbares Prüfschema, als das
  der gekürzte Satz im CR erscheint. Der Adressat, der die Stelle nachschlägt,
  findet weder die Reihenfolge noch die Zeichenkette so vor.
- **verifizierbar:** ja — `grep -n "drei Prüfungen ohne Urteil" .harness/baseline/v5.11.0/regelwerk/modul-06-roadmap.md`
  liefert Zeile 104, `grep -n "Welches Werkzeug"` dieselbe Datei Zeile 103;
  kein Gate prüft Zitattreue gegen den vendorten Baum.
- **klasse:** zitat-reihenfolge-und-wortlaut-weichen-von-der-quelle-ab

### F-4 — „Sechs konstruierte Verstöße" zählt vier Platzhalter-Formen und zwei Fail-Closed-Proben zusammen; die gebaute Form kennt vier Formen

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §5 (`BEO-009`: „wer N Formen geprüft hat, berichtet
  N") · vorheriges Finding
  `docs/reviews/2026-08-23-slice-136-agents-34-klaerung-review.md` F-4
  („suchmuster-abhaengige-kennzahl")
- **pfad:** `docs/plan/cr/2026-08-25-cr-regelwerk-v5110.md:48-52` („Sechs
  konstruierte Verstöße wurden rot gesehen, jeder mit gelesener Fundstelle")
  gegen `tools/harness/closure-outcomes.sh:54-60` („Vier Formen: zwei
  repo-eigene und zwei der Kanon-Vorlage") sowie die Botschaften von `612a619`
  („BEWUSSTES BRECHEN, vier Formen") und `2378cb4` („Vierte Form aufgenommen …
  Probe: unlesbare Datei")
- **befund:** Die Sechs entsteht erst aus der Summe zweier Commits und mischt
  dabei zwei Proben-Klassen: **vier** davon sind Verstöße gegen die
  Platzhalter-Regel (`(bei Closure)`, `wird mit dem Closure-Body gefüllt`,
  `<…>`, `<eingetreten:`), **zwei** sind Fail-Closed-Proben am Wächter selbst
  (leere Prüfmenge in `done/`; unlesbare Datei) und keine unaufgelösten
  Ausgänge. Für diese zwei gibt es auch keine „gelesene Fundstelle" im Sinne
  von Datei und Zeile, sondern nur eine Abbruchmeldung. Ein Adressat, der den
  Baubarkeits-Beleg nachzählt, findet vier konstruierte Regelverstöße und zwei
  Robustheitsproben, nicht sechs Verstöße; die Zahl steht in keiner Datei des
  Repos.
- **verifizierbar:** ja — `git log --format=%B -1 612a619` und
  `git log --format=%B -1 2378cb4` listen die Proben einzeln auf;
  `sed -n '54,60p' tools/harness/closure-outcomes.sh` nennt vier Formen.
- **klasse:** probenzahl-mischt-regelverstoss-und-fail-closed-probe

### F-5 — MR-035 schreibt dem Kanon einen Vorgang zu, den der vendorte Baum nicht kennt

- **kategorie:** LOW
- **quelle:** [`BEO-012`](../plan/planning/observations.md) · Baseline
  `.harness/baseline/v5.11.0/regelwerk/grundlagen-harness-dateien.md:234`
  (Pflichtfelder eines Adaptions-Eintrags, darunter `Ersetzt-Baseline-Regel`)
  · Baseline
  `.harness/baseline/v5.11.0/regelwerk/modul-02-harness-bootstrap.md:218-220`
  (Freshness-Audit: Frage pro Eintrag „Regelt die neue Fassung das, wofür
  diese Adaption angelegt wurde?")
- **pfad:** `harness/conventions/MR-035-cr-ablage.md:3-6` („**Ersetzt-Baseline-Regel:**
  keine. Der Kanon kennt den **eingehenden** Change Request (Vertragsänderung
  am Lastenheft) und den Konsumenten-CR als Vorgang, aber keine Ablage für den
  **ausgehenden**")
- **befund:** Der vendorte Baum kennt den *Konsumenten-CR* nicht: eine Suche
  über `.harness/baseline/v5.11.0/` nach `Konsument` findet ausschließlich den
  Artefakt-Konsumenten (`grundlagen-harness-dateien.md` §Jedes Artefakt hat
  einen Konsumenten, `modul-10`, `modul-11`), nach `\bCR\b` ausschließlich den
  **externen, eingehenden** Change Request gegen das Lastenheft — den
  `grundlagen-source-precedence.md:181-183` und `grundlagen-begriffe.md:45`
  ausdrücklich als *„bewusst kein Harness-Konstrukt — kein `CR-*`-ID-Schema,
  keine eigene Datei, kein Gate"* führen. *Konsumenten-CR* ist eine Prägung
  dieses Repos (ältester Beleg
  `docs/plan/planning/done/slice-078-ignore-refs-quell-skopus.md:244`,
  2026-07-17). Die Prämisse des Eintrags ist damit nur zur Hälfte belegt, und
  der Freshness-Audit, der pro Eintrag genau nach dieser Prämisse fragt, hat
  für die andere Hälfte nichts zum Vergleichen.
- **verifizierbar:** ja — `grep -rn "Konsument" .harness/baseline/v5.11.0/`
  liefert keinen Treffer auf einen ausgehenden CR; kein Gate prüft
  Kanon-Zuschreibungen in Prosa.
- **klasse:** unbelegte-kanon-zuschreibung-in-adaptions-eintrag

### F-6 — Die Commit-Botschaft sagt „jeder mit Zitat und Fundstelle" zu; Punkt 4 trägt weder Zitat noch Abschnitt

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §5 („Eine Commit-Botschaft … behauptet nicht mehr,
  als die Arbeit trägt")
- **pfad:** Commit-Botschaft `9b7c932`, erster Absatz („Vier Punkte gegen
  Regelwerk v5.11.0, jeder mit Zitat und Fundstelle") gegen
  `docs/plan/cr/2026-08-25-cr-regelwerk-v5110.md:145-146`
- **befund:** Die Punkte 1, 2 und 3 eröffnen mit „**Was dasteht.**" plus
  Blockzitat und benanntem Abschnitt (`modul-05` §Offene Risiken …,
  `modul-03` §Ziel-Form: Architektur-Sicht und `AGENTS.template.md` §3.4,
  `modul-13` §Fitness Function aus einem ADR-Satz). Punkt 4 nennt nur die
  Datei `grundlagen-source-precedence.md`, ohne Abschnitt und ohne einen
  einzigen zitierten Satz, und referiert sie in eigenen Worten. Ausgerechnet
  der Punkt, dessen Gegenstand die Reichweite eines Zitats ist, ist der
  einzige ohne Zitat — und ein Adressat, der die Botschaft als Wegweiser
  benutzt, sucht die Fundstelle vergeblich.
- **verifizierbar:** ja — `grep -n "^> " docs/plan/cr/2026-08-25-cr-regelwerk-v5110.md`
  zeigt Blockzitate nur unter den Punkten 1, 2 und 3.
- **klasse:** botschaft-sagt-belegform-zu-die-ein-punkt-nicht-traegt

## Negativbefunde

- geprüft, ohne Befund: **Punkt 1, `modul-05`-Zitat** — der Blockzitat-Text
  (CR:22-27) stimmt mit
  `.harness/baseline/v5.11.0/regelwerk/modul-05-planning-harness.md:129-136`
  überein; die Auslassung in der Mitte ist mit „[…]" ausgewiesen, die
  weggelassenen Quellen-Klammern („(Modul 7)", „(ohne sie ist es stilles
  Vergessen)") und der Modul-6-Link ändern die Aussage nicht. Der
  Abschnittsname *Offene Risiken werden bei Closure aufgelöst* ist die
  wörtliche `####`-Überschrift.
- geprüft, ohne Befund: **Punkt 1, Zahl „137 Slices, null Treffer"** — die
  Angabe ist ausdrücklich auf „beim Scharfschalten" bezogen und stimmt für
  diesen Zeitpunkt exakt: `git ls-tree -r --name-only 612a619` zählt 137
  `done/slice-*.md`, und sowohl `612a619` als auch der Nachbesserungs-Commit
  `2378cb4` melden „137 Slices, null offene Platzhalter". Der heutige Bestand
  (138) widerspricht nicht — er ist ein anderer Zeitpunkt.
- geprüft, ohne Befund: **Punkt 1, die drei Fallen** — alle drei sind am
  Bestand auffindbar und treffen zu: Falle 1 und 3 sind F-1 bzw. F-2 aus
  `docs/reviews/2026-08-25-slice-139-closure-ausgang-waechter-review.md`;
  Falle 2 ist am Quelltext bestätigt — `placeholderRE` in
  `internal/hexagon/core/rules/planning.go:214` verlangt mit `[^\s<>]`
  tatsächlich whitespace-freie Winkelklammern und kann die Vorlagen-Zeile
  deshalb nicht sehen.
- geprüft, ohne Befund: **Punkt 2, beide Zitate** — `modul-03-spec.md:126-127`
  und `templates/AGENTS.template.md:130-131` sind **wortgleich** wiedergegeben,
  einschließlich Auszeichnung; die Zuschreibungen (`modul-03` §Ziel-Form:
  Architektur-Sicht bzw. Template-§3.4 mit der Überschrift *3.4 Architektur ist
  sprach- und meilensteinfrei*) sind korrekt.
- geprüft, ohne Befund: **Punkt 2, Erschöpfungs-Behauptung „genau diese zwei
  Fundstellen"** — selbst nachgezählt über den **ganzen** Baum
  `.harness/baseline/v5.11.0/` mit `grep -rniE "modul[- ]?pfad|pfad[- ]?modul"`
  (deckt „Modulpfad", „Modul Pfad", Beugungen): exakt die zwei genannten
  Treffer, kein dritter. Das §5-Risiko „Punkt 2 könnte sich als Lesefehler
  erweisen" tritt **nicht** ein — zusätzlich geprüft wurden `modul-04`,
  `modul-09:71-74`, `grundlagen-referenz-richtung.md`,
  `grundlagen-begriffe.md` und `templates/README.md`: keine Auflösung des
  Begriffs. Der unabhängige Vorgänger-Review
  (`2026-08-23-slice-136-agents-34-klaerung-review.md` F-3) kam auf demselben
  Weg zum selben Ergebnis („die Mehrdeutigkeit bleibt damit im Kanon offen").
- geprüft, ohne Befund: **Punkt 2, Aussage über `architecture.template.md`** —
  die Vorlage vollständig gelesen: §1 vergibt `ARC-001` … `ARC-006` an Rollen
  (`Types / Domain`, `Config-Layer`, …), §2 führt dieselben Kennungen mit
  Schicht-Namen, §3–§5 nennen keine Code-Pfade. Die im CR mitgelieferte Grenze
  („aber eine Vorlage ist kein Regeltext") ist genau die Einschränkung, deren
  Fehlen der Vorgänger-Review als F-3 gemeldet hatte — sie ist hier gezogen.
- geprüft, ohne Befund: **Punkt 3, `modul-13`-Zitat und „sechster Schritt"** —
  der Satz ist gegen `modul-13-quality-gates.md:162-164` **wortgleich**; die
  Zählung stimmt: die Kette nennt sechs Schritte, von denen die
  sprachkonkrete Implementierung ausdrücklich nicht dargestellt wird, und
  *Bewusstes Brechen* ist der letzte.
- geprüft, ohne Befund: **Punkt 3, Instanzen einzeln** — alle drei sind
  verschieden und am Bestand auffindbar (`slice-132:167`, `slice-134:145-146`,
  `slice-138:177-181`, gebündelt in
  `docs/plan/planning/done/welle-84-results.md:51-56`). Beanstandet ist allein
  ihre Datierung (F-1), nicht ihre Existenz oder ihre Unterscheidbarkeit.
- geprüft, ohne Befund: **Punkt 4, die drei Formen gegen die drei Belege** —
  jede der drei genannten Formen hat eine eigene Instanz: Form 1 (Satz gilt
  nur für einen benannten Fall) in
  `docs/plan/planning/done/slice-127-claude-md-pointer.md:268-276`, Form 2
  (`MR-`Eintrag nach Titel statt nach `Geltungsbereich`) in
  `docs/plan/planning/done/slice-131-reviewer-skill-waisen.md:192-199`, Form 3
  (ADR-Akt als stehendes Verbot) in
  `docs/plan/planning/done/slice-130-lastenheft-historie-form.md:159-171`. Die
  Aussage „alle an eigenem Bestand beobachtet" und „die Klasse hat hier die
  Zähl-Schwelle erreicht" deckt sich mit `BEO-012` (Zähler 3, Schwelle
  erreicht).
- geprüft, ohne Befund: **„Was der CR nicht verlangt" je Punkt** — vorhanden
  bei allen vier (CR:67, :102, :138, :171); die Kopf-Zusage „Drei bitten um
  einen geschärften Satz, einer um eine Klärung" deckt sich mit den vier
  Schluss-Absätzen (Punkt 2 ist die Klärung).
- geprüft, ohne Befund: **Unabhängigkeit der vier Punkte** (§5 Risiko 2) —
  keiner setzt einen anderen voraus: jeder nennt seine eigene Quelle, seine
  eigenen Belege und seine eigene Bitte; die Streichung eines beliebigen
  Punktes lässt die übrigen drei unverändert lesbar.
- geprüft, ohne Befund: **Zustandsfeld-Form (`AGENTS.md` §3.7)** — der CR
  trägt im Kopf nur `Absender`, `Datum`, `Gegenstand` und **kein**
  `Status`-Feld; `MR-035` trägt `Status: Accepted` in der Bestandsform aller
  `MR-`Einträge. Keine Chronik in einem Zustandsfeld, keine Herkunfts-Prosa in
  einem Kommentar (der Commit fasst weder Code noch Konfiguration noch ein
  Skript an).
- geprüft, ohne Befund: **Referenz-Richtung (SDP)** — der CR nennt keine
  einzige repo-eigene Kennung (`ADR-`/`MR-`/`DC-`/`slice-`/`welle-`); das
  einzige ID-artige Token ist `ADR-<NNNN>` **innerhalb** des `modul-13`-Zitats.
  Kein Abwärts-Zeiger, kein Provenance-Marker, also auch keine
  Marker-Ehrlichkeits-Frage. `MR-035` verweist nur aufwärts auf die Baseline.
- geprüft, ohne Befund: **Links und Anker** — der CR enthält **keinen**
  Markdown-Link und keinen Anker (`grep -nE '\]\(|https?://'` leer), damit auch
  keine positionsabhängige Auflösung; `docs/plan/cr/` steht nicht in
  `links.resolve-from` und ist auch kein Lifecycle-Verzeichnis, die
  `link-position-dependent`-Falle dieses Tages ist hier gegenstandslos. Der
  einzige Link in `MR-035` (`../../.harness/baseline/v5.11.0/regelwerk/grundlagen-harness-dateien.md`)
  löst von `harness/conventions/` korrekt auf und ist pin-gebunden (`MR-021`);
  die neue Index-Zeile in `harness/conventions.md` trägt beide Anker-Formen
  (Voll-Slug plus Kurzform) und ihr Ziel `conventions/MR-035-cr-ablage.md`
  existiert.
- geprüft, ohne Befund: **Ablage-Form gegen die eigene Setzung** — `MR-035`
  verlangt `<YYYY-MM-DD>-cr-<gegenstand>.md` unter `docs/plan/cr/`; die Datei
  heißt nach `070be0e` `2026-08-25-cr-regelwerk-v5110.md` und ist die einzige
  im Verzeichnis. Die Umbenennung ist ein `git mv` mit Similarity 99 %, das
  Kopf-Datum wurde im selben Commit nachgezogen — keine hängenden Verweise
  (der Slice-Plan verweist nicht auf den Dateinamen).
- geprüft, ohne Befund: **`Gegenstand`-Kopfzeile ohne `modul-06`** — die
  Aufzählung nennt die vier Dateien, an denen der CR eine Änderung erbittet;
  `modul-06` wird in Punkt 1 nur als **Präzedenz** zitiert und soll unverändert
  bleiben. Die Auslassung ist damit konsistent, nicht unvollständig (der
  Slice-Plan führt `modul-06` in seinem `Bezug`-Feld dagegen zu Recht mit).
- geprüft, ohne Befund: **`MR-035` gegen den Konventionsspeicher-Rang** — die
  Verzeichnisfrage ist tatsächlich eine Form-Frage, die
  `grundlagen-source-precedence.md:61-71` ausdrücklich an den
  Konventionsspeicher abtritt; der Wert `keine` im Feld
  `Ersetzt-Baseline-Regel` hat Bestandspräzedenz (`MR-026`, `MR-028`,
  `MR-029`, `MR-030`, `MR-034`). Die Kanon-Stelle *„kein `CR-*`-ID-Schema,
  keine eigene Datei, kein Gate"* ist ausweislich ihres eigenen Satzbaus auf
  den **eingehenden** CR beschränkt, und `MR-035` zieht diese Grenze selbst —
  beanstandet ist nur die zweite, unbelegte Hälfte der Prämisse (F-5).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 3 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:**
dichte-behauptung-ueber-einen-tag-durch-eigene-historie-widerlegt ·
behauptete-kanon-luecke-in-der-zitierten-datei-bereits-geschlossen ·
zitat-reihenfolge-und-wortlaut-weichen-von-der-quelle-ab ·
probenzahl-mischt-regelverstoss-und-fail-closed-probe ·
unbelegte-kanon-zuschreibung-in-adaptions-eintrag ·
botschaft-sagt-belegform-zu-die-ein-punkt-nicht-traegt

## Verdikt

**Merge-blockierend:** ja — drei MEDIUM. Kein HIGH: der CR erfindet keine
Quelle, keine Fundstelle und keine Messung; die vier Punkte tragen einzeln,
Punkt 2 hält der schärfsten Prüfung stand (die Erschöpfungs-Behauptung ist
unabhängig reproduziert), und die Zahlen von Punkt 1 sind am richtigen
Zeitpunkt korrekt.

**Der Schwerpunkt liegt bei der Klasse, die der CR selbst zum Thema macht.**
Drei der sechs Findings (F-2, F-3, F-5) sind `BEO-012`-Klasse — eine Quelle
wird für mehr oder für weniger in Anspruch genommen, als sie sagt —, und sie
sitzen in Punkt 1, Punkt 4 und im begleitenden Adaptions-Eintrag. Nach der
Kontext-Eskalations-Regel des Skills ist die dritte Wiederholung derselben
Klasse in einer Sitzung ein Steering-Loop-Signal; `BEO-012` steht bereits auf
erreichter Schwelle und trägt im Reviewer-Skill weiterhin **keinen** Anker.
F-1 ist die teuerste Einzelstelle, weil sie in einem Dokument steht, das
dieses Repo nach außen gibt, und weil die eigene Datums-Berichtigung
(`070be0e`) sie im selben Zug prüfbar gemacht hat, ohne sie anzufassen.

**Übergabe:** Findings gehen an den Implementer. F-2 und F-3 sind Kandidaten
für weitere Instanzen von `BEO-012`, F-1 und F-4 für `BEO-009`; die Einordnung
und die Frage, ob der Reviewer-Skill für `BEO-012` einen Anker bekommt,
obliegen dem Maintainer bei der Slice-Closure (§7), nicht diesem Report.
Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation (DoD-/
Plan-Konformität prüft der Verifier separat).
