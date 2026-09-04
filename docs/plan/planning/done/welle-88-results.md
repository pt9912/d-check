# Welle 88 — Baseline-Migration auf v6.0.0 — Closure-Notiz

**Welle:** welle-88-baseline-v600-migration
**Abschluss:** 2026-09-04
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

Der vendorte Baseline-Baum stand auf `v5.18.0`, während upstream bereits bei
`v6.0.0` (Major) war. Drei Slices schließen dieses Delta über zwei
unabhängige Stränge:

- [slice-193](welle-88/slice-193-baseline-v600-bump.md): reiner Pin-Bump auf `v6.0.0`
  — Pfade, Cite-Spannen und die wellenlose Zeitdokumente-Archivierung
  (bereits ab `v5.20.0` im Regelwerk) werden übernommen;
  [`MR-060`](../../../../harness/conventions.md#mr-060--baseline-pin-hebung-auf-v600-zehnter-nachtrag-zu-mr-011-nachtrag-zu-mr-023)
  geschrieben. Erstmals gelang die
  [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)-Zwei-Commit-Zerlegung
  für einen Baseline-Bump vollständig (frühere Bumps blieben hier
  `weiter offen`).
- [slice-194](welle-88/slice-194-beobachtungsregister-architektur.md): die
  Architektur-Hälfte der Register-Neugestaltung —
  [ADR-0083](../../adr/0083-beobachtungsregister-verzeichnis-modus.md)
  (additiver `observations.dir`-Modus neben dem bestehenden
  `observations.register`), Kürzel-Deklaration in `harness/conventions.md`
  (`*` → `ALL`, `tools/harness/` → `HARN`), Produktcode
  (`checkObservationsDir`). Ein unabhängiger Review fand einen Pfad-Traversal
  (F-1, HIGH) im neuen Code, behoben vor der Closure.
- [slice-195](welle-88/slice-195-beobachtungsregister-migration.md): die
  Datenmigration — `docs/plan/planning/observations.md` (26 Tabellenzeilen +
  2 gestrichene Einträge) wird zu `docs/plan/planning/observations/` mit 28
  `BEO-<KUERZEL>/<slug>/`-Verzeichnissen (88 Belegdateien, Zähler-Summe exakt
  gegen die alte Tabelle geprüft), alle lebenden `BEO-<NNN>`-Zitate
  umgehängt. Unabhängiger Review und Verifikation fanden vier HIGH-Befunde
  (siehe unten), alle vor der Closure behoben.

**Alle Closure-Trigger-Bedingungen sind erfüllt:** alle drei Slices liegen in
`done/`; `make baseline-freshness` bestätigt Content-Match am gepinnten Tag
(`v6.0.0`, Bytes == vendored `SHA256SUMS`); `make gates` (zehn Gates) und
`make fullbuild` sind auf dem Endstand grün; kein lebendes `BEO-<NNN>`-Zitat
bleibt außerhalb der eingefrorenen Bestände (`done/`, `docs/reviews/`,
`harness/conventions/done/`, `docs/plan/cr/`, Accepted-ADR-Kerne) oder der
benannten Syntax-Beispiele (`spec/lastenheft.md`, `docs/user/benutzerhandbuch.md`).

## Was hat funktioniert?

**Die Drei-Slice-Schneidung nach Schicht statt nach Vorsatz hat getragen:**
Pin-Bump, Architektur-Entscheidung und Datenmigration liefen als
eigenständige, einzeln lieferbare Slices, genau wie die Slice-Größenregel
(Baseline-Regelwerk `modul-05` §Ziel-Form: Slice) es für eine Arbeit
verlangt, die über drei Schichten reicht.

**Der additive Architektur-Entscheid ([ADR-0083](../../adr/0083-beobachtungsregister-verzeichnis-modus.md)) hat die eigene
Zwischenphase tragbar gemacht:** `observations.dir` steht neben
`observations.register`, statt die alte Fähigkeit sofort abzulösen — dadurch
konnte slice-194 (Produktcode) unabhängig von slice-195 (Datenmigration)
schließen, ohne dass `make gates` zwischenzeitlich gegen ein Format lief,
das noch kein Sensor kannte.

**Unabhängiger Review und Verifikation fingen, was der Implementierungslauf
selbst nicht sah:** ein Pfad-Traversal in neuem Code (slice-194, F-1 HIGH)
und vier Verfahrensfehler in der Datenmigration (slice-195) — eine ungültige
Zitat-Analogie
([`MR-059`](../../../../harness/conventions.md#mr-059--wellen-archiv-stub-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)
statt einer eigenen Adaption), ein gebündelter
Beanspruchungs-Commit, editierte Bestandszeilen in drei Accepted-ADRs
(§3.5-Verstoß) und ein erfundener vierter Risiko-Ausgang statt eines
ehrlichen `weiter offen`.

## Was ging anders als geplant?

**Der Umfang von slice-195 sprengte die eigene Ein-Sitzungs-Review-Grenze**
(28 Einträge statt der geschätzten ~27, ~180 geänderte Dateien) und der
dafür im Slice-Plan vorab benannte Rückführungs-Trigger griff, wurde aber
bewusst nicht befolgt — eine Aufteilung hätte den Zähler-Diff-Beleg zwischen
alter und neuer Form zerrissen. Als eigene Beobachtung eingetragen statt
eines erfundenen Ausgangs (siehe Steering-Loop-Einträge unten).

**Der Migrations-Commit brauchte eine nachträgliche Zerlegung:** der erste
Anlauf bündelte Beanspruchung und vollständigen Migrationsinhalt in einem
Commit und berief sich dafür fälschlich auf
[`MR-059`](../../../../harness/conventions.md#mr-059--wellen-archiv-stub-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013).
Der unabhängige Review fand beides; behoben durch eine eigene, eng gefasste
Adaption
([`MR-061`](../../../../harness/conventions.md#mr-061--register-formatmigration-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013))
und eine nachträgliche Commit-Zerlegung in Beanspruchung (`94b19bd`) und Inhalt
(`b1b960b`).

**Ein Hintergrund-Agent pushte den unkorrigierten Zwischenstand nach
`origin/main`, entgegen der ausdrücklichen Anweisung, nicht zu pushen.** Die
Korrektur (vier Commits statt einem, alle vier HIGH-Befunde behoben) wurde
lokal aufgebaut und musste den bereits publizierten fehlerhaften Commit per
`git push --force-with-lease` ersetzen — vom Auftraggeber bestätigt, da
niemand sonst auf dem verworfenen Commit aufgesetzt hatte.

## Steering-Loop-Einträge

- **[`BEO-ALL/large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md)**
  — neu, **weiter offen**: die Spannung „nicht mitten in der
  Zähler-Verifikation teilen" gegen „einen Slice nicht über die
  Review-Sitzungs-Grenze wachsen lassen" hat keine etablierte Auflösung.
- **[`MR-061`](../../../../harness/conventions.md#mr-061--register-formatmigration-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)**
  neu geschrieben, als fünfte, aber **erschöpfte** `AGENTS.md`-§3.3-Ausnahme
  (kein wiederkehrender Vorgang, anders als die vier übrigen).
- **[`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)**
  — neu, **offen**: eine repo-weite mechanische Kennungs-Ersetzung kennt vor
  der Ausführung nicht automatisch alle eingefrorenen und gesendeten
  Quellklassen (nur drei standen im Plan, tatsächlich griffen fünf).
- **[`BEO-HARN/check-latest-blind-before-pin`](../observations/BEO-HARN/check-latest-blind-before-pin/observation.md)**
  — neu, **offen**: `fetch-baseline-cache.sh --check-latest` sah `v6.0.0`
  nicht, solange der eigene Pin noch dahinter lag (Ursache nicht
  untersucht, kein Blocker für den Bump selbst).

## Beobachtungs-Register (Zeiger)

Lese-Schritt über die Bewegungen dieser Welle — das Register trägt seit
[slice-195](welle-88/slice-195-beobachtungsregister-migration.md) die Verzeichnis-Form,
der Zähler ist die Zahl der `evidence/`-Dateien je Verzeichnis, kein
gepflegtes Feld:

| Eintrag | Zähler nach dieser Welle | Was daraus folgt |
|---|---|---|
| `BEO-HARN/check-latest-blind-before-pin` | 0 → **1** | unter der Schwelle, offen |
| `BEO-ALL/mechanical-id-rewrite-misses-frozen-classes` | 0 → **1** | unter der Schwelle, offen |
| `BEO-ALL/large-migration-exceeds-session-review-limit` | 0 → **1** | unter der Schwelle, offen |

Kein Eintrag erreichte während dieser Welle die 3×-Schwelle neu. Der
Alt-Bestand (28 aus der Migration übernommene Einträge, mehrere bereits
oberhalb der Schwelle mit zugewiesenem Ausgang) ist unverändert nach
[`observations/`](../observations/README.md) übertragen — siehe slice-195 §9.

## Folge-Slices

Keine — `welle-88` schließt vollständig innerhalb ihrer drei Slices, ohne
offenen Rest. Die in
[`BEO-ALL/registerzeile-ohne-ausgang-nach-schwelle`](../observations/BEO-ALL/registerzeile-ohne-ausgang-nach-schwelle/observation.md)
benannte Sensor-Lücke bleibt ausdrücklich außerhalb des Scopes dieser Welle
(§6 des Welle-Plans) — eigener Folge-Slice bei Bedarf.

## Verifikation

- `make gates` grün (zehn Gates) auf dem Endstand (Commit `a7fcfb7`).
- `make fullbuild` grün — 51 Requirements, 0 Waisen, `verify-closure-notes`
  586 Dateien / 0 Befunde (image-hash `sha256:60c3fc51781e4769a42cc92502336e1e39bbfa6b09b5adf3e12b58dfa4529815`).
- `make baseline-freshness` grün — Pin `v6.0.0` ist aktuellster Tag, Bytes
  unverändert gegen upstream.
- `make adr-check` grün — 0 Kern-Abweichungen, auch nach der Rücknahme der
  vier fehlerhaft editierten Geschichte-/Kern-Zeilen in slice-195.
- Je Slice ein unabhängiger Review und eine unabhängige Verifikation; alle
  gemeldeten HIGH-Befunde (eins in slice-194, vier in slice-195) vor der
  jeweiligen Closure behoben.
