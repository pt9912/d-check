# ADR-0014 — `latest`-Tag für stabile Releases (ratifiziert ADR-0002 §4)

**Status:** Accepted
**Datum:** 2026-06-21
**Autor:** pt9912
**Bezug:**
[`DC-FA-DIST-001`](../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
(Image/Tagging),
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)
(Konsumenten-Reproduzierbarkeit via Digest),
[`DC-QA-04`](../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)
(Migrationsabdeckung — Pilot-/Konsumenten-Repos pinnen konkrete
Versionen/Digests),
[ADR-0011](0011-digest-pins-build-gate-images.md) (Digest-Pins tragen die
Reproduzierbarkeit).
**Supersedes:** [ADR-0002](0002-distribution-ghcr-image.md) §4 — **nur** die
Tagging-Klausel „kein `latest`"; §1–3 und §5 bleiben unverändert gültig.
**Schärft:** keine Spec-Stelle — Prozess-/Distributions-ADR; verbindlich für
die Tagging-Politik in
[`release.yml`](../../../.github/workflows/release.yml) und
[`releasing.md`](../../../docs/user/releasing.md).

## Kontext

[ADR-0002](0002-distribution-ghcr-image.md) §4 entschied „ausschließlich
volle Semver-Tags … **kein `latest`** und keine beweglichen
Major-/Minor-Tags". Die Release-Pipeline (slice-011) pusht für **stabile**
Releases jedoch zusätzlich `:latest` — eine bewusste Implementierungs-Wahl
(Komfort-Einstieg), mit verbindlichem Konsumenten-Verweis auf `@sha256:`-
Digest-Pins
([`releasing.md`](../../../docs/user/releasing.md)). Code und Betriebs-Doku
fahren also seit slice-011 **Richtung A**, während
[ADR-0002](0002-distribution-ghcr-image.md) §4 (Source-Precedence Rang 4)
nie nachgezogen wurde — ein Doku↔Code-Drift, bestätigt im ADR-Audit
(2026-06-20) und in Review R1/slice-033 (INFO-1).

Seit [ADR-0011](0011-digest-pins-build-gate-images.md) trägt der
`@sha256:`-Digest die Reproduzierbarkeit
([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)); ein
beweglicher `:latest`-Tag gefährdet sie nicht, solange Konsumenten
verbindlich per Digest pinnen. [ADR-0002](0002-distribution-ghcr-image.md)
ist `Accepted`/immutable (AGENTS.md §3.5) — die Versöhnung erfolgt per
**neuer** ADR, nicht per Edit.

## Entscheidung

1. **`:latest` für stabile Releases.** Für Tags ohne Prerelease-Suffix
   (`vX.Y.Z`) wird `:latest` gesetzt und gepusht; es zeigt stets auf das
   **neueste stabile** Release. Ratifiziert die bestehende
   [`release.yml`](../../../.github/workflows/release.yml)-Praxis
   (`IS_STABLE`-Verzweigung).
2. **Prereleases erhalten kein `:latest`** (bereits so).
3. **Verbindlicher Konsum per Digest.** Konsumenten pinnen auf
   `@sha256:`-Digest
   ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus),
   [ADR-0011](0011-digest-pins-build-gate-images.md)); `:latest` ist
   Komfort-Einstieg, **nicht** für CI-Pipelines. Migrations-/Konsumenten-
   Repos pinnen weiterhin konkrete Versionen/Digests
   ([`DC-QA-04`](../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Pilotpraxis;
   vgl. [ADR-0002](0002-distribution-ghcr-image.md) §Konsequenzen).
4. **Eng begrenzte Ablösung.** Diese ADR löst **ausschließlich** die
   Tagging-Klausel „kein `latest`" aus
   [ADR-0002](0002-distribution-ghcr-image.md) §4 ab; §1 (Multi-Stage),
   §2 (distroless), §3 (Entrypoint) und §5 (Binary-Distribution) bleiben
   unverändert. Weiterhin **keine** beweglichen Major-/Minor-Tags — nur
   `:latest` plus volle Semver-Tags.

## Verglichene Alternativen

| Alternative | Pro | Contra |
|---|---|---|
| **A — `:latest` für stabile ratifizieren (gewählt)** | Komfort-Tag; Code/Doku bereits konsistent; Digest trägt die Reproduzierbarkeit | ein beweglicher Tag existiert (für CI ungeeignet — dokumentiert) |
| B — `:latest` entfernen (strikt §4) | strikt reproduzierbar, kein beweglicher Tag | verliert den Komfort-Einstieg; Breaking für `:latest`-Nutzer; Code/Doku-Rückbau |
| Zusätzlich bewegliche Major-/Minor-Tags | feinere Komfort-Granularität | mehr bewegliche Tags = größere Drift-Fläche; nicht gefordert |

## Konsequenzen

- [ADR-0002](0002-distribution-ghcr-image.md) §4 Tagging-Klausel ist
  **teil-abgelöst**; die ADR-Index-Zeile von
  [ADR-0002](0002-distribution-ghcr-image.md) trägt die Notiz, §1–3 und §5
  bleiben unberührt. Die ADR-Datei selbst wird nicht editiert (immutable).
- [`release.yml`](../../../.github/workflows/release.yml) und
  [`releasing.md`](../../../docs/user/releasing.md) verweisen nun auf diese
  ADR (Traceability zur ratifizierenden Entscheidung statt zur abgelösten
  §4-Klausel).
- **Kein Verhaltens-Delta** (Ratifikation): kein Release/Version-Bump, kein
  `CHANGELOG`-Eintrag (nicht nutzersichtbar).
- `matrix`: [ADR-0002](0002-distribution-ghcr-image.md) bleibt `Accepted`
  (aktiv) — der Teil-Supersede ändert das Status-Feld der Datei nicht; der
  Verweis dieser ADR auf sie erzeugt **kein** `matrix-inactive`, ein
  `allow-supersede-lineage` ist hier nicht nötig.

## Fitness Function

- [`release.yml`](../../../.github/workflows/release.yml) setzt/pusht
  `:latest` nur bei `IS_STABLE=true` (vorhanden);
  [`releasing.md`](../../../docs/user/releasing.md) dokumentiert die
  Digest-Pin-Pflicht. `make doc-check` hält die ADR-/Doku-Verweise
  konsistent.
- Kein `make`-Gate prüft die CI-YAML selbst — die strukturelle Kontrolle
  ist die SemVer-Validate-Stufe plus die `IS_STABLE`-Verzweigung in
  [`release.yml`](../../../.github/workflows/release.yml).

## Re-Evaluierungs-Trigger

- Bedarf an beweglichen Major-/Minor-Tags → neue ADR.
- Wechsel der Registry- oder Tagging-Mechanik → Pin-Strategie prüfen.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-21 | Proposed — Audit 2026-06-20 (ADR-0002 §4 ↔ `release.yml` unversöhnt), slice-034 |
| 2026-06-21 | Proposed → Accepted (slice-034; unabhängiges Review R1: 0 HIGH/0 MEDIUM/1 LOW/2 INFO) |
