# Verifikation — slice-298: Neue Namen `wave-<NN>`, `archive.zip`, `seit wave-<NN>`; der Bestand bleibt

- **Rolle:** Verifier (Modul 11) — DoD, Plan und Verdikt gegen den Stand
- **Gegenstand:** `slice-298` §1, §2, §3 (samt den Plan-Änderungen `36f6ac5b`, `73f52639`, `21d6f711`), §4, §6; Architect-Verdikt zu wave-98 (Fragen 1 bis 3, 8, 12, Nachtrag V2); `ADR-0117`; `MR-082`; Review-Report [r1](2026-11-26-slice-298-neue-namen-wave-archive-anker-r1.md) als Kontext; `wave-98` §3 (b), (c), (g), soweit der Slice sie trägt
- **Range:** `db329ff2..HEAD` (Spitze `12fbc567`, 11 Commits); Arbeitsbaum vor diesem Report sauber; `main` liegt 53 Commits vor `origin/main` (Stand der lokalen Referenz)
- **Datum:** 2026-11-26
- **Wegwerf:** Klon im Scratchpad; alle Proben und Mutationen nur dort. Die schweren Läufe liefen nacheinander; ein Abbruch des Laufs (laut Auftraggeber Speichermangel auf dem Host) lag dazwischen, die Messungen davor stammen aus ihren erfassten Ausgaben, in keinem eigenen Lauf stand ein Speicherfehler. `d-check:latest` wurde am Ende mit `make build` aus dem sauberen Baum wiederhergestellt (die Proben hatten es aus dem Klon gebaut)

## Messungen (selbst gefahren, Exit-Code ohne Pipe erfasst)

| # | Kommando | Ergebnis |
|---|---|---|
| M1 | `make gates` | `EXIT=0`; letzte Zeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `coverage-gate: OK — Coverage 94.90% erfüllt Schwelle 93%`; `Ran 55 rules on 75 files: 0 findings.`; `fetch-baseline-cache: verify ok (55 Dateien, vollständig)` |
| M2 | `make archive-wave-test` | Exit 0 (aus dem Cache); zusätzlich ohne Cache im Klon Exit 0, `ok  archive-wave  0.024s` |
| M3 | `make doc-check` | Exit 0; `d-check: 1330 Datei(en) geprüft, 0 Befund(e)` |
| M4 | `make planning-check` | Exit 0; `1324 Datei(en) geprüft, 0 Befund(e)` |
| M5 | `make verify-closure-notes` | Exit 0; `1183 Datei(en) geprüft, 0 Befund(e)` |
| M6 | `make gate-consistency` | Exit 0; `1324 Datei(en) geprüft, 0 Befund(e)` |
| M7 | `make adr-check RANGE=db329ff2..HEAD` | Exit 0; `Range 'db329ff2..HEAD' aufgeloest, 11 Commit(s) — OK.`, `0 Befund(e)` |
| M8 | `make trace-check RANGE=db329ff2..HEAD` | Exit 0; 11 Commit(s), `0 Befund(e)` |

## Nachzählung (Kommando und Ergebnis)

| Behauptung | Nachzählung |
|---|---|
| Move `9d295355` = genau ein `R100` | `git show -M --name-status --summary 9d295355`: eine Zeile `R100` (Welle-Datei), kein `M`/`A`/`D`; über die ganze Spanne ist es das einzige Rename (`R098`) |
| 12 `archiv.zip` im Werkzeug am Vorzustand | `archive.go` 7, `stub.go` 4, `main.go` 1 an `db329ff2` |
| `seit welle-` 36 (17/19) | an `db329ff2`: 36, davon 19 mit Ziffern, 17 ohne; die 17 sind **nicht** alle die Form `<NN>` (V-3) |
| 345 Pfade, 188 `archiv.zip`, 38 Registerzeilen | `HEAD`: 345 / 188 / 38 (an `db329ff2` 346 Pfade, an `590dde3d` 345) |
| `MR-082`-Trigger 435 / 19 | 435 und 19 reproduziert; zur 19 siehe V-1 |
| 38 Ergebnisnotizen, davon ohne `done/welle-` | 38 / 0 |
| `welle-98` am Vorzustand 71 Zeilen / 13 Dateien / 100 Nennungen; Tests 194 Zeilen in 7 Dateien, 16 mit `archiv.zip` | alle reproduziert |
| `welle-98` am `HEAD` | Klassen: Dateiname des Verdikts (Links), Aussagen über den Umzug (Plan §1/§3, Welle §5 Frage 2, Verdikt, `slice-297`, r1 F-4/F-5, der wiederhergestellte Messsatz im Verify-Report zu `slice-297`), Zählkommandos (Plan §3), Fixture-Kennungen in `wave_names_test.go`. Keine lebende Nennung, die brechen müsste (`doc-check` 0) |
| Nachzug nur Namenstausch | Roadmap, `slice-295`, `slice-296`, `slice-297`, vier Reports, Verdikt: nur `welle-98` → `wave-98` |
| `AGENTS.md` §3.7 zeilenneutral | `1 1` im numstat, 362 / 362 Zeilen; die Spanne `d-check:cite AGENTS.md:287-287` in `reviewer.md` hält (`doc-check` 0) |
| `ADR-0117` | `**Status:** Proposed`, Index-Zeile `Proposed`; in der Spanne ändert nur `f30cb082` die ADR (Alternative B) |
| `MR-082` gegen die Vorlage `MR-NNN-title.template.md` (`v7.0.0`) | sechs Pflichtfelder (Datum, Geltungsbereich, Ersetzt-Baseline-Regel, Adaption, Begründung, Auflösungs-Trigger), `Löst auf` zu Recht weggelassen, Hinweis-Block entfernt, Index-Zeile in `harness/conventions.md`, Verweis auf `ADR-0117` ohne Nacherzählung |
| „Bestand bleibt“ | `git diff --name-status -M db329ff2..HEAD -- docs/plan/planning/done docs/plan/planning/observations`: nur `M` an `slice-295` und `slice-297`; `observations/` leer; kein `R`/`A`/`D` unter `done/welle-*` oder `done/wellenlos` |
| Produkt und Defaults | `internal/` ohne Diff; `-- '*.go' ':!tools' Dockerfile go.mod` leer; `spec/` ohne Diff |
| (b) alter Pin | Treffermenge `v6\.18\.1` an `HEAD` und `db329ff2` bis auf den Namen der Welle-Datei gleich; Symlinks: 0 |
| (c) Ventile | die Liste `ignore-refs` ist unverändert; geändert ist nur der Anker `seit wave-98` im Kommentar |

## Probe (g) im Wegwerf-Klon

| Probe | Ergebnis |
|---|---|
| flache `wave-99-probe.md` ohne Zeiger | `make planning-check` Exit 2, `wave-99 wave-drift … nennt es nicht` |
| dieselbe mit Zeiger in „Offene Wellen“ | Exit 0, `1325 Datei(en) geprüft, 0 Befund(e)` |
| Gegenseite: flache `welle-77-probe.md` ohne Zeiger | Exit 0, 0 Befunde (der Bestand ist ungewächtert, wie benannt) |
| `make archive-wave WELLE=wave-99` (Dry-Run) mit einem `wave-99`- und einem `wave-98`-Slice | Exit 0, Plan und **genau** der `wave-99`-Slice |
| derselbe Lauf mit `APPLY=1` | `done/wave-99/archive.zip`, Stub-Zeiger `unzip -p done/wave-99/archive.zip` |
| `welle-97`, `welle-96`, `welle-95` Dry-Run, Werkzeug von `db329ff2` gegen `HEAD`, derselbe Baum | stdout `cmp` gleich, Exit 0/0, Werkzeug-stderr leer (Unterschiede nur im Docker-Build-Log); `welle-97` findet `slice-293`, `slice-294` |
| `wave-99-results.md` und `welle-97-results.md` ohne `## Was wurde geliefert` | `make verify-closure-notes` Exit 2, zwei `section-missing`, je eine |
| `wave-98` im Fließtext von Spezifikation, Lastenheft, Architektur, `ADR-0117` | vier `matrix-forbidden` (`spec-straten`/`spec-straten`/`sicht`/`adr → wave`), Exit 2; `wave-<NN>` und `wave-drift` daneben: kein Befund |
| Gegenseite `welle-97` in der Architektur | `matrix-forbidden` `sicht → welle` |

## Bewusstes Brechen (Mutationen im Klon, je gelesene Meldung)

| # | Mutation | rot | Meldung (Auszug) |
|---|---|---|---|
| T1 | `collect.go` auf `db329ff2` | 4 neue Tests; `TestRun_PlanOhneEigeneSlices_Wave` bleibt grün, wie sein Kommentar nach der Fixrunde sagt | `wave-98 sammelt []`; `erwartet 3 Moves (Plan, zwei Slices), got 1` |
| T2 | Vergleich nur über die Ziffern | `TestCollectSlices_WaveUndWelleGleicheZiffernGetrennt` | `welle-98 sammelt [slice-701-a.md slice-702-b.md slice-703-c.md slice-704-d.md]` |
| T3 | `ArchiveName` immer `archiv.zip` | 5 Tests | `ArchiveName("wave-98") = "archiv.zip", erwartet "archive.zip"` |
| T4 | `ArchiveName` immer `archive.zip` | 8 Tests, darunter Bestands-Tests und `TestFixture_Bestand_SchreibtWeiterArchiv` | `ArchiveName("welle-97") = "archive.zip", erwartet "archiv.zip"` |
| T5 | Suffix `-archiv.zip` | 4 Tests | `erwartet genau slice-601-archive.zip, got [slice-601-archiv.zip]` |
| T7 | Gesamt-Rücknahme (T1+T3+T5) | 11 Tests: **9** neue, 2 bestehende. Grün von den **14** neuen: `ExtractTitle_Wave`, `FindWellePlan_Wave`, `DryRun_NoWrites_Wave`, `Bestand_SchreibtWeiterArchiv`, `AcceptsContextualWaveMention` | — |
| T9 | gezielt: `FindWellePlan` lehnt `wave-` ab; `h1RE` kennt `wave-` nicht; Dry-Run schreibt für `wave-` | `FindWellePlan_Wave`, `ExtractTitle_Wave`, `DryRun_NoWrites_Wave` | `erwartet wave-99-testwelle.md, got ""`; `got ""`; `Dry-Run hat den Baum veraendert` |
| T10 | Zugehörigkeit aus dem ganzen Feld statt der ersten Zeile | `AcceptsContextualWaveMention` (und das `Welle`-Gegenstück) | `wave-86 ist hier nur Kontext: slice-605 gehoert zu einer Welle` |

Jeder der 14 neuen Tests ist damit gegen eine benennbare Änderung rot gefahren; die als Regressionswächter oder Gegenprobe grün bleibenden bleiben unter der Rücknahme grün.

## Zwischenstände (je Commit: `doc-check`, `planning`-Lauf, Closure-Profil mit `d-check:latest`; das Produkt ist über die Spanne unverändert)

| Commit | Ergebnis |
|---|---|
| `73f52639`, `928cb73f` | alle drei 0 Befunde |
| `9d295355` (reiner Move) | **rot:** `doc-check` 28 `target-missing`, `planning` 1 `wave-drift` (`welle-98`); Closure-Profil 0 |
| `21d6f711` (Plan-Datei) | **rot:** `doc-check` 18 `target-missing`, `planning` 1 `wave-drift`; Closure-Profil 0 |
| `2afe2d47` bis `12fbc567` (7 Commits) | alle drei 0 Befunde |

## DoD-Punkte

- **Werkzeug — in der Sache erfüllt, in einer Klausel abweichend (V-2).** Beide Kennungs-Formen gelesen (T1, T2, Probe), Archiv nach Präfix, Slice/Review `-archive.zip` (T5), Bestands-Dry-Run byte-gleich, `make archive-wave-test` grün. „Jeder neue Test läuft ohne die Änderung rot“ trifft für 9 von 14 zu, nicht für alle.
- **Sensoren — erfüllt.** `glob`/`results-glob` auf `wave-*`, Klasse `wave` mit drei Verboten neben `welle`, beide `exempt-paths`, Closure-Profil (eine Regel `*-results.md`, Plan-Änderung vor dem Commit begründet), Grenze im Kommentar, Probe (g) in beiden Richtungen, `doc-check` und `planning-check` 0. Der `exempt-paths`-Eintrag selbst ist hier nur am Diff gelesen; die Probe dazu steht in r1.
- **Commands, Briefing, Eintrag — erfüllt bis auf den Closure-Akt.** Die drei Commands und die Werkzeug-README nennen beide Formen, `AGENTS.md` §3.7 beide Anker, beide Sensor-Dateien ihre Grenze, `MR-082` samt Index, `ADR-0117` trägt die Senkung. **Offen:** `ADR-0117` → `Accepted` (Closure); vorher V-1.
- **`make gates` und `make archive-wave-test` grün — erfüllt** (M1, M2).
- **Review durchgeführt — erfüllt.** r1 liegt vor, kein HIGH/MEDIUM; die Fixrunde ist nachgesehen: F-1/F-2 (Kommentar ohne „heute“-Zahl und ohne Alternativen-Erzählung), F-3 (Kommentare tragen, T1/T3/T9), F-5 (Messsatz nennt wieder `welle-98`), F-6 (ADR nennt 345/19 mit Verweis; siehe V-1).
- **Closure-Notiz — offen.** §7 trägt `*(Wird bei der Closure gefüllt.)*`.
- **Beobachtungs-Register — offen.** Kein Diff unter `observations/`; Risiko 4 verlangt eine neue Beobachtung mit `evidence/slice-298.md`.
- **Risiken mit Ausgang — offen (formal).** Gemessene Vorschläge:
  1. Nur eine Seite: **entfallen**, wenn die Notiz je Sensor beide Seiten nennt (`planning` nur neu, `welle-77` still; `matrix`, Closure-Profil, Werkzeug beide).
  2. Test grün ohne die Regel: die Bedingung „Mutation = Änderung zurücknehmen, jeder rot“ trifft **nicht** wörtlich zu (5 grün); mit gezielten Mutationen sind alle 14 rot. Der Planner entscheidet den Ausgang und benennt die Abweichung (V-2). Dazu r1 F-9: der **bestehende** `welle-`-Ablehnungstest des Slice-Modus bleibt schwach.
  3. Bestands-Lauf geändert: **entfallen** (byte-gleich).
  4. Wächter des Bestands abgeschaltet: **weiter offen** (Register).
  5. Falscher `wave-<NN>`-Treffer: **entfallen** (`doc-check` 0 am Endstand; `git grep -E 'wave-[0-9]{2,}' -- spec docs/plan/adr` leer).
- **Drei Paarungen — offen,** nächste Welle-Closure.

## Abgrenzung

Nichts aus `slice-296` (nur der Namenstausch in seinem Plan), kein Produkt-Merkmal, keine Default-Änderung, kein Release, kein `CHANGELOG`, kein `.claude/worktrees` im Scan-Ignore, kein Bestand umbenannt; die einzige umbenannte Datei ist die Welle.

## Befunde

Kein HIGH, kein MEDIUM. Zwei LOW, drei INFO.

### V-1 — LOW — Die „19 Anker“ in `ADR-0117` und `MR-082` zählen drei Selbstnennungen des Slice-Plans

- `git grep -h -o 'seit welle-[0-9][0-9]*' -- . ':!.harness/baseline' ':!docs/plan/planning/done' ':!docs/reviews' | wc -l` ergibt heute 19; drei Treffer stehen im Plan von `slice-298` selbst (Z. 123 zweimal, Z. 165). Mit zusätzlichem Ausschluss dieser Datei: **16**. Wandert der Plan bei der Closure nach `done/`, ergibt das Kommando 16 (ohne weitere neue Nennungen).
- Folge: `ADR-0117` ginge mit einer Zahl nach `Accepted`, die ihr Verweis-Kommando unmittelbar danach nicht mehr liefert. Solange sie `Proposed` ist, ist sie änderbar. Die Menge enthält außerdem eingefrorene Dateien (`docs/plan/cr`, `harness/conventions/done`); der Trigger „kein lebendes Dokument“ misst damit mehr als lebende.
- An: Architect, **vor** dem Übergang nach `Accepted`.

### V-2 — LOW — DoD-Klausel „jeder neue Test läuft ohne die Änderung rot“ ist wörtlich nicht erfüllt

- T7: 5 der 14 neuen Tests bleiben grün. Es sind Regressionswächter an präfix-generischem Code und zwei Gegenproben, die unter der Rücknahme grün bleiben **müssen**; jeder ist gegen eine andere Änderung rot (T4, T9, T10). Die Closure-Notiz soll das so sagen und Risiko 2 nicht mit „jeder rot ohne die Änderung“ schließen.

### V-3 — INFO — „17 die Platzhalter-Form `seit welle-<NN>`“ im Plan §3 beschreibt die 17 zu eng

- An `db329ff2` sind von den 17 Treffern ohne Ziffern 10 die Form `<NN>`, 2 `<Kennung>`, 5 Prosa oder Zählkommando ohne Platzhalter. 36 und 19 reproduzieren sich.

### V-4 — INFO — Zwei Zwischenstände sind einzeln rot, der Plan benennt einen

- `9d295355` und `21d6f711` (Tabelle oben). Der Plan nennt nur den Move; die Closure-Notiz (noch leer) soll beide nennen (r1 F-4). Alle elf Commits gehören in **einen** Push; die Spitze ist grün.

### V-5 — INFO — Die Test-Zahlen in r1 F-7 sind verzählt

- r1 nennt „fünf der dreizehn neuen Tests“ und „sieben neue, vier bestehende“; gemessen sind es **14** neue, unter der Gesamt-Rücknahme 9 neue und 2 bestehende rot. Die Closure-Notiz soll 14 / 9 / 5 führen.

## Offene Risiken am Handoff

1. **Ungewächterter Bestand** (`ADR-0117`): 38 Register-Zeilen und 38 Ergebnisnotizen melden in `planning` nichts mehr; `welle-77` ohne Zeiger blieb in der Probe still.
2. **Push und CI** am Spitzen-Commit stehen aus.
3. **`ADR-0117` → `Accepted`** mit V-1 vorher; Closure-Notiz, Register, Risiko-Ausgänge: Planner.
4. **Welle-Closure-Messungen** (`fullbuild`, `blackbox-probe`, (a), (d), (e), (f) und die `done/`-Spanne ab dem Vor-Bump-Commit) sind nicht Teil dieses Slice und hier nicht gefahren.
