# Review R2 — slice-270: Die RTM zeigt die Test-Nachweise

**Review-Art:** Code (Diff gegen Plan, ADRs, Konventionen und Hard Rules — nicht gegen die DoD)
**Gegenstand:** slice-270, Range `d7e4a028..HEAD` — `6f5cd3a1` (Einarbeitung R1 F-1 bis F-9);
dazu Stichprobe im abgeleiteten Bestand `docs/user/abdeckung-tests.md`
**Skill:** `.harness/skills/reviewer.md` @ 1.20.0
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Slice-Plan slice-270;
[ADR-0104](../plan/adr/0104-test-nachweise-entlasten-in-der-rtm.md) (Proposed);
[`MR-078`](../../harness/conventions.md#mr-078);
[`DC-FA-COV-001`](../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in),
[`DC-FA-CLI-012`](../../spec/lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben);
`AGENTS.md` §3.1, §3.2, §3.6, §3.7, §5 Regeln 13/15/16;
`harness/sensors/test.md` (Abschnitt *Abdeckungs-Dateien der RTM*).
Vorherige Findings am selben Modul: R1 desselben Slice (F-1 bis F-9).

## Messungen (Kommando und Ergebnis)

| # | Kommando | Ergebnis |
|---|---|---|
| M1 | `make abdeckung`, danach `git status --short`, `ls -l docs/user/abdeckung-*.md` | Exit 0; Arbeitsbaum sauber, Rechte `-rw-r--r--` — die committeten Dateien folgen ihrer Ableitung |
| M2 | Probe: `internal/adapter/driven/configyaml/probe_windows_test.go` und `_probe_test.go` mit je einem Test samt Kennung im Doc-Kommentar, `make abdeckung`, `grep TestProbe docs/user/abdeckung-tests.md`; danach beide Dateien entfernt, `git checkout` der Abdeckungs-Datei | Exit 0; **beide** Tests erscheinen als Zeile (`DC-FA-DIST-002`, `DC-FA-WF-001`) — `go test` baut keine der beiden Dateien unter Linux |
| M3 | Probe: Datei mit nur `// +build integration` und einem Test, dessen Rumpf ein undefiniertes Symbol ruft, `make abdeckung`; danach entfernt | Exit 0 — das Paket kompiliert, `go test` nimmt die Datei also aus; die Ableitung führt den Test trotzdem als Zeile (`DC-FA-WF-001`) |
| M4 | Zeilen mit Kennung `DC-FA-CLI-012` in `docs/user/abdeckung-tests.md` gezählt; Doc-Kommentar von `TestQA03_ClosureProfil_KeineZweiteNetzTuer`; Datei-Kommentar `cli_config_path_test.go` | genau **eine** Zeile, und sie ist dieser Repo-Selbsttest; die Produkt-Tests `TestConfigPath_*` nennen `DC-FA-CLI-012` nur im Datei-Kommentar (Zeile 13) und zählen nicht |
| M5 | Tests, die Repo-Dateien lesen (`repoRoot()`, `../../../`, `.d-check.yml`), gegen ihre Zeilen in `abdeckung-tests.md` | außer M4 nur Produkt-Tests mit Fixtures; `enforcement_layer_test.go`, `docexamples_test.go`, `handbook_examples_test.go`, `abdeckung_test.go` ohne Zeile |
| M6 | `grep -n '^\s*#' tools/image-test.sh` gegen `imageTestPhaseRE` | Prosa-Zeilen 5–14 (`#   (1) Happy: …`) tragen keine Striche und treffen nicht; nur die vier Kopfzeilen 66/88/98/108 |
| M7 | `find internal cmd -name '*_test.go'` mit GOOS/GOARCH-Suffix oder `_`/`.`-Präfix; `grep -rln '^// +build\|^//go:build'` | je 0 Treffer — die Lücken aus F-2 sind heute ohne Instanz |
| M8 | `make gates` | Exit 0: `baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; Coverage 94,70 % (Schwelle 93 %); doc-check `1102 Datei(en) geprüft, 0 Befund(e)` |
| M9 | `git diff d7e4a028..HEAD \| grep nolint`; neue Kommentare auf `F-<n>`, `R1`, `slice-`, Herkunfts-Prosa | je 0 Treffer |

## Einlösung R1

| R1 | Stand | Beleg |
|---|---|---|
| F-1 | **teilweise** — die zwei neuen Tests nennen `ADR-0104`, ihre Zeilen sind weg; dieselbe Klasse steht im Bestand weiter (F-1 unten) | Diff `abdeckung-tests.md`; M4 |
| F-2 | eingelöst — `Schärft:` ohne Spec-Stelle | Diff ADR-0104 |
| F-3 | eingelöst für die gemeldete Form (abweichende Strich-/Leerzeichen-Form ist jetzt rot); Rest-Formen siehe F-4 unten | Negativliste `TestImageTestZeilenVerlangtDenAnkerJederPhase` |
| F-4 | **teilweise** — Signatur, `//go:build`, Verzeichnis-Regel eingelöst; Datei-Name als Build-Bedingung und `// +build` nicht (F-2 unten) | M2, M3 |
| F-5 | eingelöst — `abdeckung` steht vor dem `tidy`-Kommentarblock | Diff Makefile |
| F-6 | eingelöst für `chmod` (Glob); der Listen-Abgleich hält weniger, als sein Kommentar sagt (F-3 unten) | Diff Makefile; `abdeckungsListeDecktCoverage` |
| F-7 | eingelöst — Negativtests nehmen `liveConfig(…).Trace.ReqPattern` | `abdeckung_ableitung_test.go` |
| F-8 | eingelöst — Plan §1/§2 sprechen von der Spalte `Coverage` mit zwei Labels | Diff Slice-Plan |
| F-9 | eingelöst — `nurTabellenKennungen` mit Negativtest | `TestNurTabellenKennungen` |

## Findings

### F-1 — MEDIUM — Ein Repo-Selbsttest ist der einzige Test-Nachweis von `DC-FA-CLI-012`

- **quelle:** `AGENTS.md` §5 Regel 16; Skill-Prüffrage 9; R1 F-1 (dieselbe Klasse); ADR-0104 Re-Evaluierungs-Trigger 2
- **pfad:** `docs/user/abdeckung-tests.md` · „`TestQA03_ClosureProfil_KeineZweiteNetzTuer`"; Quelle `internal/adapter/driven/configyaml/gate_consistency_test.go` · „DC-FA-CLI-012/ADR-0048). Es trägt bewusst NICHT den vollen Netzlos-Doku-Satz"
- **befund:** Der Test prüft, dass die `.d-check.closure.yml` **dieses Repos** kein Netz-Modul einschaltet; `--config` als Produkt-Fähigkeit durchläuft er nicht. Sein Doc-Kommentar nennt `DC-FA-CLI-012` als Kontext, und die Ableitung führt ihn als einzige Tests-Zeile dieser Anforderung (M4) — die Produkt-Tests `TestConfigPath_*` tragen die Kennung nur im Datei-Kommentar. Failure-Szenario: `make trace` zeigt `DC-FA-CLI-012` mit Label `Tests`, belegt allein durch einen Test, der sie nicht prüft; verliert die Anforderung ihren Slice-Bezug, bleibt sie waisenfrei. Das ist der zweite Re-Evaluierungs-Trigger von ADR-0104, heute schon in der Spalte eingetreten, nur noch nicht in der Waisen-Zählung. Die R1-Einarbeitung hat die Klasse an den zwei neuen Tests behoben, nicht im Bestand gesucht.
- **verifizierbar:** ja — `make trace`, Zeile `DC-FA-CLI-012`; `grep -n 'CLI-012' docs/user/abdeckung-tests.md`
- **klasse:** selbsttest-deklariert-produkt-anforderung (2. Auftreten im Vorgang)

### F-2 — MEDIUM — Die Ableitung zählt Testdateien, die `go test` nach Dateiname oder `// +build` nicht baut

- **quelle:** Skill-Prüffragen 13, 17, 20; `AGENTS.md` §5 Regel 13; R1 F-4
- **pfad:** `internal/adapter/driven/configyaml/abdeckung_test.go` · „func imStandardBuild(f *ast.File) bool" und „if !strings.HasSuffix(p, \"_test.go\")"; Zusage im Intro von `docs/user/abdeckung-tests.md` · „Gezählt wird, was `go test` als Test ausführt"
- **befund:** Ausgewertet wird nur die `//go:build`-Zeile. Dateien mit GOOS/GOARCH-Suffix (`x_windows_test.go`), Dateien mit `_`/`.`-Präfix und Dateien mit nur `// +build`-Zeile nimmt `go test` aus, die Ableitung zählt sie — gemessen in M2 und M3 mit `make abdeckung`, Exit 0. Failure-Szenario: ein `…_windows_test.go` mit `// TestX prüft DC-FA-DIST-002.` entlastet `DC-FA-DIST-002` in der RTM, ohne dass der Test unter `make test` je läuft; `make abdeckung` schreibt die Zeile, `make test` hält sie danach grün. Richtung ist falsche Entlastung, nicht Unterzählung. Heute ohne Instanz (M7); die Negativliste nennt keinen der drei Fälle, und die Commit-Botschaft fasst F-4 als „gezählt wird, was go test ausführt" zusammen.
- **verifizierbar:** ja — die Fälle in `TestGoTestZeilenFolgtDemBuildConstraint` / `TestGoTestZeilenUeberspringtVerzeichnisseWieGoTest` aufnehmen, sie werden rot; oder M2/M3 wiederholen
- **klasse:** testname-regel-ohne-implizite-build-bedingung

### F-3 — LOW — Der Listen-Abgleich sagt „genau" und „je Datei eine Quelle" und prüft beides nicht

- **quelle:** `AGENTS.md` §3.7 (ein Kommentar beschreibt, was da ist) und §5 Regel 13; ADR-0104 Entscheidung 1
- **pfad:** `internal/adapter/driven/configyaml/abdeckung_test.go` · „verlangt, dass trace.coverage genau die\n// Abdeckungs-Dateien dieser Liste einbindet, je Datei eine Quelle." und „strings.HasPrefix(p, \"docs/user/abdeckung-\")"
- **befund:** Die Rückrichtung meldet nur Pfade mit Präfix `docs/user/abdeckung-`; eine weitere Coverage-Datei unter anderem Namen passiert, und der Negativtest bestätigt das ausdrücklich (Fall 4, `spec/x.md`, kein Fehler). Ob jede Datei eine eigene Quelle mit dem Label aus ADR-0104 bildet, wird nicht geprüft — beide Dateien unter einer Quelle oder vertauschte Labels bleiben grün. Die Präfix-Heuristik trägt also für den Fall, den F-6 meinte (eine weitere **abgeleitete** Datei), nicht für eine von Hand gepflegte, die ADR-0104 als Alternative verwirft.
- **verifizierbar:** ja — `TestAbdeckungsListeDecktCoverage` Fall 4 zeigt es
- **klasse:** kommentar-zusage-weiter-als-pruefung

### F-4 — LOW — Phasen-Kopfzeilen außerhalb der erkannten Form bleiben still

- **quelle:** `AGENTS.md` §5 Regel 13; `harness/sensors/test.md` · „Eine Phase ohne Anker ist rot."
- **pfad:** `abdeckung_test.go` · „imageTestPhaseRE = regexp.MustCompile(`^#\s*-+\s*(\(\d+\).*?)[\s-]*$`)"
- **befund:** Das Muster verlangt `#` in Spalte 1, Striche und eine rein numerische Klammer; eine eingerückte Kopfzeile, `(4b)` oder eine Kopfzeile ohne Striche ist keine Phase, und fehlt der Anker, bleibt `make test` grün. Richtung Unterzählung; die Definition „Phase" steht nur im Code-Kommentar, `harness/sensors/test.md` und das Intro von `abdeckung-e2e.md` nennen sie nicht. Falsch-positive Treffer in Prosa wären laut (Fehler), nicht still (M6).
- **verifizierbar:** ja — Kopfzeile `# --- (4b) … ---` ohne Anker in der Negativliste
- **klasse:** erkennung-still-bei-formabweichung (Rest aus R1 F-3)

### F-5 — INFO — Das Ableitungs-Ergebnis hängt an der Laufumgebung

- **quelle:** undokumentierte Annahme
- **pfad:** `abdeckung_test.go` · „return tag == runtime.GOOS || tag == runtime.GOARCH || tag == \"gc\" || strings.HasPrefix(tag, \"go1.\")"
- **befund:** `make abdeckung` und `make test` laufen auf der Host-Architektur; eine Datei mit `//go:build amd64` ergäbe auf einem arm64-Host eine andere Datei als in CI. `unix` und andere implizite Tags gelten als falsch (Unterzählung), jedes `go1.N` als wahr. Die Fehlermeldung von `abdeckungsListeDecktCoverage` iteriert über Maps und nennt bei mehreren Abweichungen eine beliebige zuerst. Heute ohne Instanz (M7).
- **klasse:** ableitung-umgebungsabhaengig

### F-6 — INFO — Kennung im Datei-Kommentar zählt nicht

- **quelle:** ADR-0104 Entscheidung 2 (gewollte Regel)
- **pfad:** `internal/adapter/driving/cli/cli_config_path_test.go` · „// --config (DC-FA-CLI-012) end-to-end gegen ein echtes Temp-Repo"
- **befund:** Zwölf Produkt-Tests für `--config` tragen die Kennung nur einmal über der Datei und erscheinen nicht in der RTM; nach der Regel richtig, aber genau deshalb kippt die Spalte für `DC-FA-CLI-012` auf den falschen Test (F-1). Wie viele Anforderungen auf diese Weise unterzählt sind, ist nicht gemessen.
- **klasse:** deklaration-auf-dateiebene

## Negativbefunde

- **Signatur-Prüfung:** geprüft, ohne Befund — genau ein Feld, höchstens ein Name, `*<testing-Name>.T`; alle Import-Namen von `testing` (auch Alias, Mehrfach-Import); `TestMain(*testing.M)`, `*testing.B`, parameterlos, fremdes `*T` fallen heraus; Negativliste deckt alle. Abweichung vom Go-Werkzeug nur beim Punkt-Import (`*T`, Unterzählung), ohne Instanz.
- **`//go:build`-Auswertung:** geprüft, ohne Befund — alle Kommentargruppen vor der `package`-Klausel, Lizenzkommentar davor stört nicht; ungültiger Ausdruck schließt die Datei aus, baut dann aber auch unter `go test` nicht (laut, `make test` rot); mehrere `//go:build`-Zeilen: das Go-Werkzeug verweigert die Datei, die Ableitung nimmt die erste — ohne Folge, weil der Build bricht. `// +build` siehe F-2.
- **Verzeichnis-Regel:** geprüft, ohne Befund — `testdata`, `_x`, `.x` übersprungen, Startverzeichnis ausgenommen; kein Untermodul (R1 M7).
- **Kopfzeile in Prosa:** geprüft, ohne Befund — kein `# --- (` in Prosa von `tools/image-test.sh` (M6); ein solcher Treffer wäre fail-closed.
- **`nurTabellenKennungen`:** geprüft, ohne Befund — prüft die gerenderte Datei samt Titel, Intro und Phasentitel gegen die Zeilen-Kennungen; Negativtest vorhanden.
- **Mutationsbeleg der Botschaft:** geprüft auf Plausibilität, ohne Befund — zu jeder der sechs genannten Prüfungen steht ein Negativfall, der bei entfernter Prüfung rot würde; den Alias-Fehler deckt `TestAlias` ab. Nicht nachgefahren.
- **Grenz-Texte:** geprüft, bis auf F-2/F-3/F-4 ohne Befund — ADR-0104 und MR-078 nennen „Deklaration, kein Beleg" und sagen nichts, was der Code nicht trägt.
- **§3.2:** geprüft, ohne Befund (M9).
- **§3.7:** geprüft, bis auf F-3 ohne Befund — Zusage-, Kopplungs- und Grenz-Kommentare, keine Befund-Marker, keine Slice-Nummern (M9).
- **§3.1:** geprüft, ohne Befund — Erzeugung und Prüfung in Docker; die Proben M2/M3 liefen über `make abdeckung`.
- **gofmt:** nach Augenschein ohne Befund (Tabs, Feld-Ausrichtung); kein Werkzeug gelaufen — das Lint-Profil führt keinen Formatter, `make` bietet kein Target dafür.
- **Determinismus:** geprüft, ohne Befund für den grünen Fall — `make abdeckung` zweimal byte-gleich zum committeten Stand (M1, nach M2/M3); Sortierung nach Datei, Test; `WalkDir` lexikalisch. Rot-Fall siehe F-5.
- **`make gates`:** grün (M8).

## Kategorie-Summary

HIGH 0 · MEDIUM 2 (F-1, F-2) · LOW 2 (F-3, F-4) · INFO 2 (F-5, F-6).
Wiederkehrende Klasse: Deklaration über ihren Geltungsbereich — F-1 ist das zweite Auftreten im Vorgang (R1 F-1), diesmal im Bestand, den die Einarbeitung nicht durchsucht hat; F-2 ist der Rest von R1 F-4 in der Richtung falscher Entlastung.

## Verdikt

Nicht freigegeben. F-1 ist eine eingetretene falsche Zuordnung in der RTM (`DC-FA-CLI-012`),
F-2 eine gemessene Lücke in der Richtung, die ADR-0104 als Re-Evaluierungs-Trigger führt.
LOW/INFO: Annahme oder Begründung genügt.
