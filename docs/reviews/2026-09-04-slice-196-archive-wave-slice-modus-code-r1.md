# Review-Report — slice-196 (archive-wave wellenloser Einzel-Slice-Modus)

- **Review-Art:** Code (gegen Plan, `AGENTS.md` Hard Rules, Baseline-Regelwerk `modul-06-roadmap.md`)
- **Gegenstand:** Commit `50e243c` (HEAD) — `tools/archive-wave/{main.go,collect.go,archive.go,stub.go,rewrite.go,slice_mode_test.go,main_test.go}`, `tools/archive-wave/Makefile`, root `Makefile`, `AGENTS.md`
- **Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md)
- **Modell-ID:** claude-sonnet-5
- **Datum:** 2026-09-04
- **Eingangs-Kontext:** Slice-Plan [`slice-196`](../plan/planning/done/slice-196-archive-wave-slice-modus.md); [`AGENTS.md`](../../AGENTS.md) §3.1, §3.3, §3.4, §3.7, §3.8, §4 (`make archive-wave`-Zeile); Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht. Kein Zugriff auf Implementer-Kontext (Rollentrennung, Modul 8).

Repo unverändert nach dem Review (`git status --short` leer, HEAD `50e243c`).

---

## Eigener Lauf (Ausgabe, nicht behauptet)

| Lauf | Ausgabe |
|---|---|
| `make archive-wave-test` | Docker-Build `test`-Target erfolgreich (`RUN CGO_ENABLED=0 go test ./...` als Layer, cache-hit auf unverändertem Quellstand — ein Build-Fehlschlag hätte den `docker build` nicht erfolgreich terminieren lassen) |
| `make gates` | `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; semgrep `Ran 55 rules on 62 files: 0 findings`; `d-check: 653 Datei(en) geprüft, 0 Befund(e)` (×2, `targets`/`planning`-Läufe); Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` — deckt sich mit der Commit-Botschaft (zehn Gates) |
| `make fullbuild` | `image-test: OK`; `[ci] gates + image-test green`; `bench: OK` (Median 1177 ms); RTM `51 Anforderung(en), 0 Waise(n)`; Closure-Profil `d-check: 584 Datei(en) geprüft, 0 Befund(e)`; Schlusszeile `[fullbuild] green — image-hash sha256:60c3fc51781e4769a42cc92502336e1e39bbfa6b09b5adf3e12b58dfa4529815` |
| `docker run … archive-wave:latest -root=/repo -slice=slice-137` (Dry-Run gegen echten Bestand) | `archive-wave: slice-137 (wellenlos)` · 1 Review-Report gefunden · `Verweise auf die geloeschten Review-Reports: (keine)`; `git status --short` danach leer — Dry-Run schreibt tatsächlich nichts |
| `docker run … archive-wave:latest -root=/repo -slice=slice-195` | `archive-wave: slice-195 nicht gefunden` (Exit 1) — erwartet, slice-195 liegt bereits im welle-88-Archiv, kein echter Wellen-Zugehörigkeits-Fall mehr im Bestand verfügbar zum Gegenprobieren |
| Grep über `docs/plan/planning/done/*.md` | Kein flacher `done/`-Slice führt aktuell ein `**Welle:**`-Feld mit echtem `welle-\d+`-Treffer (alle welle-gebundenen Slices sind bereits archiviert) — die Ablehnungs-Logik ist damit gegen den Bestand nicht smoke-testbar, nur fixture-testbar (unvermeidbare Bestandslücke, kein Implementierungsfehler) |

---

## Findings

### F-1 · HIGH · `AGENTS.md` §3.7 (fünf Kommentar-Klassen; Herkunfts-Prosa/Slice-Nummer explizit verboten) · `tools/archive-wave/main.go:47`, `tools/archive-wave/stub.go:123,130`, `tools/archive-wave/slice_mode_test.go:12,115`

**Befund:** Fünf **neue** Go-Kommentare in diesem Diff zitieren `slice-196` als Herkunfts-Prosa statt als das einzig zulässige auflösbare Feld (`DC-*`, `ADR-*`, `MR-*`, `seit welle-<NN>`):

- `main.go:47` (`validateModeFlags`): „der alte Pflicht-Fehler von vor **slice-196**, jetzt fuer zwei Flags statt einem" — erzählt eine Vorher/Nachher-Chronik des Flags, keine der fünf Klassen (Zusage · Kopplung · Abgrenzung · Rang-Zeiger · Grenze).
- `stub.go:123` (`SliceStubStandalone`): „Einzel-Slice-Modus (**slice-196**)" — benennt den Ursprungs-Slice der Funktion, kein Kopplungs-/Abgrenzungs-Inhalt.
- `stub.go:130`: „keine neue Norm, nur ein neuer Feldwert (**slice-196** §2 Punkt 4)" — zitiert die Nummer eines Absatzes im Slice-*Plan* als Begründung; das ist Herkunfts-Prosa, keine auflösbare Baseline-Schema-Kennung.
- `slice_mode_test.go:12`: „das im Slice-Plan (**slice-196** §2 Punkt 5) verlangte Fixture" — dieselbe Form, in einem Test-Kommentar (§3.7 gilt „für Code, Konfiguration und Skripte", Tests eingeschlossen).
- `slice_mode_test.go:115`: „belegt **slice-196** §4: ein Slice mit …" — dieselbe Form.

Zu unterscheiden von den **pre-existierenden** (grandfathered) Kommentaren in denselben Dateien wie „gemessen an welle-70: slice-101s eigener Welle-Link brach genau so" (`rewrite.go:86`) oder „gemessen an slice-075" (`stub.go:71`): Diese zitieren eine **empirische Beobachtung am Bestand** (Rang-Zeiger-Klasse, Beleg für eine Grenzfall-Behandlung) und sind nicht Teil dieses Diffs. Die fünf oben genannten Stellen sind dagegen **neu** und zitieren den **eigenen** Slice als Ursprung einer Design-Entscheidung bzw. als Vorher/Nachher-Erzählung — genau die Form, die §3.7 ausdrücklich als „Herkunfts-Prosa" und „Slice-Nummern" nennt und die laut [`reviewer.md`](../../.harness/skills/reviewer.md) Zeile 36 (Frage 6) HIGH ist. §3.7s Bestandsgrenze gilt nur für **vor** der Schärfung geschriebene Zeilen; „Neuzugänge fallen überall unter den Anker."

**Verifizierbar:** ja — `grep -rn "slice-196" tools/archive-wave/*.go` zeigt genau die fünf Stellen; keine davon existierte vor Commit `50e243c` (neue Dateien/neue Funktionen).
**Klasse:** `kommentar-herkunfts-prosa-slice-nummer`

### F-2 · MEDIUM · Testabdeckung (Modul 8 „Reviewer prüft Findings/Testqualität") · `tools/archive-wave/collect.go:121-147` (`FindSlice`)

**Befund:** `FindSlice`s eigener Docstring beansprucht ausdrücklich: „keine oder mehrdeutig sind beides Fehler — anders als `FindWellePlan` gibt es keinen legitimen Nullfall." Die strukturell analoge, bereits bestehende Funktion `FindWellePlan` hat für genau diese beiden Fehlerzweige eigene Tests (`TestFindWellePlan_Mehrdeutig`, `TestFindWellePlan_Keine` in `collect_test.go`). Für `FindSlice` existiert **kein** entsprechender direkter Test: `slice_mode_test.go` ruft `FindSlice` nur einmal im Erfolgsfall auf (Zeile 134, `TestRunSlice_DanglingReviewReference`). Eine Mutation, die z. B. bei 0 Treffern statt eines Fehlers einen leeren String zurückgäbe, oder bei mehreren Treffern den ersten statt eines Fehlers zurückgäbe, würde von keinem Test erkannt — `TestRunSlice_RejectsWelleSlice` und `TestRunSlice_Apply` legen beide Fixtures mit genau einer passenden Datei an.

**Verifizierbar:** ja — `grep -n "FindSlice" tools/archive-wave/*_test.go` zeigt einen einzigen Aufruf (Happy Path); ein testweise eingefügter Bug in den `case 0`/`default`-Zweigen von `FindSlice` lässt `make archive-wave-test` grün.
**Klasse:** `fehlerzweig-ohne-direkten-test`

### F-3 · MEDIUM · Testabdeckung · `tools/archive-wave/rewrite.go:185-230` (`FindReferencesToPaths`), aufgerufen aus `main.go:185`

**Befund:** `runSlice` ruft `FindReferencesToPaths(root, reviews, append([]string{slicePath}, reviews...))` auf — die Slice-Datei selbst **und** die zu löschenden Review-Reports werden von der Verweis-Suche ausgeschlossen (`excludeAbs`). Das ist eine bewusste Design-Entscheidung (ein Verweis aus der Slice-Datei auf ihren eigenen, gleich gelöschten Review-Report wäre kein „toter" Verweis im eigentlichen Sinn, weil der Slice-Inhalt selbst durch den Stub ersetzt wird), aber sie ist **nicht** durch die Fixture belegt: Weder `docs/plan/planning/done/slice-601-eins.md` noch die beiden Review-Reports in `buildSliceFixture` enthalten einen Markdown-Link aufeinander. Eine Regression, die `slicePath` (oder die Reviews) versehentlich aus `excludeAbs` entfernt, würde `TestRunSlice_DanglingReviewReference` nicht rot machen, weil das Fixture keinen solchen Link trägt, der dann fälschlich mitgezählt würde. Stichprobe im Bestand (`grep -rn "\](\.\{0,3\}/*reviews/[^)]*\.md)" docs/plan/planning/done/*.md`) findet aktuell keine echten Markdown-Links von Slice-Körpern auf ihre eigenen Review-Reports — das Praxisrisiko ist damit gering, aber der Codepfad selbst bleibt unbelegt.

**Verifizierbar:** ja — ein `excludeAbs = reviews` statt `excludeAbs = append([]string{slicePath}, reviews...)` in `main.go:185` lässt `make archive-wave-test` weiterhin grün, weil kein Fixture-Fall diesen Unterschied misst.
**Klasse:** `exclude-branch-ohne-fixture-beleg`

---

## Negativbefunde (geprüft, ohne Befund)

- **`ApplySlice`-Reihenfolge (archive.go:133-164).** `buildZip` liest die Slice-Datei **vor** dem `os.WriteFile`, das sie durch den Stub ersetzt — das Archiv enthält nachweislich den Originaltext, nicht den Stub. Bestätigt durch `TestRunSlice_Apply`: das Archiv trägt `docs/plan/planning/done/slice-601-eins.md` mit vollem Titel-/Feld-Auszug, während der Stub-Text separat geprüft wird.
- **Kein Move für die Slice-Datei selbst.** `ApplySlice` erzeugt keinen `Move`-Eintrag und `runSlice` ruft `RewriteRepo` nicht auf — korrekt, da sich der Pfad der Slice-Datei nicht ändert; ein repo-weiter Verweis-Nachzug ist für sie unnötig. Kein Move-Objekt existiert, das versehentlich falsch angewendet werden könnte.
- **`welleIDInFieldRE`-Ablehnung.** Fixture (`TestRunSlice_RejectsWelleSlice`) und eigener Lauf gegen `slice-137` (echtes `**Welle:** — wellenlos.`-Feld ohne `welle-\d+`-Teilstring) bestätigen beide Zweige der Unterscheidung.
- **`FindSlice`-Präfix-Kollisionsschutz.** Der Präfix `sliceID + "-"` (mit angehängtem Bindestrich) verhindert, dass z. B. `slice-1370-x.md` fälschlich auf `-slice=slice-137` matcht — verifiziert durch manuelle Zeichenkettenanalyse, keine Ziffernfolgen-Verwechslung möglich.
- **`resolveLink`/`relativize` bei beliebiger Verzeichnistiefe.** Die Funktion ist unverändert aus dem bestehenden Wellen-Modus übernommen und bereits eigenständig getestet (`TestResolveLink` in `rewrite_test.go`, mehrere Tiefen). Der neue Fixture-Test (`TestRunSlice_DanglingReviewReference`) prüft zusätzlich eine vier Ebenen tiefe Quelle (`docs/plan/planning/in-progress/roadmap.md` → `../../../reviews/…`) — Auflösung von Hand nachgerechnet, korrekt.
- **`run`→`runWelle`-Umbenennung.** Diff zeigt ausschließlich die Signaturzeile geändert, der Funktionskörper ist byteidentisch — keine Verhaltensänderung am bestehenden Wellen-Modus.
- **Geteilte Helfer (`buildZip`, `RelPath`, `SliceIDFromPath`, `ReadWelleField`, `CollectReviews`, `mdLinkRE`/`resolveLink`).** Keine dieser Funktionen wurde im Diff verändert — beide Modi rufen exakt denselben, unveränderten Code auf. Kein Divergenz-Risiko.
- **Makefile-Logik (`tools/archive-wave/Makefile:run`, root-`Makefile:archive-wave`).** Beide-leer und beide-gesetzt sind zwei getrennte `@if`-Blöcke mit je eigenem Exit-2 und eigener Fehlermeldung; der `docker run`-Aufruf wählt `$(if $(WELLE),-welle=$(WELLE),-slice=$(SLICE))` — da die vorangehenden Guards Exklusivität bereits erzwingen, kann dieser Zweig nie mit falschem Flag aufgerufen werden. Von Hand durchgespielt: leer/leer → Fehler 1; WELLE/leer → `-welle=`; leer/SLICE → `-slice=`; WELLE/SLICE → Fehler 2.
- **`TestValidateModeFlags`.** Alle vier Kombinationen der Wahrheitstabelle (beide leer, nur welle, nur slice, beide gesetzt) sind abgedeckt, korrekt erwartet.
- **§3.1 (Docker/make-only).** Alle eigenen Läufe (`make archive-wave-test`, `make gates`, `make fullbuild`, `make -C tools/archive-wave build`, `docker run archive-wave:latest …`) liefen über Make-Targets bzw. das gebaute Image, kein Host-Go.
- **§3.4 (Architektur/Spec referenzieren nie abwärts).** `tools/archive-wave` ist kein d-check-Spec-Stratum; nicht einschlägig. `AGENTS.md`s Gate-Zeilen-Update zitiert Slice-Nummern wie zuvor schon üblich (AGENTS.md ist Doku, nicht durch §3.7 „Code/Konfiguration/Skripte" gebunden) — kein Finding.
- **AGENTS.md-`make archive-wave`-Zeile vs. Code.** Zeile beschreibt zwei mutually-exclusive Modi mit `APPLY=1`, `SLICE=<id>` ohne Move, flaches Archiv, Meldung (nicht Behebung) toter Verweise — deckt sich mit dem tatsächlichen Verhalten von `runSlice`/`ApplySlice`/`FindReferencesToPaths`.
- **DoD/Closure-Notiz-Ehrlichkeit.** Alle sieben abgehakten DoD-Punkte sind durch eigene Läufe bestätigt (Tests grün, Gates grün, Fullbuild grün, Ablehnungs-Logik funktioniert, Dry-Run-Smoke-Test gegen slice-137 reproduziert). Die drei unabgehakten Punkte (Review, Verifikation, Closure-Notiz-„geschrieben") sind konsistent unchecked, obwohl §9 bereits vollständigen Text trägt — keine Irreführung, nur eine noch ausstehende Schluss-Bestätigung. Die behauptete „selbst gefundene und korrigierte" Fixture-Bug-Episode lässt sich nicht mehr aus dem (bereits verdichteten) Commit rekonstruieren, ist aber plausibel: der finale Test (`TestRunSlice_DanglingReviewReference`) ist korrekt und grün, die Pfadtiefe (`../../../`) stimmt nachweislich mit der Fixture-Struktur überein.
- **WIP-Limit.** `docs/plan/planning/in-progress/` enthält nur `slice-196` (plus `roadmap.md`); `slice-197` liegt korrekt in `open/`, nicht parallel in Arbeit.

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 1 | F-1 |
| MEDIUM | 2 | F-2, F-3 |
| LOW | 0 | — |
| INFO | 0 | — |

**Finding-Klassen dieses Laufs:** `kommentar-herkunfts-prosa-slice-nummer` · `fehlerzweig-ohne-direkten-test` · `exclude-branch-ohne-fixture-beleg`

## Verdikt

**Merge-blockierend: ja (F-1).**

Die funktionale Umsetzung ist solide: `FindSlice`, `ApplySlice`, `runSlice` und `FindReferencesToPaths` verhalten sich korrekt gegen Fixtures **und** gegen einen echten Bestands-Slice (`slice-137`, per Dry-Run selbst nachgefahren); die Zip-vor-Stub-Reihenfolge stimmt, kein Move für die Slice-Datei selbst, die Ablehnung wellen-gebundener Slices funktioniert, beide Makefiles setzen die Mutual-Exclusivity korrekt um, `run`→`runWelle` ist eine reine Umbenennung ohne Verhaltensänderung, und `make archive-wave-test`/`make gates`/`make fullbuild` sind unabhängig grün nachgefahren (zehn Gates, Fullbuild-Image-Hash `sha256:60c3fc51…`).

**F-1 ist der blockierende Befund:** fünf neue Kommentare zitieren `slice-196` als Herkunfts-Prosa bzw. verweisen auf nummerierte Absätze des Slice-Plans — genau die Form, die `AGENTS.md` §3.7 ausdrücklich als verbotenen Inhalt nennt („keine Slice-Nummern", „keine Herkunfts-Prosa") und die laut [`reviewer.md`](../../.harness/skills/reviewer.md) HIGH ist. Kein Gate prüft das (§3.7: „Kein Gate prüft das … Der Reviewer-Skill trägt den HIGH-Anker dazu"), die Prüfung ist Aufgabe dieses Reviews. Empfehlung: die fünf Stellen umformulieren, sodass sie ihre jeweilige Klasse (Kopplung/Abgrenzung/Grenze) direkt tragen, ohne die eigene Slice-Nummer oder Plan-Absatznummer zu zitieren — die Herkunft „slice-196" gehört in Commit-Botschaft/Closure-Notiz, nicht in den Code-Kommentar.

F-2 und F-3 sind Testabdeckungs-Lücken an echten, aber im Bestand noch unbelegten Fehlerzweigen (kein aktueller Korrektheits-Fehler, aber eine Regression in diesen Zweigen bliebe unentdeckt) — vor Merge nachrüstbar, aber für sich genommen nicht merge-blockierend, wenn F-1 behoben wird.
