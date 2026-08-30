# ADR-0078: Die erklärte Leermenge bekommt eine Zahl, keine Erlaubnis

**Status:** Accepted

**Datum:** 2026-08-30

**Autor:** pt9912

**Bezug:**
[der eingehende CR 4](../cr/2026-08-30-cr-a-check-leermenge.md) (Antrag und
Beleg des Absenders),
[ADR-0075](0075-erklaerte-teilmenge-in-structure.md) (die Nullmengen-Härte,
deren Reichweite hier zurückgeschnitten wird),
[ADR-0073](0073-befund-erlaeuterung-fuer-menschen.md) (die `hint`-Ausnahme, an
der sich der neue Befund abgrenzt),
[`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(die erweiterte Anforderung),
[`DC-FA-CLI-003`](../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes)
(die binäre Exit-Semantik, an der die beantragte Form scheiterte)

**Schärft:** [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)

**Regeln:** Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v5.12.0/regelwerk/modul-04-adrs.md).

---

## Kontext

[ADR-0075](0075-erklaerte-teilmenge-in-structure.md) hat `exempt-section-pattern`
eine **Nullmengen-Härte** gegeben: nimmt die Ausnahme alle Abschnitte, meldet
die Regel `section-missing`, statt still grün zu werden.

**Die Härte trifft zwei verschiedene Zustände mit einem Grund-Code.**

| Zustand | Bedeutung |
|---|---|
| `section-pattern` trifft **nichts** | Konfigurationsdefekt — falsches Muster, falsche Datei, umbenannter Abschnitt |
| Muster trifft, Ausnahme nimmt **alle** | Bestandszustand — *„es gibt noch nichts Neues zu prüfen"* |

Der Adopter, der ADR-0075 beantragt hat, ist beim **Anwenden** darauf gestoßen:
19 grandfatherte Anforderungen in einer Datei mit genau 19 passenden
Abschnitten. Die Regel meldet rot, wo sein abgelöstes Skript
`0 neue AC(s) geprueft, 19 grandfathered` mit Exit 0 meldete — **das Modul
macht mehr rot als der Sensor, den es ablösen soll**.

**Die Härte bleibt richtig, ihre Reichweite war zu weit.** ADR-0075 begründet
sie mit: *„Ohne diese Antwort schaltete ein zu breites Muster die Regel still
ab."* Das trifft ein **generisches** Muster. Ein **aufzählendes**, das 19
Kennungen einzeln nennt, kann nicht versehentlich zu breit werden — es kann nur
**veralten**.

## Entscheidung

**1. Eine Deklaration statt einer Erlaubnis.** `exempt-expect-count` (int ≥ 0,
opt-in, nur mit `exempt-section-pattern`) sagt, **wie viele** Abschnitte die
Ausnahme nehmen soll. Stimmt die Zahl, ist die geleerte Menge **kein** Befund.
Weicht sie ab, entsteht einer.

**2. Der Antrag lautete anders, und der Grund für die Abweichung ist gemessen.**
Der Absender beantragte `exempt-may-empty: true` und trug die Sichtbarkeit nach
`--doctor`; er nennt den Preis selbst *„schwächer"*. **Gemessen ist er null:**
`--doctor` läuft in **keinem** Gate dieses Repos — nicht im `Makefile`, nicht in
`.github/workflows/`, nicht im `pre-commit`-Hook —, und `--print-mk` verteilt
kein Doctor-Target an Konsumenten. Eine Zeile dort erreicht nur, wer ohnehin
nachsieht; für einen Bestandszustand sieht niemand nach. **Die Zahl hält die
Prüfung im Gate-Lauf** und macht genau das laut, was der Absender als einziges
Risiko benennt: *„es kann nur veralten"*.

**3. Ein neuer Grund-Code, entgegen dem Antrag.** Er sagt *„kein neuer
Grund-Code"* — das galt **seiner** Form, die nur unterdrückt. Eine **Zahl** kann
nicht stimmen, und dieser Zustand verlangt eine andere Reparatur: *Aufzählung
oder Zahl nachziehen* statt *Selektor korrigieren*.
[`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
schreibt für genau das einen eigenen Code vor, weil die Befund-Deduplikation
über (Datei, Zeile, Regel, Ziel, Grund) läuft. Er heißt
`section-exempt-mismatch`; die Form folgt `section-heading-mismatch`.

**4. Die Drift ist beidseitig und unabhängig von einer Restmenge.** Mehr
ausgenommen als deklariert ist ebenso ein Befund wie weniger, und geprüft wird
**immer**, nicht nur bei geleerter Menge. Eine erweiterte Aufzählung ohne
nachgezogene Zahl ist dieselbe Lücke wie eine veraltete; eine einseitige Prüfung
wäre ein halber Wächter. Genau der wachsende Bestand ist der Fall, den der
Antrag als einzigen benennt.

**5. Eine deklarierte Null bedeutet etwas** — *„das Muster soll heute noch nichts
treffen"*, ein Bestand, der erst wächst. Deshalb ist der Schlüssel ein **Zeiger**
und `0` erlaubt: sie von *„nicht deklariert"* zu unterscheiden ist der
Unterschied zwischen einer Aussage und ihrer Abwesenheit.

**6. Der Schlüssel greift nur nach dem Abzug.** Trifft schon `section-pattern`
nichts, bleibt es `section-missing` — auch mit gesetzter Zahl. Das ist die
Trennung aus §Kontext, und sie ist der ganze Gegenstand.

**7. Der neue Befund ist NICHT von `hint` ausgenommen.** Der Nullmengen-Befund
ist es ([ADR-0073](0073-befund-erlaeuterung-fuer-menschen.md): *„die Regel hat
nicht gemessen"*). Hier **hat** sie gemessen und eine Deklaration widerlegt; ein
verfasster Hinweis darf das erklären.

**8. Er geht NICHT in die Regel-Identität ein.** `exempt-section-pattern` tut es
(ADR-0075), weil zwei Regeln mit verschiedener Ausnahme verschiedene Zusagen
sind. Zwei **erwartete Zahlen** über derselben Regel sind dagegen ein
**Widerspruch**, kein Paar — eine davon muss falsch sein. Sie als
Konfigurations-Duplikat abzuweisen ist die richtige Antwort.

## Verglichene Alternativen

**`exempt-may-empty: true` wie beantragt.** Verworfen aus Entscheidung 2: die
Sichtbarkeit landete in einem Modus, den dieses Repo in keinem Gate fährt.

**Ein Schweregrad in `model.Finding`.** Der **erste** Entwurf des Absenders
hätte ihn stillschweigend eingeführt (*„eine nicht-fatale Zeile"*); es gibt ihn
nicht, und
[`DC-FA-CLI-003`](../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes) führt
differenzierte Exit-Codes ausdrücklich als **Out-of-Scope**. Der Absender hat
das nach einem Hinweis selbst korrigiert. Wenn das Werkzeug je einen Schweregrad
bekommt, ist das ein eigener Entscheid.

**Beide Schlüssel nebeneinander** (`exempt-may-empty` **und** die Zahl).
Verworfen: zwei Schlüssel für eine Frage sind die Verdopplung, die
[ADR-0070](0070-tabellen-klammer-und-spaltenliste.md) für die
Tabellen-Bedingungen gerade zurückgebaut hat — und ein dritter Config-Rand für
den Fall, dass jemand beide setzt.

**Die Regel weglassen, bis es etwas zu prüfen gibt.** Vom Absender selbst
verworfen, und das Argument trägt: eine Regel, die man erst einschalten muss,
wenn ihr Fall eintritt, ist kein Gate. Er hat dafür einen eigenen Beleg — ein
Modul, das eingebunden war und dreizehn Minor-Versionen ins Leere lief, weil es
nie konfiguriert wurde.

## Konsequenzen

**Positiv.** Der Bestandszustand ist vom Konfigurationsdefekt unterscheidbar,
und der Fall, den der Antrag fürchtet (*die Aufzählung veraltet*), wird **im
Gate-Lauf** laut statt in einem Diagnose-Modus. Ohne den Schlüssel ist der
Befundsatz byte-identisch — gemessen gegen das Vorgänger-Image: **169 Befunde,
`diff` leer**.

**Negativ, und das ist die ehrliche Seite.** **Die geprüfte Menge zu
verkleinern bleibt eine Lockerung** — jetzt eine mit Zahl, aber der
Bestandszustand ist ab dann stumm. Wer die Zahl mitzieht, ohne die Aufzählung
zu prüfen, hat einen Wächter, der nur noch sich selbst bestätigt. **Eine
deklarierte Zahl ist Autoren-Text und altert wie jede**; der Schlüssel
verschiebt die Pflege, er nimmt sie nicht ab.

**Neue Bauform ohne Präzedenz.** Gemessen führt kein anderer Schlüssel dieses
Moduls eine **erwartete Anzahl**. Ob die Form trägt, zeigt erst der zweite Fall.

**Ein Urteil über einen fremden Bestand.** Die `--doctor`-Messung aus
Entscheidung 2 gilt **diesem** Repo; ob der Adopter `--doctor` in einem Gate
fährt, ist nicht gefragt worden. Wäre es so, trüge seine Form, und die
Abweichung wäre Bevormundung statt Schärfung.

## Fitness Function (falls maschinell prüfbar)

`make test` — `internal/hexagon/core/rules/structure_leermenge_test.go` (sieben
Tests) und die zwei neuen Config-Ränder in `configyaml_test.go`. **Jede Zusage
ist von der Mutation gefangen, gegen die sie steht**, gemessen: Zählung
deaktiviert ⇒ 5 rot · Drift nur einseitig ⇒ 3 rot · `case`-Guard entfernt
(die Nullmengen-Härte feuert wieder) ⇒ 6 rot · Mismatch als `raw` ⇒ 1 rot ·
beide Config-Ränder entfernt ⇒ 1 rot.

**Eine Mutation war zunächst wirkungslos, und das war ein Befund am Code:** ein
`return nil` für die geleerte Menge war **toter Code** — der Rest der Funktion
liefert für null Abschnitte ohnehin nichts. Er ist entfernt; die Zusage hält,
weil der `case` die Nullmengen-Härte **überspringt**, nicht weil dort etwas
unterdrückt würde.

**Nicht maschinell geprüft** ist, ob eine gesetzte Zahl **berechtigt** ist. Das
ist ein Urteil über die Konfiguration eines fremden Repos und kein Gate.

## Re-Evaluierungs-Trigger

**Wenn ein Adopter `--doctor` in einem Gate fährt.** Dann trägt die beantragte
Form, und Entscheidung 2 ist neu zu lesen.

**Wenn eine deklarierte Zahl gemessen blind mitgezogen wurde.** Dann ist der
erste Punkt unter §Konsequenzen eingetreten, und die Frage lautet, ob die Zahl
allein reicht oder an die Aufzählung gebunden gehört.

**Wenn ein zweiter Schlüssel eine erwartete Anzahl braucht.** Dann ist die
Bauform aus Entscheidung 1 keine Einzelfall-Antwort mehr und verdient eine
gemeinsame Form.
