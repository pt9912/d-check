# Harness

## Purpose

Dieser Harness verbindet bestehende Spezifikationen, ADRs,
Planning-Dokumente und Gates. Er ist **kein Ersatz** für `spec/` oder
`docs/`, sondern ein **Einstiegspunkt** für Menschen und AI-Code-Agenten.

Wenn diese Datei einer kanonischen Quelle widerspricht, **gewinnt die
kanonische Quelle**, und diese Datei wird angepasst.

Strukturregeln (Verzeichniskonvention, ID-Schemata, Modus-Deklarationen
pro Sub-Area) sowie Adaptionen ggü. der adoptierten Baseline leben in
[`conventions.md`](conventions.md). Diese Datei dupliziert sie nicht.

## Source precedence

| Rang | Datei | Charakter |
|---|---|---|
| 1 | [`spec/lastenheft.md`](../spec/lastenheft.md) | vertraglich abnahmebindend |
| 2 | `spec/spezifikation.md` | technisch fortschreibbar — **geplant** (slice-002) |
| 3 | `spec/architecture.md` | Komponenten/Sequenzen, meilensteinfrei — **geplant** (slice-002) |
| 4 | [`docs/plan/adr/`](../docs/plan/adr/) | Architekturentscheidungen |
| 5 | [`docs/plan/planning/in-progress/roadmap.md`](../docs/plan/planning/in-progress/roadmap.md) | aktuelle Welle |
| 6 | [`README.md`](../README.md) | Projekt-Überblick |
| 7 | [`AGENTS.md`](../AGENTS.md) | Agent-Briefing |
| 8 | diese Datei | Harness-Einstieg |

## Guides (Feedforward-Quellen)

| Quelle | Inhalt |
|---|---|
| [`spec/lastenheft.md`](../spec/lastenheft.md) | Anforderungen (`DC-FA-*`, `DC-QA-*`), Akzeptanzkriterien |
| [`docs/plan/adr/`](../docs/plan/adr/) | Architekturentscheidungen |
| [`docs/plan/planning/`](../docs/plan/planning/) | Slice-Pläne und Roadmap |
| [`AGENTS.md`](../AGENTS.md) | Hard Rules, Source Precedence, Workflow |
| [`conventions.md`](conventions.md) | repo-lokale Strukturregeln, Adaptions-Block (`MR-*`), Modus-Deklarationen |

## Sensors (Feedback-Gates)

Nur Targets, die im Makefile **existieren**. Lauf-Wahrheit pro Commit
liegt in CI bzw. lokal (`make gates`), nicht hier.

| Target | Vertrag | Bindung |
|---|---|---|
| `make doc-check` | jedes lokale Markdown-Linkziel in `docs/`, `spec/`, `harness/` und den Top-Level-Dokumenten existiert | Bootstrap-Sensor (vendored, siehe [`MR-003`](conventions.md#mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh)); Ablösung: slice-004 |
| `make record-gates` | Working-Tree-Hash-Nachweis für den Stop-Hook | — |
| `make gates` | aggregiert alle inneren Gates, `record-gates` als letzter Schritt | — |

**Aktueller Lauf-Status:** lokal `make gates`.
**Rote Gates:** keine (Bootstrap-Stand).
**Nicht behauptet** (geplant): `make lint`, `make test`,
`make arch-check`, `make coverage-gate` — entstehen ab slice-003.

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
