# Verifikations-Bericht — slice-287 (DoD-/Plan-Konformität)

**Rolle:** Verifier (gegen DoD §2, Plan §1/§3 samt beiden Plan-Nachträgen/§5/§6, Anforderung, Verfahren und ADR; nicht Reviewer, nicht Validator)
**Gegenstand:** Stand `ecfdf8d1` (Commits `63d65ba5`, `640ffa88`, `f8fed40a`, `ecfdf8d1`; Basis `63d65ba5~1`; 17 Dateien)
**Anforderung:** [`DC-FA-PSYNC-001`](../../spec/lastenheft.md#dc-fa-psync-001--gleichheit-eines-wertes-an-mehreren-stellen-modul-pinsync-opt-in)
**Verfahren:** [`DC-FA-PSYNC-001.a`](../../spec/spezifikation.md#dc-fa-psync-001a--gleichheit-eines-wertes-an-mehreren-stellen-pinsync)
**ADR:** [ADR-0113](../plan/adr/0113-pinsync-modul.md)
**Vorgänger:** [`2026-10-31-slice-287-code-r1.md`](2026-10-31-slice-287-code-r1.md)
**Datum:** 2026-10-31

## Selbst gefahrene Sensoren (echte Ausgabe)

- `make gates` → Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`, **Exit 0**. Darin `d-check: 1233 Datei(en) geprüft, 0 Befund(e)`, semgrep `Ran 55 rules on 74 files: 0 findings`.
- `make blackbox-probe REF=v0.90.0` → `blackbox-probe: byte-identisch über 20 Vergleiche (Vorher 37ffaf4f, stdout/stderr/Exit getrennt, Kanarienlauf vorher und nachher)`, Exit 0 (Fixtures `links`, `sauber`, `targets`, `repo`; Formen `-`, `--json`, `--yaml`, `--doctor`).
- `make image-test` → `image-test: OK — Akzeptanzkriterien erfüllt`, Exit 0 (8 s). **Der Lauf enthält keinen `pinsync`-Teil** (Teile 1 bis 7, `metrics` zuletzt).
- `make completeness-check` → `57 Anforderung(en), 0 Waise(n)`.
- Arbeitsbaum nach allen Läufen und Mutationen: `git diff` leer, `git status --short` nur `?? docs/reviews/2026-10-31-slice-287-code-r1.md` (und dieser Bericht).

## Brech-Tests (Modul 11 §Bewusstes Brechen)

Jede Mutation lokal per `sed` eingebaut, `make test` gefahren (nie Host-Go), Fehlermeldung gelesen, Datei per `git checkout` zurückgenommen. Alle Läufe Exit 2 (rot). Die Verstoßklassen laut Closure-Trigger:

| # | Verstoßklasse | Mutation | Beobachtete Fehlermeldung (Auszug) | Richtiger Grund? |
|---|---|---|---|---|
| 1 | nicht-leere Gruppe | `idx[3] <= idx[2]` → `<` (leere Gruppe trifft) | `TestPinsyncKeinTrefferLeereGruppeLeereDatei/leer_fangende_Gruppe`: `got [pinsync-mismatch@b:1], want [pinsync-no-match@b:1]` (ebenso `erste_Passung_ist_leer`) | ja |
| 2 | erste Passung | letzte statt erster Passung | `TestPinsyncMehrdeutigUndMehrfachVorkommenInEinerZeile`: `got [], want [pinsync-mismatch@b:1]`; zusätzlich `.../erste_Passung_ist_leer` | ja |
| 3a | genau ein Treffer (mehr als eine Zeile) | `switch min(len(rd.lines),1)` (ambiguous entfällt) | `TestPinsyncMehrdeutigUndMehrfachVorkommenInEinerZeile`: `got [], want [pinsync-ambiguous@b:2]`; `TestPinsyncWenigerAlsZweiGenauEinTreffer`: `got [pinsync-mismatch@b:1], want [pinsync-ambiguous@b:1]` | ja |
| 3b | mehrere Vorkommen in einer Zeile zählen einmal | zweites Vorkommen in derselben Zeile zählt als weitere Zeile (`FindAllStringIndex`) | `TestPinsyncMehrdeutigUndMehrfachVorkommenInEinerZeile`: `got [pinsync-ambiguous@b:1], want []` | ja (präzise) |
| 4 | Bezug = erste Quelle mit genau einem Treffer | `ref := exact[len(exact)-1]` | `TestPinsyncBezugIstErsteQuelleMitGenauEinemTreffer`: `got [pinsync-no-match@a:1], want [pinsync-no-match@a:1 pinsync-mismatch@c:1]`; `TestPinsyncMismatchNenntBeideWerteUndDateien`: `got [], want [pinsync-mismatch@Makefile:2]`; `TestCLIPinsyncEnable` Exit=0 statt 1 | ja |
| 5 | `\r`-Abschneiden | `TrimSuffix(l, "\r")` entfernt | `TestPinsyncZeilenRandformen/CRLF`: `got [pinsync-mismatch@b:2], want []` und `got [pinsync-no-match@a:1], want []` | ja (die in R1 F-4 bemängelte Trennschärfe ist behoben: Ausdruck jetzt `v=(.+)`) |
| 6a | Symlink als Datei | Zweig nur für Segmente | `TestPinsyncEingabefehlerBeimLesen/Symlink_als_Quelle`: `erwartet Fehler mit "Symlink" und Pin, got ... "link" ist keine Datei` | ja |
| 6b | Symlink als Segment | Zweig nur für das Ende | `.../Symlink_als_Segment`: `got ... "linkdir" ist kein Verzeichnis` | ja |
| 6c | Datei als Verzeichnis-Segment | Zweig `kein Verzeichnis` entfernt | `.../Datei_als_Verzeichnis-Segment`: `erwartet "kein Verzeichnis", got ... "file/x" fehlt` | ja |
| 6d | Verzeichnis als Quelle | Zweig `keine Datei` entfernt | `.../Verzeichnis_als_Quelle`: `erwartet "keine Datei", got ... nicht lesbar: not found: d` | ja |
| 7 | abschließender `/` | `HasSuffix(s.File, "/")`-Prüfung entfernt | `TestDecode_PinsyncEingabefehler/abschließender_Strich`: `erwartet Fehler mit "endet auf '/'", got <nil>`; `TestPinsyncEingabefehlerDerKonfiguration/abschließender_Schrägstrich`: `got ... "Dockerfile" fehlt` | ja |
| 8a | weniger als zwei Quellen | `< 2` → `< 1` | `.../eine_Quelle`: `erwartet Fehler mit "mindestens zwei", got <nil>` (Rand) bzw. `"1 Quelle(n)"` (Kern) | ja |
| 8b | gleiche Stelle zweimal | `combos`-Prüfung entfernt | `.../dieselbe_Stelle`: `erwartet "gegen sich selbst", got <nil>` | ja |
| 8c | doppelter Pin-Name | `seen`-Prüfung entfernt | `.../doppelter_Name`: `erwartet "mehrfach", got <nil>` | ja |
| 8d | aktiv ohne Pin | Prüfung `len(cfg.Pins)==0` entfernt | `.../kein_Pin`: `erwartet "pinsync.pins ist leer", got <nil>`; `TestPinsyncAktivOhnePinIstExitZwei`: `erwartet Fehler, got <nil>` | ja |
| 9 | Konfig vor Lesen | `cfg.Validate()` im Kern entfällt | `TestPinsyncFehlerReihenfolge` rot; zusätzlich `.../file_fehlt`: `erwartet "Pflicht", got ... quelle : "" ist keine Datei`, `.../Pfad_mit_..`: `got ... ".." fehlt` | ja |
| 9b | Pin-Reihenfolge (Namen) | Sortierung absteigend | `TestPinsyncReihenfolgeDeterminismus`: `got [mismatch@b:1 mismatch@d:1], want [mismatch@d:1 mismatch@b:1]`; `TestPinsyncFehlerReihenfolge`: `erwartet den Lesefehler von a/fehlt-a2, got ... pin "b"` | ja |
| 10 | Quellen außerhalb Scan-Menge / ohne Markdown-Endung | Kern liest nur `*.md` | `TestPinsyncGleicherWertKeinBefund`: `got [pinsync-no-match@Dockerfile:1 pinsync-no-match@Makefile:1], want []`; `TestCLIPinsyncEnable` Exit=1 | ja (Proxy; `TestPinsyncQuellenAusserhalbDerScanMenge` selbst rot erst bei 11-Wiring, s. u.) |
| 11 | Bezugsmenge | Zeilentext gekürzt / `Sources++` entfernt | `TestPinsyncNoteNenntDieBezugsmenge`: `Note: "pinsync: 2 Pin(s) geprüft, 4 Quelle(n)"` bzw. `... 0 Quelle(n) gelesen` | ja |
| 12 | Verdrahtung im Lauf | `active["pinsync"] && false` | `TestPinsyncQuellenAusserhalbDerScanMenge`: `Befunde: []`; `TestPinsyncAktivOhnePinIstExitZwei`: `erwartet Fehler, got <nil>`; `TestCLIPinsyncEnable` Exit=0 | ja |
| 13 | Spiegel: Registry | `pinsync` aus `validModules` | `TestPrintConfigVerfuegbarDecktRegistry`: `Verfuegbar-Zeile nennt 26 Eintraege, die Registry kennt 25`; `TestModulRegistrySpiegel_LastenheftCLI002`: `Eintrag "pinsync" ist kein bekanntes Modul (verwaist?)`; `TestCLIPinsyncEnable`: `unbekanntes Modul "pinsync"` | ja |
| 14 | Spiegel: README.md / README.de.md / operations.md / Handbuch-Tabelle / Vorlage | `pinsync` → `pinsyn` je Datei | `TestModulRegistrySpiegel_ReadmeEn` / `_ReadmeDe`: `Modul "pinsync" fehlt`; `_OperationsMdOptionen`: `Modul "pinsync" fehlt`; `_BenutzerhandbuchTabelle`: `Modul "pinsync" fehlt` (+ `TestDocExamples_ConfigBeispieleValidieren`); `TestPrintConfigVerfuegbarDecktRegistry` / `TestPrintConfigPinsyncBlockIstGueltig`: `pinsync-Block der Vorlage fehlt oder ist zu kurz: []` | ja |
| 15 | Spiegel: Lastenheft-Aufzählung | `pinsync` aus `DC-FA-CLI-002` | `TestModulRegistrySpiegel_LastenheftCLI002`: `Modul "pinsync" fehlt` | ja |

Der Fall „Quellen außerhalb `scan.roots`/`scan.ignore`" ist am Kern nicht isoliert mutierbar, weil der Kern gar keine Scan-Menge kennt; geprüft ist er über Mutation 10 (Modul liest nur `.md`) und die Verdrahtungs-Mutation 12, in der `TestPinsyncQuellenAusserhalbDerScanMenge` (Lauf über `Run`) rot wird. Nicht mutiert wurden: Symlink-Wurzel-Fall und der Lesefehler-Pfad außer „fehlend" (`failKindFS`-Test), beide nur gelesen.

## Zählung der Modul-Aufzählungen (selbst)

`grep -rl mentions` gegen `grep -rl pinsync` (ohne `docs/plan`, `docs/reviews`, Baseline, `CHANGELOG`, `.git`): Dateien mit `mentions` ohne `pinsync` sind `.claude/settings.local.json`, `.d-check.yml`, `Makefile`, `harness/README.md`, `harness/sensors/mention-coverage.md`, `gate_consistency_test.go`, `mentions_test.go`, `report_test.go`, `mentions.go`, `workflows.go` — das ist die „Treffer ohne Änderung"-Liste des Plan-Nachtrags (Selbstlauf/Gate-Index = slice-288, Modul-Code, Werkzeug-Einstellung). Alle Aufzählungen der Registry führen `pinsync` und sind durch Mutation 13 bis 15 gefangen: `validModules`, `DC-FA-CLI-002` und Glossar im Lastenheft, beide READMEs, Handbuch-Tabelle, `operations.md`, `--print-config` („Verfügbar" und Block).

## Review-Befunde R1

| Finding | Urteil | Beleg am Code/Doku |
|---|---|---|
| F-1 INFO (ungültiger Block bei nicht aktivem Modul nur im Handbuch) | behoben | Spezifikation `DC-FA-PSYNC-001.a` Grenzen: „Ein ungültiger `pinsync`-Block ist Exit 2 auch ohne aktives Modul (Schritt 1, am Konfigurations-Rand)"; Handbuch weiter benannt. Rest-Bemerkung: das Lastenheft nennt es weiter nicht (nur Spezifikation + Handbuch). |
| F-2 LOW (Selbstvergleich syntaktisch) | behoben | Spezifikation Grenzen und Handbuch (Zeile 1655 ff.); Test `TestPinsyncGleichwertigeAusdruckeSindKeinSchutz` trägt die Grenze im Kommentar. |
| F-3 INFO (Testkommentar, fehlende Kennung) | behoben | `TestPinsyncLaeufeSindGleich` behauptet nur noch „Zwei Läufe liefern dieselben Befunde"; `TestPinsyncKindFehlerIstEingabefehler` trägt Kennung und steht in `docs/user/abdeckung-tests.md` (Zeile 322). |
| F-4 INFO (CRLF-Teilfall ohne Trennschärfe) | behoben | Teilfall nutzt `v=(.+)`; Mutation 5 wird rot. |
| F-5 INFO (FIFO/Gerätedatei/Socket) | begründet offen (benannt) | Spezifikation Grenzen nennt die Klasse („der Lauf hängt, statt mit Exit 2 zu enden"); Handbuch ebenso. Nicht im Lastenheft. Kein Code-Fix, bewusst. |
| F-6 INFO (Doctor-Text-Jargon) | behoben | `pinsync-no-match` lautet jetzt „… (die Zeilenangabe 1 ist nur eine Konvention und keine Fundstelle)". |
| F-7 INFO (dritter Zeilen-Leser) | INFO, offen (vorgabegemäß) | `pinsyncLines` neben `splitLines` und `TrimSuffix` in `metrics`; kein dritter Beleg der Beobachtungsklasse. |
| F-8 INFO (Vorlagen-Block/CLI-Test) | Vorlagen-Teil behoben, CLI-Teil INFO | `TestPrintConfigPinsyncBlockIstGueltig` (Mutation 14 rot) und `TestCLIPinsyncEnable` (Exit 1/0/2, `summary.notes`) liegen vor; die Rest-Bemerkung betrifft keinen offenen Fehler. |

## Urteil je DoD-Punkt

Haken im Plan stehen alle noch `[ ]` (Closure-Schritt noch nicht erfolgt, kein Haken durch den Implementer gesetzt).

1. **Konfiguration lesen, am Rand prüfen, Exit 2 je Eingabefehler, `--print-config` und Vorlage tragen den Block:** **bestätigt.** Jede Klasse (leer/fehlend, Name, doppelter Name, <2 Quellen, Pflichtfelder, `..`/absolut/Glob/Steuerzeichen/abschließendes `/`, Ausdruck ungültig/Gruppenzahl, gleiche Stelle, fehlende Datei, Verzeichnis, Symlink, Lesefehler) hat einen Test, der unter den Mutationen 6 bis 9b aus dem richtigen Grund rot wird. Unbekannter Schlüssel: durch den strikten YAML-Decoder des Bestands (`yaml:"..."`-Strukturen, gleiche Familie wie die übrigen Module), kein pinsync-eigener Test gebrochen — Lücke nur dort, siehe Offene Punkte.
2. **Fang je Quelle zeilenweise, drei Befund-Klassen, Grund-Codes, Reihenfolge, jede Randbedingung mit Test inkl. Negativfällen aus slice-286 §3:** **bestätigt.** Mutationen 1 bis 5, 9b, 11 rot aus dem richtigen Grund. Die Negativfälle des Plans (Kommentarzeile, `FROM`-Zeile, `v`-Präfix, CRLF, ohne Schlussumbruch, leere Gruppe, leere Datei, Ausdruck ohne Gruppe) sind je mit Testnamen vorhanden.
3. **Modul registriert (Modulliste, `--enable`, Spiegel-Deckungstests, `gate-consistency`), Handbuch mit `v?`-Beispiel:** **bestätigt.** Mutationen 13 bis 15 rot; Handbuch §4.27 zeigt `v?(\S+)` (Zeilen 1601 f., 1605); Aufzählungen selbst gezählt, Treffer ohne Änderung stimmen mit dem Plan-Nachtrag überein.
4. **`make gates` grün:** **bestätigt** durch eigenen Lauf, Exit 0.
5. **Review durchgeführt, Report liegt vor:** **bestätigt** (R1, Kein HIGH/MEDIUM); der Report ist noch nicht committet (untracked), Haken im Plan offen.
6. **Doku-Update für die Modul-Aufzählungen:** **bestätigt** (README×2, Handbuch, `operations.md`, `abdeckung-tests.md`, Lastenheft, Spezifikation §4 mit `SPEC-109` bis `SPEC-111`, Lockstep-Test grün).
7. **Closure-Notiz mit Lerneintrag:** **offen** — §7 „Entsteht bei der Closure".
8. **Beobachtungs-Register fortgeschrieben:** **offen** — nichts angelegt; §8 benennt `BEO-ALL/gemeinsame-leser-sind-zwei-familien` und `BEO-ALL/spec-randbedingung-ohne-test` als Kandidaten.
9. **Risiken §6 mit Ausgang:** **offen** (Ausgänge entstehen bei der Closure). Vorschlag aus dem Befund: Risiko 1 *entfallen* (Kern liest über den vorhandenen `Filesystem`-Port, kein neuer Port); Risiko 2 *entfallen* (RE2-Fehler = Exit 2, Test belegt die Meldung); Risiko 3 *entfallen* **nicht belegt**: der Beleg „`image-test` unter `:ro`" liegt weder hier (`make image-test` ohne `pinsync`-Teil) noch bisher vor — er gehört laut Plan zu slice-288; der Ausgang muss auf slice-288 zeigen (*eingetreten*/Folge-Slice mit Kennung) oder *weiter offen* sein, nicht *entfallen*.
10. **Drei Paarungen:** **offen** (nächste Welle-Closure).

Closure-Trigger §5: „Verifikation hat die Testbehauptungen der Anforderung gebrochen und rot gesehen" — **erfüllt** durch diesen Bericht (Zeilen 1 bis 15 der Tabelle).

## Abgleich Anforderung ↔ Code (Akzeptanzkriterien)

Happy Path, ungleiche Werte, kein Treffer, mehrdeutig, Randformen, Eingabefehler, Bezugsmenge, Scan-Menge, Reihenfolge der Eingabefehler, weitere Eingabefehler, Treffer (`a(b*)`), Reihenfolge, Determinismus, Default byte-identisch: durch Test und Mutation belegt (Default byte-identisch zusätzlich durch `blackbox-probe`: 20 Vergleiche gleich). **Nicht in diesem Slice belegt:** Kriterium „Betrieb" (`:ro` und fremder Eigentümer am Image) und „Eigener Bestand" (Selbstlauf im eigenen Profil) — beides gehört laut §1 und §6 zu slice-288. „Kein Schreiben" ist hier nur durch den Quelltext belegt: `pinsync.go` nutzt allein `Kind` und `ReadFile` des Filesystem-Ports, keine Schreib-API, kein `exec`, kein Netz.

## Abweichungen und offene Punkte

- **V-1 (keine DoD-Verletzung):** Kriterium „Betrieb" der Anforderung ist durch `make image-test` nicht abgedeckt (kein `pinsync`-Teil). Plan §6 Risiko 3 verschiebt den Beleg ausdrücklich nach slice-288; bei der Closure darf das Risiko deshalb nicht *entfallen* sein, solange der Beleg dort aussteht.
- **V-2 (LOW):** Der Eingabefehler „unbekannter Schlüssel" im `pinsync`-Block (DoD 1) hat keinen eigenen `pinsync`-Test (nicht mutiert). Er folgt aus dem strikten Decoder des Bestands; wer ihn als Randbedingung zählt, braucht einen Test (`BEO-ALL/spec-randbedingung-ohne-test`).
- **V-3 (INFO):** Mutation 3b/10 wurden mit Mutanten geführt, die breiter rot werden als ihre Klasse; die präzisen Tests (`TestPinsyncMehrdeutigUndMehrfachVorkommenInEinerZeile`, `TestPinsyncQuellenAusserhalbDerScanMenge`) sind trotzdem rot gesehen (3b direkt, Scan-Menge über Mutation 12).
- **V-4 (INFO):** R1-F-1 und F-5 stehen in Spezifikation und Handbuch, nicht im Lastenheft; die Anforderung trägt die FIFO-Grenze und das Verhalten des Blocks bei nicht aktivem Modul nicht.
- **V-5 (Kosmetik):** alle DoD-Haken im Plan noch leer; Review-Report R1 untracked.

## Verdikt

Keine DoD-Verletzung. Die Punkte 1 bis 6 sind durch eigenen `make gates`-Lauf (Exit 0), `blackbox-probe` (byte-identisch) und über dreißig Mutationen der Testbehauptungen (alle rot aus dem richtigen Grund, Arbeitsbaum zurückgesetzt) bestätigt; die Review-Befunde R1 sind bis auf die benannten INFO-Reste (F-5, F-7, F-8-CLI) im Code beziehungsweise in der Doku behoben. Die Punkte 7 bis 10 sind erwartungsgemäß offen (Closure). Bei der Closure zu beachten: Risiko 3 (`:ro`-Beleg) zeigt auf slice-288 und kann nicht *entfallen* heißen, solange dort kein `image-test`-Beleg liegt.
