# Slice slice-034: `latest`-Tag-Politik versöhnen (ADR-0002 §4 ↔ release.yml)

**Status:** next (geplant — noch nicht in Arbeit).

**Welle:** welle-23-latest-tag (Trigger: ADR-Audit 2026-06-20 + Review R1 zu
slice-033, INFO-1 — `docs/reviews/2026-06-20-slice-033-digest-pins.md`).

**Bezug:**
[ADR-0002](../../adr/0002-distribution-ghcr-image.md) §4 (die zu versöhnende
Entscheidung „kein `latest`"),
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
(Docker-Image/Tagging),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(Konsumenten-Reproduzierbarkeit via Digest).

**Autor:** pt9912. **Datum:** 2026-06-20.

---

## 1. Problem

[ADR-0002](../../adr/0002-distribution-ghcr-image.md) §4 entscheidet
„ausschließlich volle Semver-Tags … **kein `latest`** und keine beweglichen
Major-/Minor-Tags". `release.yml` (slice-011) pusht aber für stabile Releases
zusätzlich `:latest` (unter `IS_STABLE`). Bewusste Implementierungs-Wahl
(Kommentar: „latest für stabile, Konsumenten auf Digest-Pins"), aber
[ADR-0002](../../adr/0002-distribution-ghcr-image.md) §4 wurde nie versöhnt —
Doku (Source-Precedence Rang 4) und Code widersprechen sich.

## 2. Zu entscheiden (im Slice)

Eine von zwei Richtungen, festzuhalten in einer **neuen ADR**
([ADR-0002](../../adr/0002-distribution-ghcr-image.md) ist Accepted/immutable):

- **A — Praxis ratifizieren (Empfehlung):** neue ADR „`latest` für stabile
  Releases", die [ADR-0002](../../adr/0002-distribution-ghcr-image.md) §4
  insoweit ablöst: `latest` zeigt auf das neueste **stabile** Release (kein
  Prerelease), Konsumenten pinnen verbindlich per `@sha256:`-Digest (bereits
  Praxis, `docs/user/releasing.md`). Begründung: Convenience-Tag, während der
  Digest-Pin die Reproduzierbarkeit trägt.
- **B — Code an §4 angleichen:** `latest`-Push aus `release.yml` entfernen,
  nur Semver-Tags. Strikter zu §4, verliert aber den Convenience-Tag.

## 3. Definition of Done (vorläufig)

- [ ] Richtung A/B entschieden; **neue ADR** geschrieben (mit
  `Supersedes`-Lineage auf [ADR-0002](../../adr/0002-distribution-ghcr-image.md),
  Teil-Supersede von §4), Index-Zeile + Status-Lineage.
- [ ] `release.yml` und `docs/user/releasing.md` konsistent zur Entscheidung;
  der [ADR-0002](../../adr/0002-distribution-ghcr-image.md)-§4-Widerspruch ist
  aufgelöst.
- [ ] `make gates` grün (Doku-Konsistenz); Review nach Bedarf.

## 4. Risiken / offene Punkte

- [ADR-0002](../../adr/0002-distribution-ghcr-image.md) ist immutable →
  Ablösung nur via neue ADR mit `Supersedes`-Klausel; das `matrix`-Modul
  kennt `allow-supersede-lineage` (slice-024) — prüfen, ob die Supersede-Kante
  eine Ausnahme braucht.
- Reine Distributions-/Doku-Arbeit; kein neuer `DC-*`-Vertrag; kein Carveout.

## 5. Trigger

ADR-Audit (2026-06-20): [ADR-0002](../../adr/0002-distribution-ghcr-image.md)
§4 ↔ `release.yml` unversöhnt; aus slice-033 bewusst ausgekoppelt
(Distributions-Tag-Politik, nicht Image-Digest).

## 6. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Distributions-/Doku-Arbeit; Greenfield-Default).
