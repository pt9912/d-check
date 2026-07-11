# Slice slice-066: RTM-Quellen und Kennungs-Muster konfigurierbar (`trace`-Block)

**Status:** in-progress (welle-55-trace-konfigurierbare-quellen). Doc-first-
Fundament gelegt (Lastenheft-CR, Spezifikation, [ADR-0034](../../adr/0034-trace-konfigurierbare-quellen.md)
Proposed); Implementierung + Review + Release folgen. Lifecycle
`in-progress`→`done` mit Roadmap-Flip §Aktuelle Welle bei Closure
([`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)).

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
- [ ] **Modell** [`config.go`](../../../../internal/hexagon/core/model/config.go):
  `TraceConfig` (Regex + Strings, Nullwert = Default); `Config` trägt `Trace`.
- [ ] **Config-Decode** [`configyaml.go`](../../../../internal/adapter/driven/configyaml/configyaml.go):
  `trace`-Top-Level-Key; Regex-Kompilierung + Capture-Gruppen-Guard + Pfad-
  Validierung (Exit 2).
- [ ] **RTM** [`trace.go`](../../../../internal/hexagon/core/app/trace.go):
  `resolveTrace` (Default-Auflösung) + `BuildTraceMatrix(fsys, tc)`; `traceRequirements`/
  `traceRefs`/`isFullReqID` parametrisiert.
- [ ] **CLI** [`cli.go`](../../../../internal/adapter/driving/cli/cli.go):
  `loadConfig` **vor** dem `--trace`-Zweig; `cfg.Trace` durchreichen; Config-Fehler
  ⇒ Exit 2.
- [ ] **`--print-config`** [`config_template.go`](../../../../internal/adapter/driving/cli/config_template.go):
  kommentierter `trace`-Block.
- [ ] **Tests**: Default byte-identisch; Fremd-Konvention (grid-gym-Gestalt,
  Vorher/Nachher); Voll-Custom-Konvention (alle 8 Achsen); Negative-Config (Regex /
  Capture / Pfad-Escape ⇒ Exit 2); `--require-complete`-Vererbung; mutations-hart.
- [ ] **Release-Prep**: Benutzerhandbuch (`trace`-Config + Beispiel), `operations.md`,
  `CHANGELOG.md`; slice-061/062-Harnesse grün; bare-Tag-Sweep `v0.39.0`→`v0.40.0`
  + `version.md`-Register nach [`releasing.md` §4](../../../../docs/user).
- [ ] `make gates` / `make ci` grün; **ein unabhängiger Impl-Review**; Closure-Move +
  Body + **Lerneintrag**. **Release v0.40.0** (Push → CI → Tag → GHCR → digest-backfill).

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
