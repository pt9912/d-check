# Benutzerhandbuch: d-check

**Handbuch-Version:** 1.28 · **Software-Version:** [v0.43.0](../../version.md#v0.43.0) ·
**Stand:** 2026-07-14 · **Autor:** pt9912

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
nie und macht außer im optionalen Modul `external` keine
Netzwerkzugriffe. Gleiche Eingabe liefert immer dieselbe Ausgabe.

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
docker pull ghcr.io/pt9912/d-check:v0.43.0
```

Das Image läuft als Nicht-root-Prozess; ein **read-only**-Mount des
Repositorys genügt, weil d-check nie schreibt.

### Versionen und Tags

- `:v0.43.0` — eine feste Version (empfohlen für reproduzierbare Läufe; die jeweils
  aktuelle steht in [version.md](../../version.md#aktuell)).
- `:latest` — die jeweils neueste **stabile** Version. Vorabversionen
  (Prereleases, z. B. `v1.0.0-rc1`) erhalten **kein** `:latest`; für
  CI-Pipelines pinnen Sie ohnehin auf eine feste Version oder den Digest
  (Komfort-Einstieg, nicht für reproduzierbare Läufe).

Für vollständig reproduzierbare CI-Läufe pinnen Sie auf den Image-Digest
(in den Release-Notes je Version angegeben):

```bash
docker run --rm -v "$PWD:/repo:ro" \
  ghcr.io/pt9912/d-check@sha256:2963f882c40a0b34d1fc03ba0e91feaf18423e55a35084bead1efa9d5500bd53
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
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.43.0
```

d-check mountet Ihr Repository nach `/repo` und prüft es. Eine typische
Ausgabe bei sauberer Doku:

```text
d-check: 42 Datei(en) geprüft, 0 Befund(e)
```

Bei kaputten Referenzen meldet d-check je einen Befund pro Zeile und endet
mit Exit-Code 1:

```text
docs/anleitung.md:12	fehlt.md	target-missing
d-check: 42 Datei(en) geprüft, 1 Befund(e)
```

Jede Befund-Zeile hat das Format `Datei:Zeile  Ziel  Grund-Code`.

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

## 4. Aufgaben

Jede Aufgabe nennt Ziel, Voraussetzung, Vorgehen und das erwartete
Ergebnis.

### 4.1 Die Dokumentation eines Repos prüfen

**Ziel:** kaputte lokale Links und Heading-Anker finden.
**Voraussetzung:** Docker; ein Repository im aktuellen Verzeichnis.

**Vorgehen:**

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.43.0
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
  ghcr.io/pt9912/d-check:v0.43.0
```

**Ergebnis:** Der Schritt ist grün bei Exit-Code 0 und rot bei 1 oder 2 —
die CI bricht bei kaputten Referenzen ab.

**Hinweise:** `--network none` ist sinnvoll, solange das Modul `external`
nicht aktiv ist (d-check braucht dann kein Netz). Pinnen Sie für
reproduzierbare Läufe auf den Image-Digest (siehe
[Abschnitt 2](#2-installation-und-zugriff)).

### 4.3 Eine Konfiguration anlegen

**Ziel:** eine `.d-check.yml` als Startgerüst erzeugen.
**Voraussetzung:** keine.

**Vorgehen:**

```bash
docker run --rm ghcr.io/pt9912/d-check:v0.43.0 --print-config > .d-check.yml
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
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.43.0 \
  --suggest-config spec/,docs/plan/adr/ > .d-check.yml
```

**Ergebnis:** Ein vorgeschlagenes Gerüst auf stdout mit abgeleiteten
`ids`-Mustern und den opt-in-Modulen, die im Repo echtes Signal liefern.

**Hinweise:** Der Vorschlag ist beratend — prüfen und verengen Sie ihn.
d-check liest dabei nur, es schreibt nichts.

**Harness-Vorlage (`ai-harness` / `ai-harness-init`):** Für ein Repo nach
dem ai-harness-course-Standard liefern zwei reservierte Quellen ein
fertiges Gerüst — ohne die Quellen einzeln aufzulisten. Es enthält die
kanonischen `ids`-Muster, die `matrix`-Referenzrichtung und das **fixe
Standard-Modulset** (`links`, `anchors`, `ids`, `matrix`, `codepaths`, `spans`,
`hostpaths`) samt einem **repo-bewussten `planning`-Block** (aktiv, sobald die
Roadmap existiert). Situative Module ohne ableitbare Konfiguration bleiben aus —
`vcs`/`commits` (die eine Commit-Range brauchen) liefert `--print-mk` als
Makefile-Target, den Rest listet `--print-config`. Welche Quelle, hängt von Ihrer
Ausgangslage ab:

- **Frisches Repo → `ai-harness-init`** (Voll-Kanon, alle Blöcke aktiv).
  Das ist Ihr Zielbild: legen Sie die Struktur an (`spec/`,
  `docs/plan/adr/`, …), dann läuft d-check.

  ```bash
  docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.43.0 \
    --suggest-config ai-harness-init > .d-check.yml
  ```

- **Bestehendes Repo → `ai-harness`** (repo-bewusst): nur Pfade, die Ihr
  Repo schon hat, sind aktiv; fehlende erscheinen auskommentiert mit
  Hinweis (Ihre TODO-Liste). Läuft sofort.

  ```bash
  docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.43.0 \
    --suggest-config ai-harness > .d-check.yml
  ```

Beide sind beratend und kombinierbar mit echten Quellen (z. B. `ai-harness`
zusätzlich zu `spec/lastenheft.md`).

**Kennungs-Präfix (`--id-prefix`):** Das Anforderungs-`ids`-Muster ist
projektspezifisch — nur sein Präfix wechselt pro Repo (d-check: `DC`,
a-check: `AC`, …). Geben Sie es mit `--id-prefix` an:

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.43.0 \
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
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.43.0 \
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
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.43.0 \
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
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.43.0 \
  --enable matrix
```

**Ergebnis:** Verbotene Klassen-Referenzen erscheinen als
`matrix-forbidden`, Referenzen auf Dokumente mit verbotenem Status (etwa
abgelöst) als `matrix-inactive`.

**Richtung innerhalb einer Klasse (`order`/`direction`).** Trägt eine Klasse
zusätzlich `order` (eine Liste von Pfad-Globs, **autoritativste Schicht
zuerst**) und `direction: no-downward`, prüft d-check auch *klasseninterne*
Referenzen: der Rang einer Datei ist der erste passende `order`-Glob; ein
Verweis von einer höher- auf eine niederrangige Datei (auch über mehrere
Stufen) erscheint als `matrix-downward`. Dateien ohne `order`-Treffer sind
rangfrei und werden nicht geprüft. So lässt sich z. B. erzwingen, dass das
Lastenheft nie „abwärts" auf Spezifikation oder Architektur verweist:

```yaml
matrix:
  classes:
    - name: spec
      paths: [spec/lastenheft.md, spec/spezifikation.md, spec/architecture.md]
      order: [spec/lastenheft.md, spec/spezifikation.md, spec/architecture.md]
      direction: no-downward
```

`order` und `direction` gehören zusammen — eines ohne das andere (oder ein
unbekannter `direction`-Wert) ist ein Konfigurationsfehler (Exit 2), damit eine
Richtungsregel nie still wirkungslos ist. Ohne beide Felder verhält sich
`matrix` unverändert.

**Token-Referenzen (`token`) und Provenance-Marker.** `matrix` sieht
standardmäßig nur **Markdown-Links**. Eine Referenz kann aber auch als **bare
ID-Token** im Fließtext stehen (eine Slice-Kennung in einem ADR-Körper). Trägt
eine Klasse ein `token`-Regex, prüft `matrix` auch solche Token: ein Token im
Körper eines Dokuments einer anderen Klasse ist eine Referenz, auf die dieselbe
Regel greift — eine verbotene Kante meldet `matrix-forbidden` (Token in
Markdown-Links und Fenced-Code zählen nicht). Eine **erlaubte Ausnahme**
deklarieren Sie mit dem Marker `<!-- d-check:status-provenance -->` auf derselben
Zeile (etwa „verifiziert in slice-042" als Verifikations-Zeiger, keine
Entscheidungsgrundlage):

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

`exempt-paths` (Globs) nimmt **ganze Dateien** von der `matrix`-Prüfung aus —
nützlich, um unveränderliche Bestandsdokumente zu grandfathern. Ein nicht
kompilierbares `token`-Regex ist ein Konfigurationsfehler (Exit 2); ohne `token`
verhält sich `matrix` unverändert.

### 4.8 Externe Links prüfen (Modul `external`)

**Ziel:** die Erreichbarkeit externer (HTTP-)Links prüfen.
**Voraussetzung:** Netzwerkzugang; `external` ausdrücklich aktiviert.

**Vorgehen** (ohne `--network none`, da Netz gebraucht wird):

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.43.0 \
  --enable external
```

**Ergebnis:** Nicht erreichbare Links erscheinen als `external-status`,
Zeitüberschreitungen als `external-timeout`, zu viele Weiterleitungen als
`external-redirects`.

**Hinweise:** Das ist das **einzige** Modul, das Netzwerkverbindungen
öffnet. Ohne Aktivierung bleibt d-check vollständig netzlos.

### 4.9 Befunde verstehen — Diagnose-Modus (`--doctor`)

**Ziel:** Befunde erklärt und nach Datei gruppiert lesen, mit konkreten
Fix-Vorschlägen.
**Voraussetzung:** keine.

**Vorgehen:**

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.43.0 \
  --enable ids --doctor
```

**Ergebnis:** Statt der knappen Befund-Zeilen erscheint eine Diagnose:

```text
d-check Diagnose — 1 Befund(e) in 1 Datei(en):

docs/anleitung.md
  Z. 12 · Kennung im Fließtext ohne Markdown-Link auf ihre Definition [ids]
      Stelle: ADR-0007
      Fix-Kandidat: ADR-0007 → [`ADR-0007`](plan/adr)
        (Kennung als Markdown-Link auf ihre Definition (docs/plan/adr/) ausführen; Anker ggf. ergänzen)
```

**Hinweise:** `--doctor` ist read-only und gibt nur auf stdout aus. Mit
zusätzlichem `--json` wird dieselbe Diagnose **maschinenlesbar**
ausgegeben (siehe unten). Nicht mit `--repair` kombinierbar (sonst
Exit-Code 2). Einen Fix-Kandidaten gibt es nur dort, wo er **eindeutig**
ableitbar ist.

**Maschinenlesbare Diagnose (`--doctor --json`):** Statt der Prosa ein
JSON-Dokument wie bei `--json` (Abschnitt 4.11),
dessen `findings` je Eintrag zusätzlich `reasonText` (Grund-Klartext) und
`fixCandidate` (`{original, replacement, note}` oder `null`) tragen:

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.43.0 \
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
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.43.0 \
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
| `anchor-missing`, `repo-escape`, `symlink`, `codepath-missing`, `matrix-inactive`, `matrix-forbidden`, `external-status`, `external-timeout`, `external-redirects`, `span-unclosed`, `span-nested-link`, `hostpath-forbidden` | — kein Auto-Fix, von Hand beheben                                                                                                                                               | —                        |

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
  docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.43.0 \
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
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.43.0 --json
```

**Ergebnis:** Ein JSON-Dokument auf stdout mit den Feldern `findings`,
`summary` und `exitCode`:

```json
{
  "findings": [
    { "file": "docs/anleitung.md", "line": 12, "target": "fehlt.md", "rule": "links", "reason": "target-missing" }
  ],
  "summary": { "filesChecked": 42, "findingCount": 1 },
  "exitCode": 1
}
```

**Dasselbe als YAML (`--yaml`):** identische Struktur, nur YAML statt JSON
(`--json` und `--yaml` schließen sich gegenseitig aus):

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.43.0 --yaml
```

<!-- d-check-test:not-config: --yaml-Ausgabe-Beispiel, kein .d-check.yml-Input -->
```yaml
findings:
  - file: docs/anleitung.md
    line: 12
    target: fehlt.md
    rule: links
    reason: target-missing
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
(ADRs), Umsetzungs-Slices und kuratierten Coverage-Quellen referenziert wird.
`WAISE` bedeutet dabei nicht „ohne jede Referenz", sondern: weder ein Slice noch
eine konfigurierte Coverage-Quelle deckt die Anforderung ab; eine ADR allein
verhindert den Waisenstatus nicht.
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

Ab v0.43.0 liest `trace.requirements.format: table` solche Markdown-Pipe-
Tabellen nativ. Die Spalten werden über ihre exakten Header-Namen gebunden:

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
    modality: {}                     # optional: MUSS/SOLLTE/KANN sichtbar machen
```

`table.id-column` und genau eine von `table.text-column` oder
`table.text-columns` sind Pflicht; jeder darin deklarierte alternative Header
muss mindestens einmal in der Quelle vorkommen.
`table.modality-column` ist optional. Die ID-Zelle muss als Ganzes auf
`id-pattern` passen. Header fehlen/doppelt, eine doppelte erkannte ID unter der
Default-Politik `duplicate-ids: error` oder eine fehlerhafte
Tabellenzeilenbreite ergeben Exit 2. Für historische Mehrfachdefinitionen sind
`first` und `last` explizite, deterministische Overrides. Escaped Pipes (`\|`) und Pipes
in einzeiligen Backtick-Code-Spans bleiben Teil derselben Zelle.

**Vorgehen:**

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.43.0 --trace
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
`orphan`; auf Matrix-Ebene `total`/`orphans`). Standardmäßig ist `--trace`
**advisory** (Exit 0, auch bei Waisen). Für ein **Gate** ergänzen Sie
`--require-complete`: Ohne aktive `modality` gatet jede Waise; mit aktiver
`modality` gaten nur Waisen der in `require-levels` gelisteten Stufen. Mindestens
eine so gatende Waise ergibt Exit 1 (die Matrix bleibt unverändert auf stdout).
Ohne `--trace` ist `--require-complete` ein Nutzungsfehler (Exit 2).

> **Versionsgrenze:** Bis v0.42.0 war eine leere RTM kein
> Konfigurationsfehler. Ab v0.43.0 gilt: Ist `requirements.source` mit einem
> nichtleeren Wert explizit konfiguriert oder `format: table` aktiv und werden
> null Anforderungen erkannt, endet `--trace` ebenso wie
> `--trace --require-complete` fail-closed mit Exit 2:

```bash
# DCHECK_IMAGE auf v0.43.0 oder neuer pinnen.
docker run --rm -v "$PWD:/repo:ro" "$DCHECK_IMAGE" --trace --require-complete
```

```text
d-check: error: trace.requirements: Quelle "spec/lastenheft.md" im Format headings ergab 0 Anforderungen
```

`source: ""` gilt weiterhin wie ein abwesender Wert und fällt auf die
Defaultquelle zurück. Ohne explizite Quelle und im Defaultformat bleibt ein
erwartet leerer Bestand kompatibel: leere Matrix, Exit 0. Für eine zusätzliche
Bestandsinvariante kann ein Konsument die erkannte Gesamtzahl weiterhin gegen
seine erwartete Zahl plausibilisieren.

**Andere Repo-Konvention (`trace`-Block).** Standardmäßig erkennt `--trace`
Anforderungen der Gestalt `<PREFIX>-FA-…`/`-QA-…` und Slice-Dateien
`slice-NNN-….md`. Weicht Ihr Repo davon ab, überschreiben Sie die Quell-Achsen
in `.d-check.yml` (jedes Feld optional; ohne `trace`-Block bleibt die Ausgabe
byte-identisch):

```yaml
trace:
  requirements:
    id-pattern: 'GG-[A-Z][A-Z0-9]*-\d{3}'   # Ihre Anforderungs-Kennung
  slices:
    file-pattern: '^(\d+)-.*\.md$'          # z. B. Slices NNN-titel.md
```

Damit lassen sich Quellpfad, Kennungsgestalt, Definitionsformat und ADR-/Slice-
Dateinamen an eine andere Repo-Konvention anpassen (Details und alle Felder in
§5). Ohne `format: table` bleibt die oben beschriebene Heading-Syntax aktiv. Ein
Muster ohne Capture-Gruppe in `file-pattern` ist ein Konfigurationsfehler
(Exit 2).

**Kuratierte Coverage-Quellen (`trace.coverage`).** Liegt die Abdeckung mancher
Anforderungen nicht in ADRs/Slices, sondern in einer **kuratierten Matrix** (z. B.
einer ausgelagerten Traceability-Datei), zählt der Referenz-Scan sie fälschlich
als Waisen. Der opt-in `trace.coverage`-Block liest solche Dateien als **eigene
Coverage-Dimension** (separate Spalte) ein — **range-aware** (`GG-QA-001..006`
deckt alle sechs) und **abschnitts-gescopt**:

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

Eine Anforderung ist dann **Waise** nur ohne Slice **und** ohne Coverage; mit
`--require-complete` bricht das Gate entsprechend. **Wichtig:** ein Sektionsname
ist der **volle Heading-Klartext** (nicht die Kurzform `27.1.1`) — ein Name ohne
Treffer, eine fehlende `files`-Datei oder eine ungültige Range (`AAA>BBB`) sind
Konfigurationsfehler (Exit 2). Ohne `trace.coverage` bleibt die RTM
byte-identisch (keine Coverage-Spalte). Alle Felder in §5.

**Modalität (`trace.requirements.modality`).** Tragen Ihre Anforderungen
RFC-2119-Modalität (MUSS/SOLLTE/KANN) im Text, klassifiziert der opt-in
`modality`-Block jede Anforderung anhand **konfigurierbarer Modal-Verb-
Schlüsselwörter** (mit DE+EN-Defaults) und zeigt die Stufe in einer eigenen
**Modality**-Spalte. `--require-complete` bricht dann **nur** bei Waisen der
`require-levels`-Stufen (Default `[must]`) — SOLLTE/KANN/`unknown` bleiben
advisory:

```yaml
trace:
  requirements:
    id-pattern: 'GG-[A-Z][A-Z0-9]*-\d{3}'
    modality: {}          # {} = Built-in DE+EN-Defaults; require-levels [must]
```

Schon die Schreibweise `modality: {}` **aktiviert** die Klassifikation mit den
Built-in-Defaults. Klassifiziert wird über den ersten Modal-Verb-Treffer im
Body **unterhalb der erkannten Requirement-Überschrift bis zur nächsten gleich-
oder höherrangigen Überschrift** (längster Treffer zuerst: `MUSS NICHT`=may vor
`MUSS`=must, `DARF NICHT`=must; Body wird whitespace-/emphasis-normalisiert).
Die Requirement-Überschrift selbst wird nicht ausgewertet. Im Tabellenformat
liefert eine konfigurierte `table.modality-column` stattdessen ausschließlich
deren Zelleninhalt; ohne diese Spalte wird die jeweils gebundene Textspalte klassifiziert.
Ohne Treffer ⇒ Stufe `unknown` (sichtbar, advisory bei
Default). **Wichtig:** ein `unknown` unter `require-levels: [must]` gatet
**nicht** — ein echtes MUSS mit unaufgeführtem Verb entkäme so dem Gate;
Gegenmittel: die Spalte prüfen und ggf. `levels` ergänzen oder strikt
`require-levels: [must, unknown]`. Fail-closed (Exit 2): leeres Keyword,
gleiches Keyword in zwei Stufen, `unknown` als Stufen-Name, ungültiges
`require-levels`. Ohne `modality` byte-identisch (keine Spalte). Details in §5.

**Tabellenbasiertes Lastenheft migrieren.** Ab v0.43.0 ist die native
Tabellenquelle der direkte Weg: `format: table` setzen und die vorhandenen
ID-/Text-/Modalitätsheader wie oben binden. Verwendet ein Lastenheft mehrere
Text-Header für dieselbe Rolle, werden sie explizit gelistet; historische
Mehrfachdefinitionen brauchen eine bewusste Duplikatpolitik:

```yaml
trace:
  requirements:
    source: spec/lastenheft.md
    id-pattern: '(F|NF|MVP|AK|RAK)-[0-9]+'
    format: table
    table:
      id-column: Kennung
      text-columns: [Anforderung, Akzeptanzkriterium]
      modality-column: Prioritaet
      duplicate-ids: last
```

Zwei Alternativen bleiben für
Tabellen außerhalb der unterstützten Pipe-Tabellen-Grammatik:

1. Die Tabelle in Heading-plus-Body-Abschnitte wie im ersten Beispiel dieses
   Kapitels überführen. ID und Titel stehen in der Überschrift, der normative
   Text samt Modalverb im Body.
2. Aus der Tabelle deterministisch ein separates Heading-plus-Body-Dokument
   erzeugen und dieses als `trace.requirements.source` konfigurieren. Damit die
   Projektion nicht unbemerkt driftet, muss ein eigener Konsistenzsensor die
   erzeugten IDs und Texte gegen die autoritative Tabelle prüfen.

Eine ID-Regex allein migriert das Format nicht; für native Tabellen müssen
`format: table`, `table.id-column`, genau eine von `table.text-column` oder
`table.text-columns` und optional `table.modality-column` konfiguriert sein.

### 4.13 Ein Makefile-Fragment einbinden (`--print-mk`)

**Ziel:** d-check als `doc-check`-Schritt ins eigene Makefile einbinden, ohne
ein Recipe oder Skript zu kopieren — der Image-Pin bleibt bei d-check.
**Voraussetzung:** ein Projekt mit Makefile; eine eigene `.d-check.yml`.

**Vorgehen** (Fragment erzeugen, einbinden):

```bash
docker run --rm ghcr.io/pt9912/d-check:v0.43.0 --print-mk > d-check.mk
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
DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v0.43.0
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

## 5. Konfiguration

Eine optionale `.d-check.yml` in der Repo-Wurzel passt d-check an. Ohne
Datei gelten die Standardwerte (Module `links` und `anchors`, Scan-Wurzeln
`docs/` und `spec/` plus Wurzel-`*.md`). Das vollständige, kommentierte
Schema liefert `--print-config` (siehe [Abschnitt 4.3](#4-aufgaben)).

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
modules: [links, anchors, ids, matrix]
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
  ignore-refs: ["tools/altes-skript.sh"] # Ziel-Pfade nie existenz-prüfen (referenz-weit)
hostpaths:
  prefixes: [home, Users]       # Verzeichnisnamen OHNE / (führender /… ⇒ Exit 2)
external:
  timeout-seconds: 10
diagrams:                      # Kennungen in Diagramm-Fences prüfen
  fences: [mermaid]            # zu öffnende Diagramm-Sprachen (Default mermaid)
  patterns:
    - regex: 'ARC-\d{2}'
      defined-in: spec/architecture.md  # Token muss hier (außerhalb Fences) vorkommen
versions:                      # gepinnte Versions-Verweise gegen die aktuelle Version
  pin-pattern: 'ghcr\.io/[^\s:]+:(v\d+\.\d+\.\d+)'  # Version in Capture-Gruppe 1
  current-from: version.md#aktuell  # Datei#Anker (Span) mit der aktuellen Version
vcs:                           # git-Diff-Immutabilität (opt-in; braucht .git + Range)
  paths: ["docs/plan/adr/[0-9]*.md"]    # geschützte Datei-Klasse (Glob)
  immutable-when: '^\*\*Status:\*\* Accepted'        # BASE gilt ab dieser Zeile als immutabel
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

Das Modul `codepaths` kennt **drei** Ventil-Achsen, um einen Pfad von der
Existenz-Prüfung auszunehmen: den zeilenweisen Marker `d-check:ignore` (eine
einzelne Zeile), `exempt-paths` (ganze **Dateien**, Glob) und `ignore-refs`
(bestimmte **Ziel-Pfade**, Glob — **referenz-weit**, also datei- und
zeilen-unabhängig). `ignore-refs` ist als **Tombstone-Register** gedacht: Wird
eine Datei refaktoriert oder gelöscht, die unveränderliche oder historische Doku
(z. B. akzeptierte ADRs) noch in Inline-Code zitiert, dangelt der eingefrorene
Verweis sonst als `codepath-missing` — nicht behebbar, ohne die immutable Doku zu
editieren. Ein Eintrag in `ignore-refs` erklärt den Pfad als bewusst entfernt; er
unterdrückt **nur** diesen einen Pfad (alle übrigen Verweise der Datei bleiben
geprüft). Bewusster Akt **mit Gate**: vergessen Sie den Eintrag beim Löschen,
meldet `codepath-missing` weiter — nichts verschwindet still. Das Glob matcht den
**aufgelösten, Wurzel-relativen** Ziel-Pfad — nicht die rohe `./`- oder
`../`-Schreibweise, die der Befund unter `Target` zeigt. Tragen Sie den Pfad also so
ein, wie er von der Repo-Wurzel aus lautet.

Das Modul `vcs` vergleicht den **Core** einer immutablen Datei über zwei
git-Stände — es braucht daher eine **Commit-Range** und ist nie ein Default-Modul.
Aufruf über `--range <base>..<head>` (CI/Push) oder `--staged` (lokaler
pre-commit); ohne `.git`/Range bricht es ab (fail-closed). Das `--print-mk`-Target
`doc-immutable` verteilt diese Prüfung an eigene Repos (kein kopiertes Skript):

```bash
# Range aus dem CI; STAGED=1 für einen lokalen pre-commit-Lauf:
make doc-immutable RANGE="$BASE..$HEAD"
```

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

Das Modul `planning` prüft **hermetisch** (nur der Arbeitsbaum, kein git) eine
Planning-Lifecycle-Invariante: die Roadmap trägt den Ruhe-Marker (`marker`, Default
„Keine aktive Welle") in ihrem `heading`-Block (Default `## Aktuelle Welle`) genau
dann, wenn **kein** `slice-*` (`slice-glob`) im Roadmap-Verzeichnis liegt — sonst
`planning-drift`. Bei fehlender oder **mehrdeutiger** (mehrfacher) Überschrift bzw.
fehlender Roadmap-Datei bricht es fail-closed ab. Nur `roadmap` ist Pflicht;
`heading`/`marker`/`slice-glob` sind für abweichende Layouts überschreibbar:

```yaml
planning:
  roadmap: docs/plan/planning/in-progress/roadmap.md
```

Das `--print-mk`-Target `doc-planning` (hermetisch, **ohne** Range) verteilt die
Prüfung an eigene Repos mit demselben Roadmap-Layout:

```bash
make doc-planning
```

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

### Das Versions-Register `version.md` aufbauen (für Modul `versions`)

Das Modul `versions` liest die aktuelle Version aus `version.md#aktuell`
(Schlüssel `current-from`) und prüft damit die gepinnten `ghcr`-Image-Verweise.
Diese Datei ist zugleich ein **auflösbares Link-Ziel** für Erwähnungen der
eigenen Releases. Ein eigenes Repo baut sie nach diesem Muster nach — zwei
Abschnitte und eine Anker-Mechanik:

- **`## Aktuell`** — eine Zeile, die die aktuelle Version nennt. Der
  Heading-Slug `#aktuell` ist der **stabile** Verweispunkt
  (`current-from: version.md#aktuell`); er zeigt immer auf die jeweils aktuelle
  Version, nie auf eine feste Nummer.
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

Eine **nichtleer explizite** `requirements.source` oder `format: table`
aktiviert den Nullmengen-Guard: Fehlt die Quelldatei oder werden darin null
Anforderungen erkannt, endet `--trace` mit Exit 2 vor jeder RTM-Ausgabe.
`source: ""` gilt dagegen wie ein abwesendes Feld und hält den bisherigen
Default-Fallback einschließlich leerer RTM/Exit 0 kompatibel.

Der ADR-/Slice-Referenzscan läuft rekursiv über die gefundenen Dateien unter
dem jeweiligen `dir`. `file-pattern` wird gegen den **Basisdateinamen** und
nicht gegen den vollständigen Pfad geprüft. Dateien ohne Match werden
übersprungen; bei einem Match bildet Capture-Gruppe 1 zusammen mit
`id-prefix` die Owner-Kennung. Danach wird der gesamte Dateiinhalt nach
Treffern von `requirements.id-pattern` durchsucht. Owner je Requirement werden
dedupliziert und deterministisch sortiert.

Zusätzlich liest der opt-in **`trace.coverage`**-Block (eine **Liste** benannter
Quellen) **kuratierte Matrizen** als **eigene Coverage-Dimension** ein — für
Anforderungen, deren Abdeckung weder in ADRs noch Slices liegt (siehe §4.12). Je
Quelle: `files` (explizite Pfade — keine `dir`/`file-pattern`, gegen
ADR-Kontamination), `label` (Owner-Kennung in der eigenen Coverage-Spalte),
`ranges` (Default true; expandiert `<FAM>-AAA..BBB` und `<FAM>-AAA/BBB/CCC`
breiten-erhaltend, gegen `requirements.id-pattern` validiert) sowie `sections`
(Whitelist) / `exclude-sections` (Blacklist) — beide über die
`matrix.exclude-sections`-Span-Semantik, mit dem **vollen Heading-Klartext** als
Namen. Eine Anforderung ist dann **Waise** nur ohne Slice **und** ohne Coverage.
Fail-closed (Exit 2): fehlende `files`-Datei, leeres `label`, ein Sektionsname
ohne Heading-Treffer, eine ungültige Range (`AAA>BBB`/Breite). Ist **keine**
Coverage-Quelle konfiguriert, erscheint **keine** Coverage-Spalte und die RTM
bleibt byte-identisch.

Der opt-in **`trace.requirements.modality`**-Block klassifiziert jede Anforderung
nach RFC-2119-Stufe (siehe §4.12). `levels` (Map Stufe → Modal-Verb-Keywords;
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

## 6. Regelmodule

| Modul       | Standard      | Prüft                                                                                    | Grund-Codes                                                 |
| ----------- | ------------- | ---------------------------------------------------------------------------------------- | ----------------------------------------------------------- |
| `links`     | aktiv         | lokale Links/Bilder: Ziel existiert, innerhalb des Repos                                 | `target-missing`, `repo-escape`, `symlink`                  |
| `anchors`   | aktiv         | Heading-Anker (GitHub-Slugs), inkl. Inline-HTML-Anker                                    | `anchor-missing`                                            |
| `ids`       | opt-in        | Linkpflicht für Kennungen im Fließtext                                                   | `id-unlinked`                                               |
| `matrix`    | opt-in        | erlaubte Referenzrichtung und -status zwischen Dokumentklassen                           | `matrix-forbidden`, `matrix-inactive`                       |
| `codepaths` | opt-in        | explizite Pfade in Inline-Code existieren                                                | `codepath-missing`                                          |
| `spans`     | opt-in        | ungeschlossene Code-Spans, verschachtelte Links                                          | `span-unclosed`, `span-nested-link`                         |
| `hostpaths` | opt-in        | host-lokale absolute Pfade (Maschinen-Layout-Leak)                                       | `hostpath-forbidden`                                        |
| `diagrams`  | opt-in        | Kennungen in Diagramm-Fences (Default `mermaid`) existieren in ihrer `defined-in`-Quelle | `diagram-id-undefined`                                      |
| `versions`  | opt-in        | gepinnte `ghcr`-Image-Verweise tragen die aktuelle Version (aus `version.md#aktuell`), auch in Fences | `version-stale`                                             |
| `pins`      | opt-in        | Content-Pin (`<!-- dpin: … -->`): Ziel-Span eines Links unverändert seit dem Verlinken | `link-stale`                                                |
| `immutable` | opt-in        | Immutabilitäts-Pin (`<!-- immutable: … -->`): normalisierter **Core** einer Datei (ohne Marker-Zeile + `exclude-sections`) unverändert seit dem Pinnen; hermetisch (kein git) | `core-drift`                                                |
| `vcs`       | opt-in (git)  | git-Diff-Immutabilität: **Core** einer immutablen Datei (`immutable-when`) unverändert über eine Commit-Range (`--range`/`--staged`); liest `.git` read-only (kein git-Binary, kein Netz) | `core-drift-vcs`                                            |
| `commits`   | opt-in (git)  | Traceability-Kennung (`id-patterns`) in jeder Commit-Message einer Range (`--range`) bzw. der Pending-Message (`--commit-msg`); liest `.git` read-only (kein git-Binary, kein Netz) | `commit-untraceable`                                        |
| `planning`  | opt-in        | Roadmap-↔-in-progress-Lifecycle-Konsistenz: der Ruhe-Marker (`marker`) steht im `## Aktuelle Welle`-Block genau dann, wenn kein `slice-*` (`slice-glob`) im Verzeichnis liegt; **hermetisch** (kein git), fail-closed bei fehlender/mehrdeutiger Überschrift | `planning-drift`                                            |
| `tracked`   | opt-in (git)  | Getrackt-Status auflösbarer, **existierender** Link-/Bild-Ziele gegen den git-**Index** (gestagt = getrackt, keine `.gitignore`-Interpretation); liest `.git` read-only, **ohne** Range; fail-closed ohne `.git` | `target-untracked`                                          |
| `targets`   | opt-in        | Deklarations-Konsistenz Doku ↔ Build-Targets: jedes in einer Doku-**Tabellenzeile** behauptete `make X` ist eine Makefile-Regel (`makefiles`), und jede Regel steht in der Autoritäts-Doku (`authority`); **hermetisch** (kein git, kein Makefile-Ausführen), fail-closed bei fehlender Datei | `gate-phantom`, `gate-undocumented`                         |
| `external`  | opt-in (Netz) | Erreichbarkeit externer Links                                                            | `external-status`, `external-timeout`, `external-redirects` |

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

### Häufige Fehler (Exit-Code 2)

**„Scan-Wurzel nicht gefunden … wurde das Repository nach /repo
gemountet?"**
Ursache: kein (oder falscher) Mount.
Lösung: `-v "$PWD:/repo:ro"` ergänzen.

**Ungültige Konfiguration (mit Zeilenangabe).**
Ursache: Tippfehler oder unbekannter Schlüssel in `.d-check.yml`.
Lösung: gegen das Gerüst aus `--print-config` abgleichen.

## 8. FAQ

**Verändert d-check meine Dateien?**
Nein. d-check ist ein Lese-Werkzeug; selbst `--repair` gibt nur einen
Patch aus.

**Braucht d-check Netzwerk?**
Nein — außer das Modul `external` ist aktiv. Sonst läuft es vollständig
netzlos (`--network none`).

**Warum meldet d-check nichts in meinem `README`?**
Ohne Konfiguration werden `docs/`, `spec/` und Wurzel-`*.md` geprüft.
Andere Verzeichnisse nehmen Sie über `scan.roots` auf.

**Wie schalte ich eine einzelne Zeile von der `ids`-Prüfung frei?**
Setzen Sie den Marker `d-check:ignore` als HTML-Kommentar an die Zeile.

**Ist die Ausgabe stabil/wiederholbar?**
Ja. Gleiche Eingabe liefert byte-identische Ausgabe (stabile Sortierung).

## 9. Glossar

- **Befund:** eine einzelne Beanstandung (`Datei:Zeile  Ziel  Grund-Code`).
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
- **Sicherheit:** d-check schreibt nie ins Repository und öffnet außer im
  Modul `external` keine Netzwerkverbindungen; übergeben Sie keine
  Geheimnisse als Argumente.

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
| 1.4              | v0.22.0          | 2026-06-22 | `--id-prefix` für `--suggest-config` (§4.4); Traceability-Matrix `--trace` (§4.12); Makefile-Fragment `--print-mk` (§4.13); `:latest` = neueste **stabile** Version (§2)                |
| 1.5              | v0.23.0          | 2026-06-22 | Modul `codepaths`: Datei-Ventil `exempt-paths` ergänzt (§5 Weitere Module) — ganze Dateien von der Inline-Code-Pfad-Prüfung ausnehmen, wie `ids`                                        |
| 1.6              | v0.24.0          | 2026-06-23 | `--print-mk`-Fragment um `doc-trace`/`doc-complete`-Targets + `TRACE_FLAGS` erweitert (§4.13); opt-in `--trace --require-complete` (Vollständigkeits-Gate, Waise ⇒ Exit 1, §4.12)       |
| 1.7              | v0.25.0          | 2026-06-23 | Modul `diagrams` (opt-in): Kennungs-Existenz in Diagramm-Fences (§5 Weitere Module) — `mermaid`-Diagramme auf undefinierte Kennungen prüfen (Befund `diagram-id-undefined`)             |
| 1.8              | v0.26.0          | 2026-06-23 | `--suggest-config ai-harness[-init]`: Kommentar-Hinweis auf die nicht aktivierten situativen opt-in-Module (`external`/`spans`/`hostpaths`/`diagrams`) mit Verweis auf `--print-config` |
| 1.9              | v0.27.0          | 2026-06-23 | `--print-mk`-Fragment (§4.13) um `doc-doctor`/`doc-repair`/`doc-help`-Targets + `DCHECK_DIGEST` (Digest-Override, sticht den Tag) erweitert; alle Targets `##`-annotiert                |
| 1.10             | v0.28.0          | 2026-06-24 | Modul `versions` (opt-in): Versions-Pin-Konsistenz (§5 Weitere Module, §6) — gepinnte `ghcr`-Image-Verweise müssen die aktuelle Version (aus `version.md#aktuell`) tragen, auch in Fences; Befund `version-stale`                |
| 1.11             | v0.29.0          | 2026-06-24 | Modul `pins` (opt-in): Content-Pin gegen inhaltlichen Drift (§6) — ein Link mit `<!-- dpin: sha256:… -->` wird gegen den Hash seines (whitespace-normalisierten, rohen) Ziel-Spans geprüft; Drift → `link-stale`                |
| 1.12             | v0.30.0          | 2026-06-28 | Modul `matrix`: klasseninterne Verweisrichtung (§4.7) — eine Klasse mit `order` (Glob-Rang, autoritativste Schicht zuerst) + `direction: no-downward` meldet klasseninterne Abwärtsverweise als `matrix-downward`                |
| 1.13             | v0.31.0          | 2026-06-28 | Modul `matrix`: token-basierte Referenz-Richtung (§4.7) — eine Klasse mit `token`-Regex fängt verbotene Referenzen auch als bare ID-Token im Körper (`matrix-forbidden`); Provenance-Marker `<!-- d-check:status-provenance -->` nimmt aus, `exempt-paths` grandfathered ganze Dateien                |
| 1.14             | v0.32.0          | 2026-06-28 | Neues opt-in-Modul `immutable` (§6): Immutabilitäts-Pin gegen Core-Drift — eine Datei mit `<!-- immutable: sha256:… -->` wird gegen den whitespace-normalisierten **Core** (Datei ohne Marker-Zeile + `immutable.exclude-sections`) gehasht; Abweichung → `core-drift`. Hermetisch (kein git, read-only-Arbeitsbaum), diagnose-only, opt-in pro Datei                |
| 1.15             | v0.33.0          | 2026-06-29 | Neues opt-in-Modul `vcs` (13., §5/§6): git-Diff-Immutabilität des Core über eine Commit-Range (`--range`/`--staged`, `core-drift-vcs`) — reine-Go-git im read-only `.git` (distroless bleibt, fail-closed, diagnose-only); löst die Skript-Mechanik des `adr-check`-Gates ab. `--print-mk`-Target `doc-immutable` verteilt die git-Garantie an Konsumenten                |
| 1.16             | v0.34.0          | 2026-06-29 | Modul `codepaths`: Referenz-Ventil `ignore-refs` (§5/§6) — bestimmte **aufgelöste Ziel-Pfade** referenz-weit (datei-/zeilen-unabhängig) von der Existenz-/Escape-/Anker-Prüfung ausnehmen, als Tombstone-Register bewusst entfernter Artefakte; löst die Frozen-Doc-Refactoring-Falle. Dritte Ventil-Achse neben `d-check:ignore` (Zeile) und `exempt-paths` (Datei)                |
| 1.17             | v0.35.0          | 2026-07-01 | Neues opt-in-Modul `commits` (14., §5/§6): Traceability-Kennung (`id-patterns`) in jeder Commit-Message über eine Range (`--range`) bzw. der Pending-Message (`--commit-msg`, commit-msg-Hook via stdin); Befund `commit-untraceable`, `exempt-pattern` nimmt Merge/Revert-Betreffe aus — reine-Go-git im read-only `.git` (distroless bleibt, fail-closed, diagnose-only); löst die Skript-Mechanik des `trace-check`-Gates ab. `--print-mk`-Target `doc-commits` verteilt die Range-Prüfung an Konsumenten                |
| 1.18             | v0.36.0          | 2026-07-01 | Neues opt-in-Modul `planning` (15., §5/§6): Roadmap-↔-in-progress-Lifecycle-Konsistenz — der Ruhe-Marker steht im `## Aktuelle Welle`-Block genau dann, wenn kein `slice-*` im Verzeichnis liegt (`planning-drift`); **hermetisch** (nur Roadmap + Verzeichnis-Listing, kein git), fail-closed bei fehlender/mehrdeutiger Überschrift, diagnose-only. Löst die Skript-Mechanik des `planning-check`-Gates ab (letztes Familien-Skript). `--print-mk`-Target `doc-planning` (hermetisch, ohne Range) verteilt die Prüfung                |
| 1.19             | v0.37.0          | 2026-07-03 | Neues opt-in-Modul `tracked` (16., §5/§6): Getrackt-Status auflösbarer, **existierender** Link-/Bild-Ziele gegen den git-**Index** (`target-untracked` — untracked/gitignoriertes Ziel fehlt auf jedem frischen Klon); Index-Wahrheit (gestagt = getrackt), kein Doppelbefund, Ventil `exempt-targets` (aufgelöster Zielpfad, segmentweise validiert), fail-closed ohne `.git`, **ohne** Commit-Range. `--print-mk`-Target `doc-tracked` (zehn Targets, §4.13) verteilt die Prüfung                |
| 1.20             | v0.37.1          | 2026-07-04 | Fix: `--doctor` zeigt für die sieben seit v0.25 hinzugekommenen Grund-Codes (`diagram-id-undefined`, `version-stale`, `link-stale`, `core-drift`, `core-drift-vcs`, `commit-untraceable`, `planning-drift`) den Klartext statt des rohen Codes — auch als `reasonText` in `--doctor --json`/`--yaml`; die Klartext-Liste ist testseitig beidseitig gegen die Grund-Code-Tabelle der Spezifikation verriegelt (fail-closed)                |
| 1.21             | v0.37.1          | 2026-07-04 | Anleitung „Das Versions-Register `version.md` aufbauen" (§5) — Aufbau (`## Aktuell`/`## Verlauf`), die `<a id="vX.Y.Z">`-Anker-Mechanik samt Anker-Wanderung beim Release und ein kopierbares Muster zum Nachbau in eigenen Repos                |
| 1.22             | v0.38.0          | 2026-07-05 | Neues opt-in-Modul `targets` (17., §5/§6): Deklarations-Konsistenz Doku ↔ Build-Targets — jedes in einer Doku-**Tabellenzeile** behauptete `make X` ist eine Makefile-Regel (`gate-phantom`), jede Regel steht in der Autoritäts-Doku (`gate-undocumented`); **hermetisch** (kein git, kein Makefile-Ausführen), fail-closed. Löst den Doku-↔-Makefile-Kern des `gate-consistency.sh`-Meta-Gates ab; `--print-mk`-Target `doc-targets` (elf Targets, §4.13) verteilt die Prüfung                |
| 1.23             | v0.39.0          | 2026-07-06 | `--suggest-config ai-harness[-init]` (§4.4): Vorlage an die **gelebte** Konvention angeglichen — `spans`/`hostpaths` ins fixe Standard-Modulset, repo-bewusster `planning`-Block; `vcs`/`commits` (Commit-Range) via `--print-mk`, `versions`/`targets` bewusst vertagt; kanonische Vorlage der Spezifikation deckt die emittierte Ausgabe 1:1                |
| 1.24             | v0.40.0          | 2026-07-11 | Requirements Traceability Matrix (`--trace`, §4.12) über einen opt-in `trace`-Block quell-/kennungs-konfigurierbar (§5): Anforderungs-Quelldatei + Kennungs-Regex sowie je Referenzklasse Verzeichnis + Dateimuster + Owner-Präfix; Default = Konvention ⇒ byte-identisch, fail-closed bei ungültiger Regex / Muster ohne Capture-Gruppe. Passt diese Quellachsen an abweichende Repo-Konventionen an; die Heading-basierte Definitionssyntax bleibt bestehen                |
| 1.25             | v0.41.0          | 2026-07-11 | Kuratierte Coverage-Quellen `trace.coverage` (§4.12/§5): eine Liste benannter `files` liest Deckungs-Matrizen als **eigene Coverage-Spalte** ein — `ranges` (`GG-QA-001..006` → alle sechs, `/`-Aufzählung), `sections`/`exclude-sections` (voller Heading-Text). Waise = ohne Slice **und** ohne Coverage; fail-closed (fehlende Datei / leeres label / Sektion ohne Treffer / ungültige Range ⇒ Exit 2); ohne Block byte-identisch                |
| 1.26             | v0.42.0          | 2026-07-11 | Modalitäts-Klassifikation `trace.requirements.modality` (§4.12/§5): aus konfigurierbaren Modalverb-Stichwörtern (DE+EN-Defaults, opt-in) klassifiziert die RTM jede Anforderung als MUSS/SOLLTE/KANN in einer **eigenen Modalitäts-Spalte** (längster Treffer, Wortgrenze, Markup-normalisiert); optional gatet `require-levels`, welche Stufen einen Waisen zum Exit-1-Fehler machen (Default: nur MUSS). Fail-closed bei leerem Level/Stichwort, reserviertem `unknown`, Stichwort in zwei Stufen, ungültigem `require-levels` ⇒ Exit 2; ohne Block byte-identisch                |
| 1.27             | v0.42.0          | 2026-07-14 | Trace-Präzisierung (§4.12/§5): Requirement-Definition nur als ATX-Heading mit ID im ersten Token; Tabellen-/Listen-/Fließtextgrenze, Body-only-Modalität, exakte Waisen- und Referenzscan-Semantik, Warnung vor leerer RTM mit Exit 0 sowie Brownfield-Migration für tabellarische Lastenhefte dokumentiert                |
| 1.28             | v0.43.0          | 2026-07-14 | Native Trace-Tabellenquellen (§4.12/§5): `trace.requirements.format: table` bindet ID-, alternative Text- und optionale Modalitätsspalten über exakte Header-Namen; Pipe-Escaping, deterministische Duplikatpolitiken und gemeinsames RTM-Modell. Nichtleer explizite Quelle oder Tabellenmodus endet bei fehlender Quelle, ungültiger Tabellenstruktur oder null Anforderungen fail-closed mit Exit 2; unkonfigurierter Heading-Default bleibt kompatibel                |
