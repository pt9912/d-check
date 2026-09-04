# Welle 89 — Wellenlose Review-Archivierung nachgerüstet — Closure-Notiz

**Welle:** welle-89-wellenlose-review-archivierung
**Abschluss:** 2026-09-04
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

`docs/reviews/` trug ~65 Review-Reports ohne archivierte Welle — Reports zu
wellenlos geschlossenen Slices, für die `tools/archive-wave` bisher keinen
Träger hatte (das Werkzeug verlangte zwingend eine `-welle`-Kennung). Der
Baseline-Bump auf `v6.0.0` (welle-88) hatte die dafür gemeldete Kanon-Lücke
bereits geschlossen (`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht:
die Slice-Closure selbst archiviert, Schlüssel ist der Slice). Zwei Slices,
nach dem welle-87-Muster (Werkzeug bauen, dann anwenden):

- [slice-196](welle-89/slice-196-archive-wave-slice-modus.md): `tools/archive-wave`
  bekam einen zweiten Modus (`-slice=<id>`), mutually exclusive zu `-welle`.
  An einem konstruierten Fixture bewiesen, nicht am echten Bestand.
- [slice-197](welle-89/slice-197-wellenlosen-bestand-archivieren.md): der neue Modus
  auf alle 45 wellenlosen `done/`-Slices mit noch unarchivierten Reviews
  angewendet — mehr als die geschätzten ~43.

**Alle Closure-Trigger-Bedingungen sind erfüllt:** beide Slices liegen in
`done/`; `docs/reviews/` enthält keinen Report mehr zu einem archivierbaren
wellenlosen Slice (ein permanenter, benannter Rest per `ignore-refs`-Tombstone
gedeckt); `make gates` (zehn Gates) und `make fullbuild` sind auf dem
Endstand grün; `make archive-wave-test` ist grün.

## Was hat funktioniert?

**Die Zwei-Slice-Schneidung (Werkzeug bauen, dann anwenden) hat wieder
getragen** — dieselbe Disziplin wie bei welle-87: slice-196 blieb klein und
an einem Fixture beweisbar, die eigentliche Belastungsprobe kam erst in
slice-197, am gewachsenen, realen Bestand.

**Drei unabhängige Prüf-Runden (Review + Verifikation je Slice, plus die
eigene erneute Anwendung in slice-197) fingen zusammen drei echte
Werkzeug-Fehler, die am kleinen Test-Fixture unsichtbar geblieben waren:**
eine Pfad-Kollision mit den Closure-Struktur-Prüfungen, ein abgeschnittenes
mehrzeiliges Feld, und eine Zugehörigkeits-Prüfung, die auf Kontext-
Erwähnungen fremder Wellen ansprang. Keiner der drei fehlerhaften
Zwischenstände wurde dauerhaft committet.

## Was ging anders als geplant?

**slice-196s ursprüngliches Design (Stub bleibt am flachen, unveränderten
Pfad) erwies sich erst am echten Bestand als unzureichend** — kollidierte
mit `structure`/`planning.closure` (kein DoD/keine Closure-Notiz im Stub)
und erzeugte eine `ids`-Welle (bare Kennungen). Auf Nutzer-Rücksprache
korrigiert: ein gemeinsames Verzeichnis `docs/plan/planning/done/wellenlos/`
für alle wellenlosen Archive statt eines Unterverzeichnisses je Slice —
entgeht denselben nicht-rekursiven Scan-Mustern wie `done/<welle-id>/` es
für den Wellen-Modus tut, ohne 45 fast leere Ordner zu erzeugen.

**Zwei weitere Werkzeug-Fehler traten erst bei der Anwendung auf den realen
Bestand auf** (mehrzeiliges `**Welle:**`-Feld wurde abgeschnitten; die
Zugehörigkeits-Prüfung griff danach auf Kontext-Erwähnungen fremder Wellen)
— beide führten zu einer vollständigen Neuanwendung aller 45 Slices, jeweils
vor dem endgültigen Commit gefunden und behoben.

**Die im Plan zitierte [`MR-059`](../../../../harness/conventions.md#mr-059--wellen-archiv-stub-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)-Analogie
für die Commit-Granularität des neuen Einzel-Slice-Modus trug nicht** —
jene Regel grenzt ihren Geltungsbereich ausdrücklich auf den Wellen-Modus
ein, derselbe Zitat-Überdehnungs-Fehler wie bei slice-195/welle-88. Behoben durch eine
eigene, permanente Adaption ([`MR-062`](../../../../harness/conventions.md#mr-062--wellenloser-einzel-slice-archiv-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)).

## Steering-Loop-Einträge

- **[`BEO-ALL/large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md)**
  — zweites Auftreten (slice-195, jetzt slice-197), Zähler 1 → **2**, weiter
  offen: dieselbe Spannung wie bei welle-88 benannt, weiterhin ohne
  etablierte Auflösung.
- **[`BEO-ALL/review-collection-misses-non-slice-filenames`](../observations/BEO-ALL/review-collection-misses-non-slice-filenames/observation.md)**
  — neu, offen: `tools/archive-wave`s Sammel-Logik (beide Modi) findet
  Review-Reports nur über ein `slice-<NNN>`-Dateinamens-Muster; ein Report
  ohne dieses Muster bleibt strukturell unsichtbar.
- **[`MR-062`](../../../../harness/conventions.md#mr-062--wellenloser-einzel-slice-archiv-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)**
  neu geschrieben, als sechste `AGENTS.md`-§3.3-Ausnahme (permanent,
  wiederkehrend wie [`MR-059`](../../../../harness/conventions.md#mr-059--wellen-archiv-stub-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)
  — anders als die erschöpfte [`MR-061`](../../../../harness/conventions.md#mr-061--register-formatmigration-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)).

## Beobachtungs-Register (Zeiger)

Lese-Schritt über die Bewegungen dieser Welle: kein Eintrag erreichte
während dieser Welle die 3×-Schwelle neu (beide oben genannten Einträge
stehen bei Zähler 2 bzw. 1, unter der Schwelle).

| Eintrag | Zähler nach dieser Welle | Was daraus folgt |
|---|---|---|
| `BEO-ALL/large-migration-exceeds-session-review-limit` | 1 → **2** | unter der Schwelle, weiter offen |
| `BEO-ALL/review-collection-misses-non-slice-filenames` | 0 → **1** | unter der Schwelle, offen |

## Folge-Slices

Keine — `welle-89` schließt vollständig innerhalb ihrer zwei Slices, ohne
offenen Rest.

## Verifikation

- `make gates` grün (zehn Gates) auf dem Endstand.
- `make fullbuild` grün — 532 Dateien, 0 Befunde.
- `make archive-wave-test` grün, inklusive der Regressionstests für alle
  drei am Bestand gefundenen Werkzeug-Fehler.
- `make baseline-freshness` grün.
- Je Slice ein unabhängiger Review und eine unabhängige Verifikation; alle
  gemeldeten Befunde (1 HIGH in slice-196, 1 HIGH + 2 Verifikations-Befunde
  in slice-197) vor der jeweiligen Closure behoben.
