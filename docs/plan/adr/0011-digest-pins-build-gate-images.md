# ADR-0011 — Digest-Pins aller Build- und Gate-Images

**Status:** Accepted
**Datum:** 2026-06-20
**Autor:** pt9912
**Bezug:** [ADR-0002](0002-distribution-ghcr-image.md) (Build-/Runtime-Image-Digests),
[ADR-0006](0006-lint-profil-solid.md) (lint-Stage-Image),
[ADR-0010](0010-semgrep-hermetisches-gate.md) (semgrep-Gate-Image),
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)
**Schärft:** keine Spec-Stelle — Prozess-/Reproduzierbarkeits-ADR;
verbindlich für jede `FROM`-Zeile im `Dockerfile`, das Scanner-Image in
`tools/semgrep.sh` und den Pin-Beleg in `make versions`.

## Kontext

[ADR-0002](0002-distribution-ghcr-image.md) §1 entscheidet bereits
Digest-Pins für die Build-Stage (`golang`) und die Runtime-Stage
(`distroless`) und nennt in den Konsequenzen „Dockerfile … mit
Digest-Pins". Der `Dockerfile` löste das nie ein — die `FROM`-Zeilen
waren nur Tag-gepinnt (Drift, festgestellt im slice-032-Review, INFO-1).
Zwei Images sind von [ADR-0002](0002-distribution-ghcr-image.md) zudem
gar nicht erfasst: das **lint-Stage-Image**
(`golangci/golangci-lint`, [ADR-0006](0006-lint-profil-solid.md)) und das
**semgrep-Gate-Image** ([ADR-0010](0010-semgrep-hermetisches-gate.md)
pinnt es nur per Tag). Ein veränderlicher Tag kann unter gleichem Namen
neu gepusht werden ⇒ gleicher Source-Tree, andere Toolchain, andere
Befunde — eine Determinismus-Restlücke
([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).

## Entscheidung

1. **Jedes extern bezogene Image wird per `@sha256:`-Digest gepinnt** —
   alle `FROM`-Zeilen des `Dockerfile` (`golang`,
   `golangci/golangci-lint`, `gcr.io/distroless/static-debian12`) und das
   Scanner-Image in `tools/semgrep.sh` (`semgrep/semgrep`).
2. **Tag bleibt lesbar daneben** (`image:tag@sha256:…`): der Digest ist
   die Wahrheit (Docker löst darüber auf), der Tag dient der Lesbarkeit im
   Diff. Vereinheitlicht [ADR-0002](0002-distribution-ghcr-image.md) §1
   auf den lint-Stage und hebt
   [ADR-0010](0010-semgrep-hermetisches-gate.md) vom Tag auf den Digest.
3. **Multi-Arch:** gepinnt wird der **Manifest-Listen-Digest** (nicht ein
   plattformspezifischer), damit der Build auf amd64/arm64 funktioniert.
4. **Update-Politik:** ein Digest-Hoist (Tag UND Digest gemeinsam) ist ein
   **bewusster** Commit mit Begründung im Body — analog zu
   `GO_VERSION`/`GOLANGCI_LINT_VERSION`, kein Auto-Update.
5. **Beleg:** `make versions` gibt die `FROM`-Digests (Dockerfile-grep)
   und das semgrep-Image inkl. Digest aus.

## Verglichene Alternativen

| Alternative | Pro | Contra |
|---|---|---|
| **Digest inline neben Tag, alle Images (gewählt)** | reproduzierbar; `make versions` zeigt Digests direkt; Tag lesbar | Hebung = bewusster Doppel-Edit (Version+Digest) |
| Nur Tag (Status quo) | minimaler Diff | Tag neu pushbar ⇒ nicht reproduzierbar (Restlücke [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)) |
| Digest-Build-Args zentral im Makefile | alle Pins an einem Ort | `make versions` zeigt nur Var-Namen; `FROM`-Zeile unleserlicher |
| Nur `golang`+`distroless` ([ADR-0002](0002-distribution-ghcr-image.md)-Minimum) | weniger Pflege | lint-/semgrep-Image bleiben Tag-gepinnt — inkonsistent |

## Konsequenzen

- Drift zu [ADR-0002](0002-distribution-ghcr-image.md) §1 ist behoben;
  lint- und semgrep-Image sind erstmals digest-gepinnt.
- Eine Image-/Version-Hebung berührt nun zwei Stellen (Tag + Digest) —
  bewusst, macht den Engine-Wechsel im Diff sichtbar.
- `make image-test`/`make gates` ziehen über den Digest dieselbe Engine
  wie CI — host-/zeitunabhängig.
- `latest`-Tag-Politik (Distribution) bleibt unberührt — das ist ein
  separater [ADR-0002](0002-distribution-ghcr-image.md)-§4-Punkt (eigene
  Folge-ADR).

## Fitness Function

`make gates`/`make image-test` bauen aus den digest-gepinnten Images;
`make versions` belegt die Pins (Dockerfile-`FROM`-Digests +
semgrep-Image).

## Re-Evaluierungs-Trigger

- Ein Image-Update (Version) → Digest mitheben (bewusster Commit).
- Ein Registry stellt Multi-Arch-Manifeste ein → Pin-Strategie prüfen.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-20 | Proposed → Accepted (slice-033) |
