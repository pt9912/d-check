# Slice slice-033: Image-Pins auf Digest (Dockerfile + semgrep)

**Status:** done (Closure 2026-06-20).

**Welle:** welle-22-digest-pins (Trigger: Review R1 zu slice-032,
INFO-1 — `docs/reviews/2026-06-20-slice-032-semgrep-gate.md`; plus
[ADR-0002](../../adr/0002-distribution-ghcr-image.md)↔Dockerfile-Drift).

**Bezug:**
[ADR-0011](../../adr/0011-digest-pins-build-gate-images.md) (Entscheidung —
Digest-Pins aller Build-/Gate-Images),
[ADR-0002](../../adr/0002-distribution-ghcr-image.md) (Distribution/Build,
§1 Digest-Pins),
[ADR-0010](../../adr/0010-semgrep-hermetisches-gate.md) (semgrep-Image-Pin,
INFO-1),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(Determinismus — identische Eingabe ⇒ identische Befunde).

**Autor:** pt9912. **Datum:** 2026-06-20.

---

## 1. Ziel

Alle Container-Image-Pins von **veränderlichen Tags** auf
**`@sha256:`-Digests** umstellen, damit die Reproduzierbarkeit
([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)) nicht
an einem unter gleichem Tag neu gepushten Image hängt. Schließt die
Restlücke aus Review R1/slice-032 (INFO-1) und die [ADR-0002](../../adr/0002-distribution-ghcr-image.md)↔Dockerfile-Drift.

## 2. Ausgangslage (Ist)

- **`Dockerfile`** pinnt die Basis-Images nur über Tags, **kein** Digest:
  `golang:${GO_VERSION}`, `golangci/golangci-lint:${GOLANGCI_LINT_VERSION}`,
  `gcr.io/distroless/static-debian12:nonroot`.
- **`tools/semgrep.sh`** pinnt das Scanner-Image über den Tag
  `semgrep/semgrep:1.167.0` (Review-INFO-1).
- **`docs/plan/adr/0002-distribution-ghcr-image.md`** beschreibt den
  Build-Stage bereits als `FROM golang:<ver>@sha256:…` — der `Dockerfile`
  löst das aber nicht ein (Drift Doku↔Code).
- **Bereits digest-gepinnt:** nur das *ausgelieferte* d-check-Image
  (`docs/user/releasing.md`, jede Roadmap-Release-Zeile).

## 3. Definition of Done (vorläufig)

- [x] `Dockerfile`: jede `FROM`-Zeile mit `@sha256:<digest>` **inline** neben
  dem Tag (Tag bleibt lesbar, Digest = Wahrheit); `make versions` zeigt sie
  via `FROM`-grep.
- [x] `tools/semgrep.sh`: Scanner-Image `semgrep/semgrep:1.167.0@${SEMGREP_DIGEST}`.
- [x] Digest-Politik in **neuer** [ADR-0011](../../adr/0011-digest-pins-build-gate-images.md)
  festgehalten ([ADR-0002](../../adr/0002-distribution-ghcr-image.md)/[ADR-0010](../../adr/0010-semgrep-hermetisches-gate.md)
  sind Accepted/immutable — Vereinheitlichung statt Edit); §1-Drift aufgelöst.
- [x] `make versions` zeigt die `FROM`-Digests + die semgrep-Image-Zeile
  (inkl. Digest); Update-Politik (bewusster Doppel-Edit Version+Digest)
  dokumentiert — analog `GO_VERSION`/`GOLANGCI_LINT_VERSION`.
- [x] `make gates` grün (inkl. `make image-test`, `make semgrep` offline);
  reproduzierbar (Build aus Manifest-Listen-Digests).
- [x] Unabhängiges Review R1; Closure-Notiz;
  [ADR-0011](../../adr/0011-digest-pins-build-gate-images.md) auf `Accepted`.

## 4. Risiken und offene Punkte

- **Digest-Beschaffung netzlos?** Das Auflösen Tag→Digest braucht einmal
  Netz (`docker buildx imagetools inspect` o. Ä.) — Setup, nicht Gate.
  Multi-Arch: ggf. den Manifest-Listen-Digest pinnen, sonst bricht der
  Build auf abweichender Host-Architektur.
- **Pflegelast:** Digest-Hebung bei jedem Basis-Image-Update; Tag als
  Klartext daneben halten, damit Diffs lesbar bleiben.
- **Abgrenzung:** rein Reproduzierbarkeits-Härtung, kein neuer
  `DC-*`-Vertrag; kein Carveout (kein gebrochenes Gate, INFO-Niveau).

## 5. Trigger

Review R1/slice-032 INFO-1 (Tag statt Digest) +
[ADR-0002](../../adr/0002-distribution-ghcr-image.md)↔`Dockerfile`-Drift; vom
semgrep-Gate (slice-032) bewusst entkoppelt.

## 6. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Tooling-/Build-/Doku-Arbeit; Greenfield-Default).

## 7. Closure-Notiz (nach `done/`)

**Umsetzung.** Alle vier extern bezogenen Images per `@sha256:`-Digest
gepinnt, **inline neben dem Tag** (Tag lesbar, Digest = Wahrheit): die drei
`Dockerfile`-`FROM` (golang, golangci-lint, distroless) und das
Scanner-Image in `tools/semgrep.sh` (`SEMGREP_DIGEST`). `make versions` um
die semgrep-Image-Zeile ergänzt (portables `sed` statt `grep -oP`). Politik
in der **neuen** [ADR-0011](../../adr/0011-digest-pins-build-gate-images.md)
festgehalten — die Accepted/immutablen
[ADR-0002](../../adr/0002-distribution-ghcr-image.md)/[ADR-0010](../../adr/0010-semgrep-hermetisches-gate.md)
wurden vereinheitlicht, nicht editiert; die §1-Digest-Drift ist behoben.

**Belege.**
- Digests sind **Manifest-Listen** (`oci.image.index`) mit linux/amd64 +
  linux/arm64 (+ weitere) — kein plattformgebundener Pin, arm64-Build sicher.
- `make ci` grün: gates inkl. semgrep offline + `make image-test`
  (nativ == Container byte-identisch).
- `make versions` zeigt alle Pins (drei `FROM`-Digests + semgrep-Image inkl.
  Digest + runtime-image-ID).

**Review-Runde R1** (`docs/reviews/2026-06-20-slice-033-digest-pins.md`):
0 HIGH, 1 MEDIUM, 2 LOW, 1 INFO. MEDIUM-1 (Manifest-Listen-Eigenschaft
unbelegt) → **verifiziert** (alle vier Index-Digests, amd64+arm64). LOW-1
(`grep -oP`/PCRE → stiller leerer Beleg auf Nicht-GNU-grep) → behoben
(portables `sed`). LOW-2 (Tag/Digest-Drift bei Einzel-Hebung) → won't-fix:
durch die Doppel-Edit-Politik gedeckt, der Tag ist nur Lesehilfe. INFO-1
(`latest`-Drift) → Folge-Slice (s. u.). Immutability sauber.

**Gefundene Folge-Drift (delegiert).** Im ADR-Audit bestätigt:
[ADR-0002](../../adr/0002-distribution-ghcr-image.md) §4 „kein `latest`-Tag"
widerspricht `release.yml` (pusht `latest` für stabile Releases). Außerhalb
dieses Slices (Distributions-Tag-Politik, nicht Image-Digest) → eigener
Folge-Slice mit eigener ADR.

**Lerneintrag.** Ein Pin-**Beleg** (`make versions`) darf nicht still leer
werden (PCRE-Abhängigkeit) — portables `sed`, im Geist von „grün=Boden".
Inline-Digest + Versions-Var hält `make versions` ohne Build-Args lesbar.
