# d-check

> 🇬🇧 **English** · 🇩🇪 [Deutsch](README.de.md)

Consistency checker that verifies a repository's documentation against its
actual state — Markdown references (links, images, anchors, identifiers) plus
declarations (build targets, version pins, commit trace, planning, workflow
references).
Deterministic, side-effect-free, shipped as a container image.

## What is d-check?

**d-check** checks Markdown documentation as a verifiable **invariant network**:
every machine-decidable documentation invariant is a rule module that can be
enabled individually, with its own requirement in the
[requirements spec](spec/lastenheft.md) —
from the **reference network** (links, anchors, ID link obligations, reference
matrix) through Markdown hygiene (span artifacts, host-path leaks), content drift
and immutability (content/core pins, git diff) to version-pin, commit-traceability,
planning-lifecycle and tracked-status consistency, up to structure invariants
**within** a document:

- `links` — local link and image references: target exists, no
  repo escape; opt-in `resolve-from`: files in **moving** lifecycle
  directories must resolve every relative target from every location of
  their group — reported **before** the `git mv`
  ([`DC-FA-LINK-001`](spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links))
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
- `sources` — content pin of external sources against upstream drift, opt-in and
  — besides `external` — the **second** network door: a source pinned to a
  `sha256` (via the marker `<!-- source-pin: … -->` on its `http(s)` link or the
  config block `sources:`) is fetched, hashed and compared (`source-drift` with
  the full actual hash, `source-unreachable` on a network failure); single file
  or archive (`unpack: zip`)
  ([`DC-FA-SRC-001`](spec/lastenheft.md#dc-fa-src-001--upstream-content-drift-externer-quellen-modul-sources-opt-in-netz))
- `codepaths` — explicit paths in inline code, opt-in; the opt-in `check-lines`
  verifies `file:<from>-<to>` line references (`citation-out-of-range`,
  `citation-inverted-range`)
  ([`DC-FA-CODE-001`](spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in))
- `spans` — Markdown span artifacts (unclosed code spans,
  nested links), opt-in
  ([`DC-FA-SPAN-001`](spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in))
- `hostpaths` — host-local absolute paths (machine-layout leaks),
  opt-in
  ([`DC-FA-HOST-001`](spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in))
- `diagrams` — identifier existence in diagram fences (e.g. `mermaid`): every
  identifier found in the diagram must be defined in its `defined-in` source,
  opt-in; **both valves** like `ids`/`codepaths` — `exempt-paths` (file-wide)
  and the line marker, here a **token** rather than an HTML comment, and on the
  fence's **opening line** for the whole block
  ([`DC-FA-DIAG-001`](spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in))
- `versions` — version-pin consistency: pinned `ghcr` image references must
  carry the current version (from `version.md#aktuell`), also reads
  fenced code, opt-in; the **anchor** of `current-from` follows the same answer
  as in `anchors`, not the fence exception. `versions.patterns` carries
  **several** pattern/source pairs (your own release **and** a pinned foreign
  version) — the short form is the single-element list, and exemptions are
  per pair
  ([`DC-FA-VER-001`](spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in))
- `pins` — content pin against content drift: a link with
  `<!-- dpin: … -->` is checked against the hash of its target span (finding
  `link-stale` on drift), opt-in per link; **which anchor addresses the span is
  decided by the same answer as in `anchors`** (duplicate slug, percent
  decoding, case-sensitive, outside code only)
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
- `planning` — lifecycle consistency, all three sides. **Entry:** the idle marker sits
  in the `## Aktuelle Welle` block exactly when no `slice-*` is in the directory
  (`planning-drift`). **Exit** (additionally opt-in via `closure.dir`): the
  **structure** of closure notes on completed work items — section present,
  enough sentence terminators outside code blocks, no declared boilerplate
  phrase (`closure-note-missing`/`-thin`/`-boilerplate`). Checks structure, not
  meaning. Hermetic (no git), fail-closed on a missing/ambiguous heading, a
  missing closure directory and on zero candidates. **Wave registers**
  (additionally opt-in via `waves.dir`): the roadmap's wave sections against
  the wave files — active wave ⟺ flat wave document (`waves.mode: one`,
  default) **or** identifier bijection active block ⟺ files with the rest
  marker staying out of it (`waves.mode: many`); preview without a file,
  closed register ⟺ results notes both ways (`wave-drift`,
  `wave-preview-exists`, `wave-results-missing`, `wave-unregistered`), opt-in.
  **Register coverage** (additionally opt-in via `observations.register`): a
  cited observation identifier without a register line
  (`observation-unregistered`) — prose and link text count, a plain
  inline-code example does not
  ([`DC-FA-PLAN-001`](spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in))
- `tracked` — tracked status of resolvable, **existing** link/image targets
  against the git **index** (`target-untracked`: an untracked/gitignored
  target is missing on every fresh clone); index truth (staged = tracked, no
  `.gitignore` interpretation), reads `.git` read-only without a range, opt-in
  ([`DC-FA-TRK-001`](spec/lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in))
- `targets` — declaration consistency between docs and build targets: a `make X`
  claimed in a doc **table row** without a Makefile rule (`gate-phantom`), or a
  Makefile rule without an entry in the authority doc (`gate-undocumented`);
  **hermetic** (no git, no Makefile execution), fail-closed, opt-in
  ([`DC-FA-TGT-001`](spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in))
- `citations` — verbatim quote verification: the directive
  `<!-- d-check:cite <path>:<from>-<to> -->` marks the following quote (a `>`-blockquote or
  inline `„…"`/`"…"`; blank lines in between are harmless, a fenced block in
  between separates); the whitespace-normalized quote must be a contiguous substring
  of the source span (`citation-mismatch`); the directive counts only **outside**
  fenced blocks **and outside inline code** — in backticks it is a mention, not a
  directive —, while the quoted source span is read **raw**, opt-in
  ([`DC-FA-CITE-001`](spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in))
- `structure` — structure invariants **within** a document: each rule defines a
  document class via **its own** globs, a section (literal or RE2) and up to
  **ten** conditions, each with its own reason code — non-empty (`section-empty`),
  minimum sentences (`section-thin`), task ceiling (`section-oversized`),
  **open** task items on the **raw** lines (`max-open-tasks` ⇒
  `section-tasks-open`, one finding per box on **its** line — unlike
  `max-tasks` it is immune to the paragraph-wide inline-code pairing, while
  staying fence-true),
  forbidden and required patterns (`section-forbidden`,
  `section-pattern-missing`), required markers (`section-marker-missing`),
  chronological monotonicity of the key column (`section-unordered`,
  `section-cell-untyped`), the shape of **every** heading in the section
  (`headings-match`/`headings-level` ⇒ `section-heading-mismatch`, reported per
  heading on its own line) and the **cell length** of a column addressed by its
  header **name** (`table.column` ⇒
  `section-cell-oversized`/`section-cell-undersized` on **its own** line,
  `section-column-missing` for a column that cannot be addressed);
  a missing section — or a rule matching no file — yields `section-missing`,
  several matches under `sections: one` yield `section-ambiguous`. Every rule may
  **declare its base set** — `tasks-ignore-pattern` removes task items from the
  `max-tasks` count (matched against the **item text behind the checkbox**, so
  that `^` means „the item starts like this"), `exempt-section-pattern` removes
  **sections** from the rule (matched against the **same** raw heading line as
  `section-pattern`, `#` sequence included). Both only **shrink** the checked
  set, carry no reason code of their own, and without them the finding set is
  byte-identical; if the section exemption empties the set, `section-missing`
  is reported rather than silent green — **unless the rule states how many it
  exempts:** `exempt-expect-count` (int ≥ 0, only alongside the pattern) keeps
  the declared empty set silent as long as the number holds, and otherwise
  reports `section-exempt-mismatch` — in **both** directions and
  **independently** of whether sections remain. The configuration defect (the
  selector matches nothing) stays `section-missing`, even with the number set.
  The number applies **per file**, its finding **aborts that file**, and it ages
  like any other authored text. And `hint` lets a rule **write its own
  explanation** — what the reader should do is known only to the rule; the
  reason code states the **kind** of defect. **Hermetic**,
  opt-in; the closure-note structure of module `planning` is a **preset** of the
  same semantics
  ([`DC-FA-STRUCT-001`](spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in))
- `workflows` — declaration consistency of the `uses:` references of CI
  workflows below a **configured** directory (`workflows.dir` — not hard-wired,
  because the location is CI-system specific), opt-in. A **foreign** reference
  names a full 40-character commit SHA (`uses-pin-missing`) with a tag comment
  behind it (`uses-pin-untagged`) — a tag can be moved, a SHA cannot; the
  **form** is checked, not the validity (that would be network). If the same
  SHA carries more than one tag comment within the scan set — across files —
  every affected line is reported (`uses-pin-tag-conflict`); repetition stays
  clean, and which comment is correct remains network. A **local**
  reference (`./…`) needs no pin, but two other guarantees: the target exists
  (`uses-local-missing`), and the calling **job** carries the permissions the
  target demands (`uses-local-perms-undeclared`, `uses-local-perms-narrow`) — a
  job without its own `permissions:` inherits the workflow header but cannot
  pass on what it does not declare. Unparsable YAML is reported
  (`workflow-unparsable`) rather than skipped. **Hermetic** (no git, no network,
  no execution); references come from the **YAML tree**, not from a text search
  ([`DC-FA-WF-001`](spec/lastenheft.md#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in))

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

The principle: **report, never repair.** d-check is a purely read-only
tool ([`DC-QA-03`](spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit));
deterministic findings are fixed, not suppressed. An opt-out
marker exists only where a non-existent target or an illustrative identifier
can be documented intent (`d-check:ignore`, per line) — it silences four
modules: `codepaths`, `ids`, `versions` and `diagrams`. That is a **named
list**, not a derivable criterion: `matrix`, `structure` and `citations` also
report on lines and do not know the marker. At `codepaths` and `ids` the marker counts only **outside inline code** — in backticks it is a mention, and it must sit inside an HTML comment; at `versions` and `diagrams` it stays a raw **token** — structurally for `diagrams`, as a named boundary for `versions` (it does read prose lines)
([`DC-FA-CODE-001`](spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in),
[`DC-FA-ID-001`](spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids),
[`DC-FA-VER-001`](spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
[`DC-FA-DIAG-001`](spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in));
all others do not know it.

## What makes it trustworthy?

**A checking tool whose own statements wobble is worthless** —
determinism and side-effect freedom are therefore core contracts of the
spec, and both are measured, not asserted:

- **Determinism:** identical input ⇒ byte-identical output,
  stably sorted; tested over ten repeated runs with
  hash comparison
  ([`DC-QA-02`](spec/lastenheft.md#dc-qa-02--determinismus)).
- **Side-effect-free and network-free:** never writes into the checked
  repository, opens no network connections except in the opt-in modules
  `external` and `sources` — measured in the gate run with a read-only mount
  and `--network none`
  ([`DC-QA-03`](spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **No silent defaults:** every invalid `.d-check.yml` aborts
  with exit 2; checking never proceeds with guessed configuration
  ([`DC-FA-CONF-001`](spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)).
- **Dogfooding:** d-check validates its own docs on every
  gate run — with the [self-configuration](.d-check.yml) fully
  built out (**eleven** modules incl. reference matrix, span artifacts,
  host-path hygiene, version-pin consistency, structure invariants,
  diagram identifiers and verbatim citations).
- **Container native-identical:** the image's finding output and exit code
  are byte-identical to native execution, tested automatically
  ([`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image));
  CI consumption via digest pin.

## Usage

Distributed as a container image via GHCR
([`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image)); **since `v0.67.0`**
also mirrored to Docker Hub as `pt9912/d-check` —
the same image, not a second build, same **config** digest (the **manifest**
digest is registry-local: when pinning by digest, use the one from the registry
you pull from)
([`DC-FA-DIST-002`](spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel)):

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.72.0
```

CI pipelines pin to the digest from the release notes rather than to
moving tags. The task-oriented entry point — from the first check to every
single module — is the [user handbook](docs/user/benutzerhandbuch.md) (German);
the terse invocation reference with options and exit codes is
[`docs/user/operations.md`](docs/user/operations.md), the release and
digest-pin path [`docs/user/releasing.md`](docs/user/releasing.md).

## Configuration (`.d-check.yml`)

Optional in the repo root; without a file the default modules
`links` + `anchors` run over `docs/`, `spec/` and the root `*.md`:

```yaml
scan:
  roots: ["."]                  # entire repo root
modules: [links, anchors, ids]  # external + sources stay strictly opt-in (network)
ids:
  patterns:
    - regex: 'ADR-\d{4}'
      target: docs/adr/         # identifiers must link here
```

The full schema with all keys, defaults and
validation constraints is in the
[specification §`.d-check.yml`](spec/spezifikation.md#spec-005--d-checkyml);
a living example fully built out (incl. reference matrix) is
[this repo's self-configuration](.d-check.yml). Every invalid
configuration aborts with exit 2 — checking never proceeds with
silent defaults
([`DC-FA-CONF-001`](spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)).

## Getting started

| Document | Contents |
|---|---|
| [`docs/user/benutzerhandbuch.md`](docs/user/benutzerhandbuch.md) | **User handbook** (German): task-oriented, every rule module with example configuration |
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
