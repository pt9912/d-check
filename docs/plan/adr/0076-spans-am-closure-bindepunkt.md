# ADR-0076: `spans` am Closure-Bindepunkt — der Wächter über den Text, auf dem die anderen urteilen

**Status:** Accepted

**Datum:** 2026-08-30

**Autor:** pt9912

**Bezug:**
[ADR-0048](0048-closure-note-struktur-im-planning-modul.md) (der
Closure-Bindepunkt und sein zweites Prüf-Profil),
[ADR-0059](0059-closure-waechter-weicht-structure-regel.md) (der **ausgewiesene
Preis** der Bereinigung, den dieser Entscheid zur Hälfte einlöst),
[ADR-0042](0042-markdown-lexik-folgt-commonmark.md) (dessen offener Punkt (a)
hier erneut gemessen ist),
[`DC-FA-SPAN-001`](../../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
(das hinzukommende Modul),
[`DC-FA-CLI-012`](../../../spec/lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben)
(der Profil-Pfad)

**Schärft:** —. Diese Entscheidung ändert **keine** Produkt-Zusage; sie ändert,
welche Module dieses Repo an seinem eigenen Closure-Bindepunkt fährt.

**Regeln:** Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v5.12.0/regelwerk/modul-04-adrs.md).

---

## Kontext

`make verify-closure-notes` fährt `planning` und `structure` über das
`done/`-Profil. **Beide lesen den bereinigten Abschnitts-Text** — Fenced-Code
entfernt, Inline-Code geleert. Das ist richtig für Prosa-Bedingungen und der
Grund, warum [ADR-0059](0059-closure-waechter-weicht-structure-regel.md) den
Preis dieser Bereinigung ausdrücklich ausgewiesen hat.

**Der Preis hat eine zweite Hälfte, die dort nicht stand.** Verschluckt die
Bereinigung den Text, wird die Zusage darüber **still wahr**: kein Befund, Exit
0, und niemand hat etwas gesehen. Vier Proben, gegen den heutigen Stand
gefahren:

| Probe | `structure` allein | `structure` + `spans` |
|---|---|---|
| offener Haken, nackt | `section-forbidden` | ebenso |
| ungerader Backtick im Absatz | `section-forbidden` | ebenso **plus** `span-unclosed` |
| wohlgeformter Span, der das Item **umschließt** | still | still |
| vergessener Schluss-Fence | **still** | `fence-unclosed` |

**Die dritte Zeile ist kein Defekt.** Ein wohlgeformter Code-Span, der ein
Task-Item umschließt, macht das Item nach CommonMark zu **Code**. Die
Bereinigung hat recht; nur der Autor meinte etwas anderes. Kein Modul kann das
unterscheiden, weil am Dokument nichts falsch ist.

**Die vierte Zeile ist einer**, und sie ist die einzige, die heute niemand
sieht — genau am Bindepunkt, an dem die Closure-Zusage hängt.

**Die naheliegende größere Antwort ist gemessen und verworfen.** Der Verdacht
lautete, das Produkt brauche eine **einheitliche Fence-Lexik**: es führt zwei
Lesarten — den Toggle (`FenceToggle`, benutzt von `proseLines`,
`PreprocessMarkdown` und der Abschnitts-Findung) und die CommonMark-Lesart
(`FenceRun` + `FenceCloses`, benutzt von `spans` und dem Tabellen-Reader).
Beide über **alle 676** Markdown-Dateien dieses Repos gegeneinander gefahren:
**null Abweichungen**. Die Messung kann finden — drei konstruierte Fälle
(Infostring hinter dem Schluss-Fence, Zeichenwechsel `` ``` ``→`~~~`, zu kurzer
Schluss-Run) meldet sie, die saubere Datei nicht.
[ADR-0042](0042-markdown-lexik-folgt-commonmark.md) hat genau diesen Punkt als
*„bewusst offen gelassen … unbelegt — kein Realfall in den 522 Dateien"*
geführt und den Trigger benannt: erst wenn einer existiert. Er existiert nicht.

## Entscheidung

**1. `spans` läuft am Closure-Bindepunkt mit.** Das Rezept von
`verify-closure-notes` schaltet es neben `planning` und `structure` scharf. Es
prüft nicht die Closure-**Aussage**, sondern die **Lesbarkeit des Textes, über
den die beiden anderen urteilen** — der Wächter über die Wächter.

**2. Es steht im Makefile-Rezept, nicht in `modules:` des Profils.** Die leere
Liste ist eine Aussage: *„dieses Profil schaltet nichts von sich aus scharf"*.
Sie bleibt wahr, und **welche** Module laufen, steht an **einem** Ort. Zwei
Orte drifteten — und ein Go-Test hält bereits, dass das Profil keine Netz-/
Range-Tür öffnet; ihn um eine Positiv-Liste zu erweitern hieße, dieselbe
Aussage ein zweites Mal zu führen.

**3. Der Bestand ist vorher gemessen, nicht danach entdeckt.** `spans` meldet
über den ganzen Baum **0** Befunde; der Zuwachs am Bindepunkt ist rausch-frei.
Und die neue Deckung ist **diskriminierend** belegt: ein echter `done/`-Slice,
unverändert **grün** (Exit 0), mit angehängtem offenem Fence **rot**
(`fence-unclosed`, Exit 1).

**4. Die Fence-Lexik bleibt unangetastet.** Beide Lesarten bleiben, wie sie
sind. Der Anlass, den [ADR-0042](0042-markdown-lexik-folgt-commonmark.md) als
Bedingung benannt hat, ist nach erneuter Messung nicht eingetreten.

**5. Der umschließende Span bekommt keine Antwort, und das ist ausgesprochen.**
Er ist kein Defekt (§Kontext). Wer die Tabelle quer liest, könnte *„der
Bindepunkt sieht jetzt alles"* mitnehmen — er sieht eine Klasse mehr, nicht
alle.

## Verglichene Alternativen

**Die Fence-Lexik vereinheitlichen** (`proseLines` auf `FenceCloses`
umstellen). Verworfen: kein belegter Anlass in 676 Dateien, und die Änderung
träfe **jedes** scannende Modul — die Byte-Identität des Befundsatzes wäre die
zentrale Frage, für einen Fall, den niemand hat.

**`non-empty: true` an die bestehenden Closure-Regeln hängen.** Macht denselben
Fall laut (gemessen: `section-empty`), ist aber eine **Bedingungs**-Änderung an
Zusagen, die dieser Entscheid nicht führt — und sie deckte nur die Regeln, die
sie trägt, statt die Datei.

**`spans` in `modules:` des Profils.** Verworfen aus Entscheidung 2.

**Den Bindepunkt in `gates` ziehen**, damit der Defekt früher auffällt.
Verworfen: das ist die Bindepunkt-Trennung aus
[ADR-0048](0048-closure-note-struktur-im-planning-modul.md), und sie hat einen
eigenen Grund — eine Closure-Frage gehört nicht in den Inner-Loop. Der Preis
bleibt: der Defekt fällt erst bei `make fullbuild` auf.

## Konsequenzen

**Positiv.** Die still wahre Zusage ist am Closure-Bindepunkt nicht mehr
möglich, ohne dass etwas meldet. `spans` prüft dabei die **ganze** Datei, nicht
nur den Closure-Abschnitt — die Deckung ist breiter als der Anlass.

**Negativ.** Ein drittes Modul am Bindepunkt ist ein drittes, das rot werden
kann: ein künftiger Slice mit unbalanciertem Absatz blockiert die Closure,
statt nur die Prosa zu trüben. Das ist gewollt und trotzdem ein Preis.

**Die Messung ist eine Aussage über diesen Bestand.** *„0 Abweichungen in 676
Dateien"* gilt für dieses Repo und seine Schreibweise; ein Adopter kann die
Divergenz haben, und dieser Entscheid sagt über ihn nichts.

## Fitness Function (falls maschinell prüfbar)

`make verify-closure-notes` — der Bindepunkt selbst. **Nicht** maschinell
geprüft ist, dass er `spans` führt: das steht im Makefile-Rezept, und ein Test,
der ein Rezept gegen eine Liste hält, wäre die zweite Quelle aus Entscheidung 2.
`make gate-consistency` hält dafür den `##`-Hilfetext gegen die Doku.

## Re-Evaluierungs-Trigger

**Wenn eine Divergenz der beiden Fence-Lesarten in einem realen Bestand
auftritt.** Dann ist Entscheidung 4 hinfällig und der offene Punkt (a) aus
[ADR-0042](0042-markdown-lexik-folgt-commonmark.md) fällig — der Trigger ist
dort formuliert, hier ist er erneut gemessen und nicht eingetreten.

**Wenn `spans` am Bindepunkt Bestands-Befunde erzeugt, die niemand einplanen
wollte.** Dann ist zuerst der Bestand die Frage und nicht das Gate.
