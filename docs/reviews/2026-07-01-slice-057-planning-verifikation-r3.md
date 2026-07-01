# Verifikations-Review R3 — slice-057 Modul `planning` (Fix-Verifikation)

## Kopf

- **Datum:** 2026-07-01
- **Rolle:** unabhängiger Verifikations-Reviewer (R3) — adversariale Prüfung, ob
  die Fixes zu R1/R2 **korrekt + vollständig** sind und **nichts Neues**
  eingeschleppt wurde. **Kein** dritter Voll-Review, **nicht** die DoD.
- **Gegenstand:** die eingearbeiteten Fixes zu
  `docs/reviews/2026-07-01-slice-057-planning-doc-r1.md` (MEDIUM-1, LOW-1) und
  `docs/reviews/2026-07-01-slice-057-planning-code-r2.md` (MEDIUM-1, LOW-1, LOW-2).
- **Baseline:** `.harness/skills/reviewer.md` v1.2.0 (Kategorien/Schema/
  Negativbefund-Pflicht), `AGENTS.md` §3.
- **Verifikations-Mittel:** Statik (`grep`/`Read`); **echte Mutationen** gegen
  `make test` (Test-Image) für R2-MEDIUM-1 und R2-LOW-1; E2E-Proben gegen das
  gebaute `d-check:latest` (`04b40f3deb93`, `docker run --rm --network none -v
  <wegwerf>:/repo:ro`, Wegwerf-Repos im Scratchpad). Kein neues Runtime-Image
  gebaut. Arbeitsbaum nach dem Lauf **byte-identisch** (SHA-256 beider mutierter
  Dateien == Baseline; `git status` deckungsgleich).

---

## Fix-für-Fix-Verifikation

### R1-MEDIUM-1 (CHANGELOG behauptete un-existente Handbuch-Doku) — BESTÄTIGT

- **§6-Modultabelle:** `docs/user/benutzerhandbuch.md:837` trägt jetzt eine
  `planning`-Zeile (opt-in, Grund `planning-drift`, hermetisch, fail-closed bei
  fehlender/mehrdeutiger Überschrift).
- **§5-Block:** `benutzerhandbuch.md:789-806` — Prosa („prüft **hermetisch** …
  `hasActive == hasSlices`") + Konfig-Beispiel + Aufruf `make doc-planning`.
- **§4.13 / `--print-mk`:** `benutzerhandbuch.md:620-622` nennt jetzt „**neun**"
  Targets und listet `doc-planning`.
- **Header + Register:** `benutzerhandbuch.md:3` = „1.18 · v0.36.0" mit Link
  `../../version.md#v0.36.0`; der Anker existiert (`version.md:35`,
  `<a id="v0.36.0"></a>`, plus `version.md:22`). Verlaufszeile 1.18 vorhanden
  (`benutzerhandbuch.md:946`).
- **CHANGELOG-Aussage jetzt WAHR:** `CHANGELOG.md:27` („Benutzerhandbuch §5/§6
  dokumentiert das Modul") ist durch §5-Block + §6-Zeile gedeckt.
- **verifizierbar:** ja (`doc-check` grün, s. u.; `grep -n planning
  docs/user/benutzerhandbuch.md` liefert Modul-, Target- und Historie-Treffer).

### R1-LOW-1 (Boundary-AK ohne `doc-planning`) — BESTÄTIGT

- `spec/lastenheft.md:411` (Boundary-AK von `DC-FA-CLI-010`) nennt jetzt
  `doc-planning` in der **Target-Liste** (`… /doc-commits/doc-planning`) **und**
  in der **Modus-Liste** (`doc-planning → --enable planning + Fokus-Disable,
  **hermetisch ohne Range**`).
- Konsistent mit Beschreibung („sowie **neun** … Targets", `lastenheft.md:378`),
  Happy-AK (`lastenheft.md:410`, `doc-planning … ohne Range`) und Out-of-Scope
  („gelisteten **neun**", `lastenheft.md:414`).
- **verifizierbar:** nein (Prosa-Vollständigkeit; kein Gate).

### R2-MEDIUM-1 (ungültiges `slice-glob` → fail-open Silent-Grün) — BESTÄTIGT (Mutation killt)

- **Fix im Code:** `internal/adapter/driven/configyaml/configyaml.go:347-351` —
  `applyPlanning` prüft ein gesetztes `slice-glob` per `path.Match(...)` und gibt
  bei `ErrBadPattern` Exit 2 zurück. `TestDecode_PlanningFehler`
  (`configyaml_test.go:196`) deckt `slice-glob: 'slice-[.md'` ab.
- **Mutations-Verifikation:** den `slice-glob`-Guard **entfernt** (Block +
  danach ungenutzten `path`-Import) ⇒ `make test`:
  `--- FAIL: TestDecode_PlanningFehler` /
  `FAIL github.com/pt9912/d-check/internal/adapter/driven/configyaml` / Exit 1.
  Der Fix ist **load-bearing** (Test killt die Mutation). Exakt aus Backup
  wiederhergestellt: SHA-256 == Baseline, `git diff`-Hash == Baseline,
  `make test` wieder grün.
- **E2E:** Wegwerf-Repo mit `slice-glob: "slice-[.md"` ⇒ `d-check:latest` meldet
  `error: … planning.slice-glob "slice-[.md" ist kein gültiges Glob: syntax
  error in pattern`, **Exit 2** (kein stilles Grün). Fix auch im ausgelieferten
  Binary wirksam.
- **verifizierbar:** ja (Mutation + E2E, beide reproduziert).

### R2-LOW-2 (toter `p.Roadmap == ".."`-Zweig) — BESTÄTIGT

- `configyaml.go:340` lautet jetzt
  `if strings.HasPrefix(p.Roadmap, "/") || strings.Contains(p.Roadmap, "..")` —
  der subsumierte Teilzweig `p.Roadmap == ".."` ist entfernt; fail-closed bleibt
  fail-closed (E2E-Proben `/abs` und `..` weiter über `TestDecode_PlanningFehler`
  gedeckt).
- **verifizierbar:** nein (statisch).

### R2-LOW-1 (Doppel-Überschrift fail-closed) — FIX im Produkt BESTÄTIGT, **REGRESSIONSTEST WIRKUNGSLOS** (Mutation überlebt) → siehe NEUER BEFUND

- **Fix im Code vorhanden + im Binary wirksam:** `planning.go:36-40` failt bei
  `headingCount > 1` fail-closed (`planning-drift`, „Aktiv-Status mehrdeutig").
  E2E-Probe (c) — zwei `## Aktuelle Welle` (Block 1 aktiv, Block 2 idle) **plus**
  Slice — liefert `d-check:latest` **planning-drift / Exit 1**. Ohne den Guard
  wäre genau dieser Fall (Block 1 aktiv + Slice = „konsistent") **Silent-Grün /
  Exit 0**; der Exit 1 beweist, dass der Guard im ausgelieferten Binary greift.
- **Aber der zugehörige Unit-Test killt die Mutation NICHT** — Details im
  MEDIUM-Befund unten.

---

## Findings

### MEDIUM-1 (NEU) — `TestCheckPlanningDuplicateHeading` überlebt die Mutation: der R2-LOW-1-Fix ist ungetestet (Regressionslücke im Gate-Fail-closed-Pfad)

- **kategorie:** MEDIUM (fehlender/wirkungsloser Negativtest für einen
  Fail-closed-Guard; `make planning-check` ∈ `gates` — Kontext-Eskalation aus
  LOW, weil der ungeschützte Guard genau der Silent-Grün-Schutz ist)
- **quelle:** `DC-QA-04` / reviewer.md „fehlende Negativtests bei neuem
  öffentlichen Vertrag"
- **pfad:** `internal/hexagon/core/rules/planning_test.go:67-73`
  (`TestCheckPlanningDuplicateHeading`) vs. Guard `planning.go:36-40`
- **befund:** Der Test speist eine Doppel-Überschrift-Roadmap **ohne Slice** ein
  (`planCfg()` + MemFS nur mit der Roadmap). Ohne Slice ist `hasSlices=false`,
  und der erste (aktive) Block liefert `hasActive=true`, sodass der Befund
  `planning-drift` bereits über den **normalen** „aktiv-ohne-Slice"-Pfad
  (`hasActive != hasSlices`) entsteht — unabhängig vom `headingCount > 1`-Guard.
  Beide Pfade liefern `len==1` mit Reason `planning-drift`, und der Test prüft
  nur diese beiden Felder. Der Test kann den Guard also nicht von der
  Normal-Drift unterscheiden.
- **Mutations-Beleg:** Der `headingCount > 1`-Block wurde **entfernt** ⇒
  `make test` bleibt **grün** (`ok github.com/pt9912/d-check/internal/hexagon/
  core/rules`, Exit 0) — **kein** Test im gesamten Suite-Lauf fällt. Der
  Regressionstest, der den R2-LOW-1-Fix absichern soll, killt die Mutation
  **nicht**. (Exakt aus Backup wiederhergestellt, SHA-256 == Baseline,
  `make test` wieder grün.)
- **Failure-Szenario:** Entfernt eine spätere Änderung den `headingCount >
  1`-Block, kehrt das in R2-LOW-1 dokumentierte Silent-Grün zurück (Doppel-
  Überschrift + aktiver erster Block + Slice ⇒ fälschlich Exit 0), **ohne** dass
  `make test`/CI es meldet. Der Fix ist heute im Binary korrekt (E2E-Probe (c)
  belegt es), aber nicht durch einen Test verriegelt.
- **verifizierbar:** ja — die obige Mutation ist reproduzierbar; ein Test mit
  Slice + zwei Überschriften (Block 1 aktiv, Block 2 idle) würde die Mutation
  killen.

---

## Negativbefunde (geprüft, ohne blockierenden Befund)

- **`doc-check` (kein neuer Schaden):** `d-check:latest --network none
  -v $PWD:/repo:ro` = **168 Datei(en) geprüft, 0 Befund(e), Exit 0**. Der Anstieg
  gegenüber R1/R2 (166) sind exakt die zwei neuen Review-Report-Dateien
  (`…-doc-r1.md`, `…-code-r2.md`) — sie werden gescannt und sind **span/anker-
  sauber**. Die neuen Spec-/Register-Anker (`#v0.36.0`,
  `#dc-fa-cli-010-…`, `#dc-fa-plan-001-…`) lösen auf.
- **Stale-Zählungen:** `grep -rn "acht|sieben|sechs" spec/ docs/user/` liefert
  nur Wort-Fehltreffer (`beachten`, `gedacht`) — **keine** veraltete
  Target-/Modul-Zählung mehr.
- **E2E-Vollständigkeit (kein Silent-Grün-Pfad):** (a) Drift idle+Slice ⇒
  `planning-drift` / Exit 1; (b) konsistent idle+kein-Slice ⇒ Exit 0; (c)
  Doppel-Überschrift+Slice ⇒ `planning-drift` (fail-closed) / Exit 1; (d)
  ungültiges `slice-glob` ⇒ Config-Fehler / Exit 2. Alle vier wie erwartet.
- **Konsistenz der `neun`-Zählung:** `DC-FA-CLI-010` Beschreibung (378), Happy
  (410), Boundary (411), Out-of-Scope (414) alle auf „neun" + `doc-planning`.
- **Arbeitsbaum unverändert:** SHA-256 beider mutierter Dateien nach
  Wiederherstellung == Baseline; `git status --short` deckungsgleich mit dem
  Eingangszustand (staged `D tools/planning-consistency.sh`, 20 × `M`, untracked
  inkl. R1/R2-Reports + neue Quelldateien). Keine Streu-Artefakte.

---

## Kategorie-Summary

| Kategorie | Anzahl |
| --------- | ------ |
| HIGH      | 0      |
| MEDIUM    | 1      |
| LOW       | 0      |
| INFO      | 0      |

## Verdikt

**NACHBESSERN.**

Vier der fünf Fixes sind **bestätigt** und im ausgelieferten `d-check:latest`
verhaltensgeprüft: R1-MEDIUM-1 (Handbuch §5/§6 nachgezogen, CHANGELOG-Aussage nun
wahr), R1-LOW-1 (Boundary-AK enthält `doc-planning`), R2-MEDIUM-1 (slice-glob-
Validierung — Mutation killt den Test **und** E2E Exit 2) und R2-LOW-2 (toter
Zweig entfernt). Der R2-LOW-1-**Fix selbst** ist im Produktcode korrekt und per
E2E-Probe (c) als Fail-closed/Exit 1 belegt.

Blockierend ist der **neue MEDIUM-1**: der als R2-LOW-1-Absicherung hinzugefügte
`TestCheckPlanningDuplicateHeading` **killt die Mutation nicht** — sein
Slice-loser Eingang lässt die Drift über den Normal-Pfad entstehen, sodass das
Entfernen des `headingCount > 1`-Guards `make test` grün lässt. Damit ist der
Fail-closed-Guard (Silent-Grün-Schutz im Gate) **nicht** durch einen Test
verriegelt; er sollte um einen Slice-behafteten Doppel-Überschrift-Fall ergänzt
werden, der die Mutation killt. Kein neuer HIGH; kein Schaden am Default-Pfad
(`doc-check` grün, E2E sauber).
