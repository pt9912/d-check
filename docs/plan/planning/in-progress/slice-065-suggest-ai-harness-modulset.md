# Slice slice-065: ai-harness-Vorlage — Modulset an die gelebte Konvention angleichen

**Status:** in-progress (welle-54-suggest-ai-harness-modulset). Move
`next/`→`in-progress/` + Roadmap-Flip §Aktuelle Welle vollzogen
([`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise));
[ADR-0033](../../adr/0033-ai-harness-template-modulset.md) bleibt **Proposed** bis
zur Closure (ADR-Annotation erst bei Closure).

**Welle:** welle-54-suggest-ai-harness-modulset (unabhängig von
[`slice-064`](../next/slice-064-gate-consistency-tombstone.md), Reihenfolge offen).

**Bezug:** [ADR-0033](../../adr/0033-ai-harness-template-modulset.md)
(Template-Eignungs-Kriterium K1–K4 + die konkreten Modul-Zuordnungen) +
[`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)
(geänderte Anforderung; Spezifikation
[§`DC-FA-CLI-006.a`](../../../../spec/spezifikation.md#dc-fa-cli-006a--konfigurations-vorschlag)).
**Lastenheft-CR** (geändertes Modulset ist nutzersichtbar) — **kein** neues
Modul, **keine** neue `DC-*`-ID, **kein** neuer Grund-Code. **Release** geplant
(v0.39.0, geänderte `--suggest-config`-Ausgabe).

**Autor:** pt9912. **Datum:** 2026-07-06.

---

## 1. Ziel

Die `--suggest-config ai-harness[-init]`-Vorlage emittiert nur fünf Module
(`links, anchors, ids, matrix, codepaths`), während d-checks eigener Dogfood
(`.d-check.yml`) im Default-Scan **acht** führt (`+ spans, hostpaths, versions`)
und `planning`/`targets`/`vcs`/`commits` als Gates konfiguriert. Die Vorlage
predigt also weniger, als die Konvention lebt. Zugleich deckt der operative
§-Body die tatsächlich emittierte Ausgabe (Kommentarzeile + `codepaths`-Block)
nicht — der Code läuft der Norm voraus. Ziel: das Modulset an die gelebte
Konvention angleichen und die Norm die Ausgabe 1:1 decken lassen.

## 2. Entscheidungen (aus [ADR-0033](../../adr/0033-ai-harness-template-modulset.md))

- **Eignungs-Kriterium K1–K4** (konventions-kanonisch · ableitungsfrei/konventions-
  feste Config · Baum-Scan-tauglich · hermetisch/netzlos) macht die
  Modul-Zuordnung explizit statt sie im Code zu vergraben.
- **`spans` + `hostpaths` ins fixe Modulset** (beide Modi): erfüllen K1–K4
  (konfigfrei, hermetisch, d-check führt sie). Neues fixes Set:
  `links, anchors, ids, matrix, codepaths, spans, hostpaths`. **Revidiert** die
  slice-046-Einordnung dieser zwei als „situativ nicht aktiviert".
- **`planning` als repo-bewusster Block** (analog `matrix`): einzige Config der
  `roadmap`-Pfad (`docs/plan/planning/in-progress/roadmap.md`),
  `heading`/`marker`/`slice-glob` sind Defaults. Aktiv bei vorhandener Roadmap +
  `docs/plan/planning/` (`ai-harness`), sonst auskommentiert mit Hinweis; im
  Voll-Kanon (`ai-harness-init`) aktiv.
- **`vcs`/`commits` → `--print-mk`-Verweis, NICHT ins Modulset**: scheitern an K3
  (Laufzeit-Range). Die „Weitere opt-in-Module"-Kommentarzeile verweist
  zusätzlich auf `--print-mk` (`doc-immutable`/`doc-commits`).
- **`versions`/`targets` draußen, aber dokumentiert vertagt**: `versions` scheitert
  an K2 (repo-spezifisches `pin-pattern`), `targets` liegt an der K2-Grenze
  (`authority`/`doc-tables` semi-spezifisch). Beide werden im §-Body als **geprüft
  und begründet vertagt** benannt (nicht still im Code entschieden). `targets` ist
  späterer repo-bewusster-Block-Kandidat.
- **`external`/`diagrams`/`pins`/`immutable`/`tracked` bleiben draußen** (K4 Netz /
  K2 unableitbar / pro-Marker nichts zu deklarieren / K3-K4 fail-closed) —
  unverändert zu slice-046.
- **Norm deckt Ausgabe 1:1**: die kanonische Vorlage der Spezifikation bekommt
  die „Weitere opt-in-Module"-Kommentarzeile (inkl. `--print-mk`-Verweis) und den
  `codepaths`-Block; der operative Body benennt K1–K4 + die **geschlossene**
  Aktiv-Menge.

## 3. Definition of Done

- [ ] **Lastenheft-CR** in [`spec/lastenheft.md`](../../../../spec/lastenheft.md):
  der Anforderungs-Body — Standard-Modulset auf
  `links, anchors, ids, matrix, codepaths, spans, hostpaths` + repo-bewusster
  `planning`-Block; K1–K4 + geschlossene Aktiv-Menge benannt; `versions`/`targets`
  als vertagt vermerkt. **Versions-Bump + Historie-Zeile** (v0.39.0). Anlege-/
  Änderungs-Prozess nach [`harness/conventions.md` §Anforderungs-Anlege-Prozess](../../../../harness/conventions.md#anforderungs-anlege-prozess).
- [ ] **Spezifikation** [§`DC-FA-CLI-006.a`](../../../../spec/spezifikation.md#dc-fa-cli-006a--konfigurations-vorschlag):
  fixes Modulset in Prosa + kanonischer YAML-Vorlage aktualisiert (spans/hostpaths
  in `modules`, `planning`-Block, `codepaths`-Block, „Weitere opt-in-Module"-
  Kommentarzeile mit `--print-mk`-Verweis) — die emittierte Ausgabe 1:1.
- [ ] **Renderer** [`internal/hexagon/core/app/suggest.go`](../../../../internal/hexagon/core/app/suggest.go):
  `renderHarness` nimmt `spans`/`hostpaths` ins fixe `modules:`; neuer
  `renderHarnessPlanning` (repo-bewusst, analog `renderHarnessMatrix`); die
  Kommentarzeile verweist auf `--print-mk` und nennt `versions`/`targets` als
  vertagt.
- [ ] **Tests** (`suggest_test.go` / `cli_acceptance_test.go`): fixes Set enthält
  spans/hostpaths (beide Modi); `planning`-Block aktiv **mit** Roadmap-Fixture,
  auskommentiert **ohne** (`ai-harness`); im Voll-Kanon (`ai-harness-init`) aktiv;
  Round-Trip — das dekodierte Gerüst läuft weiter über `configyaml.Decode`
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- [ ] **[ADR-0033](../../adr/0033-ai-harness-template-modulset.md) auf Accepted**
  + ADR-Index. **Kein** neuer Grund-Code, **keine** neue `DC-*`-ID.
- [ ] **Release-Prep** (geänderte Ausgabe ist nutzersichtbar): Benutzerhandbuch
  `--suggest-config`-Beispiel(e) + `CHANGELOG.md` nachziehen; die
  slice-061-Config-Beispiel- und slice-062-E2E-Harnesse gegen die neue Vorlage
  grün halten (die emittierte YAML dekodiert weiter). [`release-prep`-Blindspots](../../../../docs/user)
  prüfen.
- [ ] `make gates` / `make ci` grün; **ein unabhängiger Impl-Review**;
  Closure-Move + Body + **Lerneintrag**. **Release v0.39.0** (Push → CI → Tag →
  GHCR → digest-backfill).

## 4. Trigger

Freigabe durch den Auftraggeber (2026-07-06): aus der Analyse „welche Module
werden in `--suggest-config ai-harness` nicht verwendet, und warum nicht" wurde
der Normativitäts-Spalt (Code-Kommentar ≠ Norm) und die Diskrepanz Vorlage↔Dogfood
sichtbar; Zuschnitt bestätigt: **spans+hostpaths fix, planning repo-bewusst,
`--print-mk`-Verweis für vcs/commits, versions/targets dokumentiert vertagt**.

## 5. Offene Punkte / Risiken

- **Reihenfolge zu [`slice-064`](../next/slice-064-gate-consistency-tombstone.md):** beide
  liegen im Backlog, unabhängig. Der Implementer/Auftraggeber wählt die Welle-
  Aufnahme; kein technischer Zwang.
- **`ai-harness-init` fürs leere Repo:** `planning` aktiv im Voll-Kanon setzt die
  angelegte Roadmap-Struktur voraus (wie `ids`-Targets) — konsistent mit dem
  bestehenden „läuft, sobald die Struktur existiert"-Vertrag, in Body/Test benennen.
- **Doku-Harnesse (slice-061/062):** die neue emittierte YAML muss weiter
  dekodieren; falls ein `--suggest-config`-Beispiel im Handbuch steht, wandert es
  in die Verankerung — im Release-Prep prüfen.
