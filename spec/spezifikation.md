# Spezifikation — d-check

**Status:** Aktiv. **Letzte Änderung:** 2026-07-03.

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
Verzeichnis-Symlinks werden beim Scan weder verfolgt noch
als Dateien gewertet (Symlink-Ablehnung,
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
Pfad-Reihenfolge, keine Map-Iteration für die Ausgabe). Das Modulset ist
fix (`links`, `anchors`, `ids`, `matrix`, `codepaths`; kein Probelauf);
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
modules: [links, anchors, ids, matrix, codepaths]
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

1. **Anforderungen** aus `spec/lastenheft.md`: führende Heading-Kennungen
   der Anforderungs-Gestalt `<PREFIX>-FA-<BEREICH>-NNN` bzw. `<PREFIX>-QA-NN`
   (präfix-agnostisch); Titel = Heading-Klartext ohne Kennung/Trenner.
2. **Referenzen** aus `docs/plan/adr/` (ADR-Kennung über den Dateinamen
   `NNNN-…md` → `ADR-NNNN`) und `docs/plan/planning/` (Slice-Kennung
   `slice-NNN-…md` → `slice-NNN`): je Datei alle Vorkommen einer
   Anforderungs-Kennung sammeln; Dateien ohne eigene Kennung (z. B.
   `README.md`) übersprungen.
3. **Zeile je Anforderung**: ID, Titel, sortierte ADR-/Slice-Kennungen,
   **Waise** = keine referenzierende Slice-Kennung. Default-Rendering
   Markdown-Tabelle; `--json`/`--yaml` serialisieren dieselbe Struktur
   (`requirements[]`, `total`, `orphans`).

Fehlende Quellen (kein Lastenheft / kein `adr`/`planning`-Verzeichnis)
liefern eine leere bzw. teilbefüllte Matrix, **kein** Fehler. Nur mit
`--doctor`/`--repair` ist `--trace` ein Nutzungsfehler (Exit 2).

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

### DC-FA-LINK-001.a — Markdown-Vorverarbeitung und Link-Extraktion

1. **Fences:** Zeilen, deren erste Nicht-Leerzeichen-Folge mit
   ` ``` ` oder `~~~` beginnt, schalten den Fence-Zustand um;
   Zeilen im Fence-Zustand werden von allen Modulen ignoriert.
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

**Referenz-Ventil (`ignore-refs`):** Ein in Schritt 5 Wurzel-relativ aufgelöster
Ziel-Pfad, der ein Glob aus `codepaths.ignore-refs` matcht (Syntax wie `scan.ignore`),
wird **nicht** existenz-, escape- oder anker-geprüft (kein `codepath-missing`,
`repo-escape` oder `anchor-missing`) — **referenz-weit**, unabhängig von Datei
und Zeile (anders als das datei-weite `exempt-paths` und das zeilenweise
`d-check:ignore`). Register bewusst entfernter/historischer Artefakte (Tombstones),
deren Pfad immutable/historische Doku noch zitiert; es unterdrückt nur die Prüfung
*dieses* Pfads, keine anderen Befunde. Ohne gesetztes `ignore-refs` byte-identisch.

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
   ein Zeilen-Suffix `:NNN` abtrennen (Datei:Zeile-Konvention),
   umschließende einfache/doppelte Anführungszeichen und schließende
   Satzzeichen (`.,;:`) entfernen.
4. Konservative Pfad-Erkennung — **kein** Pfad ist ein Wert, der leer
   ist, Whitespace oder Platzhalter-/Glob-Zeichen (`{}<>|*?=`)
   enthält, Ellipsen/Pfeile (`…`, `->`, `→`) enthält, mit `//` oder
   `#` beginnt oder ein externes Schema trägt. **Pfad** ist, was mit
   `./` oder `../` beginnt (Datei-relativ) oder mit einem der
   konfigurierten Präfixe aus `codepaths.roots` (Wurzel-relativ;
   Vergleich gegen `präfix/`).
5. Auflösung wie im Modul `links` (inkl. RFC-3986-Dekodierung):
   Fragment abtrennen; matcht der Wurzel-relative Pfad ein
   `codepaths.ignore-refs`-Glob, wird er übersprungen (kein Befund, s. o.
   Referenz-Ventil); sonst: Escape → `repo-escape`; fehlendes Ziel →
   `codepath-missing`. Trägt der Wert ein Fragment und ist das Ziel
   eine Markdown-Datei, wird der Anker gegen die gültige Anker-Menge
   der Zieldatei geprüft (Heading-Slugs und Inline-HTML-Anker; Verfahren
   und Cache wie
   [DC-FA-ANCH-001.a](#dc-fa-anch-001a--github-slug-algorithmus) /
   [DC-FA-ANCH-001.b](#dc-fa-anch-001b--inline-html-anker);
   Treffer fehlt → `anchor-missing`). Nicht lesbare Ziele: das Modul
   schweigt zum Anker (Existenz wurde bereits geprüft).

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

Beide Prüfungen kennen keinen Opt-out-Marker; das Modul akzeptiert
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
   Anker. Aus dem Span wird das **erste** Vorkommen eines Versions-Teilmusters
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
   höherrangigen Überschrift; HTML-Anker ab dessen Zeile). Lässt sich Datei/Anker
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
   Bedingung `vcs.immutable-when` erfüllt (Zeilen-Regex, erstes Vorkommen — z. B.
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
   Grün — sonst gälte still „aktiv" bzw. gälte nur der erste Block). Fehlende/unlesbare
   `planning.roadmap` ⇒ `planning-drift`.
4. **hasActive.** Der Aktiv-Status-Block reicht von der `planning.heading`-Zeile bis
   zur nächsten `## `-H2 (exklusive). Trägt er den `planning.marker` (Default „Keine
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
   **ausschließlich aus Tabellenzeilen** extrahiert — eine Tabellenzeile ist eine
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
   Ohne aktives `targets` ist der Befundsatz byte-identisch.

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
```

Jede Verletzung eines Constraints der folgenden Tabelle führt zu
Exit 2 ohne Prüfung
([`DC-FA-CONF-001`](lastenheft.md#dc-fa-conf-001--konfigurationsdatei)).

| Schlüssel | Typ | Default | Constraint |
|---|---|---|---|
| `scan.roots` | string[] | `DEFAULT_SCAN_ROOTS` | alle hier deklarierten Wurzeln müssen existieren und innerhalb der Repo-Wurzel liegen (Exit 2); nur die Default-Wurzeln (kein `scan.roots` gesetzt) sind optional; `"."` steht für die gesamte Repo-Wurzel (rekursiv; die `SKIP_DIRS` aus [§3](#3-defaults-und-konstanten) gelten immer und sind nicht konfigurierbar) |
| `scan.ignore` | string[] | leer | Glob-Syntax; Muster prunen auch den Verzeichnis-Abstieg — ein vollständig ignorierter Teilbaum (`pfad/**` oder direkt matchendes Muster) wird nicht betreten, unlesbare ignorierte Verzeichnisse sind dadurch kein Laufzeitfehler |
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
| `codepaths.roots` | string[] | leer | Präfixe relativ zur Repo-Wurzel: nicht leer, nicht absolut, kein `..` (Exit 2); `./`/`../` werden immer erkannt |
| `codepaths.exempt-paths` | string[] | leer | Glob (wie `scan.ignore`, relativ zur Repo-Wurzel); Dateien ganz ohne `codepaths`-Prüfung — datei-weit, unabhängig von `roots` |
| `codepaths.ignore-refs` | string[] | leer | Glob (wie `scan.ignore`); **aufgelöste Ziel-Pfade**, die matchen, werden weder existenz-, escape- noch anker-geprüft (kein `codepath-missing`/`repo-escape`/`anchor-missing`) — **referenz-weit** (Datei-/Zeilen-unabhängig), Tombstone-Register bewusst entfernter Artefakte; ohne Eintrag byte-identisch |
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
| `tracked.exempt-targets` | string[] | leer | Glob (wie `scan.ignore`); **aufgelöste Ziel-Pfade**, die matchen, werden nicht auf Getrackt-Status geprüft — **referenz-weit** (analog `codepaths.ignore-refs`), für absichtlich untrackte Ziele; jedes Glob **segmentweise** gültig und nicht leer (sonst Exit 2); ohne Eintrag byte-identisch ([`DC-FA-TRK-001`](lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in)) |
| `targets.makefiles` | string[] | leer | Wurzel-relative Makefile-Dateien, aus denen Regelnamen per statischer Zeilen-Heuristik extrahiert werden; leer ⇒ Modul inert; eine fehlende/unlesbare Datei ⇒ Exit 2 ([`DC-FA-TGT-001`](lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)) |
| `targets.doc-tables` | string[] | leer | Wurzel-relative Doku-Dateien; ihre `make X`-**Tabellenzeilen** (nur Zeilen mit Pipe in Spalte 0, keine Prosa) werden gegen die Makefile-Regelmenge geprüft (Richtung 1 `gate-phantom`); leer ⇒ Richtung 1 entfällt; fehlende Datei ⇒ Exit 2 |
| `targets.authority` | string | leer | Wurzel-relative Doku-Datei; **jede** nicht-exempte Makefile-Regel muss dort als `make X`-Tabellenzeile stehen (Richtung 2 `gate-undocumented`); leer ⇒ Richtung 2 entfällt; fehlende Datei ⇒ Exit 2 |
| `targets.exempt-targets` | string[] | leer | Regelnamen (exakt), die von der Doku-Pflicht (Richtung 2) ausgenommen sind (Utility-Targets); ohne Eintrag prüft Richtung 2 jede Regel |

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
| `anchor-missing` | anchors, codepaths | Anker entspricht keinem Heading-Slug |
| `id-unlinked` | ids | Kennung im Fließtext ohne Markdown-Link |
| `matrix-forbidden` | matrix | Referenz zwischen Klassen nicht erlaubt (**Link** oder, bei gesetztem `matrix.classes[].token`, **bare ID-Token** im Körper; Token-Form via `<!-- d-check:status-provenance -->` deklarierbar) |
| `matrix-inactive` | matrix | Referenz auf Dokument mit verbotenem Status |
| `matrix-downward` | matrix | klasseninterner Abwärtsverweis gegen die deklarierte Rangordnung (`order`/`direction: no-downward`) |
| `external-status` | external | HTTP-Status ≥ 400 oder Transportfehler (DNS/Verbindung) |
| `external-timeout` | external | Timeout überschritten |
| `codepath-missing` | codepaths | Ziel eines Inline-Code-Pfads existiert nicht |
| `hostpath-forbidden` | hostpaths | host-lokaler absoluter Pfad (Maschinen-Layout-Leak) in Prosa oder Inline-Code |
| `diagram-id-undefined` | diagrams | Kennung in einem geöffneten Diagramm-Fence ohne Definition in ihrer `defined-in`-Quelle |
| `span-unclosed` | spans | ungeschlossene Code-Span-Öffnung klebt an Nicht-Whitespace (Absatz-Parität gekippt) |
| `span-nested-link` | spans | Link-Syntax im Linktext eines weiteren Links (rendert zerrissen) |
| `external-redirects` | external | mehr als `REDIRECT_MAX` Redirects |
| `version-stale` | versions | Versions-Pin weicht von der aktuellen Version (`versions.current-from`) ab |
| `link-stale` | pins | normalisierter Ziel-Span eines gepinnten Links weicht vom hinterlegten `dpin`-Hash ab |
| `core-drift` | immutable | normalisierter Core einer gepinnten Datei (ohne Marker-Zeile + `exclude-sections`) weicht vom hinterlegten `immutable`-Hash ab |
| `core-drift-vcs` | vcs | Core einer immutablen Datei (BASE erfüllt `vcs.immutable-when`) hat sich über die Commit-Range geändert, ihr Status-Übergang ist unzulässig (`vcs.head-allow`), oder die immutable Datei wurde gelöscht/umbenannt |
| `commit-untraceable` | commits | bereinigte Commit-Message trägt keine Kennung nach `commits.id-patterns` und ist nicht per `commits.exempt-pattern` (Betreff) ausgenommen |
| `planning-drift` | planning | Roadmap-Aktiv-Status (`planning.marker` im `planning.heading`-Block) und Präsenz von `planning.slice-glob`-Slices sind inkonsistent (`hasActive ≠ hasSlices`), oder die kanonische Überschrift fehlt/ist mehrdeutig bzw. die Roadmap-Datei fehlt (fail-closed) |
| `target-untracked` | tracked | aufgelöstes, **existierendes** Link-/Bild-Ziel ist nicht im git-Index getrackt (untracked/gitignoriert) — die Referenz wäre auf jedem frischen Klon `target-missing` |

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

| Datum | Änderung | Verweis |
|---|---|---|
| 2026-06-10 | Initiale Fassung | slice-002 |
| 2026-06-10 | Review R1: `scan.roots`-Constraint präzisiert (nur deklarierte Wurzeln pflichtig), Symlink-Prüf-Scope präzisiert, unspezifizierter Grund-Code `nested-link` entfernt | slice-002 |
| 2026-06-10 | Review R2: Status-Extraktions-Reihenfolge fixiert (`**Status:**` vor `Status`-Heading), `exclude-sections`-Matching definiert (getrimmt, case-sensitiv), Exit-2-Hinweis an Config-Constraint-Tabelle | slice-002 |
| 2026-06-10 | Referenzrichtungs-Korrektur: ADR-Abwärtsverweise entfernt — Spec-Straten verweisen nie abwärts; Traceability über die `Schärft:`-Felder der ADRs (Kurs-Baseline-Korrektur, [`MR-006`](../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)) | slice-002 |
| 2026-06-10 | Review-Runde Implementierung (Black-Box): Optionen vor/nach Pfad-Argument; gänzlich leere Wurzel → Exit 2 mit Mount-Hinweis vs. „ohne Markdown" → Exit 0; leere `.d-check.yml` = Defaults; explizit leere Listen; absolute Ziele; Verzeichnis-Symlinks beim Scan | slice-003 |
| 2026-06-10 | Review R2 (Black-Box): hängendes wertnehmendes Flag = Nutzungsfehler; `d-check: error:`-Präfix für Flag-Fehler; `-h` → Usage auf stderr, Exit 0 | slice-003 |
| 2026-06-10 | Modul `anchors` normiert umgesetzt; `scan.roots`-Wert `"."` = gesamte Repo-Wurzel; Slug-Schritt 1 präzisiert: Emphasis-Sterne entfallen, literale Unterstriche bleiben (GitHub-Verhalten) | slice-004 |
| 2026-06-11 | Modul `ids` normiert umgesetzt; §[`DC-FA-ID-001.a`](spezifikation.md#dc-fa-id-001a--kennungs-prüfung) präzisiert: Ziel-Klammern von Links und Bildreferenzen (Alt-Text, Ziel) sind kein Fließtext (linkpflichtfrei); Überlappungs-Semantik der Muster-Präzedenz expliziert; Target-Existenz-Constraint im Algorithmus-Text verankert | slice-006 |
| 2026-06-11 | Review R1 zu slice-006: Inline-Code-Stripping positionserhaltend (Leerzeichen statt Entfernen — keine Schein-Vorkommen); Repo-Escape-Verbot für `scan.roots` und `ids.patterns[].target`; Leerstring-matchende ids-Regexe als Konfigurationsfehler; zeilenbasierte Link-Extraktion als normative Grenze dokumentiert | slice-006 |
| 2026-06-11 | Modul `matrix` normiert umgesetzt; §[`DC-FA-ID-001.a`](spezifikation.md#dc-fa-id-001a--kennungs-prüfung) fortgeschrieben (Befund der Dogfooding-Selbstkonfiguration): ATX-Heading-Zeilen und Vorkommen im deklarierten Muster-Target (Definitions-Ort) sind linkpflichtfrei | slice-007 |
| 2026-06-11 | Review R1 zu slice-007: Status-Extraktion liest nur Prosa-Zeilen (Fence-Inhalt ist kein Statuswert) und nur Markdown-Ziele (andere gelten als aktiv); Regel- und Status-Prüfung als unabhängig expliziert (ein Link kann zwei Befunde erzeugen) | slice-007 |
| 2026-06-11 | Modul `external` normiert umgesetzt; §[`DC-FA-EXT-001.a`](spezifikation.md#dc-fa-ext-001a--externe-erreichbarkeit) präzisiert: Transportfehler (DNS/Verbindung) → `external-status` (Status 0); Dedupe-Semantik expliziert (eine Prüfung pro URL, Befund an jedem Vorkommen) | slice-008 |
| 2026-06-11 | Review R1 zu slice-008: Fragment-Teil vor Prüfung/Dedupe entfernt (Befund nennt Original-Linkziel); Schema-Vergleich case-insensitiv; Timeout gilt pro Request (Fallback: bis zu zwei Requests); explizit gesetzte 0 in `external.timeout-seconds`/`parallel` ist Konfigurationsfehler | slice-008 |
| 2026-06-11 | Spez-Schuld eingelöst: §[`DC-QA-01.a`](spezifikation.md#dc-qa-01a--benchmark) Benchmark-Definition (Fixture, Messprotokoll, Pass-Kriterium) | slice-009 |
| 2026-06-11 | Review R1 zu slice-009/010: §[`DC-QA-01.a`](spezifikation.md#dc-qa-01a--benchmark)-Messprotokoll um die 2-vCPU-Begrenzung aus dem Lastenheft präzisiert (`--cpus 2`); N ungerade (Median = mittleres Element) | slice-009 |
| 2026-06-11 | Modul `codepaths` normiert (§[`DC-FA-CODE-001.a`](spezifikation.md#dc-fa-code-001a--pfade-in-inline-code): rohe Prosa-Zeilen, Marker-Semantik, Normalisierung, konservative Erkennung, Anker-Prüfung); Schema um `codepaths.roots`, Grund-Code `codepath-missing`, `repo-escape`/`anchor-missing` auch für codepaths; Modul-Aufzählungen ergänzt | slice-013 |
| 2026-06-12 | Kalibrierungs-Verengungen `hostpaths` (slice-016): UNC-Servername beginnt alphanumerisch (escapte Regex-Beispiele matchen nicht mehr); Default-Präfixliste ohne tmp (Lastenheft 0.7.2) | slice-016 |
| 2026-06-12 | Modul `hostpaths` normiert (§[`DC-FA-HOST-001.a`](spezifikation.md#dc-fa-host-001a--host-pfad-erkennung): Prosa inkl. Inline-Code, Fences ausgenommen; Unix-Präfixliste konfigurierbar via `hostpaths.prefixes`, Windows-/UNC-Muster fest; Wortgrenzen-Vorbedingung, Satzzeichen-Normalisierung; bekannte Grenze Repo-Verzeichnis mit Präfix-Namen dokumentiert); Schema + Grund-Code ergänzt | slice-016 |
| 2026-06-12 | Modul `spans` normiert (§[`DC-FA-SPAN-001.a`](spezifikation.md#dc-fa-span-001a--span-artefakt-erkennung): `span-unclosed` absatzweise mit Folgezeichen-Bedingung und 30-Zeichen-Kappung, `span-nested-link` lexikalisch auf vorverarbeiteten Zeilen mit 40-Zeichen-Kappung; kein Opt-out, nur generischer `scope`); Grund-Codes ergänzt | slice-015 |
| 2026-06-12 | Modul-lokaler Scan-Scope normiert (§[`DC-FA-CONF-002.a`](spezifikation.md#dc-fa-conf-002a--effektiver-scan-scope-pro-modul): `<modul>.scope` ersetzt den globalen Scope je Modul, `roots` Pflicht innerhalb `scope`, Lauf über die Vereinigungsmenge mit Einmal-Lese-Garantie, Zusammenfassung zählt die Union); Schema um `<modul>.scope.roots`/`.ignore` | slice-017 |
| 2026-06-15 | §[`DC-FA-ANCH-001.b`](spezifikation.md#dc-fa-anch-001b--inline-html-anker) ergänzt: Inline-HTML-Anker (`id` an beliebigem Element, `name` an `<a>`) zählen wörtlich zur gültigen Anker-Menge; konservativ, zeilenbasiert, außerhalb Fences. §[`DC-FA-CODE-001.a`](spezifikation.md#dc-fa-code-001a--pfade-in-inline-code) fortgeschrieben: `codepaths`-Anker-Prüfung gegen die gemeinsame Anker-Menge statt nur Heading-Slugs | slice-022 |
| 2026-06-17 | §[`DC-FA-MTX-001.a`](spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung) Schritt 4 ergänzt: Supersede-Lineage-Ausnahme (`allow-supersede-lineage`, `supersede-fields`) nimmt die deklarierte Lineage-Kante von der `matrix-inactive`-Prüfung aus (Match über Linktext bzw. Zielpfad der Referenz); Default aus ⇒ byte-identisch. Schema-Tabelle + Beispiel ergänzt | slice-024 |
| 2026-06-12 | Scan-Härtung aus der pkcs11-course-Adoption (slice-014): `scan.ignore`-Muster prunen den Verzeichnis-Abstieg (vollständig ignorierte Teilbäume werden nicht betreten — unlesbare ignorierte Verzeichnisse wie root-eigene Build-Reste sind kein Laufzeitfehler mehr); `SKIP_DIRS` um `.gradle` ergänzt (Parität zur JS-Alt-Familie) | slice-014 |
| 2026-06-12 | Inline-Code-Erkennung absatzweise statt zeilenweise (§[`DC-FA-LINK-001.a`](spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion) Schritt 2): mehrzeilige Code-Spans gemäß CommonMark, Absatzgrenzen Leerzeile/Fence, ungeschlossene Folge literal. Anlass: [`DC-QA-04`](lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Gegentest u-boot — über Zeilenumbrüche gebrochene Befehls-Spans invertierten die Backtick-Parität der Folgezeile und erzeugten False-Positive-`id-unlinked`-Befunde auf korrekt verlinkten Kennungen. Zeilenbasierte **Link**-Extraktion (Schritt 3) bleibt normative Grenze | slice-012 |
| 2026-06-18 | §[`DC-FA-CLI-007.a`](spezifikation.md#dc-fa-cli-007a--diagnose-modus) ergänzt: Diagnose-Modus `--doctor` — Lese-Lauf, nach Datei gruppierte Klartext-Diagnose auf stdout (statt Befund-Zeilen), Fix-Kandidat nur für `id-unlinked` (Link auf das ids-`target`); Grund-Klartext-Mapping über alle 14 Grund-Codes mit Vollständigkeits-Prüfung gegen die Reason-Konstanten; `--doctor`+`--json` = Nutzungsfehler (Exit 2); Determinismus über die sortierte Befundliste | slice-025 |
| 2026-06-18 | §[`DC-FA-CLI-008.a`](spezifikation.md#dc-fa-cli-008a--reparatur-patch) ergänzt: Reparatur-Modus `--repair` — unified diff auf stdout (`git apply`-kompatibel), zwei Stufen (`--repair`/`--repair-broad`); konservativ nur eindeutige `id-unlinked`-Fixes auf nackte Prosa-Vorkommen, breit Best-Guess `target-missing` (eindeutiger Basisname) mit review-pflichtig-Marker auf stderr; nicht mit `--json`/`--doctor` kombinierbar; Determinismus über sortierte Edits | slice-026 |
| 2026-06-19 | §[`DC-FA-CLI-007.a`](spezifikation.md#dc-fa-cli-007a--diagnose-modus) Schritt 6 + [JSON-Diagnose](spezifikation.md#json-diagnose---doctor---json)-Schema (§2) ergänzt: `--doctor --json` rendert dieselbe Diagnose maschinenlesbar — `findings` zusätzlich mit `reasonText` und `fixCandidate` (`{original,replacement,note}` oder explizit `null`), `file`-Gruppierung; nur noch `--repair`+`--json` und `--doctor`+`--repair` sind Nutzungsfehler | slice-029 |
| 2026-06-22 | §[`DC-FA-CODE-001.a`](spezifikation.md#dc-fa-code-001a--pfade-in-inline-code) + §2-Schema ergänzt: Datei-Ventil `codepaths.exempt-paths` (Glob wie `scan.ignore`) nimmt ganze Dateien von der `codepaths`-Prüfung aus — datei-weit, unabhängig von `codepaths.roots`; Vorbild das gleichnamige ids-Ventil. Abwärtskompatibel: ohne gesetztes `exempt-paths` byte-identisch | slice-043 |
| 2026-06-24 | §[`DC-FA-VER-001.a`](spezifikation.md#dc-fa-ver-001a--versions-pin-konsistenz-versions) + §2-Schema (`versions.pin-pattern`/`versions.current-from`/`versions.exempt-paths`) + Grund-Code `version-stale` ergänzt: opt-in Modul `versions` prüft Versions-Pins gegen die aus `versions.current-from` (Default `version.md#aktuell`) gelesene aktuelle Version; liest Pins **auch in Fences** (gescopte Ausnahme, Muster-Scan ohne Parser), Ventile `exempt-paths`/`d-check:ignore`, fail-closed bei unauflösbarer Quelle, diagnose-only (Auto-Bump-`--repair` als Folge-CR). Default-aus byte-identisch | slice-048 |
| 2026-06-24 | §[`DC-FA-PIN-001.a`](spezifikation.md#dc-fa-pin-001a--content-pin-gegen-inhaltlichen-drift-pins) + Grund-Code `link-stale` (§4) ergänzt: opt-in Modul `pins` hasst den whitespace-normalisierten **rohen** Ziel-Span (Datei/Heading-Section inkl. Fenced-Code) eines gepinnten Links (`<!-- dpin: sha256:… -->`, gebunden an den unmittelbar vorausgehenden Link derselben Zeile, sonst inert) und meldet `link-stale` bei Drift; nur auflösbare repo-interne Ziele (struktureller Befund bleibt `links`/`anchors`, kein Doppelbefund), Scope-treu (nur Quell-Dateien), diagnose-only; §2-`rule`-Feld zeigt jetzt auf die Modulliste statt einer Enum | slice-049 |
| 2026-06-28 | §[`DC-FA-MTX-001.a`](spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung) Schritt 5 + §2-Schema (`matrix.classes[].order`/`.direction`) + Grund-Code `matrix-downward` (§4) + Config-Beispiel ergänzt: klasseninterne Verweisrichtung ([`DC-FA-MTX-002`](lastenheft.md#dc-fa-mtx-002--verweisrichtung-innerhalb-einer-geordneten-dokumentklasse-modul-matrix)) — eine Klasse mit `order` (Glob-Rang, First-Match) + `direction: no-downward` meldet klasseninterne Abwärtsverweise (Rang *i* → *j > i*, auch transitiv) als `matrix-downward`; rangfreie Mitglieder und klassenübergreifende Referenzen ausgenommen; fail-closed-Config (`order`/`direction` nur zusammen, unbekannter `direction`-Wert ⇒ Exit 2); Default-aus byte-identisch | slice-050 |
| 2026-06-28 | §2 „Glob-Auswertung" ergänzt: alle Glob-Felder (`scan.ignore`, `<modul>.scope.ignore`, `matrix.classes[].paths`/`.order`, `*.exempt-paths`) werden segmentweise über Go-`path.Match` ausgewertet (`**` segmentübergreifend); negierte Zeichenklasse `[^…]` (Go), **nicht** `[!…]` (fnmatch). Reine Klarstellung des Bestands (`matchGlob`), kein Verhaltens-/Schema-Change | — |
| 2026-06-28 | §[`DC-FA-MTX-001.a`](spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung) Schritt 6 + §2-Schema (`matrix.classes[].token`, `matrix.exempt-paths`) + §4 (`matrix-forbidden` Token-Form) ergänzt: token-basierte Referenz-Richtung ([`DC-FA-MTX-003`](lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)) — `matrix` fängt verbotene Referenzen auch als bare ID-Token im Prosa-Körper (außer Links/Fences/`exclude-sections`); Provenance-Marker `<!-- d-check:status-provenance -->` auf der rohen Zeile nimmt aus; `exempt-paths` grandfathered ganze Dateien. Fail-closed (`token` kompiliert/Leerstring). Default-aus byte-identisch. Außerdem §[`DC-FA-SPAN-001.a`](lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in): Slice-Token aus dem Spec-Körper entfernt (Provenance gehört in die Historie) | slice-051 |
| 2026-06-29 | §[`DC-FA-VCS-001.a`](spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) + §2-Schema (`vcs.paths`/`immutable-when`/`exclude-sections`/`status-line`/`head-allow`) + Grund-Code `core-drift-vcs` (§4) ergänzt: opt-in Modul `vcs` vergleicht `core(BASE)` vs. `core(HEAD)` über eine Commit-Range (`--range <base>..<head>` / `--staged`), liest das read-only `.git` über einen reine-Go-VCS-Port (ohne git-Binary, ohne Netz); erweiterter Eingabe-Scope (git + Range), aber lokal/lesend/deterministisch — Determinismus/Read-only gehalten; fail-closed ohne `.git`/Range, diagnose-only. Core-Semantik in Parität zum abgelösten `adr-immutable-check.sh` (nur Kopf-Status-Zeile gestrippt, `exclude-sections`-Abschnitte). Default-aus byte-identisch. Außerdem §[`DC-FA-CLI-010.a`](spezifikation.md#dc-fa-cli-010a--makefile-fragment) (6→7 Targets): `--print-mk` trägt `doc-immutable` (`--enable vcs` + aus `ValidModules` abgeleitete Fokus-`--disable`-Liste, `RANGE`/`STAGED`) — verteilt die git-Garantie an Konsumenten | slice-053 |
| 2026-06-29 | §[`DC-FA-CODE-001.a`](spezifikation.md#dc-fa-code-001a--pfade-in-inline-code) + §2-Schema (`codepaths.ignore-refs`) ergänzt: Referenz-Ventil `ignore-refs` — ein aufgelöster Ziel-Pfad, der ein Glob matcht, wird nicht existenz-/anker-geprüft (**referenz-weit**, Tombstone-Register entfernter Artefakte); dritte Ventil-Achse neben dem zeilenweisen `d-check:ignore` und dem datei-weiten `exempt-paths`, in Schritt 5 vor `codepath-missing`, ohne Eintrag byte-identisch. Anlass: Frozen-Doc-Refactoring-Falle (immutable ADRs zitieren entfernte Pfade) — zugleich das in slice-053 behaltene `adr-immutable-check.sh` entfernt | slice-054 |
| 2026-07-01 | §[`DC-FA-COMMITS-001.a`](spezifikation.md#dc-fa-commits-001a--traceability-kennung-in-commit-messages-über-eine-commit-range-commits) + §2-Schema (`commits.id-patterns`/`commits.exempt-pattern`) + Grund-Code `commit-untraceable` (§4) ergänzt: opt-in Modul `commits` prüft, dass jede Commit-Message eine Kennung nach `commits.id-patterns` trägt (`commit-untraceable`), über **denselben VCS-Port** wie `vcs` (reine-Go-git, ohne git-Binary, ohne Netz), erweitert um Message-Lesen; zwei Quellen: `--range <base>..<head>` (Nicht-Merge-Commits, `--no-merges`-Parität) und `--commit-msg <datei\|->` (Kurzschluss-Modus für den commit-msg-Hook, einzelne Pending-Message). Uniforme `#`-/scissors-Bereinigung (Kennung auf Inhalts-Zeile), Betreff-Ausnahme `commits.exempt-pattern`; fail-closed ohne `.git`/Range/Message-Datei, diagnose-only. Default-aus byte-identisch. Portierung des abgelösten `tools/trace-check.sh` (dieselbe VCS-Port-Präzedenz wie `vcs`, auf Commit-Messages statt Datei-Inhalt). Außerdem §[`DC-FA-CLI-010.a`](spezifikation.md#dc-fa-cli-010a--makefile-fragment) (7→8 Targets): `--print-mk` trägt `doc-commits` (`--enable commits` + aus `ValidModules` abgeleitete Fokus-`--disable`-Liste, `--range`) — verteilt die Commit-Traceability an Konsumenten | slice-056 |
| 2026-07-01 | §[`DC-FA-PLAN-001.a`](spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) + §2-Schema (`planning.roadmap`/`heading`/`marker`/`slice-glob`) + Grund-Code `planning-drift` (§4) ergänzt: opt-in Modul `planning` prüft **hermetisch** (nur Arbeitsbaum, kein git/Netz) die Roadmap-↔-in-progress-Invariante — Ruhe-Marker im `## Aktuelle Welle`-Block genau dann, wenn kein `slice-*` im Verzeichnis (`hasActive == hasSlices`), sonst `planning-drift`; Heading-Guard fail-closed, Post-Pass wie `codepaths`-Existenz, diagnose-only. Default-aus byte-identisch. Portierung des abgelösten `tools/planning-consistency.sh`. Außerdem §[`DC-FA-CLI-010.a`](spezifikation.md#dc-fa-cli-010a--makefile-fragment) (8→9 Targets): `--print-mk` trägt `doc-planning` (`--enable planning`, hermetisch ohne `--range`) | slice-057 |
| 2026-07-03 | §[`DC-FA-TRK-001.a`](spezifikation.md#dc-fa-trk-001a--getrackt-status-auflösbarer-referenz-ziele-tracked) + §2-Schema (`tracked.exempt-targets`) + Grund-Code `target-untracked` (§4) ergänzt: opt-in Modul `tracked` prüft die Datei-Ebene der von `links` aufgelösten, **existierenden** repo-internen Ziele gegen den **git-Index** — dritte VCS-Port-Nutzung (`vcs` Range-Diff, `commits` Messages, `tracked` Index), **ohne** Range/`--staged`; Index statt `.gitignore`-Interpretation (gestagte Dateien gelten als getrackt), kein Doppelbefund (`target-missing` bleibt `links`), Ventil referenz-weit analog `codepaths.ignore-refs`; fail-closed ohne lesbares `.git` (Exit 2), diagnose-only, default-aus byte-identisch. Außerdem §[`DC-FA-CLI-010.a`](spezifikation.md#dc-fa-cli-010a--makefile-fragment) (9→10 Targets): `--print-mk` trägt `doc-tracked` (`--enable tracked` + fokussierte `--disable`-Liste, ohne Range) | slice-059 |
| 2026-07-03 | Review R1/R2 zu §[`DC-FA-TRK-001.a`](spezifikation.md#dc-fa-trk-001a--getrackt-status-auflösbarer-referenz-ziele-tracked), präzisiert: Prüfung **je gescannter Quell-Datei** statt „Post-Pass" (Auflösungs-Mechanik unabhängig von der Aktivierung des Moduls `links` — ein fokussierter Lauf prüft vollständig); Verzeichnis-Ziele und Symlink-Referenzen (Ziel ist/durchläuft einen Symlink) explizit kein Kandidat ([`DC-FA-LINK-002`](lastenheft.md#dc-fa-link-002--symlink-ablehnung)-Domäne — false-positive hinter getrackten Verzeichnis-Symlinks vermieden); `target` = aufgelöster Pfad bekräftigt (Ventil-Parität); `exempt-targets` segmentweise validiert (Exit 2) | slice-059 |
| 2026-07-05 | §[`DC-FA-TGT-001.a`](spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets) + §2-Schema (`targets.makefiles`/`doc-tables`/`authority`/`exempt-targets`) ergänzt: opt-in Modul `targets` prüft **hermetisch** (Filesystem-Port `ReadFile`, **kein** git/Netz/Makefile-Ausführen) die Doku-↔-Makefile-Deklarations-Konsistenz — ein nur aus **Tabellenzeilen** (Pipe-Präfix, keine Prosa) dokumentiertes `make X` ohne Makefile-Regel ⇒ `gate-phantom`, eine Makefile-Regel (nicht `targets.exempt-targets`) ohne Eintrag in `targets.authority` ⇒ `gate-undocumented`; statische Zeilen-Heuristik für Regelnamen (keine Pattern-Rules/variablen Targets, in Parität zu `tools/gate-consistency.sh`), fail-closed bei fehlender konfigurierter Datei, diagnose-only, default-aus byte-identisch; Analogie zu `planning` (Doku-Behauptung ↔ Repo-Struktur, hermetisch). Außerdem §[`DC-FA-CLI-010.a`](spezifikation.md#dc-fa-cli-010a--makefile-fragment) (10→11 Targets): `--print-mk` trägt `doc-targets` (`--enable targets`, hermetisch ohne Range). Die Grund-Codes `gate-phantom`/`gate-undocumented` (§4) folgen mit der Modul-Implementierung (AllReasons-↔-§4-Lockstep) | slice-063 |
