# `make gate-consistency` — hält Doku und Makefile über die Gate-Targets in beiden Richtungen deckungsgleich

## Vertrag

Meta-Gate gegen die häufigste Form der Harness-Lüge: ein dokumentiertes Gate,
das es nicht gibt. Der cross-repo-Kern — in [`AGENTS.md`](../../AGENTS.md) §4
und in der Sensors-Tabelle dokumentierte `make X` ↔ Makefile-Regeln, **beide
Richtungen** — läuft via Modul `targets` (Image, dogfood). Der

## Grenze — was das Grün nicht abdeckt

1. **Geprüft ist die Deklarations-Deckung, nicht der Anspruch** — dass ein
   Target existiert und dokumentiert ist, sagt nichts darüber, ob es als Gate
   *behauptet* wird oder als Werkzeug geführt gehört.
2. **Der frühere Skript-Rest ist voll abgelöst** — die
   [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Modulliste
   der [`.d-check.yml`](../../.d-check.yml) prüft ein getippter
   `configyaml.Decode`-Go-Test unter `make test`
   ([ADR-0032](../../docs/plan/adr/0032-gate-consistency-tombstone.md)); dieses
   Target trägt sie nicht mehr.

## Bindung

Bestandteil von `make gates`.
[ADR-0031](../../docs/plan/adr/0031-targets-deklarations-konsistenz-modul.md) ·
[`DC-FA-TGT-001`](../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
