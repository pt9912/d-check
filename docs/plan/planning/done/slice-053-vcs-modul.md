# Slice slice-053: Modul `vcs` — Git-Diff-Immutabilität des Core über eine Commit-Range

**Status:** in-progress (welle-42-vcs).

**Welle:** welle-42-vcs (Trigger: Auftraggeber — „können wir `adr-immutable-check`
regelkonform vollständig mechanisieren?" → die von
[ADR-0023](../../adr/0023-immutable-core-pin.md) bewusst vertagte **git-Stufe**:
ein nicht-skriptgebundener, im Image **verteilbarer** git-Diff statt eines
kopierten Shell-Skripts).

**Bezug:** Führt eine **neue VCS-Anforderung** im Lastenheft ein (Modul `vcs`,
git-Diff/`core-drift-vcs`) plus [ADR-0024](../../adr/0024-vcs-immutable-gate.md)
(neuer reine-Go-VCS-Port; **Teil-Supersede der Skript-Mechanik** von
[ADR-0016](../../adr/0016-adr-immutable-gate.md) — Policy/Gate bleiben). Schwester-
und Boden-Hälfte ist der hermetische `immutable`-Pin
([`slice-052`](../done/slice-052-immutable-modul.md),
[ADR-0023](../../adr/0023-immutable-core-pin.md)); der VCS-Port ist von
[ADR-0008](../../adr/0008-reparatur-ableitbarkeit.md) als „künftiges VCS-Modul …
analog `external`" vorgezeichnet. Verteilung wie der Rest des Werkzeugs (gepinntes
Image, kein kopiertes Skript —
[`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Linie).

**Autor:** pt9912. **Datum:** 2026-06-29.

---

## 1. Ziel

`adr-immutable-check.sh` ([ADR-0016](../../adr/0016-adr-immutable-gate.md))
erzwingt die ADR-Immutabilität über einen **git-Diff** (`core(BASE)` vs.
`core(HEAD)` über eine Commit-Range) — volle Garantie, aber als **Skript** nur per
Kopie in Schwester-Repos nutzbar (Copy-Drift). Der hermetische `immutable`-Pin
([`slice-052`](../done/slice-052-immutable-modul.md)) löste die **verteilbare**
Hälfte (neu-pinn-bar, schwächere Garantie). **Neu:** dieselbe **harte** git-Garantie
verteilbar im Image — ein opt-in Modul `vcs` vergleicht `core(BASE)` ≟ `core(HEAD)`
über `--range <base>..<head>` (CI/Push) bzw. `--staged` (pre-commit) und meldet
Abweichung / unzulässigen Status-Übergang / Löschung als `core-drift-vcs`. Es liest
`.git` über einen reine-Go-**VCS-Port** — **ohne** git-Binary (distroless bleibt),
**ohne** Netz, read-only; `.git` ist im bestehenden Mount schon da. d-check
**dogfooded** das Modul für seine eigenen `Accepted`-ADRs (das Gate `adr-check`
läuft darauf um).

## 2. Entscheidungen

- **Eigenes Modul `vcs`** (13.), nicht Strategie in `immutable`: anderer
  Bindepunkt (Commit-Range statt Arbeitsbaum), anderer Eingabe-Scope (git statt
  Datei), anderer Lebenspunkt (Integration/PR statt `make gates`). `immutable`
  bleibt sauber hermetisch; beide koexistieren als **Defense-in-Depth**.
- **Volle Parität `--range` + `--staged`.** Echter Drop-in-Ersatz für
  `adr-immutable-check.sh`: CI-Range **und** lokaler `pre-commit`-staged-Modus,
  beide Quadranten wie das Skript.
- **Reine-Go-VCS-Port, distroless bleibt.** Ein neuer Driven-Port
  (`FileAtRef`/`ChangedFiles`) + einziger git-berührender Adapter (`git`),
  parallel zu `httpcheck`/`external`. Der Adapter liest `.git` über eine reine-Go
  git-Objekt-Bibliothek — **kein** git-Binary im Image
  ([ADR-0002](../../adr/0002-distribution-ghcr-image.md) unangetastet);
  `make arch-check` (R1–R5) hält die Bibliothek im Adapter isoliert, der Kern
  bleibt rein/fakebar.
- **Lokal/deterministisch/read-only** — `vcs` relativiert **weder**
  [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus) **noch**
  [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit):
  Eingabe ist das lokale `.git` (reproduzierbar, netzlos, schreibfrei), nur der
  **Eingabe-Scope** ist erweitert. Darum strikt opt-in + **fail-closed** ohne
  `.git`/Range (Exit 2, kein stilles Grün).
- **Core in Skript-Parität.** `vcs.exclude-sections` strippt Abschnitte
  (`Geschichte`); `vcs.status-line` strippt **nur die Kopf-Status-Zeile** (eine
  gleichlautende Körper-Zeile bleibt Core); `vcs.head-allow` erzwingt den
  zulässigen Status-Übergang. Normalisierung + SHA-256 wie `immutable`/`pins`.
- **Dogfood-Ersatz.** `make adr-check`, `pre-commit`-Hook und CI rufen das
  d-check-Image mit `--enable vcs` (`--range`/`--staged`). Die `vcs`-Klasse liegt
  als `vcs:`-Block in [`.d-check.yml`](../../../../.d-check.yml), **außerhalb** der
  Default-`modules`-Liste → `make doc-check` aktiviert `vcs` nicht (default-aus
  byte-identisch).
- **Skript bleibt pfad-stabil.** `tools/adr-immutable-check.sh` wird **nicht
  gelöscht**: [ADR-0016](../../adr/0016-adr-immutable-gate.md) ist immutable und
  referenziert es in Inline-Code (Löschung bräche `make doc-check` mit
  `codepath-missing`). Es bleibt als Referenz/Fallback, aus dem Gate ausgehängt.

## 3. Definition of Done

- [x] **Spec/Doc (doc-first):** neue VCS-Anforderung
  [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  (Bereichskürzel `VCS` in §3, Versions-Bump 0.33.0 + §7-Historie, `vcs` in
  [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
  + Glossar) + [ADR-0024](../../adr/0024-vcs-immutable-gate.md) (Proposed, + ADR-Index-Zeile;
  die [ADR-0016](../../adr/0016-adr-immutable-gate.md)-Teil-Supersede-Annotation folgt
  erst bei Closure mit dem Accept) + spezifikation-`.a`-Sektion + Grund-Code
  `core-drift-vcs` (§4) + Schema-Keys `vcs.*` (§2) + architecture VCS-Port-Rolle.
- [ ] **Code:** VCS-Port (`port/driven/vcs.go`), git-Adapter (`adapter/driven/git/`,
  reine-Go, einzige git-Tür), Kern-Regel (`rules/vcs.go`: Diff der Range,
  Klassen-/`immutable-when`-Filter, Core in Skript-Parität, Vergleich →
  `core-drift-vcs`), Wiring in `run.go` + CLI-Flags `--range`/`--staged`,
  `model.VCSConfig` + `validModules()` + `configyaml`-Kompilierung, Grund-Code in
  `diagnose.go`, neue Dependency in `go.mod`/`go.sum` (Pin-Disziplin
  [ADR-0011](../../adr/0011-digest-pins-build-gate-images.md)). Tests (alle 7
  Skript-Selbsttest-Klassen): Geschichte-Anhang/Superseded-Übergang/Körper-Edit auf
  `Proposed`-BASE (`immutable-when` nicht erfüllt) feuern **nicht**; Körper-Edit/
  Status-Rückfall/Körper-`**Status:**`-Edit/Abschnitt-nach-Geschichte feuern; dazu
  Löschung/Umbenennung; fail-closed (kein `.git`/Range); Modul-aus byte-identisch
  (Fake-VCS-Port).
- [ ] **Gate-Umbau:** `make adr-check`/`pre-commit`/`ci.yml` auf `--enable vcs`
  umstellen; `vcs:`-Block in `.d-check.yml`; `gate-consistency` + Doku
  (`harness/README.md` §Sensors, `AGENTS.md` §4) nachziehen; Negativ-Selbsttest
  erhalten.
- [ ] `make ci`/`make fullbuild` grün; zwei unabhängige Reviews; Closure (Move
  nach `done/` + Roadmap-Flip,
  [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise));
  bei Closure: [ADR-0024](../../adr/0024-vcs-immutable-gate.md) → Accepted **und** die
  [ADR-0016](../../adr/0016-adr-immutable-gate.md)-Teil-Supersede-Annotation nachziehen
  (Geschichte-Append + Index — erst dann ist der Supersede real); Release v0.33.0.

## 4. Risiken / offene Punkte

- **Neue Dependency (reine-Go-git).** Vergrößert `go.mod`/`go.sum`, Image und
  `semgrep`-Scope; unterliegt der Pin-Disziplin
  ([ADR-0011](../../adr/0011-digest-pins-build-gate-images.md)). Import **nur** im
  `git`-Adapter (arch-check R2-Klasse). Mitigation: hinter dem VCS-Port gekapselt,
  fakebar.
- **`adr-check` braucht jetzt das Image** (Hook + `make adr-check` bauen/rufen das
  Image statt eines Shell-Skripts) — schwergewichtiger pro Commit; dafür verteilte
  Wahrheit. Die CI baut das Image ohnehin.
- **Bootstrap-Reihenfolge des Gate-Umbaus:** der neue `adr-check` darf erst grün
  greifen, wenn das Image das `vcs`-Modul trägt; bis dahin Übergang sauber halten
  (Skript bleibt als Fallback verfügbar).
- **fail-closed vs. Komfort:** ohne `.git`/Range bricht `vcs` laut ab (Exit 2) —
  bewusst, damit eine fehlende git-Eingabe nicht still grün ist.
- **Read-only `.git` + reine-Go-Bibliothek:** nur lesende Plumbing-Operationen
  (Commit-Auflösung, Blob-an-Pfad, Tree-Diff); keine Index-/Working-Tree-Schreibung.

## 5. Trigger

Auftraggeber 2026-06-29: „Können regelkonform `adr-immutable-check` angehen" — die
von [ADR-0023](../../adr/0023-immutable-core-pin.md) als „nicht verworfen, nur
vertagt" markierte git-Stufe (Slice B / „git-Weg") wird eingelöst. Vier
Entscheidungen geklärt (eigenes Modul; volle Parität; Dogfood-Ersatz; reine-Go-
git-Zugang).

## 6. Sub-Area-Modus-Begründung

GF (Produkt-Code + Spec; „Doc führt, Code folgt"). Der git-Adapter ist eine
GF-Erweiterung der Hexagon-Adapter-Schicht (analog `httpcheck`); keine BF-Sub-Area.

## 7. Closure-Notiz (nach `done/`)

_(folgt bei Closure)_
