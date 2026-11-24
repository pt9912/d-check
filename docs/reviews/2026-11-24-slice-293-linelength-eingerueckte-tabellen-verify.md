# Verifikations-Report — slice-293 (Anforderung, Spezifikation und Folge-ADR: eingerückte Tabellenzeilen in `linelength`)

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?“, gegen DoD und Spec, nicht gegen den Plan-Stil
**Gegenstand:** [`slice-293`](../plan/planning/in-progress/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md), `git log 686c9dbd..HEAD` auf `main` (lokal): `78b0131e`, `835907f5`, `11bc17d8`, `fe741033` (Review R1), `a3b06428` (Nachbesserung F-1), `c1b63fcf` (Plan-Nachtrag dazu). Die Folge-ADR [`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) liegt seit `6cb9c261`, also **vor** dem Start-Move `686c9dbd`, und ist im Diff unverändert.
**Eingang:** Plan §1 bis §3 und §6, Review [R1](2026-11-24-slice-293-linelength-eingerueckte-tabellen-r1.md), Architect-Verdikt [welle-97](2026-11-24-welle-97-architect-verdikt.md), Lastenheft [`DC-FA-LLEN-001`](../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in) 0.115.1, [`DC-FA-LLEN-001.a`](../../spec/spezifikation.md#dc-fa-llen-001a--zeilenlänge-in-markdown-dateien-linelength).
**Datum:** 2026-11-24 · **Modell-ID:** claude-sonnet-5-5

**Belege dieses Laufs (selbst gefahren, nicht übernommen):**

- `make gates` auf `HEAD` (Arbeitsbaum clean): Exit 0, letzte Zeilen unten.
- `git diff 686c9dbd..HEAD --stat`: vier Dateien (Slice-Plan, Review-Report, `spec/lastenheft.md`, `spec/spezifikation.md`).
- Messbild am Image `d-check:latest` mit einem Fixture im Scratchpad (`max-chars: 10`, `exempt-unbreakable: false`, damit nur (b) antwortet), je ein Lauf im Default, mit `exempt-tables: false` und mit `exempt-code-blocks: false`.
- Code gelesen: `tableRowLine` und ihre Konsumenten (`grep -rn 'tableRowLine(' internal`), `CheckLinelength`, `SplitPipeTableLine`, `declaredObservationIDs`, `trace_table.go`.

## 1. DoD-Punkt für DoD-Punkt

| # | DoD-Punkt | Urteil | Beleg |
|---|---|---|---|
| 1 | [`DC-FA-LLEN-001`](../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in) im Lastenheft geändert | **bestätigt** | Version 0.114.3 → 0.115.1; Historie: **eine** Zeile je Version, die Zeile 0.115.0 ist byte-gleich zu der aus `835907f5` (`cmp` auf `grep '^| 0.115.0 '`: gleich), 0.115.1 steht als eigene Zeile darüber; Verweis-Spalte `—` ([`MR-032`](../../harness/conventions.md#mr-032)). Die Stellen (b), Einleitungssatz, Punkt „Tabellenzeile heißt …“, Out-of-Scope (4) und die Kriterien sind umgeschrieben, keine alte Aussage steht daneben (Suche unten). Das alte Kriterium „Ausnahme (b), Grenzen“ ist ersetzt, nicht dupliziert (im Abschnitt genau ein Treffer auf „Ausnahme (b), Grenzen“). |
| 2 | Spezifikation führt Verfahren und Schema | **bestätigt** | Schritt 4, „Nicht Gegenstand“, „Grenzen, am Gegenstand gelesen“, Schema-Zeile `linelength.exempt-tables` und zwei Historie-Zeilen (0.115.0 und Nachzug) liegen vor; vier Grenz-Aussagen am Image und die Konsumenten-Menge am Code gelesen (Abschnitt 3). |
| 3 | Folge-ADR `Proposed`, im Index, mit Trigger | **bestätigt** | [`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md): `**Status:** Proposed`; Index-Zeile in [`README.md`](../plan/adr/README.md) (Status `Proposed`); `## Re-Evaluierungs-Trigger` mit fünf Triggern; `## Verglichene Alternativen` A bis H, darunter B („`tableRowLine` selbst weiten“) und C („eigener Schalter“); `**Supersedes:**` benennt **nur** den Satz „eingerückte Zeilen sind keine“ und die Konsequenz zur eingerückten Tabellenzeile von Entscheidung 5 (b) in [`ADR-0114`](../plan/adr/0114-linelength-modul.md) — der Träger der Teil-Ablösung. Der Übergang nach `Accepted` und die Geschichte-/Index-Zeile von ADR-0114 stehen bewusst bei `slice-294` (§1, Verdikt F-6); `git diff --stat` zeigt ADR-0114 unberührt. |
| 4 | `make gates` grün | **bestätigt** | selbst gefahren: Exit 0 (Zeilen unten). |
| 5 | Review durchgeführt, Report liegt vor | **bestätigt** | [R1](2026-11-24-slice-293-linelength-eingerueckte-tabellen-r1.md): 0 HIGH, 1 MEDIUM (F-1), 1 LOW, 1 INFO. F-1 ist durch `a3b06428` eingearbeitet (Abschnitt 4). Kein offener HIGH oder MEDIUM. |
| 6 | Closure-Notiz mit Lerneintrag | **offen, Closure-Handlung** | §7 steht auf „*(Wird bei der Closure gefüllt.)*“ — vor der Closure erwartet, kein Befund. |
| 7 | Beobachtungs-Register fortgeschrieben | **offen, Closure-Handlung** | noch nichts im Diff; die Klassen aus R1 und dieser Verifikation (Abschnitt 5) gehören in die Closure. |
| 8 | Jedes Risiko aus §6 trägt einen Ausgang | **offen, Closure-Handlung** | Vorschläge in Abschnitt 6. |
| 9 | Die drei Paarungen | **offen, Closure-Handlung** | im Repo ohne Wellen-Betrieb bei der Closure zu prüfen. |

Die Punkte 6 bis 9 sind nach dem Plan Closure-Pflichten, die der Planner nach dieser Verifikation erfüllt; sie sind **nicht** verletzt, aber auch noch nicht bestätigt. Die vier inhaltlichen Punkte 1 bis 5 sind bestätigt.

**Bewusstes Brechen (Modul 11):** Dieser Slice liefert Text, keinen Test; kein DoD-Punkt beruft sich auf einen grünen Test, der ohne einen Fix rot laufen müsste. Der Sensor-Beleg ist hier die Gegenprobe am Gegenstand (Abschnitt 3), nicht ein Mutationstest. Die Tests der Anforderung (Mutation je Test) trägt `slice-294` ([DoD dort](../plan/planning/open/slice-294-linelength-eingerueckte-tabellen-kern.md)).

### Zu Punkt 1: Entscheid → Kriterium, selbst zugeordnet

Die sieben Kriterien der (b)-Familie im Lastenheft: „an und aus“ (unverändert), „eingerückt, an und aus“ (K1), „Grenzen“ (K2, ersetzt), „eingerückt im Zaun“ (K3), „Bezugsmenge“ (K4), „andere Module unverändert“ (K5), „Rest unverändert“ (K6). Das sind fünf neue, ein ersetztes und ein unverändertes, also 7 ≤ 8 (Obergrenze der Welle).

| Entscheid ([`ADR-0116`](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md)) | Träger im Lastenheft | Kriterium |
|---|---|---|
| 1 — kein neuer Schlüssel, `false` misst | Absatz „Drei Ausnahmen“ nennt nur `exempt-tables`; Schema ohne neue Zeile | K1 („with `exempt-tables: false`, then je ein Befund“), K4 („with … false, then zählt sie“) |
| 2 — `> \|` bleibt gemessen | Punkt „Tabellenzeile heißt …“ | K2 (`> |`-Zeile) |
| 3 — höchstens drei Leerzeichen; Tab, U+00A0, vier gemessen; im Zaun gemessen | (b), Punkt „Tabellenzeile heißt …“ | K1 (1, 2, 3), K2 (4, Tab, U+00A0, Spalte-0-Zeile im Zaun), K3 (3 Leerzeichen im Zaun) |
| 4 — eigene, benannte Frage in der geteilten Lexik; andere Module unverändert | Einleitungssatz, Out-of-Scope (4) | K5 (die drei Module), K6 (Rest unverändert) |

Zusätzlich tragen die drei Grenzen und die Prosa-Zeile **Sätze ohne eigenes Kriterium** (Verdikt F-4, vom Plan so entschieden); die Tests dafür stehen in `slice-294` §3 (Prosa-Zeile, `- | a | b |`, vier Leerzeichen über K2). Das ist eine benannte Abweichung von der Faustregel „jede Randbedingung ein Kriterium“ (INFO V-3).

**Alte Aussagen weg (eigene Suche, dasselbe Kommando wie im Plan-Nachtrag, auf `HEAD`):** Von den Treffern zu „eingerückte Tabellenzeile“, „dieselbe Antwort“, „keine zweite“, „zweite Antwort“ ist jeder entweder die umgeschriebene Fassung (Lastenheft 4604, 4618 bis 4624, 4641 bis 4651, 4720 bis 4725, 4744, 4748; Spezifikation 3786, 4130), eine unberührte Zeile eines anderen Gegenstands (Lastenheft 3246 und 4356; Spezifikation 2419, 2801, 3040, 3555, 4135) oder eine Historie-Zeile. Zweite Suche („übrigen Module“, „andere Module“ …): die verbleibenden Treffer betreffen Opt-out-Marker, Vorverarbeitung, Eingabefehler und §5, nicht die Tabellenzeile.

## 2. Konsistenz Lastenheft ↔ Spezifikation ↔ ADR-0116, je Pflichteingabe

„Image heute“ ist das Verhalten **vor** `slice-294` (der Code ist unverändert); die Erwartung der Spec ist die Spalte „Soll“. Fixture: eine Datei mit je einer überlangen Zeile pro Fall plus Zaun.

| Eingabe | Soll (LH · Spez · ADR) | Image heute | Stellen einig |
|---|---|---|---|
| 1 bis 3 Leerzeichen vor `\|` | ausgenommen | Befund (Zeilen 2, 3; ändert sich mit `slice-294`) | ja (LH (b), Punkt, K1; Spez Schritt 4, Schema; ADR Entscheid 3) |
| 4 Leerzeichen | gemessen | Befund (Z. 4) | ja |
| Tab vor `\|` | gemessen | Befund (Z. 5) | ja (Schema-Zeile nennt ihn) |
| U+00A0 vor `\|` | gemessen | Befund (Z. 6) | ja, Schema-Zeile nennt ihn nicht (INFO, schon R1 F-3) |
| `> \|` | gemessen | Befund (Z. 7) | ja |
| `a \| b \| …` ohne führende Pipe | gemessen | Befund (Z. 8) | ja, Schema-Zeile nennt ihn nicht (R1 F-3) |
| `- \| a \| b \|` | gemessen | Befund (Z. 9) | ja (LH Punkt, Spez „Grenzen“, ADR Konsequenzen); ohne eigenes Kriterium (V-3) |
| `10. \| …` | gemessen | Befund (Z. 10) | ja (LH Punkt, Spez „Grenzen“, ADR); die Aussage „Tabelle unter `10. ` steht mit vier und mehr Leerzeichen“ ist CommonMark-Wissen, **kein** Messergebnis des Werkzeugs; der Plan benennt das selbst |
| Prosa-Zeile `\|` ohne Tabelle, Spalte 0 | ausgenommen | kein Befund (Z. 12), mit `exempt-tables: false` Befund | ja |
| Prosa-Zeile `\|` ohne Tabelle, 3 Leerzeichen | ausgenommen | Befund (Z. 11; ändert sich) | ja |
| eingerückte und Spalte-0-`\|`-Zeile im Zaun, `exempt-code-blocks: false` | gemessen | Befund beide (Z. 14, 15) | ja (LH Punkt, K2, K3; Spez „Grenzen“: „auch eine eingerückte“; ADR Konsequenzen) |
| dieselben Zeilen im Zaun, Default | nicht gemessen | kein Befund | ja |
| `exempt-tables: false` | jede `\|`-Zeile gemessen | 12 Befunde (alle Fälle inkl. Z. 1 und 12) | ja (LH K1, K4; Spez Schema; ADR Entscheid 1) |

Bezugsmenge am Image: Default „10 Zeile(n) geprüft“ (Zeilen 2 bis 11), mit `exempt-tables: false` „12“ (die zwei Spalte-0-Zeilen kommen dazu) — die Zusage aus K4 („zählt nicht zu den geprüften Zeilen“ im Default, „zählt“ mit `false`) folgt der Zählung in `CheckLinelength` (`TablesExempt` vor `count.Lines++`). Es gibt **keinen** Widerspruch zwischen den drei Stellen.

## 3. Spezifikation am Gegenstand (Regel 13)

Stichprobe, mindestens vier Grenz-Aussagen — alle gegen Code bzw. Image gehalten:

1. **„Tab, U+00A0, vier Leerzeichen, `> |`, `a | b`, `- |`, `10. |` bleiben gemessen“** — am Image: je Befund im Default (Z. 4 bis 10). `tableRowLine` ist `prose[i+1] && strings.HasPrefix(lines[i], "|")`, trimmt nichts; ein neuer Helfer, der null bis drei U+0020 zulässt, lässt diese Fälle unberührt. Stimmt.
2. **„Spalte-0-Prosa-Zeile mit `|` ohne Tabelle ist ausgenommen, wie heute“** — Z. 12 im Default kein Befund, mit `exempt-tables: false` Befund. Stimmt.
3. **„Mit (a) aus und (b) an liegt eine `|`-Zeile im Zaun, auch eine eingerückte, nicht in `proseLineSet` und wird gemessen“** — Image, `exempt-code-blocks: false`: Z. 14 (3 Leerzeichen) und Z. 15 (Spalte 0) je ein Befund; Default: keiner. Stimmt; `CheckLinelength` prüft `!prose[i+1]` erst bei `CodeBlocksExempt`.
4. **Bezugsmenge „geprüft ist, was nach (a) und (b) bleibt“** — Image: 10 bzw. 12 (oben). Stimmt.
5. **Wortgrenze von (c): Tab trennt, U+00A0 nicht, „Leerzeichen“ in (b) ist nur U+0020** — `longestWord` trennt an `' '` und `'\t'`, an keinem anderen Zeichen; die Unterscheidung im Kriterium K2 („anders als der Leerraum der Wortgrenze in (c), zu dem auch der Tab gehört“) trifft den Code.
6. **Menge „`targets`, `planning.waves` und `structure`“ (Lastenheft und Spezifikation nennen **dieselbe** Menge):** `grep -rn 'tableRowLine(' internal` (ohne Tests): `targets.go:246`, `planning_waves.go:277`, `structure_tableorder.go:129`, `structure_tablecell.go:43` und `:93`, `linelength.go:72`, dazu `markdown.go:205` (`tableHeaderOrSeparator`, Konsumenten `planning.waves` und `structure`). Das sind genau `targets`, `planning.waves`, `structure` und `linelength`; `tableSeparatorRow` prüft das führende `|` roh. Die Aussage „diese drei lesen eine eingerückte Zeile weiter nicht als Tabellenzeile“ trifft also die Konsumenten von `tableRowLine`. Stimmt, und Lastenheft (Absatz, K5), Spezifikation (Schritt 4, „Grenzen“, Historie), ADR-0116 (Entscheid 4) nennen dieselbe Menge.
7. **Die Nachbesserung 0.115.1 („der RTM-Leser und der Leser des Beobachtungs-Registers in `planning` trimmen die Zeile“)** — `declaredObservationIDs` macht `strings.TrimSpace(ln)` vor dem Test auf `|`; `SplitPipeTableLine` (vom `trace`-Leser in `trace_table.go` benutzt) beginnt mit `strings.TrimSpace(line)`. Beides stimmt. Die Menge ist damit auch gegen die Gegenprobe richtig: `metrics` hat seine eigene Spalte-0-Leseweise (`atColumnZero`) und wird weder genannt noch berührt.

Nicht am Werkzeug messbar: das Soll-Verhalten (1 bis 3 Leerzeichen ausgenommen) — es entsteht mit `slice-294`; der Slice gibt es als Zusage, nicht als Ist aus, und seine Nachträge nennen den Ist-Stand ehrlich.

## 4. Review R1 und die Nachbesserung `a3b06428`

- **Schließt die Nachbesserung F-1?** Ja. Die drei Stellen im Absatz zur geteilten Lexik nennen jetzt `targets`, `planning.waves` und `structure`; der Spezifikations-Schritt 4 ebenso; die Zeile 0.115.0 bleibt als Chronik byte-gleich, die Präzisierung steht in 0.115.1 mit dem Grund (RTM- und Register-Leser). Die zweite Suche des Plan-Nachtrags liefert keine weitere Stelle, die die Ausweitung trägt. Gegen den Code: Abschnitt 3, Punkt 6 und 7.
- **Neuer Widerspruch?** Keiner. Die Version 0.115.1 trägt `präzisiert — keine Verhaltensänderung` (Patch-Stufe, wie die Reihe 0.112.x bis 0.114.x), und die Kopfzeile `**Version:** 0.115.1` entspricht der obersten Historie-Zeile.
- **Verbleibende R1-Befunde:** F-2 (LOW, Zählung „zwei/drei Ausformungen“) und F-3 (INFO, kürzere Grenzen-Aufzählung der Schema-Zeile) sind nach der Review-Vorgabe der Welle nicht eingearbeitet und gehören ins Register; beide bleiben unverändert wahr (F-3 hier in der Tabelle in Abschnitt 2 wieder sichtbar).

## 5. Weitere Prüfungen

- **Referenz-Richtung ([`AGENTS.md`](../../AGENTS.md) §3.4):** hinzugefügte Zeilen von `spec/` gegen `ADR|slice|welle|Hash|issue` gelesen. Treffer: das Wort „waves“ im Modulnamen `planning.waves` und, in der Historie-Zeile 0.115.0, „Begründung in begleitender ADR“ (keine Kennung) sowie „GitHub-Issue 20“. Keine ADR-Kennung, kein Slice, keine Welle, kein Hash. „GitHub-Issue“ ist eine externe Feldmeldung, keine der fünf Kategorien; der Bestand des Lastenhefts führt drei solche Nennungen vor diesem Slice (`git show 686c9dbd:spec/lastenheft.md`), und `make doc-check` ist in `make gates` mit `matrix` ohne Befund. Kein Spec-Satz zeigt abwärts.
- **Abgrenzung §1 (`git diff 686c9dbd..HEAD --stat`):** vier Dateien — Slice-Plan, Review-Report, `spec/lastenheft.md` (+45/−), `spec/spezifikation.md` (+34/−). Unberührt: ADR-0116 (zuletzt `6cb9c261`, vor dem Start-Move), ADR-0114, ADR-Index, jeder Go-Code, `.d-check.yml`, READMEs, Benutzerhandbuch, CHANGELOG, `config_template.go`, Lastenheft §5 (Zeilen 4786, 4814 unverändert). Kein Zitat-Fall, kein neuer Schalter, keine Listen-Erkennung als Leistung (sie steht in Out-of-Scope (4) als Ausschluss). Die Abgrenzung hielt.
- **Pflicht „Plan von `slice-294` aktualisieren“ (§3, letzte Zeile der Tabelle):** kein Diff am Plan von `slice-294` in diesem Bereich; er war schon mit `ba1f55e6` (vor dem Start) angepasst und trägt die Testfälle aus den bestätigten Entscheiden (Prosa-Zeile, `- | a | b |`, zwei Assertionen im Kopplungs-Test, Mutationen). Kein Befund, nur die Notiz, dass die Plan-Zeile schon vor dem Slice erledigt war.
- **Lastenheft §5:** Der Satz „mit derselben Markdown-Lexik wie die übrigen Module“ bleibt wahr, solange der Helfer in der geteilten Lexik wohnt ([ADR-0116](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) Entscheid 4). Wanderte er nach `linelength`, wäre §5 berührt (Entscheidung des Auftraggebers).

## 6. Vorgeschlagene Ausgänge der Risiken aus §6

| Risiko (Kurz) | Vorschlag | Beleg |
|---|---|---|
| 1 — eine Stelle steht gegen die andere (Zusage „dieselbe Antwort … keine zweite“, Punkt, Kriterium, Out-of-Scope (4)) | **entfallen** | Abschnitt 1 (Suche: jede alte Aussage umgeschrieben), Abschnitt 2 (zwölf Pflichteingaben, drei Stellen einig). R1 F-1 ist **kein** Widerspruch der Stellen untereinander, sondern ein Satz mit zu weitem Umfang; er war **eingetreten** als eigene Klasse und ist durch `a3b06428` behoben. Für das Register: F-1 ist eine Evidenz-Datei für [`BEO-ALL/grenz-aussage-ohne-ihre-bedingung`](../plan/planning/observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/observation.md) (2.), nicht für `plan-entscheid-steht-gegen-lastenheft`. |
| 2 — Anforderung wächst über sechs Kriterien (Obergrenze acht) | **entfallen** | sieben Kriterien in der (b)-Familie (fünf neu, eines ersetzt, eines unverändert), ≤ 8; vier Entscheide der ADR unverändert (Entscheid 1 bis 4). |
| 3 — enge Antwort zu Frage 3 lässt `10. `/verschachtelte Tabellen als Befund stehen; Bestand des Melders nicht gemessen | **weiter offen** | der Auftraggeber hat die enge Fassung am 2026-11-24 bestätigt (Plan §3, ADR Kontext); die Anforderung nennt die Grenze mit ihrer Bedingung (Punkt „Tabellenzeile heißt …“, Spezifikation „Grenzen“, ADR Konsequenzen); der Bestand des Melders ist weiter nicht gemessen, der erste Re-Evaluierungs-Trigger der ADR trägt ihn. Register-Eintrag bei der Closure. |
| 4 — Lastenheft und Spezifikation führen dieselbe Aussage in eigenen Worten; eine Stelle steht gegen die andere | **entfallen** | zweite Suche nach dem Editieren: keine Stelle trägt die alte Aussage; die „übrigen Module“-Suche nach F-1 ebenso; Menge der drei Module in Lastenheft, Spezifikation und ADR gleich (Abschnitt 3, Punkt 6). |

## 7. Befunde dieses Laufs

Keine HIGH, kein MEDIUM, keine DoD-Verletzung. Drei kleine Punkte, alle **nicht blockierend** und nach der Welle-Vorgabe ins Register, nicht in den Text:

| ID | Kategorie | Kurz |
|---|---|---|
| V-1 | INFO | Lastenheft, Absatz zur geteilten Lexik: „Für Fenced-Code ist es dieselbe Antwort, die die übrigen Module auf diese Frage geben“ — „die übrigen Module“ ohne Eingrenzung, derselbe Wortlaut wie die Klasse von R1 F-1. Der Satz ist inhaltlich unverändert aus 0.112/0.113 und betrifft eine andere Frage (Fenced-Code); `versions` und `diagrams` lesen Fences bewusst, ohne die Frage „was ist Fenced-Code“ anders zu beantworten. Kein Widerspruch, nur ein Satz, der bei strengem Lesen dieselbe Eingrenzung bräuchte. Eintrag als weitere Evidenz zu `grenz-aussage-ohne-ihre-bedingung`. |
| V-2 | INFO | Die Historie-Zeile 0.115.0 des Lastenhefts trägt weiter den zu weit gefassten Satz („die übrigen Module lesen eine eingerückte Zeile weiter nicht als Tabellenzeile“); sie bleibt nach [`MR-032`](../../harness/conventions.md#mr-032) byte-gleich, und die direkt darüber stehende Zeile 0.115.1 korrigiert sie ausdrücklich. Das ist Chronik, kein Befund; nur für den Leser: wer die 0.115.0-Zeile allein liest, bekommt die weite Aussage. |
| V-3 | INFO | Die drei Grenzen (Prosa-Zeile, `- \| a \| b \|`, `10. `) und die Prosa-Zeile tragen **keine eigenen Lastenheft-Kriterien**; das ist Verdikt F-4 und Plan-Entscheid, ihre Tests stehen in `slice-294` §3. Das Beobachtungs-Register führt [`BEO-ALL/spec-randbedingung-ohne-test`](../plan/planning/observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md) (2×) mit der Faustregel „jede Randbedingung ein Kriterium“; §8 des Plans zitiert sie. Wer sie als Zähler-Evidenz lesen will, trägt diese bewusste Abweichung ein (kein 3. Treffer, weil die Tests geplant sind). Dazu bleibt R1 F-3 (Schema-Zeile kürzer als der Vertrag) bestehen. |

## 8. Plan-vs-Text-Diff

| Plan-Zusage | Text |
|---|---|
| (b) und Punkt „Tabellenzeile heißt …“ nennen die Einrückung bis drei Leerzeichen | erfüllt (LH 4604, 4641 bis 4651) |
| „dieselbe Antwort … keine zweite“ und Out-of-Scope (4) umgeschrieben | erfüllt (LH 4617 bis 4625, 4741 bis 4748) |
| Kriterium „(b), Grenzen“ geändert, kein zweites daneben | erfüllt (ersetzt; LH 4721) |
| sechs Kriterien aus §3 | erfüllt (K1 bis K6; zusätzlich „(b), an und aus“ unverändert) |
| drei Grenzen mit Bedingung als Sätze, keine neuen Kriterien | erfüllt (LH 4641 bis 4651, Spez 3781 bis 3794) |
| Zusage Spezifikation: Schritt 4, „Nicht Gegenstand“, „Grenzen“, Schema, Historie | erfüllt |
| Folge-ADR `Proposed`, Index, Trigger, Alternativen B und C | erfüllt (liegt seit `6cb9c261`) |
| Lastenheft §5 unverändert | erfüllt |
| Handbuch, Code, CHANGELOG, READMEs, ADR-0114 unberührt | erfüllt |
| Kriterium 2 behält den Spalte-0-Zaun-Fall | erfüllt (LH 4721, letzter Fall) |
| Out-of-Scope (4) nimmt die Listen-Erkennung auf, Bedingung „kein Listenkontext“ | erfüllt (LH 4745 bis 4748) |

Eine Abweichung vom Wortlaut der Plan-Tabelle: Der Plan nennt sechs Kriterien „(neu oder ersetzt)“, der Lastenheft-Text hat dazu das unveränderte „(b), an und aus“; die Plan-Zählung benennt das selbst („das Kriterium ‚an und aus‘ bleibt“). Keine Abweichung in der Sache.

## 9. Gesamturteil

**closure-reif**, sobald der Planner die Closure-Pflichten erfüllt (Notiz mit Lerneintrag, Register: F-1 samt Klasse `modul-aussage-ueber-gegenstand-hinaus` und V-1, F-2 und F-3 und V-3 als LOW/INFO, Risiko-Ausgänge nach Abschnitt 6, drei Paarungen). Die inhaltlichen DoD-Punkte 1 bis 5 sind am Text, am Code und am Image bestätigt; es gibt keine DoD-Verletzung und keinen Rest aus R1. Die Spezifikation behauptet an keiner der gelesenen Stellen mehr, als der Gegenstand trägt, und die Aussagen über Soll-Verhalten sind als Zusage von `slice-294` gekennzeichnet, nicht als Ist.

**Übergabe:** Verifier → Planner. ADR-0116 bleibt `Proposed`; `Accepted` setzt erst `slice-294`.

## Letzte Zeilen von `make gates` (selbst gefahren, Exit 0)

```
d-check: 1296 Datei(en) geprüft, 0 Befund(e)
[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green
exit 0
```
