# Review-Report: slice-195 — 2026-09-03

**Review-Art:** Code — geprüft gegen Plan, ADR und Konventionen (Modul 10
§Drei Review-Arten).

**Gegenstand:** Commit `4d3a386` (feat(planning): slice-195 --
Beobachtungs-Register-Datenmigration Tabelle -> Verzeichnis, slice-195,
welle-88, ADR-0083, MR-013)

**Nachtrag:** `4d3a386` ist nach diesem Report in zwei Commits zerlegt
worden (`94b19bd` Beanspruchung, `b1b960b` Migrationsinhalt) und trägt
die Behebung der beiden HIGH-Befunde (F-1, F-2) — siehe slice-195 §9.

**Skill:** `.harness/skills/reviewer.md`
**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-03

**Eingangs-Kontext:**

- [slice-195](../plan/planning/done/slice-195-beobachtungsregister-migration.md)
- `AGENTS.md` §3.3 (git mv + Inhaltsänderung = zwei Commits, inkl. der vier
  benannten Ausnahmen), §3.5 (ADR-Immutabilität), §3.7 (Kommentarklassen),
  §3.4 (Spec-Straten-Abwärtssperre)
- [`MR-013`](../../harness/conventions/MR-013-lifecycle-move-buendelung.md)
  (Lifecycle-Move-Bündelung), [`MR-059`](../../harness/conventions.md#mr-059--wellen-archiv-stub-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)
- `.harness/baseline/v6.0.0/regelwerk/modul-06-roadmap.md`
  §Das Beobachtungs-Register
- `git show HEAD~1:docs/plan/planning/observations.md` (Alt-Tabelle, als
  Vergleichsgrundlage für die Datenmigration)

---

## Findings

### F-1 — Drei Accepted-ADRs bekommen inhaltliche Änderungen an bereits
bestehenden `## Geschichte`-Zeilen statt nur Anhängen

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.5 ("erlaubt bleiben `## Geschichte`-Anhänge") /
  [`MR-013`](../../harness/conventions/MR-013-lifecycle-move-buendelung.md)
  (siehe F-2, dieselbe Quelle demonstriert die etablierte Alternative) /
  eigene Präzedenz derselben `.d-check.yml` (ADR-0043-Tombstone-Kommentar:
  "die ADR ist Accepted und immutable … den Verweis nachzuziehen würde den
  Beleg verfälschen")
- `pfad`: `docs/plan/adr/0074-offene-tasks-auf-rohen-zeilen.md:165,170`,
  `docs/plan/adr/0075-erklaerte-teilmenge-in-structure.md:217,219`,
  `docs/plan/adr/0078-erklaerte-leermenge-mit-zahl.md:189,191,193`
- `befund`: Alle drei ADRs führen `**Status:** Accepted` und bereits vor
  diesem Commit bestehende `## Geschichte`-Zeilen (datiert `2026-08-30`),
  die je eine `[`BEO-<NNN>`](../planning/observations.md)`-Zitierform
  trugen. Dieser Commit hat den Zitierteil dieser bereits vorhandenen
  Zeilen auf die neue Pfad-Form umgeschrieben (`[`BEO-ALL/…`](../planning/observations/BEO-ALL/…/observation.md)`)
  — das ist kein Anhängen einer neuen Zeile, sondern ein Body-Edit an
  Text, der vor diesem Commit bereits committet war. Innerhalb derselben
  Dateien blieben BEO-Zitate im **Kern** (vor `## Geschichte`, z. B.
  0074:59 `[`BEO-003`](../planning/observations.md)`) unangetastet — der
  Implementer selbst hat diesen Unterschied getroffen (Commit-Botschaft:
  "eine erste blinde sed-Ersetzung traf 14 ADR-Kerne und wurde … zurückgerollt,
  danach gezielt nur innerhalb der jeweiligen `## Geschichte`-Sektion
  wiederholt"), weil `adr-check`s Core-Hash `## Geschichte` per
  `exclude-sections` **komplett** ausklammert (`internal/hexagon/core/rules/vcs.go:134`
  `vcsCore`), nicht nur auf Anhänge prüft.
  Das ist derselbe Fehlschluss, den §3.1 an anderer Stelle ausdrücklich
  benennt: eine Regel gilt unabhängig davon, ob ihr Wächter sie durchsetzt
  ("Wer sich auf den Wächter verlässt, verlässt sich auf nichts"). Dass der
  Hash `## Geschichte` nicht mitzählt, heißt nicht, dass ihr Inhalt änderbar
  ist — AGENTS.md §3.5 nennt als Ausnahme ausdrücklich nur "Anhänge", nicht
  Edits an bereits bestehenden Zeilen. Die im selben `.d-check.yml`
  dokumentierte Präzedenz für exakt diesen Fall (ADR-0043 zitiert einen
  inzwischen archivierten Review-Report; Verweis **nicht** nachgezogen,
  stattdessen `ignore-refs`-Tombstone) wurde hier nicht angewendet, obwohl
  dieselbe Tombstone-Mechanik (`docs/plan/adr/**` gegen
  `docs/plan/planning/observations.md`) in diesem Commit ohnehin schon für
  die unangetasteten Kern-Zitate angelegt wurde und dieselbe Deckung auch
  für die drei Geschichte-Zeilen getragen hätte.
- `verifizierbar`: ja — `make adr-check RANGE=HEAD~1..HEAD` meldet **0
  Befunde** (selbst ausgeführt, siehe unten), weil der Core-Hash
  `## Geschichte` vollständig ausklammert; das ist die gemessene Bestätigung,
  dass der Wächter diese Klasse strukturell nicht sieht. `git diff
  HEAD~1 HEAD -- docs/plan/adr/0074-*.md docs/plan/adr/0075-*.md
  docs/plan/adr/0078-*.md` zeigt die editierten Bestandszeilen.
- `klasse`: Inhaltsänderung an bereits committetem Text einer `Accepted`-ADR,
  gedeckt nur durch eine Gate-Lücke (vollständige `exclude-sections`-
  Ausklammerung statt Append-only-Prüfung), nicht durch eine dokumentierte
  Ausnahme.

### F-2 — Der Beanspruchungs-Commit (`open/` → `in-progress/`) bündelt die
gesamte Slice-Ausführung statt nur Move + Roadmap-Flip; Rename-Detection
fällt unter die 50-%-Schwelle

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.3 / [`MR-013`](../../harness/conventions/MR-013-lifecycle-move-buendelung.md)
  §Ausnahme Beanspruchung — Wortlaut: "der Slice-**Body** (DoD-Haken +
  Closure-Notiz …) bleibt im **zweiten** Commit; weil die Slice-Datei im
  Move-Commit unverändert ist, hält die Rename-Detection (`R100`) und damit
  die `git log --follow`-Begründung des Baseline-§3.3."
- `pfad`: Commit `4d3a386` selbst; `docs/plan/planning/in-progress/slice-195-beobachtungsregister-migration.md`
  (neu, 228 Zeilen) vs. `docs/plan/planning/open/slice-195-beobachtungsregister-migration.md`
  (gelöscht, 156 Zeilen, Vorgänger aus Commit `8c2080e`)
- `befund`: Dieser eine Commit ist zugleich (a) die Beanspruchung
  (`open/` → `in-progress/`, Roadmap-Ruhe-Marker entfernt), (b) die
  vollständige Migrations-Implementierung (~180 Dateien), (c) das Ausfüllen
  aller DoD-Haken in §4 und (d) das Schreiben der Closure-Notiz §9 — vier
  Dinge, die laut MR-013 explizit in **zwei** getrennten Commits laufen
  müssen (Move-Commit trägt nur Roadmap-Flip + Pfad-Verweise; Slice-Body
  bleibt im zweiten Commit). Die beiden Schwester-Slices derselben Welle
  (slice-193, slice-194) haben genau dieses Muster korrekt eingehalten:
  je ein separater `… beansprucht`-Commit (`62c5a06`, `155dc9b`) vor dem
  `feat(...)`-Inhalt und einem eigenen `… Closure-Body`-Commit
  (`f5fa899`, `5370c60`).
  Gemessene Folge: `git show HEAD~1 HEAD --find-renames=30%` zeigt die
  Slice-Datei nur bei **45 %** Similarity als Rename — unter der
  50-%-Standardschwelle. `git log --follow --oneline -- docs/plan/planning/in-progress/slice-195-beobachtungsregister-migration.md`
  zeigt nur `4d3a386` und **nicht** den Anlage-Commit `8c2080e` — exakt der
  Effekt, den §3.3/MR-013 mit der Zwei-Commit-Regel verhindern will.
  Die Commit-Botschaft rechtfertigt das Bündeln mit einer Analogie zu
  [`MR-059`](../../harness/conventions.md#mr-059--wellen-archiv-stub-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)
  (Wellen-Archiv-Stub-Move). Diese Analogie trägt nicht: MR-059 gilt einem
  mechanischen Werkzeug-Akt (Volltext wird im selben Akt durch einen Stub
  ersetzt, der ihn verschiebt — es gibt dort keine Phase mit unverändertem
  Inhalt). Die Beanspruchung eines Slice-Plans hat dagegen eine etablierte,
  funktionierende Zwei-Commit-Form (siehe slice-193/194 in derselben
  Welle), und MR-013 nennt sie ausdrücklich als das, was **nicht** gebündelt
  werden darf.
- `verifizierbar`: ja — beide Befehle oben sind selbst ausgeführt und
  reproduzierbar ohne Netz.
- `klasse`: Commit-Granularität widerspricht einer explizit zitierten
  Konvention (MR-013), mit gemessenem Schaden (Rename-Detection/`--follow`
  gebrochen) — kein Gate deckt das (§3.3: "Die Commit-Zerlegung selbst …
  sieht kein Gate").

## Negativbefunde

- geprüft, ohne Befund: **Daten-Treue der Migration.** Acht Einträge
  stichprobenartig gegen `git show HEAD~1:docs/plan/planning/observations.md`
  geprüft — `chronological-table-silently-reverses` (BEO-005, gestrichen,
  1 Beleg), `registry-vs-authority-table-drift` (BEO-001, gestrichen,
  2 Belege), `eigene-menge-gemessen-fremde-behauptet` (BEO-020, Zähler 6,
  6 Belege inkl. der im Fließtext versteckten sechsten Instanz slice-172),
  `invented-fourth-closure-outcome` (BEO-015, Zähler 3), `pin-bump-mirrors-ungated`
  (BEO-008, Zähler 5), `registerzeile-ohne-ausgang-nach-schwelle` (BEO-027,
  Zähler 1), `wortlaut-behauptet-pruefung-die-fehlt` (BEO-023, Zähler 6),
  `check-latest-blind-before-pin` (BEO-028, Zähler 1, Sub-Area `tools/harness/`,
  Kürzel `HARN`). In allen acht Fällen: Zähler-Summe exakt gegen die
  `evidence/`-Dateizahl, jedes Vorgangs-Beleg-Fund korrekt zugeordnet,
  `state.md`s Stand-Zusammenfassung inhaltlich treu zur alten Stand-Zelle
  (inkl. der komplexen Mehr-Instanzen-Differenzierung bei
  `eigene-menge-gemessen-fremde-behauptet`, wo zwei Instanzen korrekt als
  "kein formgültiger Ausgang" an `registerzeile-ohne-ausgang-nach-schwelle`
  weiterverwiesen werden statt hier einen erfundenen Ausgang zu tragen — das
  spiegelt exakt den alten `BEO-027`-Verweis "Beide bei `BEO-027` mitgezählt
  statt hier erfunden"). Keine stille Bedeutungsverschiebung gefunden.
- geprüft, ohne Befund: **Verzeichnis-/Datei-Zahlen.** 29 Beobachtungs-
  Verzeichnisse vorhanden (28 migrierte + 1 neu angelegte
  `mechanical-id-rewrite-misses-frozen-classes`, die Migration selbst
  betreffend) — deckt sich mit der Commit-Botschaft. 27× `BEO-ALL`, 1×
  `BEO-HARN`, exakt wie behauptet.
- geprüft, ohne Befund: **Citation-Rewrite-Auflösung.** Stichprobe über
  `AGENTS.md`, `docs/plan/adr/README.md`, sechs `harness/conventions/MR-*.md`-
  Dateien, `.harness/skills/reviewer.md`, `.harness/skills/closure-note-reviewer.md`
  und `docs/plan/adr/0055-*.md` (Proposed, Kern-Edit erlaubt) — jeder neue
  relative Link löst per `realpath -m` auf eine existierende
  `observation.md` auf.
- geprüft, ohne Befund: **Frozen-Klassen korrekt übersprungen.**
  Repo-weiter `grep` nach verbliebenen `BEO-[0-9]+`-Vorkommen außerhalb der
  fünf genannten Frozen-Klassen findet nur: (a) den absichtlich fiktiven
  Beispiel-`BEO-999` in `spec/lastenheft.md` (Boundary-Beispiel ohne
  Registerzeile, seit ADR-0079 so gedacht), (b) Fixture-Literale in
  `internal/hexagon/core/rules/planning_observations_test.go`, (c) die
  eigene historische Prosa des Slice-Plans selbst (beschreibt die alten
  IDs als Ist-Zustand vor der Migration). Kein lebendes Zitat übrig.
- geprüft, ohne Befund: **§3.4 (Spec-Straten nie abwärts).** Die editierte
  `spec/lastenheft.md`-Zeile (0.63.0-Historie) ersetzt nur eine
  BEO-Zitierform, referenziert kein ADR/Slice/Welle neu; `make doc-check`
  (Modul `matrix`) lief im vollen Gate-Durchlauf grün.
- geprüft, ohne Befund: **§3.1 (kein Host-Toolchain-Leck).** Die
  Commit-Botschaft nennt `sed` für die mechanische Ersetzung — `sed` wird
  bereits von mehreren Repo-eigenen Gate-Skripten (`tools/harness/*.sh`,
  Makefile) gerufen und fällt damit unter die erlaubte
  POSIX-Standardwerkzeug-Klasse aus §3.1.
- geprüft, ohne Befund: **`.d-check.yml`/`.d-check.closure.yml`-
  Wohlgeformtheit.** Beide Dateien wurden von `d-check` selbst im
  gesamten Gate-Lauf (`make gates`, `make fullbuild`) fehlerfrei geparst;
  der neue `dir`-Schlüssel ersetzt additiv `register`, ohne dass eine
  Konfigurations-Kollision auftritt.
- geprüft, ohne Befund: **Tombstone-Scope.** Die fünf neuen
  `ignore-refs`-Einträge (`done/**`, `docs/reviews/**`, `docs/plan/adr/**`,
  `harness/conventions/done/**`, `docs/plan/cr/**`, alle gegen genau
  `docs/plan/planning/observations.md`) sind exakt scope-parallel zum
  bestehenden `v5.18.0`-Tombstone und nicht breiter formuliert. Kontrolle:
  kein Proposed-ADR trägt noch eine alte `observations.md`-Referenz, die
  der breite `docs/plan/adr/**`-Scope fälschlich stumm geschaltet hätte.
- geprüft, ohne Befund: **`make gates` und `make fullbuild`, selbst
  ausgeführt.** `make gates`: "[gates] baseline-verify + workflow-pins +
  doc-check + lint + test + arch-check + coverage-gate + semgrep +
  gate-consistency + planning-check green" (zehn Gates, Coverage 94,70 % ≥
  93 %, semgrep 0 Findings). `make fullbuild`: zusätzlich `image-test` OK,
  `bench` OK (Median 833 ms < 5000 ms), `--trace --require-complete`
  meldet 51 Anforderungen, 0 Waisen, und der Closure-Bindepunkt
  (`--config .d-check.closure.yml --enable planning --enable structure
  --enable spans --enable reviews`) meldet **581 Datei(en) geprüft, 0
  Befund(e)** — exakt die im Commit behaupteten Zahlen. `[fullbuild] green
  — image-hash sha256:60c3fc51…`.
- geprüft, ohne Befund: **DoD/Closure-Notiz-Ehrlichkeit (mechanischer
  Teil).** Alle mechanisch nachprüfbaren Zahlen der Closure-Notiz (28
  Verzeichnisse, 88 Belegdateien, zehn Gates, 581-Dateien-Lauf) stimmen mit
  der beobachteten Realität überein — keine Überziehung der gemessenen
  Menge im Sinne von §5/AGENTS.md ("eine Commit-Botschaft … behauptet
  nicht mehr, als die Arbeit trägt").

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 1 |

### F-3 (MEDIUM) — Das Register-eigene `README.md` nennt nur drei der fünf
tatsächlich eingefrorenen Quellklassen

- `kategorie`: MEDIUM
- `quelle`: Konsistenz innerhalb desselben Commits — `.d-check.yml` trägt
  fünf Tombstone-Einträge (`done/`, `docs/reviews/`, `docs/plan/adr/`,
  `harness/conventions/done/`, `docs/plan/cr/`), die Closure-Notiz benennt
  dieselben fünf explizit als Fund
- `pfad`: `docs/plan/planning/observations/README.md` (letzter Absatz)
- `befund`: Der letzte Absatz des neuen Register-Index sagt: "Eingefrorene
  Bestände (`done/`, `docs/reviews/`, `harness/conventions/done/`) zitieren
  die Register-Form ihrer Zeit und werden nicht nachgezogen." — das sind
  die drei im **Slice-Plan** ursprünglich benannten Klassen, nicht die
  fünf, die dieselbe Ausführung als tatsächlich notwendig entdeckt und in
  `.d-check.yml` verankert hat (Accepted-ADR-Kerne, `docs/plan/cr/`).
  Ausgerechnet das Dokument, das als "Regeln dieses Registers"-Referenz
  gedacht ist, spiegelt die eigene Erkenntnis dieses Commits nicht wider.
- `verifizierbar`: ja — Textvergleich `docs/plan/planning/observations/README.md`
  gegen die fünf `ignore-refs`-Einträge in `.d-check.yml` und die
  Closure-Notiz §9 desselben Slice-Plans.
- `klasse`: Dokument-interne Inkonsistenz innerhalb desselben Commits —
  keine Daten-/Gate-Auswirkung, aber die künftige Nachschlagequelle für
  genau diese Frage ist unvollständig.

### F-4 (INFO) — `BEO-<NNN>`-artige Slice-Nummern in `.d-check.yml`-
Kommentaren (pre-existent, nicht neu eingeführt)

- `kategorie`: INFO
- `pfad`: `.d-check.yml:220` (neuer Kommentar: "Tombstone der
  Register-Migration … (slice-195, ADR-0083)")
- `befund`: §3.7 verbietet Slice-Nummern in Kommentaren außerhalb des
  Baseline-Schemas (`DC-*`, `ADR-*`, `MR-*`, `seit welle-<NN>`). Der neue
  Kommentar nennt `slice-195` direkt. Das ist jedoch **keine neue
  Abweichung dieses Slice** — dieselbe Form ("Tombstone … (slice-NNN /
  ADR-…)") ist in `.d-check.yml` bereits dutzendfach etabliert (z. B.
  Zeile 56, 107, 150, 157, 175 — alle vor diesem Commit). Wird hier nur
  benannt, nicht als eigenständiger Befund gegen dieses Slice gezählt.
- `verifizierbar`: ja — `grep -n "slice-[0-9]" .d-check.yml` zeigt den
  Bestand.
- `klasse`: Fortführung einer bestehenden, wahrscheinlich grandfathered
  Datei-Konvention — kein neues Risiko.

## Verdikt

**Merge-blockierend: ja (F-1 und F-2, beide HIGH).**

Die eigentliche Datenmigration ist ausgezeichnet ausgeführt: alle
stichprobenartig geprüften Einträge (8 von 29, inklusive der komplexesten
Mehrfach-Instanz-Fälle) sind ohne Bedeutungsverschiebung übertragen, die
Zähler-Summen stimmen exakt, alle geprüften Citation-Rewrites lösen aufs
richtige Ziel auf, und alle zehn Gates plus `fullbuild` laufen selbst
ausgeführt grün mit exakt den in der Commit-Botschaft behaupteten Zahlen
(28 Verzeichnisse, 88 Belege, 581-Dateien-Closure-Lauf). Die
Frozen-Klassen-Behandlung (Accepted-ADR-Kerne, `docs/plan/cr/`) ist in der
Sache richtig erkannt.

Aber zwei Hard-Rule-Verstöße bleiben, beide von keinem Gate gefangen (das
ist in diesem Repo kein Freibrief — §3.1: "Die Regel gilt unabhängig von
ihrer Durchsetzung"):

1. **F-1:** Drei bereits committete `## Geschichte`-Zeilen dreier
   `Accepted`-ADRs wurden inhaltlich verändert statt nur ergänzt — ein
   Edit an Accepted-ADR-Inhalt, den §3.5 nicht erlaubt und den die eigene,
   im selben Commit demonstrierte Alternative (Tombstone statt Nachziehen)
   sauber vermieden hätte.
2. **F-2:** Der Beanspruchungs-Commit bündelt Move, komplette
   Implementierung, DoD und Closure-Notiz entgegen dem ausdrücklichen
   Wortlaut von MR-013 (das dieser Slice selbst als Bezug zitiert) — mit
   gemessenem Schaden: `git log --follow` verliert die Historie zum
   `open/`-Anlage-Commit.

**Übergabe:** Beide Findings sind an dieser Stelle nicht mehr ohne
Git-Historien-Eingriff sauber korrigierbar (der Commit ist bereits auf dem
Hauptzweig). Empfehlung an Planner/Architect: F-1 als Folge-Slice oder
-Commit heilen (die drei Geschichte-Zeilen auf die alte Zitierform
zurücksetzen und stattdessen vom bereits vorhandenen `docs/plan/adr/**`-
Tombstone decken lassen, analog ADR-0043); F-2 als Lerneintrag im
Closure-Abschnitt dieses Slice ausdrücklich benennen (fehlt dort aktuell
vollständig — der Lerneintrag spricht nur von den Frozen-Klassen, nicht
von der Commit-Granularität), da MR-013 sonst beim nächsten
Beanspruchungs-Commit erneut stillschweigend unterlaufen werden kann. F-3
ist vor der Closure leicht behebbar (README um die zwei fehlenden Klassen
ergänzen). Dieser Report ist ein Lauf-Beleg und ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat, in
eigenem Kontext.
