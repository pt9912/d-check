# Review R6: slice-267, `vcs` lässt einen reinen Pfad-Nachzug durch

- **Review-Art:** Code. Geprüft wird `8a94cb86..0dc2242c`: `0dc2242c` (`fix(vcs)`, Antwort auf R5
  M-1, M-2, M-3, L-1 und I-1), dazu die drei Form-Änderungen am R5-Report in `8a94cb86`. Geprüft
  gegen den Slice-Plan `slice-267`, die Findings aus R1 bis R5
  (`2026-10-09-slice-267-vcs-pfad-nachzug-r1.md` bis `-r5.md`), `ADR-0103`, `DC-FA-VCS-001` und
  die Hard Rules `AGENTS.md` §3.7, §3.8 sowie §5 Regel 13 und 15. Die DoD-Abhakung prüft dieses
  Review nicht.
- **Gegenstand:** `slice-267` · `8a94cb86..0dc2242c`
- **Skill:** `reviewer.md` @ 1.19.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** `DC-FA-VCS-001` im Lastenheft 0.102.4 (Absatz „Pfad-Nachzug (opt-in)",
  Historie 0.102.4), Spezifikation `DC-FA-VCS-001.a` Schritt 4 und §8-Historie, `ADR-0103`,
  `harness/sensors/adr-check.md` Grenze 4. Im Code: `vcs_link_targets.go` (`opaqueLines`,
  `strictFenceLines`, `lineCodeSpans`, `linkDestRE`, `listCodeRE`, `htmlBlockEnd`,
  `linkTargetCuts`); aus `markdown.go` `PreprocessMarkdown`, `proseLines`, `FenceRun`,
  `forEachInlineCodeSpan`, `ExtractLinkSpans`, `parseLinkAt`, `matchBracket` und `definitionRe`;
  dazu `vcs_pfad_nachzug_test.go`. Vorherige Findings am Modul: R1 bis R5 zu `slice-267`.
- **Proben:** Das Image ist per `make build IMAGE=d-check-r6` vom Stand `0dc2242c` gebaut. Für jede
  Probe gibt es ein Wegwerf-git-Repo im Scratchpad: eine `Accepted`-ADR, die geänderte Stelle unter
  `## Kontext`, zwei Commits. Der Lauf ist `--enable vcs` mit `FOCUS_DISABLE` aus dem `Makefile`,
  `--range HEAD~1..HEAD` und dem `vcs`-Block aus `.d-check.yml` (`ignore-link-targets: true`).
  „Lesart" heißt CommonMark/GFM. Mehrzeilige Formen stehen mit ` / ` als Zeilentrenner, `␣` steht
  für ein Leerzeichen.

  | # | BASE (HEAD ersetzt `Verboten` durch `Erlaubt`, wo nicht anders genannt) | Ergebnis | Lesart |
  |---|---|---|---|
  | 00 | `Wir nehmen A.` → `B.` | Drift, Exit 1 | Kontrolle |
  | 01 | Link-Ziel `a.md` → `b/a.md` | Exit 0 | reiner Nachzug |
  | 02 | `Wahl [a [b](c.md) d](Verboten) Ende.` | **Exit 0, still** | äußere Klammer ist kein Link, `(Verboten)` ist sichtbarer Text (= R5-Probe 05) |
  | 03 | Tabellenzeile `\| z \| [L](Verboten\|ja) \|` | **Exit 0, still** | GFM trennt die Zelle am `\|`, sichtbarer Text (= R5-Probe 06) |
  | 04 | `[B](a.md\(x)Verboten) gilt.` | **Exit 0, still** | `\(` zählt nicht zur Balance; Link endet nach `x`, `Verboten)` ist sichtbarer Text |
  | 05 | `[a\[b]c](Verboten) gilt.` | **Exit 0, still** | `\[` zählt nicht; `[a\[b]` gefolgt von `c` ist kein Link, sichtbarer Text |
  | 06 | `- ~~~` / `␣␣[a](Verboten)` / `␣␣~~~` | **Exit 0, still** | Fenced Code im Listenpunkt, Öffner auf der Markenzeile |
  | 07 | `1. ` + Backtick-Fence / `␣␣␣[a](Verboten)` / `␣␣␣` + Backtick-Fence | **Exit 0, still** | wie 06, geordnete Liste |
  | 08 | `- <pre>` / `␣␣[a](Verboten)` / `␣␣</pre>` | **Exit 0, still** | HTML-Block im Listenpunkt, Inhalt sichtbar |
  | 09 | `- <div>` / `␣␣[a](Verboten)` / `␣␣</div>` | **Exit 0, still** | wie 08, Typ 6 |
  | 10 | `- -␣␣␣␣␣[a](Verboten)` | **Exit 0, still** | eingerückter Code in verschachtelter Liste |
  | 11 | `- >␣␣␣␣␣[a](Verboten)` | **Exit 0, still** | eingerückter Code im Zitat im Listenpunkt |
  | 12 | `Siehe <https://e.org/[a](Verboten)>.` | **Exit 0, still** | Autolink, die URL ist sichtbarer Linktext |
  | 13 | `[x](a.md).` + CRLF → `[x](b/a.md).` + CRLF | Exit 0 | Inline-Link mit CRLF **wird** geleert |
  | 14 | `[x]: a.md` + CRLF → `[x]: b/a.md` + CRLF | Drift, Exit 1 | Definition mit CRLF bleibt Drift |
  | 16 | `Text <span title="[a](Verboten)">x</span>` | Exit 0 | Attributwert, nur als Tooltip sichtbar |
  | 17 | Bild im Linktext, Nachzug des **Bild**-Ziels | Drift, Exit 1 | fail-safe, das innere Ziel wird nie gelesen |
  | 28 | Definition `[r]: a.md` mit Titel `"[a](Verboten)"` | **Exit 0, still** | Titeländerung; die Linkerkennung läuft über die ganze Definitionszeile |
  | 29 | Titel `"Verb)oten"`, Ziel nachgezogen | Drift, Exit 1 | fail-safe |
  | 30 | Titel `"nicht"` → `"doch"` | Drift, Exit 1 | Titel bleibt Vergleich |
  | 40 | Vier-Backtick-Fence, darin eine Drei-Backtick-Zeile und `[a](Verboten)` | Drift, Exit 1 | der strenge Automat hält |
  | 41 | `~~~`-Fence, darin eine Backtick-Fence-Zeile und `[a](Verboten)` | Drift, Exit 1 | der strenge Automat hält |
  | 42 | `<?x` / Leerzeile / `[a](Verboten)` / `?>` | Drift, Exit 1 | Typ-3-Ende hält |
  | 43 | Zitat, reiner Nachzug | Drift, Exit 1 | fail-safe, benannt |
  | 44 | `-␣␣␣␣␣[a](Verboten)` | Drift, Exit 1 | R5-Probe 09 aufgelöst |
  | 45 | Definition mit Titel, reiner Nachzug | Exit 0 | Titel bleibt stehen |
  | 50 | Listenpunkt `- ` + Backtick + `a` / Listenpunkt `- x` + Backtick + ` [b](y` + Backtick + `z` + Backtick + `)`, `y` geändert | Drift, Exit 1 | versetzte Paarung, `lineCodeSpans` deckt die `[` |
  | 51 | `> Text` / `[a](a.md)`, reiner Nachzug | Exit 0 | Lazy Continuation, echter Link |
  | 53 | Tabelle mit Kopf- und Trennzeile, Zelle `[B](abgelehnt\|v2)` → `(angenommen\|v2)` | **Exit 0, still** | Failure-Szenario zu M-1 |

  **Bestand** (Wegwerf-Klon dieses Repos, derselbe Lauf auf dem Klon):

  | # | Änderung | Ergebnis |
  |---|---|---|
  | B1 | Anlass: die vier `releasing.md`-Ziele in `ADR-0014` und das in `ADR-0067` auf `docs/user/maintainer/releasing.md` | **0 Befunde, Exit 0** |
  | B2 | in **allen 96** `Accepted`-ADRs jedes Link-Ziel ohne `#`, `:` und Leerraum mit `pre/` präfigiert (`sed -E 's/\]\(([^)#: <][^): ]*)\)/](pre\/\1)/g'`, 1174 Zeilen samt `## Geschichte`) | 5 Befunde: `ADR-0016`, `ADR-0028`, `ADR-0039`, `ADR-0040`, `ADR-0058` |
  | B2a | B2 zeilenweise zerlegt (je geänderte Zeile vor `## Geschichte` ein eigenes Commit-Paar über der BASE-Fassung) | `ADR-0016`:44, `ADR-0040`:75, `ADR-0058`:99 und :127: Listen-Folgeabsatz mit sechs Spalten Einzug (benannte Grenze). `ADR-0039` (6 Zeilen) und `ADR-0040`:39: die Ersetzung traf Code-Span bzw. Fenced Code, also eine echte Inhaltsänderung, richtig gemeldet. **`ADR-0028`:85: neu, siehe L-1** |
  | B3 | auf B2: je ein Linktext in `ADR-0027` und `ADR-0028` geändert | **2 Befunde**, je Datei einer |

  **Wie viele echte Nachzüge die neuen Filter zu Drift machen.** Zitat und Listen-Code: **0**. Die
  `Accepted`-ADRs haben vor `## Geschichte` 17 Zitatzeilen (`ADR-0063`, `-0077`, `-0081`, `-0082`,
  `-0085`), keine davon trägt `](`. Eine Zeile mit `listCodeRE`-Form und Link gibt es nicht
  (`awk` über die Zeilen vor `## Geschichte`: Zeilen mit `](` **und** führendem `>` oder fünf
  Leerzeichen bzw. Tab hinter einer Listenmarke). Der zeilenlokale Code-Span-Filter macht **einen**
  Nachzug zu Drift (`ADR-0028`:85, L-1). Strenger Fence-Automat und HTML-Endmarker: 0, denn B2a
  findet keine weitere Zeile.

  **Die Bypass-Formen im Bestand.** Gezählt mit `grep -cE` über alle 96 `Accepted`-ADRs: verschachtelter
  Link, `|` im Ziel einer Tabellenzeile, `\(` im Ziel, `\[` im Linktext, Fence oder `<` auf einer
  Listenmarkenzeile, verschachtelte Listen- bzw. Zitatmarke mit fünf Leerzeichen, Autolink mit `](`,
  Definitionstitel mit `](`. Jede Form ergibt **0 Treffer**. Die Muster sind Näherungen, Form 02
  zum Beispiel nur mit einer Verschachtelungsebene.

  **Mutation** (eigener Wegwerf-Klon, `make test IMAGE=d-check-r6m`): Wird der Escape-Term
  `escapedAt(text, sp.End-1)` entfernt, wird „escapte zielklammer geaendert" rot (R5 M-3
  aufgelöst). Wird die Schließer-Bedingung von `strictFenceLines` auf
  `(c == char || true) && (run >= n || run >= 3)` gelockert, bleibt `make test` **grün** (M-3).

  Images (`d-check-r6:latest`, `d-check-r6m:*`), Probe-Repos und Klone sind entfernt. Der
  Arbeitsbaum ist bis auf diesen Report unverändert.

## Findings

### M-1 — Klammertext ohne Link wird weiter geleert, wo das Ziel syntaktisch gültig aussieht; zwei R5-Proben sind unverändert offen

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 1 und 18 · `DC-FA-VCS-001` („Jede Änderung … an der übrigen Zeile, am Titel
  eines Links … bleibt Drift") · `AGENTS.md` §5 Regel 13
- **Pfad:** `internal/hexagon/core/rules/vcs_link_targets.go` · „`if
  !linkDestRE.MatchString(text[sp.TextEnd+2 : sp.End-1]) {`" und „`for _, sp := range
  ExtractLinkSpans(text) {`"; `internal/hexagon/core/rules/markdown.go` · „`destEnd, ok :=
  matchBracket(full, textEnd+1, '(', ')')`"; `internal/hexagon/core/rules/vcs_pfad_nachzug_test.go`
  · „`tabellenzelle mit pipe geaendert`"
- **Befund:** Der neue Filter prüft nur die Form des Zielausdrucks. Er prüft nicht, ob das
  Klammerpaar überhaupt ein Link ist. Still passieren deshalb: ein äußerer Klammertext um einen
  inneren Link (02, R5-Probe 05), ein `|` ohne Leerraum im Ziel einer Tabellenzeile (03, R5-Probe 06),
  ein `\(` im Ziel und ein `\[` im Linktext (04, 05), denn `matchBracket` kennt keine Escapes und
  `escapedAt` prüft nur die drei Randklammern. Dazu kommen ein Link-förmiger Text in einer
  Autolink-URL (12) und ein Link-förmiger Text im **Titel** einer Referenz-Definition (28), weil
  `ExtractLinkSpans` über die ganze Definitionszeile läuft. Der Test „tabellenzelle mit pipe
  geaendert" prüft eine Form mit Leerraum. Sie scheitert schon am `\s` von `linkDestRE`, nicht am
  `|`, deshalb deckt der Test die R5-Probe 06 nicht ab.
- **Failure-Szenario:** Eine künftige `Accepted`-ADR enthält die Tabellenzeile
  `| Variante | [B](abgelehnt|v2) |`. Ein Commit macht daraus `[B](angenommen|v2)`. `make adr-check`
  bleibt grün, obwohl die Zelle sichtbar „abgelehnt" gegen „angenommen" tauscht (Probe 53, mit
  Kopf- und Trennzeile: Exit 0). Ebenso eine
  Titeländerung an einer Referenz-Definition, obwohl das Lastenheft den Titel ausdrücklich im
  Vergleich hält.
- **Warum nicht HIGH:** Im `Accepted`-Bestand gibt es keine dieser Formen (Zählform oben). Für einen
  Bypass muss die Form schon in BASE stehen. Gleiche Stufe wie R5 M-1.
- **Verifizierbar:** ja. Proben 02 bis 05, 12 und 28. Ein Kern-Test je Form mit erwartetem Befund
  liefe heute rot.
- **Klasse:** `muster-leert-mehr-als-gegenstand`

### M-2 — Code und HTML hinter einer Listenmarke bleiben ein stiller Pfad

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 1, 15 und 18 · `DC-FA-VCS-001` („Code, eingerückte Zeilen, Zitate,
  HTML-Blöcke … bleiben unverändert Teil des Vergleichs") · `AGENTS.md` §3.8, §5 Regel 13
- **Pfad:** `internal/hexagon/core/rules/vcs_link_targets.go` · „`if indentColumns(l) >= 4 ||
  listCodeRE.MatchString(l) ||`", „`if lead <= 3 && run >= 3 && (c == '~' ||`" und
  „`htmlBlockStartRE = regexp.MustCompile(`"; `spec/spezifikation.md` · „Fenced-Code nach einem
  strengen Automaten"
- **Befund:** Alle drei Zeilen-Filter lesen die Rohzeile ab Spalte 0. Steht das Block-Konstrukt
  hinter einer Listenmarke, greift keiner von ihnen. Das gilt für einen Fence-Öffner auf der
  Markenzeile (06, 07), einen HTML-Block auf der Markenzeile (08, 09) und eingerückten Code hinter
  einer zweiten Marke, Listen- oder Zitatmarke (10, 11). In allen sechs Fällen ist `[a](Verboten)`
  sichtbarer Code oder sichtbarer HTML-Text, und seine Änderung passiert still. Die Formen aus R5
  M-2 sind aufgelöst (Proben 43, 44, Tests). Die Klasse dahinter bleibt offen: ein Container-Präfix
  vor einem Block, der den Link zu Text macht. Spezifikation, Lastenheft, `adr-check.md` und der
  Kommentar zu `opaqueLines` sagen „Fenced-Code", „HTML-Blöcke" und „eingerückter Code … auch
  hinter einer Listenmarke" ohne diese Einschränkung.
- **Failure-Szenario:** Eine `Accepted`-ADR zeigt ein Konfigurationsbeispiel als Fenced Code direkt
  im Listenpunkt (`- ~~~` …). Ein Commit ändert darin eine Link-förmige Zeile inhaltlich, und das
  Gate bleibt grün.
- **Warum nicht HIGH:** Im Bestand gibt es kein Vorkommen (Zählform oben).
- **Verifizierbar:** ja. Proben 06 bis 11. Kern-Tests mit diesen Formen und erwartetem Befund liefen
  heute rot.
- **Klasse:** `muster-leert-mehr-als-gegenstand`

### M-3 — Die Schließer-Bedingungen des strengen Fence-Automaten haben keinen Test

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 13 · Spezifikation `DC-FA-VCS-001.a` Schritt 4 („Schließer aus demselben
  Zeichen in mindestens derselben Zahl")
- **Pfad:** `internal/hexagon/core/rules/vcs_link_targets.go` · „`if lead <= 3 && c == char && run
  >= n && run == len(trimmed) {`"
- **Befund:** Die Tests decken nur den Öffner ab („eingerueckter fence verschiebt nichts"). Gleiches
  Zeichen und Mindestlänge des Schließers sind die Bedingungen, an denen der strenge Automat vom
  Automaten der Vorverarbeitung abweicht. Diese schaltet an **jeder** Fence-Zeile um. Beide
  Bedingungen sind ungetestet: Die Mutation `(c == char || true) && (run >= n || run >= 3)` lässt
  `make test` grün. Am Image hält der Filter heute (Proben 40, 41).
- **Failure-Szenario:** Eine Vereinfachung des Schließer-Tests geht durch die Suite. Danach passiert
  die Änderung an `[a](Verboten)` in einem Vier-Backtick-Block, der eine Drei-Backtick-Zeile enthält
  (Probe 40), still. Das ist sichtbarer Code.
- **Verifizierbar:** ja. Siehe die Mutation oben.
- **Klasse:** `filterzweig-ohne-negativtest`

### L-1 — Der zeilenlokale Code-Span-Filter liest einen mehrzeiligen Span falsch und macht einen echten Nachzug im Bestand zu Drift; die Grenze ist nicht benannt

- **Kategorie:** LOW
- **Quelle:** `AGENTS.md` §5 Regel 13 (Grenze gegen den Gegenstand prüfen)
- **Pfad:** `internal/hexagon/core/rules/vcs_link_targets.go` · „`// lineCodeSpans liefert die
  Code-Spans, die vollständig in der Zeile liegen.`"; `docs/plan/adr/0028-planning-lifecycle-modul.md`
  · „``planning`, hermetisch **ohne** `--range`; [`DC-FA-CLI-010`]``"
- **Befund:** `lineCodeSpans` scannt die Rohzeile für sich allein. Die Zeile 85 von `ADR-0028` beginnt
  mit dem schließenden Backtick eines Spans, der auf Zeile 84 öffnet. Der Scanner paart deshalb
  versetzt, legt einen Schein-Span über die `[` von `DC-FA-CLI-010`, und der reine Nachzug dieses
  Ziels wird Drift (B2a). Die Richtung ist fail-safe. Die Fail-safe-Liste in Spezifikation,
  `adr-check.md` und Kommentar nennt diese Form aber nicht, und der Kommentar beschreibt die
  Funktion so, als fände sie die Spans der Zeile.
- **Verifizierbar:** ja. B2a: nur Zeile 85 von `ADR-0028` nachgezogen ergibt Exit 1.
- **Klasse:** `grenze-nur-eine-richtung`

### L-2 — Die CRLF-Grenze stimmt nur für Definitionen; ein Inline-Link auf einer CRLF-Zeile wird geleert

- **Kategorie:** LOW
- **Quelle:** `AGENTS.md` §5 Regel 13
- **Pfad:** `internal/hexagon/core/rules/vcs_link_targets.go` · „`// GRENZE: Ein Ziel auf der
  Folgezeile und eine Zeile mit CRLF-Ende werden`"; `spec/spezifikation.md` · „beginnt (sie gilt
  als HTML-Block), eine Zeile mit CRLF-Ende."; `harness/sensors/adr-check.md` · „Inline-HTML
  beginnt, CRLF-Ende),"
- **Befund:** Alle drei Stellen sagen, eine Zeile mit CRLF-Ende werde nicht geleert. Das stimmt für
  die Referenz-Definition, weil `definitionRe` am `\r` scheitert (Probe 14). Ein Inline-Link auf
  einer CRLF-Zeile wird dagegen geleert, weil das `\r` hinter der `)` liegt und den Zielausdruck nie
  erreicht (Probe 13). Die Grenze verspricht also mehr Drift, als der Code liefert. Einen stillen
  Pfad gibt es nicht, weil der geleerte Teil ein echtes Ziel ist. R5 I-1 ist damit benannt, aber
  zu breit.
- **Verifizierbar:** ja. Proben 13 und 14.
- **Klasse:** `grenze-gegen-gegenstand`

### I-1 — Dieselbe Finding-Klasse in der vierten Runde: Steering-Loop-Signal

- **Kategorie:** INFO
- **Quelle:** dieser Skill, §Kontext-Eskalation („Die dritte Wiederholung derselben Klasse in einer
  Sitzung ist ein Steering-Loop-Signal")
- **Pfad:** `spec/lastenheft.md` · „gültiges Ziel und ein Link mit escapter Klammer bleiben
  unverändert Teil des"
- **Befund:** `muster-leert-mehr-als-gegenstand` trägt R4 M-1, R5 M-1 und M-2 und jetzt R6 M-1 und
  M-2. Jede Runde löst die gemeldeten Formen auf, und die nächste Probe-Runde findet Nachbarformen
  derselben Klasse. Das Lastenheft sagt absolut „jede Änderung … bleibt Drift", die Erkennung bleibt
  eine Näherung an CommonMark/GFM. Ob die Zusage so gehalten werden kann oder ihre Reichweite
  benannt werden muss, entscheidet das Review nicht. Das ist eine Frage an Architect und `ADR-0103`.
- **Verifizierbar:** nein (Urteil).
- **Klasse:** `steering-loop-signal`

## Status der R5-Findings

- **R5 M-1:** teilweise aufgelöst. Die Formen 02 bis 04 aus R5 (Ziel mit Leerraum, offener Titel,
  `<…>` mit Rest) sind Drift und durch Tests gedeckt. Die Formen 05 (verschachtelt) und 06 (`|` ohne
  Leerraum) passieren weiter still. Weiter in M-1.
- **R5 M-2:** für die gemeldeten Formen aufgelöst: Zitat-Code, Listen-Code mit fünf Leerzeichen,
  `~~~` im Zitat, eingerückter Fence, Backtick-Paarung über Listenzeilen. Tests sind vorhanden, das
  Bewusste Brechen hat die Commit-Botschaft benannt. Nachbarformen stehen in M-2.
- **R5 M-3:** aufgelöst, Mutation rot.
- **R5 L-1:** aufgelöst für den Kommentar (Probe 42 und Test „html-kommentar mit leerzeile"). Die
  Enden der Typen `<?`, `<![CDATA[` und `<!X` haben keinen eigenen Test. Ihr Inhalt ist unsichtbar,
  deshalb hier kein Finding.
- **R5 I-1:** benannt, aber zu breit (L-2).

## Form-Änderungen am R5-Report (`8a94cb86`)

Die Zeilen 05, 06 und 10 stehen jetzt in Code-Spans, Zeile 12 ist in Worte gefasst. Die
Ursprungsfassung liegt weder in git noch im Scratchpad (`r5.md` dort ist die committete Fassung),
ein Byte-Vergleich ist deshalb nicht möglich. Geprüft ist, ob der Text in sich stimmig ist und ob
die Formen sich reproduzieren lassen. Die Formen 05 und 06 laufen am neuen Image als Proben 02 und
03 und ergeben dasselbe Verhalten wie in R5 beschrieben. Die Formen 10 und 12 decken sich mit der
Beschreibung in R5 M-2 („paart den Backtick über die Blockgrenze", „eingerückten ``` um") und mit
den neuen Tests „backticks ueber listenzeilen geaendert" und „eingerueckter fence verschiebt
nichts". Kein Widerspruch zum Rest des Reports. Kein Befund.

## Negativbefunde (geprüft, ohne Befund)

- **Anlass und Gegenrichtung:** B1 ergibt 0 Befunde, eine Linktext-Änderung wird gemeldet (B3). Kein
  Befund.
- **`linkDestRE` gegen `definitionRe`:** Gleiche Bauart, Titel in drei Formen. Ein escapter Quote im
  Titel bricht den Treffer, das ist fail-safe. Kein Befund.
- **`overlapsAny` nur auf der öffnenden Klammer:** Ein Code-Span **im** Ziel ist in der
  Vorverarbeitung geleert. `targetToken` endet am ersten Leerraum und schneidet ihn nie mit, der
  Rohtext des Spans bleibt also Vergleich. Bei versetzter Paarung aus einem Nachbarblock deckt
  `lineCodeSpans` die `[` (Probe 50: Listenpunkt mit offenem Backtick, danach Listenpunkt mit
  versetzter Paarung bis ins Ziel; Änderung im Ziel ergibt Drift). Kein stiller Pfad gefunden.
- **`strictFenceLines`:** Öffner mit Backtick in der Info-Zeile wird abgelehnt, wie es `FenceToggle`
  tut. Ein Tab-Einzug ist schon über `indentColumns` opak. Ohne Schließer reicht der Block bis zum
  Dateiende, das ist fail-safe. Kein Befund über M-3 hinaus.
- **`htmlBlockEnd`:** Der Endmarker wird auch auf der Startzeile geprüft (Test „link nach
  einzeiligem kommentar nachgezogen"). Ein Block ohne Endmarker ist bis zum Dateiende opak, das ist
  fail-safe. Kein Befund.
- **Zitat-Filter:** Eine Lazy-Continuation-Zeile ohne `>` wird geleert. Sie trägt einen echten Link,
  denn Fence und eingerückter Code sind als Lazy Continuation nicht möglich (Probe 51: Exit 0 bei
  reinem Nachzug). Kein Befund.
- **Kommentare §3.7:** `opaqueLines`, `strictFenceLines`, `lineCodeSpans`, `linkDestRE`,
  `listCodeRE`, `overlapsAny`, `htmlBlockEnd`, `htmlBlockLines` und die neuen Test-Kommentare tragen
  Zusage, Abgrenzung oder `GRENZE`. Keine Review-Historie, keine Befund- oder Slice-Nummern. Zur
  Klasse kein Befund, zur Wahrheit siehe L-1 und L-2.
- **Lastenheft 0.102.4:** Versionskopf und Historie-Zeile passen, die neueste steht oben, der Anker
  löst auf. Zur Form kein Befund, zum Inhalt siehe M-1 und M-2.
- **Spezifikation Schritt 4 und §8-Historie:** Die Historie ist ergänzt, nicht umgeschrieben. Die
  Filterbeschreibung entspricht dem Code: Listenmarke mit fünf Leerzeichen oder Tab, Zitat,
  strenger Automat, HTML-Endmarker, Zielausdruck. Kein Befund über M-1, M-2, L-1 und L-2 hinaus.
- **`ADR-0005`/Hexagon, `DC-QA-03`, Suppressions:** Nur `regexp`, `sort` und `strings`, kein Netz,
  kein `//nolint`. Kein Befund.
- **Opt-in:** Ohne Schlüssel bleibt der Zweig unverändert, jeder Test fährt beide Modi. Kein Befund.
- **Botschaft `0dc2242c` (§5 Regel 15):** Der Bestandsbeleg ist auf `ADR-0014` begrenzt und so
  benannt. Das Bewusste Brechen nennt sieben Filter. Reproduziert ist es für den Escape-Term, die
  sechs anderen sind nicht nachgefahren. Die Botschaft behauptet nicht, dass R5-Probe 05 und 06
  aufgelöst seien. Kein Befund.
- **Traceability:** Der Betreff trägt `slice-267`, `DC-FA-VCS-001` und `ADR-0103`. Kein Befund.
- **Plan-Abgrenzung:** Kein Umzug von `releasing.md`, kein Nachzug am Template-Feld. Kein Befund.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 (M-1, M-2, M-3) |
| LOW | 2 (L-1, L-2) |
| INFO | 1 (I-1) |

## Verdikt

**Kein HIGH. Drei MEDIUM sind vor Closure zu klären.** Der Anlassfall passiert, und ein breiter
Präfix-Nachzug über alle 96 `Accepted`-ADRs meldet nur an benannten Grenzen, an echten
Inhaltsänderungen und an einer neuen fail-safe-Stelle (L-1). Die neuen Filter Zitat und Listen-Code
kosten am Bestand keinen einzigen Nachzug. slice-268 kann aus Sicht dieses Reviews anlaufen. In der
Umgehungsrichtung bleibt dieselbe Klasse offen wie in R4 und R5. Klammertext, der kein Link ist,
dessen Ziel aber gültig aussieht, wird geleert, darunter zwei unverändert offene R5-Proben (M-1).
Ebenso Block-Code und HTML hinter einer Listenmarke (M-2). Keine dieser Formen steht im Bestand.
Die Wiederholung der Klasse (I-1) gehört vor einer weiteren Filterrunde an den Architect.
