# MR-082 — Neue Wellen-Namen stehen neben dem Bestand: `wave-<NN>`, `archive.zip`, `seit wave-<NN>`

- **Datum:** 2026-11-26
- **Geltungsbereich:** [`.d-check.yml`](../../.d-check.yml) (`planning.waves`, `matrix`, `ids`),
  [`.d-check.closure.yml`](../../.d-check.closure.yml), [`tools/archive-wave/`](../../tools/archive-wave/),
  [`.claude/commands/`](../../.claude/commands/), [`AGENTS.md`](../../AGENTS.md) §3.7,
  [`planning-check.md`](../sensors/planning-check.md) und [`archive-wave.md`](../sensors/archive-wave.md)
- **Ersetzt-Baseline-Regel:**
  [`grundlagen-harness-dateien.md` §Verzeichniskonvention, Absatz *Dateinamen sind englisch*](../../.harness/baseline/v7.0.0/regelwerk/grundlagen-harness-dateien.md#verzeichniskonvention)
  — die Regel gilt in diesem Repo für **Neues**; der Bestand behält seine Namen.
- **Adaption:** Die zwei Schreibweisen stehen nebeneinander, je Gegenstand eine
  Zeile; die letzte Spalte sagt, welcher Sensor welche Form liest.

  | Gegenstand | Neu | Bestand | Wer liest welche Form |
  | --- | --- | --- | --- |
  | Kennung und Datei der Welle | `wave-<NN>-<titel>.md` | `welle-<NN>-<titel>.md` | `planning` (`make planning-check`) **nur die neue**; `matrix` beide (Klassen `wave` und `welle`, je mit eigenen Verboten); `tools/archive-wave` beide |
  | Ergebnisnotiz | `wave-<NN>-results.md` | `welle-<NN>-results.md` | `planning` (`wave-results-missing`) **nur die neue**; das Closure-Profil (`make verify-closure-notes`) beide |
  | Archiv | Welle: `done/<wave-id>/archive.zip`; Slice und Review: `<kennung>-archive.zip` | `archiv.zip`, `<slice-id>-archiv.zip` | `tools/archive-wave` **schreibt** die Form nach dem Präfix der Kennung (Slice und Review immer die neue) und liest kein Archiv |
  | Herkunfts-Anker | `seit wave-<NN>` | `seit welle-<NN>` | kein Sensor; die Anker-Paarung ist ein Schritt der Closure (`/close-welle`) |

  **Die Senkung von `planning`** (der Bestand ist für `wave-drift` und
  `wave-results-missing` ungewächtert) steht in
  [`ADR-0117`](../../docs/plan/adr/0117-wellen-namen-ueber-die-konfiguration.md); dieser
  Eintrag erzählt sie nicht nach. Das Produkt und seine Defaults (`welle-*`)
  bleiben.
- **Begründung:** Der Auftraggeber hat am 2026-11-25 entschieden: Dateinamen
  englisch, Prosa deutsch, **ab jetzt für Neues, der Bestand bleibt**
  ([`wave-98`](../../docs/plan/planning/wave-98-baseline-v700.md) §1, Entscheid 1). Der
  Bestand ist Beleg, nicht Gestaltungsmasse; gemessen am Arbeitsbaum von
  [`slice-298`](../../docs/plan/planning/in-progress/slice-298-neue-namen-wave-archive-anker.md):
  `git ls-files | grep -c 'welle-[0-9]'` ergibt 345,
  `git ls-files | grep -c 'archiv\.zip$'` ergibt 188 und
  `grep -c '^| welle-' docs/plan/planning/in-progress/roadmap.md` ergibt 38.
  Die Umbenennung wäre ein eigener Beschluss mit Eingriff in Belege.
- **Auflösungs-Trigger:** der Bestand ist weg — `git ls-files | grep -c 'welle-[0-9]\|archiv\.zip'`
  ergibt 0 (heute 435) **und**
  `git grep -h -o 'seit welle-[0-9][0-9]*' -- . ':!.harness/baseline' ':!docs/plan/planning/done' ':!docs/reviews' | wc -l`
  ergibt 0 (heute 19). Nicht *permanent*: es gibt einen Zähler. Er feuert nur auf
  einen Beschluss des Auftraggebers; bis dahin gilt dieser Eintrag. Die Texte, die
  den Bestand als Form nennen (`.claude/commands/`, `AGENTS.md` §3.7,
  die beiden Sensor-Dateien), gehen mit der Auflösung.
