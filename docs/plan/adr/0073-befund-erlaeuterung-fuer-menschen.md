# ADR-0073: Die Erläuterung eines Befunds erreicht den Menschen — und eine `structure`-Regel darf sie verfassen

**Status:** Accepted

**Datum:** 2026-08-29

**Autor:** pt9912

**Bezug:** [`DC-FA-CLI-004`](../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)
(das Zeilen-Format, das hier eine vierte Spalte bekommt),
[`DC-FA-CLI-007`](../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)
(die Diagnose und der Fix-Kandidat, gegen den abzugrenzen ist),
[`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(das Modul, das den neuen Schlüssel bekommt),
[`SPEC-001`](../../../spec/spezifikation.md#spec-001--befund) (das Feld, das es
längst gibt),
[ADR-0070](0070-tabellen-klammer-und-spaltenliste.md) (dieselbe Bewegung: eine
Frage, ein Ort — statt zweier paralleler Slots)

**Schärft:** [`DC-FA-CLI-004`](../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate),
[`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)

**Regeln:** Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v5.12.0/regelwerk/modul-04-adrs.md).

---

## Kontext

**Ein Befund sagt heute, was falsch ist, aber nicht, was der Regel-Autor
wollte.** `section-forbidden` heißt „hier steht eine verbotene Wendung". Welche
Zusage die Regel hütet und was der Leser jetzt tun soll, steht bestenfalls als
Kommentar in der Konfiguration — an einem Ort, den niemand öffnet, wenn ein
Gate rot wird.

**Der Grund-Code kann das nicht tragen, und das ist kein Versäumnis.**
[`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
führt zu jeder Bedingung eine **Reparatur**-Spalte; für `forbid-pattern` lautet
sie *„die Wendung ersetzen"*. Der Grund-Code sagt die **Art** des Defekts, und
**eine** Art bedient viele Regeln. Jede `forbid`-Regel dieses Repos bekäme
denselben Satz, obwohl die eine einen erfundenen Risiko-Ausgang hütet und die
andere einen offenen DoD-Haken.

**Beim Bau gemessen, und es hat den Entwurf verschoben.**
[`SPEC-001`](../../../spec/spezifikation.md#spec-001--befund) führt seit jeher
`message` als *„menschenlesbare Erläuterung (nicht stabilitätsgarantiert)"*.
**22 von 31** Regel-Dateien unter `internal/hexagon/core/rules/` setzen es —
und gerendert wurde es für Menschen **nirgends**: weder in der Befund-Zeile
noch in `--doctor`. Es erreichte ausschließlich `--json`/`--yaml`, also den
Maschinen-Konsumenten. Wer `make doc-check` rot laufen ließ, sah die
Erläuterung nicht, die das Produkt für ihn geschrieben hatte.

**Der Anlass ist ein zurückgeführter Slice.** Ein Sensor über die
DoD-Häkchen abgeschlossener Slices war fertig entworfen und beidseitig
gemessen, als sein eigenes Vorgehen auf genau diese Grenze traf: es verlangte
eine Meldung, die sagt, *was zu tun ist*.

## Entscheidung

**Drei Entscheide, die zusammengehören.**

1. **`message` wird gerendert.** Vierte tab-getrennte Spalte der Befund-Zeile,
   **nur wenn gefüllt**; eigene `Hinweis:`-Zeile in `--doctor` unter `Stelle:`.
   `--json`/`--yaml` bleiben unverändert — dort stand das Feld schon.
2. **Eine vierte Spalte, keine Fortsetzungszeile.**
   [`DC-FA-CLI-004`](../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)
   sagt **ein Befund pro Zeile** zu, und eines seiner Akzeptanzkriterien
   **zählt Zeilen** (*„genau zwei Befund-Zeilen"*). Eine Fortsetzungszeile
   bräche beides. Eine vierte Spalte bricht es nicht: wer auf Tab trennt und
   die Felder 1–3 liest, liest weiter dasselbe. **Nur wenn gefüllt**, damit ein
   Befund ohne Erläuterung byte-identisch bleibt.
3. **`structure[].hint` schreibt `message`, statt ein zweites Feld zu
   eröffnen.** Ein zweiter Slot für dieselbe Frage wäre die Redundanz, die
   [ADR-0070](0070-tabellen-klammer-und-spaltenliste.md) an anderer Stelle
   zurückgebaut hat. Ein explizit leerer Wert ⇒ Exit 2: ein leerer Hinweis sagt
   nichts zu — dieselbe Härte, die `planning.closure.boilerplate` für den
   leeren Eintrag führt.

**Vorrang, ausdrücklich:** setzt die Bedingung selbst schon ein `message`, so
**gewinnt der `hint`**. Die modul-eigene Meldung ist die Aussage des Werkzeugs
über die Art des Defekts; der `hint` ist die Zusage des Konfigurations-Autors,
und die steht näher am Leser.

**Zwei Befunde sind ausgenommen, und das ist die interessante Hälfte:** der
unlesbare Dateibaum und die leer laufende Regel. Sie verletzen **keine
Bedingung** — die Regel hat dort gar nicht gemessen. Ein Hinweis auf die
gehütete Zusage führte in die Irre, denn die Zusage steht gar nicht zur
Debatte; was fehlt, ist die Messung. Sie behalten ihre eigene Meldung.

**Abgrenzung gegen `fixCandidate`.** Der Fix-Kandidat ist **abgeleitet** und
nur dort, wo er *eindeutig ableitbar* ist
([`DC-FA-CLI-007`](../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus));
er speist `--repair` und wird zu einem **anwendbaren Patch**. Ein `hint` ist
**verfasst** und wird nie angewendet. Beides in ein Feld zu legen hieße,
Autoren-Prosa in die Patch-Pipeline zu geben.

## Verglichene Alternativen

| Option | Pro | Contra | Verworfen, weil |
|---|---|---|---|
| **A: neues Feld `hint` neben `message`** | klare Herkunft je Feld | zwei Slots für dieselbe Frage; `--json` bekäme zwei Erläuterungs-Schlüssel | die Herkunft interessiert den Leser nicht — er will **einen** Satz. [ADR-0070](0070-tabellen-klammer-und-spaltenliste.md) hat dieselbe Doppelung schon einmal zurückgebaut |
| **B: Hinweis nur in `--doctor`** | Zeilen-Format unberührt | wer rot läuft, sieht ihn erst nach einem zweiten Lauf mit anderem Schalter | ein Hinweis, den man suchen muss, ist bei einem roten Gate kein Hinweis |
| **C: Fortsetzungszeile unter dem Befund** | mehr Platz für Text | bricht *„ein Befund pro Zeile"* und ein zählendes Akzeptanzkriterium | die Zusage ist älter als der Wunsch |
| **D: `fixCandidate.note` nutzen** | Feld existiert, `--doctor` rendert es | verfasster Text landete in der Ableitung, die `--repair` zu einem Patch macht | ein Patch aus Autoren-Prosa ist die Fehlerklasse, gegen die die Trennung gebaut ist |
| **E: Grund-Codes verfeinern** | keine neue Fläche | ein Code je Zusage: der Code-Raum wüchse mit jeder Konfiguration | Grund-Codes sind ein **geschlossener** Vertrag (§4 der Spezifikation), keine Konfigurations-Fläche |

## Konsequenzen

**Positiv.** Die Erläuterung erreicht sofort **22 Regel-Dateien**, nicht nur
die neue Zusage — der Gewinn ist größer als der Anlass. Eine `structure`-Regel
kann sagen, welche Zusage sie hütet, ohne dass der Grund-Code-Raum wächst. Der
Konfigurations-Kommentar bleibt, was er ist: die Begründung für den Leser der
Konfiguration, nicht der einzige Ort der Reparatur-Anweisung.

**Negativ, und benannt.** Die Ausgabe-Zusage wird **geweitet**: ein Konsument,
der auf genau drei Tab-Felder besteht, sieht bei Befunden mit Erläuterung vier.
Gemessen: ein grüner Lauf ist unverändert (keine Befund-Zeile), ein Befund ohne
`message` ist unverändert, ein Befund mit `message` gewinnt die Spalte —
`links` `target-missing` etwa trägt ab jetzt *„Linkziel existiert nicht"*.

**Der Text altert wie ein Kommentar.** Ob ein `hint` noch stimmt, prüft kein
Gate; [`SPEC-001`](../../../spec/spezifikation.md#spec-001--befund) sagt für
`message` ohnehin **keine Stabilität** zu. Das ist
die Klasse, die dieses Repo als
[`BEO-013`](../planning/observations.md) führt — hier ohne Wächter, mit Absicht:
die Alternative wäre, den Hinweis zu verbieten, und dann steht wieder nichts da.

**Das Feld lädt zur Ausrede ein.** Ein schlecht benannter Grund-Code lässt sich
mit einem Hinweis zudecken, statt geschärft zu werden. Dagegen steht die
Arbeitsteilung aus der Entscheidung — Code = Art, Hinweis = Zusage —, nicht ein
Sensor.

## Fitness Function (falls maschinell prüfbar)

**Ja, dreifach, und gefahren:**

1. `TestStructureHintSchreibtMessageUndGewinnt` (Kern) — der `hint` schreibt
   `message` und gewinnt gegen die modul-eigene Meldung; der Grund-Code bleibt
   unberührt.
2. `TestStructureHintGiltNichtFuerLeerlaufendeRegel` (Kern) — die benannte
   Grenze: ein Befund ohne Bedingungs-Verletzung behält seine eigene Meldung.
3. `TestHint_VierteSpalteInDerBefundZeile`, `TestHint_DoctorZeigtHinweis`,
   `TestHint_JSONUnveraendert` (CLI, end-to-end) — vier Tab-Spalten in **einer**
   Zeile, die `Hinweis:`-Zeile in der Diagnose, `--json` unverändert.

Der Config-Rand (`hint` gesetzt, aber leer bzw. nur Whitespace ⇒ Exit 2) liegt
in `TestDecode_StructureFehler`, die Durchreichung in
`TestDecode_StructureHint`.

## Re-Evaluierungs-Trigger

**Wenn ein zweites Modul eine verfasste Erläuterung braucht.** `structure` ist
das einzige, dessen Konfiguration eine **Regel-Liste** mit benannten Zusagen
führt; die übrigen leiten ihre Befunde aus dem Dokument ab. Meldet sich ein
zweiter Kandidat, ist zu entscheiden, ob `hint` ein modul-übergreifender
Schlüssel wird — dann gehört er nicht in `structure[]`.

**Wenn die vierte Spalte einen Konsumenten bricht.** Dann ist Option B
(`--doctor`-only) die Rückfallposition, und diese ADR bekommt eine Nachfolgerin.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-08-29 | **Entscheidung 1 ergänzt: der Reporter saniert.** Der Review fand die Zusage „ein Befund bleibt eine Zeile" **brechbar** — ein `hint` mit `\n` erzeugte zwei Zeilen, einer mit `\t` fünf Felder; gemessen. Die Erläuterung kommt aus der Konfiguration **oder aus dem geprüften Material** (`commits` trägt den Commit-Betreff), also aus einer Eingabe, über die der Reporter keine Kontrolle hat (`AGENTS.md` §3.8). Zwei Antworten, bewusst beide: die Konfiguration **weist ab**, was ein Autor anders gemeint hätte als das Ergebnis; der Reporter **ersetzt** Tab, Zeilenumbruch und Wagenrücklauf durch ein Leerzeichen, weil der modul-eigene Weg keinen Config-Rand hat. Drei Mutationsproben gefahren: bedingungslose vierte Spalte, entfernte Sanierung, `hint` auch für Nicht-Mess-Befunde — jede wird jetzt von genau einem Test gefangen |
| 2026-08-29 | **Zahl im Kontext präzisiert.** Dort steht „22 von 31 Regel-Dateien"; das zählt **Dateien** mit dem Literal `Message:`, darunter acht geteilte Helfer ohne eigene Regel. Der reichweiten-relevante Schnitt lautet **21 von 23** Regel-Einstiegs-Dateien (`func Check*`); ungedeckt sind genau `hostpaths` und `spans`. Der Schluss — der Gewinn ist größer als der Anlass — trägt in beiden Zählungen; die Zahl selbst führte einen späteren Leser zu `run.go`/`scan.go` statt zu den zwei Regeln |
