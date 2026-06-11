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

## Source precedence

| Rang | Datei | Charakter |
|---|---|---|
| 1 | [`spec/lastenheft.md`](../spec/lastenheft.md) | vertraglich abnahmebindend |
| 2 | [`spec/spezifikation.md`](../spec/spezifikation.md) | technisch fortschreibbar |
| 3 | [`spec/architecture.md`](../spec/architecture.md) | Komponenten/Sequenzen, meilensteinfrei |
| 4 | [`docs/plan/adr/`](../docs/plan/adr/) | Architekturentscheidungen |
| 5 | [`docs/plan/planning/in-progress/roadmap.md`](../docs/plan/planning/in-progress/roadmap.md) | aktuelle Welle |
| 6 | [`docs/user/`](../docs/user/) | Operations, Releasing (seit slice-011 — löst [`MR-009`](conventions.md#mr-009--source-precedence-ohne-docsuser-rang) auf) |
| 7 | [`README.md`](../README.md) | Projekt-Überblick |
| 8 | [`AGENTS.md`](../AGENTS.md) | Agent-Briefing |
| 9 | diese Datei | Harness-Einstieg |

## Guides (Feedforward-Quellen)

| Quelle | Inhalt |
|---|---|
| [`spec/lastenheft.md`](../spec/lastenheft.md) | Anforderungen (`DC-FA-*`, `DC-QA-*`), Akzeptanzkriterien |
| [`spec/spezifikation.md`](../spec/spezifikation.md) | Algorithmen, Schemas (`--json`, `.d-check.yml`), Defaults, Grund-Codes |
| [`spec/architecture.md`](../spec/architecture.md) | Hexagon-Schnitt (Rollen), Zugriffs-Constraints, Sequenzen |
| [`docs/plan/adr/`](../docs/plan/adr/) | Architekturentscheidungen |
| [`docs/plan/planning/`](../docs/plan/planning/) | Slice-Pläne und Roadmap |
| [`AGENTS.md`](../AGENTS.md) | Hard Rules, Source Precedence, Workflow |
| [`conventions.md`](conventions.md) | repo-lokale Strukturregeln, Adaptions-Block (`MR-*`), Modus-Deklarationen |
| [`agents-digest.md`](https://raw.githubusercontent.com/pt9912/ai-harness-course/main/kurs/de/agents-digest.md) | adoptiertes Betriebsregelwerk der Baseline in Agenten-Kurzform; derivativ — Stand siehe [`conventions.md` §Baseline](conventions.md#baseline) |
| [`.harness/skills/reviewer.md`](../.harness/skills/reviewer.md) | Reviewer-Skill: Kategorien-Anker, Output-Schema, Negativbefund-Pflicht; ein Report pro Lauf unter [`docs/reviews/`](../docs/reviews/) |

## Sensors (Feedback-Gates)

Nur Targets, die im Makefile **existieren**. Lauf-Wahrheit pro Commit
liegt in CI bzw. lokal (`make gates`), nicht hier.

| Target | Vertrag | Bindung |
|---|---|---|
| `make lint` | golangci-lint, SOLID-nahes Profil (5 Default- + 23 Linter, kalibriert; Ausnahmen zentral mit Why; Inline-Suppressions verboten) | [ADR-0006](../docs/plan/adr/0006-lint-profil-solid.md) |
| `make test` | Akzeptanzkriterien der bezogenen `DC-FA-*` als Tests; Determinismus-Test | [`DC-QA-02`](../spec/lastenheft.md#dc-qa-02--determinismus) (DC-Bindung) |
| `make arch-check` | Import-Regeln R1–R5 des Hexagon-Schnitts (`tools/arch-check.sh`, Dockerfile-Stage) | [ADR-0005](../docs/plan/adr/0005-modul-layout-hexagon-ordner.md); [`DC-QA-03`](../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) (DC-Bindung) |
| `make doc-check` | Links, **Anker**, **Kennungs-Linkpflicht** und **Referenzmatrix** der gesamten Repo-Doku via `d-check` selbst (Dogfooding-Selbstkonfiguration: Runtime-Image, read-only-Mount, `--network none`, [`.d-check.yml`](../.d-check.yml)) — zugleich die automatisierte `DC-QA-03`-Messmethode (netzloser Lauf aller Module außer `external`) | [`DC-FA-LINK-001`](../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)/[`DC-FA-ANCH-001`](../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors)/[`DC-FA-ID-001`](../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)/[`DC-FA-MTX-001`](../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)/[`DC-QA-03`](../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) (DC-Bindung); [`MR-006`](conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs) (maschinell kodiert), [`MR-007`](conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding) |
| `make coverage-gate` | Gesamt-Coverage über `./internal/...` (`-coverpkg` über Paketgrenzen, Dockerfile-Stage + `tools/coverage-gate.sh`) | Kalibrierungs-Bindung: Schwelle 93 % seit 2026-06-11 (Kalibrierung nach Test-Ausbau, Ist 95,1 %; zuvor Ramp 85 % → 90 % bei welle-03 done, slice-009); Verfehlung ⇒ Carveout-Pflicht (Kurs-Modul 13); Senkung nur per ADR (`AGENTS.md` §3.6) |
| `make gate-consistency` | Meta-Gate (`tools/gate-consistency.sh`): in `AGENTS.md` §4 / dieser Tabelle dokumentierte Targets ↔ Makefile (beide Richtungen) + `DC-QA-03`-Modulliste der [`.d-check.yml`](../.d-check.yml); Selbsttest mit Phantom-Target bei jedem Lauf | [`DC-QA-03`](../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) (DC-Bindung) |
| `make record-gates` | Working-Tree-Hash-Nachweis für den Stop-Hook | — |
| `make gates` | aggregiert doc-check + lint + test + arch-check + coverage-gate + gate-consistency, `record-gates` als letzter Schritt | — |
| `make image-test` | [`DC-FA-DIST-001`](../spec/lastenheft.md#dc-fa-dist-001--docker-image)-Akzeptanzkriterien gegen das lokal gebaute Image (`tools/image-test.sh`): Befund-Ausgabe und Exit-Code nativ vs. Container **byte-identisch**, read-only-Mount vollständig, fehlender Mount → Exit 2 mit Hinweis | [`DC-FA-DIST-001`](../spec/lastenheft.md#dc-fa-dist-001--docker-image)/[`DC-QA-02`](../spec/lastenheft.md#dc-qa-02--determinismus) (DC-Bindung) |
| `make ci` | CI-äquivalenter Lauf (gates + image-test) — das Target der Release-Pipeline (slice-011) | — |
| `make fullbuild` | volle Closure vor Welle-Merge/Release (gates + image-test + bench); schließt mit dem Image-Hash des Runtime-Builds ab | Reproduzierbarkeits-Bindung: Image-Hash (`sha256:…`) im Lauf-Abschluss; Pins via `make versions` (Kurs-Modul 14) |
| `make versions` | Reproduzierbarkeits-Pins: `GO_VERSION`, `GOLANGCI_LINT_VERSION`, alle `FROM`-Basis-Images, Runtime-Image-ID | — |

**Aktueller Lauf-Status:** lokal `make gates`; Releases laufen über
[`release.yml`](../.github/workflows/release.yml) (Tag-Push `v*` →
`make ci` → GHCR-Push mit Digest-Pin).
**Rote Gates:** keine.
**Nicht behauptet:** — keine — (alle geplanten Targets existieren;
`make gate-consistency` bewacht die Tabelle in beide Richtungen).

## Traceability rules

- PRs/Commits **müssen** mindestens eine `DC-*`- oder `ADR-*`-ID nennen.
- Neue oder geänderte Anforderungen brauchen einen Beleg: Test, Gate, Demo oder ADR.
- Neue ADRs müssen im [ADR-Index](../docs/plan/adr/README.md) ergänzt werden.
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
