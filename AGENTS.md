# AGENTS.md — Briefing für AI-Coding-Agenten

## 1. Was diese Datei ist

Onboarding-Briefing für jede AI-Session, die in diesem Repo Code oder
Dokumentation ändert. Sie verweist auf die kanonischen Quellen und
formuliert die Hard Rules, die der Implementation-Agent immer
einhalten muss.

**Bei Konflikt zwischen dieser Datei und einer kanonischen Quelle gilt
die kanonische Quelle** (Source Precedence — siehe
[`harness/README.md`](harness/README.md)).

Strukturregeln (ID-Schemata, Verzeichniskonvention, Adaptionen ggü.
Baseline, Modus-Deklarationen pro Sub-Area) leben in
[`harness/conventions.md`](harness/conventions.md).

## 2. Kanonische Quellen (Source Precedence)

In dieser Reihenfolge:

1. [`spec/lastenheft.md`](spec/lastenheft.md) — vertraglich abnahmebindend.
2. `spec/spezifikation.md` — technisch verbindlich, fortschreibbar. **Geplant** (slice-002), existiert noch nicht.
3. `spec/architecture.md` — Komponenten- und Sequenzsicht. **Geplant** (slice-002), existiert noch nicht.
4. [`docs/plan/adr/README.md`](docs/plan/adr/README.md) — ADR-Index.
5. [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md) — aktuelle Welle.
6. [`README.md`](README.md) — Projekt-Überblick.
7. **AGENTS.md (diese Datei).**
8. [`harness/README.md`](harness/README.md) — Harness-Einstieg.

## 3. Harte Regeln

### 3.1 make-only, keine Host-Toolchain

Die Implementierungssprache ist noch nicht entschieden (ADR-0001,
slice-001). Bis dahin — und auch danach — gilt: **keine
Host-Toolchain-Installationen und keine Host-Paketmanager** (`pip`,
`npm`, `cargo`, `apt`, `brew`, …). Alle Checks laufen über `make`;
nach der Sprachentscheidung läuft die Toolchain in Docker. Der Host
braucht nur `git`, GNU `make`, `bash` und Docker.

**Falsch:** `pip install …`, `npm install …`, `apt-get install …`
**Richtig:** `make doc-check`, `make gates`

**Begründung:** Toolchain-Reproduzierbarkeit + Supply-Chain-Defense.

### 3.2 Suppression-Verbot

Inline-Suppressions (Linter-Ausnahmen im Code) sind verboten.
Ausnahmen leben in einer zentralen Konfigurationsdatei mit Begründung.
Die sprachkonkreten Marker (z. B. `# noqa`, `//nolint`) werden mit
ADR-0001 hier ergänzt.

### 3.3 git mv + Inhaltsänderung = zwei Commits

Wenn eine Datei verschoben **und** der Inhalt umgeschrieben wird:

1. `git mv source target` → eigener Commit (reiner Move, Git erkennt R-Rename).
2. Inhalt umschreiben → zweiter Commit.

**Begründung:** Sonst fällt die Rename-Detection unter die
50%-Similarity-Schwelle und `git log --follow` wird unzuverlässig.

### 3.4 Architektur ist sprach- und meilensteinfrei

`spec/architecture.md` (sobald vorhanden) referenziert ADRs und
Modul-Pfade, aber **keine** Wellen, Slices, Commit-Hashes oder
Closure-Daten. Die zeitliche Schicht lebt in `docs/plan/planning/`.

### 3.5 ADRs sind nach `Accepted` immutable

Eine ADR mit Status `Accepted` wird nicht inhaltlich überschrieben.
Korrekturen entstehen als neue ADR mit `Supersedes ADR-NN`.

### 3.6 Gates dürfen nicht ohne ADR gelockert werden

Jede Schwellen-Senkung (Coverage, Linter-Strenge, Prüfregel) ist ein
ADR, kein PR-Kommentar.

## 4. Quality Gates

Nur hier gelistete Targets existieren im Makefile. Halluzinierte
Gates sind die häufigste Form von Harness-Lüge.

| Target | Zweck |
|---|---|
| `make doc-check` | interne Markdown-Linkziele existieren (Bootstrap-Sensor) |
| `make record-gates` | Nachweis schreiben: Working-Tree-Hash für den Stop-Hook |
| `make gates` | alle inneren Gates (mandatory vor Handoff) |
| `make help` | Targets anzeigen |

**Nicht behauptet** (geplant): `make lint`, `make test`,
`make arch-check`, `make coverage-gate` — entstehen ab slice-003 mit
der Implementierung.

## 5. Dokumentations-Regeln

- Commits/PRs müssen mindestens eine `DC-*`- oder `ADR-*`-ID nennen.
- Neue ADRs müssen den ADR-Index aktualisieren.
- Roadmap/Status-Geschichte lebt in `docs/plan/planning/`, nicht in der Architektur-Spec.
- Slice-Lifecycle (`open → next → in-progress → done`) ist reine Datei-Bewegung (`git mv`, siehe §3.3).
- `CHANGELOG.md` wird bei nutzersichtbaren Änderungen gepflegt.

## 6. Minimal Agent Workflow

Pro Slice:

1. [`harness/README.md`](harness/README.md) lesen.
2. Relevante kanonische Quelle lesen (Source Precedence beachten).
3. Betroffene Requirement-/ADR-IDs identifizieren.
4. Kleinste sinnvolle Änderung planen.
5. Engsten nützlichen Sensor laufen lassen.
6. Repo-weiten Gate-Lauf vor Handoff (`make gates`).
7. Doku/Indizes aktualisieren, falls ein öffentlicher Vertrag berührt.
8. Ausgeführte Sensors und verbleibende Risiken berichten — keine Erfolgsmeldung ohne Gate-Ausführung.
