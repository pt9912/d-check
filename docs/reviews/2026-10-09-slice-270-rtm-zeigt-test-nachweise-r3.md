# Review R3 — slice-270: Die RTM zeigt die Test-Nachweise

**Review-Art:** Code (Diff gegen Plan, ADRs, Konventionen und Hard Rules — nicht gegen die DoD)
**Gegenstand:** slice-270, Range `cd5a19be..HEAD` — `641650e7` (Einarbeitung R2 F-1 bis F-5);
dazu Stichprobe im abgeleiteten Bestand `docs/user/abdeckung-tests.md`
**Skill:** `.harness/skills/reviewer.md` @ 1.20.0
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Slice-Plan slice-270;
[ADR-0104](../plan/adr/0104-test-nachweise-entlasten-in-der-rtm.md) (Proposed);
[`MR-078`](../../harness/conventions.md#mr-078);
[`DC-FA-COV-001`](../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in);
`AGENTS.md` §3.1, §3.2, §3.7, §5 Regel 13; `harness/sensors/test.md`
(Abschnitt *Abdeckungs-Dateien der RTM*); `Dockerfile` (Stages `test`, `abdeckung-gen`).
Vorherige Findings am selben Modul: R1 (F-1 bis F-9) und R2 (F-1 bis F-6) desselben Slice.

## Messungen (Kommando und Ergebnis)

| # | Kommando | Ergebnis |
|---|---|---|
| M1 | `make abdeckung` zweimal, `sha256sum docs/user/abdeckung-*.md`, `git status --short` | Exit 0 beide Male; Hashes identisch (`d4d5a4aa…`, `e571c62e…`), Arbeitsbaum sauber, Rechte `-rw-r--r--` |
| M2 | Probe: drei Testdateien im Paket `configyaml` — `//go:build unix` (Paket `configyaml_test`), `//go:build cgo` mit undefiniertem Symbol im Rumpf, `//go:build amd64.v1` (Paket `configyaml`, internes Testpaket) — `make abdeckung`, `grep TestProbe`; danach entfernt, `git checkout` der Abdeckungs-Datei | Exit 0; `TestProbeUnix` und `TestProbeV1` erscheinen, `TestProbeCgo` nicht — und das Paket baute, also nahm auch `go test` (CGO_ENABLED=0) die cgo-Datei aus. Internes und externes Testpaket zählen beide |
| M3 | `Dockerfile` Zeilen 38, 67–80 | `deps` läuft `FROM --platform=$BUILDPLATFORM`; `test` und `abdeckung-gen` erben davon, beide mit `CGO_ENABLED=0`, ohne `-tags`, ohne `GOARCH`/`GOAMD64`/`GOEXPERIMENT` |
| M4 | Testdateien unter `internal/`/`cmd/`, die Repo-Dateien lesen (`repoRoot`, `liveConfig`, `"../../../"`, `runtime.Caller`, `os.Getwd`), gegen ihre Zeilen in `abdeckung-tests.md`; dazu alle Dateien der Tabelle auf `os.ReadFile("`/`"../`-Pfade | sieben Selbsttest-Dateien, **keine** mit Zeile; die Treffer in Tabellen-Dateien sind ausnahmslos Fixture-Pfade (Escape-Proben, Temp-Repos) |
| M5 | `make trace`, Spalte Coverage gezählt; Zeile `DC-FA-CLI-012` | `Tests` 43 · `E2E, Tests` 5 · `—` 6; `DC-FA-CLI-012` ohne Coverage, Slices `slice-093`, `slice-180`, `ok`; `54 Anforderung(en), 0 Waise(n)` — die Zahlen der Commit-Botschaft stimmen |
| M6 | `grep -n -E '^\s*#.*(---\|\([0-9])' tools/image-test.sh` | Prosa-Zeilen 5–11 ohne Striche (keine Phase), vier Kopfzeilen 66/88/98/108 |
| M7 | Alter Testname und `imStandardBuild` repo-weit außerhalb `docs/reviews/` | 0 Treffer |
| M8 | `git diff cd5a19be..HEAD` auf `nolint`, Befund-Marker, Slice-Nummern in Kommentaren | 0 echte Treffer |
| M9 | `make gates` | Exit 0: `baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; Coverage 94,70 % (Schwelle 93 %); doc-check `1103 Datei(en) geprüft, 0 Befund(e)` |

## Einlösung R2

| R2 | Stand | Beleg |
|---|---|---|
| F-1 | eingelöst — der Selbsttest nennt nur `ADR-0048`, seine Zeile ist weg; im Bestand keine weitere Instanz | Diff `gate_consistency_test.go`; M4, M5 |
| F-2 | eingelöst — `go/build` (`Context.MatchFile`) entscheidet; Dateiname, `_`/`.`-Präfix, `// +build` in der Negativliste | `TestGoTestZeilenFolgtDemBuildDesGoWerkzeugs`; M2 |
| F-3 | eingelöst — Kommentar sagt, was geprüft wird, und nennt die Grenze; Meldung sortiert | `abdeckungsListeDecktCoverage`, `sortierteSchluessel` |
| F-4 | eingelöst für die gemeldeten Formen (Einzug, `(4b)`); die Grenz-Zeile dazu ist enger als die Regel (F-2 unten) | `imageTestPhaseRE`; Negativliste |
| F-5 | eingelöst für GOOS/GOARCH und Map-Reihenfolge; zwei Ränder bleiben (F-1 unten) | `gebautVonGoTest`; Fall `p/x_arm64_test.go` |
| F-6 | angenommen — das Intro nennt jetzt „oder im Datei-Kommentar zählt nicht" | Diff Intro |

## Findings

### F-1 — LOW — „wie `make test` läuft" und „hängt nicht vom Rechner ab" treffen zwei Ränder nicht

- **quelle:** `AGENTS.md` §5 Regel 13 (Grenze gegen den Gegenstand prüfen), §3.7 (ein Kommentar beschreibt, was da ist); R2 F-5
- **pfad:** `internal/adapter/driven/configyaml/abdeckung_test.go` · „Der Kontext ist fest linux/amd64 ohne\n// cgo und ohne `-tags`, wie `make test` läuft; das Ergebnis hängt so nicht vom\n// Rechner ab."; Gegenstand `Dockerfile` · „FROM --platform=$BUILDPLATFORM golang"
- **befund:** `make test` läuft auf der Build-Plattform (M3), auf einem arm64-Host also unter linux/arm64 — „wie `make test` läuft" gilt für amd64-Hosts und CI, nicht allgemein; das Intro der Abdeckungs-Datei sagt es richtig („unter linux/amd64"). Und `ctx := build.Default` übernimmt `ToolTags`, die `go/build` aus der Architektur der laufenden Toolchain bildet (`amd64.v1` auf diesem Host, M2), nicht aus dem gesetzten `GOARCH` — eine Datei mit `//go:build amd64.v1` oder `arm64.v8.0` ergäbe auf arm64 und amd64 verschiedene Abdeckungs-Dateien, `make test` wäre auf einem der beiden rot. Heute ohne Instanz; Richtung laut, nicht still.
- **verifizierbar:** nein auf diesem Host (x86_64) — die arm64-Hälfte braucht einen arm64-Lauf; die amd64-Hälfte zeigt M2 (`TestProbeV1` gezählt)
- **klasse:** ableitung-umgebungsabhaengig (Rest aus R2 F-5)

### F-2 — LOW — Die Grenz-Zeile der Phasen-Erkennung nennt nur die Striche

- **quelle:** `AGENTS.md` §5 Regel 13; Skill-Prüffrage 18; R2 F-4
- **pfad:** `abdeckung_test.go` · „GRENZE: Eine Phase, deren Kopfzeile keine drei Striche trägt, ist keine\n// Phase — ohne Anker fällt sie still aus der Ableitung."; Gegenstand · „(\(\d+[a-z]?\).*?)"; `harness/sensors/test.md` · „Eine Phase ohne Anker ist rot."
- **befund:** Still fällt auch jede Kopfzeile mit Strichen, deren Klammer nicht `Ziffern[ein Kleinbuchstabe]` ist — `# --- (4.1) …`, `# --- (4B) …`, `# --- Phase 5 ---`; die Grenz-Zeile nennt davon nichts. `harness/sensors/test.md` sagt weiter unbedingt „Eine Phase ohne Anker ist rot." ohne die Definition, die jetzt im Intro der E2E-Datei steht. Richtung Unterzählung, heute ohne Instanz (M6).
- **verifizierbar:** ja — `# --- (4.1) x ---\ncode\n` in der Negativliste bleibt ohne Fehler
- **klasse:** grenz-zeile-enger-als-regel

## Negativbefunde

- **`MatchFile` mit `OpenFile` in-memory:** geprüft, ohne Befund — `OpenFile` liefert für jeden Pfad den Inhalt der einen Datei, `MatchFile` öffnet genau eine; Dateiname-Regeln (GOOS/GOARCH-Suffix, `_`/`.`-Präfix) prüft `go/build` vor dem Öffnen, Header-Regeln (`//go:build`, `// +build`, nur vor der `package`-Klausel) danach — derselbe Code, den `go test` benutzt. Fehlerhafte oder doppelte Constraint-Zeilen machen die Ableitung laut, wie den Build.
- **Implizite Tags:** geprüft, ohne Befund — `unix` wahr, `cgo` falsch, passend zu `CGO_ENABLED=0` (M2); `gc` über `Compiler`, `go1.N` über `ReleaseTags` derselben Toolchain wie `make test`; `goexperiment.*` aus derselben Umgebung (M3). Rest siehe F-1.
- **Internes/externes Testpaket:** geprüft, ohne Befund — `MatchFile` fragt nicht nach dem Paketnamen, `go test` baut beide (M2).
- **Testdateien außerhalb `internal/`/`cmd/`:** geprüft, ohne Befund — nur `tools/archive-wave/` (eigenes `go.mod`, nicht in `go test ./...` des Hauptmoduls).
- **Selbsttests mit Produkt-Kennung:** geprüft, ohne Befund (M4) — keine Datei, die Repo-Dateien liest, hat eine Zeile; die Tabellen-Dateien lesen nur Fixtures.
- **Listen-Abgleich:** geprüft, ohne Befund — Kommentar, GRENZE und Code decken sich; Ausgabe über sortierte Schlüssel.
- **Phasen-Regex gegen den Bestand:** geprüft, ohne Befund — trifft genau die vier Kopfzeilen, Prosa-Aufzählung bleibt draußen (M6); eingerückte Kopfzeile ohne Anker ist rot.
- **Grenz-Texte ADR-0104, MR-078, Intros:** geprüft, bis auf F-1/F-2 ohne Befund — „Deklaration, kein Beleg" steht überall; das Intro der Go-Datei nennt Plattform, `-tags` und Datei-Kommentar richtig; die in ADR-0104 genannten Testnamen existieren.
- **§3.2:** geprüft, ohne Befund (M8).
- **§3.7:** geprüft, bis auf F-1 ohne Befund — Zusage- und Grenz-Kommentare, keine Befund-Marker, keine Slice-Nummern; die entfernte `DC-FA-CLI-012`-Nennung lässt den Kommentar sonst unverändert.
- **§3.1:** geprüft, ohne Befund — Ableitung und Prüfung in Docker; Proben M2 über `make abdeckung`.
- **Determinismus:** geprüft, ohne Befund — byte-gleich über zwei Läufe und zum committeten Stand (M1).
- **Commit-Botschaft:** geprüft, ohne Befund — Zahlen 43/5/6 nachgemessen (M5); „dieselbe Antwort wie das Go-Werkzeug" trägt für Dateiname und Constraint, F-1 betrifft den Kontext, nicht die Regel.
- **`make gates`:** grün (M9).

## Kategorie-Summary

HIGH 0 · MEDIUM 0 · LOW 2 (F-1, F-2) · INFO 0.
Keine wiederkehrende Klasse mit neuer Instanz; beide LOW sind Reste aus R2 (F-5, F-4) in den Grenz-Texten, nicht im Verhalten.

## Verdikt

Freigegeben. Die beiden MEDIUM aus R2 sind eingelöst und gemessen (M2, M4, M5).
F-1 und F-2: Annahme oder Begründung genügt.
