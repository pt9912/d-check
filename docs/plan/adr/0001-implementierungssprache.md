# ADR-0001 — Implementierungssprache: Go

**Status:** Accepted
**Datum:** 2026-06-10
**Autor:** pt9912
**Bezug:** [`DC-QA-01`](../../../spec/lastenheft.md#dc-qa-01--performance),
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit),
[`DC-FA-DIST-001`](../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
**Schärft:** `spec/spezifikation.md` / `spec/architecture.md` —
entstehen mit slice-002 und setzen diese Entscheidung voraus. Das
Lastenheft bleibt unberührt (die Sprache ist dort bewusst nicht
festgelegt).

## Kontext

`d-check` ist ein CLI-Tool, das als kompaktes GHCR-Container-Image
verteilt wird und identisch nativ wie im Container laufen soll
([`DC-FA-DIST-001`](../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)). Die NFAs verlangen schnelle Läufe ohne
Runtime-Warmup ([`DC-QA-01`](../../../spec/lastenheft.md#dc-qa-01--performance)), keinerlei Laufzeit-Seiteneffekte und
strukturell kontrollierbare Netzwerknutzung ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)). Die zwölf
Quell-Tools existieren in Shell, Python und JavaScript — keine dieser
Basen ist gesetzt; das Lastenheft fordert Konsolidierung, nicht
Portierung.

## Entscheidung

Implementierungssprache ist **Go** (aktuelle stabile Version, im
Build-Image per Digest gepinnt — siehe
[ADR-0002](0002-distribution-ghcr-image.md)):

- `CGO_ENABLED=0` erzeugt ein vollständig **statisches Binary** —
  Basis-Image kann distroless/scratch-klein sein.
- **Cross-Compilation** über `GOOS`/`GOARCH` ohne Zusatz-Toolchain
  (linux/amd64, linux/arm64, darwin, windows) — hält die im Lastenheft
  vertagte Binary-Distribution trivial nachrüstbar.
- Die Standard-Library deckt Dateisystem, HTTP-Client, JSON und Regex
  ab; der Netz-Code bleibt sauber in einem Adapter isolierbar
  ([ADR-0004](0004-architektur-pattern-hexagonal.md), [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **Ökosystem-Konsistenz:** u-boot traf dieselbe Entscheidung
  (dortiges ADR-0001), k-deskflight ist die erprobte
  Go/Distroless-Referenz; Konventionen, Lint-Profile und
  depguard-artige Architektur-Checks existieren im eigenen Stack.

## Verglichene Alternativen

| Alternative | Pro | Contra |
|---|---|---|
| **Go (gewählt)** | statisches Binary, Cross-Compile eingebaut, schneller Start, Stack-Precedent | kein YAML in der Stdlib → gepinnte Dependency nötig ([ADR-0003](0003-config-format.md)) |
| Python | stdlib-reich, zwei Quell-Tool-Familien als Vorlage (`check_refs.py`) | Image ≥ ~50 MB (Interpreter), Startzeit, kein statisches Binary |
| Node/TypeScript | `docs-check.js` (461 Z.) als feature-reiche Vorlage | Runtime + `node_modules` im Image, Startzeit, Supply-Chain-Fläche |
| Rust | statisch und sehr klein, starke Korrektheit | höchster Implementierungsaufwand, kaum Vorlagen im eigenen Stack für CLI-Textanalyse |

Shell scheidet als vierte Option aus: Anker-, ID- und Matrix-Logik
sind in awk/grep nicht wartbar (Erfahrung der `verify-doc-refs.sh`-Familie).

## Konsequenzen

- YAML-Parsing braucht eine gepinnte Dependency
  ([ADR-0003](0003-config-format.md)); ansonsten Stdlib-only anstreben.
- Lint-Profil: `golangci-lint`; der Suppression-Marker `//nolint`
  fällt unter das Suppression-Verbot (`AGENTS.md` §3.2).
- `make arch-check` ist mit depguard-/golangci-lint-Regeln umsetzbar
  (Fitness Function aus [ADR-0004](0004-architektur-pattern-hexagonal.md)).
- Die Go-Toolchain läuft ausschließlich in Docker (`AGENTS.md` §3.1);
  der PreToolUse-Guard blockiert Host-`go`.

## Fitness Function

- Build-Gate (ab slice-003) prüft, dass das erzeugte Binary statisch
  gelinkt ist (`CGO_ENABLED=0`-Nachweis).
- `make versions` (welle-04) belegt die gepinnte Go-Version gegen das
  committete Toolchain-Manifest.

## Re-Evaluierungs-Trigger

- [`DC-QA-01`](../../../spec/lastenheft.md#dc-qa-01--performance)
  (1.000 Dateien < 5 s) wird trotz Profiling-Iteration verfehlt.
- Die YAML-Dependency wird unwartbar/unsicher und hat keinen
  gleichwertigen Ersatz.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-10 | Proposed → Accepted (slice-001) |
