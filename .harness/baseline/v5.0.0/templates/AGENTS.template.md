# AGENTS.md — Briefing für AI-Coding-Agenten

> **Template-Hinweis.** Diese Datei ist eine Vorlage für die
> Repo-Root-`AGENTS.md`. Kopiere nach `AGENTS.md` deines Repos,
> ersetze `<Platzhalter>` und lösche diesen Block. AGENTS.md *trägt
> Hard Rules und Pointer auf kanonische Quellen*, sie *dupliziert deren
> Inhalt nicht* — sonst entsteht Drift (siehe
> [Baseline-Regelwerk §AGENTS.md-Regeln](../regelwerk/modul-09-implementierung.md#agentsmd-regeln-modul-9)).
> **Pointer-Artefakt:** verweist auf andere kanonische Quellen — zuletzt
> füllen bzw. re-syncen, sobald die Ziele stehen; veraltete
> `(folgt)`/Klartext-Verweise fängt kein Linter (Reviewer-Sache).

---

## 1. Was diese Datei ist

Onboarding-Briefing für jede AI-Session, die in diesem Repo Code oder
Dokumentation ändert. Sie verweist auf die kanonischen Quellen und
formuliert die Hard Rules, die der Implementer-Agent immer
einhalten muss.

Regeln dieser Datei: Baseline-Regelwerk `modul-09-implementierung.md`
§Ziel-Form: AGENTS.md — sie trägt Hard Rules und Pointer auf kanonische
Quellen, sie dupliziert deren Inhalt nicht; sonst entsteht Drift.

**Bei Konflikt zwischen dieser Datei und einer kanonischen Quelle gilt
die kanonische Quelle** (Source Precedence — siehe
`harness/README.md`).

Strukturregeln (ID-Schemata, Verzeichniskonvention, Adaptionen ggü.
Baseline, Modus-Deklarationen pro Sub-Area, Zusatzklassen für
Sensors-Bindung) leben in
[`harness/conventions.md`](harness/conventions.md).

Das **Regelwerk der adoptierten Baseline** ist die **präsente,
nachschlagbare Vertiefung** zu diesem Briefing: ein self-navigierbares
**Modul-Bundle** (`README.md` = Index). Beim Bootstrap wird das
self-contained Release-ZIP
(<https://github.com/pt9912/ai-harness-course/releases/download/v5.0.0/lab-regelwerk.zip>)
**committet vendored** unter `.harness/baseline/<tag>/{regelwerk,templates}/`
(Regelwerk *und* Templates parallel, netzlos materialisiert samt `SHA256SUMS`
— Vorgehen siehe
Baseline-Regelwerk `modul-02-harness-bootstrap.md`;
Quelle/Stand in [`harness/conventions.md`](harness/conventions.md) §Baseline).

Die verkörperte Form (dieses Briefing, die Konventionen, deine
ausgefüllten Artefakte) **führt**; das Regelwerk wird **pro Entscheidung
nachgeschlagen, deren
operative Detailtiefe das Briefing nicht trägt** — Trigger-Klassen,
Sub-Area-Qualifikation, Carveout-vs-Reconciliation, Modus-Diagnose. Dabei
**nur den benötigten Abschnitt** laden (README ist der Index), **nicht das
ganze Regelwerk im Kontext halten**. Breiterer Pflicht-Blick bleibt bei:
Bootstrap, Änderung an [`harness/conventions.md`](harness/conventions.md)
(Adaptionen `MR-<NNN>`, Source-Precedence, ID-Schema), Drift-Audit gegen die
Baseline. Derivativ: bei Konflikt gelten die kanonischen Quellen.

Die **Skelett-Vorlagen** der Baseline liegen **vendored** unter
`.harness/baseline/<tag>/templates/` (aus demselben Baseline-Bundle) und
tragen zwei Rollen: als **Referenz-Form**, auf die das Regelwerk mit
`../templates/…` als „Ziel-Form" verweist (netzlos, weil parallel zu
`regelwerk/` vendored), und als Vorlage, die beim Anlegen neuer Artefakte
(ADR, Slice, Welle, Carveout, Review-Report) **kopiert und ausgefüllt** wird
statt frei zu formulieren.

## 2. Kanonische Quellen (Source Precedence)

In dieser Reihenfolge:

1. [`spec/lastenheft.md`](spec/lastenheft.md) — vertraglich abnahmebindend.
2. [`spec/spezifikation.md`](spec/spezifikation.md) — technisch verbindlich, fortschreibbar.
3. [`spec/architecture.md`](spec/architecture.md) — Komponenten- und Sequenzsicht.
4. [`docs/plan/adr/`](docs/plan/adr/) — ADR-Verzeichnis und -Index.
5. [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md) — aktuelle Welle.
6. `docs/user/*` *(falls vorhanden)* — Operations, Quality, Releasing. <!-- d-check:ignore (Verzeichnis optional; entlinkt, da im frischen Repo selten vorhanden) -->
7. [`README.md`](README.md) — Projekt-Überblick.
8. **AGENTS.md (diese Datei).**
9. [`harness/README.md`](harness/README.md) — Harness-Einstieg.

## 3. Harte Regeln

<!--
Eigene Hard Rules ergänzen, basierend auf der Repo-Klasse. Beispiele
zur Inspiration:
-->

### 3.1 Docker-only

<!-- Wenn das Repo Docker-only ist (typisch für Multi-Toolchain-Repos): -->

Kein lokales <venv/SDK/Toolchain-Install>. Alles läuft über `make`
(das Docker nutzt). Host braucht nur Docker und GNU `make`.

**Falsch:** <z.B. `pip install ...`>
**Richtig:** <z.B. `make test`>

**Begründung:** Toolchain-Reproduzierbarkeit + Supply-Chain-Defense.

### 3.2 Suppression-Verbot

<!--
Pro Sprache eine Variante. Beispiele:
- Python: # noqa
- Go: //nolint
- C#: #pragma warning disable, [SuppressMessage]
- Kotlin: @Suppress
- Java: @SuppressWarnings
-->

Inline-Suppression bricht das `<suppression>-gate`. Ausnahmen leben in
<zentraler Konfigurations-Datei> mit Begründung.

### 3.3 git mv + Inhaltsänderung = zwei Commits

Wenn eine Datei verschoben **und** der Inhalt umgeschrieben wird, sind das
zwei Commits — der Move-Commit bleibt rein (Git erkennt R-Rename). Welcher
zuerst kommt, sagt der Vorgang:

1. Regelfall: `git mv source target` → eigener Commit, dann Inhalt umschreiben.
2. Lifecycle-Übergang nach `done/`: erst der Inhalt (DoD-Häkchen,
   Closure-Notiz), dann der reine `git mv` — die Notiz ist die Bedingung für
   `done/`, nicht ihre Folge.

**Begründung:** Sonst fällt die Rename-Detection unter die 50%-
Similarity-Schwelle und `git log --follow` wird unzuverlässig.

### 3.4 Architektur ist sprach- und meilensteinfrei

`spec/architecture.md` referenziert Modul-Pfade, aber **keine** Wellen,
Slices, Commit-Hashes oder Closure-Daten. Die zeitliche Schicht lebt in
`docs/plan/planning/` und den späteren Closure-Notizen. Auch **keine
ADR-Bezüge**: Die Sicht steht im Stabilitäts-Rang über der ADR; welche ADR
eine Aussage verbindlich macht, deklariert die ADR in ihrem `Schärft:`-Feld.

Diese Regel ist *verkörpert*, nicht hier entschieden — sie folgt aus dem
Sicht-Stratum (Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: Architektur-Sicht).

### 3.5 ADRs sind nach `Accepted` immutable

Eine ADR mit Status `Accepted` wird nicht inhaltlich überschrieben.
Korrekturen entstehen als neue ADR mit `Supersedes ADR-NNNN`.

### 3.6 Gates dürfen nicht ohne ADR gelockert werden

Jede Schwellen-Senkung (Coverage, Linter-Strenge, Architekturregel)
ist ein ADR, kein PR-Kommentar.

<!--
Repo-spezifische Hard Rules ergänzen, z.B. für Safety/Control:
- "Optimierer darf nie direkt aufs Gerät schreiben."
- "Protokoll-Adapter dürfen keine Marktentscheidungen enthalten."
- "Produktion-Profile müssen fail-closed sein."
-->

## 4. Quality Gates

Regeln dieser Sektion: Nur Targets aufzählen, die im Makefile **existieren**.
Halluzinierte Gates sind die häufigste Form von Harness-Lüge
(Baseline-Regelwerk `modul-13-quality-gates.md`).

| Target | Zweck |
|---|---|
| `make lint` | <…> |
| `make test` | <…> |
| `make arch-check` | <…> |
| `make coverage-gate` | <…> |
| `make gates` | alle inneren Gates (mandatory vor PR) |
| `make ci` | CI-äquivalent (gates + zusätzliche) |
| `make fullbuild` | volle Closure (vor Welle-Merge) |

Diese Tabelle **listet auf**; definiert wird hier nichts. Die *Bindung* eines
Targets — welche Anforderung oder Entscheidung es durchsetzt — steht in
`harness/README.md` §Sensors; von dort führt der Weg zur `LH-*`-ID, zur ADR
oder zum Carveout.

## 5. Dokumentations-Regeln

- Requirement- und Architektur-IDs müssen in PRs/Commits referenziert
  sein. Vergeben werden IDs beim Spec-/ADR-Schreiben nach dem in
  `harness/conventions.md` deklarierten ID-Schema (Default:
  `<PREFIX>-FA-<NN>` / `<PREFIX>-QA-<NN>` aus dem Lastenheft, ADR-Nummern
  über den ADR-Index) — nie ad hoc im PR.
- Neue ADRs müssen den ADR-Index aktualisieren.
- Roadmap/Status-Geschichte lebt in `docs/plan/planning/`, nicht in `spec/architecture.md`.

## 6. Minimal Agent Workflow

Pro Slice:

1. `harness/README.md` lesen.
2. Relevante kanonische Quelle lesen (Source Precedence beachten).
3. Betroffene Requirement-/ADR-IDs identifizieren.
4. Kleinste sinnvolle Änderung planen.
5. Engsten nützlichen Sensor laufen lassen.
6. Repo-weiten Gate-Lauf vor Handoff (`make gates`).
7. Doku/Indizes aktualisieren, falls ein öffentlicher Vertrag berührt.
8. Ausgeführte Sensors und verbleibende Risiken berichten — keine Erfolgsmeldung ohne Gate-Ausführung.
