# `make doc-check` — prüft die gesamte Repo-Doku mit d-check selbst (Dogfooding)

## Vertrag

**Elf Module** über die gesamte Repo-Doku: Links, **Anker**,
**Kennungs-Linkpflicht**, **Referenzmatrix**, **Inline-Code-Pfade**,
**Abschnitts-Invarianten** (`structure`), **Kennungen in Diagramm-Fences**
(`diagrams`), **Markdown-Span-Artefakte** (`spans`), **Host-Pfade**
(`hostpaths`), **Versions-Pins** (`versions`) und **wortgleiche Zitate gegen
ihre Quelle** (`citations`).

**Dogfooding-Selbstkonfiguration:** Runtime-Image, read-only-Mount,
`--network none`, [`.d-check.yml`](../../.d-check.yml). Der Lauf ist zugleich
die automatisierte Messmethode für
[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
— netzloser Lauf aller Module außer `external`.

## Grenze — was das Grün nicht abdeckt

1. **`citations` ist fail-closed und läuft im inneren Loop** — eine
   strukturell unbrauchbare Direktive nimmt den ganzen Lauf mit, auch den
   `pre-commit`-Hook. Tragbar, seit die Direktive nur außerhalb von
   Inline-Code zählt ([ADR-0060](../../docs/plan/adr/0060-citations-marker-scan-geteilte-prosa-antwort.md));
   der fail-closed-Rest ist ein Autoren-Fehler an einer **freien** Direktive.
2. **Eine planmäßige Rot-Quelle ist der Baseline-Bump** — er verschiebt die
   Zeilenspannen; das Neu-Ankern liegt in der Bump-Prozedur
   ([`MR-051`](../conventions/MR-051-cite-spannen-beim-bump.md)).
3. **Ein Merge trägt eine fremde Direktive am Hook vorbei**, und
   `d-check:ignore` greift dort **nicht** — grob wirken nur `scan.ignore` und
   `citations.scope`. Der Modul-Scope nimmt seit
   [`MR-054`](../conventions/MR-054-vorpruefungen-belegen-ihre-regel.md) die
   drei eingefrorenen Verzeichnisse aus.

**Wie groß der Ausschnitt ist, sagt das Kommando:** `make doc-check` nennt die
geprüfte Dateizahl in seiner Schluss-Zeile.

## Bindung

Bestandteil von `make gates`; zusätzlich im `pre-commit`-Hook als Doku-Guard.
[`DC-FA-LINK-001`](../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links) ·
[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) ·
[`MR-007`](../conventions/MR-007-aufloesung-mr-003.md)
