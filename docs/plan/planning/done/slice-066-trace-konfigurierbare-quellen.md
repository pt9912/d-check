# Slice slice-066: RTM-Quellen und Kennungs-Muster konfigurierbar (`trace`-Block)

**Status:** done (welle-55-trace-konfigurierbare-quellen). Lifecycle
abgeschlossen (`in-progress`→`done`, Roadmap-Flip §Aktuelle Welle,
[`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise));
[ADR-0034](../../adr/0034-trace-konfigurierbare-quellen.md) auf **Accepted**
(ADR-Annotation bei Closure). Ergebnis + Belege in §7.

**Welle:** welle-55-trace-konfigurierbare-quellen.

**Bezug:** [ADR-0034](../../adr/0034-trace-konfigurierbare-quellen.md)
(opt-in `trace`-Block, Design + Fail-closed) +
[`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
(geänderte Anforderung; Spezifikation
[§`DC-FA-CLI-009.a`](../../../../spec/spezifikation.md#dc-fa-cli-009a--requirements-traceability-matrix)).
Mit-berührt [`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
(erbt die konfigurierten Quellen). **Lastenheft-CR** (nutzersichtbare
Config-Fläche) — **kein** neues Modul, **keine** neue `DC-*`-ID, **kein** neuer
Grund-Code. **Release** geplant (v0.40.0, neue `trace`-Config).

**Autor:** pt9912. **Datum:** 2026-07-11.

---

## 1. Ziel

Die RTM (`--trace`,
[`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix))
ist an drei Achsen hart an d-checks eigene Konvention gebunden: den
Anforderungs-Kennungs-Regex (`-FA-`/`-QA-`-Mittelsegmente), den Slice-Dateinamen
(`slice-NNN-…md`) und den ADR-Dateinamen (`NNNN-…md`); die Quell-Pfade sind
Konstanten. Der Konsument **grid-gym** (243 Anforderungen `GG-<FAMILIE>-NNN`,
Slices `NNN-…md`) sieht darum via `make doc-trace` nur **6** Anforderungen —
allein die `GG-QA-*`-Familie trifft zufällig den `-QA-`-Zweig — und alle als
Waisen (kein `slice-NNN-…md` erkannt). Ein via `--print-mk` an jeden Konsumenten
verteiltes Feature funktioniert nur für Repos, die d-checks Schema nachahmen.

Ziel: einen **opt-in `trace`-Config-Block** einführen, der die vier
Konventions-Achsen überschreibt, mit Konventions-Defaults ⇒ d-checks eigener
Lauf byte-identisch, sodass die RTM auch Fremd-Konventions-Repos vollständig
abbildet.

## 2. Entscheidungen (aus [ADR-0034](../../adr/0034-trace-konfigurierbare-quellen.md))

- **Opt-in `trace`-Block** in `.d-check.yml`:
  `requirements.source`/`.id-pattern`; `adrs.dir`/`.file-pattern`/`.id-prefix`;
  `slices.dir`/`.file-pattern`/`.id-prefix`. Der `id-pattern`-Regex erkennt eine
  Kennung als **Ganz-Token** im Heading **und** als Referenz in ADR-/Slice-
  Dateien (ein Regex, beide Rollen — wie das heutige `reqIDFull`).
- **Jedes Feld optional; Default = heutige Konstante** über `Effective*()`-
  Methoden im Modell (analog `PlanningConfig`/`VersionsConfig`). Kein `trace`-
  Block ⇒ RTM **byte-identisch** ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- **Owner-Kennung** = Capture-Gruppe 1 der `file-pattern` + `id-prefix`
  (`ADR-` + `0013`; Slice-Präfix Default leer).
- **Design spiegelt `ids.patterns`** (dediziert, explizit gesetzt) — **nicht**
  aus `ids.patterns` abgeleitet (implizit/brüchig).
- **Fail-closed zur Config-Zeit** (Exit 2, vor jedem Scan): ungültige Regex
  (`id-pattern`/`file-pattern`) oder `file-pattern` **ohne Capture-Gruppe**
  (sonst `m[1]`-Panic).
- **`--print-config`** ergänzt einen kommentierten `trace`-Block
  (Verfügbar-Vollständigkeit). **Kein** neuer Eingabe-Scope; kein VCS/Netz —
  [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
  unberührt.

## 3. Definition of Done

- [x] **Lastenheft-CR** in [`spec/lastenheft.md`](../../../../spec/lastenheft.md):
  [`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)-Body
  um den `trace`-Block; AKs (Konfiguriert/Default-byte-identisch/Negative-Config);
  Out-of-Scope-Umkehr der 0.21.0-Zeile. **Versions-Bump + Historie** (0.40.0; der
  stale Header 0.38.0 → 0.40.0 mitkorrigiert — slice-065 hatte den Header-Bump auf
  0.39.0 versäumt).
- [x] **Spezifikation** [§`DC-FA-CLI-009.a`](../../../../spec/spezifikation.md#dc-fa-cli-009a--requirements-traceability-matrix):
  Schritte 1/2 auf konfigurierbare Quellen + Config-Auflösungs-Absatz;
  §2-Schema-Tabelle um die acht `trace.*`-Keys; §7-Historie.
- [x] **[ADR-0034](../../adr/0034-trace-konfigurierbare-quellen.md)** (Proposed) +
  ADR-Index-Zeile.
- [x] **Modell** [`config.go`](../../../../internal/hexagon/core/model/config.go):
  `TraceConfig` (Regex + Strings, Nullwert = Default); `Config` trägt `Trace`.
  Default-Auflösung in `trace.go`, nicht im globals-freien Modell (lint-sicher).
- [x] **Config-Decode** [`configyaml.go`](../../../../internal/adapter/driven/configyaml/configyaml.go):
  `trace`-Top-Level-Key; `compileTracePattern` (Regex + Capture-Guard) +
  `validateTracePath` — alles fail-closed (Exit 2).
- [x] **RTM** [`trace.go`](../../../../internal/hexagon/core/app/trace.go):
  `resolveTrace` + `BuildTraceMatrix(fsys, tc)`; `traceRequirements`/`traceRefs`/
  `isFullReqID` parametrisiert.
- [x] **CLI** [`cli.go`](../../../../internal/adapter/driving/cli/cli.go):
  `loadConfig` **vor** dem `--trace`-Zweig; `cfg.Trace` durchgereicht; Config-Fehler
  ⇒ Exit 2.
- [x] **`--print-config`** [`config_template.go`](../../../../internal/adapter/driving/cli/config_template.go):
  kommentierter `trace`-Block (in `TestCLI053` verankert).
- [x] **Tests**: Default byte-identisch; Fremd-Konvention (Vorher/Nachher);
  Voll-Custom (alle 8 Achsen, mutations-hart); Negative-Config (Regex / Capture /
  Pfad-Escape ⇒ Exit 2); `--require-complete`-Vererbung.
- [x] **Release-Prep**: Benutzerhandbuch §4.12+§5, `operations.md`, `CHANGELOG.md`
  `[0.40.0]`; bare-Tag-Sweep `v0.39.0`→`v0.40.0` + `version.md`-Register;
  slice-061/062-Harnesse grün.
- [x] **Verifikation** gegen grid-gym-Realdaten: **6 → 243** Anforderungen sichtbar.
- [x] `make gates` grün; **unabhängiger Impl-Review R1 ACCEPT-WITH-NITS** (§6).
- [ ] **Offen:** `make ci` (image-test) + Closure-Move + Body + **Lerneintrag**;
  **Release v0.40.0** (Push → CI → Tag → GHCR → digest-backfill).

## 4. Trigger

Konsumenten-Befund grid-gym (2026-07-11): `make doc-trace` meldete 6 von 243
Anforderungen (nur die `GG-QA-*`-Familie traf d-checks `-QA-`-Default). Nutzer-
Entscheid: RTM-Quellen konfigurierbar machen (Richtung „d-check konfigurierbar"
statt „Feature in grid-gym entfernen").

## 5. Offene Punkte / Risiken

- **Byte-Identität im Default:** der stärkste Regressionsschutz ist ein Test, der
  d-checks **eigene** RTM (bzw. eine Fixture-RTM) mit und ohne leeren `trace`-
  Block vergleicht — beide identisch. Priorität.
- **Capture-Gruppen-Guard:** `file-pattern` ohne `(...)` würde heute `m[1]`
  panicen — der Guard ist Pflicht (Negative-Test), nicht Kosmetik.
- **Prefix-Default-Semantik:** leerer Wert = Default (`ADR-` bzw. leer für
  Slices). Ein leerer **ADR**-Präfix ist bewusst nicht ausdrückbar (Out-of-Scope,
  Lastenheft) — konsistent mit `PlanningConfig`s Empty-⇒-Default-Muster.
- **`--require-complete`:** erbt die konfigurierten Quellen automatisch (gleicher
  `BuildTraceMatrix`-Lauf) — im Test mit-abdecken.
- **Doku-`yaml`-Harness (slice-061):** ein `trace`-Beispiel in der Spezifikation/
  im Handbuch muss über `configyaml.Decode` laufen; darum das Schema als **Tabelle**
  dokumentiert (kein gefencter Config-Block in der Spezifikation), Handbuch-Beispiel
  erst nach der Decode-Implementierung.

## 6. Review-Nachtrag (Impl-R1)

Unabhängiger Impl-Review (24 Tool-Uses, alle geänderten/neuen Dateien + voller
Diff real gelesen) — **Verdikt ACCEPT-WITH-NITS**, beide nicht-blockierenden
Befunde eingearbeitet (Report:
[`docs/reviews/2026-07-11-slice-066-trace-konfigurierbare-quellen.md`](../../../reviews/2026-07-11-slice-066-trace-konfigurierbare-quellen.md)):

- **F-1 (MEDIUM):** Nur 2 der 8 Config-Achsen waren differenziell getestet; die
  übrigen 6 (`requirements.source`, `adrs.dir`/`.file-pattern`/`.id-prefix`,
  `slices.dir`/`.id-prefix`) liefen im Default-Test zwar durch, waren aber
  mutations-inert (ein gelöschter Override-Zweig bliebe ungefangen). Behoben mit
  `TestCLI066_Trace_VollCustomKonvention` — eine vollständig eigene Konvention
  (custom Quelle/Regex/Verzeichnisse/Dateimuster/Präfixe) assertet die
  aufgelösten Owner-Kennungen (`DEC-0007`/`T5`); jeder Default-Rückfall macht ihn
  rot.
- **F-2 (LOW):** der kommentierte `--print-config`-`trace`-Block war in
  `TestCLI053`s Präsenzliste nicht verankert (Entfernen bliebe ungefangen).
  `"# --- trace:"` ergänzt.
- **F-3 (INFO):** `loadConfig` läuft nun VOR dem `--trace`-Zweig — eine defekte,
  trace-fremde `.d-check.yml` liefert bei `--trace` jetzt Exit 2 statt einer RTM.
  Beabsichtigt und in [ADR-0034](../../adr/0034-trace-konfigurierbare-quellen.md)
  §Konsequenzen dokumentiert (fail-closed); nur notiert.

Explizit sauber verifiziert: Byte-Identität im Default (zeilenweise gegen die
alten hart kodierten `traceRefs`-Aufrufe), Capture-Guard greift für **beide**
`file-pattern` (kein `m[1]`-Panic-Pfad), `--require-complete` erbt die Config,
Norm↔Code-Parität (8 `rawTrace`-Keys == 8 §2-Schema-Zeilen, Default-Werte
deckungsgleich), Prefix-Semantik konsistent, Determinismus/Read-only,
Harness-Regeln (ADR nennt den Slice nur in `## Geschichte`, Versions-Header
0.40.0 == §7-Top == CHANGELOG), gocyclo unter Schwelle.

## 7. Closure-Notiz (nach done)

**Umgesetzt:** Die Requirements Traceability Matrix (`--trace`,
[`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix))
liest ihre vier Quell-Achsen jetzt aus einem opt-in `trace`-Block der
`.d-check.yml`: Anforderungs-Quelldatei + Kennungs-Regex sowie je Referenzklasse
(ADR/Slice) Verzeichnis, Basisnamen-Regex (Capture-Gruppe 1 = Owner-Kennung) und
Owner-Präfix. Jedes Feld ist optional; der Nullwert löst auf d-checks
Konventions-Default auf ⇒ RTM **byte-identisch**
([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
Fail-closed zur Config-Zeit (ungültige Regex / `file-pattern` ohne Capture-Gruppe
/ Pfad-Escape ⇒ Exit 2). Kehrt die 0.21.0-Out-of-Scope-Zeile „frei
konfigurierbare Quell-Pfade" um; kein neues Modul, keine neue `DC-*`-ID, kein
neuer Grund-Code (`--trace` bleibt advisory).

**Belege:** `make gates` grün (doc-check + lint + test + arch-check +
coverage-gate 93,7 % + semgrep + gate-consistency + planning-check).
**End-to-End gegen grid-gyms Realdaten** verifiziert: ohne `trace`-Block **6**
Anforderung(en) (nur die `GG-QA-*`-Familie trifft d-checks `-QA-`-Default), mit
`trace`-Block (`id-pattern: 'GG-…-\d{3}'`, `slices.file-pattern: '^(\d+)-…'`)
**243** Anforderung(en) mit aufgelösten ADRs/Slices. Impl-Review R1
ACCEPT-WITH-NITS (§6; MEDIUM-Test-Deckung + LOW eingearbeitet).

**Commit-Kette:** `4971b11` (doc-first) · `029fb61` (feat) · `09b945a`
(release-prep v0.40.0) · `bc8411c` (Review) · `1bb5e3f` (Closure-Move) ·
`dcfa44c` (Closure-Body) · digest-backfill. **Release v0.40.0** auf GHCR
(Release-Run 29143442890 grün), Digest-Pin
`ghcr.io/pt9912/d-check@sha256:e691053abd820f85e652a343f3d700ba135f2d8d66523151e1388c353af2ba75`.

**Lehre:** Ein via `--print-mk` an alle Konsumenten verteiltes Feature
(`doc-trace`/`doc-complete`) muss konfigurierbar sein, sonst bedient es nur die
Repo-Familie, die d-checks internes Schema zufällig teilt — der Konsumenten-Befund
(grid-gym 6/243) machte die Konventions-Bindung sichtbar. `file-pattern` mit
Capture-Gruppe braucht einen `NumSubexp()`-Guard, sonst `m[1]`-Panic. Ein
Default-„byte-identisch"-Test ist mutations-inert, wenn er alle Achsen auf ihre
Defaults setzt — ein **Voll-Custom**-Test (jede Achse non-default) fängt einen
weggebrochenen Override-Zweig (R1-MEDIUM).
