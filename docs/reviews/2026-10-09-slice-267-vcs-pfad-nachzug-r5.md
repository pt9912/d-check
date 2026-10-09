# Review R5: slice-267, `vcs` lässt einen reinen Pfad-Nachzug durch

- **Review-Art:** Code. Geprüft wird `d9a6fa89..0be94741` — `0be94741` (`fix(vcs)`, Antwort auf R4
  H-1, M-1, I-3) — gegen den Slice-Plan `slice-267`, gegen die Findings aus R1 bis R4
  (`2026-10-09-slice-267-vcs-pfad-nachzug-r1.md` bis `-r4.md`), gegen `ADR-0103`,
  `DC-FA-VCS-001` und die Hard Rules `AGENTS.md` §3.5/§3.7 sowie §5 Regel 13 und 15. Die
  DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-267` · `d9a6fa89..0be94741`
- **Skill:** `reviewer.md` @ 1.19.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** `DC-FA-VCS-001` im Lastenheft 0.102.3 (Absatz „Pfad-Nachzug (opt-in)",
  Historie 0.102.3); Spezifikation `DC-FA-VCS-001.a` Schritt 4, §2-Zeile `vcs.ignore-link-targets`,
  §8-Historie; `ADR-0103`; `harness/sensors/adr-check.md` Grenze 4. Im Code:
  `vcs_link_targets.go` (neu), `vcsCore` in `vcs.go`; in `markdown.go` `PreprocessMarkdown`,
  `proseLines`/`TrimFenceIndent`, `proseParagraphs`, `forEachInlineCodeSpan`, `ExtractLinkSpans`,
  `parseLinkAt`, `matchBracket`, `definitionRe`; `vcs_pfad_nachzug_test.go`. Vorherige Findings
  am Modul: R1 bis R4 zu `slice-267`.
- **Proben** (Image per `make build` vom Stand `0be94741`, für die Läufe als eigener Tag geführt
  und danach entfernt; `d-check:latest` war schon vorhanden und bleibt). Je Probe ein
  Wegwerf-git-Repo im Scratchpad mit einer `Accepted`-ADR unter `## Kontext`, zwei Commits,
  `--enable vcs --range HEAD~1..HEAD` mit der `vcs`-Konfiguration aus `.d-check.yml`
  (`ignore-link-targets: true`). „Lesart" ist CommonMark/GFM. Alle Repos und der Klon sind
  entfernt, der Arbeitsbaum ist unverändert.

  | # | BASE → HEAD (nur die geänderte Stelle) | Ergebnis | Lesart |
  |---|---|---|---|
  | 00 | `Wir nehmen A.` → `B.` | Drift, Exit 1 | Kontrolle |
  | 01 | `[x](a.md)` → `[x](b/a.md)` | Exit 0 | Link, reiner Nachzug |
  | 02 | `Variante [B](nicht empfohlen) gilt.` → `(stark empfohlen)` | **Exit 0, still** | kein Link (Ziel mit Leerraum, kein Titel): sichtbarer Text |
  | 03 | `[B](Verboten "x) gilt.` → `(Erlaubt "x)` | **Exit 0, still** | kein Link (Titel offen): sichtbarer Text |
  | 04 | `[B](<Verboten> weil) gilt.` → `(<Erlaubt> weil)` | **Exit 0, still** | kein Link: sichtbarer Text |
  | 05 | `Wahl [a [b](c.md) d](Verboten) Ende.` → `(Erlaubt)` | **Exit 0, still** | innerer Link gewinnt, `](Verboten)` ist sichtbarer Text |
  | 06 | Tabellenzeile `\| x \| [L](Verboten\|ja) \|` → `(Erlaubt\|ja)` | **Exit 0, still** | GFM trennt die Zelle am `\|`: sichtbarer Text |
  | 07 | Backtick-Fence im Zitat, darin `> [a](Verboten)` → `(Erlaubt)` | Drift, Exit 1 | Fence im Zitat (zufällig als Code-Span geleert) |
  | 08 | `Text` / Leerzeile / `>␣␣␣␣␣[a](Verboten)` → `(Erlaubt)` | **Exit 0, still** | eingerückter Code im Zitat |
  | 09 | `-␣␣␣␣␣[a](Verboten)` → `(Erlaubt)` | **Exit 0, still** | eingerückter Code im Listenpunkt |
  | 10 | ``- Punkt mit ` Backtick`` / ``- `[a](Verboten)` Code`` → `(Erlaubt)` | **Exit 0, still** | Code-Span im zweiten Listenpunkt; der Absatz-Scanner paart den Backtick über die Blockgrenze |
  | 11 | `<!--` / Leerzeile / `[a](Verboten)` / `-->` → `(Erlaubt)` | Exit 0 | HTML-Kommentar, unsichtbar |
  | 12 | mit vier Leerzeichen eingerückter Backtick-Fence, eine Zeile `␣␣␣␣x`, Leerzeile, danach ein echter Backtick-Fence (Info `text`) mit `[a](Verboten)` darin → `(Erlaubt)` | **Exit 0, still** | Fenced Code; der eingerückte Fence verschiebt den Fence-Zustand |
  | 13 | `[x](a.md).\r\n` → `[x](b/a.md).\r\n` | Exit 0 | CRLF, reiner Nachzug |
  | 14 | `[x]: a.md\r\n` → `[x]: b/a.md\r\n` | Drift, Exit 1 | CRLF-Definition, reiner Nachzug (fail-safe) |
  | 15 | `[Regel]: gilt.nicht "x` → `gilt.immer "x` | Drift, Exit 1 | keine Definition |
  | 16 | `Text` / `[Status]: Abgelehnt.` → `Angenommen.` | Exit 0 | benannte Grenze |
  | 17 | `[x][alt]` + `[alt]: a.md` → `[x][neu]` + `[neu]: a.md` | Drift, Exit 1 | Label-Änderung |
  | 18 | `> ~~~` / `> [a](Verboten)` / `> ~~~` → `(Erlaubt)` | **Exit 0, still** | Fenced Code im Zitat |
  | 19 | Listenpunkt, Leerzeile, `␣␣~~~` / `␣␣[a](Verboten)` / `␣␣~~~` → `(Erlaubt)` | Drift, Exit 1 | Fence im Listenpunkt erkannt |
  | 22 | ``## Titel `[a](Verboten)` `` → `(Erlaubt)` | Drift, Exit 1 | Code-Span in Überschrift |
  | 23 | `Text <span>[a](Verboten)</span>` → `(Erlaubt)` | Exit 0 | Inline-HTML, Link bleibt Link |
  | 25 | `[B](Verboten\) gilt).` → `(Erlaubt\) gilt)` | Drift, Exit 1 | escapte schließende Klammer, Filter greift |

  **Bestand** (Wegwerf-Klon dieses Repos, Lauf mit `FOCUS_DISABLE` aus dem `Makefile`):

  | # | Änderung | Ergebnis |
  |---|---|---|
  | B1 | Anlass: die vier `releasing.md`-Ziele in `ADR-0014` und das in `ADR-0067` auf `docs/user/maintainer/releasing.md` | **0 Befunde, Exit 0** |
  | B2 | in den sieben `Accepted`-ADRs mit den meisten Links (`0024`–`0029`, `0072`) jedes relative Link-Ziel mit `pre/` präfigiert — 94 Zeilen, 80 Hunks vor `## Geschichte` | **0 Befunde, Exit 0** |
  | B3 | auf B2: je Datei ein Linktext (`…Z](pre/…`) geändert | **7 Befunde, je Datei einer**, Exit 1 |
  | B4 | dasselbe Präfix in **allen 95** `Accepted`-ADRs (653 Zeilen) | 3 Befunde: `ADR-0016`, `ADR-0058` (Links in Listen-Folgeabsätzen mit fünf Leerzeichen Einzug — benannte Grenze) und `ADR-0039` (die Ersetzung traf einen Link **in** einem Code-Span — echte Inhaltsänderung, richtig gemeldet) |
  | B5 | B4 ohne die eingerückten Zeilen und die Code-Span-Stelle | **0 Befunde, Exit 0** |

  Zählform der Bestandsmessung: `grep -o '](' <ADR> | wc -l` je Datei für die Auswahl; die
  geänderten Zeilen zählt `git diff --shortstat`. Ob die Bypass-Formen 02, 06, 08, 09, 12 im
  `Accepted`-Bestand vorkommen: `grep -noE '\]\([^() ]+ [^()]*\)'` (Ziel mit Leerraum ohne
  Titel-Form, Treffer mit Titel herausgefiltert), `'\]\([^() |]*\|[^() ]*\)'`,
  `'^>( {5,}|\s*(~~~|```))'`, `'^( {4,}|\t)(```|~~~)'`, `'^\s*([-*]|[0-9]+\.) {5,}\S'` — jeweils
  **0 Treffer**. Form 10 (einzelner Backtick in einem Nachbarblock ohne Leerzeile) ist nicht
  gemessen.

## Findings

### M-1 — Die Link-Erkennung von `links` liest Klammertext als Link, den Markdown als Text rendert; eine inhaltliche Umkehr passiert dort still

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 1 und 18 · `DC-FA-VCS-001` („Jede Änderung … an der übrigen Zeile … bleibt
  Drift") · `AGENTS.md` §5 Regel 13
- **Pfad:** `internal/hexagon/core/rules/vcs_link_targets.go` · „`if c, ok := targetToken(text,
  sp.TextEnd+2, sp.End-1); ok {`"; `internal/hexagon/core/rules/markdown.go` · „`destEnd, ok :=
  matchBracket(full, textEnd+1, '(', ')')`"; `spec/spezifikation.md` · „**Grenze:** Eine
  Absatz-Folgezeile in der Form einer Referenz-Definition wird geleert"
- **Befund:** `parseLinkAt` nimmt jede balancierte `[…](…)`-Folge als Link, ohne die
  Ziel-Syntax zu prüfen; `targetToken` leert davon das erste Token. Ist das Klammerpaar nach
  CommonMark/GFM kein Link — Ziel mit Leerraum ohne gültigen Titel, offener Titel, `<…>` mit Rest,
  äußere Klammer um einen inneren Link, `|` in einer Tabellenzelle —, ist dieses Token sichtbarer
  Text, und seine Änderung passiert mit dem Schlüssel still (Proben 02–06). Die Spezifikation nennt
  als einzige Grenze in Leerungs-Richtung die Absatz-Folgezeile in Referenz-Form; das
  Inline-Gegenstück, das dieselbe Art Fehler über eine viel häufigere Syntax trägt, steht in keiner
  der vier Grenzen-Stellen (Spezifikation Schritt 4, Lastenheft, `adr-check.md` Grenze 4,
  Code-Kommentar).
- **Failure-Szenario:** Eine künftige `Accepted`-ADR schreibt „Variante `[B](nicht empfohlen)`";
  ein späterer Commit macht daraus „Variante `[B](stark empfohlen)`". `make adr-check` (Hook und
  CI-Range) meldet nichts, obwohl das Lastenheft zusagt, dass jede Änderung an der übrigen Zeile
  Drift bleibt.
- **Warum nicht HIGH:** im `Accepted`-Bestand kein Vorkommen der Formen 02 und 06 (Zählform
  oben); für einen Bypass muss die Form schon in BASE stehen. Gleiche Stufe wie R4 M-1.
- **Verifizierbar:** ja — Proben 02–06; ein Kern-Test `"[B](nicht empfohlen)"` →
  `"[B](stark empfohlen)"` mit erwartetem Befund liefe rot.
- **Klasse:** `muster-leert-mehr-als-gegenstand`

### M-2 — „Code bleibt Teil des Vergleichs" gilt nicht für Code in Zitat und Listenpunkt und nicht hinter einem verschobenen Fence- oder Backtick-Zustand

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 1 und 18 · `DC-FA-VCS-001` („eine Fußnote, Code, eingerückte Zeilen,
  HTML-Blöcke … bleiben unverändert Teil des Vergleichs") · `AGENTS.md` §3.8, §5 Regel 13
- **Pfad:** `internal/hexagon/core/rules/vcs_link_targets.go` · „`if html[ln.No] ||
  indentColumns(r) >= 4 || len(ln.Text) != len(r) {`"; `internal/hexagon/core/rules/markdown.go`
  · „`return strings.TrimLeft(raw, " \t")`"; `spec/lastenheft.md` · „eine Fußnote, Code,
  eingerückte Zeilen, HTML-Blöcke und ein Link mit escapter"; `harness/sensors/adr-check.md` ·
  „eingerückte Zeilen, HTML-Blöcke, Code und Links mit escapter"
- **Befund:** Der Einzug-Filter misst die Spalte der Rohzeile; Code hinter einem Container-Präfix
  (`>␣␣␣␣␣…` im Zitat, `-␣␣␣␣␣…` im Listenpunkt) hat Spalte 0 und wird geleert (Proben 08, 09). Die
  geteilte Vorverarbeitung erkennt keinen Fence hinter `>` (Probe 18 mit `~~~`), schaltet den
  Fence-Zustand an einem eingerückten ```` ``` ```` um, der nach CommonMark Code-Inhalt ist (Probe
  12: der folgende echte Fenced-Block wird Prosa), und paart Backticks über Listen-/Tabellenzeilen
  ohne Leerzeile hinweg, sodass ein echter Code-Span ungeleert bleibt (Probe 10). In allen fünf
  Fällen ist `[a](Verboten)` sichtbarer Code und seine Änderung passiert still. Die Zusage
  „Code … bleibt unverändert Teil des Vergleichs" steht in Lastenheft und `adr-check.md` ohne
  Einschränkung; die Spezifikation sagt nur „dieselbe Vorverarbeitung (Fenced-Code entfällt …)".
  Für `links` sind diese Formen Grenzen der Erkennung, für `vcs` werden sie zu einem stillen Pfad —
  das ist die Frage von §3.8: dieselbe Eingabe, eine andere Zusage.
- **Failure-Szenario:** Eine `Accepted`-ADR zeigt ein Beispiel als Code in einem Zitat oder
  Listenpunkt (`> ~~~` … `> [Status](Abgelehnt)` …); ein Commit ändert das Beispiel inhaltlich, das
  Gate bleibt grün.
- **Warum nicht HIGH:** im Bestand kein Vorkommen von 08, 09, 12, 18 (Zählform oben); 10 nicht
  gemessen.
- **Verifizierbar:** ja — Proben 08, 09, 10, 12, 18; Kern-Tests mit diesen Formen und erwartetem
  Befund liefen rot.
- **Klasse:** `muster-leert-mehr-als-gegenstand`

### M-3 — Der Escape-Filter auf der schließenden runden Klammer hat keinen Test

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 13 · Spezifikation `DC-FA-VCS-001.a` Schritt 4 („ein Link, dessen öffnende
  oder schließende Klammer escapt ist, ist keiner")
- **Pfad:** `internal/hexagon/core/rules/vcs_link_targets.go` · „`escapedAt(text, sp.End-1)`";
  `internal/hexagon/core/rules/vcs_pfad_nachzug_test.go` · „`escapte schliessende klammer
  geaendert`"
- **Befund:** Die Bedingung prüft drei Stellen: `[`, `]` und `)`. Die Tests decken `[` („escapte
  oeffnende klammer") und `]` („escapte schliessende klammer", `[Ausnahme\](keine)`); für `)` gibt
  es keinen Fall, und das Bewusste Brechen in `0be94741` nennt nur „Escape der schließenden
  Klammer" — gemeint ist der `]`-Fall. Ohne den dritten Term liest `matchBracket` in
  `[B](Verboten\) gilt)` die escapte `)` als Ende, `targetToken` leert `Verboten\`, und eine
  Änderung von „Verboten" passiert still. Am Image ist der Filter heute wirksam (Probe 25: Drift).
- **Failure-Szenario:** Eine spätere Vereinfachung entfernt `escapedAt(text, sp.End-1)`; die
  Suite bleibt grün, der stille Pfad aus Probe 25 ist offen.
- **Verifizierbar:** ja — den Term entfernen, `make test` bleibt grün.
- **Klasse:** `filterzweig-ohne-negativtest`

### L-1 — Die Grenze von `htmlBlockLines` sagt „weiter als CommonMark"; für Kommentar, `<?`, `<!X` und CDATA ist sie enger

- **Kategorie:** LOW
- **Quelle:** `AGENTS.md` §5 Regel 13 (Grenze gegen den Gegenstand prüfen)
- **Pfad:** `internal/hexagon/core/rules/vcs_link_targets.go` · „`// GRENZE: weiter als CommonMark,
  das nicht jedes Tag als Blockanfang liest --`"
- **Befund:** Der Block endet an der nächsten Leerzeile, außer bei `pre`/`script`/`style`/`textarea`.
  CommonMark beendet einen HTML-Kommentar erst an `-->` (ebenso `?>`, `>`, `]]>`); nach einer
  Leerzeile im Kommentar gelten die Zeilen hier wieder als Prosa und werden geleert (Probe 11). Der
  Inhalt ist unsichtbar, deshalb kein Failure-Szenario über LOW; die Grenze behauptet aber nur die
  eine Richtung.
- **Verifizierbar:** ja — Probe 11.
- **Klasse:** `grenze-nur-eine-richtung`

### I-1 — Eine Referenz-Definition mit CRLF wird nie geleert

- **Kategorie:** INFO
- **Quelle:** Maintainability
- **Pfad:** `internal/hexagon/core/rules/markdown.go` · „`(?:"[^"\n]*"|'[^'\n]*'|\([^()\n]*\)))?[
  \t]*$`"
- **Befund:** `splitLines` trennt nur an `\n`; das `\r` am Zeilenende verhindert den Treffer von
  `definitionRe`, ein reiner Nachzug einer Definition bleibt Drift (Probe 14), ein Inline-Link mit
  CRLF wird geleert (Probe 13). Fail-safe; die Fehlschlag-Liste der Spezifikation nennt es nicht.
  Im ADR-Bestand keine Datei mit `\r`.
- **Verifizierbar:** ja — Probe 14.
- **Klasse:** `crlf-asymmetrie`

## Status der R4-Findings

- **R4 H-1:** aufgelöst. Ein Link mit Code-Span im Linktext ist wieder ein Link (Test „code-span im
  linktext nachgezogen"); der Anlass passiert am Image (B1), und der breite Nachzug über alle 95
  `Accepted`-ADRs ergibt nur Befunde an benannten Grenzen oder echten Inhaltsänderungen (B4, B5).
  Die Gegenrichtung hält: jede Linktext-Änderung wird gemeldet (B3).
- **R4 M-1:** für die gemeldeten Formen aufgelöst (Leerzeichen-Tab-Einzug, `<div>`, `<pre>` mit
  Leerzeile; Tests vorhanden). Weitere Code-Formen: M-2.
- **R4 I-1, I-2:** unberührt, bleiben INFO.
- **R4 I-3:** eingelöst in der Bauart (geteilte Erkennung statt eigener Muster) und in der
  Bestandsrichtung — die Bestandsprobe dieser Runde ist über den ganzen `Accepted`-Bestand grün.
  In der Umgehungsrichtung bringt die geteilte Erkennung ihre eigenen Abweichungen von Markdown mit
  (M-1, M-2); sie waren in `links` unschädlicher, weil dort ein falsch erkannter Link einen
  Link-Befund erzeugt und keinen Vergleich unterdrückt.

## Negativbefunde (geprüft, ohne Befund)

- **Anlass und Bestand:** B1–B5 oben. Kein Befund.
- **Filter Einzug, HTML-Block, Fußnote, Pfadzeichen:** je ein Test mit Erwartung „Drift"; am Image
  wirksam. Kein Befund.
- **`targetToken`/`cutRanges`:** Titel bleibt stehen (Test „titel geaendert"), `<…>`-Ziel mit
  Leerraum als Ganzes geleert (Probe 09 aus R4, Test „spitzklammer-ziel"); überlappende Bereiche
  (Definition und Inline-Link auf einer Zeile) werden übersprungen, nicht doppelt geschnitten.
  Kein Befund.
- **Positionstreue:** Schnitte werden auf der vorverarbeiteten Zeile berechnet und auf die Rohzeile
  angewandt; `stripInlineCodeByLine` ersetzt längengleich und lässt `\n` stehen. Die Bedingung
  `len(ln.Text) != len(r)` ist damit nie wahr — ein Schutz, kein Fehler. Kein Befund.
- **Escapes und Klammer-Balance:** escapte öffnende `[`, schließende `]` und `)` verwerfen den Link
  (Proben 25, Tests); `\\[` bleibt Link; escapte Klammern im Linktext und Ziel auf der Folgezeile
  bleiben Drift (fail-safe, benannt). Kein Befund über M-3 hinaus.
- **Kommentare §3.7:** `blankedLinkTargetLines`, `linkTargetCuts`, `targetToken`, `cutRanges`,
  `escapedAt`, `indentColumns`, `htmlBlockLines`, `vcsCore` und die Test-Kommentare tragen Zusage
  bzw. `GRENZE`; keine Review-Historie, keine Befund-Nummern, keine Slice-Nummern. Kein Befund zur
  Klasse (Wahrheit: L-1).
- **Lastenheft 0.102.3:** Versionskopf und Historie-Zeile passen, neueste oben; Verweis auf
  `DC-FA-LINK-001` löst auf. Kein Befund zur Form (Inhalt: M-2).
- **Spezifikation Schritt 4, §2-Zeile, §8-Historie:** beschreiben den Code (Vorverarbeitung,
  vier Filter, Titel bleibt) zutreffend; Historie ergänzt statt umgeschrieben. Kein Befund über
  M-1/M-2 hinaus.
- **Opt-in/byte-identisch:** `blanked` bleibt ohne Schlüssel `nil`, der Zweig ohne ihn ist
  unverändert; jeder Test fährt beide Modi. Kein Befund.
- **`ADR-0005`/Hexagon, `DC-QA-03`, Suppressions:** neue Datei im Kern nutzt nur `regexp`, `sort`,
  `strings`; kein Netz, kein `//nolint`. Kein Befund.
- **Botschaft `0be94741` (§5 Regel 15):** der Bestandsbeleg ist auf `ADR-0014` begrenzt und so
  benannt; „jeder nur in Richtung Drift" gilt für die Filter. Kein Befund (die Reichweite der
  Erkennung trägt M-1/M-2).
- **Traceability:** Betreff trägt `slice-267`, `DC-FA-VCS-001`, `ADR-0103`. Kein Befund.
- **Plan-Abgrenzung:** kein Umzug von `releasing.md`, kein Template-Feld-Nachzug. Kein Befund.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 (M-1, M-2, M-3) |
| LOW | 1 (L-1) |
| INFO | 1 (I-1) |

## Verdikt

**Nicht blockiert durch HIGH; drei MEDIUM vor Closure zu klären.** Der Umbau auf die geteilte
Erkennung löst R4 H-1: der Anlass passiert, und ein reiner Nachzug über alle 95 `Accepted`-ADRs
meldet nur an benannten Grenzen. slice-268 kann aus Sicht dieses Reviews anlaufen. In der
Umgehungsrichtung übernimmt `vcs` die Abweichungen der `links`-Erkennung von Markdown: Klammertext,
der kein Link ist (M-1), und Code hinter Container-Präfixen oder verschobenem Fence-/Backtick-Zustand
(M-2) werden geleert, obwohl sie sichtbar sind. Keine der Formen steht im `Accepted`-Bestand; die
Grenzen-Stellen (Lastenheft, Spezifikation, `adr-check.md`, Code-Kommentar) nennen sie nicht. Dazu
fehlt der Test für einen der drei Escape-Terme (M-3).
