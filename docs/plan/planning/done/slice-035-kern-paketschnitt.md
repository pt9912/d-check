# Slice slice-035: Kern-Paketschnitt `model` / `rules` / `app`

**Status:** done (Closure 2026-06-20).

**Welle:** welle-24-kern-paketschnitt (Trigger: Nutzer-Entscheid „3-Wege
model/rules/app", 2026-06-20 — der Kern ist auf 5.212 Z. in einem Paket
gewachsen).

**Bezug:**
[ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md)
(Entscheidung),
[ADR-0004](../../adr/0004-architektur-pattern-hexagonal.md) („Hexagon
light", geschärft),
[ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md) (arch-check,
erweitert),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(kein Befund-Delta),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(I/O-Freiheit bleibt).

**Autor:** pt9912. **Datum:** 2026-06-20.

---

## 1. Ziel

`internal/hexagon/core` (ein Paket, 5.212 Nicht-Test-Z.) in drei Pakete mit
strikt einbahniger Importrichtung schneiden — **rein strukturell, kein
Verhaltens-/Befund-Delta**. Schichtrichtung Daten ← Regeln ← Orchestrierung
wird maschinell erzwungen.

## 2. Validierter Schnitt (Analyse 2026-06-20)

- **`model`** — `finding.go`, `config.go` (Blatt, importiert nichts aus dem
  Kern).
- **`rules`** — `links`/`anchors`/`ids`/`matrix`/`codepaths`/`spans`/
  `hostpaths`/`external.go` + `markdown.go` + `paths.go` (+ Glob-Helfer
  `ignored`) **+ `run.go` + `scan.go`** (Prüf-Orchestrierung/Discovery — bei
  der Umsetzung hierher statt `app`: die Modul-Tests koppeln `Run` mit
  Modul-Interna, so bleiben sie ohne Interna-Export testbar). Importiert
  `model` + Ports.
- **`app`** — `diagnose.go`, `repair.go`, `suggest.go` (Anwendungs-Modi).
  Importiert `model` + `rules` + Ports.
- **Belegt azyklisch:** `app` nutzt die `check*`-Regelfunktionen
  (`app → rules`); kein `rules`-Symbol referenziert `app` **außer**
  `ids.go → ignored()` (deshalb der `ignored`-Umzug nach `paths.go`);
  `model` ist Blatt; `port/driven` importiert keine Kern-Typen.

## 3. Definition of Done

- [x] `model`/`rules`/`app`-Unterpakete angelegt, alle 17 Nicht-Test-Dateien
  umpaketiert; `ignored` (+ ggf. weitere von `rules` geteilte reine Helfer)
  nach `paths.go` verschoben.
- [x] ~13 `_test.go` folgen ihrem Paket (white-box); der geteilte
  `memfs_test.go`-FS-Helfer wird ein von allen drei Test-Paketen nutzbarer
  exportierter Helfer.
- [x] Adapter (`report`, `configyaml`, …) + CLI auf die neuen Importpfade
  umgestellt (`…/core/model` bzw. `…/core/app`).
- [x] `tools/arch-check.sh` kodiert die neue Importrichtung
  (`app → rules → model`, `rules/app → port`; verboten `model→*`,
  `rules→app`, `port→{rules,app}`); R1–R5 nachgezogen; Selbsttest grün.
- [x] [`spec/architecture.md`](../../../../spec/architecture.md) §Kern auf
  die drei Ringe nachgezogen.
- [x] **Kein Befund-Delta:** Befundsatz auf einem Fixture vor/nach
  byte-identisch ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- [x] `make gates` grün; unabhängiges Review R1;
  [ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md) auf
  `Accepted`; Closure-Notiz.

## 4. Plan (vor Code)

| Datei/Bereich | Art | Begründung |
|---|---|---|
| neue Pakete `model`/`rules`/`app` unter `internal/hexagon/core/` | neu/move | Drei-Paket-Schnitt; Dateien per `git mv` (Lifecycle-Spur des Codes) |
| `internal/hexagon/core/*_test.go` | move | white-box-Tests folgen ihrem Paket; `memfs`-Helfer exportieren |
| `internal/adapter/**`, `cmd/d-check` | update | neue Importpfade |
| `tools/arch-check.sh` | update | Importrichtung der drei Pakete |
| `spec/architecture.md` | update | §Kern in drei Ringe |
| `docs/plan/adr/0012-…`, `docs/plan/adr/README.md` | update | ADR auf Accepted, Index |

## 5. Risiken und offene Punkte

- **`memfs_test.go` (geteilter Helfer):** von Tests mehrerer künftiger
  Pakete genutzt → entweder exportierter Helfer in einem neutralen
  (Test-)Paket oder je Paket dupliziert; Variante in der Umsetzung
  festlegen.
- **Import-Zyklus-Risiko:** der Go-Compiler ist die harte Schranke; vor dem
  großen Move den `ignored`-Umzug zuerst, dann paketweise umstellen.
- **Kein Verhaltens-Delta:** der Refactor darf den Befundsatz nicht ändern —
  Vorher/Nachher-Beleg auf einem Fixture ist Pflicht
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- **Größe:** der umfangreichste Einzel-Refactor des Repos — eigener Slice,
  in einem fokussierten Lauf umzusetzen.

## 6. Closure-Trigger

DoD vollständig inkl. grüner Gates (Befundsatz byte-identisch), Review R1
und [ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md) auf
`Accepted`.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung.** Der Kern (5.212 Z., ein Paket) ist in drei Pakete unter
`internal/hexagon/core/` geschnitten: `model` (finding/config),
`rules` (8 Module + markdown/paths + **run/scan**), `app`
(diagnose/repair/suggest). Importrichtung `app → rules → model` (+ Blatt
`port/driven`), per `tools/arch-check.sh` R6 erzwungen. Vorstufen committet:
`memFS` → Paket `coretest` (`5e78237`); Config-Typen-Zyklus
(`CodepathsConfig`/`HostpathsConfig` → config.go, `bf567d4`); 17 Symbole
exportiert + `ignored` nach paths.go (`72f38e8`); der Split selbst
(`4e7d26d`).

**Abweichung vom Plan.** `run.go`/`scan.go` liegen in `rules`, nicht `app`:
die Modul-Tests (White-box) koppeln den Orchestrator `Run` mit Modul-Interna
(`statusOf`, `classifyCodepath`, `htmlAnchors`, …); in `rules` bleiben sie
testbar, ohne Interna zu exportieren. „Engine" (Module + Ausführung) vs.
„Modi" (doctor/repair/suggest) ist die vom Code getragene Naht.
[ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md) +
`spec/architecture.md` §Kern nachgezogen.

**Belege.**
- `make gates` grün inkl. `arch-check` R1–R6 und `coverage-gate` 93 %.
- **Alle Tests bestehen** → kein Verhaltens-Delta
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus));
  der `doc-check`-Selbstscan liefert 0 Befunde wie vor dem Schnitt.
- 40 stale Code-Pfade in historischen Docs (done-Slices/Reviews) auf die
  neuen Paketpfade nachgezogen.

**Lerneintrag.**
- Der Make-only-Guard blockiert `gofmt -r` (auch containerisiert) → die
  ~80 Cross-Refs liefen per `sed` (geprüft kollisions-/string-sicher,
  Compiler als Netz). Ein AST-Tool-Target wäre für künftige Refactors nützlich.
- Tests zeigen die natürlichen Paket-Nähte: White-box-Kopplung Test↔Interna
  bestimmt, was zusammen bleiben muss (run/scan zu den Modulen).
- Reine Verschiebungen müssen die Referenzen in **historischer** Doku
  mitziehen (codepaths-Gate) — „Referenz folgt dem Code".

**Review-Runde R1** (`docs/reviews/2026-06-20-slice-035-kern-paketschnitt.md`):
**merge-fähig** (0 HIGH / 0 MEDIUM, 1 LOW, 3 INFO). F1/F2 (`MatchGlob`/
`Ignored` über-exportiert — 0 Konsumenten außerhalb `rules`, da `run`/`scan`
dort liegen) → un-exportiert; F3 (stale `.golangci.yml`-Kommentar
`stripInlineCode/matchGlob`) → behoben; F4 (architecture.md nennt Modus-Flags)
won't-fix (von
[ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md) sanktioniert).

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Areas GF (Kern-/Tooling-/Doku-Arbeit; Greenfield-Default);
rein struktureller Umbau mit Vorher/Nachher-Beleg.
