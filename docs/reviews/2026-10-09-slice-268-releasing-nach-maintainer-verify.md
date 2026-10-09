# Verifikation: slice-268, `releasing.md` zieht nach docs/user/maintainer/

- **Rolle:** Verifier (Modul 8/11). Geprüft wird der Stand gegen DoD (§2), Plan (§3 samt
  Plan-Änderung vor dem Code), Abgrenzung (§1) und §6, **nicht** der Diff gegen die
  Entscheidungen. Das hat R1 getan.
- **Gegenstand:** `slice-268`, Commits `bcb104fe`, `420e401a`, `1be76fc2`, `2f6a9c49`
  (`git log a9d62bab..HEAD`, alle tragen `slice-268`). Code-Stand `HEAD` = `2f6a9c49`,
  Arbeitsbaum sauber, `origin/main..HEAD` leer.
- **Eingang:** Slice-Plan in `in-progress/`; Report R1; Spezifikation `DC-FA-VCS-001.a` Schritt 4
  (`vcs.ignore-link-targets`); [ADR-0025](../plan/adr/0025-codepaths-ignore-refs.md),
  [ADR-0103](../plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md); `AGENTS.md` §3.3;
  `.githooks/pre-commit`.
- **Datum:** 2026-10-09 · **Modell-ID:** claude-opus-5-5

## Verdikt

**Der Liefer-Punkt der DoD (§2, Punkt 1) ist in der Sache bestätigt:** Die Datei zog als
`R100` um, jeder eingehende Verweis zeigt auf die neue Datei, keine lebende Stelle nennt den
alten Pfad noch als Verweis, und alle drei Gates sind grün (selbst gefahren). Die Begründung,
warum die eingehenden Verweise im Move-Commit mitreisen, trägt: durch bewusstes Brechen
bestätigt. **Abweichung vom Wortlaut** der DoD und von §1 (Commit-Schnitt), im Plan nicht
nachgetragen (V-1, deckungsgleich mit R1 LOW-1, noch unbeantwortet). Punkt 2 ist mit R1 und
diesem Bericht erfüllt; der Fix-Commit `2f6a9c49` liegt hinter R1 (V-3). **Offen ist Punkt 3**,
die Closure. Keine DoD-Verletzung in der Sache. Ein Befund LOW, sonst INFO.

## Sensor-Belege (selbst gefahren)

| Lauf | Ergebnis |
|---|---|
| `make doc-check` (auf `2f6a9c49`) | Exit 0, `d-check: 1085 Datei(en) geprüft, 0 Befund(e)` |
| `make adr-check RANGE=a9d62bab..HEAD` | Exit 0; `history-range-guard: Range 'a9d62bab..HEAD' aufgeloest, 4 Commit(s) — OK.`; `d-check: 1079 Datei(en) geprüft, 0 Befund(e)` |
| `make gates` (auf `2f6a9c49`) | Exit 0. Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; semgrep `Ran 55 rules on 66 files: 0 findings.`; a-check `gesamt: 0 Befund(e)` |
| CI (`gh run list`) | `ci` grün auf `a9d62bab`, `420e401a`, `2f6a9c49`; **kein** Lauf auf `bcb104fe` — der Zwischenstand mit den noch alten eigenen Links war nie Push-Spitze |

## DoD Punkt 1: Move, Nachzug, Zählung, Gates

**Rename.** `git show -M --summary bcb104fe` → `rename docs/user/{ => maintainer}/releasing.md (100%)`;
die Datei in `bcb104fe` ist byte-identisch mit `a9d62bab:docs/user/releasing.md` (`cmp`).

**Eigene Links (`420e401a`).** 23 Ziele, alt relativ zu `docs/user/`, neu relativ zu
`docs/user/maintainer/` aufgelöst (`realpath -m`): **alle 23 zeigen auf dieselbe Datei wie
vorher** — keine Verwechslung gleichnamiger Ziele
(`BEO-ALL/pfad-nachzug-gleicher-name-andere-datei`). Mit Link-Zielen maskiert
(`sed 's/](…)/](X)/'`) sind alte und neue Fassung identisch: **keine Inhaltsänderung**.

**Eingehende Links (`bcb104fe`).** 16 Link-Ziele in 10 Dateien, je aufgelöst: alle 16 →
`docs/user/maintainer/releasing.md` (CHANGELOG 1, README 2, README.de 2, Benutzerhandbuch 2,
`packaging/dockerhub/README.md` 1, ADR-0014 4, ADR-0067 1 in `## Geschichte`, Wellen-Ergebnisse
73/74/82 je 1).

**Zählung nachgefahren.** Wegwerf-Klon auf `a9d62bab`, nur `git mv` committet, dann
`make doc-check`: Exit 2, `58 Befund(e)` = **39 `target-missing`** (23 eigene + 16 eingehende)
+ **19 `codepath-missing`**, alle 19 auf `docs/user/releasing.md`. Exakt die Zahlen der
Plan-Änderung §3.

**Restbestand des alten Pfads** (`git grep` auf `docs/user/releasing`, `user/releasing.md`,
`(releasing.md`, ohne `maintainer/`): **keine lebende Stelle nennt ihn als Verweis.** Übrig sind
der Tombstone selbst (`.d-check.yml`), Inline-Code in eingefrorenen Dateien (ADR-0097, ADR-0103,
`done/`-Slices 230/236/248/256/267, ein Register-Beleg, MR-010, CHANGELOG-Abschnitt 0.84.0 —
veröffentlichter Eintrag), Review-Reports, synthetische Test-Fixtures
(`vcs_pfad_nachzug_test.go`) und das Beispiel in Spezifikation Schritt 4. Keine absolute
GitHub-URL auf den alten Pfad im Repo. Workflows, Skripte, Konfiguration: alle nachgezogen
(unten).

**Tombstone.** Der Eintrag `docs/user/releasing.md` in `codepaths.ignore-refs` fällt unter
ADR-0025 (Tombstone für einen am alten Ort entfernten Pfad); die Grenze (repo-weite Wirkung)
steht seit `2f6a9c49` im Kommentar.

## Pre-commit-Weg und die Begründung der Commit-Botschaft

Behauptung `bcb104fe`: getrennt gestagt liefe der gestagte ADR-Check im zweiten Commit gegen
einen BASE-Stand ohne den alten Pfad. Gegen Spezifikation Schritt 4: normiert wird nur ein Ziel,
das in der **Vereinigung der Pfad-Bäume von BASE und HEAD** auflöst. Bei `--staged` ist BASE der
letzte Commit. Liegt der reine Move schon im Commit, existiert `docs/user/releasing.md` in keinem
der beiden Stände, das alte Ziel in ADR-0014 löst nicht auf, bleibt roh — Drift.

Nachgefahren in zwei Wegwerf-Klonen (`make adr-check STAGED=1`, der Weg des Hooks):

| Klon | Zustand | Ergebnis |
|---|---|---|
| A | auf `a9d62bab`, gesamter Diff von `bcb104fe` gestagt (wie geliefert) | `0 Befund(e)`, Exit 0 |
| B | auf `a9d62bab` reiner `git mv` committet, dann nur der ADR-0014-Nachzug gestagt | `1 Befund(e)`: `docs/plan/adr/0014-latest-tag-fuer-stabile-releases.md:3 … core-drift-vcs`, Exit 2 |

Die Begründung trifft am Gegenstand zu; der gelieferte Schnitt ist der, den der Hook durchlässt.
Der Range-Lauf in CI wäre in beiden Schnitten grün — betroffen ist nur der lokale Hook.

## Release-relevante Stellen

- `.github/workflows/release.yml` Zeile 11 (Kommentar) und 191 (`::error::`-Meldung),
  `.github/workflows/hub-description.yml` Zeile 106 (`::warning::`), `tools/image-test.sh`
  Zeile 71 (`fail`-Meldung), `.github/dependabot.yml` Zeile 23: alle nennen
  `docs/user/maintainer/releasing.md`.
- **Docker-Hub-Upload unberührt:** `hub-description.yml` lädt `packaging/dockerhub/description.txt`
  und `overview.md`; `git diff a9d62bab HEAD -- packaging/` zeigt nur `packaging/dockerhub/README.md`
  (Repo-Rendering). `overview.md` nennt `releasing` nicht.
- `matrix.exempt-paths` trägt den neuen Pfad (sonst `matrix-inactive` auf der Datei).

## Plan-vs-Code-Diff

- **§3 geliefert:** Move, alle Verweise; die Plan-Änderung (Tombstone statt Edit in
  Inline-Code, die gate-losen Stellen) eingelöst.
- **§1-Abgrenzung gehalten:** `docs/user/maintainer/` enthält nur `releasing.md`; an
  `releasing.md` ändern sich nur Link-Ziele (maskierter Vergleich identisch).
- **Abweichung:** der Commit-Schnitt (V-1).
- **Rückführung aus §4** (`adr-check` lässt den ADR-Nachzug nicht durch) nicht eingetreten.

## Befunde

### V-1 — LOW: Commit-Schnitt weicht von DoD-Wortlaut und §1 ab, der Plan nennt es nicht

DoD Punkt 1 und §1 sagen „reiner Move-Commit …, danach ein Commit, der jeden Verweis nachzieht".
Geliefert: Die 16 eingehenden Links und die gate-losen Stellen reisen im Move-Commit mit, der
Folge-Commit zieht nur die eigenen Links nach. Die Abweichung ist **begründet und bestätigt**
(Klon B oben), und die Rename-Erkennung ist unberührt (`R100`). Aber sie steht nur in der
Commit-Botschaft, nicht als Plan-Änderung (`AGENTS.md` §6 Schritt 4), und die Botschaft beruft
sich auf `AGENTS.md` §3.3, das Mitreisen ausdrücklich nur für Fall 2 (Übergang nach `done/`)
erlaubt. Der Zweck der Regel (Rename-Erkennung) trägt den Fall, der Wortlaut nicht. Deckungsgleich
mit R1 LOW-1, dort noch nicht angenommen oder begründet. **Vorschlag:** Bei der Closure unter
„Was ging anders als geplant" nennen, samt dem Grund (BASE ohne alten Pfad beim gestagten
Check); ob §3.3 den Fall 1 mit gekoppelten Verweisen ausdrücklich nennen soll, ist eine
Beobachtung fürs Register.

### V-2 — INFO: Das Zählkommando im Plan ist elidiert

§3 nennt `docker run … d-check:latest` mit Auslassung. Nachvollziehbar wurde die Zählung erst
mit `make doc-check` in einem Wegwerf-Klon mit reinem Move; dort stimmen 39 + 19 exakt. Die
DoD verlangt „das Kommando im Plan" — es steht der Form nach da, nicht reproduzierbar ausgeschrieben.

### V-3 — INFO: `2f6a9c49` liegt hinter R1

Der Commit fügt nur zwei Kommentarzeilen (Grenze des Tombstones, Antwort auf R1 LOW-2) in die
`.d-check.yml`; kein Wert ändert sich. `make gates` ist auf diesem Stand grün.

### V-4 — INFO: Risiko §6 für die Closure

Der Hub-Teil des Risikos trägt nicht (`overview.md` verlinkt `releasing.md` nicht, R1 INFO-2).
Übrig bleiben Adopter-Links auf den alten GitHub-Pfad — von hier nicht messbar.

## Hinweis zum Lauf

Die Klon-Läufe haben `d-check:latest` aus dem Stand `a9d62bab` neu getaggt (Go-Quellen in der
Range unverändert); danach im Repo `make build` gefahren, Exit 0.
