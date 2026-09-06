# `make trace-check` — hält, dass jede Commit-Botschaft eine Traceability-Kennung nennt

## Vertrag

Via Modul `commits` (Image, dogfood). Jede Commit-Message nennt eine
`DC-*`/`ADR-*`/`MR-*`/`slice-*`-Kennung, sonst `commit-untraceable`.
Ausgenommen sind Merge- und Revert-Commits.

**Zwei Modi über dasselbe Modul:** `--range` für CI und Push, `--commit-msg -`
für den Hook über stdin.

## Grenze — was das Grün nicht abdeckt

1. **Geprüft ist die Nennung, nicht der Bezug** — dass eine Botschaft eine
   Kennung trägt, sagt nicht, dass der Commit sie betrifft. Permanent; das ist
   Review-Territorium.
2. **Der lokale Hook ist opt-in pro Klon** (`make hooks`); `--no-verify`
   umgeht ihn, nicht die PR-/Push-CI.

**Dependabot braucht dafür keine Ausnahme:** Seine Botschaften tragen die
Kennung im Präfix und erfüllen die Regel wie jeder andere Commit
([ADR-0067](../../docs/plan/adr/0067-dependabot-als-hebender-kanal.md)) — eine
erweiterte `exempt-pattern` hätte den Gate für eine **ganze Commit-Klasse**
blind gemacht.

## Bindung

**nicht** Teil von `gates`/`ci` — Commit-Zeit-Bindepunkt; gerufen vom
`commit-msg`-Hook und der PR-/Push-CI.
[ADR-0013](../../docs/plan/adr/0013-pr-ci-und-traceability-gate.md) ·
[ADR-0027](../../docs/plan/adr/0027-commits-traceability-modul.md) ·
[`DC-FA-COMMITS-001`](../../spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in)
