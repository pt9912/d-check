# Verifikation — slice-297: Die Spezifikation heißt `specification.md`

- **Rolle:** Verifier (Modul 11) — DoD, Plan und Verdikt gegen den Stand
- **Gegenstand:** `docs/plan/planning/in-progress/slice-297-spezifikation-heisst-specification.md` §1, §2, §3 (samt beiden Plan-Änderungen), §4, §6; Architect-Verdikt zu welle-98 (Fragen 4–7, Nachtrag V1); Review-Report [r1](2026-11-26-slice-297-spezifikation-heisst-specification-r1.md) als Kontext; `welle-98-baseline-v700.md` §3 (b), (c), (d1), (g)-Anteile
- **Range:** `f44c6e65..HEAD` (Spitze `5e607ad3`, 9 Commits); Arbeitsbaum sauber, nichts gepusht (`main` 34 Commits vor `origin/main`, davon 9 dieses Slice)
- **Datum:** 2026-11-26
- **Wegwerf:** alles im Scratchpad (`…/scratchpad/v297/`, Klon `…/v297/clone`); Mutationen und Zwischenstände nur im Klon, der Arbeitsbaum blieb unberührt

## Messungen (selbst gefahren, Exit-Code ohne Pipe erfasst)

| # | Kommando | Ergebnis |
|---|---|---|
| M1 | `make gates` (Ausgabe in Datei) | `EXIT=0`; letzte Zeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `coverage-gate: OK — Coverage 94.90% erfüllt Schwelle 93%`; `Ran 55 rules on 75 files: 0 findings.`; `fetch-baseline-cache: verify ok (55 Dateien, vollständig)` |
| M2 | `make doc-check` | Exit 0; `d-check: 1325 Datei(en) geprüft, 0 Befund(e)` |
| M3 | `make adr-check RANGE=f44c6e65..HEAD` | Exit 0; `Range 'f44c6e65..HEAD' aufgeloest, 9 Commit(s) — OK.` / `0 Befund(e)`; zusätzlich je Teilspanne im Klon (`f44c6e65..56c8f62f`, `56c8f62f..HEAD`) `0 Befund(e)` |
| M4 | `make trace-check RANGE=f44c6e65..HEAD` | Exit 0; 9 Commit(s) aufgelöst, `0 Befund(e)` |
| M5 | `make planning-check` | Exit 0; `0 Befund(e)` |
| M6 | `make test` | Exit 0 |
| M7 | `make image-test` | Exit 0; `image-test: (5) --manual — netzlos im Container, nativ == Container, Exit 0`; `image-test: OK — Akzeptanzkriterien erfüllt (DC-FA-DIST-001, DC-FA-CLI-007, DC-FA-CLI-008, DC-FA-CLI-013, DC-FA-MET-001, DC-FA-PSYNC-001, DC-FA-LLEN-001)` |
| M8 | `docker run … d-check:latest --manual manual` | Kopfzeilen `==> docs/user/benutzerhandbuch.md:1400` und `==> spec/specification.md:907`; Zeile 920 der Spezifikation trägt den neuen Pfad |

## Nachzählung §3 (Kommando und Ergebnis)

| Behauptung | Nachzählung |
|---|---|
| Move = ein `R100`, sonst nur `M` | `git show -M --name-status 56c8f62f`: eine Zeile `R100 spec/spezifikation.md spec/specification.md`, 103 `M`, kein `D`/`A` |
| 64 ADRs / 119 Zeilen | `git grep -h '\](\([^)]*/\)\?spezifikation\.md' f44c6e65 -- docs/plan/adr`: 119 Zeilen in 64 Dateien |
| 39 eingefrorene Nicht-ADR-Dateien / 80 Zeilen (23 `done/`, 13 Reviews, 2 CR, 1 `conventions/done`) | Dateizahlen gleich; numstat des Move-Commits 199 + / 199 − = 119 + 80 |
| nur Link-Ziele in den 103 Dateien | je Datei: alter Stand mit `sed` (`](` + optionaler Pfad bis `/` + `spezifikation.md` → `specification.md`) gegen neuen Stand `diff`: 103 von 103 identisch, 0 Abweichung; ADR-Dateien in keinem späteren Commit berührt (`git diff --name-only 56c8f62f..HEAD -- docs/plan/adr` leer) |
| Wortgrenzen-Falle | `git ls-files \| grep 'spezifikation\.md$'`: genau `slice-002-architektur-und-spezifikation.md` und `slice-114-spec-vergabe-spezifikation.md`; beide nicht im Diff |
| 0 Link-Ziele mit altem Namen am Ende | `git grep -c '\](\([^)]*/\)\?spezifikation\.md' \| wc -l` = 0 |
| Summe 406 | `git grep -h '\](…spezifikation\.md' f44c6e65 \| wc -l` = 406 |
| Spezifikation: 161 Zeilen, 158 Link-Ziele, 159 geändert, danach 2 | 161 / 158 am Vorzustand; `diff` (alter Stand normalisiert gegen neuen) nennt genau Zeile 920; `git grep -n 'spezifikation\.md' -- spec/specification.md`: Z. 233, 234 (Konsumenten-Beispiel); Zeile 8 (Pin) unberührt |
| 49 lebende Zeilen in 27 Markdown-Dateien | numstat Commit 2 minus Pläne und Register: 27 Dateien, 49 Zeilen |
| Handbuch nur Z. 2792 | `git diff -U0 f44c6e65..HEAD -- docs/user/benutzerhandbuch.md`: ein Hunk `@@ -2792 +2792 @@`; Z. 372, 373 tragen den alten Namen weiter |
| `.d-check.yml` 6 Zeilen und ein Eintrag | Diff: Z. 446, 517, 521, 760, 848 (Pfade), 991 (Kommentar), Eintrag `"spec/spezifikation.md"` am Ende von `codepaths.ignore-refs` samt Kommentar mit `GRENZE:` |
| Quelltext (d1) | `git diff --stat f44c6e65..HEAD -- '*.go' ':!tools' Dockerfile go.mod`: `manual.go` (2 Zeilen), `cli_manual_test.go` (Z. 21), `diagnose_test.go` (Z. 73); sonst nur `tools/image-test.sh` Z. 163; `Dockerfile`, `Makefile`, `go.mod`, `config_template.go`, `suggest.go` ohne Diff |
| 95 Kommentarzeilen | Zählkommando des Plans an `f44c6e65` und `HEAD`: je 95; die 6 Fixture-Zeilen (`cli_acceptance_test.go`, `matrix_test.go`) und 4 Konsumenten-Zeilen unverändert |
| Lastenheft | Diff 3 Hunks: `**Version:** 0.115.2`, Z. 4762 (Messmethode zu `DC-QA-01`), **eine** Historie-Zeile oben (Form „Redaktionell: …; keine inhaltliche Änderung, kein `DC-*` ändert sich“, Verweis `—`); `DC-FA-CLI-013` unberührt; die Zeile nennt weder ADR noch Slice noch Welle (§3.4) |
| `spec/architecture.md` Z. 26 | ein Hunk, Link-Text und Ziel |

## Ventil (`codepaths.ignore-refs`)

Klon bei `5e607ad3`, Eintrag `"spec/spezifikation.md"` aus der Liste entfernt, `d-check:latest` mit `--network none`: **`1325 Datei(en) geprüft, 120 Befund(e)`**, alle `codepath-missing`, in **80** Dateien (ohne die Kopfzeile des Laufs gezählt). Aufteilung nach Pfad: `done/` 51, ADRs 44, Register-Belege 5, `conventions/done` 6, `CHANGELOG.md` 2, `docs/plan/cr` 1 (= **109** eingefroren); lebend **11**: `welle-98` 4, `slice-297` 4, `slice-296` 2, `MR-081` 1. Mit dem Eintrag: 0 (M2). Das stimmt mit der Endstand-Messung in `welle-98` (Commit `0315fc57`) überein.

**Endzahlen für die Closure-Notiz (F-3):** 119 → **120** Befunde; die Pläne dieser Welle 9 → **10** Nennungen (4 + 4 + 2), lebende Nennungen insgesamt 10 → **11** (mit `MR-081`); eingefroren unverändert 109. Die Formulierung „9 → 11“ des Auftrags trifft die Pläne nicht (9 → 10) und die lebenden insgesamt nicht (10 → 11); die Closure-Notiz sollte die zwei Paare getrennt nennen.

## Bewusstes Brechen (DoD: „beide Tests laufen ohne den Nachzug aus dem richtigen Grund rot“)

Klon bei `5e607ad3`, je **eine** Mutation, `make test` bzw. `make image-test` im Klon:

| Mutation | Beobachtung (gelesen) |
|---|---|
| `cli_manual_test.go` Z. 21 auf `==> spec/spezifikation.md:` | Exit 2; `--- FAIL: TestManual_GibtAbschnitteAus` / `cli_manual_test.go:23: Ausgabe ohne "==> spec/spezifikation.md:"` — die Kopfzeile des Produkts (jetzt `specification.md`) fehlt, aus dem behaupteten Grund |
| `diagnose_test.go` Z. 73 auf `spezifikation.md` | Exit 2; `--- FAIL: TestAllReasonsDeckungGegenSpezifikationGrundCodes` / `diagnose_test.go:42: Spezifikation nicht lesbar (fail-closed): open ../../../../spec/spezifikation.md: no such file or directory` — Meldung wie im Plan |
| `tools/image-test.sh` Z. 163 auf `spezifikation.md` | Exit 2; `image-test: FAIL — --manual: kein Abschnitt aus der Spezifikation`, Abbruch vor `(5)` |
| `manual.go` `//go:embed` auf den alten Pfad (zusätzlich) | Exit 2; `manual.go:22:12: pattern spec/spezifikation.md: no matching files found` — der Build bricht, wie §3 sagt |

## Zwischenstände (§3.3-Zerlegung, `d-check:latest` je Commit im Klon)

| Commit | doc-check |
|---|---|
| `f44c6e65`, `e0a48bb9`, `8a458a27` | 0 Befunde (Pläne vor dem Code) |
| `56c8f62f` (Move) | **rot**: `d-check: error: konfiguriertes ids-Target existiert nicht: spec/spezifikation.md`, Exit 2 (dazu nicht baubar: `go:embed`, durch die Mutation oben belegt) |
| `340f3f15` | **rot**: 212 Befunde (205 `target-missing`, 7 `matrix-forbidden`, alle in der Spezifikation; die Selbst-Links zeigen noch auf den alten Namen), Exit 1 |
| `061ebfc2`, `501dc194`, `00d4f352`, `0315fc57`, `5e607ad3` | 0 Befunde |

Genau die beiden vom Plan vorab benannten Stände sind einzeln rot; die Plan-Änderungen (`e0a48bb9`, `8a458a27`) liegen zeitlich vor dem ersten Code-Commit (`56c8f62f`, 07:22 / 07:26 gegen 07:28) und stehen im Plan.

## DoD-Punkte

- **Rename und Link-Ziele in einem Commit — erfüllt.** Ein `R100`, 103 `M`, nur Link-Ziele (103 von 103 per Substitutionsprobe identisch); `adr-check` über die Spanne und beide Teilspannen 0 Befunde. `make adr-check STAGED=1` am Commit-Moment ist nachträglich nicht fahrbar; der Hook hat den Commit durchgelassen, die Spanne ist belegt.
- **Inhalt und Code — erfüllt.** Selbst-Verweise 158 + Z. 920; lebende Dateien, Konfiguration, `manual.go`, zwei Tests, `tools/image-test.sh` nachgezogen; **ein** `ignore-refs`-Eintrag mit `GRENZE:`-Kommentar in der Form des Vorgängers; `doc-check` 0, `make test` grün; beide Tests laufen ohne den Nachzug aus dem richtigen Grund rot (Meldungen gelesen).
- **Lastenheft — erfüllt.** Version 0.115.2, Z. 4762, eine Historie-Zeile in der Form der übrigen, kein `DC-*` geändert.
- **`make gates` grün; `make image-test` grün — erfüllt** (M1, M7; die Zeile 163 prüft die Kopfzeile, bei Mutation rot).
- **Review durchgeführt — erfüllt.** Report r1 liegt vor (1 MEDIUM an die Welle, 2 LOW, 3 INFO); F-1 durch `0315fc57` (welle-98 führt fünf Klassen, „Eine sechste Klasse ist ein Befund“, Endzahlen des Ventils), F-2 durch `5e607ad3` (`MR-006` gegen `f44c6e65` netto unverändert, `git diff --stat` leer); F-3 steht für die Closure-Notiz an (Endzahlen oben bestätigt).
- **Closure-Notiz mit Lerneintrag — offen.** §7 trägt `*(Wird bei der Closure gefüllt.)*`; Sache des Planners.
- **Beobachtungs-Register fortgeschrieben — offen.** Keine Evidenz-Datei für `slice-297` im Diff (`docs/plan/planning/observations/` nur die vier Link-Ziel-Nachzüge); „keine Beobachtung“ oder eine Datei zu `BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand` (Zahl zählt ihre eigene Nennung, F-3) ist in §7 zu antworten.
- **Jedes Risiko aus §6 trägt einen Ausgang — offen (formal), Inhalt gemessen:**
  1. Mechanische Ersetzung trifft Falsches: **entfallen** (`doc-check` 0; Handbuch-Diff nur Z. 2792; beide Wortgrenzen-Dateien unberührt).
  2. Move-Commit allein nicht baubar: **weiter offen** bis zum Push — nichts gepusht, `git log origin/main..HEAD` zeigt 34 Commits, die Spitze (`5e607ad3`) ist grün; CI am Spitzen-Commit nach dem Push abzulesen.
  3. Test läuft grün ohne die Regel: **entfallen** (drei Mutationen plus Build-Bruch, je rot mit gelesener Meldung).
  4. Repo-weites Ventil verdeckt künftig einen lebenden Verweis: **entfallen** (Kommentar nennt die Grenze mit Bedingung, in der Form des Vorgängers).
  5. Eine Zahl des Plans stimmt nicht: **eingetreten**, benannt und behoben — die Ventil-Zahl 119 → 120 und die lebenden Nennungen 10 → 11 (die Nennungen des Plans zählen sich selbst mit); alle übrigen Zahlen aus §3 reproduzieren sich (Tabelle oben). Die Tabelle in §3 trägt weiter 104/123/161, die Plan-Änderungen überschreiben sie ausdrücklich; die Closure-Notiz nennt die Endzahl.
- **Drei Paarungen — offen, nicht Gegenstand dieser Verifikation** (nächste Welle-Closure).

## Abgrenzung (nichts aus slice-296/298)

Keine Datei zu `slice-296`/`slice-298`, `archive`, `stub`, `wave-` oder einer neuen ADR im Diff (`git diff --name-only f44c6e65..HEAD` gefiltert: leer; `slice-296` und `welle-98` nur als Link-Ziel-Nachzug bzw. die dokumentierte Plan-Änderung am Wellen-Plan). `CHANGELOG.md`, READMEs außer den zwei Link-Zielen, Handbuch-Kopf, `Dockerfile`, `Makefile`, `config_template.go`, `suggest.go`, Pin-Zeile der Spezifikation, `go.mod` unberührt. Die Messungen der Welle (`make blackbox-probe REF=738d2bca`, Vergleich von `--manual` je Begriff, `fullbuild`) sind **nicht** gefahren und gehören der Welle-Closure.

## Befunde

Keine HIGH, kein MEDIUM, kein LOW. Drei INFO.

### V-1 — INFO — Zwei Zwischenstände einzeln rot; Push als eine Einheit

- `56c8f62f` (Exit 2, nicht baubar) und `340f3f15` (212 Befunde) sind einzeln rot, wie der Plan vorab sagt (§3.3-Lesart, Verdikt Frage 4). Nichts ist gepusht; alle neun Commits gehören in **einen** Push (`git log origin/main..HEAD` zählt 34, die neun sind die Spitze). Nach dem Push ist die CI am Spitzen-Commit zu lesen (Risiko 2).

### V-2 — INFO — Die Zahlen des Plans tragen ihre eigene Nennung nicht mit

- Stelle: Plan §3 (zweite Plan-Änderung, „119 Befunde“, „Pläne dieser Welle (9)“) gegen Endstand 120 / 10 / 11. Klasse: Zahl zählt sich selbst mit (Evidenz-Kandidat zu `BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`). Behoben im Wellen-Plan (`0315fc57`); die Closure-Notiz soll die Endzahl und das Paar 9 → 10 / 10 → 11 nennen, nicht „9 → 11“.

### V-3 — INFO — Das Ventil deckt 11 lebende Nennungen mit, die ein Zeilen-Marker nicht kleiner gemacht hätte

- Gleichlautend mit r1 F-4; gemessen: ohne den Eintrag 120, mit 0; der Eintrag ist wegen der 109 eingefrorenen Nennungen ohnehin repo-weit. Kein neuer Befund, nur bestätigt.

## Offene Risiken am Handoff

1. **Push und CI am Spitzen-Commit** (V-1, Risiko 2): bis dahin ist die Zerlegung ein Versprechen, kein Beleg.
2. **Formal offene DoD-Punkte** (Closure-Notiz, Register-Antwort, Risiko-Ausgänge, Paarungen): Sache des Planners; Vorschläge für die Ausgänge stehen oben.
3. **Welle-Closure-Messungen** (`blackbox-probe`, `--manual`-Vergleich je Begriff mit fünf Klassen, `fullbuild`) sind dem Slice nicht aufgetragen und hier nicht gefahren; M8 zeigt nur, dass Kopfzeile und Z. 920 den neuen Pfad tragen.
4. **Kommentarzeilen mit altem Namen bleiben** (95 Go-Zeilen, `Dockerfile` 1, `Makefile` 1, `tools/` 2 Zeilen): Verdikt Frage 7, benanntes Negativ, nicht Teil dieses Slice.
