# ADR-0091: Die gemeinsame `[]Line`-Link-Extraktion erkennt einen Zeilenumbruch hinter `](` absatzweise — die string-basierte Extraktion bleibt zeilenbasiert

**Status:** Superseded by ADR-0092

**Datum:** 2026-09-27

**Autor:** claude-sonnet-5

**Bezug:** Change Request eines Konsumenten; Auftraggeber-Entscheide zum
Zuschnitt (Reichweite: gemeinsame Extraktion; Aktivierung: standardmäßig an)
vom 2026-09-27; slice-232 <!-- d-check:status-provenance -->.

**Schärft:**
[`DC-FA-LINK-001`](../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
(Erweiterung — kein neues Kürzel, dieselbe Prüfung, verengte
Extraktions-Grenze).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`parseLinkAt`/`forEachLink` (`internal/hexagon/core/rules/markdown.go`)
suchen die schließende Klammer eines Linkziels innerhalb **einer** Zeile
(`Line.Text` bzw. ein einzelner `string`). Ein gültiger Markdown-Link, dessen
Zieladresse hinter `](` einen Zeilenumbruch trägt, wird deshalb nicht
gefunden — dokumentierte Grenze in
[`DC-FA-LINK-001.a`](../../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
Schritt 3 („zeilenbasiert … normative Grenze für alle Module"). Ein
Konsument bittet um Erkennung dieser Form; der Auftraggeber hat vorab
entschieden, dass die Erkennung **gemeinsam** (nicht nur im Modul `links`)
und **standardmäßig aktiv** sein soll.

Die Extraktion hat zwei Signaturen mit unterschiedlichem Adressraum:

- `ExtractLinks(lines []Line) []LinkRef` — hat Zugriff auf **alle** Zeilen
  der Datei; sechs Konsumenten: `links`, `links.resolve-from` (dieselbe
  Extraktion, zusätzliche Ortsauflösung), `anchors`, `matrix`, `external`,
  `tracked`.
- `ExtractLinkSpans(text string) []LinkSpan` — bekommt **eine** Zeile als
  reinen `string`, strukturell ohne Kenntnis der Nachbarzeilen; vier
  Konsumenten: `ids` (Linkpflicht-Ausnahme im Linktext), `pins`
  (`dpin`-Marker-Bindung „derselbe Zeile"), `--repair`
  (`internal/hexagon/core/app/repair.go`, schreibt Fix-Kandidaten auf
  Byte-Spannen einer Zeile), `planning` (Zitat-Erkennung der
  Beobachtungs-Registrierung).

Nur die erste Signatur kann eine mehrzeilige Form überhaupt sehen, ohne ihre
Schnittstelle zu ändern.

## Entscheidung

### (a) Absatzweise Erkennung, aber nur für `ExtractLinks`

`ExtractLinks` gruppiert die übergebenen `Line`s zu Absätzen (Leerzeile bzw.
Fenced-Block-Lücke trennt — dieselbe Grenzziehung, die `proseParagraphs`
bereits für die absatzweise Inline-Code-Erkennung nutzt,
[`DC-FA-LINK-001.a`](../../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
Schritt 2), fügt die Zeilen eines Absatzes mit `\n` zusammen und ruft
`forEachLink` **einmal auf dem zusammengefügten Text** auf, statt einmal je
Zeile. Ein gefundener Link/ein Bild wird der Zeile zugeordnet, auf der seine
öffnende Klammer (`[`/`![`) liegt — für einen einzeiligen Treffer ist das
exakt die bisherige Zeile (verhaltenserhaltend), für den neuen Fall die
`](`-Zeile.

`ExtractLinkSpans` bleibt **unverändert** auf ihrer String-Signatur. Ihre
vier Konsumenten sehen die neue Form strukturell nicht — das ist keine
Lücke, sondern die Konsequenz einer Schnittstelle, die nie mehr als eine
Zeile bekommt. Im Einzelnen:

- **`ids`:** Eine Kennung im Linktext eines mehrzeiligen Links bleibt von
  der Ausnahme unerfasst — der seltene Fall einer Kennung *im Linktext*
  eines Links, dessen *Adresse* umbricht, bleibt wie bisher zeilenbasiert
  behandelt. Kein Vertragsbruch:
  [`DC-FA-ID-001.a`](../../../spec/spezifikation.md#dc-fa-id-001a--kennungs-prüfung)
  verspricht die Ausnahme
  nie über den bereinigten Text einer Zeile hinaus.
- **`pins`:** Der `dpin`-Marker bindet weiterhin an einen Link **derselben
  Zeile** — ein Link, dessen Adresse jetzt über zwei Zeilen läuft, hat keine
  Zeile mehr, auf der `forEachLink` (zeilenweise, unverändert) ihn fände;
  die Bindung schlägt fehl wie bei jedem heute schon nicht erkannten Link.
  Unverändertes Verhalten, keine neue Lücke.
- **`--repair`:** Erzeugt keinen Fix-Kandidaten für die neue Form — sie
  entsteht als `target-missing`-Befund über `ExtractLinks`, aber
  `repairLine` sucht ihre Reparatur-Spanne über `ExtractLinkSpans` auf
  **einer** Zeile und findet dort keinen passenden Span. Automatische
  Konsequenz der getrennten Signaturen, keine gesonderte Ausschluss-Logik
  nötig.
- **`planning`:** Zitat-Erkennung der Beobachtungs-Registrierung bleibt
  zeilenbasiert; dieselbe Begründung wie bei `ids`.

### (b) Fund-Zeile: die öffnende Zeile

Ein Befund (z. B. `target-missing`) nennt die Zeile, auf der der Link
**öffnet** (`[text](` bzw. `![alt](`) — nicht die Zeile der Adresse. Das ist
sowohl für einzeilige Treffer **unverändertes** Verhalten (die einzige Zeile
ist ohnehin die öffnende) als auch die intuitivere Reparaturstelle: Wer den
Link reparieren will, sucht ihn dort, wo `[…](` beginnt.

### Bilder ziehen mit, unpromised

`parseLinkAt` behandelt Links und Bilder über denselben Codepfad
(`isImage`-Flag). Ein Bild, dessen Zieladresse ebenfalls über einen
Zeilenumbruch geht, wird deshalb automatisch mit erkannt — **keine eigene
Zusage**, nur eine Konsequenz des geteilten Parsers. Ebenso ungemessen und
unpromised bleibt ein Zeilenumbruch **im Linktext** oder zwischen Adresse
und Titel (`](ziel` ⏎ `"titel")`) — beides kann durch die
Absatz-Zusammenfügung technisch mit erfasst werden, ist aber vom
Change Request nicht verlangt und nicht getestet.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Nichts tun** (Bestand) | keine Änderung | löst den Change Request nicht; die Grenze ist real beobachtetes Konsumenten-Verhalten |
| **Separater, isolierter Pfad nur im Modul `links`** (kein Fold in `ExtractLinks`, ein eigener regex-/Sonderfall-Codepfad, der nur `CheckLinks` zusätzlich speist) | kleinster Blast-Radius: null Konsequenzen für `anchors`/`matrix`/`external`/`tracked` | widerspricht der Auftraggeber-Entscheidung „gemeinsame Extraktion"; erzeugt eine ZWEITE Link-Erkennungs-Logik neben `parseLinkAt` — genau die Art von Doppel-Definition, die `LinkSuffixEnd`s Dokumentationskommentar bereits als Fehlerquelle nennt (Klammern im Ziel rissen bei einem regex-Sonderpfad in slice-073 <!-- d-check:status-provenance --> echte Waisen unsichtbar) |
| **Absatzweise Faltung in `ExtractLinks`, `ExtractLinkSpans` unverändert** (gewählt) | folgt der etablierten, bereits geprüften Absatz-Technik (`stripInlineCodeByLine`/`inlineSpansByLine`); erfüllt „gemeinsame Extraktion" für die sechs `[]Line`-Konsumenten; keine zweite Link-Grammatik; verhaltenserhaltend für jeden einzeiligen Bestandslink (Faltung ändert an der Fundzeile nichts, wenn nichts zu falten ist) | vier String-Konsumenten bleiben bewusst asymmetrisch zurück (benannt, nicht verschwiegen); Bilder und Linktext-Umbruch „ziehen mit", ohne eigens getestet zu sein |

## Konsequenzen

- `ExtractLinks` gruppiert `[]Line` jetzt zu Absätzen (Leerzeile/Fenced-Lücke
  trennt) statt Zeile für Zeile zu iterieren; `ExtractLinkSpans` und
  `forEachLink`s Signatur bleiben unverändert.
- **Bestandsmessung Pflicht vor Closure:** `make doc-check` auf diesem Repo
  vor/nach der Änderung vergleichen — die Absatz-Faltung könnte theoretisch
  einen Bestandslink anders zuordnen, wenn ein Absatz eine unbalancierte
  `[`/`]`-Folge über mehrere Zeilen trägt, die vorher pro Zeile isoliert
  betrachtet wurde. Jeder neue Befund wird behoben oder mit Grund benannt.
- Keine Änderung an `--repair`, `ids`, `pins`, `planning` — ihr Verhalten ist
  oben einzeln benannt, nicht implementiert.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `go test` | `TestExtractLinks_*` (Zeilenumbruch-Form: Happy/Negative/Rot-Beleg gegen den alten Stand, unveränderte Inline-/Titel-/Spitzklammer-/Klammer-Kontrollen, Bild-Form gemessen) | `make test` |
| `make doc-check` | Bestandsmessung: Befundzahl dieses Repos vor/nach der Änderung | Closure-Notiz |

## Re-Evaluierungs-Trigger

Zwei Bedingungen, jede für sich hinreichend:

1. **Ein Konsument bittet um dieselbe Erkennung für einen der vier
   String-Konsumenten** (`ids`, `pins`, `--repair`, `planning`). Dann ist zu
   prüfen, ob ihre Signatur auf absatzweise Eingabe umgestellt wird, oder ob
   die Asymmetrie bestehen bleibt.
2. **Die Bestandsmessung zeigt einen falschen Bestandslink**, den die
   Absatz-Faltung neu (fehl-)interpretiert. Dann ist die Grenzziehung
   (Leerzeile/Fenced-Lücke) neu zu bewerten.

Ohne eines von beiden: permanent.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-09-27 | Proposed → Accepted (`slice-232`) |
| 2026-09-27 | Accepted → Superseded by ADR-0092: unabhängiger Review (R1-H1, HIGH) hat die absatzweise Faltung als fehlerhaft nachgewiesen — ein unbalanciertes `[` in gewöhnlicher Prosa konnte mit einer späteren, unabhängigen `](…)`-Sequenz im selben Absatz zu einem erfundenen Link verschmelzen. ADR-0092 begrenzt den Lookahead auf eine Zeile und die Adress-Klammer |
