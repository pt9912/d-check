# ADR-0058 — Drei Konfigurations-Flächen additiv weiten, statt Ersatz-Konstruktionen zu pflegen

**Status:** Proposed
**Datum:** 2026-08-22
**Autor:** pt9912
**Schärft:**
[`spec/spezifikation.md` §DC-FA-VER-001.a](../../../spec/spezifikation.md#dc-fa-ver-001a--versions-pin-konsistenz-versions)
(Entscheidung 1),
[`spec/spezifikation.md` §DC-FA-STRUCT-001.a](../../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
(Entscheidung 2),
[`spec/spezifikation.md` §DC-FA-DIAG-001.a](../../../spec/spezifikation.md#dc-fa-diag-001a--kennungs-konsistenz-in-diagramm-fences-diagrams)
(Entscheidung 3), die Konfigurations-Schema-Zeilen unter
[`SPEC-005`](../../../spec/spezifikation.md#spec-005--d-checkyml)
und die Grund-Code-Zeile
[`SPEC-067`](../../../spec/spezifikation.md#4-grund--und-fehler-codes)
**Bezug:**
[`DC-FA-VER-001`](../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
[`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
[`DC-FA-DIAG-001`](../../../spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in)
(je Erweiterung); Schnitt-Kriterium
[ADR-0044](0044-geteiltes-referenz-ventil-quell-skopus.md);
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus) (Byte-Identität
ohne den neuen Schlüssel); Ventil-Präzedenz
[ADR-0019](0019-versions-pin-fence-ausnahme.md).

## Kontext

Drei Module tragen eine Konfigurations-Fläche, die **eine Kerbe zu schmal**
ist. Das ist nicht theoretisch aufgefallen, sondern in den beiden
vorangegangenen Wellen, je an einem Schaden:

- **`versions`** kennt genau **ein** `pin-pattern` gegen **eine**
  `current-from`-Quelle. Die Beobachtung BEO-008 (die Spiegel einer Pin-Hebung
  sind drei Klassen, gehoben wird nur die grep-bare) steht damit bei Zähler 3 —
  Schwelle erreicht — und ihr benanntes mechanisches Gegenmittel ist **nicht
  baubar**: der fremde Baseline-Tag müsste **zusätzlich** zum eigenen Image-Pin
  geprüft werden, also ein zweites Muster gegen eine zweite Quelle.
- **`structure`** kennt keinen Schlüssel „**jede** Überschrift dieses
  Abschnitts matcht dieses Muster". Die Ersatz-Konstruktion — eine
  ausgeschriebene Präfix-Negation, weil RE2 keinen Lookahead kennt — hatte ein
  **stilles Falsch-Negativ**: eine eingerückte Sektion entkam der
  Kennungs-Pflicht, und der Wächter behauptete eine Deckung, die er nicht
  hatte.
- **`diagrams`** hat weder Datei- noch Zeilen-Ventil und im
  Konfigurations-Schema **gar keine Zeile**. Ein Beispiel-Diagramm mit
  erfundener Kennung hätte über den `pre-commit`-Hook jeden Commit blockiert;
  umgangen wurde das, indem das eigene Profil das Modul auf `spec/` scopte —
  eine Vermeidung, keine Lösung.

Die drei Fälle teilen eine Form: **die Fläche ist zu schmal, und die
Vermeidung ist teurer als die Erweiterung.** Eine Ersatz-Konstruktion muss man
pflegen, erklären und bei jeder Lexik-Änderung des Moduls nachziehen; ein
Scope, der ein Ventil ersetzt, nimmt ganze Dateibäume aus der Prüfung statt
einer Zeile. Deshalb **eine** ADR mit drei Entscheidungen — nicht drei ADRs mit
derselben Begründung.

Die gemeinsame Zusage, an der jede der drei gemessen wird: **ohne den neuen
Schlüssel ist der Befundsatz byte-identisch**
([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)). Jede
Erweiterung ist additiv und opt-in; keine bestehende Konfiguration ändert ihr
Verhalten.

## Entscheidung

1. **`versions` trägt eine Liste von Muster-Quellen-Paaren
   (`versions.patterns`), und die Kurzform *ist* die einelementige Liste.**
   Kein neues Kürzel, kein neues Modul — das Kriterium ist
   [ADR-0044](0044-geteiltes-referenz-ventil-quell-skopus.md) entnommen
   (querschnittlich ⇒ eigenes Kürzel, Einzelmodul-Frage ⇒ bestehende
   Anforderung ändern). Vier Festlegungen tragen sie:

   - **Das Paar ist die Einheit, nicht der Schlüssel.** Muster, Quelle und
     `exempt-paths` gehören zusammen; eine modulweite Ausnahmeliste über
     mehreren Reihen wäre nicht ablesbar (welcher Pfad ist wegen welcher Reihe
     ausgenommen?) und beim Hinzufügen einer Reihe eine stille Falle.
   - **Die Selbst-Ausnahme der Quell-Datei ist paar-lokal.** Das
     Release-Register ist von **seiner** Reihe ausgenommen, weil sein Verlauf
     bewusst alle Versionen listet — nicht von einer fremden. Wäre sie
     modulweit, machte jede neue Quelle zugleich ein neues blindes Feld auf.
   - **Kurzform und Liste zugleich sind fail-closed (Exit 2).** Nicht, weil
     eine Zusammenführung unmöglich wäre, sondern weil sie **geraten** werden
     müsste: gilt die Kurzform als erstes oder als letztes Paar, und wirken
     ihre `exempt-paths` auf die anderen? Eine Voreinstellung, die man raten
     muss, ist keine. Die Kurzform wird intern in die Ein-Paar-Liste übersetzt
     — es gibt genau **einen** Auswertungspfad, nicht zwei. Ob eine
     Schreibweise gesetzt ist, entscheidet die **Anwesenheit** des Schlüssels,
     nicht sein Wert: sonst schaltete ein leer gelassener Kurzform-Schlüssel
     die Prüfung still auf die andere Schreibweise um. Ein Schlüssel **ohne
     Wert** bleibt im YAML von einem fehlenden ununterscheidbar — als Grenze
     benannt, nicht als Zusage überdehnt.
   - **Eine Befund-Adresse, alle Erwartungen.** Die Befund-Adresse (Datei,
     Zeile, Regel, `target`, Grund-Code) unterscheidet zwei Befunde an
     derselben Stelle **nicht**, und die geteilte Nachrunde verwirft den
     zweiten — samt seiner Erwartung. Statt die Adresse zu erweitern (das
     änderte die Befund-Form aller Module) trägt die **Nachricht** jede
     Erwartung mit ihrer Quelle, in Deklarationsreihenfolge der Paare. Die
     **Ausgabe**-Reihenfolge bleibt die der geteilten Sortierung; sie ist
     bereits deterministisch
     ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)), und
     die Deklarationsreihenfolge hat auf sie keine Wirkung.
   - **Der Wortlaut hängt an der Paar-Zahl, nicht an der Schreibweise.** Bei
     genau einem Paar nennt jede Meldung den Kurzform-Schlüssel wie bisher, ab
     zwei Paaren die Fundstelle `versions.patterns[i]`. Damit ist die
     Byte-Identität bestehender Konfigurationen eine Eigenschaft des Codes und
     nicht eine Absicht — und ein Adopter mit mehreren Reihen erfährt, welches
     Paar spricht.

   Die Befund-**Form** bleibt unberührt: ein Grund-Code (`version-stale`),
   dieselben Felder, dasselbe `target`. Was die Reihen unterscheidbar macht,
   ist allein die Nachricht — und genau deshalb muss sie vollständig sein.

2. **`structure` prüft Überschriften positiv und je Überschrift
   (`heading-pattern`/`heading-level`) — nicht als Negation über den
   Abschnitts-Text.** Auch das ist eine achte Bedingung derselben Anforderung,
   nach demselben Kriterium
   ([ADR-0044](0044-geteiltes-referenz-ventil-quell-skopus.md)). Vier
   Festlegungen:

   - **Die Bedingung liest die Überschriften, nicht den Abschnitts-Text.** Sie
     ist damit die einzige neben der Chronologie-Bedingung, die vom
     bereinigten Text abweicht — und die einzige, die dieselbe
     Heading-Erkennung benutzt, mit der das Modul den Abschnitt findet. Das
     ist keine Bequemlichkeit, sondern die Reparatur des Anlasses: die
     abgelöste Negation hatte ihre eigene Lexik nachgebaut (Zeilenanfang,
     genau ein Leerzeichen), während das Modul beliebigen Weißraum trimmt und
     Tab akzeptiert. Eine eingerückte Sektion entkam still
     ([ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md): wer eine
     Lexik-Frage selbst beantwortet, ist ein Defekt, keine Variante).
   - **Der Schlüssel heißt `headings-match`, nicht `heading-pattern`.** Unter
     `planning.closure` trägt `heading-pattern` bereits einen **Selektor**
     („welche Überschrift eröffnet den geprüften Abschnitt"); hier wäre es
     eine **Bedingung** („welche Form müssen die Überschriften haben"). Beide
     Blöcke leben in diesem Repo im **selben** Profil
     (`.d-check.closure.yml`), und die Spezifikation weist die
     Closure-Fähigkeit als *Preset derselben Semantik* aus — Namensgleichheit
     läse sich dort als Bedeutungsgleichheit. Wer den Selektor-Wert
     (`'^#{2,3} .*Closure-Notiz'`) in eine `structure`-Regel überträgt, prüft
     ihn gegen den Überschriften-**Text**, der nie ein `#` enthält: jede
     Unterüberschrift rot, aus einem Grund, den die Meldung nicht erklärt. Der
     Plural sagt zugleich die Semantik — **alle** Überschriften, nicht eine.
   - **Geprüft wird der Überschriften-Text, nicht die rohe Zeile.** Sonst
     müsste jedes Muster die `#`-Folge mitkodieren — und damit die Ebene, die
     bereits ein eigener Schlüssel ist. Zwei Orte für dieselbe Aussage sind
     die bekannte Drift-Quelle. Der Abschnitts-**Selektor** vergleicht
     weiterhin die getrimmte Zeile; er wählt einen Abschnitt, diese Bedingung
     prüft eine Form.
   - **Ein Befund je Überschrift, auf ihrer Zeile.** Die übrigen Bedingungen
     melden am Abschnittskopf, weil ihre Aussage dem Abschnitt gilt; diese
     gilt einer Überschrift, und die Reparatur steht dort. Das ändert die
     Befund-**Zahl** gegenüber einer Abschnitts-Bedingung und ist deshalb
     zugesagt statt stillschweigend. Eigener Grund-Code
     `section-heading-mismatch` nach dem Kriterium der Bedingungs-Tabelle
     (andere Reparatur ⇒ eigener Code); ein Sammel-Code fiele zudem mit den
     übrigen Bedingungen desselben Abschnitts unter die Befund-Deduplikation
     zusammen.
   - **Default-Ebene ist Abschnitts-Ebene + 1.** Gemessen am eigenen Bestand
     macht das heute keinen Unterschied (der betroffene Abschnitt trägt
     ausschließlich Überschriften einer Ebene); der Unterschied entsteht in
     der Zukunft. „Alle Ebenen" hieße, dass eine später ergänzte vierte Ebene
     rückwirkend unter eine Kennungs-Pflicht fiele, die ihr nie zugedacht war
     — ein Falsch-Positiv mit Verzögerung. Wer sie meint, nennt sie.

   **Drei Grenzen sind benannt, nicht geschlossen:** ohne Überschrift der
   geprüften Ebene ist die Bedingung vacuously wahr; ein `headings-level`
   flacher als der Abschnitt kann in ihm nicht vorkommen; und **zwei Ebenen in
   einer Regel sind heute nicht bloß unvorgesehen, sondern gesperrt** — eine
   zweite Regel über demselben Abschnitt teilt dessen Regel-Identität
   (`files :: Selektor`) und bricht mit Exit 2. Wer zwei Ebenen prüfen will,
   braucht die Ebenen-**Menge** aus dem Re-Evaluierungs-Trigger, keinen
   Workaround. Beides sind keine
   Defekte, aber beides kann eine Deckung vortäuschen — deshalb wird die
   Umstellung des eigenen Profils an einem Fall gemessen, an dem die abgelöste
   Negation **still** ist: eine Überschrift innerhalb eines mehrzeiligen
   Inline-Code-Spans. Die Bereinigung räumt die Zeile leer, die Negation sieht
   nichts; die Überschriften-Erkennung sieht die Überschrift — wie `anchors`
   auch, das ihr einen Slug gibt.

3. **`diagrams` bekommt beide Ventile — und der Zeilen-Marker wirkt auch auf
   der Fence-Öffnungszeile.** Dritte Erweiterung nach demselben Kriterium.
   Drei Festlegungen:

   - **Ventil-Parität ist ein Prinzip, keine Nachrüstung.** Ein Modul, das
     **eigene Muster** konfiguriert und Befunde an Zeilen hängt, braucht ein
     Zeilen-Ventil — sonst bleibt dem Nutzer nur der Scope, und der ist keine
     Ausnahme, sondern eine Vermeidung: er nimmt ganze Dateibäume aus der
     Prüfung, wo eine einzelne Referenz gemeint war. Gemessen am eigenen
     Profil: das Modul musste auf `spec/` gescopt werden, weil ein
     Beispiel-Diagramm sonst jeden Commit blockiert hätte. **Die Reichweite
     ist bewusst eng gezogen:** `hostpaths`, `pins` und `spans` hängen Befunde
     ebenfalls an Zeilen und tragen weiterhin kein Ventil. Sie konfigurieren
     keine eigenen Muster — ihr Befund folgt aus einer festen Lexik, und ob
     dieselbe Antwort dort richtig ist, ist eine eigene Frage mit eigener
     Messung. Als **offene Fläche** benannt, nicht als Zusage überdehnt.
   - **Der Marker ist ein Token, kein HTML-Kommentar.** In Prosa steht er in
     `<!-- … -->`; in einem `mermaid`-Fence ist das kein Kommentar, sondern
     Diagramm-Text. Das Modul sucht deshalb das **Token** auf der Zeile — wie
     der Autor es vor dem Renderer versteckt, ist Sache der Diagramm-Sprache
     (Mermaid: `%%`). Die Alternative — eine sprach-spezifische
     Kommentar-Erkennung je Fence-Sprache — wäre ein Grammatik-Parser durch
     die Hintertür und widerspräche der Modul-Grenze („reine Token-Extraktion
     über Rohtext").
   - **Der Marker wirkt auf der Öffnungszeile für den ganzen Block.** Ohne
     diese zweite Stelle wäre die **intuitive** Platzierung — am
     Diagramm-Anfang — wirkungslos, und der Fall, der die Erweiterung
     auslöste, bräuchte den Marker auf jeder Kennungs-Zeile. Ein Ventil, das
     man N-mal setzen muss, benutzt niemand; ein Ventil, dessen naheliegende
     Platzierung still nichts tut, ist schlimmer als keines. Es bleibt bei
     **einem** Mechanismus mit zwei Orten, nicht bei zwei Schlüsseln.

   **Nachgetragen wird zugleich eine ältere Lücke:** die Schlüssel des Moduls
   standen bis hierher **nur** im Algorithmus-Abschnitt, nicht im
   §2-Konfigurations-Schema — als einziges Modul. Das ist derselbe Vertrag,
   nur an der Stelle, an der ein Adopter ihn sucht.

## Konsequenzen

**Positiv.** Die 3×-Form von BEO-008 wird **baubar** — ob das eigene Profil sie
fährt, bleibt ein eigener Entscheid mit eigener Messung. Repos mit mehreren
Versions-Reihen (eigenes Release **und** ein gepinnter fremder Stand — der
Normalfall in diesem Harness) brauchen keine zweite Prüf-Konfiguration mehr.

**Negativ / Kosten.** Zwei Schreibweisen für denselben Sachverhalt sind eine
Drift-Quelle. Sie wird auf den Config-Rand begrenzt: der Kern kennt nur die
Liste, die Übersetzung passiert einmal beim Laden, und die Byte-Identität
zwischen beiden Schreibweisen ist ein Akzeptanzkriterium mit Test — nicht eine
Absicht.

**Grenze.** Mehrere Quellen **je Paar** bleiben draußen: die Prüfung vergleicht
auf Gleichheit, und zwei Quellen hätten keine eindeutige Erwartung. Ebenso
bleibt es bei Gleichheit statt semantischer Ordnung.

**Positiv (Entscheidung 2).** Die abgelöste Ersatz-Konstruktion verschwindet
aus dem eigenen Profil — mit ihr eine 240 Zeichen lange, von Hand hergeleitete
Negation, die niemand ohne Gegenprobe lesen kann. Die Aussage steht jetzt so
da, wie sie gemeint ist.

**Negativ / Kosten (Entscheidung 2).** Eine achte Bedingung ist eine achte
Fläche: Vertrag, Schema-Zeile, Grund-Code, Klartext, Vorlage. Und die Bedingung
weicht als zweite vom bereinigten Abschnitts-Text ab — die Aussage „jede
Bedingung liest denselben Text" gilt damit für sechs von acht und ist an
beiden Ausnahmen ausdrücklich benannt.

**Positiv (Entscheidung 3).** Wer `diagrams` aktiviert, muss nicht mehr
zwischen „ganzen Baum ausnehmen" und „Gate abschalten" wählen. Ob das eigene
Profil danach ohne Scope auskommt, ist **gemessen zu entscheiden** und nicht
Teil dieser Entscheidung.

**Grenze der Additivität (Entscheidung 3).** Nur das Datei-Ventil hängt an
einem **Schlüssel**; der Zeilen-Marker hängt am **Inhalt** — wie bei
`ids`/`codepaths`/`versions`. Eine Diagramm-Zeile, die die Zeichenfolge
ohnehin trägt (etwa ein Knoten-Label, das den Mechanismus dokumentiert), wird
damit **ohne** Konfigurations-Änderung stumm. Die Byte-Identitäts-Zusage
dieser Welle gilt für die Marker-Hälfte deshalb nur für Bäume ohne diese
Zeichenfolge in einer gelisteten Fence. Das ist der Preis dafür, denselben
Marker zu benutzen statt einen eigenen zu erfinden — und er ist hier benannt,
statt in der Zusage unterzugehen.

**Negativ / Kosten (Entscheidung 3).** Ein Marker, der auf zwei Orten wirkt,
ist eine Regel mehr, die man lesen muss. Sie steht dafür in Anforderung,
Algorithmus und Schema — und die Alternative (zwei Schlüssel für Zeile und
Block) wäre eine Fläche mehr statt einer Regel mehr.

**Zweite Grenze, gemessen statt vermutet.** Zwei Paare, die auf derselben Zeile
denselben Pin-Wert treffen, teilen eine Befund-Adresse und ergeben **einen**
Befund. Ein maschineller Konsument, der die Ausgabe nach Regeln filtert, sieht
dort eine Zeile statt zweier; die zweite Erwartung liest er nur aus der
Nachricht. Das ist der Preis dafür, die Befund-Form nicht anzufassen — die
Alternative (ein Adress-Feld für die Quelle) beträfe alle Module und ist eine
eigene Entscheidung mit eigener Messung.

## Alternativen

- **Ein zweites Modul** (`versions2` o. ä.) — verworfen: dieselbe Frage,
  dieselbe Eingabe-Achse, dieselbe Befund-Form. Das ist die Definition einer
  Erweiterung nach [ADR-0044](0044-geteiltes-referenz-ventil-quell-skopus.md).
- **Ein Muster mit Alternation** (`(a|b)`) gegen eine Quelle — verworfen: die
  Reihen haben **verschiedene erwartete Versionen**; ein gemeinsames Muster
  hätte nur eine Erwartung und meldete die zweite Reihe durchgehend rot.
- **Kurzform als implizites erstes Paar** — verworfen: dann müsste man wissen,
  ob die Kurzform-`exempt-paths` nur für sie oder für alle gelten. Genau die
  Frage, die niemand aus der Datei ablesen kann.
- **Kurzform abschaffen und alle Konsumenten migrieren** — verworfen für diese
  Version: ein Breaking Change ohne Gegenwert, wo die Übersetzung eine Zeile
  Code ist.
- **Die Befund-Adresse um die Quelle erweitern**, damit zwei Paare zwei Befunde
  ergeben — verworfen für diese Version: das Feld läge in
  [`SPEC-001`](../../../spec/spezifikation.md#spec-001--befund) und beträfe
  jedes Modul und jede Ausgabe. Zwei Repo-Module umgehen dieselbe Enge heute
  über eigene Grund-Codes; ob das die richtige Antwort bleibt, ist eine
  querschnittliche Frage — siehe Re-Evaluierungs-Trigger.

## Re-Evaluierungs-Trigger

- Ein **dritter** Konsument braucht zwei Befunde an derselben Adresse (heute:
  die eigenen Grund-Code-Umgehungen in `structure` und `planning`, dazu diese
  Nachricht) — dann ist die Befund-Adresse selbst zu prüfen, nicht ein
  weiterer Umweg.
- Ein Adopter meldet, dass er die **zweite Erwartung** maschinell braucht
  (Filter über Regel und Ziel statt über die Nachricht) — dann trägt die
  Grenze aus den Konsequenzen nicht mehr.
- Die Kurzform ist im Bestand der Konsumenten **nicht mehr in Gebrauch** —
  dann ist Entscheidung 1 auf die Liste allein zu verengen und der zweite
  Schreibweg zu entfernen.
- Ein Paar braucht wiederholt **mehrere Quellen** — dann ist die erste Grenze
  gegen die gelebte Praxis zu prüfen; die Antwort wäre eine Semantik für
  „welche der Quellen gilt", keine zweite Liste.
- Eine **dritte** Bedingung weicht vom bereinigten Abschnitts-Text ab — dann
  ist nicht die Bedingung die Ausnahme, sondern der Satz „jede Bedingung liest
  denselben Text"; er braucht dann eine ehrlichere Formulierung seiner
  Reichweite.
- Ein Adopter braucht wiederholt **mehrere Ebenen** in einer Regel — dann ist
  die Ein-Ebenen-Wahl aus Entscheidung 2 gegen die Praxis zu prüfen; die
  Antwort wäre eine Ebenen-**Menge**, nicht ein zweiter Schlüssel.
- Eine Diagramm-Sprache mit **eigener Kommentar-Syntax** wird so verbreitet
  konfiguriert, dass Nutzer den Marker regelmäßig sichtbar im Diagramm haben —
  dann ist die Token-Entscheidung aus Entscheidung 3 gegen die Praxis zu
  prüfen; die Antwort wäre eine Kommentar-Lexik je Sprache, kein Parser.
- Ein Modul **ohne** konfigurierbare Muster (`hostpaths`, `pins`, `spans`, `citations`)
  braucht wiederholt eine Ausnahme — dann ist die enge Reichweite aus
  Entscheidung 3 gegen die Praxis zu prüfen; die Antwort wäre dieselbe
  Ventil-Achse, nicht eine dritte.
- Der **Block**-Ort des Markers wird häufiger gebraucht als der Zeilen-Ort —
  dann ist zu prüfen, ob der Block-Fall einen eigenen, sichtbareren Ausdruck
  verdient statt einer zweiten Bedeutung derselben Marke.
- Die **vacuously wahre** Bedingung trifft real (ein Abschnitt verliert seine
  Unterabschnitte, ohne dass es jemand merkt) — dann braucht die Bedingung
  eine Mindest-Zahl, und die Grenze aus Entscheidung 2 ist keine mehr.
