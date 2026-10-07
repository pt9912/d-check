# MR-074 — Baseline-Pin-Hebung auf `v6.17.0` (sechzehnter Nachtrag zu MR-011, Nachtrag zu MR-023)

- **Datum:** 2026-10-07
- **Geltungsbereich:** §Baseline, pin-gebundene Verweise,
  `.harness/baseline/v6.17.0/`
- **Ersetzt-Baseline-Regel:** — *(Pin-Fortschreibung im Bundle-Layout des
  Vorgängers; keine inhaltliche Abweichung)*
- **Adaption:** Der vendorte Bestand steht auf
  [`v6.17.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v6.17.0).
  Vier Zwischen-Tags (`v6.14.0`, `v6.14.1`, `v6.15.0`, `v6.16.0`) liegen
  zwischen dem alten Pin und diesem — ein Bump zielt auf den jeweils aktuellen
  Tag, nicht auf jeden dazwischen.

  **Der Delta, gemessen und nicht geschätzt** (`diff -I '<!-- Quelle:'` gegen
  den Vorzustand, beide Bäume noch vorhanden): 55 Dateien vorher wie nachher,
  keine neu, keine entfallen. **22** tragen ein inhaltliches Delta (8
  Regelwerk-Dateien, 14 Templates), dazu `SHA256SUMS`. Größte Bewegungen:
  `regelwerk/grundlagen-harness-dateien.md` mit 65 Diff-Zeilen,
  `regelwerk/grundlagen-referenz-richtung.md`,
  `templates/spec/spezifikation.template.md` und
  `templates/harness/README.template.md` mit je 17, `modul-13-quality-gates.md`
  mit 15.

  **Inhaltliche Bewegungen — vollständig nach Datei, keine davon übernommen**
  (Abgrenzung, slice-253; Adoptieren ist Sache eigener Folge-Slices):

  1. **Werkzeug-eigener Teil des Gate-Index** (`harness/mk/<werkzeug>.md`,
     fünf Bedingungen, darunter die Disjunktheit) —
     `grundlagen-harness-dateien.md`, `modul-13-quality-gates.md`,
     `grundlagen-begriffe.md`; nachgezogen in `templates/AGENTS.template.md`
     §4, `templates/harness/README.template.md`, `templates/Makefile` und
     `templates/.d-check.yml`. **Dieses Repo ist betroffen:** das `Makefile`
     bindet das Werkzeug-Fragment `a-check.mk` ein (erzeugt aus
     `a-check --print-mk`), und `targets.makefiles` liest nur `Makefile` —
     dessen einziges Target `a-check` steht bisher weder im Index noch in
     einem Werkzeug-Teil (gemessen). Die Frage, ob d-check einen Werkzeug-Teil
     für `a-check.mk` führt, ist ein Kandidat für einen Folge-Slice.
     **Eingelöst durch slice-254:** kein Werkzeug-Teil (a-check schreibt keinen),
     das Target steht in der `arch-check`-Zeile von `harness/README.md`
     §Sensors, und `targets.makefiles` liest das Fragment mit.
  2. **Festlegungen der Harness-Werkzeuge gehören in die Spezifikation**, nicht
     in ein viertes Stratum — `grundlagen-referenz-richtung.md`,
     `modul-03-spec.md`, dazu die Spec-Anteile von
     `grundlagen-harness-dateien.md` und `grundlagen-begriffe.md`;
     nachgezogen in `templates/spec/spezifikation.template.md` (neuer §7
     *Festlegungen der Harness-Werkzeuge*, die Historie rückt von §7 nach
     **§8**), `templates/harness/sensors/gate.template.md` (Sensor-Datei
     verlinkt die Spec-Kennung statt Schwelle und Randform zu führen), die
     Spec-Anteile von `templates/harness/README.template.md`,
     `templates/docs/plan/adr/NNNN-titel.template.md` und
     `templates/docs/plan/adr/README.template.md` (die ADR eines solchen Gates
     schärft dessen Spec-Stelle). Die Umnummerierung koppelt an die Ausnahme
     `"7. Historie"` ([`MR-0098`](../conventions.md#mr-0098)) — eine spätere
     Adoption fiele dort laut auf.
  3. **Reviewer** — `modul-10-review-harness.md`: kein Stil-Polizist
     (Formatierung oder Benennung ohne Konventions-Anker ist kein Finding),
     kein HIGH- oder MEDIUM-Finding ohne Failure-Szenario; die Fundstelle wird
     als wörtliches, eindeutig auffindbares Kurzzitat geankert, die Zeile ist
     Lesehilfe. Nachgezogen in `templates/.harness/skills/reviewer.template.md`
     (dort **zusätzlich**: LOW nur mit Konventions-Anker),
     `templates/docs/reviews/review-report.template.md` und
     `templates/.harness/skills/closure-note-reviewer.template.md` (dort nur
     die Kurzzitat-Regel).
  4. **Register-Kennung** `BEO-<NNN>` → `BEO-<KUERZEL>/<slug>` bzw. „die
     Beobachtung" — `modul-10-review-harness.md` sowie
     `templates/docs/plan/planning/welle-results.template.md` (nur dieser
     Nachzug), `templates/docs/plan/planning/slice.template.md`,
     `templates/harness/conventions.template.md` und
     `templates/docs/reviews/review-report.template.md`. Dieses Repo führt die
     Verzeichnis-Form bereits.
  5. **Audit-Span-Pflichtfelder** — `modul-15-observability.md`: liefert die
     Quelle den Wert eines Pflichtfelds nicht, bleibt es Pflicht und wird
     ausdrücklich als nicht bekannt gekennzeichnet (nicht `0`, nicht `false`),
     mit Nennung der Quelle. Für dieses Repo ohne Gegenstand (keine
     Audit-Spans).
  6. **Bundle-Stand** — `regelwerk/README.md` (Kurs-Welle, Datum, Pin-URLs).

  **Spiegel-Klassen, gemessen am echten Vorzustand** (HEAD vor dem Swap): 72
  Dateien außerhalb des vendorten Baums nannten `v6.13.0` — 180 Vorkommen,
  ohne den Plan dieses Slice (der den Tag als Gegenstand nennt; mit ihm 73 und
  184) —, **dazu zwei Skills unter `.harness/skills/`** mit 8 Vorkommen, die die
  erste Messung nicht sah, weil sie `.harness/` ganz ausnahm; erst das rote
  `doc-check` nach dem Entfernen des alten Baums zeigte sie. **Für die nächste
  Hebung:** die Spiegel-Messung schließt nur `.harness/baseline/` aus, nicht
  `.harness/`. **Lebend
  retargetet:** `AGENTS.md`, `harness/README.md`, die drei `harness/rules/`-
  Dateien, alle aktiven `MR-*` außer dem Vorgänger, `.claude/agents/reviewer.md`,
  die beiden Skills, `roadmap.md`, `planning/README.md`, `spec/architecture.md`
  und `spec/spezifikation.md` (Rolle-Zeile). **Zeilenweise getrennt:**
  `harness/conventions.md` (§Baseline und §Adoptierte Konventions-Quellen
  retargetet, die Index-Zeile des Vorgängers wandert nach §Aufgelöste
  Adaptionen) und `observations/README.md` (Zeilen 5 und 18 retargetet,
  Zeile 21 bleibt — sie nennt, aus welcher Baseline-Form die Trigger-Audit-
  Adoption stammt). **Frozen, byte-stabil:** die `done/`-Slices und die
  Welle-91-Datei, die Review-Reports, [ADR-0100](../../docs/plan/adr/0100-targets-authority-liste.md),
  der eingehende CR vom 2026-10-06, die CHANGELOG-Historie, der Kommentar zur
  Vorgänger-Stufe in `.d-check.yml` und der Vorgänger-Eintrag selbst (Zug nach
  `conventions/done/`). **Anderer Gegenstand:** `tools/harness/selbstpruefung.sh`
  nennt die Version des Schwester-Werkzeugs, aus dem es adoptiert ist. Die 8
  Baseline-Symlinks unter `.claude/rules/` lösen gegen `v6.17.0`.

  **`ignore-refs` wächst diesmal nicht** — gemessen: nach dem Entfernen des
  alten Baums meldet `doc-check` 0 Befunde. Keine eingefrorene Datei trägt den
  `v6.13.0`-Baum als Markdown-Link; die `d-check:cite`-Direktiven in `done/`
  liegen in Verzeichnissen, die `citations.scope` ausnimmt, die übrigen
  Nennungen stehen in Inline-Code oder Prosa.

  **`d-check:cite`-Neu-Ankern** — gemessen **14** lebende Direktiven, **4**
  davon neu geankert, je reine Zeilenverschiebung bei wortgleichem Text
  (`MR-031` `grundlagen-harness-dateien.md` 194→196; `MR-035`
  `grundlagen-begriffe.md` 49→50; Reviewer-Skill `modul-10-review-harness.md`
  82→87; Closure-Note-Skill
  `templates/.harness/skills/closure-note-reviewer.template.md` 83→84) — und
  **10** tag-only.

  **Zitat-Delta** ([`MR-039`](../conventions.md#mr-039)): die wörtlichen
  Zitate ohne Direktive in Markdown (`AGENTS.md` „Halluzinierte Gates …" aus
  `modul-13`, `MR-056` zu `modul-05` §Lifecycle) stehen im `v6.17.0`-Wortlaut
  unverändert. **Ein Zitat hat ein Delta:** der Kommentar zum
  `targets`-Block in `.d-check.yml` zitiert `templates/AGENTS.template.md` §4
  mit Auslassung („Der Gate-Index steht einmal ... Diese Datei fuehrt die
  Liste nicht."); genau die ausgelassene Mitte trägt in `v6.17.0` den Satz,
  dass Targets aus Werkzeug-Fragmenten im Teil des Werkzeugs stehen. Der
  zitierende Kommentar bleibt nach MR-039 unverändert; seine Aussage hängt an
  Bewegung 1 oben.
- **Begründung:** Der Nachtlauf meldete neuere Tags (`make
  baseline-freshness`); das Release `v6.17.0` lag vor. Die Hebung ist eine
  reine Fortschreibung; die inhaltlichen Bewegungen (oben) warten auf eigene
  Folge-Slices.
- **Auflösungs-Trigger:** die nächste Pin-Hebung.
- **Löst auf:** [`MR-073`](../conventions.md#mr-073)
- **Ausgelöst durch Baseline-Stand:** `v6.17.0`
