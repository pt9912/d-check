# ADR-0041 — Komma-Kurzform ist fail-closed, nicht still und nicht geraten

**Status:** Proposed
**Datum:** 2026-07-17
**Autor:** pt9912
**Schärft:** [`DC-FA-COV-001.a`](../../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)
**Bezug:** [`DC-FA-COV-001`](../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in), [`DC-FA-XREF-001`](../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in), [ADR-0035](0035-trace-coverage-quellen.md), [ADR-0039](0039-link-transparente-range-fortsetzung.md), [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus).

## Kontext

Zugesagt sind zwei Notationen: die Range `<FAM>-AAA..BBB` und die Aufzählung
`<FAM>-AAA/BBB/CCC`. Ein Konsument schreibt real aber auch `GG-SCN-001, 007` —
eine **Komma-Kurzform**. Sie war **nie** Vertrag; d-check las bisher nur
`GG-SCN-001` und ließ `007` **still** fallen.

Der Defekt ist nicht die fehlende Unterstützung — es ist der **stille Drop**.
Belegt (Konsumenten-Report grid-gym gegen v0.45.1, nachgemessen):

```
GG-SCN-001, 007   →  2 Anforderungen, 1 Waise    ← 007 fällt raus
GG-SCN-001/007    →  2 Anforderungen, 0 Waisen
GG-SCN-001..007   →  2 Anforderungen, 0 Waisen
```

Real: `§27.1` führt `GG-SCN-007/008`, deren Coverage-Spalte bleibt leer — das
Mapping wird nicht gezählt. Die Wirkung ist in `trace.coverage` **laut** (eine
falsche Waise, Gate rot) — die sichere Richtung; still wäre sie nur in der
Rück-Sicht unter `mode: superset`. Laut oder nicht: der Autor bekommt kein Signal,
dass seine Notation nicht gelesen wurde.

## Entscheidung

1. **Die Gestalt triggert, der Inhalt entscheidet: Exit 2.** Folgt einer
   `id-pattern`-Fundstelle ein Komma und **unmittelbar darauf Ziffern**, ist das
   eine Aufzählungs-Gestalt ohne zugesagte Notation ⇒ **Exit 2** mit Hinweis auf
   `/` und `..`. Exakt die Logik von `AAA>BBB`: die Syntax hat getriggert, der
   Inhalt ist ungültig — also rot, nicht still.

2. **Komma vor einer vollständigen Kennung ist unberührt.**
   `<FAM>-AAA, <FAM>-BBB` ist keine Kurzform; beide Kennungen werden regulär
   gefunden. Die Regel feuert **nur** bei Komma **plus Ziffern**. Das ist
   konstitutiv: komma-getrennte volle Kennungen sind in realen Design-Zellen der
   Normalfall (`GG-AR-COMP-CORE, GG-AR-COMP-DOMAIN`), und eine Regel, die sie
   bricht, wäre unbrauchbar.

3. **Kein Ausbau der Notation.** Komma-Enum zu unterstützen hieße raten: d-check
   kann eine gemeinte Kurzform nicht von einer Zahl im Fließtext unterscheiden
   (`GG-QA-001, 007 Sekunden`). Jede expandierte ID wird zwar gegen `id-pattern`
   validiert, aber eine breiten-gleiche Prosa-Zahl passiert diesen Filter.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Fail-closed auf der Gestalt** (gewählt) | rät nichts; der Autor bekommt ein Signal statt eines stillen Drops; dieselbe Logik wie `AAA>BBB` | eine Prosa-Zahl hinter einer Kennung (`GG-QA-001, 2026`) wird rot — laut und in Sekunden behebbar, aber ein neuer Falsch-Rot-Fall |
| Komma-Enum unterstützen | träfe die Autor-Absicht im Normalfall | **rät**: `GG-QA-001, 007 Sekunden` expandierte fälschlich; die id-pattern-Validierung fängt breiten-gleiche Prosa-Zahlen nicht; erweitert die Notations-Fläche für einen Konsumenten-Sonderweg |
| Status quo (stiller Drop) | kein Code | der Autor erfährt nie, dass seine Notation nicht gelesen wurde — die schlechteste der drei Optionen (so auch der Konsument selbst) |

**Fitness-Funktion:**

- `GG-SCN-001, 007` ⇒ Exit 2 mit Hinweis auf `/` und `..` — nicht 1 Waise, nicht 2 Deckungen.
- `GG-AR-COMP-CORE, GG-AR-COMP-DOMAIN` ⇒ beide Kennungen gelesen, **kein** Fehler.
- `GG-QA-001, siehe X` ⇒ kein Fehler (kein Ziffern-Suffix).
- `..`/`/`-Notation und die Fail-closed-Fälle unverändert
  ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).

## Konsequenzen

- **Positiv:** eine dritte Klasse „Notation kollidiert mit Realität" wird laut statt
  still; der Autor bekommt einen umsetzbaren Hinweis; die Notations-Fläche wächst
  **nicht**.
- **Negativ / Kosten:** ein neuer Falsch-Rot-Fall (Prosa-Zahl hinter einer Kennung).
  Bewusst in Kauf genommen: laut und behebbar sticht still und unbemerkt. Wer die
  Kurzform nutzt, muss seine Quellen anfassen — **das ist der Punkt**, nicht der
  Preis.
- **Verhaltensänderung für Bestandskonsumenten:** eine Quelle mit Komma-Kurzform
  oder Prosa-Zahl hinter einer Kennung läuft künftig auf Exit 2. Das ist ein
  **Vertrags-Zuwachs** (neues Akzeptanzkriterium), daher Lastenheft-CR 0.46.0 und
  SemVer-**Minor** — anders als die Defekt-Fixes
  [ADR-0039](0039-link-transparente-range-fortsetzung.md)/[ADR-0040](0040-kommentar-suffix-in-tabellenzeilen.md),
  die nur zugesagte Lesbarkeit herstellten.
- **Verworfen:** Komma-Enum-Ausbau, Status quo (jeweils oben begründet).

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-07-17 | Proposed. Anlass: Konsumenten-Report grid-gym gegen v0.45.1 — `GG-SCN-001, 007` ließ das produktiv verdrahtete `trace.coverage` das Mapping nicht zählen. Der Konsument benannte die Klasse selbst („stiller Drop ist die schlechteste der drei Optionen") und stellte fest, dass die Kurzform out of spec ist. Vierter Realdaten-Befund; Nutzer-Entscheid: fail-closed statt Notations-Ausbau. Umsetzender Slice slice-075. |
