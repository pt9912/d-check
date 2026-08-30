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
     letzte**, nicht oben anfügen), ggf. neue Feature-Abschnitte (§5/§6) und die
     **bare Tags ohne `ghcr`-Präfix** — davon gibt es zwei: das Beispiel in
     §Versionen und Tags (`:vX.Y.Z`) und der Docker-Hub-Pull
     `pt9912/d-check:vX.Y.Z` in §Docker-Image. Beide sind vom `versions`-Gate
     **nicht** erfasst und driften still. **Braucht das Feature
     eine §4-Aufgabe, schreibe eine eigene** — nach dem
     [Benutzerhandbuch-Standard](benutzerhandbuch-standard.md) §5
     (Ausgangslage/Ziel/Vorgehen/Ergebnis) — und **hänge sie nicht an eine
     bestehende §4-Aufgabe an**. Kein Gate erzwingt Aufgabenorientierung; genau
     durch Anhängen wuchs §4.12 auf ~330 Zeilen / 8 Themen, bis es aufgetrennt
     werden musste.
     **Die Regel gilt der Klasse, nicht dem Kapitel:** *jeder* gegliederte
     Fließtext dieses Repos, den ein Release anfasst — Handbuch in **allen**
     Kapiteln, [`operations.md`](operations.md), diese Datei und der
     [Benutzerhandbuch-Standard](benutzerhandbuch-standard.md). Der Prüfsatz
     ist eine Frage an die Überschrift: *nennt sie alles, was unter ihr steht?*
     Nicht die Länge ist der Defekt, sondern die **Unauffindbarkeit** — ein
     Abschnitt darf lang sein, wenn seine Überschrift ehrlich ist.
     **Zweimal gemessen, zweimal zugetroffen:** §4.12 wuchs durch Anhängen auf
     ~330 Zeilen / 8 Themen; §5 „Zitate und Zeilen-Referenzen gegen ihre Quelle
     prüfen" nannte zwei Fähigkeiten und trug sechs Module über 183 Zeilen,
     darunter `vcs`, `commits`, `planning` und `tracked`. Beide Male stand die
     Regel danach **nur für den Ort** da, an dem sie wehtat — deshalb steht sie
     jetzt für die Klasse.
   - **README — beide Sprachfassungen synchron halten:** `README.de.md` (Deutsch,
     **kanonische Quelle — zuerst ändern**) und danach `README.md` (Englisch,
     nachübersetzen). Bei einem neuen Modul in **jeder** Fassung: (a) die
     **Dogfooding-Zeile** unter §Warum d-check — sie nennt eine **Zahl** („im
     Vollausbau (N Module inkl. …)") **und** eine Enumeration, und beide zählen
     die Module der [`.d-check.yml`](../../.d-check.yml), nicht alle
     existierenden; und (b) die **Modul-Liste** unter §Was ist d-check
     ergänzen — bei einer neuen **Bedingung** des Moduls `structure` zusätzlich
     deren Anzahl („bis zu N Bedingungen") in derselben Aufzählung. Das
     Intro-Framing prüfen, falls ein
     Nicht-Referenz-Modul (Content-Drift/Immutabilität/Versions/Traceability/Planning)
     neu ist. d-check prüft `links`/`anchors`/`ids`/`versions` in beiden READMEs, aber
     **nicht** die inhaltliche DE↔EN-Synchronität der Prosa.
     **Diese Zeilen standen bis v0.70.0 falsch da, und der Grund war diese
     Checkliste:** sie nannte eine „Status-Zeile (*alle N Regelmodule …*)", die
     in **keiner** der beiden Fassungen existiert. Wer sie suchte, fand nichts
     und hielt das für Erledigung — die Dogfooding-Zahl blieb zwei Monate bei
     „acht Module", während die Konfiguration elf führte. **Eine Checkliste, die
     auf eine nicht existierende Zeile zeigt, ist schlimmer als keine.**
   - **Das Datum kommt aus dem Kalender, nicht aus der Zeile darüber.** Ein
     Release berührt sieben Datumsstempel (`version.md` §Aktuell **und** §Verlauf,
     `CHANGELOG`-Überschrift, Handbuch-Kopf **und** §11-Zeile, die Historien von
     Lastenheft und Spezifikation) — und wer sie schreibt, schreibt sie von der
     Vorgängerzeile ab. Zieht sich eine Welle über mehr als einen Tag, datiert
     das Release damit auf den Tag ihres Beginns. Kein Gate sieht das;
     `git log -1 --date=short --format=%ad` ist das Orakel.
   - **Eine GEÄNDERTE Zusage zählt wie ein neues Feature.** Die Punkte oben
     lesen sich, als gehe es nur um Zuwachs (neues Modul, neue Option, neuer
     Abschnitt). Ändert ein Release die **Bedeutung** einer bestehenden Zusage,
     steht sie nicht in einem neuen Abschnitt, sondern mitten im alten — im
     Handbuch dort, wo das Modul beschrieben ist, und in derselben Zeile beider
     READMEs. Genau diese Stellen liest ein Konsument, dessen bisher grüner Lauf
     nach dem Update rot wird.
   - **Operations-Referenz** ([`operations.md`](operations.md)) — bei einem
     neuen Modul die Modul-Enumeration der `--enable`/`--disable`-Zeile, bei
     einer neuen CLI-Option die Optionen-Tabelle ergänzen. Kein Gate prüft
     diese Enumerationen — die Modul-Liste blieb so von v0.25 bis v0.37 still
     bei acht Modulen stehen.
5. **`make ci`** lokal grün fahren (Pre-Tag-De-Risk, „grün = Boden"), erst dann
   taggen.

Der Digest-Pin in Handbuch §2 entsteht **nach** dem Tag (er existiert erst nach
dem GHCR-Push) als Folge-Commit.


## Vorbedingungen (einmalig, im Konto des Betreibers)

Der Docker-Hub-Spiegel ist **fail-closed**
([`DC-FA-DIST-002`](../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel),
[ADR-0064](../plan/adr/0064-dockerhub-spiegel-fail-closed.md)) — fehlt eine
dieser Vorbedingungen, **schlägt jedes Release fehl**, und zwar erst *nach* dem
GHCR-Push. Der Abbruch nennt dann den bereits veröffentlichten GHCR-Digest.

- **Repository** `pt9912/d-check` auf Docker Hub. Es entsteht beim ersten Push;
  Sichtbarkeit (öffentlich) und Beschreibung setzt der Betreiber dort.
- **Zwei Secrets** im GitHub-Repository: `DOCKERHUB_USERNAME` und
  `DOCKERHUB_TOKEN` — ein Access-Token, **nicht** das Konto-Passwort.
  **Scope `read/write/delete`, nicht nur `read/write`.** Der Spiegel-Push
  gelingt schon mit `read/write`; der **Metadaten-PATCH** der
  Beschreibungs-Seite nicht — er wird mit `Forbidden` abgelehnt, und weil der
  Schritt `continue-on-error` trägt, bleibt das Release grün und die Seite
  leer. Gemessen beim ersten Spiegel-Lauf (`v0.65.0`); die Anforderung steht in
  der Action-Doku (*„Docker Hub password or Personal Access Token with
  `read/write/delete` scope"*). Der Scope eines bestehenden Tokens lässt sich in
  der Docker-Hub-Oberfläche **ändern** — ein neues Token und ein Ersetzen des
  Secrets sind dafür **nicht** nötig; der Token-Wert bleibt derselbe.
- **Zwei Repository-Einstellungen für den Dependabot-Kanal**
  ([ADR-0067](../plan/adr/0067-dependabot-als-hebender-kanal.md)):
  **Dependabot-Alerts** und **Dependabot-Security-Updates**. Sie stehen in
  keiner Datei — sie sind Schalter im Repository. **Seit 2026-08-28 beide an**
  (`automated-security-fixes` meldet `enabled:true`, `vulnerability-alerts`
  HTTP 204); davor waren sie aus, und der Kanal hob nur, was **veraltet** war,
  nicht was **verwundbar** ist. Wer dieses Repo forkt oder neu aufsetzt, setzt
  sie erneut: ohne sie öffnet ein CVE, für das es kein neues Upstream-Release
  gibt, **keinen** PR.
- Optional die Repository-Variable `DOCKERHUB_IMAGE`, falls ein Fork woandershin
  spiegeln soll; ohne sie gilt `pt9912/d-check`.

**Die Hub-Beschreibungsseite wird aus dem Repo gesetzt** — Kurztext und
Overview-Seite kommen aus
[`packaging/dockerhub/`](../../packaging/dockerhub/README.md), die **Kategorie**
bleibt manuell im Web-UI und steht dort als Text. Beide Release-Schritte tragen
`continue-on-error`: ihr Fehlschlag lässt das Release grün
([ADR-0065](../plan/adr/0065-spiegel-gleichheit-ist-der-config-digest.md)
Punkt 5).

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
5. **Docker-Hub-Spiegel** — dasselbe lokale Bild wird nach
   `docker.io/pt9912/d-check` getaggt und gepusht, dieselbe Tag-Disziplin
   wie Schritt 4. Danach vergleicht der Schritt die **Config-Digests** beider
   Registries — aus den Registries gelesen, nicht aus dem lokalen Daemon — und
   bricht bei Ungleichheit ab
   ([`DC-FA-DIST-002`](../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel)).
   **Fail-closed:** jeder Fehlschlag hier macht das Release rot, obwohl
   GHCR bereits trägt — die Meldung nennt deshalb den veröffentlichten
   GHCR-Digest. Die Zugangsdaten werden **vor** dem Push geprüft, sonst
   scheiterte der Login mit seinem eigenen Text statt mit dieser Meldung.
6. **Hub-Darstellung** — Kurztext und Overview-Seite aus
   `packaging/dockerhub/`. Beide Schritte tragen `continue-on-error` und
   können das Release **nicht** rot machen; das Zeichen-Limit prüft
   stattdessen `make gates`.
7. **Digest-Pin** landet im Job-Summary und in den Notes des
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
