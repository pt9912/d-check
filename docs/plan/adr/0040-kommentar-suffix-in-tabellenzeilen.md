# ADR-0040 — Direktiven-Zelle in Tabellenzeilen: Header GFM-streng, Datenzeilen verengt nachsichtig

**Status:** Proposed — **Entscheid ausgesetzt, Implementierung zurückgenommen**
(2026-07-17). Die unten getroffene Entscheidung ist **falsifiziert**: Review R3
hat an ihr einen Stilles-Grün-Pfad belegt (F-1), und ihre zentrale Zusage („die
Toleranz greift **nie** über eine Tabellengrenze") ebenso widerlegt wie die
Mutations-Zusage aus Entscheidung 3 (F-2). Der Kontext (§Kontext) gilt
unverändert — der Defekt besteht. Was fehlt, ist eine tragende Regel. Bis dahin
ist der Reader byte-identisch v0.45.1; Details in
[slice-074](../planning/open/slice-074-kommentar-suffix-tabellenzeilen.md) §2.
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

**Vorbemerkung zur verworfenen Prämisse (2026-07-17, nach Review):** Die erste
Fassung dieser ADR nannte den abschließenden Kommentar ein „Suffix, keine Zelle".
Das ist **falsch**. In GFM **ist** er eine Zelle; GFM ignoriert bei **Body**-Zeilen
lediglich die überzähligen und verlangt beim **Header** Gleichstand mit der
Trennzeile. `| a | b | <!-- x -->` mit einer N+1-Trennzeile ist damit die **einzige
renderbare** Header-Form. Ein echter GFM-Parser gäbe N+1 Zellen zurück und ließe
die Frage, wo die Direktive legal wohnt, exakt so offen.

Die eigentliche Klasse ist deshalb **nicht** „der lexikalische Splitter liegt
daneben", sondern: d-check will zugleich **GFM-nachsichtig** sein (damit die
Direktive rendert) und **fail-closed** (damit nichts still verschwindet) — und
jeder Patch wählte für einen anderen Codepfad einen anderen Punkt auf dieser
Achse. Es gab keine durchgehende Regel.

1. **Die Regel folgt GFM, getrennt nach Zeilenart — das ist die durchgehende Regel.**
   - **Header: GFM-streng.** Header-Zellenzahl muss der Trennzeile entsprechen
     (`tableHeaderAt`, unverändert). Trägt ein Header die Direktive, ist sie eine
     **Spalte**: N+1 gegen N+1 wird regulär erkannt, und die Extra-Spalte bindet an
     keine Rolle. Der Header-Pfad wird **nicht angefasst**.
   - **Datenzeilen: GFM-nachsichtig, aber verengt.** GFM ignoriert **jede**
     überzählige Zelle; d-check toleriert **genau eine** — und nur, wenn sie
     ganzzellig ein HTML-Kommentar ist. Alles andere bleibt fail-closed
     ([ADR-0037](0037-trace-tabellenquellen-nullmengen-guard.md): eine still
     abreißende Tabelle verlöre Anforderungen lautlos).

2. **Die Regel wohnt dort, wo der Header-Kontext bekannt ist** (`cellCountOK` in
   der Zeilen-Schleife), **nicht** im Splitter. Der Splitter kennt den Unterschied
   zwischen Header und Datenzeile nicht; eine Regel dort ist eine **Body-Regel am
   Header** — genau der Defekt der ersten Fassung: der Header fiel unter die
   Trennzeilen-Breite, die Tabelle wurde **wortlos übersprungen**, ihre
   Anforderungen verschwanden (Exit 1 ⇒ Exit 0). `splitPipeTableLine` bleibt ein
   **reiner** Splitter und entfernt nichts.

3. **Nichts wird gestrippt, nur toleriert** — keine Zeilenart verliert Zellen, und
   die Nachsicht ist auf eine Kommentar-Zelle in einer Datenzeile begrenzt.
   **Die Toleranz endet am nächsten Tabellen-Header.** Eine Zeile, der eine
   passende Trennzeile folgt, ist der Header einer **neuen** Tabelle und keine
   tolerierbare Datenzeile der laufenden. Ohne diese Grenze fräße die Toleranz
   einen Direktiven-Header der Folgetabelle (N+1 mit Kommentar-Zelle sieht wie eine
   tolerierbare Datenzeile aus) und dessen Anforderungen verschwänden **lautlos**
   (Review R2-F-1: Exit 2 ⇒ Exit 0).
   *Zur Reichweite, ehrlich:* Die erste Neufassung behauptete, der stille
   Übersprung sei „strukturell unmöglich". Das war eine **Universal-Zusage ohne
   Beweis** — R2 hat sie an genau diesem Pfad falsifiziert. Die Zusage lautet jetzt
   enger und prüfbar: **keine Zeilenart verliert Zellen, und die Toleranz greift
   nie über eine Tabellengrenze**; beide Grenzen sind per Mutation gepinnt.

4. **Kein Lastenheft-Change-Request.** Das Lastenheft definiert die Zellenzahl
   nicht; „was eine Zelle ist" ist Spezifikations-Sache (Rang 2, fortschreibbar).
   SemVer-**Patch**: gegenüber v0.45.1 wird **keine** Zeile anders gelesen, die
   heute gelesen wird — es kommt nur die N+1-Datenzeile mit Direktive hinzu.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Überzählige Kommentar-Zelle in Datenzeilen tolerieren, Header unberührt** (gewählt) | folgt GFM getrennt nach Zeilenart — eine durchgehende Regel statt drei Punkten auf der Achse; stiller Übersprung strukturell unmöglich; Guard bleibt sonst scharf | eine Sonderregel in der Zeilen-Schleife; die Direktive braucht im Header eine N+1-Trennzeile |
| Kommentar im Splitter abstreifen (**erste Fassung, verworfen**) | eine Stelle für alle Leser | der Splitter kennt Header ≠ Datenzeile nicht ⇒ **Body-Regel am Header**: Header fällt unter die Trennzeilen-Breite, Tabelle wird wortlos übersprungen, echte Waisen verschwinden (Exit 1 ⇒ Exit 0). Zudem: eine Datenzeile mit legitimer Kommentar-**Spalte** (N Zellen) wurde fälschlich Exit 2 |
| **GFM-Semantik voll: jede überzählige Body-Zelle ignorieren** | standard-treu; keine Sonderregel | macht **jede** verrutschte Zeile still — genau die Klasse, gegen die [ADR-0037](0037-trace-tabellenquellen-nullmengen-guard.md) den Guard baute |
| Echter Markdown-/GFM-Parser statt lexikalischem Splitter | konform | **schließt die Klasse nicht**: ein GFM-Parser gäbe N+1 Zellen zurück und ließe die Frage, wo die Direktive legal wohnt, exakt so offen. Große ADR, kein Gewinn an dieser Achse |
| Nichts ändern — Konsument entfernt den Marker aus der Tabellenzeile | kein Code | verlagert das Problem auf jeden Konsumenten; der Marker ist **unsere** Empfehlung; in einer Zelle wäre er Inhalt |

**Fitness-Funktion:**

- Eine **Datenzeile** mit `<!-- d-check:ignore (…) -->` hinter der letzten Pipe wird
  gelesen wie dieselbe Zeile ohne Marker — für `trace.requirements.format: table`
  **und** `trace.cross-consistency`.
- Eine Tabelle, deren **Header** die Direktive trägt (N+1 gegen N+1-Trennzeile),
  wird **erkannt**; ihre Anforderungen bleiben sichtbar. Diese Richtung ist die
  kritische: ein übersprungener Header **versteckt** Waisen (stilles Grün).
- Eine Datenzeile mit legitimer Kommentar-**Spalte** (N Zellen) bleibt unverändert
  lesbar — keine Regression gegenüber v0.45.1.
- Eine echt verrutschte Zeile (Zelle zu viel/zu wenig, **kein** Kommentar) bleibt
  Exit 2 — der Guard ist nicht aufgeweicht.
- Folgt der laufenden Tabelle **ohne Leerzeile** eine neue mit Direktiven-Header,
  wird deren Header **nicht** als Datenzeile toleriert; das Verhalten ist
  byte-identisch zu v0.45.1 (laut, nicht still).
- Der Realdatenlauf gegen grid-gyms `spec/architecture.md` läuft durch, statt an
  Zeile 913 (einer **Daten**zeile) abzubrechen.

## Konsequenzen

- **Positiv:** der Selbstwiderspruch zwischen Direktiven-Konvention und
  Tabellen-Reader ist weg; grid-gyms Rück-Sicht wird lesbar; die Regel gilt über den
  geteilten Reader für beide Konsumenten aus einem Fix.
- **Negativ / Kosten:** die Zeilen-Schleife trägt eine Nachsichts-Regel; der
  Splitter bleibt rein. Die Grenze „genau eine, nur ganzzellig Kommentar" ist
  gesetzt — zwei überzählige Zellen stoßen an. Bewusst: die Alternative wäre Raten.
  Wer die Direktive im **Header** nutzt, braucht dort eine N+1-Trennzeile — die
  GFM-Form; das ist keine d-check-Eigenheit.
- **Verhaltensänderung für Bestandskonsumenten:** eine Quelle, die heute mit
  `Tabellenzeile N hat X statt Y Zellen` an einer Direktiven-**Datenzeile**
  abbricht, läuft danach durch. Sonst ändert sich nichts — verifiziert gegen
  v0.45.1 für den Header-Fall und die legitime Kommentar-Spalte. **Kein Lauf wird
  stiller**; diese Zusage trug die erste Fassung nicht (sie war der Defekt) und ist
  jetzt durch die Zeilenart-Trennung strukturell gedeckt.
- **Verworfen:** Strippen im Splitter (die erste Fassung — sie erzeugte den stillen
  Übersprung), GFM-weite Aufweichung, echter GFM-Parser, Nichtstun,
  Marker-Namens-Kopplung (jeweils oben begründet).

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-07-17 | **Entscheid ausgesetzt, Implementierung zurückgenommen** (`05e1889`, `806051f`); slice-074 `in-progress/` → `open/`. Anlass: Review R3 (BLOCK) — F-1 (HIGH) belegt gegen das **ausgelieferte** v0.45.1 einen Stilles-Grün-Pfad in **beiden** Konsumenten: eine tolerierte Direktiven-**Datenzeile** entfernt den Wiederaufsetz-Punkt des Header-Scans und verschluckt die **gesamte** Folgetabelle (Exit 1 ⇒ Exit 0). Damit sind die in der vorigen Zeile **neu** gegebene Zusage („die Toleranz greift nie über eine Tabellengrenze") und die SemVer-Begründung („keine Zeile wird anders gelesen") falsifiziert. F-2 (MEDIUM): die Zusage „beide Grenzen sind per Mutation gepinnt" war für beide Grenzen aus `670ebaf` **falsch** — Rückdrehen und Panic-Guard-Entfernen lassen die Suite grün. **Fünfte Wiederholung derselben Klasse** in dieser Code-Region; ein sechster Vorschlag (`isNewTableHeader` unbedingt) wurde vor dem Code an Fixture `fx-t` widerlegt. Konsequenz: nicht der sechste Anlauf, sondern Rücknahme — ein Gate-Werkzeug trägt keinen bekannten stillen Pfad auf `main`, auch ungetaggt nicht. Nachträglich belegt (`fx-s`, gegen v0.45.1): der stille Pfad ist **älter** als dieser Slice und braucht **keinen** Marker ⇒ eigener Defekt. Spike goldmark v1.8.4 (522 reale Dateien): ein echter GFM-Parser schließt diese Klasse **nicht** — er stimmt auf `fx-s`/`fx-p` exakt mit dem heutigen Reader überein; die Achse ist Policy, nicht Grammatik. Die Alternativen-Tabelle bleibt insoweit gültig, ihre Begründung „schließt die Klasse nicht" ist jetzt **gemessen** statt behauptet |
| 2026-07-17 | Toleranz um die **Tabellengrenze** verengt (`isNewTableHeader`) und die Reichweiten-Zusage korrigiert. Anlass: Review R2 (BLOCK) — die Neufassung behauptete „stiller Übersprung strukturell unmöglich"; R2 falsifizierte das: die Toleranz fraß den Direktiven-Header einer unmittelbar folgenden Tabelle, deren Anforderungen verschwanden lautlos (Exit 2 ⇒ Exit 0). Der angebotene Ausweg „ehrliche Reichweiten-Angabe statt Code-Verengung" wurde **verworfen**: eine Zusage zu entschärfen, damit ein stiller Verlust hineinpasst, ist die Bewegung, gegen die dieses Werkzeug gebaut ist. Verhalten jetzt byte-identisch zu v0.45.1 an allen geprüften Achsen. |
| 2026-07-17 | **Prämisse verworfen und Entscheid neu gefasst**, Status weiterhin `Proposed`. Anlass: unabhängiger Review VOR dem Release (BLOCK, HIGH). Die erste Fassung („Kommentar ist ein Suffix, keine Zelle") war sachlich falsch — in GFM ist er eine Zelle, N+1 gegen N+1-Trennzeile ist die einzige renderbare Header-Form. Sie strippte deshalb im Splitter, wandte damit eine Body-Regel auf den Header an, ließ dessen Zellenzahl unter die Trennzeile fallen und übersprang die Tabelle **wortlos**: echte Waisen verschwanden, Exit 1 wurde Exit 0 — der eigene Satz „Kein Lauf wird stiller" war widerlegt. Neu: Header GFM-streng (unberührt), Datenzeilen verengt nachsichtig, Regel dort wo der Header-Kontext bekannt ist. Kein Strippen, nur Tolerieren — der stille Übersprung ist damit strukturell unmöglich statt behoben. |
| 2026-07-17 | Proposed. Anlass: Realdatenlauf gegen grid-gym (der von [ADR-0038](0038-trace-cross-consistency.md) Entscheidung 7 geforderte Beleg) — v0.45.1 brach an einer realen `architecture.md`-Zeile ab, die d-checks eigenen Ignore-Marker trägt. Dritter Defekt, den erst die Realdaten zeigten; die ersten beiden lagen in den Mustern, dieser in der Grammatik. Umsetzender Slice slice-074. |
