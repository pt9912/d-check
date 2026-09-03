# Review-Report: slice-194 — 2026-09-03

**Review-Art:** Code — geprüft gegen Plan, ADR und Konventionen (Modul 10
§Drei Review-Arten).

**Gegenstand:** Commit `e186b6f` (feat(planning): Beobachtungs-Register-
Verzeichnis-Modus — fünfte Fähigkeit, additiv, slice-194, welle-88,
ADR-0083, DC-FA-PLAN-001)

**Skill:** `.harness/skills/reviewer.md` @ v1.13.0
**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-03

**Eingangs-Kontext:**

- [slice-194](../plan/planning/in-progress/slice-194-beobachtungsregister-architektur.md)
- [ADR-0083](../plan/adr/0083-beobachtungsregister-verzeichnis-modus.md)
- [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
- `AGENTS.md` §3.1 (Filesystem-Port als einzige Dateisystem-Tür), §3.4
  (Spec-Straten-Abwärtssperre), §3.7 (Kommentarklassen)
- `internal/hexagon/port/driven/filesystem.go` §Port-Vertrag ("alle Pfade
  sind '/'-getrennt und relativ zur Repo-Wurzel")

---

## Findings

### F-1 — `declaredObservationIDsFromDir.Has` verlässt sich auf den
Default-Pattern, um Pfad-Traversal zu verhindern; ein konfiguriertes
`observations.pattern` kann das unterlaufen

- `kategorie`: HIGH
- `quelle`: Filesystem-Port-Vertrag (`internal/hexagon/port/driven/filesystem.go:27`) /
  ADR-0083
- `pfad`: `internal/hexagon/core/rules/planning_observations.go:118-121`
  (`declaredObservationIDsFromDir.Has`)
- `befund`: `Has(id)` bildet `path.Join(d.root, id, "observation.md")` und
  reicht das ungeprüft an `fsys.Kind` weiter, dessen Adapter
  (`internal/adapter/driven/fs/fs.go:23-24`) den Pfad per `filepath.Join`
  gegen die Repo-Wurzel auflöst und dabei `..`-Segmente kollabiert. `id`
  stammt aus einem Regex-Treffer gegen beliebigen zitierenden
  Markdown-Fließtext (`citedWithoutRow`). Mit dem eingebauten Default
  (`BEO-[A-Z][A-Z0-9]*/[a-z][a-z0-9-]*`, ohne `.`) ist `id` harmlos — es
  kann kein `..`-Segment entstehen. Der Default ist aber nur der
  **Default**: `observations.pattern` überschreibt ihn, und
  `configyaml.applyObservations` prüft NUR `Register`/`Dir`/`Dirs` auf
  führendes `/` bzw. enthaltenes `..` (`configyaml.go:2283-2288`) — das
  Pattern selbst wird nur auf `regexp.Compile`-Gültigkeit geprüft
  (`configyaml.go:2289-2294`), nicht auf die Zeichen, die es matchen darf.
  Ein Repo, das `observations.pattern` versehentlich oder bewusst
  permissiver setzt (z. B. um zusammengesetzte Slugs zuzulassen), und eine
  zitierende Markdown-Datei mit einem Treffer wie
  `BEO-X/../../../../etc/passwd`, lässt `Has()` einen Lstat außerhalb von
  `Dir` und potenziell außerhalb der Repo-Wurzel ausführen — ein
  Existenz-Orakel für beliebige Dateisystempfade, im Widerspruch zum
  dokumentierten Versprechen des Filesystem-Ports ("alle Pfade … relativ
  zur Repo-Wurzel"). Der Tabellen-Modus hat dieses Risiko strukturell
  nicht: `declared` ist dort eine reine In-Memory-Map ohne
  Dateisystemzugriff pro Kennung.
- `verifizierbar`: ja — Unit-Test mit `Dir: "obs"`,
  `Pattern: "BEO-.+"` (o. ä. permissiv) und einer zitierten Kennung mit
  `../`-Segmenten gegen einen `MemFS`/echten `fs.Adapter`; erwartet wäre
  eine Ablehnung am Config-Rand oder eine Pfad-Normalisierung/-Prüfung in
  `Has()`, nicht ein Lstat außerhalb der Wurzel.
- `klasse`: Abgeleiteter Dateisystempfad aus regex-erfasstem Fließtext ohne
  eigene Traversal-Prüfung — Sicherheits-Invariante gilt nur für
  konfigurierte Pfade, nicht für aus Zitaten abgeleitete.

## Negativbefunde

- geprüft, ohne Befund: mit dem **Default**-Pattern ist `id` nachweislich
  frei von `.` und einem zweiten `/` — Traversal ist über den
  ausgelieferten Default nicht erreichbar; das Risiko ist auf eine
  bewusste `observations.pattern`-Überschreibung beschränkt.
- geprüft, ohne Befund: `Register`/`Dir` mutually exclusive, Exit 2 bei
  Kollision (`configyaml.go:2278-2281`), Test
  `TestDecode_ObservationsFehler["register und dir zugleich (ADR-0083)"]`
  grün.
- geprüft, ohne Befund: `Dir`, `Dirs`, `Register` selbst sind gegen
  führendes `/` und `..` validiert — dieselbe Prüfung wie vor diesem
  Slice, jetzt auch auf `Dir` angewendet.
- geprüft, ohne Befund: Tabellen-Modus unverändert — alle
  Alt-Tests (`TestObservationsCitedWithoutRow`,
  `TestObservationsCodeSpanIsExampleButLinkTextCounts`,
  `TestObservationsOnlyFirstCellDeclares`,
  `TestObservationsFailClosedAndInert`,
  `TestObservationsTwoIDsOnOneLineDoNotCollide`,
  `TestObservationsTwoBadDirsDoNotCollide`,
  `TestObservationsFirstCellMayCarryBackticks`) laufen unverändert grün
  gegen den umgebauten `switch`-Dispatch in `CheckPlanningObservations`.
- geprüft, ohne Befund: Symlink-Konsistenz — `Kind()` löst per `Lstat`
  nicht auf; ein `observation.md`, das ein Symlink ist, gilt korrekt als
  NICHT nachgewiesen (`KindSymlink != KindFile`), dieselbe Disziplin wie
  `DC-FA-LINK-002` andernorts.
- geprüft, ohne Befund: die kanonische Zitierform
  `` [`BEO-ALL/foo`](../../observations/BEO-ALL/foo/observation.md) ``
  matcht die Kennung sowohl im Linktext als auch (unbeabsichtigt) innerhalb
  des Link-Ziels; beide Fundstellen erzeugen bei fehlendem Nachweis
  potenziell zwei `Finding`-Werte mit identischem
  (Datei, Zeile, Regel, Ziel, Grund)-Tupel — `model.SortFindings` dedupt
  sie vor der Ausgabe (`run.go:115`) zu einem sichtbaren Befund. Kein
  Doppel-Report, verifiziert durch Nachvollzug der Dedup-Logik
  (`finding.go:121-149`).
- geprüft, ohne Befund: `spec/lastenheft.md`s neuer Fließtext (fünfte
  Fähigkeit, 0.85.0-Historie-Zeile) enthält **keine** `ADR-`/`slice-`/
  `welle-`-Erwähnung — AGENTS.md §3.4 (spec-straten referenzieren nie
  abwärts) eingehalten; `make doc-check` (Modul `matrix`) läuft mit 0
  Befunden gegen den Endstand.
- geprüft, ohne Befund: `harness/conventions.md`s neue Spalte
  `BEO-Kürzel` ist eine wohlgeformte 5-Spalten-Tabelle (Header, Trenner,
  zwei Datenzeilen je 7 `awk -F'|'`-Felder); `ALL`/`HARN` sind
  kollisionsfrei und beide konform zum neuen
  `DefaultObservationDirPattern`-Kürzel-Segment (`[A-Z][A-Z0-9]*`).
- geprüft, ohne Befund: ADR-0083 folgt der vollen Vorlage (Kontext,
  Entscheidung, Verglichene Alternativen, Konsequenzen, Fitness Function,
  Re-Evaluierungs-Trigger); die drei verglichenen Alternativen
  (vollständige Ablösung, eigenständiges Modul, Auto-Migrations-Skript)
  sind eigenständig begründet und keine Strohmänner.
- geprüft, ohne Befund: `make gates` (zehn Gates: baseline-verify,
  workflow-pins, doc-check, lint, test, arch-check, coverage-gate,
  semgrep, gate-consistency, planning-check), `make test`,
  `make lint`, `make coverage-gate` (94,6 % ≥ 93 %) und
  `docker run --rm --network none … d-check:latest` (500 Dateien, 0
  Befunde) laufen selbst ausgeführt grün — die Commit-Botschafts-Zusage
  ist gedeckt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 0 |

### F-2 (MEDIUM) — Der implizite `Dirs`-Default im Verzeichnis-Modus ist ungetestet

- `kategorie`: MEDIUM
- `quelle`: Maintainability / Prüffrage 13 (Negativtest zu neuem Vertrag)
- `pfad`: `internal/hexagon/core/rules/planning_observations.go:80-82`
  (`scanCitedWithoutDeclaration`, Zweig `len(dirs) == 0`)
- `befund`: Ist `observations.dirs` leer, verwendet der Verzeichnis-Modus
  `path.Dir(o.Dir)` als Scan-Wurzel — anders als im Tabellen-Modus, wo
  `path.Dir(o.Register)` auf einer DATEI operiert, operiert dieser Aufruf
  auf einem VERZEICHNIS und liefert dessen ELTERN-Verzeichnis, nicht das
  Verzeichnis selbst. Das trifft die im Slice-Plan genannte Ziel-Form
  (`observations.dir = docs/plan/planning/observations` ⇒ Default-Scan
  `docs/plan/planning`) zufällig richtig, ist aber in keinem Test
  belegt: `obsDirCfg` in `planning_observations_test.go` übergibt in
  jedem Aufruf explizit mindestens ein `dirs`-Element, der
  `len(dirs)==0`-Zweig läuft für den neuen Modus nie durch einen Test.
- `verifizierbar`: ja — `go test` mit `model.ObservationsConfig{Dir: "…"}`
  ohne `Dirs` und einer Coverage-Messung auf
  `scanCitedWithoutDeclaration:81` (aktuell laut `make coverage-gate`
  über den Tabellen-Modus-Aufruf mitabgedeckt, nicht über den
  Verzeichnis-Modus-Aufruf separat).
- `klasse`: Verzweigung mit geänderter Argument-Semantik (Datei- vs.
  Verzeichnis-Pfad an `path.Dir`), die für den neuen Zweig keinen eigenen
  Test hat.

### F-3 (LOW) — ADR-0083s Fitness-Function-Text überzeichnet den Kollisions-Test

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `docs/plan/adr/0083-beobachtungsregister-verzeichnis-modus.md:108`
- `befund`: Der Satz „der Kollisionsfall (`Dir` und `Register` beide
  gesetzt) bricht mit dem erwarteten Grund-Code" suggeriert eine Prüfung
  auf einen strukturierten `Reason`/Grund-Code wie bei einem
  `model.Finding`. Tatsächlich ist ein Config-Ladefehler ein einfacher
  `error` ohne Grund-Code-Feld, und der zugehörige Test
  (`TestDecode_ObservationsFehler`) prüft nur generisch `err != nil`,
  nicht den Fehlertext oder einen spezifischen Code.
- `verifizierbar`: ja — Lesen von `applyObservations` (gibt `error`
  zurück, kein `model.Finding`) und des Testkörpers (gemeinsame
  `err == nil`-Prüfung über alle Fälle der Map).
- `klasse`: Terminologie einer Kategorie (Grund-Code) auf einen Kontext
  übertragen, der sie nicht führt — Präzisions-, keine Korrektheitslücke.

**Finding-Klassen dieses Laufs:** Abgeleiteter Dateisystempfad aus
regex-erfasstem Fließtext ohne Traversal-Prüfung; ungetesteter
Default-Zweig bei geänderter Argument-Semantik; Terminologie-Überzeichnung
im ADR-Beleg.

## Verdikt

**Merge-blockierend: ja (F-1, HIGH).** Die additive Fähigkeit selbst,
ihre Rückwärtskompatibilität, die ADR-Struktur, die Spec-Straten-Disziplin
und alle zehn Gates sind sauber — aber `declaredObservationIDsFromDir.Has`
bricht das Filesystem-Port-Versprechen ("Pfade relativ zur Repo-Wurzel")
für einen aus Zitat-Text abgeleiteten Pfad, sobald `observations.pattern`
vom sicheren Default abweicht. Das ist vor der Closure zu schließen, nicht
danach — der Tabellen-Modus, den diese Fähigkeit erweitert, hatte diese
Angriffsfläche nie. F-2 und F-3 sind vor der Closure klärbar, aber für
sich genommen nicht blockierend.

**Übergabe:** Alle drei Findings gehen an den Implementer für einen
Folge-Commit auf demselben Slice (Commit `e186b6f` bleibt gemäß Auftrag
unverändert). Dieser Report ist ein Lauf-Beleg und ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat, in
eigenem Kontext.
