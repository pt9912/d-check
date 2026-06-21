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
Baseline, Modus-Deklarationen pro Sub-Area, Zusatzklassen für
Sensors-Bindung) leben in
[`harness/conventions.md`](harness/conventions.md).

Das Betriebsregelwerk der adoptierten Baseline in Agenten-Kurzform
einmal pro Session lesen, bevor der Workflow (§6) startet. Lese-Form:
das nach Modulen und Grundlagen-Abschnitten aufgeteilte Release-Bundle
[`lab-regelwerk.zip`](https://github.com/pt9912/ai-harness-course/releases/download/v1.3.0/lab-regelwerk.zip)
(`v1.3.0`) — so lädt ein Agent einen einzelnen Abschnitt, ohne das
gesamte Regelwerk im Kontext zu halten. Das Bundle ist eine derivative
Sicht auf die Quelldatei
[`agents-regelwerk.md`](https://raw.githubusercontent.com/pt9912/ai-harness-course/v1.3.0/kurs/de/agents-regelwerk.md);
bei Konflikt gilt die Quelldatei, über ihr die kanonischen Quellen
(Source Precedence). Der adoptierte Stand steht in
[`harness/conventions.md`](harness/conventions.md) §Baseline.

## 2. Kanonische Quellen (Source Precedence)

In dieser Reihenfolge:

1. [`spec/lastenheft.md`](spec/lastenheft.md) — vertraglich abnahmebindend.
2. [`spec/spezifikation.md`](spec/spezifikation.md) — technisch verbindlich, fortschreibbar.
3. [`spec/architecture.md`](spec/architecture.md) — Komponenten- und Sequenzsicht.
4. [`docs/plan/adr/README.md`](docs/plan/adr/README.md) — ADR-Index.
5. [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md) — aktuelle Welle.
6. [`docs/user/`](docs/user/) — Operations, Releasing.
7. [`README.md`](README.md) — Projekt-Überblick.
8. **AGENTS.md (diese Datei).**
9. [`harness/README.md`](harness/README.md) — Harness-Einstieg.

## 3. Harte Regeln

### 3.1 Docker/make-only

Implementierungssprache ist **Go**
([ADR-0001](docs/plan/adr/0001-implementierungssprache.md)). Es gilt:
**kein Host-Go und keine Host-Paketmanager** (`go`, `pip`, `npm`,
`cargo`, `apt`, `brew`, …). Alle Checks laufen über `make`; die
Go-Toolchain läuft in Docker (Multi-Stage gemäß
[ADR-0002](docs/plan/adr/0002-distribution-ghcr-image.md), entsteht
mit slice-003). Der Host braucht nur `git`, GNU `make`, `bash` und
Docker.

**Falsch:** `go build ./…`, `go test ./…`, `pip install …`
**Richtig:** `make gates` (Implementierungs-Gates entstehen mit slice-003)

**Begründung:** Toolchain-Reproduzierbarkeit + Supply-Chain-Defense.

### 3.2 Suppression-Verbot

Inline-Suppressions sind verboten: `//nolint`-Direktiven im Code
brechen das künftige Suppression-Gate. Ausnahmen leben zentral in
`.golangci.yml` (exclude-rules) mit Begründung — die Datei entsteht
mit slice-003.

### 3.3 git mv + Inhaltsänderung = zwei Commits

Wenn eine Datei verschoben **und** der Inhalt umgeschrieben wird:

1. `git mv source target` → eigener Commit (reiner Move, Git erkennt R-Rename).
2. Inhalt umschreiben → zweiter Commit.

**Begründung:** Sonst fällt die Rename-Detection unter die
50%-Similarity-Schwelle und `git log --follow` wird unzuverlässig.

**Ausnahme Slice-Lifecycle-Move (`in-progress/` → `done/`):** Der
`git mv`-Commit trägt hier **zusätzlich** den Roadmap-Flip §Aktuelle Welle
(zurück auf „Keine aktive Welle") und alle Pfad-Verweise auf den Slice
(Roadmap, §4, `harness/README.md` §Sensors) von `in-progress/` nach
`done/`. Sonst ist der Commit gate-rot: `make planning-check` koppelt
in-progress-Stand und Roadmap atomar, und die alten Verweise laufen ins
Leere (`target-missing`). Nur der **Slice-Body** (Status-Zeile,
Closure-Notiz) bleibt Commit 2 — die Slice-Datei selbst ist im Move-Commit
unverändert, also hält die Rename-Detection. Kanonisch:
[`MR-013`](harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise).

### 3.4 Architektur sprach-/meilensteinfrei; Spec-Straten nie abwärts

[`spec/architecture.md`](spec/architecture.md) benennt Schichten und
Rollen statt Technologie — keine Sprach-/Modul-Pfade. Kein
Spec-Stratum (auch [`spec/spezifikation.md`](spec/spezifikation.md))
referenziert ADRs, Wellen, Slices, Commit-Hashes oder Closure-Daten.
Die sprachkonkrete Übersetzung (Modul-Pfade, Import-Regeln) und die
Begründungen leben in den ADRs, deren `Schärft:`-Feld aufwärts zeigt;
die zeitliche Schicht lebt in `docs/plan/planning/`.

### 3.5 ADRs sind nach `Accepted` immutable

Eine ADR mit Status `Accepted` wird nicht inhaltlich überschrieben.
Korrekturen entstehen als neue ADR mit `Supersedes ADR-NN`.

### 3.6 Gates dürfen nicht ohne ADR gelockert werden

Jede Schwellen-Senkung (Coverage, Linter-Strenge, Prüfregel) ist ein
ADR, kein PR-Kommentar.

## 4. Quality Gates

Nur hier gelistete Targets existieren im Makefile. Halluzinierte
Gates sind die häufigste Form von Harness-Lüge.

| Target                       | Zweck                                                                                                                                                                                                                                              |
| ---------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `make lint`                  | golangci-lint mit dem Projekt-Profil (§3.2)                                                                                                                                                                                                        |
| `make test`                  | `go test ./...` — Akzeptanzkriterien der `DC-FA-*`                                                                                                                                                                                                 |
| `make arch-check`            | Import-Regeln des Hexagon-Schnitts + Kern-Paket-Richtung ([ADR-0005](docs/plan/adr/0005-modul-layout-hexagon-ordner.md), [ADR-0012](docs/plan/adr/0012-kern-paketschnitt-model-rules-app.md))                                                      |
| `make coverage-gate`         | Coverage-Schwelle über `./internal/...` (Kalibrierungs-Bindung, siehe [`harness/README.md`](harness/README.md) §Sensors)                                                                                                                           |
| `make gate-consistency`      | Meta-Gate: dokumentierte Targets ↔ Makefile + [`DC-QA-03`](spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Modulliste (Harness-Lügen-Schutz)                                                                           |
| `make planning-check`        | Meta-Gate: Roadmap §Aktuelle Welle ↔ `in-progress/slice-*` (Planning-Drift-Schutz) ([slice-040](docs/plan/planning/done/slice-040-planning-consistency-gate.md))                                                                            |
| `make doc-check`             | Doku-Links, Anker, Kennungs-Linkpflicht, Referenzmatrix + Inline-Code-Pfade via `d-check` selbst (Dogfooding; netzlos — zugleich [`DC-QA-03`](spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Messmethode)             |
| `make gates`                 | alle inneren Gates (mandatory vor Handoff)                                                                                                                                                                                                         |
| `make ci`                    | CI-äquivalenter Lauf: gates + image-test (fährt die Release-Pipeline)                                                                                                                                                                              |
| `make trace-check`           | Traceability-Gate: DC-/ADR-/slice-ID in Commits (`commit-msg`-Hook + PR-CI; bewusst **nicht** Teil von `gates`/`ci`) ([ADR-0013](docs/plan/adr/0013-pr-ci-und-traceability-gate.md))                                                               |
| `make adr-check`             | ADR-Immutable-Gate: `Accepted`-ADRs nicht inhaltlich ändern (`pre-commit`-Hook + PR-CI; bewusst **nicht** Teil von `gates`/`ci`) ([ADR-0016](docs/plan/adr/0016-adr-immutable-gate.md))                                                            |
| `make hooks`                 | git-Hooks installieren (`core.hooksPath` → `.githooks`; aktiviert `commit-msg`-Traceability + `pre-commit`-ADR-Immutable) ([ADR-0013](docs/plan/adr/0013-pr-ci-und-traceability-gate.md), [ADR-0016](docs/plan/adr/0016-adr-immutable-gate.md))    |
| `make fullbuild`             | volle Closure: gates + image-test + bench, schließt mit dem Image-Hash                                                                                                                                                                             |
| `make image-test`            | [`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image)-Akzeptanzkriterien gegen das lokale Image (nativ vs. Container)                                                                                                                |
| `make bench`                 | [`DC-QA-01`](spec/lastenheft.md#dc-qa-01--performance)-Benchmark gegen generiertes Fixture (Median aus 3 Läufen, kein Gate in `gates`)                                                                                                             |
| `make semgrep`               | Security-/Static-Analysis-**Gate**: gepinntes semgrep-Image + gepinntes, lokal gecachtes `go/lang/security`-Regelset, netzloser Scan (`--network none`); **Bestandteil von `gates`** ([ADR-0010](docs/plan/adr/0010-semgrep-hermetisches-gate.md)) |
| `make versions`              | Reproduzierbarkeits-Pins ausgeben (Go, Lint, Basis-Images, Runtime-Image-ID)                                                                                                                                                                       |
| `make build` / `make run`    | Runtime-Image bauen / Selbst-Smoke-Test                                                                                                                                                                                                            |
| `make deps` / `make compile` | Cache-Layer / schnelles Compile-Feedback                                                                                                                                                                                                           |
| `make record-gates`          | Nachweis schreiben: Working-Tree-Hash für den Stop-Hook                                                                                                                                                                                            |
| `make help` / `make clean`   | Targets anzeigen / Images entfernen                                                                                                                                                                                                                |

Alle dokumentiert-geplanten Targets existieren; Details und Bindungen:
Sensors-Tabelle in [`harness/README.md`](harness/README.md).

## 5. Dokumentations-Regeln

- Commits/PRs müssen mindestens eine `DC-*`-, `ADR-*`-, `MR-*`- oder
  `slice-*`-ID nennen (maschinell erzwungen: `make trace-check` /
  `commit-msg`-Hook / PR-CI,
  [ADR-0013](docs/plan/adr/0013-pr-ci-und-traceability-gate.md);
  Ausnahme: Merge-/Revert-Commits). Vergeben werden IDs nur beim
  Spec-/ADR-Schreiben nach dem deklarierten Schema
  ([`MR-008`](harness/conventions.md#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage))
  — nie ad hoc im Commit/PR; Agenten referenzieren IDs, sie erfinden
  keine.
- Neue oder geänderte `DC-*`-Anforderungen entstehen nur in
  [`spec/lastenheft.md`](spec/lastenheft.md) — nie per ADR (ADRs
  schärfen die Spezifikation, nicht das Lastenheft). Der
  Anlege-Prozess (Schema, Akzeptanzkriterien-Trio, Versions-Bump +
  Historie, Beleg-Pflicht) steht in
  [`harness/conventions.md` §Anforderungs-Anlege-Prozess](harness/conventions.md#anforderungs-anlege-prozess).
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
