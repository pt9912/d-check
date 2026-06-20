# Slice slice-033: Image-Pins auf Digest (Dockerfile + semgrep)

**Status:** next (geplant — noch nicht in Arbeit).

**Welle:** welle-22-digest-pins (Trigger: Review R1 zu slice-032,
INFO-1 — `docs/reviews/2026-06-20-slice-032-semgrep-gate.md`; plus
[ADR-0002](../../adr/0002-distribution-ghcr-image.md)↔Dockerfile-Drift).

**Bezug:**
[ADR-0002](../../adr/0002-distribution-ghcr-image.md) (Distribution/Build —
beschreibt den Build-Stage bereits als `FROM golang:<ver>@sha256:…`),
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

- [ ] `Dockerfile`: jede `FROM`-Zeile mit `@sha256:<digest>` neben dem Tag
  (Tag bleibt als lesbarer Hinweis); Digests als Build-Args/Variablen, damit
  `make versions` sie ausgibt.
- [ ] `tools/semgrep.sh`: `semgrep/semgrep:1.167.0@sha256:<digest>`.
- [ ] Die betroffenen ADRs ([ADR-0002](../../adr/0002-distribution-ghcr-image.md),
  ggf. [ADR-0010](../../adr/0010-semgrep-hermetisches-gate.md)) auf den realen
  Digest-Pin-Stand nachgezogen; Drift aufgelöst.
- [ ] `make versions` zeigt Digests; Update-Politik (bewusste Pin-Hebung,
  eigener Commit) dokumentiert — analog `GO_VERSION`/`GOLANGCI_LINT_VERSION`.
- [ ] `make gates` grün (inkl. `make image-test`, `make semgrep` offline);
  reproduzierbar.
- [ ] Unabhängiges Review R1; Closure-Notiz; ADR-Status.

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
