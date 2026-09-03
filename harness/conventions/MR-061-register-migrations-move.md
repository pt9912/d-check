# MR-061 — Register-Formatmigration ist ein einziger, deklarierter Commit (Nachtrag zu MR-013)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. `AGENTS.md` §3.3 (Baseline-§3.3-Adaption
  über [MR-013](MR-013-lifecycle-move-buendelung.md)) verlangt Move und
  Inhaltsänderung als zwei Commits, **weil** ein reiner `git mv` die
  Rename-Detection (`R100`) über der 50 %-Schwelle hält. Diese Regel benennt
  einen weiteren Fall, in dem diese Prämisse strukturell nicht zutrifft — und
  korrigiert zugleich eine fehlgeschlagene Analogie: der Migrations-Commit
  von [slice-195](../../docs/plan/planning/done/slice-195-beobachtungsregister-migration.md)
  berief sich in seiner Botschaft auf [MR-059](MR-059-wellen-archiv-stub-move.md),
  dessen eigener Geltungsbereich (`tools/archive-wave`) und dessen eigene
  Grenzklausel (*„keine Blankovollmacht für beliebige Content-Move-Commits"*)
  diese Anwendung ausdrücklich ausschließen. Der unabhängige
  Verifikations-Report zu slice-195 hat das als Formfehler benannt: nicht die
  Ein-Commit-Entscheidung selbst, sondern ihre Begründung. Diese Regel liefert
  die fehlende, eigenständige Grundlage nachträglich.
- **Datum:** 2026-09-03
- **Geltungsbereich:** ausschließlich der eine Commit, der das
  Beobachtungs-Register von der Tabellen-Form auf die Verzeichnis-Form
  ([ADR-0083](../../docs/plan/adr/0083-beobachtungsregister-verzeichnis-modus.md))
  umgestellt hat — kein wiederkehrender Vorgang. Anders als die drei
  MR-013-Fälle und anders als MR-059 (die je einen fortlaufenden
  Operations-Typ decken) deckt diese Regel einen **einmaligen, bereits
  abgeschlossenen** Formatwechsel.
- **Adaption:** Ein Registerformat-Wechsel dieser Art ersetzt Inhalt und
  Ablageform im selben Akt: die alte Tabellendatei wird gelöscht, während an
  ihrer Stelle eine komplett neue Verzeichnisstruktur mit templatierten
  Dateien entsteht ([`observation.template.md`](../../.harness/baseline/v6.0.0/templates/docs/plan/planning/observation.template.md)),
  und jedes lebende Zitat der alten Kennungsform wird im selben Zug auf die
  neue umgehängt. Es gibt **keine Phase**, in der die alte Tabellendatei
  bereits gelöscht ist, während die neue Verzeichnisform noch fehlt oder ein
  Zitat noch auf die alte Form zeigt, und `make gates` dabei grün bliebe — jede
  Zwischenteilung schüfe einen Commit, gegen den `doc-check` (fehlende
  Zitat-Ziele) oder `planning-check` (unbekannte Register-Form) sofort rot
  liefe. Der Migrations-Commit bleibt deshalb ein einziger, in der Botschaft
  ausdrücklich deklarierter Commit; `git diff-tree -r --name-status -M`
  zeigt dafür reine `D`/`A`/`M`-Paare, keine erkannten Renames (gemessen: 2
  D / 149 A / 34 M, 0 Renames).
- **Begründung:** Der Unterschied zu den drei MR-013-Fällen ist derselbe wie
  bei MR-059: dort existiert eine Phase, in der die bewegte Datei ihren
  Inhalt unverändert behält (deshalb zwei Commits möglich und verlangt);
  hier — wie bei der Wellen-Archivierung — ist die Inhaltsersetzung der
  ganze Zweck der Operation, nicht ihr Nebeneffekt. Der Unterschied zu
  MR-059 ist der Geltungsbereich: MR-059 deckt einen **wiederkehrenden**
  Werkzeug-Aufruf (`tools/archive-wave`, jede künftige Wellen-Archivierung);
  diese Regel deckt einen **einmaligen** Registerformat-Wechsel, der mit
  seinem Vollzug erschöpft ist — es gibt kein zweites Mal, das
  `observations.md` in ein Verzeichnis migriert.
- **Grenze:** Diese Regel deckt **ausschließlich** den einen benannten
  Commit. Sie ist keine Blankovollmacht für beliebige
  Registerformat-Wechsel oder andere einmalige Struktur-Umbauten — ein
  künftiger, andersartiger Formatwechsel prüft erneut gegen die
  Zwei-Commit-Grundregel und braucht, falls dieselbe Prämisse nicht zutrifft,
  eine eigene MR-Adaption.
- **Auflösungs-Trigger:** keiner — die Regel ist mit dem Vollzug des einen
  Commits, den sie deckt, bereits erschöpft; sie bleibt als Beleg stehen,
  nicht als wirksame Erlaubnis für einen künftigen Fall.
