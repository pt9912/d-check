# `make hooks` — installiert die lokalen git-Hooks, die Commit und Übergang an grüne Gates binden

## Vertrag

Setzt `core.hooksPath` auf [`.githooks/`](../../.githooks/) und bindet damit
Commit und Slice-Übergang an grüne Prüfungen:
[`make trace-check`](trace-check.md) an jede Commit-Botschaft,
[`make adr-check`](adr-check.md) und [`make doc-check`](doc-check.md) an
jeden Commit, [`make verify-closure-notes`](verify-closure-notes.md) an jeden
Übergang eines Slice nach `done/`. Welcher Hook welche Prüfung ruft und
welcher Diff den Übergang erkennt, legt
[`SPEC-095`](../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)
fest — für einen Slice direkt unter `done/` hängen die Vorbedingungen damit am
**Übergang** selbst, nicht nur an
einer gelegentlichen `fullbuild`-Prüfung; für die Unterverzeichnisse gilt das
nicht (Grenze 4).

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
4. **Der Übergang löst `verify-closure-notes` aus, aber der Lauf prüft nur
   `done/` selbst.** Ein Slice, der nach `done/wellenlos/` oder unter ein
   Wellen-Verzeichnis wandert, wird erkannt und nicht geprüft; ein grüner
   Commit sagt über seine Closure-Notiz nichts.

## Bindung

Kein Gate über den Repo-Zustand, sondern Einrichtung: das Target **installiert**
die Bindung, es prüft nichts.
[ADR-0013](../../docs/plan/adr/0013-pr-ci-und-traceability-gate.md) ·
[ADR-0016](../../docs/plan/adr/0016-adr-immutable-gate.md) ·
[ADR-0024](../../docs/plan/adr/0024-vcs-immutable-gate.md)
