# Review slice-297 — Die Spezifikation heißt `specification.md`

- **Review-Art:** Code/Doku (gegen Plan samt beiden Plan-Änderungen, Architect-Verdikt, `ADR-0025`/`ADR-0103`/`ADR-0115`/`ADR-0107`, `MR-069`/`MR-070`/`MR-032`, Hard Rules)
- **Gegenstand:** `slice-297`, Diff `f44c6e65..HEAD` (Plan-Änderungen `e0a48bb9`, `8a458a27`; Move `56c8f62f`; lebende Verweise `340f3f15`; Selbst-Verweise `061ebfc2`; Lastenheft `501dc194`)
- **Skill:** `reviewer.md` @ 1.20.0
- **Modell-ID:** claude-sonnet-5-5
- **Datum:** 2026-11-26
- **Eingangs-Kontext:** `slice-297` (§1 Abgrenzung, §3 Plan samt Plan-Änderung 1 und 2, §6 Risiken), `wave-98` (§3 (d), §5), das Architect-Verdikt zu `wave-98` (Fragen 4 bis 7, Nachtrag V1), `ADR-0025`, `ADR-0103`, `ADR-0115`, `ADR-0107`, `MR-069`, `MR-070`, `MR-032`, `AGENTS.md` §3 und §5, der Report `2026-11-26-slice-295-baseline-v700-bump-r1` als Form-Vorbild. **Nicht** gelesen und nicht geprüft: DoD-Abhakung (Sache der Verifikation).
- **Zusätzlich gefahren (Wegwerf-Klon im Scratchpad, nichts im Repo):** Ventil-Eintrag entfernt (Zählung der Befunde) und je eine Mutationsprobe mit lebender Inline-Code-Nennung, lebendem Link und einem anderen toten Pfad; `doc-check` an den Ständen `56c8f62f`, `340f3f15`, `061ebfc2` (Zwischenstände).

## Findings

Keine HIGH. Ein MEDIUM, zwei LOW, drei INFO.

### F-1 — MEDIUM — Die Welle trägt die „fünfte Klasse“ von `--manual` nicht, die der Slice ihr zuschreibt

- **quelle:** Prüffrage 10 (Messmethode klafft gegen die Stelle, die sie erfüllen soll), `slice-297` Plan-Änderung 2 („die Welle führt sie als fünfte Klasse“), `AGENTS.md` §5 Regel 18 (eine Festlegung nennt die Stelle, die sie trägt)
- **pfad:** `docs/plan/planning/wave-98-baseline-v700.md` · „Eine fünfte Klasse ist ein Befund.“
- **befund:** Der Slice ändert die Prosa-Zeile 920 der Spezifikation (der Pfad zu `--manual`), die in der Ausgabe zu `manual` und `Spezifikation` steht, und hält fest, die Welle führe sie als fünfte Klasse von (d3). Der Wellen-Plan trägt weiter genau vier Klassen und den Satz „Eine fünfte Klasse ist ein Befund“ (Z. 233 bis 237, auch im Satz „Rot ohne die Welle“); die Zeile 920 kommt dort nicht vor. Dieselbe Welle behauptet in §3 Punkt 4 „**0** in lebenden Dokumenten“ hinter dem Ventil und „**104**“ eingefrorene Nennungen; am Stand `HEAD` stehen elf lebende Nennungen (drei Pläne, `MR-081`) hinter dem Ventil und 109 eingefrorene (siehe Negativbefunde, Ventil).
- **Failure-Szenario:** Die Gegenmessung (d3) am Endstand findet die Zeile 920 als Abweichung. Nach dem geschriebenen Wortlaut ist sie ein Befund, der Closure-Trigger (d) wäre rot; wer die Welle schließt, muss den Trigger im Moment der Closure umdeuten statt ihn abzulesen.
- **Hinweis zur Stufe:** Die Änderung der Zeile 920 selbst ist vertretbar (sonst nennt die Spezifikation eine Datei, die es nicht gibt, und der Vertrag von `DC-FA-CLI-013.a` Schritt 2 hielte den Pfad falsch). Gemeldet ist die offene Übergabe: Der Slice-Plan benennt einen Folgeakt am Wellen-Plan, den es noch nicht gibt.
- **verifizierbar:** ja — `git grep -n 'fünfte\|920' -- docs/plan/planning/wave-98-baseline-v700.md` (die erste Fundstelle ist der Satz „Eine fünfte Klasse ist ein Befund“, die Zeile 920 fehlt)
- **klasse:** folgeakt-benannt-ohne-traeger

### F-2 — LOW — `MR-006` nennt den Namen als Aussage über das Template-Set 2026-06, der Nachzug schreibt sie um

- **quelle:** `MR-070` (Eigenschafts-Test: „Würde ein korrigierter Wert verfälschen, was dieses Artefakt zu seinem Datum festgehalten hat?“)
- **pfad:** `harness/conventions/MR-006-referenzrichtung-matrix.md` · „in `specification.md` (ADR-Spalte in Defaults/Historie)“
- **befund:** Der Eintrag (Datum 2026-06-10) sagt, das adoptierte Template-Set 2026-06 habe ADR-Verweise „in `spezifikation.md`“ vorgesehen; `340f3f15` setzt dort `specification.md`. Das ist eine Vergangenheits-Aussage über einen Fremdbestand (die Namensregel „Dateinamen englisch“ gilt erst ab `v7.0.0`), kein lebender Verweis auf diese Datei; der Plan führt die Zeile unter den 49 geänderten lebenden Zeilen, ohne sie Zeile für Zeile gelesen zu haben (die anderen drei aktiven `MR` — `MR-025`, `MR-076`, `MR-0098` — beschreiben den heutigen Bestand und sind richtig nachgezogen).
- **verifizierbar:** ja — `git show 340f3f15 -- harness/conventions/MR-006-referenzrichtung-matrix.md`
- **klasse:** vergangenheitsaussage-ueber-fremdbestand-mitersetzt

### F-3 — LOW — Zahlen des Plans (119 Befunde, neun Pläne) stimmen am Endstand nicht mehr

- **quelle:** `AGENTS.md` §5 Regel 14 (Form vor Zahl), `slice-297` §6 (Risiko „Eine Zahl dieses Plans stimmt nicht“)
- **pfad:** `docs/plan/planning/in-progress/slice-297-spezifikation-heisst-specification.md` · „Das Ventil deckt 119 Befunde, nicht 104.“
- **befund:** Mit entferntem `ignore-refs`-Eintrag meldet `HEAD` 120 `codepath-missing` statt 119; die zehn Befunde in lebenden Klassen sind elf (drei Pläne: `wave-98` 4, `slice-296` 2, `slice-297` 4; dazu `MR-081` 1). Die 109 eingefrorenen Nennungen stimmen Klasse für Klasse. Die Abweichung ist die eigene Nennung des Plans (die zweite Plan-Änderung schreibt `spec/spezifikation.md` in Inline-Code und gehört damit selbst zur Menge). Die Tabelle in §3 trägt außerdem weiter 104, 123 und 161; die Plan-Änderungen überschreiben sie ausdrücklich, aber die Closure-Notiz sollte die Endzahl nennen.
- **verifizierbar:** ja — Wegwerf-Klon, Eintrag `"spec/spezifikation.md"` aus `codepaths.ignore-refs` entfernt, `d-check` im Image: `1324 Datei(en) geprüft, 120 Befund(e)`
- **klasse:** zahl-zaehlt-sich-selbst-mit

### F-4 — INFO — Lebende Gegenstand-Nennungen laufen über das repo-weite Ventil, `ADR-0025` führt dafür den Zeilen-Marker

- **quelle:** `ADR-0025` §Entscheidung (Klammersatz: der Zeilen-Marker `d-check:ignore` bleibt für historische Verweise in **lebenden** Dateien)
- **pfad:** `.d-check.yml` · „die Pläne der Welle und MR-081 als Gegenstand“
- **befund:** Elf der 120 Nennungen stehen in lebenden Dateien (Pläne, `MR-081`). Die zweite Plan-Änderung lässt sie hinter dem ohnehin nötigen repo-weiten Eintrag und nennt das im Kommentar; sie sagt nicht, warum nicht der Zeilen-Marker. Auswirkung gemessen null: Der Eintrag ist für die 109 eingefrorenen Nennungen ohnehin repo-weit, ein Marker verkleinerte die stille Fläche nicht.
- **verifizierbar:** nein — Designnotiz
- **klasse:** ventil-form-weicht-von-adr-klammersatz-ab-ohne-wirkung

### F-5 — INFO — Spezifikation §8 Historie führt keine Zeile zum Rename, das Lastenheft schon

- **quelle:** `MR-032` (gilt dem Lastenheft; für die Spezifikation trägt keine Regel eine Zeile)
- **pfad:** `spec/specification.md` · „## 8. Historie“
- **befund:** Die Spezifikation trägt einen redaktionellen Pfad-Wechsel an einem Algorithmus-Satz (Zeile 920, `DC-FA-CLI-013.a` Schritt 2) und ihren neuen Dateinamen, die §8 nennt beides nicht. Frühere redaktionelle Nachzüge standen dort (Zeilen vom 2026-08-23); der Pin-Wechsel von `slice-295` stand dort nicht. Ohne Konventions-Anker kein Finding; gemeldet, weil Lastenheft und Spezifikation ungleich behandelt sind.
- **verifizierbar:** ja — `git diff f44c6e65..HEAD -- spec/specification.md | grep -c 'Historie'` ist 0
- **klasse:** historie-zeile-nur-an-einem-stratum

### F-6 — INFO — Das repo-eigene Konsumenten-Beispiel nennt jetzt eine Datei, die dieses Repo nicht mehr hat

- **quelle:** `slice-297` §1 („Keine Konsumenten-Beispiele“), Plan-Änderung 2, `MR-025` (Spiegel ändern sich an allen Stellen oder an keiner)
- **pfad:** `spec/specification.md` · „paths: [spec/lastenheft.md, spec/spezifikation.md, spec/architecture.md]“
- **befund:** Das Beispiel steht an vier Spiegel-Stellen gleich (Spezifikation Z. 233 und 234, Handbuch Z. 372 und 373, `config_template.go` Z. 56 und 57, `suggest.go` Z. 379 und 380) und bleibt an allen vier; das ist die Auslegung der Plan-Änderung 2 und sie ist vertretbar. Offen bleibt, dass `--suggest-config` einem Repo mit Spec-Verzeichnis die deutschen Namen vorschlägt, dieses Repo sie aber nicht mehr trägt; die Sichtung (`slice-296`, Frage B) ist als Träger benannt.
- **verifizierbar:** ja — `git grep -n 'spec/spezifikation.md' -- '*.go' docs/user spec/specification.md` nennt genau diese acht Zeilen
- **klasse:** spiegel-beispiel-folgt-dem-eigenen-namen-nicht

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Move `56c8f62f` ist ein reiner Rename | geprüft, ohne Befund — `git show -M --name-status --format= 56c8f62f`: eine Zeile `R100 spec/spezifikation.md spec/specification.md`, dazu 103 `M` (64 ADRs, 39 weitere), kein `D`/`A`-Paar; Dateizahl 104 laut `--stat` |
| Link-Nachzug ist reine Substitution | geprüft, ohne Befund — je Commit die gelöschten Zeilen mit `sed 's/\(\](\([^)]*\/\)\?\)spezifikation\.md/\1specification.md/g'` normalisiert und gegen die hinzugefügten Zeilen gehalten, `diff` leer: `56c8f62f` ADRs 119 / 119, übrige eingefrorene 80 / 80 (zusammen 199 / 199, Spezifikation ausgenommen); `061ebfc2` 159 / 159 mit **einer** Ausnahme, der Prosa-Zeile 920 (Inline-Code, geplant); `340f3f15` nur die benannten Dateien, Abweichungen von der Normalisierung allein die Inline-Code-Nennungen lebender Dateien (Tabellen, `harness/conventions.md`, drei `MR`, Handbuch Z. 2792, `AGENTS.md`-Linktexte) |
| `adr-check` über die Spanne | geprüft, ohne Befund — `make adr-check RANGE=f44c6e65..HEAD`: `d-check: 1318 Datei(en) geprüft, 0 Befund(e)`, Exit 0; keine inhaltliche Änderung in eingefrorenen Dateien außer dem Link-Ziel |
| Wortgrenzen-Falle | geprüft, ohne Befund — `architektur-und-spezifikation.md` und `spec-vergabe-spezifikation.md`: 7 Zeilen vor wie nach, 0 Zeilen im Diff, beide Dateien vorhanden |
| Link-Ziele ganz erledigt | geprüft, ohne Befund — `git grep -c '\](\([^)]*/\)\?spezifikation\.md' \| wc -l`: `0`; `git grep -h` ebenso 0 Zeilen; Summe vor dem Slice 406 (`git grep -h … f44c6e65 \| wc -l`), 119 + 80 + 4 + 3 + 42 + 158 geht auf (die 42 am Diff von `340f3f15` nachgezählt) |
| Ventil `codepaths.ignore-refs` — Form | geprüft, ohne Befund — ein Eintrag, exakter Pfad statt Glob (per-Pfad, `ADR-0025`), Kommentar in der Form des Vorgängers `docs/user/releasing.md` (Zeile je Eintrag, danach „GRENZE: … bleibt das still … einen Link meldet links“), Herkunft ein Feld (`seit wave-98`); `AGENTS.md` §3.6: die Gate-Senkung trägt `ADR-0025`, das Verdikt zu Frage 5 hat sie ohne eigenen `MR` entschieden; `MR-069` gilt den Einträgen auf entfernte Baseline-Bäume und trägt diesen Eintrag nicht |
| Ventil — Wirkung, Mutationsprobe | geprüft, ohne Befund — Wegwerf-Klon am Stand `501dc194`: lebende Inline-Code-Nennung in `docs/user/operations.md` (`spec/spezifikation.md`, `../../spec/spezifikation.md`, `spezifikation.md`, `spec/spezifikation.md#spec-005`) bleibt still (`0 Befund(e)`), ein lebender Link darauf meldet `target-missing`, ein anderer toter Inline-Code-Pfad (`spec/spezifikationX.md`) meldet `codepath-missing`; die GRENZE im Kommentar trifft das Verhalten |
| Ventil — Zahl 119 gegen 104 des Verdikts | geprüft, mit F-3 — Eintrag entfernt: 120 Befunde; eingefroren 109 = `done/` 51, ADRs 44, Register 5, `conventions/done` 6, `CHANGELOG.md` 2, `docs/plan/cr` 1 (stimmt mit der Plan-Änderung); die Abweichung zu 104 liegt in `done/` (+4) und `conventions/done` (+1) gegenüber dem Probelauf der Welle; mit Ventil 0 (`make doc-check`: `1324 Datei(en) geprüft, 0 Befund(e)`, Exit 0) |
| Auslegung Plan-Änderung 2 (Konsumenten-Beispiel) | geprüft, ohne Befund — vertretbar und vor dem Code festgehalten (`8a458a27` zeitlich vor `56c8f62f`); die Spezifikation trägt 159 geänderte Zeilen (158 Link-Ziele plus Zeile 920), `grep -c 'spezifikation\.md' spec/specification.md` meldet danach 2 (Z. 233, 234); Handbuch Z. 372 und 373 unberührt (`git diff -U0` am Handbuch nennt genau eine Zeile, 2792); die vier Spiegel-Stellen des Beispiels bleiben gleich, die `--manual`-Ausgabe des Beispiels ändert sich damit nicht |
| Produkt-Berührung | geprüft, ohne Befund — Quelltext geändert: `manual.go` (Pfad-Konstante, `go:embed`), `cli_manual_test.go` Z. 21, `diagnose_test.go` Z. 73 (`git diff --stat f44c6e65..HEAD -- '*.go' Dockerfile Makefile`: 3 Dateien, 4 Zeilen); `tools/image-test.sh` Z. 163; keine Logik, kein Standardwert, kein Grund-Code |
| `DC-FA-CLI-013` Wortlaut | geprüft, ohne Befund — Lastenheft-Diff 3 Zeilen (Version, Zeile 4762, Historie-Zeile), die Anforderung `DC-FA-CLI-013` unberührt; die Spezifikation §`DC-FA-CLI-013.a` ändert nur Link-Ziele und den Pfad in Schritt 2 |
| Kommentarzeilen, Zählung | geprüft, ohne Befund — `git grep -n 'spezifikation\.md' -- '*.go' \| grep -v 'go:embed' \| grep -c '^[^:]*:[0-9]*:[[:space:]]*//'`: 95 Zeilen in 38 Dateien; ohne `*_test.go` 90 in 33 (Plan-Änderung 1: Differenz 5 Zeilen in Testdateien); dazu `Dockerfile` 1, `Makefile` 1, `tools/` 2 (4 Zeilen); sechs Fixture-Zeilen in `cli_acceptance_test.go` (2) und `matrix_test.go` (4) und vier Konsumenten-Zeilen in `config_template.go` und `suggest.go`; alle unverändert (Verdikt Frage 7; die Folge für §3.7 ist dort als akzeptiertes Negativ benannt) |
| Lastenheft 0.115.1 → 0.115.2 | geprüft, ohne Befund — Version, genau eine Historie-Zeile oben in `§7` im Zeitstrahl der übrigen („Redaktionell: …; keine inhaltliche Änderung, kein `DC-*` ändert sich“, Verweis `—`, Datum 2026-11-26), Zeile 4762 (Messmethode zu `DC-QA-01`) nennt `spec/specification.md`; kein weiterer Satz, keine Anforderung geändert; kein Spec-Stratum nennt Slice oder Welle (`AGENTS.md` §3.4) |
| Abgrenzung | geprüft, ohne Befund — keine Datei zu `slice-296`/`slice-298`, `archive`, `wave`, `stub.go` oder `ADR-0117` im Diff; `CHANGELOG.md`, READMEs außerhalb der zwei Link-Ziele, Handbuch-Kopf unberührt; `docs/user/benutzerhandbuch.md` nur Z. 2792 |
| Frozen-Stabilität | geprüft, ohne Befund — in `done/`, `docs/reviews/`, `docs/plan/cr/`, `conventions/done/` und den ADRs ändert der Diff ausschließlich Link-Ziele (siehe Substitutionsprobe); `CHANGELOG.md` nicht berührt |
| Zwischenstände (§3.3-Zerlegung) | geprüft, ohne Befund — Plan und Botschaften benennen sie richtig: `doc-check` am Stand `56c8f62f` endet mit `d-check: error: konfiguriertes ids-Target existiert nicht: spec/spezifikation.md` (zusätzlich nicht baubar, `go:embed`), am Stand `340f3f15` mit `212 Befund(e)`, am Stand `061ebfc2` mit `0 Befund(e)`; nur die Commits 1 und 2 der Reihenfolge sind einzeln rot, die sechs Commits der Spanne liegen ungepusht vor `origin/main` (ein Push) |
| Botschaften (Regel 15) | geprüft, ohne Befund — die Zahlen der Botschaften (64 ADRs mit 119 Zeilen, 39 mit 80, 38 Dateien mit 26 lebenden Markdown-Dateien plus Handbuch, 158 plus Zeile 920) stimmen mit dem Diff; „Dieser Stand ist einzeln noch rot“ ist gemessen (212 Befunde); die Botschaft von `8a458a27` nennt 119 für den Stand ihrer Messung (siehe F-3) |
| Kommentare (§3.7) | geprüft, ohne Befund — der Kommentar am Ventil trägt Zusage und Grenze, keine Review-Historie, Herkunft ein Feld; die Kommentarzeile in `.d-check.yml` Z. 991 ist mit dem Pfad nachgezogen; neue Go-Kommentare gibt es nicht |
| Sensor-Dateien und lebende Doku | geprüft, ohne Befund — `git grep -n 'spezifikation\.md'` außerhalb der Frozen-Klassen, der Pläne und `*.go` nennt nur Ventil, `Dockerfile` 86, `Makefile` 35, `tools/bench-fixture.sh` 2, `tools/coverage-gate.sh` 3 (Kommentare, Verdikt Frage 7), Handbuch 372/373, Spezifikation 233/234, `MR-081` 182 (Gegenstand); `.claude/` und `.harness/skills/` nennen den Namen nicht |
| Gates, die ich gefahren habe | geprüft, ohne Befund — `make doc-check` (`1324 Datei(en) geprüft, 0 Befund(e)`, Exit 0), `make adr-check RANGE=f44c6e65..HEAD` (`1318 Datei(en) geprüft, 0 Befund(e)`, Exit 0), `make trace-check RANGE=f44c6e65..HEAD` (`1318 Datei(en) geprüft, 0 Befund(e)`, Exit 0), `make test` (Exit 0); nicht gefahren: `make gates`, `make image-test` (Sache der Verifikation) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 (F-1) |
| LOW | 2 (F-2, F-3) |
| INFO | 3 (F-4, F-5, F-6) |

## Verdikt

Der Diff trägt den Plan: der Move ist ein reiner `R100`, der Link-Nachzug ist an allen drei Stellen (ADRs, übrige eingefrorene Dateien, lebende Dateien) eine gemessene Substitution des Ziels hinter `](`, `adr-check` über die ganze Spanne endet mit Exit 0, das Ventil verhält sich so, wie sein Kommentar es sagt, das Lastenheft ändert genau Version, Zeile 4762 und eine Historie-Zeile. **F-1 blockiert** die Welle-Closure, nicht diesen Slice: Der Wellen-Plan muss die fünfte Klasse von (d3) und die geänderten Zahlen des Ventils führen, bevor die Gegenmessung läuft — Übergabe an den **Planner** (Plan-Änderung am Wellen-Plan, kein Code). F-2 geht an den **Implementer** (eine Zeile; die Entscheidung, ob die Aussage über das Template-Set stehen bleibt, ist seine mit Verweis auf `MR-070`), F-3 an den **Implementer** (Endzahl 120 in die Closure-Notiz). F-4 bis F-6 brauchen keine Nachbesserung. Eingearbeitet werden nach der Vorgabe der Welle nur HIGH und MEDIUM.
