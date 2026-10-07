# MR-074 — Baseline-Pin-Hebung auf `v6.17.0` (sechzehnter Nachtrag zu MR-011, Nachtrag zu MR-023)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** — *(Pin-Fortschreibung im Bundle-Layout des
  Vorgängers; keine inhaltliche Abweichung)*
- **Datum:** 2026-10-07
- **Geltungsbereich:** §Baseline, pin-gebundene Verweise,
  `.harness/baseline/v6.17.0/`
- **Adaption:** Der vendorte Bestand steht auf
  [`v6.17.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v6.17.0).
  Vier Zwischen-Tags (`v6.14.0`, `v6.14.1`, `v6.15.0`, `v6.16.0`) liegen
  zwischen dem alten Pin und diesem — ein Bump zielt auf den jeweils aktuellen
  Tag, nicht auf jeden dazwischen (Präzedenz:
  [`MR-073`](../conventions.md#mr-073) und seine Vorgänger).

  **Der Delta, gemessen und nicht geschätzt** (`diff -I '<!-- Quelle:'` gegen
  den Vorzustand, beide Bäume noch vorhanden): 55 Dateien vorher wie nachher,
  keine neu, keine entfallen. **22** tragen ein inhaltliches Delta (8
  Regelwerk-Dateien, 14 Templates), dazu `SHA256SUMS`. Größte Bewegungen:
  `regelwerk/grundlagen-harness-dateien.md` mit 65 Diff-Zeilen,
  `regelwerk/grundlagen-referenz-richtung.md`,
  `templates/spec/spezifikation.template.md` und
  `templates/harness/README.template.md` mit je 17, `modul-13-quality-gates.md`
  mit 15.

  **Vier inhaltliche Schlagzeilen, keine davon übernommen** (Abgrenzung,
  slice-253): Erstens führen `grundlagen-harness-dateien.md`,
  `modul-13-quality-gates.md` und `grundlagen-begriffe.md` einen
  **werkzeug-eigenen Teil des Gate-Index** (`harness/mk/<werkzeug>.md`) unter
  fünf Bedingungen, darunter die Disjunktheit — die d-check als Produkt bereits
  opt-in prüft; dieses Repo hat keine Werkzeug-Fragmente. Zweitens legen
  `grundlagen-referenz-richtung.md` und `modul-03-spec.md` fest, dass die
  **Festlegungen der Harness-Werkzeuge** (was ein Gate prüft, wie es an
  Randformen entscheidet) in die Spezifikation gehören, nicht in ein viertes
  Stratum. Drittens schärft `modul-10-review-harness.md` den Reviewer: **kein
  HIGH- oder MEDIUM-Finding ohne Failure-Szenario**, und die Fundstelle wird
  als wörtliches, eindeutig auffindbares Kurzzitat geankert, die Zeile ist nur
  Lesehilfe. Viertens ziehen die Templates (Harness-README, Sensor-Vorlage,
  Reviewer-Skill, Welle-Ergebnis, `.d-check.yml`)
  diese drei nach. Ob und wie d-check davon etwas adoptiert, ist Sache eigener
  Folge-Slices.

  **Spiegel-Klassen, gemessen am echten Vorzustand** (HEAD vor dem Swap): 72
  Dateien außerhalb des vendorten Baums nannten `v6.13.0` — 180 Vorkommen —,
  **dazu zwei Skills unter `.harness/skills/`** mit 8 Vorkommen, die die erste
  Messung nicht sah, weil sie `.harness/` ganz ausnahm; erst das rote
  `doc-check` nach dem Entfernen des alten Baums zeigte sie. **Lebend
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

  **`d-check:cite`-Neu-Ankern, kein Zitat-Delta** — gemessen **14** lebende
  Direktiven, **4** davon neu geankert, je reine Zeilenverschiebung bei
  wortgleichem Text (`MR-031` `grundlagen-harness-dateien.md` 194→196;
  `MR-035` `grundlagen-begriffe.md` 49→50; Reviewer-Skill
  `modul-10-review-harness.md` 82→87; Closure-Note-Skill
  `templates/.harness/skills/closure-note-reviewer.template.md` 83→84) — und
  **10** tag-only. Kein Zitat-Delta ([`MR-039`](../conventions.md#mr-039)):
  die wörtlichen Zitate ohne Direktive (`AGENTS.md` „Halluzinierte Gates …"
  aus `modul-13`, `MR-056` zu `modul-05` §Lifecycle) stehen im
  `v6.17.0`-Wortlaut unverändert.
- **Begründung:** Der Nachtlauf meldete neuere Tags (`make
  baseline-freshness`); das Release `v6.17.0` lag vor. Die Hebung ist eine
  reine Fortschreibung; die inhaltlichen Deltas (oben) warten auf eigene
  Folge-Slices.
- **Auflösungs-Trigger:** die nächste Pin-Hebung.
