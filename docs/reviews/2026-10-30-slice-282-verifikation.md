# Verifikations-Report — slice-282 (Spec und ADR für das Modul `metrics`)

**Rolle:** Verifier (Modul 8/11), Frage „Bauen wir es richtig?" gegen Plan und DoD
**Gegenstand:** slice-282 in `in-progress/`, Stand `HEAD` = `674e4234` (Diff `fccf31bd..HEAD`); Plan `slice-282-metrics-anforderung-und-adr.md`, Welle welle-93
**Modell-ID:** claude-sonnet-5-5
**Datum:** 2026-10-30
**Eingang:** DoD-Bestätigung des Implementers, Review-Report `2026-10-30-slice-282-r1.md` (F-1 bis F-12)

## Selbst gefahrene Sensoren

| Sensor | Ergebnis |
|---|---|
| `make gates` | Exit 0; Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; `doc-check` 0 Befunde |
| `make adr-check` | 0 Befunde (Range `HEAD~1..HEAD`, 1 Commit) |
| `make gate-consistency` | 0 Befunde |
| `make planning-check` | 0 Befunde |
| `make review-coverage` | 0 Befunde |
| Formel von Hand | siehe unten |

`make verify-closure-notes` ist nicht gefahren: der Slice liegt in `in-progress/`, die Closure-Notiz (§7) ist erwartungsgemäß leer.

## Nachrechnung der Formel (ganzzahlig, `100 · ⌈135·S / 10000⌉`, Schwelle 300000 für `135·S`)

| S | 135·S | Marke | Soll laut Lastenheft | Schwelle überschritten |
|---|---|---|---|---|
| 6790 | 916650 | 9200 | 9200 | ja |
| 5036 | 679860 | 6800 | 6800 | ja |
| 5430 | 733050 | 7400 | 7400 | ja |
| 9215 | 1244025 | 12500 | 12500 | ja |
| 4670 | 630450 | 6400 | 6400 | ja |
| 2000 | 270000 | 2700 | 2700 | nein |
| 2001 | 270135 | 2800 | 2800 | nein |
| 2222 | 299970 | 3000 | 3000 | nein |
| 2223 | 300105 | 3100 | 3100 | ja |

Alle neun Paare stimmen (Shell-Ganzzahlrechnung, Aufrundung als `(n+9999)/10000`). Das Lastenheft-Kriterium „Formel" nennt 300105 (2223) und 299970 (2222) richtig und nennt jetzt die Konfiguration (Zähler 135, Nenner 10000, Einheit 100, Schwelle 300000).

## Die neun Entscheide: Plan §1 gegen Lastenheft, Spezifikation, ADR

| # | Entscheid | Lastenheft | Spezifikation | ADR | Urteil |
|---|---|---|---|---|---|
| 1 | Zeilen-Syntax, kein YAML, Plan-Suche `slice-*.md` unmittelbar | Absatz „Der Block", Beschreibung | Schritte 2, 3, 4 | Entscheid 1 | deckungsgleich |
| 2 | Basis als Block-Feld, `HEAD` fest, Eingabefehler bei Nicht-Vorfahr | „Der Block", Eingabefehler-Liste | Schritt 4 (7 bis 40 `0-9a-f`), Schritt 9 | Entscheid 2 | deckungsgleich |
| 3 | Zwei-Stufen-Formel, Schwellenregel, Zahlen = Konfiguration | „Die Formel…", Kriterium „Formel" | Schritt 7, Schema `metrics.surcharge.*`, `rounding`, `cut-threshold` | Entscheid 3 | deckungsgleich |
| 4 | `--no-renames`, literale Pfade, Binär 0, Teilbaum, nur Eingechecktes | „Die Messung" | Schritt 9 | Entscheid 4 | deckungsgleich; das Plan-§1 erwähnt den Hinweis „Pfad ohne Treffer" nicht (siehe Befund V-2) |
| 5 | Eingabefehler Exit 2, fünf Befund-Codes, Hinweise | „Die fünf Prüfungen", Eingabefehler-Absatz | Schritte 8, 10, 12 | Entscheid 5 | deckungsgleich |
| 6 | Lesen im Kern, kein Parse-Port | „kein YAML" | Schritte 3, 4 | Entscheid 6, Alternative B | deckungsgleich |
| 7 | Fence-Ausnahme begrenzt auf `metrics`, ungeschlossen = Eingabefehler | Eingabefehler-Liste | Schritt 3 | Entscheid 7 | deckungsgleich |
| 8 | Zeilen-Summe am VCS-Port | §3.8-Absatz, „Die Messung" | Schritt 9 | Entscheid 8, Alternative C verworfen | deckungsgleich; `spec/architecture.md` ARC-005, ARC-010 und Fehlerzeile tragen `metrics` und die Zeilen-Zählung (Dateistand geprüft, Zeilen 83, 99, 153) |
| 9 | Tabelle im Modul, Prüfung 2 nur bei Messung 0, „Marke ausstehend" | „Die Tabelle", Zustände | Schritte 5, 6, 7, 8 | Entscheid 9, Alternative E | deckungsgleich |

## DoD-Punkte

| DoD-Punkt (§2) | Urteil | Beleg |
|---|---|---|
| Lastenheft: Anforderung, §3.8-Grenze, Out-of-Scope, Akzeptanzkriterien, Kürzel `MET`, Version, Historie, neun Entscheide gesetzt, ohne Modulsatz und Glossar | **bestätigt** | `spec/lastenheft.md` Z. 4187 bis 4304 (Beschreibung, Akzeptanzkriterien, Out-of-Scope mit fünf Punkten); `MET` in der Kürzel-Liste Z. 80; Version `0.110.0` im Kopf Z. 3 und Historie-Zeile Z. 4361 (`MR-032` erfüllt, keine Slice-/ADR-Kennung im Text); `DC-FA-CLI-002` und Glossar nennen `metrics` nicht (`grep` zeigt `metrics` nur in `DC-FA-MET-001`, Kürzel- und Historie-Zeile), `make test` im Gate-Lauf grün (Spiegel-Test) |
| Spezifikation: Schema `metrics.*`, Block-Syntax, Formel samt Rundung, Messsemantik, Grund-Codes im Verfahren, ohne Tabellenzeilen in §4; Architektur-Schärfung | **bestätigt** | `spec/spezifikation.md` Verfahren `DC-FA-MET-001.a` Z. 3385 bis 3491 (12 Schritte, fünf Codes in Schritt 8); Schema-Zeilen Z. 3801 bis 3811; Historie Z. 4009; keine `metrics-*`-Zeile in der Grund-Code-Tabelle (`grep` nur Z. 3455 bis 3476 im Verfahren); `spec/architecture.md` ARC-005, ARC-010, Z. 153 |
| ADR: Alternativen-Tabelle zu Entscheid 1 bis 9, Re-Evaluierungs-Trigger, `Schärft:`, Index; Status `Proposed` | **bestätigt** | `docs/plan/adr/0112-metrics-modul.md`: Status `Proposed`; Tabelle A bis F (siehe V-1 zur Abdeckung); `## Re-Evaluierungs-Trigger` (a bis c); `Schärft:` mit Abschnitts-Ankern auf ARC-005, ARC-010 und §5; Index-Zeile in `docs/plan/adr/README.md` Z. 122; `make adr-check` und `doc-check` grün |
| Bestands-Gegenprobe dokumentiert; Prüfvorlage als Tabelle in der Closure-Notiz; Herkunft aus Angaben gekennzeichnet | **offen** (teilweise) | Die Kennzeichnung „aus Angaben, nicht am Repo gemessen" steht im Lastenheft (Z. 4272 bis 4278), in der Spezifikations-Historie und in der ADR. Die Tabelle gehört in die Closure-Notiz; §7 trägt „Wird bei der Closure gefüllt". Die Rechnung der neun Paare ist oben nachgeprüft und stimmt |
| `make gates` grün | **bestätigt** | eigener Lauf, Exit 0 (siehe oben) |
| Review durchgeführt, Report unter `docs/reviews/` | **bestätigt** (mit V-3) | `2026-10-30-slice-282-r1.md` liegt vor, `make review-coverage` grün; Verdikt dort „nicht freigabefähig vor Behebung"; eine Re-Review der Einarbeitung liegt nicht vor |
| Doku-Update: Handbuch und `operations.md` nicht | **bestätigt** | Diff `fccf31bd..HEAD` berührt weder `docs/user/` noch README |
| Closure-Notiz mit Lerneintrag | **offen** | §7 leer; erwartet vor dem Move nach `done/` |
| Beobachtungs-Register fortgeschrieben | **offen** | kein Diff unter `docs/plan/planning/observations/`; Closure-Pflicht |
| Jedes Risiko aus §6 trägt einen Ausgang | **offen** | alle drei Risiken tragen „wird bei der Closure gesetzt" |
| Drei Paarungen | **offen** | getragen von der nächsten Closure; per Definition nach dem Move |

Die Plan-Grenzen aus §1 sind eingehalten: kein Produktcode im Diff (13 Dateien, alle Doku/Spec/Planung), kein Modulsatz und keine Glossar-Zeile, keine Grund-Code-Zeilen in der Tabelle von §4, keine Release-Prep-Dateien.

## Einarbeitung der Review-Befunde F-1 bis F-12 am Dateistand

| Befund | Urteil | Beleg am Dateistand |
|---|---|---|
| F-1 Zaun-Grammatik | **eingearbeitet** | Spezifikation Schritt 3 stellt auf den gemeinsamen Öffner-Test um und benennt die Abweichung vom Anlass-Repo; Lastenheft „Der Block" legt „gilt die Leseweise des Produkts" fest. Ob die Leser in der Praxis gleich antworten, zeigt erst der Bau (nicht belegbar in einem Spec-Slice) |
| F-2 Zell-Teiler | **eingearbeitet** | Schritt 6 nutzt den gemeinsamen Zeilen-Teiler (Code-Spans, optionaler Anfangsstrich), Abweichung benannt |
| F-3 Wirkung von Exit 2 | **eingearbeitet** | Lastenheft Eingabefehler-Absatz und Schritt 12: erster Fehler in Plan-Namen-Reihenfolge, „kein Befund und kein Hinweis aus keinem Plan" |
| F-4 Pfad ohne Treffer | **eingearbeitet** | Lastenheft „Die Messung" und §3.8-Absatz („Ein Tippfehler im Pfad liefert still 0"), Schritt 9 mit Hinweis; „Messung 0" ausdrücklich als Stellvertreter |
| F-5 Konfiguration des Kriteriums | **eingearbeitet** | Kriterium „Formel" nennt Zähler, Nenner, Einheit, Schwelle und die Produkte 300105 und 299970 |
| F-6 Schwellengröße | **eingearbeitet** | Lastenheft Formel-Absatz und Prüfung 5 sowie Spezifikation Schritt 7 und 8 sagen einheitlich `Zähler · Schätzung des Blocks` |
| F-7 Zeile des Befunds | **eingearbeitet** | Schritt 10: bei `metrics-cut-reason-missing` die Zeile der Schätzung, bei `metrics-over-limit` die Zeile der Marke |
| F-8 `scan.ignore`/Skip-Liste | **eingearbeitet** | Schritt 2 und Lastenheft §3.8-Absatz nennen beide Achsen |
| F-9 `Schärft:`-Feld | **eingearbeitet** | ADR-Feld verlinkt Abschnitts-Anker (`#2-schichten-und-constraints`, `#3-externe-abhängigkeiten`, `#5-fehlermodelle-und-resilienz`) und nennt die Fehlerzeile; `doc-check` löst die Anker auf |
| F-10 Entscheid-Nummern | **eingearbeitet** | ADR: „die Nummern dienen dieser Liste als Bezug"; Lastenheft und Spezifikation führen keine Entscheid-Nummern (`grep "Entscheid" spec/` ohne Treffer in den neuen Abschnitten) |
| F-11 Protokoll-Eintrag | **eingearbeitet** | Schritt 4: „auch ein Protokoll-Eintrag über frühere Marken, den ein Repo mitführt" wird als unbekannter Schlüssel ignoriert |
| F-12 Singular/Plural | **eingearbeitet** | Lastenheft „die konfigurierten Teilbäume", Schema `metrics.excluded` als Liste; ADR „ausgenommene Teilbäume" |

Alle zwölf Einarbeitungen stehen im Dateistand, nicht nur in der Commit-Botschaft.

## Befunde der Verifikation

- **V-1 (LOW, DoD-Wortlaut).** Die Alternativen-Tabelle der ADR (A bis F) deckt die Entscheide 1 bis 9 als Gesamt-Alternativen ab (B zu 1/6, C zu 8, D zu 3, E zu 9). Zu den Entscheiden 2, 4, 5 und 7 steht keine eigene Alternativenzeile; die Wahl wird nur festgestellt. Der DoD fordert „Alternativen-Tabelle zu den Entscheiden 1 bis 9"; formal erfüllt, inhaltlich dünn bei 2 (Basis aus `--range`), 5 (Befund statt Exit 2 für Eingabefehler) und 7 (Fence). Keine Blockade für `Proposed`; vor `Accepted` zu erwägen.
- **V-2 (LOW, Plan-Drift).** Plan §1 Entscheid 4 und 5 (Liste der Hinweise) kennt den Hinweis „Pfad ohne Treffer" nicht; er kam mit der Einarbeitung von F-4 in Lastenheft, Spezifikation und ADR. Das ist eine Schärfung innerhalb des Ziels, keine Ausweitung der Abgrenzung; für die Closure-Notiz zu vermerken.
- **V-3 (MEDIUM, Prozess).** Der Review-Report schließt mit „nicht freigabefähig vor Behebung der MEDIUM-Findings". Die Behebung ist geprüft (oben), aber ein zweiter Review-Lauf (R2) existiert nicht; die Bestätigung hier ersetzt ihn nicht. Der Verifier prüft die Konformität mit der Zusage, nicht die Wartbarkeit der neuen Formulierungen (etwa ob Schritt 3 und 6 mit dem tatsächlichen Verhalten von `FenceToggle` und `SplitPipeTableLine` übereinstimmen; das ist ein Sachverhalt am Code, den slice-283 mit Tests belegen muss).
- **V-4 (INFO).** Die Zusage „Zeilen-Summe gleicht `git diff --numstat --no-renames`" ist an allen drei Stellen als unbelegt gekennzeichnet (Lastenheft, Spezifikations-Historie, ADR-Folgepflicht); kein Fall von behauptet-statt-gemessen. Ein „Bewusstes Brechen" entfällt, weil der Slice keinen Test und keinen Code liefert.

## Gesamturteil

Die Lieferung von slice-282 entspricht dem Plan und den Liefer-Punkten der DoD (Lastenheft, Spezifikation samt Architektur, ADR, Gates). Offen sind ausschließlich die Closure-Pflichten (Gegenprobe-Tabelle und Lerneintrag in §7, Register, Risiko-Ausgänge, Paarungen) sowie eine empfohlene Re-Review der in R1 beanstandeten Stellen. Kein DoD-Punkt ist ohne Beleg nur behauptet. Kein Befund hindert den Übergang nach `done/`, sobald §7 gefüllt ist.
