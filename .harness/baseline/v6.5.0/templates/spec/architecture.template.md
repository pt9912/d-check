# Architektur — <Projektname>

> **Template-Hinweis.** Diese Datei ist eine Vorlage. Sie ist
> **sprach- und meilensteinfrei** (siehe
> [Baseline-Regelwerk Modul 4](../../regelwerk/modul-04-adrs.md)).
> Kopiere sie nach `spec/architecture.md`, ersetze `<Platzhalter>` und
> lösche diesen Block.

**Status:** Aktiv. **Letzte Änderung:** YYYY-MM-DD.

**Rolle:** Sicht-Stratum — *keine* eigenen Anforderungen, derivativ. Regeln:
Baseline-Regelwerk `modul-03-spec.md` §Ziel-Form: Architektur-Sicht.

**Hard Rule:** Diese Datei enthält *keine* Wellen, Slices, Commit-Hashes
oder Closure-Daten, **keine ADR-Bezüge** — die Sicht steht im
Stabilitäts-Rang über der ADR — und **keine Historie**: `Letzte Änderung`
oben ist ein Frische-Marker, kein Protokoll. Die zeitliche Schicht lebt in
`docs/plan/planning/in-progress/roadmap.md` und den späteren Closure-Notizen.
Baseline-Regelwerk `modul-03-spec.md` §Ziel-Form: Architektur-Sicht.

---

## 1. Komponenten-Übersicht

Regeln dieser Sektion: **Hier werden die `ARC-*` für Komponenten vergeben** —
eine Zeile je Kasten des Diagramms, damit es *eine* Stelle gibt. Die Kennung
ist eine Adresse, damit ein Slice sagen kann, welche Komponente er
berührt; sie ist **keine** Anforderung. Gezählt wird fortlaufend je Datei —
§3 setzt die Reihe fort, statt neu zu beginnen (Baseline-Regelwerk
`grundlagen-source-precedence.md` §ID-Schema als Klammer, §Vergabe).

<!--
Ein Diagramm (Mermaid oder ASCII) der Top-Level-Komponenten.
Jeder Kasten benennt die Schicht/Rolle, nicht die Technologie.

Mermaid-Beispiel siehe unten — zwingend ist nur die Klarheit,
nicht das Format.
-->

```mermaid
flowchart TB
    UI[UI / API-Layer]
    Runtime[Runtime / Bootstrap]
    Service[Service-Layer]
    Repo[Repository-Layer]
    Config[Config-Layer]
    Types[Types / Domain]

    UI --> Service
    Runtime --> Service
    Service --> Repo
    Service --> Types
    Repo --> Config
    Repo --> Types
    Config --> Types
```

| ID | Komponente | Rolle |
|---|---|---|
| `ARC-001` | Types / Domain | <…> |
| `ARC-002` | Config-Layer | <…> |
| `ARC-003` | Repository-Layer | <…> |
| `ARC-004` | Service-Layer | <…> |
| `ARC-005` | Runtime / Bootstrap | <…> |
| `ARC-006` | UI / API-Layer | <…> |

## 2. Schichten und Constraints

Regeln dieser Sektion: Welche ADR eine Layering-Regel verbindlich macht,
deklariert die ADR aufwärts in ihrem `Schärft:`-Feld — kein ADR-Bezug in dieser
Sicht (Baseline-Regelwerk `grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)).

<!--
Pro Schicht: was sie tut, was sie *nicht* tut. Layering-Regeln, die
durch ArchUnit / depguard / import-linter durchgesetzt werden. Welche
ADR eine Regel verbindlich macht, deklariert die ADR aufwärts in ihrem
Schärft:-Feld (Baseline-Regelwerk `grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)) — kein ADR-Bezug in dieser Sicht.

Beispiel-Schema (aus OpenAI-Layering, siehe Modul 4):
Types → Config → Repo → Service → Runtime → UI
-->

Eine Schicht ist eine *Gruppierung* über Komponenten, keine eigene Sache:
Fällt sie mit einer Komponente zusammen, nennt die Zeile deren `ARC-*` aus §1;
umfasst sie mehrere, bleibt die Spalte leer und die Constraint gilt für alle.

| Komponente(n) | Schicht | Verantwortlichkeit | Darf importieren | Darf NICHT importieren |
|---|---|---|---|---|
| `ARC-001` | Types | Domain-Modell, Pure | — | alles andere |
| `ARC-002` | Config | Konfiguration laden/validieren | Types | Service, Runtime, UI |
| `ARC-003` | Repo | Datenzugriff | Types, Config | Service, Runtime, UI |
| `ARC-004` | Service | Geschäftslogik | Types, Config, Repo | Runtime, UI |
| `ARC-005` | Runtime | Bootstrap, DI | alles oben | — |
| `ARC-006` | UI | API / CLI / GUI | alles oben außer Repo | Repo direkt |

## 3. Externe Abhängigkeiten

Regeln dieser Sektion: Auch die Schnittstelle zu einem externen System trägt
eine `ARC-*` — die Kennung benennt den *Berührungspunkt*, nicht das fremde
System (Baseline-Regelwerk `grundlagen-source-precedence.md` §ID-Schema als Klammer).

<!--
Welche externen Systeme/Bibliotheken sind Teil der Architektur. Die
Wahl-Begründung steht in der ADR (die ihre Schärft auf diese Sicht
deklariert), nicht hier.
-->

| ID | System | Rolle | Substituierbarkeit |
|---|---|---|---|
| `ARC-007` | <…> | <…> | <…> |

## 4. Sequenz-Diagramme

<!--
Für jeden kritischen Use-Case (aus Lastenheft) eine Sequenz.
Schichten als Lanes, Aktionen mit IDs aus dem Lastenheft.
-->

### Use-Case: <LH-FA-NN — Titel>

```mermaid
sequenceDiagram
    participant UI
    participant Service
    participant Repo
    UI->>Service: <Anfrage>
    Service->>Repo: <Datenzugriff>
    Repo-->>Service: <Antwort>
    Service-->>UI: <Ergebnis>
```

## 5. Fehlermodelle und Resilienz

<!-- Wo werden Fehler abgefangen, propagiert, geloggt. -->

| Fehlerquelle | Behandlung-Schicht | Logging |
|---|---|---|
| <…> | <…> | <…> |
