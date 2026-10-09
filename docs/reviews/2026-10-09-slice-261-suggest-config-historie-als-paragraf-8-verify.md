# Verifikation slice-261 — `--suggest-config` kennt die Spezifikations-Historie als §8

**Rolle:** Verifier („Bauen wir es richtig?" — gegen Plan und DoD, nicht gegen den Diff-Stil).
**Gegenstand:** `docs/plan/planning/in-progress/slice-261-suggest-config-historie-als-paragraf-8.md`, Stand `fa995ce8`. Geprüft in einem frischen Klon, nicht im Arbeitsbaum.
**Datum:** 2026-10-09.

## Summary

Verdikt: **konform** für die DoD-Punkte 1 und 2. Punkt 3 ist mit diesem Bericht zur Hälfte erfüllt, die Verifikation liegt vor. Ein Review-Report fehlt noch. Punkt 4 (Closure) ist erwartungsgemäß offen. Den Fix habe ich zurückgenommen: Dann wird genau der neue Test rot, und zwar aus dem behaupteten Grund. Black-Box-Probe und `make gates` habe ich selbst gefahren. Das Image erzeugt den Vorschlag mit der neuen Zeile, und die Ausgabe weicht gegenüber dem Vorher-Image nur in dieser einen Zeile ab. Es gibt keine DoD-Verletzung, nur zwei INFO.

## Sensor-Belege (selbst gefahren)

| Sensor | Ergebnis |
|---|---|
| `make test` mit zurückgenommenem Fix | `EXIT=2`. Als einziger Fehlschlag: `--- FAIL: TestCLI006_AiHarness_SchlaegtBeideHistorienVor`, Meldung `cli_acceptance_test.go:1033: exclude-sections ohne "8. Historie": [Historie 7. Historie Geschichte]` |
| `make blackbox-probe REF=HEAD~1` (Fix eingesetzt, Baum sauber) | `EXIT=0`, `byte-identisch über 20 Vergleiche (Vorher e107a476, stdout/stderr/Exit getrennt, Kanarienlauf vorher und nachher)` |
| `make gates` | `EXIT=0`, `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`. Einzelbelege: `verify ok (54 Dateien, vollständig)`, `0 issues.`, `coverage-gate: OK — Coverage 94.80% erfüllt Schwelle 93%`, `Ran 55 rules on 69 files: 0 findings.`, `d-check: 1153`/`1159 Datei(en) geprüft, 0 Befund(e)` |
| Black-Box `--suggest-config ai-harness` (Klon unter `/repo`, `--network none`) | `d-check:latest` (aus `fa995ce8`): Zeile 54 `  exclude-sections: [Historie, "7. Historie", "8. Historie", Geschichte]`. `d-check:probe-vorher` (aus `e107a476`): `  exclude-sections: [Historie, "7. Historie", Geschichte]`. Der `diff` der beiden Ausgaben zeigt nur diese Zeile (`54c54`). Beide Läufe enden mit Exit 0. |

## DoD-Abgleich

| DoD-Punkt | Befund | Beleg |
|---|---|---|
| 1. `8. Historie` vorgeschlagen; Test ohne Fix aus richtigem Grund rot; Spezifikation §2 nachgezogen | erfüllt | Die Zeile in `renderHarnessMatrix` (`internal/hexagon/core/app/suggest.go`) ist ergänzt. Wird der Fix zurückgenommen, schlägt nur der neue Test fehl, und die Meldung nennt das fehlende `"8. Historie"` samt dem tatsächlichen Listeninhalt. Damit ist es kein Dekodier- oder Exit-Fehler. Der Test dekodiert die Ausgabe über `configyaml.Decode`, prüft also den Wert und nicht den Text. `spec/spezifikation.md` Zeile 243 ist byte-gleich mit der Ausgabezeile des Images. Die Historie-Zeile (§8, 2026-10-09) ist ergänzt. Weitere Stellen mit der vorgeschlagenen Liste habe ich gesucht (README, Handbuch, operations.md, Lastenheft) und keine gefunden. |
| 2. Ohne Option unverändert (Probe gegen Vorher-Stand); `make gates` grün | erfüllt | Probe und Gates siehe oben |
| 3. Review + Verifikation | teilweise | Dieser Bericht liegt vor, ein Review-Report zu slice-261 unter `docs/reviews/` noch nicht |
| 4. Closure | offen (erwartet) | §7 ist leer, ebenso der Ausgang für §6 („Keine bekannt"). Beides ist vor dem `git mv` nach `done/` zu leisten. |

Plan §3 nennt zwei Zeilen, und der Diff berührt genau diese. Hinzu kommen der Test und die abgeleitete Zeile in `docs/user/abdeckung-tests.md`, die `make test` gegen ihre Ableitung hält. Die Abgrenzung aus §1 ist eingehalten: Kein anderer Vorschlag wurde geändert, und `.d-check.yml` ist unberührt.

## Befunde

- **INFO-1 — Der Modus `ai-harness-init` ändert sich mit, ohne Test.** `renderHarnessMatrix` erzeugt den Block für beide Harness-Modi. `--suggest-config ai-harness-init` gibt jetzt ebenfalls `"8. Historie"` aus (Zeile 55, selbst gefahren). Im Sinne des Plans ist das richtig, denn das Ziel nennt `--suggest-config` ohne Modus. Die Commit-Botschaft nennt aber nur `ai-harness`, und der neue Test fährt nur diesen Modus. Ein künftiges Auseinanderziehen der beiden Pfade fiele für `ai-harness-init` nicht auf. Das ist keine DoD-Verletzung.
- **INFO-2 — Die Black-Box-Probe sagt nichts über die Option selbst.** Ihre 20 Vergleiche fahren die Prüfläufe der Fixtures, nicht `--suggest-config`. Für DoD-Punkt 2 („ohne die Option unverändert") ist das die richtige Menge. Dass sich die Ausgabe mit der Option nur in der einen Zeile ändert, belegt der eigene `diff` oben, nicht die Probe.

## Grenzen dieses Berichts

Den Review-Gegenstand habe ich nicht geprüft: Stil, Kommentar-Klassen und Hard Rules am Diff bleiben beim Reviewer. Zur Validierung gegen den realen Bedarf, also ob Adopter die Vorlage `v6.17.0` tatsächlich mit §8 nutzen, sagt dieser Bericht ebenfalls nichts.
