# ADR-0074: Offene Task-Items werden auf den rohen Abschnitts-Zeilen gezählt

**Status:** Accepted

**Datum:** 2026-08-29

**Autor:** pt9912

**Bezug:** [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(das erweiterte Modul);
[ADR-0059](0059-closure-waechter-weicht-structure-regel.md) (der ausgewiesene
Preis der absatzweisen Paarung — hier wird er zu teuer);
[ADR-0057](0057-structure-tabellen-monotonie.md) (die Präzedenz: eine Bedingung,
die **rohe** Zeilen liest, mit Begründung);
[ADR-0069](0069-zellenlaenge-als-strukturbedingung.md) (dieselbe Bauform: neue
Bedingung, eigener Grund-Code, Befund auf **ihrer** Zeile);
[ADR-0073](0073-befund-erlaeuterung-fuer-menschen.md) (der `hint`, den eine
solche Regel trägt)

**Schärft:** [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)

**Regeln:** Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v5.12.0/regelwerk/modul-04-adrs.md).

---

## Kontext

**Eine Zusage über offene Task-Items ist auf dem bereinigten Abschnitts-Text
nicht haltbar.** `forbid-pattern` und `max-tasks` lesen den Text, aus dem
Fenced-Code entfernt und **Inline-Code geleert** wurde. Die Paarung ist
**absatzweise** — der in
[ADR-0059](0059-closure-waechter-weicht-structure-regel.md) ausgewiesene Preis.
Ein **einzelner** überzähliger Backtick irgendwo im Absatz macht damit alles
dahinter unsichtbar.

**Gemessen an einem Wächter über offene DoD-Haken:**

| Eingabe | Ergebnis |
|---|---|
| offener Haken | 1 Befund, Exit 1 |
| derselbe Haken, ein Backtick weiter oben im Absatz | **0 Befunde, Exit 0** |

Nicht teilweise blind, sondern ganz. Und die Exposition war real: zwei
`done/`-Slices des eigenen Bestands tragen ungerade Backtick-Zahlen in ihrem
DoD-Abschnitt.

**Der Preis war für einen Vorlagen-Platzhalter ausgewiesen; hier zahlt ihn eine
Vorbedingung.** ADR-0059 hat die Paarung gegen eine Falsch-Positiv-Klasse
entschieden — ein Slice, der **über** Platzhalter schreibt, sollte nicht seine
eigene Dokumentation melden. Diese Abwägung bleibt für Prosa richtig. Für eine
Bedingung, an der ein **Übergang** hängt, kippt sie: ein Wächter, den ein
Tippfehler abschaltet, meldet grün, wo er nichts gesehen hat.

**Ein zweiter Befund derselben Klasse:** ein `forbid-pattern` aus der
Konfiguration deckt nur die Bullet-Form, die sein Autor aufgeschrieben hat.
Gemessen liefen `* [ ]` und `+ [ ]` durch ein `- \[ \]` still hindurch. Das ist
die Gestalt, die dieses Repo als
[`BEO-003`](../planning/observations.md) führt — eine geteilte Lexik driftet,
weil jeder Konsument sie selbst vorbereitet.

**Die Hälfte der Fähigkeit lag bereits im Modul.** `taskItemRE` erkennt ein
Task-Item nach der vollen CommonMark-Form: Bindestrich, Stern, Plus **und**
geordnete Liste, mit führendem Weißraum und Leerzeichen oder Tab als Trenner.
Was fehlte, war eine Bedingung, die sie auf die rohen Zeilen anwendet.

## Entscheidung

**Eine neue Bedingung `max-open-tasks` (int ≥ 0)** zählt die **offenen**
Task-Items eines Abschnitts auf den **rohen** Zeilen und meldet
`section-tasks-open`.

1. **Roh statt bereinigt — und nur hier.** Die Bedingung ist die **dritte**, die
   einen anderen Text liest, neben der Chronologie-Monotonie
   ([ADR-0057](0057-structure-tabellen-monotonie.md)) und der
   Überschriften-Bedingung. Der Kommentar über `structureConditions` nennt sie
   alle drei, statt die neue stillschweigend danebenzustellen.
2. **Kein generischer „roh"-Schalter.** Ein `raw: true`, das jedes Muster an der
   Bereinigung vorbeiführte, holte genau die Falsch-Positiv-Klasse zurück, gegen
   die ADR-0059 entschieden hat. Die Antwort ist eine **benannte** Bedingung mit
   benanntem Preis, keine Generalvollmacht.
3. **Die Lexik kommt aus dem Modul, nicht aus der Konfiguration.**
   `openTaskItemRE` teilt die Form von `taskItemRE` und verengt sie auf die
   **leere** Box. Damit kann die Bedingung den Bullet-Form-Fehler nicht machen,
   den ein Konfigurations-Muster gemacht hat.
4. **Ein Befund je offenem Item, auf seiner Zeile.** Die Reparatur ist dort, wo
   der Haken steht — dieselbe Wahl wie bei der Zellenlängen-Bedingung
   ([ADR-0069](0069-zellenlaenge-als-strukturbedingung.md)), und der Unterschied
   zur Chronologie-Bedingung, die auf die Abschnitts-Überschrift zeigt, wenn sie
   leer läuft.
5. **Der Fence bleibt außen vor, die Inline-Code-Spanne nicht.** Die
   fence-bewusste Zeilen-Auswahl ist dieselbe, die die Tabellen-Bedingungen
   benutzen. Der Fence deckt den Fall, der die Falsch-Positiv-Sorge trägt: ein
   Dokument, das **über** Task-Items schreibt, illustriert sie im Fence. Roh
   heißt ansonsten roh — und das ist der ausgewiesene Preis dieser Bedingung.
6. **`max-tasks` bleibt unverändert.** Es zählt **alle** Items auf dem
   bereinigten Text und teilt dessen Blindstelle. Es umzustellen wäre eine
   Verhaltensänderung an einem ausgelieferten Schlüssel; sie ist nicht gemessen
   und deshalb nicht Teil dieses Entscheids.

## Verglichene Alternativen

| Option | Pro | Contra | Verworfen, weil |
|---|---|---|---|
| **A: `raw: true` je Regel** | eine Zeile, wirkt für jedes Muster | jede Bedingung könnte an der Bereinigung vorbei; die Falsch-Positiv-Klasse aus ADR-0059 kehrt zurück | eine Generalvollmacht ist keine Antwort auf einen benannten Fall |
| **B: die Bereinigung reparieren** (Paarung nicht absatzweise) | löste es für **alle** Bedingungen | ADR-0059 hat die Paarung gegen eine gemessene Falsch-Positiv-Klasse entschieden; der Preis war bekannt und akzeptiert | dieselbe Abwägung noch einmal zu treffen, ohne neue Messung, wäre ein Rückbau auf Verdacht |
| **C: `forbid-pattern` um alle Bullet-Formen erweitern** | keine neue Fläche | löst nur den **zweiten** Befund; der Backtick-Fall bleibt, und die Lexik lebt weiter in der Konfiguration | die kleinere Hälfte des Problems, mit der Bauform, die ihn erzeugt hat |
| **D: `max-tasks` auf roh umstellen** | kein neuer Schlüssel | Verhaltensänderung an einem ausgelieferten Schlüssel, ungemessen; und es zählt **alle** Items, nicht die offenen | ein stiller Semantik-Wechsel ist teurer als ein benannter neuer Schlüssel |

## Konsequenzen

**Positiv.** Die Zusage über offene Task-Items ist haltbar: der Backtick-Fall
meldet, gemessen. Alle vier Listen-Marker zählen, ohne dass ein
Konfigurations-Autor sie kennen muss. Der Befund zeigt auf die Zeile, an der die
Reparatur stattfindet, und es gibt einen je Item statt eines je Datei.

**Negativ, und benannt.** Roh heißt roh: ein Task-Item in einer
**Inline-Code**-Spanne zählt jetzt mit. Der Fence deckt den häufigen Fall — ein
Dokument, das über Task-Items schreibt, illustriert im Fence —, die Inline-Form
nicht. Wer sie braucht, schreibt sie in einen Fence.

**Zwei Bedingungen über verwandte Fragen stehen nebeneinander.** `max-tasks`
(bereinigt, alle Items) und `max-open-tasks` (roh, offene Items). Wer den
falschen greift, bekommt stillschweigend die schwächere Zusage. Dagegen steht
die Abgrenzung in beiden Schema-Zeilen, nicht ein Sensor.

**Der Grund-Code-Raum wächst um einen weiteren `section-*`-Code.** Ein
Sammel-Code schiede aus: die Befund-Deduplikation vergleicht (Datei, Zeile,
Regel, Ziel, Grund), und zwei verletzte Bedingungen desselben Abschnitts fielen
darunter zusammen.

## Fitness Function (falls maschinell prüfbar)

**Ja, und die Proben sind gefahren.** Sieben Kern-Tests decken den
Backtick-Fall, alle vier Listen-Marker (eingerückt und mit Tab-Trenner), den
gesetzten Haken, die Fence-Treue, einen Befund je Item auf seiner Zeile, die
Schwelle als Obergrenze und den abwesenden Schlüssel als „Bedingung aus". Der
Config-Rand (`max-open-tasks` < 0 ⇒ Exit 2) und die Durchreichung der expliziten
Null liegen im Decode-Test.

**Drei Mutationsproben belegen, dass die Tests beißen** — die Prozedur aus
[`BEO-023`](../planning/observations.md): Fence-Gate entfernt ⇒ der
Fence-Test wird rot; Lexik auf den Bindestrich verengt ⇒ der Marker-Test wird
rot; ein Befund je Datei statt je Item ⇒ zwei Tests werden rot.

## Re-Evaluierungs-Trigger

**Wenn die Inline-Code-Grenze im eigenen Bestand zuschlägt.** Meldet die
Bedingung ein Task-Item, das nur eine Illustration in Backticks war, ist zu
entscheiden, ob die Inline-Form mit ausgenommen wird — und damit ein Stück der
Blindstelle zurückkehrt.

**Wenn eine zweite Bedingung die rohen Zeilen braucht.** Bei drei benannten
Ausnahmen von der Bereinigung ist die Frage neu zu stellen, ob die Bereinigung
selbst der falsche Default ist — dann wäre Option B zu messen statt zu
verwerfen.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-08-29 | **Die Umsetzung ist zurückgenommen, die Entscheidung steht.** Review und Verifikation fanden, dass die Bedingung einen **stillen** Ausfall derselben Gestalt behält, gegen die sie gebaut ist: ein einzelner vergessener Schluss-Fence schaltet sie ab — dasselbe Muster wie der Backtick, nur mit anderem Zeichen. Verschärfend fährt das Closure-Profil, für das sie gedacht ist, `spans` nicht, wo `fence-unclosed` den häufigen Fall sonst laut meldete; ein naiv ausgeglichener Fence entkommt auch dem. Entscheid des Auftraggebers: **erst die Fence-Lexik**, dann die Bedingung. Der Code, die Spec-Änderungen und die Handbuch-Zeilen sind entfernt; diese ADR bleibt, weil eine `Accepted`-Entscheidung ein Beleg ist und nicht mit ihrer Umsetzung verschwindet (`AGENTS.md` §3.5, maschinell gehalten) |
| 2026-08-29 | **Zwei Zahlen im Text sind falsch, und die Korrektur gehört hierher statt in den Kern.** §Entscheidung 1 und §Kontext sprechen von der „**dritten**" Bedingung, die einen anderen Text liest, und zählen Chronologie und Überschriften daneben. Es sind **vier**: die **Zellenlängen**-Bedingung liest ebenfalls die rohen Zeilen und sagt das selbst. Folge für §Re-Evaluierungs-Trigger: *„bei **drei** benannten Ausnahmen"* war im Moment der Annahme bereits erreicht — der Trigger konnte nie feuern. Er ist als **erfüllt** zu lesen: die Frage, ob die Bereinigung der falsche Default ist, steht damit offen und wird im Zuge der Fence-Arbeit beantwortet |
| 2026-08-29 | **Die Fitness Function beschreibt Tests, die es nicht mehr gibt** — sie sind mit der Umsetzung entfernt. Zwei von ihnen taugten ohnehin nicht: das Backtick-Fixture trug eine **wohlgeformte** Spanne (zwei Backticks) und wurde vom Vorgänger gefunden, belegte also nicht die Blindstelle; und die Abschnittsgrenze hielt kein Test (die Mutation „über die ganze Datei zählen" lief grün durch). Beim Wiederaufbau sind sie neu zu schreiben, nicht zu übernehmen |
| 2026-08-30 | **Die Umsetzung folgt Entscheidung 3 in ihrer Absicht und nicht in ihrer benannten Mechanik — und das ist die stärkere Form.** Dort steht, `openTaskItemRE` *„teilt die Form von `taskItemRE` und verengt sie auf die leere Box"*. Genau daran scheiterte der erste Bau: das zweite RE2 war ein **wörtliches Präfix** des ersten, ohne Kopplungs-Test, also exakt die [`BEO-003`](../planning/observations.md)-Form, gegen die die Entscheidung argumentiert. Es gibt jetzt **kein zweites Muster**: die Zeile wird mit `taskItemRE` erkannt, und die Box liest sich aus dessen Treffer (`strings.HasSuffix(…, "[ ]")`). Damit kann die Lexik nicht driften, statt dass ein Test ihre Kopplung überwacht. Die Absicht — *die Lexik kommt aus dem Modul, nicht aus der Konfiguration* — ist unverändert eingelöst |
| 2026-08-30 | **Die Schwellen-Semantik war offen und ist entschieden:** die ersten `max-open-tasks` offenen Items in **Dokument-Reihenfolge** sind erlaubt und melden nicht, gemeldet wird der **Überhang**. Der erste Bau meldete bei einer Schwelle > 0 **alle** Items, auch die erlaubten — eine Verletzung ergab drei Befunde, und keiner davon war die Reparaturstelle. Für den gedachten Anwendungsfall (`max-open-tasks: 0`) sind beide Lesarten gleich; der Unterschied entsteht erst über null, und dort ist die erlaubte Menge kein Defekt |
| 2026-08-30 | **Die ausgewiesene Inline-Code-Grenze war überzeichnet und ist berichtigt.** §Konsequenzen sagt *„ein Task-Item in einer **Inline-Code**-Spanne zählt jetzt mit"*. Gemessen gilt das nur für die **mehrzeilige** Spanne: das Muster ist zeilen-verankert, und bei einer einzeiligen steht der Backtick **vor** dem Listen-Marker — sie meldet gar nicht. Der genannte Preis war höher als der echte. Zwei Grenzen kamen dazu, die keine Fassung nannte: ein Task-Item im **Blockquote** (`> - [ ]`) und eine Box mit **Tabulator** (`- [\t]`) zählen für **keine** der beiden Bedingungen — Eigenschaft der geteilten Lexik, nicht dieser Bedingung, und deshalb hier nur benannt statt geändert |
| 2026-08-30 | **Der Fence-Pfad, der zur Rücknahme führte, ist gemessen zerlegt — und nur eine Hälfte ist die Lexik-Frage.** Ein **vergessener** Schluss-Fence blendet alles Folgende aus; das tut die CommonMark-Lesart **genauso** (eine offene Fence läuft bis Dateiende), ein Lexik-Wechsel behebt ihn also nicht. Ihn fängt `fence-unclosed`, und das Closure-Profil fährt `spans` seit [ADR-0077](0077-spans-am-bindepunkt-die-begruendung-traegt-anders.md). Nur der **naiv ausgeglichene** Fence ist die Divergenz zwischen dem Toggle und der CommonMark-Lesart — und die hat **null** Realfälle: beide Automaten über alle 620 Markdown-Dateien gegeneinander gefahren, einziger Treffer ist die Prosa von [ADR-0042](0042-markdown-lexik-folgt-commonmark.md), die den Unterschied beschreibt. Dieselbe ADR hat die Frage ausdrücklich offen gelassen und ihre Bedingung benannt: *„erst eine Regel, wenn ein Realfall existiert"*. Es existiert keiner, also bleibt sie offen und diese Bedingung nennt die Grenze |
| 2026-08-30 | **Die Fitness Function ist gegen acht Mutationen gemessen, nicht drei.** Fence-Gate entfernt ⇒ 1 rot · Box nicht auf leer verengt ⇒ 3 rot · Schwelle ignoriert ⇒ 1 rot · über die ganze Datei statt im Abschnitt gezählt ⇒ 1 rot · ein Befund je Datei statt je Item ⇒ 2 rot · bereinigt statt roh gelesen ⇒ 2 rot · negativer Wert geschluckt ⇒ 1 rot · **Verdrahtung Config → Modell entfernt ⇒ 1 rot**. Die letzte fehlte dem ersten Bau und ist die Lücke, an der ein Schlüssel im Produkt wirkungslos wird, während alle Regel-Tests grün bleiben — sie bauen ihre Regel direkt und erreichen den Adapter nie |
