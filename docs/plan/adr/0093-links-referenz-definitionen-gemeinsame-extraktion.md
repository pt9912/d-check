# ADR-0093: Link-Referenz-Definitionen werden über die gemeinsame Extraktion geprüft, unabhängig von ihrer Verwendung — mit einer Ausnahme (`anchors`)

**Status:** Superseded by ADR-0094

**Datum:** 2026-09-27

**Autor:** claude-sonnet-5

**Bezug:** Change Request eines Konsumenten; Auftraggeber-Entscheid zum
Zuschnitt (jede Definition mit Datei-Ziel, standardmäßig an) vom
2026-09-27; slice-233 <!-- d-check:status-provenance -->.

**Schärft:**
[`DC-FA-LINK-001`](../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
(Erweiterung — kein neues Kürzel, verengter Out-of-Scope-Satz),
[`DC-FA-ANCH-001`](../../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors)
(neuer Out-of-Scope-Satz).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Eine Link-Referenz-Definition (`[label]: ziel "titel"`) bindet ein Label an
eine Zieladresse; an anderer Stelle im Dokument kann sie über
`[text][label]`, `[label][]` oder `[label]` **verwendet** werden. Das
Lastenheft schloss Reference-Style-Links bisher vollständig aus
([`DC-FA-LINK-001`](../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
Out-of-Scope). Ein Konsument bittet darum, wenigstens die Definition selbst
zu prüfen — unabhängig davon, ob sie benutzt wird: eine unbenutzte
Definition mit totem Ziel ist ebenso ein Bestandsfehler wie eine benutzte.

Die gemeinsame Extraktion `ExtractLinks([]Line) []LinkRef`
([ADR-0092](0092-links-zeilenumbruch-begrenzter-lookahead.md)) speist sechs
Konsumenten: `links`, `links.resolve-from`, `anchors`, `matrix`, `external`,
`tracked`. Eine Definition ist syntaktisch klar von einem Inline-Link
unterscheidbar (`:` statt `(` hinter der schließenden Klammer) und lässt
sich ohne zweite Link-Grammatik in dieselbe Liste einspeisen — dieselbe aus
ADR-0091 übernommene Präferenz „keine zweite Definitions-Sprache neben
`parseLinkAt`".

## Entscheidung

### Definitionen speisen dieselbe `ExtractLinks`-Liste — mit einem Discriminator

`LinkRef` bekommt ein Feld `IsDefinition bool`. `ExtractLinks` erkennt pro
Zeile zusätzlich zu Inline-Links eine Definitions-Zeile (Regex: bis zu drei
führende Leerzeichen, `[label]:`, Whitespace, Ziel + optionaler Titel auf
derselben Zeile) und liefert dafür einen `LinkRef` mit `IsDefinition: true`,
`Text` = Label, `Target` = normalisierte Zieladresse (`NormalizeTarget`,
identisch zur Titel-Abtrennung eines Inline-Links).

**Fünf der sechs Konsumenten behandeln eine Definition wie jeden anderen
`LinkRef`** — `links` prüft ihr Ziel (Existenz, Escape, Symlink,
`ignore-refs`), `links.resolve-from` ihre Ortsfestigkeit, `matrix` ihre
Klassen-/Status-Regeln, `external`/`tracked` ihre externe Erreichbarkeit
bzw. ihren Git-Tracked-Status — konsistent mit einem Inline-Link auf
dasselbe Ziel: Ein Inline-Link auf eine `Superseded`-ADR meldet
`matrix-inactive`, eine Definition mit demselben Ziel soll das ebenfalls
tun, sonst entstünde eine stille Lücke in genau der Prüfung, die der Slice
eigentlich schärft.

**`anchors` ist die eine Ausnahme.** Es prüft `if ref.IsDefinition {
continue }` vor jeder anderen Prüfung — eine Definition mit
Fragment-Ziel (`[label]: datei.md#anker`) bleibt unbeanstandet. Begründung:
Der Change Request maß nur die Datei-Ziel-Prüfung; eine falsch beanstandete
Definition wäre ein neuer Befund ohne zugehörige Anforderung.

### Drei benannte Grenzen der erkannten Form

1. **Kein Blockquote-/Listen-Präfix.** Die Definition muss (nach bis zu drei
   führenden Leerzeichen) mit `[` beginnen — eine Zeile mit `>` oder einem
   Listen-Marker davor wird nicht erkannt, auch wenn CommonMark eine
   eingebettete Definition dort zuließe.
2. **Keine Backslash-Escapes im Label.** `[foo\]bar]: ziel` wird beim ersten
   `]` beendet — das ist eine falsche Label-Grenze, aber die Zieladresse
   bleibt in beiden Lesarten dieselbe Zeichenfolge dahinter, sodass die
   Ziel-Prüfung selbst nicht falsch wird.
3. **Keine mehrzeilige Form.** Ziel oder Titel hinter einem Zeilenumbruch
   werden nicht erkannt — dieselbe zeilenlokale Grenze wie die Linktext-
   Klammer eines Inline-Links (ADR-0092): die Definition ist strukturell
   einzeilig.

Alle drei sind eine bewusst **enger** gefasste Form als das volle
CommonMark-Definitions-Grammatik, ausgeschrieben **vor** der Bestandsmessung
(`AGENTS.md` §5).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Separater, isolierter Pfad nur im Modul `links`** (kein Einspeisen in `ExtractLinks`, eigener Erkennungs-Codepfad nur für `CheckLinks`) | kleinster Blast-Radius: null Konsequenz für `matrix`/`external`/`tracked`/`anchors` | zweite Link-Grammatik neben `parseLinkAt` — dieselbe Fehlerklasse, die ADR-0091 bereits als Grund gegen diese Option nannte (Klammern im Ziel rissen in slice-073 <!-- d-check:status-provenance --> echte Waisen unsichtbar, weil ein Sonder-Regex die Haupt-Extraktion dupliziert hatte) |
| **Gemeinsame Extraktion, alle sechs Konsumenten gleich** (auch `anchors` prüft Definitionen) | einfachste Regel: keine Ausnahme | `anchors` erzeugt Befunde für eine Prüfung, die der Change Request nicht verlangt hat — ungemessene Ausweitung des Vertrags |
| **Gemeinsame Extraktion mit `IsDefinition`-Discriminator, `anchors` ausgenommen** (gewählt) | eine Grammatik, konsistente Durchsetzung über `links`/`matrix`/`external`/`tracked`/`resolve-from`; die eine Ausnahme ist explizit und lokalisiert (ein `continue` in `anchors.go`) | ein `LinkRef`-Feld mehr; jeder künftige `ExtractLinks`-Konsument muss bewusst entscheiden, ob er Definitionen mitsehen will |

## Konsequenzen

- `LinkRef.IsDefinition bool` (neues Feld); `ExtractLinks` erkennt
  Definitions-Zeilen zusätzlich zu Inline-Links, unabhängig vom
  Zeilenumbruch-Lookahead aus ADR-0092 (eine Definition ist immer einzeilig,
  der Lookahead-Pfad greift nicht).
- `CheckAnchors` (`anchors.go`) überspringt `IsDefinition`-Refs als erste
  Bedingung der Schleife.
- Kein Konsument außer `links`/`anchors` ändert Code — `matrix`, `external`,
  `tracked`, `links.resolve-from` iterieren bereits generisch über jeden
  `LinkRef`.
- Bestandsmessung gegen alle sechs Konsumenten (wie in ADR-0092
  nachgezogen).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `go test` | Definitions-Erkennung: totes Ziel → `target-missing` auf der Definitions-Zeile; lebendes Ziel → kein Befund; Definition im Fenced-Block → kein Befund; Blockquote-/Listen-Präfix, Backslash-Escape im Label, Zeilenumbruch vor Ziel → nicht erkannt (Rot-Beleg: dieselbe Zeile ohne Präfix/Escape/Umbruch erzeugt den Fund) | `make test` |
| `go test` | `anchors` überspringt eine Definition mit totem Fragment-Ziel — kein Befund, obwohl `links` bei fehlender Datei melden würde | `make test` |
| `make doc-check` | Bestandsmessung über alle sechs Konsumenten (Closure-Notiz) | Closure-Notiz |

## Re-Evaluierungs-Trigger

Drei Bedingungen, jede für sich hinreichend:

1. Ein Konsument bittet um Fragment-/Anker-Prüfung für Definitionen
   (`anchors` müsste die Ausnahme aufheben).
2. Ein Konsument bittet um Auflösung der **Verwendung**
   (`[text][label]`/`[label][]`/`[label]`) auf ihre Definition.
3. Ein Konsument bittet um eine der drei benannten Grenzen (Blockquote/
   Liste, Backslash-Escape, Zeilenumbruch).

Ohne eines von diesen: permanent.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-09-27 | Proposed → Accepted (`slice-233`) |
| 2026-09-27 | Accepted → Superseded by ADR-0094: eigene Verifikation der Regex vor dem ersten Test widerlegte die Rationale zu Grenze 2 — ein `\]` im Label führt nicht zu einer falschen, aber unschädlichen Label-Grenze, sondern lässt die **ganze** Zeile unerkannt (der ankernde regulär Ausdruck scheitert vollständig). ADR-0094 korrigiert nur diesen Satz |
