# Slice slice-011: GHCR-Release-Pipeline

**Status:** open.

**Welle:** welle-04-distribution-und-migration.

**Bezug:** [`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
(Veröffentlichungs-Teil: `ghcr.io/pt9912/d-check`, Semver-Tags),
[ADR-0002](../../adr/0002-distribution-ghcr-image.md);
[`MR-009`](../../../../harness/conventions.md#mr-009--source-precedence-ohne-docsuser-rang)
(Auflösungs-Trigger: Release-/Betriebs-Doku entsteht hier);
Meilenstein M3 (Teil 1: Image veröffentlicht).

**Autor:** pt9912. **Datum:** 2026-06-11.

---

## 1. Ziel

`d-check` wird automatisiert auf GHCR veröffentlicht: Tag-Push
(`v<semver>`) → Pipeline baut, prüft (`make ci`), pusht das Image mit
Semver-Tag und dokumentiertem Digest-Pin. Der erste Release existiert.

## 2. Definition of Done

- [ ] `.github/workflows/release.yml`: läuft bei Tag-Push `v*`;
  Schritte: Checkout → `make ci` (Gates aus slice-010) → Build →
  Push nach `ghcr.io/pt9912/d-check` mit Semver-Tag (zusätzlich
  `latest` nur für stabile Releases); der Image-Digest wird im
  Job-Log und im GitHub-Release festgehalten (Digest-Pin für
  Konsumenten).
- [ ] Release-/Betriebs-Doku unter `docs/user/` (Releasing: wie wird
  getaggt, wie konsumieren Repos das Image per Digest-Pin;
  Operations: Aufruf-Referenz) — damit ist der Auflösungs-Trigger von
  [`MR-009`](../../../../harness/conventions.md#mr-009--source-precedence-ohne-docsuser-rang)
  erfüllt: der `docs/user`-Rang wird in die Source-Precedence-Tabellen
  ([`harness/README.md`](../../../../harness/README.md),
  [`AGENTS.md`](../../../../AGENTS.md) §2) eingefügt, der MR-Eintrag
  als aufgelöst markiert.
- [ ] Erster Release `v0.1.0` ist tatsächlich veröffentlicht: Image
  auf GHCR abrufbar, Digest in der Closure-Notiz dokumentiert —
  Meilenstein M3, Teil 1.
- [ ] CI-Badge bzw. Lauf-Status-Verweis in
  [`harness/README.md`](../../../../harness/README.md) §Sensors
  („Aktueller Lauf-Status") auf die Pipeline erweitert.
- [ ] `make gates` grün; [`CHANGELOG.md`](../../../../CHANGELOG.md)
  (Release-Sektion `0.1.0`); Closure-Notiz.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.github/workflows/release.yml` | neu | Tag-getriggerter Build + Push (GITHUB_TOKEN, packages:write) |
| `docs/user/releasing.md`, `docs/user/operations.md` | neu | Release-Prozess + Aufruf-Referenz (löst [`MR-009`](../../../../harness/conventions.md#mr-009--source-precedence-ohne-docsuser-rang)) |
| [`harness/README.md`](../../../../harness/README.md), [`AGENTS.md`](../../../../AGENTS.md), [`harness/conventions.md`](../../../../harness/conventions.md) | update | Source-Precedence-Rang `docs/user`, `MR-009`-Auflösung, Lauf-Status |
| [`CHANGELOG.md`](../../../../CHANGELOG.md) | update | `0.1.0`-Release-Schnitt |

## 4. Trigger

slice-010 done — die Pipeline fährt `make ci` und die
Image-Integrationstests aus slice-010.

## 5. Closure-Trigger

DoD vollständig **inklusive real veröffentlichtem Release** (extern
verifizierbar: `docker pull ghcr.io/pt9912/d-check:v0.1.0`) +
Closure-Notiz.

## 6. Risiken und offene Punkte

- Externe Abhängigkeit: GitHub-Actions-Umgebung, `packages:write`-
  Permission, GHCR-Sichtbarkeit (public/private) — nicht lokal per
  `make gates` beweisbar; der Slice bleibt offen, bis der echte
  Release-Lauf grün war (keine Closure auf Verdacht — Harness-Lüge).
- Versionswahl `v0.1.0` vor Lastenheft-Status „Draft": bewusst 0.x —
  der Statuswechsel des Lastenhefts ist davon unabhängig
  (User-Entscheidung).
- `latest`-Tag-Politik konservativ halten (nur stabile Releases),
  Konsumenten werden auf Digest-Pins verwiesen (u-boot-Pin-Politik).

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Pipeline-/Doku-Arbeit; siehe Kurs
Modul 5 §Worked Mini-Example).
