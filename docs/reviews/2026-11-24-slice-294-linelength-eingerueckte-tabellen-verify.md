# Verifikations-Report — slice-294 (Kern und Tests)

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?“, gegen DoD, Plan und Anforderung; nicht gegen den Plan-Stil, nicht gegen den realen Bedarf.
**Gegenstand:** [`slice-294`](../plan/planning/in-progress/slice-294-linelength-eingerueckte-tabellen-kern.md), `git log 482e0dc8..HEAD` auf `main`, lokal (`HEAD` = `7e134483`). Code in `222b5bcb` (Helfer, Aufruf, Tests) und `a4a91511` (Handbuch, `--print-config`-Vorlage, Zaun-Test); `b10dbb98` ist der Plan-Nachtrag, `7e134483` der Review-Report.
**Maßstab:** [`DC-FA-LLEN-001`](../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in) in der Fassung 0.115.1, [`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) (`Proposed`), [`ADR-0054`](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 4, §1/§2/§3/§6 des Slice.
**Eingang:** der Review-Report [`2026-11-24-slice-294-linelength-eingerueckte-tabellen-r1.md`](2026-11-24-slice-294-linelength-eingerueckte-tabellen-r1.md) (kein HIGH, kein MEDIUM); die DoD-Bestätigung des Implementers wurde **nicht übernommen**, jeder Punkt ist unten nachgefahren.
**Datum:** 2026-11-24 · **Modell-ID:** claude-sonnet-5-5

**Messbasis.** Alle Sensoren auf dem Arbeitsbaum von `HEAD` (`git status --short` vor und nach dem Lauf leer). Mutationen liefen in einer Kopie von `git archive HEAD` im Scratchpad (Go in Docker, `golang:1.27.4`, `go test ./internal/hexagon/core/rules -run "Linelength|Tabellenzeilen"`, 72 Teiltests; Ausgangslage grün), nicht im Repo; nach jeder Mutation wurden die zwei Quelldateien aus der unveränderten Kopie zurückgelegt und mit `diff` als gleich bestätigt. Das Image `d-check:latest` stammt aus dem `make build` des Gate-Laufs auf `HEAD` (`--print-config` zeigt den neuen Kommentar). Das Vorher-Image `d-check:v294-vorher` ist aus `git archive 69c53e0f` mit `docker build --target runtime` gebaut.

## Ergebnis auf einen Blick

| # | DoD-/Prüfpunkt | Urteil |
|---|---|---|
| 1 | Helfer und Aufruf | bestätigt |
| 2 | Tests je Kriterium der Familie (b), jeder rot ohne die Regel (19 Mutationen gefahren, 17 mit Wirkung, 2 äquivalent) | bestätigt |
| 3 | Kopplung in beiden Richtungen, andere Module unbewegt | bestätigt |
| 4 | Spiegel: Handbuch §4.28, `--print-config`, Abdeckungs-Datei, Image-Aussagen | bestätigt (zwei INFO) |
| 5 | `make gates`, `make coverage-gate`, `make image-test` | bestätigt, alle grün |
| 6 | Form aus Issue 20 am gebauten Image | bestätigt |
| 7 | „Ohne die Änderung nichts Neues“ — Vorher-/Nachher-Vergleich | bestätigt: 120 Vergleiche, 0 Abweichungen; die gewollte Abweichung ist die einzige |
| 8 | Risiken aus §6 | alle fünf vorgeschlagen *entfallen* |
| 9 | Abgrenzung §1 | bestätigt, nichts Ausgeschlossenes berührt |
| — | Rest für den Architect/Planner (kein Befund) | siehe [Rest](#rest-kein-befund) |

**Gesamturteil: closure-reif — ja**, vorbehaltlich der unten benannten Closure-Schritte (Closure-Notiz mit den Mutationen aus [Punkt 2](#2-tests-je-kriterium-und-mutationen), Register, `Accepted`-Übergang der Folge-ADR samt Vermerk an ADR-0114). Kein Befund blockiert.

## 1. Helfer und Aufruf

- `unbreakableTableRow` steht in `internal/hexagon/core/rules/markdown.go` (Z. 208). Die erste Anweisung ist `if tableRowLine(lines, prose, i) { return true }`; danach zählt sie Leerzeichen (`' '`, also U+0020, kein Tab, kein U+00A0) und verlangt `indent >= 1 && indent <= maxTableIndent && prose[i+1] && strings.HasPrefix(line[indent:], "|")`. `maxTableIndent = 3`.
- `tableRowLine` ist zeichengleich: `git diff 482e0dc8..HEAD -- internal/hexagon/core/rules/markdown.go` zeigt **ausschließlich** hinzugefügte Zeilen (28 `+`, kein `-`, keine Änderung innerhalb von `tableRowLine`).
- `linelength` ruft für (b) nur den Helfer: `internal/hexagon/core/rules/linelength.go` Z. 75, `cfg.TablesExempt() && unbreakableTableRow(lines, prose, i)`; im Diff ist das die einzige Code-Zeile außerhalb des Dateikopf-Kommentars. Die Mutation „`linelength` ruft wieder `tableRowLine`“ (M6) wird rot — der Aufruf ist also tragend, nicht dekorativ.
- Konsumenten von `tableRowLine` (`grep -rn 'tableRowLine(' internal`, ohne Tests): `planning_waves.go` Z. 277, `structure_tableorder.go` Z. 129, `structure_tablecell.go` Z. 43 und 93, `targets.go` Z. 246, `markdown.go` Z. 233 (`tableHeaderOrSeparator`), dazu der Helfer selbst. **Keine dieser Dateien steht im Diff.**

## 2. Tests je Kriterium und Mutationen

Zuordnung **eigenhändig** aus dem Lastenheft (sieben Kriterien der Familie (b) plus die Grenzen des Punkts „Tabellenzeile heißt …“), nicht aus der Tabelle des Implementers. Alle Tests in `internal/hexagon/core/rules/linelength_test.go` (`TestLinelength…`) bzw. `lexikon_kopplung_test.go` (`TestTabellenzeilenFrage…`).

| Kriterium im Lastenheft | Test, der es hält | Mutation, die ihn rot macht (gelesene Meldung) |
|---|---|---|
| (b) an und aus, Spalte 0, ohne Trennzeile | `TestLinelengthAusnahmeTabellenzeile` (Bestand, unverändert) | M5 (Helfer ruft `tableRowLine` nicht): `Befund-Zeilen [1], erwartet []` |
| (b) eingerückt, ein/zwei/drei Leerzeichen, an und aus | `TestLinelengthEingerueckteTabellenzeile` | M6 (`linelength` ruft `tableRowLine`): `Befund-Zeilen [1 2 3], erwartet []`; M2 (Schwelle 2): rot |
| (b) Grenze vier Leerzeichen | `TestLinelengthTabellenzeileDreiGegenVierLeerzeichen` (drei/vier/fünf je allein) | M1 (Schwelle 4): `Befund-Zeilen [], erwartet [1]` an `vier_Leerzeichen`, `drei_Leerzeichen` bleibt grün; M2 (Schwelle 2) macht umgekehrt `drei_Leerzeichen` rot — beide Seiten der Schwelle |
| (b) Grenze Tab (auch `␠⇥`, `␠␠␠⇥`) | `TestLinelengthTabellenzeileGrenzen/Tab`, `/Leerzeichen_dann_Tab`, `/drei_Leerzeichen_dann_Tab` | M7 (Tab zählt als Leerzeichen): `Tab` und `Leerzeichen_dann_Tab` rot |
| (b) Grenze U+00A0 | `TestLinelengthTabellenzeileGrenzen/geschuetztes_Leerzeichen` (Bytes `c2 a0` im Quelltext geprüft) | M8 (U+00A0 → Leerzeichen): rot |
| (b) Grenze `> \|` | `…Grenzen/Zitat`, `/Zitat_mit_Einrueckung` | M9b (Helfer schneidet `> ` ab): rot an beiden Fällen und im Kopplungs-Fall `Zitat` |
| (b) Grenze `a \| b` ohne führende Pipe | `…Grenzen/ohne_fuehrende_Pipe` | M14 (`tableRowLine` nimmt `Contains`): rot (zusätzlich 16 Kopplungs-Teilfälle) |
| (b) Grenze Zeile im Zaun, (a) aus, (b) an | `TestLinelengthTabellenzeileImZaunWirdGemessen/Spalte_0` | M15 (`tableRowLine` ohne `prose`-Prüfung): rot an `Spalte_0` |
| (b) eingerückt im Zaun | `…ImZaunWirdGemessen/drei_Leerzeichen`, `/ein_Leerzeichen`; `TestLinelengthEingerueckteTabellenzeileVorUndHinterZaun` | M3 (Zaun-Prüfung im Helfer entfällt): `docs/b.md: Befund-Zeilen [], erwartet [2]`; `docs/a.md: Befund-Zeilen [], erwartet [3]` |
| (b) Bezugsmenge | `TestLinelengthBezugsmengeEingerueckteTabellenzeile`, `…ZaehltNichtInU` | M13 (`count.Lines++` vor die Ausnahme): beide rot (dazu die Bestandstests der Bezugsmenge) |
| (b) andere Module unverändert | `TestTabellenzeilenFrageHatEineAntwort` (`targets`, `planning.waves`, `structure`, `tableRowLine` auf einer Antwort; Fälle ein/zwei/drei/vier Leerzeichen, Tab) | M4 (`tableRowLine` selbst trimmt): `targets antwortet true, gemeinsam erwartet false (Antworten: map[planning.waves:true structure:true tableRowLine:true targets:true])` an `ein_Leerzeichen`, `eingerueckt`, `drei_Leerzeichen`, `vier_Leerzeichen` |
| (b) Rest unverändert (Dateien ohne eingerückte Zeile) | `TestLinelengthTabellenzeileSpalteNullUnveraendert` (Kopf-, Trenn-, Datenzeile Spalte 0, Zeilen darum, (b) aus) | M5: `Befund-Zeilen [2 3 4 5 6], erwartet [2 6]` |
| Grenze: Prosa-Zeile mit `\|` ohne Tabelle ausgenommen | `TestLinelengthProsaZeileMitPipeOhneTabelleIstAusgenommen` (Spalte 0 und drei Leerzeichen) | G1 (Helfer verlangt Trennzeile): rot |
| Grenze: `- \| a \| b \|` gemessen | `TestLinelengthTabellenzeileHinterListenmarkeWirdGemessen/Bindestrich` | M10b (Helfer schneidet `- ` ab): rot; Kopplungs-Fall `Listenmarke` rot |
| Grenze: Tabelle unter `10. ` gemessen | `…HinterListenmarke…/zweistellig_nummeriert`, `/nummeriert`, `TestLinelengthTabelleUnterNummeriertemListenpunkt` | M10c (Helfer schneidet `10. ` ab): `zweistellig_nummeriert` rot |
| Grenze: verschachtelte Liste (vier und mehr Leerzeichen) | `…DreiGegenVier…/vier_Leerzeichen`, `/fuenf_Leerzeichen` | M1, M4 wie oben |
| Text mit Einrückung ohne `\|` | `…Grenzen/Einrueckung_ohne_Pipe`, `/Pipe_nicht_erstes_Zeichen` | M12 (kein `\|` verlangt): `Einrueckung_ohne_Pipe` rot; M11 (`Contains` statt `HasPrefix`): `Pipe_nicht_erstes_Zeichen` rot |
| Kopplung linelength ↔ Helfer, beide Richtungen | `TestTabellenzeilenFrageVonLinelengthHatEineAntwort` (11 Fälle; Helfer schließt `tableRowLine` ein) | M1, M3, M5, M6, M7, M8, M9b, M10b: je mit der Meldung `linelength antwortet … erwartet …` |

**Siebzehn Mutationen mit Wirkung gefahren, alle rot aus dem behaupteten Grund** (die Meldungen oben sind gelesen, nicht erschlossen): M1 bis M8, M9b, M10b, M10c, M11, M12, M13, M14, M15, G1. Die Mutationstabelle des Plans (§3, zehn Zeilen) ist damit vollständig nachgefahren; M2, M4 und M5 stehen dort, die übrigen sind Zusatz zu „je Seite ein Fall“.

**Zwei Mutationen blieben grün — und sind äquivalent, kein Testloch.** M9 (`> ` vorab abschneiden) und M10 (`- ` vorab abschneiden) ändern allein nichts, weil `indent >= 1` den Pfad abfängt; erst zusammen mit dem Entfernen von `indent >= 1` (M9b, M10b) wirken sie und werden rot. Die Bedingung `indent >= 1` ist damit **gegen jede Mutation immun**: ohne sie antwortet der Helfer für Einrückung null genauso (`tableRowLine` hat dort schon geantwortet). Siehe [V-1](#v-1--info--indent--1-ist-eine-tote-bedingung).

`git status` und `git diff` im Repo nach allen Läufen: leer.

## 3. Kopplung

- `TestTabellenzeilenFrageHatEineAntwort` hält die **erste** Assertion von [`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) Entscheid 4: `targets`, `planning.waves`, `structure` und `tableRowLine` auf einer Antwort, die eingerückte Zeile (ein, zwei, drei, vier Leerzeichen, Tab) eingeschlossen. Der Rot-Beleg ist M4/M14/M15: eine Weitung von `tableRowLine` kippt `targets` zuerst (Meldung oben).
- `TestTabellenzeilenFrageVonLinelengthHatEineAntwort` hält die **zweite**: `linelength` und der Helfer antworten gleich, in beiden Richtungen (bis drei Leerzeichen ausgenommen; vier Leerzeichen, Tab, U+00A0, `>`, Listenmarke, Zaun gemessen), und `tableRowLine` ⇒ Helfer (der Helfer schließt sie ein) ist eigens geprüft. Das ist die Form von [`ADR-0054`](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 4 („je Konsument eine Assertion, an ihm festgenagelt“).
- Der Kopf-Kommentar des Tests ist umgeschrieben (zwei Fragen, ein Hinweis auf [`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md)); der Plan verlangte das.
- `targets`, `planning.waves` und `structure` bewegen sich nicht: kein Konsument im Diff (Punkt 1), und der Vorher-/Nachher-Lauf der Konsumenten ist in Punkt 7 für `linelength` erbracht; die drei anderen Module ändern ihr Ergebnis wegen `tableRowLine` nicht, was M4 von der Gegenseite beweist (ändert man `tableRowLine`, wird genau diese Assertion rot).

## 4. Spiegel

- **Handbuch §4.28.** Die `benutzerhandbuch.md` trägt die Regel an den geplanten Stellen: YAML-Zeile (Z. 1679), Beschreibung von (b) (Z. 1706–1707), der neue Punkt „Eine Tabellenzeile fragt das Modul anders …“ (Z. ab 1714) mit der Grenzliste und der Bezugsmengen-Aussage, die Modul-Tabelle in §6 (Z. 3148). Die Überschrift heißt „Acht Dinge“, und es stehen **acht** Punkte darunter (gezählt: Codepunkte · drei Ausnahmen · Fenced-Code · Tabellenzeile · offener Zaun · kein Zeilen-Marker · Eingabefehler · leere Menge). Die alte Zeile „Fenced-Code und Tabellenzeile kennt das Modul so, wie es die übrigen Module kennen“ ist geteilt; die Tabellenzeilen-Hälfte sagt ausdrücklich, dass `targets`, `planning.waves` und `structure` eine eingerückte Zeile weiter nicht lesen. Eine Suche über den Inhalt (`git grep` auf `eingerückte`, `mit \`|\` beginn`, `exempt-tables` außerhalb von `docs/plan`, `docs/reviews`, `.harness`, `spec`, Tests) findet nur noch die erwarteten Reste: `README.md`/`README.de.md`/`CHANGELOG.md` (Release-Prep, Regel 17 von [`AGENTS.md`](../../AGENTS.md)), Handbuch Z. 1738 (Exit-2-Aussage, unberührt), Z. 3382 (Historie), `tools/image-test.sh` Z. 331 (Test, kein Text).
- **`--print-config`-Vorlage.** `internal/adapter/driving/cli/config_template.go` Z. 314 trägt `# Zeilen mit | (bis drei Leerzeichen davor) nicht messen`; am Image: `docker run … d-check:latest --print-config` zeigt dieselbe Zeile (Z. 306 der Ausgabe).
- **Abdeckungs-Datei.** `docs/user/abdeckung-tests.md` führt neun der zehn neuen Testfunktionen (alle in `linelength_test.go`; `TestLinelengthEingerueckteTabellenzeileZaehltNichtInU` nennt in seinem Kommentar keine `DC-*`-Kennung und steht deshalb nicht in der Ableitung — die übrigen neun tragen die Kennung); `make test` hält sie gegen ihre Ableitung (`TestAbdeckungsDateienFolgenIhrerAbleitung`) und lief grün (Punkt 5).
- **Image-Aussagen.** Mehr als vier Grenz-Aussagen am gebauten Image `d-check:latest` geprüft, Fixture mit `max-chars: 20`, `exempt-unbreakable: false`, 15 Zeilen (Zeilenzahlen unten sind die der Fixture). **Default:** Befund an Zeile 4 (vier Leerzeichen), 5 (Tab), 6 (U+00A0, Bytes `c2 a0` geschrieben), 7 (`> |`), 8 (`a | b`), 9 (`- |`), 10 (`10. |`), 12 (Text mit Einrückung ohne `|`); **kein** Befund an 2 und 3 (ein und drei Leerzeichen), 11 (Prosa-Zeile `   | …` ohne Tabelle) und 14 (drei Leerzeichen im Zaun); Bezugsmenge `9 Zeile(n) geprüft`. **`exempt-tables: false`:** zusätzlich 2, 3, 11, Bezugsmenge `12`. **`exempt-code-blocks: false`:** zusätzlich 14 (die eingerückte Zeile im Zaun), Bezugsmenge `12`. Das Vorher-Image meldet im Default 2, 3 und 11 zusätzlich (Bezugsmenge `12`) — der Unterschied ist genau die Regel.
- **Zwei INFO** am Spiegel: [V-1](#v-1--info--indent--1-ist-eine-tote-bedingung) und [V-2](#v-2--info--eine-handbuch-zeile-ist-länger-als-ihr-absatz).

## 5. Gates, selbst gefahren

| Lauf | Letzte Zeile (echt, aus dem eigenen Log) | Exit |
|---|---|---|
| `make gates` | `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` | 0 |
| `make coverage-gate` | `coverage-gate: OK — Coverage 94.90% erfüllt Schwelle 93%` | 0 |
| `make image-test` | `image-test: OK — Akzeptanzkriterien erfüllt (DC-FA-DIST-001, DC-FA-CLI-007, DC-FA-CLI-008, DC-FA-CLI-013, DC-FA-MET-001, DC-FA-PSYNC-001, DC-FA-LLEN-001)` | 0 |
| `make doc-check` (nach Anlage dieses Reports) | `d-check: 1310 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make verify-closure-notes` | `d-check: 1164 Datei(en) geprüft, 0 Befund(e)` | 0 — nicht aussagekräftig für slice-294: der Slice liegt noch in `in-progress/`, §7 ist leer; der Sensor fährt am Closure-Move |

Die `doc-check`-Stufe von `make gates` lief vor Anlage dieses Reports; der Lauf danach steht in der Tabelle.

## 6. Die Form aus Issue 20 am gebauten Image

Fixture: `1. Schritt:`, Leerzeile, Tabelle (Kopf, Trennzeile, eine überlange Datenzeile mit 43 Zeichen) mit **drei** Leerzeichen davor, Leerzeile, `2. Weiter`; `max-chars: 30`.

| Fall | Ergebnis `d-check:latest` |
|---|---|
| drei Leerzeichen, Default | Exit 0, kein Befund, `4 Zeile(n) geprüft` |
| drei Leerzeichen, `exempt-tables: false` | Exit 1, ein Befund an Zeile 5 (`Zeile hat 43 Zeichen, erlaubt sind 30`), `7 Zeile(n) geprüft` |
| **vier** Leerzeichen, Default | Exit 1, ein Befund an Zeile 5 (44 Zeichen) — bleibt Befund |
| vier Leerzeichen, `exempt-tables: false` | Exit 1, derselbe Befund |
| drei Leerzeichen, Default, Vorher-Image `69c53e0f` | Exit 1, Befund an Zeile 5 — der Anlass ist reproduziert, die Änderung behebt ihn |

## 7. „Ohne die Änderung nichts Neues“ — Vorher-/Nachher-Vergleich

(Einlösung von welle-97 §3 (a), auf Anforderung des Auftrags hier als einzelner Beleg geführt.)

**Aufbau.** Vorher: `d-check:v294-vorher` aus `69c53e0f` (Self-Close von welle-96). Nachher: `d-check:latest` auf `HEAD`. Jeweils `docker run --rm --network none -v <Fixture>:/repo:ro <Image> --enable linelength [--json|--yaml]`, **stdout, stderr und Exit getrennt** mit `cmp`.

**Fixtures ohne eingerückte `|`-Zeile (1 bis 3 Leerzeichen):**

- f1: die Fixture des `image-test` (Zeile 40/41, Umlaut-/Emoji-Zeilen, Zaun, Tabellenzeile Spalte 0, URL-Zeile);
- f2: die Tabellen-Ränder — Kopf-, Trenn-, Datenzeile Spalte 0, vier Leerzeichen, Tab, U+00A0, `> |`, `a | b`, `- |`, `10. |`, `~~~`- und ```` ``` ````-Zaun mit Tabellenzeile darin;
- f3: Prosa — CRLF-Überschrift, Listenpunkt, Zitat, `d-check:ignore`-Marker, lange URL, eingerückter Code-Block ohne Zaun, letzte Zeile ohne `\n`;
- f4: eine Kopie von `docs/user/` (fünf Dateien, 4425 Zeilen; `grep -nE '^ {1,3}\|' docs/user/*.md` findet keine eingerückte Tabellenzeile).

**Matrix.** je Fixture `max-chars` 40 und 120; je Grenze fünf Schalter-Stellungen (Default, `exempt-code-blocks: false`, `exempt-tables: false`, `exempt-unbreakable: false`, alle drei aus); je Kombination drei Ausgabeformen (Text, `--json`, `--yaml`).

**Ergebnis: 4 × 2 × 5 × 3 = 120 Vergleiche (je stdout, stderr, Exit — 360 Einzelvergleiche), 0 Abweichungen.** Zur Probe, dass die Läufe nicht leer vergleichen: mit allen Ausnahmen aus und `max-chars: 40` meldet das Nachher-Image f1 4, f2 12, f3 7 und f4 3849 Befunde, identisch zum Vorher-Image.

**Die gewollte Abweichung** (eine Kopie von f2 plus eine angehängte Zeile `  | aaaa … nnnn |` mit zwei Leerzeichen, `max-chars: 40`): Exit beide 1; stdout unterscheidet sich **nur** in der Zeile `probe.md:24 … Zeile hat 75 Zeichen`, die im Nachher-Lauf fehlt; stderr nur in `14 → 13 Zeile(n) geprüft` und `8 → 7 Befund(e)`. Das ist die einzige beobachtete Abweichung: 1 gewollte Abweichung in einem Zusatzlauf, 0 ungewollte in 120.

## 8. Risiken aus §6 — vorgeschlagene Ausgänge

| Risiko (§6) | Ausgang | Beleg |
|---|---|---|
| Helfer in der geteilten Lexik verschiebt andere Module | **entfallen** | Der Diff von `markdown.go` ist rein additiv (28 `+`, kein `-`); keine Datei der Konsumenten von `tableRowLine` steht im Diff; der Kopplungs-Test trägt beide Richtungen (M4, M14, M15 rot in der ersten, M1/M3/M5–M10b in der zweiten Assertion) |
| Ein Test läuft grün, ohne die Regel zu fahren | **entfallen** | Jeder Test der Tabelle in Punkt 2 hat eine Mutation, unter der er aus dem benannten Grund rot wird (Meldungen gelesen). Die Closure-Notiz muss die Mutationen und Meldungen noch nennen — das ist ihre Bedingung aus §6; der Beleg liegt hier. Zwei grüne Mutanten (M9, M10) sind äquivalent, siehe V-1 |
| Eine Randbedingung der Anforderung hat keinen Test | **entfallen** | Zu jeder Grenze steht ein Test, der sie umdreht: Grenze 3/4 (M1, M2), Grenze „Tab“ (M7), U+00A0 (M8), Zaun (M3, M15), Bezugsmenge (M13), Prosa ohne Tabelle (G1), `- \|`/`10. ` (M10b, M10c), `> \|` (M9b), `a \| b` (M14). Kein Eintrag im Register nötig |
| Die Negativfälle fehlen, weil die Regel zwei Seiten hat | **entfallen** | Einrückungs-Seite: M1, M2, M7, M8; Erstes-Zeichen-Seite: M9b, M10b, M10c, M11, M12 — je Seite ein eigener Fall rot, und die Schwelle ist von beiden Seiten gefahren (M1 von oben, M2 von unten) |
| Die Bezugsmenge zählt anders | **entfallen** | `make gates` grün; Kriterium 5 hält die Zahl (M13 rot); am Image: Default `9` gegen `12` mit `exempt-tables: false`, Vorher-Image `12`; die Abweichung ist gewollt und in Punkt 7 als eigene Zeile geführt (`14 → 13`) |

## 9. Abgrenzung §1

`git diff 482e0dc8..HEAD --stat` nennt neun Dateien: der Slice-Plan (Nachtrag), der Review-Report, `docs/user/abdeckung-tests.md`, `docs/user/benutzerhandbuch.md`, `config_template.go`, `lexikon_kopplung_test.go`, `linelength.go`, `linelength_test.go`, `markdown.go`. **Nicht im Diff:** [`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) (bleibt `Proposed`, Index unberührt), [`ADR-0114`](../plan/adr/0114-linelength-modul.md), `.d-check.yml`, `README.md`, `README.de.md`, `CHANGELOG.md`, `tools/image-test.sh`, `harness/`, `configyaml.go` (kein neuer Schlüssel), `tableRowLine` und seine Konsumenten. Kein Mitnehmen, keine stille Plan-Änderung der Abgrenzung.

## Plan-vs-Code-Diff

| Plan (§3) | Code | Urteil |
|---|---|---|
| ein neuer benannter Helfer neben `tableRowLine`, ruft sie, erweitert um null bis drei U+0020, Zaun verlangt | `unbreakableTableRow`, `maxTableIndent = 3` | wie geplant; Name vom Implementer vergeben, wie vorgesehen |
| `tableRowLine` zeichengleich | nur `+`-Zeilen im Diff | wie geplant |
| `linelength` ruft nur den Helfer; Dateikopf nennt die Grenze | Z. 75; Dateikopf nennt `ab vier Leerzeichen, mit Tab, U+00A0 oder > davor wird gemessen` | wie geplant (siehe Review F-1 zur Begründung im Kommentar) |
| Tests Kriterien 1–4, 6, Grenzen 1–3 | zehn neue Testfunktionen in `linelength_test.go` | wie geplant, vollständiger als verlangt (Zitat mit Einrückung, `␠⇥`, fünf Leerzeichen, Listenmarken) |
| Kopplung: zwei Assertionen | zwei Tests (Plan-Nachtrag Punkt 3 nennt die Aufteilung) | wie geplant |
| Handbuch §4.28 samt Modul-Tabelle, Vorlage, Abdeckung | geändert; Aufzählung „Acht Dinge“ statt „Sieben“ (Plan-Nachtrag Punkt 1) | Abweichung benannt, ohne Wirkung auf die Abgrenzung |
| `Accepted`-Übergang der ADR, ADR-0114-Vermerk | **noch nicht** (Plan: erst nach Review und Verifikation) | Rest, kein Befund |

Der Plan-Nachtrag `b10dbb98` steht **nach** den Code-Commits. Er nennt vier Abweichungen (Zählung, Name, Aufteilung der Kopplungs-Tests, Abdeckungs-Datei im ersten Code-Commit); keine berührt §1, daher greift die Regel „Plan-Änderung vor den Code“ ([`AGENTS.md`](../../AGENTS.md) §6 Schritt 4) nicht — der Nachtrag ist für Zählungen und Namen die richtige Form.

## Befunde

### V-1 — INFO — `indent >= 1` ist eine tote Bedingung

- **pfad:** `internal/hexagon/core/rules/markdown.go` · Z. 208 ff., `return indent >= 1 && indent <= maxTableIndent && …`
- **befund:** Für Einrückung null hat `tableRowLine` im ersten Zweig schon geantwortet; ein zweites „ja“ mit `indent == 0` ist ausgeschlossen, weil es dasselbe `prose[i+1] && HasPrefix(line, "|")` verlangte. Die Bedingung ändert kein Ergebnis und ist gegen Mutation immun (M9 und M10 blieben deshalb grün, siehe Punkt 2). Das ist kein Defekt, aber ein Leser kann sie für tragend halten.
- **Aktion:** ins Register, keine Nachbesserung in dieser Runde (Welle-Vorgabe: nur HIGH und MEDIUM).

### V-2 — INFO — eine Handbuch-Zeile ist länger als ihr Absatz

- **pfad:** `docs/user/benutzerhandbuch.md` · Z. 1707 hat 133 Zeichen, die Nachbarzeilen höchstens 96 (Umbruch beim Einfügen des Satzteils „nach null bis drei Leerzeichen (U+0020)“ nicht nachgezogen).
- **befund:** rein typografisch; das Repo schaltet `linelength` für sich nicht ein ([`ADR-0114`](../plan/adr/0114-linelength-modul.md), Konsequenzen), kein Gate meldet es.
- **Aktion:** ins Register (oder bei der Release-Prep mitnehmen).

### Zum Review-Report R1

F-1 (LOW) und F-2 (INFO) des Reviews habe ich am Quelltext **bestätigt gelesen**: der Kommentar am Helfer führt „Tab, U+00A0 oder vier und mehr Leerzeichen“ auf den fehlenden Listenkontext zurück, den nur die vier-Leerzeichen-Grenze trägt; `tableRowLine` trägt „die EINE Antwort“ ohne Zeiger auf die Frage daneben. Beide sind **nicht** Befunde dieser Verifikation und werden nicht doppelt gezählt; sie gehören wie dort vermerkt ins Register.

## Rest (kein Befund)

Das sind die in der DoD offenen Punkte, die erst **nach** dieser Verifikation fällig sind — sie verhindern die Closure-Reife nicht, sie sind ihr Programm:

1. **Architect:** [`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) `Proposed` → `Accepted`, im selben Commit eine Zeile in `## Geschichte` von [`ADR-0114`](../plan/adr/0114-linelength-modul.md) („Entscheidung 5 (b) teil-superseded durch ADR-0116“) und im ADR-Index die Zeile von ADR-0114 mit „Accepted (teil-superseded: ADR-0116)“; die Status-Zeile von ADR-0114 bleibt `Accepted` ([Verdikt](2026-11-24-welle-97-architect-verdikt.md) F-6).
2. **Planner:** Closure-Notiz (§7) mit Lerneintrag; sie nennt je Test die Mutation und die gelesene Meldung (Tabelle in Punkt 2 trägt sie), die Ausgänge der fünf Risiken (Punkt 8), und die Register-Einträge für V-1, V-2, F-1, F-2; danach `make verify-closure-notes` und der reine Move nach `done/`.
3. **Welle-97:** Closure-Trigger der Welle (Probe am Image, Vorher-/Nachher-Lauf) — die Belege in Punkt 6 und 7 sind dafür verwendbar, die Einlösung gehört der Welle ([`AGENTS.md`](../../AGENTS.md) §5 Regel 10).
4. **Release-Prep:** README-Beschreibung „Zeilen, die mit `|` beginnen“ (beide Sprachen), CHANGELOG, Handbuch-Kopf und -§11 — bewusst nicht im Slice.

## Prozess-Anmerkung zur Messung

Während des Laufs hat ein zweiter `make gates`-Lauf aus einer anderen Sitzung auf demselben Arbeitsbaum und Docker-Daemon gearbeitet (Log `gates-w97h.log` im gemeinsamen Scratchpad). Die zitierten Gate-Zeilen stammen aus dem eigenen Log dieser Sitzung; beide Läufe enden grün. Beim Anlegen des Mutations-Arbeitsbaums im gemeinsamen Scratchpad habe ich ein bereits vorhandenes Verzeichnis `mut` überschrieben (Scratch der früheren Sitzung, nichts im Repo); alle eigenen Arbeitsdateien dieser Verifikation liegen seither in einem eigenen Unterverzeichnis.
