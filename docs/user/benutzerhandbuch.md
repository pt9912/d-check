# Benutzerhandbuch: d-check

**Handbuch-Version:** 1.61 · **Software-Version:** [v0.67.0](../../version.md#v0.67.0) ·
**Stand:** 2026-08-29 · **Autor:** pt9912

Dieses Handbuch folgt dem
[Benutzerhandbuch-Standard](benutzerhandbuch-standard.md): aufgabenbasiert,
zielgruppenorientiert, eindeutig. Eine knappe Optionsreferenz steht
zusätzlich in [operations.md](operations.md). d-check ist ein
Kommandozeilen-Werkzeug; „Bedienen" heißt hier **Befehle ausführen**, die
Beispiele zeigen die echte **Terminal-Ausgabe**.

## Inhalt

- [1. Einleitung](#1-einleitung)
- [2. Installation und Zugriff](#2-installation-und-zugriff)
- [3. Erste Schritte](#3-erste-schritte)
- [4. Aufgaben](#4-aufgaben)
- [5. Konfiguration](#5-konfiguration)
- [6. Regelmodule](#6-regelmodule)
- [7. Fehlerbehebung](#7-fehlerbehebung)
- [8. FAQ](#8-faq)
- [9. Glossar](#9-glossar)
- [10. Support und Lizenz](#10-support-und-lizenz)
- [11. Änderungshistorie](#11-änderungshistorie)

## 1. Einleitung

### Zweck der Software

d-check prüft die **Markdown-Dokumentation** eines Repositorys auf kaputte
Referenzen: lokale Links und Bildverweise, Heading-Anker, die Linkpflicht
für Kennungen, die Referenzrichtung zwischen Dokumentklassen und weitere
Doku-Hygiene. Ein Werkzeug für die gesamte Doku-Prüfung: ein Image, ein
Update-Pfad, repo-spezifisches Verhalten per Konfiguration (`.d-check.yml`).

d-check ist ein **Lese-Werkzeug**: Es verändert das geprüfte Repository
nie und macht außer in den optionalen Modulen `external` und `sources`
(beide opt-in, nie im Default-Lauf) keine Netzwerkzugriffe. Gleiche
Eingabe liefert immer dieselbe Ausgabe.

### Zielgruppe

Dieses Handbuch richtet sich an:

- **Repo-Maintainer**, die ihre Doku lokal sauber halten wollen.
- **CI-Integratoren**, die d-check als Prüfschritt einbinden.
- **AI-Harness-Projekte und Automatisierung**, die einen verlässlichen,
  maschinenlesbaren Doku-Sensor brauchen.

Vorkenntnisse: Kommandozeile und Docker auf Grundniveau. Kein Go nötig.

### Voraussetzungen

- **Docker** (empfohlener Weg) oder die Fähigkeit, ein OCI-Image zu ziehen.
- Ein Repository mit Markdown-Dateien.
- Optional eine Konfigurationsdatei `.d-check.yml` in der Repo-Wurzel.

## 2. Installation und Zugriff

### Docker-Image (empfohlen)

d-check wird als Container-Image über die GitHub Container Registry (GHCR)
verteilt. Es braucht keine Installation — Sie ziehen und starten das Image:

```bash
docker pull ghcr.io/pt9912/d-check:v0.67.0
```

Das Image läuft als Nicht-root-Prozess; ein **read-only**-Mount des
Repositorys genügt, weil d-check nie schreibt.

**Zweiter Bezugsweg — Docker Hub, seit v0.67.0.** Dasselbe
Bild wird zusätzlich als `docker.io/pt9912/d-check` gespiegelt — **kein zweiter
Bau**: der **Config**-Digest ist auf beiden Registries gleich
([`DC-FA-DIST-002`](../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel)).
Der **Manifest**-Digest ist es **nicht** — er hängt an der Blob-Kompression des
jeweiligen Registrys; wer per Digest pinnt, nimmt den Digest **der Registry, aus
der er zieht**. GHCR bleibt die Quelle; Docker Hub folgt ihr.

```bash
docker pull pt9912/d-check:v0.67.0
```

### Versionen und Tags

- `:v0.67.0` — eine feste Version (empfohlen für reproduzierbare Läufe; die jeweils
  aktuelle steht in [version.md](../../version.md#aktuell)).
- `:latest` — die jeweils neueste **stabile** Version. Vorabversionen
  (Prereleases, z. B. `v1.0.0-rc1`) erhalten **kein** `:latest`; für
  CI-Pipelines pinnen Sie ohnehin auf eine feste Version oder den Digest
  (Komfort-Einstieg, nicht für reproduzierbare Läufe).

Für vollständig reproduzierbare CI-Läufe pinnen Sie auf den Image-Digest
(in den Release-Notes je Version angegeben):

```bash
docker run --rm -v "$PWD:/repo:ro" \
  ghcr.io/pt9912/d-check@sha256:c6c1465b94f07ab24439665be40a3107df51a3c0c62d0159a4e4a915fb03ca7c
```

### Native Nutzung

Das ausgelieferte Binary ist statisch und liegt im Image unter `/d-check`.
Der unterstützte und dokumentierte Weg ist der Container-Aufruf; ein
nativer Aufruf verhält sich byte-identisch (das wird vor jeder
Veröffentlichung geprüft).

## 3. Erste Schritte

### Schneller Einstieg

Prüfen Sie das aktuelle Verzeichnis:

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0
```

d-check mountet Ihr Repository nach `/repo` und prüft es. Eine typische
Ausgabe bei sauberer Doku:

```text
d-check: 42 Datei(en) geprüft, 0 Befund(e)
```

Bei kaputten Referenzen meldet d-check je einen Befund pro Zeile und endet
mit Exit-Code 1:

```text
docs/anleitung.md:12	fehlt.md	target-missing	Linkziel existiert nicht
d-check: 42 Datei(en) geprüft, 1 Befund(e)
```

Jede Befund-Zeile hat das Format `Datei:Zeile  Ziel  Grund-Code`. Trägt der
Befund eine **Erläuterung**, folgt sie als **vierte** tab-getrennte Spalte —
wie oben. Ein Befund bleibt dabei **eine** Zeile, und die ersten drei Felder
bleiben, wie sie waren: wer die Ausgabe auf Tab trennt und Feld 1–3 liest,
liest weiter dasselbe. Ein Befund ohne Erläuterung ist dreispaltig.

### Grundkonzepte

- **Regelmodule:** Die Prüfung ist in Module gegliedert. Standardmäßig
  laufen `links` und `anchors`; weitere Module schalten Sie bei Bedarf zu
  (siehe [Abschnitt 6](#6-regelmodule)).
- **Exit-Codes:** `0` = keine Befunde, `1` = mindestens ein Befund,
  `2` = Nutzungs- oder Umgebungsfehler. CI bricht bei `1` und `2` ab.
- **Ausgabe-Trennung:** Befunde stehen auf **stdout**, die Zusammenfassung
  und Meldungen auf **stderr**. So lassen sich Befunde sauber
  weiterverarbeiten.
- **Read-only:** d-check schreibt nie ins Repository. Selbst die
  Reparatur-Funktion gibt nur einen Patch aus, den Sie selbst anwenden.
- **Markdown-Grammatik (CommonMark/GFM):** d-check liest Ihre Dateien wie ein
  gängiger Renderer. Fenced-Code — Blöcke aus drei Backticks oder drei Tilden
  (`~~~`) — überspringen **alle** Module. Eine Feinheit von CommonMark: öffnet
  eine Backtick-Fence-Zeile ihren Block mit einer Infozeile, die selbst noch
  einen **Backtick** enthält, ist sie **kein** Block-Öffner, sondern Fließtext —
  Links und Kennungen darauf und dahinter werden geprüft. So bleibt ein
  erklärender Satz *über* einem Codeblock, der dessen Sprache in Backticks nennt,
  sichtbar (bis v0.46.0 verschwand er samt allem Folgenden aus jedem Modul). Für
  Tilden-Fences gilt die Ausnahme nicht.

## 4. Aufgaben

Jede Aufgabe nennt Ziel, Voraussetzung, Vorgehen und das erwartete
Ergebnis.

### 4.1 Die Dokumentation eines Repos prüfen

**Ziel:** kaputte lokale Links und Heading-Anker finden.
**Voraussetzung:** Docker; ein Repository im aktuellen Verzeichnis.

**Vorgehen:**

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0
```

**Ergebnis:** Exit-Code 0 und „0 Befund(e)" bei sauberer Doku; sonst die
Befund-Zeilen und Exit-Code 1.

**Hinweise:** Geprüft werden standardmäßig `docs/`, `spec/` und alle
`*.md` direkt in der Repo-Wurzel. Build- und Fremdverzeichnisse (z. B.
`.git`, `node_modules`) werden immer übersprungen.

### 4.2 d-check in CI einbinden

**Ziel:** die Doku-Prüfung als CI-Schritt erzwingen.
**Voraussetzung:** eine CI, die Docker-Container ausführen kann.

**Vorgehen** (Beispiel als Shell-Schritt):

```bash
docker run --rm --network none -v "$PWD:/repo:ro" \
  ghcr.io/pt9912/d-check:v0.67.0
```

**Ergebnis:** Der Schritt ist grün bei Exit-Code 0 und rot bei 1 oder 2 —
die CI bricht bei kaputten Referenzen ab.

**Hinweise:** `--network none` ist sinnvoll, solange kein Netz-Modul
(`external` oder `sources`) aktiv ist (d-check braucht dann kein Netz). Pinnen Sie für
reproduzierbare Läufe auf den Image-Digest (siehe
[Abschnitt 2](#2-installation-und-zugriff)).

### 4.3 Eine Konfiguration anlegen

**Ziel:** eine `.d-check.yml` als Startgerüst erzeugen.
**Voraussetzung:** keine.

**Vorgehen:**

```bash
docker run --rm ghcr.io/pt9912/d-check:v0.67.0 --print-config > .d-check.yml
```

**Ergebnis:** Eine kommentierte `.d-check.yml` im aktuellen Verzeichnis.

**Hinweise:** `--print-config` liest das Repository **nicht** und schreibt
selbst nichts — es gibt nur das Gerüst auf stdout aus; die Umleitung
(`>`) legt die Datei an. Das Gerüst dokumentiert die verfügbaren Module
und Optionen als Kommentare.

### 4.4 Eine Konfiguration aus dem Repo vorschlagen lassen

**Ziel:** Kennungs-Muster (Modul `ids`) automatisch aus den Quellen
ableiten, in denen Kennungen definiert sind.
**Voraussetzung:** Dokumente, die Kennungen als Überschriften definieren.

**Vorgehen** (Quellen kommagetrennt):

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0 \
  --suggest-config spec/,docs/plan/adr/ > .d-check.yml
```

**Ergebnis:** Ein vorgeschlagenes Gerüst auf stdout mit abgeleiteten
`ids`-Mustern und den opt-in-Modulen, die im Repo echtes Signal liefern.

**Hinweise:** Der Vorschlag ist beratend — prüfen und verengen Sie ihn.
d-check liest dabei nur, es schreibt nichts.

**Harness-Vorlage (`ai-harness` / `ai-harness-init`):** Für ein Repo nach
dem ai-harness-course-Standard liefern zwei reservierte Quellen ein fertiges
Gerüst — ohne die Quellen einzeln aufzulisten. Welche Sie nehmen, hängt von Ihrer
Ausgangslage ab:

- **Frisches Repo → `ai-harness-init`** (Voll-Kanon, alle Blöcke aktiv).
  Das ist Ihr Zielbild: legen Sie die Struktur an (`spec/`,
  `docs/plan/adr/`, …), dann läuft d-check.

  ```bash
  docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0 \
    --suggest-config ai-harness-init > .d-check.yml
  ```

- **Bestehendes Repo → `ai-harness`** (repo-bewusst): nur Pfade, die Ihr
  Repo schon hat, sind aktiv; fehlende erscheinen auskommentiert mit
  Hinweis (Ihre TODO-Liste). Läuft sofort.

  ```bash
  docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0 \
    --suggest-config ai-harness > .d-check.yml
  ```

Beide sind beratend und kombinierbar mit echten Quellen (z. B. `ai-harness`
zusätzlich zu `spec/lastenheft.md`).

Das Gerüst enthält die kanonischen `ids`-Muster, die `matrix`-Referenzrichtung
und das **fixe Standard-Modulset** (`links`, `anchors`, `ids`, `matrix`,
`codepaths`, `spans`, `hostpaths`) samt einem **repo-bewussten `planning`-Block**
(aktiv, sobald die Roadmap existiert). Situative Module ohne ableitbare
Konfiguration bleiben aus — `vcs`/`commits` (die eine Commit-Range brauchen)
liefert `--print-mk` als Makefile-Target, den Rest listet `--print-config`.

**Kennungs-Präfix (`--id-prefix`):** Das Anforderungs-`ids`-Muster ist
projektspezifisch — nur sein Präfix wechselt pro Repo (d-check: `DC`,
a-check: `AC`, …). Geben Sie es mit `--id-prefix` an:

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0 \
  --suggest-config ai-harness-init --id-prefix AC > .d-check.yml
```

Im Modus `ai-harness` wird das Präfix **ohne** `--id-prefix` aus dem
vorhandenen `spec/lastenheft.md` abgeleitet (mehrere verschiedene Präfixe ⇒
Fehler — dann `--id-prefix` setzen). Ohne Angabe **und** ohne Ableitung
(typisch `ai-harness-init` im leeren Repo) erscheint ein markierter
Platzhalter `<PREFIX>` mit `# TODO` — **kein** stilles `DC-`; ersetzen Sie
ihn durch Ihr Projekt-Präfix.

### 4.5 Regelmodule zu- und abschalten

**Ziel:** ein Modul aktivieren oder deaktivieren.
**Voraussetzung:** keine.

**Vorgehen** (Optionen sind wiederholbar; Kommandozeile schlägt
Konfiguration):

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0 \
  --enable ids --disable anchors
```

**Ergebnis:** Nur die gewünschten Module laufen.

**Hinweise:** Dauerhaft schalten Sie Module in der `.d-check.yml`
(`modules: [...]`). Ein unbekannter Modulname endet mit Exit-Code 2 und
einer Liste gültiger Namen.

### 4.6 Kennungs-Linkpflicht prüfen (Modul `ids`)

**Ziel:** sicherstellen, dass nackte Kennungen (z. B. ADR- oder
Anforderungs-Nummern) im Fließtext als Link auf ihre Definition
ausgeführt sind.
**Voraussetzung:** `ids` aktiv und Muster konfiguriert (siehe
[Abschnitt 5](#5-konfiguration)).

**Vorgehen:**

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0 \
  --enable ids
```

**Ergebnis:** Jede nackte, unverlinkte Kennung wird als `id-unlinked`
gemeldet.

**Hinweise:** Ausnahmen steuern Sie je Muster über `exempt-paths` (Globs)
oder zeilenweise über den Marker `d-check:ignore` in einem
HTML-Kommentar.

### 4.7 Referenzrichtung prüfen (Modul `matrix`)

**Ziel:** erzwingen, dass Referenzen nur in erlaubter Richtung zwischen
Dokumentklassen zeigen (z. B. Spezifikation nicht „abwärts" auf
Architekturentscheidungen) und nicht auf abgelöste Dokumente.
**Voraussetzung:** `matrix` aktiv, Klassen und Regeln konfiguriert.

**Vorgehen:**

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0 \
  --enable matrix
```

**Ergebnis:** Verbotene Klassen-Referenzen erscheinen als
`matrix-forbidden`, Referenzen auf Dokumente mit verbotenem Status (etwa
abgelöst) als `matrix-inactive`.

**Auch *innerhalb* einer Klasse eine Rangfolge erzwingen (`order`/`direction`).**
Manchmal sollen selbst Dokumente **derselben** Klasse nicht „abwärts" verweisen —
etwa das Lastenheft nicht auf Spezifikation oder Architektur, obwohl alle drei
die Klasse `spec` sind. Geben Sie der Klasse dazu eine Rangfolge `order`
(Pfad-Globs, **autoritativste Schicht zuerst**) und `direction: no-downward`:

```yaml
matrix:
  classes:
    - name: spec
      paths: [spec/lastenheft.md, spec/spezifikation.md, spec/architecture.md]
      order: [spec/lastenheft.md, spec/spezifikation.md, spec/architecture.md]
      direction: no-downward
```

d-check prüft dann auch klasseninterne Referenzen: der Rang einer Datei ist ihr
erster passender `order`-Glob; ein Verweis von einer höher- auf eine
niederrangige Datei (auch über mehrere Stufen) erscheint als `matrix-downward`
(Dateien ohne `order`-Treffer sind rangfrei und werden nicht geprüft). `order`
und `direction` gehören zusammen — eines ohne das andere (oder ein unbekannter
`direction`-Wert) ist ein Konfigurationsfehler (Exit 2), damit eine
Richtungsregel nie still wirkungslos ist. Ohne beide Felder verhält sich
`matrix` unverändert.

**Auch bare ID-Tokens im Fließtext als Referenz prüfen (`token`).** Eine Referenz
muss kein Markdown-Link sein — eine Slice-Kennung, die **roh** im Text eines ADR
steht, ist ebenso eine Referenz, die `matrix` standardmäßig aber nicht sieht.
Damit solche Tokens erfasst werden, geben Sie der Klasse ein `token`-Regex; für
eine **erlaubte** Ausnahme (etwa „verifiziert in slice-042" als reinen
Verifikations-Zeiger, keine Entscheidungsgrundlage) setzen Sie den Marker
`<!-- d-check:status-provenance -->` auf dieselbe Zeile:

```yaml
matrix:
  classes:
    - {name: adr, paths: ["docs/plan/adr/[0-9]*.md"]}
    - name: slice
      paths: ["docs/plan/planning/**/slice-*.md"]
      token: 'slice-\d{3}'                       # Slice-Kennung im Text erkennen
  rules:
    - {from: adr, to: slice, allow: false}       # ADR nennt Slice nur als Provenance
  exempt-paths: ["docs/plan/adr/0001-*.md"]      # Alt-Dateien ganz ausnehmen
```

`matrix` behandelt dann ein solches Token im Körper eines Dokuments einer
anderen Klasse als Referenz, auf die dieselbe Regel greift (Token in
Markdown-Links und Fenced-Code zählen nicht) — eine verbotene Kante meldet
`matrix-forbidden`. `exempt-paths` (Globs) nimmt **ganze Dateien** von der
`matrix`-Prüfung aus — nützlich, um unveränderliche Bestandsdokumente zu
grandfathern. Ein nicht kompilierbares `token`-Regex ist ein
Konfigurationsfehler (Exit 2); ohne `token` verhält sich `matrix` unverändert.

### 4.8 Externe Links prüfen (Modul `external`)

**Ziel:** die Erreichbarkeit externer (HTTP-)Links prüfen.
**Voraussetzung:** Netzwerkzugang; `external` ausdrücklich aktiviert.

**Vorgehen** (ohne `--network none`, da Netz gebraucht wird):

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0 \
  --enable external
```

**Ergebnis:** Nicht erreichbare Links erscheinen als `external-status`,
Zeitüberschreitungen als `external-timeout`, zu viele Weiterleitungen als
`external-redirects`.

**Hinweise:** `external` und `sources` sind die **einzigen** Module, die
Netzwerkverbindungen öffnen (beide opt-in, nie im Default-Lauf; siehe
[Abschnitt 5](#5-konfiguration) zum Content-Pin `sources`). Ohne eines davon
bleibt d-check vollständig netzlos.

### 4.9 Befunde verstehen — Diagnose-Modus (`--doctor`)

**Ziel:** Befunde erklärt und nach Datei gruppiert lesen, mit konkreten
Fix-Vorschlägen.
**Voraussetzung:** keine.

**Vorgehen:**

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0 \
  --enable ids --doctor
```

**Ergebnis:** Statt der knappen Befund-Zeilen erscheint eine Diagnose:

```text
d-check Diagnose — 1 Befund(e) in 1 Datei(en):

docs/anleitung.md
  Z. 12 · Kennung im Fließtext ohne Markdown-Link auf ihre Definition [ids]
      Stelle: ADR-0007
      Hinweis: Kennung ohne Link auf ihre Definition
      Fix-Kandidat: ADR-0007 → [`ADR-0007`](plan/adr)
        (Kennung als Markdown-Link auf ihre Definition (docs/plan/adr/) ausführen; Anker ggf. ergänzen)
```

**Hinweise:** Die `Hinweis:`-Zeile trägt die **Erläuterung** des Befunds — bei
den meisten Modulen eine feste Meldung, bei `structure` der selbst verfasste
`hint`. Sie ist die Gegengröße zum **Fix-Kandidaten**: jene ist geschrieben,
dieser abgeleitet. **Bei manchen Modulen wiederholt sie den Grund-Klartext
fast wörtlich** — im Beispiel oben —, weil beide Texte historisch dieselbe
Aussage tragen; in der knappen Befund-Zeile, die keinen Klartext hat, trägt
sie in jedem Fall. `--doctor` ist read-only und gibt nur auf stdout aus. Mit
zusätzlichem `--json` wird dieselbe Diagnose **maschinenlesbar**
ausgegeben (siehe unten). Nicht mit `--repair` kombinierbar (sonst
Exit-Code 2). Einen Fix-Kandidaten gibt es nur dort, wo er **eindeutig**
ableitbar ist.

**Maschinenlesbare Diagnose (`--doctor --json`):** Statt der Prosa ein
JSON-Dokument wie bei `--json` (Abschnitt 4.11),
dessen `findings` je Eintrag zusätzlich `reasonText` (Grund-Klartext) und
`fixCandidate` (`{original, replacement, note}` oder `null`) tragen:

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0 \
  --enable ids --doctor --json
```

```json
{
  "findings": [
    {
      "file": "docs/anleitung.md",
      "line": 12,
      "rule": "ids",
      "target": "ADR-0007",
      "reason": "id-unlinked",
      "reasonText": "Kennung im Fließtext ohne Markdown-Link auf ihre Definition",
      "fixCandidate": {
        "original": "ADR-0007",
        "replacement": "[`ADR-0007`](plan/adr)",
        "note": "Kennung als Markdown-Link auf ihre Definition (docs/plan/adr/) ausführen; Anker ggf. ergänzen"
      }
    }
  ],
  "summary": { "filesChecked": 1, "findingCount": 1 },
  "exitCode": 1
}
```

`fixCandidate` ist `null`, wo kein Fix eindeutig ableitbar ist (nicht
weggelassen — das `null` ist die Aussage „kein eindeutiger Fix").

**Die drei Ausgaben im Vergleich:**

| Aufruf            | stdout                                                     | wofür                                        |
| ----------------- | ---------------------------------------------------------- | -------------------------------------------- |
| `--json`          | JSON: knappe Befunde (`findings`/`summary`/`exitCode`)     | CI/Maschine, reine Befundliste               |
| `--doctor`        | Prosa: gruppierte Diagnose mit Klartext + Fix-Kandidaten   | Mensch, zum Verstehen                        |
| `--doctor --json` | JSON: Befunde zusätzlich mit `reasonText` + `fixCandidate` | Maschine, die die Diagnose weiterverarbeitet |

Dieselben maschinenlesbaren Varianten gibt es als **YAML** (`--yaml` bzw.
`--doctor --yaml`) — gleiche Struktur, nur YAML statt JSON.

### 4.10 Befunde als Patch beheben (`--repair`)

**Ziel:** behebbare Befunde als anwendbaren Patch ausgeben.
**Voraussetzung:** `git` zum Anwenden des Patches.

**Vorgehen** (Patch erzeugen, sichten, anwenden, aufräumen):

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0 \
  --enable ids --repair > fix.patch
# fix.patch sichten (besonders bei --repair-broad), dann anwenden:
git apply fix.patch
rm fix.patch
```

**Ergebnis:** Ein `git apply`-kompatibler unified diff auf stdout. d-check
selbst schreibt nichts — Sie wenden den Patch an.

**Was `--repair` behebt** (alle anderen Befundarten bleiben Befund):

| Befundart                                                                                                                                                                                                                     | Reparatur                                                                                                                                                                       | Stufe                    |
| ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------ |
| `id-unlinked`                                                                                                                                                                                                                 | nackte Kennung → Markdown-Link auf ihre Definition; nur **nackte Prosa**-Vorkommen (Inline-Code oder bereits verlinkt bleiben unangetastet)                                     | konservativ (`--repair`) |
| `target-missing`                                                                                                                                                                                                              | Link → **verschobene** Markdown-Datei gleichen, im Repo **eindeutigen** Namens; **keine** Umbenennung (anderer Name), keine Nicht-Markdown-Ziele, mehrdeutige Namen → kein Edit | breit (`--repair-broad`) |
| `anchor-missing`, `repo-escape`, `symlink`, `codepath-missing`, `matrix-inactive`, `matrix-forbidden`, `external-status`, `external-timeout`, `external-redirects`, `span-unclosed`, `span-nested-link`, `fence-unclosed`, `closure-note-placeholder`, `closure-note-ambiguous`, **jeder** `section-*`, `hostpath-forbidden` | — kein Auto-Fix, von Hand beheben                                                                                                                                               | —                        |

**Hinweise:**

- Best-Guess-Reparaturen der breiten Stufe (`--repair-broad`) sind
  **review-pflichtig**; ihre Markierung erscheint auf **stderr**, damit der
  Patch auf stdout `git apply`-rein bleibt. Lesen Sie den Patch vor dem
  Anwenden.
- Nicht reparierte Befundarten erscheinen weiter in der normalen Ausgabe
  bzw. unter `--doctor` (mit Klartext-Grund, aber ohne Fix-Kandidat) und
  sind von Hand zu beheben.
- Nicht mit `--json` oder `--doctor` kombinierbar.
- **`fix.patch` ist temporär:** nach dem Anwenden löschen
  (`rm fix.patch`) — sonst bleibt die Datei im Repo liegen und kann
  versehentlich mitcommittet werden. Die Pipe-Variante unten kommt ganz
  ohne Zwischendatei aus.
- **Einzeiler (Pipe):** Weil der Patch auf stdout und nur
  Markierung/Zusammenfassung auf stderr gehen, können Sie direkt pipen:

  ```bash
  docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0 \
    --enable ids --repair | git apply
  ```

  Das umgeht aber die Sichtprüfung. Bei `--repair-broad`
  (review-pflichtig) und in CI-Shells mit `set -o pipefail` — dort lässt
  der Befund-Exit-Code 1 die Pipe scheitern, obwohl `git apply` ok war —
  ist die Datei-Variante (oben) die bessere Wahl.

### 4.11 Maschinenlesbare Ausgabe (`--json` / `--yaml`)

**Ziel:** Befunde automatisiert weiterverarbeiten.
**Voraussetzung:** keine.

**Vorgehen:**

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0 --json
```

**Ergebnis:** Ein JSON-Dokument auf stdout mit den Feldern `findings`,
`summary` und `exitCode`:

```json
{
  "findings": [
    { "file": "docs/anleitung.md", "line": 12, "target": "fehlt.md", "rule": "links", "reason": "target-missing", "message": "Linkziel existiert nicht" }
  ],
  "summary": { "filesChecked": 42, "findingCount": 1 },
  "exitCode": 1
}
```

**Dasselbe als YAML (`--yaml`):** identische Struktur, nur YAML statt JSON
(`--json` und `--yaml` schließen sich gegenseitig aus):

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0 --yaml
```

<!-- d-check-test:not-config: --yaml-Ausgabe-Beispiel, kein .d-check.yml-Input -->
```yaml
findings:
  - file: docs/anleitung.md
    line: 12
    target: fehlt.md
    rule: links
    reason: target-missing
    message: Linkziel existiert nicht
summary:
  filesChecked: 42
  findingCount: 1
exitCode: 1
```

**Hinweise:** Mit `--json` oder `--yaml` enthält stdout ausschließlich das
Dokument. Auch `--doctor --json` bzw. `--doctor --yaml` geben die Diagnose
maschinenlesbar aus (je Eintrag zusätzlich `reasonText`/`fixCandidate`,
siehe Abschnitt 4.9).

### 4.12 Anforderungs-Abdeckung prüfen — Traceability-Matrix (`--trace`)

**Ziel:** sehen, welche Anforderung von welchen Architekturentscheidungen
(ADRs), Umsetzungs-Slices und kuratierten Coverage-Quellen referenziert wird —
und welche Anforderung niemand abdeckt (`WAISE`; die genaue Bedeutung steht
unten im Ergebnis).
**Voraussetzung:** ein Repo nach Harness-Konvention (Anforderungen im
Lastenheft, ADRs unter `docs/plan/adr/`, Slices unter `docs/plan/planning/`).

**Unterstützte Definitionssyntax:** Standardmäßig erkennt `--trace` eine
Anforderung aus einer ATX-Markdown-Überschrift (`#` bis `######`, außerhalb von
Fenced-Code-Blöcken). Die ID muss das **erste vollständige Token** der
Überschrift sein und als Ganzes auf `trace.requirements.id-pattern` passen:

```markdown
### F-1 — Repository-Struktur

Das Repository MUSS alle Hauptbestandteile enthalten.
```

Im Defaultformat `headings` gelten Tabellenzeilen, Listen, Fließtext und
Setext-Überschriften nicht als Anforderungsdefinition. Eine angepasste ID-Regex
ändert nur die zulässige Kennungsgestalt, **nicht** diese Dokumentgrammatik.
Beispielsweise definiert dieser verbreitete Brownfield-Bestand ohne weitere
Formatkonfiguration keine Anforderung:

```markdown
| ID  | Modalität | Anforderung                                      |
| --- | ---------- | ------------------------------------------------ |
| F-1 | Muss       | Das Repository enthält alle Hauptbestandteile. |
```

Ein tabellarisches Lastenheft (`format: table` mit Spalten-Bindung) und die
Sonderfälle der Tabellen-Grammatik (Tabellengrenze, Direktiven-Marker in einer
Zeile) sind **Konfiguration**; die Felder, Regeln und Fehlerbilder stehen in §5.

**Vorgehen:**

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0 --trace
```

**Ergebnis:** eine Markdown-Tabelle auf stdout — je Anforderung Titel,
referenzierende ADRs, referenzierende Slices und eine Status-Spalte. `WAISE`
bedeutet exakt: **weder eine Slice- noch eine konfigurierte Coverage-Referenz**.
Eine ADR-Referenz allein verhindert den Waisenstatus nicht:

```text
# Requirements Traceability Matrix

| Anforderung   | Titel                            | ADRs | Slices    | Status |
| ------------- | -------------------------------- | ---- | --------- | ------ |
| DC-FA-CLI-009 | Requirements Traceability Matrix | —    | slice-036 | ok     |
| DC-FA-CLI-010 | Makefile-Fragment ausgeben       | —    | slice-038 | ok     |

2 Anforderung(en), 0 Waise(n).
```

**Hinweise:** `--trace` ist read-only und arbeitet nur auf der Doku
(Lastenheft/ADRs/Planning) — keine Code-Prüfung. Mit `--trace --json` bzw.
`--trace --yaml` kommt dieselbe Matrix maschinenlesbar (`requirements` mit
`id`/`title`/`adrs`/`slices`, konditional `coverage` und `modality`, sowie
`orphan`; auf Matrix-Ebene `total`/`orphans` und — bei aktivem
`trace.cross-consistency` **mit** Differenzen — `crossConsistency` mit
`requirement`/`artifact`/`direction`/`file`/`line`). Standardmäßig ist `--trace`
**advisory** (Exit 0, auch bei Waisen). Für ein **Gate** ergänzen Sie
`--require-complete`: Ohne aktive `modality` gatet jede Waise; mit aktiver
`modality` gaten nur Waisen der in `require-levels` gelisteten Stufen. Mindestens
eine so gatende Waise ergibt Exit 1 (die Matrix bleibt unverändert auf stdout).
Ohne `--trace` ist `--require-complete` ein Nutzungsfehler (Exit 2).

**Andere Repo-Konvention / tabellarisches Lastenheft.** Folgt Ihr Repo einer
anderen Kennungs-/Dateikonvention oder pflegt es ein tabellarisches Lastenheft,
passen Sie die Quell-Achsen (`trace.requirements`, `format: table`, `table.*`,
je Referenzklasse `adrs`/`slices`) in `.d-check.yml` an. Alle Felder,
Migrationswege und Fehlerbilder stehen in §5.

### 4.13 Abdeckung aus einer kuratierten Matrix mitzählen (`trace.coverage`)

**Ausgangslage:** Die Abdeckung mancher Anforderungen liegt nicht in ADRs oder
Slices, sondern in einer **kuratierten Matrix** (z. B. einer ausgelagerten
Traceability-Datei). Der Referenz-Scan kennt sie nicht und zählt diese
Anforderungen fälschlich als Waisen.
**Ziel:** die kuratierte Deckung als eigene Dimension mitzählen, damit nur
wirklich ungedeckte Anforderungen als Waise erscheinen.

**Vorgehen:** den opt-in `trace.coverage`-Block auf die kuratierte(n) Datei(en)
richten — er liest sie als **eigene Coverage-Spalte** ein (range-aware,
abschnitts-gescopt):

```yaml
trace:
  requirements:
    id-pattern: 'GG-[A-Z][A-Z0-9]*-\d{3}'
  coverage:
    - files: [docs/plan/traceability.md]   # EXPLIZITE Dateien (keine dir/pattern)
      label: Trace                          # Owner-Kennung in der Coverage-Spalte
      ranges: true                          # AAA..BBB + AAA/BBB expandieren
      exclude-sections: ["27.1.1 Anforderungen ohne Design-Artefakt"]  # voller Heading-Text!
```

**Ergebnis:** eine zusätzliche **Coverage**-Spalte in der RTM. Eine Anforderung
ist jetzt nur noch **Waise**, wenn sie **weder** ein Slice **noch** eine
Coverage-Quelle deckt; `--require-complete` gatet entsprechend. Ohne
`trace.coverage` bleibt die RTM byte-identisch. Die Felder, die zugesagten
Range-/Aufzählungs-Notationen und die Fehlerbilder (fehlende `files`-Datei,
leeres `label`, Sektion ohne Treffer, ungültige Range) stehen in §5.

### 4.14 Nur MUSS-Anforderungen als Gate erzwingen (`trace.requirements.modality`)

**Ausgangslage:** Ihre Anforderungen tragen RFC-2119-Modalität (MUSS/SOLLTE/KANN)
im Text, aber nicht jede unerfüllte Anforderung soll ein Gate brechen — eine
offene KANN-Anforderung ist kein Release-Blocker.
**Ziel:** nur Waisen der Stufe **MUSS** zum Fehler machen, SOLLTE/KANN advisory
lassen.

**Vorgehen:** den opt-in `modality`-Block aktivieren — er klassifiziert jede
Anforderung (DE+EN-Modalverb-Defaults) und zeigt die Stufe in einer eigenen
**Modality**-Spalte:

```yaml
trace:
  requirements:
    id-pattern: 'GG-[A-Z][A-Z0-9]*-\d{3}'
    modality: {}          # {} = Built-in DE+EN-Defaults; require-levels [must]
```

**Ergebnis:** eine eigene **Modality**-Spalte, und `--require-complete` bricht
**nur** bei Waisen der `require-levels`-Stufen (Default `[must]`) —
SOLLTE/KANN/`unknown` bleiben advisory. Ohne `modality` bleibt die RTM
byte-identisch. **Wichtig:** ein `unknown` unter `require-levels: [must]` gatet
**nicht** — ein echtes MUSS mit unaufgeführtem Verb entkäme so dem Gate; prüfen
Sie die Spalte und ergänzen ggf. `levels` oder setzen strikt
`require-levels: [must, unknown]`. Die Klassifikations-Mechanik (Trefferwahl,
Body- vs. Tabellenspalten-Quelle), die Keyword-Konfiguration und die weiteren
Fail-closed-Fälle stehen in §5.

### 4.15 Prüfen, ob Ihre RTM-Tabelle noch zum Design passt (`trace.cross-consistency`)

**Ausgangslage:** Sie pflegen Anforderung→Design an **zwei** Orten — einer
kuratierten RTM-Tabelle (je Anforderung die Design-Artefakte) und am Design
selbst (je Artefakt eine `Bezug`-Spalte mit seinen Anforderungen). Beide sollen
dasselbe sagen. Nach ein paar Monaten sagen sie es nicht mehr, und niemandem
fällt es auf.
**Ziel:** die Stellen finden, an denen die beiden auseinanderlaufen.

**Schritt 1 — die beiden Sichten benennen.** Sagen Sie d-check, wo sie stehen und
woran es Anforderungen und Artefakte erkennt:

```yaml
trace:
  requirements:
    id-pattern: 'GG-[A-Z]+-\d{3}'
  cross-consistency:
    forward:
      file: docs/traceability.md
      sections: ["27.1 Anforderung zu Design"]
      req-column: Anforderung
      design-column: Design-Artefakte
      design-pattern: 'GG-AR-[A-Z0-9-]+'
    backward:
      file: spec/architecture.md
      edge-column: Bezug
      req-pattern: 'GG-[A-Z]+-\d{3}'
    exclude-req: '^GG-SPEC-'
```

Die `Bezug`-Spalte allein macht eine Tabelle zur Rück-Sicht; die Artefakt-ID
nimmt d-check aus deren **erster** Spalte, weil die je Tabelle anders heißt
(`Komponente`, `Port-ID`, `Kennung`). `exclude-req` blendet Anforderungen aus,
auf die Rück-Kanten zeigen, die aber keine eigene RTM-Zeile haben (etwa eine
Spezifikations-Zwischenschicht).

**Schritt 2 — laufen lassen.** Als Gate mit `--require-complete`, sonst nur
`--trace` (dann meldet der Lauf, ändert aber den Exit-Code nicht):

```text
$ docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0 \
    --trace --require-complete
…
## Kreuzverweis-Konsistenz

| Anforderung | Artefakt | Richtung | Quelle |
|---|---|---|---|
| GG-ARCH-006 | GG-AR-COMP-CORE | in RTM, ohne Rück-Kante | docs/traceability.md:7 |
| GG-ARCH-006 | GG-AR-COMP-SCHED | Rück-Kante, ohne RTM-Eintrag | spec/architecture.md:7 |

2 Differenz(en).
d-check: 2 Kreuzverweis-Differenz(en) zwischen Vorwärts- und Rück-Sicht (--require-complete)
```

**Schritt 3 — die Befunde abarbeiten.** Jede Zeile nennt die Datei und Zeile, an
der Sie ansetzen, und die Richtung sagt Ihnen, **welche** Seite Sie anfassen:

- *„in RTM, ohne Rück-Kante"* — Ihre Tabelle behauptet ein Artefakt, das die
  Anforderung nicht (mehr) nennt. Meist ist die Tabelle veraltet.
- *„Rück-Kante, ohne RTM-Eintrag"* — das Design nennt die Anforderung, Ihre
  Tabelle kennt das Artefakt nicht. Meist fehlt die Tabellen-Zeile.

Im Zweifel gewinnt die **Rück-Kante**: sie steht dort, wo das Design gepflegt
wird, und wird beim Ändern mitgeführt. Die RTM-Tabelle ist der Spiegel — sie
driftet. Die `GG-SPEC-042`-Rück-Kante taucht übrigens nicht auf: sie fällt per <!-- d-check:ignore -->
`exclude-req` heraus.

Die RTM oberhalb bleibt unverändert; der Abgleich ist eine **eigene** Ausgabe,
keine zusätzliche Spalte.

**Wenn Ihre Tabelle noch Prosa enthält.** Steht dort „alle Scheduler-Komponenten
(siehe Architektur)" statt konkreter IDs, meldet der Lauf schlicht *alle*
Rück-Kanten als „ohne RTM-Eintrag". Das ist kein Fehler, sondern Ihre
Arbeitsliste: Sie können die Tabelle daran entlang auf konkrete IDs umbauen.
Die vollständige Feld- und Fehlerliste steht in §5.

### 4.16 Ein Makefile-Fragment einbinden (`--print-mk`)

**Ziel:** d-check als `doc-check`-Schritt ins eigene Makefile einbinden, ohne
ein Recipe oder Skript zu kopieren — der Image-Pin bleibt bei d-check.
**Voraussetzung:** ein Projekt mit Makefile; eine eigene `.d-check.yml`.

**Vorgehen** (Fragment erzeugen, einbinden):

```bash
docker run --rm ghcr.io/pt9912/d-check:v0.67.0 --print-mk > d-check.mk
# im eigenen Makefile:  include d-check.mk
```

**Ergebnis:** ein include-bares `d-check.mk` auf stdout — eine überschreibbare
`DCHECK_IMAGE`-Variable (auf die ausgelieferte Release-Version gepinnt), die
Komfort-Variable `DCHECK_DIGEST` (sticht den Tag), `TRACE_FLAGS` und elf
`##`-annotierte Targets (`doc-check`, `doc-trace`, `doc-complete`, `doc-doctor`,
`doc-repair`, `doc-immutable`, `doc-commits`, `doc-planning`, `doc-tracked`,
`doc-targets`, `doc-help`):

<!-- d-check-test:not-replayable: abgekürzte Illustration (Elision mit # …), nicht die wörtliche --print-mk-Ausgabe -->
```text
DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v0.67.0
DCHECK_DIGEST ?=
TRACE_FLAGS ?=

ifeq ($(strip $(DCHECK_DIGEST)),)
DCHECK_REF := $(DCHECK_IMAGE)
else
DCHECK_REF := ghcr.io/pt9912/d-check@$(DCHECK_DIGEST)
endif

.PHONY: doc-check
doc-check: ## Doku-Referenzen prüfen (Befund-Gate)
	docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF)

.PHONY: doc-doctor
doc-doctor: ## erklärende Diagnose mit Fix-Kandidaten
	docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --doctor

.PHONY: doc-repair
doc-repair: ## Reparatur-Patch (unified diff) auf stdout, git-apply-rein
	@docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --repair

# … doc-trace, doc-complete (RTM/Gate, --trace[ --require-complete] $(TRACE_FLAGS))
# … doc-help (listet die doc-*-Targets über die ##-Annotationen)
```

**Hinweise:** `--print-mk` liest das Repository **nicht** und schreibt selbst
nichts (die Umleitung `>` legt die Datei an). Der eingebettete Image-Ref ist die
Version des aufrufenden d-check; für strikte Reproduzierbarkeit setzen Sie den
Digest aus den Release-Notes — `make doc-check DCHECK_DIGEST=sha256:<digest>`
(sticht den Tag) oder überschreiben den vollen `DCHECK_IMAGE`. `make doc-trace`
gibt die RTM advisory aus, `make doc-complete` als Gate (Waise ⇒ Exit 1); `make
doc-doctor` zeigt eine erklärende Diagnose, `make doc-repair > fix.patch` einen
`git apply`-reinen Reparatur-Patch, `make doc-help` listet die `doc-*`-Targets.
Mit `TRACE_FLAGS=--json` werden die RTM-Targets maschinenlesbar.

### 4.17 Closure-Notizen auf Substanz prüfen (Modul `planning`)

**Ausgangslage:** Ihr Prozess verlangt zu jedem abgeschlossenen Arbeitspaket eine
Closure-Notiz. Niemand merkt, wenn eine davon ein zurückgelassenes
„_Ausstehend._" ist — der Audit-Record hat dann eine Lücke, die erst auffällt,
wenn jemand ihn braucht.

**Ziel:** ein Gate, das die **Struktur** dieser Notizen prüft: Abschnitt
vorhanden, genug Substanz, keine bekannte Floskel.
**Voraussetzung:** ein Verzeichnis mit abgeschlossenen Arbeitspaketen (Konvention
`docs/plan/planning/done/`) und eine Roadmap-Datei für das Modul `planning`.

**Vorgehen** (Prüf-Profil anlegen, Gate fahren):

```yaml
# .d-check.closure.yml
modules: []
planning:
  roadmap: docs/plan/planning/in-progress/roadmap.md
  closure:
    dir: docs/plan/planning/done
    boilerplate: ["wie geplant umgesetzt", "war ganz okay"]
```

```bash
docker run --rm --network none -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.67.0 \
  --config .d-check.closure.yml --enable planning
```

**Ergebnis:** je abgeschlossenem Paket wird der erste Abschnitt geprüft, dessen
Überschrift auf `heading-pattern` passt (Default `^#{2,3} .*Closure-Notiz`) —
bis zur nächsten gleich- oder höherrangigen Überschrift:

<!-- d-check-test:not-replayable: illustrativer Befund an einem erfundenen Paket, nicht die wörtliche Ausgabe dieses Repos -->
```text
d-check: 301 Datei(en) geprüft, 1 Befund(e)
docs/plan/planning/done/slice-042-x.md:3	docs/plan/planning/done	closure-note-thin
```

Drei Grund-Codes, weil sie drei verschiedene Reparaturen verlangen:
`closure-note-missing` (Abschnitt schreiben), `closure-note-thin` (Substanz
ergänzen — Satzende-Zeichen **außerhalb** von Code-Blöcken, Schwelle
`min-sentences`, Default 4) und `closure-note-boilerplate` (Floskel ersetzen).

**Was als Satz und als Floskel zählt** (seit v0.56.0): Der Abschnitt wird
**einmal** bereinigt — Code-Blöcke **und** Inline-Code fallen weg —, und alle
Bedingungen lesen diesen einen Text. Ein Satzende zählt nur, wenn ihm ein
Leerzeichen oder das Zeilenende folgt; die Punkte in `0.56.0` und in einem
Link-Pfad sind damit keine Sätze mehr. Floskeln werden an **Wortgrenzen**
verglichen: `ok` trifft das eigenständige Wort, nicht *dokumentiert*. Damit sind
kurze Phrasen überhaupt erst benutzbar.

**Zwei Nebenwirkungen, die Sie kennen sollten.** Eine Floskel **in Backticks**
trifft nicht mehr — eine zitierte Floskel ist keine benutzte, aber die Regel
reicht weiter: geleert wird jeder Inline-Code-Span, auch ein unbeabsichtigter
aus zwei einzelnen Backticks im selben Absatz. Und als Wortzeichen gilt die
ASCII-Menge; ein Umlaut ist keines, eine Phrase mit angrenzendem Umlaut gilt
also als grenzständig und trifft.

Eine **vierte** Bedingung ist per Default aus und wird mit
`placeholder: true` zugeschaltet: `closure-note-placeholder` meldet den
**unausgefüllten Rumpf** einer Vorlage. Sie schließt eine Lücke, die die anderen
drei nicht sehen — `Ergebnis: <ergebnis>. Belege: <belege>.` ist syntaktisch
vollständig, erreicht die Schwelle und enthält keine Floskel.

Erkannt wird die Auszeichnungs-Form: eine öffnende Winkelklammer ohne
vorausgehendes Wortzeichen, deren Inneres **kein Whitespace** enthält — ein
Platzhalter ist ein Feldname, kein Satz. Damit bleibt technische Prosa grün:
Vergleichszeichen (`p95 < 1 s` ebenso wie die enge Form), Generics,
Winkelklammer-Linkziele, Autolinks, Adressen und HTML-Tags. **Inline-Code zählt
nicht** — dort wird Syntax gezeigt (`` `<PREFIX>` ``), und genau dort steht sie
absichtlich. Gemeldet wird der **erste** Treffer je Notiz, an seiner Zeile.

**Zwei Grenzen sollten Sie kennen:** ein mit vier Leerzeichen **eingerückter**
Code-Block ist in d-check nirgends modelliert — ein Platzhalter darin meldet.
Und eine ungerade Backtick-Zahl im Absatz verschiebt die Inline-Code-Paarung,
hier wie überall. Beides sind Eigenschaften der Markdown-Vorverarbeitung.

**Hinweise:** Das eigene Profil ist kein Zufall — der `closure`-Block gehört
**nicht** in Ihre reguläre `.d-check.yml`, sonst prüft jeder gewöhnliche Lauf
Ihre Closure-Notizen mit. Eine Closure-Frage gehört an den Abschluss, nicht in
den inneren Loop; `--config` (§5) trennt beide Profile.

Das Gate ist **fail-closed**: ein gesetztes, aber fehlendes Verzeichnis meldet
ebenso wie eines **ohne einen einzigen** Kandidaten — den Schlüssel zu setzen ist
die Behauptung, dass dort Notizen liegen, und ein leer laufendes Gate darf nicht
Erfolg melden. Ohne `closure.dir` ist die Fähigkeit inert und öffnet keine Datei.

Und es prüft **Struktur, nicht Bedeutung**: „war ganz okay, läuft jetzt" ist
syntaktisch ein vollständiger Satz. Grün heißt „Form erfüllt", nicht „Notizen
sind gut" — die inhaltliche Beurteilung bleibt einem Review vorbehalten.

### 4.18 Eine Tabellenspalte kurz und gefüllt halten (Modul `structure`)

**Ziel:** Ein Register-Dokument — ein ADR-Index, eine Modul-Tabelle, ein
Release-Register — hat eine Spalte, die *kurz* sein soll: ein Titel, ein
Status. Genau dort wächst der Text, weil niemand nachmisst. Aus einer Sicht,
die auf andere Dokumente **zeigen** soll, wird so eine zweite Quelle, die von
der ersten abdriftet.

**So geht es:**

```yaml
structure:
  - files: docs/plan/adr/README.md
    section: "# ADR-Index"
    table:
      column:                       # beliebig viele Spalten, EIN Selektor
        - name: Titel               # die Spalte über ihren KOPFZEILEN-Namen
          cell-max-chars: 200       # Obergrenze in Zeichen
          cell-min-chars: 10        # Untergrenze — sonst darf sie leer sein
        - name: Datum
          cell-min-chars: 10        # ISO-Datum, exakt zehn Zeichen
          cell-max-chars: 10
```

```bash
docker run --rm --network none -v "$PWD":/repo:ro \
  ghcr.io/pt9912/d-check:latest --enable structure
```

**Ergebnis:** je zu langer oder zu kurzer Zelle ein Befund — auf **ihrer**
Zeile, nicht am Abschnittskopf:

<!-- d-check-test:not-replayable: gekürzte Illustration (Ziel-Spalte elidiert, Erläuterung gekürzt), nicht die wörtliche Ausgabe dieses Repos -->
```text
docs/plan/adr/README.md:41  …	section-cell-oversized	Zelle der Spalte "Titel" hat …
docs/plan/adr/README.md:64  …	section-column-missing	keine Tabelle des Abschnitts …
```

**Vier Dinge, die dabei zählen:**

- **Der Spaltenname, nicht die Position.** Wer `name: 2` schreiben
  könnte, hätte ein Gate, das nach dem Einfügen einer Spalte **still** die
  falsche misst. Über den Kopfzeilen-Namen fällt derselbe Umbau **laut** auf.
- **Ein Selektor, beliebig viele Spalten.** Jede weitere Spalte ist ein
  Listen-Eintrag, keine zweite Regel mit wiederholtem `files`/`section` — und
  dieselbe Spalte zweimal zu nennen ist ein Config-Fehler, kein Doppel-Gate.
- **Ohne `cell-min-chars` darf die Zelle leer sein** — null Zeichen liegen
  unter jeder Obergrenze. Wer „gefüllt" meint, muss es sagen.
- **Eine Zeile, die gar nicht bis zur Spalte reicht, meldet.** Das ist der Fall,
  den man sonst nie sieht: die vorderen Spalten stehen ja da. Im ADR-Index
  dieses Repos hat genau das **vier** Zeilen ohne Datum und ohne Bezug
  aufgedeckt.

**Mehrere Spalten** begrenzen Sie mit mehreren Regeln über denselben Abschnitt
— das geht, weil die Regel-Identität die benannte Spalte mitträgt.

**Eine Grenze:** gemessen wird die Zelle, **wie sie dasteht**, Markdown-Syntax
eingeschlossen. Auf eine Spalte voller Links angewandt misst die Bedingung
etwas anderes, als Sie meinen.

### 4.19 Workflow-Referenzen gepinnt und rechte-gedeckt halten (Modul `workflows`)

**Ausgangslage:** Ihre CI-Workflows rufen fremde Actions und eigene, lokale
Workflows auf. Beim Fremden ist ein beweglicher Tag eine Supply-Chain-Fläche —
er lässt sich umhängen. Beim Eigenen ist die Falle eine andere: ein
aufgerufener Workflow bekommt nur die Rechte, die der aufrufende **Job** selbst
führt. Verlangt er mehr, bricht der Lauf **vor dem ersten Job** ab — ohne Log,
ohne Job, ohne Hinweis darauf, welche Zeile schuld ist.

**Ziel:** Beide Zusagen als Gate, bevor ein Tag-Push sie prüft.

**So geht es:**

```yaml
modules: [workflows]
workflows:
  dir: .github/workflows        # Aktivierungs-Schalter; ohne ihn ist das Modul inert
  # exempt-paths: ["**/experimental-*.yml"]
```

```bash
docker run --rm --network none -v "$PWD":/repo:ro \
  ghcr.io/pt9912/d-check:latest --enable workflows
```

**Ergebnis:** je Referenz ein Befund auf **ihrer** Zeile:

<!-- d-check-test:not-replayable: gekürzte Illustration, nicht die wörtliche Ausgabe dieses Repos -->
```text
.github/workflows/ci.yml:31       …	uses-pin-missing
.github/workflows/release.yml:268 …	uses-local-perms-undeclared
```

**Vier Dinge, die dabei zählen:**

- **Die lokale Referenz braucht keinen Pin — und trotzdem eine Prüfung.** Sie
  löst auf denselben Commit auf wie ihr Aufrufer und ist damit *stärker*
  gebunden als ein SHA. An die Stelle des Pin-Checks treten die zwei Fragen,
  die dort Sinn ergeben: existiert das Ziel, und bekommt es seine Rechte?
- **Ein Job ohne eigenes `permissions:` ist der häufige Fall.** Er erbt den
  Workflow-Kopf — und wenn der `permissions: {}` führt, gibt er nichts weiter.
  Das meldet `uses-local-perms-undeclared`, **eine** Meldung je Referenz, weil
  die Reparatur eine ist: dem Job die Rechte deklarieren.
- **Ein Scope, den der Aufrufer nicht nennt, ist `none`** — nicht „egal".
  `read-all` und `write-all` setzen dagegen jeden Scope auf ihre Stufe.
- **Der Ort ist konfigurierbar, nicht verdrahtet.** `.github/workflows` ist
  eine GitHub-Konvention; ein festkodierter Pfad wäre für jede andere Ablage
  wertlos.

**Was das Modul nicht sagt:** ob der Workflow läuft. Es deckt **eine**
Deklarations-Klasse — ein grüner Lauf heißt „diese Klasse liegt nicht vor".

## 5. Konfiguration

Eine optionale `.d-check.yml` in der Repo-Wurzel passt d-check an. Ohne
Datei gelten die Standardwerte (Module `links` und `anchors`, Scan-Wurzeln
`docs/` und `spec/` plus Wurzel-`*.md`). Das vollständige, kommentierte
Schema liefert `--print-config` (siehe [Abschnitt 4.3](#4-aufgaben)).

**Eine andere Datei verwenden (`--config`).** `--config <datei>` **ersetzt** für
diesen Lauf die konventionelle `.d-check.yml` — der Pfad ist relativ zur
Scan-Wurzel und muss **innerhalb** von ihr liegen (im read-only-Mount ist alles
andere ohnehin unerreichbar). So fährt ein Repo mehrere **Prüf-Profile**, etwa
ein schnelles für den inneren Loop und ein engeres für den Abschluss
(Aufgabe §4.17). Die Datei ersetzt, sie ergänzt nicht: eine vorhandene
`.d-check.yml` wird dann **nicht** gelesen und nicht zusammengeführt.

Es gibt **keinen** stillen Rückfall. Fehlt die Datei, zeigt der Pfad auf ein
Verzeichnis, führt er über einen Symlink aus der Wurzel heraus oder ist der Wert
**leer** (klassisch: eine nicht expandierte Make-Variable), bricht der Lauf mit
Exit 2 ab — statt unbemerkt mit einem anderen Prüfumfang grün zu werden.

### Scan-Bereich

```yaml
scan:
  roots: ["docs", "spec"]      # zu prüfende Wurzeln
  ignore: ["docs/archiv/**"]   # zusätzliche Ignorier-Muster (Globs)
```

**Glob-Syntax.** Pro Pfad-Segment gilt Go-`path.Match` — `*` und `?` matchen
**nicht** über `/` hinweg; `**` steht für beliebig viele Segmente. Eine
**negierte** Zeichenklasse schreibt sich `[^…]` (Go), **nicht** `[!…]`
(Shell/fnmatch): `[!a]` matcht in d-check die Literale `!` und `a`, nicht
„alles außer `a`". Das gilt für alle Glob-Felder — `scan.ignore`,
`matrix.classes[].paths`/`.order` und die `exempt-paths` der Module.

### Module wählen

```yaml
modules: [links, anchors, ids, matrix]  # Netz-Module external/sources separat zuschalten
```

### Kennungs-Muster (Modul `ids`)

```yaml
ids:
  patterns:
    - regex: 'ADR-\d{4}'
      target: docs/plan/adr/   # wo die Kennung definiert ist
      link-policy: always      # auch Inline-Code-Vorkommen linkpflichtig
      exempt-paths: [CHANGELOG.md]
```

### Referenzmatrix (Modul `matrix`)

```yaml
matrix:
  classes:
    - name: spec
      paths: [spec/lastenheft.md]
    - name: adr
      paths: ["docs/plan/adr/[0-9]*.md"]
  rules:
    - {from: spec, to: adr, allow: false}
  status:
    forbidden: [superseded, deprecated]
  exclude-sections: [Historie]
```

### Weitere Module

```yaml
codepaths:
  roots: [docs, spec]          # Pfade in Inline-Code prüfen
  exempt-paths: ["docs/reviews/**"] # Dateien ganz ausnehmen (Glob, wie ids)
  ignore-refs: ["tools/altes-skript.sh"] # Alias des geteilten ignore-refs (nur codepaths)
hostpaths:
  prefixes: [home, Users]       # Verzeichnisnamen OHNE / (führender /… ⇒ Exit 2)
external:
  timeout-seconds: 10
diagrams:                      # Kennungen in Diagramm-Fences prüfen
  fences: [mermaid]            # zu öffnende Diagramm-Sprachen (Default mermaid)
  exempt-paths: ["docs/reviews/**"] # Dateien ganz ausnehmen (Glob, wie ids)
  patterns:
    - regex: 'ARC-\d{2}'
      defined-in: spec/architecture.md  # Token muss hier (außerhalb Fences) vorkommen
versions:                      # gepinnte Versions-Verweise gegen die aktuelle Version
  pin-pattern: 'ghcr\.io/[^\s:]+:(v\d+\.\d+\.\d+)'  # Version in Capture-Gruppe 1
  current-from: version.md#aktuell  # Datei#Anker (Span) mit der aktuellen Version
  exempt-paths: [CHANGELOG.md] # historische Pins ausnehmen (Glob)
  # ODER — mehrere Muster-Quellen-Paare statt eines; NICHT zusammen mit der
  # Kurzform oben (beides gesetzt ⇒ Exit 2), exempt-paths gilt je Paar:
  # patterns:
  #   - pin-pattern: 'ghcr\.io/[^\s:]+:(v\d+\.\d+\.\d+)'
  #     current-from: version.md#aktuell
  #     exempt-paths: [CHANGELOG.md]
  #   - pin-pattern: 'ai-harness-course/tree/(v\d+\.\d+\.\d+)'
  #     current-from: harness/conventions.md#baseline
vcs:                           # git-Diff-Immutabilität (opt-in; braucht .git + Range)
  paths: ["docs/plan/adr/[0-9]*.md"]    # geschützte Datei-Klasse (Glob)
  immutable-when: '^\*\*Status:\*\* Accepted'        # BASE immutabel ab dieser Zeile (nur ausserhalb von Code-Bloecken)
  exclude-sections: [Geschichte]        # Abschnitte, die nicht zum Core zählen
  status-line: '^\*\*Status:\*\*'       # Kopf-Status-Zeile (aus dem Core gestrippt)
  head-allow: '^\*\*Status:\*\* (Accepted|Superseded by ADR-[0-9]{4})'  # erlaubter Status-Übergang
commits:                       # Traceability-Kennung in Commit-Messages (opt-in; git)
  id-patterns:                 # gültige Kennungen; eine Message ohne Match ⇒ commit-untraceable
    - 'ADR-\d{4}'
    - 'slice-\d+'
  exempt-pattern: '^(Merge |Revert )'   # kennungsfreier Betreff (Merge/Revert)
tracked:                       # Getrackt-Status von Link-Zielen (opt-in; git-Index, ohne Range)
  exempt-targets: ["build/**"] # absichtlich untrackte Ziele (Glob über den aufgelösten Pfad)
targets:                       # Deklarations-Konsistenz Doku ↔ Build-Targets (opt-in; hermetisch, kein git)
  makefiles: [Makefile]        # Regelnamen-Quelle(n)
  doc-tables: [AGENTS.md]      # make-X-Tabellen (Richtung 1 ⇒ gate-phantom)
  authority: AGENTS.md         # Vollständigkeits-Quelle (Richtung 2 ⇒ gate-undocumented)
  exempt-targets: []           # Regelnamen EXAKT (kein Glob, anders als tracked) — Utility ohne Doku-Pflicht
trace:                         # konfigurierbare RTM-Quellen (KEIN Modul; steuert nur --trace)
  requirements:
    source: spec/lastenheft.md              # nichtleer explizit: fehlend/0 Treffer => Exit 2
    id-pattern: 'GG-[A-Z][A-Z0-9]*-\d{3}'  # Anforderungs-Kennung (Default: <PREFIX>-FA-…/-QA-…)
    format: headings                         # headings (Default) oder table
    # table:                                 # Pflicht bei format: table; exakte Header-Namen
    #   id-column: Kennung
    #   text-column: Anforderung
    #   # text-columns: [Anforderung, Akzeptanzkriterium] # alternativ zu text-column
    #   modality-column: Prioritaet          # optional; sonst Textspalte
    #   duplicate-ids: error                 # error (Default), first oder last
    modality: {}                          # RFC-2119-Stufe je Anforderung (eigene Spalte); {} = DE+EN-Defaults, require-levels [must]
  slices:
    file-pattern: '^(\d+)-.*\.md$'          # Slice-Dateiname NNN-titel.md (Capture 1 = Owner-Kennung)
  coverage:                                 # kuratierte Coverage-Quellen (eigene RTM-Spalte)
    - files: [docs/plan/traceability.md]    # EXPLIZITE Dateien (keine dir/pattern → keine ADR-Kontamination)
      label: Trace                          # Owner-Kennung in der Coverage-Spalte
      ranges: true                          # GG-QA-001..006 → alle sechs; gegen id-pattern validiert
      # sections: ["27.1 Anforderung zu Design"]                      # Whitelist (voller Heading-Text)
      # exclude-sections: ["27.1.1 Anforderungen ohne Design-Artefakt"]  # Blacklist (voller Heading-Text)
```

d-check kennt **vier** Ventil-Achsen, um eine Referenz von der Prüfung auszunehmen.
Sie greifen an unterschiedlichen Stellen und ergänzen sich — vom gröbsten zum
feinsten Schnitt:

- **`scan.ignore`** (Quelle, ganze Datei): eine Datei wird gar nicht erst gescannt,
  ihre Referenzen erreichen **kein** Modul. Der gröbste Schnitt.
- **`exempt-paths`** (ganze Datei, je Modul): nimmt **ganze Dateien** von der Prüfung
  **eines** Moduls aus — `ids` (je Muster), `versions` (je Muster-Quellen-Paar),
  `structure` (je Regel), `matrix`, `codepaths` und `diagrams` (jeweils modulweit);
  der Rest wird weiter gescannt. Die übrigen Module haben den Schlüssel nicht.
- **`d-check:ignore`** (eine Zeile): der Marker nimmt **eine einzelne Zeile** aus —
  in `codepaths`, `ids`, `versions` und `diagrams`. Bei `ids` gilt er für nackte
  wie für Inline-Code-Vorkommen, unabhängig von der Link-Politik; bei `versions`
  nimmt er die Zeile **allen** Muster-Quellen-Paaren aus, während `exempt-paths`
  dort paar-lokal wirkt. **Die übrigen Module kennen ihn nicht** — ein
  `matrix`-Befund etwa wird behoben oder **strukturell** ausgenommen
  (`exclude-sections`, `allow-supersede-lineage`), nicht zeilenweise
  stummgeschaltet. Diese vier sind eine **benannte Liste, kein ableitbares
  Kriterium**: `matrix`, `structure` und `citations` melden ebenfalls auf
  Zeilen und tragen den Marker nicht — `citations` sogar keine der drei
  **feinen** Achsen (`exempt-paths`, `ignore-refs`, Marker), obwohl es
  dieselben Grund-Codes meldet wie das ventil-tragende `codepaths`. Die
  **groben** Achsen (`scan.ignore`, modul-lokaler `scope`) gelten dagegen für
  jedes Modul.
- **`ignore-refs`** (Ziel, querschnittlich): nimmt bestimmte **aufgelöste
  Ziel-Pfade** von der Existenz-/Anker-Prüfung aus — **referenz-weit** (datei- und
  zeilen-unabhängig) und geteilt von `links`, `anchors` und `codepaths`.

  **Wo der Marker wirkt — und wo er nur erwähnt ist.** Bei `codepaths` und `ids`
  zählt er nur **außerhalb** von Code-Blöcken **und außerhalb von Inline-Code**:
  in Backticks ist er eine **Erwähnung**, keine Direktive. Ohne das könnte diese
  Seite hier ihn nicht beschreiben, ohne sich selbst auszunehmen — genau das war
  im Lastenheft an fünf Stellen passiert. Bei `versions` und `diagrams` wird er
  weiterhin **roh** erkannt — aus zwei **verschiedenen** Gründen. Bei `diagrams`
  ist es strukturell: es liest die Zeilen **innerhalb** eines Fence und die
  Öffnungszeile, und dort ist ein Backtick literaler Inhalt. Bei `versions` ist
  es eine **benannte Grenze**: es liest **alle** Zeilen, also auch Prosa-Zeilen,
  und auf einer solchen antwortet `d-check` damit **zweifach** — `codepaths`
  meldet dort, wo `versions` schweigt. Wer eine Zeile mit einer Erwähnung in
  Backticks schreibt, muss das wissen.
  **Zwei Ränder, die aus der Absatz-Erkennung folgen:** eine Backtick-Spanne
  desselben Absatzes verschluckt auch einen **gesetzten** Marker (er wirkt dann
  nicht — Falsch-Rot), und ein **unpaariger** Backtick weiter oben kippt die
  Parität, sodass eine Erwähnung doch als Direktive wirkt (stilles Grün). Die
  Regel oben gilt also nur bei gerader Backtick-Parität des Absatzes.
  **Auch die Form folgt der Eingabe.** Bei `codepaths` und `ids` muss der
  Marker in einem **HTML-Kommentar** stehen — `<!-- d-check:ignore (Grund) -->`;
  eine **blanke** Erwähnung in Prosa wirkt nicht. Bei `versions` und `diagrams`
  ist er ein **Token**, weil dort kein Markdown-Kommentar die Lexik bildet — im
  `mermaid`-Fence versteckt ihn die Diagramm-Sprache (`%% d-check:ignore`).
  **Konservativ:** ein `>` im Kommentar vor dem Marker lässt ihn nicht gelten.
  **Der Preis:** wer den Marker setzt, muss wissen, für welches Modul.
**Im Modul `diagrams` ist der Marker ein Token, kein HTML-Kommentar** — und er
wirkt an genau **zwei** Orten. In Prosa schreiben Sie ihn als
`<!-- d-check:ignore -->`; innerhalb eines Diagramm-Fence ist das kein
Kommentar, sondern Diagramm-Text. Das Modul sucht deshalb die **Zeichenfolge**
selbst, unabhängig davon, was sie umgibt.

**Erster Ort — die Öffnungszeile des Fence.** Der Marker nimmt dort den
**ganzen** Block aus. Er steht im Infostring, ist also weder HTML- noch
Diagramm-Kommentar und braucht kein Kommentar-Zeichen:

```mermaid d-check:ignore (Beispiel-Diagramm, erfundene Kennungen)
flowchart LR
    A["ARC-42 Frontend"] --> B["ARC-43 Backend"]
```

**Zweiter Ort — eine einzelne Diagramm-Zeile.** Dort nimmt der Marker nur diese
eine Zeile aus, und dort ist er echter Diagramm-Inhalt: Sie müssen ihn vor dem
Renderer verstecken, in Mermaid als `%% d-check:ignore`. Die übrigen Zeilen
desselben Blocks bleiben geprüft.

**Auf der schließenden Fence-Zeile wirkt er nicht.** Sie ist keiner der beiden
Orte; wer den Marker dort symmetrisch zur Öffnungszeile setzt, bekommt den
Befund weiterhin. Ohne Marker bleibt für `diagrams` nur
`diagrams.exempt-paths` (ganze Datei), der modul-lokale `diagrams.scope` oder
`scan.ignore` — Schnitte, die mehr wegnehmen, als gemeint war.

`ignore-refs` ist das **Ziel-Achsen-Pendant** zu `scan.ignore` (Quell-Achse): jenes
entfernt eine ganze Datei vom Scan, dieses ein einzelnes **Ziel** von der Prüfung —
gleich, in welcher Datei die Referenz steht. Es steht als **Top-Level-Liste**; jeder
Eintrag trägt drei Felder:

```yaml
ignore-refs:
  - in: "lab/templates/**"                    # nur Referenzen AUS diesen Dateien (optional; ohne in = repo-weit)
    refs: ["lab/templates/**"]                # diese aufgelösten Ziele ignorieren (Template-Platzhalter)
    keep: ["lab/templates/**/*.template.md"]  # aber diese echten Verweise scharf lassen
```

- **`refs`** (Pflicht): Globs auf den **aufgelösten** Ziel-Pfad. Leer/fehlend ⇒ der
  Eintrag ignoriert nichts.
- **`keep`** (optional): Ausnahmen — ein von `refs` getroffenes Ziel bleibt geprüft,
  wenn ein `keep`-Glob es trifft. `keep` gewinnt **reihenfolge-unabhängig**.
- **`in`** (optional): Quell-Skopus — ein Glob auf die **Quelldatei**. Gesetzt, gilt
  der Eintrag nur für Referenzen in passenden Dateien.

Zwei typische Fälle. Ein **Template-Verzeichnis**, dessen Platzhalter-Pfade erst im
Ziel-Repo auflösen: `in` skopiert auf die Templates, `refs` ignoriert die Platzhalter,
`keep` hält die echten (Kurs-/Template-internen) Verweise scharf — so wird das ganze
Verzeichnis geprüft, ohne die Platzhalter-Fehlalarme. Oder ein **Tombstone**: ein
bewusst entfernter Pfad, den immutable/historische Doku (z. B. eine akzeptierte ADR)
noch zitiert und der sonst als `codepath-missing` dangelte, ohne dass man die immutable
Doku editieren dürfte.

Bewusster Akt **mit Gate**: ohne passenden Eintrag meldet ein fehlendes Ziel weiter
(`target-missing`/`codepath-missing`/`anchor-missing`) — nichts verschwindet still.
Die Globs matchen den **aufgelösten, Wurzel-relativen** Pfad (bzw. die Quelldatei),
nicht die rohe `./`-/`../`-Schreibweise, die der Befund unter `Target` zeigt; tragen
Sie den Pfad also so ein, wie er von der Repo-Wurzel aus lautet. Die Symlink-Ablehnung
bleibt unberührt. Der frühere modul-lokale Schlüssel `codepaths.ignore-refs` bleibt als
**Alias** gültig (kein Config-Bruch): er wirkt wie ein `ignore-refs`-Eintrag ohne
`in`/`keep`, skopiert auf `codepaths`.

### Ortsfeste Verweise in Lifecycle-Verzeichnissen (`links.resolve-from`)

Wo Dateien per `git mv` zwischen Geschwister-Verzeichnissen **wandern** — ein
Planning-Lifecycle wie `open/ → in-progress/ → done/` ist der Anlassfall —, ist
ein relativer Nachbar-Verweis am Ist-Ort grün und bricht beim nächsten
Verzeichniswechsel. `links.resolve-from` meldet ihn **vor** dem Move: jede
Datei in einem wandernden Verzeichnis muss jedes relative Ziel von **jedem**
Ort ihrer Gruppe auflösen, und überall auf **dasselbe** Ziel.

```yaml
links:
  resolve-from:
    - dirs:                                # wandernde Geschwister-Orte (min. 2) — nur deren Dateien sind Quellen
        - docs/plan/planning/open
        - docs/plan/planning/in-progress
      fixed-dirs:                          # ortsfeste Ziele (z. B. Ruheort): zählen als Ort, prüfen aber nicht
        - docs/plan/planning/done
```

- **Quellen sind nur Dateien in `dirs`-Verzeichnissen** (exakter Verzeichnis-Treffer,
  keine Unterverzeichnisse). Dateien in `fixed-dirs` sind am Endzustand — ihre
  Verweise müssen nur vom Ist-Ort auflösen.
- **Befund ist `link-position-dependent`**, und die Reparatur ist das
  **Präfixieren des Pfads**, nicht das Anlegen des Ziels — am Ist-Ort ist nichts
  kaputt. Ein Ziel, das schon am Ist-Ort fehlt, meldet weiter nur
  `target-missing` (kein Doppelbefund).
- **Absichtlich ortsgebundene Verweise** nehmen Sie über das bestehende
  [`ignore-refs`](#weitere-module) aus — dasselbe Ventil, referenz-weit.
- **Fail-closed:** eine Gruppe, deren `dirs`-Orte **alle** fehlen, und ein Ort,
  der als **Datei** existiert, melden über denselben Code. Ein **einzelner**
  fehlender Ort meldet bewusst nicht — git überträgt leere Verzeichnisse nicht,
  ein legitim geleertes Lifecycle-Verzeichnis fehlt auf jedem frischen Klon.
- **Grenzen:** die Gruppen-Orte müssen im wirksamen Scan-Bereich liegen (eine
  nie gescannte Datei ist still keine Quelle), und wandert das **Ziel** statt
  der Quelle, meldet die Prüfung nichts. Config-Fehler (unter zwei `dirs`,
  absolute Pfade, ein Verzeichnis in zwei Gruppen) brechen mit Exit 2 ab.

Ohne den Block ist der Befundsatz byte-identisch.

### Zitate und Zeilen-Referenzen gegen ihre Quelle prüfen (`codepaths.check-lines` / Modul `citations`)

Zwei opt-in-Fähigkeiten prüfen `datei:zeile`-**Zitate**, die sonst still ins Leere
zeigen können (etwa nach einem Tag-Bump einer committet-vendored Baseline).

`codepaths.check-lines` (bool, Default `false`) verifiziert die **Zeilen-Referenz**
eines Inline-Code-Pfads: trägt `` `datei:<von>-<bis>` `` (oder `` `datei:<zeile>` ``)
einen Bereich, muss das Ziel existieren und mindestens `<bis>` Zeilen haben (sonst
`citation-out-of-range`), und `<von>` darf nicht größer als `<bis>` sein (sonst
`citation-inverted-range`). Ohne den Schalter wird das Zeilen-Suffix wie bisher
abgetrennt und verworfen — **byte-identisch** zum bisherigen Verhalten.

Das Modul `citations` (opt-in, `--enable citations`) prüft **wortgleiche Zitate**. Die
Direktive `<!-- d-check:cite <pfad>:<von>-<bis> -->` markiert das folgende Zitat — die
nächste nicht-leere Zeile als `>`-Blockquote **oder** den nächsten inline-Zitat-Span
(`„…"` bzw. `"…"`). Leerzeilen dazwischen sind unschädlich — gesucht wird der
nächste nicht-leere Kandidat. **Ein Code-Block dazwischen trennt dagegen:** dann
folgt der Direktive kein Zitat (fail-closed, Exit 2, siehe unten). Dasselbe gilt
innerhalb eines `>`-Blocks — ein Code-Block beendet ihn. Quell-Spanne und Zitattext werden
**whitespace-normalisiert** (jeder Umbruch-/Leerraum-Lauf zu einem Leerzeichen); das
Zitat muss ein zusammenhängender **Teilstring** der Quelle sein, sonst
`citation-mismatch`. So besteht ein re-wrapped, mitten in der Zeile beginnendes Zitat,
während jede echte Wort-Abweichung bricht. Zitate unter 16 Zeichen bleiben ungeprüft
(ein sehr kurzer Teilstring träfe zufällig). Ein Repo-Escape des Ziels ist ein Befund;
eine **malformte** Direktive oder ein **fehlendes** Zitat ist fail-closed (Exit 2) —
eine kaputte Direktive ist ein Autoren-Fehler, kein Schweigen.

**Wo die Direktive wirkt — und wo sie nur erwähnt ist.** Sie zählt nur
**außerhalb** von Code-Blöcken **und außerhalb von Inline-Code**: in Backticks
geschrieben ist sie eine **Erwähnung**, keine Direktive. Ohne das könnte diese
Seite hier ihre eigene Syntax nicht zeigen, ohne den Lauf abzubrechen. Trägt
eine Zeile beides, gilt die freie Direktive. **Drei Grenzen, die daraus folgen:**
die Erkennung wirkt **absatzweit**, eine freie Direktive wird also auch
verschluckt, wenn eine Backtick-Spanne des Absatzes sie umschließt (öffnet eine
Zeile vorher, schließt eine Zeile später); ein **Pfad** in Backticks fällt
fail-closed statt zu einem Befund mit leerem Ziel; und eine Backtick-Spanne
zwischen `<!--` und dem Marker verschwindet, kann eine Direktive also auch
**erzeugen**. Die **Ventil**-Direktive `d-check:ignore` folgt derselben Regel,
aber **nur bei zwei ihrer vier Konsumenten**: bei `codepaths` und `ids` zählt
sie seit v0.64.0 ebenfalls nur außerhalb von Inline-Code (und muss dort
zusätzlich in einem HTML-Kommentar stehen), bei `versions` und `diagrams` wird
sie weiterhin auf der **rohen** Zeile erkannt und wirkt auch aus Backticks
heraus. Wer den Marker setzt, muss also wissen, für welches Modul — der
Ventil-Überblick weiter oben führt die Unterscheidung aus.

**Die zitierte Quelle ist nicht Teil des Scans.** Gelesen wird sie **roh und
typunabhängig**: sie muss weder Markdown sein noch in der Scan-Menge liegen.
`scan.ignore` und `citations.scope` skopieren die **prüfende** Datei, nicht das
Ziel — eine ausgeschlossene Datei kann also weiterhin zitiert werden.

```yaml
codepaths:
  check-lines: true       # datei:<von>-<bis>-Zeilen-Referenzen verifizieren
# citations wird über --enable citations bzw. modules: [..., citations] aktiviert
```

**Achtung, dieselben Grund-Codes — verschiedene Ventil-Lage.**
`citation-out-of-range` und `citation-inverted-range` können aus **beiden**
Fähigkeiten kommen, und der Grund-Code sagt Ihnen nicht, aus welcher. Für die
Reparatur ist das der entscheidende Unterschied:

| Ventil-Achse | `codepaths.check-lines` | Modul `citations` |
|---|---|---|
| Auslöser | `` `datei:<von>-<bis>` `` in Inline-Code | Direktive `<!-- d-check:cite … -->` |
| `scan.ignore` (Quelldatei) | ja | ja |
| modul-lokaler `scope` | ja (`codepaths.scope`) | ja (`citations.scope`) |
| `exempt-paths` (Datei, modulweit) | ja (`codepaths.exempt-paths`) | **nein** |
| `ignore-refs` (Ziel, querschnittlich) | ja | **nein** |
| Zeilen-Marker `d-check:ignore` | ja | **nein** |

Die beiden **groben** Achsen wirken also für beide: `scan.ignore` nimmt die
Quelldatei aus dem Scan, `citations.scope` (Wurzeln plus `ignore`) skopiert das
Modul modul-lokal. Was `citations` **nicht** hat, sind die **feinen** Achsen —
`exempt-paths`, das geteilte `ignore-refs` und den Zeilen-Marker. Ein
`citation-out-of-range` aus einer `d-check:cite`-Direktive lässt sich deshalb
nicht **einzeln** stummschalten: entweder Sie beheben ihn (Bereich korrigieren,
Direktive entfernen) oder Sie nehmen die **ganze Datei** heraus. Wer den Marker
an eine `d-check:cite`-Zeile schreibt, ändert nichts; derselbe Marker an einer
Inline-Code-Zeilen-Referenz wirkt sehr wohl.

Ein Wort zur Erwartungshaltung: dieser Unterschied ist eine **Eigenschaft des
heutigen Stands**, keine zugesagte Semantik — weder Lastenheft noch
Spezifikation begründen ihn. Verlassen Sie sich für die groben Achsen auf die
Tabelle, nicht auf eine Regel, die Sie sich daraus herleiten.

### Immutabilität über eine Commit-Range prüfen (Modul `vcs`)

Das Modul `vcs` vergleicht den **Core** einer immutablen Datei über zwei
git-Stände — es braucht daher eine **Commit-Range** und ist nie ein Default-Modul.
Aufruf über `--range <base>..<head>` (CI/Push) oder `--staged` (lokaler
pre-commit); ohne `.git`/Range bricht es ab (fail-closed). Das `--print-mk`-Target
`doc-immutable` verteilt diese Prüfung an eigene Repos (kein kopiertes Skript):

```bash
# Range aus dem CI; STAGED=1 für einen lokalen pre-commit-Lauf:
make doc-immutable RANGE="$BASE..$HEAD"
```

### Traceability-Kennungen in Commit-Messages prüfen (Modul `commits`)

Das Modul `commits` prüft, dass jede **Commit-Message** eine Traceability-Kennung
(`id-patterns`) trägt — die verteilbare Form der früheren `trace-check`-Prüfung. Es
liest die Messages über denselben git-Port wie `vcs` (kein git-Binary, kein Netz)
und kennt zwei Quellen: `--range <base>..<head>` (jede Nicht-Merge-Commit-Message
der Range, für CI/Push) und `--commit-msg <datei>` (eine einzelne Pending-Message,
für den commit-msg-Hook; `-` liest von stdin). Merge-/Revert-Betreffe nimmt
`exempt-pattern` aus. Das `--print-mk`-Target `doc-commits` verteilt die
Range-Prüfung an eigene Repos:

```bash
# Range aus dem CI (base..head):
make doc-commits RANGE="$BASE..$HEAD"
```

### Planning-Lifecycle, Closure-Notizen und Wellen-Register prüfen (Modul `planning`)

Das Modul `planning` prüft **hermetisch** (nur der Arbeitsbaum, kein git) eine
Planning-Lifecycle-Invariante: die Roadmap trägt den Ruhe-Marker (`marker`, Default
„Keine aktive Welle") in ihrem `heading`-Block (Default `## Aktuelle Welle`) genau
dann, wenn **kein** `slice-*` (`slice-glob`) im Roadmap-Verzeichnis liegt — sonst
`planning-drift`. Bei fehlender oder **mehrdeutiger** (mehrfacher) Überschrift bzw.
fehlender Roadmap-Datei bricht es fail-closed ab. **Gezählt wird nur außerhalb
von Code-Blöcken:** eine Roadmap, die ihren eigenen Abschnitt als Beispiel zeigt,
gilt nicht als mehrdeutig, und eine Raute-Zeile in einem Beispiel-Block beendet
den Aktiv-Block nicht. Wo der Block endet, entscheidet dieselbe
Abschnitts-Mechanik wie bei den Closure-Notizen — auch eine **eingerückte** H2
und eine **H1** beenden ihn.

**Dritte Fähigkeit, opt-in über `waves.dir`: die Wellen-Register gegen die
Wellen-Dateien.** Dieselbe Invariante eine Ebene höher — vier Aussagen, vier
Grund-Codes: der Aktiv-Status-Abschnitt nennt eine Welle genau dann, wenn
unter den flachen Wellendokumenten **genau eine** Wellen-Kennung liegt
(`wave-drift`); eine Vorschau-Zeile
nennt keine Welle, die schon eine Datei hat (`wave-preview-exists`); jede Zeile
des Abschluss-Registers hat ihre **Ergebnisnotiz** (`wave-results-missing`) —
und jede Ergebnisnotiz ihre Zeile (`wave-unregistered`). Gelesen wird die
**erste Spalte** der Registerzeilen; Zeilen ohne Kennung sind erlaubt (geplante
Wellen tragen Namen), und die Trigger-Spalte darf andere Wellen nennen.
Fail-closed über `wave-drift`: ein unlesbares Wellen-Verzeichnis und eine
fehlende Register-Überschrift. **Setzen Sie beim Einführen den
konsument-gerechten `marker`** — ein fremder Ruhe-Marker erzeugt sonst einen
Artefakt-Befund; und rechnen Sie mit echten Rückständen, das ist der Zweck.

**Zwei Kardinalitäts-Modelle (`waves.mode: one | many`).** Die erste
Wellen-Aussage kennt zwei Betriebsarten. **`one`** (Default) ist das
Singleton-Prädikat oben — genau eine Welle oder Ruhe-Marker; ohne den
Schlüssel ist der Befundsatz byte-identisch. **`many`** deckt das
Offene-Wellen-Modell: der Aktiv-Abschnitt führt je offener Welle einen
Zeiger, mehrere sind legitim, und verglichen werden **Kennungs-Mengen** —
die im `heading`-Block genannten Wellen-Kennungen (Prosa-Zeilen; Fences
zählen nicht, Mehrfachnennung einmal, ob Liste oder Tabelle ist
gleichgültig) gegen die Kennungs-Menge der flachen Wellendokumente, in
beiden Richtungen. Vergleichsgröße sind in beiden Modi **Kennungen**, nicht
Dateien — zwei flache Dokumente derselben Kennung sind ein Element, auch für
das Singleton unter `one`. Das Befund-`target` ist die betroffene Kennung.
Der Ruhe-Marker geht **nicht** ein — er sagt unter diesem Modell etwas anderes (nichts ist beansprucht)
und bleibt Sache von `planning-drift`; der Zustand „Welle eröffnet, noch
nichts beansprucht" trägt Marker **und** Zeiger und ist unter `many` grün,
unter `one` absichtlich rot. **Beachten:** unter `many` zählt jede
Wellen-Kennung in der Abschnitts-Prosa als Zeiger — halten Sie erklärenden
Text dort kennungsfrei. Ein anderer Wert als `one`/`many` — auch der
explizit leere String — bricht mit Exit 2 ab.

Nur `roadmap` ist Pflicht;
`heading`/`marker`/`slice-glob` sind für abweichende Layouts überschreibbar:

```yaml
planning:
  roadmap: docs/plan/planning/in-progress/roadmap.md
```

Dasselbe Modul trägt eine **zweite Fähigkeit**, die die andere Seite derselben
Invariante prüft: die **Struktur** der Closure-Notizen abgeschlossener Pakete
(Aufgabe §4.17). Sie ist opt-in **innerhalb** des opt-in Moduls — ohne
`closure.dir` wird keine Paket-Datei geöffnet:

```yaml
planning:
  roadmap: docs/plan/planning/in-progress/roadmap.md
  closure:
    dir: docs/plan/planning/done       # Aktivierungs-Schalter; leer ⇒ inert
    # glob: '*.md'                     # eigener Kandidaten-Filter; WEGLASSEN ⇒ slice-glob
    heading-pattern: '^#{2,3} .*Closure-Notiz'  # RE2, Default
    min-sentences: 4                   # Satzende-Zeichen außerhalb Code, Default
    boilerplate: []                    # literale Floskeln, case-insensitiv; Default leer
    # placeholder: true                # unausgefüllte Vorlagen-Platzhalter melden; Default AUS
```

`min-sentences` < 1, ein **explizit** leerer oder ungültiger `glob` und ein
leerer `boilerplate`-Eintrag brechen mit Exit 2 ab (jedes wäre eine Aussage, die
nie greift). Diese Schlüssel werden am Config-Rand geprüft, **auch wenn**
`closure.dir` fehlt und die Fähigkeit damit inert ist — sonst entstünde ein
stilles Grün für den Tag, an dem Sie `dir` setzen.

**Zum `glob`:** er ist der Kandidaten-Filter **dieser** Fähigkeit und getrennt
von `slice-glob`. Die beiden stellen verschiedene Fragen — `slice-glob` zählt,
was noch **in Arbeit** ist, `glob` prüft, was **abgeschlossen** ist. Lassen Sie
ihn weg, gilt `slice-glob`; das ist ein Verweis, kein kopierter Wert, und Sie
pflegen genau ein Muster. Setzen Sie ihn erst, wenn die beiden Mengen
auseinandergehen sollen — etwa weil in Ihrem Ruheort auch Wellen- oder
Etappen-Dokumente liegen. **`slice-glob` zu weiten ist dafür der falsche Knopf:**
liegt Ihre Roadmap-Datei selbst im gezählten Verzeichnis, meldet die
Lifecycle-Invariante danach dauerhaft falsch-rot. Die Floskel-Liste ist bewusst leer
vorbelegt: sie entscheidet über rot und grün, und mitgelieferte Phrasen wären in
einem anderssprachigen Repo entweder wirkungslos oder falsch. Nehmen Sie eine
Phrase nur auf, wenn sie im Bestand **null** Treffer hat — sonst färbt sie
Notizen rot, die sehr wohl tragen.

Das `--print-mk`-Target `doc-planning` (hermetisch, **ohne** Range) verteilt die
Prüfung an eigene Repos mit demselben Roadmap-Layout:

```bash
make doc-planning
```

### Getrackt-Status von Link-Zielen prüfen (Modul `tracked`)

Das Modul `tracked` prüft, dass jedes **auflösbare, existierende**
Link-/Bild-**Datei**-Ziel im **git-Index getrackt** ist — ein untracktes oder
gitignoriertes Ziel ist beim Erzeuger grün, wäre aber auf jedem frischen
Klon ein `target-missing`; der Befund `target-untracked` fängt diese
Umgebungs-Drift am Entstehungsort. Die Wahrheit ist der **Index**, nicht
die `.gitignore`-Syntax: eine frisch per `git add` gestagte Datei gilt als
getrackt. Fehlende Ziele bleiben Sache von `links` (kein Doppelbefund);
`exempt-targets` (Globs über den **aufgelösten** Zielpfad — dieselbe Form,
die auch der Befund nennt) nimmt absichtlich untrackte Ziele aus; die Globs
werden beim Laden segmentweise validiert (ungültig ⇒ Exit 2).
Verzeichnis-Ziele und Symlink-Referenzen prüft `tracked` nicht (Symlinks
meldet `links` ohnehin kategorisch). Braucht ein lesbares `.git` unter der
Scan-Wurzel (sonst Exit 2, fail-closed), aber **keine** Commit-Range:

```yaml
tracked:
  exempt-targets: []   # z. B. ["build/**"] für lokal generierte Ziele
```

Das `--print-mk`-Target `doc-tracked` (ohne Range; das `.git` liegt im
read-only-Mount) verteilt die Prüfung:

```bash
make doc-tracked
```

### Einen repo-internen Link-Inhalt gegen Drift pinnen (Modul `pins`)

Das Modul `pins` (opt-in) friert den **Inhalt** eines Link-Ziels **innerhalb des
Repos** auf einen `sha256` ein: ändert sich der Ziel-Span (die verlinkte Datei
bzw. die Heading-Section) nach dem Verlinken, meldet es `link-stale`. So bleibt
eine Referenz, deren Aussage vom zitierten Inhalt abhängt, an genau diesen Inhalt
gebunden. (Das Netz-Gegenstück für **externe** Quellen ist `sources`, s. u.)

Der Pin steht als HTML-Kommentar-Marker **unmittelbar hinter dem Link** derselben
Zeile; gehasht wird der whitespace-normalisierte **rohe** Ziel-Span (ganze Datei
ohne Anker, sonst die Heading-Section inkl. Fenced-Code) — reine Reflow-/Umbruch-
Änderungen absorbiert die Normalisierung.

**Welcher Anker gemeint ist, entscheidet dieselbe Antwort wie im Modul
`anchors`:** Heading-Slug samt Duplikat-Zähler (`#alt-1`), prozent-kodierte
Fragmente (`#a%20b`) und Inline-HTML-Anker — jeweils **außerhalb** von
Fenced- und Inline-Code, und **case-sensitiv**. Was `anchors` nicht als Anker
kennt, macht das Ziel unauflösbar; `pins` schweigt dann, und der fehlende Anker
erscheint als `anchor-missing` statt als stiller Verlust des Drift-Schutzes.

**Einen Pin anlegen — „einmal laufen, Hash kopieren".** Sie kennen den `sha256`
des Ziel-Spans zunächst nicht. Setzen Sie den Marker mit einem Platzhalter-Hash
unmittelbar hinter den Link:

```markdown
[Abschnitt](../spec/ziel.md#abschnitt) <!-- dpin: sha256:0000000000000000000000000000000000000000000000000000000000000000 -->
```

Dann in zwei Läufen:

1. `d-check --enable pins` laufen lassen. Der `link-stale`-Befund führt in seiner
   Meldung den **vollen** errechneten Hash: `errechnet sha256:<hex> (voller
   Ist-Hash zum Re-Pinnen)`.
2. Diesen `<hex>` in den Marker kopieren und erneut laufen lassen → **kein
   Befund**: der Pin ist gesetzt.

Danach schlägt jede inhaltliche Drift des Ziel-Spans als `link-stale` an, bis Sie
bewusst neu pinnen (denselben Weg). `pins` ist strikt opt-in und **netzlos** und
prüft nur repo-interne, auflösbare Ziele (ein struktureller Fehler bleibt
`links`/`anchors`, kein Doppelbefund); ohne aktives `pins` ist der Befundsatz
byte-identisch.

### Externe Quellen auf ihren Inhalt pinnen (Modul `sources`)

Das Modul `sources` (opt-in, **Netz**) ist die **zweite** Netz-Tür neben `external`
und prüft **Inhalt** statt bloßer Erreichbarkeit: eine auf einen `sha256` gepinnte
externe `http(s)`-Quelle wird geholt, gehasht und mit dem Pin verglichen. Weicht
der Hash ab, meldet es `source-drift` — die Meldung führt den **vollen errechneten
`sha256`** mit, sodass ein Re-Pin „einmal laufen, gemeldeten Hash kopieren" ist. Ist
die Quelle nicht materialisierbar (Netzfehler, HTTP-Status ≥ 400, Timeout,
Größenlimit **oder** unter `unpack: zip` kein gültiges Zip), meldet es
`source-unreachable` (getrennt von `source-drift`).

Der Pin wird auf **zwei** Wegen deklariert:

- **Marker** am Link — `[text](https://…) <!-- source-pin: sha256:<hex> -->` bindet
  an den unmittelbar vorausgehenden `http(s)`-Link derselben Zeile; für ein
  Archiv-Ziel trägt der Marker zusätzlich das Schlüsselwort `zip`
  (`<!-- source-pin: zip sha256:<hex> -->`).
- **Config** — ein Top-Level-Block `sources:` mit Einträgen
  `{url, sha256, unpack: none|zip}`.

Zwei Quelltypen: eine **Einzeldatei** (`unpack: none`, Default) wird über ihre rohen
Antwort-Bytes gehasht; ein **Archiv** (`unpack: zip`) über ein pfad-sortiertes
**Content-Manifest** — damit ist der Pin invariant gegen die Zip-Eintrags-Reihenfolge.
`sources` prüft **nur** absolute `http`/`https`-Ziele; repo-interne Verweise bleiben
Domäne von `links`/`pins` (kein Doppelbefund). Es gibt **kein** `sources.scope`: die
Marker greifen über den globalen Scan-Bereich (`scan.roots`/`scan.ignore`). Der
`sha256` wird case-insensitiv geführt; der Fetch ist größenbegrenzt (Body- und
Entpack-Obergrenze, Zip-Bomben-Schutz, begrenzte Redirects). Wie `external` öffnet
`sources` nur bei ausdrücklicher Aktivierung Netz und ist nie Default-Modul; ohne
aktives `sources` ist der Befundsatz byte-identisch. Eine malformte Direktive (kein
`sha256:<hex>`) oder ein ungültiger Config-Eintrag (fehlende `url`/`sha256`,
unbekanntes `unpack`) bricht fail-closed mit Exit 2 ab.

```yaml
sources:                                     # opt-in; NETZZUGRIFF, zweite Netz-Tür neben external
  - url: https://example.org/regelwerk.md    # Einzeldatei: Hash der Roh-Bytes
    sha256: 0000000000000000000000000000000000000000000000000000000000000000
  - url: https://example.org/bundle.zip      # Archiv: Hash des Content-Manifests
    sha256: 0000000000000000000000000000000000000000000000000000000000000000
    unpack: zip                              # none (Default) | zip
# aktiviert über --enable sources bzw. modules: [..., sources]
```

### Das Versions-Register `version.md` aufbauen (für Modul `versions`)

Das Modul `versions` liest die aktuelle Version aus `version.md#aktuell`
(Schlüssel `current-from`) und prüft damit die gepinnten `ghcr`-Image-Verweise.
Diese Datei ist zugleich ein **auflösbares Link-Ziel** für Erwähnungen der
eigenen Releases. Ein eigenes Repo baut sie nach diesem Muster nach — zwei
Abschnitte und eine Anker-Mechanik:

- **`## Aktuell`** — eine Zeile, die die aktuelle Version nennt. Der
  Heading-Slug `#aktuell` ist der **stabile** Verweispunkt
  (`current-from: version.md#aktuell`); er zeigt immer auf die jeweils aktuelle
  Version, nie auf eine feste Nummer. Welcher Anker das ist, entscheidet
  dieselbe Antwort wie im Modul `anchors` — inklusive Duplikat-Zähler,
  Prozent-Dekodierung und **Groß-/Kleinschreibung**; ein Anker in einem
  Code-Block ist keiner, und `current-from` bricht dann fail-closed ab.
- **`## Verlauf`** — eine Tabelle mit einer Zeile pro Release (Version, Datum,
  Release-Link).
- **Der explizite HTML-Anker `<a id="vX.Y.Z"></a>`** — nur die **aktuelle**
  Versions-Zeile trägt ihn. Grund: der automatische Markdown-Slug verschluckt
  die Punkte einer Versionsnummer (`v1.2.0` → `v120`), sodass ein fester Pin
  `…#v1.2.0` sonst nicht auflöst; der explizite Anker macht die Version wörtlich
  (mit Punkten) verlinkbar.

**Die Anker-Wanderung ist Absicht:** Bei jedem Release wandert der `<a id>`-Anker
auf die neue aktuelle Version; die bisherige Zeile **verliert** ihn. Dadurch
**bricht** jeder feste Pin auf eine veraltete Version (`anchor-missing` über das
Modul `anchors`) — ein vergessener Versions-Bump fällt so auf, statt still auf
eine alte Version zu zeigen. Der zugehörige Release-Ablauf steht in
[releasing.md](releasing.md).

Minimales Muster zum Übernehmen:

```markdown
# <projekt> — Release-Register

## Aktuell

Aktuelle Version: [`vX.Y.Z`](#vX.Y.Z) — <datum>.

Aus anderen Dokumenten stabil referenzierbar als `version.md#aktuell`
(zeigt immer hierher, nie auf eine feste Nummer).

## Verlauf

| Version                      | Datum   | Release                     |
| ---------------------------- | ------- | --------------------------- |
| `vX.Y.Z` <a id="vX.Y.Z"></a> | <datum> | [Tag vX.Y.Z](<release-url>) |
| `vX.Y.W`                     | <datum> | [Tag vX.Y.W](<release-url>) |
```

Beim nächsten Release: die `## Aktuell`-Zeile auf die neue Version ziehen, eine
neue `## Verlauf`-Zeile einfügen und den `<a id>`-Anker von der bisherigen auf
die neue Version verschieben (die bisherige Zeile behält nur ihren Text).

### Modul-lokaler Scan-Bereich

Ein Modul kann einen eigenen, vom globalen abweichenden Scan-Bereich
bekommen:

```yaml
ids:
  scope:
    roots: [spec, docs/user]
```

### Konfigurierbare RTM-Quellen (`trace`)

Die Traceability-Matrix (`--trace`, siehe §4.12) erkennt Anforderungen und ihre
ADR-/Slice-Referenzen standardmäßig nach d-checks eigener Konvention (Kennungen
`<PREFIX>-FA-…`/`-QA-…` im Lastenheft, Slice-Dateien `slice-NNN-….md`,
ADR-Dateien `NNNN-….md`). Folgt Ihr Repo einer anderen Konvention, überschreibt
der optionale `trace`-Block (oben im Beispiel) die Quell-Achsen:
`requirements` (`source`, `id-pattern`, `format` und bei Tabellen `table.*`)
sowie je Referenzklasse `adrs`/`slices`
(`dir` + `file-pattern` + `id-prefix`; die **Capture-Gruppe** der `file-pattern`
liefert die Owner-Kennung, der `id-prefix` wird ihr vorangestellt). Jedes Feld
ist optional; ein abwesendes fällt auf den Default zurück, und **ohne
`trace`-Block ist die RTM-Ausgabe byte-identisch** zum bisherigen Verhalten. Ein
nicht kompilierbares Muster oder eine `file-pattern` ohne Capture-Gruppe ist ein
Konfigurationsfehler (Exit 2, fail-closed).

`requirements.format` ist `headings` (Default) oder `table`. Bei `headings`
definiert nur eine ATX-Überschrift mit einer vollständig passenden ID im ersten
Token eine Anforderung. Bei `table` sind `table.id-column` und
genau eine von `table.text-column` oder der nichtleeren Alternativliste
`table.text-columns` Pflicht; jeder Listenwert muss quellweit mindestens einmal
gebunden werden. `table.modality-column` ist optional. Die Namen
werden nach Trimmen exakt an die Tabellen-Header gebunden. Fehlende, leere oder
doppelte konfigurierte Header, eine relevante Zeile mit falscher Zellenzahl und
doppelte IDs sind unter `table.duplicate-ids: error` (Default)
Konfigurationsfehler (Exit 2). `first` behält bei historischen Duplikaten die
erste, `last` die letzte Definition. Nur eine vollständig auf
`id-pattern` passende ID-Zelle definiert eine Anforderung. Escaped Pipes und
Pipes in einzeiligen, passend begrenzten Code-Spans bleiben Bestandteil der
Zelle; Tabellen in Code-Fences und mehrzeilige Zellen werden nicht gelesen.

Ein tabellarisches Lastenheft binden Sie mit `format: table` und den vorhandenen
ID-/Text-/Modalitätsheadern:

```yaml
trace:
  requirements:
    source: spec/lastenheft.md
    id-pattern: 'F-[0-9]+'
    format: table
    table:
      id-column: ID
      text-column: Anforderung
      modality-column: Modalität   # optional; sonst wird die Textspalte klassifiziert
      duplicate-ids: error         # Default; alternativ first oder last
```

Verwendet ein Lastenheft mehrere Text-Header für dieselbe Rolle, listen Sie sie
in `table.text-columns` explizit; historische Mehrfachdefinitionen brauchen eine
bewusste `duplicate-ids`-Politik (`duplicate-ids: error` ist der Default,
`first`/`last` sind die deterministischen Overrides). Eine ID-Regex allein migriert das Format
nicht; für native Tabellen müssen `format: table`, `table.id-column`, genau eine
von `table.text-column` oder `table.text-columns` und optional
`table.modality-column` konfiguriert sein. Für Tabellen außerhalb dieser
Pipe-Grammatik überführen Sie die Zeilen in Heading-plus-Body-Abschnitte (ID und
Titel in der Überschrift, der normative Text samt Modalverb im Body) oder
erzeugen daraus deterministisch ein separates Heading-Dokument, dessen Drift ein
eigener Konsistenzsensor gegen die autoritative Tabelle absichert.

Zur Tabellen-Grammatik im Detail: die Trennzeile folgt GFM — **ein** Bindestrich
je Zelle genügt (`| - |`, `| :-: |`), Doppelpunkte für die Ausrichtung bleiben
erlaubt. Folgen zwei Tabellen ohne Leerzeile aufeinander, beendet der Beginn der
**zweiten** die erste, sofern deren Header eine konfigurierte Rolle bindet —
beide werden erkannt; ein Header ohne gebundene Rolle beendet nicht. Steht
d-checks eigene Ignore-Direktive `<!-- d-check:ignore … -->` hinter der letzten
Pipe einer Datenzeile, zählt die überzählige Kommentar-Zelle nicht als Spalte;
genau **eine** ganzzellige Kommentar-Zelle wird toleriert, zwei oder eine
Nicht-Kommentar-Zelle bleiben Exit 2.

Eine **nichtleer explizite** `requirements.source` oder `format: table`
aktiviert den Nullmengen-Guard: Fehlt die Quelldatei oder werden darin null
Anforderungen erkannt, endet `--trace` mit Exit 2 vor jeder RTM-Ausgabe.
`source: ""` gilt dagegen wie ein abwesendes Feld und hält den bisherigen
Default-Fallback einschließlich leerer RTM/Exit 0 kompatibel. Für eine
zusätzliche Bestandsinvariante kann ein Konsument die erkannte Gesamtzahl gegen
seine erwartete Zahl plausibilisieren.

Der ADR-/Slice-Referenzscan läuft rekursiv über die gefundenen Dateien unter
dem jeweiligen `dir`. `file-pattern` wird gegen den **Basisdateinamen** und
nicht gegen den vollständigen Pfad geprüft. Dateien ohne Match werden
übersprungen; bei einem Match bildet Capture-Gruppe 1 zusammen mit
`id-prefix` die Owner-Kennung. Danach wird der gesamte Dateiinhalt nach
Treffern von `requirements.id-pattern` durchsucht. Owner je Requirement werden
dedupliziert und deterministisch sortiert.

Zusätzlich liest der opt-in **`trace.coverage`**-Block (eine **Liste** benannter
Quellen) **kuratierte Matrizen** als **eigene Coverage-Dimension** ein — für
Anforderungen, deren Abdeckung weder in ADRs noch Slices liegt (siehe §4.13). Je
Quelle: `files` (explizite Pfade — keine `dir`/`file-pattern`, gegen
ADR-Kontamination), `label` (Owner-Kennung in der eigenen Coverage-Spalte),
`ranges` (Default true; expandiert `<FAM>-AAA..BBB` und `<FAM>-AAA/BBB/CCC`
breiten-erhaltend, gegen `requirements.id-pattern` validiert — **auch wenn die
Kennung verlinkt ist**: `` [`GG-QA-001`](…)..006 `` wird wie `GG-QA-001..006`
gelesen, siehe unten) sowie `sections`
(Whitelist) / `exclude-sections` (Blacklist) — beide über die
`matrix.exclude-sections`-Span-Semantik, mit dem **vollen Heading-Klartext** als
Namen. Eine Anforderung ist dann **Waise** nur ohne Slice **und** ohne Coverage.
Fail-closed (Exit 2): fehlende `files`-Datei, leeres `label`, ein Sektionsname
ohne Heading-Treffer, eine ungültige Range (`AAA>BBB`/Breite). Ist **keine**
Coverage-Quelle konfiguriert, erscheint **keine** Coverage-Spalte und die RTM
bleibt byte-identisch.

**Range-Notation unter Linkpflicht.** Erzwingt Ihr Repo Links für Kennungen
(`ids` mit `link-policy: always`, siehe §4.6), steht die Kennung als
`` [`GG-QA-001`](…) `` — die Fortsetzung `..006` folgt ihr dann nicht mehr
unmittelbar. d-check überspringt deshalb **genau ein** Link-Suffix und liest die
Fortsetzung dahinter; verlinkt und unverlinkt ergeben dasselbe. Bewusst nur
eines: **nicht** übersprungen werden Whitespace, Emphasis, ein zweites
Link-Suffix oder Text zwischen `)` und der Fortsetzung — dort wird nicht
expandiert, statt die Absicht zu raten. Klammern **im Linkziel** sind dagegen
unproblematisch: das Ziel wird klammer-balanciert abgegrenzt, genau wie bei
`links`/`ids` — `` [`GG-QA-001`](https://x.org/A_(b))..003 `` expandiert. Gilt
gleichermaßen für `trace.cross-consistency`. (Bis v0.44.0 expandierte eine
verlinkte Range gar nicht — das erzeugte in `trace.coverage` falsche Waisen.)

Genau zwei Expansionsformen sind zugesagt: die Range `<FAM>-AAA..BBB` und die
Aufzählung `<FAM>-AAA/BBB/CCC`. Eine **Komma-Kurzform** (`GG-SCN-001, 007`, auch
hinter einer Range wie `GG-SCN-001..005, 007`) ist **nicht** zugesagt und bricht
mit Exit 2 ab, statt `007` still fallen zu lassen; ein Komma **vor** einer
vollständigen Kennung ist erlaubt.

Der opt-in **`trace.requirements.modality`**-Block klassifiziert jede Anforderung
nach RFC-2119-Stufe (siehe §4.14). `levels` (Map Stufe → Modal-Verb-Keywords;
leer/`{}` ⇒ Built-in DE+EN-Defaults inkl. `DARF NICHT`→must, `MUSS NICHT`→may)
und `require-levels` (welche Stufen `--require-complete` gaten, Default `[must]`).
Ein gesetztes `levels` **ersetzt** die Defaults vollständig (kein Merge): wer
`levels` konfiguriert, muss **jede** gewünschte Stufe und **jedes** gewünschte
Verb selbst auflisten (auch die gewünschten Built-ins) — die vollständige
Default-Menge steht in `spec/spezifikation.md` (Abschnitt
`trace.requirements.modality`).
Klassifiziert wird über den ersten (frühesten), bei Gleichstand **längsten**
Keyword-Treffer. Im Heading-Format ist die Quelle der **normalisierte**
Body-Abschnitt, im Tabellenformat exklusiv `table.modality-column`, sofern
konfiguriert, andernfalls die Textzelle (Emphasis raus, Whitespace/
Umbrüche zu einem Leerzeichen), case-insensitiv und wortgrenzen-genau (`\b` ist
ASCII — ein konfiguriertes Umlaut-Rand-Keyword träfe nicht); kein Treffer ⇒ Stufe
`unknown`. **Aktiv schon bei bloßer Präsenz** (`modality: {}`). Fail-closed
(Exit 2): leerer Stufen-Name/leeres Keyword, reservierter Name `unknown` in
`levels`, gleiches Keyword in zwei Stufen, ungültiges `require-levels`. Ohne
`modality` byte-identisch (keine Spalte, kein Feld, `--require-complete` gatet
alle Waisen).

Der opt-in **`trace.cross-consistency`**-Block vergleicht **zwei unabhängig
gepflegte Sichten derselben Anforderung→Design-Relation** (siehe §4.15): die
**Vorwärts**-Sicht (eine RTM-Tabelle: je Anforderung die Design-Artefaktmenge)
gegen die **Rückwärts**-Kanten (je Design-Artefakt seine Anforderungen — dort
authort, wo das Design lebt, und damit die **Quelle der Wahrheit**). Gemeldet
werden je Anforderung die beiden Mengendifferenzen mit Richtungslabel und
`Datei:Zeile`. Beide Sichten sind kuratierte Markdown-Tabellen und werden über
denselben header-gebundenen Reader gelesen wie `requirements.format: table`.

- **`forward`** (Pflicht): `file`, `req-column`, `design-column`,
  `design-pattern` (extrahiert die Artefakt-IDs aus der Design-Zelle); optional
  `req-pattern` (erkennt die Anforderungs-IDs der `req-column` — Default
  `requirements.id-pattern`, siehe **Scope-Falle** unten),
  `sections`/`exclude-sections` (Span-Semantik wie `trace.coverage`) und
  `ranges` (Default true, expandiert `<FAM>-AAA..BBB` in der ID-Spalte).
- **`backward`** (Pflicht): `file`, `edge-column` (z. B. `Bezug` — sie **allein**
  macht eine Tabelle relevant), `req-pattern` (erkennt die Anforderungs-IDs in
  der Kanten-Zelle); `artifact-id-column` ist der Sentinel `first` (Default:
  erste Spalte — für über die Tabellen **heterogene** ID-Header wie
  `Kennung`/`Port-ID`/`Komponente`) **oder** ein Header-Name, der dann je
  relevanter Tabelle genau einmal vorkommen muss; optional `sections`, `ranges`.
- **`mode`**: `equal` (Default; beide Differenzen gaten) oder `superset` (nur
  „Rück-Kante ohne RTM-Eintrag").
- **`exclude-req`**: Regex für **Ableitungssprünge** — Mittelschicht-IDs, die
  Rück-Kanten nennen, ohne eine eigene RTM-Zeile zu haben. Ein benanntes Ventil
  mit eigener Drift-Gefahr (wie `matrix.exclude-sections`), kein gelöstes Problem.

**Namensraum:** die Artefakt-IDs beider Sichten werden bewusst mit **demselben**
`forward.design-pattern` erkannt — nur so ist der Mengen-Diff bedeutungsvoll.

**Scope-Falle: der Vergleich ist nicht die RTM.** Welche Anforderungen verglichen
werden, entscheiden `forward.req-pattern` und `backward.req-pattern` — **nicht**,
ob eine Anforderung in Ihrer RTM steht. Beide fallen per Default auf
`requirements.id-pattern` zurück. Das ist bequem, solange beide Mengen
zusammenfallen — und eine Falle, sobald Sie Ihre RTM **bewusst scopen**:
Schließt `requirements.id-pattern` etwa Architektur-Meta aus (weil es keine
Anforderung ist), findet die Vorwärts-Sicht dort **nichts**, und **jede**
Rück-Kante dieser Familie wird als „Rück-Kante, ohne RTM-Eintrag" gemeldet —
während die echte Differenz in der Gegenrichtung *verschwindet*. Das sieht wie ein
Treffer aus und ist keiner. **Gegenmittel:** `forward.req-pattern` explizit auf die
Familien setzen, die verglichen werden sollen; die RTM bleibt davon unberührt.

Der Abgleich ist **advisory**: `--trace` bleibt Exit 0 und die RTM unverändert
(keine zusätzliche Spalte — die Differenzen erscheinen als eigener Abschnitt
darunter). Den Exit-Code ändert allein das globale `--require-complete`
(≥ 1 Differenz ⇒ Exit 1); einen block-lokalen Schalter gibt es nicht.

**Fail-closed (Exit 2):** fehlendes `forward`/`backward`, unbekannter `mode`,
nicht kompilierbares Regex, leeres Pflichtfeld, fehlende Sicht-Datei, ein
Sektionsname ohne Heading-Treffer, keine Tabelle mit den konfigurierten Headern,
ein mehrfach vorkommender Rollen-Header, eine relevante Zeile mit falscher
Zellenzahl — sowie ein **vakuumer Abgleich**: greifen die Muster am Inhalt vorbei
(beide Sichten kantenleer) oder verschluckt `exclude-req` jede Anforderung, dann
könnte der Lauf konstruktionsbedingt nie eine Differenz melden, und ein
`0 Differenz(en)` behauptete eine nie geprüfte Konsistenz. Dasselbe gilt für eine
kantenleere Rück-Sicht unter `mode: superset`. **Nicht** fail-closed ist eine
einseitig leere **Vorwärts**-Sicht bei gepflegten Rück-Kanten: das ist der
erwartete Zustand, solange die RTM-Tabelle noch nicht auf konkrete IDs
restrukturiert ist — sie meldet dann jede Rück-Kante laut. Ohne den Block ist die
RTM byte-identisch.

### Struktur-Invarianten in Dokumenten (Modul `structure`)

Das Modul prüft nicht Referenzen, sondern **Form**: ob ein Abschnitt existiert,
ob er Substanz trägt, ob eine Vorlage ausgefüllt wurde. Jede Regel benennt ihre
Dateien **selbst** — unabhängig von `scan.roots`, deshalb gibt es kein `scope`.

```yaml
structure:
  - files: "docs/plan/planning/done/slice-*.md"   # Pflicht; Glob, Wurzel-relativ
    section-pattern: '^#{2,3} .*Closure-Notiz'    # RE2 …
    # section: "## 9. Closure-Notiz"              # … oder Klartext, EXAKT inkl. #-Folge
    # sections: one                               # one (Default) | each
    non-empty: true
    # min-sentences: 4
    # max-tasks: 0
    # max-open-tasks: 0                          # OFFENE Task-Items (`[ ]`) auf den ROHEN Zeilen:
    #                                             # alle vier Listen-Marker, Fence aussen vor,
    #                                             # Inline-Code NICHT — ein Befund je Item auf seiner
    #                                             # Zeile ⇒ sonst section-tasks-open
    # forbid-pattern: 'TODO'
    # require-pattern: 'Beleg'
    # require-all: ["Beleg", "Lernsignal"]
    # table:                                      # die beiden TABELLEN-Bedingungen unter einer Klammer
    #   order: desc                               # asc|desc: Chronologie-Monotonie der Schlüsselspalte
    #   order-column: 1                           # 1-basierte Schlüsselspalte (Default 1; nur mit order)
    #   column:                                   # Zellengrenzen je Spalte — beliebig viele, EIN Selektor
    #     - name: Titel                           # Spalte über ihren KOPFZEILEN-Namen (nicht Position)
    #       cell-max-chars: 200                   # Obergrenze in ZEICHEN
    #       cell-min-chars: 10                    # Untergrenze — ohne sie passiert die LEERE Zelle
    # headings-match: '^SPEC-[0-9]{3} '           # JEDE Überschrift im Abschnitt matcht dieses Muster
    # headings-level: 4                           # geprüfte ATX-Ebene (Default: Abschnitts-Ebene + 1);
    #                                             # MUSS tiefer als der Abschnitt liegen — eine gleiche
    #                                             # oder flachere Ebene kommt in ihm nicht vor, die
    #                                             # Bedingung wäre dann wirkungslos wahr
    # hint: "Haken setzen oder Slice zurückführen"   # VERFASSTE Erläuterung dieser Regel:
    #                                             # sie wird zur vierten Spalte des Befunds und
    #                                             # zur Hinweis-Zeile in --doctor. Explizit leer
    #                                             # ⇒ Exit 2
    # exempt-paths: []
```

**Sechs Dinge, die überraschen können.** Erstens: eine Regel, die **keine** Datei
trifft, meldet `section-missing` — auch dann, wenn erst `exempt-paths` die Menge
geleert hat. Eine Regel zu schreiben ist die Behauptung, dass es die Dateien
gibt; ein leer laufendes Gate darf nicht Erfolg melden. Zweitens: bei `section`
gehört die **`#`-Folge** zum Vergleich, `## E` trifft also kein `### E`.
Drittens: eine Marke aus `require-all` zählt nur am **Zeilen-Anfang** als
hervorgehobener Textlauf (`- **Beleg:**`), nicht im Fließtext.

Viertens — und das entscheidet über Ihre Schwelle: der Abschnitt wird vor jeder
Bedingung **einmal bereinigt**, Fenced-Code entfernt **und Inline-Code geleert**,
genau wie beim Closure-Gate des Moduls `planning`. Ein `forbid-pattern` auf ein
Wort in Backticks trifft deshalb **nicht**, und `min-sentences` zählt nur
Satzende-Zeichen **vor Whitespace oder Zeilenende**: die Punkte in einer
Versionsnummer wie 0.57.0 oder in einem Dateipfad tragen keinen Satz und zählen
nicht mit.

Fünftens: die **Chronologie-Bedingung** (`table.order: asc|desc`) ist die eine
Ausnahme von dieser Bereinigung — sie liest die Zellen der Schlüsselspalte
(`table.order-column`, 1-basiert) **roh**, weil reale Schlüssel in Release-Registern
in Inline-Code stehen. Verglichen wird **typisiert** (ISO-Datum `JJJJ-MM-TT`
oder Punkt-Version wie `1.10`, segmentweise numerisch — `1.10` kommt **nach**
`1.9`) und **nicht-strikt** monoton je zusammenhängender Tabelle; Kopf- und
Trennzeile zählen nicht. Eine Zeile, die die Richtung bricht, meldet
`section-unordered`; eine Zelle ohne typisierbaren Schlüssel — oder mit
anderem Typ als ihr Vorgänger — meldet `section-cell-untyped` an **genau
dieser** Zeile, dahinter setzt der Vergleich neu auf. Ein Abschnitt ganz ohne
Tabellen-Datenzeile meldet den Leerlauf: die Bedingung zu setzen ist die
Behauptung, dass hier eine chronologische Tabelle steht.

Sechstens: die **Überschriften-Bedingung** (`headings-match`) ist die **zweite**
Ausnahme von der Bereinigung — sie liest nicht den Abschnitts-Text, sondern die
Überschriften selbst, mit **derselben** Erkennung, mit der das Modul den
Abschnitt findet (beliebig viel führender Weißraum, Leerzeichen **oder** Tab als
Trenner, außerhalb von Fenced-Code). Sie sagt positiv „**jede** Überschrift
dieses Abschnitts genügt diesem Muster", verglichen wird der Überschriften-**Text**
ohne die `#`-Folge, und gemeldet wird `section-heading-mismatch` **je**
verletzender Überschrift auf **ihrer** Zeile — nicht am Abschnittskopf, denn dort
steht die Reparatur nicht. Geprüft wird genau eine Ebene: `headings-level`, per
Default die Abschnitts-Ebene + 1 (die unmittelbaren Unterabschnitte). Zwei
Grenzen sind benannt, nicht geschlossen: trägt der Abschnitt **keine** Überschrift
der geprüften Ebene, ist die Bedingung wirkungslos wahr, und eine Ebene
**flacher** als der Abschnitt kann in ihm gar nicht vorkommen. Verwechseln Sie
den Schlüssel nicht mit `planning.closure.heading-pattern`: jener ist ein
**Selektor** (welcher Abschnitt), dieser eine **Bedingung** (welche Form).

Siebtens: die **Zellenlängen-Bedingung** (`table.column`, je Eintrag ein `name`
mit `cell-max-chars` und/oder `cell-min-chars`) begrenzt eine Tabellenspalte auf
eine Spanne in **Zeichen** — nicht Bytes, ein Umlaut ist ein Zeichen. Sie ist
die dritte Bedingung auf den rohen Abschnitts-Zeilen und die einzige, die ihre
Spalte über einen **Namen** adressiert: `name` nennt die
**Kopfzelle**, nicht eine Position. Das ist Absicht — eine eingefügte Spalte
verschöbe eine Positions-Angabe **still** auf die falsche Spalte, während ein
umbenannter Kopf **laut** meldet (`section-column-missing`).

**`column` ist eine Liste, und das ist der Punkt.** Mehrere Spalten desselben
Abschnitts sind Einträge unter **einem** Selektor — vorher kostete jede weitere
Spalte eine eigene Regel mit wiederholtem `files`/`section`. Dieselbe Spalte
zweimal zu nennen ist ein Config-Fehler (Exit 2): zwei Zusagen über dieselbe
Zelle trügen dasselbe Befund-Ziel.

Gebunden wird **je Tabelle**: trägt eine Tabelle des Abschnitts den Namen
nicht, ist sie für diese Bedingung irrelevant — eine Nebentabelle schaltet die
Messung der anderen nicht ab. Trägt eine Kopfzeile ihn **mehrfach**, reicht
eine Datenzeile nicht bis zur Spalte, oder bindet **keine** Tabelle sie, ist
das jeweils `section-column-missing`.

**Setzen Sie eine Untergrenze, wenn Sie „gefüllt" meinen.** Eine Obergrenze
allein lässt die **leere** Zelle passieren, denn null Zeichen liegen unter
jeder Schwelle. `cell-min-chars` schließt das, mit eigenem Grund-Code
`section-cell-undersized` — die Reparatur ist ausfüllen, nicht kürzen.

Zwei Dinge sind benannt, nicht geschlossen: gemessen wird die Zelle, **wie sie
dasteht** — eine Zelle aus einem einzigen langen Link ist lang, auch wenn ihr
sichtbarer Text kurz ist. Und die Bedingung misst **eine** Spalte je Regel; wer
mehrere begrenzen will, schreibt mehrere Regeln über denselben Abschnitt. Das
geht, weil die Regel-Identität die benannte Spalte mitträgt — sonst wären zwei
Befunde derselben Zeile nicht unterscheidbar.

**Messen Sie, bevor Sie eine Regel aktivieren.** In diesem Repo hat genau das
eine Regel verhindert, die plausibel klang und falsch war: „abgeschlossener
Slice ohne offene Task-Boxen“ meldete 32-mal — und jedes Mal zu Recht offen,
weil die **Welle** den Punkt einlöst, nicht der Slice.

## 6. Regelmodule

| Modul       | Standard      | Prüft                                                                                    | Grund-Codes                                                 |
| ----------- | ------------- | ---------------------------------------------------------------------------------------- | ----------------------------------------------------------- |
| `links`     | aktiv         | lokale Links/Bilder: Ziel existiert, innerhalb des Repos. **Opt-in** `resolve-from`: Dateien in wandernden Lifecycle-Verzeichnissen lösen jedes relative Ziel von **jedem** Ort ihrer Gruppe auf — überall auf dasselbe Ziel; fail-closed bei Gruppe ohne existenten Ort / Ort als Datei | `target-missing`, `repo-escape`, `symlink`, `link-position-dependent` |
| `anchors`   | aktiv         | Heading-Anker (GitHub-Slugs), inkl. Inline-HTML-Anker                                    | `anchor-missing`                                            |
| `ids`       | opt-in        | Linkpflicht für Kennungen im Fließtext                                                   | `id-unlinked`                                               |
| `matrix`    | opt-in        | erlaubte Referenzrichtung und -status zwischen Dokumentklassen                           | `matrix-forbidden`, `matrix-inactive`                       |
| `codepaths` | opt-in        | explizite Pfade in Inline-Code existieren; opt-in `check-lines`: `datei:<von>-<bis>`-Zeilen-Referenzen verifizieren | `codepath-missing`, `citation-out-of-range`, `citation-inverted-range` |
| `citations` | opt-in        | **keine feine Ventil-Achse** (kein `exempt-paths`, kein `ignore-refs`, kein Zeilen-Marker; grob wirken `scan.ignore` und `citations.scope`) — wortgleiche Zitate: `<!-- d-check:cite <pfad>:<von>-<bis> -->` markiert das folgende Zitat (`>`-Block oder inline `„…"`/`"…"`; Leerzeilen dazwischen sind unschädlich, ein Code-Block dazwischen trennt), das ein whitespace-normalisierter **Teilstring** der Quell-Spanne sein muss; fail-closed bei malformter Direktive | `citation-mismatch`, `citation-out-of-range`, `citation-inverted-range` |
| `spans`     | opt-in        | ungeschlossene Code-Spans, verschachtelte Links, unbalancierte Fences                    | `span-unclosed`, `span-nested-link`, `fence-unclosed`       |
| `hostpaths` | opt-in        | host-lokale absolute Pfade (Maschinen-Layout-Leak)                                       | `hostpath-forbidden`                                        |
| `diagrams`  | opt-in        | Kennungen in Diagramm-Fences (Default `mermaid`) existieren in ihrer `defined-in`-Quelle. **Beide Ventile** wie `ids`/`codepaths`: `exempt-paths` (datei-weit) und der Zeilen-Marker `d-check:ignore` — hier als **Token** statt HTML-Kommentar, und auf der **Öffnungszeile** des Fence für den ganzen Block | `diagram-id-undefined`                                      |
| `versions`  | opt-in        | gepinnte `ghcr`-Image-Verweise tragen die aktuelle Version (aus `version.md#aktuell`), auch in Fences. **Mehrere Muster-Quellen-Paare** über `versions.patterns` (eigenes Release **und** ein fremder gepinnter Stand); die Kurzform `pin-pattern`/`current-from` **ist** die einelementige Liste, beide Schreibweisen zugleich ⇒ Exit 2. Ausnahmen sind paar-lokal, der Zeilen-Marker ist es nicht | `version-stale`                                             |
| `pins`      | opt-in        | Content-Pin (`<!-- dpin: … -->`): Ziel-Span eines Links unverändert seit dem Verlinken | `link-stale`                                                |
| `immutable` | opt-in        | Immutabilitäts-Pin (`<!-- immutable: … -->`): normalisierter **Core** einer Datei (ohne Marker-Zeile + `exclude-sections`) unverändert seit dem Pinnen; hermetisch (kein git) | `core-drift`                                                |
| `vcs`       | opt-in (git)  | git-Diff-Immutabilität: **Core** einer immutablen Datei (`immutable-when`) unverändert über eine Commit-Range (`--range`/`--staged`); liest `.git` read-only (kein git-Binary, kein Netz). `immutable-when` und die Kopf-Status-Zeile gelten nur **außerhalb von Code-Blöcken** — eine Datei, die ihren eigenen Kopf als Beispiel zeigt, wird dadurch nicht immutabel | `core-drift-vcs`                                            |
| `commits`   | opt-in (git)  | Traceability-Kennung (`id-patterns`) in jeder Commit-Message einer Range (`--range`) bzw. der Pending-Message (`--commit-msg`); liest `.git` read-only (kein git-Binary, kein Netz) | `commit-untraceable`                                        |
| `planning`  | opt-in        | Zwei Seiten derselben Lifecycle-Invariante. **Eintritt:** der Ruhe-Marker (`marker`) steht im `heading`-Block (Default `## Aktuelle Welle`) genau dann, wenn kein `slice-*` (`slice-glob`) im Verzeichnis liegt. **Austritt** (zusätzlich opt-in über `closure.dir`): die **Struktur** der Closure-Notizen abgeschlossener Pakete — Abschnitt vorhanden, genug Satzende-Zeichen außerhalb von Code-Blöcken, keine deklarierte Floskel und — opt-in über `placeholder` — kein unausgefüllter Vorlagen-Platzhalter. **Hermetisch** (kein git), fail-closed bei fehlender/mehrdeutiger Überschrift, fehlendem Closure-Verzeichnis und bei null Kandidaten. Überschriften und Marker zählen nur **außerhalb von Code-Blöcken**; die Block-Grenze ist die geteilte Abschnittsgrenze. **Dritte Fähigkeit** (opt-in über `waves.dir`): die Wellen-Register der Roadmap gegen die Wellen-Dateien — Plan-Dokument flach ⟺ aktive Welle (`waves.mode: one`, Default) **oder** Kennungs-Bijektion Aktiv-Block ⟺ flache Dokumente, Marker außen vor (`waves.mode: many`); Vorschau ohne Datei, Abschluss-Register ⟺ Ergebnisnotizen (beidseitig) | `planning-drift`, `closure-note-missing`, `closure-note-thin`, `closure-note-boilerplate`, `closure-note-placeholder`, `closure-note-ambiguous`, `wave-drift`, `wave-preview-exists`, `wave-results-missing`, `wave-unregistered` |
| `tracked`   | opt-in (git)  | Getrackt-Status auflösbarer, **existierender** Link-/Bild-Ziele gegen den git-**Index** (gestagt = getrackt, keine `.gitignore`-Interpretation); liest `.git` read-only, **ohne** Range; fail-closed ohne `.git` | `target-untracked`                                          |
| `targets`   | opt-in        | Deklarations-Konsistenz Doku ↔ Build-Targets: jedes in einer Doku-**Tabellenzeile** behauptete `make X` ist eine Makefile-Regel (`makefiles`), und jede Regel steht in der Autoritäts-Doku (`authority`); **hermetisch** (kein git, kein Makefile-Ausführen), fail-closed bei fehlender Datei. Tabellenzeilen zählen nur **außerhalb von Code-Blöcken** — ein Beispiel-Block dokumentiert kein Target | `gate-phantom`, `gate-undocumented`                         |
| `structure` | opt-in        | Struktur-Invarianten **innerhalb** eines Dokuments. Je Regel eine Dokumentklasse über **eigene** Globs (unabhängig vom Scan-Bereich, daher kein `scope`), ein Abschnitt (Klartext **oder** RE2) und bis zu zehn Bedingungen mit je eigenem Grund-Code — die siebte ist die **Chronologie-Monotonie** (`table.order`/`table.order-column`): typisierte Schlüsselspalte (ISO-Datum, Punkt-Version), rohe Zellen, nicht-strikt je zusammenhängender Tabelle; die achte die **Überschriften-Form** (`headings-match`/`headings-level`): **jede** Überschrift der geprüften Ebene innerhalb des Abschnitts matcht das Muster, geprüft auf ihrem **Text**, gemeldet **je** Überschrift auf **ihrer** Zeile; die neunte die **Zellenlänge** (`table.column` je Eintrag `name` mit `cell-max-chars`/`cell-min-chars`): jede Zelle einer über ihren **Kopfzeilen-Namen** benannten Spalte liegt in einer Spanne aus **Zeichen**, gemeldet auf **ihrer** Zeile — eine Obergrenze allein ließe die leere Zelle passieren; die zehnte die **offenen Task-Items** (`max-open-tasks`): gezählt auf den **rohen** Zeilen über die Listen-Lexik des Moduls, weil die absatzweise Inline-Code-Paarung eine Zusage über offene Haken sonst nicht trägt — ein Befund **je** Item auf **seiner** Zeile, Fence außen vor, Inline-Code nicht. `sections: one` (Default) erwartet genau einen Treffer, `each` prüft jeden. **Hermetisch** (kein git), fail-closed bei leerer Kandidaten-Menge — auch wenn erst `exempt-paths` sie geleert hat | `section-missing`, `section-ambiguous`, `section-empty`, `section-thin`, `section-oversized`, `section-forbidden`, `section-pattern-missing`, `section-marker-missing`, `section-unordered`, `section-cell-untyped`, `section-heading-mismatch`, `section-cell-oversized`, `section-cell-undersized`, `section-column-missing`, `section-tasks-open` |
| `workflows` | opt-in        | Deklarations-Konsistenz der `uses:`-Referenzen von CI-Workflows unterhalb eines **konfigurierten** Verzeichnisses (`workflows.dir` — der Ort ist nicht verdrahtet, weil er CI-System-spezifisch ist). **Fremde** Referenz: voller 40-stelliger Commit-SHA plus Tag-Kommentar dahinter — ein Tag lässt sich umhängen, ein SHA nicht; geprüft wird die **Form**, nicht die **Gültigkeit** (das wäre Netz). **Lokale** Referenz (`./…`): kein Pin nötig — sie löst auf denselben Commit auf wie ihr Aufrufer —, dafür zwei andere Fragen: **existiert** das Ziel, und **bekommt es die Rechte, die es verlangt?** Ein aufgerufener Workflow erhält nur, was der aufrufende **Job** selbst führt; ein Job ohne eigenes `permissions:` erbt zwar den Workflow-Kopf, kann aber nichts weitergeben, was er nicht deklariert. Stufen: `none` < `read` < `write`; ein nicht genannter Scope ist `none`, `read-all`/`write-all` setzen jeden. Die Referenzen kommen aus dem **YAML-Baum**, nicht aus einer Textsuche. **Hermetisch** (kein git, kein Netz, kein Ausführen), fail-closed bei leerer Prüfmenge und bei unlesbarem YAML. **Benannte Grenze:** das Modul liest die **Ziele** lokaler Referenzen, die es nicht scannt — dieselbe Parse-Zusage gilt dort —, und es deckt **eine** Deklarations-Klasse, nicht die Lauffähigkeit | `uses-pin-missing`, `uses-pin-untagged`, `uses-local-missing`, `uses-local-perms-undeclared`, `uses-local-perms-narrow`, `workflow-unparsable` |
| `external`  | opt-in (Netz) | Erreichbarkeit externer Links                                                            | `external-status`, `external-timeout`, `external-redirects` |
| `sources`   | opt-in (Netz) | Content-Pin externer Quellen gegen Upstream-Drift: eine auf `sha256` gepinnte `http(s)`-Quelle (Marker `<!-- source-pin: [zip] sha256:… -->` am Link **oder** Config-Block `sources:`; Einzeldatei oder Archiv `unpack: zip`) wird geholt, gehasht, verglichen — Meldung mit vollem Ist-Hash; **zweite** Netz-Tür neben `external`, nie im Default | `source-drift`, `source-unreachable`                        |

## 7. Fehlerbehebung

### Exit-Codes

| Code | Bedeutung                                                                        |
| ---- | -------------------------------------------------------------------------------- |
| `0`  | Prüfung gelaufen, keine Befunde                                                  |
| `1`  | Prüfung gelaufen, mindestens ein Befund                                          |
| `2`  | Nutzungs- oder Umgebungsfehler — die Prüfung lieferte keine verlässliche Aussage |

### Häufige Befunde

**`target-missing` — Linkziel existiert nicht.**
Ursache: Datei verschoben, umbenannt oder Tippfehler im Link.
Lösung: Pfad korrigieren. `--repair-broad` liefert einen Best-Guess-Patch
(review-pflichtig) **nur für Verschiebungen** — eine Markdown-Datei
gleichen, im Repo eindeutigen Namens an neuem Ort; bei **Umbenennung**
(anderer Name) gibt es keinen Vorschlag. `--doctor` zeigt die Diagnose
(für `target-missing` ohne Fix-Kandidat).

**`anchor-missing` — Anker entspricht keinem Heading.**
Ursache: Überschrift umbenannt oder falscher `#anker`.
Lösung: Anker an den (kleingeschriebenen, mit Bindestrichen
verbundenen) Heading-Slug anpassen.

**`id-unlinked` — Kennung ohne Link auf ihre Definition.**
Ursache: nackte Kennung im Fließtext.
Lösung: als Markdown-Link ausführen; `--repair` erzeugt den Fix
automatisch.

**`fence-unclosed` — Code-Block ohne Schluss bis zum Dateiende.**
Ursache: eine mit ``` oder ~~~ geöffnete Zeile findet keinen passenden
Schluss. Das ist mehr als ein Schönheitsfehler: **alles dahinter gilt als
Code**, und die Prüfungen, die darauf aufsetzen, überspringen es. Ein
kaputter Link hinter einem offenen Block wird nicht gemeldet — die Datei
sieht sauber aus, obwohl sie nicht geprüft wurde.

Lösung: den Block schließen. **Die gemeldete Zeile ist eine Fundstelle,
nicht zwingend die Reparaturstelle** — welche von mehreren Öffnungen den
Schluss vermissen lässt, ist grundsätzlich nicht entscheidbar. Sind alle
Blöcke gleich lang, zeigt der Befund auf die letzte Zeile, an der die
Zählung kippte; die fehlende Klammer kann davor liegen. Praktisch: ab der
gemeldeten Zeile rückwärts lesen, bis die Paare aufgehen.

Zwei Fälle sehen aus wie ein Fehlalarm und sind keiner: ein
Backtick-Block, den eine Tilden-Zeile „schließt" (verschiedene Zeichen
schließen einander nicht), und ein längerer Block, der eine kürzere
Fence-Zeile zeigt (die Zählung kippt, auch wenn CommonMark den Block für
geschlossen hält). In beiden Fällen liest mindestens eine Prüfung den
Rest der Datei als Code — der Befund ist berechtigt.

### Häufige Fehler (Exit-Code 2)

**„Scan-Wurzel nicht gefunden … wurde das Repository nach /repo
gemountet?"**
Ursache: kein (oder falscher) Mount.
Lösung: `-v "$PWD:/repo:ro"` ergänzen.

**Ungültige Konfiguration (mit Zeilenangabe).**
Ursache: Tippfehler oder unbekannter Schlüssel in `.d-check.yml`.
Lösung: gegen das Gerüst aus `--print-config` abgleichen.

**„… ergab 0 Anforderungen" — die konfigurierte RTM-Quelle ist leer.**
Ursache: `trace.requirements.source` ist explizit (nichtleer) gesetzt oder
`format: table` aktiv, aber d-check erkennt null Anforderungen. `--trace` endet
dann fail-closed mit Exit 2, statt eine leere Matrix zu behaupten:

```bash
# DCHECK_IMAGE auf v0.43.1 oder neuer pinnen.
docker run --rm -v "$PWD:/repo:ro" "$DCHECK_IMAGE" --trace --require-complete
```

```text
d-check: error: trace.requirements: Quelle "spec/lastenheft.md" im Format headings ergab 0 Anforderungen
```

Lösung: Quelldatei, `format` und `id-pattern` prüfen. `source: ""` gilt wie ein
abwesendes Feld und fällt auf die Default-Quelle zurück (leere Matrix, Exit 0).

## 8. FAQ

**Verändert d-check meine Dateien?**
Nein. d-check ist ein Lese-Werkzeug; selbst `--repair` gibt nur einen
Patch aus.

**Braucht d-check Netzwerk?**
Nein — außer eines der Netz-Module `external` oder `sources` ist aktiv (beide
opt-in). Sonst läuft es vollständig netzlos (`--network none`).

**Warum meldet d-check nichts in meinem `README`?**
Ohne Konfiguration werden `docs/`, `spec/` und Wurzel-`*.md` geprüft.
Andere Verzeichnisse nehmen Sie über `scan.roots` auf.

**Wie schalte ich eine einzelne Zeile von der `ids`-Prüfung frei?**
Setzen Sie den Marker `d-check:ignore` als HTML-Kommentar an die Zeile. Denselben
Marker honorieren `codepaths`, `versions` und `diagrams` — in `diagrams` als
**Token** (in einem Fence gibt es keine HTML-Kommentare), dort auch auf der
Öffnungszeile für den ganzen Block. Die übrigen Module kennen ihn nicht.

**Ist die Ausgabe stabil/wiederholbar?**
Ja. Gleiche Eingabe liefert byte-identische Ausgabe (stabile Sortierung).

## 9. Glossar

- **Befund:** eine einzelne Beanstandung (`Datei:Zeile  Ziel  Grund-Code`, plus
  eine vierte Spalte, wo eine Erläuterung vorliegt).
- **Grund-Code:** stabiler, maschinenlesbarer Code eines Befunds (z. B.
  `target-missing`).
- **Regelmodul:** ein zuschaltbarer Prüfbereich (z. B. `links`, `ids`).
- **Anker:** das `#fragment` eines Links, das auf eine Überschrift zeigt.
- **Kennung:** eine stabile ID (z. B. eine ADR- oder Anforderungs-Nummer).
- **Read-only-Mount:** ein Docker-Mount mit `:ro`; d-check braucht nicht
  mehr.

## 10. Support und Lizenz

- **Quellcode, Issues, Releases:** <https://github.com/pt9912/d-check>
- **Lizenz:** Code unter MIT, Texte unter CC BY 4.0.
- **Sicherheit:** d-check schreibt nie ins Repository und öffnet außer in
  den Modulen `external` und `sources` keine Netzwerkverbindungen; übergeben
  Sie keine Geheimnisse als Argumente.

## 11. Änderungshistorie

Die Software-Änderungen je Version stehen im
[Changelog](../../CHANGELOG.md). Dieses Handbuch ist an die
Software-Version gekoppelt und wird mit den Releases fortgeschrieben.

| Handbuch-Version | Software-Version | Stand      | Änderung                                                                                                                                                                                |
| ---------------- | ---------------- | ---------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1.0              | v0.12.0          | 2026-06-18 | Erstfassung: alle Use Cases inkl. `--doctor`/`--repair`                                                                                                                                 |
| 1.1              | v0.17.0          | 2026-06-19 | `--doctor --json`: maschinenlesbare Diagnose ergänzt (§4.9), Gegenüberstellung der drei Ausgaben                                                                                        |
| 1.2              | v0.18.0          | 2026-06-19 | `--suggest-config ai-harness`/`ai-harness-init`: Harness-Vorlage in zwei Modi ergänzt (§4.4)                                                                                            |
| 1.3              | v0.19.0          | 2026-06-20 | YAML-Ausgabe `--yaml` (auch `--doctor --yaml`): maschinenlesbare Ausgabe um YAML erweitert (§4.11)                                                                                      |
| 1.4              | v0.22.0          | 2026-06-22 | `--id-prefix` für `--suggest-config` (§4.4); Traceability-Matrix `--trace` (§4.12); Makefile-Fragment `--print-mk` (§4.16); `:latest` = neueste **stabile** Version (§2)                |
| 1.5              | v0.23.0          | 2026-06-22 | Modul `codepaths`: Datei-Ventil `exempt-paths` ergänzt (§5 Weitere Module) — ganze Dateien von der Inline-Code-Pfad-Prüfung ausnehmen, wie `ids`                                        |
| 1.6              | v0.24.0          | 2026-06-23 | `--print-mk`-Fragment um `doc-trace`/`doc-complete`-Targets + `TRACE_FLAGS` erweitert (§4.16); opt-in `--trace --require-complete` (Vollständigkeits-Gate, Waise ⇒ Exit 1, §4.12)       |
| 1.7              | v0.25.0          | 2026-06-23 | Modul `diagrams` (opt-in): Kennungs-Existenz in Diagramm-Fences (§5 Weitere Module) — `mermaid`-Diagramme auf undefinierte Kennungen prüfen (Befund `diagram-id-undefined`)             |
| 1.8              | v0.26.0          | 2026-06-23 | `--suggest-config ai-harness[-init]`: Kommentar-Hinweis auf die nicht aktivierten situativen opt-in-Module (`external`/`spans`/`hostpaths`/`diagrams`) mit Verweis auf `--print-config` |
| 1.9              | v0.27.0          | 2026-06-23 | `--print-mk`-Fragment (§4.16) um `doc-doctor`/`doc-repair`/`doc-help`-Targets + `DCHECK_DIGEST` (Digest-Override, sticht den Tag) erweitert; alle Targets `##`-annotiert                |
| 1.10             | v0.28.0          | 2026-06-24 | Modul `versions` (opt-in): Versions-Pin-Konsistenz (§5 Weitere Module, §6) — gepinnte `ghcr`-Image-Verweise müssen die aktuelle Version (aus `version.md#aktuell`) tragen, auch in Fences; Befund `version-stale`                |
| 1.11             | v0.29.0          | 2026-06-24 | Modul `pins` (opt-in): Content-Pin gegen inhaltlichen Drift (§6) — ein Link mit `<!-- dpin: sha256:… -->` wird gegen den Hash seines (whitespace-normalisierten, rohen) Ziel-Spans geprüft; Drift → `link-stale`                |
| 1.12             | v0.30.0          | 2026-06-28 | Modul `matrix`: klasseninterne Verweisrichtung (§4.7) — eine Klasse mit `order` (Glob-Rang, autoritativste Schicht zuerst) + `direction: no-downward` meldet klasseninterne Abwärtsverweise als `matrix-downward`                |
| 1.13             | v0.31.0          | 2026-06-28 | Modul `matrix`: token-basierte Referenz-Richtung (§4.7) — eine Klasse mit `token`-Regex fängt verbotene Referenzen auch als bare ID-Token im Körper (`matrix-forbidden`); Provenance-Marker `<!-- d-check:status-provenance -->` nimmt aus, `exempt-paths` grandfathered ganze Dateien                |
| 1.14             | v0.32.0          | 2026-06-28 | Neues opt-in-Modul `immutable` (§6): Immutabilitäts-Pin gegen Core-Drift — eine Datei mit `<!-- immutable: sha256:… -->` wird gegen den whitespace-normalisierten **Core** (Datei ohne Marker-Zeile + `immutable.exclude-sections`) gehasht; Abweichung → `core-drift`. Hermetisch (kein git, read-only-Arbeitsbaum), diagnose-only, opt-in pro Datei                |
| 1.15             | v0.33.0          | 2026-06-29 | Neues opt-in-Modul `vcs` (13., §5/§6): git-Diff-Immutabilität des Core über eine Commit-Range (`--range`/`--staged`, `core-drift-vcs`) — reine-Go-git im read-only `.git` (distroless bleibt, fail-closed, diagnose-only); löst die Skript-Mechanik des `adr-check`-Gates ab. `--print-mk`-Target `doc-immutable` verteilt die git-Garantie an Konsumenten                |
| 1.16             | v0.34.0          | 2026-06-29 | Modul `codepaths`: Referenz-Ventil `ignore-refs` (§5/§6) — bestimmte **aufgelöste Ziel-Pfade** referenz-weit (datei-/zeilen-unabhängig) von der Existenz-/Escape-/Anker-Prüfung ausnehmen, als Tombstone-Register bewusst entfernter Artefakte; löst die Frozen-Doc-Refactoring-Falle. Dritte Ventil-Achse neben `d-check:ignore` (Zeile) und `exempt-paths` (Datei)                |
| 1.17             | v0.35.0          | 2026-07-01 | Neues opt-in-Modul `commits` (14., §5/§6): Traceability-Kennung (`id-patterns`) in jeder Commit-Message über eine Range (`--range`) bzw. der Pending-Message (`--commit-msg`, commit-msg-Hook via stdin); Befund `commit-untraceable`, `exempt-pattern` nimmt Merge/Revert-Betreffe aus — reine-Go-git im read-only `.git` (distroless bleibt, fail-closed, diagnose-only); löst die Skript-Mechanik des `trace-check`-Gates ab. `--print-mk`-Target `doc-commits` verteilt die Range-Prüfung an Konsumenten                |
| 1.18             | v0.36.0          | 2026-07-01 | Neues opt-in-Modul `planning` (15., §5/§6): Roadmap-↔-in-progress-Lifecycle-Konsistenz — der Ruhe-Marker steht im `## Aktuelle Welle`-Block genau dann, wenn kein `slice-*` im Verzeichnis liegt (`planning-drift`); **hermetisch** (nur Roadmap + Verzeichnis-Listing, kein git), fail-closed bei fehlender/mehrdeutiger Überschrift, diagnose-only. Löst die Skript-Mechanik des `planning-check`-Gates ab (letztes Familien-Skript). `--print-mk`-Target `doc-planning` (hermetisch, ohne Range) verteilt die Prüfung                |
| 1.19             | v0.37.0          | 2026-07-03 | Neues opt-in-Modul `tracked` (16., §5/§6): Getrackt-Status auflösbarer, **existierender** Link-/Bild-Ziele gegen den git-**Index** (`target-untracked` — untracked/gitignoriertes Ziel fehlt auf jedem frischen Klon); Index-Wahrheit (gestagt = getrackt), kein Doppelbefund, Ventil `exempt-targets` (aufgelöster Zielpfad, segmentweise validiert), fail-closed ohne `.git`, **ohne** Commit-Range. `--print-mk`-Target `doc-tracked` (zehn Targets, §4.16) verteilt die Prüfung                |
| 1.20             | v0.37.1          | 2026-07-04 | Fix: `--doctor` zeigt für die sieben seit v0.25 hinzugekommenen Grund-Codes (`diagram-id-undefined`, `version-stale`, `link-stale`, `core-drift`, `core-drift-vcs`, `commit-untraceable`, `planning-drift`) den Klartext statt des rohen Codes — auch als `reasonText` in `--doctor --json`/`--yaml`; die Klartext-Liste ist testseitig beidseitig gegen die Grund-Code-Tabelle der Spezifikation verriegelt (fail-closed)                |
| 1.21             | v0.37.1          | 2026-07-04 | Anleitung „Das Versions-Register `version.md` aufbauen" (§5) — Aufbau (`## Aktuell`/`## Verlauf`), die `<a id="vX.Y.Z">`-Anker-Mechanik samt Anker-Wanderung beim Release und ein kopierbares Muster zum Nachbau in eigenen Repos                |
| 1.22             | v0.38.0          | 2026-07-05 | Neues opt-in-Modul `targets` (17., §5/§6): Deklarations-Konsistenz Doku ↔ Build-Targets — jedes in einer Doku-**Tabellenzeile** behauptete `make X` ist eine Makefile-Regel (`gate-phantom`), jede Regel steht in der Autoritäts-Doku (`gate-undocumented`); **hermetisch** (kein git, kein Makefile-Ausführen), fail-closed. Löst den Doku-↔-Makefile-Kern des `gate-consistency.sh`-Meta-Gates ab; `--print-mk`-Target `doc-targets` (elf Targets, §4.16) verteilt die Prüfung                |
| 1.23             | v0.39.0          | 2026-07-06 | `--suggest-config ai-harness[-init]` (§4.4): Vorlage an die **gelebte** Konvention angeglichen — `spans`/`hostpaths` ins fixe Standard-Modulset, repo-bewusster `planning`-Block; `vcs`/`commits` (Commit-Range) via `--print-mk`, `versions`/`targets` bewusst vertagt; kanonische Vorlage der Spezifikation deckt die emittierte Ausgabe 1:1                |
| 1.24             | v0.40.0          | 2026-07-11 | Requirements Traceability Matrix (`--trace`, §4.12) über einen opt-in `trace`-Block quell-/kennungs-konfigurierbar (§5): Anforderungs-Quelldatei + Kennungs-Regex sowie je Referenzklasse Verzeichnis + Dateimuster + Owner-Präfix; Default = Konvention ⇒ byte-identisch, fail-closed bei ungültiger Regex / Muster ohne Capture-Gruppe. Passt diese Quellachsen an abweichende Repo-Konventionen an; die Heading-basierte Definitionssyntax bleibt bestehen                |
| 1.25             | v0.41.0          | 2026-07-11 | Kuratierte Coverage-Quellen `trace.coverage` (§4.13/§5): eine Liste benannter `files` liest Deckungs-Matrizen als **eigene Coverage-Spalte** ein — `ranges` (`GG-QA-001..006` → alle sechs, `/`-Aufzählung), `sections`/`exclude-sections` (voller Heading-Text). Waise = ohne Slice **und** ohne Coverage; fail-closed (fehlende Datei / leeres label / Sektion ohne Treffer / ungültige Range ⇒ Exit 2); ohne Block byte-identisch                |
| 1.26             | v0.42.0          | 2026-07-11 | Modalitäts-Klassifikation `trace.requirements.modality` (§4.14/§5): aus konfigurierbaren Modalverb-Stichwörtern (DE+EN-Defaults, opt-in) klassifiziert die RTM jede Anforderung als MUSS/SOLLTE/KANN in einer **eigenen Modalitäts-Spalte** (längster Treffer, Wortgrenze, Markup-normalisiert); optional gatet `require-levels`, welche Stufen einen Waisen zum Exit-1-Fehler machen (Default: nur MUSS). Fail-closed bei leerem Level/Stichwort, reserviertem `unknown`, Stichwort in zwei Stufen, ungültigem `require-levels` ⇒ Exit 2; ohne Block byte-identisch                |
| 1.27             | v0.42.0          | 2026-07-14 | Trace-Präzisierung (§4.12/§5): Requirement-Definition nur als ATX-Heading mit ID im ersten Token; Tabellen-/Listen-/Fließtextgrenze, Body-only-Modalität, exakte Waisen- und Referenzscan-Semantik, Warnung vor leerer RTM mit Exit 0 sowie Brownfield-Migration für tabellarische Lastenhefte dokumentiert                |
| 1.28             | v0.43.1          | 2026-07-14 | Native Trace-Tabellenquellen (§4.12/§5): `trace.requirements.format: table` bindet ID-, alternative Text- und optionale Modalitätsspalten über exakte Header-Namen; Pipe-Escaping, deterministische Duplikatpolitiken und gemeinsames RTM-Modell. Nichtleer explizite Quelle oder Tabellenmodus endet bei fehlender Quelle, ungültiger Tabellenstruktur oder null Anforderungen fail-closed mit Exit 2; unkonfigurierter Heading-Default bleibt kompatibel                |
| 1.29             | v0.43.1          | 2026-07-15 | Keine inhaltliche Änderung — Software-Version auf v0.43.1 nachgezogen (die Zeile 1.28 wurde dabei von v0.43.0 auf v0.43.1 umgebogen). Zeile nachgetragen: der Kopf trug 1.29 seit dem v0.43.1-Release-Prep ohne eigenen Historien-Eintrag (Prosa-Currency-Lücke, kein Gate erfasst sie — siehe [releasing.md](releasing.md) §Release-Prep)                                        |
| 1.30             | v0.44.0          | 2026-07-17 | Kreuzverweis-Konsistenz `trace.cross-consistency` (§4.15/§5): der `--trace`-Lauf vergleicht die Vorwärts-RTM-Tabelle (Anforderung → Design) gegen die Rückwärts-`Bezug`-Kanten (Design → Anforderung) und meldet je Anforderung beide Mengendifferenzen mit Richtung und `Datei:Zeile`; Modi `equal`/`superset`, `exclude-req`-Ventil für Ableitungssprünge, `artifact-id-column: first` für heterogene ID-Header. Advisory unter `--trace`, gatend allein über `--require-complete`; fail-closed inkl. vakuumem Abgleich; ohne Block RTM byte-identisch |
| 1.31             | v0.44.1          | 2026-07-17 | Range-Notation unter Linkpflicht (§5): eine verlinkte Range/Enum-Fortsetzung (`` [`GG-QA-001`](…)..006 ``) wird wie die unverlinkte gelesen — d-check überspringt genau ein Link-Suffix. Bis v0.44.0 expandierte sie gar nicht und erzeugte in `trace.coverage` falsche Waisen; die eng gefasste Grenze (kein Whitespace/Emphasis/zweites Suffix) ist dokumentiert                          |
| 1.32             | v0.45.0          | 2026-07-17 | Scope-Trennung des Kreuzverweis-Abgleichs (§5): `forward.req-pattern` (Default `requirements.id-pattern`) erkennt die Anforderungs-IDs der Vorwärts-Sicht — welche Anforderungen verglichen werden, entscheidet das Muster, nicht die RTM-Mitgliedschaft. Die bis v0.44.1 stille Kopplung ist als **Scope-Falle** dokumentiert: bei bewusst gescopter RTM meldete sie jede Rück-Kante als Falschbefund, während die echte Gegenrichtung verschwand                                                    |
| 1.33             | v0.45.1          | 2026-07-17 | Richtigstellung zur Range-Notation unter Linkpflicht (§5): Klammern **im Linkziel** sind unproblematisch — das Ziel wird klammer-balanciert abgegrenzt wie bei `links`/`ids`. Die Fassung 1.31 behauptete das Gegenteil; tatsächlich expandierte v0.44.1/v0.45.0 dort Pfadsegmente als Enum und versteckte Waisen                                                                                             |
| 1.34             | v0.46.0          | 2026-07-17 | Komma-Kurzform in Coverage-Quellen dokumentiert (§5): `GG-SCN-001, 007` — und `GG-SCN-001..005, 007, 008` (hinter einer Range) — ist keine zugesagte Notation und bricht mit Exit 2 ab, statt `007` still fallen zu lassen. Zugesagt bleiben nur `..BBB` und `/BBB`; ein Komma **vor** einer vollständigen Kennung ist erlaubt                                                                                       |
| 1.35             | v0.47.0          | 2026-07-18 | Markdown-Lexik an CommonMark/GFM angeglichen (Grundkonzepte, §5): eine Tabellen-Trennzeile braucht nur **einen** Bindestrich (`\| - \|`), und eine Backtick-Fence-Zeile mit Backtick in der Infozeile ist Fließtext, kein Block-Öffner — d-check findet danach **mehr** (SemVer-Minor, ein grüner Konsumentenlauf kann rot werden)                                                                                       |
| 1.36             | v0.48.0          | 2026-07-18 | Tabellengrenze am relevanten Header (§5): eine irrelevante Tabelle verschluckt die unmittelbar folgende **relevante** nicht mehr still — ein Header, der eine konfigurierte Rolle bindet, beendet die laufende Tabelle. d-check findet danach **mehr** (SemVer-Minor, ein grüner Konsumentenlauf kann rot werden)                                                                                       |
| 1.37             | v0.48.1          | 2026-07-18 | Direktiven-Toleranz in Tabellenzeilen (§5): eine Datenzeile mit genau einer überzähligen, ganzzelligen HTML-Kommentar-Zelle (`<!-- d-check:ignore … -->` hinter der letzten Pipe) wird auf Header-Breite gelesen, statt fail-closed mit Exit 2 abzubrechen. Zwei Extra-Zellen oder Nicht-Kommentar bleiben Exit 2. Patch (rot→grün)                                                                                       |
| 1.38             | v0.49.0          | 2026-07-18 | Geteiltes Referenz-Ventil `ignore-refs` (§5): das bisher modul-lokale `codepaths.ignore-refs` wird zur **querschnittlichen** Top-Level-Fähigkeit, die `links`/`anchors`/`codepaths` gemeinsam honorieren. Neue Felder je Eintrag: `in` (Quell-Skopus, Glob auf die Quelldatei), `refs` (aufgelöste Ziele) und `keep` (Ausnahmen, reihenfolge-unabhängig). Ziel-Achsen-Pendant zu `scan.ignore`; löst die Template-Verzeichnis-Falle. `codepaths.ignore-refs` bleibt Alias (kein Config-Bruch), ohne Block byte-identisch; ungültiges Glob ⇒ Exit 2. Die vier Ventil-Achsen sind jetzt gegeneinander erklärt                                                                    |
| 1.39             | v0.50.0          | 2026-07-18 | Zitat-Verifikation (§5/§6): opt-in `codepaths.check-lines` verifiziert `datei:<von>-<bis>`-Zeilen-Referenzen (`citation-out-of-range`/`citation-inverted-range`), Default aus byte-identisch. Neues 18. Modul `citations` (opt-in): die Direktive `<!-- d-check:cite <pfad>:<von>-<bis> -->` markiert das folgende Zitat (`>`-Block oder inline `„…"`/`"…"`), das ein whitespace-normalisierter Teilstring der Quell-Spanne sein muss (`citation-mismatch`); Mindestlänge 16 Zeichen, malformte Direktive fail-closed. d-check findet danach **mehr** (SemVer-Minor)                                                                    |
| 1.40             | v0.51.0          | 2026-07-19 | Neues opt-in-Modul `sources` (19., **Netz**, §5/§6): Content-Pin externer Quellen gegen Upstream-Drift — eine auf einen `sha256` gepinnte `http(s)`-Quelle (Marker `<!-- source-pin: [zip] sha256:… -->` am Link **oder** Config-Block `sources:`) wird geholt, gehasht und verglichen; Abweichung → `source-drift` (Meldung mit vollem Ist-Hash zum Re-Pinnen), Netzfehler/HTTP ≥ 400/Timeout/Größenlimit/kein gültiges Zip → `source-unreachable`. Einzeldatei (Roh-Byte-Hash) oder Archiv (`unpack: zip`, reihenfolge-invariantes Content-Manifest). **Zweite Netz-Tür** neben `external` (beide opt-in, nie im Default-Lauf); malformte Direktive/ungültige Config fail-closed (Exit 2), ohne aktives `sources` byte-identisch                |
| 1.41             | v0.51.1          | 2026-07-19 | dpin-Ergonomie (§5, Modul `pins`): der `link-stale`-Befund führt jetzt den **vollen** errechneten `sha256` (statt nur `shortHash`) — damit ist dpin „einmal laufen, gemeldeten Hash in den `<!-- dpin: … -->`-Marker kopieren"-benutzbar. Neue Aufgaben-Sektion §5 „Einen repo-internen Link-Inhalt gegen Drift pinnen". Nur die nicht stabilitätsgarantierte Befund-Meldung, kein Verhaltens-/Grund-Code-Delta                |
| 1.42             | v0.51.1          | 2026-07-19 | §4 aufgabenorientiert nachgezogen (Benutzerhandbuch-Standard §2/§5): §4.7 (`order`/`direction`/`token`) auf die Lesersituation; der §4.12-`--trace`-Monolith in §4.12 RTM · §4.13 Coverage (`trace.coverage`) · §4.14 Modalität (`trace.requirements.modality`) · §4.15 Kreuzverweis (`trace.cross-consistency`) aufgetrennt, `--print-mk` → §4.16. Tabellen-Grammatik/Modalitäts-Mechanik → §5, „0 Anforderungen"-Fehlerbild → §7, verstreute Versions-Prosa entfernt (§11 führt die Historie), doppelte WAISE-Definition raus. Rein redaktionell, kein Verhaltens-/Grund-Code-Delta; beide Handbuch-Harnesse grün                |
| 1.43             | v0.52.0          | 2026-08-09 | **Closure-Notizen auf Substanz prüfen** — neue Aufgabe §4.17 (Modul `planning`, zweite Fähigkeit): opt-in über `closure.dir` wird je abgeschlossenem Paket der erste `heading-pattern`-Abschnitt strukturell geprüft — vorhanden (`closure-note-missing`), genug Satzende-Zeichen **außerhalb** von Code-Blöcken (`closure-note-thin`, `min-sentences` Default 4), keine deklarierte Floskel (`closure-note-boilerplate`, Liste Default leer). Fail-closed auch bei **null Kandidaten**; ohne `closure.dir` inert. Zugesagt ist **Struktur, nicht Bedeutung**. Dazu die neue Option **`--config <datei>`** (§5): sie ersetzt für den Lauf die konventionelle `.d-check.yml` (innerhalb der Scan-Wurzel, ersetzt statt ergänzt) und trennt damit Prüf-Profile — fehlende Datei, Verzeichnis-Ziel, Symlink-Ausbruch und **leerer Wert** brechen mit Exit 2 ab, kein stiller Rückfall. §5 und §6 nachgezogen |
| 1.44             | v0.53.0          | 2026-08-10 | **Unbalancierte Code-Blöcke melden** — dritte Artefakt-Klasse des Moduls `spans` (§6/§7): eine Fence-Öffnung ohne Schluss bis zum **Dateiende** meldet `fence-unclosed`. Anlass war ein stiller Grün-Pfad — hinter einem offenen Block übersprang d-check den Rest der Datei und meldete grün, ohne geprüft zu haben. Ausgewertet werden **beide** Schluss-Lesarten; ein Befund entsteht, sobald eine von beiden offen endet. Genau einer je Datei, an der Öffnungszeile — eine **Fundstelle**, nicht zwingend die Reparaturstelle. Kein neuer Schalter (`spans` ist als Ganzes opt-in). **Beachten:** die Fence-Erkennung trimmt jetzt überall gleich; am Closure-Gate des Moduls `planning` können dadurch Befunde wegfallen, die auf der abweichenden Trimmung beruhten |
| 1.45             | v0.54.0          | 2026-08-10 | **Eigener Kandidaten-Filter für Closure-Notizen** — neuer Schlüssel `planning.closure.glob` (§5): er entscheidet, welche Dateien im Ruheort auf ihre Notiz geprüft werden, getrennt von `slice-glob`. Weglassen ⇒ es gilt `slice-glob` (**Verweis**, kein kopierter Wert). Nötig wurde er, weil die beiden `planning`-Fähigkeiten verschiedene Fragen stellen und sich einen Knopf teilten: wer ihn weitet, verbiegt die jeweils andere — im Grenzfall zählt die Roadmap-Datei selbst als „Arbeit in Arbeit“. Explizit leerer/ungültiger Glob ⇒ Exit 2, geprüft auch bei inerter Fähigkeit. **Beachten:** zwei Meldungstexte haben sich geändert (der Glob steht in der Nullmengen-Meldung jetzt in Anführungszeichen, und der `--doctor`-Klartext sagt „Kandidat“ statt „Slice“); der Befundsatz selbst ist unverändert |
| 1.46             | v0.55.0          | 2026-08-10 | **Unausgefüllte Vorlagen-Platzhalter melden** — vierte, per Default **ausgeschaltete** Bedingung der Closure-Note-Struktur (§4.17/§5): `placeholder: true` meldet `closure-note-placeholder`, wenn eine Notiz noch einen Rumpf wie „Ergebnis: Feld-in-Winkelklammern“ trägt. Die drei bestehenden Bedingungen sehen ihn nicht: er ist syntaktisch vollständig. Erkannt wird die Auszeichnungs-Form, deren Inneres **kein Whitespace** enthält; **Inline-Code zählt nicht**, und Autolinks/Adressen, HTML-Tags sowie Winkelklammer-Linkziele sind ausgenommen. Erster Treffer je Notiz, an seiner Zeile. Zwei Grenzen sind benannt: eingerückte Code-Blöcke und ungerade Backtick-Parität |
| 1.47             | v0.56.0          | 2026-08-10 | **Was als Satz und als Floskel zählt** (§4.17): der Closure-Abschnitt wird **einmal** bereinigt (Code-Blöcke **und** Inline-Code), ein Satzende zählt nur vor Leerzeichen oder Zeilenende, und Floskeln werden an **Wortgrenzen** verglichen. **Drei Richtungen:** `closure-note-thin` findet **mehr** (Punkte aus Pfaden und Versionen sind keine Sätze mehr — ein grüner Lauf kann rot werden), `closure-note-boilerplate` findet **weniger** (Floskel in Backticks trifft nicht mehr; kurze Phrasen dafür überhaupt erst benutzbar). Anlass ist die Parität zu einem abzulösenden Prüfskript, belegt über 84 fremde Closure-Notizen. Behoben: bei CRLF zählte kein zeilenschließendes Satzende |
| 1.48             | v0.57.0          | 2026-08-15 | **Struktur-Invarianten prüfen** — neues opt-in-Modul `structure` (20., §5/§6): je Regel eine Dokumentklasse über **eigene** Globs, ein Abschnitt (Klartext oder RE2) und bis zu sechs Bedingungen mit je eigenem Grund-Code (leer · zu dünn · zu viele Task-Items · verbotenes bzw. gefordertes Muster · fehlende Marke). Fehlender Abschnitt und Mehrdeutigkeit sind eigene Codes; eine Regel, die **keine** Datei trifft, meldet ebenfalls — auch wenn erst `exempt-paths` die Menge geleert hat. Die Closure-Note-Struktur des Moduls `planning` ist ein **Preset** derselben Semantik und teilt die Mechanik. Dazu `closure-note-ambiguous` (mehrere Closure-Überschriften) und ein zwölftes `--print-mk`-Target `doc-structure` |
| 1.49             | v0.58.0          | 2026-08-16 | **Geteilte Lexik: fünf Module beantworten ihre Lexik-Fragen jetzt gleich** (§5/§6). Wer fragt „ist das eine Überschrift / ein Anker / derselbe Absatz“, bekommt überall dieselbe Antwort: `citations` trennt am **Code-Block** (Leerzeilen trennen weiterhin nicht) und bricht dann fail-closed ab; `versions`/`pins` lösen Anker so auf wie `anchors` — Duplikat-Slug, Prozent-Dekodierung, case-sensitiv, nur außerhalb von Code; `planning` zählt Überschriften und den Ruhe-Marker nur außerhalb von Fences und beendet den Aktiv-Block an der geteilten Abschnittsgrenze; `vcs` liest Status-Zeile und `immutable-when` nicht mehr aus einem Code-Beispiel; `targets` wertet nur Tabellenzeilen außerhalb von Fences. **Beachten — die Änderung wirkt in beide Richtungen:** sie findet **mehr** (fail-closed statt Zufalls-Paarung, zwei geschlossene stille Grün-Pfade in Gates) und **weniger** (ein Anker im Fence oder anders geschrieben löst nicht mehr auf, bei `pins` ohne Ausgabe; zwei `planning`-Falsch-Rot entfallen). Die Aufzählung ist offen. Am eigenen Bestand ist kein Fall betroffen. Dazu eine benannte Grenze: die Abschnitts-Maske von `vcs` läuft auf git-Blobs, die kein Wächter scannt |
| 1.50             | v0.59.0          | 2026-08-16 | **Wellen-Register gegen Wellen-Dateien** — dritte Fähigkeit des Moduls `planning` (§5/§6, opt-in über `planning.waves.dir`): die Lifecycle-Invariante eine Ebene höher, vier Aussagen mit je eigenem Grund-Code (`wave-drift` · `wave-preview-exists` · `wave-results-missing` · `wave-unregistered`). Das verpflichtende Artefakt einer geschlossenen Welle ist die **Ergebnisnotiz** — gemessen, nicht gesetzt: gegen das Plan-Dokument geprüft meldete die Aussage 19-mal über zwei reale Bäume. Die Vorschau-Aussage liest die **erste Spalte** und überspringt Zeilen ohne Kennung; Kopf- und Trennzeile zählen strukturell nicht. Fail-closed: unlesbare Wellen-Verzeichnisse und fehlende Register-Überschriften (je eigenes Befund-Ziel gegen die Deduplikation). **Beachten:** beim Einführen den konsument-gerechten `marker` setzen und mit echten Rückständen rechnen (im Schwester-Repo: elf) |
| 1.51             | v0.60.0          | 2026-08-16 | **Ortsfeste Verweise in Lifecycle-Verzeichnissen** — opt-in `links.resolve-from` (§5/§6): wo Dateien per `git mv` zwischen Geschwister-Verzeichnissen wandern, muss jede Datei eines `dirs`-Verzeichnisses jedes relative Ziel von **jedem** Ort ihrer Gruppe auflösen, und überall auf **dasselbe** Ziel — sonst `link-position-dependent`, gemeldet **vor** dem Move; die Reparatur ist das Präfixieren des Pfads, nicht das Anlegen des Ziels. `fixed-dirs` (etwa der Ruheort) zählen als Orte, prüfen aber nicht; ein am Ist-Ort fehlendes Ziel bleibt allein `target-missing` (kein Doppelbefund); absichtlich ortsgebundene Verweise über das bestehende `ignore-refs`. Fail-closed über denselben Code: eine Gruppe ohne einen einzigen existierenden `dirs`-Ort und ein Ort, der als Datei existiert — ein **einzelner** fehlender Ort meldet bewusst nicht (git überträgt leere Verzeichnisse nicht). Ohne den Block byte-identisch |
| 1.52             | v0.61.0          | 2026-08-21 | **Chronologische Tabellen auf Richtung prüfen** — siebte Bedingung des Moduls `structure` (§5/§6): `table-order` (`asc`/`desc`) mit `table-column` (1-basiert, Default 1) vergleicht die Schlüsselspalte jeder zusammenhängenden Tabelle des Abschnitts **typisiert** (ISO-Datum, Punkt-Version — `1.10` kommt nach `1.9`; ein zeichenweiser Vergleich meldete korrekt sortierte Tabellen rot) und **nicht-strikt** monoton. Die Zellen werden **roh** gelesen (einzige Ausnahme von der Abschnitts-Bereinigung — Release-Register führen ihre Schlüssel in Inline-Code); Kopf-/Trennzeile zählen nicht. Neue Grund-Codes `section-unordered` (Bruch-Zeile; auch Leerlauf ohne Datenzeile) und `section-cell-untyped` (untypisierbare Zelle/Typ-Mischung — Befund statt stillem Übersprung, Anker-Reset dahinter). Drei neue Exit-2-Config-Ränder. Ohne `table-order` byte-identisch |
| 1.53             | v0.62.0          | 2026-08-21 | **Kennungs-Bijektion für offene Wellen** — `planning.waves.mode: one`\|`many` (§5/§6, Lastenheft 0.62.0 auf **formalen Konsumenten-CR**): unter `many` vergleicht `wave-drift` die im `heading`-Block genannten Wellen-**Kennungen** gegen die flachen Wellendokumente (beide Richtungen, `target` = Kennung, Fences zählen nicht) statt Marker gegen Datei-Zahl — der Ruhe-Marker bleibt Sache von `planning-drift`, und der Zustand „Welle eröffnet, nichts beansprucht" (Marker **und** Zeiger) ist grün. `one` (Default) unverändert, ohne Schlüssel byte-identisch; unbekannter/explizit leerer Modus ⇒ Exit 2. **Beachten:** unter `many` zählt jede Kennung in der Abschnitts-Prosa als Zeiger — erklärenden Text dort kennungsfrei halten. Nebenbei präzisiert: die §6-planning-Zeile nennt den Eintritts-Block jetzt config-treu (`heading`, Default `## Aktuelle Welle`) |
| 1.54             | v0.62.0          | 2026-08-22 | Doku-Präzisierung ohne Software-Änderung (Lastenheft 0.62.1): die Wellen-Invariante vergleicht in **beiden** Modi Kennungs-Mengen — zwei flache Wellendokumente derselben Kennung sind ein Element, auch für das Singleton unter `one` (§6, Absatz „Zwei Kardinalitäts-Modelle"). Software-Version bleibt v0.62.0 |
| 1.55             | v0.62.0          | 2026-08-22 | Doku-Korrektur ohne Software-Änderung (§5 Ventil-Überblick): `exempt-paths` gibt es in **fünf** Modulen (`ids` je Muster, `matrix`, `codepaths`, `versions`, `structure` je Regel) — genannt waren drei; und die Abgrenzung des Zeilen-Markers `d-check:ignore` ist vollständig: er wirkt **nur** für `codepaths` und `ids`, alle übrigen Module kennen ihn nicht. Ausdrücklich benannt ist jetzt `diagrams`, das **weder** Datei- **noch** Zeilen-Ventil hat — dort schneidet nur `diagrams.scope` bzw. `scan.ignore`. Software-Version bleibt v0.62.0 |
| 1.56             | v0.63.0          | 2026-08-23 | **Drei Konfigurations-Flächen additiv geweitet** (Lastenheft 0.63.0/0.64.0/0.65.0). `versions.patterns`: eine **Liste** von Muster-Quellen-Paaren statt genau eines Paares — die Kurzform `pin-pattern`/`current-from` **ist** die einelementige Liste, beide Schreibweisen zugleich ⇒ Exit 2, Ausnahmen und Selbst-Ausnahme der Quell-Datei sind **paar-lokal**, und weil eine Befund-Adresse zwei Paare nicht unterscheidet, nennt **eine** Meldung alle Erwartungen mit ihrer Quelle (§5/§6). `structure.headings-match`/`headings-level`: achte Bedingung — **jede** Überschrift der Ebene innerhalb des Abschnitts matcht das Muster, **positiv** formuliert statt als Präfix-Negation, geprüft auf dem Überschriften-**Text** mit der Erkennung des Moduls, gemeldet **je** Überschrift auf **ihrer** Zeile; neuer Grund-Code `section-heading-mismatch`, Default-Ebene Abschnitts-Ebene + 1 (§5/§6). `diagrams.exempt-paths` **und** der Zeilen-Marker `d-check:ignore` — Ventil-Parität zu den übrigen Modulen; der Marker ist hier ein **Token** statt eines HTML-Kommentars und wirkt auf der **Öffnungszeile** für den ganzen Block (§5/§6). **Korrektur der Zeile 1.55:** der Zeilen-Marker wirkt in **vier** Modulen (`codepaths`, `ids`, `versions`, `diagrams`), nicht in zweien — `versions` honoriert ihn seit v0.30.0 —, und `exempt-paths` gibt es in **sechs** Modulen, nicht in fünf. Ohne die neuen Schlüssel ist der Befundsatz byte-identisch |
| 1.57             | v0.63.0          | 2026-08-23 | Doku-Korrektur ohne Software-Änderung, zwei Befunde. **(a) Ventil-Lage bei `citations` ausgesprochen (§5/§6):** `citation-out-of-range` und `citation-inverted-range` entstehen in **zwei** Fähigkeiten, und der Grund-Code sagt nicht, in welcher. Aus `codepaths.check-lines` sind sie über die **feinen** Achsen stummschaltbar (`exempt-paths`, das geteilte `ignore-refs`, der Zeilen-Marker `d-check:ignore`), aus dem Modul `citations` **nicht** — dessen Befund lässt sich nur über die **groben** Achsen loswerden (`scan.ignore`, modul-lokaler `citations.scope`), die für beide gelten. Gegenüberstellung als Tabelle über alle fünf Achsen, samt der ausdrücklichen Einschränkung, dass dieser Unterschied eine **Eigenschaft des heutigen Stands** ist und keine zugesagte Semantik: weder Lastenheft noch Spezifikation begründen ihn. **(b) §5-Abschnitt geschnitten:** *„Zitate und Zeilen-Referenzen gegen ihre Quelle prüfen"* trug **183 Zeilen** und **sechs** Module — `vcs`, `commits`, `planning` (drei Fähigkeiten) und `tracked` standen nicht in seiner Überschrift und waren dort nicht auffindbar. Jetzt fünf Abschnitte, jede Überschrift nennt ihren Inhalt; der Text ist **unverändert bewegt**, nicht umgeschrieben. Der Zensus über alle dreizehn §5-Abschnitte zeigte genau diesen einen Fall — Länge allein ist kein Defekt. Software-Version bleibt v0.63.0 |
| 1.58             | v0.64.0          | 2026-08-27 | **Der Zeilen-Marker `d-check:ignore` verengt sich — bei zwei seiner vier Konsumenten** (Lastenheft 0.68.0/0.69.0). Bei `codepaths` und `ids` zählt er nur noch **außerhalb von Inline-Code** und muss dort in einem **HTML-Kommentar** stehen; eine blanke oder in Backticks geschriebene Erwähnung wirkt nicht mehr. Bei `versions` und `diagrams` bleibt er ein roh gelesenes **Token** — bei `diagrams` strukturell (im Fence ist ein Backtick literaler Inhalt), bei `versions` als **benannte Grenze**. **Wer den Marker setzt, muss ab hier wissen, für welches Modul.** Ebenso verengt: das Modul `citations` erkennt seine Direktive `d-check:cite` nur außerhalb von Fences **und** Inline-Code (Lastenheft 0.67.0) — erst dadurch wurde es überhaupt aktivierbar. Dazu `matrix`: eine Klasse darf ihr `token`-Muster **ausschließlich** tragen, ohne Pfade (Lastenheft 0.66.0). **Ein nicht-kosmetischer Bestands-Effekt ohne Code-Änderung:** das Runtime-Basis-Image trägt einen neuen Bau desselben Tags und damit einen anderen **CA-Trust-Store** — 142 → 150 Wurzelzertifikate; wer `external` gegen HTTPS-Ziele fährt, vertraut ab diesem Image einer anderen Menge. **Korrektur der Zeile 1.57:** deren §5-Satz *„die Ventil-Direktive `d-check:ignore` folgt dieser Regel **nicht**"* galt bei ihrem Schreiben für alle vier Konsumenten und gilt jetzt nur noch für zwei — er stand unqualifiziert 145 Zeilen entfernt von der Stelle, die das Gegenteil sagte, und ist auf `versions`/`diagrams` eingeschränkt |
| 1.59             | v0.65.0          | 2026-08-28 | **Der Docker-Hub-Spiegel ist real, und ein neuer Sensor kam dazu.** §2 nennt den zweiten Bezugsweg jetzt als bestehend statt als angekündigt: `docker.io/pt9912/d-check` trägt dasselbe Bild, und die **Gleichheit ist der Config-Digest** — der **Manifest**-Digest ist registry-lokal, wer per Digest pinnt, nimmt den der Registry, aus der er zieht ([`DC-FA-DIST-002`](../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel), Lastenheft 0.70.0/0.71.0). **Sicherheitsrelevant:** dieses Release behebt **vierzehn** CVEs mit verfügbarem Fix im ausgelieferten Image — neun in `golang.org/x/crypto`, vier in `golang.org/x/net`, eine in `github.com/go-git/go-git/v5`; gehoben auf `x/crypto` v0.53.0, `x/net` v0.56.0, `x/sys` v0.46.0, `go-git` v5.19.2. Wer auf einen älteren Digest gepinnt hat, zieht sie weiter. Gefunden hat sie der neue Nachtlauf-Sensor `make image-scan` (Trivy gegen die **publizierten** Images, digest-gepinnter Scanner, aktuelle Vuln-DB, außerhalb von `gates`) — kein Gate dieses Repos hätte sie je gemeldet, weil alle an einem Commit hängen und CVEs ohne Commit auftauchen. **Keine Software-Änderung am Prüf-Verhalten:** Modulsatz, Grund-Codes und Konfigurations-Fläche sind unverändert; der Befundsatz bleibt byte-identisch |
| 1.60             | v0.66.1          | 2026-08-29 | **Eine Tabellenspalte kurz und gefüllt halten — und die Tabellen-Schlüssel ziehen um** (Lastenheft 0.73.0). Neu ist die **neunte** `structure`-Bedingung (§4.18, §5, §6): jede Zelle einer über ihren **Kopfzeilen-Namen** benannten Spalte liegt in einer Spanne aus **Zeichen** (nicht Bytes), gemeldet auf **ihrer** Zeile statt am Abschnittskopf — `section-cell-oversized` (kürzen), `section-cell-undersized` (ausfüllen) und `section-column-missing` für die nicht adressierbare Spalte. Der Name statt der Position ist Absicht: eine eingefügte Spalte verschöbe eine Positions-Angabe **still**, ein umbenannter Kopf meldet **laut**. **Eine Obergrenze allein lässt die leere Zelle passieren** — null Zeichen liegen unter jeder Schwelle; wer „gefüllt" zusagt, setzt `cell-min-chars`. **Breaking, und es betrifft auch bestehende Konfigurationen:** die tabellenbezogenen Schlüssel stehen jetzt unter der Klammer `table` — `table-order` → `table.order`, `table-column` → `table.order-column`, und die Zellengrenzen leben als **Liste** in `table.column[]` (`name`, `cell-max-chars`, `cell-min-chars`). Mehrere Spalten desselben Abschnitts brauchen damit **keine** zweite Regel mit wiederholtem `files`/`section` mehr. Jeder der fünf alten Schlüssel bricht mit **Exit 2 und dem neuen Ort in der Meldung**; still ignoriert wird keiner. Kein Grund-Code und keine Befund-Form der bestehenden Bedingungen ändert sich |
| 1.61             | v0.67.0          | 2026-08-29 | **Neues Modul `workflows`: Workflow-Referenzen gepinnt und rechte-gedeckt halten** (Lastenheft 0.74.0, §4.19 als eigene Aufgabe, §5-Konfiguration, §6-Modul-Tabelle). Es prüft die `uses:`-Referenzen von CI-Workflows unterhalb eines **konfigurierten** Verzeichnisses (`workflows.dir` — der Ort ist nicht verdrahtet, weil er CI-System-spezifisch ist). Eine **fremde** Referenz nennt einen vollen 40-stelligen Commit-SHA (`uses-pin-missing`) mit Tag-Kommentar dahinter (`uses-pin-untagged`); geprüft wird die **Form**, nicht die Gültigkeit des SHA. Eine **lokale** Referenz (`./…`) braucht keinen Pin — sie löst auf denselben Commit auf wie ihr Aufrufer —, dafür zwei andere Zusagen: das Ziel existiert (`uses-local-missing`), **und** der aufrufende **Job** führt die Rechte, die es verlangt (`uses-local-perms-undeclared`, `uses-local-perms-narrow`). **Das ist kein theoretischer Fall:** ein Job ohne eigenes `permissions:` erbt den Workflow-Kopf und kann nichts weitergeben, was er nicht deklariert — genau daran brach ein Release dieses Projekts **vor dem ersten Job** ab, während das damalige Gate grün meldete. Stufen `none` < `read` < `write`; ein nicht genannter Scope ist `none`, `read-all`/`write-all` setzen jeden. Unlesbares YAML meldet (`workflow-unparsable`), statt übersprungen zu werden. **Hermetisch** (kein git, kein Netz, kein Ausführen), opt-in, fail-closed bei leerer Prüfmenge. **Benannte Grenze:** das Modul liest die **Ziele** lokaler Referenzen, die es nicht scannt — dieselbe Parse-Zusage gilt dort —, und es deckt **eine** Deklarations-Klasse, nicht die Lauffähigkeit eines Workflows |
