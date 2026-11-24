# Review-Report — slice-293 (Code-Review-Zuschnitt auf Spec-Text), Runde 1

**Review-Art:** Code (Diff gegen Plan, Folge-ADR, Architect-Verdikt und Hard Rules; der Diff ist reiner Spec- und Plan-Text, kein Go-Code)
**Gegenstand:** slice-293, Diff `686c9dbd..11bc17d8` (lokal auf `main`): `78b0131e` (Plan-Nachtrag vor dem Inhalt), `835907f5` (`spec/lastenheft.md` 0.115.0, `spec/spezifikation.md`), `11bc17d8` (Plan-Nachtrag nach dem Inhalt)
**Skill:** `reviewer.md` @ 1.20.0
**Modell-ID:** claude-sonnet-5-5
**Datum:** 2026-11-24
**Eingangs-Kontext:**
- Plan slice-293 (§1 Abgrenzung, §2 DoD, §3 Entscheide und Kriterien, §6 Risiken, beide Plan-Nachträge)
- [`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) (`Proposed`), das [Architect-Verdikt zu welle-97](2026-11-24-welle-97-architect-verdikt.md) (F-1 bis F-6), [`ADR-0114`](../plan/adr/0114-linelength-modul.md), [`ADR-0054`](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
- [`DC-FA-LLEN-001`](../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in), Lastenheft [§5](../../spec/lastenheft.md#5-globale-out-of-scope-punkte), [`MR-032`](../../harness/conventions.md#mr-032)
- [`AGENTS.md`](../../AGENTS.md) §3 (besonders §3.4 und §3.8) und §5 Regel 13, 14, 15, 18
- Code als Gegenstand: `internal/hexagon/core/rules/markdown.go`, `linelength.go`, `lexikon_kopplung_test.go`, die Konsumenten von `tableRowLine`
- Vorherige Findings am Modul: keine Review-Reports zu `linelength` nach dem Modul-Review; die Beobachtungen `BEO-ALL/grenz-aussage-ohne-ihre-bedingung` und `BEO-ALL/neue-anforderung-gegen-globalen-ausschluss-ungeprueft` wurden als Prüfanker gelesen.

**Messbasis dieses Laufs:** Das Image `d-check:latest` (gebaut 2026-11-24 09:08:58; der Code von `linelength` und der Lexik ist im Diff unverändert) und eine Fixture mit `max-chars: 10`, `exempt-unbreakable: false` (damit nur (b) antwortet), aufgerufen als `docker run --rm -v "$FIXTURE":/repo:ro d-check:latest -enable linelength`. Zusätzlich ein Lauf von `--trace` über ein Lastenheft mit eingerückter Tabelle (F-1). `make doc-check` auf `HEAD` vor diesem Report: 0 Befunde.

## Übersicht

| ID | Kategorie | Kurz | Aktion |
|---|---|---|---|
| F-1 | MEDIUM | Lastenheft sagt „die übrigen Module lesen eine eingerückte Zeile weiter nicht als Tabellenzeile“; `trace` und das Register von `planning` lesen sie | einarbeiten |
| F-2 | LOW | Plan-Nachtrag zählt „zwei Ausformungen“ über vier Punkten, die Botschaft nennt „drei“ | ins Register |
| F-3 | INFO | Die Schema-Zeile zu `exempt-tables` führt eine kürzere Grenzen-Aufzählung als Lastenheft und Spezifikation | ins Register |

## Findings

### F-1 — MEDIUM — Die Aussage über „die übrigen Module“ ist weiter gefasst, als der Gegenstand sie trägt

- **kategorie:** MEDIUM · **quelle:** [`AGENTS.md`](../../AGENTS.md) §5 Regel 13 (eine aufgeschriebene Grenze am Gegenstand prüfen); Eingangsfrage 2 dieses Laufs (`BEO-ALL/grenz-aussage-ohne-ihre-bedingung`) · **klasse:** modul-aussage-ueber-gegenstand-hinaus · **verifizierbar:** ja (Lauf unten; ein Test mit eingerückter Tabelle gegen `--trace`)
- **pfad:** `spec/lastenheft.md` · „die übrigen Module lesen eine eingerückte Zeile weiter nicht als Tabellenzeile“ (Absatz „Was Fenced-Code und was eine Tabellenzeile ist …“, kommt wörtlich ein zweites Mal in der Historie-Zeile 0.115.0 vor); dazu im selben Absatz „Für die Tabellenzeile fragt `linelength` etwas anderes als die übrigen Module — nicht „deklariert diese Zeile etwas in einer Tabelle ab Spalte 0““
- **befund:** Das Architect-Verdikt (F-1) und [`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) tragen die Aussage nur für **drei** Konsumenten von `tableRowLine` (`targets`, `planning.waves`, `structure`); das Lastenheft sagt sie für **alle** übrigen Module. Am Gegenstand gelesen stimmt das nicht: `declaredObservationIDs` in `planning_observations.go` trimmt die Zeile (`strings.TrimSpace`) vor dem Test auf `|` und liest eine eingerückte Registerzeile als Tabellenzeile; der RTM-Leser (`SplitPipeTableLine` trimmt) liest sie ebenso. Gemessen: ein Lastenheft mit `   | ID | Anforderung |`, `   |---|---|`, `   | SFC-QA-01 | foo |` liefert unter `--trace` die Zeile `SFC-QA-01 | foo | — | — | WAISE`.
- **Failure-Szenario:** Ein Adopter oder der Auftraggeber liest das höherrangige Lastenheft als produktweite Zusage („außerhalb von `linelength` ist eine eingerückte `|`-Zeile nie eine Tabelle“) und schreibt eine Tabelle in einen Listenpunkt, die `--trace` oder das Beobachtungs-Register dann als Anforderungs- bzw. Register-Zeilen liest. Das Kriterium „andere Module unverändert“ und der Kopplungs-Test in `slice-294` fangen es nicht, weil beide nur die drei Module nennen. Die Spezifikation ist an dieser Stelle richtig (sie nennt die drei Module); das Lastenheft steht darüber (Source Precedence) und sagt mehr.

### F-2 — LOW — Zählung der Ausformungen: Plan-Text und Commit-Botschaft sagen Verschiedenes

- **kategorie:** LOW · **quelle:** [`AGENTS.md`](../../AGENTS.md) §5 Regel 15 (eine Botschaft behauptet nicht mehr, als die Arbeit trägt) · **klasse:** zaehlangabe-widerspricht-aufzaehlung · **verifizierbar:** ja (`git show 78b0131e`)
- **pfad:** `docs/plan/planning/done/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md` · „keine Plan-Änderung, zwei Ausformungen“ (vier Spiegelstriche darunter); Botschaft von `78b0131e` · „drei Ausformungen des Wortlauts gegenueber dem Plan“
- **befund:** Derselbe Nachtrag heißt im Plan „zwei Ausformungen“, führt vier Punkte (Kriterium 2, Out-of-Scope (4), die drei Grenzen als Sätze, die Prüfung gegen §5) und heißt in der Botschaft „drei“. Der Inhalt der Punkte ist richtig; die Zahl stimmt an keiner der beiden Stellen mit der Aufzählung.

### F-3 — INFO — Die Schema-Zeile führt eine kürzere Grenzen-Aufzählung

- **kategorie:** INFO · **quelle:** Maintainability (kein Konventions-Anker; daher nicht LOW) · **klasse:** grenzen-aufzaehlung-kuerzer-als-vertrag · **verifizierbar:** ja (Texte nebeneinander lesen)
- **pfad:** `spec/spezifikation.md` · „vier und mehr Leerzeichen, ein Tab und `> \|` davor lassen die Zeile gemessen“ (Schema-Zeile `linelength.exempt-tables`)
- **befund:** Die Schema-Zeile nennt drei Fälle; der Punkt „Tabellenzeile heißt …“ und „Grenzen, am Gegenstand gelesen“ nennen zusätzlich das geschützte Leerzeichen (U+00A0), die Zeile ohne führende Pipe und die Tabellen in Listenpunkten (`10. `, verschachtelt, `- \| a \| b \|`), die den Anlass des Issues bilden. Die Zeile liest sich als vollständige Liste und ist eine Auswahl; ihr Gegenstück vor der Änderung war ebenfalls kürzer, die Lücke ist nicht neu.

## Negativbefunde (geprüft, ohne Befund)

### Schwerpunkt 1 — Widerspruchsfreiheit

Jede Pflichteingabe an den fünf Stellen gelesen (Absatz „Drei Ausnahmen“, Einleitungssatz und Punkt „Tabellenzeile heißt …“, die sechs Kriterien, Out-of-Scope (4), Spezifikation Schritt 4, „Nicht Gegenstand“, „Grenzen“, Schema-Zeile) und gegen das Messbild des Images gehalten (Spalte „Image heute“ ist das Verhalten **vor** `slice-294`; die Erwartung des Lastenhefts ist die Spalte „Soll“):

| Eingabe | Soll (alle Stellen übereinstimmend) | Image heute |
|---|---|---|
| ein, zwei, drei Leerzeichen vor der Pipe | ausgenommen | Befund (ändert sich mit `slice-294`) |
| vier Leerzeichen | gemessen | Befund |
| Tab vor der Pipe | gemessen | Befund |
| U+00A0 vor der Pipe | gemessen | Befund |
| `> ` vor der Pipe | gemessen | Befund |
| `a \| b \| c \| d \| e` (ohne führende Pipe) | gemessen | Befund |
| `- \| a \| b \| c \|` | gemessen | Befund |
| `10. \| a \| b \| c \|` | gemessen | Befund |
| Prosa-Zeile, drei Leerzeichen, Pipe, keine Tabelle | ausgenommen | Befund (ändert sich) |
| Prosa-Zeile in Spalte 0 | ausgenommen | kein Befund |
| eingerückte Pipe-Zeile im Zaun, `exempt-code-blocks: false` | gemessen | Befund |
| `exempt-tables: false`, jede Pipe-Zeile | gemessen | Befund |

Keine Stelle sagt für eine dieser Eingaben etwas anderes als eine andere. Die Aussage „eine eingerückte Tabellenzeile wird gemessen“ und „keine zweite Antwort“ für die Tabellenzeile steht nirgends mehr (zweite Suche aus slice-293 §3 am `HEAD` wiederholt: jeder Treffer ist umgeschrieben, unberührt und begründet oder eine Historie-Zeile früherer Fassungen). Der Wortlaut „die rohe Zeile beginnt mit `|`, davor stehen null bis drei Leerzeichen“ liest sich im ersten Halbsatz wörtlich (beginnt mit der Pipe), wird aber im zweiten aufgelöst; die Spezifikation sagt „beginnt nach null bis drei Leerzeichen“. Kein Widerspruch, nur zwei Formulierungen.

### Schwerpunkt 2 — Grenz-Aussagen mit Bedingung

Gesucht über den Inhalt, nicht über die zitierte Stelle: `git grep -n -E 'eingerückte[rn]? (Tabellenzeile|`\|`)|eingerückte oder zitierte|dieselbe Antwort|keine zweite|zweite Antwort|mit `\\?\|` beginn|Spalte 0|Spalte-0' -- spec/lastenheft.md spec/spezifikation.md`, dazu `git grep -n -i -E 'eingerückt|Listenkontext|Prosa-Zeile'` über `spec/` und das Lastenheft außerhalb des Abschnitts. Die vier Grenzen (Prosa-Zeile ohne Tabelle, Marke auf der Listenzeile, Marke von vier Spalten oder Verschachtelung, Zaun bei `exempt-code-blocks: false`) tragen ihre Bedingung in Lastenheft (Punkt, Kriterium, Historie), Spezifikation („Grenzen“) und ADR. Die zwei Stellen mit anderem Umfang sind F-1 und F-3. Unberührt und richtig bleiben die Spalte-0-Aussagen für `targets` (Spezifikation 3040), `planning.waves` (2419), `structure` (2801), `metrics` (Lastenheft 4356, Spezifikation 3555) — sie beschreiben `tableRowLine` bzw. die eigene Leseweise von `metrics` (`atColumnZero`), und beides ändert sich nicht.

### Schwerpunkt 3 — Grenz-Aussagen am Gegenstand

Am Image nachgemessen (`max-chars: 10`, eine zweite Datei mit einer Prosa-Zeile `\| prosa ohne tabelle` in Spalte 0) und mit dem Plan-Nachtrag verglichen:

1. Ein bis drei Leerzeichen, vier Leerzeichen, Tab, U+00A0, `> `, Zeile ohne führende Pipe, `- \|`, `10. \|`: im Default je ein Befund (11 Befunde, darunter die Prosa-Zeile mit drei Leerzeichen) — wie die Tabelle des Nachtrags.
2. Spalte-0-Prosa-Zeile: im Default kein Befund, mit `exempt-tables: false` ein Befund — die Grenze „wie heute in Spalte 0“ stimmt.
3. Eingerückte und unabgerückte `\|`-Zeile im Zaun: im Default kein Befund, mit `exempt-code-blocks: false` je ein Befund — Kriterium „im Zaun“ stimmt.
4. Bezugsmenge: im Default `f.md` 25 Zeilen, mit `exempt-tables: false` 26; die zweite Datei steuert 2 bzw. 3 bei (Gesamtzeile 27 bzw. 29) — der Nachtrag nennt 25 und 26 für `f.md`, das stimmt.
5. Zählung vor dem Editieren: `git grep -c 'exempt-tables'` mit den Ausschlüssen aus §3 auf `686c9dbd` liefert 12 Dateien, 35 Zeilen — wie der Nachtrag.
6. Kopplungs-Test: `tabellenFaelle` führt den Fall „eingerueckt“ mit zwei Leerzeichen, `gueltig: false`, für alle vier Konsumenten — wie der Nachtrag. Ein, drei Leerzeichen sind dort nicht gemessen; der Nachtrag sagt das selbst („entsteht mit `slice-294`“), und `HasPrefix(lines[i], "|")` in `tableRowLine` entscheidet sie ohnehin.
7. „Nicht gemessen: Tabelle unter `10. ` mit vier und mehr Leerzeichen“ — der Nachtrag benennt es als Aussage der ADR über CommonMark, nicht als Messung. Die Einordnung ist ehrlich; die Aussage selbst folgt aus der Einrück-Breite einer Vier-Spalten-Marke.

### Schwerpunkt 4 — „andere Module unverändert“

`grep -n tableRowLine internal/hexagon/core/rules/*.go` ohne Tests: `targets.go:246`, `planning_waves.go:277`, `structure_tablecell.go:43` und `:93`, `structure_tableorder.go:129`, dazu `tableHeaderOrSeparator` (`markdown.go:205`, Konsumenten `planning.waves` und `structure`) und `linelength.go:72`. `tableSeparatorRow` prüft das führende `|` selbst an der rohen Zeile. Die Beschreibung der Spezifikation („Antwort von `tableRowLine`, die die übrigen Module lesen“; „Für `targets`, `planning.waves` und `structure` bleibt eine eingerückte Zeile keine Tabellenzeile“) trifft diese Menge und ist richtig. Das Kriterium „Ausnahme (b), andere Module unverändert“ nennt dieselben drei Module. Die Ausweitung im Lastenheft ist F-1.

### Schwerpunkt 5 — Referenz-Richtung (`AGENTS.md` §3.4)

Hinzugefügte Zeilen von `spec/` gegen `ADR|slice|welle|Hash|issue` gelesen: ein Treffer, die Historie-Zeile 0.115.0 mit „Begründung in begleitender ADR“ und „GitHub-Issue 20“. Das ist die bestehende Bauform (mehrere Stellen im Lastenheft tragen „Begründung in begleitender ADR“), nennt keine ADR-Kennung, keinen Slice, keine Welle, keinen Hash; `make doc-check` meldet keinen `matrix`-Befund. Kein Spec-Satz zeigt abwärts.

### Schwerpunkt 6 — Lastenheft-Version

Der Minor-Schritt 0.114.3 → 0.115.0 passt zu einer Verhaltensänderung eines ausgelieferten Moduls; die Patch-Stufen der Reihe 0.112.x und 0.113.x tragen ausdrücklich „keine Verhaltensänderung“, 0.114.0 (neuer Schlüssel) ist ein Minor. Die Historie trägt genau **eine** Zeile je Version ([`MR-032`](../../harness/conventions.md#mr-032)); die Verweis-Spalte steht auf `—`, wie MR-032 sie bis `Accepted` verlangt. „Fünf Kriterien kommen hinzu“ stimmt (eingerückt an und aus; im Zaun; Bezugsmenge; andere Module; Rest) und „ersetzt“ trifft das Kriterium „Ausnahme (b), Grenzen“ (sein erster Fall entfällt, die übrigen bleiben). In der Spezifikation steht eine Historie-Zeile, in `spec/lastenheft.md` keine zweite Version.

### Schwerpunkt 7 — Globale Ausschlüsse (Lastenheft §5)

Alle sechs Ausschlüsse (Architektur-/Build-/Mathematik-Prüfung, Auto-Fix, Nicht-Markdown-Formate, Rechtschreib-/Stil-/Lint-Prüfung) gegen die Änderung gehalten: keiner trifft sie. Die Ausnahme der Zeilenlänge „mit derselben Markdown-Lexik wie die übrigen Module, die Markdown lesen“ bleibt wahr, weil die Frage der Tabellenzeile weiter in der geteilten Lexik beantwortet wird ([Verdikt](2026-11-24-welle-97-architect-verdikt.md) F-3, [`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) Entscheid 4). Lastenheft §5 ist im Diff unverändert. Diese Lesart ist eine Entscheidung des Architect; würde der Helfer in `linelength` wandern, wäre §5 berührt.

### Schwerpunkt 8 — Botschaften und Reihenfolge

`78b0131e` trägt die Plan-Notiz **vor** dem Inhalt und sagt „Noch kein Eingriff in Lastenheft oder Spezifikation“ — der Diff dieses Commits berührt nur den Plan. `835907f5` nennt, was der Diff trägt (Lastenheft, Spezifikation, kein Code, kein Handbuch, keine ADR-Änderung, §5 unverändert); `11bc17d8` sagt „Kein Eingriff in Lastenheft oder Spezifikation“ und berührt nur den Plan. Beide Botschaften nennen eine Kennung. Einzige Abweichung ist die Zählung in F-2.

### Weitere Bereiche

| Bereich | Ergebnis |
|---|---|
| Abgrenzung slice-293 §1 | geprüft, ohne Befund: kein Code, kein Test, kein Handbuch, kein CHANGELOG, kein ADR-Text, kein Zitat-Fall, kein Schalter; die Listenkontext-Erkennung steht als Ausschluss, nicht als Leistung |
| Entscheide 1 bis 4 aus [`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) | geprüft, ohne Befund: jeder trägt in Lastenheft und Spezifikation eine Stelle; die Zuordnung im Nachtrag (Entscheid 1 zu Kriterium 1 und 4, 2 zu 2, 3 zu 1 bis 3, 4 zu 5 und 6) stimmt mit den Kriterien überein |
| Verdikt F-3 und F-4 | geprüft, ohne Befund: §5 geprüft und unverändert; die drei Grenzen stehen als Sätze im Punkt „Tabellenzeile heißt …“ mit ihrer Bedingung und ohne eigenes Kriterium |
| Plan-Nachträge, Regel 18 | geprüft, ohne Befund: jede Festlegung nennt ihre Stelle (ADR-0116 Alternative D, Entscheid 3 und 4, Verdikt F-3 und F-4, §3 Kriterium 3) |
| Zeilen-/Spiegel-Zählung (Regel 14) | geprüft, ohne Befund: Form (Kommando) vor Zahl; Zahl neu gezählt und reproduziert |
| Eingaben, die das Modul nicht scannt (`AGENTS.md` §3.8) | geprüft, ohne Befund: der Diff ändert die Eingabemenge nicht; `Außer der Konfiguration liest das Modul keine Datei außerhalb …` steht unverändert |
| Hexagon-Richtung, Netz, Suppression, Kommentare (Prüffragen 3 bis 7) | nicht berührt: kein Code im Diff |
| `make doc-check` | 0 Befunde auf `HEAD` vor diesem Report |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** modul-aussage-ueber-gegenstand-hinaus · zaehlangabe-widerspricht-aufzaehlung · grenzen-aufzaehlung-kuerzer-als-vertrag

## Verdikt

**Merge-blockierend:** ja, wegen F-1 (MEDIUM) — eine Zeile im Lastenheft (Absatz und Historie-Zeile 0.115.0) auf die drei Module einzugrenzen, die die Spezifikation und das Kriterium schon nennen; ohne neues Kriterium und ohne Eingriff in den Code. F-2 und F-3 blockieren nicht (LOW und INFO, ins Register).

**Gesamturteil:** Die Anforderung ist, bis auf F-1, widerspruchsfrei und an jeder geprüften Stelle am Gegenstand gedeckt: alle zwölf Pflichteingaben sagen an allen Stellen dasselbe, jede Grenze trägt ihre Bedingung, die Messungen des Implementers ließen sich am Image wiederholen, die Spezifikation beschreibt die Konsumenten von `tableRowLine` richtig, die Spec-Straten zeigen nicht abwärts, Version und Historie folgen [`MR-032`](../../harness/conventions.md#mr-032), und kein globaler Ausschluss ist getroffen. Das Risiko 1 aus slice-293 §6 (eine Stelle stünde gegen die andere) ist gelesen und entfällt: F-1 ist keine Stelle, die der Änderung widerspricht, sondern ein Satz, der über den Gegenstand hinausgeht.

**Reife für die Verifikation:** nach dem Einarbeiten von F-1 ja. Bis dahin steht im Lastenheft eine Zusage, die `--trace` und das Beobachtungs-Register des Moduls `planning` widerlegen; die Verifikation würde sie als DoD-Behauptung („andere Module unverändert“) gegen den Gegenstand prüfen und dasselbe finden.

**Review-Vorgabe der Welle:** höchstens zwei Runden. Dies ist Runde 1; eingearbeitet werden F-1 (HIGH und MEDIUM), F-2 und F-3 gehen ins Register.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Closure-Notiz von slice-293 §7 und von dort in den Zähler des Beobachtungs-Registers.
