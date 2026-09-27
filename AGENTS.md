# AGENTS.md — Briefing für AI-Coding-Agenten

## 1. Was diese Datei ist

Onboarding-Briefing für jede AI-Session, die in diesem Repo Code oder
Dokumentation ändert. Sie verweist auf die kanonischen Quellen und
formuliert die Hard Rules, die der Implementation-Agent immer
einhalten muss.

Diese Datei trägt **Hard Rules und Pointer** auf die kanonischen Quellen und
**dupliziert deren Inhalt nicht** — sonst entsteht Drift (Kanon:
[`modul-09-implementierung.md` §AGENTS.md-Regeln](.harness/baseline/v6.9.0/regelwerk/modul-09-implementierung.md#agentsmd-regeln-modul-9)).

**Bei Konflikt gilt die höherrangige Quelle, und die niedriger rangierte wird
angepasst** (Source Precedence — siehe
[`harness/README.md`](harness/README.md)). Das gilt zwischen **dieser Datei**
und einer kanonischen Quelle ebenso wie zwischen **zwei kanonischen Quellen**.
**Melde den Widerspruch**, statt ihn stillschweigend nach einer Seite
aufzulösen — wer ihn nur befolgt, lässt die falsche Stelle stehen.

Strukturregeln (ID-Schemata, Verzeichniskonvention, Adaptionen ggü.
Baseline, Modus-Deklarationen pro Sub-Area, Zusatzklassen für
Sensors-Bindung) leben in
[`harness/conventions.md`](harness/conventions.md) — **vor jeder Änderung an
Code oder Dokumentation zu lesen**, nicht nur vor Doku-Änderungen.

Das Betriebsregelwerk der adoptierten Baseline ist **committet vendored**:
das nach Modulen und Grundlagen-Abschnitten aufgeteilte Regelwerk liegt
entpackt unter `.harness/baseline/<tag>/regelwerk/` (die dortige `README.md`
ist der Index), samt `.harness/baseline/<tag>/SHA256SUMS`-Integritätsmanifest —
**netzlos auf jedem Checkout präsent**, offline materialisier-/verifizierbar
per `tools/harness/fetch-baseline-cache.sh` (`--verify` offline-Integrität;
`--check-latest` = Currency- + Content-Drift-Audit ggü. Upstream, informativ/kein Gate,
[`MR-022`](harness/conventions.md#mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019); Tag aus §Baseline;
Quelle ist das derivative Release-Bundle
[`lab-regelwerk.zip`](https://github.com/pt9912/ai-harness-course/releases/download/v6.9.0/lab-regelwerk.zip);
Pfadschema/Provenance siehe
[`harness/conventions.md`](harness/conventions.md) §Adoptierte Konventions-Quellen,
[`MR-019`](harness/conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)).
Die **verkörperte Form** (dieses Briefing, die Konventionen, die
ausgefüllten Artefakte) **führt**; das Regelwerk ist die präsente,
nachschlagbare **Vertiefung** und wird **pro Entscheidung** nachgeschlagen,
deren operative Detailtiefe das Briefing nicht trägt — Trigger-Klassen,
Sub-Area-Qualifikation, Carveout-vs-Reconciliation, Modus-Diagnose. Dabei
pro Session **nur den benötigten Abschnitt** lesen, bevor der Workflow (§6)
startet — nicht das gesamte Regelwerk im Kontext halten. **Breiterer
Pflicht-Blick** bleibt bei: Bootstrap, Änderung an
[`harness/conventions.md`](harness/conventions.md) (Adaptionen `MR-<NNN>`,
Source-Precedence, ID-Schema) und dem Drift-Audit gegen die Baseline
([`modul-02-harness-bootstrap.md` §Freshness-Audit](.harness/baseline/v6.9.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2)
— darunter die **Bestands-Stichprobe, die auch bei aktuellem Pin läuft**).
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
5. [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md) — Wellen-Sequenzierung (offene Wellen derivativ).
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
[ADR-0002](docs/plan/adr/0002-distribution-ghcr-image.md)). Der Host braucht `git`, GNU `make`, `bash`,
Docker und die POSIX-Standardwerkzeuge, die die Gate-Skripte rufen
(coreutils, findutils, `grep`, `awk`) — als **Klasse**, nicht als Liste.

**Auch keine Host-Skript-Interpreter** (`python`, `perl`, `ruby`, `node`, `uv`, …)
([`MR-040`](harness/conventions.md#mr-040)).

**Falsch:** `go build ./…`, `go test ./…`, `pip install …`, `python3 - <<EOF`
**Richtig:** `make gates`

**Begründung:** Toolchain-Reproduzierbarkeit + Supply-Chain-Defense — gilt
unabhängig von ihrer Durchsetzung.

**Durchsetzung, zwei unabhängige Schichten:** ein Tool-Call-Wächter
([`.claude/hooks/pretooluse-command-guard.sh`](.claude/hooks/pretooluse-command-guard.sh),
`make guard-probe`, §4) und eine Permission-Sperrliste in
[`.claude/settings.json`](.claude/settings.json). Beides ein
**Stolperdraht, keine Sandbox** — Begründung, Netz-Targets, Guard-Grenzen
und -Tabelle: [`harness/rules/docker-make-only.md`](harness/rules/docker-make-only.md).

### 3.2 Suppression-Verbot

Inline-Suppressions sind verboten: `//nolint`-Direktiven im Code
brechen das künftige Suppression-Gate. Ausnahmen leben zentral in
`.golangci.yml` (exclude-rules) mit Begründung.

**Teilweise durchgesetzt, und die Grenze:** `nolintlint` im Profil
meldet eine Direktive ohne benannten Linter, ohne Begründung oder ohne Wirkung —
sie wird damit sichtbar und zurechenbar. Eine **wohlgeformte** `//nolint`
unterdrückt einen echten Verstoß weiterhin, und `make lint` bleibt grün: der
Linter prüft die **Form** der Direktive, nicht ihre Berechtigung. Verboten
bleibt sie durch diese Regel, nicht durch das Gate. *(Auflösungs-Trigger:
permanent — die Berechtigungsfrage ist ein Urteil.)*

### 3.3 git mv + Inhaltsänderung = zwei Commits

Wenn eine Datei verschoben **und** der Inhalt umgeschrieben wird, sind das
zwei Commits — der Move-Commit bleibt rein (Git erkennt R-Rename). Welcher
zuerst kommt, sagt der Vorgang:

1. Regelfall: `git mv source target` → eigener Commit, dann Inhalt umschreiben.
2. Lifecycle-Übergang nach `done/`: erst der Inhalt (DoD-Häkchen,
   Closure-Notiz), dann der reine `git mv` — die Notiz ist die Bedingung für
   `done/`, nicht ihre Folge. Andere Dateien (Roadmap-Ruhe-Marker,
   Pfad-Verweise) dürfen im selben Move-Commit mitreisen, ohne die
   Rename-Erkennung zu berühren — die bewegte Datei selbst bleibt
   unverändert.

**Begründung:** Sonst fällt die Rename-Detection unter die
50%-Similarity-Schwelle und `git log --follow` wird unzuverlässig.

**Teilweise durchgesetzt:** `make planning-check` hält die **Kopplung**
Lifecycle-Verzeichnis ↔ Roadmap-Ruhe-Marker in beide Richtungen. Die
**Commit-Zerlegung** selbst — Move und Inhaltsänderung in einem Commit — sieht
kein Gate. *(Auflösungs-Trigger: permanent.)*

**Wo die Zerlegung nichts schützt, greift sie nicht:** Ersetzt ein Stub den
Volltext einer Datei im selben Akt vollständig, der sie verschiebt — kein
Zwischenzustand mit unverändertem Inhalt (Wellen-/Slice-/Review-Archiv-Stub-
Moves, die einmalige Register-Formatmigration) —, bleibt es bewusst bei
**einem** deklarierten Commit; git zeigt reine `D`/`A`-Paare, keine Renames.
Einzelfälle im Konventionsspeicher (`harness/conventions/done/`).

**MR-/Wellen-Lifecycle-Move** (`conventions/` → `conventions/done/`, flaches
Wellendokument → `done/`) ist der **Regelfall** (Fall 1): reiner `git mv`
zuerst — die relativen Verweise der bewegten Datei lösen für den Moment
nicht mehr auf, was Kanon ausdrücklich zulässt, solange dieser
Zwischenstand nicht die Spitze eines Push wird —, dann die
Link-Tiefen-Korrektur als eigener Commit. Der lokale `pre-commit`-Hook prüft
den **Arbeitsbaum**, nicht den git-Diff des jeweiligen Commits: ein
Zwei-Commit-Vorgang, dessen Korrektur bereits im Arbeitsbaum vorliegt, bevor
der reine Move committet wird, passiert ihn, obwohl der Move-Commit für sich
genommen inkonsistent bleibt. [`MR-013`](harness/conventions.md#mr-013) ist
vollständig aufgelöst
([`conventions/done/`](harness/conventions/done/MR-013-lifecycle-move-buendelung.md)).

### 3.4 Architektur sprach-/meilensteinfrei; Spec-Straten nie abwärts

[`spec/architecture.md`](spec/architecture.md) benennt Schichten und
Rollen statt Technologie — keine Sprach-/Modul-Pfade. Kein
Spec-Stratum (auch [`spec/spezifikation.md`](spec/spezifikation.md))
referenziert ADRs, Wellen, Slices, Commit-Hashes oder Closure-Daten.
Die sprachkonkrete Übersetzung (Modul-Pfade, Import-Regeln) und die
Begründungen leben in den ADRs, deren `Schärft:`-Feld aufwärts zeigt;
die zeitliche Schicht lebt in `docs/plan/planning/`.

**Die Abwärts-Sperre nennt fünf Kategorien; gedeckt sind vier.** `make doc-check`
(Modul `matrix`) hält **ADRs** über die Link-Prüfung, **Slices**, **Wellen** und
**Commit-Hashes** zusätzlich als Token im Körper
([`DC-FA-MTX-003`](spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)).
Erkennungs-Formen und ihre benannten Grenzen stehen bei der Regel, die sie
trägt — im `matrix`-Block der [`.d-check.yml`](.d-check.yml).

**Ungedeckt bleibt das Closure-Datum:** es ist von einem legitimen Datum nicht
unterscheidbar, und die Spec-Straten führen eigene Historie-Zeilen voller Daten.
*(Auflösungs-Trigger: keiner — die Kategorie bleibt Urteil.)*

**Die Sprachfreiheit der Sicht zerfällt in zwei ungleiche Hälften.** Ob eine
Zeile **Rollen statt Technologie** benennt, ist ein Urteil — *permanent*. Ob sie
einen **Modul-Pfad** trägt, ist ein detektierbarer Zustand: `matrix` führt die
Sicht als eigene Klasse und die Code-Modul-Pfade als Token-Ziel daneben.
Gegenstand sind **Modul**-Pfade, nicht Pfade überhaupt — Dokument- und
Skript-Pfade bleiben erlaubt, denn weder ein Dokument noch ein Gate-Skript ist
ein Modul.

**Das Pfad-Verbot ist eine Verschärfung gegenüber der Baseline** — sie erlaubt
der Sicht Modul-Pfade ausdrücklich. Geführt als
[`MR-033`](harness/conventions.md#mr-033).

### 3.5 ADRs sind nach `Accepted` immutable

Eine ADR mit Status `Accepted` wird nicht inhaltlich überschrieben.
Korrekturen entstehen als neue ADR mit `Supersedes ADR-NNNN` (vierstellig).
Maschinell erzwungen über `make adr-check` (`pre-commit`-Hook + PR-/Push-CI;
erlaubt bleiben `## Geschichte`-Anhänge + der `**Status:**`-Übergang;
[ADR-0016](docs/plan/adr/0016-adr-immutable-gate.md)).

### 3.6 Gates dürfen nicht ohne ADR gelockert werden

Jede Schwellen-Senkung (Coverage, Linter-Strenge, Prüfregel) ist ein
ADR, kein PR-Kommentar.

**Kein Gate prüft das** — die Regel gilt einem **Akt**, nicht einem ruhenden
Zustand: ob eine gesenkte Schwelle eine ADR *hat*, steht in keiner Datei, die
ein Sensor gegen die Senkung halten könnte. *(Auflösungs-Trigger: permanent.)*

### 3.7 Kommentare tragen eine der fünf Klassen

Gilt für Code, Konfiguration und Skripte — und, mit **eigener** Form, für
Zustandsfelder (unten). **Ein Kommentar beschreibt, was da ist**
(Baseline-Merksatz). Er
beantwortet in Code, Konfiguration oder Skript, was der Code nicht
beantworten kann — **Zusage · Kopplung · Abgrenzung · Rang-Zeiger ·
Grenze** ([Baseline §Was ein Kommentar trägt](.harness/baseline/v6.9.0/regelwerk/grundlagen-harness-dateien.md#was-ein-kommentar-trägt--code-konfiguration-skripte)).
Keine Review-Historie und keine Review-Befund-Marker, keine Deliberation
über Verworfenes, keine Herkunfts-Prosa, keine Slice-Nummern und keine
Mess-Labels; Herkunft nur als **ein** auflösbares Feld nach dem
Baseline-Schema (`DC-*` — die Baseline-Form `LH-*` —, `ADR-*`, `MR-*`,
`seit welle-<NN>`). Der Reviewer-Skill trägt den HIGH-Anker dazu.

**Zustandsfelder** (Roadmap-/Register-/Meilenstein-Zellen, das
`**Status:**`-Feld einer ADR) tragen **nicht** die fünf Klassen, sondern
eine eigene Form — Zustand und Beleg als auflösbarer Anker, keine Chronik.
Kein Gate prüft eines von beidem; Bestandsgrenze und Zustandsfeld-Details:
[`harness/rules/kommentare-fuenf-klassen.md`](harness/rules/kommentare-fuenf-klassen.md).
Kanon:
[Baseline §Was ein Kommentar trägt](.harness/baseline/v6.9.0/regelwerk/grundlagen-harness-dateien.md#was-ein-kommentar-trägt--code-konfiguration-skripte).

### 3.8 Ein Modul verspricht nur über das, was es scannt

Jedes Modul gibt seine Zusagen über seine **Scan-Menge**. Daneben liest es
Eingaben, die es nie scannt: Zieldateien außerhalb der Scan-Wurzeln, selbst
benannte Verzeichnisse eines Post-Passes, git-Revisionen. Dort gilt **keine**
dieser Zusagen — und die Folge kann **still** sein: ein verdecktes Heading
macht einen Anker unauflösbar, die Prüfung entfällt kommentarlos. Wer ein Modul
anlegt oder ändert, beantwortet deshalb: **welche Eingaben liest es, die es
nicht scannt — und gilt dort dieselbe Zusage?** Wo sie nicht gilt, gehört die
Grenze in die Anforderung.

**Begründung:** Eine Liste der Achsen wäre selbst unvollständig — deshalb steht
hier die Frage und keine Liste. Kein Gate fängt das; der Reviewer-Skill trägt
den MEDIUM-Anker dazu. *(Hard Rule aus dem
Steering Loop, [`BEO-ALL/module-promise-only-on-scan-axis`](docs/plan/planning/observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md),
seit welle-73; Auflösungs-Trigger: permanent.)*

### 3.9 GitHub-Action-Referenzen sind SHA-gepinnt

Jeder `uses:`-Eintrag in [`.github/workflows/`](.github/workflows) nennt einen
**vollen Commit-SHA** mit Tag-Kommentar dahinter, nie einen beweglichen Tag.
Das gilt für jeden Workflow gleich und für jeden Neuzugang.

**Eine Ausnahme, und sie ist keine Lockerung** ([ADR-0068](docs/plan/adr/0068-lokale-workflow-referenzen-ohne-pin.md)):
eine **lokale** Workflow-Referenz (`uses: ./.github/workflows/x.yml`) kann keinen
SHA tragen und **braucht keinen** — sie löst auf denselben Commit auf wie der
aufrufende Workflow und kann per Konstruktion nicht driften. Geprüft wird
stattdessen Existenz und Rechte-Anforderung des Ziels. Details, Grund-Codes
und die Grenze der Prüfung:
[`harness/rules/github-action-sha-pin.md`](harness/rules/github-action-sha-pin.md).

**Begründung:** Supply-Chain-Härtung — ein Tag lässt sich umhängen, ein SHA
nicht; dieselbe Härtung wie der Docker/make-only-Pfad in §3.1.

**Durchgesetzt:** `make workflow-pins` in `make gates` — über das Modul
`workflows` ([ADR-0072](docs/plan/adr/0072-workflows-modul.md); Dogfooding über
das eigene Image).

## 4. Quality Gates

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-harness-dateien.md`
§harness/README.md als Einstiegspunkt. **Der Gate-Index steht einmal**, in
[`harness/README.md`](harness/README.md) §Sensors — dort steht auch die
*Bindung* jedes Targets, und von dort führt der Weg zur `DC-*`-ID, zur ADR
oder zum Carveout. **Diese Datei führt die Liste nicht.**

**Kein Target nennen, das im Makefile nicht existiert — auch nicht in Prosa.**
Halluzinierte Gates sind die häufigste Form von Harness-Lüge
(Baseline-Regelwerk `modul-13-quality-gates.md`). **Gedeckt ist davon die
Tabellen-Hälfte**, und `make gate-consistency` hält sie in **beide**
Richtungen: ein im Index behauptetes `make X` ohne Makefile-Regel meldet
`gate-phantom`, eine Makefile-Regel ohne Index-Eintrag `gate-undocumented`
([`DC-FA-TGT-001`](spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)).

**Die Prosa-Hälfte trägt kein Mechanismus** — das Modul `targets` liest als
Doku-Target nur `` `make X` `` in Zeilen, deren erstes Zeichen `|` ist, und
zwar nur in [`harness/README.md`](harness/README.md); **diese Datei ist weder Scan-Ziel
noch Autorität**. Ein erfundenes Target im Fließtext — hier oder anderswo —
erzeugt null Befunde. Das ist §3.8 auf diesen Sensor angewandt: er verspricht
nur über seine Scan-Menge. *(Auflösungs-Trigger: permanent — die Prosa-Hälfte
zu decken hieße, jede Backtick-Nennung im Repo als Deklaration zu lesen.)*

## 5. Dokumentations-Regeln

Index-Tabelle nach dem Muster des Adaptions-Blocks in
[`harness/conventions.md`](harness/conventions.md) — vier kurze Regeln
stehen vollständig in der Tabelle, die übrigen als Volltext-Pointer
([ADR-0096](docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md)).

| # | Regel | Datei |
|---|---|---|
| 1 | Commits/PRs nennen mindestens eine `DC-*`/`ADR-*`/`MR-*`/`slice-*`-Kennung (`make trace-check`); IDs werden nur beim Spec-/ADR-Schreiben vergeben, nie ad hoc. | [Volltext](harness/rules/dokumentations-regeln/01-commit-id-pflicht.md) |
| 2 | Dependabot-Commits tragen die Kennung im Commit-Präfix ([ADR-0067](docs/plan/adr/0067-dependabot-als-hebender-kanal.md)), nicht als Gate-Ausnahme. | [Volltext](harness/rules/dokumentations-regeln/02-dependabot-kennung-praefix.md) |
| 3 | Neue/geänderte `DC-*`-Anforderungen entstehen nur in `spec/lastenheft.md`, nie per ADR. | [Volltext](harness/rules/dokumentations-regeln/03-dc-anforderungen-nur-lastenheft.md) |
| 4 | Neue ADRs müssen den [ADR-Index](docs/plan/adr/README.md) aktualisieren. | — |
| 5 | Neue ADRs tragen `## Re-Evaluierungs-Trigger` (oder „permanent"); vor Einführung `Accepted`-ADRs bleiben grandfathered. | [Volltext](harness/rules/dokumentations-regeln/05-adr-re-evaluierungs-trigger.md) |
| 6 | Roadmap/Status-Geschichte lebt in `docs/plan/planning/`, nicht in der Architektur-Spec. | — |
| 7 | Slice-Lifecycle (`open → next → in-progress → done`) ist reine Datei-Bewegung (`git mv`, siehe §3.3). | — |
| 8 | Neue Slice-Köpfe tragen `**Verantwortlich:**`, gesetzt spätestens bei Beanspruchung; Bestand kein Retrofit. | [Volltext](harness/rules/dokumentations-regeln/08-slice-verantwortlich-feld.md) |
| 9 | Slice-Kopf-Feld `**Berührte Spec-Stellen:**` nennt die Kennung des Zielelements, sonst den Abschnitt. | [Volltext](harness/rules/dokumentations-regeln/09-slice-beruehrte-spec-stellen.md) |
| 10 | Was eine Welle einlöst, gehört in ihren Closure-Trigger, nicht in die Slice-DoD; ein `structure`-Sensor hält offene DoD-Haken in `done/`. | [Volltext](harness/rules/dokumentations-regeln/10-welle-einloesung-nicht-in-slice-dod.md) |
| 11 | Slice-Pläne tragen **kein** `**Status:**`-Feld — der Lifecycle-Zustand **ist** die Verzeichnis-Position; neue Slices führen den `**Lifecycle:**`-Hinweis. Alt-Slices in `done/` behalten ihr historisches Feld. | — |
| 12 | Jeder Slice-Plan trägt vor der Modus-Begründung drei Vorprüfungen: Sub-Area · Beobachtungs-Register · Nachtlauf-Stand. | [Volltext](harness/rules/dokumentations-regeln/12-drei-vorpruefungen-slice-plan.md) |
| 13 | Wer eine Grenze aufschreibt, prüft sie gegen den Gegenstand (Code/Config), nicht gegen ihre eigene Beschreibung. | [Volltext](harness/rules/dokumentations-regeln/13-grenzen-gegen-gegenstand-pruefen.md) |
| 14 | Vor einer Messung steht die ausgeschriebene Form ihres Gegenstands, nicht nur ihre Zahl. | [Volltext](harness/rules/dokumentations-regeln/14-messung-form-vor-zaehlung.md) |
| 15 | Eine Commit-Botschaft/Closure-Notiz behauptet nicht mehr, als die gemessene Arbeit trägt. | [Volltext](harness/rules/dokumentations-regeln/15-commit-botschaft-nicht-mehr-behaupten.md) |
| 16 | Eine zitierte Quelle trägt nur, was in ihrem Geltungsbereich steht — das Feld lesen, nicht den Titel. | [Volltext](harness/rules/dokumentations-regeln/16-zitierte-quelle-nur-geltungsbereich.md) |
| 17 | `CHANGELOG.md`/READMEs/Handbuch-Kopf werden in der Release-Prep gepflegt, nicht im Feature-Commit. | [Volltext](harness/rules/dokumentations-regeln/17-changelog-release-prep.md) |

## 6. Minimal Agent Workflow

Pro Slice:

1. [`harness/README.md`](harness/README.md) lesen.
2. Relevante kanonische Quelle lesen (Source Precedence beachten).
3. Betroffene Requirement-/ADR-IDs identifizieren — und **vor der
   Implementierung benennen**: Slice-ID, betroffene `DC-*`-IDs, ADR-IDs,
   betroffene Module, auszuführende Gates
   ([`MR-031`](harness/conventions.md#mr-031)).
4. Kleinste sinnvolle Änderung planen — **samt Abgrenzung.** Was der Slice
   ausdrücklich nicht tut, steht in seinem Plan, und der Lauf bleibt daran
   gebunden: Er darf die Abgrenzung nicht ausweiten, weder still noch
   begründet. Wer im Lauf etwas mitnimmt, das der Plan ausschließt, hat den
   **Plan geändert** — und das gehört vor den Code, nicht in den Bericht
   danach. Der Abschnitt heißt in der Baseline-Form `v6.9.0` §1 *Ziel und
   Abgrenzung*; der eingefrorene Bestand führt ihn als §3 (Kanon:
   Baseline-Regelwerk `modul-09-implementierung.md` §Minimal Agent Workflow).
5. Engsten nützlichen Sensor laufen lassen.
6. Repo-weiten Gate-Lauf vor Handoff (`make gates`).
7. Doku/Indizes aktualisieren, falls ein öffentlicher Vertrag berührt.
8. Ausgeführte Sensors und verbleibende Risiken berichten — keine Erfolgsmeldung ohne
   Gate-Ausführung **und ohne ihre echte Ausgabe**: ein behaupteter Exit-Code ist keiner.

**Schritt 8 ist der Rollenwechsel, kein Abschluss.** Dieser Workflow deckt
ausschließlich die **Implementer**-Rolle ab; nach dem Bericht folgt der Handoff
an den Reviewer ([`.harness/skills/reviewer.md`](.harness/skills/reviewer.md),
siehe [`harness/README.md`](harness/README.md) §Guides) und danach an den
Verifier. **Kein Self-Review** — anderer Kontext findet andere Findings,
derselbe Kontext dieselben blinden Flecken (Baseline-Regelwerk
`modul-08-agentenrollen.md`). Rollen-Trennung ist Kontext-Trennung, nicht
Personen-Trennung: dieselbe Person darf mehrere Rollen füllen, aber nicht im
selben Kontextfenster.
