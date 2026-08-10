# ADR-0053 — Eine Bereinigung für alle Closure-Bedingungen: die Floskel-Prüfung wird dabei bewusst gelockert

**Status:** Proposed
**Datum:** 2026-08-10
**Autor:** pt9912
**Bezug:** [`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(Substanz und Floskel), [§`DC-FA-PLAN-001.a`](../../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
Schritt C4. **Verfeinert** [ADR-0048](0048-closure-note-struktur-im-planning-modul.md).
**ADR-pflichtig** nach [`AGENTS.md` §3.6](../../../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden):
eine Prüfregel wird gelockert.

## Kontext

Die Substanz-Zählung soll zum abzulösenden Adopter-Skript **paritätisch** werden:
Inline-Code zählt nicht, und ein Satzende zählt nur vor Whitespace oder
Zeilenende. Der Anlass ist eine gemessene Lücke in der **gefährlichen** Richtung
— eine Notiz, die der Adopter-Sensor als zu dünn meldet, lief bei d-check durch.

**Die Zählung und die Floskel-Prüfung lesen denselben bereinigten Text.** Wer
Inline-Code entfernt, verschärft deshalb nicht nur die Zählung, sondern
**lockert** zugleich die Floskel-Prüfung: eine Phrase in Backticks wird danach
nicht mehr gefunden. In beide Richtungen am Lauf belegt — dieselbe Notiz, alte
Fassung `closure-note-boilerplate`, neue Fassung befundfrei.

**Bestandsmessung (2026-08-10)** über die 97 eigenen Closure-Notizen, mit dem
Produkt gefahren (Schwelle hochdrehen, bis der erste Slice rot wird):

| | Minimum | bei Schwelle 4 |
|---|---|---|
| v0.55.0 | 7 | grün, Abstand 3 |
| nach der Angleichung | **5** | grün, Abstand **1** |

Und die Verteilung dessen, was einem Satzende-Zeichen folgt (4066 Vorkommen
nach Bereinigung): **1320** Whitespace oder Zeilenende, rund **2400** Punkte aus
Link-Pfaden und Versionsnummern, **170** ein `*` — also ein fett gesetztes
Satzende wie `**Umsetzung.**`.

## Entscheidung

1. **Eine Bereinigung für alle Bedingungen.** Der Abschnitt wird **einmal**
   bereinigt (Fences entfernt, Inline-Code geleert), und Substanz, Floskel und
   Platzhalter lesen **diesen einen** Text. Die Alternative — zwei getrennte
   bereinigte Texte, einer je Richtung — wäre eine zweite Semantik für
   „derselbe Abschnitt“ und genau die Klasse, die das Beobachtungs-Register
   als **BEO-003** führt.

2. **Die Lockerung der Floskel-Prüfung wird angenommen, nicht umgangen.** Eine
   *zitierte* Floskel ist keine benutzte. Die heutige Fassung meldet jede Notiz,
   die **über** Floskeln schreibt — inklusive dieser Repo-Dokumentation selbst.
   Sachlich ist die neue Semantik die bessere; sie ist trotzdem eine Lockerung
   und darum hier begründet.

3. **Ein Satzende zählt nur vor Whitespace oder Zeilenende — ohne Ausnahme für
   schließende Auszeichnung.** Die Messung zeigt, dass diese Regel 170 echte,
   fett gesetzte Satzenden mitverliert. Sie bleibt trotzdem: es ist die zugesagte
   **Parität** (der Zweck der Angleichung), die Richtung ist die sichere (zählt
   weniger ⇒ Gate strenger), und eine Ausnahmeliste für `*`, `_`, `` ` ``, `)`
   wäre eine **dritte** Semantik, die weder der Adopter noch CommonMark
   definiert.

4. **Kein Config-Ventil.** Ein Schalter wäre eine zweite Semantik für dieselbe
   Frage — und einer, den niemand bewusst setzt. Dieselbe Begründung, mit der die
   Zähl-Schwelle ohne Ventil auskommt.

5. **SemVer: Minor, mit Notiz in beide Richtungen.** `closure-note-thin` findet
   **mehr** (ein grüner Konsumentenlauf kann rot werden — dieselbe Einordnung wie
   [ADR-0042](0042-markdown-lexik-folgt-commonmark.md) und
   [ADR-0043](0043-tabellengrenze-am-relevanten-header.md)).
   `closure-note-boilerplate` findet **weniger** (wer eine Floskel zitiert stehen
   hat, verliert einen bestehenden Befund, ohne es zu merken). **Beide Richtungen
   gehören in die Release-Notiz** — die zweite ist die, die stillschweigend
   verschwindet, wenn man sie nicht hinschreibt.

## Verglichene Alternativen

| Alternative | Warum verworfen |
|---|---|
| Zwei getrennte bereinigte Texte (Zählung ohne Inline-Code, Floskel mit) | Zwei Semantiken für denselben Abschnitt; genau die Drift-Klasse, gegen die dieses Repo mehrfach angetreten ist (BEO-003). Und sie hielte eine Prüfung aufrecht, die sachlich falsch ist — eine zitierte Floskel ist keine benutzte |
| Satzende auch vor schließender Auszeichnung zählen | Verlässt die zugesagte Parität und erfindet eine dritte Semantik; die Messung zeigt, dass der Bestand die strengere Regel trägt (Minimum 5 bei Schwelle 4) |
| Config-Ventil für die Zähl-Semantik | Ein Schalter, den niemand bewusst setzt, und eine zweite Semantik für dieselbe Frage |
| Nur die Zählung ändern, Floskel-Prüfung unberührt lassen | Technisch nur über zwei Texte machbar — siehe oben. Die Kopplung ist keine Nachlässigkeit, sondern die Konsequenz **einer** Bereinigung |
| Als Patch ausliefern | Ein grüner Lauf kann rot werden **und** ein roter grün; beides ist mehr als eine Fehlerkorrektur |

## Konsequenzen

- **Der eigene Bestand bleibt grün**, aber der Abstand zur Schwelle fällt von 3
  auf 1. Die nächste dünne Notiz fällt eher auf — das ist der Zweck.
- **Ein stiller Verlust ist möglich:** ein Repo mit einer zitierten Floskel
  verliert einen bestehenden Befund. Deshalb steht die Richtung ausdrücklich in
  der Release-Notiz und nicht nur hier.
- **Die Bereinigung ist ab jetzt ein einziger Ort.** Wer eine Bedingung ergänzt,
  bekommt dieselbe Sicht automatisch — und wer sie ändern will, ändert sie für
  alle sichtbar.

## Fitness Function

- **Parität in der Zähl-Semantik:** eine Fixture, die das Adopter-Skript wegen
  der Satzzählung rot macht, macht auch das Modul rot.
- **Beide Richtungen sind testgehalten:** die zitierte Floskel trifft **nicht**,
  dieselbe Phrase im Fließtext trifft **weiter**.
- **Der eigene Bestand bleibt bei null** bei unverändertem `min-sentences: 4`.
- **Inline-Code trägt keine Substanz:** eine Notiz, deren Sätze in Backticks
  stehen, meldet `closure-note-thin`.

## Re-Evaluierungs-Trigger

- Wenn ein Konsument meldet, dass ihm durch die gelockerte Floskel-Prüfung ein
  **echter** Fall entgeht, ist zu prüfen, ob die Floskel-Bedingung einen eigenen
  Text braucht — dann ist Entscheidung 1 neu zu bewerten.
- Wenn die 170 fett gesetzten Satzenden in einem Konsumenten-Repo den Ausschlag
  geben, ist Entscheidung 3 neu zu bewerten — **mit Messung**, nicht auf Zuruf.
- Wenn eine vierte Bedingung eine **andere** Sicht auf den Abschnitt braucht,
  ist die „eine Bereinigung“-Entscheidung neu zu bewerten.

## Geschichte

- 2026-08-10: Proposed (doc-first, `slice-094`).
