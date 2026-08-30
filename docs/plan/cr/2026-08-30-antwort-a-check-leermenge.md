# Antwort auf den `a-check`-CR 4 vom 2026-08-30 — die legitime Leermenge einer erklärten Teilmenge

**Absender der Antwort:** d-check
**Datum:** 2026-08-30
**Bezug:** [eingehender CR 4](2026-08-30-cr-a-check-leermenge.md); Vorgänger
[CR 3](2026-08-30-cr-a-check-structure-teilmenge.md) samt
[Antwort](2026-08-30-antwort-a-check-structure-teilmenge.md)
**Ergebnis:** **in der Sache angenommen, in der Form nicht** — drei Abweichungen,
jede mit ihrer Messung
**Einordnung:** MINOR — Lastenheft `0.78.0`, Begründung in
[ADR-0078](../adr/0078-erklaerte-leermenge-mit-zahl.md)

---

## Vorab

Der Befund trägt, und er trägt genau so, wie ihr ihn stellt: die
Nullmengen-Härte aus [ADR-0075](../adr/0075-erklaerte-teilmenge-in-structure.md)
wirft zwei verschiedene Zustände in einen Grund-Code. Eure Zwei-Zustände-Tabelle
ist die richtige Zerlegung, und eure Eingrenzung — die Härte trifft ein
**generisches** Muster, nicht ein **aufzählendes** — ist die richtige
Einschränkung. *„Das Modul macht mehr rot als der Sensor, den es ablöst"* ist
ein Bruch, und er ist behoben.

Und: dass ihr euren ersten Entwurf selbst korrigiert habt, statt ihn zu
verteidigen, hat diese Runde kurz gemacht. Die beiden Spec-Zitate der zweiten
Fassung sind nachgelesen — **wortgleich und im Geltungsbereich**.

Drei Dinge sind trotzdem anders geworden, als ihr sie beantragt habt.

## Abweichung 1 — eine Deklaration statt einer Erlaubnis

Statt `exempt-may-empty: true` gibt es **`exempt-expect-count: <n>`**.

```yaml
exempt-section-pattern: '^### (AC-FA-RULE-…|…)'
exempt-expect-count: 19
```

| Lage | Ergebnis |
|---|---|
| 19 Abschnitte, 19 ausgenommen, 19 deklariert | **kein Befund**, Exit 0 |
| 20 Abschnitte, 19 ausgenommen, 19 deklariert | der zwanzigste wird normal geprüft |
| 19 Abschnitte, 18 ausgenommen, 19 deklariert | **`section-exempt-mismatch`** |
| `section-pattern` trifft nichts | **`section-missing`** — bleibt ein Defekt |

**Der Grund ist eure eigene Risiko-Aussage.** Ihr schreibt, das Muster *„kann
nur veralten"*. Genau dieser Fall ist mit `exempt-may-empty` **stumm** und mit
einer Zahl **laut** — und zwar dort, wo ihr ohnehin hinseht.

## Abweichung 2 — die Sichtbarkeit bleibt im Gate-Lauf, und der Grund ist gemessen

Ihr wählt `--doctor` als Träger und nennt den Preis *„schwächer"*. **In diesem
Repo ist er null**, gemessen:

| Ort | `--doctor` |
|---|---|
| `Makefile` (alle Targets) | 0 Fundstellen |
| `.github/workflows/` | 0 |
| `.githooks/pre-commit` | 0 |
| `--print-mk` (verteilte Targets) | kein `doc-doctor` im Gate-Pfad |

Ein Modus, den kein Gate fährt, erreicht nur den, der ohnehin nachsieht — und
für einen **Bestandszustand** sieht niemand nach. Eure Einordnung des dritten
Weges als *„ohne Vertragsbruch"* stimmt; nur ist der Preis größer, als *„etwas
schwächer"* nahelegt.

**Das ist eine Aussage über dieses Repo, nicht über eures.** Fahrt ihr
`--doctor` in einem Gate, trägt eure Form, und dann ist diese Abweichung eine
Bevormundung statt einer Schärfung — sagt uns das, dann ist es ein
Re-Evaluierungs-Trigger und kein Streit.

## Abweichung 3 — es gibt doch einen neuen Grund-Code

Euer Vertrag sagt *„kein neuer Grund-Code"*. Das galt **eurer** Form, die nur
unterdrückt. Eine **Zahl** kann nicht stimmen, und dieser Zustand verlangt eine
andere Reparatur: *Aufzählung oder Zahl nachziehen* statt *Selektor
korrigieren*. [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) schreibt dafür einen eigenen Code vor — die
Befund-Deduplikation läuft über (Datei, Zeile, Regel, Ziel, Grund), zwei
Bedeutungen unter einem Code fielen zusammen. Er heißt
**`section-exempt-mismatch`**.

## Was ihr bekommt, das ihr nicht beantragt habt

- **Die Prüfung läuft immer, nicht nur bei leerer Restmenge.** Ein zwanzigster
  Abschnitt kommt dazu und ihr zieht die Zahl nicht — das meldet. Eine
  einseitige Prüfung wäre ein halber Wächter.
- **Beidseitige Drift.** Mehr ausgenommen als deklariert meldet ebenso wie
  weniger.
- **`0` ist erlaubt und bedeutet etwas** — *„das Muster soll heute noch nichts
  treffen"*, für einen Bestand, der erst wächst. Deshalb ist der Schlüssel ein
  Zeiger: die deklarierte Null ist von *„nicht deklariert"* unterscheidbar.
- **Der Befund trägt euren `hint`.** Anders als der Nullmengen-Befund ist er
  nicht davon ausgenommen: dort hat die Regel nicht gemessen, hier hat sie eure
  Deklaration widerlegt — und ihr dürft erklären, was zu tun ist.

## Zwei Fragen, die euer Antrag offenlässt, und ihre Antwort

**Geht der Schlüssel in die Regel-Identität ein?** Nein. `exempt-section-pattern`
tut es (zwei Regeln mit verschiedener Ausnahme sind verschiedene Zusagen); zwei
**erwartete Zahlen** über derselben Regel sind dagegen ein **Widerspruch** —
eine davon muss falsch sein. Sie als Konfigurations-Duplikat abzuweisen ist die
richtige Antwort.

**Greift der Schlüssel auch, wenn schon der Selektor nichts trifft?** Nein — und
das ist eure eigene Zerlegung, mechanisch gemacht.

## Angenommen, unverändert

- **Ein optionaler Schlüssel an derselben Bedingung**; ohne ihn
  **byte-identisches** Verhalten. Gemessen gegen das Vorgänger-Image: 169
  Befunde, `diff` leer.
- **Nichts für `exempt-paths`** — eure Abgrenzung ist die richtige: ein
  Datei-Glob ist generisch, eine Abschnitts-Aufzählung in einer Datei nicht.
- **Kein Schweregrad in `model.Finding`.** Ihr habt das selbst zurückgezogen;
  wir bestätigen es: [`DC-FA-CLI-003`](../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes) führt differenzierte Exit-Codes als
  Out-of-Scope, und ein Schweregrad wäre ein eigener Entscheid des Werkzeugs.
- **Fence-Treue** unverändert.
- **Euer Argument gegen „die Regel weglassen"** trägt und steht in der ADR.

## Was wir nicht tun

**Wir schalten die Fähigkeit nicht über unseren eigenen Bestand scharf** — dieses
Repo führt heute keine Regel mit `exempt-section-pattern`.

Und die ehrliche Zeile: **eine erklärte Menge zu verkleinern bleibt eine
Lockerung**, jetzt eine mit Zahl. Wer die Zahl mitzieht, ohne die Aufzählung zu
prüfen, hat einen Wächter, der nur noch sich selbst bestätigt. Die Zahl
verschiebt die Pflege; sie nimmt sie nicht ab.
