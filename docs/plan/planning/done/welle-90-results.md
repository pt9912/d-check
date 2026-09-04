# Welle 90 — Eigenständige Review-Archivierung — Closure-Notiz

**Welle:** welle-90-eigenstaendige-review-archivierung
**Abschluss:** 2026-09-04
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

`docs/reviews/` trug nach welle-89 noch 11 Review-Reports ohne
`slice-<NNN>` im Dateinamen — CR-, ADR-, Baseline-/MR-, Backlog-, Wellen-
und Release-Prep-Reviews. `tools/archive-wave`s beide bisherige Modi finden
Review-Reports ausschließlich über ein `slice-<NNN>`-Dateinamens-Muster;
diese 11 waren dafür strukturell unsichtbar
([`BEO-ALL/review-collection-misses-non-slice-filenames`](../observations/BEO-ALL/review-collection-misses-non-slice-filenames/observation.md)).
Zwei Slices, nach dem welle-87/89-Muster (Werkzeug bauen, dann anwenden):

- [slice-198](welle-90/slice-198-archive-wave-review-modus.md): `tools/archive-wave`
  bekam einen dritten Modus (`-review=<dateiname>`, mutually exclusive zu
  `-welle`/`-slice`). Ein eigenständiger Review bekommt — anders als ein
  Review zu einem Slice — einen eigenen Stub, weil er selbst der
  abgeschlossene Vorgang ist. An einem konstruierten Fixture bewiesen.
- [slice-199](welle-90/slice-199-eigenstaendige-reviews-archivieren.md): der
  neue Modus auf alle 11 eigenständigen Reviews angewendet, die nach
  welle-89 noch in `docs/reviews/` lagen.

**Alle Closure-Trigger-Bedingungen sind erfüllt:** beide Slices liegen in
`done/`; `docs/reviews/` enthält außerhalb `docs/reviews/archiv/` keinen
Report mehr ohne `slice-<NNN>` im Dateinamen; `make gates` (zehn Gates) und
`make fullbuild` sind auf dem Endstand grün; `make archive-wave-test` ist
grün.

## Was hat funktioniert?

**Die Zwei-Slice-Schneidung (Werkzeug bauen, dann anwenden) trug ein
drittes Mal** — slice-198 blieb an einem Fixture klein und beweisbar,
slice-199 wandte ihn unverändert auf den realen Bestand an. Anders als bei
welle-89 (drei Werkzeug-Fehler erst am realen Bestand gefunden) hielt der
Voll-Dry-Run vor der Anwendung diesmal ohne Überraschung — der neu gebaute
Modus war beim Anwenden bereits ausgereift.

**Die Titel-Extraktion war der einzige tatsächlich benannte Risikopunkt und
trug korrekt:** eine eigene `ExtractFullHeading`-Funktion statt der
Slice-/Welle-spezifischen `ExtractTitle` verhinderte das antizipierte Risiko
(führendes Wort einer uneinheitlichen Review-Überschrift verschluckt) —
gegen drei reale Überschriftenformen verifiziert.

## Was ging anders als geplant?

**Dieselbe Lücken-Klasse wie bei welle-89 trat ein zweites Mal auf, diesmal
beim Schreiben statt erst bei der Anwendung:** der neue `-review=<datei>`-
Modus hatte keine der fünf in `AGENTS.md` §3.3 enumerierten
Ein-Commit-Ausnahmen — unabhängiger Review und unabhängige Verifikation
fanden das übereinstimmend (F-1, HIGH) und verschärften die eigene, zunächst
als „Beobachtungswert" eingestufte Einschätzung der Verifikation auf
closure-blockierend. Behoben durch
[`MR-063`](../../../../harness/conventions.md#mr-063--eigenständiger-review-archiv-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)
plus sechste `AGENTS.md`-§3.3-Ausnahme — dieselbe Korrektur-Form wie bei
[`MR-062`](../../../../harness/conventions.md#mr-062--wellenloser-einzel-slice-archiv-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)
(slice-197), nur diesmal **vor** statt nach der Slice-Closure geschlossen.

**Zwei weitere, kleinere Lücken fielen erst der unabhängigen Prüfung auf:**
vier durch die Archivierung verwaiste `ignore-refs`-Einträge in
`.d-check.yml` (F-2) und ein von keinem Modul gescannter Klartext-
Pfadverweis in der immutablen [ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md)
(F-3) — beide behoben, Letzteres über
einen Geschichte-Anhang statt einer Kern-Änderung (§3.5-konform).

## Steering-Loop-Einträge

- **[`BEO-ALL/review-collection-misses-non-slice-filenames`](../observations/BEO-ALL/review-collection-misses-non-slice-filenames/observation.md)**
  — auf **gestrichen** gesetzt, unabhängig vom 3×-Zähler (stand bei 1): der
  neue `-review`-Modus deckt genau die Lücke, die die Beobachtung benannte,
  auch wenn die Sammel-Logik der beiden anderen Modi selbst unverändert
  bleibt.
- **[`MR-063`](../../../../harness/conventions.md#mr-063--eigenständiger-review-archiv-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)**
  neu geschrieben, als siebte `AGENTS.md`-§3.3-Ausnahme (permanent,
  wiederkehrend wie [`MR-059`](../../../../harness/conventions.md#mr-059--wellen-archiv-stub-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)/[`MR-062`](../../../../harness/conventions.md#mr-062--wellenloser-einzel-slice-archiv-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)
  — anders als die erschöpfte [`MR-061`](../../../../harness/conventions.md#mr-061--register-formatmigration-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)),
  mit expliziter Bündelungs-Klärung für mehrere unabhängige Review-Moves in
  einem Commit.
- **Kein neuer Beobachtungs-Register-Eintrag** aus dieser Welle — beide
  Lücken (F-1 fehlende MR, F-2 verwaiste `ignore-refs`, F-3 übersehener
  Klartext-Verweis) wurden direkt im Slice behoben, nicht als offene
  Beobachtung eingetragen.

## Beobachtungs-Register (Zeiger)

Lese-Schritt über die Bewegungen dieser Welle: kein Eintrag erreichte
während dieser Welle die 3×-Schwelle neu.

| Eintrag | Zähler nach dieser Welle | Was daraus folgt |
|---|---|---|
| `BEO-ALL/review-collection-misses-non-slice-filenames` | 1 (unverändert) | **gestrichen** — Ursache behoben, unabhängig vom Zähler |
| `BEO-ALL/large-migration-exceeds-session-review-limit` | 2 (unverändert) | unter der Schwelle, weiter offen — Gegenprüfung in slice-199 bestätigt: Umfang (11 Läufe, ein Commit) sprengte die Ein-Sitzungs-Grenze nicht |

## Folge-Slices

Keine — `welle-90` schließt vollständig innerhalb ihrer zwei Slices, ohne
offenen Rest.

## Verifikation

- `make gates` grün (zehn Gates) auf dem Endstand.
- `make fullbuild` grün, 0 Trace-Waisen bei unveränderten 51 Anforderungen.
- `make archive-wave-test` grün, inklusive der neuen Regressionstests für
  den Review-Modus (`review_mode_test.go`).
- Je Slice ein unabhängiger Review und eine unabhängige Verifikation; alle
  gemeldeten Befunde (slice-198: 1 HIGH + 2 LOW + 1 INFO; slice-199: 1 HIGH
  + 2 MEDIUM + 1 LOW, in beiden Fällen von Review und Verifikation
  unabhängig übereinstimmend gefunden) vor der jeweiligen Closure behoben.
