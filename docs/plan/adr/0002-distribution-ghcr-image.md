# ADR-0002 — Distribution: GHCR-Image auf distroless/static

**Status:** Accepted
**Datum:** 2026-06-10
**Autor:** pt9912
**Bezug:** [`DC-FA-DIST-001`](../../../spec/lastenheft.md#dc-fa-dist-001--docker-image),
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit),
[ADR-0001](0001-implementierungssprache.md)
**Schärft:** `spec/spezifikation.md` (entsteht mit slice-002:
Image-Aufruf-Konventionen) — nicht das Lastenheft.

## Kontext

[`DC-FA-DIST-001`](../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
fordert ein GHCR-Image mit Semver-Tags, Non-root-Ausführung,
read-only-Mount und identischem Verhalten zur nativen Ausführung.
Offen sind Basis-Image, Tagging-Disziplin, Entrypoint-Gestaltung und
der vertagte Binary-Distributionsweg. ADR-0001 (Go, `CGO_ENABLED=0`)
macht ein statisches Binary verfügbar.

## Entscheidung

1. **Multi-Stage-Build:** Build-Stage `FROM golang:<ver>@sha256:…`
   (Digest-gepinnt), Runtime-Stage
   `FROM gcr.io/distroless/static-debian12:nonroot@sha256:…`.
2. **Basis-Image distroless/static:** bringt CA-Zertifikate (nötig für
   TLS im Modul `external`) und einen Non-root-User mit — beides
   müsste bei `scratch` manuell nachgebaut werden; kein Shell/
   Paketmanager im Runtime-Image (minimale Angriffsfläche).
3. **Entrypoint:** `ENTRYPOINT ["/d-check", "/repo"]` — der
   Lastenheft-Default „Prüfung von `/repo`, Optionen als
   Container-Argumente angehängt" funktioniert damit wörtlich
   (`docker run … ghcr.io/pt9912/d-check:<tag> --disable anchors`).
4. **Tagging:** ausschließlich volle Semver-Tags (`vX.Y.Z`),
   **kein `latest`** und keine beweglichen Major-/Minor-Tags.
   Konsumenten-Repos pinnen zusätzlich per Digest
   (`@sha256:…`) — etablierte Praxis im Stack (b-cad pinnt d-migrate
   so).
5. **Binary-Distribution: vertagt mit Trigger.** Goreleaser-basierte
   Release-Binaries werden nachgerüstet, sobald ein konkreter
   Bedarf außerhalb Docker-fähiger Umgebungen auftritt (z. B.
   pre-commit-Hook auf Hosts ohne Docker). Bis dahin ist das Image der
   einzige Distributionsweg.

## Verglichene Alternativen

| Basis-Image | Pro | Contra |
|---|---|---|
| **distroless/static (gewählt)** | CA-Bundle + nonroot eingebaut, kein Shell/apk, ~2 MB Overhead | kein Debug-Shell im Image (bewusst) |
| scratch | minimal | keine CA-Zertifikate (Modul `external` bräuchte Handarbeit), kein nonroot-User-Eintrag |
| alpine | klein, Shell für Debugging | apk + Shell = Angriffsfläche; musl-Randfälle unnötig, da statisches Binary |
| debian-slim | vertraut | ~80 MB, Paketmanager im Runtime-Image |

## Konsequenzen

- Dockerfile (slice-003) ist Multi-Stage mit Digest-Pins; `make versions`
  (welle-04) belegt die Pins.
- Kein `latest`-Tag bedeutet: Migrations-Doku der Konsumenten-Repos
  nennt immer eine konkrete Version (`DC-QA-04`-Pilotläufe).
- Debugging im Container erfolgt über lokal gebaute Debug-Stages,
  nicht im Release-Image.

## Fitness Function

- CI-Gate (welle-04): Image-Build aus identischem Source-Tree ist
  digest-stabil; `harness/image-hash.txt` als Beleg-Manifest. <!-- d-check:ignore (geplantes Beleg-Manifest, entsteht mit fullbuild-Praxis) -->
- Negative-Kriterium aus `DC-FA-DIST-001` (kein Mount → Exit 2) läuft
  als Integrationstest gegen das gebaute Image.

## Re-Evaluierungs-Trigger

- Erster konkreter Bedarfsfall für Binary-Distribution (hebt
  Punkt 5 auf — neue ADR).
- distroless/static stellt CA-Bundle- oder nonroot-Garantien ein.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-10 | Proposed → Accepted (slice-001) |
