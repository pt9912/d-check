# ADR-0039 — Link-transparente Range-/Enum-Fortsetzung im geteilten Range-Parser

**Status:** Accepted
**Datum:** 2026-07-17
**Autor:** pt9912
**Schärft:** [`DC-FA-COV-001.a`](../../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage), [`DC-FA-XREF-001.a`](../../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
**Bezug:** [`DC-FA-COV-001`](../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in), [`DC-FA-XREF-001`](../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in), [`DC-FA-ID-001`](../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids), [ADR-0035](0035-trace-coverage-quellen.md), [ADR-0038](0038-trace-cross-consistency.md), [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus).

## Kontext

Das Lastenheft verspricht die Range-Expansion **unqualifiziert**:
[`DC-FA-COV-001`](../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
sagt „`<FAM>-AAA..BBB` deckt alle …". Die Verengung steht allein in der
Spezifikation ([§`DC-FA-COV-001.a`](../../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)
Schritt 3): expandiert wird nur, wenn `..<Ziffern>` der Fundstelle **unmittelbar**
folgt. Der Code hält sich exakt daran — die Lücke ist im Vertrag, nicht in der
Implementierung.

Diese Verengung kollidiert **strukturell** mit dem eigenen Modul `ids`: wo
[`DC-FA-ID-001`](../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
Linkpflicht erzwingt, steht die Kennung als `` [`GG-UI-001`](…) `` — und die
Fortsetzung `..009` folgt dann eben **nicht** mehr unmittelbar, sondern hinter dem
Link-Suffix. Ein Repo, das beide Fähigkeiten nutzt (Linkpflicht **und**
Range-Notation), verliert die Expansion still. Das ist kein Randfall eines
Konsumenten, sondern die Regel für jeden mit `link-policy: always`.

**Der Defekt ist ausgeliefert.** Er trifft
[`DC-FA-COV-001`](../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
seit v0.41.0, nicht erst die Kreuzverweis-Fähigkeit: identische Quelle, nur die
Range verlinkt, ergibt statt 0 plötzlich 2 Waisen — **Falschbefunde**, die unter
`--require-complete` fälschlich gaten. Sichtbar wurde er erst über den
Kreuzverweis-Abgleich (Realdaten-Lauf grid-gym, 2026-07-17), weil dort beide
Sichten verglichen werden und die Asymmetrie auffällt.

Realdaten-Befund des Konsumenten (ausgezählt, nicht geschätzt): **40 von 40**
Fortsetzungen in den `Bezug`-Zellen haben exakt **eine** Form —
`` [`ID`](url) `` unmittelbar gefolgt von `..NNN` bzw. `/NNN(/NNN)*`. Null
Vorkommen ohne Code-Span, null mit einem Zeichen zwischen `)` und der Fortsetzung.

## Entscheidung

1. **Genau ein Link-Suffix darf die Fortsetzung unterbrechen.** Der geteilte
   Range-Parser überspringt nach einer `id-pattern`-Fundstelle **höchstens einmal**
   ein unmittelbar folgendes Markdown-Link-Suffix der Form `](…)` — inklusive des
   schließenden Code-Span-Backticks und der `]`, die zwischen ID und Suffix stehen
   — und prüft **dahinter** wieder **unmittelbar** auf `..<Ziffern>` bzw.
   `/<Ziffern>`. Sonst gilt die bisherige Regel unverändert.

2. **Kein weiteres Peeling.** Nicht übersprungen werden: Whitespace, Emphasis,
   mehrere aufeinanderfolgende Link-Suffixe, Text zwischen `)` und der
   Fortsetzung. Begründung: jede weitere Toleranz **rät**, was der Autor meinte,
   und macht aus einem deterministischen Parser eine Heuristik. Die Realdaten
   belegen, dass genau ein Suffix den gesamten Bestand deckt; mehr zu erlauben
   erkauft nichts und riskiert Falsch-Expansionen.

3. **Die ID darf im Code-Span des Linktexts sitzen** — das kann der Parser bereits
   (er findet `GG-UI-001` in `` [`GG-UI-001`](…) ``). Neu ist allein: nach `)`
   weiterscannen. Der Entscheid fügt **keine** neue ID-Erkennung hinzu.

4. **Ein Fix, zwei Konsumenten.** Der Parser ist geteilt
   ([ADR-0038](0038-trace-cross-consistency.md): Reader-Reuse statt zweitem
   Parser), daher wirkt der Entscheid zugleich auf
   [`DC-FA-COV-001.a`](../../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)
   und [`DC-FA-XREF-001.a`](../../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency).
   Das ist der Nutzen der geteilten Mechanik — und zugleich ihr Preis: derselbe
   Defekt saß in beiden.

5. **Kein Lastenheft-Change-Request.** Das Lastenheft verspricht die Expansion
   bereits unqualifiziert; der Entscheid **stellt die Zusage her**, statt sie zu
   ändern. Er ist eine Spezifikations-Schärfung (Rang 2, fortschreibbar) und ein
   **Defekt-Fix** — SemVer-Patch, kein Minor.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Genau ein Link-Suffix überspringen, Ziel klammer-balanciert über den kanonischen Reader** (gewählt) | deckt 100 % des belegten Bestands; bleibt deterministisch; ein Fix für beide Konsumenten; **eine** Link-Definition im Repo | ein neuer Sonderweg im Parser (Markdown-Wissen in einer bisher rein lexikalischen Funktion) |
| Ziel per Regex bis zur ersten `)` abgrenzen (**erste Fassung, verworfen**) | kein Rückgriff auf den Reader | **zweite** Link-Definition neben `rules.parseLinkAt`; bei einer Klammer im Ziel landet der URL-Rest im Range-Parser und expandiert Pfadsegmente als Enum — versteckt Waisen (still) statt sie zu melden (laut) |
| Beliebiges Peeling (Whitespace/Emphasis/Mehrfach-Links) | träfe auch ungesehene Schreibweisen | rät die Autor-Absicht; Falsch-Expansionen; aus dem Parser wird eine Heuristik |
| Zelle vor dem Scan Markdown-normalisieren (Links zu Klartext) | konzeptionell sauber | zweiter Markdown-Pfad neben dem Reader; ändert die Fundstellen-Offsets und damit `Datei:Zeile`; große Fläche für einen kleinen Defekt |
| Nichts tun — Konsument schreibt die Range in den Linktext (`` [`ID..009`](…) ``) | kein Code | verlagert die Last auf jeden Konsumenten; die Lastenheft-Zusage bleibt gebrochen; die Kollision mit der Linkpflicht bleibt strukturell bestehen |

**Fitness-Funktion:**

- Eine verlinkte Range expandiert wie die unverlinkte: identische Quelle, einmal
  `GG-UI-001..003`, einmal `` [`GG-UI-001`](…)..003 `` ⇒ **gleiches** Ergebnis.
- Eine Zelle **ohne** Range-Notation expandiert **nie** — auch nicht, wenn das
  Linkziel Klammern oder ziffern-artige Pfadsegmente trägt
  (`` [`GG-QA-001`](…/Rev(2)/002/003.md) ``). Diese Richtung ist die kritische:
  eine Falsch-Expansion **versteckt** Waisen (stilles Grün), eine fehlende meldet
  zu viele (lautes Rot).
- Der ausgelieferte `trace.coverage`-Falschbefund (2 Waisen statt 0) ist weg.
- Zwei Link-Suffixe hintereinander oder ein Zeichen zwischen `)` und `..` werden
  **nicht** expandiert (kein Raten) — negativ belegt.
- Unverlinkte Ranges, Enum-Notation und die Fail-closed-Fälle (`AAA>BBB`,
  Breiten-Mismatch) bleiben unverändert
  ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).

## Konsequenzen

- **Positiv:** eine ausgelieferte Falschbefund-Klasse verschwindet; die
  Lastenheft-Zusage gilt wieder; die strukturelle Kollision zwischen Range-Notation
  und Linkpflicht ist aufgelöst; beide Konsumenten des Parsers profitieren aus
  einem Fix.
- **Negativ / Kosten:** der Parser trägt jetzt ein Stück Markdown-Wissen (das
  Link-Suffix), war vorher rein lexikalisch. Die Grenze „genau eins" ist eine
  gesetzte Zahl — sie ist durch Realdaten gedeckt, aber ein Konsument mit einer
  anderen Schreibweise stößt erneut an. Das ist bewusst: die Alternative wäre
  Raten.
- **Verhaltensänderung für Bestandskonsumenten:** wer heute verlinkte Ranges in
  `trace.coverage` führt, sieht **weniger** Waisen — ein Lauf, der fälschlich rot
  war, wird grün. Kein Konsument verliert Deckung; kein Befund entsteht neu.
- **Verworfen:** beliebiges Peeling, Markdown-Vor-Normalisierung, Nichtstun
  (jeweils oben begründet).

## Geschichte

| 2026-07-17 | **Accepted** mit der Closure von slice-073. Die Entscheidung ist umgesetzt und als v0.45.1 ausgeliefert; der unabhängige R2-Review (ACCEPT-WITH-NITS) hat alle R1-Befunde durch eigene Messung geschlossen, der Rest R2-F-1 (Cross-Test-Lücke) ist in `b8c503a` behoben. Kein Inhalt geändert — nur der Status folgt der Auslieferung nach (Statusdrift vermeiden, vgl. ADR-0037 bei slice-070). |
| Datum | Ereignis |
|---|---|
| 2026-07-17 | Ziel-Abgrenzung auf den **klammer-balancierten** kanonischen Reader (`rules.LinkSuffixEnd`) gezogen; die erste Fassung grenzte per Regex bis zur ersten `)` ab und war damit eine zweite Link-Definition. Anlass: nachgeholter unabhängiger Review zu slice-073 (R1-F-1, HIGH) — bei einer Klammer im Ziel landete der URL-Rest im Range-Parser, `/002/003` expandierte als Enum, und eine Zelle ohne jede Range-Notation versteckte Waisen (Exit 1 → Exit 0). Fitness-Funktion um die kritische Richtung ergänzt. Status weiterhin `Proposed`. |
| 2026-07-17 | Proposed. Anlass: Realdaten-Lauf des Konsumenten grid-gym gegen v0.44.0 (Defekt 2 von zweien) — reproduziert und dabei als **ausgelieferter** `trace.coverage`-Defekt seit v0.41.0 erkannt, den der Konsumenten-Report nicht sehen konnte. Umsetzender Slice slice-073. |
