# `make hooks` — installiert die lokalen git-Hooks, die Commit und Übergang an grüne Gates binden

## Vertrag

Setzt `core.hooksPath` auf [`.githooks/`](../../.githooks/) und aktiviert damit
drei Bindungen:

- **`commit-msg`** — Traceability: jede Commit-Botschaft nennt eine
  `DC-*`/`ADR-*`/`MR-*`/`slice-*`-Kennung (`make trace-check`).
- **`pre-commit`, drei Teile** — ADR-Immutabilität über das Modul `vcs`
  ([ADR-0024](../../docs/plan/adr/0024-vcs-immutable-gate.md)) **und** der
  volle `make doc-check` als Doku-Guard (seit welle-79): der Commit ist an
  einen grünen Doku-Stand gebunden, und ein roter Gate-Exit kann keine
  Shell-Kette mehr passieren.
  **Dritter `pre-commit`-Teil, der Slice-Closure-Übergangs-Wächter** (seit
  welle-86): ein gestagter
  Rename/Add nach `docs/plan/planning/done/slice-*.md` — **nicht rekursiv**,
  ein archivierter Stub eine Ebene tiefer zählt nicht — löst zusätzlich
  [`make verify-closure-notes`](verify-closure-notes.md) aus. Die
  Vorbedingungen hängen damit am **Übergang** selbst, nicht nur an einer
  gelegentlichen `fullbuild`-Prüfung.

## Grenze — was das Grün nicht abdeckt

1. **Opt-in pro Klon** — `core.hooksPath` ist lokale git-Konfiguration; aus
   einem fremden Klon sind die Hooks nicht erzwingbar. Permanent. Der
   klon-unabhängige Boden ist die PR-/Push-CI, die dieselbe Bindung über die
   Commit-Range fährt (`git diff --diff-filter=AR "$RANGE"`) — `--no-verify`
   umgeht den lokalen Hook, nicht die CI.
2. **Der `pre-commit`-Hook prüft den Arbeitsbaum, nicht den Commit-Stand.**
   Ein Commit, der weniger enthält als der Arbeitsbaum — etwa weil ein
   `git add` scheiterte —, kann grün passieren und trotzdem einen roten
   Zwischenstand hinterlassen. Gemessen in slice-202. **Sichtbar** nur durch
   Gegenlesen (`git show HEAD:<datei>`) — geheilt wird per `--amend`, nicht
   durch den Hook.
3. **Die CI blockiert einen Merge nur mit Branch Protection** — ein
   Pflicht-Status-Check auf dem Default-Branch liegt **außerhalb** des Repos
   und ist aus dem Klon nicht auditierbar. Ohne sie ist die CI *advisory*.

## Bindung

Kein Gate über den Repo-Zustand, sondern Einrichtung: das Target **installiert**
die Bindung, es prüft nichts.
[ADR-0013](../../docs/plan/adr/0013-pr-ci-und-traceability-gate.md) ·
[ADR-0016](../../docs/plan/adr/0016-adr-immutable-gate.md) ·
[ADR-0024](../../docs/plan/adr/0024-vcs-immutable-gate.md)
