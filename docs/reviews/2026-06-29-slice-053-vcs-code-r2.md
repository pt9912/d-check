# Review-Report — slice-053 Modul `vcs` (Code- + Gate-Hälfte), R2

| Feld | Wert |
| --- | --- |
| **Gegenstand** | Code + Gate-Mechanik von slice-053 (opt-in Modul `vcs`, git-Diff-Immutabilität via go-git, Dogfood-Ersatz des `adr-check`-Gates) |
| **Datum** | 2026-06-29 |
| **Reviewer** | unabhängig (Subagent) |
| **Rolle** | Code-Review (kein Verifier, kein Stil-Polizist); die Doku-Hälfte ist separat reviewt (R1) |
| **Diff-Umfang** | uncommitted gegen HEAD `6eab73b`: `internal/hexagon/port/driven/vcs.go`, `internal/adapter/driven/git/{git.go,git_test.go}`, `internal/hexagon/core/rules/{vcs.go,vcs_test.go,run.go}`, `internal/hexagon/core/model/{config.go,finding.go}`, `internal/adapter/driven/configyaml/configyaml.go`, `internal/adapter/driving/cli/{cli.go,cli_vcs_test.go}`, `tools/arch-check.sh`, `Makefile`, `.githooks/pre-commit`, `.github/workflows/ci.yml`, `.d-check.yml`, `go.mod`/`go.sum` |
| **Quellen** | `DC-FA-VCS-001` (lastenheft §, spezifikation §.a), ADR-0024, slice-053, `AGENTS.md` §3, `tools/adr-immutable-check.sh` (Parität), `.harness/skills/reviewer.md` |

## Findings

### F-1 · MEDIUM · ADR-0024 (Fitness Function) / Maintainability · `Makefile:176` (+ `.d-check.yml:16`)

`make adr-check` ist `--enable vcs --disable links --disable anchors`. Die
Modul-Basis ist aber `.d-check.yml` `modules:` =
`[links, anchors, ids, matrix, codepaths, spans, hostpaths, versions]`; der
effektive Satz (EffectiveModules: Basis ∪ enable ∖ disable) ist daher
`{ids, matrix, codepaths, spans, hostpaths, versions, vcs}` — **nicht** „nur
`vcs`". Beim `pre-commit`-Hook (`STAGED=1`) scannen die sechs Nicht-`vcs`-Module
über den `:ro`-Mount den **Arbeitsbaum** (inkl. ungestaged WIP), während `vcs`
nur HEAD↔Index difft; ein nicht-ADR-Befund (z. B. `codepaths`/`versions`/`ids`)
färbt damit den als „ADR-Immutable" deklarierten Hook bzw. die CI rot, **ohne
dass eine Accepted-ADR berührt wurde** — entgegen der ADR-0024-Fitness „grün
sonst". Warum genau `links`+`anchors` (und nicht die übrigen sechs) abgewählt
sind, ist nirgends begründet. **Kein** Silent-Grün: `vcs` feuert weiterhin auf
echte ADR-Drift, der Gate über-feuert lediglich.
**Verifizierbar:** ja — ein bewusster Nicht-ADR-`codepaths`/`versions`-Verstoß im
Arbeitsbaum + `make adr-check` (ohne RANGE) ⇒ Exit 1 mit Nicht-`vcs`-Reason; bzw.
EffectiveModules-Ausgabe.

### F-2 · INFO · DC-FA-VCS-001 (fail-closed) · `internal/adapter/driven/git/git.go:51–67`

Der Kommentar verspricht, ein „vorhandenes, aber unlesbares HEAD" sei „ein echter
Fehler (fail-closed)". `hasHead()` kollabiert jedoch **jeden** `repo.Head()`-Fehler
(auch ein unlesbares HEAD) auf `false` ⇒ `--staged` gibt `nil,nil` (grün) zurück.
Korrekt für den unborn-HEAD-Erstcommit (Skript-Parität), aber die im Kommentar
behauptete Unterscheidung existiert im Code nicht; ein präsentes, aber korruptes
HEAD wäre damit theoretisch still grün statt Exit 2.
**Verifizierbar:** nein — exotischer Korruptions-Edge; PlainOpen scheiterte real
zuvor.

### F-3 · INFO · Maintainability · `internal/hexagon/core/rules/vcs.go:92–107` + `internal/hexagon/core/model/finding.go:40` (SortFindings)

Treten an **einer** Datei zugleich Körper-Drift **und** unzulässiger
Status-Übergang auf, teilen beide Befunde `File`/`Line`/`Rule`/`Target`/`Reason`
(`core-drift-vcs`); SortFindings dedupliziert sie zu **einem** Befund und
verschluckt die zweite `message`. Gate-Semantik bleibt korrekt (≥1 Befund ⇒ rot),
nur die Diagnose-Vollständigkeit leidet. Kein Test deckt den Kombinationsfall ab.
**Verifizierbar:** ja — `make test` mit einem Fall „Accepted→Proposed + Körper-Edit"
zeigte `len==1` trotz zweier Ursachen.

### F-4 · INFO · Maintainability (Coverage) · `internal/hexagon/core/rules/vcs.go:149–152`; `internal/adapter/driving/cli/cli_vcs_test.go`

Der `statusLine == nil`-Zweig von `vcsHeadStatusLineNo` wird nie getroffen
(Dogfood-Config + alle Tests setzen `status-line`). Zudem nutzen die e2e-Tests ein
**beschreibbares** Temp-Repo (`PlainInit`/`Add`/`Commit`); ein hypothetischer
go-git-Schreibversuch auf `:ro` würde erst im `:ro`-Image-Pfad (`make adr-check`/
`image-test`) sichtbar, nicht im Unit/e2e-Test.
**Verifizierbar:** ja — `make coverage-gate` zeigt die ungedeckten Zweige.

### F-5 · INFO · Skript-Parität (sichere Richtung) · `internal/hexagon/core/rules/vcs.go:129–143` (via `excludedRanges`)

Die `exclude-sections`-Abgrenzung endet beim nächsten Heading der Stufe ≤ der von
`## Geschichte` (Stufe 2); das abgelöste Skript skippte dagegen über ein
**Level-1**-Heading (`# `) nach `## Geschichte` hinweg. Divergenz nur, falls ein
`# `-Heading auf Geschichte folgt (nicht das ADR-Layout) — dann ist das Modul
**strenger** (fängt einen Edit, den das Skript weggestrippt hätte). Kein
Silent-Grün, nur eine Parität-Differenz in die sichere Richtung.
**Verifizierbar:** nein — kein realer ADR-Aufbau trifft den Fall.

## Negativbefunde (geprüft, ohne Befund)

- **Silent-Grün / fail-closed:** fehlendes `.git` (cli `resolveVCS`→`gitadapter.Open`),
  nicht auflösbare/leere/`..`-lose Range (`vcsRefs` + `treeAt`→`ResolveRevision`),
  Port-Fehler (CheckVCS reicht `error` durch → run.go → cli Exit 2) sind durchgängig
  fail-closed; der `--staged`-„kein HEAD ⇒ leer" ist die gewollte
  „nichts-zu-schützen"-Parität zum Skript. Abgedeckt von `TestVCS_FailClosed` /
  `TestVCSFailClosed`.
- **DC-QA-03 read-only / netzlos:** der git-Adapter nutzt ausschließlich lesende
  Plumbing-Aufrufe (`PlainOpen`, `ResolveRevision`, `CommitObject`, `Tree`, `Diff`,
  `Storer.Index`, `BlobObject`, `Files`); keine Index-/Working-Tree-Schreibung, kein
  Remote/Netz.
- **DC-QA-02 Determinismus:** `diffTreeIndex` iteriert Maps in unbestimmter
  Reihenfolge, aber run.go stabilisiert den Befundsatz via `SortFindings`.
- **Skript-Parität (7 Klassen):** `TestVCSModified` reproduziert Geschichte-Anhang,
  Superseded-Übergang, Körper-Edit, Status-Rückfall, Proposed-BASE-Grandfathering,
  Edit-nach-Geschichte, Körper-`**Status:**`-Edit (+ Reflow); `TestVCSDeleteAddClass`
  Löschung/Hinzufügung/Klassen-Filter. Kopf-Status-Strip stoppt korrekt an der ersten
  `## `-H2 (`vcsHeadStatusLineNo`); Rename = Delete+Add wird über den Delete-Pfad des
  alten immutablen Pfads gefangen; `vcs.paths`-Glob (`[0-9]*.md`) schließt
  `README.md` aus (path.Match).
- **go-git-Isolation (ADR-0005/0024, arch-check R2):** der erweiterte R2-Zweig flaggt
  jeden go-git-Import außer `internal/adapter/driven/git` (deckt auch den Kern ab);
  Port (`port/driven/vcs.go`) und `rules/vcs.go` importieren kein go-git; das
  cli-Wiring (Composition Root) ist regelkonform.
- **§3.2 Suppression-Verbot:** kein `//nolint` in den neuen/geänderten Go-Dateien.
- **§3.1 Docker/make-only:** `make tidy` führt `go mod tidy` in Docker aus; keine
  Host-Toolchain im Pfad.
- **Supply-Chain (§3.6/ADR-0011):** go-git `v5.19.1` + ~20 transitive Indirects in
  `go.sum` gepinnt; `GO_VERSION` 1.26.4 ≥ `go.mod` `go 1.25.0` (konsistent);
  `github.com/pkg/errors` nur transitiv (gomodguard prüft Direkt-Importe) — durch
  ADR-0024 bewusst akzeptiert.

## Kategorie-Summary

| Kategorie | Anzahl |
| --- | --- |
| HIGH | 0 |
| MEDIUM | 1 (F-1) |
| LOW | 0 |
| INFO | 4 (F-2…F-5) |

## Verdikt

**Bedingt mergebar.** Kein HIGH, kein Silent-Grün-Pfad: die fail-closed-Kette,
DC-QA-02/03-Zusagen, die Skript-Parität (7 Klassen) und die go-git-Isolation sind
sauber. **F-1 (MEDIUM) vor Merge klären/begründen** — `make adr-check` läuft nicht
„nur `vcs`", sondern den halben Default-`doc-check`-Satz mit, was der ADR-0024-Fitness
„grün sonst" widerspricht und den als „ADR-Immutable" deklarierten Hook an
fremd-Befunden rot färben kann; entweder den Modulsatz auf `vcs` verengen oder die
bewusste Kopplung in ADR/Doku festschreiben. Die vier INFO-Punkte sind nicht
blockierend (Diagnose-Vollständigkeit, Coverage, Kommentar-Präzision,
Parität-in-sichere-Richtung).
