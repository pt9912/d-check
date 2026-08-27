# Review slice-159 — Der `d-check:ignore`-Marker hat keine Syntax

**Gegenstand:** [slice-159](../plan/planning/done/slice-159-ignore-marker-syntax.md), Stand des Feat-Commits vor der Nacharbeit.
**Datum:** 2026-08-27. **Reviewer:** unabhängiger Subagent.

---

## Urteil

**Schließbar nach Nacharbeit.** Die Verhaltensänderung ist richtig, eng und
**monoton** — die neue Bedingung impliziert die alte, es kann kein neuer stiller
Grün-Pfad entstehen. Der Widerspruch zur `diagrams`-Festlegung ist
**entschieden statt umgangen**, und die Auflösung *„dieselbe Regel auf
verschiedene Eingaben"* trägt. Die Nacharbeit betrifft die **Behauptungen um
die Änderung herum**.

## Befunde

| # | Rang | Befund |
|---|---|---|
| F-1 | **HIGH** | Die eine berührte Zeile trägt **nicht** die bare Form — der Marker steht in **Backticks** und wirkt nur durch das Paritäts-Leck (ungerade Backtick-Parität des Absatzes). Im Bestand gibt es **null** wirksame bare Marker; **65 von 66** tragen die Form bereits. ADR-0063 §Kontext sagt *„jeder"* |
| F-2 | **HIGH** | [ADR-0062](../plan/adr/0062-ventil-marker-versions-ist-eine-benannte-grenze.md)s Re-Evaluierungs-Trigger für Entscheidung 4 war **eingetreten** (*„heute tut es keine"* war beim Schreiben falsch), und dieser Commit hat die Instanz **stillschweigend beseitigt** — eine Wiedervorlage verfiel, ohne dass jemand entschied |
| F-3 | MEDIUM | `249`/`183` stammen aus der Grundgesamtheit von ADR-0062, nicht aus dieser (`BEO-020`). Richtig für HEAD: **558 / 259 / 66 / 193**. Zusätzlich nennen ADR und Historie zwei verschiedene Dateizahlen für dieselbe Messung |
| F-4 | MEDIUM | Unbenannte Grenze: ein Marker im **mehrzeiligen** HTML-Kommentar wirkt nicht — die Prüfung ist zeilenweise, das lexikalische Objekt zeilenübergreifend (§3.8) |
| F-5 | MEDIUM | Der Vertrag fordert die **Klammer**, der Code nur den **Öffner** — ein nie geschlossener Kommentar wirkt. Der Code fordert **weniger** als die Spec, also die Richtung zum stillen Grün |
| F-6 | MEDIUM | *„**jede** Erwähnungs-Form entschärft"* ist ein Overclaim (`BEO-009` b). Weiter wirksam: escapte Erwähnung, Erwähnung im HTML-Attribut, im eingerückten Code-Block. Die stärkere und haltbare Aussage ist **Monotonie** |
| F-7 | MEDIUM | `diagrams`: **keine Assertion**, obwohl [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) je Konsument eine verlangt — und das neue Kriterium war in einer Konstellation formuliert, die es bei `diagrams` nie gibt |
| F-8 | MEDIUM | Der Lexik-Wächter deckt die neue Achse **nicht**: wer den Marker auf dem **gestrippten** Text sucht, trifft die Lage und verfehlt die Form; sein Kommentar zeigt noch auf `stripInlineCodeByLine` |
| F-9 | MEDIUM | Ein Kommentar zieht `versions` in denselben Topf wie `diagrams` (*„andere Frage"*) — ADR-0062 führt es ausdrücklich als **benannte Grenze** (`BEO-012`) |
| F-10 | LOW | *„sieben neue Grenz-Proben"* — es sind **sechs** |
| F-11 | LOW | Der Nachtrag saniert §1, lässt §2/§3 mit demselben überholten Stand stehen — ausgerechnet dort, wo die DoD den **Preis** misst |
| F-12 | LOW | *„null zusätzliche Befunde über 558 Dateien"* rahmt den Beleg zu breit: der Repo-Lauf **konnte** nichts finden |
| F-13 | LOW | Beide `README`-Fassungen verlieren die Begründung der benannten Grenze |
| F-14 | LOW | Der zentrale Marker-Kommentar an der Konstanten nennt die **Form** nicht |

**Was geprüft wurde und trägt:** Monotonie der Bedingung; ein Dutzend
Regex-Kanten gemessen (Marker nach geschlossenem Kommentar → korrekt nicht
anerkannt; Marker in Backticks im Kommentar → korrekt Erwähnung); die
Bestands-Fundstelle auf der `mermaid`-Öffnungszeile geht **nicht** verloren;
`make gates` Exit 0 (eigener Lauf); Straten-Richtung §3.4 gewahrt; §3 nicht
verletzt; §2 Schritt 5 erfüllt — der Plan hatte den Fall *„Schritt 1 fällt
zugunsten der Spec aus"* vorgesehen und verlangt dann Scoping plus Begründung
im selben ADR, und ADR-0063 leistet genau das.

## Erledigung

Alle vierzehn Befunde sind eingearbeitet. **F-1 und F-2 sind eigens
nachgemessen** (221 Backticks im Absatz), nicht übernommen.

- **F-1**, **F-3**, **F-6** als `## Geschichte`-Anhang an
  [ADR-0063](../plan/adr/0063-marker-form-folgt-der-kommentar-lexik-der-eingabe.md)
  — sie war gepusht und `Accepted`, und die Befunde sind Tatsachen- und
  Reichweiten-Korrekturen, für die eine Geschichte-Zeile das richtige
  Instrument ist. Dazu an allen lebenden Stellen: Slice §1,
  Lastenheft-Historie, `CHANGELOG`.
- **F-2** als `## Geschichte`-Anhang an
  [ADR-0062](../plan/adr/0062-ventil-marker-versions-ist-eine-benannte-grenze.md):
  der Trigger war eingetreten, die Grenze **besteht fort**, nur ihre eine
  Fundstelle ist weg.
- **F-4**, **F-5** in `spec/spezifikation.md` benannt und als Proben
  festgenagelt; die Spec ist auf die **tatsächliche** Bedingung gezogen.
- **F-7** durch ein eigenes Akzeptanzkriterium in der Fence-Konstellation und
  drei `diagrams`-Proben.
- **F-8**, **F-9**, **F-14** als Kommentar-Korrekturen im Regel-Paket.
- **F-10**, **F-11**, **F-12**, **F-13** direkt behoben.
