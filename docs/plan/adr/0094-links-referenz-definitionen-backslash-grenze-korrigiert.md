# ADR-0094: Ein Backslash im Definitions-Label lässt die ganze Zeile unerkannt, nicht nur die Label-Grenze (supersedes ADR-0093)

**Status:** Accepted

**Supersedes:** ADR-0093

**Datum:** 2026-09-27

**Autor:** claude-sonnet-5

**Bezug:** eigene Verifikation der Erkennungs-Regex vor dem ersten Test;
slice-233 <!-- d-check:status-provenance -->.

**Schärft:**
[`DC-FA-LINK-001`](../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
(präzisiert die Erweiterung aus ADR-0093, keine neue Anforderung).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

ADR-0093 benannte drei Grenzen der Definitions-Erkennung, darunter: „ein
`[foo\]bar]: ziel` wird beim ersten `]` beendet — das ist eine falsche
Label-Grenze, aber die Zieladresse bleibt in beiden Lesarten dieselbe
Zeichenfolge dahinter, sodass die Ziel-Prüfung selbst nicht falsch wird."

Vor dem ersten Test wurde die Erkennungs-Regex
(`^ {0,3}\[([^\]\n]+)\]:[ \t]+(\S.*)$`) gegen genau diesen Fall verifiziert:
`[foo\]bar]: ziel.md` liefert **keinen** Treffer. Die Zeichenklasse
`[^\]\n]` kann kein `]` konsumieren; die Erfassungsgruppe endet zwingend vor
dem ersten literalen `]` (nach „foo\"). Der Regex verlangt **unmittelbar**
danach `]:` — tatsächlich folgt aber `bar]:`. Der verankerte
(`^`/`$`) Ausdruck kann diesen Fehlschlag nicht durch Rückwärtssuche an
anderer Stelle beheben: **die ganze Zeile bleibt unerkannt**, keine
Definition entsteht, keine Ziel-Prüfung läuft — der Gegenteil der
ursprünglichen Behauptung.

## Entscheidung

**Grenze 2 wird korrigiert, nicht die Erkennungs-Regex.** Die Regex bleibt
unverändert (`internal/hexagon/core/rules/markdown.go`,
`definitionRe`/`parseDefinitionLine`) — das Verhalten „ganze Zeile
unerkannt bei einem `\]` vor dem echten Label-Ende" ist akzeptabel,
dieselbe Klasse Grenze wie „kein Blockquote-/Listen-Präfix": die Definition
wird schlicht nicht erkannt, still, ohne Fehlbefund. Korrigiert wird nur
die **Beschreibung** dieser Grenze in ADR und Spec-Straten: statt „falsche
Label-Grenze, Ziel bleibt korrekt" heißt es „die ganze Zeile bleibt
unerkannt".

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Regex so ändern, dass sie einen escapten `]` überspringt** | würde die ursprünglich behauptete (aber falsche) Semantik nachträglich herstellen | mehr Komplexität für eine Form, die der Change Request nicht verlangt hat (§1 Abgrenzung nennt Backslash-Escapes ausdrücklich als offenen, zurückgestellten Punkt) |
| **Nur die Rationale-Prosa korrigieren, Regex unverändert** (gewählt) | kleinstmögliche Änderung; das tatsächliche, jetzt korrekt beschriebene Verhalten ist bereits sicher (still unerkannt, kein Fehlbefund) | ein Label mit `\]` bleibt für die Definitions-Prüfung vollständig unsichtbar — benannt, nicht behoben |
| **Nichts korrigieren, falsche Rationale stehen lassen** | kein Aufwand | eine falsche Beschreibung im Vertrag ist eine Harness-Lüge, sobald sie in einem Fall zutrifft, den ein Konsument tatsächlich hat |

## Konsequenzen

- ADR-0093 Grenze 2 gilt inhaltlich als durch diese ADR ersetzt; der
  Discriminator-Mechanismus (`LinkRef.IsDefinition`), die
  Konsumenten-Behandlung und die anderen beiden Grenzen aus ADR-0093
  bleiben unverändert gültig.
- `spec/lastenheft.md` und `spec/spezifikation.md` übernehmen die
  korrigierte Formulierung.
- Kein Code ändert sich — nur die Beschreibung eines bereits vorhandenen,
  sicheren Verhaltens.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `go test` | `[foo\]bar]: ziel.md` erzeugt **keinen** `LinkRef` (weder mit korrektem noch mit falschem Label) | `make test` |

## Re-Evaluierungs-Trigger

Ein Konsument bittet um Backslash-Escape-Unterstützung im Label (dieselbe
Bedingung wie ADR-0093 Re-Evaluierungs-Trigger 3). Ohne das: permanent.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-09-27 | Proposed → Accepted (`slice-233`, Nachzug vor dem ersten Test) |
