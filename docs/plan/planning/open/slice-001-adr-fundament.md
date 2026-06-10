# Slice slice-001: ADR-Fundament — Sprache, Distribution, Config-Format

**Status:** open.

**Welle:** welle-01-fundament.

**Bezug:** [`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image),
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei),
[`DC-QA-01`](../../../../spec/lastenheft.md#dc-qa-01--performance),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).

**Autor:** pt9912. **Datum:** 2026-06-10.

---

## 1. Ziel

Die drei Fundament-Entscheidungen liegen als akzeptierte ADRs vor:
Implementierungssprache, Distributionsweg/Basis-Image und
Config-Format/-Parser für `.d-check.yml`.

## 2. Definition of Done

- [ ] ADR-0001 Implementierungssprache: `Accepted`, mindestens 3 verglichene Alternativen (Kandidaten u. a. Go, Python, Node — Bewertung gegen [`DC-QA-01`](../../../../spec/lastenheft.md#dc-qa-01--performance)–[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) und die Erfahrung der 12 Quell-Tools).
- [ ] ADR-0002 Distribution: `Accepted` — GHCR-Image, Basis-Image-Wahl (distroless/scratch/slim), Tagging-Schema; Binary-Distribution entschieden oder explizit vertagt mit Trigger.
- [ ] ADR-0003 Config-Format: `Accepted` — YAML-Schema-Ansatz für `.d-check.yml`, Validierungsstrategie ([`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei) Negative-Kriterium).
- [ ] ADR-Index aktualisiert.
- [ ] [`AGENTS.md`](../../../../AGENTS.md) §3.1/§3.2 um die sprachkonkreten Festlegungen ergänzt (Toolchain-Image, Suppression-Marker).
- [ ] `make gates` grün.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/0001-implementierungssprache.md` | neu | Kernentscheidung, blockiert Architektur und alle Implementierungs-Slices |
| `docs/plan/adr/0002-distribution-ghcr-image.md` | neu | [`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image) braucht Basis-Image- und Tagging-Festlegung |
| `docs/plan/adr/0003-config-format.md` | neu | [`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei) braucht Schema- und Validierungs-Festlegung |
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
- ADR-0002 und ADR-0003 hängen teilweise von ADR-0001 ab (verfügbare
  Parser/Basis-Images) — Reihenfolge innerhalb des Slice beachten.

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (reine Doku-/Entscheidungsarbeit; siehe
Kurs Modul 5 §Worked Mini-Example).
