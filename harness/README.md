# Harness

## Purpose

Dieser Harness verbindet bestehende Spezifikationen, ADRs,
Planning-Dokumente und Gates. Er ist **kein Ersatz** für `spec/` oder
`docs/`, sondern ein **Einstiegspunkt** für Menschen und AI-Code-Agenten.

Wenn diese Datei einer kanonischen Quelle widerspricht, **gewinnt die
kanonische Quelle**, und diese Datei wird angepasst.

Strukturregeln (Verzeichniskonvention, ID-Schemata, Modus-Deklarationen
pro Sub-Area, Zusatzklassen für Sensors-Bindung) sowie Adaptionen ggü.
der adoptierten Baseline leben in [`conventions.md`](conventions.md).
Diese Datei dupliziert sie nicht.

## Leseordnung

Die Menschen-Hälfte des Einstiegs (Baseline v5.5.0) — was zuerst, was bei
Bedarf; die Sektionen darunter sind Referenzfläche zum Nachschlagen.
Bewusste Auslegung: die Ordnung steht hier **vorn** statt wie im
Baseline-Skelett am Ende — der Einstieg gehört an den Anfang:

1. [`AGENTS.md`](../AGENTS.md) — Hard Rules und Workflow (jeder Lauf).
2. [`harness/conventions.md`](conventions.md) — Baseline-Stand und aktive
   Adaptionen (vor jeder Doku-/Konventions-Änderung).
3. Der aktive Slice unter
   [`docs/plan/planning/in-progress/`](../docs/plan/planning/in-progress/)
   samt [Roadmap](../docs/plan/planning/in-progress/roadmap.md) — woran
   gerade gearbeitet wird.
4. Bei Bedarf: das vendorte Regelwerk
   ([Index](../.harness/baseline/v6.3.1/regelwerk/README.md)) — nur den
   benötigten Abschnitt.

## Source precedence

| Rang | Datei                                                                                       | Charakter                                                                                                                 |
| ---- | ------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| 1    | [`spec/lastenheft.md`](../spec/lastenheft.md)                                               | vertraglich abnahmebindend                                                                                                |
| 2    | [`spec/spezifikation.md`](../spec/spezifikation.md)                                         | technisch fortschreibbar                                                                                                  |
| 3    | [`spec/architecture.md`](../spec/architecture.md)                                           | Komponenten/Sequenzen, meilensteinfrei                                                                                    |
| 4    | [`docs/plan/adr/`](../docs/plan/adr/)                                                       | Architekturentscheidungen                                                                                                 |
| 5    | [`docs/plan/planning/in-progress/roadmap.md`](../docs/plan/planning/in-progress/roadmap.md) | Wellen-Sequenzierung (offene Wellen derivativ)                                                                                                            |
| 6    | [`docs/user/`](../docs/user/)                                                               | Operations, Releasing (löst [`MR-009`](conventions.md#mr-009--source-precedence-ohne-docsuser-rang) auf) |
| 7    | [`README.md`](../README.md)                                                                 | Projekt-Überblick                                                                                                         |
| 8    | [`AGENTS.md`](../AGENTS.md)                                                                 | Agent-Briefing                                                                                                            |
| 9    | diese Datei                                                                                 | Harness-Einstieg                                                                                                          |

## Guides (Feedforward-Quellen)

| Quelle                                                                                                                 | Inhalt                                                                                                                                                                                                                                                                                                                                     |
| ---------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [`spec/lastenheft.md`](../spec/lastenheft.md)                                                                          | Anforderungen (`DC-FA-*`, `DC-QA-*`), Akzeptanzkriterien                                                                                                                                                                                                                                                                                   |
| [`spec/spezifikation.md`](../spec/spezifikation.md)                                                                    | Algorithmen, Schemas (`--json`, `.d-check.yml`), Defaults, Grund-Codes                                                                                                                                                                                                                                                                     |
| [`spec/architecture.md`](../spec/architecture.md)                                                                      | Hexagon-Schnitt (Rollen), Zugriffs-Constraints, Sequenzen                                                                                                                                                                                                                                                                                  |
| [`docs/plan/adr/`](../docs/plan/adr/)                                                                                  | Architekturentscheidungen                                                                                                                                                                                                                                                                                                                  |
| [`docs/plan/planning/`](../docs/plan/planning/)                                                                        | Slice-Pläne und Roadmap                                                                                                                                                                                                                                                                                                                    |
| [`AGENTS.md`](../AGENTS.md)                                                                                            | Hard Rules, Source Precedence, Workflow                                                                                                                                                                                                                                                                                                    |
| [`conventions.md`](conventions.md)                                                                                     | repo-lokale Strukturregeln, Adaptions-Block (`MR-*`), Modus-Deklarationen                                                                                                                                                                                                                                                                  |
| [`.harness/baseline/v6.3.1/regelwerk/`](../.harness/baseline/v6.3.1/regelwerk/) | adoptiertes Betriebsregelwerk der Baseline (committet vendored, netzlos; die dortige `README.md` ist der Index), pro Session **nur den benötigten Abschnitt** lesen — nach Modulen + Grundlagen aufgeteilt, vendored aus dem self-contained [`lab-regelwerk.zip`](https://github.com/pt9912/ai-harness-course/releases/download/v6.3.1/lab-regelwerk.zip) ([`MR-019`](conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017), Layout [`MR-023`](conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)); das frühere separate `agents-regelwerk.md` ist im Kurs abgelöst; derivativ — Stand siehe [`conventions.md` §Baseline](conventions.md#baseline) |
| [`.harness/skills/reviewer.md`](../.harness/skills/reviewer.md)                                                        | Reviewer-Skill: Kategorien-Anker, Output-Schema, Negativbefund-Pflicht; ein Report pro Lauf unter [`docs/reviews/`](../docs/reviews/)                                                                                                                                                                                                      |
| [`.harness/skills/closure-note-reviewer.md`](../.harness/skills/closure-note-reviewer.md)                              | Closure-Note-Reviewer-Skill: die **semantische** Schicht über dem strukturellen `make verify-closure-notes` — prüft *Inhalt vs. Floskel* (Lernsignal · Folge-Slice · Architektur-Beobachtung) und meldet ausdrücklich **nicht** doppelt, was das Gate bereits abdeckt |

## Sensors (Feedback-Gates)

Nur Targets, die im Makefile **existieren**. Lauf-Wahrheit pro Commit
liegt in CI bzw. lokal (`make gates`), nicht hier.

| Target | Vertrag | Bindung |
| --- | --- | --- |
| [`make lint`](sensors/lint.md) | fährt das SOLID-nahe golangci-lint-Profil dieses Repos | [ADR-0006](../docs/plan/adr/0006-lint-profil-solid.md) (Profil); [`AGENTS.md`](../AGENTS.md) §3.2 (`nolintlint`) |
| [`make test`](sensors/test.md) | fährt die Go-Testsuite des Hauptmoduls, inklusive zweier Zusagen des Repos über sich selbst | [`DC-QA-02`](../spec/lastenheft.md#dc-qa-02--determinismus)/[`DC-QA-03`](../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) |
| [`make arch-check`](sensors/arch-check.md) | hält die Import-Regeln des Hexagon-Schnitts über das Schwester-Tool a-check | [ADR-0005](../docs/plan/adr/0005-modul-layout-hexagon-ordner.md), [ADR-0012](../docs/plan/adr/0012-kern-paketschnitt-model-rules-app.md), [ADR-0029](../docs/plan/adr/0029-arch-check-via-a-check.md) (löst die Skript-/Stage-Mechanik ab); [`DC-QA-03`](../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) (DC-Bindung) |
| [`make doc-check`](sensors/doc-check.md) | prüft die gesamte Repo-Doku mit d-check selbst (Dogfooding) | [`DC-FA-LINK-001`](../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)/[`DC-FA-ANCH-001`](../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors)/[`DC-FA-ID-001`](../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)/[`DC-FA-MTX-001`](../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)/[`DC-FA-CODE-001`](../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)/[`DC-QA-03`](../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) (DC-Bindung); [`MR-006`](conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs) (maschinell kodiert), [`MR-007`](conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding) |
| `make coverage-gate` | Gesamt-Coverage über `internal/...` (`-coverpkg` über Paketgrenzen, Dockerfile-Stage + `tools/coverage-gate.sh`) | Kalibrierungs-Bindung: Schwelle 93 % seit 2026-06-11 (Kalibrierung nach Test-Ausbau, Ist 95,1 %); Verfehlung ⇒ Carveout-Pflicht ([Kurs-Modul 13](../.harness/baseline/v6.3.1/regelwerk/modul-13-quality-gates.md)); Senkung nur per ADR (`AGENTS.md` §3.6) |
| [`make semgrep`](sensors/semgrep.md) | hermetisches Security-/Static-Analysis-Gate über den Go-Code | [ADR-0010](../docs/plan/adr/0010-semgrep-hermetisches-gate.md); [`DC-QA-02`](../spec/lastenheft.md#dc-qa-02--determinismus)/[`DC-QA-03`](../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) (DC-Bindung) |
| [`make image-scan`](sensors/image-scan.md) | CVE-Scan gegen die publizierten Images, nicht gegen den Arbeitsbaum | [ADR-0066](../docs/plan/adr/0066-cve-scan-gegen-das-publizierte-image.md) |
| [`make gate-consistency`](sensors/gate-consistency.md) | hält Doku und Makefile über die Gate-Targets in beiden Richtungen deckungsgleich | [`DC-FA-TGT-001`](../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)/[ADR-0031](../docs/plan/adr/0031-targets-deklarations-konsistenz-modul.md) (DC-Bindung) |
| [`make planning-check`](sensors/planning-check.md) | hält den Ruhe-Marker der Roadmap gegen das Lifecycle-Verzeichnis | [ADR-0028](../docs/plan/adr/0028-planning-lifecycle-modul.md) (löst die Skript-Mechanik ab); [`DC-FA-PLAN-001`](../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) |
| [`make workflow-pins`](sensors/workflow-pins.md) | hält die Deklarations-Form aller `uses:`-Referenzen in den Workflows | [ADR-0072](../docs/plan/adr/0072-workflows-modul.md); [`DC-FA-WF-001`](../spec/lastenheft.md#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in) |
| [`make review-coverage`](sensors/review-coverage.md) | hält, dass ein geschlossener Slice mit Review-Zusage auch einen Report hat | [ADR-0081](../docs/plan/adr/0081-reviews-modul.md); [`DC-FA-RVW-001`](../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in) |
| [`make baseline-verify`](sensors/baseline-verify.md) | prüft den committeten vendorten Baseline-Bestand auf Unversehrtheit | [`MR-011`](conventions.md#mr-011)-Kette; [`MR-021`](conventions.md#mr-021) (pin-gebundene Verweise); [`MR-055`](conventions.md#mr-055) (der Symlink als Träger) |
| `make gates` | aggregiert baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check, `record-gates` als letzter Schritt | — |
| [`make image-test`](sensors/image-test.md) | prüft die Distributions-Akzeptanzkriterien gegen das lokal gebaute Image | [`DC-FA-DIST-001`](../spec/lastenheft.md#dc-fa-dist-001--docker-image)/[`DC-QA-02`](../spec/lastenheft.md#dc-qa-02--determinismus) (DC-Bindung) |
| `make ci` | CI-äquivalenter Lauf (gates + image-test) — das Target der Release-Pipeline | — |
| [`make trace-check`](sensors/trace-check.md) | hält, dass jede Commit-Botschaft eine Traceability-Kennung nennt | [ADR-0027](../docs/plan/adr/0027-commits-traceability-modul.md) (löst die Skript-Mechanik von [ADR-0013](../docs/plan/adr/0013-pr-ci-und-traceability-gate.md) ab); [`DC-FA-COMMITS-001`](../spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in) |
| [`make adr-check`](sensors/adr-check.md) | hält, dass eine `Accepted`-ADR nicht inhaltlich überschrieben wird | [ADR-0024](../docs/plan/adr/0024-vcs-immutable-gate.md) (löst die Skript-Mechanik von [ADR-0016](../docs/plan/adr/0016-adr-immutable-gate.md) ab); [ADR-0025](../docs/plan/adr/0025-codepaths-ignore-refs.md) (entfernt das Alt-Skript); [`AGENTS.md` §3.5](../AGENTS.md#35-adrs-sind-nach-accepted-immutable) |
| [`make completeness-check`](sensors/completeness-check.md) | meldet Anforderungen ohne referenzierenden Slice am Closure-Punkt | [ADR-0026](../docs/plan/adr/0026-completeness-in-product-gate.md) (löst die Skript-Mechanik von [ADR-0017](../docs/plan/adr/0017-requirements-completeness-gate.md) ab); [`DC-FA-CLI-011`](../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code) |
| [`make verify-closure-notes`](sensors/verify-closure-notes.md) | hält die Struktur des `done/`-Bestands am Closure-Bindepunkt | [ADR-0048](../docs/plan/adr/0048-closure-note-struktur-im-planning-modul.md); [`DC-FA-PLAN-001`](../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) |
| `make fullbuild` | volle Closure vor Welle-Merge/Release (gates + image-test + bench + completeness-check + verify-closure-notes); schließt mit dem Image-Hash des Runtime-Builds ab | Reproduzierbarkeits-Bindung: Image-Hash (`sha256:…`) im Lauf-Abschluss; Pins via `make versions` ([Kurs-Modul 14](../.harness/baseline/v6.3.1/regelwerk/modul-14-docker-harness.md)) |


**Werkzeuge — genannt, weil der Lauf sie braucht, aber kein Gate.** Das
Kriterium ist nicht, ob ein Target in `gates` läuft, sondern **worüber es
urteilt**: Ein Gate prüft den Zustand des Repos, ein Werkzeug die
Vorbedingungen seines eigenen Laufs — es *bewegt*, *misst*, *sagt* oder
*meldet Fremdes*. Das schärfste Unterscheidungsmerkmal ist **fail-open**:
was im Zweifel durchlässt, urteilt nicht (Baseline-Regelwerk
`modul-13-quality-gates.md` §Hard Rule, *Die dritte Lage*). Weglassen wäre
nur für ein Target richtig, das niemand braucht; unmarkiert in der
Gate-Tabelle stehen hieße, es als Gate zu behaupten.

| Target | Tut was | Bindung |
| --- | --- | --- |
| `make freshness-trivy` · `make trivy-digest` | Scanner-Pin auf beiden Achsen (Version, Digest). **Netz**, fail-open, Nachtlauf | kein Gate |
| `.github/dependabot.yml` | **kein Target** — der Kanal, der hebt, was `image-scan` meldet; `gomod` und `github-actions`, **nicht** `docker` | kein Gate · [ADR-0067](../docs/plan/adr/0067-dependabot-als-hebender-kanal.md) |
| [`make freshness-go`](sensors/freshness-go.md) · `make freshness-golangci` · `make freshness-semgrep` · `make freshness-a-check` | meldet, ob upstream ein neuerer Release existiert als der gepinnte | kein Gate · [`AGENTS.md`](../AGENTS.md) §4; Nutzer-Regel: Go-Bump zieht das `golangci`-Pendant nach |
| [`make checkout-pin-freshness`](sensors/checkout-pin-freshness.md) · `make login-pin-freshness` · `make hubdesc-pin-freshness` | meldet, ob die drei Action-Pins der Workflows veraltet sind | kein Gate · [`AGENTS.md`](../AGENTS.md) §3.9, §4 |
| [`make runtime-base-digest`](sensors/runtime-base-digest.md) · `make go-base-digest` · `make lint-base-digest` · `make semgrep-digest` · `make a-check-digest` | meldet, ob ein digest-gepinntes Fremd-Image unter demselben Tag neu gebaut wurde | kein Gate · [ADR-0011](../docs/plan/adr/0011-digest-pins-build-gate-images.md), [`AGENTS.md`](../AGENTS.md) §4 |
| [`make baseline-freshness`](sensors/baseline-freshness.md) | auditiert den Baseline-Pin gegen upstream, in zwei getrennten Teilen | kein Gate · [`MR-011`](conventions.md#mr-011)-Kette. **Nicht** [`DC-QA-03`](../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit): jene Zusage gilt dem **Produkt**, nicht dem inneren Lauf — dass `gates` netzlos bleibt, ist eine Eigenschaft dieses Repos |
| [`make nightly-state`](sensors/nightly-state.md) | liest den Ausgang der beiden Nachtläufe und sagt, ob er gelesen werden muss | kein Gate · [`MR-053`](conventions.md#mr-053); [`AGENTS.md`](../AGENTS.md) §4, §5 |
| [`make guard-probe`](sensors/guard-probe.md) | fährt den Tool-Call-Wächter gegen seine Proben | kein Gate · [`MR-005`](conventions.md#mr-005), [`MR-040`](conventions.md#mr-040), [`MR-042`](conventions.md#mr-042), [`MR-044`](conventions.md#mr-044) |
| `make record-gates` | Working-Tree-Hash-Nachweis für den Stop-Hook | kein Gate |
| [`make hooks`](sensors/hooks.md) | installiert die lokalen git-Hooks, die Commit und Übergang an grüne Gates binden | kein Gate · [ADR-0013](../docs/plan/adr/0013-pr-ci-und-traceability-gate.md), [ADR-0024](../docs/plan/adr/0024-vcs-immutable-gate.md) |
| `make versions` | Reproduzierbarkeits-Pins: `GO_VERSION`, `GOLANGCI_LINT_VERSION`, alle `FROM`-Basis-Images, Runtime-Image-ID | kein Gate |
| `make bench` | misst die Performance gegen ein generiertes Fixture (Median aus drei Läufen) | kein Gate, [`DC-QA-01`](../spec/lastenheft.md#dc-qa-01--performance) |
| `make baseline-probe` | fährt die Alias-Auflösung von [`baseline-verify`](sensors/baseline-verify.md) gegen neun Proben | kein Gate · [`MR-055`](conventions.md#mr-055) |
| `make trace` | gibt die Requirements-Traceability-Matrix auf stdout aus | kein Gate · [`DC-FA-CLI-009`](../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix) |
| `make doc-complete` | sagt, welche Anforderungen ohne referenzierenden Slice sind — ohne Closure-Bindung | kein Gate · [`DC-FA-CLI-011`](../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code) |
| `make archive-wave` | bewegt geschlossene Zeitdokumente ins Archiv und ersetzt sie durch Stubs; ohne `APPLY=1` wird nichts geschrieben | kein Gate |
| `make archive-wave-test` | fährt die Testsuite von `tools/archive-wave/` (eigenes `go.mod`, nicht Teil von `make test`) | kein Gate |
| `make tidy` | pflegt `go.mod`/`go.sum` in Docker — bewusster Akt am Dependency-Stand | kein Gate |

**Aktueller Lauf-Status:** lokal `make gates`; Releases laufen über
[`release.yml`](../.github/workflows/release.yml) (Tag-Push `v*` →
`make ci` → GHCR-Push mit Digest-Pin).
**Rote Gates:** keine.
**Nicht behauptet:** — keine — (alle geplanten Targets existieren;
`make gate-consistency` bewacht die Tabelle in beide Richtungen).

### Gate-Taxonomie und Durchsetzungsgrenzen

Die Sensors zerfallen in zwei Klassen — ehrlich benannt, weil ihre
Durchsetzungskraft unterschiedlich ist. **Die Werkzeuge aus der zweiten
Tabelle oben stehen hier nicht**: Sie urteilen nicht über den Repo-Zustand und
sind deshalb keine Gates, auch keine Meta-Gates.

| Klasse | Was geprüft wird | Targets | Bindepunkt |
| --- | --- | --- | --- |
| **Produkt-Gates** | Eigenschaften des Arbeitsprodukts (Code, Doku, Image) | `lint`, `test`, `arch-check`, `coverage-gate`, `semgrep`, `doc-check`, `image-test` | Arbeitsbaum-Inhalt (in `make gates`/`ci`) |
| **Meta-/Governance-Gates** | Harness-Integrität & Prozess-Invarianten („keine Harness-Lüge") | `gate-consistency`, `planning-check`, `baseline-verify`, `workflow-pins` (in `gates`); `trace-check`, `adr-check` (Commit-/Diff-Bindepunkt, **nicht** in `gates`/`ci`); `review-coverage` (eigenständiger Fokus-Lauf, **nicht** in `gates`/`ci`); `completeness-check`, `verify-closure-notes` (**Closure-Bindepunkte**: in `fullbuild`, **nicht** `gates`/`ci`) | Doku↔Makefile, Roadmap↔Lifecycle, Commit↔ID, ADR-Immutability, Requirements-Waisen, Closure-Note-Substanz, Risiko-Ausgänge, Vendor-Integrität, `uses:`-Pin-Form, Review-Report-Deckung |

Meta-Gates sind **ausführbare, fail-closed Gates mit Negativ-Selbsttest** —
keine aspirativen Texte. Aber ihre Kraft ist real begrenzt:

- **Lokale Hooks** (`commit-msg`-Traceability, `pre-commit`-ADR-Immutable)
  sind **opt-in pro Klon** (`make hooks` setzt `core.hooksPath`); aus einem
  fremden Klon sind sie nicht erzwingbar.
- Der **klon-unabhängige Boden** ist die PR-/Push-CI
  ([`ci.yml`](../.github/workflows/ci.yml)): sie fährt `make ci` **und** die
  Range-Gates (`trace-check`, `adr-check`) auf jede Integration.
- Die CI **blockiert** einen Merge aber nur, wenn **Branch Protection**
  (Pflicht-Status-Checks auf dem Default-Branch) gesetzt ist. Das liegt
  **außerhalb des Repos** und ist aus dem Klon **nicht auditierbar** — ohne
  sie ist die CI nur *advisory*. **Betriebsempfehlung:** den `ci`-Check als
  *required status check* für den Default-Branch konfigurieren
  ([ADR-0013](../docs/plan/adr/0013-pr-ci-und-traceability-gate.md),
  [ADR-0016](../docs/plan/adr/0016-adr-immutable-gate.md) benennen dieselbe
  Restlücke).
- Der `Stop`-Hook
  ([`stop-require-gates.sh`](../.claude/hooks/stop-require-gates.sh)) ist
  Claude-spezifisch und gibt frische, cleane Klone ohne lokalen Gate-State
  frei — dafür ist die CI das Netz.
- Dasselbe gilt für das zweite Closure-Gate `verify-closure-notes`
  (Closure-Note-Struktur,
  [ADR-0048](../docs/plan/adr/0048-closure-note-struktur-im-planning-modul.md)):
  es hängt an `make fullbuild`, nicht an `gates`/`ci`, und prüft **Struktur, nicht
  Bedeutung** — eine grüne Ausgabe sagt „Form erfüllt", nicht „Notizen sind gut".
- Das **Closure-Gate** `completeness-check` (Requirements-Waisen,
  [ADR-0017](../docs/plan/adr/0017-requirements-completeness-gate.md)) hängt an
  `make fullbuild`, **nicht** an `gates`/`ci`/[`release.yml`](../.github/workflows/release.yml)
  — es ist der **manuelle** Welle-/Release-Abschluss-Check (GF erlaubt
  transiente Waisen im Inner-Loop) und wird von der CI bewusst **nicht**
  erzwungen; Disziplin liegt am Closure-Punkt.

## Traceability rules

- PRs/Commits **müssen** mindestens eine `DC-*`-, `ADR-*`-, `MR-*`- oder
  `slice-*`-ID nennen — maschinell erzwungen über `make trace-check`
  (lokaler `commit-msg`-Hook via `make hooks` + PR-/Push-CI; Ausnahme:
  Merge-/Revert-Commits;
  [ADR-0013](../docs/plan/adr/0013-pr-ci-und-traceability-gate.md)).
- Neue oder geänderte Anforderungen brauchen einen Beleg: Test, Gate, Demo oder ADR.
- Neue ADRs müssen im [ADR-Index](../docs/plan/adr/README.md) ergänzt werden.
- Eine `Accepted`-ADR wird **nicht inhaltlich überschrieben** (`AGENTS.md`
  §3.5) — maschinell erzwungen über `make adr-check` (`pre-commit`-Hook +
  PR-/Push-CI; erlaubt bleiben `## Geschichte`-Anhänge + der
  `**Status:**`-Übergang;
  [ADR-0016](../docs/plan/adr/0016-adr-immutable-gate.md)).
- Änderungen an Planning-Dokumenten folgen den Lifecycle-Regeln (`open → next → in-progress → done`; reine `git mv`-Commits, siehe `AGENTS.md` §3.3).

## Safety and scope boundaries

- `d-check` ist ein **Lese-Tool**: Es schreibt nie in das geprüfte
  Repository (Kernvertrag
  [`DC-QA-03`](../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- Netzwerkzugriffe nur im explizit aktivierten Modul `external` —
  niemals im Default.
- Determinismus ist Kernvertrag
  ([`DC-QA-02`](../spec/lastenheft.md#dc-qa-02--determinismus)):
  identische Eingabe ⇒ identische Ausgabe.
- Dieses Repo ist kein produktiver Service; das Produkt ist ein
  CLI-Tool/Container-Image.

## Minimal agent workflow

1. Diese Datei lesen.
2. Relevante kanonische Quelle lesen.
3. Betroffene IDs identifizieren.
4. Kleinste Änderung planen.
5. Engsten nützlichen Sensor laufen lassen.
6. Repo-weiten Gate-Lauf vor Handoff (`make gates`).
7. Doku/Indizes aktualisieren, falls ein öffentlicher Vertrag berührt.
8. Ausgeführte Sensors und verbleibende Risiken berichten.
