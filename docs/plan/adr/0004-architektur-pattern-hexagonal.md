# ADR-0004 — Architektur-Pattern: Hexagon light

**Status:** Accepted
**Datum:** 2026-06-10
**Autor:** pt9912
**Bezug:** [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit),
[`DC-FA-CLI-002`](../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
**Schärft:** `spec/architecture.md` (entsteht mit slice-002 und setzt
diesen Schnitt um) — nicht das Lastenheft.

## Kontext

Zwei NFAs sind strukturell erzwingbar statt nur versprechbar:
Seiteneffektfreiheit mit Netz nur im Modul `external`
([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit))
und Determinismus
([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).
Außerdem verlangt
[`DC-FA-CLI-002`](../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
einzeln schaltbare Regelmodule. Die Schwester-Repos (u-boot, m-trace,
b-cad, grid-guide) nutzen durchgängig hexagonale Schnitte mit
`arch-check`-Gates. d-check ist zugleich ein kleines Tool — das Risiko
ist Über-Zeremonie, nicht Unterstrukturierung.

## Entscheidung

**Hexagonale Architektur in bewusst leichter Ausprägung**
(„Hexagon light"): drei Bereiche, vier schmale Ports, keine
Application-/Domain-Unterteilung im Kern.

| Bereich | Inhalt | Regel |
|---|---|---|
| `internal/core` | Markdown-Analyse (Links, Headings, Fences), GitHub-Slug, Regelmodule (`links`, `anchors`, `ids`, `matrix`, `external`-Logik) als Strategien hinter einem gemeinsamen Interface, Befund-Modell, deterministische Sortierung, Pfad-/Escape-Regeln; **Port-Interfaces werden hier definiert** | keine I/O-Imports (`os`, `net`, `net/http`, `syscall`), kein Import aus `internal/adapter` | <!-- d-check:ignore (Layout-Skizze dieser ADR, real wurde ADR-0005) -->
| `internal/adapter/{fs,httpcheck,config,report}` | Datei-Discovery/Lesen/Symlink-Erkennung; HTTP-Erreichbarkeit; `.d-check.yml` laden/validieren ([ADR-0003](0003-config-format.md)); Text-/JSON-Ausgabe | nur `adapter/httpcheck` importiert `net/http`; nur `adapter/fs` importiert `os`; Adapter importieren einander nicht |
| `cmd/d-check` | CLI-Parsing, Composition Root (verdrahtet Adapter an Ports) | einziger driving Adapter — kein HTTP-Server o. ä. |

## Verglichene Alternativen

| Alternative | Pro | Contra |
|---|---|---|
| **Hexagon light (gewählt)** | NFAs als Import-Regeln maschinell prüfbar; Kern-Tests gegen In-Memory-FS-Port schnell und deterministisch; Stack-Konsistenz | minimal mehr Struktur als nötig für ein kleines Tool |
| Hexagonal voll (mit Application-Layer, Domain/UseCase-Trennung) | maximale Schichtung | Zeremonie ohne Nutzen — d-check hat einen Use Case („prüfe Repo") |
| Klassisch geschichtet (cli → service → util) | vertraut, wenig Begriffe | I/O sickert erfahrungsgemäß in die Mitte; keine scharfe, prüfbare Grenze für [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) |
| Flaches Paket mit Ad-hoc-Interfaces | minimaler Aufwand | keine maschinell prüfbare Architektur-Aussage → kein `arch-check`, NFA bleibt Versprechen |

## Konsequenzen

- `spec/architecture.md` (slice-002) dokumentiert genau diesen Schnitt
  und keine weitere Schichtung.
- Die Akzeptanztests der Regelmodule laufen gegen einen
  In-Memory-Filesystem-Port — kein echtes Dateisystem nötig, was die
  Determinismus-Messung ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)) trivialisiert.
- Neue Regelmodule (welle-03) sind Kern-Strategien, keine
  Architekturänderungen.

## Fitness Function

`make arch-check` (ab slice-003, depguard-/golangci-lint-Regeln):

1. `internal/core` importiert weder `os`, `net`, `net/http`, `syscall` <!-- d-check:ignore (Layout-Skizze dieser ADR) -->
   noch `internal/adapter/*`.
2. `net/http` ausschließlich in `internal/adapter/httpcheck`; <!-- d-check:ignore (Layout-Skizze dieser ADR) -->
   `gopkg.in/yaml.v3` ausschließlich in `internal/adapter/config`. <!-- d-check:ignore (Layout-Skizze dieser ADR) -->
3. `internal/adapter/*` importieren einander nicht.

Rot = diese ADR ist verletzt; eine Lockerung der Regeln ist eine neue
ADR (`AGENTS.md` §3.6).

## Re-Evaluierungs-Trigger

- Ein Regelmodul lässt sich nachweislich nicht ohne direkten I/O-Zugriff
  im Kern umsetzen.
- d-check erhält einen zweiten driving Adapter (z. B. LSP/Server-Modus)
  — dann ist die light-Ausprägung neu zu bewerten.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-10 | Proposed → Accepted (slice-001) |
