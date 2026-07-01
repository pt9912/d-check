# d-check

> 🇬🇧 **English** · 🇩🇪 [Deutsch](README.de.md)

Documentation reference checker for Markdown — deterministic,
side-effect-free, shipped as a container image.

**Status: released** — all fifteen rule modules (`links`, `anchors`, `ids`,
`matrix`, `codepaths`, `spans`, `hostpaths`, `diagrams`, `versions`, `pins`,
`immutable`, `vcs`, `commits`, `planning`, `external`) are in the GHCR image.
The authoritative source is the [requirements spec](spec/lastenheft.md); the
most recent changes (most recently the hermetic opt-in module `planning` for
roadmap-↔-in-progress lifecycle consistency — the last mechanized gate script
of the `tools/*.sh` audit) are tracked in [CHANGELOG.md](CHANGELOG.md).

## What is d-check?

**d-check** checks Markdown documentation as a verifiable **invariant network**:
every machine-decidable documentation invariant is an individually enablable rule
module with its own requirement in the [requirements spec](spec/lastenheft.md) —
from the **reference network** (links, anchors, ID link obligations, reference
matrix) through Markdown hygiene (span artifacts, host-path leaks), content drift
and immutability (content/core pins, git diff) to version-pin, commit-traceability
and planning-lifecycle consistency:

- `links` — local link and image references: target exists, no
  repo escape ([`DC-FA-LINK-001`](spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links))
- `anchors` — heading anchors (GitHub slug procedure) and inline HTML anchors
  (`<a name>`, `id=`)
  ([`DC-FA-ANCH-001`](spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors))
- `ids` — link obligation for identifiers (e.g. `ADR-NNNN`) per
  declared patterns ([`DC-FA-ID-001`](spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids))
- `matrix` — reference-direction rules between document classes
  ([`DC-FA-MTX-001`](spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)),
  **within** an ordered class (`order`/`direction` ⇒ `matrix-downward`,
  [`DC-FA-MTX-002`](spec/lastenheft.md#dc-fa-mtx-002--verweisrichtung-innerhalb-einer-geordneten-dokumentklasse-modul-matrix))
  and also as a **bare ID token** in the body (`token` ⇒ `matrix-forbidden`,
  [`DC-FA-MTX-003`](spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix))
  plus status conditions; exceptions are structural (`exclude-sections`,
  `allow-supersede-lineage`) or declared via the provenance marker
  `<!-- d-check:status-provenance -->`
- `external` — reachability of external URLs, strictly opt-in
  ([`DC-FA-EXT-001`](spec/lastenheft.md#dc-fa-ext-001--externe-links-modul-external-opt-in))
- `codepaths` — explicit paths in inline code, opt-in
  ([`DC-FA-CODE-001`](spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in))
- `spans` — Markdown span artifacts (unclosed code spans,
  nested links), opt-in
  ([`DC-FA-SPAN-001`](spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in))
- `hostpaths` — host-local absolute paths (machine-layout leaks),
  opt-in
  ([`DC-FA-HOST-001`](spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in))
- `diagrams` — identifier existence in diagram fences (e.g. `mermaid`): every
  identifier found in the diagram must be defined in its `defined-in` source,
  opt-in
  ([`DC-FA-DIAG-001`](spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in))
- `versions` — version-pin consistency: pinned `ghcr` image references must
  carry the current version (from `version.md#aktuell`), also reads
  fenced code, opt-in
  ([`DC-FA-VER-001`](spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in))
- `pins` — content pin against content drift: a link with
  `<!-- dpin: … -->` is checked against the hash of its target span (finding
  `link-stale` on drift), opt-in per link
  ([`DC-FA-PIN-001`](spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in))
- `immutable` — immutability pin against core drift: a file with
  `<!-- immutable: … -->` is checked against the hash of its normalized **core**
  (without the marker line + `exclude-sections`) (finding `core-drift` on
  drift), opt-in per file; hermetic (no git, read-only working tree)
  ([`DC-FA-IMM-001`](spec/lastenheft.md#dc-fa-imm-001--immutabilitäts-pin-gegen-core-drift-modul-immutable-opt-in))
- `vcs` — git-diff immutability of the core over a commit range: mechanizes the
  ADR immutability as a distributable module (`core-drift-vcs`), pure-Go git in
  the read-only `.git` (**no** git binary, **no** network), opt-in
  ([`DC-FA-VCS-001`](spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in))
- `commits` — traceability identifier in commit messages over a range (`--range`)
  or the pending message (`--commit-msg`): every commit message carries a
  `DC-`/`ADR-`/`MR-`/`slice-` ID (`commit-untraceable`), shares the VCS port with
  `vcs`, opt-in
  ([`DC-FA-COMMITS-001`](spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in))
- `planning` — roadmap-↔-in-progress lifecycle consistency: the idle marker sits
  in the `## Aktuelle Welle` block exactly when no `slice-*` is in the directory
  (`planning-drift`); hermetic (no git), fail-closed on a missing/ambiguous
  heading, opt-in
  ([`DC-FA-PLAN-001`](spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in))

Every finding names file, line, target and reason; exit codes:
`0` clean, `1` findings, `2` environment or configuration error.

## Why d-check?

Twelve functionally overlapping tool copies from three families
(`verify-doc-refs.sh` — Shell, `check_refs.py` — Python,
`docs-check.js` — JavaScript) grew across the sister repositories of the
development workspace — each with its own feature set, its own drift, its
own maintenance. d-check replaces them with **one** tool:

- **Configuration instead of fork:** repo-specific behavior lives
  declaratively in `.d-check.yml`
  ([`DC-FA-CONF-001`](spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)),
  not in copied scripts.
- **Replaceability is measurable:** every legacy tool must be replaceable by
  d-check with matching configuration — at least the same real findings, no
  false positives that break a green CI
  ([`DC-QA-04`](spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)).
- **One distribution path:** container image with digest pin instead of
  n maintained copies
  ([`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image)).

## Core idea

**Documentation is a reference graph with verifiable invariants.**
Whether a link reaches its target, an anchor exists, an identifier links
to its definition, or a document class points downward is
machine-decidable — d-check turns these invariants into a gate
instead of a review opinion.

The rule: **report, never repair.** d-check is a pure read
tool ([`DC-QA-03`](spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit));
deterministic findings are fixed, not suppressed. An opt-out
marker exists only where a non-existent target or an illustrative identifier
can be documented intent (`d-check:ignore`, per line) — it silences
exclusively the modules `codepaths` and `ids`
([`DC-FA-CODE-001`](spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in),
[`DC-FA-ID-001`](spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)).

## What makes it trustworthy?

**A checking tool whose own statements wobble is worthless** —
determinism and side-effect freedom are therefore core contracts of the
spec, and both are measured, not asserted:

- **Determinism:** identical input ⇒ byte-identical output,
  stably sorted; tested over ten repeated runs with
  hash comparison
  ([`DC-QA-02`](spec/lastenheft.md#dc-qa-02--determinismus)).
- **Side-effect-free and network-less:** never writes into the checked
  repository, opens no network connections except in the opt-in module
  `external` — measured in the gate run with a read-only mount
  and `--network none`
  ([`DC-QA-03`](spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **No silent defaults:** every invalid `.d-check.yml` aborts
  with exit 2; checking never proceeds with guessed configuration
  ([`DC-FA-CONF-001`](spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)).
- **Dogfooding:** d-check validates its own docs on every
  gate run — with the [self-configuration](.d-check.yml) fully
  built out (eight modules incl. reference matrix, span artifacts,
  host-path hygiene and version-pin consistency).
- **Container native-identical:** the image's finding output and exit code
  are byte-identical to native execution, tested automatically
  ([`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image));
  CI consumption via digest pin.

## Usage

Distributed as a container image via GHCR
([`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image)):

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.36.0
```

CI pipelines pin to the digest from the release notes rather than to
moving tags — details, options and exit codes:
[`docs/user/operations.md`](docs/user/operations.md) and
[`docs/user/releasing.md`](docs/user/releasing.md).

## Configuration (`.d-check.yml`)

Optional in the repo root; without a file the default modules
`links` + `anchors` run over `docs/`, `spec/` and the root `*.md`:

```yaml
scan:
  roots: ["."]                  # entire repo root
modules: [links, anchors, ids]  # external stays strictly opt-in
ids:
  patterns:
    - regex: 'ADR-\d{4}'
      target: docs/adr/         # identifiers must link here
```

The full schema with all keys, defaults and
validation constraints is in the
[specification §`.d-check.yml`](spec/spezifikation.md#d-checkyml);
a living example fully built out (incl. reference matrix) is
[this repo's self-configuration](.d-check.yml). Every invalid
configuration aborts with exit 2 — checking never proceeds with
silent defaults
([`DC-FA-CONF-001`](spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)).

## Getting started

| Document | Contents |
|---|---|
| [`docs/user/operations.md`](docs/user/operations.md) | Invocation reference: options, exit codes, configuration |
| [`docs/user/releasing.md`](docs/user/releasing.md) | Release process, digest-pin consumption |
| [`spec/lastenheft.md`](spec/lastenheft.md) | Requirements (`DC-FA-*`, `DC-QA-*`), acceptance criteria |
| [`harness/README.md`](harness/README.md) | Harness entry: source precedence, guides, sensors |
| [`AGENTS.md`](AGENTS.md) | Briefing for AI coding agents, hard rules |
| [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md) | Waves and slices |
| [`CHANGELOG.md`](CHANGELOG.md) | Change history |

## Development

The host needs only `git`, GNU `make`, `bash` and Docker
(see [`AGENTS.md`](AGENTS.md) §3.1).

```bash
make help     # available targets
make gates    # all inner gates (mandatory before handoff)
```

## License

This project is licensed under the [MIT License](LICENSE).
