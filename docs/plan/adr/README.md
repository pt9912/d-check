# ADR-Index

Architecture Decision Records dieses Repos. Konventionen:

- **Dateiname:** `<NNNN>-<kurzer-titel-kebab>.md` (vierstellig, zero-padded).
- **Status:** `Proposed` → `Accepted`; danach immutable. Ablösung nur
  via neue ADR mit `Supersedes ADR-NN` (Status der alten wird
  `Superseded by ADR-NN`).
- Jede ADR deklariert im `**Schärft:**`-Feld aufwärts, welche
  Spec-Stelle sie verbindlich macht (nie das Lastenheft).
- Neue ADRs werden in der Tabelle unten ergänzt.

| ID | Titel | Status | Datum | Bezug |
|---|---|---|---|---|
| [ADR-0001](0001-implementierungssprache.md) | Implementierungssprache: Go | Accepted | 2026-06-10 | [`DC-QA-01`](../../../spec/lastenheft.md#dc-qa-01--performance), [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit), [`DC-FA-DIST-001`](../../../spec/lastenheft.md#dc-fa-dist-001--docker-image) |
| [ADR-0002](0002-distribution-ghcr-image.md) | Distribution: GHCR-Image auf distroless/static | Accepted | 2026-06-10 | [`DC-FA-DIST-001`](../../../spec/lastenheft.md#dc-fa-dist-001--docker-image), ADR-0001 |
| [ADR-0003](0003-config-format.md) | Config-Parsing: striktes YAML via gopkg.in/yaml.v3 | Accepted | 2026-06-10 | [`DC-FA-CONF-001`](../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei), ADR-0001 |
| [ADR-0004](0004-architektur-pattern-hexagonal.md) | Architektur-Pattern: Hexagon light | Accepted (Pfad-Festlegung superseded durch ADR-0005) | 2026-06-10 | [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus), [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit), [`DC-FA-CLI-002`](../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl) |
| [ADR-0005](0005-modul-layout-hexagon-ordner.md) | Modul-Layout: hexagon-/adapter-Ordner nach u-boot-Konvention | Accepted | 2026-06-10 | ADR-0004, [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) |
| [ADR-0006](0006-lint-profil-solid.md) | Lint-Profil: SOLID-nah nach u-boot-Vorbild (ohne depguard) | Accepted | 2026-06-10 | ADR-0001, ADR-0005 |
| [ADR-0007](0007-repository-lizenz-mit.md) | Repository-Lizenz: MIT | Accepted | 2026-06-11 | Repository-Veröffentlichung und Weiterverwendung |
| [ADR-0008](0008-reparatur-ableitbarkeit.md) | Reparatur-Modus: nur deterministisch ableitbare Fixes; Best-Guess review-pflichtig | Accepted | 2026-06-19 | [`DC-FA-CLI-008`](../../../spec/lastenheft.md#dc-fa-cli-008--reparatur-patch), [`DC-FA-CLI-007`](../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus) |
| [ADR-0009](0009-yaml-im-report-adapter.md) | `gopkg.in/yaml.v3` auch im report-Adapter (Reporter-Serialisierung) | Accepted | 2026-06-19 | [ADR-0005](0005-modul-layout-hexagon-ordner.md), [`DC-FA-CLI-004`](../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate) |
| [ADR-0010](0010-semgrep-hermetisches-gate.md) | semgrep als hermetisches Gate (lokal gecachtes, gepinntes `go/lang/security`-Regelset) | Accepted | 2026-06-20 | [ADR-0006](0006-lint-profil-solid.md), [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus), [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) |
| [ADR-0011](0011-digest-pins-build-gate-images.md) | Digest-Pins aller Build-/Gate-Images (Dockerfile-`FROM` + semgrep) | Accepted | 2026-06-20 | [ADR-0002](0002-distribution-ghcr-image.md), [ADR-0006](0006-lint-profil-solid.md), [ADR-0010](0010-semgrep-hermetisches-gate.md), [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus) |
| [ADR-0012](0012-kern-paketschnitt-model-rules-app.md) | Kern-Paketschnitt: `model`/`rules`/`app` (drei Pakete statt einem) | Accepted | 2026-06-20 | [ADR-0004](0004-architektur-pattern-hexagonal.md), [ADR-0005](0005-modul-layout-hexagon-ordner.md), [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) |
| [ADR-0013](0013-pr-ci-und-traceability-gate.md) | PR-/Push-CI und Traceability-Gate (DC-*/ADR-* in Commits) | Proposed | 2026-06-21 | [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus), [`DC-FA-DIST-001`](../../../spec/lastenheft.md#dc-fa-dist-001--docker-image), [ADR-0011](0011-digest-pins-build-gate-images.md), [ADR-0002](0002-distribution-ghcr-image.md) |
