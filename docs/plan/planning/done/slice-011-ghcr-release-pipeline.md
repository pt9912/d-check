# Slice slice-011: GHCR-Release-Pipeline

**Status:** done.

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

- [x] `.github/workflows/release.yml`: läuft bei Tag-Push `v*`;
  Schritte: Checkout → `make ci` (Gates aus slice-010) → Build →
  Push nach `ghcr.io/pt9912/d-check` mit Semver-Tag (zusätzlich
  `latest` nur für stabile Releases); der Image-Digest wird im
  Job-Log und im GitHub-Release festgehalten (Digest-Pin für
  Konsumenten).
- [x] Release-/Betriebs-Doku unter `docs/user/` (Releasing: wie wird
  getaggt, wie konsumieren Repos das Image per Digest-Pin;
  Operations: Aufruf-Referenz) — damit ist der Auflösungs-Trigger von
  [`MR-009`](../../../../harness/conventions.md#mr-009--source-precedence-ohne-docsuser-rang)
  erfüllt: der `docs/user`-Rang wird in die Source-Precedence-Tabellen
  ([`harness/README.md`](../../../../harness/README.md),
  [`AGENTS.md`](../../../../AGENTS.md) §2) eingefügt, der MR-Eintrag
  als aufgelöst markiert.
- [x] Erster Release `v0.1.0` ist tatsächlich veröffentlicht: Image
  auf GHCR abrufbar, Digest in der Closure-Notiz dokumentiert —
  Meilenstein M3, Teil 1.
- [x] CI-Badge bzw. Lauf-Status-Verweis in
  [`harness/README.md`](../../../../harness/README.md) §Sensors
  („Aktueller Lauf-Status") auf die Pipeline erweitert.
- [x] `make gates` grün; [`CHANGELOG.md`](../../../../CHANGELOG.md)
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

**Umsetzung:** Commit `4180a83` (Workflow, docs/user,
[`MR-010`](../../../../harness/conventions.md#mr-010--auflösung-von-mr-009-docsuser-rang-eingefügt),
0.1.0-Schnitt). **Release-Beweis (extern verifiziert):** Run
`27364508002` grün in 1m52s (alle Steps inkl. `make ci` und
OCI-Label-Pin); `docker pull ghcr.io/pt9912/d-check:v0.1.0`
erfolgreich; Smoke-Lauf per Digest-Pin gegen das eigene Repo: 35
Dateien, 0 Befunde; Versions-Label `0.1.0` matcht den Tag.
**Digest-Pin:**
`ghcr.io/pt9912/d-check@sha256:5710b54bc4712af9769d7a820fd3fe62621451daeb43f3e9737b382099137b9e`
— Meilenstein M3, Teil 1 (Image veröffentlicht).

- **Was hat funktioniert:** Das u-boot-publish-Muster (SemVer-Validate
  fail-fast, SHA-gepinnte Actions, latest-nur-stabil, Digest in die
  Release-Notes) ließ sich nahezu unverändert übernehmen — der erste
  echte Pipeline-Lauf war sofort grün, weil `make ci` lokal und in
  Actions identisch ist (Docker-only zahlt sich exakt hier aus).
- **Anders als geplant:** Die Tag↔Image-Versionskette
  (`ARG VERSION` → OCI-Label `image.version`, vom Workflow gegen den
  Tag gepinnt; plus `licenses=MIT`) war im Plan §3 nicht vorgesehen —
  übernommen aus dem u-boot-Vorbild, weil Digest-Pins nur
  vertrauenswürdig sind, wenn das Artefakt seine Version beweisbar
  trägt.
- **Steering-Loop-Lerneintrag:** Ein extern abhängiger Slice wird
  beherrschbar, wenn der lokale Sensor (`make ci`) exakt das ist, was
  die Pipeline fährt — die einzige echte Unbekannte war dann nur noch
  die GitHub-Permission. Die Closure-Disziplin „offen bis zum realen
  Pull-Beweis" hat sich als billig erwiesen (eine Stunde Latenz) und
  verhindert die Harness-Lüge „Release behauptet, nie gezogen".
- **Folge-Slices:** keine neuen; slice-012 (Pilot-Migrationen) ist
  entriegelt und konsumiert den obigen Digest-Pin.

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Pipeline-/Doku-Arbeit; siehe Kurs
Modul 5 §Worked Mini-Example).
