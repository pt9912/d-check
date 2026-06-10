# Spezifikation — d-check

**Status:** Aktiv. **Letzte Änderung:** 2026-06-10.

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
Exit 2 (kein Teilergebnis als Erfolg).

### DC-FA-CLI-002.a — Modul-Auflösung

Effektive Module = (`modules` aus Config, sonst `DEFAULT_MODULES`)
∪ `--enable`-Angaben ∖ `--disable`-Angaben. CLI-Angaben werden nach
der Config angewandt (CLI-Präzedenz). Unbekannte Modulnamen sind
Nutzungsfehler (Exit 2) mit Auflistung der gültigen Namen.

### DC-FA-LINK-001.a — Markdown-Vorverarbeitung und Link-Extraktion

1. **Fences:** Zeilen, deren erste Nicht-Leerzeichen-Folge mit
   ` ``` ` oder `~~~` beginnt, schalten den Fence-Zustand um;
   Zeilen im Fence-Zustand werden von allen Modulen ignoriert.
2. **Inline-Code:** Backtick-Spans werden zeilenweise entfernt,
   längste Backtick-Folge zuerst (Mehrfach-Backticks).
3. **Extraktion:** Inline-Links `[text](ziel)` und Bilder
   `![alt](ziel)` per Klammer-balancierter Suche; mehrere Links pro
   Zeile werden alle erfasst. Ziele in `<…>` werden entquotet; ein
   Titel-Suffix (` "…"`) wird abgetrennt.
4. **Ziel-Normalisierung:** Prozent-Dekodierung (RFC 3986, vollständig)
   → Auflösung relativ zum Verzeichnis der enthaltenden Datei →
   lexikalische Normalisierung. Die Repo-Escape-Prüfung erfolgt **nach**
   der Dekodierung (`DC-FA-LINK-001`).
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
   Emphasis-Marker (`*`, `_`), Links → Linktext.
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

### DC-FA-ID-001.a — Kennungs-Prüfung

Pro Zeile (außerhalb Fence/Inline-Code) werden die konfigurierten
Muster in Deklarationsreihenfolge gematcht; das erste passende Muster
gewinnt pro Vorkommen. Ein Vorkommen gilt als verlinkt, wenn es
innerhalb des Linktexts eines Markdown-Links liegt; sonst Befund
`id-unlinked`.

### DC-FA-MTX-001.a — Klassen- und Status-Auflösung

1. **Klassenzuordnung:** Glob-Muster der `classes` in
   Deklarationsreihenfolge; die erste passende Klasse gilt. Dateien
   ohne Klasse nehmen nicht an der Matrix-Prüfung teil.
2. **Status-Extraktion** in fester Reihenfolge: (1) erste Zeile, die
   mit `**Status:**` beginnt; (2) sonst erste nicht-leere Textzeile
   unter einem `Status`-Heading (beliebige Ebene, Heading-Vergleich
   case-insensitiv). Sind beide Formen vorhanden, zählt die
   `**Status:**`-Zeile. Der Wert wird case-insensitiv als Präfix-Match
   gegen `status.forbidden` verglichen (so matcht
   `Superseded by ADR-0007` den Wert `superseded`). Ohne Status-Feld
   gilt das Dokument als aktiv.
3. **Sektions-Ausnahme:** Links innerhalb von Sektionen, deren
   Heading-Text (getrimmt, ohne Markdown-Auszeichnung, case-sensitiv)
   in `exclude-sections` steht (z. B. „Historie"),
   werden von `matrix` nicht geprüft (Provenance-Ausnahme gemäß
   [`DC-FA-MTX-001`](lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
   Out-of-Scope).

### DC-FA-EXT-001.a — Externe Erreichbarkeit

HEAD-Request; bei HTTP 405/501 Fallback auf GET (Body verworfen).
Redirects bis `REDIRECT_MAX` Stationen gefolgt. Pro URL genau eine
Prüfung pro Lauf (Dedupe), begrenzte Parallelität
(`EXTERNAL_PARALLEL`). Ergebnis-Auswertung: Status < 400 → kein
Befund; ≥ 400 → `external-status`; Timeout → `external-timeout`;
Redirect-Kette > `REDIRECT_MAX` → `external-redirects`.

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
| `rule` | string | Regelmodul (`links`, `anchors`, `ids`, `matrix`, `external`) |
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
          "rule":    {"type": "string", "enum": ["links", "anchors", "ids", "matrix", "external"]},
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

### `.d-check.yml`

Unbekannte Schlüssel sind Fehler (striktes Decoding).
Kommentiertes Vollbeispiel:

```yaml
scan:
  roots: [docs, spec]            # ersetzt die Default-Wurzeln; müssen existieren
  ignore: ["docs/archive/**"]    # Glob, relativ zur Repo-Wurzel
modules: [links, anchors, ids]   # ersetzt DEFAULT_MODULES
ids:
  patterns:                      # Reihenfolge = Präzedenz
    - regex: 'ADR-\d{4}'
      target: docs/plan/adr/     # Definition (Datei oder Verzeichnis)
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
  exclude-sections: [Historie, Geschichte]
external:
  timeout-seconds: 10
  parallel: 4
```

Jede Verletzung eines Constraints der folgenden Tabelle führt zu
Exit 2 ohne Prüfung
([`DC-FA-CONF-001`](lastenheft.md#dc-fa-conf-001--konfigurationsdatei)).

| Schlüssel | Typ | Default | Constraint |
|---|---|---|---|
| `scan.roots` | string[] | `DEFAULT_SCAN_ROOTS` | alle hier deklarierten Wurzeln müssen existieren (Exit 2); nur die Default-Wurzeln (kein `scan.roots` gesetzt) sind optional |
| `scan.ignore` | string[] | leer | Glob-Syntax |
| `modules` | string[] | `DEFAULT_MODULES` | nur gültige Modulnamen |
| `ids.patterns[].regex` | string | — | muss kompilieren (Exit 2) |
| `ids.patterns[].target` | string | — | muss existieren |
| `matrix.classes[].name` | string | — | eindeutig |
| `matrix.classes[].paths` | string[] | — | Glob |
| `matrix.rules[]` | {from,to,allow} | — | Klassen müssen deklariert sein |
| `matrix.status.forbidden` | string[] | `[superseded, deprecated]` | case-insensitiv |
| `matrix.exclude-sections` | string[] | leer | Vergleich gegen den getrimmten Heading-Text ohne Markdown-Auszeichnung, case-sensitiv |
| `external.timeout-seconds` | integer | 10 | 1–300 |
| `external.parallel` | integer | 4 | 1–16 |

## 3. Defaults und Konstanten

| Name | Wert | Begründung | Bezug |
|---|---|---|---|
| `DEFAULT_SCAN_ROOTS` | `docs/`, `spec/` (rekursiv, optional) + `*.md` der Repo-Wurzel | [`DC-FA-SCAN-001`](lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln) | — |
| `SKIP_DIRS` | `.git`, `node_modules`, `build`, `target`, `dist`, `vendor`, `.venv`, `__pycache__`, `.idea`, `.vscode` | immer übersprungen | [`DC-FA-SCAN-001`](lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln) |
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
| `repo-escape` | links | aufgelöstes Ziel verlässt die Repo-Wurzel |
| `symlink` | links | Ziel ist/enthält Symlink (Vorrang vor `repo-escape`) |
| `anchor-missing` | anchors | Anker entspricht keinem Heading-Slug |
| `id-unlinked` | ids | Kennung im Fließtext ohne Markdown-Link |
| `matrix-forbidden` | matrix | Referenz zwischen Klassen nicht erlaubt |
| `matrix-inactive` | matrix | Referenz auf Dokument mit verbotenem Status |
| `external-status` | external | HTTP-Status ≥ 400 |
| `external-timeout` | external | Timeout überschritten |
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
| 2026-06-10 | Referenzrichtungs-Korrektur: ADR-Abwärtsverweise entfernt — Spec-Straten verweisen nie abwärts; Traceability über die `Schärft:`-Felder der ADRs (Kurs-Baseline-Korrektur, MR-006) | slice-002 |
