# Lastenheft — d-check

**Version:** 0.46.0

**Status:** Draft

**Autor:** pt9912, **Datum:** 2026-06-10.

---

## 1. Zweck und Geltungsbereich

`d-check` ist ein Kommandozeilen-Tool, das Markdown-Dokumentation eines
Repositories auf kaputte Referenzen prüft: lokale Links und
Bildreferenzen, Heading-Anker, Linkpflicht für Anforderungs- und
Entscheidungs-Kennungen sowie Referenzrichtungs-Regeln zwischen
Dokumentklassen (Referenzmatrix). Es konsolidiert dreizehn funktional
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
> Traceability-Sichten), `CONF` (Konfiguration), `DIST`
> (Distribution).

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
- **Hilfe:** Given `d-check --help`, when aufgerufen, then Exit-Code 0 und die Nutzung auf stderr enthält die Synopsis mit `[pfad]`, die Pfad-Argument-Beschreibung und einen Verweis auf `--print-config`.

**Out-of-Scope:** Prüfung mehrerer Repos in einem Aufruf.

---

### DC-FA-CLI-002 — Regelmodul-Auswahl

**Beschreibung:** Die Prüf-Funktionalität ist in benannte Regelmodule
gegliedert: `links`, `anchors`, `ids`, `matrix`, `external`,
`codepaths`, `spans`, `hostpaths`, `diagrams`, `versions`, `pins`,
`immutable`, `vcs`, `commits`, `planning`, `tracked`, `targets`. Ohne Konfiguration sind `links` und `anchors` aktiv. Module werden über
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
Zeile, Pfade relativ zur Scan-Wurzel); Zusammenfassung und Diagnose
gehen auf stderr. Mit `--json` **oder** `--yaml` erfolgt die gesamte
Ausgabe auf stdout als ein maschinenlesbares Dokument (JSON bzw. YAML) mit
mindestens den Feldern `findings` (Liste mit `file`, `line`, `target`,
`rule`, `reason`), `summary` (`filesChecked`, `findingCount`) und
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
ohne `.git`); `vcs`/`commits` brauchen eine Commit-Range (K3) und werden als
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
- **`ai-harness` Auffindbarkeit:** Given `d-check --suggest-config ai-harness` (oder `ai-harness-init`), when der Lauf endet, then nennt die Ausgabe die nicht aktivierten situativen opt-in-Module (`external`, `diagrams`, `versions`, `pins`, `immutable`, `tracked`, `targets`) in einem Kommentar mit Verweis auf `--print-config` und verweist für die Range-Module `vcs`/`commits` auf `--print-mk` (Auffindbarkeit ohne Aktivierung), Exit 0; das Repository wird nicht geschrieben.
- **`ai-harness` Modulset (`spans`/`hostpaths`):** Given ein Repo mit den Harness-Artefakten, when `d-check --suggest-config ai-harness` läuft, then enthält die `modules`-Zeile `spans` und `hostpaths` (Teil der fixen Aktiv-Menge), das Gerüst dekodiert über den eigenen Parser, Exit 0.
- **`ai-harness` `planning` repo-bewusst:** Given ein Repo **mit** `docs/plan/planning/in-progress/roadmap.md`, when `d-check --suggest-config ai-harness` läuft, then ist `planning` in `modules` **und** ein `planning`-Block mit `roadmap:` aktiv; **fehlt** die Roadmap, sind Modul und Block **auskommentiert** (repo-bewusst); im Modus `ai-harness-init` ist `planning` **aktiv** (Voll-Kanon). Exit 0.

**Out-of-Scope:** Schreiben der Datei (immer stdout); Muster-Ableitung aus beliebigem Fließtext (nur aus definierten Headings benannter Autoritäts-Quellen); Garantie eines minimalen/perfekten `regex` (Best-Guess-Generalisierung + Quell-Kennungs-Kommentar — der Mensch verengt); automatisches Ableiten von `link-policy`, `matrix`-Regeln oder `exempt-paths`. Die reservierten Quellen `ai-harness`/`ai-harness-init` sind hiervon ausgenommen — sie liefern `matrix`-Regeln, `link-policy` und `exempt-paths` aus der **bekannten Konvention** (nicht aus dem Repo abgeleitet); sie sind an **eine** adoptierte Baseline-Version gebunden (im Kommentar-Header genannt) — automatische Erkennung oder Hebung der Baseline-Version, das Anlegen der Verzeichnisstruktur (read-only) und das automatische Aktivieren der nicht qualifizierten opt-in-Module (`external`, `diagrams`, `versions`, `pins`, `immutable`, `vcs`, `commits`, `tracked`, `targets`) sind Out-of-Scope — die nicht aktivierten situativen Module nennt die Ausgabe stattdessen im Kommentar mit Verweis auf [`--print-config`](#dc-fa-cli-005--konfigurations-gerüst-ausgeben) bzw. für die Range-Module `vcs`/`commits` auf [`--print-mk`](#dc-fa-cli-010--makefile-fragment-ausgeben) (Auffindbarkeit ohne stilles Aktivieren eines inerten Moduls). Aufgenommen sind dagegen `spans`/`hostpaths` (fixe Aktiv-Menge) und der repo-bewusste `planning`-Block (Eignungs-Kriterium K1–K4, s. o.).

---

### DC-FA-CLI-007 — Diagnose-Modus

**Beschreibung:** Mit der Option `--doctor` macht `d-check` einen
**Lese**-Durchgang wie eine normale Prüfung, gibt aber statt der knappen
Befund-Zeilen ([`DC-FA-CLI-004`](#dc-fa-cli-004--ausgabeformate)) eine
**erklärende, nach Datei und Regel gruppierte Diagnose** auf stdout aus:
je Befund den Grund-Code in Klartext und — wo aus dem Befund **eindeutig
ableitbar** — einen **Fix-Kandidaten** (die vorgeschlagene Änderung,
**nicht angewendet**). Das Werkzeug schreibt niemals selbst
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
[`DC-FA-DIST-001`](#dc-fa-dist-001--docker-image)) — sowie elf
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

- **Happy Path:** Given ein installiertes `d-check`, when `d-check --print-mk` läuft, then liegt auf stdout ein Makefile-Fragment mit `DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v…`, den Variablen `DCHECK_DIGEST` und `TRACE_FLAGS` und den `##`-annotierten Targets `doc-check`, `doc-trace`, `doc-complete`, `doc-doctor`, `doc-repair`, `doc-immutable` (mit `--enable vcs` und einer auf `vcs` fokussierten `--disable`-Liste), `doc-commits` (mit `--enable commits` und einer auf `commits` fokussierten `--disable`-Liste), `doc-planning` (mit `--enable planning` und einer auf `planning` fokussierten `--disable`-Liste, ohne Range), `doc-tracked` (mit `--enable tracked` und einer auf `tracked` fokussierten `--disable`-Liste, ohne Range), `doc-targets` (mit `--enable targets` und einer auf `targets` fokussierten `--disable`-Liste, ohne Range) und `doc-help` (jeweils `docker run … --network none … :/repo:ro`), Exit 0; kein Repo-Zugriff nötig.
- **Boundary:** Given das Fragment wird per `d-check --print-mk > d-check.mk` umgeleitet und in ein Makefile `include`-t, when `make doc-check` (bzw. `doc-trace`/`doc-complete`/`doc-doctor`/`doc-repair`/`doc-immutable`/`doc-commits`/`doc-planning`/`doc-tracked`/`doc-targets`) läuft, then ruft das Target das gepinnte Image (bzw. den per `DCHECK_IMAGE`/`DCHECK_DIGEST` gesetzten Override) im passenden Modus (`doc-trace` → `--trace`, `doc-complete` → `--trace --require-complete`, `doc-doctor` → `--doctor`, `doc-repair` → `--repair` mit unterdrücktem Recipe-Echo, `doc-immutable` → `--enable vcs` + Fokus-`--disable` mit `--range $(RANGE)` bzw. `--staged`, `doc-commits` → `--enable commits` + Fokus-`--disable` mit `--range $(RANGE)`, `doc-planning` → `--enable planning` + Fokus-`--disable`, **hermetisch ohne Range**, `doc-tracked` → `--enable tracked` + Fokus-`--disable`, ohne Range, `doc-targets` → `--enable targets` + Fokus-`--disable`, **hermetisch ohne Range**); d-check selbst schreibt dabei nichts.
- **Negative:** Given `d-check --print-mk` mit einem unbekannten Flag, when aufgerufen, then Exit-Code 2 (Nutzungsfehler, [`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)).

**Out-of-Scope:** Schreiben der `d-check.mk` (immer stdout); Einbetten des eigenen Image-**Digests** ins Binary (Henne-Ei — der Digest hasht das Binary selbst; der Konsument pinnt per `DCHECK_IMAGE`-Override); weitere Targets jenseits der gelisteten elf (`doc-check`/`doc-trace`/`doc-complete`/`doc-doctor`/`doc-repair`/`doc-immutable`/`doc-commits`/`doc-planning`/`doc-tracked`/`doc-targets`/`doc-help` — Konsumenten komponieren weitere `gates` selbst); ein `help`-Target (Namens-Kollision mit dem Konsumenten — daher namespaced `doc-help`); Nicht-`@sha256:`-Digest-Formen in `DCHECK_DIGEST`; die Exit-Code-Semantik der RTM-Targets selbst (in [`DC-FA-CLI-011`](#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code) bzw. [`DC-FA-CLI-009`](#dc-fa-cli-009--requirements-traceability-matrix) festgelegt); Nicht-Make-Build-Systeme.

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
als Fließtext).

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Link auf eine existierende Datei im Repo, when das Modul `links` läuft, then kein Befund.
- **Boundary:** Given ein Link auf eine nicht existierende Datei innerhalb eines Fenced-Code-Blocks, when das Modul läuft, then kein Befund.
- **Negative:** Given ein Link `../../etc/passwd`, dessen Ziel die Repository-Wurzel verlässt, when das Modul läuft, then ein Befund mit Grund „verlässt Repository" — auch wenn der Zielpfad existiert. <!-- d-check:ignore (AK-Beispiel: Angriffs-Pfad) -->
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
- **`always` Negative:** Given `link-policy: always` und ein Vorkommen `` `ADR-0042` `` ohne Link (außerhalb `exempt-paths` und ohne `d-check:ignore`), when das Modul läuft, then ein Befund `id-unlinked`.
- **`always` Boundary (Ventile):** Given `link-policy: always`, ein `` `ADR-0042` `` in einer `exempt-paths`-Datei und ein zweites `` `ADR-0099` `` auf einer Zeile mit `d-check:ignore`, when das Modul läuft, then kein Befund für beide.
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
erkennt. Tritt im Körper eines Dokuments der Klasse A ein `token` der Klasse B
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
nicht geprüft. Ein HTML-Kommentar `d-check:ignore` auf derselben
Zeile (Begründung in Klammern empfohlen) nimmt die Zeile von **genau
dieser** Prüfung aus — bewusst nicht existierende Beispiel-Pfade
(etwa Angriffs-Beispiele in Lehrtexten) sind kein Fehler, sondern
eine zu dokumentierende Absicht. Für alle anderen Module existiert
kein Opt-out-Marker: deterministische Befunde werden behoben, nicht
unterdrückt. Zusätzlich nimmt eine optionale `exempt-paths`-Glob-Liste
(Syntax wie `scan.ignore`, relativ zur Repository-Wurzel) **ganze Dateien**
von der `codepaths`-Prüfung aus — datei-weit, unabhängig von den
Wurzel-Präfixen; dasselbe Datei-Ventil wie bei
[`DC-FA-ID-001`](#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids),
komplementär zum zeilenweisen `d-check:ignore`-Marker (typischer Fall:
Review-Reports, die naturgemäß `Datei:Zeile`/Pfade zitieren). Schließlich nimmt
eine optionale `ignore-refs`-Glob-Liste **bestimmte Ziel-Pfade** von der
Existenz-Prüfung aus — **referenz-weit**, egal in welcher Datei/Zeile der Pfad als
Inline-Code-Verweis vorkommt. Sie ist ein Register **bewusst entfernter/historischer
Artefakte** (Tombstones): immutable/historische Doku (z. B. eine `Accepted`-ADR, die
nicht editierbar ist) darf einen entfernten Pfad weiter zitieren, ohne dass er
existieren muss. Bewusster Akt **mit Gate** — ohne Eintrag meldet ein gelöschter,
noch zitierter Pfad weiter `codepath-missing` (nichts dangelt still); `ignore-refs`
unterdrückt **nur** die Existenz-/Anker-Prüfung des genannten Pfads, keine anderen
Befunde. Die drei Ventile decken drei Achsen ab: Zeile (`d-check:ignore`), Datei
(`exempt-paths`), Ziel-Pfad (`ignore-refs`).

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Inline-Code-Span `` `docs/plan/adr/` `` auf ein existierendes Verzeichnis und das konfigurierte Präfix `docs/`, when das Modul `codepaths` läuft, then kein Befund.
- **Boundary:** Given ein Inline-Code-Span mit nicht existierendem Pfad und ein Kommentar `d-check:ignore` auf derselben Zeile, when das Modul läuft, then kein Befund — und der Marker hat keinerlei Wirkung auf Befunde anderer Module derselben Zeile.
- **Negative:** Given ein Inline-Code-Span `` `../fehlt.md` ``, dessen Ziel nicht existiert (oder die Repository-Wurzel verlässt), when das Modul läuft, then ein Befund mit Datei, Zeile, Ziel und Grund, Exit-Code 1.
- **Ventil (exempt-paths):** Given ein Inline-Code-Span mit nicht existierendem Pfad in einer Datei, die ein `codepaths.exempt-paths`-Glob matcht, when das Modul läuft, then kein Befund — und ohne gesetztes `exempt-paths` ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)).
- **Ventil (ignore-refs):** Given ein Inline-Code-Span auf einen nicht existierenden Pfad, der ein `codepaths.ignore-refs`-Glob matcht, when das Modul läuft, then **kein** Befund — referenz-weit, auch in einer sonst geprüften Datei; ein nicht von `ignore-refs` gedeckter fehlender Pfad bleibt `codepath-missing`, und ohne gesetztes `ignore-refs` ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)).

**Out-of-Scope:** Pfad-Erkennung im Fließtext ohne Inline-Code; Pfade in Fenced-Code-Blöcken; Opt-out-Marker für andere Module; semantische Prüfung, ob der referenzierte Pfad inhaltlich passt; Zeilen-Granularität von `exempt-paths` (datei-weit wie `ids`; für Zeilen der `d-check:ignore`-Marker).

---

### DC-FA-SPAN-001 — Markdown-Span-Artefakte (Modul `spans`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `spans` werden zwei
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

Der Befund nennt Datei, Zeile (Opener- bzw. Muster-Zeile), die
betroffene Backtick-Folge bzw. das Muster und den Grund. Es gilt die
allgemeine Opt-out-Regel unverändert: deterministische Befunde werden
behoben, nicht unterdrückt — der Zeilen-Marker `d-check:ignore`
bleibt auf das Modul `codepaths` beschränkt
([`DC-FA-CODE-001`](#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)).

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Absatz mit ausschließlich balancierten Code-Spans — auch solchen, die einen Zeilenumbruch enthalten —, when das Modul `spans` läuft, then kein Befund.
- **Boundary:** Given eine alleinstehende literale Backtick-Folge (beidseitig Whitespace), when das Modul läuft, then kein Befund.
- **Negative:** Given ein Listenpunkt, dessen öffnende Backtick-Folge unmittelbar vor Nicht-Whitespace steht und im Absatz ungeschlossen bleibt, when das Modul läuft, then ein Befund `span-unclosed` mit Datei, Zeile und Grund, Exit-Code 1.

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

**Akzeptanzkriterien:**

- **Happy Path:** Given `diagrams` aktiv (`fences: [mermaid]`) und ein `mermaid`-Block, dessen Kennungen alle in ihrer `defined-in`-Quelle vorkommen, when `d-check` läuft, then kein Befund, Exit 0; ein read-only gemountetes Repository genügt.
- **Boundary:** Given ein `mermaid`-Block mit genau einer Kennung ohne Definition in `defined-in`, when das Modul läuft, then genau ein `diagram-id-undefined` (Datei, Zeile im Fence, Kennung), Exit 1; dieselbe Kennung in einem **nicht** gelisteten Fence (z. B. `bash`) oder außerhalb jedes Fence bleibt für dieses Modul unberührt.
- **Negative:** Given **kein** `diagrams`-Block in der Konfiguration, when `d-check` läuft, then ist der Befundsatz byte-identisch zum Lauf ohne das Modul ([`DC-QA-02`](#dc-qa-02--determinismus)) und es wird nichts geschrieben ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

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
Muster-Scan). Zwei Ventile wie bei
[`DC-FA-ID-001`](#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids): die
Glob-Liste `exempt-paths` (historische Pins in Planning-`done/`-Slices,
`CHANGELOG.md` und der Lastenheft-Historie bleiben unberührt) und der
Zeilen-Marker `d-check:ignore`.

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
- **Boundary (Ventile):** Given ein Pin `…:v0.1.0` in einer `exempt-paths`-Datei
  (z. B. ein Planning-`done/`-Slice) und ein zweiter Pin auf einer Zeile mit
  `d-check:ignore`, when `d-check` läuft, then kein Befund für beide.
- **Modul-aus:** Given **kein** `versions`-Block in der Konfiguration, when
  `d-check` läuft, then ist der Befundsatz byte-identisch zum Lauf ohne das
  Modul ([`DC-QA-02`](#dc-qa-02--determinismus)) und es wird nichts geschrieben
  ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

**Out-of-Scope:** `--repair`/Auto-Bump von `version-stale` in dieser Version
(deterministisch ableitbar, aber Folge-CR an
[`DC-FA-CLI-008`](#dc-fa-cli-008--reparatur-patch)); Ableiten der aktuellen
Version aus Git-Tags (außerhalb des read-only gemounteten Baums, bräche
[`DC-QA-02`](#dc-qa-02--determinismus)); mehrere unabhängige Versions-Reihen pro
Repo in der ersten Fassung (genau ein `pin-pattern` und ein `current-from`);
semantische Versions-Ordnung (nur Gleichheit, keine „neuer als"-Prüfung);
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
nächsten gleich-/höherrangigen) —, normalisiert ihn **whitespace-/reflow-invariant**
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

**Akzeptanzkriterien:**

- **Happy Path:** Given `planning` aktiv mit gesetzter `planning.roadmap`, einer Roadmap mit dem kanonischen `## Aktuelle Welle`-Abschnitt **ohne** den Ruhe-Marker und ≥1 `slice-*` im Verzeichnis, when `d-check --enable planning` läuft, then kein Befund, Exit 0 (ebenso konsistent: Ruhe-Marker vorhanden **und** kein Slice).
- **Boundary (Modul-aus):** Given **kein** aktives `planning`, when `d-check` in einer netzlosen, read-only Umgebung läuft, then ist der Befundsatz byte-identisch ([`DC-QA-02`](#dc-qa-02--determinismus)) und nichts wird geschrieben ([`DC-QA-03`](#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **Negative:** Given `planning` aktiv und eine Drift (Slice im Verzeichnis, aber die Roadmap trägt den Ruhe-Marker; oder kein Slice, aber die Roadmap benennt eine aktive Welle), when `d-check --enable planning` läuft, then ein Befund `planning-drift`, Exit 1.
- **fail-closed (Heading/Datei fehlt):** Given `planning` aktiv und eine Roadmap **ohne** die kanonische `planning.heading`-Überschrift (bzw. eine fehlende `planning.roadmap`-Datei), when `d-check --enable planning` läuft, then `planning-drift`, Exit 1 — kein stilles Grün.

**Out-of-Scope:** eine git-/VCS-basierte Lifecycle-Prüfung (rein hermetisch, nur
Arbeitsbaum); mehr als eine Roadmap bzw. ein Slice-Verzeichnis pro Lauf; das Prüfen
des Slice-**Inhalts** (nur Datei-Existenz, wie `codepaths`); die Roadmap-Prosa
jenseits des Aktiv-Status-Markers; ein `--repair`-Hunk; das Erzwingen der
Lifecycle-Move-Commit-Bündelung selbst ([`MR-013`](../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)
bleibt Commit-Zeit-Disziplin).

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

- **Happy Path:** Given eine `.d-check.yml` mit einem Ignorier-Muster `docs/archive/**`, when `d-check` läuft, then erzeugen kaputte Links unter `docs/archive/` keine Befunde. <!-- d-check:ignore (AK-Beispiel) -->
- **Boundary:** Given keine `.d-check.yml`, when `d-check` läuft, then laufen genau die Default-Module (`links`, `anchors`) auf den Default-Scan-Wurzeln.
- **Negative:** Given eine syntaktisch ungültige `.d-check.yml`, when `d-check` läuft, then Exit-Code 2 mit Fehlermeldung inkl. Zeilenangabe; es wird keine Prüfung mit stillschweigenden Defaults durchgeführt.

**Out-of-Scope:** Konfigurations-Vererbung über mehrere Verzeichnisebenen; Migrations-Assistent von Alt-Tool-Konfigurationen.

---

### DC-FA-CONF-002 — Modul-lokaler Scan-Scope

**Beschreibung:** Jedes Regelmodul akzeptiert in der
Konfigurationsdatei optional einen Schlüssel `<modul>.scope` mit den
Unterschlüsseln `roots` und `ignore`. Ist er gesetzt, **ersetzt** er
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

**Out-of-Scope:** Distributionswege jenseits GHCR (Homebrew, Paketmanager); deren Bewertung erfolgt später per ADR.

## 4. Nichtfunktionale Anforderungen

### DC-QA-01 — Performance

- **Anforderung:** Ein vollständiger Lauf der Default-Module über ein Repository mit 1.000 Markdown-Dateien (insgesamt ≤ 20 MB) terminiert in < 5 s auf Standard-Entwickler-Hardware (2 vCPU); das Modul `external` ist von dieser Schranke ausgenommen.
- **Messmethode:** Benchmark-Lauf gegen ein generiertes Fixture-Repo, Definition in `spec/spezifikation.md`.

### DC-QA-02 — Determinismus

- **Anforderung:** Identische Eingabe (Repo-Stand + Konfiguration + Optionen) liefert byte-identische Ausgabe; Befunde sind stabil sortiert (Pfad, dann Zeilennummer).
- **Messmethode:** Test führt denselben Lauf 10× aus und vergleicht Hashes der Ausgabe.

### DC-QA-03 — Seiteneffektfreiheit und Netzwerk-Sparsamkeit

- **Anforderung:** Das Tool schreibt nie in das geprüfte Repository und öffnet außer im explizit aktivierten Modul `external` keine Netzwerkverbindungen.
- **Messmethode:** Integrationstest mit read-only-Mount und netzwerkloser Umgebung (`docker run --network none`), alle Module außer `external` und `vcs` aktiv (beide brauchen eine explizite Eingabe jenseits des gescannten Markdown-Baums — `external` das Netz, `vcs` eine git-Range — und sind nie Teil des netzlosen Default-Laufs; `vcs` liest `.git` zwar read-only/netzlos, fail-closed ohne Range).

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
| Regelmodul | Benannte, einzeln aktivierbare Prüf-Einheit (`links`, `anchors`, `ids`, `matrix`, `external`, `codepaths`, `spans`, `hostpaths`, `diagrams`, `versions`, `pins`, `immutable`, `vcs`, `commits`, `planning`, `tracked`, `targets`). |
| Scan-Wurzel | Verzeichnis, unterhalb dessen Markdown-Dateien gesucht werden; zugleich Bezugspunkt der Pfadauflösung. |
| Anker | Fragment-Teil eines Links (`#…`), das auf ein Heading der Zieldatei zeigt (GitHub-Slug-Verfahren). |
| Repo-Escape | Linkziel, dessen aufgelöster Pfad außerhalb der Repository-Wurzel liegt. |
| Kennung | Textuelle ID nach deklariertem Muster (z. B. `ADR-0042`), für die Linkpflicht gelten kann. | <!-- d-check:ignore (Beispiel-ID, fiktiv) -->
| Dokumentklasse | Über Pfad-Muster definierte Gruppe von Dokumenten (z. B. Contract-Spec, ADR, Slice) als Knoten der Referenzmatrix. |
| Referenzmatrix | Deklaration, welche Dokumentklasse auf welche verweisen darf, inkl. Status-Bedingungen. |
| Aktives ADR | ADR, dessen Status-Feld keinen verbotenen Wert (`superseded`, `deprecated`) trägt. |
| Quell-Tools | Die dreizehn konsolidierten Alt-Tool-Vorkommen in den Schwester-Repositories des Entwicklungs-Workspace: zwölf aus drei Familien (Shell: `verify-doc-refs.sh`, Python: `check_refs.py`, JavaScript: `docs-check.js`) plus eine eigenständige Python-Linie (`check_markdown_links.py`, Inventur-Nachtrag 2026-06-12). |

## 7. Historie

| Version | Datum | Änderung | Verweis |
|---|---|---|---|
| 0.46.0 | 2026-07-17 | Change Request (Auftraggeber): die **Komma-Kurzform** `<FAM>-AAA, BBB` wird **fail-closed** (Exit 2) statt still verschluckt — betrifft `DC-FA-COV-001` und, über den geteilten Reader, `DC-FA-XREF-001`. Zugesagt waren nur `..`-Range und `/`-Aufzählung; die Kurzform war nie Vertrag, ihr **stiller Drop** aber auch nicht: `GG-SCN-001, 007` deckte nur `GG-SCN-001` und erzeugte eine falsche Waise. Verworfen: (a) Komma-Enum unterstützen — d-check kann eine gemeinte Kurzform nicht von einer Zahl im Fließtext unterscheiden (`GG-QA-001, 007 Sekunden`), das wäre Raten; (b) Status quo — der stille Drop ist die schlechteste Option. Gewählt: die **Gestalt** triggert (Kennung + Komma + Ziffern), der Inhalt ist keine unterstützte Notation ⇒ Exit 2 mit Hinweis — dieselbe Logik wie bei `AAA>BBB`. Ein Komma vor einer **vollständigen** Kennung (`GG-AR-COMP-CORE, GG-AR-COMP-DOMAIN`) bleibt unberührt. Neues Negative- und Boundary-Kriterium; Schritt 3 in `DC-FA-COV-001.a`; Begründung in der begleitenden ADR. Anlass: Konsumenten-Report grid-gym gegen v0.45.1 — reale §27.1-Zeilen `GG-SCN-001, 007` ließen das produktiv verdrahtete `trace.coverage` das Mapping nicht zählen | slice-075 |
| 0.45.0 | 2026-07-17 | Change Request (Auftraggeber): `DC-FA-XREF-001` um **`forward.req-pattern`** erweitert (RE2, Default `requirements.id-pattern`) — symmetrisch zum vorhandenen `backward.req-pattern`, ein Denkmodell statt zwei. Bis dahin las die Vorwärts-Sicht ihre Anforderungs-IDs **still** über `requirements.id-pattern`, die Rück-Sicht über ihr eigenes Muster; die Kopplung stand weder im Vertrag noch in der Config-Oberfläche. Folge bei einer bewusst gescopten RTM (Architektur-Meta ausgeschlossen): die Vorwärts-Sicht ist leer, **jede** Rück-Kante wird als „ohne RTM-Eintrag“ gemeldet, und die eigentliche `F \ B`-Drift verschwindet — ein Falschbefund, der wie echter Drift aussieht. Neu festgehalten: **die Vergleichs-Schlüsselmenge ist nicht die RTM-Anforderungsmenge** (belegt: eine ID, die das Muster trifft, aber keine RTM-Zeile hat, wird verglichen). Neues Boundary-Kriterium (RTM-Scope ≠ Vergleichs-Scope); Schema-Zeile + Algorithmus-Schritt 2 in `DC-FA-XREF-001.a`; Begründung in der begleitenden ADR. Anlass: Realdaten-Lauf des Konsumenten grid-gym gegen v0.44.0 (Defekt 1 von zweien), belegt durch Umschalten **nur** dieses Musters bei identischen Dateien | slice-071 |
| 0.44.2 | 2026-07-17 | Change Request (Auftraggeber): die Vakuitäts-Stufe aus 0.44.1 von einer **Ursachen-** auf eine **Wirkungs**-Fassung gezogen. 0.44.1 band Vakuum an die Muster-Ursache („Given ein `design-pattern`, das … vorbeigreift“); ein **übergriffiges `exclude-req`** — fehlerfreie Muster, aber ein Ventil, das jede Anforderung verschluckt — schaltete das Gate bei realem Drift ebenso still ab, war aber von keinem Akzeptanzkriterium gedeckt. Vakuität wird daher **nach** dem Ausschluss gemessen (das Ventil ist selbst eine kuratierte, drift-fähige Kante), und geprüft wird die **Wirkung** — ob der Abgleich konstruktionsbedingt überhaupt einen Befund liefern könnte — statt einer Ursachenliste, die bei der nächsten unbekannten Ursache erneut risse. Neues Negative-Kriterium (Ventil, samt ventil-benennender Meldung); Ausschluss-Stufe in der Fehlerpräzedenz von `DC-FA-XREF-001.a`; Begründung in der begleitenden ADR. Anlass: unabhängiges Closure-Review zu `slice-071` — R3 reproduzierte `exclude-req: '.'` mit `0 Differenz(en)`/Exit 0 bei echtem Drift, R4 wies die AK-Lücke nach | slice-071 |
| 0.44.1 | 2026-07-17 | Change Request (Auftraggeber): `DC-FA-XREF-001` um eine **Vakuitäts-Stufe** geschärft — ein Abgleich, aus dem keine Kante fällt, ist **kein bestandener** Abgleich, sondern Exit 2. Vakuum ist definiert als **beide** Sichten kantenleer (typischer Anlass: ein `design-pattern`, das kompiliert, aber am Artefakt-Namensraum vorbeigreift — die Namensraum-Kongruenz war bis dahin nur als Vorbedingung *beschrieben*, nicht mechanisiert) **oder** die Rück-Sicht kantenleer unter `mode: superset` (dann kann `B \ F` konstruktionsbedingt nie einen Befund liefern). Abgegrenzt: eine **einseitig** leere Vorwärts-Sicht bleibt ein wohldefiniertes Ergebnis (Diff über `keys(F) ∪ keys(B)`) und meldet `B \ F` laut — der erwartete Bootstrap-Zustand vor der Restrukturierung der Vorwärts-Sicht, den ein Guard sonst mit einer Config-Fehldiagnose abwürgte. Neues Negative- und Boundary-Kriterium; Vakuitäts-Stufe in der Fehlerpräzedenz von `DC-FA-XREF-001.a` der Spezifikation; Begründung + Abgrenzung zum Nullmengen-Guard der Tabellenquellen in der begleitenden ADR. Anlass: unabhängiges Closure-Review zu `slice-071` — R1 reproduzierte das stille Grün, R2 wies den ersten, symmetrisch je Sicht feuernden Fix als vertragswidrig nach | slice-071 |
| 0.44.0 | 2026-07-16 | Change Request (Auftraggeber): neue Anforderung `DC-FA-XREF-001` — **Kreuzverweis-Konsistenz zweier Traceability-Sichten** (`trace.cross-consistency`, opt-in): der `--trace`-Lauf vergleicht zusätzlich die Vorwärts-RTM-Tabelle (Anforderung → Design-Menge) gegen die Rückwärts-`Bezug`-Kanten (Design → Anforderung) und meldet je Anforderung die beiden Mengendifferenzen `F(R) \ B(R)`/`B(R) \ F(R)` mit Richtungslabel; Modus `equal`/`superset`, `1:N` Normalfall. Beide Sichten kuratierte Tabellen über den header-gebundenen Reader (`DC-FA-REQ-001`) + range-aware Span-Semantik (`DC-FA-COV-001`); Rück-Kanten = Quelle der Wahrheit (Artefakt-ID = erste Spalte via `design-pattern`, `Bezug`-Zelle range-aware). Ableitungssprünge per `exclude-req` (RE2) ausgenommen. Advisory unter `--trace`; Exit-Änderung nur über das globale `--require-complete` (`DC-FA-CLI-011`). Fail-closed (ungültiges Regex/fehlende Spalte/ID-Header nicht genau einmal/unbekannter `mode` ⇒ Exit 2); ohne `trace.cross-consistency`-Block RTM byte-identisch (`DC-QA-02`/`DC-QA-03`). Bereich `XREF` in §3; additive Erweiterung des `--trace`-Laufs (`DC-FA-CLI-009`) **ohne** RTM-Änderung (separate advisory Befunde, keine RTM-Spalte); `DC-FA-XREF-001.a` + Schema-Keys (`trace.cross-consistency.*`) in der Spezifikation; Implementierung, Realdatenbeleg und Release folgen in `slice-071`. Anlass: Konsument grid-gym (Trigger 088) — §27.1 nannte {GG-AR-COMP-CORE, GG-AR-COMP-DOMAIN}, die `Bezug`-Rück-Kanten {GG-AR-P-005, GG-AR-P-009, GG-AR-COMP-SCHED}, Schnittmenge null, von keinem Gate bemerkt | slice-071 |
| 0.43.0 | 2026-07-14 | Change Request (Auftraggeber): neue Anforderung `DC-FA-REQ-001` — **native tabellarische Anforderungsquellen + Nullmengen-Guard** für die RTM. `trace.requirements.format: table` liest Markdown-Pipe-Tabellen über konfigurierte Header-Namen (`table.id-column`, genau eine von `table.text-column`/`.text-columns`, optional `.modality-column`); die ID-Zelle muss als Ganzes auf `id-pattern` passen, die gebundene Textzelle liefert RTM-Titel/Body, eine gesetzte Modalitätsspalte ist die alleinige Keyword-Quelle für `DC-FA-MOD-001`. `table.duplicate-ids` ist `error` (Default), `first` oder `last`. Beide Formate münden in dasselbe RTM-Modell. **Fail-closed:** nichtleer explizite `requirements.source` oder Tabellenmodus mit null erkannten Anforderungen sowie unbekanntes Format/fehlende oder doppelte konfigurierte Spalten ⇒ Exit 2 (auch mit `--require-complete`) statt irreführendem `0 Anforderungen, 0 Waisen`/Exit 0. `source: ""` gilt wie abwesend. Ohne `trace`-Block bleibt der Heading-Default samt Deduplizierung/leerer RTM/Exit 0 byte-identisch (`DC-QA-02`/`DC-QA-03`). Mit-Modifikation `DC-FA-CLI-009`; Bereich `REQ` in §3; Spezifikation, ADR, Implementierung, Realdatenbeleg und Release folgen in `slice-070`. Anlass: reproduzierter Konsumentenbefund `m-trace` mit d-check v0.42.0 — 371 eindeutige IDs in Tabellen mit den Text-Headern `Anforderung`/`Akzeptanzkriterium`, explizite Quelle und passende Regex ergaben null Anforderungen bei Exit 0 | slice-070 |
| 0.42.0 | 2026-07-11 | Change Request (Auftraggeber): neue Anforderung `DC-FA-MOD-001` — **Modalitäts-Klassifikation der Anforderungen** (`trace.requirements.modality`, opt-in): d-check klassifiziert jede RTM-Anforderung anhand **konfigurierbarer Modal-Verb-Schlüsselwörter** (Built-in DE+EN-RFC-2119-Defaults; `modality.levels` Stufe→Keywords, `modality.require-levels` welche Stufen gaten, Default `[must]`) über den **ersten** Treffer im Anforderungs-Body (Längster-Treffer-zuerst gegen `MUSS NICHT`≠`DARF NICHT`, case-insensitiv/wortgrenzen-genau; kein Treffer ⇒ Stufe `unknown`). Neue **Modality-Spalte** (konditional); `--require-complete` ([`DC-FA-CLI-011`](#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)) bricht dann **nur** bei Waisen der `require-levels`-Stufen (SOLLTE/KANN/`unknown` advisory). Fail-closed (leerer Stufen-Name/Keyword, `require-levels`-Eintrag weder Stufe noch `unknown` ⇒ Exit 2), strikt opt-in, default-aus **byte-identisch** (`DC-QA-02`/`DC-QA-03`; ohne `modality` gaten alle Waisen wie bisher). Bereich `MOD` in §3; `DC-FA-MOD-001.a` + Schema-Keys (`trace.requirements.modality.levels`/`require-levels`) in der Spezifikation; Mit-Modifikation `DC-FA-CLI-009` (Modality-Spalte). `--print-config` führt den `modality`-Block. Begründung + Matching-/Unknown-Semantik in begleitender ADR. Anlass: Konsumenten-Analyse grid-gym — die 10 Coverage-Rest-„Waisen" sind 5× KANN (Future) + 4× Nicht-Ziele + 1× DARF NICHT; die slice-zentrische RTM behandelte MUSS und KANN gleich | slice-068 |
| 0.41.0 | 2026-07-11 | Change Request (Auftraggeber): neue Anforderung `DC-FA-COV-001` — **kuratierte Coverage-Quellen der RTM** (`trace.coverage`, opt-in): eine dritte Referenzklasse (Liste benannter `files` + `label` + `ranges` + `sections`/`exclude-sections`), die eine ausgelagerte Traceability-Matrix **range-aware** als Coverage einliest, **ohne** `adrs`/`slices` zu berühren. `<FAM>-AAA..BBB`/`<FAM>-AAA/BBB/CCC` expandieren (breiten-erhaltend, gegen `requirements.id-pattern` validiert); Abschnitts-Whitelist/Blacklist (dieselbe Span-Semantik wie `matrix.exclude-sections` — gegen die „…ohne Design-Artefakt"-Falle). Mit-Modifikationen: `DC-FA-CLI-009` (RTM trägt bei aktiver `trace.coverage` eine **Coverage-Spalte** + `coverage`-Feld in `--json`/`--yaml`) und `DC-FA-CLI-011` (**Waise = ¬slice ∧ ¬coverage**; Slice **oder** Coverage deckt ab, bloße ADR-Referenz weiterhin nicht). Fail-closed (fehlende Datei / ungültige Range `AAA>BBB` / Breiten-Mismatch ⇒ Exit 2), strikt opt-in, default-aus **byte-identisch** (`DC-QA-02`/`DC-QA-03`). Bereich `COV` in §3; `DC-FA-COV-001.a` + Schema-Keys (`trace.coverage[].files`/`label`/`ranges`/`sections`/`exclude-sections`) in der Spezifikation; `--print-config` führt den `coverage`-Block. Begründung + Range-/Sektions-Semantik in begleitender ADR. Anlass: Konsumenten-Analyse grid-gym — 171 „Waisen" waren zu ≥122 anderswo (ADR/traceability.md/Wellen) belegt; d-checks slice-zentrische RTM verfehlte die kuratierte Deckungs-Matrix | slice-067 |
| 0.40.0 | 2026-07-11 | Change Request (Auftraggeber): `DC-FA-CLI-009` (Requirements Traceability Matrix) um einen **opt-in `trace`-Config-Block** erweitert — die vier bislang hart an d-checks Konvention gebundenen RTM-Annahmen sind überschreibbar: Anforderungs-Quelldatei + Kennungs-Gestalt (`trace.requirements.source`/`.id-pattern`) sowie je Referenzklasse ADR/Slice das Verzeichnis, die Dateinamen-Gestalt (Capture-Gruppe = Owner-Kennung) und das Owner-Präfix (`trace.adrs.*`/`trace.slices.*`). Jedes Feld optional; abwesend ⇒ d-checks Konventions-Default ⇒ RTM **byte-identisch** (`DC-QA-02`), nichts geschrieben (`DC-QA-03`). Fail-closed: ungültige Regex oder `file-pattern` ohne Capture-Gruppe ⇒ Exit 2. **Kehrt die 0.21.0-Out-of-Scope-Zeile „frei konfigurierbare Quell-Pfade jenseits der adoptierten Harness-Konvention" um** (durch bewusst begrenzte, explizit gesetzte Config ersetzt — kein Ableiten, je eine ADR-/Slice-Klasse, kein leerer ADR-Präfix, kein VCS). Config-Schema `trace.*` + Algorithmus-Konfigurierbarkeit in `DC-FA-CLI-009.a` der Spezifikation; Begründung (Konsumenten-Konventions-Bindung, Design spiegelt `ids.patterns`) in begleitender ADR. `DC-FA-CLI-011` (`--require-complete`) erbt die konfigurierten Quellen unverändert (gleicher RTM-Lauf). Anlass: Konsumenten-Befund grid-gym 2026-07-11 — `make doc-trace` sah nur 6 von 243 Anforderungen (allein die `GG-QA-*`-Familie traf zufällig d-checks `-QA-`-Default-Gestalt), die übrigen 40 Familien und alle `NNN-…md`-Slices blieben unsichtbar | slice-066 |
| 0.39.0 | 2026-07-06 | Schärfung `DC-FA-CLI-006` (Auftraggeber): die `--suggest-config ai-harness[-init]`-Vorlage an die **gelebte** Dogfood-Konvention angeglichen — `spans`/`hostpaths` ins fixe Standard-Modulset aufgenommen (revidiert die 0.26.0-„situativ nicht aktiviert"-Einordnung; d-checks eigene `.d-check.yml` führte sie längst) und ein **repo-bewusster `planning`-Block** ergänzt (aktiv bei vorhandener `roadmap`, sonst auskommentiert; im Voll-Kanon `ai-harness-init` aktiv). `vcs`/`commits` (Commit-Range) werden auf `--print-mk` verwiesen statt ins statische `modules` gelegt; `versions`/`targets` bleiben bewusst vertagt (repo-spezifische `pin-pattern`/`authority`). Neu benannt: **Eignungs-Kriterium K1–K4** (konventions-kanonisch · ableitungsfrei/konventions-feste Config · Baum-Scan-tauglich · hermetisch) + die **geschlossene** Aktiv-Menge im Body; die kanonische Vorlage in der Spezifikation (`DC-FA-CLI-006.a`) deckt jetzt die emittierte Ausgabe **1:1** (Kommentarzeile inkl. `--print-mk`-Verweis + `codepaths`-Block ergänzt) — der Normativitäts-Spalt (Code-Kommentar ≠ Norm) ist geschlossen. Zugleich die situativen-Modul-Enumeration (AKs/Out-of-Scope) um das seit 0.38.0 fehlende `targets` vervollständigt. Betrifft nur die `ai-harness`-Vorlage (nicht den generischen Quellen-Modus, nicht die eigene `.d-check.yml`); `DC-QA-02`/`DC-QA-03` unberührt (nur mehr Ausgabe). Begründung + Eignungs-Kriterium in begleitender ADR. Anlass: Nutzer-Analyse „welche Module nutzt `--suggest-config ai-harness` nicht, und warum nicht" (2026-07-06) | slice-065 |
| 0.38.0 | 2026-07-05 | Neue Anforderung `DC-FA-TGT-001` (Modul `targets`, opt-in): Deklarations-Konsistenz zwischen Doku und Build-Targets — jedes in einer Doku-**Tabellenzeile** als ` `make X` ` behauptete Target muss eine Makefile-Regel sein (sonst `gate-phantom`), und jede Makefile-Regel (minus `targets.exempt-targets`) muss in der Autoritäts-Doku als ` `make X` ` stehen (sonst `gate-undocumented`) — die maschinelle Form des Vertrags „kein halluziniertes / kein undokumentiertes Gate". **Tabellen-Scoping** (`^\|`, nur Tabellenzeilen — Prosa zählt nicht, sonst spuriöse `gate-phantom`). **Hermetisch** wie `DC-FA-PLAN-001`: Filesystem-Port `ReadFile` auf die konfigurierten Dateien, **kein** Makefile-Ausführen, **kein** git, **kein** Netz (`DC-QA-02`/`DC-QA-03` unberührt); Zeilen-Heuristik für Regelnamen (keine Pattern-Rules/variablen Targets, ≡ dem abgelösten Skript). Strikt opt-in, fail-closed (fehlende konfigurierte Datei ⇒ Exit 2; leeres `makefiles`/`doc-tables` ⇒ inert), diagnose-only, default-aus byte-identisch. 17. Regelmodul; Bereichskürzel `TGT` in §3, `targets` in `DC-FA-CLI-002` + Glossar, Algorithmus-Sektion `DC-FA-TGT-001.a` + Schema-Keys (`targets.makefiles`/`doc-tables`/`authority`/`exempt-targets`) in der Spezifikation; die Grund-Codes `gate-phantom`/`gate-undocumented` (§4) landen mit der Modul-Implementierung (AllReasons-↔-§4-Lockstep). Zugleich **`DC-FA-CLI-010`-Erweiterung** (10→11 Targets): `--print-mk` trägt ein `doc-targets`-Target (`--enable targets` + Fokus-`--disable`, hermetisch ohne Range); `--print-config`/`--suggest-config` führen `targets`. Anlass: Auftraggeber 2026-07-05 — `tools/gate-consistency.sh` driftet über die Repo-Familie (a-check ≠ d-check), Klasse-A-Mechanisierung [`MR-007`](../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding) „verteilen statt kopieren": die verteilbare Modul-Form löst den cross-repo-Kern (Doku↔Makefile), d-check dogfoodet das Modul für sein eigenes `make gate-consistency`-Gate. Identitäts-Ausweitung „Doku-Checker → Deklarations-Konsistenz-Checker" (Makefile-Lesen via Filesystem-Port, wie `planning` die Verzeichnis-Struktur) in begleitender ADR | slice-063 |
| 0.37.1 | 2026-07-03 | Review R1 (doc) + R2 (code) zum `DC-FA-TRK-001`-CR, präzisiert: Verzeichnis-Ziele (Index führt nur Dateien) und Symlink-Referenzen (kategorisch [`DC-FA-LINK-002`](#dc-fa-link-002--symlink-ablehnung)-Domäne — sonst false-positive hinter getrackten Verzeichnis-Symlinks) explizit kein Kandidat/Out-of-Scope; Auflösungs-Mechanik ausdrücklich unabhängig von der Aktivierung des Moduls `links`; intent-to-add + linked worktree als Ränder benannt. `DC-FA-CLI-010`-AKs/Out-of-Scope von neun auf **zehn** Targets nachgezogen (`doc-tracked` — Selbstwiderspruch behoben); `DC-FA-CLI-006`-Enumeration der situativen Module um `commits`/`planning`/`tracked` vervollständigt | slice-059 |
| 0.37.0 | 2026-07-03 | Neue Anforderung `DC-FA-TRK-001` (Modul `tracked`, opt-in): Getrackt-Status auflösbarer Referenz-Ziele — jedes auflösbare, **existierende** Link-/Bild-Ziel muss im **git-Index getrackt** sein, sonst `target-untracked` (die Referenz wäre auf jedem frischen Klon `target-missing` — Umgebungs-Drift zwischen Arbeitsbäumen, am Entstehungsort gefangen statt erst in der CI des nächsten Checkouts). Wahrheit ist der **Index**, nicht die `.gitignore`-Syntax (kein zweiter Regel-Interpreter; frisch gestagte Dateien gelten als getrackt — WIP-tauglich). Liest `.git` über den **VCS-Port** (dritte Nutzung: `vcs` Range-Diff, `commits` Messages, `tracked` **Index** — ohne Range), reine-Go/ohne Netz, lokal/lesend/deterministisch (`DC-QA-02`/`DC-QA-03` wie bei `DC-FA-VCS-001` formuliert). Kein Doppelbefund (nur existierende Ziele; `target-missing` bleibt `links`, Prinzip von `DC-FA-PIN-001`), Ventil `tracked.exempt-targets`, strikt opt-in/fail-closed (aktiv ohne lesbares `.git` ⇒ Exit 2)/diagnose-only/default-aus byte-identisch. 16. Regelmodul; Bereichskürzel `TRK` in §3, `tracked` in `DC-FA-CLI-002` + Glossar, Algorithmus-Sektion `DC-FA-TRK-001.a` + Grund-Code `target-untracked` + Schema-Key (`tracked.exempt-targets`) in der Spezifikation. Zugleich **`DC-FA-CLI-010`-Erweiterung** (9→10 Targets): `--print-mk` trägt ein `doc-tracked`-Target. Anlass: Auftraggeber-Frage 2026-07-03 („Was passiert, wenn ein Dokument ein gitignoriertes Dokument referenziert?") + Demo-Beleg (Erzeuger-Checkout grün, frischer Klon `target-missing`); heutige Ventile (CI-Netz, gitignore+`scan.ignore`-Doppel nach [`MR-017`](../harness/conventions.md#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen)-Muster, Vendoring nach [`MR-019`](../harness/conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)-Muster) fangen die Falle erst spät oder nur je Einzelfall | slice-059 |
| 0.36.0 | 2026-07-01 | Neue Anforderung `DC-FA-PLAN-001` (Modul `planning`, opt-in): maschinelle Durchsetzung der Planning-Lifecycle-Invariante — die Roadmap trägt den Ruhe-Marker (`planning.marker`, „Keine aktive Welle") in ihrem `planning.heading`-Abschnitt (`## Aktuelle Welle`) genau dann, wenn kein `slice-*` (`planning.slice-glob`) im Verzeichnis liegt (`hasActive == hasSlices`), sonst `planning-drift`. **Hermetisch** (nur Roadmap-Datei + Verzeichnis-Listing, **kein** git, **kein** Netz, read-only) — normales Modul wie `codepaths`, `DC-QA-02`/`DC-QA-03` unberührt; fail-closed bei fehlender kanonischer Überschrift bzw. Roadmap-Datei (Heading-Guard), strikt opt-in, diagnose-only, default-aus byte-identisch. 15. Regelmodul; Bereichskürzel `PLAN` in §3, `planning` als Modul in `DC-FA-CLI-002` + Glossar, Algorithmus-Sektion `DC-FA-PLAN-001.a` + Grund-Code `planning-drift` + Schema-Keys (`planning.roadmap`/`marker`/`heading`/`slice-glob`) in der Spezifikation. Anlass: Auftraggeber — `tools/planning-consistency.sh` mechanisieren (letzter Kandidat des `tools/*.sh`-Audits, [`MR-007`](../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Klasse „verteilen statt kopieren"); bewusst nachrangig (die „Keine aktive Welle"-Konvention ist harness-layout-spezifisch, kleinerer Verteilungswert) — im Gegensatz zu `vcs`/`commits` **ohne** git/VCS-Port, rein hermetisch. `make planning-check` dogfoodet das Modul. Zugleich **`DC-FA-CLI-010`-Erweiterung** (8→9 Targets): `--print-mk` trägt ein `doc-planning`-Target; `--print-config`/`--suggest-config` führen `planning` | slice-057 |
| 0.35.0 | 2026-07-01 | Neue Anforderung `DC-FA-COMMITS-001` (Modul `commits`, opt-in): maschinelle Durchsetzung der Traceability-Regel — jede geprüfte Commit-Message muss eine Kennung nach `commits.id-patterns` (`ADR-`/`MR-`/`DC-`/`slice-`) auf einer Inhalts-Zeile tragen, sonst `commit-untraceable`. Liest die Commit-**Messages** über **denselben VCS-Port** wie `DC-FA-VCS-001` (reine-Go-git, **ohne git-Binary** → distroless bleibt, **ohne Netz**), erweitert um Message-Lesen; zwei Quellen für dieselbe Prüfung: `--range <base>..<head>` (CI/Push, Nicht-Merge-Commits) und `--commit-msg <datei\|->` (commit-msg-Hook, einzelne Pending-Message). Uniforme `#`-/scissors-Bereinigung (Kennung auf Inhalts-Zeile), Betreff-Ausnahme `commits.exempt-pattern` (Selbstkonfig `^(Merge \|Revert )`). Strikt opt-in (nie Default, wie `external`/`vcs`), fail-closed ohne `.git`/Range/Message, diagnose-only; erweiterter Eingabe-Scope, aber lokal/lesend/deterministisch — `DC-QA-02`/`DC-QA-03` unberührt. 14. Regelmodul; Bereichskürzel `COMMITS` in §3, `commits` als Modul in `DC-FA-CLI-002` + Glossar, Algorithmus-Sektion `DC-FA-COMMITS-001.a` + Grund-Code `commit-untraceable` + Schema-Keys (`commits.id-patterns`/`exempt-pattern`) in der Spezifikation. Anlass: Auftraggeber — `tools/trace-check.sh` vollständig mechanisieren (Copy-Drift über die Repo-Familie, [`MR-007`](../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Klasse „verteilen statt kopieren"): die verteilbare Modul-Form löst es im Image; d-check dogfoodet das Modul für sein eigenes `make trace-check`-Gate. Das „nächste `adr-check`" — dieselbe VCS-Port-Präzedenz wie `DC-FA-VCS-001`, hier auf Commit-Messages statt Datei-Inhalt. Zugleich **`DC-FA-CLI-010`-Erweiterung** (7→8 Targets): `--print-mk` trägt ein `doc-commits`-Target (`--enable commits` + aus `ValidModules` abgeleitete Fokus-`--disable`-Liste, `--range`) — die **verteilte** Commit-Traceability für Konsumenten ohne Skript-Kopie („verteilen statt kopieren", parallel zu `doc-immutable`); `--print-config`/`--suggest-config` führen `commits` in der Verfügbar-/opt-in-Liste | slice-056 |
| 0.34.0 | 2026-06-29 | `DC-FA-CODE-001` (Modul `codepaths`) um `ignore-refs` erweitert: eine Glob-Liste nimmt **bestimmte Ziel-Pfade** von der Existenz-Prüfung aus — **referenz-weit** (egal Datei/Zeile), als Register bewusst entfernter/historischer Artefakte (Tombstones); dritte Ventil-Achse neben dem zeilenweisen `d-check:ignore` und dem datei-weiten `exempt-paths`. Bewusster Akt mit Gate (ohne Eintrag bleibt `codepath-missing`), Default-leer byte-identisch. Schema-Key `codepaths.ignore-refs` + §DC-FA-CODE-001.a-Schritt in der Spezifikation. Anlass: die Frozen-Doc-Refactoring-Falle — immutable ADRs zitieren Code-Pfade, die bei Refactoring/Löschung dangeln, aber nicht editierbar sind; gelöst und zugleich das in slice-053 nur „pfad-stabil behaltene" `tools/adr-immutable-check.sh` entfernt; begleitende ADR (Supersedes die „Skript-behalten"-Teilentscheidung der VCS-ADR) | slice-054 |
| 0.33.0 | 2026-06-29 | Neue Anforderung `DC-FA-VCS-001` (Modul `vcs`, opt-in): git-historienbasierte Immutabilität des Core über eine Commit-Range — `core(BASE)` ≟ `core(HEAD)`, geliefert über `--range <base>..<head>` (CI/Push) oder `--staged` (pre-commit). Liest das read-only `.git` über einen eigenen **VCS-Port** (reine-Go-git, **ohne git-Binary** → distroless bleibt, **ohne Netz**); erweiterter Eingabe-Scope (git + Range statt nur Markdown-Baum), aber lokal/lesend/deterministisch — `DC-QA-02`/`DC-QA-03` unberührt. Geprüft wird jede in der Range geänderte, der Klasse (`vcs.paths`) entsprechende Datei mit immutabler BASE (`vcs.immutable-when`): Körper-Drift, unzulässiger Status-Übergang (`vcs.head-allow`) oder Löschung/Umbenennung → `core-drift-vcs`. Strikt opt-in (nie Default, wie `external`), fail-closed bei fehlendem `.git`/Range, diagnose-only. **Die volle git-Garantie, die `DC-FA-IMM-001` (hermetischer Pin) bewusst der späteren VCS-Stufe überließ** — beide koexistieren als Defense-in-Depth. 13. Regelmodul; Bereichskürzel `VCS` in §3, `vcs` als Modul in `DC-FA-CLI-002` + Glossar, Algorithmus-Sektion `DC-FA-VCS-001.a` + Grund-Code `core-drift-vcs` + Schema-Keys (`vcs.paths`/`immutable-when`/`exclude-sections`/`status-line`/`head-allow`) in der Spezifikation. Anlass: Auftraggeber — `adr-immutable-check.sh` vollständig mechanisieren (Copy-Drift über die Repo-Familie, [`MR-007`](../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Klasse): die verteilbare git-Form löst es im Image; d-check dogfooded das Modul für die eigenen Accepted-ADRs. Zugleich **DC-FA-CLI-010-Erweiterung** (6→7 Targets): `--print-mk` trägt ein `doc-immutable`-Target (`--enable vcs` + aus dem Modulsatz abgeleitete Fokus-`--disable`-Liste, RANGE/STAGED) — die **verteilte** git-Garantie für Konsumenten ohne Skript-Kopie („verteilen statt kopieren"); `model.ValidModules` dafür exportiert. Außerdem **Config-Surface-Bereinigung**: `--print-config` (`DC-FA-CLI-005`) führt jetzt **alle** Module (die zuvor fehlenden `pins`/`immutable` + `vcs` ergänzt) und `--suggest-config ai-harness` (`DC-FA-CLI-006`) nennt die situativen opt-in-Module vollständig (`versions`/`pins`/`immutable`/`vcs` nachgezogen) — die „Verfügbar"-Liste war eine Harness-Ehrlichkeits-Lücke über drei Slices (049/052/053) | slice-053 |
| 0.32.0 | 2026-06-28 | Neue Anforderung `DC-FA-IMM-001` (Modul `immutable`, opt-in): Immutabilitäts-Pin gegen Core-Drift — eine Datei mit Inline-Marker `<!-- immutable: sha256:… -->` wird gegen den whitespace-normalisierten **Core** (Datei ohne die Marker-Zeile und ohne die per `immutable.exclude-sections` benannten Abschnitte) gehasht; Abweichung → `core-drift`. Hermetische, im read-only-Arbeitsbaum entscheidbare Immutabilitäts-Prüfung (kein git; die git-historienbasierte Diff-Form bleibt Out-of-Scope bzw. spätere opt-in-Stufe, in begleitender ADR festgehalten). Opt-in pro Datei, default-off byte-identisch (`DC-QA-02`), diagnose-only. 12. Regelmodul; Bereichskürzel `IMM` in §3, `immutable` als Modul in `DC-FA-CLI-002` + Glossar, Algorithmus-Sektion `DC-FA-IMM-001.a` + Grund-Code `core-drift` + Schema-Key (`immutable.exclude-sections`) in der Spezifikation. Anlass: Auftraggeber — das ADR-Immutable-Gate war nur ein Skript (Copy-Drift über die Repo-Familie, [`MR-007`](../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Klasse „verteilen statt kopieren"); die verteilbare Pin-Form löst das hermetisch, die volle git-Garantie bleibt einem späteren VCS-Adapter vorbehalten | slice-052 |
| 0.31.0 | 2026-06-28 | Neue Anforderung `DC-FA-MTX-003` (Modul `matrix`): **Token-basierte** Referenz-Richtung + Provenance-Marker + Grandfathering. Eine Klasse kann ein `token`-Regex tragen → `matrix` fängt verbotene Referenzen auch als bare ID-Token im Körper (nicht nur als Link), `matrix-forbidden` in Token-Form. Provenance-Marker `<!-- d-check:status-provenance -->` auf der Zeile nimmt eine verbotene Token-Referenz aus (deklarierte Provenance/Verifikations-Zeiger) — `matrix`' **erster** Zeilen-Marker, kehrt die „nur strukturelle Ausnahmen"-Haltung von `DC-FA-MTX-001` bewusst um (benannt/semantisch, nicht generisches Muting); Ehrlichkeit bleibt Reviewer. Neues `matrix.exempt-paths` grandfathered immutable `Accepted`-ADRs (Regelwerk §Referenz-Richtung). Default-aus byte-identisch (`DC-QA-02`). §DC-FA-MTX-001.a-Schritt + Schema-Keys (`matrix.classes[].token`, `matrix.exempt-paths`) in der Spezifikation. Anlass: Auftraggeber — d-check mechanisiert die Referenz-Richtung, die das adoptierte Regelwerk bewusst dem Reviewer überließ (wie schon [`MR-006`](../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)/`matrix`); der Marker macht die Provenance-vs-Entscheidungsgrundlage-Unterscheidung grep-bar | slice-051 |
| 0.30.0 | 2026-06-28 | Neue Anforderung `DC-FA-MTX-002` (Modul `matrix`): Verweisrichtung **innerhalb** einer geordneten Dokumentklasse — eine Klasse kann `order` (Pfad-Globs, autoritativste Schicht zuerst; Rang = erster Treffer) plus `direction: no-downward` tragen; ein klasseninterner Abwärtsverweis (höherrangig → niederrangig, auch über mehrere Stufen) ⇒ neuer Grund-Code `matrix-downward`. Additiv zu den Klassen-Paar-Regeln (`DC-FA-MTX-001`); generalisiert die Spec-Straten-Schichtung auf viele Dateien je Schicht (Globs statt Einzeldateien-Listing). Fail-closed: `order`/`direction` nur zusammen, unbekannter `direction`-Wert ⇒ Config-Fehler; rangfreie Mitglieder nehmen nicht teil; Default beide leer ⇒ byte-identisch (`DC-QA-02`). Algorithmus-Schritt in `DC-FA-MTX-001.a`, Grund-Code `matrix-downward` + Schema-Keys (`matrix.classes[].order`/`.direction`) in der Spezifikation. Anlass: Auftraggeber — die alternative Einzelklassen-Aufzählung war als Richtung nicht erkennbar und verschattete (First-Match) tote Regeln; zugleich Konsumenten-Bedarf d-migrate (23 Spec-Dateien ⇒ Glob-Schichten statt 23-Zeilen-Listing) | slice-050 |
| 0.29.0 | 2026-06-24 | Neue Anforderung `DC-FA-PIN-001` (Modul `pins`, opt-in): Content-Pin gegen inhaltlichen Drift — ein Link mit Inline-Marker `<!-- dpin: sha256:… -->` (bindet an den unmittelbar vorausgehenden Link derselben Zeile, sonst inert) wird gegen den whitespace-normalisierten **rohen** Ziel-Span (ganze Datei oder Heading-Section, inkl. Fenced-Code) gehasst; Mismatch → `link-stale`. Nur auflösbare Links (struktureller Befund bleibt `DC-FA-LINK-001`/`DC-FA-ANCH-001`, kein Doppelbefund, auch pins-only); opt-in pro Link, default-off byte-identisch (`DC-QA-02`), diagnose-only (`--bless` spätere CR). Bereichskürzel `PIN` in §3, `pins` als Modul in `DC-FA-CLI-002` + Glossar, Algorithmus-Sektion `DC-FA-PIN-001.a` + Grund-Code `link-stale` in der Spezifikation. Anlass: Auftraggeber-Idee 2 (stale citation) + Spike (Drift real, Rauschen ~0 bei Normalisierung) | slice-049 |
| 0.28.0 | 2026-06-24 | Neue Anforderung `DC-FA-VER-001` (Modul `versions`, opt-in): Versions-Pin-Konsistenz — alle Pins (`versions.pin-pattern`) müssen die aktuelle Version aus `versions.current-from` (Default `version.md#aktuell`) tragen, sonst `version-stale`; liest dafür Pins **auch in Fenced-Code** (gescopte Fence-Ausnahme), Ventile `exempt-paths`/`d-check:ignore` für historische Pins; opt-in/default-off (ohne Block byte-identisch, `DC-QA-02`), diagnose-only (Auto-Bump-`--repair` als Folge-CR an `DC-FA-CLI-008`). Bereichskürzel `VER` in §3, `versions` als gültiges Modul in `DC-FA-CLI-002` + Glossar, Algorithmus-Sektion `DC-FA-VER-001.a`, Grund-Code `version-stale` und Config-Schema (`versions.pin-pattern`/`current-from`/`exempt-paths`) in der Spezifikation ergänzt. Anlass: Auftraggeber-Idee „nicht vergessen, die Version zu bumpen" + Spike (Meta-Gate-Skript wegen Copy-Drift über die Repo-Familie verworfen) | slice-048 |
| 0.27.0 | 2026-06-23 | Change Request (Auftraggeber): `DC-FA-CLI-010` (`--print-mk`-Fragment) um drei Targets + eine Variable erweitert — `doc-doctor` (`--doctor`), `doc-repair` (`--repair`, Recipe-Echo unterdrückt für `git apply`-reine stdout) und `doc-help` (namespaced, listet die `doc-*`-Targets via `##`-Annotationen; **kein** `help` wegen Konsumenten-Kollision) sowie `DCHECK_DIGEST` (Digest-Override per `ifeq`, sticht den Tag). Alle Targets `##`-annotiert (greift das `help` des Konsumenten auf). Read-only/deterministisch unverändert. Anlass: Auftraggeber-Wunsch nach `doc-doctor`/`doc-repair`/Self-Doc/Digest-Komfort | slice-047 |
| 0.26.0 | 2026-06-23 | Schärfung `DC-FA-CLI-006` (Auftraggeber): das opt-in-Modul `diagrams` zur Out-of-Scope-Liste der **nicht** auto-aktivierten situativen Module ergänzt (`external`/`spans`/`hostpaths`/`diagrams`); die `--suggest-config ai-harness[-init]`-Ausgabe nennt diese situativen Module stattdessen in einem **Kommentar mit Verweis auf `--print-config`** (Auffindbarkeit ohne Aktivieren eines inerten Moduls — `diagrams` braucht repo-spezifische `patterns`/`defined-in`, lässt sich nicht ableiten). Read-only/advisory unverändert. Anlass: Nutzer-Frage nach slice-045 (wird `diagrams` in `--suggest-config ai-harness` berücksichtigt?) | slice-046 |
| 0.25.0 | 2026-06-23 | Change Request (Auftraggeber): neue Anforderung `DC-FA-DIAG-001` — opt-in Modul `diagrams` öffnet gezielt benannte Diagramm-Fences (Default `mermaid`) und prüft die darin gefundenen Kennungen auf **Existenz** in ihrer `defined-in`-Quelle (Befund `diagram-id-undefined`); reine Token-Extraktion ohne Mermaid-Parser, read-only/netzlos (`DC-QA-03`), deterministisch (`DC-QA-02`), Default aus (byte-identisch). Link-Policy gilt in Fences nicht (keine Markdown-Links möglich) → Existenz statt Linkpflicht. Bereich `DIAG` in der Schema-Konvention deklariert; Modul-Liste in `DC-FA-CLI-002` ergänzt. Anlass: belief-agent-Architektur — `ARC-NN`/`LH-*`-Kennungen in `mermaid`-Diagrammen entgehen heute allen Modulen, weil Fences opak sind | slice-045 |
| 0.24.0 | 2026-06-23 | Change Request (Auftraggeber): neue Anforderung `DC-FA-CLI-011` — opt-in `--trace --require-complete` bindet die Waisen-Markierung an den Exit-Code (≥1 Requirements-Waise ⇒ Exit 1, sonst 0); Default-`--trace` bleibt advisory (Exit 0). Mit-Erweiterung `DC-FA-CLI-010`: das `--print-mk`-Fragment trägt zusätzlich `doc-trace` (advisory RTM) und `doc-complete` (Vollständigkeits-Gate) plus eine `TRACE_FLAGS`-Variable. Read-only (`DC-QA-03`), deterministisch (`DC-QA-02`), kein Dokument geschrieben. Anlass: Konsumenten (a-check-Bootstrap) sollen die Vollständigkeits-Invariante als Makefile-Gate binden, ohne die `completeness-check.sh`-Parsing-Logik zu kopieren | slice-044 |
| 0.23.0 | 2026-06-22 | Change Request (Auftraggeber): `DC-FA-CODE-001` um das Datei-Ventil `exempt-paths` (Glob-Liste, Syntax wie `scan.ignore`) erweitert — nimmt **ganze Dateien** von der `codepaths`-Prüfung aus, datei-weit und unabhängig von `roots`; dasselbe Ventil wie `DC-FA-ID-001` (slice-018/023), komplementär zum zeilenweisen `d-check:ignore`-Marker. Abwärtskompatibel (`DC-QA-02`): ohne gesetztes `exempt-paths` byte-identisch. Anlass: slice-042-Nebenbefund — Review-Reports unter `docs/reviews/` zitieren naturgemäß `Datei:Zeile`/Pfade und lösten `codepath-missing` aus (`ids` exemptet `docs/reviews/**` längst; `codepaths` zieht nach) | slice-043 |
| 0.22.0 | 2026-06-21 | Change Request (Auftraggeber): neue Anforderung `DC-FA-CLI-010` — Option `--print-mk` gibt ein include-bares `d-check.mk` (überschreibbare `DCHECK_IMAGE`-Variable mit version-gepinntem Image + `doc-check`-Target) auf stdout aus; read-only (`DC-QA-03`), deterministisch (`DC-QA-02`), kein Dokument geschrieben. Konsumenten `include`-n statt Recipe/Skript zu kopieren; der Image-Ref ist die ins Binary eingebettete Release-Version (Digest via `DCHECK_IMAGE`-Override — Henne-Ei: das Binary kennt seinen eigenen Digest nicht). Anlass: a-check-Bootstrap 2026-06-20 — `d-check.mk` wurde handgepflegt; `--print-mk` verlagert den Pin nach d-check | slice-038 |
| 0.21.0 | 2026-06-21 | Change Request (Auftraggeber): neue Anforderung `DC-FA-CLI-009` — Option `--trace` gibt eine **Requirements Traceability Matrix** (Anforderung → referenzierende ADRs/Slices, Waisen-Markierung) auf stdout aus (Default Markdown-Tabelle, optional `--trace --json`/`--yaml` über den format-neutralen Reporter); read-only (`DC-QA-03`), deterministisch (`DC-QA-02`), kein Dokument erzeugt; **Doku-Domäne** (Lastenheft/ADR/Planning), kein Code/keine Go-Toolchain (arch-check bewusst ausgeklammert). Anlass: RTM als d-check-Modus statt separatem Skript (Nutzer-Entscheid 2026-06-20, Prototyp-Beleg) | slice-036 |
| 0.20.0 | 2026-06-21 | Change Request (Auftraggeber): `DC-FA-CLI-006` um einen **Anforderungs-Präfix-Parameter** erweitert — `--suggest-config ai-harness[-init]` backt nicht mehr fix d-checks `DC-` ein: `--id-prefix <PREFIX>` setzt es explizit, der Modus `ai-harness` leitet es aus dem Lastenheft ab (eindeutige FA-/QA-Kennung; mehrere verschiedene ⇒ Nutzungsfehler), ohne Angabe/Ableitung erscheint ein markierter Platzhalter `<PREFIX>` + TODO statt eines stillen `DC-`. **Breaking** ggü. 0.18.1 (`ai-harness-init` ohne Präfix lieferte zuvor `DC-`); Begründung in eigener ADR. Anlass: a-check-Bootstrap 2026-06-20 — das Init-Template emittierte d-checks Eigen-Präfix in ein Fremd-Repo | slice-037 |
| 0.19.0 | 2026-06-19 | Change Request (Auftraggeber): `DC-FA-CLI-004` um das Ausgabeformat **YAML** (`--yaml`) erweitert — gleiche Struktur wie `--json` (`findings`/`summary`/`exitCode`), volle Parität inkl. `--doctor --yaml` (mit `reasonText`/`fixCandidate`); `--json`+`--yaml` und `--repair`+`--yaml` sind Nutzungsfehler (Exit 2). Serialisierung im report-Adapter braucht `gopkg.in/yaml.v3` dort — die Modul-Import-Regeln werden dafür per Folge-ADR erweitert. Anlass: YAML als lesbareres maschinenlesbares Format neben JSON | slice-031 |
| 0.18.1 | 2026-06-19 | Schärfung `DC-FA-CLI-006` (Auftraggeber, vor Release): die ai-harness-Vorlage in **zwei explizite Modi** aufgeteilt (Henne-Ei — nicht aus der Repo-Existenz ableitbar) — `ai-harness-init` (Voll-Kanon, alle Blöcke aktiv; Zielbild für ein leeres/frisches Repo, läuft nach Struktur-Anlage) und `ai-harness` (repo-bewusst, fehlende Blöcke auskommentiert; läuft sofort gegen ein bestehendes Repo). Anlass: der einzelne hybride Modus kommentierte in einem leeren Repo alles aus und taugte nicht als Bootstrap-Vorlage | slice-030 |
| 0.18.0 | 2026-06-19 | Change Request (Auftraggeber): `DC-FA-CLI-006` um die reservierte Quelle `ai-harness` erweitert — `--suggest-config ai-harness` schlägt ein an die adoptierte ai-harness-course-Konvention angelehntes `.d-check.yml` vor (kanonische `ids`-Muster, `matrix`-Klassen samt Referenzrichtung, Standard-Modulset, Scan-Scope). Hybrid: strukturelle Konventionen immer, konkrete Pfade repo-bewusst (nur existierende roots; fehlende Artefakte auskommentiert mit Hinweis). Read-only/advisory, deterministisch (`DC-QA-02`); liefert — als Ausnahme zur reinen Quellen-Ableitung — `matrix`-Regeln/`link-policy`/`exempt-paths` aus der bekannten Konvention. Anlass: Adoptions-Start für Harness-Repos ohne manuelles Quellen-Auflisten | slice-030 |
| 0.17.1 | 2026-06-19 | Change Request (Auftraggeber): Out-of-Scope von `DC-FA-CLI-008` präzisiert — reparierbar sind nur `id-unlinked` (konservativ) und `target-missing` als eindeutiger Datei-Move (breit; Markdown, im Scan-Bestand eindeutiger Basisname); alle übrigen Befundarten bleiben Befund unter `DC-FA-CLI-007`. VCS-/git-historienbasierte Move-/Rename-Erkennung explizit ausgeschlossen (erweiterte die Eingabe über den gescannten read-only-Baum hinaus, `DC-QA-02`; wäre ein eigenes opt-in-Modul analog `external`). Klarstellung des bestehenden Vertrags — keine Verhaltens-/Akzeptanzkriterien-Änderung; Begründung in eigener ADR festgehalten | — |
| 0.17.0 | 2026-06-18 | Change Request (Auftraggeber): `--doctor` wird mit `--json` kombinierbar — Schärfung `DC-FA-CLI-007` (maschinenlesbare Diagnose: `findings` je Eintrag um `reasonText` und `fixCandidate` erweitert, Gruppierung über das `file`-Feld) und `DC-FA-CLI-004` (Kombinierbarkeit `--doctor`+`--json` definiert statt verboten; `--repair`+`--json` und `--doctor`+`--repair` bleiben Nutzungsfehler). Dieselbe Grund-Klartext-/Fix-Kandidaten-Ableitung, drittes Rendering neben Prosa und Patch. Anlass: Auftraggeber-Wunsch nach maschinenlesbarer Diagnose | slice-029 |
| 0.16.0 | 2026-06-18 | Change Request (Auftraggeber): neue Anforderung `DC-FA-CLI-008` (Option `--repair`: unified diff auf stdout, `git apply`-kompatibel, read-only; konservative Stufe als Default mit eindeutig ableitbaren Fixes — v1 v. a. `id-unlinked` → Definitions-Link —, breite Best-Guess-Stufe opt-in, Kennzeichnung review-pflichtiger Hunks auf stderr, damit der Patch `git apply`-rein bleibt). Baut auf den Ausgabe-Modi aus 0.15.0 auf; deterministisch (`DC-QA-02`), read-only-Kernvertrag (`DC-QA-03`); In-place-Schreiben bleibt Out-of-Scope (wäre eigene Anforderung + ADR) | slice-026 |
| 0.15.0 | 2026-06-18 | Change Request (Auftraggeber): neue Anforderung `DC-FA-CLI-007` (Option `--doctor`: erklärende, nach Datei/Regel gruppierte Diagnose mit Fix-Kandidaten; read-only, stdout-only, deterministisch `DC-QA-02`/`DC-QA-03`). Mit-Schärfung `DC-FA-CLI-003` (die Ausgabe-Modi folgen den Codes 0/1/2, Ausgabe auf stdout unabhängig vom Code) und `DC-FA-CLI-004` (die Modi `--doctor`/`--repair` ersetzen das Default-stdout-Format; untereinander und mit `--json` nicht kombinierbar → Nutzungsfehler exit 2; JSON-Varianten out of scope). Anlass: Machbarkeitsfrage des Auftraggebers (2026-06-18) — beratende Diagnose/Reparatur ohne Bruch der Seiteneffektfreiheit | slice-025 |
| 0.14.0 | 2026-06-17 | Change Request (Auftraggeber): `DC-FA-MTX-001` — opt-in `allow-supersede-lineage` (+ `supersede-fields`) nimmt die Supersede-Lineage-Kante X → Y von der Status-Prüfung aus: ein ablösendes Dokument darf auf das von ihm abgelöste (inaktive) Dokument verweisen, ohne `matrix-inactive` zu erzeugen, sofern X über ein deklariertes Feld die Ablösung von Y benennt (Match über Linktext bzw. Zielpfad der Referenz). Carve-out beschränkt auf die deklarierte Lineage-Kante; `matrix-forbidden` (Klassen-Regeln) unberührt. Abwärtskompatibel (`DC-QA-02`): Default aus ⇒ Befundsatz byte-identisch. Marker-Politik (B2 der CR): `matrix` bleibt bewusst ohne Zeilen-Opt-out-Marker `d-check:ignore` (legitime Lineage strukturell ausgenommen, nicht stummgeschaltet). Anlass: reproduzierter Fremd-Repo-Befund (`grid-gym`, v0.10.0) — normative ADR→ADR-Lineage (`Aenderungstyp: Supersedes …`) als `matrix-inactive` gemeldet, brach `make docs-check` | slice-024 |
| 0.13.0 | 2026-06-16 | Change Request (Auftraggeber): `DC-FA-ID-001` — Geltungsbereich der beiden Ventile `exempt-paths` und `d-check:ignore` auf **nackte Fließtext-Vorkommen** erweitert. Bisher griffen sie nur auf die `always`-Inline-Code-Vorkommen; eine nackte Kennung in einer `exempt-paths`-Datei (bzw. auf einer `d-check:ignore`-Zeile) wurde weiterhin als `id-unlinked` gemeldet. Neu: beide Ventile sind ein Ganzdatei- bzw. Ganzzeilen-Carve-out für alle Vorkommen des Musters, unabhängig von der `link-policy`. Abwärtskompatibel (`DC-QA-02`): Configs ohne gesetzte Ventile bleiben byte-identisch; Wirkung nur in Richtung *weniger* Befunde in explizit ausgenommenen Dateien/Zeilen. Anlass: reproduzierter Fremd-Repo-Befund (`ai-harness-init`, v0.8.0/v0.9.0) — nacktes `MR-004` in `docs/reviews/_t.md` trotz `exempt-paths: ["docs/reviews/**"]` gemeldet | slice-023 |
| 0.12.0 | 2026-06-15 | Change Request (Auftraggeber): `DC-FA-ANCH-001` erweitert — Inline-HTML-Anker innerhalb von Markdown zählen zur gültigen Anker-Menge (`id` an beliebigem Element, `name` an `<a>`; GitHub-Parität, wörtlicher Vergleich). Hebt die bisherige HTML-Anker-Out-of-Scope-Zeile auf; Abgrenzung zu §5 (HTML als Dateiformat bleibt out-of-scope) explizit. Gilt mittelbar auch für `codepaths` (geteiltes Anker-Verfahren). Anlass: Falsch-Befunde `anchor-missing` auf manuell gesetzte HTML-Anker in der Doku | slice-022 |
| 0.11.0 | 2026-06-13 | Schärfung `DC-FA-CLI-001` (Auftraggeber): neues Akzeptanzkriterium für `--help` — die Hilfe nennt die Synopsis `d-check [optionen] [pfad]`, beschreibt das Pfad-Argument (Scan-Wurzel, Default cwd) und verweist für das Config-Format auf `--print-config` (kein Format-Duplikat). Anlass: die nackte `flag`-Default-Usage verschwieg das Pfad-Argument | slice-021 |
| 0.10.0 | 2026-06-13 | Change Request (Auftraggeber): neue Anforderung `DC-FA-CLI-006` — Option `--suggest-config` leitet die `ids`-Config aus benannten Autoritäts-Dokumenten ab (definierte Kennungen → Muster + target, Round-Trip-Garantie; opt-in-Module nach Signal; Ausgabe-only, read-only). Dazu Schärfung von `DC-FA-ID-001`: das Muster-Ableitungs-Out-of-Scope gilt für die **Prüfung**; ein advisory Scaffold-Modus darf aus benannten Autoritäts-Quellen ableiten. Anlass: Adoptions-Reibung — neue Repos brauchen einen treffsicheren Config-Start statt eines rein statischen Gerüsts | slice-020 |
| 0.9.0 | 2026-06-13 | Change Request (Auftraggeber): neue Anforderung `DC-FA-CLI-005` — Option `--print-config` gibt ein statisches, kommentiertes `.d-check.yml`-Startgerüst auf stdout aus (kein Repo-Zugriff, kein Schreiben — read-only-Vertrag bleibt; Umleiten via `> .d-check.yml` macht der Aufrufer). Anlass: Adoptions-Reibung in neuen Repos ohne Config; macht zugleich die verfügbaren Optionen sichtbar. Ableitung aus Repo-Inhalt bewusst nicht Teil (späterer eigener Modus) | slice-019 |
| 0.8.0 | 2026-06-13 | Change Request (Auftraggeber): `DC-FA-ID-001` um konfigurierbare `link-policy: prose\|always` (je Muster) erweitert — `always` macht auch Inline-Code-Vorkommen linkpflichtig, Default `prose` (opt-in, abwärtskompatibel). Zwei Ventile: `exempt-paths` (Glob-Liste je Muster) und der Zeilen-Marker `d-check:ignore`, dessen Geltungsbereich von `codepaths`-only auf `ids` erweitert wird (illustrative Beispiel-IDs). Anlass: der `ids`-Sensor maß bislang nicht das Ziel „gut verlinkt" — ein Code-Span konnte stillschweigend einen fehlenden Link verbergen (Ausgangsbefund `DC-QA-03` in slice-017). Kalibrierung über die drei `ids`-Repos (d-check 155, u-boot 9, b-trace 2) bestimmte die Ventil-Form | slice-018 |
| 0.1.0 | 2026-06-10 | Initiale Fassung (Konsolidierung von 12 Quell-Tools, Modul-Schnitt, Docker-Distribution) | — |
| 0.2.0 | 2026-06-10 | Review-Runde R1: Modul-Schnitt `links`/`anchors` präzisiert (Fragment-Zuständigkeit, fehlende Zieldatei), Slug-Duplikat-Reihenfolge, Symlink-Vorrang, RFC-3986-Dekodierung vor Escape-Prüfung, Redirect-Regel `external`, Muster-Präzedenz `ids`, Status-Default `matrix`, Scan-Wurzel- und Config-Vollvalidierung, Out-of-Scope Reference-Style-Links, Image-Default-Befehl | — |
| 0.2.1 | 2026-06-10 | Redaktionell: Beispiel-Kennungen in DC-FA-ID-001/DC-FA-MTX-001/Glossar auf fiktive Nummern (`ADR-0042`, `ADR-0099`) umgestellt — Kollision mit real entstandenen/zukünftigen eigenen ADRs vermeiden; keine inhaltliche Änderung | — | <!-- d-check:ignore (Beispiel-ID, fiktiv) -->
| 0.2.2 | 2026-06-10 | Redaktionell: absolute Workspace-Pfade entfernt („Schwester-Repositories des Entwicklungs-Workspace" statt konkreter Pfade); keine inhaltliche Änderung | — |
| 0.2.3 | 2026-06-11 | Redaktionell: „(folgt)" in der `DC-QA-01`-Messmethode entfernt — die Benchmark-Definition existiert in der Spezifikation; keine inhaltliche Änderung | — |
| 0.7.2 | 2026-06-12 | Review zum Change Request `DC-FA-HOST-001` (Kalibrierungs-Befund slice-016): tmp aus der Default-Präfixliste gestrichen — der Kalibrierungslauf zeigte überwiegend legitime Laufzeit-Doku (Log-/Output-Pfade) statt Maschinen-Layout-Leaks; wer tmp prüfen will, konfiguriert es | slice-016 |
| 0.7.1 | 2026-06-12 | Review zum Change Request `DC-FA-SPAN-001` (Kalibrierungs-Befund slice-015): Bildreferenzen als Linktext (`[![…](…)](…)` — Badge-Muster, z. B. Shields in vendorten Paket-READMEs) sind legales Markdown und kein `span-nested-link`-Treffer | slice-015 |
| 0.7.0 | 2026-06-12 | Change Request (Erst-Bedarfsträger grid-gym): neue Anforderung `DC-FA-CONF-002` — optionaler modul-lokaler Scan-Scope `<modul>.scope` (`roots`/`ignore`), ersetzt für genau dieses Modul den globalen Scope; additiv und abwärtskompatibel, Constraints spiegeln `scan.*`. Anlass (gemessen, v0.2.0): `ids`-Aktivierung in grid-gym liefert global 2776 Befunde (Masse: Retro-Verlinkung historischer done-Planning-Docs = Umschreiben des Audit-Trails) vs. 312 echte, fixbare Befunde im kuratierten Scope `spec/` + `docs/user/` — der Linkpflicht-Nutzen ist heute nur um den Preis des globalen Sweeps oder des Verzichts auf breite `links`/`anchors`-Abdeckung zu haben | slice-017 |
| 0.6.0 | 2026-06-12 | Change Request (Auftraggeber): neue Anforderung `DC-FA-HOST-001` — Modul `hostpaths` meldet host-lokale absolute Pfade in Prosa und Inline-Code (opt-in; Unix-Präfixliste konfigurierbar, Windows-/UNC-Muster fest; Fences ausgenommen, kein Opt-out-Marker). Anlass: der bess-ems-Rest-Sensor generalisiert (dort als eigenes Tool gebaut) plus die 8 Host-Pfad-Links aus dem d-migrate-Vergleichslauf und die eigene 0.2.2-Hygiene-Korrektur — dieselbe Leak-Klasse dreimal unabhängig. Bereich `HOST` deklariert; Modul-Listen in `DC-FA-CLI-002` und Glossar ergänzt | slice-016 |
| 0.5.0 | 2026-06-12 | Change Request (Auftraggeber): neue Anforderung `DC-FA-SPAN-001` — Modul `spans` meldet ungeschlossene Code-Spans und verschachtelte Link-Artefakte (opt-in, konservative Erkennung: alleinstehende Backticks bleiben befundfrei; kein Opt-out-Marker). Anlass: die `DC-QA-04`-Vergleichsläufe aus slice-012/slice-014 — ~100 u-boot-Artefakte sowie Fälle in grid-gym und bess-ems wurden nur indirekt als `id-unlinked`-Folgefehler sichtbar; ein direkter Sensor meldet die Ursache statt der Symptome. Bereich `SPAN` in der Schema-Konvention deklariert; Modul-Listen in `DC-FA-CLI-002` und Glossar ergänzt | slice-015 |
| 0.3.0 | 2026-06-11 | Change Request (Auftraggeber): neue Anforderung `DC-FA-CODE-001` — Modul `codepaths` prüft explizite Pfade in Inline-Code (opt-in, konservative Erkennung, Zeilen-Opt-out `d-check:ignore` nur für dieses Modul). Anlass: `DC-QA-04`-Vergleichslauf gegen die JS-Familie (`docs-check.js`) zeigte die Prüfklasse als Konsolidierungs-Lücke. Bereich `CODE` in der Schema-Konvention deklariert; Modul-Listen in `DC-FA-CLI-002` und Glossar ergänzt | slice-013 |
| 0.4.0 | 2026-06-12 | Change Request (Auftraggeber): Inventur-Nachtrag zu `DC-QA-04` — dreizehntes Alt-Tool-Vorkommen entdeckt (`check_markdown_links.py` in bess-ems, eigenständige Python-Linie, entstanden 2026-05-24 und damit vor der Inventur vom 2026-06-10 übersehen); Anforderungs-Text, Einleitung, Stakeholder-Tabelle und Glossar von zwölf auf dreizehn fortgeschrieben. Messmethode (drei Familien-Piloten) unverändert | slice-014 |
| 0.3.1 | 2026-06-11 | Review R1 zum Change Request: `DC-FA-CODE-001` präzisiert — Wert-Normalisierung (Anführungszeichen, schließende Satzzeichen) und Anker-Prüfung bei Markdown-Zielen (`DC-QA-04`-Parität zur JS-Familie); `DC-FA-LINK-001`-Inline-Code-Aussage auf Modul-Bezug eingegrenzt | slice-013 |
