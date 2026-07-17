# ADR-0040 — Abschließender HTML-Kommentar in Tabellenzeilen ist ein Suffix, keine Zelle

**Status:** Proposed
**Datum:** 2026-07-17
**Autor:** pt9912
**Schärft:** [`DC-FA-REQ-001.a`](../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
**Bezug:** [`DC-FA-REQ-001`](../../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen), [`DC-FA-XREF-001`](../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in), [ADR-0037](0037-trace-tabellenquellen-nullmengen-guard.md), [ADR-0038](0038-trace-cross-consistency.md), [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus).

## Kontext

Der header-gebundene Tabellen-Reader ist **fail-closed**: eine Datenzeile muss
dieselbe Zellenzahl wie ihr Header tragen, sonst Exit 2
([ADR-0037](0037-trace-tabellenquellen-nullmengen-guard.md) — der Guard verhindert,
dass eine Tabelle still abreißt und Anforderungen lautlos verschwinden).

d-check kennt zugleich eine **Direktiven-Konvention**: `<!-- d-check:ignore (Grund) -->`
nimmt eine Zeile von der Prüfung aus. In einer Tabellenzeile **muss** dieser Marker
hinter der schließenden Pipe stehen — in einer Zelle wäre er Zellinhalt und damit
Teil des Titels bzw. der Kanten-Menge.

Beides zusammen ist ein Selbstwiderspruch: der Reader zählt den Marker als
zusätzliche Zelle und bricht mit Exit 2 ab. **d-checks eigene Konvention macht
d-checks eigenen Reader blind.**

Belegt am Realdatenlauf gegen den Konsumenten grid-gym (2026-07-17): eine reale
3-Spalten-Tabelle in `spec/architecture.md` trägt

```markdown
| E2E / Demo-Abnahme | `tests/e2e/demo` | [`GG-TESTTYPE-005`](…) | <!-- d-check:ignore (geplant: …) -->
```

und ließ v0.45.1 mit `Tabellenzeile 913 hat 4 statt 3 Zellen` abbrechen — die
gesamte Rück-Sicht war unlesbar, der Kreuzverweis-Abgleich nicht fahrbar. Isoliert
reproduziert für `trace.cross-consistency` **und** für
`trace.requirements.format: table` (dort **ausgeliefert seit v0.43.0**).

Verschärfend: **GFM ignoriert überzählige Zellen** — die Zeile rendert in jedem
Markdown-Viewer normal. d-check ist hier strenger als der Standard, ausgerechnet
gegenüber der eigenen Empfehlung.

## Entscheidung

1. **Ein abschließender HTML-Kommentar, der die ganze letzte Zelle ausmacht, ist
   ein Zeilen-Suffix und keine Zelle.** Er wird vor dem Zellen-Zählen entfernt —
   **genau einer**, **nur** am Zeilenende, mit oder ohne nachfolgende Pipe. Damit
   trägt eine Tabellenzeile die Direktiven-Konvention, ohne den Reader zu brechen.

2. **Nur am Ende, nur ganzzellig.** Ein Kommentar **innerhalb** einer Zelle
   (`| a <!-- x --> | b |`) oder in der Zeilenmitte bleibt Zellinhalt. Der Schnitt
   ist damit lexikalisch eindeutig und rät nichts.

3. **Der Zellenzahl-Guard bleibt scharf.** Nach dem Entfernen des Suffixes gilt
   die Regel unverändert: abweichende Zellenzahl ⇒ Exit 2. Die Ausnahme ist genau
   eine — die eigene, dokumentierte Direktive —, nicht eine allgemeine Aufweichung.

4. **Kein Lastenheft-Change-Request.** Das Lastenheft definiert die Zellenzahl
   nicht; „was eine Zelle ist" ist Spezifikations-Sache (Rang 2, fortschreibbar).
   Der Entscheid **stellt die Lesbarkeit her**, die die Grammatik zusagt —
   SemVer-Patch, kein Minor.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Abschließenden Kommentar als Suffix abstreifen** (gewählt) | trifft genau die eigene Konvention; der Guard bleibt sonst scharf; lexikalisch eindeutig | eine Sonderregel im Splitter; weniger standard-treu als GFM |
| **GFM-Semantik: überzählige Zellen ignorieren** | standard-treu; jeder Renderer tut es; keine Sonderregel | macht **jede** verrutschte Zeile still — genau die Klasse, gegen die [ADR-0037](0037-trace-tabellenquellen-nullmengen-guard.md) den Guard gebaut hat; eine kaputte Tabelle verlöre lautlos Anforderungen |
| **Nichts ändern — Konsument entfernt den Marker aus der Tabellenzeile** | kein Code | verlagert das Problem auf jeden Konsumenten; der Marker ist **unsere** Empfehlung; in einer Zelle wäre er Inhalt, es gibt also keinen legalen Platz |
| Marker-Name hart erkennen (`d-check:ignore` statt „HTML-Kommentar") | maximal eng | koppelt den Splitter an eine Marker-Semantik, die er nicht kennt; fremde Trailing-Kommentare blieben grundlos rot |

**Fitness-Funktion:**

- Eine Tabellenzeile mit `<!-- d-check:ignore (…) -->` hinter der letzten Pipe wird
  gelesen wie dieselbe Zeile ohne Marker — für `trace.requirements.format: table`
  **und** `trace.cross-consistency`.
- Der Realdatenlauf gegen grid-gyms `spec/architecture.md` läuft durch, statt an
  Zeile 913 abzubrechen.
- Eine echt verrutschte Zeile (Zelle zu viel/zu wenig, **kein** Kommentar) bleibt
  Exit 2 — der Guard ist nicht aufgeweicht.
- Ein Kommentar **in** einer Zelle bleibt Zellinhalt (Titel/Kante unverändert).

## Konsequenzen

- **Positiv:** der Selbstwiderspruch zwischen Direktiven-Konvention und
  Tabellen-Reader ist weg; grid-gyms Rück-Sicht wird lesbar; die Regel gilt über den
  geteilten Reader für beide Konsumenten aus einem Fix.
- **Negativ / Kosten:** der Splitter trägt ein weiteres Stück Markdown-Wissen (nach
  Code-Span und Escape jetzt auch der Kommentar). Die Grenze „genau einer, nur am
  Ende" ist gesetzt — ein Konsument mit zwei Trailing-Kommentaren stößt an. Bewusst:
  die Alternative wäre Raten.
- **Verhaltensänderung für Bestandskonsumenten:** eine Quelle, die heute mit
  `Tabellenzeile N hat X statt Y Zellen` abbricht, läuft danach durch. Kein Lauf
  wird stiller: der Guard feuert weiter, nur nicht mehr auf die eigene Direktive.
- **Verworfen:** GFM-weite Aufweichung, Nichtstun, Marker-Namens-Kopplung (jeweils
  oben begründet).

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-07-17 | Proposed. Anlass: Realdatenlauf gegen grid-gym (der von [ADR-0038](0038-trace-cross-consistency.md) Entscheidung 7 geforderte Beleg) — v0.45.1 brach an einer realen `architecture.md`-Zeile ab, die d-checks eigenen Ignore-Marker trägt. Dritter Defekt, den erst die Realdaten zeigten; die ersten beiden lagen in den Mustern, dieser in der Grammatik. Umsetzender Slice slice-074. |
