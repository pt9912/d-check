# `make completeness-check` — meldet Anforderungen ohne referenzierenden Slice am Closure-Punkt

## Vertrag

Requirements-Completeness **via in-Produkt-Flag** `--trace --require-complete`:
mindestens eine **Requirements-Waise** — eine Anforderung, die kein Slice
referenziert — ergibt Exit 1, mit Anzahl-Meldung und den sichtbaren
`WAISE`-Zeilen der `--trace`-Tabelle.

Dieselbe verteilte Mechanik wie `doc-complete` und die übrigen Konsumenten —
**keine Skript-Kopie mehr**
([ADR-0026](../../docs/plan/adr/0026-completeness-in-product-gate.md) löst die
Skript-Mechanik ab). Netzlos, read-only.

## Grenze — was das Grün nicht abdeckt

1. **Gemessen ist die Referenz, nicht die Erfüllung** — dass ein Slice eine
   Anforderung nennt, sagt nicht, dass er sie einlöst.
2. **Bewusst nicht in `gates`/`ci`** — Greenfield erlaubt **transiente**
   Waisen im inneren Loop; ein Gate dort machte das Arbeiten an einer neuen
   Anforderung unmöglich. Die Disziplin liegt am Closure-Punkt.

## Bindung

**Closure-Bindepunkt** — in `make fullbuild`, bewusst **nicht** in `gates`/`ci`
und auch nicht in der Release-Pipeline.
[ADR-0026](../../docs/plan/adr/0026-completeness-in-product-gate.md) ·
[`DC-FA-CLI-011`](../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code) ·
[`DC-FA-CLI-009`](../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
