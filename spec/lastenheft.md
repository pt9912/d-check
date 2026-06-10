# Lastenheft — d-check

**Version:** 0.2.0

**Status:** Draft

**Autor:** pt9912, **Datum:** 2026-06-10.

---

## 1. Zweck und Geltungsbereich

`d-check` ist ein Kommandozeilen-Tool, das Markdown-Dokumentation eines
Repositories auf kaputte Referenzen prüft: lokale Links und
Bildreferenzen, Heading-Anker, Linkpflicht für Anforderungs- und
Entscheidungs-Kennungen sowie Referenzrichtungs-Regeln zwischen
Dokumentklassen (Referenzmatrix). Es konsolidiert zwölf funktional
überlappende Einzeltools (`check_refs.py`, `docs-check.js`,
`verify-doc-refs.sh`), die heute als Kopien in den Repositories unter
`/Development` gepflegt werden.

Das Tool wird als Docker-Image über GHCR verteilt und in CI-Pipelines
sowie lokal als einzelner Prüfschritt aufgerufen. Die Funktionalität
ist in Regelmodule gegliedert, die pro Repository über
Kommandozeilen-Optionen und eine Konfigurationsdatei aktiviert werden —
ein Image, ein Update-Pfad, repo-spezifisches Verhalten per Config
statt per Code-Kopie.

## 2. Stakeholder

| Stakeholder | Rolle | Erwartung |
|---|---|---|
| Repo-Maintainer (pt9912) | Auftraggeber | Ein gepflegtes Tool statt zwölf driftender Kopien; Fixes wirken überall |
| CI-Pipelines der Repos | Konsument | Ein Docker-Step mit deterministischen Exit-Codes und stabiler Ausgabe |
| AI-Agenten (Harness-Sensorik) | Konsument | Verlässlicher, maschinenlesbarer Doku-Sensor (z. B. als `make`-Gate eingebunden) |
| Harness-Projekte nach Kurs-Methodik | Nutznießer | Referenzmatrix- und ID-Regeln der Spec-Stratifizierung maschinell prüfbar |

## 3. Funktionale Anforderungen

> **Schema-Konvention.** Funktionale Anforderungen verwenden von Beginn
> an Bereichskürzel: `DC-FA-<BEREICH>-<NNN>`. Bereiche: `CLI`
> (Aufruf/Ausgabe), `SCAN` (Datei-Auswahl), `LINK` (Link-Modul), `ANCH`
> (Anker-Modul), `ID` (ID-Linkpflicht-Modul), `MTX`
> (Referenzmatrix-Modul), `EXT` (externe Links), `CONF`
> (Konfiguration), `DIST` (Distribution).

### DC-FA-CLI-001 — Aufruf und Scan-Wurzel

**Beschreibung:** Das Tool wird als `d-check [pfad]` aufgerufen. Ohne
Argument ist die Scan-Wurzel das aktuelle Arbeitsverzeichnis; mit
Argument das angegebene Verzeichnis. Die Scan-Wurzel gilt als
Repository-Wurzel für alle Pfadauflösungen.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Repo ohne kaputte Referenzen, when `d-check` in der Repo-Wurzel aufgerufen wird, then Exit-Code 0 und eine Zusammenfassung „N Datei(en) geprüft, 0 Befunde".
- **Boundary:** Given ein Verzeichnis ohne Markdown-Dateien, when `d-check` aufgerufen wird, then Exit-Code 0, Zusammenfassung „0 Datei(en) geprüft", kein Fehler.
- **Negative:** Given ein nicht existierendes Verzeichnis als Argument, when `d-check /gibt/es/nicht` aufgerufen wird, then Exit-Code 2 und eine Fehlermeldung auf stderr.

**Out-of-Scope:** Prüfung mehrerer Repos in einem Aufruf.

---

### DC-FA-CLI-002 — Regelmodul-Auswahl

**Beschreibung:** Die Prüf-Funktionalität ist in benannte Regelmodule
gegliedert: `links`, `anchors`, `ids`, `matrix`, `external`. Ohne
Konfiguration sind `links` und `anchors` aktiv. Module werden über
Kommandozeilen-Optionen (`--enable <modul>`, `--disable <modul>`)
und über die Konfigurationsdatei ([`DC-FA-CONF-001`](#dc-fa-conf-001--konfigurationsdatei))
aktiviert; Kommandozeilen-Optionen haben Vorrang vor der Konfiguration.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Repo mit einem kaputten Anker, aber korrekten Links, when `d-check --disable anchors` aufgerufen wird, then Exit-Code 0 (das Anker-Modul läuft nicht).
- **Boundary:** Given eine Konfigurationsdatei, die `ids` aktiviert, when `d-check --disable ids` aufgerufen wird, then läuft `ids` nicht (Kommandozeile schlägt Konfiguration).
- **Negative:** Given ein unbekannter Modulname, when `d-check --enable foo` aufgerufen wird, then Exit-Code 2 und eine Meldung mit der Liste gültiger Modulnamen.

**Out-of-Scope:** Dynamisches Nachladen externer Module (Plugins).

---

### DC-FA-CLI-003 — Exit-Codes

**Beschreibung:** Das Tool liefert deterministische Exit-Codes:
`0` = Prüfung gelaufen, keine Befunde; `1` = Prüfung gelaufen,
mindestens ein Befund; `2` = Nutzungs- oder Umgebungsfehler (ungültige
Option, ungültige Konfiguration, Scan-Wurzel nicht lesbar) — die
Prüfung hat dann keine verlässliche Aussage geliefert.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Repo ohne kaputte Referenzen, when `d-check` läuft, then Exit-Code 0.
- **Boundary:** Given genau ein Befund, when `d-check` läuft, then Exit-Code 1 (nicht 2), und der Befund wird ausgegeben.
- **Negative:** Given eine ungültige Option, when `d-check --foo` aufgerufen wird, then Exit-Code 2 und keine Prüf-Ausgabe.

**Out-of-Scope:** Differenzierte Exit-Codes pro Befund-Kategorie.

---

### DC-FA-CLI-004 — Ausgabeformate

**Beschreibung:** Befunde werden zeilenweise im Format
`<pfad>:<zeile>\t<ziel>\t<grund>` auf stdout ausgegeben (ein Befund pro
Zeile, Pfade relativ zur Scan-Wurzel); Zusammenfassung und Diagnose
gehen auf stderr. Mit `--json` erfolgt die gesamte Ausgabe auf stdout
als ein maschinenlesbares JSON-Dokument mit mindestens den Feldern
`findings` (Liste mit `file`, `line`, `target`, `rule`, `reason`),
`summary` (`filesChecked`, `findingCount`) und `exitCode`; stdout
enthält dann keine unstrukturierten Textzeilen.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Repo mit zwei kaputten Links, when `d-check` läuft, then genau zwei Befund-Zeilen auf stdout, je mit Datei, Zeilennummer, Ziel und Grund.
- **Boundary:** Given dieselben Befunde, when `d-check --json` läuft, then ist stdout als JSON parsbar, `summary.findingCount` ist 2 und stdout enthält keine Nicht-JSON-Zeilen.
- **Negative:** Given die unbekannte Option `--format` (kein Teil des CLI), when `d-check --json --format xml` aufgerufen wird, then Exit-Code 2 (ungültige Nutzung, vgl. [`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)).

**Out-of-Scope:** Weitere Formate (SARIF, JUnit-XML) in dieser Version.

---

### DC-FA-SCAN-001 — Datei-Auswahl und Ignorier-Regeln

**Beschreibung:** Geprüft werden Markdown-Dateien (`*.md`) unterhalb
konfigurierter Scan-Wurzeln. Default ohne Konfiguration: `docs/`,
`spec/` (jeweils rekursiv, sofern vorhanden) sowie alle `*.md` direkt
in der Repository-Wurzel. Build- und Fremdverzeichnisse (`.git/`,
`node_modules/`, `build/`, `target/`, `.venv/`, `__pycache__/` u. ä.)
werden immer übersprungen. Zusätzliche Ignorier-Muster sind
konfigurierbar. Alle explizit in der Konfiguration deklarierten
Scan-Wurzeln müssen existieren; nur die Default-Wurzeln sind optional.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein kaputter Link in `docs/a.md`, when `d-check` ohne Konfiguration läuft, then wird der Befund gemeldet.
- **Boundary:** Given ein kaputter Link in `node_modules/x/README.md`, when `d-check` läuft, then wird kein Befund gemeldet.
- **Negative:** Given eine Konfiguration, die explizit die Scan-Wurzel `handbuch/` deklariert, die nicht existiert, when `d-check` läuft, then Exit-Code 2 mit Hinweis auf die fehlende Wurzel (Default-Wurzeln dagegen sind optional).

**Out-of-Scope:** Prüfung anderer Markup-Formate (reStructuredText, AsciiDoc, HTML).

---

### DC-FA-LINK-001 — Lokale Link- und Bildreferenzen (Modul `links`)

**Beschreibung:** Markdown-Links `[text](ziel)` und Bildreferenzen
`![alt](ziel)` mit relativem Ziel werden geprüft: das Ziel muss nach
Pfadauflösung existieren und innerhalb der Repository-Wurzel liegen.
Die Pfadauflösung dekodiert Prozent-Escapes vollständig (RFC 3986,
z. B. `%20`); die Dekodierung erfolgt **vor** der Repo-Escape-Prüfung.
Ziele mit externem Schema (`http:`, `https:`, `mailto:`, `ftp:`)
werden von diesem Modul ignoriert; reine Anker-Ziele (`#...`) behandelt
das Modul `anchors`. Fragment-Teile gemischter Ziele (`ziel.md#anker`)
ignoriert dieses Modul ebenfalls — Anker-Validierung ist ausschließlich
Aufgabe des Moduls `anchors` (ist `anchors` deaktiviert, bleiben
fehlerhafte Anker unbemängelt). Vorkommen in Fenced-Code-Blöcken
(``` / ~~~) und in Inline-Code werden nicht geprüft; HTML-Kommentare
werden nicht gesondert behandelt (Links darin gelten als Fließtext).

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Link auf eine existierende Datei im Repo, when das Modul `links` läuft, then kein Befund.
- **Boundary:** Given ein Link auf eine nicht existierende Datei innerhalb eines Fenced-Code-Blocks, when das Modul läuft, then kein Befund.
- **Negative:** Given ein Link `../../etc/passwd`, dessen Ziel die Repository-Wurzel verlässt, when das Modul läuft, then ein Befund mit Grund „verlässt Repository" — auch wenn der Zielpfad existiert.
- **Negative:** Given ein Link auf eine nicht existierende Datei, when das Modul läuft, then ein Befund mit Datei, Zeile, Ziel und Grund, Exit-Code 1.

**Out-of-Scope:** Reference-Style-Links (`[text][ref]` mit separater `[ref]: ziel`-Definition); semantische Prüfung, ob das verlinkte Dokument inhaltlich passt.

---

### DC-FA-LINK-002 — Symlink-Ablehnung

**Beschreibung:** Linkziele, die ein symbolischer Link sind oder deren
Pfad eine Symlink-Komponente enthält, erzeugen einen Befund —
unabhängig davon, wohin der Symlink zeigt (Defense-in-depth gegen
Verweise, die effektiv aus dem Repo herausführen). Die Symlink-Prüfung
hat Vorrang vor der Repo-Escape-Prüfung aus
[`DC-FA-LINK-001`](#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links);
pro Linkziel entsteht genau ein Befund.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Linkziel, das eine reguläre Datei ist, when das Modul `links` läuft, then kein Befund.
- **Boundary:** Given ein Symlink, der auf eine Datei *innerhalb* des Repos zeigt, when das Modul läuft, then trotzdem ein Befund mit Grund „Symlink".
- **Negative:** Given ein Symlink, der auf einen Pfad außerhalb des Repos zeigt, when das Modul läuft, then ein Befund; der Inhalt des Ziels wird nicht gelesen.

**Out-of-Scope:** Eine Konfigurations-Option, die Symlinks erlaubt.

---

### DC-FA-ANCH-001 — Heading-Anker-Validierung (Modul `anchors`)

**Beschreibung:** Links mit Fragment (`ziel.md#anker` sowie `#anker`
innerhalb derselben Datei) werden gegen die tatsächlichen Headings der
Zieldatei geprüft. Die Anker-Bildung folgt dem GitHub-Slug-Verfahren
(Kleinschreibung, Sonderzeichen-Entfernung, Leerzeichen → `-`,
Duplikat-Suffixe `-1`, `-2`, …). Headings aller Ebenen (`#`–`######`)
zählen; bei gleichlautenden Headings erhält das erste Vorkommen in
Dokumentreihenfolge den Basis-Slug, weitere die Suffixe `-1`, `-2`, ….
Existiert die Zieldatei nicht, wird die Anker-Prüfung für diesen Link
übersprungen — der Befund kommt vom Modul `links`
([`DC-FA-LINK-001`](#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)).

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Link `a.md#zweck-und-geltungsbereich` und ein Heading „Zweck und Geltungsbereich" in `a.md`, when das Modul läuft, then kein Befund.
- **Boundary:** Given zwei gleichlautende Headings „Beispiel" in der Zieldatei, when ein Link auf `#beispiel-1` zeigt, then kein Befund (Duplikat-Suffix korrekt aufgelöst).
- **Boundary:** Given ein Link `fehlt.md#x` auf eine nicht existierende Datei, when die Module `links` und `anchors` aktiv sind, then genau ein Befund — aus `links`; das Modul `anchors` schweigt.
- **Negative:** Given ein Link auf `a.md#gibt-es-nicht`, when das Modul läuft, then ein Befund mit Grund „Anker nicht gefunden".

**Out-of-Scope:** Anker in HTML-Tags (`<a name=…>`, `id=…`) als Linkziele.

---

### DC-FA-ID-001 — Linkpflicht für Kennungen (Modul `ids`)

**Beschreibung:** Die Konfiguration deklariert Kennungs-Muster (Regex,
z. B. `ADR-\d{4}`, `LH(-[A-Z0-9]+)+-\d{3}[A-Z]?`) und je Muster das
Ziel-Dokument bzw. Ziel-Verzeichnis, in dem die Kennung definiert ist.
Nackte Vorkommen einer Kennung im Fließtext — außerhalb von Inline-Code
und Fenced-Code-Blöcken — müssen als Markdown-Link auf ihre Definition
ausgeführt sein. Überlappen mehrere Muster für dasselbe Vorkommen,
gilt die Deklarationsreihenfolge in der Konfiguration: das erste
passende Muster gewinnt.

**Akzeptanzkriterien:**

- **Happy Path:** Given das Muster `ADR-\d{4}` und ein Vorkommen `[ADR-0003](docs/plan/adr/0003-x.md)`, when das Modul `ids` läuft, then kein Befund.
- **Boundary:** Given ein Vorkommen `` `ADR-0003` `` in Inline-Code, when das Modul läuft, then kein Befund (Code-Vorkommen sind linkpflichtfrei).
- **Negative:** Given ein nacktes `ADR-0003` im Fließtext, when das Modul läuft, then ein Befund mit Grund „Kennung ohne Link".

**Out-of-Scope:** Automatisches Ermitteln der Muster aus dem Repo-Inhalt; Prüfung, ob die verlinkte Definition inhaltlich zur Kennung passt.

---

### DC-FA-MTX-001 — Referenzmatrix zwischen Dokumentklassen (Modul `matrix`)

**Beschreibung:** Die Konfiguration deklariert Dokumentklassen über
Pfad-Muster (z. B. Contract-Spec `spec/lastenheft.md`, View-Spec
`spec/architecture.md`, ADR `docs/plan/adr/*.md`, Slice
`docs/plan/planning/**/slice-*.md`) und je Klassen-Paar, ob Referenzen
erlaubt sind. Zusätzlich sind Status-Bedingungen deklarierbar:
Referenzen auf Dokumente, deren Status-Feld (`## Status`-Heading oder
`**Status:**`-Zeile) einen verbotenen Wert hat (z. B. `superseded`,
`deprecated`), erzeugen einen Befund; fehlt das Status-Feld, gilt das
Dokument als aktiv. Damit sind die
Referenzrichtungs-Regeln der Spec-Stratifizierung (Lastenheft darf
nicht abwärts auf ADR/Slice verweisen; Slices nur auf aktive ADRs)
maschinell prüfbar.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Slice mit Link auf ein ADR mit Status `Accepted`, when das Modul `matrix` läuft, then kein Befund.
- **Boundary:** Given ein ADR mit Status `Superseded by ADR-0007`, when ein Slice darauf verlinkt, then ein Befund mit Grund „Referenz auf inaktives ADR".
- **Negative:** Given ein Link aus `spec/lastenheft.md` auf eine ADR-Datei, when das Modul läuft, then ein Befund mit Grund „verbotene Abwärtsreferenz" und Angabe beider Dokumentklassen.

**Out-of-Scope:** Semantische Unterscheidung von Verweis-Zwecken (z. B. Verifikations-Zeiger vs. Entscheidungsgrundlage) — das bleibt Review-Aufgabe; Provenance-/Historie-Sektionen können per Konfiguration von der Prüfung ausgenommen werden, eine automatische Erkennung solcher Sektionen ist nicht gefordert.

---

### DC-FA-EXT-001 — Externe Links (Modul `external`, opt-in)

**Beschreibung:** Nur bei explizit aktiviertem Modul werden `http:`-
und `https:`-Links per HTTP-Request (HEAD, Fallback GET) auf
Erreichbarkeit geprüft, mit konfigurierbarem Timeout (Default 10 s,
einstellbar über die Konfigurationsdatei) und begrenzter Parallelität.
Redirects (HTTP 301/302/303/307/308) werden bis zu fünf Stationen
gefolgt; längere Ketten erzeugen einen Befund. Das Modul ist nie Teil
der Default-Module —
Netzwerkzugriffe erfolgen ausschließlich opt-in.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Link auf eine erreichbare URL (HTTP-Status < 400), when das Modul `external` aktiviert ist, then kein Befund.
- **Boundary:** Given das Modul ist *nicht* aktiviert, when `d-check` in einer Umgebung ohne Netzwerk läuft, then erfolgt kein einziger Netzwerkzugriff und externe Links erzeugen keine Befunde.
- **Negative:** Given ein Link, der HTTP 404 liefert oder ins Timeout läuft, when das Modul aktiviert ist, then ein Befund mit Grund (Statuscode bzw. „Timeout").

**Out-of-Scope:** Link-Rotting-Archivierung, Redirect-Ketten-Bewertung, Authentifizierung gegen geschützte URLs.

---

### DC-FA-CONF-001 — Konfigurationsdatei

**Beschreibung:** Eine optionale Datei `.d-check.yml` in der
Repository-Wurzel deklariert: Scan-Wurzeln, Ignorier-Muster, aktive
Module, Kennungs-Muster (für `ids`), Dokumentklassen und
Referenzregeln (für `matrix`) sowie Modul-Parameter (z. B. Timeout für
`external`). Ohne Konfigurationsdatei gelten die in den jeweiligen
Anforderungen definierten Defaults; eine vorhandene Datei wird beim
Start vollständig validiert — Syntax **und** Inhalt (z. B. ungültige
Regexe, unbekannte Modulnamen). Jeder Konfigurationsfehler führt zu
Exit-Code 2; es findet keine Prüfung mit stillschweigenden Defaults
statt.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine `.d-check.yml` mit einem Ignorier-Muster `docs/archive/**`, when `d-check` läuft, then erzeugen kaputte Links unter `docs/archive/` keine Befunde.
- **Boundary:** Given keine `.d-check.yml`, when `d-check` läuft, then laufen genau die Default-Module (`links`, `anchors`) auf den Default-Scan-Wurzeln.
- **Negative:** Given eine syntaktisch ungültige `.d-check.yml`, when `d-check` läuft, then Exit-Code 2 mit Fehlermeldung inkl. Zeilenangabe; es wird keine Prüfung mit stillschweigenden Defaults durchgeführt.

**Out-of-Scope:** Konfigurations-Vererbung über mehrere Verzeichnisebenen; Migrations-Assistent von Alt-Tool-Konfigurationen.

---

### DC-FA-DIST-001 — Docker-Image

**Beschreibung:** Das Tool wird als Container-Image
`ghcr.io/pt9912/d-check` mit Semver-Tags veröffentlicht. Aufruf:
`docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:<tag>`
prüft das nach `/repo` gemountete Repository. Das Image läuft als
Non-root-Prozess; ein read-only-Mount genügt, da das Tool das geprüfte
Repository nie beschreibt. Default-Befehl des Images ist die Prüfung
von `/repo`; CLI-Optionen werden als Container-Argumente angehängt.
Ergebnis und Exit-Code sind identisch zur nativen Ausführung.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Repo mit einem kaputten Link, when der Container mit gemountetem Repo läuft, then identische Befund-Ausgabe und Exit-Code 1 wie beim nativen Aufruf.
- **Boundary:** Given ein read-only gemountetes Repo (`:ro`), when der Container läuft, then keine Schreibfehler und vollständige Prüfung.
- **Negative:** Given kein Mount auf `/repo` (Verzeichnis leer oder fehlt), when der Container läuft, then Exit-Code 2 mit Hinweis auf den erwarteten Mount.

**Out-of-Scope:** Distributionswege jenseits GHCR (Homebrew, Paketmanager); deren Bewertung erfolgt später per ADR.

## 4. Nichtfunktionale Anforderungen

### DC-QA-01 — Performance

- **Anforderung:** Ein vollständiger Lauf der Default-Module über ein Repository mit 1.000 Markdown-Dateien (insgesamt ≤ 20 MB) terminiert in < 5 s auf Standard-Entwickler-Hardware (2 vCPU); das Modul `external` ist von dieser Schranke ausgenommen.
- **Messmethode:** Benchmark-Lauf gegen ein generiertes Fixture-Repo, Definition in `spec/spezifikation.md` (folgt).

### DC-QA-02 — Determinismus

- **Anforderung:** Identische Eingabe (Repo-Stand + Konfiguration + Optionen) liefert byte-identische Ausgabe; Befunde sind stabil sortiert (Pfad, dann Zeilennummer).
- **Messmethode:** Test führt denselben Lauf 10× aus und vergleicht Hashes der Ausgabe.

### DC-QA-03 — Seiteneffektfreiheit und Netzwerk-Sparsamkeit

- **Anforderung:** Das Tool schreibt nie in das geprüfte Repository und öffnet außer im explizit aktivierten Modul `external` keine Netzwerkverbindungen.
- **Messmethode:** Integrationstest mit read-only-Mount und netzwerkloser Umgebung (`docker run --network none`), alle Module außer `external` aktiv.

### DC-QA-04 — Migrationsabdeckung der Alt-Tools

- **Anforderung:** Jedes der zwölf Quell-Tools ist durch `d-check` mit passender Konfiguration ersetzbar: Auf dem jeweiligen Repo-Stand meldet `d-check` mindestens dieselben echten Befunde wie das Alt-Tool und erzeugt keine False-Positives, die eine bislang grüne CI brechen.
- **Messmethode:** Pilot-Migration in mindestens drei Repos — je ein Vertreter der Shell-Familie (`verify-doc-refs.sh`), der Python-Familie (`check_refs.py`, inkl. u-boot-Vollausbau) und der JS-Familie (`docs-check.js`) — mit Vergleichslauf Alt-Tool vs. `d-check`.

## 5. Globale Out-of-Scope-Punkte

- Architektur-/Boundary-Checks (C++-Include-Regeln, Go-Package-Graph, Rust-Crate-Grenzen) — bleiben repo-lokale, sprachspezifische Tools.
- Build-Reproduzierbarkeits-Prüfung (Funktionalität von `repro-check.sh`).
- Mathematik-/Formel-Validierung (MathJax-Rendering-Checks aus `euler-fourier-hilbert`) — Kandidat für eine spätere Version, nicht Teil von 0.x.
- Automatisches Reparieren kaputter Referenzen (Auto-Fix) — das Tool berichtet nur.
- Prüfung von Nicht-Markdown-Formaten (reStructuredText, AsciiDoc, HTML, Jupyter-Notebooks).
- Rechtschreib-, Stil- oder Markdown-Lint-Prüfungen (dafür existieren etablierte Tools).

## 6. Glossar

| Begriff | Bedeutung im Lastenheft |
|---|---|
| Befund | Eine einzelne festgestellte Regelverletzung mit Datei, Zeile, Ziel und Grund. |
| Regelmodul | Benannte, einzeln aktivierbare Prüf-Einheit (`links`, `anchors`, `ids`, `matrix`, `external`). |
| Scan-Wurzel | Verzeichnis, unterhalb dessen Markdown-Dateien gesucht werden; zugleich Bezugspunkt der Pfadauflösung. |
| Anker | Fragment-Teil eines Links (`#…`), das auf ein Heading der Zieldatei zeigt (GitHub-Slug-Verfahren). |
| Repo-Escape | Linkziel, dessen aufgelöster Pfad außerhalb der Repository-Wurzel liegt. |
| Kennung | Textuelle ID nach deklariertem Muster (z. B. `ADR-0003`), für die Linkpflicht gelten kann. |
| Dokumentklasse | Über Pfad-Muster definierte Gruppe von Dokumenten (z. B. Contract-Spec, ADR, Slice) als Knoten der Referenzmatrix. |
| Referenzmatrix | Deklaration, welche Dokumentklasse auf welche verweisen darf, inkl. Status-Bedingungen. |
| Aktives ADR | ADR, dessen Status-Feld keinen verbotenen Wert (`superseded`, `deprecated`) trägt. |
| Quell-Tools | Die zwölf konsolidierten Alt-Tool-Vorkommen aus drei Familien (Shell: `verify-doc-refs.sh`, Python: `check_refs.py`, JavaScript: `docs-check.js`) in den Repos unter `/Development`. |

## 7. Historie

| Version | Datum | Änderung | Verweis |
|---|---|---|---|
| 0.1.0 | 2026-06-10 | Initiale Fassung (Konsolidierung von 12 Quell-Tools, Modul-Schnitt, Docker-Distribution) | — |
| 0.2.0 | 2026-06-10 | Review-Runde R1: Modul-Schnitt `links`/`anchors` präzisiert (Fragment-Zuständigkeit, fehlende Zieldatei), Slug-Duplikat-Reihenfolge, Symlink-Vorrang, RFC-3986-Dekodierung vor Escape-Prüfung, Redirect-Regel `external`, Muster-Präzedenz `ids`, Status-Default `matrix`, Scan-Wurzel- und Config-Vollvalidierung, Out-of-Scope Reference-Style-Links, Image-Default-Befehl | — |
