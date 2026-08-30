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

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-08-30 | **Die tragende Messung von §Entscheidung 2 ist zur Hälfte falsch — die Entscheidung steht, ihr Preis-Urteil wird schwächer.** Dort steht, `--print-mk` verteile *„kein Doctor-Target an Konsumenten"*, und daraus folgt der Preis sei **null**. Gemessen mit **einem** Produktlauf (`d-check --print-mk`): das Fragment trägt `.PHONY: doc-doctor` samt `--doctor`-Rezept, und [`DC-FA-CLI-010`](../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben) schreibt dieses Target in ihrem **Happy Path** ausdrücklich vor — die Behauptung widerspricht also einem Akzeptanzkriterium desselben Repos. Der Adopter **hat** das Target. Was bleibt: ein Target ist kein Gate-Lauf, und ob er es in seine Kette hängt, ist unbekannt und wurde nicht gefragt. Der Preis ist damit nicht *null*, sondern **unbestimmt**, und der verbliebene Grund für die gewählte Form ist der schwächere Satz *„die Zahl hängt an keiner Unbekannten"*. §Entscheidung 1 ist davon **unberührt** — sie steht auf dem anderen Bein: eine Erlaubnis kann eine veraltende Aufzählung überhaupt nicht laut machen. Von Review und Verifikation unabhängig gefunden; Klasse [`BEO-020`](../planning/observations.md) — gemessen wurde die eigene Menge (die Gate-Dateien dieses Repos), ausgesagt wurde über die fremde (das verteilte Fragment) |
| 2026-08-30 | **Zwei Eigenschaften des Schlüssels waren nirgends deklariert.** Erstens gilt die Zahl **je Datei**, nicht je Regel: eine Regel über einen Glob hält sie gegen **jede** getroffene Datei einzeln. Für den Bestand, um den der Antrag gestellt wurde — eine Datei —, ist das der gemeinte Fall; über einen Mehr-Datei-Glob ist es fast nie die gemeinte Aussage. Zweitens **bricht der Befund die Datei ab**: was hinter dem Mismatch stünde, wird nicht mehr geprüft, eine falsche Zahl verdeckt also den Rest, bis sie stimmt. Beides steht seither in Lastenheft, Spezifikation und Handbuch, und die Reichweite ist ein eigenes Akzeptanzkriterium |
| 2026-08-30 | **Die Fitness Function versprach mehr, als ihre Proben hielten.** Sie sagt *„jede Zusage ist von der Mutation gefangen, gegen die sie steht"*; zwei Mutationen liefen grün durch. (a) Die Zuweisung `ExemptExpectCount` im Config-Adapter löschen — das Feature ist im Produkt vollständig wirkungslos, `make test` Exit 0: die Regel-Tests bauen ihre Regel direkt und erreichen den Adapter nie. (b) Die Zählung auf `> 0` einschränken — die deklarierte Null wird von *„nicht deklariert"* ununterscheidbar, also genau die Zusage aus §Entscheidung 5 entwertet, ebenfalls grün. Beide Lücken sind mit je einer Zusage geschlossen: die Verdrahtung wird im Adapter-Test mitgeprüft, und die Null-Semantik hat eine Gegenprobe (null deklariert, einer ausgenommen ⇒ Mismatch). Beide Mutationen sind danach rot. Klasse [`BEO-023`](../planning/observations.md) |
| 2026-08-30 | **Die Byte-Identitäts-Zahl nennt jetzt ihren Korpus.** *„169 Befunde, `diff` leer"* stand ohne Prüfmenge und Profil da und war damit nicht nachvollziehbar. Der Lauf ist: `max-tasks: 3` mit dem Selektor auf den DoD-Abschnitt über den Glob der Slice-Pläne; die 169 zerfallen in 86 `section-missing` und 83 `section-oversized`. Die Verifikation hat ihn gegen ein aus `a062fe8` gebautes Vorgänger-Image reproduziert — **exakt**, dazu in drei weiteren Prüfmengen und vier Ausgabeformen |
| 2026-08-30 | **Ein Zitat in §Entscheidung 3 reicht über seinen Geltungsbereich.** Dort steht, [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) schreibe *„für genau das einen eigenen Code vor"*. Der Grundsatz *jede mit eigenem Grund-Code* gilt dort den **Bedingungen im Abschnitt**; die Ventil-Familie, zu der der neue Schlüssel gehört, trägt nach derselben Anforderung ausdrücklich **keinen** eigenen Code. Die Entscheidung bleibt richtig — sie steht auf der abweichenden **Reparatur**, nicht auf einem Satz, der von anderen Schlüsseln handelt. Klasse [`BEO-012`](../planning/observations.md) |
| 2026-08-30 | **Der erste Re-Evaluierungs-Trigger ist geprüft und nicht gezogen — durch eine Messung der Gegenseite.** Er lautet *„wenn ein Adopter `--doctor` in einem Gate fährt"*. Nachdem die Korrektur oben die Frage überhaupt erst stellbar machte, hat der Adopter auf seiner Seite gemessen: Makefile-Aggregate **0**, Workflows **0**, Hooks **0**; `doc-doctor` ist dort **advisory**. Die Antwort ist **nein**, die Abweichung von seiner beantragten Form bleibt richtig — und sie steht damit zum ersten Mal auf einer Messung **seines** Bestands statt auf einer Vermutung darüber. **Das heilt den Fehler nicht:** dass die unbelegte Behauptung im Ergebnis zutraf, macht sie nicht nachträglich zu einem Beleg; der Zähler bleibt erhöht. Der Trigger bleibt stehen — ein Repo kann `doc-doctor` jederzeit verdrahten |
