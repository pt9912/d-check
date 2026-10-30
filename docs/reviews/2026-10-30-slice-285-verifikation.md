# Verifikations-Bericht — slice-285 (Tabellenzählung des Moduls `metrics`)

**Rolle:** Verifier (Modul 8/11, „Bauen wir es richtig?" gegen Plan und DoD)
**Gegenstand:** Arbeitsbaum auf `bb792e8e` plus gestagter Fix zu den Befunden F-1 bis F-11 aus [`2026-10-30-slice-285-r1.md`](2026-10-30-slice-285-r1.md)
**Eingang:** Plan `docs/plan/planning/in-progress/slice-285-metrics-tabellenzaehlung-und-schaetzung.md` (§1, §2, §3 samt Negativliste), [`welle-93`](../plan/planning/welle-93-metrics-modul.md), [`ADR-0112`](../plan/adr/0112-metrics-modul.md) Entscheid 9, [`DC-FA-MET-001`](../../spec/lastenheft.md#dc-fa-met-001--kennzahl-eines-planungsdokuments-gegen-messung-aus-git-modul-metrics-opt-in) (Lastenheft 0.110.3), [`DC-FA-MET-001.a`](../../spec/spezifikation.md#dc-fa-met-001a--kennzahl-eines-planungsdokuments-gegen-messung-aus-git-metrics) Schritte 5 bis 8, Tabellenformat aus GitHub-Issue 13 (nur gelesen).
**Datum:** 2026-10-30

## Befund vorab: Stand des Gegenstands

Bei Beginn der Verifikation lagen Fix und Review-Report gestagt vor; während des Laufs wurde der Fix als `c8c1a084` („fix(metrics): Review-Befunde R1 …") committet. Verifiziert ist der Inhalt des gestagten Stands, der dem Inhalt von `c8c1a084` entspricht (kein unstaged Diff, Arbeitsbaum nach dem Commit sauber).

## Sensor-Belege (selbst gefahren)

| Sensor | Ergebnis |
|---|---|
| `make gates` | Exit 0; „baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green"; doc-check „55 rules on 73 files: 0 findings"; coverage-gate „94.80 % erfüllt Schwelle 93 %" |
| `make image-test` | Exit 0; Fälle (1) bis (7) „OK", darunter (6) „Messung unter :ro und fremdem Benutzer, Exit 1; ohne .git Exit 2" und (7) „Schätzung aus der Tabelle: Befund bei Messung 0, Hinweis bei Messung >= 1" |
| `make blackbox-probe REF=v0.89.0` | Exit 0; „byte-identisch über 20 Vergleiche (stdout/stderr/Exit getrennt, Kanarienlauf vorher und nachher)" — Läufe ohne `metrics`-Konfiguration unverändert |
| `make test` mit Mutation | siehe Punkt 4 |
| `make verify-closure-notes` | nicht gefahren: der Slice liegt in `in-progress/`, §7 ist „Wird bei der Closure gefüllt." |

## (1) Tabellen-Leser, Konfiguration, Nachrechnung der Issue-Zahlen

Wegwerf-Fixtures (git-Repos unter `/tmp/claude-1000/-Development-d-check/884309e6-148a-4f7a-bcf2-99cddd4d10c5/scratchpad/verif/`, Skript `mk.sh`), Lauf `docker run --rm --network none -v <fx>:/repo:ro d-check:latest --enable metrics`, Kopf mit fünf Zellen und `ä` wie im Issue, Faktoren 85/80/200/1, Plan-Commit ohne Änderung unter `src/` (Messung 0):

| Slice | Zeilen (nebenläufig) | Fakes | Produktcode | Schätzung / Marke | Ergebnis |
|---|---|---|---|---|---|
| A | 44 (2) | 0 | 1506 | 5036 / 6800 | Exit 0, 0 Befunde |
| B | 49 (2) | 0 | 1500 | 5430 / 7400 | Exit 0, 0 Befunde |
| C | 70 (3) | 8 | 2000 | 9215 / 12500 | Exit 0, 0 Befunde |
| B, Block mit 5570 / 7600 (der „Wert am Ende"-Fehler) | wie B | | | | Exit 1, `metrics-estimate-mismatch` „estimate 5570 … (erwartet 5430)" |
| C, Block mit 9017 / 12200 | wie C | | | | Exit 1, `metrics-estimate-mismatch` „erwartet 9215" |

Die Schätzungen 5036, 5430, 9215 und die Marken 6800, 7400, 12500 stimmen mit dem Issue überein; die falschen Schätzungen fallen nur durch Prüfung 2 (Marke ist zu der falschen Schätzung konsistent gewählt).

Fehlerfälle am Image (Exit 2, Klartext, kein Befund desselben Plans), jeweils eine Mutation der Tabelle von Slice A:

| Fall | Meldung (gekürzt) |
|---|---|
| Tabelle fehlt | „tabelle fehlt (kopf …)" |
| zweite Tabelle mit passendem Kopf | „tabelle mehrdeutig: 2 Tabellen …" |
| Trennzeile gelöscht | „trennzeile nach dem tabellenkopf fehlt" |
| Zeile mit vier Zellen | „zeile 2 der tabelle: 4 zellen statt 5" |
| `Nr.` leer | „zeile 3 der tabelle: Nr. leer" |
| `Nr.` doppelt | „zeile 4 der tabelle: Nr. "3" doppelt" |
| Klasse `Andere` | „zeile 2: klasse "Andere" unbekannt" |
| Tabelle ohne Datenzeile | „tabelle ohne zeilen" |
| Tabelle im `~~~`-Zaun / eingerückt | „tabelle fehlt" (Exit 2) |
| CRLF-Plan | Exit 0 |
| Strich im Code-Span einer Zelle | Exit 0 (eine Zelle, wie die geänderte Plan-Zeile und Schritt 6 es sagen) |

Config-Rand am Image (alle Exit 2): `metrics` ohne `table` („… sind Pflicht"), leere Klassenliste, Klasse ohne Faktor, negativer Faktor, Schlüsselspalte nicht im Kopf, `estimate-inputs`-Feld ohne Faktor, doppelte Kopfzelle; ein `estimate-inputs`-Feld, das der Block nicht trägt, endet mit „eingabe-feld "extra" fehlt". Der Rest der Negativliste (Tilden, Einrückung, Zitat, `>`-Ende, Leerraum hinter dem Schlussstrich, Kommentar-Zelle, `nebenlaeufig`, Überlauf) steht in `TestMetricsTabellenGrenzen`, `TestMetricsStrichImCodeSpanIstEineZelle` und `TestMetricsSchaetzungUeberlaufUndNull` (durch `make gates` grün).

**Urteil DoD „Tabellen-Leser und Konfiguration": bestätigt.**

## (2) Prüfung 2 und „Marke ausstehend"

- Messung 0, Schätzung ≠ Tabelle: Exit 1, Befund `metrics-estimate-mismatch` (siehe B, C oben).
- Messung ≥ 1 (Fixture B mit anschließendem Code-Commit in `src/`, Messung 51): Exit 0, nur Hinweis „estimate 5570 weicht von der Tabelle ab (erwartet 5430; Messung 51, nur ein Hinweis)".
- Block ohne Basis, Schätzung, Marke: Slice A Hinweis „Marke ausstehend (estimate 5036, limit 6800)", Slice C „(estimate 9215, limit 12500)", jeweils Exit 0 und 0 Befunde.

**Urteil DoD „Schätzung und Prüfung 2" und „Hinweis Marke ausstehend": bestätigt** (belegte Zahlen der Anforderung am Image und als Testfälle).

## (3) Gates

`make gates` Exit 0, `make image-test` Exit 0, siehe Tabelle oben. **DoD „make gates grün": bestätigt.**

## (4) Bewusstes Brechen (`make test`, danach `git diff` auf die Datei leer)

1. Klassen-Check in `checkMetricsRow` (Zweig „klasse … unbekannt") entfernt: `make test` Exit 2. Rot: `TestMetricsTabellenfehler` (Fall „Klasse unbekannt: erwartet Fehler, got … Hint: estimate 1000 weicht von der Tabelle ab (erwartet 0 …)") und `TestMetricsTabellenGrenzen` (Fall „nebenlaeufig: erwartet Fehler"). Aus dem richtigen Grund: ohne den Check zählt eine unbekannte Klasse still mit Faktor 0 und die Schätzung bricht auf 0.
2. `atColumnZero` auf `return true` gesetzt: `make test` Exit 2, `TestMetricsTabellenGrenzen` rot mit „Kopf eingerückt: erwartet Fehler" (die Map-Reihenfolge zeigt je Lauf den ersten Treffer; „Kopf im Zitat" und „Zitat beendet Tabelle" hängen an derselben Funktion). Richtiger Grund: der eingerückte Kopf wird als Tabelle gelesen.

Beide Dateien sind aus der Sicherung wiederhergestellt; `git diff --quiet -- internal/hexagon/core/rules/metrics_table.go` Exit 0. Der Überlauf-Zweig (F-2) wurde nicht gebrochen; der Test `TestMetricsSchaetzungUeberlaufUndNull` ist im Gate grün, sein Rot-Beleg steht aus.

## (5) Abgrenzung §1

Der Feature-Commit `5629eb5a` berührt elf Dateien, alle in der Plan-Tabelle §3 oder dort durch die Plan-Änderung `6ff7ff3c` benannt (`model/config.go`, `diagnose.go`, `image-test.sh`, `abdeckung-e2e.md`, Spezifikation). Nicht mitgenommen: Prosa-Form, Summen/mehrere Tabellen, `CHANGELOG.md`, README, Handbuch-Kopf (kein Diff seit dem Slice-Start). Ohne `metrics`-Konfiguration byte-gleich zu v0.89.0 (blackbox-probe, 20 Vergleiche).

**Abweichung (nicht bestätigt, LOW bis MEDIUM):** der gestagte Fix ändert `spec/lastenheft.md` (Version 0.110.3, Text der Anforderung: Abweichungen der Leseweise, Historie). Plan §3 führt `spec/lastenheft.md` nicht, und die „Berührten Spec-Stellen" sagen „gelesen, nicht verändert, außer dem Fassungsstand der Historie". Eine Plan-Änderungs-Zeile dazu fehlt (AGENTS.md §6 Schritt 4). Ebenso ist der Fix an `spec/spezifikation.md` Schritt 6 (Spalte 0, Tabelle bei jedem Plan mit Block) inhaltlich größer als „Zeile in Tabelle §4, Tabellenende bei Anfangsstrich, Historie". Empfehlung: Plan-Änderung nachtragen, bevor der Fix committet wird.

## (6) „Pflicht table" gegenüber Gerüst und Spezifikation

- Spezifikation §2 (Schema-Tabelle) und Schritt 1: `table.header`, `.key-column`, `.class-column`, `.classes` Pflicht, sobald ein `metrics`-Block steht; Image: Exit 2 „… sind Pflicht". Das Gerüst (`config_template.go`) kennzeichnet `table:` als „Pflicht" — deckungsgleich.
- `estimate-inputs`: die Schema-Zeile der Spezifikation führt den Schlüssel mit Default „leer" und ohne „Pflicht"; Schritt 5 verlangt aber, dass der Block alle dort genannten Felder trägt. Das Gerüst sagt jetzt „jedes steht in jedem Block (Pflicht)". Verhalten stimmt (Image: fehlendes Feld ⇒ Exit 2), aber das Wort „Pflicht" meint im Gerüst die Block-Felder, in der Schema-Zeile hat es bei den übrigen Schlüsseln die Bedeutung „Config-Schlüssel erforderlich". INFO: die Schema-Zeile der Spezifikation nennt die Block-Pflicht nicht (Review F-10 teilweise offen).
- Die Anforderung gilt nur für ein Repo, das `metrics` konfiguriert; das Modul ist unveröffentlicht, keine Bestandskonfiguration betroffen.

**Urteil: bestätigt**, mit der INFO-Anmerkung.

## DoD-Punkte im Überblick

| DoD-Punkt (§2) | Urteil | Beleg |
|---|---|---|
| Config-Schlüssel und Zustand | bestätigt | Config-Rand am Image (Punkt 1), „Marke ausstehend" mit Zahlen (Punkt 2) |
| Tabellen-Leser und Konfiguration | bestätigt | Fehlerfälle am Image, Mutationen 1 und 2 |
| Schätzung und Prüfung 2 | bestätigt | A, B, C und Fehlvarianten am Image; Messung 0 / ≥ 1 |
| Hinweis „Marke ausstehend" | bestätigt | Image, A und C |
| `make gates` grün | bestätigt | Exit 0, Coverage 94,80 % |
| Review durchgeführt, Report liegt vor | teilweise: Report R1 liegt vor (nur gestagt); die Befunde sind im Fix bearbeitet, ein Folge-Lauf auf den MEDIUM-Befunden F-1 bis F-4 liegt nicht vor | `git status`: `A docs/reviews/2026-10-30-slice-285-r1.md` |
| Doku-Update nur, wo ein Test es hält; Handbuch-Kopf, CHANGELOG nicht | bestätigt | kein Diff an `CHANGELOG.md`, README, Handbuch |
| Grund-Code `metrics-estimate-mismatch` (Tabelle §4, `AllReasons()`, `--doctor`-Klartext) | bestätigt | `gates` grün, darunter `TestAllReasonsDeckungGegenSpezifikationGrundCodes`; Feature-Commit ändert `diagnose.go` und `spezifikation.md` |
| Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | offen | §7 „Wird bei der Closure gefüllt", §6 Ausgänge „wird bei der Closure gesetzt"; Slice liegt in `in-progress/` |

## Offene Punkte für den Planner

1. Der Review-Report R1 war beim Verifikationsstart nur gestagt; vor dem Closure prüfen, dass er in `c8c1a084` oder einem Folge-Commit liegt.
2. Plan-Änderungs-Zeile für `spec/lastenheft.md` und die erweiterte Schritt-6-Fassung nachtragen (Punkt 5).
3. Folge-Lauf des Reviewers auf F-1 bis F-4 (Fix liegt vor, kein R2-Report).
4. Rot-Beleg für den Überlauf-Zweig (F-2) bei Gelegenheit nachtragen; er blieb hier ungebrochen.
5. INFO: Schema-Zeile `metrics.estimate-inputs` in der Spezifikation um die Block-Pflicht ergänzen (Punkt 6).
6. Closure: Risiko 1 (Zaun-/Zeilen-Erkennung weicht vom Markdown-Leser ab) ist nicht eingetreten, soweit hier gemessen; Risiko 3 (Paritäts-Probe vergleicht nur Inhalt) bleibt bis zur Abnahme des Anlass-Repos offen — die Ausgänge setzt der Planner.
