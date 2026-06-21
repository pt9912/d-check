# Review — slice-038 Implementierung (`--print-mk` → include-bares `d-check.mk`)

## Kopf-Metadaten

- **Review-Art:** Code-Review (Commit-Diff gegen Plan/Anforderungen/Hard Rules
  — **kein Verifier**: DoD-Abhaken ist nicht Gegenstand). Gates wurden hier
  zur Faktenuntermauerung selbst gefahren (Befunde unten mit Kommando belegt),
  nicht als „grün angenommen".
- **Datum:** 2026-06-21
- **Gegenstand:** Commit `1864372` (`feat(print-mk): slice-038 --print-mk
  include-bares d-check.mk (DC-FA-CLI-010)`). Inhalt: neuer read-only-Generator
  `--print-mk` (`internal/adapter/driving/cli/print_mk.go` — NEU), Verdrahtung
  in `internal/adapter/driving/cli/cli.go` (`earlyGenerators`, `printMK`-Flag,
  Dispatch), Akzeptanztests `TestCLI038_PrintMK(+_UnbekanntesFlag)`,
  `Dockerfile` build-Stage `ARG VERSION` + `-ldflags -X`, Spec-CR
  (`spec/lastenheft.md` 0.22.0 §DC-FA-CLI-010, `spec/spezifikation.md`
  §DC-FA-CLI-010.a), `docs/user/operations.md`, `CHANGELOG.md`, Slice-Plan
  nach `in-progress`.
- **Reviewer:** Claude (Agent), Skill `.harness/skills/reviewer.md` v1.0.0.
- **Modell:** `claude-opus-4-8[1m]`.
- **Eingangs-Kontext:** Slice-Plan
  `docs/plan/planning/done/slice-038-print-mk.md`; Anforderung
  `DC-FA-CLI-010` (`spec/lastenheft.md`) + Spezifikation §DC-FA-CLI-010.a;
  read-only-Vertrag `DC-QA-03`, Determinismus `DC-QA-02`; Distribution
  `DC-FA-DIST-001`/ADR-0002/ADR-0014 (Tag-/Digest-Politik); `MR-006`
  (Spec-Straten verweisen nie abwärts auf ADRs); Hard Rules `AGENTS.md` §3
  (insb. §3.4 Spec-nie-abwärts, §3.6 Gate-Lockerung nur per ADR). Kein neues
  ADR — vom Auftraggeber entschieden (Version-Tag-Default, keine neue
  Konsum-Pin-Politik). **Docker verfügbar:** `--print-mk` real getestet am
  gebauten `d-check:latest` (Dev-Version) und an einem Test-Build mit
  `--build-arg VERSION=9.9.9-test`.

## Findings

### HIGH

Keine.

### MEDIUM

Keine.

### LOW

#### LOW-1 — `--print-mk` kombiniert still mit anderen Ausgabe-Modi statt als Nutzungsfehler abgewiesen zu werden

- **Kategorie:** LOW
- **Quelle:** Maintainability (Reviewer-Skill LOW-Anker „latente Wartungsfalle");
  Konsistenz zum etablierten `comboError`-Vertrag (`cli.go:126`)
- **Pfad:** `internal/adapter/driving/cli/cli.go:143` (`earlyGenerators`),
  `cli.go:333` (Kurzschluss vor `comboError`-losen Modi)
- **Befund:** `earlyGenerators` schaltet allein an `o.printConfig`/`o.printMK`
  kurz und ignoriert dabei jeden anderen gesetzten Modus.
  `docker run d-check:latest --print-mk --json` gibt **stillschweigend** das
  `.mk`-Fragment aus und ignoriert `--json` (Exit 0, beobachtet); ebenso
  `--print-mk --doctor`. `--print-config --print-mk` gibt nur `--print-config`
  aus (Switch-Reihenfolge entscheidet, `--print-mk` verschwindet lautlos).
  Die übrigen einander ausschließenden Ausgabe-Modi (`--json`+`--yaml`,
  `--repair`+`--json`, `--trace`+`--doctor`) werden dagegen in `comboError`
  als Nutzungsfehler (Exit 2) abgewiesen. Ein Nutzer, der `--print-mk --json`
  tippt (in der Erwartung maschinenlesbarer Ausgabe oder eines Fehlers),
  bekommt weder das eine noch das andere, sondern ein still „gewinnendes"
  `--print-mk`. **Kein Spec-Verstoß:** §DC-FA-CLI-010 (Negative-AK) fordert nur
  „unbekanntes Flag → Exit 2"; Kombinationen sind nicht spezifiziert. Daher
  LOW (Konsistenzlücke / Überraschung), nicht MEDIUM.
- **Verifizierbar:** ja — `docker run --rm d-check:latest --print-mk --json`
  liefert das `.mk`-Fragment bei Exit 0 (kein Exit 2); `--print-config
  --print-mk` liefert nur das Config-Gerüst.

### INFO

#### INFO-1 — `--print-mk` ist mit der Default-Entrypoint-Form `/repo`-tolerant, ohne dass je ein Repo gemountet wird (gewollte Folge des Kurzschlusses, undokumentiert)

- **Kategorie:** INFO
- **Quelle:** `DC-QA-03` (read-only); dokumentationswürdige Annahme
  (Reviewer-Skill INFO-Anker)
- **Pfad:** `internal/adapter/driving/cli/cli.go:333` (Kurzschluss vor
  `openRoot`), `Dockerfile:136` (`ENTRYPOINT ["/d-check","/repo"]`)
- **Befund:** Der ENTRYPOINT stellt jedem Aufruf `/repo` voran, also wird
  `docker run d-check:latest --print-mk` zu `/d-check /repo --print-mk`.
  `earlyGenerators` schaltet **vor** `openRoot` kurz, daher wird `/repo`
  weder gestattet noch gelesen — `--print-mk` funktioniert **ohne `-v`-Mount**
  (beobachtet: `docker run --rm d-check:latest --print-mk` liefert das Fragment
  bei Exit 0, ohne Mount). Das ist die korrekte, gewollte read-only-Eigenschaft
  und genau das, was a-check für seinen Bootstrap braucht; sie ist nirgends als
  zugesicherte Eigenschaft notiert (im Gegensatz zu `--print-config`, das den
  ENTRYPOINT-`/repo`-Effekt teilt). Kein Fehler — bewusst festgehaltene
  Annahme.
- **Verifizierbar:** ja — `docker run --rm d-check:latest --print-mk` (ohne
  `-v`/`--network`) gibt das Fragment bei Exit 0 aus.

## Negativbefunde (geprüft, ohne Befund)

- **Makefile-Validität / Parse + Expansion (`DC-FA-CLI-010` Boundary-AK):**
  `docker run --rm d-check:latest --print-mk > d-check.mk` und
  `make -n -f d-check.mk doc-check` parst und expandiert sauber zu
  `docker run --rm --network none -v "/tmp/t037:/repo:ro"
  ghcr.io/pt9912/d-check:v0.0.0-dev` (Exit 0). `$(CURDIR)` und
  `$(DCHECK_IMAGE)` lösen korrekt auf — kein Finding.
- **`DCHECK_IMAGE`-Override (CLI + Umgebung):** `make -n -f d-check.mk
  doc-check DCHECK_IMAGE=…@sha256:deadbeef` UND
  `DCHECK_IMAGE=…@sha256:cafe make -n -f d-check.mk doc-check` ersetzen beide
  den Ref im Recipe (`?=` greift, der Override gewinnt). Der dokumentierte
  Digest-Pin-Pfad funktioniert — kein Finding.
- **TAB-Einrückung des Recipes:** `cat -A` zeigt die `doc-check`-Recipe-Zeile
  mit führendem `^I` (TAB), nicht Spaces — make-konform; deshalb parst das
  Fragment überhaupt. Der Test `cli_acceptance_test.go` assertet `"\n\t"`
  (`\ndoc-check:\n\tdocker run`) — deckt das ab. Kein Finding.
- **Sprintf-`%`-Hazard:** `mkTemplate` enthält **genau ein** `%s` und kein
  weiteres `%`; die emittierte Ausgabe enthält **kein** literales `%`. Kein
  versehentliches Doppel-`%`/Format-Verb — `makefileFragment()` ist sauber.
  Kein Finding.
- **Version-Einbettung — `-X`-Pfad korrekt (`DC-FA-CLI-010` Beschreibung):**
  `go.mod` = `module github.com/pt9912/d-check`; der `-ldflags`-Pfad
  `github.com/pt9912/d-check/internal/adapter/driving/cli.version`
  (`Dockerfile:112`) trifft exakt die package-level `var version`
  (`print_mk.go:14`, einzige Definition im Paket, keine Kollision). Build mit
  `--build-arg VERSION=9.9.9-test --target runtime` → `--print-mk` zeigt
  `DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v9.9.9-test` (real verifiziert).
  Dev-Default `0.0.0-dev` greift ohne Build-Arg. Kein Finding.
- **Release-Pipeline reicht VERSION bis in die build-Stage durch:**
  `release.yml:71` ruft `make ci VERSION="$VERSION"` (Version ohne `v`-Präfix,
  bei `:60` aus dem Tag gestrippt). Dry-Run `make -n ci VERSION=7.7.7-citest`
  zeigt die Kette bis `docker build … --build-arg VERSION=7.7.7-citest --target
  runtime` (`ci`→`gates`→`doc-check`→`build`; `build`-Recipe `Makefile:98`
  reicht `--build-arg VERSION=$(VERSION)`). `--build-arg` gilt für alle Stages
  mit `ARG VERSION` (build **und** runtime), beide aus derselben make-Variable
  — keine Drift zwischen eingebetteter Version und OCI-Label möglich. Das
  Template prependet `v` (`:v%s`), die Pipeline strippt `v` — Ergebnis
  `:v0.22.0`, **kein** Doppel-`v`. Kein `release.yml`-Edit nötig. Kein Finding.
- **Henne-Ei sauber benannt (Digest nicht selbst-einbettbar):** Code-Kommentar
  (`print_mk.go:9–13`), Spec (`lastenheft.md:367–371`, `spezifikation.md:330`)
  und Out-of-Scope (`lastenheft.md:386`) benennen konsistent: das Binary kennt
  seine Version, nicht seinen eigenen Digest (der hasht das Binary selbst); der
  Konsument pinnt per `DCHECK_IMAGE`-Override auf `@sha256:`. Der
  Kommentar-Kopf im Fragment zeigt diesen Override-Pfad. Kein Finding.
- **Read-only / Determinismus (`DC-QA-03`/`DC-QA-02`):** zwei aufeinander
  folgende `--print-mk`-Läufe sind byte-identisch (551 B, `cmp` still); der
  Generator schreibt nichts (`fmt.Fprint(stdout, …)`, kein Datei-Handle).
  `--print-mk` schaltet vor `openRoot` kurz — kein Repo-Zugriff (mit
  überschriebenem Entrypoint + nicht-existentem Pfad dennoch Exit 0 + Ausgabe).
  Kein Finding.
- **Kurzschluss-Refactor — keine `--print-config`-Regression:**
  `earlyGenerators` ersetzt den vormaligen Inline-`printConfig`-Block 1:1
  (`configTemplate` unverändert ausgegeben); `--print-config` liefert real das
  Config-Gerüst (Exit 0), `--print-config --print-mk` → Config-Gerüst (Switch-
  Reihenfolge `printConfig` zuerst). Der Refactor lagert nur aus, ohne den
  `--print-config`-Pfad zu verändern. Kein Finding (Kombinations-Überraschung
  separat: LOW-1).
- **Negative-AK — unbekanntes Flag → Exit 2:** `TestCLI038_PrintMK_Unbekanntes`
  `Flag` deckt `--print-mk --bogus` → Exit 2 ab; `flag.Parse` läuft vor dem
  Kurzschluss (`parseOptions` in `Run`), daher schlägt ein unbekanntes Flag
  zuerst zu. Entspricht der Spec-Negative-AK. Kein Finding.
- **MR-006 / AGENTS.md §3.4 (Spec verweist nie abwärts auf ADR/Slice):** Der
  Körper von §DC-FA-CLI-010 (`lastenheft.md`) und §DC-FA-CLI-010.a
  (`spezifikation.md`) enthält **keinen** ADR-Verweis (per `awk`-Sektion-Scan
  bestätigt — insb. kein ADR-0014 im Spec-Körper); die Historie nennt nur
  `slice-038` in der Verweis-Spalte (Planning-Verweis in der Historie ist
  zulässig, der Slice-Bezug ist kein ADR-Abwärtslink im Anforderungs-Körper).
  Kein Finding.
- **Anker-Auflösung / Doku-Verweise:** Das Ziel
  `### DC-FA-CLI-010 — Makefile-Fragment ausgeben` (`lastenheft.md:359`) und
  alle Verweise (`CHANGELOG.md:17`, `operations.md:30`, `spezifikation.md:316`,
  Slice-Plan:12) nutzen einheitlich den Slug
  `#dc-fa-cli-010--makefile-fragment-ausgeben`. `make doc-check` (Dogfooding)
  ist grün — alle Anker/Links lösen auf. Kein Finding.
- **Spec-Schicht-Treue (Lastenheft vs. Spezifikation vs. Code):** Die
  Spezifikation (`spezifikation.md:329`) zitiert das Recipe **exakt**
  (`docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_IMAGE)`) —
  deckungsgleich mit der Emission. Der Lastenheft-Körper (`:373`) zeigt es
  bewusst prosaisch/abgekürzt (`$PWD`, ohne `--rm`), im selben illustrativen
  Stil wie `DC-FA-DIST-001` (`:795`) — Anforderungs- vs. Mess-Stratum, keine
  Drift. Kein Finding.
- **Version-Bump / Historie:** Lastenheft 0.21.0 → 0.22.0 (`:3`), neue
  Historien-Zeile (`:860`) mit Datum 2026-06-21, CR-Begründung, Slice-Verweis;
  korrekt über 0.21.0 eingefügt. CHANGELOG `### Added`-Eintrag oben (vor
  slice-036). Regelkonform. Kein Finding.
- **Kein ADR — vertretbar (`AGENTS.md` §3.6):** Es wird kein Gate gelockert,
  keine bestehende ADR geändert, keine neue Import-Kante eingeführt
  (`make arch-check` grün, R1–R6). Der Version-Tag-Default ist mit der
  ratifizierten Tag-/Digest-Politik (ADR-0014/ADR-0002 §1, Konsument pinnt per
  Digest) konsistent — eine eigene ADR ist nicht erzwungen. Kein Finding.
- **Faktencheck Gates (kein Regress):** `make test` grün (cli-Paket
  `ok 0.046s`, inkl. `TestCLI038*`); `make arch-check` grün
  („Import-Regeln R1–R5 (ADR-0005) + R6 (ADR-0012) eingehalten"); `make
  doc-check` grün („97 Datei(en) geprüft, 0 Befund(e)"). Kein Finding.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 0 | 1 | 1 |

## Verdikt

**Closure-fähig — keine blockierenden Punkte.** Die Umsetzung ist sauber und
real verifiziert: Das Fragment ist ein valides, include-bares Makefile (parst
+ expandiert `$(CURDIR)`/`$(DCHECK_IMAGE)`, TAB-eingerücktes Recipe, kein
Sprintf-`%`-Hazard); `DCHECK_IMAGE` ist per CLI **und** Umgebung
überschreibbar (Digest-Pin-Pfad funktioniert); die Ausgabe ist deterministisch
(byte-identisch) und read-only (kein Mount nötig). Die Version-Einbettung ist
end-to-end belegt — der `-X`-Pfad trifft die Var exakt, ein Build mit
`--build-arg VERSION=9.9.9-test` ergibt `:v9.9.9-test`, und die Release-
Pipeline reicht VERSION lückenlos durch (`make ci VERSION` → `build`-Stage,
kein Doppel-`v`, kein `release.yml`-Edit nötig). `MR-006` ist eingehalten (kein
ADR-Abwärtslink im Spec-Körper), Anker lösen auf, Version-Bump/Historie sind
regelkonform, und alle drei gefahrenen Gates (`test`, `arch-check`,
`doc-check`) sind grün. Offen bleibt **LOW-1**: `--print-mk` kombiniert still
mit `--json`/`--doctor`/`--print-config`, statt analog zu `comboError` als
Nutzungsfehler abgewiesen zu werden — eine Konsistenz-/Überraschungslücke ohne
Spec-Verstoß (die Negative-AK fordert nur Exit 2 bei unbekanntem Flag).
**INFO-1** hält die gewollte Mount-Freiheit des Generators als undokumentierte
Annahme fest. Beides blockiert die Closure nicht; LOW-1 ist als Folgepunkt
oder bewusster Won't-Fix zu vermerken.

<!-- d-check:ignore (illustrative /tmp- und Beispiel-Pfade im Befund: /tmp/t037, sha256:deadbeef/cafe, d-check:t037-vtest — keine realen Repo-Artefakte) -->
