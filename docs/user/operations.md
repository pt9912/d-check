# Operations — d-check Aufruf-Referenz

Kurzreferenz für den Betrieb; verbindliche Verträge stehen im
[Lastenheft](../../spec/lastenheft.md) (Anforderungen) und in der
[Spezifikation](../../spec/spezifikation.md) (Schemata, Defaults,
Grund-Codes) — diese Seite dupliziert sie nicht.

## Aufruf

```sh
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check@sha256:<digest>
```

Das Image prüft das nach `/repo` gemountete Repository (read-only
genügt — das Tool schreibt nie,
[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
CLI-Optionen werden als Container-Argumente angehängt.

## Optionen

| Option | Wirkung |
|---|---|
| `--enable <modul>` / `--disable <modul>` | Regelmodule zu-/abschalten (`links`, `anchors`, `ids`, `matrix`, `external`, `codepaths`, `spans`, `hostpaths`); CLI schlägt Konfiguration ([`DC-FA-CLI-002`](../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)) |
| `--json` | maschinenlesbare Gesamt-Ausgabe ([Schema](../../spec/spezifikation.md)) |
| `--yaml` | wie `--json`, aber als YAML — gleiche Struktur (`findings`/`summary`/`exitCode`), nur Serialisierung; `--json` und `--yaml` schließen sich aus ([`DC-FA-CLI-004`](../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)) |
| `--doctor` | erklärende, nach Datei gruppierte Diagnose auf stdout statt der Befund-Zeilen, mit Fix-Kandidaten wo eindeutig (v1: `id-unlinked` → Definitions-Link) — **liest, schreibt nichts**; mit `--json` oder `--yaml` kombinierbar (maschinenlesbare Diagnose: `findings` zusätzlich mit `reasonText`/`fixCandidate`), nicht mit `--repair` ([`DC-FA-CLI-007`](../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)) |
| `--repair` | **konservativer** Reparatur-Patch (unified diff) auf stdout, `git apply`-kompatibel — nur eindeutige Fixes (v1: `id-unlinked` auf nackte Prosa-Vorkommen). **Liest, schreibt nichts**; Anwenden via `d-check --repair > fix.patch && git apply fix.patch`. Nicht mit `--json`/`--yaml`/`--doctor` kombinierbar ([`DC-FA-CLI-008`](../../spec/lastenheft.md#dc-fa-cli-008--reparatur-patch)) |
| `--repair-broad` | wie `--repair`, zusätzlich **Best-Guess** (z. B. `target-missing` → Datei gleichen Basisnamens) — review-pflichtig, die Markierung erscheint auf **stderr** (Patch auf stdout bleibt `git apply`-rein) |
| `--print-config` | kommentiertes `.d-check.yml`-Startgerüst auf stdout, dann Exit 0 — **kein Scan, schreibt nichts**; Anlegen via Umleitung: `d-check --print-config > .d-check.yml` ([`DC-FA-CLI-005`](../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben)) |
| `--suggest-config <quelle>[,…]` | liest die benannten Autoritäts-Quellen und schlägt ein `.d-check.yml` vor (abgeleitete `ids`-Muster + opt-in-Module nach Signal) — **liest, schreibt nichts**; Umleiten via `> .d-check.yml`. Die reservierten Quellen `ai-harness` (repo-bewusst) und `ai-harness-init` (Voll-Kanon fürs leere Repo) schlagen stattdessen eine an die ai-harness-course-Konvention angelehnte Vorlage vor (kanonische `ids`-Muster, `matrix` samt Referenzrichtung, Standard-Modulset; Anforderungs-Präfix via `--id-prefix` bzw. abgeleitet, sonst Platzhalter) ([`DC-FA-CLI-006`](../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)) |
| `--trace` | **Requirements Traceability Matrix** auf stdout: je Anforderung (`DC-*` im Lastenheft) die referenzierenden ADRs/Slices + Waisen-Markierung — **liest, schreibt nichts**, kein Dokument erzeugt; Default Markdown-Tabelle, mit `--json`/`--yaml` maschinenlesbar; nicht mit `--doctor`/`--repair` ([`DC-FA-CLI-009`](../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)) |

Default-Module ohne Konfiguration: `links` + `anchors`. Das Modul
`external` ist strikt opt-in (einzige Netzwerk-Tür).

## Exit-Codes

`0` = keine Befunde, `1` = mindestens ein Befund, `2` = Nutzungs-
oder Umgebungsfehler
([`DC-FA-CLI-003`](../../spec/lastenheft.md#dc-fa-cli-003--exit-codes)).

## Konfiguration

Optionale `.d-check.yml` in der Repo-Wurzel; Schema und Defaults in
der [Spezifikation §`.d-check.yml`](../../spec/spezifikation.md#d-checkyml),
Vollvalidierung gemäß
[`DC-FA-CONF-001`](../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(jeder Konfigurationsfehler → Exit 2). Ein Beispiel im Vollausbau ist
die [Selbstkonfiguration dieses Repos](../../.d-check.yml).

## Linkdichte erzwingen — `link-policy: always`

Das Modul `ids` prüft per Default nur, dass *nackte* Kennungen im
Fließtext verlinkt sind; eine Kennung in Inline-Code (`` `…` ``) bleibt
frei. Wer **gut verlinkte Dokumente** als gemessenes Property will
(jede Kennung eine navigierbare Referenz), setzt pro Muster
`link-policy: always` — dann ist auch eine Kennung in Inline-Code
linkpflichtig
([`DC-FA-ID-001`](../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)):

```yaml
ids:
  patterns:
    - regex: 'ADR-\d{4}'
      target: docs/plan/adr/
      link-policy: always            # prose (Default) | always
      exempt-paths: [CHANGELOG.md, "docs/reviews/**"]
```

Zwei Ventile nehmen eine Datei bzw. Zeile von der Linkpflicht eines
Musters aus — für **alle** Vorkommen (nackt im Fließtext wie in
Inline-Code) und unabhängig von der `link-policy`:

- **`exempt-paths`** (Glob-Liste, Syntax wie `scan.ignore`): Dateien
  ohne Linkpflicht für das Muster — typischerweise literal-schwere
  Artefakte wie Changelogs oder Review-Reports. Gilt gleich, ob die
  Kennung dort nackt oder in Backticks steht.
- **`d-check:ignore`** (HTML-Kommentar auf der Zeile, Begründung
  empfohlen): nimmt die ganze Zeile aus — für bewusst illustrative
  Beispiel-Kennungen. Der Marker wirkt auf `ids` und `codepaths`.

Beide Ventile gelten auch unter der Default-Politik `prose`: eine nackte
Kennung in einem ausgenommenen Artefakt (etwa einem Review-Report) lässt
sich so stummschalten, ohne `always` zu aktivieren.

**Bewusst opt-in:** Der Default bleibt `prose`, damit bestehende
Konfigurationen byte-identisch laufen und kein Repo ungefragt rote
Läufe bekommt. Das heißt **nicht**, dass ungenügende Verlinkung
unsichtbar bleibt: Sie zu *entdecken* ist eine Bringschuld des
Betreibers — ein `always`-Lauf über die geprüften Repos zeigt die
Lücken, unabhängig davon, ob ein Repo `always` schon aktiviert hat
([`DC-QA-04`](../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Muster).

## Referenzmatrix — Supersede-Lineage und Marker-Politik

Das Modul `matrix`
([`DC-FA-MTX-001`](../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix))
meldet Referenzen auf inaktive Dokumente (Status `superseded`/
`deprecated`) als `matrix-inactive`. Eine **ablösende** Datei verweist
aber per Definition auf das Dokument, das sie ablöst — diese
Lineage-Kante ist legitim. Die opt-in Konfiguration nimmt sie aus:

```yaml
matrix:
  status:
    forbidden: [superseded, deprecated]
    allow-supersede-lineage: true
    supersede-fields: [Supersedes, Aenderungstyp]
```

Genau die Kante X → Y wird von der Status-Prüfung ausgenommen, wenn X
über eines der `supersede-fields` (Form `**Feld:** Wert` oder
`Feld: Wert`) deklariert, dass es Y ablöst — erkannt am Linktext oder
Zielpfad der Referenz. Alle anderen Referenzen auf Y bleiben
`matrix-inactive`, und die Klassen-Regeln (`matrix-forbidden`) sind
unberührt. Default aus: ohne das Flag ist der Befundsatz byte-identisch.

**Kein Zeilen-Marker für `matrix`.** Der `d-check:ignore`-Marker wirkt
ausschließlich auf `ids` und `codepaths` (illustrative Beispiele).
`matrix`-Befunde werden behoben oder **strukturell** ausgenommen
(`exclude-sections` für Provenance-/Historie-Sektionen,
`allow-supersede-lineage` für die Lineage-Kante) — legitime Ausnahmen
sind deklarierte Konfiguration, keine verstreuten Kommentare.

## Config vorschlagen — `--suggest-config`

`d-check --suggest-config spec/lastenheft.md,harness/conventions.md,docs/plan/adr/`
liest die benannten **Autoritäts-Quellen** (Dateien/Verzeichnisse, in
denen Kennungen *definiert* sind) und gibt ein vorgeschlagenes
`.d-check.yml` auf stdout aus — es liest das Repo, **schreibt aber nie**
(read-only-Vertrag; Umleiten macht der Aufrufer). Je Quelle wird ein
`ids`-Muster abgeleitet, dessen `regex` alle dort gefundenen Kennungen
matcht (die Quell-Kennungen stehen als Kommentar dabei); zusätzlich
werden opt-in-Module vorgeschlagen, die echtes Signal liefern.

**Harness-Vorlage (`ai-harness` / `ai-harness-init`):** Statt Pfaden geben
sie eine an die ai-harness-course-Konvention (Baseline `v1.3.0`) angelehnte
Vorlage aus: die kanonischen `ids`-Muster (`ADR-`, `MR-`, `DC-`, `slice`),
die `matrix`-Klassen samt Referenzrichtung und das Standard-Modulset. Zwei
Modi (das passende ist nicht auto-erkennbar — Henne-Ei):

- **`ai-harness-init`** gibt den **Voll-Kanon** aus (alle Blöcke aktiv) —
  Zielbild fürs **leere Repo**; läuft, sobald die Struktur (Scan-Wurzeln,
  `ids`-Targets) existiert.
- **`ai-harness`** ist **repo-bewusst** — nur existierende Pfade aktiv,
  fehlende auskommentiert mit Hinweis; läuft sofort gegen ein
  **bestehendes** Repo.

Beide read-only und mit echten Quellen kombinierbar.

**Scaffold, kein Orakel — die Grenze ist bewusst:** Erkannt werden
Kennungen, die als **führendes Token einer Überschrift** in
Großschreibung definiert sind (`### DC-FA-ID-001 — …`, `# ADR-0001 — …`). <!-- d-check:ignore (Format-Beispiele) -->
**Nicht** erkannt werden kleingeschriebene IDs (`slice-001`) und solche,
die nicht in Überschriften definiert sind (z. B. in Tabellen). Der
Vorschlag ist eine **Best-Guess-Ableitung**, die der Mensch prüft,
verengt (etwa `ADR-\d+` → `ADR-\d{4}`) und ergänzt — kein automatisch
verbindlicher Check (die Prüfung läuft stets gegen explizit
konfigurierte Muster, [`DC-FA-ID-001`](../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)).
