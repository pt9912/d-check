# Verifikation slice-283 — Modul `metrics`, hermetischer Kern

**Rolle:** Verifier (Modul 11) — "Bauen wir es richtig?" gegen DoD (§2) und Plan (§1, §3 samt Plan-Änderungen)
**Gegenstand:** HEAD `c7a23b70` (Feat `c9478b69`, Plan-Änderung `5deb73b2`, Review-Fix `c7a23b70`)
**Datum:** 2026-10-30
**Eingang:** DoD-Bestätigung des Implementers; Review-Report `docs/reviews/2026-10-30-slice-283-r1.md`. Die Sensoren sind selbst nachgefahren, behauptete Exit-Codes zählen nicht.

## Selbst gefahrene Sensoren

| Sensor | Ergebnis |
|---|---|
| `make gates` | Exit 0; letzte Zeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; semgrep 0 Findings; doc-check, targets und planning je `1225 Datei(en) geprüft, 0 Befund(e)` |
| Image-Lauf (`d-check:latest`, `--network none`, Fixture unter dem Scratchpad `verify283/`) | siehe unten, alle fünf erwarteten Ausgänge beobachtet |
| Bewusstes Brechen 1 und 2 (`make test`) | beide rot aus dem richtigen Grund, beide mit `git checkout --` zurückgenommen, `git diff --quiet` danach leer |
| `make verify-closure-notes` | nicht gefahren: der Slice liegt noch in `in-progress/`, die Closure-Notiz (§7) ist leer ("Wird bei der Closure gefüllt"). Gehört zur Closure, nicht zu dieser Verifikation |

### Image-Fixture (Konfiguration `modules: [links]` plus `metrics: plans-dir/block/surcharge/rounding/cut-threshold`)

| Fixture | Aufruf `--enable metrics` | Beobachtet |
|---|---|---|
| korrekt (E = 6790 ⇒ 9200 mit `cut-reason` und Ausnahme mit Grund; E = 2222 ⇒ 3000; ein Plan ohne Block) | Exit **0** | `metrics: 3 Plan/Pläne geprüft, 2 mit Block`, Hinweis `kein Block (nicht gemessen)`, 0 Befunde |
| falsche Marke / leerer Grund / fehlendes Pflichtfeld (je ein Plan) | Exit **1** | `slice-1.md:7 limit metrics-limit-mismatch ... (erwartet 1400)`; `slice-2.md:8 exception metrics-reason-empty`; `slice-3.md:6 cut-reason metrics-cut-reason-missing` (E = 2223) |
| Block-Zeile `path src/` ohne `:` (zweiter Plan) | Exit **2** | `error: metrics: plans/slice-2.md: zeile 4: kein ':' in "path src/" (DC-FA-MET-001, fail-closed)`; kein Befund aus slice-1 |
| korrektes Fixture ohne `--enable metrics` | Exit 0 | `0 Datei(en) geprüft, 0 Befund(e)` — das Modul bleibt ohne Aktivierung stumm |
| Config ohne `surcharge`/`rounding` | Exit **2** | `metrics.surcharge.numerator, ... sind Pflicht` |
| Config mit `metrics.excluded` | Exit **2** | `field excluded not found` (so von §1 des Plans vorgesehen) |
| `--print-config` | | Modulliste nennt `metrics`; Block `# metrics:` mit den Schlüsseln dieses Slice im Gerüst |

Nicht belegt: der in §5 verlangte Byte-Vergleich des Laufs ohne `metrics`-Konfiguration gegen das Vorher-Image (`make blackbox-probe REF=...`) wurde hier nicht gefahren; belegt ist nur, dass das Modul ohne Aktivierung nichts ausgibt. Das ist ein Closure-Kriterium (§5), kein DoD-Punkt aus §2.

## Bewusstes Brechen (Arbeitsbaum danach sauber)

1. Duplikat-Prüfung in `readMetricsFields` ausgeschaltet (`if false && seen[key] && ...`): `make test` Exit 2. Rot sind `TestMetricsEingabefehler` ("Schlüssel doppelt: erwartet Fehler, got {Findings:[] ... WithBlock:1}") und `TestMetricsDuplikatePruefungJeSchluessel` ("Duplikat \"base: 9a0efa5\\n\" muss ein Fehler sein"). Der Grund ist der behauptete.
2. Schwellen-Vergleichszeile `scaled <= cfg.CutThreshold` auf `<` geändert: `make test` Exit 2. Rot ist `TestMetricsSchwelleGenauAufDerGrenze` ("auf der Schwelle kein Befund, got ... metrics-cut-reason-missing"). Beachtlich: die belegten Zahlen 2222/2223 allein fangen die Mutation nicht (135 · 2222 = 299970 und 135 · 2223 = 300105 liegen beidseits von 300000); erst der mit R1-F-3 nachgetragene Test mit Schwelle 135 · 2000 trägt den Gleichheitsfall. Beide Rücknahmen: `git checkout -- internal/hexagon/core/rules/metrics.go`, danach `git diff --quiet` leer.

## DoD-Punkte

**1. Registrierung und Konfiguration — bestätigt, mit benannter Lücke zur Formulierung**
- `ValidModules()` enthält `metrics` (`model/config.go`), `Config.Metrics`, `MetricsConfig`.
- Config-Rand getestet (`TestDecode_MetricsRand`, 14 Fälle): fehlendes/leeres `plans-dir`, fehlender `block`, fehlendes `surcharge`/`numerator`/`rounding`, Teiler 0, `rounding` 0, negative Werte, gleiche Feldnamen, unbekannter Schlüssel (auch in `fields`), falscher Typ — alle Fehler; am Image Exit 2 belegt.
- Gerüst führt den Block (Image-Lauf, `config_template.go`). Spiegel-Test grün (`make gates`), Spiegel-Orte in Lastenheft (DC-FA-CLI-002, Glossar), `operations.md`, beide READMEs, Handbuch-§6 nennen `metrics`. `registry_mirror_test.go`, `gate_consistency_test.go`, `Makefile`, `.d-check.yml` und `harness/` sind im Diff unverändert (geprüft über `git diff --stat HEAD~6 HEAD`).
- Lücke gegenüber dem DoD-Wortlaut: "Faktor ohne Eingabe → Exit 2" ist nicht umgesetzt, weil `estimate-inputs` per Plan-Änderung `5deb73b2` nach slice-285 wanderte. Das ist in §1 begründet und mit Adresse versehen, der DoD-Text in §2 wurde aber nicht angepasst; er verlangt weiter etwas, das §1 ausschließt. Vor dem Abhaken ist §2 an §1 anzugleichen.
- Folge: Spec §2 führt die Schlüssel `estimate-inputs`, `table.*`, `excluded` als Schema, der Config-Rand lehnt sie ab (Exit 2, Image belegt). Die Spec darf slice-284/285 nicht nennen (§3.4), ein Nutzerhinweis fehlt also. Das Gerüst und die Handbuch-Zeile sagen es nicht. Als Restrisiko in §6/§7 zu führen (Ausgang: weiter offen oder Folge-Slice).

**2. Prüfung "Obergrenze = Formel" — bestätigt**
- Die belegten Zahlen stehen als Tests (`TestMetricsFormelBelegteZahlen`: 6790⇒9200, 5036⇒6800, 5430⇒7400, 9215⇒12500, 4670⇒6400, 2000⇒2700, 2001⇒2800, 2222⇒3000, 2223⇒3100 mit `cut-reason`, 0⇒0); Mismatch-Befund mit Datei, Zeile 7, Ziel `limit`; Schwellenregel (`TestMetricsCutReasonMissing`, `...GenauAufDerGrenze`); Eingabefehler als Exit 2 ohne Befund (`TestMetricsEingabefehler`, 19 Fälle; `TestMetricsErsterFehlerBeendetDenLauf`), Überlauf (`TestMetricsUeberlauf`), unbekannte Block-Felder ignoriert (`TestMetricsZaeune`, CRLF-Fall). Rot-Beleg für Duplikat-Prüfung und Schwellen-Vergleich erbracht (oben). Am Image: Exit 1 bei falscher Marke und fehlendem Pflichtfeld, Exit 2 bei Zeile ohne `:`.
- Beobachtung (kein Befund): die Negativliste nennt "Block in einer Einrückung" unter den Fällen. `TestMetricsZaeune` legt fest, dass ein vierfach eingerückter Zaun als Block zählt (Öffner-Test des Produkts, §6 Risiko 1 im Plan); ein Markdown-Leser hält ihn für einen eingerückten Code-Block. Das ist die im Plan als Risiko 1 geführte Abweichung; sie braucht bei der Closure einen Ausgang.

**3. Prüfung "leerer Grund", Plan-Suche — bestätigt**
- `TestMetricsReasonEmpty` (leer und nur Leerraum, unabhängig von der Messung; der Befund trägt Ziel `exception`); `TestMetricsPlanSuche` (nur `slice-*.md` unmittelbar, `notes.md`, `sub/slice-9.md`, `slice-1.md.bak` bleiben außen; ohne Datei grün mit Hinweis "kein Plan"; fehlendes oder dateiförmiges Verzeichnis ⇒ Fehler). Symlink-Plan ⇒ Fehler (`TestMetricsPlanSymlinkIstEinFehler`, Spec Schritt 2 nachgezogen). Image: Exit 1 mit `metrics-reason-empty`.

**4. Grund-Codes — bestätigt**
- Spec §4: `SPEC-104`/`105`/`106` in der Tabelle; Konstanten in `metrics.go` (der Plan nennt `model/finding.go`, die Konstanten liegen tatsächlich in `rules/metrics.go`, wie bei den Nachbar-Modulen; die Plan-Tabelle ist an der Stelle ungenau); `AllReasons()` und `reasonTexts()` in `diagnose.go`; `TestAllReasonsDeckungGegenSpezifikationGrundCodes` ist Teil von `make test` und grün. Die Tabellenzeilen stehen vor `SPEC-084` (Reihenfolge nicht aufsteigend) — kosmetisch, kein Gate.
- Nicht verifiziert: Ob `diagnose.go` gofmt-sauber ist. Im Diff sind die drei neuen Map-Einträge anders ausgerichtet als die Nachbarn; `make lint` ist grün, ein eigenes gofmt-Gate existiert nicht, ein Host-`gofmt` ist nach §3.1 nicht erlaubt.

**5. `make gates` grün — bestätigt** (eigener Lauf, Exit 0, siehe oben).

**6. Review durchgeführt — teilweise**
- Report `docs/reviews/2026-10-30-slice-283-r1.md` liegt vor (Befunde F-1 bis F-6, 3 MEDIUM, 3 LOW, kein HIGH). `make review-coverage` ist für diesen Slice erst nach `done/` relevant.
- F-1 (Schema stückweise): Plan geändert, aber Spec/Handbuch-Hinweis offen (siehe Punkt 1). F-2 (Symlink): behoben, Test vorhanden. F-3 (Negativliste): behoben, Mutation 2 zeigt es. F-4 (Hinweis "Marke ausstehend"): per Plan-Änderung in §1 geheilt. F-5 (Herkunfts-Prosa im Test): im HEAD nicht mehr auffindbar (`grep "Issue"` leer). F-6 (Grenze in Nutzertexten): Handbuch-Zeile nennt jetzt alle drei Teile; die README-Zeilen und das Gerüst nennen nur einen Teil bzw. keinen (Gerüst nennt `scan.roots` und `scan.ignore`, nicht die Skip-Liste).
- Offen: Die Rückführungs-Bedingung aus Plan-§4 (MR-066: Grund und Ersatz-Form in zwei Durchgängen, wenn der Slice nicht zurückgeführt wird) ist im Report ausdrücklich "nicht Gegenstand dieses Durchgangs" (Diff 15 Dateien, vier Schichten). Weder Report noch Plan tragen sie bisher. Vor Closure ist sie nachzutragen oder der Slice zurückzuführen.
- Kein Re-Review der Fix-Commits; die Bestätigung der Behebungen ruht auf dieser Verifikation.

**7. Doku-Update der Spiegel (nur Aufzählung) — bestätigt, mit zwei Abweichungen**
- Nur die Aufzählung wurde geändert; `CHANGELOG.md`, Digest-Pins und Handbuch-Kopf sind unberührt.
- Abweichung A (nicht bestätigt): Plan-§3 verlangt für `spec/lastenheft.md` und `spec/spezifikation.md` "Fassungsstand und Historie" (§1: "außer dem Fassungsstand der Historie"). Im Diff fehlt beides: `**Version:** 0.110.0` unverändert, kein neuer Historie-Eintrag, obwohl DC-FA-CLI-002 (Modulsatz, Glossar) im Lastenheft geändert und in der Spezifikation drei Grund-Code-Zeilen und die Symlink-Regel ergänzt wurden. Das ist die einzige Plan-Zeile, die ohne Beleg geblieben ist.
- Abweichung B (Grenzfall): Die Handbuch-§6-Zeile trägt mehr als die Aufzählung (Beschreibung samt Scan-Grenze). Das ist der Spiegel-Ort, der Test hält ihn; die Grenze zu "Prosa über die Aufzählung hinaus" ist unscharf, hier als vertretbar gewertet (Handbuch ohne Slice-Verweis, Ist-Zustand).

**8. Grund-Codes dieses Slice** — siehe Punkt 4, bestätigt.

**9. Closure-Notiz, Register, Risiko-Ausgänge, Paarungen** — offen (erwartungsgemäß: §7 leer, §6-Risiken tragen "wird bei der Closure gesetzt", Register nicht fortgeschrieben, DoD-Haken alle leer).

## Abgrenzung (§1) gegen den Diff

- Messung aus git: nicht mitgenommen (`metrics.go` importiert nur `fmt`, `math`, `path`, `sort`, `strings`, `model`, `driven`; kein VCS-Port; Image-Läufe mit `--network none` und ohne `.git`).
- Tabelle/`estimate-inputs`/`table.*`/`excluded`: nicht mitgenommen; der Config-Rand lehnt sie ab (Image belegt).
- Eintrag im eigenen `.d-check.yml`, `FOCUS_DISABLE`, `gate_consistency_test.go`, `Makefile`, `harness/`: unverändert.
- Release-Prep außerhalb der Spiegel: nicht berührt.
- Mitnahme mit Plan-Änderung vor dem Fix (konform): Hinweis "Marke ausstehend" (F-4), Symlink-Regel (F-2, Spec nachgezogen).
- Plan-Änderungen `01fd0e4e` und `5deb73b2` stehen vor den Code-Änderungen, die sie betreffen (Reihenfolge gewahrt).

## Gesamturteil

Die technische Zusage des Slice (Formel, Schwellenregel, leerer Grund, Eingabefehler, Config-Rand, Grund-Codes) ist bestätigt: Gates grün, Image-Verhalten wie erwartet, zwei Mutationen aus dem richtigen Grund rot. Nicht closure-reif sind:

1. Plan-Tabelle §3 Zeile "Fassungsstand und Historie" von `spec/lastenheft.md` und `spec/spezifikation.md` ohne Beleg (Version 0.110.0 unverändert, kein Historie-Eintrag).
2. DoD-Text §2 ("Faktor ohne Eingabe") widerspricht der Plan-Änderung in §1; angleichen.
3. MR-066-Aussage (Grund und Ersatz-Form der Prüfung) zur nicht gezogenen Rückführung fehlt im Review-Report.
4. Nutzerseitige Lücke "Spec-Schema > Config-Rand" und Einrückungs-Zaun (Risiko 1) brauchen bei der Closure einen Ausgang.
5. Byte-Vergleich gegen das Vorher-Image (§5) ist noch zu fahren.
