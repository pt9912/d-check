# Review-Report R3: slice-103 — Geteilte Lexik, dritte Runde — 2026-08-16

**Review-Art:** dritte Runde (Code) — geprüft wird **nicht**, ob die Befunde der
beiden Vorrunden abgehakt sind, sondern ob die **Klasse geschlossen** ist: gibt
es im Produkt noch Stellen, die eine Lexik-Frage selbst und anders beantworten,
und trägt der neue Kopplungs-Test, was er verspricht? Nicht geprüft wird die
DoD-Abhakung (getrennter Kontext, Verifikation).

**Gegenstand:** Commit-Range `d2aaf90..4546a62` (acht Commits: Wellen-Eröffnung ·
Messung · Vertrag · Implementierung · Erst-Report · Heilung 1 · Re-Review ·
Heilung 2); im Besonderen der zweite Heilungs-Commit `4546a62`.
Arbeitsbaum-Stand `4546a62`.

**Skill:** `.harness/skills/reviewer.md` @ 1.4.0 ·
**Modell:** claude-opus-5[1m] · **Datum:** 2026-08-16

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`DC-FA-CITE-001`](../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in),
  [`DC-FA-VER-001`](../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
  [`DC-FA-PIN-001`](../../spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in),
  [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in),
  [`DC-FA-ANCH-001`](../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors),
  [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
  [`DC-FA-TGT-001`](../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in),
  [`DC-FA-MTX-001`](../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
  (Lastenheft-Fassung 0.58.2)
- §[`DC-FA-ANCH-001.a`](../../spec/spezifikation.md#dc-fa-anch-001a--github-slug-algorithmus),
  §[`DC-FA-ANCH-001.b`](../../spec/spezifikation.md#dc-fa-anch-001b--inline-html-anker),
  §[`DC-FA-CITE-001.a`](../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations)
  Schritt 2,
  §[`DC-FA-VER-001.a`](../../spec/spezifikation.md#dc-fa-ver-001a--versions-pin-konsistenz-versions)
  Schritt 1,
  §[`DC-FA-PIN-001.a`](../../spec/spezifikation.md#dc-fa-pin-001a--content-pin-gegen-inhaltlichen-drift-pins)
  Schritt 2,
  §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritte 3/4,
  §[`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
  Schritte 3/4,
  §[`DC-FA-TGT-001.a`](../../spec/spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets)
  Schritt 3,
  §[`DC-FA-MTX-001.a`](../../spec/spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung),
  §[`DC-FA-LINK-001.a`](../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
  Schritte 1/2
- [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  (Entscheidungen 1–4, Konsequenzen, Re-Evaluierungs-Trigger),
  [ADR-0019](../plan/adr/0019-versions-pin-fence-ausnahme.md),
  [ADR-0020](../plan/adr/0020-content-pin-fence-ausnahme.md),
  [ADR-0024](../plan/adr/0024-vcs-immutable-gate.md),
  [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md)
- [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md),
  [`AGENTS.md`](../../AGENTS.md) §3 Hard Rules,
  [`CLAUDE.md`](../../CLAUDE.md)
- Der [Erst-Report](2026-08-16-slice-103-geteilte-lexik-review.md) (L-1…L-13) und
  die [Re-Review](2026-08-16-slice-103-geteilte-lexik-re-review.md) (N-1…N-10) —
  hier als **Negativ-Liste**: was dort steht, gilt als bearbeitet, gesucht wurde
  das, was dort **nicht** steht
- Der [Slice-Plan](../plan/planning/in-progress/slice-103-geteilte-lexik-raender.md)
  (§4 Abnahme-Punkte, §4a Messung und Klassen-Abschluss, §5 Definition of Done),
  das [Wellendokument](../plan/planning/welle-74-geteilte-lexik-raender.md) und
  die Register-Zelle **BEO-003** in
  [observations.md](../plan/planning/observations.md)

**Läufe dieses Reviews.** Alle Fixtures liegen in einem Scratch-Verzeichnis
außerhalb des Repos; alle d-check-Läufe sind netzlos und read-only. Gefahren
wurden: `make build`, `make test` (Baseline plus drei Mutationsläufe), `make
gates` (grün, 374 Dateien, 0 Befunde), `make verify-closure-notes` (grün, 344/0)
sowie rund 40 Fixture-Läufe gegen zwei Images — den HEAD-Build und einen aus
Dateikopien rekonstruierten **Vor-Slice**-Build (Stand `d2aaf90`, sechs
Regel-Dateien zurückgeschrieben). Jede Mutation wurde **vor** dem Eingriff in das
Scratch-Verzeichnis kopiert und danach aus der Kopie zurückgeschrieben — **kein**
`git checkout`. Jede Mutations-Gegenprobe ist am **Exit-Code** von `make test`
bewertet worden, nicht an einem grep-Muster; ein erster Rückbau erzeugte einen
Compile-Fehler (Exit 2, `declared and not used`) und wurde als solcher verworfen
statt als Testlücke gezählt. **Der Arbeitsbaum ist wiederhergestellt:** `git
status --short` ist leer bis auf diesen Report, und der Neubau nach der letzten
Mutation liefert dieselbe Image-ID wie der erste Build vor dem ersten Eingriff.

---

## Findings

### R3-1 — Der Kopplungs-Test, der die Aufzählung ablösen soll, kennt den dritten Konsumenten nicht; `pins` lässt sich ohne einen einzigen roten Test wieder abkoppeln

- `kategorie`: HIGH
- `quelle`: [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 4 („Jeder reparierte Konsument bekommt einen Test, der die geteilte
  Antwort **an ihm** festnagelt"); das in dieser Range **neu** hinzugefügte
  Akzeptanzkriterium „Boundary (eine Anker-Antwort)" von
  [`DC-FA-PIN-001`](../../spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in)
  („when `pins` und `anchors` im selben Lauf aktiv sind, then bestimmen beide
  denselben Anker gleich"); Register-Zelle **BEO-003** („dieselbe Eingabe durch
  **alle** Konsumenten derselben Frage"); Reviewer-Skill §HIGH
  (Stilles-Grün-Pfad) und §MEDIUM (fehlende Negativtests bei neuem öffentlichen
  Vertrag)
- `pfad`: `internal/hexagon/core/rules/lexikon_kopplung_test.go:81`–`94`
  (`TestAnkerFrageHatEineAntwort` vergleicht `anchorsKenntAnker` gegen
  `versionsKenntAnker` — `pins` kommt in der Datei nicht vor) gegen
  `internal/hexagon/core/rules/pins.go:103` (`DecodeFragment`, in dieser Range
  neu)
- `befund`: Der Test ist als Verkörperung von BEO-003 deklariert und mit dem Satz
  begründet, eine Liste könne eine Stelle vergessen, eine Kopplung nicht — er
  vergisst selbst eine: von den **drei** reparierten Konsumenten fährt er zwei.
  `pins` teilt mit `versions` zwar `headingSection`, hat aber eine **eigene**
  Aufruf-Stelle der Fragment-Dekodierung; wird sie zurückgebaut, bleibt `make
  test` bei Exit 0, und derselbe gepinnte Link kippt von `link-stale` (Exit 1)
  auf 0 Befunde (Exit 0) — der Drift-Schutz fällt kommentarlos weg. Damit hat das
  neue Akzeptanzkriterium von `DC-FA-PIN-001` keine Assertion, und die
  Beleg-Form, die die widerlegte Aufzählung ersetzen soll, hat denselben Defekt
  wie die Aufzählung.
- `verifizierbar`: ja — Rückbau von `DecodeFragment` auf den rohen
  Fragment-String in `internal/hexagon/core/rules/pins.go:103` über eine
  Dateikopie: `make test` **Exit 0** (kein roter Test). Aus demselben Stand ein
  Image gebaut und gegen ein Fixture außerhalb des Repos gefahren (Zieldatei mit
  einem HTML-Anker, dessen Wert ein Leerzeichen trägt, Link mit prozent-kodiertem
  Fragment, dazu ein absichtlich falscher `dpin`-Hash): HEAD-Build Exit 1
  `link-stale`, Mutations-Build **Exit 0, 0 Befunde**. Gegenprobe zur Echtheit
  des Tests: derselbe Rückbau in `versions` (eigene Slug-Antwort statt der
  geteilten in `internal/hexagon/core/rules/versions.go:121`–`141`) macht `make
  test` rot, und zwar in `TestAnkerFrageHatEineAntwort/Duplikat-Slug` **und**
  `TestVersionsCurrentFromDuplikatSlug` — für den Konsumenten, den er kennt,
  funktioniert er also genau wie versprochen. Arbeitsbaum danach
  wiederhergestellt.
- `klasse`: „Kopplung als Beleg, ein Konsument nicht gekoppelt"

### R3-2 — Die Richtungs-Aussage ist zum dritten Mal geschlossen formuliert und unvollständig; drei belegte Auslassungen, und die einzige Stelle, die „offen" behauptet, ist eine Meta-Zeile

- `kategorie`: HIGH
- `quelle`: [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  §Konsequenzen („**Weniger:** … und ein `dpin` …" — zwei Stellen, abschließend);
  `CHANGELOG.md:23` („Sie **findet weniger** an **zwei** Stellen") und
  `CHANGELOG.md:39` („**Findet weniger** — und zwar Fehlmessungen");
  `spec/lastenheft.md:2650` (Historien-Zeile 0.58.0: „findet **weniger** an
  **zwei** konstruierbaren Stellen");
  `docs/plan/planning/in-progress/slice-103-geteilte-lexik-raender.md:174`–`177`
  (Definition of Done: „**SemVer: Minor** … und **weniger** an keiner Stelle") —
  gegen `spec/lastenheft.md:2648` (0.58.2: „Sie ist jetzt offen formuliert");
  [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus); Reviewer-Skill
  §HIGH (Stilles-Grün-Pfad) mit §Kontext-Eskalation (dritte Wiederholung
  derselben Klasse)
- `pfad`: `docs/plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md:90`–`99`,
  `CHANGELOG.md:23` und `:39`, `spec/lastenheft.md:2650`,
  `docs/plan/planning/in-progress/slice-103-geteilte-lexik-raender.md:176`
- `befund`: Keine der vier release-sichtbaren Flächen ist umformuliert worden —
  die ADR, das `CHANGELOG`, die Historien-Zeile 0.58.0 und die Definition of Done
  zählen ihre Richtungen weiterhin abschließend auf, die DoD sogar noch in der
  Fassung „weniger an **keiner** Stelle", die schon der Erst-Report widerlegt
  hat. Ich habe in einem Nachmittag **drei** weitere Stellen belegt, die dort
  fehlen: (a) die case-sensitive Anker-Auflösung nimmt `pins` einen Befund
  **ohne jede Ausgabe** (die Stelle, die 0.58.2 selbst benennt — nur eben nirgends
  in der Aufzählung); (b) das reparierte `planning` findet an einer
  konstruierbaren Stelle auch **mehr**, während der neue `CHANGELOG`-Eintrag nur
  „findet weniger" sagt; (c) die Duplikat-Slug-Vereinheitlichung ändert, welchen
  **Span** ein `-1`-Anker adressiert, wenn eine gleichnamige Überschrift mit der
  automatisch erzeugten Kennung konkurriert — mit einem neuen `version-stale` als
  Folge. Die einzige Stelle, an der die Aufzählung „offen" heißt, ist die
  Meta-Aussage einer Historien-Zeile über sich selbst.
- `verifizierbar`: ja — drei Fixtures außerhalb des Repos, je gegen den
  HEAD-Build und den rekonstruierten Vor-Slice-Build (`d2aaf90`). (a) Zieldatei
  mit einem HTML-Anker in Großschreibung, Referenz und `current-from` in
  Kleinschreibung: Vor-Slice `pins` **Exit 1** `link-stale`, HEAD `pins` **Exit
  0, 0 Befunde**; `versions` kippt am selben Fixture von Exit 0 auf Exit 2
  („Anker nicht auflösbar"). (b) Roadmap, deren kanonische Überschrift
  **ausschließlich** in einem Beispiel-Block steht: Vor-Slice **Exit 0**, HEAD
  **Exit 1** `planning-drift` (fail-closed, „kanonische Überschrift fehlt"). (c)
  Datei mit den Überschriften `Alt`, `Alt` und `Alt-1` und `current-from` auf den
  Duplikat-Anker: Vor-Slice **Exit 0**, HEAD **Exit 1** `version-stale` — die
  Adresse zeigt jetzt auf die zweite `Alt`-Sektion statt auf `Alt-1`.
- `klasse`: „Verhaltensänderung eines ausgelieferten Moduls ohne gedeckte Richtung"

### R3-3 — Das Modul `vcs` beantwortet „welche Zeile trägt den Status" roh, während `matrix` dieselbe Frage fence-bewusst beantwortet; die Folge ist ein stilles Grün im Immutabilitäts-Gate

- `kategorie`: HIGH
- `quelle`: [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 1 („Ein Konsument, der eine Lexik-Frage selbst beantwortet, ist ein
  Defekt — keine Variante… Die Reparatur ist die **Übernahme der vorhandenen
  Antwort**") und Entscheidung 3 (die `vcs`-Grenze sei „in **beiden**
  Ausprägungen" benannt);
  §[`DC-FA-MTX-001.a`](../../spec/spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung)
  gegen §[`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
  Schritte 3/4; [ADR-0024](../plan/adr/0024-vcs-immutable-gate.md);
  Reviewer-Skill §HIGH (Stilles-Grün-Pfad in einem Gate oder Gate-Skript)
- `pfad`: `internal/hexagon/core/rules/vcs.go:149`–`162`
  (`vcsHeadStatusLineNo`: `strings.HasPrefix(raw, "## ")` und
  `statusLine.MatchString(raw)` über **alle** Rohzeilen) und
  `internal/hexagon/core/rules/vcs.go:113`–`123` (`baseImmutable`: dasselbe für
  `immutable-when`) gegen `internal/hexagon/core/rules/matrix.go:436`–`445`
  (`statusOf`: „Beide Formen lesen nur Prosa-Zeilen — Inhalte in
  Fenced-Code-Blöcken sind keine Statuswerte") und gegen
  `internal/hexagon/core/rules/vcs.go:129`–`143` (`vcsCore`, dessen
  Abschnitts-Maske über `excludedRanges` **fence-bewusst** ist)
- `befund`: Dasselbe Modul zieht in derselben Funktion zwei verschiedene
  Lexik-Antworten: die Abschnitts-Maske ist fence-bewusst, die Status-Zeile und
  die Immutabilitäts-Entscheidung sind es nicht. Eine ADR, die ihre Kopfzeile als
  **Beispiel in einem Code-Block** zeigt — eine Vorlagen- oder Konventionsdatei
  ist der Normalfall dafür —, verschiebt damit die gestrippte Zeile: eine reale
  Änderung an dieser Beispielzeile ist aus dem Core herausgerechnet und passiert
  das Gate **mit Exit 0 und ohne Ausgabe**, obwohl die Datei `Accepted` ist. In
  der Gegenrichtung meldet dieselbe Konstruktion einen „unzulässigen
  Status-Übergang", der keiner ist, und ein `Proposed`-Dokument, das ein
  `Accepted`-Beispiel zeigt, gilt als immutabel und blockiert legitime
  Weiterarbeit. Der Fence ist in allen Fällen **balanciert** und die Datei liegt
  im Arbeitsbaum — der benannte Re-Evaluierungs-Trigger („`fence-unclosed` in
  einer `vcs.paths`-Datei") kann hier per Konstruktion nicht feuern, und die
  benannte Grenze beschreibt einen anderen Mechanismus. Die richtige Antwort
  existiert im selben Binary, also greift ausdrücklich das
  Erreichbarkeits-Kriterium, mit dem Entscheidung 3 den Verzicht begründet.
- `verifizierbar`: ja — vier git-Fixtures außerhalb des Repos mit der
  `vcs`-Konfiguration dieses Repos (`immutable-when` auf die `Accepted`-Kopfzeile,
  `status-line`, `head-allow`, `exclude-sections`), je zwei Commits, Lauf über
  `--enable vcs --range HEAD~1..HEAD`. (1) **Stilles Grün:** ADR mit einem
  Beispiel-Block, dessen Zeile die `status-line` trifft; geändert wird
  ausschließlich diese Zeile ⇒ **0 Befunde, Exit 0**. (2) **Kontrolle:**
  identische Datei, deren Beispielzeile die `status-line` **nicht** trifft,
  dieselbe Art Änderung ⇒ `core-drift-vcs`, Exit 1. (3) **Falsch-Rot:** dieselbe
  Bauform mit einer abweichenden Beispiel-Kopfzeile ⇒ „unzulässiger
  Status-Übergang einer immutablen Datei", obwohl die echte Kopfzeile `Accepted`
  ist. (4) **Falsch-Immutabel:** ein `Proposed`-Dokument mit einem
  `Accepted`-Beispiel im Block, echte Körper-Änderung ⇒ `core-drift-vcs`, Exit 1.
  Die Gegenprobe an derselben Frage in `matrix`: Beispiel-Kopfzeile im Block ⇒ 0
  Befunde; dieselbe Zeile außerhalb des Blocks ⇒ `matrix-inactive`, Exit 1.
  Bestandsmessung: alle **54** ADR-Dateien tragen genau **eine**
  `status-line`-Zeile (Zeile 3) — der Fall ist heute latent. `make adr-check`
  läuft im CI-Workflow und als `pre-commit`-Hook.
- `klasse`: „geteilte Lexik, vom Konsumenten selbst vorbereitet"

### R3-4 — Das Modul `planning` beantwortet „ist das die nächste H2" weiter roh — in genau der Funktion, die dieser Commit an die geteilte Zeilen-Menge gebunden hat

- `kategorie`: HIGH
- `quelle`: [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 2 („Wer … fragt ‚ist das eine **Überschrift**' … bekommt die
  geteilte Antwort");
  §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritt 4 (in dieser Range neu geschrieben) gegen
  §[`DC-FA-ANCH-001.a`](../../spec/spezifikation.md#dc-fa-anch-001a--github-slug-algorithmus);
  Reviewer-Skill §HIGH (Stilles-Grün-Pfad in einem Gate)
- `pfad`: `internal/hexagon/core/rules/planning.go:377`–`390`
  (`planningBlockHasMarker`, Zeile 382: `strings.HasPrefix(lines[i], "## ")`)
  gegen `internal/hexagon/core/rules/markdown.go:340`–`352` (`parseATXHeading`)
  und `internal/hexagon/core/rules/sections.go:19`–`57` (`FindSectionHeads` /
  `SectionEnd`, von der Closure-Fähigkeit **desselben Moduls** benutzt)
- `befund`: Der Commit bindet diese Funktion an `proseLineSet` und heilt damit
  die Fence-Hälfte der Frage — die **Erkennung** bleibt eine eigene: ein roher
  Präfix-Vergleich statt des geteilten ATX-Parsers. Drei Formen, die das Produkt
  an anderer Stelle als Überschrift liest, beenden den Aktiv-Block deshalb nicht:
  eine eingerückte H2, eine tab-getrennte H2 und eine H1 (die nach der geteilten
  Abschnitts-Semantik einen H2-Abschnitt sehr wohl beendet). Steht der
  Ruhe-Marker hinter einer solchen Zeile, gilt er als Teil des Aktiv-Blocks, und
  die Lifecycle-Verletzung wird **nicht gemeldet** — Exit 0 ohne Ausgabe, wo die
  Kontrolle Exit 1 liefert. Das ist dieselbe Datei, dieselbe Frage, zwei
  Antworten, wie in der Re-Review N-3, nur eine Ebene weiter; `planning-check`
  ist Bestandteil von `make gates`.
- `verifizierbar`: ja — vier Roadmap-Fixtures außerhalb des Repos, je mit
  `--enable planning`, Verzeichnis ohne Slice-Datei, Ruhe-Marker **hinter** der
  jeweiligen Trennzeile. Eingerückte H2 ⇒ **0 Befunde, Exit 0**; H1 ⇒ **0
  Befunde, Exit 0**; tab-getrennte H2 ⇒ **0 Befunde, Exit 0**; Kontrolle mit
  gewöhnlicher H2 ⇒ `planning-drift`, **Exit 1**. Gegenprobe zur Zwei-Antworten-
  Aussage: dieselbe eingerückte Überschrift als Anker-Ziel adressiert ⇒ `links` +
  `anchors` melden **0 Befunde**, das Produkt liest sie also als Überschrift.
- `klasse`: „geteilte Lexik, vom Konsumenten selbst vorbereitet"

### R3-5 — Das Modul `targets` beantwortet „ist das eine Tabellenzeile" roh; ein Beispiel im Code-Block dokumentiert damit ein Target, und `gate-undocumented` entfällt

- `kategorie`: MEDIUM
- `quelle`: §[`DC-FA-TGT-001.a`](../../spec/spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets)
  Schritt 3 („eine Tabellenzeile ist eine Zeile, deren **erstes Zeichen** ein
  Pipe ist … in Parität zu `grep -E '^\|'`") gegen die fence-bewusste
  Tabellen-Lexik in `internal/hexagon/core/app/trace_table.go`;
  [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 1; das [Wellendokument](../plan/planning/welle-74-geteilte-lexik-raender.md)
  §6 („Die Tabellen-Lexik … ist damit selbst ein geteilter Rand — aber ein
  **neuer**, nicht ein driftender"); Reviewer-Skill §MEDIUM (Konsistenz-Lücke
  zwischen Modulen derselben Eingabe-Klasse)
- `pfad`: `internal/hexagon/core/rules/targets.go:148`–`163` (Zeile 155:
  `strings.HasPrefix(line, "|")` über `splitLines`) gegen
  `internal/hexagon/core/app/trace_table.go:320`–`347` (`markdownTableLines`,
  Fence-Zustand über die geteilten Prädikate, „damit `trace` und `links` dasselbe
  Dokument sehen")
- `befund`: Eine Beispiel-Tabellenzeile **innerhalb** eines Code-Blocks der
  Autoritäts-Doku zählt für `targets` als Deklaration; jedes andere Modul des
  Produkts liest dieselbe Zeile gar nicht. Ein Makefile-Target, das nur in einem
  Beispiel-Block auftaucht, gilt damit als dokumentiert, und `gate-undocumented`
  bleibt aus — ein stilles Grün im Meta-Gate `gate-consistency`, das Teil von
  `make gates` ist; in der Gegenrichtung erzeugt dieselbe Zeile ohne
  Makefile-Regel ein `gate-phantom`, das keines ist. Gemeldet wird nicht die
  `grep`-Parität als solche (sie ist per Vertrag gewollt), sondern dass ihre
  Fence-Folge in keinem Vertragstext steht und dass die Aufzählung des
  Klassen-Abschlusses die Tabellen-Achse gar nicht kennt. Das Wellendokument
  nimmt die Tabellen-Lexik ausdrücklich mit der Begründung aus, sie sei ein
  neuer und kein driftender Rand — diese Begründung ist falsifiziert.
- `verifizierbar`: ja — Fixture außerhalb des Repos mit zwei Makefile-Regeln, von
  denen nur eine in der Autoritäts-Tabelle steht. Mit der zweiten als
  Beispiel-Tabellenzeile **in einem Code-Block**: `--enable targets` ⇒ **0
  Befunde, Exit 0**. Dieselbe Datei ohne den Beispiel-Block ⇒
  `gate-undocumented`, **Exit 1**. Bestandsmessung im eigenen Repo: in
  `AGENTS.md` und `harness/README.md` gibt es heute **null** solche Zeilen — der
  Fall ist latent.
- `klasse`: „geteilte Lexik, vom Konsumenten selbst vorbereitet"

### R3-6 — Die Historien-Zeile 0.58.2 behauptet eine Grenze, die nirgends steht; der Satz, den sie einschränken soll, ist unverändert und gegen einen Lauf falsch

- `kategorie`: MEDIUM
- `quelle`: `spec/lastenheft.md:2648` (0.58.2: „die permissive Anker-Menge ist bei
  ihrem neuen Konsumenten `versions` **nicht** folgenlos … — als Grenze benannt,
  nicht mitgenommen") gegen
  §[`DC-FA-ANCH-001.b`](../../spec/spezifikation.md#dc-fa-anch-001b--inline-html-anker)
  Schluss-Absatz („Eine zu groß geratene Anker-Menge führt nur dazu, dass das
  Modul schweigt (nie ein Falsch-Befund)"); [`AGENTS.md`](../../AGENTS.md) §2
  (das Lastenheft ist vertraglich abnahmebindend); N-5 der
  [Re-Review](2026-08-16-slice-103-geteilte-lexik-re-review.md); Reviewer-Skill
  §MEDIUM (Spec-Treue-Lücke einer Messmethode)
- `pfad`: `spec/spezifikation.md:943`–`946` (unverändert) und
  `spec/lastenheft.md:2648`; keine Grenz-Aussage in
  [`DC-FA-VER-001`](../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
  in [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  (§Re-Evaluierungs-Trigger trägt weiterhin drei Punkte, keiner davon dieser), im
  Handbuch oder im `CHANGELOG`
- `befund`: Der Diff berührt weder den Satz, der die permissive Anker-Menge mit
  „nie ein Falsch-Befund" begründet, noch fügt er irgendwo eine Grenze oder einen
  Re-Evaluierungs-Trigger hinzu — die Historien-Zeile schreibt eine Lieferung
  fest, die es nicht gibt. Der Satz bleibt damit als generelle Zusage stehen,
  obwohl sein neuer Konsument genau den ausgeschlossenen Fall erzeugt. Das
  Failure-Szenario ist der Leser, der die Historie als Delta-Register benutzt und
  die Grenze sucht, die dort als geliefert vermerkt ist.
- `verifizierbar`: ja — Volltextsuche über alle Nicht-Report-Markdown-Dateien
  nach den Zeichenfolgen „zu groß geratene", „Schweige-Charakter",
  „HTML-Kommentar" und „eingerückt" im Zusammenhang mit der Anker-Menge: genau
  **eine** Fundstelle, der unveränderte Satz in der Spezifikation. Verhalten
  nachgemessen: Fixture außerhalb des Repos mit einem gleichnamigen Anker in
  einer HTML-Kommentar-Zeile **vor** dem echten Anker ⇒ `--enable versions`
  meldet `version-stale`, **Exit 1**, obwohl der Pin die Version des echten
  Anker-Spans trägt; dieselbe Datei ohne die Kommentarzeile ⇒ 0 Befunde, Exit 0.
- `klasse`: „Rand behauptet mehr als der Lauf"

### R3-7 — §4a trägt die widerlegte Aufzählung weiter im Indikativ, direkt über dem Absatz, der sie widerlegt; die Definition of Done steht auf dem Stand vor beiden Heilungen

- `kategorie`: LOW
- `quelle`: `docs/plan/planning/in-progress/slice-103-geteilte-lexik-raender.md`
  §4a und §5; Reviewer-Skill §LOW (Doku-Drift, latente Wartungsfalle)
- `pfad`: `docs/plan/planning/in-progress/slice-103-geteilte-lexik-raender.md:134`–`136`
  („**Übrig bleiben genau zwei** Stellen, die eine Lexik-Frage **roh**
  beantworten") gegen `:146`–`162` (derselbe Abschnitt, „Die Aufzählung hat ein
  zweites Mal nicht getragen"), sowie `:174` („`make gates` grün (372/0)")
- `befund`: Der Abschnitt, der den Klassen-Abschluss belegen soll, behauptet in
  einem Satz weiterhin eine abgeschlossene Zahl und widerruft sie zwölf Zeilen
  später — wer den Beleg zitiert, zitiert je nach Absatz das Gegenteil. Die
  Definition of Done nennt zudem den Gate-Stand `372/0`, während der Arbeitsbaum
  `374/0` liefert, und wiederholt die schon zweimal widerlegte Richtungs-Zusage
  (siehe R3-2).
- `verifizierbar`: ja — Volltext von §4a und §5 gegen den eigenen Lauf `make
  gates` (Exit 0, 374 Dateien, 0 Befunde).
- `klasse`: „Rand auf der alten Fassung stehengeblieben"

### R3-8 — Die Nachzieh-Regel für geänderte Zusagen ist erneut nur zur Hälfte angewandt: Handbuch ja, beide READMEs nein

- `kategorie`: LOW
- `quelle`: `docs/user/releasing.md` §4 (die in dieser Range eingeführte Regel,
  dass eine **geänderte** Zusage mitten in einem bestehenden Abschnitt
  nachzuziehen ist); [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  (Spiegel „Nutzer-Doku"); N-8 der
  [Re-Review](2026-08-16-slice-103-geteilte-lexik-re-review.md); Reviewer-Skill
  §LOW (Doku-Drift)
- `pfad`: `README.de.md:60`–`67` und `README.md:61`–`68` (die Modul-Punkte zu
  `versions` und `pins`) gegen `docs/user/benutzerhandbuch.md:1281`–`1287` und
  `:1364`–`1367` (dort ist die Anker-Antwort nachgezogen)
- `befund`: Der Commit zieht die geänderte Anker-Zusage im Handbuch an **beiden**
  Anker-Konsumenten nach und lässt die Kurzform in beiden READMEs unverändert.
  Dort steht zu `versions` weiterhin nur „liest auch Fenced-Code" — genau die
  Formulierung, aus der nach L-5 der Widerspruch zwischen Pin-Ausnahme und
  Anker-Regel entsteht, und die die Anforderung selbst inzwischen mit „Sie gilt
  dem Pin, nicht dem Anker" auflöst. Bewusst LOW: der Stand ist unreleased, und
  das Repo zieht die READMEs planmäßig im Release-Prep nach.
- `verifizierbar`: ja — Volltext der vier Modul-Punkte gelesen; die
  Zeichenfolgen „Anker", „Duplikat" und „Inline-Code" kommen dort nicht vor. Der
  Diff der Range berührt in beiden READMEs nur den `citations`-Punkt.
- `klasse`: „Rand auf der alten Fassung stehengeblieben"

### R3-9 — Der Kopplungs-Test misst „Anker auflösbar" am Fehlen eines Fehlers, nicht am adressierten Span

- `kategorie`: INFO
- `quelle`: [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 4; Reviewer-Skill §INFO (dokumentationswürdige, aber
  undokumentierte Annahme)
- `pfad`: `internal/hexagon/core/rules/lexikon_kopplung_test.go:67`–`77`
  (`versionsKenntAnker` liefert `err == nil`)
- `befund`: Die Kopplung vergleicht zwei **Booleans**: kennt `anchors` den Anker,
  und läuft `versions` ohne Fehler. Eine Divergenz in der **Span-Grenze** — also
  darin, *welchen* Ausschnitt derselbe Anker adressiert — ist damit unsichtbar,
  obwohl genau das der Gegenstand von
  §[`DC-FA-VER-001.a`](../../spec/spezifikation.md#dc-fa-ver-001a--versions-pin-konsistenz-versions)
  Schritt 1 ist und obwohl R3-2 (c) zeigt, dass sich diese Grenze in dieser Range
  verschoben hat. Festgehalten als Reichweiten-Annahme, nicht als Defekt: für die
  zehn geführten Schreibweisen ist die Boolesche Kopplung ausreichend.
- `verifizierbar`: ja — Lesart des Testkörpers; die Span-Verschiebung aus R3-2
  (c) lässt `make test` grün und ist nur am Fixture sichtbar.
- `klasse`: „Kopplung misst weniger, als der Vertrag zusagt"

## Negativbefunde

- geprüft, ohne Befund: **`citations` — Absatz-Frage.** Absatz-Sammler,
  Blockquote-Sammler und Eintritts-Prüfung ziehen die Grenze über
  `fencedBlockBetween`, dasselbe Prädikat wie `proseParagraphs`. Randfall
  nachgefahren: Direktive in Zeile 1 (kleinster erreichbarer Index nach dem
  Entfernen der Wache) — kein Panik-Pfad. Die neue Zusage „Leerzeilen trennen
  nicht" stimmt mit dem Code überein (Leerzeilen werden vor der
  Kandidatenwahl übersprungen).
- geprüft, ohne Befund: **`anchors` / `versions` — HTML- und Slug-Hälfte.** Beide
  Konsumenten gehen jetzt durch `htmlAnchorLines`, `headingSlugsOrdered` und
  `DecodeFragment`; der Rückbau der Slug-Antwort in `versions` macht genau zwei
  Tests rot, davon den Kopplungs-Test. Ein prozent-Zeichen ohne gültige
  Kodierung (`#100%` auf `id="100%"`) fällt sauber auf den Rohwert zurück: 0
  Befunde, Exit 0.
- geprüft, ohne Befund: **`immutable`.** Der Pin-Marker wird auf den
  **vorverarbeiteten** Zeilen gesucht, die Abschnitts-Maske läuft über
  `excludedRanges` — beide Antworten geteilt. Der rohe Core ist die per
  [ADR-0020](../plan/adr/0020-content-pin-fence-ausnahme.md)-Familie gescopte
  andere **Frage**, keine andere Antwort.
- geprüft, ohne Befund: **`matrix`.** Klassen-Zuordnung, Status-Extraktion,
  Token-Referenzen und Sektions-Masken laufen sämtlich über `proseLines` bzw.
  `extractHeadingLines`; die Status-Antwort ist ausdrücklich fence-bewusst
  dokumentiert und am Fixture bestätigt (siehe R3-3).
- geprüft, ohne Befund: **`structure`.** Abschnitts-Kopf, Abschnitts-Ende und
  Bereinigung kommen aus `FindSectionHeads` / `SectionEnd` / `SectionProse`; die
  Bedingungen arbeiten auf dem bereinigten Text.
- geprüft, ohne Befund: **`ids`, `codepaths`, `hostpaths`, `diagrams`, `spans`,
  `links`, `external`.** Alle sieben beziehen ihre Zeilen-Menge aus `proseLines`
  bzw. `PreprocessMarkdown`; `codepaths` holt die Anker-Frage aus `AnchorSet`,
  `diagrams` liest die Fence-Innenseite über denselben Automaten.
- geprüft, ohne Befund: **`sources`, `tracked`, `pins` (Marker-Bindung).** Alle
  drei arbeiten auf den vorverarbeiteten Zeilen und greifen nur für die
  **Positions**-Prüfung auf die rohe Zeile zurück — das ist längenerhaltend und
  ausdrücklich begründet, keine zweite Antwort.
- geprüft, ohne Befund: **`commits`.** Die Raute-Behandlung dort betrifft
  git-Kommentarzeilen einer Commit-Message, nicht Markdown — andere
  Eingabe-Klasse, keine Lexik-Divergenz.
- geprüft, ohne Befund: **App-Schicht.** `trace_table.go` führt den einzigen
  weiteren Fence-Automaten und speist ihn aus den geteilten Prädikaten;
  `repair.go` liest die Zeilen über `PreprocessMarkdown` und benutzt die rohe
  Zeile nur als Schreibziel; `suggest.go` liest Kennungen ausschließlich über
  `ExtractHeadings` / `StripHeadingLinks`; `diagnose.go` enthält keine Lexik.
- geprüft, ohne Befund: **Der `--doctor`-Spiegel.** Der Klartext von
  `anchor-missing` ist an allen vier Fundstellen wortgleich (beide Emitter,
  `reasonTexts`, §4-Tabelle); `--doctor` über das eigene Repo läuft befundfrei.
- geprüft, ohne Befund: **Keine unterschlagene neue Vertragsfläche.** Der Diff
  führt keinen neuen Grund-Code, keinen Config-Schlüssel und kein Modul ein;
  §2-Schema, §4-Grund-Code-Tabelle und die `--print-config`-Vorlage bleiben zu
  Recht unverändert. `docs/user/operations.md` trägt nur Modul-Enumerationen (20
  Einträge, vollständig) und keine der geänderten Zusagen.
- geprüft, ohne Befund: **Widerspruchsfreiheit der neuen Formulierungen
  untereinander.** Anforderung, Algorithmus, Handbuch und beide READMEs sagen zur
  `citations`-Paarung jetzt dasselbe (Leerzeilen trennen nicht, ein
  Fenced-Block trennt in beiden Zweigen); die Akzeptanzkriterien von
  `DC-FA-CITE-001`, `DC-FA-VER-001` und `DC-FA-PIN-001` tragen die neuen
  Boundary-Fälle. Der einzige gefundene Widerspruch ist R3-6.
- geprüft, ohne Befund: **Hard Rules.** Kein Spec-Stratum nennt eine ADR-Kennung;
  [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) steht auf
  `Proposed`, ist im Index eingetragen und trägt den geforderten
  Re-Evaluierungs-Trigger-Abschnitt; keine Inline-Suppression, keine
  Gate-Lockerung, kein Netzzugriff außerhalb der Netz-Module; die Import-Richtung
  des Kerns ist unverändert (nur `model` und Ports).
- geprüft, ohne Befund: **Referenz-Richtung (SDP) und Marker-Ehrlichkeit.** Die
  Range führt keinen Provenance-Marker ein; die Slice-Kennung steht in der ADR
  ausschließlich im ausgenommenen Geschichts-Abschnitt. `matrix` ist im grünen
  Gate-Lauf aktiv.
- geprüft, ohne Befund: **Gate-Stand.** `make gates` Exit 0 (374 Dateien, 0
  Befunde: `doc-check`, `lint`, `test`, `arch-check`, `coverage-gate`, `semgrep`,
  `gate-consistency`, `planning-check`) und `make verify-closure-notes` Exit 0
  (344/0), beide nach der Wiederherstellung des Arbeitsbaums.
- geprüft, ohne Befund: **Arbeitsbaum.** Jede Mutation aus der Scratch-Kopie
  zurückgeschrieben, kein `git checkout`; `git status --short` ist am Ende leer
  (dieser Report ist die einzige neue Datei), und der Neubau liefert dieselbe
  Image-ID wie vor dem ersten Eingriff.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 4 |
| MEDIUM | 2 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Kopplung als Beleg, ein Konsument nicht
gekoppelt · Verhaltensänderung eines ausgelieferten Moduls ohne gedeckte Richtung
· geteilte Lexik, vom Konsumenten selbst vorbereitet · Rand behauptet mehr als der
Lauf · Rand auf der alten Fassung stehengeblieben · Kopplung misst weniger, als
der Vertrag zusagt.

**Wiederholungs-Signal.** Beide blockierenden Klassen der Vorrunden treten ein
drittes Mal auf, und zwar in genau den Artefakten, die sie schließen sollen: die
geteilte Antwort ist an zwei weiteren Modulen halb übernommen (`vcs`, `planning`
— jeweils Fence-Hälfte geheilt oder nie gestellt, Erkennung eigen), und die
Richtungs-Aufzählung ist zum dritten Mal geschlossen formuliert, diesmal
einschließlich der Definition of Done. Nach §Kontext-Eskalation des
Reviewer-Skills ist das kein Melde-, sondern ein Steering-Fall: die Antwort auf
die Aufzählung war ein Kopplungs-Test, und der Kopplungs-Test hat dieselbe Lücke
wie die Aufzählung (R3-1).

## Verdikt

**Merge-blockierend:** ja — vier HIGH und zwei MEDIUM.

**Was geliefert wurde, hält.** Die Anker-Frage ist für `anchors` und `versions`
jetzt wirklich **eine** Antwort — HTML-Zweig, Slug-Zweig samt Duplikat-Zähler und
Fragment-Dekodierung —, und der Kopplungs-Test fängt eine Wieder-Divergenz an
diesem Konsumenten sofort und punktgenau: der Rückbau der Slug-Antwort macht
genau zwei Tests rot, einer davon der Kopplungs-Test. Die zwei
`planning`-Falsch-Rot sind repariert und assertiert, die falsche Absatz-Zusage
ist in allen fünf Dokumenten korrigiert und stimmt jetzt mit dem Lauf überein,
die Akzeptanzkriterien der drei Anforderungen tragen die neuen Boundary-Fälle,
und die Spezifikation hat ihre eigene Historien-Zeile bekommen. Die Idee, eine
Aufzählung durch eine Kopplung zu ersetzen, ist die richtige — sie ist nur nicht
zu Ende geführt.

**Die Klasse ist nicht geschlossen.** An zwei Modulen, die in keinem der beiden
Vorgänger-Reports vorkommen, beantwortet das Produkt weiterhin eine Lexik-Frage
selbst und anders, und beide Male ist die Folge ein **stilles Grün in einem
Gate**: `vcs` liest die Status-Zeile und die Immutabilitäts-Bedingung roh,
während `matrix` dieselbe Frage fence-bewusst beantwortet und das ausdrücklich
begründet — eine echte Änderung an einer `Accepted`-ADR passiert `make adr-check`
mit Exit 0 ohne Ausgabe (R3-3). `planning` liest die Block-Grenze weiterhin als
rohen Präfix, in genau der Funktion, die dieser Commit angefasst hat — drei
Konstruktionen verlieren den `planning-drift` still (R3-4). `targets` beantwortet
die Tabellen-Frage roh, und ein Beispiel im Code-Block dokumentiert damit ein
Target (R3-5); das Wellendokument nimmt diese Achse mit der Begründung aus, sie
sei ein neuer und kein driftender Rand — sie driftet bereits.

**Der Beleg, der die Aufzählung ersetzen sollte, hat ihren Defekt geerbt.** Der
Kopplungs-Test fährt zwei der drei reparierten Konsumenten; `pins` fehlt, und
genau für `pins` hat derselbe Commit ein neues Akzeptanzkriterium geschrieben.
Der Rückbau seiner Fragment-Dekodierung ist ein Einzeiler, lässt `make test` bei
Exit 0 und kostet einen Konsumenten seinen Drift-Befund (R3-1). **Und die
Richtungs-Aussage ist zum dritten Mal geschlossen und unvollständig** — auf ADR,
`CHANGELOG`, Historien-Zeile 0.58.0 und Definition of Done; drei weitere
Verlust- bzw. Zugewinn-Stellen sind gegen den Vor-Slice-Build belegt, darunter
eine, die bei `pins` schweigt und die 0.58.2 zwar benennt, aber in keiner
Aufzählung nachträgt (R3-2).

**Release-Empfehlung: noch nicht releasen.** Die Einordnung als **Minor** bleibt
richtig — sie ist es sogar deutlicher als bisher dokumentiert, denn die Änderung
wirkt an mehr Stellen in beide Richtungen als die vier Ränder behaupten. Ihre
**Begründung** trägt nicht: die DoD-Zeile, aus der die SemVer-Einordnung
abgeleitet wird, steht wörtlich auf der Fassung, die der Erst-Report widerlegt
hat. Vor den Tag gehören R3-1 (sonst ist die Verkörperung von BEO-003 eine
Behauptung), R3-2 (sonst geht ein Konsument mit einer falschen Release-Notiz ins
Update) und R3-3 (ein stilles Grün im Immutabilitäts-Gate ist die
Harness-Integrität selbst). R3-4 liegt auf demselben Diff und derselben Zeile
Code. R3-5 und R3-6 sind Vertrags- und Entscheidungsarbeit; R3-7, R3-8 und R3-9
sind klein, R3-8 ist Release-Prep-Material.

**Übergabe:** Die Findings gehen an den Implementer. Die tragende Beobachtung
dieser Runde ist nicht ein einzelner Befund, sondern dass die **Form** des
Klassen-Abschlusses zum dritten Mal an derselben Kante gescheitert ist: erst eine
Aufzählung, die zwei Stellen nicht kannte, dann eine Aufzählung, die zwei weitere
nicht kannte, jetzt eine Kopplung, die einen von drei Konsumenten nicht kennt —
und daneben drei Module, die nie jemand gefragt hat. Für das Register heißt das:
die Verkörperung ist die richtige Antwort, aber sie ist noch nicht vollständig
verkörpert. Dieser Report ist ein Lauf-Beleg (diese Range, dieser Skill, dieses
Modell, dieses Verdikt) und ersetzt keine Verifikation — DoD- und
Plan-Konformität prüft der Verifier separat.
