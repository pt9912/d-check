# ADR-0058 — Drei Konfigurations-Flächen additiv weiten, statt Ersatz-Konstruktionen zu pflegen

**Status:** Proposed
**Datum:** 2026-08-22
**Autor:** pt9912
**Schärft:**
[`spec/spezifikation.md` §DC-FA-VER-001.a](../../../spec/spezifikation.md#dc-fa-ver-001a--versions-pin-konsistenz-versions)
und die Konfigurations-Schema-Zeilen unter
[`SPEC-005`](../../../spec/spezifikation.md#spec-005--d-checkyml)
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
