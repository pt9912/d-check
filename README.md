# d-check

Doc-Referenz-Checker für Markdown-Dokumentation: prüft lokale Links
und Bildreferenzen, Heading-Anker, Linkpflicht für Kennungen (z. B.
`ADR-NNNN`) und Referenzrichtungs-Regeln zwischen Dokumentklassen
(Referenzmatrix). `d-check` konsolidiert zwölf funktional überlappende
Einzeltools (`check_refs.py`, `docs-check.js`, `verify-doc-refs.sh`)
aus den Schwester-Repositories des Entwicklungs-Workspace in ein
konfigurierbares Tool.

**Status: MVP läuft** — die Module `links` und `anchors` prüfen
produktiv (Dogfooding: d-check validiert die eigene Doku); `ids`,
`matrix` und `external` entstehen in welle-03. Verbindlich ist das
[Lastenheft](spec/lastenheft.md).

## Geplante Nutzung

Verteilung als Container-Image über GHCR
([`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image)):

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:<tag>
```

Regelmodule (`links`, `anchors`, `ids`, `matrix`, `external`) werden
per CLI-Option oder `.d-check.yml` pro Repo aktiviert.

## Einstieg

| Dokument | Inhalt |
|---|---|
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
