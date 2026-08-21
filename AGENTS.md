# AGENTS.md — Briefing für AI-Coding-Agenten

## 1. Was diese Datei ist

Onboarding-Briefing für jede AI-Session, die in diesem Repo Code oder
Dokumentation ändert. Sie verweist auf die kanonischen Quellen und
formuliert die Hard Rules, die der Implementation-Agent immer
einhalten muss.

Diese Datei trägt **Hard Rules und Pointer** auf die kanonischen Quellen und
**dupliziert deren Inhalt nicht** — sonst entsteht Drift (Kanon:
[`modul-09-implementierung.md` §AGENTS.md-Regeln](.harness/baseline/v5.6.0/regelwerk/modul-09-implementierung.md#agentsmd-regeln-modul-9)).

**Bei Konflikt zwischen dieser Datei und einer kanonischen Quelle gilt
die kanonische Quelle** (Source Precedence — siehe
[`harness/README.md`](harness/README.md)).

Strukturregeln (ID-Schemata, Verzeichniskonvention, Adaptionen ggü.
Baseline, Modus-Deklarationen pro Sub-Area, Zusatzklassen für
Sensors-Bindung) leben in
[`harness/conventions.md`](harness/conventions.md).

Das Betriebsregelwerk der adoptierten Baseline ist **committet vendored**:
das nach Modulen und Grundlagen-Abschnitten aufgeteilte Regelwerk liegt
entpackt unter `.harness/baseline/<tag>/regelwerk/` (die dortige `README.md`
ist der Index), samt `.harness/baseline/<tag>/SHA256SUMS`-Integritätsmanifest —
**netzlos auf jedem Checkout präsent**, offline materialisier-/verifizierbar
per `tools/harness/fetch-baseline-cache.sh` (`--verify` offline-Integrität;
`--check-latest` = Currency- + Content-Drift-Audit ggü. Upstream, informativ/kein Gate,
[`MR-022`](harness/conventions.md#mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019); Tag aus §Baseline;
Quelle ist das derivative Release-Bundle
[`lab-regelwerk.zip`](https://github.com/pt9912/ai-harness-course/releases/download/v5.6.0/lab-regelwerk.zip);
Pfadschema/Provenance siehe
[`harness/conventions.md`](harness/conventions.md) §Adoptierte Konventions-Quellen,
[`MR-019`](harness/conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)).
Pro Session **nur den benötigten Abschnitt** lesen, bevor der Workflow (§6)
startet — nicht das gesamte Regelwerk im Kontext halten.
Die **Skelett-Vorlagen** der Baseline liegen aus demselben self-contained Bundle
**committet vendored** unter `.harness/baseline/<tag>/templates/` (parallel zum
`.harness/baseline/<tag>/regelwerk/`-Baum, netzlos) und tragen zwei Rollen: als
**Referenz-Form**, auf die das Regelwerk als „Ziel-Form" verweist, und als **Vorlage**
beim Anlegen neuer Artefakte (ADR, Slice, Welle, …). d-checks gelebte Slice-/ADR-Struktur
folgt dabei einer **Haus-Stil-Form** — in Etappe C als baseline-konforme Form-Wahl
aufgelöst, nicht als Fork. Das Bundle ist derivativ; bei Konflikt sticht die Quelldatei
das Bundle, über ihr die kanonischen Quellen (Source Precedence). Stand/Provenance führt
[`harness/conventions.md`](harness/conventions.md) (§Adoptierte Konventions-Quellen bzw.
§Baseline).

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
Leere (`target-missing`). Nur der **Slice-Body** (DoD-Haken + Closure-Notiz;
historische Slices auch die Status-Zeile) bleibt Commit 2 — die Slice-Datei selbst ist im Move-Commit
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
Korrekturen entstehen als neue ADR mit `Supersedes ADR-NNNN` (vierstellig).
Maschinell erzwungen über `make adr-check` (`pre-commit`-Hook + PR-/Push-CI;
erlaubt bleiben `## Geschichte`-Anhänge + der `**Status:**`-Übergang;
[ADR-0016](docs/plan/adr/0016-adr-immutable-gate.md)).

### 3.6 Gates dürfen nicht ohne ADR gelockert werden

Jede Schwellen-Senkung (Coverage, Linter-Strenge, Prüfregel) ist ein
ADR, kein PR-Kommentar.

## 4. Quality Gates

Nur hier gelistete Targets existieren im Makefile. Halluzinierte
Gates sind die häufigste Form von Harness-Lüge.

| Target                       | Zweck                                                                                                                                                                                                                                              |
| ---------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `make lint`                  | golangci-lint mit dem Projekt-Profil (§3.2)                                                                                                                                                                                                        |
| `make test`                  | `go test ./...` — Akzeptanzkriterien der `DC-FA-*`; [`DC-QA-03`](spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Netzlos-Modullisten-Integrität der `.d-check.yml` (Go-Test, [ADR-0032](docs/plan/adr/0032-gate-consistency-tombstone.md))                                                                                                                                                                                                 |
| `make arch-check`            | Import-Regeln des Hexagon-Schnitts + Kern-Paket-Richtung **via digest-gepinntes a-check-Image** (Schwester-Tool, `a-check.mk` + `.a-check.yml`, netzlos/read-only) ([ADR-0005](docs/plan/adr/0005-modul-layout-hexagon-ordner.md), [ADR-0012](docs/plan/adr/0012-kern-paketschnitt-model-rules-app.md), [ADR-0029](docs/plan/adr/0029-arch-check-via-a-check.md) löst die Skript-/Stage-Mechanik ab)                                                      |
| `make coverage-gate`         | Coverage-Schwelle über `./internal/...` (Kalibrierungs-Bindung, siehe [`harness/README.md`](harness/README.md) §Sensors)                                                                                                                           |
| `make gate-consistency`      | Meta-Gate: Deklarations-Konsistenz Doku↔Makefile via Modul `targets` (Image, dogfood; [ADR-0031](docs/plan/adr/0031-targets-deklarations-konsistenz-modul.md); [ADR-0032](docs/plan/adr/0032-gate-consistency-tombstone.md) löst das Rest-Skript voll ab). Die [`DC-QA-03`](spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Modullisten-Integrität prüft jetzt ein getippter Go-Test in `make test`                                                                           |
| `make planning-check`        | Meta-Gate **via Modul `planning`** (Image, dogfood): Roadmap §Aktuelle Welle ↔ `in-progress/slice-*` (`planning-drift`, hermetisch — kein git, in `gates`) ([ADR-0028](docs/plan/adr/0028-planning-lifecycle-modul.md) löst die Skript-Mechanik von [slice-040](docs/plan/planning/done/slice-040-planning-consistency-gate.md) ab; [`DC-FA-PLAN-001`](spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)) |
| `make doc-check`             | Doku-Links, Anker, Kennungs-Linkpflicht, Referenzmatrix + Inline-Code-Pfade via `d-check` selbst (Dogfooding; netzlos — zugleich [`DC-QA-03`](spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Messmethode)             |
| `make gates`                 | alle inneren Gates (mandatory vor Handoff)                                                                                                                                                                                                         |
| `make ci`                    | CI-äquivalenter Lauf: gates + image-test (fährt die Release-Pipeline)                                                                                                                                                                              |
| `make trace-check`           | Traceability-Gate **via Modul `commits`** (Image, dogfood): DC-/ADR-/MR-/slice-ID in Commit-Messages (`commit-untraceable`; `RANGE=`-Range für CI, `MSGFILE=` für den `commit-msg`-Hook via stdin; bewusst **nicht** Teil von `gates`/`ci`) ([ADR-0027](docs/plan/adr/0027-commits-traceability-modul.md) löst die Skript-Mechanik von [ADR-0013](docs/plan/adr/0013-pr-ci-und-traceability-gate.md) ab; [`DC-FA-COMMITS-001`](spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in)) |
| `make adr-check`             | ADR-Immutable-Gate **via Modul `vcs`** (Image, dogfood): `Accepted`-ADRs nicht inhaltlich ändern (`RANGE=`/`STAGED=`-Modi; `pre-commit`-Hook + PR-CI; bewusst **nicht** Teil von `gates`/`ci`) ([ADR-0024](docs/plan/adr/0024-vcs-immutable-gate.md) löst die Skript-Mechanik von [ADR-0016](docs/plan/adr/0016-adr-immutable-gate.md) ab, [ADR-0025](docs/plan/adr/0025-codepaths-ignore-refs.md) entfernt das Alt-Skript)                                                            |
| `make hooks`                 | git-Hooks installieren (`core.hooksPath` → `.githooks`; aktiviert `commit-msg`-Traceability + `pre-commit`-ADR-Immutable via Modul `vcs`) ([ADR-0013](docs/plan/adr/0013-pr-ci-und-traceability-gate.md), [ADR-0016](docs/plan/adr/0016-adr-immutable-gate.md), [ADR-0024](docs/plan/adr/0024-vcs-immutable-gate.md))    |
| `make completeness-check`    | Requirements-Completeness-Gate **via in-Produkt-Flag** `--trace --require-complete` (≥1 Waise ⇒ Exit 1, mit `WAISE`-Zeilen + Anzahl); **Closure-Bindepunkt** (in `make fullbuild`, **nicht** `gates`/`ci`) ([ADR-0026](docs/plan/adr/0026-completeness-in-product-gate.md) löst die Skript-Mechanik von [ADR-0017](docs/plan/adr/0017-requirements-completeness-gate.md) ab; [`DC-FA-CLI-011`](spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)) |
| `make verify-closure-notes`  | Struktur des `done/`-Bestands: die Closure-Notizen **via Modul `planning`** (Abschnitt vorhanden, Substanz außerhalb Code, keine deklarierte Floskel, opt-in kein Vorlagen-Platzhalter) **und** Abschnitts-Invarianten **via Modul `structure`** (`section-*`-Codes) — beide über dasselbe `--config`-Profil. Fährt ein **eigenes** Prüf-Profil über `--config` ([`DC-FA-CLI-012`](spec/lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben)); **Closure-Bindepunkt** (in `make fullbuild`, bewusst **nicht** `gates`/`ci`) ([ADR-0048](docs/plan/adr/0048-closure-note-struktur-im-planning-modul.md), [`DC-FA-PLAN-001`](spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)) |
| `make fullbuild`             | volle Closure: gates + image-test + bench + completeness-check + verify-closure-notes, schließt mit dem Image-Hash                                                                                                                                                                             |
| `make image-test`            | [`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image)-Akzeptanzkriterien gegen das lokale Image (nativ vs. Container)                                                                                                                |
| `make bench`                 | [`DC-QA-01`](spec/lastenheft.md#dc-qa-01--performance)-Benchmark gegen generiertes Fixture (Median aus 3 Läufen, kein Gate in `gates`)                                                                                                             |
| `make trace`                 | Requirements Traceability Matrix via `d-check` selbst auf stdout (`--trace`, Dogfooding; netzlos, **kein Gate** — informativ) ([`DC-FA-CLI-009`](spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)) |
| `make doc-complete`          | Vollständigkeits-Dogfood via `d-check` selbst (`--trace --require-complete`, Requirements-Waise ⇒ Exit 1; Dogfooding, netzlos) — **kein** Gate-Bindepunkt; die Closure-Wahrheit bleibt `make completeness-check` ([`DC-FA-CLI-011`](spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code), [ADR-0017](docs/plan/adr/0017-requirements-completeness-gate.md))                              |
| `make semgrep`               | Security-/Static-Analysis-**Gate**: gepinntes semgrep-Image + gepinntes, lokal gecachtes `go/lang/security`-Regelset, netzloser Scan (`--network none`); **Bestandteil von `gates`** ([ADR-0010](docs/plan/adr/0010-semgrep-hermetisches-gate.md)) |
| `make versions`              | Reproduzierbarkeits-Pins ausgeben (Go, Lint, Basis-Images, Runtime-Image-ID)                                                                                                                                                                       |
| `make build` / `make run`    | Runtime-Image bauen / Selbst-Smoke-Test                                                                                                                                                                                                            |
| `make tidy`                  | `go.mod`/`go.sum` pflegen (`go mod tidy` in Docker; Dependency-Aufnahme/-Hebung — **kein** Gate, bewusster Akt am Dependency-Stand) |
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
  Anlege-Prozess (Akzeptanzkriterien-Trio, Versions-Bump + Historie,
  Beleg-Pflicht) folgt dem Baseline-Regelwerk
  ([`modul-03-spec`](.harness/baseline/v5.6.0/regelwerk/modul-03-spec.md)); das
  repo-spezifische ID-Schema steht in `spec/lastenheft.md` §3.
- Neue ADRs müssen den ADR-Index aktualisieren.
- Neue ADRs tragen die Sektion `## Re-Evaluierungs-Trigger` (oder „permanent");
  die vor Einführung `Accepted`-ADRs sind immutable und **grandfathered** (das
  Trigger-Feld liegt im ADR-Core, nachträgliches Ergänzen bräche `make adr-check`).
  Der Welle-Closure-Trigger-Audit (Baseline-Regelwerk Modul 6) bestätigt oder
  revidiert sie (Folge-ADR mit `supersedes`).
- Roadmap/Status-Geschichte lebt in `docs/plan/planning/`, nicht in der Architektur-Spec.
- Slice-Lifecycle (`open → next → in-progress → done`) ist reine Datei-Bewegung (`git mv`, siehe §3.3).
- Slice-Pläne tragen **kein** `**Status:**`-Feld — der Lifecycle-Zustand **ist** die
  Verzeichnis-Position; neue Slices führen stattdessen den `**Lifecycle:**`-Hinweis
  (Baseline-`slice.template.md`). Alt-Slices in `done/` behalten ihr historisches Feld.
- Jeder Slice-Plan trägt **vor** der Sub-Area-Modus-Begründung die zwei
  **Vorprüfungen** (Sub-Area prüfen · offene Beobachtungen im Register
  `observations.md` sichten) — Baseline-Regelwerk Modul 5/6, unabhängig vom
  Sub-Area-Modus.
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
