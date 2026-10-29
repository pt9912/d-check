# Verifikation — slice-279: Der Trace erkennt Slices mit Namens-Kennung

- **Rolle:** Verifier (Modul 11) — DoD, Plan und Sensor-Seite gegen den Stand
- **Gegenstand:** `docs/plan/planning/done/welle-92/slice-279-trace-erkennt-namens-slices.md` §1, §2, §3 (Negativliste, Spiegel), §5, §6; `ADR-0110`
- **Range:** `cebdfb39..HEAD` (Spitze `ae8a7a42`); der Review-Report `2026-10-29-slice-279-trace-namens-slices.md` nur als Kontext
- **Datum:** 2026-10-29

## Messungen (selbst gefahren)

| Nr | Prüfung | Behauptet | Gemessen | |
|---|---|---|---|---|
| 1 | `make gates` auf sauberem Arbeitsbaum | grün | Exit 0, Schlusszeile `[gates] … green`, Semgrep 0 Befunde, doc-check/planning-check 0 Befunde, coverage-gate-Stage gebaut | ok |
| 2 | Mutation A: Namens-Fall abgeschaltet (`&& false` am `sliceNameShape`-Treffer) | Namens-Test rot aus dem richtigen Grund | rot: `TestCLI009_Trace_NamensSlice` („Namens-Slice fehlt in der RTM"), `TestTraceSliceNameDefaultCovers` (`got []`), `TestTraceSliceDigitCaseUnchanged` (`got [slice-039]`, fehlen `slice-273`, `slice-qa-sicherheit`), `TestTraceSliceNameIdentity`, `TestTraceSliceDirNamedObservations`, `TestTraceSliceDirUncleaned`; Negativ-Test bleibt grün (erwartet) | ok |
| 3 | Mutation B: `observations`-Ausschluss im Namens-Fall abgeschaltet (`if false && underSegment…`) | je einzeln rot | rot: `TestCLI009_Trace_NamensSlice` („Beleg unter observations/ darf keine Anforderung decken"), `TestTraceSliceNameNegatives` (`got [slice-qa-sicherheit]`), `TestTraceSliceDirNamedObservations`/`TestTraceSliceDirUncleaned` (`slice-x-y` wird fälschlich gezählt); Namens-Happy bleibt grün | ok |
| 4 | Mutation C: Ausschluss VOR dem Ziffern-Fall | `TestTraceSliceDigitCaseIgnoresObservationsExclusion` rot | genau dieser Test rot: „Ziffern-Fall unter observations bleibt Slice: []" (`slice-007` fehlt); alle anderen grün | ok |
| 5 | Wiederherstellung | `git diff` leer | `git status --short` leer, `git diff` 0 Zeilen (Original aus Kopie zurückgespielt) | ok |
| 6 | Default byte-identisch am Ziffern-Fall | `--trace`, `--json`, `--yaml` gleich zu `v0.87.1` | Vorher-Image aus `git archive v0.87.1` (`--target runtime`), Nachher-Image aus HEAD; Korpus = Kopie von `docs/plan/planning` (549 `slice-*.md`, davon 270 Belege unter `observations/`) + `spec/lastenheft.md` + `docs/plan/adr`, ohne `.d-check.yml`, `--network none`, `:ro`. `cmp`: 12188 B / 25253 B / 19795 B identisch, stderr identisch, alle Exit 0 | ok |
| 7 | Fixture Issue 11 (`TDS-QA-03`, `slice-qa-sicherheit.md`) | `ok`, Exit 0 | Vorher-Image: `WAISE`, Exit 1; Nachher-Image: Zeile `… \| slice-qa-sicherheit \| ok`, „0 Waise(n)", Exit 0 unter `--trace --require-complete` | ok |
| 8 | Zählung §3 (Spiegel/Risiko 1) | 549 / 270, alle im Rest unter `evidence/` | `find`: 549; Rest ohne Ziffern-Form: 270 Zeilen, alle `observations/<BEO>/evidence`; ein zweiter Ort in `docs/plan/planning` existiert nicht | ok |

## DoD-Punkte

| Punkt | Befund |
|---|---|
| Test referenziert (Happy/Boundary/Negative) + Rot-Beleg | Erfüllt. Happy `TestTraceSliceNameDefaultCovers`/`TestCLI009_Trace_NamensSlice`, Boundary `TestTraceSliceDigitCaseUnchanged`, Negative `TestTraceSliceNameNegatives`/`…FilePatternOverrides`. Rot-Beleg nachgefahren (Messung 2–4): Fehlermeldung ist die behauptete (Kennung fehlt / Beleg deckt), nicht ein Folgefehler. Akzeptanz ist in `docs/user/abdeckung-tests.md` eingetragen. |
| Fassung nachgezogen | Erfüllt. Lastenheft 0.107.0 (Version, Historie-Zeile, Kriterium *Default byte-identisch* auf den Ziffern-Fall gefasst, neues Kriterium *Namens-Slice*), Spezifikation Schritt 2, Schema-Zeile und Historie, ADR-0110 samt Index-Zeile. |
| Gerüst-Spiegel | Erfüllt. `config_template.go`-Kommentar nennt Ziffern-Fall, Namens-Fall, Ausschluss `observations/` und „gesetzt gilt es unverändert". |
| `make gates` grün | Erfüllt (Messung 1). |
| Review-Report liegt vor | Erfüllt. `2026-10-29-slice-279-trace-namens-slices.md`; F-1 (MEDIUM) und F-2 (LOW) im Fix `26d014d7` behoben, wie die Mutationen C und die Tests `…DirUncleaned`/`…IgnoresObservationsExclusion` zeigen. |
| Doku-Update Schnittstelle | Erfüllt (Spec, Gerüst; Handbuch laut Plan Release-Prep). |
| Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | Noch offen — der Slice liegt in `in-progress/`, §7 „Wird bei der Closure gefüllt", Haken stehen leer. Das ist der Zustand vor Closure, keine Verletzung; `make verify-closure-notes` ist erst mit dem `git mv` nach `done/` sinnvoll und wurde nicht gefahren. |
| Closure-Trigger §5 | Beide Hälften gemessen (Messung 6 und 7). |

## Prüfung (b): „je einzeln rot" (ADR-0110, Fitness-Funktion)

Bestätigt. Ohne Namens-Fall (A) und ohne `observations`-Ausschluss (B) wird der Namens-Test je für sich rot, jeweils mit eigener, passender Meldung; ohne Reihenfolge Ziffern-vor-Ausschluss (C) wird der eigens dafür geschriebene Test rot. Die Mutation A bricht nur den Namens-Fall, B nur den Ausschluss.

## Prüfung (c): Spiegel

Code, Lastenheft, Spezifikation, Gerüst und ADR tragen dieselbe Form: Ziffern-Fall `slice-NNN-…md` zuerst; Namens-Fall nur `[a-z0-9]+(-[a-z0-9]+)*`, Kennung = Basisname ohne `.md`; Ausschluss nur im Namens-Fall, nur unterhalb von `slices.dir` (das Segment in `dir` selbst zählt nicht, `dir` wird bereinigt); gesetztes `file-pattern` gilt unverändert ohne Ausschluss. Die Versionsangabe 0.107.0 steht in Kopf und Historie.

## Befunde / DoD-Verletzungen

Keine DoD-Verletzung. Hinweise ohne Verletzung:

- ADR-0110 steht auf `Proposed` (laut Geschichte „bis die Umsetzung geprüft ist"); der Übergang nach `Accepted` ist Closure-Arbeit und gehört vor den `git mv` des Slice.
- Die Boundary-Aussage „Default byte-identisch" ist am Korpus des Eigenbestands mit Ziffern-Slices und 270 Namens-Belegen belegt, nicht an einem Repo mit Ziffern-Slices *und* fremden Namens-Dateien außerhalb von `observations/` (dort ist die RTM beabsichtigt anders, Risiko 2 im Plan, im Lastenheft gefasst).
- Die Risiken §6 tragen noch keinen Ausgang (Closure). Messung 8 liefert für Risiko 1 die Datenlage: kein zweiter Ort im Eigenbestand.

## Nicht geprüft

- `make verify-closure-notes`, Closure-Notiz, Register-Eintrag und die drei Paarungen (Slice noch nicht geschlossen).
- Die Image-Gegenprobe am `--print-config`-Gerüst-Kommentar zur Laufzeit (nur Quelltext gelesen); `make image-test` nicht gefahren, da nicht Teil der DoD.
- `reviews`, `--suggest-config`, `tools/archive-wave` (laut Plan §1 ausdrücklich ausgeschlossen).
