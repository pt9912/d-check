# Review-Report: slice-200 — 2026-09-04

**Review-Art:** Code-Review (gegen Plan + Konventionen, Modul 10 §Drei Review-Arten) —
kein Plan-Review (der Plan liegt bereits vollständig vor, Gegenstand ist seine Umsetzung).

**Gegenstand:** Commit `ce5fb50ed7207c39d86dcbd8725ea9c7d611e07e`
(„feat(planning): slice-200 -- sieben verbliebene wellenlose Slices archiviert") —
sieben `tools/archive-wave -slice=<id> -apply`-Läufe (slice-141, slice-168,
slice-169, slice-170, slice-183, slice-184, slice-188) in einem Commit, plus
drei nachgezogene Referenz-Fixes (`AGENTS.md`, `welle-87-results.md`,
`MR-057-baseline-v5150.md`).

**Skill:** `.harness/skills/reviewer.md` @ 1.13.0
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-04

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-200-verbliebene-wellenlose-slices-archivieren.md` (vollständig)
- `AGENTS.md` §3.3 (git mv + Inhaltsänderung), §3.7 (Kommentar-Klassen)
- `harness/conventions.md` / `harness/conventions/MR-062-wellenloser-slice-archiv-move.md`,
  `harness/conventions/MR-063-eigenstaendiger-review-archiv-move.md`, `MR-059`, `MR-013`
- `tools/archive-wave/archive.go`, `main.go`, `slice_mode_test.go` (Verhalten von `ApplySlice`)
- `.d-check.closure.yml` §`reviews`
- Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur Schritt 4, §Wann Arbeit
  eine Welle braucht

---

## Findings

### F-1 — Commit bündelt sieben unabhängige MR-062-Akte ohne Geltungsbereichs-Deckung für die Bündelung

- `kategorie`: MEDIUM
- `quelle`: MR-062 (§Grenze), MR-063 (§Adaption/§Grenze)
- `pfad`: `docs/plan/planning/in-progress/slice-200-verbliebene-wellenlose-slices-archivieren.md:50-52`
  (§2 Punkt 5, „Commit-Granularität")
- `befund`: Der Plan begründet den einen gebündelten Commit für alle 7 Archivierungen mit
  „dieselbe Begründung wie bei slice-197/slice-199". `MR-062` selbst hält seine §Grenze
  jedoch ausdrücklich fest: „Ob mehrere unabhängige Einzel-Slice-Archive in einem
  gemeinsamen Commit gebündelt werden, ist eine andere Frage — diese Regel begründet nur,
  dass je Slice genau ein Commit die Untergrenze ist". `MR-063` beantwortet die
  Bündelungsfrage zwar positiv, aber **ausdrücklich nur für den Review-Modus** und nennt
  wörtlich: „diese Regel beantwortet sie hier nur für den Review-Modus, nicht rückwirkend
  für MR-062". Für den hier verwendeten Einzel-Slice-Modus (`ApplySlice()`, exakt MR-062s
  Geltungsbereich) bleibt die Bündelungsfrage damit **weiterhin formal offen** — der Plan
  zitiert eine Präzedenz (slice-197s faktische Praxis, vor MR-062/063 entstanden) für eine
  Aussage, die die zuständige, später geschriebene Quelle für genau diesen Geltungsbereich
  explizit nicht trifft. Das ist dieselbe Fehlerklasse, die MR-062 selbst korrigiert
  („genau dieser Fehler unterlief slice-197s eigenem Plan"): eine Zitat-Analogie über den
  erklärten Geltungsbereich der Quelle hinaus.
- `verifizierbar`: nein — kein Gate prüft Commit-Granularität oder Zitat-Geltungsbereiche;
  die Prüfung ist Urteil (Reviewer-Skill, Prüffrage #9).
- `klasse`: citation-stretched-beyond-scope (Bündelungsfrage MR-062)

**Einordnung:** Kein Verstoß gegen eine bestehende Verbots-Norm (MR-062 verbietet die
Bündelung nicht, sie lässt sie offen) und keine Wiederholung eines bereits gemeldeten,
unkorrigierten Fehlers — aber exakt die Lücke, die MR-063 für den Review-Modus bereits
geschlossen hat und für den Einzel-Slice-Modus ausdrücklich ungeschlossen ließ. Ein
Folge-Slice/eine Folge-MR (analog MR-063, aber für MR-062s Geltungsbereich) sollte die
Frage schließen, bevor ein weiterer Bündel-Commit im Einzel-Slice-Modus entsteht — sonst
bleibt jeder solche Commit auf derselben Lücke stehen.

### F-2 — Zwei der sieben archivierten Volltexte sind nicht byte-identisch mit dem committeten Original

- `kategorie`: MEDIUM
- `quelle`: Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur Schritt 4
  („unveränderliches Archiv"); Reviewer-Skill Prüffrage #2 (Korrektheit)
- `pfad`: `docs/plan/planning/done/wellenlos/slice-184-archiv.zip`,
  `docs/plan/planning/done/wellenlos/slice-188-archiv.zip`
- `befund`: `unzip -p …/slice-184-archiv.zip … | diff - <(git show ce5fb50^:…/slice-184-uses-tag-kohaerenz.md)`
  und dieselbe Probe für slice-188 zeigen je einen Unterschied: der im Archiv gespeicherte
  Volltext trägt bereits den nachgezogenen Link `wellenlos/slice-183-baseline-v5150.md`
  statt des ursprünglich committeten `../done/slice-183-baseline-v5150.md`. Ursache:
  `ApplySlice()` baut die Zip-Datei aus dem **aktuellen Arbeitsverzeichnis-Stand**
  (`buildZip` liest `slicePath` direkt von der Platte, vor dem `os.ReadFile` für den Stub).
  Da die sieben Läufe sequenziell in derselben Sitzung liefen (§2 Punkt 3 des Plans) und
  `RewriteRepo()` nach jedem Einzel-Lauf **repo-weit** greift, hat der slice-183-Lauf den
  noch unarchivierten slice-184/188-Volltext bereits umgeschrieben, bevor deren eigener
  Archivierungs-Lauf sie zippte. Die fünf übrigen Archive (141, 168, 169, 170) sind
  byte-identisch (per `diff` bestätigt) — das Muster tritt ausschließlich dort auf, wo ein
  Slice im selben Batch auf einen anderen, bereits archivierten Slice desselben Batches
  verweist.
- `verifizierbar`: ja — die genannten `diff`-Kommandos reproduzieren den Befund
  deterministisch; `TestRunSlice_Apply` und die übrige Testsuite decken diesen
  Cross-Reference-Fall nicht ab (kein Test mit zwei sich gegenseitig referenzierenden
  Slices im selben Lauf).
- `klasse`: archiv-batch-cross-reference-link-drift

**Einordnung:** Kein Datenverlust und keine inhaltliche Verfälschung der Entscheidung —
der Link zeigt weiterhin korrekt auf slice-183, nur eben auf dessen *neuen* statt
*historischen* Pfad. Trotzdem eine Abweichung von der in Modul 6 versprochenen
Unveränderlichkeit des Archivs: der archivierte Text entspricht nicht mehr dem, was zu
diesem Slice tatsächlich im Repository committet war. Da dies der erste beobachtete Fall
ist (Steering-Loop: 1×), reicht eine Beobachtungs-Registrierung; bei einer dritten
Wiederholung wäre `tools/archive-wave` entsprechend zu härten (z. B. Zip-Reihenfolge
so wählen, dass innerhalb eines Batches zuerst referenzierte, dann referenzierende Slices
gezippt werden, oder Zip-Aufbau vor jedem `RewriteRepo`-Lauf des Batches vorziehen).

### F-3 — Untracked Arbeitskopie von slice-200 bereits unter `done/` vorhanden

- `kategorie`: LOW
- `quelle`: AGENTS.md §3.3 / Modul 5 §Lifecycle als State Machine (Zustand = Verzeichnis, `git mv`)
- `pfad`: `docs/plan/planning/done/slice-200-verbliebene-wellenlose-slices-archivieren.md` (untracked)
- `befund`: Im Arbeitsverzeichnis liegt zum Review-Zeitpunkt zusätzlich zur getrackten
  Datei in `in-progress/` eine **untracked**, inhaltsgleiche Kopie unter `done/` (`git
  status --short` zeigt nur diese eine Zeile, `diff` gegen die `in-progress/`-Fassung ist
  leer). Sie wurde offenbar kopiert statt per `git mv` verschoben und nie gestaged — ein
  `git add -A` an dieser Stelle würde beide Pfade parallel committen, statt den
  Lifecycle-Übergang als reinen Move abzubilden.
- `verifizierbar`: ja — `git status --short` und `diff` wie oben.
- `klasse`: stray-untracked-lifecycle-copy

**Einordnung:** Betrifft nicht den geprüften Commit `ce5fb50` selbst, sondern den
aktuellen Arbeitsbaum zum Review-Zeitpunkt. Vor dem eigentlichen Closure-Move (`git mv
in-progress/... done/...`) sollte diese Datei gelöscht werden, sonst entsteht daraus ein
zweiter, nicht per Rename erkannter Pfad.

### F-4 — Fünf `exempt-paths`-Einträge in `.d-check.closure.yml` zeigen jetzt auf nicht mehr existierende Pfade

- `kategorie`: INFO
- `quelle`: `spec/lastenheft.md` DC-FA-RVW-001; `internal/hexagon/core/rules/reviews.go`
- `pfad`: `.d-check.closure.yml:244-248` (`reviews.exempt-paths`, fünf Zeilen für
  slice-141/168/169/170/188)
- `befund`: Die fünf gelisteten Pfade (`docs/plan/planning/done/slice-141-…md` usw.)
  existieren nach der Archivierung nicht mehr an dieser Stelle. Funktional harmlos:
  `reviewCandidates()` (`internal/hexagon/core/rules/reviews.go:129`) scannt `done-dir`
  nicht rekursiv und findet die verschobenen Dateien ohnehin nicht mehr als Kandidaten —
  `exempt-paths` wird nur gegen tatsächlich gefundene Kandidaten gematcht
  (`matchAnyGlob(cfg.ExemptPaths, rel)`), ein Treffer auf einen nicht existierenden Pfad
  bleibt folgenlos. Die fünf Einträge sind damit tote, aber ungefährliche Konfiguration.
- `verifizierbar`: ja — `make review-coverage` bleibt grün; Lesen von `reviews.go:129-152`
  bestätigt die Nicht-Rekursion und die reine Filter-Semantik von `exempt-paths`.
- `klasse`: orphaned-exempt-path-entry

**Einordnung:** Kein Fehler, reine Aufräum-Gelegenheit — die AGENTS.md-Zeile zu
`make review-coverage` beschreibt diese Fünferliste bereits ausdrücklich als
„Bestands-Ausnahme, feste Dateiliste" für den einmaligen Scharfschalt-Moment; ihr
Weiterbestehen nach der Archivierung ist dieselbe Klasse wie ein „Migrations-Fangnetz,
das niemanden mehr fängt" (vgl. slice-168, `BEO-013`). Kein Blocker.

---

## Negativbefunde

- geprüft, ohne Befund: Stub-Form aller sieben `done/wellenlos/*.md` — Titel erhalten,
  kein Volltext, `unzip -p done/wellenlos/slice-<NNN>-archiv.zip`-Zeiger korrekt, `Welle:`-Feld
  unverändert übernommen, `Archiviert:`-Feld im etablierten Platzhalter-Format (identisch zu
  älteren Einzel-Slice-Stubs wie slice-083).
- geprüft, ohne Befund: `Hervorgegangen:`-Felder — die genannten ADR-/DC-Kennungen sind in
  den jeweiligen archivierten Volltexten tatsächlich enthalten (stichprobenartig gegen
  slice-168/169/170/183/184/188 geprüft); alle tragen denselben, bereits etablierten
  `d-check:ignore`-Kommentar ohne Slice-Provenienz-Verstoß (§3.7).
- geprüft, ohne Befund: fünf der sieben Zip-Archive (slice-141, 168, 169, 170) —
  `unzip -l` liest sie an, Inhalt byte-identisch mit `git show ce5fb50^:<pfad>` (slice-183
  ist die siebte, ebenfalls identische — sie referenziert keinen anderen Slice des Batches).
- geprüft, ohne Befund: die drei nachgezogenen Referenz-Fixes (`AGENTS.md` → slice-170,
  `welle-87-results.md` → slice-188, `MR-057-baseline-v5150.md` → slice-188) — alle drei
  neuen Pfade existieren und lösen korrekt auf.
- geprüft, ohne Befund: Commit-Form — `git diff-tree --name-status -M -r ce5fb50` zeigt
  ausschließlich `D`/`A`-Paare, keine erkannten Renames; deckt sich mit MR-062s
  beschriebener Form für den Einzel-Slice-Modus.
- geprüft, ohne Befund: keine neuen Kommentare mit Slice-Nummer-/Review-Befund-Provenienz
  im Diff (§3.7) — der einzige neu hinzukommende Kommentartyp ist der bereits etablierte,
  wiederholte `d-check:ignore`-Boilerplate der `Hervorgegangen:`-Zeile.
- geprüft, ohne Befund: `ApplySlice()`/`ReadWelleField`/`ExtractTitle` lehnen keinen der
  sieben Slices fälschlich ab (kein Slice trägt ein echtes `Welle:`-Feld); alle sieben
  waren laut Plan als `— wellenlos` deklariert und sind es auch im archivierten Volltext.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** citation-stretched-beyond-scope (Bündelungsfrage MR-062) ·
archiv-batch-cross-reference-link-drift · stray-untracked-lifecycle-copy ·
orphaned-exempt-path-entry

## Verdikt

**Merge-blockierend:** nein, mit zwei Auflagen für die Closure-Notiz. Die eigentliche
Werkzeug-Anwendung ist korrekt: fünf von sieben Archiven sind exakt byte-identisch, alle
sieben Stubs sind formkonform, alle drei Referenz-Nachzüge lösen auf, und die Commit-Form
erfüllt MR-062. Beide MEDIUM-Findings sind aber substanziell genug, um in die
Closure-Notiz (§9) bzw. ins Beobachtungs-Register aufgenommen zu werden, statt kommentarlos
zu verschwinden:

- F-1 (Bündelungsfrage) ist eine echte, von den eigenen Konventionsdokumenten
  (MR-062/MR-063) benannte offene Frage, die dieser Slice erneut unbeantwortet lässt.
- F-2 (Archiv-Drift bei Cross-Referenzen) ist ein erster gemessener Fund einer Klasse, die
  `tools/archive-wave`s bisherige Testsuite nicht abdeckt; sie betrifft ausschließlich die
  Reihenfolge-Empfindlichkeit von Batch-Läufen mit gegenseitig referenzierenden Slices.

F-3 ist vor dem tatsächlichen `git mv`-Closure-Schritt zu bereinigen (Arbeitsbaum-Hygiene,
kein Commit-Defekt). F-4 ist reine Aufräum-Gelegenheit ohne Dringlichkeit.

**Übergabe:** Findings gehen an den Implementer/Planner dieses Slice. Die Finding-Klassen
gehen in die Slice-Closure §7 und von dort in den Zähler des Beobachtungs-Registers.
Dieser Report ist Lauf-Beleg und ersetzt keine Verifikation — DoD-/Spec-Konformität
(insbesondere `make gates`/`make fullbuild` auf dem Endstand) prüft die unabhängige
Verifikation separat.
