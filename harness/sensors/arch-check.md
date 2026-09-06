# `make arch-check` — hält die Import-Regeln des Hexagon-Schnitts über das Schwester-Tool a-check

## Vertrag

Import-Regeln R1–R5 des Hexagon-Schnitts plus R6 (Kern-Paket-Richtung
`model` ← `rules` ← `app`), **via digest-gepinntes a-check-Image**: ein
include-bares `a-check.mk` plus [`.a-check.yml`](../../.a-check.yml), Lauf mit
`--network none` und read-only-Mount.

**Test-Dateien sind via `exclude` ausgenommen.** Die Ablösung des früheren
Skripts folgt [ADR-0029](../../docs/plan/adr/0029-arch-check-via-a-check.md);
womit sie gedeckt ist, steht dort.

## Grenze — was das Grün nicht abdeckt

1. **Der Gegenstand sind Import-Kanten** — die sechs Regeln R1–R6, nicht die
   Rollen-Treue der Schichten: dass ein Paket importieren *darf*, sagt nicht, dass es die
   richtige Rolle spielt. Permanent — das ist Review-Territorium.

## Bindung

Bestandteil von `make gates`. Netzlos, read-only.
[ADR-0005](../../docs/plan/adr/0005-modul-layout-hexagon-ordner.md) ·
[ADR-0012](../../docs/plan/adr/0012-kern-paketschnitt-model-rules-app.md) ·
[ADR-0029](../../docs/plan/adr/0029-arch-check-via-a-check.md) ·
[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
