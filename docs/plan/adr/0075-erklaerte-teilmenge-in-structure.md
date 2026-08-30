# ADR-0075: Eine `structure`-Regel erklärt ihre Grundmenge — zwei Muster, zwei verschiedene Ziele

**Status:** Accepted

**Datum:** 2026-08-30

**Autor:** pt9912

**Bezug:** [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(die erweiterte Anforderung),
[der eingehende CR](../cr/2026-08-30-cr-a-check-structure-teilmenge.md)
(Antrag und Beleg des Absenders),
[ADR-0070](0070-tabellen-klammer-und-spaltenliste.md) (die Präzedenz: eine
Option an einer bestehenden Bedingung statt eines neuen Moduls),
[ADR-0073](0073-befund-erlaeuterung-fuer-menschen.md) (die zuletzt so
verortete Option `hint`)

**Schärft:** [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)

**Regeln:** Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v5.12.0/regelwerk/modul-04-adrs.md).

---

## Kontext

`structure` wendet seine Bedingungen auf die **vorgefundene** Menge an: alle
Task-Items eines Abschnitts, alle Abschnitte, die der Selektor trifft. Ein
[eingehender CR](../cr/2026-08-30-cr-a-check-structure-teilmenge.md) eines
Adopters verlangt an zwei Stellen eine **erklärte** Teilmenge und legt seine
Messung vor: `max-tasks: 3` über seine Slice-Pläne liefert neun Befunde, acht
davon auf Plänen mit je sieben Task-Items, von denen zwei bis drei
Liefer-Punkte sind. Der Bestand ist regelkonform; der Zähler misst das
Falsche.

**Der Anlass gilt hier genauso, und er ist nachgemessen.** Dieselbe Regel über
die **86** Slice-Pläne dieses Repos, die einen `## 4. Definition of Done`
führen (**444** Task-Items): **80** `section-oversized`.

**Die Vorlage der Baseline liefert den Defekt mit aus.** Ihre DoD trägt neun
Checkboxen, von denen sechs pro Slice konstant sind, und sie schreibt eine
Zeile darüber selbst: *„Gate-Läufe und die vier Closure-Pflichten darunter
zählen nicht mit."* Wer die Vorlage benutzt **und** die Größen-Regel prüft,
bekommt Falsch-Positive auf jedem neuen Slice, während der Altbestand grün
bleibt — der Sensor wird über die Zeit unbrauchbar, nicht sofort.

**Der zweite Bedarf ist hier ebenso belegt.** Grandfathering läuft in diesem
Repo über **Datei**-Globs ([`MR-049`](../../../harness/conventions.md#mr-049)),
weil es nichts Feineres gibt. Eine Regel mit `sections: each` über **eine**
Datei — etwa das Lastenheft — hätte keinen Hebel; die Alternative wäre, den
Bestand umzuschreiben, und das träfe die Form statt der Substanz.

**Der Antrag hat zwei gemessene Defekte, und sie prägen diese Entscheidung
mehr als der Antrag selbst.** Gefahren über dieselben 444 Items:

| Musterform | Items ignoriert | davon falsch |
|---|---|---|
| frei, wie im CR | 13 | **2** |
| verankert (das Item **beginnt** mit dem Ausdruck) | 26 | 0 |

Erstens nimmt das freie Muster **echte Zusagen** mit: zwei der 13 sind
Liefer-Punkte, die die Closure-Notiz nur als ihren **Ort** nennen. Zweitens —
und das ist der schärfere — trifft die erste Alternative des CR, `make gates`,
**null** von 444 Items. Der Grund ist die Zusage, die der CR selbst verlangt:
gezählt wird der **bereinigte** Text, und in allen 86 Abschnitten steht die
Wendung in Inline-Code (null Fundstellen ohne Backticks). Was fence- und
inline-code-treu gezählt wird, ist auch so **ignorierbar** — wer ein Muster auf
einen Ausdruck in Backticks richtet, richtet es auf Leerzeichen.

## Entscheidung

**1. Zwei Optionen an bestehenden Bedingungen, kein neues Modul und kein neuer
Grund-Code.** `tasks-ignore-pattern` gehört zu `max-tasks` wie
`table.order-column` zu `table.order`; `exempt-section-pattern` ist das
Geschwister von `exempt-paths` eine Granularitätsstufe tiefer. Beide
**verkleinern** nur die geprüfte Menge und beantworten dieselbe Frage wie
zuvor. Die Verortung folgt [ADR-0070](0070-tabellen-klammer-und-spaltenliste.md).

**2. Der Schlüssel heißt `exempt-section-pattern`, nicht `exempt-sections`.**
Zwei Gründe, beide aus dem Bestand: In `structure` bedeutet das Suffix
`-pattern` **RE2** (`section-pattern`, `forbid-pattern`, `require-pattern`) —
und `tasks-ignore-pattern` des Antrags folgt dieser Konvention selbst. Und
`exclude-sections` ist in **zwei** anderen Modulen (`vcs`, `sources`) bereits
vergeben, dort als **Liste literaler** Überschriften. Ein dritter, sprechender
Name kostet vier Zeichen und nimmt beide Verwechslungen heraus.

**3. Die beiden Muster sehen bewusst VERSCHIEDENE Zeichenketten, und das ist
kein Versehen.**

- `exempt-section-pattern` sieht **dieselbe** Zeichenkette wie
  `section-pattern`: die getrimmte Überschriften-Zeile **einschließlich** der
  `#`-Folge. Zwei RE2 in **einer** Regel mit zwei verschiedenen Zielen wären
  die Falle, in die der Antrag selbst gelaufen ist — sein Beispiel
  (`'^AC-…'`) trifft am realen Lastenheft **nichts**, weil dort `### AC-…`
  steht. Ein Muster, das analog zum Nachbarn geschrieben wird und still nicht
  greift, ist genau die Klasse, gegen die dieser Slice gebaut wurde.
- `tasks-ignore-pattern` sieht den **Item-Text** hinter Listen-Marker und
  Checkbox. Gegen die rohe Zeile bezeichnete `^` immer den Listen-Marker, und
  die **verankerte** Form wäre unschreibbar — die aber ist gemessen die
  tragfähige (26 gegen 13 Treffer, 0 gegen 2 falsche). Hier gibt es, anders als
  bei der Überschrift, **keinen** Nachbar-Schlüssel, an dem sich ein Autor
  orientieren könnte; die Asymmetrie kostet also keine Analogie, sie
  ermöglicht die Verankerung.

**4. Die Überdeckung ist sichtbar: die Meldung nennt die Zahl der ignorierten
Items** — *„Abschnitt trägt 4 Task-Items (3 ignoriert), erlaubt sind 3"*. Sie
steht **nur**, wenn ein Muster gesetzt ist (sonst wäre die Meldung jeder
Bestandsregel eine andere), und sie steht dann **auch bei null**: ein Muster,
das nichts trifft, ist eine Zusage, die nicht wirkt, und gehört sichtbar.

**5. Leert `exempt-section-pattern` die Abschnitts-Menge, meldet die Regel
`section-missing`** — dieselbe Nullmengen-Härte wie bei `exempt-paths`, eine
Stufe tiefer, und die Meldung nennt den Schlüssel, der es tat. Ohne diese
Antwort schaltete ein zu breites Muster die Regel **still** ab.

**6. Die Ausnahme läuft VOR der Kardinalitäts-Prüfung.** Was ausgenommen ist,
kann `sections: one` nicht mehrdeutig machen: die Ausnahme erklärt die
**Grundmenge**, und über eine Menge, die es nicht gibt, sagt eine Zahl nichts.
Dieselbe Reihenfolge wie bei `exempt-paths`.

**7. `tasks-ignore-pattern` ohne `max-tasks` ist Exit 2.** Ein Ignorier-Muster
ohne die Zählung, die es verkleinern soll, ist eine Zusage, die dasteht und
nicht wirkt — dieselbe halbe Aktivierung, die `table.order-column` ohne
`table.order` und `headings-level` ohne `headings-match` bereits abweist.

**8. Beide Schlüssel sind einfache Zeichenketten, kein Zeiger.** Der leere Wert
ist die Abwesenheit, wie bei jedem anderen RE2-Schlüssel dieses Moduls — und
beide fallen dabei in die **strenge** Richtung: ohne Muster zählt jedes Item
und wird jeder Abschnitt geprüft. `hint` braucht seinen Zeiger, weil er in die
andere Richtung fällt.

## Verglichene Alternativen

**Zwei eigene Module.** Verworfen: beides ist dieselbe Frage mit erklärter
Grundmenge, keine neue Frage. Ein Modul, das nur eine Teilmenge deklariert,
hätte keinen eigenen Befund und keinen eigenen Grund-Code — es wäre eine
Konfigurations-Umleitung mit Modul-Kostüm.

**`exempt-sections` wie beantragt.** Verworfen aus Entscheidung 2. Der
Wortstamm ist in `vcs` und `sources` mit **anderer Form** (Liste statt RE2)
belegt; wer eine Liste schreibt, bekommt zwar einen lauten Decoder-Fehler,
aber der sagt „cannot unmarshal", nicht „hier steht ein Muster".

**`exempt-headings`.** Verworfen: `headings-match`/`headings-level` meinen im
selben Modul die Überschriften **innerhalb** des Abschnitts. Derselbe
Wortstamm für zwei Ebenen wäre die dritte Verwechslung statt keiner.

**Beide Muster gegen die rohe Zeile.** Verworfen aus Entscheidung 3: die
Verankerung ist die gemessen tragfähige Form, und gegen die rohe Item-Zeile
ist sie nicht schreibbar.

**Grandfathering als Werkzeug-Begriff („ab Nummer N").** Vom Absender selbst
abgegrenzt und hier ebenso: das ist über RE2 ausdrückbar, und ein
Werkzeug-Begriff für einen Stichtag wäre ein Datums-Modell, das das Werkzeug
nicht hat.

**Den eigenen Bestand umschreiben statt der Fähigkeit.** Verworfen: 444 Items
in 86 Dateien anzufassen träfe die Form statt der Substanz, und die
Slice-Pläne in `done/` sind Lauf-Belege.

## Konsequenzen

**Positiv.** Eine Größen-Regel wird auf einer Vorlage mit konstanten
DoD-Punkten überhaupt erst brauchbar. Ein Stichtag steht in der Konfiguration
statt in einem Skript. Der Default ist byte-identisch — ein Test hält die
Meldung ohne Muster als Literal.

**Negativ, und das ist die ehrliche Seite.** **Die geprüfte Menge zu
verkleinern ist per Konstruktion eine Lockerung.** Sie ist opt-in, aber jede
Konfiguration, die sie setzt, senkt ihre eigene Zusage, und **kein Gate meldet
das** — es ist eine Konfigurations-Entscheidung, kein Zustand.

**Die Sichtbarkeits-Zusage hat eine benannte Grenze:** sie greift, **solange
die Regel meldet**. Ein Muster, das so breit ist, dass die Schwelle nie
überschritten wird, bleibt stumm. Wer alles ignoriert, sieht nichts.

**Ein Muster, das auf Inline-Code zielt, trifft Leerzeichen.** Das ist die
Kehrseite der Fence-Treue und wird beim ersten Versuch jeden treffen, der die
Wendung in Backticks schreibt — dieses Repo tut das in **allen** 86
Abschnitten. Die Meldung *„(0 ignoriert)"* ist die Diagnose dafür; das
Handbuch führt den Fall ausdrücklich.

**Kein Anwenden auf den eigenen Bestand.** Die 80 Befunde sind der Anlass,
nicht der Auftrag: ob dieses Repo `max-tasks` über seine Slice-Pläne scharf
schaltet, ist ein eigener Entscheid nach dieser Fähigkeit.

## Fitness Function (falls maschinell prüfbar)

`make test` — die neuen Tests in
`internal/hexagon/core/rules/structure_teilmenge_test.go` und der erweiterte
Config-Rand in `internal/adapter/driven/configyaml/configyaml_test.go`. Sie
halten jede Zusage dieser ADR einzeln, und die tragenden drei tragen ihre
**Umkehr** in derselben Funktion — der Vorzustand ohne den Schlüssel ist
mitgeprüft, ein Regressions-Test ohne belegte Regression ist keiner.

**Nicht maschinell geprüft** ist die Konsequenz, die zählt: ob eine gesetzte
Teilmenge **berechtigt** ist. Das ist ein Urteil über die Konfiguration eines
fremden Repos und kein Gate.

## Re-Evaluierungs-Trigger

**Wenn eine dritte Bedingung eine Teilmenge braucht.** Zwei Optionen an zwei
Bedingungen sind eine Erweiterung; ab der dritten ist die Frage, ob die
Teilmengen-Erklärung eine eigene Klammer verdient — dieselbe Frage, die
[ADR-0070](0070-tabellen-klammer-und-spaltenliste.md) für die Tabelle
beantwortet hat.

**Wenn ein Muster gemessen zu breit war.** Tritt der Fall ein, dass eine
Teilmengen-Erklärung eine echte Zusage still entfernt hat, ist die
Sichtbarkeits-Zusage aus Entscheidung 4 zu schwach und die Frage lautet, ob
eine Obergrenze für den Ignorier-Anteil dazugehört.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-08-30 | **Die Zahl „verankert 26 Treffer, 0 falsche" ist falsch gemessen — die Entscheidung steht, ihr Beleg ändert sich.** §Kontext und §Entscheidung 3 stellen sie neben „frei 13" unter die Überschrift *„über dieselben 444 Items"*; tatsächlich messen die beiden Zeilen **verschiedene Ausdrücke**. Ein `^` vor demselben Ausdruck kann die Treffermenge nur verkleinern, 13 → 26 ist also nicht bloß ungenau, sondern unmöglich. Nachgemessen: das freie CR-Muster trifft **13** (davon **2** echte Liefer-Zusagen), **dasselbe** Muster verankert trifft **1** — und dieses eine ist ebenfalls eine Liefer-Zusage. Die **26** stammen von `^grün \(Exit explizit\)`, einem korpus-eigenen Ausdruck, und dort sind **0** falsch. **Die Empfehlung „verankern" wird dadurch nicht schwächer, sondern präziser:** verankern allein genügt nicht — das Muster muss zusätzlich auf Text zielen, den die **Bereinigung übrig lässt**. Genau daran scheitert das CR-Beispiel doppelt. Gefunden von Review und Verifikation unabhängig; die Klasse ist [`BEO-020`](../planning/observations.md) |
| 2026-08-30 | **Die Korpus-Zahl „86" bezeichnet die falsche Menge.** §Kontext, §Konsequenzen und die Fitness Function sprechen von den *„**86** Slice-Plänen dieses Repos, die einen `## 4. Definition of Done` führen"*. Gemessen sind es **89** Abschnitte in 175 Slice-Dateien; **86** ist die Zahl der Dateien **ohne** den Abschnitt (der `section-missing`-Zähler desselben Laufs) — und, zufällig, die Zahl der Abschnitte, die `make gates` tragen. Die abgeleiteten Zahlen (444 Items, 80 `section-oversized`, 166 im Byte-Identitäts-Lauf) sind davon **unberührt** und reproduzieren. Ebenfalls zu eng: *„in **allen** 86 Abschnitten steht die Wendung in Inline-Code"* — es sind **86 von 89**; drei Abschnitte tragen sie gar nicht. Die tragende Aussage (null von 444 Items, weil durchgängig in Backticks) bleibt |
| 2026-08-30 | **Der „dritte gemessene Punkt" gegen den Antrag ist keiner.** §Entscheidung 3 schreibt, das Beispiel des Absenders sei *„die Falle, in die der Antrag selbst gelaufen ist"*. Der CR deklariert seinen Vergleichsgegenstand aber ausdrücklich als **Überschriftstext**, und sein Muster ist gegen **diese** Semantik korrekt. Dass es hier nicht greift, folgt aus **unserer** Entscheidung, gegen die rohe Zeile zu vergleichen — das ist eine Entscheidung, kein Defekt des Antrags, und es wurde am Antrag nichts gemessen. Die Entscheidung selbst bleibt richtig; ihre Begründung steht auf der Konsistenz innerhalb der Regel, nicht auf einem fremden Fehler. Klasse: [`BEO-012`](../planning/observations.md) |
| 2026-08-30 | **Zwei Lücken der Umsetzung, beide geschlossen, beide ohne Änderung an einer Entscheidung dieser ADR.** Erstens löschte ein `hint` (ADR-0073) **beide** Sichtbarkeits-Zusagen aus §Entscheidung 4/5. Der Nullmengen-Befund gehört zur Klasse *„die Regel hat nicht gemessen"*, die ADR-0073 vom verfassten Hinweis ausdrücklich **ausnimmt** — er trägt jetzt die modul-eigene Meldung. Für die `section-oversized`-Zahl gewinnt der Hinweis dagegen zu Recht; das ist eine **zweite benannte Grenze** der Sichtbarkeits-Zusage und steht seither in Lastenheft, Spezifikation und Handbuch. Zweitens ging das **gesetzte** `exempt-section-pattern` nicht in die Regel-Identität ein, womit die Paarung *Bedingung A für alle Abschnitte, B nur für die übrigen* — die Form, die Grandfathering braucht — als Konfigurations-Duplikat mit Exit 2 endete. Es geht jetzt ein, wie `table.order-column` seit [ADR-0069](0069-zellenlaenge-als-strukturbedingung.md); ein leeres Muster nicht, damit kein Bestands-`target` wandert |
