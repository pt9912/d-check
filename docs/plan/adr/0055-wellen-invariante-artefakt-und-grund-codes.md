# ADR-0055 — Wellen-Invariante: die Ergebnisnotiz ist das Artefakt, und vier Reparaturen brauchen vier Grund-Codes

**Status:** Proposed
**Datum:** 2026-08-16
**Autor:** pt9912
**Bezug:** [`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(dritte Fähigkeit); Modul-Fundament [ADR-0028](0028-planning-lifecycle-modul.md);
Trennungs-Begründung für Grund-Codes [ADR-0049](0049-structure-modul-schnitt-und-preset.md);
geteilte Lexik [ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md);
Tabellenzeilen-Lexik [`DC-FA-TGT-001`](../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)

## Kontext

Das Modul `planning` prüft die Lifecycle-Invariante heute auf der **Slice**-Ebene:
der Ruhe-Marker steht im Aktiv-Status-Block genau dann, wenn kein Slice im
Verzeichnis liegt. Eine Wellen-Datei trägt ihren Zustand im **Ort**, exakt wie
ein Slice — die Invariante ist eine Ebene höher also ebenso entscheidbar.

Vier Aussagen sind formuliert worden, und die Bestandsmessung über drei
Planungs-Bäume hat zwei davon **widerlegt, wie sie dastanden**:

- Gegen das **Plan-Dokument** gemessen meldet die Abschluss-Aussage **19-mal**
  über zwei Bäume — jedes Mal zu Unrecht: ältere Wellen sind geschlossen worden,
  bevor es die Konvention des flachen Wellendokuments gab.
- Die Vorschau-Aussage ist als Zeilen-Scan **sofort falsch**: die Trigger-Spalte
  einer Vorschau-Zeile darf andere Wellen nennen, und die eigene Roadmap tut das.

Dazu kommt eine Richtung, die beim Schnitt des Slice nicht formuliert war und
seither **eingetreten** ist: in einer Wellen-Closure fehlten drei Zeilen im
Register, obwohl alle drei Ergebnisnotizen im Ruheort lagen (Beobachtungs-Register
**BEO-001**, Zähler 2).

## Entscheidung

1. **Das verpflichtende Artefakt einer geschlossenen Welle ist die
   Ergebnisnotiz, nicht das Plan-Dokument.** Die Abschluss-Aussage prüft gegen
   sie. Begründung ist der Bestand, nicht die Ästhetik: die Notiz verlangt die
   Closure-Prozedur und ist über den ganzen Bestand vorhanden, das Plan-Dokument
   erst seit einer späteren Konvention.

2. **Plan-Dokument und Ergebnisnotiz sind zwei Rollen mit zwei Globs.** Die
   Aussagen zum **Aktiv-Status** fragen nach dem Plan-Dokument (es liegt flach,
   solange die Welle läuft), die Abschluss-Aussagen nach der Notiz. Der
   Ergebnis-Glob wird vom Plan-Glob **abgezogen**, sonst zählt jede Notiz als
   Plan-Dokument.

3. **Die Vorschau-Aussage greift nur auf der Welle-Spalte und nur bei
   Kennungen.** Zwei der drei vermessenen Bäume schreiben dort **Namen** — eine
   geplante Welle hat noch keine Kennung, sie bekommt sie bei der Eröffnung. Wo
   eine Kennung steht, greift die Aussage scharf; wo ein Name steht, gibt es
   nichts zu prüfen. Ein Token-Scan über die ganze Zeile ist ausgeschlossen.

4. **Vier Reparaturen brauchen vier Grund-Codes.** `wave-drift` (Aktiv-Status
   gegen Plan-Dokument, beide Richtungen in einer Meldung — wie beim
   Slice-Pendant), `wave-preview-exists`, `wave-results-missing`,
   `wave-unregistered`. Die Trennung ist erzwungen, nicht gewählt: zwei dieser
   Verletzungen können **dieselbe Roadmap-Zeile** treffen, und die
   Befund-Deduplikation über (Datei, Zeile, Regel, Ziel, Grund) ließe sie sonst
   zusammenfallen.

5. **Die dritte Fähigkeit ruft die bestehende Aktiv-Status-Bestimmung auf,
   statt sie zu wiederholen.** Nach
   [ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 1 ist
   eine zweite Antwort auf dieselbe Frage ein Defekt. Sind Slice- und
   Wellen-Invariante gleichzeitig verletzt, entstehen **zwei** Befunde mit
   verschiedenen Codes — kein Widerspruch, sondern zwei Reparaturen.

## Alternativen

- **Nur die Aussagen 1 und 2 liefern** (die Rückfallebene des Slice). Verworfen,
  weil die Tabellenzeilen-Lexik seit der Vorgänger-Welle entdriftet vorliegt: die
  Aussagen 3 und 4 brauchen nur noch eine Spalten-Adresse, und die Richtung, die
  **eingetreten** ist, liegt gerade in Aussage 4.
- **Ein Grund-Code für alle vier Aussagen.** Verworfen nach Entscheidung 4 — die
  Deduplikation verlöre Befunde, und die Meldung könnte die Reparatur nicht
  benennen.
- **Die Abschluss-Aussage gegen Plan-Dokument *und* Notiz prüfen.** Verworfen:
  das erzeugt am gemessenen Bestand 19 Befunde für einen Zustand, den die
  Closure-Prozedur nie verlangt hat.

## Konsequenzen

- Die Fähigkeit ist **opt-in innerhalb** des opt-in Moduls (wie die
  Closure-Fähigkeit): ohne den Aktivierungs-Schlüssel wird kein Wellen-Dokument
  geöffnet und der Befundsatz ist byte-identisch.
- **Sie findet beim ersten Lauf echte Rückstände** — im Schwester-Repo elf
  fehlende Ergebnisnotizen und eine Roadmap, die eine aktive Welle ohne
  Wellendokument nennt. Das ist der Zweck, aber es macht die Einführung dort zu
  einem eigenen Schritt.
- Die Tabellenzeilen-Lexik bekommt ihren **zweiten** Konsumenten. Nach
  [ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 4 heißt
  das: geteilte Antwort, und je Konsument eine Assertion. Beim **dritten**
  Konsumenten wird daraus ein Kopplungs-Test.

## Re-Evaluierungs-Trigger

- Ein Baum schreibt Ergebnisnotizen unter einem anderen Muster, das sich vom
  Plan-Glob nicht abziehen lässt — dann trägt die Zwei-Glob-Form nicht mehr.
- Eine geplante Welle bekommt ihre Kennung **vor** der Eröffnung (etwa durch
  eine Nummern-Reservierung). Dann ist Entscheidung 3 neu zu stellen, weil die
  Vorschau-Zeile dann regulär eine Kennung trägt.
- Eine dritte Stelle liest Tabellenzeilen. Dann ist der Kopplungs-Test fällig,
  nicht eine dritte Einzel-Assertion.

## Geschichte

- 2026-08-16: Proposed (`slice-102`, nach der Bestandsmessung über drei Bäume).
