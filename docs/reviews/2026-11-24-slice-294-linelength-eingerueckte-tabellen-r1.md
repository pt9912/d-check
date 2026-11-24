# Review-Report — slice-294 (Kern und Tests), Runde 1

**Review-Art:** Code (Diff gegen Plan, Folge-ADR, Architect-Verdikt, Spezifikation und Hard Rules)
**Gegenstand:** slice-294, Diff `482e0dc8..b10dbb98` (lokal auf `main`): `222b5bcb` (Stufe 1: Helfer `unbreakableTableRow`, Aufruf in `linelength`, Kopplungs-Test, Tests), `a4a91511` (Stufe 2: Handbuch §4.28, `--print-config`-Vorlage, ein Zaun-Test), `b10dbb98` (Plan-Nachtrag)
**Skill:** `reviewer.md` @ 1.20.0
**Modell-ID:** claude-sonnet-5-5
**Datum:** 2026-11-24
**Eingangs-Kontext:**
- Plan slice-294 (§1 Abgrenzung, §3 Tabelle, Testfälle und Mutationstabelle, §6 Risiken, Plan-Nachtrag)
- [`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) (`Proposed`), Entscheide 1 bis 4; das [Architect-Verdikt zu welle-97](2026-11-24-welle-97-architect-verdikt.md); [`ADR-0054`](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) (Kopplung der geteilten Lexik); [`ADR-0114`](../plan/adr/0114-linelength-modul.md)
- [`DC-FA-LLEN-001`](../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in) (Lastenheft 0.115.1) und [§DC-FA-LLEN-001.a](../../spec/spezifikation.md#dc-fa-llen-001a--zeilenlänge-in-markdown-dateien-linelength) Schritt 4 samt „Nicht Gegenstand“ und „Grenzen, am Gegenstand gelesen“
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.2, §3.7, §3.8 und §5 Regel 13 bis 18
- Code als Gegenstand: `internal/hexagon/core/rules/markdown.go`, `linelength.go`, `lexikon_kopplung_test.go`, `linelength_test.go`, die Konsumenten von `tableRowLine`
- Vorherige Findings am Modul: der Report [`2026-11-24-slice-293-linelength-eingerueckte-tabellen-r1.md`](2026-11-24-slice-293-linelength-eingerueckte-tabellen-r1.md) (Spec-Text, nicht Code); dessen Klasse „Aussage über den Gegenstand hinaus“ wurde als Prüfanker für die Handbuch-Aussagen gelesen.

**Messbasis dieses Laufs:** `make test` auf dem Arbeitsbaum von `HEAD` grün, danach je Mutation ein eigener `make test`-Lauf mit Rücknahme (`git status` und `git diff` nach dem Lauf leer). Das Image `d-check:latest` wurde mit `make build` neu gebaut und trägt den Stand von `HEAD` (`--print-config` zeigt den neuen Kommentar). Fixture: `max-chars: 20`, `exempt-unbreakable: false` (damit nur (b) antwortet), 17 Zeilen, aufgerufen als `docker run --rm --network none -v "$FIXTURE":/repo:ro d-check:latest --enable linelength`.

## Übersicht

| ID | Kategorie | Kurz | Aktion |
|---|---|---|---|
| [F-1](#f-1--low--die-begründung-keinen-listenkontext-hängt-an-fällen-die-sie-nicht-trägt) | LOW | Kommentare in `markdown.go` und `linelength.go` begründen „Tab, U+00A0, `>`“ mit dem fehlenden Listenkontext | ins Register |
| [F-2](#f-2--info--der-kommentar-von-tablerowline-nennt-seine-zweite-frage-nicht) | INFO | `tableRowLine` trägt weiter „die EINE Antwort“ ohne Zeiger auf den Helfer daneben | ins Register |

## Findings

### F-1 — LOW — Die Begründung „keinen Listenkontext“ hängt an Fällen, die sie nicht trägt

- **kategorie:** LOW · **quelle:** [`AGENTS.md`](../../AGENTS.md) §3.7 (Klasse *Grenze*: Aussage mit ihrer Bedingung) und §5 Regel 13 (eine aufgeschriebene Grenze gegen den Gegenstand prüfen); Prüffrage 20 dieses Skills, **bewusst eine Stufe tiefer** (siehe Begründung unten) · **klasse:** begruendung-trifft-teilmenge · **verifizierbar:** ja (Kommentare gegen [`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) Entscheid 2 und 3 und gegen die Fixture lesen)
- **pfad 1:** `internal/hexagon/core/rules/linelength.go` · „U+00A0 oder `>` davor wird gemessen“ (Dateikopf, Satz „… wird gemessen, weil die Lexik keinen Listenkontext hergibt“)
- **pfad 2:** `internal/hexagon/core/rules/markdown.go` · „weil die Lexik keinen Listenkontext trägt“ (Kommentar von `unbreakableTableRow`, Satz „ein Tab, ein U+00A0 oder vier und mehr Leerzeichen vor dem `|` zählen nicht, …“)
- **befund:** Beide Sätze führen die Grenze für vier und mehr Leerzeichen, Tab, U+00A0 und (im Dateikopf) `>` auf einen Grund zurück, den nur die Grenze „vier und mehr Leerzeichen“ trägt (kein Listenkontext, dazu der Schnitt der Anforderung beim eingerückten Code-Block). Ein U+00A0 ist kein Einrückungs-Leerzeichen der Markdown-Lexik, und `> |` bleibt gemessen, weil [`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) Entscheid 2 es so entscheidet, nicht wegen einer fehlenden Listen-Erkennung. Die Entscheidung bliebe mit falscher Begründung dieselbe (Prüffrage 20). Handbuch, Spezifikation und ADR sind an dieser Stelle richtig: sie hängen den Listenkontext nur an die Listenmarke und an `10. `.
- **Einordnung:** Prüffrage 20 steht im Skill auf MEDIUM und verlangt ein Failure-Szenario. Das einzige erzählbare ist, dass ein späterer Bearbeiter nach dem Re-Evaluierungs-Trigger „die Lexik bekommt einen Listenkontext“ ([`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md)) aus dem Kommentar schließt, U+00A0 und `>` fielen mit, und die Grenze falsch weitet. Die kanonischen Quellen (ADR, Spezifikation, Handbuch) halten dagegen, der Kommentar ist nur die nähere. Deshalb LOW mit dem Konventions-Anker §3.7, nicht MEDIUM.

### F-2 — INFO — Der Kommentar von `tableRowLine` nennt seine zweite Frage nicht

- **kategorie:** INFO · **quelle:** Maintainability (kein Konventions-Anker; daher nicht LOW) · **klasse:** kopplungs-zeiger-nur-in-einer-richtung · **verifizierbar:** ja (Kommentar lesen)
- **pfad:** `internal/hexagon/core/rules/markdown.go` · „Dies ist die EINE Tabellenzeilen-Antwort des Produkts (targets, planning.waves, structure)“
- **befund:** Der neue Helfer verweist in seinem Kommentar auf `tableRowLine` („nicht die von tableRowLine“), aber wer `tableRowLine` liest, erfährt nicht, dass daneben eine Frage steht, die sie einschließt. Die Zusage des Kommentars stimmt für die drei genannten Konsumenten; sie ist wegen des Plans („`tableRowLine` bleibt zeichengleich“) auch nicht angefasst worden.

## Negativbefunde (geprüft, ohne Befund)

### Schwerpunkt 1 — Korrektheit des Helfers

`unbreakableTableRow` ruft zuerst `tableRowLine` (Spalte 0, außerhalb von Fenced-Code) und prüft sonst `indent >= 1 && indent <= maxTableIndent && prose[i+1] && HasPrefix(line[indent:], "|")`. Gezählt wird nur U+0020 (`line[indent] == ' '`). Die Fence-Lage prüft der Helfer selbst über `prose[i+1]`, und zwar für die eingerückte Zeile; für Spalte 0 trägt sie `tableRowLine`.

Fixture (`max-chars: 20`, Zeile 3 bis 17), Messung am Image, Default:

| Zeile | Form | Default | `exempt-tables: false` |
|---|---|---|---|
| 3 | `\|` in Spalte 0 | kein Befund | Befund |
| 4, 5 | ein, drei Leerzeichen | kein Befund | Befund |
| 6 | vier Leerzeichen | Befund | Befund |
| 7 | Tab | Befund | Befund |
| 8 | U+00A0 | Befund | Befund |
| 9 | `> \|` | Befund | Befund |
| 10 | `a \| b …` | Befund | Befund |
| 11 | `- \| …` | Befund | Befund |
| 12 | `10. \| …` | Befund | Befund |
| 13 | Prosa-Zeile mit drei Leerzeichen und `\|`, keine Tabelle | kein Befund | Befund |
| 14 | Leerzeichen, dann Tab | Befund | Befund |
| 16 | drei Leerzeichen im Zaun | kein Befund (Zaun); mit `exempt-code-blocks: false` Befund | — |

Bezugsmenge im Default: 10 geprüfte Zeilen (17 minus vier Tabellen-/Prosa-Zeilen minus drei Zaun-Zeilen); mit `exempt-tables: false` 14. Beides ist die Spezifikation, Schritt 4 und 8.

**`tableRowLine` zeichengleich:** `git diff 482e0dc8..b10dbb98 -- internal/hexagon/core/rules/markdown.go` zeigt nur hinzugefügte Zeilen (28 `+`, keine `-`). Konsumenten gelesen (`grep -rn 'tableRowLine(' internal`): `targets.go`, `planning_waves.go`, `structure_tablecell.go` (zwei Stellen), `structure_tableorder.go`, in `markdown.go` der Helfer und `tableHeaderOrSeparator` (`tableRowLine(lines, prose, i+1)`), im Test die Kopplung. Keine dieser Dateien steht im Diff (bis auf den Kopplungs-Test); `tableHeaderOrSeparator` ruft weiter `tableRowLine`, nicht den Helfer. Keiner der Konsumenten bewegt sich.

### Schwerpunkt 2 — Aufruf in `linelength`

`linelength.go` ruft für (b) nur `unbreakableTableRow`; `tableRowLine` kommt in der Datei nicht mehr vor. Reihenfolge (a) → (b) → `count.Lines++` → Länge → (c) ist im Diff unverändert und deckt sich mit Spezifikation Schritt 4 bis 6 und 8 (Bezugsmenge nach (a) und (b), `U` nach (c)).

### Schwerpunkt 3 — Mutationen

Zwölf Läufe, je einzeln mit `make test` und Rücknahme: zehn Mutationen wurden rot, der Doppelaufruf V blieb grün (äquivalenter Mutant, siehe unten), V + M10 wurde rot. Die Spalte „Plan“ nennt die Zeile der Mutationstabelle in §3, soweit es eine gibt; T1, T2, C1 und L1 sind Zusatzproben.

| Nr. | Mutation | Plan | Ergebnis |
|---|---|---|---|
| M1 | `linelength.go` ruft für (b) wieder `tableRowLine` | Kriterium 1 | rot: `TestLinelengthEingerueckteTabellenzeile`, `…DreiGegenVierLeerzeichen/drei_Leerzeichen`, `…ProsaZeileMitPipeOhneTabelleIstAusgenommen`, `…TabelleUnterNummeriertemListenpunkt`, `…VorUndHinterZaun`, `…Bezugsmenge…`, `…ZaehltNichtInU`, Kopplung `…VonLinelength…/ein_Leerzeichen` bis `drei_Leerzeichen` |
| M2 | `maxTableIndent` 3 → 2 | Grenze 3 | rot: Kopplung `…/drei_Leerzeichen`, `…DreiGegenVier/drei_Leerzeichen`, sieben weitere; ein Leerzeichen bleibt grün |
| M3a | `maxTableIndent` 3 → 4 | Grenze 3/4 | rot, genau zwei Fälle: Kopplung `…/vier_Leerzeichen`, `…DreiGegenVier/vier_Leerzeichen` (der Fall „fünf Leerzeichen“ bleibt grün, wie er soll) |
| M5 | `prose[i+1]` im Helfer entfällt | Kriterium 3 | rot: Kopplung `…/drei_Leerzeichen_im_Fence`, `…ImZaunWirdGemessen/drei_Leerzeichen` und `/ein_Leerzeichen`, `…VorUndHinterZaun` |
| M9 | `tableRowLine` selbst trimmt Leerzeichen | Kriterium 5 | rot: `TestTabellenzeilenFrageHatEineAntwort/ein_Leerzeichen`, `/eingerueckt`, `/drei_Leerzeichen`, `/vier_Leerzeichen`, dazu `TestCheckTargetsColumnZero` und die vier-Leerzeichen-Fälle von `linelength` |
| M10 | der Helfer ruft `tableRowLine` nicht mehr (Block entfernt) | Kriterium 6 | rot: Kopplung `…VonLinelength…/ausserhalb_Fence`, `TestLinelengthAusnahmeTabellenzeile`, `…SpalteNullUnveraendert`, `…ProsaZeile…`, `TestLinelengthBezugsmenge` |
| T1 | ein Tab zählt als Einrückung | Kriterium 2 | rot: Kopplung `…/Tab`, `…TabellenzeileGrenzen/Tab` und `/Leerzeichen_dann_Tab` |
| T2 | der Helfer nimmt `> ` hinter der Einrückung weg | Kriterium 2 | rot: `…TabellenzeileGrenzen/Zitat_mit_Einrueckung` |
| C1 | `Contains` statt `HasPrefix` für `\|` | Kriterium 2 | rot: `…Grenzen/Pipe_nicht_erstes_Zeichen`, `/Zitat_mit_Einrueckung`, `/Leerzeichen_dann_Tab`, `/drei_Leerzeichen_dann_Tab`, `…Listenmarke…/Marke_mit_Einrueckung` |
| L1 | der Helfer nimmt eine Listenmarke (`- `, `10. `) vor dem `\|` weg, ohne Mindest-Einrückung | Grenze 2 | rot: Kopplung `…/Listenmarke`, `…HinterListenmarkeWirdGemessen/Bindestrich`, `/zweistellig_nummeriert`, `/Marke_mit_Einrueckung` |
| V | `linelength` ruft Helfer **und** `tableRowLine` | die vom Implementer nicht gefahrene Variante | **grün, `make test` Exit 0** — äquivalenter Mutant, siehe unten |
| V + M10 | Variante V und der Helfer ruft `tableRowLine` nicht | dieselbe Variante mit der Mutation zu Kriterium 6 | rot: Kopplung `…VonLinelength…/ausserhalb_Fence` (der Helfer wird dort direkt gefragt) |

**Zur Variante V:** `unbreakableTableRow` ist eine Obermenge von `tableRowLine`; ein zusätzlicher Aufruf von `tableRowLine` in `linelength` ändert keine Antwort und ist beobachtbar nicht zu unterscheiden. Ein Verhaltenstest kann sie nicht rot machen; das ist kein Mangel der Tests. Was V gefährlich machen würde, ist, dass sie die Mutation „der Helfer ruft `tableRowLine` nicht“ auf Ebene `linelength` verdecken könnte — das tut sie nicht, weil die zweite Kopplungs-Assertion den Helfer **direkt** fragt und die Zeile `tableRowLine ja, Helfer nein` ausdrücklich prüft (V + M10 ist rot). Die Zusage des Plans („ruft nur den Helfer“) ist damit durch den Code im Diff gehalten, durch die Tests nur im Zusammenspiel; ich lese das als ausreichend.

Ein zusätzlicher Blick auf zwei Plan-Mutationen: „der Helfer verlangt eine Trennzeile“ (Grenze 1) ist nicht einzeln gefahren; `TestLinelengthProsaZeileMitPipeOhneTabelleIstAusgenommen` ist nach M1 und M10 rot und deckt die Prosa-Zeile mit und ohne Einrückung, die Mutation würde die Zeile im Test messen.

### Schwerpunkt 4 — Kopplungs-Test

`TestTabellenzeilenFrageHatEineAntwort` hält `targets`, `planning.waves`, `structure` und `tableRowLine` auf einer Antwort über `tabellenFaelle` (außerhalb Fence, im Fence, ein, zwei, drei, vier Leerzeichen, Tab). `TestTabellenzeilenFrageVonLinelengthHatEineAntwort` hält `linelength` (über die Ausnahme (b) allein) und den Helfer auf einer zweiten über elf Fälle (Spalte 0, Fence, ein bis drei Leerzeichen, drei Leerzeichen im Zaun, vier, Tab, U+00A0, Zitat, Listenmarke) und prüft zusätzlich „`tableRowLine` ja ⇒ Helfer ja“. Das sind die zwei Assertionen aus [`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) Entscheid 4, in beiden Richtungen. Der Kopf-Kommentar ist wahr: er nennt drei Konsumenten für `tableRowLine`, die fünf Schreibweisen „ein bis vier Leerzeichen, Tab“ stehen als Fälle im Test, und die andere Frage von `linelength` trägt ihre ADR. Die Reichweiten-Notiz darunter ist unverändert und gilt weiter.

### Schwerpunkt 5 — Grenz-Aussagen im Handbuch §4.28 und in der Vorlage

Gegen das Image nachgemessen (Fixture oben):

1. Die YAML-Zeile „Zeilen mit | (bis drei Leerzeichen davor) nicht messen“ stimmt mit der Vorlage überein; `--print-config` des neu gebauten Images zeigt sie.
2. „null bis drei Leerzeichen (U+0020) … eine Trennzeile wird nicht verlangt“: Zeile 3 bis 5 ausgenommen, keine Trennzeile in der Fixture.
3. „Gemessen werden: vier und mehr Leerzeichen, Tab, U+00A0, `> |`, `a | b`“: Zeilen 6 bis 10 und 14 je ein Befund.
4. „`exempt-code-blocks: false` … eine `|`-Zeile im Zaun, auch eine eingerückte“: Zeile 16 ein Befund, im Default keiner.
5. „`- | a | b |` … und `10. ` … gemessen“: Zeilen 11 und 12 je ein Befund.
6. „Eine Prosa-Zeile, die nach ein bis drei Leerzeichen mit `|` beginnt, ohne dass eine Tabelle dasteht, ist ausgenommen“: Zeile 13 kein Befund, mit `exempt-tables: false` ein Befund.
7. „Eine ausgenommene Zeile zählt nicht zu den geprüften Zeilen“: 10 im Default, 14 mit `exempt-tables: false`.
8. „Acht Dinge“: die Aufzählung trägt acht Spiegelstriche (Codepunkte, drei Ausnahmen, Fenced-Code, Tabellenzeile, offener Zaun, Zeilen-Marker, Eingabefehler, leere Menge). Die Überschrift „Sieben Dinge“ von §4.26 (Zeile 1624) ist ein anderer Abschnitt und unberührt.
9. „jene drei Module lesen eine eingerückte Zeile weiter nicht als Tabellenzeile“ nennt `targets`, `planning.waves` und `structure` — wie das Lastenheft 0.115.1 und die Spezifikation, nicht „die übrigen Module“.
10. Modul-Tabelle §6 (Zeile 3148) und die Abdeckungs-Datei tragen die neuen Testnamen; `make test` hält die Datei gegen ihre Ableitung (grün).

Gesucht über den Inhalt der alten Aussage, nicht über die zitierte Stelle: `git grep -n -E 'exempt-tables|Tabellenzeile' -- ':!docs/plan' ':!docs/reviews' ':!CHANGELOG.md' ':!.harness' ':!*_test.go' ':!spec'`. Jede verbleibende Aussage ist eine Spalte-0-Zusage eines anderen Moduls, eine Historie-Zeile, `harness/sensors/image-test.md` (Fixture mit Spalte 0, unverändert richtig) oder die Release-Prep-Stellen (README, Handbuch-Kopf und -Historie), die der Plan benennt.

### Schwerpunkt 6 — Kommentare (`AGENTS.md` §3.7)

Hinzugefügte Kommentare gelesen: Zusage (`unbreakableTableRow`, `maxTableIndent`), Abgrenzung (zu `tableRowLine`), Grenze (je mit Bedingung, siehe aber F-1), Rang-Zeiger (`ADR-0116`). Keine Review-Historie, keine Slice-Nummern, keine Mess-Labels, keine Deliberation über Verworfenes. Die Herkunftsangaben „(DC-FA-LLEN-001, ADR-0116)“ in den Test-Kommentaren und im Helfer entsprechen der Form des Bestands in `linelength_test.go`. Kein `//nolint` im Diff. (Der Verweis „Review F-3“ in der Reichweiten-Notiz von `lexikon_kopplung_test.go` steht im unveränderten Bestand und ist nicht Teil des Diffs.)

### Schwerpunkt 7 — Abgrenzung §1

`git diff 482e0dc8..b10dbb98 --stat`: acht Dateien — der Plan, `docs/user/abdeckung-tests.md`, `docs/user/benutzerhandbuch.md`, `config_template.go`, `lexikon_kopplung_test.go`, `linelength.go`, `linelength_test.go`, `markdown.go`. Das entspricht der Tabelle in §3. Unberührt sind [`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md), [`ADR-0114`](../plan/adr/0114-linelength-modul.md), `.d-check.yml`, die READMEs, `CHANGELOG.md`, `tools/image-test.sh`, `configyaml.go` und die Konsumenten von `tableRowLine`. Die vier Abweichungen des Plan-Nachtrags (acht statt sieben Punkte im Handbuch, Name des Helfers, zwei Kopplungs-Tests, Abdeckungs-Datei im ersten Commit) berühren die Abgrenzung nicht.

### Schwerpunkt 8 — Commit-Botschaften (Regel 15)

`222b5bcb`: „Tests je Kriterium und je Grenze der Anforderung“ — jede der Grenzen (Prosa-Zeile, Listenmarke, vier Leerzeichen, Tab, U+00A0, Zaun, Bezugsmenge) hat einen eigenen Test, der in den Mutationen oben rot wurde. „Nicht in diesem Commit: Handbuch, Vorlage, ADR-0116-Status“ stimmt mit der Datei-Liste. `a4a91511`: „Die Aufzaehlung heisst nun ‚Acht Dinge‘“, „Grenzen mit ihrer Bedingung“ und „Neuer Test: eingerueckte |-Zeile vor, in und hinter einem geschlossenen Zaun“ stimmen mit dem Diff. `b10dbb98`: „vier Abweichungen“ — vier stehen im Nachtrag; „Kein neuer Plan-Inhalt“ trifft zu, weil jede Abweichung als ohne Wirkung auf die Abgrenzung begründet ist. Keine Botschaft nennt eine Probe, die nicht stattfand; keine behauptet die Image-Probe der Welle.

### Weitere Bereiche

| Bereich | Ergebnis |
|---|---|
| Prüffragen 1 bis 5 (stilles Grün, falscher Befund, Hexagon-Richtung, Suppression, Netz) | geprüft, ohne Befund: der Helfer liegt in `core/rules`, importiert nichts Neues; kein Gate wird berührt; kein Netz |
| Prüffrage 13 (Negativtest zum öffentlichen Vertrag) | geprüft, ohne Befund: je Seite der Regel (Einrückung und erstes Zeichen) ein eigener Fall, der ohne die Regel rot wird (C1, T1, T2, M2, M3a) |
| Prüffrage 15 / §3.8 (Eingaben, die das Modul nicht scannt) | geprüft, ohne Befund: der Diff ändert die Eingabemenge nicht |
| Prüffrage 11 (zwei Module derselben Eingabe-Klasse, verschieden) | geprüft, ohne Befund: der Unterschied zwischen `linelength` und den drei Spalte-0-Modulen ist benannt ([`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) Entscheid 4, Lastenheft 0.115.1) |
| Prüffrage 21 (eigenes Muster statt vorhandener Erkennung) | geprüft, ohne Befund: der Helfer ruft `tableRowLine` und liest die Fence-Lage aus `proseLineSet`; das Zählen führender Leerzeichen ist die neue Frage selbst |
| Prüffrage 18 (Grenzen-Liste ohne ihre größte Lücke) | geprüft, ohne Befund: die größte Lücke der Regel (Tabellen unter `10. ` und in verschachtelten Listen bleiben gemessen) steht in Handbuch, Spezifikation und ADR |
| Docker/make-only (§3.1) | geprüft, ohne Befund: alle Läufe über `make` und `docker run` |
| `make doc-check` | 0 Befunde auf `HEAD` vor diesem Report |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** begruendung-trifft-teilmenge · kopplungs-zeiger-nur-in-einer-richtung

## Verdikt

**Merge-blockierend:** nein. Kein HIGH und kein MEDIUM; F-1 (LOW) und F-2 (INFO) gehen ins Register.

**Gesamturteil:** Der Helfer ist korrekt und eng: `tableRowLine` bleibt zeichengleich und alle seine Konsumenten bewegen sich nicht, `linelength` ruft nur den Helfer, Reihenfolge und Bezugsmenge decken sich mit Lastenheft und Spezifikation, die Kopplung hält beide Fragen in beiden Richtungen, und jede der zehn gefahrenen Mutationen wurde rot (der Doppelaufruf ist ein äquivalenter Mutant, und die Mutation „Helfer ruft `tableRowLine` nicht“ bleibt auch mit ihm rot, weil die zweite Assertion den Helfer direkt fragt). Die Grenz-Aussagen des Handbuchs und der Vorlage stimmen am Image mit dem Verhalten überein; „Acht Dinge“ stimmt.

**Reife für die Verifikation:** ja.

**Review-Vorgabe der Welle:** höchstens zwei Runden. Dies ist Runde 1; es gibt nichts einzuarbeiten, eine Runde 2 ist nicht nötig.

**Übergabe:** Die Finding-Klassen gehen in die Closure-Notiz von slice-294 §7 und von dort in den Zähler des Beobachtungs-Registers; F-1 ist ein Kandidat für die Evidenz-Datei zu `BEO-ALL/grenz-aussage-ohne-ihre-bedingung`.
