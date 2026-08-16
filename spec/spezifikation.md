# Spezifikation — d-check

**Status:** Aktiv. **Letzte Änderung:** 2026-08-09.

**Bezug zum Lastenheft:** Diese Spezifikation präzisiert die in
[`lastenheft.md`](lastenheft.md) formulierten Anforderungen
(`DC-*`-IDs). Bei Konflikt gewinnt das Lastenheft.

---

## 1. Algorithmen und Datenflüsse

### DC-FA-CLI-001.a — Ablauf eines Prüflaufs

**Eingabe:** Scan-Wurzel (Argument oder cwd), CLI-Optionen, optionale
`.d-check.yml`. **Ausgabe:** Befundliste + Exit-Code. **Schritte:**

1. Konfiguration laden und vollständig validieren; jeder Fehler →
   Exit 2, keine Prüfung
   ([`DC-FA-CONF-001`](lastenheft.md#dc-fa-conf-001--konfigurationsdatei)).
2. Effektive Module bestimmen (siehe
   [DC-FA-CLI-002.a](#dc-fa-cli-002a--modul-auflösung)).
3. Markdown-Dateien gemäß
   [`DC-FA-SCAN-001`](lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln)
   ermitteln (deterministische Reihenfolge: Pfade bytewise sortiert).
4. Pro Datei: Vorverarbeitung (Fences, Inline-Code), Extraktion von
   Links/Headings/Kennungen; aktive Module erzeugen Befunde.
5. Befunde sammeln, deduplizieren (identisches Tupel aus Datei, Zeile,
   Regel, Ziel, Grund), sortieren
   ([DC-QA-02.a](#dc-qa-02a--determinismus-und-sortierung)), ausgeben.
6. Exit-Code gemäß
   [`DC-FA-CLI-003`](lastenheft.md#dc-fa-cli-003--exit-codes).

**Fehlermodi:** nicht lesbare Datei oder Scan-Wurzel → Umgebungsfehler,
Exit 2 (kein Teilergebnis als Erfolg). Eine **gänzlich leere**
Scan-Wurzel (keinerlei Einträge) ist ebenfalls ein Umgebungsfehler —
Exit 2 mit Mount-Hinweis
([`DC-FA-DIST-001`](lastenheft.md#dc-fa-dist-001--docker-image)
Negative); eine Wurzel ohne Markdown-Dateien, aber mit Inhalt, liefert
„0 Datei(en) geprüft" und Exit 0
([`DC-FA-CLI-001`](lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel)
Boundary). CLI-Optionen dürfen vor oder nach dem Pfad-Argument stehen
(Container-Aufrufmuster: ENTRYPOINT setzt `/repo`, Optionen werden
angehängt); ein wertnehmendes Flag ohne Wert ist ein Nutzungsfehler
(Exit 2), Nutzungsfehler tragen den Präfix `d-check: error:`, und
`-h`/`--help` zeigt die Nutzung auf stderr und endet mit Exit 0. Die
Usage-Ausgabe führt (in dieser Reihenfolge) eine Kurzbeschreibung, die
Synopsis `d-check [optionen] [pfad]`, eine Zeile zum Pfad-Argument
(Scan-Wurzel, Default: aktuelles Verzeichnis), die Flag-Liste
(`flag.PrintDefaults`) und einen Konfigurations-Hinweis, der auf
[`DC-FA-CLI-005`](lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben)
(`--print-config`) und
[`DC-FA-CLI-006`](lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)
(`--suggest-config`) verweist — das Config-Format wird nicht dupliziert.
Symlinks werden beim Scan weder verfolgt noch als Dateien gewertet — der
Baum-Walk kennt nur Verzeichnisse und reguläre Dateien. Das gilt für **beide**
Formen: ein Verzeichnis-Symlink wird nicht betreten, und eine nur über einen
Datei-Symlink erreichbare Markdown-Datei ist **keine** Kandidatin, in keinem
Modul (Symlink-Ablehnung,
[`DC-FA-LINK-002`](lastenheft.md#dc-fa-link-002--symlink-ablehnung)).

### DC-FA-CLI-002.a — Modul-Auflösung

Effektive Module = (`modules` aus Config, sonst `DEFAULT_MODULES`)
∪ `--enable`-Angaben ∖ `--disable`-Angaben. CLI-Angaben werden nach
der Config angewandt (CLI-Präzedenz). Unbekannte Modulnamen sind
Nutzungsfehler (Exit 2) mit Auflistung der gültigen Namen.

### DC-FA-CLI-005.a — Konfigurations-Gerüst

`--print-config` ist ein Kurzschluss-Modus wie `-h`: Nach dem
Optionen-Parsing — **vor** Scan-Wurzel-Öffnung, Config-Laden und
Prüflauf — wird ein statisches, in das Binary eingebettetes
`.d-check.yml`-Gerüst unverändert auf stdout geschrieben, Exit-Code 0.
Es findet **kein** Zugriff auf das geprüfte Repository statt (weder
Lesen noch Schreiben — die Seiteneffektfreiheit aus
[`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
ist hier trivial, da nicht einmal gelesen wird). Da das Gerüst eine
Konstante ist, ist die Ausgabe deterministisch
([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)). Das Gerüst ist
gültiges YAML, das der eigene Konfigurations-Parser fehlerfrei
dekodiert ([`DC-FA-CONF-001`](lastenheft.md#dc-fa-conf-001--konfigurationsdatei));
es führt die aktiven Default-Einstellungen und die übrigen Module mit
ihren Optionen als kommentierte Beispiele (`ids.link-policy`,
`exempt-paths`, `<modul>.scope` u. a.). Trifft `--print-config` mit
weiteren Optionen oder einem Pfad-Argument zusammen, gewinnt der
Kurzschluss (die übrigen Angaben bleiben wirkungslos).

### DC-FA-CLI-006.a — Konfigurations-Vorschlag

`--suggest-config <quelle>[,<quelle>…]` ist ein **Lese**-Modus, der ein
Gerüst statt Befunde ausgibt. Schreibzugriff entsteht nie
([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
Schritte (deterministisch, alle Mengen bytewise sortiert —
[`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)):

1. **Quellen auflösen:** jede `<quelle>` (Datei oder Verzeichnis,
   relativ zur Repo-Wurzel, kein Repo-Escape — sonst Exit 2). Eine nicht
   existierende Quelle ist ein Nutzungsfehler (Exit 2).
2. **Kennungen extrahieren:** je Quelle alle Markdown-Dateien lesen; aus
   jeder ATX-Heading-Zeile das **führende** Token nehmen, das der
   allgemeinen Kennungs-Gestalt entspricht
   (`[A-Z][A-Z0-9]*(?:-[A-Z0-9]+)*-\d+[A-Za-z]?`). Ergebnis je Quelle:
   die sortierte Menge der dort *definierten* Kennungen (kein
   Fließtext-Mining).
3. **Muster ableiten (je Quelle):** jede Kennung in Präfix + `-\d+` +
   optionalen Buchstaben zerlegen; die distinkten Präfixe ergeben
   `(?:<präfix₁>|<präfix₂>|…)-\d+` (mit `[A-Za-z]?`, falls eine Kennung
   einen Suffix-Buchstaben trug). **Round-Trip-Invariante:** der
   abgeleitete `regex` matcht jede Quell-Kennung (Best-Guess-Verengung
   ist Sache des Menschen — die Quell-Kennungen werden als Kommentar
   mitgegeben). `target` = die Quelle.
4. **Opt-in-Module nach Signal:** die opt-in-Module (`codepaths`,
   `spans`, `hostpaths`) probeweise über das Repo laufen lassen; jene
   mit ≥1 Befund werden als Vorschlag vermerkt (die Default-Module
   `links`/`anchors` immer aktiv).
5. **Scan-Scope** konservativ vorschlagen (Default-Wurzeln).
6. Das zusammengesetzte, kommentierte Gerüst auf stdout schreiben,
   Exit 0; es dekodiert über den eigenen Parser
   ([`DC-FA-CONF-001`](lastenheft.md#dc-fa-conf-001--konfigurationsdatei)).

**Reservierte Quellen `ai-harness` / `ai-harness-init`.** Beide sind statt
eines Pfads **reservierte Schlüsselwörter**: Schritt 1 löst sie nicht als
Quelle auf (ein fehlendes Verzeichnis dieses Namens ist **kein** Fehler),
stattdessen wird ein an die ai-harness-course-Konvention (Baseline
**v1.3.0**) angelehntes Gerüst erzeugt. Kombinierbar mit echten Quellen
(`<token>,<quelle>,…`); deren Muster (Schritte 2–4) werden unter
`ids.patterns` angehängt. Read-only, deterministisch (feste Block- und
Pfad-Reihenfolge, keine Map-Iteration für die Ausgabe). Das fixe Modulset ist
`links`, `anchors`, `ids`, `matrix`, `codepaths`, `spans`, `hostpaths` (kein
Probelauf), dazu ein **repo-bewusster `planning`-Block** (Modul + `roadmap:`-Block
aktiv, wenn `docs/plan/planning/in-progress/roadmap.md` existiert bzw. im
Voll-Kanon, sonst auskommentiert). Aufnahme-Kriterium (K1–K4) und geschlossene
Aktiv-Menge: [`DC-FA-CLI-006`](lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten);
die situativen Range-Module `vcs`/`commits` werden über
[`--print-mk`](lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben) verteilt,
nicht ins statische `modules` gelegt.
`CO-\d{3}` (Carveouts) hat kein festes Definitions-`target` und bleibt in
beiden Modi auskommentiert. **Zwei Modi** (Henne-Ei — nicht aus der
Repo-Existenz ableitbar, daher explizit gewählt):

1. **`ai-harness-init` (Voll-Kanon):** alle Blöcke **aktiv**, ohne
   Existenzprüfung — `scan.roots: [spec, docs, harness]`, sämtliche
   ids-Muster und matrix-Klassen/-Regeln aktiv. Zielbild für ein
   leeres/frisches Repo; der erste Lauf erfordert die angelegte Struktur
   (sonst Exit 2 wegen fehlender Scan-Wurzel / fehlendem ids-`target`).
2. **`ai-harness` (repo-bewusst):** nur Blöcke, deren Pfad/Target im
   gescannten Baum existiert, sind aktiv (`scan.roots` nur vorhandene);
   fehlende werden **auskommentiert mit Hinweis** ausgegeben
   (Existenz-Prüfung nur über den gescannten Baum — kein git, kein Netz).

Sind beide Tokens zugleich angegeben, gewinnt `ai-harness-init` (Voll-Kanon).

**Anforderungs-Präfix.** Nur das Anforderungs-`ids`-Muster trägt ein
projektspezifisches Präfix (`<PREFIX>`); `ADR-`/`MR-`/`slice`/Carveout sind
konventions-fest. Quelle des Präfix: die Option `--id-prefix <PREFIX>`
(explizit, gewinnt immer); im Modus `ai-harness` ohne Option das
**eindeutige** Projekt-Präfix der FA-/QA-Kennungs-Headings in
`spec/lastenheft.md` (mehrere verschiedene ⇒ Nutzungsfehler). Ohne Option
**und** ohne Ableitung (insbesondere `ai-harness-init` fürs leere Repo)
bleibt der markierte Platzhalter `<PREFIX>` mit `# TODO`-Hinweis stehen —
**kein** stiller `DC-`.

Kanonische Vorlage (Spiegel der Repo-Konvention; `ai-harness-init` gibt sie
vollständig aktiv aus, `ai-harness` nur die im Baum vorhandenen Teile):

```yaml
scan:
  roots: [spec, docs, harness]     # nur vorhandene
modules: [links, anchors, ids, matrix, codepaths, spans, hostpaths, planning]
# Weitere opt-in-Module sind situativ und hier nicht vorab aktiviert:
# external, diagrams, versions, pins, immutable, tracked, targets — Voll-Schema: d-check --print-config.
# vcs/commits brauchen eine Commit-Range und werden als Makefile-Target verteilt: d-check --print-mk.
ids:
  scope:
    roots: [spec, docs/user]       # Linkpflicht nicht über den ganzen Audit-Trail
  patterns:
    - regex: 'ADR-\d{4}'
      target: docs/plan/adr/
      link-policy: always
      exempt-paths: [CHANGELOG.md, "docs/reviews/**"]
    - regex: 'MR-\d{3}'
      target: harness/conventions.md
      link-policy: always
      exempt-paths: [CHANGELOG.md, "docs/reviews/**"]
    # <PREFIX>: via --id-prefix bzw. im ai-harness-Modus aus dem Lastenheft
    # abgeleitet; ohne beides bleibt der Platzhalter + TODO stehen.
    - regex: '<PREFIX>-(FA-[A-Z]+|QA)-\d+'
      target: spec/lastenheft.md
      link-policy: always
      exempt-paths: [CHANGELOG.md, "docs/reviews/**"]
    - regex: 'slice-\d{3}'
      target: docs/plan/planning/
      link-policy: always
      exempt-paths: [CHANGELOG.md, "docs/reviews/**"]
matrix:
  classes:
    - name: spec-straten
      paths: [spec/lastenheft.md, spec/spezifikation.md, spec/architecture.md]
      order: [spec/lastenheft.md, spec/spezifikation.md, spec/architecture.md]  # autoritativste zuerst
      direction: no-downward  # klasseninterner Abwärtsverweis ⇒ matrix-downward
    - {name: adr, paths: ["docs/plan/adr/[0-9]*.md"]}
    - {name: slice, paths: ["docs/plan/planning/**/slice-*.md"]}
  rules:
    - {from: spec-straten, to: adr, allow: false}
    - {from: spec-straten, to: slice, allow: false}
  status:
    forbidden: [superseded, deprecated]
  exclude-sections: [Historie, "7. Historie", Geschichte]
planning:
  roadmap: docs/plan/planning/in-progress/roadmap.md
codepaths:
  exempt-paths: [CHANGELOG.md, "docs/reviews/**"]
  # ignore-refs: ["tools/altes-skript.sh"]  # Ziel-Pfade entfernter Artefakte (referenz-weit; leer starten, DC-FA-CODE-001)
```

Der aktive (nicht auskommentierte) Teil dekodiert über den eigenen Parser
([`DC-FA-CONF-001`](lastenheft.md#dc-fa-conf-001--konfigurationsdatei)).

### DC-FA-CLI-007.a — Diagnose-Modus

`--doctor` ist ein **Lese**-Modus, der statt der knappen Befund-Zeilen
([`DC-FA-CLI-004`](lastenheft.md#dc-fa-cli-004--ausgabeformate)) eine
erklärende, nach Datei gruppierte Diagnose ausgibt; Schreibzugriff
entsteht nie
([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
Mit zusätzlichem `--json` wird dieselbe Diagnose **maschinenlesbar**
gerendert (Schritt 6) — ein drittes Rendering desselben Modells neben
Prosa und Patch; nur `--repair`+`--json` und `--doctor`+`--repair`
bleiben Nutzungsfehler (Exit 2). Schritte
(deterministisch — [`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus):
es wird nur über die stabil sortierte Befundliste und die geordneten
ids-Muster iteriert, keine Map-Reihenfolge):

1. **Prüflauf** wie im Default (Module nach Konfiguration/CLI; gleiche
   Befundmenge und Exit-Code-Logik 0/1/2 wie
   [`DC-FA-CLI-003`](lastenheft.md#dc-fa-cli-003--exit-codes)).
2. **Gruppieren** der Befunde nach Datei (in der bestehenden
   Sortierreihenfolge Datei → Zeile → Regel → Ziel → Grund).
3. **Klartext je Befund:** der Grund-Code wird über ein festes Mapping
   in einen Klartext übersetzt — für jeden Grund-Code aus
   [§4](#4-grund--und-fehler-codes) genau ein Eintrag, abgesichert durch
   eine Vollständigkeits-Prüfung gegen die Reason-Konstanten. Ausgegeben
   werden Zeile, Klartext, Regelmodul und Stelle (Ziel).
4. **Fix-Kandidat** (nur wo EINDEUTIG ableitbar, sonst keiner): in dieser
   Version ausschließlich für `id-unlinked` — die nackte Kennung wird als
   Markdown-Link auf das in der passenden ids-Regel deklarierte
   Definitions-`target` vorgeschlagen (relativ zum Verzeichnis der
   Befund-Datei, Datei-Ebene; den genauen Anker setzt der Mensch). Der
   Kandidat wird **nicht angewendet**; er ist zugleich die
   wiederverwendbare Eingabe des Patch-Modus (`--repair`, folgt). Best-
   Guess-Fälle (`target-missing`, `span-*` …) liefern hier bewusst keinen
   Kandidaten.
5. Bei **null Befunden** weist die Diagnose das aus (kein Kandidat). Die
   Diagnose geht auf stdout, die Zusammenfassung (geprüfte Dateien,
   Befundzahl) auf stderr — analog zum Default-Reporter.
6. **JSON-Rendering** (`--doctor --json`): statt der Prosa wird ein
   JSON-Dokument wie die [JSON-Ausgabe](#json-ausgabe---json) auf stdout
   geschrieben, dessen `findings` je Eintrag zusätzlich `reasonText`
   (der Grund-Klartext aus Schritt 3) und `fixCandidate` tragen — das
   Objekt `{original, replacement, note}` aus Schritt 4 oder explizit
   `null`, wo kein eindeutiger Kandidat existiert (nicht weggelassen,
   sonst verschwindet die Aussage „kein eindeutiger Fix"). Die
   Datei-Gruppierung trägt das bereits vorhandene `file`-Feld (keine
   Prosa-Einrückung); `summary`/`exitCode` wie die JSON-Ausgabe. stdout
   enthält nur das JSON-Dokument (Felder in fester Reihenfolge, keine
   Map-Iteration — [`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)).

### DC-FA-CLI-008.a — Reparatur-Patch

`--repair` ist ein **Lese**-Modus, der statt der Befund-Zeilen einen
**unified diff** auf stdout ausgibt (`git apply`-kompatibel); das Werkzeug
schreibt nie
([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
Zwei Stufen über Flags: **konservativ** (Default, `--repair`) und **breit**
(opt-in, `--repair-broad`, impliziert `--repair`); `--repair` ist mit
`--json` und `--doctor` nicht kombinierbar (Nutzungsfehler, Exit 2).
Schritte (deterministisch,
[`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus) — Befunde stabil
sortiert, Ausgabe nach Datei/Zeile):

1. **Prüflauf** wie im Default; Exit-Code-Logik 0/1/2 wie
   [`DC-FA-CLI-003`](lastenheft.md#dc-fa-cli-003--exit-codes). Der Patch
   erscheint auf stdout unabhängig vom Code.
2. **Edits ableiten** je Befund:
   - **konservativ** (beide Stufen): nur EINDEUTIGE Fixes — in dieser
     Version `id-unlinked` → Markdown-Link auf das ids-Definitions-`target`
     ([`DC-FA-CLI-007.a`](#dc-fa-cli-007a--diagnose-modus), gemeinsame
     Ableitung), angewandt **nur auf nackte Prosa-Vorkommen**: das
     Vorkommen wird im vorverarbeiteten Zeilentext (Inline-Code
     positionserhaltend geleert) gesucht — Vorkommen innerhalb von
     Inline-Code oder bereits in einem Link bleiben unangetastet (kein
     zerrissener Span). Mehrere nackte Vorkommen derselben Kennung auf
     einer Zeile werden alle ersetzt.
   - **breit** (`--repair-broad`): zusätzlich **Best-Guess** — in dieser
     Version `target-missing` → die im Scan-Bestand EINDEUTIG mit gleichem
     Basisnamen vorkommende Datei (relativ zur Befund-Datei). Mehrdeutige
     oder fehlende Treffer liefern keinen Edit. Best-Guess-Edits sind
     **review-pflichtig**.
3. **Hunks rendern:** je geänderte Zeile ein Null-Kontext-Hunk
   (`@@ -L,1 +L,1 @@`) unter `--- a/<datei>` / `+++ b/<datei>`; der Patch
   ist gegen den unveränderten Baum sauber anwendbar. Die
   review-pflichtig-Markierung der breiten Stufe geht auf **stderr** (wie
   Diagnose/Zusammenfassung) — damit bleibt der stdout-Patch
   `git apply`-rein. Befunde ohne Edit in der gewählten Stufe bleiben
   unangetastet und erscheinen unter
   [`DC-FA-CLI-007`](lastenheft.md#dc-fa-cli-007--diagnose-modus).

### DC-FA-CLI-009.a — Requirements Traceability Matrix

`--trace` ([`DC-FA-CLI-009`](lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix))
ist ein **Lese**-Modus; Schreibzugriff entsteht nie
([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
Ableitung (deterministisch,
[`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus); feste Sortierung, keine
Map-Iteration in der Ausgabe):

1. **Anforderungen** aus `trace.requirements.source` (Default
   `spec/lastenheft.md`) gemäß
   [§`DC-FA-REQ-001.a`](spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen):
   `trace.requirements.format` wählt Heading- oder Tabellenextraktion; beide
   liefern Kennung, Titel und Modalitäts-Eingabetext. Die Kennung muss als
   Ganz-Token/-Zelle auf `trace.requirements.id-pattern` passen (Default
   `<PREFIX>-FA-<BEREICH>-NNN` bzw. `<PREFIX>-QA-NN`, präfix-agnostisch).
2. **Referenzen** aus `trace.adrs.dir` (Default `docs/plan/adr/`; Owner-Kennung
   über den Basisnamen `trace.adrs.file-pattern`, Default `NNNN-…md` → Capture 1,
   Präfix `trace.adrs.id-prefix` Default `ADR-` → `ADR-NNNN`) und
   `trace.slices.dir` (Default `docs/plan/planning/`; Owner-Kennung
   `trace.slices.file-pattern`, Default `slice-NNN-…md` → Capture 1, Präfix
   `trace.slices.id-prefix` Default leer → `slice-NNN`): je Datei alle Vorkommen
   einer Anforderungs-Kennung (`trace.requirements.id-pattern`) sammeln; Dateien,
   deren Basisname das `file-pattern` **nicht** matcht (z. B. `README.md`),
   übersprungen.
3. **Zeile je Anforderung**: ID, Titel, sortierte ADR-/Slice-Kennungen,
   **Waise** = keine referenzierende Slice-Kennung. Default-Rendering
   Markdown-Tabelle; `--json`/`--yaml` serialisieren dieselbe Struktur
   (`requirements[]`, `total`, `orphans`).

**Config-Auflösung (opt-in `trace`-Block, [`DC-FA-CLI-009`](lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)).**
Jedes `trace`-Feld ist optional; ein abwesender/leerer Wert fällt auf den
Default oben zurück. **Kein `trace`-Block ⇒ RTM byte-identisch**
([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)) zur Fassung ohne diese
Erweiterung. **Fail-closed** (Config-Zeit, vor jedem Scan, Exit 2 mit
erklärender Meldung): jede gesetzte `id-pattern`/`file-pattern` muss ein
gültiges Regex sein, und jede gesetzte `file-pattern` muss **mindestens eine
Capture-Gruppe** tragen (sonst wäre die Owner-Kennung undefiniert). Die
konfigurierten Quellen gelten unverändert auch für
`--require-complete` ([`DC-FA-CLI-011`](lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code),
gleicher RTM-Lauf).

Fehlende Quellen: Ohne nichtleer explizite `requirements.source` liefert eine
fehlende Default-Quelle weiterhin eine leere Matrix, **kein** Fehler. Eine
nichtleer explizite Anforderungsquelle oder `format: table` aktiviert dagegen
den fail-closed Quellenvertrag aus [§`DC-FA-REQ-001.a`](spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen).
Fehlende `adrs.dir`/`slices.dir` liefern weiterhin nur keine Referenzen. Nur mit
`--doctor`/`--repair` ist `--trace` ein Nutzungsfehler (Exit 2).

### DC-FA-REQ-001.a — Anforderungsquellen (Headings und Tabellen)

Der Requirements-Extraktor wird vor Schritt 2 der RTM ausgeführt und liefert
eine nach Kennung sortierte Menge `{id, title, modalityText}`. Danach ist der
Ablauf formatunabhängig.

1. **Format-Auflösung.** Leer/abwesend oder `headings` wählt den bestehenden
   ATX-Heading-Pfad. `table` wählt den Pipe-Tabellen-Pfad. Jeder andere Wert ist
   config-zeitig Exit 2. Ein `table`-Block ist bei `format: table` Pflicht und
   bei `headings`/leer verboten (kein inert akzeptierter Block).
2. **Heading-Pfad.** Außerhalb von Fenced-Code zählt eine ATX-Überschrift nur,
   wenn `id-pattern` ihr erstes Token als Ganzes matcht
   (`FindString(tok) == tok`); Titel = Heading-Klartext ohne Kennung/Trenner,
   `modalityText` = Body bis zur nächsten gleich-/höherrangigen Überschrift.
3. **Tabellen-Lexik.** Außerhalb von Fenced-Code beginnt eine Tabelle mit einer
   Headerzeile und der unmittelbar folgenden Trennzeile. Jede Trennzelle matcht
   `^:?-+:?$` — **ein** Bindestrich genügt (GFM); führende/abschließende Pipes
   sind optional. Ein `\|` und eine
   Pipe in einem auf derselben Zeile korrekt geschlossenen Backtick-Code-Span
   teilen keine Zelle. Zellen werden getrimmt; `\|` wird danach zu `|`.
   Mehrzeilige Zellen/Fences innerhalb einer Zelle, HTML-Tabellen und
   Block-Markdown sind keine Tabelle dieses Extraktors.
4. **Header-Bindung.** `table.id-column` und genau eine von
   `table.text-column` (ein Name) oder `table.text-columns` (nichtleere,
   duplikatfreie Liste alternativer Namen) sind Pflicht;
   `table.modality-column` ist optional. Vergleich gegen die
   normalisierte Headerzelle ist exakt und case-sensitiv. Eine relevante
   Tabelle muss ID, Modalität (wenn gesetzt) und **genau einen** der
   Text-Header je genau einmal tragen; fehlt eine Rolle über alle Tabellen,
   kommt ein Header doppelt oder kommen zwei Text-Alternativen in derselben
   Tabelle vor ⇒ Exit 2. Zusätzlich muss **jede** in `text-columns`
   deklarierte Alternative über die gesamte Quelle mindestens einmal gebunden
   werden; ein unbenutzter Name ist ein Tippfehler-/Teilmenge-Guard und Exit 2.
   Mehrere relevante Tabellen werden in Quellreihenfolge zusammengeführt.
5. **Datenzeilen und Tabellengrenze.** Eine Zeile gehört bis zur ersten
   Leer-/Nicht-Tabellenzeile zur Tabelle. **Grenze am relevanten Header:** bildet
   eine Zeile mit ihrer Folgezeile einen gültigen Header + Trennzeile (Schritt 3)
   **und** bindet ihr Header eine konfigurierte Rolle (Schritt 4), ist sie der
   Header einer **neuen** Tabelle und beendet die laufende — auch bei passender
   Zellenzahl; die neue Tabelle wird ab dieser Zeile erneut erkannt. Ein Header
   **ohne** gebundene Rolle beendet die laufende Tabelle nicht (er wird als
   Datenzeile gelesen). Damit wird jede relevante Tabelle erkannt und kann nicht
   still in einer vorangehenden Tabelle verschwinden. Jede sonstige Datenzeile
   muss dieselbe Zellenzahl wie der Header tragen; sonst Exit 2. **Ausnahme
   (Direktiven-Toleranz):** eine Datenzeile mit **genau einer** überzähligen Zelle,
   die ganz aus einem HTML-Kommentar besteht (`<!-- … -->` hinter der letzten
   Pipe), wird auf Header-Breite gelesen — GFM-nachsichtig für Body-Zeilen; die
   eigene Direktiven-Konvention `<!-- d-check:ignore … -->` in einer Tabellenzeile
   bricht den Reader damit nicht mehr ab. Zwei überzählige Zellen oder eine
   überzählige Nicht-Kommentar-Zelle bleiben Exit 2. Nur wenn die
   getrimmte ID-Zelle als Ganzes auf `id-pattern` passt, definiert sie eine
   Anforderung; andere Datenzeilen werden ignoriert. Titel = Inhalt der in dieser
   Tabelle gebundenen Textspalte. `modalityText` = Inhalt von `modality-column`,
   wenn gesetzt, sonst Inhalt derselben Textspalte.
6. **Duplikate und Nullmenge.** `table.duplicate-ids` ist `error` (Default),
   `first` oder `last`. Bei `first` bleibt die erste Definition, bei `last`
   überschreibt die spätere Titel und `modalityText`; die ID-Reihenfolge bleibt
   wegen der abschließenden Sortierung unabhängig davon deterministisch. Jeder
   andere Wert ist config-zeitig Exit 2. `strictSource = (nichtleerer, expliziter
   requirements.source) ∨ (format == table)`. Bei `strictSource` gilt:
   fehlende/unlesbare Quelle, doppelte erkannte ID unter Politik `error` oder null erkannte
   Anforderungen ⇒ Exit 2. Die Nullmengenmeldung nennt Quelle und Format.
   Ohne `strictSource` behält der Heading-Pfad die alte Semantik: fehlende
   Default-Quelle/null Treffer ⇒ leere RTM; doppelte ID ⇒ erster Treffer
   gewinnt. `source: ""` gilt wie abwesend und aktiviert `strictSource` nicht.
7. **Modalität.** Bei aktivem `trace.requirements.modality` klassifiziert der
   bestehende Matcher aus [§`DC-FA-MOD-001.a`](spezifikation.md#dc-fa-mod-001a--modalitäts-klassifikation-tracerequirementsmodality)
   `modalityText`; Keywords, Normalisierung, `unknown` und `require-levels`
   bleiben unverändert. Ohne aktive Modalität wird der Text nicht ausgewertet.

**Fehlerpräzedenz:** striktes YAML/Format-/Block-Schema → Quelle lesen →
Tabellenstruktur/Header → Duplicate-ID → Nullmenge → Referenz-/Coverage-
Scans. Der erste Fehler beendet den Lauf vor dem Reporter; `--trace` und
`--trace --require-complete` liefern Exit 2 und keine RTM.

### DC-FA-XREF-001.a — Kreuzverweis-Konsistenz (`cross-consistency`)

Der opt-in `trace.cross-consistency`-Abgleich
([`DC-FA-XREF-001`](lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in))
läuft nach der RTM-Verrechnung (Schritt 2 von
[§`DC-FA-CLI-009.a`](spezifikation.md#dc-fa-cli-009a--requirements-traceability-matrix))
und ändert die RTM selbst nicht. Verrechnung (deterministisch,
[`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)):

1. **Aktivierung.** Nur bei vorhandenem `trace.cross-consistency`. `forward` und
   `backward` sind beide Pflicht; fehlt einer, ist die Config Exit 2. Ohne Block
   bleibt die RTM byte-identisch.
2. **Vorwärts-Menge `F`.** Aus `forward.file`, auf `forward.sections` /
   `forward.exclude-sections` gespannt (Span-Semantik wie
   [§`DC-FA-COV-001.a`](spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)),
   wird jede relevante Tabelle über den header-gebundenen Reader
   ([§`DC-FA-REQ-001.a`](spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen))
   gelesen: `req-column` liefert über `forward.req-pattern` die Anforderungs-IDs
   (bei `ranges: true` range-/enum-expandiert), `design-column` die
   Design-Artefakt-IDs über alle `design-pattern`-Treffer der Zelle.
   **Scope-Trennung:** `forward.req-pattern` fällt per Default auf
   `requirements.id-pattern` zurück, ist aber **eigenständig** — welche
   Anforderungen der Abgleich vergleicht, entscheidet allein das Muster, nicht die
   RTM-Mitgliedschaft. Eine ID, die das Muster trifft, aber keine RTM-Zeile hat,
   wird verglichen; eine RTM-Anforderung, die das Muster nicht trifft, nicht. Ergebnis: je Anforderung `R` die Menge
   `F(R)` der genannten Design-Artefakte.
3. **Rückwärts-Menge `B`.** Aus `backward.file`, auf `backward.sections` gespannt,
   zählt jede Tabelle mit einem `edge-column`-Header. Die Artefakt-ID ist der erste
   Treffer von `forward.design-pattern` in der **ersten Spalte**
   (`artifact-id-column: first`; alternativ ein Header-Name), da der ID-Spalten-
   Header über die Tabellen variiert; die `edge-column`-Zelle liefert die
   Anforderungs-IDs (`req-pattern`, bei `ranges: true` expandiert). Invertiert: je
   so genannter Anforderung `R` wird die Artefakt-ID in `B(R)` aufgenommen.
   **Vorbedingung:** Vorwärts- (`design-column`) und Rückwärts-Artefakt-IDs werden
   bewusst mit **demselben** `forward.design-pattern` extrahiert und liegen damit im
   selben Namensraum — sonst wäre der Mengen-Diff (Schritt 6) inhärent leer/voll und
   bedeutungslos; die Vakuitäts-Prüfung (Schritt 5) fängt genau diesen Fall.
4. **Ausschluss.** Anforderungs-IDs, die auf `exclude-req` (RE2) passen, werden vor
   dem Diff aus den Schlüsselmengen von `F` und `B` entfernt (Ableitungssprünge in
   Mittelschichten).
5. **Vakuitäts-Prüfung** — **nach** dem Ausschluss (Schritt 4) gemessen, denn
   maßgeblich ist, was am Ende tatsächlich verglichen wird: ein `exclude-req`, das
   alle Anforderungen verschluckt, schaltet das Gate ebenso still ab wie ein
   fehlgreifendes Muster (das Ventil ist selbst eine kuratierte, drift-fähige
   Kante). Ein Abgleich, aus dem **keine Kante** fällt, ist kein
   bestandener Abgleich: die Muster greifen dann am Inhalt vorbei, und ein
   `0 Differenz(en)`/Exit 0 behauptete eine nie geprüfte Konsistenz. **Exit 2**,
   wenn (a) `F` **und** `B` leer sind — der typische Anlass ist ein
   `design-pattern`, das kompiliert, aber am Artefakt-Namensraum vorbeigreift und
   damit **beide** Sichten leerräumt (es ist geteilt, Schritt 3) — oder (b) `B`
   leer ist **und** `mode: superset` gilt, weil dann allein `B \ F` gatet und
   konstruktionsbedingt nie ein Befund entstehen kann. **Nicht** vakuum ist eine
   **einseitig** leere Sicht mit nicht-leerer Gegenseite: der Diff (Schritt 6)
   läuft über `keys(F) ∪ keys(B)` und ist dafür wohldefiniert — eine noch
   unrestrukturierte Vorwärts-Sicht bei gepflegten Rück-Kanten meldet jede
   Rück-Kante als `B \ F` und ist der erwartete Bootstrap-Zustand, kein Fehler.
6. **Diff.** Für jede Anforderung `R ∈ keys(F) ∪ keys(B)`: `F(R) \ B(R)` und
   `B(R) \ F(R)`. `mode: equal` meldet beide Richtungen, `mode: superset` nur
   `B(R) \ F(R)`. Jede Differenz ist ein Befund mit Quell-`Datei:Zeile`,
   Anforderungs-ID, fehlender/überzähliger Artefaktmenge und Richtungslabel;
   Befunde werden nach `(R, Artefakt-ID, Richtung)` deterministisch sortiert.
7. **Gatung.** Der Abgleich ist advisory: `--trace` bleibt Exit 0. Nur das globale
   `--require-complete`
   ([`DC-FA-CLI-011`](lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code))
   wertet die Befunde als Gate (Exit 1 bei ≥1 Differenz).

**Fehlerpräzedenz:** Config-Schema (Pflichtblöcke, `mode`, kompilierende
`design-pattern`/`req-pattern`/`exclude-req`) → Quellen lesen (fehlende
`forward.file`/`backward.file` ⇒ Exit 2,
[`DC-FA-CLI-003`](lastenheft.md#dc-fa-cli-003--exit-codes)) → Abschnitts-Spannung
(konfigurierter Sektionsname ohne Heading-Treffer ⇒ Exit 2) → Header-Bindung
(`req-column`/`design-column`/`edge-column` sowie `artifact-id-column`, wenn ≠
`first`, je genau einmal) → Range-Expansion (ungültige Range ⇒ Exit 2) →
Ausschluss (Schritt 4) → Vakuität (Schritt 5) → Diff (Schritt 6). Jede Stufe läuft
**über beide Sichten**, bevor die nächste beginnt (sonst verdeckte ein
Vorwärts-Fehler einen früherstufigen Rückwärts-Fehler); der erste Fehler beendet
den Lauf vor dem Reporter.

### DC-FA-COV-001.a — Kuratierte Coverage-Quellen (`trace.coverage`)

Die opt-in Liste `trace.coverage`
([`DC-FA-COV-001`](lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in))
erweitert die RTM-Verrechnung (Schritt 2 von
[§`DC-FA-CLI-009.a`](spezifikation.md#dc-fa-cli-009a--requirements-traceability-matrix))
um eine dritte Referenzklasse. Verrechnung (deterministisch,
[`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)):

1. **Je Quelle** (`{files, label, ranges, sections, exclude-sections}`) je Datei
   in `files`: den Datei-Inhalt **abschnitts-filtern** — ein Sektionsname ist der
   **volle Heading-Klartext** (exakter Vergleich wie `matrix.exclude-sections`,
   z. B. `"27.1 Anforderung zu Design"`, **nicht** die Kurzform `"27.1"`):
   - `sections` gesetzt ⇒ **nur** die Zeilen der gelisteten Abschnitte behalten
     (Span wie `matrix.exclude-sections`: Überschrift bis zur nächsten
     gleich-/höherrangigen); leer ⇒ ganze Datei;
   - `exclude-sections` ⇒ die Zeilen der gelisteten Abschnitte **entfernen**
     (dieselbe Span-Mechanik, `excludedRanges`). Reihenfolge: erst Whitelist, dann
     Blacklist. So ergibt `sections: ["27.1 Anforderung zu Design"]` +
     `exclude-sections: ["27.1.1 Anforderungen ohne Design-Artefakt"]` genau
     „§27.1 ohne §27.1.1".
   - **Fail-closed:** ein konfigurierter Sektionsname (in `sections` **oder**
     `exclude-sections`), der über alle `files` der Quelle **keine** Überschrift
     trifft ⇒ **Exit 2** (Tippfehler-/Kurzform-Guard — sonst blankt eine
     Whitelist ohne Treffer still die ganze Datei bzw. greift eine Blacklist
     nicht, und Anforderungen würden falsch (un)gedeckt).
2. **ID-Extraktion** aus dem gefilterten Text: alle exakten
   `requirements.id-pattern`-Treffer; bei `ranges: true` zusätzlich die
   **Range-/Enum-Expansion** (Schritt 3).
3. **Range-Parser** (isolierte, deterministische Funktion, parametrisiert über
   `id-pattern`): für jede `id-pattern`-Fundstelle prüfen, ob unmittelbar — oder
   **genau ein** Markdown-Link-Suffix `](…)` später (**link-transparent**, siehe
   unten) —
   - `..<Ziffern>` folgt (`<FAM>-AAA..BBB`): Familie = Fundstelle ohne
     Trailing-`-<Ziffern>`, Start `AAA` = diese Trailing-Ziffern (Breite `W`),
     Ende `BBB` = die Folge-Ziffern; expandiere `<FAM>-<i>` für `i ∈ [AAA, BBB]`,
     **`W`-breit null-aufgefüllt**; `AAA>BBB` oder `len(BBB)≠W` ⇒ **Exit 2**;
   - `/<Ziffern>`-Folgen folgen (`<FAM>-AAA/BBB/CCC`): je `<FAM>-<Ziffern>`;
   jede expandierte ID wird gegen `id-pattern` geprüft und bei Nicht-Treffer
   **verworfen**.
   **Komma-Kurzform ⇒ Exit 2.** Folgt der Fundstelle — **oder der von ihr
   ausgehenden zugesagten Notation** (`..BBB`/`/BBB`, nachdem diese konsumiert ist)
   — nach der Link-Transparenz ein Komma und **unmittelbar darauf Ziffern**
   (`<FAM>-AAA, BBB`; ebenso `<FAM>-AAA..CCC, DDD` und `<FAM>-AAA/CCC, DDD`), ist
   das eine Aufzählungs-**Gestalt** ohne zugesagte Notation: **Exit 2** mit Hinweis
   auf `/` und `..`. Weder still verschlucken (die Kurzform verschwände und erzeugte
   eine falsche Waise) noch expandieren (`GG-QA-001, 007 Sekunden` wäre eine
   geratene Absicht) — dieselbe Logik wie bei `AAA>BBB`: die Gestalt triggert, der
   Inhalt ist ungültig. Die Prüfung greift **an jeder der drei Positionen**
   (hinter der nackten Kennung, hinter einer Range, hinter einem Enum) auf dem
   nach der Notation verbleibenden Rest. Ein Komma vor einer **vollständigen**
   Kennung (`<FAM>-AAA, <FAM>-BBB`, auch `<FAM>-AAA..CCC, <FAM>-DDD`) ist **keine**
   Kurzform und unberührt; beide Kennungen werden regulär gefunden.

   **Link-Transparenz.** Steht die Kennung unter Linkpflicht
   ([`DC-FA-ID-001`](lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)),
   trägt sie ein Link-Suffix, und die Fortsetzung folgt der Fundstelle nicht mehr
   unmittelbar: `` [`GG-UI-001`](…)..009 ``. Der Parser überspringt daher hinter
   der Fundstelle **höchstens einmal** ein unmittelbar anschließendes Suffix der
   Form `` `](…) `` (optionales schließendes Code-Span-Backtick, `]`, geklammertes
   Ziel) und prüft **dahinter** wieder **unmittelbar** auf `..`/`/`. Ohne diese
   Regel bräche die unqualifizierte Range-Zusage des Lastenhefts genau dort, wo
   d-checks eigene Linkpflicht greift.
   **Nicht** übersprungen werden — bewusst, weil jede weitere Toleranz die
   Autor-Absicht rät statt sie zu lesen: Whitespace, Emphasis, ein zweites
   Link-Suffix, beliebiger Text zwischen `)` und der Fortsetzung. Die ID selbst
   darf im Code-Span des Linktexts stehen (das leistet bereits die
   `id-pattern`-Suche).
4. **Zuordnung**: jede abgedeckte Anforderungs-Kennung erhält das `label` der
   Quelle in ihrer **Coverage**-Menge (dedupliziert, sortiert; mehrere Quellen ⇒
   mehrere Labels).
5. **Waise-Neubestimmung** ([`DC-FA-CLI-011`](lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)):
   `orphan = (keine Slice-Referenz) ∧ (keine Coverage-Referenz)`. ADR-Referenzen
   zählen wie bisher **nicht** zur Waisen-Freiheit.

**Rendering/Config-Auflösung.** Ist ≥1 `trace.coverage`-Quelle konfiguriert,
trägt die Markdown-RTM eine zusätzliche **Coverage**-Spalte **vor** der
`Status`-Spalte (`… | Slices | Coverage | Status |`; Labels, sonst „—");
`--json`/`--yaml` tragen je Anforderung ein `coverage`-Feld (`omitempty` — leer
⇒ entfällt). Die Spalten-Entscheidung hängt an einem expliziten Aktiv-Flag der
Matrix (≥1 Quelle konfiguriert), nicht am Zeilen-Inhalt. **Kein `trace.coverage`
⇒ kein Flag, keine Spalte, kein Feld ⇒ RTM byte-identisch.** **Fail-closed**
(Config-/Lauf-Zeit, Exit 2): leere `files`-Liste, ein `files`-Pfad außerhalb der
Repo-Wurzel (führendes `/` oder `..`) oder eine **fehlende** `files`-Datei
(anders als `adrs.dir`/`slices.dir`, wo Fehlen = Skip — `files` sind **explizit
benannt**), leeres `label`, ein Sektionsname ohne Heading-Treffer, oder eine
ungültige Range (`AAA>BBB`/Breite). `trace.coverage` führt **kein** eigenes Regex
(nutzt `requirements.id-pattern`, bereits config-zeitig geprüft). Read-only
([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)) —
nur Markdown-Lesen der benannten `files` im gemounteten Baum, kein neuer
Eingabe-Scope darüber hinaus.

### DC-FA-MOD-001.a — Modalitäts-Klassifikation (`trace.requirements.modality`)

Der opt-in Block `trace.requirements.modality`
([`DC-FA-MOD-001`](lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in))
weist jeder RTM-Anforderung eine **Stufe** zu. Ableitung (deterministisch,
[`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)):

0. **Aktiv-Prädikat.** `modality` gilt als aktiv, sobald der **Schlüssel
   vorhanden** ist — auch als leere Map `modality: {}` (dann greifen die
   Defaults). **Nicht** an `len(levels)>0` hängen (sonst wäre `modality: {}`
   fälschlich inaktiv). Fehlt der Schlüssel ⇒ inaktiv, byte-identisch.
1. **Keyword-Menge.** `modality.levels` (Map Stufe→Keywords); leer ⇒ die
   **kanonische Built-in-Default-Menge** (fix, reproduzierbar):
   - `must`:   `MUSS`, `MUESSEN`, `MÜSSEN`, `DARF NICHT`, `DÜRFEN NICHT`, `MUST`, `SHALL`, `MUST NOT`, `SHALL NOT`
   - `should`: `SOLLTE`, `SOLLTEN`, `SOLLTE NICHT`, `SOLLTEN NICHT`, `SHOULD`, `SHOULD NOT`
   - `may`:    `KANN`, `KÖNNEN`, `MUSS NICHT`, `MÜSSEN NICHT`, `MAY`, `OPTIONAL`

   Config-zeitig **fail-closed** (Exit 2): leerer Stufen-Name; leeres Keyword;
   der reservierte Stufen-Name **`unknown`** in `levels` (kollidiert mit dem
   Kein-Treffer-Sentinel); **dasselbe Keyword in mehr als einer Stufe** (sonst
   hinge die gewinnende Stufe an der Map-Iteration ⇒ Nondeterminismus, verletzt
   [`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)).
2. **Body-Abschnitt + Normalisierung.** Je Anforderung der Text von ihrer
   Heading-Zeile bis zur nächsten gleich-/höherrangigen Überschrift (Span-Mechanik
   wie `matrix.exclude-sections`/`SelectSections`), dann **normalisiert**:
   Markdown-Emphasis (`*`/`` ` ``) entfernt (wie `StripHeadingLinks` bei Headings)
   und **Whitespace-Folgen inkl. Zeilenumbrüchen zu je einem Leerzeichen**
   zusammengezogen. Ohne diese Normalisierung matchte eine umbrochene/emphasierte
   Phrase `**MUSS** NICHT` bzw. `MUSS\nNICHT` **nicht** und fiele still auf `MUSS`
   (must) statt `may` zurück.
3. **Erster/längster Treffer.** Über den normalisierten Abschnitt wird die
   **früheste Position** gesucht, an der **irgendein** Keyword matcht; bei
   mehreren Keywords an **derselben** Position gewinnt die **längste** Phrase
   (`DARF NICHT` vor `DARF`; `MUSS NICHT` [may] vor `MUSS` [must]). Matching
   **case-insensitiv** und **wortgrenzen-genau** — implementiert als
   `(?i)\bKEYWORD\b` mit regex-gequotetem Keyword. **Caveat:** RE2-`\b` ist
   **ASCII**-basiert; die Default-Keywords haben ausschließlich ASCII-Ränder
   (`MÜSSEN`/`KÖNNEN`/`DÜRFEN NICHT` beginnen/enden ASCII) und sind sicher; ein
   **konfiguriertes** Keyword mit führendem/schließendem Umlaut träfe wegen der
   ASCII-`\b` nicht — dann Wortgrenze weglassen bzw. das Keyword anders fassen.
   Die Stufe des gewinnenden Keywords ist die Modalität; **kein** Treffer ⇒
   `unknown`.
4. **Rendering.** Bei aktivem `modality` trägt die Markdown-RTM eine **Modality**-
   Spalte (Stufe, sonst „—") **vor** `Status` (nach `Coverage`; Voll-Ordnung
   `Anforderung|Titel|ADRs|Slices|[Coverage]|[Modality]|Status`); `--json`/`--yaml`
   ein `modality`-Feld (`omitempty`). Ohne `modality`: keine Spalte, kein Feld ⇒
   byte-identisch.
5. **Gating** ([`DC-FA-CLI-011`](lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)):
   `--require-complete` zählt eine Waise (¬slice ∧ ¬coverage) nur dann als
   gatend (Exit 1), wenn ihre Stufe in `modality.require-levels` liegt (Default
   `[must]`; `unknown` gatet nur, wenn explizit gelistet). Ohne `modality` gaten
   **alle** Waisen (byte-identisch). Die **Status**-Spalte bleibt für jede Waise
   „WAISE" (die Modalität steht in ihrer eigenen Spalte); die stderr-Zeile von
   `--require-complete` nennt die **gatende** Zahl (`N gatende von M Waisen`),
   nicht nur die Gesamtzahl. Fail-closed: ein `require-levels`-Eintrag, der weder
   deklarierte Stufe noch `unknown` ist ⇒ Exit 2.

Read-only ([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)) —
nur Body-Lesen im gemounteten Baum. **Fail-open-Grenze** (bewusst): eine
Anforderung mit **unaufgeführtem** Modal-Verb fällt auf `unknown` (advisory bei
Default-`require-levels`) — sichtbar in der Spalte, per `require-levels: [must, unknown]`
strikt gatbar.

### DC-FA-CLI-010.a — Makefile-Fragment

`--print-mk` ([`DC-FA-CLI-010`](lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben))
ist ein **Lese**-Modus wie `--print-config`: kein Repo-Zugriff, keine Datei
geschrieben ([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit));
Kurzschluss vor der Wurzel-Auflösung. Ausgabe (deterministisch,
[`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus) — hängt nur an der
eingebetteten Version):

1. Kommentar-Kopf: Einbindung via `include`, Hinweis zum Digest-Pin via
   `DCHECK_IMAGE` bzw. bequemer `DCHECK_DIGEST`.
2. `DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v<version>` — `<version>` ist die
   beim Tag-Build via `-ldflags -X …/cli.version=<tag>` eingebettete
   Release-Version (Default `0.0.0-dev` für lokale/Gate-Builds).
3. `DCHECK_DIGEST ?=` plus ein `ifeq`-Block: ein gesetzter `DCHECK_DIGEST` (ein
   `@sha256:`-Digest) leitet `DCHECK_REF` auf `…/d-check@$(DCHECK_DIGEST)` um und
   **sticht** so den Tag von `DCHECK_IMAGE`; sonst `DCHECK_REF := $(DCHECK_IMAGE)`.
4. `TRACE_FLAGS ?=` — überschreibbare Flag-Variable für die RTM-Targets.
5. Elf `.PHONY`-Targets, jeweils mit `##`-Annotation (die das `help` des
   Konsumenten aufgreift) und **TAB**-eingerücktem
   `docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) …`:
   - `doc-check` (Befund-Gate, ohne Zusatz-Flag),
   - `doc-trace` (`--trace $(TRACE_FLAGS)`, advisory RTM,
     [`DC-FA-CLI-009`](lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)),
   - `doc-complete` (`--trace --require-complete $(TRACE_FLAGS)`,
     [`DC-FA-CLI-011`](lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)),
   - `doc-doctor` (`--doctor`,
     [`DC-FA-CLI-007`](lastenheft.md#dc-fa-cli-007--diagnose-modus)),
   - `doc-repair` (`--repair`,
     [`DC-FA-CLI-008`](lastenheft.md#dc-fa-cli-008--reparatur-patch)) — das Recipe ist
     mit `@` **echo-unterdrückt**, damit stdout ein `git apply`-reiner Patch bleibt
     (sonst verunreinigte die von `make` auf stdout gedruckte Recipe-Zeile den Patch),
   - `doc-immutable` (`--enable vcs` + auf `vcs` fokussierte `--disable`-Liste,
     [`DC-FA-VCS-001`](lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in));
     `$(if $(STAGED),--staged,--range $(RANGE))` — die Range liefert der Konsument
     aus seinem CI, `STAGED=1` für einen lokalen pre-commit-Lauf. Die `--disable`-
     Liste wird aus dem Modulsatz (`ValidModules` ohne `vcs`) **abgeleitet**, nicht
     dupliziert, damit das Target nur `vcs` läuft (sonst über-feuerten die
     Doc-Module des Konsumenten auf Nicht-ADR-Befunde). Die **verteilte** git-Form
     der Immutabilität — kein kopiertes Skript (der Antrieb hinter dem Modul `vcs`),
   - `doc-commits` (`--enable commits` + auf `commits` fokussierte `--disable`-Liste,
     [`DC-FA-COMMITS-001`](lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in));
     `--range $(RANGE)` — der Konsument liefert die Range aus seinem CI. Die
     `--disable`-Liste wird analog zu `doc-immutable` aus dem Modulsatz
     (`ValidModules` ohne `commits`) **abgeleitet**, damit das Target nur `commits`
     läuft. Die **verteilte** Commit-Traceability — kein kopiertes Skript (der
     Antrieb hinter dem Modul `commits`, dieselbe Linie wie `doc-immutable`),
   - `doc-planning` (`--enable planning` + auf `planning` fokussierte `--disable`-Liste,
     [`DC-FA-PLAN-001`](lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in));
     **kein** `--range`/`--staged` (hermetisch — nur der Arbeitsbaum). Die
     `--disable`-Liste analog aus `ValidModules` ohne `planning` abgeleitet; verteilt
     die Planning-Lifecycle-Konsistenz an Konsumenten mit demselben Roadmap-Layout,
   - `doc-tracked` (`--enable tracked` + auf `tracked` fokussierte `--disable`-Liste,
     [`DC-FA-TRK-001`](lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in));
     **kein** `--range`/`--staged` (der Index-Stand des gemounteten Arbeitsbaums —
     das `.git` liegt im read-only-Mount). Die `--disable`-Liste analog aus
     `ValidModules` ohne `tracked` abgeleitet; verteilt den
     Getrackt-Status-Check an Konsumenten ohne Skript-Kopie,
   - `doc-targets` (`--enable targets` + auf `targets` fokussierte `--disable`-Liste,
     [`DC-FA-TGT-001`](lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in));
     **kein** `--range`/`--staged` (hermetisch — nur der Arbeitsbaum: Makefile +
     Doku-Tabellen). Die `--disable`-Liste analog aus `ValidModules` ohne `targets`
     abgeleitet; verteilt die Deklarations-Konsistenz-Prüfung an Konsumenten ohne
     Skript-Kopie,
   - `doc-help` (listet die `doc-*`-Targets via `grep … $(MAKEFILE_LIST) | sed`
     über die `##`-Annotationen; **namespaced** statt `help`, um die
     Namens-Kollision mit dem Konsumenten-Makefile zu vermeiden).

Der eigene Image-**Digest** wird NICHT eingebettet (er hasht das Binary
selbst — Henne-Ei); der Konsument pinnt per `DCHECK_IMAGE`- bzw.
`DCHECK_DIGEST`-Override.

### DC-FA-CLI-011.a — Vollständigkeits-Prüfung (`--require-complete`)

`--require-complete`
([`DC-FA-CLI-011`](lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code))
ist ein **Modifikator von `--trace`**, kein eigener Modus: der RTM-Lauf
([`DC-FA-CLI-009`](lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix))
und seine Ausgabe (Markdown bzw. `--json`/`--yaml`) bleiben **byte-identisch**
([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)); allein der Exit-Code
ändert sich:

1. `matrix.orphans > 0` ⇒ **Exit 1** (Befund-Code,
   [`DC-FA-CLI-003`](lastenheft.md#dc-fa-cli-003--exit-codes)); zusätzlich eine
   deterministische Zähl-Zeile auf **stderr** (die RTM bleibt auf stdout rein).
2. `matrix.orphans == 0` ⇒ Exit 0.
3. `--require-complete` ohne `--trace` ⇒ Nutzungsfehler (Exit 2), behandelt in
   der Kombinations-Prüfung vor jedem Repo-Zugriff.

Der Default-`--trace` ohne dieses Flag bleibt **advisory** (Exit 0 auch bei
Waisen) — die Durchsetzung ist strikt opt-in und ändert den
[`DC-FA-CLI-009`](lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)-Vertrag
nicht.

### DC-FA-CLI-012.a — Konfigurations-Pfad (`--config`)

`--config <datei>`
([`DC-FA-CLI-012`](lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben))
verschiebt **ausschließlich die Herkunft** der Konfiguration; Format, Validierung
und Fehlerverhalten sind die von
[`DC-FA-CONF-001`](lastenheft.md#dc-fa-conf-001--konfigurationsdatei). Der Ablauf
setzt vor Schritt „Konfiguration laden" von
[`DC-FA-CLI-001.a`](#dc-fa-cli-001a--ablauf-eines-prüflaufs) an:

1. **Auflösen.** Ohne `--config` bleibt die Quelle der konventionelle Dateiname in
   der Scan-Wurzel; die Auflösung ist unverändert und der Befundsatz
   byte-identisch ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)). Mit
   `--config` wird der Wert **relativ zur Scan-Wurzel** aufgelöst (ein absoluter
   Pfad wird auf die Wurzel bezogen geprüft).
2. **Wurzel-Constraint (fail-closed).** Liegt der aufgelöste Pfad — nach
   Auflösung von `..`-Segmenten — **außerhalb** der Scan-Wurzel, ⇒ **Exit 2**.
   Begründung: der Lauf ist auf den read-only-Mount beschränkt; eine Datei
   außerhalb wäre im Container ohnehin unerreichbar, und ein stiller Fehlschlag
   dort wäre ein Silent-Grün. Dieselbe Constraint wie bei `scan.roots` und
   `planning.roadmap`.
3. **Existenz (fail-closed).** Existiert die Datei nicht oder ist sie unlesbar ⇒
   **Exit 2** mit Hinweis auf stderr. Es gibt **keinen** Rückfall — weder auf die
   Defaults noch auf eine vorhandene konventionelle Datei. Ein vertipptes Profil
   darf nicht unbemerkt einen anderen Prüfumfang fahren.
4. **Ersetzen, nicht ergänzen.** Die angegebene Datei ist die **einzige** Quelle:
   eine konventionelle Datei in der Scan-Wurzel wird nicht gelesen und nicht
   zusammengeführt. Es gibt keine Vererbung über Verzeichnisebenen
   ([`DC-FA-CONF-001`](lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
   §Out-of-Scope bleibt gültig).
5. **Validierung unverändert.** Der Inhalt durchläuft dieselbe strikte
   Schema-/Inhalts-Validierung; jeder Konfigurationsfehler ⇒ Exit 2 mit
   Zeilenangabe. Module, Defaults und Grund-Codes verhalten sich identisch — der
   Schalter ist eine **Herkunfts**-, keine Semantik-Änderung.
6. **Provenance in Befunden.** Wo ein Befund heute die Konfigurationsdatei als
   Herkunft benennt, benennt er die **tatsächlich geladene** Datei — nicht den
   konventionellen Namen (sonst zeigte die Meldung auf eine Datei, die der Lauf
   nie gelesen hat).

### DC-FA-REF-001.a — Geteiltes Referenz-Ventil (`ignore-refs`)

Anwendbar in `links`
([DC-FA-LINK-001.a](#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)),
`anchors` ([DC-FA-ANCH-001.a](#dc-fa-anch-001a--github-slug-algorithmus)) und
`codepaths` ([DC-FA-CODE-001.a](#dc-fa-code-001a--pfade-in-inline-code)). Betrachtet
werden je Referenz die **Quelldatei** `S` (Wurzel-relativ) und der im jeweiligen
Modul **aufgelöste** Ziel-Pfad `T` (der Befund-`target`; Fragmente sind bereits
abgetrennt).

**Match-Prädikat.** Die Referenz wird vom Ventil **ignoriert**, wenn **mindestens
ein** `ignore-refs`-Eintrag `E` erfüllt:

```text
(E.in fehlt  ∨  match(E.in, S))  ∧  match(E.refs, T)  ∧  ¬match(E.keep, T)
```

— der Quell-Skopus trifft (oder ist offen), ein `refs`-Glob trifft das aufgelöste
Ziel, und **kein** `keep`-Glob **desselben Eintrags** trifft es. `keep` gewinnt
innerhalb seines Eintrags **unbedingt und reihenfolge-unabhängig** (kein
gitignore-Last-Match). Alle Globs matchen wie `scan.ignore` (`matchGlob`,
segmentweises `**`), gegen den **aufgelösten** Pfad bzw. die Quelldatei. Mehrere
Einträge wirken additiv (Union der Ignorier-Prädikate).

**Wirkung.** Ignoriert → das Modul überspringt für dieses Ziel die **Existenz**-, die
**Repo-Escape**- und (bei Markdown-Zielen) die **Anker**-Prüfung (kein
`target-missing`/`codepath-missing`/`repo-escape`/`anchor-missing`) — dieselbe
Unterdrückung, die das modul-lokale `codepaths.ignore-refs` bisher leistete, jetzt
geteilt. Die **Symlink-Ablehnung**
([`DC-FA-LINK-002`](lastenheft.md#dc-fa-link-002--symlink-ablehnung), Vorrang, greift
nur an existierenden Symlink-Zielen) bleibt unberührt; andere Module ebenso. Ein
Eintrag mit leerer/fehlender `refs`-Liste ist inert. Ohne jeden `ignore-refs`-Eintrag
ist der Befundsatz byte-identisch
([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)).

**Alias `codepaths.ignore-refs`.** Die modul-lokale Liste wird wie ein
`ignore-refs`-Eintrag ohne `in`/`keep` behandelt, dessen `refs` die Liste ist,
angewandt **nur** in `codepaths` — die Fassung vor dieser Anforderung
([DC-FA-CODE-001.a](#dc-fa-code-001a--pfade-in-inline-code)) bleibt damit
byte-identisch. Top-Level-`ignore-refs` und Alias wirken additiv.

**Achsen-Präzedenz.** `scan.ignore` (Quell-Achse) entfernt ganze Dateien **vor** jeder
Modul-Prüfung; deren Referenzen erreichen das Ventil nie. Der Zeilen-Marker
`d-check:ignore` (nur `codepaths`, Schritt 1) wirkt **vor** der Pfad-Erkennung.
`ignore-refs` (Ziel-Achse) wirkt **bei der Auflösung** des Ziels. Die drei Achsen sind
orthogonal; keine überschreibt eine andere.

---

### DC-FA-LINK-001.a — Markdown-Vorverarbeitung und Link-Extraktion

1. **Fences:** Zeilen, deren erste Nicht-Leerzeichen-Folge mit
   ` ``` ` oder `~~~` beginnt, schalten den Fence-Zustand um;
   Zeilen im Fence-Zustand werden von allen Modulen ignoriert.
   **Infozeilen-Regel (CommonMark):** eine ` ``` `-Zeile, deren Rest
   einen **Backtick** enthält, ist **kein** Fence-Öffner, sondern
   Fließtext — sie schaltet nichts um. Ohne die Regel öffnet ein Satz
   **über** einen Fence (```` ```yaml-Fence (`datei.md`) — … ````) einen
   Fence-Zustand, der bis zum nächsten Öffner oder zum Dateiende reicht,
   und blendet **alle** Module für diesen Bereich: Befunde verschwinden
   lautlos (Exit 1 ⇒ Exit 0). Für `~~~`-Öffner gilt die Regel **nicht**
   (CommonMark erlaubt dort Backticks in der Infozeile).
2. **Inline-Code:** Backtick-Spans werden durch Leerzeichen gleicher
   Länge ersetzt (positionserhaltend — angrenzender Text kann nicht
   zu Schein-Vorkommen verschmelzen); die öffnende Backtick-Folge
   bestimmt die schließende (Mehrfach-Backticks). Die Erkennung ist
   **absatzweise** (CommonMark): ein Span darf Zeilenumbrüche
   enthalten; Absatzgrenzen sind Leerzeilen und Fences, eine im
   Absatz ungeschlossene Backtick-Folge ist literal. Damit invertiert
   ein über den Zeilenumbruch gebrochener Span nicht die
   Backtick-Parität der Folgezeile.
3. **Extraktion:** Inline-Links `[text](ziel)` und Bilder
   `![alt](ziel)` per Klammer-balancierter Suche; mehrere Links pro
   Zeile werden alle erfasst. Ziele in `<…>` werden entquotet; ein
   Titel-Suffix (` "…"`) wird abgetrennt. Die Extraktion ist
   **zeilenbasiert**: Inline-Links, deren Syntax sich über einen
   Zeilenumbruch erstreckt (GFM-Soft-Break im Linktext), werden nicht
   erkannt — normative Grenze für alle Module.
4. **Ziel-Normalisierung:** Prozent-Dekodierung (RFC 3986, vollständig)
   → Auflösung relativ zum Verzeichnis der enthaltenden Datei →
   lexikalische Normalisierung. Die Repo-Escape-Prüfung erfolgt **nach**
   der Dekodierung ([`DC-FA-LINK-001`](lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)). Ziele mit führendem `/` werden
   relativ zur Repo-Wurzel interpretiert (Schärfung über das
   Lastenheft-Minimum „relative Ziele" hinaus).
5. **Symlink-Prüfung** ([`DC-FA-LINK-002`](lastenheft.md#dc-fa-link-002--symlink-ablehnung)):
   alle Komponenten des lexikalisch aufgelösten Zielpfads, die
   innerhalb der Repo-Wurzel liegen, werden per Lstat geprüft
   (außerhalb liegende Komponenten sind nicht prüfbar — dort greift
   `repo-escape`). Ist eine Komponente oder das Ziel selbst ein
   Symlink ⇒ Befund `symlink`, unabhängig vom Symlink-Ziel; Vorrang
   vor `repo-escape`, genau ein Befund pro Linkziel.

6. **Ortsfeste Verweise (`links.resolve-from`, opt-in).** Je konfigurierter
   **Gruppe** (wandernde Verzeichnisse `dirs`, optional ortsfeste `fixed-dirs`):
   für jede gescannte Datei, deren Verzeichnis in `dirs` liegt, wird jedes
   relative Ziel (dieselbe Extraktion und Dekodierung wie Schritte 3–4)
   zusätzlich von **jedem** Ort der Gruppe (`dirs` ∪ `fixed-dirs`) aufgelöst.
   **Vorbedingung ist der saubere Ist-Ort:** löst das Ziel schon am Ist-Ort
   nicht auf (oder meldet dort `repo-escape`/`symlink`), melden die Schritte
   4/5 — dieser Schritt schweigt, sonst trüge dieselbe Referenz zwei Befunde.
   Die Klasse dieses Schritts ist „am Ist-Ort grün“; geprüft werden die
   **übrigen** Orte. Löst das Ziel von mindestens einem Ort **nicht** auf ein
   existierendes Ziel auf, **oder** lösen zwei Orte auf **verschiedene** Ziele
   ⇒ Grund-Code `link-position-dependent`, ein Befund je Referenz an ihrer
   Zeile; die Meldung nennt den ersten nicht auflösenden Ort (Orte in
   Config-Reihenfolge) bzw. die divergierenden Ziele (sortiert). Quellen sind
   nur Dateien **unmittelbar** in einem `dirs`-Verzeichnis — Unterverzeichnisse
   wandern nicht als Einheit mit und sind ausgenommen. Dateien in `fixed-dirs`
   sind **keine** Quellen (ihre Verweise prüft nur Schritt 4). **Mindestens ein
   `dirs`-Ort jeder Gruppe muss im Baum existieren** — existiert keiner, zeigt
   die Gruppe sicher ins Leere (Tippfehler im Stamm-Pfad); ebenso ist ein Ort,
   der als **Datei** existiert, sicher falsch. Beides meldet einmal je Lauf
   über denselben Grund-Code mit dem Ort als Ziel (fail-closed). Ein
   **einzelner** fehlender Ort meldet bewusst **nicht**: git überträgt leere
   Verzeichnisse nicht — ein legitim geleertes Lifecycle-Verzeichnis fehlt auf
   jedem frischen Klon und wäre von einem Tippfehler nicht unterscheidbar
   (benannte Grenze).
   Das Ventil `ignore-refs` gilt wie in Schritt 4 (ein dort ausgenommenes Ziel
   wird auch hier übersprungen); Anker-Fragmente bleiben außen vor. **Die
   Gruppen-Orte müssen im wirksamen Scan-Bereich liegen** (`scan.roots` bzw.
   `links.scope`): eine Datei, die der Scanner nie liest, ist still keine
   Quelle — die Prüfung läuft je gescannter Datei, und diese Kopplung ist
   benannt, nicht geprüft. Ohne den Block entfällt der Schritt und der
   Befundsatz ist byte-identisch. **Exit 2** am Config-Rand bei: leerem `dirs` (< 2 Orte —
   eine Gruppe aus einem Ort prüft nichts), einem Verzeichnis absolut oder mit
   `..`-Segment, oder einem Ort, der in **mehreren** Gruppen als `dirs`-Mitglied
   auftritt (die Zuordnung einer Quelle wäre mehrdeutig).

Vor `target-missing`/`repo-escape` greift das geteilte Referenz-Ventil
([DC-FA-REF-001.a](#dc-fa-ref-001a--geteiltes-referenz-ventil-ignore-refs)): ein
aufgelöstes Ziel, das ein `ignore-refs`-Eintrag ignoriert, wird übersprungen; die
Symlink-Prüfung bleibt.

### DC-FA-ANCH-001.a — GitHub-Slug-Algorithmus

**Eingabe:** Heading-Text (ATX, `#`–`######`). **Schritte:**

1. Markdown-Inline-Auszeichnung entfernen: Code-Span-Backticks,
   Emphasis-Sterne (`*`), Links → Linktext; literale Unterstriche
   bleiben erhalten (GitHub-Verhalten, Schritt 3 erlaubt `_`).
2. Unicode-Kleinschreibung.
3. Alle Zeichen entfernen, die nicht Unicode-Buchstabe, Ziffer,
   Leerzeichen, `-` oder `_` sind (Umlaute bleiben erhalten).
4. Leerzeichen → `-` (jedes einzeln; Mehrfach-Leerzeichen ergeben
   Mehrfach-Bindestriche).
5. Duplikate: das erste Vorkommen erhält den Basis-Slug, weitere in
   Dokumentreihenfolge die Suffixe `-1`, `-2`, ….

Setext-Headings (`===`/`---`-Unterstreichung) werden in 0.x nicht
unterstützt (die Quell-Tools verwendeten ausschließlich ATX);
Fortschreibung dieser Datei genügt, falls Bedarf entsteht. Existiert
die Zieldatei eines `ziel.md#anker`-Links nicht, schweigt `anchors`
(Befund kommt von `links`,
[`DC-FA-ANCH-001`](lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors)).
Das geteilte Referenz-Ventil
([DC-FA-REF-001.a](#dc-fa-ref-001a--geteiltes-referenz-ventil-ignore-refs)) überspringt
die Anker-Prüfung eines Ziels, das ein `ignore-refs`-Eintrag ignoriert.

### DC-FA-ANCH-001.b — Inline-HTML-Anker

Zur **gültigen Anker-Menge** einer Datei zählen — zusätzlich zu den
Heading-Slugs aus
[DC-FA-ANCH-001.a](#dc-fa-anch-001a--github-slug-algorithmus) — ihre
Inline-HTML-Anker. **Eingabe:** Dateiinhalt. **Erkannt** werden,
zeilenbasiert und außerhalb von Fenced-Code-Blöcken **und
Inline-Code-Spans** (GitHub rendert HTML in Code-Auszeichnung nicht als
Sprungziel — es bleibt literaler Text; gleiche Code-Behandlung wie die
Heading-Extraktion):

1. der Wert eines `id`-Attributs an einem beliebigen HTML-Element;
2. der Wert eines `name`-Attributs an einem `<a>`-Element — der Tag-Name
   wird exakt verglichen (`a` mit Wortgrenze, kein Präfix-Treffer auf
   `<area>`, `<abbr>` o. ä.).

Attributwerte stehen in doppelten oder einfachen Anführungszeichen.
HTML-Anker werden **wörtlich** in die Anker-Menge aufgenommen (kein
Slug, keine Kleinschreibung — HTML-Fragment-Auflösung ist
case-sensitiv). Der Fragment-Vergleich (nach RFC-3986-Dekodierung)
trifft, wenn das Fragment einem Heading-Slug **oder** einem
Inline-HTML-Anker der Zieldatei entspricht. Über Zeilengrenzen
verteilte Tags werden nicht erkannt (konservativ). Eine zu groß
geratene Anker-Menge führt nur dazu, dass das Modul schweigt (nie ein
Falsch-Befund); die permissive Richtung ist daher bewusst und mit dem
Schweige-Charakter des Moduls konsistent.

### DC-FA-ID-001.a — Kennungs-Prüfung

Pro Zeile (außerhalb Fence/Inline-Code) werden die konfigurierten
Muster in Deklarationsreihenfolge gematcht; das erste passende Muster
gewinnt pro Vorkommen (auch bei überlappenden Treffern: ein von einem
früheren Muster beanspruchter Textbereich wird von späteren Mustern
nicht erneut gematcht). Ein Vorkommen gilt als verlinkt, wenn es
innerhalb des Linktexts eines Markdown-Links (keine Bildreferenz)
liegt. Vorkommen innerhalb der Link-Syntax außerhalb des Linktexts
(Ziel-Klammer) sowie innerhalb von Bildreferenzen (Alt-Text und Ziel)
sind kein Fließtext und damit linkpflichtfrei. Ebenfalls
linkpflichtfrei sind ATX-Heading-Zeilen (Headings sind Struktur-,
kein Fließtext — Definitions- und Schärfungs-Headings tragen ihre
Kennung nackt) sowie alle Vorkommen innerhalb des deklarierten
`target` des matchenden Musters (Definitions-Ort: die Target-Datei
selbst bzw. alle Dateien unterhalb eines Target-Verzeichnisses — eine
Definition muss nicht auf sich selbst verlinken). Ebenso
linkpflichtfrei ist ein nacktes Vorkommen, dessen Datei ein Glob aus
`ids.patterns[].exempt-paths` matcht oder dessen Zeile den Marker
`d-check:ignore` trägt — die beiden unten unter *Link-Politik*
beschriebenen Ventile sind ein Ganzdatei- bzw. Ganzzeilen-Carve-out und
wirken auf nackte Fließtext-Vorkommen genau wie auf
Inline-Code-Vorkommen (unabhängig von der `link-policy`). Alle übrigen
Vorkommen erzeugen den Befund `id-unlinked`. Da die Extraktion zeilenbasiert
ist ([§DC-FA-LINK-001.a](#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
Schritt 3), gilt eine Kennung in mehrzeiligem Linktext als nackt —
linkpflichtige Kennungen gehören in einzeilige Links. Das deklarierte
`ids.patterns[].target` (Datei oder Verzeichnis, relativ zur
Repo-Wurzel) muss existieren und innerhalb der Repo-Wurzel liegen —
Verletzung ist ein Konfigurationsfehler (Exit 2, [§2](#d-checkyml)).

**Link-Politik (`ids.patterns[].link-policy`).** Default `prose` = das
oben beschriebene Verhalten (Inline-Code-Vorkommen sind frei, weil die
Vorverarbeitung sie leert). Bei `always` wird das Muster **zusätzlich**
gegen die Werte der Inline-Code-Spans der rohen Prosa-Zeilen geprüft —
dieselbe fence-aware Span-Erkennung wie
[§DC-FA-CODE-001.a](#dc-fa-code-001a--pfade-in-inline-code) Schritt 2.
Ein Inline-Code-Vorkommen ist linkpflichtfrei genau dann, wenn eine der
Bedingungen gilt: (1) der Code-Span ist der Linktext eines Markdown-Links
(`` [`…`](ziel) ``); (2) das Vorkommen liegt im `target` des Musters;
(3) die geprüfte Datei matcht ein Glob aus `ids.patterns[].exempt-paths`;
(4) die Zeile trägt den Marker `d-check:ignore` (HTML-Kommentar — ab
dieser Anforderung wirkt er auf `codepaths` **und** `ids`); (5) es ist
eine ATX-Heading-Zeile. Alle übrigen Inline-Code-Vorkommen erzeugen
`id-unlinked` (kein neuer Grund-Code). Die Ventile (3) `exempt-paths`
und (4) `d-check:ignore` sind dieselben wie im Prosa-Basisalgorithmus
oben — sie nehmen Datei bzw. Zeile vollständig aus, für nackte wie für
Inline-Code-Vorkommen. Muster-Präzedenz bleibt unverändert; der Wechsel
`prose`→`always` lässt den Prosa-Befundsatz unverändert und ist rein
additiv (ein Muster mit `always` findet eine Obermenge seiner
`prose`-Befunde).
`exempt-paths` nutzt Glob-Syntax relativ zur Repo-Wurzel wie
`scan.ignore`; die Glob-Auswertung ist reihenfolgestabil
([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)).

### DC-FA-MTX-001.a — Klassen- und Status-Auflösung

1. **Klassenzuordnung:** Glob-Muster der `classes` in
   Deklarationsreihenfolge; die erste passende Klasse gilt. Dateien
   ohne Klasse nehmen nicht an der Matrix-Prüfung teil.
2. **Status-Extraktion** in fester Reihenfolge: (1) erste Zeile, die
   mit `**Status:**` beginnt; (2) sonst erste nicht-leere Textzeile
   unter einem `Status`-Heading (beliebige Ebene, Heading-Vergleich
   case-insensitiv). Sind beide Formen vorhanden, zählt die
   `**Status:**`-Zeile. Beide Formen lesen nur Zeilen außerhalb von
   Fenced-Code-Blöcken (Fence-Inhalt ist kein Statuswert); Status wird
   nur aus Markdown-Zieldateien extrahiert, andere Ziele gelten als
   aktiv. Der Wert wird case-insensitiv als Präfix-Match
   gegen `status.forbidden` verglichen (so matcht
   `Superseded by ADR-0007` den Wert `superseded`). Ohne Status-Feld <!-- d-check:ignore (Beispiel-Statuswert, fiktiv) -->
   gilt das Dokument als aktiv. Regel- und Status-Prüfung sind
   unabhängig: ein Link kann `matrix-forbidden` **und**
   `matrix-inactive` zugleich erzeugen (zwei verschiedene
   Verletzungen, zwei Befunde).
3. **Sektions-Ausnahme:** Links innerhalb von Sektionen, deren
   Heading-Text (getrimmt, ohne Markdown-Auszeichnung, case-sensitiv)
   in `exclude-sections` steht (z. B. „Historie"),
   werden von `matrix` nicht geprüft (Provenance-Ausnahme gemäß
   [`DC-FA-MTX-001`](lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
   Out-of-Scope).
4. **Supersede-Lineage-Ausnahme** (`allow-supersede-lineage`, Default
   `false`): Ist das Flag gesetzt, wird vor dem Emittieren eines
   `matrix-inactive`-Befundes geprüft, ob die Kante eine deklarierte
   Ablösung ist. Aus der Quelldatei werden einmalig die Werte aller
   Felder aus `supersede-fields` gewonnen; eine Feld-Zeile hat die Form
   `**Feld:** Wert` oder `Feld: Wert` (Feldname case-insensitiv, nur
   Zeilen außerhalb von Fenced-Code). Die Kante X → Y ist ausgenommen,
   wenn ein solcher Feldwert — nach Normalisierung (Kleinschreibung,
   Whitespace auf einzelne Leerzeichen kollabiert) — den Linktext der
   Referenz (falls nicht leer) oder den aufgelösten Zielpfad (relativ
   zur Repo-Wurzel, dessen Basename oder Basename ohne Endung) als
   Teilzeichenkette enthält. Die Ausnahme wirkt ausschließlich auf
   `matrix-inactive`; die Klassen-Regelprüfung (`matrix-forbidden`,
   Schritt 1) bleibt unberührt, und nur die deklarierte Lineage-Kante
   ist betroffen — andere Quellen ohne passendes Feld melden Y weiter
   als inaktiv. Ohne gesetztes Flag wird `supersede-fields` nicht
   konsultiert; der Befundsatz ist byte-identisch
   ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)). `matrix`
   trägt — anders als `codepaths`/`ids` — keinen `d-check:ignore`-Marker
   ([§DC-FA-CODE-001.a](#dc-fa-code-001a--pfade-in-inline-code) Schritt 1
   gilt nur für jene Module).
5. **Klasseninterne Verweisrichtung**
   ([`DC-FA-MTX-002`](lastenheft.md#dc-fa-mtx-002--verweisrichtung-innerhalb-einer-geordneten-dokumentklasse-modul-matrix)):
   Trägt die in Schritt 1 zugeordnete Quell-Klasse ein nicht-leeres `order`
   (Liste von Pfad-Globs) und `direction: no-downward`, wird — zusätzlich zur
   Klassen-Regelprüfung — geprüft, ob die Zieldatei **derselben** Klasse
   angehört. Ist das so, ist der **Rang** je Datei der Index des ersten
   `order`-Globs, den sie matcht (First-Match wie die Klassenzuordnung; eine
   Datei ohne Treffer ist rangfrei). Haben Quelle und Ziel je einen Rang und
   ist der Quell-Rang **kleiner** als der Ziel-Rang (Verweis von der
   autoritativeren auf die weniger autoritative Schicht, auch über mehrere
   Stufen), entsteht ein `matrix-downward`-Befund (Datei, Zeile, Ziel; die
   Meldung nennt beide Ränge). Rangfreie Dateien und klassenübergreifende
   Referenzen lösen kein `matrix-downward` aus; die Prüfung ist unabhängig von
   `matrix-forbidden`/`matrix-inactive`. Fehlkonfiguration ist fail-closed
   (Config-Adapter, Exit 2): unbekannter `direction`-Wert, `direction` ohne
   `order`, `order` ohne `direction`. Ohne beide Felder ist der Befundsatz
   byte-identisch ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)).
6. **Token-Referenzen und Grandfathering**
   ([`DC-FA-MTX-003`](lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)):
   Trägt eine Klasse ein `token`-Regex, scannt `matrix` zusätzlich zu den Links
   den **Prosa-Körper** (außerhalb Fenced-Code, außerhalb `exclude-sections` und
   ohne Markdown-Link-Spans — die deckt die Link-Prüfung ab) nach `token`
   *anderer* Klassen. Ein Treffer der Klasse B im Dokument der Klasse A ist eine
   Token-Referenz A → B; eine verbotene Kante (`rules`) erzeugt `matrix-forbidden`
   (Token-Form; Datei, Zeile, Token; Meldung nennt beide Klassen) — **es sei
   denn**, auf derselben **rohen** Zeile steht der Provenance-Marker
   `<!-- d-check:status-provenance -->` (Match unabhängig vom Kommentar-Whitespace),
   der die verbotene Token-Referenz als deklarierte Provenance ausnimmt. Token in
   Markdown-Links und in Fenced-Code zählen nicht. Eine Datei, die ein
   `exempt-paths`-Glob matcht, wird von `matrix` **ganz** übersprungen
   (Grandfathering immutabler Dokumente). Fehlkonfiguration ist fail-closed
   (Config-Adapter, Exit 2): `token` kompiliert nicht oder matcht den Leerstring.
   Ohne `token`/`exempt-paths` ist der Befundsatz byte-identisch
   ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)).

### DC-FA-EXT-001.a — Externe Erreichbarkeit

HEAD-Request; bei HTTP 405/501 Fallback auf GET (Body verworfen).
Das Timeout gilt **pro Request** — im Fallback-Fall (HEAD → GET)
können bis zu zwei Requests anfallen. Der Schema-Vergleich
(`http:`/`https:`) ist case-insensitiv (RFC 3986). Der Fragment-Teil
wird vor Prüfung und Dedupe entfernt (Fragmente werden nie
übertragen); der Befund nennt das Original-Linkziel.
Redirects bis `REDIRECT_MAX` Stationen gefolgt. Pro URL genau eine
Prüfung pro Lauf (Dedupe — der Befund erscheint an jedem Vorkommen),
begrenzte Parallelität (`EXTERNAL_PARALLEL`). Ergebnis-Auswertung:
Status < 400 → kein Befund; ≥ 400 → `external-status`; Timeout →
`external-timeout`; Redirect-Kette > `REDIRECT_MAX` →
`external-redirects`; Transportfehler (DNS-/Verbindungsfehler) →
`external-status` (Status 0, Grund in der Meldung).

### DC-FA-CODE-001.a — Pfade in Inline-Code

Arbeitet auf den **rohen Prosa-Zeilen** (fence-aware) — die übrige
Vorverarbeitung entfernt Inline-Code gerade.

**Datei-Ventil (vorab):** Dateien, deren Pfad ein Glob aus
`codepaths.exempt-paths` matcht (Syntax wie `scan.ignore`, relativ zur
Repo-Wurzel), werden **ganz** übersprungen — datei-weit, unabhängig von
`codepaths.roots`; dasselbe Ventil wie
[DC-FA-ID-001.a](#dc-fa-id-001a--kennungs-prüfung), komplementär zum
zeilenweisen `d-check:ignore` (Schritt 1). Ohne gesetztes `exempt-paths`
byte-identisch.

**Referenz-Ventil (`ignore-refs`):** `codepaths` honoriert das **geteilte
Referenz-Ventil** [DC-FA-REF-001.a](#dc-fa-ref-001a--geteiltes-referenz-ventil-ignore-refs):
ein in Schritt 5 aufgelöster Ziel-Pfad, den ein `ignore-refs`-Eintrag ignoriert
(Quell-Skopus `in` ∧ `refs` ∧ ¬`keep`), wird **nicht** existenz-, escape-, anker- oder
ortsfestigkeits-geprüft (kein `target-missing`, `codepath-missing`,
`repo-escape`, `anchor-missing` oder `link-position-dependent`) —
**referenz-weit**, unabhängig von Datei und Zeile (anders als das datei-weite
`exempt-paths` und das zeilenweise `d-check:ignore`). Es unterdrückt nur die Prüfung
*dieses* Pfads, keine anderen Befunde. Die modul-lokale Liste `codepaths.ignore-refs`
bleibt als **Alias** gültig (wie ein Eintrag ohne `in`/`keep`, skopiert auf
`codepaths`) — byte-identisch zur bisherigen Fassung. Ohne jeden Eintrag
byte-identisch.

**Schritte:**

1. Zeilen mit dem Marker `d-check:ignore` (HTML-Kommentar, Begründung
   in Klammern empfohlen) werden übersprungen — der Marker wirkt
   ausschließlich auf dieses Modul. ATX-Heading-Zeilen werden ebenso
   übersprungen: Titel sind keine Prosa-Referenzen (gleiche Ausnahme
   wie [DC-FA-ID-001.a](#dc-fa-id-001a--kennungs-prüfung); ein Marker
   im Heading würde zudem dessen Anker-Slug verändern).
2. Pro Zeile alle Inline-Code-Spans extrahieren (CommonMark,
   Multi-Backtick-fähig — dieselbe Span-Erkennung wie das
   positionserhaltende Stripping der übrigen Module). Spans, die der
   Linktext eines Markdown-Links sind (`` [`…`](ziel) ``), sind
   ausgenommen — deren Ziel prüft das Modul `links`, der Text ist
   Beschriftung.
3. Span-Wert normalisieren (iterativ bis stabil): Whitespace trimmen,
   umschließende einfache/doppelte Anführungszeichen und schließende Satzzeichen
   (`.,;:`) entfernen. **Zeilen-Suffix:** ohne `check-lines` wird — wie bisher — nur
   ein einzelnes `:<zahl>` abgetrennt; ein Bereich `:<von>-<bis>` bleibt dabei
   **unberührt** (byte-identisch zum heutigen Verhalten). Bei aktivem `check-lines`
   werden `:<von>` **und** `:<von>-<bis>` als Zeilen-Referenz erkannt, vom Pfad
   abgetrennt und der Bereich (`<von>`, `<bis>`; ohne `-` ist `<bis> = <von>`) für
   Schritt 6 **gemerkt**.
4. Konservative Pfad-Erkennung — **kein** Pfad ist ein Wert, der leer
   ist, Whitespace oder Platzhalter-/Glob-Zeichen (`{}<>|*?=`)
   enthält, Ellipsen/Pfeile (`…`, `->`, `→`) enthält, mit `//` oder
   `#` beginnt oder ein externes Schema trägt. **Pfad** ist, was mit
   `./` oder `../` beginnt (Datei-relativ) oder mit einem der
   konfigurierten Präfixe aus `codepaths.roots` (Wurzel-relativ;
   Vergleich gegen `präfix/`).
5. Auflösung wie im Modul `links` (inkl. RFC-3986-Dekodierung):
   Fragment abtrennen; ignoriert das geteilte Referenz-Ventil (`ignore-refs`
   inkl. Alias `codepaths.ignore-refs`) den aufgelösten Pfad, wird er übersprungen
   (kein Befund, s. o. Referenz-Ventil); sonst: Escape → `repo-escape`; fehlendes
   Ziel → `codepath-missing`. Trägt der Wert ein Fragment und ist das Ziel
   eine Markdown-Datei, wird der Anker gegen die gültige Anker-Menge
   der Zieldatei geprüft (Heading-Slugs und Inline-HTML-Anker; Verfahren
   und Cache wie
   [DC-FA-ANCH-001.a](#dc-fa-anch-001a--github-slug-algorithmus) /
   [DC-FA-ANCH-001.b](#dc-fa-anch-001b--inline-html-anker);
   Treffer fehlt → `anchor-missing`). Nicht lesbare Ziele: das Modul
   schweigt zum Anker (Existenz wurde bereits geprüft).
6. **Zeilen-Check (nur bei `check-lines`):** trägt der Wert einen in
   Schritt 3 gemerkten Zeilen-Bereich und existiert das aufgelöste,
   nicht-ignorierte Ziel (Schritt 5 lieferte keinen `codepath-missing`),
   wird der Bereich geprüft: `von < 1` oder `von > bis` ⇒ `citation-inverted-range`
   (die Zeilennummern sind 1-basiert); sonst
   hat die Zieldatei weniger als `bis` Zeilen ⇒ `citation-out-of-range`. Ohne
   `check-lines` entfällt der Schritt (byte-identisch); ein fehlendes Ziel
   bleibt `codepath-missing` (der Zeilen-Check setzt Existenz voraus, nicht
   umgekehrt). Der Wert ohne Zeilen-Bereich ist unverändert.

### DC-FA-CITE-001.a — Verbatim-Zitat-Verifikation (`citations`)

Das Modul `citations` prüft **nur** per Direktive ausgezeichnete Zitate — kein
Prosa- oder Voll-Scan. Arbeitet auf den rohen Zeilen (fence-aware wie die übrigen
Module).

**Schritte:**

1. Zeilen mit dem Marker `<!-- d-check:cite <pfad>:<von>-<bis> -->` finden
   (HTML-Kommentar; `<pfad>` Datei-/Wurzel-relativ wie in
   [DC-FA-CODE-001.a](#dc-fa-code-001a--pfade-in-inline-code), `<von>`/`<bis>`
   **1-basierte** Zeilennummern, `<bis>` optional = `<von>`). Fehlt die Direktive,
   prüft das Modul nichts. Ein **malformter** Span (nicht-numerisch, fehlend) ⇒
   **fail-closed** (Exit 2).
2. Der **Zitattext** ist das der Direktive folgende Zitat. **Leerzeilen trennen
   nicht** — gesucht wird der nächste nicht-leere Kandidat, auch über mehrere
   Leerzeilen hinweg. **Ein Fenced-Block dagegen trennt in beiden Zweigen**: er
   ist eine Absatzgrenze
   ([DC-FA-LINK-001.a](#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
   Schritt 2), und liegt einer zwischen Direktive und Kandidat, folgt der
   Direktive **kein** Zitat — der fail-closed-Fall unten greift. Ist die **nächste nicht-leere Zeile** ein
   `>`-Blockquote, gilt dieser Block (zusammenhängende `>`-Zeilen, jeweils nach
   Abtrennen von `> ` bzw. `>` ohne Leerzeichen; eine Leer-, eine Nicht-`>`-Zeile
   **oder ein Fenced-Block** beendet ihn). Andernfalls der **nächste
   inline-Zitat-Span im selben Absatz** nach der Direktive — er darf Prosa vor
   sich haben und über mehrere Zeilen laufen, begrenzt durch ein
   Anführungs-Paar: `„` öffnet + `"` schließt, oder das erste `"` öffnet + das
   nächste `"` schließt. Findet sich weder ein `>`-Block (als nächste nicht-leere
   Zeile) noch ein schließendes Anführungspaar im Absatz ⇒ **fail-closed**
   (Exit 2) — die Direktive ist unbrauchbar (Autoren-Fehler, kein Schweigen).
3. Auflösung von `<pfad>` wie im Modul `links`; verlässt das Ziel die Wurzel ⇒ Befund
   `repo-escape` (Exit 1, dieselbe Sicherheits-Prüfung wie `codepaths`/`links`).
   Danach die **Zitat-Fäule** (Befund, **nicht** fail-closed — kohärent zum
   Zeilen-Check aus [DC-FA-CODE-001.a](#dc-fa-code-001a--pfade-in-inline-code)
   Schritt 6): `von < 1` oder `von > bis` ⇒ `citation-inverted-range` (die Zeilennummern
   sind 1-basiert); sonst fehlt die Zieldatei oder
   hat sie weniger als `bis` Zeilen ⇒ `citation-out-of-range` (kein Vergleich, da die
   Spanne nicht existiert).
4. **Normalisierung.** Quell-Spanne (Zeilen `von`–`bis` verbunden) und Zitattext
   werden je **whitespace-normalisiert**: jeder Lauf aus Leerzeichen/Tab/Zeilenumbruch
   wird zu **einem** Leerzeichen, führend/schließend getrimmt. Sonst **keine**
   Normalisierung (Markdown-Auszeichnung, Satzzeichen, Groß-/Kleinschreibung bleiben).
   **Mindestlänge:** ist der normalisierte Zitattext kürzer als **16 Zeichen**, wird
   er **nicht** geprüft (kein Befund) — ein sehr kurzer Teilstring träfe zufällig
   (schwache Diskriminierung); dokumentierter Trade-off, keine Falsch-Rot-Gefahr.
5. **Teilstring-Vergleich.** Ist der normalisierte Zitattext ein **zusammenhängender
   Teilstring** der normalisierten Quell-Spanne ⇒ kein Befund; sonst `citation-mismatch`
   (Datei = prüfende Datei, Zeile = Direktiven-Zeile, `target` = `<pfad>:<von>-<bis>`).
   So bleibt der Vergleich **zeichengenau auf dem Inhalt**, tolerant nur gegenüber
   Umbruch/Whitespace — ein inline, re-wrapped, mitten in der Zeile beginnendes/endendes
   Zitat besteht; jede echte Wort-Abweichung bricht den Teilstring.

Hermetisch (nur Datei-Lesen über den Filesystem-Port, kein git/Netz —
[`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit));
strikt opt-in, ohne aktives `citations`-Modul byte-identisch
([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)).

### DC-FA-SPAN-001.a — Span-Artefakt-Erkennung

1. **`span-unclosed`:** absatzweise, mit derselben Absatz-Semantik
   wie die Vorverarbeitung
   ([§DC-FA-LINK-001.a](#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
   Schritt 2: Grenzen sind Leerzeile und Fence). Im Absatz werden
   Backtick-Folgen von links nach rechts gepaart (die öffnende Folge
   bestimmt die gleich lange schließende; eine ungeschlossene Folge
   ist literal, die Suche läuft dahinter weiter). Für jede
   ungeschlossene Folge, auf die **unmittelbar ein
   Nicht-Whitespace-Zeichen** folgt (kein Leerzeichen/Tab, kein
   Zeilenumbruch, kein Absatz-Ende), entsteht ein Befund an der
   Zeile der Folge; das gemeldete Ziel ist die Backtick-Folge samt
   der unmittelbar folgenden Nicht-Whitespace-Zeichen, gekappt auf
   30 Zeichen. Ungeschlossene Folgen mit Whitespace dahinter sind
   beabsichtigt literal — kein Befund (konservative Erkennung,
   [`DC-FA-SPAN-001`](lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)).
2. **`span-nested-link`:** auf den vorverarbeiteten Zeilen (Fences
   entfernt, Code-Spans geleert) wird jedes Vorkommen des Musters
   „Linktext-Schließung mit Ziel, unmittelbar gefolgt von einer
   weiteren Linktext-Schließung mit Ziel-Öffnung" gemeldet —
   lexikalisch: `](` … `)` direkt gefolgt von `](`. Benachbarte
   eigenständige Links sind kein Treffer (zwischen ihnen steht die
   öffnende Linktext-Klammer), und **Bildreferenzen als Linktext**
   (`[![…](…)](…)` — das Badge-Muster) sind legales Markdown und
   ebenfalls kein Treffer (Kalibrierungs-Befund: vendorte
   Paket-READMEs mit Shields-Badges). Das gemeldete Ziel ist der
   Treffer, gekappt auf 40 Zeichen.

3. **`fence-unclosed`:** dateiweit, **nicht** absatzweise. Jede Zeile
   wird — mit derselben Trimmung wie die Vorverarbeitung (Space und Tab,
   **nicht** unicode-weit) — durch **beide** Schluss-Lesarten geführt,
   die das Produkt kennt:
   - **naiv:** jede Fence-Zeile (`FenceToggle`) kippt den Zustand. Das
     ist die Lesart der Vorverarbeitung
     ([§DC-FA-LINK-001.a](#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
     Schritt 1) und damit aller Module, die auf ihr aufsetzen.
   - **streng:** geschlossen wird nur zeichen- und längenabgeglichen
     (gleiches Fence-Zeichen, mindestens gleiche Länge, dahinter nur
     Whitespace); jede andere Fence-Zeile innerhalb eines offenen Blocks
     ist Inhalt. Das ist die Lesart des Tabellen-Lesers aus
     [§DC-FA-REQ-001.a](#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen).

   Steht **mindestens eine** der beiden am Dateiende noch innerhalb eines
   Blocks, entsteht **ein** Befund — denn genau dann überspringt mindestens
   ein Modul den Rest der Datei. Ist die strenge Lesart offen, steht der
   Befund an **ihrer** Öffnungszeile, sonst an der zuletzt öffnend
   gewerteten. Beides ist eine **Fundstelle**: welche Öffnung fehlt, ist
   grundsätzlich nicht entscheidbar — gleich lange Öffner sind
   ununterscheidbar, und auch die strenge Lesart zeigt daneben, wenn eine
   längere Fence-Zeile eine kürzere Öffnung geschlossen hat und erst eine
   spätere Öffnung offen bleibt. Das gemeldete Ziel ist diese Zeile, links
   wie rechts getrimmt (Space, Tab und das CR einer CRLF-Zeile) und auf
   30 **Runen** gekappt.

   **Bewusst nicht gelöst:** die Paarung selbst. Welche der beiden
   Lesarten die richtige ist, bleibt offen — die Vorverarbeitung nach
   [§DC-FA-LINK-001.a](#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
   Schritt 1 ändert sich **nicht**. Gemeldet wird der **Zustand**: dass
   eine Lesart offen endet, unabhängig davon, welche recht hat.
   Genau **ein** Befund je Datei.

Die drei Prüfungen kennen keinen Opt-out-Marker; das Modul akzeptiert
den generischen Schlüssel `spans.scope`
([§DC-FA-CONF-002.a](#dc-fa-conf-002a--effektiver-scan-scope-pro-modul)),
weitere Konfiguration existiert nicht.

### DC-FA-HOST-001.a — Host-Pfad-Erkennung

1. **Geltungsbereich:** rohe Prosa-Zeilen außerhalb von
   Fenced-Code-Blöcken, **einschließlich Inline-Code** (dort leben
   die Leaks typischerweise — bewusste Abweichung von den
   Strip-Konsumenten; Beispiel-Inhalte gehören in Fences,
   [`DC-FA-HOST-001`](lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in)).
2. **Muster** (jeweils mit Wortgrenzen-Vorbedingung — dem Treffer
   darf kein Buchstabe, keine Ziffer und keines der Zeichen
   Unterstrich, Punkt, Doppelpunkt, Schrägstrich oder Bindestrich
   unmittelbar vorausgehen; URL-Pfade hinter Schemata sind damit
   ausgenommen):
   - **Unix:** Schrägstrich + Präfix-Verzeichnisname + Schrägstrich
     + Restpfad (Zeichen bis Whitespace bzw. `<`, `>`, `)`, `]`,
     Anführungszeichen oder Backtick — Inline-Code-Delimiter gehören
     nicht zum Pfad). Die Präfixliste ist konfigurierbar
     (`hostpaths.prefixes`, **ersetzt** den Default); Default-Namen:
     Development, home, Users, Volumes, mnt, media (tmp bewusst
     nicht — Lastenheft 0.7.2).
   - **Windows-Laufwerk:** Buchstabe + Doppelpunkt + Backslash +
     Restpfad (Vorbedingung hier: kein Wortzeichen davor).
   - **UNC:** doppelter Backslash + Servername + Backslash +
     Restpfad (Vorbedingung wie Windows-Laufwerk; der Servername
     beginnt mit einem alphanumerischen Zeichen — Regex-Beispiele
     wie escapte Punkt-Folgen matchen damit nicht).
3. **Normalisierung:** schließende Satzzeichen (`.`, `,`, `;`, `:`)
   werden vom gemeldeten Pfad abgetrennt (wie
   `codepaths`-Normalisierung). Befund `hostpath-forbidden` mit
   Datei, Zeile und dem normalisierten Pfad; kein Opt-out-Marker,
   das Modul akzeptiert den generischen Schlüssel `hostpaths.scope`
   ([§DC-FA-CONF-002.a](#dc-fa-conf-002a--effektiver-scan-scope-pro-modul)).

**Bekannte Grenze:** ein Repo-Wurzel-relatives Linkziel, dessen
erstes Segment zufällig ein Präfix-Name ist (etwa ein Verzeichnis
`home/` im Repo), wird gemeldet — Ausweg ist die konfigurierbare
Präfixliste; konservativer Default vor Vollständigkeit.

### DC-FA-CONF-002.a — Effektiver Scan-Scope pro Modul

1. **Auflösung:** Für jedes aktive Modul gilt der globale Scan-Scope
   (`scan.roots`/`scan.ignore`), außer das Modul deklariert
   `<modul>.scope` — dann **ersetzt** dieser den globalen Scope für
   genau dieses Modul (eigener Discover-Lauf; ein Modul-Scope kann
   Dateien umfassen, die der globale Scan nicht enthält). Innerhalb
   von `scope` ist `roots` Pflicht (fehlend = Konfigurationsfehler,
   Exit 2 — keine stille Vererbung), `ignore` ist optional. Es gelten
   unverändert die Scan-Regeln aus
   [`DC-FA-SCAN-001`](lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln):
   deklarierte Wurzeln müssen existieren und innerhalb der Repo-Wurzel
   liegen (Exit 2), `"."` steht für die gesamte Repo-Wurzel,
   Ignorier-Muster prunen den Verzeichnis-Abstieg, die `SKIP_DIRS`
   gelten immer, eine explizit leere `roots`-Liste prüft nichts.
2. **Lauf:** Geprüft wird die **Vereinigungsmenge** aller effektiven
   Modul-Scopes in deterministischer Reihenfolge (Pfade bytewise
   sortiert); jede Datei wird genau einmal gelesen und vorverarbeitet,
   jedes Modul prüft nur Dateien seines effektiven Scopes. Die
   Zusammenfassung („N Datei(en) geprüft") und das `files`-Feld der
   JSON-Ausgabe zählen die Vereinigungsmenge. Die Befund-Sortierung
   bleibt global ([DC-QA-02.a](#dc-qa-02a--determinismus-und-sortierung)).
3. **Abwärtskompatibilität:** Ohne `scope`-Schlüssel ist das Verhalten
   byte-identisch zum Lauf vor dieser Anforderung (ein globaler
   Discover-Lauf, alle aktiven Module auf allen Dateien).

### DC-QA-01.a — Benchmark

**Fixture** (deterministisch generiert): 1.000 Markdown-Dateien unter
`docs/`, je ein H1- und zehn H2-Headings; pro Abschnitt ein
Datei-Querverweis auf die zyklisch nächste Datei und ein Anker-Link
auf deren gleichnamigen Abschnitt, plus Fülltext; Gesamtgröße ≤ 20 MB.
**Messprotokoll:** Default-Module (`links`, `anchors`) ohne
Konfigurationsdatei, Runtime-Container mit read-only-Mount, ohne
Netz und auf **2 vCPU begrenzt** (`--cpus 2` — die
Hardware-Normierung aus
[`DC-QA-01`](lastenheft.md#dc-qa-01--performance)), N ≥ 3 Läufe
(ungerade); der **Median** zählt (inklusive Container-Start).
**Pass-Kriterium:** Median < 5 s
([`DC-QA-01`](lastenheft.md#dc-qa-01--performance)).

### DC-QA-02.a — Determinismus und Sortierung

Befunde werden nach vollständiger Sammlung stabil sortiert:
(1) Datei-Pfad bytewise aufsteigend, (2) Zeile aufsteigend,
(3) Regelmodul-Name, (4) Ziel, (5) Grund-Code. Interne Parallelität
ist erlaubt, darf aber die Ausgabe nicht beeinflussen. Das Modul
`external` ist von der Byte-Identitäts-Garantie ausgenommen, soweit
Server-Antworten variieren (Netz-Nichtdeterminismus); Sortierung gilt
auch dort.

### DC-FA-DIAG-001.a — Kennungs-Konsistenz in Diagramm-Fences (`diagrams`)

Das Modul `diagrams`
([`DC-FA-DIAG-001`](lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in))
ist opt-in und arbeitet **gegenläufig** zur Vorverarbeitung
([§DC-FA-LINK-001.a](#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)),
die Fences für alle übrigen Module entfernt:

1. **Fence-Öffnung.** Ein eigener Fence-Automat sammelt die Roh-Zeilen
   (1-basiert) innerhalb von Fenced-Blöcken, deren Info-String (erstes Token der
   Öffner-Zeile nach den Fence-Zeichen, kleingeschrieben) in `diagrams.fences`
   steht (Default `mermaid`). Alle übrigen Fences und der Prosa-Text bleiben außen
   vor; die gemeinsame Vorverarbeitung der anderen Module bleibt unverändert
   (keine `spans`-Interaktion).
2. **Token-Extraktion.** Je Muster (`diagrams.patterns`, Gestalt wie
   `ids.patterns`; Deklarationsreihenfolge = Präzedenz, überlappende Vorkommen
   gehören dem früheren Muster) werden die Kennungs-Token der Fence-Zeile per
   Regex gefunden. **Kein** Mermaid-/Grammatik-Parsing.
3. **Existenz-Auflösung.** Eine Kennung gilt als **definiert**, wenn sie in der
   `defined-in`-Datei des Musters als Token derselben Regex **außerhalb von
   Fences** vorkommt (die Quelle wird dafür mit entfernten Fences gelesen —
   Heading wie Tabelle/Fließtext zählen, bewusst nicht heading-zentriert wie
   `ids`). Token ohne Definition ⇒ Grund-Code `diagram-id-undefined` (Datei, Zeile
   im Fence, Kennung). Die Token-Menge je (`defined-in`, Regex) wird gecacht
   (einmaliges Lesen).
4. **Validierung beim Lauf-Start.** Jede `defined-in`-Datei muss existieren und
   innerhalb der Repo-Wurzel liegen (sonst Exit 2 ohne Prüfung, analog
   ids-Targets). **Link-Policy gilt nicht** (in Fences kein Markdown-Link
   möglich). Ohne `diagrams`-Block ist der Befundsatz byte-identisch
   ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)).

### DC-FA-VER-001.a — Versions-Pin-Konsistenz (`versions`)

Das Modul `versions`
([`DC-FA-VER-001`](lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in))
ist opt-in und arbeitet — wie `diagrams` — **gegenläufig** zur Vorverarbeitung
([§DC-FA-LINK-001.a](#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)),
die Fences für alle übrigen Module entfernt:

1. **Aktuelle Version auflösen (`current-from`).** Default `version.md#aktuell`:
   Datei links vom `#`, Anker rechts. Der adressierte Span ist die
   Heading-Section des Ankers (Heading bis zur nächsten gleich-/höherrangigen
   Überschrift; bei einem HTML-Anker ab dessen Zeile) bzw. die ganze Datei ohne
   Anker. **Welcher Anker das ist, entscheidet die Anker-Antwort aus
   [DC-FA-ANCH-001.a](#dc-fa-anch-001a--github-slug-algorithmus) und
   [DC-FA-ANCH-001.b](#dc-fa-anch-001b--inline-html-anker) — vollständig:**
   Heading-Slug **einschließlich Duplikat-Suffix** (`-1`, `-2`, …), das Fragment
   **prozent-dekodiert**, Inline-HTML-Anker nur außerhalb von Fenced- und
   Inline-Code, und der Vergleich **wörtlich** (case-sensitiv). Ein Anker in
   einem Code-Beispiel löst nicht auf; sonst gäbe derselbe Lauf zwei Antworten
   auf dieselbe Frage. Der adressierte **Span** bleibt davon unberührt und wird
   roh gelesen (Schritt 2). Aus dem Span wird das **erste** Vorkommen eines Versions-Teilmusters
   `v?\d+\.\d+\.\d+` als *erwartete* Version extrahiert. Existiert die Datei
   nicht, löst der Anker nicht auf, oder trägt der Span keine erkennbare Version
   ⇒ **Exit 2** ohne Prüfung (fail-closed, analog zur Target-Validierung von
   `ids`/`diagrams`), keine stille Grün-Meldung.
2. **Fence-Öffnung + Pin-Scan.** Ein eigener Fence-Automat liefert **alle**
   Roh-Zeilen (1-basiert) — Prosa **und** Fenced-Code (Versions-Pins leben in
   Kommando-Beispielen); die gemeinsame Vorverarbeitung der anderen Module bleibt
   unverändert (keine `spans`-Interaktion). Je Zeile werden die Vorkommen des
   konfigurierten `versions.pin-pattern` per Regex gesucht (Beispiel
   `ghcr\.io/[^\s:]+:(v?\d+\.\d+\.\d+)`); die Version steht in Capture-Gruppe 1,
   fehlt eine Gruppe, zählt der ganze Treffer.
3. **Vergleich.** Weicht die im Pin gefundene Version von der erwarteten ab ⇒
   Grund-Code `version-stale` (Datei, Zeile, `target` = gefundener Pin-Wert;
   `message` nennt erwartet vs. gefunden). Gleichheit ⇒ kein Befund.
4. **Ventile.** Wie `ids`/`codepaths`: die Glob-Liste `versions.exempt-paths`
   (datei-weit — historische Pins in Planning-`done/`-Slices, `CHANGELOG.md`,
   Lastenheft-Historie) und der Zeilen-Marker `d-check:ignore` nehmen Vorkommen
   aus (nackt wie in Code, in Fence wie außerhalb). Die `current-from`-Datei ist
   von der Pin-Prüfung **selbst ausgenommen** (ihr Verlauf listet bewusst alle
   Versionen).
5. **Read-only/Determinismus.** Ohne `versions`-Block ist der Befundsatz
   byte-identisch ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)); es wird
   nichts geschrieben ([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
   `version-stale` ist diagnose-only — kein `--repair`-Hunk in dieser Version.

### DC-FA-PIN-001.a — Content-Pin gegen inhaltlichen Drift (`pins`)

Das Modul `pins`
([`DC-FA-PIN-001`](lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in))
ist opt-in und liest den **Ziel-Span gegenläufig** zur Vorverarbeitung
([§DC-FA-LINK-001.a](#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)),
die Fences sonst entfernt:

1. **Marker-Erkennung + Bindung.** Je vorverarbeiteter Zeile werden
   **Nicht-Bild-Links** (`ExtractLinkSpans` ohne Bilder — ein Content-Pin
   adressiert Text-Zitate, kein Bild) und `<!-- dpin: sha256:<hex> -->`-Marker
   gesucht. Die „nur Whitespace dazwischen"-Bindung wird auf der **rohen** Zeile
   geprüft (die Vorverarbeitung leert Inline-Code zu Leerzeichen). Ein
   Marker bindet an den Link, dessen schließendes `)` ihm **unmittelbar**
   vorausgeht (nur Whitespace dazwischen, gleiche Zeile); bei mehreren
   Links/Markern je Zeile gilt die nächste-vorausgehende-Bindung. Ein Marker ohne
   unmittelbar vorausgehenden Link ist **inert**. Nur gebundene Links werden
   geprüft (opt-in pro Link).
2. **Ziel-Span bestimmen.** Linkziel `datei` bzw. `datei#anker` (relativ zur
   prüfenden Datei, repo-intern). Ohne Anker (oder leerem Anker `datei#`): ganze
   Ziel-Datei. Mit Anker: die
   Heading-Section (Heading mit passendem Slug bis zur nächsten gleich-/
   höherrangigen Überschrift; HTML-Anker ab dessen Zeile). **Welcher Anker das
   ist, entscheidet dieselbe vollständige Anker-Antwort wie bei `versions`**
   (Duplikat-Suffix, Prozent-Dekodierung, nur außerhalb von Fenced-/Inline-Code,
   wörtlicher Vergleich —
   [DC-FA-ANCH-001.b](#dc-fa-anch-001b--inline-html-anker)); der **gehashte
   Span** bleibt der rohe, einschließlich Fenced-Code — das ist die gescopte
   Ausnahme dieses Moduls und keine Lexik-Abweichung. Lässt sich Datei/Anker
   nicht auflösen oder liegt das Ziel außerhalb der **Repo-Wurzel** (repo-escape)
   ⇒ **kein** `pins`-Befund (struktureller Befund bleibt `links`/`anchors`; auch im
   pins-only-Lauf kein Ersatz-Befund, kein Doppelbefund). Ein repo-internes Ziel
   wird **unabhängig vom `pins`-Scope** gehasht (s. Punkt 5).
3. **Normalisierung + Hash.** Der Span wird **roh** gelesen (inkl. Fenced-Code —
   anders als die übrige Vorverarbeitung), dann normalisiert: alle
   Whitespace-Folgen (inkl. Zeilenumbrüche) zu **einem** Leerzeichen kollabiert,
   Rand-Whitespace getrimmt. SHA-256 über das Ergebnis (hex).
4. **Vergleich.** Hash ≠ Pin-Wert ⇒ Grund-Code `link-stale` (Datei, Zeile des
   Links, `target` = Linkziel; `message` nennt erwartet vs. errechnet, gekürzt).
   Gleichheit ⇒ kein Befund. **Diagnose-only**: kein `--repair`-Hunk (Re-Pinnen
   ist menschliche Annahme der Drift).
5. **Scope/Determinismus.** `pins` respektiert den Modul-Scope
   ([§DC-FA-CONF-002.a](#dc-fa-conf-002a--effektiver-scan-scope-pro-modul)) wie
   jedes Modul: **nur die gepinnten Quell-Dateien** folgen dem effektiven Scope.
   Die **Ziel-Datei** wird davon unabhängig zum Hashen gelesen — sie muss nur
   repo-intern und auflösbar sein (der Scope begrenzt nicht, *welche* Ziele
   gehasht werden). Ohne `pins` ist der
   Befundsatz byte-identisch ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus));
   nichts wird geschrieben
   ([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

### DC-FA-IMM-001.a — Immutabilitäts-Pin gegen Core-Drift (`immutable`)

Das Modul `immutable`
([`DC-FA-IMM-001`](lastenheft.md#dc-fa-imm-001--immutabilitäts-pin-gegen-core-drift-modul-immutable-opt-in))
ist opt-in und prüft den **Core** einer Datei gegen einen im Dokument
hinterlegten Pin — **ohne git**, rein im gescannten Arbeitsbaum (die
git-historienbasierte `core(BASE)`-vs-`core(HEAD)`-Form bleibt Out-of-Scope, in
begleitender ADR festgehalten):

1. **Marker-Erkennung.** Eine Datei trägt höchstens **einen** wirksamen Pin: den
   **ersten** `<!-- immutable: sha256:<hex> -->` (case-insensitiver Hex). Die
   Erkennung läuft auf der **vorverarbeiteten** Datei
   ([§DC-FA-LINK-001.a](#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)) —
   ein Marker in Fenced- oder Inline-Code (z. B. ein Syntax-Beispiel in der Doku)
   ist **kein** Pin. Ohne Marker wird die Datei nicht geprüft (opt-in pro Datei);
   weitere Marker nach dem ersten sind inert.
2. **Core bestimmen.** Der Core ist der **rohe** Datei-Inhalt ohne (a) die
   **Marker-Zeile** selbst (sonst Selbstbezug — der Marker trägt den Hash des
   Core) und (b) die Abschnitte, deren getrimmter Heading-Text (ohne
   Markdown-Auszeichnung) ein Eintrag aus `immutable.exclude-sections` ist
   (Heading-Zeile bis zur nächsten gleich-/höherrangigen Überschrift bzw. EOF —
   dieselbe Section-Abgrenzung wie `matrix.exclude-sections`). Der verbleibende
   Inhalt wird **roh** behandelt (inkl. Fenced-Code — eine Änderung an einem
   Code-Beispiel im Core ist eine Core-Änderung, wie bei `pins`).
3. **Normalisierung + Hash.** Wie
   [§DC-FA-PIN-001.a Schritt 3](#dc-fa-pin-001a--content-pin-gegen-inhaltlichen-drift-pins):
   alle Whitespace-Folgen (inkl. Zeilenumbrüche) zu **einem** Leerzeichen
   kollabiert, Rand-Whitespace getrimmt, SHA-256 über das Ergebnis (hex).
4. **Vergleich.** Hash ≠ Pin-Wert ⇒ Grund-Code `core-drift` (Datei, Zeile des
   Markers, `target` = `sha256:<gekürzter-Pin>`; `message` nennt erwartet vs.
   errechnet, gekürzt). Gleichheit ⇒ kein Befund. **Diagnose-only**: kein
   `--repair`-Hunk (Neu-Pinnen ist menschliche Annahme der Änderung).
5. **Scope/Determinismus.** `immutable` respektiert den Modul-Scope
   ([§DC-FA-CONF-002.a](#dc-fa-conf-002a--effektiver-scan-scope-pro-modul)) wie
   jedes Modul; nur Dateien **mit** Marker im effektiven Scope werden geprüft.
   **Kein git**, kein Netz — die Prüfung liest nur die gescannte Datei selbst
   ([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
   Ohne `immutable` ist der Befundsatz byte-identisch
   ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)).

### DC-FA-VCS-001.a — Git-Diff-Immutabilität über eine Commit-Range (`vcs`)

Das Modul `vcs`
([`DC-FA-VCS-001`](lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in))
ist opt-in und vergleicht den **Core** einer immutablen Datei über zwei git-Stände
(`core(BASE)` vs. `core(HEAD)`). Es liest dazu das read-only `.git` über einen
**VCS-Port** (reine-Go-git, **ohne** git-Binary, **ohne** Netz) — die Eingabe ist
gegenüber den hermetischen Modulen um git-Historie + Range **erweitert**, bleibt
aber lokal, lesend und deterministisch. Es ist die git-historienbasierte Hälfte
der Immutabilität, deren hermetische Pin-Hälfte
[§DC-FA-IMM-001.a](#dc-fa-imm-001a--immutabilitäts-pin-gegen-core-drift-immutable)
abdeckt:

1. **Range/Modus (vom CLI geliefert).**
   - `--range <base>..<head>`: zwei Commit-Refs. Der `..`-Separator ist Pflicht —
     ohne ihn ist `base == head`, der Diff leer und der Lauf still grün; das ist
     ein **fail-closed** Nutzungsfehler (Exit 2). Eine leere oder nicht auflösbare
     Basis (nur Nullen, oder kein `^{commit}`) → Fehler (Exit 2).
   - `--staged`: BASE = `HEAD`, HEAD = der staged Index. Existiert kein `HEAD`
     (erster Commit), ist nichts zu schützen (kein Befund).
   - Fehlt das `.git` oder ist es unlesbar → Fehler (Exit 2). **fail-closed.**
2. **Geänderte Kandidaten.** Aus dem Diff `BASE..HEAD` (bzw. staged) werden die
   Pfade gewählt, die der Klasse `vcs.paths` entsprechen (Glob, `matchGlob` wie
   `scan.ignore`). Pro Eintrag zählt der Diff-Status: Modifikation/Typänderung
   (M/T) → Core-Vergleich; Löschung (D) / Umbenennung (R) →
   Pfad-Stabilitäts-Prüfung; Hinzufügung (A) → frei (eine neue Datei ist noch
   nicht immutabel — wie eine frisch reifende ADR).
3. **Immutabilitäts-Bedingung.** Geprüft wird nur, wenn die **BASE**-Version die
   Bedingung `vcs.immutable-when` erfüllt (Zeilen-Regex, erstes Vorkommen
   **außerhalb** von Fenced-Code — eine Datei, die ihren eigenen Kopf als
   Beispiel zeigt, wird dadurch nicht immutabel; dieselbe Zeilen-Menge gilt für
   die Kopf-Status-Zeile aus Schritt 4 und für `matrix.statusOf` — z. B.
   `^\*\*Status:\*\* Accepted`). Trägt BASE die Bedingung nicht (z. B.
   `Proposed`), ist die Datei frei (auch ihre `Proposed → Accepted`-Reifung — das
   automatische Grandfathering des abgelösten Skripts).
4. **Core bestimmen.** Wie
   [§DC-FA-IMM-001.a Schritt 2](#dc-fa-imm-001a--immutabilitäts-pin-gegen-core-drift-immutable):
   roher Inhalt ohne die `vcs.exclude-sections`-Abschnitte (Section-Abgrenzung wie
   `matrix.exclude-sections`); zusätzlich wird die **Kopf**-Status-Zeile
   (`vcs.status-line`, erstes Vorkommen **vor** der ersten `## `-H2) entfernt —
   eine gleichlautende Zeile im Körper (nach der ersten H2) **bleibt** Teil des
   Core (sonst rutschte ein Edit an ihr durch). Anders als `immutable` gibt es
   **keine** Marker-Zeile (der „Pin" ist hier der BASE-Commit selbst, kein im
   Dokument hinterlegter Hash). Normalisierung + SHA-256 wie
   [§DC-FA-PIN-001.a Schritt 3](#dc-fa-pin-001a--content-pin-gegen-inhaltlichen-drift-pins).
5. **Vergleich.**
   - `core(BASE)` ≠ `core(HEAD)` ⇒ Grund-Code `core-drift-vcs` (Körper geändert).
   - HEAD-Status-Zeile (erstes Vorkommen) erfüllt `vcs.head-allow` **nicht** ⇒
     `core-drift-vcs` (unzulässiger Status-Übergang einer immutablen Datei).
   - D/R einer immutablen BASE ⇒ `core-drift-vcs` (Pfad stabil; `message` nennt
     Löschung bzw. das Umbenennungs-Ziel).
   Befund: Datei, Zeile (die `status-line` bzw. 1), `target` = der geprüfte Pfad,
   `message` = Klartext der Verletzung. **Diagnose-only**: kein `--repair`-Hunk.
6. **Determinismus/Read-only.** Identische git-Historie + identische Range ⇒
   identischer, stabil sortierter Befundsatz
   ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)). Der Port liest `.git`
   **nur lesend**, ohne Netz und ohne Schreiben ins Repository
   ([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit));
   die Determinismus-/Read-only-Zusage gilt also auch hier, lediglich der
   **Eingabe-Scope** ist um `.git` + Range erweitert (darum strikt opt-in und
   fail-closed ohne `.git`). Ohne `vcs` ist der Befundsatz byte-identisch.

### DC-FA-COMMITS-001.a — Traceability-Kennung in Commit-Messages über eine Commit-Range (`commits`)

Das Modul `commits`
([`DC-FA-COMMITS-001`](lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in))
ist opt-in und prüft, dass jede geprüfte Commit-Message eine Traceability-Kennung
trägt. Es liest die Commit-**Messages** über **denselben VCS-Port** wie
[§DC-FA-VCS-001.a](#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
(reine-Go-git, **ohne** git-Binary, **ohne** Netz), erweitert um eine
Message-Lese-Operation. Es ist die Portierung des abgelösten Traceability-Skripts:

1. **Quelle/Modus (vom CLI geliefert).**
   - `--range <base>..<head>`: **dieselbe** Range-Semantik wie `vcs`
     ([§DC-FA-VCS-001.a](#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
     Schritt 1) — der `..`-Separator ist Pflicht (ohne ihn `base == head`, still
     grün ⇒ **fail-closed** Exit 2), leere/nicht auflösbare Basis ⇒ Exit 2, fehlendes
     oder unlesbares `.git` ⇒ Exit 2. Der Port liefert die **Nicht-Merge**-Commit-
     Messages der Range (`git rev-list --no-merges`-Parität) in deterministischer
     Reihenfolge (nach Commit-SHA; die Befunde sortiert der Kern ohnehin,
     [`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)).
     Eine **gültige** Range mit 0 Commits ist kein Fehler (nichts zu prüfen ⇒ Exit 0,
     wie `vcs` ohne geänderte Datei).
   - `--commit-msg <datei>`: **Kurzschluss-Modus** (wie `--print-config`/`--trace`;
     nach dem Optionen-Parsing, **ohne** Repo-Scan und **ohne** VCS-Port) — liest
     **eine** Message aus der Datei (`-` = stdin) und prüft nur sie. Nicht lesbare
     Datei/stdin ⇒ Exit 2. Bei leeren `commits.id-patterns` ist nichts zu prüfen ⇒
     **fail-closed** Exit 2 (Misskonfiguration, kein stilles Grün). Dieser Modus
     dient dem lokalen `commit-msg`-Hook (die Pending-Message existiert noch nicht
     als Commit, ist also über keine Range erreichbar).
2. **Message-Bereinigung.** Jede Message wird **uniform** wie git-`strip` bereinigt:
   alles ab der ersten scissors-Zeile (`^#.*>8`, verbose-Diff) entfällt, danach alle
   `#`-Kommentarzeilen (Annahme `core.commentChar='#'`, Default). Range- und
   Message-Quelle wenden **dieselbe** Bereinigung an — gleiche Bewertung, keine
   Divergenz je git-Cleanup-Modus (`-m`=whitespace behält `#`, Editor=strip entfernt
   `#`); eine Kennung muss auf einer **Inhalts**-Zeile stehen, nicht in einem
   Kommentar.
3. **Ausnahme.** Trägt der **Betreff** (erste Zeile der bereinigten Message) einen
   Match auf `commits.exempt-pattern` (Zeilen-Regex), ist die Message frei
   (Selbstkonfiguration `^(Merge |Revert )` — Merge-/Revert-Commits). Leeres
   `exempt-pattern` ⇒ keine Ausnahme. (Merge-Commits sind zudem schon aus der
   Range-Aufzählung ausgenommen; die Betreff-Ausnahme deckt Revert und den
   Message-Modus.)
4. **Kennungs-Prüfung.** Die bereinigte Message muss auf **irgendeiner** Zeile
   mindestens ein `commits.id-patterns`-Muster matchen (ODER über alle Regexe). Kein
   Match ⇒ Grund-Code `commit-untraceable`. Befund: `file` = `target` =
   Commit-Kurz-SHA (Range) bzw. `pending` (Message-Modus) — der „Ort" ist der Commit
   (parallel zu `vcs`, wo `file` == `target`), `line` = 1, `message` = der Betreff.
   **Diagnose-only**: kein `--repair`-Hunk (die Korrektur ist ein neuer Commit / ein
   menschliches `--amend`).
5. **Determinismus/Read-only.** Identische Historie + identische Range ⇒ identischer,
   stabil sortierter Befundsatz
   ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)); der Port liest `.git` **nur
   lesend**, netzlos, ohne Schreiben ins Repository
   ([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
   Der Eingabe-Scope ist um `.git` + Range (bzw. die Message-Datei) erweitert — darum
   strikt opt-in und fail-closed. Ohne `commits` ist der Befundsatz byte-identisch.

### DC-FA-PLAN-001.a — Planning-Lifecycle-Konsistenz (`planning`)

Das Modul `planning`
([`DC-FA-PLAN-001`](lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in))
ist opt-in und prüft die Roadmap-↔-in-progress-Invariante **hermetisch** (nur der
read-only-Arbeitsbaum, **kein** git, **kein** Netz). Es ist ein **Post-Pass**, der
nicht den Markdown-Baum scannt, sondern **eine** deklarierte Roadmap-Datei plus das
Listing ihres Verzeichnisses prüft:

1. **Inert/Config.** Leere `planning.roadmap` ⇒ Modul inert (kein Befund). Sonst
   wird die Datei `planning.roadmap` (Wurzel-relativ, innerhalb der Repo-Wurzel)
   gelesen; ihr Verzeichnis ist das **Slice-Verzeichnis**.
2. **hasSlices.** Das Verzeichnis-Listing wird nach Dateien gefiltert, deren
   Basisname `planning.slice-glob` (Default `slice-*.md`, per `path.Match`) matcht;
   ≥1 Treffer ⇒ `hasSlices`.
3. **Heading-Guard (fail-closed).** Fehlt die kanonische Überschrift
   `planning.heading` (Default `## Aktuelle Welle`) **exakt** als eigene (getrimmte)
   Zeile, **oder kommt sie mehrfach vor** (Aktiv-Status mehrdeutig — welcher Block
   gilt?), ist der Aktiv-Status nicht bestimmbar ⇒ `planning-drift` (kein stilles
   Grün — sonst gälte still „aktiv“ bzw. gälte nur der erste Block). Fehlende/unlesbare
   `planning.roadmap` ⇒ `planning-drift`. **Gezählt wird nur außerhalb von
   Fenced-Code** ([DC-FA-LINK-001.a](#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
   Schritt 1): eine Roadmap, die ihren eigenen Abschnitt in einem Beispiel-Block
   zeigt, gilt sonst als mehrdeutig und meldet ein Falsch-Rot.
4. **hasActive.** Zeilen **innerhalb** eines Fenced-Blocks zählen auch hier nicht:
   sie beenden den Block nicht und tragen keinen Marker — eine Raute-Zeile in einem
   Beispiel-Block beendete den Aktiv-Block sonst vorzeitig und verlöre den
   Ruhe-Marker dahinter. Der Aktiv-Status-Block reicht von der `planning.heading`-Zeile bis
   zur nächsten Überschrift **gleicher oder höherer Ordnung** — bestimmt über
   dieselbe Abschnitts-Mechanik wie die Closure-Fähigkeit desselben Moduls, nicht
   über einen rohen `## `-Präfix-Vergleich: eine eingerückte oder tab-getrennte
   H2 ist eine Überschrift, und eine H1 beendet einen H2-Abschnitt ebenfalls
   (exklusive). Trägt er den `planning.marker` (Default „Keine
   aktive Welle", literaler Teilstring), ist die Welle **ruhend** (`hasActive` =
   false); sonst **aktiv**. Der Marker wird **nur** in diesem Block gesucht (ein
   erklärender Verweis anderswo verfälscht den Status nicht).
5. **Vergleich.** `hasActive ≠ hasSlices` ⇒ Grund-Code `planning-drift`. Befund:
   `file` = `planning.roadmap`, `line` = Zeile der `planning.heading` (bzw. 1),
   `target` = das Slice-Verzeichnis, `message` = Klartext der Drift-Richtung (Slice
   vorhanden + Ruhe-Marker; kein Slice + aktive Welle; oder Heading/Datei fehlt).
   **Diagnose-only** (kein `--repair`-Hunk).
6. **Determinismus/Read-only.** Identischer Arbeitsbaum ⇒ identischer Befund
   ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)); nur lesend, netzlos, ohne
   Schreiben ([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
   Ohne `planning` ist der Befundsatz byte-identisch.

**Wellen-Invariante** (dritte Fähigkeit desselben Moduls, opt-in **innerhalb**
des opt-in Moduls). Sie prüft dieselbe Invariante eine Ebene höher und liest
dazu **denselben** Aktiv-Status wie Schritt 4 — sie bestimmt ihn nicht neu:

- **W1. Inert/Config.** Leere `planning.waves.dir` ⇒ Fähigkeit inert: **kein**
  Wellendokument wird geöffnet, der Befundsatz ist byte-identisch. **Exit 2**
  (Config-Fehler, vor dem Lauf) bei: `dir`/`done-dir` absolut oder mit
  `..`-Segment; **explizit** leerem oder ungültigem `glob`/`results-glob`;
  **explizit** leerem `next-heading`/`closed-heading` (ein abwesender Schlüssel
  liefert den Konventions-Default — die Unterscheidung ist Teil der Zusage,
  sonst wäre die Null-Aussage unerreichbar prüfbar); einem `glob`, der mit
  einem **Platzhalter beginnt** (die Wellen-Kennung braucht ein literales
  Präfix), und einem `results-glob`, der **nicht** mit dem Kennungs-Präfix des
  `glob` beginnt — sonst laufen Datei-Zuordnung und Zeilen-Erkennung
  auseinander. `done-dir` ohne Wert ⇒ `<dir>/done`.
- **W2. Wellen-Mengen.** Aus dem Listing von `dir` (nicht rekursiv) die
  **flachen** Wellendokumente: Basisnamen nach `glob` (Default `welle-*.md`)
  **abzüglich** `results-glob` (Default `welle-*-results.md`) — die
  Ergebnisnotiz ist kein Plan-Dokument. Aus dem Listing von `done-dir` die
  **Ergebnisnotizen** nach `results-glob`. Die Kennung ist das **literale
  Glob-Präfix plus die unmittelbar folgende Ziffernfolge** (beim Default
  `welle-<n>`); zwei Dateien mit derselben Kennung zählen als **eine** Welle.
  **Ein unlesbares `dir` oder `done-dir` ist fail-closed** ⇒ `wave-drift` mit
  dem Verzeichnis als Ziel — im Ruhe-Zustand wäre die leere Menge konsistent,
  und ein Tippfehler im Pfad schaltete die Fähigkeit sonst dauerhaft und
  wortlos ab (dieselbe Disziplin wie `closure.dir`, C2). Eine unlesbare
  Roadmap meldet bereits die Aktiv-Status-Prüfung (fail-closed); die
  Wellen-Fähigkeit schweigt dann — kein Doppelbefund.
- **W3. Aktiv-Status gegen flaches Dokument.** `hasActive` aus Schritt 4
  (dieselbe Größe, dieselbe Quelle). `hasActive` ⇔ **genau ein** flaches
  Wellendokument. Sonst ⇒ `wave-drift`, Befund an der `planning.heading`-Zeile,
  `target` = `dir`; die Meldung nennt beide Richtungen (Roadmap nennt eine Welle
  ohne Dokument / Dokument liegt, aber die Roadmap trägt den Ruhe-Marker) und
  bei mehr als einem Dokument deren Zahl.
- **W4. Register-Tabellen.** In der Roadmap werden die Abschnitte
  `next-heading` (Default `## Nächste Wellen`) und `closed-heading` (Default
  `## Abgeschlossene Wellen`) gelesen. Tabellenzeilen sind Zeilen, die
  **außerhalb** von Fenced-Code mit `|` beginnen (dieselbe Lexik wie
  [`DC-FA-TGT-001.a`](#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets)),
  ohne Kopf- und Trennzeile. Gelesen wird die **erste Spalte**; enthält sie
  keine Kennung (Präfix + Ziffernfolge), wird die Zeile **übersprungen** (eine
  geplante Welle trägt einen Namen, noch keine Kennung — und die
  Trigger-Spalte darf andere Wellen nennen). **Fehlt eine der beiden
  Register-Überschriften, ist das fail-closed** ⇒ `wave-drift` mit der
  **fehlenden Überschrift als Ziel** — die Aktivierung der Fähigkeit ist die
  Behauptung, dass die Roadmap beide Register führt; sonst schaltete ein
  Tippfehler in der Überschrift die Register-Aussagen wortlos ab (dieselbe
  Disziplin wie der Heading-Guard aus Schritt 3). Die Ziele der drei
  `wave-drift`-Bedeutungen sind bewusst verschieden (Verzeichnis · fehlende
  Überschrift), damit die Befund-Deduplikation über (Datei, Zeile, Regel, Ziel,
  Grund) nie zwei Reparaturen zu einem Befund zusammenfallen lässt.
- **W5. Die drei Register-Aussagen.** (a) Eine Kennung in der Vorschau, zu der
  eine Datei existiert (flach **oder** im Ruheort) ⇒ `wave-preview-exists`,
  Befund an der Zeile. (b) Eine Kennung im Abschluss-Register **ohne**
  Ergebnisnotiz ⇒ `wave-results-missing`, Befund an der Zeile. (c) Eine
  Ergebnisnotiz **ohne** Kennung im Abschluss-Register ⇒ `wave-unregistered`,
  Befund an der `closed-heading`-Zeile mit der Notiz als `target`. Alle drei
  sind unabhängig und können nebeneinander melden; die Grund-Codes sind
  getrennt, weil (a) und (b) dieselbe Zeile treffen können und die
  Deduplikation über (Datei, Zeile, Regel, Ziel, Grund) sie sonst
  zusammenfallen ließe. Das `target` der Gegenrichtung ist die
  **Ergebnisnotiz** (Datei-Pfad, Ventil-Parität wie in den übrigen Modulen).

**Closure-Note-Struktur** (zweite Fähigkeit desselben Moduls, opt-in **innerhalb**
des opt-in Moduls). Sie läuft unabhängig von den Schritten 1–5 — die
Aktiv-Status-Prüfung bleibt unberührt, auch wenn die Closure-Prüfung Befunde
liefert (und umgekehrt):

- **C1. Inert/Config.** Leere `planning.closure.dir` ⇒ Fähigkeit inert: **keine**
  Slice-Datei wird geöffnet, der Befundsatz ist byte-identisch zum Stand ohne die
  Fähigkeit. **Exit 2** (Config-Fehler, vor dem Lauf) bei: `dir` absolut oder mit
  `..`-Segment; nicht kompilierendes `heading-pattern`; explizit gesetztes
  `min-sentences` < 1; leerer (oder nur aus Whitespace bestehender)
  `boilerplate`-Eintrag. Ein **abwesendes** `min-sentences` ist kein Fehler,
  sondern der Default — die Unterscheidung „nicht gesetzt" vs. „auf 0 gesetzt"
  ist Teil der Zusage, sonst wäre die Null-Schwelle unerreichbar prüfbar.
- **C2. Kandidaten.** Das Listing von `planning.closure.dir` wird nach Dateien
  gefiltert, deren Basisname `planning.closure.glob` matcht — ein **eigener**
  Filter dieser Fähigkeit. Ist er nicht gesetzt, gilt `planning.slice-glob`
  (Default als **Verweis**, nicht als kopiertes Literal: solange niemand die
  Mengen trennt, gibt es genau ein zu pflegendes Muster, und der Befundsatz ist
  byte-identisch zum Stand ohne den Schlüssel). Ist er **explizit** leer oder
  kein gültiges `path.Match`-Muster ⇒ Exit 2; ein Rückfall auf den Default wäre
  ein stilles Übergehen einer gesetzten Aussage. Ein YAML-`null` (`glob:` ohne
  Wert) gilt als **abwesend**, nicht als leer — wie bei `min-sentences`.
  Die beiden Filter sind getrennt, weil die beiden Fähigkeiten verschiedene
  Grundmengen haben: Schritt 2 zählt, was **noch in Arbeit** ist, C2 prüft, was
  **abgeschlossen** ist — Letzteres kann auch Wellen- oder Etappen-Dokumente
  umfassen. Ein gemeinsames Muster koppelt beide so, dass ein Weiten der einen
  Menge die andere verbiegt (im Grenzfall matcht die Roadmap-Datei selbst und
  der Ruhe-Marker ist dauerhaft falsch-rot).
  Fehlt das gesetzte Verzeichnis oder ist es unlesbar ⇒ `closure-note-missing`
  mit `file` = `planning.closure.dir` (fail-closed, kein stilles Grün).
  **Ebenso fail-closed: null Kandidaten.** `planning.closure.dir` zu setzen
  **ist** die Behauptung, dass dort Closure-Notizen liegen; findet der Lauf
  keine, stimmt die Behauptung nicht mehr (typischer Auslöser: der Bestand
  wandert in Unterordner) und das Gate liefe fortan leer und grün — dieselbe
  Nullmengen-Logik wie bei den Anforderungsquellen der RTM
  ([`DC-FA-REQ-001`](lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen)).
  Ein Repo ohne abgeschlossene Slices setzt den Schlüssel schlicht noch nicht.
  Die Kandidaten werden in stabiler Namens-Reihenfolge geprüft — die Sortierung
  liegt im Kern, nicht beim Dateisystem.
- **C3. Abschnitt bestimmen.** Je Kandidat wird die **erste** Zeile gesucht, die
  (a) **außerhalb** eines Fenced-Code-Blocks liegt — eine `#`-Zeile in einem
  Beispielblock ist keine Überschrift —, (b) eine **echte ATX-Überschrift** nach
  derselben Lexik ist wie in
  [`DC-FA-ANCH-001.a`](#dc-fa-anch-001a--github-slug-algorithmus) (`#`-Folge der
  Länge 1–6, gefolgt von Whitespace; `#1 war ein Thema` ist damit Fließtext,
  **keine** H1), und (c) deren getrimmte Fassung auf
  `planning.closure.heading-pattern` passt. Kein Treffer ⇒
  `closure-note-missing` (`line` = 1). **Mehr als ein Treffer ⇒
  `closure-note-ambiguous`** (`line` = Zeile des **zweiten**) und Abbruch für
  diese Datei: ohne eindeutigen Abschnitt sagt eine Satzzahl nichts, und ein
  zweiter, stehengebliebener Abschnitt ist der typische Rest einer Vorlage.
  Dieselbe Härte trägt der Aktiv-Status-Guard (Schritt 3) längst; die Asymmetrie
  war ein Versäumnis. Sonst reicht der Abschnitt von der Zeile
  **nach** der Überschrift bis zur nächsten **echten** Überschrift gleicher oder
  höherer Ebene bzw. bis zum Dateiende. Die Ebene ist die Länge der `#`-Folge.
  Dass beide Grenzen dieselbe Heading-Lexik nutzen, ist konstitutiv: eine
  eigene `#`-Heuristik beendete den Abschnitt an Fließtext-Zeilen und
  verstünde alles dahinter als „nicht Teil der Notiz" — ein stilles Grün.
- **C4. Bereinigen und messen.** Der Abschnitt wird **einmal** bereinigt, und
  alle folgenden Bedingungen lesen **diesen einen** Text: Fenced-Code-Blöcke
  entfallen, Inline-Code-Spans werden positionserhaltend geleert — dieselbe
  Vorverarbeitung wie im übrigen Scanner
  ([`DC-FA-LINK-001.a`](#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
  Schritte 1–2). Inline-Code ist **Code, kein Fließtext**; er trägt weder Sätze
  noch Floskeln.

  Im bereinigten Text werden die Satzende-Zeichen `.`, `!`, `?` gezählt — aber
  nur, wenn ihnen **Whitespace oder das Zeilenende** folgt. Ohne diese
  Einschränkung zählte jeder Punkt einer Versionsnummer und jeder in einem
  Link-Pfad mit; gemessen am eigenen Bestand sind das mehr als die Hälfte aller
  Vorkommen. Zahl < `min-sentences` ⇒ `closure-note-thin`.

  Anschließend wird **derselbe** bereinigte Text **case-insensitiv** gegen jede
  `planning.closure.boilerplate`-Phrase geprüft — an **Wortgrenzen**: unmittelbar
  vor und hinter dem Treffer darf kein Wortzeichen (`0-9A-Za-z_`) stehen,
  Zeilenanfang und -ende zählen als Grenze. Ein Treffer ⇒
  `closure-note-boilerplate` (der **erste** Treffer benennt die Meldung).

  Als Wortzeichen gilt die **ASCII**-Menge; Umlaute, Anführungszeichen und
  Bindestriche sind damit Grenzen. Für die gelebten Formen ist das richtig
  (`„Ok“`, `Ok.`, `Ok—`), und eine Unicode-Buchstaben-Prüfung wäre eine
  eigene Semantik ohne gemessenen Anlass.

  **Die gemeinsame Bereinigung wirkt in zwei Richtungen.** Sie **verschärft** die
  Zählung und **lockert** zugleich die Floskel-Prüfung: eine Phrase in Backticks
  wird nicht mehr gefunden. Das ist beabsichtigt — eine *zitierte* Floskel ist
  keine benutzte — und ausdrücklich entschieden, statt zwei getrennte bereinigte
  Texte zu führen.
- **C4b. Platzhalter** (nur bei `planning.closure.placeholder: true`; sonst
  entfällt der Schritt vollständig und der Befundsatz ist byte-identisch). Auf dem
  Abschnitts-Text aus C4 werden zusätzlich die **Inline-Code-Spans geleert**
  (positionserhaltend, gleiche Paarungsregel wie
  [§DC-FA-LINK-001.a](#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
  Schritt 2). Auf dem Ergebnis wird der **erste** Treffer gesucht:
  eine öffnende Winkelklammer, der **kein** Wortzeichen und **kein** Schrägstrich
  vorausgeht (Zeilenanfang zählt), deren erstes inneres Zeichen **nicht**
  Whitespace, `!` oder `/` ist und deren Inneres **kein Whitespace enthält** (ein
  Feldname, kein Satz), bis zur nächsten schließenden Klammer **derselben
  Zeile**. Drei **Nachfilter** verwerfen einen Treffer — sie sind Code, nicht Teil
  des Musters, weil sie in ein Muster gepresst unlesbar und unprüfbar würden:
  - das Innere enthält `://` oder `@` (Autolink, Adresse);
  - das Innere bis zum ersten `/`, kleingeschrieben, ist ein bekannter
    HTML-Tag-Name;
  - das konsumierte Vorzeichen ist eine öffnende runde Klammer — dann ist es ein
    Markdown-Linkziel (`](<ziel>)`), nicht eine Vorlage.

  Bleibt ein Treffer übrig ⇒ **ein** Befund `closure-note-placeholder` je
  Kandidat, an der **Zeile des Treffers**; der Treffer selbst steht in der
  `message` (auf 40 Runen gekappt), `target` bleibt wie bei allen
  Closure-Befunden das Verzeichnis (C5). Mehrere Platzhalter derselben Notiz sind dieselbe Reparatur.
  Gearbeitet wird auf dem **einmal** bereinigten Text aus C4; eine eigene,
  engere Sicht führt dieser Schritt nicht. Konsumierte Vorzeichen zweier benachbarter
  Platzhalter können sich überlappen; für den Vertrag ist das folgenlos, weil nur
  der erste Treffer gemeldet wird.
  Die Vorlage des Konsumenten nutzt Lookbehind und Lookahead und ist damit in RE2
  **nicht** ausdrückbar; beide sind feste Zeichen-Prüfungen und hier durch
  **Konsumieren** des Vorzeichens bzw. eine Zeichenklasse ersetzt — eine
  Portierung mit eigener Testlast, keine Übernahme.

- **C5. Befund-Form.** `file` = die Slice-Datei (bzw. das Verzeichnis bei C2),
  `rule` = `planning`, `target` = `planning.closure.dir`, `message` = Klartext
  der verletzten Bedingung (fehlender Abschnitt · Ist-/Soll-Satzzahl ·
  getroffene Floskel · getroffener Platzhalter).
  `line` = Zeile der Abschnitts-Überschrift (bzw. 1) — **außer** bei C4b: dort
  die Zeile des **Treffers**, weil dort die Reparatur liegt und der Abschnitt
  mehrzeilig ist.
  **Diagnose-only** (kein `--repair`-Hunk): eine Closure-Notiz schreibt der Autor.
  Ein Kandidat kann `closure-note-thin`, `closure-note-boilerplate` und
  `closure-note-placeholder` **nebeneinander** tragen (verschiedene
  Bedingungen), aber `closure-note-missing` schließt alle drei aus — ohne
  Abschnitt gibt es nichts zu messen.

Die Fähigkeit prüft **Struktur**, nicht Bedeutung: ob eine formal ausreichende
Notiz inhaltlich trägt, ist semantisch und ausdrücklich nicht zugesagt
([`DC-FA-PLAN-001`](lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
§Out-of-Scope).

### DC-FA-STRUCT-001.a — Struktur-Invarianten innerhalb eines Dokuments (`structure`)

Das opt-in Modul `structure`
([`DC-FA-STRUCT-001`](lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in))
prüft **hermetisch** je Regel eine Dokumentklasse. Es ist ein **Post-Pass** wie
`planning`: es scannt nicht den Markdown-Baum, sondern die von den Regel-Globs
getroffenen Dateien.

1. **Inert/Config.** Leere Regel-Liste ⇒ Modul inert (kein Befund, keine Datei
   geöffnet). **Exit 2** vor dem Lauf bei: Regel ohne `files`; ungültiges Glob in
   `files`/`exempt-paths`; weder `section` noch `section-pattern` gesetzt **oder**
   beide; `sections` weder `one` noch `each`; nicht kompilierendes
   `section-pattern`/`forbid-pattern`/`require-pattern`; **explizit** gesetztes
   `min-sentences` < 1 oder `max-tasks` < 0; leerer Eintrag in `require-all`.
   Ein **abwesender** Zahlen-Schlüssel ist kein Fehler, sondern „Bedingung
   aus" — die Unterscheidung „nicht gesetzt" vs. „auf 0 gesetzt" ist Teil der
   Zusage, sonst wäre die Null-Schwelle unerreichbar prüfbar.
2. **Kandidaten.** Die `files`-Globs werden gegen **Wurzel-relative Pfade über
   den gesamten Baum** ausgewertet — **unabhängig** von `scan.roots`/`scan.ignore`
   (die `SKIP_DIRS` gelten weiterhin). Eine Regel benennt ihre Dateien selbst;
   `structure` kennt deshalb **kein** `<modul>.scope`. Kandidaten sind **nur
   Markdown-Dateien** (Endungs-Menge des Scanners) und **keine Symlinks** (§1) —
   eine nur über einen Symlink erreichbare Datei prüft auch dieses Modul nicht.
   Abgezogen wird `exempt-paths`. Die verbleibenden Dateien werden stabil sortiert geprüft
   ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)).
   **Null Kandidaten ⇒ `section-missing`** (`file` = der Glob, `line` = 1,
   `target` = die **Regel-Identität** aus `files` und Abschnitts-Selektor, damit
   zwei leer laufende Regeln über derselben Dateimenge nicht unter der
   Befund-Deduplikation zusammenfallen; identische Identität ⇒ Exit 2) —
   **auch dann, wenn erst `exempt-paths` die Menge geleert hat**: sonst schaltete
   ein Ventil die Regel still ab (dieselbe Nullmengen-Logik wie bei den
   Anforderungsquellen der RTM). Ist schon der **Dateibaum** nicht lesbar, meldet
   **jede** Regel `section-missing` mit ihrer Identität (fail-closed) — ein
   leerer Befundsatz wäre von „alle Regeln erfüllt" nicht zu unterscheiden.
3. **Abschnitte finden.** Gesucht werden alle Zeilen, die (a) **außerhalb** eines
   Fenced-Code-Blocks liegen, (b) **echte ATX-Überschriften** nach der Lexik aus
   [`DC-FA-ANCH-001.a`](#dc-fa-anch-001a--github-slug-algorithmus) sind und
   (c) passen: bei `section` ist der Vergleichsgegenstand die **getrimmte
   Überschriften-Zeile einschließlich der `#`-Folge**, exakt; bei
   `section-pattern` dieselbe Zeile gegen RE2.
4. **Kardinalität** — `sections` entscheidet, wie viele Treffer erwartet werden:
   - **`one`** (Default): 0 ⇒ `section-missing` (`line` = 1); > 1 ⇒
     `section-ambiguous` (`line` = **zweiter** Treffer) und **Abbruch für diese
     Datei**, keine Bedingungs-Prüfung. Der Abbruch gilt der **Datei innerhalb
     dieser Regel** — andere Regeln und andere Dateien bleiben unberührt.
   - **`each`**: 0 ⇒ `section-missing` (`line` = 1); sonst wird **jeder** Treffer
     einzeln nach Schritt 5–6 geprüft. Mehrfachtreffer sind kein Befund.
5. **Bereinigen.** Je Abschnitt (Zeile **nach** der Überschrift bis zur nächsten
   echten Überschrift gleicher oder höherer Ebene bzw. Dateiende) werden die
   Fenced-Code-Blöcke entfernt **und die Inline-Code-Spans geleert** — dieselbe
   Vorverarbeitung wie im übrigen Scanner
   ([`DC-FA-LINK-001.a`](#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
   Schritte 1–2) und dieselbe wie beim Preset-Partner
   ([`DC-FA-PLAN-001.a`](#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
   Schritt C4). Alle Bedingungen arbeiten auf **diesem** Text. Die Folge gehört
   dazugesagt: ein `forbid-pattern`, das auf ein Wort in Backticks zielt, trifft
   **nicht** — Inline-Code ist Code, kein Fließtext.
6. **Bedingungen prüfen**, jede optional, jede mit **eigenem** Grund-Code, damit
   zwei Verletzungen desselben Abschnitts nicht unter der Befund-Deduplikation
   (Datei, Zeile, Regel, Ziel, Grund) zusammenfallen:

   | Bedingung | erfüllt, wenn | sonst |
   |---|---|---|
   | `non-empty` | der bereinigte Text mindestens ein Nicht-Whitespace-Zeichen trägt | `section-empty` |
   | `min-sentences` | die Zahl der Satzende-Zeichen (`.`, `!`, `?`) ≥ Schwelle ist — gezählt wird nur, was **vor Whitespace oder Zeilenende** steht, wie beim Preset-Partner | `section-thin` |
   | `max-tasks` | die Zahl der **Task-Items** ≤ Schwelle ist | `section-oversized` |
   | `forbid-pattern` | das Muster **nicht** matcht | `section-forbidden` |
   | `require-pattern` | das Muster matcht | `section-pattern-missing` |
   | `require-all` | **jede** Marke vorhanden ist | `section-marker-missing` |

   **Task-Item** ist eine Zeile, die — nach optionalem Whitespace — mit einem
   Listen-Marker (`-`, `*`, `+` oder `<ziffern>.`), Whitespace und `[ ]` bzw.
   `[x]`/`[X]` beginnt.
   **Marke `M` vorhanden** heißt: eine Zeile beginnt — nach optionalem
   Listen-Marker und Whitespace — mit einem **hervorgehobenen Textlauf**
   (`**…**`), dessen Inhalt mit `M` anfängt und dort **endet oder mit einem
   nicht-alphanumerischen Zeichen weitergeht**. Alphanumerisch heißt hier
   **unicode-weit** (Buchstabe oder Ziffer beliebiger Schrift): in
   `- **Ergebnisüberblick:**` setzt das `ü` den Textlauf fort, die Marke
   `Ergebnis` ist damit **nicht** erfüllt.
   Das ist bewusst ein **anderer** Wort-Begriff als bei der Floskel-Prüfung des
   Preset-Partners ([`DC-FA-PLAN-001.a`](#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
   Schritt C4, ASCII einschließlich `_`): die Floskel vergleicht eine
   **konfigurierte Phrase** in Fließtext, die Marke einen **Feldnamen** in einer
   Auszeichnung — ein Umlaut setzt hier ein Wort fort, dort beendet er eines,
   und ein Unterstrich umgekehrt. Damit gelten `- **M:**`,
   `**M:**` und `- **M (Zusatz):**` als Treffer, `M` im Fließtext nicht.
   `forbid-pattern`/`require-pattern` matchen gegen den **gesamten** bereinigten
   Abschnitts-Text; `^`/`$` binden dort an Text-, nicht an Zeilen-Grenzen (wer
   Zeilen meint, schreibt `(?m)`).
7. **Befund-Form.** `file` = die geprüfte Datei (bzw. der Glob in Schritt 2),
   `rule` = `structure`, `target` = die **Regel-Identität** (`files`-Glob und
   Abschnitts-Selektor), `line` = Zeile der Abschnitts-Überschrift (bzw. 1).
   Die Identität im `target` ist konstitutiv: die Deduplikation vergleicht
   (Datei, Zeile, Regel, Ziel, Grund), und ohne sie verlöre man je Datei den
   Befund der zweiten Regel. **Diagnose-only** (kein `--repair`-Hunk).
8. **Determinismus/Read-only.** Identischer Arbeitsbaum ⇒ identischer, stabil
   sortierter Befundsatz; nur lesend, netzlos
   ([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
   Ohne aktives `structure` ist der Befundsatz byte-identisch.

**Die Closure-Note-Struktur ist ein Preset dieser Schritte im Modus `one`.** Die
Schritte C3–C5 von
[`DC-FA-PLAN-001.a`](#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
sind dieselbe Mechanik mit fester Bedingungs-Auswahl und eigenen Grund-Codes;
beide teilen Abschnitts-Findung, Kardinalitäts-Behandlung, Bereinigung und
Zählung. Eine Änderung an einer der beiden Stellen ohne die andere ist ein
Spezifikations-Bug; ein Akzeptanzkriterium beider Anforderungen hält sie
zusammen.

### DC-FA-TRK-001.a — Getrackt-Status auflösbarer Referenz-Ziele (`tracked`)

Opt-in-Prüfung **je gescannter Quell-Datei** ([`DC-FA-TRK-001`](lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in)):
prüft die **Datei-Ebene** der aufgelösten repo-internen Link-/Bild-Ziele
gegen den **git-Index** — dieselbe Auflösungs-**Mechanik** wie
[`DC-FA-LINK-001`](lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links),
aber **unabhängig** von der Aktivierung des Moduls `links` (ein fokussierter
`--enable tracked`-Lauf prüft vollständig). Die dritte Nutzung des VCS-Ports
(`vcs`: Range-Diff, `commits`: Commit-Messages, `tracked`: Index), diesmal
**ohne** Range/`--staged`; die Index-Menge wird **einmal je Lauf** geladen.

1. **Aktivierung/fail-closed.** Nur bei explizit aktiviertem `tracked`. Aktiv,
   aber kein lesbares `.git` unter der Repo-Wurzel ⇒ **Exit 2** mit
   stderr-Hinweis (kein stilles Grün). Ohne aktives `tracked` ist der Befundsatz
   byte-identisch und `.git` wird nicht gelesen.
2. **Index-Menge.** Der VCS-Port liefert die Menge der **getrackten Pfade**
   (git-Index des Arbeitsbaums; eine gestagte, noch nie committete Datei ist
   enthalten — der Index ist die Wahrheit, keine `.gitignore`-Interpretation).
3. **Kandidaten.** Alle **erfolgreich aufgelösten, existierenden**
   repo-internen **Datei**-Ziele der gescannten Quell-Dateien (Fragment/Anker
   wird verworfen; Bild-Ziele eingeschlossen). **Kein** Kandidat sind:
   Verzeichnis-Ziele (der Index führt nur Dateien), Symlink-Referenzen (Ziel
   ist oder durchläuft einen Symlink — kategorisch
   [`DC-FA-LINK-002`](lastenheft.md#dc-fa-link-002--symlink-ablehnung)-Domäne;
   der Index führt den realen Pfad, nicht den Alias) und nicht existierende
   Ziele — deren Befund bleibt `target-missing`
   ([`DC-FA-LINK-001`](lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links),
   **kein Doppelbefund** — Prinzip von
   [`DC-FA-PIN-001`](lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in)).
   Der Modul-Scope ([`DC-FA-CONF-002`](lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope))
   gilt für die **Quell**-Dateien wie bei jedem Modul.
4. **Ventil.** Matcht der aufgelöste Zielpfad ein `tracked.exempt-targets`-Glob
   (`matchGlob`, wie `scan.ignore`), wird die Referenz übersprungen —
   **referenz-weit** (analog `codepaths.ignore-refs`), für absichtlich
   untrackte Ziele. Jedes Glob wird config-zeitig **segmentweise** validiert
   (ungültig/leer ⇒ Exit 2 — kein still wirkungsloses Ventil).
5. **Vergleich.** Ziel existiert ∧ Pfad ∉ Index ⇒ Grund-Code `target-untracked`.
   Befund: `file`/`line` = Quell-Datei/Link-Zeile, `target` = der **aufgelöste**
   Zielpfad (dieselbe Form, die das Ventil matcht — nicht die rohe
   Link-Schreibweise), `message` benennt, dass das Ziel auf einem frischen Klon fehlt
   (untracked/gitignoriert). **Diagnose-only** (kein `--repair`-Hunk; `git add`
   bzw. Committen ist eine menschliche Entscheidung). Pro Referenz maximal ein
   Befund; mehrere Referenzen auf dasselbe untrackte Ziel melden je Referenz
   (Zeilen-genau, wie `target-missing`).
6. **Determinismus/Read-only.** Identischer Arbeitsbaum + identischer Index ⇒
   identischer Befundsatz ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus));
   `.git` wird ausschließlich gelesen, kein Netz, kein Schreiben
   ([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
   in der `vcs`-Lesart: erweiterte, aber lokale/lesende Eingabe).

### DC-FA-TGT-001.a — Deklarations-Konsistenz Doku und Build-Targets (`targets`)

Das Modul `targets`
([`DC-FA-TGT-001`](lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in))
ist opt-in und prüft **hermetisch** (Filesystem-Port `ReadFile`, **kein** git,
**kein Ausführen** des Makefile, **kein** Netz), ob die in der Doku als
` `make X` ` behaupteten Build-Targets und die real definierten Makefile-Regeln
übereinstimmen. Wie `planning` ein **Post-Pass**, der nicht den Markdown-Baum
scannt, sondern **deklarierte** Dateien liest:

1. **Inert/Config/fail-closed.** Leeres `targets.makefiles` ⇒ Modul inert
   (keine Regelmenge, kein Befund). Sonst wird jede konfigurierte Datei
   (Wurzel-relativ, innerhalb der Repo-Wurzel) gelesen; eine fehlende/unlesbare
   konfigurierte Datei ⇒ **Exit 2** mit stderr-Hinweis (kein stilles Grün).
   **Richtung 1** (Schritt 4) läuft nur bei nicht-leerem `targets.doc-tables`,
   **Richtung 2** (Schritt 5) nur bei nicht-leerem `targets.authority` — die
   beiden Richtungen sind voneinander unabhängig.
2. **Makefile-Regelmenge.** Aus jeder `targets.makefiles`-Datei werden die
   Regelnamen extrahiert: eine Zeile, die am **Zeilenanfang** mit einem oder
   mehreren durch Leerzeichen getrennten Namen (`[A-Za-z][A-Za-z0-9_-]*`)
   gefolgt von `:` beginnt, wobei nach dem `:` **kein** `=` folgt (Zuweisungen
   `X :=`/`X ?=` ausgenommen) und kein Name mit `.` beginnt (`.PHONY`/
   `.DEFAULT_GOAL` ausgenommen). Pattern-Rules (`%`), variabel benannte und über
   `include` eingezogene Targets sind **kein** Kandidat — rein statische
   Zeilen-Heuristik, in Parität zum abgelösten `tools/gate-consistency.sh`.
3. **Dokumentierte-Target-Menge (Tabellen-Scoping).** Aus jeder Doku-Datei
   werden ` `make X` `-Tokens (`X` = `[a-z][a-z0-9_-]*`, in Parität zum Skript)
   **ausschließlich aus Tabellenzeilen außerhalb von Fenced-Code** extrahiert
   (ein Beispiel-Block, der eine Ziel-Tabelle zeigt, dokumentiert nichts) — eine Tabellenzeile ist eine
   Zeile, deren **erstes Zeichen** ein Pipe `|` ist (**Spalte 0**, in Parität zu
   `grep -E '^\|'`; **Einrückung zählt nicht**). Prosa-Erwähnungen (z. B.
   „Richtig: ` `make gates` `") zählen **nicht**, sonst entstünden aus
   entfernten, nur in Prosa erwähnten Targets spuriöse `gate-phantom`.
4. **Richtung 1 (Phantom).** Für jede `targets.doc-tables`-Datei: jedes
   dokumentierte ` `make X` `, dessen `X` **nicht** in der Makefile-Regelmenge
   (Schritt 2) ist ⇒ Grund-Code `gate-phantom`. Befund: `file`/`line` =
   Doku-Datei/Tabellenzeile, `target` = `X`, `message` = dokumentiertes Target
   ohne Makefile-Regel.
5. **Richtung 2 (undokumentiert).** Jede Makefile-Regel `X` (Schritt 2), die
   **nicht** in `targets.exempt-targets` steht und **nicht** in der aus
   `targets.authority` (Tabellen-Scoping) gewonnenen Menge enthalten ist ⇒
   Grund-Code `gate-undocumented`. Befund: `file`/`line` = Makefile/Regelzeile,
   `target` = `X`, `message` = Makefile-Regel ohne Doku-Deklaration. Leeres
   `targets.authority` ⇒ Richtung 2 entfällt.
6. **Diagnose-only / Determinismus / Read-only.** Kein `--repair`-Hunk;
   identischer Arbeitsbaum ⇒ identischer Befundsatz
   ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)), nur lesend, netzlos,
   ohne Schreiben
   ([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
   Ohne aktives `targets` ist der Befundsatz byte-identisch. Befunde entstehen
   **je Vorkommen** (mehrere Tabellenzeilen bzw. Regelzeilen mit demselben Namen
   ⇒ je ein Befund an seiner Datei:Zeile) — dieselbe **Detektion** wie das
   abgelöste `tools/gate-consistency.sh`, aber ohne dessen `sort -u`-Namens-
   Deduplizierung (fundstellen-präzise statt namens-eindeutig).

### DC-FA-SRC-001.a — Upstream-Content-Drift externer Quellen (`sources`)

Das Modul `sources`
([`DC-FA-SRC-001`](lastenheft.md#dc-fa-src-001--upstream-content-drift-externer-quellen-modul-sources-opt-in-netz))
ist opt-in und — neben `external` — die einzige Netz-Tür
([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
Es prüft den **Inhalt** einer auf einen `sha256` gepinnten externen Quelle
(Content-Hash), nicht bloße Erreichbarkeit (`external`) und nicht ein
repo-internes Ziel (`pins`). **Post-Pass** nach dem Datei-Scan (wie `external`)
über die eingesammelten Pins:

1. **Aktivierung/Sammeln.** Nur bei explizit aktiviertem `sources`. Beim Scan
   werden zwei Pin-Quellen gesammelt: (a) **Marker**
   `<!-- source-pin: [zip] sha256:<hex> -->` (das optionale Schlüsselwort `zip`
   markiert ein Archiv, parallel zu `unpack: zip`), gebunden an den **unmittelbar
   links** stehenden `http(s)`-Link derselben Zeile — bei mehreren Links den
   nächstgelegenen; ohne vorausgehenden `http(s)`-Link inert. Und (b) die
   **Config**-Einträge `sources[]`. Beide ergeben je einen Pin
   `{url, sha256, unpack ∈ {none, zip}}`. `sources` akzeptiert **keinen**
   `sources.scope` (die Config-Fläche ist eine bare Liste); die Marker-Sammlung
   nutzt den **globalen** Scan-Scope — bewusste Ausnahme zu
   [`DC-FA-CONF-002`](lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope)
   (`sources.scope` ⇒ Exit 2, unbekannter Schlüssel).
2. **fail-closed / Skopus.** Eine malformte Direktive (`source-pin` ohne
   `sha256:<hex>`) oder ein ungültiger Config-Eintrag (fehlende `url`/`sha256`,
   `url` nicht `http(s)`, `sha256` nicht genau 64 Hex-Zeichen, unbekanntes
   `unpack`) ⇒ **Exit 2** mit stderr-Hinweis (kein stilles Grün). Der `sha256`
   wird **case-insensitiv** geführt (Pin und errechneter Hash zu Kleinbuchstaben
   normalisiert). Ein **wohlgeformter** `source-pin` an einem repo-internen oder
   Nicht-`http(s)`-Link ist **inert** (kein Befund, kein `sources`-Kandidat —
   repo-interne Ziele sind `pins`/`links`-Domäne, kein Doppelbefund); die
   Direktive ist dort bewusst wirkungslos (mehrdeutig mit `pins`), **nicht**
   fail-closed.
3. **Fetch (begrenzt).** Je Pin wird die `url` über den Netz-Port geholt (GET;
   Timeout/Parallelität wie `external`; Redirects werden wie `external` bis zu
   **fünf** verfolgt, gehasht wird der Inhalt der **finalen** Antwort). Der Body
   wird nur bis zu einer festen Obergrenze gelesen (**≤ 64 MiB**). **Jeder** der
   folgenden Fälle ⇒ Grund-Code `source-unreachable` (Befund: `target` = URL,
   `message` benennt die Ursache), **ohne** Hash-Schritt: Netzfehler,
   HTTP-Status ≥ 400, Timeout, mehr als fünf Redirects, überschrittenes
   Body-Limit **oder** (bei `unpack: zip`) eine 2xx-Antwort, die **kein gültiges
   Zip** ist (Nicht-Zip / HTML-Fehlerseite / abgeschnittenes Archiv).
   „Unerreichbar" umfasst damit jede Lage, in der sich **kein** vergleichbarer
   Inhalt materialisieren lässt — bewusst getrennt von `source-drift` (Inhalt
   da, aber ≠ Pin).
4. **Hash (byte-genaues Content-Manifest).** `unpack: none` (Default): `sha256`
   über die **rohen Antwort-Bytes**, als Kleinbuchstaben-Hex. `unpack: zip`: die
   Antwort wird als Zip gelesen; **nur reguläre Datei-Einträge**
   (Verzeichnis-Einträge — Name endet auf `/` — ausgenommen).
   Zip-Bomben-/Ressourcen-Schutz: die **entpackte Gesamtgröße** ist auf
   **≤ 256 MiB** und die **Eintragszahl** auf **≤ 10 000** begrenzt;
   Überschreitung ⇒ `source-unreachable`. Das kanonische **Content-Manifest**
   ist byte-genau definiert:
   - je Datei-Eintrag eine Zeile `<hex>` + **zwei Leerzeichen** + `<pfad>`, wobei
     `<hex>` der `sha256` des **Eintrags-Inhalts** in Kleinbuchstaben-Hex ist und
     `<pfad>` der **volle Zip-interne Pfad** normalisiert (Backslashes → `/`,
     führendes `./` und `/` entfernt; **kein** Basisname — verschachtelte
     Verzeichnisse bleiben im Pfad, daher keine Basisnamen-Kollision);
   - die Zeilen **aufsteigend nach `<pfad>`** sortiert (byteweise, `LC_ALL=C` —
     **nicht** nach der ganzen Zeile; bei identischem `<pfad>` — Zip erlaubt
     Duplikate — sekundär nach `<hex>`, ein stabiler eindeutiger Tie-Break), je
     mit `\n` terminiert und konkateniert;
   - der Archiv-Hash ist der `sha256` dieser Manifest-Bytes (Kleinbuchstaben-Hex).

   Damit ist der Hash **reihenfolge-invariant** gegen die
   Zip-Eintrags-Reihenfolge und vollständig determiniert
   ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)). Er folgt
   **konzeptionell** dem committet-vendored `SHA256SUMS`-Muster (Per-Datei-Hash-
   Manifest), ist aber durch die Pfad-Sortierung und -Normalisierung
   **eigenständig** kanonisiert — **nicht** byte-identisch zur unsortierten
   `sha256sum <glob>`-Ausgabe.
5. **Vergleich.** Ist-Hash und Pin werden **case-insensitiv** (beide
   Kleinbuchstaben) verglichen. Gleich ⇒ kein Befund. Ungleich ⇒ Grund-Code
   `source-drift`; `message` führt den **vollständigen** Ist-`sha256` in
   Kleinbuchstaben (Re-Pin-Vorlage, schreibweisen-stabil). Befund-`file`/`line`:
   für einen **Marker-Pin** die Marker-Datei und -Zeile; für einen **Config-Pin**
   `file` = die **tatsächlich geladene** Konfigurationsdatei (konventionell
   `.d-check.yml`, unter `--config` die dort genannte —
   [`DC-FA-CLI-012.a`](#dc-fa-cli-012a--konfigurations-pfad---config) Schritt 6)
   und `line` = die Zeile des `url`-Feldes des Eintrags
   (die YAML-Dekodierung führt sie; sonst `1`); `target` = URL. **Diagnose-only**
   (kein `--repair`-Hunk — ein Re-Pin ist eine menschliche Entscheidung).
6. **Determinismus/Read-only.** Identische Antwort-Bytes ⇒ identischer
   Hash/Befund ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)); das
   geprüfte Repository wird nie geschrieben, Netz nur zu den explizit gepinnten
   `url` ([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
   — die zweite, wie `external` opt-in netz-öffnende Ausnahme).

Die Grund-Codes `source-drift`/`source-unreachable` landen mit der
Modul-Implementierung in [§4](#4-grund--und-fehler-codes)
(AllReasons-↔-§4-Lockstep).

## 2. Datenstrukturen und Schemas

### Befund

| Feld | Typ | Bedeutung |
|---|---|---|
| `file` | string | Pfad relativ zur Repo-Wurzel, `/`-getrennt |
| `line` | integer ≥ 1 | Zeile des Vorkommens |
| `rule` | string | Regelmodul-Name; gültige Werte siehe [`DC-FA-CLI-002`](lastenheft.md#dc-fa-cli-002--regelmodul-auswahl) |
| `target` | string | geprüftes Ziel (Linkziel, Kennung, URL) |
| `reason` | string | Grund-Code (siehe [§4](#4-grund--und-fehler-codes)) |
| `message` | string | menschenlesbare Erläuterung (nicht stabilitätsgarantiert) |

**Text-Format** (Default, stdout, ein Befund pro Zeile,
[`DC-FA-CLI-004`](lastenheft.md#dc-fa-cli-004--ausgabeformate)):

```
<file>:<line>	<target>	<reason>
```

Zusammenfassung auf stderr: `d-check: <N> Datei(en) geprüft, <M> Befund(e)`.

### JSON-Ausgabe (`--json`)

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["findings", "summary", "exitCode"],
  "properties": {
    "findings": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["file", "line", "rule", "target", "reason"],
        "properties": {
          "file":    {"type": "string"},
          "line":    {"type": "integer", "minimum": 1},
          "rule":    {"type": "string"},
          "target":  {"type": "string"},
          "reason":  {"type": "string"},
          "message": {"type": "string"}
        }
      }
    },
    "summary": {
      "type": "object",
      "required": ["filesChecked", "findingCount"],
      "properties": {
        "filesChecked": {"type": "integer", "minimum": 0},
        "findingCount": {"type": "integer", "minimum": 0}
      }
    },
    "exitCode": {"type": "integer", "enum": [0, 1]}
  }
}
```

Nutzungs-/Umgebungsfehler (Exit 2) erzeugen **kein** JSON-Dokument,
sondern eine stderr-Meldung (siehe [§4](#4-grund--und-fehler-codes)).

### JSON-Diagnose (`--doctor --json`)

Dasselbe Dokument wie die [JSON-Ausgabe](#json-ausgabe---json) — gleiche
`summary`/`exitCode` und dieselben Befund-Basisfelder —, aber jeder
`findings`-Eintrag trägt zusätzlich `reasonText` und `fixCandidate`
([`DC-FA-CLI-007`](lastenheft.md#dc-fa-cli-007--diagnose-modus)). Delta
zum `items`-Schema:

```json
{
  "required": ["file", "line", "rule", "target", "reason", "reasonText", "fixCandidate"],
  "properties": {
    "reasonText": {"type": "string"},
    "fixCandidate": {
      "type": ["object", "null"],
      "required": ["original", "replacement", "note"],
      "properties": {
        "original":    {"type": "string"},
        "replacement": {"type": "string"},
        "note":        {"type": "string"}
      }
    }
  }
}
```

`fixCandidate` ist `null` (nicht weggelassen), wo kein eindeutiger
Kandidat existiert — das `null` ist die Aussage „kein eindeutiger Fix".
Das optionale `message`-Feld des Basis-`items`-Schemas der
[JSON-Ausgabe](#json-ausgabe---json) (nicht stabilitätsgarantiert) bleibt
unverändert Teil jedes Eintrags und kann zusätzlich erscheinen.

### YAML-Ausgabe (`--yaml`)

`--yaml` serialisiert **dasselbe Dokument** wie `--json` (knappe Befunde)
bzw. `--doctor --yaml` dasselbe wie `--doctor --json` (Diagnose mit
`reasonText`/`fixCandidate`) — identische Struktur und Feldnamen
(`findings`/`summary`/`exitCode`, `filesChecked`/`findingCount`, …), nur
als YAML statt JSON. `--json` und `--yaml` schließen sich aus
([`DC-FA-CLI-004`](lastenheft.md#dc-fa-cli-004--ausgabeformate)); ein
fehlender Fix-Kandidat steht auch in YAML explizit als `fixCandidate:
null`. Deterministisch (feste Feld-Reihenfolge, keine Map) wie die
JSON-Ausgabe.

### `.d-check.yml`

Unbekannte Schlüssel sind Fehler (striktes Decoding). Eine leere oder
nur kommentierte Datei ist ein YAML-Null-Dokument und wirkt wie keine
Datei (Defaults). Explizit leere Listen ersetzen die Defaults durch
die leere Menge (`scan.roots: []` prüft nichts, `modules: []` lässt
kein Modul laufen — bewusste Setzung, kein Fehler). Kommentiertes
Vollbeispiel:

```yaml
scan:
  roots: [docs, spec]            # ersetzt die Default-Wurzeln; müssen existieren
  ignore: ["docs/archive/**"]    # Glob, relativ zur Repo-Wurzel
modules: [links, anchors, ids]   # ersetzt DEFAULT_MODULES
ids:
  scope:                         # optional: ersetzt den globalen Scan
    roots: [spec, docs/user]     #   nur für dieses Modul (DC-FA-CONF-002)
    ignore: []
  patterns:                      # Reihenfolge = Präzedenz
    - regex: 'ADR-\d{4}'
      target: docs/plan/adr/     # Definition (Datei oder Verzeichnis)
      link-policy: always        # prose (Default) | always: auch Inline-Code linkpflichtig
      exempt-paths: [CHANGELOG.md]  # Globs ohne Linkpflicht (nackt wie Inline-Code)
matrix:
  classes:                       # Reihenfolge = Präzedenz
    - name: contract
      paths: [spec/lastenheft.md]
    - name: adr
      paths: ["docs/plan/adr/[0-9]*.md"]
  rules:
    - {from: contract, to: adr, allow: false}
  status:
    forbidden: [superseded, deprecated]
    allow-supersede-lineage: true        # ablösende Datei darf auf ihr abgelöstes Ziel zeigen
    supersede-fields: [Supersedes, Aenderungstyp]   # Felder, aus denen die Ablösung gelesen wird
  exclude-sections: [Historie, Geschichte]
external:
  timeout-seconds: 10
  parallel: 4
codepaths:
  roots: [docs, tools]           # Präfixe für Wurzel-relative Inline-Code-Pfade
  exempt-paths: [CHANGELOG.md, "docs/reviews/**"]   # Dateien ganz ohne codepath-Prüfung (datei-weit)
sources:                         # Netz, opt-in; Upstream-Content-Drift externer Quellen
  - url: https://example.org/regelwerk.md   # Einzeldatei: Hash der Roh-Bytes
    sha256: <64-hex>
  - url: https://example.org/bundle.zip     # Archiv: Hash des Content-Manifests
    sha256: <64-hex>
    unpack: zip
```

Jede Verletzung eines Constraints der folgenden Tabelle führt zu
Exit 2 ohne Prüfung
([`DC-FA-CONF-001`](lastenheft.md#dc-fa-conf-001--konfigurationsdatei)).

| Schlüssel | Typ | Default | Constraint |
|---|---|---|---|
| `links.resolve-from` | list | leer | Gruppen ortsfester Verweise; je Gruppe `dirs` (string[], **wandernde** Geschwister-Verzeichnisse, Wurzel-relativ — deren Dateien sind Quellen; mindestens 2, sonst Exit 2) und optional `fixed-dirs` (string[], ortsfeste Orte: hypothetische Ziele, deren Dateien keine Quellen sind). Ein Verzeichnis darf nur in **einer** Gruppe `dirs`-Mitglied sein (Exit 2). Ohne den Block byte-identisch |
| `scan.roots` | string[] | `DEFAULT_SCAN_ROOTS` | alle hier deklarierten Wurzeln müssen existieren und innerhalb der Repo-Wurzel liegen (Exit 2); nur die Default-Wurzeln (kein `scan.roots` gesetzt) sind optional; `"."` steht für die gesamte Repo-Wurzel (rekursiv; die `SKIP_DIRS` aus [§3](#3-defaults-und-konstanten) gelten immer und sind nicht konfigurierbar) |
| `scan.ignore` | string[] | leer | Glob-Syntax; Muster prunen auch den Verzeichnis-Abstieg — ein vollständig ignorierter Teilbaum (`pfad/**` oder direkt matchendes Muster) wird nicht betreten, unlesbare ignorierte Verzeichnisse sind dadurch kein Laufzeitfehler |
| `ignore-refs` | object[] | leer | Geteiltes Referenz-Ventil ([`DC-FA-REF-001`](lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)), honoriert von `links`/`anchors`/`codepaths`; jeder Eintrag `{in?, refs, keep?}`. Ohne Einträge byte-identisch |
| `ignore-refs[].in` | string | — (offen) | optionaler Glob (wie `scan.ignore`) auf die **Quelldatei** (die Datei, in der die Referenz steht); ungesetzt = repo-weit |
| `ignore-refs[].refs` | string[] | leer | Glob-Liste auf den **aufgelösten Ziel-Pfad**; leer/fehlend = Eintrag inert (ignoriert nichts) |
| `ignore-refs[].keep` | string[] | leer | Glob-Liste; Ausnahmen — ein Ziel wird nur ignoriert, wenn `refs` matcht **und** `keep` **nicht** (`keep` gewinnt reihenfolge-unabhängig, kein gitignore-Last-Match) |
| `modules` | string[] | `DEFAULT_MODULES` | nur gültige Modulnamen |
| `<modul>.scope.roots` | string[] | — (globaler Scope) | Pflicht, wenn `scope` gesetzt ist; Constraints wie `scan.roots` (Existenz, kein Repo-Escape, `"."` = Repo-Wurzel, leere Liste = nichts) |
| `<modul>.scope.ignore` | string[] | leer | wie `scan.ignore` (Glob, Abstiegs-Pruning) |
| `hostpaths.prefixes` | string[] | Development, home, Users, Volumes, mnt, media | ersetzt die Default-Liste; Einträge sind nicht-leere Verzeichnisnamen ohne `/` (Exit 2) |
| `ids.patterns[].regex` | string | — | muss kompilieren und darf den Leerstring nicht matchen (Exit 2) |
| `ids.patterns[].target` | string | — | muss existieren und innerhalb der Repo-Wurzel liegen |
| `ids.patterns[].link-policy` | string | `prose` | nur `prose` oder `always` (Exit 2); `always` macht auch Inline-Code-Vorkommen linkpflichtig |
| `ids.patterns[].exempt-paths` | string[] | leer | Glob (wie `scan.ignore`, relativ zur Repo-Wurzel); Dateien ohne Linkpflicht für das Muster — nackte wie Inline-Code-Vorkommen, unabhängig von der `link-policy` |
| `matrix.classes[].name` | string | — | eindeutig |
| `matrix.classes[].paths` | string[] | — | Glob (Mitgliedschaft) |
| `matrix.classes[].order` | string[] | leer | Glob-Liste, autoritativste Schicht zuerst (Rang = Index des ersten Treffers); nur zusammen mit `direction` (Exit 2) |
| `matrix.classes[].direction` | string | leer | nur `no-downward` (Exit 2 bei unbekanntem Wert); verlangt nicht-leeres `order` (Exit 2) — keine still wirkungslose Richtungs-Deklaration |
| `matrix.classes[].token` | string | leer | Regex; erkennt Referenzen auf diese Klasse als **bare ID-Token** im Fließtext ([`DC-FA-MTX-003`](lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)). Muss kompilieren und darf den Leerstring nicht matchen (Exit 2). Ohne `token` nur Link-Erkennung |
| `matrix.rules[]` | {from,to,allow} | — | Klassen müssen deklariert sein |
| `matrix.status.forbidden` | string[] | `[superseded, deprecated]` | case-insensitiv |
| `matrix.status.allow-supersede-lineage` | bool | `false` | nimmt die deklarierte Supersede-Lineage-Kante von der Status-Prüfung aus (nur `matrix-inactive`); ohne `true` byte-identisch |
| `matrix.status.supersede-fields` | string[] | leer | Feldnamen (z. B. `Supersedes`, `Aenderungstyp`), aus denen die Ablösung gelesen wird; Einträge nicht leer (Exit 2); nur wirksam bei `allow-supersede-lineage: true` |
| `matrix.exclude-sections` | string[] | leer | Vergleich gegen den getrimmten Heading-Text ohne Markdown-Auszeichnung, case-sensitiv |
| `matrix.exempt-paths` | string[] | leer | Glob (wie `scan.ignore`, relativ zur Repo-Wurzel); Dateien ganz ohne `matrix`-Prüfung — Grandfathering immutabler Dokumente ([`DC-FA-MTX-003`](lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)) |
| `external.timeout-seconds` | integer | 10 | 1–300 |
| `external.parallel` | integer | 4 | 1–16 |
| `sources[].url` | string | — | Pflicht; absolute `http`/`https`-URL (sonst Exit 2) |
| `sources[].sha256` | string | — | Pflicht; genau 64 Hex-Zeichen, **case-insensitiv** verglichen (sonst Exit 2) |
| `sources[].unpack` | string | `none` | nur `none` oder `zip` (Exit 2); `zip` hasht das `LC_ALL=C`-sortierte Content-Manifest statt der Roh-Bytes |
| `codepaths.roots` | string[] | leer | Präfixe relativ zur Repo-Wurzel: nicht leer, nicht absolut, kein `..` (Exit 2); `./`/`../` werden immer erkannt |
| `codepaths.exempt-paths` | string[] | leer | Glob (wie `scan.ignore`, relativ zur Repo-Wurzel); Dateien ganz ohne `codepaths`-Prüfung — datei-weit, unabhängig von `roots` |
| `codepaths.ignore-refs` | string[] | leer | **Alias** des geteilten `ignore-refs` ([`DC-FA-REF-001`](lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)), skopiert auf `codepaths` (wie ein Eintrag ohne `in`/`keep`): aufgelöste Ziel-Pfade, die matchen, werden weder existenz-, escape- noch anker-geprüft (kein `codepath-missing`/`repo-escape`/`anchor-missing`) — **referenz-weit**; byte-identisch zur bisherigen Fassung |
| `codepaths.check-lines` | bool | `false` | Zeilen-Referenz eines Inline-Code-Pfads (`datei:<von>-<bis>`) verifizieren: existierendes Ziel mit ≥ `<bis>` Zeilen (sonst `citation-out-of-range`) und `<von> ≤ <bis>` (sonst `citation-inverted-range`). Default aus ⇒ das Suffix wird wie bisher abgetrennt und verworfen (byte-identisch, [`DC-FA-CODE-001.a`](#dc-fa-code-001a--pfade-in-inline-code) Schritt 6) |
| `versions.pin-pattern` | string | — | Regex; muss kompilieren und darf den Leerstring nicht matchen (Exit 2); die gefundene Version steht in Capture-Gruppe 1, sonst zählt der ganze Treffer |
| `versions.current-from` | string | `version.md#aktuell` | `datei#anker` oder `datei`; die Datei muss existieren und innerhalb der Repo-Wurzel liegen, der adressierte Span muss eine Version (`v?\d+\.\d+\.\d+`) tragen (sonst Exit 2) |
| `versions.exempt-paths` | string[] | leer | Glob (wie `scan.ignore`, relativ zur Repo-Wurzel); Dateien ganz ohne `versions`-Prüfung — datei-weit; die `current-from`-Datei ist stets ausgenommen |
| `immutable.exclude-sections` | string[] | leer | Heading-Titel, deren Abschnitte **nicht** zum gehashten Core zählen (Vergleich gegen den getrimmten Heading-Text ohne Markdown-Auszeichnung, case-sensitiv — wie `matrix.exclude-sections`); für ADRs typisch `[Geschichte]` ([`DC-FA-IMM-001`](lastenheft.md#dc-fa-imm-001--immutabilitäts-pin-gegen-core-drift-modul-immutable-opt-in)) |
| `vcs.paths` | string[] | leer | Glob-Klasse (wie `scan.ignore`) der zu schützenden Dateien; nur in der Range geänderte Pfade dieser Klasse werden geprüft; leer ⇒ Modul inert ([`DC-FA-VCS-001`](lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)) |
| `vcs.immutable-when` | string | — | Zeilen-Regex; die **BASE**-Version gilt als immutabel, wenn ihr **erstes** Vorkommen matcht (z. B. `^\*\*Status:\*\* Accepted`); muss kompilieren (sonst Exit 2) |
| `vcs.exclude-sections` | string[] | leer | wie `immutable.exclude-sections` — Heading-Titel, deren Abschnitte nicht zum Core zählen (für ADRs `[Geschichte]`) |
| `vcs.status-line` | string | leer | Zeilen-Regex der **Kopf**-Status-Zeile; ihr erstes Vorkommen **vor** der ersten `## `-H2 wird aus dem Core entfernt (eine gleichlautende Körper-Zeile bleibt); leer ⇒ keine Status-Zeile gestrippt |
| `vcs.head-allow` | string | leer | Zeilen-Regex; die HEAD-Status-Zeile (erstes Vorkommen) muss matchen, sonst `core-drift-vcs` (unzulässiger Status-Übergang); leer ⇒ keine Status-Übergangs-Prüfung |
| `commits.id-patterns` | string[] | leer | Regex-Liste der gültigen Traceability-Kennungen; eine bereinigte Message ohne Match auf **irgendein** Muster ⇒ `commit-untraceable`; jedes Muster muss kompilieren (sonst Exit 2); leer ⇒ Modul inert (Range-Modus) bzw. Exit 2 (Message-Modus, nichts zu prüfen) ([`DC-FA-COMMITS-001`](lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in)) |
| `commits.exempt-pattern` | string | leer | Zeilen-Regex gegen den **Betreff** (erste Zeile); Match ⇒ Message kennungs-frei erlaubt (Selbstkonfig `^(Merge \|Revert )`); muss kompilieren (sonst Exit 2); leer ⇒ keine Ausnahme |
| `planning.roadmap` | string | leer | `datei` (Wurzel-relativ, innerhalb der Repo-Wurzel); die Roadmap-Datei mit dem Aktiv-Status-Abschnitt. Ihr Verzeichnis ist das Slice-Verzeichnis. Leer ⇒ Modul inert ([`DC-FA-PLAN-001`](lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)) |
| `planning.heading` | string | `## Aktuelle Welle` | die kanonische H2-Überschrift des Aktiv-Status-Abschnitts (exakter, getrimmter Zeilen-Vergleich); fehlt sie ⇒ `planning-drift` (fail-closed) |
| `planning.marker` | string | `Keine aktive Welle` | literaler Ruhe-Marker-Teilstring; nur im `planning.heading`-Block gesucht; vorhanden ⇒ ruhende Welle |
| `planning.slice-glob` | string | `slice-*.md` | Glob (`path.Match` auf den Basisnamen) der Slice-Dateien im Roadmap-Verzeichnis; ≥1 Treffer ⇒ aktive-Welle-erwartet; ein explizit gesetztes Muster muss ein gültiges `path.Match`-Glob sein (sonst Exit 2 — verhindert ein fail-open Silent-Grün) |
| `planning.waves.dir` | string | leer | `verzeichnis` (Wurzel-relativ, innerhalb der Repo-Wurzel); Ort der **flachen** Wellendokumente. Gesetzt ⇒ dritte Fähigkeit aktiv; leer ⇒ inert (kein Wellendokument wird geöffnet) |
| `planning.waves.done-dir` | string | `<dir>/done` | Ruheort der Ergebnisnotizen |
| `planning.waves.glob` | string | `welle-*.md` | Basisnamen-Glob (`path.Match`) der **Plan**-Dokumente; `results-glob` wird davon **abgezogen** |
| `planning.waves.results-glob` | string | `welle-*-results.md` | Basisnamen-Glob der **Ergebnisnotizen** — das verpflichtende Artefakt einer geschlossenen Welle |
| `planning.waves.next-heading` | string | `## Nächste Wellen` | H2 des Vorschau-Registers (exakter, getrimmter Zeilen-Vergleich) |
| `planning.waves.closed-heading` | string | `## Abgeschlossene Wellen` | H2 des Abschluss-Registers |
| `planning.closure.dir` | string | leer | `verzeichnis` (Wurzel-relativ, innerhalb der Repo-Wurzel); das Verzeichnis der abgeschlossenen Slices, deren Closure-Notiz **strukturell** geprüft wird. Leer ⇒ Closure-Fähigkeit inert (keine Slice-Datei wird gelesen, Befundsatz byte-identisch); gesetzt, aber fehlend/unlesbar ⇒ `closure-note-missing` auf dem Verzeichnis (fail-closed) ([`DC-FA-PLAN-001`](lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)) |
| `planning.closure.glob` | string | Wert von `planning.slice-glob` | Basisnamen-Glob (`path.Match`) der Kandidaten **dieser** Fähigkeit — getrennt von `planning.slice-glob`, weil die Lifecycle-Invariante „noch in Arbeit“ und die Closure-Struktur „abgeschlossen“ zählt. Nicht gesetzt ⇒ Default-Verweis (Befundsatz byte-identisch); **explizit** leer oder kein gültiges Glob ⇒ Exit 2; kein Kandidat unter dem gesetzten Glob ⇒ `closure-note-missing` auf dem Verzeichnis (fail-closed) |
| `planning.closure.heading-pattern` | string | `^#{2,3} .*Closure-Notiz` | RE2 gegen die **getrimmte** Überschriften-Zeile; **genau ein** Treffer eröffnet den geprüften Abschnitt (mehrere ⇒ `closure-note-ambiguous`) (bis zur nächsten Überschrift gleicher oder höherer Ebene). Kein Treffer ⇒ `closure-note-missing`; ein nicht kompilierendes Muster ⇒ Exit 2 |
| `planning.closure.min-sentences` | int | `4` | Mindestzahl der Satzende-Zeichen (`.`, `!`, `?`) im Abschnitt **nach** Entfernen der Fenced-Code-Blöcke; darunter ⇒ `closure-note-thin`. Wert < 1 ⇒ Exit 2 (eine Schwelle von 0 wäre ein stilles Grün) |
| `planning.closure.placeholder` | bool | `false` | schaltet die vierte Struktur-Bedingung frei: der unausgefüllte Rumpf einer Vorlage ⇒ `closure-note-placeholder`. Aus ⇒ Schritt C4b entfällt, Befundsatz byte-identisch. Die Erkennung ignoriert **Inline-Code** und verwirft Autolinks/Adressen sowie HTML-Tags per Nachfilter |
| `planning.closure.boilerplate` | string[] | leer | literale Floskel-Phrasen, **case-insensitiv** und an **Wortgrenzen** gegen den bereinigten Abschnitts-Text geprüft; ein Treffer ⇒ `closure-note-boilerplate`. Bewusst **leer** per Default — der Vertrag bringt keine sprach-spezifischen Phrasen mit; ein leerer Eintrag ⇒ Exit 2 (er träfe jeden Text) |
| `structure[].files` | string | — | Glob (Pfad, wie `scan.ignore`) über **Wurzel-relative** Pfade des gesamten Baums, unabhängig von `scan.roots`/`scan.ignore`; Pflicht je Regel. Null Kandidaten — auch nach Abzug von `exempt-paths` — ⇒ `section-missing` auf dem Glob ([`DC-FA-STRUCT-001`](lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)) |
| `structure[].section` | string | leer | Heading-**Klartext**, exakter Vergleich der getrimmten Überschriften-Zeile **einschließlich `#`-Folge**. Genau eines von `section`/`section-pattern` ist Pflicht — beide oder keines ⇒ Exit 2 |
| `structure[].section-pattern` | string | leer | RE2 gegen dieselbe getrimmte Zeile; Alternative zu `section` |
| `structure[].sections` | string | `one` | `one` = genau ein Treffer erwartet (0 ⇒ `section-missing`, > 1 ⇒ `section-ambiguous` + Abbruch für diese Datei); `each` = jeder Treffer wird geprüft (für **wiederkehrende** Abschnitte, 0 ⇒ `section-missing`). Anderer Wert ⇒ Exit 2 |
| `structure[].non-empty` | bool | `false` | der bereinigte Abschnitts-Text muss mindestens ein Nicht-Whitespace-Zeichen tragen ⇒ sonst `section-empty` |
| `structure[].min-sentences` | int | abwesend (aus) | Mindestzahl der Satzende-Zeichen (`.`, `!`, `?`) im bereinigten Abschnitts-Text — Fenced-Code entfernt **und Inline-Code geleert** —, gezählt nur, was **vor Whitespace oder Zeilenende** steht (`a.b.c.d.` ⇒ **eins**, nicht vier) ⇒ sonst `section-thin`; **explizit** < 1 ⇒ Exit 2 |
| `structure[].max-tasks` | int | abwesend (aus) | Obergrenze der Task-Items **im Abschnitt**, nicht dateiweit ⇒ sonst `section-oversized`; **explizit** < 0 ⇒ Exit 2 |
| `structure[].forbid-pattern` | string | leer | RE2 gegen den **gesamten** bereinigten Abschnitts-Text; Treffer ⇒ `section-forbidden` |
| `structure[].require-pattern` | string | leer | RE2, Spiegelbild von `forbid-pattern`; **kein** Treffer ⇒ `section-pattern-missing`. Deckt zugesagte Aussagen, die **innerhalb** einer Auszeichnung stehen und deshalb keine Marke sind |
| `structure[].require-all` | string[] | leer | benannte Marken, die **alle** vorkommen müssen — als hervorgehobener Textlauf am Zeilen-Anfang nach optionalem **Listen-Marker** (`- **M:**`, `**M:**`, `- **M (Zusatz):**`) ⇒ sonst `section-marker-missing`; leerer Eintrag ⇒ Exit 2 |
| `structure[].exempt-paths` | string[] | leer | Glob (wie `scan.ignore`) über die Quell-Pfade; Treffer werden von **dieser** Regel nicht geprüft — hebeln den Leerlauf-Befund aber nicht aus |
| `tracked.exempt-targets` | string[] | leer | Glob (wie `scan.ignore`); **aufgelöste Ziel-Pfade**, die matchen, werden nicht auf Getrackt-Status geprüft — **referenz-weit** (analog `codepaths.ignore-refs`), für absichtlich untrackte Ziele; jedes Glob **segmentweise** gültig und nicht leer (sonst Exit 2); ohne Eintrag byte-identisch ([`DC-FA-TRK-001`](lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in)) |
| `targets.makefiles` | string[] | leer | Wurzel-relative Makefile-Dateien, aus denen Regelnamen per statischer Zeilen-Heuristik extrahiert werden; leer ⇒ Modul inert; eine fehlende/unlesbare Datei ⇒ Exit 2 ([`DC-FA-TGT-001`](lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)) |
| `targets.doc-tables` | string[] | leer | Wurzel-relative Doku-Dateien; ihre `make X`-**Tabellenzeilen** (nur Zeilen mit Pipe in Spalte 0, keine Prosa) werden gegen die Makefile-Regelmenge geprüft (Richtung 1 `gate-phantom`); leer ⇒ Richtung 1 entfällt; fehlende Datei ⇒ Exit 2 |
| `targets.authority` | string | leer | Wurzel-relative Doku-Datei; **jede** nicht-exempte Makefile-Regel muss dort als `make X`-Tabellenzeile stehen (Richtung 2 `gate-undocumented`); leer ⇒ Richtung 2 entfällt; fehlende Datei ⇒ Exit 2 |
| `targets.exempt-targets` | string[] | leer | Regelnamen (**exakt**-Vergleich, **kein** Glob — anders als `tracked.exempt-targets`, das Pfad-Globs matcht), die von der Doku-Pflicht (Richtung 2) ausgenommen sind (Utility-Targets); ohne Eintrag prüft Richtung 2 jede Regel |
| `trace.requirements.source` | string | `spec/lastenheft.md` | Wurzel-relative Anforderungsdatei; muss innerhalb der Repo-Wurzel liegen; leer/abwesend ⇒ Default und aktiviert keinen Strict-Guard. Ein **nichtleerer expliziter** Wert aktiviert fail-closed bei fehlender Quelle/null erkannten Anforderungen ([`DC-FA-REQ-001`](lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen)) |
| `trace.requirements.id-pattern` | string | `[A-Z][A-Z0-9]*-(?:FA-[A-Z]+\|QA)-\d+[A-Za-z]?` | Regex; erkennt eine Anforderungs-Kennung als **Ganz-Token** im Heading bzw. **Ganzzelle** der ID-Spalte und als Vorkommen in ADR-/Slice-Dateien; muss kompilieren (sonst Exit 2); leer ⇒ Default |
| `trace.requirements.format` | string | `headings` | `headings` oder `table`; leer ⇒ `headings`; anderer Wert ⇒ Exit 2. `table` aktiviert den Strict-Guard und verlangt `trace.requirements.table`; ohne `trace`-Block bleibt die RTM byte-identisch |
| `trace.requirements.table` | map | — | Pflicht genau bei `format: table`, bei `headings`/leer verboten (sonst Exit 2); bindet Tabellenzellen nach exaktem Header-Namen |
| `trace.requirements.table.id-column` | string | — | nichtleerer Header-Name der ID-Spalte; muss in einer relevanten Tabelle genau einmal vorkommen, sonst Exit 2 |
| `trace.requirements.table.text-column` | string | — | einzelner nichtleerer Header-Name der Text-/Titelspalte; alternativ zu `text-columns` |
| `trace.requirements.table.text-columns` | string[] | — | nichtleere, duplikatfreie Liste alternativer exakter Text-Header; je relevanter Tabelle muss genau einer und jeder deklarierte Name quellweit mindestens einmal vorkommen; alternativ zu `text-column` |
| `trace.requirements.table.modality-column` | string | leer | optionaler Header-Name der alleinigen Modalitätsquelle; gesetzt ⇒ muss genau einmal vorkommen, sonst Exit 2; ohne ihn klassifiziert `modality` die Textspalte |
| `trace.requirements.table.duplicate-ids` | string | `error` | `error`, `first` oder `last`; Default weist doppelte IDs ab, die Overrides lösen historische Mehrfachdefinitionen deterministisch |
| `trace.adrs.dir` | string | `docs/plan/adr` | Wurzel-relatives Referenz-Verzeichnis der ADRs (rekursiv gescannt); fehlt es ⇒ keine ADR-Referenzen (kein Fehler); leer ⇒ Default |
| `trace.adrs.file-pattern` | string | `^(\d{4})-.*\.md$` | Regex auf den **Basisnamen**; **Capture-Gruppe 1** = Owner-Kennung der Datei; muss kompilieren **und ≥1 Capture-Gruppe** haben (sonst Exit 2 — sonst wäre die Owner-Ableitung undefiniert); Dateien ohne Treffer übersprungen; leer ⇒ Default |
| `trace.adrs.id-prefix` | string | `ADR-` | der Owner-Kennung (Capture-Gruppe 1) vorangestellt (`ADR-` + Capture `NNNN` ⇒ die ADR-Kennung `ADR-NNNN`); leer ⇒ Default `ADR-` (ein leerer ADR-Präfix ist nicht ausdrückbar, [`DC-FA-CLI-009`](lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix) Out-of-Scope) |
| `trace.slices.dir` | string | `docs/plan/planning` | Wurzel-relatives Referenz-Verzeichnis der Slices (rekursiv); fehlt es ⇒ keine Slice-Referenzen (alle Anforderungen Waisen); leer ⇒ Default |
| `trace.slices.file-pattern` | string | `^(slice-\d+)-.*\.md$` | Regex auf den **Basisnamen**; **Capture-Gruppe 1** = Owner-Kennung; muss kompilieren **und ≥1 Capture-Gruppe** haben (sonst Exit 2); Dateien ohne Treffer (z. B. `README.md`) übersprungen; leer ⇒ Default |
| `trace.slices.id-prefix` | string | leer | der Owner-Kennung vorangestellt; Default **leer** (der Slice-Dateiname trägt die volle Kennung `slice-NNN`); setzbar (z. B. `Slice ` ⇒ `Slice 063`) |
| `trace.coverage` | list | leer | opt-in **Liste** kuratierter Coverage-Quellen ([`DC-FA-COV-001`](lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)); leer/abwesend ⇒ **keine** Coverage-Spalte, RTM byte-identisch ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)) |
| `trace.coverage[].files` | string[] | — | **explizite** Wurzel-relative Pfad-Liste (keine `dir`/`file-pattern`-Ableitung — gegen ADR-/Slice-Kontamination); nicht leer; jede Datei muss existieren (sonst Exit 2, fail-closed) |
| `trace.coverage[].label` | string | — | feste Owner-Kennung in der **Coverage**-Spalte (z. B. `Trace`); nicht leer (Exit 2) |
| `trace.coverage[].ranges` | bool | `true` | expandiert Range-/Enum-Notation: `<FAM>-AAA..BBB` (inklusiv, breiten-erhaltend) + `<FAM>-AAA/BBB/CCC`; jede expandierte ID gegen `trace.requirements.id-pattern` validiert (Nicht-Treffer verworfen); `AAA>BBB` oder Breiten-Mismatch ⇒ Exit 2. `false` ⇒ nur exakte IDs |
| `trace.coverage[].sections` | string[] | leer | **Whitelist**: nur diese H2/H3-Abschnitte zählen als Coverage (Span wie `matrix.exclude-sections` — Überschrift bis zur nächsten gleich-/höherrangigen); leer ⇒ ganze Datei |
| `trace.coverage[].exclude-sections` | string[] | leer | **Blacklist**: diese Abschnitte zählen **nicht** (Span wie `matrix.exclude-sections`); gegen „…ohne Design-Artefakt"-Listen, die sonst nicht gedeckte IDs kreditierten |
| `trace.requirements.modality` | map | leer | opt-in Modalitäts-Klassifikation ([`DC-FA-MOD-001`](lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in)); abwesend ⇒ **keine** Modality-Spalte, `--require-complete` gatet alle Waisen, RTM byte-identisch ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)) |
| `trace.requirements.modality.levels` | map | Built-in DE+EN | Stufen-Name → Schlüsselwort-Liste; leer ⇒ Default (`must`/`should`/`may` mit DE+EN-RFC-2119-Verben inkl. Negationen `DARF NICHT`→must, `MUSS NICHT`→may). Ein leerer Stufen-Name/leeres Keyword ⇒ Exit 2. Matching: erster Treffer im Anforderungs-Body (Span wie `matrix.exclude-sections`), **längste Phrase zuerst**, case-insensitiv, wortgrenzen-genau |
| `trace.requirements.modality.require-levels` | string[] | `[must]` | Stufen, deren Waisen `--require-complete` gaten ([`DC-FA-CLI-011`](lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)); jeder Eintrag muss ein deklarierter Stufen-Name **oder** `unknown` sein (sonst Exit 2). Anforderung ohne Keyword-Treffer ⇒ Stufe `unknown` |
| `trace.cross-consistency` | map | leer | opt-in Kreuzverweis-Abgleich zweier Traceability-Sichten ([`DC-FA-XREF-001`](lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)); leer/abwesend ⇒ kein Abgleich, RTM byte-identisch ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)); verlangt `forward` **und** `backward` (sonst Exit 2) |
| `trace.cross-consistency.forward.file` | string | — | Wurzel-relative Vorwärts-Tabellendatei (Anforderung → Design); muss innerhalb der Repo-Wurzel liegen und existieren (sonst Exit 2) |
| `trace.cross-consistency.forward.sections` | string[] | leer | **Whitelist** der H2/H3-Abschnitte (Span wie `matrix.exclude-sections`); leer ⇒ ganze Datei |
| `trace.cross-consistency.forward.exclude-sections` | string[] | leer | **Blacklist** der Abschnitte (Span wie `matrix.exclude-sections`) |
| `trace.cross-consistency.forward.req-column` | string | — | header-gebundener Name der Anforderungs-Spalte; muss je relevanter Tabelle genau einmal vorkommen (sonst Exit 2) |
| `trace.cross-consistency.forward.design-column` | string | — | header-gebundener Name der Design-Spalte; genau einmal je Tabelle (sonst Exit 2) |
| `trace.cross-consistency.forward.design-pattern` | string | — | Regex; extrahiert Design-Artefakt-IDs aus der Design-Zelle (alle Treffer); muss kompilieren (sonst Exit 2) |
| `trace.cross-consistency.forward.req-pattern` | string | `trace.requirements.id-pattern` | Regex; erkennt die Anforderungs-IDs in der `req-column`-Zelle — symmetrisch zu `backward.req-pattern`. Muss kompilieren (sonst Exit 2). **Der Default ist eine bewusste Kopplung, keine Ableitung:** wer die RTM auf eine Teilmenge scopt, aber weiter über die volle Menge vergleichen will, setzt ihn explizit — die Vergleichs-Schlüsselmenge ist **nicht** die RTM-Anforderungsmenge |
| `trace.cross-consistency.forward.ranges` | bool | `true` | Range-/Enum-Expansion der Anforderungs-IDs der ID-Spalte (`<FAM>-AAA..BBB`/`/`-Enum, wie `trace.coverage[].ranges`); `false` ⇒ nur exakte IDs |
| `trace.cross-consistency.backward.file` | string | — | Wurzel-relative Rück-Kanten-Datei (Design → Anforderung); muss innerhalb der Repo-Wurzel liegen und existieren (sonst Exit 2) |
| `trace.cross-consistency.backward.sections` | string[] | leer | **Whitelist** der Abschnitte, deren `edge-column`-Tabellen zählen; leer ⇒ ganze Datei |
| `trace.cross-consistency.backward.artifact-id-column` | string | `first` | Header-Name (dann genau einmal je relevanter Tabelle, sonst Exit 2) **oder** Sentinel `first` (erste Spalte — für über die Tabellen heterogene ID-Header); die Artefakt-ID wird per `forward.design-pattern` aus der Zelle extrahiert (bewusst geteilt — sichert den gemeinsamen Namensraum mit der Vorwärts-Sicht) |
| `trace.cross-consistency.backward.edge-column` | string | — | header-gebundener Name der Kanten-Spalte (z. B. `Bezug`); genau einmal je relevanter Tabelle (sonst Exit 2) |
| `trace.cross-consistency.backward.req-pattern` | string | — | Regex; erkennt Anforderungs-IDs in der `edge-column`-Zelle; muss kompilieren (sonst Exit 2) |
| `trace.cross-consistency.backward.ranges` | bool | `true` | Range-/Enum-Expansion der Anforderungs-IDs der Kanten-Zelle; `false` ⇒ nur exakte IDs |
| `trace.cross-consistency.mode` | string | `equal` | `equal` (beide Differenzen `F\B` und `B\F` gaten) oder `superset` (nur `B\F`); anderer Wert ⇒ Exit 2 |
| `trace.cross-consistency.exclude-req` | string | leer | Regex; Anforderungs-IDs, die vor dem Diff aus beiden Sichten fallen (Ableitungssprünge in Mittelschichten); muss kompilieren (sonst Exit 2); leer ⇒ keine Ausnahme |

**Glob-Auswertung.** Alle Glob-Felder (`scan.ignore`, `<modul>.scope.ignore`,
`matrix.classes[].paths`/`.order`, die `*.exempt-paths`) werden **segmentweise
über Go-`path.Match`** ausgewertet (`*`/`?` matchen nicht über `/`), zusätzlich
`**` für beliebig viele Segmente. Eine **negierte** Zeichenklasse ist `[^…]`
(Go-Syntax), **nicht** `[!…]` (Shell/fnmatch) — `[!a]` matcht die Literale `!`
und `a`, nicht „alles außer `a`".

## 3. Defaults und Konstanten

| Name | Wert | Begründung | Bezug |
|---|---|---|---|
| `DEFAULT_SCAN_ROOTS` | `docs/`, `spec/` (rekursiv, optional) + `*.md` der Repo-Wurzel | [`DC-FA-SCAN-001`](lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln) | — |
| `SKIP_DIRS` | `.git`, `node_modules`, `build`, `target`, `dist`, `vendor`, `.venv`, `__pycache__`, `.idea`, `.vscode`, `.gradle` | immer übersprungen | [`DC-FA-SCAN-001`](lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln) |
| `DEFAULT_MODULES` | `links`, `anchors` | [`DC-FA-CLI-002`](lastenheft.md#dc-fa-cli-002--regelmodul-auswahl) | — |
| `EXTERNAL_TIMEOUT` | 10 s | [`DC-FA-EXT-001`](lastenheft.md#dc-fa-ext-001--externe-links-modul-external-opt-in) | konfigurierbar |
| `EXTERNAL_PARALLEL` | 4 | begrenzte Parallelität | konfigurierbar |
| `REDIRECT_MAX` | 5 | [`DC-FA-EXT-001`](lastenheft.md#dc-fa-ext-001--externe-links-modul-external-opt-in) | fest |
| Exit-Codes | 0 / 1 / 2 | [`DC-FA-CLI-003`](lastenheft.md#dc-fa-cli-003--exit-codes) | fest |

## 4. Grund- und Fehler-Codes

Grund-Codes der Befunde (stabil, maschinenlesbar):

| Code | Modul | Bedingung |
|---|---|---|
| `target-missing` | links | Linkziel existiert nicht |
| `repo-escape` | links, codepaths | aufgelöstes Ziel verlässt die Repo-Wurzel |
| `symlink` | links | Ziel ist/enthält Symlink (Vorrang vor `repo-escape`) |
| `anchor-missing` | anchors, codepaths | Anker entspricht keinem Heading-Slug und keinem HTML-Anker der Zieldatei |
| `id-unlinked` | ids | Kennung im Fließtext ohne Markdown-Link |
| `matrix-forbidden` | matrix | Referenz zwischen Klassen nicht erlaubt (**Link** oder, bei gesetztem `matrix.classes[].token`, **bare ID-Token** im Körper; Token-Form via `<!-- d-check:status-provenance -->` deklarierbar) |
| `matrix-inactive` | matrix | Referenz auf Dokument mit verbotenem Status |
| `matrix-downward` | matrix | klasseninterner Abwärtsverweis gegen die deklarierte Rangordnung (`order`/`direction: no-downward`) |
| `external-status` | external | HTTP-Status ≥ 400 oder Transportfehler (DNS/Verbindung) |
| `external-timeout` | external | Timeout überschritten |
| `codepath-missing` | codepaths | Ziel eines Inline-Code-Pfads existiert nicht |
| `citation-out-of-range` | codepaths, citations | Zeilen-Referenz (`datei:<von>-<bis>`) hinter dem Datei-Ende — oder (bei `citations`) die Zieldatei fehlt; nur bei `codepaths.check-lines` bzw. `d-check:cite` |
| `citation-inverted-range` | codepaths, citations | Zeilen-Referenz invertiert (`<von> > <bis>`) |
| `citation-mismatch` | citations | per `d-check:cite` markiertes Zitat ist kein zusammenhängender Teilstring der whitespace-normalisierten Quell-Spanne (Zitat-Fäule) |
| `hostpath-forbidden` | hostpaths | host-lokaler absoluter Pfad (Maschinen-Layout-Leak) in Prosa oder Inline-Code |
| `diagram-id-undefined` | diagrams | Kennung in einem geöffneten Diagramm-Fence ohne Definition in ihrer `defined-in`-Quelle |
| `span-unclosed` | spans | ungeschlossene Code-Span-Öffnung klebt an Nicht-Whitespace (Absatz-Parität gekippt) |
| `fence-unclosed` | spans | Fenced-Code-Block ohne Schluss bis zum **Dateiende** (nicht Absatz-Ende — ein Fence *ist* eine Absatzgrenze): mindestens eine der beiden Fence-Lesarten endet offen, und jedes Modul, das ihr folgt, überspringt alles dahinter. Befund an der Öffnungszeile — eine **Fundstelle**, nicht zwingend die Reparaturstelle —, genau einer je Datei |
| `span-nested-link` | spans | Link-Syntax im Linktext eines weiteren Links (rendert zerrissen) |
| `external-redirects` | external | mehr als `REDIRECT_MAX` Redirects |
| `version-stale` | versions | Versions-Pin weicht von der aktuellen Version (`versions.current-from`) ab |
| `link-stale` | pins | normalisierter Ziel-Span eines gepinnten Links weicht vom hinterlegten `dpin`-Hash ab |
| `core-drift` | immutable | normalisierter Core einer gepinnten Datei (ohne Marker-Zeile + `exclude-sections`) weicht vom hinterlegten `immutable`-Hash ab |
| `core-drift-vcs` | vcs | Core einer immutablen Datei (BASE erfüllt `vcs.immutable-when`) hat sich über die Commit-Range geändert, ihr Status-Übergang ist unzulässig (`vcs.head-allow`), oder die immutable Datei wurde gelöscht/umbenannt |
| `commit-untraceable` | commits | bereinigte Commit-Message trägt keine Kennung nach `commits.id-patterns` und ist nicht per `commits.exempt-pattern` (Betreff) ausgenommen |
| `planning-drift` | planning | Roadmap-Aktiv-Status (`planning.marker` im `planning.heading`-Block) und Präsenz von `planning.slice-glob`-Slices sind inkonsistent (`hasActive ≠ hasSlices`), oder die kanonische Überschrift fehlt/ist mehrdeutig bzw. die Roadmap-Datei fehlt (fail-closed) |
| `closure-note-missing` | planning | Kandidat im `planning.closure.dir` (Filter: `planning.closure.glob`, sonst `planning.slice-glob`) **ohne** einen auf `planning.closure.heading-pattern` passenden Abschnitt — oder das gesetzte `planning.closure.dir` fehlt, ist unlesbar oder enthält **keinen** Kandidaten unter dem effektiven Filter (fail-closed); schließt `closure-note-thin`/`-boilerplate` aus (ohne Abschnitt gibt es nichts zu messen) |
| `closure-note-thin` | planning | Closure-Notiz-Abschnitt trägt weniger als `planning.closure.min-sentences` Satzende-Zeichen **außerhalb** der Fenced-Code-Blöcke (Platzhalter, Einzeiler) |
| `closure-note-boilerplate` | planning | bereinigter Closure-Notiz-Text enthält (case-insensitiv, an Wortgrenzen) eine literale Phrasg aus `planning.closure.boilerplate`; der erste Treffer benennt die Meldung |
| `closure-note-placeholder` | planning | Closure-Notiz-Abschnitt trägt einen unausgefüllten Vorlagen-Platzhalter in Auszeichnungs-Form (opt-in über `planning.closure.placeholder`); Inline-Code, Autolinks/Adressen und HTML-Tags sind ausgenommen, gemeldet wird der **erste** Treffer je Kandidat |
| `closure-note-ambiguous` | planning | mehrere auf `planning.closure.heading-pattern` passende Überschriften — ohne eindeutigen Abschnitt wird **nicht** gemessen (schließt `-thin`/`-boilerplate`/`-placeholder` aus); `line` = **zweiter** Treffer |
| `link-position-dependent` | links | relativer Verweis einer Datei in einem `resolve-from`-`dirs`-Verzeichnis löst von mindestens einem Ort der Gruppe nicht auf — oder von verschiedenen Orten auf **verschiedene** Ziele; Reparatur ist das Präfixieren des Pfads. Fail-closed über denselben Code: eine Gruppe ohne einen einzigen existierenden `dirs`-Ort und ein Ort, der als Datei existiert |
| `wave-drift` | planning | Aktiv-Status der Roadmap (`planning.heading`-Block) und Präsenz eines **flachen** Wellendokuments (`planning.waves.glob` abzüglich `results-glob`) widersprechen sich; auch bei **mehr als einem** flachen Dokument. Fail-closed über denselben Code: unlesbares `waves.dir`/`done-dir` und fehlende Register-Überschrift |
| `wave-preview-exists` | planning | eine Zeile des Vorschau-Registers (`planning.waves.next-heading`) nennt in ihrer **ersten Spalte** eine Welle, für die bereits eine Datei existiert (flach oder im Ruheort) — die geplante Welle hätte drei Positionen statt zwei |
| `wave-results-missing` | planning | eine Zeile des Abschluss-Registers (`planning.waves.closed-heading`) nennt eine Welle **ohne** Ergebnisnotiz (`planning.waves.results-glob`) im Ruheort |
| `wave-unregistered` | planning | eine **Ergebnisnotiz** im Ruheort hat **keine** Zeile im Abschluss-Register — die Richtung „Artefakt ⇒ Register" |
| `section-missing` | structure | kein Abschnitt passt auf den Selektor der Regel — **oder** die Regel trifft keine Datei (auch nach Abzug von `exempt-paths`, fail-closed); `file` = Datei bzw. Glob, `line` = 1 |
| `section-ambiguous` | structure | Abschnitt kommt mehrfach vor, obwohl `sections: one` genau einen erwartet; Abbruch für **diese** Datei in **dieser** Regel, `line` = zweiter Treffer |
| `section-empty` | structure | bereinigter Abschnitts-Text ohne ein einziges Nicht-Whitespace-Zeichen (`non-empty`) |
| `section-thin` | structure | weniger Satzende-Zeichen als `min-sentences` verlangt |
| `section-oversized` | structure | mehr Task-Items als `max-tasks` erlaubt (Zählung **im Abschnitt**, nicht dateiweit) |
| `section-forbidden` | structure | `forbid-pattern` trifft den bereinigten Abschnitts-Text |
| `section-pattern-missing` | structure | `require-pattern` trifft **nicht** — das Spiegelbild von `section-forbidden` |
| `section-marker-missing` | structure | eine Marke aus `require-all` fehlt (Auszeichnungs-Marke am Zeilen-Anfang, nach optionalem Listen-Marker) |
| `target-untracked` | tracked | aufgelöstes, **existierendes** Link-/Bild-Ziel ist nicht im git-Index getrackt (untracked/gitignoriert) — die Referenz wäre auf jedem frischen Klon `target-missing` |
| `gate-phantom` | targets | in einer Doku-Tabellenzeile als `make X` behauptetes Target ohne zugehörige Makefile-Regel (halluziniertes Gate) |
| `gate-undocumented` | targets | Makefile-Regel (nicht in `targets.exempt-targets`) ohne Deklaration als `make X` in der `targets.authority`-Doku (undokumentiertes Gate) |
| `source-drift` | sources | gepinnte externe Quelle (Marker `source-pin` oder Config `sources[]`) inhaltlich gedriftet — Content-Hash der Roh-Bytes bzw. des `unpack: zip`-Content-Manifests weicht vom hinterlegten `sha256` ab; die Meldung trägt den vollen Ist-`sha256` (Re-Pin-Vorlage) |
| `source-unreachable` | sources | gepinnte externe Quelle nicht materialisierbar (Netzfehler, HTTP-Status ≥ 400, Timeout, > `REDIRECT_MAX` Redirects, Body-/Entpack-Limit oder unter `unpack: zip` kein gültiges Zip) — bewusst getrennt von `source-drift` |

Nutzungs-/Umgebungsfehler (Exit 2) melden auf stderr mit Präfix
`d-check: error:`; Konfigurationsfehler nennen Datei und Zeile.

## 5. Metriken und Tracing-Felder

Keine — d-check ist ein CLI-Tool ohne Telemetrie; außerhalb des
Moduls `external` finden keine Netzwerkzugriffe statt
([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

## 6. Externe Verträge

| System | Version/Stand | Vertrag |
|---|---|---|
| `gopkg.in/yaml.v3` | gepinnt via `go.sum` | striktes Decoding (`KnownFields`); vollständig im Config-Adapter gekapselt |
| GitHub Flavored Markdown (Slug-/Anker-Verhalten) | Referenzverhalten, Stand 2026-06 | [§1, DC-FA-ANCH-001.a](#dc-fa-anch-001a--github-slug-algorithmus) |
| Runtime-Basis-Image distroless/static | Digest-gepinnt | Multi-Stage-Build; nur volle Semver-Tags, kein `latest` |

## 7. Historie

| Datum | Änderung |
|---|---|
| 2026-08-16 | §[`DC-FA-LINK-001.a`](spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion) Schritt 6 nach unabhängigem Review und einem **CI-Realfund** nachgezogen: die Ist-Ort-Vorbedingung ist normativ (die Klasse ist „am Ist-Ort grün", kein Doppelbefund), die fail-closed-Zusage der Gruppen-Orte ist an der git-Realität justiert (ein **einzelner** fehlender Ort ist von einem legitim geleerten Verzeichnis nicht unterscheidbar — gemeldet wird die Gruppe ohne einen einzigen existierenden Ort und der Ort, der als Datei existiert), die Scan-Bereichs-Kopplung ist benannt, und die Ventil-Wirkung von §[`DC-FA-REF-001.a`](spezifikation.md#dc-fa-ref-001a--geteiltes-referenz-ventil-ignore-refs) zählt alle fünf unterdrückten Codes |
| 2026-08-16 | §[`DC-FA-LINK-001.a`](spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion) um Schritt **6** + §2-Schema (`links.resolve-from`) erweitert: **ortsfeste Verweise** — je Gruppe wandernder Geschwister-Verzeichnisse wird jedes relative Ziel zusätzlich von jedem Ort der Gruppe aufgelöst; nicht überall auflösbar **oder** nicht überall dasselbe Ziel ⇒ `link-position-dependent`. Dateien in `fixed-dirs` sind keine Quellen (die Bestandsmessung zeigt sonst 108 Falsch-Positive auf Ruheort-Dokumenten). Grund-Code (§4) folgt mit der Implementierung (AllReasons-↔-§4-Lockstep). Begründung in begleitender ADR |
| 2026-08-16 | §[`DC-FA-PLAN-001.a`](spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) um die Schritte **W1–W5** + §2-Schema (`planning.waves.*`) erweitert: die Lifecycle-Invariante eine Ebene höher, die **Wellen**-Abschnitte der Roadmap gegen die Wellen-Dateien. W3 liest den Aktiv-Status **aus Schritt 4**, statt ihn neu zu bestimmen. Zwei Festlegungen kommen aus der Bestandsmessung: das Artefakt einer geschlossenen Welle ist die **Ergebnisnotiz** (gegen das Plan-Dokument geprüft meldet die Aussage 19-mal über zwei Bäume), und die Vorschau-Aussage liest die **erste Spalte** und überspringt Zeilen ohne Kennung (geplante Wellen tragen Namen; die Trigger-Spalte darf andere Wellen nennen). Tabellenzeilen nach derselben Lexik wie §[`DC-FA-TGT-001.a`](spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets) — deren **zweiter** Konsument. Vier Grund-Codes folgen mit der Implementierung. Begründung in begleitender ADR |
| 2026-08-16 | Nachzug nach **dritter** Review-Runde: drei weitere Module beantworteten eine Lexik-Frage roh. §[`DC-FA-VCS-001.a`](spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) — `immutable-when` und die Kopf-Status-Zeile gelten nur **außerhalb** von Fenced-Code (vorher: eine Datei, die ihren Kopf als Beispiel zeigt, galt als immutabel bzw. verschob die gestrippte Zeile ⇒ **stilles Grün** im Immutabilitäts-Gate). §[`DC-FA-PLAN-001.a`](spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) Schritt 4 — die Block-Grenze ist die **geteilte** Abschnitts-Grenze statt eines rohen `## `-Präfix-Vergleichs (eingerückte H2, tab-getrennte H2, H1). §[`DC-FA-TGT-001.a`](spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets) — Tabellenzeilen zählen nur außerhalb von Fenced-Code |
| 2026-08-16 | Nachzug nach bestätigender Re-Review: §[`DC-FA-CITE-001.a`](spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations) Schritt 2 sagt jetzt **beides** — Leerzeilen trennen **nicht** (die Direktive paart mit dem nächsten nicht-leeren Kandidaten), ein Fenced-Block trennt in beiden Zweigen; die Vorfassung erklärte die Fence-Regel mit einem Vergleich, den ein Lauf widerlegt. Ferner: §[`DC-FA-VER-001.a`](spezifikation.md#dc-fa-ver-001a--versions-pin-konsistenz-versions) Schritt 1 und §[`DC-FA-PIN-001.a`](spezifikation.md#dc-fa-pin-001a--content-pin-gegen-inhaltlichen-drift-pins) Schritt 2 sagen die Anker-Antwort jetzt **ganz** zu (Duplikat-Slug, Prozent-Dekodierung, Groß-/Kleinschreibung), nicht nur ihre Fence-Hälfte, und §[`DC-FA-PLAN-001.a`](spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) Schritte 3 und 4 binden den Aktiv-Status-Guard an dieselbe Zeilen-Menge (zwei belegte Falsch-Rot) |
| 2026-08-16 | Drei Konsumenten der geteilten Lexik nachgezogen (Begründung in begleitender ADR): §[`DC-FA-CITE-001.a`](spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations) Schritt 2 nennt den **Absatz** ausdrücklich als den geteilten (Leerzeile **und** Fenced-Block begrenzen) — ein Fenced-Block zwischen Direktive und Zitat trennt und führt in den fail-closed-Fall; §[`DC-FA-VER-001.a`](spezifikation.md#dc-fa-ver-001a--versions-pin-konsistenz-versions) Schritt 1 und §[`DC-FA-PIN-001.a`](spezifikation.md#dc-fa-pin-001a--content-pin-gegen-inhaltlichen-drift-pins) Schritt 2 sagen, dass **beide Anker-Formen nur außerhalb von Fences** erkannt werden (§[`DC-FA-ANCH-001.b`](spezifikation.md#dc-fa-anch-001b--inline-html-anker)), der adressierte bzw. gehashte **Span** aber roh bleibt. Beide Zusagen standen dem Kopfsatz „fence-aware wie die übrigen Module“ nach schon zu; gemessen wurde eine andere Antwort |
| 2026-08-15 | Nachzug aus der Closure des Struktur-Moduls, vier Vertragsflächen: §1 sagt die Symlink-Grenze jetzt für **beide** Formen (auch eine nur über einen **Datei**-Symlink erreichbare Markdown-Datei ist keine Kandidatin — bisher stand dort nur der Verzeichnis-Symlink), §[`DC-FA-STRUCT-001.a`](spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure) Schritt 2 zeigt darauf und sagt zusätzlich den **unlesbaren Dateibaum** fail-closed zu (je Regel ein Befund statt eines leeren Befundsatzes), Schritt 6 nennt den Wort-Begriff der **Marke** ausdrücklich unicode-weit und grenzt ihn gegen die ASCII-Wortgrenze der Floskel aus Schritt C4 ab (dieselbe geteilte Mechanik, zwei Wort-Begriffe — ein Umlaut setzt hier ein Wort fort, dort beendet er eines), und die §2-Schema-Zeile zu `structure[].min-sentences` trägt die Zähl-Semantik aus Schritt 6 („Inline-Code geleert“, „nur vor Whitespace oder Zeilenende“) statt der überholten Kurzfassung |
| 2026-08-15 | §[`DC-FA-STRUCT-001.a`](spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure) Schritte 5 und 6 an die **geteilte** Mechanik angeglichen: die Bereinigung leert auch die Inline-Code-Spans, und ein Satzende zählt nur vor Whitespace oder Zeilenende — beides stand seit Lastenheft 0.55.0 für den Preset-Partner im Vertrag, für `structure` nirgends. Die Folge ist ausdrücklich benannt: ein `forbid-pattern` auf ein Wort in Backticks trifft **nicht** |
| 2026-08-10 | §[`DC-FA-PLAN-001.a`](spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) Schritt C2: YAML-`null` bei `closure.glob` ausdrücklich als **abwesend** ausgewiesen (nicht als leer — die Exit-2-Zusage gilt dem explizit leeren String); §4-Zeile zu `closure-note-missing` nennt den geprüften Gegenstand jetzt „Kandidat“ statt „Slice“ samt effektivem Filter, weil die Kandidaten-Menge seit `planning.closure.glob` nicht mehr aus Slice-Dateien bestehen muss |
| 2026-08-10 | §[`DC-FA-PLAN-001.a`](spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) Schritt C2 + §2-Schema (`planning.closure.glob`): die Closure-Fähigkeit bekommt einen **eigenen** Kandidaten-Filter, dessen Default ein **Verweis** auf `planning.slice-glob` ist (nicht ein kopiertes Literal — ein zweites Muster wäre eine zweite Pflegestelle). Anlass: die beiden Fähigkeiten zählen verschiedene Mengen („noch in Arbeit“ gegen „abgeschlossen“) und teilten sich einen Schlüssel; wer die eine weitet, verbiegt die andere. Ein **explizit** leerer oder ungültiger Glob ⇒ Exit 2 statt stillem Rückfall auf den Default. Kein neuer Grund-Code |
| 2026-08-10 | §2-Schema und §4-Grund-Code-Zeile zu `closure-note-boilerplate` nachgezogen (Wortgrenzen statt Teilstring) sowie Schritt C4b: er sagte „die Zählung aus C4 bleibt unberührt, sie sieht Inline-Code weiterhin“ — seit der Zähl-Angleichung falsch, beide lesen denselben einmal bereinigten Text |
| 2026-08-10 | §[`DC-FA-PLAN-001.a`](spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) Schritt **C4**: die Floskel-Prüfung vergleicht an **Wortgrenzen** statt als Teilstring (Wortzeichen = ASCII `0-9A-Za-z_`, Zeilenränder sind Grenzen). Anlass: kurze Phrasen waren unbrauchbar — `ok` traf auch in *dokumentiert* und *Protokoll*; am eigenen Bestand 68 Treffer, davon einer echt. Mehrwortige Phrasen sind verhaltensgleich. Kein neuer Grund-Code |
| 2026-08-10 | §[`DC-FA-PLAN-001.a`](spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) Schritt **C4** angeglichen (Parität zum abzulösenden Adopter-Skript): der Abschnitt wird **einmal** bereinigt — Fences **und** Inline-Code-Spans —, und alle Bedingungen lesen diesen einen Text; ein Satzende zählt nur vor **Whitespace oder Zeilenende**. Gemessen am eigenen Bestand tragen mehr als die Hälfte aller Satzende-Vorkommen keinen Satz (Link-Pfade, Versionsnummern). Die Änderung wirkt in **zwei** Richtungen: `closure-note-thin` wird **schärfer**, `closure-note-boilerplate` **lockerer** (eine zitierte Floskel in Backticks trifft nicht mehr) — die Lockerung ist begleitend per ADR entschieden. C4b liest damit denselben Text und leert Inline-Code nicht mehr selbst |
| 2026-08-10 | §[`DC-FA-PLAN-001.a`](spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) nach unabhängigem Review nachgezogen: **C5 war auf der Fassung vor C4b stehengeblieben** und widersprach ihr (Zeilen-Angabe, Kombinierbarkeit) — C5 nennt jetzt die C4b-Ausnahme („Zeile des Treffers“) und alle drei kombinierbaren Bedingungen; C4b sagt umgekehrt, dass der Treffer in der `message` steht und `target` das Verzeichnis bleibt. Fachlich verengt: das Innere eines Platzhalters muss **frei von Whitespace** sein, und ein Winkelklammer-Linkziel ist ein dritter Nachfilter |
| 2026-08-10 | §[`DC-FA-PLAN-001.a`](spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) um Schritt **C4b** + §2-Schema (`planning.closure.placeholder`) erweitert: opt-in Erkennung unausgefüllter Vorlagen-Platzhalter. Auf dem C4-Abschnittstext werden zusätzlich die **Inline-Code-Spans geleert** — am eigenen Bestand gemessen liegen alle zwölf Treffer dort und keiner außerhalb, ohne die Einschränkung wäre jeder ein Falsch-Positiv. Zwei Nachfilter (Autolink/Adresse, HTML-Tag) sind ausdrücklich **Code**, nicht Teil des Musters. Erster Treffer je Kandidat; die Zählung aus C4 bleibt unberührt. Die Konsumenten-Vorlage nutzt Lookarounds und ist in RE2 nicht ausdrückbar — portiert durch Konsumieren des Vorzeichens. Grund-Code (§4) folgt mit der Implementierung (AllReasons-↔-§4-Lockstep) |
| 2026-08-10 | §4-Grund-Code-Zeile zu `fence-unclosed` nachgezogen: sie trug weiter die mit Lastenheft 0.52.1 widerrufene Reichweite („von **allen** Modulen übersprungen. Befund an der Öffnungszeile") und widersprach damit der Anforderung — ein Lauf widerlegte sie in derselben Ausgabe. Jetzt „mindestens eine Lesart endet offen" und **Fundstelle** statt Reparaturstelle. Der `--doctor`-Klartext desselben Grund-Codes ist mitgezogen |
| 2026-08-10 | §[`DC-FA-SPAN-001.a`](spezifikation.md#dc-fa-span-001a--span-artefakt-erkennung) Schritt 3 nach bestätigender Re-Review nachgezogen: die Trimmung ist als **geteiltes** Prädikat festgeschrieben (das Modul `planning` trimmte weiter unicode-weit und trug den Anlassfall des Slice unverändert weiter), das Befund-Ziel verliert auch das CR einer CRLF-Zeile, und die Fundstellen-Klausel gilt für **beide** Lesarten — auch die strenge zeigt daneben, wenn eine längere Fence-Zeile eine kürzere Öffnung geschlossen hat |
| 2026-08-09 | §[`DC-FA-SPAN-001.a`](spezifikation.md#dc-fa-span-001a--span-artefakt-erkennung) Schritt 3 nach unabhängigem Code-Review überarbeitet: der Algorithmus führt jede Zeile durch **beide** Schluss-Lesarten des Produkts (naiver Toggle und längenabgeglichener CommonMark-Schluss) und meldet, sobald **eine** von beiden am Dateiende offen steht — vorher nur die naive, wodurch der Tabellen-Leser aus §[`DC-FA-REQ-001.a`](spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen) unbewacht blieb. Die Trimmung ist auf Space/Tab festgeschrieben (identisch zur Vorverarbeitung, nicht unicode-weit); die Zeilenwahl ist präzisiert (Öffnungszeile der strengen Lesart, sonst letzte öffnend gewertete Zeile — unter reiner Paritätskippung ist die fehlende Öffnung nicht bestimmbar); die Kappung zählt **Runen** |
| 2026-08-09 | §[`DC-FA-SPAN-001.a`](spezifikation.md#dc-fa-span-001a--span-artefakt-erkennung) um eine **dritte** Artefakt-Klasse erweitert: `fence-unclosed` — eine Fence-Öffnung ohne Schluss bis zum **Dateiende**. Dateiweit statt absatzweise (ein Fence *ist* eine Absatzgrenze), Befund an der Öffnungszeile, genau einer je Datei. Anlass ist ein belegter **stiller Grün-Pfad**: hinter einem offenen Fence überspringt jede Vorverarbeitung den Rest, und ein Gate meldet grün, ohne geprüft zu haben. **Bewusst nicht** mitgelöst wird die Fence-**Paarung** (längenabgeglichener CommonMark-Schluss) — sie ist auf dem gemessenen Bestand wirkungslos und löst den Fall ohnehin nicht, weil ein nie geschlossener Fence unter jeder Paarungsregel offen bleibt. Grund-Code (§4) folgt mit der Implementierung (AllReasons-↔-§4-Lockstep). Begründung (Modul-Wahl, Reichweite, Nicht-Lösung der Paarung) in begleitender ADR |
| 2026-08-09 | Neue §[`DC-FA-STRUCT-001.a`](spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure) + §2-Schema (`structure[]`) — **Struktur-Invarianten innerhalb eines Dokuments** (Modul `structure`, opt-in, hermetisch, Post-Pass): je Regel eine Dokumentklasse über **eigene** Globs (unabhängig vom globalen Scan-Scope, daher kein `<modul>.scope`), ein Abschnitts-Typ (`section` Klartext **oder** `section-pattern` RE2) und ein **Kardinalitäts-Modus** `sections`: `one` (Default; mehrere Treffer ⇒ `section-ambiguous` und Abbruch für diese Datei) oder `each` (jeder Treffer wird geprüft — für **wiederkehrende** Abschnitte). Null Kandidaten ⇒ `section-missing`, auch wenn erst `exempt-paths` die Menge geleert hat. Nach Fence-Bereinigung sechs Bedingungen mit **je eigenem** Grund-Code (`section-empty`/`-thin`/`-oversized`/`-forbidden`/`-pattern-missing`/`-marker-missing`), weil die Befund-Deduplikation zwei Verletzungen desselben Abschnitts sonst zusammenfallen ließe. Marken sind **Auszeichnungs-Marken** am Zeilen-Anfang nach optionalem **Listen-Marker**; `require-pattern` ist das Spiegelbild von `forbid-pattern` und deckt Aussagen **innerhalb** einer Auszeichnung. §[`DC-FA-PLAN-001.a`](spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) Schritt C3 um dieselbe Mehrdeutigkeits-Härte ergänzt (`closure-note-ambiguous`) und als **Preset** im Modus `one` ausgewiesen — beide teilen Abschnitts-Findung, Kardinalitäts-Behandlung, Bereinigung und Zählung; eine Änderung an nur einer Stelle ist ein Spezifikations-Bug. Grund-Codes (§4) folgen mit der Implementierung (AllReasons-↔-§4-Lockstep). Begründung in begleitender ADR |
| 2026-08-09 | §[`DC-FA-PLAN-001.a`](spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) um die **Closure-Note-Struktur** (Schritte C1–C5) erweitert + neue §[`DC-FA-CLI-012.a`](spezifikation.md#dc-fa-cli-012a--konfigurations-pfad---config) + §2-Schema (`planning.closure.dir`/`heading-pattern`/`min-sentences`/`boilerplate`). Die Closure-Fähigkeit ist opt-in **innerhalb** des opt-in Moduls (leere `dir` ⇒ keine Slice-Datei wird geöffnet, byte-identisch): erster auf `heading-pattern` passender Abschnitt (bis zur nächsten gleich-/höherrangigen Überschrift), Fenced-Code entfernt, Satzende-Zeichen gezählt (< `min-sentences` ⇒ `closure-note-thin`), Rest case-insensitiv gegen die literalen `boilerplate`-Teilstrings (⇒ `closure-note-boilerplate`); kein passender Abschnitt bzw. fehlendes gesetztes Verzeichnis ⇒ `closure-note-missing`. Leeres Closure-Verzeichnis ist befundfrei, `min-sentences` < 1 und leerer Floskel-Eintrag ⇒ Exit 2; diagnose-only. `--config <datei>` verschiebt **nur die Herkunft** der Konfiguration (Wurzel-Constraint + Existenz fail-closed, ersetzt statt ergänzt, Validierung und Semantik unverändert, Befund-Provenance nennt die geladene Datei) — es trennt die Prüf-Profile zweier Bindepunkte, ohne die Modulwahl auf der Kommandozeile nachzubauen. Grund-Codes `closure-note-missing`/`closure-note-thin`/`closure-note-boilerplate` in §4 (im Lockstep mit der Implementierung ergänzt). Begründung (Modul-Schnitt statt neues Modul, Bindepunkt-Trennung, Struktur-vs-Bedeutung-Grenze, Schwellenwahl) in begleitender ADR |
| 2026-07-19 | Neue §[`DC-FA-SRC-001.a`](spezifikation.md#dc-fa-src-001a--upstream-content-drift-externer-quellen-sources) + §2-Schema (`sources[]`: `url`/`sha256`/`unpack`) — **Upstream-Content-Drift** (Modul `sources`, opt-in, **Netz**): eine auf `sha256` gepinnte externe Quelle (Marker `source-pin` am Link **oder** Config `sources:`) wird geholt, gehasht, verglichen; Abweichung → `source-drift` (voller Ist-Hash), Fetch-Fehler → `source-unreachable`. Einzeldatei (Roh-Byte-Hash) oder Archiv (`unpack: zip` → `LC_ALL=C`-sortiertes, reihenfolge-invariantes Content-Manifest). Post-Pass wie `external`; zweite Netz-Tür ([`DC-QA-03`](lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)). Grund-Codes `source-drift`/`source-unreachable` (§4) folgen mit der Modul-Implementierung (AllReasons-↔-§4-Lockstep). Begründung (Netz-Modul-Design, Amendment der Netz-Sparsamkeit, Manifest-Hash, Marker + Config) in begleitender ADR |
| 2026-07-18 | §[`DC-FA-CODE-001.a`](spezifikation.md#dc-fa-code-001a--pfade-in-inline-code) Schritt 3/6 + neue §[`DC-FA-CITE-001.a`](spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations) + §2-Schema (`codepaths.check-lines`) — **Zitat-Verifikation** in zwei Teilen: (1) opt-in `codepaths.check-lines` verifiziert die Zeilen-Referenz eines `datei:<von>-<bis>`-Pfads (Ziel existiert + ≥ `<bis>` Zeilen ⇒ sonst `citation-out-of-range`; `<von> ≤ <bis>` ⇒ sonst `citation-inverted-range`), das Suffix wurde bisher abgetrennt und verworfen, default-aus **byte-identisch**. (2) neues Modul `citations`: die Direktive `<!-- d-check:cite <pfad>:<von>-<bis> -->` vor einem `>`-Zitatblock ⇒ **zeichengenauer** Vergleich der Zeilen `<von>`–`<bis>` gegen den Zitatblock (`citation-mismatch`), hermetisch, fail-closed (fehlende Datei/Spanne/ungültiger Bereich ⇒ Exit 2), greift die `codepaths`-`datei:zeile`-Erkennung auf. Grund-Codes (§4) folgen mit der Modul-Implementierung (AllReasons-↔-§4-Lockstep). Begründung (Erweiterung vs. eigenes Modul, `verbatim`/Direktiven-Design) in begleitender ADR |
| 2026-07-18 | §[`DC-FA-REF-001.a`](spezifikation.md#dc-fa-ref-001a--geteiltes-referenz-ventil-ignore-refs) — neues **geteiltes Referenz-Ventil** `ignore-refs` mit **Quell-Skopus** (`in`) und Zwei-Feld-Semantik (`refs ∧ ¬keep`; `keep` reihenfolge-unabhängig, kein gitignore-Last-Match), honoriert von `links`/`anchors`/`codepaths`. §[`DC-FA-CODE-001.a`](spezifikation.md#dc-fa-code-001a--pfade-in-inline-code) Referenz-Ventil **und** Schritt 5 auf das geteilte Ventil umgestellt; die modul-lokale Liste `codepaths.ignore-refs` bleibt **Alias** (byte-identisch). §[`DC-FA-LINK-001.a`](spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion) und §[`DC-FA-ANCH-001.a`](spezifikation.md#dc-fa-anch-001a--github-slug-algorithmus) honorieren es (Existenz-/Escape-/Anker-Prüfung übersprungen, Symlink-Prüfung bleibt). §2-Schema: Top-Level `ignore-refs[].in`/`refs`/`keep` + `codepaths.ignore-refs` als Alias annotiert. **CR** (Konsument `ai-harness-course`); Begründung (Zwei-Feld vs. `!`-Negation, Alias-Pfad) in der begleitenden ADR |
| 2026-07-18 | §[`DC-FA-REQ-001.a`](spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen) Schritt 5 — **Direktiven-Toleranz** ergänzt: eine Datenzeile mit **genau einer** überzähligen, ganzzelligen HTML-Kommentar-Zelle (`<!-- … -->` hinter der letzten Pipe) wird auf Header-Breite gelesen (GFM-nachsichtig für Body); d-checks eigene `<!-- d-check:ignore … -->`-Konvention in einer Tabellenzeile bricht den Reader nicht mehr ab. Zwei überzählige Zellen oder eine Nicht-Kommentar-Zelle bleiben Exit 2. **Sicherer Aufsatz** auf der Tabellengrenze (Schritt 5): eine relevante Folgetabelle wird nicht mehr verschluckt — der zuvor fünffach reproduzierte stille Übersprung ist strukturell zu. Wirkt über den geteilten Reader zugleich auf §[`DC-FA-XREF-001.a`](spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency). **Defekt-Fix (rot→grün: ein spurioser Exit 2 wird lesbar), kein CR; SemVer-Patch.** |
| 2026-07-18 | §[`DC-FA-REQ-001.a`](spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen) Schritt 5 — **Tabellengrenze am relevanten Header** ergänzt: bildet eine Zeile mit ihrer Folgezeile einen gültigen Header + Trennzeile und **bindet ihr Header eine Rolle** (Schritt 4), beendet sie die laufende Tabelle — auch bei passender Zellenzahl. Damit wird jede relevante Tabelle erkannt und kann nicht mehr still in einer vorangehenden (irrelevanten) verschwinden; ein rollenloser Header (z. B. all-dashes) beendet **nicht** (Gegenprobe `fx-t`). Wirkt über den geteilten Reader zugleich auf §[`DC-FA-XREF-001.a`](spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency). **Defekt-Fix, kein CR** (das Lastenheft definiert keine Tabellengrenze), aber **SemVer-Minor**: d-check findet danach **mehr**. Belegt gegen das ausgelieferte v0.47.0-Image (`fx-s`: still `1 Waise` statt 3). |
| 2026-07-17 | §[`DC-FA-COV-001.a`](spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage) Schritt 3 — **Komma-Kurzform-Regel geschärft** (bei der Implementierung): sie greift nicht nur direkt hinter der Kennung, sondern **auch hinter einer konsumierten Range/Enum** (`<FAM>-AAA..CCC, DDD`). Anlass: der reale Auslöser in grid-gyms `traceability.md` ist `GG-SCN-001..005, 007, 008` — die enge Erst-Formulierung („folgt der Fundstelle ein Komma") ließ `007, 008` **still** fallen (Exit 0, 3 Waisen, gemessen), obwohl genau dieser stille Drop der Zweck des Slice war. Ein Komma vor einer vollständigen Kennung bleibt auch nach einer Range unberührt. |
| 2026-07-17 | **Markdown-Lexik an CommonMark/GFM angeglichen**, zwei Regeln: §[`DC-FA-REQ-001.a`](spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen) Schritt 3 — Trennzelle `^:?-{3,}:?$` → `^:?-+:?$` (GFM verlangt **einen** Bindestrich, wir verlangten drei; jede reale Tabelle mit `\| -- \|` war für d-check keine Tabelle). §[`DC-FA-LINK-001.a`](spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion) Schritt 1 — **Infozeilen-Regel**: eine ` ``` `-Zeile mit Backtick im Rest ist kein Fence-Öffner (CommonMark), sondern Fließtext; ohne sie blendet ein Satz **über** einen Fence **alle** Module bis zum Dateiende (Exit 1 ⇒ Exit 0, gemessen). Beides **still** und **ausgeliefert**; belegt per Differential-Spike gegen goldmark v1.8.4 über 522 reale Dateien (490 Tabellen ⇒ 8 Abweichungen, alle „d-check ist blind"). **Defekt-Fix, kein CR** (das Lastenheft sagt weder, was eine Trennzeile ist, noch was einen Fence öffnet), aber **SemVer-Minor**: d-check findet danach **mehr**. |
| 2026-07-17 | §[`DC-FA-COV-001.a`](spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage) Schritt 3 um die **Komma-Kurzform** ergänzt: Kennung + Komma + Ziffern ⇒ Exit 2 mit Hinweis auf die zugesagten Notationen, statt stillem Drop oder geratener Expansion. Komma vor einer vollständigen Kennung bleibt unberührt. Lastenheft-CR 0.46.0; Begründung in der begleitenden ADR |
| 2026-07-17 | §[`DC-FA-XREF-001.a`](spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency) Schritt 2 + §2-Schema um **`forward.req-pattern`** ergänzt (Default `trace.requirements.id-pattern`, symmetrisch zu `backward.req-pattern`). Die Vorwärts-Sicht las ihre IDs bis dahin **still** über das RTM-Muster — die Kopplung war nirgends ausgesprochen. Festgehalten: die **Vergleichs-Schlüsselmenge ist nicht die RTM-Anforderungsmenge** (das Muster entscheidet, nicht die RTM-Mitgliedschaft). Lastenheft-CR 0.45.0. |
| 2026-07-17 | §[`DC-FA-COV-001.a`](spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage) Schritt 3 (Range-Parser) um **Link-Transparenz** geschärft: die Fortsetzung `..NNN`/`/NNN` darf **genau einmal** durch ein Markdown-Link-Suffix `](…)` unterbrochen sein, dahinter gilt wieder „unmittelbar“; weitergehendes Peeling (Whitespace, Emphasis, zweites Suffix) bleibt ausgeschlossen. Wirkt über den geteilten Parser zugleich auf §[`DC-FA-XREF-001.a`](spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency). **Defekt-Fix, kein CR:** das Lastenheft verspricht die Expansion unqualifiziert — die Verengung „unmittelbar“ stand allein hier und kollidierte strukturell mit der Linkpflicht ([`DC-FA-ID-001`](lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)); betroffen war `trace.coverage` seit v0.41.0 (verlinkte Range ⇒ falsche Waisen). |
| 2026-07-17 | §[`DC-FA-XREF-001.a`](spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency) um die **Vakuitäts-Stufe** (Schritt 5) geschärft: ein Abgleich ohne eine einzige Kante ist Exit 2 statt `0 Differenz(en)`/Exit 0 — beide Sichten kantenleer (geteiltes `design-pattern` greift am Namensraum vorbei) oder Rück-Sicht kantenleer unter `mode: superset`. Eine **einseitig** leere Vorwärts-Sicht bleibt wohldefiniert (Diff über `keys(F) ∪ keys(B)`) und meldet `B \ F` laut. Vakuität wird **nach** dem Ausschluss (Schritt 4) gemessen — ein `exclude-req`, das alle Anforderungen verschluckt, schaltet das Gate ebenso still ab wie ein fehlgreifendes Muster. Fehlerpräzedenz um die Abschnitts-Spannungs- und die Ausschluss-Stufe ergänzt und als **stufenweise über beide Sichten** präzisiert. Lastenheft-CR 0.44.1. |
| 2026-07-16 | §[`DC-FA-XREF-001.a`](spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency) + §2-Schema (`trace.cross-consistency.*`) ergänzt: opt-in Mengenabgleich zweier Traceability-Sichten — Vorwärts-RTM-Tabelle (Anforderung→Design) gegen Rückwärts-`Bezug`-Kanten (Design→Anforderung), je Anforderung `F\B`/`B\F` mit Richtung, Modus `equal`/`superset`. Beide über den header-gebundenen Reader ([§`DC-FA-REQ-001.a`](spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)) + range-aware Span-Semantik ([§`DC-FA-COV-001.a`](spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)); Rück-Artefakt-ID = erste Spalte via `design-pattern`, `exclude-req`-Ventil (RE2). Advisory unter `--trace`; Gatung über globales `--require-complete` ([`DC-FA-CLI-011`](lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)). Fail-closed, ohne Block byte-identisch ([`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)). |
| 2026-07-14 | §[`DC-FA-REQ-001.a`](spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen) + §2-Schema (`trace.requirements.format`/`table.*`) ergänzt: `headings` bleibt Default, `table` extrahiert Markdown-Pipe-Tabellen außerhalb von Fences über exakt benannte ID-/Text-/optionale Modalitätsheader; `\|` und Pipes in einzeiligen Code-Spans teilen keine Zelle. Beide Extraktoren liefern dasselbe `{id,title,modalityText}`-Modell vor Referenzscan/Gating. Nichtleer explizite Quelle oder Tabellenmodus aktiviert fail-closed (fehlende Quelle, Header-/Zeilenstruktur, Duplicate-ID, null Treffer ⇒ Exit 2); `source: ""` und unkonfigurierter Heading-Pfad behalten Empty⇒Default/Deduplizierung/leere RTM byte-identisch. Modalität klassifiziert im Tabellenmodus die konfigurierte Modalitätsspalte, sonst die Textspalte. |
| 2026-07-11 | §[`DC-FA-CLI-009.a`](spezifikation.md#dc-fa-cli-009a--requirements-traceability-matrix) Schritte 1/2 + Config-Auflösungs-Absatz + §2-Schema (`trace.requirements.source`/`.id-pattern`, `trace.adrs.dir`/`.file-pattern`/`.id-prefix`, `trace.slices.dir`/`.file-pattern`/`.id-prefix`) ergänzt: die vier RTM-Konventions-Annahmen ([`DC-FA-CLI-009`](lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)) sind über einen opt-in `trace`-Block überschreibbar — Anforderungs-Quelldatei + Kennungs-Regex (Ganz-Token im Heading, Referenz-Zählung in ADR/Slice) sowie je Referenzklasse Verzeichnis + Basisnamen-Regex (Capture-Gruppe 1 = Owner-Kennung) + Owner-Präfix. Jedes Feld optional, abwesend ⇒ Default (byte-identisch, [`DC-QA-02`](lastenheft.md#dc-qa-02--determinismus)); fail-closed zur Config-Zeit (ungültige Regex bzw. `file-pattern` ohne Capture-Gruppe ⇒ Exit 2). Gilt unverändert für `--require-complete` ([`DC-FA-CLI-011`](lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)). Kein neuer Grund-Code (`--trace` bleibt advisory). Design spiegelt `ids.patterns` (begleitende ADR) |
| 2026-07-11 | §[`DC-FA-COV-001.a`](spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage) + §2-Schema (`trace.coverage[].files`/`label`/`ranges`/`sections`/`exclude-sections`) ergänzt: dritte opt-in RTM-Referenzklasse `trace.coverage` ([`DC-FA-COV-001`](lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)) liest kuratierte Matrizen range-aware als Coverage — Abschnitts-Whitelist/Blacklist über die bestehende `matrix`-Span-Semantik (`excludedRanges`), Range-Parser (`<FAM>-AAA..BBB` breiten-erhaltend + `/`-Enum, gegen `id-pattern` validiert, `AAA>BBB`/Breite ⇒ Exit 2), Waise = ¬slice ∧ ¬coverage. Konditionale Coverage-Spalte + `coverage`-Feld (omitempty) ⇒ ohne Quelle byte-identisch. Mit-Modifikation §[`DC-FA-CLI-009.a`](spezifikation.md#dc-fa-cli-009a--requirements-traceability-matrix) (Spalte) + [`DC-FA-CLI-011`](lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code) (Waisen-Definition). `--print-config` führt den `coverage`-Block. Fail-closed |
| 2026-07-11 | §[`DC-FA-MOD-001.a`](spezifikation.md#dc-fa-mod-001a--modalitäts-klassifikation-tracerequirementsmodality) + §2-Schema (`trace.requirements.modality.levels`/`require-levels`) ergänzt: opt-in Modalitäts-Klassifikation ([`DC-FA-MOD-001`](lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in)) — je Anforderung die RFC-2119-Stufe aus konfigurierbaren Modal-Verb-Keywords (DE+EN-Defaults, erster/längster Treffer im Body-Abschnitt via Span-Mechanik, `unknown`-Fallback), konditionale Modality-Spalte; `--require-complete` gatet nur `require-levels`-Stufen (Default `[must]`, `unknown` explizit). Byte-identisch ohne Block; fail-closed (leerer Stufen-Name/Keyword, `require-levels` weder Stufe noch `unknown` ⇒ Exit 2). Mit-Modifikation §[`DC-FA-CLI-009.a`](spezifikation.md#dc-fa-cli-009a--requirements-traceability-matrix) (Spalte) + [`DC-FA-CLI-011`](lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code) (gatende Stufen). `--print-config` führt den `modality`-Block |
| 2026-07-05 | §[`DC-FA-TGT-001.a`](spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets) + §2-Schema (`targets.makefiles`/`doc-tables`/`authority`/`exempt-targets`) ergänzt: opt-in Modul `targets` prüft **hermetisch** (Filesystem-Port `ReadFile`, **kein** git/Netz/Makefile-Ausführen) die Doku-↔-Makefile-Deklarations-Konsistenz — ein nur aus **Tabellenzeilen** (Pipe-Präfix, keine Prosa) dokumentiertes `make X` ohne Makefile-Regel ⇒ `gate-phantom`, eine Makefile-Regel (nicht `targets.exempt-targets`) ohne Eintrag in `targets.authority` ⇒ `gate-undocumented`; statische Zeilen-Heuristik für Regelnamen (keine Pattern-Rules/variablen Targets, in Parität zu `tools/gate-consistency.sh`), fail-closed bei fehlender konfigurierter Datei, diagnose-only, default-aus byte-identisch; Analogie zu `planning` (Doku-Behauptung ↔ Repo-Struktur, hermetisch). Außerdem §[`DC-FA-CLI-010.a`](spezifikation.md#dc-fa-cli-010a--makefile-fragment) (10→11 Targets): `--print-mk` trägt `doc-targets` (`--enable targets`, hermetisch ohne Range). Die Grund-Codes `gate-phantom`/`gate-undocumented` (§4) folgen mit der Modul-Implementierung (AllReasons-↔-§4-Lockstep) |
| 2026-07-03 | Review R1/R2 zu §[`DC-FA-TRK-001.a`](spezifikation.md#dc-fa-trk-001a--getrackt-status-auflösbarer-referenz-ziele-tracked), präzisiert: Prüfung **je gescannter Quell-Datei** statt „Post-Pass" (Auflösungs-Mechanik unabhängig von der Aktivierung des Moduls `links` — ein fokussierter Lauf prüft vollständig); Verzeichnis-Ziele und Symlink-Referenzen (Ziel ist/durchläuft einen Symlink) explizit kein Kandidat ([`DC-FA-LINK-002`](lastenheft.md#dc-fa-link-002--symlink-ablehnung)-Domäne — false-positive hinter getrackten Verzeichnis-Symlinks vermieden); `target` = aufgelöster Pfad bekräftigt (Ventil-Parität); `exempt-targets` segmentweise validiert (Exit 2) |
| 2026-07-03 | §[`DC-FA-TRK-001.a`](spezifikation.md#dc-fa-trk-001a--getrackt-status-auflösbarer-referenz-ziele-tracked) + §2-Schema (`tracked.exempt-targets`) + Grund-Code `target-untracked` (§4) ergänzt: opt-in Modul `tracked` prüft die Datei-Ebene der von `links` aufgelösten, **existierenden** repo-internen Ziele gegen den **git-Index** — dritte VCS-Port-Nutzung (`vcs` Range-Diff, `commits` Messages, `tracked` Index), **ohne** Range/`--staged`; Index statt `.gitignore`-Interpretation (gestagte Dateien gelten als getrackt), kein Doppelbefund (`target-missing` bleibt `links`), Ventil referenz-weit analog `codepaths.ignore-refs`; fail-closed ohne lesbares `.git` (Exit 2), diagnose-only, default-aus byte-identisch. Außerdem §[`DC-FA-CLI-010.a`](spezifikation.md#dc-fa-cli-010a--makefile-fragment) (9→10 Targets): `--print-mk` trägt `doc-tracked` (`--enable tracked` + fokussierte `--disable`-Liste, ohne Range) |
| 2026-07-01 | §[`DC-FA-PLAN-001.a`](spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) + §2-Schema (`planning.roadmap`/`heading`/`marker`/`slice-glob`) + Grund-Code `planning-drift` (§4) ergänzt: opt-in Modul `planning` prüft **hermetisch** (nur Arbeitsbaum, kein git/Netz) die Roadmap-↔-in-progress-Invariante — Ruhe-Marker im `## Aktuelle Welle`-Block genau dann, wenn kein `slice-*` im Verzeichnis (`hasActive == hasSlices`), sonst `planning-drift`; Heading-Guard fail-closed, Post-Pass wie `codepaths`-Existenz, diagnose-only. Default-aus byte-identisch. Portierung des abgelösten `tools/planning-consistency.sh`. Außerdem §[`DC-FA-CLI-010.a`](spezifikation.md#dc-fa-cli-010a--makefile-fragment) (8→9 Targets): `--print-mk` trägt `doc-planning` (`--enable planning`, hermetisch ohne `--range`) |
| 2026-07-01 | §[`DC-FA-COMMITS-001.a`](spezifikation.md#dc-fa-commits-001a--traceability-kennung-in-commit-messages-über-eine-commit-range-commits) + §2-Schema (`commits.id-patterns`/`commits.exempt-pattern`) + Grund-Code `commit-untraceable` (§4) ergänzt: opt-in Modul `commits` prüft, dass jede Commit-Message eine Kennung nach `commits.id-patterns` trägt (`commit-untraceable`), über **denselben VCS-Port** wie `vcs` (reine-Go-git, ohne git-Binary, ohne Netz), erweitert um Message-Lesen; zwei Quellen: `--range <base>..<head>` (Nicht-Merge-Commits, `--no-merges`-Parität) und `--commit-msg <datei\|->` (Kurzschluss-Modus für den commit-msg-Hook, einzelne Pending-Message). Uniforme `#`-/scissors-Bereinigung (Kennung auf Inhalts-Zeile), Betreff-Ausnahme `commits.exempt-pattern`; fail-closed ohne `.git`/Range/Message-Datei, diagnose-only. Default-aus byte-identisch. Portierung des abgelösten `tools/trace-check.sh` (dieselbe VCS-Port-Präzedenz wie `vcs`, auf Commit-Messages statt Datei-Inhalt). Außerdem §[`DC-FA-CLI-010.a`](spezifikation.md#dc-fa-cli-010a--makefile-fragment) (7→8 Targets): `--print-mk` trägt `doc-commits` (`--enable commits` + aus `ValidModules` abgeleitete Fokus-`--disable`-Liste, `--range`) — verteilt die Commit-Traceability an Konsumenten |
| 2026-06-29 | §[`DC-FA-CODE-001.a`](spezifikation.md#dc-fa-code-001a--pfade-in-inline-code) + §2-Schema (`codepaths.ignore-refs`) ergänzt: Referenz-Ventil `ignore-refs` — ein aufgelöster Ziel-Pfad, der ein Glob matcht, wird nicht existenz-/anker-geprüft (**referenz-weit**, Tombstone-Register entfernter Artefakte); dritte Ventil-Achse neben dem zeilenweisen `d-check:ignore` und dem datei-weiten `exempt-paths`, in Schritt 5 vor `codepath-missing`, ohne Eintrag byte-identisch. Anlass: Frozen-Doc-Refactoring-Falle (immutable ADRs zitieren entfernte Pfade) — zugleich das behaltene `adr-immutable-check.sh` entfernt |
| 2026-06-29 | §[`DC-FA-VCS-001.a`](spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) + §2-Schema (`vcs.paths`/`immutable-when`/`exclude-sections`/`status-line`/`head-allow`) + Grund-Code `core-drift-vcs` (§4) ergänzt: opt-in Modul `vcs` vergleicht `core(BASE)` vs. `core(HEAD)` über eine Commit-Range (`--range <base>..<head>` / `--staged`), liest das read-only `.git` über einen reine-Go-VCS-Port (ohne git-Binary, ohne Netz); erweiterter Eingabe-Scope (git + Range), aber lokal/lesend/deterministisch — Determinismus/Read-only gehalten; fail-closed ohne `.git`/Range, diagnose-only. Core-Semantik in Parität zum abgelösten `adr-immutable-check.sh` (nur Kopf-Status-Zeile gestrippt, `exclude-sections`-Abschnitte). Default-aus byte-identisch. Außerdem §[`DC-FA-CLI-010.a`](spezifikation.md#dc-fa-cli-010a--makefile-fragment) (6→7 Targets): `--print-mk` trägt `doc-immutable` (`--enable vcs` + aus `ValidModules` abgeleitete Fokus-`--disable`-Liste, `RANGE`/`STAGED`) — verteilt die git-Garantie an Konsumenten |
| 2026-06-28 | §[`DC-FA-MTX-001.a`](spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung) Schritt 6 + §2-Schema (`matrix.classes[].token`, `matrix.exempt-paths`) + §4 (`matrix-forbidden` Token-Form) ergänzt: token-basierte Referenz-Richtung ([`DC-FA-MTX-003`](lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)) — `matrix` fängt verbotene Referenzen auch als bare ID-Token im Prosa-Körper (außer Links/Fences/`exclude-sections`); Provenance-Marker `<!-- d-check:status-provenance -->` auf der rohen Zeile nimmt aus; `exempt-paths` grandfathered ganze Dateien. Fail-closed (`token` kompiliert/Leerstring). Default-aus byte-identisch. Außerdem §[`DC-FA-SPAN-001.a`](lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in): Slice-Token aus dem Spec-Körper entfernt (Provenance gehört in die Historie) |
| 2026-06-28 | §2 „Glob-Auswertung" ergänzt: alle Glob-Felder (`scan.ignore`, `<modul>.scope.ignore`, `matrix.classes[].paths`/`.order`, `*.exempt-paths`) werden segmentweise über Go-`path.Match` ausgewertet (`**` segmentübergreifend); negierte Zeichenklasse `[^…]` (Go), **nicht** `[!…]` (fnmatch). Reine Klarstellung des Bestands (`matchGlob`), kein Verhaltens-/Schema-Change |
| 2026-06-28 | §[`DC-FA-MTX-001.a`](spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung) Schritt 5 + §2-Schema (`matrix.classes[].order`/`.direction`) + Grund-Code `matrix-downward` (§4) + Config-Beispiel ergänzt: klasseninterne Verweisrichtung ([`DC-FA-MTX-002`](lastenheft.md#dc-fa-mtx-002--verweisrichtung-innerhalb-einer-geordneten-dokumentklasse-modul-matrix)) — eine Klasse mit `order` (Glob-Rang, First-Match) + `direction: no-downward` meldet klasseninterne Abwärtsverweise (Rang *i* → *j > i*, auch transitiv) als `matrix-downward`; rangfreie Mitglieder und klassenübergreifende Referenzen ausgenommen; fail-closed-Config (`order`/`direction` nur zusammen, unbekannter `direction`-Wert ⇒ Exit 2); Default-aus byte-identisch |
| 2026-06-24 | §[`DC-FA-PIN-001.a`](spezifikation.md#dc-fa-pin-001a--content-pin-gegen-inhaltlichen-drift-pins) + Grund-Code `link-stale` (§4) ergänzt: opt-in Modul `pins` hasst den whitespace-normalisierten **rohen** Ziel-Span (Datei/Heading-Section inkl. Fenced-Code) eines gepinnten Links (`<!-- dpin: sha256:… -->`, gebunden an den unmittelbar vorausgehenden Link derselben Zeile, sonst inert) und meldet `link-stale` bei Drift; nur auflösbare repo-interne Ziele (struktureller Befund bleibt `links`/`anchors`, kein Doppelbefund), Scope-treu (nur Quell-Dateien), diagnose-only; §2-`rule`-Feld zeigt jetzt auf die Modulliste statt einer Enum |
| 2026-06-24 | §[`DC-FA-VER-001.a`](spezifikation.md#dc-fa-ver-001a--versions-pin-konsistenz-versions) + §2-Schema (`versions.pin-pattern`/`versions.current-from`/`versions.exempt-paths`) + Grund-Code `version-stale` ergänzt: opt-in Modul `versions` prüft Versions-Pins gegen die aus `versions.current-from` (Default `version.md#aktuell`) gelesene aktuelle Version; liest Pins **auch in Fences** (gescopte Ausnahme, Muster-Scan ohne Parser), Ventile `exempt-paths`/`d-check:ignore`, fail-closed bei unauflösbarer Quelle, diagnose-only (Auto-Bump-`--repair` als Folge-CR). Default-aus byte-identisch |
| 2026-06-22 | §[`DC-FA-CODE-001.a`](spezifikation.md#dc-fa-code-001a--pfade-in-inline-code) + §2-Schema ergänzt: Datei-Ventil `codepaths.exempt-paths` (Glob wie `scan.ignore`) nimmt ganze Dateien von der `codepaths`-Prüfung aus — datei-weit, unabhängig von `codepaths.roots`; Vorbild das gleichnamige ids-Ventil. Abwärtskompatibel: ohne gesetztes `exempt-paths` byte-identisch |
| 2026-06-19 | §[`DC-FA-CLI-007.a`](spezifikation.md#dc-fa-cli-007a--diagnose-modus) Schritt 6 + [JSON-Diagnose](spezifikation.md#json-diagnose---doctor---json)-Schema (§2) ergänzt: `--doctor --json` rendert dieselbe Diagnose maschinenlesbar — `findings` zusätzlich mit `reasonText` und `fixCandidate` (`{original,replacement,note}` oder explizit `null`), `file`-Gruppierung; nur noch `--repair`+`--json` und `--doctor`+`--repair` sind Nutzungsfehler |
| 2026-06-18 | §[`DC-FA-CLI-008.a`](spezifikation.md#dc-fa-cli-008a--reparatur-patch) ergänzt: Reparatur-Modus `--repair` — unified diff auf stdout (`git apply`-kompatibel), zwei Stufen (`--repair`/`--repair-broad`); konservativ nur eindeutige `id-unlinked`-Fixes auf nackte Prosa-Vorkommen, breit Best-Guess `target-missing` (eindeutiger Basisname) mit review-pflichtig-Marker auf stderr; nicht mit `--json`/`--doctor` kombinierbar; Determinismus über sortierte Edits |
| 2026-06-18 | §[`DC-FA-CLI-007.a`](spezifikation.md#dc-fa-cli-007a--diagnose-modus) ergänzt: Diagnose-Modus `--doctor` — Lese-Lauf, nach Datei gruppierte Klartext-Diagnose auf stdout (statt Befund-Zeilen), Fix-Kandidat nur für `id-unlinked` (Link auf das ids-`target`); Grund-Klartext-Mapping über alle 14 Grund-Codes mit Vollständigkeits-Prüfung gegen die Reason-Konstanten; `--doctor`+`--json` = Nutzungsfehler (Exit 2); Determinismus über die sortierte Befundliste |
| 2026-06-17 | §[`DC-FA-MTX-001.a`](spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung) Schritt 4 ergänzt: Supersede-Lineage-Ausnahme (`allow-supersede-lineage`, `supersede-fields`) nimmt die deklarierte Lineage-Kante von der `matrix-inactive`-Prüfung aus (Match über Linktext bzw. Zielpfad der Referenz); Default aus ⇒ byte-identisch. Schema-Tabelle + Beispiel ergänzt |
| 2026-06-15 | §[`DC-FA-ANCH-001.b`](spezifikation.md#dc-fa-anch-001b--inline-html-anker) ergänzt: Inline-HTML-Anker (`id` an beliebigem Element, `name` an `<a>`) zählen wörtlich zur gültigen Anker-Menge; konservativ, zeilenbasiert, außerhalb Fences. §[`DC-FA-CODE-001.a`](spezifikation.md#dc-fa-code-001a--pfade-in-inline-code) fortgeschrieben: `codepaths`-Anker-Prüfung gegen die gemeinsame Anker-Menge statt nur Heading-Slugs |
| 2026-06-12 | Inline-Code-Erkennung absatzweise statt zeilenweise (§[`DC-FA-LINK-001.a`](spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion) Schritt 2): mehrzeilige Code-Spans gemäß CommonMark, Absatzgrenzen Leerzeile/Fence, ungeschlossene Folge literal. Anlass: [`DC-QA-04`](lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Gegentest u-boot — über Zeilenumbrüche gebrochene Befehls-Spans invertierten die Backtick-Parität der Folgezeile und erzeugten False-Positive-`id-unlinked`-Befunde auf korrekt verlinkten Kennungen. Zeilenbasierte **Link**-Extraktion (Schritt 3) bleibt normative Grenze |
| 2026-06-12 | Scan-Härtung aus der pkcs11-course-Adoption: `scan.ignore`-Muster prunen den Verzeichnis-Abstieg (vollständig ignorierte Teilbäume werden nicht betreten — unlesbare ignorierte Verzeichnisse wie root-eigene Build-Reste sind kein Laufzeitfehler mehr); `SKIP_DIRS` um `.gradle` ergänzt (Parität zur JS-Alt-Familie) |
| 2026-06-12 | Modul-lokaler Scan-Scope normiert (§[`DC-FA-CONF-002.a`](spezifikation.md#dc-fa-conf-002a--effektiver-scan-scope-pro-modul): `<modul>.scope` ersetzt den globalen Scope je Modul, `roots` Pflicht innerhalb `scope`, Lauf über die Vereinigungsmenge mit Einmal-Lese-Garantie, Zusammenfassung zählt die Union); Schema um `<modul>.scope.roots`/`.ignore` |
| 2026-06-12 | Modul `spans` normiert (§[`DC-FA-SPAN-001.a`](spezifikation.md#dc-fa-span-001a--span-artefakt-erkennung): `span-unclosed` absatzweise mit Folgezeichen-Bedingung und 30-Zeichen-Kappung, `span-nested-link` lexikalisch auf vorverarbeiteten Zeilen mit 40-Zeichen-Kappung; kein Opt-out, nur generischer `scope`); Grund-Codes ergänzt |
| 2026-06-12 | Modul `hostpaths` normiert (§[`DC-FA-HOST-001.a`](spezifikation.md#dc-fa-host-001a--host-pfad-erkennung): Prosa inkl. Inline-Code, Fences ausgenommen; Unix-Präfixliste konfigurierbar via `hostpaths.prefixes`, Windows-/UNC-Muster fest; Wortgrenzen-Vorbedingung, Satzzeichen-Normalisierung; bekannte Grenze Repo-Verzeichnis mit Präfix-Namen dokumentiert); Schema + Grund-Code ergänzt |
| 2026-06-12 | Kalibrierungs-Verengungen `hostpaths`: UNC-Servername beginnt alphanumerisch (escapte Regex-Beispiele matchen nicht mehr); Default-Präfixliste ohne tmp (Lastenheft 0.7.2) |
| 2026-06-11 | Modul `codepaths` normiert (§[`DC-FA-CODE-001.a`](spezifikation.md#dc-fa-code-001a--pfade-in-inline-code): rohe Prosa-Zeilen, Marker-Semantik, Normalisierung, konservative Erkennung, Anker-Prüfung); Schema um `codepaths.roots`, Grund-Code `codepath-missing`, `repo-escape`/`anchor-missing` auch für codepaths; Modul-Aufzählungen ergänzt |
| 2026-06-11 | Review R1: §[`DC-QA-01.a`](spezifikation.md#dc-qa-01a--benchmark)-Messprotokoll um die 2-vCPU-Begrenzung aus dem Lastenheft präzisiert (`--cpus 2`); N ungerade (Median = mittleres Element) |
| 2026-06-11 | Spez-Schuld eingelöst: §[`DC-QA-01.a`](spezifikation.md#dc-qa-01a--benchmark) Benchmark-Definition (Fixture, Messprotokoll, Pass-Kriterium) |
| 2026-06-11 | Review R1: Fragment-Teil vor Prüfung/Dedupe entfernt (Befund nennt Original-Linkziel); Schema-Vergleich case-insensitiv; Timeout gilt pro Request (Fallback: bis zu zwei Requests); explizit gesetzte 0 in `external.timeout-seconds`/`parallel` ist Konfigurationsfehler |
| 2026-06-11 | Modul `external` normiert umgesetzt; §[`DC-FA-EXT-001.a`](spezifikation.md#dc-fa-ext-001a--externe-erreichbarkeit) präzisiert: Transportfehler (DNS/Verbindung) → `external-status` (Status 0); Dedupe-Semantik expliziert (eine Prüfung pro URL, Befund an jedem Vorkommen) |
| 2026-06-11 | Review R1: Status-Extraktion liest nur Prosa-Zeilen (Fence-Inhalt ist kein Statuswert) und nur Markdown-Ziele (andere gelten als aktiv); Regel- und Status-Prüfung als unabhängig expliziert (ein Link kann zwei Befunde erzeugen) |
| 2026-06-11 | Modul `matrix` normiert umgesetzt; §[`DC-FA-ID-001.a`](spezifikation.md#dc-fa-id-001a--kennungs-prüfung) fortgeschrieben (Befund der Dogfooding-Selbstkonfiguration): ATX-Heading-Zeilen und Vorkommen im deklarierten Muster-Target (Definitions-Ort) sind linkpflichtfrei |
| 2026-06-11 | Review R1: Inline-Code-Stripping positionserhaltend (Leerzeichen statt Entfernen — keine Schein-Vorkommen); Repo-Escape-Verbot für `scan.roots` und `ids.patterns[].target`; Leerstring-matchende ids-Regexe als Konfigurationsfehler; zeilenbasierte Link-Extraktion als normative Grenze dokumentiert |
| 2026-06-11 | Modul `ids` normiert umgesetzt; §[`DC-FA-ID-001.a`](spezifikation.md#dc-fa-id-001a--kennungs-prüfung) präzisiert: Ziel-Klammern von Links und Bildreferenzen (Alt-Text, Ziel) sind kein Fließtext (linkpflichtfrei); Überlappungs-Semantik der Muster-Präzedenz expliziert; Target-Existenz-Constraint im Algorithmus-Text verankert |
| 2026-06-10 | Modul `anchors` normiert umgesetzt; `scan.roots`-Wert `"."` = gesamte Repo-Wurzel; Slug-Schritt 1 präzisiert: Emphasis-Sterne entfallen, literale Unterstriche bleiben (GitHub-Verhalten) |
| 2026-06-10 | Review R2 (Black-Box): hängendes wertnehmendes Flag = Nutzungsfehler; `d-check: error:`-Präfix für Flag-Fehler; `-h` → Usage auf stderr, Exit 0 |
| 2026-06-10 | Review-Runde Implementierung (Black-Box): Optionen vor/nach Pfad-Argument; gänzlich leere Wurzel → Exit 2 mit Mount-Hinweis vs. „ohne Markdown" → Exit 0; leere `.d-check.yml` = Defaults; explizit leere Listen; absolute Ziele; Verzeichnis-Symlinks beim Scan |
| 2026-06-10 | Referenzrichtungs-Korrektur: ADR-Abwärtsverweise entfernt — Spec-Straten verweisen nie abwärts; Traceability über die `Schärft:`-Felder der ADRs (Kurs-Baseline-Korrektur, [`MR-006`](../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)) |
| 2026-06-10 | Review R2: Status-Extraktions-Reihenfolge fixiert (`**Status:**` vor `Status`-Heading), `exclude-sections`-Matching definiert (getrimmt, case-sensitiv), Exit-2-Hinweis an Config-Constraint-Tabelle |
| 2026-06-10 | Review R1: `scan.roots`-Constraint präzisiert (nur deklarierte Wurzeln pflichtig), Symlink-Prüf-Scope präzisiert, unspezifizierter Grund-Code `nested-link` entfernt |
| 2026-06-10 | Initiale Fassung |
