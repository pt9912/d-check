# MR-062 — Wellenloser Einzel-Slice-Archiv-Move ist ein einziger, deklarierter Commit (Nachtrag zu MR-013)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. `AGENTS.md` §3.3 (Baseline-§3.3-Adaption
  über [MR-013](MR-013-lifecycle-move-buendelung.md)) verlangt Move und
  Inhaltsänderung als zwei Commits, **weil** ein reiner `git mv` die
  Rename-Detection (`R100`) über der 50 %-Schwelle hält. Diese Regel benennt
  einen weiteren Fall, in dem diese Prämisse strukturell nicht zutrifft.
- **Datum:** 2026-09-04
- **Geltungsbereich:** `tools/archive-wave`s Einzel-Slice-Modus
  (`ApplySlice()`, jeder `-slice=<id> -apply`-Lauf) — **nicht** der
  Wellen-Modus (der trägt bereits [MR-059](MR-059-wellen-archiv-stub-move.md))
  und nicht die Beanspruchung oder sonstige Slice-Lifecycle-Moves (die drei
  MR-013-Fälle).
- **Adaption:** Genau wie beim Wellen-Modus ersetzt ein wellenloser
  Einzel-Slice-Archiv-Move Inhalt und Ablageort im selben Akt: die Slice-Datei
  wandert nach `docs/plan/planning/done/wellenlos/`, ihr Volltext wird durch
  einen komplett neuen, templatierten Stub ersetzt
  ([`archiv-stub-slice.template.md`](../../.harness/baseline/v6.0.0/templates/docs/plan/planning/archiv-stub-slice.template.md)),
  und ihre Review-Reports verschwinden ersatzlos ins Archiv. Es gibt **keine
  Phase**, in der die bewegte Datei ihren Inhalt unverändert behält — die
  Inhaltsersetzung ist der Zweck der Operation, nicht ihr Nebeneffekt. Der
  Archiv-Commit bleibt deshalb bewusst ein einziger, in der Botschaft
  ausdrücklich deklarierter Commit; `git diff-tree --name-status -M` zeigt
  dafür reine `D`/`A`-Paare, keine erkannten Renames.
- **Begründung:** Der Unterschied zu den drei MR-013-Fällen ist derselbe wie
  bei MR-059: dort existiert eine Phase mit unverändertem Inhalt (deshalb
  zwei Commits verlangt und möglich); hier nicht. Der Unterschied zu MR-059
  selbst ist nur der Modus (Einzel-Slice statt Welle) — dieselbe Prämisse,
  dieselbe Konsequenz, eine eigene Regel, weil MR-059 seinen Geltungsbereich
  ausdrücklich auf `tools/archive-wave` (`Apply()`) beschränkt und eine
  Anwendung darüber hinaus als Zitat-Analogie explizit ausschließt — genau
  dieser Fehler unterlief slice-197s eigenem Plan (§2 Punkt 3), gefunden von
  der unabhängigen Verifikation. Diese Regel liefert die fehlende,
  eigenständige Grundlage nachträglich.
- **Grenze:** Diese Regel deckt **ausschließlich** den Einzel-Slice-Archiv-Move
  selbst (`ApplySlice()`). Sie ist keine Blankovollmacht für beliebige
  Content-Move-Commits — jeder andere Fall prüft weiterhin gegen die
  Zwei-Commit-Grundregel und die vier bereits benannten Ausnahmen
  (MR-013s drei Fälle, MR-059). Ob **mehrere** unabhängige Einzel-Slice-Archive
  in einem gemeinsamen Commit gebündelt werden, ist eine andere Frage — diese
  Regel begründet nur, dass **je Slice** genau ein Commit die Untergrenze ist,
  nicht ob mehrere solcher Commits zusammengefasst werden dürfen.
- **Auflösungs-Trigger:** permanent, solange `tools/archive-wave`s
  Einzel-Slice-Modus Stubs durch Volltext-Ersetzung erzeugt.
