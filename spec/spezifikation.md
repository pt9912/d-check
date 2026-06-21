# Spezifikation — d-check

**Status:** Aktiv. **Letzte Änderung:** 2026-06-19.

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
    - {name: spec-straten, paths: [spec/lastenheft.md, spec/spezifikation.md, spec/architecture.md]}
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

1. Kommentar-Kopf: Einbindung via `include`, Hinweis zum
   `DCHECK_IMAGE`-Override auf einen `@sha256:`-Digest.
2. `DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v<version>` — `<version>` ist die
   beim Tag-Build via `-ldflags -X …/cli.version=<tag>` eingebettete
   Release-Version (Default `0.0.0-dev` für lokale/Gate-Builds).
3. `.PHONY: doc-check` und ein `doc-check`-Target mit **TAB**-eingerücktem
   `docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_IMAGE)`.

Der eigene Image-**Digest** wird NICHT eingebettet (er hasht das Binary
selbst — Henne-Ei); der Konsument pinnt per `DCHECK_IMAGE`-Override.

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
Vorverarbeitung entfernt Inline-Code gerade. **Schritte:**

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
   Fragment abtrennen; Escape → `repo-escape`; fehlendes Ziel →
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
   ebenfalls kein Treffer (Kalibrierungs-Befund slice-015: vendorte
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

## 2. Datenstrukturen und Schemas

### Befund

| Feld | Typ | Bedeutung |
|---|---|---|
| `file` | string | Pfad relativ zur Repo-Wurzel, `/`-getrennt |
| `line` | integer ≥ 1 | Zeile des Vorkommens |
| `rule` | string | Regelmodul (`links`, `anchors`, `ids`, `matrix`, `external`, `codepaths`) |
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
          "rule":    {"type": "string", "enum": ["links", "anchors", "ids", "matrix", "external", "codepaths"]},
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
| `matrix.classes[].paths` | string[] | — | Glob |
| `matrix.rules[]` | {from,to,allow} | — | Klassen müssen deklariert sein |
| `matrix.status.forbidden` | string[] | `[superseded, deprecated]` | case-insensitiv |
| `matrix.status.allow-supersede-lineage` | bool | `false` | nimmt die deklarierte Supersede-Lineage-Kante von der Status-Prüfung aus (nur `matrix-inactive`); ohne `true` byte-identisch |
| `matrix.status.supersede-fields` | string[] | leer | Feldnamen (z. B. `Supersedes`, `Aenderungstyp`), aus denen die Ablösung gelesen wird; Einträge nicht leer (Exit 2); nur wirksam bei `allow-supersede-lineage: true` |
| `matrix.exclude-sections` | string[] | leer | Vergleich gegen den getrimmten Heading-Text ohne Markdown-Auszeichnung, case-sensitiv |
| `external.timeout-seconds` | integer | 10 | 1–300 |
| `external.parallel` | integer | 4 | 1–16 |
| `codepaths.roots` | string[] | leer | Präfixe relativ zur Repo-Wurzel: nicht leer, nicht absolut, kein `..` (Exit 2); `./`/`../` werden immer erkannt |

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
| `matrix-forbidden` | matrix | Referenz zwischen Klassen nicht erlaubt |
| `matrix-inactive` | matrix | Referenz auf Dokument mit verbotenem Status |
| `external-status` | external | HTTP-Status ≥ 400 oder Transportfehler (DNS/Verbindung) |
| `external-timeout` | external | Timeout überschritten |
| `codepath-missing` | codepaths | Ziel eines Inline-Code-Pfads existiert nicht |
| `hostpath-forbidden` | hostpaths | host-lokaler absoluter Pfad (Maschinen-Layout-Leak) in Prosa oder Inline-Code |
| `span-unclosed` | spans | ungeschlossene Code-Span-Öffnung klebt an Nicht-Whitespace (Absatz-Parität gekippt) |
| `span-nested-link` | spans | Link-Syntax im Linktext eines weiteren Links (rendert zerrissen) |
| `external-redirects` | external | mehr als `REDIRECT_MAX` Redirects |

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
