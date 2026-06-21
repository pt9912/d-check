# Verifikation R2 — slice-038 (`--print-mk`) — unabhängige Bestätigung

## Kopf-Metadaten

- **Review-Art:** Unabhängige Verifikation (Runde R2) — adversariale
  Reproduktion + Suche nach Übersehenem. **Kein Fix-Commit** zwischen R1 und
  R2: diese Runde bestätigt oder widerlegt R1
  (`docs/reviews/2026-06-21-slice-038-print-mk.md`, 0 HIGH/0 MEDIUM/1 LOW/1
  INFO) und sucht Neues. Frischer Kontext.
- **Datum:** 2026-06-21
- **Gegenstand:** Commit `1864372` (`feat(print-mk): slice-038 --print-mk
  include-bares d-check.mk (DC-FA-CLI-010)`), Stand `HEAD = 1864372` auf
  `main`. Betroffen: `internal/adapter/driving/cli/print_mk.go` (NEU),
  `cli.go` (`earlyGenerators`, `--print-mk`-Flag, Dispatch),
  `cli_acceptance_test.go` (`TestCLI038_PrintMK(+_UnbekanntesFlag)`),
  `Dockerfile` (build-Stage `ARG VERSION` + `-ldflags -X`), `spec/lastenheft.md`
  §DC-FA-CLI-010 (0.22.0), `spec/spezifikation.md` §DC-FA-CLI-010.a,
  `docs/user/operations.md`, `CHANGELOG.md`,
  `docs/plan/planning/done/slice-038-print-mk.md`.
- **Reviewer:** Claude (Agent), Skill `.harness/skills/reviewer.md` v1.0.0.
- **Modell:** `claude-opus-4-8[1m]`.
- **Eingangs-Kontext:** R1-Report; Slice-Plan; `DC-FA-CLI-010` + §DC-FA-CLI-010.a;
  read-only-Vertrag `DC-QA-03`, Determinismus `DC-QA-02`; Distribution
  `DC-FA-DIST-001`/ADR-0002/ADR-0011 (Digest-Pin-Politik); `MR-006`
  (Spec-Straten verweisen nie abwärts auf ADRs); Hard Rules `AGENTS.md` §3.
  **Docker verfügbar** (über `make`-Targets, AGENTS.md §3.1): `--print-mk` real
  am gebauten `d-check:latest` getestet — Dev-Build (`0.0.0-dev`), Test-Build
  `--build-arg VERSION=9.9.9-test` und ein Adversarial-Build mit Sonderzeichen
  (`1.2.3-rc.1+meta_x%y`); Mutationsprobe in einem Wegwerf-Worktree
  (`git worktree`, Repo unverändert).

## Findings

### HIGH

Keine.

### MEDIUM

Keine.

### LOW

#### LOW-1 (bestätigt R1) — `--print-mk` kombiniert still mit anderen Ausgabe-Modi statt als Nutzungsfehler abgewiesen zu werden

- **Kategorie:** LOW
- **Quelle:** Maintainability (Reviewer-Skill LOW-Anker „latente Wartungsfalle");
  Konsistenz zum `comboError`-Vertrag (`cli.go:126`)
- **Pfad:** `internal/adapter/driving/cli/cli.go:143` (`earlyGenerators`),
  `cli.go:333` (Kurzschluss vor den `comboError`-losen Modi)
- **Befund:** `earlyGenerators` schaltet allein an `o.printConfig`/`o.printMK`
  kurz und ignoriert jeden anderen gesetzten Modus, **sofern dieser kein
  `comboError` ist**. Real beobachtet (Entrypoint überschrieben, alle Combos):
  `--print-mk --json`, `--print-mk --yaml`, `--print-mk --doctor`,
  `--print-mk --trace`, `--print-mk --repair`, `--print-mk --repair-broad`
  geben **alle** still das `.mk`-Fragment bei Exit 0 aus; `--print-config
  --print-mk` und `--print-mk --print-config` geben **beide** nur das
  Config-Gerüst aus (Switch-Reihenfolge `printConfig` zuerst — `--print-mk`
  verschwindet lautlos). Kein Spec-Verstoß: §DC-FA-CLI-010 (Negative-AK)
  fordert nur „unbekanntes Flag → Exit 2". Daher LOW.
- **Verifizierbar:** ja — `docker run --rm --network none --entrypoint
  /d-check d-check:latest --print-mk --json` liefert das `.mk`-Fragment bei
  Exit 0 (kein Exit 2).
- **R2-Einordnung (LOW-1 ist *identisch* zu `--print-config`, nicht
  schlimmer):** Die in der R2-Vorgabe genannten Verdachtsfälle wurden
  reproduziert und gegen `--print-config` gespiegelt — **kein Fall ist
  schlimmer als bei `--print-config`:**
  - `--print-mk --repair` / `--print-mk --trace` / `--print-mk --doctor` →
    Exit 0, Fragment; `--print-config --repair`/`--trace` ebenso (Baseline
    identisch).
  - `--print-mk` + Pfad-Argument (`--print-mk /nonexistent`, `--print-mk
    /repo`) → Exit 0, Fragment — der Kurzschluss liegt **vor** `openRoot`,
    der Pfad wird nie gelesen (genau die gewollte read-only-/Mount-Freiheit).
  - **Wichtige Nuance (R2-Neufund zur Mechanik, nicht zum Verdikt):**
    `comboError` läuft in `parseOptions` **vor** `earlyGenerators`. Combos,
    die `comboError` *fängt*, gewinnen daher gegen `--print-mk`:
    `--print-mk --json --yaml` → Exit 2 („--json und --yaml…"), `--print-mk
    --repair --json` → Exit 2, `--print-mk /a /b` → Exit 2 („höchstens ein
    Pfad-Argument"). Das stille Kombinieren tritt also nur für Modus-Paare auf,
    die selbst kein `comboError` sind — und für genau diese verhält sich
    `--print-config` Byte-für-Byte gleich. Die in der R1-LOW-1 vermutete
    Asymmetrie zu `--print-config` existiert **nicht**; LOW-1 bleibt eine
    konsistente, geräuscharme Konventionslücke ohne Eskalationsanlass.

### INFO

#### INFO-1 (bestätigt R1) — Mount-Freiheit des Generators ist gewollte, undokumentierte zugesicherte Eigenschaft

- **Kategorie:** INFO
- **Quelle:** `DC-QA-03` (read-only); dokumentationswürdige Annahme
- **Pfad:** `internal/adapter/driving/cli/cli.go:333` (Kurzschluss vor
  `openRoot`), `Dockerfile:136` (`ENTRYPOINT ["/d-check","/repo"]`)
- **Befund:** `docker run --rm --network none d-check:latest --print-mk`
  liefert das Fragment bei Exit 0 **ohne `-v`-Mount und ohne Netz** (real
  beobachtet, leeres stderr) — der ENTRYPOINT prependet `/repo`, doch
  `earlyGenerators` schaltet vor `openRoot` kurz, sodass `/repo` weder
  verlangt noch gelesen wird. Korrekte, gewollte read-only-Eigenschaft (genau
  das, was a-checks Bootstrap braucht), nirgends als zugesicherte Eigenschaft
  notiert. Kein Fehler.
- **Verifizierbar:** ja — wie oben, Exit 0 ohne Mount/Netz.

#### INFO-2 (R2-neu) — `--print-mk`-Test fixiert nur das `:v`-Präfix, nicht den Versions-*Wert*

- **Kategorie:** INFO
- **Quelle:** dokumentationswürdige Test-Scope-Annahme (Reviewer-Skill
  INFO-Anker)
- **Pfad:** `internal/adapter/driving/cli/cli_acceptance_test.go:1354`
  (`"DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v"`)
- **Befund:** Der Happy-Test assertet das Präfix `…/d-check:v`, **nicht** den
  konkreten eingebetteten Versionswert. Eine Mutationsprobe belegt die
  Aussagekraft *innerhalb dieses Scopes*: das Streichen des `v`-Präfix
  (`:v%s` → `:%s`) und das Ersetzen des TAB durch Spaces im Recipe machen
  `TestCLI038_PrintMK` jeweils rot (real reproduziert im Wegwerf-Worktree).
  Eine fehlerhafte `-ldflags`-Verdrahtung, die einen *falschen Wert* statt
  `0.0.0-dev` einbettete, würde der Test **nicht** fangen — diese Achse ist
  bewusst über den Build-Arg-Pfad verifiziert (Negativbefund unten), nicht
  über den Unit-Test. Kein Fehler — bewusster Test-Scope (der Wert ist auf dem
  Dev-Build ohnehin `0.0.0-dev`).
- **Verifizierbar:** ja — Mutation `:v%s`→`:%s` bzw. `\t`→Spaces ⇒
  `make test` rot (FAIL `TestCLI038_PrintMK`).

## Negativbefunde (geprüft, ohne Befund)

- **R1-Kern (a) Parse + Expansion / TAB (`DC-FA-CLI-010` Boundary-AK):**
  `docker run --rm d-check:latest --print-mk > d-check.mk` (551 B) +
  `make -n -f d-check.mk doc-check` expandiert zu `docker run --rm --network
  none -v "/tmp/t037:/repo:ro" ghcr.io/pt9912/d-check:v0.0.0-dev` (Exit 0);
  `cat -A` zeigt die Recipe-Zeile mit führendem `^I` (TAB). **Bestätigt.**
  <!-- d-check:ignore (illustrativer /tmp-Expansionspfad, kein Repo-Artefakt) -->
- **R1-Kern (b) `DCHECK_IMAGE`-Override CLI + Umgebung:** `make -n -f
  d-check.mk doc-check DCHECK_IMAGE=…@sha256:deadbeef` UND
  `DCHECK_IMAGE=…@sha256:cafe make -n -f d-check.mk doc-check` ersetzen beide
  den Ref im Recipe (`?=` gibt CLI- und Env-Override nach). **Bestätigt.**
  <!-- d-check:ignore (illustrative Beispiel-Digests sha256:deadbeef/cafe) -->
- **R1-Kern (c) Version-Einbettung, `-X`-Pfad, kein Doppel-`v`:** Build
  `make build VERSION=9.9.9-test` → `--print-mk` zeigt `DCHECK_IMAGE ?=
  ghcr.io/pt9912/d-check:v9.9.9-test` (einfaches `v`, `-X`-Pfad
  `…/cli.version` trifft die Var). **Bestätigt.**
- **R1-Kern (d) Determinismus + read-only:** zwei `--print-mk`-Läufe
  byte-identisch (`cmp` still, Exit 0); `--print-mk` mit `--network none`
  **ohne** Mount → Exit 0, Fragment, leeres stderr (kein Repo-Zugriff).
  **Bestätigt.**
- **R1-Kern (e) MR-006 + Anker:** Sektion-Scan über §DC-FA-CLI-010
  (`lastenheft.md:359`) und §DC-FA-CLI-010.a (`spezifikation.md:316`) — **kein**
  ADR-Verweis im Körper (unabhängig per `awk`-Scan reproduziert). Anker-Slug
  `#dc-fa-cli-010--makefile-fragment-ausgeben` ist in `operations.md:30`,
  `CHANGELOG.md:17`, `spezifikation.md:316`, Slice-Plan:12 einheitlich und löst
  auf (`make doc-check` grün). **Bestätigt.**
- **Template-Edge-Cases (Sprintf-`%`, Trailing-Newline, CRLF):** Ausgabe
  enthält **kein** literales `%`; `mkTemplate` enthält genau **ein** `%s` und
  kein weiteres `%` (`print_mk.go:30`). Letztes Byte ist `0a` (Trailing-Newline
  vorhanden — include-/Konkatenations-sicher); **kein** CR (kein CRLF). Kein
  Finding.
- **Sonderzeichen im Tag / `-ldflags`-Quoting (adversarial):** Build mit
  `VERSION='1.2.3-rc.1+meta_x%y'` → `--print-mk` emittiert
  `…/d-check:v1.2.3-rc.1+meta_x%y` **verbatim** — das einfach-gequotete
  `-X '…cli.version=${VERSION}'` (`Dockerfile:112`) hält `+`/`%`/`_`, und das
  `%` korrumpiert das `Sprintf` **nicht** (es ist das *Argument*, nicht der
  Format-String). Ein `%`-Tag ergäbe einen nicht-ziehbaren Docker-Ref — das ist
  jedoch ein Malformed-Tag-In-Problem; reale Release-Tags sind SemVer-rein (auf
  `X.Y.Z` gestrippt). Kein Code-Defekt. Kein Finding.
- **Henne-Ei / Dev-Default dokumentiert:** Code-Kommentar (`print_mk.go:9–14`),
  Spec (`spezifikation.md:327` „Default `0.0.0-dev` für lokale/Gate-Builds")
  und Out-of-Scope (`lastenheft.md`) benennen konsistent: das Binary kennt
  seine Version, nicht seinen Digest; der Dev-Build emittiert das
  nicht-ziehbare `:v0.0.0-dev`, der Konsument pinnt per `DCHECK_IMAGE`-Override
  auf `@sha256:`. Annahme dokumentiert. Kein Finding.
- **arch-check / Hexagon (neue Datei `print_mk.go`):** Datei liegt im
  cli-Adapter (driving) und importiert ausschließlich `fmt` (`print_mk.go:7`) —
  **keine neue Import-Kante**, keine Kern-/Port-Verletzung. `make arch-check`
  grün („Import-Regeln R1–R5 (ADR-0005) + R6 (ADR-0012) eingehalten"). Kein
  Finding.
- **Refactor-Grenze (`suggestConfig` NICHT in `earlyGenerators`):** Nur die
  zwei *repo-freien* Generatoren (`--print-config`/`--print-mk`) sind in
  `earlyGenerators` (`cli.go:143`); `--suggest-config` (liest das Repo) bleibt
  korrekt **nach** `openRoot` (`cli.go:342`). Die Grenze des Refactors stimmt —
  kein versehentlicher Repo-Lese-Modus in den Kurzschluss gezogen. Kein
  Finding.
- **`--print-config`-Regression:** `--print-config` liefert weiterhin das
  Config-Gerüst (Exit 0); `earlyGenerators` lagert den vormaligen Inline-Block
  1:1 aus. Kein Finding.
- **Negative-AK (`DC-FA-CLI-010`):** `TestCLI038_PrintMK_UnbekanntesFlag`
  (`--print-mk --bogus` → Exit 2) deckt die Negative-AK; `flag.Parse` läuft vor
  dem Kurzschluss, daher schlägt das unbekannte Flag zuerst zu (real:
  `--print-mk --bogus` → Exit 2). Kein Finding.
- **Test-Mutations-Aussagekraft:** Wegwerf-Worktree-Probe — Mutation TAB→Spaces
  im Recipe ⇒ `make test` rot (FAIL `TestCLI038_PrintMK`: die Newline-TAB-
  Recipe-Assertion schlägt fehl); Mutation `:v%s`→`:%s` ⇒ rot (die
  `DCHECK_IMAGE`-`:v`-Präfix-Assertion schlägt fehl). Beide Kern-Eigenschaften
  (TAB, `v`-Ref) werden vom Test gekillt. Kein Finding (Wertschärfe siehe
  INFO-2).
- **Regressions-Gates (HEAD = 1864372):** `make test` grün (cli-Paket
  `ok 0.054s`, inkl. `TestCLI038*`); `make arch-check` grün; `make doc-check`
  grün („98 Datei(en) geprüft, 0 Befund(e)" — 98 statt R1-97, weil der
  R1-Report seither im Repo liegt). Kein Finding.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 0 | 1 | 2 |

## Verdikt

**Closure-fähig — R1 bestätigt, keine neuen blockierenden Punkte.** Alle fünf
R1-Kernpunkte sind unabhängig reproduziert: (a) das Fragment parst + expandiert
`$(CURDIR)`/`$(DCHECK_IMAGE)`, Recipe TAB-eingerückt; (b) `DCHECK_IMAGE` ist per
CLI **und** Umgebung überschreibbar; (c) `--build-arg VERSION=9.9.9-test` ergibt
`:v9.9.9-test` (einfaches `v`, `-X`-Pfad trifft die Var); (d) byte-identisch und
mount-/netzfrei; (e) MR-006 eingehalten, Anker lösen auf. LOW-1 ist nach
adversarialer Spiegelung **identisch** zu `--print-config` (kein Fall —
`--repair`/`--trace`/Pfad-Arg — ist schlimmer; die vermutete Asymmetrie
existiert nicht), bestätigt als geräuscharme Konventionslücke ohne
Eskalationsanlass. **R2-neu, nicht-blockierend:** INFO-2 (Test fixiert nur das
`:v`-Präfix, nicht den Versionswert — bewusster Scope, per Mutationsprobe
abgesichert) und die Mechanik-Nuance, dass `comboError` vor `earlyGenerators`
greift (combos, die comboError fängt, ergeben Exit 2 — auch mit `--print-mk`).
Sonderzeichen im Tag (`%`/`+`) brechen weder das `-ldflags`-Quoting noch das
`Sprintf`. Drei Regressions-Gates (`test`/`arch-check`/`doc-check`) grün auf
HEAD. **Kein HIGH/MEDIUM. slice-038 ist closure-fähig**; LOW-1 als bewusster
Won't-Fix oder Folgepunkt zu vermerken.

<!-- d-check:ignore (illustrative /tmp- und Beispiel-Pfade/Digests: /tmp/t037, sha256:deadbeef/cafe, VERSION 9.9.9-test, 1.2.3-rc.1+meta_x%y — keine realen Repo-Artefakte) -->
