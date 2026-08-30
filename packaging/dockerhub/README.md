# Docker-Hub-Repository-Metadaten

Quelle für die Darstellung von [`pt9912/d-check`](https://hub.docker.com/r/pt9912/d-check)
auf Docker Hub. Docker Hub hat **drei getrennte Felder**; sie werden hier
unterschiedlich gepflegt:

| Feld auf Docker Hub | Quelle | Limit | Gepflegt durch |
|---|---|---|---|
| **Description** (Kurztext unter dem Repo-Namen) | [`description.txt`](description.txt) | 100 Zeichen | Release-Build — **beim ersten Lauf `Forbidden`, siehe unten** |
| **Repository overview** (Markdown-Seite) | [`overview.md`](overview.md) | 25.000 Zeichen | Release-Build, derselbe Step — **ebenso** |
| **Category** | dieses Dokument (siehe unten) | — | **manuell im Web-UI** |

Die ersten beiden **soll** `peter-evans/dockerhub-description` bei jedem
Release-Build setzen. Die Action hat **keinen** Input für die Kategorie —
deshalb steht die Entscheidung hier als Text, statt still im UI zu leben.

**Die beiden hochgeladenen Dateien sind ENGLISCH**, dieses Dokument nicht. Die
Hub-Seite ist die Außensicht und folgt darin
[`README.md`](../../README.md) (englisch, mit
[`README.de.md`](../../README.de.md) daneben); `description.txt` und
`overview.md` sind das einzige Repo-Material, das ungelesen auf einer fremden
Plattform steht. Was **hier** steht, ist Betriebswissen für dieses Repo und
bleibt deutsch wie der Rest der Doku.

## Der erste Lauf hat sie NICHT gesetzt — gemessen

Beim Release `v0.65.0` (Run 33139273535) meldete der Schritt `success`, sein Log
aber `##[error]Forbidden` beim `PATCH` an die Docker-Hub-API. Die API bestätigt
es: `description` ist leer, `full_description` fehlt ganz.

**Der `success` war ein Artefakt von `continue-on-error`** — der Schritt darf
das Release nicht rot machen (das ist gewollt,
[ADR-0065](../../docs/plan/adr/0065-spiegel-gleichheit-ist-der-config-digest.md)
Punkt 5), aber er verschluckte dabei auch die Meldung, dass er nichts bewirkt
hat. **Ab dem nächsten Release verschluckt er sie nicht mehr** — der Melder ist nach v0.65.0 entstanden und hat damit noch keinen Lauf gesehen: ein Folgeschritt liest
`steps.hubdesc.outcome` — `continue-on-error` setzt `conclusion: success`, lässt
`outcome` aber auf `failure` — und setzt eine **Warnung** mit der wahrscheinlichen
Ursache. Warnung, nicht Fehler: das Release soll an der Darstellung nicht
scheitern, nur nicht schweigen.

**Die Ursache ist belegt.** Der Push desselben Tokens funktioniert (das Bild
liegt auf Docker Hub); abgelehnt wird nur der Metadaten-`PATCH`. Die Action
verlangt dafür einen anderen Scope als der Push: *„Docker Hub password or
Personal Access Token with `read/write/delete` scope"* — ein Token mit
`read/write` pusht, darf die Beschreibung aber nicht ändern. Im Schwester-Repo
`pt9912/d-migrate` setzt dieselbe Action ihre Beschreibung erfolgreich; dort
trägt das Token den weiteren Scope.

**Behebung:** den Scope des bestehenden Tokens in der Docker-Hub-Oberfläche auf
`read/write/delete` **ändern** — das geht ohne Neuanlage, der Token-Wert bleibt
derselbe, und `DOCKERHUB_TOKEN` muss nicht ersetzt werden. Die Beschreibung
setzt die Action dann beim **nächsten** Release; rückwirkend geschieht nichts.
Bis dahin sind beide Felder **manuell** zu setzen, aus genau diesen zwei
Dateien. Der Inhalt ist damit weiter reviewbar; nur der Transport fehlt.

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
wc -m < packaging/dockerhub/description.txt   # muss <= 100 sein
```

**`-m`, nicht `-c`:** Docker Hub misst **Zeichen**, der Go-Test in
[`dockerhub_description_test.go`](../../internal/adapter/driven/configyaml/dockerhub_description_test.go)
ebenso (`utf8.RuneCountInString`), und der Release-Build auch. Eine
Byte-Messung wäre bei Nicht-ASCII **strenger** als die Regel und meldete einen
Text rot, den Docker Hub annimmt. Hier stand `-c`, solange der Kurztext
deutsche Umlaute trug — also genau in der Zeit, in der die beiden Zahlen
auseinandergingen.
