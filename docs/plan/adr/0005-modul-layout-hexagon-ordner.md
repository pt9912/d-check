# ADR-0005 — Modul-Layout: hexagon-/adapter-Ordnerkonvention nach u-boot

**Status:** Accepted
**Datum:** 2026-06-10
**Autor:** pt9912
**Bezug:** [ADR-0004](0004-architektur-pattern-hexagonal.md) (ersetzt
dessen Pfad-Festlegung — Teil-Supersede; das Pattern „Hexagon light"
und die Rollen-Regeln bleiben unverändert),
[ADR-0001](0001-implementierungssprache.md)
**Schärft:** `spec/architecture.md` §2 (liefert die sprachkonkrete
Übersetzung der dortigen Rollen-Constraints) — nicht das Lastenheft.

## Kontext

ADR-0004 legte neben dem Pattern auch konkrete Modul-Pfade fest
(`internal/core`, `internal/adapter/{fs,httpcheck,config,report}`).
Das Ökosystem (u-boot als Referenz) verwendet eine andere, etablierte
Ordnerkonvention: `internal/hexagon/…` für Kern und Ports,
`internal/adapter/{driven,driving}/…` für Adapter. Einheitliche
Layouts über die Schwester-Repos senken Einarbeitungs- und
Review-Kosten (Entscheidung User, 2026-06-10). ADRs sind nach
`Accepted` immutable — die Pfad-Revision ist daher eine neue ADR.

## Entscheidung

Modul-Layout nach u-boot-Konvention, „light" gefüllt (kein
Domain-/Application-Split — ADR-0004 bleibt dafür maßgeblich):

| Rolle (spec/architecture.md §2) | Pfad |
|---|---|
| Kern | `internal/hexagon/core` |
| Ports (driven) | `internal/hexagon/port/driven` (Paket `driven`) |
| Filesystem-Adapter | `internal/adapter/driven/fs` |
| HTTP-Adapter (ab Modul `external`) | `internal/adapter/driven/httpcheck` |
| Config-Adapter | `internal/adapter/driven/configyaml` |
| Reporter-Adapter | `internal/adapter/driven/report` |
| CLI (driving, Composition Root) | `internal/adapter/driving/cli` |
| Einstiegspunkt (dünn) | `cmd/d-check` |

Import-Regeln (sprachkonkrete Übersetzung der Rollen-Constraints):

1. `internal/hexagon/*` importiert weder I/O-APIs (`os`, `syscall`,
   `io/fs`, `net`-Sockets, `net/http`) noch `internal/adapter/*` noch
   `gopkg.in/yaml.v3`. Reine Parser ohne I/O (`net/url`) sind erlaubt.
2. `net/http` ausschließlich in `internal/adapter/driven/httpcheck`.
3. `gopkg.in/yaml.v3` ausschließlich in
   `internal/adapter/driven/configyaml`.
4. `os` ausschließlich in `internal/adapter/driven/fs`,
   `internal/adapter/driving/cli` und `cmd/*` (Composition Root).
5. driven Adapter importieren einander nicht; das driving CLI darf
   alle Adapter verdrahten.

## Verglichene Alternativen

| Alternative | Pro | Contra |
|---|---|---|
| **u-boot-Konvention, light gefüllt (gewählt)** | Ökosystem-einheitlich; Ports als eigenes Paket testbar; driving/driven explizit | minimal mehr Verzeichnistiefe |
| Layout aus ADR-0004 (`internal/core` + flache Adapter) | kürzeste Pfade | weicht vom Schwester-Repo-Standard ab; driving/driven nicht sichtbar |
| u-boot-Konvention voll (domain/application-Split) | maximale Strukturgleichheit mit u-boot | Über-Zeremonie für ein Ein-Use-Case-Tool (Begründung in ADR-0004) |
| Flaches `pkg/`-Layout | simpel | keine maschinell prüfbare Hexagon-Grenze |

## Konsequenzen

- `make arch-check` (Fitness Function) prüft die Regeln 1–5 dieser
  ADR; `spec/architecture.md` bleibt sprachneutral und unverändert.
- Slices referenzieren für Modul-Pfade diese ADR statt ADR-0004.
- Das Modul `external` erhält bei Implementierung den Pfad
  `internal/adapter/driven/httpcheck`.

## Fitness Function

`tools/arch-check.sh` (Dockerfile-Stage `arch-check`, `make
arch-check`): übersetzt die Regeln 1–5 in `go list`-Prüfungen;
Verstoß bricht den Build. Bindung:
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).

## Re-Evaluierungs-Trigger

- u-boot ändert seine Ordnerkonvention (Gleichlauf neu bewerten).
- d-check erhält einen zweiten driving Adapter oder echte
  Domain-Komplexität (dann ADR-0004-Trigger: light-Ausprägung neu
  bewerten).

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-10 | Proposed → Accepted (slice-003) |
