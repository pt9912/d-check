## Durchsetzungsschicht
*Quelle: [grundlagen/durchsetzungsschicht.md](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/durchsetzungsschicht.md)*

Konventionen, Hard Rules und Sensors sind *aspirativ*, bis etwas sie an
die Agent-Schleife **bindet**. Diese Seite beschreibt die
**Durchsetzungsschicht** — die fail-closed Mechanik, die aus „die Doku
sagt X" ein „der Harness erzwingt X" macht. Erst das tool-neutrale
Prinzip, dann eine konkrete Realisierung (Claude-Code-Hooks).

### Die Lücke: aspirativ vs. bindend

Ein Guide, der „make/Docker-only" oder „Gates vor dem Handoff" sagt, ist
*inferential feedforward* — er **informiert**. Ein driftender oder
vergesslicher Agent kann ihn ignorieren, ohne dass etwas passiert. Die
Durchsetzungsschicht verschiebt dieselbe Regel in die **computational**-
Spalte der [2×2-Matrix](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/konzeptkarte.md#2x2-schnellanker): die falsche
Handlung wird *technisch erschwert* (feedforward) oder *deterministisch
erkannt* (feedback) — nicht bloß abgeraten. Es ist dieselbe Bewegung wie
beim `check-references`-Gate ([Traceability-Constraint](grundlagen-konventionen.md#traceability-constraint)):
eine Doku-Regel bekommt einen mechanischen Wächter.

### Drei Bindepunkte

Jeder bindet einen anderen Punkt der Agent-Schleife — und fällt in einen
Quadranten, den du schon kennst:

| Bindepunkt | Wann | 2×2-Quadrant | Wirkung |
|---|---|---|---|
| **Tool-Call-Gate** | vor jedem Tool-Call | computational **feedforward** | falsche Handlung technisch verhindern (Tool-Allowlist / Befehls-Guard) |
| **Handoff-Gate** | bevor der Agent „fertig" meldet | computational **feedback** | deterministisch prüfen, dass die Gates wirklich liefen |
| **Workflow-Skelett** | beim Start einer Aufgabe | inferential feedforward | den Ablauf vorgeben (Slice-Workflow als feste Schrittfolge) |

Zwei der drei *erzwingen* (computational); das **Workflow-Skelett ist der
schwächste Bindepunkt** — es gibt den Ablauf vor, erzwingt ihn aber nicht
(inferential), und bleibt das einzige der drei, das ein Agent noch
ignorieren kann. Die fail-closed Klammer aus der Einleitung sind die zwei
Gates; das Skelett ist das Gerüst, das sie absichern.

Realisierung in Claude Code: ein `PreToolUse`-Hook (Tool-Call-Gate), ein
`Stop`-Hook (Handoff-Gate) und ein Slash-Command (Workflow-Skelett),
verdrahtet in `.claude/settings.json`. Portierbar: andere Harnesses haben
äquivalente Punkte (Pre-/Post-Tool-Hooks, Pre-Commit-Hooks, Pflicht-CI-
Jobs). Der Bindepunkt ist das Konzept, der Hook nur eine Form.

### Vier Design-Eigenschaften

1. **fail-closed.** Fehlt das Prüfmittel (Interpreter nicht da, Input
   unlesbar, zu tiefe Verschachtelung), wird **blockiert**, nicht
   durchgewunken. Ein Gate, das im Zweifel passieren lässt, ist keiner.
2. **Nachweis über Inhalt, nicht Diff.** Ein Content-Hash des Arbeitsbaums
   belegt „die Gates liefen auf *genau diesem* Stand". Inhaltsbasiert (statt
   diff-/status-basiert) hält der Nachweis über Commits hinweg — und ein
   Commit *ohne* vorherigen Gate-Lauf bleibt trotzdem erkennbar. Das ist
   die Mechanik gegen die [Harness-Lüge](grundlagen-konventionen.md#kernbegriffe)
   „ich hab die Gates laufen lassen".
3. **Loop-Guard.** Ein Handoff-Gate muss erkennen, ob es sich in derselben
   Runde schon einmal blockiert hat — sonst Endlosschleife bei dauerhaft
   rotem Gate. Der Hook gibt sich beim zweiten Anlauf frei.
4. **bootstrap-aware.** Der Gate erzwingt nur die Gates, die *schon
   existieren*. Ein [Bootstrap-aware Gate](grundlagen-konventionen.md#kernbegriffe)
   wächst mit der Reife; ein harter Handoff-Gate ab Schritt 0 bekämpft die
   weiche Frühphase des [Harness-Bootstraps](grundlagen-konventionen.md#harness-bootstrap).
   Erst binden, wenn es etwas zu binden gibt.

### Grenzen — ehrlich benannt

- Ein Befehls-Guard, der nur **Befehlspositionen** prüft, ist ein
  *Stolperdraht, keine Sandbox*: Interpreter-Umwege (`python -c "…"`)
  bleiben möglich. Sein Wert ist, *versehentliche* Drift zu verhindern,
  nicht böswillige.
- Der Inhalts-Nachweis hat eine Lücke bei frischem Klon bzw. gelöschtem
  State mit cleanem Tree (kein Nachweis prüfbar) — dort ist **CI das Netz**.
- Diese Grenzen zu *benennen* ist Pflicht. Ein Gate, das so tut, als decke
  es mehr ab, als es tut, ist selbst eine [Harness-Lüge](grundlagen-konventionen.md#kernbegriffe)
  — dieselbe Klasse wie ein undeklariertes Gate.

### Die Schicht wird selbst gesteuert

Die Durchsetzungsschicht ist Code *im* Harness — also unterliegt sie
demselben [Steering-Loop](grundlagen-konventionen.md#verbindung-zum-steering-loop)
wie alles andere. Ein Befehls-Guard etwa reift in Wellen: zuerst nur die
Befehlsposition, dann Sub-Shell-Rekursion (`bash -c "…"`), dann
kombinierte Flags (`-lc`, `-ec`). Genau diese Härtung *am Wächter selbst*
ist der Steering-Loop, auf den Harness angewandt. Der ausführliche Worked
Example dazu gehört nach [Modul 13](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/04-qualitaet/modul-13-quality-gates.md)
(folgt).

### Referenz-Implementierung

Das vollständige Artefakt-Set einer Durchsetzungsschicht:

- `.claude/settings.json` — Hook-Verdrahtung (welcher Hook an welchem Punkt)
- `.claude/hooks/*.sh` — Tool-Call-Gate (Befehls-Guard) und Handoff-Gate
  (Stop-/Gate-Nachweis)
- `.claude/commands/*.md` — Workflow-Skelett als Slash-Command
- `tools/harness/working-tree-hash.sh` + `record-gates.sh` — gemeinsame, <!-- d-check:ignore (Referenz-Artefakt im Fallstudien-Repo) -->
  inhaltsbasierte Nachweis-Quelle für Gate-Lauf *und* Handoff-Gate (eine
  Wahrheit, keine Logik-Dopplung)

Die Gate- und Traceability-Mechanik, an die das andockt, läuft in den
Fallstudien-Repos ([`fallstudien.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/fallstudien.md)) bereits real; die
Skripte werden hier **referenziert, nicht inline ausgerollt** — die
mechanische Hälfte bleibt klein und lesbar, die Details holt man aus der
Referenz.
