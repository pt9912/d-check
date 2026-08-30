# Lastenheft — d-check

**Version:** 0.77.0

**Status:** Draft

> **Bewusst strenger als verlangt.** Vor `Accepted` wäre dieses Dokument frei
> änderbar — ohne Change Request, ohne Versions-Bump, ohne Historie-Zeile.
> Dieses Repo führt Bump und Historie trotzdem seit der ersten Fassung: die
> Historie ist der einzige Ort, an dem eine Anforderungs-Änderung ihren Beleg
> hinterlässt, und sie erst ab `Accepted` zu führen hieße, den Bestand bis
> dahin spurlos wachsen zu lassen. Die Spalte **Verweis** bleibt deshalb
> durchgehend `—`: solange die CR-Pflicht nicht begonnen hat, gibt es keinen
> externen Vorgang, den sie nennen könnte. Als Adaption geführt in
> [`MR-032`](../harness/conventions.md#mr-032) — dort steht auch ihr
> Auflösungs-Trigger.

**Autor:** pt9912, **Datum:** 2026-06-10.

---

## 1. Zweck und Geltungsbereich

`d-check` ist ein Kommandozeilen-Tool, das Markdown-Dokumentation eines
Repositories auf kaputte Referenzen prüft: lokale Links und
Bildreferenzen, Heading-Anker, Linkpflicht für Anforderungs- und
Entscheidungs-Kennungen sowie Referenzrichtungs-Regeln zwischen
Dokumentklassen (Referenzmatrix). **Und ob ein Dokument selbst richtig
gebaut ist:** ob es die Abschnitte trägt, die seine Klasse verlangt, und ob
diese Abschnitte die zugesagten Bausteine enthalten. Beide Fragen sind
Gegenstand des Tools. Die zweite ist nicht neu — `spans` und `hostpaths`
prüfen längst den Text selbst, nicht seine Verweise —, aber sie war nie
ausgesprochen: die Module wuchsen entlang der Referenz-Frage, und was
daneben entstand, sah wie Einzelfall aus. Mit dem Modul `structure` wird sie
als eigene Frage benannt, statt weiter mitzulaufen. Es konsolidiert dreizehn funktional
überlappende Einzeltools (`check_refs.py`, `docs-check.js`,
`verify-doc-refs.sh`), die heute als Kopien in den
Schwester-Repositories des Entwicklungs-Workspace gepflegt werden.

Das Tool wird als Docker-Image über GHCR verteilt und in CI-Pipelines
sowie lokal als einzelner Prüfschritt aufgerufen. Die Funktionalität
ist in Regelmodule gegliedert, die pro Repository über
Kommandozeilen-Optionen und eine Konfigurationsdatei aktiviert werden —
ein Image, ein Update-Pfad, repo-spezifisches Verhalten per Config
statt per Code-Kopie.

## 2. Stakeholder

| Stakeholder | Rolle | Erwartung |
|---|---|---|
| Repo-Maintainer (pt9912) | Auftraggeber | Ein gepflegtes Tool statt dreizehn driftender Kopien; Fixes wirken überall |
| CI-Pipelines der Repos | Konsument | Ein Docker-Step mit deterministischen Exit-Codes und stabiler Ausgabe |
| AI-Agenten (Harness-Sensorik) | Konsument | Verlässlicher, maschinenlesbarer Doku-Sensor (z. B. als `make`-Gate eingebunden) |
| Harness-Projekte nach Kurs-Methodik | Nutznießer | Referenzmatrix- und ID-Regeln der Spec-Stratifizierung maschinell prüfbar |

## 3. Funktionale Anforderungen

> **Schema-Konvention.** Funktionale Anforderungen verwenden von Beginn
> an Bereichskürzel: `DC-FA-<BEREICH>-<NNN>`. Bereiche: `CLI`
> (Aufruf/Ausgabe), `SCAN` (Datei-Auswahl), `LINK` (Link-Modul), `ANCH`
> (Anker-Modul), `ID` (ID-Linkpflicht-Modul), `MTX`
> (Referenzmatrix-Modul), `EXT` (externe Links), `CODE`
> (Inline-Code-Pfade), `SPAN` (Markdown-Span-Artefakte), `HOST`
> (Host-Pfad-Hygiene), `DIAG` (Diagramm-Kennungen), `VER`
> (Versions-Pin-Konsistenz), `PIN` (Content-Pin/Drift), `IMM`
> (Immutabilitäts-/Core-Pin), `VCS` (git-historienbasierte
> Core-Immutabilität), `COMMITS` (Traceability-Kennung in
> Commit-Messages), `PLAN` (Planning-Lifecycle-Konsistenz), `TRK`
> (Getrackt-Status von Referenz-Zielen), `TGT` (Deklarations-Konsistenz
> Doku ↔ Build-Targets), `COV` (kuratierte Coverage-Quellen der RTM),
> `MOD` (Modalitäts-Klassifikation der Anforderungen), `REQ`
> (Anforderungsquellen der RTM), `XREF` (Kreuzverweis-Konsistenz zweier
> Traceability-Sichten), `REF` (geteiltes Referenz-Ventil `ignore-refs`,
> Ziel-Achse zu `SCAN`), `CITE` (Verbatim-Zitat-Verifikation, Modul
> `citations`), `SRC` (Upstream-Content-Drift externer Quellen, Modul
> `sources`), `STRUCT` (Struktur-Invarianten
> innerhalb eines Dokuments, Modul `structure`), `WF`
> (Workflow-Deklarations-Konsistenz, Modul `workflows`), `CONF`
> (Konfiguration), `DIST` (Distribution).

### DC-FA-CLI-001 — Aufruf und Scan-Wurzel

**Beschreibung:** Das Tool wird als `d-check [pfad]` aufgerufen. Ohne
Argument ist die Scan-Wurzel das aktuelle Arbeitsverzeichnis; mit
Argument das angegebene Verzeichnis. Die Scan-Wurzel gilt als
Repository-Wurzel für alle Pfadauflösungen. Die Hilfe (`-h`/`--help`)
nennt eine Synopsis (`d-check [optionen] [pfad]`), beschreibt das
Pfad-Argument als Scan-Wurzel (Default: aktuelles Verzeichnis) und
verweist für das Konfigurations-Format auf
[`DC-FA-CLI-005`](#dc-fa-cli-005--konfigurations-gerüst-ausgeben)
(`--print-config`) — sie dupliziert das Format nicht.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Repo ohne kaputte Referenzen, when `d-check` in der Repo-Wurzel aufgerufen wird, then Exit-Code 0 und eine Zusammenfassung „N Datei(en) geprüft, 0 Befunde".
- **Boundary:** Given ein Verzeichnis ohne Markdown-Dateien, when `d-check` aufgerufen wird, then Exit-Code 0, Zusammenfassung „0 Datei(en) geprüft", kein Fehler.
- **Negative:** Given ein nicht existierendes Verzeichnis als Argument, when `d-check /gibt/es/nicht` aufgerufen wird, then Exit-Code 2 und eine Fehlermeldung auf stderr.
- **Hilfe:** Given `d-check --help`, when aufgerufen, then Exit-Code 0 und die Nutzung auf stderr enthält die Synopsis mit `[pfad]`, die Pfad-Argument-Beschreibung, einen Verweis auf `--print-config` **und die URL des Benutzerhandbuchs** auf dem **Hauptzweig** (ohne Versionsangabe).

**Out-of-Scope:** Prüfung mehrerer Repos in einem Aufruf.

---

### DC-FA-CLI-002 — Regelmodul-Auswahl

**Beschreibung:** Die Prüf-Funktionalität ist in benannte Regelmodule
gegliedert: `links`, `anchors`, `ids`, `matrix`, `external`,
`codepaths`, `spans`, `hostpaths`, `diagrams`, `versions`, `pins`,
`immutable`, `vcs`, `commits`, `planning`, `tracked`, `targets`, `citations`, `sources`, `structure`. Ohne Konfiguration sind `links` und `anchors` aktiv. Module werden über
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
Prüfung hat dann keine verlässliche Aussage geliefert. Die Ausgabe-Modi
[`DC-FA-CLI-007`](#dc-fa-cli-007--diagnose-modus) (`--doctor`) und
[`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch) (`--repair`) folgen
denselben Codes; ihre Diagnose- bzw. Patch-Ausgabe erscheint auf stdout
unabhängig vom Code.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Repo ohne kaputte Referenzen, when `d-check` läuft, then Exit-Code 0.
- **Boundary:** Given genau ein Befund, when `d-check` läuft, then Exit-Code 1 (nicht 2), und der Befund wird ausgegeben.
- **Negative:** Given eine ungültige Option, when `d-check --foo` aufgerufen wird, then Exit-Code 2 und keine Prüf-Ausgabe.

**Out-of-Scope:** Differenzierte Exit-Codes pro Befund-Kategorie.

---

### DC-FA-CLI-004 — Ausgabeformate

**Beschreibung:** Befunde werden zeilenweise im Format
`<pfad>:<zeile>\t<ziel>\t<grund>` auf stdout ausgegeben (ein Befund pro
Zeile, Pfade relativ zur Scan-Wurzel); trägt der Befund eine **Erläuterung**
(das Befund-Feld `message`), folgt sie als
**vierte** tab-getrennte Spalte — ein Befund bleibt damit **eine** Zeile, und
ein Befund ohne Erläuterung ist unverändert. **Der Text der Erläuterung ist
nicht stabilitätsgarantiert** — seine Anwesenheit als viertes Feld ist es, sein
Wortlaut nicht. Zusammenfassung und Diagnose
gehen auf stderr. Mit `--json` **oder** `--yaml` erfolgt die gesamte
Ausgabe auf stdout als ein maschinenlesbares Dokument (JSON bzw. YAML) mit
mindestens den Feldern `findings` (Liste mit `file`, `line`, `target`,
`rule`, `reason` sowie `message`, wo der Befund eine Erläuterung trägt),
`summary` (`filesChecked`, `findingCount`) und
`exitCode` — **gleiche Struktur, nur Serialisierung verschieden**; stdout
enthält dann keine unstrukturierten Textzeilen. `--json` und `--yaml`
schließen sich gegenseitig aus (Nutzungsfehler, Exit-Code 2). Die
Ausgabe-Modi [`DC-FA-CLI-007`](#dc-fa-cli-007--diagnose-modus)
(`--doctor`) und [`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch)
(`--repair`) ersetzen das Default-stdout-Format durch eine Diagnose bzw.
einen unified diff; `--doctor` ist mit `--json` **oder** `--yaml`
**kombinierbar** — dann gibt es die Diagnose maschinenlesbar (JSON bzw.
YAML statt Prosa, siehe
[`DC-FA-CLI-007`](#dc-fa-cli-007--diagnose-modus)). `--repair` mit
`--json`/`--yaml` und `--doctor` mit `--repair` bleiben **nicht
kombinierbar** (Nutzungsfehler, Exit-Code 2 nach
[`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)); eine JSON-/YAML-Variante
des Patches ist out of scope.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Repo mit zwei kaputten Links, when `d-check` läuft, then genau zwei Befund-Zeilen auf stdout, je mit Datei, Zeilennummer, Ziel und Grund.
- **Boundary:** Given dieselben Befunde, when `d-check --json` läuft, then ist stdout als JSON parsbar, `summary.findingCount` ist 2 und stdout enthält keine Nicht-JSON-Zeilen.
- **Negative:** Given die unbekannte Option `--format` (kein Teil des CLI), when `d-check --json --format xml` aufgerufen wird, then Exit-Code 2 (ungültige Nutzung, vgl. [`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)).
- **YAML:** Given dieselben zwei Befunde, when `d-check --yaml` läuft, then ist stdout als YAML parsbar mit derselben Struktur wie `--json` (`summary.findingCount` ist 2, keine unstrukturierten Zeilen); `d-check --doctor --yaml` ergänzt je `findings`-Eintrag `reasonText` und `fixCandidate` (analog `--doctor --json`).
- **Vierte Spalte:** Given ein Befund mit Erläuterung (etwa `target-missing` des Moduls `links`), when `d-check` läuft, then trägt seine Zeile **vier** tab-getrennte Felder, ist **eine** Zeile, und die Felder 1–3 sind unverändert; ein Befund ohne Erläuterung bleibt dreispaltig.
- **YAML Negative:** Given die Kombination `d-check --json --yaml`, when aufgerufen, then Exit-Code 2 (sich ausschließende Formate); ebenso `d-check --repair --yaml`.

**Out-of-Scope:** Weitere Formate (SARIF, JUnit-XML) in dieser Version (YAML ist mit 0.19.0 ergänzt).

---

### DC-FA-CLI-005 — Konfigurations-Gerüst ausgeben

**Beschreibung:** Mit der Option `--print-config` gibt `d-check` ein
kommentiertes `.d-check.yml`-Startgerüst auf **stdout** aus und endet
mit Exit-Code 0 — **ohne das geprüfte Repository zu lesen oder zu
beschreiben** (kein Scan). Der Aufrufer leitet selbst um
(`d-check --print-config > .d-check.yml`); das Werkzeug schreibt
niemals selbst (read-only-Kernvertrag
[`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
Das Gerüst ist **statisch** (nicht aus Repo-Inhalt abgeleitet),
deterministisch ([`DC-QA-02`](#dc-qa-02--determinismus)) und ist
gültiges, vom eigenen Konfigurations-Parser vollständig validierendes
YAML ([`DC-FA-CONF-001`](#dc-fa-conf-001--konfigurationsdatei)); es
dokumentiert die verfügbaren Module und Optionen als Kommentare, damit
sie sichtbar sind.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein beliebiger Aufruf-Kontext, when `d-check --print-config` läuft, then liegt auf stdout ein kommentiertes YAML, Exit-Code 0, und der eigene Konfigurations-Parser akzeptiert es ohne Fehler.
- **Boundary:** Given ein Repo, das bereits eine `.d-check.yml` enthält, when `d-check --print-config` läuft, then wird diese weder gelesen noch verändert und es findet kein Scan statt (gleiche Ausgabe wie ohne vorhandene Datei).
- **Negative:** Given ein read-only gemountetes Repository, when `d-check --print-config` läuft, then entsteht kein Schreibzugriff und Exit-Code 0 (Seiteneffektfreiheit wie [`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

**Out-of-Scope:** Schreiben der Datei durch das Werkzeug selbst (immer stdout); Ableiten der Konfiguration aus dem Repo-Inhalt (eigener, späterer Modus — die Muster-Ableitung bleibt für die *Prüfung* Out-of-Scope, [`DC-FA-ID-001`](#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)).

---

### DC-FA-CLI-006 — Konfigurations-Vorschlag aus Autoritäts-Dokumenten

**Beschreibung:** Mit `--suggest-config <quelle>[,<quelle>…]` macht
`d-check` einen **Lese**-Durchgang über das Repository und gibt ein
*vorgeschlagenes* `.d-check.yml`-Gerüst auf **stdout** aus — wie
[`DC-FA-CLI-005`](#dc-fa-cli-005--konfigurations-gerüst-ausgeben)
schreibt es niemals selbst (read-only-Kernvertrag
[`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
bleibt; der Aufrufer leitet um). Jede `<quelle>` ist eine Datei oder ein
Verzeichnis, in dem Kennungen **definiert** sind (z. B. das Lastenheft,
die Konventionen, `adr/`). Aus den dort definierten Kennungen
(führendes Token von ATX-Headings, das der allgemeinen Kennungs-Gestalt
entspricht) wird je Quelle ein `ids`-Muster abgeleitet: ein `regex`,
der **mindestens alle** dort gefundenen Kennungen matcht
(Round-Trip-Garantie), mit `target` = der Quelle, und die Quell-Kennungen
als Kommentar (Nachvollziehbarkeit). Zusätzlich werden die
opt-in-Module vorgeschlagen, die in einem Probelauf echtes Signal
liefern, sowie ein Scan-Scope. Der Vorschlag ist **advisory** — der
Mensch prüft und verengt. Die Ableitung ist deterministisch
([`DC-QA-02`](#dc-qa-02--determinismus)) und das Gerüst dekodiert über
den eigenen Parser ([`DC-FA-CONF-001`](#dc-fa-conf-001--konfigurationsdatei)).

**Reservierte Quellen `ai-harness` / `ai-harness-init`:** Statt eines
Pfads schlagen sie ein an die adoptierte **ai-harness-course-Konvention**
angelehntes Gerüst vor — die strukturellen Konventionen, die eine reine
Quellen-Ableitung nicht erzeugen kann: die kanonischen `ids`-Muster (für
`ADR-`-, `DC-`-, `MR-`-, `slice`- und Carveout-Kennungen mit ihren
Definitions-`target`s), die `matrix`-Klassen
(Vertrag/Spezifikation/ADR/Planning) samt **Referenzrichtung** (kein
Spec-Stratum verweist abwärts auf ADRs oder Planning-Artefakte), das
**Standard-Modulset** (`links`, `anchors`, `ids`, `matrix`, `codepaths`,
`spans`, `hostpaths`) samt einem **repo-bewussten `planning`-Block** sowie einen
Scan-Scope (`spec`, `docs`, `harness`; modul-lokaler `ids`-Scope auf
`spec`+`docs/user`, `exempt-paths` für Review-Reports).

**Aufnahme ins Modulset (geschlossene Menge).** Ein Modul gehört in die Vorlage
genau dann, wenn es (K1) die adoptierte Konvention tatsächlich nutzt, (K2)
ableitungsfrei oder mit konventions-festen Werten läuft, (K3) über den Baum-Scan
statt über eine Laufzeit-Range arbeitet und (K4) hermetisch/netzlos ist. Das
trifft auf `links`/`anchors`/`ids`/`matrix`/`codepaths`/`spans`/`hostpaths` (fixe
Aktiv-Menge) und `planning` (repo-bewusster Block, aktiv bei vorhandener Roadmap)
zu. **Alle übrigen** Module bleiben inaktiv: `external` (Netz, K4), `diagrams`
(nicht ableitbare `patterns`/`defined-in`, K2), `versions`/`targets`
(repo-spezifische `pin-pattern`/`authority`, K2 — bewusst vertagt),
`pins`/`immutable` (pro-Marker, nichts zu deklarieren), `tracked` (fail-closed
ohne `.git`), `citations` (pro-Marker wie `pins`), `structure` (die Regel-Liste
ist repo-spezifisch und nicht ableitbar, K2 — jede Regel benennt ihre
Dokumentklasse selbst), `sources` (Netz, K4); `vcs`/`commits` brauchen eine Commit-Range (K3) und werden als
Makefile-Target über [`--print-mk`](#dc-fa-cli-010--makefile-fragment-ausgeben)
verteilt statt ins statische `modules` aufgenommen. Die Vorlage nennt die nicht
aktivierten Module in einem Kommentar mit Verweis auf
[`--print-config`](#dc-fa-cli-005--konfigurations-gerüst-ausgeben) (Auffindbarkeit
ohne stilles Aktivieren eines inerten Moduls).

**Zwei Modi** (Henne-Ei — der passende Modus ist nicht aus der
Repo-Existenz ableitbar, daher explizit gewählt):

- **`ai-harness-init` (Voll-Kanon):** **alle** Blöcke aktiv, ohne
  Existenzprüfung und ohne Auskommentieren — das vollständige Regelwerk
  als Zielbild für ein **leeres/frisches** Repo. Es läuft erst, wenn die
  Struktur existiert (eine deklarierte Scan-Wurzel bzw. ein `ids`-`target`,
  das fehlt, ist beim Lauf ein Nutzungsfehler —
  [`DC-FA-SCAN-001`](#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln) /
  [`DC-FA-ID-001`](#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)):
  man legt die Struktur an und wächst hinein. (Carveout ohne festes
  Target bleibt auch hier auskommentiert.)
- **`ai-harness` (repo-bewusst):** nur Blöcke, deren Pfad/Target im
  gescannten Baum existiert, sind aktiv; fehlende erscheinen
  **auskommentiert mit Hinweis** (statt stillem Weglassen). Läuft sofort
  gegen das **bestehende/teil-konforme** Repo; die Kommentare sind die
  TODO-Liste Richtung Konformität.

**Anforderungs-Präfix (`--id-prefix`).** Das Anforderungs-`ids`-Muster ist
projektspezifisch: nur sein **Präfix** (für d-check `DC`) wechselt pro Repo;
die `ADR-`/`MR-`/`slice`/Carveout-Muster sind konventions-fest. Das Präfix
ist **parametrisierbar** statt fix `DC`: `--id-prefix <PREFIX>` (z. B. `AC`)
setzt es explizit; im Modus `ai-harness` wird es ohne Angabe aus dem
Lastenheft **abgeleitet** — das eindeutige Projekt-Präfix der FA-/QA-Kennungen
(mehrere verschiedene ⇒ Nutzungsfehler, der Mensch gibt `--id-prefix` an). Ist
kein Präfix angegeben **und** keins ableitbar (insbesondere `ai-harness-init`
fürs leere Repo), wird ein **markierter Platzhalter `<PREFIX>` mit
TODO-Hinweis** emittiert — **kein** stiller `DC-`-Default (der wäre in einem
Fremd-Repo schlicht falsch).

Beide sind **read-only und advisory** (der Mensch verengt), deterministisch
([`DC-QA-02`](#dc-qa-02--determinismus)), mit echten Quellen kombinierbar
(`…,<quelle>`); der aktive Teil dekodiert über den eigenen Parser
([`DC-FA-CONF-001`](#dc-fa-conf-001--konfigurationsdatei)). Ein
Kommentar-Header nennt die zugrunde gelegte Baseline-Version.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine Autoritäts-Quelle mit definierten Kennungen (Headings `ADR-0042`, `ADR-0099`), when `d-check --suggest-config docs/plan/adr/` läuft, then enthält die stdout-Ausgabe ein `ids`-Muster, dessen `regex` `ADR-0042` und `ADR-0099` matcht, mit `target: docs/plan/adr/`, Exit 0. <!-- d-check:ignore (fiktive Beispiel-IDs) -->
- **Boundary:** Given eine benannte Quelle ohne Kennungs-Headings, when `--suggest-config` läuft, then entsteht für sie kein `ids`-Muster (sondern ein Hinweis-Kommentar), kein Absturz, und das Repository wird nicht geschrieben.
- **Negative:** Given eine nicht existierende Quelle, when `d-check --suggest-config gibt/es/nicht` läuft, then Exit-Code 2 (Nutzungsfehler); ein read-only gemountetes Repository genügt (kein Schreibzugriff).
- **`ai-harness` Happy:** Given ein Repo mit `docs/plan/adr/` und `spec/lastenheft.md`, when `d-check --suggest-config ai-harness` läuft, then enthält die stdout-Ausgabe ein vom eigenen Parser akzeptiertes `.d-check.yml` mit den kanonischen `ids`-Mustern (u. a. ein `regex`, der `ADR-`-Kennungen matcht, `target: docs/plan/adr/`), den `matrix`-Klassen samt Referenzrichtungs-Regel und dem Standard-Modulset, Exit 0.
- **`ai-harness` Boundary:** Given ein Repo **ohne** `docs/plan/adr/`, when `d-check --suggest-config ai-harness` läuft, then erscheint der ADR-bezogene Block (zugehöriges `ids`-Muster und die `matrix`-Klasse für ADRs) **auskommentiert mit Hinweis** statt aktiv, `scan.roots` enthält nur existierende Pfade, und das Repository wird nicht geschrieben.
- **`ai-harness` Abgrenzung:** Given es existiert keine Datei und kein Verzeichnis namens `ai-harness`, when `d-check --suggest-config ai-harness` läuft, then wird `ai-harness` als **reservierter Modus** erkannt (nicht als fehlende Quelle → **kein** Exit-Code 2 nach der Negative-Regel oben), Exit 0; ein read-only gemountetes Repository genügt.
- **`ai-harness-init` Voll-Kanon:** Given ein Repo **ohne** `docs/plan/adr/`, when `d-check --suggest-config ai-harness-init` läuft, then sind das ADR-`ids`-Muster und die `matrix`-Klasse `adr` **aktiv** (nicht auskommentiert — im Gegensatz zu `ai-harness`), der aktive Teil dekodiert über den eigenen Parser (das Decoding prüft keine Target-Existenz), Exit 0.
- **`--id-prefix` Happy:** Given `d-check --suggest-config ai-harness-init --id-prefix AC`, when der Lauf endet, then matcht das Anforderungs-`ids`-Muster `AC-FA-…`/`AC-QA-…` und **nicht** `DC-…`, Exit 0.
- **`--id-prefix` Ableitung:** Given ein Repo, dessen `spec/lastenheft.md`-FA-/QA-Kennungen das Präfix `AC` tragen, when `d-check --suggest-config ai-harness` ohne `--id-prefix` läuft, then ist das Anforderungs-Muster auf `AC` abgeleitet, Exit 0.
- **`--id-prefix` Boundary (Platzhalter):** Given `d-check --suggest-config ai-harness-init` ohne `--id-prefix` und ohne ableitbares Präfix, when der Lauf endet, then enthält die Ausgabe den Platzhalter `<PREFIX>` mit TODO-Hinweis (kein `DC-`) und dekodiert über den eigenen Parser, Exit 0.
- **`--id-prefix` Negative (Konflikt):** Given ein `spec/lastenheft.md` mit **mehreren** verschiedenen FA-/QA-Präfixen, when `d-check --suggest-config ai-harness` ohne `--id-prefix` läuft, then Exit-Code 2 (Nutzungsfehler), kein Gerüst.
- **`ai-harness` Auffindbarkeit:** Given `d-check --suggest-config ai-harness` (oder `ai-harness-init`), when der Lauf endet, then nennt die Ausgabe die nicht aktivierten situativen opt-in-Module (`external`, `diagrams`, `versions`, `pins`, `immutable`, `tracked`, `targets`, `structure`) in einem Kommentar mit Verweis auf `--print-config` und verweist für die Range-Module `vcs`/`commits` auf `--print-mk` (Auffindbarkeit ohne Aktivierung), Exit 0; das Repository wird nicht geschrieben.
- **`ai-harness` Modulset (`spans`/`hostpaths`):** Given ein Repo mit den Harness-Artefakten, when `d-check --suggest-config ai-harness` läuft, then enthält die `modules`-Zeile `spans` und `hostpaths` (Teil der fixen Aktiv-Menge), das Gerüst dekodiert über den eigenen Parser, Exit 0.
- **`ai-harness` `planning` repo-bewusst:** Given ein Repo **mit** `docs/plan/planning/in-progress/roadmap.md`, when `d-check --suggest-config ai-harness` läuft, then ist `planning` in `modules` **und** ein `planning`-Block mit `roadmap:` aktiv; **fehlt** die Roadmap, sind Modul und Block **auskommentiert** (repo-bewusst); im Modus `ai-harness-init` ist `planning` **aktiv** (Voll-Kanon). Exit 0.

**Out-of-Scope:** Schreiben der Datei (immer stdout); Muster-Ableitung aus beliebigem Fließtext (nur aus definierten Headings benannter Autoritäts-Quellen); Garantie eines minimalen/perfekten `regex` (Best-Guess-Generalisierung + Quell-Kennungs-Kommentar — der Mensch verengt); automatisches Ableiten von `link-policy`, `matrix`-Regeln oder `exempt-paths`. Die reservierten Quellen `ai-harness`/`ai-harness-init` sind hiervon ausgenommen — sie liefern `matrix`-Regeln, `link-policy` und `exempt-paths` aus der **bekannten Konvention** (nicht aus dem Repo abgeleitet); sie sind an **eine** adoptierte Baseline-Version gebunden (im Kommentar-Header genannt) — automatische Erkennung oder Hebung der Baseline-Version, das Anlegen der Verzeichnisstruktur (read-only) und das automatische Aktivieren der nicht qualifizierten opt-in-Module (`external`, `diagrams`, `versions`, `pins`, `immutable`, `vcs`, `commits`, `tracked`, `targets`) sind Out-of-Scope — die nicht aktivierten situativen Module nennt die Ausgabe stattdessen im Kommentar mit Verweis auf [`--print-config`](#dc-fa-cli-005--konfigurations-gerüst-ausgeben) bzw. für die Range-Module `vcs`/`commits` auf [`--print-mk`](#dc-fa-cli-010--makefile-fragment-ausgeben) (Auffindbarkeit ohne stilles Aktivieren eines inerten Moduls). Aufgenommen sind dagegen `spans`/`hostpaths` (fixe Aktiv-Menge) und der repo-bewusste `planning`-Block (Eignungs-Kriterium K1–K4, s. o.).

---

### DC-FA-CLI-007 — Diagnose-Modus

**Beschreibung:** Mit der Option `--doctor` macht `d-check` einen
**Lese**-Durchgang wie eine normale Prüfung, gibt aber statt der knappen
Befund-Zeilen ([`DC-FA-CLI-004`](#dc-fa-cli-004--ausgabeformate)) eine
**erklärende, nach Datei und Regel gruppierte Diagnose** auf stdout aus:
je Befund den Grund-Code in Klartext, die **Erläuterung** des Befunds
(das Befund-Feld `message`) als eigene
`Hinweis:`-Zeile, wo sie gesetzt ist, und — wo aus dem Befund **eindeutig
ableitbar** — einen **Fix-Kandidaten** (die vorgeschlagene Änderung,
**nicht angewendet**). **Hinweis und Fix-Kandidat sind verschiedene Größen:**
jener ist **verfasst** und wird nie angewendet, dieser ist **abgeleitet** und
wird von [`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch) zum anwendbaren
Patch gerendert. Das Werkzeug schreibt niemals selbst
(read-only-Kernvertrag
[`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
Die Fix-Kandidaten entstehen aus derselben Mechanik, die
[`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch) zum anwendbaren Patch
rendert — eine Quelle, mehrere Ausgaben (Prosa, Patch, JSON). Die Diagnose ist deterministisch
([`DC-QA-02`](#dc-qa-02--determinismus)); die Exit-Codes folgen
[`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes) (die Diagnose erscheint auf
stdout unabhängig vom Code). Mit zusätzlichem `--json` wird die Diagnose
**maschinenlesbar** ausgegeben: ein JSON-Dokument wie
[`DC-FA-CLI-004`](#dc-fa-cli-004--ausgabeformate), dessen `findings` je
Eintrag zusätzlich `reasonText` (Grund-Klartext) und `fixCandidate`
(Objekt mit `original`/`replacement`/`note`, sonst `null`) tragen; die
Gruppierung nach Datei trägt das `file`-Feld (keine Prosa-Einrückung). Es
ist dieselbe Ableitung, ein drittes Rendering neben Prosa und Patch.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Repo mit aktivem Modul `ids` und einem `id-unlinked`-Befund (eine nackte Kennung, deren Definition bekannt ist), when `d-check --doctor` läuft, then enthält stdout eine Diagnose-Gruppe für die betroffene Datei mit dem Grund in Klartext und einem Fix-Kandidaten (Kennung → Link auf ihre Definition), Exit-Code 1.
- **Boundary:** Given ein Repo ohne Befunde, when `d-check --doctor` läuft, then Exit-Code 0 und eine Diagnose, die „0 Befunde" ausweist, ohne Fix-Kandidaten.
- **JSON-Variante:** Given derselbe `id-unlinked`-Befund, when `d-check --doctor --json` läuft, then ist stdout als JSON parsbar (keine Nicht-JSON-Zeilen), der betroffene `findings`-Eintrag trägt `reasonText` und ein `fixCandidate` mit `replacement`, und `exitCode` ist 1.
- **Negative:** Given die Kombination `d-check --doctor --repair`, when aufgerufen, then Exit-Code 2 (Nutzungsfehler, vgl. [`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)) und keine Ausgabe auf stdout.

**Out-of-Scope:** Anwenden der Fix-Kandidaten (das leistet als Patch [`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch)); Fix-Kandidaten für Befunde ohne eindeutige Ableitung (sie werden erklärt, aber ohne Vorschlag); eine JSON-Variante des **Patches** (`--repair --json` bleibt Nutzungsfehler).

---

### DC-FA-CLI-008 — Reparatur-Patch

**Beschreibung:** Mit der Option `--repair` gibt `d-check` einen
**unified diff auf stdout** aus, der ableitbare Befunde behebt —
`git apply`-kompatibel; das Werkzeug schreibt selbst **nichts**
(read-only-Kernvertrag
[`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit);
der Aufrufer wendet an: `d-check --repair > fix.patch`, dann
`git apply fix.patch`). Es gibt **zwei Stufen**, per Schalter wählbar:
die **konservative** Stufe (Default) emittiert nur Hunks für Befunde mit
**eindeutig** ableitbarem Fix (in dieser Version insbesondere
`id-unlinked` → Link auf die bekannte Definition); die **breite** Stufe
(opt-in) schließt **Best-Guess**-Reparaturen ein (z. B. `target-missing`
→ nächstliegende Überschrift bzw. Datei), die als **review-pflichtig**
gekennzeichnet werden — die Kennzeichnung erscheint auf **stderr** (wie
Diagnose/Zusammenfassung, [`DC-FA-CLI-004`](#dc-fa-cli-004--ausgabeformate)),
damit der Patch auf stdout `git apply`-rein bleibt. Befunde ohne Fix in der
gewählten Stufe bleiben unangetastet und erscheinen unter
[`DC-FA-CLI-007`](#dc-fa-cli-007--diagnose-modus). Der Patch ist
deterministisch ([`DC-QA-02`](#dc-qa-02--determinismus)): gleicher Input
und gleiche Stufe → byte-identischer Patch; er ist gegen den
unveränderten Arbeitsbaum sauber anwendbar.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Repo mit einer nackten Kennung (`id-unlinked`, Definition bekannt), when `d-check --repair` in der konservativen Stufe läuft, then liegt auf stdout ein unified diff, der genau diese Zeile in die verlinkte Form ändert (Exit-Code 1, da der Befund im gescannten, ungepatchten Baum besteht), `git apply` nimmt ihn sauber an, und ein erneuter `d-check`-Lauf auf dem gepatchten Baum meldet den Befund nicht mehr.
- **Boundary:** Given ein Repo, dessen einzige Befunde nur best-guess-fähig sind (z. B. `target-missing` — kein eindeutiger Fix), when `d-check --repair` konservativ läuft, then ist der Patch leer (keine Hunks) und das Repository bleibt ungeschrieben; in der breiten Stufe erscheint derselbe Befund als Best-Guess-Hunk (nächstliegende Überschrift bzw. Datei), dessen review-pflichtig-Kennzeichnung auf stderr erscheint.
- **Negative:** Given einen unbekannten Wert für die Stufen-Wahl (bzw. die Kombination `d-check --repair --json`), when aufgerufen, then Exit-Code 2 (Nutzungsfehler, vgl. [`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)) und kein Patch; ein read-only gemountetes Repository genügt (kein Schreibzugriff, [`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

**Out-of-Scope:** In-place-Schreiben der Dateien durch das Werkzeug selbst (immer stdout-Patch; ein In-place-Modus wäre ein Bruch von [`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) und bräuchte eigene Anforderung samt ADR); eine JSON-Variante des Patches; die Garantie, dass die breite (Best-Guess-)Stufe semantisch korrekte Ziele trifft (deshalb review-pflichtig). **Reparierbar sind nur eindeutig ableitbare Befunde:** in dieser Version `id-unlinked` (konservativ) und `target-missing` als eindeutiger Datei-Move (breit) — alle übrigen Befundarten (`anchor-missing`, `repo-escape`, `symlink`, `codepath-missing`, `matrix-inactive`, `matrix-forbidden`, `external-status`, `external-timeout`, `external-redirects`, `span-unclosed`, `span-nested-link`, `hostpath-forbidden`) liefern keinen Edit und bleiben Befund unter [`DC-FA-CLI-007`](#dc-fa-cli-007--diagnose-modus). **Move-/Rename-Erkennung über Datei-Historie:** der `target-missing`-Best-Guess matcht ausschließlich über einen im Scan-Bestand eindeutigen, gleichen Basisnamen einer Markdown-Datei — er erkennt also nur Verschiebungen, keine Umbenennungen und keine Nicht-Markdown-Ziele; eine VCS-/git-historienbasierte Erkennung ist ausgeschlossen, weil sie die Eingabe über den gescannten, read-only gemounteten Baum hinaus erweitern würde ([`DC-QA-02`](#dc-qa-02--determinismus): gleiche Eingabe ⇒ gleiche Ausgabe) und ein eigenes opt-in-Modul (analog `external`) als eigene Anforderung bräuchte.

---

### DC-FA-CLI-009 — Requirements Traceability Matrix

**Beschreibung:** Mit `--trace` macht `d-check` einen **Lese**-Durchgang
über die kanonischen Quellen und gibt eine **Requirements Traceability
Matrix** (RTM) auf **stdout** aus — read-only
([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)),
deterministisch ([`DC-QA-02`](#dc-qa-02--determinismus)), **kein Dokument
erzeugt**, immer frisch abgeleitet. Je Anforderung (Anforderungs-Kennung im
Lastenheft; Default-Gestalt `<PREFIX>-FA-*`/`-QA-*`, per
`trace.requirements.id-pattern` konfigurierbar; Definitionsformat gemäß
[`DC-FA-REQ-001`](#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen))
zeigt die Matrix Titel, die
referenzierenden **ADRs** und **Slices** sowie eine **Lücken-Markierung**
(Anforderung ohne referenzierenden Slice **und** ohne Coverage = Waise). Ist
**mindestens eine** kuratierte Coverage-Quelle konfiguriert (opt-in
`trace.coverage`, [`DC-FA-COV-001`](#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)),
trägt die Matrix zusätzlich eine **Coverage**-Spalte (die Quell-Labels);
**ohne** `trace.coverage` erscheint **keine** Coverage-Spalte und die RTM ist
byte-identisch zur Fassung vor dieser Erweiterung
([`DC-QA-02`](#dc-qa-02--determinismus)). Analog trägt sie bei aktivem opt-in
`trace.requirements.modality` eine **Modality**-Spalte (RFC-2119-Stufe je
Anforderung, [`DC-FA-MOD-001`](#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in));
ohne sie **keine** Modality-Spalte, byte-identisch. Default-Format ist eine
**Markdown-Tabelle**; mit `--json`/`--yaml` wird dieselbe Matrix
strukturgleich maschinenlesbar (je Anforderung ein `coverage`-Feld, das ohne
Coverage-Referenz entfällt — byte-identisch ohne aktive Quelle; format-neutraler
Reporter, [`DC-FA-CLI-004`](#dc-fa-cli-004--ausgabeformate)). **Doku-Domäne:**
Anforderungen aus `spec/lastenheft.md`, Referenzen aus `docs/plan/adr/`
und `docs/plan/planning/` (Default-Pfade; kein Code/keine Go-Toolchain).

**Konfigurierbare Quellen (opt-in `trace`-Block):** Die Konventions-
Annahmen der RTM — die Quell-Datei der Anforderungen samt ihrer
Kennungs-Gestalt und ihres Definitionsformats sowie je Referenzklasse
(ADR, Slice) das Verzeichnis, die
Dateinamen-Gestalt und das Owner-Präfix — sind über einen `trace`-Block in
`.d-check.yml` überschreibbar (`requirements.source`/`.id-pattern`;
`adrs.dir`/`.file-pattern`/`.id-prefix`;
`slices.dir`/`.file-pattern`/`.id-prefix`). **Jedes Feld ist optional; jeder
abwesende Wert fällt auf d-checks eigene Konvention zurück** — eine Config
ohne `trace`-Block (oder mit leerem Block) liefert eine **byte-identische**
RTM ([`DC-QA-02`](#dc-qa-02--determinismus)). So bildet die RTM auch
Konsumenten-Repos mit abweichender Kennungs-/Datei-Konvention **vollständig**
ab, statt nur die Familie zu treffen, die d-checks Schema zufällig teilt.
Reiht sich in die
Advisory-Modi (`--print-config`/`--suggest-config`/`--doctor`) ein; nicht
mit `--doctor`/`--repair` kombinierbar (Nutzungsfehler, Exit 2).

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Repo mit einer Lastenheft-Anforderung, einer sie referenzierenden ADR und einem sie referenzierenden Slice, when `d-check --trace` läuft, then enthält die Markdown-RTM eine Zeile für die Anforderung mit der ADR- und der Slice-Kennung und Status „ok", Exit 0; ein read-only gemountetes Repository genügt.
- **Boundary:** Given ein Repo ohne `spec/lastenheft.md` (oder ohne Anforderungs-Kennungen) und **ohne explizit konfigurierte** `trace.requirements.source` bzw. Tabellenformat, when `--trace` läuft, then eine leere Matrix (0 Anforderungen), kein Absturz, Exit 0, kein Schreibzugriff.
- **Negative:** Given `d-check --trace --repair`, when aufgerufen, then Exit-Code 2 (Nutzungsfehler, [`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)), keine RTM.
- **Konfiguriert (Fremd-Konvention):** Given ein Repo, dessen Anforderungs-Kennungen und Slice-Dateinamen einer von d-checks Default abweichenden Konvention folgen (z. B. `GG-QA-001` und `063-titel.md`), und eine `.d-check.yml` mit `trace.requirements.id-pattern` und `trace.slices.file-pattern` für diese Konvention, when `d-check --trace` läuft, then listet die RTM **alle** so definierten Anforderungen (nicht nur die zufällig zur Default-Gestalt passenden) mit ihren referenzierenden ADRs/Slices, Exit 0; ein read-only gemountetes Repository genügt.
- **Default byte-identisch:** Given **kein** `trace`-Block (oder ein leerer) in der Konfiguration, when `d-check --trace` läuft, then ist die RTM byte-identisch zur Fassung vor dieser Erweiterung ([`DC-QA-02`](#dc-qa-02--determinismus)) und es wird nichts geschrieben ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **Negative (Config):** Given ein `trace`-Block mit einem ungültigen Regex (`id-pattern`/`file-pattern`) **oder** einer `file-pattern` ohne Capture-Gruppe, when `d-check --trace` läuft, then Exit-Code 2 (Konfigurationsfehler, [`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)) mit erklärender Meldung, keine RTM (fail-closed).

**Out-of-Scope:** Erzeugen/Schreiben eines RTM-Dokuments (immer stdout); Code-/Test-Abdeckung (Scan von Go-Testdateien o. Ä. — verlässt die Markdown-Domäne und bräuchte eine Go-Toolchain, unvereinbar mit dem I/O-armen distroless-Design); Status-/Aktualitäts-Bewertung der referenzierenden ADRs (z. B. superseded); **Ableiten** der `trace`-Config aus dem Repo (der Block wird explizit gesetzt, nicht wie `--suggest-config` erraten); ein leerer **ADR**-Owner-Präfix (leerer Wert = Default `ADR-`; der Slice-Owner-Präfix ist per Default leer und bleibt setzbar); mehr als je eine ADR- und eine Slice-Referenzklasse; VCS-/git-historienbasierte Referenz-Erkennung (bliebe außerhalb des read-only Markdown-Baums, [`DC-QA-02`](#dc-qa-02--determinismus)).

---

### DC-FA-CLI-010 — Makefile-Fragment ausgeben

**Beschreibung:** Mit `--print-mk` gibt `d-check` ein **include-bares
Makefile-Fragment** (`d-check.mk`) auf **stdout** aus — kein Repo-Zugriff,
kein Schreiben (read-only-Kernvertrag
[`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit); wie
[`DC-FA-CLI-005`](#dc-fa-cli-005--konfigurations-gerüst-ausgeben)). Das
Fragment trägt eine überschreibbare Image-Variable
`DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v<version>` — die **eingebettete
Release-Version** des laufenden Binaries (das Binary kennt seine Version,
nicht seinen eigenen Digest; für strikte Reproduzierbarkeit überschreibt der
Konsument `DCHECK_IMAGE` mit einem `@sha256:`-Digest aus den Release-Notes,
konsistent mit der Konsum-Pin-Politik aus
[`DC-FA-DIST-001`](#dc-fa-dist-001--docker-image)) — sowie zwölf
`##`-annotierte Targets: `doc-check` (Doku-Gate), `doc-trace` (advisory RTM,
[`DC-FA-CLI-009`](#dc-fa-cli-009--requirements-traceability-matrix)), `doc-complete`
(Vollständigkeits-Gate,
[`DC-FA-CLI-011`](#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)),
`doc-doctor` ([`DC-FA-CLI-007`](#dc-fa-cli-007--diagnose-modus)-Diagnose),
`doc-repair` ([`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch)-Patch — Recipe-Echo
unterdrückt, damit stdout `git apply`-rein bleibt), `doc-immutable`
(git-Diff-Immutabilität via Modul `vcs`,
[`DC-FA-VCS-001`](#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
— `--enable vcs` mit auf `vcs` fokussierter `--disable`-Liste, `RANGE=base..head`
bzw. `STAGED=1`; die **verteilte** Form der git-Garantie für Konsumenten, ohne
Skript-Kopie), `doc-commits` (Commit-Message-Traceability via Modul `commits`,
[`DC-FA-COMMITS-001`](#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in)
— `--enable commits` mit auf `commits` fokussierter `--disable`-Liste, `--range $(RANGE)`;
die **verteilte** Commit-Traceability für Konsumenten ohne Skript-Kopie, parallel zu
`doc-immutable`), `doc-planning` (Planning-Lifecycle-Konsistenz via Modul `planning`,
[`DC-FA-PLAN-001`](#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
— `--enable planning` mit auf `planning` fokussierter `--disable`-Liste, **hermetisch**
ohne `RANGE`/`STAGED`), `doc-tracked` (Getrackt-Status auflösbarer Referenz-Ziele
via Modul `tracked`,
[`DC-FA-TRK-001`](#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in)
— `--enable tracked` mit auf `tracked` fokussierter `--disable`-Liste, ohne
`RANGE`/`STAGED`; liest das `.git` im read-only-Mount), `doc-targets`
(Deklarations-Konsistenz Doku ↔ Build-Targets via Modul `targets`,
[`DC-FA-TGT-001`](#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
— `--enable targets` mit auf `targets` fokussierter `--disable`-Liste,
**hermetisch** ohne `RANGE`/`STAGED`), `doc-structure` (Struktur-Invarianten
**innerhalb** eines Dokuments via Modul `structure`,
[`DC-FA-STRUCT-001`](#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
— `--enable structure` mit auf `structure` fokussierter `--disable`-Liste,
**hermetisch** ohne `RANGE`/`STAGED`) und `doc-help` (listet die
`doc-*`-Targets), jeweils `docker run --network none -v "$PWD:/repo:ro"`. Dazu die
Variablen `TRACE_FLAGS` (Flags der RTM-Targets) und `DCHECK_DIGEST` (ein
`@sha256:`-Digest, der den Tag von `DCHECK_IMAGE` sticht — strikte
Reproduzierbarkeit, ohne den vollen Ref zu überschreiben).
Konsumenten `include d-check.mk` und legen ihre eigene `.d-check.yml`
daneben — **keine Recipe-/Skript-Kopie**, der Image-Pin lebt in d-check.
Deterministisch ([`DC-QA-02`](#dc-qa-02--determinismus)): hängt nur an der
eingebetteten Version. Reiht sich in die read-only-Generatoren
(`--print-config`/`--suggest-config`) ein.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein installiertes `d-check`, when `d-check --print-mk` läuft, then liegt auf stdout ein Makefile-Fragment mit `DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v…`, den Variablen `DCHECK_DIGEST` und `TRACE_FLAGS` und den `##`-annotierten Targets `doc-check`, `doc-trace`, `doc-complete`, `doc-doctor`, `doc-repair`, `doc-immutable` (mit `--enable vcs` und einer auf `vcs` fokussierten `--disable`-Liste), `doc-commits` (mit `--enable commits` und einer auf `commits` fokussierten `--disable`-Liste), `doc-planning` (mit `--enable planning` und einer auf `planning` fokussierten `--disable`-Liste, ohne Range), `doc-tracked` (mit `--enable tracked` und einer auf `tracked` fokussierten `--disable`-Liste, ohne Range), `doc-targets` (mit `--enable targets` und einer auf `targets` fokussierten `--disable`-Liste, ohne Range) und `doc-structure` (mit `--enable structure` und einer auf `structure` fokussierten `--disable`-Liste), `doc-help` (jeweils `docker run … --network none … :/repo:ro`), Exit 0; kein Repo-Zugriff nötig. **Der Kopfkommentar nennt die URL des Benutzerhandbuchs** — das Fragment reist in ein **fremdes** Repo, und sein Kopf ist der einzige Ort, an dem ein Zeiger dorthin dauerhaft mitfährt; auch sie zeigt auf den **Hauptzweig**, ohne Versionsangabe.
- **Boundary:** Given das Fragment wird per `d-check --print-mk > d-check.mk` umgeleitet und in ein Makefile `include`-t, when `make doc-check` (bzw. `doc-trace`/`doc-complete`/`doc-doctor`/`doc-repair`/`doc-immutable`/`doc-commits`/`doc-planning`/`doc-tracked`/`doc-targets`/`doc-structure`) läuft, then ruft das Target das gepinnte Image (bzw. den per `DCHECK_IMAGE`/`DCHECK_DIGEST` gesetzten Override) im passenden Modus (`doc-trace` → `--trace`, `doc-complete` → `--trace --require-complete`, `doc-doctor` → `--doctor`, `doc-repair` → `--repair` mit unterdrücktem Recipe-Echo, `doc-immutable` → `--enable vcs` + Fokus-`--disable` mit `--range $(RANGE)` bzw. `--staged`, `doc-commits` → `--enable commits` + Fokus-`--disable` mit `--range $(RANGE)`, `doc-planning` → `--enable planning` + Fokus-`--disable`, **hermetisch ohne Range**, `doc-tracked` → `--enable tracked` + Fokus-`--disable`, ohne Range, `doc-targets` → `--enable targets` + Fokus-`--disable`, **hermetisch ohne Range**); d-check selbst schreibt dabei nichts.
- **Negative:** Given `d-check --print-mk` mit einem unbekannten Flag, when aufgerufen, then Exit-Code 2 (Nutzungsfehler, [`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)).

**Out-of-Scope:** Schreiben der `d-check.mk` (immer stdout); Einbetten des eigenen Image-**Digests** ins Binary (Henne-Ei — der Digest hasht das Binary selbst; der Konsument pinnt per `DCHECK_IMAGE`-Override); weitere Targets jenseits der gelisteten zwölf (`doc-check`/`doc-trace`/`doc-complete`/`doc-doctor`/`doc-repair`/`doc-immutable`/`doc-commits`/`doc-planning`/`doc-tracked`/`doc-targets`/`doc-structure`/`doc-help` — Konsumenten komponieren weitere `gates` selbst); ein `help`-Target (Namens-Kollision mit dem Konsumenten — daher namespaced `doc-help`); Nicht-`@sha256:`-Digest-Formen in `DCHECK_DIGEST`; die Exit-Code-Semantik der RTM-Targets selbst (in [`DC-FA-CLI-011`](#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code) bzw. [`DC-FA-CLI-009`](#dc-fa-cli-009--requirements-traceability-matrix) festgelegt); Nicht-Make-Build-Systeme.

---

### DC-FA-CLI-011 — Vollständigkeits-Prüfung als opt-in Exit-Code

**Beschreibung:** Mit `--trace --require-complete` macht `d-check` denselben
read-only RTM-Lauf wie
[`DC-FA-CLI-009`](#dc-fa-cli-009--requirements-traceability-matrix)
(Anforderung → ADRs/Slices, Waisen-Markierung; read-only
[`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit),
deterministisch [`DC-QA-02`](#dc-qa-02--determinismus), kein Dokument erzeugt),
**bindet das Ergebnis aber an den Exit-Code**: ≥1 Requirements-Waise
(Anforderung ohne referenzierenden Slice **und** — falls `trace.coverage`
konfiguriert ist — ohne Coverage-Referenz, [`DC-FA-COV-001`](#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in))
⇒ **Exit 1** (Befund-Code,
[`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)); 0 Waisen ⇒ Exit 0. Ist zusätzlich
`trace.requirements.modality` aktiv ([`DC-FA-MOD-001`](#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in)),
gatet eine Waise **nur**, wenn ihre Modalitäts-Stufe in `require-levels` liegt
(Default `[must]`) — SOLLTE/KANN/`unknown`-Waisen bleiben advisory. Ohne
`modality` gaten wie bisher **alle** Waisen. Die RTM
erscheint unverändert auf stdout (Default Markdown, mit `--json`/`--yaml`
maschinenlesbar) — `--require-complete` ändert nur den Exit-Code, nicht die
Ausgabe. Es ist ein **Modifikator von `--trace`**: ohne `--trace` ein
Nutzungsfehler (Exit 2,
[`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)). Der Default-`--trace` bleibt
**advisory** (Exit 0 bei Waisen,
[`DC-FA-CLI-009`](#dc-fa-cli-009--requirements-traceability-matrix)
unangetastet); die Durchsetzung ist opt-in. So binden Konsumenten die
Vollständigkeits-Invariante als Makefile-Gate (vgl. das `doc-complete`-Target
aus [`DC-FA-CLI-010`](#dc-fa-cli-010--makefile-fragment-ausgeben)), ohne die
RTM-Parsing-Logik zu kopieren.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Repo, dessen Lastenheft-Anforderungen alle von ≥1 Slice referenziert sind, when `d-check --trace --require-complete` läuft, then erscheint die RTM auf stdout und der Exit-Code ist 0; ein read-only gemountetes Repository genügt.
- **Boundary:** Given ein Repo mit genau einer Requirements-Waise (Anforderung ohne referenzierenden Slice und ohne Coverage-Referenz), when `d-check --trace --require-complete` (auch mit `--json`/`--yaml`) läuft, then ist der Exit-Code 1, die vollständige RTM erscheint auf stdout (Waise markiert), und es wird nichts geschrieben ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **Coverage deckt ab:** Given ein Repo mit einer Anforderung **ohne** Slice, aber **mit** einer `trace.coverage`-Referenz, when `d-check --trace --require-complete` läuft, then ist diese Anforderung **keine** Waise; sind alle Anforderungen so abgedeckt, Exit 0.
- **Negative:** Given `d-check --require-complete` ohne `--trace`, when aufgerufen, then Exit-Code 2 (Nutzungsfehler, [`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)), keine RTM.

**Out-of-Scope:** Ändern des Default-`--trace`-Exit-Codes (bleibt advisory Exit 0, [`DC-FA-CLI-009`](#dc-fa-cli-009--requirements-traceability-matrix)); Bewertung des Slice-**Status** (done vs. in-progress — eine von ≥1 Slice beanspruchte Anforderung zählt als abgedeckt, unabhängig vom Fortschritt); Erzwingung anderer RTM-Eigenschaften als Waisenfreiheit (Slice **oder** Coverage genügt; eine bloße ADR-Referenz ohne Slice/Coverage deckt weiterhin **nicht** ab); Schreiben eines Reports.

---

### DC-FA-CLI-012 — Konfigurations-Pfad überschreiben

**Beschreibung:** Die Option `--config <datei>` ersetzt für **diesen** Lauf die
konventionelle Konfigurations-Quelle — die `.d-check.yml` in der Scan-Wurzel
([`DC-FA-CONF-001`](#dc-fa-conf-001--konfigurationsdatei)) — durch die angegebene
Datei. Der Pfad wird relativ zur Scan-Wurzel aufgelöst und muss **innerhalb** der
Scan-Wurzel liegen; damit bleibt die Datei im read-only-Mount lesbar und der Lauf
hermetisch (dieselbe Wurzel-Constraint wie `scan.roots` und `planning.roadmap`).
Die Datei durchläuft **exakt** dieselbe strikte Validierung wie eine
konventionelle `.d-check.yml` (Syntax **und** Inhalt, jeder Fehler Exit 2). Fehlt
sie, ist das ein **Nutzungsfehler** (Exit 2,
[`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)) — es gibt **keinen** stillen
Rückfall auf Defaults und **keinen** auf die konventionelle Datei; ein getipptes
Profil darf nicht unbemerkt zu einem anderen Prüfumfang führen. Die angegebene
Datei **ersetzt**, sie ergänzt nicht.

Der Zweck ist die saubere Trennung von **Prüf-Profilen** an unterschiedlichen
Bindepunkten desselben Repos: ein Profil für den inneren Loop und ein zweites,
enger gefasstes für einen Abschluss-Bindepunkt (Wellen-/Release-Closure) — ohne
die Modul- und Parameterwahl über die Kommandozeile nachzubauen und ohne beide
Profile in einer Datei zu vermischen. Ohne `--config` ist das Verhalten
unverändert und byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)).

**Akzeptanzkriterien:**

- **Happy Path:** Given eine gültige alternative Konfigurationsdatei innerhalb der Scan-Wurzel, when `d-check --config <datei>` läuft, then gilt ausschließlich deren Inhalt, eine vorhandene `.d-check.yml` der Scan-Wurzel wird **nicht** gelesen, und der Exit-Code folgt der Befundlage.
- **Boundary (ohne Flag):** Given **kein** `--config`, when `d-check` läuft, then unverändertes Verhalten (`.d-check.yml` aus der Scan-Wurzel) und byte-identischer Befundsatz ([`DC-QA-02`](#dc-qa-02--determinismus)).
- **Negative (fehlt):** Given einen `--config`-Pfad, der nicht existiert, when `d-check` läuft, then Exit-Code 2 mit Hinweis auf stderr — **kein** Rückfall auf Defaults oder auf die konventionelle Datei.
- **Negative (leerer Wert):** Given `--config` **mit leerem Wert** (typischer Auslöser: eine nicht expandierte Make-/CI-Variable), when `d-check` läuft, then Exit-Code 2 — ein gesetztes Flag ohne Pfad ist ein Nutzungsfehler und **kein** stiller Rückfall auf die konventionelle Datei.
- **Negative (außerhalb der Wurzel):** Given einen `--config`-Pfad außerhalb der Scan-Wurzel, when `d-check` läuft, then Exit-Code 2 — auch dann, wenn die Datei dort existiert und gültig ist.
- **Negative (Symlink-Ausbruch):** Given einen `--config`-Pfad, der über einen **Verzeichnis-Symlink** aus der Scan-Wurzel herausführt, when `d-check` läuft, then Exit-Code 2 mit Nennung des Symlinks — die Wurzel-Grenze ist nicht nur lexikalisch.
- **Negative (kein Datei-Ziel):** Given einen `--config`-Pfad, der auf ein Verzeichnis zeigt, when `d-check` läuft, then Exit-Code 2.
- **Boundary (absolut innerhalb):** Given einen **absoluten** `--config`-Pfad, der innerhalb der Scan-Wurzel liegt, when `d-check` läuft, then wird er auf die Wurzel bezogen aufgelöst und gilt.
- **Provenance:** Given eine `--config`-Datei mit einem Konfigurationsfehler, dessen Meldung die Datei benennt, when `d-check` läuft, then nennt die Meldung die **tatsächlich geladene** Datei, nicht die konventionelle.
- **Negative (ungültig):** Given eine `--config`-Datei mit ungültigem Schema, when `d-check` läuft, then Exit-Code 2 mit Zeilenangabe, wie bei einer ungültigen `.d-check.yml` ([`DC-FA-CONF-001`](#dc-fa-conf-001--konfigurationsdatei)).

**Out-of-Scope:** Zusammenführen oder Vererben mehrerer Konfigurationsdateien (die
angegebene Datei ersetzt vollständig); mehrere `--config`-Angaben in einem Lauf;
Konfigurationen außerhalb der Scan-Wurzel (der read-only-Mount ist die Grenze);
Profil-Namen oder eine Profil-Registry innerhalb einer Datei.

---

### DC-FA-COV-001 — Kuratierte Coverage-Quellen der RTM (`trace.coverage`, opt-in)

**Beschreibung:** Die RTM
([`DC-FA-CLI-009`](#dc-fa-cli-009--requirements-traceability-matrix)) leitet
Abdeckung aus **Referenz-Scans** von ADR-/Slice-Dateien ab. Anforderungen, deren
Abdeckung in einer **kuratierten Matrix** (eine ausgelagerte
Traceability-Datei, außerhalb `adrs`/`slices`) nachgewiesen ist,
erscheinen daher als **falsche Waisen** — die Matrix ist weder
ein ADR noch ein Slice, und sie nutzt oft **Bereichs-Notation**
(`<FAM>-001..006`), die der Scanner nur als erste ID erkennt. `trace.coverage`
ist eine **dritte, opt-in Referenzklasse**: eine **Liste** benannter kuratierter
Quellen, die **range-aware** als Coverage eingelesen werden, **ohne**
`adrs`/`slices` zu berühren.

Je Quelle: **`files`** (explizite Pfad-Liste — **keine** `dir`+`file-pattern`-Ableitung,
damit ADR-/Slice-Dateien nicht als Coverage kontaminieren), **`label`** (feste
Owner-Kennung in der eigenen **Coverage**-Spalte der RTM, z. B. `Trace`), optional
**`ranges`** (bool, Default **true**), optional **`sections`** (Whitelist: nur die
genannten H2/H3-Abschnitte zählen) und **`exclude-sections`** (Blacklist: die
genannten Abschnitte zählen nicht — z. B. eine „…ohne Design-Artefakt"-Liste, die
sonst **nicht** gedeckte IDs fälschlich krediten würde). Beide Abschnitts-Mengen
nutzen dieselbe Heading-Span-Semantik wie `matrix.exclude-sections` (Abschnitt bis
zur nächsten gleich-/höherrangigen Überschrift).

Eine Anforderung gilt als von der Quelle **abgedeckt**, wenn ihre Kennung im
(abschnitts-gefilterten) Text vorkommt — **exakt** (`requirements.id-pattern`)
**oder**, bei `ranges: true`, **range-expandiert**: `<FAM>-AAA..BBB` deckt alle
IDs `AAA…BBB` **inklusive** ab, die Aufzählung `<FAM>-AAA/BBB/CCC` die genannten;
jede expandierte ID muss zum `requirements.id-pattern` passen, sonst wird sie
verworfen. Abgedeckte Anforderungen tragen das `label` der Quelle in der
Coverage-Spalte; eine Anforderung ist **Waise** nur, wenn sie **weder** Slice-
**noch** Coverage-Referenz trägt ([`DC-FA-CLI-011`](#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)).
Read-only ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)),
deterministisch ([`DC-QA-02`](#dc-qa-02--determinismus)). **Fail-closed** (Exit 2):
leere `files`-Liste, ein `files`-Pfad außerhalb der Repo-Wurzel (führendes `/`
oder `..`), eine **fehlende** `files`-Datei (`files` sind **explizit** benannt —
Fehlen ist Fehler, nicht Skip wie bei `adrs.dir`/`slices.dir`), leeres `label`,
ein **Sektionsname ohne Heading-Treffer** (Tippfehler-/Kurzform-Guard), oder eine
Range mit `AAA>BBB` bzw. abweichender Ziffern-Breite; eine **Komma-Kurzform**
(`<FAM>-AAA, BBB` — eine Kennung unmittelbar gefolgt von Komma und Ziffern) ist
**keine** zugesagte Notation und daher ebenfalls Exit 2: d-check kann eine
gemeinte Kurzform nicht von einer Zahl im Fließtext unterscheiden, und beides zu
raten wäre schlimmer als beides abzulehnen. Ein Komma vor einer **vollständigen**
Kennung (`<FAM>-AAA, <FAM>-BBB`) ist unberührt — das ist keine Kurzform. Ein Sektionsname ist der
**volle Heading-Klartext** (exakt wie `matrix.exclude-sections`, z. B.
`"27.1.1 Anforderungen ohne Design-Artefakt"`, nicht die Kurzform). `trace.coverage`
führt **kein** eigenes Regex (nutzt `requirements.id-pattern`). **Strikt opt-in:** ohne
`trace.coverage` ist die RTM byte-identisch (keine Coverage-Spalte, kein
`coverage`-Feld), und `adrs`/`slices`-Semantik bleibt unberührt.

**Akzeptanzkriterien:**

- **Happy Path:** Given `trace.coverage: [{files: [docs/plan/traceability.md], label: Trace}]` und eine Anforderung **ohne** Slice, deren Kennung **exakt** in `traceability.md` steht, when `d-check --trace` läuft, then trägt ihre RTM-Zeile das Label `Trace` in der Coverage-Spalte und ist **keine** Waise, Exit 0; ein read-only gemountetes Repository genügt.
- **Range:** Given `ranges: true` und die Zeile `GG-QA-001..006` in der Quelle (bei `requirements.id-pattern` `GG-…-\d{3}`), when `d-check --trace` läuft, then sind **alle sechs** `GG-QA-001`…`GG-QA-006` als von `Trace` abgedeckt markiert (nicht nur `001`).
- **Sektionen (Blacklist):** Given eine Quelle mit `exclude-sections: ["27.1.1 Anforderungen ohne Design-Artefakt"]`, deren §27.1.1 sonst nicht gedeckte IDs listet, when `d-check --trace` läuft, then werden die **nur** in §27.1.1 genannten IDs **nicht** als Coverage gewertet.
- **Keine ADR-Kontamination:** Given `trace.coverage` mit `files`, when `d-check --trace` läuft, then werden ADR-/Slice-Dateien **nicht** als Coverage gezählt (nur die gelisteten `files`), und die Coverage-Labels stehen in ihrer **eigenen** Spalte (keine Kollision mit Slice-Namen).
- **Negative (Config/Range):** Given eine fehlende `files`-Datei **oder** eine ungültige Range `GG-RT-009..003` (`AAA>BBB`), when `d-check --trace` läuft, then Exit-Code 2 (Konfigurationsfehler, [`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)) mit erklärender Meldung.
- **Negative (Komma-Kurzform):** Given eine Zelle `GG-SCN-001, 007`, when
  `d-check --trace` mit `ranges: true` läuft, then Exit 2 mit erklärender Meldung
  (Verweis auf die zugesagten Notationen `/` und `..`) — **nicht** ein stiller
  Drop von `007` und **nicht** eine geratene Expansion.
- **Boundary (Komma vor voller Kennung):** Given eine Zelle
  `GG-AR-COMP-CORE, GG-AR-COMP-DOMAIN`, when der Lauf läuft, then werden **beide**
  Kennungen gelesen, kein Fehler — die Kurzform-Regel greift nur bei Komma **plus
  Ziffern**.
- **Modul-aus:** Given **kein** `trace.coverage`, when `d-check --trace` läuft, then ist die RTM byte-identisch zur Fassung vor dieser Anforderung ([`DC-QA-02`](#dc-qa-02--determinismus)) — keine Coverage-Spalte, kein `coverage`-Feld — und es wird nichts geschrieben ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

**Out-of-Scope:** Parsen **semantischer Design-Mappings** (z. B. `GG-AR-*`-Architektur-Kennungen) — Coverage ist reine **ID-Präsenz** in der Quelle; ein **Erfüllungs-/Status-Urteil** (Coverage ≠ „erfüllt", nur „in kuratierter Quelle referenziert"; kein `✓`); **Auto-Generierung** oder Schreiben der Coverage-Datei; `dir`+`file-pattern`-Ableitung für Coverage-Quellen (bewusst nur explizite `files`, gegen ADR-Kontamination); range-/enum-Erkennung außerhalb `trace.coverage` (eine spätere `commits`/`ids`-Range wäre eigener CR); fam-qualifizierte Range-Enden (`<FAM>-AAA..<FAM>-BBB` — nur die Kurzform `..BBB`).

---

### DC-FA-MOD-001 — Modalitäts-Klassifikation der Anforderungen (`trace.requirements.modality`, opt-in)

**Beschreibung:** Anforderungen tragen eine **Modalität** (RFC-2119-Stufe:
MUSS/SOLLTE/KANN bzw. MUST/SHOULD/MAY) im Anforderungs-**Text**. Die RTM
([`DC-FA-CLI-009`](#dc-fa-cli-009--requirements-traceability-matrix)) liest heute
nur die Heading-Kennung und behandelt jede Anforderung gleich — eine reine
**KANN**-Anforderung ohne Slice/Coverage erscheint wie eine unabgedeckte
**MUSS**-Pflicht. Bei explizit aktiviertem opt-in `trace.requirements.modality`
klassifiziert d-check jede Anforderung anhand von **Modal-Verb-Schlüsselwörtern**
im Body und macht die Stufe in einer eigenen **Modality**-Spalte sichtbar; die
Vollständigkeits-Prüfung ([`DC-FA-CLI-011`](#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code))
kann so **nur auf verpflichtende Stufen** brechen.

Die Schlüsselwörter sind **konfigurierbar mit Built-in-Defaults**:
`modality.levels` bildet Stufen-Namen → Schlüsselwort-Liste ab (Default eine
DE+EN-RFC-2119-Menge — `must`/`should`/`may`); `modality.require-levels` nennt die
Stufen, die `--require-complete` gaten (Default `[must]`). Klassifiziert wird über
den Schlüsselwort-Treffer an der **frühesten Position** im Body-Abschnitt der
Anforderung (Überschrift bis zur nächsten gleich-/höherrangigen); mehrere Phrasen
an **derselben** Position: **längster** Treffer zuerst (damit `DARF NICHT` vor
`DARF`, und `MUSS NICHT` [may] nicht als `MUSS` [must] verschluckt wird),
Schlüsselwörter case-insensitiv und wortgrenzen-genau (`MUSS` matcht nicht
`musste`). Der Body wird vor dem Matching **normalisiert** (Markdown-Emphasis
entfernt, Whitespace-/Umbruch-Folgen zu einem Leerzeichen) — sonst matchte eine
umbrochene/emphasierte Phrase `**MUSS** NICHT`/`MUSS\nNICHT` nicht und fiele still
auf `MUSS` zurück. Eine Anforderung **ohne** Treffer erhält die Stufe
**`unknown`** (in der Spalte sichtbar); ob `unknown` gatet, entscheidet
`require-levels` (Default: nicht — `unknown` ∉ `[must]`).

**Waise-Interaktion:** `--require-complete` bricht (Exit 1) nur bei Waisen
(¬slice ∧ ¬coverage), **deren Stufe in `require-levels` liegt** (die stderr-Zeile
nennt die gatende von der Gesamt-Waisenzahl). Ist `trace.requirements.modality`
**nicht** gesetzt, gaten wie bisher **alle** Waisen (byte-identisch); **aktiv** ist
`modality` schon bei bloßer **Schlüssel-Präsenz** (`modality: {}` ⇒ Defaults + Spalte).
Read-only ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)),
deterministisch ([`DC-QA-02`](#dc-qa-02--determinismus)). **Fail-closed** (Exit 2):
leerer Stufen-Name/leeres Schlüsselwort; der **reservierte** Stufen-Name
`unknown` in `levels`; **dasselbe Keyword in mehr als einer Stufe** (sonst
nondeterministisch); ein `require-levels`-Eintrag, der weder deklarierter
Stufen-Name noch `unknown` ist. **Strikt opt-in:** ohne `modality` keine
Modality-Spalte, kein `modality`-Feld, `--require-complete`-Semantik unverändert
— byte-identisch.

**Akzeptanzkriterien:**

- **Happy Path:** Given `trace.requirements.modality: {}` (Defaults) und eine Anforderung, deren Body `… MUSS …` enthält, when `d-check --trace` läuft, then trägt ihre RTM-Zeile die Modality `must`, Exit 0; ein read-only gemountetes Repository genügt.
- **Gate nach Stufe:** Given eine **KANN**-Anforderung ohne Slice/Coverage und `require-levels: [must]` (Default), when `d-check --trace --require-complete` läuft, then ist sie **keine** gatende Waise (Exit 0, sofern keine MUSS-Waise); dieselbe Anforderung mit `require-levels: [must, may]` ⇒ Exit 1.
- **Längster Treffer / Negation:** Given ein Body mit `… MUSS NICHT …` (Default-`may` enthält `MUSS NICHT`) und ein zweiter mit `… DARF NICHT …` (Default-`must`), when `d-check --trace` läuft, then wird der erste als `may` (nicht `must`) und der zweite als `must` klassifiziert — auch wenn die Phrase **zeilenumbrochen** (`MUSS\nNICHT`) oder **emphasiert** (`**MUSS** NICHT`) ist (Body-Normalisierung).
- **Unknown:** Given eine Anforderung ohne Modal-Verb (z. B. ein deklaratives Nicht-Ziel), when `d-check --trace` läuft, then Modality `unknown`; sie gatet nur, wenn `require-levels` `unknown` enthält.
- **Negative (Config):** Given eine ungültige `modality`-Config — ein `require-levels`-Eintrag, der weder deklarierte Stufe noch `unknown` ist; **oder** ein leerer Stufen-Name/leeres Keyword; **oder** der reservierte Stufen-Name `unknown` in `levels`; **oder** dasselbe Keyword in zwei Stufen — when `d-check --trace` läuft, then Exit-Code 2 (Konfigurationsfehler, [`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)).
- **Modul-aus:** Given **kein** `trace.requirements.modality`, when `d-check --trace [--require-complete]` läuft, then keine Modality-Spalte, kein Feld, alle Waisen gaten wie bisher — byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)), nichts geschrieben ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

**Out-of-Scope:** Semantische Sprach-Analyse jenseits der Schlüsselwort-Präsenz (keine NLP; „der Satz **verweist** auf eine MUSS-Regel" wird nicht erkannt); Modalität aus dem Heading-Titel statt dem Body; mehrere Modalitäten je Anforderung (genau **eine** Stufe = erster Treffer); automatische Sprach-Erkennung der Defaults (die Default-Menge ist DE+EN-fix, sonst konfiguriert der Mensch); Ableitung von `require-levels` (explizit gesetzt); Bewertung, **ob** die Modalität semantisch korrekt vergeben ist (nur Klassifikation, kein Urteil).

---

### DC-FA-REQ-001 — Anforderungsquellen als Headings oder Tabellen

**Beschreibung:** Die Requirements Traceability Matrix
([`DC-FA-CLI-009`](#dc-fa-cli-009--requirements-traceability-matrix)) liest
Anforderungsdefinitionen standardmäßig aus ATX-Überschriften: Die Kennung muss
als erstes vollständiges Heading-Token auf
`trace.requirements.id-pattern` passen; der nachfolgende Abschnitt liefert
Titel und Body. Zusätzlich kann `trace.requirements.format: table` ein
tabellenbasiertes Lastenheft nativ einlesen. Der zugehörige `table`-Block
benennt die Spalten über ihre Header statt über instabile Positionen:
`id-column` (Kennung), genau eine von `text-column` (ein Text-Header) oder
`text-columns` (explizite alternative Text-Header, von denen jeder mindestens
einmal in der Quelle vorkommen muss) sowie
optional `modality-column` (Modalitätsquelle). d-check scannt Markdown-
Pipe-Tabellen der konfigurierten Quelldatei; eine Datenzeile definiert genau
dann eine Anforderung, wenn ihre ID-Zelle als Ganzes auf `id-pattern` passt.

Ist `trace.requirements.modality` aktiv
([`DC-FA-MOD-001`](#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in)),
klassifiziert d-check im Tabellenformat bei gesetzter `modality-column`
**ausschließlich deren Zelleninhalt** mit denselben konfigurierten Keywords;
ohne `modality-column` dient der Inhalt der jeweils gebundenen Textspalte als Body. Damit ist
etwa `Kennung | Prioritaet | Anforderung` durch `id-column: Kennung`,
`modality-column: Prioritaet` und `text-column: Anforderung` abbildbar. Die
Heading- und Tabellenform liefern danach dasselbe RTM-Modell; Referenzscan,
Coverage, Waisenstatus und Ausgabeformate bleiben identisch.

Die Quellenwahl ist **fail-closed**: Ist `trace.requirements.source`
mit einem **nichtleeren** Wert explizit konfiguriert oder `format: table`
aktiviert und werden **null**
Anforderungen erkannt, endet `--trace` ebenso wie
`--trace --require-complete` mit Exit 2 und einer erklärenden Meldung statt
einer irreführend grünen leeren RTM. Ebenso sind ein unbekanntes `format`, ein
fehlender `table`-Block, leere Spaltennamen, ein in keiner Tabelle vorhandener
oder dort doppelt vorkommender konfigurierter Header sowie standardmäßig eine
doppelte erkannte Anforderungs-ID im Tabellenmodus oder bei nichtleer expliziter
Quelle Konfigurations-/Quelldatenfehler (Exit 2). Für historische
Mehrfachdefinitionen wählt `table.duplicate-ids` explizit `first` oder `last`;
Default `error` bleibt fail-closed. `source: ""` gilt wie bisher als
abwesend und fällt auf den Default zurück. Ohne `trace`-Block bleibt
`format: headings` der Default; eine fehlende Default-Quelle bzw. null
Heading-Treffer liefert aus Kompatibilitätsgründen weiterhin die leere RTM
mit Exit 0 ([`DC-QA-02`](#dc-qa-02--determinismus)). Beide Formate sind
read-only ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

**Akzeptanzkriterien:**

- **Happy Path:** Given eine explizit konfigurierte Tabellenquelle mit `Kennung | Prioritaet | Anforderung` sowie `Kennung | Prioritaet | Akzeptanzkriterium`, `format: table`, `text-columns: [Anforderung, Akzeptanzkriterium]`, 372 passende Datenzeilen und 371 eindeutige IDs (eine historische Hochstufung doppelt), when `duplicate-ids: last` gesetzt ist und `d-check --trace` läuft, then enthält die RTM genau 371 Anforderungen; die spätere Definition gewinnt, Referenzen/Waisen werden wie bei Heading-Definitionen berechnet, Exit 0.
- **Boundary (Modalität):** Given eine Tabellenzeile `F-1 | Muss | Das Repository muss …`, eine konfigurierte `modality-column: Prioritaet` und aktive Default-`modality`, when `d-check --trace` läuft, then ist die Anforderung als `must` klassifiziert; der Text der Anforderung beeinflusst diese Klassifikation nicht. Ohne `modality-column` wird stattdessen `text-column` klassifiziert.
- **Negative (Nullmenge):** Given eine mit einem **nichtleeren** Wert explizit konfigurierte `trace.requirements.source`, deren Format oder Inhalte keine einzige Definition gemäß der gewählten Grammatik und `id-pattern` ergeben, when `d-check --trace [--require-complete]` läuft, then Exit 2 mit Quelle, Format und Hinweis auf null erkannte Anforderungen; keine grüne leere RTM. `source: ""` gilt wie Abwesenheit und behält die Default-Semantik.
- **Negative (Schema):** Given `format: table` ohne genau eine der Alternativen `text-column`/`text-columns`, mit einem fehlenden/leer benannten oder in der Quelle nicht vorhandenen Rollen-Header beziehungsweise mit unbekanntem `duplicate-ids`, when `d-check --trace` läuft, then Exit 2 mit dem betroffenen Feld oder Spaltennamen.
- **Default byte-identisch:** Given kein `trace`-Block oder `trace.requirements.source: ""`, when `d-check --trace [--require-complete]` auf demselben Repo läuft, then bleiben Heading-Erkennung, Deduplizierung, Ausgabe und Exit-Code byte-identisch zur Fassung vor dieser Anforderung; insbesondere bleibt die unkonfigurierte Nullmenge Exit 0 ([`DC-QA-02`](#dc-qa-02--determinismus)).

**Out-of-Scope:** HTML-Tabellen; Tabellen ohne Markdown-Header-/Trennzeile;
zusammengeführte Zellen, mehrzeilige Tabellenzeilen oder Block-Markdown in
einer Zelle; automatische Spalten-Ableitung über Synonyme oder Positionen;
gleichzeitiges Mischen von Heading- und Tabellen-Definitionen in einem Lauf
(genau ein `format`); mehrere Anforderungs-Quelldateien; semantische Bewertung
des Anforderungstextes jenseits der bestehenden Keyword-Klassifikation.

---

### DC-FA-XREF-001 — Kreuzverweis-Konsistenz zweier Traceability-Sichten (`trace.cross-consistency`, opt-in)

**Beschreibung:** Anforderung→Design wird in zwei Sichten geführt: **Rückwärts-
Kanten**, die — dort authort, wo das Design lebt — je Design-Artefakt aufwärts die
umgesetzten Anforderungen nennen (**Quelle der Wahrheit**, driftarm), und einer
**Vorwärts-Sicht** (die RTM-Tabelle nennt je Anforderung die Design-Artefaktmenge),
die als kuratierter Spiegel driftet. Mit einem opt-in `trace.cross-consistency`-
Block macht [`DC-FA-CLI-009`](#dc-fa-cli-009--requirements-traceability-matrix)
zusätzlich zum RTM-Lauf einen **Mengenabgleich**: Für jede Anforderungs-ID `R` sei
`F(R)` die Design-Artefaktmenge aus der Vorwärts-Sicht und `B(R)` die Menge der
Artefakte, deren Rück-Kante `R` nennt. Gemeldet werden je `R` die Differenzen
`F(R) \ B(R)` („in der RTM, ohne Rück-Kante") und `B(R) \ F(R)` („Rück-Kante ohne
RTM-Eintrag"), je mit Richtungslabel. `1:N` ist Normalfall; Modus `equal` (`F = B`,
beide Differenzen gaten) oder `superset` (`F ⊇ B`, nur `B \ F` gatet).

Beide Sichten sind **kuratierte Markdown-Tabellen**, gelesen über den vorhandenen
header-gebundenen Tabellen-Reader
([`DC-FA-REQ-001`](#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen))
plus die Abschnitts-/Span- und **range-aware** ID-Semantik von
[`DC-FA-COV-001`](#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
(Range-/Enum-Notation `<FAM>-AAA..BBB` und `/`-Aufzählung). Die einzige neue Logik
ist: eine Sicht invertieren und die Mengen diffen — **kein neuer Parser**.

- **Vorwärts:** header-gebundene Anforderungs- und Design-Spalte; Design-Artefakte
  per `design-pattern` (RE2) aus der Design-Zelle, Anforderungen range-aware per
  `forward.req-pattern` (RE2) aus der ID-Spalte — **symmetrisch** zum
  `backward.req-pattern` der Rück-Sicht.
- **Rückwärts:** je Tabelle mit einer `Bezug`-Spalte (header-gebunden) ist die
  **Artefakt-ID die erste Spalte** (positionell, da deren Header je Tabelle variiert
  — `Kennung`/`Port-ID`/`Tabu-ID`/`Komponente` —, per `design-pattern` extrahiert);
  die `Bezug`-Zelle liefert die Anforderungs-IDs (`req-pattern`, range-aware).
- **Namensraum-Vorbedingung:** Vorwärts- und Rückwärts-Artefakt-IDs werden mit
  **demselben** `design-pattern` erkannt und liegen damit im selben Namensraum;
  andernfalls wäre der Mengen-Diff inhärent leer/voll und bedeutungslos.
- **Die Vergleichs-Schlüsselmenge ist nicht die RTM-Anforderungsmenge.** Welche
  Anforderungen verglichen werden, entscheiden allein `forward.req-pattern` und
  `backward.req-pattern` — **nicht**, ob eine Anforderung in der RTM steht. Beide
  fallen per Default auf `requirements.id-pattern` zurück; wer die RTM bewusst auf
  eine Teilmenge scopt (etwa Architektur-Meta ausschließt), setzt sie **explizit**
  weiter. Die Kopplung ist damit eine sichtbare Konfigurationsentscheidung: still
  erzeugte sie eine leere Vorwärts-Sicht, deren Rück-Kanten wie echter Drift
  aussehen.

Bei mehrschichtigen Spec-Modellen (Vertrag → Spezifikation → Architektur) zeigen
Rück-Kanten teils auf **Mittelschicht-IDs** ohne eigene RTM-Vorwärts-Zeile. Diese
**Ableitungssprünge** nimmt `exclude-req` (RE2) aus dem Abgleich. Das Ventil ist
selbst eine kuratierte Kante, die mit der Schicht-Struktur synchron bleiben muss
und driften kann (wie `matrix.exclude-sections`) — bewusst benannter Trade-off.

Ein **vakuumer Abgleich** ist kein bestandener Abgleich: bleibt nach allen Stufen
keine Kante übrig, behauptete ein `0 Differenz(en)` eine nie geprüfte Konsistenz.
Maßgeblich ist die **Wirkung**, nicht die Ursache — geprüft wird, ob der Abgleich
konstruktionsbedingt überhaupt einen Befund liefern könnte. Vakuum heißt: **beide**
Sichten kantenleer — oder, unter `mode: superset`, die **Rück**-Sicht kantenleer
(dann kann `B \ F` nie einen Befund liefern). Beides ⇒ Exit 2. Gemessen wird
**nach** dem `exclude-req`-Ausschluss: ein Ventil, das alle Anforderungen
verschluckt, schaltet das Gate ebenso still ab wie ein Muster, das am Inhalt
vorbeigreift — es ist selbst eine kuratierte, drift-fähige Kante.
Eine **einseitig** leere Sicht ist dagegen ein wohldefiniertes Ergebnis, kein
Fehler: der Diff läuft über `keys(F) ∪ keys(B)`, und eine noch unrestrukturierte
Vorwärts-Sicht bei gepflegten Rück-Kanten ist der erwartete Bootstrap-Zustand —
sie meldet `B \ F` laut, statt den Lauf abzubrechen.

Der Abgleich ist **fail-closed** (ungültiges Regex, fehlende Spalte, ID-Header
nicht genau einmal, unbekannter `mode`, vakuumer Abgleich ⇒ Exit 2,
[`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)) und **advisory**: ohne Block ist die
RTM byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)), nichts wird geschrieben
([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)); der Exit-
Code ändert sich nur unter dem globalen
[`--require-complete`](#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
(**kein** block-lokaler Schalter). Befund-Ausgabe ist deterministisch sortiert.

**Akzeptanzkriterien:**

- **Happy Path (konsistentes 1:N):** Given eine Anforderung, die vorwärts fünf
  Design-Artefakte nennt, und Rück-Kanten derselben fünf Artefakte auf sie, when
  `d-check --trace` mit `trace.cross-consistency` (`mode: equal`) läuft, then keine
  Differenz, Exit 0; ein read-only gemountetes Repository genügt.
- **Boundary (Richtungs-Differenz):** Given `F(R) = {A, B}`, `B(R) = {B, C}`, when
  der Abgleich läuft, then zwei Befunde — `A` „in RTM, ohne Rück-Kante", `C`
  „Rück-Kante, ohne RTM-Eintrag" — je mit `Datei:Zeile`, Anforderungs-ID, Richtung.
- **Range-aware:** Given eine Rück-Kante `<FAM>-001..009` und eine Vorwärts-Zeile
  `<FAM>-001..004`, when der Abgleich läuft, then werden beide Bereiche expandiert
  und je Einzel-ID verglichen.
- **Superset-Modus:** Given `mode: superset` und `F(R) ⊋ B(R)`, then ist
  `F(R) \ B(R)` kein Befund, nur `B(R) \ F(R)` würde einen erzeugen.
- **Ableitungssprung ausgeschlossen:** Given eine Rück-Kante auf eine
  `exclude-req`-ID (Mittelschicht-Familie), then fällt sie aus dem Vergleich —
  weder Differenz noch Waise.
- **Negative (Config, fail-closed):** Given ungültiges Regex, unbekannter `mode`,
  fehlende header-gebundene Spalte oder ID-Header nicht genau einmal, when
  `d-check --trace` läuft, then Exit 2, kein Abgleich.
- **Negative (Vakuum, fail-closed):** Given ein `design-pattern`, das kompiliert,
  aber am Artefakt-Namensraum vorbeigreift (beide Sichten kantenleer) — oder eine
  kantenleere Rück-Sicht unter `mode: superset` —, when `d-check --trace` läuft,
  then Exit 2 statt `0 Differenz(en)`/Exit 0.
- **Negative (übergriffiges Ventil, fail-closed):** Given fehlerfreie Muster, aber
  ein `exclude-req`, das **jede** Anforderung beider Sichten verschluckt, when
  `d-check --trace` gegen einen realen Drift läuft, then Exit 2 — der Abgleich
  könnte nie einen Befund liefern; die Meldung benennt das **Ventil**, nicht die
  korrekten Muster.
- **Boundary (RTM-Scope ≠ Vergleichs-Scope):** Given ein `requirements.id-pattern`,
  das eine Familie bewusst ausschließt (Architektur-Meta), und ein
  `forward.req-pattern`, das sie einschließt, when der Abgleich läuft, then wird
  diese Familie **verglichen** — beide Richtungsdifferenzen erscheinen, obwohl die
  RTM sie nicht führt.
- **Boundary (einseitig leere Sicht):** Given eine kantenleere **Vorwärts**-Sicht
  und `B(R) ≠ ∅`, when der Abgleich läuft, then ist das **kein** Fehler: jede
  Rück-Kante wird als `B \ F` gemeldet (der Bootstrap-Zustand vor der
  Restrukturierung der Vorwärts-Sicht).
- **Default byte-identisch:** Given **kein** `trace.cross-consistency`-Block, when
  `d-check --trace [--require-complete]` läuft, then RTM und Exit byte-identisch zur
  Fassung vor dieser Anforderung ([`DC-QA-02`](#dc-qa-02--determinismus)), nichts
  geschrieben.

**Out-of-Scope:** (1) Der **Generator** (Vorwärts-RTM aus den Rück-Kanten erzeugen
statt hand-pflegen) — Schreib-Pfad, näher an
[`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch); Ziel-Zustand, aber eigene
spätere Anforderung/ADR (dieses Gate ist sein Korrektheits-Harness). (2) **Prosa**-
`Bezug:`-Zeilen ohne Tabelle (v1 liest nur `Bezug`-**Spalten**; Konsument
tabellarisiert oder akzeptiert Nicht-Gatung). (3) **Transitive Mehr-Hop-Auflösung**
über mehr als eine Zwischenschicht. (4) Kuratierte, bewusst designlose Anforderungen
(Abschnitts-Blacklist, nicht invertierbar).

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

- **Happy Path:** Given ein kaputter Link in `docs/a.md`, when `d-check` ohne Konfiguration läuft, then wird der Befund gemeldet. <!-- d-check:ignore (AK-Beispiel) -->
- **Boundary:** Given ein kaputter Link in `node_modules/x/README.md`, when `d-check` läuft, then wird kein Befund gemeldet.
- **Negative:** Given eine Konfiguration, die explizit die Scan-Wurzel `handbuch/` deklariert, die nicht existiert, when `d-check` läuft, then Exit-Code 2 mit Hinweis auf die fehlende Wurzel (Default-Wurzeln dagegen sind optional).

**Out-of-Scope:** Prüfung anderer Markup-Formate (reStructuredText, AsciiDoc, HTML).

---

### DC-FA-REF-001 — Geteiltes Referenz-Ventil (`ignore-refs` mit Quell-Skopus)

**Beschreibung:** Ein optionales, geteiltes Ventil nimmt **bestimmte
Referenz-Ziele** von der Existenz- und Anker-Prüfung der Module `links`
([`DC-FA-LINK-001`](#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)),
`anchors`
([`DC-FA-ANCH-001`](#dc-fa-anch-001--heading-anker-validierung-modul-anchors)) und
`codepaths`
([`DC-FA-CODE-001`](#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in))
aus — **referenz-weit**, unabhängig davon, in welcher Datei/Zeile die Referenz
steht. Es ist das **Ziel-Achsen-Pendant** zu `scan.ignore`
([`DC-FA-SCAN-001`](#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln),
Quell-Achse): jenes entfernt ganze Dateien vom Scan, dieses einzelne Ziele von der
Prüfung. Die Deklaration steht als Top-Level-Liste `ignore-refs`; jeder Eintrag
trägt drei Felder:

- **`refs`** (Glob-Liste, Pflicht): Muster auf den **aufgelösten Ziel-Pfad**, so
  wie das jeweilige Modul heute auflöst (der Befund-`target` ist der aufgelöste
  Pfad, gleiche Ventil-Parität wie bei den `exempt-*`-Ventilen). Ein leeres oder
  fehlendes `refs` macht den Eintrag inert (nichts wird ignoriert).
- **`keep`** (Glob-Liste, optional): Ausnahmen. **Ein Ziel wird genau dann
  ignoriert, wenn `refs` matcht UND `keep` nicht** (`refs ∧ ¬keep`). `keep` gewinnt
  **unbedingt und reihenfolge-unabhängig** — bewusst **nicht** gitignore-Last-Match.
- **`in`** (Glob, optional): **Quell-Skopus** — Muster auf den Pfad der
  **Quelldatei** (der Datei, in der die Referenz vorkommt). Gesetzt, gilt der
  Eintrag **nur** für Referenzen in passenden Dateien; ohne `in` gilt er repo-weit.

Das Ventil ist ein **Register bewusst behandelter Referenzen** — etwa
Template-Platzhalter, die auf Ziel-Repo-Positionen zeigen (im Quell-Repo nicht
auflösbar), oder historische/immutable Doku, die einen entfernten Pfad zitiert.
Bewusster Akt **mit Gate**: ohne passenden Eintrag meldet ein fehlendes Ziel weiter
(nichts dangelt still); das Ventil unterdrückt **nur** die Auflösungs-Klasse
des genannten Ziels — Existenz-, Anker- und Ortsfestigkeits-Prüfung
(`target-missing`/`anchor-missing`/`codepath-missing` und
`link-position-dependent`) —, keine anderen Befunde (Symlink-Ablehnung,
Repo-Escape u. Ä. bleiben). Zwei Felder statt `!`-Negations-Syntax: gemessen braucht kein realer Fall
die Re-Ignore-Alternierung, die eine Ordnungssemantik böte — der Verzicht kauft
Reihenfolge-Unabhängigkeit in einer YAML-Liste.

**Alias (kein Config-Bruch):** die bestehende modul-lokale Liste
`codepaths.ignore-refs`
([`DC-FA-CODE-001`](#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in))
bleibt als **Alias** gültig — sie wirkt wie ein `ignore-refs`-Eintrag ohne
`in`/`keep`, skopiert auf `codepaths`-Ziele; bestehende Konfigurationen laufen
unverändert.

**Achsen-Abgrenzung:** `ignore-refs` (Ziel) ist orthogonal zu `scan.ignore`
(Quelldatei, ganze Datei) und zum Zeilen-Marker `d-check:ignore` (Zeile;
honoriert von `codepaths`, [`ids`](#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids),
[`versions`](#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in) und
[`diagrams`](#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in)
— eine benannte Liste, kein ableitbares Kriterium; `matrix`, `structure` und
[`citations`](#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in)
melden ebenfalls auf Zeilen und tragen ihn nicht). Keine Achse überschreibt eine andere; ein Ziel, das `scan.ignore`
bereits vom Scan entfernt, erreicht `ignore-refs` gar nicht.

**Wo der Zeilen-Marker gilt — und wo er nur erwähnt ist.** *„Ist diese Zeile
eine Direktive"* ist eine **Prosa-Frage** und wird geteilt beantwortet: der
Marker wirkt nur **außerhalb** von Fenced-Blöcken **und außerhalb von
Inline-Code**. In Backticks ist er eine **Erwähnung** — sonst nähme eine Zeile,
die das Ventil beschreibt, sich selbst aus der Prüfung. Das gilt für die zwei
Konsumenten mit Prosa-Eingabe, `codepaths` und
[`ids`](#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids). Bei
[`versions`](#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in) und
[`diagrams`](#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in)
wird er weiterhin **roh** erkannt, und das ist eine **benannte Skopierung**:
jene lesen alle Zeilen einschließlich Fences bzw. die Zeilen **innerhalb** eines
Fence samt Öffnungszeile, wo ein Backtick literaler Inhalt ist. Bei `diagrams`
ist das strukturell; bei `versions` ist es eine **benannte Grenze** — seine
Eingabe **enthält** Prosa-Zeilen, und auf einer solchen antwortet das Produkt
damit zweifach.

**Auch die FORM folgt der Eingabe.** Sie ist je Konsument verschieden, weil die
Kommentar-Lexik es ist: bei `codepaths` und
[`ids`](#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids) — Eingabe ist
Markdown-Prosa — muss der Marker in einem **HTML-Kommentar** stehen; eine blanke
Erwähnung wirkt nicht. Bei
[`diagrams`](#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in)
ist er ein **Token**, weil in einem Fence die Diagramm-Sprache den Kommentar
bildet und eine Lexik je Fence-Sprache ein Grammatik-Parser wäre; bei
[`versions`](#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in)
ebenso, weil es alle Zeilen sprachgemischt liest — bei `diagrams` ist das strukturell, bei `versions` dieselbe **benannte Grenze** wie auf der Lage-Achse. **Preis, ausgewiesen:** wer
den Marker setzt, muss wissen, für welches Modul.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine Referenz auf ein nicht existierendes Ziel in einer Datei, die ein `in`-Glob matcht, und das aufgelöste Ziel matcht `refs`, aber nicht `keep`, when ein Ventil-tragendes Modul (`links`/`anchors`/`codepaths`) läuft, then kein Befund.
- **Boundary (`keep` gewinnt):** Given denselben `refs`-Match, aber das aufgelöste Ziel matcht zusätzlich ein `keep`-Glob (ein real auflösender Verweis im selben Baum), when das Modul läuft, then wird die Referenz **geprüft** — `keep` hebt die Ausnahme reihenfolge-unabhängig auf (ist das Ziel real, kein Befund; ist es verfälscht, ein Befund).
- **Negative (Tippfehler):** Given ein verfälschter Pfad in einer `in`-passenden Datei, dessen aufgelöstes Ziel **nicht** `refs` matcht (oder von `keep` zurückgeholt wird), when das Modul läuft, then ein Befund mit Datei, Zeile, Ziel und Grund, Exit-Code 1 — eine „ignoriere, was nicht auflöst"-Heuristik bestünde diesen Test nicht, deshalb Muster statt Heuristik.
- **Ventil (Skopus-Isolation):** Given ein `ignore-refs`-Eintrag mit `in`-Glob, when eine Referenz auf **dasselbe** Ziel-Muster in einer Datei **außerhalb** des `in`-Globs vorkommt, then bleibt sie voll geprüft.
- **Ventil (Alias):** Given eine bestehende `codepaths.ignore-refs`-Liste ohne Top-Level-`ignore-refs`, when `codepaths` läuft, then verhält es sich byte-identisch zur Fassung vor dieser Anforderung ([`DC-QA-02`](#dc-qa-02--determinismus)).
- **Regression (default-aus):** Given weder Top-Level-`ignore-refs` noch `codepaths.ignore-refs`, when die Module laufen, then ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)).

**Out-of-Scope:** Negations-Syntax (`!`-Globs) als zweiter Glob-Dialekt;
gitignore-Last-Match-Ordnung (`keep` ist reihenfolge-unabhängig);
Re-Ignore-Alternierung (`ignore`→`keep`→`ignore` — zwei Felder können das bewusst
nicht, gemessen kein Bedarf); Zeilen-Granularität (dafür der `d-check:ignore`-Marker
in `codepaths`); semantische Prüfung, ob das ignorierte Ziel inhaltlich passt.

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
(``` / ~~~) und in Inline-Code werden **von diesem Modul** nicht
geprüft (explizite Pfade in Inline-Code prüft das opt-in-Modul
`codepaths`,
[`DC-FA-CODE-001`](#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in));
HTML-Kommentare werden nicht gesondert behandelt (Links darin gelten
als Fließtext). Das geteilte Referenz-Ventil `ignore-refs`
([`DC-FA-REF-001`](#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus))
nimmt bestimmte Ziel-Pfade referenz-weit von der Existenz-Prüfung aus.

**Ortsfeste Verweise (opt-in, `links.resolve-from`).** Wo Dateien per `git mv`
zwischen Geschwister-Verzeichnissen **wandern** (der Planning-Lifecycle ist der
Anlassfall), muss ein relativer Verweis von **jedem** Ort der Gruppe auflösen —
nicht nur vom Ist-Ort: der Lifecycle-Wechsel ist eine Bewegung ohne
Inhaltsänderung, und ein präfixloser Nachbar-Verweis bricht mit ihr. Eine
Gruppe benennt ihre **wandernden** Verzeichnisse (`dirs`, deren unmittelbare
Dateien Quellen sind — Unterverzeichnisse wandern nicht als Einheit mit) und
optional **ortsfeste** (`fixed-dirs`, etwa den Ruheort: hypothetische Ziele der
Wanderung, deren eigene Dateien am Endzustand sind — sonst wären am gemessenen
Bestand 108 Befunde Falsch-Positive). **Positionsabhängig** ist ein Verweis,
der von mindestens einem Ort der Gruppe **nicht** auflöst **oder** von
verschiedenen Orten auf **verschiedene** Ziele; beide Fälle melden den eigenen
Grund-Code `link-position-dependent` (nicht `target-missing`: am Ist-Ort ist
nichts kaputt, die Reparatur ist das Präfixieren des Pfads). **Vorbedingung ist
der saubere Ist-Ort** — löst das Ziel schon dort nicht auf, melden die
bestehenden Prüfungen, und diese schweigt (kein Doppelbefund). **Zwei Grenzen
sind benannt:** die **Ziel**-Wanderung (Quelle ortsfest, Ziel wandert — etwa
der Review-Report auf einen `in-progress/`-Slice) prüft diese Fähigkeit nicht,
sie prüft hypothetische Quell-Orte; und ein **einzelner** fehlender Gruppen-Ort
ist von einem legitim geleerten Verzeichnis nicht unterscheidbar (git überträgt
leere Verzeichnisse nicht) — fail-closed meldet erst, wenn **kein** `dirs`-Ort
der Gruppe existiert oder ein Ort als **Datei** existiert. Und die Gruppen-Orte
müssen im wirksamen **Scan-Bereich** liegen: eine nie gescannte Datei ist still
keine Quelle. Anker bleiben außen
vor; Ziel-Menge, Vorverarbeitung, Dekodierung und Ventile sind die der
bestehenden Prüfung. Ohne den Block ist der Befundsatz byte-identisch
([`DC-QA-02`](#dc-qa-02--determinismus)). Begründung in begleitender ADR.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Link auf eine existierende Datei im Repo, when das Modul `links` läuft, then kein Befund.
- **Boundary:** Given ein Link auf eine nicht existierende Datei innerhalb eines Fenced-Code-Blocks, when das Modul läuft, then kein Befund.
- **Negative:** Given ein Link `../../etc/passwd`, dessen Ziel die Repository-Wurzel verlässt, when das Modul läuft, then ein Befund mit Grund „verlässt Repository" — auch wenn der Zielpfad existiert. <!-- d-check:ignore (AK-Beispiel: Angriffs-Pfad) -->
- **Negative:** Given ein Link auf eine nicht existierende Datei, when das Modul läuft, then ein Befund mit Datei, Zeile, Ziel und Grund, Exit-Code 1.

- **Happy Path (resolve-from):** Given eine Gruppe aus drei wandernden Verzeichnissen, eine Quell-Datei in einem davon und einen Verweis der Form von einem Geschwister-Präfix (`../<ort>/<datei>`) auf ein existierendes Ziel, when das Modul läuft, then kein Befund — der Verweis löst von jedem Ort der Gruppe auf dasselbe Ziel auf.
- **Negative (resolve-from, Auflösbarkeit):** Given denselben Aufbau, aber einen **präfixlosen** Nachbar-Verweis (`<datei>` im selben Verzeichnis), when das Modul läuft, then ein Befund `link-position-dependent`, dessen Meldung einen nicht auflösenden Ort nennt — am Ist-Ort existiert das Ziel.
- **Negative (resolve-from, Ziel-Identität):** Given einen Verweis, der von jedem Ort der Gruppe auflöst, aber auf **verschiedene** existierende Dateien, when das Modul läuft, then ein Befund `link-position-dependent`, dessen Meldung die divergierenden Ziele nennt.
- **Boundary (ortsfeste Datei):** Given eine Datei in einem `fixed-dirs`-Verzeichnis mit einem Verweis, der nur vom Ist-Ort auflöst, when das Modul läuft, then **kein** Befund — ortsfeste Dateien sind keine Quellen.
- **Modul-aus (resolve-from):** Given **kein** `links.resolve-from`-Block, when `d-check` läuft, then ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)).

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
innerhalb derselben Datei) werden gegen die **gültige Anker-Menge** der
Zieldatei geprüft. Zu ihr zählen (1) die Heading-Slugs nach dem
GitHub-Slug-Verfahren (Kleinschreibung, Sonderzeichen-Entfernung,
Leerzeichen → `-`, Duplikat-Suffixe `-1`, `-2`, …; Headings aller Ebenen
`#`–`######`, das erste gleichlautende Vorkommen in Dokumentreihenfolge
trägt den Basis-Slug, weitere die Suffixe) und (2) die **Inline-HTML-Anker**
der Datei: der Wert eines `id`-Attributs an einem beliebigen HTML-Element
sowie der Wert eines `name`-Attributs an einem `<a>`-Element
(GitHub-Render-Verhalten). HTML-Anker werden **wörtlich** verglichen (kein
Slug, keine Kleinschreibung). Ein Fragment gilt als aufgelöst, wenn es
einen Heading-Slug **oder** einen Inline-HTML-Anker der Zieldatei trifft.
Existiert die Zieldatei nicht, wird die Anker-Prüfung für diesen Link
übersprungen — der Befund kommt vom Modul `links`
([`DC-FA-LINK-001`](#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)).
Das geteilte Referenz-Ventil `ignore-refs`
([`DC-FA-REF-001`](#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus))
nimmt bestimmte Ziele referenz-weit von der Anker-Prüfung aus.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Link `a.md#zweck-und-geltungsbereich` und ein Heading „Zweck und Geltungsbereich" in `a.md`, when das Modul läuft, then kein Befund.
- **Boundary:** Given zwei gleichlautende Headings „Beispiel" in der Zieldatei, when ein Link auf `#beispiel-1` zeigt, then kein Befund (Duplikat-Suffix korrekt aufgelöst).
- **Boundary:** Given ein Link `fehlt.md#x` auf eine nicht existierende Datei, when die Module `links` und `anchors` aktiv sind, then genau ein Befund — aus `links`; das Modul `anchors` schweigt.
- **Negative:** Given ein Link auf `a.md#gibt-es-nicht`, when das Modul läuft, then ein Befund mit Grund „Anker nicht gefunden".
- **Happy Path (HTML-Anker):** Given ein Link `a.md#abschnitt-2` und in `a.md` der Inline-HTML-Anker `<a name="abschnitt-2"></a>` (ohne gleichnamigen Heading), when das Modul läuft, then kein Befund.
- **Boundary (HTML-Anker):** Given in der Zieldatei `<div id="Übersicht">` und ein Link auf `#Übersicht`, when das Modul läuft, then kein Befund (HTML-Anker wörtlich/case-sensitiv, nicht geslugged).
- **Negative (HTML-Anker):** Given `id="x"` ausschließlich innerhalb einer Code-Auszeichnung (Fenced-Code-Block oder Inline-Code-Span, nicht gerendert) und ein Link auf `#x` ohne passenden Heading-Slug oder HTML-Anker außerhalb von Code-Auszeichnung, when das Modul läuft, then ein Befund mit Grund „Anker nicht gefunden".

**Out-of-Scope:**

- HTML-Anker in über mehrere Zeilen verteilten Tags (Erkennung ist konservativ und zeilenbasiert).
- `name`-Attribut an anderen Elementen als `<a>` (GitHub honoriert `name` nur dort).
- Prüfung von HTML als eigenständigem Dateiformat — bleibt global Out-of-Scope (§5); erkannt werden ausschließlich Inline-HTML-Anker **innerhalb von Markdown-Dateien**.

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

**Link-Politik (`link-policy`, je Muster konfigurierbar).** Pro Muster
ist wählbar, wie streng die Linkpflicht greift:

- `prose` (**Default**): wie oben — nur nackte Vorkommen im Fließtext
  sind linkpflichtig; Inline-Code-Vorkommen sind frei.
- `always`: auch Vorkommen **innerhalb von Inline-Code** müssen im
  Linktext eines Markdown-Links stehen (`` [`ADR-0042`](ziel) `` ist
  erfüllt, `` `ADR-0042` `` allein nicht). Fenced-Code-Blöcke, <!-- d-check:ignore (Beispiel-ID, fiktiv) -->
  Heading-Zeilen, das `target` des Musters sowie die beiden unten
  genannten Ventile (`exempt-paths`, `d-check:ignore`) bleiben frei.

Der Default ist `prose`, damit bestehende Konfigurationen
byte-identisch bleiben und kein Repo ungefragt rote Läufe bekommt
(opt-in fürs Gating). Die Entdeckung ungenügender Verlinkung ist davon
unabhängig eine Bringschuld des Werkzeug-Betreibers (fleet-weiter
`always`-Lauf, [`DC-QA-04`](#dc-qa-04--migrationsabdeckung-der-alt-tools)-Muster);
„opt-in" heißt nicht „unsichtbar". Zwei Ventile nehmen eine Datei bzw.
eine Zeile von der Linkpflicht eines Musters aus. Beide gelten für
**alle** Vorkommen des Musters — nackt im Fließtext **wie** in
Inline-Code — und unabhängig von der `link-policy`; sie sind ein
Ganzdatei- bzw. Ganzzeilen-Carve-out, kein bloßes Dämpfen der
`always`-Strenge:

- `exempt-paths` (Glob-Liste je Muster): Dateien, in denen das Muster
  **keine** Linkpflicht erzeugt (literal-schwere Artefakte wie
  Changelogs oder Review-Reports) — gleich, ob die Kennung dort nackt
  oder in Backticks steht.
- Der Zeilen-Marker `d-check:ignore` (HTML-Kommentar, Begründung in
  Klammern empfohlen) nimmt die **ganze Zeile** von der `ids`-Prüfung
  aus — für bewusst illustrative Beispiel-IDs (gleiche Begründung wie
  bei [`DC-FA-CODE-001`](#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in);
  der Marker wirkt ab dieser Anforderung auf `codepaths` **und** `ids`),
  ebenfalls für nacktes wie für Inline-Code-Vorkommen.

**Akzeptanzkriterien:**

- **Happy Path:** Given das Muster `ADR-\d{4}` und ein Vorkommen `[ADR-0042](docs/plan/adr/0042-beispiel.md)`, when das Modul `ids` läuft, then kein Befund.
- **Boundary:** Given ein Vorkommen `` `ADR-0042` `` in Inline-Code und `link-policy: prose` (Default), when das Modul läuft, then kein Befund (Code-Vorkommen sind linkpflichtfrei). <!-- d-check:ignore (Beispiel-ID, fiktiv) -->
- **Negative:** Given ein nacktes `ADR-0042` im Fließtext (außerhalb `exempt-paths`, ohne `d-check:ignore`), when das Modul läuft, then ein Befund mit Grund „Kennung ohne Link". <!-- d-check:ignore (Beispiel-ID, fiktiv) -->
- **`always` Happy Path:** Given `link-policy: always` und ein Vorkommen `` [`ADR-0042`](docs/plan/adr/0042-beispiel.md) ``, when das Modul läuft, then kein Befund (Code-Span im Linktext zählt als verlinkt).
- **`always` Negative:** Given `link-policy: always` und ein Vorkommen `` `ADR-0042` `` ohne Link (außerhalb `exempt-paths` und ohne `d-check:ignore`), when das Modul läuft, then ein Befund `id-unlinked`. <!-- d-check:ignore (die Kennungen dieser Zeile sind Beispiele des Kriteriums, keine Verweise) -->
- **`always` Boundary (Ventile):** Given `link-policy: always`, ein `` `ADR-0042` `` in einer `exempt-paths`-Datei und ein zweites `` `ADR-0099` `` auf einer Zeile mit `d-check:ignore`, when das Modul läuft, then kein Befund für beide. <!-- d-check:ignore (dito) -->
- **Boundary (Marker in Inline-Code ist Erwähnung):** Given eine Zeile, die den Zeilen-Marker in **Backticks** nennt und daneben einen nicht auflösbaren Pfad bzw. eine unverlinkte Kennung trägt, when `codepaths` bzw. `ids` läuft, then **ein Befund** — dieselbe Zeile mit **freiem** Marker meldet keinen (Gegenprobe). Given dieselbe Konstellation bei `versions` oder `diagrams`, then **kein** Befund: dort wird der Marker weiterhin roh erkannt, weil ihre Eingabe kein Inline-Code kennt.
- **Boundary (Form folgt der Eingabe):** Given eine Zeile mit **blankem** Marker (ohne Kommentar-Klammer) und einem nicht auflösbaren Pfad bzw. einer unverlinkten Kennung, when `codepaths` bzw. `ids` läuft, then **ein Befund** — dieselbe Zeile mit `<!-- d-check:ignore -->` meldet keinen (Gegenprobe). Given dieselbe Zeile bei `versions`, then **kein** Befund: dort ist der Marker ein **Token**. Given eine **Diagramm-Fence-Zeile** (oder die Öffnungszeile) mit blankem Token und einer undefinierten Kennung, when `diagrams` läuft, then **kein** Befund — eine Prosa-Zeile ist dort nie Eingabe, das Kriterium braucht deshalb seine eigene Konstellation. **Konservativ:** ein `>` im Kommentar vor dem Marker lässt ihn nicht gelten.
- **Ventile für nackte Vorkommen:** Given ein **nacktes** `ADR-0042` im Fließtext einer `exempt-paths`-Datei und ein zweites nacktes `ADR-0099` auf einer Zeile mit `d-check:ignore`, when das Modul läuft, then kein Befund für beide — die Ventile gelten für alle Vorkommen des Musters (nackt wie Inline-Code), unabhängig von der `link-policy`. <!-- d-check:ignore (Beispiel-IDs, fiktiv) -->

**Out-of-Scope:** Automatisches Ermitteln der Muster aus dem Repo-Inhalt **für die Prüfung** (die Prüfung läuft stets gegen explizit konfigurierte Muster — deterministisch, Vertrag); ein *advisory* Scaffold-Modus darf Muster aus **benannten Autoritäts-Quellen** ableiten (Ausgabe-only, vom Menschen bestätigt — [`DC-FA-CLI-006`](#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)). Prüfung, ob die verlinkte Definition inhaltlich zur Kennung passt; ein `link-policy`-Default abweichend von `prose` (bewusst opt-in).

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

**Supersede-Lineage-Ausnahme (opt-in).** Ein ablösendes Dokument
verweist per Definition auf das Dokument, das es ablöst (z. B.
ADR→ADR-Lineage); dieser Verweis zielt notwendig auf ein inaktives
Dokument und ist dennoch legitim. Ist `allow-supersede-lineage`
aktiviert, wird genau die Kante X → Y von der Status-Prüfung
ausgenommen, wenn die Quelle X über ein Feld aus `supersede-fields`
(z. B. `Supersedes`, `Aenderungstyp`) deklariert, dass sie Y ablöst —
der Feldwert nennt Y über den Linktext oder den Zielpfad der Referenz.
Alle anderen Referenzen auf Y bleiben `matrix-inactive`, und die
Klassen-Regelprüfung (`matrix-forbidden`) ist davon unberührt. Default:
aus — ohne gesetztes Flag ist der Befundsatz byte-identisch
([`DC-QA-02`](#dc-qa-02--determinismus)). Wie alle Module außer
`codepaths`/`ids` trägt `matrix` **keinen** Zeilen-Opt-out-Marker
(`d-check:ignore`): deterministische Befunde werden behoben oder
strukturell ausgenommen (`exclude-sections`, `allow-supersede-lineage`),
nicht zeilenweise stummgeschaltet.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Slice mit Link auf ein ADR mit Status `Accepted`, when das Modul `matrix` läuft, then kein Befund.
- **Boundary:** Given ein ADR mit Status `Superseded by ADR-0099`, when ein Slice darauf verlinkt, then ein Befund mit Grund „Referenz auf inaktives ADR". <!-- d-check:ignore (Beispiel-ID, fiktiv) -->
- **Negative:** Given ein Link aus `spec/lastenheft.md` auf eine ADR-Datei, when das Modul läuft, then ein Befund mit Grund „verbotene Abwärtsreferenz" und Angabe beider Dokumentklassen.
- **Lineage Happy:** Given `allow-supersede-lineage: true`, `supersede-fields: [Aenderungstyp]` und ein ADR X, das ein Feld `Aenderungstyp: Supersedes ADR-0099` sowie einen Link `[ADR-0099](0099-x.md)` auf das mit `superseded` markierte ADR Y trägt, when das Modul läuft, then **kein** `matrix-inactive`-Befund für diese Kante. <!-- d-check:ignore (Beispiel-IDs, fiktiv) -->
- **Lineage Boundary:** Given dieselbe Konfiguration und eine **andere** Datei ohne Supersede-Deklaration auf Y, die auf Y verlinkt, when das Modul läuft, then bleibt für diese Datei der `matrix-inactive`-Befund — der Carve-out gilt nur für die deklarierte Lineage-Kante.
- **Token-Ziel Happy:** Given eine Klasse ohne `paths` mit einem `token`-Muster und eine Regel, die eine Dokumentklasse als Quelle und diese Klasse als Ziel nennt, when im Körper eines Quell-Dokuments ein Treffer des Musters steht, then ein `matrix-forbidden`-Befund in der Token-Form.
- **Token-Ziel Boundary:** Given dieselbe Klasse als **Quelle** einer Regel, when das Modul läuft, then kein Befund aus dieser Regel — die Klasse hat keine Mitglieder, also nie ein Quell-Dokument.
- **Lineage Negative (Default):** Given keine `allow-supersede-lineage`-Angabe (Default aus), when das Modul über dieselben Dateien läuft, then erzeugt auch die Lineage-Kante `matrix-inactive` — der Befundsatz ist bit-genau wie ohne das Feature.

**Out-of-Scope:** Semantische Unterscheidung von Verweis-Zwecken (z. B. Verifikations-Zeiger vs. Entscheidungsgrundlage) — das bleibt Review-Aufgabe; Provenance-/Historie-Sektionen können per Konfiguration von der Prüfung ausgenommen werden, eine automatische Erkennung solcher Sektionen ist nicht gefordert; die Supersede-Lineage-Ausnahme erkennt Ablösungen ausschließlich über deklarierte `supersede-fields` (keine semantische Schlussfolgerung) und gilt nur für die Status-Prüfung (`matrix-inactive`), nicht für die Klassen-Regeln; ein zeilenweiser Opt-out-Marker für `matrix` ist nicht vorgesehen.

---

### DC-FA-MTX-002 — Verweisrichtung innerhalb einer geordneten Dokumentklasse (Modul `matrix`)

**Beschreibung:** Eine Dokumentklasse
([`DC-FA-MTX-001`](#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix))
kann zusätzlich eine **geordnete Rangfolge** und eine **Richtungspolitik**
tragen: `order` ist eine Liste von Pfad-Globs (autoritativste Schicht zuerst),
`direction` die Politik. Bei `direction: no-downward` ist jede
**klasseninterne** Referenz von einer höherrangigen (früher gelisteten,
autoritativeren) Datei auf eine niederrangige verboten; der Rang einer Datei
ist der Index des **ersten** `order`-Globs, den sie matcht (First-Match wie die
Klassenzuordnung). Damit wird die Source-Precedence-Schichtung **innerhalb**
eines Stratums maschinell prüfbar — kein Dokument verweist abwärts auf eine
weniger autoritative Schicht —, additiv zu den Klassen-Paar-Regeln aus
[`DC-FA-MTX-001`](#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix).
`order` listet **Globs**, nicht zwingend Einzeldateien: eine Schicht kann viele
Dateien fassen (eine Spec-Familie aus vielen Dateien braucht kein vollständiges
Auflisten). Klassen-Mitglieder ohne `order`-Treffer sind **rangfrei** und nehmen
an der Richtungsprüfung nicht teil. Fehlkonfiguration ist **fail-closed**:
unbekannter `direction`-Wert, `direction` ohne `order` und `order` ohne
`direction` sind Konfigurationsfehler — eine Richtungs-Deklaration darf nicht
still wirkungslos sein. Default (beide Felder leer): das Modul verhält sich
byte-identisch zu [`DC-FA-MTX-001`](#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
([`DC-QA-02`](#dc-qa-02--determinismus)).

**Akzeptanzkriterien:**

- **Happy Path:** Given eine Klasse mit `order: [a.md, b.md]` und `direction: no-downward` und ein Link in `b.md` auf `a.md` (aufwärts), when das Modul `matrix` läuft, then kein Befund.
- **Boundary:** Given dieselbe Klasse und ein Link in `a.md` auf `b.md` (abwärts — auch über mehr als eine Rang-Stufe), when das Modul läuft, then ein Befund mit Grund `matrix-downward` und Nennung beider Ränge.
- **Negative:** Given `order` ohne `direction` (oder `direction` ohne `order`, oder ein unbekannter `direction`-Wert), when die Konfiguration geladen wird, then ein Konfigurationsfehler (Exit 2) statt eines stillen No-op.

**Out-of-Scope:** Richtungsregeln **zwischen** verschiedenen Klassen (das sind die `rules`-Paare aus [`DC-FA-MTX-001`](#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)); automatische Herleitung der Rangfolge aus Dateiinhalt oder Source-Precedence-Tabelle (die Ordnung wird explizit deklariert); andere Politiken als `no-downward` (z. B. `no-upward`, `strict-adjacent`) bleiben einer späteren CR; rangfreie Mitglieder lösen keinen Befund aus.

---

### DC-FA-MTX-003 — Token-basierte Referenz-Richtung mit Provenance-Marker (Modul `matrix`)

**Beschreibung:** `matrix` erkennt verbotene Referenzen
([`DC-FA-MTX-001`](#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix))
bisher nur als **Links**; eine Referenz kann aber auch als **bare ID-Token** im
Fließtext stehen (etwa eine Slice-Kennung in einem ADR-Körper), die der
link-basierte Scan nicht sieht. Eine Klasse kann daher zusätzlich ein
**`token`-Muster** (Regex) tragen, das Referenzen auf diese Klasse als Token
erkennt. Sie kann ein `token`-Muster auch **ausschließlich** tragen, ohne
Pfad-Muster — ein **Token-Ziel**: die Klasse hat dann keine Mitglieder, ihr
Gegenstand ist eine Zeichenkette statt eines Dokuments (etwa ein Commit-Hash
oder ein Code-Modul-Pfad), und eine Regel, die sie als **Quelle** nennt, kann
nie feuern. Ohne Pfad-Muster **und** ohne `token`-Muster ist eine Klasse
inert. Tritt im Körper eines Dokuments der Klasse A ein `token` der Klasse B
auf — außerhalb Fenced-Code, außerhalb `exclude-sections` und außerhalb von
Markdown-Links (die deckt der Link-Scan ab) —, gilt das als Referenz A → B;
eine verbotene Kante erzeugt `matrix-forbidden` in der **Token-Form**.

**Provenance-Marker.** Eine verbotene Token-Referenz wird **ausgenommen**, wenn
auf derselben Zeile der Marker `<!-- d-check:status-provenance -->` steht — die
explizite Autor-Deklaration „dies ist Provenance/Verifikations-Zeiger, keine
Entscheidungsgrundlage". Das ist `matrix`' **erster** Zeilen-Marker und kehrt die
bisherige „nur strukturelle Ausnahmen"-Haltung von
[`DC-FA-MTX-001`](#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
bewusst um — als **benannte, semantische** Ausnahme (näher an
`allow-supersede-lineage` als an einem generischen Stummschalt-Marker); der
begleitende Architektur-Entscheid hält die Begründung fest. Die Deklaration ist
Selbst-Auskunft → ihre **Ehrlichkeit** (echte Provenance vs. getarnte
Entscheidungsgrundlage) bleibt Reviewer-Sache. **Links** bleiben markerlos:
legitime Provenance ist ein stabiler Token, kein Datei-Link, und Provenance
unter `## Geschichte`/Historie deckt bereits `exclude-sections` ab.

**Grandfathering.** Eine Datei, die ein `matrix.exempt-paths`-Glob matcht, wird
von `matrix` ganz übersprungen. Zweck: bereits `Accepted`-ADRs sind immutabel
und können nicht nachträglich markiert werden; sie werden grandfathered, neue
Dokumente ab Einführung tragen die Deklaration. Ohne `token`/`exempt-paths` ist
der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)).

**Akzeptanzkriterien:**

- **Happy Path:** Given eine Klasse `slice` mit `token: 'slice-\d{3}'`, eine verbotene Regel `{from: adr, to: slice}` und ein ADR-Körper, der eine Slice-Kennung **mit** `<!-- d-check:status-provenance -->` auf derselben Zeile nennt, when das Modul `matrix` läuft, then kein Befund.
- **Boundary:** Given dieselbe Konfiguration und derselbe ADR-Körper **ohne** Marker, when das Modul läuft, then ein `matrix-forbidden`-Befund (Token-Form) mit Nennung beider Klassen.
- **Negative:** Given eine Datei, die ein `matrix.exempt-paths`-Glob matcht und einen unmarkierten verbotenen Token trägt, when das Modul läuft, then kein Befund (grandfathered).

**Out-of-Scope:** Marker für **Link**-Referenzen (Links bleiben strukturell, `matrix-forbidden` ohne Ausnahme); die semantische Beurteilung „Provenance vs. Entscheidungsgrundlage" (sie steht hinter dem Marker und bleibt Reviewer-Sache, kein Linter-Urteil); Token-Erkennung für ohnehin link-pflichtige Kennungen (`ADR-`/`MR-`/`DC-` sind über `ids` Links — die Token-Erkennung zielt auf bare Kennungen ohne Linkpflicht).

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

### DC-FA-CODE-001 — Explizite Pfade in Inline-Code (Modul `codepaths`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `codepaths` werden
Inline-Code-Spans geprüft, deren Inhalt ein expliziter relativer Pfad
ist: beginnend mit `./` oder `../` (immer) oder mit einem der
konfigurierten Wurzel-Präfixe (z. B. `docs/`, `tools/` — relativ zur
Repository-Wurzel). Vor der Prüfung wird der Wert normalisiert:
umschließende Anführungszeichen und schließende Satzzeichen entfallen
(`` `../foo.md,` `` prüft `../foo.md`). Das Ziel muss nach <!-- d-check:ignore (Beispiel im Anforderungstext) -->
Pfadauflösung existieren und innerhalb der Repository-Wurzel liegen;
Fragment-Teile (`#…`) werden abgetrennt, und bei Markdown-Zielen wird
der Anker zusätzlich gegen die gültige Anker-Menge der Zieldatei geprüft
(Heading-Slugs und Inline-HTML-Anker) — gleiches Verfahren wie
[`DC-FA-ANCH-001`](#dc-fa-anch-001--heading-anker-validierung-modul-anchors).
Werte mit Whitespace,
Platzhalter-/Glob-Zeichen oder Ellipsen gelten nicht als Pfad
(konservative Erkennung); Vorkommen in Fenced-Code-Blöcken werden
nicht geprüft. Der Marker `d-check:ignore` auf derselben Zeile nimmt sie von
**genau dieser** Prüfung aus. Er gilt nur **in einem HTML-Kommentar** und nur
**außerhalb** von Fences und Inline-Code — blank oder in Backticks ist er eine
**Erwähnung** (§Achsen-Abgrenzung); die Begründung in der Klammer ist empfohlen.
Bewusst nicht existierende Beispiel-Pfade
(etwa Angriffs-Beispiele in Lehrtexten) sind kein Fehler, sondern
eine zu dokumentierende Absicht. Für alle anderen Module existiert
kein Opt-out-Marker: deterministische Befunde werden behoben, nicht
unterdrückt. Zusätzlich nimmt eine optionale `exempt-paths`-Glob-Liste
(Syntax wie `scan.ignore`, relativ zur Repository-Wurzel) **ganze Dateien**
von der `codepaths`-Prüfung aus — datei-weit, unabhängig von den
Wurzel-Präfixen; dasselbe Datei-Ventil wie bei
[`DC-FA-ID-001`](#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids),
komplementär zum zeilenweisen `d-check:ignore`-Marker (typischer Fall:
Review-Reports, die naturgemäß `Datei:Zeile`/Pfade zitieren). Schließlich gilt für
`codepaths` das **geteilte Referenz-Ventil** `ignore-refs`
([`DC-FA-REF-001`](#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)),
das **bestimmte Ziel-Pfade** referenz-weit von der Existenz-/Anker-Prüfung ausnimmt
(Register bewusst entfernter/historischer Artefakte); die frühere modul-lokale Liste
`codepaths.ignore-refs` bleibt als **Alias** gültig (kein Config-Bruch). Die drei
Ventil-Achsen dieses Moduls: Zeile (`d-check:ignore`), Datei (`exempt-paths`), Ziel
(`ignore-refs`).

Optional verifiziert `codepaths` bei aktivem `check-lines` (Default aus,
byte-identisch) zusätzlich die **Zeilen-Referenz** eines Pfads: trägt der Wert ein
`:<von>`- oder `:<von>-<bis>`-Suffix (`datei:zeile`-Konvention), muss die Zieldatei
existieren und mindestens `<bis>` Zeilen haben — sonst `citation-out-of-range` — und
`<von> ≤ <bis>` erfüllen — sonst `citation-inverted-range`. Das fängt **Zitat-Fäule**
nach einem Tag-Bump (die Zeile zeigt hinter das Datei-Ende) und invertierte/vertippte
Bereiche. Ohne `check-lines` bleibt das Verhalten unverändert: wie bisher wird nur ein
einzelnes `:<zahl>` abgetrennt (ein Bereich `:<von>-<bis>` bleibt dabei unberührt) —
byte-identisch. Der Zeilen-Check setzt die bestehende Existenz-Prüfung voraus (fehlende
Datei bleibt `codepath-missing`).

**Akzeptanzkriterien:**

- **Boundary (Marker in Inline-Code ist Erwähnung):** Given eine Zeile, die den Zeilen-Marker in **Backticks** nennt und daneben einen nicht auflösbaren Pfad in Inline-Code trägt, when das Modul läuft, then **ein Befund** `codepath-missing` — dieselbe Zeile mit **freiem** Marker meldet keinen (Gegenprobe). **Benannte Grenze:** eine Code-Spanne desselben Absatzes verschluckt auch einen freien Marker (Falsch-Rot), und ein unpaariger Backtick weiter oben im Absatz kippt die Parität, sodass eine Erwähnung doch wirkt (stilles Grün).
- **Boundary (Form folgt der Eingabe):** Given eine Zeile mit **blankem** Marker (ohne Kommentar-Klammer) und einem nicht auflösbaren Pfad bzw. einer unverlinkten Kennung, when `codepaths` bzw. `ids` läuft, then **ein Befund** — dieselbe Zeile mit `<!-- d-check:ignore -->` meldet keinen (Gegenprobe). Given dieselbe Zeile bei `versions`, then **kein** Befund: dort ist der Marker ein **Token**. Given eine **Diagramm-Fence-Zeile** (oder die Öffnungszeile) mit blankem Token und einer undefinierten Kennung, when `diagrams` läuft, then **kein** Befund — eine Prosa-Zeile ist dort nie Eingabe, das Kriterium braucht deshalb seine eigene Konstellation. **Konservativ:** ein `>` im Kommentar vor dem Marker lässt ihn nicht gelten.
- **Happy Path:** Given ein Inline-Code-Span `` `docs/plan/adr/` `` auf ein existierendes Verzeichnis und das konfigurierte Präfix `docs/`, when das Modul `codepaths` läuft, then kein Befund.
- **Boundary:** Given ein Inline-Code-Span mit nicht existierendem Pfad und ein Kommentar `d-check:ignore` auf derselben Zeile, when das Modul läuft, then kein Befund — und der Marker hat keinerlei Wirkung auf Befunde anderer Module derselben Zeile.
- **Negative:** Given ein Inline-Code-Span `` `../fehlt.md` ``, dessen Ziel nicht existiert (oder die Repository-Wurzel verlässt), when das Modul läuft, then ein Befund mit Datei, Zeile, Ziel und Grund, Exit-Code 1.
- **Ventil (exempt-paths):** Given ein Inline-Code-Span mit nicht existierendem Pfad in einer Datei, die ein `codepaths.exempt-paths`-Glob matcht, when das Modul läuft, then kein Befund — und ohne gesetztes `exempt-paths` ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)).
- **Ventil (ignore-refs):** Given ein Inline-Code-Span auf einen nicht existierenden Pfad, der ein `codepaths.ignore-refs`-Glob matcht, when das Modul läuft, then **kein** Befund — referenz-weit, auch in einer sonst geprüften Datei; ein nicht von `ignore-refs` gedeckter fehlender Pfad bleibt `codepath-missing`, und ohne gesetztes `ignore-refs` ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)).
- **Zeilen-Check (out-of-range):** Given `check-lines: true` und ein Inline-Code-Span `` `docs/a.md:999` `` auf eine existierende Datei mit weniger als 999 Zeilen, when das Modul läuft, then ein Befund `citation-out-of-range`; ohne `check-lines` (Default) ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)).
- **Zeilen-Check (inverted):** Given `check-lines: true` und ein Span `` `docs/a.md:20-10` `` (`von > bis`), when das Modul läuft, then ein Befund `citation-inverted-range`.

**Out-of-Scope:** Pfad-Erkennung im Fließtext ohne Inline-Code; Pfade in Fenced-Code-Blöcken; Opt-out-Marker für andere Module; semantische Prüfung, ob der referenzierte **Inhalt** an der Zeile passt (das prüft das Modul `citations`, [`DC-FA-CITE-001`](#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in)); Zeilen-Granularität von `exempt-paths` (datei-weit wie `ids`; für Zeilen der `d-check:ignore`-Marker).

---

### DC-FA-CITE-001 — Verbatim-Zitat-Verifikation (Modul `citations`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `citations` wird ein per Direktive
`d-check:cite` ausgezeichnetes **Zitat** gegen die von ihm zitierte Quell-Spanne
geprüft — die Zusage „wortgleich" wird von einer Behauptung zu einem **gemessenen
Property**. Die Direktive ist ein HTML-Kommentar
`<!-- d-check:cite <pfad>:<von>-<bis> -->` unmittelbar vor dem Zitat: einem
`>`-Blockquote **oder** dem folgenden inline-Zitat-Span (`„…"` / `"…"`) — die realen
Zitate stehen **inline** und beim Verfassen umgebrochen, nicht als `>`-Blöcke.
**Was „unmittelbar" heißt, ist genau bestimmt:** Leerzeilen trennen **nicht** —
die Direktive paart mit dem nächsten nicht-leeren Kandidaten, auch über mehrere
Leerzeilen hinweg. Ein **Fenced-Block** dagegen trennt: liegt einer zwischen
Direktive und Kandidat, folgt der Direktive kein Zitat, und das ist der
fail-closed-Fall, keine Paarung über den Block hinweg. Das
Modul liest die Zeilen `<von>`–`<bis>` der Zieldatei und vergleicht sie mit dem
Zitattext **whitespace-normalisiert** (Läufe aus Leerzeichen/Umbruch → ein
Leerzeichen): der normalisierte Zitattext muss ein **zusammenhängender Teilstring**
der normalisierten Quell-Spanne sein. So bleibt der Vergleich zeichengenau auf dem
**Inhalt** (Markdown-Auszeichnung, Satzzeichen, Groß-/Kleinschreibung zählen),
tolerant nur gegenüber Umbruch — ein mitten in der Zeile beginnendes/endendes,
re-wrapptes Zitat besteht, jede echte Wort-Abweichung bricht (`citation-mismatch`).
Das Modul hat eine **eigene** Erkennung (es parst die Direktive) und teilt die
Pfad-/Zeilen-Auflösung von `codepaths`/`links`
([`DC-FA-CODE-001`](#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in));
strikt opt-in, hermetisch (nur Datei-Lesen, kein git/Netz —
[`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)) und
default-aus **byte-identisch** ([`DC-QA-02`](#dc-qa-02--determinismus)).
**Was als Direktive gilt, ist eine Prosa-Frage und wird geteilt beantwortet:**
die Direktive wirkt nur **außerhalb** von Fenced-Blöcken **und außerhalb von
Inline-Code** — die Syntax in Backticks ist eine **Erwähnung**, keine Direktive,
und trägt eine Zeile beides, gilt die freie. So bleibt die Dokumentation der
Direktive dokumentierbar, ohne den Lauf abzubrechen. Der **Zitattext** dagegen
wird **roh** verglichen: dort sind die Bytes der Gegenstand, nicht ihre
Prosa-Rolle. **Preis, benannt, und er reicht weiter als „steht in Backticks":**
eine **echte** Direktive wird still übersprungen, sobald eine Code-Spanne
**desselben Absatzes** sie umschließt — die Spanne darf eine Zeile vorher öffnen
und eine Zeile später schließen, die Direktive selbst unverklammert dazwischen
stehen. Es ist derselbe Preis, den jedes prosa-lesende Modul zahlt. **Zwei
weitere Grenzen, die derselben Wahl folgen:** ein Pfad in Backticks wird durch
das Strippen zum **fehlenden** Pfad und damit fail-closed statt zu einem Befund
mit leerem Ziel; und eine Code-Spanne zwischen Kommentar-Öffner und Marker
verschwindet ebenso — sie kann eine Direktive damit **erzeugen**. **Die
Ziel-Seite trägt keine dieser Zusagen:** die zitierte Quell-Spanne wird **roh
und typunabhängig** gelesen, ohne Fence-Bewusstsein und ohne Strippen, und die
Zieldatei muss weder Markdown sein noch in der Scan-Menge liegen — `scan.ignore`
gilt der **prüfenden** Datei, nicht dem Ziel.
**Zitat-Fäule** (Zieldatei fehlt bzw. Spanne über das Datei-Ende ⇒
`citation-out-of-range`; `von > bis` ⇒ `citation-inverted-range`) ist ein **Befund**
(kohärent zum `codepaths`-Zeilen-Check), Exit 1. **Fail-closed** (Exit 2) nur bei
strukturell **unbrauchbarer** Direktive (malformter Span, kein folgendes Zitat); ein
Repo-Escape des Ziels ist wie bei `codepaths`/`links` ein Befund. `d-check` kennt
damit **zwei** Direktiven. Ihre Platzierung ist **nicht** dieselbe, und das ist
eine benannte Grenze: die Zitat-Direktive wirkt nur außerhalb von Fences und
außerhalb von Inline-Code (oben), die Ventil-Direktive des Moduls
[`DC-FA-CODE-001`](#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
wird weiterhin auf der **rohen** Zeile erkannt und wirkt damit auch aus einer
Backtick-Spanne heraus. Dieselbe Frage — *ist das eine Direktive* — hat im
Produkt derzeit zwei Antworten; welche bleibt, ist offen.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine Direktive `<!-- d-check:cite a.md:3-4 -->` vor einem Zitat, dessen whitespace-normalisierter Text ein Teilstring der normalisierten Zeilen 3–4 von `a.md` ist — auch wenn das Zitat mitten in Zeile 3 beginnt, in Zeile 4 endet oder anders umbricht, when das Modul `citations` läuft, then kein Befund.
- **Negative:** Given denselben Aufbau, aber der Zitattext weicht in mindestens einem Wort/Zeichen ab (nicht nur Whitespace), when das Modul läuft, then ein Befund `citation-mismatch` (Datei, Zeile, Ziel), Exit-Code 1.
- **Boundary (Fenced-Block trennt):** Given eine `d-check:cite`-Direktive, zwischen der und dem Zitat ein **Fenced-Block** liegt, when das Modul läuft, then **kein** Zitat gilt als gefunden und der Lauf bricht **fail-closed** ab (Exit 2) — dieselbe Datei **ohne** den Block paart normal. Leerzeilen dazwischen trennen dagegen **nicht**.
- **Boundary (Zitat-Fäule = Befund, nicht fail-closed):** Given eine `d-check:cite`-Direktive, deren Spanne die Zieldatei überschreitet (Zitat-Fäule nach einem Tag-Bump), when das Modul läuft, then ein Befund `citation-out-of-range`, Exit-Code 1 — **kohärent** zum `codepaths`-Zeilen-Check, **nicht** Exit 2; ein ungültiger Bereich (`von > bis`) ⇒ `citation-inverted-range`. Fail-closed (Exit 2) bleibt der malformten Direktive bzw. dem fehlenden folgenden Zitat vorbehalten; ohne `citations`-Modul ist jeder Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)).
- **Boundary (Inline-Code ist Erwähnung, nicht Direktive):** Given eine Datei, die die Direktiv-Syntax in **Inline-Code** schreibt — auch wohlgeformt und auch über eine absatzweite Backtick-Spanne umgebrochen —, when das Modul läuft, then **kein** Befund und **kein** fail-closed; dieselbe Direktive **ohne** Backticks wird geprüft (Gegenprobe), und trägt eine Zeile beides, gilt die **freie** Direktive — nicht die zuerst stehende Erwähnung.
- **Boundary (der Preis reicht bis zur Absatzgrenze):** Given eine **freie**, unverklammerte Direktive, die von einer Code-Spanne **desselben Absatzes** umschlossen wird (Spanne öffnet eine Zeile vorher, schließt eine Zeile später), when das Modul läuft, then **kein** Befund — dieselbe Datei **ohne** die zwei Backticks meldet ihn (Gegenprobe). Given eine Direktive, deren **Pfad** in Backticks steht, then **fail-closed** (fehlender Pfad), nicht ein Befund mit leerem Ziel.

**Out-of-Scope:** Zitate ohne `d-check:cite`-Direktive (das Modul prüft nur ausgezeichnete Zitate — kein Prosa-Scanning); sehr kurze Zitate (< 16 Zeichen normalisiert) — zu schwache Teilstring-Diskriminierung, ungeprüft; Normalisierung über Whitespace/Umbruch hinaus (Markdown-Auszeichnung, Satzzeichen, Groß-/Kleinschreibung zählen); freie Zahlen und Prosa-Quantoren mit externer Grundwahrheit („42 Dateien im ZIP", „fast alle") — bleiben Review-Territorium.

---

### DC-FA-SPAN-001 — Markdown-Span-Artefakte (Modul `spans`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `spans` werden **drei**
deterministische Artefakt-Klassen gemeldet, bei denen Markdown auf
GitHub nachweislich anders rendert, als der Quelltext nahelegt:

1. **Ungeschlossene Code-Spans** (`span-unclosed`): eine öffnende
   Backtick-Folge, die unmittelbar von Nicht-Whitespace gefolgt wird
   und im Absatz keine gleich lange schließende Folge findet
   (Absatzgrenzen: Leerzeile und Fence — dieselbe Absatz-Semantik wie
   die Vorverarbeitung der übrigen Module). Solche Opener kippen die
   Backtick-Parität des restlichen Absatzes; nachfolgende Links und
   Kennungen rendern zerrissen. Alleinstehende Backtick-Folgen
   (beidseitig Whitespace oder Zeilenrand) gelten als beabsichtigt
   literal und sind kein Befund (konservative Erkennung).
2. **Verschachtelte Link-Artefakte** (`span-nested-link`): Link-Syntax
   im Linktext eines weiteren Links (schließende Ziel-Klammer
   unmittelbar gefolgt von einer weiteren Linktext-Schließung mit
   Ziel-Öffnung) außerhalb von Fences und Code-Spans — für Leser
   sieht das wie ein Link aus, wird aber nur teilweise als Link
   geparst. Eine **Bildreferenz als Linktext** (Badge-Muster) ist
   legales Markdown und kein Treffer.
3. **Unbalancierter Fenced-Code-Block** (`fence-unclosed`): eine
   Fence-Öffnung, die bis zum **Dateiende** keinen Schluss findet.
   Alles dahinter gilt der Vorverarbeitung als Code und wird
   übersprungen — der Autor schaltet damit unbemerkt die Prüfung des
   Restdokuments ab. Es ist dieselbe Aussage wie Klasse 1, eine Ebene
   höher: eine Öffnung ohne Schluss, die alles Folgende umdeutet.
   **Reichweite ist die Datei, nicht der Absatz** (ein Fence kennt
   keine Absatzgrenze — er *ist* eine).

   Was ein *Schluss* ist, liest d-check an zwei Stellen verschieden:
   der naive Toggle zählt jede Fence-Zeile, der CommonMark-Schluss
   verlangt gleiches Zeichen und mindestens gleiche Länge; welche die
   richtige ist, bleibt bewusst offen. Der Befund wertet
   **beide** aus und meldet, sobald **eine** von beiden am Dateiende
   offen steht — denn genau dann überspringt mindestens ein Modul den
   Rest. Nur eine auszuwerten hieße, die andere unbewacht zu lassen.

   Der Befund steht an der Öffnungszeile der strengen Lesart, wenn diese
   offen ist, sonst an der zuletzt öffnend gewerteten Zeile. In beiden Fällen
   ist das eine **Fundstelle**, nicht zwingend die Reparaturstelle: welche von
   mehreren Öffnungen fehlt, ist grundsätzlich nicht entscheidbar. Bei reiner
   Paritätskippung sind gleich lange Öffner ununterscheidbar; aber auch die
   strenge Lesart zeigt daneben, sobald eine längere Fence-Zeile eine kürzere
   Öffnung geschlossen hat und erst eine spätere Öffnung offen bleibt.

   **Grenze:** die Prüfung greift nur auf Dateien im Scan-Scope. Zwei Klassen
   von Dateien liegen systematisch daneben — Module, die ihre Eingabe selbst
   benennen (Post-Pässe über deklarierte Verzeichnisse), und **Zieldateien**
   außerhalb der Scan-Wurzeln, aus denen Module lesen: `matrix` den Status,
   `anchors` und `codepaths` die Slugs, `diagrams` und `versions` ihre
   deklarierten Quellen, `pins` die gehashte Heading-Section. Bei `pins` ist die
   Folge **still**: verdeckt ein offener Fence die Überschrift, ist der Anker
   unauflösbar, und die Drift-Prüfung entfällt kommentarlos statt zu melden.

   **Ebenfalls außerhalb:** der eingerückte Code-Block (vier Leerzeichen). Er
   ist in keiner Lesart modelliert — Fence-Zeilen darin zählen wie überall, und
   eine gerade Anzahl davon verschluckt den Text dazwischen, ohne dass eine
   Lesart unbalanciert endet.

Der Befund nennt Datei, Zeile (Opener- bzw. Muster-Zeile), den Grund
und als Ziel je nach Klasse die betroffene Backtick-Folge (Klasse 1),
das Muster (Klasse 2) oder die getrimmte Fence-Zeile samt Infozeile
(Klasse 3). Es gilt die
allgemeine Opt-out-Regel unverändert: deterministische Befunde werden
behoben, nicht unterdrückt — der Zeilen-Marker `d-check:ignore`
bleibt auf das Modul `codepaths` beschränkt
([`DC-FA-CODE-001`](#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)).

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Absatz mit ausschließlich balancierten Code-Spans — auch solchen, die einen Zeilenumbruch enthalten —, when das Modul `spans` läuft, then kein Befund.
- **Boundary:** Given eine alleinstehende literale Backtick-Folge (beidseitig Whitespace), when das Modul läuft, then kein Befund.
- **Negative:** Given ein Listenpunkt, dessen öffnende Backtick-Folge unmittelbar vor Nicht-Whitespace steht und im Absatz ungeschlossen bleibt, when das Modul läuft, then ein Befund `span-unclosed` mit Datei, Zeile und Grund, Exit-Code 1.
- **Negative (unbalancierter Fence):** Given eine Datei, deren Fence-Öffnung bis zum **Dateiende** keinen Schluss findet, when das Modul läuft, then ein Befund `fence-unclosed` an der **Öffnungszeile**, Exit-Code 1 — unabhängig davon, ob dahinter noch Inhalt steht.
- **Negative (nur eine Lesart offen):** Given eine Datei, in der eine Backtick-Fence von einer Tilden-Zeile „geschlossen“ wird — für den naiven Toggle balanciert, für den CommonMark-Schluss offen —, when das Modul läuft, then ein Befund `fence-unclosed`, Exit-Code 1.
- **Boundary (balanciert):** Given eine Datei mit beliebig vielen, jeweils **unter beiden Lesarten** geschlossenen Fences — einschließlich einer legalen Verschachtelung, in der ein längerer Fence einen kürzeren zeigt —, when das Modul läuft, then **kein** `fence-unclosed`.
- **Boundary (Einrückung):** Given eine Fence-Zeile, die mit Unicode-Whitespace statt Space/Tab eingerückt ist, when das Modul läuft, then zählt sie **nicht** als Fence — dieselbe Trimmung wie die Vorverarbeitung, sonst bewachte der Befund einen anderen Automaten als den bewachten.

**Out-of-Scope:** Emphasis-/Bold-Artefakte (`*`, `_`); Syntaxfehler in Fenced-Code-Blöcken; Reference-Style-Definitionen; automatische Korrektur; ein Opt-out-Marker für dieses Modul.

---

### DC-FA-HOST-001 — Host-lokale absolute Pfade (Modul `hostpaths`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `hostpaths` werden
absolute Pfade gemeldet, die ein Maschinen-Layout statt einer
Repo-Struktur beschreiben — sie funktionieren nur auf einem konkreten
Host und leaken dessen Verzeichnis-Aufbau. Erkannt werden, jeweils
als erstes Segment eines absoluten Pfads:

1. **Unix-Host-Präfixe** (Default-Liste, als Wurzel-Verzeichnisnamen
   deklariert: Development, home, Users, Volumes, mnt, media;
   per `hostpaths.prefixes` ersetzbar — tmp gehört bewusst nicht
   dazu: ein POSIX-Standard-Laufzeitort, dessen Erwähnung legitime
   Betriebs-Doku ist),
2. **Windows-Laufwerkspfade** (Laufwerksbuchstabe, Doppelpunkt,
   Backslash) und **UNC-Pfade** (doppelter Backslash plus
   Servername) — beide immer, nicht konfigurierbar.

Geprüft werden Prosa-Zeilen **einschließlich Inline-Code** (dort
leben solche Pfade typischerweise); Fenced-Code-Blöcke sind
ausgenommen — Beispiel- und Lehrinhalte mit bewussten Host-Pfaden
gehören in Fences. Vorbedingung ist eine Wortgrenze (kein
unmittelbar vorangehendes URL-, Pfad- oder Wortzeichen); schließende
Satzzeichen werden vom gemeldeten Pfad abgetrennt. Grund-Code
`hostpath-forbidden` mit Datei, Zeile und gefundenem Pfad. Es gibt
**keinen Opt-out-Marker** (der Zeilen-Marker `d-check:ignore` bleibt
auf das Modul `codepaths` beschränkt,
[`DC-FA-CODE-001`](#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)) —
die Auswege sind Fences für beabsichtigte Beispiele und die
Präfixliste für repo-spezifische Sonderfälle.

**Akzeptanzkriterien:**

- **Happy Path:** Given Dokumentation, deren Pfad-Angaben relativ oder Repo-Wurzel-absolut sind, when das Modul `hostpaths` läuft, then kein Befund.
- **Boundary:** Given ein host-lokaler absoluter Pfad innerhalb eines Fenced-Code-Blocks, when das Modul läuft, then kein Befund.
- **Negative:** Given eine Prosa-Zeile oder ein Inline-Code-Span mit einem Pfad unterhalb eines deklarierten Host-Präfixes, when das Modul läuft, then ein Befund `hostpath-forbidden` mit Datei, Zeile, Pfad und Grund, Exit-Code 1.

**Out-of-Scope:** Erkennung relativer Pfade (Aufgabe von `links`/`codepaths`); URL-Pfade hinter Schemata; Pfade in Fenced-Code-Blöcken; ein Opt-out-Marker für dieses Modul; automatische Umschreibung.

---

### DC-FA-DIAG-001 — Kennungs-Konsistenz in Diagramm-Fences (Modul `diagrams`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `diagrams` öffnet d-check
**gezielt benannte Diagramm-Fences** (Default `mermaid`; per `diagrams.fences`
gesetzt) und extrahiert daraus **Kennungen per Regex** (Muster wie beim Modul
[`ids`](#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids), je Muster eine
Definitionsquelle `defined-in`). Geprüft wird **Existenz, nicht Linkpflicht**:
jede im Diagramm gefundene Kennung muss in ihrer `defined-in`-Quelle **definiert**
sein. **Definiert** heißt: die Kennung kommt in der `defined-in`-Datei als
eigenständiges Kennungs-Token **außerhalb von Fences** vor — als Heading **oder**
in Fließtext/Tabelle; bewusst **nicht** heading-zentriert (anders als das Modul
`ids`), damit auch **tabellen-definierte** Komponenten-Kennungen (z. B.
Architektur-IDs) erfasst werden. Eine Kennung ohne solche Definition ist der
Befund `diagram-id-undefined` (Datei, Zeile im Fence, Kennung). **Link-Policy gilt hier nicht** (Mermaid kennt keine
Markdown-Links; eine Diagramm-Kennung *kann* nicht verlinkt werden), daher das
von [`DC-FA-ID-001`](#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
abweichende Existenz-Kriterium.

Das Modul **öffnet als bislang einziges** gezielt benannte Diagramm-Fences; alle
übrigen Module behandeln Fence-Inhalt als opak (vgl. die Vorverarbeitung des
Moduls [`links`](#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)).
**Nur** die in `diagrams.fences` gelisteten Sprachen werden geöffnet, **kein** Mermaid-/
Grammatik-Parsing (reine Token-Extraktion über Rohtext) — deterministisch
([`DC-QA-02`](#dc-qa-02--determinismus)), read-only und netzlos
([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)), ohne
zusätzliche Toolchain. Strikt opt-in (Default aus): ohne `diagrams`-Block ist der
Befundsatz **byte-identisch** zum Lauf ohne das Modul. Reiht sich in die
opt-in-Module ([`DC-FA-CLI-002`](#dc-fa-cli-002--regelmodul-auswahl)) ein.

**Ventile — Parität zu den Modulen mit konfigurierbaren Mustern.** `diagrams`
war das einzige Modul, das **eigene Muster** konfiguriert und **kein** Ventil
trug: weder eine Datei-Ausnahme noch den Zeilen-Marker. Wer es aktiviert,
konnte ein Beispiel-Diagramm mit erfundener Kennung nur loswerden, indem er
ganze Dateibäume aus dem Scan-Bereich nahm — eine Vermeidung, keine Ausnahme.
Es trägt deshalb dieselben zwei Ventile wie
[`DC-FA-ID-001`](#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids),
[`DC-FA-CODE-001`](#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
und [`DC-FA-VER-001`](#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in):
die Glob-Liste `diagrams.exempt-paths` (datei-weit) und den Zeilen-Marker
`d-check:ignore`. **Nicht** gedeckt sind damit die Module ohne konfigurierbare
Muster, die ebenfalls kein Ventil tragen
([`DC-FA-HOST-001`](#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in),
[`DC-FA-PIN-001`](#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in),
[`DC-FA-SPAN-001`](#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)) — eine
offene Fläche, keine Zusage dieser Version.

**Die Ventile wirken scan-seitig, nicht auf der Definitionsmenge.** Weder
`exempt-paths` noch der Marker beeinflussen, was in der `defined-in`-Quelle als
**definiert** gilt: eine dort als illustrativ markierte Kennung definiert
weiter. Das ist eine Festlegung, keine Auslassung — die Ventile unterdrücken
**Befunde**, sie ändern nicht die Wahrheit über die Definitionsmenge. Wer eine
Beispiel-Kennung aus der Quelle heraushalten will, schreibt sie nicht als
Token hinein.

**Der Marker ist ein Token, kein HTML-Kommentar.** Das ist keine Feinheit,
sondern der Unterschied zwischen Prosa und Fence: in Prosa steht der Marker
üblicherweise in `<!-- … -->`, in einem `mermaid`-Fence ist das kein Kommentar,
sondern Diagramm-Text. Das Modul sucht das **Token** auf der Zeile; wie es der
Autor vor dem Renderer versteckt, ist Sache der Diagramm-Sprache (in Mermaid
`%% d-check:ignore`). Und er wirkt an **zwei** Orten: auf einer Diagramm-Zeile
nimmt er **diese Zeile** aus, auf der **Öffnungszeile** des Fence den
**ganzen** Block. Ohne die zweite Stelle wäre die intuitive Platzierung — der
Marker am Diagramm-Anfang — wirkungslos, und ein Beispiel-Diagramm bräuchte
den Marker auf jeder Kennungs-Zeile.

**Akzeptanzkriterien:**

- **Happy Path:** Given `diagrams` aktiv (`fences: [mermaid]`) und ein `mermaid`-Block, dessen Kennungen alle in ihrer `defined-in`-Quelle vorkommen, when `d-check` läuft, then kein Befund, Exit 0; ein read-only gemountetes Repository genügt.
- **Boundary:** Given ein `mermaid`-Block mit genau einer Kennung ohne Definition in `defined-in`, when das Modul läuft, then genau ein `diagram-id-undefined` (Datei, Zeile im Fence, Kennung), Exit 1; dieselbe Kennung in einem **nicht** gelisteten Fence (z. B. `bash`) oder außerhalb jedes Fence bleibt für dieses Modul unberührt.
- **Negative:** Given **kein** `diagrams`-Block in der Konfiguration, when `d-check` läuft, then ist der Befundsatz byte-identisch zum Lauf ohne das Modul ([`DC-QA-02`](#dc-qa-02--determinismus)) und es wird nichts geschrieben ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **fail-closed (`defined-in`):** Given ein **Verzeichnis** statt einer Datei als `defined-in`, when `d-check` startet, then Exit 2 ohne Prüfung — das Modul liest die Quelle, ein Verzeichnis lieferte eine leere Definitionsmenge und damit einen Befund-Sturm ohne erkennbare Ursache (Abgrenzung zu [`DC-FA-ID-001`](#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids), wo ein Verzeichnis als Ziel zulässig ist).
- **Ventil (Datei):** Given eine Datei, deren Pfad einem `diagrams.exempt-paths`-Glob entspricht und die ein Diagramm mit undefinierter Kennung trägt, when `d-check` läuft, then kein Befund für diese Datei — datei-weit, unabhängig vom Scan-Bereich.
- **Ventil (Zeile):** Given eine Diagramm-Zeile mit undefinierter Kennung und dem Token `d-check:ignore` (in Mermaid als `%%`-Kommentar), when `d-check` läuft, then kein Befund für **diese** Zeile; eine Kennung auf einer **anderen** Zeile desselben Blocks wird weiter gemeldet.
- **Ventil (ganzer Block):** Given den Marker auf der **Öffnungszeile** eines Diagramm-Fence, when `d-check` läuft, then bleibt der **ganze** Block ungeprüft — auch mehrere Kennungs-Zeilen.
- **Ventil (Nicht-Wirkung):** Given weder `exempt-paths` noch die Zeichenfolge `d-check:ignore` in einer gelisteten Fence, when `d-check` läuft, then ist der Befundsatz byte-identisch zur Fassung ohne die Ventile ([`DC-QA-02`](#dc-qa-02--determinismus)); ein ungültiges oder leeres `exempt-paths`-Glob ⇒ Exit 2 vor dem Lauf. **Der Marker ist nicht schlüssel-gebunden:** wie bei `ids`/`codepaths`/`versions` wirkt er allein über den Zeilen-**Inhalt** — eine Diagramm-Zeile, die die Zeichenfolge ohnehin trägt, wird ohne jede Konfigurations-Änderung stumm. Die Byte-Identität gilt deshalb für Bäume **ohne** diese Zeichenfolge in einer gelisteten Fence.

**Out-of-Scope:** Mermaid-**Syntax**-/Rendering-Validierung (gehört in einen Diagramm-Linter, braucht den Original-Parser); Grammatik-Parsing der Diagramm-Sprache; Link-Policy innerhalb von Fences (technisch unmöglich); Vollständigkeit „jede definierte Kennung erscheint im Diagramm" (separater Folge-Check); Referenzrichtungs-Prüfung von Diagramm-Kennungen (Sache einer späteren `matrix`-Erweiterung); Nicht-Mermaid-Diagrammsprachen in der ersten Fassung.

---

### DC-FA-VER-001 — Versions-Pin-Konsistenz (Modul `versions`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `versions` prüft `d-check`,
dass alle **Versions-Pins** im Repo mit der **deklarierten aktuellen Version**
übereinstimmen. Ein Versions-Pin ist ein Vorkommen des konfigurierten Musters
`versions.pin-pattern` (z. B. `ghcr.io/…/d-check:v\d+\.\d+\.\d+`); die aktuelle
Version stammt aus `versions.current-from` (Default `version.md#aktuell` — der
Versionsstring im so adressierten Span). Weicht ein Pin von der aktuellen
Version ab → Befund `version-stale` (Datei:Zeile, gefundene vs. erwartete
Version).

Anders als die übrigen Module liest `versions` die Pins **auch innerhalb von
Fenced-Code-Blöcken** — Versions-Pins leben fast ausschließlich in
Kommando-Beispielen. Das ist eine bewusste, auf das `pin-pattern` gescopte
Ausnahme von der Fence-Opazität der Vorverarbeitung (kein Sprach-Parser, reiner
Muster-Scan). **Grenze der Anker-Menge:** die geteilte Anker-Erkennung ist
**permissiv** — sie liest HTML-Tags auf allen Prosa-Zeilen, also auch in einer
HTML-**Kommentar**-Zeile und in einem **eingerückten** Code-Block (den die
Vorverarbeitung nicht kennt). Für [`DC-FA-ANCH-001`](#dc-fa-anch-001--heading-anker-validierung-modul-anchors)
ist eine zu große Anker-Menge folgenlos (das Modul schweigt), für `versions`
**nicht**: der adressierte Span beginnt beim **ersten** Vorkommen, und die
„aktuelle Version“ kann dadurch aus einem auskommentierten Beispiel stammen —
ein Falsch-Rot. Die Verengung ist eine eigene Entscheidung mit eigener Messung
und hier ausdrücklich **nicht** getroffen. **Die Fence-Ausnahme gilt dem Pin,
nicht dem Anker:** welcher Span
`versions.current-from` adressiert, entscheidet dieselbe Anker-Erkennung wie in
[`DC-FA-ANCH-001`](#dc-fa-anch-001--heading-anker-validierung-modul-anchors) —
ein Heading-Slug oder ein Inline-HTML-Anker, beide **außerhalb** von Fenced-Code
und Inline-Code. Sonst gäbe derselbe Lauf zwei Antworten auf dieselbe Frage.
Der adressierte Span selbst wird danach **roh** gelesen, einschließlich seiner
Fences. Zwei Ventile wie bei
[`DC-FA-ID-001`](#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids): die
Glob-Liste `exempt-paths` (historische Pins in Planning-`done/`-Slices,
`CHANGELOG.md` und der Lastenheft-Historie bleiben unberührt) und der
Zeilen-Marker `d-check:ignore`.

**Mehrere Muster-Quellen-Paare.** `pin-pattern`, `current-from` und
`exempt-paths` beschreiben **ein** Paar aus Pin-Muster, Versions-Quelle und
Ausnahmen. Ein Repo führt aber oft mehrere unabhängige Versions-Reihen — das
eigene Image gegen das Release-Register, ein fremder Baseline-Tag gegen den Pin
im Konventionsspeicher. Dafür trägt `versions.patterns` eine **Liste** solcher
Paare. Jedes Paar wird für sich ausgewertet: eigenes Muster, eigene Quelle (mit
demselben Default), eigene `exempt-paths` — und **seine** Quell-Datei ist nur
von **seiner** Prüfung ausgenommen, nicht von der der anderen. Die Kurzform
(die drei Schlüssel direkt unter `versions`) bleibt gültig und **ist** die
einelementige Liste; es gibt genau einen Auswertungspfad, nicht zwei. **Beide
Schreibweisen zugleich sind ein Nutzungsfehler** (Exit 2): welche
`exempt-paths` dann für welches Muster gälten, wäre nicht ablesbar, und eine
Voreinstellung, die man raten muss, ist keine. Ob eine Schreibweise gesetzt
ist, entscheidet die **Anwesenheit** des Schlüssels, nicht sein Wert:
`pin-pattern: ""` neben `patterns` ist eine Mischform. Ein Schlüssel **ohne
Wert** ist im YAML von einem fehlenden nicht unterscheidbar und zählt als
fehlend — eine benannte Grenze.

**Eine Befund-Adresse, alle Erwartungen.** Je Zeile entsteht höchstens **ein**
Befund pro gefundenem Pin-Wert: die Befund-Adresse (Datei, Zeile, Regel,
`target`, Grund-Code) kann zwei Befunde an derselben Stelle nicht
unterscheiden, und die geteilte Nachrunde verwürfe den zweiten samt seiner
Erwartung. Treffen mehrere Paare denselben Wert, nennt die **Nachricht**
deshalb alle Erwartungen — in Deklarationsreihenfolge der Paare, jede mit
ihrer Quelle. Die **Ausgabe**-Reihenfolge ist die geteilte Sortierung, nicht
die Deklarationsreihenfolge. Bei genau einem Paar bleibt der Wortlaut der
Nachricht unverändert; ab zwei Paaren benennt er die Fundstelle
(`versions.patterns[i].current-from`), auch im fail-closed-Fall.

Strikt opt-in (Default aus): ohne `versions`-Block ist der Befundsatz
byte-identisch zum Lauf ohne das Modul
([`DC-QA-02`](#dc-qa-02--determinismus)) und es wird nichts geschrieben
([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
**Fail-closed:** ist `versions.current-from` nicht auflösbar oder trägt der
adressierte Span keine erkennbare Version, bricht der Lauf mit einem
Nutzungsfehler (Exit 2, vgl.
[`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)) statt einer stillen
Grün-Meldung. `version-stale` ist in dieser Version **diagnose-only** — es
liefert keinen `--repair`-Hunk; der konservative Auto-Bump (Pin → aktuelle
Version, deterministisch ableitbar) folgt als eigener Change Request an
[`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch).

**Akzeptanzkriterien:**

- **Happy Path:** Given `versions` aktiv (`pin-pattern` gesetzt, `current-from`
  nennt `v0.27.0`) und alle Pins zeigen `v0.27.0`, when `d-check` läuft, then
  kein Befund, Exit 0; ein read-only gemountetes Repository genügt.
- **Negative:** Given ein Pin `…:v0.26.0` (auch innerhalb eines
  Fenced-Code-Blocks), während `current-from` `v0.27.0` nennt, außerhalb
  `exempt-paths` und ohne `d-check:ignore`, when `d-check` läuft, then ein
  Befund `version-stale` (Datei:Zeile, gefunden `v0.26.0`, erwartet `v0.27.0`),
  Exit 1.
- **Boundary (eine Anker-Antwort):** Given eine `current-from`-Adresse, deren Anker ein **Duplikat-Slug** (`#alt-1`), ein **prozent-kodiertes** Fragment (`#a%20b`) oder ein Inline-HTML-Anker ist, when `versions` und `anchors` im **selben** Lauf aktiv sind, then lösen beide denselben Anker gleich auf — was `anchors` für gültig hält, löst `current-from` auf, und was `anchors` nicht kennt (Anker in Fenced-/Inline-Code, `data-id`, `name` an einem Nicht-`a`-Element, anker-förmige Prosa ohne Tag), führt hier zum fail-closed-Abbruch statt zu einer stillen Auflösung.
- **Boundary (Ventile):** Given ein Pin `…:v0.1.0` in einer `exempt-paths`-Datei
  (z. B. ein Planning-`done/`-Slice) und ein zweiter Pin auf einer Zeile mit
  `d-check:ignore`, when `d-check` läuft, then kein Befund für beide.
- **Happy Path (zwei Paare):** Given `versions.patterns` mit zwei Paaren —
  verschiedene Muster, verschiedene `current-from`-Quellen — und alle Pins
  beider Reihen tragen die je eigene aktuelle Version, when `d-check` läuft,
  then kein Befund, Exit 0.
- **Negative (je Paar):** Given einen veralteten Pin der **zweiten** Reihe,
  während die erste Reihe aktuell ist, when `d-check` läuft, then genau ein
  Befund `version-stale`, dessen erwartete Version aus der Quelle **des
  zweiten** Paares stammt.
- **Boundary (Kurzform ist die Ein-Paar-Liste):** Given dieselbe Konfiguration
  einmal als Kurzform und einmal als einelementige `versions.patterns`-Liste,
  when `d-check` je einmal läuft, then sind beide Befundsätze byte-identisch.
- **Boundary (Mischform):** Given Kurzform **und** `versions.patterns`
  zugleich — auch mit leerem Wert (`pin-pattern: ""`, `exempt-paths: []`,
  `patterns: []`) —, when `d-check` läuft, then Nutzungsfehler (Exit 2,
  [`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)) ohne Prüfung — keine der
  beiden Schreibweisen gewinnt still.
- **Boundary (Ventile gelten je Paar):** Given eine Datei in den
  `exempt-paths` des **ersten** Paares, die einen veralteten Pin des
  **zweiten** trägt, und zusätzlich die `current-from`-Datei des ersten
  Paares mit einem veralteten Pin des zweiten, when `d-check` läuft, then
  meldet das zweite Paar beide — Ausnahmen und Selbst-Ausnahme sind
  paar-lokal, nicht modulweit.
- **Boundary (dieselbe Zeile, derselbe Wert):** Given zwei Paare, die auf
  derselben Zeile denselben Pin-Wert treffen, when `d-check` läuft, then
  **genau ein** Befund — und seine Nachricht nennt **beide** Erwartungen mit
  ihrer je eigenen Quelle, in Deklarationsreihenfolge, auch wenn die
  Erwartungen verschieden sind. Trifft **dasselbe** Paar den Wert mehrfach auf
  der Zeile, steht seine Erwartung trotzdem einmal darin.
- **Boundary (Wortlaut bei genau einem Paar):** Given eine Ein-Paar-Konfiguration
  (Kurzform oder einelementige Liste), when ein Pin veraltet ist, then nennt
  die Nachricht `versions.current-from` wie bisher — der Befundsatz
  bestehender Konfigurationen bleibt byte-identisch
  ([`DC-QA-02`](#dc-qa-02--determinismus)).
- **Boundary (fail-closed je Paar):** Given eine leere `versions.patterns`-
  Liste, ein Paar ohne `pin-pattern` oder ein Paar, dessen `current-from`
  nicht auflöst, when `d-check` läuft, then Exit 2 ohne Prüfung — ein
  unvollständiges Paar schaltet nicht still die übrigen scharf.
- **Modul-aus:** Given **kein** `versions`-Block in der Konfiguration, when
  `d-check` läuft, then ist der Befundsatz byte-identisch zum Lauf ohne das
  Modul ([`DC-QA-02`](#dc-qa-02--determinismus)) und es wird nichts geschrieben
  ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

**Out-of-Scope:** `--repair`/Auto-Bump von `version-stale` in dieser Version
(deterministisch ableitbar, aber Folge-CR an
[`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch)); Ableiten der aktuellen
Version aus Git-Tags (außerhalb des read-only gemounteten Baums, bräche
[`DC-QA-02`](#dc-qa-02--determinismus)); **mehrere Quellen je Paar** (ein
Muster gegen zwei `current-from`-Quellen — die Gleichheits-Prüfung hätte dann
keine eindeutige Erwartung); semantische Versions-Ordnung (nur Gleichheit, keine „neuer als"-Prüfung);
Schreiben oder Anlegen des Registers durch das Werkzeug selbst (read-only).

---

### DC-FA-PIN-001 — Content-Pin gegen inhaltlichen Drift (Modul `pins`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `pins` prüft d-check
Markdown-Links (Bild-Links ausgenommen), die einen **Content-Pin** tragen — einen Inline-HTML-Kommentar
`<!-- dpin: sha256:<hex> -->`, der **unmittelbar** (nur Whitespace dazwischen)
auf der **gleichen Zeile** dem schließenden `)` des Links folgt. Der Pin bindet
an genau diesen Link; ein Marker, der keinem Link unmittelbar folgt (allein, auf
der Folgezeile oder durch Text getrennt), ist **inert** (kein Befund). Bei
mehreren Links je Zeile bindet jeder Marker an den ihm direkt vorausgehenden Link.

d-check löst das Linkziel auf, bestimmt den **Ziel-Span** — die ganze Ziel-Datei
(Link ohne Anker) oder die **Heading-Section** des Ankers (Überschrift bis zur
nächsten gleich-/höherrangigen) — **welcher Anker das ist**, entscheidet dieselbe
Erkennung wie in
[`DC-FA-ANCH-001`](#dc-fa-anch-001--heading-anker-validierung-modul-anchors)
(Heading-Slug oder Inline-HTML-Anker, beide außerhalb von Fenced- und
Inline-Code) —, normalisiert ihn **whitespace-/reflow-invariant**
(Whitespace-Folgen → ein Leerzeichen, Zeilenenden vereinheitlicht,
führende/abschließende Leerzeilen entfernt), bildet den SHA-256 und vergleicht mit
dem Pin. Abweichung → Befund `link-stale`. Der Span wird **roh** gehasht — inkl.
Fenced-Code-Inhalt (Drift in Code-Beispielen soll gefangen werden; eng begrenzte
Ausnahme von der Fence-Opazität).

**Nur auflösbare Links:** Lässt sich das Ziel des gepinnten Links nicht auflösen
(Datei fehlt, Anker ohne passende Section), erzeugt `pins` **keinen** Befund — der
strukturelle Befund ist Sache von
[`DC-FA-LINK-001`](#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)/[`DC-FA-ANCH-001`](#dc-fa-anch-001--heading-anker-validierung-modul-anchors)
(`target-missing`/`anchor-missing`), auch im `pins`-only-Lauf (kein eigener, kein
doppelter Befund).

**Opt-in pro Link:** nur Links mit Pin werden geprüft. Strikt opt-in Modul
(Default aus): ohne aktives `pins` byte-identisch
([`DC-QA-02`](#dc-qa-02--determinismus)), nichts geschrieben
([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)); der
Modul-Scope ([`DC-FA-CONF-002`](#dc-fa-conf-002--modul-lokaler-scan-scope)) gilt
wie für jedes Modul. **Diagnose-only:** `link-stale` liefert keinen
`--repair`-Hunk — Re-Pinnen ist menschliche Annahme der Drift, kein eindeutig
ableitbarer Fix; ein `--bless` wäre eine spätere, eigene Anforderung (berührt
[`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch)).

**Akzeptanzkriterien:**

- **Happy Path:** Given `pins` aktiv und `[…](ziel.md#abschnitt) <!-- dpin: sha256:<H> -->`, dessen normalisierter Ziel-Span `<H>` ergibt, when `d-check` läuft, then kein Befund, Exit 0; ein read-only gemountetes Repository genügt.
- **Boundary (Reflow):** Given korrekter Pin und eine **nur**-Whitespace/Umbruch-Änderung am Ziel-Span (Wort-Inhalt identisch), when `d-check` läuft, then **kein** Befund (die Normalisierung absorbiert Reflow).
- **Negative:** Given Pin und eine **inhaltliche** Änderung am Ziel-Span, when `d-check` läuft, then ein Befund `link-stale` (Datei:Zeile des Links, erwarteter vs. errechneter Hash gekürzt), Exit 1.
- **Boundary (Marker-Bindung):** Given `[a](a.md) <!-- dpin: sha256:<Ha> --> [b](b.md) <!-- dpin: sha256:<Hb> -->` (zwei Links, je ein unmittelbar — nur Whitespace — folgender Marker in Normsyntax) und zusätzlich einen Marker allein auf der Folgezeile, when `d-check` läuft, then prüft `pins` **beide** Links (der Pin nach `a` bindet an `a`, der nach `b` an `b`); der Folgezeilen-Marker folgt keinem Link unmittelbar und ist inert (kein Befund daraus).
- **Boundary (eine Anker-Antwort):** Given ein `dpin`-Link, dessen Fragment ein **Duplikat-Slug**, ein **prozent-kodiertes** Fragment oder ein Inline-HTML-Anker ist, when `pins` und `anchors` im selben Lauf aktiv sind, then bestimmen beide denselben Anker gleich; ein Anker, den `anchors` nicht kennt, macht das Ziel **unauflösbar** — `pins` schweigt dann (kein Ersatz-Befund), und der Drift-Schutz entfällt sichtbar über `anchor-missing` statt still.
- **Boundary (Ziel weg):** Given ein gepinnter Link mit fehlender Ziel-Datei oder fehlendem Anker, when `d-check --enable pins` (auch ohne `links`/`anchors`) läuft, then **kein** `link-stale`; mit aktivem `links`/`anchors` erscheint `target-missing`/`anchor-missing` **einmal** (kein Doppelbefund durch `pins`).
- **Modul-aus:** Given **kein** aktives `pins`, when `d-check` läuft, then ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)) und es wird nichts geschrieben ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

**Out-of-Scope:** semantische Prüfung, ob das Zitat *sinngemäß* passt
(unentscheidbar); Pinnen/Re-Pinnen durch das Werkzeug selbst (read-only; ein
`--bless`-Emissionsmodus wäre eine eigene Anforderung, berührt
[`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch)); **Absatz-Ebene** als Span
(in Markdown nicht stabil adressierbar); Pins auf Ziele **außerhalb der
Repo-Wurzel** (repo-escape, nicht hashbar — der `pins`-Scope begrenzt nur die
gescannten Quell-Dateien, nicht das Hashen repo-interner Ziele;
[`DC-QA-02`](#dc-qa-02--determinismus));
Default-on oder eine Pflicht zu pinnen; mehrere Hash-Algorithmen (nur `sha256`).

---

### DC-FA-IMM-001 — Immutabilitäts-Pin gegen Core-Drift (Modul `immutable`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `immutable` prüft d-check, ob
der **unveränderliche Core** einer Datei seit dem Setzen eines Pins inhaltlich
unverändert ist. Eine Datei trägt dazu einen Inline-HTML-Kommentar
`<!-- immutable: sha256:<hex> -->` (den **Pin**). Ist er vorhanden (der **erste**
Marker je Datei zählt; weitere sind inert), bildet d-check den **Core-Hash** der
Datei und vergleicht ihn mit dem Pin. Abweichung → Befund `core-drift`.

Der **Core** ist der Datei-Inhalt **ohne** (a) die Pin-Marker-Zeile selbst (sonst
Selbstbezug — der Marker trägt den Hash des Core) und (b) die per
`immutable.exclude-sections` benannten Abschnitte (Überschrift bis zur nächsten
gleich-/höherrangigen; exakter Heading-Vergleich wie `matrix`). Der verbleibende
Core wird **whitespace-/reflow-invariant** normalisiert (Whitespace-Folgen → ein
Leerzeichen, führende/abschließende Leerzeichen entfernt) und per SHA-256 gehasht
— dieselbe Normalisierung wie
[`DC-FA-PIN-001`](#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in).

**Zweck:** Eine *immutabel* erklärte Datei — etwa eine `Accepted`-ADR, deren
Körper laut Prozess nicht mehr inhaltlich überschrieben werden darf (nur Anhänge
unter einem ausgenommenen Abschnitt und der Status-Übergang) — erzeugt bei jeder
versehentlichen oder unausgewiesenen Körper-Änderung einen Befund. Anders als eine
git-historienbasierte Diff-Prüfung (läge außerhalb des read-only gemounteten
Baums, bräche [`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit))
arbeitet `immutable` rein auf dem gescannten Arbeitsbaum: der Pin **ist** die im
Dokument hinterlegte Wahrheit; ein bewusstes Neu-Pinnen ist die menschliche
Annahme einer legitimen Änderung (dieselbe Disziplin wie `pins`/`versions`).

**Opt-in pro Datei:** nur Dateien **mit** Pin werden geprüft (der Pin ist die
bewusste „diese Datei ist eingefroren"-Markierung). Strikt opt-in Modul (Default
aus): ohne aktives `immutable` byte-identisch
([`DC-QA-02`](#dc-qa-02--determinismus)), nichts geschrieben
([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)); der
Modul-Scope ([`DC-FA-CONF-002`](#dc-fa-conf-002--modul-lokaler-scan-scope)) gilt
wie für jedes Modul. **Diagnose-only:** `core-drift` liefert keinen
`--repair`-Hunk — Neu-Pinnen ist menschliche Annahme der Änderung, kein eindeutig
ableitbarer Fix (wie `link-stale`/`version-stale`); ein `--bless` wäre eine
spätere, eigene Anforderung (berührt
[`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch)).

**Akzeptanzkriterien:**

- **Happy Path:** Given `immutable` aktiv und eine Datei mit `<!-- immutable: sha256:<H> -->`, deren normalisierter Core `<H>` ergibt, when `d-check` läuft, then kein Befund, Exit 0; ein read-only gemountetes Repository genügt.
- **Boundary (Reflow):** Given korrekter Pin und eine **nur**-Whitespace/Umbruch-Änderung am Core (Wort-Inhalt identisch), when `d-check` läuft, then **kein** Befund (die Normalisierung absorbiert Reflow).
- **Boundary (ausgenommener Abschnitt):** Given korrekter Pin und ein Anhang **innerhalb** eines `exclude-sections`-Abschnitts (z. B. eine neue Zeile unter `## Geschichte`), when `d-check` läuft, then **kein** Befund (der ausgenommene Abschnitt zählt nicht zum Core).
- **Negative:** Given Pin und eine **inhaltliche** Änderung am Core **außerhalb** der ausgenommenen Abschnitte, when `d-check` läuft, then ein Befund `core-drift` (Datei:Zeile des Markers, erwarteter vs. errechneter Hash gekürzt), Exit 1.
- **Boundary (kein Marker):** Given eine Datei **ohne** Pin-Marker, when `d-check --enable immutable` läuft, then **kein** Befund (nur gepinnte Dateien werden geprüft).
- **Modul-aus:** Given **kein** aktives `immutable`, when `d-check` läuft, then ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)) und es wird nichts geschrieben ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

**Out-of-Scope:** eine **VCS-/git-historienbasierte** Immutabilitäts-Prüfung
(`core(BASE)` vs. `core(HEAD)` über eine Commit-Range) — sie erweiterte die
Eingabe über den gescannten read-only-Baum hinaus
([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)) und
bräuchte einen eigenen nicht-hermetischen Port (analog `external`); eine spätere
opt-in-Stufe als eigene Anforderung/ADR (die hier gewählte Pin-Form ist die
hermetische, im Arbeitsbaum entscheidbare Hälfte, festgehalten in eigener ADR);
Pinnen/Neu-Pinnen durch das Werkzeug selbst (read-only; ein `--bless`-Emissionsmodus
wäre eine eigene Anforderung, berührt
[`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch)); das Erzwingen, **welche**
Status-Übergänge an einer gepinnten Datei zulässig sind (die Pin-Form sieht nur
„Core unverändert / nicht" — ein den Core berührender Status-Wechsel verlangt ein
bewusstes Neu-Pinnen); Default-on oder eine Pflicht zu pinnen; mehrere
Hash-Algorithmen (nur `sha256`).

---

### DC-FA-VCS-001 — Git-Diff-Immutabilität des Core über eine Commit-Range (Modul `vcs`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `vcs` prüft d-check über eine
git-**Commit-Range**, ob der **unveränderliche Core** einer als immutabel
erklärten Datei zwischen zwei Ständen inhaltlich gleich geblieben ist:
`core(BASE)` ≟ `core(HEAD)`. Das ist die **volle** git-Garantie, die
[`DC-FA-IMM-001`](#dc-fa-imm-001--immutabilitäts-pin-gegen-core-drift-modul-immutable-opt-in)
(hermetischer Pin im Arbeitsbaum) bewusst offenließ — sie fängt auch eine
Körper-Änderung, die *mit* einem gefälschten Neu-Pin getarnt würde, weil sie den
Stand zweier Commits vergleicht statt einer im Dokument hinterlegten Zahl.

Dafür liest `vcs` die git-Historie über einen eigenen **VCS-Port**: es öffnet das
im read-only-Mount vorhandene `.git` über eine **reine-Go**-git-Implementierung —
**ohne git-Binary** (das distroless-Laufzeit-Image bleibt unangetastet) und
**ohne Netz**. Die Eingabe ist damit gegenüber den hermetischen Modulen
**erweitert** (das `.git` und eine Range, nicht nur der gescannte Markdown-Baum),
bleibt aber **lokal, lesend und deterministisch**: dieselbe git-Historie + dieselbe
Range ⇒ derselbe Befundsatz ([`DC-QA-02`](#dc-qa-02--determinismus)), kein
Netzzugriff und kein Schreiben ins Repository
([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)). Die Range
wird geliefert über `--range <base>..<head>` (PR-/Push-CI) oder `--staged`
(lokaler pre-commit-Hook: BASE = `HEAD`, HEAD = der staged Index).

Geprüft wird jede in der Range **geänderte** Datei, die der konfigurierten Klasse
(`vcs.paths`, Glob) entspricht und deren **BASE**-Version die
Immutabilitäts-Bedingung erfüllt (`vcs.immutable-when`, Zeilen-Regex — z. B.
`**Status:** Accepted`):

- `core(BASE)` ≠ `core(HEAD)` → Befund `core-drift-vcs` (Körper geändert).
- HEAD-Status-Zeile außerhalb des erlaubten Musters (`vcs.head-allow`) → Befund
  `core-drift-vcs` (unzulässiger Status-Übergang einer immutablen Datei).
- **gelöschte** oder **umbenannte** immutable Datei → Befund `core-drift-vcs`
  (der Pfad einer immutablen Datei ist stabil).

Der **Core** ist — wie bei
[`DC-FA-IMM-001`](#dc-fa-imm-001--immutabilitäts-pin-gegen-core-drift-modul-immutable-opt-in)
— der normalisierte Inhalt ohne die per `vcs.exclude-sections` benannten
Abschnitte; zusätzlich wird **nur die Kopf-Status-Zeile** (`vcs.status-line`,
erstes Vorkommen **vor** der ersten H2) aus dem Core entfernt, eine gleichlautende
Zeile im Körper bleibt Teil des Core (sonst rutschte ein Edit an ihr durch). Die
verbleibende Eingabe wird **whitespace-/reflow-invariant** normalisiert und per
SHA-256 gehasht — dieselbe Normalisierung wie
[`DC-FA-IMM-001`](#dc-fa-imm-001--immutabilitäts-pin-gegen-core-drift-modul-immutable-opt-in)
/ [`DC-FA-PIN-001`](#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in).

**Strikt opt-in, fail-closed, diagnose-only:** `vcs` ist nie Default-Modul (wie
`external`); ohne aktives `vcs` ist der Befundsatz byte-identisch
([`DC-QA-02`](#dc-qa-02--determinismus)) und nichts wird gelesen, was über den
Markdown-Baum hinausgeht ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
Der Modul-Scope ([`DC-FA-CONF-002`](#dc-fa-conf-002--modul-lokaler-scan-scope))
gilt wie für jedes Modul. **fail-closed:** ein fehlendes oder unlesbares `.git`,
eine fehlende oder nicht auflösbare Range (z. B. unbekannte Basis) → Fehler
(Exit 2) mit Hinweis auf stderr, keine stille Grün-Meldung. `core-drift-vcs` ist
**diagnose-only** — es liefert keinen `--repair`-Hunk (die Korrektur ist eine
menschliche Entscheidung, kein eindeutig ableitbarer Fix; vgl.
[`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch)).

**Akzeptanzkriterien:**

- **Happy Path:** Given `vcs` aktiv mit `vcs.paths` auf eine Datei-Klasse und einer Range, in der eine immutable (`vcs.immutable-when`) Datei **nur** einen Anhang innerhalb eines `vcs.exclude-sections`-Abschnitts erhält, when `d-check --enable vcs --range <base>..<head>` läuft, then kein Befund, Exit 0.
- **Boundary (Modul-aus / git-frei):** Given **kein** aktives `vcs`, when `d-check` ohne Range in einer netzlosen, read-only Umgebung läuft, then ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)), es erfolgt kein git-Zugriff über den Scan hinaus und nichts wird geschrieben ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **Negative:** Given `vcs` aktiv und eine Range, in der der Körper einer immutablen Datei **außerhalb** der ausgenommenen Abschnitte geändert (oder die Datei gelöscht/umbenannt) wird, when `d-check --enable vcs --range <base>..<head>` läuft, then ein Befund `core-drift-vcs` (Datei, Grund), Exit 1.
- **fail-closed (git-Eingabe fehlt):** Given `vcs` aktiv, aber **kein** lesbares `.git` oder eine **fehlende/unauflösbare** Range (leere Basis, fehlender `..`-Separator), when `d-check --enable vcs` läuft, then **Exit 2** mit Hinweis auf stderr — kein stilles Grün (kein Exit 0) und kein Befund-Exit (kein Exit 1).

**Grenze (benannt, nicht mechanisiert):** Die Abschnitts-Maske von
`vcs.exclude-sections` wird auf **git-Blobs** gerechnet — auf einer Eingabe, die
kein scannendes Modul je sieht. Die Maske selbst folgt der geteilten
Überschriften-Lexik; ein **unbalancierter Fence in einer Revision** verschiebt
sie dort trotzdem, und kein Wächter meldet das, weil die Wächter den
**Arbeitsbaum** scannen. Für jeden bereits existierenden Commit ist diese Sicht
rückwirkend blind. **Die gefährliche Ausprägung ist nicht das Falsch-Rot,
sondern das stille Grün:** liegt die Öffnung **innerhalb** des ausgenommenen
Abschnitts, findet die Maske keine folgende Überschrift mehr, die Ausnahme läuft
bis zum **Dateiende**, und eine reale Core-Änderung passiert das Gate mit
**Exit 0 und ohne Ausgabe**. Beobachtbar ist das **nicht** an `vcs` selbst,
sondern nur an der Arbeitsbaum-Fassung: meldet das Modul `spans` ein
`fence-unclosed` in einer Datei, die `vcs.paths` trifft, ist der stille Pfad
erreichbar geworden — das ist der Trigger, an dem diese Grenze neu zu bewerten
ist. Gemessen wurde die Klasse über den eigenen Bestand — **null**
von 152 immutablen Revisions-Blobs —, und der Vertrag benennt sie hier, statt
einen zweiten Wächter auf einer dritten Eingabe-Achse zu fordern (Begründung in
begleitender ADR).

**Out-of-Scope:** die hermetische Pin-Form — das ist
[`DC-FA-IMM-001`](#dc-fa-imm-001--immutabilitäts-pin-gegen-core-drift-modul-immutable-opt-in),
die im Arbeitsbaum entscheidbare Schwester-Hälfte (`vcs` ist die
git-historienbasierte Hälfte, beide koexistieren als Defense-in-Depth); ein
git-**Binary** im Laufzeit-Image (der Port liest `.git` rein in Go; das
distroless-Image bleibt unangetastet, in begleitender ADR festgehalten);
Forge-/Netz-API-Aufrufe (nur lokale git-Objekte); Umbenennungs-Erkennung über
Inhalts-Ähnlichkeit statt Pfad (eine immutable Datei behält ihren Pfad);
Schreiben/Neu-Pinnen durch das Werkzeug (read-only); mehrere Hash-Algorithmen
(nur `sha256`).

---

### DC-FA-COMMITS-001 — Traceability-Kennung in Commit-Messages über eine Commit-Range (Modul `commits`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `commits` prüft d-check, dass
jede geprüfte **Commit-Message** mindestens eine **Traceability-Kennung** nach
deklariertem Muster (`commits.id-patterns` — z. B. `ADR-`-, `MR-`-, `DC-`- oder
`slice-`-Kennungen) auf einer Inhalts-Zeile trägt. Das ist die maschinelle
Durchsetzung der Regel „PRs/Commits **müssen** mindestens eine `DC-*`-, `ADR-*`-,
`MR-*`- oder `slice-*`-ID nennen"
([`harness/README.md` §Traceability rules](../harness/README.md#traceability-rules)),
als Regelmodul verkörpert — die verteilbare Form derselben Prüfung, die zuvor ein
kopiertes Shell-Skript leistete.

Dafür liest `commits` die Commit-Historie über **denselben VCS-Port** wie
[`DC-FA-VCS-001`](#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
(reine-Go-git, **ohne git-Binary** → das distroless-Laufzeit-Image bleibt
unangetastet, **ohne Netz**), erweitert um das Lesen der Commit-**Messages** einer
Range. Zwei Eingabe-Quellen speisen **dieselbe** Prüfung:

- `--range <base>..<head>` (PR-/Push-CI, lokales Gate): jede **Nicht-Merge**-Commit-Message der Range.
- `--commit-msg <datei>` (lokaler commit-msg-Hook; `-` = stdin): die **einzelne** Pending-Message, noch bevor der Commit existiert.

Die Eingabe ist damit gegenüber den hermetischen Modulen **erweitert** (git-Historie
+ Range bzw. eine Message), bleibt aber **lokal, lesend und deterministisch**:
dieselbe Historie + dieselbe Range ⇒ derselbe Befundsatz
([`DC-QA-02`](#dc-qa-02--determinismus)), kein Netzzugriff und kein Schreiben ins
Repository ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

Vor der Prüfung wird jede Message **uniform bereinigt** (git-`strip`-Cleanup: alles
ab der scissors-Zeile `# … >8 …` sowie `#`-Kommentarzeilen entfällt) — damit die
Range- und die Message-Quelle **dieselbe** (kommentar-bereinigte) Message bewerten
und nicht je nach git-Cleanup-Modus divergieren; eine Kennung muss auf einer
**Inhalts**-Zeile stehen, nicht in einem Kommentar. **Ausgenommen** (kennungs-frei
erlaubt): eine Message, deren **Betreff** (erste Zeile) `commits.exempt-pattern`
matcht — die Selbstkonfiguration setzt `^(Merge |Revert )` (Merge-/Revert-Commits,
wie schon der Skript-Vorläufer).

- Message ohne Kennung und **nicht** ausgenommen → Befund `commit-untraceable`
  (`target` = Commit-Kurz-SHA bzw. `pending` für die Message-Quelle, `message` =
  der Betreff).

**Strikt opt-in, fail-closed, diagnose-only:** `commits` ist nie Default-Modul (wie
`external`/`vcs`); ohne aktives `commits` ist der Befundsatz byte-identisch
([`DC-QA-02`](#dc-qa-02--determinismus)) und nichts wird gelesen, was über den
Markdown-Baum hinausgeht ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
Der Modul-Scope ([`DC-FA-CONF-002`](#dc-fa-conf-002--modul-lokaler-scan-scope))
gilt wie für jedes Modul. **fail-closed:** eine fehlende oder nicht auflösbare
Range (leere/unauflösbare Basis, fehlender `..`-Separator), ein fehlendes oder
unlesbares `.git` im Range-Modus, oder eine nicht lesbare Message-Datei → Fehler
(Exit 2) mit Hinweis auf stderr, keine stille Grün-Meldung. `commit-untraceable`
ist **diagnose-only** — es liefert keinen `--repair`-Hunk (die Korrektur ist ein
neuer Commit bzw. ein menschliches `--amend`; vgl.
[`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch)).

**Akzeptanzkriterien:**

- **Happy Path:** Given `commits` aktiv mit gesetzten `commits.id-patterns` und einer Range, in der jede Nicht-Merge-Commit-Message eine Kennung auf einer Inhalts-Zeile trägt, when `d-check --enable commits --range <base>..<head>` läuft, then kein Befund, Exit 0.
- **Boundary (Modul-aus / git-frei):** Given **kein** aktives `commits`, when `d-check` ohne Range in einer netzlosen, read-only Umgebung läuft, then ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)), es erfolgt kein git-Zugriff über den Scan hinaus und nichts wird geschrieben ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **Negative:** Given `commits` aktiv und eine Range mit einer Nicht-Merge-Commit-Message **ohne** Kennung (und ohne Betreff-Ausnahme), when `d-check --enable commits --range <base>..<head>` läuft, then ein Befund `commit-untraceable` (Commit, Grund), Exit 1; eine `commits.exempt-pattern`-Message (Merge/Revert) in derselben Range erzeugt keinen Befund.
- **fail-closed (git-/Range-Eingabe fehlt):** Given `commits` aktiv, aber eine **fehlende/unauflösbare** Range (leere Basis, fehlender `..`-Separator) bzw. ein **unlesbares** `.git`, when `d-check --enable commits` läuft, then **Exit 2** mit Hinweis auf stderr — kein stilles Grün (kein Exit 0) und kein Befund-Exit (kein Exit 1).

**Out-of-Scope:** die semantische Prüfung, ob die genannte Kennung *fachlich passt*
(nur Existenz eines Musters, keine Auflösung gegen Lastenheft/ADR-Index — das
leisten [`DC-FA-ID-001`](#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)/[`DC-FA-MTX-001`](#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
und die RTM [`DC-FA-CLI-009`](#dc-fa-cli-009--requirements-traceability-matrix));
das Prüfen des Arbeitsbaum-/Index-**Inhalts** (das ist
[`DC-FA-VCS-001`](#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in));
Forge-/Netz-API-Aufrufe (nur lokale git-Objekte); das *Vergeben* von Kennungen
(Agenten referenzieren, erfinden nicht); das Prüfen von **Merge**-Commits (per
Range nicht aufgezählt, `--no-merges`-Parität); Schreiben durch das Werkzeug
(read-only); ein git-**Binary** im Laufzeit-Image (der Port liest `.git` rein in
Go; in begleitender ADR festgehalten).

---

### DC-FA-PLAN-001 — Planning-Lifecycle-Konsistenz (Modul `planning`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `planning` prüft d-check eine
**Planning-Lifecycle-Invariante**: die Roadmap darf in ihrem Aktiv-Status-Abschnitt
(`planning.heading`, per Konvention `## Aktuelle Welle`) genau dann den Ruhe-Marker
(`planning.marker`, per Konvention „Keine aktive Welle") tragen, wenn **kein**
`slice-*` (`planning.slice-glob`) im Slice-Verzeichnis liegt — und umgekehrt.
Formal: `hasActive == hasSlices`, sonst Befund `planning-drift`. Es ist die
maschinelle Durchsetzung der Lifecycle-/Roadmap-Kopplung aus
[`AGENTS.md` §3.3](../AGENTS.md#33-git-mv--inhaltsänderung--zwei-commits) — als
Regelmodul verkörpert, die verteilbare Form derselben Prüfung, die zuvor ein
kopiertes Shell-Skript leistete.

Die Prüfung ist **hermetisch**: sie liest nur den read-only-Arbeitsbaum — die
Datei `planning.roadmap` und das Verzeichnis-Listing ihres Verzeichnisses; **kein**
git, **kein** Netz, **kein** Schreiben. `planning` relativiert damit **weder**
[`DC-QA-02`](#dc-qa-02--determinismus) **noch**
[`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit); es ist ein
normales hermetisches Modul (wie `codepaths`), **kein** eingabe-scope-erweiterndes
wie `vcs`/`commits`.

- `hasSlices`: ≥1 Datei im Verzeichnis der `planning.roadmap` matcht `planning.slice-glob` (per `path.Match`).
- `hasActive`: der `planning.heading`-Abschnitt (kanonische H2 bis zur nächsten H2) enthält den `planning.marker` **nicht**.
- `hasActive ≠ hasSlices` → Befund `planning-drift` — beide Drift-Richtungen: Slice vorhanden + Ruhe-Marker; oder kein Slice + benannte aktive Welle.

**fail-closed (Heading-Guard):** fehlt die kanonische `planning.heading`-Überschrift
**exakt** (Tippfehler, Zusatztext, Umbenennung), **kommt sie mehrfach vor**
(mehrdeutiger Aktiv-Status) oder fehlt die `planning.roadmap`-Datei, ist der
Aktiv-Status nicht bestimmbar → Befund `planning-drift` (kein stilles Grün —
sonst gälte still „aktiv", auch wenn ein Slice vorliegt). Ein explizit gesetztes
`planning.slice-glob` muss ein gültiges `path.Match`-Muster sein (sonst Exit 2).
**Strikt opt-in,
diagnose-only:** `planning` ist nie Default-Modul; ohne aktives `planning` ist der
Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)). Leere
`planning.roadmap` ⇒ Modul inert. `planning-drift` liefert keinen `--repair`-Hunk
(die Korrektur ist eine Roadmap-/Lifecycle-Entscheidung; vgl.
[`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch)).

**Closure-Note-Struktur (zweite Fähigkeit, `planning.closure`).** Dieselbe
Lifecycle-Invariante hat eine zweite Seite: ein Slice, der den Lifecycle
**verlässt**, trägt eine Closure-Notiz. `planning` prüft deren **Struktur** —
aber **nur**, wenn `planning.closure.dir` gesetzt ist (opt-in **innerhalb** des
opt-in Moduls; ohne den Schlüssel ist die Fähigkeit inert und der Befundsatz
byte-identisch, [`DC-QA-02`](#dc-qa-02--determinismus)). Geprüft wird jede Datei
in `planning.closure.dir`, deren Basisname `planning.closure.glob` matcht — ein
**eigener** Kandidaten-Filter, dessen Default `planning.slice-glob` ist; in ihr
**genau ein** Abschnitt, dessen Überschrift auf `planning.closure.heading-pattern`
passt (RE2, Default `^#{2,3} .*Closure-Notiz`), reichend bis zur nächsten
Überschrift gleicher oder höherer Ebene — mehrere Treffer sind mehrdeutig und
werden nicht gemessen (siehe unten). Drei Struktur-Bedingungen:

- **Abschnitt vorhanden** — keine passende Überschrift ⇒ `closure-note-missing`.
- **Substanz** — im Abschnitt stehen weniger als
  `planning.closure.min-sentences` Satzende-Zeichen (`.`, `!`, `?`; Default 4)
  ⇒ `closure-note-thin`. Gezählt wird im **einmal bereinigten** Abschnitt
  (Fenced-Code entfernt, Inline-Code geleert), und ein Satzende zählt nur vor
  **Whitespace oder Zeilenende**. Beides zusammen trennt Sätze von dem, was nur
  wie eines aussieht: die Punkte in `0.55.0` und die in einem Link-Pfad sind
  keine.
  Ein **kurzer** Platzhalter fällt damit auf; ein vollständiger Vorlagen-Rumpf
  nicht — er erreicht die Schwelle, dafür gibt es die vierte Bedingung.
- **Floskel** — der **gleiche** bereinigte Abschnitts-Text enthält
  (case-insensitiv) eine der Phrasen aus `planning.closure.boilerplate` **als
  eigenständigen Wortlaut** ⇒ `closure-note-boilerplate`. Vor und hinter dem
  Treffer darf kein Wortzeichen stehen; Satzzeichen, Anführungszeichen,
  Bindestrich und Zeilenrand sind Grenzen. Als Wortzeichen gilt die
  **ASCII**-Menge (`0-9A-Za-z_`); ein Umlaut zählt damit **nicht** als
  Wortzeichen. Für die gelebten Formen ist das richtig (`„Ok“`, `Ok.`,
  `Ok—`), es heißt aber auch: eine Phrase, an die unmittelbar ein Umlaut grenzt,
  gilt als grenzständig und **trifft**. Am eigenen Bestand gemessen kommt der
  Fall nicht vor.

  **Warum nicht als Teilstring:** kurze Phrasen wären sonst unbrauchbar — und
  kurze Phrasen sind genau die, für die die Bedingung gedacht ist. Eine Notiz,
  die nur „Ok.“ sagt, ist der Kern-Fall; als Teilstring trifft `ok`
  aber auch *dokumentiert*, *Dokument* und *Protokoll*. Am eigenen Bestand
  gemessen (97 Notizen): 68 Teilstring-Treffer für `ok`, davon **einer** echt.
  Bei mehrwortigen Phrasen ändert sich nichts. Dass es
  **derselbe** Text ist, hat eine Folge: eine Floskel **in Backticks** trifft
  nicht mehr. Eine *zitierte* Floskel ist keine benutzte — eine Notiz, die
  *über* Floskeln schreibt, soll nicht daran scheitern. Die Liste ist **per Default leer**: der Vertrag
  bringt keine sprach-spezifischen Phrasen mit, das Repo deklariert seine eigenen.

**Warum ein eigener Filter.** Die beiden Fähigkeiten stellen **verschiedene**
Fragen und haben deshalb verschiedene Grundmengen: die Lifecycle-Invariante fragt
„liegt hier noch Arbeit?" und meint Slice-Dateien im Roadmap-Verzeichnis; die
Closure-Struktur fragt „ist jedes abgeschlossene Paket dokumentiert?" und kann
auch Wellen- oder Etappen-Dokumente meinen. Sie teilten sich einen Schlüssel,
solange beide zufällig dasselbe trafen; wer die eine Menge weitet, verbiegt die
andere — im Grenzfall so, dass die Roadmap-Datei selbst als Slice zählt und der
Ruhe-Marker dauerhaft falsch-rot meldet. Der Default als **Verweis** (nicht als
wiederholtes Literal) hält beide Mengen zusammen, solange niemand sie trennt:
ohne den Schlüssel ist der Befundsatz byte-identisch, und es gibt kein zweites
Muster, das gepflegt werden müsste.

**Vierte Bedingung, opt-in: unausgefüllte Platzhalter** (`closure-note-placeholder`,
Schalter `planning.closure.placeholder`, bool, Default `false`). Ein
Template-Rumpf ist syntaktisch vollständig und passiert die drei Bedingungen
oben: `Ergebnis: <ergebnis>. Belege: <belege>. Offen: <offen>.` hat einen
Abschnitt, genug Satzende-Zeichen und keine deklarierte Floskel — und bleibt
grün. Erkannt wird deshalb die **Auszeichnungs-Form** des Platzhalters: eine
öffnende Winkelklammer, der kein Wortzeichen und kein Schrägstrich vorausgeht,
deren erstes Inneres kein Whitespace, `!` oder `/` ist, bis zur nächsten
schließenden Klammer **derselben Zeile**.

**Winkelklammern sind in technischer Prosa häufig**, und ein Falsch-Positiv ist
hier teurer als ein übersehener Platzhalter — es macht das Gate unglaubwürdig.
Drei Einschränkungen halten die Erkennung eng:

1. **Inline-Code zählt nicht.** Backtick-Spans werden vor der Suche geleert.
   Meta-Syntax wird in Inline-Code **gezeigt** (`<PREFIX>`, `<a id>`,
   `<datei>`) — dort steht sie absichtlich; Platzhalter stehen in der Prosa.
   Gemessen am eigenen Bestand (96 Closure-Notizen): **12** Treffer, **alle
   zwölf** in Inline-Code, **null** außerhalb. Ohne diese Einschränkung wären
   alle zwölf Falsch-Positive.
2. **Autolinks und Adressen fallen raus** — enthält das Innere `://` oder `@`,
   ist es kein Platzhalter.
3. **HTML-Tags fallen raus** — ist das erste Token des Inneren ein bekannter
   Tag-Name, ist es Markup.

Zusätzlich muss das Innere **frei von Whitespace** sein — ein Platzhalter ist ein
Feldname, kein Satz. Das trägt die Messwert-Prosa: `p95 < 1 s` scheidet über das
erste Zeichen aus, die enge Form `<1 s und der Recall >0,9` über das Whitespace
im Inneren. Generics (`vector<float>`) scheiden aus, weil ihnen ein Wortzeichen
vorausgeht; ein Winkelklammer-**Linkziel** (`](<ziel>)`) an der öffnenden
Klammer als Vorzeichen.

**Zwei Grenzen sind benannt, nicht geschlossen.** Ein **eingerückter**
Code-Block (vier Leerzeichen) ist in d-check **nirgends** modelliert — weder
hier noch in `links`, `ids` oder `codepaths`; ein Platzhalter darin meldet.
Und eine ungerade Backtick-Zahl im Absatz verschiebt die Inline-Code-Paarung,
wie überall sonst auch. Beides sind Eigenschaften der geteilten Lexik; sie hier
allein zu reparieren hieße, eine sechste Sicht auf denselben Text zu bauen.

Gemeldet wird **der erste** Treffer je Kandidat, wie bei der Floskel — mehrere
Platzhalter derselben Notiz sind dieselbe Reparatur. Diese Bedingung liest
denselben **einmal bereinigten** Abschnitt wie die drei anderen; eine eigene,
engere Sicht führt sie nicht mehr.

Geprüft wird **ausschließlich** `planning.closure.dir` (per Konvention das
Verzeichnis der abgeschlossenen Slices) — ein Slice in Arbeit trägt noch keine
Closure-Pflicht. Als Abschnitts-Überschrift zählt nur eine **echte
ATX-Überschrift außerhalb von Fenced-Code**; eine Fließtext-Zeile, die bloß mit
`#` beginnt, eröffnet und beendet keinen Abschnitt. Kommt die Überschrift
**mehrfach** vor, ist der Abschnitt mehrdeutig ⇒ `closure-note-ambiguous`, und
es wird **nicht** gemessen: ohne eindeutigen Abschnitt sagt eine Satzzahl
nichts, und ein zweiter, stehengebliebener Abschnitt ist der typische Rest einer
Vorlage. Dieselbe Härte gilt in dieser Anforderung längst für den
Aktiv-Status-Abschnitt; die Asymmetrie war ein Versäumnis.

**Diese Fähigkeit ist ein Preset der allgemeinen Struktur-Semantik**
([`DC-FA-STRUCT-001`](#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)):
gleiche Abschnitts-Bestimmung, gleiche Fence-Behandlung, gleiche Zählung. Sie
behält ihre eigenen Grund-Codes — die sind stabil zugesagt — und bleibt damit
für bestehende Konsumenten unverändert; verändert werden darf nur beides
zugleich.

**fail-closed** in zwei Stufen. Zur **Laufzeit** (`closure-note-missing`, Exit 1):
ein gesetztes, aber fehlendes oder unlesbares `planning.closure.dir` — **und
ebenso ein Verzeichnis ohne einen einzigen Kandidaten.** Den Schlüssel zu setzen
**ist** die Behauptung, dass dort Closure-Notizen liegen; findet der Lauf keine,
liefe das Gate fortan leer und grün. Am **Config-Rand** (Exit 2): ein `dir`
außerhalb der Repo-Wurzel (absolut oder mit `..`), ein nicht kompilierendes
`heading-pattern`, ein **explizit** gesetztes `min-sentences` < 1, ein
**explizit** gesetzter leerer oder ungültiger `closure.glob` und ein leerer
`boilerplate`-Eintrag. `placeholder` kennt keinen ungültigen Wert — ein Bool ist
gesetzt oder nicht. Ein explizit leerer Glob bricht ab, statt still auf den
Default zurückzufallen: den Schlüssel zu setzen ist eine Aussage, und eine
Aussage, die nichts trifft, ist ein Konfigurationsfehler — dieselbe Begründung
wie bei `min-sentences`.
Die Fähigkeit bleibt **hermetisch** (Arbeitsbaum, kein git, kein Netz) und
**diagnose-only** — kein `--repair`-Hunk, denn eine Closure-Notiz schreibt der
Autor, nicht das Werkzeug.

Zugesagt ist **Struktur, nicht Bedeutung**: ob eine formal ausreichende Notiz
inhaltlich Substanz trägt oder eine Floskel ist, ist semantisch und bleibt einem
inferentiellen Nachlauf überlassen. d-check verspricht hier keine
Bedeutungs-Erkennung — und meldet umgekehrt nichts doppelt, was die Struktur
bereits abdeckt.

**Wellen-Invariante (dritte Fähigkeit, `planning.waves`).** Dieselbe
Lifecycle-Invariante **eine Ebene höher**: nicht der Slice, sondern die **Welle**
trägt ihren Zustand im Ort. Ein laufendes Wellendokument liegt **flach** im
Wellen-Verzeichnis, ein geschlossenes im Ruheort; eine **geplante** Welle hat
noch **keine** Datei und steht nur in der Vorschau. Damit sind vier Aussagen
maschinell entscheidbar, und die Fähigkeit prüft sie gegen die Roadmap:

| # | Aussage | Grund-Code | Reparatur |
|---|---|---|---|
| 1/2 | **`mode: one`** (Default): der Aktiv-Status-Abschnitt nennt eine Welle ⟺ unter den flachen Wellendokumenten liegt **genau eine** Wellen-Kennung (Ruhe-Marker ⟺ **keine**) · **`mode: many`**: die im Aktiv-Status-Abschnitt genannte **Kennungs-Menge** ist gleich der **Kennungs-Menge der flachen Wellendokumente** — beide Richtungen, jede Kardinalität, der Ruhe-Marker geht **nicht** ein. Vergleichsgröße sind in beiden Modi **Kennungen**, nicht Dateien: zwei flache Dokumente derselben Kennung sind **ein** Element | `wave-drift` | Roadmap nachziehen **oder** Datei verschieben |
| 3 | eine Zeile der **Vorschau** nennt eine Welle **ohne** Datei — flach wie im Ruheort | `wave-preview-exists` | Vorschau-Zeile entfernen (die Welle läuft oder ist geschlossen) |
| 4 | eine Zeile des **Abschluss-Registers** nennt eine Welle **mit Ergebnisnotiz** im Ruheort | `wave-results-missing` | Ergebnisnotiz nachtragen |
| 4′ | jede **Ergebnisnotiz** im Ruheort hat ihre Zeile im Abschluss-Register | `wave-unregistered` | Register-Zeile nachtragen |

**Das verpflichtende Artefakt einer geschlossenen Welle ist die Ergebnisnotiz,
nicht das Plan-Dokument** — und das ist gemessen, nicht gesetzt: gegen das
Plan-Dokument geprüft meldet Aussage 4 über zwei reale Planungs-Bäume **19-mal**,
weil ältere Wellen geschlossen wurden, bevor es das flache Wellendokument als
Konvention gab. Die Notiz dagegen verlangt die Closure-Prozedur und liegt über
den ganzen Bestand vor. Plan-Dokument (`planning.waves.glob`) und Ergebnisnotiz
(`planning.waves.results-glob`) sind deshalb **zwei Rollen**; die zweite Menge
wird von der ersten abgezogen.

**Die Vorschau-Aussage greift nur auf der Welle-Spalte und nur bei Kennungen.**
Eine geplante Welle trägt in der gelebten Praxis noch keine Kennung — sie bekommt
sie bei der Eröffnung —, und die **Trigger**-Spalte einer Vorschau-Zeile darf
andere Wellen nennen. Geprüft wird darum die **erste Spalte** der Tabellenzeile,
und nur, wenn dort eine Kennung steht. Verglichen wird über das Zahlen-Präfix
`welle-<n>`: die Zeile trägt die volle Kennung, die Notiz den kurzen Namen.

**Vier Reparaturen, vier Grund-Codes.** Die Trennung ist erzwungen: zwei dieser
Verletzungen können **dieselbe** Roadmap-Zeile treffen, und die
Befund-Deduplikation über (Datei, Zeile, Regel, Ziel, Grund) ließe sie sonst
zusammenfallen. Der Aktiv-Status wird **nicht** ein zweites Mal bestimmt — die
Fähigkeit liest denselben `planning.heading`-Block wie die Slice-Invariante; sind
beide verletzt, entstehen zwei Befunde mit verschiedenen Codes, keine
Widersprüche. **Opt-in innerhalb** des opt-in Moduls: ohne `planning.waves.dir`
wird kein Wellendokument geöffnet und der Befundsatz ist byte-identisch.

**Zwei Kardinalitäts-Modelle, ein Prädikat (`planning.waves.mode`).** Die erste
Aussage kennt zwei Adopter-Wirklichkeiten. Im **Ein-Wellen-Betrieb** (`one`,
Default) ist der Aktiv-Status ein Bool: entweder läuft eine Welle, oder der
Ruhe-Marker steht — verglichen wird gegen **genau eine** Wellen-Kennung unter
den flachen Dokumenten. Im
**Offene-Wellen-Modell** (`many`) ist der Aktiv-Abschnitt eine **Liste**: Er nennt
je eine Kennung pro offener Welle, mehrere sind legitim, und der Ruhe-Marker sagt
etwas **anderes** — nämlich dass gerade nichts beansprucht ist. Beides in einem
Prädikat zu führen hieße, zwei unabhängige Aussagen aneinander zu binden; der
baseline-legitime Zustand „Welle eröffnet, noch nichts beansprucht" trägt dann
Marker **und** Zeiger und wäre unter `one` unvermeidlich rot.

Unter `many` vergleicht die Fähigkeit deshalb **Mengen**, nicht Marker gegen Anzahl:
Kennungen aus dem `planning.heading`-Block gegen die Kennungen der flachen
Dokumente, in beiden Richtungen. Die **Marker-Aussage bleibt vollständig bei der
Slice-Invariante** (`planning-drift`) — dieselbe Trennung, die der Adopter in seinem
Regelwerk führt: die Liste folgt den **Dateien**, der Marker folgt dem **Anspruch**.
Der Aktiv-Status wird weiterhin **nicht** ein zweites Mal bestimmt; beide
Fähigkeiten lesen denselben Block.

**Akzeptanzkriterien:**

- **Wellen-Happy-Path:** Given `planning.waves.dir` gesetzt, eine Roadmap, deren Aktiv-Status-Abschnitt eine Welle nennt, **genau ein** flaches Wellendokument, eine Vorschau ohne Kennungen und ein Abschluss-Register, dessen Zeilen und Ergebnisnotizen sich **beidseitig** decken, when `d-check --enable planning` läuft, then **kein** Befund.
- **Wellen-Negative (vier Richtungen):** Given je eine Verletzung — Roadmap nennt eine Welle ohne flaches Dokument · eine Vorschau-Zeile nennt eine Welle, die bereits eine Datei hat · eine Abschluss-Zeile ohne Ergebnisnotiz · eine Ergebnisnotiz ohne Abschluss-Zeile —, when das Modul läuft, then je **ein** Befund mit dem zugehörigen Grund-Code (`wave-drift` · `wave-preview-exists` · `wave-results-missing` · `wave-unregistered`), Exit-Code 1.
- **Wellen-Boundary (Rollen-Trennung):** Given ein Ruheort, in dem Plan-Dokumente **und** Ergebnisnotizen demselben `waves.glob` genügen, when das Modul läuft, then zählt eine Ergebnisnotiz **nicht** als Plan-Dokument (`results-glob` wird abgezogen) — und eine Abschluss-Zeile, deren Welle nur ein Plan-Dokument, aber keine Notiz hat, meldet `wave-results-missing`.
- **Wellen-Boundary (Namens-Vorschau):** Given eine Vorschau-Zeile, deren erste Spalte einen **Namen** statt einer Kennung trägt, während ihre Trigger-Spalte eine laufende Welle nennt, when das Modul läuft, then **kein** Befund — geprüft wird die erste Spalte, nicht die Zeile.
- **Modul-aus (dritte Fähigkeit):** Given **kein** `planning.waves.dir`, when `d-check --enable planning` läuft, then wird kein Wellendokument geöffnet und der Befundsatz ist byte-identisch zum Lauf ohne die Fähigkeit ([`DC-QA-02`](#dc-qa-02--determinismus)).
- **Wellen-Happy-Path (`mode: many`):** Given `planning.waves.mode: many`, ein Aktiv-Status-Abschnitt, der **zwei** Wellen-Kennungen nennt, und **genau diese zwei** flachen Wellendokumente, when `d-check --enable planning` läuft, then **kein** Befund — unabhängig davon, ob der Ruhe-Marker im Abschnitt steht.
- **Wellen-Negative (`mode: many`, beide Richtungen):** Given `mode: many` und je eine Verletzung — der Abschnitt nennt eine Kennung ohne flaches Dokument · ein flaches Dokument ohne Kennung im Abschnitt —, when das Modul läuft, then je **ein** Befund `wave-drift`, dessen **Ziel** die betroffene Kennung nennt, Exit-Code 1.
- **Wellen-Boundary (gleiche Kennung, beide Modi):** Given **zwei** flache Wellendokumente **derselben** Kennung, when das Modul läuft, then zählen sie als **ein** Element — unter `mode: one` mit genanntem Aktiv-Status **kein** `wave-drift`, unter `mode: many` mit **einem** Zeiger dieser Kennung **kein** `wave-drift`; die Kennung, nicht die Datei, ist die Vergleichsgröße (ob ein solches Doppel-Dokument selbst meldepflichtig sein sollte, ist nicht Teil dieser Anforderung).
- **Wellen-Boundary (Marker orthogonal):** Given `mode: many`, **eine** offene Welle als Zeiger im Abschnitt, **kein** Slice im Lifecycle-Verzeichnis und den Ruhe-Marker **zusätzlich** im selben Abschnitt, when das Modul läuft, then **kein** Befund — weder `wave-drift` noch `planning-drift`. Unter `mode: one` meldet derselbe Zustand `wave-drift` (das ist die Abgrenzung der beiden Modelle, kein Fehler).
- **Wellen-Boundary (Modus-Default):** Given **kein** `planning.waves.mode`, when das Modul läuft, then ist der Befundsatz byte-identisch zu einem Lauf vor dieser Erweiterung ([`DC-QA-02`](#dc-qa-02--determinismus)).
- **fail-closed (unbekannter Modus):** Given `planning.waves.mode` auf einen anderen Wert als `one`/`many` gesetzt — einschließlich des **explizit leeren** Strings —, when `d-check` startet, then Abbruch mit Exit-Code 2 und einer Meldung, die den Schlüssel nennt (Zeiger-Disziplin wie bei `closure.glob` und den übrigen `waves`-Schlüsseln, 0.59.1).
- **Happy Path:** Given `planning` aktiv mit gesetzter `planning.roadmap`, einer Roadmap mit dem kanonischen `## Aktuelle Welle`-Abschnitt **ohne** den Ruhe-Marker und ≥1 `slice-*` im Verzeichnis, when `d-check --enable planning` läuft, then kein Befund, Exit 0 (ebenso konsistent: Ruhe-Marker vorhanden **und** kein Slice).
- **Boundary (Modul-aus):** Given **kein** aktives `planning`, when `d-check` in einer netzlosen, read-only Umgebung läuft, then ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)) und nichts wird geschrieben ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **Negative:** Given `planning` aktiv und eine Drift (Slice im Verzeichnis, aber die Roadmap trägt den Ruhe-Marker; oder kein Slice, aber die Roadmap benennt eine aktive Welle), when `d-check --enable planning` läuft, then ein Befund `planning-drift`, Exit 1.
- **fail-closed (Heading/Datei fehlt):** Given `planning` aktiv und eine Roadmap **ohne** die kanonische `planning.heading`-Überschrift (bzw. eine fehlende `planning.roadmap`-Datei), when `d-check --enable planning` läuft, then `planning-drift`, Exit 1 — kein stilles Grün.
- **Happy Path (Closure-Struktur):** Given `planning` aktiv mit gesetztem `planning.closure.dir` und Slices, deren Closure-Notiz-Abschnitt mindestens `min-sentences` Satzende-Zeichen außerhalb von Fenced-Code-Blöcken trägt, when `d-check --enable planning` läuft, then kein Befund, Exit 0.
- **Boundary (Closure-Fähigkeit aus):** Given `planning` aktiv **ohne** `planning.closure.dir`, when `d-check --enable planning` läuft, then ist der Befundsatz byte-identisch zur Aktiv-Status-Prüfung allein ([`DC-QA-02`](#dc-qa-02--determinismus)) — die Closure-Fähigkeit ist inert und liest keine Slice-Datei.
- **Negative (Platzhalter/dünn):** Given ein Slice im Closure-Verzeichnis, dessen Closure-Notiz-Abschnitt weniger als `min-sentences` Satzende-Zeichen außerhalb von Code-Blöcken trägt, when `d-check --enable planning` läuft, then ein Befund `closure-note-thin` (Datei, Zeile der Abschnitts-Überschrift), Exit 1.
- **Negative (Abschnitt fehlt):** Given ein Slice im Closure-Verzeichnis **ohne** eine auf `heading-pattern` passende Überschrift, when `d-check --enable planning` läuft, then ein Befund `closure-note-missing`, Exit 1.
- **Negative (Floskel):** Given einen konfigurierten `planning.closure.boilerplate`-Eintrag, dessen Phrase case-insensitiv und **an Wortgrenzen** im bereinigten Abschnitts-Text vorkommt, when `d-check --enable planning` läuft, then ein Befund `closure-note-boilerplate`, Exit-Code 1.
- **Boundary (Code-Block zählt nicht):** Given eine Closure-Notiz, deren Satzende-Zeichen überwiegend in einem Fenced-Code-Block stehen, when `d-check --enable planning` läuft, then zählen nur die Zeichen außerhalb des Blocks — bleibt die Zahl unter der Schwelle, dann `closure-note-thin`.
- **Boundary (eigener Kandidaten-Filter, Default):** Given `planning.closure.dir` gesetzt und **kein** `planning.closure.glob`, when `d-check --enable planning` läuft, then ist der Befundsatz byte-identisch zu einem Lauf, der `planning.slice-glob` als Kandidaten-Filter verwendet.
- **Happy Path (Filter entkoppelt):** Given `planning.closure.glob: "*.md"` bei unverändertem `planning.slice-glob`, when `d-check --enable planning` läuft, then werden **alle** Markdown-Dateien im Closure-Verzeichnis auf ihre Notiz geprüft, **ohne** dass sich die Aussage der Lifecycle-Invariante ändert.
- **fail-closed (leerer Glob):** Given `planning.closure.glob` **explizit** auf den leeren String gesetzt, when `d-check` startet, then Abbruch mit Exit-Code 2 und einer Meldung, die den Schlüssel nennt.
- **Boundary (Platzhalter-Erkennung aus):** Given `planning.closure.placeholder` nicht gesetzt oder `false`, when das Modul läuft, then ist der Befundsatz byte-identisch zu einem Lauf ohne diese Bedingung.
- **Negative (Template-Rumpf):** Given eine Closure-Notiz mit vier Platzhalter-Sätzen und aktivem Schalter, when das Modul läuft, then **genau ein** Befund `closure-note-placeholder` (erster Treffer), Exit-Code 1.
- **Boundary (technische Prosa):** Given eine Closure-Notiz mit `p95 < 1 s`, `Recall > 0,9`, einem Generic, einem Autolink, einer Mail-Adresse, einem HTML-Tag und einer in **Inline-Code** gezeigten Meta-Syntax, when das Modul mit aktivem Schalter läuft, then **kein** Befund.
- **Boundary (kombinierbar):** Given eine Notiz, die zugleich unter der Substanz-Schwelle liegt und einen Platzhalter trägt, when das Modul läuft, then **beide** Befunde.
- **Negative (Satzende-Form):** Given eine Closure-Notiz, deren Satzende-Zeichen überwiegend in Versionsnummern und Link-Pfaden stehen, when das Modul läuft, then zählen **nur** die vor Whitespace oder Zeilenende stehenden — liegt die Zahl darunter, ein Befund `closure-note-thin`.
- **Boundary (Inline-Code trägt keine Substanz):** Given eine Notiz, deren Sätze fast ausschließlich in Inline-Code stehen, when das Modul läuft, then zählen sie **nicht** mit.
- **Boundary (zitierte Floskel):** Given eine deklarierte Floskel, die im Abschnitt **in Backticks** steht, when das Modul läuft, then **kein** `closure-note-boilerplate`; dieselbe Phrase im Fließtext meldet weiter.
- **Negative (kurze Phrase, eigenständig):** Given eine deklarierte Phrase, die im Abschnitt als eigenständiges Wort steht, when das Modul läuft, then ein Befund `closure-note-boilerplate`.
- **Boundary (kurze Phrase in einem Wort):** Given dieselbe Phrase **nur** als Teil eines längeren Wortes, when das Modul läuft, then **kein** Befund.
- **fail-closed (kein Kandidat):** Given ein existierendes `planning.closure.dir`, in dem **keine** Datei den **effektiven** Kandidaten-Filter matcht (gesetzter `planning.closure.glob`, sonst `planning.slice-glob`), when `d-check --enable planning` läuft, then ein Befund `closure-note-missing` auf das Verzeichnis, Exit-Code 1.
- **fail-closed (Closure-Verzeichnis):** Given `planning.closure.dir` gesetzt, aber fehlend oder unlesbar, when `d-check --enable planning` läuft, then `closure-note-missing`, Exit 1.
- **fail-closed (Config-Rand):** Given ein `planning.closure.dir` außerhalb der Repo-Wurzel, ein nicht kompilierendes `heading-pattern`, ein explizit gesetztes `min-sentences` < 1 **oder** einen leeren `boilerplate`-Eintrag, when `d-check` startet, then Exit 2 vor dem Lauf; ein **abwesendes** `min-sentences` ist dagegen der Default (kein Fehler).
- **Boundary (Rauten-Fließtext):** Given eine Closure-Notiz, die eine Zeile wie `#1 war ein Thema` enthält, when `d-check --enable planning` läuft, then gilt diese Zeile **nicht** als Überschrift — der Abschnitt reicht über sie hinaus und eine dahinter stehende Floskel wird gefunden.
- **Preset-Kopplung:** Given dieselbe Datei, einmal über diese Fähigkeit und einmal über eine gleichwertig konfigurierte `structure`-Regel im Modus `one`, when beide laufen, then melden sie an **denselben Zeilen** — die Semantiken fallen aus einer Mechanik.
- **Negative (mehrdeutig):** Given einen Slice mit **zwei** auf `heading-pattern` passenden Überschriften, when `d-check --enable planning` läuft, then `closure-note-ambiguous` mit der Zeile der **zweiten**, Exit 1 — und **kein** `closure-note-thin`/`-boilerplate` (ohne eindeutigen Abschnitt wird nicht gemessen).

**Out-of-Scope:** eine git-/VCS-basierte Lifecycle-Prüfung (rein hermetisch, nur
Arbeitsbaum); mehr als eine Roadmap bzw. ein Slice- oder Closure-Verzeichnis pro
Lauf; die Roadmap-Prosa jenseits des Aktiv-Status-Markers; ein `--repair`-Hunk; das
Erzwingen der Lifecycle-Move-Commit-Bündelung selbst
([`MR-013`](../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)
bleibt Commit-Zeit-Disziplin). Für die Aktiv-Status-Invariante zählt weiterhin nur
die **Existenz** der Slice-Dateien, nicht ihr Inhalt — Slice-Inhalt liest
ausschließlich die opt-in Closure-Fähigkeit, und auch sie nur im
Closure-Verzeichnis und nur strukturell: die **inhaltliche** Bewertung einer
Closure-Notiz (trägt sie Lernsignal, Folge-Slice oder Architektur-Beobachtung?)
ist semantisch und ausdrücklich **nicht** zugesagt.

---

### DC-FA-STRUCT-001 — Struktur-Invarianten innerhalb eines Dokuments (Modul `structure`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `structure` prüft d-check, ob
Dokumente einer **Klasse** die Abschnitte tragen, die ihre Klasse verlangt, und
ob diese Abschnitte die zugesagten Bausteine enthalten — die Frage „ist dieses
Dokument selbst richtig gebaut?" statt „zeigt es korrekt auf andere?".

Konfiguriert wird eine **Liste** von Regeln. Jede Regel bindet eine
Dokumentklasse an einen Abschnitts-Typ und an Bedingungen in ihm.

**Kandidaten-Menge.** `files` ist ein Glob (Syntax wie
[`DC-FA-SCAN-001`](#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln)) über
**Wurzel-relative Pfade**, ausgewertet über den gesamten Baum — **unabhängig**
von `scan.roots`/`scan.ignore`: eine Regel benennt ihre Dateien selbst, wie es
die Deklarations-Module bereits tun. Kandidaten sind **nur Markdown-Dateien** (dieselbe Endungs-Menge wie der
Scanner) — `structure` prüft Abschnitts-Struktur, und die gibt es außerhalb von
Markdown nicht; ein Glob, der ausschließlich anderes trifft, hat damit null
Kandidaten. Abgezogen wird `exempt-paths` der Regel.
Trifft eine Regel **null** Dateien ⇒ `section-missing` auf dem Glob (`line` = 1):
eine Regel, die auf nichts zeigt, läuft leer und meldet nicht Erfolg. Der
Befund-`target` trägt dabei die **Regel-Identität** (`files`-Glob **und**
Abschnitts-Selektor), nicht nur den Glob — sonst fielen zwei leer laufende
Regeln über derselben Dateimenge unter der Befund-Deduplikation zu einem Befund
zusammen, und eine der beiden bliebe unsichtbar. Zwei Regeln mit **identischer**
Identität sind ein Konfigurations-Duplikat ⇒ Exit 2. Weil
`structure` seine Eingabe so selbst bestimmt, kennt es **kein** `<modul>.scope`
([`DC-FA-CONF-002`](#dc-fa-conf-002--modul-lokaler-scan-scope) gilt für
scannende Module).

**Abschnitte bestimmen.** `section` (Heading-Klartext, exakter Vergleich der
getrimmten Überschriften-Zeile) **oder** `section-pattern` (RE2 auf dieselbe
Zeile); genau eines von beiden ist Pflicht. Als Überschrift zählt nur eine
**echte ATX-Überschrift außerhalb von Fenced-Code**; ein Abschnitt reicht bis zur
nächsten Überschrift gleicher oder höherer Ebene.

Wie viele Treffer erwartet werden, sagt `sections`:

- **`one`** (Default) — genau einer. 0 ⇒ `section-missing` (`line` = 1);
  **mehr als einer** ⇒ `section-ambiguous` (`line` = zweiter Treffer) und **keine
  Messung für diese Datei**: über dem falschen Abschnitt sagt eine Zahl nichts,
  und ein zweiter, stehengebliebener Abschnitt ist der typische Rest einer
  Vorlage.
- **`each`** — **jeder** Treffer wird einzeln geprüft. Für Klassen, deren
  Abschnitte sich **wiederholen** (Anforderungen in einem Lastenheft, Einträge in
  einem Register). 0 Treffer ⇒ `section-missing`; Mehrfachtreffer sind hier der
  Normalfall und **kein** Befund. Ohne diesen Modus bliebe jede **neu
  hinzukommende** Anforderung ungeprüft — eine Regel je Abschnitt zu schreiben
  ist kein Ausweg, sondern ein stiller.

**Bedingungen im Abschnitt**, je Abschnitt geprüft, jede optional, jede
**fence-treu** (der Abschnitts-Text wird zuvor um Fenced-Code bereinigt; **zwei**
Bedingungen lesen einen anderen Text und sind unten je benannt: die
Chronologie-Bedingung die rohen Abschnitts-Zeilen, die Überschriften-Bedingung
die Überschriften selbst) — und jede mit **eigenem** Grund-Code, weil jede eine andere Reparatur
verlangt:

| Bedingung | verletzt ⇒ | Reparatur |
|---|---|---|
| `non-empty` (bool) | `section-empty` | Inhalt schreiben |
| `min-sentences` (int ≥ 1) | `section-thin` | Substanz ergänzen |
| `max-tasks` (int ≥ 0; erklärte Teilmenge über `tasks-ignore-pattern`) | `section-oversized` | zerlegen statt dehnen |
| `forbid-pattern` (RE2) | `section-forbidden` | die Wendung ersetzen |
| `require-pattern` (RE2) | `section-pattern-missing` | die zugesagte Aussage nachtragen |
| `require-all` (Marken-Liste) | `section-marker-missing` | den fehlenden Baustein ergänzen |
| `table.order` (`asc`/`desc`; Schlüsselspalte `table.order-column`, Default 1) | `section-unordered` | die Zeile chronologisch einsortieren (bzw. die zugesagte Tabelle führen) |
| — dieselbe Bedingung, untypisierbare Schlüsselzelle | `section-cell-untyped` | die Schlüsselzelle bzw. die Spaltenwahl korrigieren |
| `headings-match` (RE2; Ebene `headings-level`, Default Abschnitts-Ebene + 1) | `section-heading-mismatch` | die Überschrift in die zugesagte Form bringen |
| `table.column[].name` (Kopfzeilen-Name) + `cell-max-chars` (int ≥ 1) | `section-cell-oversized` | die Zelle kürzen — bzw. den Inhalt dorthin bringen, wo er hingehört |
| — dieselbe Bedingung, Untergrenze `cell-min-chars` | `section-cell-undersized` | die Zelle ausfüllen |
| — dieselbe Bedingung, Spalte nicht adressierbar | `section-column-missing` | den Spaltennamen bzw. die Kopfzeile korrigieren |


**Eine Regel darf ihre Grundmenge erklären — zweimal, und die beiden Muster
sehen VERSCHIEDENE Zeichenketten.** Beide **verkleinern** nur, beide sind
optional, beide tragen **keinen** eigenen Grund-Code.

- `tasks-ignore-pattern` (RE2) nimmt Task-Items aus der `max-tasks`-Zählung
  heraus. Geprüft wird der **Item-Text** hinter Listen-Marker und Checkbox,
  nicht die rohe Zeile: gegen die rohe Zeile bezeichnete `^` immer den
  Listen-Marker, und die **verankerte** Muster-Form wäre unschreibbar. Sie ist
  die tragfähige, aber nur zusammen mit der zweiten Hälfte — gemessen an 444
  Items eines gewachsenen Bestands: ein freies Muster traf 13 Items, davon
  **zwei** echte Liefer-Zusagen; **dasselbe** Muster verankert traf **eines**,
  und auch das war eine. Erst ein Muster, das **verankert ist und auf Text
  zielt, den die Bereinigung übrig lässt**, traf 26 Items ohne einen falschen.
  Der Schlüssel ohne `max-tasks` ⇒ Exit 2.
- `exempt-section-pattern` (RE2) nimmt **Abschnitte** aus dieser Regel heraus —
  das Geschwister von `exempt-paths` eine Granularitätsstufe tiefer, für
  Bestände, die **innerhalb einer Datei** leben. Geprüft wird **dieselbe**
  Zeichenkette, die `section-pattern` sieht: die getrimmte
  Überschriften-Zeile **einschließlich** der `#`-Folge. Zwei RE2 in einer Regel
  mit zwei verschiedenen Zielen wären eine Falle — ein Muster, analog zum
  Nachbarn geschrieben, träfe still nichts. Ein **gesetztes** Muster gehört zur
  **Regel-Identität**: zwei Regeln über denselben Selektor, von denen eine
  einen Bestand ausnimmt, sind verschiedene Zusagen — und genau diese Paarung
  (Bedingung A für alle Abschnitte, Bedingung B nur für die übrigen) ist die
  Form, die Grandfathering braucht; ohne den Zusatz wäre sie ein
  Konfigurations-Duplikat ⇒ Exit 2 und damit unschreibbar. Ein **leeres**
  Muster geht nicht ein, sonst änderte sich das `target` jeder Bestandsregel.

Es läuft **vor** der Kardinalitäts-Prüfung: was ausgenommen ist, kann
`sections: one` nicht mehrdeutig machen. **Leert es die Menge ⇒
`section-missing`**, und die Meldung nennt den Schlüssel, der es tat — dieselbe
Nullmengen-Härte wie bei `exempt-paths`; ohne sie schaltete ein zu breites
Muster die Regel still ab.

**Die Überdeckung ist sichtbar:** ist `tasks-ignore-pattern` gesetzt, nennt die
`section-oversized`-Meldung die Zahl der ignorierten Items — *„Abschnitt trägt
4 Task-Items (3 ignoriert), erlaubt sind 3"* —, **auch bei null**: ein Muster,
das nichts trifft, ist eine Zusage, die nicht wirkt. Ohne den Schlüssel bleibt
die Meldung unverändert. **Zwei Grenzen, benannt statt zugesagt:** die
Sichtbarkeit greift, **solange die Regel meldet** — wer so breit ignoriert,
dass die Schwelle nie fällt, sieht nichts; und ein **`hint`** an derselben
Regel **ersetzt** die Meldung samt dieser Zahl (er gewinnt gegen jede
modul-eigene Meldung einer verletzten Bedingung). Der Nullmengen-Befund der
Abschnitts-Ausnahme ist davon **ausgenommen** und behält seinen Text: dort hat
die Regel nicht gemessen.

**Fence- und Inline-Code-Treue: für die beiden Muster verschieden, und der
Grund ist ihr Vergleichsgegenstand.** `tasks-ignore-pattern` liest den
bereinigten Abschnitts-Text und ist damit **fence- und inline-code-treu** — ein
Muster, das auf einen Ausdruck in **Backticks** zielt, trifft Leerzeichen
(gemessen: `make gates` trifft null von 444 Items, weil die Wendung
durchgängig in Inline-Code steht). `exempt-section-pattern` ist **fence-treu**
(eine Überschrift im Fenced-Block ist keine), **sieht aber Inline-Code** — weil
`section-pattern`, dessen Zeichenkette es teilt, ihn ebenfalls sieht. Eine
andere Wahl hier hieße, dass zwei Muster derselben Regel dieselbe Zeile
verschieden lesen.


**Die Regel sagt, welche Zusage sie hütet (`hint`).** Der Grund-Code sagt die
**Art** des Defekts — *„die Wendung ersetzen"* gilt für jede `forbid`-Regel
gleich. Welche Zusage **diese** Regel hütet und was der Leser tun soll, kann
nur die Regel selbst sagen: `hint` ist freier, vom Konfigurations-Autor
**verfasster** Text, der die **Erläuterung** des Befunds
(das Befund-Feld `message`) schreibt und damit
in der vierten Spalte der Befund-Zeile und in `--doctor` erscheint
([`DC-FA-CLI-004`](#dc-fa-cli-004--ausgabeformate),
[`DC-FA-CLI-007`](#dc-fa-cli-007--diagnose-modus)). Er **gewinnt** gegen die
modul-eigene Meldung. Ein **explizit leerer** Wert ⇒ Exit 2: ein leerer Hinweis
sagt nichts zu.


**Drei Befunde sind ausgenommen, und der Grund ist nicht technisch:** der
unlesbare Dateibaum, die **leer laufende** Regel und die unlesbare
**Einzeldatei** verletzen **keine Bedingung** — dort hat die Regel gar nicht
gemessen, und ein Hinweis auf die gehütete Zusage verdrängte die
fail-closed-Ursache. Sie behalten ihre eigene Meldung.
**Der Hinweis wird zur vierten Spalte einer tab-getrennten Zeile**
([`DC-FA-CLI-004`](#dc-fa-cli-004--ausgabeformate)). Ein `hint` mit **Tab oder
Zeilenumbruch** wird deshalb abgewiesen (Exit 2) statt still begradigt: der
Autor meinte etwas anderes als das Ergebnis.

**Abgegrenzt gegen den Fix-Kandidaten**
([`DC-FA-CLI-007`](#dc-fa-cli-007--diagnose-modus)): der ist **abgeleitet**,
nur dort vorhanden, wo er eindeutig ableitbar ist, und wird von
[`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch) zu einem **anwendbaren**
Patch. Ein `hint` wird nie angewendet.

Ein Sammel-Code schiede aus: die Befund-Deduplikation vergleicht (Datei, Zeile,
Regel, Ziel, Grund) — zwei verletzte Bedingungen desselben Abschnitts fielen
darunter zu **einem** Befund zusammen, und die Meldung müsste die Unterscheidung
tragen, die laut §4 der Spezifikation gerade **nicht** stabil zugesagt ist.

**Jede Überschrift des Abschnitts, positiv geprüft.** `headings-match` sagt:
*jede* Überschrift **innerhalb** des Abschnitts genügt diesem Muster. Geprüft
wird der **Überschriften-Text** (ohne die `#`-Folge, getrimmt) jeder
Überschrift, deren Ebene `headings-level` nennt — Default ist die **Ebene des
Abschnitts + 1**, also seine unmittelbaren Unterabschnitte. Verletzt eine
Überschrift das Muster ⇒ ein Befund **je Überschrift**, auf **ihrer** Zeile
(nicht auf der des Abschnitts): die Reparatur ist dort, wo die Überschrift
steht.

Diese Bedingung ist — neben der Chronologie-Bedingung — die **zweite**, die
nicht den bereinigten Abschnitts-Text liest: sie liest die Überschriften
selbst, mit **derselben** Erkennung, die den Abschnitt findet ([`DC-FA-ANCH-001`](#dc-fa-anch-001--heading-anker-validierung-modul-anchors)-Lexik:
beliebig viel führender Weißraum, Leerzeichen **oder** Tab als Trenner,
außerhalb von Fenced-Code). Der Anlass ist gemessen: dieselbe Aussage war
vorher nur als **Negation** über den Abschnitts-Text formulierbar (RE2 kennt
keinen Lookahead), und diese Ersatz-Konstruktion sprach die Lexik des Moduls
**nicht** — eine eingerückte Sektion entkam still. Eine Bedingung, die ihre
eigene Heading-Lexik nachbaut, behauptet eine Deckung, die sie nicht hat.

**Grenzen, benannt statt zugesagt:** trägt ein Abschnitt **keine** Überschrift
der geprüften Ebene, ist die Bedingung wirkungslos (kein Befund) — deshalb
gehört zu ihrer Einführung dieselbe Messung wie zu jedem Sensor: vorher rot,
nachher grün. Und ein `headings-level` **flacher als der Abschnitt** kann
innerhalb des Abschnitts nicht vorkommen (der Abschnitt endet dort); die
Bedingung ist dann ebenfalls wirkungslos.

**Marken sind Auszeichnungs-Marken, nicht Wörter.** Eine Marke `M` aus
`require-all` gilt als vorhanden, wenn eine Zeile des bereinigten Textes — nach
optionalem **Listen-Marker** und Whitespace — mit einem **hervorgehobenen**
Textlauf beginnt, dessen Inhalt mit `M` anfängt und dort endet oder mit einem
nicht-alphanumerischen Zeichen weitergeht. Gemessen an den beiden Repos, die den
Antrag tragen, deckt genau das die gelebten Formen ab: `- **Boundary:**`
(Listen-Item), `**Boundary:**` (bare) und `- **Boundary (Modul-aus):**`
(qualifiziert). Ein `Boundary` im Fließtext genügt **nicht**: die Zusage ist eine
Gliederungs-Marke.

`require-pattern` ist das Spiegelbild von `forbid-pattern` und deckt den Fall,
den Marken **nicht** decken: eine zugesagte Aussage, die **innerhalb** einer
Auszeichnung steht (gemessen: 37 von 61 Lerneintrag-Marken tragen die geforderte
Form-Angabe im Inneren des Textlaufs, nicht an seinem Anfang). Eine
Marken-Alternative („eine von dreien") wird dadurch überflüssig und ist
**nicht** Teil der Zusage.

**Chronologie-Monotonie (`table.order`).** Eine chronologische Tabelle kippt
still ihre Richtung — wer eine Zeile anfügt, schaut auf die Zeile daneben statt
auf die Regel, und danach führt dieselbe Tabelle zwei gegenläufige Blöcke; kein
anderes Gate liest Reihenfolge. Ist `table.order` gesetzt (`asc` oder `desc` —
die Richtung ist **Regel**-Konfiguration, weil legitime Bestände beide führen),
prüft d-check jede **zusammenhängende** Tabelle des Abschnitts: die Zellen der
Schlüsselspalte (`table.order-column`, 1-basiert, Default 1) werden **typisiert** und
über benachbarte Datenzeilen **nicht-strikt** monoton verglichen (gleiche
Schlüssel erlaubt — mehrere Releases an einem Tag sind der Normalfall). Kopf-
und Trennzeile deklarieren keine Daten. Jede Zeile, die die Richtung bricht ⇒
`section-unordered` an dieser Zeile; **null Datenzeilen im Abschnitt** ⇒
derselbe Code als Leerlauf-Befund an der Abschnitts-Überschrift (die Bedingung
zu setzen **ist** die Behauptung, dass hier eine chronologische Tabelle steht).

Die Typ-Menge ist **geschlossen**: **ISO-Datum** (`JJJJ-MM-TT`) und
**Punkt-Version** (optionales `v`-Präfix, mindestens zwei numerische Segmente,
segmentweise numerisch verglichen — `1.10` > `1.9`, `0.60.2` > `0.60.0`; bei
gleichem Präfix ist die kürzere Segmentfolge kleiner). Der Typ ist Pflicht,
kein Komfort: ein zeichenweiser Vergleich meldet gemessen drei **korrekt**
sortierte Bestandstabellen rot. Getypt wird der **erste** Treffer beider Muster
in der **rohen** Zelle — die Bedingung liest die rohen
Abschnitts-Zeilen, weil die Bereinigung Inline-Code leert und reale
Schlüsselspalten (Release-Register) genau dort stehen; das ist die **erste** von zwei
ausdrücklich benannten Ausnahmen (die zweite ist die Überschriften-Bedingung). Die **Tabellenzeilen-Auswahl** bleibt
fence-treu (eine Tabelle in einem Fenced-Code-Block deklariert nichts). Eine
Zelle ohne typisierbaren Schlüssel, eine Zeile mit zu wenigen Zellen oder eine
**Typ-Mischung** innerhalb der Spalte ⇒ `section-cell-untyped` an dieser Zeile
— kein stilles Überspringen, sonst schaltete ein Tippfehler die Prüfung der
restlichen Tabelle wortlos ab; hinter der gemeldeten Zelle setzt der Vergleich
beim nächsten typisierbaren Nachbar-Paar wieder auf.

**Zellenlänge (`table.column`).** Eine Tabellenspalte, die
einen **Titel** trägt, nimmt still die ganze Entscheidung auf: wer eine Zeile
anfügt, schreibt so viel hinein, wie er gerade weiß, und niemand misst nach.
Aus einer derivativen Sicht wird so eine **zweite Quelle**, die von der
eigentlichen abdriften kann. Nennt `table.column` eine Spalte, prüft d-check jede
**zusammenhängende** Tabelle des Abschnitts: die Zellen dieser Spalte tragen
höchstens `cell-max-chars` **Zeichen** ⇒ sonst `section-cell-oversized` an
**dieser Zeile**, nicht am Abschnittskopf — dort ist die Reparatur.

**`table.column` ist eine Liste, und jeder Eintrag ist eine eigene Zusage.**
Mehrere Spalten desselben Abschnitts stehen damit unter **einem** Selektor —
sie kosten keine zweite Regel, die `files` und `section` wortgleich
wiederholt. Jeder Eintrag wird unabhängig ausgewertet, und das Befund-Ziel
trägt die Spalte; ohne diesen Zusatz fielen zwei zu lange Zellen **derselben
Zeile** unter die Deduplikation zusammen. Dieselbe Spalte zweimal zu nennen ist
deshalb ein Config-Fehler, kein doppeltes Gate.

Die Spalte wird über ihren **Kopfzeilen-Namen** adressiert, nicht über eine
Position: eine eingefügte Spalte verschöbe eine Positions-Angabe **still** auf
die falsche Spalte, während ein umbenannter Kopf **laut** meldet. Adressierbar
heißt: **genau eine** Kopfzelle trägt den Namen. Trägt ihn keine, ist die
Tabelle für diese Bedingung irrelevant (eine Nebentabelle im selben Abschnitt
schaltet die Messung nicht ab); trägt ihn eine Kopfzeile **mehrfach**, oder
bindet **keine** Tabelle des Abschnitts die Spalte, oder reicht eine Datenzeile
nicht bis zu ihr ⇒ `section-column-missing`. Der Leerlauf-Fall ist dieselbe
Doppel-Rolle wie bei der Chronologie: die Bedingung zu setzen **ist** die
Behauptung, dass es diese Spalte gibt.

**Eine Obergrenze allein lässt die leere Zelle durch.** Null Zeichen liegen
unter jeder Schwelle; eine Spalte, die nur nach oben begrenzt ist, darf leer
sein. Wer zusagt, dass sie **gefüllt** ist, sagt das mit `cell-min-chars` ⇒
sonst `section-cell-undersized`. Die beiden Grenzen tragen **eigene**
Grund-Codes, weil die Reparatur eine andere ist: ausfüllen statt kürzen.
Mindestens eine der beiden ist Pflicht, sobald ein Eintrag eine Spalte nennt.

**Gezählt werden Zeichen, nicht Bytes** — die Schwelle beschreibt einen Text,
und ein Umlaut ist ein Zeichen. Die Zell-Zerlegung ist **escape- und
backtick-bewusst**: `\|` ist ein Zeichen der Zelle und **kein** Zelltrenner. Das
ist keine Bequemlichkeit, sondern die Bedingung dafür, dass die Zusage trägt —
ein naiver Split zerteilte genau die Zellen, die ein Autor mit einer Pipe
füllt, und ließe die zu lange Zelle **unbemerkt** durch. **Grenze, benannt
statt zugesagt:** gemessen wird die Zelle, **wie sie dasteht**, Markdown-Syntax
eingeschlossen — eine Zelle aus einem einzigen langen Link ist lang, auch wenn
ihr sichtbarer Text kurz ist.

**Verhältnis zur Closure-Note-Struktur.** Die zweite Fähigkeit des Moduls
`planning` ([`DC-FA-PLAN-001`](#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in))
ist ein **Preset** dieser Semantik im Modus `one`: gleiche Abschnitts-Bestimmung,
gleiche Bereinigung, gleiche Zählung. Sie behält ihre eigenen Grund-Codes —
Grund-Codes sind stabil zugesagt — und bleibt für bestehende Konsumenten
unverändert. Doppelt ist die Konfigurations-Oberfläche, nicht die Mechanik; ein
Akzeptanzkriterium hält beide zusammen.

**Strikt opt-in, hermetisch, diagnose-only:** ohne aktives `structure` bzw. ohne
Regeln ist der Befundsatz byte-identisch
([`DC-QA-02`](#dc-qa-02--determinismus)); nur lesend, netzlos
([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)); kein
`--repair`-Hunk. **fail-closed** (Exit 2, vor dem Lauf): Regel ohne `files`;
ungültiges Glob in `files`/`exempt-paths`; weder `section` noch
`section-pattern` **oder** beide; unbekannter `sections`-Wert; nicht
kompilierendes `section-pattern`/`forbid-pattern`/`require-pattern`/`tasks-ignore-pattern`/`exempt-section-pattern`;
ein `tasks-ignore-pattern` **ohne** `max-tasks` (dieselbe halbe Aktivierung wie
`table.order-column` ohne `table.order`); explizit
gesetztes `min-sentences` < 1 oder `max-tasks` < 0; leerer Eintrag in
`require-all`; **explizit leerer** `hint` oder einer, der nur Whitespace trägt;
ein `hint` mit **Tab oder Zeilenumbruch** (die Befund-Zeile ist tab-getrennt und
einzeilig); `table.order` außerhalb `asc`/`desc`; explizit gesetztes
`table.order-column` < 1; `table.order-column` ohne `table.order` (eine halbe Aktivierung
ist ein Config-Fehler, kein Zustand); ein `table` **ohne** `order` und ohne
`column`-Eintrag (eine leere Klammer sagt nichts zu); ein `table.column[].name`
aus lauter Weißraum; **derselbe Name zweimal** in einer Spalten-Liste (beide
Einträge trügen dasselbe Befund-Ziel);
explizit gesetztes `cell-max-chars`/`cell-min-chars` < 1; eine Untergrenze
**über** der Obergrenze (keine Zelle erfüllt beides); und ein Spalten-Eintrag
**ohne** jede Grenze — hier gibt es, anders als
bei `table.order-column`, auf **keiner** Seite einen Default, den man still annehmen
könnte. **Ebenfalls Exit 2, und mit dem neuen Ort in der Meldung:** die fünf
flachen Vorgänger-Schlüssel auf Regel-Ebene (`table-order`, `table-column`,
`cell-max-column`, `cell-max-chars`, `cell-min-chars`) — eine Zusage, die
dasteht und nicht wirkt, wäre der schlimmste Ausgang.
**`exempt-paths` hebelt den Leerlauf-Befund nicht aus:** bleiben
nach Abzug null Kandidaten, ist das derselbe `section-missing` — sonst schaltete
ein Ventil die Regel still ab.

**Akzeptanzkriterien:**

- **Happy Path:** Given `structure` aktiv mit einer Regel im Default-Modus, deren Dokumente den geforderten Abschnitt samt erfüllter Bedingungen tragen, when `d-check --enable structure` läuft, then kein Befund, Exit 0.
- **Boundary (Modul-aus):** Given **kein** aktives `structure` oder eine leere Regel-Liste, when `d-check` läuft, then ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)) und keine Datei wird zusätzlich gelesen.
- **`sections: each`:** Given eine Datei mit **drei** passenden Abschnitten, von denen einer eine Bedingung verletzt, when `d-check --enable structure` läuft, then **genau ein** Befund — für den verletzenden Abschnitt, mit dessen Zeile — und **kein** `section-ambiguous`.
- **`sections: one` (Default):** Given eine Datei mit **zwei** passenden Abschnitten, when `d-check --enable structure` läuft, then `section-ambiguous` mit der Zeile des **zweiten**, Exit 1 — und **kein** Bedingungs-Befund für diese Datei.
- **Negative (Abschnitt fehlt):** Given ein Dokument der Klasse **ohne** passende Überschrift, when `d-check --enable structure` läuft, then `section-missing` (`line` = 1), Exit 1.
- **Negative (je Bedingung ein Code):** Given einen Abschnitt, der **zwei** Bedingungen zugleich verletzt, when `d-check --enable structure` läuft, then **zwei** Befunde mit **verschiedenen** Grund-Codes — die Deduplikation fasst sie nicht zusammen.
- **Boundary (fence-treu):** Given einen Abschnitt, dessen Sätze, Task-Items oder Marken **ausschließlich** in einem Fenced-Code-Block stehen, when `d-check --enable structure` läuft, then zählen sie nicht — die Bedingung gilt als verletzt.
- **Teilmenge (`tasks-ignore-pattern`):** Given einen Abschnitt mit sieben Task-Items, von denen vier dem Muster genügen, und `max-tasks: 3`, when `d-check --enable structure` läuft, then kein Befund; **ohne** den Schlüssel `section-oversized`, und die Meldung ist dann **byte-identisch** zu der vor dieser Fähigkeit. Ein fünfter, **nicht** getroffener Punkt ⇒ wieder `section-oversized`, dessen Meldung die Zahl der ignorierten nennt — auch wenn sie null ist.
- **Teilmenge (Muster-Ziel):** Given dieselben Items, when das Muster am **Item-Text** verankert ist (`^…`), then trifft es; when es am **Listen-Marker** verankert ist (`^- \[ \]…`), then trifft es **nicht** — geprüft wird der Text hinter der Checkbox. Given ein Item, dessen Wendung in **Inline-Code** steht, when ein Muster auf diese Wendung zielt, then trifft es **nicht** (die Zählung liest den bereinigten Text).
- **Teilmenge (`exempt-section-pattern`):** Given eine Datei mit drei gleichartigen Abschnitten, von denen zwei dem Muster genügen, `sections: each` und eine verletzte Bedingung in allen dreien, when `d-check --enable structure` läuft, then **genau ein** Befund; **ohne** den Schlüssel drei. Given ein Muster **ohne** die `#`-Folge, when es läuft, then trifft es **keinen** Abschnitt — geprüft wird dieselbe Zeile wie bei `section-pattern`.
- **Teilmenge (Nullmenge und Kardinalität):** Given ein `exempt-section-pattern`, das **alle** Treffer nimmt, when der Lauf endet, then `section-missing` (`line` = 1) und die Meldung nennt `exempt-section-pattern` samt Zahl — kein stilles Grün. Given zwei Treffer bei `sections: one`, von denen einer ausgenommen ist, when der Lauf endet, then **kein** `section-ambiguous`: die Ausnahme läuft vor der Kardinalitäts-Prüfung.
- **Marken-Formen:** Given einen Abschnitt, der eine `require-all`-Marke als Listen-Item (`- **M:**`), bare (`**M:**`) oder qualifiziert (`- **M (Zusatz):**`) trägt, when `d-check --enable structure` läuft, then gilt sie in **allen drei** Formen als vorhanden; trägt er sie nur als Wort im Fließtext, dann `section-marker-missing`.
- **`require-pattern`:** Given einen Abschnitt, dessen zugesagte Aussage **innerhalb** einer Auszeichnung steht, when eine `require-pattern`-Regel sie beschreibt, then kein Befund; fehlt sie, dann `section-pattern-missing`.
- **Chronologie (Happy):** Given eine Regel mit `table.order: desc` und einen Abschnitt, dessen Tabelle in der Schlüsselspalte absteigend sortiert ist — gleiche Schlüssel in Folge eingeschlossen —, when `d-check --enable structure` läuft, then kein Befund.
- **Chronologie (Bruch):** Given dieselbe Regel und eine Datenzeile, deren Schlüssel größer als der ihrer Vorgänger-Zeile ist, when der Lauf endet, then `section-unordered` mit der Zeile der **brechenden** Datenzeile, Exit 1.
- **Chronologie (Typ-Pflicht):** Given eine absteigend sortierte Versions-Spalte, in der `0.10.0` über `0.9.0` steht, when der Lauf endet, then **kein** Befund — segmentweise numerisch ist `0.10.0` größer; ein zeichenweiser Vergleich meldete rot.
- **Chronologie (Roh-Lesung):** Given eine Schlüsselspalte, deren Zellen in Inline-Code stehen — eine davon zusätzlich mit HTML-Anker —, when der Lauf endet, then wird jede Zelle getypt und die Ordnung geprüft; auf dem bereinigten Abschnitts-Text wäre jede Zelle leer.
- **Chronologie (untypisierbar):** Given eine Datenzeile, deren Schlüsselzelle kein Datums-/Versions-Token trägt (oder zu wenige Zellen hat, oder deren Typ vom typisierbaren **Vorgänger** abweicht), when der Lauf endet, then `section-cell-untyped` an **genau dieser** Zeile — der Vergleichs-Anker wird nach jedem Befund zurückgesetzt: eine gesunde Folge-Zeile hinter einer Misch-Zelle meldet **nicht**, und benachbarte typisierbare Paare werden weiterhin verglichen.
- **Chronologie (Kopf/Fence/Leerlauf):** Given eine Tabelle, deren Kopfzeile ein Wort in der Schlüsselspalte trägt, when der Lauf endet, then zählen Kopf- und Trennzeile nicht als Datenzeilen; Given ein Abschnitt, dessen einzige Tabellenzeilen in einem Fenced-Code-Block stehen oder der gar keine trägt, then `section-unordered` als Leerlauf-Befund an der Abschnitts-Überschrift.
- **Ventil:** Given ein Dokument, dessen Pfad einem `exempt-paths`-Glob der Regel entspricht, when `d-check --enable structure` läuft, then kein Befund für dieses Dokument.
- **Preset-Kopplung:** Given dieselbe Datei, einmal über eine `structure`-Regel und einmal über die Closure-Fähigkeit mit gleichwertiger Konfiguration, when beide laufen, then melden sie an **denselben Zeilen** — die Semantiken fallen aus einer Mechanik.
- **fail-closed (Leerlauf):** Given eine Regel, deren `files`-Glob keine Datei trifft — auch nach Abzug von `exempt-paths` —, when `d-check --enable structure` läuft, then `section-missing` auf dem Glob, Exit 1.
- **fail-closed (Config-Rand):** Given eine Regel mit einem der oben aufgezählten Config-Fehler, when `d-check` startet, then Exit 2 vor dem Lauf.
- **Überschriften-Muster (Happy Path):** Given eine Regel mit `headings-match` und einen Abschnitt, dessen sämtliche Unterabschnitte dem Muster genügen, when `d-check --enable structure` läuft, then kein Befund, Exit 0.
- **Überschriften-Muster (Negative, je Überschrift):** Given einen Abschnitt mit **zwei** verletzenden Unterabschnitten, when der Lauf endet, then **zwei** Befunde `section-heading-mismatch`, jeder auf der **Zeile seiner** Überschrift, jeder mit dem Überschriften-Text in der Meldung.
- **Überschriften-Muster (Lexik des Moduls):** Given verletzende Überschriften in den Formen, an denen eine nachgebaute Lexik scheitert — **eingerückt**, mit **Tab** statt Leerzeichen getrennt, und eine gleichlautende Zeile **innerhalb** eines Fenced-Blocks —, when der Lauf endet, then melden die ersten beiden, die dritte nicht.
- **Überschriften-Muster (Ebene):** Given einen Abschnitt der Ebene 2 mit einer verletzenden Überschrift der Ebene 3 **und** einer der Ebene 4, when die Regel keine `headings-level` nennt, then meldet nur die der Ebene 3 (Default = Abschnitts-Ebene + 1); mit `headings-level: 4` nur die der Ebene 4.
- **Überschriften-Muster (wirkungslos, benannt):** Given einen Abschnitt **ohne** Überschrift der geprüften Ebene, when der Lauf endet, then kein Befund — die Bedingung ist dann vacuously wahr, und genau deshalb wird ihre Einführung vorher rot gemessen.
- **Überschriften-Muster (Modul-aus):** Given eine Regel **ohne** `headings-match`, when `d-check` läuft, then ist der Befundsatz byte-identisch zum Lauf ohne den Schlüssel ([`DC-QA-02`](#dc-qa-02--determinismus)); ein nicht kompilierendes Muster ⇒ Exit 2, ein `headings-level` außerhalb 1–6 ⇒ Exit 2.
- **Zellenlänge (Happy Path):** Given eine Regel mit einem `table.column`-Eintrag (`name` plus `cell-max-chars`) und einen Abschnitt, dessen Zellen dieser Spalte unter der Schwelle bleiben — während eine **Nachbar**-Spalte sie deutlich überschreitet —, when `d-check --enable structure` läuft, then kein Befund, Exit 0: die Bedingung ist **spalten**-gebunden, nicht zeilen-pauschal.
- **Zellenlänge (Negative, auf ihrer Zeile):** Given dieselbe Regel und eine Datenzeile, deren Zelle die Schwelle überschreitet, when der Lauf endet, then `section-cell-oversized` mit der Zeile **dieser Datenzeile** (nicht der Abschnitts-Überschrift), und die Meldung nennt Ist- und Soll-Zahl.
- **Zellenlänge (Grenzfall):** Given eine Zelle mit **genau** `cell-max-chars` Zeichen, when der Lauf endet, then kein Befund; Given eine mit **einem** Zeichen mehr, then genau ein Befund — die Schwelle ist die größte zulässige Länge.
- **Zellenlänge (Zeichen, nicht Bytes):** Given eine Zelle aus 20 Umlauten (40 Byte) bei `cell-max-chars: 20`, when der Lauf endet, then **kein** Befund — eine Byte-Zählung meldete sie rot.
- **Zellenlänge (escapte Pipe):** Given eine zu lange Zelle, die eine escapte Pipe (`\|`) enthält, when der Lauf endet, then wird sie als **eine** Zelle gemessen und gemeldet, und die escapte Pipe zählt als **ein** Zeichen — ein naiver Split an `|` zerteilte sie und ließe sie unbemerkt durch.
- **Zellenlänge (Spalte nicht adressierbar):** Given einen Abschnitt, dessen Tabellen den Namen **nicht** tragen, when der Lauf endet, then `section-column-missing` an der Abschnitts-Überschrift; Given eine Kopfzeile, die ihn **zweimal** trägt, then derselbe Code auf der **Kopfzeile**; Given eine Datenzeile mit zu wenigen Zellen, then derselbe Code auf **dieser** Zeile.
- **Zellenlänge (Leerlauf und Fence-Treue):** Given einen Abschnitt **ohne** Tabelle — oder einen, dessen einzige Tabelle in einem Fenced-Code-Block steht —, when der Lauf endet, then `section-column-missing` als Leerlauf-Befund: eine Tabelle im Fence deklariert nichts, und eine Bedingung, die nichts misst, meldet nicht Erfolg.
- **Zellenlänge (zwei Tabellen):** Given einen Abschnitt mit zwei Tabellen, von denen nur die zweite die Spalte trägt, when der Lauf endet, then meldet **nur** die zu lange Zelle der zweiten — eine Tabelle ohne die Spalte schaltet die Messung der anderen **nicht** ab.
- **Zellenlänge (leere Zelle):** Given eine Regel mit **nur** `cell-max-chars` und eine **leere** Zelle der benannten Spalte, when der Lauf endet, then **kein** Befund — null Zeichen liegen unter jeder Obergrenze; Given dieselbe Regel **mit** `cell-min-chars`, then `section-cell-undersized` auf **ihrer** Zeile, mit einem **anderen** Grund-Code als die zu lange Zelle.
- **Zellenlänge (zwei Spalten, eine Zeile):** Given **eine** Regel, deren `table.column` **zwei** Einträge über verschiedene Spalten trägt, und eine Zeile, die in **beiden** verletzt, when der Lauf endet, then **zwei** Befunde — das Befund-Ziel trägt die Spalte, sonst fielen sie unter die Deduplikation (Datei, Zeile, Regel, Ziel, Grund) zusammen. Die **Regel-Identität** trägt die Zellen-Spalten **nicht**: sie leben innerhalb einer Regel und können keine zwei Regeln kollidieren lassen.
- **Hinweis (Happy Path):** Given eine Regel mit `hint` und einen Abschnitt, der ihre Bedingung verletzt, when `d-check --enable structure` läuft, then trägt der Befund den **verfassten** Text als Erläuterung — nicht die modul-eigene Meldung —, und der Grund-Code bleibt unverändert.
- **Hinweis (Modul-aus):** Given dieselbe Regel **ohne** `hint`, when der Lauf endet, then ist der Befundsatz byte-identisch zum Lauf vor der Einführung des Schlüssels ([`DC-QA-02`](#dc-qa-02--determinismus)): gleiche Zeile, gleicher Grund-Code, modul-eigene Meldung.
- **Hinweis (Nicht-Mess-Befunde):** Given eine Regel mit `hint`, deren `files`-Glob **keine** Datei trifft, deren Dateibaum **unlesbar** ist oder deren Kandidat als **Einzeldatei** unlesbar ist, when der Lauf endet, then trägt der Befund die **fail-closed-Ursache**, nicht den Hinweis — die Regel hat dort nicht gemessen.
- **Hinweis (Config-Rand):** Given `hint` **explizit leer** oder nur Whitespace, when `d-check` startet, then Exit 2 vor dem Lauf; Given einen `hint` mit **Tab oder Zeilenumbruch**, then ebenfalls Exit 2 — die Befund-Zeile ist tab-getrennt und einzeilig ([`DC-FA-CLI-004`](#dc-fa-cli-004--ausgabeformate)).
- **Zellenlänge (ein Selektor, viele Spalten):** Given einen Abschnitt und **vier** zu begrenzende Spalten, when sie als vier Einträge **einer** Regel konfiguriert sind, then gilt jede Zusage einzeln und `files`/`section` stehen **einmal** da; Given denselben `name` **zweimal** in einer Liste, then Exit 2 vor dem Lauf — zwei Einträge derselben Spalte trügen dasselbe Befund-Ziel.
- **Zellenlänge (Modul-aus / Config-Rand):** Given eine Regel **ohne** `table`, when `d-check` läuft, then ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)); Given einen Spalten-Eintrag ohne jede Grenze, ein `table` ohne `order` und ohne `column`, einen `name` aus lauter Weißraum, ein explizit gesetztes `cell-max-chars`/`cell-min-chars` < 1 oder eine Untergrenze **über** der Obergrenze, then Exit 2 vor dem Lauf.
- **Zellenlänge (Vorgänger-Schlüssel):** Given eine Regel, die einen der fünf flachen Vorgänger-Schlüssel auf Regel-Ebene trägt (`table-order`, `table-column`, `cell-max-column`, `cell-max-chars`, `cell-min-chars`), when `d-check` startet, then Exit 2 **mit dem neuen Ort in der Meldung** — nicht stilles Ignorieren und nicht die generische Unbekannt-Feld-Meldung des Decoders.

**Out-of-Scope:** Aussagen über den **Ort** eines Dokuments (Dateinamens-,
Verzeichnis- oder Nummerierungs-Konventionen) — keine Struktur *innerhalb* eines
Dokuments; eine **Stichtags**-Regel („erst ab Kennung N"), die die
Kennungs-Konvention des Adopters interpretieren müsste (d-check führt keinen
zweiten Regel-Interpreter, dieselbe Grenze wie die Index-Wahrheit in
[`DC-FA-TRK-001`](#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in));
**namentliche Ausnahmen innerhalb einer Datei** (einzelne Abschnitte einer Klasse
ausnehmen) — ausdrückbar ist nur die Pfad-Ausnahme, alles andere wäre eine
zweite Kennungs-Semantik; die **Bedeutung** des Inhalts; verschachtelte Regeln
und mehr als ein Abschnitts-Typ pro Regel (wer zwei Typen prüfen will, schreibt
zwei Regeln — das ist hier kein Ausweg, sondern der Normalfall, weil die Typen
verschieden sind); weitere **Schlüssel-Typen** der Chronologie-Bedingung
jenseits von ISO-Datum und Punkt-Version (geschlossene Menge — ein dritter Typ
ist ein Change Request, kein konfigurierbares Muster: kein zweiter
Regel-Interpreter, dieselbe Grenze wie oben); **Ordnungs-Aussagen über
Tabellen- oder Spalten-Grenzen hinweg** (zwei getrennt sortierte, gegenläufige
Tabellen im selben Abschnitt bleiben unerkannt — benannte Grenze; der belegte
Anlassfall ist der Richtungs-Bruch **innerhalb** einer Tabelle).

**Aufgelöste Grenze (0.72.0, Form seit 0.73.0):** *zwei Zusagen über
verschiedene Spalten desselben Abschnitts* waren bis 0.71.0 ein
Konfigurations-Duplikat (Exit 2), weil die Regel-Identität nur aus Glob und
Abschnitts-Selektor bestand. Seit 0.73.0 sind die Zellen-Spalten eine **Liste
innerhalb einer Regel** und können zwei Regeln gar nicht mehr kollidieren
lassen; getrennt werden ihre **Befunde**, über das Ziel `… :: Spalte <name>` —
ohne das fielen zwei Befunde über **dieselbe Zeile** mit **demselben**
Grund-Code unter die Deduplikation zusammen, und einer bliebe unsichtbar. Die
**Regel-Identität** trägt weiterhin ein **explizit** gesetztes
`table.order-column`: zwei Chronologie-Zusagen über denselben Abschnitt sind
zwei Regeln und brauchen zwei Identitäten. Eine Regel, die **keine** Spalte
nennt, behält ihre Identität unverändert.

---

### DC-FA-TRK-001 — Getrackt-Status auflösbarer Referenz-Ziele (Modul `tracked`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `tracked` prüft d-check für
jedes **auflösbare, existierende** repo-interne **Datei**-Ziel eines Links oder
Bildes (die Datei-Ebene
der [`DC-FA-LINK-001`](#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)-Auflösung
— dieselbe Auflösungs-**Mechanik**, unabhängig davon, ob das Modul `links`
aktiv ist; Anker sind hier irrelevant), ob die Zieldatei im **git-Index
getrackt** ist. Verzeichnis-Ziele sind kein Kandidat (der Index führt nur
Dateien); Symlink-Referenzen — das Ziel ist oder durchläuft einen Symlink —
sind kategorisch die Domäne von
[`DC-FA-LINK-002`](#dc-fa-link-002--symlink-ablehnung) (`links`
meldet `symlink`), `tracked` prüft sie nicht (der Index führt den realen
Pfad, nicht den Alias — sonst false-positive hinter getrackten
Verzeichnis-Symlinks).
Existiert das Ziel nur im Arbeitsbaum — untracked, gleich ob vergessen oder
gitignoriert —, entsteht der Befund `target-untracked`: die Referenz ist beim
Erzeuger grün, wäre aber auf **jedem frischen Klon** ein `target-missing`
(Umgebungs-Drift zwischen Arbeitsbäumen). Das Modul fängt diese Falle am
Entstehungsort statt erst in der CI des nächsten Checkouts.

Die **Wahrheit ist der Index**, nicht die `.gitignore`-Syntax: d-check parst
keine ignore-Regeln (kein zweiter Regel-Interpreter); eine frisch per `git add`
gestagte, noch nie committete Datei gilt als **getrackt** — neue Doku wird mit
dem Staging grün, der Inner-Loop bleibt arbeitsfähig.

Dafür liest `tracked` das im read-only-Mount vorhandene `.git` über den
**VCS-Port** ([`DC-FA-VCS-001`](#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)/[`DC-FA-COMMITS-001`](#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in):
reine-Go-Implementierung, **ohne git-Binary**, **ohne Netz**) — aber **ohne
Commit-Range**: geprüft wird ausschließlich der Index-Stand des Arbeitsbaums.
Die Eingabe ist gegenüber den hermetischen Modulen erweitert (das `.git`),
bleibt aber **lokal, lesend und deterministisch**: derselbe Arbeitsbaum +
derselbe Index ⇒ derselbe Befundsatz ([`DC-QA-02`](#dc-qa-02--determinismus)),
kein Netzzugriff, kein Schreiben
([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

**Kein Doppelbefund:** ein nicht existierendes Ziel bleibt `target-missing`
([`DC-FA-LINK-001`](#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links));
`tracked` prüft nur Ziele, die die Link-Auflösung erfolgreich fand (dasselbe
Prinzip wie [`DC-FA-PIN-001`](#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in)
— struktureller Befund bleibt beim Struktur-Modul). Ventil:
`tracked.exempt-targets` (Glob über den **aufgelösten** Zielpfad, Syntax wie
`scan.ignore`) nimmt absichtlich untrackte Ziele aus (z. B. lokal generierte
Artefakte, deren Erzeugung der Konsument selbst verantwortet).

**Strikt opt-in, fail-closed, diagnose-only:** `tracked` ist nie Default-Modul
(wie `vcs`/`commits`); ohne aktives `tracked` ist der Befundsatz byte-identisch
und nichts wird über den Markdown-Baum hinaus gelesen. **fail-closed:** aktives
`tracked` ohne lesbares `.git` ⇒ Fehler (Exit 2) mit Hinweis auf stderr — kein
stilles Grün. `target-untracked` liefert keinen `--repair`-Hunk (`git add` bzw.
Committen ist eine menschliche Entscheidung; vgl.
[`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch)).

**Akzeptanzkriterien:**

- **Happy Path:** Given `tracked` aktiv und ein Link auf eine existierende, im git-Index getrackte Datei, when `d-check --enable tracked` läuft, then kein Befund, Exit 0.
- **Boundary (Index-Wahrheit):** Given ein Link auf eine existierende Datei, die frisch per `git add` gestagt und noch nie committet ist, when `d-check --enable tracked` läuft, then kein Befund — der Index entscheidet, nicht die Commit-Historie.
- **Boundary (Modul-aus / git-frei):** Given **kein** aktives `tracked`, when `d-check` in einer netzlosen, read-only Umgebung läuft, then ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)) und es erfolgt kein `.git`-Zugriff ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **Negative:** Given `tracked` aktiv und ein Link auf eine Datei, die im Arbeitsbaum existiert, aber nicht im Index getrackt ist (untracked oder gitignoriert), when `d-check --enable tracked` läuft, then ein Befund `target-untracked` (Datei, Zeile, aufgelöstes Ziel), Exit 1.
- **Kein Doppelbefund:** Given `tracked` aktiv und ein Link auf eine **nicht existierende** Datei, when `d-check --enable tracked` läuft, then nur `target-missing` (Modul `links`), kein zusätzliches `target-untracked`.
- **Ventil:** Given ein untracktes, existierendes Ziel, dessen aufgelöster Pfad einem `tracked.exempt-targets`-Glob entspricht, when `d-check --enable tracked` läuft, then kein Befund für dieses Ziel.
- **fail-closed (git-Eingabe fehlt):** Given `tracked` aktiv, aber **kein** lesbares `.git` unter der Scan-Wurzel, when `d-check --enable tracked` läuft, then **Exit 2** mit Hinweis auf stderr — kein stilles Grün.

**Out-of-Scope:** Inline-Code-Pfade
([`DC-FA-CODE-001`](#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)-Ziele)
in der ersten Ausbaustufe (Re-Evaluierung nach Praxis-Evidenz); der
Getrackt-Status **gescannter Dateien selbst** (würde jede WIP-Datei vor dem
ersten `git add` flaggen — die Gegenrichtung der Drift bleibt bei der
`scan.ignore`-Disziplin und der CI als Netz); **Verzeichnis**-Ziele (der
Index führt nur Dateien — ein gitignoriertes Verzeichnis ohne getrackte
Dateien als Linkziel bleibt ungefangen, Re-Evaluierung nach
Praxis-Evidenz) und **Symlink**-Referenzen (kategorisch
[`DC-FA-LINK-002`](#dc-fa-link-002--symlink-ablehnung)-Domäne);
Interpretation der `.gitignore`-Syntax (der Index ist die einzige
Wahrheit); intent-to-add-Feinheiten (`git add -N` gilt als getrackt);
Submodule und verschachtelte Arbeitsbäume (linked worktree ⇒ Exit 2);
Schreiben/`git add` durch das Werkzeug (read-only).

---

### DC-FA-TGT-001 — Deklarations-Konsistenz zwischen Doku und Build-Targets (Modul `targets`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `targets` prüft d-check, dass
die in der Doku als ` `make X` ` **behaupteten** Build-Targets und die real im
Makefile definierten Targets übereinstimmen — die maschinelle Form des
Vertrags „kein halluziniertes Gate / kein undokumentiertes Gate" (Meta-Gate
gegen Harness-Lügen). Zwei Richtungen mit je eigenem Grund-Code:

- **Phantom-Gate (`gate-phantom`):** ein ` `make X` `-Token in einer
  konfigurierten Doku-Datei (`targets.doc-tables`), zu dem **keine** Regel in
  einer konfigurierten Makefile-Quelle (`targets.makefiles`) existiert — die
  Doku behauptet ein Gate, das es nicht gibt (Befund an Datei:Zeile der
  Doku-Behauptung).
- **Undeklariertes Gate (`gate-undocumented`):** eine Makefile-Regel (minus
  `targets.exempt-targets`), die in der Autoritäts-Doku (`targets.authority`)
  **nicht** als ` `make X` ` steht — ein Gate ohne Deklaration (Befund an
  Datei:Zeile der Makefile-Regel).

**Tabellen-Scoping (Erkennungs-Vertrag):** ` `make X` ` gilt **nur in
Tabellenzeilen** (Zeilen mit `|`-Präfix) als Existenz-/Vollständigkeits-
Behauptung — Prosa-Erwähnungen (z. B. „Richtig: ` `make gates` `") zählen
**nicht**. Ohne diese Regel würden in Prosa erwähnte, entfernte Targets zu
spuriösen `gate-phantom`-Befunden.

Dafür liest `targets` die konfigurierten Doku- und Makefile-Dateien über den
**Filesystem-Port** (`ReadFile`) — **hermetisch** wie das Modul `planning`
([`DC-FA-PLAN-001`](#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)):
**kein** git, **kein Ausführen** des Makefile, **kein** Netz. Die
Target-Erkennung im Makefile ist eine Zeilen-Heuristik (literale Regelnamen am
Zeilenanfang, ohne Zuweisungen und `.PHONY`); Pattern-Rules und variabel
benannte Targets sind kein Kandidat. Die Eingabe bleibt lokal, lesend und
deterministisch ([`DC-QA-02`](#dc-qa-02--determinismus)/[`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
Dieselbe Klasse wie `planning` — „Doku-Behauptung ↔ Repo-Struktur": dort
Roadmap ↔ Verzeichnis, hier Doku-Tabelle ↔ Makefile-Regeln.

**Strikt opt-in, fail-closed, diagnose-only:** `targets` ist nie Default-Modul;
ohne aktives `targets` ist der Befundsatz byte-identisch
([`DC-QA-02`](#dc-qa-02--determinismus)) und nichts über den Markdown-Baum
hinaus wird gelesen. **fail-closed:** aktiv mit einer fehlenden konfigurierten
Datei (Makefile oder Doku-Datei) ⇒ **Exit 2** mit Hinweis auf stderr — kein
stilles Grün. Leeres `targets.makefiles` ⇒ Modul inert (keine Regelmenge);
leeres `targets.doc-tables` ⇒ Richtung 1 (Phantom) entfällt, leeres
`targets.authority` ⇒ Richtung 2 entfällt — die beiden Richtungen sind
voneinander unabhängig.
`gate-phantom`/`gate-undocumented` liefern keinen `--repair`-Hunk
([`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch)) — Doku- bzw. Makefile-Edit
ist eine menschliche Entscheidung.

**Akzeptanzkriterien:**

- **Happy Path:** Given `targets` aktiv, eine Doku-Tabellenzeile mit ` `make foo` `, ein Makefile mit der Regel `foo:` und ein Eintrag ` `make foo` ` in der Autoritäts-Doku, when `d-check --enable targets` läuft, then kein Befund, Exit 0.
- **Boundary (Prosa zählt nicht):** Given `targets` aktiv und ein ` `make bar` ` **nur in Prosa** (keine Tabellenzeile), zu dem keine Makefile-Regel existiert, when `d-check --enable targets` läuft, then **kein** `gate-phantom` — nur Tabellenzeilen sind Existenz-Behauptung.
- **Boundary (Modul-aus):** Given **kein** aktives `targets`, when `d-check` in einer netzlosen, read-only Umgebung läuft, then ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)) und weder das Makefile noch die Doku-Tabellen-Konsistenz werden gelesen/geprüft ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **Negative (Phantom):** Given `targets` aktiv und eine Doku-Tabellenzeile mit ` `make ghost` `, zu der **keine** Makefile-Regel existiert, when `d-check --enable targets` läuft, then ein Befund `gate-phantom` (Datei, Zeile, Target `ghost`), Exit 1.
- **Negative (undokumentiert):** Given `targets` aktiv und eine Makefile-Regel `secret:` (nicht in `targets.exempt-targets`), die in der Autoritäts-Doku fehlt, when `d-check --enable targets` läuft, then ein Befund `gate-undocumented` (Makefile, Zeile, Target `secret`), Exit 1.
- **Boundary (exempt):** Given `targets` aktiv und eine Makefile-Regel `clean:`, die in `targets.exempt-targets` steht und in der Autoritäts-Doku **fehlt**, when `d-check --enable targets` läuft, then **kein** `gate-undocumented` für `clean` (Utility-Regeln sind von der Doku-Pflicht ausgenommen).
- **fail-closed (fehlende Datei):** Given `targets` aktiv und eine konfigurierte Makefile-/Doku-Datei existiert nicht, when `d-check --enable targets` läuft, then **Exit 2** mit Hinweis auf stderr — kein stilles Grün.

**Out-of-Scope:** Ausführen des Makefile oder Auflösen von `include`-Direktiven,
Variablen- oder Pattern-Rule-Targets (rein statische Zeilen-Heuristik, identisch
zum abgelösten `tools/gate-consistency.sh`); Section-/Heading-Scoping innerhalb
der Autoritäts-Doku (die **ganze** Datei wird nach Tabellenzeilen gescannt — ein
optionaler `authority`-Section-Anker bleibt spätere Anforderung); die
`.d-check.yml`-Modul-Listen-Selbstkonsistenz des netzlosen doc-check
(repo-spezifische Prüfung der Netzlos-Gate-Config, **kein** cross-repo-Kern,
verbleibt im repo-lokalen Rest von `gate-consistency.sh`); Nicht-Make-Build-
Systeme; ein Auto-Fix (`--repair`).

---

### DC-FA-SRC-001 — Upstream-Content-Drift externer Quellen (Modul `sources`, opt-in, Netz)

**Beschreibung:** Bei explizit aktiviertem Modul `sources` prüft d-check, dass
eine auf einen `sha256` **gepinnte externe Quelle** seit dem Pinnen inhaltlich
unverändert ist: d-check holt die Quelle über den Netz-Port, errechnet ihren
Hash und vergleicht ihn mit dem hinterlegten Pin. Weicht er ab, meldet es
`source-drift`; ist die Quelle nicht erreichbar, `source-unreachable`. Das ist
das **Netz-Gegenstück** zum in-repo-Content-Pin `pins`
([`DC-FA-PIN-001`](#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in))
und teilt die Netz-Natur mit `external`
([`DC-FA-EXT-001`](#dc-fa-ext-001--externe-links-modul-external-opt-in)) — es
prüft aber **Inhalt** (Hash), nicht bloße Erreichbarkeit.

**Zwei Deklarations-Flächen (beide gültig):**

- **Marker** am externen Link: `[text](https://…) <!-- source-pin: sha256:<hex> -->`
  bindet an den unmittelbar vorausgehenden `http(s)`-Link derselben Zeile (wie
  `pins`), sonst inert. Für Archiv-Ziele trägt der Marker das explizite
  Schlüsselwort `zip` (parallel zu `unpack: zip`):
  `<!-- source-pin: zip sha256:<hex> -->`.
- **Config** in `.d-check.yml`: ein Top-Level-Block `sources:` mit Einträgen
  `{url, sha256, unpack: none|zip}`.

**Zwei Quelltypen:**

- **Einzeldatei (`unpack: none`, Default):** Hash über die **rohen
  Antwort-Bytes**.
- **Archiv (`unpack: zip`):** Hash über ein **Content-Manifest**, nicht über die
  Zip-Roh-Bytes (die durch Recompression/Framing instabil sind): je enthaltene
  Datei ihr Inhalts-`sha256`, die Zeilen **nach dem normalisierten Pfad**
  sortiert konkateniert, davon der `sha256` (byte-genaue Kanonisierung in der
  Spezifikation). Damit ist der Pin **invariant gegen die
  Zip-Eintrags-Reihenfolge**; konzeptionell wie das committet-vendored
  `SHA256SUMS`-Muster, aber eigenständig kanonisiert (kein byte-identischer
  `sha256sum`-Abzug).

**Nur externe http(s)-Ziele, kein Doppelbefund:** `sources` prüft ausschließlich
absolute `http`/`https`-URLs; repo-interne Ziele bleiben Domäne von
`links`/`pins` (kein zweiter Befund). Strukturelle Erreichbarkeit ohne Pin
bleibt Sache von `external`.

**Ergonomie (voller Hash):** Der `source-drift`-Befund führt den
**vollständigen** errechneten `sha256` in der Meldung — Pinnen ist damit „einmal
laufen, gemeldeten Hash in Marker/Config kopieren", ohne externes Hash-Werkzeug.

**Strikt opt-in, Netz nur hier, fail-closed, diagnose-only:** `sources` ist nie
Default-Modul und — neben `external` — die einzige Netz-Tür
([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)); ohne
aktives `sources` ist der Befundsatz byte-identisch
([`DC-QA-02`](#dc-qa-02--determinismus)) und **keine** Netzverbindung wird
geöffnet. Der `sha256` wird **case-insensitiv** geführt. Der Fetch ist
**größenbegrenzt** (Body- und Entpack-Obergrenze, Zip-Bomben-Schutz); jede nicht
materialisierbare Antwort — Netzfehler, HTTP-Status ≥ 400, Timeout,
überschrittenes Limit **oder** ein unter `unpack: zip` ungültiges Zip — ⇒
`source-unreachable` (getrennt von `source-drift`). **fail-closed:** eine
malformte Direktive (kein `sha256:<hex>`) oder ein ungültiger Config-Eintrag
(fehlende `url`/`sha256`, unbekanntes `unpack`) ⇒ **Exit 2** mit Hinweis auf
stderr — kein stilles Grün.
`source-drift`/`source-unreachable` liefern keinen `--repair`-Hunk
([`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch)) — ein Re-Pin ist eine
menschliche Entscheidung.

**Akzeptanzkriterien:**

- **Happy Path:** Given `sources` aktiv, eine per Marker gepinnte URL und ein per Config gepinntes Archiv (`unpack: zip`), deren Inhalte den hinterlegten `sha256` treffen, when `d-check --enable sources` läuft, then kein Befund, Exit 0.
- **Boundary (Archiv-Determinismus):** Given ein gepinntes Archiv, dessen Zip-Einträge upstream **umsortiert**, inhaltlich aber identisch sind, when `d-check --enable sources` läuft, then **kein** `source-drift` — der Manifest-Hash ist eintrags-reihenfolge-invariant.
- **Boundary (Modul-aus / netzlos):** Given **kein** aktives `sources`, when `d-check` in einer netzlosen, read-only Umgebung läuft, then ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)) und es wird **keine** Netzverbindung geöffnet ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **Negative (Drift):** Given `sources` aktiv und ein gepinntes Ziel, dessen Inhalt sich seit dem Pinnen geändert hat, when `d-check --enable sources` läuft, then genau ein Befund `source-drift` (Datei/Zeile bzw. Config-Eintrag, URL, **vollständiger** errechneter `sha256`), Exit 1.
- **Boundary (unerreichbar ≠ Drift):** Given `sources` aktiv und ein gepinntes Ziel, das nicht materialisierbar ist (Netzfehler, HTTP-Status ≥ 400, Timeout, Größenlimit **oder** unter `unpack: zip` kein gültiges Zip), when `d-check --enable sources` läuft, then ein Befund `source-unreachable` (nicht `source-drift`), Exit 1.
- **Boundary (kein Doppelbefund / repo-intern):** Given `sources` aktiv und ein `source-pin`-Marker an einem **repo-internen** Link, when `d-check --enable sources` läuft, then **kein** `sources`-Befund — repo-interne Ziele bleiben `pins`/`links`-Domäne.
- **Boundary (Marker-Bindung):** Given `sources` aktiv und ein `source-pin`-Marker **ohne** unmittelbar vorausgehenden `http(s)`-Link auf derselben Zeile, when `d-check --enable sources` läuft, then ist der Marker **inert** (kein Befund); bei **mehreren** Links auf der Zeile bindet der Marker an den **unmittelbar links** stehenden.
- **fail-closed (malform):** Given `sources` aktiv und eine malformte `source-pin`-Direktive (kein `sha256:<hex>`) oder ein ungültiger `sources`-Config-Eintrag (unbekanntes `unpack`), when `d-check` läuft, then **Exit 2** mit Hinweis auf stderr — kein stilles Grün.

**Out-of-Scope:** Currency/„neuerer Tag verfügbar" (bleibt der Bash-Helfer
`tools/harness/fetch-baseline-cache.sh --check-latest`,
[`MR-022`](../harness/conventions.md#mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019));
Nicht-`http`/`https`-Schemata (`file:`, `git:`, …); Authentifizierung,
Custom-Header oder Credentials; repo-interne Ziele (Domäne `links`/`pins`);
Verzeichnis-/Listing-Ziele ohne Archiv-Bundle; andere Archiv-Formate als `zip`
(spätere Anforderung); ein Auto-Fix/Re-Pin (`--repair`); die Hash-Emission des
Bestandsmoduls `pins`/`dpin` (separater Kleinst-Change am Bestandsmodul); ein
modul-lokaler `sources.scope` (die Config-Fläche ist eine bare Liste `sources[]`
— `source-pin`-Marker werden über den **globalen** Scan-Scope gesammelt;
`sources.scope` ⇒ Exit 2 als unbekannter Schlüssel, bewusste Ausnahme zu
[`DC-FA-CONF-002`](#dc-fa-conf-002--modul-lokaler-scan-scope)).

---

### DC-FA-WF-001 — Deklarations-Konsistenz von Workflow-Referenzen (Modul `workflows`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `workflows` prüft d-check die
`uses:`-Referenzen von CI-Workflow-Dateien auf zwei Deklarations-Zusagen: dass
eine **fremde** Referenz unbeweglich gepinnt ist, und dass eine **lokale**
Referenz auflöst **und** vom aufrufenden Job die Rechte bekommt, die sie
verlangt. Geprüft wird die **Deklaration**, nicht die Lauffähigkeit.

| Bedingung | Grund-Code | Reparatur |
|---|---|---|
| fremde Referenz nennt einen **vollen 40-stelligen** Commit-SHA | `uses-pin-missing` | den beweglichen Tag durch den SHA ersetzen |
| … und dahinter einen Kommentar, der den Tag nennt | `uses-pin-untagged` | den Tag-Kommentar ergänzen |
| lokale Referenz (`./…`) zeigt auf eine existierende Datei | `uses-local-missing` | den Pfad korrigieren |
| der aufrufende **Job** trägt ein eigenes `permissions:`, wenn das Ziel Rechte verlangt | `uses-local-perms-undeclared` | dem Job die Rechte deklarieren |
| … und führt jeden geforderten Scope mindestens so hoch | `uses-local-perms-narrow` | den Scope anheben |
| jede gelesene Datei ist gültiges YAML | `workflow-unparsable` | die Datei reparieren |

**Warum der Pin.** Ein Tag lässt sich umhängen, ein SHA nicht; eine bewegliche
Referenz ist eine Supply-Chain-Fläche. Geprüft wird die **Form** des Pins, nicht
seine **Gültigkeit** — ob der SHA existiert und den Commit bezeichnet, den der
Tag-Kommentar behauptet, ist eine Netz-Frage und ausdrücklich außerhalb.

**Warum die lokale Referenz anders behandelt wird.** Sie kann keinen SHA tragen
und **braucht keinen**: sie löst auf denselben Commit auf wie ihr Aufrufer und
ist damit **stärker** gebunden als ein Pin. An die Stelle des Pin-Checks treten
zwei Fragen, die dort Sinn ergeben — **existiert** das Ziel, und **bekommt es
die Rechte, die es verlangt?** Ein aufgerufener Workflow erhält nur, was der
aufrufende **Job** selbst führt; verlangt er mehr, bricht der ganze Lauf vor dem
ersten Job ab. Ein Job **ohne** eigenes `permissions:` erbt den Workflow-Kopf —
genau diese stille Vererbung ist der belegte Ausfall.

**Rechte-Vergleich.** Die Stufen sind geordnet: `none` < `read` < `write`. Ein
Scope, den der Aufrufer **nicht nennt**, zählt als `none`. `read-all` und
`write-all` setzen alle Scopes auf die jeweilige Stufe; ein leeres Mapping
(`{}`) verlangt bzw. gewährt nichts.

**Scan-Menge und die Grenze, die daraus folgt.** Das Modul liest die
Workflow-Dateien unterhalb eines **konfigurierten** Verzeichnisses
(`workflows.dir`) — der Ort ist nicht verdrahtet, weil er CI-System-spezifisch
ist. **Es liest darüber hinaus Dateien, die es nicht scannt:** das Ziel jeder
lokalen Referenz. Für diese Ziele gilt **dieselbe** Zusage wie für die
Scan-Menge — sie werden geparst, und ein Parse-Fehler ist ein Befund, kein
Übersprung. Ein Ziel **außerhalb** des konfigurierten Verzeichnisses wird
ebenso gelesen; die Referenz bestimmt die Eingabe, nicht die Scan-Wurzel.

**Hermetisch und netzlos:** nur der Filesystem-Port, **kein** git, **kein**
Netz, **kein Ausführen** eines Workflows.

**Strikt opt-in, fail-closed:** `workflows` ist nie Default-Modul; ohne
`workflows.dir` ist es **inert** (keine Datei wird geöffnet, der Befundsatz ist
byte-identisch, [`DC-QA-02`](#dc-qa-02--determinismus)). Ist es aktiv und das
Verzeichnis fehlt, enthält keine Workflow-Datei oder keine einzige
`uses:`-Referenz, ist das ein **Befund**, kein stilles Grün: ein Wächter, der
nichts gefunden hat, und einer, der nichts zu prüfen hatte, sähen im Exit-Code
sonst identisch aus.

**Exit 2 vor dem Lauf** bei: leerem `workflows.dir` (nur Weißraum) und
ungültigem Glob in `workflows.exempt-paths`.

**Akzeptanzkriterien:**

- **Happy Path:** Given `workflows` aktiv, eine fremde Referenz mit vollem SHA plus Tag-Kommentar und eine lokale Referenz, deren Ziel existiert und deren aufrufender Job die geforderten Rechte deklariert, when `d-check --enable workflows` läuft, then kein Befund, Exit 0.
- **Negative (Pin):** Given eine fremde Referenz auf einen beweglichen Tag (`actions/checkout@v4`), when der Lauf endet, then `uses-pin-missing` auf **ihrer** Zeile; Given einen vollen SHA **ohne** Tag-Kommentar, then `uses-pin-untagged`.
- **Negative (lokale Referenz):** Given eine lokale Referenz auf eine nicht existierende Datei, when der Lauf endet, then `uses-local-missing` auf ihrer Zeile.
- **Negative (stille Vererbung):** Given einen Workflow-Kopf mit `permissions: {}`, einen Job **ohne** eigenes `permissions:` und ein lokales Ziel, das `contents: read` verlangt, when der Lauf endet, then `uses-local-perms-undeclared` — der Job erbt sonst den Kopf, und der Lauf bricht vor dem ersten Job ab.
- **Negative (zu enger Scope):** Given einen Job mit `contents: read` und ein Ziel, das `contents: write` verlangt, then `uses-local-perms-narrow`; Given ein Ziel, das einen vom Aufrufer **nicht genannten** Scope verlangt, then derselbe Code — ein nicht genannter Scope ist `none`.
- **Boundary (`read-all`):** Given einen Aufrufer mit `permissions: read-all` und ein Ziel, das `contents: read` verlangt, when der Lauf endet, then **kein** Befund; verlangt das Ziel `contents: write`, then `uses-local-perms-narrow`.
- **Boundary (Ziel ohne Forderung):** Given ein lokales Ziel **ohne** `permissions`, when der aufrufende Job ebenfalls keines trägt, then **kein** Befund — es gibt nichts zu vergleichen.
- **Boundary (fremdes Repository):** Given eine Referenz der Form `owner/repo/.github/workflows/x.yml@<sha>`, when der Lauf endet, then wird sie wie jede Action auf den Pin geprüft und **nicht** auf Rechte — ihr Inhalt liegt nicht im Repo.
- **fail-closed (leere Prüfmenge):** Given `workflows` aktiv und ein Verzeichnis ohne Workflow-Datei — oder ohne eine einzige `uses:`-Referenz —, when der Lauf endet, then ein Befund, nicht Exit 0.
- **fail-closed (unlesbares YAML):** Given eine Workflow-Datei oder ein Referenz-Ziel, das kein gültiges YAML ist, when der Lauf endet, then `workflow-unparsable` — nicht Übersprung.
- **Boundary (Modul-aus):** Given **kein** aktives `workflows`, when `d-check` läuft, then ist der Befundsatz byte-identisch zum Lauf ohne den Konfigurations-Block ([`DC-QA-02`](#dc-qa-02--determinismus)), und keine Workflow-Datei wird geöffnet.

**Out-of-Scope:** die **Gültigkeit** eines SHA (Netz — gehört zur
Freshness-Familie); die **Rechte** einer Referenz in ein fremdes Repository
(deren Inhalt liegt nicht vor); das Ausführen oder Simulieren eines Workflows;
jede weitere Semantik von GitHub Actions (Matrix-Auflösung,
Bedingungs-Auswertung, Secrets-Verfügbarkeit) — das Modul deckt **eine**
Deklarations-Klasse, und ein grüner Lauf sagt „diese Klasse liegt nicht vor",
nicht „der Workflow läuft".

---

### DC-FA-CONF-001 — Konfigurationsdatei

**Beschreibung:** Eine optionale Datei `.d-check.yml` in der
Repository-Wurzel deklariert: Scan-Wurzeln, Ignorier-Muster, aktive
Module, Kennungs-Muster (für `ids`), Dokumentklassen und
Referenzregeln (für `matrix`) sowie Modul-Parameter (z. B. Timeout für
`external`). Der Datei**pfad** ist Konvention, nicht Zwang: er lässt sich pro
Lauf überschreiben ([`DC-FA-CLI-012`](#dc-fa-cli-012--konfigurations-pfad-überschreiben)) —
Format, Validierung und Fehlerverhalten bleiben davon unberührt.
Ohne Konfigurationsdatei gelten die in den jeweiligen
Anforderungen definierten Defaults; eine vorhandene Datei wird beim
Start vollständig validiert — Syntax **und** Inhalt (z. B. ungültige
Regexe, unbekannte Modulnamen). Jeder Konfigurationsfehler führt zu
Exit-Code 2; es findet keine Prüfung mit stillschweigenden Defaults
statt.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine `.d-check.yml` mit einem Ignorier-Muster `docs/archive/**`, when `d-check` läuft, then erzeugen kaputte Links unter `docs/archive/` keine Befunde. <!-- d-check:ignore (AK-Beispiel) -->
- **Boundary:** Given keine `.d-check.yml`, when `d-check` läuft, then laufen genau die Default-Module (`links`, `anchors`) auf den Default-Scan-Wurzeln.
- **Negative:** Given eine syntaktisch ungültige `.d-check.yml`, when `d-check` läuft, then Exit-Code 2 mit Fehlermeldung inkl. Zeilenangabe; es wird keine Prüfung mit stillschweigenden Defaults durchgeführt.

**Out-of-Scope:** Konfigurations-Vererbung über mehrere Verzeichnisebenen; Migrations-Assistent von Alt-Tool-Konfigurationen.

---

### DC-FA-CONF-002 — Modul-lokaler Scan-Scope

**Beschreibung:** Jedes **scannende** Regelmodul akzeptiert in der
Konfigurationsdatei optional einen Schlüssel `<modul>.scope` mit den
Unterschlüsseln `roots` und `ignore`. **Nicht** scannende Module kennen ihn
nicht und lehnen ihn ab: `planning` und `targets` sind Post-Pässe über
benannte Dateien, `sources` arbeitet über eine Pin-Liste, und
[`DC-FA-STRUCT-001`](#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
bestimmt seine Eingabe je Regel selbst. Ein `scope` wäre dort wirkungslos —
und ein wirkungsloser Schlüssel, den der Decoder annimmt, ist ein stilles Grün. Ist er gesetzt, **ersetzt** er
für genau dieses Modul den globalen Scan-Scope
(`scan.roots`/`scan.ignore`) — bewusst keine Schnittmengen- oder
Vererbungs-Semantik (einfachste ehrliche Form). Module ohne `scope`
erben den globalen Scan unverändert; bestehende Konfigurationen
bleiben unverändert gültig (additiv, abwärtskompatibel). Für
`scope.roots` und `scope.ignore` gelten dieselben Regeln und
Constraints wie für die globalen Schlüssel: deklarierte Wurzeln
müssen existieren und innerhalb der Repository-Wurzel liegen (jede
Verletzung Exit 2,
[`DC-FA-CONF-001`](#dc-fa-conf-001--konfigurationsdatei)),
Ignorier-Muster prunen den Verzeichnis-Abstieg, die `SKIP_DIRS`
gelten immer, und eine explizit leere `roots`-Liste prüft nichts
(Leere-Menge-Semantik wie `scan.roots`).

**Akzeptanzkriterien:**

- **Happy Path:** Given global `scan.roots: ["."]` und ein `ids`-`scope` mit `roots: ["spec"]`, when `d-check` läuft, then stammen `id-unlinked`-Befunde ausschließlich aus `spec/`, während `links`/`anchors` weiterhin den globalen Scope prüfen.
- **Boundary:** Given ein Modul ohne `scope`-Schlüssel, when `d-check` läuft, then verhält es sich byte-identisch zum Verhalten ohne diese Anforderung.
- **Boundary:** Given ein Modul-`scope` mit explizit leerer `roots`-Liste, when `d-check` läuft, then prüft dieses Modul nichts — bewusste Setzung, kein Fehler.
- **Negative:** Given ein `scope.roots`-Eintrag, der nicht existiert oder die Repository-Wurzel verlässt, when `d-check` startet, then Exit-Code 2 mit Fehlermeldung; es findet keine Prüfung mit stillschweigenden Defaults statt.

**Out-of-Scope:** Scope je Muster (`ids.patterns[].roots` — erst bei zweitem Bedarfsträger); Schnittmengen-/Unions-Semantik zwischen globalem und Modul-Scope; alternative Konfigurationsdateien, Mehrfach-Configs und Vererbung (Out-of-Scope-Linie von [`DC-FA-CONF-001`](#dc-fa-conf-001--konfigurationsdatei)).

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

**Out-of-Scope:** Distributionswege jenseits der Container-Registries (Homebrew, Paketmanager, Release-Binaries); deren Bewertung erfolgt später per ADR. Der **Spiegel** nach Docker Hub ist seit [`DC-FA-DIST-002`](#dc-fa-dist-002--docker-hub-spiegel) **nicht** mehr out-of-scope; das Bezugs-Ziel dieser Anforderung bleibt GHCR.

### DC-FA-DIST-002 — Docker-Hub-Spiegel

**Beschreibung:** Jedes nach GHCR veröffentlichte Image wird zusätzlich als
`docker.io/pt9912/d-check` gespiegelt — **dasselbe Bild unter denselben Tags**,
nicht ein zweiter Bau. Prüfgröße ist der **Config-Digest**, aus **beiden
Registries** gelesen: er ist die Identität des Bild-Inhalts und über Registries
hinweg stabil. Der **Manifest**-Digest ist es ausdrücklich **nicht** — er hängt
an der Blob-Kompression des jeweiligen Registrys und fällt je Registry
verschieden aus; wer `docker.io/…@sha256:…` pinnt, nimmt deshalb den
Docker-Hub-Digest, nicht den von GHCR. GHCR bleibt die **Quelle**, Docker Hub
der Spiegel; die Richtung ist Teil der Zusage. Die Tagging-Disziplin ist die von
[`DC-FA-DIST-001`](#dc-fa-dist-001--docker-image): volle Semver-Tags, und
`:latest` bewegt sich ausschließlich für stabile Releases. Die Spiegelung ist
**fail-closed** — schlägt sie fehl, ist das Release fehlgeschlagen, auch wenn
GHCR bereits trägt; die Meldung nennt dann, was schon veröffentlicht ist, damit
der Zustand nicht geraten werden muss.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein veröffentlichtes Release `vX.Y.Z`, when der Config-Digest von `docker.io/pt9912/d-check:vX.Y.Z` und der von `ghcr.io/pt9912/d-check:vX.Y.Z` **aus den Registries** gelesen werden, then sind beide identisch.
- **Boundary:** Given ein Release mit Prerelease-Suffix (`vX.Y.Z-rc1`), when die Veröffentlichung läuft, then trägt Docker Hub den Versions-Tag, und `:latest` bleibt dort unverändert auf dem letzten stabilen Release.
- **Negative:** Given fehlende oder ungültige Docker-Hub-Zugangsdaten, when die Veröffentlichung läuft, then bricht sie **vor** dem Spiegel-Push mit einer Fehlermeldung ab, die den bereits veröffentlichten GHCR-Stand ausdrücklich benennt — kein stiller Durchlauf und keine Teil-Veröffentlichung ohne Aussage.

**Out-of-Scope:** Spiegel auf weitere Registries; ein vom GHCR-Bild **abweichender** Docker-Hub-Bau (etwa andere Basis oder andere Plattform-Matrix) — die Zusage ist Inhalts-Gleichheit, nicht Parallelbau; **Gleichheit des Manifest-Digests** (registry-lokal, siehe oben); der **Inhalt** der Hub-Beschreibungsseite — er wird aus [`packaging/dockerhub/`](../packaging/dockerhub/README.md) gesetzt, ist aber nicht Teil der Distributions-Zusage: sein Fehlschlag lässt das Release grün.

## 4. Nichtfunktionale Anforderungen

### DC-QA-01 — Performance

- **Anforderung:** Ein vollständiger Lauf der Default-Module über ein Repository mit 1.000 Markdown-Dateien (insgesamt ≤ 20 MB) terminiert in < 5 s auf Standard-Entwickler-Hardware (2 vCPU); das Modul `external` ist von dieser Schranke ausgenommen.
- **Messmethode:** Benchmark-Lauf gegen ein generiertes Fixture-Repo, Definition in `spec/spezifikation.md`.

### DC-QA-02 — Determinismus

- **Anforderung:** Identische Eingabe (Repo-Stand + Konfiguration + Optionen) liefert byte-identische Ausgabe; Befunde sind stabil sortiert (Pfad, dann Zeilennummer).
- **Messmethode:** Test führt denselben Lauf 10× aus und vergleicht Hashes der Ausgabe.

### DC-QA-03 — Seiteneffektfreiheit und Netzwerk-Sparsamkeit

- **Anforderung:** Das Tool schreibt nie in das geprüfte Repository und öffnet außer in den explizit aktivierten Modulen `external` und `sources` keine Netzwerkverbindungen.
- **Messmethode:** Integrationstest mit read-only-Mount und netzwerkloser Umgebung (`docker run --network none`), alle Module außer `external`, `sources` und `vcs` aktiv (alle drei brauchen eine explizite Eingabe jenseits des gescannten Markdown-Baums — `external` und `sources` das Netz, `vcs` eine git-Range — und sind nie Teil des netzlosen Default-Laufs; `sources` fetcht nur die explizit gepinnten Quellen, nie den Markdown-Baum; `vcs` liest `.git` zwar read-only/netzlos, fail-closed ohne Range).

### DC-QA-04 — Migrationsabdeckung der Alt-Tools

- **Anforderung:** Jedes der dreizehn Quell-Tool-Vorkommen ist durch `d-check` mit passender Konfiguration ersetzbar: Auf dem jeweiligen Repo-Stand meldet `d-check` mindestens dieselben echten Befunde wie das Alt-Tool und erzeugt keine False-Positives, die eine bislang grüne CI brechen.
- **Messmethode:** Pilot-Migration in mindestens drei Repos — je ein Vertreter der Shell-Familie (`verify-doc-refs.sh`), der Python-Familie (`check_refs.py`, inkl. u-boot-Vollausbau) und der JS-Familie (`docs-check.js`) — mit Vergleichslauf Alt-Tool vs. `d-check`.

## 5. Globale Out-of-Scope-Punkte

- Architektur-/Boundary-Checks (C++-Include-Regeln, Go-Package-Graph, Rust-Crate-Grenzen) — bleiben repo-lokale, sprachspezifische Tools.
- Build-Reproduzierbarkeits-Prüfung (Funktionalität von `repro-check.sh`).
- Mathematik-/Formel-Validierung (MathJax-Rendering-Checks aus `euler-fourier-hilbert`) — Kandidat für eine spätere Version, nicht Teil von 0.x.
- Automatisches Reparieren kaputter Referenzen (Auto-Fix) — das Tool berichtet nur.
- Prüfung von Nicht-Markdown-Formaten (reStructuredText, AsciiDoc, HTML, Jupyter-Notebooks) als eigenständige Dateien. Inline-HTML-Anker *innerhalb* von Markdown sind hiervon ausgenommen und werden geprüft ([`DC-FA-ANCH-001`](#dc-fa-anch-001--heading-anker-validierung-modul-anchors)).
- Rechtschreib-, Stil- oder Markdown-Lint-Prüfungen (dafür existieren etablierte Tools).

## 6. Glossar

| Begriff | Bedeutung im Lastenheft |
|---|---|
| Befund | Eine einzelne festgestellte Regelverletzung mit Datei, Zeile, Ziel und Grund. |
| Regelmodul | Benannte, einzeln aktivierbare Prüf-Einheit (`links`, `anchors`, `ids`, `matrix`, `external`, `codepaths`, `spans`, `hostpaths`, `diagrams`, `versions`, `pins`, `immutable`, `vcs`, `commits`, `planning`, `tracked`, `targets`, `citations`, `sources`, `structure`). |
| Scan-Wurzel | Verzeichnis, unterhalb dessen Markdown-Dateien gesucht werden; zugleich Bezugspunkt der Pfadauflösung. |
| Anker | Fragment-Teil eines Links (`#…`), das auf ein Heading der Zieldatei zeigt (GitHub-Slug-Verfahren). |
| Repo-Escape | Linkziel, dessen aufgelöster Pfad außerhalb der Repository-Wurzel liegt. |
| Kennung | Textuelle ID nach deklariertem Muster (z. B. `ADR-0042`), für die Linkpflicht gelten kann. | <!-- d-check:ignore (Beispiel-ID, fiktiv) -->
| Dokumentklasse | Knoten der Referenzmatrix. **Regelfall:** eine über Pfad-Muster definierte Gruppe von Dokumenten (z. B. Contract-Spec, ADR, Slice). **Token-Ziel:** eine Klasse darf allein ein `token`-Muster tragen, ohne Pfade ([`DC-FA-MTX-003`](#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)); sie hat dann keine Mitglieder, und ihr Gegenstand ist eine Zeichenkette statt eines Dokuments — eine Regel, die sie als **Quelle** nennt, kann deshalb nie feuern. |
| Referenzmatrix | Deklaration, welche Dokumentklasse auf welche verweisen darf, inkl. Status-Bedingungen. |
| Aktives ADR | ADR, dessen Status-Feld keinen verbotenen Wert (`superseded`, `deprecated`) trägt. |
| Quell-Tools | Die dreizehn konsolidierten Alt-Tool-Vorkommen in den Schwester-Repositories des Entwicklungs-Workspace: zwölf aus drei Familien (Shell: `verify-doc-refs.sh`, Python: `check_refs.py`, JavaScript: `docs-check.js`) plus eine eigenständige Python-Linie (`check_markdown_links.py`, Inventur-Nachtrag 2026-06-12). |

## 7. Historie

| Version | Datum | Änderung | Verweis |
|---|---|---|---|
| 0.77.0 | 2026-08-30 | [`DC-FA-CLI-001`](#dc-fa-cli-001--aufruf-und-scan-wurzel) und [`DC-FA-CLI-010`](#dc-fa-cli-010--makefile-fragment-ausgeben): **beide Ausgaben nennen das Benutzerhandbuch mit voller URL** — die Hilfe als eigene Trailer-Zeile, das erzeugte `d-check.mk` im Kopfkommentar. **Anlass ist eine Lücke, die gemessen ist:** beide Ausgaben trugen **null** URLs und verwiesen ausschließlich auf **andere Ausgaben** (`--print-config`/`--suggest-config` bzw. die Release-Notes); wer d-check aus dem Werkzeug heraus kennenlernt, fand den aufgabenorientierten Einstieg nicht. Beim Fragment wiegt das schwerer als bei der Hilfe: es reist in ein **fremdes** Repo, und sein Kopf ist der einzige Ort, an dem ein Zeiger dauerhaft mitfährt. **Die URL zeigt auf den Hauptzweig, nicht auf eine Version, und das ist die Entscheidung:** eine versionierte Form wäre eine **neue Release-Prep-Fläche** — und zwar eine, die grundsätzlich kein Sensor dieses Repos erreichen kann: der Fundort ist eine **Go-Konstante**, und d-check scannt ausschließlich Markdown. Der naheliegende Vergleich mit [`DC-FA-VER-001`](#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in) trägt dabei **nicht**: dessen Muster ist frei konfigurierbar, und dass dieses Repo es nur auf `ghcr`-präfixierte Pins richtet, ist eine Eigenschaft seiner `.d-check.yml`, nicht der Anforderung. **Der Preis ist benannt statt verschwiegen:** wer ein älteres Image fährt, liest ein neueres Handbuch — und eine URL im Werkzeug ist ein Versprechen über ein fremdes System, das **kein** Gate dieses Repos hält (`external` ist strikt opt-in, die Zeichenkette steht zudem unverbunden ein zweites Mal in `packaging/dockerhub/overview.md`). **Kein fremdsprachiger Marker:** anders als `README.md`, das englisch ist und das Handbuch als *(German)* führt, sind Hilfe und Fragment **selbst deutsch** — sie nennen die Sprache trotzdem, nur in eigener Sprache (*„aufgabenorientiert, deutsch"*). Kein neues Verhalten, keine neue Option, kein Grund-Code — zwei Textzeilen und zwei Akzeptanzkriterien | — |
| 0.76.0 | 2026-08-30 | [`DC-FA-STRUCT-001`](#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in): eine Regel darf ihre **Grundmenge erklären** — `tasks-ignore-pattern` nimmt Task-Items aus der `max-tasks`-Zählung, `exempt-section-pattern` nimmt **Abschnitte** aus der Regel. Beide **verkleinern** nur, beide sind opt-in, beide tragen **keinen** neuen Grund-Code; ohne sie ist das Verhalten byte-identisch. **Anlass ist ein eingehender CR** eines Adopters (`docs/plan/cr/2026-08-30-…`), und er gilt hier nachgemessen genauso: `max-tasks: 3` über die **89** DoD-Abschnitte dieses Repos (in 175 Slice-Plänen) liefert **80** `section-oversized` bei **444** Task-Items — der Bestand ist regelkonform, der Zähler misst das Falsche, und die Slice-Vorlage der Baseline liefert den Defekt mit aus. **Zwei Defekte des Antrags sind dabei gemessen und behoben:** das freie Beispiel-Muster nimmt **zwei** echte Liefer-Zusagen mit — und **dasselbe Muster verankert** trifft nur noch ein einziges Item, das ebenfalls eine ist; tragfähig ist erst ein Muster, das **verankert ist und auf Text zielt, den die Bereinigung übrig lässt** (26 Treffer, kein falscher), und seine erste Alternative trifft **null** von 444 Items, weil die Zählung den **bereinigten** Text liest und die Wendung durchgängig in Inline-Code steht — was fence-treu gezählt wird, ist auch fence-treu ignorierbar. **Die beiden Muster sehen deshalb verschiedene Zeichenketten, und das ist die Entscheidung:** `exempt-section-pattern` dieselbe Zeile wie `section-pattern` (samt `#`-Folge, sonst wäre ein analog geschriebenes Muster still wirkungslos), `tasks-ignore-pattern` den **Item-Text** hinter der Checkbox (sonst wäre die verankerte Form unschreibbar). **Die Überdeckung ist sichtbar:** die Meldung nennt die Zahl der ignorierten Items, auch bei null — mit der benannten Grenze, dass sie nur greift, solange die Regel meldet. **Leert `exempt-section-pattern` die Abschnitts-Menge ⇒ `section-missing`** wie bei `exempt-paths`, und es läuft **vor** der Kardinalitäts-Prüfung. Zwei neue Config-Ränder: nicht kompilierendes RE2 je Schlüssel und `tasks-ignore-pattern` **ohne** `max-tasks`. **Der Schlüssel heißt bewusst nicht `exempt-sections`:** `-pattern` bedeutet in diesem Modul RE2, und `exclude-sections` ist in `vcs` und `sources` als **Liste literaler** Überschriften vergeben. Begründung in begleitender ADR | — |
| 0.75.0 | 2026-08-29 | [`DC-FA-CLI-004`](#dc-fa-cli-004--ausgabeformate) und [`DC-FA-STRUCT-001`](#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in): die **Erläuterung eines Befunds** (das Befund-Feld `message`) erreicht den Menschen — als **vierte** tab-getrennte Spalte der Befund-Zeile, **nur wenn gefüllt**, und als eigene Zeile in `--doctor`. **Der Befund bleibt EINE Zeile:** die Zusage *„ein Befund pro Zeile"* und das zählende Akzeptanzkriterium bleiben unberührt, ein Befund ohne Erläuterung ist byte-identisch zu vorher. **Gemessen war der Anlass:** **21 von 23** Regel-Einstiegs-Dateien setzen das Feld (ungedeckt nur `hostpaths` und `spans`), und gerendert wurde es für Menschen nirgends — es erreichte ausschließlich `--json`/`--yaml`. **Die vierte Spalte ist gegen fremden Text abgesichert:** eine Erläuterung stammt aus der Konfiguration oder aus dem geprüften Material (`commits` trägt den Commit-Betreff); ein `hint` mit Tab oder Zeilenumbruch wird **abgewiesen** (Exit 2), und der Reporter ersetzt beides im modul-eigenen Weg durch ein Leerzeichen — sonst wäre die Zusage „ein Befund pro Zeile" über einen Text brechbar, den das Werkzeug nicht kontrolliert. Neu ist außerdem `structure[].hint`: eine Regel **verfasst** ihre Erläuterung selbst und sagt damit, welche **Zusage** sie hütet, während der Grund-Code die **Art** des Defekts sagt; sie gewinnt gegen die modul-eigene Meldung, **außer** bei den **drei** Befunden, die keine Bedingung verletzen (unlesbarer Dateibaum, leer laufende Regel, unlesbare Einzeldatei — dort hat die Regel nicht gemessen). Ein explizit leerer `hint` ⇒ Exit 2. **Abgegrenzt gegen den Fix-Kandidaten** ([`DC-FA-CLI-007`](#dc-fa-cli-007--diagnose-modus)): der ist **abgeleitet** und speist `--repair`, der Hinweis ist **verfasst** und wird nie angewendet. Begründung in begleitender ADR | — |
| 0.74.0 | 2026-08-29 | Neue Anforderung [`DC-FA-WF-001`](#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in) (Modul `workflows`, opt-in) samt neuem Bereichskürzel `WF` in §3: die `uses:`-Referenzen von CI-Workflows tragen **zwei** Deklarations-Zusagen — die fremde Referenz einen vollen 40-stelligen SHA mit Tag-Kommentar (`uses-pin-missing`/`uses-pin-untagged`), die **lokale** ein existierendes Ziel (`uses-local-missing`) **und** einen aufrufenden Job, der die geforderten Rechte führt (`uses-local-perms-undeclared`/`uses-local-perms-narrow`); unlesbares YAML ist ein Befund (`workflow-unparsable`), kein Übersprung. **Die Anforderung entsteht aus einem belegten Ausfall:** ein Job erbte `permissions: {}` vom Workflow-Kopf, sein lokales Ziel verlangte `contents: read`, und der ganze Lauf brach vor dem ersten Job ab — während das damalige Skript-Gate grün meldete. **Sie beantwortet die Grenz-Frage ausdrücklich:** das Modul liest die **Ziele** lokaler Referenzen, die es nicht scannt, und dieselbe Zusage gilt dort. Die Scan-Menge ist **konfigurierbar** (`workflows.dir`), weil der Ort CI-System-spezifisch ist; ohne sie ist das Modul inert. Löst die Mechanik des früheren Harness-Skripts ab. Begründung in begleitender ADR | — |
| 0.73.0 | 2026-08-29 | [`DC-FA-STRUCT-001`](#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in): die **tabellenbezogenen Bedingungen** stehen jetzt unter **einer Klammer** `table` (`order`, `order-column`, `column[]`), und die Zellengrenzen sind eine **Liste je Spalte** statt eines Schlüssel-Tripels je Regel. **Anlass war der eigene Bestand:** vier Spalten desselben Abschnitts kosteten vier Regeln mit viermal wortgleichem `files`/`section` — die Redundanz war die Form, die jede künftige Mehrspalten-Zusage geerbt hätte. Zugleich behoben ist ein **Namens-Defekt**: `cell-max-column` benannte die Spalte und schaltete die Bedingung scharf, trug aber `max` im Namen — auch dort, wo nur eine Untergrenze stand (der eigene Bestand belegte genau das). Die fünf flachen Vorgänger-Schlüssel sind **entfallen** und werden mit dem **neuen Ort** abgewiesen (Exit 2), nicht still ignoriert und nicht mit der generischen Unbekannt-Feld-Meldung des Decoders. **Zwei neue Config-Ränder**, die erst die Liste möglich macht: eine leere `table`-Klammer und derselbe Spaltenname zweimal. **Keine Verhaltensänderung an einem Befund:** Grund-Codes, Zeilen-Zuordnung und die Ziel-Form `… :: Spalte <name>` bleiben, wie sie waren — was sich ändert, ist, wer die Spalte trägt: nicht mehr die Regel-Identität, sondern das Befund-Ziel. Begründung in begleitender ADR | — |
| 0.72.0 | 2026-08-28 | [`DC-FA-STRUCT-001`](#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) um die **Zellenlängen-Bedingung** erweitert (`cell-max-column`/`cell-max-chars`, neunte Bedingung, opt-in; Erweiterung statt neues Kürzel nach dem etablierten Schnitt-Kriterium — Einzelmodul-Frage ⇒ bestehende Anforderung ändern): jede Zelle einer **benannten** Spalte trägt höchstens N **Zeichen**, gemeldet auf **ihrer** Zeile. Zwei neue Grund-Codes `section-cell-oversized`/`section-column-missing`, vier neue fail-closed Config-Ränder. **Die Spalte wird über ihren Kopfzeilen-Namen adressiert, nicht über eine Position:** eine eingefügte Spalte verschöbe eine Positions-Angabe **still**, ein umbenannter Kopf meldet **laut** — dieselbe Wahl, die die achte Bedingung bei den Überschriften traf. Anlass ist ein gemessener Bestand: eine Titel-Spalte mit Median 442 und Maximum 2206 Zeichen neben H1-Titeln von Median 77, wodurch eine als **derivativ** deklarierte Sicht zur zweiten, driftfähigen Quelle wurde. **Ein `forbid-pattern` konnte die Länge ausdrücken und taugte trotzdem nicht:** sein Befund nennt den Abschnitt statt Zeile und Spalte, er hängt an der Form einer Nachbarzelle, und eine Zelle mit einer Pipe entkam ihm — dieselbe Bauform, die dieses Repo zuletzt zugunsten einer `structure`-Regel aufgegeben hat. **Mit-Wirkung auf drei Bestands-Flächen:** die Zell-Zerlegung des Produkts ist jetzt **eine** statt zwei — escape- und backtick-bewusst (`\|` ist ein Zeichen, kein Trenner) —, womit die dazu benannte Grenze der Chronologie-Bedingung und die von `planning.waves` **entfällt** statt vererbt zu werden; kein Grund-Code und keine Befund-Form dieser beiden ändert sich. Zugleich sagt die Befund-Form jetzt zu, dass die drei zeilen-gebundenen Bedingungen auf **der Zeile melden, an der repariert wird**, und nur im Leerlauf auf die Überschrift zurückfallen — das galt seit der siebten Bedingung und stand so nicht da. Begründung in begleitender ADR | — |
| 0.71.0 | 2026-08-27 | [`DC-FA-DIST-002`](#dc-fa-dist-002--docker-hub-spiegel) sagt die Gleichheit jetzt am **Config-Digest** zu statt am **Manifest**-Digest — die Vorfassung war **am Bestand widerlegt**: derselbe lokale Bild-Push liefert je Registry verschiedene Manifest-Digests (gleicher `mediaType`, gleicher Config-Digest, aber neu komprimierte Layer-Blobs; drei von drei geprüften Tag-Paaren des Schwester-Repos). Die fail-closed-Prüfung der Vorfassung hätte damit **jedes** Release gebrochen. Der Config-Digest ist die Identität des Bild-**Inhalts** und registry-stabil; gelesen wird er aus **beiden Registries**, nicht aus dem lokalen Daemon — zwei aus demselben lokalen Bild abgeleitete Werte wären trivial gleich, und ein Vergleich, der nicht fehlschlagen kann, prüft nichts. **Neu ausgesprochen:** der Manifest-Digest ist registry-lokal, ein `docker.io`-Pin nimmt den Docker-Hub-Digest — verschwiegen wäre das eine Falle, deshalb steht es in der Beschreibung und im Out-of-Scope. **Negative-Kriterium geschärft:** der Abbruch erfolgt **vor** dem Spiegel-Push, weil eine `uses:`-Login-Action mit ihrem eigenen Text scheitert und eine nachgelagerte Meldung nie anliefe. **Out-of-Scope präzisiert:** die Hub-Beschreibungsseite wird aus dem Repo gesetzt, ihr Fehlschlag lässt das Release grün — die Prüfung ihres Zeichen-Limits wandert dafür in `make gates`, wo sie beim Schreiben greift statt beim Veröffentlichen. Anlass: unabhängiger Review. Begründung in begleitender ADR | — |
| 0.70.0 | 2026-08-27 | Neue Anforderung [`DC-FA-DIST-002`](#dc-fa-dist-002--docker-hub-spiegel): jedes nach GHCR veröffentlichte Image wird zusätzlich als `docker.io/pt9912/d-check` **gespiegelt** — dasselbe Bild unter denselben Tags, **kein zweiter Bau**. Prüfgröße ist der **Manifest-Digest**; er ist auf beiden Registries gleich, und damit trägt ein `docker.io`-Pin so weit wie ein `ghcr.io`-Pin. GHCR bleibt die **Quelle**, die Richtung ist Teil der Zusage, und die Tagging-Disziplin ist die von [`DC-FA-DIST-001`](#dc-fa-dist-001--docker-image) (volle Semver-Tags, `:latest` nur stabil). **Fail-closed, und das ist die eigentliche Entscheidung:** eine fehlgeschlagene Spiegelung ist ein fehlgeschlagenes Release, auch wenn GHCR bereits trägt — der Preis ist die Bindung an eine fremde Verfügbarkeit, ausdrücklich gewählt gegen die fail-open-Variante, die ein Schwester-Repo fährt. Damit die Teil-Veröffentlichung nicht geraten werden muss, verlangt das Negative-Kriterium, dass die Meldung den bereits veröffentlichten GHCR-Stand **benennt**. Zugleich **geändert**: der Out-of-Scope-Satz von [`DC-FA-DIST-001`](#dc-fa-dist-001--docker-image) schloss *„Distributionswege jenseits GHCR"* aus und hätte diese Anforderung verboten; er ist auf die verbleibenden Wege (Homebrew, Paketmanager, Release-Binaries) zurückgeschnitten, statt stillschweigend überholt zu werden. **Nicht** zugesagt: weitere Registries, ein vom GHCR-Bild abweichender Hub-Bau, der Inhalt der Hub-Beschreibungsseite (er wird aus dem Repo gesetzt, ist aber nicht Teil der Zusage — sein Fehlschlag lässt das Release grün). Anlass: Auftraggeber. Begründung in begleitender ADR | — |
| 0.69.0 | 2026-08-27 | Die **Form** des Zeilen-Markers folgt der **Kommentar-Lexik seiner Eingabe** und ist damit je Konsument verschieden. Bei [`DC-FA-CODE-001`](#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in) und [`DC-FA-ID-001`](#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids) — Eingabe ist Markdown-Prosa — muss er in einem **HTML-Kommentar** stehen; eine **blanke** Erwähnung wirkt nicht mehr. Bei [`DC-FA-DIAG-001`](#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in) bleibt er ein **Token** — das ist eine **ausdrückliche Festlegung der Spezifikation** (in einem `mermaid`-Fence bildet die Diagramm-Sprache den Kommentar; eine Lexik je Fence-Sprache wäre ein Grammatik-Parser) —, und bei [`DC-FA-VER-001`](#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in) ebenso, weil es alle Zeilen sprachgemischt liest. Damit ist die **blanke** Erwähnung in Prosa entschärft; die in Backticks trug schon 0.68.0 über die **Lage**. Die Bedingung ist **monoton verengend** — sie kann keinen neuen stillen Grün-Pfad erzeugen —, aber sie entschärft nicht *jede* Form: eine escapte oder in einem eingerückten Code-Block stehende Erwähnung wirkt weiter. **Im Bestand kostenlos, und der Grund dafür ist eng:** von **66** wirksamen Markern über 558 Dateien tragen **65** die Form bereits; der 66. ist eine Erwähnung in Backticks, die nur durch das Paritäts-Leck wirkt (benannte Grenze der Vorfassung) und in einem Verzeichnis liegt, das beide Konsumenten datei-weit ausnehmen. Der Repo-Lauf **konnte** also nichts finden — der einzige Zähne-Beleg ist die konstruierte Gegenprobe. Einen **baren** wirksamen Marker gibt es im Bestand nicht. **Konservativ:** ein `>` im Kommentar vor dem Marker lässt ihn nicht gelten (verpasster Marker = Falsch-Rot und laut; erfundener = stilles Grün). Neues Akzeptanzkriterium (Boundary) bei **beiden** Anforderungen. **Preis, ausgewiesen:** wer den Marker setzt, muss wissen, für welches Modul. Begründung in begleitender ADR | — |
| 0.68.0 | 2026-08-27 | Der **Zeilen-Marker** `d-check:ignore` wirkt bei [`DC-FA-CODE-001`](#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in) und [`DC-FA-ID-001`](#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids) nur noch **außerhalb** von Fenced-Blöcken **und außerhalb von Inline-Code** — in Backticks ist er eine Erwähnung. Bei [`DC-FA-VER-001`](#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in) und [`DC-FA-DIAG-001`](#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in) bleibt die Erkennung **roh**, als benannte Skopierung: ihre Eingabe ist keine Prosa (alle Zeilen inkl. Fences bzw. die Zeilen **innerhalb** eines Fence), dort ist ein Backtick literaler Inhalt. Anlass ist ein gemessener stiller Grün-Pfad: von **249** Marker-Prosa-Zeilen über **553** getrackte Dateien tragen **183** ihn nur in Inline-Code (**66** wirken), und der Angleich legt **fünf echte** Befunde frei — Zeilen dieses Dokuments, die das Ventil beschreiben und sich dadurch selbst ausnahmen; sie tragen jetzt einen **gesetzten** Marker. **Erweiterung, kein Ersatz** — die Ventil-Wirkung, die Achsen-Präzedenz und alle Grund-Codes bleiben unverändert; kein Test hatte die alte Roh-Lesung je behauptet. Neues Akzeptanzkriterium (Boundary) mit Gegenprobe in beide Richtungen. **Benannte Grenze:** die **Form** des Markers bleibt ungeregelt — er gilt als gesetzt, sobald die Zeichenkette außerhalb von Code steht, auch ohne Kommentar-Klammer (im Bestand trägt genau **einer** der 66 wirksamen Marker die bare Form). Begründung in begleitender ADR | — |
| 0.67.0 | 2026-08-27 | [`DC-FA-CITE-001`](#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in) sagt jetzt zu, **was als Direktive gilt**: sie wirkt nur außerhalb von Fenced-Blöcken **und außerhalb von Inline-Code** — die Syntax in Backticks ist eine Erwähnung —, während der **Zitattext** ausdrücklich **roh** verglichen wird. Anlass ist ein gemessener Zustand: über 544 getrackte Dateien tragen **25** Zeilen einen geöffneten Marker, **24** davon in Inline-Code (20 malformt, weil sie die Syntax dokumentieren), **eine** im Fence, **keine** frei; **null** ist eine produktive Direktive — und weil ein malformter Marker fail-closed den ganzen Lauf abbricht, war das Modul **nicht aktivierbar**. Eine Doku-Konvention wäre kein Ausweg gewesen: 12 der 25 liegen in fünf Dateien, die per Regel unantastbar sind. **Zwei** neue Akzeptanzkriterien (Boundary), je mit Gegenprobe. **Erweiterung, kein Ersatz** — der Zitat-Vergleich, die Fail-closed-**Regel** und alle Grund-Codes bleiben unverändert; die **Menge**, die die Regel trifft, fällt von 25 auf null, und genau das ist der Zweck. Ohne aktives Modul ist jeder Befundsatz weiterhin byte-identisch. **Drei Grenzen benannt:** eine echte Direktive wird verschluckt, sobald eine Code-Spanne desselben **Absatzes** sie umschließt; ein Pfad in Backticks fällt fail-closed; und das Strippen kann eine Direktive auch **erzeugen**. Zusätzlich ausgewiesen: die **Ziel-Seite** trägt keine dieser Zusagen (roh, typunabhängig, außerhalb von `scan.ignore`), und die **Ventil-Direktive** wird weiterhin roh erkannt — dieselbe Frage hat im Produkt derzeit zwei Antworten. Begründung in begleitender ADR | — |
| 0.66.0 | 2026-08-26 | [`DC-FA-MTX-003`](#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix) sagt das **Token-Ziel** zu: eine Klasse darf ein `token`-Muster **ausschließlich** tragen, ohne Pfade — sie hat dann keine Mitglieder, ihr Gegenstand ist eine Zeichenkette, und als **Quelle** einer Regel kann sie nie feuern. Zwei Akzeptanzkriterien decken beide Richtungen. Die Glossar-Definition der **Dokumentklasse** ist entsprechend **geweitet** (Regelfall Pfad-Muster, dazu das Token-Ziel) statt einen Gegenbegriff danebenzustellen: so bleibt jede Bestandszeile wahr, die von Dokumentklassen spricht. **Kein Verhaltens-Delta:** die Fähigkeit bestand in der Umsetzung, war aber nicht zugesagt | — |
| 0.65.4 | 2026-08-25 | **Form eines Verweises, keine Anforderungs-Änderung.** Eine Historie-Zeile nannte eine Kennung der zeitlichen Schicht und war damit ein Abwärtsverweis, den §3.4 verbietet; die berichtete Messung bleibt unverändert, der Zeiger auf die zeitliche Schicht entfällt ersatzlos. Sichtbar wurde sie erst, als die Referenzmatrix eine Klasse für diese Schicht bekam — bis dahin fehlte schlicht der Wächter, der sie melden konnte. **Kein Protokoll-Eingriff im Sinne des Kanons:** geändert ist die **Form** eines Verweises, nicht die berichtete Tatsache | — |
| 0.65.3 | 2026-08-23 | **Form der Historie, keine Anforderungs-Änderung.** Die Tabelle trägt die kanonische vierte Spalte `Verweis`; alle Bestandszeilen tragen `—`, weil ohne begonnene CR-Pflicht kein externer Vorgang existiert, den sie nennen könnte. Der Kopf deklariert, dass dieses Repo Bump und Historie **schon vor `Accepted`** führt, obwohl der Kanon vor diesem Status *keine* verlangt — als Adaption in [`MR-032`](../harness/conventions.md#mr-032) geführt. **Zugleich Tatsachenberichtigung:** die Zeilen 0.65.1 und 0.65.2 waren Berichtigungen im Sinne des Kanons, trugen den verlangten Ausweis aber nicht — nachgetragen. Die Spezifikations-Historie bleibt bei zwei Spalten: eine `Verweis`-Spalte trägt nach dem Kanon nur, **was sonst nirgends im Repo steht** — beim Vertrag der externe Change Request, der kein anderes Zuhause hat, während die Technik ihre Aufwärts-Bezüge bereits im Körper verankert; dieselbe Kopplung ein zweites Mal in der Historie zu führen erzeugt keine Information, sondern eine zweite Fassung, die driftet | — |
| 0.65.2 | 2026-08-23 | Nachzug nach unabhängigem Review: **Tatsachenberichtigung** (als solche ausgewiesen,  keine Anforderungs-Änderung.** Die Achsen-Abgrenzung des geteilten Referenz-Ventils ([`DC-FA-REF-001`](#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)) wies die vier marker-tragenden Module in 0.65.1 als **benannte Liste** aus, nannte aber keinen Gegen-Beleg. Sie nennt jetzt die drei Module, die ebenfalls auf Zeilen melden und den Marker **nicht** tragen — `matrix`, `structure` und [`DC-FA-CITE-001`](#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in). Ohne diesen Gegen-Beleg liest sich die Liste weiter wie ein herleitbares Kriterium, und genau das ist sie nicht. Kein Grund-Code, kein Schema, kein Befundsatz betroffen | — |
| 0.65.1 | 2026-08-23 | Nachzug nach unabhängigem Review, vor dem Release: **Tatsachenberichtigung** (als solche ausgewiesen,  keine Anforderungs-Änderung.** Die Achsen-Abgrenzung des geteilten Referenz-Ventils ([`DC-FA-REF-001`](#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)) führte den Zeilen-Marker `d-check:ignore` als „Zeile, **nur** `codepaths`“ — die Aussage war seit 0.8.0 falsch ([`DC-FA-ID-001`](#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids) honoriert ihn seither, [`DC-FA-VER-001`](#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in) seit seiner Einführung, [`DC-FA-DIAG-001`](#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in) seit 0.65.0). Sie nennt jetzt alle **vier** Module und sagt ausdrücklich, dass die Menge eine **benannte Liste** ist und **kein** ableitbares Kriterium: `matrix` und `structure` konfigurieren ebenfalls eigene Muster und melden ebenfalls auf Zeilen, tragen den Marker aber nicht. Anlass war die Nutzer-Doku-Korrektur desselben Satzes — der Spiegel im ranghöheren Stratum wäre sonst stehengeblieben und hätte per Source Precedence gewonnen | — |
| 0.65.0 | 2026-08-22 | [`DC-FA-DIAG-001`](#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in) um die **beiden Ventile** erweitert (`diagrams.exempt-paths` und der Zeilen-Marker `d-check:ignore`, opt-in): das Modul war das einzige mit **konfigurierbaren Mustern**, das **kein** Ventil trug (Module ohne eigene Muster — `hostpaths`, `pins`, `spans` — tragen weiterhin keines; offene Fläche, keine Zusage) — ein Beispiel-Diagramm mit erfundener Kennung war nur durch Herausnehmen ganzer Dateibäume aus dem Scan-Bereich loszuwerden, eine Vermeidung statt einer Ausnahme (am eigenen Profil gemessen). **Zwei Festlegungen aus der Fence-Natur:** der Marker ist ein **Token**, kein HTML-Kommentar — wie der Autor ihn vor dem Renderer versteckt, ist Sache der Diagramm-Sprache (Mermaid: `%%`) —, und er wirkt auf einer Diagramm-Zeile für **diese Zeile**, auf der **Öffnungszeile** für den **ganzen** Block; ohne die zweite Stelle wäre die intuitive Platzierung wirkungslos. Kein neuer Grund-Code. **Ein fail-closed-Rand ist zugleich enger geworden:** ein **Verzeichnis** als `defined-in` bricht jetzt (Exit 2) statt eine leere Definitionsmenge zu liefern und jede Kennung als undefiniert zu melden. Begründung in begleitender ADR | — |
| 0.64.0 | 2026-08-22 | [`DC-FA-STRUCT-001`](#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) um die **Überschriften-Bedingung** erweitert (`headings-match`/`headings-level`, achte Bedingung, opt-in; Erweiterung statt neues Kürzel nach dem etablierten Schnitt-Kriterium — Einzelmodul-Frage ⇒ bestehende Anforderung ändern): *jede* Überschrift des Abschnitts genügt einem Muster, **positiv** formuliert und **je Überschrift** gemeldet (eigener Grund-Code `section-heading-mismatch`, Befund auf der Zeile der Überschrift — dort ist die Reparatur). Anlass ist eine gemessene Ersatz-Konstruktion: dieselbe Aussage war nur als ausgeschriebene Präfix-Negation über den Abschnitts-Text formulierbar (RE2 kennt keinen Lookahead), und die sprach die Heading-Lexik des Moduls nicht — eine **eingerückte** Sektion entkam still. Die Bedingung ist als **zweite** nicht auf dem bereinigten Abschnitts-Text definiert, sondern auf den Überschriften selbst, mit **derselben** Erkennung, die den Abschnitt findet. Zwei Grenzen benannt statt zugesagt: ohne Überschrift der geprüften Ebene ist sie wirkungslos, und ein `headings-level` flacher als der Abschnitt kann in ihm nicht vorkommen. **Der Schlüssel heißt bewusst nicht `heading-pattern`:** unter `planning.closure` trägt dieser Name bereits einen **Selektor** („welcher Abschnitt"), hier wäre es eine **Bedingung** („welche Form") — und beide Blöcke leben im selben Profil. Begründung in begleitender ADR | — |
| 0.63.1 | 2026-08-22 | Nachzug nach unabhängigem Review, vor dem Release: die Zusage „zwei Paare, zwei Befunde" war **nicht haltbar** — die Befund-Adresse (Datei, Zeile, Regel, `target`, Grund-Code) unterscheidet zwei Befunde an derselben Stelle nicht, und die geteilte Nachrunde verwarf den zweiten samt seiner Erwartung (gemessenes stilles Falsch-Negativ). Jetzt gilt: **eine Befund-Adresse, alle Erwartungen** — die Nachricht nennt jede Erwartung mit ihrer Quelle, in Deklarationsreihenfolge; die **Ausgabe**-Reihenfolge ist die geteilte Sortierung, nicht die Deklarationsreihenfolge (die alte Aussage war als beobachtbare Eigenschaft falsch). Ferner: die Mischform-Erkennung hängt an der **Anwesenheit** des Schlüssels, nicht an seinem Wert (`pin-pattern: ""` neben `patterns` ist eine Mischform), mit der Grenze eines wertlosen Schlüssels; und ab zwei Paaren benennt jede Meldung die Fundstelle, bei genau einem Paar bleibt der Wortlaut byte-identisch | — |
| 0.63.0 | 2026-08-22 | [`DC-FA-VER-001`](#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in) um **mehrere Muster-Quellen-Paare** erweitert (`versions.patterns`, opt-in; Erweiterung statt neues Kürzel nach dem etablierten Schnitt-Kriterium — Einzelmodul-Frage ⇒ bestehende Anforderung ändern): ein Repo führt mehrere unabhängige Versions-Reihen, das Modul kannte bisher genau **eine**. Anlass ist die Beobachtung BEO-008 bei Zähler 3, deren benannte mechanische Form (Baseline-Tag in URLs und Prosa gegen den §Baseline-Pin, **zusätzlich** zum Image-Pin gegen das Release-Register) mit einem einzigen `pin-pattern` nicht baubar war. Jedes Paar trägt eigene Quelle und eigene Ausnahmen, die Selbst-Ausnahme der Quell-Datei ist **paar-lokal**; die Kurzform bleibt gültig und **ist** die einelementige Liste (ein Auswertungspfad), beide Schreibweisen zugleich sind fail-closed. Kein neuer Grund-Code, keine Änderung an der Befund-Form. Begründung in begleitender ADR | — |
| 0.62.1 | 2026-08-22 | Präzisierung ohne Verhaltensänderung (Review-Hinweis zur 0.62.0-Erweiterung): [`DC-FA-PLAN-001`](#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) §Wellen-Invariante (Zeile 1/2 und der Absatz „Zwei Kardinalitäts-Modelle") nennt jetzt die Vergleichsgröße des Prädikats in **beiden** Modi — **Kennungen**, nicht Dateien: zwei flache Wellendokumente derselben Kennung sind ein Element (so zählt die Fähigkeit seit 0.59.0, die Spezifikation sagte es in W2 bereits; der Lastenheft-Wortlaut las sich als Datei-Menge). Neues Akzeptanzkriterium **Wellen-Boundary (gleiche Kennung, beide Modi)** pinnt das Ist-Verhalten; ob ein Doppel-Dokument selbst meldepflichtig sein sollte, bleibt ausdrücklich offen (eigener Change Request). Kein Release-Anlass | — |
| 0.62.0 | 2026-08-21 | **CR des Konsumenten ai-harness-course** („planning.waves: Bijektion statt Singleton"): [`DC-FA-PLAN-001`](#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) §Wellen-Invariante bekommt den opt-in Kardinalitäts-Modus `planning.waves.mode: one`\|`many` (Default `one`; ohne den Schlüssel ist der Befundsatz byte-identisch, [`DC-QA-02`](#dc-qa-02--determinismus)). **Anlass:** die Baseline v5.7.0 des Adopters (Kurs-Welle 81) schärft §Offene Wellen auf zwei unabhängige Aussagen — der Ruhe-Marker folgt dem **Anspruch** und steht **zusätzlich** zur Liste, nicht an ihrer Stelle; damit ist der Marker kein Stellvertreter für „läuft eine Welle" mehr, und die `one`-Kopplung (Marker gegen Datei-Zahl) widerspricht dem Modell, das sie stützen soll. **Messung statt Herleitung:** Replay-Harness des Konsumenten, 11/11 PASS — der Singleton beißt bei zwei gelisteten offenen Wellen (s04b), die Marker-Hälfte hält beidseitig über `planning-drift` (s04c/s04d) und bleibt unberührt. **Unter `many`** vergleicht `wave-drift` **Kennungs-Mengen** (dasselbe Verfahren wie in den Registern: literales Glob-Präfix + Ziffernfolge, zeilenweise über die Prosa-Zeilen des Blocks — Fence-Inhalte zählen nicht, Mehrfachnennung zählt einmal, layout-agnostisch für Tabellen- wie Listen-Form), beide Richtungen, jede Kardinalität einschließlich null; das Befund-`target` ist die betroffene Kennung, **kein** zweiter Grund-Code (die Reparatur ist dieselbe, die Richtungen unterscheidet das `target` wie bei den Register-Aussagen). **Nicht-Ziele (CR §6):** keine Änderung an `planning-drift`, keine Default-Änderung, keine Festlegung, ob der Aktiv-Block Tabelle oder Liste führt. **fail-closed:** unbekannter oder explizit leerer Modus ⇒ Exit 2 mit Schlüssel-Nennung (Zeiger-Disziplin der übrigen `waves`-Schlüssel, 0.59.1). **Schnitt:** Einzelmodul-Frage ⇒ bestehende Anforderung erweitert statt neues Kürzel (Kriterium seit 0.47.0, zuletzt 0.60.0) | — |
| 0.61.1 | 2026-08-21 | Nachzug nach unabhängigem Review, vor dem Release. **Die Typ-Mischungs-Semantik ist auf genau eine Lesart gepinnt** (der MEDIUM-Befund): Paar-Lesart mit Anker-Reset — die Misch-Zelle meldet sich selbst, die gesunde Folge-Zeile dahinter meldet **nicht**; das Akzeptanzkriterium sagt es jetzt ausdrücklich, der Kaskaden-Fall ist als Test gepinnt. Ferner: ein Versions-Segment außerhalb des Zahlbereichs ist **untypisierbar** statt still kleinstmögliche Version (der LOW-Befund), und die Ausdrucks-Grenze **„eine Chronologie-Zusage je Abschnitt"** ist als Out-of-Scope benannt (die Regel-Identität trägt keine Spalte) | — |
| 0.61.0 | 2026-08-21 | [`DC-FA-STRUCT-001`](#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) um die **Chronologie-Monotonie** erweitert (`table-order`/`table-column`, siebte Bedingung; Erweiterung statt neues Kürzel nach dem etablierten Schnitt-Kriterium — Einzelmodul-Frage ⇒ bestehende Anforderung ändern): die Schlüsselspalte jeder zusammenhängenden Tabelle des Abschnitts wird **typisiert** (geschlossene Menge: ISO-Datum, Punkt-Version — segmentweise numerisch, ein String-Vergleich meldet gemessen drei korrekt sortierte Bestandstabellen rot) und nicht-strikt monoton in der je Regel konfigurierten Richtung verglichen. Zwei neue Grund-Codes: `section-unordered` (Bruch-Zeile; auch Leerlauf ohne Datenzeile) und `section-cell-untyped` (untypisierbare Zelle/Typ-Mischung — Befund statt stillem Übersprung). **Erste benannte Roh-Lese-Ausnahme** von der Abschnitts-Bereinigung: reale Schlüsselspalten stehen in Inline-Code. Drei neue fail-closed Config-Ränder. Begründung in begleitender ADR | — |
| 0.60.2 | 2026-08-16 | Nachzug der bestätigenden Re-Review (Text-Auflagen, keine Verhaltensänderung): die **Scan-Bereichs-Kopplung** ist als Grenze benannt (Gruppen-Orte müssen im wirksamen Scan-Bereich liegen — eine nie gescannte Datei ist still keine Quelle), die Ventil-Aufzählung von [`DC-FA-REF-001`](#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus) zählt spezifikationsseitig alle **fünf** unterdrückten Codes, und die letzten beiden Flächen des „zeichengenau"-Überclaims (CHANGELOG, Config-Kommentar) sind auf die 15/19-Aussage gezogen | — |
| 0.60.1 | 2026-08-16 | Nachzug nach unabhängigem Review **und** einem CI-Realfund, vor dem Release. **Review:** die fail-closed-Klasse zum dritten Mal (ein `resolve-from`-Verzeichnis mit Tippfehler schaltete die Quellen-Rolle still ab), die Ist-Ort-Vorbedingung stand nur im Code, der Ventil-Wortlaut von [`DC-FA-REF-001`](#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus) widersprach der neuen Wirkung aktiv, und „zeichengenau“ war ein Überclaim — die Retro-19 überlappen den realen Bruch nur zu 15/19; die übrigen vier waren **Ziel**-Wanderungen, jetzt als Grenze benannt. **CI-Realfund:** der erste fail-closed-Zuschnitt meldete auf jedem frischen Klon das legitim **geleerte** `open/`-Verzeichnis — git überträgt leere Verzeichnisse nicht, ein einzelner fehlender Ort ist von einem Tippfehler nicht unterscheidbar. Jetzt meldet fail-closed, wenn **kein** `dirs`-Ort existiert oder ein Ort als Datei existiert; die Rest-Grenze ist benannt. **Und die Beschreibung selbst fehlte:** ein abbrechender Batch-Editor hatte beim 0.60.0-Schnitt nur Historie und Akzeptanzkriterien geschrieben — der Anforderungs-Text ist mit dieser Version nachgeliefert | — |
| 0.60.0 | 2026-08-16 | [`DC-FA-LINK-001`](#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links) um **ortsfeste Verweise** erweitert (`links.resolve-from`, opt-in; Change Request des Konsumenten a-check, Erweiterung statt neues Kürzel nach dem etablierten Schnitt-Kriterium (Einzelmodul-Frage ⇒ bestehende Anforderung ändern)): wo Dateien zwischen Geschwister-Verzeichnissen wandern, muss ein relativer Verweis von **jedem** Ort der Gruppe auflösen — und überall auf **dasselbe** Ziel. Eigener Grund-Code `link-position-dependent` (die Reparatur ist Präfixieren, nicht Ziel-Anlegen). **Zwei Festlegungen aus der Bestandsmessung:** Quellen sind nur die **wandernden** Verzeichnisse (`dirs`) — über alle vier Lifecycle-Orte gerechnet wären heute 108 Befunde Falsch-Positive auf ortsfesten Ruheort-Dokumenten (mit der Einschränkung: 0 von 79); und der **Retro-Beleg** reproduziert einen historischen 19-Link-Bruch zeichengenau. Begründung in begleitender ADR | — |
| 0.59.1 | 2026-08-16 | Nachzug nach unabhängigem Review, vor dem Release. **Der blockierende Befund traf die Motivations-Richtung des Slice selbst:** ein unlesbares `waves.dir` schaltete die Fähigkeit im Ruhe-Zustand **still** ab — mit einem Pfad-Tippfehler wäre genau die zweimal real eingetretene Aussage-2-Verletzung dauerhaft unsichtbar gewesen; jetzt fail-closed über `wave-drift`, ebenso eine **fehlende Register-Überschrift** (die Aktivierung ist die Behauptung, dass die Roadmap beide Register führt). Ferner: die Wellen-Kennung ist ans **literale Glob-Präfix** gebunden statt an ein hartes „welle-“ (beide Globs müssen dasselbe Präfix tragen, Exit 2 sonst), explizit leere `waves`-Schlüssel brechen mit Exit 2 statt still auf den Default zu fallen (Zeiger-Disziplin wie bei `closure.glob`), und das `target` von `wave-unregistered` ist die **Ergebnisnotiz** statt der Kennung (Ventil-Parität). **Und eine eigene Messaussage war falsch:** der zwölfte Befund des Schwester-Repo-Laufs war ein Artefakt des Default-Markers der Probe-Konfiguration, keine Bestands-Verletzung — die elf fehlenden Ergebnisnotizen sind robust | — |
| 0.59.0 | 2026-08-16 | [`DC-FA-PLAN-001`](#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) bekommt eine **dritte Fähigkeit**: die Lifecycle-Invariante eine Ebene höher, die **Wellen**-Abschnitte der Roadmap gegen die Wellen-Dateien (`planning.waves`, opt-in innerhalb des opt-in Moduls; vier Aussagen, vier Grund-Codes). **Zwei Entwurfsannahmen hat die Bestandsmessung über drei Planungs-Bäume widerlegt:** das verpflichtende Artefakt einer geschlossenen Welle ist die **Ergebnisnotiz**, nicht das Plan-Dokument (gegen dieses geprüft meldet die Aussage 19-mal, weil ältere Wellen vor der Wellendokument-Konvention geschlossen wurden), und die Vorschau-Aussage greift nur auf der **Welle-Spalte** und nur bei Kennungen (zwei der drei Bäume schreiben dort Namen; die Trigger-Spalte darf andere Wellen nennen). Dazu die Gegenrichtung „Ergebnisnotiz ohne Register-Zeile“, die im eigenen Bestand **dreimal eingetreten** war. Begründung in begleitender ADR | — |
| 0.58.3 | 2026-08-16 | Nachzug nach **dritter** Review-Runde. **Die Klasse war nicht geschlossen:** drei Module, die in keinem der beiden Vorgänger-Reports vorkamen, beantworteten eine Lexik-Frage roh — [`DC-FA-VCS-001`](#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in) las `immutable-when` und die Kopf-Status-Zeile über **alle** Zeilen (eine Datei, die ihren Kopf als Beispiel zeigt, galt als immutabel bzw. verschob die gestrippte Zeile ⇒ **stilles Grün** im Immutabilitäts-Gate), [`DC-FA-PLAN-001`](#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) beendete den Aktiv-Block an einem rohen `## `-Präfix (eingerückte H2, tab-getrennte H2 und H1 wurden übersehen, `planning-drift` entfiel **still**), und [`DC-FA-TGT-001`](#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in) las Tabellenzeilen roh (ein Beispiel-Block ließ ein undokumentiertes Target als dokumentiert gelten). Alle drei repariert, alle drei mit einer Gegenprobe bewacht. **Die Richtungs-Aufzählung ist jetzt ausdrücklich offen** — in drei Runden war sie dreimal unvollständig; ferner trägt die Anker-Menge ihre **Grenze** jetzt in der Anforderung, nicht nur in der ADR | — |
| 0.58.2 | 2026-08-16 | Nachzug nach **bestätigender Re-Review**, vor dem Release — beide blockierenden Befunde der Vorrunde waren nur **teilweise** geheilt, und zwar in derselben Bauform wie zuvor. (1) Die Anker-Vereinheitlichung galt nur der **HTML**-Hälfte; die **Slug**-Hälfte divergierte weiter (Duplikat-Slug `#alt-1`, prozent-kodiertes Fragment `#a%20b`) — jetzt eine Antwort für alle drei Module, mit Akzeptanzkriterien bei [`DC-FA-VER-001`](#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in) und [`DC-FA-PIN-001`](#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in). (2) Die Richtungs-Aufzählung war erneut **geschlossen** formuliert und übersah eine dritte, bei `pins` stille Stelle: die Auflösung ist jetzt case-sensitiv. Sie ist jetzt offen formuliert. **Ferner:** die Absatz-Zusage aus 0.58.1 war gegen einen Lauf **falsch** — eine Leerzeile trennt gerade nicht, die Direktive paart über Leerzeilen hinweg; das Modul `planning` beantwortete „ist das eine Überschrift“ zweimal roh und meldete dadurch zwei **Falsch-Rot**; und die permissive Anker-Menge ist bei ihrem neuen Konsumenten `versions` **nicht** folgenlos (Anker in HTML-Kommentar oder eingerücktem Block ⇒ Falsch-Rot) — als Grenze benannt, nicht mitgenommen | — |
| 0.58.1 | 2026-08-16 | Nachzug nach unabhängigem Review, vor dem Release. **Zwei blockierende Befunde:** (1) die Anker-Erkennung war nur zur **Fence-Hälfte** vereinheitlicht — `versions`/`pins` behielten ihre eigene Regex und hielten vier Zeichenfolgen für Anker, die [`DC-FA-ANCH-001`](#dc-fa-anch-001--heading-anker-validierung-modul-anchors) nie als solche gelesen hat (Anker in Inline-Code, `data-id`, `name` an beliebigem Element, anker-förmige Prosa ohne Tag); gemessen trug die reparierte Achse **1** Vorkommen, die stehengebliebene **40**. Jetzt **eine** Erkennung für alle drei Module. (2) Die Richtungs-Zusage „findet mehr, und weniger an keiner Stelle“ war **falsch** — zwei Fälle verlieren einen Befund, einer davon still. Beide Richtungen stehen jetzt in Anforderung, ADR und Release-Notiz. Ferner: die Zusagen stehen jetzt in den **Anforderungen** selbst statt nur in Algorithmus und ADR (der Widerspruch zur Fence-Ausnahme von `versions` entstand genau dort), der Fenced-Block ist als **Blockquote-Terminator** benannt, und die `vcs`-Grenze beschreibt die **stille** Richtung samt eines beobachtbaren Re-Evaluierungs-Triggers | — |
| 0.58.0 | 2026-08-16 | Drei Konsumenten der geteilten Lexik: [`DC-FA-CITE-001`](#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in) trennt am **Fenced-Block** wie am Leerzeile-Absatz — ein Code-Block zwischen Direktive und Zitat führt in den zugesagten fail-closed-Abbruch statt in eine Zufalls-Paarung; [`DC-FA-VER-001`](#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in) und [`DC-FA-PIN-001`](#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in) erkennen **beide Anker-Formen nur außerhalb von Fences**, während ihr roher Span unberührt bleibt (die gescopten Roh-Lesungen sind eine andere **Frage**, keine andere **Antwort**). [`DC-FA-VCS-001`](#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in) bekommt die **benannte Grenze** der Revisions-Achse statt eines zweiten Wächters. **Die Änderung wirkt in beide Richtungen:** sie **findet mehr** (ein grüner Konsumentenlauf kann fail-closed werden, wo ein Fenced-Block trennt) und **findet weniger** an zwei konstruierbaren Stellen — ein von einem Fence unterbrochenes Blockquote meldet kein `citation-mismatch` mehr, und ein `dpin` auf einen Anker im Fence verliert seinen Drift-Schutz **kommentarlos** (das Modul schweigt zu unauflösbaren Zielen). Am eigenen Bestand gemessen ist heute **kein** Fall betroffen (0 wirksame Anker, 0 von 152 Revisions-Blobs). Begründung in begleitender ADR | — |
| 0.57.2 | 2026-08-15 | Nachzug nach bestätigender Re-Review: die Aufnahme-Klausel von [`DC-FA-CLI-006`](#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten) nennt sich eine „geschlossene Menge“, klassifizierte aber drei Module gar nicht — `citations`, `sources` und das neue `structure`. Das Akzeptanzkriterium derselben Anforderung tat es; die Klausel jetzt auch | — |
| 0.57.1 | 2026-08-15 | Nachzug nach unabhängigem Review, vor dem Release: die **Akzeptanzkriterien** von [`DC-FA-CLI-010`](#dc-fa-cli-010--makefile-fragment-ausgeben) führten weiter elf Targets ohne `doc-structure` — Beschreibung und Out-of-Scope waren nachgezogen, die Kriterien nicht. Dieselbe Stelle war schon in 0.37.1 als Selbstwiderspruch saniert worden. Ferner nennt [`DC-FA-CLI-006`](#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten) `structure` jetzt in der Enumeration der nicht aktivierten situativen Module | — |
| 0.57.0 | 2026-08-15 | Umsetzung von [`DC-FA-STRUCT-001`](#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) — das **20. Regelmodul** `structure` liegt vor, samt der von [`DC-FA-PLAN-001`](#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) seit 0.51.0 zugesagten, aber nicht implementierten Mehrdeutigkeits-Härte (`closure-note-ambiguous`). Die Abschnitts-Mechanik ist **geteilt**, nicht kopiert: die Closure-Fähigkeit ist ein Preset derselben Semantik, und ein Kopplungs-Test fährt dieselbe Eingabe durch beide Oberflächen. Neun neue Grund-Codes. **`DC-FA-CLI-010`-Erweiterung (11 → 12 Targets):** `--print-mk` trägt ein `doc-structure`-Target | — |
| 0.56.1 | 2026-08-10 | Nachzug nach unabhängigem Review, vor dem Release: drei Vertragsflächen standen noch auf der Fassung **vor** der Zähl-Angleichung und waren gegen die Umsetzung falsifizierbar — die Aussage, die Zählung sehe Inline-Code weiterhin (sie tut es nicht mehr), und das Akzeptanzkriterium zur Floskel, das einen **Teilstring**-Treffer verlangte (der Lauf meldet nichts). Ferner stand die ASCII-Wortgrenze nur in ihrer günstigen Richtung; sie steht jetzt in beiden: ein Umlaut ist **kein** Wortzeichen, eine Phrase mit angrenzendem Umlaut gilt damit als grenzständig und trifft | — |
| 0.56.0 | 2026-08-10 | Change Request des Auftraggebers: die Floskel-Bedingung von [`DC-FA-PLAN-001`](#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) vergleicht an **Wortgrenzen** statt als Teilstring. Anlass: kurze Phrasen waren unbrauchbar, und kurze Phrasen sind genau die, für die die Bedingung gedacht ist — eine Notiz, die nur „Ok.“ sagt, ist der Kern-Fall. Am eigenen Bestand gemessen (97 Notizen): `ok` fällt von **68** Treffern auf **1**, `n/a` von 2 auf 0, `fertig` von 3 auf 0; mehrwortige Phrasen sind verhaltensgleich, und die fünf **konfigurierten** Phrasen ändern sich nicht. Die Änderung **findet weniger** — ein roter Konsumentenlauf kann grün werden; sie gehört damit in dieselbe Release-Notiz wie die Lockerung aus 0.55.0. Wortgrenzen machen kurze Phrasen brauchbar, nicht automatisch sicher: der eine verbleibende `ok`-Treffer ist echt | — |
| 0.55.0 | 2026-08-10 | Change Request (Parität zum abzulösenden Adopter-Skript des Schwester-Repos): die **Substanz-Zählung** von [`DC-FA-PLAN-001`](#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) ist angeglichen — der Abschnitt wird **einmal** bereinigt (Fences **und** Inline-Code), und ein Satzende zählt nur vor **Whitespace oder Zeilenende**. Anlass: eine Notiz, die der Adopter-Sensor als zu dünn meldet, lief bei d-check durch. Gemessen am eigenen Bestand tragen **mehr als die Hälfte** aller Satzende-Vorkommen keinen Satz (Link-Pfade, Versionsnummern); das Minimum fällt von 7 auf 5, bei Schwelle 4 bleibt der Bestand grün. **Die Änderung wirkt in zwei Richtungen:** `closure-note-thin` wird **schärfer** (ein grüner Konsumentenlauf kann rot werden), `closure-note-boilerplate` **lockerer** — eine Floskel in Backticks trifft nicht mehr, weil beide Bedingungen denselben Text lesen. Die Lockerung ist begleitend per ADR entschieden ([`AGENTS.md` §3.6](../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden)) und gehört in die Release-Notiz: wer eine Floskel zitiert stehen hat, verliert einen bestehenden Befund | — |
| 0.54.1 | 2026-08-10 | Nachzug nach unabhängigem Review, vor dem Release der 0.54.0: die Zusage „Vergleichszeichen sind durch die Form ausgeschlossen“ hielt nur für die Schreibweise **mit** Leerzeichen — `Die Latenz blieb <1 s und der Recall >0,9` meldete. Das Innere eines Platzhalters muss jetzt **frei von Whitespace** sein (ein Feldname, kein Satz); das erledigt zugleich HTML-Tags mit Attributen. Neu ausgeschlossen: das Winkelklammer-**Linkziel** `](<ziel>)`. Der Substanz-Bullet behauptete weiter, ein zurückgelassener Platzhalter falle der Zählung auf — er tut es nur, wenn er kurz ist; ein vollständiger Vorlagen-Rumpf erreicht die Schwelle, und genau deshalb gibt es die vierte Bedingung. Zwei Grenzen sind jetzt **benannt**: der eingerückte Code-Block (in d-check nirgends modelliert) und die ungerade Backtick-Parität im Absatz | — |
| 0.54.0 | 2026-08-10 | Change Request (Konsument `ai-harness-course`, CR 2): [`DC-FA-PLAN-001`](#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) bekommt eine **vierte**, opt-in Struktur-Bedingung — `closure-note-placeholder` (Schalter `planning.closure.placeholder`, Default `false`) meldet den unausgefüllten Rumpf einer Vorlage. Anlass: ein Template-Rumpf ist syntaktisch vollständig und passiert alle drei bestehenden Bedingungen (gemessen gegen v0.52.0: 0 Befunde). Erkannt wird die Auszeichnungs-Form; drei Einschränkungen halten sie eng, davon die tragende: **Inline-Code zählt nicht** — am eigenen Bestand gemessen (96 Closure-Notizen) liegen **alle zwölf** Treffer in Inline-Code und **null** außerhalb, ohne die Einschränkung wäre jeder ein Falsch-Positiv. Autolinks/Adressen und HTML-Tags fallen per Nachfilter raus; Vergleichszeichen und Generics bereits durch die Form. Erster Treffer je Kandidat, wie bei der Floskel. **Die Substanz-Zählung bleibt unberührt** — ihre Angleichung an dieselbe engere Sicht ist eine eigene Frage, weil sie eine ausgelieferte Schwelle bewegt | — |
| 0.53.1 | 2026-08-10 | Nachzug nach unabhängigem Review: das Akzeptanzkriterium zur Nullmengen-Härte stand **zweimal** in der Anforderung — die ältere Fassung nannte weiter `planning.slice-glob` und war gegen die Umsetzung falsifizierbar (Closure-Verzeichnis nur mit Wellen-Dateien plus gesetztem `closure.glob` ⇒ das Kriterium verlangt einen Befund, der Lauf meldet Exit 0). Beide zu **einem** Kriterium über den **effektiven** Filter zusammengezogen | — |
| 0.53.0 | 2026-08-10 | Change Request (Konsument `ai-harness-course`, CR 1): [`DC-FA-PLAN-001`](#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) bekommt mit `planning.closure.glob` einen **eigenen** Kandidaten-Filter für die Closure-Fähigkeit; Default ist ein **Verweis** auf `planning.slice-glob`, nicht ein wiederholtes Literal — ohne den Schlüssel ist der Befundsatz byte-identisch. Anlass: die beiden Fähigkeiten stellen verschiedene Fragen mit verschiedenen Grundmengen („liegt hier noch Arbeit?“ gegen „ist jedes abgeschlossene Paket dokumentiert?“) und teilten sich einen Schlüssel, solange beide zufällig dasselbe trafen. Wer die Menge weitet, verbiegt die jeweils andere — im Grenzfall zählt die Roadmap-Datei selbst als Slice und der Ruhe-Marker meldet dauerhaft falsch-rot (gemessen gegen v0.52.0). Ein **explizit** leerer oder ungültiger Glob bricht am Config-Rand mit Exit 2 ab statt still auf den Default zurückzufallen | — |
| 0.52.3 | 2026-08-10 | Nachzug nach dritter Review-Runde. **Konsumenten-relevant:** die Heilung der Fence-Trimmung hat auch das mit 0.52.0 ausgelieferte Closure-Gate ([`DC-FA-PLAN-001`](#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)) geändert — die Richtung ist **nicht nur** „findet mehr": hinter einer mit Unicode-Whitespace eingerückten Fence-Zeile verstummen `closure-note-thin` und `closure-note-missing`, während `closure-note-boilerplate` neu meldet. Fachlich war das immer die zugesagte Lexik (Spezifikation Schritt C4: Fence-Lexik wie im übrigen Scanner); gemessen wurde bis dahin eine andere. Wer nach dem Update an einem Closure-Gate Befunde **verliert**, verliert Fehlmessungen. Ferner: die §4-Grund-Code-Zeile und der `--doctor`-Klartext trugen weiter die mit 0.52.1 widerrufene Reichweite („von **allen** Modulen übersprungen") und sind nachgezogen | — |
| 0.52.2 | 2026-08-10 | Nachzug nach bestätigender Re-Review: die **Grenze** von `DC-FA-SPAN-001` Klasse 3 nennt jetzt auch `codepaths` und `pins` — bei `pins` ist die Folge **still** (ein offener Fence verdeckt die Überschrift, der Anker wird unauflösbar, die Drift-Prüfung entfällt kommentarlos statt zu melden) — sowie den eingerückten Code-Block, der in keiner Lesart modelliert ist. Die Fundstellen-Klausel war zu eng an die Paritätskippung gebunden: auch die strenge Lesart zeigt daneben, wenn eine längere Fence-Zeile eine kürzere Öffnung geschlossen hat. Das Befund-Ziel verliert zusätzlich das CR einer CRLF-Zeile | — |
| 0.52.1 | 2026-08-09 | Korrektur nach unabhängigem Code-Review, noch vor dem Release der 0.52.0: `DC-FA-SPAN-001` Klasse 3 sagte eine **Reichweite** zu, die die Umsetzung nicht hielt („jede Vorverarbeitung“, „alle Module“). d-check trägt **zwei** Schluss-Lesarten — den naiven Toggle und den längenabgeglichenen CommonMark-Schluss —, und der Wächter wertete nur die erste aus: ein von `~~~` „geschlossener“ Backtick-Block war für ihn balanciert, für den Tabellen-Leser bis Dateiende offen (belegt: ein Vollständigkeits-Gate meldete Exit 0 über eine ungedeckte Anforderung hinter dem Fence). Der Befund wertet jetzt **beide** aus und meldet, sobald eine von beiden offen endet. Zusätzlich: die Trimmung folgt exakt der Vorverarbeitung (Space/Tab, nicht unicode-weit — eine mit U+00A0 eingerückte Fence-Zeile schaltete sonst NUR den Wächter um und machte ihn blind für den echten offenen Fence dahinter); die gemeldete Zeile heißt ehrlich **Fundstelle**, weil unter reiner Paritätskippung nicht bestimmbar ist, welche Öffnung fehlt; der Sammelsatz nennt das Ziel je Klasse; die **Grenze** deckt neben Post-Pässen jetzt auch **Zieldateien** außerhalb der Scan-Wurzeln, aus denen Module lesen | — |
| 0.52.0 | 2026-08-09 | Change Request (eigener Review-Befund): `DC-FA-SPAN-001` um eine **dritte** Artefakt-Klasse erweitert — `fence-unclosed`, eine Fence-Öffnung ohne Schluss bis zum **Dateiende**. Anlass ist ein **ausgelieferter stiller Grün-Pfad**: hinter einem offenen Fence überspringt jede Vorverarbeitung den Rest des Dokuments, und ein Gate meldet grün, ohne geprüft zu haben (reproduziert gegen das veröffentlichte Image). **Gemeldet wird der Zustand, nicht die Paarung:** der längenabgeglichene CommonMark-Schluss bleibt ausdrücklich offen — er löst den Fall nicht (ein nie geschlossener Fence bleibt unter jeder Paarungsregel offen) und wäre auf dem gemessenen Bestand wirkungslos (776 Dateien in drei Repos: null unbalancierte Fences, null gemischte Fence-Längen, null `~~~`). Reichweite ist die **Datei**, nicht der Absatz (ein Fence *ist* eine Absatzgrenze), der Befund steht an der **Öffnungszeile** (dort liegt die Reparatur), genau einer je Datei; kein neuer Config-Schlüssel — `spans` ist als Ganzes opt-in. Einordnung als **Erweiterung** statt neues Kürzel nach dem Kriterium „dieselbe Frageform, eine Ebene höher". Algorithmus (Schritt 3) in der Spezifikation; Begründung (Modul-Wahl, Reichweite, ausdrückliche Nicht-Lösung der Paarung) in begleitender ADR; Implementierung und Release folgen separat | — |
| 0.51.0 | 2026-08-09 | Change Request (Schwester-Repo `a-check`): neue Anforderung `DC-FA-STRUCT-001` — **20. Regelmodul `structure`** (opt-in, hermetisch, Post-Pass): Struktur-Invarianten **innerhalb** eines Dokuments. Eine **Liste** von Regeln bindet je eine Dokumentklasse (`files`-Glob über Wurzel-relative Pfade, unabhängig vom globalen Scan-Scope) an einen Abschnitts-Typ (`section` Klartext **oder** `section-pattern` RE2) und an Bedingungen in ihm. **`sections`** entscheidet die Kardinalität: `one` (Default — genau einer erwartet, mehrere ⇒ `section-ambiguous` und keine Messung) oder `each` (**jeder** Treffer wird geprüft — für Klassen mit **wiederkehrenden** Abschnitten wie Anforderungen; ohne diesen Modus bliebe jede neu hinzukommende Anforderung ungeprüft). Bedingungen mit je **eigenem** Grund-Code, weil jede eine andere Reparatur verlangt: `non-empty`→`section-empty`, `min-sentences`→`section-thin`, `max-tasks` (Task-Items **im Abschnitt**, nicht dateiweit)→`section-oversized`, `forbid-pattern`→`section-forbidden`, `require-pattern`→`section-pattern-missing`, `require-all`→`section-marker-missing`; ein Sammel-Code schiede aus, weil die Befund-Deduplikation zwei Verletzungen desselben Abschnitts zusammenfallen ließe und die Unterscheidung im nicht stabil zugesagten Meldungstext läge. **Marken sind Auszeichnungs-Marken:** hervorgehobener Textlauf am Zeilen-Anfang nach optionalem **Listen-Marker** — gemessen an den beiden Repos, die den Antrag tragen (108× Listen-Item, 44× bare, dazu qualifizierte Formen). `require-pattern` ist das Spiegelbild von `forbid-pattern` und deckt zugesagte Aussagen **innerhalb** einer Auszeichnung; eine Marken-Alternative wird dadurch überflüssig. Fence-treu, opt-in, diagnose-only, fail-closed am Config-Rand **und** bei Leerlauf (null Kandidaten, auch nach `exempt-paths`). **Mit-Modifikation `DC-FA-PLAN-001`:** die Closure-Fähigkeit wird als **Preset** derselben Semantik im Modus `one` ausgewiesen (nicht superseded — ihre Grund-Codes sind stabil zugesagt) und additiv um `closure-note-ambiguous` ergänzt, das die bis dahin stille Übernahme des ersten Treffers beendet; ein Akzeptanzkriterium beider Anforderungen hält die Semantiken zusammen. §1 benennt die Form-Frage als eigene Kategorie (sie lief mit `spans`/`hostpaths` längst mit, war aber nie ausgesprochen), Bereich `STRUCT` in §3, `structure` in `DC-FA-CLI-002`. Ausdrücklich Out-of-Scope: Aussagen über den **Ort** eines Dokuments, eine **Stichtags**-Mechanik und **namentliche** Ausnahmen innerhalb einer Datei — alle drei hießen, die Kennungs-Konvention des Adopters zu interpretieren, und d-check führt keinen zweiten Regel-Interpreter. Algorithmen und Schema-Schlüssel in der Spezifikation; Begründung (Modul-Schnitt, Nicht-Supersede wegen Code-Stabilität, Kardinalitäts-Modus, gemessene Marken-Form, Modul-Grenze) in begleitender ADR; Implementierung, Paritäts-Beleg und Release folgen separat. Anlass: ein Adopter hat 480 Zeilen Shell in drei Prüfskripten vermessen — 11 Prüfungen, davon 2 gedeckt, 3 nach Kalibrierung, 4 ungedeckt, 2 außerhalb | — |
| 0.50.0 | 2026-08-09 | Change Request (Auftraggeber): der **Closure-Note-Qualitäts-Nachlauf** wird mechanisiert, in zwei Teilen. (1) `DC-FA-PLAN-001` geschärft — das Modul `planning` bekommt eine **zweite Fähigkeit**, opt-in über `planning.closure.dir`: für jeden Slice im Closure-Verzeichnis wird der Closure-Notiz-Abschnitt (erste Überschrift auf `heading-pattern`, Default `^#{2,3} .*Closure-Notiz`) **strukturell** geprüft — Abschnitt vorhanden (`closure-note-missing`), mindestens `min-sentences` Satzende-Zeichen **außerhalb Fenced-Code** (Default 4, `closure-note-thin`), keine der literalen `boilerplate`-Phrasen (`closure-note-boilerplate`, Liste per Default **leer** — der Vertrag bringt keine sprach-spezifischen Phrasen mit). Ohne den Schlüssel inert und byte-identisch; fail-closed bei fehlendem Verzeichnis bzw. ungültigem Muster; diagnose-only. Zugesagt ist **Struktur, nicht Bedeutung** — die Floskel-*Semantik* bleibt ausdrücklich unzugesagt und einem inferentiellen Nachlauf überlassen. (2) Neue Anforderung `DC-FA-CLI-012` — `--config <datei>` überschreibt pro Lauf den konventionellen Konfigurations-Pfad (innerhalb der Scan-Wurzel, gleiche strikte Validierung, **kein** stiller Rückfall, ersetzt statt ergänzt). Sie ist die Voraussetzung dafür, dass die Closure-Prüfung an einem **eigenen** Bindepunkt hängt statt im inneren Loop mitzulaufen: ein Repo fährt zwei disjunkte Prüf-Profile, ohne die Modulwahl auf der Kommandozeile nachzubauen. `DC-FA-CONF-001` notiert den Pfad entsprechend als Konvention statt Zwang. Bereiche `PLAN`/`CLI` bestehen, kein neues Modul. Algorithmen, Grund-Codes und Schema-Schlüssel in der Spezifikation; Begründung (Modul-Schnitt statt neues Modul, Bindepunkt-Trennung, Struktur-vs-Bedeutung-Grenze, Schwellenwahl) in begleitender ADR; Implementierung, Realdatenbeleg und Release folgen separat. Anlass: die Baseline vendored beide Ziel-Formen (Struktur-Gate + inferentieller Reviewer-Skill), d-check hatte **keine** von beiden — und die eigene Closure-Note-Pflicht stand ohne jede maschinelle Entsprechung | — |
| 0.49.0 | 2026-07-19 | Neue Anforderung `DC-FA-SRC-001` — **19. Regelmodul `sources`** (opt-in, **Netz**): Upstream-Content-Drift externer Quellen. Ein auf einen `sha256` gepinnter externer Verweis (per Marker `<!-- source-pin: sha256:… -->` am externen Link **oder** per Config-Block `sources:`) wird geholt, gehasht und verglichen; Abweichung → `source-drift` (mit **vollem Ist-Hash** zum Re-Pinnen), Fetch-Fehler → `source-unreachable`. Zwei Quelltypen: Einzeldatei (Roh-Byte-Hash) und Archiv (`unpack: zip` → pfad-sortierter, byte-genau definierter **Content-Manifest**-Hash, Zip-Reihenfolge-invariant; konzeptionell wie `SHA256SUMS`). Content-Hash-Geschwister von `pins` (in-repo) und `external` (Netz-Erreichbarkeit); produktisiert das Kurs-Beispiel `check_regelwerk_drift.py`. **Erweitert `DC-QA-03`** um eine zweite Netz-Tür (`external` **und** `sources`; nie im netzlosen Default-Lauf, Netzlos-Modullisten-Test um `sources` ergänzt). Bereich `SRC` in §3, `sources` in `DC-FA-CLI-002` + Glossar; `DC-FA-SRC-001.a` + Grund-Codes `source-drift`/`source-unreachable` (§4) in der Spezifikation; Begründung (Netz-Modul-Design, DC-QA-03-Amendment, Manifest-Hash, Marker + Config) in begleitender ADR. Implementierung/Realdatenbeleg/Release folgen separat. Anlass: Nutzer-Frage „Drift gegen Upstream in d-check einbauen" — der Bash-Helfer `fetch-baseline-cache.sh --check-latest` deckt d-checks eigene Baseline, das Modul macht es reusable für jeden Adopter (u. a. `ai-harness-course`, dessen `check_regelwerk_drift.py` genau darauf DEFERRED ist) | — |
| 0.48.0 | 2026-07-18 | Change Request (Adopter `ai-harness-init`): **Zitat-Verifikation** von `datei:zeile`-Zitaten, zwei Teile. (1) `DC-FA-CODE-001`-Erweiterung — opt-in `codepaths.check-lines` prüft die Zeilen-Referenz eines Inline-Code-Pfads (`datei:<von>-<bis>`): Ziel existiert **und** hat ≥ `bis` Zeilen (sonst `citation-out-of-range`) **und** `von ≤ bis` (sonst `citation-inverted-range`); Default aus **byte-identisch** (die Zeile wurde bisher erkannt und verworfen). (2) Neue Anforderung `DC-FA-CITE-001` — 18. Regelmodul `citations` (opt-in): ein per Direktive `d-check:cite` ausgezeichneter Zitatblock wird **zeichengenau** gegen die zitierte Quell-Spanne geprüft (`citation-mismatch`), greift die `codepaths`-`datei:zeile`-Erkennung auf (kein zweiter Detektor), hermetisch, fail-closed. Bereich `CITE` in §3, `citations` in `DC-FA-CLI-002` + Glossar; `DC-FA-CODE-001.a`-Erweiterung + `DC-FA-CITE-001.a` + Grund-Codes in der Spezifikation; Begründung (Erweiterung vs. eigenes Modul, `verbatim`/Direktiven-Design) in begleitender ADR. Zweite Direktive `d-check:cite`. §4-Vorfragen entschieden: Adopter-Rückfrage empirisch (33/33 Zitate in Inline-Code ⇒ `codepaths`-Erweiterung, kein Prosa-Scanning), Zuschnitt Form (c). Anlass: die committet-vendored Baseline erzeugt Zeilen-Zitate, die beim nächsten Tag-Bump still verfaulen | — |
| 0.47.0 | 2026-07-18 | Change Request (Konsument `ai-harness-course`): neue Anforderung `DC-FA-REF-001` — **geteiltes Referenz-Ventil** `ignore-refs` mit **Quell-Skopus** (`in:`) und Zwei-Feld-Semantik (`refs ∧ ¬keep`; `keep` gewinnt unbedingt und reihenfolge-unabhängig, **nicht** gitignore-Last-Match). Das bisher **modul-lokal** in `DC-FA-CODE-001` wohnende `ignore-refs` wird zur **querschnittlichen** Anforderung: `links`/`anchors`/`codepaths` verweisen darauf; `codepaths.ignore-refs` bleibt **Alias** (kein Config-Bruch). Ziel-Achsen-Pendant zu `scan.ignore` (`DC-FA-SCAN-001`, Quell-Achse — jene entfernt Dateien, dieses Ziele). Bereich `REF` in §3; `DC-FA-CODE-001`/`DC-FA-LINK-001`/`DC-FA-ANCH-001` an das Ventil angebunden. §4-Vorfrage vom Auftraggeber entschieden (Option a: neues geteiltes Kürzel statt Änderung dreier Anforderungen; **gemeinsames Kriterium**: querschnittlich → neues Kürzel, Einzelmodul → bestehende Anforderung ändern). `DC-FA-REF-001.a` + Schema-Keys (`ignore-refs[].in`/`refs`/`keep`) + Alias-Semantik in der Spezifikation, Begründung (Zwei-Feld vs. `!`-Negation, Alias-Pfad) in begleitender ADR, Implementierung/Realdatenbeleg/Release folgen separat. Anlass: Template-Verzeichnisse mit Ziel-Repo-Platzhaltern zwingen heute das ganze Verzeichnis in `scan.ignore` und machen damit auch die **echten** Verweise blind, deren Auflösung beim Release unveränderlich eingefroren wird | — |
| 0.46.0 | 2026-07-17 | Change Request (Auftraggeber): die **Komma-Kurzform** `<FAM>-AAA, BBB` wird **fail-closed** (Exit 2) statt still verschluckt — betrifft `DC-FA-COV-001` und, über den geteilten Reader, `DC-FA-XREF-001`. Zugesagt waren nur `..`-Range und `/`-Aufzählung; die Kurzform war nie Vertrag, ihr **stiller Drop** aber auch nicht: `GG-SCN-001, 007` deckte nur `GG-SCN-001` und erzeugte eine falsche Waise. Verworfen: (a) Komma-Enum unterstützen — d-check kann eine gemeinte Kurzform nicht von einer Zahl im Fließtext unterscheiden (`GG-QA-001, 007 Sekunden`), das wäre Raten; (b) Status quo — der stille Drop ist die schlechteste Option. Gewählt: die **Gestalt** triggert (Kennung + Komma + Ziffern), der Inhalt ist keine unterstützte Notation ⇒ Exit 2 mit Hinweis — dieselbe Logik wie bei `AAA>BBB`. Ein Komma vor einer **vollständigen** Kennung (`GG-AR-COMP-CORE, GG-AR-COMP-DOMAIN`) bleibt unberührt. Neues Negative- und Boundary-Kriterium; Schritt 3 in `DC-FA-COV-001.a`; Begründung in der begleitenden ADR. Anlass: Konsumenten-Report grid-gym gegen v0.45.1 — reale §27.1-Zeilen `GG-SCN-001, 007` ließen das produktiv verdrahtete `trace.coverage` das Mapping nicht zählen | — |
| 0.45.0 | 2026-07-17 | Change Request (Auftraggeber): `DC-FA-XREF-001` um **`forward.req-pattern`** erweitert (RE2, Default `requirements.id-pattern`) — symmetrisch zum vorhandenen `backward.req-pattern`, ein Denkmodell statt zwei. Bis dahin las die Vorwärts-Sicht ihre Anforderungs-IDs **still** über `requirements.id-pattern`, die Rück-Sicht über ihr eigenes Muster; die Kopplung stand weder im Vertrag noch in der Config-Oberfläche. Folge bei einer bewusst gescopten RTM (Architektur-Meta ausgeschlossen): die Vorwärts-Sicht ist leer, **jede** Rück-Kante wird als „ohne RTM-Eintrag“ gemeldet, und die eigentliche `F \ B`-Drift verschwindet — ein Falschbefund, der wie echter Drift aussieht. Neu festgehalten: **die Vergleichs-Schlüsselmenge ist nicht die RTM-Anforderungsmenge** (belegt: eine ID, die das Muster trifft, aber keine RTM-Zeile hat, wird verglichen). Neues Boundary-Kriterium (RTM-Scope ≠ Vergleichs-Scope); Schema-Zeile + Algorithmus-Schritt 2 in `DC-FA-XREF-001.a`; Begründung in der begleitenden ADR. Anlass: Realdaten-Lauf des Konsumenten grid-gym gegen v0.44.0 (Defekt 1 von zweien), belegt durch Umschalten **nur** dieses Musters bei identischen Dateien | — |
| 0.44.2 | 2026-07-17 | Change Request (Auftraggeber): die Vakuitäts-Stufe aus 0.44.1 von einer **Ursachen-** auf eine **Wirkungs**-Fassung gezogen. 0.44.1 band Vakuum an die Muster-Ursache („Given ein `design-pattern`, das … vorbeigreift“); ein **übergriffiges `exclude-req`** — fehlerfreie Muster, aber ein Ventil, das jede Anforderung verschluckt — schaltete das Gate bei realem Drift ebenso still ab, war aber von keinem Akzeptanzkriterium gedeckt. Vakuität wird daher **nach** dem Ausschluss gemessen (das Ventil ist selbst eine kuratierte, drift-fähige Kante), und geprüft wird die **Wirkung** — ob der Abgleich konstruktionsbedingt überhaupt einen Befund liefern könnte — statt einer Ursachenliste, die bei der nächsten unbekannten Ursache erneut risse. Neues Negative-Kriterium (Ventil, samt ventil-benennender Meldung); Ausschluss-Stufe in der Fehlerpräzedenz von `DC-FA-XREF-001.a`; Begründung in der begleitenden ADR. Anlass: unabhängiges Closure-Review — R3 reproduzierte `exclude-req: '.'` mit `0 Differenz(en)`/Exit 0 bei echtem Drift, R4 wies die AK-Lücke nach | — |
| 0.44.1 | 2026-07-17 | Change Request (Auftraggeber): `DC-FA-XREF-001` um eine **Vakuitäts-Stufe** geschärft — ein Abgleich, aus dem keine Kante fällt, ist **kein bestandener** Abgleich, sondern Exit 2. Vakuum ist definiert als **beide** Sichten kantenleer (typischer Anlass: ein `design-pattern`, das kompiliert, aber am Artefakt-Namensraum vorbeigreift — die Namensraum-Kongruenz war bis dahin nur als Vorbedingung *beschrieben*, nicht mechanisiert) **oder** die Rück-Sicht kantenleer unter `mode: superset` (dann kann `B \ F` konstruktionsbedingt nie einen Befund liefern). Abgegrenzt: eine **einseitig** leere Vorwärts-Sicht bleibt ein wohldefiniertes Ergebnis (Diff über `keys(F) ∪ keys(B)`) und meldet `B \ F` laut — der erwartete Bootstrap-Zustand vor der Restrukturierung der Vorwärts-Sicht, den ein Guard sonst mit einer Config-Fehldiagnose abwürgte. Neues Negative- und Boundary-Kriterium; Vakuitäts-Stufe in der Fehlerpräzedenz von `DC-FA-XREF-001.a` der Spezifikation; Begründung + Abgrenzung zum Nullmengen-Guard der Tabellenquellen in der begleitenden ADR. Anlass: unabhängiges Closure-Review — R1 reproduzierte das stille Grün, R2 wies den ersten, symmetrisch je Sicht feuernden Fix als vertragswidrig nach | — |
| 0.44.0 | 2026-07-16 | Change Request (Auftraggeber): neue Anforderung `DC-FA-XREF-001` — **Kreuzverweis-Konsistenz zweier Traceability-Sichten** (`trace.cross-consistency`, opt-in): der `--trace`-Lauf vergleicht zusätzlich die Vorwärts-RTM-Tabelle (Anforderung → Design-Menge) gegen die Rückwärts-`Bezug`-Kanten (Design → Anforderung) und meldet je Anforderung die beiden Mengendifferenzen `F(R) \ B(R)`/`B(R) \ F(R)` mit Richtungslabel; Modus `equal`/`superset`, `1:N` Normalfall. Beide Sichten kuratierte Tabellen über den header-gebundenen Reader (`DC-FA-REQ-001`) + range-aware Span-Semantik (`DC-FA-COV-001`); Rück-Kanten = Quelle der Wahrheit (Artefakt-ID = erste Spalte via `design-pattern`, `Bezug`-Zelle range-aware). Ableitungssprünge per `exclude-req` (RE2) ausgenommen. Advisory unter `--trace`; Exit-Änderung nur über das globale `--require-complete` (`DC-FA-CLI-011`). Fail-closed (ungültiges Regex/fehlende Spalte/ID-Header nicht genau einmal/unbekannter `mode` ⇒ Exit 2); ohne `trace.cross-consistency`-Block RTM byte-identisch (`DC-QA-02`/`DC-QA-03`). Bereich `XREF` in §3; additive Erweiterung des `--trace`-Laufs (`DC-FA-CLI-009`) **ohne** RTM-Änderung (separate advisory Befunde, keine RTM-Spalte); `DC-FA-XREF-001.a` + Schema-Keys (`trace.cross-consistency.*`) in der Spezifikation; Implementierung, Realdatenbeleg und Release folgen separat. Anlass: Konsument grid-gym (Trigger 088) — §27.1 nannte {GG-AR-COMP-CORE, GG-AR-COMP-DOMAIN}, die `Bezug`-Rück-Kanten {GG-AR-P-005, GG-AR-P-009, GG-AR-COMP-SCHED}, Schnittmenge null, von keinem Gate bemerkt | — |
| 0.43.0 | 2026-07-14 | Change Request (Auftraggeber): neue Anforderung `DC-FA-REQ-001` — **native tabellarische Anforderungsquellen + Nullmengen-Guard** für die RTM. `trace.requirements.format: table` liest Markdown-Pipe-Tabellen über konfigurierte Header-Namen (`table.id-column`, genau eine von `table.text-column`/`.text-columns`, optional `.modality-column`); die ID-Zelle muss als Ganzes auf `id-pattern` passen, die gebundene Textzelle liefert RTM-Titel/Body, eine gesetzte Modalitätsspalte ist die alleinige Keyword-Quelle für `DC-FA-MOD-001`. `table.duplicate-ids` ist `error` (Default), `first` oder `last`. Beide Formate münden in dasselbe RTM-Modell. **Fail-closed:** nichtleer explizite `requirements.source` oder Tabellenmodus mit null erkannten Anforderungen sowie unbekanntes Format/fehlende oder doppelte konfigurierte Spalten ⇒ Exit 2 (auch mit `--require-complete`) statt irreführendem `0 Anforderungen, 0 Waisen`/Exit 0. `source: ""` gilt wie abwesend. Ohne `trace`-Block bleibt der Heading-Default samt Deduplizierung/leerer RTM/Exit 0 byte-identisch (`DC-QA-02`/`DC-QA-03`). Mit-Modifikation `DC-FA-CLI-009`; Bereich `REQ` in §3; Spezifikation, ADR, Implementierung, Realdatenbeleg und Release folgen separat. Anlass: reproduzierter Konsumentenbefund `m-trace` mit d-check v0.42.0 — 371 eindeutige IDs in Tabellen mit den Text-Headern `Anforderung`/`Akzeptanzkriterium`, explizite Quelle und passende Regex ergaben null Anforderungen bei Exit 0 | — |
| 0.42.0 | 2026-07-11 | Change Request (Auftraggeber): neue Anforderung `DC-FA-MOD-001` — **Modalitäts-Klassifikation der Anforderungen** (`trace.requirements.modality`, opt-in): d-check klassifiziert jede RTM-Anforderung anhand **konfigurierbarer Modal-Verb-Schlüsselwörter** (Built-in DE+EN-RFC-2119-Defaults; `modality.levels` Stufe→Keywords, `modality.require-levels` welche Stufen gaten, Default `[must]`) über den **ersten** Treffer im Anforderungs-Body (Längster-Treffer-zuerst gegen `MUSS NICHT`≠`DARF NICHT`, case-insensitiv/wortgrenzen-genau; kein Treffer ⇒ Stufe `unknown`). Neue **Modality-Spalte** (konditional); `--require-complete` ([`DC-FA-CLI-011`](#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)) bricht dann **nur** bei Waisen der `require-levels`-Stufen (SOLLTE/KANN/`unknown` advisory). Fail-closed (leerer Stufen-Name/Keyword, `require-levels`-Eintrag weder Stufe noch `unknown` ⇒ Exit 2), strikt opt-in, default-aus **byte-identisch** (`DC-QA-02`/`DC-QA-03`; ohne `modality` gaten alle Waisen wie bisher). Bereich `MOD` in §3; `DC-FA-MOD-001.a` + Schema-Keys (`trace.requirements.modality.levels`/`require-levels`) in der Spezifikation; Mit-Modifikation `DC-FA-CLI-009` (Modality-Spalte). `--print-config` führt den `modality`-Block. Begründung + Matching-/Unknown-Semantik in begleitender ADR. Anlass: Konsumenten-Analyse grid-gym — die 10 Coverage-Rest-„Waisen" sind 5× KANN (Future) + 4× Nicht-Ziele + 1× DARF NICHT; die slice-zentrische RTM behandelte MUSS und KANN gleich | — |
| 0.41.0 | 2026-07-11 | Change Request (Auftraggeber): neue Anforderung `DC-FA-COV-001` — **kuratierte Coverage-Quellen der RTM** (`trace.coverage`, opt-in): eine dritte Referenzklasse (Liste benannter `files` + `label` + `ranges` + `sections`/`exclude-sections`), die eine ausgelagerte Traceability-Matrix **range-aware** als Coverage einliest, **ohne** `adrs`/`slices` zu berühren. `<FAM>-AAA..BBB`/`<FAM>-AAA/BBB/CCC` expandieren (breiten-erhaltend, gegen `requirements.id-pattern` validiert); Abschnitts-Whitelist/Blacklist (dieselbe Span-Semantik wie `matrix.exclude-sections` — gegen die „…ohne Design-Artefakt"-Falle). Mit-Modifikationen: `DC-FA-CLI-009` (RTM trägt bei aktiver `trace.coverage` eine **Coverage-Spalte** + `coverage`-Feld in `--json`/`--yaml`) und `DC-FA-CLI-011` (**Waise = ¬slice ∧ ¬coverage**; Slice **oder** Coverage deckt ab, bloße ADR-Referenz weiterhin nicht). Fail-closed (fehlende Datei / ungültige Range `AAA>BBB` / Breiten-Mismatch ⇒ Exit 2), strikt opt-in, default-aus **byte-identisch** (`DC-QA-02`/`DC-QA-03`). Bereich `COV` in §3; `DC-FA-COV-001.a` + Schema-Keys (`trace.coverage[].files`/`label`/`ranges`/`sections`/`exclude-sections`) in der Spezifikation; `--print-config` führt den `coverage`-Block. Begründung + Range-/Sektions-Semantik in begleitender ADR. Anlass: Konsumenten-Analyse grid-gym — 171 „Waisen" waren zu ≥122 anderswo (ADR/traceability.md/Wellen) belegt; d-checks slice-zentrische RTM verfehlte die kuratierte Deckungs-Matrix | — |
| 0.40.0 | 2026-07-11 | Change Request (Auftraggeber): `DC-FA-CLI-009` (Requirements Traceability Matrix) um einen **opt-in `trace`-Config-Block** erweitert — die vier bislang hart an d-checks Konvention gebundenen RTM-Annahmen sind überschreibbar: Anforderungs-Quelldatei + Kennungs-Gestalt (`trace.requirements.source`/`.id-pattern`) sowie je Referenzklasse ADR/Slice das Verzeichnis, die Dateinamen-Gestalt (Capture-Gruppe = Owner-Kennung) und das Owner-Präfix (`trace.adrs.*`/`trace.slices.*`). Jedes Feld optional; abwesend ⇒ d-checks Konventions-Default ⇒ RTM **byte-identisch** (`DC-QA-02`), nichts geschrieben (`DC-QA-03`). Fail-closed: ungültige Regex oder `file-pattern` ohne Capture-Gruppe ⇒ Exit 2. **Kehrt die 0.21.0-Out-of-Scope-Zeile „frei konfigurierbare Quell-Pfade jenseits der adoptierten Harness-Konvention" um** (durch bewusst begrenzte, explizit gesetzte Config ersetzt — kein Ableiten, je eine ADR-/Slice-Klasse, kein leerer ADR-Präfix, kein VCS). Config-Schema `trace.*` + Algorithmus-Konfigurierbarkeit in `DC-FA-CLI-009.a` der Spezifikation; Begründung (Konsumenten-Konventions-Bindung, Design spiegelt `ids.patterns`) in begleitender ADR. `DC-FA-CLI-011` (`--require-complete`) erbt die konfigurierten Quellen unverändert (gleicher RTM-Lauf). Anlass: Konsumenten-Befund grid-gym 2026-07-11 — `make doc-trace` sah nur 6 von 243 Anforderungen (allein die `GG-QA-*`-Familie traf zufällig d-checks `-QA-`-Default-Gestalt), die übrigen 40 Familien und alle `NNN-…md`-Slices blieben unsichtbar | — |
| 0.39.0 | 2026-07-06 | Schärfung `DC-FA-CLI-006` (Auftraggeber): die `--suggest-config ai-harness[-init]`-Vorlage an die **gelebte** Dogfood-Konvention angeglichen — `spans`/`hostpaths` ins fixe Standard-Modulset aufgenommen (revidiert die 0.26.0-„situativ nicht aktiviert"-Einordnung; d-checks eigene `.d-check.yml` führte sie längst) und ein **repo-bewusster `planning`-Block** ergänzt (aktiv bei vorhandener `roadmap`, sonst auskommentiert; im Voll-Kanon `ai-harness-init` aktiv). `vcs`/`commits` (Commit-Range) werden auf `--print-mk` verwiesen statt ins statische `modules` gelegt; `versions`/`targets` bleiben bewusst vertagt (repo-spezifische `pin-pattern`/`authority`). Neu benannt: **Eignungs-Kriterium K1–K4** (konventions-kanonisch · ableitungsfrei/konventions-feste Config · Baum-Scan-tauglich · hermetisch) + die **geschlossene** Aktiv-Menge im Body; die kanonische Vorlage in der Spezifikation (`DC-FA-CLI-006.a`) deckt jetzt die emittierte Ausgabe **1:1** (Kommentarzeile inkl. `--print-mk`-Verweis + `codepaths`-Block ergänzt) — der Normativitäts-Spalt (Code-Kommentar ≠ Norm) ist geschlossen. Zugleich die situativen-Modul-Enumeration (AKs/Out-of-Scope) um das seit 0.38.0 fehlende `targets` vervollständigt. Betrifft nur die `ai-harness`-Vorlage (nicht den generischen Quellen-Modus, nicht die eigene `.d-check.yml`); `DC-QA-02`/`DC-QA-03` unberührt (nur mehr Ausgabe). Begründung + Eignungs-Kriterium in begleitender ADR. Anlass: Nutzer-Analyse „welche Module nutzt `--suggest-config ai-harness` nicht, und warum nicht" (2026-07-06) | — |
| 0.38.0 | 2026-07-05 | Neue Anforderung `DC-FA-TGT-001` (Modul `targets`, opt-in): Deklarations-Konsistenz zwischen Doku und Build-Targets — jedes in einer Doku-**Tabellenzeile** als ` `make X` ` behauptete Target muss eine Makefile-Regel sein (sonst `gate-phantom`), und jede Makefile-Regel (minus `targets.exempt-targets`) muss in der Autoritäts-Doku als ` `make X` ` stehen (sonst `gate-undocumented`) — die maschinelle Form des Vertrags „kein halluziniertes / kein undokumentiertes Gate". **Tabellen-Scoping** (`^\|`, nur Tabellenzeilen — Prosa zählt nicht, sonst spuriöse `gate-phantom`). **Hermetisch** wie `DC-FA-PLAN-001`: Filesystem-Port `ReadFile` auf die konfigurierten Dateien, **kein** Makefile-Ausführen, **kein** git, **kein** Netz (`DC-QA-02`/`DC-QA-03` unberührt); Zeilen-Heuristik für Regelnamen (keine Pattern-Rules/variablen Targets, ≡ dem abgelösten Skript). Strikt opt-in, fail-closed (fehlende konfigurierte Datei ⇒ Exit 2; leeres `makefiles`/`doc-tables` ⇒ inert), diagnose-only, default-aus byte-identisch. 17. Regelmodul; Bereichskürzel `TGT` in §3, `targets` in `DC-FA-CLI-002` + Glossar, Algorithmus-Sektion `DC-FA-TGT-001.a` + Schema-Keys (`targets.makefiles`/`doc-tables`/`authority`/`exempt-targets`) in der Spezifikation; die Grund-Codes `gate-phantom`/`gate-undocumented` (§4) landen mit der Modul-Implementierung (AllReasons-↔-§4-Lockstep). Zugleich **`DC-FA-CLI-010`-Erweiterung** (10→11 Targets): `--print-mk` trägt ein `doc-targets`-Target (`--enable targets` + Fokus-`--disable`, hermetisch ohne Range); `--print-config`/`--suggest-config` führen `targets`. Anlass: Auftraggeber 2026-07-05 — `tools/gate-consistency.sh` driftet über die Repo-Familie (a-check ≠ d-check), Klasse-A-Mechanisierung [`MR-007`](../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding) „verteilen statt kopieren": die verteilbare Modul-Form löst den cross-repo-Kern (Doku↔Makefile), d-check dogfoodet das Modul für sein eigenes `make gate-consistency`-Gate. Identitäts-Ausweitung „Doku-Checker → Deklarations-Konsistenz-Checker" (Makefile-Lesen via Filesystem-Port, wie `planning` die Verzeichnis-Struktur) in begleitender ADR | — |
| 0.37.1 | 2026-07-03 | Review R1 (doc) + R2 (code) zum `DC-FA-TRK-001`-CR, präzisiert: Verzeichnis-Ziele (Index führt nur Dateien) und Symlink-Referenzen (kategorisch [`DC-FA-LINK-002`](#dc-fa-link-002--symlink-ablehnung)-Domäne — sonst false-positive hinter getrackten Verzeichnis-Symlinks) explizit kein Kandidat/Out-of-Scope; Auflösungs-Mechanik ausdrücklich unabhängig von der Aktivierung des Moduls `links`; intent-to-add + linked worktree als Ränder benannt. `DC-FA-CLI-010`-AKs/Out-of-Scope von neun auf **zehn** Targets nachgezogen (`doc-tracked` — Selbstwiderspruch behoben); `DC-FA-CLI-006`-Enumeration der situativen Module um `commits`/`planning`/`tracked` vervollständigt | — |
| 0.37.0 | 2026-07-03 | Neue Anforderung `DC-FA-TRK-001` (Modul `tracked`, opt-in): Getrackt-Status auflösbarer Referenz-Ziele — jedes auflösbare, **existierende** Link-/Bild-Ziel muss im **git-Index getrackt** sein, sonst `target-untracked` (die Referenz wäre auf jedem frischen Klon `target-missing` — Umgebungs-Drift zwischen Arbeitsbäumen, am Entstehungsort gefangen statt erst in der CI des nächsten Checkouts). Wahrheit ist der **Index**, nicht die `.gitignore`-Syntax (kein zweiter Regel-Interpreter; frisch gestagte Dateien gelten als getrackt — WIP-tauglich). Liest `.git` über den **VCS-Port** (dritte Nutzung: `vcs` Range-Diff, `commits` Messages, `tracked` **Index** — ohne Range), reine-Go/ohne Netz, lokal/lesend/deterministisch (`DC-QA-02`/`DC-QA-03` wie bei `DC-FA-VCS-001` formuliert). Kein Doppelbefund (nur existierende Ziele; `target-missing` bleibt `links`, Prinzip von `DC-FA-PIN-001`), Ventil `tracked.exempt-targets`, strikt opt-in/fail-closed (aktiv ohne lesbares `.git` ⇒ Exit 2)/diagnose-only/default-aus byte-identisch. 16. Regelmodul; Bereichskürzel `TRK` in §3, `tracked` in `DC-FA-CLI-002` + Glossar, Algorithmus-Sektion `DC-FA-TRK-001.a` + Grund-Code `target-untracked` + Schema-Key (`tracked.exempt-targets`) in der Spezifikation. Zugleich **`DC-FA-CLI-010`-Erweiterung** (9→10 Targets): `--print-mk` trägt ein `doc-tracked`-Target. Anlass: Auftraggeber-Frage 2026-07-03 („Was passiert, wenn ein Dokument ein gitignoriertes Dokument referenziert?") + Demo-Beleg (Erzeuger-Checkout grün, frischer Klon `target-missing`); heutige Ventile (CI-Netz, gitignore+`scan.ignore`-Doppel nach [`MR-017`](../harness/conventions.md#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen)-Muster, Vendoring nach [`MR-019`](../harness/conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)-Muster) fangen die Falle erst spät oder nur je Einzelfall | — |
| 0.36.0 | 2026-07-01 | Neue Anforderung `DC-FA-PLAN-001` (Modul `planning`, opt-in): maschinelle Durchsetzung der Planning-Lifecycle-Invariante — die Roadmap trägt den Ruhe-Marker (`planning.marker`, „Keine aktive Welle") in ihrem `planning.heading`-Abschnitt (`## Aktuelle Welle`) genau dann, wenn kein `slice-*` (`planning.slice-glob`) im Verzeichnis liegt (`hasActive == hasSlices`), sonst `planning-drift`. **Hermetisch** (nur Roadmap-Datei + Verzeichnis-Listing, **kein** git, **kein** Netz, read-only) — normales Modul wie `codepaths`, `DC-QA-02`/`DC-QA-03` unberührt; fail-closed bei fehlender kanonischer Überschrift bzw. Roadmap-Datei (Heading-Guard), strikt opt-in, diagnose-only, default-aus byte-identisch. 15. Regelmodul; Bereichskürzel `PLAN` in §3, `planning` als Modul in `DC-FA-CLI-002` + Glossar, Algorithmus-Sektion `DC-FA-PLAN-001.a` + Grund-Code `planning-drift` + Schema-Keys (`planning.roadmap`/`marker`/`heading`/`slice-glob`) in der Spezifikation. Anlass: Auftraggeber — `tools/planning-consistency.sh` mechanisieren (letzter Kandidat des `tools/*.sh`-Audits, [`MR-007`](../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Klasse „verteilen statt kopieren"); bewusst nachrangig (die „Keine aktive Welle"-Konvention ist harness-layout-spezifisch, kleinerer Verteilungswert) — im Gegensatz zu `vcs`/`commits` **ohne** git/VCS-Port, rein hermetisch. `make planning-check` dogfoodet das Modul. Zugleich **`DC-FA-CLI-010`-Erweiterung** (8→9 Targets): `--print-mk` trägt ein `doc-planning`-Target; `--print-config`/`--suggest-config` führen `planning` | — |
| 0.35.0 | 2026-07-01 | Neue Anforderung `DC-FA-COMMITS-001` (Modul `commits`, opt-in): maschinelle Durchsetzung der Traceability-Regel — jede geprüfte Commit-Message muss eine Kennung nach `commits.id-patterns` (`ADR-`/`MR-`/`DC-`/`slice-`) auf einer Inhalts-Zeile tragen, sonst `commit-untraceable`. Liest die Commit-**Messages** über **denselben VCS-Port** wie `DC-FA-VCS-001` (reine-Go-git, **ohne git-Binary** → distroless bleibt, **ohne Netz**), erweitert um Message-Lesen; zwei Quellen für dieselbe Prüfung: `--range <base>..<head>` (CI/Push, Nicht-Merge-Commits) und `--commit-msg <datei\|->` (commit-msg-Hook, einzelne Pending-Message). Uniforme `#`-/scissors-Bereinigung (Kennung auf Inhalts-Zeile), Betreff-Ausnahme `commits.exempt-pattern` (Selbstkonfig `^(Merge \|Revert )`). Strikt opt-in (nie Default, wie `external`/`vcs`), fail-closed ohne `.git`/Range/Message, diagnose-only; erweiterter Eingabe-Scope, aber lokal/lesend/deterministisch — `DC-QA-02`/`DC-QA-03` unberührt. 14. Regelmodul; Bereichskürzel `COMMITS` in §3, `commits` als Modul in `DC-FA-CLI-002` + Glossar, Algorithmus-Sektion `DC-FA-COMMITS-001.a` + Grund-Code `commit-untraceable` + Schema-Keys (`commits.id-patterns`/`exempt-pattern`) in der Spezifikation. Anlass: Auftraggeber — `tools/trace-check.sh` vollständig mechanisieren (Copy-Drift über die Repo-Familie, [`MR-007`](../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Klasse „verteilen statt kopieren"): die verteilbare Modul-Form löst es im Image; d-check dogfoodet das Modul für sein eigenes `make trace-check`-Gate. Das „nächste `adr-check`" — dieselbe VCS-Port-Präzedenz wie `DC-FA-VCS-001`, hier auf Commit-Messages statt Datei-Inhalt. Zugleich **`DC-FA-CLI-010`-Erweiterung** (7→8 Targets): `--print-mk` trägt ein `doc-commits`-Target (`--enable commits` + aus `ValidModules` abgeleitete Fokus-`--disable`-Liste, `--range`) — die **verteilte** Commit-Traceability für Konsumenten ohne Skript-Kopie („verteilen statt kopieren", parallel zu `doc-immutable`); `--print-config`/`--suggest-config` führen `commits` in der Verfügbar-/opt-in-Liste | — |
| 0.34.0 | 2026-06-29 | `DC-FA-CODE-001` (Modul `codepaths`) um `ignore-refs` erweitert: eine Glob-Liste nimmt **bestimmte Ziel-Pfade** von der Existenz-Prüfung aus — **referenz-weit** (egal Datei/Zeile), als Register bewusst entfernter/historischer Artefakte (Tombstones); dritte Ventil-Achse neben dem zeilenweisen `d-check:ignore` und dem datei-weiten `exempt-paths`. Bewusster Akt mit Gate (ohne Eintrag bleibt `codepath-missing`), Default-leer byte-identisch. Schema-Key `codepaths.ignore-refs` + §DC-FA-CODE-001.a-Schritt in der Spezifikation. Anlass: die Frozen-Doc-Refactoring-Falle — immutable ADRs zitieren Code-Pfade, die bei Refactoring/Löschung dangeln, aber nicht editierbar sind; gelöst und zugleich das nur „pfad-stabil behaltene" `tools/adr-immutable-check.sh` entfernt; begleitende ADR (Supersedes die „Skript-behalten"-Teilentscheidung der VCS-ADR) | — |
| 0.33.0 | 2026-06-29 | Neue Anforderung `DC-FA-VCS-001` (Modul `vcs`, opt-in): git-historienbasierte Immutabilität des Core über eine Commit-Range — `core(BASE)` ≟ `core(HEAD)`, geliefert über `--range <base>..<head>` (CI/Push) oder `--staged` (pre-commit). Liest das read-only `.git` über einen eigenen **VCS-Port** (reine-Go-git, **ohne git-Binary** → distroless bleibt, **ohne Netz**); erweiterter Eingabe-Scope (git + Range statt nur Markdown-Baum), aber lokal/lesend/deterministisch — `DC-QA-02`/`DC-QA-03` unberührt. Geprüft wird jede in der Range geänderte, der Klasse (`vcs.paths`) entsprechende Datei mit immutabler BASE (`vcs.immutable-when`): Körper-Drift, unzulässiger Status-Übergang (`vcs.head-allow`) oder Löschung/Umbenennung → `core-drift-vcs`. Strikt opt-in (nie Default, wie `external`), fail-closed bei fehlendem `.git`/Range, diagnose-only. **Die volle git-Garantie, die `DC-FA-IMM-001` (hermetischer Pin) bewusst der späteren VCS-Stufe überließ** — beide koexistieren als Defense-in-Depth. 13. Regelmodul; Bereichskürzel `VCS` in §3, `vcs` als Modul in `DC-FA-CLI-002` + Glossar, Algorithmus-Sektion `DC-FA-VCS-001.a` + Grund-Code `core-drift-vcs` + Schema-Keys (`vcs.paths`/`immutable-when`/`exclude-sections`/`status-line`/`head-allow`) in der Spezifikation. Anlass: Auftraggeber — `adr-immutable-check.sh` vollständig mechanisieren (Copy-Drift über die Repo-Familie, [`MR-007`](../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Klasse): die verteilbare git-Form löst es im Image; d-check dogfooded das Modul für die eigenen Accepted-ADRs. Zugleich **DC-FA-CLI-010-Erweiterung** (6→7 Targets): `--print-mk` trägt ein `doc-immutable`-Target (`--enable vcs` + aus dem Modulsatz abgeleitete Fokus-`--disable`-Liste, RANGE/STAGED) — die **verteilte** git-Garantie für Konsumenten ohne Skript-Kopie („verteilen statt kopieren"); `model.ValidModules` dafür exportiert. Außerdem **Config-Surface-Bereinigung**: `--print-config` (`DC-FA-CLI-005`) führt jetzt **alle** Module (die zuvor fehlenden `pins`/`immutable` + `vcs` ergänzt) und `--suggest-config ai-harness` (`DC-FA-CLI-006`) nennt die situativen opt-in-Module vollständig (`versions`/`pins`/`immutable`/`vcs` nachgezogen) — die „Verfügbar"-Liste war eine Harness-Ehrlichkeits-Lücke über drei Slices hinweg | — |
| 0.32.0 | 2026-06-28 | Neue Anforderung `DC-FA-IMM-001` (Modul `immutable`, opt-in): Immutabilitäts-Pin gegen Core-Drift — eine Datei mit Inline-Marker `<!-- immutable: sha256:… -->` wird gegen den whitespace-normalisierten **Core** (Datei ohne die Marker-Zeile und ohne die per `immutable.exclude-sections` benannten Abschnitte) gehasht; Abweichung → `core-drift`. Hermetische, im read-only-Arbeitsbaum entscheidbare Immutabilitäts-Prüfung (kein git; die git-historienbasierte Diff-Form bleibt Out-of-Scope bzw. spätere opt-in-Stufe, in begleitender ADR festgehalten). Opt-in pro Datei, default-off byte-identisch (`DC-QA-02`), diagnose-only. 12. Regelmodul; Bereichskürzel `IMM` in §3, `immutable` als Modul in `DC-FA-CLI-002` + Glossar, Algorithmus-Sektion `DC-FA-IMM-001.a` + Grund-Code `core-drift` + Schema-Key (`immutable.exclude-sections`) in der Spezifikation. Anlass: Auftraggeber — das ADR-Immutable-Gate war nur ein Skript (Copy-Drift über die Repo-Familie, [`MR-007`](../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Klasse „verteilen statt kopieren"); die verteilbare Pin-Form löst das hermetisch, die volle git-Garantie bleibt einem späteren VCS-Adapter vorbehalten | — |
| 0.31.0 | 2026-06-28 | Neue Anforderung `DC-FA-MTX-003` (Modul `matrix`): **Token-basierte** Referenz-Richtung + Provenance-Marker + Grandfathering. Eine Klasse kann ein `token`-Regex tragen → `matrix` fängt verbotene Referenzen auch als bare ID-Token im Körper (nicht nur als Link), `matrix-forbidden` in Token-Form. Provenance-Marker `<!-- d-check:status-provenance -->` auf der Zeile nimmt eine verbotene Token-Referenz aus (deklarierte Provenance/Verifikations-Zeiger) — `matrix`' **erster** Zeilen-Marker, kehrt die „nur strukturelle Ausnahmen"-Haltung von `DC-FA-MTX-001` bewusst um (benannt/semantisch, nicht generisches Muting); Ehrlichkeit bleibt Reviewer. Neues `matrix.exempt-paths` grandfathered immutable `Accepted`-ADRs (Regelwerk §Referenz-Richtung). Default-aus byte-identisch (`DC-QA-02`). §DC-FA-MTX-001.a-Schritt + Schema-Keys (`matrix.classes[].token`, `matrix.exempt-paths`) in der Spezifikation. Anlass: Auftraggeber — d-check mechanisiert die Referenz-Richtung, die das adoptierte Regelwerk bewusst dem Reviewer überließ (wie schon [`MR-006`](../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)/`matrix`); der Marker macht die Provenance-vs-Entscheidungsgrundlage-Unterscheidung grep-bar | — |
| 0.30.0 | 2026-06-28 | Neue Anforderung `DC-FA-MTX-002` (Modul `matrix`): Verweisrichtung **innerhalb** einer geordneten Dokumentklasse — eine Klasse kann `order` (Pfad-Globs, autoritativste Schicht zuerst; Rang = erster Treffer) plus `direction: no-downward` tragen; ein klasseninterner Abwärtsverweis (höherrangig → niederrangig, auch über mehrere Stufen) ⇒ neuer Grund-Code `matrix-downward`. Additiv zu den Klassen-Paar-Regeln (`DC-FA-MTX-001`); generalisiert die Spec-Straten-Schichtung auf viele Dateien je Schicht (Globs statt Einzeldateien-Listing). Fail-closed: `order`/`direction` nur zusammen, unbekannter `direction`-Wert ⇒ Config-Fehler; rangfreie Mitglieder nehmen nicht teil; Default beide leer ⇒ byte-identisch (`DC-QA-02`). Algorithmus-Schritt in `DC-FA-MTX-001.a`, Grund-Code `matrix-downward` + Schema-Keys (`matrix.classes[].order`/`.direction`) in der Spezifikation. Anlass: Auftraggeber — die alternative Einzelklassen-Aufzählung war als Richtung nicht erkennbar und verschattete (First-Match) tote Regeln; zugleich Konsumenten-Bedarf d-migrate (23 Spec-Dateien ⇒ Glob-Schichten statt 23-Zeilen-Listing) | — |
| 0.29.0 | 2026-06-24 | Neue Anforderung `DC-FA-PIN-001` (Modul `pins`, opt-in): Content-Pin gegen inhaltlichen Drift — ein Link mit Inline-Marker `<!-- dpin: sha256:… -->` (bindet an den unmittelbar vorausgehenden Link derselben Zeile, sonst inert) wird gegen den whitespace-normalisierten **rohen** Ziel-Span (ganze Datei oder Heading-Section, inkl. Fenced-Code) gehasst; Mismatch → `link-stale`. Nur auflösbare Links (struktureller Befund bleibt `DC-FA-LINK-001`/`DC-FA-ANCH-001`, kein Doppelbefund, auch pins-only); opt-in pro Link, default-off byte-identisch (`DC-QA-02`), diagnose-only (`--bless` spätere CR). Bereichskürzel `PIN` in §3, `pins` als Modul in `DC-FA-CLI-002` + Glossar, Algorithmus-Sektion `DC-FA-PIN-001.a` + Grund-Code `link-stale` in der Spezifikation. Anlass: Auftraggeber-Idee 2 (stale citation) + Spike (Drift real, Rauschen ~0 bei Normalisierung) | — |
| 0.28.0 | 2026-06-24 | Neue Anforderung `DC-FA-VER-001` (Modul `versions`, opt-in): Versions-Pin-Konsistenz — alle Pins (`versions.pin-pattern`) müssen die aktuelle Version aus `versions.current-from` (Default `version.md#aktuell`) tragen, sonst `version-stale`; liest dafür Pins **auch in Fenced-Code** (gescopte Fence-Ausnahme), Ventile `exempt-paths`/`d-check:ignore` für historische Pins; opt-in/default-off (ohne Block byte-identisch, `DC-QA-02`), diagnose-only (Auto-Bump-`--repair` als Folge-CR an `DC-FA-CLI-008`). Bereichskürzel `VER` in §3, `versions` als gültiges Modul in `DC-FA-CLI-002` + Glossar, Algorithmus-Sektion `DC-FA-VER-001.a`, Grund-Code `version-stale` und Config-Schema (`versions.pin-pattern`/`current-from`/`exempt-paths`) in der Spezifikation ergänzt. Anlass: Auftraggeber-Idee „nicht vergessen, die Version zu bumpen" + Spike (Meta-Gate-Skript wegen Copy-Drift über die Repo-Familie verworfen) | — |
| 0.27.0 | 2026-06-23 | Change Request (Auftraggeber): `DC-FA-CLI-010` (`--print-mk`-Fragment) um drei Targets + eine Variable erweitert — `doc-doctor` (`--doctor`), `doc-repair` (`--repair`, Recipe-Echo unterdrückt für `git apply`-reine stdout) und `doc-help` (namespaced, listet die `doc-*`-Targets via `##`-Annotationen; **kein** `help` wegen Konsumenten-Kollision) sowie `DCHECK_DIGEST` (Digest-Override per `ifeq`, sticht den Tag). Alle Targets `##`-annotiert (greift das `help` des Konsumenten auf). Read-only/deterministisch unverändert. Anlass: Auftraggeber-Wunsch nach `doc-doctor`/`doc-repair`/Self-Doc/Digest-Komfort | — |
| 0.26.0 | 2026-06-23 | Schärfung `DC-FA-CLI-006` (Auftraggeber): das opt-in-Modul `diagrams` zur Out-of-Scope-Liste der **nicht** auto-aktivierten situativen Module ergänzt (`external`/`spans`/`hostpaths`/`diagrams`); die `--suggest-config ai-harness[-init]`-Ausgabe nennt diese situativen Module stattdessen in einem **Kommentar mit Verweis auf `--print-config`** (Auffindbarkeit ohne Aktivieren eines inerten Moduls — `diagrams` braucht repo-spezifische `patterns`/`defined-in`, lässt sich nicht ableiten). Read-only/advisory unverändert. Anlass: Nutzer-Frage (wird `diagrams` in `--suggest-config ai-harness` berücksichtigt?) | — |
| 0.25.0 | 2026-06-23 | Change Request (Auftraggeber): neue Anforderung `DC-FA-DIAG-001` — opt-in Modul `diagrams` öffnet gezielt benannte Diagramm-Fences (Default `mermaid`) und prüft die darin gefundenen Kennungen auf **Existenz** in ihrer `defined-in`-Quelle (Befund `diagram-id-undefined`); reine Token-Extraktion ohne Mermaid-Parser, read-only/netzlos (`DC-QA-03`), deterministisch (`DC-QA-02`), Default aus (byte-identisch). Link-Policy gilt in Fences nicht (keine Markdown-Links möglich) → Existenz statt Linkpflicht. Bereich `DIAG` in der Schema-Konvention deklariert; Modul-Liste in `DC-FA-CLI-002` ergänzt. Anlass: belief-agent-Architektur — `ARC-NN`/`LH-*`-Kennungen in `mermaid`-Diagrammen entgehen heute allen Modulen, weil Fences opak sind | — |
| 0.24.0 | 2026-06-23 | Change Request (Auftraggeber): neue Anforderung `DC-FA-CLI-011` — opt-in `--trace --require-complete` bindet die Waisen-Markierung an den Exit-Code (≥1 Requirements-Waise ⇒ Exit 1, sonst 0); Default-`--trace` bleibt advisory (Exit 0). Mit-Erweiterung `DC-FA-CLI-010`: das `--print-mk`-Fragment trägt zusätzlich `doc-trace` (advisory RTM) und `doc-complete` (Vollständigkeits-Gate) plus eine `TRACE_FLAGS`-Variable. Read-only (`DC-QA-03`), deterministisch (`DC-QA-02`), kein Dokument geschrieben. Anlass: Konsumenten (a-check-Bootstrap) sollen die Vollständigkeits-Invariante als Makefile-Gate binden, ohne die `completeness-check.sh`-Parsing-Logik zu kopieren | — |
| 0.23.0 | 2026-06-22 | Change Request (Auftraggeber): `DC-FA-CODE-001` um das Datei-Ventil `exempt-paths` (Glob-Liste, Syntax wie `scan.ignore`) erweitert — nimmt **ganze Dateien** von der `codepaths`-Prüfung aus, datei-weit und unabhängig von `roots`; dasselbe Ventil wie `DC-FA-ID-001`, komplementär zum zeilenweisen `d-check:ignore`-Marker. Abwärtskompatibel (`DC-QA-02`): ohne gesetztes `exempt-paths` byte-identisch. Anlass: ein Nebenbefund — Review-Reports unter `docs/reviews/` zitieren naturgemäß `Datei:Zeile`/Pfade und lösten `codepath-missing` aus (`ids` exemptet `docs/reviews/**` längst; `codepaths` zieht nach) | — |
| 0.22.0 | 2026-06-21 | Change Request (Auftraggeber): neue Anforderung `DC-FA-CLI-010` — Option `--print-mk` gibt ein include-bares `d-check.mk` (überschreibbare `DCHECK_IMAGE`-Variable mit version-gepinntem Image + `doc-check`-Target) auf stdout aus; read-only (`DC-QA-03`), deterministisch (`DC-QA-02`), kein Dokument geschrieben. Konsumenten `include`-n statt Recipe/Skript zu kopieren; der Image-Ref ist die ins Binary eingebettete Release-Version (Digest via `DCHECK_IMAGE`-Override — Henne-Ei: das Binary kennt seinen eigenen Digest nicht). Anlass: a-check-Bootstrap 2026-06-20 — `d-check.mk` wurde handgepflegt; `--print-mk` verlagert den Pin nach d-check | — |
| 0.21.0 | 2026-06-21 | Change Request (Auftraggeber): neue Anforderung `DC-FA-CLI-009` — Option `--trace` gibt eine **Requirements Traceability Matrix** (Anforderung → referenzierende ADRs/Slices, Waisen-Markierung) auf stdout aus (Default Markdown-Tabelle, optional `--trace --json`/`--yaml` über den format-neutralen Reporter); read-only (`DC-QA-03`), deterministisch (`DC-QA-02`), kein Dokument erzeugt; **Doku-Domäne** (Lastenheft/ADR/Planning), kein Code/keine Go-Toolchain (arch-check bewusst ausgeklammert). Anlass: RTM als d-check-Modus statt separatem Skript (Nutzer-Entscheid 2026-06-20, Prototyp-Beleg) | — |
| 0.20.0 | 2026-06-21 | Change Request (Auftraggeber): `DC-FA-CLI-006` um einen **Anforderungs-Präfix-Parameter** erweitert — `--suggest-config ai-harness[-init]` backt nicht mehr fix d-checks `DC-` ein: `--id-prefix <PREFIX>` setzt es explizit, der Modus `ai-harness` leitet es aus dem Lastenheft ab (eindeutige FA-/QA-Kennung; mehrere verschiedene ⇒ Nutzungsfehler), ohne Angabe/Ableitung erscheint ein markierter Platzhalter `<PREFIX>` + TODO statt eines stillen `DC-`. **Breaking** ggü. 0.18.1 (`ai-harness-init` ohne Präfix lieferte zuvor `DC-`); Begründung in eigener ADR. Anlass: a-check-Bootstrap 2026-06-20 — das Init-Template emittierte d-checks Eigen-Präfix in ein Fremd-Repo | — |
| 0.19.0 | 2026-06-19 | Change Request (Auftraggeber): `DC-FA-CLI-004` um das Ausgabeformat **YAML** (`--yaml`) erweitert — gleiche Struktur wie `--json` (`findings`/`summary`/`exitCode`), volle Parität inkl. `--doctor --yaml` (mit `reasonText`/`fixCandidate`); `--json`+`--yaml` und `--repair`+`--yaml` sind Nutzungsfehler (Exit 2). Serialisierung im report-Adapter braucht `gopkg.in/yaml.v3` dort — die Modul-Import-Regeln werden dafür per Folge-ADR erweitert. Anlass: YAML als lesbareres maschinenlesbares Format neben JSON | — |
| 0.18.1 | 2026-06-19 | Schärfung `DC-FA-CLI-006` (Auftraggeber, vor Release): die ai-harness-Vorlage in **zwei explizite Modi** aufgeteilt (Henne-Ei — nicht aus der Repo-Existenz ableitbar) — `ai-harness-init` (Voll-Kanon, alle Blöcke aktiv; Zielbild für ein leeres/frisches Repo, läuft nach Struktur-Anlage) und `ai-harness` (repo-bewusst, fehlende Blöcke auskommentiert; läuft sofort gegen ein bestehendes Repo). Anlass: der einzelne hybride Modus kommentierte in einem leeren Repo alles aus und taugte nicht als Bootstrap-Vorlage | — |
| 0.18.0 | 2026-06-19 | Change Request (Auftraggeber): `DC-FA-CLI-006` um die reservierte Quelle `ai-harness` erweitert — `--suggest-config ai-harness` schlägt ein an die adoptierte ai-harness-course-Konvention angelehntes `.d-check.yml` vor (kanonische `ids`-Muster, `matrix`-Klassen samt Referenzrichtung, Standard-Modulset, Scan-Scope). Hybrid: strukturelle Konventionen immer, konkrete Pfade repo-bewusst (nur existierende roots; fehlende Artefakte auskommentiert mit Hinweis). Read-only/advisory, deterministisch (`DC-QA-02`); liefert — als Ausnahme zur reinen Quellen-Ableitung — `matrix`-Regeln/`link-policy`/`exempt-paths` aus der bekannten Konvention. Anlass: Adoptions-Start für Harness-Repos ohne manuelles Quellen-Auflisten | — |
| 0.17.1 | 2026-06-19 | Change Request (Auftraggeber): Out-of-Scope von `DC-FA-CLI-008` präzisiert — reparierbar sind nur `id-unlinked` (konservativ) und `target-missing` als eindeutiger Datei-Move (breit; Markdown, im Scan-Bestand eindeutiger Basisname); alle übrigen Befundarten bleiben Befund unter `DC-FA-CLI-007`. VCS-/git-historienbasierte Move-/Rename-Erkennung explizit ausgeschlossen (erweiterte die Eingabe über den gescannten read-only-Baum hinaus, `DC-QA-02`; wäre ein eigenes opt-in-Modul analog `external`). Klarstellung des bestehenden Vertrags — keine Verhaltens-/Akzeptanzkriterien-Änderung; Begründung in eigener ADR festgehalten | — |
| 0.17.0 | 2026-06-18 | Change Request (Auftraggeber): `--doctor` wird mit `--json` kombinierbar — Schärfung `DC-FA-CLI-007` (maschinenlesbare Diagnose: `findings` je Eintrag um `reasonText` und `fixCandidate` erweitert, Gruppierung über das `file`-Feld) und `DC-FA-CLI-004` (Kombinierbarkeit `--doctor`+`--json` definiert statt verboten; `--repair`+`--json` und `--doctor`+`--repair` bleiben Nutzungsfehler). Dieselbe Grund-Klartext-/Fix-Kandidaten-Ableitung, drittes Rendering neben Prosa und Patch. Anlass: Auftraggeber-Wunsch nach maschinenlesbarer Diagnose | — |
| 0.16.0 | 2026-06-18 | Change Request (Auftraggeber): neue Anforderung `DC-FA-CLI-008` (Option `--repair`: unified diff auf stdout, `git apply`-kompatibel, read-only; konservative Stufe als Default mit eindeutig ableitbaren Fixes — v1 v. a. `id-unlinked` → Definitions-Link —, breite Best-Guess-Stufe opt-in, Kennzeichnung review-pflichtiger Hunks auf stderr, damit der Patch `git apply`-rein bleibt). Baut auf den Ausgabe-Modi aus 0.15.0 auf; deterministisch (`DC-QA-02`), read-only-Kernvertrag (`DC-QA-03`); In-place-Schreiben bleibt Out-of-Scope (wäre eigene Anforderung + ADR) | — |
| 0.15.0 | 2026-06-18 | Change Request (Auftraggeber): neue Anforderung `DC-FA-CLI-007` (Option `--doctor`: erklärende, nach Datei/Regel gruppierte Diagnose mit Fix-Kandidaten; read-only, stdout-only, deterministisch `DC-QA-02`/`DC-QA-03`). Mit-Schärfung `DC-FA-CLI-003` (die Ausgabe-Modi folgen den Codes 0/1/2, Ausgabe auf stdout unabhängig vom Code) und `DC-FA-CLI-004` (die Modi `--doctor`/`--repair` ersetzen das Default-stdout-Format; untereinander und mit `--json` nicht kombinierbar → Nutzungsfehler exit 2; JSON-Varianten out of scope). Anlass: Machbarkeitsfrage des Auftraggebers (2026-06-18) — beratende Diagnose/Reparatur ohne Bruch der Seiteneffektfreiheit | — |
| 0.14.0 | 2026-06-17 | Change Request (Auftraggeber): `DC-FA-MTX-001` — opt-in `allow-supersede-lineage` (+ `supersede-fields`) nimmt die Supersede-Lineage-Kante X → Y von der Status-Prüfung aus: ein ablösendes Dokument darf auf das von ihm abgelöste (inaktive) Dokument verweisen, ohne `matrix-inactive` zu erzeugen, sofern X über ein deklariertes Feld die Ablösung von Y benennt (Match über Linktext bzw. Zielpfad der Referenz). Carve-out beschränkt auf die deklarierte Lineage-Kante; `matrix-forbidden` (Klassen-Regeln) unberührt. Abwärtskompatibel (`DC-QA-02`): Default aus ⇒ Befundsatz byte-identisch. Marker-Politik (B2 der CR): `matrix` bleibt bewusst ohne Zeilen-Opt-out-Marker `d-check:ignore` (legitime Lineage strukturell ausgenommen, nicht stummgeschaltet). Anlass: reproduzierter Fremd-Repo-Befund (`grid-gym`, v0.10.0) — normative ADR→ADR-Lineage (`Aenderungstyp: Supersedes …`) als `matrix-inactive` gemeldet, brach `make docs-check` | — |
| 0.13.0 | 2026-06-16 | Change Request (Auftraggeber): `DC-FA-ID-001` — Geltungsbereich der beiden Ventile `exempt-paths` und `d-check:ignore` auf **nackte Fließtext-Vorkommen** erweitert. Bisher griffen sie nur auf die `always`-Inline-Code-Vorkommen; eine nackte Kennung in einer `exempt-paths`-Datei (bzw. auf einer `d-check:ignore`-Zeile) wurde weiterhin als `id-unlinked` gemeldet. Neu: beide Ventile sind ein Ganzdatei- bzw. Ganzzeilen-Carve-out für alle Vorkommen des Musters, unabhängig von der `link-policy`. Abwärtskompatibel (`DC-QA-02`): Configs ohne gesetzte Ventile bleiben byte-identisch; Wirkung nur in Richtung *weniger* Befunde in explizit ausgenommenen Dateien/Zeilen. Anlass: reproduzierter Fremd-Repo-Befund (`ai-harness-init`, v0.8.0/v0.9.0) — nacktes `MR-004` in `docs/reviews/_t.md` trotz `exempt-paths: ["docs/reviews/**"]` gemeldet <!-- d-check:ignore (Kennung und Pfad sind der reproduzierte Fremd-Befund, keine Verweise) --> | — |
| 0.12.0 | 2026-06-15 | Change Request (Auftraggeber): `DC-FA-ANCH-001` erweitert — Inline-HTML-Anker innerhalb von Markdown zählen zur gültigen Anker-Menge (`id` an beliebigem Element, `name` an `<a>`; GitHub-Parität, wörtlicher Vergleich). Hebt die bisherige HTML-Anker-Out-of-Scope-Zeile auf; Abgrenzung zu §5 (HTML als Dateiformat bleibt out-of-scope) explizit. Gilt mittelbar auch für `codepaths` (geteiltes Anker-Verfahren). Anlass: Falsch-Befunde `anchor-missing` auf manuell gesetzte HTML-Anker in der Doku | — |
| 0.11.0 | 2026-06-13 | Schärfung `DC-FA-CLI-001` (Auftraggeber): neues Akzeptanzkriterium für `--help` — die Hilfe nennt die Synopsis `d-check [optionen] [pfad]`, beschreibt das Pfad-Argument (Scan-Wurzel, Default cwd) und verweist für das Config-Format auf `--print-config` (kein Format-Duplikat). Anlass: die nackte `flag`-Default-Usage verschwieg das Pfad-Argument | — |
| 0.10.0 | 2026-06-13 | Change Request (Auftraggeber): neue Anforderung `DC-FA-CLI-006` — Option `--suggest-config` leitet die `ids`-Config aus benannten Autoritäts-Dokumenten ab (definierte Kennungen → Muster + target, Round-Trip-Garantie; opt-in-Module nach Signal; Ausgabe-only, read-only). Dazu Schärfung von `DC-FA-ID-001`: das Muster-Ableitungs-Out-of-Scope gilt für die **Prüfung**; ein advisory Scaffold-Modus darf aus benannten Autoritäts-Quellen ableiten. Anlass: Adoptions-Reibung — neue Repos brauchen einen treffsicheren Config-Start statt eines rein statischen Gerüsts | — |
| 0.9.0 | 2026-06-13 | Change Request (Auftraggeber): neue Anforderung `DC-FA-CLI-005` — Option `--print-config` gibt ein statisches, kommentiertes `.d-check.yml`-Startgerüst auf stdout aus (kein Repo-Zugriff, kein Schreiben — read-only-Vertrag bleibt; Umleiten via `> .d-check.yml` macht der Aufrufer). Anlass: Adoptions-Reibung in neuen Repos ohne Config; macht zugleich die verfügbaren Optionen sichtbar. Ableitung aus Repo-Inhalt bewusst nicht Teil (späterer eigener Modus) | — |
| 0.8.0 | 2026-06-13 | Change Request (Auftraggeber): `DC-FA-ID-001` um konfigurierbare `link-policy: prose\|always` (je Muster) erweitert — `always` macht auch Inline-Code-Vorkommen linkpflichtig, Default `prose` (opt-in, abwärtskompatibel). Zwei Ventile: `exempt-paths` (Glob-Liste je Muster) und der Zeilen-Marker `d-check:ignore`, dessen Geltungsbereich von `codepaths`-only auf `ids` erweitert wird (illustrative Beispiel-IDs). Anlass: der `ids`-Sensor maß bislang nicht das Ziel „gut verlinkt" — ein Code-Span konnte stillschweigend einen fehlenden Link verbergen (Ausgangsbefund `DC-QA-03`). Kalibrierung über die drei `ids`-Repos (d-check 155, u-boot 9, b-trace 2) bestimmte die Ventil-Form | — |
| 0.7.2 | 2026-06-12 | Review zum Change Request `DC-FA-HOST-001` (Kalibrierungs-Befund): tmp aus der Default-Präfixliste gestrichen — der Kalibrierungslauf zeigte überwiegend legitime Laufzeit-Doku (Log-/Output-Pfade) statt Maschinen-Layout-Leaks; wer tmp prüfen will, konfiguriert es | — |
| 0.7.1 | 2026-06-12 | Review zum Change Request `DC-FA-SPAN-001` (Kalibrierungs-Befund): Bildreferenzen als Linktext (`[![…](…)](…)` — Badge-Muster, z. B. Shields in vendorten Paket-READMEs) sind legales Markdown und kein `span-nested-link`-Treffer | — |
| 0.7.0 | 2026-06-12 | Change Request (Erst-Bedarfsträger grid-gym): neue Anforderung `DC-FA-CONF-002` — optionaler modul-lokaler Scan-Scope `<modul>.scope` (`roots`/`ignore`), ersetzt für genau dieses Modul den globalen Scope; additiv und abwärtskompatibel, Constraints spiegeln `scan.*`. Anlass (gemessen, v0.2.0): `ids`-Aktivierung in grid-gym liefert global 2776 Befunde (Masse: Retro-Verlinkung historischer done-Planning-Docs = Umschreiben des Audit-Trails) vs. 312 echte, fixbare Befunde im kuratierten Scope `spec/` + `docs/user/` — der Linkpflicht-Nutzen ist heute nur um den Preis des globalen Sweeps oder des Verzichts auf breite `links`/`anchors`-Abdeckung zu haben | — |
| 0.6.0 | 2026-06-12 | Change Request (Auftraggeber): neue Anforderung `DC-FA-HOST-001` — Modul `hostpaths` meldet host-lokale absolute Pfade in Prosa und Inline-Code (opt-in; Unix-Präfixliste konfigurierbar, Windows-/UNC-Muster fest; Fences ausgenommen, kein Opt-out-Marker). Anlass: der bess-ems-Rest-Sensor generalisiert (dort als eigenes Tool gebaut) plus die 8 Host-Pfad-Links aus dem d-migrate-Vergleichslauf und die eigene 0.2.2-Hygiene-Korrektur — dieselbe Leak-Klasse dreimal unabhängig. Bereich `HOST` deklariert; Modul-Listen in `DC-FA-CLI-002` und Glossar ergänzt | — |
| 0.5.0 | 2026-06-12 | Change Request (Auftraggeber): neue Anforderung `DC-FA-SPAN-001` — Modul `spans` meldet ungeschlossene Code-Spans und verschachtelte Link-Artefakte (opt-in, konservative Erkennung: alleinstehende Backticks bleiben befundfrei; kein Opt-out-Marker). Anlass: die `DC-QA-04`-Vergleichsläufe — ~100 u-boot-Artefakte sowie Fälle in grid-gym und bess-ems wurden nur indirekt als `id-unlinked`-Folgefehler sichtbar; ein direkter Sensor meldet die Ursache statt der Symptome. Bereich `SPAN` in der Schema-Konvention deklariert; Modul-Listen in `DC-FA-CLI-002` und Glossar ergänzt | — |
| 0.4.0 | 2026-06-12 | Change Request (Auftraggeber): Inventur-Nachtrag zu `DC-QA-04` — dreizehntes Alt-Tool-Vorkommen entdeckt (`check_markdown_links.py` in bess-ems, eigenständige Python-Linie, entstanden 2026-05-24 und damit vor der Inventur vom 2026-06-10 übersehen); Anforderungs-Text, Einleitung, Stakeholder-Tabelle und Glossar von zwölf auf dreizehn fortgeschrieben. Messmethode (drei Familien-Piloten) unverändert | — |
| 0.3.1 | 2026-06-11 | Review R1 zum Change Request: `DC-FA-CODE-001` präzisiert — Wert-Normalisierung (Anführungszeichen, schließende Satzzeichen) und Anker-Prüfung bei Markdown-Zielen (`DC-QA-04`-Parität zur JS-Familie); `DC-FA-LINK-001`-Inline-Code-Aussage auf Modul-Bezug eingegrenzt | — |
| 0.3.0 | 2026-06-11 | Change Request (Auftraggeber): neue Anforderung `DC-FA-CODE-001` — Modul `codepaths` prüft explizite Pfade in Inline-Code (opt-in, konservative Erkennung, Zeilen-Opt-out `d-check:ignore` nur für dieses Modul). Anlass: `DC-QA-04`-Vergleichslauf gegen die JS-Familie (`docs-check.js`) zeigte die Prüfklasse als Konsolidierungs-Lücke. Bereich `CODE` in der Schema-Konvention deklariert; Modul-Listen in `DC-FA-CLI-002` und Glossar ergänzt | — |
| 0.2.3 | 2026-06-11 | Redaktionell: „(folgt)" in der `DC-QA-01`-Messmethode entfernt — die Benchmark-Definition existiert in der Spezifikation; keine inhaltliche Änderung | — |
| 0.2.2 | 2026-06-10 | Redaktionell: absolute Workspace-Pfade entfernt („Schwester-Repositories des Entwicklungs-Workspace" statt konkreter Pfade); keine inhaltliche Änderung | — |
| 0.2.1 | 2026-06-10 | Redaktionell: Beispiel-Kennungen in DC-FA-ID-001/DC-FA-MTX-001/Glossar auf fiktive Nummern (`ADR-0042`, `ADR-0099`) umgestellt — Kollision mit real entstandenen/zukünftigen eigenen ADRs vermeiden; keine inhaltliche Änderung | — | <!-- d-check:ignore (Beispiel-ID, fiktiv) -->
| 0.2.0 | 2026-06-10 | Review-Runde R1: Modul-Schnitt `links`/`anchors` präzisiert (Fragment-Zuständigkeit, fehlende Zieldatei), Slug-Duplikat-Reihenfolge, Symlink-Vorrang, RFC-3986-Dekodierung vor Escape-Prüfung, Redirect-Regel `external`, Muster-Präzedenz `ids`, Status-Default `matrix`, Scan-Wurzel- und Config-Vollvalidierung, Out-of-Scope Reference-Style-Links, Image-Default-Befehl | — |
| 0.1.0 | 2026-06-10 | Initiale Fassung (Konsolidierung von 12 Quell-Tools, Modul-Schnitt, Docker-Distribution) | — |
