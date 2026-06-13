# Slice slice-001: ADR-Fundament — Sprache, Distribution, Config-Format, Architektur-Pattern

**Status:** done.

**Welle:** welle-01-fundament.

**Bezug:** [`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image),
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei),
[`DC-QA-01`](../../../../spec/lastenheft.md#dc-qa-01--performance),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).

**Autor:** pt9912. **Datum:** 2026-06-10.

---

## 1. Ziel

Die vier Fundament-Entscheidungen liegen als akzeptierte ADRs vor:
Implementierungssprache, Distributionsweg/Basis-Image,
Config-Format/-Parser für `.d-check.yml` und Architektur-Pattern.

## 2. Definition of Done

- [x] [ADR-0001](../../adr/0001-implementierungssprache.md) Implementierungssprache: `Accepted`, mindestens 3 verglichene Alternativen (Kandidaten u. a. Go, Python, Node — Bewertung gegen [`DC-QA-01`](../../../../spec/lastenheft.md#dc-qa-01--performance)–[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) und die Erfahrung der 12 Quell-Tools).
- [x] [ADR-0002](../../adr/0002-distribution-ghcr-image.md) Distribution: `Accepted` — GHCR-Image, Basis-Image-Wahl (distroless/scratch/slim), Tagging-Schema; Binary-Distribution entschieden oder explizit vertagt mit Trigger.
- [x] [ADR-0003](../../adr/0003-config-format.md) Config-Format: `Accepted` — YAML-Schema-Ansatz für `.d-check.yml`, Validierungsstrategie ([`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei) Negative-Kriterium).
- [x] [ADR-0004](../../adr/0004-architektur-pattern-hexagonal.md) Architektur-Pattern: `Accepted` — hexagonaler Schnitt (Kern ohne I/O; Ports: Filesystem, HTTP, Config, Reporter; einziger driving Adapter: das CLI), inkl. Fitness-Function-Skizze für das spätere `make arch-check`-Gate (strukturelle Durchsetzung von [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit): Netz nur im HTTP-Adapter, Kern ohne I/O-Imports).
- [x] [ADR-Index](../../adr/README.md) aktualisiert.
- [x] [`AGENTS.md`](../../../../AGENTS.md) §3.1/§3.2 um die sprachkonkreten Festlegungen ergänzt (Toolchain-Image, Suppression-Marker).
- [x] `make gates` grün.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`docs/plan/adr/0001-implementierungssprache.md`](../../adr/0001-implementierungssprache.md) | neu | Kernentscheidung, blockiert Architektur und alle Implementierungs-Slices |
| [`docs/plan/adr/0002-distribution-ghcr-image.md`](../../adr/0002-distribution-ghcr-image.md) | neu | [`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image) braucht Basis-Image- und Tagging-Festlegung |
| [`docs/plan/adr/0003-config-format.md`](../../adr/0003-config-format.md) | neu | [`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei) braucht Schema- und Validierungs-Festlegung |
| [`docs/plan/adr/0004-architektur-pattern-hexagonal.md`](../../adr/0004-architektur-pattern-hexagonal.md) | neu | sprachunabhängig (parallel zu [`ADR-0001`](../../adr/0001-implementierungssprache.md) entscheidbar); Zielbild für slice-002-Architektur und `arch-check`-Fitness-Function |
| [`docs/plan/adr/README.md`](../../adr/README.md) | update | Index-Pflicht |
| [`AGENTS.md`](../../../../AGENTS.md) | update | Hard Rules §3.1/§3.2 sprachkonkret machen |

## 4. Trigger

Sofort — welle-01-fundament ist die aktive Welle.

## 5. Closure-Trigger

DoD vollständig + Commit(s) auf `main` + Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Die Sprachentscheidung prägt Image-Größe und Bootstrap-Aufwand;
  eine spätere Revision wäre teuer (ganzer Code-Bestand). Daher
  Alternativen-Vergleich explizit gegen die NFAs führen.
- [`ADR-0002`](../../adr/0002-distribution-ghcr-image.md) und [`ADR-0003`](../../adr/0003-config-format.md) hängen teilweise von [`ADR-0001`](../../adr/0001-implementierungssprache.md) ab (verfügbare
  Parser/Basis-Images) — Reihenfolge innerhalb des Slice beachten;
  [`ADR-0004`](../../adr/0004-architektur-pattern-hexagonal.md) ist sprachunabhängig und kann zuerst entschieden werden.
- Hexagon-Risiko bei einem kleinen Tool ist Über-Zeremonie — [`ADR-0004`](../../adr/0004-architektur-pattern-hexagonal.md)
  soll explizit „Hexagon light" festschreiben (vier schmale Ports,
  keine Application-/Domain-Unterteilung).

## 7. Closure-Notiz (nach `done/`)

**Umsetzung:** Commit `78081bd` (ADRs 0001–0004, Index, `AGENTS.md`
§3.1/§3.2, Guard-Erweiterung um Host-`go`).

- **Was hat funktioniert:** Alle vier Entscheidungen ließen sich direkt
  gegen die `DC-QA`-NFAs führen; [`ADR-0004`](../../adr/0004-architektur-pattern-hexagonal.md) zuerst zu entscheiden
  (sprachunabhängig) hat die übrigen Abwägungen vereinfacht. Die
  User-Setzung „Go wegen kompakter Multi-Plattform-Binaries" deckte
  sich mit der NFA-Abwägung.
- **Anders als geplant:** Zusätzlich zum geplanten `AGENTS.md`-Update
  wurde der PreToolUse-Guard um Host-`go`/`golangci-lint` erweitert
  (folgt aus §3.1, stand nicht in der Plan-Tabelle). [`ADR-0002`](../../adr/0002-distribution-ghcr-image.md) legt
  auch den Image-Entrypoint fest (knüpft an die R1-Präzisierung
  „Image-Default-Befehl" im Lastenheft an).
- **Steering-Loop-Lerneintrag:** Das ADR-Feld `Schärft:` zeigt vor
  slice-002 zwangsläufig auf noch nicht existierende Spec-Straten —
  hier als „entstehen mit slice-002" formuliert. Bei slice-002-Closure
  prüfen, ob die `Schärft:`-Felder auf konkrete Abschnitte
  präzisiert werden sollen.
- **Folge-Slices:** keine neuen; slice-002 ist durch [`ADR-0001`](../../adr/0001-implementierungssprache.md)…0004
  entsperrt (Trigger erfüllt).

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (reine Doku-/Entscheidungsarbeit; siehe
Kurs Modul 5 §Worked Mini-Example).
