# MR-063 — Eigenständiger Review-Archiv-Move ist ein einziger, deklarierter Commit (Nachtrag zu MR-013)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. `AGENTS.md` §3.3 (Baseline-§3.3-Adaption
  über [MR-013](MR-013-lifecycle-move-buendelung.md)) verlangt Move und
  Inhaltsänderung als zwei Commits, **weil** ein reiner `git mv` die
  Rename-Detection (`R100`) über der 50 %-Schwelle hält. Diese Regel benennt
  einen weiteren Fall, in dem diese Prämisse strukturell nicht zutrifft.
- **Datum:** 2026-09-04
- **Geltungsbereich:** `tools/archive-wave`s Review-Modus (`ApplyReview()`,
  jeder `-review=<datei> -apply`-Lauf) — **nicht** der Wellen-Modus (trägt
  bereits [MR-059](MR-059-wellen-archiv-stub-move.md)) und nicht der
  Einzel-Slice-Modus (trägt bereits
  [MR-062](MR-062-wellenloser-slice-archiv-move.md)) oder sonstige
  Slice-Lifecycle-Moves (die drei MR-013-Fälle).
- **Adaption:** Genau wie bei Welle und Einzel-Slice ersetzt ein
  Review-Archiv-Move Inhalt und Ablageort im selben Akt: die Review-Datei
  wandert nach `docs/reviews/archiv/`, ihr Volltext wird durch einen
  komplett neuen, templatierten Stub ersetzt (`ReviewStub()`,
  `tools/archive-wave/stub.go`). Es gibt **keine Phase**, in der die bewegte
  Datei ihren Inhalt unverändert behält — die Inhaltsersetzung ist der Zweck
  der Operation, nicht ihr Nebeneffekt. Der Archiv-Commit bleibt deshalb
  bewusst ein einziger, in der Botschaft ausdrücklich deklarierter Commit;
  `git diff-tree --name-status -M` zeigt dafür reine `D`/`A`-Paare, keine
  erkannten Renames.
- **Begründung:** Dieselbe Prämisse wie bei MR-059/MR-062: dort existiert
  keine Phase mit unverändertem Inhalt, deshalb keine Zwei-Commit-Zerlegung
  möglich. Der Unterschied zu beiden ist nur der Modus (eigenständiger
  Review statt Welle oder Slice) — MR-059 beschränkt seinen Geltungsbereich
  ausdrücklich auf `tools/archive-wave` (`Apply()`, Wellen-Modus) und
  schließt eine Anwendung darüber hinaus als Zitat-Analogie explizit aus;
  MR-062 tut dasselbe für seinen eigenen Modus. Diese Regel liefert die für
  den Review-Modus fehlende, eigenständige Grundlage — unabhängiger
  Code-Review von slice-199 (2026-09-04) fand denselben Lücken-Typ, der bei
  slice-197 bereits einen eigenen Korrektur-Commit brauchte, hier aber
  direkt beim Schreiben dieser Regel geschlossen wird, statt einen zweiten
  Korrektur-Commit zu erzwingen.
- **Grenze:** Diese Regel deckt **ausschließlich** den Review-Archiv-Move
  selbst (`ApplyReview()`). Sie ist keine Blankovollmacht für beliebige
  Content-Move-Commits — jeder andere Fall prüft weiterhin gegen die
  Zwei-Commit-Grundregel und die bereits benannten Ausnahmen (MR-013s drei
  Fälle, MR-059, MR-062). **Bündelung mehrerer unabhängiger Review-Moves in
  einem gemeinsamen Commit ist zulässig**, wenn die Operation mechanisch
  uniform und werkzeug-verifiziert ist — ohne inhaltliches
  Unterscheidungsurteil je Datei. Das ist dieselbe Praxis, die slice-197
  bereits für den Einzel-Slice-Modus anwendete (45 Einzel-Slice-Archive in
  einem Commit gebündelt), obwohl MR-062 diese Frage für seinen eigenen
  Geltungsbereich formal offen ließ; diese Regel beantwortet sie hier nur
  für den Review-Modus, nicht rückwirkend für MR-062.
- **Auflösungs-Trigger:** permanent, solange `tools/archive-wave`s
  Review-Modus Stubs durch Volltext-Ersetzung erzeugt.
