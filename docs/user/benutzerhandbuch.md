# Benutzerhandbuch: d-check

**Handbuch-Version:** 1.0 · **Software-Version:** v0.12.0 ·
**Stand:** 2026-06-18 · **Autor:** pt9912

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
Doku-Hygiene. Es ersetzt eine Familie handgepflegter Einzelskripte durch
ein Werkzeug — ein Image, ein Update-Pfad, repo-spezifisches Verhalten per
Konfiguration.

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
docker pull ghcr.io/pt9912/d-check:v0.12.0
```

Das Image läuft als Nicht-root-Prozess; ein **read-only**-Mount des
Repositorys genügt, weil d-check nie schreibt.

### Versionen und Tags

- `:v0.12.0` — eine feste Version (empfohlen für reproduzierbare Läufe).
- `:latest` — die jeweils neueste Version.

Für vollständig reproduzierbare CI-Läufe pinnen Sie auf den Image-Digest
(in den Release-Notes je Version angegeben):

```bash
docker run --rm -v "$PWD:/repo:ro" \
  ghcr.io/pt9912/d-check@sha256:e65654ef8b35c9329f01eeee693bd0c10f583c9e6e01c89f24dd3c2615de32ac
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
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.12.0
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
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.12.0
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
  ghcr.io/pt9912/d-check:v0.12.0
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
docker run --rm ghcr.io/pt9912/d-check:v0.12.0 --print-config > .d-check.yml
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
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.12.0 \
  --suggest-config spec/,docs/plan/adr/ > .d-check.yml
```

**Ergebnis:** Ein vorgeschlagenes Gerüst auf stdout mit abgeleiteten
`ids`-Mustern und den opt-in-Modulen, die im Repo echtes Signal liefern.

**Hinweise:** Der Vorschlag ist beratend — prüfen und verengen Sie ihn.
d-check liest dabei nur, es schreibt nichts.

### 4.5 Regelmodule zu- und abschalten

**Ziel:** ein Modul aktivieren oder deaktivieren.
**Voraussetzung:** keine.

**Vorgehen** (Optionen sind wiederholbar; Kommandozeile schlägt
Konfiguration):

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.12.0 \
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
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.12.0 \
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
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.12.0 \
  --enable matrix
```

**Ergebnis:** Verbotene Klassen-Referenzen erscheinen als
`matrix-forbidden`, Referenzen auf Dokumente mit verbotenem Status (etwa
abgelöst) als `matrix-inactive`.

### 4.8 Externe Links prüfen (Modul `external`)

**Ziel:** die Erreichbarkeit externer (HTTP-)Links prüfen.
**Voraussetzung:** Netzwerkzugang; `external` ausdrücklich aktiviert.

**Vorgehen** (ohne `--network none`, da Netz gebraucht wird):

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.12.0 \
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
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.12.0 \
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

**Hinweise:** `--doctor` ist read-only und gibt nur auf stdout aus. Es ist
nicht mit `--json` oder `--repair` kombinierbar (sonst Exit-Code 2). Einen
Fix-Kandidaten gibt es nur dort, wo er **eindeutig** ableitbar ist.

### 4.10 Befunde als Patch beheben (`--repair`)

**Ziel:** behebbare Befunde als anwendbaren Patch ausgeben.
**Voraussetzung:** `git` zum Anwenden des Patches.

**Vorgehen** (Patch erzeugen, sichten, anwenden, aufräumen):

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.12.0 \
  --enable ids --repair > fix.patch
# fix.patch sichten (besonders bei --repair-broad), dann anwenden:
git apply fix.patch
rm fix.patch
```

**Ergebnis:** Ein `git apply`-kompatibler unified diff auf stdout. d-check
selbst schreibt nichts — Sie wenden den Patch an.

**Hinweise:**

- **Konservativ** (Standard, `--repair`): nur eindeutige Fixes. In dieser
  Version vor allem nackte Kennungen, die verlinkt werden.
- **Breit** (`--repair-broad`): zusätzlich Best-Guess-Reparaturen (etwa
  ein fehlendes Linkziel auf eine Datei gleichen Namens). Diese sind
  **review-pflichtig**; ihre Markierung erscheint auf **stderr**, damit
  der Patch auf stdout sauber anwendbar bleibt. Lesen Sie den Patch vor
  dem Anwenden.
- Nicht mit `--json` oder `--doctor` kombinierbar.
- **`fix.patch` ist temporär:** nach dem Anwenden löschen
  (`rm fix.patch`) — sonst bleibt die Datei im Repo liegen und kann
  versehentlich mitcommittet werden. Die Pipe-Variante unten kommt ganz
  ohne Zwischendatei aus.
- **Einzeiler (Pipe):** Weil der Patch auf stdout und nur
  Markierung/Zusammenfassung auf stderr gehen, können Sie direkt pipen:

  ```bash
  docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.12.0 \
    --enable ids --repair | git apply
  ```

  Das umgeht aber die Sichtprüfung. Bei `--repair-broad`
  (review-pflichtig) und in CI-Shells mit `set -o pipefail` — dort lässt
  der Befund-Exit-Code 1 die Pipe scheitern, obwohl `git apply` ok war —
  ist die Datei-Variante (oben) die bessere Wahl.

### 4.11 Maschinenlesbare Ausgabe (`--json`)

**Ziel:** Befunde automatisiert weiterverarbeiten.
**Voraussetzung:** keine.

**Vorgehen:**

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.12.0 --json
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

**Hinweise:** Mit `--json` enthält stdout ausschließlich das JSON-Dokument.

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
hostpaths:
  prefixes: ["/home", "/Users"] # host-lokale Pfade melden
external:
  timeout-seconds: 10
```

### Modul-lokaler Scan-Bereich

Ein Modul kann einen eigenen, vom globalen abweichenden Scan-Bereich
bekommen:

```yaml
ids:
  scope:
    roots: [spec, docs/user]
```

## 6. Regelmodule

| Modul | Standard | Prüft | Grund-Codes |
|---|---|---|---|
| `links` | aktiv | lokale Links/Bilder: Ziel existiert, innerhalb des Repos | `target-missing`, `repo-escape`, `symlink` |
| `anchors` | aktiv | Heading-Anker (GitHub-Slugs), inkl. Inline-HTML-Anker | `anchor-missing` |
| `ids` | opt-in | Linkpflicht für Kennungen im Fließtext | `id-unlinked` |
| `matrix` | opt-in | erlaubte Referenzrichtung und -status zwischen Dokumentklassen | `matrix-forbidden`, `matrix-inactive` |
| `codepaths` | opt-in | explizite Pfade in Inline-Code existieren | `codepath-missing` |
| `spans` | opt-in | ungeschlossene Code-Spans, verschachtelte Links | `span-unclosed`, `span-nested-link` |
| `hostpaths` | opt-in | host-lokale absolute Pfade (Maschinen-Layout-Leak) | `hostpath-forbidden` |
| `external` | opt-in (Netz) | Erreichbarkeit externer Links | `external-status`, `external-timeout`, `external-redirects` |

## 7. Fehlerbehebung

### Exit-Codes

| Code | Bedeutung |
|---|---|
| `0` | Prüfung gelaufen, keine Befunde |
| `1` | Prüfung gelaufen, mindestens ein Befund |
| `2` | Nutzungs- oder Umgebungsfehler — die Prüfung lieferte keine verlässliche Aussage |

### Häufige Befunde

**`target-missing` — Linkziel existiert nicht.**
Ursache: Datei umbenannt/verschoben oder Tippfehler im Link.
Lösung: Pfad korrigieren; oder `--doctor` für eine Diagnose und
`--repair-broad` für einen Best-Guess-Patch (review-pflichtig).

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

| Handbuch-Version | Software-Version | Stand | Änderung |
|---|---|---|---|
| 1.0 | v0.12.0 | 2026-06-18 | Erstfassung: alle Use Cases inkl. `--doctor`/`--repair` |
