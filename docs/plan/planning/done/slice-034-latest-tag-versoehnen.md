# Slice slice-034: `latest`-Tag-Politik versöhnen (ADR-0002 §4 ↔ release.yml)

**Status:** done (Closure 2026-06-21; **Richtung A** — Praxis ratifiziert via
[ADR-0014](../../adr/0014-latest-tag-fuer-stabile-releases.md)).

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

- [x] Richtung **A** (Praxis ratifizieren) entschieden; neue
  [ADR-0014](../../adr/0014-latest-tag-fuer-stabile-releases.md) geschrieben
  (`Supersedes` [ADR-0002](../../adr/0002-distribution-ghcr-image.md) §4,
  Teil-Supersede), ADR-Index-Zeile + Teil-Supersede-Notiz an
  [ADR-0002](../../adr/0002-distribution-ghcr-image.md).
- [x] `release.yml` und `docs/user/releasing.md` verweisen auf
  [ADR-0014](../../adr/0014-latest-tag-fuer-stabile-releases.md); der
  [ADR-0002](../../adr/0002-distribution-ghcr-image.md)-§4-Widerspruch ist
  aufgelöst (Code/Doku fuhren A bereits — nur die ADR-Ebene nachgezogen).
- [x] `make gates` grün (Doku-Konsistenz); unabhängiges Review R1
  (0 HIGH/0 MEDIUM/1 LOW/2 INFO);
  [ADR-0014](../../adr/0014-latest-tag-fuer-stabile-releases.md) → `Accepted`;
  Closure-Notiz + `git mv` nach `done/`.

## 4. Risiken / offene Punkte

- [ADR-0002](../../adr/0002-distribution-ghcr-image.md) ist immutable →
  Ablösung via neue
  [ADR-0014](../../adr/0014-latest-tag-fuer-stabile-releases.md) mit
  `Supersedes`-Klausel (Datei unverändert). **Geprüft:** weil nur §4
  *teil*-abgelöst wird, bleibt
  [ADR-0002](../../adr/0002-distribution-ghcr-image.md) `Accepted`/aktiv —
  der Verweis erzeugt **kein** `matrix-inactive`, ein
  `allow-supersede-lineage` ist **nicht** nötig.
- Reine Distributions-/Doku-Arbeit; kein neuer `DC-*`-Vertrag; kein Carveout.

## 5. Trigger

ADR-Audit (2026-06-20): [ADR-0002](../../adr/0002-distribution-ghcr-image.md)
§4 ↔ `release.yml` unversöhnt; aus slice-033 bewusst ausgekoppelt
(Distributions-Tag-Politik, nicht Image-Digest).

## 6. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Distributions-/Doku-Arbeit; Greenfield-Default).

## 7. Closure-Notiz (nach `done/`)

**Umsetzung.** Richtung **A** (Praxis ratifizieren): neue
[ADR-0014](../../adr/0014-latest-tag-fuer-stabile-releases.md) (`Accepted`)
löst die Tagging-Klausel „kein `latest`" aus
[ADR-0002](../../adr/0002-distribution-ghcr-image.md) §4 **teil**-ab —
`:latest` → neuestes **stabiles** Release, verbindlicher Konsum per
`@sha256:`-Digest
([ADR-0011](../../adr/0011-digest-pins-build-gate-images.md),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
[ADR-0002](../../adr/0002-distribution-ghcr-image.md) §1–3/5 bleiben gültig;
die Datei ist **byte-unverändert** (Immutabilität, AGENTS.md §3.5) — die
Versöhnung trägt allein [ADR-0014](../../adr/0014-latest-tag-fuer-stabile-releases.md)
plus die Teil-Supersede-Notiz in der ADR-Index-Zeile (gleiche Form wie die
bestehende Teil-Supersede-Präzedenz im Index). `release.yml` und
`docs/user/releasing.md` verweisen nun auf
[ADR-0014](../../adr/0014-latest-tag-fuer-stabile-releases.md) (sie fuhren
Richtung A schon bewusst — **kein Verhaltens-Delta**, kein
Release/`CHANGELOG`).

**Belege.**
- `make gates` grün (doc-check 90/0, lint, test, arch-check, coverage
  94,20 %, semgrep 55/0, gate-consistency).
- `matrix`: kein `matrix-inactive` —
  [ADR-0002](../../adr/0002-distribution-ghcr-image.md) bleibt
  `Accepted`/aktiv (Teil-Supersede ändert das Status-Feld nicht);
  `allow-supersede-lineage` bestätigt **nicht** nötig.

**Review-Runde R1** (`docs/reviews/2026-06-21-adr-0014-latest-tag.md`):
0 HIGH, 0 MEDIUM, 1 LOW, 2 INFO. INFO-1 (die
[`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Annotation
schrieb dem Vertrag eine
[ADR-0002](../../adr/0002-distribution-ghcr-image.md)-§Konsequenzen-Aussage
zu) → **vor Accept behoben**
(Bezug + Entscheidung §3 präzisiert). LOW-1 (Kopffeld `Supersedes:` nicht
vorbenutzt; Token-Kollision mit künftigem `allow-supersede-lineage`) →
**won't-fix:** folgt der Index-Konvention „Ablösung via `Supersedes ADR-NN`";
das Flag ist in d-checks `.d-check.yml` aus, die Kollision rein hypothetisch
und semantisch ohnehin korrekt. INFO-2 (`:latest` erstmals als gewollter
Distributions-Vertrag dokumentiert) → akzeptiert, genau der Zweck der ADR.

**Lerneintrag.** Ein Doku↔Code-Drift, bei dem der **Code** die bessere
Entscheidung schon getroffen hatte, wird über eine **ratifizierende** ADR
(neue Autorität) aufgelöst — nicht durch Rückbau des Codes; die immutable
Vorgänger-ADR bleibt unangetastet, der Teil-Supersede lebt im Index.
Independent Review fängt Faktentreue-Nuancen (INFO-1) vor der
Immutabilität.
