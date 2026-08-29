# ADR-0057 — Chronologie-Monotonie als siebte `structure`-Bedingung: typisiert, roh gelesen, geteilte Tabellen-Lexik

**Status:** Accepted (teil-superseded: ADR-0070)
**Datum:** 2026-08-21
**Autor:** pt9912
**Schärft:** [`spec/spezifikation.md` §DC-FA-STRUCT-001.a](../../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
**Bezug:** [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(Erweiterung); Modul-Schnitt [ADR-0049](0049-structure-modul-schnitt-und-preset.md)
(**nicht** superseded — eine additive Bedingung, kein neuer Schnitt);
Schnitt-Kriterium [ADR-0044](0044-geteiltes-referenz-ventil-quell-skopus.md);
Lexik-Bindung [ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md);
gescopte Roh-Lesungen als Präzedenz
[ADR-0019](0019-versions-pin-fence-ausnahme.md),
[ADR-0020](0020-content-pin-fence-ausnahme.md).

## Kontext

Eine chronologische Tabelle kippt still ihre Richtung: wer eine Zeile anfügt,
schaut auf die Zeile daneben statt auf die Regel — aus „unten anhängen" wird
irgendwann „oben einfügen", und danach führt dieselbe Tabelle zwei gegenläufige
Blöcke. Sichtbar wird das erst, wenn jemand nach einem Datum sucht: kein Gate
liest Reihenfolge, und beide Fassungen sind für sich plausibel. Im eigenen
Bestand ist die Klasse am selben Tag an **drei** Tabellen gefunden worden
(die §7-Historien beider Spec-Straten und die Drift-Tabelle der Roadmap); der
Bruch war jedes Mal die Stelle, an der die Pflege-Gewohnheit wechselte.

Gemessen statt behauptet (sechs Bestandstabellen, 318 Zeilen zum
Messzeitpunkt): ein **typisierter** Monotonie-Vergleich der Schlüsselspalte
findet am Stand **vor** der Heilung alle drei gekippten Tabellen
(14 · 6 · 7 Verletzungen) und am Stand danach null. Ein **naiver**
String-Vergleich dagegen meldet drei **korrekt** sortierte Tabellen rot
(`0.10.0 → 0.9.0`, `v0.10.0 → v0.9.0`, `1.9 → 1.10`) — der Typ ist Pflicht,
kein Komfort.

Drei Eigenschaften der Eingabe sind Entwurfs-Zwänge, keine Details: die
Schlüsselzelle von `version.md` steht in **Inline-Code** (die aktuelle Zeile
zusätzlich mit wanderndem HTML-Anker), die Bedingung braucht als einzige eine
**Zell-Adresse** (Tabellenzeile, Kopf-/Trennzeile, Spalte), und das
Benutzerhandbuch ist als einzige Bestandstabelle **aufsteigend** sortiert.

## Entscheidung

1. **Siebte Bedingung in [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) — kein neues Modul, kein neues
   Kürzel.** Das Kriterium ist
   [ADR-0044](0044-geteiltes-referenz-ventil-quell-skopus.md) entnommen:
   querschnittlich ⇒ eigenes Kürzel, Einzelmodul-Frage ⇒ bestehende Anforderung
   ändern. Die Chronologie-Frage hat dieselbe Eingabe-Achse wie jede
   `structure`-Regel (eine Regel benennt Dateien und Abschnitt selbst); sie ist
   eine weitere Aussage über die Form **innerhalb** eines Dokuments — exakt die
   Modul-Grenze aus [ADR-0049](0049-structure-modul-schnitt-und-preset.md)
   Entscheidung 2.

2. **Zwei Schlüssel je Regel: `table-order` und `table-column`.**
   `table-order` (`asc` | `desc`, kein Default) schaltet die Bedingung scharf —
   die Richtung ist **Regel**-Konfiguration, weil der eigene Bestand beide
   Richtungen legitim führt (das Handbuch aufsteigend, alles andere
   absteigend); ein Modul-Default wäre eine stille Behauptung über fremde
   Repos. `table-column` (1-basiert, Default 1) adressiert die Schlüsselspalte;
   gesetzt ohne `table-order` ⇒ Exit 2 (eine halbe Aktivierung ist ein
   Config-Fehler, kein Zustand), explizit < 1 ⇒ Exit 2, `table-order` außerhalb
   `asc`/`desc` ⇒ Exit 2.

3. **Typisierter Vergleich über eine geschlossene Typ-Menge.** Zwei Typen:
   **ISO-Datum** (`JJJJ-MM-TT`) und **Punkt-Version** (optionales `v`-Präfix,
   mindestens zwei numerische Segmente, segmentweise numerisch verglichen;
   kürzere Segmentfolge bei gleichem Präfix ist kleiner). Getypt wird der
   **erste** Treffer beider Muster in der **rohen** Zelle — so überleben
   Inline-Code-Backticks und der HTML-Anker der `version.md`, ohne dass die
   Bedingung Markup interpretiert. Keine konfigurierbare Typ-Registry und kein
   frei konfigurierbares Schlüssel-Muster: das wäre ein zweiter
   Regel-Interpreter im Werkzeug — dieselbe Grenze, die
   [ADR-0049](0049-structure-modul-schnitt-und-preset.md) Entscheidung 6 für
   Stichtags-Mechaniken gezogen hat. Weitere Typen sind ein Change Request
   gegen diese Entscheidung, kein stilles Wachstum.

4. **Rohe Abschnitts-Zeilen — die erste benannte Ausnahme von „alle
   Bedingungen arbeiten auf diesem Text".** Die Bereinigung aus
   [§`DC-FA-STRUCT-001.a`](../../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
   Schritt 5 leert Inline-Code; damit wäre die
   `version.md`-Schlüsselspalte leer und jede Zelle untypisierbar. Die
   Bedingung liest deshalb die rohen Zeilen des Abschnitts. Das ist die
   Trennlinie aus [ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md)
   Entscheidung 2 — **andere Frage, nicht andere Antwort**: „was steht in
   dieser Zelle?" ist keine Lexik-Frage, die anderswo anders beantwortet würde
   (Präzedenz: die gescopten Roh-Lesungen
   [ADR-0019](0019-versions-pin-fence-ausnahme.md)/[ADR-0020](0020-content-pin-fence-ausnahme.md)).
   Die **Lexik**-Fragen bleiben geteilt: ob eine Zeile eine Tabellenzeile ist,
   entscheidet die fence-bewusste Tabellen-Antwort des Produkts, nicht ein
   roher `^|`-Test.

5. **Zell-Adresse über die geteilte Tabellen-Lexik — und die bekommt mit dem
   dritten Konsumenten ihre Kopplung.** Tabellenzeilen-Erkennung ist bereits
   die eine Antwort des Produkts (`targets`, `planning.waves`);
   Trennzeilen-/Kopfzeilen-Erkennung und Zell-Splitting leben bisher privat
   beim zweiten Konsumenten und wandern an den geteilten Ort. Nach
   [ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 4
   genügt das geteilte Prädikat nicht: ein **Kopplungs-Test** fährt dieselbe
   Eingabe durch alle drei Konsumenten und schlägt bei abweichender Antwort
   fehl — eine Aufzählung kann eine Stelle vergessen, eine Kopplung nicht.

6. **Nicht-strikte Monotonie über benachbarte Datenzeilen, je
   zusammenhängender Tabelle.** Gleiche Schlüssel sind erlaubt — mehrere
   Releases an einem Tag sind der Normalfall, eine strikte Ordnung meldete den
   gesunden Bestand rot. Geprüft wird jede **zusammenhängende** Folge von
   Tabellenzeilen im Abschnitt für sich; je Zeile, die die Richtung ihres
   Vorgängers bricht, **ein** Befund (`section-unordered`, `line` = die
   brechende Zeile). Kopf- und Trennzeile deklarieren keine Daten und werden
   übersprungen — dieselbe Lesart wie beim Wellen-Register. **Null Datenzeilen
   im Abschnitt ⇒ `section-unordered` als Leerlauf-Befund** (`line` = die
   Abschnitts-Überschrift; dieselbe Doppel-Rolle wie `section-missing` beim
   Kandidaten-Leerlauf): die Bedingung zu setzen **ist** die Behauptung, dass
   der Abschnitt eine chronologische Tabelle trägt — sonst schaltete eine in
   eine Liste umgebaute Tabelle die Zusage wortlos ab.

7. **Eine untypisierbare Zelle ist ein Befund, kein Übersprung.** Fehlt der
   Spalte eine typisierbare Zelle (kein Muster trifft, die Zeile hat zu wenige
   Zellen) oder mischt die Spalte zwei Typen (Datum neben Version —
   unvergleichbar), meldet `section-cell-untyped` an dieser Zeile. Ein stilles
   Auslassen schaltete die Prüfung der restlichen Tabelle wortlos ab — dieselbe
   fail-closed-Disziplin wie der Leerlauf-Befund des Moduls. Der
   Vergleichs-Anker wird nach **jedem** Befund zurückgesetzt — auch nach
   einer Typ-Mischung — und setzt beim nächsten typisierbaren Nachbar-Paar
   wieder auf: eine kaputte Zelle macht nicht die ganze Tabelle unprüfbar,
   sie meldet sich **selbst**, und die gesunde Folge-Zeile dahinter meldet
   **nicht** (Paar-Lesart, im Review als einzig verbleibende Lesart gepinnt).

8. **Zwei neue Grund-Codes, kein Sammel-Code.** `section-unordered` (Reparatur:
   die Zeile einsortieren) und `section-cell-untyped` (Reparatur: die
   Schlüsselzelle bzw. Spaltenwahl korrigieren) verlangen verschiedene
   Reparaturen; die Befund-Deduplikation vergleicht (Datei, Zeile, Regel, Ziel,
   Grund) — die Begründung aus
   [ADR-0049](0049-structure-modul-schnitt-und-preset.md) Entscheidung 8 gilt
   unverändert.

## Verglichene Alternativen

| Alternative | Warum verworfen |
|---|---|
| Neues Modul „chronology" | Gleiche Eingabe-Achse wie `structure` (Regel benennt Dateien + Abschnitt); das ADR-0044-Kriterium sagt Einzelmodul-Erweiterung, und ein Zwillings-Modul verdoppelte Kandidaten-Menge, Kardinalität und Ventile |
| Naiver String-Vergleich | Gemessen falsch: drei korrekt sortierte Bestandstabellen würden rot (`0.10.0 → 0.9.0`, `v0.10.0 → v0.9.0`, `1.9 → 1.10`) |
| Konfigurierbares Schlüssel-Muster / offene Typ-Registry | Zweiter Regel-Interpreter im Werkzeug — die Grenze aus ADR-0049 Entscheidung 6; die geschlossene Zwei-Typen-Menge deckt alle sechs Bestandstabellen |
| Vergleich auf dem bereinigten Abschnitts-Text | Leert die `version.md`-Schlüsselspalte (Inline-Code) — jede Zelle wäre untypisierbar; die Ausnahme muss in den Vertrag, nicht in den Code |
| Untypisierbare Zellen still überspringen | Ein Tippfehler in einer Zelle schaltete die Prüfung wortlos ab — der stille Grün-Pfad, den das Modul überall sonst fail-closed vermeidet |
| Strikte Monotonie | Mehrere Releases am selben Tag sind der Normalfall des eigenen Bestands — der gesunde Zustand würde rot |
| Richtung als Modul-Default (`desc`) | Der eigene Bestand führt beide Richtungen legitim; ein Default wäre eine stille Behauptung über fremde Repos und machte die Handbuch-Regel zur Ausnahme-Syntax |
| Ein Sammel-Code für beide Fälle | Verschiedene Reparaturen; und zwei Befunde derselben Zeile fielen unter der Deduplikation zusammen (ADR-0049 Entscheidung 8) |
| Tabellen-Lexik lokal nachbauen | Die dritte Kopie derselben Antwort — exakt die Klasse aus ADR-0054, mit Kopplungs-Test statt Aufzählung zu binden |

## Konsequenzen

- **Opt-in und additiv** ⇒ SemVer-Minor: ohne `table-order` ist der Befundsatz
  byte-identisch; kein bestehender Konsument ändert sein Verhalten.
- **Die Roh-Lese-Ausnahme steht im Vertrag**, nicht nur im Code — sie ist die
  erste Bedingung, die Schritt 5 nicht konsumiert, und genau deshalb in
  Anforderung und Algorithmus ausdrücklich benannt. Ein zweiter Anwender dieser
  Ausnahme wäre ein Re-Evaluierungs-Fall, kein Präzedenz-Zitat.
- **Zwei getrennt sortierte, gegenläufige Tabellen im selben Abschnitt bleiben
  unerkannt** (benannte Grenze): geprüft wird je zusammenhängender Tabelle. Der
  Anlassfall war ein Richtungs-Bruch **innerhalb** einer Tabelle; wer eine
  Tabelle in zwei spaltet, hat ein anderes Problem als eine stille Kipp-Stelle.
- **Die Lexik-Hebung berührt ausgeliefertes Verhalten** (`planning.waves` seit
  v0.59.0, `targets` seit v0.44): reine Verschiebung an den geteilten Ort,
  verhaltensgleich — bewacht vom Kopplungs-Test und den Bestandstests beider
  Module.
- Der `--doctor`-Klartext und die §4-Tabelle wachsen um zwei Zeilen
  (AllReasons-↔-§4-Lockstep).
- **Eine Chronologie-Zusage je Abschnitt:** die Regel-Identität besteht aus
  Glob und Abschnitts-Selektor und trägt keine Spalte — zwei Regeln gleicher
  Identität mit verschiedenem `table-column` sind ein Konfigurations-Duplikat
  (Exit 2, laut statt still). Wer zwei Spalten derselben Tabelle monoton
  zusagen will, stellt einen Change Request gegen diese benannte Grenze; sie
  still über eine Spalten-erweiterte Identität zu öffnen änderte die
  Deduplikations-Semantik aller sieben Bedingungen.
- **Eine zweite Zell-Lesart existiert im Produkt** und bleibt getrennt: der
  RTM-/trace-Leser zerlegt Tabellenzeilen escape- und backtick-bewusst
  (andere Vertragsfläche, andere Frage im Sinn von
  [ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 2).
  Die Konsequenz gehört benannt: ein Pipe in einem Backtick-Span der
  Schlüsselspalte verschöbe hier die Spaltenadresse — im heutigen Bestand der
  sechs aktivierten Tabellen kommt das nicht vor; träte es ein, ist das der
  Latenz-Fall der geteilten-Lexik-Klasse und nach deren Regeln zu behandeln.

## Fitness Function

- **Retro-Beleg mit dem Produkt:** am Stand vor der Heilung melden die drei
  gekippten Tabellen Verletzungen (Größenordnung 14 · 6 · 7 der
  Skript-Messung — Abweichungen sind zu erklären, nicht zu glätten), am
  heutigen Stand null.
- **Typ-Pflicht:** die naive Gegenprobe ist ein Testfall — eine korrekt
  absteigende Versions-Spalte mit `0.10.0` über `0.9.0` bleibt grün.
- **Roh-Lesung:** eine Schlüsselspalte in Inline-Code (`version.md`-Form) wird
  getypt; auf dem bereinigten Text wäre sie leer.
- **Kopplung:** dieselbe Eingabe durch `targets`, `planning.waves` und
  `structure` — die Tabellenzeilen-Antwort ist eine.
- **fail-closed:** eine untypisierbare Zelle bzw. Typ-Mischung meldet; ein
  Abschnitt ohne Tabellen-Datenzeile meldet den Leerlauf; die drei Config-Ränder
  brechen mit Exit 2.
- **Byte-Identität:** ohne `table-order` ist der Befundsatz unverändert.

## Re-Evaluierungs-Trigger

- Ein Adopter braucht wiederholt einen **dritten Schlüssel-Typ** (Zeitstempel,
  Kennungs-Nummern) — dann ist Entscheidung 3 gegen die gelebte Praxis zu
  prüfen; die Antwort wäre eine bewusste Erweiterung der geschlossenen Menge,
  kein konfigurierbares Muster.
- Eine **zweite Bedingung** braucht die rohen Abschnitts-Zeilen — dann ist die
  Ausnahme aus Entscheidung 4 keine Ausnahme mehr, und Schritt 5 braucht eine
  ehrlichere Formulierung seiner Reichweite.
- Der Grenz-Fall aus den Konsequenzen tritt ein (eine real gekippte Ordnung
  über **zwei** Tabellen desselben Abschnitts) — dann ist die
  Je-Tabelle-Semantik aus Entscheidung 6 zu erweitern.

## Geschichte

- 2026-08-21: Proposed (doc-first, `slice-105`).
- 2026-08-21: nach unabhängigem Review revidiert, **bevor** der Status
  wechselt: die Typ-Mischungs-Semantik trug zwei Lesarten (Bedingungs-Tabelle
  gegen Fließtext) und ist auf die Paar-Lesart mit Anker-Reset gepinnt
  (Entscheidung 7 geschärft); ein Überlauf-Segment ist untypisierbar statt
  still kleinste Version; zwei Grenzen in den Konsequenzen nachgetragen (eine
  Chronologie-Zusage je Abschnitt; die zweite, escape-bewusste Zell-Lesart des
  trace-Lesers als andere Frage). Die bestätigende Re-Review war APPROVE ohne
  Auflagen.
- 2026-08-21: Accepted (Closure `slice-105`, Release v0.61.0; Retro-Beleg mit
  dem Produkt 27 = 14 · 6 · 7 am Vor-Heilungs-Stand, heutiger Bestand null).
