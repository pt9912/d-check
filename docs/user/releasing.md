# Releasing — d-check

Release-Prozess für `ghcr.io/pt9912/d-check`
([`DC-FA-DIST-001`](../../spec/lastenheft.md#dc-fa-dist-001--docker-image),
[ADR-0002](../plan/adr/0002-distribution-ghcr-image.md),
[ADR-0014](../plan/adr/0014-latest-tag-fuer-stabile-releases.md)). Diese
Datei beschreibt den Prozess; die Pipeline selbst ist
[`.github/workflows/release.yml`](../../.github/workflows/release.yml).

## Versionsquelle

Versionen folgen SemVer; die menschlich kuratierte Begründung jedes
Releases ist der zugehörige Abschnitt in
[`CHANGELOG.md`](../../CHANGELOG.md). Vor dem Tag wird dort der
`[Unreleased]`-Stand unter die neue Version geschnitten.

Die **aktuelle** Version führt zusätzlich das Release-Register
[`version.md`](../../version.md#aktuell) (§Aktuell). Das opt-in Modul `versions`
([`DC-FA-VER-001`](../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in))
prüft im Dogfooding-Lauf, dass alle gepinnten `ghcr`-Image-Verweise genau diese
Version tragen — der Bump beim Release ist daher nicht optional
(siehe [Release-Prep](#release-prep-vor-dem-tag)).

## Release-Prep (vor dem Tag)

In **einem** Commit vor dem Tag (kein Slice-Commit), sonst läuft `make ci` rot:

1. **`version.md`** — §Aktuell auf die neue Version, neue §Verlauf-Zeile, und den
   `<a id>`-Anker **auf die neue Version verschieben** (die bisherige Zeile
   verliert ihn; nur die aktuelle Version ist Anker-Ziel — veraltete
   Markdown-Link-Pins brechen so via `anchors`-Gate).
2. **Alle gepinnten `ghcr`-Image-Verweise** (README, Benutzerhandbuch) auf die
   neue Version ziehen — das aktive `versions`-Gate meldet sonst `version-stale`
   für jeden vergessenen Pin. Historische Pins in `done/`-Slices, `CHANGELOG.md`
   und der Lastenheft-Historie sind per `exempt-paths` ausgenommen.
3. **`CHANGELOG.md`** — den `[Unreleased]`-Stand unter die neue Version schneiden.
4. **Prosa-Currency von Hand nachziehen — kein Gate erzwingt sie.** Der
   `versions`-Gate prüft nur `ghcr`-**präfixierte** Pins gegen `version.md#aktuell`,
   nicht Fließtext oder nackte Tags. Betroffen:
   - **Benutzerhandbuch** — Header-Stempel (Handbuch-/Software-Version),
     Versionsverlauf-Zeile (§11; die neue Zeile **chronologisch unter die
     letzte**, nicht oben anfügen), ggf. neue Feature-Abschnitte (§5/§6) und das
     **bare-Tag-Beispiel** in §Versionen und Tags (`:vX.Y.Z` **ohne** `ghcr`-Präfix
     ⇒ vom `versions`-Gate nicht erfasst, driftet still). **Braucht das Feature
     eine §4-Aufgabe, schreibe eine eigene** — nach dem
     [Benutzerhandbuch-Standard](benutzerhandbuch-standard.md) §5
     (Ausgangslage/Ziel/Vorgehen/Ergebnis) — und **hänge sie nicht an eine
     bestehende §4-Aufgabe an**. Kein Gate erzwingt Aufgabenorientierung; genau
     durch Anhängen wuchs §4.12 auf ~330 Zeilen / 8 Themen, bis es aufgetrennt
     werden musste.
   - **README — beide Sprachfassungen synchron halten:** `README.de.md` (Deutsch,
     **kanonische Quelle — zuerst ändern**) und danach `README.md` (Englisch,
     nachübersetzen). Bei einem neuen Modul in **jeder** Fassung: (a) die
     **Status-Zeile** („alle N Regelmodule (…)" — Zahl *und* Enumeration *und* das
     „zuletzt das Modul X"-Fragment) und (b) die **Modul-Liste** unter
     §Was ist d-check ergänzen; das Intro-Framing prüfen, falls ein
     Nicht-Referenz-Modul (Content-Drift/Immutabilität/Versions/Traceability/Planning)
     neu ist. d-check prüft `links`/`anchors`/`ids`/`versions` in beiden READMEs, aber
     **nicht** die inhaltliche DE↔EN-Synchronität der Prosa.
   - **Das Datum kommt aus dem Kalender, nicht aus der Zeile darüber.** Ein
     Release trägt fünf Datumsstempel (`version.md` §Aktuell **und** §Verlauf,
     `CHANGELOG`-Überschrift, Handbuch-Kopf **und** §11-Zeile, die Historien von
     Lastenheft und Spezifikation) — und wer sie schreibt, schreibt sie von der
     Vorgängerzeile ab. Zieht sich eine Welle über mehr als einen Tag, datiert
     das Release damit auf den Tag ihres Beginns. Kein Gate sieht das;
     `git log -1 --date=short --format=%ad` ist das Orakel.
   - **Operations-Referenz** ([`operations.md`](operations.md)) — bei einem
     neuen Modul die Modul-Enumeration der `--enable`/`--disable`-Zeile, bei
     einer neuen CLI-Option die Optionen-Tabelle ergänzen. Kein Gate prüft
     diese Enumerationen — die Modul-Liste blieb so von v0.25 bis v0.37 still
     bei acht Modulen stehen.
5. **`make ci`** lokal grün fahren (Pre-Tag-De-Risk, „grün = Boden"), erst dann
   taggen.

Der Digest-Pin in Handbuch §2 entsteht **nach** dem Tag (er existiert erst nach
dem GHCR-Push) als Folge-Commit.

## Release auslösen

```sh
git tag v0.1.0
git push origin v0.1.0
```

Die Pipeline (`release.yml`) läuft bei jedem `v*`-Tag-Push:

1. **SemVer-Validate** (fail-fast): nur `vMAJOR.MINOR.PATCH` oder
   `…-PRERELEASE`; Build-Metadaten (`+`) werden abgelehnt.
2. **`make ci`** — alle Gates plus Image-Integrationstests; baut das
   Runtime-Image mit `VERSION` aus dem Tag.
3. **OCI-Label-Pin** — `org.opencontainers.image.version` muss exakt
   der Tag-Version entsprechen (Version-Drift shippt nicht).
4. **Push** nach `ghcr.io/pt9912/d-check:v<version>`; `:latest`
   **nur** für stabile Releases (kein Prerelease-Suffix) —
   [ADR-0014](../plan/adr/0014-latest-tag-fuer-stabile-releases.md).
5. **Digest-Pin** landet im Job-Summary und in den Notes des
   automatisch angelegten GitHub-Releases. Existiert das Release zum
   Tag bereits (z. B. Workflow-Re-Run), wird es wiederverwendet — der
   aktuelle Digest steht dann nur im Job-Summary.

## Konsum (Digest-Pin)

Konsumenten pinnen auf den Digest, nicht auf bewegliche Tags
(u-boot-Pin-Politik):

```sh
docker run --rm -v "$PWD:/repo:ro" \
  ghcr.io/pt9912/d-check@sha256:<digest-aus-den-release-notes>
```

`:latest` existiert als Komfort-Einstieg, ist aber für CI-Pipelines
ungeeignet (beweglich). Tag-Hebungen in konsumierenden Repos sind
Routine-Pins: Digest austauschen, Begründung in den Commit-Body.

## Aufruf-Referenz

Siehe [`operations.md`](operations.md).
