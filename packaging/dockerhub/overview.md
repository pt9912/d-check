# d-check

A deterministic, network-free reference gate for Markdown repositories.
d-check checks whether a documentation's references hold: local links and
images, heading anchors, identifier link obligations, permitted reference
directions, paths in inline code, verbatim quotes against their source — and
reports what does not hold. **It repairs nothing and never writes into the
checked repository.**

## Usage

```bash
docker run --rm -v "$PWD:/repo:ro" pt9912/d-check:__VERSION__
```

The mounted directory is checked as `/repo`; CLI options are appended. The
process runs as non-root, and a **read-only** mount is sufficient.

```bash
docker run --rm -v "$PWD:/repo:ro" pt9912/d-check:__VERSION__ --enable ids --disable external
```

## Exit codes

| Code | Meaning |
|---|---|
| `0` | no findings |
| `1` | findings reported |
| `2` | usage or configuration error (including: nothing mounted at `/repo`) |

## Reproducible runs

For CI, pin to the **digest** rather than to a moving tag. This image is a
**mirror** of `ghcr.io/pt9912/d-check` — the same image, not a second build:
the **config** digest is identical on both registries, and the release pipeline
verifies it after the push.

**The manifest digest, in contrast, is registry-local** — it depends on each
registry's blob compression. So take the digest **from the registry you pull
from**; one copied from GHCR will not resolve here.

```bash
docker run --rm -v "$PWD:/repo:ro" pt9912/d-check@sha256:<digest>
```

The digest for every version is listed in the
[release notes](https://github.com/pt9912/d-check/releases).

`:latest` moves for stable releases only; pre-releases never receive it.

## Modules

21 rule modules, most of them opt-in: `links`, `anchors`, `ids`, `matrix`,
`external`, `sources`, `codepaths`, `spans`, `hostpaths`, `diagrams`,
`versions`, `pins`, `immutable`, `vcs`, `commits`, `planning`, `tracked`,
`targets`, `citations`, `structure`, `workflows`. Two of them reach the network
(`external`, `sources`) and are never active by default.

## Documentation

- [README](https://github.com/pt9912/d-check#readme) — overview
- [User handbook](https://github.com/pt9912/d-check/blob/main/docs/user/benutzerhandbuch.md) — task-oriented (German)
- [Invocation reference](https://github.com/pt9912/d-check/blob/main/docs/user/operations.md) — options and exit codes

Source and issues: <https://github.com/pt9912/d-check> · Licence: MIT
