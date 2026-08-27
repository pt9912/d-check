# Review slice-158 — Der `citations`-Scan sieht Inline-Code nicht

**Gegenstand:** [slice-158](../plan/planning/done/slice-158-citations-inline-code.md), Stand des Feat-Commits vor der Nacharbeit.
**Datum:** 2026-08-27. **Reviewer:** unabhängiger Subagent.

---

## Urteil

**Schließbar nach Nacharbeit.** Die Code-Änderung selbst ist verifiziert
korrekt; vier HIGH-Befunde stehen im **Vertragstext** und in der **Messung**,
auf der die Entscheidung ruht.

**Was der Reviewer selbst grün gefahren hat:** `make test`, `make doc-check`
(547 Dateien, 0 Befunde), `make adr-check`, `make completeness-check` (48/0).
Die zentrale Behauptung *„vorher Abbruch an `CHANGELOG.md:592`, nachher 546
Dateien / 0 Befunde / Exit 0"* ist mit einem aus dem Vorzustand **neu gebauten**
Image nachgefahren und **exakt reproduziert**, beide Seiten.

## Befunde

| # | Rang | Befund |
|---|---|---|
| F-1 | **HIGH** | *„`citations` führt keinen einzigen Konfigurations-Schlüssel"* ist falsch — das Modul trägt `citations.scope`, und das Benutzerhandbuch sagt es ausdrücklich richtig. Das Beobachtungs-Register nennt genau diesen Fall im Ableiter-Satz von `BEO-011`, den der Slice als gesichtet deklariert |
| F-2 | **HIGH** | Die ADR-Konsequenz zum Backtick-Pfad ist **widerlegt**: gemessen fällt er nicht fail-closed, sondern in einen Befund mit **leerem** Ziel. Und die tatsächlich existierende „findet mehr"-Richtung fehlt ganz — das Strippen kann eine Direktive **erzeugen** |
| F-3 | **HIGH** | Nach der Umstellung beantwortet das Produkt *„ist das eine Direktive"* **zweifach**: die vier Konsumenten des Ventils lesen weiter roh. **173** Prosa-Zeilen tragen den Ventil-Marker ausschließlich in Inline-Code. Das ist der Defekt aus ADR-0054 Entscheidung 1; der Lastenheft-Satz *„die Platzierungsregeln folgen der bestehenden Ventil-Konvention"* ist damit falsch geworden |
| F-4 | **HIGH** | Die Messtabelle reproduziert mit keiner Methode, und *„Marker-Vorkommen"* trägt **zwei unvereinbare Grundgesamtheiten**. Richtig ist unter der Marker-Basis 25/24/1/0; *„65 malformt"* ist eine Kategorien-Verwechslung, und *„74 → null"* heißt in Wahrheit *„25 → null"*. Dieselbe Wendung steht in ADR-0054 für die andere Basis |
| F-5 | MEDIUM | Der benannte Preis ist enger als das Verhalten: auch eine **freie** Direktive verschwindet still, wenn eine Spanne **desselben Absatzes** sie umschließt |
| F-6 | MEDIUM | Die Lastenheft-Historie sagt *„die Fail-closed-Menge bleibt unverändert"* — sie ist der Gegenstand der Änderung |
| F-7 | MEDIUM | §3.8: die **Ziel-Achse** bleibt unbenannt, obwohl das Modul jede Datei liest, unabhängig von Typ und Scan-Menge |
| F-8 | MEDIUM | ADR-0054s eigener Re-Evaluierungs-Trigger (*„eine vierte Stelle"*) wird nicht adressiert |
| F-9 | MEDIUM | *„ein `done/`-Slice"* / *„zwölf Fundstellen"* halten nicht — die Zahlen stammen aus einer dritten Basis |
| F-10 | MEDIUM | Doku-Currency: `CHANGELOG`, Handbuch und beide `README`-Fassungen nicht nachgezogen; die DoD führte dafür keinen Haken |
| F-11 | LOW | gemessene Menge (544 getrackte `.md`) ≠ gescannte Menge (546) |
| F-12 | LOW | Offener §3.5-Konflikt im Arbeitsbaum: der ADR-Kern war bereits umgeschrieben |
| F-13 | LOW | DoD *„ein konstruierter Verstoß je gedeckter Form"* nicht erfüllt — zwei von der ADR selbst benannte Formen waren nie gefahren |

**Zu §3 „Ausdrücklich NICHT": nichts verletzt.** Zu §5: Risiko 1 positiv
beantwortbar, **Risiko 2 eingetreten** — in beiden Ausprägungen von `BEO-011`.

## Erledigung

Alle dreizehn Befunde sind eingearbeitet:

- **F-1** korrigiert in ADR, Slice-Plan und dem wartenden Nachbar-Slice.
- **F-2** durch **Härten** statt Streichen: der Pfad im Muster beginnt mit einem
  Nicht-Whitespace-Zeichen, womit die Aussage stimmt; die Erzeugungs-Richtung
  ist als Konsequenz benannt und als Test festgenagelt.
- **F-3** als ausdrückliche Skopierung in ADR-0060 Entscheidung 7, mit
  gemessenen Rändern und
  [slice-162](../plan/planning/in-progress/slice-162-ignore-marker-geteilte-antwort.md)
  als Ausgang; der Lastenheft-Satz ist ersetzt.
- **F-4** durch eine eigene Nachmessung mit der Produkt-Lexik, deckungsgleich
  mit dem Reviewer und mit `git grep`; beide Basen stehen jetzt getrennt da.
- **F-5**, **F-7** als Vertrags-Grenzen in beiden Straten und im Handbuch.
- **F-6**, **F-9**, **F-11** durch Umformulierung auf die gemessene Größe.
- **F-8** als eingetretener Trigger in den Konsequenzen; die Gate-Frage ist eine
  DoD-Position des Folge-Slice.
- **F-10** in einem eigenen `docs(user)`-Commit, plus DoD-Position.
- **F-12** durch Neuaufbau des unveröffentlichten Feat-Commits — die ADR geht
  einmal und richtig in die Historie, statt eine falsche Fassung zu
  konservieren.
- **F-13** durch drei zusätzliche Tests.
