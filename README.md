# d-check

Doc-Referenz-Checker für Markdown-Dokumentation — deterministisch,
seiteneffektfrei, ausgeliefert als Container-Image.

**Status: released** — alle elf Regelmodule (`links`, `anchors`, `ids`,
`matrix`, `codepaths`, `spans`, `hostpaths`, `diagrams`, `versions`, `pins`,
`external`) sind im GHCR-Image. Verbindlich ist das [Lastenheft](spec/lastenheft.md);
die jeweils jüngsten Änderungen (zuletzt die token-basierte Referenz-Richtung im
Modul `matrix`) führt die [CHANGELOG.md](CHANGELOG.md).

## Was ist d-check?

**d-check** prüft Markdown-Dokumentation als Referenz-Netz: Jede prüfbare
Verweis- oder Kennungs-Invariante ist ein einzeln aktivierbares Regelmodul mit
eigener Anforderung im [Lastenheft](spec/lastenheft.md):

- `links` — lokale Link- und Bildreferenzen: Ziel existiert, kein
  Repo-Escape ([`DC-FA-LINK-001`](spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links))
- `anchors` — Heading-Anker (GitHub-Slug-Verfahren) und Inline-HTML-Anker
  (`<a name>`, `id=`)
  ([`DC-FA-ANCH-001`](spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors))
- `ids` — Linkpflicht für Kennungen (z. B. `ADR-NNNN`) nach
  deklarierten Mustern ([`DC-FA-ID-001`](spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids))
- `matrix` — Referenzrichtungs-Regeln zwischen Dokumentklassen
  ([`DC-FA-MTX-001`](spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)),
  **innerhalb** einer geordneten Klasse (`order`/`direction` ⇒ `matrix-downward`,
  [`DC-FA-MTX-002`](spec/lastenheft.md#dc-fa-mtx-002--verweisrichtung-innerhalb-einer-geordneten-dokumentklasse-modul-matrix))
  und auch als **bare ID-Token** im Körper (`token` ⇒ `matrix-forbidden`,
  [`DC-FA-MTX-003`](spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix))
  plus Status-Bedingungen; Ausnahmen sind strukturell (`exclude-sections`,
  `allow-supersede-lineage`) oder per Provenance-Marker
  `<!-- d-check:status-provenance -->` deklariert
- `external` — Erreichbarkeit externer URLs, strikt opt-in
  ([`DC-FA-EXT-001`](spec/lastenheft.md#dc-fa-ext-001--externe-links-modul-external-opt-in))
- `codepaths` — explizite Pfade in Inline-Code, opt-in
  ([`DC-FA-CODE-001`](spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in))
- `spans` — Markdown-Span-Artefakte (ungeschlossene Code-Spans,
  verschachtelte Links), opt-in
  ([`DC-FA-SPAN-001`](spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in))
- `hostpaths` — host-lokale absolute Pfade (Maschinen-Layout-Leaks),
  opt-in
  ([`DC-FA-HOST-001`](spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in))
- `diagrams` — Kennungs-Existenz in Diagramm-Fences (z. B. `mermaid`): jede
  im Diagramm gefundene Kennung muss in ihrer `defined-in`-Quelle definiert
  sein, opt-in
  ([`DC-FA-DIAG-001`](spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in))
- `versions` — Versions-Pin-Konsistenz: gepinnte `ghcr`-Image-Verweise müssen
  die aktuelle Version (aus `version.md#aktuell`) tragen, liest auch
  Fenced-Code, opt-in
  ([`DC-FA-VER-001`](spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in))
- `pins` — Content-Pin gegen inhaltlichen Drift: ein Link mit
  `<!-- dpin: … -->` wird gegen den Hash seines Ziel-Spans geprüft (Befund
  `link-stale` bei Drift), opt-in pro Link
  ([`DC-FA-PIN-001`](spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in))

Jeder Befund nennt Datei, Zeile, Ziel und Grund; Exit-Codes:
`0` sauber, `1` Befunde, `2` Umgebungs- oder Konfigurationsfehler.

## Warum d-check?

Zwölf funktional überlappende Tool-Kopien aus drei Familien
(`verify-doc-refs.sh` — Shell, `check_refs.py` — Python,
`docs-check.js` — JavaScript) sind in den Schwester-Repositories des
Entwicklungs-Workspace gewachsen — jede mit eigenem Funktionsstand,
eigener Drift, eigener Pflege. d-check ersetzt sie durch **ein** Tool:

- **Konfiguration statt Fork:** repo-spezifisches Verhalten lebt
  deklarativ in `.d-check.yml`
  ([`DC-FA-CONF-001`](spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)),
  nicht in kopierten Skripten.
- **Ersetzbarkeit ist messbar:** Jedes Alt-Tool muss durch d-check
  mit passender Konfiguration ablösbar sein — mindestens dieselben
  echten Befunde, keine False-Positives, die eine grüne CI brechen
  ([`DC-QA-04`](spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)).
- **Ein Distributionsweg:** Container-Image mit Digest-Pin statt
  n gepflegter Kopien
  ([`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image)).

## Kerngedanke

**Dokumentation ist ein Referenz-Graph mit prüfbaren Invarianten.**
Ob ein Link sein Ziel erreicht, ein Anker existiert, eine Kennung auf
ihre Definition verlinkt oder eine Dokumentklasse abwärts zeigt, ist
maschinell entscheidbar — d-check macht diese Invarianten zum Gate
statt zur Review-Meinung.

Dabei gilt: **berichten, nie reparieren.** d-check ist ein reines
Lese-Tool ([`DC-QA-03`](spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit));
deterministische Befunde werden behoben, nicht unterdrückt. Einen
Opt-out-Marker gibt es nur dort, wo ein nicht existierendes Ziel oder
eine illustrative Kennung dokumentierte Absicht sein kann
(`d-check:ignore`, zeilenweise) — er stellt ausschließlich die Module
`codepaths` und `ids` still
([`DC-FA-CODE-001`](spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in),
[`DC-FA-ID-001`](spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)).

## Was macht es vertrauenswürdig?

**Ein Prüf-Tool, dessen eigene Aussagen wackeln, ist wertlos** —
Determinismus und Seiteneffektfreiheit sind deshalb Kernverträge des
Lastenhefts, und beide werden gemessen, nicht behauptet:

- **Determinismus:** identische Eingabe ⇒ byte-identische Ausgabe,
  stabil sortiert; getestet über zehn wiederholte Läufe mit
  Hash-Vergleich
  ([`DC-QA-02`](spec/lastenheft.md#dc-qa-02--determinismus)).
- **Seiteneffektfrei und netzlos:** schreibt nie in das geprüfte
  Repository, öffnet außer im opt-in-Modul `external` keine
  Netzwerkverbindungen — gemessen im Gate-Lauf mit read-only-Mount
  und `--network none`
  ([`DC-QA-03`](spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **Keine stillen Defaults:** Jede ungültige `.d-check.yml` bricht
  mit Exit 2 ab; geprüft wird nie mit geratener Konfiguration
  ([`DC-FA-CONF-001`](spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)).
- **Dogfooding:** d-check validiert die eigene Doku bei jedem
  Gate-Lauf — mit der [Selbstkonfiguration](.d-check.yml) im
  Vollausbau (acht Module inkl. Referenzmatrix, Span-Artefakten,
  Host-Pfad-Hygiene und Versions-Pin-Konsistenz).
- **Container nativ-identisch:** Befund-Ausgabe und Exit-Code des
  Images sind byte-identisch zur nativen Ausführung, automatisiert
  getestet ([`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image));
  CI-Konsum per Digest-Pin.

## Nutzung

Verteilung als Container-Image über GHCR
([`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image)):

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.31.0
```

CI-Pipelines pinnen auf den Digest aus den Release-Notes statt auf
bewegliche Tags — Details, Optionen und Exit-Codes:
[`docs/user/operations.md`](docs/user/operations.md) und
[`docs/user/releasing.md`](docs/user/releasing.md).

## Konfiguration (`.d-check.yml`)

Optional in der Repo-Wurzel; ohne Datei laufen die Default-Module
`links` + `anchors` über `docs/`, `spec/` und die Root-`*.md`:

```yaml
scan:
  roots: ["."]                  # gesamte Repo-Wurzel
modules: [links, anchors, ids]  # external bleibt strikt opt-in
ids:
  patterns:
    - regex: 'ADR-\d{4}'
      target: docs/adr/         # Kennungen müssen hierhin verlinken
```

Das vollständige Schema mit allen Schlüsseln, Defaults und
Validierungs-Constraints steht in der
[Spezifikation §`.d-check.yml`](spec/spezifikation.md#d-checkyml);
ein lebendes Beispiel im Vollausbau (inkl. Referenzmatrix) ist die
[Selbstkonfiguration dieses Repos](.d-check.yml). Jede ungültige
Konfiguration bricht mit Exit 2 ab — geprüft wird nie mit
stillschweigenden Defaults
([`DC-FA-CONF-001`](spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)).

## Einstieg

| Dokument | Inhalt |
|---|---|
| [`docs/user/operations.md`](docs/user/operations.md) | Aufruf-Referenz: Optionen, Exit-Codes, Konfiguration |
| [`docs/user/releasing.md`](docs/user/releasing.md) | Release-Prozess, Digest-Pin-Konsum |
| [`spec/lastenheft.md`](spec/lastenheft.md) | Anforderungen (`DC-FA-*`, `DC-QA-*`), Akzeptanzkriterien |
| [`harness/README.md`](harness/README.md) | Harness-Einstieg: Source Precedence, Guides, Sensors |
| [`AGENTS.md`](AGENTS.md) | Briefing für AI-Coding-Agenten, Hard Rules |
| [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md) | Wellen und Slices |
| [`CHANGELOG.md`](CHANGELOG.md) | Änderungshistorie |

## Entwicklung

Der Host braucht nur `git`, GNU `make`, `bash` und Docker
(siehe [`AGENTS.md`](AGENTS.md) §3.1).

```bash
make help     # verfügbare Targets
make gates    # alle inneren Gates (mandatory vor Handoff)
```

## Lizenz

Dieses Projekt steht unter der [MIT-Lizenz](LICENSE).
