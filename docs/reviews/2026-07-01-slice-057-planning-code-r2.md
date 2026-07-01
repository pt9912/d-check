# Review R2 — slice-057 Modul `planning` (Go-Code / Gate-Umbau / Tests)

## Kopf

- **Datum:** 2026-07-01
- **Gegenstand:** slice-057 — hermetisches opt-in Modul `planning`
  (mechanisiert das entfernte `tools/planning-consistency.sh`), Gate-Umbau
  `make planning-check`, zwei gocyclo-Refactors, Tests.
- **Rolle:** unabhängiger, adversarialer Code-Reviewer (Fokus: Go-Code +
  Gate + Tests + adversariale Korrektheit; NICHT die Doc-Straten, NICHT DoD).
- **Baseline:** `.harness/skills/reviewer.md` v1.2.0, `AGENTS.md` §3,
  `tools/arch-check.sh` (R1–R6).
- **Diff:** `git diff HEAD` + neue Dateien `internal/hexagon/core/rules/planning.go`,
  `internal/hexagon/core/rules/planning_test.go`,
  `internal/adapter/driving/cli/cli_planning_test.go`.
- **Fachliches Vorbild:** `git show HEAD:tools/planning-consistency.sh` (Skript-Semantik).
- **Verifikation:** gezielte Proben gegen das gebaute `d-check:latest`
  (`37339d4c6214`, 27 min alt) mit `docker run --rm --network none -v <wegwerf>:/repo:ro`,
  Wegwerf-Repos unter dem Scratchpad. Kein neues Image gebaut. `make gates`
  NICHT (Verifier-Rolle).

---

## Findings

### MEDIUM-1 — ungültiges `planning.slice-glob` verschluckt den `path.Match`-Fehler → fail-open Silent-Grün im Gate

- **kategorie:** MEDIUM (Kontext-Eskalation LOW→MEDIUM: Silent-Grün-Pfad in
  einem Gate, `make planning-check` ∈ `gates`)
- **quelle:** `DC-QA-04` / „Stilles-Grün-Pfad in einem Gate" (reviewer.md HIGH-Anker,
  hier durch Misconfig-Bedingung auf MEDIUM gedämpft)
- **pfad:** `internal/hexagon/core/rules/planning.go:97` (`if ok, _ := path.Match(glob, e.Name); ok`)
  + `internal/adapter/driven/configyaml/configyaml.go:334` (`applyPlanning` validiert `slice-glob` nicht)
- **befund:** `planningHasSlices` verschluckt den `path.Match`-`ErrBadPattern`
  (`ok, _`). Ein ungültiges `slice-glob` (z. B. `slice-[.md`) liefert für jeden
  Eintrag `ok=false` ⇒ `hasSlices=false`. Kombiniert mit einer idle-Roadmap bei
  vorhandenem Slice wird `hasActive==hasSlices` (`false==false`) und die reale
  Drift **nicht** gemeldet. `applyPlanning` prüft die Muster-Gültigkeit nicht,
  obwohl die Familie andere Silent-Grün-Muster am Config-Rand fail-closed abfängt
  (`applyCommits`/`compileMatrixToken`/`applyVersions` lehnen leer-matchende oder
  nicht kompilierende Muster ab).
- **verifizierbar:** ja — Wegwerf-Repo mit `slice-glob: "slice-[.md"`, idle-Roadmap,
  einem `slice-001-x.md`; `docker run … --enable planning` liefert
  **`exit=0`, „0 Befund(e)"** statt `planning-drift`/Exit 1 (Probe 7).
- **Mitigation (Kontext für die Kategorie):** die **Default**-Konfiguration
  (`slice-glob` leer ⇒ `slice-*.md`) ist gültig — das ausgelieferte d-check-Gate
  selbst ist grün-ehrlich; der Fehler zündet nur bei bewusster Misconfig des
  optionalen Feldes. Das Fehler-Schlucken spiegelt zudem den Haus-Stil von
  `matchGlob` (`paths.go:115`, `ok, _`). Aber: bei `scan.ignore`/`exempt-paths`
  ist die Schluck-Richtung fail-**safe** (nichts ausgenommen = mehr geprüft),
  bei `slice-glob` fail-**open** (nichts erkannt = Drift verpasst). Das Skript
  kannte das Feld nicht (`slice-*.md` hartkodiert) — die neue Konfigurierbarkeit
  öffnet die Fläche, nicht das Skript.

### LOW-1 — Doppel-`## Aktuelle Welle`-Überschrift: Modul weicht vom Skript ab (grün, wo das Skript rot meldet)

- **kategorie:** LOW (echte DC-QA-04-Erkennungs-Differenz in Silent-Grün-Richtung,
  aber strukturell missgebildete Eingabe)
- **quelle:** `DC-QA-04` (Parität zur Alt-Tool-Familie)
- **pfad:** `internal/hexagon/core/rules/planning.go:77-87` (`planningBlockHasMarker`
  terminiert am ersten `## `-H2) vs. Skript-`active_section` (awk `next` auf die
  Überschrift ⇒ **Wieder-Eintritt** bei einer zweiten `## Aktuelle Welle`)
- **befund:** Bei zwei `## Aktuelle Welle`-Überschriften — Block 1 aktiv, Block 2
  mit Ruhe-Marker — terminiert `planningBlockHasMarker` an der zweiten Überschrift
  (Prefix `## `) und sieht den Marker nicht ⇒ `hasActive=true`. Das awk-Skript
  re-entriert den Block über die zweite Überschrift (`next` statt `exit`), findet
  den Marker ⇒ idle. Bei vorhandenem Slice: Modul = konsistent (grün), Skript =
  Drift (rot).
- **verifizierbar:** ja — Wegwerf-Repo mit zwei `## Aktuelle Welle` (Block1 aktiv,
  Block2 idle) + Slice: Modul liefert **`exit=0`** (Probe 6); das Skript hätte
  Exit 1 geliefert.
- **Einordnung:** die Eingabe verletzt die Roadmap-Invariante (genau eine
  `## Aktuelle Welle`); die „first-block-wins"-Semantik des Moduls ist für sich
  vertretbar. Divergenz dennoch dokumentiert, weil sie in der gefährlichen
  Richtung (grün) liegt.

### LOW-2 — toter Teil-Zweig im Escape-Guard (`p.Roadmap == ".."` von `Contains(…, "..")` subsumiert)

- **kategorie:** LOW (Maintainability / latente Leseverwirrung)
- **quelle:** „Maintainability"
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go:339`
- **befund:** `strings.HasPrefix(p.Roadmap, "/") || p.Roadmap == ".." || strings.Contains(p.Roadmap, "..")`
  — der mittlere Teil-Ausdruck ist vom letzten vollständig subsumiert (wenn
  `Roadmap == ".."`, ist `Contains(…, "..")` bereits true). Kein Korrektheitsfehler
  (fail-closed bleibt fail-closed), aber toter Code. Spiegelt den identischen
  Guard in `applyCodepaths` (Haus-Stil), also konsistent — deshalb nur LOW.
- **verifizierbar:** nein (statisch).

### INFO-1 — `planning.roadmap` wird direkt gelesen, ohne `ResolveConfigPath`/`Kind`-Symlink-Prüfung

- **kategorie:** INFO (dokumentationswürdige Annahme)
- **pfad:** `internal/hexagon/core/rules/planning.go:24` (`fsys.ReadFile(cfg.Roadmap)`)
- **befund:** Anders als `ids`-Targets / `diagrams.defined-in` (die beim Lauf-Start
  über `ResolveConfigPath` + `Kind` Escape/Symlink prüfen) liest `CheckPlanning`
  die Roadmap direkt. Escape/Absolut ist am Config-Rand (`applyPlanning`) bereits
  fail-closed abgefangen; ein **symlink**-Roadmap würde `os.ReadFile` jedoch folgen.
  Rest-Risiko minimal: read-only Mount, `--network none`, Ausgabe ist nur
  drift/kein-drift (kein Inhalt exfiltriert), und ein Angreifer müsste den Symlink
  im Repo platzieren. Gleicht der Direkt-Lese-Semantik von `vcs`/`immutable` auf
  konfigurierte Pfade. Bewusst notiert, nicht als Defekt gewertet.
- **verifizierbar:** nein.

### INFO-2 — `planning-check` löst den vollen globalen Scan aus, obwohl planning nur Roadmap + ein Listing braucht

- **kategorie:** INFO (Effizienz-Notiz, kein Korrektheits-/Sicherheitsbelang)
- **pfad:** `internal/hexagon/core/rules/run.go:64` (`discoverScopes` läuft immer)
- **befund:** planning ist ein Post-Pass, doch der globale Discover-Scan läuft
  trotzdem (Live-Probe: „166 Datei(en) geprüft" bei `--enable planning` +
  FOCUS_DISABLE). Harmlos (Wiederverwendung der Scan-Infra, `FilesChecked` nur
  informativ), aber unnötige I/O für ein Modul, das nur eine Datei + ein
  Verzeichnis-Listing braucht.
- **verifizierbar:** ja (Ausgabe-Zeile), aber ohne Befund-Folge.

---

## Negativbefunde (geprüft, ohne blockierenden Befund je Bereich)

- **Parität zum Skript (Heading/Marker/Glob/`hasActive==hasSlices`):**
  `planningHeadingLine` (TrimRight `" \t\r"` == Heading) deckt den Skript-Guard
  `^## Aktuelle Welle[[:space:]]*$` (keine führenden Zeichen, nur nachlaufender
  Whitespace, CRLF-tolerant) — verifiziert Proben 1–5. `planningBlockHasMarker`
  (Block bis erste `## `-H2, `strings.Contains` je Zeile) spiegelt awk
  `active_section` + `grep -qF` — Marker in **späterer** Sektion zählt nicht
  (Block-Scoping, Test „Block-Scoping"). `planningHasSlices` matcht **Basisnamen**
  (`DirEntry.Name`, fs-Adapter `e.Name()` bzw. MemFS-Segment) gegen `slice-*.md`
  via `path.Match` — inkl. `slice-.md` (`*` matcht leer, Probe 10 ⇒ Drift, wie das
  Shell-Glob). `## Aktuelle Wellen` (Zusatzbuchstabe) ⇒ fail-closed (Probe 4).
  Marker im **Fenced-Code-Block** innerhalb des Blocks ⇒ beide naiv = idle
  (Probe 8, Parität; kein `PreprocessMarkdown`, roher Read wie das Skript).
  **Einzige** Divergenz: Doppel-Überschrift (LOW-1).
- **fail-closed:** fehlende Roadmap-Datei ⇒ `planning-drift` (Probe 5); leere
  Roadmap-Datei (existiert, keine Heading) ⇒ `planning-drift` (Probe 9);
  kaputte Heading ⇒ `planning-drift` (Probe 4); no-slice+aktiv ⇒ Drift (Probe 3).
  Kein Silent-Grün auf dem Default-Pfad. Einzige Silent-Grün-Fläche: MEDIUM-1
  (Misconfig-`slice-glob`). Inert-Fälle (leere `roadmap`-Config, `nil`-fs) korrekt
  (Unit-Test `TestCheckPlanningInert`).
- **Hermetik (DC-QA-03):** `CheckPlanning` berührt ausschließlich den
  Filesystem-Port (`fsys.ReadFile`, `fsys.List`) — kein git, kein Netz, kein
  Schreiben. Alle Proben liefen unter `--network none`. In `runPostPasses` ist
  planning mit **`fsys`** verdrahtet (`run.go:121`), NICHT mit dem VCS-Port —
  funktioniert auch über den `Run`-Einstieg mit `vcs=nil` (E2E-Tests grün ohne
  `--range`).
- **Refactor-Regressionen:** `applyModules` (configyaml) ruft alle Apply-Funktionen
  in **unveränderter** Reihenfolge (IDs→Matrix→External→Codepaths→Hostpaths→
  Diagrams→Versions→**Immutable (kein error, korrekt eingebettet)**→VCS→Commits→
  **Planning (neu, vor Scopes)**→Scopes) mit identischer First-Error-Propagation.
  `runPostPasses` (run.go) behält vcs/commits als fail-closed-`error` (`return nil, err`
  ⇒ Exit 2) und hängt planning als diagnose-only-Befund an; `SortFindings` weiterhin
  am Ende. Keine ausgelassene/umsortierte Prüfung.
- **arch-check (R1/R6):** `planning.go` importiert nur `path`, `strings`,
  `core/model`, `port/driven` — kein `os`/`net`/`syscall`/`io/fs`, kein Adapter,
  kein `go-git`, kein `yaml`. Nutzt den **Port** `driven.Filesystem`, keinen
  Adapter. `path` ist reiner Parser (nicht in der R1-Verbotsliste). d-check:arch-check-
  Stage frisch gebaut (15 min).
- **Config-Guard-Negativtest:** `TestDecode_PlanningFehler` deckt absolut
  (`/abs/…`), `..` und `planning.scope`-Ablehnung (KnownFields) ab. Mutations-Denke:
  Entfernen des Escape-Guards ⇒ beide Bad-Cases würden akzeptiert ⇒ Test fällt;
  Hinzufügen eines `Scope`-Feldes ⇒ scope-Case akzeptiert ⇒ Test fällt. Killt die
  Mutationen. `TestDecode_PlanningHappy` prüft Übernahme + Defaults.
- **Gate-Umbau:** `planning-check: build` + `$(DCHECK_RUN) --enable planning
  $(FOCUS_DISABLE)`. FOCUS_DISABLE = exakt die acht `.d-check.yml`-`modules`
  (`links,anchors,ids,matrix,codepaths,spans,hostpaths,versions`) — spiegelt die
  Liste; `planning` ist **nicht** dabei, also (EffectiveModules: enable vor disable)
  bleibt genau planning aktiv. Live-Probe gegen das echte Repo: **`exit=0`,
  „166 Datei(en), 0 Befund(e)"** — grün-ehrlich (slice-057 in `in-progress/` +
  Roadmap benennt welle-46 = konsistent). `planning-check` bleibt in `gates`
  (Zeile 152) und via `ci: gates image-test`. `tools/planning-consistency.sh`
  vollständig gelöscht (git `D`), `codepaths.ignore-refs`-Tombstone gesetzt.
- **print-mk (`doc-planning`):** hermetisch, ohne `--range`/`--staged`,
  `--enable planning` + `disableAllExcept("planning")` wählt alle 14 übrigen Module
  (inkl. `external`, `commits`, `vcs`) ab, **nicht** planning. Live-`--print-mk`
  bestätigt das Recipe. `TestCLI057_PrintMK_PlanningTarget` deckt es ab.

---

## Kategorie-Summary

| Kategorie | Anzahl |
| --------- | ------ |
| HIGH      | 0      |
| MEDIUM    | 1      |
| LOW       | 2      |
| INFO      | 2      |

---

## Verdikt

**ACCEPT mit Auflage (MEDIUM-1 vor Merge klären).**

Kein HIGH: keine Harness-Lüge auf dem ausgelieferten Default-Pfad, kein
ADR-0005-Import-Verstoß, kein Netz außerhalb `external`, beide Refactors
verhaltensgleich, fail-closed auf allen Default-Klassen verifiziert, Gate-Umbau
korrekt (nur planning läuft, in `gates`, Skript atomar entfernt).

MEDIUM-1 (Silent-Grün bei ungültigem `slice-glob`, verifiziert) ist die einzige
merge-relevante Auflage: die Familie fängt vergleichbare Silent-Grün-Muster am
Config-Rand fail-closed ab — `slice-glob` sollte konsistent validiert werden (der
verschluckte `path.Match`-Fehler ist in dieser Richtung fail-**open**). Da nur die
optionale Misconfig zündet und das Default-Gate ehrlich grün ist, blockiert es die
grundsätzliche Annahme nicht, gehört aber vor Merge adressiert (oder als bewusstes
Won't-Fix mit Begründung dokumentiert). LOW-1 (Doppel-Heading-Divergenz) und die
INFO-Punkte sind nice-to-fix bzw. Awareness.
