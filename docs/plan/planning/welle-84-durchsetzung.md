# Welle welle-84-durchsetzung: Neun Hard Rules — und ein Sensor, den niemand aufruft

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** und liegt **flach**
unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/`
(neben ihre `welle-84-results.md`). Der Zustand ist die Verzeichnis-Position —
kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug. Ob die Welle mit einem Release
schließt, entscheidet der Zensus — durchgesetzt wird zuerst mit dem, was schon
da ist.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

[`modul-09-implementierung.md` §AGENTS.md-Regeln](../../../.harness/baseline/v5.11.0/regelwerk/modul-09-implementierung.md#agentsmd-regeln-modul-9)
sagt: *„Jede Hard Rule liegt in **zwei** Quadranten: inferential feedforward
(steht in AGENTS.md) + computational feedback (Fitness Function/Linter-Gate).
Hard Rule nur in einem Quadranten ist **halb durchgesetzt**."*

[`AGENTS.md`](../../../AGENTS.md) §3 trägt **neun** Hard Rules, dazu die
Botschafts-Regel in §5. Für **drei** davon ist die Einseitigkeit belegt — sie
sind in welle-83 zugezogen und tragen kein Gate. Für die übrigen ist sie
**weder belegt noch widerlegt**: niemand hat je Regel für Regel gefragt, welcher
Gate-Lauf sie trägt.

**Der Anlass ist ein Fund, der zeigt, dass die Frage nötig ist.** Beim Sichten
der Werkzeuge des Schwester-Repos hat sich herausgestellt, dass
[`tools/harness/fetch-baseline-cache.sh`](../../../tools/harness/fetch-baseline-cache.sh)
drei Fähigkeiten trägt — Integritätsprüfung des vendorten Bestands,
Currency-Audit gegen die Release-Liste, Content-Drift am gepinnten Tag — und
dass **kein `make`-Target, kein Workflow und kein Hook es aufruft**. Ein
gebauter, nie eingesteckter Sensor ist die reinste Form von *halb durchgesetzt*:
die Feedback-Hälfte existiert und wirkt trotzdem nicht.

Ziel der Welle ist deshalb **eine Antwort je Regel**, nicht ein Gate mehr:
entweder ein benannter Gate-Lauf, der sie trägt — oder die ausdrückliche
Ausweisung, dass sie einseitig bleibt.

## 2. Trigger (Welle startet)

- [welle-83](done/welle-83-baseline-v5110-migration.md) ist geschlossen, ihre
  Ergebnisnotiz liegt in `done/`.
- `in-progress/` trägt keinen Slice.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen.

- **Der Zensus liegt vor:** je Hard Rule in [`AGENTS.md`](../../../AGENTS.md) §3
  **und** der Botschafts-Regel in §5 eine Antwort **mit Beleg** — der tragende
  Gate-Lauf ist benannt, oder die Regel ist als **einseitig ausgewiesen**.
  **Keine Regel ohne Zeile.** Das ist das *Mehr*: eine Aussage über eine ganze
  Datei, die keine einzelne Slice-DoD belegen kann.
- **Was der Zensus als baubar ausweist, ist entweder in dieser Welle gebaut oder
  als Folge-Welle benannt** — nicht stillschweigend offen.
- `make fullbuild` grün (Exit explizit).
- Ergebnisnotiz `welle-84-results.md` mit Register-Lese-Schritt.

## 4. Slices in dieser Welle

| Slice | Rolle |
|---|---|
| [slice-132](in-progress/slice-132-hard-rule-zensus.md) | **Der Zensus:** je Hard Rule eine Antwort mit Beleg; weist die einseitigen aus und **schneidet** den Rest |
| [slice-133](done/slice-133-baseline-sensor-verdrahten.md) | **Den verwaisten Sensor anschließen:** `--verify` als Gate, `--check-latest` als Target plus Nachtlauf |

slice-133 hängt **nicht** am Zensus — er behebt einen bereits belegten Fund und
läuft unabhängig. Was der Zensus darüber hinaus schneidet, kommt als weitere
Slices in diese Welle oder als Folge-Welle.

## 5. Abhängigkeiten

- **Blockiert:** eine mögliche **Produkt-Welle** (Freshness als Modul). Ihr
  Zuschnitt kommt aus dem Zensus und wird hier nicht vorweggenommen — die
  entscheidende Messung liegt vor: d-check scannt **nur Markdown** (467 geprüfte
  Dateien = 467 `.md`-Dateien), und die zwei ungewachten Toolchain-Pins stehen
  im `Dockerfile`, also außerhalb der gescannten Menge.
- **Wird blockiert von:** nichts.

## 6. Out-of-Scope für diese Welle

- **Kein Lastenheft-Delta und keine neuen `DC-FA-*`.** Ein Produkt-Delta hätte
  eine eigene Closure-Bedingung — ein Release —, und die gehört nicht in eine
  Durchsetzungs-Welle.
- **Kein Release.**
- **Keine Heuristik-Wächter.** Wo kein baubares Gate existiert, wird die
  Einseitigkeit **ausgewiesen**, nicht mit einem Wortlisten-Lint übertüncht. Der
  Kurs hat genau diese Forderung an unserem eigenen Konsumenten-CR abgelehnt:
  *ohne Baubarkeit wäre sie ein behauptetes Gate.*
- **Kein Mutations-Sensor.** Die Einführungs-Hälfte (*„Bewusstes Brechen"*)
  steht bereits gerankt in
  [`modul-13-quality-gates.md`](../../../.harness/baseline/v5.11.0/regelwerk/modul-13-quality-gates.md#adr-zur-fitness-function);
  die **Haltbarkeits**-Frage wartet auf einen belegten Vorfall — das
  [Beobachtungs-Register](observations.md) führt dazu **null von zwölf**
  Einträgen, und `modul-13` §Guard-Härtung verlangt Beobachtung statt
  Bedrohungsmodell.
