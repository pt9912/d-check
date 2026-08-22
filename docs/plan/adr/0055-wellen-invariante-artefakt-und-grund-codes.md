# ADR-0055 — Wellen-Invariante: die Ergebnisnotiz ist das Artefakt, und vier Reparaturen brauchen vier Grund-Codes

**Status:** Proposed
**Datum:** 2026-08-16
**Autor:** pt9912
**Bezug:** [`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(dritte Fähigkeit); Modul-Fundament [ADR-0028](0028-planning-lifecycle-modul.md);
Trennungs-Begründung für Grund-Codes [ADR-0049](0049-structure-modul-schnitt-und-preset.md);
geteilte Lexik [ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md);
Tabellenzeilen-Lexik [`DC-FA-TGT-001`](../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
**Schärft:** die Wellen-Schritte W1–W5 in
[`spec/spezifikation.md` §DC-FA-PLAN-001.a](../../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
und die vier Grund-Code-Festlegungen [`SPEC-045`](../../../spec/spezifikation.md#4-grund--und-fehler-codes), [`SPEC-046`](../../../spec/spezifikation.md#4-grund--und-fehler-codes),
[`SPEC-047`](../../../spec/spezifikation.md#4-grund--und-fehler-codes) und [`SPEC-048`](../../../spec/spezifikation.md#4-grund--und-fehler-codes).

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

6. **Die erste Aussage kennt zwei Kardinalitäts-Modelle
   (`planning.waves.mode: one | many`, Fortschreibung 2026-08-21).**
   `one` (Default) ist das unveränderte Singleton-Prädikat aus
   Entscheidung 4 — ohne den Schlüssel bleibt der Befundsatz
   byte-identisch. `many` deckt das Offene-Wellen-Modell des Adopters
   (Baseline v5.7.0: die Liste folgt den **Dateien**, der Ruhe-Marker
   folgt dem **Anspruch** und steht zusätzlich): verglichen werden
   **Kennungs-Mengen** — die im `planning.heading`-Block genannten
   Kennungen gegen die Kennungs-Menge der flachen Wellendokumente (W2:
   zwei Dateien derselben Kennung sind ein Element — das gilt ebenso für
   das Singleton unter `one`), beide Richtungen, jede
   Kardinalität einschließlich null; der Marker geht **nicht** ein, seine
   Aussage liegt vollständig bei `planning-drift` (Entscheidung 5 bleibt:
   die Aktiv-Status-Bestimmung wird nicht wiederholt; die Kennungs-Liste
   liest ein **eigener Scan über die geteilte Block-Grenze** — dieselbe
   Grenze, die auch die Marker-Suche nutzt, als die eine Antwort auf die
   Block-Frage). **Erkennungs-Verfahren** wie in den
   Registern: literales Glob-Präfix plus Ziffernfolge (`waveID`),
   zeilenweise über die **Prosa-Zeilen** des Blocks — Fence-Inhalte zählen
   nicht, Mehrfachnennung zählt einmal, layout-agnostisch für Tabellen-
   wie Listen-Form. **Kein zweiter Grund-Code:** die Reparatur ist
   dieselbe („Roadmap nachziehen oder Datei verschieben"), und die
   Dedup-Begründung aus Entscheidung 4 trägt weiter, weil das
   Befund-`target` unter `many` die **betroffene Kennung** ist — zwei
   Richtungen an derselben Zeile bleiben im Tupel (Datei, Zeile, Regel,
   Ziel, Grund) verschieden, wie es die Register-Aussagen praktizieren.
   **fail-closed:** ein unbekannter oder explizit leerer Modus bricht mit
   Exit 2 und Schlüssel-Nennung (Zeiger-Disziplin der übrigen
   `waves`-Schlüssel).

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
- **Das Singleton-Prädikat im Default zur Bijektion umbauen** (statt
  opt-in-Modus). Verworfen mit der Fortschreibung: bestehende Konsumenten
  können sich auf die `one`-Strenge verlassen (der Zustand „Welle offen,
  nichts beansprucht" **soll** dort rot sein), und der CR des Konsumenten
  beantragt ausdrücklich Default-Treue
  ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus):
  ohne Schlüssel byte-identisch).
- **Ein zweiter Grund-Code je Bijektions-Richtung.** Verworfen: die
  Reparatur-Klasse ist identisch, die Richtungen unterscheidet das
  `target` (Kennung) — ein zweiter Code widerspräche der eigenen
  Vier-Codes-Begründung, die Codes an **Reparaturen** bindet, nicht an
  Richtungen.

## Konsequenzen

- Die Fähigkeit ist **opt-in innerhalb** des opt-in Moduls (wie die
  Closure-Fähigkeit): ohne den Aktivierungs-Schlüssel wird kein Wellen-Dokument
  geöffnet und der Befundsatz ist byte-identisch.
- **Sie findet beim ersten Lauf echte Rückstände** — im Schwester-Repo elf
  fehlende Ergebnisnotizen (robust, mit konsument-gerechtem Marker
  nachgemessen). Das ist der Zweck, aber es macht die Einführung dort zu einem
  eigenen Schritt — samt konsument-gerechtem Ruhe-Marker: der erste Probe-Lauf
  meldete einen **zwölften** Befund, der ein Artefakt des Default-Markers war,
  nicht des Bestands.
- Die Tabellenzeilen-Lexik bekommt ihren **zweiten** Konsumenten. Nach
  [ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 4 heißt
  das: geteilte Antwort, und je Konsument eine Assertion. Beim **dritten**
  Konsumenten wird daraus ein Kopplungs-Test.
- **Unter `many` wird die Sektions-Prosa des `planning.heading`-Blocks zur
  Messfläche** (Fortschreibung): jede dort genannte Wellen-Kennung zählt
  als Zeiger. Ein Adopter, der den Modus setzt, hält seine Sektionsregel
  kennungsfrei — dieselbe Paraphrase-Disziplin, die für den
  Marker-Wortlaut bereits gilt. Dieses Repo stellt seine eigenen
  Prüf-Profile erst **nach** Release und Digest-Backfill auf `many` um
  (der gepinnte Prüfer muss den Schlüssel kennen, sonst Exit 2).

## Re-Evaluierungs-Trigger

- Ein Baum schreibt Ergebnisnotizen unter einem anderen Muster, das sich vom
  Plan-Glob nicht abziehen lässt — dann trägt die Zwei-Glob-Form nicht mehr.
- Eine geplante Welle bekommt ihre Kennung **vor** der Eröffnung (etwa durch
  eine Nummern-Reservierung). Dann ist Entscheidung 3 neu zu stellen, weil die
  Vorschau-Zeile dann regulär eine Kennung trägt.
- Eine dritte Stelle liest Tabellenzeilen. Dann ist der Kopplungs-Test fällig,
  nicht eine dritte Einzel-Assertion.
- Ein Adopter braucht eine **dritte** Kardinalitäts-Semantik oder eine
  Kennungs-Erkennung jenseits des literalen Glob-Präfixes — dann ist
  Entscheidung 6 neu zu stellen, nicht still zu erweitern.

## Geschichte

- 2026-08-16: Proposed (`slice-102`, nach der Bestandsmessung über drei Bäume).
- 2026-08-21: **Fortgeschrieben** (`slice-111`): Entscheidung 6 —
  `planning.waves.mode: one | many` auf **formalen Konsumenten-CR**
  (ai-harness-course, „planning.waves: Bijektion statt Singleton";
  Anlass: Baseline v5.7.0 „Zwei Hälften, ein Wächter" macht den Marker
  zur Anspruchs-Aussage **zusätzlich** zur Liste; Messung team-sim
  s04a–s04d, 11/11 PASS — der Singleton beißt bei zwei offenen Wellen,
  die Marker-Hälfte hält beidseitig). Vertrags-Text im Lastenheft 0.62.0;
  Alternativen um Default-Umbau und Zweit-Code ergänzt (beide verworfen),
  Konsequenzen um die Prosa-Messfläche, Trigger um die dritte
  Kardinalitäts-Semantik. **Wortlaut-Korrektur auf Review-Auflage F-1**
  (vor dem Release): die Kennungs-Liste liest ein eigener Scan über die
  geteilte Block-Grenze — nicht der Rückgabewert der
  Aktiv-Status-Bestimmung, wie die Erst-Fassung behauptete.
- 2026-08-22: **Präzisiert** (`slice-112`): Entscheidung 6 — die
  Vergleichsgröße ist in **beiden** Modi die Kennungs-Menge (W2: zwei Dateien
  derselben Kennung sind ein Element; Lastenheft 0.62.1). Kein
  Verhaltens-Unterschied, kein Release; der Beleg ist ein Pinning-Test mit
  zwei gleich-kennigen Dateien.
- 2026-08-22: `Schärft:`-Feld **ergänzt** (es fehlte): die vier Grund-Code-
  Festlegungen bei ihrer Struktur-Kennung, der Algorithmus-Abschnitt bei
  seiner Verfeinerungs-Kennung. Adressierungs-Form der Baseline, Kern
  unverändert, Status bleibt `Proposed`.
