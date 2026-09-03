# MR-059 — Wellen-Archiv-Stub-Move ist ein einziger, deklarierter Commit (Nachtrag zu MR-013)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. `AGENTS.md` §3.3 (Baseline-§3.3-Adaption
  über [MR-013](MR-013-lifecycle-move-buendelung.md)) verlangt Move und
  Inhaltsänderung als zwei Commits, **weil** ein reiner `git mv` die
  Rename-Detection (`R100`) über der 50 %-Schwelle hält. Diese Regel benennt
  den Fall, in dem diese Prämisse strukturell nicht zutrifft.
- **Datum:** 2026-09-03
- **Geltungsbereich:** `tools/archive-wave` (`Apply()`), jeder
  Wellen-Archivierungs-Commit (`make archive-wave WELLE=<id> APPLY=1`).
- **Adaption:** Ein Slice- oder Welle-Stub **ersetzt** den Volltext durch
  einen komplett neuen, templatierten Text ([`archiv-stub-slice.template.md`](../../.harness/baseline/v5.18.0/templates/docs/plan/planning/archiv-stub-slice.template.md)) —
  Identität, Archiv-Zeiger, Zustand, überlebende Kennungen, sonst nichts.
  Die alte Datei wird zugleich gelöscht, die neue liegt eine Verzeichnisebene
  tiefer (`done/<welle-id>/`). Anders als bei den drei MR-013-Fällen gibt es
  hier **keine Phase**, in der die bewegte Datei ihren Inhalt unverändert
  behält — Move und Inhaltsersetzung sind ein und derselbe Akt, weil der
  Zweck der Operation genau diese Ersetzung ist. Eine Zwei-Commit-Zerlegung
  (erst `git mv`, dann Stub-Text) wäre kein Feintuning, sondern
  Selbstwiderspruch: der Move-Commit enthielte eine Kopie des vollen
  Alt-Textes am neuen Pfad, die der zweite Commit sofort wieder verwirft.
  **Ein Wellen-Archivierungs-Commit bleibt deshalb bewusst ein einziger
  Commit** — Zip-Bau, Stub-Schreiben, Alt-Löschen, Verweis-Nachzug — und
  git zeigt dafür reine `D`/`A`-Paare, keine erkannten Renames (gemessen:
  `git diff-tree --name-status -M` über mehrere Wellen-Archivierungs-Commits
  dieses Repos). Die Commit-Botschaft deklariert das ausdrücklich, dieselbe
  Form wie beim MR-/Wellen-Lifecycle-Move in MR-013, aus demselben Grund —
  „git mv trifft hier nicht zu" statt eines stillschweigenden Bruchs der
  Zwei-Commit-Erwartung.
- **Begründung:** Gemessen bei [slice-191](../../docs/plan/planning/done/slice-191-alt-bestand-archivieren.md)
  (Anwendung von [slice-190](../../docs/plan/planning/done/slice-190-wellen-archiv-werkzeug.md)s
  Werkzeug auf welle-01…85, 85 Archivierungs-Commits): jeder Commit trägt
  D/A-Paare statt Renames, keiner löste `make planning-check` oder eine
  andere Lifecycle-Kopplung aus (die betroffenen Dateien sind bereits in
  `done/`, keine Lifecycle-Grenze wird überschritten). Ohne diese Regel
  bliebe unklar, ob ein solcher Commit gegen `AGENTS.md` §3.3 verstößt oder
  eine weitere, unbenannte Ausnahme davon ist.
- **Grenze:** Diese Regel deckt nur die **Wellen-Archivierung** selbst
  (Modul 6 Schritt 4). Sie ist keine Blankovollmacht für beliebige
  Content-Move-Commits — jeder andere Fall prüft weiterhin gegen die
  Zwei-Commit-Grundregel und die drei MR-013-Ausnahmen.
- **Auflösungs-Trigger:** permanent, solange `tools/archive-wave` Stubs
  durch Volltext-Ersetzung erzeugt.
