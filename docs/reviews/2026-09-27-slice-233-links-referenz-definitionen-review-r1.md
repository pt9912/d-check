# Review-Report — slice-233 (`links`: Referenz-Definitionen `[label]: ziel`), R1

**Review-Art:** Code (gegen Slice-Plan, betroffene `DC-*`-Anforderung, ADR-0093/ADR-0094, Hard Rules).

**Gegenstand:** `slice-233`, Commit-Range `6f8d6906..82c941ab` (fünf Commits:
`ae8b1a6c`/`64be56b3` Lifecycle-Move/Beanspruchung, `955a8dc1` Vertrag/ADR-0093,
`1fad1b64` Selbstkorrektur/ADR-0094 (supersedes ADR-0093), `82c941ab`
Feature-Code).

**Skill:** `.harness/skills/reviewer.md` @ `8dcc0452` (Version 1.16.0, 2026-09-07).

**Modell-ID:** `claude-sonnet-5`.

**Datum:** 2026-09-27.

**Eingangs-Kontext:**

- Slice-Plan `docs/plan/planning/in-progress/slice-233-links-referenz-definitionen.md`
  (§1–§8, vollständig gelesen).
- `DC-FA-LINK-001` (`spec/lastenheft.md`, Version 0.93.1), `DC-FA-ANCH-001`
  Out-of-Scope-Ergänzung; Spec-Stelle `DC-FA-LINK-001.a` Schritt 3
  (`spec/spezifikation.md`).
- `docs/plan/adr/0093-links-referenz-definitionen-gemeinsame-extraktion.md`
  (Status `Superseded by ADR-0094`, vollständig gelesen) und
  `docs/plan/adr/0094-links-referenz-definitionen-backslash-grenze-korrigiert.md`
  (vollständig gelesen).
- `AGENTS.md` §3.5 (ADR-Immutabilität), §3.7 (Kommentar-Klassen), §5
  (Grenzen-Regel, Commit-Overclaim-Regel, Zitat-Geltungsbereich).
- `docs/reviews/` nach früheren Reviews zu `markdown.go`/`ExtractLinks`
  durchsucht: `2026-09-27-slice-232-links-zeilenumbruch-review-r1.md`/`-r2.md`
  gefunden (dieselbe Datei, dieselbe gemeinsame Extraktion) — Wiederholungs-Muster
  aus R1 dieser Vorgänger-Review übernommen: (a) Bestandsmessung deckt nicht
  automatisch alle sechs `ExtractLinks`-Konsumenten ab, wenn `.d-check.yml` nur
  einen Teil aktiviert; (b) eine Grenzen-Liste in ADR/Spec kann eine reale,
  über die genannten Grenzen hinausgehende Verhaltens-Facette übersehen, auch
  wenn der aktuelle Bestand 0 Befunde zeigt. Beide Muster wurden in diesem Lauf
  gezielt gegengeprüft (siehe unten).
- **Nicht** erhalten (laut Skill): die DoD-Abhakung selbst.

**Gates unabhängig nachgefahren** (Docker/`make`, kein Host-Go — AGENTS.md §3.1):

| Gate | Ergebnis |
|---|---|
| `make gates` (voll, zehn Glieder) | grün — `baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; Coverage `94.30 %` (Schwelle `93 %`); `semgrep`: 55 Regeln/65 Dateien, 0 Befunde; `doc-check`: `847 Datei(en) geprüft, 0 Befund(e)` |
| `make adr-check RANGE=6f8d6906..82c941ab` (unabhängig, volle Slice-Range statt nur `HEAD~1..HEAD`) | grün — `847 Datei(en) geprüft, 0 Befund(e)` (Modul `vcs`) |
| Runtime-Image (`d-check:latest`, HEAD `82c941ab`) gegen dieses Repo, `--enable external --enable tracked` | `847 Datei(en) geprüft, 107 Befund(e)` |
| Runtime-Image aus `git worktree` auf HEAD `6f8d6906` (vor slice-233), Ziel: `--target runtime`, gegen **denselben** aktuellen Bestand, `--enable external --enable tracked` | `847 Datei(en) geprüft, 107 Befund(e)` — reproduziert die Behauptung „0 neue Befunde über alle sechs Konsumenten" exakt (identische Zahl, Alt- vs. Neu-Binary, gleicher Bestand) |
| Dieselben zwei Images gegen ein Fixture mit toter Definition (`[label]: fehlt.md`) | Alt: `0 Befund(e)` · Neu: `1 Befund(e)` (`target-missing`) — reproduziert den Rot-Beleg des DoD-Punkts „Prüfung" unabhängig |

---

## Zentraler Prüfpunkt: Erkennungs-Regex — Korrektheit und tatsächliche Grenzen

### 1. Die drei benannten Grenzen — jede einzeln adversariell nachgestellt (nicht nur an den vorhandenen Tests)

Gegen den laufenden Container (`d-check:latest`), nicht nur gegen die
vorhandenen Unit-/Akzeptanztests:

| Fixture | Erwartung | Ergebnis |
|---|---|---|
| `> [label]: ziel.md` (Blockquote-Präfix, totes Ziel) | kein Fund | bestätigt |
| `    [label]: ziel.md` (vier Leerzeichen, über der CommonMark-Grenze) | kein Fund | bestätigt (drei Leerzeichen matchen noch, vier nicht mehr — beide Seiten der Grenze geprüft) |
| `[foo\]bar]: ziel.md` (Backslash vor dem eigentlichen Label-Ende) | **ganze Zeile** unerkannt (ADR-0094-Korrektur, nicht nur falsche Label-Grenze) | bestätigt — `0` Refs, auch mit PCRE-Trace der Regex unabhängig nachvollzogen (`[^\]\n]+` kann kein `]` konsumieren, das Backtracking findet keine Position, an der `]:` unmittelbar folgt) |
| `[label]: ziel.md` ⏎ `weiter` (Ziel/Titel hinter Zeilenumbruch) | nicht erkannt | durch Code bestätigt: `parseDefinitionLine` operiert ausschließlich auf `ln.Text`, es gibt **keinen** Lookahead-Pfad für Definitionen (anders als bei Inline-Links/ADR-0092) — strukturell ausgeschlossen, nicht nur ungemessen |

Alle drei benannten Grenzen sind korrekt beschrieben und halten. Die
ADR-0094-Korrektur selbst ist richtig: Der verankerte, greedy `[^\]\n]+`
kann nach einem literalen `]` nicht zurücksetzen und an anderer Stelle neu
ansetzen — ein Fehlschlag der Erfassungsgruppe lässt die **ganze** Zeile
unerkannt, es entsteht kein `LinkRef` mit falscher, aber harmloser
Label-Grenze. Das entspricht exakt der jetzt korrigierten Beschreibung.

### 2. Eigene adversarielle Fälle über die Tests hinaus — eine vierte, unbenannte Grenze gefunden

Vier zusätzliche, nicht in ADR-0093/0094 oder der bestehenden Testsuite
geprüfte Formen, jeweils gegen den laufenden Container verifiziert:

| Fixture | Ergebnis | Bewertung |
|---|---|---|
| `[[nested]]: ziel.md` (Label mit eckigen Klammern ohne Backslash) | kein Fund | sicher — dieselbe „ganze Zeile unerkannt"-Klasse wie Grenze 2, aus demselben strukturellen Grund (`[^\]\n]+` stoppt vor dem ersten `]`, kein Backtrack findet `]:`) |
| `\| [label]: ziel.md \|` (sieht wie eine Tabellenzeile mit Definition aus) | kein Fund | sicher — `|` ist kein `[`, die Verankerung `^ {0,3}\[` greift nicht |
| `[a]: fehlt1.md` / `[b]: fehlt2.md` (zwei Definitionen in Folge) | zwei Funde, je eigene Zeile | korrekt |
| `[label]: fehlt-a.md and see [text](fehlt-b.md)` (Definition, gefolgt von einem Inline-Link auf **derselben** Zeile) | **zwei** Funde: `fehlt-a.md` (Definition) **und** `fehlt-b.md` (Inline-Link) | funktioniert, aber zeigt denselben Mechanismus wie Finding H1: die Definitions-Regex nimmt den gesamten Zeilenrest als „Ziel + optionaler Titel" und `NormalizeTarget` trennt nur am ersten Whitespace — sie prüft nicht, ob der Rest überhaupt wie ein gültiger Titel aussieht |
| **`[TERM]: First In, First Out`** (glossar-artige Prosazeile, **keine** beabsichtigte Link-Referenz-Definition) | **Fund:** `docs/glossary.md:1  First  target-missing  Linkziel existiert nicht` | **Fehlfund — siehe R1-H1** |

Der letzte Fall ist der entscheidende Befund: `definitionRe`
(`^ {0,3}\[([^\]\n]+)\]:[ \t]+(\S.*)$`) validiert **nicht**, dass der
Rest der Zeile hinter der Zieladresse entweder leer/Whitespace ist oder ein
nach CommonMark gültig delimitiertes Titel-Suffix (`"…"`, `'…'`, `(…)`,
gefolgt von nur Whitespace bis Zeilenende) trägt. Sie akzeptiert **jeden**
nicht-leeren Rest als „Ziel + Titel" und schneidet mit `NormalizeTarget`
einfach am ersten Whitespace ab. Nach echtem CommonMark ist
„`[TERM]: First In, First Out`" **keine** gültige Link-Referenz-Definition
(der Rest nach dem Ziel ist weder leer noch ein gültig delimitiertes
Titel-Suffix) — ein CommonMark-konformer Parser würde die Zeile als
gewöhnlichen Prosa-Absatz lesen, in dem `[TERM]` nur literaler Text ist.
d-check erkennt sie trotzdem als Definition und meldet für das Wort „First"
(zufällig das erste Token nach dem Doppelpunkt) `target-missing`. Der
Code-Kommentar über `definitionRe` behauptet ausdrücklich „nach CommonMark"
— das trifft für diese Form nicht zu, und es ist keine der drei in
ADR-0093/0094 benannten Grenzen.

**Warum das über eine Doku-Ungenauigkeit hinausgeht:** Der Auftraggeber-Entscheid
zu diesem Slice aktiviert die Prüfung **standardmäßig, für jede Definition mit
Datei-Ziel**, in **jedem** Repo, das d-check einsetzt — nicht nur in diesem.
Eine Zeile der Form `[BEGRIFF]: Erklärungstext` ist ein in technischer Prosa
nicht seltenes Muster (Glossare, Abkürzungslisten), das mit dieser Erweiterung ab
sofort einen erfundenen `target-missing`-Befund erzeugen kann, dessen "Ziel" ein
zufälliges erstes Wort der Erklärung ist. Auf dem eigenen Bestand dieses Repos
tritt der Fall nicht auf (bestätigt: 107/107 Befunde vor/nach, siehe Gates-Tabelle
und §3 unten) — das schließt ihn für andere/künftige Dokumente nicht aus, exakt
dieselbe Argumentationsfigur wie das HIGH-Finding der Vorgänger-Review
(slice-232, R1-H1).

### 3. Bestandsmessung nachgefahren — bestätigt für den gemessenen Bestand

- Neues Image (`82c941ab`) gegen diesen Bestand, `--enable external --enable
  tracked`: `847 Datei(en) geprüft, 107 Befund(e)`.
- Altes Image (Stand `6f8d6906`, per `git worktree` + `docker build --target
  runtime`) gegen **denselben** Bestand, dieselben Module: `847 Datei(en)
  geprüft, 107 Befund(e)` — identisch. Die Behauptung „0 neue Befunde über alle
  sechs Konsumenten" ist für den tatsächlichen Bestand dieses Repos bestätigt,
  und zwar (anders als bei der Vorgänger-Review) **inklusive** `external` und
  `tracked` — diese beiden werden von `.d-check.yml` nicht standardmäßig
  aktiviert und wurden hier gezielt dazugenommen, weil R1 der Vorgänger-Review
  genau dort eine Lücke fand. Für diesen Slice besteht sie nicht: `external`/
  `tracked` iterieren unverändert generisch über `ExtractLinks` (Code bestätigt,
  §„Weitere Prüfpunkte"), eine neue Definition ändert an ihrem Verhalten nichts,
  was nicht schon durch die Ziel-Menge (Datei-Existenz/Git-Tracked-Status)
  erklärt wäre.
- Die Fehlfund-Klasse aus §2 (Finding H1) tritt auf dem aktuellen Bestand nicht
  auf — eine Zahlen-Messung allein hätte sie ohnehin nicht sichtbar gemacht,
  falls sie in einem der 847 Dokumente vorläge, ohne dass jemand gezielt danach
  sucht (dieselbe Einschränkung wie bei R1-M1 der Vorgänger-Review, dort für
  Fund-**Ort** statt Fund-**Existenz**).

---

## Weitere Prüfpunkte

- **`anchors`-Ausnahme (Code, nicht nur Test):** `CheckAnchors`
  (`internal/hexagon/core/rules/anchors.go:242`) hat genau **eine** Aufrufstelle
  von `ExtractLinks` im gesamten Modul; der `if ref.IsDefinition { continue }`
  steht als **erste** Anweisung der Schleife, vor jeder Auflösung
  (`resolveAnchorRef`). Es gibt keinen zweiten Pfad, der Definitionen doch an
  eine Anker-Prüfung heranträgt. Eigene Probe: Definition mit totem
  Fragment-Ziel auf **fehlender** Basisdatei (`[label]: fehlt.md#anker`) meldet
  genau **einen** Befund (`links`s `target-missing` mit dem vollen, fragment-
  tragenden Ziel-String) — `anchors` bleibt still, keine Doppelzählung über den
  `#`-Split.
- **Fünf generische Konsumenten:** `grep -n "IsDefinition\|ExtractLinks("` über
  `links.go`, `links_resolvefrom.go`, `external.go`, `tracked.go` zeigt je
  **eine** Aufrufstelle, keine Referenz auf `IsDefinition` — die ADR-Behauptung
  „kein Konsument außer `links`/`anchors` ändert Code" ist zutreffend.
- **`matrix`-Lineage-Match über `ref.Text` (Prüfpunkt 3 aus dem Auftrag):**
  `lineageExempt(supersedeValues, ref.Text, rel)` (`matrix.go:77`) nimmt für
  eine Definition jetzt auch deren **Label** als Kandidat in die
  Supersede-Ausnahme auf. Das ist konsistent mit der ADR-Entscheidung „`matrix`
  behandelt eine Definition wie jeden anderen `LinkRef`" — ein Inline-Link mit
  Linktext „ADR-0001" auf eine `Superseded`-Datei erhält dieselbe Ausnahme wie
  eine Definition mit Label „ADR-0001" auf dasselbe Ziel. Kein Bug, aber eine
  Konsequenz, die weder ADR-0093 noch die Spezifikation für `matrix` konkret
  ausbuchstabiert (sie nennt nur „Klassen-/Status-Regeln" pauschal) — nicht
  merge-blockierend, siehe Negativbefunde.
- **`NormalizeTarget` auf dem Definitions-Ziel (Prüfpunkt 4 aus dem Auftrag):**
  CommonMark erlaubt für Definitionen drei Titel-Delimiter (`"…"`, `'…'`,
  `(…)`), für Inline-Link-Titel nur die ersten beiden. `NormalizeTarget`
  parst den Titel gar nicht inhaltlich, sondern trennt nur am ersten
  Whitespace ab — das Ergebnis ist unabhängig vom tatsächlich verwendeten
  Delimiter korrekt (eigene Probe: `[label]: ziel.md (ein Titel)` liefert
  `Target = "ziel.md"`, genau wie mit `"…"`). Die engere Frage aus dem Auftrag
  hat also **keinen** Befund; die **breitere** Frage (wird überhaupt geprüft,
  ob der Rest ein gültiges Titel-Suffix **ist**) ist genau R1-H1.
- **Kommentar-Disziplin (§3.7):** alle neuen Kommentare in `markdown.go`
  (Feld-Kommentar, Regex-Kommentar, Funktions-Kommentar),
  `anchors.go` (`// ADR-0093: anchors prueft Referenz-Definitionen nicht`) und
  den beiden Testdateien tragen ausschließlich `ADR-0093`/`ADR-0094` als
  Herkunfts-Feld — keine Slice-Nummern, keine Review-Befund-Marker, keine
  Deliberation über Verworfenes. Sauber (im Unterschied zu R1-H2 der
  Vorgänger-Review, wo `// slice-232 …`-Präfixe gefunden wurden).
- **ADR-Immutabilität (§3.5):** `git show 1fad1b64 -- 0093-…md` ändert
  **ausschließlich** die `**Status:**`-Zeile und hängt eine Zeile an die
  `## Geschichte`-Tabelle an — der Kern (Kontext, Entscheidung, Alternativen,
  Konsequenzen, Fitness Function, Re-Evaluierungs-Trigger) bleibt byte-identisch.
  `make adr-check RANGE=6f8d6906..82c941ab` bestätigt das unabhängig (0
  Befunde).
- **ADR-0094 Form:** drei verglichene Alternativen mit Trade-offs, Fitness
  Function (Tabelle), `Re-Evaluierungs-Trigger` (eine Bedingung, explizit
  „permanent" sonst), `Schärft:` zeigt aufwärts auf `DC-FA-LINK-001`,
  `Supersedes: ADR-0093` gesetzt — vollständig.
- **ADR-Index:** `docs/plan/adr/README.md` trägt beide Zeilen (ADR-0093 auf
  `Superseded by ADR-0094`, neue ADR-0094-Zeile) im selben Commit wie der
  Status-Übergang.
- **MR-032/Decken-Regel:** Historie-Zeilen 0.93.0/0.93.1 nennen weder ADR- noch
  Slice-Nummer im Fließtext, `Verweis`-Spalte `—`; Spezifikation ebenso
  („Begründung in begleitender ADR" ohne Nummer).
- **§3.4/MR-006 Referenzrichtung:** `grep -n "ADR-009[34]\|slice-233"
  spec/lastenheft.md spec/spezifikation.md` liefert keinen Treffer im Körper.
- **Abgrenzung §1 „Ausdrücklich NICHT":** Verwendungs-Auflösung, `anchors` bei
  Definitionen, Zeilenumbruch-Form, `--repair`, Handbuch/README/CHANGELOG —
  keines davon im Diff berührt.
- **Commit-Zerlegung:** `82c941ab` (Feature) ändert ausschließlich
  `internal/`; `955a8dc1`/`1fad1b64` (Vertrag) ändern ausschließlich
  `spec/`/`docs/plan/adr/`, kein Code; `ae8b1a6c`/`64be56b3` sind reine
  Lifecycle-Moves/Beanspruchung. Sauber.

---

## Findings

| # | Kategorie | Quelle | Pfad | Befund | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| R1-H1 | HIGH | `DC-FA-LINK-001.a` Schritt 3 / Prüffrage 2 (Kern-Modul meldet falsch) | `internal/hexagon/core/rules/markdown.go:575` (`definitionRe`), `:581` (`parseDefinitionLine`) | `definitionRe` validiert nicht, dass der Zeilenrest hinter der Zieladresse entweder leer/Whitespace ist oder ein nach CommonMark gültig delimitiertes Titel-Suffix trägt — sie akzeptiert jeden nicht-leeren Rest. Eine glossar-artige Prosazeile wie `[TERM]: First In, First Out`, die nach CommonMark **keine** gültige Link-Referenz-Definition ist (kein gültiges Titel-Suffix, kein reines Whitespace nach dem Ziel), wird trotzdem als Definition erkannt; `NormalizeTarget` schneidet am ersten Whitespace ab und erzeugt ein erfundenes Ziel („First"), für das `links` `target-missing` meldet. Das ist eine vierte, nicht in ADR-0093/ADR-0094 oder der Spezifikation benannte Grenze — der Code-Kommentar über `definitionRe` behauptet „nach CommonMark", was für diese Form nicht zutrifft. Da die Erweiterung laut Slice-Auftrag „standardmäßig an" für jedes d-check-Repo ist, betrifft das nicht nur diesen Bestand: eine in technischer Prosa nicht seltene Zeilenform (Glossar-/Abkürzungs-Stil `[BEGRIFF]: Erklärung`) erzeugt ab sofort potenziell einen erfundenen Befund. Auf dem aktuellen 847-Datei-Bestand dieses Repos tritt der Fall nicht auf (Bestandsmessung 107/107 Befunde vor/nach), das schließt ihn für andere/künftige Dokumente aber nicht aus. | ja — Fixture `docs/glossary.md` mit Inhalt `[TERM]: First In, First Out`, `d-check` (kein Flag nötig) darauf ausführen: meldet `docs/glossary.md:1  First  target-missing`, obwohl die Zeile nach CommonMark keine Link-Referenz-Definition ist | definitionRe-title-syntax-nicht-validiert |
| R1-L1 | LOW | Doku-Präzision, `matrix.go:77` | `docs/plan/adr/0093-links-referenz-definitionen-gemeinsame-extraktion.md` (Abschnitt „Entscheidung", Aufzählung der Konsumenten-Behandlung) | Die ADR sagt für `matrix` pauschal „ihre Klassen-/Status-Regeln", ohne zu erwähnen, dass eine Definition über `ref.Text` (ihr Label) jetzt auch in die Supersede-Lineage-Ausnahme (`lineageExempt`) eingeht — dieselbe Ausnahme, die für Inline-Link-Text gilt. Kein fehlerhaftes Verhalten (konsistent mit „wie jeden anderen `LinkRef`"), aber eine Konsequenz, die beim Lesen der ADR nicht offensichtlich ist. | ja — `grep -n "ref.Text" internal/hexagon/core/rules/matrix.go` zeigt die Verwendung in `lineageExempt` | adr-konsequenz-nicht-ausbuchstabiert |

## Negativbefunde (geprüft, ohne Befund)

- **Gates:** `make gates` (zehn Glieder) und unabhängig `make adr-check
  RANGE=6f8d6906..82c941ab` grün, echte Ausgabe oben zitiert (Coverage
  94,30 %, `semgrep` 0/55, `doc-check` 847/0).
- **`anchors`-Ausnahme:** Code-Pfad bestätigt (einzige Aufrufstelle, `continue`
  vor jeder Auflösung); eigene Probe mit totem Fragment auf fehlender
  Basisdatei erzeugt genau einen (`links`-)Befund, keinen doppelten.
- **Fünf generische Konsumenten (`links`, `links.resolve-from`, `matrix`,
  `external`, `tracked`):** kein Code-Diff außer `links`/`anchors`, bestätigt
  per `grep` auf allen vier unveränderten Dateien.
- **Drei benannte Grenzen (Blockquote/Liste, Backslash, Zeilenumbruch):** je
  eigene Fixture gegen den laufenden Container getestet, alle drei halten;
  ADR-0094-Korrektur der Backslash-Grenze inhaltlich per PCRE-Nachvollzug der
  Regex-Semantik bestätigt.
- **Vier zusätzliche adversarielle Fälle** (verschachtelte Klammern ohne
  Backslash, Tabellenzeilen-Optik, mehrere Definitionen in Folge, Definition +
  Inline-Link auf derselben Zeile): alle außer dem Titel-Syntax-Fall (R1-H1)
  verhalten sich sicher bzw. wie erwartet.
- **Bestandsmessung, inklusive `external`/`tracked`:** 107/107 Befunde auf
  Alt- vs. Neu-Binary gegen denselben aktuellen Bestand — reproduziert exakt,
  keine der beiden von `.d-check.yml` nicht standardmäßig aktivierten Module
  zeigt eine Verhaltensänderung (anders als bei der Vorgänger-Review
  slice-232, wo genau hier eine echte Differenz auftrat).
- **Rot-Beleg gegen den alten Stand (DoD-Testbehauptung):** unabhängig
  nachgefahren — Alt-Image 0 Befunde, Neu-Image 1 `target-missing` auf
  demselben Fixture mit totem Definitions-Ziel.
- **`NormalizeTarget`/Titel-Delimiter-Frage (enge Fassung):** Ziel-Extraktion
  ist unabhängig vom verwendeten Delimiter (`"…"`, `(…)`) korrekt, da nur am
  ersten Whitespace getrennt wird — kein Befund für die enge Frage (siehe
  R1-H1 für die breitere).
- **ADR-Immutabilität (§3.5):** Diff von `1fad1b64` auf ADR-0093 betrifft
  ausschließlich Status-Feld und Geschichte-Anhang; `make adr-check` über die
  volle Range bestätigt unabhängig.
- **ADR-0094 Form (Modul 4):** drei Alternativen, Fitness Function,
  Re-Evaluierungs-Trigger, `Schärft:` aufwärts, `Supersedes`-Feld —
  vollständig; ADR-Index aktualisiert.
- **Kommentar-Disziplin (§3.7):** alle neuen Kommentare tragen ausschließlich
  `ADR-0093`/`ADR-0094` als Herkunfts-Feld, keine Slice-Nummern, keine
  Review-Marker.
- **§3.4/MR-006 Referenzrichtung:** kein ADR-/Slice-Token im Körper von
  Lastenheft oder Spezifikation.
- **MR-032 Decken-Regel:** Historie-Zeilen ohne ADR-/Slice-Nennung.
- **Commit-Zerlegung:** Feature-Commit nur `internal/`, Vertrags-Commits ohne
  Code, Lifecycle-Moves ohne Inhaltsänderung an den bewegten Dateien.
- **Hexagon-Import-Richtung (ADR-0005):** keine neuen Imports; Diff berührt
  keine Adapter-Schicht.
- **Gate-Suppression/Schwellen-Senkung ohne ADR (§3.2/§3.6):** keine
  `//nolint`, keine Schwellen-Änderung im Diff.
- **Netzzugriff außerhalb `external`:** keiner im Feature-Commit.
- **Zustandsfelder:** keine Roadmap-/Register-Zustandszeile inhaltlich
  verändert außer dem etablierten Lifecycle-Move-Begleiteffekt.
- **`DC-QA-04` Alt-Tool-Migrationsabdeckung:** nicht berührt.
- **Modul-Scan-Grenze (§3.8):** nicht einschlägig — `ExtractLinks` ist eine
  interne Parser-Funktion, kein scannendes Audit-Modul mit eigener
  Scan-Wurzel.
- **Provenance-Marker:** zwei `<!-- d-check:status-provenance -->`-Marker in
  ADR-0093/ADR-0094 zeigen jeweils nur, *wo* die Entscheidung/Verifikation
  entstand — keine getarnte Entscheidungsgrundlage.

## Kategorie-Summary

- HIGH: 1 (R1-H1)
- MEDIUM: 0
- LOW: 1 (R1-L1)
- INFO: 0

## Verdikt

**Merge-blockierend: ja, wegen R1-H1.**

Begründung: R1-H1 ist ein nachgewiesener (nicht nur hypothetischer)
Korrektheits-Fehlfund in der Kern-Extraktion, die diese Erweiterung
standardmäßig für jedes d-check-Repo aktiviert — ein CommonMark-invalides
Muster (`[TERM]: First In, First Out`) wird als gültige Definition erkannt
und erzeugt ein erfundenes Link-Ziel samt `target-missing`-Befund. Das ist
weder in ADR-0093 noch in ADR-0094 noch in der Spezifikation als Grenze
benannt, obwohl beide ADRs ausdrücklich „drei benannte Grenzen" als
vollständige Liste der Abweichungen von echtem CommonMark führen und die
zweite ADR eigens entstand, um eine dieser drei Grenzen nach eigener
Verifikation zu korrigieren — dieselbe Verifikations-Disziplin, die die
vierte Grenze hier nicht erfasst hat. Auf dem aktuellen Bestand dieses Repos
zeigt sich der Fehlfund nicht (Bestandsmessung 107/107 bestätigt), das ist
aber kein hinreichender Beleg, da die Erweiterung fremde/künftige Dokumente
außerhalb dieses Repos trifft. Vor Merge sollte entweder die Regex um eine
Prüfung des Titel-Suffixes (leer, oder gültig delimitiert und danach nur
Whitespace) geschärft werden, oder die Grenze muss explizit in ADR-0093/94
und der Spezifikation benannt und mit einem Negativtest belegt werden — die
im Auftraggeber-Entscheid gewählte „standardmäßig an"-Voreinstellung macht
diese Wahl zu einer bewussten Risikoabwägung, keiner stillen Lücke.
R1-L1 ist nicht merge-blockierend (Doku-Präzision, kein Fehlverhalten).
