# d-check

Ein deterministisches, netzloses Referenz-Gate für Markdown-Repositories.
d-check prüft, ob die Verweise in einer Dokumentation halten: lokale Links und
Bilder, Heading-Anker, Kennungs-Linkpflicht, erlaubte Referenzrichtungen,
Pfade in Inline-Code, wortgleiche Zitate gegen ihre Quelle — und meldet, was
nicht hält. **Es repariert nichts und schreibt nie in das geprüfte Repository.**

## Aufruf

```bash
docker run --rm -v "$PWD:/repo:ro" pt9912/d-check:__VERSION__
```

Das gemountete Verzeichnis wird als `/repo` geprüft; CLI-Optionen werden
angehängt. Der Prozess läuft als Nicht-root, und ein **read-only**-Mount
genügt.

```bash
docker run --rm -v "$PWD:/repo:ro" pt9912/d-check:__VERSION__ --enable ids --disable external
```

## Exit-Codes

| Code | Bedeutung |
|---|---|
| `0` | keine Befunde |
| `1` | Befunde gefunden |
| `2` | Nutzungs-/Konfigurationsfehler (auch: kein Mount auf `/repo`) |

## Reproduzierbare Läufe

Für CI pinnen Sie auf den **Digest** statt auf einen beweglichen Tag. Dieses
Image ist ein **Spiegel** von `ghcr.io/pt9912/d-check` — dasselbe Bild, kein
zweiter Bau: der **Config**-Digest ist auf beiden Registries gleich, und die
Release-Pipeline misst das nach dem Push.

**Der Manifest-Digest ist dagegen registry-lokal** — er hängt an der
Blob-Kompression der jeweiligen Registry. Nehmen Sie deshalb den Digest **der
Registry, aus der Sie ziehen**; ein von GHCR kopierter Digest löst hier nicht
auf.

```bash
docker run --rm -v "$PWD:/repo:ro" pt9912/d-check@sha256:<digest>
```

Der Digest jeder Version steht in den
[Release-Notes](https://github.com/pt9912/d-check/releases).

`:latest` bewegt sich ausschließlich für stabile Releases; Vorabversionen
erhalten es nicht.

## Module

20 Regelmodule, die meisten opt-in: `links`, `anchors`, `ids`, `matrix`,
`external`, `sources`, `codepaths`, `spans`, `hostpaths`, `diagrams`,
`versions`, `pins`, `immutable`, `vcs`, `commits`, `planning`, `tracked`,
`targets`, `citations`, `structure`. Zwei davon gehen ins Netz (`external`,
`sources`) und sind nie im Default aktiv.

## Dokumentation

- [README](https://github.com/pt9912/d-check#readme) — Überblick
- [Benutzerhandbuch](https://github.com/pt9912/d-check/blob/main/docs/user/benutzerhandbuch.md) — aufgabenorientiert (Deutsch)
- [Aufruf-Referenz](https://github.com/pt9912/d-check/blob/main/docs/user/operations.md) — Optionen und Exit-Codes

Quelle und Issues: <https://github.com/pt9912/d-check> · Lizenz: MIT
