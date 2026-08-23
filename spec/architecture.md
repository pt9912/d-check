# Architektur — d-check

**Status:** Aktiv. **Letzte Änderung:** 2026-08-22.

**Rolle:** Sicht-Stratum — *keine* eigenen Anforderungen, derivativ. Regeln:
Baseline-Regelwerk
[`modul-03-spec.md` §Ziel-Form: Architektur-Sicht](../.harness/baseline/v5.11.0/regelwerk/modul-03-spec.md#ziel-form-architektur-sicht).

**Hard Rule:** Diese Datei ist **sprach- und meilensteinfrei**: Sie
benennt Schichten und Rollen, keine Technologie, und enthält keine
Wellen, Slices, Commit-Hashes oder Closure-Daten (`AGENTS.md` §3.4).
Die zeitliche Schicht lebt in `docs/plan/planning/`. Und **keine Historie**:
`Letzte Änderung` oben ist ein **Frische-Marker** — die Aussage, wann die
Sicht zuletzt **gegen den Code** gehalten wurde —, kein Protokoll. Deshalb
bleibt die Zeile hier, wo lebende Register und das Technik-Stratum keine
tragen. Sie wird gesetzt, wenn jemand die Sicht gegen den Code hält; steht
sie still, ist das die ehrliche Aussage „seither nicht geprüft" und kein
Versäumnis einer automatischen Pflege.

**Kennungen:** Komponenten und externe Berührungspunkte tragen eine
**Struktur-ID** `ARC-<NNN>` — fortlaufend je Datei, vergeben in §1 und in §3
fortgesetzt; §2 nennt dieselben Kennungen und führt keine neuen ein. Eine
Struktur-ID macht adressierbar; sie ist **keine** Anforderung und trägt keine
Zusage — Zusagen leben im Lastenheft.

**Bezug:** [`lastenheft.md`](lastenheft.md) und
[`spezifikation.md`](spezifikation.md) — bei Konflikt gewinnen diese.
Der hexagonale Schnitt („Hexagon light") ist per ADR entschieden; die
Begründung und die sprachkonkrete Übersetzung (Modul-Pfade,
Import-Regeln) leben dort, und die ADRs deklarieren ihre Schärfung
dieser Datei aufwärts (`Schärft:`-Feld). Spec-Straten verweisen nicht
abwärts auf ADRs oder Planning-Artefakte.

---

## 1. Komponenten-Übersicht

```mermaid
flowchart TB
    CLI["ARC-001 CLI — einziger driving Adapter,<br/>Composition Root"]
    CORE["ARC-002 Kern — Regelmodule, Markdown-Analyse,<br/>Befund-Modell; definiert die Ports"]
    FS["ARC-003 Filesystem-Adapter"]
    HTTP["ARC-004 HTTP-Adapter"]
    VCS["ARC-005 VCS-/git-Adapter"]
    CFG["ARC-006 Config-Adapter"]
    REP["ARC-007 Reporter-Adapter"]

    CLI -->|"ruft auf, verdrahtet"| CORE
    CLI --> CFG
    CLI --> REP
    FS -.->|"implementiert Filesystem-Port"| CORE
    HTTP -.->|"implementiert HTTP-Port"| CORE
    VCS -.->|"implementiert VCS-Port (opt-in)"| CORE
    CFG -.->|"liefert validierte Config"| CORE
    REP -.->|"konsumiert Befundliste"| CORE
```

| Kennung | Komponente | Rolle |
|---|---|---|
| `ARC-001` | CLI | einziger driving Adapter; Argument-Parsing, Composition Root, Exit-Code |
| `ARC-002` | Kern | Regelmodule, Markdown-Analyse, Befund-Modell; definiert die Ports |
| `ARC-003` | Filesystem-Adapter | Datei-Discovery, Lesen, Symlink-Erkennung |
| `ARC-004` | HTTP-Adapter | Erreichbarkeit externer Verweise |
| `ARC-005` | VCS-/git-Adapter | rein lesender Zugriff auf die Versionsgeschichte |
| `ARC-006` | Config-Adapter | Konfiguration dekodieren und validieren |
| `ARC-007` | Reporter-Adapter | Rendern der Befundliste |

Das CLI ist der **einzige driving Adapter** (kein Server-Modus). Die
Regelmodule (`links`, `anchors`, `ids`, `matrix`, `external`) sind
Strategien im Kern hinter einem gemeinsamen Interface
([`DC-FA-CLI-002`](lastenheft.md#dc-fa-cli-002--regelmodul-auswahl));
neue Module sind Kern-Erweiterungen, keine Architekturänderungen. Der Kern
ist intern in drei Pakete mit einbahniger Importrichtung geschnitten —
`model` (Daten) ← `rules` (Prüf-Engine: Module + Orchestrierung) ← `app`
(Modi `--doctor`/`--repair`/`--suggest-config`).

## 2. Schichten und Constraints

| Kennung | Schicht / Rolle | Verantwortlichkeit | Darf nutzen | Darf NICHT nutzen |
|---|---|---|---|---|
| `ARC-002` | Kern | Markdown-Vorverarbeitung, Link-/Heading-/Kennungs-Extraktion, Slug, Regelmodule, Befund-Modell, deterministische Sortierung, Pfad-/Escape-Regeln; definiert die Ports | reine Standardbibliothek ohne I/O; Port-Interfaces | Dateisystem-, Netzwerk-, Prozess-APIs; Adapter; YAML-Bibliothek |
| `ARC-003` | Filesystem-Adapter | Datei-Discovery, Lesen, Symlink-Erkennung (Lstat) | Kern-Ports; Dateisystem-API | andere Adapter; Netzwerk |
| `ARC-004` | HTTP-Adapter | HEAD/GET-Erreichbarkeit, Timeout, Redirect-Limit | Kern-Ports; HTTP-Client | andere Adapter; Dateisystem |
| `ARC-005` | VCS-/git-Adapter | Lesen der git-Historie aus `.git`: Datei-Inhalt an einem Commit-Ref, geänderte Pfade einer Commit-Range, Commit-Messages einer Range; **rein lesend**, ohne externes git-Binary, ohne Netz (opt-in Module `vcs`, `commits`) | Kern-Ports; git-Objekt-Bibliothek; read-only `.git` | andere Adapter; Netzwerk; Schreiben ins Repository |
| `ARC-006` | Config-Adapter | `.d-check.yml` strikt dekodieren, zweistufig validieren — den Datei-Inhalt beschafft das CLI über den Filesystem-Adapter, der die einzige Dateisystem-Tür bleibt | Kern-Typen; YAML-Bibliothek | andere Adapter; Dateisystem; Netzwerk |
| `ARC-007` | Reporter-Adapter | Text-/JSON-Rendering auf stdout/stderr | Kern-Typen; Serialisierung | andere Adapter; Netzwerk; Dateizugriffe jenseits stdout/stderr |
| `ARC-001` | CLI | Argument-Parsing, Composition Root, Exit-Code | alles oben | — |

Diese Tabelle ist die rollenbasierte Quelle der
`arch-check`-Fitness-Function; ihre sprachkonkrete Übersetzung
(Modul-Pfade, Import-Regeln) trägt die zugehörige ADR und deklariert
sie aufwärts. Eine Lockerung ist eine neue ADR (`AGENTS.md` §3.6).

## 3. Externe Abhängigkeiten

| Kennung | System | Rolle | Substituierbarkeit |
|---|---|---|---|
| `ARC-008` | YAML-Bibliothek | Decoding im Config-Adapter | hoch — vollständig im Adapter gekapselt |
| `ARC-009` | HTTP-Client der Standardbibliothek | Erreichbarkeits-Checks im HTTP-Adapter | hoch — hinter dem HTTP-Port |
| `ARC-010` | git-Objekt-Bibliothek (rein in der Implementierungssprache, **kein** externes git-Binary) | Lesen von `.git` im VCS-Adapter (opt-in Module `vcs`, `commits`) | mittel — hinter dem VCS-Port; read-only, netzlos |
| `ARC-011` | Minimal-Runtime ohne Shell/Paketmanager | Auslieferung | mittel — CA-Bundle-/Non-root-Annahmen |

## 4. Sequenz-Diagramme

### Use-Case: [`DC-FA-CLI-001`](lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel) — Prüflauf

```mermaid
sequenceDiagram
    participant CLI
    participant CFG as Config-Adapter
    participant CORE as Kern
    participant FS as Filesystem-Adapter
    participant HTTP as HTTP-Adapter
    participant VCS as VCS-Adapter
    participant REP as Reporter

    CLI->>FS: Read(.d-check.yml)
    FS-->>CLI: Datei-Inhalt | nicht vorhanden (→ Defaults)
    CLI->>CFG: dekodieren + validieren
    CFG-->>CLI: Config | Fehler (Exit 2)
    CLI->>CORE: Run(Config, Ports)
    CORE->>FS: Discover(Wurzeln, Ignores)
    FS-->>CORE: Dateiliste (sortiert)
    loop pro Datei
        CORE->>FS: Read / Lstat
        CORE->>CORE: Module: links, anchors, ids, matrix
    end
    opt Modul external aktiv
        CORE->>HTTP: Head/Get(URL)
        HTTP-->>CORE: Status | Timeout
    end
    opt Modul vcs aktiv (--range/--staged)
        CORE->>VCS: FileAtRef(BASE/HEAD, Pfad) / ChangedFiles(Range)
        VCS-->>CORE: Inhalt | Diff | Fehler (.git/Range fehlt → Exit 2)
    end
    opt Modul commits aktiv (--range)
        CORE->>VCS: CommitMessages(Range)
        VCS-->>CORE: Messages | Fehler (.git/Range fehlt → Exit 2)
    end
    CORE-->>CLI: Befundliste (sortiert, DC-QA-02)
    CLI->>REP: Render(Text | JSON)
    CLI-->>CLI: Exit 0 | 1
```

## 5. Fehlermodelle und Resilienz

| Fehlerquelle | Behandelnde Schicht | Wirkung |
|---|---|---|
| ungültige CLI-Nutzung | CLI | stderr, Exit 2 |
| ungültige `.d-check.yml` (Syntax/Semantik) | Config-Adapter | stderr mit Zeilenangabe, Exit 2, keine Prüfung |
| Scan-Wurzel/Datei nicht lesbar | Filesystem-Adapter → CLI | stderr, Exit 2 (kein Teilergebnis als Erfolg) |
| Regelverletzung in Doku | Kern (Regelmodule) | Befund, Exit 1 |
| HTTP-Fehler/Timeout (`external`) | HTTP-Adapter → Kern | Befund, kein Abbruch |
| `.git` fehlt/unlesbar, Range unauflösbar oder Message-Datei unlesbar (`vcs`, `commits`) | VCS-Adapter → CLI | stderr, Exit 2 (fail-closed; eine fehlende git-Eingabe ist kein stilles Grün) |

Der Kern wirft keine I/O-Fehler selbst — sie erreichen ihn als
Port-Ergebnisse; das geprüfte Repository wird nie beschrieben
([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
