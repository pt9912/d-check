# Review R1 — slice-270: Die RTM zeigt die Test-Nachweise

**Review-Art:** Code (Diff gegen Plan, ADRs, Konventionen und Hard Rules — nicht gegen die DoD)
**Gegenstand:** slice-270, Range `7753f4d0..HEAD` — `3f79cf3b` (ADR-0104, MR-078),
`d768523d` (Umsetzung)
**Skill:** `.harness/skills/reviewer.md` @ 1.20.0
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Slice-Plan slice-270 samt §3 *Entscheidungen beim Beanspruchen*;
[`DC-FA-COV-001`](../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
und
[`DC-FA-COV-001.a`](../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage);
[`DC-FA-CLI-011`](../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code),
[`DC-FA-DIST-001`](../../spec/lastenheft.md#dc-fa-dist-001--docker-image),
[`DC-FA-CLI-007`](../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus),
[`DC-FA-CLI-008`](../../spec/lastenheft.md#dc-fa-cli-008--reparatur-patch),
[`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit);
[ADR-0104](../plan/adr/0104-test-nachweise-entlasten-in-der-rtm.md) (Proposed),
[ADR-0026](../plan/adr/0026-completeness-in-product-gate.md),
[ADR-0013](../plan/adr/0013-pr-ci-und-traceability-gate.md) (Vergleich der `Schärft:`-Form);
[`MR-078`](../../harness/conventions.md#mr-078); `AGENTS.md` §3.1, §3.2, §3.6, §3.7, §4;
Baseline `v6.17.0` · `templates/docs/plan/adr/NNNN-titel.template.md`,
`templates/harness/conventions/MR-NNN-titel.template.md`,
`regelwerk/grundlagen-traceability.md` §Die zweite Richtung.
Vorherige Findings am selben Modul: keine (neuer Mechanismus).

## Messungen (Kommando und Ergebnis)

| # | Kommando | Ergebnis |
|---|---|---|
| M1 | `make abdeckung`, danach `git status --short` | Exit 0; Arbeitsbaum sauber — die committeten Dateien sind byte-gleich mit der Ableitung; Rechte danach `-rw-r--r--` |
| M2 | Wegwerf-Probe im Scratchpad: `FROM scratch` + `COPY a.md /`, `docker build --output type=local,dest=out` gegen ein `out/` mit `keep.md`, `sub/x` und einem schreibgeschützten `a.md` | `keep.md` und `sub/x` bleiben, `a.md` überschrieben — der local-Exporter ist additiv, er räumt `docs/user/` nicht |
| M3 | `make trace`, Coverage-Spalte gezählt (`awk -F'\|'` über die `DC`-Zeilen) | `44 Tests` · `5 E2E, Tests` · `5 —`; `54 Anforderung(en), 0 Waise(n)` — deckt sich mit der Commit-Botschaft. Die fünf ohne Test: `DC-FA-DIST-002`, `DC-FA-FILE-001`, `DC-FA-RVW-001`, `DC-FA-WF-001`, `DC-QA-01` |
| M4 | `grep -c '^\| \[' docs/user/abdeckung-tests.md` | 358 Zeilen = 356 (Plan) + 2 neue Tests mit Kennung im Doc-Kommentar |
| M5 | alle Kennungen aus `abdeckung-tests.md` gegen die `###`-Überschriften des Lastenhefts (`comm -23`); Kennungen ohne `DC-`-Präfix | keine Kennung ohne Lastenheft-Eintrag; keine fremde Familie |
| M6 | `grep -B3 '^func TestMain\|^func Fuzz\|^func Example'`, `grep -rln '^//go:build'` über `*_test.go` unter `internal/` und `cmd/`; Signaturen `func Test…` ohne `(t *testing.T)` | je 0 Treffer — die Lücken aus F-4 sind heute ohne Instanz |
| M7 | `find internal cmd -name go.mod`; `_test.go` außerhalb `internal/`/`cmd/` | kein Untermodul; außerhalb nur `tools/archive-wave/` (eigenes `go.mod`, nicht in `make test`) — der Scan-Ausschnitt passt zu `go test ./...` |
| M8 | `make lint` · `make gate-consistency` | beide grün (`1101 Datei(en) geprüft, 0 Befund(e)`) |
| M9 | `git diff 7753f4d0..HEAD \| grep nolint` | nur die bestehende Index-Zeile zu `nolintlint`, keine Direktive |

## Findings

### F-1 — MEDIUM — Repo-Selbsttests deklarieren die Produkt-Anforderung `DC-FA-COV-001`

- **quelle:** `AGENTS.md` §5 Regel 16 (Quelle über ihren Geltungsbereich); Skill-Prüffrage 9; Slice-Plan §6 (Risiko „Die Spalte behauptet mehr, als der Test prüft")
- **pfad:** `docs/user/abdeckung-tests.md` · „`TestAbdeckungsDateienFolgenIhrerAbleitung`" (Zeile 15) und „`TestGoTestZeilenZaehltNurDenDocKommentarEinerTestfunktion`" (Zeile 14); Quelle: `internal/adapter/driven/configyaml/abdeckung_test.go` · „RTM gegen ihre Ableitung aus den Testquellen (DC-FA-COV-001)."; `abdeckung_ableitung_test.go` · „jede andere Funktion zählt nicht (DC-FA-COV-001)."
- **befund:** Die zwei neuen Tests prüfen die Ableitungs-Mechanik dieses Repos, nicht die Produkt-Fähigkeit `trace.coverage`; ihr Doc-Kommentar nennt trotzdem `DC-FA-COV-001`, und die eigene Ableitung trägt sie damit als Test-Nachweis dieser Produkt-Anforderung in die RTM ein — und nach ADR-0104 als Entlastung. Das in §6 benannte Risiko tritt im eigenen Diff ein. Failure-Szenario: Verlieren die Produkt-Tests (`TestCLI067_Coverage_RangeDecktAb`, `TestExpandRange`, …) ihre Kennung, bleibt `DC-FA-COV-001` über zwei Tests waisenfrei, die keinen Pfad des Produkts durchlaufen. Dieselbe Zuordnung tragen Makefile-Hilfezeile und Dockerfile-Kommentar.
- **verifizierbar:** ja — `make trace` (Spalte Coverage der Zeile `DC-FA-COV-001`), `grep -n COV-001 docs/user/abdeckung-tests.md`
- **klasse:** selbsttest-deklariert-produkt-anforderung

### F-2 — MEDIUM — ADR-0104 `Schärft:` eine Lastenheft-Anforderung, die sie nicht schärft

- **quelle:** Skill-Prüffragen 9 und 16; Baseline `v6.17.0` · `templates/docs/plan/adr/NNNN-titel.template.md` · Feld `Schärft:` mit Platzhalter `<SPEC-NNN>`
- **pfad:** `docs/plan/adr/0104-test-nachweise-entlasten-in-der-rtm.md` · „— welche Quelle in diesem Repo eine Anforderung von der Waise entlastet."
- **befund:** Die ADR entscheidet eine Konfiguration und Prozess-Regel dieses Repos und ändert am Verhalten von `DC-FA-COV-001` nichts, führt die Anforderung aber unter `Schärft:` — und zwar die Lastenheft-Kennung, wo die Vorlage Spezifikations-/Struktur-Kennungen vorsieht. Die vergleichbaren Prozess-ADRs schreiben „Schärft: keine Spec-Stelle — Prozess-ADR" (ADR-0013, ADR-0026). Failure-Szenario: Die RTM führt ADR-0104 seitdem in der ADR-Spalte von `DC-FA-COV-001` neben ADR-0035/0038/0039/0041 (M3); wer die Architektur-Begründung der Produkt-Fähigkeit sucht, liest dort eine Repo-Policy.
- **verifizierbar:** ja — `make trace`, Zeile `DC-FA-COV-001`, Spalte ADRs
- **klasse:** prozess-adr-schaerft-produkt-anforderung

### F-3 — MEDIUM — Eine Phase in abweichender Kopfzeilen-Form fällt still aus der E2E-Ableitung

- **quelle:** Skill-Prüffrage 18 (Grenzen-Liste ohne ihre größte Lücke); `harness/sensors/test.md`
- **pfad:** `internal/adapter/driven/configyaml/abdeckung_test.go` · „imageTestPhaseRE = regexp.MustCompile(`^# --- (\(\d+\) .*?) -+$`)"; Vertrag in `harness/sensors/test.md` · „Eine Phase ohne Anker ist rot."
- **befund:** Phase ist nur, was exakt diese Form trifft; eine Kopfzeile ohne abschließende Striche oder mit Leerzeichen am Ende ist keine Phase, und fehlt darunter der Anker, entsteht weder Zeile noch Fehler — der Wächter bleibt grün. Failure-Szenario: eine neue Phase `# --- (5) Plattform: arm64` ohne Anker läuft in `make image-test`, fehlt in `abdeckung-e2e.md`, und `make test` meldet nichts. Die Richtung ist Unterzählung, nie falsche Entlastung — deshalb nicht HIGH; der Vertrag nennt die Grenze aber nicht, und die Negativliste prüft nur wohlgeformte Kopfzeilen.
- **verifizierbar:** ja — Kopfzeile in einer Kopie von `tools/image-test.sh` umformen, Anker entfernen, `make test` bleibt grün
- **klasse:** erkennung-still-bei-formabweichung

### F-4 — MEDIUM — Die Testname-Regel zählt, was `go test` nicht als Test fährt

- **quelle:** Skill-Prüffragen 13 und 17; Workflow-Skelett Schritt 19 (`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`)
- **pfad:** `internal/adapter/driven/configyaml/abdeckung_test.go` · „istTestName folgt der Regel von `go test`"
- **befund:** Geprüft werden Präfix und Folgezeichen, nicht die Signatur `func(*testing.T)` und nicht der Build-Constraint der Datei: `func TestMain(m *testing.M)` mit Kennung im Doc-Kommentar zählt als Test-Nachweis, ebenso ein Test in einer Datei mit `//go:build integration` oder `//go:build ignore`, die `go test ./...` nie baut; Verzeichnisse mit `_`/`.`-Präfix, die das Go-Werkzeug überspringt, liest der Walk mit. Failure-Szenario: `// TestMain sperrt das Netz für DC-QA-03.` über `func TestMain` entlastet `DC-QA-03` über eine Funktion, die nichts prüft. Heute ohne Instanz (M6); die Negativliste nennt keinen der Fälle.
- **verifizierbar:** ja — die Fälle in `TestGoTestZeilenZaehltNurDenDocKommentarEinerTestfunktion` aufnehmen, sie wird rot
- **klasse:** testname-regel-ohne-signatur-und-build-constraint

### F-5 — LOW — Der `tidy`-Kommentar steht jetzt über `abdeckung`

- **quelle:** `AGENTS.md` §3.7 (ein Kommentar beschreibt, was da ist)
- **pfad:** `Makefile` · „# go.mod/go.sum pflegen: die Go-Toolchain läuft in Docker" — direkt darunter `abdeckung:`, `tidy:` erst danach
- **befund:** Das neue Target ist zwischen den Kommentarblock von `tidy` und sein Target eingefügt; wer das Makefile liest, ordnet die go.mod-Beschreibung `abdeckung` zu, `tidy` steht ohne Kommentar.
- **verifizierbar:** nein — Urteil, kein Gate
- **klasse:** kommentar-vom-target-getrennt

### F-6 — LOW — Die Liste der Abdeckungs-Dateien steht an drei Stellen

- **quelle:** dieser Skill, LOW-Anker „latente Wartungsfalle (hart verdrahteter Wert, der erst bei künftigem Edit zündet)"
- **pfad:** `Makefile` · „chmod 644 docs/user/abdeckung-tests.md docs/user/abdeckung-e2e.md"; `abdeckung_test.go` · „func abdeckungsDateien()"; `.d-check.yml` · „- files: [docs/user/abdeckung-tests.md]"
- **befund:** Eine dritte Abdeckungs-Datei in der Go-Liste wird geschrieben und gewächtert, bleibt aber `0600` (kein `chmod`) und taucht ohne Nachtrag in `trace.coverage` nicht in der RTM auf — keine der drei Stellen kennt die andere.
- **verifizierbar:** nein
- **klasse:** dateiliste-dreifach-verdrahtet

### F-7 — INFO — Die Negativliste läuft gegen eine Kopie des `id-pattern`

- **quelle:** `BEO-ALL/shared-lexicon-drifts-at-edges`; Slice-Plan §3 („dasselbe `id-pattern` wie die RTM, kein zweites")
- **pfad:** `internal/adapter/driven/configyaml/abdeckung_ableitung_test.go` · „func abdeckungTestPattern() *regexp.Regexp"
- **befund:** Der Wächter liest das Muster live aus der `.d-check.yml`, die zwei Negativtests nutzen eine wörtliche Kopie ohne Kopplungs-Kommentar. Ein Schaden ist nicht erzählbar (die Negativliste prüft Positionen, nicht das Muster); die Plan-Aussage gilt nur für den Wächter.
- **klasse:** muster-kopie-im-fixture

### F-8 — INFO — Plan §1 sagt „zwei Spalten", die RTM zeigt eine

- **quelle:** Slice-Plan §1
- **pfad:** Slice-Plan §1 · „in zwei Spalten: **Tests** (die Go-Suite von `make test`)"
- **befund:** Das Produkt rendert eine Spalte `Coverage` mit den Labels als Werten (`E2E, Tests`, M3); Commit-Botschaft und `.d-check.yml` beschreiben das richtig, der Plan nicht.
- **klasse:** plan-text-gegen-produktform

### F-9 — INFO — Die RTM liest jede Kennung der Abdeckungs-Datei, nicht nur die Kennungs-Spalte

- **quelle:** `DC-FA-COV-001.a` Schritt 2 („alle exakten `requirements.id-pattern`-Treffer")
- **pfad:** `docs/user/abdeckung-e2e.md` · „| Kennung | Phase | Datei |"
- **befund:** Eine Kennung im Phasentitel von `tools/image-test.sh` oder in Einleitung/Grenze der Dateien entlastete ebenfalls; heute keine (M5), Einleitung und Grenze sind feste Strings im Test. Undokumentierte Annahme.
- **klasse:** coverage-liest-ganze-datei

## Negativbefunde

- **Testname-Regel (Präfix/Folgezeichen), Methoden, Doc-Zuordnung:** geprüft, ohne Befund — `Test_unterstrich` zählt, `Testing` nicht, Methoden (`fd.Recv != nil`) nicht; `fd.Doc` hängt nur an der unmittelbar vorangehenden Kommentargruppe (Leerzeile trennt), Rumpf/`t.Run`/Meldung/Testname zählen nicht; Negativliste deckt diese Fälle.
- **testdata, Scan-Ausschnitt:** geprüft, ohne Befund — `testdata` übersprungen; kein Untermodul unter `internal/`/`cmd/` (M7).
- **Determinismus der Ausgabe:** geprüft, ohne Befund — Kennungen je Zeile dedupliziert und sortiert, Zeilen nach Datei, dann Testname; `WalkDir` lexikalisch.
- **Wächter bei fehlender oder abweichender Datei:** geprüft, ohne Befund — fehlende Datei `t.Fatalf`, Abweichung `t.Errorf`, Ableitungsfehler `t.Fatalf`; `.d-check.yml` unlesbar oder ohne `id-pattern` fail-closed.
- **`make abdeckung`:** geprüft, ohne Befund — schreibt nur die zwei Dateien (Scratch-Stage mit `/out`), local-Exporter additiv (M2), reproduziert den committeten Stand (M1); `--no-cache-filter` betrifft nur `abdeckung-gen`, `deps` bleibt gecacht; als Werkzeug mit „kein Gate" im Index, `make gate-consistency` grün (M8).
- **Anker in `tools/image-test.sh`:** geprüft gegen das Lastenheft, ohne Befund — (1) Happy-Path von `DC-FA-DIST-001` plus Identität nativ/Container (`DC-QA-02`); (2) Boundary `:ro` und `--network none` (`DC-FA-DIST-001`, `DC-QA-03`); (3) Negative ohne Mount, Exit 2 mit Hinweis (`DC-FA-DIST-001`); (4) `--doctor`/`--repair` (`DC-FA-CLI-007`/`-008`) byte-identisch (`DC-QA-02`).
- **Zahlen:** geprüft, ohne Befund — 44/5/5 und 0 Waisen (M3); 356 deckt sich mit 358 − 2 (M4). Die 944 aus dem Plan (Commit vor der Range) habe ich nicht nachgemessen; `grep '^func Test'` liefert heute 954, die Zählformen unterscheiden sich.
- **MR-078 gegen die Vorlage:** geprüft, ohne Befund — alle Pflichtfelder, `Ersetzt-Baseline-Regel` als Link mit Anker auf genau eine Regel, Grenze, Auflösungs-Trigger; Index-Zeile vorhanden.
- **ADR-0104 Form, Grenze, Trigger:** geprüft, bis auf F-2 ohne Befund — Abschnitte nach Vorlage, Grenze in den Konsequenzen, zwei beobachtbare Re-Evaluierungs-Trigger, Index-Zeile; Provenance-Marker bei `slice-270` im `Bezug:` zeigt Herkunft, begründet nichts.
- **Spiegel der Waisen-Regel:** geprüft, ohne Befund — `harness/README.md` (beide Zeilen), `harness/sensors/completeness-check.md` (Vertrag und Grenze 1), `harness/sensors/test.md`; Benutzerhandbuch und `operations.md` beschreiben die Produkt-Semantik („ohne Slice **und** ohne Coverage"), die sich nicht ändert; README ohne Waisen-Aussage über dieses Repo.
- **§3.2:** geprüft, ohne Befund (M9).
- **§3.7 übrige Kommentare:** geprüft, bis auf F-5 ohne Befund — Zusage-, Kopplungs- und Grenz-Kommentare; keine Slice-Nummern, keine Review-Historie.
- **§3.1:** geprüft, ohne Befund — Erzeugung in Docker, kein Host-Go.
- **ADR-0005 (Hexagon):** geprüft, ohne Befund — die Testdateien importieren nur `configyaml` und die Standardbibliothek; `make lint` grün (M8). Coverage-Gate: nur `_test.go` geändert, Messbasis unberührt.

## Kategorie-Summary

HIGH 0 · MEDIUM 4 (F-1, F-2, F-3, F-4) · LOW 2 (F-5, F-6) · INFO 3 (F-7, F-8, F-9).
Wiederkehrende Klasse: Deklaration über ihren Geltungsbereich (F-1, F-2) — zweimal im selben Vorgang, das in §6 benannte Risiko.

## Verdikt

Nicht freigegeben. Die vier MEDIUM-Findings sind vor der Closure zu klären — F-1 und F-2
berühren direkt, was die RTM als Entlastung zeigt, F-3 und F-4 die Zusage der Ableitung.
LOW/INFO: Annahme oder Begründung genügt.
