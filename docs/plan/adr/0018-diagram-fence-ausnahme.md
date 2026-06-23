# ADR-0018 — Diagramm-Fence-Ausnahme: das Modul `diagrams` liest gelistete Fences

**Status:** Accepted
**Datum:** 2026-06-23
**Autor:** pt9912
**Bezug:** [`DC-FA-DIAG-001`](../../../spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in)
(opt-in Modul `diagrams`), Mechanik-Präzedenz der opt-in-Module
[ADR-0010](0010-semgrep-hermetisches-gate.md) (distroless/netzlos als
Härtungslinie).
**Schärft:** [`spec/spezifikation.md` §DC-FA-LINK-001.a](../../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
(Fence-Opazität der Vorverarbeitung) und die neue Diagramm-Algorithmus-Sektion
(`.a`) der Spezifikation.

## Kontext

Die Markdown-Vorverarbeitung entfernt Fenced-Code-Blöcke für **alle** Module
([`spec/spezifikation.md` §DC-FA-LINK-001.a](../../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion))
— eine bewusste, stark gelebte Invariante (selbst `hostpaths` exemptiert Fences,
[`DC-FA-HOST-001`](../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in)).
Dadurch entgehen Kennungen, die in Diagrammen leben (z. B. Architektur-IDs in
`mermaid`-Flowcharts/Sequenzdiagrammen, Anlass: belief-agent), jeder Prüfung: ein
getipptes Token oder Drift gegen die definierende Tabelle bleibt unsichtbar. Eine
echte Diagramm-**Syntax**-Validierung scheidet aus — sie bräuchte den
Mermaid-Parser bzw. eine JS-Toolchain, unvereinbar mit dem distroless/I-O-armen
Kernvertrag
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).
Die d-check-förmige Frage ist Referenz-/**Existenz**-Konsistenz der Kennungen,
nicht Grammatik.

## Entscheidung

Das opt-in Modul `diagrams`
([`DC-FA-DIAG-001`](../../../spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in))
durchbricht die Fence-Opazität **scoped**:

1. **Nur gelistete Sprachen, nur für dieses Modul.** Ein eigener Fence-Automat
   (`diagramFenceLines`) öffnet ausschließlich Fences, deren Info-String in
   `diagrams.fences` steht (Default `mermaid`). Die gemeinsame Vorverarbeitung
   (`PreprocessMarkdown`/`proseLines`) bleibt **unverändert** — alle übrigen
   Module sehen Fence-Inhalt weiterhin nicht. Andere Fence-Sprachen (`bash`/`go`)
   bleiben opak (False-Positive-Schutz: ID-ähnliche Tokens in Code-Fences sind
   keine Referenzen).
2. **Reine Token-Extraktion, kein Parser.** Der Fence-Rohtext wird mit den
   `diagrams.patterns`-Regex gescannt; keine Mermaid-Grammatik. Deterministisch
   ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)),
   read-only/netzlos, keine zusätzliche Toolchain.
3. **Existenz statt Link-Policy.** In Fences ist kein Markdown-Link möglich; eine
   Diagramm-Kennung *kann* nicht verlinkt werden. Daher prüft das Modul Existenz:
   das Token muss in seiner `defined-in`-Quelle als eigenständiges Kennungs-Token
   **außerhalb von Fences** vorkommen (Heading **oder** Tabelle/Fließtext) — Befund
   `diagram-id-undefined`. Bewusst weiter als die heading-zentrierte Auflösung von
   [`DC-FA-ID-001`](../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids).
4. **`spans`-Isolation.** Das aktive Modul `spans` (ungeschlossene Code-Spans,
   [`DC-FA-SPAN-001`](../../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in))
   wird **nicht** berührt: `diagrams` arbeitet über seinen eigenen Fence-Automaten
   auf dem Rohtext und ändert den Fence-Zustand der gemeinsamen Vorverarbeitung
   nicht — keine `span-unclosed`-False-Positives durch geöffnete Diagramm-Fences.
5. **Default aus.** Ohne `diagrams`-Block ist der Befundsatz byte-identisch zum
   Lauf ohne das Modul ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).

## Verglichene Alternativen

| Alternative | Pro | Contra |
|---|---|---|
| **Scoped Fence-Öffnung + Existenz (gewählt)** | d-check-förmig, distroless-treu, kein Default-Delta | bricht die Fence-Opazität — scoped + opt-in + hier begründet |
| Mermaid-Parser / `mmdc` einbetten | echte Syntax-Validierung | JS-Toolchain/Browser, bricht den distroless/netzlosen Kernvertrag; Syntax ist nicht d-checks Domäne |
| Modul `ids` in alle Fences schauen lassen | wiederverwendet `ids` | Link-Policy in Fences unerfüllbar (unbehebbares `id-unlinked`); False Positives in Code-Fences |
| Nichts tun | kein Aufwand | Diagramm-Drift bleibt unsichtbar (der Anlass) |

## Konsequenzen

- Die Fence-Opazität erhält **eine** scoped, opt-in Ausnahme; Default und alle
  übrigen Module bleiben byte-identisch — die Invariante gilt weiter für jeden
  nicht gelisteten Fence und jedes andere Modul.
- „Definiert" ist vorkommen-basiert (nicht heading-zentriert), damit
  tabellen-definierte IDs erfasst werden — bewusst weiter als `ids`.
- Diagramm-**Syntax**-Validierung bleibt out of scope (Diagramm-Linter/CI, nicht
  d-checks distroless-Domäne).

## Fitness Function

- `diagrams` aktiv: undefinierte Diagramm-Kennung ⇒ `diagram-id-undefined`
  (Exit 1); ohne `diagrams`-Block byte-identisch (Negativtest,
  `internal/hexagon/core/rules/diagrams_test.go`); ein `bash`-Fence bleibt
  unberührt (Gegentest).

## Re-Evaluierungs-Trigger

- Bedarf nach Diagramm-**Syntax**-Prüfung → eigener Diagramm-Linter, nicht d-check.
- Bedarf nach Referenz-**Richtung** von Diagramm-Kennungen (matrix-artig) →
  Folge-ADR (Semantik-Erweiterung), nicht still im Modul.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-23 | Proposed → Accepted (Umsetzung slice-045: Modul `diagrams`, scoped Fence-Öffnung + Existenz-Auflösung; doc-first-Fundament R1→R2 ACCEPT; Anlass belief-agent-Architektur) |
