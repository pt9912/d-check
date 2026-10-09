# Verifikation — slice-270: Die RTM zeigt die Test-Nachweise

- **Rolle:** Verifier (Modul 11) — DoD, Plan und Spec gegen den Stand, nicht der Diff gegen Plan/ADR
- **Gegenstand:** `docs/plan/planning/in-progress/slice-270-rtm-zeigt-test-nachweise.md` §1, §2, §3 (Entscheidungen beim Beanspruchen)
- **Range:** `7753f4d0..HEAD` (8 Commits, Spitze `cd2a4d86`)
- **Eingang:** Commits, Reports R1–R3 unter `docs/reviews/`;
  [`DC-FA-COV-001`](../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
  und [§`DC-FA-COV-001.a`](../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage) Schritt 4/5;
  [ADR-0104](../plan/adr/0104-test-nachweise-entlasten-in-der-rtm.md) (Proposed);
  [`MR-078`](../../harness/conventions.md#mr-078)
- **Datum:** 2026-10-09

## Messungen (Kommando und Ergebnis)

Alle Läufe selbst gefahren, nur über make/Docker und Lesewerkzeuge. Brech-Proben
liefen in einem Wegwerf-Klon von `HEAD` im Scratchpad. Nach jeder Probe wurde der
Klon per `git checkout -- .` und `git clean -fd` zurückgesetzt. Der Arbeitsbaum
des Repos blieb unberührt.

| # | Kommando | Ergebnis |
|---|---|---|
| M1 | `make gates` | Exit 0; `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; lint `0 issues.`; `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; doc-check `1104 Datei(en) geprüft, 0 Befund(e)`; die Test-Stage lief mit `--no-cache-filter test`, das Paket `configyaml` war `ok` |
| M2 | `make trace`, Spalte Coverage gezählt | Kopf `… \| Slices \| Coverage \| Status \|`; `Tests` 43 · `E2E, Tests` 5 (`DC-FA-CLI-007`, `DC-FA-CLI-008`, `DC-FA-DIST-001`, `DC-QA-02`, `DC-QA-03`) · `—` 6 (`DC-FA-CLI-012`, `DC-FA-DIST-002`, `DC-FA-FILE-001`, `DC-FA-RVW-001`, `DC-FA-WF-001`, `DC-QA-01`); `54 Anforderung(en), 0 Waise(n).`; keine Anforderung mit Slices `—` |
| M3 | Klon: `make abdeckung` auf dem unveränderten Stand, danach `git status --short --untracked-files=all` | Exit 0, keine Änderung — die committeten Dateien sind die Ableitung; Rechte `-rw-r--r--` |
| M4 | Klon, Mutation A: Doc-Kommentar mit `DC-FA-ANCH-001` über `TestHeadingSlugs_DuplikateUndFences` (`anchors_test.go`) eingefügt, ohne `make abdeckung`; `make test` | rot: `--- FAIL: TestAbdeckungsDateienFolgenIhrerAbleitung` · `abdeckung_test.go:119: docs/user/abdeckung-tests.md weicht von der Ableitung ab — make abdeckung schreibt sie neu`; kein anderer Test rot |
| M5 | Klon, Mutation B: `DC-FA-VER-001` aus dem Doc-Kommentar von `TestDecode_Versions` (`configyaml_test.go:594`) entfernt; `make test` | rot, dieselbe Meldung für `abdeckung-tests.md` (entfernte Deklaration) |
| M6 | Klon, Mutation C: Anker `# abdeckung: DC-FA-DIST-001` unter Phase (3) von `tools/image-test.sh` gelöscht; `make test` | rot: `abdeckung_test.go:101: docs/user/abdeckung-e2e.md: Ableitung scheitert: tools/image-test.sh:98: Phase "(3) Negative: kein Mount → Exit 2 + Mount-Hinweis" ohne Anker darunter` |
| M7 | Klon, Mutation D: Phasen-Text in `docs/user/abdeckung-e2e.md` von Hand gekürzt; `make test` | rot: `docs/user/abdeckung-e2e.md weicht von der Ableitung ab` |
| M8 | Klon, Mutation E: Kopfzeile von Phase (3) in `image-test.sh` umbenannt, Anker bleibt; `make test` | rot: `abdeckung-e2e.md weicht von der Ableitung ab` |
| M9 | Klon, Mutation A, dann `make abdeckung` zweimal, `sha256sum`, `git status --short --untracked-files=all`, `make test` | beide Läufe Exit 0, Hashes identisch; geändert nur `docs/user/abdeckung-tests.md` (genau eine neue Zeile, `DC-FA-ANCH-001` · `TestHeadingSlugs_DuplikateUndFences`) und die mutierte Testdatei; keine weitere Datei angelegt; `make test` danach Exit 0 |
| M10 | Klon: Anforderung `### DC-FA-VRF-001 — Verifier-Probe ohne Slice` in `spec/lastenheft.md` eingefügt; `make completeness-check` | Fehlschlag (make-Exit 2, d-check-Meldung `1 Requirements-Waise(n) ohne referenzierenden Slice und ohne Coverage (--require-complete)`); Zeile `\| DC-FA-VRF-001 \| … \| — \| — \| — \| WAISE \|`; `55 Anforderung(en), 1 Waise(n).` |
| M11 | wie M10, dazu Doc-Kommentar `… prüft DC-FA-VRF-001.` über `TestHeadingSlugs_DuplikateUndFences`, `make abdeckung`, `make completeness-check` | Exit 0; Zeile `\| DC-FA-VRF-001 \| … \| — \| — \| Tests \| ok \|`; `55 Anforderung(en), 0 Waise(n).` — eine Anforderung ohne Slice, deren einziger Nachweis ein Test ist, der sie gar nicht prüft |
| M12 | wie M10, aber `DC-FA-VRF-001` nur als Kommentar im **Rumpf** des Tests; `make abdeckung`, `make completeness-check` | `make abdeckung` lässt die Abdeckungs-Datei unverändert (`git status`: nur Testdatei und Lastenheft); `completeness-check` rot, `1 Waise(n)` — die Stelle im Rumpf entlastet nicht |
| M13 | `git ls-files '*_test.go'` außerhalb `internal/`/`cmd/`; `ls tools/*/go.mod` | nur `tools/archive-wave/*_test.go`, eigenes `go.mod`, nicht Teil von `make test` — der Scan-Ausschluss im Intro stimmt |
| M14 | `grep -n '^\s*# *-\{3,\}' tools/image-test.sh` | vier Kopfzeilen (66, 88, 98, 108), jede mit Anker; keine Kopfzeile mit Strichen, die außerhalb der Ableitung bleibt |
| M15 | `git diff --stat 7753f4d0..HEAD` auf Nicht-Test-Dateien unter `internal/`/`cmd/` | keine — der Produkt-Code ist unberührt; geändert sind nur die zwei neuen Testdateien und ein Doc-Kommentar in `gate_consistency_test.go` |

## DoD-Punkte

**DoD 1 — zwei Abdeckungs-Dateien, abgeleitet, über `trace.coverage` eingebunden; `make trace` zeigt beide Labels.** **Erfüllt.**
`.d-check.yml` bindet `docs/user/abdeckung-tests.md` (Label `Tests`) und
`docs/user/abdeckung-e2e.md` (Label `E2E`) ein. `make trace` zeigt beide Labels in
der Spalte `Coverage`, die vor `Status` steht, wie §`DC-FA-COV-001.a` *Rendering*
es zusagt (M2). Dass die Dateien abgeleitet sind und nicht von Hand gepflegt, belegt
M3: `make abdeckung` auf dem unveränderten Stand ändert nichts. Die Tabellen haben
keine Spalte `Zeile`, weil §3 entschieden hat, Datei und Testname zu binden statt
einer Zeilennummer (siehe V-1).

**DoD 2 — ein Test in `make test` hält jede Datei gegen ihre Ableitung; neue oder entfernte Deklaration rot (bewusstes Brechen belegt); `make gates` grün.** **Erfüllt.**
Ich habe alle Brech-Proben selbst nachgefahren:
- neue Deklaration: M4
- entfernte Deklaration: M5
- Phase ohne Anker: M6
- Handänderung an einer Abdeckungs-Datei: M7
- umbenannte Phase: M8

Jede Probe wird rot, und zwar durch `TestAbdeckungsDateienFolgenIhrerAbleitung` mit
einer Meldung, die den richtigen Grund nennt. Kein anderer Test wird rot.

`make abdeckung` ist deterministisch und schreibt nur die zwei Dateien (M9). Danach
ist `make test` wieder grün. `make gates` ist grün (M1).

**DoD 3 — Review durchgeführt, Report unter `docs/reviews/`; Verifikation.** **Erfüllt mit diesem Bericht.**
R1 und R2 waren nicht freigegeben, R3 ist freigegeben. Die Einarbeitungen liegen in
`6f5cd3a1`, `641650e7` und `cd2a4d86`. Mit diesem Report ist der Punkt erfüllt.

**DoD 4 — Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen.** **Offen, wie es vor der Closure sein muss.**
§7 ist noch leer, und das Risiko in §6 steht auf *(offen)*. Das ist ein Closure-Punkt
und kein Lieferpunkt. Für seinen Ausgang siehe V-2.

## Zusätzliche Prüfpunkte des Auftrags

**Entlastung (ADR-0104 Entscheidung 3, §`DC-FA-COV-001.a` Schritt 5).** **Bestätigt.**
Ich habe eine konstruierte Anforderung ohne Slice in drei Varianten geprüft:
- ohne Test-Nachweis: Waise, `completeness-check` rot (M10)
- mit Kennung im Doc-Kommentar eines Tests: waisenfrei, `Tests`, Exit 0 (M11)
- mit Kennung nur im Rumpf: Waise (M12)

Damit gilt `orphan = (keine Slice-Referenz) ∧ (keine Coverage-Referenz)` wie in der
Spec. Die Aussage von ADR-0104 „Heute ohne Wirkung: die RTM meldet 0 Waisen" stimmt:
Jede der 54 Anforderungen hat einen Slice (M2).

**Plan gegen Code.**

*Was §1 ausschließt, wurde nicht geliefert:*
- **Produkt-Code:** keiner (M15).
- **Neue Tests zum Schließen von Waisen:** keine. Neu sind nur die zwei Selbsttest-Dateien des Wächters. Sie nennen im Doc-Kommentar `ADR-0104` und keine Anforderung, deshalb entlasten sie auch nichts (R1 F-1).
- **Bench- und Bestandsproben als eigene Spalten:** keine.

*Eine Änderung an einem bestehenden Test:* In `gate_consistency_test.go` wurde
`DC-FA-CLI-012` aus dem Doc-Kommentar entfernt (R2 F-1). Das korrigiert eine
Deklaration und schließt keine Waise. Siehe V-3.

*Was §3 verlangt, ist vollständig:*
- beide Dateien
- Ableitung und Wächter im Paket `configyaml`
- `make abdeckung`, als Werkzeug-Zeile `kein Gate`
- Anker in `tools/image-test.sh`, alle vier Phasen
- `trace.coverage`
- ADR-0104 samt Index-Zeile
- `MR-078` samt Index-Zeile
- `harness/sensors/test.md`, Titel jetzt *dreier Zusagen*
- `harness/README.md` (Sensors-Zeile und Werkzeug-Zeile)

`harness/sensors/completeness-check.md` und die Zeilen `completeness-check`/`doc-complete`
im Index sind über §3 hinaus nachgezogen. Das ist folgerichtig, weil ADR-0104 deren
Vertrag ändert.

*Die Entscheidungen aus §3 im Code:*
- **Deklarations-Form:** Gelesen wird mit `go/parser`/`go/ast`, nur `fd.Doc` einer Funktion, die `go test` ausführt; die Kennung über `trace.requirements.id-pattern` aus der dekodierten `.d-check.yml` (`liveConfig`). Das ist eingehalten.
- **Ort:** `docs/user/`. Eingehalten.
- **Bindung:** Datei und Testname, ohne Zeilennummer. Eingehalten.

**Grenz-Texte gegen den Code.** **Stimmen.**
- **Intro `abdeckung-tests.md`:**
  - „unter `internal/` und `cmd/`" stimmt (M13).
  - „was `go test` unter linux/amd64 ohne `-tags` als Test ausführt" entspricht `gebautVonGoTest`, `istTestName` und `nimmtTestingT`.
  - „Kennung im Datei-Kommentar zählt nicht" entspricht `fd.Doc`.
- **Intro `abdeckung-e2e.md` und `harness/sensors/test.md`:**
  - Die Phasen-Definition („mindestens drei Strichen vor einer Nummer in Klammern") entspricht `imageTestPhaseRE`.
  - „ohne Anker rot" belegt M6.
  - „andere Kopfzeilen-Form fällt still aus der Ableitung" ist als Grenze genannt, und heute gibt es keine solche Kopfzeile (M14).
- **`harness/sensors/completeness-check.md` Grenze 1, ADR-0104 Konsequenzen und `MR-078` Grenze:** „eine einzige Kennung im Doc-Kommentar macht sie waisenfrei" ist durch M11 wörtlich bestätigt. Die Grenze ist also nicht nur behauptet, sie ist real.

## Befunde

- **V-1 — LOW — §1 nennt noch `Datei:Zeile`.** Das Ziel in §1 beschreibt die Tabelle als
  `Kennung → Test → Datei:Zeile`. §3 hat entschieden, keine Zeilennummer zu binden,
  und geliefert ist auch keine. Der Plan widerspricht sich in diesem Punkt. Die
  Lieferung folgt der späteren und begründeten Entscheidung. Vor der Closure sollte
  §1 nachgezogen werden, damit niemand das Fehlen der Zeile als DoD-Lücke liest.
- **V-2 — INFO (Closure-Eingang) — das Risiko aus §6 ist nicht hypothetisch.** M11
  zeigt es konkret: Eine Anforderung ohne Slice wird durch die Kennung im
  Doc-Kommentar eines Tests waisenfrei, der sie nicht prüft. Das entspricht dem
  zweiten Re-Evaluierungs-Trigger von ADR-0104. Heute ist es nicht eingetreten
  (M2: jede Anforderung hat einen Slice). Als Ausgang trägt deshalb *weiter offen*
  ins Register, nicht *entfallen*.
- **V-3 — INFO — die Spalte zählt unter, und das Maß ist ungemessen.** Wer eine
  Kennung nur im Datei-Kommentar führt, erscheint nicht in der Spalte. Beispiel:
  Die `TestConfigPath_*`-Tests nennen `DC-FA-CLI-012` nur im Datei-Kommentar, also
  steht dort `—` (M2, R2 F-6). Das ist die Regel und durch §1 gedeckt („schreibt
  keine neuen Tests", keine Umschreibung). Der Leser der RTM muss aber wissen, dass
  `—` in der Spalte Coverage nicht „ungetestet" heißt. Die Datei-Intros sagen das
  nur indirekt („eine Kennung an anderer Stelle … zählt nicht"). Kein Handlungsbedarf
  in diesem Slice, aber ein Kandidat für die Closure-Notiz.

Keine DoD-Verletzung.

## Kategorie-Summary

- **DoD-Verletzung:** 0
- **Plan-Inkonsistenz:** 1 (V-1, LOW)
- **Closure-Eingänge:** 2 (V-2, V-3, INFO)

## Verdikt

**DoD 1–3 sind bestätigt**, DoD 4 ist ein offener Closure-Punkt. Den Wächter habe
ich in fünf Mutationen bewusst gebrochen, jede schlug aus dem richtigen Grund fehl.
Die Entlastung habe ich im Wegwerf-Klon an einer konstruierten Anforderung belegt.
`make gates` ist grün. Der Slice kann in die Closure, nachdem §1 nachgezogen ist
(V-1, optional). Das Risiko aus §6 geht *weiter offen* ins Register (V-2).
