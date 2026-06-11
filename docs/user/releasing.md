# Releasing — d-check

Release-Prozess für `ghcr.io/pt9912/d-check`
([`DC-FA-DIST-001`](../../spec/lastenheft.md#dc-fa-dist-001--docker-image),
[ADR-0002](../plan/adr/0002-distribution-ghcr-image.md)). Diese Datei
beschreibt den Prozess; die Pipeline selbst ist
[`.github/workflows/release.yml`](../../.github/workflows/release.yml).

## Versionsquelle

Versionen folgen SemVer; die menschlich kuratierte Begründung jedes
Releases ist der zugehörige Abschnitt in
[`CHANGELOG.md`](../../CHANGELOG.md). Vor dem Tag wird dort der
`[Unreleased]`-Stand unter die neue Version geschnitten.

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
   **nur** für stabile Releases (kein Prerelease-Suffix).
5. **Digest-Pin** landet im Job-Summary und in den Notes des
   automatisch angelegten GitHub-Releases.

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
