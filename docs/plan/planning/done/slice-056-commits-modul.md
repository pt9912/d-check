# Slice slice-056: Modul `commits` — Traceability-Kennung in Commit-Messages über eine Commit-Range

**Status:** in-progress (welle-45-commits-modul).

**Welle:** welle-45-commits-modul (Trigger: Auftraggeber-Audit „welche `tools/*.sh`
noch in d-check mechanisieren?" → das **letzte** Gate-Skript der Familie,
`tools/trace-check.sh`; der reine-Go-VCS-Port aus
[`slice-053`](../done/slice-053-vcs-modul.md) macht es mechanisierbar — „das nächste
`adr-check`").

**Bezug:** Führt eine **neue Anforderung** im Lastenheft ein (Modul `commits`,
`commit-untraceable`) plus [ADR-0027](../../adr/0027-commits-traceability-modul.md)
(VCS-Port um Commit-Message-Lesen erweitert; **Teil-Supersede der Skript-Mechanik** von
[ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md) — Policy/Bindepunkt/CI-Topologie
bleiben). Dieselbe VCS-Port-Präzedenz wie
[`slice-053`](../done/slice-053-vcs-modul.md)/[ADR-0024](../../adr/0024-vcs-immutable-gate.md)
(dort Datei-Inhalt, hier Commit-**Messages**) und dieselbe Skript-Ablösungs-Linie wie
[`slice-055`](../done/slice-055-completeness-rueckbau.md)/[ADR-0026](../../adr/0026-completeness-in-product-gate.md).
Verteilung wie der Rest des Werkzeugs (gepinntes Image, kein kopiertes Skript —
[`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Linie).

**Autor:** pt9912. **Datum:** 2026-07-01.

---

## 1. Ziel

`tools/trace-check.sh` ([ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md))
erzwingt die Traceability-Regel („jede Commit-Message nennt eine `DC-`/`ADR-`/`MR-`/`slice-`-ID") —
volle Garantie, aber als **Skript** nur per Kopie in Schwester-Repos nutzbar
(Copy-Drift, [`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Klasse).
Es ist das **letzte** unmechanisierte Gate-Skript der Familie (nach
`adr-immutable-check.sh` → Modul `vcs`, `completeness-check.sh` → in-Produkt-Flag).
**Neu:** dieselbe Prüfung **verteilbar im Image** — ein opt-in Modul `commits` prüft,
dass jede geprüfte Commit-Message eine Kennung nach `commits.id-patterns` trägt, sonst
`commit-untraceable`. Es liest die Commit-**Messages** über **denselben** reine-Go-VCS-Port
wie `vcs` (**ohne** git-Binary → distroless bleibt, **ohne** Netz, read-only), erweitert um
`CommitMessages`. d-check **dogfooded** das Modul für seine eigenen Commits (das Gate
`make trace-check` läuft darauf um).

## 2. Entscheidungen

- **Modulname `commits`**, nicht `trace`: `--trace` und die RTM
  ([`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix))
  belegen `trace` bereits (Doc-/Requirements-Achse). `commits` beschreibt den geprüften
  **Gegenstand** (analog `links`/`anchors`/`pins`/`versions`) und hält `trace` eindeutig
  für die RTM frei. Das **Gate-Target bleibt `make trace-check`** (etablierter Begriff),
  delegiert intern an das Modul. Bereichskürzel `COMMITS`, Grund-Code `commit-untraceable`.
- **Eigenes Modul** (14.), nicht Strategie in `vcs`: anderer Gegenstand (Commit-**Message**
  statt Datei-**Inhalt**), anderer Grund-Code/Semantik. Beide teilen den VCS-Port, nicht die
  Regel.
- **VCS-Port um `CommitMessages(base, head)` erweitert** (Nicht-Merge-Commits,
  `--no-merges`-Parität); die **einzige** git-berührende Stelle bleibt der reine-Go-`git`-Adapter
  (`arch-check` R2). **Keine** neue Dependency — go-git ist seit
  [ADR-0024](../../adr/0024-vcs-immutable-gate.md) da.
- **Zwei Quellen, ein Kern.** Ein reiner Modul-Kern prüft eine Liste von Messages; ihn
  speisen der **Range**-Modus (`--enable commits --range <base>..<head>`, via Port — CI/Push
  + `make trace-check`) und der **Message**-Modus (`--commit-msg <datei|->`, **Kurzschluss-Modus**
  wie `--print-config`/`--trace`, ohne Repo-Scan/Port — der `commit-msg`-Hook pipet die
  Pending-Message über stdin). Beide wenden **dieselbe** git-`strip`-Bereinigung
  (`#`-/scissors-Zeilen) und Prüfung an — keine Divergenz je git-Cleanup-Modus.
- **Lokal/deterministisch/read-only, strikt opt-in, fail-closed** — wie `vcs`: `commits`
  relativiert **weder**
  [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus) **noch**
  [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit);
  Eingabe ist das lokale `.git` / eine Message-Datei, nur der Eingabe-Scope ist erweitert.
  Nie Default-Modul; ohne aktives `commits` byte-identisch. **fail-closed:**
  fehlende/unauflösbare Range, unlesbares `.git`, nicht lesbare Message-Datei → Exit 2.
- **ID-Muster als Config, eine Datei.** `commits.id-patterns` (Liste) + `commits.exempt-pattern`
  (Betreff-Regex, Selbstkonfig `^(Merge |Revert )`) liegen im `commits:`-Block der
  [`.d-check.yml`](../../../../.d-check.yml) **neben** `ids.patterns` — die ID-Definition
  wandert vom hartcodierten Skript-`ID_RE` in dieselbe, gemeinsam reviewbare Datei.
- **Dogfood-Ersatz + Skript entfernt.** `make trace-check`, `commit-msg`-Hook und CI rufen
  das d-check-Image mit `--enable commits` statt `bash tools/trace-check.sh`. Der `commits:`-Block
  bleibt **außerhalb** der Default-`modules`-Liste → `make doc-check` aktiviert ihn nicht
  (default-aus byte-identisch). `tools/trace-check.sh` wird **entfernt** (`git rm`) — die
  immutable [ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md)-Inline-Referenz wird
  ein `codepaths.ignore-refs`-Tombstone ([ADR-0025](../../adr/0025-codepaths-ignore-refs.md),
  **dritter** Anwendungsfall). Der Umweg „Skript zunächst pfad-stabil behalten" (slice-053)
  entfällt — die Lehre ist gezogen.

## 3. Definition of Done

- [x] **Spec/Doc (doc-first):** neue Anforderung
  [`DC-FA-COMMITS-001`](../../../../spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in)
  (Bereichskürzel `COMMITS` in §3, Versions-Bump 0.35.0 + §7-Historie, `commits` in
  [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
  + Glossar) + [ADR-0027](../../adr/0027-commits-traceability-modul.md) (Proposed, + ADR-Index-Zeile;
  die [ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md)-Teil-Supersede-Annotation folgt
  erst bei Closure mit dem Accept) + spezifikation-`.a`-Sektion + Grund-Code
  `commit-untraceable` (§4) + Schema-Keys `commits.*` (§2) + architecture VCS-Port-Rolle.
  Stale ADR-Index-Zeile [ADR-0026](../../adr/0026-completeness-in-product-gate.md) (`Proposed`→`Accepted`) mitkorrigiert.
- [x] **Code:** VCS-Port `CommitMessages` (`port/driven/vcs.go`), go-git-Adapter
  (`adapter/driven/git/git.go`, einzige git-Tür), Kern-Regel (`rules/commits.go`: Bereinigung,
  Betreff-Ausnahme, ID-Prüfung → `commit-untraceable`), Wiring in `run.go` (Range-Post-Pass)
  + CLI-Modus `--commit-msg` (Kurzschluss, stdin/`-`) + `--range`-Mitnutzung,
  `model.CommitsConfig` + `validModules()` + `configyaml`-Kompilierung, Grund-Code in
  `finding.go`. Tests (Skript-Selbsttest-Klassen als Modul-Tests): ID erkannt; fehlende ID
  gefangen; Merge/Revert ausgenommen; Kennung nur auf `#`-Zeile ⇒ Befund (Bereinigung);
  leere gültige Range ⇒ Exit 0; fail-closed (kein `.git`/Range, unlesbare Message); Modul-aus
  byte-identisch (Fake-VCS-Port); E2E `--enable commits --range` + `--commit-msg -`.
- [x] **Gate-Umbau:** `make trace-check`/`.githooks/commit-msg`/`ci.yml` auf `--enable commits`
  (bzw. `--commit-msg -`) umstellen; `commits:`-Block in `.d-check.yml`; `tools/trace-check.sh`
  per `git rm` entfernt + `codepaths.ignore-refs`-Tombstone; `COMMITS_DISABLE` (nur-`commits`,
  aus `ValidModules` abgeleitet, analog `VCS_DISABLE`); Doku (`harness/README.md` §Sensors,
  `AGENTS.md` §4) nachziehen; `gate-consistency` grün.
- [x] `make ci`/`make fullbuild` grün; `make trace-check` (Dogfood, beide Modi); drei
  unabhängige Reviews (R1 doc + R2 code + R3 verifikation); CHANGELOG; Closure (Move nach `done/` + Roadmap-Flip,
  [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise));
  bei Closure: [ADR-0027](../../adr/0027-commits-traceability-modul.md) → Accepted **und** die
  [ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md)-Teil-Supersede-Annotation nachziehen
  (Geschichte-Append + Index-Status-Spalte); Release v0.35.0.

## 4. Risiken / offene Punkte

- **`make trace-check` braucht jetzt das Image** (Hook + CI + lokal bauen/rufen das Image
  statt eines Shell-Skripts) — schwergewichtiger pro Commit; dafür verteilte Wahrheit. Die CI
  baut das Image ohnehin (wie `adr-check` seit [ADR-0024](../../adr/0024-vcs-immutable-gate.md)).
- **commit-msg-Hook über stdin ins Image:** die Pending-Message wird per `docker run -i` an
  `--commit-msg -` gepipet — kein Pfad-Mapping vom Host in den Container nötig (robust für jede
  Message-Datei-Lage). Mitigation: E2E-Test des stdin-Modus.
- **Verlust des Per-Lauf-Selbsttests:** das Skript testete sich bei jedem Aufruf; der
  Selbsttest wandert in `make test` (Modul-Akzeptanztest) — derselbe bewusste Trade wie
  [ADR-0024](../../adr/0024-vcs-immutable-gate.md)/[ADR-0026](../../adr/0026-completeness-in-product-gate.md).
- **Bootstrap-Reihenfolge des Gate-Umbaus:** der neue `make trace-check` darf erst grün
  greifen, wenn das Image das `commits`-Modul trägt; die Umstellung + Skript-Löschung erfolgt
  atomar mit dem Modul.
- **fail-closed vs. Komfort:** ohne `.git`/Range/Message bricht `commits` laut ab (Exit 2) —
  bewusst, damit eine fehlende Eingabe nicht still grün ist.
- **Kein `--repair`.** `commit-untraceable` ist diagnose-only (die Korrektur ist ein neuer
  Commit / ein menschliches `--amend`).

## 5. Trigger

Auftraggeber-Audit 2026-06-29 („welche `tools/*.sh` noch mechanisieren?", Roadmap
welle-45): `tools/trace-check.sh` ist das letzte Familien-Skript; der VCS-Port
([ADR-0024](../../adr/0024-vcs-immutable-gate.md)) macht es mechanisierbar. Modulname
`commits` per Auftraggeber-Entscheid 2026-07-01 (Kollision `trace`/RTM vermieden).

## 6. Sub-Area-Modus-Begründung

GF (Produkt-Code + Spec; „Doc führt, Code folgt"). Die Port-Erweiterung ist eine
GF-Erweiterung des bestehenden VCS-Ports/`git`-Adapters (slice-053); keine BF-Sub-Area.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung.** slice-056 als **14. Regelmodul `commits`**: prüft, dass jede geprüfte
Commit-Message eine Traceability-Kennung (`commits.id-patterns`) trägt, sonst
`commit-untraceable`. Zwei Quellen, ein Kern (`CheckCommits`/`CheckCommitMessage`): der
**Range**-Modus (`--enable commits --range`, via VCS-Port `CommitMessages` — CI +
`make trace-check`) und der **Message**-Modus (`--commit-msg <datei|->`, Kurzschluss,
commit-msg-Hook via stdin). Uniforme git-`strip`-Bereinigung (`#`-/scissors-Zeilen),
Betreff-Ausnahme `commits.exempt-pattern`; strikt opt-in, fail-closed, diagnose-only,
default-aus byte-identisch. Der [ADR-0024](../../adr/0024-vcs-immutable-gate.md)-VCS-Port
um `CommitMessages` erweitert (go-git, **kein** git-Binary → distroless unangetastet;
`arch-check` R2 hält go-git im Adapter). Doc-first:
[`DC-FA-COMMITS-001`](../../../../spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in)
(Lastenheft 0.35.0, Bereich `COMMITS`) + [ADR-0027](../../adr/0027-commits-traceability-modul.md)
+ spezifikation `.a` + Grund-Code gingen dem Code voraus.

**Dogfood + Verteilung.** Das `trace-check`-Gate (Makefile/`commit-msg`-Hook/CI) läuft jetzt
über das Modul `commits` (Image, `--enable commits`/`--commit-msg -`, `FOCUS_DISABLE`) statt
`tools/trace-check.sh`; das Skript ist per `git rm` **entfernt** (drittes Familien-Skript nach
`adr-immutable-check.sh`/`completeness-check.sh`), seine immutable
[ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md)-Inline-Referenz als
`codepaths.ignore-refs`-Tombstone. **Der Dogfood bewies sich an den slice-056-Commits selbst** —
jeder Closure-Commit dieses Slices lief durch das neue `commits`-Modul (commit-msg-Hook).
`--print-mk doc-commits` verteilt die Range-Prüfung an Konsumenten ([`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
7→8 Targets); `--print-config`/`--suggest-config`/Benutzerhandbuch §5/§6 führen `commits`.

**Belege.**
- `make ci` **grün** (doc-check 164/0, lint, test, arch-check, Coverage, semgrep 0/55,
  gate-consistency, planning-check; image-test nativ == Container).
- `make fullbuild` **grün** (36 Anforderungen/**0 Waisen**, lokaler Image-Hash `sha256:9ebb6ead…6139`).
- `make trace-check` (Dogfood) **grün** in beiden Modi; adversariale Probe: kennungslose Message
  → `commit-untraceable` + non-zero (Gate rot), ID-tragend → Exit 0.
- **Drei** unabhängige Reviews (Reports
  [r1](../../../reviews/2026-07-01-slice-056-commits-doc-r1.md)/[r2](../../../reviews/2026-07-01-slice-056-commits-code-r2.md)/[r3](../../../reviews/2026-07-01-slice-056-commits-verifikation-r3.md)):
  R1 (doc) 1 MEDIUM ([`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)-Körper-Zählung) + 2 LOW; R2 (code, adversarial) 1 MEDIUM
  (applyCommits-Negativtests) + 1 LOW + 3 INFO, **keine Paritäts-Divergenz** zum abgelösten
  Skript (14-Fall-Batterie `--message` vs. `--commit-msg -`); R3 (verifikation) **VERIFIED** +
  Mutations-Beleg (Guard-Entfernung ⇒ `TestDecode_CommitsFehler` fällt). Alle Befunde eingearbeitet.
- Tests: `rules/commits_test.go` (Skript-Selbsttest-Klassen via Fake-Port), `git/git_test.go`
  (on-disk go-git: Range-Ausschluss/Merge-Filter/fail-closed), `cli/cli_commits_test.go` (E2E
  Range + `--commit-msg` Datei/stdin + fail-closed), `configyaml_test.go` (Guard-Negativtests,
  mutations-verifiziert), `print_mk`-Test (`doc-commits`).
- Release **v0.35.0** auf GHCR (Digest-Pin folgt via digest-backfill).

**Lerneintrag.** Das **letzte** Familien-Skript (`trace-check.sh`) ist mechanisiert — der
VCS-Port (slice-053) trägt jetzt zwei Modi: Datei-Inhalt für `vcs`, Commit-Messages für
`commits`. Modulname `commits` statt `trace` (Auftraggeber-Entscheid) hält `--trace`/RTM
eindeutig — der geprüfte **Gegenstand** benennt das Modul, nicht die Prüfabsicht. Die
Skript-Löschung erfolgte **atomar** mit dem Modul (der Umweg „pfad-stabil behalten" von
slice-053 entfällt — Lehre gezogen). Config-Surface (print-config/suggest/print-mk/Handbuch)
auf Auftraggeber-Nachtrag vollständig nachgezogen ([`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben) 7→8). Der `commit-msg`-Hook
dogfoodet das Modul an den eigenen Commits — die stärkste Form des „iss dein eigenes Futter".
