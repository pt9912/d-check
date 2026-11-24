# ADR-0116: `linelength.exempt-tables` nimmt auch eine `|`-Zeile mit bis zu drei Leerzeichen davor aus

**Status:** Accepted

**Datum:** 2026-11-24

**Autor:** pt9912

**Bezug:** [Issue #20](https://github.com/pt9912/d-check/issues/20) (Tabellen in Listenpunkten),
[`DC-FA-LLEN-001`](../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in),
[Lastenheft §5](../../../spec/lastenheft.md#5-globale-out-of-scope-punkte) (die Ausnahme für die Zeilenlänge
„mit derselben Markdown-Lexik“),
[ADR-0114](0114-linelength-modul.md) (das Modul; Entscheidung 5 (b) wird hier teilweise abgelöst),
[ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) (Entscheidung 1, 2 und 4: andere Antwort oder
andere Frage, je Konsument eine Assertion),
slice-293 <!-- d-check:status-provenance -->.

**Supersedes:** [ADR-0114](0114-linelength-modul.md) Entscheidung 5 (b) — **nur** der Satz „eingerückte
Zeilen sind keine“ und die Konsequenz „Eine eingerückte Tabellenzeile … wird gemessen“; Entscheidung 5 (b)
im Übrigen (Zeilen außerhalb von Fenced-Code, die mit `|` beginnen; keine Trennzeile verlangt) und alle
übrigen Entscheidungen von ADR-0114 bleiben gültig.

**Schärft:**
[`DC-FA-LLEN-001.a`](../../../spec/spezifikation.md#dc-fa-llen-001a--zeilenlänge-in-markdown-dateien-linelength)
Schritt 4, „Nicht Gegenstand“ und „Grenzen, am Gegenstand gelesen“, dazu die Zeile `linelength.exempt-tables`
in [`SPEC-005`](../../../spec/spezifikation.md#spec-005--d-checkyml).

**Regeln:** Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v6.18.1/regelwerk/modul-04-adrs.md).

---

## Kontext

Ein Repo führt Tabellen in nummerierten Listenpunkten; unter `1. ` stehen sie mit drei Leerzeichen davor.
`linelength` misst diese Zeilen, weil seine Ausnahme (b) die Antwort `tableRowLine` der geteilten Lexik
übernimmt: Tabellenzeile ist eine Zeile außerhalb von Fenced-Code, deren **rohe** Zeile mit `|` beginnt
(ADR-0114 Entscheidung 5 (b)). Eine Tabellenzeile lässt sich nicht umbrechen; der Melder hat 41 von 85
Befunden in dieser Form (Angabe des Issues, von hier nicht gemessen).

`tableRowLine` hat außer `linelength` drei Konsumenten. `targets` liest daraus Doku-Targets, die
Spezifikation sagt ihm „erstes Zeichen ein Pipe, Spalte 0“ zu
([`DC-FA-TGT-001.a`](../../../spec/spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets)
und die Schema-Zeile `targets.doc-tables`). Weitere Konsumenten sind `planning.waves` (Registerzeilen) und
`structure` (Zell- und Reihenfolge-Bedingungen). `tableHeaderOrSeparator` baut auf `tableRowLine` auf;
`tableSeparatorRow` prüft das führende `|` selbst an der rohen Zeile. Die Kopplung
`TestTabellenzeilenFrageHatEineAntwort` hält heute alle vier Konsumenten auf **einer** Antwort, auch für eine
Zeile mit zwei Leerzeichen davor (Antwort: keine Tabellenzeile).

Für den Wunsch lassen sich die übrigen Module also nicht mitbewegen, ohne dass sich ihre Zusagen ändern. Und
ADR-0054 Entscheidung 1 verbietet, dass ein Konsument eine Lexik-Frage selbst beantwortet. Entscheidung 2 lässt
eine **andere Frage** zu, wenn eine ADR sie scopt. `linelength` fragt nicht „deklariert diese Zeile etwas in
einer Spalte-0-Tabelle“, sondern „ist diese Zeile Teil einer Tabelle und deshalb nicht umbrechbar“.

Die Entscheide 1 bis 3 hat der Auftraggeber am 2026-11-24 bestätigt (Fragen 1 bis 3 der Welle, Frage 3 als
„höchstens drei Leerzeichen, kein Listenkontext“). Entscheid 4 ist die Antwort des Architect auf Frage 4 der
Welle (Verdikt unter `docs/reviews/`).

## Entscheidung

1. **Kein neuer Schlüssel.** `linelength.exempt-tables` selbst wird erweitert; im Default an, mit `false` wird
   die eingerückte Zeile gemessen wie jede andere.
2. **`> |` bleibt gemessen.** Eine Zeile hinter einem Zitat-Zeichen ist auch künftig keine Tabellenzeile.
3. **Höchstens drei Leerzeichen, kein Listenkontext.** Ausgenommen wird eine Zeile außerhalb von Fenced-Code,
   vor deren erstem `|` null bis drei Leerzeichen (U+0020) stehen und sonst nichts. Ein Tab, ein anderes
   Leerraum-Zeichen (etwa U+00A0) und vier oder mehr Leerzeichen davor lassen die Zeile gemessen. Der Schnitt
   bei vier ist derselbe, den die Anforderung für den eingerückten Code-Block ohne Zaun schon zieht.
4. **Die Antwort wohnt in der geteilten Lexik, als eigene, benannte Frage.** Ein neuer Helfer in der
   geteilten Lexik neben `tableRowLine` beantwortet die Frage aus Entscheid 3. Er **schließt `tableRowLine`
   ein**: Er ruft sie, statt die Spalte-0-Antwort ein zweites Mal zu schreiben, und erweitert sie um die
   Einrückung. `linelength` ruft für (b) nur diesen Helfer. `tableRowLine` bleibt zeichengleich, und `targets`,
   `planning.waves` und `structure` lesen eine eingerückte Zeile weiter nicht als Tabellenzeile. Nach
   ADR-0054 Entscheidung 2 ist das eine andere Frage, gescopt durch diese ADR, keine zweite Antwort. Die
   Kopplung hält beide Fragen (ADR-0054 Entscheidung 4):
   - die drei Konsumenten von `tableRowLine` auf einer Antwort, die eingerückte Zeile eingeschlossen;
   - `linelength` auf der Antwort des Helfers, in beiden Richtungen: bis drei Leerzeichen ausgenommen; vier
     Leerzeichen, Tab und `>` gemessen; im Zaun und außerhalb wie `tableRowLine`.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; der Melder löst die Tabellen aus den Listenpunkten | kein Code | eine Tabellenzeile lässt sich nicht umbrechen, die Abhilfe ist eine andere Gliederung des Dokuments; der Melder ersetzt sein Gate nicht |
| B — `tableRowLine` selbst weiten | eine Antwort für alle vier Konsumenten | verschiebt `targets` gegen seine Spalte-0-Zusage in der Spezifikation, dazu `planning.waves` und `structure`; jedes Modul bräuchte eigene Kriterien. `tableSeparatorRow` prüft das führende `\|` selbst, also wäre die Weitung nur halb, wenn sie nicht auch dort ansetzt |
| C — eigener Schalter (etwa `exempt-indented-tables`) | wer `exempt-tables` fährt, behält das alte Verhalten | eine Schema-Zeile, ein Fall am Konfigurations-Rand und Kriterien je Wert, ohne dass der Anlass sie braucht |
| D — zwei bis fünf Leerzeichen, nur im Listenkontext (Wortlaut des Issues) | deckt `10. ` und verschachtelte Listen | braucht eine Listen-Erkennung in der geteilten Lexik, die es nicht gibt; die Listen-Muster im Code sind private Einzelmuster mit eigenem Zweck |
| E — jede Einrückung (die Zeile vor `\|` links trimmen) | einfachste Regel | nähme einen eingerückten Code-Block aus, der mit `\|` beginnt, gegen die bestehende Grenze „kein Zaun, kein Fenced-Code“ |
| F — eine private Erkennung in `linelength` | berührt die geteilte Lexik nicht | ein Konsument beantwortet eine Lexik-Frage selbst (ADR-0054 Entscheidung 1). Lastenheft §5 sagt dem Modul zu, mit derselben Markdown-Lexik zu laufen wie die übrigen; eine private Antwort hielte das nicht |
| G — ein Zeilen-Marker (ADR-0114, Alternative I) | trifft jede Zeile | ADR-0114 Entscheidung 9 schließt ihn aus; der Melder müsste 41 Zeilen markieren |
| **H — gewählt: `exempt-tables` erweitert, bis drei Leerzeichen, ein benannter Helfer in der Lexik, der `tableRowLine` einschließt** | kleinste Antwort, die den Anlass trägt; kein anderes Modul bewegt sich; die Spalte-0-Antwort bleibt einmal geschrieben | Tabellen mit vier und mehr Leerzeichen bleiben gemessen (Konsequenzen) |

## Konsequenzen

- Positiv: eine Tabelle unter `1. ` oder unter `- ` (zwei Leerzeichen) erzeugt im Default keinen Befund mehr.
  Kein anderes Modul ändert sein Ergebnis.
- Negativ — die Bezugsmenge nennt bei solchen Dateien weniger geprüfte Zeilen. Gewollt: geprüft ist, was nach
  (a) und (b) bleibt. Eine eingerückte Tabellenzeile mit einem Wort über der Grenze nimmt (b) aus, bevor (c)
  sie sieht; sie zählt deshalb auch nicht in `U`, wie die Spalte-0-Zeile heute.
- Negativ — Grenzen der Regel, je mit ihrer Bedingung:
  - Eine Tabelle unter einer Marke von vier Spalten (`10. `) oder in einem verschachtelten Listenpunkt steht mit
    vier und mehr Leerzeichen und bleibt gemessen.
  - Eine Tabellenzeile auf der Zeile der Listenmarke (`- | a | b |`) beginnt nicht mit Leerzeichen und `|` und
    bleibt gemessen.
  - Ein Tab oder ein geschütztes Leerzeichen vor `|` lässt die Zeile gemessen, ebenso `> |` und eine Zeile
    ohne führende Pipe (`a | b`).
  - Mit `exempt-code-blocks: false` und `exempt-tables: true` liegt eine eingerückte `|`-Zeile im Zaun nicht
    außerhalb von Fenced-Code und wird gemessen, wie die Spalte-0-Zeile.
  - Eine Trennzeile wird weiter nicht verlangt. Eine Prosa-Zeile, die nach ein bis drei Leerzeichen mit `|`
    beginnt, ohne dass eine Tabelle dasteht, ist ausgenommen, wie heute in Spalte 0.
- Negativ — die Zusage der Anforderung „dieselbe Antwort, die die übrigen Module auf diese Frage geben, keine
  zweite“ und ihr Out-of-Scope-Punkt (4) gelten für die Tabellenzeile nicht mehr wörtlich und werden
  umgeschrieben: Die Antwort kommt weiter aus der geteilten Lexik, auf eine eigene, benannte Frage.
- Lastenheft §5 bleibt wahr, weil der Helfer in der geteilten Lexik wohnt. Wanderte er in `linelength`, wäre
  die dortige Zusage berührt, und das wäre eine Entscheidung des Auftraggebers.
- [ADR-0114](0114-linelength-modul.md) bleibt `Accepted`; ihre Zeile im Index führt
  „teil-superseded: ADR-0116“, und ihre `## Geschichte` bekommt eine Zeile, wenn diese ADR `Accepted` wird.
- Folgepflicht: jede Grenze oben steht mit ihrer Bedingung in der Anforderung; jeder Test der Anforderung läuft
  ohne die Regel aus dem richtigen Grund rot.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Kopplungs-Test der Tabellenzeilen-Frage | `targets`, `planning.waves` und `structure` antworten auf eine Zeile mit ein bis drei Leerzeichen vor `\|` „keine Tabellenzeile“, `linelength` „ausgenommen“; auf vier Leerzeichen, Tab und `>` antworten alle vier „gemessen“ bzw. „keine“; rot, wenn `tableRowLine` geweitet wird oder der Helfer `tableRowLine` nicht mehr einschließt | `make test` |
| Tests des Moduls | drei Leerzeichen an und aus; die Grenze bei drei gegen vier; Tab, U+00A0, `> \|`, `a \| b`; die eingerückte Zeile im Zaun; die Bezugsmenge | `make test` |
| Probe der Welle am gebauten Image | die Tabelle des Issues unter `1. `: Default kein Befund, mit `exempt-tables: false` ein Befund; vier Leerzeichen, Tab und `>` je ein Befund | — (Closure der Welle) |

## Re-Evaluierungs-Trigger

- Der Bestand des Melders behält nach der Hebung Befunde an Tabellen mit vier oder mehr Leerzeichen in
  Listenpunkten — dann Alternative D als eigener Vorgang.
- Die geteilte Lexik bekommt eine Antwort auf den Listenkontext — dann ist der Schnitt bei drei Leerzeichen
  neu zu prüfen.
- Ein zweiter Konsument will die eingerückte Tabellenzeile — dann bekommt die Frage des Helfers ihre eigene
  Kopplung (ADR-0054 Entscheidung 4) und die Spalte-0-Konsumenten werden einzeln gefragt.
- Jemand verlangt `> |` als Tabellenzeile.
- Eine Prosa-Zeile mit ein bis drei Leerzeichen vor `|` wird als zu Unrecht ausgenommen gemeldet.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-11-24 | Proposed: Entscheide 1–3 wie vom Auftraggeber bestätigt, Entscheid 4 aus dem Architect-Verdikt zur Welle | slice-293 <!-- d-check:status-provenance --> |
| 2026-11-24 | Accepted: der Mechanismus steht (slice-294). Der Helfer `unbreakableTableRow` in der geteilten Lexik ruft `tableRowLine` und lässt davor bis drei Leerzeichen (U+0020) zu, außerhalb von Fenced-Code; `tableRowLine` ist zeichengleich, `linelength` ruft für (b) nur den Helfer, und der Kopplungs-Test hält beide Fragen. Review R1 ohne HIGH und MEDIUM, die Verifikation hat jede Grenze mit einer Mutation rot gesehen; die Entscheide 1–4 tragen ohne inhaltliche Änderung gegenüber der Fassung als `Proposed` | slice-294 <!-- d-check:status-provenance --> |
