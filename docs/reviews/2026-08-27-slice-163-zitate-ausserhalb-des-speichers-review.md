# Review slice-163 — Baseline-Zitate außerhalb des Konventionsspeichers

**Gegenstand:** [slice-163](../plan/planning/done/slice-163-zitate-ausserhalb-des-speichers.md), Stand der zwei Feat-Commits vor der Nacharbeit.
**Datum:** 2026-08-27. **Reviewer:** unabhängiger Subagent.

---

## Urteil

**Schließbar nach Nacharbeit.** Die vier Direktiven sind handwerklich sauber —
alle vier live, richtig gepaart, Spannen minimal, mit eigenen Bruchproben
verifiziert. `make gates` grün (eigener Lauf). **Die Zahlen 44 / 9 / 6 stimmen;
die Zahl 4 nicht**, und eine der drei geforderten Klassen wurde gar nicht
berichtet.

## Befunde

| # | Rang | Befund |
|---|---|---|
| F-1 | **HIGH** | *„nur vier sind adressierbar"* ist durch Messung widerlegt — es sind **sechs**. Die Begründung hält dreifach nicht: eine Tabelle **ist** ein Absatz; `AGENTS.md:351` trägt **genau ein** Zitat; und die Form, die eine Tabellenzeile nicht bricht, hat [ADR-0040](../plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md) längst entschieden |
| F-2 | **HIGH** | Die Klasse **„abweichend"** wurde nicht berichtet, und sie ist nicht leer: das `BEO-015`-Kanon-Zitat lässt die Fettung des zweiten Halbsatzes weg. Der Trichter las sich vollständig |
| F-3 | MEDIUM | Die gemessene Menge ist nicht die, über die geredet wird (`BEO-020`): gemessen sind **sechs** Dateien, [`MR-039`](../../harness/conventions.md#mr-039) nennt die lebenden **Planungs**-Dokumente mit. Richtig: **49**, nicht 44 |
| F-4 | MEDIUM | Das erste **Nicht-Baseline-Ziel** hat keine Regel — [`MR-051`](../../harness/conventions.md#mr-051) deckt nur Baseline-Ziele, und zwei Spiegel behaupten weiter *„die planmäßige Rot-Quelle ist der Bump"* |
| F-5 | MEDIUM | `MR-052` schreibt [`MR-039`](../../harness/conventions.md#mr-039) eine Annahme zu, die dort nicht steht (`BEO-012`) — Träger war die **Praxis** in slice-152, und die ist dort bereits protokolliert |
| F-6 | LOW | `MR-052` führt ein Feld `Schärft:` ein, das die Vorlage nicht kennt und kein anderer Eintrag benutzt — der Titel sagt es bereits |
| F-7 | LOW | Die zwei Direktiven unter Überschriften haben die Leerzeile aufgefressen — unnötig, sie hätte bleiben können |
| F-8 | LOW | Die Zahl 44 gehört zum **rohen** Lauf, §2 Schritt 1 deklariert die gestrippte Lexik. Die Wahl ist sachlich richtig (der Zitattext wird roh gelesen), die Methodenangabe sagt das Gegenteil |
| F-9 | LOW | *„Ein konstruierter Verstoß je gedeckter Form"* ist zur Hälfte belegt — der Slice führt **zwei** Platzierungsformen ein, dokumentiert ist eine |

**Was geprüft wurde und trägt:** die vier Direktiven (live, richtig gepaart,
minimal); die Delta-Behauptung von `MR-052` ist **belegt** — `git show` liefert
die `v5.11.0`-Zeile, `MR-039`s Tabelle ist wörtlich, `MR-033` hat die Fettung
verschoben; die zweite Delta-Zeile ist ebenfalls wörtlich; §3 ist eingehalten;
`MR-052` als **neuer Eintrag** ist die richtige Form (ein
`## Geschichte`-Anhang an `MR-039` wäre falsch gewesen); keine Spec-/Lastenheft-
Änderung nötig, keine `CHANGELOG`-Pflicht.

## Erledigung

Alle neun Befunde sind eingearbeitet. **F-1 und F-2 sind eigens nachgemessen.**

- **F-1** durch **Auszeichnen** statt Umformulieren: beide Tabellenzeilen tragen
  jetzt ihre Direktive in der [ADR-0040](../plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md)-Form,
  beide mit Bruchprobe belegt. Der *richtige* Nicht-Adressierbar-Fall (`BEO-015`,
  viertes Zitat seiner Zeile) ist als solcher benannt.
- **F-2** korrigiert; damit ist die Klasse **wirklich** leer, und die Zahl der
  wörtlichen Treffer steigt auf zehn.
- **F-3** nachgemessen: **49 / 10 / 3 Marker / 7 Zitate / 6 ausgezeichnet**.
- **F-4** in [`MR-051`](../../harness/conventions.md#mr-051) ergänzt, samt beiden
  Spiegeln.
- **F-5**, **F-6** in `MR-052` behoben — er war noch nicht veröffentlicht.
- **F-7**, **F-9** direkt behoben; **F-8** in der Closure-Notiz benannt.
