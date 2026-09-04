# Review-Report — slice-197 (Wellenlosen Review-Bestand archivieren)

- **Review-Art:** Code (gegen Plan, `AGENTS.md` Hard Rules, Baseline-Regelwerk `modul-06-roadmap.md`)
- **Gegenstand:** drei Commits — `d057f19` (Werkzeug-Fix `tools/archive-wave/{archive.go,main.go,stub.go,slice_mode_test.go}`, Nachtrag zu slice-196), `0fe76a7` (45 wellenlose Slices archiviert, ~215 Dateien, `.d-check.yml`-Ignore-Refs-Tombstone), `ba750b7` (Closure-Body slice-197)
- **Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md)
- **Modell-ID:** claude-sonnet-5
- **Datum:** 2026-09-04
- **Eingangs-Kontext:** [`AGENTS.md`](../../AGENTS.md) vollständig (§3.1, §3.3/MR-013, §3.4, §3.5, §3.7, §3.8, §4, §5); Baseline-Regelwerk `modul-05-planning-harness.md` (Closure-/Risiko-Ausgangs-Regeln), `modul-06-roadmap.md` (§Wann Arbeit eine Welle braucht, §Wellen-Closure-Prozedur), `modul-08-agentenrollen.md`. Kein Zugriff auf Implementer-Kontext (Rollentrennung, Modul 8). Vorheriger Review-Report `2026-09-04-slice-196-archive-wave-slice-modus-code-r1.md` (jetzt archiviert, `slice-196-archiv.zip`) als Referenz für den F-1-Kommentar-Provenienz-Befund, den `d057f19` vermeiden musste.

Repo unverändert nach dem Review (`git status` leer, HEAD `ba750b7`).

---

## Eigener Lauf (Ausgabe, nicht behauptet)

| Lauf | Ausgabe |
|---|---|
| `make gates` | `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; semgrep `Ran 55 rules on 62 files: 0 findings`; `d-check: 595 Datei(en) geprüft, 0 Befund(e)` (×2, `targets`/`planning`); Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` — deckt sich mit der Commit-Botschaft (zehn Gates, 595 Dateien) |
| `make fullbuild` | `image-test: OK — DC-FA-DIST-001-Akzeptanzkriterien erfüllt`; `bench: OK` (Median 1172 ms); RTM `51 Anforderung(en), 0 Waise(n)`; Closure-Profil `d-check: 526 Datei(en) geprüft, 0 Befund(e)`; Schlusszeile `[fullbuild] green — image-hash sha256:1c9a353425ff…` — deckt sich mit „526 Dateien, 0 Befunde" |
| `make archive-wave-test` (`docker build --target test`, Layer-Cache) | `RUN CGO_ENABLED=0 go test ./...` als cache-hit-Layer terminiert erfolgreich |
| `docker build --no-cache --target test` (unabhängig vom Make-Cache erzwungen) | `ok  archive-wave  0.022s` — frischer, unabhängiger Beweis, kein Cache-Vertrauen nötig |
| `make adr-check RANGE=5e57c00..0fe76a7` | `d-check: 595 Datei(en) geprüft, 0 Befund(e)` — die vielen Geschichte-Verweis-Rewrites in `Accepted`-ADRs (0060/0061/0064/0065/0066/0067) verletzen die Immutable-Prüfung nicht |
| Byte-Diff Zip-Inhalt vs. `git show 0fe76a7~1:<pfad>` für alle 45 archivierten Slices | 22/45 identisch, 23/45 differieren **nur** in Cross-Referenz-Pfaden (`../done/slice-X.md` → `wellenlos/slice-X.md`) — Ergebnis der sequentiellen Verarbeitung: `RewriteRepo` zieht bei jedem Einzel-Lauf bereits archivierte Geschwister-Links nach, bevor der jeweils nächste Slice selbst archiviert wird. Kein Byte-Verlust, keine Bedeutungsverschiebung in diesen 23 Fällen — stichprobenartig an 3 Dateien (`slice-145`, `slice-163`, `slice-159`) Zeile für Zeile verifiziert |
| Titel-Vergleich Original-H1 vs. Stub-H1, alle 45 | 0 Abweichungen |
| **Welle-Feld-Vergleich Original vs. Stub, alle 45** | **44/45 mit mehrzeiligem Original-Feld → Stub trägt nur die erste Zeile, mitten im Satz abgeschnitten** (siehe F-1) |

---

## Findings

### F-1 · HIGH · Datenfidelität / stille Bedeutungsverschiebung — genau die vom Slice selbst benannte Top-Risikoklasse (§5 „stille Bedeutungsverschiebung wie bei einer Datenmigration") · `tools/archive-wave/collect.go:14,21-29` (`welleFieldRE`, `ReadWelleField`), ausgewirkt in 44 von 45 Stubs unter `docs/plan/planning/done/wellenlos/`

**Befund:** `welleFieldRE = regexp.MustCompile(`(?m)^\*\*Welle:\*\*\s*(.*)$`)` (unverändert von `d057f19`, aber erstmals am echten Bestand angewendet) erfasst nur die **erste Zeile** des `**Welle:**`-Feldes. Der Hausstil dieses Repos schreibt für wellenlose Slices jedoch **grundsätzlich mehrzeilige Begründungs-Prosa** in dieses Feld (Template `archiv-stub-slice.template.md` sieht dafür nur einen kurzen Wert vor, `<welle-id | ohne Welle>` — die tatsächliche Praxis weicht seit Langem davon ab). Gemessen: von den 45 archivierten Slices hat **44** ein mehrzeiliges Original-Feld; **alle 44** Stubs enden mitten im Satz. Beispiele (Original → Stub):

- slice-095: „bei Start zu eröffnen. Ein Slice „ohne Welle" ist in diesem Repo **nicht einlösbar: `make planning-check` koppelt Ruhe-Marker und `in-progress/` atomar (…), ein Slice in Arbeit verlangt also eine benannte Welle — auch wenn er inhaltlich…**" → Stub: „bei Start zu eröffnen. Ein Slice „ohne Welle" ist in diesem Repo" (bricht mitten im Satz ab, das Gegenteil der eigentlichen Aussage bleibt unausgesprochen stehen).
- slice-083: „keiner Welle zugeordnet — die Roadmap (…) **steht in Ruhe (welle-66 abgeschlossen); die Einplanung ist Teil der Abnahme (§5).**" → Stub: „keiner Welle zugeordnet — die Roadmap (…)" (kein Verb, kein Subsatz — unvollständiger Gedanke).
- slice-121, 165, 186, 189, 095, 143–182 (Stichprobe: 44 von 45) zeigen dasselbe Muster.

Der Autor des Werkzeugs teilt selbst dieses (falsche) Modell des Feldes: `SliceStub`s Doc-Kommentar sagt ausdrücklich „welleField ist der Wert des urspruenglichen **Welle:**-Feldes (z. B. `welle-87` oder `ohne Welle`)" — beides kurze, einzeilige Beispiele. Im **Wellen-Modus** manifestiert sich der Bug nie, weil ein welle-gebundener Slice dort tatsächlich nur einen kurzen Wert (`welle-87` oder einen Link) trägt (verifiziert an `welle-73`/`welle-83`/`welle-85`-Archiven). Er trifft ausschließlich den wellenlosen Modus, wo der Hausstil lange Prosa erzwingt — und der **eigene Test-Fixture** (`buildSliceFixture`, `slice_mode_test.go:17`) verwendet mit `"**Welle:** — wellenlos."` bewusst oder unbewusst genau die kurze Form, die den Bug nicht auslöst. Deshalb ist er weder in slice-196 (konstruiertes Fixture) noch in slice-197s eigenem Dry-Run/Gate-Lauf aufgefallen: kein Gate prüft Prosa-Kohärenz, nur Struktur (Links, IDs, Abschnitte).

**Bezug zur eigenen Risiko-Bewertung des Slice (§5 / §9):** Das Slice benennt genau diese Fundklasse als „Größtes Risiko" und behauptet in §9 („Was anders lief") sowie implizit in §5 ausdrücklich, sie sei bis auf den bereits gefundenen und korrigierten Design-Fehler (flacher Stub-Pfad) **nicht** eingetreten („entfallen, mit Korrektur unterwegs"). Das ist eine zweite, unabhängige Instanz **derselben** Fundklasse, die weder der vorgelagerte Dry-Run-Scan noch `make gates`/`make fullbuild` (die nur Struktur, nicht Prosa-Kohärenz prüfen) noch der Zwischen-Test erkannt hat — und die die Closure-Notiz nicht nennt.

**Grenze der Schwere:** Kein Datenverlust im strengen Sinn — der Volltext bleibt byte-identisch im jeweiligen `<id>-archiv.zip` erreichbar (siehe Eigener Lauf, Zeile „Byte-Diff"). Aber der Stub selbst — das einzige Artefakt, das ein Leser ohne `unzip` sieht — ist in 44 von 45 Fällen eine grammatikalisch unvollständige, teils sinnentstellende Paraphrase des ursprünglichen Feldes, nicht die vom Template zugesagte Kürzung.

**Verifizierbar:** ja — Skript-Vergleich Original-`**Welle:**`-Block (git-Historie vor `0fe76a7`) gegen Stub-`**Welle:**`-Zeile für alle 45 IDs, Ergebnis 44/45 abgeschnitten, 1/45 (`slice-140`, bereits einzeilig im Original) unbetroffen.
**Klasse:** `mechanische-migration-verschluckt-mehrzeiliges-feld`

**Empfehlung:** `welleFieldRE`/`ReadWelleField` auf einen absatzweiten Capture umstellen (bis zur nächsten Leerzeile statt bis zum Zeilenende) — dieselbe Technik, die `RewriteFieldForMove` bereits für Links im Feld beherrscht, träfe dann auch den Text davor/danach. Die 44 betroffenen Stubs brauchen danach eine Korrektur-Iteration (Werkzeug erneut mit korrigierter Extraktion über die bereits archivierten Fälle laufen lassen, oder die Stubs von Hand aus den bereits vorhandenen Zip-Inhalten nachziehen) **vor** dem `git mv` von slice-197 nach `done/` — sonst schließt der Slice mit einem unaufgelösten Fund genau der Klasse, die er selbst als größtes Risiko führt.

### F-2 · LOW · DoD-Formulierungs-Ungenauigkeit (§6 Dokumentations-Regeln: „behauptet nicht mehr, als die Arbeit trägt") · `docs/plan/planning/in-progress/slice-197-wellenlosen-bestand-archivieren.md:63-65` (§4 DoD, erster Punkt)

**Befund:** Die DoD listet den Bestand als „slice-083, 095, 102, 112, 121, 127, 137–147, 151–167, 171, 176–182, 185–187, 189" — der Bindestrich-Bereich „137–147" suggeriert eine lückenlose Folge, tatsächlich fehlt darin `slice-141` (bleibt korrekt in `done/`, weil kein zugehöriger Review-Report mehr in `docs/reviews/` lag — legitim außerhalb des in §1/§2 definierten Sammel-Kriteriums). Die Gesamtzahl (45) und die einzelnen IDs, wo ausgeschrieben, stimmen; nur die Bereichs-Klammerung „137–147" ist eine minimal irreführende Verdichtung, die eine Kontiguität behauptet, die nicht besteht.

**Verifizierbar:** ja — `find docs/plan/planning -iname "*slice-141*"` zeigt die Datei weiterhin unter `docs/plan/planning/done/` (nicht unter `wellenlos/`); `ls docs/reviews | grep 141` liefert keinen Treffer, bestätigt die legitime Nichtaufnahme.
**Klasse:** `bereichs-notation-verschleiert-luecke`

---

## Negativbefunde (geprüft, ohne Befund)

- **`ApplySlice`-Reihenfolge korrekt.** `buildZip` archiviert den Original-Volltext, **bevor** `os.ReadFile`/`os.WriteFile` den Stub erzeugen und das Original per `os.Remove` gelöscht wird (`archive.go:143-171`) — bestätigt sowohl im Code als auch am Bestand: alle 45 Zips enthalten den vollständigen (Stand-zum-Zeitpunkt-des-eigenen-Laufs-)Volltext, niemals den Stub.
- **Move-Rückgabe und Verweis-Nachzug korrekt.** `ApplySlice` liefert `[]Move{{Old, New}}`; `runSlice` reicht diesen an `RewriteRepo` weiter (Apply) bzw. `PreviewRewrites` (Dry-Run, schreibt nachweislich nichts — vom neuen Test `TestRunSlice_DryRun_NoWrites` mit Snapshot-Vergleich abgesichert). Repo-weite Rewrites (ADR-Geschichte-Tabellen, Welle-Results-Dateien, MR-Konventionsdateien, andere Review-Reports) sind stichprobenartig geprüft und durchweg korrekt: alter Pfad `slice-<NNN>-….md` bzw. `../done/slice-<NNN>-….md` wird zu `wellenlos/slice-<NNN>-….md` bzw. `done/wellenlos/slice-<NNN>-….md`, relativ zum jeweiligen Quelldokument neu berechnet.
- **Kein neuer §3.7-Kommentar-Herkunfts-Verstoß.** Die Vorgänger-Review (`slice-196` R1, F-1 HIGH) hatte fünf neue Kommentare mit Slice-Nummer-Herkunfts-Prosa gefunden; `d057f19` vermeidet dieses Muster vollständig — die neuen/geänderten Doc-Kommentare (`WellenlosArchiveDir`, `ApplySlice`, `SliceStubStandalone`) beschreiben Zusage/Abgrenzung direkt, ohne Slice-Nummer als Ursprungs-Erzählung. Die zwei neuen Kommentare in `slice_mode_test.go` ebenso.
- **`.d-check.yml`-Tombstone korrekt und notwendig.** `docs/reviews/2026-08-31-release-prep-v0.71.0-review.md` zitiert nachweislich (`grep`) genau die beiden Pfade `2026-08-31-slice-185-code-r1.md`/`-verifikation.md`; beide existieren nicht mehr unter `docs/reviews/`, sind aber byte-identisch im neuen `slice-185-archiv.zip` enthalten. Form und Platzierung des neuen Eintrags entsprechen exakt dem etablierten Schwester-Muster (`docs/reviews/2026-08-31-slice-174-register-deckung-code-r1.md`-Eintrag, dieselbe Begründungsklasse „Review-Report ist Lauf-Beleg, wird nicht editiert").
- **ADR-Immutabilität nicht verletzt.** Sechs `Accepted`-ADRs (0060/0061/0064/0065/0066/0067) bekommen Pfad-Updates in ihren Geschichte-Tabellen (`Verweis`-Spalte) — `make adr-check` über die genaue Commit-Range bestätigt 0 Befunde; die Kern-Aussagen der Einträge sind unverändert, nur die Ziel-Pfade der bereits verlinkten Slices folgen deren Move.
- **§3.3/MR-013-Bündelungs-Begründung trägt.** Die Ausnahme „ein Commit statt 45" ist keine der fünf explizit benannten §3.3-Ausnahmen wörtlich, aber MR-059 (Wellen-Archiv-Stub-Move ist ein Commit) trägt bereits die Prämisse, dass Move+Stub-Ersatz je Vorgang atomar in einem Commit liegen darf — die verbleibende Frage (dürfen mehrere solcher Ein-Commit-Vorgänge zu einem Commit gebündelt werden) ist eine Repo-Ermessensfrage ohne eigenes Gate, im Slice-Plan selbst vorab als Entscheidungsoption benannt (§2 Punkt 3) und in der Commit-Botschaft mit nachvollziehbarem Grund versehen (mechanische, werkzeug-verifizierte Operation; Gates liefen gegen den Endstand). Konsistent mit dem bereits akzeptierten Präzedenzfall der Wellen-Archivierungen.
- **§3.1 Docker/make-only, soweit rekonstruierbar.** Kein Hinweis auf einen rohen Host-`go`/`docker`-Aufruf in Commit-Botschaften oder Plan; die genutzte CLI-Form (`tools/archive-wave -slice=<id> -apply`) beschreibt exakt die Flags, die `make archive-wave SLICE=<id> APPLY=1` an das Docker-Image durchreicht (`tools/archive-wave/Makefile:run`), und ist derselbe Beschreibungsstil, den slice-196 bereits für dasselbe Werkzeug verwendet hat. Nicht beweisbar aus der Historie allein, aber keine Gegenevidenz.
- **`docs/reviews/`-Editier-Fälle sind Redirects, keine Tombstone-würdigen Verstöße.** Zwei bereits bestehende Review-Reports (`2026-08-09-backlog-schnitt-review.md`, `2026-08-30-release-prep-v0.70.0-review.md`) werden in `0fe76a7` inhaltlich verändert (Link-Pfad-Fix auf `slice-095`/`slice-178`). Das widerspricht der im selben Commit zitierten Regel „Review-Report wird nicht nachträglich editiert" nicht: jene Regel gilt Links auf **andere, ersatzlos gelöschte Review-Reports** (kein Redirect-Ziel vorhanden); hier existiert das Ziel (der Slice) weiterhin unter neuem Pfad — ein Redirect ist die korrekte Behandlung, konsistent mit der Behandlung in ADRs/MR-Dateien/Wellen-Results oben.
- **DoD/Closure-Zahlen (bis auf F-2) korrekt.** „45 Slices" (nicht 43) — exakt nachgezählt (`D`-Zeilen `docs/plan/planning/done/slice-*.md` = 45, `A`-Zeilen `wellenlos/*.md` = 45, `A`-Zeilen `wellenlos/*.zip` = 45). „`make gates` grün (zehn Gates, 595 Dateien)" und „`make fullbuild` grün (526 Dateien, 0 Befunde)" — beide unabhängig reproduziert, exakte Zahlen bestätigt.
- **WIP-Limit.** `docs/plan/planning/in-progress/` enthält nur `slice-197` (und `roadmap.md`); kein Parallel-Slice.

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 1 | F-1 |
| MEDIUM | 0 | — |
| LOW | 1 | F-2 |
| INFO | 0 | — |

**Finding-Klassen dieses Laufs:** `mechanische-migration-verschluckt-mehrzeiliges-feld` · `bereichs-notation-verschleiert-luecke`

## Verdikt

**Merge-blockierend: ja (F-1).**

Der Werkzeug-Fix selbst (`d057f19`) ist korrekt und sauber: `ApplySlice` baut das Zip nachweislich vor jeder Mutation, verschiebt korrekt nach `docs/plan/planning/done/wellenlos/`, liefert einen echten `Move` für den repo-weiten Verweis-Nachzug, löscht Review-Reports weiterhin ersatzlos, und `runSlice` verdrahtet `RewriteRepo`/`PreviewRewrites` richtig (durch einen neuen, gezielten Test mit Vorher/Nachher-Snapshot abgesichert). Die 45 Einzelanwendungen sind mechanisch korrekt durchgeführt: Titel, Hervorgegangen-Kennungen und Zip-Inhalte sind über alle 45 Fälle hinweg fehlerfrei bzw. durch die erwartete, korrekte sequentielle Verweis-Aktualisierung erklärbar; die neue `.d-check.yml`-Tombstone ist notwendig und formgleich zum etablierten Muster; alle drei eigenen Gate-Läufe (`make gates`, `make fullbuild`, `make archive-wave-test` inkl. eines `--no-cache`-erzwungenen frischen Testlaufs) sind unabhängig grün reproduziert.

**F-1 ist der blockierende Befund:** Das Slice benennt „stille Bedeutungsverschiebung wie bei einer Datenmigration" explizit als sein größtes Risiko und behauptet in §5/§9, dieses Risiko sei — abgesehen vom bereits gefundenen und korrigierten Pfad-Kollisions-Fehler — nicht eingetreten. Tatsächlich tritt es ein zweites Mal ein, in 44 von 45 archivierten Stubs, unbemerkt von jedem gelaufenen Gate (die vorgelagerte Dry-Run-Prüfung, `make gates` und `make fullbuild` prüfen Struktur/Links/IDs, keine Prosa-Kohärenz eines Freitext-Felds). Der Ursprung ist ein einzeiliger Regex (`welleFieldRE`), der die Praxis dieses Repos (mehrzeilige Begründungsprosa im `**Welle:**`-Feld) nicht abbildet — dieselbe Fixture-vs-Bestand-Lücke, die das Slice selbst als Lerneintrag für die *andere* gefundene Kollision festhält, hier aber ein zweites Mal und unentdeckt. Datenverlust im strengen Sinn liegt nicht vor (Volltext bleibt im jeweiligen `.zip` erreichbar), aber der Stub — das Artefakt, das ohne `unzip` gelesen wird — ist in 44 von 45 Fällen eine mitten im Satz abgebrochene, teils sinnentstellende Paraphrase. Empfehlung: `ReadWelleField` auf einen absatzweiten statt zeilenweiten Capture umstellen und die 44 betroffenen Stub-Dateien vor dem `git mv` von slice-197 nach `done/` nachziehen (Werkzeug erneut anwendbar, da die Zip-Archive den vollständigen Originaltext tragen).

F-2 ist eine kosmetische Ungenauigkeit in der DoD-Kurzschreibweise, für sich genommen nicht merge-blockierend.
