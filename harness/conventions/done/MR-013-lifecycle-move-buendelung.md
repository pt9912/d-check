# MR-013 — Lifecycle-Move-Commit bündelt gekoppelte Verweise

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. `AGENTS.template.md` §3.3 deckt seit
  `v6.6.0` zwei Fälle (Regelfall · Lifecycle-Übergang nach `done/`, Inhalt
  vor reinem `git mv`) und damit **alle** Lifecycle-Move-Bündelungen dieses
  Repos vollständig ab — Slice-Lifecycle (Beanspruchung, Closure) **und**
  MR-/Wellen-Lifecycle-Moves. Die ursprüngliche Aussage dieses Eintrags
  nannte `modul-05-planning-harness.md` als ersetzte Stelle — das war
  durchgehend die **falsche** Stelle, wie
  [slice-223](../../../docs/plan/planning/in-progress/slice-223-commit-zerlegung-ausnahmen-aufloesen.md)
  gezeigt hat; die tatsächlich einschlägige Kanon-Stelle war
  `grundlagen-traceability.md` §Herkunfts-Anker: *„Der `git mv` zieht die
  Pfad-Berichtigung nach sich, als eigener Commit nach dem Umzug."* — der
  **Regelfall** (Fall 1) selbst, nicht eine Ausnahme davon.
- **Datum:** 2026-06-21 (getrimmt 2026-09-16, vollständig aufgelöst
  2026-09-16 nach Review-Korrektur —
  [slice-223](../../../docs/plan/planning/in-progress/slice-223-commit-zerlegung-ausnahmen-aufloesen.md))
- **Geltungsbereich:** keiner mehr. War: `AGENTS.md` §3.3, Slice- und
  MR-/Wellen-Lifecycle-Moves.
- **Adaption (Verlauf, damit die Kette nachvollziehbar bleibt):**

  **2026-06-21 — Einträg entsteht.** Seit `make planning-check` (slice-040)
  den Roadmap-Zustand atomar an den `in-progress`-Stand koppelt, wäre ein
  *byte*-reiner Move-Commit beim Slice-Lifecycle-Übergang zwangsläufig
  gate-rot. Adaption: Der `git mv`-Commit trägt **zusätzlich** den
  Roadmap-Flip und alle Pfad-Verweise auf den Slice; der Slice-Body selbst
  bleibt Commit 2. Später erweitert auf die Beanspruchung (`open/` →
  `in-progress/`, seit slice-185) und auf MR-/Wellen-Lifecycle-Moves (seit
  welle-79) — drei Fälle unter einem Eintrag.

  **2026-09-16, erster Schritt — [slice-223](../../../docs/plan/planning/in-progress/slice-223-commit-zerlegung-ausnahmen-aufloesen.md)
  trimmt auf einen Fall.** Die Slice-Lifecycle-Hälfte (zwei der drei Fälle)
  erwies sich als reine Git-Semantik ohne Kanon-Bedarf: Die bewegte
  Slice-Datei selbst bleibt beim Move unverändert, nur fremde Dateien
  (Roadmap, Sensoren-Index) reisen mit — das berührt die Rename-Erkennung
  nicht und braucht keine Abweichung. Für den dritten Fall (MR-/Wellen-Move,
  wo die bewegte Datei ihre **eigenen** relativen Verweise trägt) blieb der
  Eintrag zunächst bestehen, mit der Begründung, der lokale
  `pre-commit`-Hook mache den kanonischen Zwei-Commit-Weg (Move, dann
  Korrektur, beide im selben Push) lokal uncommittierbar.

  **2026-09-16, zweiter Schritt — der unabhängige Review widerlegt die
  verbleibende Begründung, derselbe Slice löst vollständig auf.** Der
  Review von slice-223 (Finding F-1, HIGH) stellte den eigenen Move-Commit
  der fünf `MR`-Dateien (`a148466d`) isoliert nach (`git worktree add … &&
  make doc-check`) und fand ihn **rot** — neun `target-missing`-Befunde,
  weil die externen Rückverweise aus `harness/conventions.md` und
  `harness/sensors/archive-wave.md` erst im **nächsten** Commit korrigiert
  wurden. Trotzdem hatte der lokale `pre-commit`-Hook diesen Commit
  **durchgelassen**. Der Grund: Der Hook prüft `make doc-check` gegen den
  **Arbeitsbaum** (`docker run -v "$PWD:/repo:ro"`), nicht gegen den
  git-Diff des jeweiligen Commits. Zum Zeitpunkt des Commits lagen die
  externen Korrekturen bereits unstaged im Arbeitsbaum — der Hook sah einen
  grünen Baum, unabhängig davon, was tatsächlich in diesem einen Commit
  landete. **Ein reiner Move-Commit ist damit lokal sehr wohl
  committierbar**, solange die Korrektur vor dem `git commit`-Aufruf schon
  im Arbeitsbaum geschrieben, aber gezielt aus der Staging-Area
  ausgeschlossen ist (`git add <nur-die-Move-Dateien>`) — exakt die
  Technik, mit der dieses Repo in [slice-224](../../../docs/plan/planning/done/slice-224-baseline-v690-bump.md)
  bereits einen DoD-Haken vor einem reinen Move committet hatte, nur hier
  nicht bewusst als Beleg gegen die eigene `MR-013`-Prämisse erkannt.
  Die Prämisse *„lokal nicht committierbar"* war damit falsch, nicht nur
  unbelegt — der einzige verbleibende Fall braucht keine Abweichung mehr.

- **Begründung:** Sichtbar 2026-06-21 — die PR-/Push-CI prüft den Push-Tip,
  der ein Zwischen-Commit sein kann; sie lief auf dem reinen Move-Commit von
  slice-040 rot. Für **alle drei** ursprünglichen Fälle war das im Kern
  dasselbe Missverständnis: Ein Konflikt zwischen Rename-Erkennung und
  Gate-Grün wurde angenommen, wo entweder (a) die bewegte Datei ohnehin
  unverändert bleibt (Slice-Lifecycle) oder (b) der Kanon selbst schon die
  richtige Form vorgibt und der lokale Hook sie nicht verhindert
  (MR-/Wellen-Move) — geprüft erst am 2026-09-16, fünf Wellen nach der
  Einführung des Eintrags.
- **Auflösungs-Trigger:** eingetreten — vollständig durch den Kanon gedeckt,
  kein verbleibender Geltungsbereich.
