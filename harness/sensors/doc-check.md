# `make doc-check` — prüft die gesamte Repo-Doku mit d-check selbst (Dogfooding)

## Vertrag

**Die Module der Liste `modules`** der `.d-check.yml` über die gesamte Repo-Doku, darunter Links, **Anker**,
**Kennungs-Linkpflicht**, **Referenzmatrix**, **Inline-Code-Pfade**,
**Abschnitts-Invarianten** (`structure`), **Kennungen in Diagramm-Fences**
(`diagrams`), **Markdown-Span-Artefakte** (`spans`), **Host-Pfade**
(`hostpaths`), **Versions-Pins** (`versions`), **wortgleiche Zitate gegen
ihre Quelle** (`citations`) und **die Gleichheit eines Wertes an mehreren
Stellen** (`pinsync`).

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

4. **Der Prüfbereich ist über drei Ventile verkleinert, und das ist eine
   deklarierte Gate-Senkung** — `ignore-refs` trägt **45** Einträge, und
   **248** `d-check:ignore`-Marker stehen über die Doku verteilt. Beide sind
   bewusste Akte mit Begründung an Ort und Stelle
   ([`MR-069`](../conventions.md#mr-069) nennt das Ventil selbst so), aber ein
   grüner Lauf sagt „nichts außerhalb dieser Ausnahmen". Ein drittes Ventil
   nimmt ganze Pfade aus dem Scan: `scan.ignore` in der
   [`.d-check.yml`](../../.d-check.yml) — die vendorte Baseline und ihr Cache
   (Fremdinhalt) und die Fixtures der [Black-Box-Probe](blackbox-probe.md),
   die absichtlich kaputte Links tragen.
5. **`pinsync` bewacht zwei Paare, nicht alle Pins** — `GO_VERSION` und
   `GOLANGCI_LINT_VERSION` in `Dockerfile` und `Makefile`
   ([`DC-FA-PSYNC-001`](../../spec/lastenheft.md#dc-fa-psync-001--gleichheit-eines-wertes-an-mehreren-stellen-modul-pinsync-opt-in),
   [ADR-0113](../../docs/plan/adr/0113-pinsync-modul.md)). Das Modul vergleicht die
   Werte beider Dateien auf Gleichheit; ob der Wert aktuell ist, sagen die
   `freshness-*`-Targets, und die Digests der Basis-Images und weitere Spiegel einer
   Pin-Hebung bleiben ungeprüft. Die Quellen liest es außerhalb der Scan-Menge
   (`scan.roots`, `scan.ignore` gelten dort nicht). Es liest den Dateistand: das
   `Makefile` setzt beide Werte mit `?=`, und ein Aufruf mit `GO_VERSION=…` oder eine
   Umgebungsvariable überschreibt sie zur Laufzeit, ohne dass das Modul es sieht.
**Wie groß der Ausschnitt ist, sagt das Kommando:** `make doc-check` nennt die
geprüfte Dateizahl in seiner Schluss-Zeile.

## Bindung

Bestandteil von `make gates`; zusätzlich im `pre-commit`-Hook als Doku-Guard.
[`DC-FA-LINK-001`](../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links) ·
[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) ·
[`MR-007`](../conventions/MR-007-aufloesung-mr-003.md)
