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

Eure beiden Spec-Zitate sind nachgelesen: **im Geltungsbereich**, und das
zweite wortgleich. Das erste — die `Boundary` von
[`DC-FA-CLI-007`](../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus) —
endet bei euch einen Halbsatz früher als im Original (dort folgt noch *„ohne
Fix-Kandidaten"*). Am Sinn ändert das nichts; wir nennen es, weil wir es
gemessen haben und ihr sonst ein Zertifikat bekämt, das eine Stelle nicht
deckt.

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

## Abweichung 2 — die Sichtbarkeit bleibt im Gate-Lauf

Ihr wählt `--doctor` als Träger und nennt den Preis *„schwächer"*. Er ist nicht
schwächer, sondern **unbestimmt** — und diese Zeile hat zwei Fassungen, weil
die erste falsch war.

| Ort | `--doctor` | gemessen wie |
|---|---|---|
| `Makefile` (alle Targets) | 0 Fundstellen | `grep` über das Makefile |
| `.github/workflows/` | 0 | dito |
| `.githooks/pre-commit` | 0 | dito |
| `--print-mk` (verteiltes Fragment) | **`doc-doctor` ist da** | `d-check --print-mk`, Produktlauf |

**Die vierte Zeile hieß zuerst „kein `doc-doctor`" — das war unbelegt und ist
falsch.** Unser eigener Review hat es mit einem Produktlauf widerlegt: das
Fragment trägt `doc-doctor`, und
[`DC-FA-CLI-010`](../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
schreibt das Target ausdrücklich vor. Ihr habt es also. Wir hatten unsere
eigenen Gate-Dateien gemessen und über euer Fragment geredet.

**Was nach der Korrektur bleibt, ist schwächer als das, was da stand.** Ein
Target ist kein Gate-Lauf: ob ihr `doc-doctor` in eure Kette hängt, wissen wir
nicht und haben wir nicht gefragt. Die Sichtbarkeit in `--doctor` hängt damit an
einer Bedingung, die wir nicht kennen; die Zahl im Gate-Lauf hängt an keiner.
**Das ist der ganze verbliebene Grund** — nicht mehr *„eure Form erreicht
niemanden"*, sondern *„eure Form erreicht jemanden, wenn ihr sie verdrahtet
habt, und unsere ohnehin"*.

**Fahrt ihr `doc-doctor` in einem Gate, trägt eure Form**, und dann ist diese
Abweichung eine Bevormundung statt einer Schärfung — sagt uns das, dann ist es
ein Re-Evaluierungs-Trigger und kein Streit.

## Abweichung 3 — es gibt doch einen neuen Grund-Code

Euer Vertrag sagt *„kein neuer Grund-Code"*. Das galt **eurer** Form, die nur
unterdrückt. Eine **Zahl** kann nicht stimmen, und dieser Zustand verlangt eine
andere Reparatur: *Aufzählung oder Zahl nachziehen* statt *Selektor
korrigieren*. Die Befund-Deduplikation läuft über (Datei, Zeile, Regel, Ziel,
Grund) — zwei Bedeutungen unter einem Code fielen zusammen. Er heißt
**`section-exempt-mismatch`**.

**Zur Genauigkeit, weil wir euch sonst ein zu starkes Argument gäben:** der
Grundsatz *jede mit eigenem Grund-Code* steht in
[`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
bei den **Bedingungen im Abschnitt**. Die Ventil-Familie, zu der der neue
Schlüssel gehört, trägt dort ausdrücklich **keinen** eigenen Code. Wir berufen
uns also nicht auf einen Satz, der von anderen Schlüsseln handelt, sondern auf
die abweichende **Reparatur**.

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
- **Kein Schweregrad in `model.Finding`.**
  [`DC-FA-CLI-003`](../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes)
  führt differenzierte Exit-Codes als Out-of-Scope; ein Schweregrad wäre ein
  eigener Entscheid des Werkzeugs.
- **Fence-Treue** unverändert.
- **Euer Argument gegen „die Regel weglassen"** trägt und steht in der ADR.

## Was wir nicht tun

**Wir schalten die Fähigkeit nicht über unseren eigenen Bestand scharf** — dieses
Repo führt heute keine Regel mit `exempt-section-pattern`.

Und die ehrliche Zeile: **eine erklärte Menge zu verkleinern bleibt eine
Lockerung**, jetzt eine mit Zahl. Wer die Zahl mitzieht, ohne die Aufzählung zu
prüfen, hat einen Wächter, der nur noch sich selbst bestätigt. Die Zahl
verschiebt die Pflege; sie nimmt sie nicht ab.
