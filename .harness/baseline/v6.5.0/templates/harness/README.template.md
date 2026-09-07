# Harness

> **Template-Hinweis.** Diese Datei ist eine Vorlage für `harness/README.md`
> deines Repos. Kopiere sie nach `harness/README.md`, ersetze
> `<Platzhalter>` und lösche diesen Block. Pflichtgliederung folgt
> [Baseline-Regelwerk §harness/README.md als Einstiegspunkt](../../regelwerk/grundlagen-harness-dateien.md#harnessreadmemd-als-einstiegspunkt).
> **Pointer-Artefakt:** verweist auf andere kanonische Quellen — zuletzt
> füllen bzw. re-syncen, sobald die Ziele stehen; veraltete
> `(folgt)`/Klartext-Verweise fängt kein Linter (Reviewer-Sache).

---

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
| 5 | [`docs/plan/planning/in-progress/roadmap.md`](../docs/plan/planning/in-progress/roadmap.md) | Wellen-Sequenz |
| 6 | `docs/user/*` *(falls vorhanden)* | Operations, Quality, Releasing | <!-- d-check:ignore (Verzeichnis optional; entlinkt, da im frischen Repo selten vorhanden) -->
| 7 | [`README.md`](../README.md) | Projekt-Überblick |
| 8 | [`AGENTS.md`](../AGENTS.md) | Agent-Briefing |
| 9 | diese Datei | Harness-Einstieg |

> Die Ränge 1–3 sind die **drei Spec-Straten** — Vertrag, Technik, Sicht —,
> und alle drei sind obligatorisch (Baseline-Regelwerk
> `grundlagen-referenz-richtung.md` §Spec-Straten). **Adaption ist die
> Zwei-Straten-Form, nicht die Drei-Straten-Form**: Wer Rang 2 streicht,
> deklariert das als `MR-<NNN>` in [`conventions.md`](conventions.md) und
> nummeriert neu (dann acht Ränge).

## Guides (Feedforward-Quellen)

<!--
Was lenkt den Agenten *vor* der Handlung? Pointer, kein Inhalt.
-->

| Quelle | Inhalt |
|---|---|
| [`spec/lastenheft.md`](../spec/lastenheft.md) | Anforderungen, IDs, Akzeptanzkriterien |
| [`spec/spezifikation.md`](../spec/spezifikation.md) | technische Details, Defaults |
| [`spec/architecture.md`](../spec/architecture.md) | Komponenten, Schichten, Constraints |
| [`docs/plan/adr/`](../docs/plan/adr/) | Architekturentscheidungen |
| [`docs/plan/planning/`](../docs/plan/planning/) | Slice-Pläne und Roadmap |
| [`AGENTS.md`](../AGENTS.md) | Hard Rules, Source Precedence, Workflow |
| [`conventions.md`](conventions.md) | repo-lokale Strukturregeln, Adaptions-Block (`MR-*`), Modus-Deklarationen |
| `.harness/skills/reviewer.md` | Reviewer-Skill: HIGH-Liste, Kategorien-Regeln, Negativbefund-Pflicht, Output-Schema (Modul 10) — nächste Rolle nach Schritt 8 des Minimal Agent Workflow, nicht Teil der Implementer-Eingabe |
| `.harness/baseline/<tag>/regelwerk/` (vendored; `README.md` = Index) | adoptiertes Betriebsregelwerk in Agenten-Kurzform — **präsente nachschlagbare Vertiefung**, pro Entscheidung abschnittsweise (siehe [`AGENTS.md`](../AGENTS.md) §1); derivativ, Stand/Tag siehe [`conventions.md`](conventions.md) §Baseline |
| `.harness/baseline/<tag>/templates/` (vendored, parallel) | Referenz-Form der Skelette, auf die das Regelwerk mit `../templates/…` als „Ziel-Form" verweist (netzlos, weil parallel zu `regelwerk/`); Vorlagen zum Kopieren-und-Ausfüllen |

## Sensors (Feedback-Gates)

<!--
WICHTIG: Nur Befehle aufzählen, die im Makefile *existieren*.
Halluzinierte Gates sind die häufigste Form von Harness-Lüge (Modul 13).

Drei Spalten — kein Lauf-Status:
- Target:  der Make-Befehl.
- Vertrag: was prüft das Gate (was wäre verletzt, wenn es rot wird).
- Bindung: strukturelle Referenzen — Carveout-ID (`CO-<NNN>`),
  Slice-ID, Schwelle, Image-Hash, ADR-ID. NICHT der Lauf-Status,
  sondern was das Gate *strukturell trägt*.

Lauf-Wahrheit pro Commit liegt in CI (Badge/Dashboard), nicht hier
(`harness/README.md` ist Rang 9 in der Source Precedence).
Strukturell rote Gates (dauerhaft rot) bekommen einen Carveout in
`docs/plan/carveouts/CO-<NNN>-…` mit Auflösungs-Trigger und Folge-Slice
(Modul 7); die Bindung-Spalte verweist auf die `CO-<NNN>`-ID, die
Begründung lebt im Carveout, nicht hier.

Bei d-check-Einsatz (≥ v0.73.0) deckt Modul `reviews` (Ziel `doc-reviews`,
Review-Report-Deckung für `done/`-Slices mit Review-DoD-Haken) die
Code→Review-Kante ab, Modul `planning` (Ziel `doc-planning`,
Planning-Lifecycle-Konsistenz) die Verify→Closure-Kante — beide aus dem
Lebenszyklus-Diagramm in Modul 1. Nur eintragen, wenn das Ziel im
Makefile existiert (siehe oben).

NICHT-GATES: Ein Target, das der Agent braucht, aber das nichts über den
Zustand des Repos urteilt — es *bewegt* (Slice-Move), *misst* (Latenz) oder
*sagt*, was ein schreibender Lauf täte — steht in der zweiten Tabelle und
trägt `kein Gate` IN DER ZEILE SELBST, in der Spalte, die hier die Bindung
führt — nicht in Prosa daneben. Weglassen ist nur für Targets richtig, die
niemand braucht (Modul 13 §Vorhanden ≠ behauptet).

WÄCHST DIE SEKTION: Die Tabelle bleibt klein, die Prosa darunter nicht.
Braucht ein Gate mehr als EINEN SATZ — Deckungsgrenze, Ausgabe-Bedeutung,
Exit-Codes, Abbruch-Bedingungen —, wandert das nach. Ob der Überhang schon
unter der Tabelle steht oder in die Zelle gedrängt wurde, ist dieselbe Sache:
eine Zelle, die zum Absatz geworden ist, ist der Fund, nicht die Ausnahme.
`harness/sensors/<target>.md`, und die **Target-Zelle wird zum Link darauf**
— wie die `MR`-Zelle im Adaptions-Block. Der Link ist kein Komfort: Er ist die
einzige Fassung dieser Zuordnung, die der Link-Sensor prüft. Eine bloße
Namenskonvention (`make X` -> `sensors/X.md`) bleibt still grün, wenn die
Datei verschwindet und die Zeile stehen bleibt. Seine Grenze: Er prüft EINE
Richtung — ob das Ziel existiert; eine Datei ohne Index-Zeile und eine Zeile
auf die falsche Datei bleiben still grün, und geprüft wird nur, wo ein
Link-Sensor über `harness/` läuft. Kein `sensors/done/`:
ein retiriertes Gate verschwindet, `git` hält seine Geschichte. Was das
Werkzeug selbst deckt (welcher Test welche Hälfte trägt), gehört NICHT
dorthin, sondern in seine ADR/Spec-Zeile/seinen Skriptkopf.
-->

| Target | Vertrag | Bindung |
|---|---|---|
| `make lint` | <was prüft es> | — |
| `make test` | <…> | — |
| `make arch-check` | <…> | ADR-<NNNN> |
| `make coverage-gate` | <…>, bootstrap-aware | Schwelle X %, M<n> → Y % |
| `make coverage-gate-critical` | <…> | bootstrap via `CO-<NNN>` bis <Slice/Welle> |
| `make gates` | alle inneren Gates | — |
| `make ci` | gates + extras | — |
| `make fullbuild` | volle Closure | Image-Hash `sha256:…` (Modul 14) |
| [`make <gate-mit-grenze>`](sensors/<target>.md) | <…>; Grenze und Ausgänge in der verlinkten Datei | ADR-<NNNN> |

**Werkzeuge — genannt, weil der Lauf sie braucht, aber kein Gate:**

| Target | Tut was | Bindung |
|---|---|---|
| `make <mover>` | bewegt <…>, prüft nichts | kein Gate |
| `make <messung>` | misst <…> gegen <Schwelle> | kein Gate, ADR-<NNNN> |
| [`make <vorschau>`](sensors/<vorschau>.md) | sagt, was <schreibender Lauf> täte; Ausgänge und Sperren in der verlinkten Datei | kein Gate |

**Aktueller Lauf-Status:** CI-Badge bzw. lokal `make help` / `make gates`.
**Rote Gates:** Begründung im verlinkten `CO-<NNN>` (siehe Bindung-Spalte), Modul 7.
**Nicht behauptet** (geplant): `<make-target-1>`, `<make-target-2>` (Welle <n>).

<!-- Domänenspezifische Gates ergänzen, je nach Repo-Klasse: -->

## Traceability rules

- PRs/Commits **müssen** mindestens eine `<LH-*>` oder `ADR-*`-ID nennen.
- Neue oder geänderte Anforderungen brauchen einen Beleg: Test, Gate, Demo oder ADR.
- Neue ADRs müssen im ADR-Index ergänzt werden.
- Änderungen an Planning-Dokumenten müssen die Lifecycle-Regeln beachten (open → next → in-progress → done; reine `git mv`-Commits siehe AGENTS.md §3.3).

## Safety and scope boundaries

<!--
Repo-spezifisch formulieren. Beispiele:

Für ein Referenz-Repo:
- Dies ist kein produktiver Service.
- Externer Cloud-Zugriff darf nicht für lokale Demo-Abnahme vorausgesetzt werden.
- Determinismus und Replayability sind Kernverträge.

Für ein Safety/Control-Repo:
- Markt-/Optimierungs-Output muss durch Statemachine, Constraint-Limiter, Ramp-Limiter fließen.
- Software-Stop ersetzt keine Hardware-Sicherheitsfunktionen.
- Produktion-Profile müssen fail-closed sein.

Für ein Policy/Compliance-Repo:
- Dieses Werkzeug ist keine Rechts-/Steuer-/Fachberatung.
- KI-Funktionen liefern Vorschläge, keine verbindlichen Entscheidungen.
-->

- <…>
- <…>

## Minimal agent workflow

1. Diese Datei lesen.
2. Relevante kanonische Quelle lesen.
3. Betroffene IDs identifizieren.
4. Kleinste Änderung planen.
5. Engsten nützlichen Sensor laufen lassen.
6. Repo-weiten Gate-Lauf vor Handoff (`make gates`).
7. Doku/Indizes aktualisieren, falls ein öffentlicher Vertrag berührt.
8. Ausgeführte Sensors und verbleibende Risiken berichten.

Dieser Workflow deckt ausschließlich die Implementer-Rolle ab. Schritt 8
ist der Rollenwechsel, kein Abschluss: Bericht → Handoff an Reviewer
(`.harness/skills/reviewer.md`, siehe §Guides) → Verifier. Kein
Self-Review — anderer Kontext findet andere Findings, derselbe Kontext
dieselben blinden Flecken (Baseline-Regelwerk `modul-08-agentenrollen.md`).

## Leseordnung

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-harness-dateien.md`
§harness/README.md als Einstiegspunkt — die Menschen-Hälfte des Einstiegs:
drei bis fünf **geordnete** Zeiger, was ein neuer Mensch zuerst liest und was
bei Bedarf; eine Leseordnung, die alles nennt, ist keine.

1. [<zuerst — z. B. `AGENTS.md` §Hard Rules>](<pfad>)
2. [<dann — z. B. `spec/lastenheft.md`>](<pfad>)
3. [<bei Bedarf — z. B. `harness/conventions.md`>](<pfad>)
