# ADR-0095: Der Rest der Zeile hinter dem Ziel muss ein korrekt abgegrenzter Titel sein, sonst bleibt die ganze Definition unerkannt (supersedes ADR-0094)

**Status:** Accepted

**Supersedes:** ADR-0094

**Datum:** 2026-09-27

**Autor:** claude-sonnet-5

**Bezug:** Unabhängiger Review-Report
`docs/reviews/2026-09-27-slice-233-links-referenz-definitionen-review-r1.md` <!-- d-check:status-provenance -->,
Befund R1-H1; slice-233 <!-- d-check:status-provenance -->.

**Schärft:**
[`DC-FA-LINK-001`](../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
(präzisiert die Erweiterung aus ADR-0093, keine neue Anforderung).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

ADR-0093 behauptete: „`Target` = normalisierte Zieladresse (`NormalizeTarget`,
**identisch zur Titel-Abtrennung eines Inline-Links**)." Die ursprüngliche
Regex (`^ {0,3}\[([^\]\n]+)\]:[ \t]+(\S.*)$`) erfasste als zweite Gruppe
jedoch den **gesamten Rest der Zeile** und überließ `NormalizeTarget` die
Abtrennung — jene Funktion schneidet aber nur naiv am **ersten** Whitespace
ab, ohne zu prüfen, ob danach überhaupt ein gültig delimitierter Titel
(`"…"`, `'…'` oder `(…)`) folgt.

Der unabhängige Review (R1-H1, HIGH) konstruierte den Gegenbeweis: eine
gewöhnliche Glossar-Prosazeile `[TERM]: First In, First Out` — nach
CommonMark **keine** gültige Link-Referenz-Definition, weil der Rest hinter
dem ersten Token weder leer noch ein delimitierter Titel ist — wurde
trotzdem als Definition erkannt, mit dem **erfundenen** Ziel „First". Das
Modul `links` meldete dafür `target-missing`, ein Fehlbefund auf einer
Zeile, die nach dem Vertrag gar keine Definition ist.

## Entscheidung

### Der Rest der Zeile wird strukturell validiert, nicht mehr naiv abgeschnitten

`definitionRe` bekommt eine dritte, **optionale** Gruppe für den Titel und
verankert das Zeilenende (`$`) direkt dahinter:

```
^ {0,3}\[([^\]\n]+)\]:[ \t]+(<[^<>\n]*>|\S+)(?:[ \t]+(?:"[^"\n]*"|'[^'\n]*'|\([^()\n]*\)))?[ \t]*$
```

- Die **Ziel-Gruppe** ist entweder `<…>`-umschlossen oder ein Token **ohne
  eingebetteten Whitespace** (`\S+`) — dieselbe Grundform wie bei einem
  Inline-Link.
- Danach darf **nichts** mehr folgen außer optionalem Whitespace oder
  **genau einem** korrekt abgegrenzten Titel (Anführungszeichen `"…"`/`'…'`
  oder Klammern `(…)`, CommonMark-Titel-Delimiter für Definitionen).
- Passt der Rest der Zeile in **keine** dieser Formen — wie bei
  `[TERM]: First In, First Out` —, scheitert der gesamte, verankerte Ausdruck:
  **keine Definition entsteht**, dieselbe „ganze Zeile bleibt unerkannt"-
  Disziplin wie bei den Grenzen aus ADR-0093/ADR-0094 (Blockquote-Präfix,
  Backslash-Escape).

`NormalizeTarget` bleibt auf der jetzt bereits sauber abgegrenzten
Ziel-Gruppe im Einsatz (nur noch für die `<…>`-Entquotung relevant — ihr
Whitespace-Split-Zweig greift für ein `\S+`-Match nie, weil dort per
Konstruktion kein Whitespace mehr enthalten ist).

**Korrigierte Aussage:** Die Titel-**Erkennung** einer Definition ist
**strenger** als die eines Inline-Links (echte Delimiter-Prüfung statt
naivem Whitespace-Schnitt) — nicht „identisch", wie ADR-0093 behauptete.
Die **Ziel-Normalisierung** der bereits abgegrenzten Ziel-Gruppe
(`NormalizeTarget`) bleibt dieselbe Funktion.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Naiver Whitespace-Schnitt beibehalten** (ADR-0093, verworfen) | einfachste Umsetzung, bereits implementiert | nachweislich fehlerhaft (R1-H1): erfundene Ziele aus gewöhnlicher Prosa |
| **Ganze Zeile erfassen, Titel-Validität nachträglich prüfen und bei Fehlschlag verwerfen** | bliebe näher an der bisherigen Struktur | Reparatur eines fehlerhaften Fundes ist fehleranfälliger als ihn gar nicht erst zuzulassen — dieselbe Erwägung wie in ADR-0092 gegen einen nachträglichen Filter |
| **Titel-Delimiter-Prüfung direkt in der Regex, Zeilenende dahinter verankert** (gewählt) | eliminiert die Fehlerklasse durch Konstruktion; keine zweite Prüf-Stelle, die aus dem Takt geraten könnte | eine Definition mit einem **nicht** delimitierten, aber vom Autor gemeinten Titel (z. B. `[l]: ziel.md Titel ohne Anfuehrungszeichen`) bleibt unerkannt — nach CommonMark ohnehin ungültig, keine neue Grenze |

## Konsequenzen

- `definitionRe` (`internal/hexagon/core/rules/markdown.go`) um die
  Titel-Delimiter-Gruppe erweitert; `parseDefinitionLine` unverändert
  (arbeitet weiterhin auf der Ziel-Gruppe des Regex-Treffers).
- Neuer Test: `[TERM]: First In, First Out` erzeugt **keinen** `LinkRef`
  (Rot-Beleg: schlug an der ADR-0093/0094-Fassung fehl, ist an dieser
  Fassung grün).
- Kein Konsument außer `links`/`ExtractLinks` selbst ändert Code.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `go test` | `[TERM]: First In, First Out` → kein `LinkRef` (Rot-Beleg gegen die Vorfassung) | `make test` |
| `go test` | Bestehende Definitions-Tests (Happy, Titel, Blockquote, Backslash) bleiben grün | `make test` |

## Re-Evaluierungs-Trigger

Ein Konsument bittet um eine Definition mit einem Titel, der nicht durch
`"…"`/`'…'`/`(…)` abgegrenzt ist (CommonMark kennt dafür ohnehin keine
gültige Form). Ohne das: permanent.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-09-27 | Proposed → Accepted (`slice-233`, Nachzug nach unabhängigem Review R1-H1) |
