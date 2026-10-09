# Verifikation — slice-269: `releasing.md` zieht nach docs/maintainer/

- **Rolle:** Verifier (Modul 11) — DoD und Plan gegen den Stand, nicht der Diff gegen Plan/ADR
- **Gegenstand:** `docs/plan/planning/in-progress/slice-269-releasing-nach-docs-maintainer.md` §1, §2, §3 samt Plan-Änderung
- **Range:** `2b088bb1..HEAD` (7 Commits, Spitze `7e72816e`)
- **Eingang:** Commits, Reports R1 und R2 unter `docs/reviews/`
- **Datum:** 2026-10-09

## Messungen (Kommando und Ergebnis)

Alle Läufe selbst gefahren; Docker/make, keine Host-Toolchain.

| # | Kommando | Ergebnis |
|---|---|---|
| M1 | `git show -M --summary 877af2d4` | `rename docs/{user => }/maintainer/releasing.md (100%)` — die bewegte Datei unverändert im Move-Commit |
| M2 | `git show --stat 877af2d4` · `git show --stat d0e64cdb` | Move-Commit: 15 weitere Dateien (eingehende Verweise, Tombstone, Workflows, Skript); Folge-Commit: nur `docs/maintainer/releasing.md`, 23 Zeilen geändert |
| M3 | `git show d0e64cdb --word-diff=porcelain` | 23 Änderungen, ausschließlich Link-Ziele; je altes und neues Ziel dieselbe Datei (`version.md`, `spec/lastenheft.md`, ADRs, `operations.md`, `benutzerhandbuch-standard.md`, …) |
| M4 | `git diff fceb796d:docs/user/releasing.md HEAD:docs/maintainer/releasing.md` | gegen den Stand **vor** slice-268 nur 5 Zeilen, alle Link-Ziele auf Geschwister unter `docs/user/` — `docs/user/` und `docs/maintainer/` liegen gleich tief; kein Inhalt geändert |
| M5 | Wegwerf-Klon auf `2b088bb1`, `git mv docs/user/maintainer/releasing.md docs/maintainer/releasing.md`, `docker run --rm --network none -v <klon>:/repo:ro d-check:latest`, Befunde per `awk -F'\t'` nach Grund und Datei gezählt | **39** Befunde, solange das leere Verzeichnis im Klon stehen bleibt; nach `rmdir docs/user/maintainer` **40**. Aufteilung: `docs/maintainer/releasing.md` 12 `target-missing` + 11 `repo-escape`; CHANGELOG 1, README 2, README.de 2, Benutzerhandbuch 2, `packaging/dockerhub/README.md` 1, ADR-0014 4, ADR-0067 1, welle-73/74/82-results je 1; ein `codepath-missing` in der Closure-Notiz von slice-268 (Zeile 106). Deckt sich mit Plan §3 |
| M6 | `git ls-files \| xargs grep -n -I -e 'docs/user/releasing\.md' -e 'docs/user/maintainer' -e 'user/maintainer'` ohne `done/` und `docs/reviews/` | kein Link mehr auf einen der alten Orte. Treffer nur: Tombstone-Kommentar und -Einträge in `.d-check.yml`; Inline-Code in historischen Einträgen (`CHANGELOG.md:37`, ADR-0097, ADR-0103 — beide `Accepted`, `harness/conventions/done/MR-010`, eine Evidence-Datei); der Plan von slice-269 selbst; Test-Fixtures in `internal/hexagon/core/rules/vcs_pfad_nachzug_test.go` (Testdaten, keine Doku) |
| M7 | `make doc-check` | `1097 Datei(en) geprüft, 0 Befund(e)`, Exit 0 |
| M8 | `make adr-check RANGE=2b088bb1..HEAD` | `history-range-guard: … 7 Commit(s) — OK`; `1091 Datei(en) geprüft, 0 Befund(e)`, Exit 0 |
| M9 | `make gates` | `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`, Exit 0; `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%` |
| M10 | frischer Klon von `HEAD`; `ls -d docs/user/maintainer` | Verzeichnis existiert nicht (lokal ebenfalls nicht mehr) |
| M11 | im frischen Klon `d-check` wie `doc-check` und mit `--enable vcs` (alle anderen Module abgewählt, wie das Rezept) `--range 2b088bb1..HEAD` | `1097 … 0 Befund(e)` Exit 0 · `1091 … 0 Befund(e)` Exit 0 |
| M12 | Probe-Datei `docs/user/probe.md` im Klon mit Inline-Code `docs/user/releasing.md`, `docs/user/maintainer`, `docs/user/maintainer/releasing.md`, `docs/user/maintainer/` und Links `[x](maintainer/releasing.md)`, `[y](releasing.md)` | gemeldet: `docs/user/maintainer/releasing.md` `codepath-missing`; beide Links `target-missing`. Still: `docs/user/releasing.md`, `docs/user/maintainer`, `docs/user/maintainer/` |
| M13 | Bewusstes Brechen im Klon: Tombstone `docs/user/maintainer` aus `ignore-refs` entfernt, ein README-Link auf `docs/user/maintainer/releasing.md` zurückgesetzt | 3 Befunde: README.md:326 und :358 `target-missing` (der `sed` traf beide Links der Datei), slice-268 Zeile 106 `docs/user/maintainer/` `codepath-missing` — der Tombstone-Eintrag ist nötig und wirkt aus dem richtigen Grund; Klon danach per `git checkout -- .` zurückgesetzt |

## DoD-Punkte

### DoD 1 — Move-Commit, eigene Links, Liste gezählt, Gates grün auch im frischen Klon

**Verdikt: erfüllt, mit einem Form-Befund (V-1).**

- Rename erkannt, bewegte Datei unverändert (M1); eingehende Verweise im
  selben Commit (M2); eigene Links im Folge-Commit, nur Link-Ziele (M3, M4).
- 39 Link-Ziele: 23 eigene + 16 eingehende, alle auf dieselben Zieldateien
  bzw. `docs/maintainer/releasing.md` (M3, M5); am Stand lösen sie auf (M7,
  M11: 0 Befunde). Keine lebende Stelle verlinkt die alten Orte (M6).
- `make doc-check`, `make adr-check RANGE=2b088bb1..HEAD`, `make gates` grün
  (M7–M9); im frischen Klon ohne das leere Verzeichnis ebenfalls (M10, M11).
- Bewusstes Brechen (M13): ohne den Tombstone-Eintrag und mit einem alten
  Link-Ziel wird `doc-check` aus genau den erwarteten Gründen rot. Die
  Gate-Aussage hängt damit nicht an einem leeren Prüfbereich.
- **V-1:** Die DoD verlangt „die Liste am Repo gezählt, **das Kommando im
  Plan**". §3 beschreibt die Messung in Prosa („Gemessen im Wegwerf-Klon …
  Befunde nach Grund und Datei gezählt") und nennt die Zahlen, aber **kein
  ausführbares Kommando**. Die Zahl ist richtig (M5 reproduziert 40 genau),
  die Form der Zusage nicht erfüllt.

### DoD 2 — Review durchgeführt, Report unter `docs/reviews/`; Verifikation

**Verdikt: erfüllt.** R1 (`c81f3545`) und R2 (`cbcccbce`) liegen unter
`docs/reviews/`; R2: HIGH 0 · MEDIUM 0 · LOW 1 · INFO 2, alle in `7e72816e`
aufgenommen. Die Verifikation ist dieser Bericht.

### DoD 3 — Closure-Notiz, Register, Risiko-Ausgänge, Paarungen

**Verdikt: offen, wie erwartet** — Closure-Arbeit vor dem `git mv` nach
`done/`, nicht Gegenstand dieser Verifikation. Hinweis für den Ausgang des
Risikos aus §6: M5 belegt, dass es real war (lokal 39, ohne Verzeichnis 40
Befunde); M10/M11 belegen, dass die Prüfung im frischen Klon es abfängt.

## Zusätzliche Prüfpunkte des Auftrags

### Source Precedence und MR-077

**Verdikt: erfüllt.**

- `AGENTS.md` §2 Rang 6 nennt `docs/user/` (Operations) und
  `docs/maintainer/` (Releasing) mit Verweis auf MR-077; `harness/README.md`
  §Source precedence Rang 6 nennt beide Verzeichnisse, MR-009-Bezug bleibt.
- `harness/conventions/MR-077-releasing-doku-unter-docs-maintainer.md` trägt
  Datum, Geltungsbereich, Ersetzt-Baseline-Regel (verlinkt, pin-gebunden
  auf `v6.17.0`), Adaption, Grenze, Begründung, Auflösungs-Trigger — dieselben
  Felder wie der Nachbar MR-076. Index-Zeile in `harness/conventions.md`
  §Aktive Adaptionen mit `<a id="mr-077">`; die Anker lösen auf (M7).
- Tombstone-Einträge `docs/user/releasing.md` und `docs/user/maintainer` und
  ihr Grenz-Kommentar stimmen mit dem gemessenen Verhalten (M12): die beiden
  ausgenommenen Pfade bleiben in Inline-Code still, der volle Zwischen-Pfad
  `docs/user/maintainer/releasing.md` wird weiter gemeldet, Links meldet
  `links`. Ergänzend gemessen: auch die Schreibweise mit Schrägstrich
  (`docs/user/maintainer/`) bleibt still — sie ist dasselbe Verzeichnis, der
  Kommentar deckt sie (V-2, INFO).

### Release-relevante Stellen

**Verdikt: erfüllt.** `release.yml` (Kopfkommentar und `::error::`-Meldung),
`hub-description.yml` (`::warning::`-Meldung), `tools/image-test.sh`
(`fail`-Meldung) und `.github/dependabot.yml` (Kommentar) nennen
`docs/maintainer/releasing.md` (Diff `877af2d4`; M6 ohne Resttreffer).
`matrix.exempt-paths` nennt den neuen Pfad.

### Plan vs. Code

**Verdikt: keine Lieferung außerhalb von §1.**

- `docs/maintainer/` enthält nur `releasing.md` (Ausschluss 1 gehalten).
- `releasing.md` inhaltlich unverändert, nur Link-Ziele (M3, M4; Ausschluss 2).
- slice-268 und seine Reports in der Range nicht berührt
  (`git diff --stat 2b088bb1..HEAD`; Ausschluss 3).
- `AGENTS.md`, `harness/README.md`, MR-077 kamen per Plan-Änderung in §3
  (`1633d299`) **vor** dem Code-Commit `28552fcb` — Plan geändert, nicht still
  ausgeweitet. Das §1-Ziel nennt die Rang-6-Änderung nicht (bereits R2 INFO-2),
  die §3-Tabelle führt sie.

## Befunde

| ID | Kategorie | Befund | Ort |
|---|---|---|---|
| V-1 | LOW | DoD 1 verlangt das Zähl-Kommando im Plan; §3 beschreibt die Messung nur in Prosa. Vor dem DoD-Haken das Kommando nachtragen (z. B. die Form aus M5) oder den Haken mit dieser Lücke begründen | slice-269 §3, Absatz „Gemessen im Wegwerf-Klon …" |
| V-2 | INFO | Der Tombstone `docs/user/maintainer` lässt auch `docs/user/maintainer/` in Inline-Code still; der Grenz-Kommentar („das Verzeichnis") deckt das, kein Handlungsbedarf | `.d-check.yml` `codepaths.ignore-refs` |

## Kategorie-Summary

HIGH 0 · MEDIUM 0 · LOW 1 · INFO 1

## Verdikt

DoD 1 im Gegenstand erfüllt, in der Form ein offener Punkt (V-1); DoD 2
erfüllt; DoD 3 offen bis zur Closure. Source Precedence, MR-077, Tombstone und
Release-Stellen ohne DoD-Verletzung; keine Lieferung außerhalb von §1.
