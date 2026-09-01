# Templates

Skelett-Vorlagen für die Dokumenttypen des Kurses. **Sprachneutral** —
unabhängig davon, ob dein Repo Go, Python, Kotlin, Java oder C# nutzt.

## Übersicht

Diese Tabelle listet die **20 Dokument-Skelette** (Phase 0 → 1 beim
Bootstrap — das Repo füllt sie). Das Verzeichnis trägt **24 Dateien in drei
Klassen**, nicht 24 gleichartige Vorlagen:

| Klasse | Anzahl | Wo |
|---|---|---|
| **Dokument-Skelette** — Artefakte, die das Repo kopiert und ausfüllt | 20 | Tabelle unten |
| **Skill-Dateien** — Urteilsgrundlage eines Agenten, nicht Artefakt | 2 | [§Skill-Dateien](#skill-dateien) |
| **Tooling-Dateien** (`Makefile`, `.d-check.yml`) | 2 | [§Gate-Baseline](#gate-baseline) |

| Template | Wofür | Regelwerk-Abschnitt |
|---|---|---|
| [`spec/lastenheft.template.md`](spec/lastenheft.template.md) | Vertraglich abnahmebindende Anforderungen (`LH-*`-IDs) | [Modul 3](../regelwerk/modul-03-spec.md) |
| [`spec/spezifikation.template.md`](spec/spezifikation.template.md) | Technisch verbindlich, fortschreibbar — Algorithmen, Defaults, Codes | [Modul 3](../regelwerk/modul-03-spec.md) (Spec-Stratifizierung) |
| [`spec/architecture.template.md`](spec/architecture.template.md) | Komponenten- und Sequenzsicht, sprach- und meilensteinfrei | [Modul 3](../regelwerk/modul-03-spec.md) |
| [`docs/plan/adr/NNNN-titel.template.md`](docs/plan/adr/NNNN-titel.template.md) | Architecture Decision Record im MADR/Nygard-Stil | [Modul 4](../regelwerk/modul-04-adrs.md) |
| [`docs/plan/adr/README.template.md`](docs/plan/adr/README.template.md) | ADR-Index (derivativ; Liste aller ADRs mit Status) | [Modul 4](../regelwerk/modul-04-adrs.md) |
| [`docs/plan/planning/slice.template.md`](docs/plan/planning/slice.template.md) | Slice-Plan mit DoD, Trigger, Closure | [Modul 5](../regelwerk/modul-05-planning-harness.md) |
| [`docs/plan/planning/welle.template.md`](docs/plan/planning/welle.template.md) | Welle als Bündel von Slices | [Modul 5](../regelwerk/modul-05-planning-harness.md) + [Modul 6](../regelwerk/modul-06-roadmap.md) |
| [`docs/plan/planning/welle-results.template.md`](docs/plan/planning/welle-results.template.md) | Welle-Closure-Notiz: Ergebnis, Steering-Loop-Einträge, Zeiger aufs Beobachtungs-Register | [Modul 6](../regelwerk/modul-06-roadmap.md) |
| [`docs/plan/planning/roadmap.template.md`](docs/plan/planning/roadmap.template.md) | Roadmap als Reihenfolge von Wellen, nicht Termine | [Modul 6](../regelwerk/modul-06-roadmap.md) |
| [`docs/plan/planning/observations.template.md`](docs/plan/planning/observations.template.md) | Beobachtungs-Register: der Steering-Loop-Zähler, stehend statt in der Welle-Closure | [Modul 6](../regelwerk/modul-06-roadmap.md#das-beobachtungs-register-modul-6) |
| [`docs/plan/planning/archiv-stub-slice.template.md`](docs/plan/planning/archiv-stub-slice.template.md) | Gekürzter Stub an der Stelle eines archivierten Slice-Volltexts | [Modul 6](../regelwerk/modul-06-roadmap.md) |
| [`docs/plan/planning/archiv-stub-welle.template.md`](docs/plan/planning/archiv-stub-welle.template.md) | Gekürzter Stub an der Stelle eines archivierten Welle-Plans | [Modul 6](../regelwerk/modul-06-roadmap.md) |
| [`docs/plan/planning/reconciliation.template.md`](docs/plan/planning/reconciliation.template.md) | Reconciliation-Register: die offenen Funde der Brownfield-Inventur; nur für Repos im BF-Bootstrap | [Modul 2](../regelwerk/modul-02-harness-bootstrap.md) |
| [`docs/plan/planning/README.template.md`](docs/plan/planning/README.template.md) | Planning-Index: Slice-Lifecycle + Slice-vs-Welle-Konvention | [Modul 5](../regelwerk/modul-05-planning-harness.md) |
| [`docs/plan/carveouts/carveout.template.md`](docs/plan/carveouts/carveout.template.md) | Dokumentierte Ausnahme mit Auflösungs-Trigger | [Modul 7](../regelwerk/modul-07-carveouts.md) |
| [`docs/plan/carveouts/README.template.md`](docs/plan/carveouts/README.template.md) | Carveout-Index (derivativ; aktive/aufgelöste Carveouts) | [Modul 7](../regelwerk/modul-07-carveouts.md) |
| [`docs/reviews/review-report.template.md`](docs/reviews/review-report.template.md) | Review-Report: Kopf-Metadaten, Findings nach Output-Schema, Negativbefunde, Verdikt | [Modul 10](../regelwerk/modul-10-review-harness.md) |
| [`project-readme.template.md`](project-readme.template.md) | Projekt-Root-`README.md`: Überblick, Ist-Stand, Vertrauens-Signale (Rang 6) | [Modul 2](../regelwerk/modul-02-harness-bootstrap.md) |
| [`AGENTS.template.md`](AGENTS.template.md) | Repo-weite Hard Rules und Source Precedence | [Modul 9](../regelwerk/modul-09-implementierung.md) |
| [`harness/README.template.md`](harness/README.template.md) | Repo-Einstiegspunkt mit Guides, Sensors, Safety | [Konventionen](../regelwerk/grundlagen-harness-dateien.md#harnessreadmemd-als-einstiegspunkt) |
| [`harness/conventions.template.md`](harness/conventions.template.md) | Repo-lokale Strukturregeln, Adaptions-Block (`MR-*`), Zusatzklassen-Deklaration, Modus-Deklaration pro Sub-Area | [Konventionen](../regelwerk/grundlagen-harness-dateien.md#harnessconventionsmd-als-konventionsspeicher) |
| [`harness/conventions/MR-NNN-titel.template.md`](harness/conventions/MR-NNN-titel.template.md) | Ein Adaptions-Eintrag (`MR-<NNN>`); Index in `conventions.md` | [Modul 2](../regelwerk/modul-02-harness-bootstrap.md) |

## Download als ZIP

**Stabiler Link (kein Login nötig):** der Workflow `templates-release`
hängt bei jedem Release-Tag *ein* self-contained Baseline-Asset an:

> Baseline-Bundle: <https://github.com/pt9912/ai-harness-course/releases/latest/download/lab-regelwerk.zip>
> — `regelwerk/` (self-navigierbares Modul-Bundle) + `templates/` parallel,
> interne Verweise auf den Tag gepinnt. Nach `.harness/baseline/<tag>/`
> entpacken; die Skelette liegen unter `templates/`.

Zusätzlich lädt der Workflow `templates-zip` diesen Ordner (Artifact
`lab-templates`) bei jeder Änderung als Vorschau-Stand von `main` hoch: auf
GitHub unter **Actions →
templates-zip → neuester Lauf → Artifacts**. Artifacts erfordern einen
GitHub-Login und verfallen nach 90 Tagen; über **Run workflow**
(workflow_dispatch) lässt sich alles jederzeit neu erzeugen.

## Verwendung

1. **Abschnitt lesen** im Baseline-Regelwerk.
2. **Template kopieren** in dein eigenes Repo:
   ```bash
   cp lab/templates/spec/lastenheft.template.md mein-repo/spec/lastenheft.md
   ```
3. **`<Platzhalter>`-Stellen ersetzen.**
4. **Template-Hinweis-Block oben entfernen** (er beginnt mit `> **Template-Hinweis.**`).
5. **HTML-Kommentar-Hilfen entfernen** (`<!-- ... -->`) — **außer**
   `<!-- d-check:ignore … -->`-Marker: die unterdrücken Falsch-Positive
   des Referenz-Gates für bewusst illustrative Pfade und müssen bleiben.
6. **Mit dem entsprechenden Pfad in `lab/example/` vergleichen** —
   so siehst du, wie ein voll ausgefülltes Artefakt aussieht.

## Ein- vs. wiederkehrende Templates

Die Templates haben zwei Lebenszyklen:

- **Singletons** — einmal beim Bootstrap zu `.md` füllen, dann das
  `.template.md` verwerfen: `project-readme`, `spec/lastenheft`,
  `spec/spezifikation`, `spec/architecture`, `AGENTS`, `harness/README`,
  `harness/conventions`, `roadmap`, `observations`, `reconciliation`. *Verwerfen* meint die Kopie im
  Arbeitsbaum — die **vendored** Referenz-Form unter
  `.harness/baseline/<tag>/templates/` bleibt und ist beim Baseline-Update
  die Vergleichsgrundlage für die Form deiner gefüllten Artefakte.
- **Wiederkehrend** — als `.template.md` **co-located** im Repo behalten;
  jede neue Instanz wird daneben kopiert: `adr/NNNN-titel`, `slice`,
  `welle`, `carveout`, `review-report`, `archiv-stub-slice`,
  `archiv-stub-welle` (die beiden entstehen bei jeder Wellen-Closure neu).

Wiederkehrende Templates bleiben also dauerhaft im Repo (z. B.
`docs/plan/adr/NNNN-titel.template.md` neben den echten ADRs). Damit ihre
Platzhalter den Gate nicht rot färben, ignoriert die mitgelieferte
`.d-check.yml` sie per Suffix (`**/*.template.md`). `/tmp` ist nur die
kurzlebige Entpack-Station — der `harness/`-Ordner ist **kein**
Template-Lager.

**Adoptions-Reihenfolge:** Singletons in Abhängigkeitsfolge füllen
(Lastenheft → Architektur → harness → …). **Pointer-Artefakte**
(`AGENTS.md`, `README.md`, `harness/README.md`) verweisen auf die anderen
— sie **zuletzt** füllen bzw. re-syncen, sobald die Ziele stehen. Sonst
veraltet ihr `(folgt)`/Link-Stand: Drift, die der Referenz-Gate nicht
fängt (er prüft Existenz verlinkter Ziele, nicht ob Vorhandenes als
vorhanden beschrieben wird) — Reviewer-Sache.

## Skill-Dateien

Skill-Dateien sind **keine Dokument-Skelette**: Sie werden nicht zu einem
Repo-Artefakt ausgefüllt, sondern tragen die **Urteilsgrundlage einer
Agenten-Rolle** — das repo-spezifische „worauf achtest du hier". Zielort im
Adopter-Repo ist `.harness/skills/`, nicht `docs/` oder `spec/`.

| Skill-Template | Wofür | Regelwerk-Abschnitt |
|---|---|---|
| [`.harness/skills/reviewer.template.md`](.harness/skills/reviewer.template.md) | Reviewer-Skill: HIGH-Liste, Kategorien-Regeln, Negativbefund-Pflicht, Output-Schema | [Modul 10](../regelwerk/modul-10-review-harness.md#ziel-form-reviewer-skill) |
| [`.harness/skills/closure-note-reviewer.template.md`](.harness/skills/closure-note-reviewer.template.md) | Spezialisierter Reviewer für Closure-Notizen | [Modul 11](../regelwerk/modul-11-verification.md) |

**Warum nur zwei bei sechs Rollen?** Eine Rolle braucht genau dann eine
Skill-Datei, wenn ihr Urteil *inferential* ist **und** auf repo-spezifischem
Wissen beruht, das aus keinem Artefakt ableitbar ist — das trifft nur auf den
Reviewer zu. Planner und Architect laufen über Templates, der Implementer über
das Briefing (`AGENTS.md`), Verifier und Validator über die Prüfgrundlage im
Slice. Die zweite Datei ist **keine siebte Rolle**, sondern derselbe Reviewer
mit einem anderen *Urteilstyp*: **Skills wachsen pro Urteilstyp, nicht pro
Rolle.** Zusätzliche Skill-Dateien für Template- oder Briefing-Rollen wären
Attrappen ohne nicht-ableitbaren Inhalt. Kriterium und Zuordnungstabelle:
[Modul 8 §Welche Rolle braucht welche Artefaktklasse](../regelwerk/modul-08-agentenrollen.md#artefaktklasse-pro-rolle).

## Gate-Baseline

Zwei mitgelieferte Dateien geben dir den Doku-Referenz-Gate
out-of-the-box (ins Repo-Root kopieren). Sie sind **Werkzeug-Startgerüste**,
keine Phase-0→1-Dokument-Skelette — das `Makefile` trägt den d-check-
Doku-Gate direkt, `.d-check.yml` die Modul-Auswahl:

| Datei | Rolle |
|---|---|
| [`.d-check.yml`](.d-check.yml) | Modul-Auswahl + Suffix-Ignore; `ids`/`codepaths` wachsen mit den Artefakten |
| [`Makefile`](Makefile) | ruft d-check direkt (`docs-check`-Target, Image per Digest gepinnt); `gates: docs-check`, Code-Gates ergänzt der Adopter |

Danach läuft `make docs-check` sofort (`links`/`anchors`). `ids` und
`codepaths` im `.d-check.yml` einkommentieren, sobald die Ziele bzw.
Verzeichnisse existieren — sonst behauptet der Gate eine Dimension, die er
nicht durchsetzt (Modul 13). Gerüst neu erzeugen: `d-check --print-config`
(leer) oder — für ein Repo nach diesem Kurs-Standard —
`d-check --suggest-config ai-harness-init --id-prefix <PRÄFIX>`, das
`ids`/`matrix`/`codepaths` mit den Kurs-Kennungen (`ADR-…`, `MR-…`,
`slice-…`, `<PRÄFIX>-FA-…`/`-QA-…`) vorbelegt; ohne `--id-prefix` bleibt der
Platzhalter `<PREFIX>` plus `# TODO` stehen.

**Gate-Fragment neu erzeugen.** Das `docs-check`-Target steht direkt im
`Makefile` (gepinnte d-check-Image-Zeile) — genauso dogfooded dieser
Kurs-Repo seinen Doku-Gate. Wer das Fragment lieber frisch generiert (immer
aktueller Pin, kein statisches Duplikat), nutzt `d-check --print-mk` statt
es von Hand zu pflegen.

## Pflichtgliederung vs. freie Form

Die Templates geben **Pflichtgliederung** vor (Abschnitte, IDs,
Verlinkung). Innerhalb der Abschnitte hast du Freiraum — was am
besten zu deinem Projekt passt. Pflicht-Strukturen sind:

- ID-Schema (z.B. `LH-*`) konsistent durchziehen.
- ADRs nach Accepted nicht überschreiben; Schärfung als Folge-ADR mit
  `Supersedes` (Hard Rule, Baseline-Regelwerk `modul-04-adrs.md`).
- Carveouts brauchen immer Trigger + Folge-Slice.
- Slices brauchen DoD mit prüfbaren Kriterien.
- §8 in `slice.template.md` steht in **jedem** Slice-Plan: Die beiden
  *Vorgelagert*-Blöcke (Sub-Area-Wahl prüfen · offene Beobachtungen sichten)
  sind immer auszufüllen. Der **Modus-Begründungsblock** darin ist Pflicht bei
  mindestens einer in BF oder Hybrid berührten Sub-Area, pro berührter
  Sub-Area einer; bei reinem GF genügt der Hinweis "alle berührten Sub-Areas
  GF".
  Voraussetzung-Wissen: Baseline-Regelwerk `modul-05-planning-harness.md`
  §Ziel-Form: Sub-Area-Modus-Begründung.

## Ergänzungen

Wenn du eigene Template-Varianten brauchst (z.B. für ein
Compliance-Repo mit zusätzlichem Disclaimer-Block oder ein
Safety-Repo mit HIL-Test-Plan), lege sie in einem eigenen
`templates/`-Unterordner deines Repos an, *nicht* hier — diese Datei
ist die Referenz-Quelle.
