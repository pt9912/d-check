# d-check

Doc-Referenz-Checker für Markdown-Dokumentation: prüft lokale Links
und Bildreferenzen, Heading-Anker, Linkpflicht für Kennungen (z. B.
`ADR-NNNN`) und Referenzrichtungs-Regeln zwischen Dokumentklassen
(Referenzmatrix). `d-check` konsolidiert zwölf funktional überlappende
Einzeltools (`check_refs.py`, `docs-check.js`, `verify-doc-refs.sh`)
aus den Schwester-Repositories des Entwicklungs-Workspace in ein
konfigurierbares Tool.

**Status: released** — die Regelmodule `links`, `anchors`, `ids`,
`matrix` und `external` sind implementiert und getestet; seit
`v0.1.0` als Image auf GHCR (Dogfooding: d-check validiert die eigene
Doku bei jedem Gate-Lauf). Das sechste Modul `codepaths` (Change
Request 0.3.0) ist in Umsetzung. Verbindlich ist das
[Lastenheft](spec/lastenheft.md).

## Nutzung

Verteilung als Container-Image über GHCR
([`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image)):

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.1.0
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
