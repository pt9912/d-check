# Docker-Hub-Repository-Metadaten

Quelle für die Darstellung von [`pt9912/d-check`](https://hub.docker.com/r/pt9912/d-check)
auf Docker Hub. Docker Hub hat **drei getrennte Felder**; sie werden hier
unterschiedlich gepflegt:

| Feld auf Docker Hub | Quelle | Limit | Gepflegt durch |
|---|---|---|---|
| **Description** (Kurztext unter dem Repo-Namen) | [`description.txt`](description.txt) | 100 Bytes | Release-Build ([`release.yml`](../../.github/workflows/release.yml)) |
| **Repository overview** (Markdown-Seite) | [`overview.md`](overview.md) | 25.000 Bytes | Release-Build, derselbe Step |
| **Category** | dieses Dokument (siehe unten) | — | **manuell im Web-UI** |

Die ersten beiden setzt `peter-evans/dockerhub-description` bei jedem
Release-Build. Die Action hat **keinen** Input für die Kategorie — deshalb steht
die Entscheidung hier als Text, statt still im UI zu leben.

## Category

**Zu setzen: „Developer tools"** (`developer-tools`).

Begründung: d-check prüft Dokumentations-Referenzen in Entwickler-Repositories
und läuft als Gate in CI-Pipelines — die Kategorie beschreibt den Nutzen für den
Suchenden, nicht die Implementierungssprache. „Languages & frameworks" wäre
falsch: Go ist Implementierungsdetail
([ADR-0001](../../docs/plan/adr/0001-implementierungssprache.md)) und für
niemanden ein Grund, das Image zu ziehen.

Docker Hubs Taxonomie ist eine feste Liste; freie Schlagworte gibt es nicht.

## Warum die Darstellung nicht fail-closed ist

Der **Spiegel** ist fail-closed
([`DC-FA-DIST-002`](../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel)):
liegt das Bild nicht auf Docker Hub, ist das Release fehlgeschlagen. Die
**Darstellung** ist es nicht — sie läuft mit `continue-on-error`. Der
Unterschied ist der Gegenstand: das Bild ist die Zusage, der Beschreibungstext
ist Präsentation. Ein Release rot zu machen, weil ein Kurztext nicht gesetzt
werden konnte, verwechselte beides.

## `__VERSION__`

[`overview.md`](overview.md) trägt den Platzhalter `__VERSION__` in den
Aufruf-Beispielen; der Release-Build ersetzt ihn durch die Tag-Version, bevor er
die Seite hochlädt. Die Datei im Repo bleibt damit versionsfrei und muss bei
keinem Release angefasst werden — anders als die
`ghcr`-Pins in README und Handbuch, die
[`releasing.md`](../../docs/user/releasing.md) §Release-Prep aufzählt.

## Änderungen prüfen

Beide Dateien sind reiner Repo-Inhalt und laufen durch `make doc-check` wie jede
andere Doku. Das Byte-Limit der Description prüft der Release-Build
**fail-fast**, statt der Action das stille Abschneiden zu überlassen:

```bash
wc -c < packaging/dockerhub/description.txt   # muss <= 100 sein
```
