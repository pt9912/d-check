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
     — es gibt genau **einen** Auswertungspfad, nicht zwei.
   - **Reihenfolge und Dedup gehören in die Spezifikation, nicht in den
     Zufall.** Befunde entstehen je Zeile in **Deklarationsreihenfolge** der
     Paare; ein Befund-Tupel, das zwei Paare identisch erzeugen, erscheint
     **einmal**. Ohne diese beiden Sätze hinge der Befundsatz an einer
     Map-Iteration — und [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)
     verlangt das Gegenteil.

   Die Befund-Form bleibt unberührt: ein Grund-Code (`version-stale`), dieselbe
   Nachricht, dasselbe `target`. Was die Paare unterscheidbar macht, ist die
   **erwartete** Version — sie steht bereits in der Nachricht.

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
